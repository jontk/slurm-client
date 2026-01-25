// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"context"
	"fmt"
	"strings"

	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/jontk/slurm-client/pkg/errors"
)

// AccountAdapter implements the AccountAdapter interface for v0.0.40
type AccountAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewAccountAdapter creates a new Account adapter for v0.0.40
func NewAccountAdapter(client *api.ClientWithResponses) *AccountAdapter {
	return &AccountAdapter{
		BaseManager: base.NewBaseManager("v0.0.40", "Account"),
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
	params := &api.SlurmdbV0040GetAccountsParams{}

	// Apply filters from options
	// Note: v0.0.40 API has limited parameter support
	if opts != nil {
		// Account and Organization filtering is not supported by v0.0.40 GetAccounts params
		// Only Description and with_* flags are available
		if len(opts.Descriptions) > 0 {
			descStr := strings.Join(opts.Descriptions, ",")
			params.Description = &descStr
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
		// TODO: Names and Organizations filtering would need to be done client-side
		// since v0.0.40 API doesn't support these parameters
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0040GetAccountsWithResponse(ctx, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.40"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "List Accounts"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Accounts, "List Accounts - accounts field"); err != nil {
		return nil, err
	}

	// Convert the response to common types
	accountList := make([]types.Account, 0, len(resp.JSON200.Accounts))
	for _, apiAccount := range resp.JSON200.Accounts {
		account := a.convertAPIAccountToCommon(apiAccount)
		accountList = append(accountList, *account)
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
	if start >= len(accountList) {
		return &types.AccountList{
			Accounts: []types.Account{},
			Total:    len(accountList),
		}, nil
	}

	end := len(accountList)
	if listOpts.Limit > 0 {
		end = start + listOpts.Limit
		if end > len(accountList) {
			end = len(accountList)
		}
	}

	return &types.AccountList{
		Accounts: accountList[start:end],
		Total:    len(accountList),
	}, nil
}

// Get retrieves a specific account by name
func (a *AccountAdapter) Get(ctx context.Context, accountName string) (*types.Account, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.ValidateResourceName(accountName, "account name"); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmdbV0040GetAccountParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0040GetAccountWithResponse(ctx, accountName, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.40"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "Get Account"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Accounts, "Get Account - accounts field"); err != nil {
		return nil, err
	}

	// Check if we got any account entries
	if len(resp.JSON200.Accounts) == 0 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("Account %s not found", accountName))
	}

	// Convert the first account (should be the only one)
	account := a.convertAPIAccountToCommon(resp.JSON200.Accounts[0])

	return account, nil
}

// Create creates a new account
func (a *AccountAdapter) Create(ctx context.Context, account *types.AccountCreate) (*types.AccountCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.validateAccountCreate(account); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert to API format
	apiAccount := a.convertCommonAccountCreateToAPI(account)

	// Create request body
	reqBody := api.SlurmdbV0040PostAccountsJSONRequestBody{
		Accounts: []api.V0040Account{*apiAccount},
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0040PostAccountsWithResponse(ctx, reqBody)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.40"); err != nil {
		return nil, err
	}

	return &types.AccountCreateResponse{
		AccountName: account.Name,
	}, nil
}

// Update updates an existing account
func (a *AccountAdapter) Update(ctx context.Context, accountName string, update *types.AccountUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(accountName, "account name"); err != nil {
		return err
	}
	if err := a.validateAccountUpdate(update); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// First, get the existing account to merge updates
	existingAccount, err := a.Get(ctx, accountName)
	if err != nil {
		return err
	}

	// Convert to API format and apply updates
	apiAccount := a.convertCommonAccountUpdateToAPI(existingAccount, update)

	// Create request body
	reqBody := api.SlurmdbV0040PostAccountsJSONRequestBody{
		Accounts: []api.V0040Account{*apiAccount},
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0040PostAccountsWithResponse(ctx, reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "v0.0.40")
}

// Delete deletes an account
func (a *AccountAdapter) Delete(ctx context.Context, accountName string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(accountName, "account name"); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0040DeleteAccountWithResponse(ctx, accountName)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	// Create adapter with special handling for 204 (No Content) status
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)

	// For DELETE operations, 204 is also a success
	if resp.StatusCode() == 204 {
		return nil
	}

	return common.HandleAPIResponse(responseAdapter, "v0.0.40")
}

// validateAccountCreate validates account creation request
func (a *AccountAdapter) validateAccountCreate(account *types.AccountCreate) error {
	if account == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "account creation data is required", "account", nil, nil)
	}
	if account.Name == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "name", account.Name, nil)
	}
	return nil
}

// validateAccountUpdate validates account update request
func (a *AccountAdapter) validateAccountUpdate(update *types.AccountUpdate) error {
	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "account update data is required", "update", nil, nil)
	}
	// At least one field should be provided for update
	if update.Description == nil && update.Organization == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "at least one field must be provided for update", "update", update, nil)
	}
	return nil
}

// CreateAssociation creates associations for accounts (not supported in v0.0.40)
func (a *AccountAdapter) CreateAssociation(ctx context.Context, req *types.AccountAssociationRequest) (*types.AssociationCreateResponse, error) {
	return nil, errors.NewNotImplementedError("CreateAssociation", a.GetVersion())
}
