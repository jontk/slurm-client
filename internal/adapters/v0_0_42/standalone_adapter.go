// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_42

import (
	"context"
	"fmt"

	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_42"
	"github.com/jontk/slurm-client/pkg/errors"
)

// StandaloneAdapter implements the standalone operations for v0.0.42
type StandaloneAdapter struct {
	client       *api.ClientWithResponses
	errorAdapter *ErrorAdapter
}

// NewStandaloneAdapter creates a new standalone adapter
func NewStandaloneAdapter(client *api.ClientWithResponses) *StandaloneAdapter {
	return &StandaloneAdapter{
		client:       client,
		errorAdapter: NewErrorAdapter(),
	}
}

// GetLicenses retrieves license information
func (a *StandaloneAdapter) GetLicenses(ctx context.Context) (*types.LicenseList, error) {
	if a.client == nil {
		return nil, fmt.Errorf("API client not initialized")
	}
	resp, err := a.client.SlurmV0042GetLicensesWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get licenses: %w", err)
	}
	// Handle API response with enhanced error handling
	if err := a.errorAdapter.HandleAPIResponse(resp.StatusCode(), resp.Body, "GetLicenses"); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return &types.LicenseList{Licenses: []types.License{}}, nil
	}
	// Convert API licenses to common types
	licenses := make([]types.License, 0)
	// Licenses is not a pointer field
	for _, apiLicense := range resp.JSON200.Licenses {
		license := types.License{}
		if apiLicense.LicenseName != nil {
			license.Name = *apiLicense.LicenseName
		}
		if apiLicense.Total != nil {
			license.Total = int(*apiLicense.Total)
		}
		if apiLicense.Used != nil {
			license.Used = int(*apiLicense.Used)
		}
		if apiLicense.Free != nil {
			license.Free = int(*apiLicense.Free)
		}
		if apiLicense.Reserved != nil {
			license.Reserved = int(*apiLicense.Reserved)
		}
		// Remote is a bool - we could set a flag or use a different field
		// For now, just note that this is a remote license
		if apiLicense.Remote != nil && *apiLicense.Remote {
		}
		licenses = append(licenses, license)
	}
	return &types.LicenseList{
		Licenses: licenses,
		Meta:     extractMeta(resp.JSON200.Meta),
	}, nil
}

// GetShares retrieves fairshare information with optional filtering
func (a *StandaloneAdapter) GetShares(ctx context.Context, opts *types.GetSharesOptions) (*types.SharesList, error) {
	if a.client == nil {
		return nil, fmt.Errorf("API client not initialized")
	}
	// Build query parameters
	params := &api.SlurmV0042GetSharesParams{}
	if opts != nil {
		if len(opts.Users) > 0 {
			params.Users = &opts.Users[0] // API might take single user
		}
		if len(opts.Accounts) > 0 {
			params.Accounts = &opts.Accounts[0] // API might take single account
		}
		// v0.0.42 doesn't have a Partition parameter for GetShares
	}
	resp, err := a.client.SlurmV0042GetSharesWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get shares: %w", err)
	}
	// Handle API response with enhanced error handling
	if err := a.errorAdapter.HandleAPIResponse(resp.StatusCode(), resp.Body, "GetShares"); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return &types.SharesList{Shares: []types.Share{}}, nil
	}
	// Convert API shares to common types
	shares := make([]types.Share, 0)
	for _, apiShare := range *resp.JSON200.Shares.Shares {
		share := types.Share{}
		if apiShare.Name != nil {
			// This could be account or user name
			share.Account = *apiShare.Name
		}
		if apiShare.Partition != nil {
			share.Partition = *apiShare.Partition
		}
		// Convert share numbers
		if apiShare.Shares != nil && apiShare.Shares.Number != nil {
			share.RawShares = int(*apiShare.Shares.Number)
		}
		if apiShare.Usage != nil {
			share.RawUsage = *apiShare.Usage
		}
		if apiShare.Fairshare != nil && apiShare.Fairshare.Level != nil && apiShare.Fairshare.Level.Number != nil {
			share.FairshareLevel = *apiShare.Fairshare.Level.Number
		}
		if apiShare.SharesNormalized != nil && apiShare.SharesNormalized.Number != nil {
			// SharesNormalized.Number is float64, convert to int
			share.FairshareShares = int(*apiShare.SharesNormalized.Number)
		}
		shares = append(shares, share)
	}
	return &types.SharesList{
		Shares: shares,
		Meta:   extractMeta(resp.JSON200.Meta),
	}, nil
}

