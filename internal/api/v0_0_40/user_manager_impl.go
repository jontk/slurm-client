// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"context"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// UserManagerImpl implements the UserManager interface for v0.0.40
type UserManagerImpl struct {
	client *WrapperClient
}

// NewUserManagerImpl creates a new UserManagerImpl
func NewUserManagerImpl(client *WrapperClient) *UserManagerImpl {
	return &UserManagerImpl{
		client: client,
	}
}

// List retrieves a list of users with optional filtering
func (u *UserManagerImpl) List(ctx context.Context, opts *interfaces.ListUsersOptions) (*interfaces.UserList, error) {
	// Check if client is initialized
	if u.client == nil || u.client.apiClient == nil || u.client.apiClient.ClientInterface == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if err := u.client.CheckContext(ctx); err != nil {
		return nil, err
	}

	// Prepare parameters
	params := &SlurmdbV0040GetUsersParams{}
	if opts != nil {
		// v0.0.40 API doesn't support name filtering in params
		// We'll need to filter client-side
		if len(opts.Names) > 0 {
			// Store names for client-side filtering
			_ = opts.Names
		}
		// WithAssociations, WithDeleted, WithCoordinators fields don't exist in ListUsersOptions
		// Instead we have ActiveOnly, CoordinatorsOnly, etc.
		if opts.CoordinatorsOnly {
			// Convert CoordinatorsOnly to WithCoords
			coords := "true"
			params.WithCoords = &coords
		}
		if len(opts.AdminLevels) > 0 {
			// Set admin level filter
			adminLevel := opts.AdminLevels[0]
			if adminLevel == "Administrator" {
				level := Administrator
				params.AdminLevel = &level
			}
		}
	}

	// Make API call
	resp, err := u.client.apiClient.SlurmdbV0040GetUsersWithResponse(ctx, params)
	if err != nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "failed to list users")
	}

	// Check response
	if resp.StatusCode() != 200 {
		return nil, u.client.HandleErrorResponse(resp.StatusCode(), resp.Body)
	}

	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "received nil response")
	}

	// Convert response
	userList := &interfaces.UserList{
		Users: make([]interfaces.User, 0),
	}

	// resp.JSON200.Users is V0040UserList which is []V0040User, not *[]V0040User
	for _, usr := range resp.JSON200.Users {
		user := u.convertV0040UserToInterface(usr)
		userList.Users = append(userList.Users, *user)
	}

	return userList, nil
}

// Get retrieves a specific user by name
func (u *UserManagerImpl) Get(ctx context.Context, userName string) (*interfaces.User, error) {
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	// Check if client is initialized
	if u.client == nil || u.client.apiClient == nil || u.client.apiClient.ClientInterface == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if err := u.client.CheckContext(ctx); err != nil {
		return nil, err
	}

	// Make API call
	params := &SlurmdbV0040GetUserParams{}
	resp, err := u.client.apiClient.SlurmdbV0040GetUserWithResponse(ctx, userName, params)
	if err != nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "failed to get user")
	}

	// Check response
	if resp.StatusCode() != 200 {
		return nil, u.client.HandleErrorResponse(resp.StatusCode(), resp.Body)
	}

	if resp.JSON200 == nil || len(resp.JSON200.Users) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "user not found")
	}

	// Convert the first user
	return u.convertV0040UserToInterface(resp.JSON200.Users[0]), nil
}

// GetUserAccounts retrieves all account associations for a user
func (u *UserManagerImpl) GetUserAccounts(ctx context.Context, userName string) ([]*interfaces.UserAccount, error) {
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}
	return nil, errors.NewNotImplementedError("user accounts retrieval", "v0.0.40")
}

// GetUserQuotas retrieves quota information for a user
func (u *UserManagerImpl) GetUserQuotas(ctx context.Context, userName string) (*interfaces.UserQuota, error) {
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}
	return nil, errors.NewNotImplementedError("user quotas retrieval", "v0.0.40")
}

// GetUserDefaultAccount retrieves the default account for a user
func (u *UserManagerImpl) GetUserDefaultAccount(ctx context.Context, userName string) (*interfaces.Account, error) {
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}
	return nil, errors.NewNotImplementedError("user default account retrieval", "v0.0.40")
}

