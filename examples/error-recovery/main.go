// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	stderrors "errors"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	slurm "github.com/jontk/slurm-client"
	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
	"github.com/jontk/slurm-client/pkg/errors"
	"github.com/jontk/slurm-client/pkg/retry"
)

// Example: Error recovery and resilience patterns
func main() {
	// Create configuration
	cfg := config.NewDefault()
	cfg.BaseURL = "https://cluster.example.com:6820"

	// Create authentication
	authProvider := auth.NewTokenAuth("your-jwt-token")

	ctx := context.Background()

	// Example 1: Structured error handling
	fmt.Println("=== Structured Error Handling ===")
	demonstrateStructuredErrorHandling(ctx, cfg, authProvider)

	// Example 2: Retry strategies
	fmt.Println("\n=== Retry Strategies ===")
	demonstrateRetryStrategies(ctx, cfg, authProvider)

	// Example 3: Circuit breaker pattern
	fmt.Println("\n=== Circuit Breaker Pattern ===")
	demonstrateCircuitBreaker()

	// Example 4: Graceful degradation
	fmt.Println("\n=== Graceful Degradation ===")
	demonstrateGracefulDegradation(ctx, cfg, authProvider)

	// Example 5: Error recovery workflows
	fmt.Println("\n=== Error Recovery Workflows ===")
	demonstrateErrorRecoveryWorkflows(ctx, cfg, authProvider)
}

// demonstrateStructuredErrorHandling shows comprehensive error handling
func demonstrateStructuredErrorHandling(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	client, err := slurm.NewClient(ctx,
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
	)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
	defer client.Close()

	// Example: Submit job with comprehensive error handling
	job := &slurm.JobSubmission{
		Name:      "error-test-job",
		Command:   "echo 'Testing error handling'",
		Partition: "invalid-partition", // Intentional error
		CPUs:      1000,                // Excessive resources
		Memory:    1024 * 1024 * 1024,  // 1TB - likely to fail
		TimeLimit: 10080,               // 1 week - may exceed limits
	}

	resp, err := client.Jobs().Submit(ctx, job)
	if err != nil {
		handleJobSubmissionError(err, job)
		return
	}

	fmt.Printf("Job submitted successfully: %d\n", resp.JobId)
}

// handleJobSubmissionError demonstrates comprehensive error handling
func handleJobSubmissionError(err error, job *slurm.JobSubmission) {
	// Check if it's a SLURM error
	var slurmErr *errors.SlurmError
	if stderrors.As(err, &slurmErr) {
		fmt.Printf("SLURM Error Details:\n")
		fmt.Printf("  Code: %s\n", slurmErr.Code)
		fmt.Printf("  Category: %s\n", slurmErr.Category)
		fmt.Printf("  Message: %s\n", slurmErr.Message)
		fmt.Printf("  Details: %s\n", slurmErr.Details)
		fmt.Printf("  Retryable: %t\n", slurmErr.IsRetryable())

		// Handle specific error codes
		switch slurmErr.Code {
		case errors.ErrorCodeValidationFailed:
			fmt.Println("\nValidation Error - Checking specific issues:")
			// Check partition
			if slurmErr.Details != "" {
				fmt.Printf("  %s\n", slurmErr.Details)
			}
			fmt.Println("  Suggested fixes:")
			fmt.Println("    - Verify partition name is correct")
			fmt.Println("    - Check resource requirements are within limits")
			fmt.Println("    - Ensure time limit is reasonable")

		case errors.ErrorCodeResourceExhausted:
			fmt.Println("\nResource Exhaustion - Analyzing requirements:")
			fmt.Printf("  Requested: CPUs=%d, Memory=%dGB\n",
				job.CPUs, job.Memory/(1024*1024*1024))
			fmt.Println("  Suggested actions:")
			fmt.Println("    - Reduce resource requirements")
			fmt.Println("    - Check cluster capacity")
			fmt.Println("    - Try a different partition")
			fmt.Println("    - Submit during off-peak hours")

		case errors.ErrorCodeUnauthorized:
			fmt.Println("\nAuthentication Failed:")
			fmt.Println("  - Check your authentication token")
			fmt.Println("  - Verify token hasn't expired")
			fmt.Println("  - Ensure you have submit permissions")

		case errors.ErrorCodeRateLimited:
			fmt.Println("\nRate Limited:")
			fmt.Println("  - Too many requests in a short time")
			fmt.Println("  - Wait before retrying")
			fmt.Println("  - Consider implementing rate limiting")

		default:
			fmt.Printf("\nUnhandled error code: %s\n", slurmErr.Code)
		}

		// Check error category for broader handling
		switch slurmErr.Category {
		case errors.CategoryClient:
			fmt.Println("\nClient-side error - check your request")
		case errors.CategoryServer:
			fmt.Println("\nServer-side error - may be temporary")
		case errors.CategoryNetwork:
			fmt.Println("\nNetwork error - check connectivity")
		default:
			fmt.Printf("\nOther error category: %s\n", slurmErr.Category)
		}

		return
	}

	// Handle other error types
	if errors.IsNetworkError(err) {
		fmt.Println("Network Error Detected:")
		fmt.Println("  - Check network connectivity")
		fmt.Println("  - Verify SLURM REST API URL")
		fmt.Println("  - Check firewall settings")
		return
	}

	if errors.IsAuthenticationError(err) {
		fmt.Println("Authentication Error:")
		fmt.Println("  - Verify credentials")
		fmt.Println("  - Check authentication method")
		return
	}

	// Generic error
	fmt.Printf("Unexpected error: %v\n", err)
}

