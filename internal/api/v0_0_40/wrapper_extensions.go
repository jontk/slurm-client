package v0_0_40

import (
	"context"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/interfaces"
)

// Standalone operations - not implemented in this API version

// GetLicenses returns not implemented error
func (c *WrapperClient) GetLicenses(ctx context.Context) (*interfaces.LicenseList, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.40"}
	return stub.GetLicenses(ctx)
}

// GetShares returns not implemented error
func (c *WrapperClient) GetShares(ctx context.Context, opts *interfaces.GetSharesOptions) (*interfaces.SharesList, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.40"}
	return stub.GetShares(ctx, opts)
}

// GetConfig returns not implemented error
func (c *WrapperClient) GetConfig(ctx context.Context) (*interfaces.Config, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.40"}
	return stub.GetConfig(ctx)
}

// GetDiagnostics returns not implemented error
func (c *WrapperClient) GetDiagnostics(ctx context.Context) (*interfaces.Diagnostics, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.40"}
	return stub.GetDiagnostics(ctx)
}

// GetDBDiagnostics returns not implemented error
func (c *WrapperClient) GetDBDiagnostics(ctx context.Context) (*interfaces.Diagnostics, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.40"}
	return stub.GetDBDiagnostics(ctx)
}

// GetInstance returns not implemented error
func (c *WrapperClient) GetInstance(ctx context.Context, opts *interfaces.GetInstanceOptions) (*interfaces.Instance, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.40"}
	return stub.GetInstance(ctx, opts)
}

// GetInstances returns not implemented error
func (c *WrapperClient) GetInstances(ctx context.Context, opts *interfaces.GetInstancesOptions) (*interfaces.InstanceList, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.40"}
	return stub.GetInstances(ctx, opts)
}

// GetTRES returns not implemented error
func (c *WrapperClient) GetTRES(ctx context.Context) (*interfaces.TRESList, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.40"}
	return stub.GetTRES(ctx)
}

// CreateTRES returns not implemented error
func (c *WrapperClient) CreateTRES(ctx context.Context, req *interfaces.CreateTRESRequest) (*interfaces.TRES, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.40"}
	return stub.CreateTRES(ctx, req)
}

// Reconfigure returns not implemented error
func (c *WrapperClient) Reconfigure(ctx context.Context) (*interfaces.ReconfigureResponse, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.40"}
	return stub.Reconfigure(ctx)
}