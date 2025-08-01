// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jontk/slurm-client/internal/factory"
	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

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

	// Create configuration
	cfg := config.NewDefault()
	cfg.BaseURL = "http://rocky9.ar.jontk.com:6820"
	cfg.Debug = true

	// Create JWT authentication provider
	authProvider := auth.NewTokenAuth(jwtToken)

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
	
	fmt.Println("\n=== Testing Job Submission with Wrapper Client ===")
	testJob := &interfaces.JobSubmission{
		Name:      fmt.Sprintf("wrapper-test-%s-%d", version, time.Now().Unix()),
		Partition: "normal",
		Script:    "#!/bin/bash\necho 'Hello from wrapper test'\necho 'PATH=$PATH'\necho 'TEST_VAR=$TEST_VAR'\nsleep 10\necho 'Done'",
		TimeLimit: 1, // 1 minute
		Nodes:     1,
		WorkingDir: "/tmp",  // Add working directory
		Environment: map[string]string{
			"PATH":     "/usr/bin:/bin",
			"USER":     "root",
			"HOME":     "/tmp",
			"TEST_VAR": "wrapper_test",
		},
	}

	submitResp, err := client.Jobs().Submit(ctx, testJob)
	if err != nil {
		log.Fatalf("Failed to submit job: %v", err)
	}
	
	fmt.Printf("Successfully submitted job with ID: %s\n", submitResp.JobID)

	// Cancel the job
	fmt.Printf("\nCancelling job %s...\n", submitResp.JobID)
	err = client.Jobs().Cancel(ctx, submitResp.JobID)
	if err != nil {
		log.Printf("Failed to cancel job %s: %v", submitResp.JobID, err)
	} else {
		fmt.Printf("Successfully cancelled job %s\n", submitResp.JobID)
	}

	fmt.Println("\n=== Test completed successfully! ===")
}
