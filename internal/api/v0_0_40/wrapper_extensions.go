// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"context"
	"fmt"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
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

// CheckContext validates the context
func (c *WrapperClient) CheckContext(ctx context.Context) error {
	if ctx == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "context is required", "ctx", nil, nil)
	}
	return nil
}

// HandleErrorResponse handles error responses from the API
func (c *WrapperClient) HandleErrorResponse(statusCode int, body []byte) error {
	return c.handleHTTPError(statusCode, body)
}

// handleHTTPError handles HTTP error responses
func (c *WrapperClient) handleHTTPError(statusCode int, body []byte) error {
	switch statusCode {
	case 400:
		return errors.NewClientError(errors.ErrorCodeInvalidRequest, fmt.Sprintf("Bad Request: %s", string(body)))
	case 401:
		return errors.NewAuthError("HTTP", "Bearer", fmt.Errorf("authentication failed"))
	case 403:
		return errors.NewAuthError("HTTP", "Bearer", fmt.Errorf("permission denied"))
	case 404:
		return errors.NewClientError(errors.ErrorCodeResourceNotFound, "Resource not found")
	case 500:
		return errors.NewClientError(errors.ErrorCodeServerInternal, fmt.Sprintf("Internal Server Error: %s", string(body)))
	default:
		return errors.NewClientError(errors.ErrorCodeServerInternal, fmt.Sprintf("HTTP %d: %s", statusCode, string(body)))
	}
}
