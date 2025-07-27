package v0_0_43

import (
	"context"
	"fmt"
	"time"

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

	// Prepare parameters for the API call
	params := &SlurmdbV0043GetAccountsParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.apiClient.SlurmdbV0043GetAccountsWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
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
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response but got nil")
	}

	// Convert the response to our interface types
	accounts := make([]interfaces.Account, 0, len(resp.JSON200.Accounts))
	for _, apiAccount := range resp.JSON200.Accounts {
		account, err := convertAPIAccountToInterface(apiAccount)
		if err != nil {
			conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert account data")
			conversionErr.Cause = err
			conversionErr.Details = fmt.Sprintf("Error converting account %s", apiAccount.Name)
			return nil, conversionErr
		}
		accounts = append(accounts, *account)
	}

	// Apply client-side filtering if options are provided
	if opts != nil {
		accounts = filterAccounts(accounts, opts)
	}

	return &interfaces.AccountList{
		Accounts: accounts,
		Total:    len(accounts),
	}, nil
}

// Get retrieves a specific account by name
func (a *AccountManagerImpl) Get(ctx context.Context, accountName string) (*interfaces.Account, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0043GetAccountParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.apiClient.SlurmdbV0043GetAccountWithResponse(ctx, accountName, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
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
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil || len(resp.JSON200.Accounts) == 0 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("account %s not found", accountName))
	}

	// Convert the first account (should be the only one)
	account, err := convertAPIAccountToInterface(resp.JSON200.Accounts[0])
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert account data")
		conversionErr.Cause = err
		conversionErr.Details = fmt.Sprintf("Error converting account %s", accountName)
		return nil, conversionErr
	}

	return account, nil
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

	// Convert the account create request to API format
	apiAccount, err := convertAccountCreateToAPI(account)
	if err != nil {
		return nil, err
	}

	// Create request body
	reqBody := SlurmdbV0043PostAccountsJSONRequestBody{
		Accounts: V0043AccountList{*apiAccount},
	}

	// Call the generated OpenAPI client
	resp, err := a.client.apiClient.SlurmdbV0043PostAccountsWithResponse(ctx, reqBody)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
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
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
		return nil, httpErr
	}

	return &interfaces.AccountCreateResponse{
		AccountName: account.Name,
	}, nil
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

	// First, get the existing account to merge updates
	existingAccount, err := a.Get(ctx, accountName)
	if err != nil {
		return err
	}

	// Convert existing account to API format and apply updates
	apiAccount := &V0043Account{
		Name:         accountName,
		Description:  existingAccount.Description,
		Organization: existingAccount.Organization,
	}

	// Apply updates
	if update.Description != nil {
		apiAccount.Description = *update.Description
	}
	if update.Organization != nil {
		apiAccount.Organization = *update.Organization
	}

	// Handle coordinators
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

	// Create request body
	reqBody := SlurmdbV0043PostAccountsJSONRequestBody{
		Accounts: V0043AccountList{*apiAccount},
	}

	// Call the generated OpenAPI client (POST is used for updates in SLURM API)
	resp, err := a.client.apiClient.SlurmdbV0043PostAccountsWithResponse(ctx, reqBody)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
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

	// Get the root account
	rootAccountData, err := a.Get(ctx, rootAccount)
	if err != nil {
		return nil, err
	}

	// Build the hierarchy recursively
	hierarchy, err := a.buildAccountHierarchy(ctx, rootAccountData, 0, []string{rootAccount})
	if err != nil {
		return nil, err
	}

	return hierarchy, nil
}

// GetParentAccounts retrieves the parent chain for an account up to the root
func (a *AccountManagerImpl) GetParentAccounts(ctx context.Context, accountName string) ([]*interfaces.Account, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	// Get all accounts to build the parent chain
	allAccounts, err := a.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Build a map for quick lookup
	accountMap := make(map[string]*interfaces.Account)
	for i := range allAccounts.Accounts {
		accountMap[allAccounts.Accounts[i].Name] = &allAccounts.Accounts[i]
	}

	// Find the account
	account, exists := accountMap[accountName]
	if !exists {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("account %s not found", accountName))
	}

	// Build parent chain
	var parents []*interfaces.Account
	current := account
	visited := make(map[string]bool) // Prevent cycles

	for current.ParentAccount != "" {
		if visited[current.ParentAccount] {
			// Cycle detected
			break
		}
		visited[current.ParentAccount] = true

		parent, exists := accountMap[current.ParentAccount]
		if !exists {
			// Parent not found, stop here
			break
		}

		parents = append(parents, parent)
		current = parent
	}

	return parents, nil
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

	// Get all accounts
	allAccounts, err := a.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Build parent-child relationships
	childrenMap := make(map[string][]*interfaces.Account)
	accountMap := make(map[string]*interfaces.Account)

	for i := range allAccounts.Accounts {
		account := &allAccounts.Accounts[i]
		accountMap[account.Name] = account

		if account.ParentAccount != "" {
			childrenMap[account.ParentAccount] = append(childrenMap[account.ParentAccount], account)
		}
	}

	// Check if the account exists
	if _, exists := accountMap[accountName]; !exists {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("account %s not found", accountName))
	}

	// Collect children recursively
	var children []*interfaces.Account
	a.collectChildAccounts(accountName, childrenMap, &children, depth, 1)

	return children, nil
}

