// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_44

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_44"
)

// enhanceQoSWithSkippedFields adds the skipped fields to a QoS after base conversion
func (a *QoSAdapter) enhanceQoSWithSkippedFields(result *types.QoS, apiObj api.V0044Qos) {
	if result == nil {
		return
	}
	// Limits
	if apiObj.Limits != nil {
		result.Limits = convertAPIQoSLimitsToCommon(apiObj.Limits)
	}
	// Preempt
	if apiObj.Preempt != nil {
		result.Preempt = convertAPIQoSPreemptToCommon(apiObj.Preempt)
	}
}

// convertAPIQoSLimitsToCommon converts API QoS Limits to common type
func convertAPIQoSLimitsToCommon(apiLimits *struct {
	Factor    *api.V0044Float64NoValStruct `json:"factor,omitempty"`
	GraceTime *int32                       `json:"grace_time,omitempty"`
	Max       *struct {
		Accruing *struct {
			Per *struct {
				Account *api.V0044Uint32NoValStruct `json:"account,omitempty"`
				User    *api.V0044Uint32NoValStruct `json:"user,omitempty"`
			} `json:"per,omitempty"`
		} `json:"accruing,omitempty"`
		ActiveJobs *struct {
			Accruing *api.V0044Uint32NoValStruct `json:"accruing,omitempty"`
			Count    *api.V0044Uint32NoValStruct `json:"count,omitempty"`
		} `json:"active_jobs,omitempty"`
		Jobs *struct {
			ActiveJobs *struct {
				Per *struct {
					Account *api.V0044Uint32NoValStruct `json:"account,omitempty"`
					User    *api.V0044Uint32NoValStruct `json:"user,omitempty"`
				} `json:"per,omitempty"`
			} `json:"active_jobs,omitempty"`
			Count *api.V0044Uint32NoValStruct `json:"count,omitempty"`
			Per   *struct {
				Account *api.V0044Uint32NoValStruct `json:"account,omitempty"`
				User    *api.V0044Uint32NoValStruct `json:"user,omitempty"`
			} `json:"per,omitempty"`
		} `json:"jobs,omitempty"`
		Tres *struct {
			Minutes *struct {
				Per *struct {
					Account *api.V0044TresList `json:"account,omitempty"`
					Job     *api.V0044TresList `json:"job,omitempty"`
					Qos     *api.V0044TresList `json:"qos,omitempty"`
					User    *api.V0044TresList `json:"user,omitempty"`
				} `json:"per,omitempty"`
				Total *api.V0044TresList `json:"total,omitempty"`
			} `json:"minutes,omitempty"`
			Per *struct {
				Account *api.V0044TresList `json:"account,omitempty"`
				Job     *api.V0044TresList `json:"job,omitempty"`
				Node    *api.V0044TresList `json:"node,omitempty"`
				User    *api.V0044TresList `json:"user,omitempty"`
			} `json:"per,omitempty"`
			Total *api.V0044TresList `json:"total,omitempty"`
		} `json:"tres,omitempty"`
		WallClock *struct {
			Per *struct {
				Job *api.V0044Uint32NoValStruct `json:"job,omitempty"`
				Qos *api.V0044Uint32NoValStruct `json:"qos,omitempty"`
			} `json:"per,omitempty"`
		} `json:"wall_clock,omitempty"`
	} `json:"max,omitempty"`
	Min *struct {
		PriorityThreshold *api.V0044Uint32NoValStruct `json:"priority_threshold,omitempty"`
		Tres              *struct {
			Per *struct {
				Job *api.V0044TresList `json:"job,omitempty"`
			} `json:"per,omitempty"`
		} `json:"tres,omitempty"`
	} `json:"min,omitempty"`
}) *types.QoSLimits {
	if apiLimits == nil {
		return nil
	}
	result := &types.QoSLimits{}
	// Factor
	if apiLimits.Factor != nil && apiLimits.Factor.Set != nil && *apiLimits.Factor.Set && apiLimits.Factor.Number != nil {
		v := float64(*apiLimits.Factor.Number)
		result.Factor = &v
	}
	// GraceTime
	if apiLimits.GraceTime != nil {
		result.GraceTime = apiLimits.GraceTime
	}
	// Max
	if apiLimits.Max != nil {
		result.Max = convertAPIQoSLimitsMaxToCommon(apiLimits.Max)
	}
	// Min
	if apiLimits.Min != nil {
		result.Min = convertAPIQoSLimitsMinToCommon(apiLimits.Min)
	}
	return result
}

