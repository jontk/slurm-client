// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"context"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// AccountManagerImpl implements the AccountManager interface for v0.0.40
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
	if err := a.client.CheckContext(ctx); err != nil {
		return nil, err
	}

	// Prepare parameters
	params := &SlurmdbV0040GetAccountsParams{}
	if opts != nil {
		// Note: v0.0.40 doesn't have a Name parameter, need to filter client-side
		if opts.WithAssociations {
			withAssocsStr := "true"
			params.WithAssocs = &withAssocsStr
		}
		if opts.WithDeleted {
			withDeletedStr := "true"
			params.WithDeleted = &withDeletedStr
		}
	}

	// Make API call
	resp, err := a.client.apiClient.SlurmdbV0040GetAccountsWithResponse(ctx, params)
	if err != nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal, "failed to list accounts")
	}

	// Check response
	if resp.StatusCode() != 200 {
		return nil, a.client.HandleErrorResponse(resp.StatusCode(), resp.Body)
	}

	if resp.JSON200 == nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal, "received nil response")
	}

	// Convert response
	accountList := &interfaces.AccountList{
		Accounts: make([]interfaces.Account, 0),
	}

	// Apply name filtering if requested (v0.0.40 doesn't support server-side name filtering)
	nameFilter := make(map[string]bool)
	if opts != nil && opts.Names != nil && len(opts.Names) > 0 {
		for _, name := range opts.Names {
			nameFilter[name] = true
		}
	}

	if len(resp.JSON200.Accounts) > 0 {
		for _, acc := range resp.JSON200.Accounts {
			account := a.convertV0040AccountToInterface(acc)
			
			// Apply name filter if specified
			if len(nameFilter) > 0 {
				if !nameFilter[account.Name] {
					continue
				}
			}
			
			accountList.Accounts = append(accountList.Accounts, *account)
		}
	}

	return accountList, nil
}

// Get retrieves a specific account by name
func (a *AccountManagerImpl) Get(ctx context.Context, accountName string) (*interfaces.Account, error) {
	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	if err := a.client.CheckContext(ctx); err != nil {
		return nil, err
	}

	// Make API call
	params := &SlurmdbV0040GetAccountParams{}
	resp, err := a.client.apiClient.SlurmdbV0040GetAccountWithResponse(ctx, accountName, params)
	if err != nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal, "failed to get account")
	}

	// Check response
	if resp.StatusCode() != 200 {
		return nil, a.client.HandleErrorResponse(resp.StatusCode(), resp.Body)
	}

	if resp.JSON200 == nil || resp.JSON200.Accounts == nil || len(resp.JSON200.Accounts) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "account not found")
	}

	// Convert the first account
	accounts := resp.JSON200.Accounts
	return a.convertV0040AccountToInterface(accounts[0]), nil
}

// Create creates a new account
func (a *AccountManagerImpl) Create(ctx context.Context, account *interfaces.AccountCreate) (*interfaces.AccountCreateResponse, error) {
	if account == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account data is required", "account", account, nil)
	}

	if err := a.client.CheckContext(ctx); err != nil {
		return nil, err
	}

	// Convert to API format
	apiAccount := a.convertInterfaceAccountCreateToV0040(account)
	
	// Create request body
	reqBody := V0040OpenapiAccountsResp{
		Accounts: []V0040Account{*apiAccount},
	}

	// Make API call
	resp, err := a.client.apiClient.SlurmdbV0040PostAccountsWithResponse(ctx, reqBody)
	if err != nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal, "failed to create account")
	}

	// Check response
	if resp.StatusCode() != 200 {
		return nil, a.client.HandleErrorResponse(resp.StatusCode(), resp.Body)
	}

	return &interfaces.AccountCreateResponse{
		AccountName: account.Name,
	}, nil
}

