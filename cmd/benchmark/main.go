// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/internal/factory"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

type BenchmarkResult struct {
	Implementation string
	Version        string
	Operation      string
	Duration       time.Duration
	MemoryUsed     int64
	Error          error
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <version|all>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Supported versions: v0.0.40, v0.0.41, v0.0.42, v0.0.43, all\n")
		os.Exit(1)
	}

	// Get JWT token from environment
	jwtToken := os.Getenv("SLURM_JWT")
	if jwtToken == "" {
		log.Fatal("SLURM_JWT environment variable is required")
	}

	var versions []string
	if os.Args[1] == "all" {
		versions = []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}
	} else {
		versions = []string{os.Args[1]}
	}

	// Run benchmarks for each version
	for _, version := range versions {
		fmt.Printf("\n========================================\n")
		fmt.Printf("Benchmarking API Version: %s\n", version)
		fmt.Printf("========================================\n")

		// Benchmark adapter implementation
		fmt.Printf("\n--- Adapter Implementation ---\n")
		benchmarkImplementation(version, jwtToken)
	}
}

func benchmarkImplementation(version, jwtToken string) {
	implName := "Adapter"

	// Benchmark initialization
	initStart := time.Now()
	client, err := createClient(version, jwtToken)
	initDuration := time.Since(initStart)

	if err != nil {
		fmt.Printf("Failed to create %s client: %v\n", implName, err)
		return
	}

	fmt.Printf("Initialization: %v\n", initDuration)

	ctx := context.Background()

	// Benchmark operations
	operations := []struct {
		name string
		fn   func() error
	}{
		{"List Jobs", func() error {
			_, err := client.Jobs().List(ctx, &types.ListJobsOptions{Limit: 100})
			return err
		}},
		{"List Nodes", func() error {
			_, err := client.Nodes().List(ctx, &types.ListNodesOptions{Limit: 100})
			return err
		}},
		{"List Partitions", func() error {
			_, err := client.Partitions().List(ctx, &types.ListPartitionsOptions{})
			return err
		}},
		{"List Accounts", func() error {
			_, err := client.Accounts().List(ctx, &types.ListAccountsOptions{})
			return err
		}},
		{"Ping", func() error {
			return client.Info().Ping(ctx)
		}},
	}

	// Run each operation multiple times and average
	iterations := 10
	for _, op := range operations {
		var totalDuration time.Duration
		var successCount int
		var lastErr error

		for range iterations {
			start := time.Now()
			err := op.fn()
			duration := time.Since(start)

			if err == nil {
				totalDuration += duration
				successCount++
			} else {
				lastErr = err
			}
		}

		avgDuration := time.Duration(0)
		if successCount > 0 {
			avgDuration = totalDuration / time.Duration(successCount)
		}

		if successCount > 0 {
			fmt.Printf("%s: %v (avg of %d successful runs)\n", op.name, avgDuration, successCount)
		} else {
			fmt.Printf("%s: FAILED - %v\n", op.name, lastErr)
		}
	}
}

func createClient(version, jwtToken string) (types.SlurmClient, error) {
	// Create configuration
	cfg := config.NewDefault()
	cfg.BaseURL = "http://localhost:6820"
	cfg.Debug = false

	// Create JWT authentication provider
	authProvider := auth.NewTokenAuth(jwtToken)

	// Create factory
	clientFactory, err := factory.NewClientFactory(
		factory.WithConfig(cfg),
		factory.WithAuth(authProvider),
		factory.WithBaseURL(cfg.BaseURL),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create factory: %w", err)
	}

	// Create client with specific version
	client, err := clientFactory.NewClientWithVersion(context.Background(), version)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}


