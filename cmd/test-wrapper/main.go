// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

func main() {
	// Check if JWT token is provided as environment variable
	jwtToken := os.Getenv("SLURM_JWT")
	if jwtToken == "" {
		fmt.Println("Error: SLURM_JWT environment variable not set")
		fmt.Println("Usage: SLURM_JWT='your-jwt-token' go run test-wrapper.go")
		os.Exit(1)
	}

	// Configuration
	cfg := &config.Config{
		BaseURL: "http://rocky9.ar.jontk.com:6820/slurm",
		Debug:   true,
	}

	// Create token auth provider for JWT
	authProvider := auth.NewTokenAuth(jwtToken)

	// Create client with specific version for v0.0.40
	ctx := context.Background()
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.40",
		slurm.WithConfig(cfg),
		slurm.WithAuth(authProvider),
	)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		os.Exit(1)
	}

	// Test 1: Ping
	fmt.Println("\n=== Test 1: Ping ===")
	infoMgr := client.Info()
	err = infoMgr.Ping(ctx)
	if err != nil {
		fmt.Printf("Ping failed: %v\n", err)
	} else {
		fmt.Printf("Ping successful!\n")
	}

	// Test 2: Get API version
	fmt.Println("\n=== Test 2: Get API Version ===")
	fmt.Printf("Client API Version: %s\n", client.Version())
	
	// Test 2b: Get cluster info
	fmt.Println("\n=== Test 2b: Get Cluster Info ===")
	info, err := infoMgr.Get(ctx)
	if err != nil {
		fmt.Printf("Get cluster info failed: %v\n", err)
	} else {
		fmt.Printf("Cluster Info: Name=%s, Version=%s, API=%s\n", info.ClusterName, info.Version, info.APIVersion)
	}

	// Test 3: List jobs
	fmt.Println("\n=== Test 3: List Jobs ===")
	jobMgr := client.Jobs()
	jobs, err := jobMgr.List(ctx, nil)
	if err != nil {
		fmt.Printf("List jobs failed: %v\n", err)
	} else {
		fmt.Printf("Found %d jobs\n", len(jobs.Jobs))
		for i, job := range jobs.Jobs {
			if i < 5 { // Show first 5 jobs
				fmt.Printf("  Job: ID=%d, Name=%s, State=%s\n", job.ID, job.Name, job.State)
			}
		}
		if len(jobs.Jobs) > 5 {
			fmt.Printf("  ... and %d more jobs\n", len(jobs.Jobs)-5)
		}
	}

	// Test 4: List nodes
	fmt.Println("\n=== Test 4: List Nodes ===")
	nodeMgr := client.Nodes()
	nodes, err := nodeMgr.List(ctx, nil)
	if err != nil {
		fmt.Printf("List nodes failed: %v\n", err)
	} else {
		fmt.Printf("Found %d nodes\n", len(nodes.Nodes))
		for i, node := range nodes.Nodes {
			if i < 5 { // Show first 5 nodes
				fmt.Printf("  Node: Name=%s, State=%s, CPUs=%d\n", node.Name, node.State, node.CPUs)
			}
		}
		if len(nodes.Nodes) > 5 {
			fmt.Printf("  ... and %d more nodes\n", len(nodes.Nodes)-5)
		}
	}

	// Test 5: List partitions
	fmt.Println("\n=== Test 5: List Partitions ===")
	partMgr := client.Partitions()
	partitions, err := partMgr.List(ctx, nil)
	if err != nil {
		fmt.Printf("List partitions failed: %v\n", err)
	} else {
		fmt.Printf("Found %d partitions\n", len(partitions.Partitions))
		for _, partition := range partitions.Partitions {
			fmt.Printf("  Partition: Name=%s, State=%s, Nodes=%s\n", 
				partition.Name, partition.State, partition.Nodes)
		}
	}

	// Test 6: List accounts
	fmt.Println("\n=== Test 6: List Accounts ===")
	acctMgr := client.Accounts()
	accounts, err := acctMgr.List(ctx, nil)
	if err != nil {
		fmt.Printf("List accounts failed: %v\n", err)
	} else {
		fmt.Printf("Found %d accounts\n", len(accounts.Accounts))
		for i, account := range accounts.Accounts {
			if i < 5 { // Show first 5 accounts
				fmt.Printf("  Account: Name=%s, Description=%s\n", 
					account.Name, account.Description)
			}
		}
		if len(accounts.Accounts) > 5 {
			fmt.Printf("  ... and %d more accounts\n", len(accounts.Accounts)-5)
		}
	}

	// Test 7: List users
	fmt.Println("\n=== Test 7: List Users ===")
	userMgr := client.Users()
	users, err := userMgr.List(ctx, nil)
	if err != nil {
		fmt.Printf("List users failed: %v\n", err)
	} else {
		fmt.Printf("Found %d users\n", len(users.Users))
		for i, user := range users.Users {
			if i < 5 { // Show first 5 users
				fmt.Printf("  User: Name=%s, DefaultAccount=%s\n", 
					user.Name, user.DefaultAccount)
			}
		}
		if len(users.Users) > 5 {
			fmt.Printf("  ... and %d more users\n", len(users.Users)-5)
		}
	}

	fmt.Println("\n=== All Tests Completed ===")
}
