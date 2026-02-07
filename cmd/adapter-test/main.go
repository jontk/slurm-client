// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	types "github.com/jontk/slurm-client/api"
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
	cfg.BaseURL = "http://localhost:6820"
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
	jobs, err := client.Jobs().List(ctx, &types.ListJobsOptions{
		Limit: 5,
	})
	if err != nil {
		log.Printf("Failed to list jobs: %v", err)
	} else {
		fmt.Printf("Found %d jobs (showing up to 5)\n", jobs.Total)
		for i, job := range jobs.Jobs {
			fmt.Printf("  [%d] Job %v: %v (State: %v)\n", i+1, job.JobID, job.Name, job.JobState)
		}
	}

	// Test 2: Submit a test job
	fmt.Println("\n=== Testing Job Submission ===")
	testJob := &types.JobSubmission{
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
		jobIDStr := strconv.Itoa(int(submitResp.JobId))
		fmt.Printf("Successfully submitted job with ID: %s\n", jobIDStr)

		// Test 3: Get the submitted job
		fmt.Println("\n=== Testing Get Job ===")
		job, err := client.Jobs().Get(ctx, jobIDStr)
		if err != nil {
			log.Printf("Failed to get job %s: %v", jobIDStr, err)
		} else {
			fmt.Printf("Retrieved job %v:\n", job.JobID)
			fmt.Printf("  Name: %v\n", job.Name)
			fmt.Printf("  State: %v\n", job.JobState)
			fmt.Printf("  Partition: %v\n", job.Partition)
			fmt.Printf("  Working Directory: %v\n", job.CurrentWorkingDirectory)
		}

		// Test 4: Cancel the job
		fmt.Println("\n=== Testing Cancel Job ===")
		err = client.Jobs().Cancel(ctx, jobIDStr)
		if err != nil {
			log.Printf("Failed to cancel job %s: %v", jobIDStr, err)
		} else {
			fmt.Printf("Successfully cancelled job %s\n", jobIDStr)
		}
	}

	// Test 5: List Nodes
	fmt.Println("\n=== Testing List Nodes ===")
	nodes, err := client.Nodes().List(ctx, &types.ListNodesOptions{
		Limit: 5,
	})
	if err != nil {
		log.Printf("Failed to list nodes: %v", err)
	} else {
		fmt.Printf("Found %d nodes (showing up to 5)\n", nodes.Total)
		for i, node := range nodes.Nodes {
			nodeName := ""
			if node.Name != nil {
				nodeName = *node.Name
			}
			cpus := int32(0)
			if node.CPUs != nil {
				cpus = *node.CPUs
			}
			memory := int64(0)
			if node.RealMemory != nil {
				memory = *node.RealMemory
			}
			fmt.Printf("  [%d] Node %s: State=%v, CPUs=%d, Memory=%dMB\n",
				i+1, nodeName, node.State, cpus, memory)
		}
	}

	// Test 6: List Partitions
	fmt.Println("\n=== Testing List Partitions ===")
	partitions, err := client.Partitions().List(ctx, &types.ListPartitionsOptions{
		Limit: 5,
	})
	if err != nil {
		log.Printf("Failed to list partitions: %v", err)
	} else {
		fmt.Printf("Found %d partitions (showing up to 5)\n", partitions.Total)
		for i, partition := range partitions.Partitions {
			partitionName := ""
			if partition.Name != nil {
				partitionName = *partition.Name
			}
			state := "UNKNOWN"
			if partition.Partition != nil && len(partition.Partition.State) > 0 {
				state = string(partition.Partition.State[0])
			}
			totalNodes := int32(0)
			if partition.Nodes != nil && partition.Nodes.Total != nil {
				totalNodes = *partition.Nodes.Total
			}
			fmt.Printf("  [%d] Partition %s: State=%s, Nodes=%d total\n",
				i+1, partitionName, state, totalNodes)
		}
	}

	fmt.Println("\n=== All tests completed ===")
}