// GetAccountQuotas retrieves quota information for an account
func (a *AccountManagerImpl) GetAccountQuotas(ctx context.Context, accountName string) (*interfaces.AccountQuota, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	// Get account with associations
	params := &SlurmdbV0043GetAccountParams{}
	resp, err := a.client.apiClient.SlurmdbV0043GetAccountWithResponse(ctx, accountName, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("account %s not found", accountName))
	}

	if resp.JSON200 == nil || len(resp.JSON200.Accounts) == 0 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("account %s not found", accountName))
	}

	// Extract quota information from associations
	// Note: In SLURM, quotas are often stored in associations
	accountQuota := &interfaces.AccountQuota{
		// Initialize with default values
		GrpTRES:      make(map[string]int),
		GrpTRESUsed:  make(map[string]int),
		MaxTRES:      make(map[string]int),
		MaxTRESUsed:  make(map[string]int),
	}

	// Extract from associations if available
	apiAccount := resp.JSON200.Accounts[0]
	if apiAccount.Associations != nil {
		// Process associations to extract quota information
		// This is a simplified implementation as the actual quota structure
		// depends on SLURM configuration
		for _, assoc := range *apiAccount.Associations {
			// Extract quota limits from association
			// Note: This would need to be enhanced based on actual SLURM API response
			_ = assoc // Placeholder to avoid unused variable
		}
	}

	return accountQuota, nil
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

	// Get associations to calculate usage
	assocParams := &SlurmdbV0043GetAssociationsParams{
		Account: &accountName,
	}

	resp, err := a.client.apiClient.SlurmdbV0043GetAssociationsWithResponse(ctx, assocParams)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("account %s not found", accountName))
	}

	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format")
	}

	// Calculate usage statistics
	usage := &interfaces.AccountUsage{
		AccountName: accountName,
		Period:      timeframe,
		StartTime:   time.Now().Add(-24 * time.Hour), // Default to last 24 hours
		EndTime:     time.Now(),
		TRESUsage:   make(map[string]float64),
	}

	// Process associations to calculate usage
	if len(resp.JSON200.Associations) > 0 {
		userMap := make(map[string]bool)
		for _, assoc := range resp.JSON200.Associations {
			// Count unique users
			if assoc.User != "" {
				userMap[assoc.User] = true
			}
		}
		usage.UserCount = len(userMap)

		// Extract active users
		for user := range userMap {
			usage.ActiveUsers = append(usage.ActiveUsers, user)
		}
	}

	return usage, nil
}

