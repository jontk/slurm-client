// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/internal/factory"
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

	// Test various operations
	ctx := context.Background()

	// Test 1: List Jobs
	fmt.Println("\n=== Testing List Jobs ===")
	jobs, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{
		Limit: 5,
	})
	if err != nil {
		log.Printf("Failed to list jobs: %v", err)
	} else {
		fmt.Printf("Found %d jobs (showing up to 5)\n", jobs.Total)
		for i, job := range jobs.Jobs {
			fmt.Printf("  [%d] Job %s: %s (State: %s)\n", i+1, job.ID, job.Name, job.State)
		}
	}

	// Test 2: Submit a test job
	fmt.Println("\n=== Testing Job Submission ===")
	testJob := &interfaces.JobSubmission{
		Name:       fmt.Sprintf("adapter-test-%s-%d", version, time.Now().Unix()),
		Partition:  "normal",
		Script:     "#!/bin/bash\necho 'Hello from adapter test'\nsleep 10\necho 'Done'",
		TimeLimit:  1, // 1 minute
		Nodes:      1,
		CPUs:       1,
		WorkingDir: "/tmp",
		Environment: map[string]string{
			"PATH":     "/usr/bin:/bin",
			"USER":     "root",
			"HOME":     "/tmp",
			"TEST_VAR": "adapter_test",
		},
	}

	submitResp, err := client.Jobs().Submit(ctx, testJob)
	if err != nil {
		log.Printf("Failed to submit job: %v", err)
	} else {
		fmt.Printf("Successfully submitted job with ID: %s\n", submitResp.JobID)

		// Test 3: Get the submitted job
		fmt.Println("\n=== Testing Get Job ===")
		job, err := client.Jobs().Get(ctx, submitResp.JobID)
		if err != nil {
			log.Printf("Failed to get job %s: %v", submitResp.JobID, err)
		} else {
			fmt.Printf("Retrieved job %s:\n", job.ID)
			fmt.Printf("  Name: %s\n", job.Name)
			fmt.Printf("  State: %s\n", job.State)
			fmt.Printf("  Partition: %s\n", job.Partition)
			fmt.Printf("  Working Directory: %s\n", job.WorkingDir)
		}

		// Test 4: Cancel the job
		fmt.Println("\n=== Testing Cancel Job ===")
		err = client.Jobs().Cancel(ctx, submitResp.JobID)
		if err != nil {
			log.Printf("Failed to cancel job %s: %v", submitResp.JobID, err)
		} else {
			fmt.Printf("Successfully cancelled job %s\n", submitResp.JobID)
		}
	}

	// Test 5: List Nodes
	fmt.Println("\n=== Testing List Nodes ===")
	nodes, err := client.Nodes().List(ctx, &interfaces.ListNodesOptions{
		Limit: 5,
	})
	if err != nil {
		log.Printf("Failed to list nodes: %v", err)
	} else {
		fmt.Printf("Found %d nodes (showing up to 5)\n", nodes.Total)
		for i, node := range nodes.Nodes {
			fmt.Printf("  [%d] Node %s: State=%s, CPUs=%d, Memory=%dMB\n",
				i+1, node.Name, node.State, node.CPUs, node.Memory)
		}
	}

	// Test 6: List Partitions
	fmt.Println("\n=== Testing List Partitions ===")
	partitions, err := client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{
		Limit: 5,
	})
	if err != nil {
		log.Printf("Failed to list partitions: %v", err)
	} else {
		fmt.Printf("Found %d partitions (showing up to 5)\n", partitions.Total)
		for i, partition := range partitions.Partitions {
			fmt.Printf("  [%d] Partition %s: State=%s, Nodes=%d total (%d available)\n",
				i+1, partition.Name, partition.State, partition.TotalNodes, partition.TotalNodes)
		}
	}

	fmt.Println("\n=== All tests completed ===")
}
