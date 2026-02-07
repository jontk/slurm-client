// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

// Commented out - unused test file
/*
import (
	"context"
	"fmt"
	"os"

	"github.com/jontk/slurm-client/internal/factory"
	"github.com/jontk/slurm-client/pkg/auth"
)

func testV0041() {
	// Get JWT token from environment or use default value
	token := os.Getenv("SLURM_JWT_TOKEN")
	if token == "" {
		// Use default value
		token = "your-jwt-token-here"
	}

	// Create auth provider
	authProvider := auth.NewTokenAuth(token)

	// Create factory
	clientFactory, err := factory.NewClientFactory(
		factory.WithAuth(authProvider),
		factory.WithBaseURL("http://localhost:6820"),
	)
	if err != nil {
		fmt.Printf("Error creating factory: %v\n", err)
		os.Exit(1)
	}

	// Create v0.0.41 client
	client, err := clientFactory.NewClientWithVersion(context.Background(), "v0.0.41")
	if err != nil {
		fmt.Printf("Error creating v0.0.41 client: %v\n", err)
		os.Exit(1)
	}

	// Test ping
	fmt.Println("Testing ping...")
	err = client.Info().Ping(context.Background())
	if err != nil {
		fmt.Printf("Ping failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Ping successful!\n")

	// Test list nodes
	fmt.Println("\nTesting list nodes...")
	nodes, err := client.Nodes().List(context.Background(), nil)
	if err != nil {
		fmt.Printf("List nodes failed: %v\n", err)
	} else {
		fmt.Printf("List nodes successful! Found %d nodes\n", len(nodes.Nodes))
		for _, node := range nodes.Nodes {
			fmt.Printf("  - %s (state: %s, cpus: %d)\n", node.Name, node.State, node.CPUs)
		}
	}

	// Test list jobs
	fmt.Println("\nTesting list jobs...")
	jobs, err := client.Jobs().List(context.Background(), nil)
	if err != nil {
		fmt.Printf("List jobs failed: %v\n", err)
	} else {
		fmt.Printf("List jobs successful! Found %d jobs\n", len(jobs.Jobs))
	}

	// Test list partitions
	fmt.Println("\nTesting list partitions...")
	partitions, err := client.Partitions().List(context.Background(), nil)
	if err != nil {
		fmt.Printf("List partitions failed: %v\n", err)
	} else {
		fmt.Printf("List partitions successful! Found %d partitions\n", len(partitions.Partitions))
	}

	fmt.Println("\nAll tests completed!")
}
*/