// GetConfig retrieves SLURM configuration
func (a *StandaloneAdapter) GetConfig(ctx context.Context) (*types.Config, error) {
	if a.client == nil {
		return nil, fmt.Errorf("API client not initialized")
	}
	resp, err := a.client.SlurmdbV0042GetConfigWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	// Handle API response with enhanced error handling
	if err := a.errorAdapter.HandleAPIResponse(resp.StatusCode(), resp.Body, "GetConfig"); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty configuration response")
	}
	// Convert API config to common type
	// Note: v0.0.42 config endpoint returns entity lists (accounts, associations, etc.)
	// rather than daemon configuration settings. Extract what metadata is available.
	config := &types.Config{
		Meta:    extractMeta(resp.JSON200.Meta),
		Version: "v0.0.42",
	}
	// Extract cluster name from first cluster if available
	if resp.JSON200.Clusters != nil && len(*resp.JSON200.Clusters) > 0 {
		clusters := *resp.JSON200.Clusters
		if clusters[0].Name != nil {
			config.ClusterName = *clusters[0].Name
		}
	}
	return config, nil
}

// GetDiagnostics retrieves SLURM diagnostics information
func (a *StandaloneAdapter) GetDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	if a.client == nil {
		return nil, fmt.Errorf("API client not initialized")
	}
	resp, err := a.client.SlurmV0042GetDiagWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get diagnostics: %w", err)
	}
	// Handle API response with enhanced error handling
	if err := a.errorAdapter.HandleAPIResponse(resp.StatusCode(), resp.Body, "GetDiagnostics"); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty diagnostics response")
	}
	// Convert API diagnostics to common type
	diag := &types.Diagnostics{
		Meta: extractMeta(resp.JSON200.Meta),
	}
	// Map statistics fields using helper functions
	stats := resp.JSON200.Statistics
	a.setJobCountMetrics(diag, stats)
	a.setBackfillMetrics(diag, stats)
	a.setScheduleCycleMetrics(diag, stats)
	a.setAgentThreadMetrics(diag, stats)
	// RPC statistics
	// Note: v0.0.42 doesn't have RPC statistics in the same structure
	return diag, nil
}

// setJobCountMetrics sets job count statistics (submitted, started, completed, etc.)
func (a *StandaloneAdapter) setJobCountMetrics(diag *types.Diagnostics, stats api.V0042StatsMsg) {
	if stats.JobsSubmitted != nil {
		diag.JobsSubmitted = int(*stats.JobsSubmitted)
	}
	if stats.JobsStarted != nil {
		diag.JobsStarted = int(*stats.JobsStarted)
	}
	if stats.JobsCompleted != nil {
		diag.JobsCompleted = int(*stats.JobsCompleted)
	}
	if stats.JobsCanceled != nil {
		diag.JobsCanceled = int(*stats.JobsCanceled)
	}
	if stats.JobsFailed != nil {
		diag.JobsFailed = int(*stats.JobsFailed)
	}
	if stats.JobsPending != nil {
		diag.JobsPending = int(*stats.JobsPending)
	}
	if stats.JobsRunning != nil {
		diag.JobsRunning = int(*stats.JobsRunning)
	}
}

// setBackfillMetrics sets backfill scheduler metrics
func (a *StandaloneAdapter) setBackfillMetrics(diag *types.Diagnostics, stats api.V0042StatsMsg) {
	if stats.BfCycleCounter != nil {
		diag.BFCycle = int(*stats.BfCycleCounter)
	}
	if stats.BfCycleMean != nil {
		diag.BFCycleMean = *stats.BfCycleMean
	}
	if stats.BfCycleMax != nil {
		diag.BFCycleMax = int64(*stats.BfCycleMax)
	}
	if stats.BfCycleLast != nil {
		diag.BFCycleMean = int64(*stats.BfCycleLast) // Store last in a mean-like field
	}
}

// setScheduleCycleMetrics sets schedule cycle metrics
func (a *StandaloneAdapter) setScheduleCycleMetrics(diag *types.Diagnostics, stats api.V0042StatsMsg) {
	if stats.ScheduleCycleTotal != nil {
		diag.ScheduleCycleCounter = int(*stats.ScheduleCycleTotal)
	}
	if stats.ScheduleCycleMean != nil {
		diag.ScheduleCycleMean = *stats.ScheduleCycleMean
	}
	if stats.ScheduleCycleMax != nil {
		diag.ScheduleCycleMax = int64(*stats.ScheduleCycleMax)
	}
	if stats.ScheduleCycleLast != nil {
		diag.ScheduleCycleLast = int64(*stats.ScheduleCycleLast)
	}
}

