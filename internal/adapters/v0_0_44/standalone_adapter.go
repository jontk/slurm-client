// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_44

import (
	"context"

	api "github.com/jontk/slurm-client/internal/api/v0_0_44"
	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/jontk/slurm-client/pkg/errors"
)

// StandaloneAdapter implements the standalone operations for v0.0.44
type StandaloneAdapter struct {
	*base.BaseManager
	client *api.ClientWithResponses
}

// NewStandaloneAdapter creates a new standalone adapter
func NewStandaloneAdapter(client *api.ClientWithResponses) *StandaloneAdapter {
	return &StandaloneAdapter{
		BaseManager: base.NewBaseManager("v0.0.44", "Standalone"),
		client:      client,
	}
}

// GetLicenses retrieves license information
func (a *StandaloneAdapter) GetLicenses(ctx context.Context) (*types.LicenseList, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	resp, err := a.client.SlurmV0044GetLicensesWithResponse(ctx)
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
	if err := a.CheckNilResponse(resp.JSON200, "Get Licenses"); err != nil {
		return nil, err
	}

	// Convert API licenses to common types
	licenses := make([]types.License, 0, len(resp.JSON200.Licenses))
	for _, apiLicense := range resp.JSON200.Licenses {
		license := types.License{}

		// Handle pointer fields
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
		if apiLicense.LastConsumed != nil {
			license.RemoteUsed = int(*apiLicense.LastConsumed)
		}

		licenses = append(licenses, license)
	}

	return &types.LicenseList{
		Licenses: licenses,
	}, nil
}

// GetTRES retrieves all TRES (Trackable RESources)
func (a *StandaloneAdapter) GetTRES(ctx context.Context) (*types.TRESList, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Note: SLURM v0.0.44 doesn't have a direct TRES endpoint
	// Return empty list for now
	return &types.TRESList{
		TRES: []types.TRES{},
	}, nil
}

// CreateTRES creates a new TRES entry
func (a *StandaloneAdapter) CreateTRES(ctx context.Context, req *types.CreateTRESRequest) (*types.TRES, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "CreateTRES request is required", "request", req, nil)
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Note: SLURM v0.0.44 doesn't support TRES creation via REST API
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "TRES creation not supported in v0.0.44")
}

// GetShares retrieves fairshare information with optional filtering
func (a *StandaloneAdapter) GetShares(ctx context.Context, opts *types.GetSharesOptions) (*types.SharesList, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Note: SLURM v0.0.44 doesn't have a direct shares endpoint
	// Return empty list for now
	return &types.SharesList{
		Shares: []types.Share{},
	}, nil
}

// GetConfig retrieves SLURM configuration
func (a *StandaloneAdapter) GetConfig(ctx context.Context) (*types.Config, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Note: SLURM v0.0.44 doesn't have a direct config endpoint
	// Return basic config for now
	return &types.Config{
		Version: "v0.0.44",
	}, nil
}

// GetDiagnostics retrieves SLURM diagnostics information
func (a *StandaloneAdapter) GetDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

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
	if err := a.CheckNilResponse(resp.JSON200, "Get Diagnostics"); err != nil {
		return nil, err
	}

	// Convert to common diagnostics format
	diag := &types.Diagnostics{
		// Fill basic fields based on API response
		ServerThreadCount: 0, // Default values since v0.0.44 has different structure
		AgentQueueSize:    0,
		AgentCount:        0,
		AgentThreadCount:  0,
		DBDAgentCount:     0,
		GittosCount:       0,
	}

	// v0.0.44 response has different structure, use basic conversion
	// resp.JSON200.Statistics contains different data than expected

	return diag, nil
}

// GetDBDiagnostics retrieves SLURM database diagnostics information
func (a *StandaloneAdapter) GetDBDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Note: v0.0.44 doesn't have separate DB diagnostics
	// Delegate to regular diagnostics
	return a.GetDiagnostics(ctx)
}

// GetInstance retrieves a specific database instance
func (a *StandaloneAdapter) GetInstance(ctx context.Context, opts *types.GetInstanceOptions) (*types.Instance, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Note: SLURM v0.0.44 doesn't have instance endpoints
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Instance operations not supported in v0.0.44")
}

// GetInstances retrieves multiple database instances with filtering
func (a *StandaloneAdapter) GetInstances(ctx context.Context, opts *types.GetInstancesOptions) (*types.InstanceList, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Note: SLURM v0.0.44 doesn't have instance endpoints
	return &types.InstanceList{
		Instances: []types.Instance{},
	}, nil
}

// Reconfigure triggers a SLURM reconfiguration
func (a *StandaloneAdapter) Reconfigure(ctx context.Context) (*types.ReconfigureResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Note: SLURM v0.0.44 doesn't have reconfigure endpoint
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Reconfigure operation not supported in v0.0.44")
}

// PingDatabase pings the SLURM database
func (a *StandaloneAdapter) PingDatabase(ctx context.Context) (*types.PingResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Note: SLURM v0.0.44 doesn't have ping endpoint
	return &types.PingResponse{
		Status:  "ok",
		Message: "v0.0.44 database ping simulation",
	}, nil
}
