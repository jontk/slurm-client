// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package benchmarks

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jontk/slurm-client/pkg/retry"
)

func BenchmarkRetryableCheck(b *testing.B) {
	policy := retry.NewHTTPExponentialBackoff().
		WithMaxRetries(3).
		WithMinWaitTime(1 * time.Second).
		WithMaxWaitTime(30 * time.Second)

	// Different response codes to test
	codes := []int{500, 502, 503, 504, 429, 400, 401, 404}
	responses := make([]*http.Response, len(codes))
	for i, code := range codes {
		responses[i] = &http.Response{StatusCode: code}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := responses[i%len(responses)]
		_ = policy.ShouldRetry(context.Background(), resp, nil, 1)
	}
}

func BenchmarkBackoffCalculation(b *testing.B) {
	policy := retry.NewHTTPExponentialBackoff().
		WithMaxRetries(5).
		WithMinWaitTime(1 * time.Second).
		WithMaxWaitTime(30 * time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		attempt := i % 5
		_ = policy.WaitTime(attempt)
	}
}

func BenchmarkRetryWithError(b *testing.B) {
	policy := retry.NewHTTPExponentialBackoff().
		WithMaxRetries(3).
		WithMinWaitTime(1 * time.Millisecond). // Short for benchmarking
		WithMaxWaitTime(10 * time.Millisecond)

	// Simulate a function that always fails
	failingFunc := func() (*http.Response, error) {
		return nil, fmt.Errorf("connection refused")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		attempts := 0
		
		// Simulate retry loop without actual delays
		for attempts < policy.MaxRetries() {
			_, err := failingFunc()
			if err == nil {
				break
			}
			
			// Check if we should retry
			if !policy.ShouldRetry(ctx, nil, err, attempts) {
				break
			}
			attempts++
		}
	}
}