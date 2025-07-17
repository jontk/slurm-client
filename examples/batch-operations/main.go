package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

// Example: Batch job submission and monitoring
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

	// Example 1: Submit multiple jobs in batch
	fmt.Println("=== Batch Job Submission ===")
	jobIDs := submitBatchJobs(ctx, client, 5)

	// Example 2: Monitor job progress
	fmt.Println("\n=== Job Progress Monitoring ===")
	monitorJobProgress(ctx, client, jobIDs)

	// Example 3: Collect job results
	fmt.Println("\n=== Job Results Collection ===")
	collectJobResults(ctx, client, jobIDs)

	// Example 4: Cleanup completed jobs
	fmt.Println("\n=== Cleanup ===")
	cleanupJobs(ctx, client, jobIDs)
}

// submitBatchJobs submits multiple jobs concurrently
func submitBatchJobs(ctx context.Context, client slurm.SlurmClient, count int) []string {
	var wg sync.WaitGroup
	jobIDs := make([]string, 0, count)
	jobChan := make(chan string, count)
	errChan := make(chan error, count)

	// Submit jobs concurrently
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			job := &interfaces.JobSubmission{
				Name:             fmt.Sprintf("batch-job-%d", index),
				Command:          fmt.Sprintf("python process.py --input data_%d.txt", index),
				Partition:  "compute",
				CPUs:       4,
				Memory:     8192, // 8GB
				TimeLimit:  60,   // 60 minutes
				WorkingDir: "/scratch/batch",
				Environment: map[string]string{
					"JOB_INDEX": fmt.Sprintf("%d", index),
					"BATCH_ID":  "example-batch",
				},
			}

			resp, err := client.Jobs().Submit(ctx, job)
			if err != nil {
				errChan <- fmt.Errorf("failed to submit job %d: %w", index, err)
				return
			}

			jobChan <- resp.JobID
			fmt.Printf("Submitted job %d: ID=%s\n", index, resp.JobID)
		}(i)
	}

	// Wait for all submissions to complete
	wg.Wait()
	close(jobChan)
	close(errChan)

	// Collect job IDs
	for jobID := range jobChan {
		jobIDs = append(jobIDs, jobID)
	}

	// Report any errors
	for err := range errChan {
		log.Printf("Submission error: %v", err)
	}

	return jobIDs
}

// monitorJobProgress monitors the progress of submitted jobs
func monitorJobProgress(ctx context.Context, client slurm.SlurmClient, jobIDs []string) {
	// Create a ticker for periodic checks
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	completedJobs := make(map[string]bool)
	startTime := time.Now()
	timeout := 30 * time.Minute

	fmt.Println("Monitoring job progress...")

	for {
		select {
		case <-ticker.C:
			// Check each job status
			for _, jobID := range jobIDs {
				if completedJobs[jobID] {
					continue
				}

				job, err := client.Jobs().Get(ctx, jobID)
				if err != nil {
					log.Printf("Failed to get job %s: %v", jobID, err)
					continue
				}

				fmt.Printf("Job %s: State=%s", jobID, job.State)
				if job.StartTime != nil {
					fmt.Printf(", Runtime=%s", time.Since(*job.StartTime).Round(time.Second))
				}
				fmt.Println()

				// Check if job completed
				if job.State == "COMPLETED" || job.State == "FAILED" || job.State == "CANCELLED" {
					completedJobs[jobID] = true
					if job.State == "FAILED" {
						log.Printf("Job %s failed!", jobID)
					}
				}
			}

			// Check if all jobs completed
			if len(completedJobs) == len(jobIDs) {
				fmt.Println("All jobs completed!")
				return
			}

		case <-time.After(timeout - time.Since(startTime)):
			fmt.Println("Timeout reached, stopping monitoring")
			return

		case <-ctx.Done():
			fmt.Println("Context cancelled, stopping monitoring")
			return
		}
	}
}

// collectJobResults collects results from completed jobs
func collectJobResults(ctx context.Context, client slurm.SlurmClient, jobIDs []string) {
	type JobResult struct {
		JobID     string
		State     string
		ExitCode  int
		StartTime *time.Time
		EndTime   *time.Time
		Runtime   time.Duration
		CPUTime   float64
		MaxMemory int64
	}

	results := make([]JobResult, 0, len(jobIDs))

	for _, jobID := range jobIDs {
		job, err := client.Jobs().Get(ctx, jobID)
		if err != nil {
			log.Printf("Failed to get job %s: %v", jobID, err)
			continue
		}

		result := JobResult{
			JobID:    job.ID,
			State:    job.State,
			ExitCode: job.ExitCode,
		}

		// Calculate runtime if available
		if job.StartTime != nil && job.EndTime != nil {
			result.StartTime = job.StartTime
			result.EndTime = job.EndTime
			result.Runtime = job.EndTime.Sub(*job.StartTime)
		}

		// Get additional metrics from metadata if available
		if cpuTime, ok := job.Metadata["cpu_time"].(float64); ok {
			result.CPUTime = cpuTime
		}
		if maxMem, ok := job.Metadata["max_memory"].(int64); ok {
			result.MaxMemory = maxMem
		}

		results = append(results, result)
	}

	// Generate summary report
	fmt.Println("\n=== Job Results Summary ===")
	fmt.Printf("%-15s %-12s %-10s %-15s %-12s\n", 
		"Job ID", "State", "Exit Code", "Runtime", "CPU Time")
	fmt.Println(strings.Repeat("-", 70))

	var totalRuntime time.Duration
	var totalCPUTime float64
	successCount := 0

	for _, result := range results {
		runtimeStr := "N/A"
		if result.Runtime > 0 {
			runtimeStr = result.Runtime.Round(time.Second).String()
			totalRuntime += result.Runtime
		}

		cpuTimeStr := "N/A"
		if result.CPUTime > 0 {
			cpuTimeStr = fmt.Sprintf("%.2fs", result.CPUTime)
			totalCPUTime += result.CPUTime
		}

		fmt.Printf("%-15s %-12s %-10d %-15s %-12s\n",
			result.JobID, result.State, result.ExitCode, runtimeStr, cpuTimeStr)

		if result.State == "COMPLETED" && result.ExitCode == 0 {
			successCount++
		}
	}

	fmt.Println(strings.Repeat("-", 70))
	fmt.Printf("Success Rate: %d/%d (%.1f%%)\n", 
		successCount, len(results), float64(successCount)/float64(len(results))*100)
	if totalRuntime > 0 {
		fmt.Printf("Total Runtime: %s\n", totalRuntime.Round(time.Second))
		fmt.Printf("Average Runtime: %s\n", (totalRuntime / time.Duration(len(results))).Round(time.Second))
	}
	if totalCPUTime > 0 {
		fmt.Printf("Total CPU Time: %.2fs\n", totalCPUTime)
	}
}

// cleanupJobs cancels any still-running jobs and performs cleanup
func cleanupJobs(ctx context.Context, client slurm.SlurmClient, jobIDs []string) {
	fmt.Println("Cleaning up jobs...")
	
	for _, jobID := range jobIDs {
		job, err := client.Jobs().Get(ctx, jobID)
		if err != nil {
			log.Printf("Failed to get job %s: %v", jobID, err)
			continue
		}

		// Cancel if still running
		if job.State == "RUNNING" || job.State == "PENDING" {
			fmt.Printf("Cancelling job %s (state: %s)\n", jobID, job.State)
			if err := client.Jobs().Cancel(ctx, jobID); err != nil {
				log.Printf("Failed to cancel job %s: %v", jobID, err)
			}
		}
	}

	fmt.Println("Cleanup completed")
}