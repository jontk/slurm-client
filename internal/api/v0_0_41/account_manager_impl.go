// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

import (
	"context"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// AccountManagerImpl implements the AccountManager interface for v0.0.41
type AccountManagerImpl struct {
	client *WrapperClient
}

// NewAccountManagerImpl creates a new AccountManagerImpl
func NewAccountManagerImpl(client *WrapperClient) *AccountManagerImpl {
	return &AccountManagerImpl{
		client: client,
	}
}

// List lists accounts with optional filtering
func (a *AccountManagerImpl) List(ctx context.Context, opts *interfaces.ListAccountsOptions) (*interfaces.AccountList, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// v0.0.41 has very limited account support
	// Most enhanced features are not available
	return nil, errors.NewNotImplementedError("account listing not supported", "v0.0.41")
}

// Get retrieves a specific account by name
func (a *AccountManagerImpl) Get(ctx context.Context, accountName string) (*interfaces.Account, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	// v0.0.41 has minimal account support
	return nil, errors.NewNotImplementedError("account retrieval not supported", "v0.0.41")
}

// Create creates a new account
func (a *AccountManagerImpl) Create(ctx context.Context, account *interfaces.AccountCreate) (*interfaces.AccountCreateResponse, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if account == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account data is required", "account", account, nil)
	}

	if account.Name == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "account.Name", account.Name, nil)
	}

	return nil, errors.NewNotImplementedError("account creation not supported", "v0.0.41")
}

// Update updates an existing account
func (a *AccountManagerImpl) Update(ctx context.Context, accountName string, update *interfaces.AccountUpdate) error {
	if a.client == nil || a.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "update data is required", "update", update, nil)
	}

	return errors.NewNotImplementedError("account update not supported", "v0.0.41")
}

// Delete deletes an account
func (a *AccountManagerImpl) Delete(ctx context.Context, accountName string) error {
	if a.client == nil || a.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	return errors.NewNotImplementedError("account deletion not supported", "v0.0.41")
}

// GetAccountHierarchy retrieves the complete account hierarchy starting from a root account
// This feature is not available in v0.0.41
func (a *AccountManagerImpl) GetAccountHierarchy(ctx context.Context, rootAccount string) (*interfaces.AccountHierarchy, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if rootAccount == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "root account name is required", "rootAccount", rootAccount, nil)
	}

	return nil, errors.NewNotImplementedError("account hierarchy not supported", "v0.0.41")
}

// GetParentAccounts retrieves the parent chain for an account up to the root
// This feature is not available in v0.0.41
func (a *AccountManagerImpl) GetParentAccounts(ctx context.Context, accountName string) ([]*interfaces.Account, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	return nil, errors.NewNotImplementedError("parent accounts not supported", "v0.0.41")
}

// GetChildAccounts retrieves child accounts with optional depth limiting
// This feature is not available in v0.0.41
func (a *AccountManagerImpl) GetChildAccounts(ctx context.Context, accountName string, depth int) ([]*interfaces.Account, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	if depth < 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "depth must be non-negative (0 means unlimited)", "depth", depth, nil)
	}

	return nil, errors.NewNotImplementedError("child accounts not supported", "v0.0.41")
}

// GetAccountQuotas retrieves quota information for an account
// This feature is not available in v0.0.41
func (a *AccountManagerImpl) GetAccountQuotas(ctx context.Context, accountName string) (*interfaces.AccountQuota, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	return nil, errors.NewNotImplementedError("account quotas not supported", "v0.0.41")
}

// GetAccountQuotaUsage retrieves quota usage information for an account within a timeframe
// This feature is not available in v0.0.41
func (a *AccountManagerImpl) GetAccountQuotaUsage(ctx context.Context, accountName string, timeframe string) (*interfaces.AccountUsage, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	return nil, errors.NewNotImplementedError("account quota usage not supported", "v0.0.41")
}

// GetAccountUsers retrieves all users associated with an account
// This feature is not available in v0.0.41
func (a *AccountManagerImpl) GetAccountUsers(ctx context.Context, accountName string, opts *interfaces.ListAccountUsersOptions) ([]*interfaces.UserAccountAssociation, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	return nil, errors.NewNotImplementedError("account users not supported", "v0.0.41")
}

// ValidateUserAccess validates user access to an account
// This feature is not available in v0.0.41
func (a *AccountManagerImpl) ValidateUserAccess(ctx context.Context, userName, accountName string) (*interfaces.UserAccessValidation, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	return nil, errors.NewNotImplementedError("user access validation not supported", "v0.0.41")
}

// GetAccountUsersWithPermissions retrieves users with specific permissions for an account
// This feature is not available in v0.0.41
func (a *AccountManagerImpl) GetAccountUsersWithPermissions(ctx context.Context, accountName string, permissions []string) ([]*interfaces.UserAccountAssociation, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	if len(permissions) == 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "at least one permission is required", "permissions", permissions, nil)
	}

	return nil, errors.NewNotImplementedError("account users with permissions not supported", "v0.0.41")
}

// GetAccountFairShare retrieves fair-share configuration and state for an account
// This feature is not available in v0.0.41
func (a *AccountManagerImpl) GetAccountFairShare(ctx context.Context, accountName string) (*interfaces.AccountFairShare, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	return nil, errors.NewNotImplementedError("account fair-share not supported", "v0.0.41")
}

// GetFairShareHierarchy retrieves the complete fair-share tree structure
// This feature is not available in v0.0.41
func (a *AccountManagerImpl) GetFairShareHierarchy(ctx context.Context, rootAccount string) (*interfaces.FairShareHierarchy, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if rootAccount == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "root account name is required", "rootAccount", rootAccount, nil)
	}

	return nil, errors.NewNotImplementedError("fair-share hierarchy not supported", "v0.0.41")
}
