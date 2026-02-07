//go:build integration
// +build integration

// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	slurm "github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

// TestImplementationStatus tests which methods are actually implemented
func TestImplementationStatus(t *testing.T) {
	if os.Getenv("SLURM_REAL_SERVER_TEST") != "true" {
		t.Skip("Real server tests disabled. Set SLURM_REAL_SERVER_TEST=true to enable")
	}

	ctx := context.Background()
	token := os.Getenv("SLURM_JWT_TOKEN")
	if token == "" {
		// Try to fetch token
		tokenBytes, err := fetchJWTTokenViaSSH()
		if err != nil {
			t.Skipf("Could not fetch JWT token: %v", err)
		}
		token = tokenBytes
	}

	serverURL := os.Getenv("SLURM_SERVER_URL")
	if serverURL == "" {
		serverURL = "http://localhost
	}

	versions := []string{"v0.0.42", "v0.0.43"}

	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			client, err := slurm.NewClientWithVersion(ctx, version,
				slurm.WithBaseURL(serverURL),
				slurm.WithAuth(auth.NewTokenAuth(token)),
				slurm.WithConfig(&config.Config{
					Timeout:            30 * time.Second,
					MaxRetries:         3,
					Debug:              true,
					InsecureSkipVerify: true,
				}),
			)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}
			defer client.Close()

			t.Logf("Testing %s implementation status", version)

			// Test Info Manager
			t.Run("InfoManager", func(t *testing.T) {
				t.Run("Ping", func(t *testing.T) {
					err := client.Info().Ping(ctx)
					if err != nil {
						if strings.Contains(err.Error(), "nil pointer") {
							t.Error("Method not implemented (nil pointer)")
						} else {
							t.Logf("Method implemented but failed: %v", err)
						}
					} else {
						t.Log("✓ Ping implemented and working")
					}
				})

				t.Run("Get", func(t *testing.T) {
					info, err := client.Info().Get(ctx)
					if err != nil {
						if strings.Contains(err.Error(), "nil pointer") {
							t.Error("Method not implemented (nil pointer)")
						} else {
							t.Logf("Method implemented but failed: %v", err)
						}
					} else if info == nil {
						t.Error("Method returns nil without error")
					} else {
						t.Logf("✓ Get implemented and working: cluster=%s", info.ClusterName)
					}
				})

				t.Run("Version", func(t *testing.T) {
					ver, err := client.Info().Version(ctx)
					if err != nil {
						if strings.Contains(err.Error(), "nil pointer") {
							t.Error("Method not implemented (nil pointer)")
						} else {
							t.Logf("Method implemented but failed: %v", err)
						}
					} else if ver == nil {
						t.Error("Method returns nil without error")
					} else {
						t.Logf("✓ Version implemented and working: %s", ver.Version)
					}
				})

				t.Run("Stats", func(t *testing.T) {
					stats, err := client.Info().Stats(ctx)
					if err != nil {
						if strings.Contains(err.Error(), "nil pointer") {
							t.Error("Method not implemented (nil pointer)")
						} else {
							t.Logf("Method implemented but failed: %v", err)
						}
					} else if stats == nil {
						t.Error("Method returns nil without error")
					} else {
						t.Logf("✓ Stats implemented and working: nodes=%d", stats.TotalNodes)
					}
				})
			})

			// Test Job Manager
			t.Run("JobManager", func(t *testing.T) {
				t.Run("List", func(t *testing.T) {
					jobs, err := client.Jobs().List(ctx, &slurm.ListJobsOptions{Limit: 5})
					if err != nil {
						if strings.Contains(err.Error(), "nil pointer") {
							t.Error("Method not implemented (nil pointer)")
						} else {
							t.Logf("Method implemented but failed: %v", err)
						}
					} else if jobs == nil {
						t.Error("Method returns nil without error")
					} else {
						t.Logf("✓ List implemented and working: %d jobs", len(jobs.Jobs))
					}
				})
			})

			// Test Node Manager
			t.Run("NodeManager", func(t *testing.T) {
				t.Run("List", func(t *testing.T) {
					nodes, err := client.Nodes().List(ctx, &slurm.ListNodesOptions{Limit: 5})
					if err != nil {
						if strings.Contains(err.Error(), "nil pointer") {
							t.Error("Method not implemented (nil pointer)")
						} else {
							t.Logf("Method implemented but failed: %v", err)
						}
					} else if nodes == nil {
						t.Error("Method returns nil without error")
					} else {
						t.Logf("✓ List implemented and working: %d nodes", len(nodes.Nodes))
					}
				})
			})

			// Test Partition Manager
			t.Run("PartitionManager", func(t *testing.T) {
				t.Run("List", func(t *testing.T) {
					partitions, err := client.Partitions().List(ctx, &slurm.ListPartitionsOptions{Limit: 5})
					if err != nil {
						if strings.Contains(err.Error(), "nil pointer") {
							t.Error("Method not implemented (nil pointer)")
						} else {
							t.Logf("Method implemented but failed: %v", err)
						}
					} else if partitions == nil {
						t.Error("Method returns nil without error")
					} else {
						t.Logf("✓ List implemented and working: %d partitions", len(partitions.Partitions))
					}
				})
			})
		})
	}
}

