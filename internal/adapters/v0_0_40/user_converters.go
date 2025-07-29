package v0_0_40

import (
	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
)

// convertAPIUserToCommon converts a v0.0.40 API User to common User type
func (a *UserAdapter) convertAPIUserToCommon(apiUser api.V0040User) (*types.User, error) {
	user := &types.User{}

	// Basic fields
	if apiUser.Name != nil {
		user.Name = *apiUser.Name
	}
	if apiUser.Uid != nil {
		user.UID = *apiUser.Uid
	}
	if apiUser.Default != nil && apiUser.Default.Account != nil {
		user.DefaultAccount = *apiUser.Default.Account
	}
	if apiUser.Default != nil && apiUser.Default.Wckey != nil {
		user.DefaultWCKey = *apiUser.Default.Wckey
	}

	// Admin level
	if apiUser.Administrator != nil && len(*apiUser.Administrator) > 0 {
		// Convert admin level array to single admin level
		user.AdminLevel = string((*apiUser.Administrator)[0])
	}

	// Flags
	if apiUser.Flags != nil && len(*apiUser.Flags) > 0 {
		user.Flags = make([]string, len(*apiUser.Flags))
		for i, flag := range *apiUser.Flags {
			user.Flags[i] = string(flag)
		}
	}

	// Coordinators
	if apiUser.Coordinators != nil {
		user.Coordinators = make([]types.CoordinatorInfo, 0, len(*apiUser.Coordinators))
		for _, apiCoord := range *apiUser.Coordinators {
			coord := types.CoordinatorInfo{}
			if apiCoord.Name != nil {
				coord.Name = *apiCoord.Name
			}
			if apiCoord.Direct != nil && apiCoord.Direct.Count != nil {
				coord.DirectCount = int32(*apiCoord.Direct.Count)
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
				assoc.Account = *apiAssoc.Account
			}
			if apiAssoc.Cluster != nil {
				assoc.Cluster = *apiAssoc.Cluster
			}
			if apiAssoc.Partition != nil {
				assoc.Partition = *apiAssoc.Partition
			}
			if apiAssoc.User != nil {
				assoc.User = *apiAssoc.User
			}
			if apiAssoc.Defaults != nil && apiAssoc.Defaults.Qos != nil {
				assoc.DefaultQOS = *apiAssoc.Defaults.Qos
			}
			if apiAssoc.Flags != nil && len(*apiAssoc.Flags) > 0 {
				assoc.Flags = make([]string, len(*apiAssoc.Flags))
				for i, flag := range *apiAssoc.Flags {
					assoc.Flags[i] = string(flag)
				}
			}
			if apiAssoc.Shares != nil && apiAssoc.Shares.Object != nil && apiAssoc.Shares.Object.Raw != nil {
				shareVal := int32(*apiAssoc.Shares.Object.Raw)
				assoc.RawShares = &shareVal
			}
			if apiAssoc.Priority != nil && apiAssoc.Priority.Number != nil {
				priority := *apiAssoc.Priority.Number
				assoc.Priority = &priority
			}
			user.Associations = append(user.Associations, assoc)
		}
	}

	return user, nil
}

// convertCommonUserCreateToAPI converts common UserCreate to v0.0.40 API format
func (a *UserAdapter) convertCommonUserCreateToAPI(user *types.UserCreate) (*api.V0040User, error) {
	apiUser := &api.V0040User{}

	// Basic fields
	apiUser.Name = &user.Name
	
	if user.UID > 0 {
		apiUser.Uid = &user.UID
	}

	// Default settings
	defaults := &api.V0040UserDefault{}
	if user.DefaultAccount != "" {
		defaults.Account = &user.DefaultAccount
	}
	if user.DefaultWCKey != "" {
		defaults.Wckey = &user.DefaultWCKey
	}
	if user.DefaultAccount != "" || user.DefaultWCKey != "" {
		apiUser.Default = defaults
	}

	// Admin level
	if user.AdminLevel != "" {
		adminLevels := []api.V0040UserAdministrator{api.V0040UserAdministrator(user.AdminLevel)}
		apiUser.Administrator = &adminLevels
	}

	// Flags
	if len(user.Flags) > 0 {
		flags := make([]api.V0040UserFlags, len(user.Flags))
		for i, flag := range user.Flags {
			flags[i] = api.V0040UserFlags(flag)
		}
		apiUser.Flags = &flags
	}

	// Coordinators
	if len(user.Coordinators) > 0 {
		coords := make([]api.V0040Coordinator, len(user.Coordinators))
		for i, coord := range user.Coordinators {
			apiCoord := api.V0040Coordinator{
				Name: &coord.Name,
			}
			if coord.DirectCount > 0 {
				count := int(coord.DirectCount)
				apiCoord.Direct = &api.V0040CoordinatorDirect{
					Count: &count,
				}
			}
			coords[i] = apiCoord
		}
		apiUser.Coordinators = &coords
	}

	return apiUser, nil
}

// convertCommonUserUpdateToAPI converts common UserUpdate to v0.0.40 API format
func (a *UserAdapter) convertCommonUserUpdateToAPI(existingUser *types.User, update *types.UserUpdate) (*api.V0040User, error) {
	apiUser := &api.V0040User{}

	// Name (required)
	apiUser.Name = &existingUser.Name

	// Apply updates
	defaults := &api.V0040UserDefault{}
	hasDefaults := false

	if update.DefaultAccount != nil {
		defaults.Account = update.DefaultAccount
		hasDefaults = true
	}
	if update.DefaultWCKey != nil {
		defaults.Wckey = update.DefaultWCKey
		hasDefaults = true
	}

	if hasDefaults {
		apiUser.Default = defaults
	}

	// Admin level
	if update.AdminLevel != nil {
		adminLevels := []api.V0040UserAdministrator{api.V0040UserAdministrator(*update.AdminLevel)}
		apiUser.Administrator = &adminLevels
	}

	// Flags
	if update.Flags != nil {
		flags := make([]api.V0040UserFlags, len(*update.Flags))
		for i, flag := range *update.Flags {
			flags[i] = api.V0040UserFlags(flag)
		}
		apiUser.Flags = &flags
	}

	return apiUser, nil
}