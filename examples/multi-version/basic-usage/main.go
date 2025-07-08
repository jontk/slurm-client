package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

func main() {
	ctx := context.Background()

	// Example 1: Automatic version detection
	fmt.Println("=== Example 1: Automatic Version Detection ===")
	client, err := createClientWithAutoDetection(ctx)
	if err != nil {
		log.Printf("Auto-detection failed: %v", err)
	} else {
		printBasicInfo(ctx, client)
	}

	// Example 2: Explicit version selection
	fmt.Println("\n=== Example 2: Explicit Version Selection ===")
	for _, version := range slurm.SupportedVersions() {
		fmt.Printf("Testing version %s:\n", version)
		client, err := slurm.NewClientWithVersion(ctx, version, 
			slurm.WithBaseURL("https://localhost:6820"),
			slurm.WithAuth(auth.NewNoneAuth()),
		)
		if err != nil {
			log.Printf("  Error creating client for %s: %v", version, err)
			continue
		}
		
		// Test basic connectivity
		info, err := client.GetInfo(ctx)
		if err != nil {
			log.Printf("  Error getting info: %v", err)
		} else {
			fmt.Printf("  Connected! Server version: %s\n", info.Version)
		}
	}

	// Example 3: Version selection for specific Slurm releases
	fmt.Println("\n=== Example 3: Slurm Release Compatibility ===")
	slurmVersions := []string{"24.05", "24.11", "25.05"}
	for _, slurmVersion := range slurmVersions {
		fmt.Printf("Best API version for Slurm %s: ", slurmVersion)
		client, err := slurm.NewClientForSlurmVersion(ctx, slurmVersion,
			slurm.WithBaseURL("https://localhost:6820"),
			slurm.WithAuth(auth.NewNoneAuth()),
		)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Success (using API version based on compatibility)\n")
			_ = client // Use client for operations
		}
	}

	// Example 4: Configuration-based client creation
	fmt.Println("\n=== Example 4: Configuration-Based Setup ===")
	cfg := &config.Config{
		BaseURL: "https://localhost:6820",
		Timeout: "30s",
		Auth: config.AuthConfig{
			Type: "token",
			Token: os.Getenv("SLURM_TOKEN"),
		},
		Retry: config.RetryConfig{
			MaxRetries: 3,
			BaseDelay:  "1s",
			MaxDelay:   "10s",
		},
	}

	client, err = slurm.NewClientFromConfig(ctx, cfg)
	if err != nil {
		log.Printf("Config-based client creation failed: %v", err)
	} else {
		fmt.Println("Config-based client created successfully")
		printBasicInfo(ctx, client)
	}
}

func createClientWithAutoDetection(ctx context.Context) (slurm.SlurmClient, error) {
	// Configure client with automatic version detection
	return slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(auth.NewNoneAuth()),
		slurm.WithTimeout("30s"),
		slurm.WithRetryPolicy(3, "1s", "10s"),
		slurm.WithDebug(true),
	)
}

func printBasicInfo(ctx context.Context, client slurm.SlurmClient) {
	// Get server information
	info, err := client.GetInfo(ctx)
	if err != nil {
		log.Printf("  Error getting server info: %v", err)
		return
	}

	fmt.Printf("  Server Version: %s\n", info.Version)
	fmt.Printf("  Release: %s\n", info.Release)
	
	// List available partitions
	partitions, err := client.ListPartitions(ctx)
	if err != nil {
		log.Printf("  Error listing partitions: %v", err)
		return
	}

	fmt.Printf("  Available Partitions: %d\n", len(partitions))
	for _, partition := range partitions {
		fmt.Printf("    - %s (nodes: %d)\n", partition.Name, len(partition.Nodes))
	}
}