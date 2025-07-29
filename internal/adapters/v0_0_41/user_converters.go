package v0_0_41

import (
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
)

// convertAPIUserToCommon converts a v0.0.41 API User to common User type
func (a *UserAdapter) convertAPIUserToCommon(apiUser interface{}) (*types.User, error) {
	// Type assertion to handle the anonymous struct
	userData, ok := apiUser.(struct {
		AdministratorLevel *[]api.V0041OpenapiUsersRespUsersAdministratorLevel `json:"administrator_level,omitempty"`
		Associations       *[]struct {
			Account    *string `json:"account,omitempty"`
			Cluster    *string `json:"cluster,omitempty"`
			Partition  *string `json:"partition,omitempty"`
		} `json:"associations,omitempty"`
		Coordinators       *[]struct {
			Name *string `json:"name,omitempty"`
		} `json:"coordinators,omitempty"`
		Default            *struct {
			Account *string `json:"account,omitempty"`
			Wckey   *string `json:"wckey,omitempty"`
		} `json:"default,omitempty"`
		Flags              *[]api.V0041OpenapiUsersRespUsersFlags `json:"flags,omitempty"`
		Name               *string `json:"name,omitempty"`
		OldName            *string `json:"old_name,omitempty"`
		Wckeys             *[]struct {
			Cluster *string `json:"cluster,omitempty"`
			Flags   *[]api.V0041OpenapiUsersRespUsersWckeysFlags `json:"flags,omitempty"`
			Id      *api.V0041OpenapiUsersRespUsersWckeysId `json:"id,omitempty"`
			Name    *string `json:"name,omitempty"`
			User    *string `json:"user,omitempty"`
		} `json:"wckeys,omitempty"`
	})
	if !ok {
		return nil, fmt.Errorf("unexpected user data type")
	}

	user := &types.User{}

	// Basic fields
	if userData.Name != nil {
		user.Name = *userData.Name
	}
	if userData.OldName != nil {
		user.OldName = *userData.OldName
	}

	// Default account and wckey
	if userData.Default != nil {
		if userData.Default.Account != nil {
			user.DefaultAccount = *userData.Default.Account
		}
		if userData.Default.Wckey != nil {
			user.DefaultWCKey = *userData.Default.Wckey
		}
	}

	// Administrator level
	if userData.AdministratorLevel != nil && len(*userData.AdministratorLevel) > 0 {
		// Convert the first admin level to string
		adminLevel := string((*userData.AdministratorLevel)[0])
		user.AdminLevel = adminLevel
	}

	// Flags
	if userData.Flags != nil {
		var flags []string
		for _, flag := range *userData.Flags {
			flags = append(flags, string(flag))
		}
		user.Flags = flags
	}

	// Associations (if requested)
	if userData.Associations != nil {
		associations := make([]types.Association, 0, len(*userData.Associations))
		for _, apiAssoc := range *userData.Associations {
			assoc := types.Association{
				User: user.Name,
			}

			if apiAssoc.Account != nil {
				assoc.Account = *apiAssoc.Account
			}
			if apiAssoc.Cluster != nil {
				assoc.Cluster = *apiAssoc.Cluster
			}
			if apiAssoc.Partition != nil {
				assoc.Partition = *apiAssoc.Partition
			}

			associations = append(associations, assoc)
		}
		user.Associations = associations
	}

	// Coordinators (accounts this user coordinates)
	if userData.Coordinators != nil {
		coordinators := make([]string, 0, len(*userData.Coordinators))
		for _, coord := range *userData.Coordinators {
			if coord.Name != nil {
				coordinators = append(coordinators, *coord.Name)
			}
		}
		user.CoordinatorAccounts = coordinators
	}

	// WCKeys
	if userData.Wckeys != nil {
		wckeys := make([]types.WCKey, 0, len(*userData.Wckeys))
		for _, apiWckey := range *userData.Wckeys {
			wckey := types.WCKey{}

			if apiWckey.Id != nil && apiWckey.Id.Number != nil {
				wckey.ID = uint32(*apiWckey.Id.Number)
			}
			if apiWckey.Name != nil {
				wckey.Name = *apiWckey.Name
			}
			if apiWckey.User != nil {
				wckey.User = *apiWckey.User
			}
			if apiWckey.Cluster != nil {
				wckey.Cluster = *apiWckey.Cluster
			}

			// Flags
			if apiWckey.Flags != nil {
				var flags []string
				for _, flag := range *apiWckey.Flags {
					flags = append(flags, string(flag))
				}
				wckey.Flags = flags
			}

			wckeys = append(wckeys, wckey)
		}
		user.WCKeys = wckeys
	}

	return user, nil
}

