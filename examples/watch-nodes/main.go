// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

func main() {
	// Create configuration
	cfg := config.NewDefault()
	
	// Override with environment variables if needed
	if url := os.Getenv("SLURM_REST_URL"); url != "" {
		cfg.BaseURL = url
	}

	// Create authentication provider
	var authProvider auth.Provider
	if token := os.Getenv("SLURM_JWT"); token != "" {
		authProvider = auth.NewTokenAuth(token)
	} else if username := os.Getenv("SLURM_USERNAME"); username != "" {
		password := os.Getenv("SLURM_PASSWORD")
		authProvider = auth.NewBasicAuth(username, password)
	} else {
		authProvider = auth.NewNoAuth()
	}

	// Create client
	ctx := context.Background()
	client, err := slurm.NewClient(ctx,
		slurm.WithConfig(cfg),
		slurm.WithAuth(authProvider),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Get node manager
	nodeManager := client.Nodes()

	// Set up watch options
	watchOpts := &interfaces.WatchNodesOptions{
		// Watch specific node states (optional)
		// States: []string{"DOWN", "DRAINING"},
		
		// Watch specific nodes (optional)
		// NodeNames: []string{"node-001", "node-002"},
		
		// Watch nodes in specific partition (optional)
		// Partition: "gpu",
	}

	// Start watching for node events
	fmt.Println("Starting to watch for node events...")
	fmt.Println("Press Ctrl+C to stop")
	
	watchCtx, cancelWatch := context.WithCancel(ctx)
	defer cancelWatch()

	eventChan, err := nodeManager.Watch(watchCtx, watchOpts)
	if err != nil {
		log.Fatalf("Failed to start watching nodes: %v", err)
	}

	// Handle interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Process events
	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				fmt.Println("Event channel closed")
				return
			}

			// Handle different event types
			switch event.Type {
			case "node_new":
				fmt.Printf("[%s] New node detected: %s, State=%s\n",
					event.Timestamp.Format(time.RFC3339),
					event.NodeName,
					event.NewState)
				if event.Node != nil {
					fmt.Printf("  CPUs: %d, Memory: %d MB, Partitions: %v\n", 
						event.Node.CPUs, event.Node.Memory, event.Node.Partitions)
				}

			case "node_state_change":
				fmt.Printf("[%s] Node state changed: %s, %s -> %s\n",
					event.Timestamp.Format(time.RFC3339),
					event.NodeName,
					event.OldState,
					event.NewState)
				
				// Show additional details for state changes
				if event.Node != nil {
					if event.NewState == "DOWN" || event.NewState == "DRAINING" {
						fmt.Printf("  Reason: %s\n", event.Node.Reason)
					}
					fmt.Printf("  CPUs: %d, Memory: %d MB\n", 
						event.Node.CPUs, event.Node.Memory)
				}

			case "error":
				fmt.Printf("[%s] Error: %v\n",
					event.Timestamp.Format(time.RFC3339),
					event.Error)

			default:
				fmt.Printf("[%s] Unknown event type: %s\n",
					event.Timestamp.Format(time.RFC3339),
					event.Type)
			}

		case <-sigChan:
			fmt.Println("\nShutting down...")
			cancelWatch()
			// Give a moment for cleanup
			time.Sleep(100 * time.Millisecond)
			return
		}
	}
}
