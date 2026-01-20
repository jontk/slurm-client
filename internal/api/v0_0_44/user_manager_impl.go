// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_44

import (
	"context"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// UserManagerImpl provides the actual implementation for UserManager methods
type UserManagerImpl struct {
	client *WrapperClient
}

// NewUserManagerImpl creates a new UserManager implementation
func NewUserManagerImpl(client *WrapperClient) *UserManagerImpl {
	return &UserManagerImpl{client: client}
}

// List retrieves a list of users with optional filtering
func (m *UserManagerImpl) List(ctx context.Context, opts *interfaces.ListUsersOptions) (*interfaces.UserList, error) {
	return &interfaces.UserList{Users: make([]interfaces.User, 0)}, nil
}

// Get retrieves a specific user by name
func (m *UserManagerImpl) Get(ctx context.Context, userName string) (*interfaces.User, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "User management not yet implemented for v0.0.44")
}

// GetUserAccounts retrieves accounts associated with user
func (m *UserManagerImpl) GetUserAccounts(ctx context.Context, userName string) ([]*interfaces.UserAccount, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "User management not yet implemented for v0.0.44")
}

// GetUserQuotas retrieves user quotas
func (m *UserManagerImpl) GetUserQuotas(ctx context.Context, userName string) (*interfaces.UserQuota, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "User management not yet implemented for v0.0.44")
}

// GetUserDefaultAccount retrieves user default account
func (m *UserManagerImpl) GetUserDefaultAccount(ctx context.Context, userName string) (*interfaces.Account, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "User management not yet implemented for v0.0.44")
}

// GetUserFairShare retrieves user fair share information
func (m *UserManagerImpl) GetUserFairShare(ctx context.Context, userName string) (*interfaces.UserFairShare, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "User management not yet implemented for v0.0.44")
}

// CalculateJobPriority calculates job priority for user
func (m *UserManagerImpl) CalculateJobPriority(ctx context.Context, userName string, jobSubmission *interfaces.JobSubmission) (*interfaces.JobPriorityInfo, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "User management not yet implemented for v0.0.44")
}

// ValidateUserAccountAccess validates user access to account
func (m *UserManagerImpl) ValidateUserAccountAccess(ctx context.Context, userName, accountName string) (*interfaces.UserAccessValidation, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "User management not yet implemented for v0.0.44")
}

// GetUserAccountAssociations retrieves user account associations
func (m *UserManagerImpl) GetUserAccountAssociations(ctx context.Context, userName string, opts *interfaces.ListUserAccountAssociationsOptions) ([]*interfaces.UserAccountAssociation, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "User management not yet implemented for v0.0.44")
}

// GetBulkUserAccounts retrieves accounts for multiple users
func (m *UserManagerImpl) GetBulkUserAccounts(ctx context.Context, userNames []string) (map[string][]*interfaces.UserAccount, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "User management not yet implemented for v0.0.44")
}

// GetBulkAccountUsers retrieves users for multiple accounts
func (m *UserManagerImpl) GetBulkAccountUsers(ctx context.Context, accountNames []string) (map[string][]*interfaces.UserAccountAssociation, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "User management not yet implemented for v0.0.44")
}

// Create creates a new user
func (m *UserManagerImpl) Create(ctx context.Context, user *interfaces.UserCreate) (*interfaces.UserCreateResponse, error) {
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// For now, return not implemented - user creation might require database management endpoints
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "User creation not yet implemented for v0.0.44")
}

// Update updates an existing user
func (m *UserManagerImpl) Update(ctx context.Context, userName string, update *interfaces.UserUpdate) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// For now, return not implemented - user updates might require database management endpoints
	return errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "User updates not yet implemented for v0.0.44")
}

// Delete deletes a user
func (m *UserManagerImpl) Delete(ctx context.Context, userName string) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// For now, return not implemented - user deletion might require database management endpoints
	return errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "User deletion not yet implemented for v0.0.44")
}

// CreateAssociation creates a new user-account association
func (m *UserManagerImpl) CreateAssociation(ctx context.Context, accountName string, opts *interfaces.AssociationOptions) (*interfaces.AssociationCreateResponse, error) {
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// For now, return not implemented - association creation might require database management endpoints
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "User association creation not yet implemented for v0.0.44")
}
