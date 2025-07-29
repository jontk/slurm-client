package v0_0_41

import (
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
)

// convertAPIAccountToCommon converts a v0.0.41 API Account to common Account type
func (a *AccountAdapter) convertAPIAccountToCommon(apiAccount interface{}) (*types.Account, error) {
	// Type assertion to handle the anonymous struct - simplified to avoid undefined types
	accountData, ok := apiAccount.(struct {
		Name            *string   `json:"name,omitempty"`
		Description     *string   `json:"description,omitempty"`
		Organization    *string   `json:"organization,omitempty"`
		Coordinators    *[]struct {
			Name *string `json:"name,omitempty"`
		} `json:"coordinators,omitempty"`
		Flags           *[]string `json:"flags,omitempty"`
	})
	if !ok {
		return nil, fmt.Errorf("unexpected account data type")
	}

	account := &types.Account{}

	// Basic fields
	if accountData.Name != nil {
		account.Name = *accountData.Name
	}
	if accountData.Description != nil {
		account.Description = *accountData.Description
	}
	if accountData.Organization != nil {
		account.Organization = *accountData.Organization
	}

	// Handle flags as strings to avoid undefined API types
	if accountData.Flags != nil {
		for _, flag := range *accountData.Flags {
			if flag == "DELETED" || flag == "deleted" {
				account.Deleted = true
			}
		}
	}

	// Coordinators
	if accountData.Coordinators != nil {
		coordinators := make([]string, 0, len(*accountData.Coordinators))
		for _, coord := range *accountData.Coordinators {
			if coord.Name != nil {
				coordinators = append(coordinators, *coord.Name)
			}
		}
		account.Coordinators = coordinators
	}

	return account, nil
}

// convertCommonToAPIAccount converts common Account to v0.0.41 API request
func (a *AccountAdapter) convertCommonToAPIAccount(account *types.Account) interface{} {
	// Return a simplified structure that matches what the API expects
	// without relying on undefined specific types
	req := struct {
		Accounts []struct {
			Name         *string   `json:"name,omitempty"`
			Description  *string   `json:"description,omitempty"`
			Organization *string   `json:"organization,omitempty"`
			Coordinators *[]struct {
				Name *string `json:"name,omitempty"`
			} `json:"coordinators,omitempty"`
			Flags *[]string `json:"flags,omitempty"`
		} `json:"accounts"`
	}{
		Accounts: []struct {
			Name         *string   `json:"name,omitempty"`
			Description  *string   `json:"description,omitempty"`
			Organization *string   `json:"organization,omitempty"`
			Coordinators *[]struct {
				Name *string `json:"name,omitempty"`
			} `json:"coordinators,omitempty"`
			Flags *[]string `json:"flags,omitempty"`
		}{{}},
	}

	acc := &req.Accounts[0]

	// Set basic fields
	if account.Name != "" {
		acc.Name = &account.Name
	}
	if account.Description != "" {
		acc.Description = &account.Description
	}
	if account.Organization != "" {
		acc.Organization = &account.Organization
	}

	// Handle flags as strings
	if account.Deleted {
		flags := []string{"DELETED"}
		acc.Flags = &flags
	}

	// Convert coordinators
	if len(account.Coordinators) > 0 {
		coords := make([]struct {
			Name *string `json:"name,omitempty"`
		}, 0, len(account.Coordinators))
		for _, coordName := range account.Coordinators {
			coordNameCopy := coordName // Create a copy to avoid pointer issues
			coords = append(coords, struct {
				Name *string `json:"name,omitempty"`
			}{
				Name: &coordNameCopy,
			})
		}
		acc.Coordinators = &coords
	}

	return req
}

// convertCommonToAPIAccountUpdate converts common AccountUpdate to v0.0.41 API request
func (a *AccountAdapter) convertCommonToAPIAccountUpdate(update *types.AccountUpdate) interface{} {
	// For v0.0.41, updates are done by sending the full account object
	// Return a simplified structure for update operations
	req := struct {
		Accounts []struct {
			Description  *string   `json:"description,omitempty"`
			Organization *string   `json:"organization,omitempty"`
			Coordinators []string  `json:"coordinators,omitempty"`
		} `json:"accounts"`
	}{
		Accounts: []struct {
			Description  *string   `json:"description,omitempty"`
			Organization *string   `json:"organization,omitempty"`
			Coordinators []string  `json:"coordinators,omitempty"`
		}{{}},
	}

	acc := &req.Accounts[0]
	if update.Description != nil {
		acc.Description = update.Description
	}
	if update.Organization != nil {
		acc.Organization = update.Organization
	}
	if len(update.Coordinators) > 0 {
		acc.Coordinators = update.Coordinators
	}

	return req
}