// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_44

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_44"
)

// enhanceAssociationWithSkippedFields adds the skipped fields to an Association after base conversion
func (a *AssociationAdapter) enhanceAssociationWithSkippedFields(result *types.Association, apiObj api.V0044Assoc) {
	if result == nil {
		return
	}
	// Accounting
	if apiObj.Accounting != nil {
		result.Accounting = convertAPIAccountingListToCommon(*apiObj.Accounting)
	}
	// Default
	if apiObj.Default != nil {
		result.Default = &types.AssociationDefault{
			QoS: apiObj.Default.Qos,
		}
	}
	// Flags
	if apiObj.Flags != nil {
		for _, flag := range *apiObj.Flags {
			result.Flags = append(result.Flags, types.AssociationDefaultFlagsValue(flag))
		}
	}
	// Max
	if apiObj.Max != nil {
		result.Max = convertAPIAssocMaxToCommon(apiObj.Max)
	}
	// Min
	if apiObj.Min != nil {
		result.Min = convertAPIAssocMinToCommon(apiObj.Min)
	}
}

// convertAPIAssocMaxToCommon converts API Association Max to common type
func convertAPIAssocMaxToCommon(apiMax *struct {
	Jobs *struct {
		Accruing *api.V0044Uint32NoValStruct `json:"accruing,omitempty"`
		Active   *api.V0044Uint32NoValStruct `json:"active,omitempty"`
		Per      *struct {
			Accruing  *api.V0044Uint32NoValStruct `json:"accruing,omitempty"`
			Count     *api.V0044Uint32NoValStruct `json:"count,omitempty"`
			Submitted *api.V0044Uint32NoValStruct `json:"submitted,omitempty"`
			WallClock *api.V0044Uint32NoValStruct `json:"wall_clock,omitempty"`
		} `json:"per,omitempty"`
		Total *api.V0044Uint32NoValStruct `json:"total,omitempty"`
	} `json:"jobs,omitempty"`
	Per *struct {
		Account *struct {
			WallClock *api.V0044Uint32NoValStruct `json:"wall_clock,omitempty"`
		} `json:"account,omitempty"`
	} `json:"per,omitempty"`
	Tres *struct {
		Group *struct {
			Active  *api.V0044TresList `json:"active,omitempty"`
			Minutes *api.V0044TresList `json:"minutes,omitempty"`
		} `json:"group,omitempty"`
		Minutes *struct {
			Per *struct {
				Job *api.V0044TresList `json:"job,omitempty"`
			} `json:"per,omitempty"`
			Total *api.V0044TresList `json:"total,omitempty"`
		} `json:"minutes,omitempty"`
		Per *struct {
			Job  *api.V0044TresList `json:"job,omitempty"`
			Node *api.V0044TresList `json:"node,omitempty"`
		} `json:"per,omitempty"`
		Total *api.V0044TresList `json:"total,omitempty"`
	} `json:"tres,omitempty"`
}) *types.AssociationMax {
	if apiMax == nil {
		return nil
	}
	result := &types.AssociationMax{}
	// Jobs
	if apiMax.Jobs != nil {
		result.Jobs = &types.AssociationMaxJobs{}
		if apiMax.Jobs.Accruing != nil && apiMax.Jobs.Accruing.Set != nil && *apiMax.Jobs.Accruing.Set && apiMax.Jobs.Accruing.Number != nil {
			v := uint32(*apiMax.Jobs.Accruing.Number)
			result.Jobs.Accruing = &v
		}
		if apiMax.Jobs.Active != nil && apiMax.Jobs.Active.Set != nil && *apiMax.Jobs.Active.Set && apiMax.Jobs.Active.Number != nil {
			v := uint32(*apiMax.Jobs.Active.Number)
			result.Jobs.Active = &v
		}
		if apiMax.Jobs.Total != nil && apiMax.Jobs.Total.Set != nil && *apiMax.Jobs.Total.Set && apiMax.Jobs.Total.Number != nil {
			v := uint32(*apiMax.Jobs.Total.Number)
			result.Jobs.Total = &v
		}
		if apiMax.Jobs.Per != nil {
			result.Jobs.Per = &types.AssociationMaxJobsPer{}
			if apiMax.Jobs.Per.Accruing != nil && apiMax.Jobs.Per.Accruing.Set != nil && *apiMax.Jobs.Per.Accruing.Set && apiMax.Jobs.Per.Accruing.Number != nil {
				v := uint32(*apiMax.Jobs.Per.Accruing.Number)
				result.Jobs.Per.Accruing = &v
			}
			if apiMax.Jobs.Per.Count != nil && apiMax.Jobs.Per.Count.Set != nil && *apiMax.Jobs.Per.Count.Set && apiMax.Jobs.Per.Count.Number != nil {
				v := uint32(*apiMax.Jobs.Per.Count.Number)
				result.Jobs.Per.Count = &v
			}
			if apiMax.Jobs.Per.Submitted != nil && apiMax.Jobs.Per.Submitted.Set != nil && *apiMax.Jobs.Per.Submitted.Set && apiMax.Jobs.Per.Submitted.Number != nil {
				v := uint32(*apiMax.Jobs.Per.Submitted.Number)
				result.Jobs.Per.Submitted = &v
			}
			if apiMax.Jobs.Per.WallClock != nil && apiMax.Jobs.Per.WallClock.Set != nil && *apiMax.Jobs.Per.WallClock.Set && apiMax.Jobs.Per.WallClock.Number != nil {
				v := uint32(*apiMax.Jobs.Per.WallClock.Number)
				result.Jobs.Per.WallClock = &v
			}
		}
	}
	// Per
	if apiMax.Per != nil && apiMax.Per.Account != nil {
		result.Per = &types.AssociationMaxPer{
			Account: &types.AssociationMaxPerAccount{},
		}
		if apiMax.Per.Account.WallClock != nil && apiMax.Per.Account.WallClock.Set != nil && *apiMax.Per.Account.WallClock.Set && apiMax.Per.Account.WallClock.Number != nil {
			v := uint32(*apiMax.Per.Account.WallClock.Number)
			result.Per.Account.WallClock = &v
		}
	}
	// TRES
	if apiMax.Tres != nil {
		result.TRES = &types.AssociationMaxTRES{}
		if apiMax.Tres.Total != nil {
			result.TRES.Total = convertAPITresListToCommon(*apiMax.Tres.Total)
		}
		if apiMax.Tres.Group != nil {
			result.TRES.Group = &types.AssociationMaxTRESGroup{}
			if apiMax.Tres.Group.Active != nil {
				result.TRES.Group.Active = convertAPITresListToCommon(*apiMax.Tres.Group.Active)
			}
			if apiMax.Tres.Group.Minutes != nil {
				result.TRES.Group.Minutes = convertAPITresListToCommon(*apiMax.Tres.Group.Minutes)
			}
		}
		if apiMax.Tres.Minutes != nil {
			result.TRES.Minutes = &types.AssociationMaxTRESMinutes{}
			if apiMax.Tres.Minutes.Total != nil {
				result.TRES.Minutes.Total = convertAPITresListToCommon(*apiMax.Tres.Minutes.Total)
			}
			if apiMax.Tres.Minutes.Per != nil && apiMax.Tres.Minutes.Per.Job != nil {
				result.TRES.Minutes.Per = &types.AssociationMaxTRESMinutesPer{
					Job: convertAPITresListToCommon(*apiMax.Tres.Minutes.Per.Job),
				}
			}
		}
		if apiMax.Tres.Per != nil {
			result.TRES.Per = &types.AssociationMaxTRESPer{}
			if apiMax.Tres.Per.Job != nil {
				result.TRES.Per.Job = convertAPITresListToCommon(*apiMax.Tres.Per.Job)
			}
			if apiMax.Tres.Per.Node != nil {
				result.TRES.Per.Node = convertAPITresListToCommon(*apiMax.Tres.Per.Node)
			}
		}
	}
	return result
}

// convertAPIAssocMinToCommon converts API Association Min to common type
func convertAPIAssocMinToCommon(apiMin *struct {
	PriorityThreshold *api.V0044Uint32NoValStruct `json:"priority_threshold,omitempty"`
}) *types.AssociationMin {
	if apiMin == nil {
		return nil
	}
	result := &types.AssociationMin{}
	if apiMin.PriorityThreshold != nil && apiMin.PriorityThreshold.Set != nil && *apiMin.PriorityThreshold.Set && apiMin.PriorityThreshold.Number != nil {
		v := uint32(*apiMin.PriorityThreshold.Number)
		result.PriorityThreshold = &v
	}
	return result
}