// demonstrateRetryStrategies shows different retry patterns
func demonstrateRetryStrategies(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	// Strategy 1: Exponential backoff with jitter
	fmt.Println("1. Exponential Backoff with Jitter:")

	exponentialRetry := retry.NewHTTPExponentialBackoff().
		WithMaxRetries(5).
		WithMinWaitTime(100 * time.Millisecond).
		WithMaxWaitTime(10 * time.Second).
		WithJitter(true)

	client1, err := slurm.NewClient(ctx,
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
		slurm.WithRetryPolicy(exponentialRetry),
	)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
	defer client1.Close()

	// Test with operation that might fail transiently
	testRetryableOperation(ctx, client1, "Exponential Backoff")

	// Strategy 2: Linear backoff for predictable delays
	fmt.Println("\n2. Linear Backoff:")

	// Create a simple retry policy with fixed delay
	// LinearBackoff is not available in the current implementation
	// Using HTTPExponentialBackoff with minimal settings
	linearRetry := retry.NewHTTPExponentialBackoff().
		WithMaxRetries(3).
		WithMinWaitTime(1 * time.Second).
		WithMaxWaitTime(1 * time.Second)

	client2, err := slurm.NewClient(ctx,
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
		slurm.WithRetryPolicy(linearRetry),
	)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
	defer client2.Close()

	testRetryableOperation(ctx, client2, "Linear Backoff")

	// Strategy 3: Custom retry with error inspection
	fmt.Println("\n3. Custom Retry Logic:")

	customRetry := &customRetryPolicy{
		maxRetries: 3,
		shouldRetry: func(err error, attempt int) bool {
			// Only retry specific errors
			var slurmErr *errors.SlurmError
			if stderrors.As(err, &slurmErr) {
				// Retry server errors and rate limits
				if slurmErr.Category == errors.CategoryServer ||
					slurmErr.Code == errors.ErrorCodeRateLimited {
					fmt.Printf("  Attempt %d: Retrying %s error\n", attempt, slurmErr.Code)
					return true
				}
			}
			// Don't retry client errors
			return false
		},
	}

	client3, err := slurm.NewClient(ctx,
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
		slurm.WithRetryPolicy(customRetry),
	)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
	defer client3.Close()

	testRetryableOperation(ctx, client3, "Custom Retry")
}

