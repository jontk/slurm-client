package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/pkg/auth"
)

func main() {
	ctx := context.Background()

	// Create client with stable API version for production workflows
	client, err := slurm.NewClientWithVersion(ctx, slurm.StableVersion(),
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(auth.NewNoneAuth()),
		slurm.WithTimeout("30s"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Printf("Using stable API version: %s\n", slurm.StableVersion())

	// Example 1: Simple job submission
	fmt.Println("\n=== Example 1: Simple Job Submission ===")
	if err := submitSimpleJob(ctx, client); err != nil {
		log.Printf("Simple job submission failed: %v", err)
	}

	// Example 2: Batch job submission
	fmt.Println("\n=== Example 2: Batch Job Submission ===")
	if err := submitBatchJob(ctx, client); err != nil {
		log.Printf("Batch job submission failed: %v", err)
	}

	// Example 3: Job monitoring workflow
	fmt.Println("\n=== Example 3: Job Monitoring Workflow ===")
	if err := monitorJobWorkflow(ctx, client); err != nil {
		log.Printf("Job monitoring failed: %v", err)
	}

	// Example 4: Resource-aware job submission
	fmt.Println("\n=== Example 4: Resource-Aware Job Submission ===")
	if err := submitResourceAwareJob(ctx, client); err != nil {
		log.Printf("Resource-aware job submission failed: %v", err)
	}
}

func submitSimpleJob(ctx context.Context, client slurm.SlurmClient) error {
	// Define a simple job
	jobReq := &slurm.JobSubmissionRequest{
		Script: `#!/bin/bash
#SBATCH --job-name=hello-world
#SBATCH --output=hello-%j.out
#SBATCH --ntasks=1
#SBATCH --time=00:01:00

echo "Hello from Slurm job!"
hostname
date`,
		Name:      "hello-world-example",
		Partition: "debug",
		Nodes:     1,
		Tasks:     1,
		Time:      "00:01:00",
	}

	// Submit the job
	jobID, err := client.SubmitJob(ctx, jobReq)
	if err != nil {
		return fmt.Errorf("job submission failed: %w", err)
	}

	fmt.Printf("✓ Job submitted successfully! Job ID: %d\n", jobID)

	// Wait for job to complete (with timeout)
	return waitForJobCompletion(ctx, client, jobID, 2*time.Minute)
}

func submitBatchJob(ctx context.Context, client slurm.SlurmClient) error {
	// Define a batch job with multiple tasks
	jobReq := &slurm.JobSubmissionRequest{
		Script: `#!/bin/bash
#SBATCH --job-name=batch-processing
#SBATCH --output=batch-%j-%t.out
#SBATCH --array=1-5
#SBATCH --ntasks=1
#SBATCH --time=00:02:00

echo "Processing batch item $SLURM_ARRAY_TASK_ID"
sleep 10
echo "Batch item $SLURM_ARRAY_TASK_ID completed"`,
		Name:      "batch-processing",
		Partition: "compute",
		Nodes:     1,
		Tasks:     1,
		Time:      "00:02:00",
		Array:     "1-5",
	}

	jobID, err := client.SubmitJob(ctx, jobReq)
	if err != nil {
		return fmt.Errorf("batch job submission failed: %w", err)
	}

	fmt.Printf("✓ Batch job submitted! Job ID: %d (array 1-5)\n", jobID)

	// Monitor array job progress
	return monitorArrayJob(ctx, client, jobID)
}

func submitResourceAwareJob(ctx context.Context, client slurm.SlurmClient) error {
	// First, check available resources
	partitions, err := client.ListPartitions(ctx)
	if err != nil {
		return fmt.Errorf("failed to list partitions: %w", err)
	}

	// Find best partition for our job
	var selectedPartition string
	var maxNodes int
	for _, partition := range partitions {
		availableNodes := countAvailableNodes(partition.Nodes)
		if availableNodes > maxNodes {
			maxNodes = availableNodes
			selectedPartition = partition.Name
		}
	}

	if selectedPartition == "" {
		return fmt.Errorf("no suitable partition found")
	}

	fmt.Printf("Selected partition: %s (available nodes: %d)\n", selectedPartition, maxNodes)

	// Submit resource-optimized job
	jobReq := &slurm.JobSubmissionRequest{
		Script: `#!/bin/bash
#SBATCH --job-name=resource-optimized
#SBATCH --output=optimized-%j.out
#SBATCH --ntasks-per-node=4
#SBATCH --mem=8G
#SBATCH --time=00:05:00

echo "Running on optimized resources"
echo "Allocated CPUs: $SLURM_CPUS_ON_NODE"
echo "Allocated Memory: $SLURM_MEM_PER_NODE MB"
echo "Node: $SLURMD_NODENAME"

# Simulate resource-intensive work
stress --cpu 4 --timeout 60s`,
		Name:         "resource-optimized",
		Partition:    selectedPartition,
		Nodes:        min(2, maxNodes),
		TasksPerNode: 4,
		Memory:       "8G",
		Time:         "00:05:00",
	}

	jobID, err := client.SubmitJob(ctx, jobReq)
	if err != nil {
		return fmt.Errorf("resource-aware job submission failed: %w", err)
	}

	fmt.Printf("✓ Resource-optimized job submitted! Job ID: %d\n", jobID)
	return nil
}

func monitorJobWorkflow(ctx context.Context, client slurm.SlurmClient) error {
	// Submit a monitoring test job
	jobReq := &slurm.JobSubmissionRequest{
		Script: `#!/bin/bash
#SBATCH --job-name=monitor-test
#SBATCH --output=monitor-%j.out
#SBATCH --ntasks=1
#SBATCH --time=00:03:00

for i in {1..30}; do
    echo "Progress: $i/30"
    sleep 5
done
echo "Job completed successfully"`,
		Name:      "monitor-test",
		Partition: "debug",
		Nodes:     1,
		Tasks:     1,
		Time:      "00:03:00",
	}

	jobID, err := client.SubmitJob(ctx, jobReq)
	if err != nil {
		return fmt.Errorf("monitoring test job submission failed: %w", err)
	}

	fmt.Printf("✓ Monitoring test job submitted! Job ID: %d\n", jobID)

	// Monitor job with detailed status updates
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	timeout := time.After(5 * time.Minute)

	for {
		select {
		case <-ticker.C:
			job, err := client.GetJob(ctx, jobID)
			if err != nil {
				log.Printf("Error getting job status: %v", err)
				continue
			}

			fmt.Printf("Job %d status: %s (runtime: %s)\n", 
				jobID, job.State, formatDuration(job.RunTime))

			if isJobFinished(job.State) {
				fmt.Printf("✓ Job finished with state: %s\n", job.State)
				return nil
			}

		case <-timeout:
			return fmt.Errorf("job monitoring timed out after 5 minutes")

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func monitorArrayJob(ctx context.Context, client slurm.SlurmClient, jobID uint32) error {
	fmt.Printf("Monitoring array job %d...\n", jobID)

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	timeout := time.After(5 * time.Minute)

	for {
		select {
		case <-ticker.C:
			jobs, err := client.ListJobs(ctx)
			if err != nil {
				log.Printf("Error listing jobs: %v", err)
				continue
			}

			// Find our array job tasks
			var arrayTasks []slurm.Job
			for _, job := range jobs {
				if job.ArrayJobID != nil && *job.ArrayJobID == jobID {
					arrayTasks = append(arrayTasks, job)
				}
			}

			if len(arrayTasks) == 0 {
				fmt.Printf("Array job %d completed or not found\n", jobID)
				return nil
			}

			// Count task states
			stateCount := make(map[string]int)
			for _, task := range arrayTasks {
				stateCount[task.State]++
			}

			fmt.Printf("Array job %d progress: ", jobID)
			for state, count := range stateCount {
				fmt.Printf("%s:%d ", state, count)
			}
			fmt.Println()

			// Check if all tasks are finished
			allFinished := true
			for _, task := range arrayTasks {
				if !isJobFinished(task.State) {
					allFinished = false
					break
				}
			}

			if allFinished {
				fmt.Printf("✓ All array job tasks completed\n")
				return nil
			}

		case <-timeout:
			return fmt.Errorf("array job monitoring timed out")

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func waitForJobCompletion(ctx context.Context, client slurm.SlurmClient, jobID uint32, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			job, err := client.GetJob(ctx, jobID)
			if err != nil {
				return fmt.Errorf("error checking job status: %w", err)
			}

			fmt.Printf("Job %d status: %s\n", jobID, job.State)

			if isJobFinished(job.State) {
				if job.State == "COMPLETED" {
					fmt.Printf("✓ Job completed successfully\n")
				} else {
					fmt.Printf("✗ Job finished with state: %s\n", job.State)
				}
				return nil
			}

		case <-ctx.Done():
			return fmt.Errorf("job wait timed out or cancelled")
		}
	}
}

func countAvailableNodes(nodes []slurm.Node) int {
	count := 0
	for _, node := range nodes {
		if node.State == "IDLE" || node.State == "MIXED" {
			count++
		}
	}
	return count
}

func isJobFinished(state string) bool {
	finishedStates := []string{
		"COMPLETED", "FAILED", "CANCELLED", "TIMEOUT", "NODE_FAIL", "PREEMPTED",
	}
	for _, finished := range finishedStates {
		if state == finished {
			return true
		}
	}
	return false
}

func formatDuration(seconds int) string {
	duration := time.Duration(seconds) * time.Second
	return duration.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}