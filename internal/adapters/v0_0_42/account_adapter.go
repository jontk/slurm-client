package v0_0_42

import (
	"context"
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

// AccountAdapter implements the AccountAdapter interface for v0.0.42
type AccountAdapter struct {
	*base.BaseManager
	client *api.ClientWithResponses
}

// NewAccountAdapter creates a new Account adapter for v0.0.42
func NewAccountAdapter(client *api.ClientWithResponses) *AccountAdapter {
	return &AccountAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Account"),
		client:      client,
	}
}

// List retrieves a list of accounts
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
	params := &api.SlurmdbV0042GetAccountsParams{}

	// Apply filters from options
	if opts != nil {
		if len(opts.Names) > 0 {
			accountStr := strings.Join(opts.Names, ",")
			params.Account = &accountStr
		}
		if opts.WithAssociations {
			withAssoc := "true"
			params.WithAssoc = &withAssoc
		}
		if opts.WithCoordinators {
			withCoord := "true"
			params.WithCoord = &withCoord
		}
		if opts.WithDeleted {
			withDeleted := "true"
			params.WithDeleted = &withDeleted
		}
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042GetAccountsWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list accounts")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(resp.StatusCode(), resp.Body)
	}

	// Check for API response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from API")
	}

	// Convert the response to common types
	accountList := &types.AccountList{
		Accounts: make([]*types.Account, 0),
	}

	if resp.JSON200.Accounts != nil {
		for _, apiAccount := range *resp.JSON200.Accounts {
			account, err := a.convertAPIAccountToCommon(apiAccount)
			if err != nil {
				// Log conversion error but continue
				continue
			}
			accountList.Accounts = append(accountList.Accounts, account)
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

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters
	params := &api.SlurmdbV0042GetAccountParams{}

	// Call the API
	resp, err := a.client.SlurmdbV0042GetAccountWithResponse(ctx, name, params)
	if err != nil {
		return nil, a.WrapError(err, fmt.Sprintf("failed to get account %s", name))
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(resp.StatusCode(), resp.Body)
	}

	// Check for API response
	if resp.JSON200 == nil || resp.JSON200.Accounts == nil || len(*resp.JSON200.Accounts) == 0 {
		return nil, fmt.Errorf("account %s not found", name)
	}

	// Convert the first account in the response
	accounts := *resp.JSON200.Accounts
	return a.convertAPIAccountToCommon(accounts[0])
}

// Create creates a new account
func (a *AccountAdapter) Create(ctx context.Context, account *types.AccountCreateRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Convert common account to API format
	apiAccount, err := a.convertCommonAccountCreateToAPI(account)
	if err != nil {
		return a.WrapError(err, "failed to convert account create request")
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042PostAccountsWithResponse(ctx, apiAccount)
	if err != nil {
		return a.WrapError(err, "failed to create account")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleAPIError(resp.StatusCode(), resp.Body)
	}

	return nil
}

// Update updates an existing account
func (a *AccountAdapter) Update(ctx context.Context, name string, updates *types.AccountUpdateRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// v0.0.42 doesn't have a direct account update endpoint
	// Updates are typically done through associations
	return fmt.Errorf("account update not directly supported via v0.0.42 API - use association updates")
}

// Delete deletes an account
func (a *AccountAdapter) Delete(ctx context.Context, name string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042DeleteAccountWithResponse(ctx, name)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to delete account %s", name))
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleAPIError(resp.StatusCode(), resp.Body)
	}

	return nil
}