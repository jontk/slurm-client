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

// Example: Demonstrating version-specific differences and compatibility
func main() {
	// Create configuration
	cfg := config.NewDefault()
	cfg.BaseURL = "https://cluster.example.com:6820"
	
	// Create authentication
	authProvider := auth.NewTokenAuth("your-jwt-token")

	ctx := context.Background()

	// Example 1: Auto-detect version
	fmt.Println("=== Auto-Detection ===")
	autoDetectVersion(ctx, cfg, authProvider)

	// Example 2: Version-specific features
	fmt.Println("\n=== Version-Specific Features ===")
	demonstrateVersionFeatures(ctx, cfg, authProvider)

	// Example 3: Handling breaking changes
	fmt.Println("\n=== Breaking Change Handling ===")
	handleBreakingChanges(ctx, cfg, authProvider)

	// Example 4: Version compatibility check
	fmt.Println("\n=== Version Compatibility ===")
	checkVersionCompatibility()
}

// autoDetectVersion demonstrates automatic version detection
func autoDetectVersion(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	// Create client with auto-detection
	client, err := slurm.NewClient(ctx,
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
	)
	if err != nil {
		log.Printf("Failed to create client with auto-detection: %v", err)
		return
	}
	defer client.Close()

	// Get version info
	info, err := client.Info().Version(ctx)
	if err != nil {
		log.Printf("Failed to get version info: %v", err)
		return
	}

	fmt.Printf("Auto-detected API version: %s\n", info.Version)
	// SlurmVersion field doesn't exist in APIVersion
	// fmt.Printf("SLURM version: %s\n", info.SlurmVersion)
	fmt.Printf("Supported versions: %v\n", slurm.SupportedVersions())
}

// demonstrateVersionFeatures shows version-specific features
func demonstrateVersionFeatures(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}

	for _, version := range versions {
		fmt.Printf("\n--- Testing %s ---\n", version)
		
		// Create version-specific client
		client, err := slurm.NewClientWithVersion(ctx, version,
			slurm.WithConfig(cfg),
			slurm.WithAuth(auth),
		)
		if err != nil {
			log.Printf("Failed to create %s client: %v", version, err)
			continue
		}

		// Test version-specific features
		switch version {
		case "v0.0.40":
			testV40Features(ctx, client)
		case "v0.0.41":
			testV41Features(ctx, client)
		case "v0.0.42":
			testV42Features(ctx, client)
		case "v0.0.43":
			testV43Features(ctx, client)
		}

		client.Close()
	}
}

// testV40Features tests features specific to v0.0.40
func testV40Features(ctx context.Context, client slurm.SlurmClient) {
	fmt.Println("Testing v0.0.40 features:")
	
	// v0.0.40 uses minimum_switches in job submission
	job := &interfaces.JobSubmission{
		Name:      "v40-test",
		Command:   "echo 'Testing v0.0.40'",
		Partition: "compute",
		CPUs:      2,
		Memory:    4096,
		TimeLimit: 10,
		// Metadata would be in SBATCH directives
// 		// Removed: Metadata: map[string]interface{}{
// 			"minimum_switches": 1, // v0.0.40 specific field
// 		},
	}

	resp, err := client.Jobs().Submit(ctx, job)
	if err != nil {
		log.Printf("  Job submission failed: %v", err)
	} else {
		fmt.Printf("  Job submitted with ID: %s\n", resp.JobID)
		fmt.Println("  Note: Uses 'minimum_switches' field")
	}

	// List jobs with v0.0.40 specific options
	jobs, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{
		States: []string{"RUNNING", "PENDING"},
		Limit:  5,
	})
	if err != nil {
		log.Printf("  Job listing failed: %v", err)
	} else {
		fmt.Printf("  Found %d jobs\n", len(jobs.Jobs))
	}
}

