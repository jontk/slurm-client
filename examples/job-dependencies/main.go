package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

// Example: Job dependencies and workflow orchestration
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

	// Example 1: Simple job chain (A -> B -> C)
	fmt.Println("=== Simple Job Chain ===")
	runSimpleJobChain(ctx, client)

	// Example 2: Fan-out workflow (A -> B1, B2, B3 -> C)
	fmt.Println("\n=== Fan-out Workflow ===")
	runFanOutWorkflow(ctx, client)

	// Example 3: Complex DAG workflow
	fmt.Println("\n=== Complex DAG Workflow ===")
	runComplexDAGWorkflow(ctx, client)

	// Example 4: Conditional workflow with error handling
	fmt.Println("\n=== Conditional Workflow ===")
	runConditionalWorkflow(ctx, client)
}

// runSimpleJobChain creates a simple linear job dependency chain
func runSimpleJobChain(ctx context.Context, client slurm.SlurmClient) {
	fmt.Println("Creating job chain: preprocessing -> processing -> postprocessing")

	// Job 1: Preprocessing
	job1 := &interfaces.JobSubmission{
		Name:      "preprocess",
		Script: `#!/bin/bash
echo "Starting preprocessing at $(date)"
# Download and prepare data
wget -q https://example.com/dataset.tar.gz
tar -xzf dataset.tar.gz
echo "Preprocessing completed"
`,
		Partition: "compute",
		CPUs:      1,
		Memory:    2048,
		TimeLimit: 10,
	}

	resp1, err := client.Jobs().Submit(ctx, job1)
	if err != nil {
		log.Fatalf("Failed to submit preprocessing job: %v", err)
	}
	fmt.Printf("Preprocessing job submitted: %s\n", resp1.JobID)

	// Job 2: Processing (depends on Job 1)
	job2 := &interfaces.JobSubmission{
		Name:      "process",
		Script: `#!/bin/bash
echo "Starting main processing at $(date)"
# Process the data
python3 /scripts/process_data.py --input dataset/ --output results/
echo "Processing completed"
`,
		Partition: "compute",
		CPUs:      8,
		Memory:    16384,
		TimeLimit: 60,
		Metadata: map[string]interface{}{
			"dependencies": []string{fmt.Sprintf("afterok:%s", resp1.JobID)},
		},
	}

	resp2, err := client.Jobs().Submit(ctx, job2)
	if err != nil {
		log.Fatalf("Failed to submit processing job: %v", err)
	}
	fmt.Printf("Processing job submitted: %s (depends on %s)\n", resp2.JobID, resp1.JobID)

	// Job 3: Postprocessing (depends on Job 2)
	job3 := &interfaces.JobSubmission{
		Name:      "postprocess",
		Script: `#!/bin/bash
echo "Starting postprocessing at $(date)"
# Generate reports and clean up
python3 /scripts/generate_report.py --input results/
rm -rf dataset/
echo "Postprocessing completed"
`,
		Partition: "compute",
		CPUs:      2,
		Memory:    4096,
		TimeLimit: 15,
		Metadata: map[string]interface{}{
			"dependencies": []string{fmt.Sprintf("afterok:%s", resp2.JobID)},
		},
	}

	resp3, err := client.Jobs().Submit(ctx, job3)
	if err != nil {
		log.Fatalf("Failed to submit postprocessing job: %v", err)
	}
	fmt.Printf("Postprocessing job submitted: %s (depends on %s)\n", resp3.JobID, resp2.JobID)

	// Monitor the chain
	monitorJobChain(ctx, client, []string{resp1.JobID, resp2.JobID, resp3.JobID})
}

