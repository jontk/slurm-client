// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"context"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// StandaloneOperationsStub provides stub implementations for standalone operations in older API versions
type StandaloneOperationsStub struct {
	Version string
}

// GetLicenses returns not implemented error
func (s *StandaloneOperationsStub) GetLicenses(ctx context.Context) (*interfaces.LicenseList, error) {
	return nil, errors.NewNotImplementedError("license listing", s.Version)
}

// GetShares returns not implemented error
func (s *StandaloneOperationsStub) GetShares(ctx context.Context, opts *interfaces.GetSharesOptions) (*interfaces.SharesList, error) {
	return nil, errors.NewNotImplementedError("shares listing", s.Version)
}

// GetConfig returns not implemented error
func (s *StandaloneOperationsStub) GetConfig(ctx context.Context) (*interfaces.Config, error) {
	return nil, errors.NewNotImplementedError("config retrieval", s.Version)
}

// GetDiagnostics returns not implemented error
func (s *StandaloneOperationsStub) GetDiagnostics(ctx context.Context) (*interfaces.Diagnostics, error) {
	return nil, errors.NewNotImplementedError("diagnostics retrieval", s.Version)
}

// GetDBDiagnostics returns not implemented error
func (s *StandaloneOperationsStub) GetDBDiagnostics(ctx context.Context) (*interfaces.Diagnostics, error) {
	return nil, errors.NewNotImplementedError("database diagnostics retrieval", s.Version)
}

// GetInstance returns not implemented error
func (s *StandaloneOperationsStub) GetInstance(ctx context.Context, opts *interfaces.GetInstanceOptions) (*interfaces.Instance, error) {
	return nil, errors.NewNotImplementedError("instance retrieval", s.Version)
}

// GetInstances returns not implemented error
func (s *StandaloneOperationsStub) GetInstances(ctx context.Context, opts *interfaces.GetInstancesOptions) (*interfaces.InstanceList, error) {
	return nil, errors.NewNotImplementedError("instances listing", s.Version)
}

// GetTRES returns not implemented error
func (s *StandaloneOperationsStub) GetTRES(ctx context.Context) (*interfaces.TRESList, error) {
	return nil, errors.NewNotImplementedError("TRES listing", s.Version)
}

// CreateTRES returns not implemented error
func (s *StandaloneOperationsStub) CreateTRES(ctx context.Context, req *interfaces.CreateTRESRequest) (*interfaces.TRES, error) {
	return nil, errors.NewNotImplementedError("TRES creation", s.Version)
}

// Reconfigure returns not implemented error
func (s *StandaloneOperationsStub) Reconfigure(ctx context.Context) (*interfaces.ReconfigureResponse, error) {
	return nil, errors.NewNotImplementedError("reconfiguration", s.Version)
}
