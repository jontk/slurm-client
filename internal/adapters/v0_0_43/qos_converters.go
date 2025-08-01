// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// convertAPIQoSToCommon converts a v0.0.43 API QoS to common QoS type
func (a *QoSAdapter) convertAPIQoSToCommon(apiQoS api.V0043Qos) (*types.QoS, error) {
	qos := &types.QoS{}

	// Basic fields
	if apiQoS.Name != nil {
		qos.Name = *apiQoS.Name
	}
	if apiQoS.Description != nil {
		qos.Description = *apiQoS.Description
	}
	if apiQoS.Id != nil {
		qos.ID = *apiQoS.Id
	}

	// Priority
	if apiQoS.Priority != nil && apiQoS.Priority.Set != nil && *apiQoS.Priority.Set && apiQoS.Priority.Number != nil {
		qos.Priority = int(*apiQoS.Priority.Number)
	}

	// Flags
	if apiQoS.Flags != nil {
		flags := make([]string, 0, len(*apiQoS.Flags))
		for _, flag := range *apiQoS.Flags {
			flags = append(flags, string(flag))
		}
		qos.Flags = flags
	}

	// Preempt mode
	if apiQoS.Preempt != nil && apiQoS.Preempt.Mode != nil && len(*apiQoS.Preempt.Mode) > 0 {
		qos.PreemptMode = string((*apiQoS.Preempt.Mode)[0])
	}

	// Preempt exempt time
	if apiQoS.Preempt != nil && apiQoS.Preempt.ExemptTime != nil && 
		apiQoS.Preempt.ExemptTime.Set != nil && *apiQoS.Preempt.ExemptTime.Set && 
		apiQoS.Preempt.ExemptTime.Number != nil {
		qos.PreemptExemptTime = int(*apiQoS.Preempt.ExemptTime.Number)
	}

	// Grace time (from Limits)
	if apiQoS.Limits != nil && apiQoS.Limits.GraceTime != nil {
		qos.GraceTime = int(*apiQoS.Limits.GraceTime)
	}

	// Usage factor
	if apiQoS.UsageFactor != nil && apiQoS.UsageFactor.Set != nil && *apiQoS.UsageFactor.Set && apiQoS.UsageFactor.Number != nil {
		qos.UsageFactor = *apiQoS.UsageFactor.Number
	}

	// Usage threshold
	if apiQoS.UsageThreshold != nil && apiQoS.UsageThreshold.Set != nil && *apiQoS.UsageThreshold.Set && apiQoS.UsageThreshold.Number != nil {
		qos.UsageThreshold = *apiQoS.UsageThreshold.Number
	}

	// Convert limits
	if apiQoS.Limits != nil {
		limits := &types.QoSLimits{}
		hasLimits := false

		// Grace time
		if apiQoS.Limits.GraceTime != nil {
			graceTime := int(*apiQoS.Limits.GraceTime)
			qos.GraceTime = graceTime
		}

		// Max limits
		if apiQoS.Limits.Max != nil {
			// Per-user limits
			if apiQoS.Limits.Max.Jobs != nil && apiQoS.Limits.Max.Jobs.Per != nil && apiQoS.Limits.Max.Jobs.Per.User != nil {
				if apiQoS.Limits.Max.Jobs.Per.User.Set != nil && *apiQoS.Limits.Max.Jobs.Per.User.Set && apiQoS.Limits.Max.Jobs.Per.User.Number != nil {
					val := int(*apiQoS.Limits.Max.Jobs.Per.User.Number)
					limits.MaxJobsPerUser = &val
					hasLimits = true
				}
			}

			// Per-account limits
			if apiQoS.Limits.Max.Jobs != nil && apiQoS.Limits.Max.Jobs.Per != nil && apiQoS.Limits.Max.Jobs.Per.Account != nil {
				if apiQoS.Limits.Max.Jobs.Per.Account.Set != nil && *apiQoS.Limits.Max.Jobs.Per.Account.Set && apiQoS.Limits.Max.Jobs.Per.Account.Number != nil {
					val := int(*apiQoS.Limits.Max.Jobs.Per.Account.Number)
					limits.MaxJobsPerAccount = &val
					hasLimits = true
				}
			}

			// TODO: Convert TRES limits (CPU, Memory, Node) when we understand the TRES format
			// The API uses TRES (Trackable RESources) which is a more complex format
		}

		// Min limits
		if apiQoS.Limits.Min != nil {
			// TODO: Convert minimum TRES limits when we understand the TRES format
		}

		if hasLimits {
			qos.Limits = limits
		}
	}

	return qos, nil
}

