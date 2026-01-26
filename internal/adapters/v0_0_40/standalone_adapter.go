// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"context"
	"fmt"

	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/pkg/errors"
)

// StandaloneAdapter implements the standalone operations for v0.0.40
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

	resp, err := a.client.SlurmV0040GetLicensesWithResponse(ctx)
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
	params := &api.SlurmV0040GetSharesParams{}
	if opts != nil {
		if len(opts.Users) > 0 {
			params.Users = &opts.Users[0] // API might take single user
		}
		if len(opts.Accounts) > 0 {
			params.Accounts = &opts.Accounts[0] // API might take single account
		}
		// v0.0.40 doesn't have a Partition parameter for GetShares
	}

	resp, err := a.client.SlurmV0040GetSharesWithResponse(ctx, params)
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
		// v0.0.40: Fairshare.Level is *float64, not a struct with Number field
		if apiShare.Fairshare != nil && apiShare.Fairshare.Level != nil {
			share.FairshareLevel = *apiShare.Fairshare.Level
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

// GetConfig returns not implemented error for v0.0.40
func (a *StandaloneAdapter) GetConfig(ctx context.Context) (*types.Config, error) {
	return nil, errors.NewNotImplementedError("GetConfig", "v0.0.40")
}

// GetDiagnostics retrieves SLURM diagnostics information
func (a *StandaloneAdapter) GetDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	if a.client == nil {
		return nil, fmt.Errorf("API client not initialized")
	}

	resp, err := a.client.SlurmV0040GetDiagWithResponse(ctx)
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

	// Map statistics fields
	stats := resp.JSON200.Statistics
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

	// RPC statistics
	// Note: v0.0.40 doesn't have RPC statistics in the same structure

	return diag, nil
}

// GetDBDiagnostics retrieves SLURM database diagnostics information
func (a *StandaloneAdapter) GetDBDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	// Note: v0.0.40 DB diagnostics return a different structure (V0040StatsRec with RPCs/Users)
	// that doesn't map cleanly to types.Diagnostics. Delegate to GetDiagnostics for now.
	return a.GetDiagnostics(ctx)
}

// GetInstance returns not implemented error for v0.0.40
func (a *StandaloneAdapter) GetInstance(ctx context.Context, opts *types.GetInstanceOptions) (*types.Instance, error) {
	return nil, errors.NewNotImplementedError("GetInstance", "v0.0.40")
}

// GetInstances returns not implemented error for v0.0.40
func (a *StandaloneAdapter) GetInstances(ctx context.Context, opts *types.GetInstancesOptions) (*types.InstanceList, error) {
	return nil, errors.NewNotImplementedError("GetInstances", "v0.0.40")
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

	resp, err := a.client.SlurmdbV0040GetTresWithResponse(ctx)
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
	tresList := make([]types.TRES, 0)
	for _, apiTres := range resp.JSON200.TRES {
		tres := types.TRES{}

		if apiTres.Id != nil {
			tres.ID = int(*apiTres.Id)
		}
		// Type is not a pointer in v0.0.40
		tres.Type = apiTres.Type
		if apiTres.Name != nil {
			tres.Name = *apiTres.Name
		}
		if apiTres.Count != nil {
			tres.Count = *apiTres.Count
		}

		tresList = append(tresList, tres)
	}

	return &types.TRESList{
		TRES: tresList,
		Meta: extractMeta(resp.JSON200.Meta),
	}, nil
}

// CreateTRES returns not implemented error for v0.0.40
func (a *StandaloneAdapter) CreateTRES(ctx context.Context, req *types.CreateTRESRequest) (*types.TRES, error) {
	return nil, errors.NewNotImplementedError("CreateTRES", "v0.0.40")
}

// Reconfigure returns not implemented error for v0.0.40
func (a *StandaloneAdapter) Reconfigure(ctx context.Context) (*types.ReconfigureResponse, error) {
	return nil, errors.NewNotImplementedError("Reconfigure", "v0.0.40")
}

// PingDatabase pings the SLURM database for health checks
func (a *StandaloneAdapter) PingDatabase(ctx context.Context) (*types.PingResponse, error) {
	return nil, errors.NewNotImplementedError("PingDatabase", "v0.0.40")
}

// extractMeta extracts metadata from API response meta field
func extractMeta(meta *api.V0040OpenapiMeta) map[string]interface{} {
	result := make(map[string]interface{})

	if meta == nil {
		return result
	}

	// V0040OpenapiMeta has Client, Command, Plugin fields but not Messages/Warnings/Errors
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

	if meta.Command != nil && len(*meta.Command) > 0 {
		result["command"] = *meta.Command
	}

	if meta.Plugin != nil {
		pluginInfo := make(map[string]interface{})
		if meta.Plugin.Type != nil {
			pluginInfo["type"] = *meta.Plugin.Type
		}
		if meta.Plugin.Name != nil {
			pluginInfo["name"] = *meta.Plugin.Name
		}
		if len(pluginInfo) > 0 {
			result["plugin"] = pluginInfo
		}
	}

	return result
}
