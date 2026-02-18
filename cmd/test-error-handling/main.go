// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/internal/factory"
	"github.com/jontk/slurm-client/pkg/config"
	pkgerrors "github.com/jontk/slurm-client/pkg/errors"
)

// userTokenAuth implements authentication with both username and token headers
type userTokenAuth struct {
	username string
	token    string
}

func (u *userTokenAuth) Authenticate(_ context.Context, req *http.Request) error {
	req.Header.Set("X-SLURM-USER-NAME", u.username)
	req.Header.Set("X-SLURM-USER-TOKEN", u.token)
	return nil
}

func (u *userTokenAuth) Type() string {
	return "user-token"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <test-type>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Test types:\n")
		fmt.Fprintf(os.Stderr, "  invalid-account   - Test with invalid account\n")
		fmt.Fprintf(os.Stderr, "  invalid-partition - Test with invalid partition\n")
		fmt.Fprintf(os.Stderr, "  invalid-qos       - Test with invalid QoS\n")
		fmt.Fprintf(os.Stderr, "  duplicate-account - Test creating duplicate account\n")
		fmt.Fprintf(os.Stderr, "  missing-reservation - Test getting non-existent reservation\n")
		fmt.Fprintf(os.Stderr, "  bad-job-time     - Test with invalid time limit\n")
		fmt.Fprintf(os.Stderr, "  permission       - Test permission denied scenario\n")
		fmt.Fprintf(os.Stderr, "  all              - Run all tests\n")
		os.Exit(1)
	}

	// Get JWT token from environment
	jwtToken := os.Getenv("SLURM_JWT")
	if jwtToken == "" {
		log.Fatal("SLURM_JWT environment variable is required")
	}

	username := os.Getenv("SLURM_USER")
	if username == "" {
		username = "root" // Default username for testing
	}

	testType := os.Args[1]

	// Create client
	client, err := createClient(jwtToken, username)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("Enhanced Error Handling Test")
	fmt.Println("============================")
	fmt.Printf("Testing: %s\n\n", testType)

	// Run tests based on type
	switch testType {
	case "invalid-account":
		testInvalidAccount(client)
	case "invalid-partition":
		testInvalidPartition(client)
	case "invalid-qos":
		testInvalidQoS(client)
	case "duplicate-account":
		testDuplicateAccount(client)
	case "missing-reservation":
		testMissingReservation(client)
	case "bad-job-time":
		testBadJobTime(client)
	case "permission":
		testPermission(client)
	case "all":
		testInvalidAccount(client)
		fmt.Println("\n" + strings.Repeat("-", 50) + "\n")
		testInvalidPartition(client)
		fmt.Println("\n" + strings.Repeat("-", 50) + "\n")
		testInvalidQoS(client)
		fmt.Println("\n" + strings.Repeat("-", 50) + "\n")
		testDuplicateAccount(client)
		fmt.Println("\n" + strings.Repeat("-", 50) + "\n")
		testMissingReservation(client)
		fmt.Println("\n" + strings.Repeat("-", 50) + "\n")
		testBadJobTime(client)
		fmt.Println("\n" + strings.Repeat("-", 50) + "\n")
		testPermission(client)
	default:
		fmt.Fprintf(os.Stderr, "Unknown test type: %s\n", testType)
		os.Exit(1)
	}
}

func createClient(jwtToken string, username string) (types.SlurmClient, error) {
	// Create configuration
	cfg := config.NewDefault()
	cfg.BaseURL = "http://localhost:6820"
	cfg.Debug = false

	// Create JWT authentication provider
	authProvider := &userTokenAuth{
		username: username,
		token:    jwtToken,
	}

	// Create factory with adapter option
	clientFactory, err := factory.NewClientFactory(
		factory.WithConfig(cfg),
		factory.WithAuth(authProvider),
		factory.WithBaseURL(cfg.BaseURL),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create factory: %w", err)
	}

	// Create client with v0.0.43
	client, err := clientFactory.NewClientWithVersion(context.Background(), "v0.0.43")
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}

func testInvalidAccount(client types.SlurmClient) {
	fmt.Println("Test: Submit job with invalid account")
	fmt.Println("Expected: Should get INVALID_ACCOUNT error with enhanced description")

	ctx := context.Background()
	submitJob := &types.JobSubmission{
		Name:       "test-invalid-account",
		Account:    "nonexistent-account-12345", // This account doesn't exist
		Partition:  "normal",
		Script:     "#!/bin/bash\necho 'Test job'\n",
		TimeLimit:  1,
		Nodes:      1,
		WorkingDir: "/tmp",
	}

	_, err := client.Jobs().Submit(ctx, submitJob)
	analyzeError(err)
}

func testInvalidPartition(client types.SlurmClient) {
	fmt.Println("Test: Submit job with invalid partition")
	fmt.Println("Expected: Should get INVALID_PARTITION error with enhanced description")

	ctx := context.Background()
	submitJob := &types.JobSubmission{
		Name:       "test-invalid-partition",
		Account:    "root",
		Partition:  "nonexistent-partition-12345", // This partition doesn't exist
		Script:     "#!/bin/bash\necho 'Test job'\n",
		TimeLimit:  1,
		Nodes:      1,
		WorkingDir: "/tmp",
	}

	_, err := client.Jobs().Submit(ctx, submitJob)
	analyzeError(err)
}

func testInvalidQoS(client types.SlurmClient) {
	fmt.Println("Test: Get non-existent QoS")
	fmt.Println("Expected: Should get QOS_NOT_FOUND error with enhanced description")

	ctx := context.Background()
	_, err := client.QoS().Get(ctx, "nonexistent-qos-12345")
	analyzeError(err)
}

