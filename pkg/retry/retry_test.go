// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package retry

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/jontk/slurm-client/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestHTTPExponentialBackoff_Default(t *testing.T) {
	policy := NewHTTPExponentialBackoff()

	// Test default values
	helpers.AssertEqual(t, 3, policy.MaxRetries())
	helpers.AssertEqual(t, 1*time.Second, policy.minWaitTime)
	helpers.AssertEqual(t, 30*time.Second, policy.maxWaitTime)
	helpers.AssertEqual(t, 2.0, policy.backoffFactor)
	helpers.AssertEqual(t, true, policy.jitter)
}

func TestHTTPExponentialBackoff_WithMethods(t *testing.T) {
	policy := NewHTTPExponentialBackoff().
		WithMaxRetries(5).
		WithMinWaitTime(2 * time.Second).
		WithMaxWaitTime(60 * time.Second).
		WithBackoffFactor(1.5).
		WithJitter(false)

	helpers.AssertEqual(t, 5, policy.MaxRetries())
	helpers.AssertEqual(t, 2*time.Second, policy.minWaitTime)
	helpers.AssertEqual(t, 60*time.Second, policy.maxWaitTime)
	helpers.AssertEqual(t, 1.5, policy.backoffFactor)
	helpers.AssertEqual(t, false, policy.jitter)
}

func TestHTTPExponentialBackoff_ShouldRetry(t *testing.T) {
	policy := NewHTTPExponentialBackoff().WithMaxRetries(3)
	ctx := helpers.TestContext(t)

	tests := []struct {
		name        string
		resp        *http.Response
		err         error
		attempt     int
		shouldRetry bool
	}{
		{
			name:        "network error should retry",
			resp:        nil,
			err:         errors.New("network error"),
			attempt:     1,
			shouldRetry: true,
		},
		{
			name:        "max retries exceeded",
			resp:        nil,
			err:         errors.New("network error"),
			attempt:     3,
			shouldRetry: false,
		},
		{
			name:        "500 status should retry",
			resp:        &http.Response{StatusCode: 500},
			err:         nil,
			attempt:     1,
			shouldRetry: true,
		},
		{
			name:        "503 status should retry",
			resp:        &http.Response{StatusCode: 503},
			err:         nil,
			attempt:     1,
			shouldRetry: true,
		},
		{
			name:        "429 status should retry",
			resp:        &http.Response{StatusCode: 429},
			err:         nil,
			attempt:     1,
			shouldRetry: true,
		},
		{
			name:        "200 status should not retry",
			resp:        &http.Response{StatusCode: 200},
			err:         nil,
			attempt:     1,
			shouldRetry: false,
		},
		{
			name:        "404 status should not retry",
			resp:        &http.Response{StatusCode: 404},
			err:         nil,
			attempt:     1,
			shouldRetry: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := policy.ShouldRetry(ctx, tt.resp, tt.err, tt.attempt)
			helpers.AssertEqual(t, tt.shouldRetry, result)
		})
	}
}

func TestHTTPExponentialBackoff_ShouldRetryWithCancelledContext(t *testing.T) {
	policy := NewHTTPExponentialBackoff()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context

	// Should not retry when context is cancelled
	result := policy.ShouldRetry(ctx, nil, errors.New("error"), 1)
	helpers.AssertEqual(t, false, result)
}

func TestHTTPExponentialBackoff_WaitTime(t *testing.T) {
	policy := NewHTTPExponentialBackoff().
		WithMinWaitTime(1 * time.Second).
		WithMaxWaitTime(10 * time.Second).
		WithBackoffFactor(2.0).
		WithJitter(false) // Disable jitter for predictable testing

	tests := []struct {
		name        string
		attempt     int
		expectedMin time.Duration
		expectedMax time.Duration
	}{
		{
			name:        "attempt 0",
			attempt:     0,
			expectedMin: 1 * time.Second,
			expectedMax: 1 * time.Second,
		},
		{
			name:        "attempt 1",
			attempt:     1,
			expectedMin: 1 * time.Second,
			expectedMax: 1 * time.Second,
		},
		{
			name:        "attempt 2",
			attempt:     2,
			expectedMin: 2 * time.Second,
			expectedMax: 2 * time.Second,
		},
		{
			name:        "attempt 3",
			attempt:     3,
			expectedMin: 4 * time.Second,
			expectedMax: 4 * time.Second,
		},
		{
			name:        "attempt 4 (hits max)",
			attempt:     4,
			expectedMin: 8 * time.Second,
			expectedMax: 10 * time.Second, // Should be capped at max
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			waitTime := policy.WaitTime(tt.attempt)

			if tt.expectedMin == tt.expectedMax {
				helpers.AssertEqual(t, tt.expectedMin, waitTime)
			} else {
				assert.GreaterOrEqual(t, waitTime, tt.expectedMin)
				assert.LessOrEqual(t, waitTime, tt.expectedMax)
			}
		})
	}
}

