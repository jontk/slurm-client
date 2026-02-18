// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	slurm "github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/pkg/config"
)

// userTokenAuth implements authentication with both username and token headers
type userTokenAuth struct {
	username string
	token    string
}

func (u *userTokenAuth) Authenticate(_ context.Context, req *http.Request) error {
	req.Header.Set("X-SLURM-USER-NAME", u.username)
	req.Header.Set("X-SLURM-USER-TOKEN", u.token)
	return nil
}

func (u *userTokenAuth) Type() string {
	return "user-token"
}

func main() {
	serverURL := "http://localhost:6820"
	//nolint:gosec // G101 - Test token for development/testing purposes only
	token := "your-jwt-token-here"
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}

	ctx := context.Background()

	fmt.Println("Testing SLURM API Versions")
	fmt.Println("Server:", serverURL)
	fmt.Println()

	for _, version := range versions {
		testVersion(ctx, version, serverURL, token)
	}

	fmt.Println("=== Summary ===")
	fmt.Println("Tested all supported API versions against the SLURM server")
}

func testVersion(ctx context.Context, version, serverURL, token string) {
	fmt.Printf("=== Testing API Version %s ===\n", version)

	username := "root" // Default username for testing

	// Create client for this version
	client, err := slurm.NewClientWithVersion(ctx, version,
		slurm.WithBaseURL(serverURL),
		slurm.WithAuth(&userTokenAuth{
			username: username,
			token:    token,
		}),
		slurm.WithConfig(&config.Config{
			Timeout:            30 * time.Second,
			MaxRetries:         3,
			Debug:              false,
			InsecureSkipVerify: true,
		}),
	)
	if err != nil {
		log.Printf("Failed to create client for %s: %v", version, err)
		return
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
	if jobs, err := client.Jobs().List(ctx, &slurm.ListJobsOptions{Limit: 5}); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Printf("✅ Found %d jobs\n", len(jobs.Jobs))
	}

	// Test 4: Job submission
	fmt.Print("  Submit Job: ")
	submission := &slurm.JobSubmission{
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
		jobIdStr := fmt.Sprintf("%d", resp.JobId)
		fmt.Printf("✅ Job ID=%s", jobIdStr)
		// Try to cancel
		if err := client.Jobs().Cancel(ctx, jobIdStr); err != nil {
			fmt.Printf(" (cancel failed: %v)\n", err)
		} else {
			fmt.Println(" (cancelled)")
		}
	}

	// Test 5: List nodes
	fmt.Print("  List Nodes: ")
	if nodes, err := client.Nodes().List(ctx, &slurm.ListNodesOptions{Limit: 5}); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Printf("✅ Found %d nodes\n", len(nodes.Nodes))
	}

	// Test 6: List partitions
	fmt.Print("  List Partitions: ")
	if partitions, err := client.Partitions().List(ctx, &slurm.ListPartitionsOptions{Limit: 5}); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Printf("✅ Found %d partitions\n", len(partitions.Partitions))
	}

	fmt.Println()
}
