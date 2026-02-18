// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/internal/factory"
	"github.com/jontk/slurm-client/pkg/config"
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
		fmt.Fprintf(os.Stderr, "Usage: %s <version>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Supported versions: v0.0.40, v0.0.41, v0.0.42, v0.0.43\n")
		os.Exit(1)
	}

	version := os.Args[1]

	// Get JWT token from environment
	jwtToken := os.Getenv("SLURM_JWT")
	if jwtToken == "" {
		log.Fatal("SLURM_JWT environment variable is required")
	}

	username := os.Getenv("SLURM_USER")
	if username == "" {
		username = "root" // Default username for testing
	}

	// Create configuration
	cfg := config.NewDefault()
	cfg.BaseURL = "http://localhost:6820"
	cfg.Debug = true

	// Create JWT authentication provider
	authProvider := &userTokenAuth{
		username: username,
		token:    jwtToken,
	}

	// Create factory
	clientFactory, err := factory.NewClientFactory(
		factory.WithConfig(cfg),
		factory.WithAuth(authProvider),
		factory.WithBaseURL(cfg.BaseURL),
	)
	if err != nil {
		log.Fatalf("Failed to create factory: %v", err)
	}

	// Create client with specific version
	client, err := clientFactory.NewClientWithVersion(context.Background(), version)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Printf("Successfully created %s client using adapters!\n", version)

	// Test job submission
	ctx := context.Background()

	fmt.Println("\n=== Testing Job Submission with Adapter Client ===")
	testJob := &types.JobSubmission{
		Name:       fmt.Sprintf("adapter-test-%s-%d", version, time.Now().Unix()),
		Partition:  "normal",
		Script:     "#!/bin/bash\necho 'Hello from adapter test'\necho 'PATH=$PATH'\necho 'TEST_VAR=$TEST_VAR'\nsleep 10\necho 'Done'",
		TimeLimit:  1, // 1 minute
		Nodes:      1,
		WorkingDir: "/tmp", // Add working directory
		Environment: map[string]string{
			"PATH":     "/usr/bin:/bin",
			"USER":     "root",
			"HOME":     "/tmp",
			"TEST_VAR": "adapter_test",
		},
	}

	submitResp, err := client.Jobs().Submit(ctx, testJob)
	if err != nil {
		log.Fatalf("Failed to submit job: %v", err)
	}

	fmt.Printf("Successfully submitted job with ID: %d\n", submitResp.JobId)

	// Cancel the job
	fmt.Printf("\nCancelling job %d...\n", submitResp.JobId)
	err = client.Jobs().Cancel(ctx, fmt.Sprintf("%d", submitResp.JobId))
	if err != nil {
		log.Printf("Failed to cancel job %d: %v", submitResp.JobId, err)
	} else {
		fmt.Printf("Successfully cancelled job %d\n", submitResp.JobId)
	}

	fmt.Println("\n=== Test completed successfully! ===")
}
