package v0_0_43

import (
	"context"
	"fmt"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// AccountManagerImpl implements the AccountManager interface for v0.0.43
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

	// TODO: Implement actual API call when account endpoints are available in OpenAPI spec
	// For now, return NotImplementedError as the actual implementation requires
	// the generated client to have account-related methods
	return nil, errors.NewNotImplementedError("account listing", "v0.0.43")
}

// Get retrieves a specific account by name
func (a *AccountManagerImpl) Get(ctx context.Context, accountName string) (*interfaces.Account, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	// TODO: Implement actual API call when account endpoints are available in OpenAPI spec
	return nil, errors.NewNotImplementedError("account retrieval", "v0.0.43")
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

	// Validate account hierarchy
	if account.ParentAccount != "" {
		// In a real implementation, we would verify the parent account exists
		// For now, just check it's not the same as the account being created
		if account.ParentAccount == account.Name {
			return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account cannot be its own parent", "parentAccount", account.ParentAccount, nil)
		}
	}

	// TODO: Implement actual API call when account endpoints are available in OpenAPI spec
	return nil, errors.NewNotImplementedError("account creation", "v0.0.43")
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

	// TODO: Implement actual API call when account endpoints are available in OpenAPI spec
	return errors.NewNotImplementedError("account update", "v0.0.43")
}

// Delete deletes an account
func (a *AccountManagerImpl) Delete(ctx context.Context, accountName string) error {
	if a.client == nil || a.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	// Call the generated OpenAPI client
	resp, err := a.client.apiClient.SlurmdbV0043DeleteAccountWithResponse(ctx, accountName)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status (200 or 204 for successful deletion)
	if resp.StatusCode() != 200 && resp.StatusCode() != 204 {
		var responseBody []byte
		if resp.JSON200 != nil {
			// Try to extract error details from response
			if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
				apiErrors := make([]errors.SlurmAPIErrorDetail, len(*resp.JSON200.Errors))
				for i, apiErr := range *resp.JSON200.Errors {
					var errorNumber int
					if apiErr.ErrorNumber != nil {
						errorNumber = int(*apiErr.ErrorNumber)
					}
					var errorCode string
					if apiErr.Error != nil {
						errorCode = *apiErr.Error
					}
					var source string
					if apiErr.Source != nil {
						source = *apiErr.Source
					}
					var description string
					if apiErr.Description != nil {
						description = *apiErr.Description
					}

					apiErrors[i] = errors.SlurmAPIErrorDetail{
						ErrorNumber: errorNumber,
						ErrorCode:   errorCode,
						Source:      source,
						Description: description,
					}
				}
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
		return httpErr
	}

	return nil
}

// GetAccountHierarchy retrieves the complete account hierarchy starting from a root account
func (a *AccountManagerImpl) GetAccountHierarchy(ctx context.Context, rootAccount string) (*interfaces.AccountHierarchy, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if rootAccount == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "root account name is required", "rootAccount", rootAccount, nil)
	}

	// TODO: Implement actual API call to retrieve account hierarchy
	// This would involve:
	// 1. Get the root account details
	// 2. Recursively get all child accounts
	// 3. Build the hierarchy structure with aggregated quotas and usage
	return nil, errors.NewNotImplementedError("account hierarchy retrieval", "v0.0.43")
}

// GetParentAccounts retrieves the parent chain for an account up to the root
func (a *AccountManagerImpl) GetParentAccounts(ctx context.Context, accountName string) ([]*interfaces.Account, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	// TODO: Implement actual API call to retrieve parent account chain
	// This would involve traversing up the account hierarchy from the given account
	// to the root account, collecting all parent accounts
	return nil, errors.NewNotImplementedError("parent accounts retrieval", "v0.0.43")
}

// GetChildAccounts retrieves child accounts with optional depth limiting
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

	// TODO: Implement actual API call to retrieve child accounts
	// This would involve recursively getting child accounts up to the specified depth
	// depth=0 means unlimited depth, depth=1 means direct children only
	return nil, errors.NewNotImplementedError("child accounts retrieval", "v0.0.43")
}

