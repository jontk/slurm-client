// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_42

import (
	"context"

	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/internal/common"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_42"
	"github.com/jontk/slurm-client/pkg/errors"
)

// createAssociationImpl implements the CreateAssociation method for accounts
func (a *AccountAdapter) createAssociationImpl(ctx context.Context, req *types.AccountAssociationRequest) (*types.AssociationCreateResponse, error) {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "association request is required", "request", nil, nil)
	}
	if len(req.Accounts) == 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "at least one account is required", "accounts", nil, nil)
	}
	if req.Cluster == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "cluster is required", "cluster", nil, nil)
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}
	// Build associations list - one for each account
	associations := make([]api.V0042Assoc, 0, len(req.Accounts))
	for _, accountName := range req.Accounts {
		assoc := api.V0042Assoc{
			Account: &accountName,
			Cluster: &req.Cluster,
		}
		// Set optional fields
		if req.DefaultQoS != "" {
			assoc.Default = &struct {
				Qos *string `json:"qos,omitempty"`
			}{
				Qos: &req.DefaultQoS,
			}
		}
		if req.Description != "" {
			assoc.Comment = &req.Description
		}
		associations = append(associations, assoc)
	}
	// Create request body
	reqBody := api.SlurmdbV0042PostAssociationsJSONRequestBody{
		Associations: associations,
	}
	// Call the API
	resp, err := a.client.SlurmdbV0042PostAssociationsWithResponse(ctx, reqBody)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}
	// Handle response errors
	var apiErrors *api.V0042OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.42"); err != nil {
		return nil, err
	}
	// Return success response
	return &types.AssociationCreateResponse{
		Status: "created",
	}, nil
}
