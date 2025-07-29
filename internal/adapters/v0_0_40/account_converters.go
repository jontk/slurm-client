package v0_0_40

import (
	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
)

// convertAPIAccountToCommon converts a v0.0.40 API Account to common Account type
func (a *AccountAdapter) convertAPIAccountToCommon(apiAccount api.V0040Account) (*types.Account, error) {
	account := &types.Account{}

	// Basic fields
	if apiAccount.Name != nil {
		account.Name = *apiAccount.Name
	}
	if apiAccount.Description != nil {
		account.Description = *apiAccount.Description
	}
	if apiAccount.Organization != nil {
		account.Organization = *apiAccount.Organization
	}

	// Flags
	if apiAccount.Flags != nil && len(*apiAccount.Flags) > 0 {
		account.Flags = make([]string, len(*apiAccount.Flags))
		for i, flag := range *apiAccount.Flags {
			account.Flags[i] = string(flag)
		}
	}

	// Coordinators
	if apiAccount.Coordinators != nil {
		account.Coordinators = make([]types.Coordinator, 0, len(*apiAccount.Coordinators))
		for _, apiCoord := range *apiAccount.Coordinators {
			coord := types.Coordinator{}
			if apiCoord.Name != nil {
				coord.Name = *apiCoord.Name
			}
			if apiCoord.Direct != nil && apiCoord.Direct.Count != nil {
				coord.DirectCount = int32(*apiCoord.Direct.Count)
			}
			account.Coordinators = append(account.Coordinators, coord)
		}
	}

	// Associations
	if apiAccount.Associations != nil {
		account.Associations = make([]types.AssociationShort, 0, len(*apiAccount.Associations))
		for _, apiAssoc := range *apiAccount.Associations {
			assoc := types.AssociationShort{}
			if apiAssoc.User != nil {
				assoc.User = *apiAssoc.User
			}
			if apiAssoc.Cluster != nil {
				assoc.Cluster = *apiAssoc.Cluster
			}
			if apiAssoc.Partition != nil {
				assoc.Partition = *apiAssoc.Partition
			}
			if apiAssoc.Account != nil {
				assoc.Account = *apiAssoc.Account
			}
			account.Associations = append(account.Associations, assoc)
		}
	}

	return account, nil
}

// convertCommonAccountCreateToAPI converts common AccountCreate to v0.0.40 API format
func (a *AccountAdapter) convertCommonAccountCreateToAPI(account *types.AccountCreate) (*api.V0040Account, error) {
	apiAccount := &api.V0040Account{}

	// Basic fields
	apiAccount.Name = &account.Name
	
	if account.Description != "" {
		apiAccount.Description = &account.Description
	}
	if account.Organization != "" {
		apiAccount.Organization = &account.Organization
	}

	// Parent account
	if account.Parent != "" {
		// v0.0.40 may handle parent differently, might need to set via associations
		// For now, we'll skip this as it's typically handled by the accounting system
	}

	// Flags
	if len(account.Flags) > 0 {
		flags := make([]api.V0040AccountFlags, len(account.Flags))
		for i, flag := range account.Flags {
			flags[i] = api.V0040AccountFlags(flag)
		}
		apiAccount.Flags = &flags
	}

	// Coordinators
	if len(account.Coordinators) > 0 {
		coords := make([]api.V0040Coordinator, len(account.Coordinators))
		for i, coord := range account.Coordinators {
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
		apiAccount.Coordinators = &coords
	}

	return apiAccount, nil
}

// convertCommonAccountUpdateToAPI converts common AccountUpdate to v0.0.40 API format
func (a *AccountAdapter) convertCommonAccountUpdateToAPI(existingAccount *types.Account, update *types.AccountUpdate) (*api.V0040Account, error) {
	apiAccount := &api.V0040Account{}

	// Name (required)
	apiAccount.Name = &existingAccount.Name

	// Apply updates
	if update.Description != nil {
		apiAccount.Description = update.Description
	}
	if update.Organization != nil {
		apiAccount.Organization = update.Organization
	}

	// Flags
	if update.Flags != nil {
		flags := make([]api.V0040AccountFlags, len(*update.Flags))
		for i, flag := range *update.Flags {
			flags[i] = api.V0040AccountFlags(flag)
		}
		apiAccount.Flags = &flags
	}

	return apiAccount, nil
}