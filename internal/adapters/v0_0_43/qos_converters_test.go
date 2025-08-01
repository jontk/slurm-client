// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQoSAdapter_ConvertAPIQoSToCommon_WithLimits(t *testing.T) {
	adapter := &QoSAdapter{
		QoSBaseManager: base.NewQoSBaseManager("v0.0.43"),
	}

	tests := []struct {
		name     string
		apiQoS   api.V0043Qos
		expected *types.QoS
	}{
		{
			name: "QoS with grace time",
			apiQoS: api.V0043Qos{
				Name:        &[]string{"test-qos"}[0],
				Description: &[]string{"Test QoS"}[0],
				Limits: &struct {
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
				}{
					GraceTime: int32Ptr(300),
				},
			},
			expected: &types.QoS{
				Name:        "test-qos",
				Description: "Test QoS",
				GraceTime:   300,
			},
		},
		{
			name: "QoS with job limits",
			apiQoS: api.V0043Qos{
				Name: &[]string{"limited-qos"}[0],
				Limits: &struct {
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
				}{
					Max: &struct {
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
					}{
						Jobs: &struct {
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
							}{
								User: &api.V0043Uint32NoValStruct{
									Set:    boolPtr(true),
									Number: int32Ptr(10),
								},
								Account: &api.V0043Uint32NoValStruct{
									Set:    boolPtr(true),
									Number: int32Ptr(50),
								},
							},
						},
					},
				},
			},
			expected: &types.QoS{
				Name: "limited-qos",
				Limits: &types.QoSLimits{
					MaxJobsPerUser:    intPtr(10),
					MaxJobsPerAccount: intPtr(50),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.convertAPIQoSToCommon(tt.apiQoS)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestQoSAdapter_ConvertCommonQoSCreateToAPI_WithLimits(t *testing.T) {
	adapter := &QoSAdapter{
		QoSBaseManager: base.NewQoSBaseManager("v0.0.43"),
	}

	tests := []struct {
		name    string
		create  *types.QoSCreate
		wantErr bool
		check   func(t *testing.T, apiQoS *api.V0043Qos)
	}{
		{
			name: "QoS with grace time",
			create: &types.QoSCreate{
				Name:      "test-qos",
				GraceTime: intPtr(300),
			},
			check: func(t *testing.T, apiQoS *api.V0043Qos) {
				require.NotNil(t, apiQoS.Limits)
				require.NotNil(t, apiQoS.Limits.GraceTime)
				assert.Equal(t, int32(300), *apiQoS.Limits.GraceTime)
			},
		},
		{
			name: "QoS with job limits",
			create: &types.QoSCreate{
				Name: "limited-qos",
				Limits: &types.QoSLimits{
					MaxJobsPerUser:    intPtr(10),
					MaxJobsPerAccount: intPtr(50),
				},
			},
			check: func(t *testing.T, apiQoS *api.V0043Qos) {
				require.NotNil(t, apiQoS.Limits)
				require.NotNil(t, apiQoS.Limits.Max)
				require.NotNil(t, apiQoS.Limits.Max.Jobs)
				require.NotNil(t, apiQoS.Limits.Max.Jobs.Per)
				require.NotNil(t, apiQoS.Limits.Max.Jobs.Per.User)
				assert.Equal(t, int32(10), *apiQoS.Limits.Max.Jobs.Per.User.Number)
				require.NotNil(t, apiQoS.Limits.Max.Jobs.Per.Account)
				assert.Equal(t, int32(50), *apiQoS.Limits.Max.Jobs.Per.Account.Number)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.convertCommonQoSCreateToAPI(tt.create)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tt.check(t, result)
			}
		})
	}
}

// Helper functions
func int32Ptr(i int32) *int32 {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
