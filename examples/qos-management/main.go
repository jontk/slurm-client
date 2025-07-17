package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

// Example: Quality of Service (QoS) management (v0.0.43+ only)
func main() {
	// Create configuration
	cfg := config.NewDefault()
	cfg.BaseURL = "https://cluster.example.com:6820"
	
	// Create authentication
	authProvider := auth.NewTokenAuth("your-jwt-token")

	ctx := context.Background()

	// Example 1: Check QoS support
	fmt.Println("=== QoS Support Check ===")
	checkQoSSupport(ctx, cfg, authProvider)

	// Example 2: List QoS configurations
	fmt.Println("\n=== List QoS Configurations ===")
	listQoSConfigurations(ctx, cfg, authProvider)

	// Example 3: Create QoS levels
	fmt.Println("\n=== Create QoS Levels ===")
	createQoSLevels(ctx, cfg, authProvider)

	// Example 4: Update QoS
	fmt.Println("\n=== Update QoS ===")
	updateQoS(ctx, cfg, authProvider)

	// Example 5: QoS hierarchy and preemption
	fmt.Println("\n=== QoS Hierarchy and Preemption ===")
	demonstrateQoSHierarchy(ctx, cfg, authProvider)

	// Example 6: Resource limits and fair share
	fmt.Println("\n=== Resource Limits and Fair Share ===")
	demonstrateResourceLimits(ctx, cfg, authProvider)
}

