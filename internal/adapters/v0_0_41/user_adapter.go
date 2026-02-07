// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"context"
	"fmt"

	types "github.com/jontk/slurm-client/api"
	adapterbase "github.com/jontk/slurm-client/internal/adapters/base"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
	"github.com/jontk/slurm-client/pkg/errors"
)

// UserAdapter implements the UserAdapter interface for v0.0.41
type UserAdapter struct {
	*adapterbase.BaseManager
	client *api.ClientWithResponses
}

// NewUserAdapter creates a new User adapter for v0.0.41
func NewUserAdapter(client *api.ClientWithResponses) *UserAdapter {
	return &UserAdapter{
		BaseManager: adapterbase.NewBaseManager("v0.0.41", "User"),
		client:      client,
	}
}

// List retrieves a list of users with optional filtering
func (a *UserAdapter) List(ctx context.Context, opts *types.UserListOptions) (*types.UserList, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}
	// Prepare parameters for the API call
	params := &api.SlurmdbV0041GetUsersParams{}
	// Apply filters from options
	if opts != nil {
		// Names is not supported directly in v0.0.41 API params
		// Skip user name filtering for now
		if len(opts.Names) > 0 {
			_ = opts.Names
		}
		// Account field doesn't exist in UserListOptions
		// Skip account filtering
		if opts.DefaultAccount != "" {
			params.DefaultAccount = &opts.DefaultAccount
		}
		// DefaultWCKey field doesn't exist in UserListOptions
		// Skip DefaultWCKey filtering
		if opts.WithDeleted {
			withDeleted := "true"
			params.WithDeleted = &withDeleted
		}
		if opts.WithAssocs {
			withAssocs := "true"
			params.WithAssocs = &withAssocs
		}
		if opts.WithCoords {
			withCoords := "true"
			params.WithCoords = &withCoords
		}
		if opts.WithWCKeys {
			withWckeys := "true"
			params.WithWckeys = &withWckeys
		}
		// AdminLevel field doesn't exist in UserListOptions
		// Skip AdminLevel filtering
	}
	// Make the API call
	resp, err := a.client.SlurmdbV0041GetUsersWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list users")
	}
	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected nil response")
	}
	// Convert response to common types
	userList := &types.UserList{
		Users: make([]types.User, 0, len(resp.JSON200.Users)),
		Total: 0,
	}
	for _, apiUser := range resp.JSON200.Users {
		user, err := a.convertAPIUserToCommon(apiUser)
		if err != nil {
			// Log the error but continue processing other users
			continue
		}
		userList.Users = append(userList.Users, *user)
	}
	// Extract warning and error messages if any (but UserList doesn't have Meta)
	// Warnings are ignored for now as UserList structure doesn't support them
	if resp.JSON200.Warnings != nil {
		// Log warnings if needed
		_ = resp.JSON200.Warnings
	}
	// Extract error messages if any
	if resp.JSON200.Errors != nil {
		errors := make([]string, 0, len(*resp.JSON200.Errors))
		for _, error := range *resp.JSON200.Errors {
			if error.Description != nil {
				errors = append(errors, *error.Description)
			}
		}
		// UserList doesn't have Meta field
		// Skip error storage
		_ = errors
	}
	return userList, nil
}

// Get retrieves a specific user by name
func (a *UserAdapter) Get(ctx context.Context, name string) (*types.User, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	// Validate name
	if err := a.ValidateResourceName("user name", name); err != nil {
		return nil, err
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}
	// Make the API call
	params := &api.SlurmdbV0041GetUserParams{
		WithAssocs: ptr("true"),
		WithCoords: ptr("true"),
		WithWckeys: ptr("true"),
	}
	resp, err := a.client.SlurmdbV0041GetUserWithResponse(ctx, name, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to get user "+name)
	}
	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil || len(resp.JSON200.Users) == 0 {
		return nil, a.HandleNotFound("user " + name)
	}
	// Convert the first user in the response
	user, err := a.convertAPIUserToCommon(resp.JSON200.Users[0])
	if err != nil {
		return nil, a.WrapError(err, "failed to convert user "+name)
	}
	return user, nil
}

// Create creates a new user
func (a *UserAdapter) Create(ctx context.Context, req *types.UserCreate) (*types.UserCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	// Validate request
	if req == nil {
		return nil, a.HandleValidationError("user create request cannot be nil")
	}
	if err := a.ValidateResourceName("user name", req.Name); err != nil {
		return nil, err
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}
	// Convert request to user for API call
	user := &types.User{
		Name: req.Name,
	}
	// Set Default values if provided
	if req.DefaultAccount != "" || req.DefaultWCKey != "" {
		user.Default = &types.UserDefault{}
		if req.DefaultAccount != "" {
			user.Default.Account = &req.DefaultAccount
		}
		if req.DefaultWCKey != "" {
			user.Default.Wckey = &req.DefaultWCKey
		}
	}
	// Set AdminLevel if provided
	if req.AdminLevel != "" {
		user.AdministratorLevel = []types.AdministratorLevelValue{types.AdministratorLevelValue(req.AdminLevel)}
	}
	// Convert user to API request
	createReq := a.convertCommonToAPIUser(user)
	// Make the API call
	resp, err := a.client.SlurmdbV0041PostUsersWithResponse(ctx, *createReq)
	if err != nil {
		return nil, a.WrapError(err, "failed to create user "+req.Name)
	}
	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}
	return &types.UserCreateResponse{UserName: req.Name}, nil
}