// convertAPIQoSLimitsMaxToCommon converts API QoS Limits Max to common type
func convertAPIQoSLimitsMaxToCommon(apiMax *struct {
	Accruing *struct {
		Per *struct {
			Account *api.V0044Uint32NoValStruct `json:"account,omitempty"`
			User    *api.V0044Uint32NoValStruct `json:"user,omitempty"`
		} `json:"per,omitempty"`
	} `json:"accruing,omitempty"`
	ActiveJobs *struct {
		Accruing *api.V0044Uint32NoValStruct `json:"accruing,omitempty"`
		Count    *api.V0044Uint32NoValStruct `json:"count,omitempty"`
	} `json:"active_jobs,omitempty"`
	Jobs *struct {
		ActiveJobs *struct {
			Per *struct {
				Account *api.V0044Uint32NoValStruct `json:"account,omitempty"`
				User    *api.V0044Uint32NoValStruct `json:"user,omitempty"`
			} `json:"per,omitempty"`
		} `json:"active_jobs,omitempty"`
		Count *api.V0044Uint32NoValStruct `json:"count,omitempty"`
		Per   *struct {
			Account *api.V0044Uint32NoValStruct `json:"account,omitempty"`
			User    *api.V0044Uint32NoValStruct `json:"user,omitempty"`
		} `json:"per,omitempty"`
	} `json:"jobs,omitempty"`
	Tres *struct {
		Minutes *struct {
			Per *struct {
				Account *api.V0044TresList `json:"account,omitempty"`
				Job     *api.V0044TresList `json:"job,omitempty"`
				Qos     *api.V0044TresList `json:"qos,omitempty"`
				User    *api.V0044TresList `json:"user,omitempty"`
			} `json:"per,omitempty"`
			Total *api.V0044TresList `json:"total,omitempty"`
		} `json:"minutes,omitempty"`
		Per *struct {
			Account *api.V0044TresList `json:"account,omitempty"`
			Job     *api.V0044TresList `json:"job,omitempty"`
			Node    *api.V0044TresList `json:"node,omitempty"`
			User    *api.V0044TresList `json:"user,omitempty"`
		} `json:"per,omitempty"`
		Total *api.V0044TresList `json:"total,omitempty"`
	} `json:"tres,omitempty"`
	WallClock *struct {
		Per *struct {
			Job *api.V0044Uint32NoValStruct `json:"job,omitempty"`
			Qos *api.V0044Uint32NoValStruct `json:"qos,omitempty"`
		} `json:"per,omitempty"`
	} `json:"wall_clock,omitempty"`
}) *types.QoSLimitsMax {
	if apiMax == nil {
		return nil
	}
	result := &types.QoSLimitsMax{}
	// Accruing
	if apiMax.Accruing != nil && apiMax.Accruing.Per != nil {
		result.Accruing = &types.QoSLimitsMaxAccruing{
			Per: &types.QoSLimitsMaxAccruingPer{},
		}
		if apiMax.Accruing.Per.Account != nil && apiMax.Accruing.Per.Account.Set != nil && *apiMax.Accruing.Per.Account.Set && apiMax.Accruing.Per.Account.Number != nil {
			v := uint32(*apiMax.Accruing.Per.Account.Number)
			result.Accruing.Per.Account = &v
		}
		if apiMax.Accruing.Per.User != nil && apiMax.Accruing.Per.User.Set != nil && *apiMax.Accruing.Per.User.Set && apiMax.Accruing.Per.User.Number != nil {
			v := uint32(*apiMax.Accruing.Per.User.Number)
			result.Accruing.Per.User = &v
		}
	}
	// ActiveJobs
	if apiMax.ActiveJobs != nil {
		result.ActiveJobs = &types.QoSLimitsMaxActiveJobs{}
		if apiMax.ActiveJobs.Accruing != nil && apiMax.ActiveJobs.Accruing.Set != nil && *apiMax.ActiveJobs.Accruing.Set && apiMax.ActiveJobs.Accruing.Number != nil {
			v := uint32(*apiMax.ActiveJobs.Accruing.Number)
			result.ActiveJobs.Accruing = &v
		}
		if apiMax.ActiveJobs.Count != nil && apiMax.ActiveJobs.Count.Set != nil && *apiMax.ActiveJobs.Count.Set && apiMax.ActiveJobs.Count.Number != nil {
			v := uint32(*apiMax.ActiveJobs.Count.Number)
			result.ActiveJobs.Count = &v
		}
	}
	// Jobs
	if apiMax.Jobs != nil {
		result.Jobs = &types.QoSLimitsMaxJobs{}
		if apiMax.Jobs.Count != nil && apiMax.Jobs.Count.Set != nil && *apiMax.Jobs.Count.Set && apiMax.Jobs.Count.Number != nil {
			v := uint32(*apiMax.Jobs.Count.Number)
			result.Jobs.Count = &v
		}
		if apiMax.Jobs.ActiveJobs != nil && apiMax.Jobs.ActiveJobs.Per != nil {
			result.Jobs.ActiveJobs = &types.QoSLimitsMaxJobsActiveJobs{
				Per: &types.QoSLimitsMaxJobsActiveJobsPer{},
			}
			if apiMax.Jobs.ActiveJobs.Per.Account != nil && apiMax.Jobs.ActiveJobs.Per.Account.Set != nil && *apiMax.Jobs.ActiveJobs.Per.Account.Set && apiMax.Jobs.ActiveJobs.Per.Account.Number != nil {
				v := uint32(*apiMax.Jobs.ActiveJobs.Per.Account.Number)
				result.Jobs.ActiveJobs.Per.Account = &v
			}
			if apiMax.Jobs.ActiveJobs.Per.User != nil && apiMax.Jobs.ActiveJobs.Per.User.Set != nil && *apiMax.Jobs.ActiveJobs.Per.User.Set && apiMax.Jobs.ActiveJobs.Per.User.Number != nil {
				v := uint32(*apiMax.Jobs.ActiveJobs.Per.User.Number)
				result.Jobs.ActiveJobs.Per.User = &v
			}
		}
		if apiMax.Jobs.Per != nil {
			result.Jobs.Per = &types.QoSLimitsMaxJobsPer{}
			if apiMax.Jobs.Per.Account != nil && apiMax.Jobs.Per.Account.Set != nil && *apiMax.Jobs.Per.Account.Set && apiMax.Jobs.Per.Account.Number != nil {
				v := uint32(*apiMax.Jobs.Per.Account.Number)
				result.Jobs.Per.Account = &v
			}
			if apiMax.Jobs.Per.User != nil && apiMax.Jobs.Per.User.Set != nil && *apiMax.Jobs.Per.User.Set && apiMax.Jobs.Per.User.Number != nil {
				v := uint32(*apiMax.Jobs.Per.User.Number)
				result.Jobs.Per.User = &v
			}
		}
	}
	// TRES
	if apiMax.Tres != nil {
		result.TRES = &types.QoSLimitsMaxTRES{}
		if apiMax.Tres.Total != nil {
			result.TRES.Total = convertAPITresListToCommon(*apiMax.Tres.Total)
		}
		if apiMax.Tres.Minutes != nil {
			result.TRES.Minutes = &types.QoSLimitsMaxTRESMinutes{}
			if apiMax.Tres.Minutes.Total != nil {
				result.TRES.Minutes.Total = convertAPITresListToCommon(*apiMax.Tres.Minutes.Total)
			}
			if apiMax.Tres.Minutes.Per != nil {
				result.TRES.Minutes.Per = &types.QoSLimitsMaxTRESMinutesPer{}
				if apiMax.Tres.Minutes.Per.Account != nil {
					result.TRES.Minutes.Per.Account = convertAPITresListToCommon(*apiMax.Tres.Minutes.Per.Account)
				}
				if apiMax.Tres.Minutes.Per.Job != nil {
					result.TRES.Minutes.Per.Job = convertAPITresListToCommon(*apiMax.Tres.Minutes.Per.Job)
				}
				if apiMax.Tres.Minutes.Per.Qos != nil {
					result.TRES.Minutes.Per.QoS = convertAPITresListToCommon(*apiMax.Tres.Minutes.Per.Qos)
				}
				if apiMax.Tres.Minutes.Per.User != nil {
					result.TRES.Minutes.Per.User = convertAPITresListToCommon(*apiMax.Tres.Minutes.Per.User)
				}
			}
		}
		if apiMax.Tres.Per != nil {
			result.TRES.Per = &types.QoSLimitsMaxTRESPer{}
			if apiMax.Tres.Per.Account != nil {
				result.TRES.Per.Account = convertAPITresListToCommon(*apiMax.Tres.Per.Account)
			}
			if apiMax.Tres.Per.Job != nil {
				result.TRES.Per.Job = convertAPITresListToCommon(*apiMax.Tres.Per.Job)
			}
			if apiMax.Tres.Per.Node != nil {
				result.TRES.Per.Node = convertAPITresListToCommon(*apiMax.Tres.Per.Node)
			}
			if apiMax.Tres.Per.User != nil {
				result.TRES.Per.User = convertAPITresListToCommon(*apiMax.Tres.Per.User)
			}
		}
	}
	// WallClock
	if apiMax.WallClock != nil && apiMax.WallClock.Per != nil {
		result.WallClock = &types.QoSLimitsMaxWallClock{
			Per: &types.QoSLimitsMaxWallClockPer{},
		}
		if apiMax.WallClock.Per.Job != nil && apiMax.WallClock.Per.Job.Set != nil && *apiMax.WallClock.Per.Job.Set && apiMax.WallClock.Per.Job.Number != nil {
			v := uint32(*apiMax.WallClock.Per.Job.Number)
			result.WallClock.Per.Job = &v
		}
		if apiMax.WallClock.Per.Qos != nil && apiMax.WallClock.Per.Qos.Set != nil && *apiMax.WallClock.Per.Qos.Set && apiMax.WallClock.Per.Qos.Number != nil {
			v := uint32(*apiMax.WallClock.Per.Qos.Number)
			result.WallClock.Per.QoS = &v
		}
	}
	return result
}

