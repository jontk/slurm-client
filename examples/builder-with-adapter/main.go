// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"net/http"
	
	"github.com/jontk/slurm-client/internal/adapters/common"
	adapterV0043 "github.com/jontk/slurm-client/internal/adapters/v0_0_43"
	apiV0043 "github.com/jontk/slurm-client/internal/api/v0_0_43"
	"github.com/jontk/slurm-client/internal/common/builders"
)

func main() {
	// Get configuration from environment
	host := os.Getenv("SLURM_HOST")
	if host == "" {
		host = "localhost:6820"
	}

	token := os.Getenv("SLURM_TOKEN")
	if token == "" {
		log.Fatal("SLURM_TOKEN environment variable is required")
	}

	// Create the base HTTP client with authentication
	httpClient := &http.Client{
		Transport: &tokenTransport{
			Token: token,
			Base:  http.DefaultTransport,
		},
	}

	// Create version-specific client
	apiClient, err := apiV0043.NewClientWithResponses(
		fmt.Sprintf("http://%s", host),
		apiV0043.WithHTTPClient(httpClient),
	)
	if err != nil {
		log.Fatalf("Failed to create API client: %v", err)
	}

	// Create the adapter
	adapter := adapterV0043.NewAdapter(apiClient)
	ctx := context.Background()

	// Demonstrate various builder patterns
	fmt.Println("=== QoS Builder Pattern Examples ===")

	// Example 1: Simple QoS
	if err := createSimpleQoS(ctx, adapter); err != nil {
		log.Printf("Error creating simple QoS: %v", err)
	}

	// Example 2: High Priority QoS with Limits
	if err := createHighPriorityQoS(ctx, adapter); err != nil {
		log.Printf("Error creating high priority QoS: %v", err)
	}

	// Example 3: Batch Queue QoS
	if err := createBatchQoS(ctx, adapter); err != nil {
		log.Printf("Error creating batch QoS: %v", err)
	}

	// Example 4: Interactive QoS
	if err := createInteractiveQoS(ctx, adapter); err != nil {
		log.Printf("Error creating interactive QoS: %v", err)
	}

	// Example 5: Update using builder
	if err := updateQoSWithBuilder(ctx, adapter); err != nil {
		log.Printf("Error updating QoS: %v", err)
	}

	// Example 6: Clone and modify
	if err := cloneAndModifyQoS(ctx, adapter); err != nil {
		log.Printf("Error with clone example: %v", err)
	}
}

func createSimpleQoS(ctx context.Context, adapter common.VersionAdapter) error {
	fmt.Println("1. Creating Simple QoS")
	fmt.Println("---------------------")

	qos, err := builders.NewQoSBuilder("simple-qos").
		WithDescription("A simple QoS for basic jobs").
		WithPriority(50).
		Build()

	if err != nil {
		return fmt.Errorf("failed to build QoS: %w", err)
	}

	qosManager := adapter.GetQoSManager()
	resp, err := qosManager.Create(ctx, qos)
	if err != nil {
		return fmt.Errorf("failed to create QoS: %w", err)
	}

	fmt.Printf("✓ Created simple QoS: %s\n\n", resp.QoSName)
	return nil
}

func createHighPriorityQoS(ctx context.Context, adapter common.VersionAdapter) error {
	fmt.Println("2. Creating High Priority QoS with Limits")
	fmt.Println("----------------------------------------")

	qos, err := builders.NewQoSBuilder("high-priority-qos").
		AsHighPriority(). // Use preset for high priority
		WithDescription("High priority QoS for critical workloads").
		WithPreemptExemptTime(10). // 10 minutes exempt from preemption
		WithLimits().
			WithMaxCPUsPerUser(500).
			WithMaxJobsPerUser(20).
			WithMaxNodesPerUser(25).
			WithMaxWallTime(24 * time.Hour).
			WithMaxMemoryPerNode(256 * builders.GB).
			WithMaxMemoryPerCPU(8 * builders.GB).
			WithMinCPUsPerJob(4). // Require at least 4 CPUs
			Done().
		Build()

	if err != nil {
		return fmt.Errorf("failed to build QoS: %w", err)
	}

	qosManager := adapter.GetQoSManager()
	resp, err := qosManager.Create(ctx, qos)
	if err != nil {
		return fmt.Errorf("failed to create QoS: %w", err)
	}

	fmt.Printf("✓ Created high priority QoS: %s\n", resp.QoSName)
	fmt.Printf("  - Priority: %d\n", qos.Priority)
	fmt.Printf("  - Usage Factor: %.1f\n", qos.UsageFactor)
	fmt.Printf("  - Max CPUs/user: %d\n", *qos.Limits.MaxCPUsPerUser)
	fmt.Printf("  - Max wall time: %v\n\n", 24*time.Hour)
	return nil
}

