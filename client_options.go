package xxyy

import (
	"net/http"
	"time"
)

// Option configures the Client.
type Option func(*Client)

// WithBaseURL sets a custom API base URL. Default: "https://www.xxyy.io".
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient sets a custom *http.Client for requests.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithTimeout sets the default request timeout. Default: 30s.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.timeout = d
	}
}
