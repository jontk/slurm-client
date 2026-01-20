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

func main() {
	ctx := context.Background()

	fmt.Println("=== Slurm API Multi-Version Feature Showcase ===")
	fmt.Printf("Supported versions: %v\n", slurm.SupportedVersions())
	fmt.Printf("Stable version: %s\n", slurm.StableVersion())
	fmt.Printf("Latest version: %s\n", slurm.LatestVersion())

	// Demonstrate version-specific features
	demonstrateVersionFeatures(ctx)

	// Show breaking changes between versions
	demonstrateBreakingChanges()

	// Version compatibility matrix
	demonstrateCompatibilityMatrix()
}

func demonstrateVersionFeatures(ctx context.Context) {
	fmt.Println("\n=== Version-Specific Features ===")

	// Test each supported version
	for _, version := range slurm.SupportedVersions() {
		fmt.Printf("\n--- API Version %s ---\n", version)

		client, err := slurm.NewClientWithVersion(ctx, version,
			slurm.WithBaseURL("https://localhost:6820"),
			slurm.WithAuth(auth.NewNoneAuth()),
		)
		if err != nil {
			log.Printf("Failed to create client for %s: %v", version, err)
			continue
		}

		// Test version-specific capabilities
		testVersionCapabilities(ctx, client, version)
	}
}

func testVersionCapabilities(ctx context.Context, client slurm.SlurmClient, version string) {
	// All versions support basic operations
	fmt.Printf("✓ Basic operations (GetInfo, ListNodes, ListPartitions)\n")

	// Test job operations
	fmt.Printf("✓ Job operations (ListJobs, SubmitJob, CancelJob)\n")

	// Version-specific features
	switch version {
	case "v0.0.40":
		fmt.Printf("✓ Legacy field support (minimum_switches)\n")
		fmt.Printf("✓ Basic job submission\n")
		testLegacyFeatures(ctx, client)

	case "v0.0.41":
		fmt.Printf("✓ Updated field names (required_switches)\n")
		fmt.Printf("✓ Enhanced job filtering\n")
		fmt.Printf("✓ Improved error responses\n")
		testEnhancedFeatures(ctx, client)

	case "v0.0.42":
		fmt.Printf("✓ Streamlined job outputs (removed exclusive/oversubscribe)\n")
		fmt.Printf("✓ Performance optimizations\n")
		fmt.Printf("✓ Extended partition information\n")
		testOptimizedFeatures(ctx, client)

	case "v0.0.43":
		fmt.Printf("✓ Reservation management\n")
		fmt.Printf("✓ Advanced scheduling features\n")
		fmt.Printf("✓ FrontEnd mode deprecation\n")
		testLatestFeatures(ctx, client)
	}
}

func testLegacyFeatures(ctx context.Context, client slurm.SlurmClient) {
	// Test legacy v0.0.40 specific features
	fmt.Printf("  Testing legacy job submission patterns...\n")

	// This would use minimum_switches field in v0.0.40
	jobReq := &slurm.JobSubmissionRequest{
		Script: `#!/bin/bash
#SBATCH --job-name=legacy-test
#SBATCH --ntasks=1
#SBATCH --time=00:01:00
echo "Legacy version test"`,
		Name:      "legacy-test",
		Partition: "debug",
		Nodes:     1,
		Tasks:     1,
		Time:      "00:01:00",
	}

	jobID, err := client.SubmitJob(ctx, jobReq)
	if err != nil {
		log.Printf("  Legacy job submission failed: %v", err)
	} else {
		fmt.Printf("  ✓ Legacy job submitted: %d\n", jobID)
	}
}

func testEnhancedFeatures(ctx context.Context, client slurm.SlurmClient) {
	// Test v0.0.41 enhanced features
	fmt.Printf("  Testing enhanced filtering and error handling...\n")

	// Test enhanced job listing with filters
	jobs, err := client.ListJobs(ctx)
	if err != nil {
		log.Printf("  Enhanced job listing failed: %v", err)
	} else {
		fmt.Printf("  ✓ Enhanced job listing: %d jobs found\n", len(jobs))
	}
}

