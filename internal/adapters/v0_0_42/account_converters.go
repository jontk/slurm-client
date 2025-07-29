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
	if apiAccount.Flags != nil && len(*apiAccount.Flags) > 0 {
		account.Flags = *apiAccount.Flags
	}

	// Convert coordinators if present
	if apiAccount.Coordinators != nil && len(*apiAccount.Coordinators) > 0 {
		coordinators := make([]string, 0, len(*apiAccount.Coordinators))
		for _, coord := range *apiAccount.Coordinators {
			if coord.Name != nil {
				coordinators = append(coordinators, *coord.Name)
			}
		}
		account.Coordinators = coordinators
	}

	// Handle associations if present
	if apiAccount.Associations != nil && len(*apiAccount.Associations) > 0 {
		// Store association count as metadata
		account.AssociationCount = len(*apiAccount.Associations)
		
		// Extract some key association info
		for _, assoc := range *apiAccount.Associations {
			// If this is the parent association (where user is empty), capture some details
			if assoc.User == nil || (assoc.User != nil && *assoc.User == "") {
				if assoc.Qos != nil && len(*assoc.Qos) > 0 {
					account.DefaultQoS = (*assoc.Qos)[0]
				}
				if assoc.GrpJobs != nil && assoc.GrpJobs.Number > 0 {
					grpJobs := uint32(assoc.GrpJobs.Number)
					account.GrpJobs = &grpJobs
				}
				if assoc.GrpCpus != nil && assoc.GrpCpus.Number > 0 {
					grpCPUs := uint32(assoc.GrpCpus.Number)
					account.GrpCPUs = &grpCPUs
				}
				if assoc.GrpMem != nil && assoc.GrpMem.Number > 0 {
					grpMem := uint64(assoc.GrpMem.Number)
					account.GrpMem = &grpMem
				}
				if assoc.GrpNodes != nil && assoc.GrpNodes.Number > 0 {
					grpNodes := uint32(assoc.GrpNodes.Number)
					account.GrpNodes = &grpNodes
				}
				break
			}
		}
	}

	return account, nil
}

// convertCommonAccountCreateToAPI converts common account create request to v0.0.42 API format
func (a *AccountAdapter) convertCommonAccountCreateToAPI(req *types.AccountCreateRequest) (*api.SlurmdbV0042PostAccountsJSONRequestBody, error) {
	apiReq := &api.SlurmdbV0042PostAccountsJSONRequestBody{
		Accounts: &[]api.V0042Account{
			{
				Name:         req.Name,
				Description:  req.Description,
				Organization: req.Organization,
			},
		},
	}

	account := &(*apiReq.Accounts)[0]

	// Add coordinators if specified
	if len(req.Coordinators) > 0 {
		coords := make(api.V0042CoordList, len(req.Coordinators))
		for i, coordName := range req.Coordinators {
			coords[i] = api.V0042Coord{
				Name: &coordName,
			}
		}
		account.Coordinators = &coords
	}

	// Add flags if specified
	if len(req.Flags) > 0 {
		flags := api.V0042AccountFlags(req.Flags)
		account.Flags = &flags
	}

	return apiReq, nil
}