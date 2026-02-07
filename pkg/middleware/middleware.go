// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Package middleware provides HTTP middleware for the SLURM client
package middleware

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/jontk/slurm-client/pkg/logging"
	"github.com/jontk/slurm-client/pkg/retry"
)

// Middleware is a function that wraps an http.RoundTripper
type Middleware func(http.RoundTripper) http.RoundTripper

// Chain creates a single middleware from a chain of middlewares
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// RoundTripperFunc is an adapter to allow functions to be used as RoundTrippers
type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// ContextKeyRequestID is the key for storing request IDs in context
	ContextKeyRequestID ContextKey = "request_id"
)

// For backwards compatibility, keep private alias
const contextKeyRequestID = ContextKeyRequestID

// WithTimeout adds timeout handling to requests
func WithTimeout(timeout time.Duration) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			ctx := req.Context()

			// Only add timeout if context doesn't already have a deadline
			if _, hasDeadline := ctx.Deadline(); !hasDeadline && timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, timeout)
				defer cancel()
				req = req.WithContext(ctx)
			}

			return next.RoundTrip(req)
		})
	}
}

// TimeoutConfig holds operation-specific timeout configuration
type TimeoutConfig struct {
	Default time.Duration // Default timeout for all operations
	Read    time.Duration // Timeout for GET requests
	Write   time.Duration // Timeout for POST, PUT, DELETE requests
	List    time.Duration // Timeout for list operations (not distinguishable at middleware level, uses Read)
	Watch   time.Duration // Timeout for watch operations (long-polling)
}

// WithTimeoutConfig adds operation-specific timeouts based on HTTP method
func WithTimeoutConfig(config *TimeoutConfig) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			ctx := req.Context()

			// Only add timeout if context doesn't already have a deadline
			if _, hasDeadline := ctx.Deadline(); !hasDeadline {
				timeout := selectTimeout(config, req.Method)
				if timeout > 0 {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, timeout)
					defer cancel()
					req = req.WithContext(ctx)
				}
			}

			return next.RoundTrip(req)
		})
	}
}

// selectTimeout returns the appropriate timeout based on HTTP method
func selectTimeout(config *TimeoutConfig, method string) time.Duration {
	if config == nil {
		return 0
	}

	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		if config.Read > 0 {
			return config.Read
		}
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		if config.Write > 0 {
			return config.Write
		}
	}

	return config.Default
}

// WithLogging adds structured logging to requests
func WithLogging(logger logging.Logger) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			start := time.Now()

			// Log request
			reqLogger := logging.LogAPICall(logger, req.Method, req.URL.Path,
				"host", req.URL.Host,
				"content_length", req.ContentLength,
			)

			reqLogger.Debug("sending request")

			// Execute request
			resp, err := next.RoundTrip(req)

			// Log response
			duration := time.Since(start)
			if err != nil {
				logging.LogError(reqLogger, err, "request_failed",
					"duration_ms", duration.Milliseconds(),
				)
				return nil, err
			}

			reqLogger.Info("request completed",
				"status_code", resp.StatusCode,
				"duration_ms", duration.Milliseconds(),
				"content_length", resp.ContentLength,
			)

			return resp, nil
		})
	}
}

// WithRetry adds retry logic with exponential backoff
func WithRetry(maxAttempts int, shouldRetry ShouldRetryFunc) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			var lastErr error
			var lastResp *http.Response

			for attempt := range maxAttempts {
			// Clone request for retry
				reqCopy := cloneRequest(req)

				resp, err := next.RoundTrip(reqCopy)

				// Check if we should retry
				if !shouldRetry(resp, err, attempt) {
					return resp, err
				}

				// Close response body if present
				if resp != nil && resp.Body != nil {
					_, _ = io.Copy(io.Discard, resp.Body) // Intentionally ignore error during cleanup
					_ = resp.Body.Close()                 // Intentionally ignore error during cleanup
				}

				lastErr = err
				lastResp = resp

				// Calculate backoff
				if attempt < maxAttempts-1 {
					backoff := calculateBackoff(attempt)
					select {
					case <-time.After(backoff):
						// Continue to next attempt
					case <-req.Context().Done():
						return nil, req.Context().Err()
					}
				}
			}

			// Return last response/error
			if lastErr != nil {
				return nil, fmt.Errorf("all %d attempts failed: %w", maxAttempts, lastErr)
			}
			return lastResp, nil
		})
	}
}

// ShouldRetryFunc determines if a request should be retried
type ShouldRetryFunc func(resp *http.Response, err error, attempt int) bool

// DefaultShouldRetry is the default retry logic
func DefaultShouldRetry(resp *http.Response, err error, _ int) bool {
	// Don't retry if context is canceled
	if err != nil && errors.Is(err, context.Canceled) {
		return false
	}

	// Retry on network errors
	if err != nil {
		return true
	}

	// Retry on 5xx errors
	if resp != nil && resp.StatusCode >= 500 {
		return true
	}

	// Retry on 429 (Too Many Requests)
	if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
		return true
	}

	return false
}

