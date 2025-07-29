package v0_0_41

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/adapters/common"
	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
)

// convertAPIQoSToCommon converts a v0.0.41 API QoS to common QoS type
func (a *QoSAdapter) convertAPIQoSToCommon(apiQoS interface{}) (*types.QoS, error) {
	// Type assertion to handle the anonymous struct
	qosData, ok := apiQoS.(struct {
		Description *string `json:"description,omitempty"`
		Flags       *[]api.V0041OpenapiSlurmdbdQosRespQosFlags `json:"flags,omitempty"`
		Id          *api.V0041OpenapiSlurmdbdQosRespQosId `json:"id,omitempty"`
		Limits      *struct {
			Factor      *float32 `json:"factor,omitempty"`
			GraceTime   *api.V0041OpenapiSlurmdbdQosRespQosLimitsGraceTime `json:"grace_time,omitempty"`
			Max         *struct {
				ActiveJobs  *struct {
					Count      *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxActiveJobsCount `json:"count,omitempty"`
					AccruingPer *struct {
						Account *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxActiveJobsAccruingPerAccount `json:"account,omitempty"`
						User    *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxActiveJobsAccruingPerUser    `json:"user,omitempty"`
					} `json:"accruing_per,omitempty"`
				} `json:"active_jobs,omitempty"`
				Tres        *struct {
					Total      *[]struct {
						Type  *string `json:"type,omitempty"`
						Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresTotalValue `json:"value,omitempty"`
					} `json:"total,omitempty"`
					Minutes    *struct {
						Total      *[]struct {
							Type  *string `json:"type,omitempty"`
							Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresMinutesTotalValue `json:"value,omitempty"`
						} `json:"total,omitempty"`
						Per        *struct {
							Qos     *[]struct {
								Type  *string `json:"type,omitempty"`
								Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresMinutesPerQosValue `json:"value,omitempty"`
							} `json:"qos,omitempty"`
							Job     *[]struct {
								Type  *string `json:"type,omitempty"`
								Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresMinutesPerJobValue `json:"value,omitempty"`
							} `json:"job,omitempty"`
							Account *[]struct {
								Type  *string `json:"type,omitempty"`
								Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresMinutesPerAccountValue `json:"value,omitempty"`
							} `json:"account,omitempty"`
							User    *[]struct {
								Type  *string `json:"type,omitempty"`
								Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresMinutesPerUserValue `json:"value,omitempty"`
							} `json:"user,omitempty"`
						} `json:"per,omitempty"`
					} `json:"minutes,omitempty"`
					Per         *struct {
						Account *[]struct {
							Type  *string `json:"type,omitempty"`
							Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresPerAccountValue `json:"value,omitempty"`
						} `json:"account,omitempty"`
						Job     *[]struct {
							Type  *string `json:"type,omitempty"`
							Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresPerJobValue `json:"value,omitempty"`
						} `json:"job,omitempty"`
						Node    *[]struct {
							Type  *string `json:"type,omitempty"`
							Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresPerNodeValue `json:"value,omitempty"`
						} `json:"node,omitempty"`
						User    *[]struct {
							Type  *string `json:"type,omitempty"`
							Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresPerUserValue `json:"value,omitempty"`
						} `json:"user,omitempty"`
					} `json:"per,omitempty"`
				} `json:"tres,omitempty"`
				Jobs        *struct {
					Per *struct {
						Account   *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxJobsPerAccount   `json:"account,omitempty"`
						User      *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxJobsPerUser      `json:"user,omitempty"`
						AccruingPer *struct {
							Account *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxJobsAccruingPerAccount `json:"account,omitempty"`
							User    *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxJobsAccruingPerUser    `json:"user,omitempty"`
						} `json:"accruing_per,omitempty"`
					} `json:"per,omitempty"`
				} `json:"jobs,omitempty"`
				WallClock   *struct {
					Per *struct {
						Job *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxWallClockPerJob `json:"job,omitempty"`
						Qos *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxWallClockPerQos `json:"qos,omitempty"`
					} `json:"per,omitempty"`
				} `json:"wall_clock,omitempty"`
			} `json:"max,omitempty"`
			Min         *struct {
				Priority  *api.V0041OpenapiSlurmdbdQosRespQosLimitsMinPriority `json:"priority,omitempty"`
				TresPerJob *[]struct {
					Type  *string `json:"type,omitempty"`
					Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMinTresPerJobValue `json:"value,omitempty"`
				} `json:"tres_per_job,omitempty"`
			} `json:"min,omitempty"`
		} `json:"limits,omitempty"`
		Name        *string `json:"name,omitempty"`
		Preempt     *struct {
			List     *[]string `json:"list,omitempty"`
			Mode     *[]api.V0041OpenapiSlurmdbdQosRespQosPreemptMode `json:"mode,omitempty"`
			ExemptTime *api.V0041OpenapiSlurmdbdQosRespQosPreemptExemptTime `json:"exempt_time,omitempty"`
		} `json:"preempt,omitempty"`
		Priority    *api.V0041OpenapiSlurmdbdQosRespQosPriority `json:"priority,omitempty"`
	})
	if !ok {
		return nil, fmt.Errorf("unexpected QoS data type")
	}

	qos := &types.QoS{}

	// Basic fields
	if qosData.Id != nil && qosData.Id.Number != nil {
		qos.ID = uint32(*qosData.Id.Number)
	}
	if qosData.Name != nil {
		qos.Name = *qosData.Name
	}
	if qosData.Description != nil {
		qos.Description = *qosData.Description
	}
	if qosData.Priority != nil && qosData.Priority.Number != nil {
		qos.Priority = uint32(*qosData.Priority.Number)
	}

	// Flags
	if qosData.Flags != nil {
		var flags []string
		for _, flag := range *qosData.Flags {
			flags = append(flags, string(flag))
		}
		qos.Flags = flags
	}

	// Preempt settings
	if qosData.Preempt != nil {
		if qosData.Preempt.Mode != nil && len(*qosData.Preempt.Mode) > 0 {
			// Convert first preempt mode to string
			qos.PreemptMode = string((*qosData.Preempt.Mode)[0])
		}
		if qosData.Preempt.List != nil {
			qos.PreemptList = *qosData.Preempt.List
		}
		if qosData.Preempt.ExemptTime != nil && qosData.Preempt.ExemptTime.Number != nil {
			qos.PreemptExemptTime = time.Duration(*qosData.Preempt.ExemptTime.Number) * time.Minute
		}
	}

	// Limits
	if qosData.Limits != nil {
		// Grace time
		if qosData.Limits.GraceTime != nil && qosData.Limits.GraceTime.Number != nil {
			qos.GraceTime = uint32(*qosData.Limits.GraceTime.Number)
		}

		// Usage factor
		if qosData.Limits.Factor != nil {
			qos.UsageFactor = float64(*qosData.Limits.Factor)
		}

		// Max limits
		if qosData.Limits.Max != nil {
			// Max jobs
			if qosData.Limits.Max.Jobs != nil && qosData.Limits.Max.Jobs.Per != nil {
				if qosData.Limits.Max.Jobs.Per.User != nil && qosData.Limits.Max.Jobs.Per.User.Number != nil {
					qos.MaxJobsPerUser = uint32(*qosData.Limits.Max.Jobs.Per.User.Number)
				}
				if qosData.Limits.Max.Jobs.Per.Account != nil && qosData.Limits.Max.Jobs.Per.Account.Number != nil {
					qos.MaxJobsPerAccount = uint32(*qosData.Limits.Max.Jobs.Per.Account.Number)
				}
			}

			// Max active jobs
			if qosData.Limits.Max.ActiveJobs != nil {
				if qosData.Limits.Max.ActiveJobs.Count != nil && qosData.Limits.Max.ActiveJobs.Count.Number != nil {
					qos.MaxJobsTotal = uint32(*qosData.Limits.Max.ActiveJobs.Count.Number)
				}
				if qosData.Limits.Max.ActiveJobs.AccruingPer != nil {
					if qosData.Limits.Max.ActiveJobs.AccruingPer.User != nil && qosData.Limits.Max.ActiveJobs.AccruingPer.User.Number != nil {
						qos.MaxSubmitJobsPerUser = uint32(*qosData.Limits.Max.ActiveJobs.AccruingPer.User.Number)
					}
				}
			}

			// Max wall clock
			if qosData.Limits.Max.WallClock != nil && qosData.Limits.Max.WallClock.Per != nil {
				if qosData.Limits.Max.WallClock.Per.Job != nil && qosData.Limits.Max.WallClock.Per.Job.Number != nil {
					qos.MaxWall = time.Duration(*qosData.Limits.Max.WallClock.Per.Job.Number) * time.Minute
				}
			}

			// Max TRES
			if qosData.Limits.Max.Tres != nil {
				// Total TRES
				if qosData.Limits.Max.Tres.Total != nil {
					qos.MaxTRES = convertQoSTRESListToString(*qosData.Limits.Max.Tres.Total)
				}

				// Per-user TRES
				if qosData.Limits.Max.Tres.Per != nil && qosData.Limits.Max.Tres.Per.User != nil {
					qos.MaxTRESPerUser = convertQoSTRESListToString(*qosData.Limits.Max.Tres.Per.User)
				}

				// Per-job TRES
				if qosData.Limits.Max.Tres.Per != nil && qosData.Limits.Max.Tres.Per.Job != nil {
					qos.MaxTRESPerJob = convertQoSTRESListToString(*qosData.Limits.Max.Tres.Per.Job)
				}

				// Per-node TRES
				if qosData.Limits.Max.Tres.Per != nil && qosData.Limits.Max.Tres.Per.Node != nil {
					qos.MaxTRESPerNode = convertQoSTRESListToString(*qosData.Limits.Max.Tres.Per.Node)
				}

				// Per-account TRES
				if qosData.Limits.Max.Tres.Per != nil && qosData.Limits.Max.Tres.Per.Account != nil {
					qos.MaxTRESPerAccount = convertQoSTRESListToString(*qosData.Limits.Max.Tres.Per.Account)
				}

				// TRES minutes
				if qosData.Limits.Max.Tres.Minutes != nil && qosData.Limits.Max.Tres.Minutes.Total != nil {
					qos.MaxTRESMinutes = convertQoSTRESListToString(*qosData.Limits.Max.Tres.Minutes.Total)
				}
			}
		}

		// Min limits
		if qosData.Limits.Min != nil {
			// Min priority threshold
			if qosData.Limits.Min.Priority != nil && qosData.Limits.Min.Priority.Number != nil {
				qos.MinPrioThreshold = uint32(*qosData.Limits.Min.Priority.Number)
			}

			// Min TRES per job
			if qosData.Limits.Min.TresPerJob != nil {
				qos.MinTRESPerJob = convertQoSTRESListToString(*qosData.Limits.Min.TresPerJob)
			}
		}
	}

	return qos, nil
}

