// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/pkg/auth"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("=== Slurm REST API Multi-Version Client Demo ===")
	fmt.Println()

	// Display supported versions
	fmt.Printf("Supported API versions: %v\n", slurm.SupportedVersions())
	fmt.Printf("Latest version: %s\n", slurm.LatestVersion())
	fmt.Printf("Stable version: %s\n", slurm.StableVersion())
	fmt.Println()

	// Example 1: Create client with automatic version detection
	fmt.Println("--- Example 1: Automatic Version Detection ---")
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("http://localhost:6820"),
		slurm.WithAuth(auth.NewNoAuth()),
	)
	if err != nil {
		log.Printf("Failed to create client with auto-detection: %v\n", err)
	} else {
		fmt.Printf("Created client with version: %s\n", client.Version())
		client.Close()
	}
	fmt.Println()

	// Example 2: Create client for specific version
	fmt.Println("--- Example 2: Specific Version (v0.0.42) ---")
	client42, err := slurm.NewClientWithVersion(ctx, "v0.0.42",
		slurm.WithBaseURL("http://localhost:6820"),
		slurm.WithAuth(auth.NewNoAuth()),
	)
	if err != nil {
		log.Printf("Failed to create v0.0.42 client: %v\n", err)
	} else {
		fmt.Printf("Created client with version: %s\n", client42.Version())

		// Demonstrate API operations
		fmt.Println("Testing job operations...")

		// List jobs
		jobs, err := client42.Jobs().List(ctx, &slurm.ListJobsOptions{
			States: []string{"RUNNING"},
			Limit:  10,
		})
		if err != nil {
			log.Printf("Failed to list jobs: %v\n", err)
		} else {
			fmt.Printf("Found %d running jobs\n", len(jobs.Jobs))
		}

		// Test connectivity
		err = client42.Info().Ping(ctx)
		if err != nil {
			log.Printf("Ping failed: %v\n", err)
		} else {
			fmt.Println("Connectivity test: OK")
		}

		client42.Close()
	}
	fmt.Println()

	// Example 3: Version compatibility information
	fmt.Println("--- Example 3: Version Compatibility ---")
	compatibility := slurm.GetVersionCompatibility()

	for version, slurmVersions := range compatibility.SlurmVersions {
		fmt.Printf("API %s is compatible with Slurm versions: %v\n",
			version, slurmVersions)
	}
	fmt.Println()

	// Example 4: Create client for Slurm version
	fmt.Println("--- Example 4: Client for Slurm Version ---")
	clientForSlurm, err := slurm.NewClientForSlurmVersion(ctx, "25.05",
		slurm.WithBaseURL("http://localhost:6820"),
		slurm.WithAuth(auth.NewNoAuth()),
	)
	if err != nil {
		log.Printf("Failed to create client for Slurm 25.05: %v\n", err)
	} else {
		fmt.Printf("Created client with version: %s for Slurm 25.05\n",
			clientForSlurm.Version())
		clientForSlurm.Close()
	}
	fmt.Println()

	// Example 5: Job submission workflow
	fmt.Println("--- Example 5: Job Submission Workflow ---")
	if client42 != nil {
		// Submit a test job
		submission := &slurm.JobSubmission{
			Name:      "demo-job",
			Script:    "#!/bin/bash\necho 'Hello from Slurm!'\nsleep 30",
			Partition: "debug",
			CPUs:      1,
			Memory:    1024,
			TimeLimit: 300, // 5 minutes
		}

		response, err := client42.Jobs().Submit(ctx, submission)
		if err != nil {
			log.Printf("Failed to submit job: %v\n", err)
		} else {
			fmt.Printf("Submitted job with ID: %s\n", response.JobID)

			// Get job details
			job, err := client42.Jobs().Get(ctx, response.JobID)
			if err != nil {
				log.Printf("Failed to get job details: %v\n", err)
			} else {
				fmt.Printf("Job state: %s\n", job.State)
				fmt.Printf("Job submitted at: %s\n", job.SubmitTime.Format(time.RFC3339))
			}
		}
	}

	fmt.Println()
	fmt.Println("=== Demo Complete ===")
}
