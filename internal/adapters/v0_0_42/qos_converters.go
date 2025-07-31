package v0_0_42

import (
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

// QoSCreateRequest represents a request to create a QoS
type QoSCreateRequest struct {
	Name              string
	Description       string
	Priority          int
	Flags             []string
	PreemptMode       string  // comma-separated list
	PreemptList       []string
	PreemptExemptTime int
	GraceTime         int
	UsageFactor       float64
	MaxJobs           int
	MaxJobsPerUser    int
	MaxJobsPerAccount int
	MaxTRES           string
	MaxTRESPerJob     string
	MaxTRESPerUser    string
	MaxWallDurationPerJob int
	MinPriorityThreshold int
}

// convertAPIQoSToCommon converts a v0.0.42 API QoS to common QoS type
func (a *QoSAdapter) convertAPIQoSToCommon(apiQoS api.V0042Qos) (*types.QoS, error) {
	qos := &types.QoS{}

	// Basic fields
	if apiQoS.Name != nil {
		qos.Name = *apiQoS.Name
	}
	if apiQoS.Id != nil {
		qos.ID = int32(*apiQoS.Id)
	}
	if apiQoS.Description != nil {
		qos.Description = *apiQoS.Description
	}

	// Priority
	if apiQoS.Priority != nil {
		qos.Priority = int(*apiQoS.Priority.Number)
	}

	// Usage factor
	if apiQoS.UsageFactor != nil {
		qos.UsageFactor = *apiQoS.UsageFactor.Number
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
		// PreemptList field doesn't exist in common QoS type
		// Skip preempt list
		if apiQoS.Preempt.ExemptTime != nil {
			qos.PreemptExemptTime = int(*apiQoS.Preempt.ExemptTime.Number)
		}
	}

	// Limits
	if apiQoS.Limits != nil {
		// Grace time
		if apiQoS.Limits.GraceTime != nil {
			qos.GraceTime = int(*apiQoS.Limits.GraceTime)
		}

		// Factor
		if apiQoS.Limits.Factor != nil {
			qos.UsageThreshold = *apiQoS.Limits.Factor.Number
		}

		// Max limits
		if apiQoS.Limits.Max != nil {
			// Jobs limits
			if apiQoS.Limits.Max.Jobs != nil {
				if apiQoS.Limits.Max.Jobs.Count != nil {
					// MaxJobs field doesn't exist in common QoS type, use Limits
					if qos.Limits == nil {
						qos.Limits = &types.QoSLimits{}
					}
					maxJobs := int(*apiQoS.Limits.Max.Jobs.Count.Number)
					qos.Limits.MaxJobsPerUser = &maxJobs
				}
				if apiQoS.Limits.Max.Jobs.Per != nil {
					if apiQoS.Limits.Max.Jobs.Per.User != nil {
						if qos.Limits == nil {
							qos.Limits = &types.QoSLimits{}
						}
						maxJobs := int(*apiQoS.Limits.Max.Jobs.Per.User.Number)
						qos.Limits.MaxJobsPerUser = &maxJobs
					}
					if apiQoS.Limits.Max.Jobs.Per.Account != nil {
						if qos.Limits == nil {
							qos.Limits = &types.QoSLimits{}
						}
						maxJobs := int(*apiQoS.Limits.Max.Jobs.Per.Account.Number)
						qos.Limits.MaxJobsPerAccount = &maxJobs
					}
				}
			}

			// Active jobs limits
			if apiQoS.Limits.Max.ActiveJobs != nil {
				if apiQoS.Limits.Max.ActiveJobs.Count != nil {
					if qos.Limits == nil {
						qos.Limits = &types.QoSLimits{}
					}
					maxJobsSubmit := int(*apiQoS.Limits.Max.ActiveJobs.Count.Number)
					qos.Limits.MaxSubmitJobsPerUser = &maxJobsSubmit
				}
			}

			// TRES limits
			if apiQoS.Limits.Max.Tres != nil {
				// Total TRES
				if apiQoS.Limits.Max.Tres.Total != nil {
					// Total TRES doesn't have a direct mapping - skip
					_ = apiQoS.Limits.Max.Tres.Total
				}
				
				// Per-job TRES
				if apiQoS.Limits.Max.Tres.Per != nil && apiQoS.Limits.Max.Tres.Per.Job != nil {
					// Per-job TRES doesn't have a direct mapping - skip
					_ = apiQoS.Limits.Max.Tres.Per.Job
				}
				
				// Per-user TRES
				if apiQoS.Limits.Max.Tres.Per != nil && apiQoS.Limits.Max.Tres.Per.User != nil {
					// Per-user TRES doesn't have a direct mapping - skip
					_ = apiQoS.Limits.Max.Tres.Per.User
				}
				
				// Per-account TRES
				if apiQoS.Limits.Max.Tres.Per != nil && apiQoS.Limits.Max.Tres.Per.Account != nil {
					// Per-account TRES doesn't have a direct mapping - skip
					_ = apiQoS.Limits.Max.Tres.Per.Account
				}
				
				// Per-node TRES
				if apiQoS.Limits.Max.Tres.Per != nil && apiQoS.Limits.Max.Tres.Per.Node != nil {
					// Per-node TRES doesn't have a direct mapping - skip
					_ = apiQoS.Limits.Max.Tres.Per.Node
				}
			}

			// Wall clock limits
			if apiQoS.Limits.Max.WallClock != nil && apiQoS.Limits.Max.WallClock.Per != nil {
				if apiQoS.Limits.Max.WallClock.Per.Job != nil {
					if qos.Limits == nil {
						qos.Limits = &types.QoSLimits{}
					}
					maxWall := int(*apiQoS.Limits.Max.WallClock.Per.Job.Number)
					qos.Limits.MaxWallTimePerJob = &maxWall
				}
			}
		}

		// Min limits
		if apiQoS.Limits.Min != nil {
			// Priority threshold
			if apiQoS.Limits.Min.PriorityThreshold != nil {
				// Priority threshold doesn't have a direct mapping - skip
				_ = apiQoS.Limits.Min.PriorityThreshold
			}
			
			// Min TRES per job
			if apiQoS.Limits.Min.Tres != nil && apiQoS.Limits.Min.Tres.Per != nil && apiQoS.Limits.Min.Tres.Per.Job != nil {
				// Min TRES per job doesn't have a direct mapping - skip
				_ = apiQoS.Limits.Min.Tres.Per.Job
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
		if tres.Type != "" && tres.Count != nil {
			tresStrs = append(tresStrs, fmt.Sprintf("%s=%d", tres.Type, *tres.Count))
		}
	}

	return strings.Join(tresStrs, ",")
}

// convertCommonQoSCreateToAPI converts common QoS create request to v0.0.42 API format
func (a *QoSAdapter) convertCommonQoSCreateToAPI(req *QoSCreateRequest) (*api.SlurmdbV0042PostQosJSONRequestBody, error) {
	qosList := api.V0042QosList{
		{
			Name:        &req.Name,
			Description: &req.Description,
		},
	}
	apiReq := &api.SlurmdbV0042PostQosJSONRequestBody{
		Qos: qosList,
	}

	qos := &apiReq.Qos[0]

	// Priority
	if req.Priority > 0 {
		set := true
		priority := int32(req.Priority)
		priorityStruct := api.V0042Uint32NoValStruct{
			Set:    &set,
			Number: &priority,
		}
		qos.Priority = &priorityStruct
	}

	// Usage factor
	if req.UsageFactor > 0 {
		set := true
		usageFactor := req.UsageFactor
		usageFactorStruct := api.V0042Float64NoValStruct{
			Set:    &set,
			Number: &usageFactor,
		}
		qos.UsageFactor = &usageFactorStruct
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
			set := true
			exemptTimeNum := int32(req.PreemptExemptTime)
			exemptTime := api.V0042Uint32NoValStruct{
				Set:    &set,
				Number: &exemptTimeNum,
			}
			qos.Preempt.ExemptTime = &exemptTime
		}
	}

	// Limits are too complex to handle properly in v0.0.42 due to nested anonymous structs
	// Skip limits implementation for now

	return apiReq, nil
}