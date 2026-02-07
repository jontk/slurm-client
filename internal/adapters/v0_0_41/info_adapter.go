// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	types "github.com/jontk/slurm-client/api"
	adapterbase "github.com/jontk/slurm-client/internal/adapters/base"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
)

// InfoAdapter implements the InfoAdapter interface for v0.0.41
type InfoAdapter struct {
	*adapterbase.BaseManager
	client *api.ClientWithResponses
}

// NewInfoAdapter creates a new Info adapter for v0.0.41
func NewInfoAdapter(client *api.ClientWithResponses) *InfoAdapter {
	return &InfoAdapter{
		BaseManager: adapterbase.NewBaseManager("v0.0.41", "Info"),
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
	pingResp, err := a.client.SlurmV0041GetPingWithResponse(ctx)
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
		APIVersion: "v0.0.41", // Current API version
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
	diagResp, diagErr := a.client.SlurmV0041GetDiagWithResponse(ctx)
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
	resp, err := a.client.SlurmV0041GetPingWithResponse(ctx)
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
	// Call the database config endpoint (v0.0.41 feature) to test database connectivity
	resp, err := a.client.SlurmdbV0041GetConfigWithResponse(ctx)
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
	resp, err := a.client.SlurmV0041GetDiagWithResponse(ctx)
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
	a.extractJobStats(resp.JSON200.Statistics, stats)
	// Get and process node statistics
	a.processNodeStats(ctx, stats)
	return stats, nil
}

// extractJobStats extracts job statistics from the statistics response
func (a *InfoAdapter) extractJobStats(statsField interface{}, stats *types.ClusterStats) {
	if statsField == nil {
		return
	}
	// Use reflection to safely access fields on the anonymous struct
	v := reflect.ValueOf(statsField)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}
	// Extract job counts using helper function
	a.extractJobCount(v, "JobsSubmitted", &stats.TotalJobs)
	a.extractJobCount(v, "JobsPending", &stats.PendingJobs)
	a.extractJobCount(v, "JobsRunning", &stats.RunningJobs)
	a.extractJobCount(v, "JobsCompleted", &stats.CompletedJobs)
}

// extractJobCount extracts a single job count from a stats field
func (a *InfoAdapter) extractJobCount(v reflect.Value, fieldName string, target *int) {
	if field := v.FieldByName(fieldName); field.IsValid() && !field.IsNil() {
		if val, ok := field.Interface().(*int32); ok && val != nil {
			*target = int(*val)
		}
	}
}

// processNodeStats retrieves and processes node statistics from the nodes endpoint
func (a *InfoAdapter) processNodeStats(ctx context.Context, stats *types.ClusterStats) {
	nodesResp, err := a.client.SlurmV0041GetNodesWithResponse(ctx, nil)
	if err != nil || nodesResp.StatusCode() != 200 || nodesResp.JSON200 == nil {
		return
	}
	for _, node := range nodesResp.JSON200.Nodes {
		a.updateNodeStats(node, stats)
	}
	// Calculate idle CPUs if not fully allocated
	if stats.IdleCPUs == 0 && stats.TotalCPUs > 0 {
		stats.IdleCPUs = stats.TotalCPUs - stats.AllocatedCPUs
	}
}

// updateNodeStats updates statistics for a single node
func (a *InfoAdapter) updateNodeStats(node interface{}, stats *types.ClusterStats) {
	stats.TotalNodes++
	// Get CPUs from node
	v := reflect.ValueOf(node)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	cpus := a.getNodeCPUs(v)
	if cpus > 0 {
		stats.TotalCPUs += cpus
	}
	// Check node state
	state := a.getNodeState(v)
	if state != "" {
		a.processNodeState(state, cpus, stats)
	}
}

// getNodeCPUs extracts CPU count from a node
func (a *InfoAdapter) getNodeCPUs(nodeVal reflect.Value) int {
	if field := nodeVal.FieldByName("Cpus"); field.IsValid() && !field.IsNil() {
		if val, ok := field.Interface().(*int32); ok && val != nil {
			return int(*val)
		}
	}
	return 0
}

// getNodeState extracts state from a node
func (a *InfoAdapter) getNodeState(nodeVal reflect.Value) string {
	if field := nodeVal.FieldByName("State"); field.IsValid() && !field.IsNil() {
		if statePtr, ok := field.Interface().(*[]string); ok && statePtr != nil && len(*statePtr) > 0 {
			return (*statePtr)[0]
		}
	}
	return ""
}

// processNodeState processes the state of a single node
func (a *InfoAdapter) processNodeState(state string, cpus int, stats *types.ClusterStats) {
	stateLower := strings.ToLower(state)
	if strings.Contains(stateLower, "idle") {
		stats.IdleNodes++
		if cpus > 0 {
			stats.IdleCPUs += cpus
		}
		return
	}
	if strings.Contains(stateLower, "alloc") || strings.Contains(stateLower, "mixed") {
		stats.AllocatedNodes++
		// For allocated/mixed nodes, we'd need more info to get exact CPU allocation
		// This is a simplified approach
		if strings.Contains(stateLower, "alloc") && cpus > 0 {
			stats.AllocatedCPUs += cpus
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
	resp, err := a.client.SlurmV0041GetPingWithResponse(ctx)
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
		Version:     "v0.0.41",
		Description: "Slurm REST API v0.0.41",
		Deprecated:  false, // v0.0.41 is currently latest
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
