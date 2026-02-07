// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"context"
	"fmt"
	"strconv"

	types "github.com/jontk/slurm-client/api"
	adapterbase "github.com/jontk/slurm-client/internal/adapters/base"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
)

// AssociationAdapter implements the AssociationAdapter interface for v0.0.41
type AssociationAdapter struct {
	*adapterbase.BaseManager
	client *api.ClientWithResponses
}

// NewAssociationAdapter creates a new Association adapter for v0.0.41
func NewAssociationAdapter(client *api.ClientWithResponses) *AssociationAdapter {
	return &AssociationAdapter{
		BaseManager: adapterbase.NewBaseManager("v0.0.41", "Association"),
		client:      client,
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
	// Note: AssociationList doesn't have a Meta field in common types
	// Warnings and errors from the response are being ignored for now
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
		// ID is *int32, compare after converting id string to int
		if assoc.ID != nil && strconv.Itoa(int(*assoc.ID)) == id {
			return &assoc, nil
		}
	}
	return nil, a.HandleNotFound("association with ID " + id)
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

	// Convert to API request using JSON marshaling workaround
	reqBody, err := a.convertAssociationCreateToAPI(req)
	if err != nil {
		return nil, a.WrapError(err, "failed to convert association create request")
	}

	// Make the API call
	resp, err := a.client.SlurmdbV0041PostAssociationsWithResponse(ctx, reqBody)
	if err != nil {
		return nil, a.WrapError(err, "failed to create association")
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}

	// Check for errors in response
	if resp.JSON200 != nil && resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		errMsgs := make([]string, 0)
		for _, apiErr := range *resp.JSON200.Errors {
			if apiErr.Error != nil {
				errMsgs = append(errMsgs, *apiErr.Error)
			}
		}
		if len(errMsgs) > 0 {
			return nil, fmt.Errorf("association creation failed: %v", errMsgs)
		}
	}

	return &types.AssociationCreateResponse{
		Status:  "success",
		Message: fmt.Sprintf("Created association for %s:%s:%s", req.Account, req.User, req.Cluster),
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

	// Convert to API request using JSON marshaling workaround
	reqBody, err := a.convertAssociationUpdateToAPI(id, update)
	if err != nil {
		return a.WrapError(err, "failed to convert association update request")
	}

	// Make the API call
	resp, err := a.client.SlurmdbV0041PostAssociationsWithResponse(ctx, reqBody)
	if err != nil {
		return a.WrapError(err, "failed to update association")
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
	}

	// Check for errors in response
	if resp.JSON200 != nil && resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		errMsgs := make([]string, 0)
		for _, apiErr := range *resp.JSON200.Errors {
			if apiErr.Error != nil {
				errMsgs = append(errMsgs, *apiErr.Error)
			}
		}
		if len(errMsgs) > 0 {
			return fmt.Errorf("association update failed: %v", errMsgs)
		}
	}

	return nil
}

// Delete deletes an association by ID or composite key.
// Supports two formats:
//   - Numeric ID: "123"
//   - Composite key: "account:user:cluster[:partition]"
func (a *AssociationAdapter) Delete(ctx context.Context, id string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	// Validate ID
	if id == "" {
		return a.HandleValidationError("association ID cannot be empty")
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Check if this is a numeric ID or composite key
	_, numErr := strconv.Atoi(id)
	isNumericID := numErr == nil

	if isNumericID {
		// Use ID-based deletion
		params := &api.SlurmdbV0041DeleteAssociationParams{
			Id: &id,
		}
		resp, err := a.client.SlurmdbV0041DeleteAssociationWithResponse(ctx, params)
		if err != nil {
			return a.WrapError(err, "failed to delete association "+id)
		}
		if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
			return err
		}
		return nil
	}

	// Parse composite key format: "account:user:cluster[:partition]"
	parts := make([]string, 0, 4)
	current := ""
	for _, c := range id {
		if c == ':' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	parts = append(parts, current)

	if len(parts) < 3 {
		return a.HandleValidationError("invalid association ID format, expected numeric ID or 'account:user:cluster[:partition]'")
	}

	account := parts[0]
	user := parts[1]
	cluster := parts[2]

	// Use filter-based deletion
	params := &api.SlurmdbV0041DeleteAssociationsParams{
		Account: &account,
		User:    &user,
		Cluster: &cluster,
	}
	if len(parts) > 3 && parts[3] != "" {
		params.Partition = &parts[3]
	}

	resp, err := a.client.SlurmdbV0041DeleteAssociationsWithResponse(ctx, params)
	if err != nil {
		return a.WrapError(err, "failed to delete association "+id)
	}
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