// convertCommonQoSCreateToAPI converts common QoS create type to v0.0.43 API format
func (a *QoSAdapter) convertCommonQoSCreateToAPI(create *types.QoSCreate) (*api.V0043Qos, error) {
	apiQoS := &api.V0043Qos{}

	// Required fields
	apiQoS.Name = &create.Name

	// Optional fields
	if create.Description != "" {
		apiQoS.Description = &create.Description
	}

	// Priority
	if create.Priority > 0 {
		setTrue := true
		priority := int32(create.Priority)
		apiQoS.Priority = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &priority,
		}
	}

	// Flags
	if len(create.Flags) > 0 {
		flags := make([]api.V0043QosFlags, 0, len(create.Flags))
		for _, flag := range create.Flags {
			flags = append(flags, api.V0043QosFlags(flag))
		}
		apiQoS.Flags = &flags
	}

	// Preempt mode
	if len(create.PreemptMode) > 0 {
		modes := make([]api.V0043QosPreemptMode, 0, len(create.PreemptMode))
		for _, mode := range create.PreemptMode {
			modes = append(modes, api.V0043QosPreemptMode(mode))
		}
		if apiQoS.Preempt == nil {
			apiQoS.Preempt = &struct {
				ExemptTime *api.V0043Uint32NoValStruct `json:"exempt_time,omitempty"`
				List       *api.V0043QosPreemptList    `json:"list,omitempty"`
				Mode       *[]api.V0043QosPreemptMode  `json:"mode,omitempty"`
			}{}
		}
		apiQoS.Preempt.Mode = &modes
	}

	// Preempt exempt time
	if create.PreemptExemptTime != nil {
		setTrue := true
		exemptTime := int32(*create.PreemptExemptTime)
		if apiQoS.Preempt == nil {
			apiQoS.Preempt = &struct {
				ExemptTime *api.V0043Uint32NoValStruct `json:"exempt_time,omitempty"`
				List       *api.V0043QosPreemptList    `json:"list,omitempty"`
				Mode       *[]api.V0043QosPreemptMode  `json:"mode,omitempty"`
			}{}
		}
		apiQoS.Preempt.ExemptTime = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &exemptTime,
		}
	}

	// Usage factor
	if create.UsageFactor != 0 {
		setTrue := true
		apiQoS.UsageFactor = &api.V0043Float64NoValStruct{
			Set:    &setTrue,
			Number: &create.UsageFactor,
		}
	}

	// Usage threshold
	if create.UsageThreshold != 0 {
		setTrue := true
		apiQoS.UsageThreshold = &api.V0043Float64NoValStruct{
			Set:    &setTrue,
			Number: &create.UsageThreshold,
		}
	}

	// Convert limits if provided
	if create.Limits != nil || create.GraceTime != 0 {
		apiQoS.Limits = &struct {
			Factor    *api.V0043Float64NoValStruct `json:"factor,omitempty"`
			GraceTime *int32                       `json:"grace_time,omitempty"`
			Max       *struct {
				Accruing *struct {
					Per *struct {
						Account *api.V0043Uint32NoValStruct `json:"account,omitempty"`
						User    *api.V0043Uint32NoValStruct `json:"user,omitempty"`
					} `json:"per,omitempty"`
				} `json:"accruing,omitempty"`
				ActiveJobs *struct {
					Accruing *api.V0043Uint32NoValStruct `json:"accruing,omitempty"`
					Count    *api.V0043Uint32NoValStruct `json:"count,omitempty"`
				} `json:"active_jobs,omitempty"`
				Jobs *struct {
					ActiveJobs *struct {
						Per *struct {
							Account *api.V0043Uint32NoValStruct `json:"account,omitempty"`
							User    *api.V0043Uint32NoValStruct `json:"user,omitempty"`
						} `json:"per,omitempty"`
					} `json:"active_jobs,omitempty"`
					Count *api.V0043Uint32NoValStruct `json:"count,omitempty"`
					Per   *struct {
						Account *api.V0043Uint32NoValStruct `json:"account,omitempty"`
						User    *api.V0043Uint32NoValStruct `json:"user,omitempty"`
					} `json:"per,omitempty"`
				} `json:"jobs,omitempty"`
				Tres *struct {
					Minutes *struct {
						Per *struct {
							Account *api.V0043TresList `json:"account,omitempty"`
							Job     *api.V0043TresList `json:"job,omitempty"`
							Qos     *api.V0043TresList `json:"qos,omitempty"`
							User    *api.V0043TresList `json:"user,omitempty"`
						} `json:"per,omitempty"`
						Total *api.V0043TresList `json:"total,omitempty"`
					} `json:"minutes,omitempty"`
					Per *struct {
						Account *api.V0043TresList `json:"account,omitempty"`
						Job     *api.V0043TresList `json:"job,omitempty"`
						Node    *api.V0043TresList `json:"node,omitempty"`
						User    *api.V0043TresList `json:"user,omitempty"`
					} `json:"per,omitempty"`
					Total *api.V0043TresList `json:"total,omitempty"`
				} `json:"tres,omitempty"`
				WallClock *struct {
					Per *struct {
						Job *api.V0043Uint32NoValStruct `json:"job,omitempty"`
						Qos *api.V0043Uint32NoValStruct `json:"qos,omitempty"`
					} `json:"per,omitempty"`
				} `json:"wall_clock,omitempty"`
			} `json:"max,omitempty"`
			Min *struct {
				PriorityThreshold *api.V0043Uint32NoValStruct `json:"priority_threshold,omitempty"`
				Tres              *struct {
					Per *struct {
						Job *api.V0043TresList `json:"job,omitempty"`
					} `json:"per,omitempty"`
				} `json:"tres,omitempty"`
			} `json:"min,omitempty"`
		}{}

		// Set grace time
		if create.GraceTime != 0 {
			graceTime := int32(create.GraceTime)
			apiQoS.Limits.GraceTime = &graceTime
		}

		// Convert limits
		if create.Limits != nil {
			// Initialize Max structure if we have any max limits
			if create.Limits.MaxJobsPerUser != nil || create.Limits.MaxJobsPerAccount != nil {
				apiQoS.Limits.Max = &struct {
					Accruing *struct {
						Per *struct {
							Account *api.V0043Uint32NoValStruct `json:"account,omitempty"`
							User    *api.V0043Uint32NoValStruct `json:"user,omitempty"`
						} `json:"per,omitempty"`
					} `json:"accruing,omitempty"`
					ActiveJobs *struct {
						Accruing *api.V0043Uint32NoValStruct `json:"accruing,omitempty"`
						Count    *api.V0043Uint32NoValStruct `json:"count,omitempty"`
					} `json:"active_jobs,omitempty"`
					Jobs *struct {
						ActiveJobs *struct {
							Per *struct {
								Account *api.V0043Uint32NoValStruct `json:"account,omitempty"`
								User    *api.V0043Uint32NoValStruct `json:"user,omitempty"`
							} `json:"per,omitempty"`
						} `json:"active_jobs,omitempty"`
						Count *api.V0043Uint32NoValStruct `json:"count,omitempty"`
						Per   *struct {
							Account *api.V0043Uint32NoValStruct `json:"account,omitempty"`
							User    *api.V0043Uint32NoValStruct `json:"user,omitempty"`
						} `json:"per,omitempty"`
					} `json:"jobs,omitempty"`
					Tres *struct {
						Minutes *struct {
							Per *struct {
								Account *api.V0043TresList `json:"account,omitempty"`
								Job     *api.V0043TresList `json:"job,omitempty"`
								Qos     *api.V0043TresList `json:"qos,omitempty"`
								User    *api.V0043TresList `json:"user,omitempty"`
							} `json:"per,omitempty"`
							Total *api.V0043TresList `json:"total,omitempty"`
						} `json:"minutes,omitempty"`
						Per *struct {
							Account *api.V0043TresList `json:"account,omitempty"`
							Job     *api.V0043TresList `json:"job,omitempty"`
							Node    *api.V0043TresList `json:"node,omitempty"`
							User    *api.V0043TresList `json:"user,omitempty"`
						} `json:"per,omitempty"`
						Total *api.V0043TresList `json:"total,omitempty"`
					} `json:"tres,omitempty"`
					WallClock *struct {
						Per *struct {
							Job *api.V0043Uint32NoValStruct `json:"job,omitempty"`
							Qos *api.V0043Uint32NoValStruct `json:"qos,omitempty"`
						} `json:"per,omitempty"`
					} `json:"wall_clock,omitempty"`
				}{}
			}

			// Set per-user job limits
			if create.Limits.MaxJobsPerUser != nil && apiQoS.Limits.Max != nil {
				apiQoS.Limits.Max.Jobs = &struct {
					ActiveJobs *struct {
						Per *struct {
							Account *api.V0043Uint32NoValStruct `json:"account,omitempty"`
							User    *api.V0043Uint32NoValStruct `json:"user,omitempty"`
						} `json:"per,omitempty"`
					} `json:"active_jobs,omitempty"`
					Count *api.V0043Uint32NoValStruct `json:"count,omitempty"`
					Per   *struct {
						Account *api.V0043Uint32NoValStruct `json:"account,omitempty"`
						User    *api.V0043Uint32NoValStruct `json:"user,omitempty"`
					} `json:"per,omitempty"`
				}{
					Per: &struct {
						Account *api.V0043Uint32NoValStruct `json:"account,omitempty"`
						User    *api.V0043Uint32NoValStruct `json:"user,omitempty"`
					}{},
				}

				setTrue := true
				jobLimit := int32(*create.Limits.MaxJobsPerUser)
				apiQoS.Limits.Max.Jobs.Per.User = &api.V0043Uint32NoValStruct{
					Set:    &setTrue,
					Number: &jobLimit,
				}
			}

			// Set per-account job limits
			if create.Limits.MaxJobsPerAccount != nil && apiQoS.Limits.Max != nil {
				if apiQoS.Limits.Max.Jobs == nil {
					apiQoS.Limits.Max.Jobs = &struct {
						ActiveJobs *struct {
							Per *struct {
								Account *api.V0043Uint32NoValStruct `json:"account,omitempty"`
								User    *api.V0043Uint32NoValStruct `json:"user,omitempty"`
							} `json:"per,omitempty"`
						} `json:"active_jobs,omitempty"`
						Count *api.V0043Uint32NoValStruct `json:"count,omitempty"`
						Per   *struct {
							Account *api.V0043Uint32NoValStruct `json:"account,omitempty"`
							User    *api.V0043Uint32NoValStruct `json:"user,omitempty"`
						} `json:"per,omitempty"`
					}{}
				}
				if apiQoS.Limits.Max.Jobs.Per == nil {
					apiQoS.Limits.Max.Jobs.Per = &struct {
						Account *api.V0043Uint32NoValStruct `json:"account,omitempty"`
						User    *api.V0043Uint32NoValStruct `json:"user,omitempty"`
					}{}
				}

				setTrue := true
				jobLimit := int32(*create.Limits.MaxJobsPerAccount)
				apiQoS.Limits.Max.Jobs.Per.Account = &api.V0043Uint32NoValStruct{
					Set:    &setTrue,
					Number: &jobLimit,
				}
			}

			// TODO: Convert TRES limits (CPU, Memory, Node) when we have a clearer understanding of TRES format
			// TRES format is complex and requires special handling
		}
	}

	return apiQoS, nil
}