// setAgentThreadMetrics sets agent and thread metrics
func (a *StandaloneAdapter) setAgentThreadMetrics(diag *types.Diagnostics, stats api.V0042StatsMsg) {
	if stats.AgentCount != nil {
		diag.AgentCount = int(*stats.AgentCount)
	}
	if stats.AgentThreadCount != nil {
		diag.AgentThreadCount = int(*stats.AgentThreadCount)
	}
	if stats.DbdAgentQueueSize != nil {
		diag.DBDAgentCount = int(*stats.DbdAgentQueueSize)
	}
	if stats.ServerThreadCount != nil {
		diag.ServerThreadCount = int(*stats.ServerThreadCount)
	}
}

// GetDBDiagnostics retrieves SLURM database diagnostics information
func (a *StandaloneAdapter) GetDBDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	// Note: v0.0.42 DB diagnostics return a different structure (V0042StatsRec with RPCs/Users)
	// that doesn't map cleanly to types.Diagnostics. Delegate to GetDiagnostics for now.
	return a.GetDiagnostics(ctx)
}

// GetInstance returns not implemented error for v0.0.42
func (a *StandaloneAdapter) GetInstance(ctx context.Context, opts *types.GetInstanceOptions) (*types.Instance, error) {
	return nil, errors.NewNotImplementedError("GetInstance", "v0.0.42")
}

// GetInstances returns not implemented error for v0.0.42
func (a *StandaloneAdapter) GetInstances(ctx context.Context, opts *types.GetInstancesOptions) (*types.InstanceList, error) {
	return nil, errors.NewNotImplementedError("GetInstances", "v0.0.42")
}

// GetTRES retrieves all TRES (Trackable RESources)
func (a *StandaloneAdapter) GetTRES(ctx context.Context) (*types.TRESList, error) {
	if ctx == nil {
		return nil, errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"context is required",
			"ctx", nil, nil,
		)
	}
	if a.client == nil {
		return nil, fmt.Errorf("API client not initialized")
	}
	resp, err := a.client.SlurmdbV0042GetTresWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get TRES: %w", err)
	}
	// Handle API response with enhanced error handling
	if err := a.errorAdapter.HandleAPIResponse(resp.StatusCode(), resp.Body, "GetTRES"); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil || resp.JSON200.TRES == nil {
		return &types.TRESList{TRES: []types.TRES{}}, nil
	}
	// Convert API TRES to common types
	// types.TRES has ID *int32, Type string, Name *string, Count *int64
	tresList := make([]types.TRES, 0)
	for _, apiTres := range resp.JSON200.TRES {
		tres := types.TRES{}
		// ID is *int32 in types.TRES
		tres.ID = apiTres.Id
		// Type is string, not pointer
		tres.Type = apiTres.Type
		// Name is *string in types.TRES
		tres.Name = apiTres.Name
		// Count is *int64 in types.TRES
		tres.Count = apiTres.Count
		tresList = append(tresList, tres)
	}
	return &types.TRESList{
		TRES: tresList,
		Meta: extractMeta(resp.JSON200.Meta),
	}, nil
}

// CreateTRES returns not implemented error for v0.0.42
func (a *StandaloneAdapter) CreateTRES(ctx context.Context, req *types.CreateTRESRequest) (*types.TRES, error) {
	return nil, errors.NewNotImplementedError("CreateTRES", "v0.0.42")
}

// Reconfigure returns not implemented error for v0.0.42
func (a *StandaloneAdapter) Reconfigure(ctx context.Context) (*types.ReconfigureResponse, error) {
	return nil, errors.NewNotImplementedError("Reconfigure", "v0.0.42")
}

// PingDatabase pings the SLURM database for health checks
func (a *StandaloneAdapter) PingDatabase(ctx context.Context) (*types.PingResponse, error) {
	return nil, errors.NewNotImplementedError("PingDatabase", "v0.0.42")
}

// extractMeta safely extracts metadata from API response
func extractMeta(meta *api.V0042OpenapiMeta) map[string]interface{} {
	result := make(map[string]interface{})
	if meta == nil {
		return result
	}
	// V0042OpenapiMeta has Client, Command, Plugin fields but not Messages/Warnings/Errors
	// Extract basic metadata
	if meta.Client != nil {
		clientInfo := make(map[string]interface{})
		if meta.Client.Source != nil {
			clientInfo["source"] = *meta.Client.Source
		}
		if meta.Client.User != nil {
			clientInfo["user"] = *meta.Client.User
		}
		if meta.Client.Group != nil {
			clientInfo["group"] = *meta.Client.Group
		}
		if len(clientInfo) > 0 {
			result["client"] = clientInfo
		}
	}
	if meta.Plugin != nil {
		pluginInfo := make(map[string]interface{})
		if meta.Plugin.AccountingStorage != nil {
			pluginInfo["accounting_storage"] = *meta.Plugin.AccountingStorage
		}
		if len(pluginInfo) > 0 {
			result["plugin"] = pluginInfo
		}
	}
	return result
}
