// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"context"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
)

// UserAdapter implements the UserAdapter interface for v0.0.40
type UserAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewUserAdapter creates a new User adapter for v0.0.40
func NewUserAdapter(client *api.ClientWithResponses) *UserAdapter {
	return &UserAdapter{
		BaseManager: base.NewBaseManager("v0.0.40", "User"),
		client:      client,
		wrapper:     nil, // We'll implement this later
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
	params := &api.SlurmdbV0040GetUsersParams{}

	// Apply filters from options
	if opts != nil {
		// v0.0.40 doesn't support user name filtering in params
		// We'll need to filter client-side
		if opts.DefaultAccount != "" {
			params.DefaultAccount = &opts.DefaultAccount
		}
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
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0040GetUsersWithResponse(ctx, params)
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
	if err := a.CheckNilResponse(resp.JSON200, "List Users"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Users, "List Users - users field"); err != nil {
		return nil, err
	}

	// Convert the response to common types
	userList := make([]types.User, 0, len(resp.JSON200.Users))
	for _, apiUser := range resp.JSON200.Users {
		user, err := a.convertAPIUserToCommon(apiUser)
		if err != nil {
			return nil, a.HandleConversionError(err, apiUser.Name)
		}
		userList = append(userList, *user)
	}

	// Apply client-side filtering if needed
	if opts != nil {
		userList = a.filterUserList(userList, opts)
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
	if start >= len(userList) {
		return &types.UserList{
			Users: []types.User{},
			Total: len(userList),
		}, nil
	}

	end := len(userList)
	if listOpts.Limit > 0 {
		end = start + listOpts.Limit
		if end > len(userList) {
			end = len(userList)
		}
	}

	return &types.UserList{
		Users: userList[start:end],
		Total: len(userList),
	}, nil
}

// Get retrieves a specific user by name
func (a *UserAdapter) Get(ctx context.Context, userName string) (*types.User, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.ValidateResourceName(userName, "userName"); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmdbV0040GetUserParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0040GetUserWithResponse(ctx, userName, params)
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
	if err := a.CheckNilResponse(resp.JSON200, "Get User"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Users, "Get User - users field"); err != nil {
		return nil, err
	}

	// Check if we got any user entries
	if len(resp.JSON200.Users) == 0 {
		return nil, common.NewResourceNotFoundError("User", userName)
	}

	// Convert the first user (should be the only one)
	user, err := a.convertAPIUserToCommon(resp.JSON200.Users[0])
	if err != nil {
		return nil, a.HandleConversionError(err, userName)
	}

	return user, nil
}

// Create creates a new user
func (a *UserAdapter) Create(ctx context.Context, user *types.UserCreate) (*types.UserCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.validateUserCreate(user); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert to API format
	apiUser, err := a.convertCommonUserCreateToAPI(user)
	if err != nil {
		return nil, err
	}

	// Create request body
	reqBody := api.SlurmdbV0040PostUsersJSONRequestBody{
		Users: []api.V0040User{*apiUser},
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0040PostUsersWithResponse(ctx, reqBody)
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

	// Create response object with created user information
	createResponse := &types.UserCreateResponse{
		UserName: user.Name,
		// Add any additional response fields from API if available
	}

	return createResponse, nil
}

// Update updates an existing user
func (a *UserAdapter) Update(ctx context.Context, userName string, update *types.UserUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(userName, "userName"); err != nil {
		return err
	}
	if err := a.validateUserUpdate(update); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// First, get the existing user to merge updates
	existingUser, err := a.Get(ctx, userName)
	if err != nil {
		return err
	}

	// Convert to API format and apply updates
	apiUser, err := a.convertCommonUserUpdateToAPI(existingUser, update)
	if err != nil {
		return err
	}

	// Create request body
	reqBody := api.SlurmdbV0040PostUsersJSONRequestBody{
		Users: []api.V0040User{*apiUser},
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0040PostUsersWithResponse(ctx, reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "v0.0.40")
}

// Delete deletes a user
func (a *UserAdapter) Delete(ctx context.Context, userName string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(userName, "userName"); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0040DeleteUserWithResponse(ctx, userName)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	// Create adapter with special handling for 204 (No Content) status
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)

	// For DELETE operations, 204 is also a success
	if resp.StatusCode() == 204 {
		return nil
	}

	return common.HandleAPIResponse(responseAdapter, "v0.0.40")
}

// filterUserList applies client-side filtering to the user list
func (a *UserAdapter) filterUserList(users []types.User, opts *types.UserListOptions) []types.User {
	filtered := make([]types.User, 0, len(users))
	
	for _, user := range users {
		// Apply DefaultAccount filter
		if opts.DefaultAccount != "" && user.DefaultAccount != opts.DefaultAccount {
			continue
		}

		// v0.0.40 doesn't support DefaultWCKey or AdminLevel filters in options
		// Skip these filters

		filtered = append(filtered, user)
	}

	return filtered
}

// validateUserCreate validates user creation request
func (a *UserAdapter) validateUserCreate(user *types.UserCreate) error {
	if user == nil {
		return common.NewValidationError("user creation data is required", "user", nil)
	}
	if user.Name == "" {
		return common.NewValidationError("user name is required", "name", user.Name)
	}
	return nil
}

// validateUserUpdate validates user update request
func (a *UserAdapter) validateUserUpdate(update *types.UserUpdate) error {
	if update == nil {
		return common.NewValidationError("user update data is required", "update", nil)
	}
	// At least one field should be provided for update
	if update.DefaultAccount == nil && update.DefaultWCKey == nil && update.AdminLevel == nil {
		return common.NewValidationError("at least one field must be provided for update", "update", update)
	}
	return nil
}
