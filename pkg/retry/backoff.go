// Package retry provides advanced retry mechanisms with various backoff strategies
package retry

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// BackoffStrategy defines the interface for backoff strategies
type BackoffStrategy interface {
	// NextDelay returns the next delay duration and whether to continue retrying
	NextDelay(attempt int) (time.Duration, bool)
	// Reset resets the backoff strategy
	Reset()
}

// ExponentialBackoff implements exponential backoff with jitter
type ExponentialBackoff struct {
	// InitialDelay is the delay for the first retry
	InitialDelay time.Duration
	
	// MaxDelay is the maximum delay between retries
	MaxDelay time.Duration
	
	// Multiplier is the factor by which the delay increases
	Multiplier float64
	
	// Jitter adds randomness to prevent thundering herd
	Jitter float64
	
	// MaxAttempts is the maximum number of retry attempts
	MaxAttempts int
}

// NewExponentialBackoff creates a new exponential backoff strategy with defaults
func NewExponentialBackoff() *ExponentialBackoff {
	return &ExponentialBackoff{
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     30 * time.Second,
		Multiplier:   2.0,
		Jitter:       0.1,
		MaxAttempts:  5,
	}
}

// NextDelay calculates the next delay with exponential backoff and jitter
func (e *ExponentialBackoff) NextDelay(attempt int) (time.Duration, bool) {
	if attempt >= e.MaxAttempts {
		return 0, false
	}
	
	// Calculate base delay
	delay := float64(e.InitialDelay) * math.Pow(e.Multiplier, float64(attempt))
	
	// Cap at max delay
	if delay > float64(e.MaxDelay) {
		delay = float64(e.MaxDelay)
	}
	
	// Add jitter
	if e.Jitter > 0 {
		jitter := delay * e.Jitter
		delay = delay - jitter + (2 * jitter * rand.Float64())
	}
	
	return time.Duration(delay), true
}

// Reset resets the backoff strategy (no-op for stateless exponential backoff)
func (e *ExponentialBackoff) Reset() {}

// LinearBackoff implements linear backoff strategy
type LinearBackoff struct {
	// InitialDelay is the delay for the first retry
	InitialDelay time.Duration
	
	// Increment is added to the delay for each retry
	Increment time.Duration
	
	// MaxDelay is the maximum delay between retries
	MaxDelay time.Duration
	
	// Jitter adds randomness to prevent thundering herd
	Jitter float64
	
	// MaxAttempts is the maximum number of retry attempts
	MaxAttempts int
}

// NewLinearBackoff creates a new linear backoff strategy with defaults
func NewLinearBackoff() *LinearBackoff {
	return &LinearBackoff{
		InitialDelay: 100 * time.Millisecond,
		Increment:    100 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Jitter:       0.1,
		MaxAttempts:  5,
	}
}

// NextDelay calculates the next delay with linear backoff
func (l *LinearBackoff) NextDelay(attempt int) (time.Duration, bool) {
	if attempt >= l.MaxAttempts {
		return 0, false
	}
	
	// Calculate base delay
	delay := l.InitialDelay + (time.Duration(attempt) * l.Increment)
	
	// Cap at max delay
	if delay > l.MaxDelay {
		delay = l.MaxDelay
	}
	
	// Add jitter
	if l.Jitter > 0 {
		jitter := float64(delay) * l.Jitter
		delayFloat := float64(delay) - jitter + (2 * jitter * rand.Float64())
		delay = time.Duration(delayFloat)
	}
	
	return delay, true
}

// Reset resets the backoff strategy (no-op for stateless linear backoff)
func (l *LinearBackoff) Reset() {}

// FibonacciBackoff implements Fibonacci sequence backoff
type FibonacciBackoff struct {
	// InitialDelay is the base unit for the Fibonacci sequence
	InitialDelay time.Duration
	
	// MaxDelay is the maximum delay between retries
	MaxDelay time.Duration
	
	// MaxAttempts is the maximum number of retry attempts
	MaxAttempts int
	
	// fibonacci sequence cache
	fib []int
}

// NewFibonacciBackoff creates a new Fibonacci backoff strategy
func NewFibonacciBackoff() *FibonacciBackoff {
	return &FibonacciBackoff{
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     30 * time.Second,
		MaxAttempts:  10,
		fib:          []int{1, 1},
	}
}

// NextDelay calculates the next delay using Fibonacci sequence
func (f *FibonacciBackoff) NextDelay(attempt int) (time.Duration, bool) {
	if attempt >= f.MaxAttempts {
		return 0, false
	}
	
	// Calculate Fibonacci number for this attempt
	fibValue := f.fibonacci(attempt)
	delay := time.Duration(fibValue) * f.InitialDelay
	
	// Cap at max delay
	if delay > f.MaxDelay {
		delay = f.MaxDelay
	}
	
	return delay, true
}

// fibonacci returns the nth Fibonacci number
func (f *FibonacciBackoff) fibonacci(n int) int {
	// Extend cache if needed
	for len(f.fib) <= n {
		next := f.fib[len(f.fib)-1] + f.fib[len(f.fib)-2]
		f.fib = append(f.fib, next)
	}
	return f.fib[n]
}

// Reset resets the backoff strategy
func (f *FibonacciBackoff) Reset() {
	f.fib = []int{1, 1}
}

// ConstantBackoff implements a constant delay backoff
type ConstantBackoff struct {
	// Delay is the constant delay between retries
	Delay time.Duration
	
	// MaxAttempts is the maximum number of retry attempts
	MaxAttempts int
}

// NewConstantBackoff creates a new constant backoff strategy
func NewConstantBackoff(delay time.Duration, maxAttempts int) *ConstantBackoff {
	return &ConstantBackoff{
		Delay:       delay,
		MaxAttempts: maxAttempts,
	}
}

// NextDelay returns the constant delay
func (c *ConstantBackoff) NextDelay(attempt int) (time.Duration, bool) {
	if attempt >= c.MaxAttempts {
		return 0, false
	}
	return c.Delay, true
}

// Reset resets the backoff strategy (no-op for constant backoff)
func (c *ConstantBackoff) Reset() {}

// Retry executes a function with the specified backoff strategy
func Retry(ctx context.Context, backoff BackoffStrategy, fn func() error) error {
	var lastErr error
	
	for attempt := 0; ; attempt++ {
		// Execute the function
		err := fn()
		if err == nil {
			return nil
		}
		
		lastErr = err
		
		// Get next delay
		delay, shouldContinue := backoff.NextDelay(attempt)
		if !shouldContinue {
			break
		}
		
		// Wait for the delay or context cancellation
		select {
		case <-time.After(delay):
			// Continue to next attempt
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	
	return lastErr
}

// RetryWithResult executes a function with retry and returns the result
func RetryWithResult[T any](ctx context.Context, backoff BackoffStrategy, fn func() (T, error)) (T, error) {
	var zero T
	var lastErr error
	
	for attempt := 0; ; attempt++ {
		// Execute the function
		result, err := fn()
		if err == nil {
			return result, nil
		}
		
		lastErr = err
		
		// Get next delay
		delay, shouldContinue := backoff.NextDelay(attempt)
		if !shouldContinue {
			break
		}
		
		// Wait for the delay or context cancellation
		select {
		case <-time.After(delay):
			// Continue to next attempt
		case <-ctx.Done():
			return zero, ctx.Err()
		}
	}
	
	return zero, lastErr
}