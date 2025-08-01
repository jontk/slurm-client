// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package retry

import (
	"context"
	"math"
	"math/rand"
	"net/http"
	"time"
)

// Policy defines the interface for retry policies
type Policy interface {
	// ShouldRetry determines if a request should be retried
	ShouldRetry(ctx context.Context, resp *http.Response, err error, attempt int) bool

	// WaitTime returns the wait time before the next retry
	WaitTime(attempt int) time.Duration

	// MaxRetries returns the maximum number of retries
	MaxRetries() int
}

// HTTPExponentialBackoff implements exponential backoff retry policy for HTTP requests
type HTTPExponentialBackoff struct {
	maxRetries    int
	minWaitTime   time.Duration
	maxWaitTime   time.Duration
	backoffFactor float64
	jitter        bool
}

// NewHTTPExponentialBackoff creates a new exponential backoff retry policy for HTTP requests
func NewHTTPExponentialBackoff() *HTTPExponentialBackoff {
	return &HTTPExponentialBackoff{
		maxRetries:    3,
		minWaitTime:   1 * time.Second,
		maxWaitTime:   30 * time.Second,
		backoffFactor: 2.0,
		jitter:        true,
	}
}

// WithMaxRetries sets the maximum number of retries
func (e *HTTPExponentialBackoff) WithMaxRetries(maxRetries int) *HTTPExponentialBackoff {
	e.maxRetries = maxRetries
	return e
}

// WithMinWaitTime sets the minimum wait time
func (e *HTTPExponentialBackoff) WithMinWaitTime(minWaitTime time.Duration) *HTTPExponentialBackoff {
	e.minWaitTime = minWaitTime
	return e
}

// WithMaxWaitTime sets the maximum wait time
func (e *HTTPExponentialBackoff) WithMaxWaitTime(maxWaitTime time.Duration) *HTTPExponentialBackoff {
	e.maxWaitTime = maxWaitTime
	return e
}

// WithBackoffFactor sets the backoff factor
func (e *HTTPExponentialBackoff) WithBackoffFactor(backoffFactor float64) *HTTPExponentialBackoff {
	e.backoffFactor = backoffFactor
	return e
}

// WithJitter enables or disables jitter
func (e *HTTPExponentialBackoff) WithJitter(jitter bool) *HTTPExponentialBackoff {
	e.jitter = jitter
	return e
}

// ShouldRetry determines if a request should be retried
func (e *HTTPExponentialBackoff) ShouldRetry(ctx context.Context, resp *http.Response, err error, attempt int) bool {
	if attempt >= e.maxRetries {
		return false
	}

	// Check if context is cancelled
	select {
	case <-ctx.Done():
		return false
	default:
	}

	// Retry on network errors
	if err != nil {
		return true
	}

	// Retry on specific HTTP status codes
	if resp != nil {
		switch resp.StatusCode {
		case http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout:
			return true
		}
	}

	return false
}

// WaitTime returns the wait time before the next retry
func (e *HTTPExponentialBackoff) WaitTime(attempt int) time.Duration {
	if attempt <= 0 {
		return e.minWaitTime
	}

	// Calculate exponential backoff
	waitTime := time.Duration(float64(e.minWaitTime) * math.Pow(e.backoffFactor, float64(attempt-1)))

	// Apply maximum wait time
	if waitTime > e.maxWaitTime {
		waitTime = e.maxWaitTime
	}

	// Apply jitter if enabled
	if e.jitter {
		jitterAmount := time.Duration(rand.Float64() * float64(waitTime) * 0.1)
		waitTime += jitterAmount
	}

	return waitTime
}

// MaxRetries returns the maximum number of retries
func (e *HTTPExponentialBackoff) MaxRetries() int {
	return e.maxRetries
}

// FixedDelay implements fixed delay retry policy
type FixedDelay struct {
	maxRetries int
	delay      time.Duration
}

// NewFixedDelay creates a new fixed delay retry policy
func NewFixedDelay(maxRetries int, delay time.Duration) *FixedDelay {
	return &FixedDelay{
		maxRetries: maxRetries,
		delay:      delay,
	}
}

// ShouldRetry determines if a request should be retried
func (f *FixedDelay) ShouldRetry(ctx context.Context, resp *http.Response, err error, attempt int) bool {
	if attempt >= f.maxRetries {
		return false
	}

	// Check if context is cancelled
	select {
	case <-ctx.Done():
		return false
	default:
	}

	// Retry on network errors
	if err != nil {
		return true
	}

	// Retry on specific HTTP status codes
	if resp != nil {
		switch resp.StatusCode {
		case http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout:
			return true
		}
	}

	return false
}

// WaitTime returns the wait time before the next retry
func (f *FixedDelay) WaitTime(attempt int) time.Duration {
	return f.delay
}

// MaxRetries returns the maximum number of retries
func (f *FixedDelay) MaxRetries() int {
	return f.maxRetries
}

// NoRetry implements no retry policy
type NoRetry struct{}

// NewNoRetry creates a new no retry policy
func NewNoRetry() *NoRetry {
	return &NoRetry{}
}

// ShouldRetry always returns false
func (n *NoRetry) ShouldRetry(ctx context.Context, resp *http.Response, err error, attempt int) bool {
	return false
}

// WaitTime returns zero duration
func (n *NoRetry) WaitTime(attempt int) time.Duration {
	return 0
}

// MaxRetries returns zero
func (n *NoRetry) MaxRetries() int {
	return 0
}