// GetAccountUsers retrieves all users associated with an account
func (a *AccountManagerImpl) GetAccountUsers(ctx context.Context, accountName string, opts *interfaces.ListAccountUsersOptions) ([]*interfaces.UserAccountAssociation, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	// Get associations for the account
	assocParams := &SlurmdbV0043GetAssociationsParams{
		Account: &accountName,
	}

	resp, err := a.client.apiClient.SlurmdbV0043GetAssociationsWithResponse(ctx, assocParams)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("account %s not found", accountName))
	}

	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format")
	}

	// Convert associations to user-account associations
	var userAssociations []*interfaces.UserAccountAssociation

	for _, assoc := range resp.JSON200.Associations {
		if assoc.User == "" {
			continue
		}

		userAssoc := &interfaces.UserAccountAssociation{
			UserName:    assoc.User,
			AccountName: accountName,
		}

		// Set cluster if available
		if assoc.Cluster != nil {
			userAssoc.Cluster = *assoc.Cluster
		}

		// Set partition if available
		if assoc.Partition != nil {
			userAssoc.Partition = *assoc.Partition
		}

		// Check if user is a coordinator
		if assoc.IsDefault != nil {
			userAssoc.IsDefault = *assoc.IsDefault
		}

		// Set default permissions
		// Note: Admin level would need to be fetched from user info separately
		userAssoc.Role = "user"
		userAssoc.Permissions = []string{"read", "submit"}
		userAssoc.IsActive = true

		userAssociations = append(userAssociations, userAssoc)
	}

	// Apply filtering if options are provided
	if opts != nil {
		userAssociations = filterUserAssociations(userAssociations, opts)
	}

	return userAssociations, nil
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

	// Get associations for the user and account
	assocParams := &SlurmdbV0043GetAssociationsParams{
		Account: &accountName,
		User:    &userName,
	}

	resp, err := a.client.apiClient.SlurmdbV0043GetAssociationsWithResponse(ctx, assocParams)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to validate user access")
	}

	validation := &interfaces.UserAccessValidation{
		UserName:       userName,
		AccountName:    accountName,
		HasAccess:      false,
		AccessLevel:    "none",
		Permissions:    []string{},
		ValidationTime: time.Now(),
	}

	if resp.JSON200 != nil && len(resp.JSON200.Associations) > 0 {
		// User has access to the account
		validation.HasAccess = true
		validation.ValidFrom = time.Now() // Could be extracted from association if available

		// Extract access level and permissions from the first matching association
		assoc := resp.JSON200.Associations[0]

		// Set default access level and permissions
		// Note: Admin level would need to be fetched from user info separately
		validation.AccessLevel = "user"
		validation.Permissions = []string{"read", "submit"}

		// Create association details
		validation.Association = &interfaces.UserAccountAssociation{
			UserName:    userName,
			AccountName: accountName,
			Role:        validation.AccessLevel,
			Permissions: validation.Permissions,
			IsActive:    true,
		}

		if assoc.Cluster != nil {
			validation.Association.Cluster = *assoc.Cluster
		}
		if assoc.Partition != nil {
			validation.Association.Partition = *assoc.Partition
		}
		if assoc.IsDefault != nil {
			validation.Association.IsDefault = *assoc.IsDefault
		}
	} else {
		validation.Reason = "No association found between user and account"
	}

	return validation, nil
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

	// First get all users for the account
	allUsers, err := a.GetAccountUsers(ctx, accountName, nil)
	if err != nil {
		return nil, err
	}

	// Filter users by permissions
	var filteredUsers []*interfaces.UserAccountAssociation
	for _, user := range allUsers {
		// Check if user has all required permissions
		hasAllPermissions := true
		for _, requiredPerm := range permissions {
			hasPermission := false
			for _, userPerm := range user.Permissions {
				if userPerm == requiredPerm {
					hasPermission = true
					break
				}
			}
			if !hasPermission {
				hasAllPermissions = false
				break
			}
		}

		if hasAllPermissions {
			filteredUsers = append(filteredUsers, user)
		}
	}

	return filteredUsers, nil
}

// GetAccountFairShare retrieves fair-share configuration and state for an account
func (a *AccountManagerImpl) GetAccountFairShare(ctx context.Context, accountName string) (*interfaces.AccountFairShare, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	// Note: Fair-share information is typically available through the shares API
	// For now, we'll return a basic structure as the exact API endpoint may vary
	fairShare := &interfaces.AccountFairShare{
		AccountName: accountName,
		Cluster:     "default", // Would be extracted from API
		Shares:      1,          // Default shares
		RawShares:   1,
		NormalizedShares: 1.0,
		Usage:       0.0,
		EffectiveUsage: 0.0,
		FairShareFactor: 1.0,
		Level:       1,
	}

	// Get account to check parent
	account, err := a.Get(ctx, accountName)
	if err != nil {
		return nil, err
	}

	if account.ParentAccount != "" {
		fairShare.Parent = account.ParentAccount
	}

	// Get users to count
	users, err := a.GetAccountUsers(ctx, accountName, nil)
	if err == nil {
		fairShare.UserCount = len(users)
		fairShare.ActiveUsers = len(users) // Would need to filter by activity
	}

	return fairShare, nil
}

// GetFairShareHierarchy retrieves the complete fair-share tree structure
func (a *AccountManagerImpl) GetFairShareHierarchy(ctx context.Context, rootAccount string) (*interfaces.FairShareHierarchy, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if rootAccount == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "root account name is required", "rootAccount", rootAccount, nil)
	}

	// Build the account hierarchy first
	accountHierarchy, err := a.GetAccountHierarchy(ctx, rootAccount)
	if err != nil {
		return nil, err
	}

	// Convert to fair-share hierarchy
	fairShareHierarchy := &interfaces.FairShareHierarchy{
		Cluster:       "default", // Would be extracted from API
		RootAccount:   rootAccount,
		LastUpdate:    time.Now(),
		DecayHalfLife: 7 * 24,     // Default 7 days in hours
		UsageWindow:   30 * 24,    // Default 30 days in hours
		Algorithm:     "classic",  // Default algorithm
	}

	// Build fair-share tree from account hierarchy
	fairShareTree := a.buildFairShareNode(accountHierarchy)
	fairShareHierarchy.Tree = fairShareTree

	// Calculate total shares
	fairShareHierarchy.TotalShares = a.calculateTotalShares(fairShareTree)

	return fairShareHierarchy, nil
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

