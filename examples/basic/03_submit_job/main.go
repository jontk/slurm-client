// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// This example demonstrates how to submit jobs to SLURM using the client.
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

	// Example 1: Submit a simple job
	fmt.Println("=== Example 1: Simple Job ===")
	simpleJob(ctx, client)

	// Example 2: Submit a job with resources
	fmt.Println("\n=== Example 2: Job with Resources ===")
	resourceJob(ctx, client)

	// Example 3: Submit an array job
	fmt.Println("\n=== Example 3: Array Job ===")
	arrayJob(ctx, client)

	// Example 4: Submit a job with dependencies
	fmt.Println("\n=== Example 4: Job with Dependencies ===")
	dependentJob(ctx, client)
}

func simpleJob(ctx context.Context, client *slurm.Client) {
	// Create a simple job specification
	jobSpec := &types.JobSubmitRequest{
		Script: "#!/bin/bash\n#SBATCH --job-name=simple-test\n\necho 'Hello from SLURM!'\nsleep 60\necho 'Job completed'",
		Job: &types.JobProperties{
			Partition: "compute",
			Name:      "simple-test-job",
			TimeLimit: "00:05:00", // 5 minutes
			Account:   "default",
		},
	}

	// Submit the job
	jobID, err := client.Jobs().Submit(ctx, jobSpec)
	if err != nil {
		log.Printf("Failed to submit job: %v", err)
		return
	}

	fmt.Printf("Successfully submitted job with ID: %d\n", jobID)

	// Get job details
	job, err := client.Jobs().Get(ctx, jobID)
	if err != nil {
		log.Printf("Failed to get job details: %v", err)
		return
	}

	fmt.Printf("Job State: %s\n", job.State)
	fmt.Printf("Job Partition: %s\n", job.Partition)
}

func resourceJob(ctx context.Context, client *slurm.Client) {
	// Job with specific resource requirements
	jobSpec := &types.JobSubmitRequest{
		Script: `#!/bin/bash
#SBATCH --job-name=resource-test
#SBATCH --nodes=2
#SBATCH --ntasks-per-node=4
#SBATCH --cpus-per-task=2
#SBATCH --mem=16G
#SBATCH --time=01:00:00

echo "Running on $(hostname)"
echo "Total tasks: $SLURM_NTASKS"
echo "CPUs per task: $SLURM_CPUS_PER_TASK"

# Run a parallel application
srun hostname`,
		Job: &types.JobProperties{
			Partition: "compute",
			Name:      "resource-test-job",
			Nodes:     2,
			Tasks:     8,
			CPUsPerTask: 2,
			Memory:    "16G",
			TimeLimit: "01:00:00",
		},
	}

	jobID, err := client.Jobs().Submit(ctx, jobSpec)
	if err != nil {
		log.Printf("Failed to submit resource job: %v", err)
		return
	}

	fmt.Printf("Submitted resource job with ID: %d\n", jobID)
	fmt.Printf("Requested: 2 nodes, 8 tasks, 2 CPUs/task, 16GB memory\n")
}

func arrayJob(ctx context.Context, client *slurm.Client) {
	// Submit an array job
	jobSpec := &types.JobSubmitRequest{
		Script: `#!/bin/bash
#SBATCH --job-name=array-test
#SBATCH --array=1-10
#SBATCH --time=00:10:00

echo "This is array task $SLURM_ARRAY_TASK_ID"
sleep 30
echo "Task $SLURM_ARRAY_TASK_ID completed"`,
		Job: &types.JobProperties{
			Partition: "compute",
			Name:      "array-test-job",
			Array:     "1-10",
			TimeLimit: "00:10:00",
		},
	}

	jobID, err := client.Jobs().Submit(ctx, jobSpec)
	if err != nil {
		log.Printf("Failed to submit array job: %v", err)
		return
	}

	fmt.Printf("Submitted array job with ID: %d\n", jobID)
	fmt.Println("Array tasks 1-10 will run in parallel based on available resources")
}

func dependentJob(ctx context.Context, client *slurm.Client) {
	// First, submit a job that others will depend on
	prereqSpec := &types.JobSubmitRequest{
		Script: "#!/bin/bash\necho 'Prerequisite job running'\nsleep 30\necho 'Prerequisite complete'",
		Job: &types.JobProperties{
			Partition: "compute",
			Name:      "prerequisite-job",
			TimeLimit: "00:05:00",
		},
	}

	prereqID, err := client.Jobs().Submit(ctx, prereqSpec)
	if err != nil {
		log.Printf("Failed to submit prerequisite job: %v", err)
		return
	}

	fmt.Printf("Submitted prerequisite job with ID: %d\n", prereqID)

	// Submit a job that depends on the prerequisite
	dependentSpec := &types.JobSubmitRequest{
		Script: fmt.Sprintf(`#!/bin/bash
#SBATCH --dependency=afterok:%d

echo 'Dependent job started after job %d completed successfully'
date
echo 'Running dependent task...'
sleep 20
echo 'Dependent job complete'`, prereqID, prereqID),
		Job: &types.JobProperties{
			Partition:  "compute",
			Name:       "dependent-job",
			TimeLimit:  "00:05:00",
			Dependency: fmt.Sprintf("afterok:%d", prereqID),
		},
	}

	dependentID, err := client.Jobs().Submit(ctx, dependentSpec)
	if err != nil {
		log.Printf("Failed to submit dependent job: %v", err)
		return
	}

	fmt.Printf("Submitted dependent job with ID: %d\n", dependentID)
	fmt.Printf("Job %d will start after job %d completes successfully\n", dependentID, prereqID)

	// Monitor the jobs
	fmt.Println("\nMonitoring job progress...")
	for i := 0; i < 10; i++ {
		prereqJob, _ := client.Jobs().Get(ctx, prereqID)
		dependentJob, _ := client.Jobs().Get(ctx, dependentID)

		fmt.Printf("  Prerequisite job %d: %s\n", prereqID, prereqJob.State)
		fmt.Printf("  Dependent job %d: %s\n", dependentID, dependentJob.State)

		if dependentJob.State == "COMPLETED" || dependentJob.State == "FAILED" {
			break
		}

		time.Sleep(5 * time.Second)
	}
}