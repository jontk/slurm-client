// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/internal/factory"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <version>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Supported versions: v0.0.40, v0.0.41, v0.0.42, v0.0.43\n")
		os.Exit(1)
	}

	version := os.Args[1]

	// Get JWT token from environment
	jwtToken := os.Getenv("SLURM_JWT")
	if jwtToken == "" {
		log.Fatal("SLURM_JWT environment variable is required")
	}

	// Create configuration
	cfg := config.NewDefault()
	cfg.BaseURL = "http://localhost:6820"
	cfg.Debug = false

	// Create JWT authentication provider
	authProvider := auth.NewTokenAuth(jwtToken)

	fmt.Println("===========================================")
	fmt.Println("Testing SLURM Client Implementations")
	fmt.Println("===========================================")

	// Test adapter implementation
	fmt.Printf("\n2. Testing ADAPTER implementation with %s\n", version)
	fmt.Println("-------------------------------------------")
	testImplementation(cfg, authProvider, version)
}

func testImplementation(cfg *config.Config, authProvider auth.Provider, version string) {
	// Create factory with adapter option
	clientFactory, err := factory.NewClientFactory(
		factory.WithConfig(cfg),
		factory.WithAuth(authProvider),
		factory.WithBaseURL(cfg.BaseURL),
	)
	if err != nil {
		log.Printf("❌ Failed to create factory: %v", err)
		return
	}

	// Create client with specific version
	client, err := clientFactory.NewClientWithVersion(context.Background(), version)
	if err != nil {
		log.Printf("❌ Failed to create client: %v", err)
		return
	}

	fmt.Printf("✅ Successfully created %s client using adapter implementation!\n", version)

	ctx := context.Background()

	// Test 1: Ping
	fmt.Print("  • Testing Ping... ")
	err = client.Info().Ping(ctx)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Success")
	}

	// Test 2: List Jobs
	fmt.Print("  • Testing List Jobs... ")
	jobs, err := client.Jobs().List(ctx, &types.ListJobsOptions{
		Limit: 3,
	})
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Printf("✅ Success (found %d jobs)\n", jobs.Total)
	}

	// Test 3: List Nodes
	fmt.Print("  • Testing List Nodes... ")
	nodes, err := client.Nodes().List(ctx, &types.ListNodesOptions{
		Limit: 3,
	})
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Printf("✅ Success (found %d nodes)\n", nodes.Total)
	}

	// Test 4: List Partitions
	fmt.Print("  • Testing List Partitions... ")
	partitions, err := client.Partitions().List(ctx, &types.ListPartitionsOptions{
		Limit: 3,
	})
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Printf("✅ Success (found %d partitions)\n", partitions.Total)
	}

	// Test 5: Get API Version
	fmt.Print("  • Testing Get API Version... ")
	apiVersion, err := client.Info().Version(ctx)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Printf("✅ Success (version: %s)\n", apiVersion.Version)
	}
}