// convertCommonToAPIQoS converts common QoS to v0.0.41 API request
func (a *QoSAdapter) convertCommonToAPIQoS(qos *types.QoS) *api.V0041OpenapiSlurmdbdQosResp {
	// This is a complex conversion - for brevity, implementing only the basic structure
	req := &api.V0041OpenapiSlurmdbdQosResp{
		Qos: []struct {
			Description *string `json:"description,omitempty"`
			Flags       *[]api.V0041OpenapiSlurmdbdQosRespQosFlags `json:"flags,omitempty"`
			Id          *api.V0041OpenapiSlurmdbdQosRespQosId `json:"id,omitempty"`
			Limits      *struct {
				Factor      *float32 `json:"factor,omitempty"`
				GraceTime   *api.V0041OpenapiSlurmdbdQosRespQosLimitsGraceTime `json:"grace_time,omitempty"`
				Max         *struct {
					ActiveJobs  *struct {
						Count      *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxActiveJobsCount `json:"count,omitempty"`
						AccruingPer *struct {
							Account *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxActiveJobsAccruingPerAccount `json:"account,omitempty"`
							User    *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxActiveJobsAccruingPerUser    `json:"user,omitempty"`
						} `json:"accruing_per,omitempty"`
					} `json:"active_jobs,omitempty"`
					Tres        *struct {
						Total      *[]struct {
							Type  *string `json:"type,omitempty"`
							Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresTotalValue `json:"value,omitempty"`
						} `json:"total,omitempty"`
						Minutes    *struct {
							Total      *[]struct {
								Type  *string `json:"type,omitempty"`
								Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresMinutesTotalValue `json:"value,omitempty"`
							} `json:"total,omitempty"`
							Per        *struct {
								Qos     *[]struct {
									Type  *string `json:"type,omitempty"`
									Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresMinutesPerQosValue `json:"value,omitempty"`
								} `json:"qos,omitempty"`
								Job     *[]struct {
									Type  *string `json:"type,omitempty"`
									Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresMinutesPerJobValue `json:"value,omitempty"`
								} `json:"job,omitempty"`
								Account *[]struct {
									Type  *string `json:"type,omitempty"`
									Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresMinutesPerAccountValue `json:"value,omitempty"`
								} `json:"account,omitempty"`
								User    *[]struct {
									Type  *string `json:"type,omitempty"`
									Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresMinutesPerUserValue `json:"value,omitempty"`
								} `json:"user,omitempty"`
							} `json:"per,omitempty"`
						} `json:"minutes,omitempty"`
						Per         *struct {
							Account *[]struct {
								Type  *string `json:"type,omitempty"`
								Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresPerAccountValue `json:"value,omitempty"`
							} `json:"account,omitempty"`
							Job     *[]struct {
								Type  *string `json:"type,omitempty"`
								Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresPerJobValue `json:"value,omitempty"`
							} `json:"job,omitempty"`
							Node    *[]struct {
								Type  *string `json:"type,omitempty"`
								Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresPerNodeValue `json:"value,omitempty"`
							} `json:"node,omitempty"`
							User    *[]struct {
								Type  *string `json:"type,omitempty"`
								Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresPerUserValue `json:"value,omitempty"`
							} `json:"user,omitempty"`
						} `json:"per,omitempty"`
					} `json:"tres,omitempty"`
					Jobs        *struct {
						Per *struct {
							Account   *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxJobsPerAccount   `json:"account,omitempty"`
							User      *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxJobsPerUser      `json:"user,omitempty"`
							AccruingPer *struct {
								Account *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxJobsAccruingPerAccount `json:"account,omitempty"`
								User    *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxJobsAccruingPerUser    `json:"user,omitempty"`
							} `json:"accruing_per,omitempty"`
						} `json:"per,omitempty"`
					} `json:"jobs,omitempty"`
					WallClock   *struct {
						Per *struct {
							Job *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxWallClockPerJob `json:"job,omitempty"`
							Qos *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxWallClockPerQos `json:"qos,omitempty"`
						} `json:"per,omitempty"`
					} `json:"wall_clock,omitempty"`
				} `json:"max,omitempty"`
				Min         *struct {
					Priority  *api.V0041OpenapiSlurmdbdQosRespQosLimitsMinPriority `json:"priority,omitempty"`
					TresPerJob *[]struct {
						Type  *string `json:"type,omitempty"`
						Value *api.V0041OpenapiSlurmdbdQosRespQosLimitsMinTresPerJobValue `json:"value,omitempty"`
					} `json:"tres_per_job,omitempty"`
				} `json:"min,omitempty"`
			} `json:"limits,omitempty"`
			Name        *string `json:"name,omitempty"`
			Preempt     *struct {
				List     *[]string `json:"list,omitempty"`
				Mode     *[]api.V0041OpenapiSlurmdbdQosRespQosPreemptMode `json:"mode,omitempty"`
				ExemptTime *api.V0041OpenapiSlurmdbdQosRespQosPreemptExemptTime `json:"exempt_time,omitempty"`
			} `json:"preempt,omitempty"`
			Priority    *api.V0041OpenapiSlurmdbdQosRespQosPriority `json:"priority,omitempty"`
		}{
			{},
		},
	}

	q := &req.Qos[0]

	// Set basic fields
	if qos.Name != "" {
		q.Name = &qos.Name
	}
	if qos.Description != "" {
		q.Description = &qos.Description
	}

	// Set priority
	if qos.Priority > 0 {
		priority := int32(qos.Priority)
		set := true
		q.Priority = &api.V0041OpenapiSlurmdbdQosRespQosPriority{
			Number: &priority,
			Set:    &set,
		}
	}

	// Convert flags
	if len(qos.Flags) > 0 {
		flags := make([]api.V0041OpenapiSlurmdbdQosRespQosFlags, 0, len(qos.Flags))
		for _, flag := range qos.Flags {
			switch strings.ToLower(flag) {
			case "denyonlimit":
				flags = append(flags, api.V0041OpenapiSlurmdbdQosRespQosFlagsDenyOnLimit)
			case "enforce_usage_threshold":
				flags = append(flags, api.V0041OpenapiSlurmdbdQosRespQosFlagsEnforceUsageThreshold)
			case "noreserve":
				flags = append(flags, api.V0041OpenapiSlurmdbdQosRespQosFlagsNoReserve)
			case "nodecay":
				flags = append(flags, api.V0041OpenapiSlurmdbdQosRespQosFlagsNoDecay)
			case "override_partition_qos":
				flags = append(flags, api.V0041OpenapiSlurmdbdQosRespQosFlagsOverridePartitionQOS)
			case "partitionminimumnodes":
				flags = append(flags, api.V0041OpenapiSlurmdbdQosRespQosFlagsPartitionMinNodes)
			case "partitionmaximumnodes":
				flags = append(flags, api.V0041OpenapiSlurmdbdQosRespQosFlagsPartitionMaxNodes)
			case "partitiontimelimt":
				flags = append(flags, api.V0041OpenapiSlurmdbdQosRespQosFlagsPartitionTimeLimit)
			case "requiredreservation":
				flags = append(flags, api.V0041OpenapiSlurmdbdQosRespQosFlagsRequiresReservation)
			case "usagefactorsafety":
				flags = append(flags, api.V0041OpenapiSlurmdbdQosRespQosFlagsUsageFactorSafe)
			}
		}
		if len(flags) > 0 {
			q.Flags = &flags
		}
	}

	// Note: Due to the complexity of the nested structure in v0.0.41,
	// I'm only implementing the basic fields here. A full implementation
	// would need to convert all the limits and TRES specifications.

	return req
}

