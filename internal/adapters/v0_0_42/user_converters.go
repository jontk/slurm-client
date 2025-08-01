// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

// convertAPIUserToCommon converts a v0.0.42 API User to common User type
func (a *UserAdapter) convertAPIUserToCommon(apiUser api.V0042User) (*types.User, error) {
	user := &types.User{
		Name: apiUser.Name,
	}

	// Old name - not in common User type
	// Skip OldName handling

	// Default account and wckey
	if apiUser.Default != nil {
		if apiUser.Default.Account != nil {
			user.DefaultAccount = *apiUser.Default.Account
		}
		if apiUser.Default.Wckey != nil {
			user.DefaultWCKey = *apiUser.Default.Wckey
		}
	}

	// Administrator level
	if apiUser.AdministratorLevel != nil {
		// Convert to AdminLevel type
		adminStr := a.convertAdminLevelToString(apiUser.AdministratorLevel)
		user.AdminLevel = types.AdminLevel(adminStr)
	}

	// Flags - not in common User type
	// Skip Flags handling

	// Coordinators
	if apiUser.Coordinators != nil && len(*apiUser.Coordinators) > 0 {
		coordinators := make([]types.UserCoordinator, 0, len(*apiUser.Coordinators))
		for _, coord := range *apiUser.Coordinators {
			// coord has Name field as string, not pointer
			if coord.Name != "" {
				coordinators = append(coordinators, types.UserCoordinator{
					AccountName: coord.Name,
				})
			}
		}
		user.Coordinators = coordinators
	}

	// Associations
	if apiUser.Associations != nil && len(*apiUser.Associations) > 0 {
		// Convert associations to user associations
		associations := make([]types.UserAssociation, 0, len(*apiUser.Associations))
		
		// Extract account list from associations
		accountMap := make(map[string]bool)
		for _, assoc := range *apiUser.Associations {
			userAssoc := types.UserAssociation{}
			if assoc.Account != nil {
				userAssoc.AccountName = *assoc.Account
				accountMap[*assoc.Account] = true
			}
			if assoc.Cluster != nil {
				userAssoc.Cluster = *assoc.Cluster
			}
			if assoc.Partition != nil {
				userAssoc.Partition = *assoc.Partition
			}
			// User field is string
			userAssoc.UserName = assoc.User
			associations = append(associations, userAssoc)
		}
		user.Associations = associations
		
		// Convert accounts to list
		user.Accounts = make([]string, 0, len(accountMap))
		for account := range accountMap {
			user.Accounts = append(user.Accounts, account)
		}
	}

	// WCKeys
	if apiUser.Wckeys != nil && len(*apiUser.Wckeys) > 0 {
		wckeys := make([]string, 0, len(*apiUser.Wckeys))
		for _, wckey := range *apiUser.Wckeys {
			// wckey.Name is string, not pointer
			if wckey.Name != "" {
				wckeys = append(wckeys, wckey.Name)
			}
		}
		user.WCKeys = wckeys
	}

	return user, nil
}

// convertAdminLevelToString converts admin level to string representation
func (a *UserAdapter) convertAdminLevelToString(level *api.V0042AdminLvl) string {
	if level == nil || len(*level) == 0 {
		return "None"
	}
	
	// Return the first/primary admin level
	return (*level)[0]
}

// convertCommonUserCreateToAPI converts common user create request to v0.0.42 API format
func (a *UserAdapter) convertCommonUserCreateToAPI(req *types.UserCreateRequest) (*api.SlurmdbV0042PostUsersJSONRequestBody, error) {
	users := []api.V0042User{
		{
			Name: req.Name,
		},
	}
	apiReq := &api.SlurmdbV0042PostUsersJSONRequestBody{
		Users: users,
	}

	user := &apiReq.Users[0]

	// Default account and wckey
	if req.DefaultAccount != "" || req.DefaultWCKey != "" {
		user.Default = &struct {
			Account *string `json:"account,omitempty"`
			Wckey   *string `json:"wckey,omitempty"`
		}{}
		
		if req.DefaultAccount != "" {
			user.Default.Account = &req.DefaultAccount
		}
		if req.DefaultWCKey != "" {
			user.Default.Wckey = &req.DefaultWCKey
		}
	}

	// Administrator level
	if req.AdminLevel != "" && req.AdminLevel != "None" {
		// Convert AdminLevel type to string array
		adminLevel := []string{string(req.AdminLevel)}
		user.AdministratorLevel = &adminLevel
	}

	// Flags - not in UserCreate type
	// Skip Flags handling

	// Coordinators - not in UserCreate type
	// Skip Coordinators handling

	return apiReq, nil
}
