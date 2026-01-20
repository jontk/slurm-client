// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package context

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultTimeoutConfig(t *testing.T) {
	config := DefaultTimeoutConfig()

	require.NotNil(t, config)
	assert.Equal(t, DefaultTimeout, config.Default)
	assert.Equal(t, 30*time.Second, config.Read)
	assert.Equal(t, 1*time.Minute, config.Write)
	assert.Equal(t, 2*time.Minute, config.List)
	assert.Equal(t, time.Duration(0), config.Watch)
}

func TestWithTimeout(t *testing.T) {
	config := &TimeoutConfig{
		Default: 10 * time.Second,
		Read:    5 * time.Second,
		Write:   15 * time.Second,
		List:    30 * time.Second,
		Watch:   0, // No timeout
	}

	tests := []struct {
		name          string
		operationType OperationType
		expectedTime  time.Duration
		expectCancel  bool
	}{
		{
			name:          "read operation",
			operationType: OpRead,
			expectedTime:  5 * time.Second,
			expectCancel:  false,
		},
		{
			name:          "write operation",
			operationType: OpWrite,
			expectedTime:  15 * time.Second,
			expectCancel:  false,
		},
		{
			name:          "list operation",
			operationType: OpList,
			expectedTime:  30 * time.Second,
			expectCancel:  false,
		},
		{
			name:          "watch operation (no timeout)",
			operationType: OpWatch,
			expectedTime:  0,
			expectCancel:  true,
		},
		{
			name:          "default operation",
			operationType: OpDefault,
			expectedTime:  10 * time.Second,
			expectCancel:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			timeoutCtx, cancel := WithTimeout(ctx, tt.operationType, config)
			defer cancel()

			if tt.expectCancel {
				// For watch operations with no timeout, we expect a cancel context
				deadline, hasDeadline := timeoutCtx.Deadline()
				assert.False(t, hasDeadline)
				assert.True(t, deadline.IsZero())
			} else {
				// For other operations, we expect a timeout context
				deadline, hasDeadline := timeoutCtx.Deadline()
				assert.True(t, hasDeadline)

				// Check that the deadline is approximately correct
				expectedDeadline := time.Now().Add(tt.expectedTime)
				assert.WithinDuration(t, expectedDeadline, deadline, 100*time.Millisecond)
			}
		})
	}
}

func TestWithTimeoutNilConfig(t *testing.T) {
	ctx := context.Background()
	timeoutCtx, cancel := WithTimeout(ctx, OpRead, nil)
	defer cancel()

	deadline, hasDeadline := timeoutCtx.Deadline()
	assert.True(t, hasDeadline)

	// Should use default config
	expectedDeadline := time.Now().Add(30 * time.Second)
	assert.WithinDuration(t, expectedDeadline, deadline, 100*time.Millisecond)
}

func TestWithTimeoutWatchWithTimeout(t *testing.T) {
	config := &TimeoutConfig{
		Watch: 1 * time.Minute,
	}

	ctx := context.Background()
	timeoutCtx, cancel := WithTimeout(ctx, OpWatch, config)
	defer cancel()

	deadline, hasDeadline := timeoutCtx.Deadline()
	assert.True(t, hasDeadline)

	expectedDeadline := time.Now().Add(1 * time.Minute)
	assert.WithinDuration(t, expectedDeadline, deadline, 100*time.Millisecond)
}

func TestWithDeadline(t *testing.T) {
	t.Run("no existing deadline", func(t *testing.T) {
		ctx := context.Background()
		deadline := time.Now().Add(1 * time.Hour)

		deadlineCtx, cancel := WithDeadline(ctx, deadline)
		defer cancel()

		actualDeadline, hasDeadline := deadlineCtx.Deadline()
		assert.True(t, hasDeadline)
		assert.Equal(t, deadline, actualDeadline)
	})

	t.Run("existing deadline is sooner", func(t *testing.T) {
		soonerDeadline := time.Now().Add(1 * time.Hour)
		ctx, cancel := context.WithDeadline(context.Background(), soonerDeadline)
		defer cancel()

		laterDeadline := time.Now().Add(2 * time.Hour)
		deadlineCtx, cancelFunc := WithDeadline(ctx, laterDeadline)

		// Cancel function should be a no-op
		cancelFunc()

		actualDeadline, hasDeadline := deadlineCtx.Deadline()
		assert.True(t, hasDeadline)
		assert.Equal(t, soonerDeadline, actualDeadline)
		assert.Equal(t, ctx, deadlineCtx)
	})

	t.Run("existing deadline is later", func(t *testing.T) {
		laterDeadline := time.Now().Add(2 * time.Hour)
		ctx, cancel := context.WithDeadline(context.Background(), laterDeadline)
		defer cancel()

		soonerDeadline := time.Now().Add(1 * time.Hour)
		deadlineCtx, cancelFunc := WithDeadline(ctx, soonerDeadline)
		defer cancelFunc()

		actualDeadline, hasDeadline := deadlineCtx.Deadline()
		assert.True(t, hasDeadline)
		assert.Equal(t, soonerDeadline, actualDeadline)
	})
}

