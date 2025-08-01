// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

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
	stub := &common.StandaloneOperationsStub{Version: "v0.0.41"}
	return stub.GetLicenses(ctx)
}

// GetShares returns not implemented error
func (c *WrapperClient) GetShares(ctx context.Context, opts *interfaces.GetSharesOptions) (*interfaces.SharesList, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.41"}
	return stub.GetShares(ctx, opts)
}

// GetConfig returns not implemented error
func (c *WrapperClient) GetConfig(ctx context.Context) (*interfaces.Config, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.41"}
	return stub.GetConfig(ctx)
}

// GetDiagnostics returns not implemented error
func (c *WrapperClient) GetDiagnostics(ctx context.Context) (*interfaces.Diagnostics, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.41"}
	return stub.GetDiagnostics(ctx)
}

// GetDBDiagnostics returns not implemented error
func (c *WrapperClient) GetDBDiagnostics(ctx context.Context) (*interfaces.Diagnostics, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.41"}
	return stub.GetDBDiagnostics(ctx)
}

// GetInstance returns not implemented error
func (c *WrapperClient) GetInstance(ctx context.Context, opts *interfaces.GetInstanceOptions) (*interfaces.Instance, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.41"}
	return stub.GetInstance(ctx, opts)
}

// GetInstances returns not implemented error
func (c *WrapperClient) GetInstances(ctx context.Context, opts *interfaces.GetInstancesOptions) (*interfaces.InstanceList, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.41"}
	return stub.GetInstances(ctx, opts)
}

// GetTRES returns not implemented error
func (c *WrapperClient) GetTRES(ctx context.Context) (*interfaces.TRESList, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.41"}
	return stub.GetTRES(ctx)
}

// CreateTRES returns not implemented error
func (c *WrapperClient) CreateTRES(ctx context.Context, req *interfaces.CreateTRESRequest) (*interfaces.TRES, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.41"}
	return stub.CreateTRES(ctx, req)
}

// Reconfigure returns not implemented error
func (c *WrapperClient) Reconfigure(ctx context.Context) (*interfaces.ReconfigureResponse, error) {
	stub := &common.StandaloneOperationsStub{Version: "v0.0.41"}
	return stub.Reconfigure(ctx)
}

// HandleErrorResponse handles HTTP error responses
func (c *WrapperClient) HandleErrorResponse(statusCode int, body []byte) error {
	// Try to parse the response as structured error
	if len(body) > 0 {
		// Log the body for debugging (in production, you'd use a proper logger)
		// fmt.Printf("Error response body: %s\n", string(body))
	}
	
	// Map HTTP status code to error
	var code errors.ErrorCode
	switch statusCode {
	case 400:
		code = errors.ErrorCodeInvalidRequest
	case 401:
		code = errors.ErrorCodeUnauthorized
	case 403:
		code = errors.ErrorCodePermissionDenied
	case 404:
		code = errors.ErrorCodeResourceNotFound
	case 409:
		code = errors.ErrorCodeConflict
	case 422:
		code = errors.ErrorCodeValidationFailed
	case 429:
		code = errors.ErrorCodeRateLimited
	case 500:
		code = errors.ErrorCodeServerInternal
	case 502, 503, 504:
		code = errors.ErrorCodeSlurmDaemonDown
	default:
		code = errors.ErrorCodeUnknown
	}
	
	return errors.NewSlurmError(code, fmt.Sprintf("HTTP %d error", statusCode))
}
