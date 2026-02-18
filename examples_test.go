// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package slurm_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/pkg/errors"
)

// Example_createClient demonstrates how to create a basic SLURM client
// with automatic version detection.
func Example_createClient() {
	ctx := context.Background()

	// Create a client with automatic version detection
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://your-slurm-server:6820"),
		slurm.WithUserToken("your-username", "your-token"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// The client automatically detects and uses the best compatible API version
	fmt.Printf("Connected to SLURM API version: %s\n", client.Version())
}

// Example_createClientWithVersion demonstrates how to create a client
// with a specific API version.
func Example_createClientWithVersion() {
	ctx := context.Background()

	// Create a client with a specific API version
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.44",
		slurm.WithBaseURL("https://your-slurm-server:6820"),
		slurm.WithUserToken("your-username", "your-token"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Printf("Using SLURM API version: %s\n", client.Version())
}

// Example_listJobs demonstrates how to list jobs from the SLURM cluster.
func Example_listJobs() {
	ctx := context.Background()

	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://your-slurm-server:6820"),
		slurm.WithUserToken("your-username", "your-token"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// List all jobs
	jobs, err := client.Jobs().List(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d jobs\n", len(jobs.Jobs))
	for _, job := range jobs.Jobs {
		fmt.Printf("Job ID: %v, State: %v, User: %v\n",
			job.JobID, job.JobState, job.UserID)
	}
}

// Example_listJobsWithFilters demonstrates how to list jobs with filtering options.
func Example_listJobsWithFilters() {
	ctx := context.Background()

	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://your-slurm-server:6820"),
		slurm.WithUserToken("your-username", "your-token"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// List jobs for a specific user
	jobs, err := client.Jobs().List(ctx, &slurm.ListJobsOptions{
		UserID: "alice",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d jobs for user alice\n", len(jobs.Jobs))
}

// Example_submitJob demonstrates how to submit a new job to the SLURM cluster.
func Example_submitJob() {
	ctx := context.Background()

	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://your-slurm-server:6820"),
		slurm.WithUserToken("your-username", "your-token"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Prepare job submission
	job := &slurm.JobSubmission{
		Name:   "test-job",
		Script: "#!/bin/bash\n#SBATCH --job-name=test\n#SBATCH --output=test.out\nsleep 60",
		CPUs:   2,
		Memory: 4096,
	}

	// Submit the job
	response, err := client.Jobs().Submit(ctx, job)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Job submitted successfully. Job ID: %d\n", response.JobId)
}

// Example_getNode demonstrates how to get information about a specific node.
func Example_getNode() {
	ctx := context.Background()

	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://your-slurm-server:6820"),
		slurm.WithUserToken("your-username", "your-token"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Get information about a specific node
	node, err := client.Nodes().Get(ctx, "node001")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Node: %s\n", *node.Name)
	fmt.Printf("State: %s\n", node.State)
	fmt.Printf("CPUs: %d\n", node.CPUs)
	fmt.Printf("Memory: %d MB\n", node.RealMemory)
}

// Example_listNodes demonstrates how to list all nodes in the cluster.
func Example_listNodes() {
	ctx := context.Background()

	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://your-slurm-server:6820"),
		slurm.WithUserToken("your-username", "your-token"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// List all nodes
	nodes, err := client.Nodes().List(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d nodes\n", len(nodes.Nodes))
	for _, node := range nodes.Nodes {
		fmt.Printf("Node: %s, State: %s, CPUs: %d\n",
			*node.Name, node.State, node.CPUs)
	}
}

// Example_errorHandling demonstrates structured error handling with version context.
func Example_errorHandling() {
	ctx := context.Background()

	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://your-slurm-server:6820"),
		slurm.WithUserToken("your-username", "your-token"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Try to use a feature that might not be available in all versions
	reservations, err := client.Reservations().List(ctx, nil)
	if err != nil {
		// Check for authentication errors
		if errors.IsAuthenticationError(err) {
			fmt.Println("Authentication failed - please check your credentials")
			return
		}
		// Check for not implemented errors (feature not available in this version)
		if errors.IsNotImplementedError(err) {
			fmt.Printf("Reservations not available in API version %s\n", client.Version())
			return
		}
		log.Fatal(err)
	}

	fmt.Printf("Found %d reservations\n", len(reservations.Reservations))
}

// Example_withTimeout demonstrates how to set custom timeouts for operations.
func Example_withTimeout() {
	ctx := context.Background()

	// Create a client with a custom timeout
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://your-slurm-server:6820"),
		slurm.WithUserToken("your-username", "your-token"),
		slurm.WithTimeout(10*time.Second), // 10 second timeout for all operations
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// All operations will use the configured timeout
	jobs, err := client.Jobs().List(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d jobs\n", len(jobs.Jobs))
}
