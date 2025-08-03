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

// Example: WCKey (Workload Characterization Key) Management
// WCKeys are used to categorize and track different types of workloads
// in SLURM for accounting and resource management purposes.
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

	// Create configuration for v0.0.43 (WCKey support)
	cfg := config.NewDefault()
	cfg.BaseURL = baseURL
	cfg.APIVersion = "v0.0.43" // WCKey management requires v0.0.43+
	
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

	// Get WCKey manager (only available in v0.0.43+)
	wcKeyManager := client.WCKeys()
	if wcKeyManager == nil {
		log.Fatal("WCKey management not available. Ensure you're using API v0.0.43 or later")
	}

	// Example 1: List existing WCKeys
	fmt.Println("=== Listing WCKeys ===")
	listWCKeys(ctx, wcKeyManager)

	// Example 2: Create new WCKeys for different workload types
	fmt.Println("\n=== Creating WCKeys ===")
	createWorkloadKeys(ctx, wcKeyManager)

	// Example 3: Filter WCKeys by user or cluster
	fmt.Println("\n=== Filtering WCKeys ===")
	filterWCKeys(ctx, wcKeyManager)

	// Example 4: Using WCKeys in job submission
	fmt.Println("\n=== Using WCKeys in Jobs ===")
	demonstrateWCKeyUsage(ctx, client)
}

func listWCKeys(ctx context.Context, wcKeyManager interfaces.WCKeyManager) {
	// List all WCKeys
	wcKeys, err := wcKeyManager.List(ctx, nil)
	if err != nil {
		log.Printf("Failed to list WCKeys: %v", err)
		return
	}

	fmt.Printf("Found %d WCKeys:\n", len(wcKeys.WCKeys))
	for _, wckey := range wcKeys.WCKeys {
		fmt.Printf("  - ID: %s, Name: %s, User: %s, Cluster: %s\n",
			wckey.ID, wckey.Name, wckey.User, wckey.Cluster)
	}
}

func createWorkloadKeys(ctx context.Context, wcKeyManager interfaces.WCKeyManager) {
	// Define different workload types
	workloadTypes := []struct {
		name        string
		description string
		user        string
		cluster     string
	}{
		{
			name:        "ml-training",
			description: "Machine Learning Training Jobs",
			user:        "researcher1",
			cluster:     "gpu-cluster",
		},
		{
			name:        "data-preprocessing",
			description: "Data Preprocessing and ETL",
			user:        "researcher1",
			cluster:     "gpu-cluster",
		},
		{
			name:        "simulation",
			description: "Scientific Simulations",
			user:        "scientist1",
			cluster:     "hpc-cluster",
		},
		{
			name:        "analysis",
			description: "Data Analysis and Visualization",
			user:        "analyst1",
			cluster:     "general",
		},
	}

	for _, wt := range workloadTypes {
		wckey := &types.WCKeyCreate{
			Name:    wt.name,
			User:    wt.user,
			Cluster: wt.cluster,
		}

		fmt.Printf("Creating WCKey '%s' for user '%s' on cluster '%s'...\n", 
			wt.name, wt.user, wt.cluster)

		resp, err := wcKeyManager.Create(ctx, wckey)
		if err != nil {
			log.Printf("  Failed to create WCKey: %v", err)
			continue
		}

		fmt.Printf("  ✓ Created successfully: %s\n", resp.Message)
	}
}

