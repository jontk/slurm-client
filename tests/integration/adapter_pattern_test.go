package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/internal/interfaces"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
	"github.com/jontk/slurm-client/pkg/auth"
)

// TestAdapterPatternWithMockServer tests the adapter pattern implementation
func TestAdapterPatternWithMockServer(t *testing.T) {
	// Create a mock server that simulates SLURM REST API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/slurm/v0.0.43/slurmdb/qos/":
			handleListQoS(t, w, r)
		case "/slurm/v0.0.43/slurmdb/qos/high-priority":
			handleGetQoS(t, w, r)
		case "/slurm/v0.0.43/slurmdb/qos":
			if r.Method == http.MethodPost {
				handleCreateQoS(t, w, r)
			}
		case "/slurm/v0.0.43/slurmdb/qos/test-qos":
			if r.Method == http.MethodDelete {
				handleDeleteQoS(t, w, r)
			}
		default:
			t.Logf("Unhandled path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create client with our adapter pattern
	ctx := context.Background()
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
		slurm.WithBaseURL(server.URL),
		slurm.WithAuth(auth.NewTokenAuth("test-token")),
	)
	require.NoError(t, err)
	defer client.Close()

	t.Run("List QoS", func(t *testing.T) {
		qosList, err := client.QoS().List(ctx, &interfaces.ListQoSOptions{
			Limit: 10,
		})
		require.NoError(t, err)
		assert.NotNil(t, qosList)
		assert.Len(t, qosList.QoS, 2)
		
		// Check first QoS
		assert.Equal(t, "normal", qosList.QoS[0].Name)
		assert.Equal(t, 100, qosList.QoS[0].Priority)
		assert.Equal(t, 1.0, qosList.QoS[0].UsageFactor)
		
		// Check limits conversion
		assert.Equal(t, 300, qosList.QoS[0].GraceTime)
		assert.Equal(t, 10, qosList.QoS[0].MaxJobsPerUser)
		assert.Equal(t, 50, qosList.QoS[0].MaxJobsPerAccount)
	})

	t.Run("Get QoS", func(t *testing.T) {
		qos, err := client.QoS().Get(ctx, "high-priority")
		require.NoError(t, err)
		assert.NotNil(t, qos)
		assert.Equal(t, "high-priority", qos.Name)
		assert.Equal(t, 1000, qos.Priority)
		assert.Equal(t, "High priority QoS for critical jobs", qos.Description)
	})

	t.Run("Create QoS", func(t *testing.T) {
		createReq := &interfaces.QoSCreate{
			Name:        "test-qos",
			Description: "Test QoS created via adapter",
			Priority:        500,
			Flags:           []string{"DenyOnLimit"},
			UsageFactor:     2.0,
			GraceTime:       600,
			MaxJobsPerUser:  20,
			MaxJobsPerAccount: 100,
		}
		
		resp, err := client.QoS().Create(ctx, createReq)
		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "test-qos", resp.QoSName)
	})

	t.Run("Delete QoS", func(t *testing.T) {
		err := client.QoS().Delete(ctx, "test-qos")
		require.NoError(t, err)
	})
}

func handleListQoS(t *testing.T, w http.ResponseWriter, r *http.Request) {
	response := api.V0043OpenapiSlurmdbdQosResp{
		Qos: []api.V0043Qos{
			{
				Name:        apStringPtr("normal"),
				Description: apStringPtr("Normal priority QoS"),
				Priority: &api.V0043Uint32NoValStruct{
					Set:    apBoolPtr(true),
					Number: apInt32Ptr(100),
				},
				UsageFactor: &api.V0043Float64NoValStruct{
					Set:    apBoolPtr(true),
					Number: apFloat64Ptr(1.0),
				},
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
					GraceTime: apInt32Ptr(300),
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
									Set:    apBoolPtr(true),
									Number: apInt32Ptr(10),
								},
								Account: &api.V0043Uint32NoValStruct{
									Set:    apBoolPtr(true),
									Number: apInt32Ptr(50),
								},
							},
						},
					},
				},
			},
			{
				Name:        apStringPtr("high-priority"),
				Description: apStringPtr("High priority QoS for critical jobs"),
				Priority: &api.V0043Uint32NoValStruct{
					Set:    apBoolPtr(true),
					Number: apInt32Ptr(1000),
				},
			},
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGetQoS(t *testing.T, w http.ResponseWriter, r *http.Request) {
	response := api.V0043OpenapiSlurmdbdQosResp{
		Qos: []api.V0043Qos{
			{
				Name:        apStringPtr("high-priority"),
				Description: apStringPtr("High priority QoS for critical jobs"),
				Priority: &api.V0043Uint32NoValStruct{
					Set:    apBoolPtr(true),
					Number: apInt32Ptr(1000),
				},
			},
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleCreateQoS(t *testing.T, w http.ResponseWriter, r *http.Request) {
	// Just return success with empty errors
	response := api.V0043OpenapiSlurmdbdQosResp{
		Errors: &api.V0043OpenapiErrors{},
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleDeleteQoS(t *testing.T, w http.ResponseWriter, r *http.Request) {
	// Return 204 No Content for successful delete
	w.WriteHeader(http.StatusNoContent)
}

// Helper functions for adapter pattern test
func apStringPtr(s string) *string {
	return &s
}

func apIntPtr(i int) *int {
	return &i
}

func apInt32Ptr(i int32) *int32 {
	return &i
}

func apFloat64Ptr(f float64) *float64 {
	return &f
}

func apBoolPtr(b bool) *bool {
	return &b
}