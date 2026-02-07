// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/pkg/auth"
)

// Example: Quality of Service (QoS) management (v0.0.43+ only)
func main() {
	ctx := context.Background()
	baseURL := "https://cluster.example.com:6820"
	jwtToken := "your-jwt-token"

	// Example 1: Check QoS support
	fmt.Println("=== QoS Support Check ===")
	checkQoSSupport(ctx, baseURL, jwtToken)

	// Example 2: List QoS configurations
	fmt.Println("\n=== List QoS Configurations ===")
	listQoSConfigurations(ctx, baseURL, jwtToken)

	// Example 3: Create QoS levels
	fmt.Println("\n=== Create QoS Levels ===")
	createQoSLevels(ctx, baseURL, jwtToken)

	// Example 4: Update QoS
	fmt.Println("\n=== Update QoS ===")
	updateQoS(ctx, baseURL, jwtToken)

	// Example 5: QoS hierarchy and preemption
	fmt.Println("\n=== QoS Hierarchy and Preemption ===")
	demonstrateQoSHierarchy(ctx, baseURL, jwtToken)

	// Example 6: Resource limits and fair share
	fmt.Println("\n=== Resource Limits and Fair Share ===")
	demonstrateResourceLimits(ctx, baseURL, jwtToken)
}

// checkQoSSupport checks if the cluster supports QoS management
func checkQoSSupport(ctx context.Context, baseURL, jwtToken string) {
	// Try different versions
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}

	for _, version := range versions {
		client, err := slurm.NewClient(ctx,
			slurm.WithBaseURL(baseURL),
			slurm.WithAuth(auth.NewTokenAuth(jwtToken)),
			slurm.WithVersion(version),
		)
		if err != nil {
			log.Printf("Failed to create %s client: %v", version, err)
			continue
		}
		defer client.Close()

		// Check if QoS is supported
		if client.QoS() == nil {
			fmt.Printf("%s: QoS NOT supported\n", version)
		} else {
			fmt.Printf("%s: QoS supported ✓\n", version)
		}
	}

	fmt.Println("\nNote: QoS management requires SLURM REST API v0.0.43 or later")
}

// listQoSConfigurations demonstrates listing QoS configurations
func listQoSConfigurations(ctx context.Context, baseURL, jwtToken string) {
	// Create v0.0.43 client
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL(baseURL),
		slurm.WithAuth(auth.NewTokenAuth(jwtToken)),
		slurm.WithVersion("v0.0.43"),
	)
	if err != nil {
		log.Printf("Failed to create v0.0.43 client: %v", err)
		return
	}
	defer client.Close()

	if client.QoS() == nil {
		fmt.Println("QoS not supported")
		return
	}

	// List all QoS
	fmt.Println("Listing all QoS configurations:")
	qosList, err := client.QoS().List(ctx, nil)
	if err != nil {
		log.Printf("Failed to list QoS: %v", err)
		return
	}

	if len(qosList.QoS) == 0 {
		fmt.Println("No QoS configurations found")
		return
	}

	// Display QoS configurations
	for _, qos := range qosList.QoS {
		fmt.Printf("\nQoS: %v\n", qos.Name)
		fmt.Printf("  Description: %v\n", qos.Description)
		fmt.Printf("  Priority: %v\n", qos.Priority)
		if qos.Preempt != nil {
			fmt.Printf("  Preempt Mode: %v\n", qos.Preempt.Mode)
		}
		fmt.Printf("  Usage Factor: %v\n", qos.UsageFactor)

		if len(qos.Flags) > 0 {
			fmt.Printf("  Flags: %v\n", qos.Flags)
		}
		if qos.Limits != nil {
			fmt.Printf("  Has custom limits configured\n")
		}
	}

	// List QoS for specific accounts
	fmt.Println("\nListing QoS for specific accounts:")
	accountQoS, err := client.QoS().List(ctx, &slurm.ListQoSOptions{
		Accounts: []string{"research", "engineering"},
	})
	if err != nil {
		log.Printf("Failed to list account QoS: %v", err)
		return
	}

	fmt.Printf("Found %d QoS configurations for specified accounts\n",
		len(accountQoS.QoS))
}

