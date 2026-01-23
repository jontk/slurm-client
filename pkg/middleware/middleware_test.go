// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jontk/slurm-client/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock round tripper for testing
type mockRoundTripper struct {
	responses []mockResponse
	calls     []http.Request
	mu        sync.Mutex
}

type mockResponse struct {
	response *http.Response
	err      error
}

func newMockRoundTripper() *mockRoundTripper {
	return &mockRoundTripper{
		responses: make([]mockResponse, 0),
		calls:     make([]http.Request, 0),
	}
}

func (m *mockRoundTripper) addResponse(resp *http.Response, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses = append(m.responses, mockResponse{response: resp, err: err})
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Store the call
	m.calls = append(m.calls, *req)

	if len(m.responses) == 0 {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	}

	response := m.responses[0]
	if len(m.responses) > 1 {
		m.responses = m.responses[1:]
	} else {
		// Keep the last response for subsequent calls
		// Don't remove it to avoid nil pointer issues
	}

	return response.response, response.err
}

func (m *mockRoundTripper) getCalls() []http.Request {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]http.Request{}, m.calls...)
}

// Mock metrics collector
type mockMetricsCollector struct {
	requests  []requestRecord
	responses []responseRecord
	errors    []errorRecord
	mu        sync.Mutex
}

type requestRecord struct {
	method string
	path   string
}

type responseRecord struct {
	method     string
	path       string
	statusCode int
	duration   time.Duration
}

type errorRecord struct {
	method string
	path   string
	err    error
}

func newMockMetricsCollector() *mockMetricsCollector {
	return &mockMetricsCollector{
		requests:  make([]requestRecord, 0),
		responses: make([]responseRecord, 0),
		errors:    make([]errorRecord, 0),
	}
}

func (m *mockMetricsCollector) RecordRequest(method, path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requests = append(m.requests, requestRecord{method: method, path: path})
}

func (m *mockMetricsCollector) RecordResponse(method, path string, statusCode int, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses = append(m.responses, responseRecord{
		method:     method,
		path:       path,
		statusCode: statusCode,
		duration:   duration,
	})
}

func (m *mockMetricsCollector) RecordError(method, path string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errors = append(m.errors, errorRecord{method: method, path: path, err: err})
}

func TestRoundTripperFunc(t *testing.T) {
	called := false
	fn := RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		called = true
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := fn.RoundTrip(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.True(t, called)
}

func TestChain(t *testing.T) {
	mock := newMockRoundTripper()

	// Create middleware that adds headers
	middleware1 := func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			req.Header.Set("X-Middleware-1", "true")
			return next.RoundTrip(req)
		})
	}

	middleware2 := func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			req.Header.Set("X-Middleware-2", "true")
			return next.RoundTrip(req)
		})
	}

	// Chain middlewares
	chained := Chain(middleware1, middleware2)
	roundTripper := chained(mock)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := roundTripper.RoundTrip(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.NoError(t, err)

	calls := mock.getCalls()
	require.Len(t, calls, 1)

	// Both middleware headers should be present
	assert.Equal(t, "true", calls[0].Header.Get("X-Middleware-1"))
	assert.Equal(t, "true", calls[0].Header.Get("X-Middleware-2"))
}