// Update updates an existing user
func (a *UserAdapter) Update(ctx context.Context, name string, update *types.UserUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	// Validate name
	if err := a.ValidateResourceName("user name", name); err != nil {
		return err
	}
	// Validate update
	if update == nil {
		return a.HandleValidationError("user update cannot be nil")
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}
	// Get the existing user first
	existingUser, err := a.Get(ctx, name)
	if err != nil {
		return err
	}
	// Apply updates
	if update.DefaultAccount != nil {
		if existingUser.Default == nil {
			existingUser.Default = &types.UserDefault{}
		}
		existingUser.Default.Account = update.DefaultAccount
	}
	if update.DefaultWCKey != nil {
		if existingUser.Default == nil {
			existingUser.Default = &types.UserDefault{}
		}
		existingUser.Default.Wckey = update.DefaultWCKey
	}
	if update.AdminLevel != nil {
		existingUser.AdministratorLevel = []types.AdministratorLevelValue{types.AdministratorLevelValue(*update.AdminLevel)}
	}
	// Convert to API request
	updateReq := a.convertCommonToAPIUser(existingUser)
	// Make the API call
	resp, err := a.client.SlurmdbV0041PostUsersWithResponse(ctx, *updateReq)
	if err != nil {
		return a.WrapError(err, "failed to update user "+name)
	}
	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
	}
	return nil
}

// Delete deletes a user
func (a *UserAdapter) Delete(ctx context.Context, name string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	// Validate name
	if err := a.ValidateResourceName("user name", name); err != nil {
		return err
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}
	// Make the API call
	resp, err := a.client.SlurmdbV0041DeleteUserWithResponse(ctx, name)
	if err != nil {
		return a.WrapError(err, "failed to delete user "+name)
	}
	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
	}
	return nil
}

// GetAssociations gets associations for a user
func (a *UserAdapter) GetAssociations(ctx context.Context, name string) (*types.AssociationList, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	// Validate name
	if err := a.ValidateResourceName("user name", name); err != nil {
		return nil, err
	}
	// v0.0.41 doesn't have a direct method to get associations for a specific user
	// We would need to use the association manager instead
	return nil, fmt.Errorf("getting associations for a specific user is not directly supported in API v0.0.41, use the association manager instead")
}

// AddToAccount adds a user to an account
func (a *UserAdapter) AddToAccount(ctx context.Context, userName string, accountName string, opts *types.UserAccountOptions) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	// Validate names
	if err := a.ValidateResourceName("user name", userName); err != nil {
		return err
	}
	if err := a.ValidateResourceName("account name", accountName); err != nil {
		return err
	}
	// v0.0.41 user association management is complex and involves undefined API types
	// Return not implemented for now
	return errors.NewNotImplementedError("user account association", "v0.0.41")
}

// RemoveFromAccount removes a user from an account
func (a *UserAdapter) RemoveFromAccount(ctx context.Context, userName string, accountName string) error {
	// v0.0.41 doesn't have a direct method to remove a user from an account
	// This would need to be done through the association manager by deleting the association
	return fmt.Errorf("removing a user from an account is not directly supported in API v0.0.41, use the association manager instead")
}

// GetWCKeys gets work charge keys for a user
func (a *UserAdapter) GetWCKeys(ctx context.Context, name string) ([]types.WCKey, error) {
	// Get the user with WCKeys
	user, err := a.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	// Return the Wckeys field directly (already []types.WCKey)
	if user.Wckeys == nil {
		return []types.WCKey{}, nil
	}
	return user.Wckeys, nil
}

// SetCoordinatorStatus sets whether a user is a coordinator for accounts
func (a *UserAdapter) SetCoordinatorStatus(ctx context.Context, name string, accounts []string, isCoordinator bool) error {
	// v0.0.41 doesn't support setting coordinator status through the API
	return errors.NewNotImplementedError("Set User Coordinator Status", "v0.0.41")
}

// Helper function to create string pointer
func ptr(s string) *string {
	return &s
}

// CreateAssociation creates associations for users (not supported in v0.0.41)
func (a *UserAdapter) CreateAssociation(ctx context.Context, req *types.UserAssociationRequest) (*types.AssociationCreateResponse, error) {
	return nil, a.HandleNotImplemented("CreateAssociation", "v0.0.41")
}
