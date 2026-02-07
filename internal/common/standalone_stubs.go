// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"context"

	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/pkg/errors"
)

// StandaloneOperationsStub provides stub implementations for standalone operations in older API versions
type StandaloneOperationsStub struct {
	Version string
}

// GetLicenses returns not implemented error
func (s *StandaloneOperationsStub) GetLicenses(ctx context.Context) (*types.LicenseList, error) {
	return nil, errors.NewNotImplementedError("license listing", s.Version)
}

// GetShares returns not implemented error
func (s *StandaloneOperationsStub) GetShares(ctx context.Context, opts *types.GetSharesOptions) (*types.SharesList, error) {
	return nil, errors.NewNotImplementedError("shares listing", s.Version)
}

// GetConfig returns not implemented error
func (s *StandaloneOperationsStub) GetConfig(ctx context.Context) (*types.Config, error) {
	return nil, errors.NewNotImplementedError("config retrieval", s.Version)
}

// GetDiagnostics returns not implemented error
func (s *StandaloneOperationsStub) GetDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	return nil, errors.NewNotImplementedError("diagnostics retrieval", s.Version)
}

// GetDBDiagnostics returns not implemented error
func (s *StandaloneOperationsStub) GetDBDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	return nil, errors.NewNotImplementedError("database diagnostics retrieval", s.Version)
}

// GetInstance returns not implemented error
func (s *StandaloneOperationsStub) GetInstance(ctx context.Context, opts *types.GetInstanceOptions) (*types.Instance, error) {
	return nil, errors.NewNotImplementedError("instance retrieval", s.Version)
}

// GetInstances returns not implemented error
func (s *StandaloneOperationsStub) GetInstances(ctx context.Context, opts *types.GetInstancesOptions) (*types.InstanceList, error) {
	return nil, errors.NewNotImplementedError("instances listing", s.Version)
}

// GetTRES returns not implemented error
func (s *StandaloneOperationsStub) GetTRES(ctx context.Context) (*types.TRESList, error) {
	return nil, errors.NewNotImplementedError("TRES listing", s.Version)
}

// CreateTRES returns not implemented error
func (s *StandaloneOperationsStub) CreateTRES(ctx context.Context, req *types.CreateTRESRequest) (*types.TRES, error) {
	return nil, errors.NewNotImplementedError("TRES creation", s.Version)
}

// Reconfigure returns not implemented error
func (s *StandaloneOperationsStub) Reconfigure(ctx context.Context) (*types.ReconfigureResponse, error) {
	return nil, errors.NewNotImplementedError("reconfiguration", s.Version)
}

// WCKeyManagerStub provides stub implementations for WCKey operations in older API versions
type WCKeyManagerStub struct {
	Version string
}

// List returns not implemented error
func (w *WCKeyManagerStub) List(ctx context.Context, opts *types.WCKeyListOptions) (*types.WCKeyList, error) {
	return nil, errors.NewNotImplementedError("WCKey listing", w.Version)
}

// Get returns not implemented error
func (w *WCKeyManagerStub) Get(ctx context.Context, wckeyName, user, cluster string) (*types.WCKey, error) {
	return nil, errors.NewNotImplementedError("WCKey retrieval", w.Version)
}

// Create returns not implemented error
func (w *WCKeyManagerStub) Create(ctx context.Context, wckey *types.WCKeyCreate) (*types.WCKeyCreateResponse, error) {
	return nil, errors.NewNotImplementedError("WCKey creation", w.Version)
}

// Update returns not implemented error
func (w *WCKeyManagerStub) Update(ctx context.Context, wckeyName, user, cluster string, update *types.WCKeyUpdate) error {
	return errors.NewNotImplementedError("WCKey update", w.Version)
}

// Delete returns not implemented error
func (w *WCKeyManagerStub) Delete(ctx context.Context, wckeyID string) error {
	return errors.NewNotImplementedError("WCKey deletion", w.Version)
}
