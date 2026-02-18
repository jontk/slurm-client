// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// This example demonstrates how to list and filter jobs using the SLURM client.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jontk/slurm-client"
)

func main() {
	ctx := context.Background()

	// Create client with proper authentication
	// IMPORTANT: Use WithUserToken to set both username and token headers
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithUserToken("your-username", "your-jwt-token"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Example 1: List all jobs
	fmt.Println("=== Example 1: List All Jobs ===")
	listAllJobs(ctx, client)

	// Example 2: Filter jobs by state
	fmt.Println("\n=== Example 2: Filter by State ===")
	filterByState(ctx, client)

	// Example 3: Filter jobs by user
	fmt.Println("\n=== Example 3: Filter by User ===")
	filterByUser(ctx, client)

	// Example 4: Advanced filtering
	fmt.Println("\n=== Example 4: Advanced Filtering ===")
	advancedFiltering(ctx, client)

	// Example 5: Pagination
	fmt.Println("\n=== Example 5: Pagination ===")
	paginatedList(ctx, client)
}

func listAllJobs(ctx context.Context, client slurm.SlurmClient) {
	jobList, err := client.Jobs().List(ctx, nil)
	if err != nil {
		log.Printf("Failed to list jobs: %v", err)
		return
	}

	fmt.Printf("Total jobs: %d\n", jobList.Total)

	// Display first 5 jobs
	for i, job := range jobList.Jobs {
		if i >= 5 {
			break
		}
		fmt.Printf("Job %v: %v (User: %v, State: %v)\n",
			job.JobID, job.Name, job.UserID, job.JobState)
	}
}

func filterByState(ctx context.Context, client slurm.SlurmClient) {
	// Filter for running and pending jobs
	opts := &slurm.ListJobsOptions{
		States: []string{"RUNNING", "PENDING"},
	}

	jobList, err := client.Jobs().List(ctx, opts)
	if err != nil {
		log.Printf("Failed to filter jobs: %v", err)
		return
	}

	fmt.Printf("Running/Pending jobs: %d\n", jobList.Total)
	for _, job := range jobList.Jobs {
		nodesCount := 0
		if job.Nodes != nil {
			nodesCount = len(*job.Nodes)
		}
		fmt.Printf("- Job %v: %v (State: %v, Nodes: %d)\n",
			job.JobID, job.Name, job.JobState, nodesCount)
	}
}

func filterByUser(ctx context.Context, client slurm.SlurmClient) {
	// Filter jobs for specific user
	// Note: API only supports filtering by one user at a time
	opts := &slurm.ListJobsOptions{
		UserID: "alice",
	}

	jobList, err := client.Jobs().List(ctx, opts)
	if err != nil {
		log.Printf("Failed to filter by user: %v", err)
		return
	}

	fmt.Printf("Jobs for alice: %d\n", jobList.Total)
	for _, job := range jobList.Jobs {
		fmt.Printf("- Job %v: %v (User: %v)\n",
			job.JobID, job.Name, job.UserID)
	}
}

func advancedFiltering(ctx context.Context, client slurm.SlurmClient) {
	// Complex filtering with multiple criteria
	// Note: Can only filter by one partition at a time
	opts := &slurm.ListJobsOptions{
		States:    []string{"RUNNING"},
		Partition: "compute",
	}

	jobList, err := client.Jobs().List(ctx, opts)
	if err != nil {
		log.Printf("Failed to apply advanced filters: %v", err)
		return
	}

	fmt.Printf("Jobs matching advanced criteria: %d\n", jobList.Total)
	for _, job := range jobList.Jobs {
		if !job.StartTime.IsZero() {
			runtime := time.Since(job.StartTime)
			fmt.Printf("- Job %v: %v (Partition: %v, Runtime: %v)\n",
				job.JobID, job.Name, job.Partition, runtime.Round(time.Minute))
		} else {
			fmt.Printf("- Job %v: %v (Partition: %v, Not started)\n",
				job.JobID, job.Name, job.Partition)
		}
	}
}

func paginatedList(ctx context.Context, client slurm.SlurmClient) {
	// Paginate through jobs
	pageSize := 10
	offset := 0
	totalJobs := 0

	for {
		opts := &slurm.ListJobsOptions{
			Limit:  pageSize,
			Offset: offset,
		}

		jobList, err := client.Jobs().List(ctx, opts)
		if err != nil {
			log.Printf("Failed to get page: %v", err)
			break
		}

		if len(jobList.Jobs) == 0 {
			break
		}

		fmt.Printf("Page %d: Retrieved %d jobs\n", (offset/pageSize)+1, len(jobList.Jobs))
		totalJobs += len(jobList.Jobs)

		// Process jobs in this page
		for _, job := range jobList.Jobs {
			// Process each job...
			_ = job
		}

		// Move to next page
		offset += pageSize

		// Stop after 3 pages for this example
		if offset >= pageSize*3 {
			fmt.Println("(Stopping after 3 pages for demo)")
			break
		}
	}

	fmt.Printf("Total jobs processed: %d\n", totalJobs)
}
