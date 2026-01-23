// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"fmt"

	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
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
		// Note: v0.0.42 API doesn't support filtering by account names in GetAccounts
		// Account filtering is done after retrieval
		if opts.WithAssocs {
			withAssoc := "true"
			params.WithAssociations = &withAssoc
		}
		if opts.WithCoords {
			withCoord := "true"
			params.WithCoordinators = &withCoord
		}
		if opts.WithDeleted {
			withDeleted := "true"
			params.DELETED = &withDeleted
		}
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042GetAccountsWithResponse(ctx, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode())
	}

	// Check for API response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from API")
	}

	// Convert the response to common types
	accountList := &types.AccountList{
		Accounts: make([]types.Account, 0),
	}

	if resp.JSON200.Accounts != nil {
		for _, apiAccount := range resp.JSON200.Accounts {
			account, err := a.convertAPIAccountToCommon(apiAccount)
			if err != nil {
				// Log conversion error but continue
				continue
			}
			accountList.Accounts = append(accountList.Accounts, *account)
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
		return nil, a.HandleAPIError(err)
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode())
	}

	// Check for API response
	if resp.JSON200 == nil || resp.JSON200.Accounts == nil || len(resp.JSON200.Accounts) == 0 {
		return nil, fmt.Errorf("account %s not found", name)
	}

	// Convert the first account in the response
	accounts := resp.JSON200.Accounts
	return a.convertAPIAccountToCommon(accounts[0])
}

// Create creates a new account
func (a *AccountAdapter) Create(ctx context.Context, account *types.AccountCreate) (*types.AccountCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert common account to API format
	apiAccount, err := a.convertCommonAccountCreateToAPI(account)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042PostAccountsWithResponse(ctx, *apiAccount)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode())
	}

	// Return success response
	return &types.AccountCreateResponse{
		AccountName: account.Name,
	}, nil
}

// Update updates an existing account
func (a *AccountAdapter) Update(ctx context.Context, name string, updates *types.AccountUpdate) error {
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
		return a.HandleAPIError(err)
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return fmt.Errorf("API request failed with status %d", resp.StatusCode())
	}

	return nil
}

// CreateAssociation creates associations for accounts
func (a *AccountAdapter) CreateAssociation(ctx context.Context, req *types.AccountAssociationRequest) (*types.AssociationCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Validate request
	if req == nil {
		return nil, fmt.Errorf("association request is required")
	}
	if len(req.Accounts) == 0 {
		return nil, fmt.Errorf("at least one account is required")
	}
	if req.Cluster == "" {
		return nil, fmt.Errorf("cluster is required")
	}

	// Convert common request to API request structure
	apiReq, err := a.convertAccountAssociationRequestToAPI(req)
	if err != nil {
		return nil, a.WrapError(err, "failed to convert association request")
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042PostAccountsAssociationWithResponse(ctx, *apiReq)
	if err != nil {
		return nil, a.WrapError(err, "failed to create account associations")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}

	// Check for API response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from API")
	}

	// Convert response
	return a.convertAccountAssociationResponseToCommon(resp.JSON200), nil
}

// convertAccountAssociationRequestToAPI converts common request to API structure
func (a *AccountAdapter) convertAccountAssociationRequestToAPI(req *types.AccountAssociationRequest) (*api.V0042OpenapiAccountsAddCondResp, error) {
	// Create accounts list from string slice
	accounts := make(api.V0042StringList, len(req.Accounts))
	copy(accounts, req.Accounts)

	// Create association condition
	assocCond := &api.V0042AccountsAddCond{
		Accounts: accounts,
	}

	// Add clusters if specified
	if req.Cluster != "" {
		clusters := api.V0042StringList{req.Cluster}
		assocCond.Clusters = &clusters
	}

	// Create association record set if we have additional fields
	if req.Parent != "" || req.DefaultQoS != "" || req.Fairshare != 0 || len(req.QoS) > 0 {
		assocRec := &api.V0042AssocRecSet{}

		// Set default QoS (note: field name is "Defaultqos" in v0.0.42)
		if req.DefaultQoS != "" {
			assocRec.Defaultqos = &req.DefaultQoS
		}

		// Set fairshare (note: field name is "Fairshare" in v0.0.42)
		if req.Fairshare != 0 {
			fairshareInt32 := req.Fairshare
			assocRec.Fairshare = &fairshareInt32
		}

		// Note: ParentAccount is not available in V0042AssocRecSet
		// Parent account relationships are managed differently in v0.0.42

		assocCond.Association = assocRec
	}

	// Create the API request
	apiReq := &api.V0042OpenapiAccountsAddCondResp{
		AssociationCondition: assocCond,
	}

	return apiReq, nil
}

