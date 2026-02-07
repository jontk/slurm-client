// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Example demonstrating complete JobManager functionality
package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/jontk/slurm-client"
)

func main() {
	// This example shows how the complete JobManager implementation works
	// Note: This requires a real Slurm REST API server to work

	ctx := context.Background()

	// Create a v0.0.42 client (stable version)
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.42",
		slurm.WithBaseURL("https://your-slurm-server:6820"), // Replace with real URL
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Get the job manager
	jobManager := client.Jobs()

	// 1. Submit a new job
	fmt.Println("=== Job Submission ===")
	jobSubmission := &slurm.JobSubmission{
		Name:       "example-job",
		Script:     "#!/bin/bash\necho 'Hello from Slurm!'\nsleep 30",
		Partition:  "debug",
		CPUs:       2,
		Memory:     1024 * 1024 * 1024, // 1GB in bytes
		TimeLimit:  5,                  // 5 minutes
		Nodes:      1,
		WorkingDir: "/tmp",
		Environment: map[string]string{
			"MY_VAR": "example_value",
		},
	}

	submitResp, err := jobManager.Submit(ctx, jobSubmission)
	if err != nil {
		fmt.Printf("Expected error (no real server): %v\n", err)
		// Continue with example using a mock job ID
		fmt.Println("Using mock job ID for demonstration...")
		mockJobID := "12345"

		// 2. Get specific job details
		fmt.Println("\n=== Job Details ===")
		job, err := jobManager.Get(ctx, mockJobID)
		if err != nil {
			fmt.Printf("Expected error (no real server): %v\n", err)
		} else {
			fmt.Printf("Job %v: %v\n", job.JobID, job.Name)
			fmt.Printf("  State: %v\n", job.JobState)
			fmt.Printf("  User: %v\n", job.UserID)
			fmt.Printf("  CPUs: %d\n", job.CPUs)
			fmt.Printf("  Memory: %d MB\n", *job.MemoryPerNode/(1024*1024))
		}

		// 3. List all jobs
		fmt.Println("\n=== Job List ===")
		jobList, err := jobManager.List(ctx, nil)
		if err != nil {
			fmt.Printf("Expected error (no real server): %v\n", err)
		} else {
			fmt.Printf("Found %d jobs:\n", jobList.Total)
			for _, job := range jobList.Jobs {
				fmt.Printf("- Job %v: %v (%v) - %v\n", job.JobID, job.Name, job.UserID, job.JobState)
			}
		}

		// 4. List jobs with filtering
		fmt.Println("\n=== Filtered Job List ===")
		filteredJobList, err := jobManager.List(ctx, &slurm.ListJobsOptions{
			UserID: "1000",
			States: []string{"RUNNING", "PENDING"},
			Limit:  10,
		})
		if err != nil {
			fmt.Printf("Expected error (no real server): %v\n", err)
		} else {
			fmt.Printf("Found %d filtered jobs:\n", filteredJobList.Total)
			for _, job := range filteredJobList.Jobs {
				fmt.Printf("- Job %v: %v - %v\n", job.JobID, job.Name, job.JobState)
			}
		}

		// 5. Cancel the job
		fmt.Println("\n=== Job Cancellation ===")
		err = jobManager.Cancel(ctx, mockJobID)
		if err != nil {
			fmt.Printf("Expected error (no real server): %v\n", err)
		} else {
			fmt.Printf("Successfully cancelled job %s\n", mockJobID)
		}

		return
	}

	// If we get here, job submission was successful
	fmt.Printf("Successfully submitted job %s\n", strconv.Itoa(int(submitResp.JobId)))

	// 2. Get the submitted job details
	fmt.Println("\n=== Job Details ===")
	job, err := jobManager.Get(ctx, strconv.Itoa(int(submitResp.JobId)))
	if err != nil {
		fmt.Printf("Failed to get job details: %v\n", err)
	} else {
		fmt.Printf("Job %v: %v\n", job.JobID, job.Name)
		fmt.Printf("  State: %v\n", job.JobState)
		fmt.Printf("  User: %v\n", job.UserID)
		fmt.Printf("  CPUs: %d\n", job.CPUs)
		fmt.Printf("  Memory: %d MB\n", *job.MemoryPerNode/(1024*1024))
	}

	// 3. List all jobs
	fmt.Println("\n=== Job List ===")
	jobList, err := jobManager.List(ctx, nil)
	if err != nil {
		fmt.Printf("Failed to list jobs: %v\n", err)
	} else {
		fmt.Printf("Found %d jobs:\n", jobList.Total)
		for _, job := range jobList.Jobs {
			fmt.Printf("- Job %v: %v (%v) - %v\n", job.JobID, job.Name, job.UserID, job.JobState)
		}
	}

	// 4. Cancel the job (cleanup)
	fmt.Println("\n=== Job Cancellation ===")
	err = jobManager.Cancel(ctx, strconv.Itoa(int(submitResp.JobId)))
	if err != nil {
		fmt.Printf("Failed to cancel job: %v\n", err)
	} else {
		fmt.Printf("Successfully cancelled job %s\n", strconv.Itoa(int(submitResp.JobId)))
	}
}