// demonstrateCircuitBreaker shows circuit breaker pattern
func demonstrateCircuitBreaker() {
	// Create circuit breaker
	breaker := &circuitBreaker{
		failureThreshold: 3,
		recoveryTimeout:  30 * time.Second,
		state:            "closed",
	}

	fmt.Printf("Circuit breaker initial state: %s\n", breaker.state)

	// Simulate requests with circuit breaker
	for i := range 10 {
		fmt.Printf("\nRequest %d:\n", i+1)

		if breaker.isOpen() {
			fmt.Println("  Circuit OPEN - failing fast")
			continue
		}

		// Try the operation
		err := simulateOperation(i < 5) // First 5 fail, rest succeed

		if err != nil {
			breaker.recordFailure()
			fmt.Printf("  Failed: %v\n", err)
			fmt.Printf("  Circuit state: %s (failures: %d/%d)\n",
				breaker.state, breaker.failures, breaker.failureThreshold)
		} else {
			breaker.recordSuccess()
			fmt.Println("  Success!")
			fmt.Printf("  Circuit state: %s\n", breaker.state)
		}

		time.Sleep(100 * time.Millisecond)
	}

	// Show recovery
	fmt.Println("\nWaiting for circuit recovery...")
	time.Sleep(2 * time.Second)

	if breaker.state == circuitStateHalfOpen {
		fmt.Println("Circuit is HALF-OPEN - testing with single request")
		err := simulateOperation(false) // Success
		if err == nil {
			breaker.recordSuccess()
			fmt.Printf("Recovery successful! Circuit state: %s\n", breaker.state)
		}
	}
}

// demonstrateGracefulDegradation shows graceful degradation patterns
func demonstrateGracefulDegradation(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	client, err := slurm.NewClient(ctx,
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
	)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
	defer client.Close()

	// Primary operation with fallbacks
	jobID := "12345"

	// Level 1: Try to get detailed job info
	fmt.Println("Attempting primary operation (detailed job info)...")
	job, err := client.Jobs().Get(ctx, jobID)
	if err != nil {
		fmt.Printf("Primary failed: %v\n", err)

		// Level 2: Fall back to job list with filtering
		fmt.Println("\nFalling back to job list...")
		jobs, err := client.Jobs().List(ctx, &slurm.ListJobsOptions{
			UserID: "current-user",
			Limit:  10,
		})
		if err != nil {
			fmt.Printf("List fallback failed: %v\n", err)

			// Level 3: Fall back to basic connectivity check
			fmt.Println("\nFalling back to connectivity check...")
			err = client.Info().Ping(ctx)
			if err != nil {
				fmt.Printf("Connectivity check failed: %v\n", err)

				// Level 4: Use cached data if available
				fmt.Println("\nFalling back to cached data...")
				cachedJob := getCachedJobData(jobID)
				if cachedJob != nil {
					fmt.Printf("Using cached data (age: %v)\n",
						time.Since(cachedJob.CachedAt))
					displayJobInfo(cachedJob.Job)
				} else {
					fmt.Println("No cached data available - degraded to offline mode")
				}
			} else {
				fmt.Println("Connectivity OK - server may be overloaded")
				fmt.Println("Suggested action: Retry later with backoff")
			}
		} else {
			fmt.Printf("Retrieved %d jobs from list\n", len(jobs.Jobs))
			// Look for our job in the list
			for _, j := range jobs.Jobs {
				// Convert JobID (int32) to string for comparison
				if j.JobID != nil && fmt.Sprintf("%d", *j.JobID) == jobID {
					fmt.Println("Found job in list!")
					displayJobInfo(&j)
					break
				}
			}
		}
	} else {
		fmt.Println("Primary operation successful!")
		displayJobInfo(job)
	}
}