// calculateBackoff calculates exponential backoff with jitter
func calculateBackoff(attempt int) time.Duration {
	// Cap attempt to prevent integer overflow (max ~8.5 minutes at attempt 9)
	if attempt < 0 {
		attempt = 0
	}
	if attempt > 9 {
		attempt = 9
	}
	// #nosec G115 -- attempt is bounded to [0,9] range, preventing overflow
	base := time.Duration(1<<uint(attempt)) * time.Second
	jitter := time.Duration(float64(base) * 0.1)
	return base + jitter
}

// WithRetryPolicy adds retry logic using a custom retry.Policy for backoff configuration
func WithRetryPolicy(policy retry.Policy) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			var lastErr error
			var lastResp *http.Response

			maxAttempts := policy.MaxRetries() + 1 // MaxRetries is number of retries, not total attempts
			for attempt := range maxAttempts {
				// Clone request for retry
				reqCopy := cloneRequest(req)

				resp, err := next.RoundTrip(reqCopy)

				// Check if we should retry using the policy
				if !policy.ShouldRetry(req.Context(), resp, err, attempt) {
					return resp, err
				}

				// Close response body if present
				if resp != nil && resp.Body != nil {
					_, _ = io.Copy(io.Discard, resp.Body)
					_ = resp.Body.Close()
				}

				lastErr = err
				lastResp = resp

				// Use policy's wait time for backoff
				if attempt < maxAttempts-1 {
					waitTime := policy.WaitTime(attempt)
					select {
					case <-time.After(waitTime):
						// Continue to next attempt
					case <-req.Context().Done():
						return nil, req.Context().Err()
					}
				}
			}

			// Return last response/error
			if lastErr != nil {
				return nil, fmt.Errorf("all %d attempts failed: %w", maxAttempts, lastErr)
			}
			return lastResp, nil
		})
	}
}

// WithHeaders adds custom headers to requests
func WithHeaders(headers map[string]string) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			// Clone request to avoid modifying the original
			req = cloneRequest(req)

			// Add headers
			for key, value := range headers {
				req.Header.Set(key, value)
			}

			return next.RoundTrip(req)
		})
	}
}

// WithUserAgent sets a custom User-Agent header
func WithUserAgent(userAgent string) Middleware {
	return WithHeaders(map[string]string{
		"User-Agent": userAgent,
	})
}

// WithRequestID adds a unique request ID to each request
func WithRequestID(generator func() string) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			// Generate request ID
			requestID := generator()

			// Clone request and add header
			req = cloneRequest(req)
			req.Header.Set("X-Request-ID", requestID)

			// Add to context for logging
			ctx := context.WithValue(req.Context(), contextKeyRequestID, requestID)
			req = req.WithContext(ctx)

			return next.RoundTrip(req)
		})
	}
}

// WithMetrics adds metrics collection to requests
func WithMetrics(collector MetricsCollector) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			start := time.Now()

			// Record request
			collector.RecordRequest(req.Method, req.URL.Path)

			// Execute request
			resp, err := next.RoundTrip(req)

			// Record response
			duration := time.Since(start)
			if err != nil {
				collector.RecordError(req.Method, req.URL.Path, err)
			} else {
				collector.RecordResponse(req.Method, req.URL.Path, resp.StatusCode, duration)
			}

			return resp, err
		})
	}
}

// MetricsCollector is the interface for collecting metrics
type MetricsCollector interface {
	RecordRequest(method, path string)
	RecordResponse(method, path string, statusCode int, duration time.Duration)
	RecordError(method, path string, err error)
}

// cloneRequest creates a shallow copy of a request
func cloneRequest(req *http.Request) *http.Request {
	// Clone the request
	r := req.Clone(req.Context())

	// Clone body if present
	if req.Body != nil {
		bodyBytes, _ := io.ReadAll(req.Body) // Best effort body cloning
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	return r
}

// WithCircuitBreaker adds circuit breaker functionality
func WithCircuitBreaker(threshold int, timeout time.Duration) Middleware {
	breaker := &circuitBreaker{
		threshold: threshold,
		timeout:   timeout,
		failures:  0,
		lastFail:  time.Time{},
	}

	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if !breaker.Allow() {
				return nil, fmt.Errorf("circuit breaker is open")
			}

			resp, err := next.RoundTrip(req)

			if err != nil || (resp != nil && resp.StatusCode >= 500) {
				breaker.RecordFailure()
			} else {
				breaker.RecordSuccess()
			}

			return resp, err
		})
	}
}

type circuitBreaker struct {
	threshold int
	timeout   time.Duration
	mu        sync.RWMutex
	failures  int
	lastFail  time.Time
}

func (cb *circuitBreaker) Allow() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	if cb.failures < cb.threshold {
		return true
	}

	// Check if timeout has passed
	return time.Since(cb.lastFail) > cb.timeout
}

func (cb *circuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFail = time.Now()
}

func (cb *circuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures = 0
}
