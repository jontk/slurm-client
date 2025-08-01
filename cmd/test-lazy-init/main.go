// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"

	slurm "github.com/jontk/slurm-client"
)

func main() {
	ctx := context.Background()

	// Create client without specifying version (should use latest)
	// Using WithBaseURL option to specify the test URL
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://test.example.com"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Printf("Client created successfully with version: %s\n", client.Version())

	// Test lazy initialization by accessing managers
	fmt.Println("\nTesting lazy initialization...")

	// Access Jobs manager - should be lazily initialized
	jobsManager := client.Jobs()
	fmt.Printf("Jobs manager type: %T\n", jobsManager)

	// Access Nodes manager
	nodesManager := client.Nodes()
	fmt.Printf("Nodes manager type: %T\n", nodesManager)

	// Access Accounts manager
	accountsManager := client.Accounts()
	fmt.Printf("Accounts manager type: %T\n", accountsManager)

	// Try to use a manager method (will fail due to no real server, but shows lazy init works)
	err = client.Info().Ping(ctx)
	if err != nil {
		fmt.Printf("\nPing failed as expected (no real server): %v\n", err)
	}

	fmt.Println("\nLazy initialization test completed successfully!")
}
