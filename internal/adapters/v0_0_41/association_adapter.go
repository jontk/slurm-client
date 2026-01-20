// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

import (
	"context"
	"fmt"
	"strings"

	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
)

// AssociationAdapter implements the AssociationAdapter interface for v0.0.41
type AssociationAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewAssociationAdapter creates a new Association adapter for v0.0.41
func NewAssociationAdapter(client *api.ClientWithResponses) *AssociationAdapter {
	return &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Association"),
		client:      client,
		wrapper:     nil, // We'll implement this later
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
	params := &api.SlurmdbV0041GetAssociationsParams{}

	// Apply filters from options
	if opts != nil {
		// Map AssociationListOptions fields to API parameters
		if len(opts.Accounts) > 0 {
			// API takes a single account, use the first one
			params.Account = &opts.Accounts[0]
		}
		if len(opts.Users) > 0 {
			// API takes a single user, use the first one
			params.User = &opts.Users[0]
		}
		if len(opts.Clusters) > 0 {
			// API takes a single cluster, use the first one
			params.Cluster = &opts.Clusters[0]
		}
		if opts.WithDeleted {
			withDeleted := "true"
			params.WithDeleted = &withDeleted
		}
		// Note: ParentAccount and WithSubAccounts are not in the common types
		// These would need to be handled differently if needed
	}

	// Make the API call
	resp, err := a.client.SlurmdbV0041GetAssociationsWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list associations")
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected nil response")
	}

	// Convert response to common types
	assocList := &types.AssociationList{
		Associations: make([]types.Association, 0, len(resp.JSON200.Associations)),
		Total:        len(resp.JSON200.Associations),
	}

	for _, apiAssoc := range resp.JSON200.Associations {
		assoc, err := a.convertAPIAssociationToCommon(apiAssoc)
		if err != nil {
			// Log the error but continue processing other associations
			continue
		}
		assocList.Associations = append(assocList.Associations, *assoc)
	}

	// Extract warning messages if any
	if resp.JSON200.Warnings != nil {
		warnings := make([]string, 0, len(*resp.JSON200.Warnings))
		for _, warning := range *resp.JSON200.Warnings {
			if warning.Description != nil {
				warnings = append(warnings, *warning.Description)
			}
		}
		if len(warnings) > 0 {
			// Note: AssociationList doesn't have a Meta field in common types
			// Warnings are being ignored for now
		}
	}

	// Extract error messages if any
	if resp.JSON200.Errors != nil {
		errors := make([]string, 0, len(*resp.JSON200.Errors))
		for _, error := range *resp.JSON200.Errors {
			if error.Description != nil {
				errors = append(errors, *error.Description)
			}
		}
		if len(errors) > 0 {
			// AssociationList doesn't have a Meta field in common types
			// Log errors but continue
		}
	}

	return assocList, nil
}

// Get retrieves a specific association by ID
func (a *AssociationAdapter) Get(ctx context.Context, id string) (*types.Association, error) {
	// v0.0.41 doesn't support getting a single association by ID
	// We need to list all associations and filter
	assocList, err := a.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	for _, assoc := range assocList.Associations {
		if assoc.ID == id {
			return &assoc, nil
		}
	}

	return nil, a.HandleNotFound(fmt.Sprintf("association with ID %s", id))
}