// Update updates an existing account
func (a *AccountManagerImpl) Update(ctx context.Context, accountName string, update *interfaces.AccountUpdate) error {
	if accountName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}
	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "update data is required", "update", update, nil)
	}

	if err := a.client.CheckContext(ctx); err != nil {
		return err
	}

	// Get existing account first
	existing, err := a.Get(ctx, accountName)
	if err != nil {
		return err
	}

	// Apply updates to existing account
	apiAccount := a.convertInterfaceAccountToV0040Update(existing, update)

	// Create request body
	reqBody := V0040OpenapiAccountsResp{
		Accounts: []V0040Account{*apiAccount},
	}

	// Make API call
	resp, err := a.client.apiClient.SlurmdbV0040PostAccountsWithResponse(ctx, reqBody)
	if err != nil {
		return errors.NewSlurmError(errors.ErrorCodeServerInternal, "failed to update account")
	}

	// Check response
	if resp.StatusCode() != 200 {
		return a.client.HandleErrorResponse(resp.StatusCode(), resp.Body)
	}

	return nil
}

// Delete deletes an account
func (a *AccountManagerImpl) Delete(ctx context.Context, accountName string) error {
	if accountName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	if err := a.client.CheckContext(ctx); err != nil {
		return err
	}

	// Make API call
	resp, err := a.client.apiClient.SlurmdbV0040DeleteAccountWithResponse(ctx, accountName)
	if err != nil {
		return errors.NewSlurmError(errors.ErrorCodeServerInternal, "failed to delete account")
	}

	// Check response
	if resp.StatusCode() != 200 && resp.StatusCode() != 204 {
		return a.client.HandleErrorResponse(resp.StatusCode(), resp.Body)
	}

	return nil
}

// GetAccountHierarchy retrieves the complete account hierarchy starting from a root account
func (a *AccountManagerImpl) GetAccountHierarchy(ctx context.Context, rootAccount string) (*interfaces.AccountHierarchy, error) {
	if rootAccount == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "root account name is required", "rootAccount", rootAccount, nil)
	}
	return nil, errors.NewNotImplementedError("account hierarchy retrieval", "v0.0.40")
}

// GetParentAccounts retrieves the parent chain for an account up to the root
func (a *AccountManagerImpl) GetParentAccounts(ctx context.Context, accountName string) ([]*interfaces.Account, error) {
	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}
	return nil, errors.NewNotImplementedError("parent accounts retrieval", "v0.0.40")
}

// GetChildAccounts retrieves child accounts with optional depth limiting
func (a *AccountManagerImpl) GetChildAccounts(ctx context.Context, accountName string, depth int) ([]*interfaces.Account, error) {
	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}
	if depth < 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "depth must be non-negative (0 means unlimited)", "depth", depth, nil)
	}
	return nil, errors.NewNotImplementedError("child accounts retrieval", "v0.0.40")
}

// GetAccountQuotas retrieves quota information for an account
func (a *AccountManagerImpl) GetAccountQuotas(ctx context.Context, accountName string) (*interfaces.AccountQuota, error) {
	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}
	return nil, errors.NewNotImplementedError("account quotas retrieval", "v0.0.40")
}

// GetAccountQuotaUsage retrieves quota usage information for an account within a timeframe
func (a *AccountManagerImpl) GetAccountQuotaUsage(ctx context.Context, accountName string, timeframe string) (*interfaces.AccountUsage, error) {
	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}
	return nil, errors.NewNotImplementedError("account quota usage retrieval", "v0.0.40")
}

// GetAccountUsers retrieves all users associated with an account
func (a *AccountManagerImpl) GetAccountUsers(ctx context.Context, accountName string, opts *interfaces.ListAccountUsersOptions) ([]*interfaces.UserAccountAssociation, error) {
	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}
	return nil, errors.NewNotImplementedError("account users retrieval", "v0.0.40")
}

// ValidateUserAccess validates user access to an account
func (a *AccountManagerImpl) ValidateUserAccess(ctx context.Context, userName, accountName string) (*interfaces.UserAccessValidation, error) {
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}
	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}
	return nil, errors.NewNotImplementedError("user access validation", "v0.0.40")
}

