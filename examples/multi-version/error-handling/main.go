// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/retry"
)

func main() {
	ctx := context.Background()

	fmt.Println("=== Error Handling and Retry Patterns ===")

	// Example 1: Basic error handling
	fmt.Println("\n--- Example 1: Basic Error Handling ---")
	demonstrateBasicErrorHandling(ctx)

	// Example 2: Advanced retry policies
	fmt.Println("\n--- Example 2: Advanced Retry Policies ---")
	demonstrateAdvancedRetryPolicies(ctx)

	// Example 3: Circuit breaker pattern
	fmt.Println("\n--- Example 3: Circuit Breaker Pattern ---")
	demonstrateCircuitBreakerPattern(ctx)

	// Example 4: Graceful degradation
	fmt.Println("\n--- Example 4: Graceful Degradation ---")
	demonstrateGracefulDegradation(ctx)

	// Example 5: Timeout and cancellation
	fmt.Println("\n--- Example 5: Timeout and Cancellation ---")
	demonstrateTimeoutAndCancellation(ctx)
}

func demonstrateBasicErrorHandling(ctx context.Context) {
	// Create client with basic error handling
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://invalid-server:6820"),
		slurm.WithAuth(auth.NewNoneAuth()),
		slurm.WithTimeout("5s"),
	)
	if err != nil {
		log.Printf("✓ Expected error creating client: %v", err)
		return
	}

	// Test connection error handling
	_, err = client.GetInfo(ctx)
	if err != nil {
		handleSlurmError("GetInfo", err)
	}

	// Test authentication error
	authClient, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(auth.NewTokenAuth("invalid-token")),
		slurm.WithTimeout("5s"),
	)
	if err != nil {
		log.Printf("Error creating auth client: %v", err)
		return
	}

	_, err = authClient.GetInfo(ctx)
	if err != nil {
		handleSlurmError("GetInfo (auth)", err)
	}

	// Test invalid job submission
	invalidJobReq := &slurm.JobSubmissionRequest{
		Script:    "#!/bin/bash\necho 'test'",
		Name:      "",  // Invalid: empty name
		Partition: "nonexistent-partition",
		Nodes:     -1,  // Invalid: negative nodes
		Tasks:     0,   // Invalid: zero tasks
		Time:      "invalid-time-format",
	}

	_, err = authClient.SubmitJob(ctx, invalidJobReq)
	if err != nil {
		handleSlurmError("SubmitJob (validation)", err)
	}
}

func demonstrateAdvancedRetryPolicies(ctx context.Context) {
	// Example 1: Exponential backoff with jitter
	fmt.Println("Testing exponential backoff with jitter...")
	
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(auth.NewNoneAuth()),
		slurm.WithRetryPolicy(5, "100ms", "10s"),
		slurm.WithDebug(true),
	)
	if err != nil {
		log.Printf("Error creating retry client: %v", err)
		return
	}

	// This will likely fail and demonstrate retry behavior
	start := time.Now()
	_, err = client.GetInfo(ctx)
	duration := time.Since(start)
	
	if err != nil {
		fmt.Printf("✓ Retry policy executed in %s: %v\n", duration, err)
	}

	// Example 2: Custom retry policy
	fmt.Println("\nTesting custom retry policy...")
	
	customRetryPolicy := &retry.Policy{
		MaxRetries: 3,
		BaseDelay:  200 * time.Millisecond,
		MaxDelay:   5 * time.Second,
		Multiplier: 2.0,
		Jitter:     true,
		RetryableErrors: []string{
			"connection refused",
			"timeout",
			"temporary failure",
		},
	}

	retryClient, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(auth.NewNoneAuth()),
		slurm.WithCustomRetryPolicy(customRetryPolicy),
	)
	if err != nil {
		log.Printf("Error creating custom retry client: %v", err)
		return
	}

	start = time.Now()
	_, err = retryClient.ListJobs(ctx)
	duration = time.Since(start)
	
	if err != nil {
		fmt.Printf("✓ Custom retry policy executed in %s: %v\n", duration, err)
	}
}

