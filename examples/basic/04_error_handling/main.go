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
	slurmErrors "github.com/jontk/slurm-client/pkg/errors"
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
		var slurmErr *slurmErrors.SlurmError
		if errors.As(err, &slurmErr) {
			fmt.Printf("SLURM Error: %s (Code: %s)\n", slurmErr.Message, slurmErr.Code)

			// Check specific error types
			switch slurmErr.Code {
			case slurmErrors.ErrorCodeInvalidCredentials:
				fmt.Println("This is an authentication error - check your token")
			case slurmErrors.ErrorCodePermissionDenied:
				fmt.Println("This is a permission error - check your access rights")
			case slurmErrors.ErrorCodeUnauthorized:
				fmt.Println("Unauthorized access")
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
	jobID := "999999"
	job, err := client.Jobs().Get(ctx, jobID)
	if err != nil {
		var slurmErr *slurmErrors.SlurmError
		if errors.As(err, &slurmErr) {
			switch slurmErr.Code {
			case slurmErrors.ErrorCodeResourceNotFound:
				fmt.Printf("Job %s not found\n", jobID)
			case slurmErrors.ErrorCodeRateLimited:
				fmt.Printf("Rate limit exceeded\n")
			case slurmErrors.ErrorCodeServerInternal:
				fmt.Printf("Server error: %s\n", slurmErr.Message)
			default:
				fmt.Printf("API error: %s\n", slurmErr.Message)
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
	jobList, err := client.Jobs().List(timeoutCtx, nil)
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
	fmt.Printf("Retrieved %d jobs\n", jobList.Total)
}

func handleRetryAndRecovery(ctx context.Context) {
	// Create client with retry configuration
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(auth.NewTokenAuth("your-token")),
		slurm.WithMaxRetries(3),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Implement custom retry logic
	var jobList *slurm.JobList
	maxAttempts := 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		jobList, err = client.Jobs().List(ctx, nil)
		if err == nil {
			fmt.Printf("Success on attempt %d\n", attempt)
			break
		}

		var slurmErr *slurmErrors.SlurmError
		if errors.As(err, &slurmErr) {
			// Don't retry on client errors
			if slurmErr.Code == slurmErrors.ErrorCodeInvalidRequest ||
				slurmErr.Code == slurmErrors.ErrorCodeInvalidCredentials ||
				slurmErr.Code == slurmErrors.ErrorCodePermissionDenied {
				fmt.Printf("Client error (not retrying): %v\n", err)
				break
			}

			// Check if error is retryable
			if slurmErr.Code == slurmErrors.ErrorCodeServerInternal ||
				slurmErr.Code == slurmErrors.ErrorCodeRateLimited {
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
		fmt.Printf("Successfully retrieved %d jobs\n", jobList.Total)
	}
}

// Example of a robust operation wrapper - commented out as unused
/*
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
	var slurmErr *slurmErrors.SlurmError
	if errors.As(err, &slurmErr) {
		switch slurmErr.Code {
		case slurmErrors.ErrorCodeInvalidCredentials:
			// Could trigger re-authentication here
			return result, fmt.Errorf("authentication required: %w", err)
		case slurmErrors.ErrorCodeRateLimited:
			// Could implement backoff and retry
			return result, fmt.Errorf("rate limited: %w", err)
		case slurmErrors.ErrorCodeServerInternal:
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
*/
