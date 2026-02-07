// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_44

import (
	"context"

	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/internal/common"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_44"
	"github.com/jontk/slurm-client/pkg/errors"
)

// createAssociationImpl implements the CreateAssociation method for users
func (a *UserAdapter) createAssociationImpl(ctx context.Context, req *types.UserAssociationRequest) (*types.AssociationCreateResponse, error) {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "association request is required", "request", nil, nil)
	}
	if len(req.Users) == 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "at least one user is required", "users", nil, nil)
	}
	if req.Account == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account is required", "account", nil, nil)
	}
	if req.Cluster == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "cluster is required", "cluster", nil, nil)
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}
	// Build associations list - one for each user
	associations := make([]api.V0044Assoc, 0, len(req.Users))
	for _, userName := range req.Users {
		assoc := api.V0044Assoc{
			User:    userName,
			Account: &req.Account,
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
		associations = append(associations, assoc)
	}
	// Create request body
	reqBody := api.SlurmdbV0044PostAssociationsJSONRequestBody{
		Associations: associations,
	}
	// Call the API
	resp, err := a.client.SlurmdbV0044PostAssociationsWithResponse(ctx, reqBody)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}
	// Handle response errors
	var apiErrors *api.V0044OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.44"); err != nil {
		return nil, err
	}
	// Return success response
	return &types.AssociationCreateResponse{
		Status: "created",
	}, nil
}