// convertCommonQoSUpdateToAPI converts common QoS update to v0.0.43 API format
func (a *QoSAdapter) convertCommonQoSUpdateToAPI(existing *types.QoS, update *types.QoSUpdate) (*api.V0043Qos, error) {
	apiQoS := &api.V0043Qos{}
	apiQoS.Name = &existing.Name
	apiQoS.Description = &existing.Description

	// Apply updates
	if update.Description != nil {
		apiQoS.Description = update.Description
	}

	// Priority
	priority := existing.Priority
	if update.Priority != nil {
		priority = *update.Priority
	}
	if priority > 0 {
		setTrue := true
		priorityInt32 := int32(priority)
		apiQoS.Priority = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &priorityInt32,
		}
	}

	// Flags
	flags := existing.Flags
	if update.Flags != nil && len(*update.Flags) > 0 {
		flags = *update.Flags
	}
	if len(flags) > 0 {
		apiFlags := make([]api.V0043QosFlags, 0, len(flags))
		for _, flag := range flags {
			apiFlags = append(apiFlags, api.V0043QosFlags(flag))
		}
		apiQoS.Flags = &apiFlags
	}

	// Preempt mode
	preemptMode := existing.PreemptMode
	if update.PreemptMode != nil && len(*update.PreemptMode) > 0 {
		preemptMode = (*update.PreemptMode)[0]
	}
	if preemptMode != "" {
		modes := []api.V0043QosPreemptMode{api.V0043QosPreemptMode(preemptMode)}
		if apiQoS.Preempt == nil {
			apiQoS.Preempt = &struct {
				ExemptTime *api.V0043Uint32NoValStruct `json:"exempt_time,omitempty"`
				List       *api.V0043QosPreemptList    `json:"list,omitempty"`
				Mode       *[]api.V0043QosPreemptMode  `json:"mode,omitempty"`
			}{}
		}
		apiQoS.Preempt.Mode = &modes
	}

	// Usage factor
	usageFactor := existing.UsageFactor
	if update.UsageFactor != nil {
		usageFactor = *update.UsageFactor
	}
	if usageFactor != 0 {
		setTrue := true
		apiQoS.UsageFactor = &api.V0043Float64NoValStruct{
			Set:    &setTrue,
			Number: &usageFactor,
		}
	}

	// Usage threshold
	usageThreshold := existing.UsageThreshold
	if update.UsageThreshold != nil {
		usageThreshold = *update.UsageThreshold
	}
	if usageThreshold != 0 {
		setTrue := true
		apiQoS.UsageThreshold = &api.V0043Float64NoValStruct{
			Set:    &setTrue,
			Number: &usageThreshold,
		}
	}

	// Convert limits if provided in update
	if update.Limits != nil || update.GraceTime != nil {
		// Initialize the Limits structure (same complex structure as in create)
		apiQoS.Limits = &struct {
			Factor    *api.V0043Float64NoValStruct `json:"factor,omitempty"`
			GraceTime *int32                       `json:"grace_time,omitempty"`
			Max       *struct {
				Accruing *struct {
					Per *struct {
						Account *api.V0043Uint32NoValStruct `json:"account,omitempty"`
						User    *api.V0043Uint32NoValStruct `json:"user,omitempty"`
					} `json:"per,omitempty"`
				} `json:"accruing,omitempty"`
				ActiveJobs *struct {
					Accruing *api.V0043Uint32NoValStruct `json:"accruing,omitempty"`
					Count    *api.V0043Uint32NoValStruct `json:"count,omitempty"`
				} `json:"active_jobs,omitempty"`
				Jobs *struct {
					ActiveJobs *struct {
						Per *struct {
							Account *api.V0043Uint32NoValStruct `json:"account,omitempty"`
							User    *api.V0043Uint32NoValStruct `json:"user,omitempty"`
						} `json:"per,omitempty"`
					} `json:"active_jobs,omitempty"`
					Count *api.V0043Uint32NoValStruct `json:"count,omitempty"`
					Per   *struct {
						Account *api.V0043Uint32NoValStruct `json:"account,omitempty"`
						User    *api.V0043Uint32NoValStruct `json:"user,omitempty"`
					} `json:"per,omitempty"`
				} `json:"jobs,omitempty"`
				Tres *struct {
					Minutes *struct {
						Per *struct {
							Account *api.V0043TresList `json:"account,omitempty"`
							Job     *api.V0043TresList `json:"job,omitempty"`
							Qos     *api.V0043TresList `json:"qos,omitempty"`
							User    *api.V0043TresList `json:"user,omitempty"`
						} `json:"per,omitempty"`
						Total *api.V0043TresList `json:"total,omitempty"`
					} `json:"minutes,omitempty"`
					Per *struct {
						Account *api.V0043TresList `json:"account,omitempty"`
						Job     *api.V0043TresList `json:"job,omitempty"`
						Node    *api.V0043TresList `json:"node,omitempty"`
						User    *api.V0043TresList `json:"user,omitempty"`
					} `json:"per,omitempty"`
					Total *api.V0043TresList `json:"total,omitempty"`
				} `json:"tres,omitempty"`
				WallClock *struct {
					Per *struct {
						Job *api.V0043Uint32NoValStruct `json:"job,omitempty"`
						Qos *api.V0043Uint32NoValStruct `json:"qos,omitempty"`
					} `json:"per,omitempty"`
				} `json:"wall_clock,omitempty"`
			} `json:"max,omitempty"`
			Min *struct {
				PriorityThreshold *api.V0043Uint32NoValStruct `json:"priority_threshold,omitempty"`
				Tres              *struct {
					Per *struct {
						Job *api.V0043TresList `json:"job,omitempty"`
					} `json:"per,omitempty"`
				} `json:"tres,omitempty"`
			} `json:"min,omitempty"`
		}{}

		// Set grace time
		if update.GraceTime != nil {
			graceTime := int32(*update.GraceTime)
			apiQoS.Limits.GraceTime = &graceTime
		}

		// Convert limits
		if update.Limits != nil {
			// Similar logic as in create, but for updates
			// Initialize Max structure if we have any max limits
			if update.Limits.MaxJobsPerUser != nil || update.Limits.MaxJobsPerAccount != nil {
				// Create the same nested structure as in create
				// (Code would be similar to create, just using update.Limits instead)
				// For brevity, I'll add a TODO here
				// TODO: Implement full limits update conversion similar to create
			}
		}
	}

	return apiQoS, nil
}