func demonstrateCircuitBreakerPattern(ctx context.Context) {
	// Simulate circuit breaker using a wrapper
	circuitBreaker := &CircuitBreaker{
		threshold:    3,
		timeout:      10 * time.Second,
		successCount: 0,
		failureCount: 0,
		state:        "CLOSED",
	}

	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(auth.NewNoneAuth()),
		slurm.WithTimeout("2s"),
	)
	if err != nil {
		log.Printf("Error creating circuit breaker client: %v", err)
		return
	}

	// Test multiple calls to trigger circuit breaker
	for i := 1; i <= 5; i++ {
		fmt.Printf("Attempt %d - Circuit breaker state: %s\n", i, circuitBreaker.state)
		
		err := circuitBreaker.Call(func() error {
			_, err := client.GetInfo(ctx)
			return err
		})

		if err != nil {
			fmt.Printf("  ✗ Call failed: %v\n", err)
		} else {
			fmt.Printf("  ✓ Call succeeded\n")
		}
		
		time.Sleep(100 * time.Millisecond)
	}
}

func demonstrateGracefulDegradation(ctx context.Context) {
	// Create primary and fallback clients
	primaryClient, err := slurm.NewClientWithVersion(ctx, slurm.LatestVersion(),
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(auth.NewNoneAuth()),
		slurm.WithTimeout("2s"),
	)
	if err != nil {
		log.Printf("Error creating primary client: %v", err)
	}

	fallbackClient, err := slurm.NewClientWithVersion(ctx, slurm.StableVersion(),
		slurm.WithBaseURL("https://backup-server:6820"),
		slurm.WithAuth(auth.NewNoneAuth()),
		slurm.WithTimeout("5s"),
	)
	if err != nil {
		log.Printf("Error creating fallback client: %v", err)
	}

	// Try primary, fallback to stable version
	info, err := getInfoWithFallback(ctx, primaryClient, fallbackClient)
	if err != nil {
		fmt.Printf("✗ Both primary and fallback failed: %v\n", err)
		
		// Last resort: offline mode with cached data
		info = getCachedInfo()
		fmt.Printf("✓ Using cached data: %s\n", info.Version)
	} else {
		fmt.Printf("✓ Got server info: %s\n", info.Version)
	}

	// Test graceful job submission degradation
	jobReq := &slurm.JobSubmissionRequest{
		Script: "#!/bin/bash\necho 'resilient job'",
		Name:   "resilient-job",
		Time:   "00:01:00",
	}

	jobID, err := submitJobWithDegradation(ctx, primaryClient, fallbackClient, jobReq)
	if err != nil {
		fmt.Printf("✗ Job submission failed completely: %v\n", err)
	} else {
		fmt.Printf("✓ Job submitted successfully: %d\n", jobID)
	}
}

func demonstrateTimeoutAndCancellation(ctx context.Context) {
	// Example 1: Request timeout
	fmt.Println("Testing request timeout...")
	
	shortTimeoutClient, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(auth.NewNoneAuth()),
		slurm.WithTimeout("100ms"), // Very short timeout
	)
	if err != nil {
		log.Printf("Error creating timeout client: %v", err)
		return
	}

	start := time.Now()
	_, err = shortTimeoutClient.GetInfo(ctx)
	duration := time.Since(start)
	
	if err != nil {
		fmt.Printf("✓ Request timed out after %s: %v\n", duration, err)
	}

	// Example 2: Context cancellation
	fmt.Println("\nTesting context cancellation...")
	
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(auth.NewNoneAuth()),
		slurm.WithTimeout("10s"),
	)
	if err != nil {
		log.Printf("Error creating cancellation client: %v", err)
		return
	}

	// Create context with cancellation
	cancelCtx, cancel := context.WithCancel(ctx)
	
	// Start operation
	go func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("Cancelling operation...")
		cancel()
	}()

	start = time.Now()
	_, err = client.ListJobs(cancelCtx)
	duration = time.Since(start)
	
	if err != nil {
		fmt.Printf("✓ Operation cancelled after %s: %v\n", duration, err)
	}

	// Example 3: Deadline exceeded
	fmt.Println("\nTesting deadline exceeded...")
	
	deadlineCtx, deadlineCancel := context.WithDeadline(ctx, time.Now().Add(200*time.Millisecond))
	defer deadlineCancel()

	start = time.Now()
	_, err = client.GetInfo(deadlineCtx)
	duration = time.Since(start)
	
	if err != nil {
		fmt.Printf("✓ Deadline exceeded after %s: %v\n", duration, err)
	}
}

// Helper functions

