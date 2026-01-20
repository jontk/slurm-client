// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"

	"github.com/jontk/slurm-client/internal/adapters/common"
	v0_0_43 "github.com/jontk/slurm-client/internal/adapters/v0_0_43"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
	"github.com/jontk/slurm-client/internal/common/types"
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

	// Example: Coordinated operations across managers
	coordinatedManagerExample(adapter)

	fmt.Println("\n✅ Integration example completed!")
}

func demonstrateManagerUsage(adapter common.VersionAdapter) {
	_ = context.Background()

	// QoS Management
	fmt.Println("🎯 QoS Management:")
	qosManager := adapter.GetQoSManager()
	fmt.Printf("   - Manager type: %T\n", qosManager)

	// Example QoS operations (would work with real server)
	// qosList, err := qosManager.List(ctx, &types.QoSListOptions{
	//     Names: []string{"normal", "high"},
	// })
	// if err != nil {
	//     log.Printf("Failed to list QoS: %v", err)
	// }

	// Job Management
	fmt.Println("\n💼 Job Management:")
	jobManager := adapter.GetJobManager()
	fmt.Printf("   - Manager type: %T\n", jobManager)

	// Node Management
	fmt.Println("\n🖥️ Node Management:")
	nodeManager := adapter.GetNodeManager()
	fmt.Printf("   - Manager type: %T\n", nodeManager)

	// Partition Management
	fmt.Println("\n🔧 Partition Management:")
	partitionManager := adapter.GetPartitionManager()
	fmt.Printf("   - Manager type: %T\n", partitionManager)

	// Account Management
	fmt.Println("\n👥 Account Management:")
	accountManager := adapter.GetAccountManager()
	fmt.Printf("   - Manager type: %T\n", accountManager)

	// Association Management
	fmt.Println("\n🔗 Association Management:")
	associationManager := adapter.GetAssociationManager()
	fmt.Printf("   - Manager type: %T\n", associationManager)

	// Reservation Management
	fmt.Println("\n📅 Reservation Management:")
	reservationManager := adapter.GetReservationManager()
	fmt.Printf("   - Manager type: %T\n", reservationManager)
}

func coordinatedManagerExample(adapter common.VersionAdapter) {
	ctx := context.Background()

	fmt.Println("\n🎭 Coordinated Manager Operations:")

	// Get multiple managers
	qosManager := adapter.GetQoSManager()
	jobManager := adapter.GetJobManager()
	accountManager := adapter.GetAccountManager()

	// Example: Create a job with specific QoS and Account

	// 1. Verify QoS exists (mock example)
	// In a real scenario, this would fetch from the server
	qos := &types.QoS{
		Name:        "high_priority",
		Description: "High priority QoS",
		Priority:    1000,
	}
	_ = qos // Simulate using the QoS

	// Mock error handling for demonstration
	var err error
	// In a real implementation, you would check:
	// if err != nil {
	//     log.Printf("QoS not found: %v", err)
	//     return
	// }

	// 2. Verify account exists (mock example)
	account := &types.Account{
		Name:         "research_group",
		Description:  "Research Group Account",
		Organization: "University",
	}
	_ = account // Simulate using the account

	// 3. Demonstrate job submission (would work with real server)
	// jobResponse, err := jobManager.Submit(ctx, &types.JobCreate{
	//     Name:        "coordinated_job",
	//     Command:     "echo 'Hello from coordinated job'",
	//     Account:     account.Name,
	//     QoS:         qos.Name,
	//     TimeLimit:   3600, // 1 hour
	//     Partition:   "compute",
	// })
	// if err != nil {
	//     log.Printf("Job submission failed: %v", err)
	//     return
	// }
	// fmt.Printf("✅ Job submitted successfully: ID %d\n", jobResponse.JobID)

	// For this example, we'll just show the coordination concept
	fmt.Println("   - Would verify QoS 'high_priority' exists")
	fmt.Println("   - Would verify account 'research_group' exists")
	fmt.Println("   - Would submit job with verified QoS and Account")

	// Use the managers to avoid unused variable errors
	_ = qosManager
	_ = jobManager
	_ = accountManager
	_ = ctx
	_ = err

	fmt.Println("\n📝 Key Concepts:")
	fmt.Println("   - This example shows how managers work together")
	fmt.Println("   - Each manager handles its domain while sharing the same API client")
	fmt.Println("   - All operations use consistent error handling and validation")
	fmt.Println("   - In production, you would configure the client with real server details")
}