func filterWCKeys(ctx context.Context, wcKeyManager interfaces.WCKeyManager) {
	// Example 1: Filter by specific users
	fmt.Println("WCKeys for specific users:")
	userOpts := &types.WCKeyListOptions{
		Users: []string{"researcher1", "scientist1"},
	}
	
	userWCKeys, err := wcKeyManager.List(ctx, userOpts)
	if err != nil {
		log.Printf("Failed to filter by users: %v", err)
		return
	}

	for _, wckey := range userWCKeys.WCKeys {
		fmt.Printf("  - User: %s, WCKey: %s\n", wckey.User, wckey.Name)
	}

	// Example 2: Filter by cluster
	fmt.Println("\nWCKeys for GPU cluster:")
	clusterOpts := &types.WCKeyListOptions{
		Clusters: []string{"gpu-cluster"},
	}
	
	clusterWCKeys, err := wcKeyManager.List(ctx, clusterOpts)
	if err != nil {
		log.Printf("Failed to filter by cluster: %v", err)
		return
	}

	for _, wckey := range clusterWCKeys.WCKeys {
		fmt.Printf("  - Cluster: %s, WCKey: %s, User: %s\n", 
			wckey.Cluster, wckey.Name, wckey.User)
	}

	// Example 3: Get only default WCKeys
	fmt.Println("\nDefault WCKeys only:")
	defaultOpts := &types.WCKeyListOptions{
		OnlyDefaults: true,
	}
	
	defaultWCKeys, err := wcKeyManager.List(ctx, defaultOpts)
	if err != nil {
		log.Printf("Failed to filter defaults: %v", err)
		return
	}

	for _, wckey := range defaultWCKeys.WCKeys {
		fmt.Printf("  - Default WCKey: %s for user %s\n", wckey.Name, wckey.User)
	}
}

func demonstrateWCKeyUsage(ctx context.Context, client interfaces.SlurmClient) {
	// Show how to use WCKeys when submitting jobs
	jobManager := client.Jobs()

	// Create a job with specific WCKey for workload tracking
	job := &interfaces.JobSubmission{
		Name:      "ML Training Job",
		Command:   "python train_model.py",
		Partition: "gpu",
		CPUs:      8,
		Memory:    32 * 1024 * 1024 * 1024, // 32GB
		TimeLimit: 120,                     // 2 hours
		WCKey:     "ml-training",           // Specify the workload type
		Environment: map[string]string{
			"CUDA_VISIBLE_DEVICES": "0,1",
		},
	}

	fmt.Println("Submitting job with WCKey 'ml-training'...")
	resp, err := jobManager.Submit(ctx, job)
	if err != nil {
		log.Printf("Failed to submit job: %v", err)
		return
	}

	fmt.Printf("Job submitted successfully with ID: %s\n", resp.JobID)
	fmt.Printf("The job will be tracked under WCKey: ml-training\n")

	// You can later query jobs by WCKey for accounting
	fmt.Println("\nQuerying jobs by WCKey...")
	jobOpts := &interfaces.ListJobsOptions{
		WCKey:     "ml-training",
		StartTime: time.Now().Add(-24 * time.Hour), // Last 24 hours
	}

	jobs, err := jobManager.List(ctx, jobOpts)
	if err != nil {
		log.Printf("Failed to list jobs by WCKey: %v", err)
		return
	}

	fmt.Printf("Found %d jobs with WCKey 'ml-training' in the last 24 hours\n", len(jobs.Jobs))
	for _, j := range jobs.Jobs {
		fmt.Printf("  - Job %s: %s (State: %s)\n", j.ID, j.Name, j.State)
	}
}

// Example: Cleanup specific WCKey
func cleanupWCKey(ctx context.Context, wcKeyManager interfaces.WCKeyManager, wcKeyID string) {
	fmt.Printf("Deleting WCKey with ID: %s\n", wcKeyID)
	
	err := wcKeyManager.Delete(ctx, wcKeyID)
	if err != nil {
		log.Printf("Failed to delete WCKey: %v", err)
		return
	}
	
	fmt.Println("WCKey deleted successfully")
}

// Example: Get detailed information about a specific WCKey
func getWCKeyDetails(ctx context.Context, wcKeyManager interfaces.WCKeyManager, wcKeyID string) {
	wckey, err := wcKeyManager.Get(ctx, wcKeyID)
	if err != nil {
		log.Printf("Failed to get WCKey details: %v", err)
		return
	}

	fmt.Printf("WCKey Details:\n")
	fmt.Printf("  ID:      %s\n", wckey.ID)
	fmt.Printf("  Name:    %s\n", wckey.Name)
	fmt.Printf("  User:    %s\n", wckey.User)
	fmt.Printf("  Cluster: %s\n", wckey.Cluster)
	fmt.Printf("  Active:  %v\n", wckey.Active)
	
	// Print any additional metadata
	if len(wckey.Meta) > 0 {
		fmt.Printf("  Metadata:\n")
		for k, v := range wckey.Meta {
			fmt.Printf("    %s: %v\n", k, v)
		}
	}
}