// GetAccountUsersWithPermissions retrieves users with specific permissions for an account
func (a *AccountManagerImpl) GetAccountUsersWithPermissions(ctx context.Context, accountName string, permissions []string) ([]*interfaces.UserAccountAssociation, error) {
	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}
	if len(permissions) == 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "at least one permission is required", "permissions", permissions, nil)
	}
	return nil, errors.NewNotImplementedError("account users with permissions retrieval", "v0.0.40")
}

// GetAccountFairShare retrieves fair-share configuration and state for an account
func (a *AccountManagerImpl) GetAccountFairShare(ctx context.Context, accountName string) (*interfaces.AccountFairShare, error) {
	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}
	return nil, errors.NewNotImplementedError("account fair-share retrieval", "v0.0.40")
}

// GetFairShareHierarchy retrieves the complete fair-share tree structure
func (a *AccountManagerImpl) GetFairShareHierarchy(ctx context.Context, rootAccount string) (*interfaces.FairShareHierarchy, error) {
	if rootAccount == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "root account name is required", "rootAccount", rootAccount, nil)
	}
	return nil, errors.NewNotImplementedError("fair-share hierarchy retrieval", "v0.0.40")
}

// convertV0040AccountToInterface converts v0.0.40 account to interface format
func (a *AccountManagerImpl) convertV0040AccountToInterface(acc V0040Account) *interfaces.Account {
	account := &interfaces.Account{}
	
	account.Name = acc.Name
	account.Description = acc.Description
	// Organization field doesn't exist in V0040Account, skip it
	
	// Convert coordinators
	// Convert coordinators
	if acc.Coordinators != nil && len(*acc.Coordinators) > 0 {
		account.CoordinatorUsers = make([]string, len(*acc.Coordinators))
		for i, coord := range *acc.Coordinators {
			account.CoordinatorUsers[i] = coord.Name
		}
	}
	
	return account
}

// convertInterfaceAccountCreateToV0040 converts interface account create to v0.0.40 format
func (a *AccountManagerImpl) convertInterfaceAccountCreateToV0040(acc *interfaces.AccountCreate) *V0040Account {
	apiAccount := &V0040Account{
		Name:         acc.Name,
		Description:  acc.Description,
		// Organization field doesn't exist in V0040Account
	}
	
	if len(acc.CoordinatorUsers) > 0 {
		coords := make(V0040CoordList, len(acc.CoordinatorUsers))
		for i, user := range acc.CoordinatorUsers {
			coords[i] = V0040Coord{
				Name: user,
			}
		}
		apiAccount.Coordinators = &coords
	}
	
	return apiAccount
}

// convertInterfaceAccountToV0040Update converts interface account with updates to v0.0.40 format
func (a *AccountManagerImpl) convertInterfaceAccountToV0040Update(existing *interfaces.Account, update *interfaces.AccountUpdate) *V0040Account {
	apiAccount := &V0040Account{
		Name: existing.Name,
	}
	
	// Apply updates
	if update.Description != nil {
		apiAccount.Description = *update.Description
	} else {
		apiAccount.Description = existing.Description
	}
	
	// Organization field doesn't exist in V0040Account
	
	if len(update.CoordinatorUsers) > 0 {
		coords := make(V0040CoordList, len(update.CoordinatorUsers))
		for i, user := range update.CoordinatorUsers {
			coords[i] = V0040Coord{
				Name: user,
			}
		}
		apiAccount.Coordinators = &coords
	} else if len(existing.CoordinatorUsers) > 0 {
		coords := make(V0040CoordList, len(existing.CoordinatorUsers))
		for i, user := range existing.CoordinatorUsers {
			coords[i] = V0040Coord{
				Name: user,
			}
		}
		apiAccount.Coordinators = &coords
	}
	
	return apiAccount
}
