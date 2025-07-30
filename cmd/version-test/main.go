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

func main() {
	serverURL := "http://rocky9.ar.jontk.com:6820"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI2NTM4Mjk5NzYsImlhdCI6MTc1MzgyOTk3Niwic3VuIjoicm9vdCJ9.-z8Cq_wHuOxNJ7KHHTboX3l9r6JBtSD1RxQUgQR9owE"
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}

	ctx := context.Background()

	fmt.Println("Testing SLURM API Versions")
	fmt.Println("Server:", serverURL)
	fmt.Println()

	for _, version := range versions {
		fmt.Printf("=== Testing API Version %s ===\n", version)

		// Create client for this version
		client, err := slurm.NewClientWithVersion(ctx, version,
			slurm.WithBaseURL(serverURL),
			slurm.WithAuth(auth.NewTokenAuth(token)),
			slurm.WithConfig(&config.Config{
				Timeout:            30 * time.Second,
				MaxRetries:         3,
				Debug:              false,
				InsecureSkipVerify: true,
			}),
		)
		if err != nil {
			log.Printf("Failed to create client for %s: %v", version, err)
			continue
		}
		defer client.Close()

		// Test 1: Basic connectivity
		fmt.Print("  Ping: ")
		if err := client.Info().Ping(ctx); err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Println("✅ Success")
		}

		// Test 2: Get version info
		fmt.Print("  Version: ")
		if info, err := client.Info().Version(ctx); err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ %s\n", info.Version)
		}

		// Test 3: List jobs
		fmt.Print("  List Jobs: ")
		if jobs, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{Limit: 5}); err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ Found %d jobs\n", len(jobs.Jobs))
		}

		// Test 4: Job submission
		fmt.Print("  Submit Job: ")
		submission := &interfaces.JobSubmission{
			Name:       fmt.Sprintf("test-%s-%d", version, time.Now().Unix()),
			Script:     "#!/bin/bash\necho 'Testing " + version + "'\nhostname\ndate\nsleep 5",
			Partition:  "debug",
			Nodes:      1,
			CPUs:       1,
			TimeLimit:  5,
			WorkingDir: "/tmp",
		}

		if resp, err := client.Jobs().Submit(ctx, submission); err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ Job ID=%s", resp.JobID)
			// Try to cancel
			if err := client.Jobs().Cancel(ctx, resp.JobID); err != nil {
				fmt.Printf(" (cancel failed: %v)\n", err)
			} else {
				fmt.Println(" (cancelled)")
			}
		}

		// Test 5: List nodes
		fmt.Print("  List Nodes: ")
		if nodes, err := client.Nodes().List(ctx, &interfaces.ListNodesOptions{Limit: 5}); err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ Found %d nodes\n", len(nodes.Nodes))
		}

		// Test 6: List partitions
		fmt.Print("  List Partitions: ")
		if partitions, err := client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{Limit: 5}); err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ Found %d partitions\n", len(partitions.Partitions))
		}

		fmt.Println()
	}

	fmt.Println("=== Summary ===")
	fmt.Println("Tested all supported API versions against the SLURM server")
}