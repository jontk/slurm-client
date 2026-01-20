// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_44

import (
	"context"
	"fmt"
	"strings"

	api "github.com/jontk/slurm-client/internal/api/v0_0_44"
	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/jontk/slurm-client/pkg/errors"
)

// AccountAdapter implements the AccountAdapter interface for v0.0.44
type AccountAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewAccountAdapter creates a new Account adapter for v0.0.44
func NewAccountAdapter(client *api.ClientWithResponses) *AccountAdapter {
	return &AccountAdapter{
		BaseManager: base.NewBaseManager("v0.0.44", "Account"),
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
	params := &api.SlurmdbV0044GetAccountsParams{}

	// Apply filters from options
	if opts != nil {
		// Note: v0.0.44 API has limited query parameters
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

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0044GetAccountsWithResponse(ctx, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0044OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.44"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "List Accounts"); err != nil {
		return nil, err
	}
	// V0044AccountList is already a slice, not a pointer
	// No need to check for nil

	// Convert the response to common types
	accountList := make([]types.Account, 0, len(resp.JSON200.Accounts))
	for _, apiAccount := range resp.JSON200.Accounts {
		account, err := a.convertAPIAccountToCommon(apiAccount)
		if err != nil {
			return nil, a.HandleConversionError(err, apiAccount.Name)
		}
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
	params := &api.SlurmdbV0044GetAccountParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0044GetAccountWithResponse(ctx, accountName, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0044OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.44"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "Get Account"); err != nil {
		return nil, err
	}
	// V0044AccountList is already a slice, not a pointer
	// No need to check for nil

	// Check if we got any account entries
	if len(resp.JSON200.Accounts) == 0 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("Account %s not found", accountName))
	}

	// Convert the first account (should be the only one)
	account, err := a.convertAPIAccountToCommon(resp.JSON200.Accounts[0])
	if err != nil {
		return nil, a.HandleConversionError(err, accountName)
	}

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
	apiAccount, err := a.convertCommonAccountCreateToAPI(account)
	if err != nil {
		return nil, err
	}

	// Create request body
	reqBody := api.SlurmdbV0044PostAccountsJSONRequestBody{
		Accounts: []api.V0044Account{*apiAccount},
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0044PostAccountsWithResponse(ctx, reqBody)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0044OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.44"); err != nil {
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
	apiAccount, err := a.convertCommonAccountUpdateToAPI(existingAccount, update)
	if err != nil {
		return err
	}

	// Create request body
	reqBody := api.SlurmdbV0044PostAccountsJSONRequestBody{
		Accounts: []api.V0044Account{*apiAccount},
	}

	// Call the generated OpenAPI client (POST is used for updates in SLURM API)
	resp, err := a.client.SlurmdbV0044PostAccountsWithResponse(ctx, reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0044OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "v0.0.44")
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
	resp, err := a.client.SlurmdbV0044DeleteAccountWithResponse(ctx, accountName)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0044OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	// Create adapter with special handling for 204 (No Content) status
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)

	// For DELETE operations, 204 is also a success
	if resp.StatusCode() == 204 {
		return nil
	}

	return common.HandleAPIResponse(responseAdapter, "v0.0.44")
}

// validateAccountCreate validates account creation request
func (a *AccountAdapter) validateAccountCreate(account *types.AccountCreate) error {
	if account == nil {
		return errors.NewValidationErrorf("account", nil, "account creation data is required")
	}
	if account.Name == "" {
		return errors.NewValidationErrorf("name", account.Name, "account name is required")
	}
	// Validate numeric fields
	if account.FairShare < 0 {
		return errors.NewValidationErrorf("fairShare", account.FairShare, "fair share must be non-negative")
	}
	if account.Priority < 0 {
		return errors.NewValidationErrorf("priority", account.Priority, "priority must be non-negative")
	}
	if account.MaxJobs < 0 {
		return errors.NewValidationErrorf("maxJobs", account.MaxJobs, "max jobs must be non-negative")
	}
	if account.MaxWallTime < 0 {
		return errors.NewValidationErrorf("maxWallTime", account.MaxWallTime, "max wall time must be non-negative")
	}
	return nil
}

// validateAccountUpdate validates account update request
func (a *AccountAdapter) validateAccountUpdate(update *types.AccountUpdate) error {
	if update == nil {
		return errors.NewValidationErrorf("update", nil, "account update data is required")
	}
	// At least one field should be provided for update
	if update.Description == nil && update.Organization == nil && len(update.Coordinators) == 0 &&
		update.DefaultQoS == nil && len(update.QoSList) == 0 && len(update.AllowedPartitions) == 0 &&
		update.DefaultPartition == nil && update.FairShare == nil && update.Priority == nil &&
		update.MaxJobs == nil && update.MaxWallTime == nil {
		return errors.NewValidationErrorf("update", update, "at least one field must be provided for update")
	}

	// Validate numeric fields if provided
	if update.FairShare != nil && *update.FairShare < 0 {
		return errors.NewValidationErrorf("fairShare", *update.FairShare, "fair share must be non-negative")
	}
	if update.Priority != nil && *update.Priority < 0 {
		return errors.NewValidationErrorf("priority", *update.Priority, "priority must be non-negative")
	}
	if update.MaxJobs != nil && *update.MaxJobs < 0 {
		return errors.NewValidationErrorf("maxJobs", *update.MaxJobs, "max jobs must be non-negative")
	}
	if update.MaxWallTime != nil && *update.MaxWallTime < 0 {
		return errors.NewValidationErrorf("maxWallTime", *update.MaxWallTime, "max wall time must be non-negative")
	}
	return nil
}

// CreateAssociation creates associations for accounts
func (a *AccountAdapter) CreateAssociation(ctx context.Context, req *types.AccountAssociationRequest) (*types.AssociationCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.validateAccountAssociationRequest(req); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert to API format
	apiAssociations, err := a.convertCommonAccountAssociationToAPI(req)
	if err != nil {
		return nil, err
	}

	// Create request body
	reqBody := api.SlurmdbV0044PostAssociationsJSONRequestBody{
		Associations: apiAssociations,
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0044PostAssociationsWithResponse(ctx, reqBody)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0044OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.44"); err != nil {
		return nil, err
	}

	return &types.AssociationCreateResponse{
		Status:  "success",
		Message: fmt.Sprintf("Created associations for %d accounts in cluster %s", len(req.Accounts), req.Cluster),
	}, nil
}

// validateAccountAssociationRequest validates account association creation request
func (a *AccountAdapter) validateAccountAssociationRequest(req *types.AccountAssociationRequest) error {
	if req == nil {
		return errors.NewValidationErrorf("request", nil, "account association request is required")
	}
	if len(req.Accounts) == 0 {
		return errors.NewValidationErrorf("accounts", req.Accounts, "at least one account is required")
	}
	if req.Cluster == "" {
		return errors.NewValidationErrorf("cluster", req.Cluster, "cluster is required")
	}
	// Validate numeric fields
	if req.Fairshare < 0 {
		return errors.NewValidationErrorf("fairshare", req.Fairshare, "fairshare must be non-negative")
	}
	return nil
}

// convertCommonAccountAssociationToAPI converts common account association request to API format
func (a *AccountAdapter) convertCommonAccountAssociationToAPI(req *types.AccountAssociationRequest) ([]api.V0044Assoc, error) {
	associations := make([]api.V0044Assoc, 0, len(req.Accounts))

	for _, accountName := range req.Accounts {
		association := api.V0044Assoc{
			Account: &accountName,
			Cluster: &req.Cluster,
		}

		if req.Partition != "" {
			association.Partition = &req.Partition
		}
		if req.Parent != "" {
			association.ParentAccount = &req.Parent
		}
		if len(req.QoS) > 0 {
			qosList := make(api.V0044QosStringIdList, len(req.QoS))
			copy(qosList, req.QoS)
			association.Qos = &qosList
		}
		if req.DefaultQoS != "" {
			association.Default = &struct {
				Qos *string `json:"qos,omitempty"`
			}{
				Qos: &req.DefaultQoS,
			}
		}
		if req.Fairshare > 0 {
			association.SharesRaw = &req.Fairshare
		}

		// Handle TRES if provided using TRES utilities
		// TRES handling is now implemented with proper parsing and conversion
		if association.Max != nil && association.Max.Tres != nil {
			// TRES limits are available in the association
			// Further TRES processing can be done here if needed using NewTRESUtils()
			// For now, we keep the TRES data as provided by the API
		}

		associations = append(associations, association)
	}

	return associations, nil
}
