// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"context"

	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/pkg/errors"
)

// StandaloneAdapter implements the standalone operations for v0.0.40
type StandaloneAdapter struct {
	client  *api.ClientWithResponses
	version string
}

// NewStandaloneAdapter creates a new standalone adapter
func NewStandaloneAdapter(client *api.ClientWithResponses) *StandaloneAdapter {
	return &StandaloneAdapter{
		client:  client,
		version: "v0.0.40",
	}
}

// GetLicenses returns not implemented error for v0.0.40
func (a *StandaloneAdapter) GetLicenses(ctx context.Context) (*types.LicenseList, error) {
	return nil, errors.NewNotImplementedError("GetLicenses", a.version)
}

// GetShares returns not implemented error for v0.0.40
func (a *StandaloneAdapter) GetShares(ctx context.Context, opts *types.GetSharesOptions) (*types.SharesList, error) {
	return nil, errors.NewNotImplementedError("GetShares", a.version)
}

// GetConfig returns not implemented error for v0.0.40
func (a *StandaloneAdapter) GetConfig(ctx context.Context) (*types.Config, error) {
	return nil, errors.NewNotImplementedError("GetConfig", a.version)
}

// GetDiagnostics returns not implemented error for v0.0.40
func (a *StandaloneAdapter) GetDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	return nil, errors.NewNotImplementedError("GetDiagnostics", a.version)
}

// GetDBDiagnostics returns not implemented error for v0.0.40
func (a *StandaloneAdapter) GetDBDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	return nil, errors.NewNotImplementedError("GetDBDiagnostics", a.version)
}

// GetInstance returns not implemented error for v0.0.40
func (a *StandaloneAdapter) GetInstance(ctx context.Context, opts *types.GetInstanceOptions) (*types.Instance, error) {
	return nil, errors.NewNotImplementedError("GetInstance", a.version)
}

// GetInstances returns not implemented error for v0.0.40
func (a *StandaloneAdapter) GetInstances(ctx context.Context, opts *types.GetInstancesOptions) (*types.InstanceList, error) {
	return nil, errors.NewNotImplementedError("GetInstances", a.version)
}

// GetTRES returns not implemented error for v0.0.40
func (a *StandaloneAdapter) GetTRES(ctx context.Context) (*types.TRESList, error) {
	return nil, errors.NewNotImplementedError("GetTRES", a.version)
}

// CreateTRES returns not implemented error for v0.0.40
func (a *StandaloneAdapter) CreateTRES(ctx context.Context, req *types.CreateTRESRequest) (*types.TRES, error) {
	return nil, errors.NewNotImplementedError("CreateTRES", a.version)
}

// Reconfigure returns not implemented error for v0.0.40
func (a *StandaloneAdapter) Reconfigure(ctx context.Context) (*types.ReconfigureResponse, error) {
	return nil, errors.NewNotImplementedError("Reconfigure", a.version)
}

// PingDatabase pings the SLURM database for health checks
func (a *StandaloneAdapter) PingDatabase(ctx context.Context) (*types.PingResponse, error) {
	return nil, errors.NewNotImplementedError("PingDatabase", a.version)
}
