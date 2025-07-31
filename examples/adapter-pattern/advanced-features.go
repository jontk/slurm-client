package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
	"github.com/jontk/slurm-client/pkg/errors"
	"github.com/jontk/slurm-client/internal/interfaces"
)

// AdvancedAdapterFeatures demonstrates sophisticated adapter capabilities:
// - Advanced type conversion scenarios
// - Complex workflow management
// - Cross-version compatibility handling
// - Performance optimization techniques
func main() {
	ctx := context.Background()

	// Create client with advanced adapter configuration
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://your-slurm-server:6820"),
		slurm.WithAuth(auth.NewTokenAuth("your-jwt-token")),
		slurm.WithAdapterConfig(&config.AdapterConfig{
			EnableTypeCache:       true,
			EnableResponseCache:   true,
			CacheTimeout:         10 * time.Minute,
			StrictTypeConversion: false, // Allow version compatibility
			PreferZeroCopy:       true,  // Memory optimization
			AutoVersionFallback:  true,  // Graceful version handling
			ProvideVersionContext: true,
		}),
	)
	if err != nil {
		log.Fatalf("Failed to create SLURM client: %v", err)
	}
	defer client.Close()

	fmt.Printf("Advanced Adapter Features Demo - API %s\n", client.Version())

	// Example 1: Complex Reservation Management
	demonstrateReservationManagement(ctx, client)

	// Example 2: Account and Association Management
	demonstrateAccountManagement(ctx, client)

	// Example 3: Batch Operations with Error Recovery
	demonstrateBatchOperations(ctx, client)

	// Example 4: Version Migration Scenarios
	demonstrateVersionMigration(ctx, client)

	// Example 5: Performance Optimization
	demonstratePerformanceOptimization(ctx, client)
}

