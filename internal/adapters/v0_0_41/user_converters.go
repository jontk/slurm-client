package v0_0_41

import (
	"fmt"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
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

	// Default account and wckey
	if v, ok := userData["default"]; ok {
		if defaultData, ok := v.(map[string]interface{}); ok {
			if acc, ok := defaultData["account"]; ok {
				if account, ok := acc.(string); ok {
					user.DefaultAccount = account
				}
			}
			if wc, ok := defaultData["wckey"]; ok {
				if wckey, ok := wc.(string); ok {
					user.DefaultWCKey = wckey
				}
			}
		}
	}

	// Administrator level
	if v, ok := userData["administrator_level"]; ok {
		if levels, ok := v.([]interface{}); ok && len(levels) > 0 {
			if level, ok := levels[0].(string); ok {
				user.AdminLevel = types.AdminLevel(level)
			}
		}
	}

	// Associations
	if v, ok := userData["associations"]; ok {
		if assocData, ok := v.([]interface{}); ok {
			userAssocs := make([]types.UserAssociation, 0, len(assocData))
			for _, a := range assocData {
				if assoc, ok := a.(map[string]interface{}); ok {
					userAssoc := types.UserAssociation{}
					if acc, ok := assoc["account"].(string); ok {
						userAssoc.AccountName = acc
					}
					if cl, ok := assoc["cluster"].(string); ok {
						userAssoc.Cluster = cl
					}
					if part, ok := assoc["partition"].(string); ok {
						userAssoc.Partition = part
					}
					userAssocs = append(userAssocs, userAssoc)
				}
			}
			user.Associations = userAssocs
		}
	}

	// Coordinators
	if v, ok := userData["coordinators"]; ok {
		if coordData, ok := v.([]interface{}); ok {
			for _, c := range coordData {
				if coord, ok := c.(map[string]interface{}); ok {
					if acc, ok := coord["account"].(string); ok {
						user.Coordinators = append(user.Coordinators, types.UserCoordinator{
							AccountName: acc,
						})
					}
				}
			}
		}
	}

	// WCKeys
	if v, ok := userData["wckeys"]; ok {
		if wckeyData, ok := v.([]interface{}); ok {
			wckeys := make([]string, 0, len(wckeyData))
			for _, w := range wckeyData {
				if wckey, ok := w.(map[string]interface{}); ok {
					if name, ok := wckey["name"].(string); ok {
						wckeys = append(wckeys, name)
					}
				}
			}
			user.WCKeys = wckeys
		}
	}

	// Skip complex nested structures for now
	// These would need detailed mapping based on the actual API response structure

	return user, nil
}

// convertCommonToAPIUser converts common User to API format
func (a *UserAdapter) convertCommonToAPIUser(user *types.User) *api.V0041OpenapiUsersResp {
	apiUser := &api.V0041OpenapiUsersResp{
		Users: []struct {
			// AdministratorLevel AdminLevel granted to the user
			AdministratorLevel *[]api.V0041OpenapiUsersRespUsersAdministratorLevel `json:"administrator_level,omitempty"`

			// Associations Associations created for this user
			Associations *[]struct {
				// Account Account
				Account *string `json:"account,omitempty"`

				// Cluster Cluster
				Cluster *string `json:"cluster,omitempty"`

				// Id Numeric association ID
				Id *int32 `json:"id,omitempty"`

				// Partition Partition
				Partition *string `json:"partition,omitempty"`

				// User User name
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
				// Accounting Accounting records containing related resource usage
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

				// Cluster Cluster name
				Cluster string `json:"cluster"`

				// Flags Flags associated with the WCKey
				Flags *[]api.V0041OpenapiUsersRespUsersWckeysFlags `json:"flags,omitempty"`

				// Id Unique ID for this user-cluster-wckey combination
				Id *int32 `json:"id,omitempty"`

				// Name WCKey name
				Name string `json:"name"`

				// User User name
				User string `json:"user"`
			} `json:"wckeys,omitempty"`
		}{
			{
				Name: user.Name,
			},
		},
	}

	// Set admin level
	if user.AdminLevel != "" {
		level := api.V0041OpenapiUsersRespUsersAdministratorLevel(user.AdminLevel)
		apiUser.Users[0].AdministratorLevel = &[]api.V0041OpenapiUsersRespUsersAdministratorLevel{level}
	}

	// Set default account and wckey
	if user.DefaultAccount != "" || user.DefaultWCKey != "" {
		apiUser.Users[0].Default = &struct {
			// Account Default bank account name
			Account *string `json:"account,omitempty"`

			// Wckey Default wckey name
			Wckey *string `json:"wckey,omitempty"`
		}{}
		if user.DefaultAccount != "" {
			apiUser.Users[0].Default.Account = &user.DefaultAccount
		}
		if user.DefaultWCKey != "" {
			apiUser.Users[0].Default.Wckey = &user.DefaultWCKey
		}
	}

	return apiUser
}