// GetAccountQuotas retrieves quota information for an account
func (a *AccountManagerImpl) GetAccountQuotas(ctx context.Context, accountName string) (*interfaces.AccountQuota, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	// TODO: Implement actual API call to retrieve account quotas
	// This would involve querying SLURM's accounting database for quota information
	// including CPU limits, job limits, TRES quotas, etc.
	return nil, errors.NewNotImplementedError("account quotas retrieval", "v0.0.43")
}

// GetAccountQuotaUsage retrieves quota usage information for an account within a timeframe
func (a *AccountManagerImpl) GetAccountQuotaUsage(ctx context.Context, accountName string, timeframe string) (*interfaces.AccountUsage, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	if timeframe == "" {
		timeframe = "current" // Default to current usage period
	}

	// Validate timeframe format (could be "current", "daily", "weekly", "monthly", "yearly")
	validTimeframes := []string{"current", "daily", "weekly", "monthly", "yearly"}
	isValid := false
	for _, valid := range validTimeframes {
		if timeframe == valid {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "timeframe must be one of: current, daily, weekly, monthly, yearly", "timeframe", timeframe, nil)
	}

	// TODO: Implement actual API call to retrieve account usage statistics
	// This would involve querying SLURM's accounting database for usage information
	// including CPU hours, job counts, efficiency ratios, etc.
	return nil, errors.NewNotImplementedError("account quota usage retrieval", "v0.0.43")
}

// GetAccountUsers retrieves all users associated with an account
func (a *AccountManagerImpl) GetAccountUsers(ctx context.Context, accountName string, opts *interfaces.ListAccountUsersOptions) ([]*interfaces.UserAccountAssociation, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	// TODO: Implement actual API call to retrieve account users
	// This would involve querying SLURM's association database for all users
	// associated with the given account, including their roles and permissions
	return nil, errors.NewNotImplementedError("account users retrieval", "v0.0.43")
}

// ValidateUserAccess validates user access to an account
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

	// TODO: Implement actual API call to validate user access
	// This would involve checking user associations, permissions, and quotas
	return nil, errors.NewNotImplementedError("user access validation", "v0.0.43")
}

// GetAccountUsersWithPermissions retrieves users with specific permissions for an account
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

	// Validate permission names
	validPermissions := []string{"read", "write", "admin", "coordinator", "submit", "cancel", "modify"}
	for _, perm := range permissions {
		isValid := false
		for _, valid := range validPermissions {
			if perm == valid {
				isValid = true
				break
			}
		}
		if !isValid {
			return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, fmt.Sprintf("invalid permission: %s", perm), "permissions", perm, nil)
		}
	}

	// TODO: Implement actual API call to retrieve users with specific permissions
	return nil, errors.NewNotImplementedError("account users with permissions retrieval", "v0.0.43")
}

// GetAccountFairShare retrieves fair-share configuration and state for an account
func (a *AccountManagerImpl) GetAccountFairShare(ctx context.Context, accountName string) (*interfaces.AccountFairShare, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	// TODO: Implement actual API call to retrieve account fair-share information
	// This would involve querying SLURM's shares database for account-level configuration
	// including shares, usage, and hierarchical fair-share data
	return nil, errors.NewNotImplementedError("account fair-share retrieval", "v0.0.43")
}

// GetFairShareHierarchy retrieves the complete fair-share tree structure
func (a *AccountManagerImpl) GetFairShareHierarchy(ctx context.Context, rootAccount string) (*interfaces.FairShareHierarchy, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if rootAccount == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "root account name is required", "rootAccount", rootAccount, nil)
	}

	// TODO: Implement actual API call to retrieve complete fair-share hierarchy
	// This would involve querying SLURM's shares database for the complete tree structure
	// starting from the specified root account, including all child accounts and users
	return nil, errors.NewNotImplementedError("fair-share hierarchy retrieval", "v0.0.43")
}

// Helper function to validate TRES format
func validateTRES(tres map[string]int) error {
	// TRES (Trackable Resources) typically include: cpu, mem, energy, node, billing, fs/disk, vmem, pages
	// Values should be non-negative
	for resource, value := range tres {
		if value < 0 {
			return errors.NewValidationError(errors.ErrorCodeValidationFailed, fmt.Sprintf("invalid TRES value for %s: must be non-negative", resource), "tres."+resource, value, nil)
		}
	}
	return nil
}