// demonstrateErrorRecoveryWorkflows shows complex error recovery
func demonstrateErrorRecoveryWorkflows(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	client, err := slurm.NewClient(ctx,
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
	)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
	defer client.Close()

	// Workflow: Submit job with automatic recovery
	fmt.Println("Job submission with automatic recovery:")

	job := &slurm.JobSubmission{
		Name:      "recovery-test",
		Command:   "python process.py",
		Partition: "compute",
		CPUs:      8,
		Memory:    16384,
		TimeLimit: 60,
	}

	// Attempt submission with recovery logic
	var jobID string
	maxAttempts := 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("\nAttempt %d/%d:\n", attempt, maxAttempts)

		resp, err := client.Jobs().Submit(ctx, job)
		if err == nil {
			jobID = fmt.Sprintf("%d", resp.JobId)
			fmt.Printf("  Success! Job ID: %s\n", jobID)
			break
		}

		// Analyze error and adapt
		var slurmErr *errors.SlurmError
		if stderrors.As(err, &slurmErr) {
			fmt.Printf("  Failed: %s - %s\n", slurmErr.Code, slurmErr.Message)

			// Adapt based on error
			switch slurmErr.Code {
			case errors.ErrorCodeResourceExhausted:
				// Reduce requirements
				job.CPUs /= 2
				job.Memory /= 2
				fmt.Printf("  Adapting: Reduced to CPUs=%d, Memory=%dGB\n",
					job.CPUs, job.Memory/(1024*1024))

			case errors.ErrorCodeValidationFailed:
				// Try different partition
				if job.Partition == "compute" {
					job.Partition = "shared"
				} else {
					job.Partition = "compute"
				}
				fmt.Printf("  Adapting: Changed partition to '%s'\n", job.Partition)

			case errors.ErrorCodeRateLimited:
				// Wait and retry
				fmt.Println("  Adapting: Waiting 5 seconds due to rate limit...")
				time.Sleep(5 * time.Second)

			default:
				if !slurmErr.IsRetryable() {
					fmt.Println("  Error is not retryable - giving up")
					return
				}
			}
		}

		if attempt < maxAttempts {
			backoff := time.Duration(attempt) * time.Second
			fmt.Printf("  Waiting %v before retry...\n", backoff)
			time.Sleep(backoff)
		}
	}

	if jobID == "" {
		fmt.Println("\nAll attempts failed - initiating fallback workflow")

		// Fallback: Submit to local queue for later processing
		submitToLocalQueue(job)
		fmt.Println("Job queued locally for later submission")
		return
	}

	// Monitor job with error recovery
	fmt.Printf("\nMonitoring job %s with error recovery:\n", jobID)

	monitoringAttempts := 0
	for {
		monitoringAttempts++

		jobInfo, err := client.Jobs().Get(ctx, jobID)
		if err != nil {
			fmt.Printf("Monitoring error (attempt %d): %v\n", monitoringAttempts, err)

			if monitoringAttempts >= 3 {
				fmt.Println("Monitoring failed - job may still be running")
				break
			}

			time.Sleep(2 * time.Second)
			continue
		}

		// Get job state from slice (if available)
		var jobState string
		if len(jobInfo.JobState) > 0 {
			jobState = string(jobInfo.JobState[0])
		}
		fmt.Printf("Job state: %s\n", jobState)

		if jobState == "COMPLETED" || jobState == "FAILED" {
			if jobState == "FAILED" && jobInfo.ExitCode != nil && jobInfo.ExitCode.ReturnCode != nil {
				fmt.Printf("Job failed with exit code %d\n", *jobInfo.ExitCode.ReturnCode)

				// Automatic resubmission for specific exit codes
				if shouldResubmit(int(*jobInfo.ExitCode.ReturnCode)) {
					fmt.Println("Automatically resubmitting job...")
					// Resubmit logic here
				}
			}
			break
		}

		time.Sleep(5 * time.Second)
	}
}

// Helper types and functions

const (
	circuitStateClosed   = "closed"
	circuitStateOpen     = "open"
	circuitStateHalfOpen = "half-open"
)

type customRetryPolicy struct {
	maxRetries  int
	shouldRetry func(error, int) bool
}

func (c *customRetryPolicy) ShouldRetry(_ context.Context, _ *http.Response, err error, attempt int) bool {
	if attempt >= c.maxRetries {
		return false
	}
	return c.shouldRetry(err, attempt)
}

func (c *customRetryPolicy) WaitTime(attempt int) time.Duration {
	return time.Duration(attempt) * time.Second
}