func TestHTTPExponentialBackoff_WaitTimeWithJitter(t *testing.T) {
	policy := NewHTTPExponentialBackoff().
		WithMinWaitTime(1 * time.Second).
		WithMaxWaitTime(10 * time.Second).
		WithBackoffFactor(2.0).
		WithJitter(true)

	// Test that jitter adds some randomness
	waitTime1 := policy.WaitTime(2)
	waitTime2 := policy.WaitTime(2)

	// With jitter, the wait times should be at least the base wait time
	baseWaitTime := 2 * time.Second
	assert.GreaterOrEqual(t, waitTime1, baseWaitTime)
	assert.GreaterOrEqual(t, waitTime2, baseWaitTime)

	// Due to jitter, wait times might be different (though they could be the same due to randomness)
	// We can't guarantee they'll be different, but we can test the bounds
	assert.LessOrEqual(t, waitTime1, baseWaitTime+time.Duration(float64(baseWaitTime)*0.1))
	assert.LessOrEqual(t, waitTime2, baseWaitTime+time.Duration(float64(baseWaitTime)*0.1))
}

func TestFixedDelay(t *testing.T) {
	maxRetries := 3
	delay := 5 * time.Second
	policy := NewFixedDelay(maxRetries, delay)

	// Test basic properties
	helpers.AssertEqual(t, maxRetries, policy.MaxRetries())
	helpers.AssertEqual(t, delay, policy.WaitTime(1))
	helpers.AssertEqual(t, delay, policy.WaitTime(5)) // Should always return same delay

	ctx := helpers.TestContext(t)

	// Test retry logic
	helpers.AssertEqual(t, true, policy.ShouldRetry(ctx, nil, errors.New("error"), 1))
	helpers.AssertEqual(t, true, policy.ShouldRetry(ctx, &http.Response{StatusCode: 500}, nil, 2))
	helpers.AssertEqual(t, false, policy.ShouldRetry(ctx, nil, errors.New("error"), 3)) // Max retries exceeded
	helpers.AssertEqual(t, false, policy.ShouldRetry(ctx, &http.Response{StatusCode: 200}, nil, 1))
}

func TestFixedDelay_ShouldRetryWithCancelledContext(t *testing.T) {
	policy := NewFixedDelay(3, 1*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context

	// Should not retry when context is cancelled
	result := policy.ShouldRetry(ctx, nil, errors.New("error"), 1)
	helpers.AssertEqual(t, false, result)
}

func TestNoRetry(t *testing.T) {
	policy := NewNoRetry()

	// Test basic properties
	helpers.AssertEqual(t, 0, policy.MaxRetries())
	helpers.AssertEqual(t, time.Duration(0), policy.WaitTime(1))

	ctx := helpers.TestContext(t)

	// Should never retry
	helpers.AssertEqual(t, false, policy.ShouldRetry(ctx, nil, errors.New("error"), 0))
	helpers.AssertEqual(t, false, policy.ShouldRetry(ctx, &http.Response{StatusCode: 500}, nil, 0))
	helpers.AssertEqual(t, false, policy.ShouldRetry(ctx, nil, errors.New("error"), 1))
}

func TestPolicyInterface(t *testing.T) {
	// Test that all retry policies implement the Policy interface
	var _ Policy = &HTTPExponentialBackoff{}
	var _ Policy = &FixedDelay{}
	var _ Policy = &NoRetry{}

	// Test different policies
	policies := []Policy{
		NewHTTPExponentialBackoff(),
		NewFixedDelay(3, 1*time.Second),
		NewNoRetry(),
	}

	ctx := helpers.TestContext(t)

	for _, policy := range policies {
		// Each policy should have max retries
		maxRetries := policy.MaxRetries()
		assert.GreaterOrEqual(t, maxRetries, 0)

		// Each policy should return wait time
		waitTime := policy.WaitTime(1)
		assert.GreaterOrEqual(t, waitTime, time.Duration(0))

		// Each policy should respond to ShouldRetry
		shouldRetry := policy.ShouldRetry(ctx, nil, errors.New("error"), 0)
		// We don't assert a specific value since it depends on the policy
		_ = shouldRetry
	}
}

func TestRetryableHTTPStatusCodes(t *testing.T) {
	policy := NewHTTPExponentialBackoff()
	ctx := helpers.TestContext(t)

	retryableStatusCodes := []int{
		http.StatusTooManyRequests,     // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout,      // 504
	}

	nonRetryableStatusCodes := []int{
		http.StatusOK,                  // 200
		http.StatusBadRequest,          // 400
		http.StatusUnauthorized,        // 401
		http.StatusForbidden,           // 403
		http.StatusNotFound,            // 404
		http.StatusMethodNotAllowed,    // 405
		http.StatusConflict,            // 409
		http.StatusUnprocessableEntity, // 422
	}

	for _, statusCode := range retryableStatusCodes {
		t.Run("retryable_"+http.StatusText(statusCode), func(t *testing.T) {
			resp := &http.Response{StatusCode: statusCode}
			result := policy.ShouldRetry(ctx, resp, nil, 1)
			helpers.AssertEqual(t, true, result)
		})
	}

	for _, statusCode := range nonRetryableStatusCodes {
		t.Run("non_retryable_"+http.StatusText(statusCode), func(t *testing.T) {
			resp := &http.Response{StatusCode: statusCode}
			result := policy.ShouldRetry(ctx, resp, nil, 1)
			helpers.AssertEqual(t, false, result)
		})
	}
}