// convertAPIQoSLimitsMinToCommon converts API QoS Limits Min to common type
func convertAPIQoSLimitsMinToCommon(apiMin *struct {
	PriorityThreshold *api.V0044Uint32NoValStruct `json:"priority_threshold,omitempty"`
	Tres              *struct {
		Per *struct {
			Job *api.V0044TresList `json:"job,omitempty"`
		} `json:"per,omitempty"`
	} `json:"tres,omitempty"`
}) *types.QoSLimitsMin {
	if apiMin == nil {
		return nil
	}
	result := &types.QoSLimitsMin{}
	// PriorityThreshold
	if apiMin.PriorityThreshold != nil && apiMin.PriorityThreshold.Set != nil && *apiMin.PriorityThreshold.Set && apiMin.PriorityThreshold.Number != nil {
		v := uint32(*apiMin.PriorityThreshold.Number)
		result.PriorityThreshold = &v
	}
	// TRES
	if apiMin.Tres != nil && apiMin.Tres.Per != nil && apiMin.Tres.Per.Job != nil {
		result.TRES = &types.QoSLimitsMinTRES{
			Per: &types.QoSLimitsMinTRESPer{
				Job: convertAPITresListToCommon(*apiMin.Tres.Per.Job),
			},
		}
	}
	return result
}