// convertAPIAccountToInterface converts V0043Account to interfaces.Account
func convertAPIAccountToInterface(apiAccount V0043Account) (*interfaces.Account, error) {
	account := &interfaces.Account{}

	// Basic fields
	account.Name = apiAccount.Name
	account.Description = apiAccount.Description
	account.Organization = apiAccount.Organization

	// Flags
	if apiAccount.Flags != nil {
		flags := make([]string, 0, len(*apiAccount.Flags))
		for _, flag := range *apiAccount.Flags {
			flags = append(flags, string(flag))
		}
		account.Flags = flags
	}

	// Coordinators
	if apiAccount.Coordinators != nil {
		coordinators := make([]string, 0, len(*apiAccount.Coordinators))
		for _, coord := range *apiAccount.Coordinators {
			coordinators = append(coordinators, coord.Name)
		}
		account.CoordinatorUsers = coordinators
	}

	return account, nil
}

// filterAccounts applies client-side filtering to the account list
func filterAccounts(accounts []interfaces.Account, opts *interfaces.ListAccountsOptions) []interfaces.Account {
	if opts == nil {
		return accounts
	}

	filtered := make([]interfaces.Account, 0, len(accounts))
	for _, account := range accounts {
		// Filter by names
		if len(opts.Names) > 0 {
			found := false
			for _, name := range opts.Names {
				if account.Name == name {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filter by organizations
		if len(opts.Organizations) > 0 {
			found := false
			for _, org := range opts.Organizations {
				if account.Organization == org {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filter by parent accounts
		if len(opts.ParentAccounts) > 0 {
			found := false
			for _, parent := range opts.ParentAccounts {
				if account.ParentAccount == parent {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filter by include deleted
		if !opts.WithDeleted {
			// Check if account has DELETED flag
			hasDeletedFlag := false
			for _, flag := range account.Flags {
				if flag == "DELETED" {
					hasDeletedFlag = true
					break
				}
			}
			if hasDeletedFlag {
				continue
			}
		}

		filtered = append(filtered, account)
	}

	return filtered
}

// convertAccountCreateToAPI converts interfaces.AccountCreate to API format
func convertAccountCreateToAPI(create *interfaces.AccountCreate) (*V0043Account, error) {
	apiAccount := &V0043Account{
		Name:         create.Name,
		Description:  create.Description,
		Organization: create.Organization,
	}

	// Flags
	if len(create.Flags) > 0 {
		flags := make([]V0043AccountFlags, 0, len(create.Flags))
		for _, flag := range create.Flags {
			flags = append(flags, V0043AccountFlags(flag))
		}
		apiAccount.Flags = &flags
	}

	// Coordinators
	if len(create.CoordinatorUsers) > 0 {
		coords := make(V0043CoordList, 0, len(create.CoordinatorUsers))
		for _, coordName := range create.CoordinatorUsers {
			coords = append(coords, V0043Coord{
				Name:   coordName,
				Direct: &[]bool{true}[0],
			})
		}
		apiAccount.Coordinators = &coords
	}

	return apiAccount, nil
}

// convertAccountUpdateToAPI converts interfaces.AccountUpdate to API format
func convertAccountUpdateToAPI(update *interfaces.AccountUpdate) (*V0043Account, error) {
	apiAccount := &V0043Account{}

	// Description
	if update.Description != nil {
		apiAccount.Description = *update.Description
	}

	// Organization
	if update.Organization != nil {
		apiAccount.Organization = *update.Organization
	}

	// Coordinators
	if len(update.CoordinatorUsers) > 0 {
		coords := make(V0043CoordList, 0, len(update.CoordinatorUsers))
		for _, coordName := range update.CoordinatorUsers {
			coords = append(coords, V0043Coord{
				Name:   coordName,
				Direct: &[]bool{true}[0],
			})
		}
		apiAccount.Coordinators = &coords
	}

	return apiAccount, nil
}