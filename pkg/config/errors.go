package config

import "errors"

var (
	// ErrMissingBaseURL is returned when the base URL is not set
	ErrMissingBaseURL = errors.New("base URL is required")
	
	// ErrInvalidTimeout is returned when the timeout is invalid
	ErrInvalidTimeout = errors.New("timeout must be greater than 0")
	
	// ErrInvalidMaxRetries is returned when max retries is invalid
	ErrInvalidMaxRetries = errors.New("max retries must be greater than or equal to 0")
)