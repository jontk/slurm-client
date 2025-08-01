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
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/types"
)

func main() {
	ctx := context.Background()

	// Create client
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(auth.NewTokenAuth("your-jwt-token")),
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

func listAllJobs(ctx context.Context, client *slurm.Client) {
	jobs, err := client.Jobs().List(ctx, nil)
	if err != nil {
		log.Printf("Failed to list jobs: %v", err)
		return
	}

	fmt.Printf("Total jobs: %d\n", len(jobs))
	
	// Display first 5 jobs
	for i, job := range jobs {
		if i >= 5 {
			break
		}
		fmt.Printf("Job %d: %s (User: %s, State: %s)\n",
			job.JobID, job.Name, job.UserName, job.State)
	}
}

func filterByState(ctx context.Context, client *slurm.Client) {
	// Filter for running and pending jobs
	filters := &types.JobFilters{
		States: []string{"RUNNING", "PENDING"},
	}

	jobs, err := client.Jobs().List(ctx, filters)
	if err != nil {
		log.Printf("Failed to filter jobs: %v", err)
		return
	}

	fmt.Printf("Running/Pending jobs: %d\n", len(jobs))
	for _, job := range jobs {
		fmt.Printf("- Job %d: %s (State: %s, Nodes: %d)\n",
			job.JobID, job.Name, job.State, job.NodeCount)
	}
}

func filterByUser(ctx context.Context, client *slurm.Client) {
	// Filter jobs for specific users
	filters := &types.JobFilters{
		Users: []string{"alice", "bob"},
	}

	jobs, err := client.Jobs().List(ctx, filters)
	if err != nil {
		log.Printf("Failed to filter by user: %v", err)
		return
	}

	fmt.Printf("Jobs for alice and bob: %d\n", len(jobs))
	for _, job := range jobs {
		fmt.Printf("- Job %d: %s (User: %s)\n",
			job.JobID, job.Name, job.UserName)
	}
}

func advancedFiltering(ctx context.Context, client *slurm.Client) {
	// Complex filtering with multiple criteria
	filters := &types.JobFilters{
		States:     []string{"RUNNING"},
		Partitions: []string{"compute", "gpu"},
		StartTime:  time.Now().Add(-24 * time.Hour), // Jobs started in last 24h
		EndTime:    time.Now(),
	}

	jobs, err := client.Jobs().List(ctx, filters)
	if err != nil {
		log.Printf("Failed to apply advanced filters: %v", err)
		return
	}

	fmt.Printf("Jobs matching advanced criteria: %d\n", len(jobs))
	for _, job := range jobs {
		runtime := time.Since(job.StartTime)
		fmt.Printf("- Job %d: %s (Partition: %s, Runtime: %v)\n",
			job.JobID, job.Name, job.Partition, runtime.Round(time.Minute))
	}
}

func paginatedList(ctx context.Context, client *slurm.Client) {
	// Paginate through jobs
	pageSize := 10
	offset := 0
	totalJobs := 0

	for {
		opts := &types.ListOptions{
			Limit:  pageSize,
			Offset: offset,
		}

		jobs, err := client.Jobs().ListWithOptions(ctx, nil, opts)
		if err != nil {
			log.Printf("Failed to get page: %v", err)
			break
		}

		if len(jobs) == 0 {
			break
		}

		fmt.Printf("Page %d: Retrieved %d jobs\n", (offset/pageSize)+1, len(jobs))
		totalJobs += len(jobs)

		// Process jobs in this page
		for _, job := range jobs {
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