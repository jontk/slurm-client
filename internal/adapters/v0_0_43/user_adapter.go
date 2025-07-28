package v0_0_43

import (
	"context"
	"strings"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
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
		if len(opts.Names) > 0 {
			nameStr := strings.Join(opts.Names, ",")
			params.Name = &nameStr
		}
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
	if err := a.CheckNilResponse(resp.JSON200.Users, "List Users - users field"); err != nil {
		return nil, err
	}

	// Convert the response to common types
	userList := make([]types.User, 0, len(*resp.JSON200.Users))
	for _, apiUser := range *resp.JSON200.Users {
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
	if err := a.ValidateResourceName(userName, "userName"); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmdbV0043GetSingleUserParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0043GetSingleUserWithResponse(ctx, userName, params)
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
	if err := a.CheckNilResponse(resp.JSON200.Users, "Get User - users field"); err != nil {
		return nil, err
	}

	// Check if we got any user entries
	if len(*resp.JSON200.Users) == 0 {
		return nil, common.NewResourceNotFoundError("User", userName)
	}

	// Convert the first user (should be the only one)
	user, err := a.convertAPIUserToCommon((*resp.JSON200.Users)[0])
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
		Users: []api.V0043UserInfo{*apiUser},
	}

	// Prepare parameters for the API call
	params := &api.SlurmdbV0043PostUsersParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0043PostUsersWithResponse(ctx, params, reqBody)
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
	reqBody := api.SlurmdbV0043PostUsersJSONRequestBody{
		Users: []api.V0043UserInfo{*apiUser},
	}

	// Prepare parameters for the API call
	params := &api.SlurmdbV0043PostUsersParams{}

	// Call the generated OpenAPI client (POST is used for updates in SLURM API)
	resp, err := a.client.SlurmdbV0043PostUsersWithResponse(ctx, params, reqBody)
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
	if err := a.ValidateResourceName(userName, "userName"); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0043DeleteSingleUserWithResponse(ctx, userName)
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
	// At least one field should be provided for update
	if update.DefaultAccount == nil && update.DefaultQoS == nil && len(update.QoSList) == 0 &&
	   len(update.Accounts) == 0 && update.AdminLevel == nil && update.MaxJobs == nil {
		return common.NewValidationError("at least one field must be provided for update", "update", update)
	}
	return nil
}

// Simplified converter methods for user management
func (a *UserAdapter) convertAPIUserToCommon(apiUser api.V0043UserInfo) (*types.User, error) {
	user := &types.User{}
	if apiUser.Name != nil {
		user.Name = *apiUser.Name
	}
	if apiUser.DefaultAccount != nil {
		user.DefaultAccount = *apiUser.DefaultAccount
	}
	if apiUser.DefaultQos != nil {
		user.DefaultQoS = *apiUser.DefaultQos
	}
	if apiUser.Flags != nil {
		for _, flag := range *apiUser.Flags {
			if flag == api.V0043UserInfoFlagsDELETED {
				user.Deleted = true
			}
		}
	}
	return user, nil
}

func (a *UserAdapter) convertCommonUserCreateToAPI(create *types.UserCreate) (*api.V0043UserInfo, error) {
	apiUser := &api.V0043UserInfo{}
	apiUser.Name = &create.Name
	if create.DefaultAccount != "" {
		apiUser.DefaultAccount = &create.DefaultAccount
	}
	if create.DefaultQoS != "" {
		apiUser.DefaultQos = &create.DefaultQoS
	}
	return apiUser, nil
}

func (a *UserAdapter) convertCommonUserUpdateToAPI(existing *types.User, update *types.UserUpdate) (*api.V0043UserInfo, error) {
	apiUser := &api.V0043UserInfo{}
	apiUser.Name = &existing.Name

	defaultAccount := existing.DefaultAccount
	if update.DefaultAccount != nil {
		defaultAccount = *update.DefaultAccount
	}
	if defaultAccount != "" {
		apiUser.DefaultAccount = &defaultAccount
	}

	defaultQoS := existing.DefaultQoS
	if update.DefaultQoS != nil {
		defaultQoS = *update.DefaultQoS
	}
	if defaultQoS != "" {
		apiUser.DefaultQos = &defaultQoS
	}

	return apiUser, nil
}