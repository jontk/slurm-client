//go:build integration
// +build integration

// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"net/http"
	"os"
	"testing"

	v042 "github.com/jontk/slurm-client/internal/openapi/v0_0_42"
)

// TestDirectAPIWithAuth tests the v0.0.42 API directly with proper authentication
func TestDirectAPIWithAuth(t *testing.T) {
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

	// Create client with request editor for authentication
	client, err := v042.NewClientWithResponses("http://localhost
		v042.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("X-SLURM-USER-TOKEN", token)
			return nil
		}),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	t.Run("Ping", func(t *testing.T) {
		resp, err := client.SlurmV0042GetPingWithResponse(ctx)
		if err != nil {
			t.Fatalf("Ping request failed: %v", err)
		}
		if resp.StatusCode() != 200 {
			t.Fatalf("Ping returned status %d: %s", resp.StatusCode(), string(resp.Body))
		}
		t.Logf("Ping successful: %d pings", len(resp.JSON200.Pings))
		for i, ping := range resp.JSON200.Pings {
			if i < 3 {
				t.Logf("  Controller %d: %s (%s)", i, *ping.Hostname, *ping.Pinged)
			}
		}
	})

	t.Run("GetJobs", func(t *testing.T) {
		resp, err := client.SlurmV0042GetJobsWithResponse(ctx, nil)
		if err != nil {
			t.Fatalf("GetJobs request failed: %v", err)
		}
		if resp.StatusCode() != 200 {
			t.Fatalf("GetJobs returned status %d: %s", resp.StatusCode(), string(resp.Body))
		}
		t.Logf("Found %d jobs", len(resp.JSON200.Jobs))
		for i, job := range resp.JSON200.Jobs {
			if i < 3 && job.Name != nil && job.JobId != nil && job.JobState != nil {
				state := "UNKNOWN"
				if len(*job.JobState) > 0 {
					state = (*job.JobState)[0]
				}
				t.Logf("  Job %d: %s (%s)", *job.JobId, *job.Name, state)
			}
		}
	})

	t.Run("GetNodes", func(t *testing.T) {
		resp, err := client.SlurmV0042GetNodesWithResponse(ctx, nil)
		if err != nil {
			t.Fatalf("GetNodes request failed: %v", err)
		}
		if resp.StatusCode() != 200 {
			t.Fatalf("GetNodes returned status %d: %s", resp.StatusCode(), string(resp.Body))
		}
		t.Logf("Found %d nodes", len(resp.JSON200.Nodes))
		for i, node := range resp.JSON200.Nodes {
			if i < 3 && node.Name != nil && node.State != nil && node.Cpus != nil {
				state := "UNKNOWN"
				if len(*node.State) > 0 {
					state = (*node.State)[0]
				}
				t.Logf("  Node %s: %s (Cpus: %d)", *node.Name, state, *node.Cpus)
			}
		}
	})

	t.Run("GetPartitions", func(t *testing.T) {
		resp, err := client.SlurmV0042GetPartitionsWithResponse(ctx, nil)
		if err != nil {
			t.Fatalf("GetPartitions request failed: %v", err)
		}
		if resp.StatusCode() != 200 {
			t.Fatalf("GetPartitions returned status %d: %s", resp.StatusCode(), string(resp.Body))
		}
		t.Logf("Found %d partitions", len(resp.JSON200.Partitions))
		for i, partition := range resp.JSON200.Partitions {
			if i < 3 && partition.Name != nil {
				t.Logf("  Partition %s", *partition.Name)
			}
		}
	})

	t.Run("Diag", func(t *testing.T) {
		resp, err := client.SlurmV0042GetDiagWithResponse(ctx)
		if err != nil {
			t.Fatalf("Diag request failed: %v", err)
		}
		if resp.StatusCode() != 200 {
			t.Fatalf("Diag returned status %d: %s", resp.StatusCode(), string(resp.Body))
		}
		if resp.JSON200.Statistics.PartsPacked != nil {
			t.Logf("Cluster statistics: Parts packed = %d", *resp.JSON200.Statistics.PartsPacked)
		}
	})
}