// convertCommonToAPIUser converts common User to v0.0.41 API request
func (a *UserAdapter) convertCommonToAPIUser(user *types.User) *api.V0041OpenapiUsersResp {
	req := &api.V0041OpenapiUsersResp{
		Users: []struct {
			AdministratorLevel *[]api.V0041OpenapiUsersRespUsersAdministratorLevel `json:"administrator_level,omitempty"`
			Associations       *[]struct {
				Account    *string `json:"account,omitempty"`
				Cluster    *string `json:"cluster,omitempty"`
				Partition  *string `json:"partition,omitempty"`
			} `json:"associations,omitempty"`
			Coordinators       *[]struct {
				Name *string `json:"name,omitempty"`
			} `json:"coordinators,omitempty"`
			Default            *struct {
				Account *string `json:"account,omitempty"`
				Wckey   *string `json:"wckey,omitempty"`
			} `json:"default,omitempty"`
			Flags              *[]api.V0041OpenapiUsersRespUsersFlags `json:"flags,omitempty"`
			Name               *string `json:"name,omitempty"`
			OldName            *string `json:"old_name,omitempty"`
			Wckeys             *[]struct {
				Cluster *string `json:"cluster,omitempty"`
				Flags   *[]api.V0041OpenapiUsersRespUsersWckeysFlags `json:"flags,omitempty"`
				Id      *api.V0041OpenapiUsersRespUsersWckeysId `json:"id,omitempty"`
				Name    *string `json:"name,omitempty"`
				User    *string `json:"user,omitempty"`
			} `json:"wckeys,omitempty"`
		}{
			{},
		},
	}

	usr := &req.Users[0]

	// Set basic fields
	if user.Name != "" {
		usr.Name = &user.Name
	}
	if user.OldName != "" {
		usr.OldName = &user.OldName
	}

	// Set default account and wckey
	if user.DefaultAccount != "" || user.DefaultWCKey != "" {
		usr.Default = &struct {
			Account *string `json:"account,omitempty"`
			Wckey   *string `json:"wckey,omitempty"`
		}{}
		if user.DefaultAccount != "" {
			usr.Default.Account = &user.DefaultAccount
		}
		if user.DefaultWCKey != "" {
			usr.Default.Wckey = &user.DefaultWCKey
		}
	}

	// Set administrator level
	if user.AdminLevel != "" {
		adminLevels := []api.V0041OpenapiUsersRespUsersAdministratorLevel{}
		switch strings.ToLower(user.AdminLevel) {
		case "administrator":
			adminLevels = append(adminLevels, api.V0041OpenapiUsersRespUsersAdministratorLevelAdministrator)
		case "operator":
			adminLevels = append(adminLevels, api.V0041OpenapiUsersRespUsersAdministratorLevelOperator)
		case "none":
			adminLevels = append(adminLevels, api.V0041OpenapiUsersRespUsersAdministratorLevelNone)
		}
		if len(adminLevels) > 0 {
			usr.AdministratorLevel = &adminLevels
		}
	}

	// Convert flags
	if len(user.Flags) > 0 {
		flags := make([]api.V0041OpenapiUsersRespUsersFlags, 0, len(user.Flags))
		for _, flag := range user.Flags {
			// Map common flags to API flags
			switch strings.ToLower(flag) {
			case "deleted":
				flags = append(flags, api.V0041OpenapiUsersRespUsersFlagsDELETED)
			case "exact":
				flags = append(flags, api.V0041OpenapiUsersRespUsersFlagsExact)
			case "noupdate":
				flags = append(flags, api.V0041OpenapiUsersRespUsersFlagsNoUpdate)
			case "nousersarecoords":
				flags = append(flags, api.V0041OpenapiUsersRespUsersFlagsNoUsersAreCoords)
			case "usersarecoords":
				flags = append(flags, api.V0041OpenapiUsersRespUsersFlagsUsersAreCoords)
			}
		}
		if len(flags) > 0 {
			usr.Flags = &flags
		}
	}

	// Convert coordinators
	if len(user.CoordinatorAccounts) > 0 {
		coords := make([]struct {
			Name *string `json:"name,omitempty"`
		}, 0, len(user.CoordinatorAccounts))
		for _, coordName := range user.CoordinatorAccounts {
			coordNameCopy := coordName // Create a copy to avoid pointer issues
			coords = append(coords, struct {
				Name *string `json:"name,omitempty"`
			}{
				Name: &coordNameCopy,
			})
		}
		usr.Coordinators = &coords
	}

	// Convert WCKeys
	if len(user.WCKeys) > 0 {
		wckeys := make([]struct {
			Cluster *string `json:"cluster,omitempty"`
			Flags   *[]api.V0041OpenapiUsersRespUsersWckeysFlags `json:"flags,omitempty"`
			Id      *api.V0041OpenapiUsersRespUsersWckeysId `json:"id,omitempty"`
			Name    *string `json:"name,omitempty"`
			User    *string `json:"user,omitempty"`
		}, 0, len(user.WCKeys))

		for _, wckey := range user.WCKeys {
			apiWckey := struct {
				Cluster *string `json:"cluster,omitempty"`
				Flags   *[]api.V0041OpenapiUsersRespUsersWckeysFlags `json:"flags,omitempty"`
				Id      *api.V0041OpenapiUsersRespUsersWckeysId `json:"id,omitempty"`
				Name    *string `json:"name,omitempty"`
				User    *string `json:"user,omitempty"`
			}{}

			if wckey.Name != "" {
				apiWckey.Name = &wckey.Name
			}
			if wckey.User != "" {
				apiWckey.User = &wckey.User
			}
			if wckey.Cluster != "" {
				apiWckey.Cluster = &wckey.Cluster
			}

			// Convert wckey flags
			if len(wckey.Flags) > 0 {
				wckeyFlags := make([]api.V0041OpenapiUsersRespUsersWckeysFlags, 0, len(wckey.Flags))
				for _, flag := range wckey.Flags {
					switch strings.ToLower(flag) {
					case "deleted":
						wckeyFlags = append(wckeyFlags, api.V0041OpenapiUsersRespUsersWckeysFlagsDELETED)
					case "exact":
						wckeyFlags = append(wckeyFlags, api.V0041OpenapiUsersRespUsersWckeysFlagsExact)
					case "noupdate":
						wckeyFlags = append(wckeyFlags, api.V0041OpenapiUsersRespUsersWckeysFlagsNoUpdate)
					}
				}
				if len(wckeyFlags) > 0 {
					apiWckey.Flags = &wckeyFlags
				}
			}

			wckeys = append(wckeys, apiWckey)
		}
		usr.Wckeys = &wckeys
	}

	return req
}

// convertCommonToAPIUserUpdate converts common UserUpdate to v0.0.41 API request
func (a *UserAdapter) convertCommonToAPIUserUpdate(update *types.UserUpdate) *api.V0041OpenapiUsersResp {
	// For v0.0.41, updates are done by sending the full user object
	// This is a placeholder that would need the full user data
	return nil
}