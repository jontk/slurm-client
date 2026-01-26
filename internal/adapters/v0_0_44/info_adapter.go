// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_44

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	api "github.com/jontk/slurm-client/internal/api/v0_0_44"
	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
)

// InfoAdapter implements the InfoAdapter interface for v0.0.44
type InfoAdapter struct {
	*base.BaseManager
	client *api.ClientWithResponses
}

// NewInfoAdapter creates a new Info adapter for v0.0.44
func NewInfoAdapter(client *api.ClientWithResponses) *InfoAdapter {
	return &InfoAdapter{
		BaseManager: base.NewBaseManager("v0.0.44", "Info"),
		client:      client,
	}
}

// Get retrieves cluster information
func (a *InfoAdapter) Get(ctx context.Context) (*types.ClusterInfo, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Call the ping endpoint to get basic cluster information
	pingResp, err := a.client.SlurmV0044GetPingWithResponse(ctx)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0044OpenapiErrors
	if pingResp.JSON200 != nil {
		apiErrors = pingResp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(pingResp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.44"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(pingResp.JSON200, "Get Cluster Info"); err != nil {
		return nil, err
	}

	// Extract cluster info from ping response
	clusterInfo := &types.ClusterInfo{
		APIVersion: "v0.0.44", // Current API version
	}

	// Extract cluster metadata from ping response
	a.extractClusterMetadata(pingResp.JSON200, clusterInfo)

	// Try to get additional diagnostic information
	a.addDiagnosticInfo(ctx, clusterInfo)

	return clusterInfo, nil
}

// extractClusterMetadata extracts metadata from ping response
func (a *InfoAdapter) extractClusterMetadata(pingResp *api.V0044OpenapiPingArrayResp, clusterInfo *types.ClusterInfo) {
	if pingResp.Meta != nil && pingResp.Meta.Slurm != nil {
		if pingResp.Meta.Slurm.Version != nil {
			if pingResp.Meta.Slurm.Version.Major != nil &&
				pingResp.Meta.Slurm.Version.Minor != nil &&
				pingResp.Meta.Slurm.Version.Micro != nil {
				clusterInfo.Version = fmt.Sprintf("%s.%s.%s",
					*pingResp.Meta.Slurm.Version.Major,
					*pingResp.Meta.Slurm.Version.Minor,
					*pingResp.Meta.Slurm.Version.Micro)
			}
		}

		if pingResp.Meta.Slurm.Release != nil {
			clusterInfo.Release = *pingResp.Meta.Slurm.Release
		}

		if pingResp.Meta.Slurm.Cluster != nil {
			clusterInfo.ClusterName = *pingResp.Meta.Slurm.Cluster
		}
	}
}

// addDiagnosticInfo adds diagnostic information to cluster info
func (a *InfoAdapter) addDiagnosticInfo(ctx context.Context, clusterInfo *types.ClusterInfo) {
	diagResp, diagErr := a.client.SlurmV0044GetDiagWithResponse(ctx)
	if diagErr == nil && diagResp.StatusCode() == 200 && diagResp.JSON200 != nil &&
		(diagResp.JSON200.Errors == nil || len(*diagResp.JSON200.Errors) == 0) {
		// Extract uptime from diagnostic statistics
		if diagResp.JSON200.Statistics.ServerThreadCount != nil {
			// Server thread count can be used as a proxy for activity/uptime
			clusterInfo.Uptime = int(*diagResp.JSON200.Statistics.ServerThreadCount)
		}
	}
}

// Ping tests connectivity to the cluster
func (a *InfoAdapter) Ping(ctx context.Context) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the ping endpoint
	resp, err := a.client.SlurmV0044GetPingWithResponse(ctx)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0044OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.44"); err != nil {
		return err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "Ping Cluster"); err != nil {
		return err
	}

	return nil
}

// PingDatabase tests connectivity to the SLURM database
func (a *InfoAdapter) PingDatabase(ctx context.Context) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the database config endpoint (v0.0.44 feature) to test database connectivity
	resp, err := a.client.SlurmdbV0044GetConfigWithResponse(ctx)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0044OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.44"); err != nil {
		return err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "Ping Database"); err != nil {
		return err
	}

	return nil
}

// Stats retrieves cluster statistics
func (a *InfoAdapter) Stats(ctx context.Context) (*types.ClusterStats, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Call the diagnostic endpoint to get cluster statistics
	resp, err := a.client.SlurmV0044GetDiagWithResponse(ctx)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0044OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.44"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "Get Cluster Stats"); err != nil {
		return nil, err
	}

	stats := &types.ClusterStats{}

	// Extract job statistics from diagnostic response
	a.extractJobStats(resp.JSON200, stats)

	// Get node statistics by querying the nodes endpoint
	a.enrichStatsWithNodeData(ctx, stats)

	return stats, nil
}