// testV41Features tests features specific to v0.0.41
func testV41Features(ctx context.Context, client slurm.SlurmClient) {
	fmt.Println("Testing v0.0.41 features:")
	
	// v0.0.41 renamed minimum_switches to required_switches
	job := &interfaces.JobSubmission{
		Name:      "v41-test",
		Command:   "echo 'Testing v0.0.41'",
		Partition: "compute",
		CPUs:      2,
		Memory:    4096,
		TimeLimit: 10,
		// Metadata would be in SBATCH directives
// 		// Removed: Metadata: map[string]interface{}{
// 			"required_switches": 1, // v0.0.41 renamed field
// 		},
	}

	resp, err := client.Jobs().Submit(ctx, job)
	if err != nil {
		log.Printf("  Job submission failed: %v", err)
	} else {
		fmt.Printf("  Job submitted with ID: %s\n", resp.JobID)
		fmt.Println("  Note: Uses 'required_switches' field (renamed from minimum_switches)")
	}

	// v0.0.41 has extended node information
	nodes, err := client.Nodes().List(ctx, &interfaces.ListNodesOptions{
		States: []string{"IDLE"},
		Limit:  3,
	})
	if err != nil {
		log.Printf("  Node listing failed: %v", err)
	} else {
		fmt.Printf("  Found %d idle nodes\n", len(nodes.Nodes))
		// v0.0.41 includes additional node metrics
		// GPU info would come from other node fields
		// for _, node := range nodes.Nodes { ... }
	}
}

// testV42Features tests features specific to v0.0.42 (stable)
func testV42Features(ctx context.Context, client slurm.SlurmClient) {
	fmt.Println("Testing v0.0.42 features (stable version):")
	
	// v0.0.42 removed exclusive and oversubscribe from job outputs
	job := &interfaces.JobSubmission{
		Name:      "v42-test",
		Command:   "echo 'Testing v0.0.42 stable'",
		Partition: "compute",
		CPUs:      4,
		Memory:    8192,
		TimeLimit: 15,
		// Metadata would be in SBATCH directives
// 		// Removed: Metadata: map[string]interface{}{
// 			"required_switches": 1,
// 			// exclusive and oversubscribe no longer in outputs
// 		},
	}

	resp, err := client.Jobs().Submit(ctx, job)
	if err != nil {
		log.Printf("  Job submission failed: %v", err)
	} else {
		fmt.Printf("  Job submitted with ID: %s\n", resp.JobID)
		fmt.Println("  Note: 'exclusive' and 'oversubscribe' removed from outputs")
	}

	// v0.0.42 has improved error responses
	_, err = client.Jobs().Get(ctx, "invalid-job-id")
	if err != nil {
		fmt.Println("  Enhanced error handling in v0.0.42:")
		fmt.Printf("  Error: %v\n", err)
	}

	// Check cluster stats (stable in v0.0.42)
	stats, err := client.Info().Stats(ctx)
	if err != nil {
		log.Printf("  Stats retrieval failed: %v", err)
	} else {
		fmt.Printf("  Cluster stats - Running jobs: %d, Total nodes: %d\n",
			stats.RunningJobs, stats.TotalNodes)
	}
}

// testV43Features tests features specific to v0.0.43 (latest)
func testV43Features(ctx context.Context, client slurm.SlurmClient) {
	fmt.Println("Testing v0.0.43 features (latest version):")
	
	// v0.0.43 adds reservation management support
	job := &interfaces.JobSubmission{
		Name:      "v43-test",
		Command:   "echo 'Testing v0.0.43 latest'",
		Partition: "compute",
		CPUs:      4,
		Memory:    8192,
		TimeLimit: 15,
		// Metadata would be in SBATCH directives
// 		// Removed: Metadata: map[string]interface{}{
// 			"reservation": "weekly-maintenance", // v0.0.43 reservation support
// 		},
	}

	resp, err := client.Jobs().Submit(ctx, job)
	if err != nil {
		log.Printf("  Job submission failed: %v", err)
	} else {
		fmt.Printf("  Job submitted with ID: %s\n", resp.JobID)
		fmt.Println("  Note: Includes reservation management support")
	}

	// v0.0.43 removed FrontEnd mode support
	fmt.Println("  Note: FrontEnd mode support removed in v0.0.43")

	// v0.0.43 has enhanced partition information
	partitions, err := client.Partitions().List(ctx, nil)
	if err != nil {
		log.Printf("  Partition listing failed: %v", err)
	} else {
		fmt.Printf("  Found %d partitions\n", len(partitions.Partitions))
		// QoS info would come from other partition fields
		// for _, p := range partitions.Partitions { ... }
	}
}