// createQoSLevels demonstrates creating different QoS levels
func createQoSLevels(ctx context.Context, baseURL, jwtToken string) {
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL(baseURL),
		slurm.WithAuth(auth.NewTokenAuth(jwtToken)),
		slurm.WithVersion("v0.0.43"),
	)
	if err != nil {
		log.Printf("Failed to create v0.0.43 client: %v", err)
		return
	}
	defer client.Close()

	if client.QoS() == nil {
		fmt.Println("QoS not supported")
		return
	}

	// Example 1: Create high-priority QoS
	fmt.Println("Creating high-priority QoS:")

	highPriorityQoS := &slurm.QoSCreate{
		Name:        "high-priority",
		Description: "High priority for critical jobs",
		Priority:    10000,
		PreemptMode: []string{"REQUEUE"},
		GraceTime:   300, // 5 minutes
		UsageFactor: 2.0, // Double charge for priority
		Flags:       []string{"DENY_LIMIT", "REQUIRE_ASSOC"},
		// Note: Detailed limit configuration would require setting Limits struct
		// MaxJobs, MaxCPUs, MaxWallTime, MaxNodes are configured via the Limits field
		// AllowedAccounts would be managed separately via account-QoS associations
	}

	resp, err := client.QoS().Create(ctx, highPriorityQoS)
	if err != nil {
		log.Printf("Failed to create high-priority QoS: %v", err)
	} else {
		fmt.Printf("Created QoS: %s\n", resp.QoSName)
	}

	// Example 2: Create normal QoS
	fmt.Println("\nCreating normal QoS:")

	normalQoS := &slurm.QoSCreate{
		Name:        "normal",
		Description: "Standard QoS for regular jobs",
		Priority:    1000,
		PreemptMode: []string{"SUSPEND"},
		GraceTime:   600, // 10 minutes
		UsageFactor: 1.0,
		Flags:       []string{"DENY_LIMIT"},
		// Note: Detailed limit configuration would require setting Limits struct
	}

	resp2, err := client.QoS().Create(ctx, normalQoS)
	if err != nil {
		log.Printf("Failed to create normal QoS: %v", err)
	} else {
		fmt.Printf("Created QoS: %s\n", resp2.QoSName)
	}

	// Example 3: Create low-priority/scavenger QoS
	fmt.Println("\nCreating scavenger QoS:")

	scavengerQoS := &slurm.QoSCreate{
		Name:        "scavenger",
		Description: "Low priority for opportunistic jobs",
		Priority:    100,
		PreemptMode: []string{"CANCEL"},
		UsageFactor: 0.1, // 10% charge - incentivize usage
		Flags:       []string{"NO_RESERVE"},
		// Note: Detailed limit configuration would require setting Limits struct
		// Account restrictions would be managed separately
	}

	resp3, err := client.QoS().Create(ctx, scavengerQoS)
	if err != nil {
		log.Printf("Failed to create scavenger QoS: %v", err)
	} else {
		fmt.Printf("Created QoS: %s\n", resp3.QoSName)
	}
}