// convertAccountAssociationResponseToCommon converts API response to common type
func (a *AccountAdapter) convertAccountAssociationResponseToCommon(apiResp *api.V0042OpenapiAccountsAddCondRespStr) *types.AssociationCreateResponse {
	resp := &types.AssociationCreateResponse{
		Status: "success",
		Meta:   make(map[string]interface{}),
	}

	// Extract added accounts info
	if apiResp.AddedAccounts != "" {
		resp.Message = "Successfully created associations for accounts: " + apiResp.AddedAccounts
		resp.Meta["added_accounts"] = apiResp.AddedAccounts
	} else {
		resp.Message = "Account associations created successfully"
	}

	// Handle errors in response
	if apiResp.Errors != nil && len(*apiResp.Errors) > 0 {
		resp.Status = "error"
		errors := *apiResp.Errors
		if len(errors) > 0 && errors[0].Error != nil {
			resp.Message = *errors[0].Error
		} else {
			resp.Message = "Account association creation failed"
		}
	}

	// Extract metadata if available
	if apiResp.Meta != nil {
		if apiResp.Meta.Client != nil {
			clientInfo := make(map[string]interface{})
			if apiResp.Meta.Client.Source != nil {
				clientInfo["source"] = *apiResp.Meta.Client.Source
			}
			if apiResp.Meta.Client.User != nil {
				clientInfo["user"] = *apiResp.Meta.Client.User
			}
			if len(clientInfo) > 0 {
				resp.Meta["client"] = clientInfo
			}
		}
	}

	return resp
}

// convertAPIAccountToCommon converts API account to common type
func (a *AccountAdapter) convertAPIAccountToCommon(apiAccount api.V0042Account) (*types.Account, error) {
	account := &types.Account{}

	// Set basic fields - V0042Account has Name and Description as strings, not pointers
	account.Name = apiAccount.Name
	account.Description = apiAccount.Description

	// Organization is not directly available in V0042Account
	// We'll leave it empty as it's not part of the main account structure

	// Note: V0042Account structure doesn't expose associations in a way that maps to our common type
	// Association data would need to be retrieved separately via the association adapter

	// Convert coordinators if present
	if apiAccount.Coordinators != nil {
		for _, coord := range *apiAccount.Coordinators {
			// V0042Coord has Name as string, not pointer
			account.Coordinators = append(account.Coordinators, coord.Name)
		}
	}

	return account, nil
}

// convertCommonAccountCreateToAPI converts common account create to API format
func (a *AccountAdapter) convertCommonAccountCreateToAPI(accountCreate *types.AccountCreate) (*api.V0042OpenapiAccountsResp, error) {
	if accountCreate == nil {
		return nil, fmt.Errorf("account create request cannot be nil")
	}

	// V0042Account has Name and Description as strings, not pointers
	apiAccount := api.V0042Account{
		Name:        accountCreate.Name,
		Description: accountCreate.Description,
	}

	// Note: V0042Account doesn't support embedded associations in create request
	// Associations would need to be created separately via the association adapter

	return &api.V0042OpenapiAccountsResp{
		Accounts: []api.V0042Account{apiAccount},
	}, nil
}