func TestEnsureTimeout(t *testing.T) {
	t.Run("no existing deadline", func(t *testing.T) {
		ctx := context.Background()
		defaultTimeout := 30 * time.Second

		timeoutCtx, cancel := EnsureTimeout(ctx, defaultTimeout)
		defer cancel()

		deadline, hasDeadline := timeoutCtx.Deadline()
		assert.True(t, hasDeadline)

		expectedDeadline := time.Now().Add(defaultTimeout)
		assert.WithinDuration(t, expectedDeadline, deadline, 100*time.Millisecond)
	})

	t.Run("existing deadline", func(t *testing.T) {
		existingDeadline := time.Now().Add(1 * time.Hour)
		ctx, cancel := context.WithDeadline(context.Background(), existingDeadline)
		defer cancel()

		timeoutCtx, cancelFunc := EnsureTimeout(ctx, 30*time.Second)

		// Cancel function should be a no-op
		cancelFunc()

		actualDeadline, hasDeadline := timeoutCtx.Deadline()
		assert.True(t, hasDeadline)
		assert.Equal(t, existingDeadline, actualDeadline)
		assert.Equal(t, ctx, timeoutCtx)
	})

	t.Run("zero default timeout", func(t *testing.T) {
		ctx := context.Background()

		timeoutCtx, cancel := EnsureTimeout(ctx, 0)
		defer cancel()

		deadline, hasDeadline := timeoutCtx.Deadline()
		assert.True(t, hasDeadline)

		// Should use DefaultTimeout
		expectedDeadline := time.Now().Add(DefaultTimeout)
		assert.WithinDuration(t, expectedDeadline, deadline, 100*time.Millisecond)
	})
}

func TestIsContextError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "context canceled",
			err:      context.Canceled,
			expected: true,
		},
		{
			name:     "context deadline exceeded",
			err:      context.DeadlineExceeded,
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("some other error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsContextError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContextError(t *testing.T) {
	t.Run("deadline exceeded", func(t *testing.T) {
		err := &ContextError{
			Operation: "test-operation",
			Timeout:   30 * time.Second,
			Err:       context.DeadlineExceeded,
		}

		expected := "operation 'test-operation' timed out after 30s"
		assert.Equal(t, expected, err.Error())
		assert.Equal(t, context.DeadlineExceeded, err.Unwrap())
	})

	t.Run("canceled", func(t *testing.T) {
		err := &ContextError{
			Operation: "test-operation",
			Timeout:   30 * time.Second,
			Err:       context.Canceled,
		}

		expected := "operation 'test-operation' was canceled"
		assert.Equal(t, expected, err.Error())
		assert.Equal(t, context.Canceled, err.Unwrap())
	})

	t.Run("other context error", func(t *testing.T) {
		customErr := errors.New("custom context error")
		err := &ContextError{
			Operation: "test-operation",
			Timeout:   30 * time.Second,
			Err:       customErr,
		}

		expected := "context error in operation 'test-operation': custom context error"
		assert.Equal(t, expected, err.Error())
		assert.Equal(t, customErr, err.Unwrap())
	})
}

func TestWrapContextError(t *testing.T) {
	t.Run("context error", func(t *testing.T) {
		operation := "test-operation"
		timeout := 30 * time.Second

		wrappedErr := WrapContextError(context.DeadlineExceeded, operation, timeout)

		require.IsType(t, &ContextError{}, wrappedErr)
		contextErr := wrappedErr.(*ContextError)
		assert.Equal(t, operation, contextErr.Operation)
		assert.Equal(t, timeout, contextErr.Timeout)
		assert.Equal(t, context.DeadlineExceeded, contextErr.Err)
	})

	t.Run("non-context error", func(t *testing.T) {
		originalErr := errors.New("not a context error")
		operation := "test-operation"
		timeout := 30 * time.Second

		wrappedErr := WrapContextError(originalErr, operation, timeout)

		// Should return the original error unchanged
		assert.Equal(t, originalErr, wrappedErr)
	})

	t.Run("nil error", func(t *testing.T) {
		operation := "test-operation"
		timeout := 30 * time.Second

		wrappedErr := WrapContextError(nil, operation, timeout)

		// Should return nil unchanged
		assert.Nil(t, wrappedErr)
	})
}

func TestOperationType(t *testing.T) {
	// Test that the operation types have expected values
	assert.Equal(t, OperationType(0), OpRead)
	assert.Equal(t, OperationType(1), OpWrite)
	assert.Equal(t, OperationType(2), OpList)
	assert.Equal(t, OperationType(3), OpWatch)
	assert.Equal(t, OperationType(4), OpDefault)
}

func TestConstants(t *testing.T) {
	assert.Equal(t, 30*time.Second, DefaultTimeout)
	assert.Equal(t, 5*time.Minute, DefaultLongTimeout)
}