// runFanOutWorkflow creates a fan-out/fan-in workflow pattern
func runFanOutWorkflow(ctx context.Context, client slurm.SlurmClient) {
	fmt.Println("Creating fan-out workflow: setup -> parallel tasks -> merge")

	// Setup job
	setupJob := &interfaces.JobSubmission{
		Name:      "setup-fanout",
		Script: `#!/bin/bash
echo "Setting up parallel workflow"
# Split data into chunks for parallel processing
split -n 3 input.data chunk_
echo "Setup completed - data split into 3 chunks"
`,
		Partition: "compute",
		CPUs:      1,
		Memory:    1024,
		TimeLimit: 5,
	}

	setupResp, err := client.Jobs().Submit(ctx, setupJob)
	if err != nil {
		log.Fatalf("Failed to submit setup job: %v", err)
	}
	fmt.Printf("Setup job submitted: %s\n", setupResp.JobID)

	// Submit parallel processing jobs
	var parallelJobIDs []string
	for i := 0; i < 3; i++ {
		parallelJob := &interfaces.JobSubmission{
			Name:      fmt.Sprintf("parallel-task-%d", i),
			Script: fmt.Sprintf(`#!/bin/bash
echo "Processing chunk %d"
# Process specific chunk
python3 process_chunk.py --input chunk_%c --output result_%d.out
echo "Chunk %d processing completed"
`, i, 'a'+i, i, i),
			Partition: "compute",
			CPUs:      4,
			Memory:    8192,
			TimeLimit: 30,
			Metadata: map[string]interface{}{
				"dependencies": []string{fmt.Sprintf("afterok:%s", setupResp.JobID)},
			},
		}

		resp, err := client.Jobs().Submit(ctx, parallelJob)
		if err != nil {
			log.Printf("Failed to submit parallel job %d: %v", i, err)
			continue
		}
		parallelJobIDs = append(parallelJobIDs, resp.JobID)
		fmt.Printf("Parallel job %d submitted: %s\n", i, resp.JobID)
	}

	// Merge job (depends on all parallel jobs)
	dependencies := make([]string, len(parallelJobIDs))
	for i, jobID := range parallelJobIDs {
		dependencies[i] = fmt.Sprintf("afterok:%s", jobID)
	}

	mergeJob := &interfaces.JobSubmission{
		Name:      "merge-results",
		Script: `#!/bin/bash
echo "Merging results from parallel jobs"
cat result_*.out > final_result.out
# Generate summary report
python3 summarize.py --input final_result.out --output summary.pdf
echo "Merge completed"
`,
		Partition: "compute",
		CPUs:      2,
		Memory:    4096,
		TimeLimit: 10,
		Metadata: map[string]interface{}{
			"dependencies": dependencies,
		},
	}

	mergeResp, err := client.Jobs().Submit(ctx, mergeJob)
	if err != nil {
		log.Fatalf("Failed to submit merge job: %v", err)
	}
	fmt.Printf("Merge job submitted: %s (depends on %d parallel jobs)\n", 
		mergeResp.JobID, len(parallelJobIDs))
}

// runComplexDAGWorkflow creates a complex Directed Acyclic Graph workflow
func runComplexDAGWorkflow(ctx context.Context, client slurm.SlurmClient) {
	fmt.Println("Creating complex DAG workflow")
	
	// Define workflow stages
	type workflowJob struct {
		name         string
		script       string
		cpus         int
		dependencies []string
		jobID        string
	}
	
	jobs := map[string]*workflowJob{
		"data-fetch": {
			name: "data-fetch",
			script: `#!/bin/bash
echo "Fetching data from multiple sources"
wget -q https://source1.com/data1.csv
wget -q https://source2.com/data2.csv
`,
			cpus: 1,
		},
		"data-clean": {
			name: "data-clean",
			script: `#!/bin/bash
echo "Cleaning and validating data"
python3 clean_data.py --input *.csv --output cleaned/
`,
			cpus:         2,
			dependencies: []string{"data-fetch"},
		},
		"feature-eng": {
			name: "feature-eng",
			script: `#!/bin/bash
echo "Engineering features"
python3 feature_engineering.py --input cleaned/ --output features/
`,
			cpus:         4,
			dependencies: []string{"data-clean"},
		},
		"train-model1": {
			name: "train-model1",
			script: `#!/bin/bash
echo "Training model 1 (Random Forest)"
python3 train_rf.py --features features/ --output models/rf.pkl
`,
			cpus:         8,
			dependencies: []string{"feature-eng"},
		},
		"train-model2": {
			name: "train-model2",
			script: `#!/bin/bash
echo "Training model 2 (Neural Network)"
python3 train_nn.py --features features/ --output models/nn.h5
`,
			cpus:         16,
			dependencies: []string{"feature-eng"},
		},
		"evaluate": {
			name: "evaluate",
			script: `#!/bin/bash
echo "Evaluating models"
python3 evaluate.py --models models/ --output evaluation/
`,
			cpus:         4,
			dependencies: []string{"train-model1", "train-model2"},
		},
		"deploy": {
			name: "deploy",
			script: `#!/bin/bash
echo "Deploying best model"
python3 deploy.py --evaluation evaluation/ --models models/
`,
			cpus:         2,
			dependencies: []string{"evaluate"},
		},
	}
	
	// Submit jobs in dependency order
	for _, job := range jobs {
		// Build dependency list
		var deps []string
		for _, depName := range job.dependencies {
			if depJob, ok := jobs[depName]; ok && depJob.jobID != "" {
				deps = append(deps, fmt.Sprintf("afterok:%s", depJob.jobID))
			}
		}
		
		// Submit job
		submission := &interfaces.JobSubmission{
			Name:      job.name,
			Script:    job.script,
			Partition: "compute",
			CPUs:      job.cpus,
			Memory:    job.cpus * 2048, // 2GB per CPU
			TimeLimit: 30,
		}
		if len(deps) > 0 {
			submission.Metadata = map[string]interface{}{
				"dependencies": deps,
			}
		}
		
		resp, err := client.Jobs().Submit(ctx, submission)
		if err != nil {
			log.Printf("Failed to submit job %s: %v", job.name, err)
			continue
		}
		
		job.jobID = resp.JobID
		fmt.Printf("Submitted %s: %s", job.name, job.jobID)
		if len(deps) > 0 {
			fmt.Printf(" (depends on %d jobs)", len(deps))
		}
		fmt.Println()
	}
}

