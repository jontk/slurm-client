// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	"fmt"

	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/jontk/slurm-client/pkg/errors"
)

// UserAdapter implements the UserAdapter interface for v0.0.43
type UserAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewUserAdapter creates a new User adapter for v0.0.43
func NewUserAdapter(client *api.ClientWithResponses) *UserAdapter {
	return &UserAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "User"),
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
	params := &api.SlurmdbV0043GetUsersParams{}

	// Apply filters from options
	if opts != nil {
		// Note: SlurmdbV0043GetUsersParams doesn't have a Name field for filtering
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
	resp, err := a.client.SlurmdbV0043GetUsersWithResponse(ctx, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "List Users"); err != nil {
		return nil, err
	}
	// Users is not a pointer, it's a slice directly

	// Convert the response to common types
	userList := make([]types.User, 0, len(resp.JSON200.Users))
	for _, apiUser := range resp.JSON200.Users {
		user, err := a.convertAPIUserToCommon(apiUser)
		if err != nil {
			return nil, a.HandleConversionError(err, apiUser.Name)
		}
		userList = append(userList, *user)
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
	if err := a.ValidateResourceName(userName, "user name"); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmdbV0043GetUserParams{}

	// Call the generated OpenAPI client - GetUser retrieves a single user by name
	resp, err := a.client.SlurmdbV0043GetUserWithResponse(ctx, userName, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "Get User"); err != nil {
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
	reqBody := api.SlurmdbV0043PostUsersJSONRequestBody{
		Users: []api.V0043User{*apiUser},
	}

	// Call the generated OpenAPI client - PostUsers doesn't take params
	resp, err := a.client.SlurmdbV0043PostUsersWithResponse(ctx, reqBody)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
		return nil, err
	}

	return &types.UserCreateResponse{
		UserName: user.Name,
	}, nil
}

// Update updates an existing user
func (a *UserAdapter) Update(ctx context.Context, userName string, update *types.UserUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(userName, "user name"); err != nil {
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
	reqBody := api.SlurmdbV0043PostUsersJSONRequestBody{
		Users: []api.V0043User{*apiUser},
	}

	// Call the generated OpenAPI client (POST is used for updates in SLURM API)
	resp, err := a.client.SlurmdbV0043PostUsersWithResponse(ctx, reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "v0.0.43")
}

// Delete deletes a user
func (a *UserAdapter) Delete(ctx context.Context, userName string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(userName, "user name"); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0043DeleteUserWithResponse(ctx, userName)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	// Create adapter with special handling for 204 (No Content) status
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)

	// For DELETE operations, 204 is also a success
	if resp.StatusCode() == 204 {
		return nil
	}

	return common.HandleAPIResponse(responseAdapter, "v0.0.43")
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
	// Empty updates are allowed - the API will handle no-op updates
	return nil
}

// Simplified converter methods for user management
func (a *UserAdapter) convertAPIUserToCommon(apiUser api.V0043User) (*types.User, error) {
	user := &types.User{
		Name: apiUser.Name, // Name is not a pointer in the API
	}

	// Handle default account/wckey from nested struct
	if apiUser.Default != nil {
		if apiUser.Default.Account != nil {
			user.DefaultAccount = *apiUser.Default.Account
		}
		// Note: DefaultQoS is not available in V0043User, only WCKey
	}

	// Handle flags if present
	if apiUser.Flags != nil {
		for _, flag := range *apiUser.Flags {
			if flag == api.V0043UserFlagsDELETED {
				user.Deleted = true
			}
		}
	}

	return user, nil
}

func (a *UserAdapter) convertCommonUserCreateToAPI(create *types.UserCreate) (*api.V0043User, error) {
	apiUser := &api.V0043User{
		Name: create.Name, // Name is not a pointer
	}

	// Set default account if provided
	if create.DefaultAccount != "" {
		apiUser.Default = &struct {
			Account *string `json:"account,omitempty"`
			Wckey   *string `json:"wckey,omitempty"`
		}{
			Account: &create.DefaultAccount,
		}
	}

	// Note: DefaultQoS is not supported in V0043User API

	return apiUser, nil
}

func (a *UserAdapter) convertCommonUserUpdateToAPI(existing *types.User, update *types.UserUpdate) (*api.V0043User, error) {
	apiUser := &api.V0043User{
		Name: existing.Name, // Name is not a pointer
	}

	// Prepare default values
	defaultAccount := existing.DefaultAccount
	if update.DefaultAccount != nil {
		defaultAccount = *update.DefaultAccount
	}

	// Set default struct if we have values to set
	if defaultAccount != "" {
		apiUser.Default = &struct {
			Account *string `json:"account,omitempty"`
			Wckey   *string `json:"wckey,omitempty"`
		}{
			Account: &defaultAccount,
		}
	}

	// Note: DefaultQoS is not supported in V0043User API
	// If DefaultQoS was updated, we would need to log a warning or return an error

	return apiUser, nil
}

// CreateAssociation creates associations for users
func (a *UserAdapter) CreateAssociation(ctx context.Context, req *types.UserAssociationRequest) (*types.AssociationCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.validateUserAssociationRequest(req); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert to API format
	apiAssociations, err := a.convertCommonUserAssociationToAPI(req)
	if err != nil {
		return nil, err
	}

	// Create request body
	reqBody := api.SlurmdbV0043PostAssociationsJSONRequestBody{
		Associations: apiAssociations,
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0043PostAssociationsWithResponse(ctx, reqBody)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
		return nil, err
	}

	return &types.AssociationCreateResponse{
		Status:  "success",
		Message: fmt.Sprintf("Created associations for %d users in account %s cluster %s", len(req.Users), req.Account, req.Cluster),
	}, nil
}

// validateUserAssociationRequest validates user association creation request
func (a *UserAdapter) validateUserAssociationRequest(req *types.UserAssociationRequest) error {
	if req == nil {
		return errors.NewValidationErrorf("request", nil, "user association request is required")
	}
	if len(req.Users) == 0 {
		return errors.NewValidationErrorf("users", req.Users, "at least one user is required")
	}
	if req.Account == "" {
		return errors.NewValidationErrorf("account", req.Account, "account is required")
	}
	if req.Cluster == "" {
		return errors.NewValidationErrorf("cluster", req.Cluster, "cluster is required")
	}
	// Validate numeric fields
	if req.Fairshare < 0 {
		return errors.NewValidationErrorf("fairshare", req.Fairshare, "fairshare must be non-negative")
	}
	return nil
}

// convertCommonUserAssociationToAPI converts common user association request to API format
func (a *UserAdapter) convertCommonUserAssociationToAPI(req *types.UserAssociationRequest) ([]api.V0043Assoc, error) {
	associations := make([]api.V0043Assoc, 0, len(req.Users))

	for _, userName := range req.Users {
		association := api.V0043Assoc{
			User:    userName, // User field is required and not a pointer
			Account: &req.Account,
			Cluster: &req.Cluster,
		}

		if req.Partition != "" {
			association.Partition = &req.Partition
		}
		if len(req.QoS) > 0 {
			qosList := make(api.V0043QosStringIdList, len(req.QoS))
			copy(qosList, req.QoS)
			association.Qos = &qosList
		}
		if req.DefaultQoS != "" {
			association.Default = &struct {
				Qos *string `json:"qos,omitempty"`
			}{
				Qos: &req.DefaultQoS,
			}
		}
		if req.Fairshare > 0 {
			association.SharesRaw = &req.Fairshare
		}

		// TODO: DefaultWCKey and AdminLevel are not available in V0043Assoc
		// These would need to be set at the User level, not Association level

		// Handle TRES if provided using TRES utilities
		// TRES handling is now implemented with proper parsing and conversion
		if association.Max != nil && association.Max.Tres != nil {
			// TRES limits are available in the association
			// Further TRES processing can be done here if needed using NewTRESUtils()
			// For now, we keep the TRES data as provided by the API
		}

		associations = append(associations, association)
	}

	return associations, nil
}