func testDuplicateAccount(client types.SlurmClient) {
	fmt.Println("Test: Create duplicate account")
	fmt.Println("Expected: Should get ACCOUNT_ALREADY_EXISTS error with enhanced description")

	ctx := context.Background()

	// First, create an account
	createReq := &types.AccountCreate{
		Name:        fmt.Sprintf("test-dup-%d", time.Now().Unix()),
		Description: "Test account for error handling",
	}

	_, err := client.Accounts().Create(ctx, createReq)
	if err != nil {
		fmt.Printf("Failed to create initial account: %v\n", err)
		return
	}

	// Try to create the same account again
	_, err = client.Accounts().Create(ctx, createReq)
	analyzeError(err)

	// Clean up
	_ = client.Accounts().Delete(ctx, createReq.Name)
}

func testMissingReservation(client types.SlurmClient) {
	fmt.Println("Test: Get non-existent reservation")
	fmt.Println("Expected: Should get RESERVATION_NOT_FOUND error with enhanced description")

	ctx := context.Background()
	_, err := client.Reservations().Get(ctx, "nonexistent-reservation-12345")
	analyzeError(err)
}

func testBadJobTime(client types.SlurmClient) {
	fmt.Println("Test: Submit job with invalid time limit")
	fmt.Println("Expected: Should get INVALID_TIME_LIMIT error with enhanced description")

	ctx := context.Background()
	submitJob := &types.JobSubmission{
		Name:       "test-bad-time",
		Account:    "root",
		Partition:  "normal",
		Script:     "#!/bin/bash\necho 'Test job'\n",
		TimeLimit:  -1, // Negative time limit
		Nodes:      1,
		WorkingDir: "/tmp",
	}

	_, err := client.Jobs().Submit(ctx, submitJob)
	analyzeError(err)
}

func testPermission(client types.SlurmClient) {
	fmt.Println("Test: Update system QoS (permission test)")
	fmt.Println("Expected: Should get PERMISSION_DENIED error with enhanced description")

	ctx := context.Background()

	// Try to modify a system QoS (usually protected)
	updateReq := &types.QoSUpdate{
		Priority: intPtr(99999),
	}

	err := client.QoS().Update(ctx, "normal", updateReq)
	analyzeError(err)
}

func analyzeError(err error) {
	if err == nil {
		fmt.Println("❌ No error returned (expected an error)")
		return
	}

	fmt.Printf("\n📋 Error Analysis:\n")
	fmt.Printf("  Error Type: %T\n", err)
	fmt.Printf("  Error String: %v\n", err)

	// Check if it's a SLURM error using errors.As
	var slurmErr *pkgerrors.SlurmError
	var apiErr *pkgerrors.SlurmAPIError

	if errors.As(err, &slurmErr) {
		fmt.Printf("\n✅ Enhanced SLURM Error Details:\n")
		fmt.Printf("  Code: %s\n", slurmErr.Code)
		fmt.Printf("  Category: %s\n", slurmErr.Category)
		fmt.Printf("  Message: %s\n", slurmErr.Message)
		if slurmErr.Details != "" {
			fmt.Printf("  Details: %s\n", slurmErr.Details)
		}
		if slurmErr.APIVersion != "" {
			fmt.Printf("  API Version: %s\n", slurmErr.APIVersion)
		}
		if slurmErr.StatusCode != 0 {
			fmt.Printf("  HTTP Status: %d\n", slurmErr.StatusCode)
		}
		fmt.Printf("  Retryable: %v\n", slurmErr.Retryable)
		fmt.Printf("  Timestamp: %s\n", slurmErr.Timestamp.Format(time.RFC3339))
	} else if errors.As(err, &apiErr) {
		fmt.Printf("\n✅ Enhanced SLURM API Error Details:\n")
		fmt.Printf("  Code: %s\n", apiErr.Code)
		fmt.Printf("  Category: %s\n", apiErr.Category)
		fmt.Printf("  Message: %s\n", apiErr.Message)
		if apiErr.Details != "" {
			fmt.Printf("  Details: %s\n", apiErr.Details)
		}
		if apiErr.APIVersion != "" {
			fmt.Printf("  API Version: %s\n", apiErr.APIVersion)
		}
		if apiErr.StatusCode != 0 {
			fmt.Printf("  HTTP Status: %d\n", apiErr.StatusCode)
		}

		// Check for API errors
		if len(apiErr.Errors) > 0 {
			fmt.Printf("\n  API Error Details:\n")
			for i, detail := range apiErr.Errors {
				fmt.Printf("    [%d] Error Number: %d\n", i+1, detail.ErrorNumber)
				fmt.Printf("        Error Code: %s\n", detail.ErrorCode)
				fmt.Printf("        Description: %s\n", detail.Description)
				if detail.Source != "" {
					fmt.Printf("        Source: %s\n", detail.Source)
				}
			}
		}
	} else {
		fmt.Printf("\n❌ Not a SLURM error (no enhanced details available)\n")
		fmt.Printf("  Raw Error Type: %T\n", err)
	}

	// Try to extract error chain
	fmt.Printf("\n🔗 Error Chain:\n")
	printErrorChain(err, 1)
}

func printErrorChain(err error, depth int) {
	if err == nil {
		return
	}

	indent := strings.Repeat("  ", depth)
	fmt.Printf("%s- %v\n", indent, err)

	// Check if error implements Unwrap
	if unwrapper, ok := err.(interface{ Unwrap() error }); ok {
		if wrapped := unwrapper.Unwrap(); wrapped != nil {
			printErrorChain(wrapped, depth+1)
		}
	}
}

func intPtr(i int) *int {
	return &i
}