func createBatchQoS(ctx context.Context, adapter common.VersionAdapter) error {
	fmt.Println("3. Creating Batch Queue QoS")
	fmt.Println("--------------------------")

	qos, err := builders.NewQoSBuilder("batch-qos").
		AsBatchQueue(). // Use preset for batch jobs
		WithDescription("QoS for long-running batch jobs").
		WithLimits().
			WithMaxCPUsPerUser(200).
			WithMaxJobsPerUser(100).
			WithMaxSubmitJobsPerUser(500). // Can queue many jobs
			WithMaxWallTime(7 * 24 * time.Hour). // 7 days max
			WithMaxMemoryPerNode(128 * builders.GB).
			Done().
		Build()

	if err != nil {
		return fmt.Errorf("failed to build QoS: %w", err)
	}

	qosManager := adapter.GetQoSManager()
	resp, err := qosManager.Create(ctx, qos)
	if err != nil {
		return fmt.Errorf("failed to create QoS: %w", err)
	}

	fmt.Printf("✓ Created batch QoS: %s\n", resp.QoSName)
	fmt.Printf("  - Priority: %d (low)\n", qos.Priority)
	fmt.Printf("  - Usage Factor: %.1f (discounted)\n", qos.UsageFactor)
	fmt.Printf("  - Max submit jobs: %d\n\n", *qos.Limits.MaxSubmitJobsPerUser)
	return nil
}

func createInteractiveQoS(ctx context.Context, adapter common.VersionAdapter) error {
	fmt.Println("4. Creating Interactive QoS")
	fmt.Println("--------------------------")

	qos, err := builders.NewQoSBuilder("interactive-qos").
		AsInteractive(). // Use preset for interactive jobs
		WithDescription("QoS for interactive development sessions").
		WithGraceTime(300). // 5 minute grace period
		WithLimits().
			WithMaxCPUsPerUser(32).
			WithMaxJobsPerUser(5). // Limited concurrent sessions
			WithMaxWallTime(8 * time.Hour). // 8 hour sessions
			WithMaxMemoryPerNode(64 * builders.GB).
			WithMaxNodesPerJob(1). // Single node only
			Done().
		Build()

	if err != nil {
		return fmt.Errorf("failed to build QoS: %w", err)
	}

	qosManager := adapter.GetQoSManager()
	resp, err := qosManager.Create(ctx, qos)
	if err != nil {
		return fmt.Errorf("failed to create QoS: %w", err)
	}

	fmt.Printf("✓ Created interactive QoS: %s\n", resp.QoSName)
	fmt.Printf("  - Priority: %d (medium-high)\n", qos.Priority)
	fmt.Printf("  - Preempt mode: %s\n", qos.PreemptMode[0])
	fmt.Printf("  - Max concurrent jobs: %d\n\n", *qos.Limits.MaxJobsPerUser)
	return nil
}

