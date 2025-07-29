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

	// Old name (for tracking user renames)
	if apiUser.OldName != nil {
		user.OldName = *apiUser.OldName
	}

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
		user.AdminLevel = a.convertAdminLevelToString(apiUser.AdministratorLevel)
	}

	// Flags
	if apiUser.Flags != nil && len(*apiUser.Flags) > 0 {
		user.Flags = *apiUser.Flags
	}

	// Coordinators
	if apiUser.Coordinators != nil && len(*apiUser.Coordinators) > 0 {
		coordinators := make([]string, 0, len(*apiUser.Coordinators))
		for _, coord := range *apiUser.Coordinators {
			if coord.Name != nil {
				coordinators = append(coordinators, *coord.Name)
			}
		}
		user.CoordinatorOf = coordinators
	}

	// Associations
	if apiUser.Associations != nil && len(*apiUser.Associations) > 0 {
		// Store association count
		user.AssociationCount = len(*apiUser.Associations)
		
		// Extract account list from associations
		accountMap := make(map[string]bool)
		for _, assoc := range *apiUser.Associations {
			if assoc.Account != nil {
				accountMap[*assoc.Account] = true
			}
		}
		
		user.Accounts = make([]string, 0, len(accountMap))
		for account := range accountMap {
			user.Accounts = append(user.Accounts, account)
		}
		
		// Find primary association details
		for _, assoc := range *apiUser.Associations {
			// If this matches the default account, capture some details
			if assoc.Account != nil && *assoc.Account == user.DefaultAccount {
				if assoc.Qos != nil && len(*assoc.Qos) > 0 {
					user.DefaultQoS = (*assoc.Qos)[0]
				}
				if assoc.MaxJobs != nil && assoc.MaxJobs.Number > 0 {
					maxJobs := uint32(assoc.MaxJobs.Number)
					user.MaxJobs = &maxJobs
				}
				if assoc.MaxCpus != nil && assoc.MaxCpus.Number > 0 {
					maxCPUs := uint32(assoc.MaxCpus.Number)
					user.MaxCPUs = &maxCPUs
				}
				break
			}
		}
	}

	// WCKeys
	if apiUser.Wckeys != nil && len(*apiUser.Wckeys) > 0 {
		wckeys := make([]string, 0, len(*apiUser.Wckeys))
		for _, wckey := range *apiUser.Wckeys {
			if wckey.Name != nil {
				wckeys = append(wckeys, *wckey.Name)
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
	apiReq := &api.SlurmdbV0042PostUsersJSONRequestBody{
		Users: &[]api.V0042User{
			{
				Name: req.Name,
			},
		},
	}

	user := &(*apiReq.Users)[0]

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
		adminLevel := api.V0042AdminLvl{req.AdminLevel}
		user.AdministratorLevel = &adminLevel
	}

	// Flags
	if len(req.Flags) > 0 {
		flags := api.V0042UserFlags(req.Flags)
		user.Flags = &flags
	}

	// Coordinators
	if len(req.CoordinatorOf) > 0 {
		coords := make(api.V0042CoordList, len(req.CoordinatorOf))
		for i, coordName := range req.CoordinatorOf {
			coords[i] = api.V0042Coord{
				Name: &coordName,
			}
		}
		user.Coordinators = &coords
	}

	return apiReq, nil
}