// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	// "log"

	"github.com/jontk/slurm-client/internal/adapters/common"
	v0_0_43 "github.com/jontk/slurm-client/internal/adapters/v0_0_43"
	// "github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// Example demonstrating how to use the integrated adapter
func main() {
	fmt.Println("🚀 SLURM Client Adapter Integration Example")
	fmt.Println("==========================================")
	
	// Create an API client (this would normally be configured with real server details)
	client := &api.ClientWithResponses{
		// ClientInterface: httpClient, // Configure with real HTTP client
		// Server: "https://your-slurm-server.com",
	}
	
	// Create the version-specific adapter
	adapter := v0_0_43.NewAdapter(client)
	
	// The adapter implements the VersionAdapter interface
	var versionAdapter common.VersionAdapter = adapter
	
	fmt.Printf("📋 Using adapter for version: %s\n\n", versionAdapter.GetVersion())
	
	// Example: Working with different managers
	demonstrateManagerUsage(adapter)
	
	fmt.Println("✅ Integration example completed!")
}

func demonstrateManagerUsage(adapter common.VersionAdapter) {
	_ = context.Background()
	
	// QoS Management
	fmt.Println("🎯 QoS Management:")
	qosManager := adapter.GetQoSManager()
	fmt.Printf("   - Manager type: %T\n", qosManager)
	
	// Example QoS operations (would work with real server)
	/*
	qosList, err := qosManager.List(ctx, &types.QoSListOptions{
		Names: []string{"normal", "high"},
	})
	if err != nil {
		log.Printf("   - Error listing QoS: %v", err)
	} else {
		fmt.Printf("   - Found %d QoS entries\n", len(qosList.QoSList))
	}
	
	// Job Management
	fmt.Println("💼 Job Management:")
	jobManager := adapter.GetJobManager()
	fmt.Printf("   - Manager type: %T\n", jobManager)
	
	// Example job operations (would work with real server)
	/*
	jobList, err := jobManager.List(ctx, &types.JobListOptions{
		States: []types.JobState{types.JobStateRunning, types.JobStatePending},
		Limit:  10,
	})
	if err != nil {
		log.Printf("   - Error listing jobs: %v", err)
	} else {
		fmt.Printf("   - Found %d jobs\n", len(jobList.Jobs))
	}
	
	// Partition Management
	fmt.Println("🗂️  Partition Management:")
	partitionManager := adapter.GetPartitionManager()
	fmt.Printf("   - Manager type: %T\n", partitionManager)
	
	// Node Management
	fmt.Println("🖥️  Node Management:")
	nodeManager := adapter.GetNodeManager()
	fmt.Printf("   - Manager type: %T\n", nodeManager)
	
	// Account Management
	fmt.Println("👤 Account Management:")
	accountManager := adapter.GetAccountManager()
	fmt.Printf("   - Manager type: %T\n", accountManager)
	
	// User Management
	fmt.Println("👥 User Management:")
	userManager := adapter.GetUserManager()
	fmt.Printf("   - Manager type: %T\n", userManager)
	
	// Reservation Management
	fmt.Println("📅 Reservation Management:")
	reservationManager := adapter.GetReservationManager()
	fmt.Printf("   - Manager type: %T\n", reservationManager)
	
	// Association Management
	fmt.Println("🔗 Association Management:")
	associationManager := adapter.GetAssociationManager()
	fmt.Printf("   - Manager type: %T\n", associationManager)
	
	fmt.Println()
}

// Example of how to work with multiple managers in a coordinated way
func coordinatedManagerExample(adapter common.VersionAdapter) {
	_ = context.Background()
	
	fmt.Println("🎭 Coordinated Manager Operations:")
	
	// Get multiple managers
	_ = adapter.GetQoSManager()
	_ = adapter.GetJobManager()
	_ = adapter.GetAccountManager()
	
	// Example: Create a job with specific QoS and Account
	/*
	// 1. Verify QoS exists
	qos, err := qosManager.Get(ctx, "high_priority")
	if err != nil {
		log.Printf("QoS not found: %v", err)
		return
	}
	
	// 2. Verify account exists
	account, err := accountManager.Get(ctx, "research_group")
	if err != nil {
		log.Printf("Account not found: %v", err)
		return
	}
	
	// 3. Submit job with verified QoS and Account
	jobResponse, err := jobManager.Submit(ctx, &types.JobCreate{
		Name:        "coordinated_job",
		Command:     "echo 'Hello from coordinated job'",
		Account:     account.Name,
		QoS:         qos.Name,
		TimeLimit:   3600, // 1 hour
		Partition:   "compute",
	})
	if err != nil {
		log.Printf("Job submission failed: %v", err)
		return
	}
	
	fmt.Printf("✅ Job submitted successfully: ID %d\n", jobResponse.JobID)
	
	fmt.Println("   - This example shows how managers work together")
	fmt.Println("   - Each manager handles its domain while sharing the same API client")
	fmt.Println("   - All operations use consistent error handling and validation")
}