func handleSlurmError(operation string, err error) {
	fmt.Printf("Error in %s: ", operation)

	// Type switch for different error types
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		fmt.Printf("✓ Timeout error: %v\n", err)
	case errors.Is(err, context.Canceled):
		fmt.Printf("✓ Cancelled error: %v\n", err)
	case isNetworkError(err):
		fmt.Printf("✓ Network error: %v\n", err)
	case isAuthenticationError(err):
		fmt.Printf("✓ Authentication error: %v\n", err)
	case isValidationError(err):
		fmt.Printf("✓ Validation error: %v\n", err)
	case isServerError(err):
		fmt.Printf("✓ Server error: %v\n", err)
	default:
		fmt.Printf("✓ Unknown error: %v\n", err)
	}
}

func isNetworkError(err error) bool {
	return err != nil && (
		err.Error() == "connection refused" ||
		err.Error() == "no route to host" ||
		err.Error() == "network unreachable")
}

func isAuthenticationError(err error) bool {
	return err != nil && (
		err.Error() == "401 Unauthorized" ||
		err.Error() == "403 Forbidden")
}

func isValidationError(err error) bool {
	return err != nil && (
		err.Error() == "400 Bad Request" ||
		err.Error() == "422 Unprocessable Entity")
}

func isServerError(err error) bool {
	return err != nil && (
		err.Error() == "500 Internal Server Error" ||
		err.Error() == "502 Bad Gateway" ||
		err.Error() == "503 Service Unavailable")
}

func getInfoWithFallback(ctx context.Context, primary, fallback slurm.SlurmClient) (*slurm.ServerInfo, error) {
	// Try primary client first
	info, err := primary.GetInfo(ctx)
	if err == nil {
		fmt.Printf("✓ Primary client succeeded\n")
		return info, nil
	}

	fmt.Printf("✗ Primary client failed: %v\n", err)
	fmt.Printf("⚠ Falling back to secondary client...\n")

	// Try fallback client
	info, err = fallback.GetInfo(ctx)
	if err == nil {
		fmt.Printf("✓ Fallback client succeeded\n")
		return info, nil
	}

	fmt.Printf("✗ Fallback client failed: %v\n", err)
	return nil, fmt.Errorf("both primary and fallback failed: %w", err)
}

func submitJobWithDegradation(ctx context.Context, primary, fallback slurm.SlurmClient, jobReq *slurm.JobSubmissionRequest) (uint32, error) {
	// Try primary with advanced features
	jobID, err := primary.SubmitJob(ctx, jobReq)
	if err == nil {
		return jobID, nil
	}

	// Fallback: simplify job requirements
	simplifiedReq := &slurm.JobSubmissionRequest{
		Script:    jobReq.Script,
		Name:      jobReq.Name,
		Partition: "debug", // Use default partition
		Nodes:     1,       // Minimal resources
		Tasks:     1,
		Time:      "00:05:00", // Extended time
	}

	return fallback.SubmitJob(ctx, simplifiedReq)
}

func getCachedInfo() *slurm.ServerInfo {
	// Return cached server information
	return &slurm.ServerInfo{
		Version: "cached-24.05.0",
		Release: "24.05",
	}
}

// Circuit Breaker implementation
type CircuitBreaker struct {
	threshold    int
	timeout      time.Duration
	successCount int
	failureCount int
	state        string
	lastFailTime time.Time
}

func (cb *CircuitBreaker) Call(operation func() error) error {
	// Check if circuit is open
	if cb.state == "OPEN" {
		if time.Since(cb.lastFailTime) > cb.timeout {
			cb.state = "HALF_OPEN"
			fmt.Printf("  Circuit breaker: OPEN → HALF_OPEN\n")
		} else {
			return fmt.Errorf("circuit breaker is OPEN")
		}
	}

	// Execute operation
	err := operation()
	
	if err != nil {
		cb.failureCount++
		cb.lastFailTime = time.Now()
		
		if cb.failureCount >= cb.threshold {
			cb.state = "OPEN"
			fmt.Printf("  Circuit breaker: %s → OPEN (failures: %d)\n", cb.state, cb.failureCount)
		}
		
		return err
	}

	// Success
	cb.successCount++
	if cb.state == "HALF_OPEN" {
		cb.state = "CLOSED"
		cb.failureCount = 0
		fmt.Printf("  Circuit breaker: HALF_OPEN → CLOSED\n")
	}

	return nil
}