// GetUserFairShare retrieves fair-share information for a user
func (u *UserManagerImpl) GetUserFairShare(ctx context.Context, userName string) (*interfaces.UserFairShare, error) {
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}
	return nil, errors.NewNotImplementedError("user fair-share retrieval", "v0.0.40")
}

// CalculateJobPriority calculates job priority for a user and job submission
func (u *UserManagerImpl) CalculateJobPriority(ctx context.Context, userName string, jobSubmission *interfaces.JobSubmission) (*interfaces.JobPriorityInfo, error) {
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}
	if jobSubmission == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "job submission data is required", "jobSubmission", jobSubmission, nil)
	}
	return nil, errors.NewNotImplementedError("job priority calculation", "v0.0.40")
}

// ValidateUserAccountAccess validates user access to a specific account
func (u *UserManagerImpl) ValidateUserAccountAccess(ctx context.Context, userName, accountName string) (*interfaces.UserAccessValidation, error) {
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}
	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}
	return nil, errors.NewNotImplementedError("user-account access validation", "v0.0.40")
}

// GetUserAccountAssociations retrieves detailed user account associations
func (u *UserManagerImpl) GetUserAccountAssociations(ctx context.Context, userName string, opts *interfaces.ListUserAccountAssociationsOptions) ([]*interfaces.UserAccountAssociation, error) {
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}
	return nil, errors.NewNotImplementedError("user account associations retrieval", "v0.0.40")
}

// GetBulkUserAccounts retrieves accounts for multiple users in a single call
func (u *UserManagerImpl) GetBulkUserAccounts(ctx context.Context, userNames []string) (map[string][]*interfaces.UserAccount, error) {
	if len(userNames) == 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "at least one user name is required", "userNames", userNames, nil)
	}
	return nil, errors.NewNotImplementedError("bulk user accounts retrieval", "v0.0.40")
}

// GetBulkAccountUsers retrieves users for multiple accounts in a single call
func (u *UserManagerImpl) GetBulkAccountUsers(ctx context.Context, accountNames []string) (map[string][]*interfaces.UserAccountAssociation, error) {
	if len(accountNames) == 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "at least one account name is required", "accountNames", accountNames, nil)
	}
	return nil, errors.NewNotImplementedError("bulk account users retrieval", "v0.0.40")
}

// convertV0040UserToInterface converts v0.0.40 user to interface format
func (u *UserManagerImpl) convertV0040UserToInterface(usr V0040User) *interfaces.User {
	user := &interfaces.User{}

	// usr.Name is string, not *string
	user.Name = usr.Name
	// Uid field doesn't exist in V0040User - skip UID setting
	// user.UID = 0
	// DefaultAccount is in usr.Default.Account
	if usr.Default != nil && usr.Default.Account != nil {
		user.DefaultAccount = *usr.Default.Account
	}
	// DefaultWCKey is in usr.Default.Wckey
	if usr.Default != nil && usr.Default.Wckey != nil {
		user.DefaultWCKey = *usr.Default.Wckey
	}

	// Convert admin level - AdministratorLevel is *V0040AdminLvl which is []string
	if usr.AdministratorLevel != nil && len(*usr.AdministratorLevel) > 0 {
		adminLevels := *usr.AdministratorLevel
		// Take the first admin level
		user.AdminLevel = adminLevels[0]
	}

	return user
}

// Create creates a new user
func (u *UserManagerImpl) Create(ctx context.Context, user *interfaces.UserCreate) (*interfaces.UserCreateResponse, error) {
	return nil, errors.NewNotImplementedError("Create", "v0.0.40")
}

// Update updates a user
func (u *UserManagerImpl) Update(ctx context.Context, userName string, update *interfaces.UserUpdate) error {
	return errors.NewNotImplementedError("Update", "v0.0.40")
}

// Delete deletes a user
func (u *UserManagerImpl) Delete(ctx context.Context, userName string) error {
	return errors.NewNotImplementedError("Delete", "v0.0.40")
}

// CreateAssociation creates a user-account association
func (u *UserManagerImpl) CreateAssociation(ctx context.Context, accountName string, opts *interfaces.AssociationOptions) (*interfaces.AssociationCreateResponse, error) {
	return nil, errors.NewNotImplementedError("CreateAssociation", "v0.0.40")
}