// updateQoS demonstrates updating QoS configurations
func updateQoS(ctx context.Context, baseURL, jwtToken string) {
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL(baseURL),
		slurm.WithAuth(auth.NewTokenAuth(jwtToken)),
		slurm.WithVersion("v0.0.43"),
	)
	if err != nil {
		log.Printf("Failed to create v0.0.43 client: %v", err)
		return
	}
	defer client.Close()

	if client.QoS() == nil {
		fmt.Println("QoS not supported")
		return
	}

	qosName := "normal"

	// Get current QoS
	fmt.Printf("Getting QoS %s:\n", qosName)
	current, err := client.QoS().Get(ctx, qosName)
	if err != nil {
		log.Printf("Failed to get QoS: %v", err)
		return
	}

	fmt.Printf("Current settings:\n")
	fmt.Printf("  Priority: %v\n", current.Priority)
	fmt.Printf("  Usage Factor: %v\n", current.UsageFactor)
	if current.Limits != nil && current.Limits.Max != nil && current.Limits.Max.WallClock != nil {
		fmt.Printf("  Wall Clock Max: %v\n", current.Limits.Max.WallClock)
	}

	// Update QoS - update description
	newDescription := "Updated normal QoS with increased limits"
	update := &slurm.QoSUpdate{
		Description: &newDescription,
		// Note: Limit modifications would require setting via Limits struct
	}

	fmt.Println("\nUpdating QoS limits...")
	err = client.QoS().Update(ctx, qosName, update)
	if err != nil {
		log.Printf("Failed to update QoS: %v", err)
		return
	}

	fmt.Println("QoS updated successfully")

	// Verify update
	updated, err := client.QoS().Get(ctx, qosName)
	if err != nil {
		log.Printf("Failed to get updated QoS: %v", err)
		return
	}

	fmt.Printf("New settings:\n")
	fmt.Printf("  Priority: %v\n", updated.Priority)
	fmt.Printf("  Usage Factor: %v\n", updated.UsageFactor)
	if updated.Limits != nil && updated.Limits.Max != nil && updated.Limits.Max.WallClock != nil {
		fmt.Printf("  Wall Clock Max: %v\n", updated.Limits.Max.WallClock)
	}
}

// demonstrateQoSHierarchy shows QoS hierarchy and preemption
func demonstrateQoSHierarchy(ctx context.Context, baseURL, jwtToken string) {
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL(baseURL),
		slurm.WithAuth(auth.NewTokenAuth(jwtToken)),
		slurm.WithVersion("v0.0.43"),
	)
	if err != nil {
		log.Printf("Failed to create v0.0.43 client: %v", err)
		return
	}
	defer client.Close()

	if client.QoS() == nil {
		fmt.Println("QoS not supported")
		return
	}

	// Create a QoS hierarchy
	fmt.Println("Creating QoS hierarchy:")

	// 1. Executive QoS - highest priority, can preempt anything
	executive := &slurm.QoSCreate{
		Name:        "executive",
		Description: "Executive priority - preempts all",
		Priority:    100000,
		PreemptMode: []string{"REQUEUE"},
		GraceTime:   60, // 1 minute grace
		Flags: []string{
			// Note: PreemptExempt is not directly available as a string flag
			// instead, use PreemptList and PreemptExemptTime to control preemption behavior
		},
		// Note: User restrictions would be managed separately via account-QoS associations
	}

	_, err = client.QoS().Create(ctx, executive)
	if err != nil {
		log.Printf("Failed to create executive QoS: %v", err)
	} else {
		fmt.Println("  Created: executive (priority: 100000)")
	}

	// 2. Urgent QoS - high priority, can preempt normal and below
	urgent := &slurm.QoSCreate{
		Name:        "urgent",
		Description: "Urgent priority - preempts normal and low",
		Priority:    50000,
		PreemptMode: []string{"SUSPEND"},
		GraceTime:   300, // 5 minutes
		Flags:       []string{"DENY_LIMIT"},
		// Note: Account restrictions would be managed separately via account-QoS associations
	}

	_, err = client.QoS().Create(ctx, urgent)
	if err != nil {
		log.Printf("Failed to create urgent QoS: %v", err)
	} else {
		fmt.Println("  Created: urgent (priority: 50000)")
	}

	// 3. Interactive QoS - for interactive/debug jobs
	interactive := &slurm.QoSCreate{
		Name:        "interactive",
		Description: "Interactive jobs - quick turnaround",
		Priority:    5000,
		PreemptMode: []string{"DISABLED"},
		Flags:       []string{"NO_RESERVE"}, // Don't make reservations
		// Note: Job count and CPU limits would be configured via Limits struct
	}

	_, err = client.QoS().Create(ctx, interactive)
	if err != nil {
		log.Printf("Failed to create interactive QoS: %v", err)
	} else {
		fmt.Println("  Created: interactive (priority: 5000)")
	}

	fmt.Println("\nQoS Preemption Chain:")
	fmt.Println("  executive -> can preempt: urgent, normal, interactive, scavenger")
	fmt.Println("  urgent -> can preempt: normal, scavenger")
	fmt.Println("  normal -> can preempt: scavenger")
	fmt.Println("  interactive -> no preemption")
	fmt.Println("  scavenger -> can be preempted by all")
}

