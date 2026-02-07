// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_44

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_44"
)

// convertAPIShareToCommon converts API V0044AssocSharesObjWrap to common Share
func (a *StandaloneAdapter) convertAPIShareToCommon(apiShare api.V0044AssocSharesObjWrap) types.Share {
	result := types.Share{}
	if apiShare.Cluster != nil {
		result.Cluster = *apiShare.Cluster
	}
	if apiShare.Name != nil {
		// Name could be account or user depending on Type
		result.Account = *apiShare.Name
	}
	if apiShare.Partition != nil {
		result.Partition = *apiShare.Partition
	}
	if apiShare.Id != nil {
		result.AssocID = int(*apiShare.Id)
	}
	// EffectiveUsage
	if apiShare.EffectiveUsage != nil && apiShare.EffectiveUsage.Set != nil && *apiShare.EffectiveUsage.Set {
		if apiShare.EffectiveUsage.Number != nil {
			result.EffectiveUsage = *apiShare.EffectiveUsage.Number
		}
	}
	// Fairshare
	if apiShare.Fairshare != nil {
		if apiShare.Fairshare.Factor != nil && apiShare.Fairshare.Factor.Set != nil && *apiShare.Fairshare.Factor.Set {
			if apiShare.Fairshare.Factor.Number != nil {
				result.FairshareUsage = *apiShare.Fairshare.Factor.Number
			}
		}
		if apiShare.Fairshare.Level != nil && apiShare.Fairshare.Level.Set != nil && *apiShare.Fairshare.Level.Set {
			if apiShare.Fairshare.Level.Number != nil {
				result.FairshareLevel = *apiShare.Fairshare.Level.Number
			}
		}
	}
	// Shares
	if apiShare.Shares != nil && apiShare.Shares.Set != nil && *apiShare.Shares.Set {
		if apiShare.Shares.Number != nil {
			result.RawShares = int(*apiShare.Shares.Number)
			result.FairshareShares = int(*apiShare.Shares.Number)
		}
	}
	// SharesNormalized
	if apiShare.SharesNormalized != nil && apiShare.SharesNormalized.Set != nil && *apiShare.SharesNormalized.Set {
		if apiShare.SharesNormalized.Number != nil {
			result.NormalizedShares = *apiShare.SharesNormalized.Number
		}
	}
	// Usage
	if apiShare.Usage != nil {
		result.RawUsage = *apiShare.Usage
	}
	// UsageNormalized
	if apiShare.UsageNormalized != nil && apiShare.UsageNormalized.Set != nil && *apiShare.UsageNormalized.Set {
		if apiShare.UsageNormalized.Number != nil {
			result.NormalizedUsage = *apiShare.UsageNormalized.Number
		}
	}
	return result
}
