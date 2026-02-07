// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"fmt"

	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
)

// convertAPIUserToCommon converts a v0.0.41 API User to common User type
func (a *UserAdapter) convertAPIUserToCommon(apiUser interface{}) (*types.User, error) {
	// Use map interface for handling anonymous structs in v0.0.41
	userData, ok := apiUser.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected user data type: %T", apiUser)
	}
	user := &types.User{}
	// Basic fields - using safe type assertions
	if v, ok := userData["name"]; ok {
		if name, ok := v.(string); ok {
			user.Name = name
		}
	}
	// Default account and wckey (nested in Default struct)
	if v, ok := userData["default"]; ok {
		if defaultData, ok := v.(map[string]interface{}); ok {
			user.Default = &types.UserDefault{}
			if acc, ok := defaultData["account"]; ok {
				if account, ok := acc.(string); ok {
					user.Default.Account = &account
				}
			}
			if wc, ok := defaultData["wckey"]; ok {
				if wckey, ok := wc.(string); ok {
					user.Default.Wckey = &wckey
				}
			}
		}
	}
	// Administrator level (slice of AdministratorLevelValue)
	if v, ok := userData["administrator_level"]; ok {
		if levels, ok := v.([]interface{}); ok {
			adminLevels := make([]types.AdministratorLevelValue, 0, len(levels))
			for _, l := range levels {
				if level, ok := l.(string); ok {
					adminLevels = append(adminLevels, types.AdministratorLevelValue(level))
				}
			}
			user.AdministratorLevel = adminLevels
		}
	}
	// Associations (AssocShort slice)
	if v, ok := userData["associations"]; ok {
		if assocData, ok := v.([]interface{}); ok {
			userAssocs := make([]types.AssocShort, 0, len(assocData))
			for _, a := range assocData {
				assoc, ok := a.(map[string]interface{})
				if !ok {
					continue
				}
				userAssoc := types.AssocShort{}
				if acc, ok := assoc["account"].(string); ok {
					userAssoc.Account = &acc
				}
				if cl, ok := assoc["cluster"].(string); ok {
					userAssoc.Cluster = &cl
				}
				if part, ok := assoc["partition"].(string); ok {
					userAssoc.Partition = &part
				}
				if usr, ok := assoc["user"].(string); ok {
					userAssoc.User = usr
				}
				userAssocs = append(userAssocs, userAssoc)
			}
			user.Associations = userAssocs
		}
	}
	// Coordinators (Coord slice)
	if v, ok := userData["coordinators"]; ok {
		if coordData, ok := v.([]interface{}); ok {
			coords := make([]types.Coord, 0, len(coordData))
			for _, c := range coordData {
				if coord, ok := c.(map[string]interface{}); ok {
					coordEntry := types.Coord{}
					if name, ok := coord["name"].(string); ok {
						coordEntry.Name = name
					}
					if direct, ok := coord["direct"].(bool); ok {
						coordEntry.Direct = &direct
					}
					coords = append(coords, coordEntry)
				}
			}
			user.Coordinators = coords
		}
	}
	// WCKeys (WCKey slice)
	if v, ok := userData["wckeys"]; ok {
		if wckeyData, ok := v.([]interface{}); ok {
			wckeys := make([]types.WCKey, 0, len(wckeyData))
			for _, w := range wckeyData {
				if wckey, ok := w.(map[string]interface{}); ok {
					wckeyEntry := types.WCKey{}
					if name, ok := wckey["name"].(string); ok {
						wckeyEntry.Name = name
					}
					if cluster, ok := wckey["cluster"].(string); ok {
						wckeyEntry.Cluster = cluster
					}
					if userStr, ok := wckey["user"].(string); ok {
						wckeyEntry.User = userStr
					}
					if id, ok := wckey["id"].(float64); ok {
						idVal := int32(id)
						wckeyEntry.ID = &idVal
					}
					wckeys = append(wckeys, wckeyEntry)
				}
			}
			user.Wckeys = wckeys
		}
	}
	// Flags - convert using helper function pattern
	if v, ok := userData["flags"]; ok {
		if flags, ok := v.([]interface{}); ok {
			userFlags := make([]types.UserDefaultFlagsValue, 0, len(flags))
			for _, f := range flags {
				if flag, ok := f.(string); ok {
					userFlags = append(userFlags, types.UserDefaultFlagsValue(flag))
				}
			}
			user.Flags = userFlags
		}
	}
	return user, nil
}

