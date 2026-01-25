// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
	"github.com/jontk/slurm-client/internal/common/types"
)

// convertAPIAccountToCommon converts a v0.0.40 API Account to common Account type
func (a *AccountAdapter) convertAPIAccountToCommon(apiAccount api.V0040Account) *types.Account {
	account := &types.Account{}

	// Basic fields (v0.0.40 has direct string fields, not pointers)
	account.Name = apiAccount.Name
	account.Description = apiAccount.Description
	account.Organization = apiAccount.Organization

	// Set deleted flag if account has DELETED flag
	if apiAccount.Flags != nil && len(*apiAccount.Flags) > 0 {
		for _, flag := range *apiAccount.Flags {
			if flag == "DELETED" {
				account.Deleted = true
				break
			}
		}
	}

	// Coordinators - v0.0.40 API has complex coordinator structure, but common type just uses names
	if apiAccount.Coordinators != nil {
		account.Coordinators = make([]string, 0, len(*apiAccount.Coordinators))
		for _, apiCoord := range *apiAccount.Coordinators {
			// V0040Coord has Name field as a direct string
			account.Coordinators = append(account.Coordinators, apiCoord.Name)
		}
	}

	// Associations - Common Account type doesn't include associations
	// Associations would be handled separately via the Association manager
	// Skip this field for v0.0.40 basic account conversion

	return account
}

// convertCommonAccountCreateToAPI converts common AccountCreate to v0.0.40 API format
func (a *AccountAdapter) convertCommonAccountCreateToAPI(account *types.AccountCreate) *api.V0040Account {
	apiAccount := &api.V0040Account{}

	// Basic fields (v0.0.40 expects direct strings, not pointers)
	apiAccount.Name = account.Name
	apiAccount.Description = account.Description
	apiAccount.Organization = account.Organization

	// Parent account - AccountCreate uses ParentName, not Parent
	// v0.0.40 may handle parent differently, might need to set via associations
	// For now, we'll skip this as it's typically handled by the accounting system
	// TODO: Implement parent account handling for v0.0.40
	if account.ParentName != "" {
	}

	// Flags - Common AccountCreate type doesn't have Flags field
	// Skip flags for create operation in v0.0.40

	// Coordinators - convert string names to V0040Coord structs
	if len(account.Coordinators) > 0 {
		coords := make([]api.V0040Coord, len(account.Coordinators))
		for i, coordName := range account.Coordinators {
			coords[i] = api.V0040Coord{
				Name: coordName,
				// Direct can be set to true for explicitly assigned coordinators
				Direct: boolPtr(true),
			}
		}
		apiAccount.Coordinators = &coords
	}

	return apiAccount
}

// convertCommonAccountUpdateToAPI converts common AccountUpdate to v0.0.40 API format
func (a *AccountAdapter) convertCommonAccountUpdateToAPI(existingAccount *types.Account, update *types.AccountUpdate) *api.V0040Account {
	apiAccount := &api.V0040Account{}

	// Name (required) - v0.0.40 expects direct string
	apiAccount.Name = existingAccount.Name

	// Apply updates
	if update.Description != nil {
		apiAccount.Description = *update.Description
	} else {
		apiAccount.Description = existingAccount.Description
	}
	if update.Organization != nil {
		apiAccount.Organization = *update.Organization
	} else {
		apiAccount.Organization = existingAccount.Organization
	}

	// Flags - Common AccountUpdate type doesn't have Flags field
	// Skip flags for update operation in v0.0.40

	return apiAccount
}
