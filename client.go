package xxyy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://www.xxyy.io"
	defaultTimeout = 30 * time.Second
	apiPrefix      = "/api/trade/open/api"
)

// Client is the XXYY API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

// NewClient creates a new XXYY API client.
//
// apiKey is required — get one at https://www.xxyy.io/apikey
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: http.DefaultClient,
		timeout:    defaultTimeout,
	}
	for _, opt := range opts {
		opt(c)
	}
	// Trim trailing slashes from base URL
	c.baseURL = strings.TrimRight(c.baseURL, "/")
	return c
}

// apiResponse is used internally for JSON unmarshalling to check structure.
type apiResponseRaw struct {
	Code    int             `json:"code"`
	Msg     string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
	Success bool            `json:"success"`
}

// doRequest performs an HTTP request against the XXYY API.
// It handles error codes, rate-limit retries, and response validation.
func doRequest[T any](
	ctx context.Context,
	c *Client,
	method, path string,
	params map[string]string,
	body any,
	retries int,
) (*ApiResponse[T], error) {
	// Build URL
	u := c.baseURL + path
	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			q.Set(k, v)
		}
		u += "?" + q.Encode()
	}

	// Build request body
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("xxyy: failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	// Create request with timeout context
	reqCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, method, u, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("xxyy: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("xxyy: request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, newXxyyError(resp.StatusCode, formatHTTPError(resp.StatusCode, resp.Status))
	}

	// Read and parse response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("xxyy: failed to read response body: %w", err)
	}

	var raw apiResponseRaw
	if err := json.Unmarshal(respBody, &raw); err != nil {
		return nil, newXxyyError(0, "unexpected API response format — the XXYY API may have changed")
	}

	// Handle API error codes
	if raw.Code == ErrCodeAPIKeyInvalid || raw.Code == ErrCodeAPIKeyDisabled {
		return nil, newXxyyError(raw.Code, fmt.Sprintf("API Key error: %s", raw.Msg))
	}

	if raw.Code == ErrCodeRateLimited {
		if retries > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(2 * time.Second):
			}
			return doRequest[T](ctx, c, method, path, params, body, retries-1)
		}
		return nil, newXxyyError(ErrCodeRateLimited, "rate limited — please try again later")
	}

	if raw.Code == ErrCodeServerError {
		return nil, newXxyyError(ErrCodeServerError, fmt.Sprintf("server error: %s", raw.Msg))
	}

	// Parse data field
	var data T
	if len(raw.Data) > 0 && string(raw.Data) != "null" {
		if err := json.Unmarshal(raw.Data, &data); err != nil {
			return nil, fmt.Errorf("xxyy: failed to parse response data: %w", err)
		}
	}

	return &ApiResponse[T]{
		Code:    raw.Code,
		Msg:     raw.Msg,
		Data:    data,
		Success: raw.Success,
	}, nil
}

// doGet performs a GET request with automatic rate-limit retry.
func doGet[T any](ctx context.Context, c *Client, path string, params map[string]string) (*ApiResponse[T], error) {
	return doRequest[T](ctx, c, http.MethodGet, path, params, nil, 2)
}

// doPost performs a POST request with automatic rate-limit retry.
func doPost[T any](ctx context.Context, c *Client, path string, body any, params map[string]string) (*ApiResponse[T], error) {
	return doRequest[T](ctx, c, http.MethodPost, path, params, body, 2)
}

// doPostNoRetry performs a POST request without rate-limit retry.
// Used for irreversible operations like swap.
func doPostNoRetry[T any](ctx context.Context, c *Client, path string, body any, params map[string]string) (*ApiResponse[T], error) {
	return doRequest[T](ctx, c, http.MethodPost, path, params, body, 0)
}