func (c *customRetryPolicy) MaxRetries() int {
	return c.maxRetries
}

type circuitBreaker struct {
	state            string // closed, open, half-open
	failures         int
	failureThreshold int
	lastFailureTime  time.Time
	recoveryTimeout  time.Duration
}

func (cb *circuitBreaker) isOpen() bool {
	if cb.state == "open" {
		// Check if recovery timeout has passed
		if time.Since(cb.lastFailureTime) > cb.recoveryTimeout {
			cb.state = circuitStateHalfOpen
			cb.failures = 0
			return false
		}
		return true
	}
	return false
}

func (cb *circuitBreaker) recordFailure() {
	cb.failures++
	cb.lastFailureTime = time.Now()

	if cb.failures >= cb.failureThreshold {
		cb.state = "open"
	}
}

func (cb *circuitBreaker) recordSuccess() {
	if cb.state == circuitStateHalfOpen {
		cb.state = "closed"
	}
	cb.failures = 0
}

type cachedJob struct {
	Job      *slurm.Job
	CachedAt time.Time
}

func getCachedJobData(jobID string) *cachedJob {
	// Simulate cached data retrieval
	// Convert jobID string to int32
	var jobIDInt int32
	_, _ = fmt.Sscanf(jobID, "%d", &jobIDInt)

	name := "cached-job"
	partition := "compute"
	cpus := uint32(4)
	memPerNode := uint64(8192)
	return &cachedJob{
		Job: &slurm.Job{
			JobID:         &jobIDInt,
			Name:          &name,
			JobState:      []types.JobState{types.JobStatePending},
			Partition:     &partition,
			CPUs:          &cpus,
			MemoryPerNode: &memPerNode,
		},
		CachedAt: time.Now().Add(-5 * time.Minute),
	}
}

func displayJobInfo(job *slurm.Job) {
	fmt.Printf("Job Details:\n")
	jobID := ""
	if job.JobID != nil {
		jobID = fmt.Sprintf("%d", *job.JobID)
	}
	fmt.Printf("  ID: %s\n", jobID)

	name := ""
	if job.Name != nil {
		name = *job.Name
	}
	fmt.Printf("  Name: %s\n", name)

	state := ""
	if len(job.JobState) > 0 {
		state = string(job.JobState[0])
	}
	fmt.Printf("  State: %s\n", state)

	partition := ""
	if job.Partition != nil {
		partition = *job.Partition
	}
	fmt.Printf("  Partition: %s\n", partition)

	cpus := uint32(0)
	if job.CPUs != nil {
		cpus = *job.CPUs
	}
	fmt.Printf("  CPUs: %d\n", cpus)

	memPerNode := uint64(0)
	if job.MemoryPerNode != nil {
		memPerNode = *job.MemoryPerNode
	}
	fmt.Printf("  Memory: %dMB\n", memPerNode)
}

func simulateOperation(shouldFail bool) error {
	if shouldFail {
		return fmt.Errorf("simulated failure")
	}
	return nil
}

func submitToLocalQueue(job *slurm.JobSubmission) {
	// Simulate local queue submission
	fmt.Printf("Queuing job '%s' locally\n", job.Name)
}

func shouldResubmit(exitCode int) bool {
	// Define exit codes that warrant automatic resubmission
	resubmitCodes := []int{137, 143} // SIGKILL, SIGTERM - possibly OOM or timeout
	for _, code := range resubmitCodes {
		if exitCode == code {
			return true
		}
	}
	return false
}

func testRetryableOperation(ctx context.Context, client slurm.SlurmClient, strategyName string) {
	fmt.Printf("Testing %s:\n", strategyName)

	start := time.Now()
	_, err := client.Jobs().Get(ctx, "nonexistent-job")
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("  Operation failed after retries: %v\n", err)
		fmt.Printf("  Total time with retries: %v\n", elapsed)
	} else {
		fmt.Printf("  Operation succeeded after retries\n")
		fmt.Printf("  Total time: %v\n", elapsed)
	}
}
