package v0_0_41

import (
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
)

// convertAPIAccountToCommon converts a v0.0.41 API Account to common Account type
func (a *AccountAdapter) convertAPIAccountToCommon(apiAccount interface{}) (*types.Account, error) {
	// Type assertion to handle the anonymous struct
	accountData, ok := apiAccount.(struct {
		Associations    *[]struct {
			Account    *string `json:"account,omitempty"`
			Cluster    *string `json:"cluster,omitempty"`
			Id         *api.V0041OpenapiAccountsRespAccountsAssociationsId `json:"id,omitempty"`
			IsDefault  *bool   `json:"is_default,omitempty"`
			Max        *struct {
				Jobs        *struct {
					Total *api.V0041OpenapiAccountsRespAccountsAssociationsMaxJobsTotal `json:"total,omitempty"`
				} `json:"jobs,omitempty"`
				PerAccount  *struct {
					CpuMinutes    *api.V0041OpenapiAccountsRespAccountsAssociationsMaxPerAccountCpuMinutes    `json:"cpu_minutes,omitempty"`
					RunMinutes    *api.V0041OpenapiAccountsRespAccountsAssociationsMaxPerAccountRunMinutes    `json:"run_minutes,omitempty"`
					SubmittingJobs *api.V0041OpenapiAccountsRespAccountsAssociationsMaxPerAccountSubmittingJobs `json:"submitting_jobs,omitempty"`
					SuspendedJobs  *api.V0041OpenapiAccountsRespAccountsAssociationsMaxPerAccountSuspendedJobs  `json:"suspended_jobs,omitempty"`
				} `json:"per_account,omitempty"`
				TresGroup   *struct {
					Minutes *string `json:"minutes,omitempty"`
					Total   *string `json:"total,omitempty"`
				} `json:"tres_group,omitempty"`
				TresTotal   *struct {
					Minutes *string `json:"minutes,omitempty"`
					Total   *string `json:"total,omitempty"`
				} `json:"tres_total,omitempty"`
			} `json:"max,omitempty"`
			Min         *struct {
				PriorityThreshold *api.V0041OpenapiAccountsRespAccountsAssociationsMinPriorityThreshold `json:"priority_threshold,omitempty"`
			} `json:"min,omitempty"`
			Parent      *string `json:"parent,omitempty"`
			Partition   *string `json:"partition,omitempty"`
			Priority    *api.V0041OpenapiAccountsRespAccountsAssociationsPriority `json:"priority,omitempty"`
			Qos         *[]string `json:"qos,omitempty"`
			SharesRaw   *int32    `json:"shares_raw,omitempty"`
			User        *string   `json:"user,omitempty"`
		} `json:"associations,omitempty"`
		Coordinators    *[]struct {
			Name *string `json:"name,omitempty"`
		} `json:"coordinators,omitempty"`
		Description     *string   `json:"description,omitempty"`
		Flags           *[]api.V0041OpenapiAccountsRespAccountsFlags `json:"flags,omitempty"`
		Name            *string   `json:"name,omitempty"`
		Organization    *string   `json:"organization,omitempty"`
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

	// Flags
	if accountData.Flags != nil {
		var flags []string
		for _, flag := range *accountData.Flags {
			flags = append(flags, string(flag))
		}
		account.Flags = flags
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

	// Associations (if requested)
	if accountData.Associations != nil {
		associations := make([]types.Association, 0, len(*accountData.Associations))
		for _, apiAssoc := range *accountData.Associations {
			assoc := types.Association{}

			if apiAssoc.Id != nil && apiAssoc.Id.Number != nil {
				assoc.ID = uint32(*apiAssoc.Id.Number)
			}
			if apiAssoc.Account != nil {
				assoc.Account = *apiAssoc.Account
			}
			if apiAssoc.Cluster != nil {
				assoc.Cluster = *apiAssoc.Cluster
			}
			if apiAssoc.User != nil {
				assoc.User = *apiAssoc.User
			}
			if apiAssoc.Partition != nil {
				assoc.Partition = *apiAssoc.Partition
			}
			if apiAssoc.Parent != nil {
				assoc.ParentAccount = *apiAssoc.Parent
			}
			if apiAssoc.IsDefault != nil {
				assoc.IsDefault = *apiAssoc.IsDefault
			}
			if apiAssoc.SharesRaw != nil {
				assoc.SharesRaw = uint32(*apiAssoc.SharesRaw)
			}
			if apiAssoc.Priority != nil && apiAssoc.Priority.Number != nil {
				assoc.Priority = uint32(*apiAssoc.Priority.Number)
			}

			// QoS list
			if apiAssoc.Qos != nil {
				assoc.QosList = *apiAssoc.Qos
			}

			// Limits
			if apiAssoc.Max != nil {
				if apiAssoc.Max.Jobs != nil && apiAssoc.Max.Jobs.Total != nil && apiAssoc.Max.Jobs.Total.Number != nil {
					assoc.MaxJobs = uint32(*apiAssoc.Max.Jobs.Total.Number)
				}
				if apiAssoc.Max.PerAccount != nil {
					if apiAssoc.Max.PerAccount.SubmittingJobs != nil && apiAssoc.Max.PerAccount.SubmittingJobs.Number != nil {
						assoc.MaxSubmitJobs = uint32(*apiAssoc.Max.PerAccount.SubmittingJobs.Number)
					}
				}
				if apiAssoc.Max.TresTotal != nil {
					if apiAssoc.Max.TresTotal.Total != nil {
						assoc.MaxTRES = *apiAssoc.Max.TresTotal.Total
					}
					if apiAssoc.Max.TresTotal.Minutes != nil {
						assoc.MaxTRESMinutes = *apiAssoc.Max.TresTotal.Minutes
					}
				}
			}

			associations = append(associations, assoc)
		}
		account.Associations = associations
	}

	return account, nil
}

// convertCommonToAPIAccount converts common Account to v0.0.41 API request
func (a *AccountAdapter) convertCommonToAPIAccount(account *types.Account) *api.V0041OpenapiAccountsResp {
	req := &api.V0041OpenapiAccountsResp{
		Accounts: []struct {
			Associations    *[]struct {
				Account    *string `json:"account,omitempty"`
				Cluster    *string `json:"cluster,omitempty"`
				Id         *api.V0041OpenapiAccountsRespAccountsAssociationsId `json:"id,omitempty"`
				IsDefault  *bool   `json:"is_default,omitempty"`
				Max        *struct {
					Jobs        *struct {
						Total *api.V0041OpenapiAccountsRespAccountsAssociationsMaxJobsTotal `json:"total,omitempty"`
					} `json:"jobs,omitempty"`
					PerAccount  *struct {
						CpuMinutes    *api.V0041OpenapiAccountsRespAccountsAssociationsMaxPerAccountCpuMinutes    `json:"cpu_minutes,omitempty"`
						RunMinutes    *api.V0041OpenapiAccountsRespAccountsAssociationsMaxPerAccountRunMinutes    `json:"run_minutes,omitempty"`
						SubmittingJobs *api.V0041OpenapiAccountsRespAccountsAssociationsMaxPerAccountSubmittingJobs `json:"submitting_jobs,omitempty"`
						SuspendedJobs  *api.V0041OpenapiAccountsRespAccountsAssociationsMaxPerAccountSuspendedJobs  `json:"suspended_jobs,omitempty"`
					} `json:"per_account,omitempty"`
					TresGroup   *struct {
						Minutes *string `json:"minutes,omitempty"`
						Total   *string `json:"total,omitempty"`
					} `json:"tres_group,omitempty"`
					TresTotal   *struct {
						Minutes *string `json:"minutes,omitempty"`
						Total   *string `json:"total,omitempty"`
					} `json:"tres_total,omitempty"`
				} `json:"max,omitempty"`
				Min         *struct {
					PriorityThreshold *api.V0041OpenapiAccountsRespAccountsAssociationsMinPriorityThreshold `json:"priority_threshold,omitempty"`
				} `json:"min,omitempty"`
				Parent      *string `json:"parent,omitempty"`
				Partition   *string `json:"partition,omitempty"`
				Priority    *api.V0041OpenapiAccountsRespAccountsAssociationsPriority `json:"priority,omitempty"`
				Qos         *[]string `json:"qos,omitempty"`
				SharesRaw   *int32    `json:"shares_raw,omitempty"`
				User        *string   `json:"user,omitempty"`
			} `json:"associations,omitempty"`
			Coordinators    *[]struct {
				Name *string `json:"name,omitempty"`
			} `json:"coordinators,omitempty"`
			Description     *string   `json:"description,omitempty"`
			Flags           *[]api.V0041OpenapiAccountsRespAccountsFlags `json:"flags,omitempty"`
			Name            *string   `json:"name,omitempty"`
			Organization    *string   `json:"organization,omitempty"`
		}{
			{},
		},
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

	// Convert flags
	if len(account.Flags) > 0 {
		flags := make([]api.V0041OpenapiAccountsRespAccountsFlags, 0, len(account.Flags))
		for _, flag := range account.Flags {
			// Map common flags to API flags
			switch strings.ToLower(flag) {
			case "deleted":
				flags = append(flags, api.V0041OpenapiAccountsRespAccountsFlagsDELETED)
			case "nousersarecoords":
				flags = append(flags, api.V0041OpenapiAccountsRespAccountsFlagsNoUsersAreCoords)
			case "usersarecoords":
				flags = append(flags, api.V0041OpenapiAccountsRespAccountsFlagsUsersAreCoords)
			case "withassociations":
				flags = append(flags, api.V0041OpenapiAccountsRespAccountsFlagsWithAssociations)
			case "withcoordinators":
				flags = append(flags, api.V0041OpenapiAccountsRespAccountsFlagsWithCoordinators)
			}
		}
		if len(flags) > 0 {
			acc.Flags = &flags
		}
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
func (a *AccountAdapter) convertCommonToAPIAccountUpdate(update *types.AccountUpdate) *api.V0041OpenapiAccountsResp {
	// For v0.0.41, updates are done by sending the full account object
	// This is a placeholder that would need the full account data
	return nil
}