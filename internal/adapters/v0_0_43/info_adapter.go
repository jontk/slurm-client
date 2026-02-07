// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_43

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	types "github.com/jontk/slurm-client/api"
	adapterbase "github.com/jontk/slurm-client/internal/adapters/base"
	"github.com/jontk/slurm-client/internal/common"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_43"
)

// InfoAdapter implements the InfoAdapter interface for v0.0.43
type InfoAdapter struct {
	*adapterbase.BaseManager
	client *api.ClientWithResponses
}

// NewInfoAdapter creates a new Info adapter for v0.0.43
func NewInfoAdapter(client *api.ClientWithResponses) *InfoAdapter {
	return &InfoAdapter{
		BaseManager: adapterbase.NewBaseManager("v0.0.43", "Info"),
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
	pingResp, err := a.client.SlurmV0043GetPingWithResponse(ctx)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}
	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if pingResp.JSON200 != nil {
		apiErrors = pingResp.JSON200.Errors
	}
	responseAdapter := api.NewResponseAdapter(pingResp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
		return nil, err
	}
	// Check for unexpected response format
	if err := a.CheckNilResponse(pingResp.JSON200, "Get Cluster Info"); err != nil {
		return nil, err
	}
	// Extract cluster info from ping response
	clusterInfo := a.extractClusterInfoFromPing(pingResp.JSON200)
	// Try to get additional diagnostic information
	a.enrichClusterInfoWithDiagnostics(ctx, clusterInfo)
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
	resp, err := a.client.SlurmV0043GetPingWithResponse(ctx)
	if err != nil {
		return a.HandleAPIError(err)
	}
	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
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
	// Call the database config endpoint (v0.0.43 feature) to test database connectivity
	resp, err := a.client.SlurmdbV0043GetConfigWithResponse(ctx)
	if err != nil {
		return a.HandleAPIError(err)
	}
	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
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
	resp, err := a.client.SlurmV0043GetDiagWithResponse(ctx)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}
	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
		return nil, err
	}
	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "Get Cluster Stats"); err != nil {
		return nil, err
	}
	stats := &types.ClusterStats{}
	// Extract job statistics from diagnostic response
	a.extractJobStats(resp.JSON200.Statistics, stats)
	// Get node statistics by querying the nodes endpoint
	a.enrichStatsWithNodeData(ctx, stats)
	return stats, nil
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
	resp, err := a.client.SlurmV0043GetPingWithResponse(ctx)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}
	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
		return nil, err
	}
	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "Get API Version"); err != nil {
		return nil, err
	}
	apiVersion := &types.APIVersion{
		Version:     "v0.0.43",
		Description: "Slurm REST API v0.0.43",
		Deprecated:  false, // v0.0.43 is currently latest
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

// extractClusterInfoFromPing extracts cluster information from a ping response
func (a *InfoAdapter) extractClusterInfoFromPing(pingResp *api.V0043OpenapiPingArrayResp) *types.ClusterInfo {
	clusterInfo := &types.ClusterInfo{
		APIVersion: "v0.0.43", // Current API version
	}
	if pingResp.Meta == nil || pingResp.Meta.Slurm == nil {
		return clusterInfo
	}
	slurm := pingResp.Meta.Slurm
	// Extract version information
	if slurm.Version != nil {
		if slurm.Version.Major != nil && slurm.Version.Minor != nil && slurm.Version.Micro != nil {
			clusterInfo.Version = fmt.Sprintf("%s.%s.%s",
				*slurm.Version.Major,
				*slurm.Version.Minor,
				*slurm.Version.Micro)
		}
	}
	// Extract release information
	if slurm.Release != nil {
		clusterInfo.Release = *slurm.Release
	}
	// Extract cluster name
	if slurm.Cluster != nil {
		clusterInfo.ClusterName = *slurm.Cluster
	}
	return clusterInfo
}

// enrichClusterInfoWithDiagnostics adds diagnostic information to cluster info
func (a *InfoAdapter) enrichClusterInfoWithDiagnostics(ctx context.Context, clusterInfo *types.ClusterInfo) {
	diagResp, diagErr := a.client.SlurmV0043GetDiagWithResponse(ctx)
	if diagErr != nil {
		return
	}
	if diagResp.StatusCode() != 200 || diagResp.JSON200 == nil {
		return
	}
	// Check for errors in the response
	if diagResp.JSON200.Errors != nil && len(*diagResp.JSON200.Errors) > 0 {
		return
	}
	// Extract uptime from diagnostic statistics
	if diagResp.JSON200.Statistics.ServerThreadCount != nil {
		// Server thread count can be used as a proxy for activity/uptime
		clusterInfo.Uptime = int(*diagResp.JSON200.Statistics.ServerThreadCount)
	}
}

// extractJobStats extracts job statistics from diagnostic statistics
func (a *InfoAdapter) extractJobStats(diagStats api.V0043StatsMsg, stats *types.ClusterStats) {
	if diagStats.JobsSubmitted != nil {
		stats.TotalJobs = int(*diagStats.JobsSubmitted)
	}
	if diagStats.JobsPending != nil {
		stats.PendingJobs = int(*diagStats.JobsPending)
	}
	if diagStats.JobsRunning != nil {
		stats.RunningJobs = int(*diagStats.JobsRunning)
	}
	if diagStats.JobsCompleted != nil {
		stats.CompletedJobs = int(*diagStats.JobsCompleted)
	}
}

// enrichStatsWithNodeData queries and processes node statistics to enrich cluster stats
func (a *InfoAdapter) enrichStatsWithNodeData(ctx context.Context, stats *types.ClusterStats) {
	nodesResp, err := a.client.SlurmV0043GetNodesWithResponse(ctx, nil)
	if err != nil {
		return
	}
	if nodesResp.StatusCode() != 200 || nodesResp.JSON200 == nil {
		return
	}
	// Count nodes and CPUs by state
	for _, node := range nodesResp.JSON200.Nodes {
		a.processNodeForStats(node, stats)
	}
	// Calculate idle CPUs if not fully allocated
	a.calculateIdleCPUs(stats)
}

// processNodeForStats updates stats based on a single node's information
func (a *InfoAdapter) processNodeForStats(node api.V0043Node, stats *types.ClusterStats) {
	stats.TotalNodes++
	// Count CPUs
	if node.Cpus != nil {
		stats.TotalCPUs += int(*node.Cpus)
	}
	// Check node state and update stats accordingly
	if node.State != nil && len(*node.State) > 0 {
		state := string((*node.State)[0])
		a.updateStatsBasedOnNodeState(state, node, stats)
	}
}

// updateStatsBasedOnNodeState categorizes node and updates stats based on its state
func (a *InfoAdapter) updateStatsBasedOnNodeState(state string, node api.V0043Node, stats *types.ClusterStats) {
	stateLower := strings.ToLower(state)
	if strings.Contains(stateLower, "idle") {
		stats.IdleNodes++
		if node.Cpus != nil {
			stats.IdleCPUs += int(*node.Cpus)
		}
		return
	}
	if strings.Contains(stateLower, "alloc") || strings.Contains(stateLower, "mixed") {
		stats.AllocatedNodes++
		// For allocated/mixed nodes, we'd need more info to get exact CPU allocation
		// This is a simplified approach
		if node.Cpus != nil && strings.Contains(stateLower, "alloc") {
			stats.AllocatedCPUs += int(*node.Cpus)
		}
	}
}

// calculateIdleCPUs calculates idle CPU count when not fully allocated
func (a *InfoAdapter) calculateIdleCPUs(stats *types.ClusterStats) {
	if stats.IdleCPUs == 0 && stats.TotalCPUs > 0 {
		stats.IdleCPUs = stats.TotalCPUs - stats.AllocatedCPUs
	}
}
