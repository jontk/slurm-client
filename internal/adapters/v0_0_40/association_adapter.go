// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"context"
	"strings"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/jontk/slurm-client/pkg/errors"
	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
)

// AssociationAdapter implements the AssociationAdapter interface for v0.0.40
type AssociationAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewAssociationAdapter creates a new Association adapter for v0.0.40
func NewAssociationAdapter(client *api.ClientWithResponses) *AssociationAdapter {
	return &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.40", "Association"),
		client:      client,
		wrapper:     nil,
	}
}

// List retrieves a list of associations with optional filtering
func (a *AssociationAdapter) List(ctx context.Context, opts *types.AssociationListOptions) (*types.AssociationList, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmdbV0040GetAssociationsParams{}

	// Apply filters from options
	if opts != nil {
		if len(opts.Accounts) > 0 {
			accountStr := strings.Join(opts.Accounts, ",")
			params.Account = &accountStr
		}
		if len(opts.Clusters) > 0 {
			clusterStr := strings.Join(opts.Clusters, ",")
			params.Cluster = &clusterStr
		}
		if len(opts.Users) > 0 {
			userStr := strings.Join(opts.Users, ",")
			params.User = &userStr
		}
		if len(opts.Partitions) > 0 {
			partitionStr := strings.Join(opts.Partitions, ",")
			params.Partition = &partitionStr
		}
		if opts.WithDeleted {
			withDeleted := "true"
			params.WithDeleted = &withDeleted
		}
		// Note: WithSubAccounts not supported in AssociationListOptions for v0.0.40
		// if opts.WithSubAccounts {
		// 	withSubAccts := "true"
		// 	params.WithSubAccts = &withSubAccts
		// }
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0040GetAssociationsWithResponse(ctx, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.40"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "List Associations"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Associations, "List Associations - associations field"); err != nil {
		return nil, err
	}

	// Convert the response to common types
	associationList := make([]types.Association, 0, len(resp.JSON200.Associations))
	for _, apiAssociation := range resp.JSON200.Associations {
		association, err := a.convertAPIAssociationToCommon(apiAssociation)
		if err != nil {
			return nil, a.HandleConversionError(err, "association")
		}
		associationList = append(associationList, *association)
	}

	// Apply pagination
	listOpts := base.ListOptions{}
	if opts != nil {
		listOpts.Limit = opts.Limit
		listOpts.Offset = opts.Offset
	}

	// Apply pagination
	start := listOpts.Offset
	if start < 0 {
		start = 0
	}
	if start >= len(associationList) {
		return &types.AssociationList{
			Associations: []types.Association{},
			Total:        len(associationList),
		}, nil
	}

	end := len(associationList)
	if listOpts.Limit > 0 {
		end = start + listOpts.Limit
		if end > len(associationList) {
			end = len(associationList)
		}
	}

	return &types.AssociationList{
		Associations: associationList[start:end],
		Total:        len(associationList),
	}, nil
}

// Get retrieves a specific association by ID
func (a *AssociationAdapter) Get(ctx context.Context, associationID string) (*types.Association, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.ValidateResourceName(associationID, "associationID"); err != nil {
		return nil, err
	}

	// v0.0.40 doesn't have a single association GET endpoint
	// We need to list all and filter
	list, err := a.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Find the association by constructing ID from fields
	for _, assoc := range list.Associations {
		constructedID := a.constructAssociationID(assoc)
		if constructedID == associationID {
			return &assoc, nil
		}
	}

	return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, "Association '"+associationID+"' not found")
}

// Create creates a new association
func (a *AssociationAdapter) Create(ctx context.Context, association *types.AssociationCreate) (*types.AssociationCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.validateAssociationCreate(association); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert to API format
	apiAssociation, err := a.convertCommonAssociationCreateToAPI(association)
	if err != nil {
		return nil, err
	}

	// Create request body
	// Convert V0040AssocShort to V0040Assoc for the associations list
	assoc := api.V0040Assoc{
		Account: apiAssociation.Account,
		Cluster: apiAssociation.Cluster,
		Partition: apiAssociation.Partition,
		Id: apiAssociation,  // V0040AssocShort is used as the Id field in V0040Assoc
	}
	
	reqBody := api.SlurmdbV0040PostAssociationsJSONRequestBody{
		Associations: api.V0040AssocList{assoc},
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0040PostAssociationsWithResponse(ctx, reqBody)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.40"); err != nil {
		return nil, err
	}

	return &types.AssociationCreateResponse{
		Status:  "success",
		Message: "Association created successfully",
		Meta: map[string]interface{}{
			"association_id": a.constructAssociationIDFromCreate(*association),
			"account_name":   association.AccountName,
			"user_name":      association.UserName,
			"cluster":        association.Cluster,
		},
	}, nil
}

// Update updates an existing association
func (a *AssociationAdapter) Update(ctx context.Context, associationID string, update *types.AssociationUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(associationID, "associationID"); err != nil {
		return err
	}

	// v0.0.40 may not support association updates directly
	return errors.NewNotImplementedError("Update Association", "v0.0.40")
}

// Delete deletes an association
func (a *AssociationAdapter) Delete(ctx context.Context, associationID string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(associationID, "associationID"); err != nil {
		return err
	}

	// v0.0.40 doesn't have a direct association delete endpoint
	// You would typically delete by removing the user from account
	return errors.NewNotImplementedError("Delete Association", "v0.0.40")
}

// constructAssociationID constructs an ID from association fields
func (a *AssociationAdapter) constructAssociationID(assoc types.Association) string {
	parts := []string{assoc.Cluster, assoc.AccountName, assoc.UserName}
	if assoc.Partition != "" {
		parts = append(parts, assoc.Partition)
	}
	return strings.Join(parts, ":")
}

// validateAssociationCreate validates association creation request
func (a *AssociationAdapter) validateAssociationCreate(association *types.AssociationCreate) error {
	if association == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "association creation data is required", "association", nil, nil)
	}
	if association.AccountName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "account is required", "account", association.AccountName, nil)
	}
	if association.UserName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "user is required", "user", association.UserName, nil)
	}
	if association.Cluster == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "cluster is required", "cluster", association.Cluster, nil)
	}
	return nil
}

// constructAssociationIDFromCreate constructs an ID from association create fields
func (a *AssociationAdapter) constructAssociationIDFromCreate(assoc types.AssociationCreate) string {
	parts := []string{assoc.Cluster, assoc.AccountName, assoc.UserName}
	if assoc.Partition != "" {
		parts = append(parts, assoc.Partition)
	}
	return strings.Join(parts, ":")
}
