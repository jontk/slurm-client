package v0_0_41

import (
	"context"
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
)

// AccountAdapter implements the AccountAdapter interface for v0.0.41
type AccountAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewAccountAdapter creates a new Account adapter for v0.0.41
func NewAccountAdapter(client *api.ClientWithResponses) *AccountAdapter {
	return &AccountAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Account"),
		client:      client,
		wrapper:     nil, // We'll implement this later
	}
}

// List retrieves a list of accounts with optional filtering
func (a *AccountAdapter) List(ctx context.Context, opts *types.AccountListOptions) (*types.AccountList, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmdbV0041GetAccountsParams{}

	// Apply filters from options
	if opts != nil {
		if len(opts.Names) > 0 {
			accountStr := strings.Join(opts.Names, ",")
			params.Account = &accountStr
		}
		if opts.Description != "" {
			params.Description = &opts.Description
		}
		if opts.Organization != "" {
			params.Organization = &opts.Organization
		}
		if opts.WithDeleted {
			withDeleted := "true"
			params.WithDeleted = &withDeleted
		}
		if opts.WithAssociations {
			withAssocs := "true"
			params.WithAssocs = &withAssocs
		}
		if opts.WithCoordinators {
			withCoords := "true"
			params.WithCoords = &withCoords
		}
	}

	// Make the API call
	resp, err := a.client.SlurmdbV0041GetAccountsWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list accounts")
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected nil response")
	}

	// Convert response to common types
	accountList := &types.AccountList{
		Accounts: make([]types.Account, 0, len(resp.JSON200.Accounts)),
		Meta: &types.ListMeta{
			Version: a.GetVersion(),
		},
	}

	for _, apiAccount := range resp.JSON200.Accounts {
		account, err := a.convertAPIAccountToCommon(apiAccount)
		if err != nil {
			// Log the error but continue processing other accounts
			continue
		}
		accountList.Accounts = append(accountList.Accounts, *account)
	}

	// Extract warning messages if any
	if resp.JSON200.Warnings != nil {
		warnings := make([]string, 0, len(*resp.JSON200.Warnings))
		for _, warning := range *resp.JSON200.Warnings {
			if warning.Description != nil {
				warnings = append(warnings, *warning.Description)
			}
		}
		if len(warnings) > 0 {
			accountList.Meta.Warnings = warnings
		}
	}

	// Extract error messages if any
	if resp.JSON200.Errors != nil {
		errors := make([]string, 0, len(*resp.JSON200.Errors))
		for _, error := range *resp.JSON200.Errors {
			if error.Description != nil {
				errors = append(errors, *error.Description)
			}
		}
		if len(errors) > 0 {
			accountList.Meta.Errors = errors
		}
	}

	return accountList, nil
}

// Get retrieves a specific account by name
func (a *AccountAdapter) Get(ctx context.Context, name string) (*types.Account, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Validate name
	if err := a.ValidateResourceName("account name", name); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Make the API call
	params := &api.SlurmdbV0041GetAccountParams{}
	resp, err := a.client.SlurmdbV0041GetAccountWithResponse(ctx, name, params)
	if err != nil {
		return nil, a.WrapError(err, fmt.Sprintf("failed to get account %s", name))
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}

	if resp.JSON200 == nil || len(resp.JSON200.Accounts) == 0 {
		return nil, a.HandleNotFound(fmt.Sprintf("account %s", name))
	}

	// Convert the first account in the response
	account, err := a.convertAPIAccountToCommon(resp.JSON200.Accounts[0])
	if err != nil {
		return nil, a.WrapError(err, fmt.Sprintf("failed to convert account %s", name))
	}

	return account, nil
}

// Create creates a new account
func (a *AccountAdapter) Create(ctx context.Context, account *types.Account) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Validate account
	if account == nil {
		return a.HandleValidationError("account cannot be nil")
	}
	if err := a.ValidateResourceName("account name", account.Name); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Convert account to API request
	createReq := a.convertCommonToAPIAccount(account)

	// Make the API call
	resp, err := a.client.SlurmdbV0041PostAccountsWithResponse(ctx, *createReq)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to create account %s", account.Name))
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
	}

	return nil
}

// Update updates an existing account
func (a *AccountAdapter) Update(ctx context.Context, name string, update *types.AccountUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Validate name
	if err := a.ValidateResourceName("account name", name); err != nil {
		return err
	}

	// Validate update
	if update == nil {
		return a.HandleValidationError("account update cannot be nil")
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Get the existing account first
	existingAccount, err := a.Get(ctx, name)
	if err != nil {
		return err
	}

	// Apply updates
	if update.Description != nil {
		existingAccount.Description = *update.Description
	}
	if update.Organization != nil {
		existingAccount.Organization = *update.Organization
	}

	// Convert to API request
	updateReq := a.convertCommonToAPIAccount(existingAccount)

	// Make the API call
	resp, err := a.client.SlurmdbV0041PostAccountsWithResponse(ctx, *updateReq)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to update account %s", name))
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
	}

	return nil
}

// Delete deletes an account
func (a *AccountAdapter) Delete(ctx context.Context, name string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Validate name
	if err := a.ValidateResourceName("account name", name); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Make the API call
	resp, err := a.client.SlurmdbV0041DeleteAccountWithResponse(ctx, name)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to delete account %s", name))
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
	}

	return nil
}

// GetAssociations gets associations for an account
func (a *AccountAdapter) GetAssociations(ctx context.Context, name string) (*types.AssociationList, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Validate name
	if err := a.ValidateResourceName("account name", name); err != nil {
		return nil, err
	}

	// v0.0.41 doesn't have a direct method to get associations for a specific account
	// We would need to use the association manager instead
	return nil, fmt.Errorf("getting associations for a specific account is not directly supported in API v0.0.41, use the association manager instead")
}

// AddUser adds a user to an account
func (a *AccountAdapter) AddUser(ctx context.Context, accountName string, userName string, opts *types.AccountUserOptions) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Validate names
	if err := a.ValidateResourceName("account name", accountName); err != nil {
		return err
	}
	if err := a.ValidateResourceName("user name", userName); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Create association request
	assocReq := api.SlurmdbV0041PostAccountsAssociationJSONBody{
		Accounts: &[]string{accountName},
		Users:    &[]string{userName},
	}

	// Apply options
	if opts != nil {
		if opts.Cluster != "" {
			assocReq.Cluster = &opts.Cluster
		}
		if opts.Partition != "" {
			assocReq.Partition = &opts.Partition
		}
		if opts.DefaultQoS != "" {
			assocReq.DefaultQos = &opts.DefaultQoS
		}
	}

	// Make the API call
	resp, err := a.client.SlurmdbV0041PostAccountsAssociationWithResponse(ctx, assocReq)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to add user %s to account %s", userName, accountName))
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
	}

	return nil
}

// RemoveUser removes a user from an account
func (a *AccountAdapter) RemoveUser(ctx context.Context, accountName string, userName string) error {
	// v0.0.41 doesn't have a direct method to remove a user from an account
	// This would need to be done through the association manager by deleting the association
	return fmt.Errorf("removing a user from an account is not directly supported in API v0.0.41, use the association manager instead")
}

// SetCoordinators sets the coordinators for an account
func (a *AccountAdapter) SetCoordinators(ctx context.Context, name string, coordinators []string) error {
	// v0.0.41 doesn't support setting coordinators through the API
	return fmt.Errorf("setting coordinators is not supported in API v0.0.41")
}