// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/internal/common"
)

// === Manager Methods ===
// Note: All manager methods are now generated in managers.go

// Requeue requeues a job
func (m *JobManager) Requeue(ctx context.Context, jobID string) error {
	if m.impl == nil {
		m.impl = NewJobManagerImpl(m.client)
	}
	return m.impl.Requeue(ctx, jobID)
}

// Standalone operations - not implemented in this API version

// GetLicenses returns not implemented error
func (c *WrapperClient) GetLicenses(ctx context.Context) (*interfaces.LicenseList, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.42"}
	return stub.GetLicenses(ctx)
}

// GetShares returns not implemented error
func (c *WrapperClient) GetShares(ctx context.Context, opts *interfaces.GetSharesOptions) (*interfaces.SharesList, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.42"}
	return stub.GetShares(ctx, opts)
}

// GetConfig returns not implemented error
func (c *WrapperClient) GetConfig(ctx context.Context) (*interfaces.Config, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.42"}
	return stub.GetConfig(ctx)
}

// GetDiagnostics returns not implemented error
func (c *WrapperClient) GetDiagnostics(ctx context.Context) (*interfaces.Diagnostics, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.42"}
	return stub.GetDiagnostics(ctx)
}

// GetDBDiagnostics returns not implemented error
func (c *WrapperClient) GetDBDiagnostics(ctx context.Context) (*interfaces.Diagnostics, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.42"}
	return stub.GetDBDiagnostics(ctx)
}

// GetInstance returns not implemented error
func (c *WrapperClient) GetInstance(ctx context.Context, opts *interfaces.GetInstanceOptions) (*interfaces.Instance, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.42"}
	return stub.GetInstance(ctx, opts)
}

// GetInstances returns not implemented error
func (c *WrapperClient) GetInstances(ctx context.Context, opts *interfaces.GetInstancesOptions) (*interfaces.InstanceList, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.42"}
	return stub.GetInstances(ctx, opts)
}

// GetTRES returns not implemented error
func (c *WrapperClient) GetTRES(ctx context.Context) (*interfaces.TRESList, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.42"}
	return stub.GetTRES(ctx)
}

// CreateTRES returns not implemented error
func (c *WrapperClient) CreateTRES(ctx context.Context, req *interfaces.CreateTRESRequest) (*interfaces.TRES, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.42"}
	return stub.CreateTRES(ctx, req)
}

// Reconfigure returns not implemented error
func (c *WrapperClient) Reconfigure(ctx context.Context) (*interfaces.ReconfigureResponse, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.42"}
	return stub.Reconfigure(ctx)
}