// demonstrateReservationManagement shows complex reservation operations
// with automatic type conversion and version compatibility
func demonstrateReservationManagement(ctx context.Context, client slurm.Client) {
	fmt.Println("\n=== Advanced Reservation Management ===")

	// Check if reservations are supported in current API version
	if client.Reservations() == nil {
		fmt.Printf("Reservations not supported in API %s - graceful degradation\n", client.Version())
		return
	}

	// Create a complex reservation with all available features
	reservation := &interfaces.ReservationCreate{
		Name:      fmt.Sprintf("maintenance-%d", time.Now().Unix()),
		StartTime: time.Now().Add(1 * time.Hour),
		Duration:  4 * 3600, // 4 hours
		Nodes:     []string{"node[001-010]", "gpu[001-004]"},
		Users:     []string{"admin", "maintenance"},
		Accounts:  []string{"operations", "infrastructure"},
		Flags:     []string{"MAINT", "IGNORE_JOBS", "NO_HOLD_JOBS_AFTER_END"},
		Features:  []string{"maintenance", "exclusive"},
		Licenses: map[string]int{
			"matlab": 10,
			"ansys":  5,
		},
		PartitionName: "maintenance",
		BurstBuffer:   "type=nvme,capacity=1TB",
	}

	// Adapter automatically converts complex types and handles version differences
	response, err := client.Reservations().Create(ctx, reservation)
	if err != nil {
		if slurmErr, ok := err.(*errors.SlurmError); ok {
			fmt.Printf("Reservation creation failed [%s]: %s\n", slurmErr.APIVersion, slurmErr.Message)
			
			// Check for version-specific limitations
			if errors.IsVersionNotSupported(err) {
				fmt.Println("→ Feature not available in this API version")
				return
			}
		}
		log.Printf("Reservation error: %v", err)
		return
	}

	fmt.Printf("Created reservation: %s\n", response.ReservationName)

	// List and examine reservations with complex filtering
	reservations, err := client.Reservations().List(ctx, &interfaces.ListReservationsOptions{
		Names:     []string{response.ReservationName},
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		log.Printf("Failed to list reservations: %v", err)
		return
	}

	for _, res := range reservations.Reservations {
		fmt.Printf("Reservation Details:\n")
		fmt.Printf("  Name: %s\n", res.Name)
		fmt.Printf("  State: %s\n", res.State)
		fmt.Printf("  Nodes: %v\n", res.Nodes)
		fmt.Printf("  Users: %v\n", res.Users)
		fmt.Printf("  Flags: %v\n", res.Flags)
		fmt.Printf("  Licenses: %v\n", res.Licenses)
	}

	// Clean up
	err = client.Reservations().Delete(ctx, response.ReservationName)
	if err != nil {
		log.Printf("Failed to delete reservation: %v", err)
	} else {
		fmt.Printf("Cleaned up reservation: %s\n", response.ReservationName)
	}
}

// demonstrateAccountManagement shows account hierarchy management
func demonstrateAccountManagement(ctx context.Context, client slurm.Client) {
	fmt.Println("\n=== Advanced Account Management ===")

	if client.Accounts() == nil {
		fmt.Printf("Account management not supported in API %s\n", client.Version())
		return
	}

	// Create a complex account hierarchy
	rootAccount := &interfaces.AccountCreate{
		Name:         fmt.Sprintf("research-dept-%d", time.Now().Unix()%1000),
		Description:  "Research Department with Adapter Demo",
		Organization: "Adapter Test Corp",
		CoordinatorUsers: []string{"dept-head", "admin"},
		AllowedPartitions: []string{"compute", "gpu", "highmem"},
		MaxJobs:       500,
		MaxNodes:      100,
		SharesPriority: 200,
		MaxTRES: map[string]int{
			"cpu":    2000,
			"mem":    8192000, // 8TB in MB
			"gpu":    20,
			"node":   50,
		},
		GrpTRES: map[string]int{
			"cpu":    1000,
			"mem":    4096000, // 4TB in MB
			"gpu":    10,
		},
		Flags: []string{"AllowSubmit"},
		DefaultQoS: "normal",
		AllowedQoS: []string{"normal", "high", "low"},
	}

	// Adapter handles complex type conversions automatically
	resp, err := client.Accounts().Create(ctx, rootAccount)
	if err != nil {
		if slurmErr, ok := err.(*errors.SlurmError); ok {
			fmt.Printf("Account creation failed [%s]: %s\n", slurmErr.APIVersion, slurmErr.Message)
		}
		return
	}

	fmt.Printf("Created account: %s\n", resp.AccountName)

	// Create sub-accounts with different configurations
	subAccounts := []struct {
		name   string
		config *interfaces.AccountCreate
	}{
		{
			name: "ml-project",
			config: &interfaces.AccountCreate{
				Name:            "ml-project",
				Description:     "Machine Learning Project",
				ParentAccount:   resp.AccountName,
				DefaultPartition: "gpu",
				MaxJobs:         100,
				MaxJobsPerUser:  20,
				MaxTRES: map[string]int{
					"gpu": 8,
					"cpu": 400,
				},
			},
		},
		{
			name: "simulation",
			config: &interfaces.AccountCreate{
				Name:            "simulation",
				Description:     "High-Performance Simulation",
				ParentAccount:   resp.AccountName,
				DefaultPartition: "highmem",
				MaxJobs:         200,
				MaxTRES: map[string]int{
					"mem":  2048000, // 2TB
					"cpu":  800,
					"node": 20,
				},
			},
		},
	}

	for _, sub := range subAccounts {
		subResp, err := client.Accounts().Create(ctx, sub.config)
		if err != nil {
			log.Printf("Failed to create sub-account %s: %v", sub.name, err)
			continue
		}
		fmt.Printf("Created sub-account: %s\n", subResp.AccountName)
	}

	// List account hierarchy
	accounts, err := client.Accounts().List(ctx, &interfaces.ListAccountsOptions{
		WithAssociations: true,
		WithCoordinators: true,
		ParentAccounts:   []string{resp.AccountName},
	})
	if err != nil {
		log.Printf("Failed to list accounts: %v", err)
		return
	}

	fmt.Printf("Account Hierarchy:\n")
	for _, account := range accounts.Accounts {
		fmt.Printf("  Account: %s (Parent: %s)\n", account.Name, account.ParentAccount)
		fmt.Printf("    Max Jobs: %d, Max Nodes: %d\n", account.MaxJobs, account.MaxNodes)
		fmt.Printf("    TRES Limits: %v\n", account.MaxTRES)
	}
}

// demonstrateBatchOperations shows batch processing with error recovery
func demonstrateBatchOperations(ctx context.Context, client slurm.Client) {
	fmt.Println("\n=== Batch Operations with Error Recovery ===")

	// Submit multiple jobs with different configurations
	jobTemplates := []struct {
		name   string
		config *interfaces.JobSubmission
	}{
		{
			name: "cpu-intensive",
			config: &interfaces.JobSubmission{
				Name:      "cpu-job",
				Script:    "#!/bin/bash\nstress-ng --cpu 4 --timeout 60s",
				Partition: "compute",
				CPUs:      4,
				Memory:    2 * 1024 * 1024 * 1024, // 2GB
				TimeLimit: 10,
			},
		},
		{
			name: "memory-intensive",
			config: &interfaces.JobSubmission{
				Name:      "memory-job",
				Script:    "#!/bin/bash\nstress-ng --vm 1 --vm-bytes 4G --timeout 60s",
				Partition: "compute",
				CPUs:      1,
				Memory:    4 * 1024 * 1024 * 1024, // 4GB
				TimeLimit: 10,
			},
		},
		{
			name: "gpu-job",
			config: &interfaces.JobSubmission{
				Name:      "gpu-job",
				Script:    "#!/bin/bash\nnvidia-smi; sleep 30",
				Partition: "gpu",
				CPUs:      2,
				Memory:    1 * 1024 * 1024 * 1024, // 1GB
				GPUs:      1,
				TimeLimit: 5,
			},
		},
	}

	var submittedJobs []string
	successCount := 0

	for _, template := range jobTemplates {
		response, err := client.Jobs().Submit(ctx, template.config)
		if err != nil {
			// Adapter provides detailed error context
			if slurmErr, ok := err.(*errors.SlurmError); ok {
				fmt.Printf("Job %s failed [%s]: %s\n", template.name, slurmErr.APIVersion, slurmErr.Message)
				
				// Check for resource limitations and retry with adjusted parameters
				if slurmErr.Code == errors.ErrorCodeResourceExhausted {
					fmt.Printf("→ Resource exhausted for %s, could retry with lower requirements\n", template.name)
				}
			}
			continue
		}

		submittedJobs = append(submittedJobs, response.JobID)
		successCount++
		fmt.Printf("Submitted %s: JobID=%s\n", template.name, response.JobID)
	}

	fmt.Printf("Batch operation summary: %d/%d jobs submitted successfully\n", successCount, len(jobTemplates))

	// Monitor batch job status
	if len(submittedJobs) > 0 {
		fmt.Println("Monitoring job status...")
		for _, jobID := range submittedJobs {
			job, err := client.Jobs().Get(ctx, jobID)
			if err != nil {
				fmt.Printf("Failed to get status for job %s: %v\n", jobID, err)
				continue
			}
			fmt.Printf("Job %s: State=%s, StartTime=%s\n", jobID, job.State, job.StartTime)
		}
	}
}

// demonstrateVersionMigration shows handling of version differences
func demonstrateVersionMigration(ctx context.Context, client slurm.Client) {
	fmt.Println("\n=== Version Migration Scenarios ===")

	currentVersion := client.Version()
	fmt.Printf("Current API Version: %s\n", currentVersion)

	// Demonstrate feature availability checking
	features := map[string]bool{
		"Job Management":         client.Jobs() != nil,
		"Node Management":        client.Nodes() != nil,
		"Partition Management":   client.Partitions() != nil,
		"Reservation Management": client.Reservations() != nil,
		"QoS Management":         client.QoS() != nil,
		"Account Management":     client.Accounts() != nil,
		"Association Management": client.Associations() != nil,
		"User Management":        client.Users() != nil,
	}

	fmt.Println("Feature Availability Matrix:")
	for feature, available := range features {
		status := "❌ Not Available"
		if available {
			status = "✅ Available"
		}
		fmt.Printf("  %s: %s\n", feature, status)
	}

	// Demonstrate graceful degradation
	if client.Reservations() == nil {
		fmt.Println("\nGraceful Degradation Example:")
		fmt.Println("→ Reservations not available, using alternative workflow:")
		fmt.Println("  1. Create maintenance job instead of reservation")
		fmt.Println("  2. Use node drain commands via job submission")
		fmt.Println("  3. Schedule maintenance through job dependencies")
	}
}

// demonstratePerformanceOptimization shows adapter performance features
func demonstratePerformanceOptimization(ctx context.Context, client slurm.Client) {
	fmt.Println("\n=== Performance Optimization ===")

	// Demonstrate caching benefits
	start := time.Now()
	
	// First call - no cache
	info1, err := client.Info().Get(ctx)
	if err == nil {
		duration1 := time.Since(start)
		fmt.Printf("First info call (no cache): %v\n", duration1)
		
		// Second call - should use cache
		start2 := time.Now()
		info2, err := client.Info().Get(ctx)
		if err == nil {
			duration2 := time.Since(start2)
			fmt.Printf("Second info call (cached): %v\n", duration2)
			
			if duration2 < duration1 {
				fmt.Printf("Cache speedup: %.2fx faster\n", float64(duration1)/float64(duration2))
			}
			
			// Verify cache consistency
			if info1.ClusterName == info2.ClusterName {
				fmt.Println("✅ Cache consistency verified")
			}
		}
	}

	// Demonstrate batch operations performance
	fmt.Println("\nBatch vs Individual Operations:")
	
	// Individual calls
	start = time.Now()
	for i := 0; i < 5; i++ {
		_, _ = client.Info().Ping(ctx)
	}
	individualDuration := time.Since(start)
	fmt.Printf("5 individual ping calls: %v\n", individualDuration)

	// If batch operations were available, they would be faster
	fmt.Println("→ Adapter optimizations include:")
	fmt.Println("  • Connection pooling and reuse")
	fmt.Println("  • Response caching for expensive operations")
	fmt.Println("  • Zero-copy type conversions where possible")
	fmt.Println("  • Intelligent request batching")

	fmt.Println("\nAdapter Pattern Benefits Summary:")
	fmt.Println("✅ Version abstraction - works across all API versions")
	fmt.Println("✅ Type safety - compile-time guarantees")
	fmt.Println("✅ Performance - caching and optimization")
	fmt.Println("✅ Error context - version information in all errors")
	fmt.Println("✅ Future-proof - easy addition of new versions")
}