// demonstrateResourceLimits shows resource limits and fair share
func demonstrateResourceLimits(ctx context.Context, baseURL, jwtToken string) {
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL(baseURL),
		slurm.WithAuth(auth.NewTokenAuth(jwtToken)),
		slurm.WithVersion("v0.0.43"),
	)
	if err != nil {
		log.Printf("Failed to create v0.0.43 client: %v", err)
		return
	}
	defer client.Close()

	if client.QoS() == nil {
		fmt.Println("QoS not supported")
		return
	}

	// Example 1: Create QoS with strict resource limits
	fmt.Println("Creating QoS with strict resource limits:")

	limitedQoS := &slurm.QoSCreate{
		Name:            "limited-resources",
		Description:     "Strict resource limits for fair sharing",
		Priority:        1000,
		UsageFactor:     1.5, // 50% surcharge
		UsageThreshold:  0.8, // Start limiting at 80% usage
		Flags: []string{
			"DENY_LIMIT",               // Deny when limits reached
			"ENFORCE_USAGE_THRESHOLD",  // Enforce the usage threshold
			"NO_DECAY",                 // Don't decay priority
		},
		// Note: Detailed job and CPU limits would be configured via Limits struct
	}

	_, err = client.QoS().Create(ctx, limitedQoS)
	if err != nil {
		log.Printf("Failed to create limited QoS: %v", err)
	} else {
		fmt.Println("Created QoS with strict limits")
	}

	// Example 2: Create QoS for GPU resources
	fmt.Println("\nCreating GPU-specific QoS:")

	gpuQoS := &slurm.QoSCreate{
		Name:        "gpu-jobs",
		Description: "QoS for GPU-accelerated jobs",
		Priority:    2000,
		UsageFactor: 3.0, // 3x charge for GPU resources
		Flags:       []string{"REQUIRE_ASSOC", "PARTITION_QOS"},
		// Note: Actual GPU limits would be set via GRES in job submission
		// Job and CPU limits would be configured via Limits struct
		// Account restrictions would be managed separately via account-QoS associations
	}

	_, err = client.QoS().Create(ctx, gpuQoS)
	if err != nil {
		log.Printf("Failed to create GPU QoS: %v", err)
	} else {
		fmt.Println("Created GPU-specific QoS")
	}

	// Example 3: Fair share configuration
	fmt.Println("\nDemonstrating fair share with QoS:")

	fairShareQoS := &slurm.QoSCreate{
		Name:            "fair-share",
		Description:     "Fair share QoS with usage tracking",
		Priority:        1000,
		UsageFactor:     1.0,
		UsageThreshold:  0.5, // Start reducing priority at 50% share
		Flags: []string{
			"NO_DECAY",            // Maintain priority decay
			"USAGE_FACTOR_SAFE",   // Protected from usage factor changes
			"PARTITION_TIME_LIMIT", // Use partition time limits
		},
		// Note: Job limits would be configured via Limits struct
	}

	_, err = client.QoS().Create(ctx, fairShareQoS)
	if err != nil {
		log.Printf("Failed to create fair share QoS: %v", err)
	} else {
		fmt.Println("Created fair share QoS")
	}

	fmt.Println("\nResource Limit Summary:")
	fmt.Println("  limited-resources: Strict limits with usage threshold")
	fmt.Println("  gpu-jobs: Higher cost (3x) for GPU resource usage")
	fmt.Println("  fair-share: Balanced sharing with usage tracking")

	// Show how to check current usage
	fmt.Println("\nTo check current QoS usage:")
	fmt.Println("  - List jobs by QoS to see active usage")
	fmt.Println("  - Monitor MaxSubmitJobs vs current queued jobs")
	fmt.Println("  - Track UsageThreshold impacts on job priority")
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
