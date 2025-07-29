package v0_0_42

import (
	"fmt"
	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

// convertAPIAssociationToCommon converts a v0.0.42 API Association to common Association type
func (a *AssociationAdapter) convertAPIAssociationToCommon(apiAssoc api.V0042Assoc) (*types.Association, error) {
	assoc := &types.Association{}

	// Basic fields
	if apiAssoc.Id != nil {
		assoc.ID = fmt.Sprintf("%d", *apiAssoc.Id)
	}
	if apiAssoc.Account != nil {
		assoc.AccountName = *apiAssoc.Account
	}
	if apiAssoc.User != "" {
		assoc.UserName = apiAssoc.User
	}
	if apiAssoc.Cluster != nil {
		assoc.Cluster = *apiAssoc.Cluster
	}
	if apiAssoc.Partition != nil {
		assoc.Partition = *apiAssoc.Partition
	}
	if apiAssoc.Comment != nil {
		assoc.Comment = *apiAssoc.Comment
	}

	// Parent account
	if apiAssoc.ParentAccount != nil {
		assoc.ParentAccount = *apiAssoc.ParentAccount
	}

	// Default QoS
	if apiAssoc.Default != nil && apiAssoc.Default.Qos != nil {
		assoc.DefaultQoS = *apiAssoc.Default.Qos
	}

	// QoS list
	if apiAssoc.Qos != nil && len(*apiAssoc.Qos) > 0 {
		assoc.QoSList = *apiAssoc.Qos
	}

	// Is default
	if apiAssoc.IsDefault != nil {
		assoc.IsDefault = *apiAssoc.IsDefault
	}

	// Priority
	if apiAssoc.Priority != nil && apiAssoc.Priority.Number != nil {
		assoc.Priority = *apiAssoc.Priority.Number
	}

	// Shares
	if apiAssoc.SharesRaw != nil {
		assoc.SharesRaw = *apiAssoc.SharesRaw
	}

	// Limits
	if apiAssoc.Max != nil {
		// Jobs limits
		if apiAssoc.Max.Jobs != nil {
			if apiAssoc.Max.Jobs.Total != nil && apiAssoc.Max.Jobs.Total.Number != nil {
				assoc.MaxJobs = *apiAssoc.Max.Jobs.Total.Number
			}
			if apiAssoc.Max.Jobs.Accruing != nil && apiAssoc.Max.Jobs.Accruing.Number != nil {
				assoc.MaxJobsAccrue = *apiAssoc.Max.Jobs.Accruing.Number
			}
			if apiAssoc.Max.Jobs.Per != nil {
				if apiAssoc.Max.Jobs.Per.Submitted != nil && apiAssoc.Max.Jobs.Per.Submitted.Number != nil {
					assoc.MaxSubmitJobs = *apiAssoc.Max.Jobs.Per.Submitted.Number
				}
				if apiAssoc.Max.Jobs.Per.WallClock != nil && apiAssoc.Max.Jobs.Per.WallClock.Number != nil {
					assoc.MaxWallTime = *apiAssoc.Max.Jobs.Per.WallClock.Number
				}
			}
		}
	}

	return assoc, nil
}