// Create creates a new association
func (a *AssociationAdapter) Create(ctx context.Context, req *types.AssociationCreate) (*types.AssociationCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Validate request
	if req == nil {
		return nil, a.HandleValidationError("association create request cannot be nil")
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert request to association for API call
	association := &types.Association{
		AccountName:   req.AccountName,
		UserName:      req.UserName,
		Cluster:       req.Cluster,
		Partition:     req.Partition,
		DefaultQoS:    req.DefaultQoS,
		SharesRaw:     req.SharesRaw,
		Priority:      req.Priority,
		ParentAccount: req.ParentAccount,
	}

	// Convert association to API request
	createReq := a.convertCommonToAPIAssociation(association)

	// Make the API call
	resp, err := a.client.SlurmdbV0041PostAssociationsWithResponse(ctx, *createReq)
	if err != nil {
		return nil, a.WrapError(err, "failed to create association")
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}

	return &types.AssociationCreateResponse{
		Status:  "success",
		Message: "Association created successfully",
		Meta: map[string]interface{}{
			"association_id": association.ID,
			"account_name":   association.AccountName,
			"user_name":      association.UserName,
			"cluster":        association.Cluster,
		},
	}, nil
}

// Update updates an existing association
func (a *AssociationAdapter) Update(ctx context.Context, id string, update *types.AssociationUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Validate update
	if update == nil {
		return a.HandleValidationError("association update cannot be nil")
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Note: idNum conversion was here but is not needed since we're using the string ID directly

	// Get the existing association first
	existingAssoc, err := a.Get(ctx, id)
	if err != nil {
		return err
	}

	// Apply updates
	if update.DefaultQoS != nil {
		existingAssoc.DefaultQoS = *update.DefaultQoS
	}
	if update.SharesRaw != nil {
		existingAssoc.SharesRaw = *update.SharesRaw
	}
	if update.Priority != nil {
		existingAssoc.Priority = *update.Priority
	}

	// Convert to API request
	updateReq := a.convertCommonToAPIAssociation(existingAssoc)

	// Make the API call
	resp, err := a.client.SlurmdbV0041PostAssociationsWithResponse(ctx, *updateReq)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to update association %s", id))
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
	}

	return nil
}

// Delete deletes an association
func (a *AssociationAdapter) Delete(ctx context.Context, id string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// v0.0.41 doesn't support deleting by ID directly
	// Parse the composite key (format: "account:user:cluster:partition")
	parts := strings.Split(id, ":")
	if len(parts) < 3 {
		return a.HandleValidationError("invalid association ID format, expected 'account:user:cluster[:partition]'")
	}

	account := parts[0]
	user := parts[1]
	cluster := parts[2]
	partition := ""
	if len(parts) > 3 {
		partition = parts[3]
	}

	// Make the API call using account, user, cluster, partition
	params := &api.SlurmdbV0041DeleteAssociationsParams{
		Account: &account,
		User:    &user,
		Cluster: &cluster,
	}

	if partition != "" {
		params.Partition = &partition
	}

	resp, err := a.client.SlurmdbV0041DeleteAssociationsWithResponse(ctx, params)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to delete association %s", id))
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
	}

	return nil
}

// SetLimits sets resource limits for an association
// NOTE: This method is not part of the AssociationAdapter interface
// and AssociationLimits type is not defined in common types
/*
// func (a *AssociationAdapter) SetLimits(ctx context.Context, id uint32, limits *types.AssociationLimits) error {
// 	// Use the Update method to set limits
// 	update := &types.AssociationUpdate{}
//
// 	if limits.MaxJobs != nil {
// 		update.MaxJobs = limits.MaxJobs
// 	}
// 	if limits.MaxSubmitJobs != nil {
// 		update.MaxSubmitJobs = limits.MaxSubmitJobs
// 	}
// 	if limits.MaxTRES != nil {
// 		// Convert TRES map to string format
// 		tresStr := formatTRESMap(limits.MaxTRES)
// 		update.MaxTRES = &tresStr
// 	}
// 	if limits.MaxTRESPerJob != nil {
// 		tresStr := formatTRESMap(limits.MaxTRESPerJob)
// 		update.MaxTRESPerJob = &tresStr
// 	}
//
// 	return a.Update(ctx, id, update)
// }

// GetByUserAccount gets associations for a specific user and account
func (a *AssociationAdapter) GetByUserAccount(ctx context.Context, user, account, cluster string) (*types.Association, error) {
	opts := &types.AssociationListOptions{
		Users:    []string{user},
		Accounts: []string{account},
		Clusters: []string{cluster},
	}

	assocList, err := a.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	if len(assocList.Associations) == 0 {
		return nil, a.HandleNotFound(fmt.Sprintf("association for user %s in account %s", user, account))
	}

	// Return the first matching association
	return &assocList.Associations[0], nil
}
*/
