// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

// convertAPIAccountToCommon converts a v0.0.42 API Account to common Account type
func (a *AccountAdapter) convertAPIAccountToCommon(apiAccount api.V0042Account) (*types.Account, error) {
	account := &types.Account{
		Name:         apiAccount.Name,
		Description:  apiAccount.Description,
		Organization: apiAccount.Organization,
	}

	// Convert flags if present
	// Note: Account type doesn't have Flags field, skipping conversion

	// Convert coordinators if present
	if apiAccount.Coordinators != nil && len(*apiAccount.Coordinators) > 0 {
		coordinators := make([]string, 0, len(*apiAccount.Coordinators))
		for _, coord := range *apiAccount.Coordinators {
			// coord.Name is a string, not a pointer
			coordinators = append(coordinators, coord.Name)
		}
		account.Coordinators = coordinators
	}

	// Handle associations if present
	if apiAccount.Associations != nil && len(*apiAccount.Associations) > 0 {
		// Note: Account type doesn't have AssociationCount field, skipping conversion
		// Just extract some basic association info if available
		// This is a simplified conversion since the actual association fields 
		// would require detailed mapping not available in the current account structure
	}

	return account, nil
}

// convertCommonAccountCreateToAPI converts common account create request to v0.0.42 API format
func (a *AccountAdapter) convertCommonAccountCreateToAPI(req *types.AccountCreate) (*api.SlurmdbV0042PostAccountsJSONRequestBody, error) {
	apiReq := &api.SlurmdbV0042PostAccountsJSONRequestBody{
		Accounts: api.V0042AccountList{
			{
				Name:         req.Name,
				Description:  req.Description,
				Organization: req.Organization,
			},
		},
	}

	account := &apiReq.Accounts[0]

	// Add coordinators if specified
	if len(req.Coordinators) > 0 {
		coords := make(api.V0042CoordList, len(req.Coordinators))
		for i, coordName := range req.Coordinators {
			coords[i] = api.V0042Coord{
				Name: coordName,
			}
		}
		account.Coordinators = &coords
	}

	// Note: Skip flags for now as AccountCreate type doesn't have Flags field

	return apiReq, nil
}
