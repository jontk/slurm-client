// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

// Example: Array job submission and management
func main() {
	// Create configuration
	cfg := config.NewDefault()
	cfg.BaseURL = "https://cluster.example.com:6820"

	// Create authentication
	authProvider := auth.NewTokenAuth("your-jwt-token")

	// Create client
	ctx := context.Background()
	client, err := slurm.NewClient(ctx,
		slurm.WithConfig(cfg),
		slurm.WithAuth(authProvider),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Example 1: Submit an array job
	fmt.Println("=== Array Job Submission ===")
	arrayJobID := submitArrayJob(ctx, client)

	// Example 2: Monitor array job progress
	fmt.Println("\n=== Array Job Monitoring ===")
	monitorArrayJob(ctx, client, arrayJobID)

	// Example 3: Get individual array task details
	fmt.Println("\n=== Array Task Details ===")
	getArrayTaskDetails(ctx, client, arrayJobID)

	// Example 4: Cancel specific array tasks
	fmt.Println("\n=== Selective Task Cancellation ===")
	cancelArrayTasks(ctx, client, arrayJobID)
}

// submitArrayJob submits an array job with multiple tasks
func submitArrayJob(ctx context.Context, client slurm.SlurmClient) string {
	// Create array job submission
	// Array syntax: jobname[1-100:5] means tasks 1-100 with step of 5
	job := &interfaces.JobSubmission{
		Name: "array-example",
		Script: `#!/bin/bash
#SBATCH --array=1-20:2  # Create tasks 1,3,5,7,9,11,13,15,17,19

echo "Starting array task $SLURM_ARRAY_TASK_ID"
echo "Processing dataset chunk $SLURM_ARRAY_TASK_ID"

# Simulate different processing based on task ID
sleep $((SLURM_ARRAY_TASK_ID * 2))

# Write output to task-specific file
echo "Task $SLURM_ARRAY_TASK_ID completed at $(date)" > output_$SLURM_ARRAY_TASK_ID.txt
`,
		Partition:  "compute",
		CPUs:       2,
		Memory:     4096, // 4GB
		TimeLimit:  30,   // 30 minutes
		WorkingDir: "/scratch/array-jobs",
		Environment: map[string]string{
			"DATA_DIR":    "/data/datasets",
			"OUTPUT_DIR":  "/scratch/array-jobs/outputs",
			"PYTHON_PATH": "/usr/bin/python3",
		},
	}

	// Submit the array job
	resp, err := client.Jobs().Submit(ctx, job)
	if err != nil {
		log.Fatalf("Failed to submit array job: %v", err)
	}

	fmt.Printf("Array job submitted with ID: %s\n", resp.JobID)
	fmt.Printf("Total array tasks: 10 (1,3,5,7,9,11,13,15,17,19)\n")

	return resp.JobID
}

// monitorArrayJob monitors the progress of an array job
func monitorArrayJob(ctx context.Context, client slurm.SlurmClient, arrayJobID string) {
	fmt.Printf("Monitoring array job %s...\n", arrayJobID)

	// Track task states
	taskStates := make(map[string]string)
	completedCount := 0

	// Monitor for up to 5 minutes
	timeout := time.After(5 * time.Minute)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// List all jobs with array job ID prefix
			jobs, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{
				UserID: "", // Filter by current user if needed
				Limit:  100,
			})
			if err != nil {
				log.Printf("Failed to list array tasks: %v", err)
				continue
			}

			// Update task states
			for _, job := range jobs.Jobs {
				// Extract array task ID from job name or ID
				// Format typically: jobid_taskid
				taskID := extractArrayTaskID(job.ID, arrayJobID)

				oldState, exists := taskStates[taskID]
				if !exists || oldState != job.State {
					taskStates[taskID] = job.State
					fmt.Printf("Task %s: %s", taskID, job.State)

					if exists && oldState != job.State {
						fmt.Printf(" (was %s)", oldState)
					}

					if job.State == "COMPLETED" {
						completedCount++
						fmt.Printf(" - Runtime: %s", formatRuntime(job.StartTime, job.EndTime))
					} else if job.State == "FAILED" {
						fmt.Printf(" - Exit code: %d", job.ExitCode)
					}

					fmt.Println()
				}
			}

			// Check if all tasks are done
			if completedCount == 10 { // We have 10 array tasks
				fmt.Println("All array tasks completed!")
				return
			}

		case <-timeout:
			fmt.Println("Monitoring timeout reached")
			return
		case <-ctx.Done():
			fmt.Println("Context cancelled")
			return
		}
	}
}

