// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// TestJobsWithRealServer tests job endpoints that don't require slurmdbd
func TestJobsWithRealServer(t *testing.T) {
	// Skip if not explicitly enabled
	if os.Getenv("SLURM_REAL_SERVER_TEST") != "true" {
		t.Skip("Real server tests disabled. Set SLURM_REAL_SERVER_TEST=true to enable")
	}

	// Get server configuration
	serverURL := os.Getenv("SLURM_SERVER_URL")
	if serverURL == "" {
		serverURL = "http://rocky9:6820"
	}

	token := os.Getenv("SLURM_JWT_TOKEN")
	if token == "" {
		// Use the token you provided
		token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTM3MzE3MjgsImlhdCI6MTc1MzcyOTkyOCwic3VuIjoicm9vdCJ9.7DGZd7hWhJQhkIx_0wMsKGM2rDipM27CGgaZFU1z_Ns"
	}

	// Create a transport that adds the token
	transport := &tokenTransport{
		token: token,
		base:  http.DefaultTransport,
	}
	
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	// Create API client
	apiClient, err := api.NewClientWithResponses(
		serverURL,
		api.WithHTTPClient(httpClient),
	)
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("List Jobs", func(t *testing.T) {
		resp, err := apiClient.SlurmV0043GetJobsWithResponse(ctx, &api.SlurmV0043GetJobsParams{})
		require.NoError(t, err)
		
		assert.Equal(t, 200, resp.StatusCode())
		
		if resp.JSON200 != nil {
			t.Logf("Found %d jobs on the server", len(resp.JSON200.Jobs))
			
			// Log first few jobs
			for i, job := range resp.JSON200.Jobs {
				if i >= 3 {
					break
				}
				
				jobID := "unknown"
				if job.JobId != nil {
					jobID = fmt.Sprintf("%d", *job.JobId)
				}
				
				jobName := "unknown"
				if job.Name != nil {
					jobName = *job.Name
				}
				
				state := "unknown"
				if job.JobState != nil {
					states := *job.JobState
					if len(states) > 0 {
						state = string(states[0])
					}
				}
				
				t.Logf("  Job %s: name=%s, state=%s", jobID, jobName, state)
			}
		}
	})

	t.Run("List Partitions", func(t *testing.T) {
		resp, err := apiClient.SlurmV0043GetPartitionsWithResponse(ctx, &api.SlurmV0043GetPartitionsParams{})
		require.NoError(t, err)
		
		// Skip if we get a 502 due to slurmdbd not being connected
		if resp.StatusCode() == 502 {
			t.Skip("Skipping partitions test: slurmdbd is not connected (HTTP 502)")
			return
		}
		
		assert.Equal(t, 200, resp.StatusCode())
		
		if resp.JSON200 != nil {
			t.Logf("Found %d partitions on the server", len(resp.JSON200.Partitions))
			
			for _, partition := range resp.JSON200.Partitions {
				name := "unknown"
				if partition.Name != nil {
					name = *partition.Name
				}
				
				state := "unknown" 
				if partition.Partition != nil && partition.Partition.State != nil {
					states := *partition.Partition.State
					if len(states) > 0 {
						state = string(states[0])
					}
				}
				
				nodes := "unknown"
				if partition.Nodes != nil && partition.Nodes.Configured != nil {
					nodes = *partition.Nodes.Configured
				}
				
				t.Logf("  Partition %s: state=%s, nodes=%s", name, state, nodes)
			}
		}
	})

	t.Run("Get Nodes", func(t *testing.T) {
		resp, err := apiClient.SlurmV0043GetNodesWithResponse(ctx, &api.SlurmV0043GetNodesParams{})
		require.NoError(t, err)
		
		assert.Equal(t, 200, resp.StatusCode())
		
		if resp.JSON200 != nil {
			t.Logf("Found %d nodes on the server", len(resp.JSON200.Nodes))
			
			for _, node := range resp.JSON200.Nodes {
				name := "unknown"
				if node.Name != nil {
					name = *node.Name
				}
				
				state := "unknown"
				if node.State != nil {
					states := *node.State
					if len(states) > 0 {
						state = string(states[0])
					}
				}
				
				t.Logf("  Node %s: state=%s", name, state)
				
				// Log CPU info if available
				if node.Cpus != nil {
					t.Logf("    CPUs: %d", *node.Cpus)
				}
				
				// Log memory if available
				if node.RealMemory != nil {
					t.Logf("    Memory: %d MB", *node.RealMemory)
				}
			}
		}
	})

	t.Run("Get Diagnostics", func(t *testing.T) {
		resp, err := apiClient.SlurmV0043GetDiagWithResponse(ctx)
		require.NoError(t, err)
		
		assert.Equal(t, 200, resp.StatusCode())
		
		if resp.JSON200 != nil {
			stats := &resp.JSON200.Statistics
			
			t.Logf("SLURM Diagnostics:")
			
			// Log server thread count
			if stats.ServerThreadCount != nil {
				t.Logf("  Server thread count: %d", *stats.ServerThreadCount)
			}
			
			// Log job stats
			if stats.JobsSubmitted != nil {
				t.Logf("  Jobs submitted: %d", *stats.JobsSubmitted)
			}
			if stats.JobsStarted != nil {
				t.Logf("  Jobs started: %d", *stats.JobsStarted)
			}
			if stats.JobsCompleted != nil {
				t.Logf("  Jobs completed: %d", *stats.JobsCompleted)
			}
			if stats.JobsCanceled != nil {
				t.Logf("  Jobs canceled: %d", *stats.JobsCanceled)
			}
			if stats.JobsFailed != nil {
				t.Logf("  Jobs failed: %d", *stats.JobsFailed)
			}
			
			// Log agent stats
			if stats.AgentCount != nil {
				t.Logf("  Agent count: %d", *stats.AgentCount)
			}
			if stats.AgentQueueSize != nil {
				t.Logf("  Agent queue size: %d", *stats.AgentQueueSize)
			}
		}
	})
}