func TestWithTimeout(t *testing.T) {
	t.Run("adds timeout to request without deadline", func(t *testing.T) {
		mock := newMockRoundTripper()
		middleware := WithTimeout(1 * time.Second)
		roundTripper := middleware(mock)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		resp, err := roundTripper.RoundTrip(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.NoError(t, err)

		calls := mock.getCalls()
		require.Len(t, calls, 1)

		// Should have a deadline
		deadline, hasDeadline := calls[0].Context().Deadline()
		assert.True(t, hasDeadline)
		assert.WithinDuration(t, time.Now().Add(1*time.Second), deadline, 100*time.Millisecond)
	})

	t.Run("preserves existing deadline", func(t *testing.T) {
		mock := newMockRoundTripper()
		middleware := WithTimeout(1 * time.Second)
		roundTripper := middleware(mock)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		req := httptest.NewRequest(http.MethodGet, "/test", nil).WithContext(ctx)
		originalDeadline, _ := req.Context().Deadline()

		resp, err := roundTripper.RoundTrip(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.NoError(t, err)

		calls := mock.getCalls()
		require.Len(t, calls, 1)

		// Should preserve original deadline
		deadline, hasDeadline := calls[0].Context().Deadline()
		assert.True(t, hasDeadline)
		assert.Equal(t, originalDeadline, deadline)
	})

	t.Run("zero timeout does nothing", func(t *testing.T) {
		mock := newMockRoundTripper()
		middleware := WithTimeout(0)
		roundTripper := middleware(mock)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		resp, err := roundTripper.RoundTrip(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.NoError(t, err)

		calls := mock.getCalls()
		require.Len(t, calls, 1)

		// Should not have a deadline
		_, hasDeadline := calls[0].Context().Deadline()
		assert.False(t, hasDeadline)
	})
}

func TestWithLogging(t *testing.T) {
	mock := newMockRoundTripper()
	logger := logging.NoOpLogger{}
	middleware := WithLogging(logger)
	roundTripper := middleware(mock)

	// Add a successful response
	mock.addResponse(&http.Response{
		StatusCode:    http.StatusOK,
		ContentLength: 100,
		Body:          io.NopCloser(strings.NewReader("")),
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.ContentLength = 50

	resp, err := roundTripper.RoundTrip(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	calls := mock.getCalls()
	require.Len(t, calls, 1)
}

func TestWithLogging_Error(t *testing.T) {
	mock := newMockRoundTripper()
	logger := logging.NoOpLogger{}
	middleware := WithLogging(logger)
	roundTripper := middleware(mock)

	// Add an error response
	expectedErr := errors.New("network error")
	mock.addResponse(nil, expectedErr)

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)

	resp, err := roundTripper.RoundTrip(req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestWithRetry(t *testing.T) {
	t.Run("successful on first attempt", func(t *testing.T) {
		mock := newMockRoundTripper()
		middleware := WithRetry(3, DefaultShouldRetry)
		roundTripper := middleware(mock)

		mock.addResponse(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		resp, err := roundTripper.RoundTrip(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		calls := mock.getCalls()
		assert.Len(t, calls, 1) // Only one attempt
	})

	t.Run("retries on 500 error", func(t *testing.T) {
		mock := newMockRoundTripper()
		middleware := WithRetry(3, DefaultShouldRetry)
		roundTripper := middleware(mock)

		// First two attempts fail, third succeeds
		mock.addResponse(&http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader("error"))}, nil)
		mock.addResponse(&http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader("error"))}, nil)
		mock.addResponse(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		resp, err := roundTripper.RoundTrip(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		calls := mock.getCalls()
		assert.Len(t, calls, 3) // Three attempts
	})

	t.Run("fails after max attempts", func(t *testing.T) {
		mock := newMockRoundTripper()
		middleware := WithRetry(2, DefaultShouldRetry)
		roundTripper := middleware(mock)

		// All attempts fail
		networkErr := errors.New("network error")
		mock.addResponse(nil, networkErr)
		mock.addResponse(nil, networkErr)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		resp, err := roundTripper.RoundTrip(req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "all 2 attempts failed")

		calls := mock.getCalls()
		assert.Len(t, calls, 2) // Two attempts
	})

	t.Run("context cancellation", func(t *testing.T) {
		mock := newMockRoundTripper()
		middleware := WithRetry(3, DefaultShouldRetry)
		roundTripper := middleware(mock)

		// First attempt fails
		mock.addResponse(&http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader("error"))}, nil)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		req := httptest.NewRequest(http.MethodGet, "/test", nil).WithContext(ctx)
		resp, err := roundTripper.RoundTrip(req)

		assert.Nil(t, resp)
		assert.Equal(t, context.Canceled, err)
	})
}

func TestDefaultShouldRetry(t *testing.T) {
	tests := []struct {
		name     string
		resp     *http.Response
		err      error
		attempt  int
		expected bool
	}{
		{
			name:     "context canceled",
			resp:     nil,
			err:      context.Canceled,
			attempt:  0,
			expected: false,
		},
		{
			name:     "network error",
			resp:     nil,
			err:      errors.New("network error"),
			attempt:  0,
			expected: true,
		},
		{
			name:     "500 error",
			resp:     &http.Response{StatusCode: http.StatusInternalServerError},
			err:      nil,
			attempt:  0,
			expected: true,
		},
		{
			name:     "429 error",
			resp:     &http.Response{StatusCode: http.StatusTooManyRequests},
			err:      nil,
			attempt:  0,
			expected: true,
		},
		{
			name:     "200 success",
			resp:     &http.Response{StatusCode: http.StatusOK},
			err:      nil,
			attempt:  0,
			expected: false,
		},
		{
			name:     "400 error",
			resp:     &http.Response{StatusCode: http.StatusBadRequest},
			err:      nil,
			attempt:  0,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DefaultShouldRetry(tt.resp, tt.err, tt.attempt)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateBackoff(t *testing.T) {
	tests := []struct {
		attempt int
		minBase time.Duration
		maxBase time.Duration
	}{
		{attempt: 0, minBase: 1 * time.Second, maxBase: 2 * time.Second},
		{attempt: 1, minBase: 2 * time.Second, maxBase: 3 * time.Second},
		{attempt: 2, minBase: 4 * time.Second, maxBase: 5 * time.Second},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("attempt_%d", tt.attempt), func(t *testing.T) {
			backoff := calculateBackoff(tt.attempt)
			assert.GreaterOrEqual(t, backoff, tt.minBase)
			assert.LessOrEqual(t, backoff, tt.maxBase)
		})
	}
}

func TestWithHeaders(t *testing.T) {
	mock := newMockRoundTripper()
	headers := map[string]string{
		"X-Custom-Header": "custom-value",
		"Authorization":   "Bearer token",
	}
	middleware := WithHeaders(headers)
	roundTripper := middleware(mock)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := roundTripper.RoundTrip(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.NoError(t, err)

	calls := mock.getCalls()
	require.Len(t, calls, 1)

	assert.Equal(t, "custom-value", calls[0].Header.Get("X-Custom-Header"))
	assert.Equal(t, "Bearer token", calls[0].Header.Get("Authorization"))
}

func TestWithUserAgent(t *testing.T) {
	mock := newMockRoundTripper()
	middleware := WithUserAgent("test-agent/1.0")
	roundTripper := middleware(mock)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := roundTripper.RoundTrip(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.NoError(t, err)

	calls := mock.getCalls()
	require.Len(t, calls, 1)

	assert.Equal(t, "test-agent/1.0", calls[0].Header.Get("User-Agent"))
}

func TestWithRequestID(t *testing.T) {
	mock := newMockRoundTripper()

	idCounter := 0
	generator := func() string {
		idCounter++
		return fmt.Sprintf("req-%d", idCounter)
	}

	middleware := WithRequestID(generator)
	roundTripper := middleware(mock)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := roundTripper.RoundTrip(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.NoError(t, err)

	calls := mock.getCalls()
	require.Len(t, calls, 1)

	assert.Equal(t, "req-1", calls[0].Header.Get("X-Request-ID"))

	// Check context value - use exported context key from middleware package
	requestID := calls[0].Context().Value(ContextKeyRequestID)
	assert.Equal(t, "req-1", requestID)
}

func TestWithMetrics(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		mock := newMockRoundTripper()
		collector := newMockMetricsCollector()
		middleware := WithMetrics(collector)
		roundTripper := middleware(mock)

		mock.addResponse(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		resp, err := roundTripper.RoundTrip(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.NoError(t, err)

		collector.mu.Lock()
		assert.Len(t, collector.requests, 1)
		assert.Len(t, collector.responses, 1)
		assert.Len(t, collector.errors, 0)

		assert.Equal(t, "GET", collector.requests[0].method)
		assert.Equal(t, "/api/test", collector.requests[0].path)
		assert.Equal(t, 200, collector.responses[0].statusCode)
		collector.mu.Unlock()
	})

	t.Run("error request", func(t *testing.T) {
		mock := newMockRoundTripper()
		collector := newMockMetricsCollector()
		middleware := WithMetrics(collector)
		roundTripper := middleware(mock)

		expectedErr := errors.New("network error")
		mock.addResponse(nil, expectedErr)

		req := httptest.NewRequest(http.MethodPost, "/api/jobs", nil)
		_, err := roundTripper.RoundTrip(req)

		assert.Equal(t, expectedErr, err)

		collector.mu.Lock()
		assert.Len(t, collector.requests, 1)
		assert.Len(t, collector.responses, 0)
		assert.Len(t, collector.errors, 1)

		assert.Equal(t, "POST", collector.errors[0].method)
		assert.Equal(t, "/api/jobs", collector.errors[0].path)
		assert.Equal(t, expectedErr, collector.errors[0].err)
		collector.mu.Unlock()
	})
}

func TestCloneRequest(t *testing.T) {
	t.Run("request without body", func(t *testing.T) {
		original := httptest.NewRequest(http.MethodGet, "/test", nil)
		original.Header.Set("X-Original", "true")

		cloned := cloneRequest(original)

		assert.Equal(t, original.Method, cloned.Method)
		assert.Equal(t, original.URL.String(), cloned.URL.String())
		assert.Equal(t, original.Header.Get("X-Original"), cloned.Header.Get("X-Original"))
		// cloneRequest always creates a body even for requests without one
		// This is expected behavior based on the implementation
	})

	t.Run("request with body", func(t *testing.T) {
		body := "test body content"
		original := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(body))

		cloned := cloneRequest(original)

		// Both should have readable bodies
		originalBody, err := io.ReadAll(original.Body)
		assert.NoError(t, err)
		assert.Equal(t, body, string(originalBody))

		clonedBody, err := io.ReadAll(cloned.Body)
		assert.NoError(t, err)
		assert.Equal(t, body, string(clonedBody))
	})
}

func TestWithCircuitBreaker(t *testing.T) {
	t.Run("allows requests when under threshold", func(t *testing.T) {
		mock := newMockRoundTripper()
		middleware := WithCircuitBreaker(3, 1*time.Second)
		roundTripper := middleware(mock)

		mock.addResponse(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		resp, err := roundTripper.RoundTrip(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("opens circuit after threshold failures", func(t *testing.T) {
		mock := newMockRoundTripper()
		middleware := WithCircuitBreaker(2, 1*time.Second)
		roundTripper := middleware(mock)

		// Add failing responses
		mock.addResponse(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil)
		mock.addResponse(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		// First two should go through but fail
		resp1, err1 := roundTripper.RoundTrip(req)
		require.NoError(t, err1)
		defer resp1.Body.Close()
		assert.NoError(t, err1)
		assert.Equal(t, 500, resp1.StatusCode)

		resp2, err2 := roundTripper.RoundTrip(req)
		require.NoError(t, err2)
		defer resp2.Body.Close()
		assert.NoError(t, err2)
		assert.Equal(t, 500, resp2.StatusCode)

		// Third should be blocked by circuit breaker
		resp3, err3 := roundTripper.RoundTrip(req)
		assert.Nil(t, resp3)
		assert.Error(t, err3)
		assert.Contains(t, err3.Error(), "circuit breaker is open")
	})

	t.Run("circuit breaker with network error", func(t *testing.T) {
		mock := newMockRoundTripper()
		middleware := WithCircuitBreaker(1, 1*time.Second)
		roundTripper := middleware(mock)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		// Network error should trigger circuit breaker
		mock.addResponse(nil, errors.New("network error"))
		resp1, err1 := roundTripper.RoundTrip(req)
		assert.Nil(t, resp1)
		assert.Error(t, err1)
		assert.Equal(t, "network error", err1.Error())

		// Circuit should be open now
		resp2, err2 := roundTripper.RoundTrip(req)
		assert.Nil(t, resp2)
		assert.Error(t, err2)
		assert.Contains(t, err2.Error(), "circuit breaker is open")
	})
}

func TestCircuitBreaker(t *testing.T) {
	t.Run("initial state allows requests", func(t *testing.T) {
		cb := &circuitBreaker{threshold: 3, timeout: 1 * time.Second}
		assert.True(t, cb.Allow())
	})

	t.Run("allows requests under threshold", func(t *testing.T) {
		cb := &circuitBreaker{threshold: 3, timeout: 1 * time.Second}

		cb.RecordFailure()
		assert.True(t, cb.Allow())

		cb.RecordFailure()
		assert.True(t, cb.Allow())
	})

	t.Run("blocks requests at threshold", func(t *testing.T) {
		cb := &circuitBreaker{threshold: 2, timeout: 1 * time.Second}

		cb.RecordFailure()
		cb.RecordFailure()
		assert.False(t, cb.Allow())
	})

	t.Run("resets failure count on success", func(t *testing.T) {
		cb := &circuitBreaker{threshold: 2, timeout: 1 * time.Second}

		cb.RecordFailure()
		cb.RecordSuccess()
		assert.Equal(t, 0, cb.failures)
		assert.True(t, cb.Allow())
	})
}

func TestMiddlewareInterface(t *testing.T) {
	// Test that our middleware functions return the correct type
	var _ Middleware = WithTimeout(1 * time.Second)
	var _ Middleware = WithLogging(logging.NoOpLogger{})
	var _ Middleware = WithRetry(3, DefaultShouldRetry)
	var _ Middleware = WithHeaders(map[string]string{})
	var _ Middleware = WithUserAgent("test")
	var _ Middleware = WithRequestID(func() string { return "test" })
	var _ Middleware = WithCircuitBreaker(5, 1*time.Second)
}
