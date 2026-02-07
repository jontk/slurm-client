// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	types "github.com/jontk/slurm-client/api"
	adapterbase "github.com/jontk/slurm-client/internal/adapters/base"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
	"github.com/jontk/slurm-client/pkg/errors"
)

// AccountAdapter implements the AccountAdapter interface for v0.0.41
type AccountAdapter struct {
	*adapterbase.BaseManager
	client *api.ClientWithResponses
}

// NewAccountAdapter creates a new Account adapter for v0.0.41
func NewAccountAdapter(client *api.ClientWithResponses) *AccountAdapter {
	return &AccountAdapter{
		BaseManager: adapterbase.NewBaseManager("v0.0.41", "Account"),
		client:      client,
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
		if len(opts.Descriptions) > 0 {
			descStr := strings.Join(opts.Descriptions, ",")
			params.Description = &descStr
		}
		if opts.WithDeleted {
			withDeleted := "true"
			params.DELETED = &withDeleted
		}
		if opts.WithAssocs {
			withAssocs := "true"
			params.WithAssociations = &withAssocs
		}
		if opts.WithCoords {
			withCoords := "true"
			params.WithCoordinators = &withCoords
		}
	}
	// Make the API call
	resp, err := a.client.SlurmdbV0041GetAccountsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}
	// Handle response
	if resp.HTTPResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: failed to list accounts", resp.HTTPResponse.StatusCode)
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected nil response")
	}
	// Convert response to common types
	accountList := &types.AccountList{
		Accounts: make([]types.Account, 0, len(resp.JSON200.Accounts)),
		Total:    len(resp.JSON200.Accounts),
	}
	for _, apiAccount := range resp.JSON200.Accounts {
		account, err := a.convertAPIAccountToCommon(apiAccount)
		if err != nil {
			// Log the error but continue processing other accounts
			continue
		}
		accountList.Accounts = append(accountList.Accounts, *account)
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
	if err := a.ValidateResourceName(name, "name"); err != nil {
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
		return nil, fmt.Errorf("failed to get account %s: %w", name, err)
	}
	// Handle response
	if resp.HTTPResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: failed to get account %s", resp.HTTPResponse.StatusCode, name)
	}
	if resp.JSON200 == nil || len(resp.JSON200.Accounts) == 0 {
		return nil, fmt.Errorf("account %s not found", name)
	}
	// Convert the first account in the response
	account, err := a.convertAPIAccountToCommon(resp.JSON200.Accounts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to convert account %s: %w", name, err)
	}
	return account, nil
}

// Create creates a new account
func (a *AccountAdapter) Create(ctx context.Context, account *types.AccountCreate) (*types.AccountCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	// Validate account
	if account == nil {
		return nil, fmt.Errorf("account cannot be nil")
	}
	if err := a.ValidateResourceName("account name", account.Name); err != nil {
		return nil, err
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}
	// Convert to API request using JSON marshaling workaround
	reqBody, err := a.convertAccountCreateToAPI(account)
	if err != nil {
		return nil, a.WrapError(err, "failed to convert account create request")
	}

	// Make the API call
	resp, err := a.client.SlurmdbV0041PostAccountsWithResponse(ctx, reqBody)
	if err != nil {
		return nil, a.WrapError(err, "failed to create account")
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}

	// Check for errors in response
	if resp.JSON200 != nil && resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		errMsgs := make([]string, 0)
		for _, apiErr := range *resp.JSON200.Errors {
			if apiErr.Error != nil {
				errMsgs = append(errMsgs, *apiErr.Error)
			}
		}
		if len(errMsgs) > 0 {
			return nil, fmt.Errorf("account creation failed: %v", errMsgs)
		}
	}

	return &types.AccountCreateResponse{
		AccountName: account.Name,
	}, nil
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
		return fmt.Errorf("account update cannot be nil")
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

	// Convert to API request using JSON marshaling workaround
	reqBody, err := a.convertAccountUpdateToAPI(existingAccount)
	if err != nil {
		return a.WrapError(err, "failed to convert account update request")
	}

	// Make the API call
	resp, err := a.client.SlurmdbV0041PostAccountsWithResponse(ctx, reqBody)
	if err != nil {
		return a.WrapError(err, "failed to update account")
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
	}

	// Check for errors in response
	if resp.JSON200 != nil && resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		errMsgs := make([]string, 0)
		for _, apiErr := range *resp.JSON200.Errors {
			if apiErr.Error != nil {
				errMsgs = append(errMsgs, *apiErr.Error)
			}
		}
		if len(errMsgs) > 0 {
			return fmt.Errorf("account update failed: %v", errMsgs)
		}
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
		return fmt.Errorf("failed to delete account %s: %w", name, err)
	}
	// Handle response
	if resp.HTTPResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: request failed", resp.HTTPResponse.StatusCode)
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
	if err := a.ValidateResourceName(name, "name"); err != nil {
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
	// TODO: Fix this when the correct API type is found
	// Create association request
	// assocReq := api.SlurmdbV0041PostAccountsAssociationJSONBody{
	//     Accounts: &[]string{accountName},
	//     Users:    &[]string{userName},
	// }
	return fmt.Errorf("AddUser not implemented for v0.0.41 - association API type not found")
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
	return errors.NewNotImplementedError("Set Account Coordinators", "v0.0.41")
}

// CreateAssociation creates associations for accounts (not supported in v0.0.41)
func (a *AccountAdapter) CreateAssociation(ctx context.Context, req *types.AccountAssociationRequest) (*types.AssociationCreateResponse, error) {
	return nil, a.HandleNotImplemented("CreateAssociation", "v0.0.41")
}
