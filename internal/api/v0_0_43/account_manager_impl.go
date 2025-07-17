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
		return nil, errors.NewClientError("API client not initialized", nil)
	}

	// TODO: Implement actual API call when account endpoints are available in OpenAPI spec
	// For now, return NotImplementedError as the actual implementation requires
	// the generated client to have account-related methods
	return nil, errors.NewNotImplementedError("account listing", "v0.0.43")
}

// Get retrieves a specific account by name
func (a *AccountManagerImpl) Get(ctx context.Context, accountName string) (*interfaces.Account, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError("API client not initialized", nil)
	}

	if accountName == "" {
		return nil, errors.NewValidationError("account name is required", nil)
	}

	// TODO: Implement actual API call when account endpoints are available in OpenAPI spec
	return nil, errors.NewNotImplementedError("account retrieval", "v0.0.43")
}

// Create creates a new account
func (a *AccountManagerImpl) Create(ctx context.Context, account *interfaces.AccountCreate) (*interfaces.AccountCreateResponse, error) {
	if a.client == nil || a.client.apiClient == nil {
		return nil, errors.NewClientError("API client not initialized", nil)
	}

	if account == nil {
		return nil, errors.NewValidationError("account data is required", nil)
	}

	if account.Name == "" {
		return nil, errors.NewValidationError("account name is required", nil)
	}

	// Validate account hierarchy
	if account.ParentAccount != "" {
		// In a real implementation, we would verify the parent account exists
		// For now, just check it's not the same as the account being created
		if account.ParentAccount == account.Name {
			return nil, errors.NewValidationError("account cannot be its own parent", nil)
		}
	}

	// TODO: Implement actual API call when account endpoints are available in OpenAPI spec
	return nil, errors.NewNotImplementedError("account creation", "v0.0.43")
}

// Update updates an existing account
func (a *AccountManagerImpl) Update(ctx context.Context, accountName string, update *interfaces.AccountUpdate) error {
	if a.client == nil || a.client.apiClient == nil {
		return errors.NewClientError("API client not initialized", nil)
	}

	if accountName == "" {
		return errors.NewValidationError("account name is required", nil)
	}

	if update == nil {
		return errors.NewValidationError("update data is required", nil)
	}

	// TODO: Implement actual API call when account endpoints are available in OpenAPI spec
	return errors.NewNotImplementedError("account update", "v0.0.43")
}

// Delete deletes an account
func (a *AccountManagerImpl) Delete(ctx context.Context, accountName string) error {
	if a.client == nil || a.client.apiClient == nil {
		return errors.NewClientError("API client not initialized", nil)
	}

	if accountName == "" {
		return errors.NewValidationError("account name is required", nil)
	}

	// TODO: Implement actual API call when account endpoints are available in OpenAPI spec
	// Note: Account deletion may require special handling for:
	// - Child accounts (cascade vs prevent)
	// - Active users in the account
	// - Running jobs associated with the account
	return errors.NewNotImplementedError("account deletion", "v0.0.43")
}

// Helper function to validate TRES format
func validateTRES(tres map[string]int) error {
	// TRES (Trackable Resources) typically include: cpu, mem, energy, node, billing, fs/disk, vmem, pages
	// Values should be non-negative
	for resource, value := range tres {
		if value < 0 {
			return errors.NewValidationError(fmt.Sprintf("invalid TRES value for %s: must be non-negative", resource), nil)
		}
	}
	return nil
}