// buildAccountHierarchy recursively builds the account hierarchy
func (a *AccountManagerImpl) buildAccountHierarchy(ctx context.Context, account *interfaces.Account, level int, path []string) (*interfaces.AccountHierarchy, error) {
	hierarchy := &interfaces.AccountHierarchy{
		Account: account,
		Level:   level,
		Path:    path,
	}

	// Get child accounts
	children, err := a.GetChildAccounts(ctx, account.Name, 0)
	if err != nil {
		// Log error but don't fail - account might not have children
		children = []*interfaces.Account{}
	}

	hierarchy.TotalSubAccounts = len(children)

	// Build child hierarchies
	if len(children) > 0 {
		hierarchy.ChildAccounts = make([]*interfaces.AccountHierarchy, 0, len(children))
		for _, child := range children {
			childPath := append([]string{}, path...)
			childPath = append(childPath, child.Name)
			childHierarchy, err := a.buildAccountHierarchy(ctx, child, level+1, childPath)
			if err != nil {
				continue // Skip on error
			}
			hierarchy.ChildAccounts = append(hierarchy.ChildAccounts, childHierarchy)
			hierarchy.TotalSubAccounts += childHierarchy.TotalSubAccounts
		}
	}

	// Get account users
	users, err := a.GetAccountUsers(ctx, account.Name, nil)
	if err == nil {
		hierarchy.TotalUsers = len(users)
	}

	// Get quotas
	quota, err := a.GetAccountQuotas(ctx, account.Name)
	if err == nil {
		hierarchy.AggregateQuota = quota
	}

	// Get usage
	usage, err := a.GetAccountQuotaUsage(ctx, account.Name, "current")
	if err == nil {
		hierarchy.AggregateUsage = usage
	}

	return hierarchy, nil
}

// collectChildAccounts recursively collects child accounts up to the specified depth
func (a *AccountManagerImpl) collectChildAccounts(accountName string, childrenMap map[string][]*interfaces.Account, result *[]*interfaces.Account, maxDepth, currentDepth int) {
	if maxDepth > 0 && currentDepth > maxDepth {
		return
	}

	children, exists := childrenMap[accountName]
	if !exists {
		return
	}

	for _, child := range children {
		*result = append(*result, child)
		a.collectChildAccounts(child.Name, childrenMap, result, maxDepth, currentDepth+1)
	}
}

// filterUserAssociations applies filtering to user associations
func filterUserAssociations(associations []*interfaces.UserAccountAssociation, opts *interfaces.ListAccountUsersOptions) []*interfaces.UserAccountAssociation {
	if opts == nil {
		return associations
	}

	filtered := make([]*interfaces.UserAccountAssociation, 0, len(associations))
	for _, assoc := range associations {
		// Filter by roles
		if len(opts.Roles) > 0 {
			found := false
			for _, role := range opts.Roles {
				if assoc.Role == role {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filter by active only
		if opts.ActiveOnly && !assoc.IsActive {
			continue
		}

		// Filter by coordinators only
		if opts.CoordinatorsOnly && !assoc.IsCoordinator {
			continue
		}

		filtered = append(filtered, assoc)
	}

	return filtered
}

// buildFairShareNode converts account hierarchy to fair-share node
func (a *AccountManagerImpl) buildFairShareNode(hierarchy *interfaces.AccountHierarchy) *interfaces.FairShareNode {
	if hierarchy == nil || hierarchy.Account == nil {
		return nil
	}

	node := &interfaces.FairShareNode{
		Name:             hierarchy.Account.Name,
		Account:          hierarchy.Account.Name,
		Level:            hierarchy.Level,
		Shares:           1, // Default shares
		NormalizedShares: 1.0,
		Usage:            0.0,
		FairShareFactor:  1.0,
	}

	// Set parent if available
	if hierarchy.Account.ParentAccount != "" {
		node.Parent = hierarchy.Account.ParentAccount
	}

	// Add children
	if len(hierarchy.ChildAccounts) > 0 {
		node.Children = make([]*interfaces.FairShareNode, 0, len(hierarchy.ChildAccounts))
		for _, childHierarchy := range hierarchy.ChildAccounts {
			childNode := a.buildFairShareNode(childHierarchy)
			if childNode != nil {
				node.Children = append(node.Children, childNode)
			}
		}
	}

	return node
}

// calculateTotalShares recursively calculates total shares in the fair-share tree
func (a *AccountManagerImpl) calculateTotalShares(node *interfaces.FairShareNode) int {
	if node == nil {
		return 0
	}

	total := node.Shares
	for _, child := range node.Children {
		total += a.calculateTotalShares(child)
	}

	return total
}