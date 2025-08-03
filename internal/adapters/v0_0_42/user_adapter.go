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

// CreateAssociation creates associations for users
func (a *UserAdapter) CreateAssociation(ctx context.Context, req *types.UserAssociationRequest) (*types.AssociationCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Validate request
	if req == nil {
		return nil, fmt.Errorf("association request is required")
	}
	if len(req.Users) == 0 {
		return nil, fmt.Errorf("at least one user is required")
	}
	if req.Account == "" {
		return nil, fmt.Errorf("account is required")
	}
	if req.Cluster == "" {
		return nil, fmt.Errorf("cluster is required")
	}

	// Convert common request to API request structure
	apiReq, err := a.convertUserAssociationRequestToAPI(req)
	if err != nil {
		return nil, a.WrapError(err, "failed to convert association request")
	}

	// Prepare parameters (optional flags)
	params := &api.SlurmdbV0042PostUsersAssociationParams{}

	// Call the API
	resp, err := a.client.SlurmdbV0042PostUsersAssociationWithResponse(ctx, params, *apiReq)
	if err != nil {
		return nil, a.WrapError(err, "failed to create user associations")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}

	// Check for API response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from API")
	}

	// Convert response
	return a.convertUserAssociationResponseToCommon(resp.JSON200), nil
}

// convertUserAssociationRequestToAPI converts common request to API structure
func (a *UserAdapter) convertUserAssociationRequestToAPI(req *types.UserAssociationRequest) (*api.V0042OpenapiUsersAddCondResp, error) {
	// Create users list from string slice
	users := make(api.V0042StringList, len(req.Users))
	for i, user := range req.Users {
		users[i] = user
	}

	// Create association condition
	assocCond := api.V0042UsersAddCond{
		Users: users,
	}

	// Add accounts
	if req.Account != "" {
		accounts := api.V0042StringList{req.Account}
		assocCond.Accounts = &accounts
	}

	// Add clusters
	if req.Cluster != "" {
		clusters := api.V0042StringList{req.Cluster}
		assocCond.Clusters = &clusters
	}

	// Add partitions if specified
	if req.Partition != "" {
		partitions := api.V0042StringList{req.Partition}
		assocCond.Partitions = &partitions
	}

	// Create association record set if we have additional fields
	if req.DefaultQoS != "" || req.Fairshare != 0 || len(req.QoS) > 0 || req.DefaultWCKey != "" {
		assocRec := &api.V0042AssocRecSet{}

		// Set default QoS (note: field name is "Defaultqos" in v0.0.42)
		if req.DefaultQoS != "" {
			assocRec.Defaultqos = &req.DefaultQoS
		}

		// Set fairshare (note: field name is "Fairshare" in v0.0.42)
		if req.Fairshare != 0 {
			fairshareInt32 := int32(req.Fairshare)
			assocRec.Fairshare = &fairshareInt32
		}

		// Note: DefaultWCKey is set in UserShort, not AssocRecSet

		assocCond.Association = assocRec
	}

	// Create a minimal user short structure
	userShort := api.V0042UserShort{}

	// Set default WCKey (note: field name is "Defaultwckey" in v0.0.42)
	if req.DefaultWCKey != "" {
		userShort.Defaultwckey = &req.DefaultWCKey
	}

	// Set admin level if specified (note: field name is "Adminlevel" in v0.0.42, and it's an array)
	if req.AdminLevel != "" {
		adminLevel := api.V0042AdminLvl{req.AdminLevel}
		userShort.Adminlevel = &adminLevel
	}

	// Create the API request
	apiReq := &api.V0042OpenapiUsersAddCondResp{
		AssociationCondition: assocCond,
		User:                 userShort,
	}

	return apiReq, nil
}

// convertUserAssociationResponseToCommon converts API response to common type
func (a *UserAdapter) convertUserAssociationResponseToCommon(apiResp *api.V0042OpenapiUsersAddCondRespStr) *types.AssociationCreateResponse {
	resp := &types.AssociationCreateResponse{
		Status: "success",
		Meta:   make(map[string]interface{}),
	}

	// Extract added users info
	if apiResp.AddedUsers != "" {
		resp.Message = fmt.Sprintf("Successfully created associations for users: %s", apiResp.AddedUsers)
		resp.Meta["added_users"] = apiResp.AddedUsers
	} else {
		resp.Message = "User associations created successfully"
	}

	// Handle errors in response
	if apiResp.Errors != nil && len(*apiResp.Errors) > 0 {
		resp.Status = "error"
		errors := *apiResp.Errors
		if len(errors) > 0 && errors[0].Error != nil {
			resp.Message = *errors[0].Error
		} else {
			resp.Message = "User association creation failed"
		}
	}

	// Extract metadata if available
	if apiResp.Meta != nil {
		if apiResp.Meta.Client != nil {
			clientInfo := make(map[string]interface{})
			if apiResp.Meta.Client.Source != nil {
				clientInfo["source"] = *apiResp.Meta.Client.Source
			}
			if apiResp.Meta.Client.User != nil {
				clientInfo["user"] = *apiResp.Meta.Client.User
			}
			if len(clientInfo) > 0 {
				resp.Meta["client"] = clientInfo
			}
		}
	}

	return resp
}
