// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

// Example: Job Allocation API Usage
// The Allocate API allows you to request resource allocation without
// submitting a script, useful for interactive sessions and dynamic workloads
func main() {
	// Get configuration from environment or use defaults
	baseURL := os.Getenv("SLURM_BASE_URL")
	if baseURL == "" {
		baseURL = "https://cluster.example.com:6820"
	}
	
	token := os.Getenv("SLURM_TOKEN")
	if token == "" {
		log.Fatal("Please set SLURM_TOKEN environment variable")
	}

	// Create configuration for v0.0.43 (Allocate API support)
	cfg := config.NewDefault()
	cfg.BaseURL = baseURL
	cfg.APIVersion = "v0.0.43" // Job allocation requires v0.0.43+
	
	// Create authentication
	authProvider := auth.NewTokenAuth(token)

	// Create client
	ctx := context.Background()
	client, err := slurm.NewClient(ctx,
		slurm.WithConfig(cfg),
		slurm.WithAuth(authProvider),
		slurm.WithAPIVersion("v0.0.43"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	jobManager := client.Jobs()

	// Example 1: Basic resource allocation
	fmt.Println("=== Basic Resource Allocation ===")
	basicAllocation(ctx, jobManager)

	// Example 2: GPU allocation for interactive work
	fmt.Println("\n=== GPU Resource Allocation ===")
	gpuAllocation(ctx, jobManager)

	// Example 3: Memory-intensive allocation
	fmt.Println("\n=== High Memory Allocation ===")
	highMemoryAllocation(ctx, jobManager)

	// Example 4: Allocation with specific constraints
	fmt.Println("\n=== Constrained Allocation ===")
	constrainedAllocation(ctx, jobManager)

	// Example 5: Allocation workflow
	fmt.Println("\n=== Allocation Workflow ===")
	allocationWorkflow(ctx, jobManager)
}

func basicAllocation(ctx context.Context, jobManager interfaces.JobManager) {
	// Request basic compute resources
	allocReq := &types.JobAllocateRequest{
		Name:      "interactive-session",
		Account:   "research",
		Partition: "compute",
		Nodes:     "1",        // Request 1 node
		CPUs:      4,          // Request 4 CPUs
		Memory:    "8G",       // Request 8GB memory
		TimeLimit: "2:00:00",  // 2 hour time limit
	}

	fmt.Println("Requesting allocation:")
	fmt.Printf("  - Nodes: %s\n", allocReq.Nodes)
	fmt.Printf("  - CPUs: %d\n", allocReq.CPUs)
	fmt.Printf("  - Memory: %s\n", allocReq.Memory)
	fmt.Printf("  - Time: %s\n", allocReq.TimeLimit)

	resp, err := jobManager.Allocate(ctx, allocReq)
	if err != nil {
		log.Printf("Failed to allocate resources: %v", err)
		return
	}

	fmt.Printf("\n✓ Allocation successful!\n")
	fmt.Printf("  Job ID: %d\n", resp.JobID)
	fmt.Printf("  Status: %s\n", resp.Status)
	if resp.NodeList != "" {
		fmt.Printf("  Allocated nodes: %s\n", resp.NodeList)
	}
	if resp.Message != "" {
		fmt.Printf("  Message: %s\n", resp.Message)
	}
}

func gpuAllocation(ctx context.Context, jobManager interfaces.JobManager) {
	// Request GPU resources for interactive ML development
	allocReq := &types.JobAllocateRequest{
		Name:      "ml-interactive",
		Account:   "ml-research",
		Partition: "gpu",
		Nodes:     "1",
		CPUs:      8,
		Memory:    "32G",
		GPUs:      2,              // Request 2 GPUs
		TimeLimit: "4:00:00",
		QoS:       "interactive",  // Use interactive QoS for better response
		Features:  []string{"gpu_tesla_v100"}, // Request specific GPU type
	}

	fmt.Println("Requesting GPU allocation:")
	fmt.Printf("  - GPUs: %d\n", allocReq.GPUs)
	fmt.Printf("  - GPU Type: %v\n", allocReq.Features)
	fmt.Printf("  - CPUs: %d\n", allocReq.CPUs)
	fmt.Printf("  - Memory: %s\n", allocReq.Memory)

	resp, err := jobManager.Allocate(ctx, allocReq)
	if err != nil {
		log.Printf("Failed to allocate GPU resources: %v", err)
		return
	}

	fmt.Printf("\n✓ GPU allocation successful!\n")
	fmt.Printf("  Job ID: %d\n", resp.JobID)
	fmt.Printf("  Allocated nodes: %s\n", resp.NodeList)
	
	// In a real scenario, you would now connect to the allocated node
	// and start your interactive GPU workload
	fmt.Println("\nYou can now SSH to the allocated node and use the GPUs")
	fmt.Printf("Example: ssh %s\n", resp.NodeList)
}

func highMemoryAllocation(ctx context.Context, jobManager interfaces.JobManager) {
	// Request high-memory resources for data processing
	allocReq := &types.JobAllocateRequest{
		Name:         "bigdata-processing",
		Account:      "data-science",
		Partition:    "highmem",
		Nodes:        "1",
		CPUs:         16,
		Memory:       "256G",      // Request 256GB memory
		TimeLimit:    "8:00:00",
		MinMemPerCPU: "16G",       // Ensure 16GB per CPU
		Exclusive:    true,        // Request exclusive node access
	}

	fmt.Println("Requesting high-memory allocation:")
	fmt.Printf("  - Memory: %s\n", allocReq.Memory)
	fmt.Printf("  - CPUs: %d\n", allocReq.CPUs)
	fmt.Printf("  - Exclusive: %v\n", allocReq.Exclusive)

	resp, err := jobManager.Allocate(ctx, allocReq)
	if err != nil {
		log.Printf("Failed to allocate high-memory resources: %v", err)
		return
	}

	fmt.Printf("\n✓ High-memory allocation successful!\n")
	fmt.Printf("  Job ID: %d\n", resp.JobID)
	fmt.Printf("  Total allocated CPUs: %d\n", resp.AllocatedCPUs)
	fmt.Printf("  Total allocated memory: %d MB\n", resp.AllocatedMemory)
}

func constrainedAllocation(ctx context.Context, jobManager interfaces.JobManager) {
	// Request allocation with specific constraints
	allocReq := &types.JobAllocateRequest{
		Name:      "specific-node-job",
		Account:   "special-project",
		Partition: "compute",
		// Request specific nodes by name
		NodeList:  "node[001-004]",
		CPUs:      32,
		Memory:    "64G",
		TimeLimit: "1:00:00",
		// Additional constraints
		Features:     []string{"infiniband", "ssd"},
		Constraints:  "cpu_family:intel&cpu_model:skylake",
		Reservation:  "maintenance-window", // Use a specific reservation
		Dependencies: []string{},           // No dependencies for allocation
	}

	fmt.Println("Requesting constrained allocation:")
	fmt.Printf("  - Node list: %s\n", allocReq.NodeList)
	fmt.Printf("  - Features: %v\n", allocReq.Features)
	fmt.Printf("  - Constraints: %s\n", allocReq.Constraints)
	fmt.Printf("  - Reservation: %s\n", allocReq.Reservation)

	resp, err := jobManager.Allocate(ctx, allocReq)
	if err != nil {
		log.Printf("Failed to allocate with constraints: %v", err)
		return
	}

	fmt.Printf("\n✓ Constrained allocation successful!\n")
	fmt.Printf("  Job ID: %d\n", resp.JobID)
	fmt.Printf("  Allocated from nodes: %s\n", resp.NodeList)
}

func allocationWorkflow(ctx context.Context, jobManager interfaces.JobManager) {
	// Demonstrate a complete allocation workflow
	fmt.Println("Starting allocation workflow...")

	// Step 1: Request initial allocation
	allocReq := &types.JobAllocateRequest{
		Name:      "workflow-master",
		Account:   "workflow",
		Partition: "compute",
		Nodes:     "1",
		CPUs:      2,
		Memory:    "4G",
		TimeLimit: "30:00",
		WCKey:     "workflow-orchestration", // Track with WCKey
	}

	resp, err := jobManager.Allocate(ctx, allocReq)
	if err != nil {
		log.Printf("Failed to allocate master node: %v", err)
		return
	}

	masterJobID := resp.JobID
	fmt.Printf("✓ Master allocation (Job ID: %d) on node: %s\n", 
		masterJobID, resp.NodeList)

	// Step 2: Check allocation status
	time.Sleep(2 * time.Second) // Wait for allocation to be fully ready
	
	job, err := jobManager.Get(ctx, fmt.Sprintf("%d", masterJobID))
	if err != nil {
		log.Printf("Failed to get job status: %v", err)
		return
	}

	fmt.Printf("  Allocation state: %s\n", job.State)
	fmt.Printf("  Start time: %v\n", job.StartTime)

	// Step 3: Based on master allocation, request worker allocations
	if job.State == "RUNNING" {
		fmt.Println("\nMaster is running, allocating workers...")
		
		// Allocate multiple worker nodes
		for i := 1; i <= 3; i++ {
			workerReq := &types.JobAllocateRequest{
				Name:         fmt.Sprintf("workflow-worker-%d", i),
				Account:      "workflow",
				Partition:    "compute",
				Nodes:        "1",
				CPUs:         8,
				Memory:       "16G",
				TimeLimit:    "25:00", // Slightly less than master
				WCKey:        "workflow-worker",
				Dependency:   fmt.Sprintf("after:%d", masterJobID), // Depend on master
			}

			workerResp, err := jobManager.Allocate(ctx, workerReq)
			if err != nil {
				log.Printf("Failed to allocate worker %d: %v", i, err)
				continue
			}

			fmt.Printf("  ✓ Worker %d allocated (Job ID: %d)\n", i, workerResp.JobID)
		}
	}

	// Step 4: Monitor allocations
	fmt.Println("\nMonitoring allocations...")
	jobOpts := &interfaces.ListJobsOptions{
		UserID:    os.Getenv("USER"),
		States:    []string{"RUNNING", "PENDING"},
		StartTime: time.Now().Add(-5 * time.Minute),
	}

	jobs, err := jobManager.List(ctx, jobOpts)
	if err != nil {
		log.Printf("Failed to list jobs: %v", err)
		return
	}

	fmt.Printf("Active allocations: %d\n", len(jobs.Jobs))
	for _, j := range jobs.Jobs {
		if j.Name != "" && (j.Name == "workflow-master" || 
			len(j.Name) > 15 && j.Name[:15] == "workflow-worker") {
			fmt.Printf("  - %s (ID: %s): %s on %s\n", 
				j.Name, j.ID, j.State, j.NodeList)
		}
	}

	// In a real workflow, you would:
	// 1. SSH to allocated nodes
	// 2. Start your distributed application
	// 3. Monitor progress
	// 4. Release allocations when done
	
	fmt.Println("\nWorkflow allocation complete!")
	fmt.Println("Remember to release allocations when done:")
	fmt.Printf("  scancel %d  # Cancel master and dependent workers\n", masterJobID)
}

// Helper function to wait for allocation to be ready
func waitForAllocation(ctx context.Context, jobManager interfaces.JobManager, 
	jobID string, timeout time.Duration) error {
	
	deadline := time.Now().Add(timeout)
	
	for time.Now().Before(deadline) {
		job, err := jobManager.Get(ctx, jobID)
		if err != nil {
			return fmt.Errorf("failed to check job status: %w", err)
		}
		
		switch job.State {
		case "RUNNING":
			return nil // Allocation is ready
		case "FAILED", "CANCELLED", "TIMEOUT":
			return fmt.Errorf("allocation failed with state: %s", job.State)
		case "PENDING":
			// Still waiting for resources
			time.Sleep(5 * time.Second)
		default:
			time.Sleep(2 * time.Second)
		}
	}
	
	return fmt.Errorf("allocation timeout after %v", timeout)
}

// Example: Release an allocation
func releaseAllocation(ctx context.Context, jobManager interfaces.JobManager, jobID string) {
	fmt.Printf("Releasing allocation %s...\n", jobID)
	
	err := jobManager.Cancel(ctx, jobID)
	if err != nil {
		log.Printf("Failed to release allocation: %v", err)
		return
	}
	
	fmt.Println("Allocation released successfully")
}