// getArrayTaskDetails gets details for specific array tasks
func getArrayTaskDetails(ctx context.Context, client slurm.SlurmClient, arrayJobID string) {
	// Get details for specific array tasks
	taskIDs := []string{"1", "5", "10", "15", "19"}

	fmt.Printf("Getting details for array tasks: %v\n", taskIDs)

	for _, taskID := range taskIDs {
		// Construct the full job ID for the array task
		// Format depends on SLURM configuration, typically: arrayJobID_taskID
		fullJobID := fmt.Sprintf("%s_%s", arrayJobID, taskID)

		job, err := client.Jobs().Get(ctx, fullJobID)
		if err != nil {
			// Try alternative format
			fullJobID = fmt.Sprintf("%s[%s]", arrayJobID, taskID)
			job, err = client.Jobs().Get(ctx, fullJobID)
			if err != nil {
				log.Printf("Failed to get task %s: %v", taskID, err)
				continue
			}
		}

		fmt.Printf("\nArray Task %s:\n", taskID)
		fmt.Printf("  State:      %s\n", job.State)
		fmt.Printf("  Nodes:      %v\n", job.Nodes)
		fmt.Printf("  CPUs:       %d\n", job.CPUs)
		fmt.Printf("  Memory:     %d MB\n", job.Memory)

		if job.StartTime != nil {
			fmt.Printf("  Start Time: %s\n", job.StartTime.Format(time.TimeOnly))
		}
		if job.EndTime != nil && job.StartTime != nil {
			runtime := job.EndTime.Sub(*job.StartTime)
			fmt.Printf("  Runtime:    %s\n", runtime.Round(time.Second))
		}

		// Show environment variables specific to array tasks
		if arrayTaskID, ok := job.Environment["SLURM_ARRAY_TASK_ID"]; ok {
			fmt.Printf("  Array Task ID: %s\n", arrayTaskID)
		}
		if arrayJobID, ok := job.Environment["SLURM_ARRAY_JOB_ID"]; ok {
			fmt.Printf("  Array Job ID:  %s\n", arrayJobID)
		}
	}
}

// cancelArrayTasks demonstrates selective cancellation of array tasks
func cancelArrayTasks(ctx context.Context, client slurm.SlurmClient, arrayJobID string) {
	// Cancel specific array tasks (e.g., only odd-numbered tasks)
	tasksToCancel := []string{"11", "13", "15", "17", "19"}

	fmt.Printf("Cancelling array tasks: %v\n", tasksToCancel)

	for _, taskID := range tasksToCancel {
		// Try different job ID formats
		jobIDs := []string{
			fmt.Sprintf("%s_%s", arrayJobID, taskID),
			fmt.Sprintf("%s[%s]", arrayJobID, taskID),
		}

		cancelled := false
		for _, jobID := range jobIDs {
			err := client.Jobs().Cancel(ctx, jobID)
			if err == nil {
				fmt.Printf("Cancelled task %s\n", taskID)
				cancelled = true
				break
			}
		}

		if !cancelled {
			log.Printf("Failed to cancel task %s", taskID)
		}
	}

	// Alternative: Cancel entire array job
	// This would cancel all pending/running tasks but leave completed ones
	fmt.Println("\nTo cancel the entire array job, use:")
	fmt.Printf("client.Jobs().Cancel(ctx, \"%s\")\n", arrayJobID)
}

// Helper functions

func extractArrayTaskID(jobID, arrayJobID string) string {
	// Extract array task ID from full job ID
	// This is SLURM configuration dependent
	// Common formats: arrayJobID_taskID or arrayJobID[taskID]

	prefix1 := arrayJobID + "_"
	if len(jobID) > len(prefix1) && jobID[:len(prefix1)] == prefix1 {
		return jobID[len(prefix1):]
	}

	prefix2 := arrayJobID + "["
	if len(jobID) > len(prefix2) && jobID[:len(prefix2)] == prefix2 {
		for i := len(prefix2); i < len(jobID); i++ {
			if jobID[i] == ']' {
				return jobID[len(prefix2):i]
			}
		}
	}

	return jobID
}

func formatRuntime(start, end *time.Time) string {
	if start == nil || end == nil {
		return "N/A"
	}
	runtime := end.Sub(*start)
	return runtime.Round(time.Second).String()
}
