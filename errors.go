package xxyy

import "fmt"

// Error codes returned by the XXYY API.
const (
	ErrCodeAPIKeyInvalid  = 8060
	ErrCodeAPIKeyDisabled = 8061
	ErrCodeRateLimited    = 8062
	ErrCodeServerError    = 300
)

// XxyyError represents an error returned by the XXYY API.
type XxyyError struct {
	Code    int    // API error code or HTTP status code
	Message string // Human-readable error message
}

func (e *XxyyError) Error() string {
	return fmt.Sprintf("xxyy: [%d] %s", e.Code, e.Message)
}

// IsAPIKeyError returns true if the error is an API key error (invalid or disabled).
func (e *XxyyError) IsAPIKeyError() bool {
	return e.Code == ErrCodeAPIKeyInvalid || e.Code == ErrCodeAPIKeyDisabled
}

// IsRateLimited returns true if the error is a rate limit error.
func (e *XxyyError) IsRateLimited() bool {
	return e.Code == ErrCodeRateLimited
}

// IsServerError returns true if the error is a server-side error.
func (e *XxyyError) IsServerError() bool {
	return e.Code == ErrCodeServerError
}

// newXxyyError creates a new XxyyError.
func newXxyyError(code int, msg string) *XxyyError {
	return &XxyyError{Code: code, Message: msg}
}

// formatHTTPError returns a human-readable error message for HTTP status codes.
func formatHTTPError(status int, statusText string) string {
	switch {
	case status == 401 || status == 403:
		return fmt.Sprintf("Authentication failed (HTTP %d). Check your XXYY_API_KEY.", status)
	case status == 408 || status == 504:
		return fmt.Sprintf("Request timed out (HTTP %d). The XXYY server may be slow — try again later.", status)
	case status >= 500:
		return fmt.Sprintf("XXYY server error (HTTP %d). Try again later.", status)
	default:
		return fmt.Sprintf("HTTP %d: %s", status, statusText)
	}
}
