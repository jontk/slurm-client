// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_44

import (
	"fmt"

	api "github.com/jontk/slurm-client/internal/api/v0_0_44"
	"github.com/jontk/slurm-client/internal/common/types"
)

// convertAPIAccountToCommon converts a v0.0.44 API Account to common Account type
func (a *AccountAdapter) convertAPIAccountToCommon(apiAccount api.V0044Account) (*types.Account, error) {
	account := &types.Account{}

	// Basic fields - in v0.0.44 these are direct strings, not pointers
	account.Name = apiAccount.Name
	account.Description = apiAccount.Description
	account.Organization = apiAccount.Organization

	// Coordinators
	if apiAccount.Coordinators != nil && len(*apiAccount.Coordinators) > 0 {
		coordinators := make([]string, len(*apiAccount.Coordinators))
		for i, coord := range *apiAccount.Coordinators {
			// coord.Name is a string, not a pointer
			coordinators[i] = coord.Name
		}
		account.Coordinators = coordinators
	}

	// Note: The v0.0.44 API Account struct only contains:
	// - Associations, Coordinators, Description, Flags, Name, Organization
	// It does not have QoS, priority, limits, or other advanced fields.
	// These would typically come from associations or require separate API calls.

	// Flags
	if apiAccount.Flags != nil {
		for _, flag := range *apiAccount.Flags {
			switch flag {
			case api.V0044AccountFlagsDELETED:
				account.Deleted = true
				// Note: V0044AccountFlagsDEFAULT does not exist in this API version
			}
		}
	}

	return account, nil
}

// convertCommonAccountCreateToAPI converts common AccountCreate type to v0.0.44 API format
func (a *AccountAdapter) convertCommonAccountCreateToAPI(create *types.AccountCreate) (*api.V0044Account, error) {
	apiAccount := &api.V0044Account{}

	// Required fields - these are non-pointer strings in v0.0.44
	apiAccount.Name = create.Name

	// Ensure Description is not empty (required field)
	if create.Description != "" {
		apiAccount.Description = create.Description
	} else {
		// Provide default description if empty
		apiAccount.Description = fmt.Sprintf("Account %s", create.Name)
	}

	// Ensure Organization is not empty (required field)
	if create.Organization != "" {
		apiAccount.Organization = create.Organization
	} else {
		// Provide default organization if empty
		apiAccount.Organization = "default"
	}

	// Coordinators
	if len(create.Coordinators) > 0 {
		coordinators := make([]api.V0044Coord, len(create.Coordinators))
		for i, coord := range create.Coordinators {
			coordinators[i] = api.V0044Coord{
				Name: coord,
			}
		}
		apiCoordList := coordinators
		apiAccount.Coordinators = &apiCoordList
	}

	// Note: v0.0.44 API Account struct only supports basic fields:
	// - Name, Description, Organization, Coordinators, Associations, Flags
	// Other fields like QoS, limits, priority are not available in this API version.
	// These settings may need to be handled through account associations or
	// separate API endpoints.

	return apiAccount, nil
}

// convertCommonAccountUpdateToAPI converts common AccountUpdate to v0.0.44 API format
func (a *AccountAdapter) convertCommonAccountUpdateToAPI(existing *types.Account, update *types.AccountUpdate) (*api.V0044Account, error) {
	apiAccount := &api.V0044Account{}

	// Always include the account name for updates
	apiAccount.Name = existing.Name

	// Apply updates to fields
	description := existing.Description
	if update.Description != nil {
		description = *update.Description
	}
	apiAccount.Description = description

	organization := existing.Organization
	if update.Organization != nil {
		organization = *update.Organization
	}
	apiAccount.Organization = organization

	// Coordinators
	coordinators := existing.Coordinators
	if len(update.Coordinators) > 0 {
		coordinators = update.Coordinators
	}
	if len(coordinators) > 0 {
		apiCoordinators := make([]api.V0044Coord, len(coordinators))
		for i, coord := range coordinators {
			apiCoordinators[i] = api.V0044Coord{
				Name: coord,
			}
		}
		apiCoordList := apiCoordinators
		apiAccount.Coordinators = &apiCoordList
	}

	// Note: v0.0.44 API Account struct only supports basic fields:
	// - Name, Description, Organization, Coordinators, Associations, Flags
	// Other fields like QoS, limits, priority are not available in this API version.

	return apiAccount, nil
}
