// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"context"

	"github.com/jontk/slurm-client/interfaces"
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

// WCKeyManagerStub provides stub implementations for WCKey operations in older API versions
type WCKeyManagerStub struct {
	Version string
}

// List returns not implemented error
func (w *WCKeyManagerStub) List(ctx context.Context, opts *interfaces.WCKeyListOptions) (*interfaces.WCKeyList, error) {
	return nil, errors.NewNotImplementedError("WCKey listing", w.Version)
}

// Get returns not implemented error
func (w *WCKeyManagerStub) Get(ctx context.Context, wckeyName, user, cluster string) (*interfaces.WCKey, error) {
	return nil, errors.NewNotImplementedError("WCKey retrieval", w.Version)
}

// Create returns not implemented error
func (w *WCKeyManagerStub) Create(ctx context.Context, wckey *interfaces.WCKeyCreate) (*interfaces.WCKeyCreateResponse, error) {
	return nil, errors.NewNotImplementedError("WCKey creation", w.Version)
}

// Update returns not implemented error
func (w *WCKeyManagerStub) Update(ctx context.Context, wckeyName, user, cluster string, update *interfaces.WCKeyUpdate) error {
	return errors.NewNotImplementedError("WCKey update", w.Version)
}

// Delete returns not implemented error
func (w *WCKeyManagerStub) Delete(ctx context.Context, wckeyID string) error {
	return errors.NewNotImplementedError("WCKey deletion", w.Version)
}
