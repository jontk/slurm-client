// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jontk/slurm-client/internal/factory"
	"github.com/jontk/slurm-client/internal/interfaces"
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

	versions := []string{}
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

		// Benchmark wrapper implementation
		fmt.Printf("\n--- Wrapper Implementation ---\n")
		wrapperResults := benchmarkImplementation(version, jwtToken, false)
		
		// Benchmark adapter implementation
		fmt.Printf("\n--- Adapter Implementation ---\n")
		adapterResults := benchmarkImplementation(version, jwtToken, true)

		// Compare results
		fmt.Printf("\n--- Performance Comparison ---\n")
		compareResults(wrapperResults, adapterResults)
	}
}

func benchmarkImplementation(version, jwtToken string, useAdapters bool) []BenchmarkResult {
	results := []BenchmarkResult{}
	implName := "Wrapper"
	if useAdapters {
		implName = "Adapter"
	}

	// Benchmark initialization
	initStart := time.Now()
	client, err := createClient(version, jwtToken, useAdapters)
	initDuration := time.Since(initStart)
	
	results = append(results, BenchmarkResult{
		Implementation: implName,
		Version:        version,
		Operation:      "Initialization",
		Duration:       initDuration,
		Error:          err,
	})

	if err != nil {
		fmt.Printf("Failed to create %s client: %v\n", implName, err)
		return results
	}

	ctx := context.Background()

	// Benchmark operations
	operations := []struct {
		name string
		fn   func() error
	}{
		{"List Jobs", func() error {
			_, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{Limit: 100})
			return err
		}},
		{"List Nodes", func() error {
			_, err := client.Nodes().List(ctx, &interfaces.ListNodesOptions{Limit: 100})
			return err
		}},
		{"List Partitions", func() error {
			_, err := client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{})
			return err
		}},
		{"List Accounts", func() error {
			_, err := client.Accounts().List(ctx, &interfaces.ListAccountsOptions{})
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

		for i := 0; i < iterations; i++ {
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

		results = append(results, BenchmarkResult{
			Implementation: implName,
			Version:        version,
			Operation:      op.name,
			Duration:       avgDuration,
			Error:          lastErr,
		})

		if successCount > 0 {
			fmt.Printf("%s: %v (avg of %d successful runs)\n", op.name, avgDuration, successCount)
		} else {
			fmt.Printf("%s: FAILED - %v\n", op.name, lastErr)
		}
	}

	return results
}

func createClient(version, jwtToken string, useAdapters bool) (interfaces.SlurmClient, error) {
	// Create configuration
	cfg := config.NewDefault()
	cfg.BaseURL = "http://rocky9.ar.jontk.com:6820"
	cfg.Debug = false

	// Create JWT authentication provider
	authProvider := auth.NewTokenAuth(jwtToken)

	// Create factory
	clientFactory, err := factory.NewClientFactory(
		factory.WithConfig(cfg),
		factory.WithAuth(authProvider),
		factory.WithBaseURL(cfg.BaseURL),
		factory.WithUseAdapters(useAdapters),
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

func compareResults(wrapperResults, adapterResults []BenchmarkResult) {
	// Create a map for easy comparison
	wrapperMap := make(map[string]BenchmarkResult)
	for _, r := range wrapperResults {
		wrapperMap[r.Operation] = r
	}

	fmt.Printf("\n%-20s %15s %15s %10s %s\n", "Operation", "Wrapper", "Adapter", "Diff", "Winner")
	fmt.Printf("%-20s %15s %15s %10s %s\n", "---------", "-------", "-------", "----", "------")

	for _, adapterResult := range adapterResults {
		if wrapperResult, ok := wrapperMap[adapterResult.Operation]; ok {
			var diff time.Duration
			var pctDiff float64
			var winner string

			// Both succeeded
			if wrapperResult.Error == nil && adapterResult.Error == nil {
				diff = adapterResult.Duration - wrapperResult.Duration
				if wrapperResult.Duration > 0 {
					pctDiff = float64(diff) / float64(wrapperResult.Duration) * 100
				}
				
				if adapterResult.Duration < wrapperResult.Duration {
					winner = "Adapter"
				} else if wrapperResult.Duration < adapterResult.Duration {
					winner = "Wrapper"
				} else {
					winner = "Tie"
				}

				fmt.Printf("%-20s %15s %15s %9.1f%% %s\n",
					adapterResult.Operation,
					wrapperResult.Duration,
					adapterResult.Duration,
					pctDiff,
					winner,
				)
			} else {
				// Handle errors
				wrapperStatus := "OK"
				adapterStatus := "OK"
				if wrapperResult.Error != nil {
					wrapperStatus = "FAILED"
				}
				if adapterResult.Error != nil {
					adapterStatus = "FAILED"
				}
				
				fmt.Printf("%-20s %15s %15s %10s %s\n",
					adapterResult.Operation,
					wrapperStatus,
					adapterStatus,
					"N/A",
					"-",
				)
			}
		}
	}

	// Overall summary
	fmt.Printf("\n--- Summary ---\n")
	wrapperSuccesses := 0
	adapterSuccesses := 0
	var totalWrapperTime, totalAdapterTime time.Duration

	for _, r := range wrapperResults {
		if r.Error == nil && r.Operation != "Initialization" {
			wrapperSuccesses++
			totalWrapperTime += r.Duration
		}
	}

	for _, r := range adapterResults {
		if r.Error == nil && r.Operation != "Initialization" {
			adapterSuccesses++
			totalAdapterTime += r.Duration
		}
	}

	fmt.Printf("Wrapper: %d/%d operations succeeded, Total time: %v\n", 
		wrapperSuccesses, len(wrapperResults)-1, totalWrapperTime)
	fmt.Printf("Adapter: %d/%d operations succeeded, Total time: %v\n", 
		adapterSuccesses, len(adapterResults)-1, totalAdapterTime)

	if totalWrapperTime > 0 && totalAdapterTime > 0 {
		speedup := float64(totalWrapperTime) / float64(totalAdapterTime)
		if speedup > 1 {
			fmt.Printf("\nAdapter is %.2fx faster overall\n", speedup)
		} else {
			fmt.Printf("\nWrapper is %.2fx faster overall\n", 1/speedup)
		}
	}
}