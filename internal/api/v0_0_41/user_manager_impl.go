// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

import (
	"context"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// UserManagerImpl implements the UserManager interface for v0.0.41
// Most user management features are not available in v0.0.41
type UserManagerImpl struct {
	client *WrapperClient
}

// NewUserManagerImpl creates a new UserManagerImpl
func NewUserManagerImpl(client *WrapperClient) *UserManagerImpl {
	return &UserManagerImpl{
		client: client,
	}
}

// List lists users with optional filtering
// This feature is not available in v0.0.41
func (u *UserManagerImpl) List(ctx context.Context, opts *interfaces.ListUsersOptions) (*interfaces.UserList, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// v0.0.41 has very limited user management support
	// Most user features are not available
	return nil, errors.NewNotImplementedError("user listing not supported", "v0.0.41")
}

// Get retrieves a specific user by name
// This feature is not available in v0.0.41
func (u *UserManagerImpl) Get(ctx context.Context, userName string) (*interfaces.User, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	// v0.0.41 has minimal user retrieval support
	return nil, errors.NewNotImplementedError("user retrieval not supported", "v0.0.41")
}

// GetUserAccounts retrieves all account associations for a user
// This feature is not available in v0.0.41
func (u *UserManagerImpl) GetUserAccounts(ctx context.Context, userName string) ([]*interfaces.UserAccount, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	return nil, errors.NewNotImplementedError("user accounts not supported", "v0.0.41")
}

// GetUserQuotas retrieves quota information for a user
// This feature is not available in v0.0.41
func (u *UserManagerImpl) GetUserQuotas(ctx context.Context, userName string) (*interfaces.UserQuota, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	return nil, errors.NewNotImplementedError("user quotas not supported", "v0.0.41")
}

// GetUserDefaultAccount retrieves the default account for a user
// This feature is not available in v0.0.41
func (u *UserManagerImpl) GetUserDefaultAccount(ctx context.Context, userName string) (*interfaces.Account, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	return nil, errors.NewNotImplementedError("user default account not supported", "v0.0.41")
}

// GetUserFairShare retrieves fair-share information for a user
// This feature is not available in v0.0.41
func (u *UserManagerImpl) GetUserFairShare(ctx context.Context, userName string) (*interfaces.UserFairShare, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	return nil, errors.NewNotImplementedError("user fair-share not supported", "v0.0.41")
}

// CalculateJobPriority calculates job priority for a user and job submission
// This feature is not available in v0.0.41
func (u *UserManagerImpl) CalculateJobPriority(ctx context.Context, userName string, jobSubmission *interfaces.JobSubmission) (*interfaces.JobPriorityInfo, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	if jobSubmission == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "job submission data is required", "jobSubmission", jobSubmission, nil)
	}

	return nil, errors.NewNotImplementedError("job priority calculation not supported", "v0.0.41")
}

// ValidateUserAccountAccess validates user access to a specific account
// This feature is not available in v0.0.41
func (u *UserManagerImpl) ValidateUserAccountAccess(ctx context.Context, userName, accountName string) (*interfaces.UserAccessValidation, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	return nil, errors.NewNotImplementedError("user-account access validation not supported", "v0.0.41")
}

// GetUserAccountAssociations retrieves detailed user account associations
// This feature is not available in v0.0.41
func (u *UserManagerImpl) GetUserAccountAssociations(ctx context.Context, userName string, opts *interfaces.ListUserAccountAssociationsOptions) ([]*interfaces.UserAccountAssociation, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	return nil, errors.NewNotImplementedError("user account associations not supported", "v0.0.41")
}

// GetBulkUserAccounts retrieves accounts for multiple users in a single call
// This feature is not available in v0.0.41
func (u *UserManagerImpl) GetBulkUserAccounts(ctx context.Context, userNames []string) (map[string][]*interfaces.UserAccount, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if len(userNames) == 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "at least one user name is required", "userNames", userNames, nil)
	}

	return nil, errors.NewNotImplementedError("bulk user accounts not supported", "v0.0.41")
}

// GetBulkAccountUsers retrieves users for multiple accounts in a single call
// This feature is not available in v0.0.41
func (u *UserManagerImpl) GetBulkAccountUsers(ctx context.Context, accountNames []string) (map[string][]*interfaces.UserAccountAssociation, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if len(accountNames) == 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "at least one account name is required", "accountNames", accountNames, nil)
	}

	return nil, errors.NewNotImplementedError("bulk account users not supported", "v0.0.41")
}

// Create creates a new user
func (u *UserManagerImpl) Create(ctx context.Context, user *interfaces.UserCreate) (*interfaces.UserCreateResponse, error) {
	return nil, errors.NewNotImplementedError("Create", "v0.0.41")
}

// Update updates a user
func (u *UserManagerImpl) Update(ctx context.Context, userName string, update *interfaces.UserUpdate) error {
	return errors.NewNotImplementedError("Update", "v0.0.41")
}

// Delete deletes a user
func (u *UserManagerImpl) Delete(ctx context.Context, userName string) error {
	return errors.NewNotImplementedError("Delete", "v0.0.41")
}

// CreateAssociation creates a user-account association
func (u *UserManagerImpl) CreateAssociation(ctx context.Context, accountName string, opts *interfaces.AssociationOptions) (*interfaces.AssociationCreateResponse, error) {
	return nil, errors.NewNotImplementedError("CreateAssociation", "v0.0.41")
}