// convertAPIQoSPreemptToCommon converts API QoS Preempt to common type
func convertAPIQoSPreemptToCommon(apiPreempt *struct {
	ExemptTime *api.V0044Uint32NoValStruct `json:"exempt_time,omitempty"`
	List       *api.V0044QosPreemptList    `json:"list,omitempty"`
	Mode       *[]api.V0044QosPreemptMode  `json:"mode,omitempty"`
}) *types.QoSPreempt {
	if apiPreempt == nil {
		return nil
	}
	result := &types.QoSPreempt{}
	// ExemptTime
	if apiPreempt.ExemptTime != nil && apiPreempt.ExemptTime.Set != nil && *apiPreempt.ExemptTime.Set && apiPreempt.ExemptTime.Number != nil {
		v := uint32(*apiPreempt.ExemptTime.Number)
		result.ExemptTime = &v
	}
	// List
	if apiPreempt.List != nil {
		result.List = *apiPreempt.List
	}
	// Mode
	if apiPreempt.Mode != nil {
		for _, m := range *apiPreempt.Mode {
			result.Mode = append(result.Mode, types.ModeValue(m))
		}
	}
	return result
}

// convertAPITresListToCommon converts API TRES list to common type
func convertAPITresListToCommon(apiTres api.V0044TresList) []types.TRES {
	if len(apiTres) == 0 {
		return nil
	}
	result := make([]types.TRES, 0, len(apiTres))
	for _, t := range apiTres {
		result = append(result, types.TRES{
			Count: t.Count,
			ID:    t.Id,
			Name:  t.Name,
			Type:  t.Type,
		})
	}
	return result
}