// convertCommonToAPIUser converts common User to API format
func (a *UserAdapter) convertCommonToAPIUser(user *types.User) *api.V0041OpenapiUsersResp {
	apiUser := &api.V0041OpenapiUsersResp{
		Users: []struct {
			// AdministratorLevel AdminLevel granted to the user
			AdministratorLevel *[]api.V0041OpenapiUsersRespUsersAdministratorLevel `json:"administrator_level,omitempty"`
			// Associations created for this user
			Associations *[]struct {
				// Account details
				Account *string `json:"account,omitempty"`
				// Cluster details
				Cluster *string `json:"cluster,omitempty"`
				// Id Numeric association ID
				Id *int32 `json:"id,omitempty"`
				// Partition details
				Partition *string `json:"partition,omitempty"`
				// User name
				User string `json:"user"`
			} `json:"associations,omitempty"`
			// Coordinators Accounts this user is a coordinator for
			Coordinators *[]struct {
				// Direct Indicates whether the coordinator was directly assigned to this account
				Direct *bool `json:"direct,omitempty"`
				// Name User name
				Name string `json:"name"`
			} `json:"coordinators,omitempty"`
			Default *struct {
				// Account Default bank account name
				Account *string `json:"account,omitempty"`
				// Wckey Default wckey name
				Wckey *string `json:"wckey,omitempty"`
			} `json:"default,omitempty"`
			// Flags flags for a user
			Flags *[]api.V0041OpenapiUsersRespUsersFlags `json:"flags,omitempty"`
			// Name User name
			Name string `json:"name"`
			// OldName Previous user name
			OldName *string `json:"old_name,omitempty"`
			// Wckeys List of available WCKeys
			Wckeys *[]struct {
				// Accounting records containing related resource usage
				Accounting *[]struct {
					// TRES Trackable resources
					TRES *struct {
						// Count TRES count (0 if listed generically)
						Count *int64 `json:"count,omitempty"`
						// Id ID used in database
						Id *int32 `json:"id,omitempty"`
						// Name TRES name (if applicable)
						Name *string `json:"name,omitempty"`
						// Type TRES type (CPU, MEM, etc)
						Type string `json:"type"`
					} `json:"TRES,omitempty"`
					Allocated *struct {
						// Seconds Number of cpu seconds allocated
						Seconds *int64 `json:"seconds,omitempty"`
					} `json:"allocated,omitempty"`
					// Id Association ID or Workload characterization key ID
					Id *int32 `json:"id,omitempty"`
					// Start When the record was started
					Start *int64 `json:"start,omitempty"`
				} `json:"accounting,omitempty"`
				// Cluster name
				Cluster string `json:"cluster"`
				// Flags associated with the WCKey
				Flags *[]api.V0041OpenapiUsersRespUsersWckeysFlags `json:"flags,omitempty"`
				// Id Unique ID for this user-cluster-wckey combination
				Id *int32 `json:"id,omitempty"`
				// Name WCKey name
				Name string `json:"name"`
				// User name
				User string `json:"user"`
			} `json:"wckeys,omitempty"`
		}{
			{
				Name: user.Name,
			},
		},
	}
	// Set admin level
	if len(user.AdministratorLevel) > 0 {
		levels := make([]api.V0041OpenapiUsersRespUsersAdministratorLevel, 0, len(user.AdministratorLevel))
		for _, l := range user.AdministratorLevel {
			levels = append(levels, api.V0041OpenapiUsersRespUsersAdministratorLevel(l))
		}
		apiUser.Users[0].AdministratorLevel = &levels
	}
	// Set default account and wckey
	if user.Default != nil && (user.Default.Account != nil || user.Default.Wckey != nil) {
		apiUser.Users[0].Default = &struct {
			// Account Default bank account name
			Account *string `json:"account,omitempty"`
			// Wckey Default wckey name
			Wckey *string `json:"wckey,omitempty"`
		}{}
		if user.Default.Account != nil {
			apiUser.Users[0].Default.Account = user.Default.Account
		}
		if user.Default.Wckey != nil {
			apiUser.Users[0].Default.Wckey = user.Default.Wckey
		}
	}
	return apiUser
}
