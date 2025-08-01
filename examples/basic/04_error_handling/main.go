// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// This example demonstrates comprehensive error handling with the SLURM client.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/types"
)

func main() {
	ctx := context.Background()

	// Example 1: Connection errors
	fmt.Println("=== Example 1: Connection Errors ===")
	handleConnectionErrors(ctx)

	// Example 2: Authentication errors
	fmt.Println("\n=== Example 2: Authentication Errors ===")
	handleAuthErrors(ctx)

	// Example 3: API errors
	fmt.Println("\n=== Example 3: API Errors ===")
	handleAPIErrors(ctx)

	// Example 4: Timeout and context errors
	fmt.Println("\n=== Example 4: Timeout Errors ===")
	handleTimeoutErrors(ctx)

	// Example 5: Retry and recovery
	fmt.Println("\n=== Example 5: Retry and Recovery ===")
	handleRetryAndRecovery(ctx)
}

func handleConnectionErrors(ctx context.Context) {
	// Try to connect to a non-existent server
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://non-existent-server:6820"),
		slurm.WithAuth(auth.NewTokenAuth("token")),
		slurm.WithTimeout(5*time.Second),
	)
	
	if err != nil {
		// Check for specific error types
		var netErr net.Error
		if errors.As(err, &netErr) {
			if netErr.Timeout() {
				fmt.Println("Connection timed out")
			} else {
				fmt.Printf("Network error: %v\n", netErr)
			}
		} else {
			fmt.Printf("Connection error: %v\n", err)
		}
		return
	}
	defer client.Close()
}

func handleAuthErrors(ctx context.Context) {
	// Try with invalid credentials
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(auth.NewTokenAuth("invalid-token")),
	)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}
	defer client.Close()

	// Try to list jobs - this should fail with auth error
	_, err = client.Jobs().List(ctx, nil)
	if err != nil {
		var apiErr *types.APIError
		if errors.As(err, &apiErr) {
			fmt.Printf("API Error: %s (Code: %s)\n", apiErr.Message, apiErr.ErrorCode)
			
			// Check specific error types
			if apiErr.IsAuthError() {
				fmt.Println("This is an authentication error - check your token")
			} else if apiErr.IsPermissionError() {
				fmt.Println("This is a permission error - check your access rights")
			}
			
			// Get HTTP status code
			if apiErr.StatusCode == 401 {
				fmt.Println("HTTP 401 Unauthorized")
			} else if apiErr.StatusCode == 403 {
				fmt.Println("HTTP 403 Forbidden")
			}
		} else {
			fmt.Printf("Unexpected error: %v\n", err)
		}
	}
}

func handleAPIErrors(ctx context.Context) {
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(auth.NewTokenAuth("your-token")),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Try to get a non-existent job
	jobID := 999999
	job, err := client.Jobs().Get(ctx, jobID)
	if err != nil {
		var apiErr *types.APIError
		if errors.As(err, &apiErr) {
			switch {
			case apiErr.IsNotFoundError():
				fmt.Printf("Job %d not found\n", jobID)
			case apiErr.IsRateLimitError():
				fmt.Printf("Rate limit exceeded. Retry after: %v\n", apiErr.RetryAfter)
			case apiErr.IsServerError():
				fmt.Printf("Server error: %s\n", apiErr.Message)
			default:
				fmt.Printf("API error: %s\n", apiErr.Message)
			}
		}
		return
	}
	fmt.Printf("Job found: %v\n", job)
}

func handleTimeoutErrors(ctx context.Context) {
	// Create a client with short timeout
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(auth.NewTokenAuth("your-token")),
		slurm.WithTimeout(100*time.Millisecond), // Very short timeout
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// This might timeout
	jobs, err := client.Jobs().List(timeoutCtx, nil)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("Operation timed out")
		} else if errors.Is(err, context.Canceled) {
			fmt.Println("Operation was canceled")
		} else {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}
	fmt.Printf("Retrieved %d jobs\n", len(jobs))
}

func handleRetryAndRecovery(ctx context.Context) {
	// Create client with retry configuration
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(auth.NewTokenAuth("your-token")),
		slurm.WithMaxRetries(3),
		slurm.WithRetryWaitMin(1*time.Second),
		slurm.WithRetryWaitMax(10*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Implement custom retry logic
	var jobs []types.Job
	maxAttempts := 3
	
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		jobs, err = client.Jobs().List(ctx, nil)
		if err == nil {
			fmt.Printf("Success on attempt %d\n", attempt)
			break
		}

		var apiErr *types.APIError
		if errors.As(err, &apiErr) {
			// Don't retry on client errors
			if apiErr.StatusCode >= 400 && apiErr.StatusCode < 500 {
				fmt.Printf("Client error (not retrying): %v\n", err)
				break
			}
			
			// Check if error is retryable
			if apiErr.IsServerError() || apiErr.IsRateLimitError() {
				if attempt < maxAttempts {
					backoff := time.Duration(attempt) * time.Second
					fmt.Printf("Attempt %d failed: %v. Retrying in %v...\n", 
						attempt, err, backoff)
					time.Sleep(backoff)
					continue
				}
			}
		}

		fmt.Printf("Failed after %d attempts: %v\n", attempt, err)
		break
	}

	if err == nil {
		fmt.Printf("Successfully retrieved %d jobs\n", len(jobs))
	}
}

// Example of a robust operation wrapper
func robustOperation[T any](
	ctx context.Context,
	operation func() (T, error),
	operationName string,
) (T, error) {
	var result T
	var err error

	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Try the operation
	result, err = operation()
	if err == nil {
		return result, nil
	}

	// Log the error
	log.Printf("%s failed: %v", operationName, err)

	// Analyze error and decide on action
	var apiErr *types.APIError
	if errors.As(err, &apiErr) {
		switch {
		case apiErr.IsAuthError():
			// Could trigger re-authentication here
			return result, fmt.Errorf("authentication required: %w", err)
		case apiErr.IsRateLimitError():
			// Could implement backoff and retry
			return result, fmt.Errorf("rate limited: %w", err)
		case apiErr.IsServerError():
			// Could retry with exponential backoff
			return result, fmt.Errorf("server error: %w", err)
		default:
			return result, err
		}
	}

	// Check for context errors
	if errors.Is(err, context.DeadlineExceeded) {
		return result, fmt.Errorf("operation timed out: %w", err)
	}
	if errors.Is(err, context.Canceled) {
		return result, fmt.Errorf("operation canceled: %w", err)
	}

	return result, err
}