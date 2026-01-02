package geocoding

import (
	"errors"
	"fmt"
)

var (
	// ErrConcurrencyLimitExceeded is returned when the maximum number of concurrent requests is exceeded.
	ErrConcurrencyLimitExceeded = errors.New("concurrency limit exceeded")
	// ErrInvalidParameter is returned when an input parameter is invalid.
	ErrInvalidParameter = errors.New("invalid parameter")
)

// APIError represents an error returned by the Open Meteo API.
type APIError struct {
	Reason string `json:"reason"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api error: %s", e.Reason)
}
