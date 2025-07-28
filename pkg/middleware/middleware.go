// Package middleware provides HTTP middleware for the SLURM client
package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jontk/slurm-client/pkg/logging"
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
			
			for attempt := 0; attempt < maxAttempts; attempt++ {
				// Clone request for retry
				reqCopy := cloneRequest(req)
				
				resp, err := next.RoundTrip(reqCopy)
				
				// Check if we should retry
				if !shouldRetry(resp, err, attempt) {
					return resp, err
				}
				
				// Close response body if present
				if resp != nil && resp.Body != nil {
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
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
func DefaultShouldRetry(resp *http.Response, err error, attempt int) bool {
	// Don't retry if context is canceled
	if err != nil && err == context.Canceled {
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
	if resp != nil && resp.StatusCode == 429 {
		return true
	}
	
	return false
}

// calculateBackoff calculates exponential backoff with jitter
func calculateBackoff(attempt int) time.Duration {
	base := time.Duration(1<<uint(attempt)) * time.Second
	jitter := time.Duration(float64(base) * 0.1)
	return base + jitter
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
			ctx := context.WithValue(req.Context(), "request_id", requestID)
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
		bodyBytes, _ := io.ReadAll(req.Body)
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
	failures  int
	lastFail  time.Time
}

func (cb *circuitBreaker) Allow() bool {
	if cb.failures < cb.threshold {
		return true
	}
	
	// Check if timeout has passed
	return time.Since(cb.lastFail) > cb.timeout
}

func (cb *circuitBreaker) RecordFailure() {
	cb.failures++
	cb.lastFail = time.Now()
}

func (cb *circuitBreaker) RecordSuccess() {
	cb.failures = 0
}