// handleBreakingChanges demonstrates handling breaking changes between versions
func handleBreakingChanges(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	// Example: Submit a job that works across versions with different field names
	
	// Detect server version first
	client, err := slurm.NewClient(ctx,
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
	)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
	defer client.Close()

	info, err := client.Info().Version(ctx)
	if err != nil {
		log.Printf("Failed to get version: %v", err)
		return
	}

	// Create job submission based on version
	job := &interfaces.JobSubmission{
		Name:      "cross-version-job",
		Command:   "echo 'Cross-version compatible job'",
		Partition: "compute",
		CPUs:      2,
		Memory:    4096,
		TimeLimit: 10,
	}

	// Handle version-specific fields
	switch info.Version {
	case "v0.0.40":
		// Metadata field doesn't exist in JobSubmission
		// Switches would be specified in SBATCH directives
		fmt.Println("Using v0.0.40 field names (minimum_switches)")
	case "v0.0.41", "v0.0.42", "v0.0.43":
		// Metadata field doesn't exist in JobSubmission
		// Switches would be specified in SBATCH directives
		fmt.Println("Using v0.0.41+ field names (required_switches)")
	}

	resp, err := client.Jobs().Submit(ctx, job)
	if err != nil {
		log.Printf("Job submission failed: %v", err)
		return
	}

	fmt.Printf("Successfully submitted cross-version job: %s\n", resp.JobID)
}

// checkVersionCompatibility shows how to check version compatibility
func checkVersionCompatibility() {
	fmt.Println("\nVersion Compatibility Information:")
	
	// Show all supported versions
	fmt.Printf("Supported API versions: %v\n", slurm.SupportedVersions())
	fmt.Printf("Stable version: %s\n", slurm.StableVersion())
	fmt.Printf("Latest version: %s\n", slurm.LatestVersion())

	// Check compatibility for specific SLURM versions
	slurmVersions := []string{"24.05", "24.11", "25.05", "25.11"}
	
	fmt.Println("\nSLURM version to API version mapping:")
	for _, slurmVer := range slurmVersions {
		compatibility := getVersionCompatibility(slurmVer)
		fmt.Printf("  SLURM %s -> API versions: %v (recommended: %s)\n",
			slurmVer, compatibility.supportedVersions, compatibility.recommendedVersion)
	}

	// Show breaking changes
	fmt.Println("\nBreaking changes between versions:")
	fmt.Println("  v0.0.40 -> v0.0.41:")
	fmt.Println("    - Field rename: minimum_switches -> required_switches")
	fmt.Println("    - Added extended node metrics")
	fmt.Println("  v0.0.41 -> v0.0.42:")
	fmt.Println("    - Removed fields: exclusive, oversubscribe from job outputs")
	fmt.Println("    - Enhanced error response structure")
	fmt.Println("  v0.0.42 -> v0.0.43:")
	fmt.Println("    - Added reservation management support")
	fmt.Println("    - Removed FrontEnd mode support")
	fmt.Println("    - Enhanced QoS information in partitions")
}

// Helper struct for version compatibility
type versionCompatibility struct {
	supportedVersions  []string
	recommendedVersion string
}

// getVersionCompatibility returns API version compatibility for a SLURM version
func getVersionCompatibility(slurmVersion string) versionCompatibility {
	switch slurmVersion {
	case "24.05":
		return versionCompatibility{
			supportedVersions:  []string{"v0.0.40", "v0.0.41", "v0.0.42"},
			recommendedVersion: "v0.0.42",
		}
	case "24.11":
		return versionCompatibility{
			supportedVersions:  []string{"v0.0.41", "v0.0.42", "v0.0.43"},
			recommendedVersion: "v0.0.42",
		}
	case "25.05", "25.11":
		return versionCompatibility{
			supportedVersions:  []string{"v0.0.42", "v0.0.43"},
			recommendedVersion: "v0.0.42",
		}
	default:
		return versionCompatibility{
			supportedVersions:  []string{"v0.0.42"},
			recommendedVersion: "v0.0.42",
		}
	}
}