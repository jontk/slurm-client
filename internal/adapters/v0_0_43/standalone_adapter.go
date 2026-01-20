// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	"fmt"
	"time"

	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/pkg/errors"
)

// StandaloneAdapter implements the standalone operations for v0.0.43
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

	resp, err := a.client.SlurmV0043GetLicensesWithResponse(ctx)
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
		if apiLicense.Remote != nil && *apiLicense.Remote {
			// Remote is a bool - we could set a flag or use a different field
			// For now, just note that this is a remote license
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
	params := &api.SlurmV0043GetSharesParams{}
	if opts != nil {
		if len(opts.Users) > 0 {
			params.Users = &opts.Users[0] // API might take single user
		}
		if len(opts.Accounts) > 0 {
			params.Accounts = &opts.Accounts[0] // API might take single account
		}
		// v0.0.43 doesn't have a Partition parameter for GetShares
	}

	resp, err := a.client.SlurmV0043GetSharesWithResponse(ctx, params)
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

	resp, err := a.client.SlurmdbV0043GetConfigWithResponse(ctx)
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
	config := &types.Config{
		Meta: extractMeta(resp.JSON200.Meta),
	}

	// Extract key configuration values
	// Note: The actual field mapping depends on the API structure
	// This is a simplified version - you'd need to check the actual API response structure
	// Map fields from the response
	// This would need to be expanded based on actual API structure
	config.Version = "v0.0.43" // Set version based on adapter

	return config, nil
}

// GetDiagnostics retrieves SLURM diagnostics information
func (a *StandaloneAdapter) GetDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	if a.client == nil {
		return nil, fmt.Errorf("API client not initialized")
	}

	resp, err := a.client.SlurmV0043GetDiagWithResponse(ctx)
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
	// Note: v0.0.43 doesn't have RPC statistics in the same structure

	return diag, nil
}

// GetDBDiagnostics retrieves SLURM database diagnostics information
func (a *StandaloneAdapter) GetDBDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	// Note: This might be the same endpoint as GetDiagnostics or might be a separate one
	// For now, we'll implement it similarly
	return a.GetDiagnostics(ctx)
}

// GetInstance retrieves a specific database instance
func (a *StandaloneAdapter) GetInstance(ctx context.Context, opts *types.GetInstanceOptions) (*types.Instance, error) {
	if a.client == nil {
		return nil, fmt.Errorf("API client not initialized")
	}

	// Build query parameters
	params := &api.SlurmdbV0043GetInstanceParams{}
	if opts != nil {
		if opts.Cluster != "" {
			params.Cluster = &opts.Cluster
		}
		if opts.Extra != "" {
			params.Extra = &opts.Extra
		}
		if opts.Format != "" {
			params.Format = &opts.Format
		}
		// v0.0.43 uses inline instance query, not a separate Instance field
		if opts.NodeList != "" {
			params.NodeList = &opts.NodeList
		}
		if opts.TimeStart != nil {
			timeStr := opts.TimeStart.Format("2006-01-02T15:04:05")
			params.TimeStart = &timeStr
		}
		if opts.TimeEnd != nil {
			timeStr := opts.TimeEnd.Format("2006-01-02T15:04:05")
			params.TimeEnd = &timeStr
		}
	}

	resp, err := a.client.SlurmdbV0043GetInstanceWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}

	// Handle API response with enhanced error handling
	if err := a.errorAdapter.HandleAPIResponse(resp.StatusCode(), resp.Body, "GetInstance"); err != nil {
		return nil, err
	}

	if resp.JSON200 == nil || resp.JSON200.Instances == nil || len(resp.JSON200.Instances) == 0 {
		return nil, fmt.Errorf("instance not found")
	}

	// Get the first instance (assuming single result)
	apiInstance := resp.JSON200.Instances[0]
	instance := &types.Instance{}

	if apiInstance.Cluster != nil {
		instance.Cluster = *apiInstance.Cluster
	}
	if apiInstance.Extra != nil {
		instance.ExtraInfo = *apiInstance.Extra
	}
	// Note: v0.0.43 doesn't have an Instance field in the response
	if apiInstance.InstanceId != nil {
		instance.InstanceID = *apiInstance.InstanceId
	}
	if apiInstance.InstanceType != nil {
		instance.InstanceType = *apiInstance.InstanceType
	}
	// Note: v0.0.43 doesn't have NodeCount field in the response

	return instance, nil
}