// runConditionalWorkflow creates a workflow with conditional execution
func runConditionalWorkflow(ctx context.Context, client slurm.SlurmClient) {
	fmt.Println("Creating conditional workflow with error handling")
	
	// Validation job
	validationJob := &interfaces.JobSubmission{
		Name: "validate-input",
		Script: `#!/bin/bash
echo "Validating input data"
if [ ! -f input.data ]; then
    echo "ERROR: Input file not found"
    exit 1
fi
if [ $(wc -l < input.data) -lt 1000 ]; then
    echo "WARNING: Input file has less than 1000 lines"
    exit 2
fi
echo "Validation passed"
exit 0
`,
		Partition: "compute",
		CPUs:      1,
		Memory:    1024,
		TimeLimit: 5,
	}
	
	valResp, err := client.Jobs().Submit(ctx, validationJob)
	if err != nil {
		log.Fatalf("Failed to submit validation job: %v", err)
	}
	fmt.Printf("Validation job submitted: %s\n", valResp.JobID)
	
	// Success path - runs only if validation succeeds
	successJob := &interfaces.JobSubmission{
		Name: "process-valid-data",
		Script: `#!/bin/bash
echo "Processing validated data"
python3 process.py --input input.data --output output.data
echo "Processing completed successfully"
`,
		Partition: "compute",
		CPUs:      8,
		Memory:    16384,
		TimeLimit: 60,
		Metadata: map[string]interface{}{
			"dependencies": []string{fmt.Sprintf("afterok:%s", valResp.JobID)},
		},
	}
	
	successResp, err := client.Jobs().Submit(ctx, successJob)
	if err != nil {
		log.Printf("Failed to submit success job: %v", err)
	} else {
		fmt.Printf("Success job submitted: %s (runs on validation success)\n", successResp.JobID)
	}
	
	// Failure path - runs only if validation fails
	failureJob := &interfaces.JobSubmission{
		Name: "handle-invalid-data",
		Script: `#!/bin/bash
echo "Handling validation failure"
# Send notification
mail -s "Data validation failed" admin@example.com < /dev/null
# Try to fetch valid data
wget -q https://backup.example.com/valid_input.data -O input.data
echo "Fetched backup data"
`,
		Partition: "compute",
		CPUs:      1,
		Memory:    1024,
		TimeLimit: 10,
		Metadata: map[string]interface{}{
			"dependencies": []string{fmt.Sprintf("afternotok:%s", valResp.JobID)},
		},
	}
	
	failureResp, err := client.Jobs().Submit(ctx, failureJob)
	if err != nil {
		log.Printf("Failed to submit failure job: %v", err)
	} else {
		fmt.Printf("Failure job submitted: %s (runs on validation failure)\n", failureResp.JobID)
	}
	
	// Cleanup job - runs regardless of validation result
	cleanupJob := &interfaces.JobSubmission{
		Name: "cleanup",
		Script: `#!/bin/bash
echo "Performing cleanup"
# Archive logs
tar -czf logs_$(date +%Y%m%d).tar.gz *.log
# Remove temporary files
rm -f /tmp/work_*
echo "Cleanup completed"
`,
		Partition: "compute",
		CPUs:      1,
		Memory:    1024,
		TimeLimit: 5,
		Metadata: map[string]interface{}{
			"dependencies": []string{fmt.Sprintf("afterany:%s", valResp.JobID)},
		},
	}
	
	cleanupResp, err := client.Jobs().Submit(ctx, cleanupJob)
	if err != nil {
		log.Printf("Failed to submit cleanup job: %v", err)
	} else {
		fmt.Printf("Cleanup job submitted: %s (runs after validation regardless of result)\n", cleanupResp.JobID)
	}
}

// monitorJobChain monitors a chain of dependent jobs
func monitorJobChain(ctx context.Context, client slurm.SlurmClient, jobIDs []string) {
	fmt.Println("\nMonitoring job chain...")
	
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	timeout := time.After(10 * time.Minute)
	
	for {
		select {
		case <-ticker.C:
			allComplete := true
			
			for i, jobID := range jobIDs {
				job, err := client.Jobs().Get(ctx, jobID)
				if err != nil {
					log.Printf("Failed to get job %s: %v", jobID, err)
					continue
				}
				
				fmt.Printf("Job %d (%s): %s", i+1, job.Name, job.State)
				
				if job.State == "PENDING" && job.Reason != "" {
					fmt.Printf(" - %s", job.Reason)
				}
				
				if job.State != "COMPLETED" && job.State != "FAILED" && job.State != "CANCELLED" {
					allComplete = false
				}
				
				fmt.Println()
			}
			
			if allComplete {
				fmt.Println("All jobs in chain completed!")
				return
			}
			
			fmt.Println("---")
			
		case <-timeout:
			fmt.Println("Monitoring timeout reached")
			return
			
		case <-ctx.Done():
			fmt.Println("Context cancelled")
			return
		}
	}
}