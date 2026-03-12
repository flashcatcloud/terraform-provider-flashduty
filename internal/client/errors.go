package client

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrNotFound      = errors.New("resource not found")
	ErrUnauthorized  = errors.New("unauthorized: invalid or missing APP key")
	ErrForbidden     = errors.New("forbidden: insufficient permissions")
	ErrBadRequest    = errors.New("bad request: invalid parameters")
	ErrInternalError = errors.New("internal server error")
	ErrRateLimited   = errors.New("rate limited: too many requests")
	ErrConflict      = errors.New("conflict: resource already exists or was modified")
)

// IsNotFoundError checks if the error indicates a resource was not found.
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// httpStatusToError converts HTTP status codes to appropriate errors.
func httpStatusToError(statusCode int, body string) error {
	switch statusCode {
	case http.StatusNotFound:
		return fmt.Errorf("%w: %s", ErrNotFound, body)
	case http.StatusUnauthorized:
		return fmt.Errorf("%w: %s", ErrUnauthorized, body)
	case http.StatusForbidden:
		return fmt.Errorf("%w: %s", ErrForbidden, body)
	case http.StatusBadRequest:
		return fmt.Errorf("%w: %s", ErrBadRequest, body)
	case http.StatusTooManyRequests:
		return fmt.Errorf("%w: %s", ErrRateLimited, body)
	case http.StatusConflict:
		return fmt.Errorf("%w: %s", ErrConflict, body)
	default:
		if statusCode >= 500 {
			return fmt.Errorf("%w: status %d, %s", ErrInternalError, statusCode, body)
		}
		return fmt.Errorf("API request failed with status %d: %s", statusCode, body)
	}
}