// extractJobStats extracts job statistics from diagnostic response
func (a *InfoAdapter) extractJobStats(diagResp *api.V0044OpenapiDiagResp, stats *types.ClusterStats) {
	if diagResp.Statistics.JobsSubmitted != nil {
		stats.TotalJobs = int(*diagResp.Statistics.JobsSubmitted)
	}

	if diagResp.Statistics.JobsPending != nil {
		stats.PendingJobs = int(*diagResp.Statistics.JobsPending)
	}

	if diagResp.Statistics.JobsRunning != nil {
		stats.RunningJobs = int(*diagResp.Statistics.JobsRunning)
	}

	if diagResp.Statistics.JobsCompleted != nil {
		stats.CompletedJobs = int(*diagResp.Statistics.JobsCompleted)
	}
}

// enrichStatsWithNodeData adds node statistics to cluster stats
func (a *InfoAdapter) enrichStatsWithNodeData(ctx context.Context, stats *types.ClusterStats) {
	nodesResp, err := a.client.SlurmV0044GetNodesWithResponse(ctx, nil)
	if err == nil && nodesResp.StatusCode() == 200 && nodesResp.JSON200 != nil {
		// Count nodes and CPUs by state
		for _, node := range nodesResp.JSON200.Nodes {
			a.processNodeForStats(node, stats)
		}

		// Calculate idle CPUs if not fully allocated
		if stats.IdleCPUs == 0 && stats.TotalCPUs > 0 {
			stats.IdleCPUs = stats.TotalCPUs - stats.AllocatedCPUs
		}
	}
}

// processNodeForStats processes a single node and updates stats
func (a *InfoAdapter) processNodeForStats(node api.V0044Node, stats *types.ClusterStats) {
	stats.TotalNodes++

	// Count CPUs
	if node.Cpus != nil {
		stats.TotalCPUs += int(*node.Cpus)
	}

	// Check node state and update statistics
	a.updateStatsForNodeState(node, stats)
}

// updateStatsForNodeState updates statistics based on node state
func (a *InfoAdapter) updateStatsForNodeState(node api.V0044Node, stats *types.ClusterStats) {
	if node.State != nil && len(*node.State) > 0 {
		state := string((*node.State)[0])
		stateLower := strings.ToLower(state)

		switch {
		case strings.Contains(stateLower, "idle"):
			stats.IdleNodes++
			if node.Cpus != nil {
				stats.IdleCPUs += int(*node.Cpus)
			}
		case strings.Contains(stateLower, "alloc") || strings.Contains(stateLower, "mixed"):
			stats.AllocatedNodes++
			// For allocated/mixed nodes, we'd need more info to get exact CPU allocation
			// This is a simplified approach
			if node.Cpus != nil && strings.Contains(stateLower, "alloc") {
				stats.AllocatedCPUs += int(*node.Cpus)
			}
		}
	}
}

// Version retrieves API version information
func (a *InfoAdapter) Version(ctx context.Context) (*types.APIVersion, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Call the ping endpoint to get version information
	resp, err := a.client.SlurmV0044GetPingWithResponse(ctx)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0044OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.44"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "Get API Version"); err != nil {
		return nil, err
	}

	apiVersion := &types.APIVersion{
		Version:     "v0.0.44",
		Description: "Slurm REST API v0.0.44",
		Deprecated:  false, // v0.0.44 is currently latest
	}

	// Extract Slurm version from ping response
	if resp.JSON200.Meta != nil && resp.JSON200.Meta.Slurm != nil {
		if resp.JSON200.Meta.Slurm.Version != nil {
			if resp.JSON200.Meta.Slurm.Version.Major != nil &&
				resp.JSON200.Meta.Slurm.Version.Minor != nil &&
				resp.JSON200.Meta.Slurm.Version.Micro != nil {
				apiVersion.Release = fmt.Sprintf("%s.%s.%s",
					*resp.JSON200.Meta.Slurm.Version.Major,
					*resp.JSON200.Meta.Slurm.Version.Minor,
					*resp.JSON200.Meta.Slurm.Version.Micro)
			}
		}

		if resp.JSON200.Meta.Slurm.Release != nil {
			// Extract more detailed release info if available
			release := *resp.JSON200.Meta.Slurm.Release
			if release != "" {
				apiVersion.Release = release

				// Check if this is a pre-release or development version
				if matched, _ := regexp.MatchString(`(alpha|beta|rc|dev)`, release); matched {
					apiVersion.Description += " (Development/Pre-release)"
				}
			}
		}
	}

	return apiVersion, nil
}
