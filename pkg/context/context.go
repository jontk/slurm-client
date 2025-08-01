// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Package context provides utilities for context management in SLURM client operations
package context

import (
	"context"
	"time"
)

// DefaultTimeout is the default timeout for SLURM API operations
const DefaultTimeout = 30 * time.Second

// DefaultLongTimeout is used for operations that may take longer (e.g., large job listings)
const DefaultLongTimeout = 5 * time.Minute

// TimeoutConfig holds timeout configuration for different operation types
type TimeoutConfig struct {
	// Default timeout for most operations
	Default time.Duration

	// Timeout for read operations (GET requests)
	Read time.Duration

	// Timeout for write operations (POST, PUT, DELETE)
	Write time.Duration

	// Timeout for list operations that may return large datasets
	List time.Duration

	// Timeout for watch/streaming operations
	Watch time.Duration
}

// DefaultTimeoutConfig returns a timeout configuration with sensible defaults
func DefaultTimeoutConfig() *TimeoutConfig {
	return &TimeoutConfig{
		Default: DefaultTimeout,
		Read:    30 * time.Second,
		Write:   1 * time.Minute,
		List:    2 * time.Minute,
		Watch:   0, // No timeout for watch operations
	}
}

// WithTimeout adds a timeout to the context based on the operation type
func WithTimeout(ctx context.Context, operationType OperationType, config *TimeoutConfig) (context.Context, context.CancelFunc) {
	if config == nil {
		config = DefaultTimeoutConfig()
	}

	timeout := config.Default
	switch operationType {
	case OpRead:
		timeout = config.Read
	case OpWrite:
		timeout = config.Write
	case OpList:
		timeout = config.List
	case OpWatch:
		if config.Watch == 0 {
			// No timeout for watch operations
			return context.WithCancel(ctx)
		}
		timeout = config.Watch
	}

	return context.WithTimeout(ctx, timeout)
}

// OperationType represents the type of operation being performed
type OperationType int

const (
	// OpRead represents a read operation (GET)
	OpRead OperationType = iota
	// OpWrite represents a write operation (POST, PUT, DELETE)
	OpWrite
	// OpList represents a list operation that may return large datasets
	OpList
	// OpWatch represents a watch/streaming operation
	OpWatch
	// OpDefault represents any other operation
	OpDefault
)

// WithDeadline adds a deadline to the context if it doesn't already have one
func WithDeadline(ctx context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
	// Check if context already has a deadline
	if existing, ok := ctx.Deadline(); ok && existing.Before(deadline) {
		// Existing deadline is sooner, keep it
		return ctx, func() {}
	}
	return context.WithDeadline(ctx, deadline)
}

// EnsureTimeout ensures the context has a timeout, adding a default if needed
func EnsureTimeout(ctx context.Context, defaultTimeout time.Duration) (context.Context, context.CancelFunc) {
	// Check if context already has a deadline
	if _, ok := ctx.Deadline(); ok {
		// Already has a deadline, return as-is
		return ctx, func() {}
	}
	
	if defaultTimeout == 0 {
		defaultTimeout = DefaultTimeout
	}
	
	return context.WithTimeout(ctx, defaultTimeout)
}

// IsContextError checks if an error is a context-related error
func IsContextError(err error) bool {
	if err == nil {
		return false
	}
	return err == context.Canceled || err == context.DeadlineExceeded
}

// ContextError wraps context errors with more descriptive messages
type ContextError struct {
	Operation string
	Timeout   time.Duration
	Err       error
}

func (e *ContextError) Error() string {
	if e.Err == context.DeadlineExceeded {
		return "operation '" + e.Operation + "' timed out after " + e.Timeout.String()
	}
	if e.Err == context.Canceled {
		return "operation '" + e.Operation + "' was canceled"
	}
	return "context error in operation '" + e.Operation + "': " + e.Err.Error()
}

func (e *ContextError) Unwrap() error {
	return e.Err
}

// WrapContextError wraps a context error with operation details
func WrapContextError(err error, operation string, timeout time.Duration) error {
	if !IsContextError(err) {
		return err
	}
	return &ContextError{
		Operation: operation,
		Timeout:   timeout,
		Err:       err,
	}
}
