// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_40

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	types "github.com/jontk/slurm-client/api"
	adapterbase "github.com/jontk/slurm-client/internal/adapters/base"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_40"
)

// InfoAdapter implements the InfoAdapter interface for v0.0.40
type InfoAdapter struct {
	*adapterbase.BaseManager
	client *api.ClientWithResponses
}

// NewInfoAdapter creates a new Info adapter for v0.0.40
func NewInfoAdapter(client *api.ClientWithResponses) *InfoAdapter {
	return &InfoAdapter{
		BaseManager: adapterbase.NewBaseManager("v0.0.40", "Info"),
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
	pingResp, err := a.client.SlurmV0040GetPingWithResponse(ctx)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}
	// Check response status
	if pingResp.StatusCode() != 200 {
		return nil, a.HandleAPIError(fmt.Errorf("API error: status %d", pingResp.StatusCode()))
	}
	// Check for unexpected response format
	if err := a.CheckNilResponse(pingResp.JSON200, "Get Cluster Info"); err != nil {
		return nil, err
	}
	// Extract cluster info from ping response
	clusterInfo := &types.ClusterInfo{
		APIVersion: "v0.0.40", // Current API version
	}
	if pingResp.JSON200.Meta != nil && pingResp.JSON200.Meta.Slurm != nil {
		if pingResp.JSON200.Meta.Slurm.Version != nil {
			if pingResp.JSON200.Meta.Slurm.Version.Major != nil &&
				pingResp.JSON200.Meta.Slurm.Version.Minor != nil &&
				pingResp.JSON200.Meta.Slurm.Version.Micro != nil {
				clusterInfo.Version = fmt.Sprintf("%s.%s.%s",
					*pingResp.JSON200.Meta.Slurm.Version.Major,
					*pingResp.JSON200.Meta.Slurm.Version.Minor,
					*pingResp.JSON200.Meta.Slurm.Version.Micro)
			}
		}
		if pingResp.JSON200.Meta.Slurm.Release != nil {
			clusterInfo.Release = *pingResp.JSON200.Meta.Slurm.Release
		}
		if pingResp.JSON200.Meta.Slurm.Cluster != nil {
			clusterInfo.ClusterName = *pingResp.JSON200.Meta.Slurm.Cluster
		}
	}
	// Try to get additional diagnostic information
	diagResp, diagErr := a.client.SlurmV0040GetDiagWithResponse(ctx)
	if diagErr == nil && diagResp.StatusCode() == 200 && diagResp.JSON200 != nil &&
		(diagResp.JSON200.Errors == nil || len(*diagResp.JSON200.Errors) == 0) {
		// Extract uptime from diagnostic statistics
		if diagResp.JSON200.Statistics.ServerThreadCount != nil {
			// Server thread count can be used as a proxy for activity/uptime
			clusterInfo.Uptime = int(*diagResp.JSON200.Statistics.ServerThreadCount)
		}
	}
	return clusterInfo, nil
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
	resp, err := a.client.SlurmV0040GetPingWithResponse(ctx)
	if err != nil {
		return a.HandleAPIError(err)
	}
	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
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
	// Call the database config endpoint (v0.0.40 feature) to test database connectivity
	resp, err := a.client.SlurmdbV0040GetConfigWithResponse(ctx)
	if err != nil {
		return a.HandleAPIError(err)
	}
	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
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
	resp, err := a.client.SlurmV0040GetDiagWithResponse(ctx)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}
	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}
	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "Get Cluster Stats"); err != nil {
		return nil, err
	}
	stats := &types.ClusterStats{}
	// Extract job statistics from diagnostic response
	a.extractJobStats(resp.JSON200, stats)
	// Get node statistics by querying the nodes endpoint
	nodesResp, err := a.client.SlurmV0040GetNodesWithResponse(ctx, nil)
	if err == nil && nodesResp.StatusCode() == 200 && nodesResp.JSON200 != nil {
		a.extractNodeStats(nodesResp.JSON200, stats)
	}
	return stats, nil
}

// extractJobStats extracts job-related statistics from diagnostic response
func (a *InfoAdapter) extractJobStats(diag *api.V0040OpenapiDiagResp, stats *types.ClusterStats) {
	if diag.Statistics.JobsSubmitted != nil {
		stats.TotalJobs = int(*diag.Statistics.JobsSubmitted)
	}
	if diag.Statistics.JobsPending != nil {
		stats.PendingJobs = int(*diag.Statistics.JobsPending)
	}
	if diag.Statistics.JobsRunning != nil {
		stats.RunningJobs = int(*diag.Statistics.JobsRunning)
	}
	if diag.Statistics.JobsCompleted != nil {
		stats.CompletedJobs = int(*diag.Statistics.JobsCompleted)
	}
}

// extractNodeStats extracts node-related statistics from nodes response
func (a *InfoAdapter) extractNodeStats(nodes *api.V0040OpenapiNodesResp, stats *types.ClusterStats) {
	for _, node := range nodes.Nodes {
		stats.TotalNodes++
		// Count total CPUs
		if node.Cpus != nil {
			stats.TotalCPUs += int(*node.Cpus)
		}
		// Update stats based on node state
		a.updateNodeStateStats(node, stats)
	}
	// Calculate idle CPUs if not fully allocated
	if stats.IdleCPUs == 0 && stats.TotalCPUs > 0 {
		stats.IdleCPUs = stats.TotalCPUs - stats.AllocatedCPUs
	}
}

// updateNodeStateStats updates stats based on node state
func (a *InfoAdapter) updateNodeStateStats(node api.V0040Node, stats *types.ClusterStats) {
	if node.State == nil || len(*node.State) == 0 {
		return
	}
	stateKeyword := strings.ToLower((*node.State)[0])
	if a.isNodeInState(stateKeyword, "idle") {
		stats.IdleNodes++
		if node.Cpus != nil {
			stats.IdleCPUs += int(*node.Cpus)
		}
		return
	}
	if a.isNodeAllocated(stateKeyword) {
		stats.AllocatedNodes++
		if node.Cpus != nil && strings.Contains(stateKeyword, "alloc") {
			stats.AllocatedCPUs += int(*node.Cpus)
		}
	}
}

// isNodeInState checks if the node state contains the specified keyword
func (a *InfoAdapter) isNodeInState(stateKeyword, keyword string) bool {
	return strings.Contains(stateKeyword, keyword)
}

// isNodeAllocated checks if the node is allocated or mixed
func (a *InfoAdapter) isNodeAllocated(stateKeyword string) bool {
	return strings.Contains(stateKeyword, "alloc") ||
		strings.Contains(stateKeyword, "mixed")
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
	resp, err := a.client.SlurmV0040GetPingWithResponse(ctx)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}
	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}
	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "Get API Version"); err != nil {
		return nil, err
	}
	apiVersion := &types.APIVersion{
		Version:     "v0.0.40",
		Description: "Slurm REST API v0.0.40",
		Deprecated:  false, // v0.0.40 is currently latest
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
