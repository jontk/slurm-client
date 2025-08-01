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
	client := &retry.Client{
		RetryMax:     3,
		RetryWaitMin: 1 * time.Second,
		RetryWaitMax: 30 * time.Second,
	}

	// Different response codes to test
	codes := []int{500, 502, 503, 504, 429, 400, 401, 404}
	responses := make([]*http.Response, len(codes))
	for i, code := range codes {
		responses[i] = &http.Response{StatusCode: code}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := responses[i%len(responses)]
		_ = client.CheckRetry(context.Background(), resp, nil)
	}
}

func BenchmarkBackoffCalculation(b *testing.B) {
	client := &retry.Client{
		RetryMax:     5,
		RetryWaitMin: 1 * time.Second,
		RetryWaitMax: 30 * time.Second,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		attempt := i % 5
		min := float64(client.RetryWaitMin)
		max := float64(client.RetryWaitMax)
		
		// Simulate backoff calculation
		wait := min * float64(1<<uint(attempt))
		if wait > max {
			wait = max
		}
		_ = time.Duration(wait)
	}
}

func BenchmarkRetryWithError(b *testing.B) {
	client := &retry.Client{
		RetryMax:     3,
		RetryWaitMin: 1 * time.Millisecond, // Short for benchmarking
		RetryWaitMax: 10 * time.Millisecond,
	}

	// Simulate a function that always fails
	failingFunc := func() (*http.Response, error) {
		return nil, fmt.Errorf("connection refused")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		attempts := 0
		
		// Simulate retry loop without actual delays
		for attempts < client.RetryMax {
			_, err := failingFunc()
			if err == nil {
				break
			}
			attempts++
		}
	}
}