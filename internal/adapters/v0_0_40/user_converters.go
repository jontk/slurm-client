// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
)

// convertAPIUserToCommon converts a v0.0.40 API User to common User type
func (a *UserAdapter) convertAPIUserToCommon(apiUser api.V0040User) (*types.User, error) {
	user := &types.User{}

	// Basic fields
	// Name is string, not pointer in v0.0.40
	user.Name = apiUser.Name
	// v0.0.40 doesn't have Uid field
	if apiUser.Default != nil && apiUser.Default.Account != nil {
		user.DefaultAccount = *apiUser.Default.Account
	}
	if apiUser.Default != nil && apiUser.Default.Wckey != nil {
		user.DefaultWCKey = *apiUser.Default.Wckey
	}

	// Admin level
	if apiUser.AdministratorLevel != nil && len(*apiUser.AdministratorLevel) > 0 {
		// Convert admin level array to single admin level
		user.AdminLevel = types.AdminLevel((*apiUser.AdministratorLevel)[0])
	}

	// v0.0.40 common User type doesn't have Flags field

	// Coordinators
	if apiUser.Coordinators != nil {
		user.Coordinators = make([]types.UserCoordinator, 0, len(*apiUser.Coordinators))
		for _, apiCoord := range *apiUser.Coordinators {
			// Name is string, not pointer in v0.0.40
			coord := types.UserCoordinator{
				// Using coordinator name as both account and coordinator
				AccountName: apiCoord.Name,
				Coordinator: apiCoord.Name,
			}
			user.Coordinators = append(user.Coordinators, coord)
		}
	}

	// Associations
	if apiUser.Associations != nil {
		user.Associations = make([]types.UserAssociation, 0, len(*apiUser.Associations))
		for _, apiAssoc := range *apiUser.Associations {
			assoc := types.UserAssociation{}
			if apiAssoc.Account != nil {
				assoc.AccountName = *apiAssoc.Account
			}
			if apiAssoc.Cluster != nil {
				assoc.Cluster = *apiAssoc.Cluster
			}
			if apiAssoc.Partition != nil {
				assoc.Partition = *apiAssoc.Partition
			}
			// User is string, not pointer in v0.0.40
			assoc.UserName = apiAssoc.User
			// v0.0.40 AssocShort doesn't have Defaults, Flags, Shares, or Priority fields
			// These would come from a full association fetch, not the short form
			user.Associations = append(user.Associations, assoc)
		}
	}

	return user, nil
}

// convertCommonUserCreateToAPI converts common UserCreate to v0.0.40 API format
func (a *UserAdapter) convertCommonUserCreateToAPI(user *types.UserCreate) (*api.V0040User, error) {
	apiUser := &api.V0040User{}

	// Basic fields
	apiUser.Name = user.Name

	// Default settings
	if user.DefaultAccount != "" || user.DefaultWCKey != "" {
		apiUser.Default = &struct {
			Account *string `json:"account,omitempty"`
			Wckey   *string `json:"wckey,omitempty"`
		}{
			Account: &user.DefaultAccount,
			Wckey:   &user.DefaultWCKey,
		}
	}

	// Admin level
	if user.AdminLevel != "" {
		adminLevel := api.V0040AdminLvl([]string{string(user.AdminLevel)})
		apiUser.AdministratorLevel = &adminLevel
	}

	return apiUser, nil
}

// convertCommonUserUpdateToAPI converts common UserUpdate to v0.0.40 API format
func (a *UserAdapter) convertCommonUserUpdateToAPI(existingUser *types.User, update *types.UserUpdate) (*api.V0040User, error) {
	apiUser := &api.V0040User{}

	// Name (required)
	apiUser.Name = existingUser.Name

	// Apply updates
	if update.DefaultAccount != nil || update.DefaultWCKey != nil {
		apiUser.Default = &struct {
			Account *string `json:"account,omitempty"`
			Wckey   *string `json:"wckey,omitempty"`
		}{
			Account: update.DefaultAccount,
			Wckey:   update.DefaultWCKey,
		}
	}

	// Admin level
	if update.AdminLevel != nil {
		adminLevel := api.V0040AdminLvl([]string{string(*update.AdminLevel)})
		apiUser.AdministratorLevel = &adminLevel
	}

	return apiUser, nil
}
