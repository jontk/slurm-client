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
)

func main() {
	ctx := context.Background()

	// Create client
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithUserToken("your-username", "your-jwt-token"),
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

func simpleJob(ctx context.Context, client slurm.SlurmClient) {
	// Create a simple job specification
	jobSpec := &slurm.JobSubmission{
		Name:      "simple-test-job",
		Script:    "#!/bin/bash\n#SBATCH --job-name=simple-test\n\necho 'Hello from SLURM!'\nsleep 60\necho 'Job completed'",
		Partition: "compute",
		TimeLimit: 300, // 5 minutes in seconds
		Account:   "default",
	}

	// Submit the job
	response, err := client.Jobs().Submit(ctx, jobSpec)
	if err != nil {
		log.Printf("Failed to submit job: %v", err)
		return
	}

	fmt.Printf("Successfully submitted job with ID: %d\n", response.JobId)

	// Get job details
	jobIDStr := fmt.Sprintf("%d", response.JobId)
	job, err := client.Jobs().Get(ctx, jobIDStr)
	if err != nil {
		log.Printf("Failed to get job details: %v", err)
		return
	}

	if len(job.JobState) > 0 {
		fmt.Printf("Job State: %s\n", string(job.JobState[0]))
	}
	fmt.Printf("Job Partition: %v\n", job.Partition)
}

func resourceJob(ctx context.Context, client slurm.SlurmClient) {
	// Job with specific resource requirements
	jobSpec := &slurm.JobSubmission{
		Name: "resource-test-job",
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
		Partition: "compute",
		Nodes:     2,
		CPUs:      16,    // 2 nodes * 4 tasks/node * 2 cpus/task
		Memory:    16384, // 16GB in MB
		TimeLimit: 3600,  // 1 hour in seconds
	}

	response, err := client.Jobs().Submit(ctx, jobSpec)
	if err != nil {
		log.Printf("Failed to submit resource job: %v", err)
		return
	}

	fmt.Printf("Submitted resource job with ID: %d\n", response.JobId)
	fmt.Printf("Requested: 2 nodes, 16 CPUs, 16GB memory\n")
}

func arrayJob(ctx context.Context, client slurm.SlurmClient) {
	// Submit an array job
	// Note: Array job syntax is typically handled in the script itself
	jobSpec := &slurm.JobSubmission{
		Name: "array-test-job",
		Script: `#!/bin/bash
#SBATCH --job-name=array-test
#SBATCH --array=1-10
#SBATCH --time=00:10:00

echo "This is array task $SLURM_ARRAY_TASK_ID"
sleep 30
echo "Task $SLURM_ARRAY_TASK_ID completed"`,
		Partition: "compute",
		TimeLimit: 600, // 10 minutes in seconds
	}

	response, err := client.Jobs().Submit(ctx, jobSpec)
	if err != nil {
		log.Printf("Failed to submit array job: %v", err)
		return
	}

	fmt.Printf("Submitted array job with ID: %d\n", response.JobId)
	fmt.Println("Array tasks 1-10 will run in parallel based on available resources")
}

func dependentJob(ctx context.Context, client slurm.SlurmClient) {
	// First, submit a job that others will depend on
	prereqSpec := &slurm.JobSubmission{
		Name:      "prerequisite-job",
		Script:    "#!/bin/bash\necho 'Prerequisite job running'\nsleep 30\necho 'Prerequisite complete'",
		Partition: "compute",
		TimeLimit: 300, // 5 minutes in seconds
	}

	prereqResp, err := client.Jobs().Submit(ctx, prereqSpec)
	if err != nil {
		log.Printf("Failed to submit prerequisite job: %v", err)
		return
	}

	fmt.Printf("Submitted prerequisite job with ID: %d\n", prereqResp.JobId)

	// Submit a job that depends on the prerequisite
	dependentSpec := &slurm.JobSubmission{
		Name: "dependent-job",
		Script: fmt.Sprintf(`#!/bin/bash
#SBATCH --dependency=afterok:%d

echo 'Dependent job started after job %d completed successfully'
date
echo 'Running dependent task...'
sleep 20
echo 'Dependent job complete'`, prereqResp.JobId, prereqResp.JobId),
		Partition: "compute",
		TimeLimit: 300, // 5 minutes in seconds
	}

	dependentResp, err := client.Jobs().Submit(ctx, dependentSpec)
	if err != nil {
		log.Printf("Failed to submit dependent job: %v", err)
		return
	}

	fmt.Printf("Submitted dependent job with ID: %d\n", dependentResp.JobId)
	fmt.Printf("Job %d will start after job %d completes successfully\n", dependentResp.JobId, prereqResp.JobId)

	// Monitor the jobs
	fmt.Println("\nMonitoring job progress...")
	for range 10 {
		prereqJobIDStr := fmt.Sprintf("%d", prereqResp.JobId)
		dependentJobIDStr := fmt.Sprintf("%d", dependentResp.JobId)
		prereqJob, _ := client.Jobs().Get(ctx, prereqJobIDStr)
		dependentJob, _ := client.Jobs().Get(ctx, dependentJobIDStr)

		prereqState := ""
		if len(prereqJob.JobState) > 0 {
			prereqState = string(prereqJob.JobState[0])
		}
		dependentState := ""
		if len(dependentJob.JobState) > 0 {
			dependentState = string(dependentJob.JobState[0])
		}

		fmt.Printf("  Prerequisite job %d: %s\n", prereqResp.JobId, prereqState)
		fmt.Printf("  Dependent job %d: %s\n", dependentResp.JobId, dependentState)

		if dependentState == "COMPLETED" || dependentState == "FAILED" {
			break
		}

		time.Sleep(5 * time.Second)
	}
}