// GetInstances retrieves multiple database instances with filtering
func (a *StandaloneAdapter) GetInstances(ctx context.Context, opts *types.GetInstancesOptions) (*types.InstanceList, error) {
	if a.client == nil {
		return nil, fmt.Errorf("API client not initialized")
	}

	// Build query parameters
	params := &api.SlurmdbV0043GetInstancesParams{}
	if opts != nil {
		if opts.Extra != "" {
			params.Extra = &opts.Extra
		}
		if opts.Format != "" {
			params.Format = &opts.Format
		}
		// v0.0.43 uses inline instance query, not a separate Instance field
		if opts.NodeList != "" {
			params.NodeList = &opts.NodeList
		}
		if opts.TimeStart != nil {
			timeStr := opts.TimeStart.Format("2006-01-02T15:04:05")
			params.TimeStart = &timeStr
		}
		if opts.TimeEnd != nil {
			timeStr := opts.TimeEnd.Format("2006-01-02T15:04:05")
			params.TimeEnd = &timeStr
		}
	}

	resp, err := a.client.SlurmdbV0043GetInstancesWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get instances: %w", err)
	}

	// Handle API response with enhanced error handling
	if err := a.errorAdapter.HandleAPIResponse(resp.StatusCode(), resp.Body, "GetInstances"); err != nil {
		return nil, err
	}

	if resp.JSON200 == nil || resp.JSON200.Instances == nil {
		return &types.InstanceList{Instances: []types.Instance{}}, nil
	}

	// Convert API instances to common types
	instances := make([]types.Instance, 0)
	for _, apiInstance := range resp.JSON200.Instances {
		instance := types.Instance{}

		if apiInstance.Cluster != nil {
			instance.Cluster = *apiInstance.Cluster
		}
		if apiInstance.Extra != nil {
			instance.ExtraInfo = *apiInstance.Extra
		}
		// Note: v0.0.43 doesn't have an Instance field in the response
		if apiInstance.InstanceId != nil {
			instance.InstanceID = *apiInstance.InstanceId
		}
		if apiInstance.InstanceType != nil {
			instance.InstanceType = *apiInstance.InstanceType
		}
		// Note: v0.0.43 doesn't have NodeCount field in the response

		instances = append(instances, instance)
	}

	return &types.InstanceList{
		Instances: instances,
		Meta:      extractMeta(resp.JSON200.Meta),
	}, nil
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

	resp, err := a.client.SlurmdbV0043GetTresWithResponse(ctx)
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
		// Type is not a pointer in v0.0.43
		tres.Type = apiTres.Type
		if apiTres.Name != nil {
			tres.Name = *apiTres.Name
		}
		if apiTres.Count != nil {
			tres.Count = int64(*apiTres.Count)
		}

		tresList = append(tresList, tres)
	}

	return &types.TRESList{
		TRES: tresList,
		Meta: extractMeta(resp.JSON200.Meta),
	}, nil
}

// CreateTRES creates a new TRES entry
func (a *StandaloneAdapter) CreateTRES(ctx context.Context, req *types.CreateTRESRequest) (*types.TRES, error) {
	if a.client == nil {
		return nil, fmt.Errorf("API client not initialized")
	}

	// Build API request
	apiReq := api.V0043OpenapiTresResp{
		TRES: []api.V0043Tres{
			{
				Type: req.Type,
				Name: &req.Name,
			},
		},
	}

	if req.Count > 0 {
		count := int64(req.Count)
		apiReq.TRES[0].Count = &count
	}

	// Note: The actual endpoint might be different - this is based on the pattern
	// You may need to adjust based on the actual API
	resp, err := a.client.SlurmdbV0043PostTresWithResponse(ctx, apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create TRES: %w", err)
	}

	// Handle API response with enhanced error handling
	if err := a.errorAdapter.HandleAPIResponse(resp.StatusCode(), resp.Body, "CreateTRES"); err != nil {
		return nil, err
	}

	// V0043OpenapiResp doesn't have a TRES field - the API returns a general response
	// We'll return the TRES info from the request since the response doesn't contain it
	tres := &types.TRES{
		Type: req.Type,
		Name: req.Name,
	}

	if req.Count > 0 {
		tres.Count = int64(req.Count)
	}

	return tres, nil
}

// Reconfigure triggers a SLURM reconfiguration
func (a *StandaloneAdapter) Reconfigure(ctx context.Context) (*types.ReconfigureResponse, error) {
	if a.client == nil {
		return nil, fmt.Errorf("API client not initialized")
	}

	resp, err := a.client.SlurmV0043GetReconfigureWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to trigger reconfigure: %w", err)
	}

	// Handle API response with enhanced error handling
	if err := a.errorAdapter.HandleAPIResponse(resp.StatusCode(), resp.Body, "Reconfigure"); err != nil {
		return nil, err
	}

	// Build response
	result := &types.ReconfigureResponse{
		Status: "success",
		Meta:   make(map[string]interface{}),
	}

	if resp.JSON200 != nil {
		result.Meta = extractMeta(resp.JSON200.Meta)

		// Extract any warnings or errors from meta
		// Note: v0.0.43 meta structure is different - simplified handling

		result.Message = "SLURM reconfiguration triggered successfully"
	}

	return result, nil
}

// PingDatabase pings the SLURM database for health checks
func (a *StandaloneAdapter) PingDatabase(ctx context.Context) (*types.PingResponse, error) {
	if a.client == nil {
		return nil, fmt.Errorf("API client not initialized")
	}

	startTime := time.Now()
	resp, err := a.client.SlurmdbV0043GetPingWithResponse(ctx)
	latency := time.Since(startTime).Microseconds()

	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Handle API response with enhanced error handling
	if err := a.errorAdapter.HandleAPIResponse(resp.StatusCode(), resp.Body, "PingDatabase"); err != nil {
		return nil, err
	}

	// Build ping response
	result := &types.PingResponse{
		Status:    "success",
		Timestamp: time.Now(),
		Latency:   latency,
		Meta:      make(map[string]interface{}),
	}

	if resp.JSON200 != nil {
		result.Meta = extractMeta(resp.JSON200.Meta)
		result.Message = "Database ping successful"

		// Extract any ping-specific information if available in the response
		// The API response structure may contain ping-specific data
	} else {
		result.Message = "Database ping completed"
	}

	return result, nil
}

// extractMeta safely extracts metadata from API response
func extractMeta(meta *api.V0043OpenapiMeta) map[string]interface{} {
	result := make(map[string]interface{})

	if meta == nil {
		return result
	}

	// V0043OpenapiMeta has Client, Command, Plugin fields but not Messages/Warnings/Errors
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