// checkQoSSupport checks if the cluster supports QoS management
func checkQoSSupport(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	// Try different versions
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}
	
	for _, version := range versions {
		client, err := slurm.NewClientWithVersion(ctx, version,
			slurm.WithConfig(cfg),
			slurm.WithAuth(auth),
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
func listQoSConfigurations(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	// Create v0.0.43 client
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
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
		fmt.Printf("\nQoS: %s\n", qos.Name)
		fmt.Printf("  Description: %s\n", qos.Description)
		fmt.Printf("  Priority: %d\n", qos.Priority)
		fmt.Printf("  Preempt Mode: %s\n", qos.PreemptMode)
		fmt.Printf("  Max Jobs: %d (per user: %d)\n", qos.MaxJobs, qos.MaxJobsPerUser)
		fmt.Printf("  Max CPUs: %d (per user: %d)\n", qos.MaxCPUs, qos.MaxCPUsPerUser)
		fmt.Printf("  Max Wall Time: %d hours\n", qos.MaxWallTime/3600)
		fmt.Printf("  Usage Factor: %.2f\n", qos.UsageFactor)
		
		if len(qos.Flags) > 0 {
			fmt.Printf("  Flags: %v\n", qos.Flags)
		}
		if len(qos.AllowedAccounts) > 0 {
			fmt.Printf("  Allowed Accounts: %v\n", qos.AllowedAccounts)
		}
	}
	
	// List QoS for specific accounts
	fmt.Println("\nListing QoS for specific accounts:")
	accountQoS, err := client.QoS().List(ctx, &interfaces.ListQoSOptions{
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
func createQoSLevels(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
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
	
	highPriorityQoS := &interfaces.QoSCreate{
		Name:               "high-priority",
		Description:        "High priority for critical jobs",
		Priority:           10000,
		PreemptMode:        "requeue",
		GraceTime:          300, // 5 minutes
		MaxJobs:            50,
		MaxJobsPerUser:     10,
		MaxJobsPerAccount:  25,
		MaxCPUs:            500,
		MaxCPUsPerUser:     100,
		MaxNodes:           25,
		MaxWallTime:        86400, // 24 hours
		UsageFactor:        2.0,   // Double charge for priority
		Flags:              []string{"DenyOnLimit", "RequireAssoc"},
		AllowedAccounts:    []string{"critical-research", "production"},
	}
	
	resp, err := client.QoS().Create(ctx, highPriorityQoS)
	if err != nil {
		log.Printf("Failed to create high-priority QoS: %v", err)
	} else {
		fmt.Printf("Created QoS: %s\n", resp.QoSName)
	}
	
	// Example 2: Create normal QoS
	fmt.Println("\nCreating normal QoS:")
	
	normalQoS := &interfaces.QoSCreate{
		Name:               "normal",
		Description:        "Standard QoS for regular jobs",
		Priority:           1000,
		PreemptMode:        "suspend",
		GraceTime:          600, // 10 minutes
		MaxJobs:            100,
		MaxJobsPerUser:     20,
		MaxJobsPerAccount:  50,
		MaxCPUs:            200,
		MaxCPUsPerUser:     50,
		MaxNodes:           10,
		MaxWallTime:        43200, // 12 hours
		UsageFactor:        1.0,
		Flags:              []string{"DenyOnLimit"},
	}
	
	resp2, err := client.QoS().Create(ctx, normalQoS)
	if err != nil {
		log.Printf("Failed to create normal QoS: %v", err)
	} else {
		fmt.Printf("Created QoS: %s\n", resp2.QoSName)
	}
	
	// Example 3: Create low-priority/scavenger QoS
	fmt.Println("\nCreating scavenger QoS:")
	
	scavengerQoS := &interfaces.QoSCreate{
		Name:               "scavenger",
		Description:        "Low priority for opportunistic jobs",
		Priority:           100,
		PreemptMode:        "cancel",
		MaxJobs:            200,
		MaxJobsPerUser:     50,
		MaxCPUs:            100,
		MaxNodes:           5,
		MaxWallTime:        7200,  // 2 hours
		UsageFactor:        0.1,   // 10% charge - incentivize usage
		Flags:              []string{"NoReserve", "Preemptable"},
		DeniedAccounts:     []string{"critical-research"}, // Not for critical work
	}
	
	resp3, err := client.QoS().Create(ctx, scavengerQoS)
	if err != nil {
		log.Printf("Failed to create scavenger QoS: %v", err)
	} else {
		fmt.Printf("Created QoS: %s\n", resp3.QoSName)
	}
}

// updateQoS demonstrates updating QoS configurations
func updateQoS(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
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
	fmt.Printf("  Max CPUs: %d\n", current.MaxCPUs)
	fmt.Printf("  Max Wall Time: %d hours\n", current.MaxWallTime/3600)
	
	// Update QoS - increase limits
	newMaxCPUs := 300
	newMaxWallTime := 72000 // 20 hours
	update := &interfaces.QoSUpdate{
		MaxCPUs:     &newMaxCPUs,
		MaxWallTime: &newMaxWallTime,
		Description: stringPtr("Updated normal QoS with increased limits"),
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
	fmt.Printf("  Max CPUs: %d\n", updated.MaxCPUs)
	fmt.Printf("  Max Wall Time: %d hours\n", updated.MaxWallTime/3600)
}

// demonstrateQoSHierarchy shows QoS hierarchy and preemption
func demonstrateQoSHierarchy(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
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
	executive := &interfaces.QoSCreate{
		Name:        "executive",
		Description: "Executive priority - preempts all",
		Priority:    100000,
		PreemptMode: "requeue",
		GraceTime:   60, // 1 minute grace
		MaxJobs:     10,
		MaxCPUs:     1000,
		Flags: []string{
			"PreemptExempt",    // Cannot be preempted
			"OverPartQOS",      // Override partition QoS
			"UsageFactorSafe",  // Protect from usage factor
		},
		AllowedUsers: []string{"ceo", "cto"},
	}
	
	_, err = client.QoS().Create(ctx, executive)
	if err != nil {
		log.Printf("Failed to create executive QoS: %v", err)
	} else {
		fmt.Println("  Created: executive (priority: 100000)")
	}
	
	// 2. Urgent QoS - high priority, can preempt normal and below
	urgent := &interfaces.QoSCreate{
		Name:        "urgent",
		Description: "Urgent priority - preempts normal and low",
		Priority:    50000,
		PreemptMode: "suspend",
		GraceTime:   300, // 5 minutes
		MaxJobs:     25,
		MaxCPUs:     500,
		Flags:       []string{"DenyOnLimit"},
		AllowedAccounts: []string{"operations", "critical-research"},
	}
	
	_, err = client.QoS().Create(ctx, urgent)
	if err != nil {
		log.Printf("Failed to create urgent QoS: %v", err)
	} else {
		fmt.Println("  Created: urgent (priority: 50000)")
	}
	
	// 3. Interactive QoS - for interactive/debug jobs
	interactive := &interfaces.QoSCreate{
		Name:        "interactive",
		Description: "Interactive jobs - quick turnaround",
		Priority:    5000,
		PreemptMode: "off",
		MaxJobs:     5,
		MaxJobsPerUser: 2,
		MaxCPUs:     16,
		MaxNodes:    1,
		MaxWallTime: 3600, // 1 hour max
		MinCPUs:     1,
		Flags:       []string{"NoReserve"}, // Don't make reservations
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
func demonstrateResourceLimits(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
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
	
	limitedQoS := &interfaces.QoSCreate{
		Name:               "limited-resources",
		Description:        "Strict resource limits for fair sharing",
		Priority:           1000,
		MaxJobs:            20,
		MaxJobsPerUser:     5,
		MaxJobsPerAccount:  10,
		MaxSubmitJobs:      50,    // Can queue up to 50
		MaxCPUs:            100,
		MaxCPUsPerUser:     25,
		MaxNodes:           5,
		MaxWallTime:        14400,  // 4 hours
		MinCPUs:            1,
		MinNodes:           1,
		UsageFactor:        1.5,    // 50% surcharge
		UsageThreshold:     0.8,    // Start limiting at 80% usage
		Flags: []string{
			"DenyOnLimit",      // Deny when limits reached
			"EnforceUsageThreshold", // Enforce the usage threshold
			"NoDecay",          // Don't decay priority
		},
	}
	
	_, err = client.QoS().Create(ctx, limitedQoS)
	if err != nil {
		log.Printf("Failed to create limited QoS: %v", err)
	} else {
		fmt.Println("Created QoS with strict limits")
	}
	
	// Example 2: Create QoS for GPU resources
	fmt.Println("\nCreating GPU-specific QoS:")
	
	gpuQoS := &interfaces.QoSCreate{
		Name:               "gpu-jobs",
		Description:        "QoS for GPU-accelerated jobs",
		Priority:           2000,
		MaxJobs:            10,
		MaxJobsPerUser:     2,
		MaxCPUs:            32,     // Limited CPUs per GPU job
		MaxNodes:           4,      // Max 4 GPU nodes
		MaxWallTime:        86400,  // 24 hours
		UsageFactor:        3.0,    // 3x charge for GPU resources
		Flags:              []string{"RequireAssoc", "PartitionQOS"},
		AllowedAccounts:    []string{"ml-research", "gpu-users"},
		// Note: Actual GPU limits would be set via GRES in job submission
	}
	
	_, err = client.QoS().Create(ctx, gpuQoS)
	if err != nil {
		log.Printf("Failed to create GPU QoS: %v", err)
	} else {
		fmt.Println("Created GPU-specific QoS")
	}
	
	// Example 3: Fair share configuration
	fmt.Println("\nDemonstrating fair share with QoS:")
	
	fairShareQoS := &interfaces.QoSCreate{
		Name:               "fair-share",
		Description:        "Fair share QoS with usage tracking",
		Priority:           1000,
		MaxJobs:            100,
		MaxCPUs:            200,
		UsageFactor:        1.0,
		UsageThreshold:     0.5,    // Start reducing priority at 50% share
		Flags: []string{
			"NoDecay",               // Maintain priority decay
			"UsageFactorSafe",       // Protected from usage factor changes
			"PartitionTimeLimit",    // Use partition time limits
		},
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