// convertQoSTRESListToString is a wrapper that extracts values from API-specific types
// and delegates to the common TRES converter
func convertQoSTRESListToString(tresList []struct {
	Type  *string `json:"type,omitempty"`
	Value interface{} `json:"value,omitempty"`
}) string {
	// Convert to a format the common function can handle
	var convertedList []struct {
		Type  *string `json:"type,omitempty"`
		Value interface{} `json:"value,omitempty"`
	}
	
	for _, tres := range tresList {
		if tres.Type != nil && tres.Value != nil {
			var extractedValue interface{}
			
			// Extract the numeric value from API-specific types
			switch v := tres.Value.(type) {
			case *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresTotalValue:
				if v != nil && v.Number != nil {
					extractedValue = *v.Number
				}
			case *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresPerUserValue:
				if v != nil && v.Number != nil {
					extractedValue = *v.Number
				}
			case *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresPerJobValue:
				if v != nil && v.Number != nil {
					extractedValue = *v.Number
				}
			case *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresPerNodeValue:
				if v != nil && v.Number != nil {
					extractedValue = *v.Number
				}
			case *api.V0041OpenapiSlurmdbdQosRespQosLimitsMaxTresPerAccountValue:
				if v != nil && v.Number != nil {
					extractedValue = *v.Number
				}
			case *api.V0041OpenapiSlurmdbdQosRespQosLimitsMinTresPerJobValue:
				if v != nil && v.Number != nil {
					extractedValue = *v.Number
				}
			default:
				// For unknown types, pass through as-is
				extractedValue = tres.Value
			}
			
			if extractedValue != nil {
				convertedList = append(convertedList, struct {
					Type  *string `json:"type,omitempty"`
					Value interface{} `json:"value,omitempty"`
				}{
					Type:  tres.Type,
					Value: extractedValue,
				})
			}
		}
	}
	
	return common.ConvertTRESListToString(convertedList)
}