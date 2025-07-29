package v0_0_42

import (
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

// convertAPIQoSToCommon converts a v0.0.42 API QoS to common QoS type
func (a *QoSAdapter) convertAPIQoSToCommon(apiQoS api.V0042Qos) (*types.QoS, error) {
	qos := &types.QoS{}

	// Basic fields
	if apiQoS.Name != nil {
		qos.Name = *apiQoS.Name
	}
	if apiQoS.Id != nil {
		qos.ID = uint32(*apiQoS.Id)
	}
	if apiQoS.Description != nil {
		qos.Description = *apiQoS.Description
	}

	// Priority
	if apiQoS.Priority != nil {
		qos.Priority = uint32(apiQoS.Priority.Number)
	}

	// Usage factor
	if apiQoS.UsageFactor != nil {
		qos.UsageFactor = apiQoS.UsageFactor.Number
	}

	// Flags
	if apiQoS.Flags != nil && len(*apiQoS.Flags) > 0 {
		qos.Flags = *apiQoS.Flags
	}

	// Preemption settings
	if apiQoS.Preempt != nil {
		if apiQoS.Preempt.Mode != nil && len(*apiQoS.Preempt.Mode) > 0 {
			qos.PreemptMode = strings.Join(*apiQoS.Preempt.Mode, ",")
		}
		if apiQoS.Preempt.List != nil && len(*apiQoS.Preempt.List) > 0 {
			qos.PreemptList = *apiQoS.Preempt.List
		}
		if apiQoS.Preempt.ExemptTime != nil {
			qos.PreemptExemptTime = uint32(apiQoS.Preempt.ExemptTime.Number)
		}
	}

	// Limits
	if apiQoS.Limits != nil {
		// Grace time
		if apiQoS.Limits.GraceTime != nil {
			qos.GraceTime = uint32(*apiQoS.Limits.GraceTime)
		}

		// Factor
		if apiQoS.Limits.Factor != nil {
			qos.UsageThreshold = apiQoS.Limits.Factor.Number
		}

		// Max limits
		if apiQoS.Limits.Max != nil {
			// Jobs limits
			if apiQoS.Limits.Max.Jobs != nil {
				if apiQoS.Limits.Max.Jobs.Count != nil {
					qos.MaxJobs = uint32(apiQoS.Limits.Max.Jobs.Count.Number)
				}
				if apiQoS.Limits.Max.Jobs.Per != nil {
					if apiQoS.Limits.Max.Jobs.Per.User != nil {
						qos.MaxJobsPerUser = uint32(apiQoS.Limits.Max.Jobs.Per.User.Number)
					}
					if apiQoS.Limits.Max.Jobs.Per.Account != nil {
						qos.MaxJobsPerAccount = uint32(apiQoS.Limits.Max.Jobs.Per.Account.Number)
					}
				}
			}

			// Active jobs limits
			if apiQoS.Limits.Max.ActiveJobs != nil {
				if apiQoS.Limits.Max.ActiveJobs.Count != nil {
					qos.MaxJobsSubmit = uint32(apiQoS.Limits.Max.ActiveJobs.Count.Number)
				}
			}

			// TRES limits
			if apiQoS.Limits.Max.Tres != nil {
				// Total TRES
				if apiQoS.Limits.Max.Tres.Total != nil {
					qos.MaxTRES = a.convertTRESListToString(apiQoS.Limits.Max.Tres.Total)
				}
				
				// Per-job TRES
				if apiQoS.Limits.Max.Tres.Per != nil && apiQoS.Limits.Max.Tres.Per.Job != nil {
					qos.MaxTRESPerJob = a.convertTRESListToString(apiQoS.Limits.Max.Tres.Per.Job)
				}
				
				// Per-user TRES
				if apiQoS.Limits.Max.Tres.Per != nil && apiQoS.Limits.Max.Tres.Per.User != nil {
					qos.MaxTRESPerUser = a.convertTRESListToString(apiQoS.Limits.Max.Tres.Per.User)
				}
				
				// Per-account TRES
				if apiQoS.Limits.Max.Tres.Per != nil && apiQoS.Limits.Max.Tres.Per.Account != nil {
					qos.MaxTRESPerAccount = a.convertTRESListToString(apiQoS.Limits.Max.Tres.Per.Account)
				}
				
				// Per-node TRES
				if apiQoS.Limits.Max.Tres.Per != nil && apiQoS.Limits.Max.Tres.Per.Node != nil {
					qos.MaxTRESPerNode = a.convertTRESListToString(apiQoS.Limits.Max.Tres.Per.Node)
				}
			}

			// Wall clock limits
			if apiQoS.Limits.Max.WallClock != nil && apiQoS.Limits.Max.WallClock.Per != nil {
				if apiQoS.Limits.Max.WallClock.Per.Job != nil {
					qos.MaxWallDurationPerJob = uint32(apiQoS.Limits.Max.WallClock.Per.Job.Number)
				}
			}
		}

		// Min limits
		if apiQoS.Limits.Min != nil {
			// Priority threshold
			if apiQoS.Limits.Min.PriorityThreshold != nil {
				qos.MinPriorityThreshold = uint32(apiQoS.Limits.Min.PriorityThreshold.Number)
			}
			
			// Min TRES per job
			if apiQoS.Limits.Min.Tres != nil && apiQoS.Limits.Min.Tres.Per != nil && apiQoS.Limits.Min.Tres.Per.Job != nil {
				qos.MinTRESPerJob = a.convertTRESListToString(apiQoS.Limits.Min.Tres.Per.Job)
			}
		}
	}

	return qos, nil
}

// convertTRESListToString converts TRES list to string representation
func (a *QoSAdapter) convertTRESListToString(tresList *api.V0042TresList) string {
	if tresList == nil || len(*tresList) == 0 {
		return ""
	}

	tresStrs := make([]string, 0, len(*tresList))
	for _, tres := range *tresList {
		if tres.Type != nil && tres.Count != nil {
			tresStrs = append(tresStrs, *tres.Type+"="+string(*tres.Count))
		}
	}

	return strings.Join(tresStrs, ",")
}

// convertCommonQoSCreateToAPI converts common QoS create request to v0.0.42 API format
func (a *QoSAdapter) convertCommonQoSCreateToAPI(req *types.QoSCreateRequest) (*api.SlurmdbV0042PostQosJSONRequestBody, error) {
	apiReq := &api.SlurmdbV0042PostQosJSONRequestBody{
		Qos: &[]api.V0042Qos{
			{
				Name:        &req.Name,
				Description: &req.Description,
			},
		},
	}

	qos := &(*apiReq.Qos)[0]

	// Priority
	if req.Priority > 0 {
		priority := api.V0042Uint32NoValStruct{
			Set:    true,
			Number: uint64(req.Priority),
		}
		qos.Priority = &priority
	}

	// Usage factor
	if req.UsageFactor > 0 {
		usageFactor := api.V0042Float64NoValStruct{
			Set:    true,
			Number: req.UsageFactor,
		}
		qos.UsageFactor = &usageFactor
	}

	// Flags
	if len(req.Flags) > 0 {
		flags := api.V0042QosFlags(req.Flags)
		qos.Flags = &flags
	}

	// Preemption
	if req.PreemptMode != "" || len(req.PreemptList) > 0 || req.PreemptExemptTime > 0 {
		qos.Preempt = &struct {
			ExemptTime *api.V0042Uint32NoValStruct `json:"exempt_time,omitempty"`
			List       *api.V0042QosPreemptList    `json:"list,omitempty"`
			Mode       *api.V0042QosPreemptModes   `json:"mode,omitempty"`
		}{}

		if req.PreemptMode != "" {
			modes := strings.Split(req.PreemptMode, ",")
			preemptModes := api.V0042QosPreemptModes(modes)
			qos.Preempt.Mode = &preemptModes
		}

		if len(req.PreemptList) > 0 {
			preemptList := api.V0042QosPreemptList(req.PreemptList)
			qos.Preempt.List = &preemptList
		}

		if req.PreemptExemptTime > 0 {
			exemptTime := api.V0042Uint32NoValStruct{
				Set:    true,
				Number: uint64(req.PreemptExemptTime),
			}
			qos.Preempt.ExemptTime = &exemptTime
		}
	}

	// Initialize limits structure if we have any limits
	if req.GraceTime > 0 || req.MaxJobs > 0 || req.MaxJobsPerUser > 0 || req.MaxJobsPerAccount > 0 ||
		req.MaxTRES != "" || req.MaxTRESPerJob != "" || req.MaxTRESPerUser != "" ||
		req.MaxWallDurationPerJob > 0 || req.MinPriorityThreshold > 0 {
		
		qos.Limits = &struct {
			Factor    *api.V0042Float64NoValStruct `json:"factor,omitempty"`
			GraceTime *int32                       `json:"grace_time,omitempty"`
			Max       *struct {
				Accruing   *struct{} `json:"accruing,omitempty"`
				ActiveJobs *struct {
					Accruing *api.V0042Uint32NoValStruct `json:"accruing,omitempty"`
					Count    *api.V0042Uint32NoValStruct `json:"count,omitempty"`
				} `json:"active_jobs,omitempty"`
				Jobs *struct {
					ActiveJobs *struct{} `json:"active_jobs,omitempty"`
					Count      *api.V0042Uint32NoValStruct `json:"count,omitempty"`
					Per        *struct {
						Account *api.V0042Uint32NoValStruct `json:"account,omitempty"`
						User    *api.V0042Uint32NoValStruct `json:"user,omitempty"`
					} `json:"per,omitempty"`
				} `json:"jobs,omitempty"`
				Tres      *struct{} `json:"tres,omitempty"`
				WallClock *struct {
					Per *struct {
						Job *api.V0042Uint32NoValStruct `json:"job,omitempty"`
						Qos *api.V0042Uint32NoValStruct `json:"qos,omitempty"`
					} `json:"per,omitempty"`
				} `json:"wall_clock,omitempty"`
			} `json:"max,omitempty"`
			Min *struct {
				PriorityThreshold *api.V0042Uint32NoValStruct `json:"priority_threshold,omitempty"`
				Tres              *struct{}                   `json:"tres,omitempty"`
			} `json:"min,omitempty"`
		}{}
	}

	// Set individual limits
	if qos.Limits != nil {
		// Grace time
		if req.GraceTime > 0 {
			graceTime := int32(req.GraceTime)
			qos.Limits.GraceTime = &graceTime
		}

		// Job limits
		if req.MaxJobs > 0 || req.MaxJobsPerUser > 0 || req.MaxJobsPerAccount > 0 {
			qos.Limits.Max = &struct {
				Accruing   *struct{} `json:"accruing,omitempty"`
				ActiveJobs *struct {
					Accruing *api.V0042Uint32NoValStruct `json:"accruing,omitempty"`
					Count    *api.V0042Uint32NoValStruct `json:"count,omitempty"`
				} `json:"active_jobs,omitempty"`
				Jobs *struct {
					ActiveJobs *struct{} `json:"active_jobs,omitempty"`
					Count      *api.V0042Uint32NoValStruct `json:"count,omitempty"`
					Per        *struct {
						Account *api.V0042Uint32NoValStruct `json:"account,omitempty"`
						User    *api.V0042Uint32NoValStruct `json:"user,omitempty"`
					} `json:"per,omitempty"`
				} `json:"jobs,omitempty"`
				Tres      *struct{} `json:"tres,omitempty"`
				WallClock *struct {
					Per *struct {
						Job *api.V0042Uint32NoValStruct `json:"job,omitempty"`
						Qos *api.V0042Uint32NoValStruct `json:"qos,omitempty"`
					} `json:"per,omitempty"`
				} `json:"wall_clock,omitempty"`
			}{}

			if req.MaxJobs > 0 {
				qos.Limits.Max.Jobs = &struct {
					ActiveJobs *struct{} `json:"active_jobs,omitempty"`
					Count      *api.V0042Uint32NoValStruct `json:"count,omitempty"`
					Per        *struct {
						Account *api.V0042Uint32NoValStruct `json:"account,omitempty"`
						User    *api.V0042Uint32NoValStruct `json:"user,omitempty"`
					} `json:"per,omitempty"`
				}{
					Count: &api.V0042Uint32NoValStruct{
						Set:    true,
						Number: uint64(req.MaxJobs),
					},
				}
			}
		}
	}

	return apiReq, nil
}