// TestV42WithRealServer specifically tests v0.0.42 implementations
func TestV42WithRealServer(t *testing.T) {
	if os.Getenv("SLURM_REAL_SERVER_TEST") != "true" {
		t.Skip("Real server tests disabled")
	}

	ctx := context.Background()
	token := os.Getenv("SLURM_JWT_TOKEN")
	if token == "" {
		tokenBytes, err := fetchJWTTokenViaSSH()
		if err != nil {
			t.Skipf("Could not fetch JWT token: %v", err)
		}
		token = tokenBytes
	}

	// Force v0.0.42 which has implementations
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.42",
		slurm.WithBaseURL("http://localhost
		slurm.WithAuth(auth.NewTokenAuth(token)),
		slurm.WithConfig(&config.Config{
			Timeout:            30 * time.Second,
			MaxRetries:         3,
			Debug:              false,
			InsecureSkipVerify: true,
		}),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Test implemented functionality
	t.Run("Ping", func(t *testing.T) {
		err := client.Info().Ping(ctx)
		if err != nil {
			t.Errorf("Ping failed: %v", err)
		}
	})

	t.Run("ClusterInfo", func(t *testing.T) {
		info, err := client.Info().Get(ctx)
		if err != nil {
			t.Errorf("Get cluster info failed: %v", err)
		} else {
			t.Logf("Cluster: %s", info.ClusterName)
		}
	})

	t.Run("ListJobs", func(t *testing.T) {
		jobs, err := client.Jobs().List(ctx, &slurm.ListJobsOptions{
			Limit:  10,
			States: []string{"RUNNING", "PENDING"},
		})
		if err != nil {
			t.Errorf("List jobs failed: %v", err)
		} else {
			t.Logf("Found %d jobs", len(jobs.Jobs))
			for i, job := range jobs.Jobs {
				if i < 3 {
					jobID := "unknown"
					if job.JobID != nil {
						jobID = fmt.Sprintf("%d", *job.JobID)
					}
					jobName := "unknown"
					if job.Name != nil {
						jobName = *job.Name
					}
					jobState := "unknown"
					if len(job.JobState) > 0 {
						jobState = string(job.JobState[0])
					}
					t.Logf("  Job %s: %s (%s)", jobID, jobName, jobState)
				}
			}
		}
	})

	t.Run("ListNodes", func(t *testing.T) {
		nodes, err := client.Nodes().List(ctx, &slurm.ListNodesOptions{
			Limit: 10,
		})
		if err != nil {
			t.Errorf("List nodes failed: %v", err)
		} else {
			t.Logf("Found %d nodes", len(nodes.Nodes))
			for i, node := range nodes.Nodes {
				if i < 3 {
					t.Logf("  Node %s: %s (Cpus: %d)", node.Name, node.State, node.CPUs)
				}
			}
		}
	})

	t.Run("SubmitJob", func(t *testing.T) {
		submission := &slurm.JobSubmission{
			Name:      fmt.Sprintf("test-job-%d", time.Now().Unix()),
			Script:    "#!/bin/bash\necho 'Hello from Go client test'\nhostname\ndate",
			Partition: "compute",
			Nodes:     1,
			CPUs:      1,
			TimeLimit: 1, // 1 minute
		}

		resp, err := client.Jobs().Submit(ctx, submission)
		if err != nil {
			t.Errorf("Submit job failed: %v", err)
		} else {
			t.Logf("Submitted job: %d", resp.JobId)

			// Try to cancel it
			err = client.Jobs().Cancel(ctx, fmt.Sprintf("%d", resp.JobId))
			if err != nil {
				t.Logf("Cancel job failed (might be normal if job completed): %v", err)
			} else {
				t.Logf("Cancelled job: %d", resp.JobId)
			}
		}
	})
}
