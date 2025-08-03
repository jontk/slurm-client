// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"fmt"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

// UserAdapter implements the UserAdapter interface for v0.0.42
type UserAdapter struct {
	*base.BaseManager
	client *api.ClientWithResponses
}

// NewUserAdapter creates a new User adapter for v0.0.42
func NewUserAdapter(client *api.ClientWithResponses) *UserAdapter {
	return &UserAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "User"),
		client:      client,
	}
}

// List retrieves a list of users
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
	params := &api.SlurmdbV0042GetUsersParams{}

	// Apply filters from options
	if opts != nil {
		// v0.0.42 doesn't support user name filtering in the API params
		// We'll need to filter client-side
		if opts.DefaultAccount != "" {
			params.DefaultAccount = &opts.DefaultAccount
		}
		if opts.WithAssocs {
			withAssoc := "true"
			params.WithAssocs = &withAssoc
		}
		if opts.WithCoords {
			withCoord := "true"
			params.WithCoords = &withCoord
		}
		if opts.WithDeleted {
			withDeleted := "true"
			params.WithDeleted = &withDeleted
		}
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042GetUsersWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list users")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleHTTPResponse(resp.HTTPResponse, resp.Body)
	}

	// Check for API response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from API")
	}

	// Convert the response to common types
	userList := &types.UserList{
		Users: make([]types.User, 0),
	}

	if resp.JSON200.Users != nil {
		for _, apiUser := range resp.JSON200.Users {
			user, err := a.convertAPIUserToCommon(apiUser)
			if err != nil {
				// Log conversion error but continue
				continue
			}
			userList.Users = append(userList.Users, *user)
		}
		
		// Apply client-side filtering if needed
		if opts != nil && len(opts.Names) > 0 {
			filtered := make([]types.User, 0, len(userList.Users))
			for _, user := range userList.Users {
				for _, name := range opts.Names {
					if user.Name == name {
						filtered = append(filtered, user)
						break
					}
				}
			}
			userList.Users = filtered
		}
	}

	return userList, nil
}

// Get retrieves a specific user by name
func (a *UserAdapter) Get(ctx context.Context, name string) (*types.User, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters
	params := &api.SlurmdbV0042GetUserParams{}
	withAssoc := "true"
	params.WithAssocs = &withAssoc

	// Call the API
	resp, err := a.client.SlurmdbV0042GetUserWithResponse(ctx, name, params)
	if err != nil {
		return nil, a.WrapError(err, fmt.Sprintf("failed to get user %s", name))
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleHTTPResponse(resp.HTTPResponse, resp.Body)
	}

	// Check for API response
	if resp.JSON200 == nil || len(resp.JSON200.Users) == 0 {
		return nil, fmt.Errorf("user %s not found", name)
	}

	// Find the specific user by name
	for _, apiUser := range resp.JSON200.Users {
		user, err := a.convertAPIUserToCommon(apiUser)
		if err != nil {
			continue
		}
		if user.Name == name {
			return user, nil
		}
	}

	return nil, fmt.Errorf("user %s not found", name)
}

// Create creates a new user
func (a *UserAdapter) Create(ctx context.Context, user *types.UserCreate) (*types.UserCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert common user to API format
	apiUser, err := a.convertCommonUserCreateToAPI(user)
	if err != nil {
		return nil, a.WrapError(err, "failed to convert user create request")
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042PostUsersWithResponse(ctx, *apiUser)
	if err != nil {
		return nil, a.WrapError(err, "failed to create user")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleHTTPResponse(resp.HTTPResponse, resp.Body)
	}

	// Return success response
	return &types.UserCreateResponse{
		UserName: user.Name,
	}, nil
}

// Update updates an existing user
func (a *UserAdapter) Update(ctx context.Context, name string, updates *types.UserUpdateRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// v0.0.42 doesn't have a direct user update endpoint
	// Updates are typically done through associations
	return fmt.Errorf("user update not directly supported via v0.0.42 API - use association updates")
}

// Delete deletes a user
func (a *UserAdapter) Delete(ctx context.Context, name string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042DeleteUserWithResponse(ctx, name)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to delete user %s", name))
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleHTTPResponse(resp.HTTPResponse, resp.Body)
	}

	return nil
}

// CreateAssociation creates associations for users (not supported in v0.0.42)
func (a *UserAdapter) CreateAssociation(ctx context.Context, req *types.UserAssociationRequest) (*types.AssociationCreateResponse, error) {
	return nil, fmt.Errorf("CreateAssociation not supported in API v0.0.42")
}