func updateQoSWithBuilder(ctx context.Context, adapter common.VersionAdapter) error {
	fmt.Println("5. Updating QoS with Builder")
	fmt.Println("---------------------------")

	// First create a QoS to update
	original, err := builders.NewQoSBuilder("update-example").
		WithDescription("Original description").
		WithPriority(100).
		Build()

	if err != nil {
		return fmt.Errorf("failed to build original QoS: %w", err)
	}

	qosManager := adapter.GetQoSManager()
	_, err = qosManager.Create(ctx, original)
	if err != nil {
		return fmt.Errorf("failed to create original QoS: %w", err)
	}

	// Now build an update
	update, err := builders.NewQoSBuilder("update-example").
		WithDescription("Updated via builder").
		WithPriority(200).
		WithUsageFactor(1.5).
		WithLimits().
			WithMaxCPUsPerUser(150).
			WithMaxJobsPerUser(15).
			Done().
		BuildForUpdate() // Note: BuildForUpdate() instead of Build()

	if err != nil {
		return fmt.Errorf("failed to build update: %w", err)
	}

	err = qosManager.Update(ctx, "update-example", update)
	if err != nil {
		return fmt.Errorf("failed to update QoS: %w", err)
	}

	fmt.Printf("✓ Updated QoS: update-example\n")
	fmt.Printf("  - New description: %s\n", *update.Description)
	fmt.Printf("  - New priority: %d\n\n", *update.Priority)
	return nil
}

func cloneAndModifyQoS(ctx context.Context, adapter common.VersionAdapter) error {
	fmt.Println("6. Clone and Modify Pattern")
	fmt.Println("--------------------------")

	// Create a base template
	baseTemplate := builders.NewQoSBuilder("template").
		WithDescription("Base GPU QoS template").
		WithPriority(300).
		WithFlags("gpu", "DenyOnLimit").
		WithLimits().
			WithMaxCPUsPerUser(64).
			WithMaxNodesPerUser(4).
			WithMaxMemoryPerNode(128 * builders.GB).
			Done()

	// Clone and create variants
	variants := []struct {
		name        string
		modifier    func(*builders.QoSBuilder) *builders.QoSBuilder
	}{
		{
			name: "gpu-small",
			modifier: func(b *builders.QoSBuilder) *builders.QoSBuilder {
				return b.WithLimits().
					WithMaxCPUsPerUser(16).
					WithMaxNodesPerUser(1).
					Done()
			},
		},
		{
			name: "gpu-large",
			modifier: func(b *builders.QoSBuilder) *builders.QoSBuilder {
				return b.WithPriority(400).
					WithLimits().
						WithMaxCPUsPerUser(128).
						WithMaxNodesPerUser(8).
						Done()
			},
		},
		{
			name: "gpu-debug",
			modifier: func(b *builders.QoSBuilder) *builders.QoSBuilder {
				return b.WithPriority(500).
					WithFlags("debug").
					WithLimits().
						WithMaxWallTime(1 * time.Hour).
						WithMaxJobsPerUser(1).
						Done()
			},
		},
	}

	qosManager := adapter.GetQoSManager()

	for _, variant := range variants {
		// Clone the base template
		cloned := baseTemplate.Clone()
		
		// Apply variant-specific modifications
		builder := variant.modifier(cloned)
		
		// Get the built QoS from the modified builder
		builtQoS, _ := builder.Build()
		
		// Update the name
		newBuilder := builders.NewQoSBuilder(variant.name).
			WithDescription(fmt.Sprintf("GPU QoS variant: %s", variant.name)).
			WithPriority(builtQoS.Priority). // Preserve modified priority
			WithFlags(builtQoS.Flags...)    // Preserve all flags

		// Build the QoS
		qos, err := newBuilder.Build()
		if err != nil {
			return fmt.Errorf("failed to build %s: %w", variant.name, err)
		}

		resp, err := qosManager.Create(ctx, qos)
		if err != nil {
			log.Printf("Failed to create %s: %v", variant.name, err)
			continue
		}

		fmt.Printf("✓ Created variant: %s\n", resp.QoSName)
	}

	fmt.Println("")
	return nil
}

// tokenTransport adds authentication token to requests
type tokenTransport struct {
	Token string
	Base  http.RoundTripper
}

func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-SLURM-USER-TOKEN", t.Token)
	return t.Base.RoundTrip(req)
}