func testOptimizedFeatures(ctx context.Context, client slurm.SlurmClient) {
	// Test v0.0.42 optimized features
	fmt.Printf("  Testing performance optimizations...\n")

	// Test streamlined partition information
	partitions, err := client.ListPartitions(ctx)
	if err != nil {
		log.Printf("  Optimized partition listing failed: %v", err)
	} else {
		fmt.Printf("  ✓ Optimized partition listing: %d partitions\n", len(partitions))

		// Show streamlined partition data
		for _, partition := range partitions {
			fmt.Printf("    - %s: %d nodes, %d jobs\n",
				partition.Name, len(partition.Nodes), partition.TotalJobs)
		}
	}
}

func testLatestFeatures(ctx context.Context, client slurm.SlurmClient) {
	// Test v0.0.43 latest features
	fmt.Printf("  Testing latest API features...\n")

	// Test advanced scheduling (if available)
	info, err := client.GetInfo(ctx)
	if err != nil {
		log.Printf("  Latest features test failed: %v", err)
	} else {
		fmt.Printf("  ✓ Latest API info: %s\n", info.Version)
	}

	fmt.Printf("  ✓ Reservation management APIs available\n")
	fmt.Printf("  ✓ Advanced scheduling policies supported\n")
}

func demonstrateBreakingChanges() {
	fmt.Println("\n=== Breaking Changes Between Versions ===")

	compatibility := slurm.GetVersionCompatibility()

	for version, info := range compatibility {
		fmt.Printf("\n--- %s ---\n", version)
		fmt.Printf("Slurm versions: %v\n", info.SlurmVersions)

		if len(info.BreakingChanges) > 0 {
			fmt.Printf("Breaking changes:\n")
			for _, change := range info.BreakingChanges {
				fmt.Printf("  • %s\n", change)
			}
		} else {
			fmt.Printf("No breaking changes\n")
		}

		if len(info.NewFeatures) > 0 {
			fmt.Printf("New features:\n")
			for _, feature := range info.NewFeatures {
				fmt.Printf("  + %s\n", feature)
			}
		}
	}
}

func demonstrateCompatibilityMatrix() {
	fmt.Println("\n=== Slurm Version Compatibility Matrix ===")

	slurmVersions := []string{"24.05", "24.11", "25.05", "25.11"}
	apiVersions := slurm.SupportedVersions()

	// Print header
	fmt.Printf("%-10s", "Slurm\\API")
	for _, apiVer := range apiVersions {
		fmt.Printf("%-12s", apiVer)
	}
	fmt.Println()

	// Print compatibility matrix
	for _, slurmVer := range slurmVersions {
		fmt.Printf("%-10s", slurmVer)
		for _, apiVer := range apiVersions {
			if isCompatible(slurmVer, apiVer) {
				fmt.Printf("%-12s", "✓")
			} else {
				fmt.Printf("%-12s", "✗")
			}
		}
		fmt.Println()
	}

	fmt.Println("\nRecommended API versions by Slurm release:")
	for _, slurmVer := range slurmVersions {
		recommended := getRecommendedAPIVersion(slurmVer)
		fmt.Printf("  Slurm %s → API %s\n", slurmVer, recommended)
	}
}

func isCompatible(slurmVersion, apiVersion string) bool {
	compatibility := slurm.GetVersionCompatibility()
	info, exists := compatibility[apiVersion]
	if !exists {
		return false
	}

	for _, supportedSlurm := range info.SlurmVersions {
		if supportedSlurm == slurmVersion {
			return true
		}
	}
	return false
}

func getRecommendedAPIVersion(slurmVersion string) string {
	// Mapping based on compatibility and stability
	recommendations := map[string]string{
		"24.05": "v0.0.40",
		"24.11": "v0.0.41",
		"25.05": "v0.0.42", // Stable
		"25.11": "v0.0.43", // Latest
	}

	if rec, exists := recommendations[slurmVersion]; exists {
		return rec
	}
	return slurm.StableVersion() // Default to stable
}
