package v0_0_43

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQoSManagerImpl_List(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/slurmdb/v0.0.43/qos", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			// Return mock response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"meta": {
					"slurm": {
						"version": {
							"major": 25,
							"micro": 5,
							"minor": 0
						},
						"release": "25.05.0"
					}
				},
				"qos": [
					{
						"name": "normal",
						"description": "Normal priority QoS",
						"priority": {
							"set": true,
							"infinite": false,
							"number": 100
						},
						"flags": ["PARTITION_QOS"],
						"preempt": {
							"mode": "DISABLED",
							"list": []
						},
						"limits": {
							"max_jobs": {
								"total": {
									"set": true,
									"infinite": false,
									"number": 1000
								}
							},
							"max_submit_jobs": {
								"total": {
									"set": true,
									"infinite": false,
									"number": 2000
								}
							},
							"max_wall_clock": {
								"per": {
									"job": {
										"set": true,
										"infinite": false,
										"number": 86400
									}
								}
							},
							"max_nodes": {
								"per": {
									"job": {
										"set": true,
										"infinite": false,
										"number": 100
									}
								}
							},
							"max_cpus": {
								"per": {
									"job": {
										"set": true,
										"infinite": false,
										"number": 1000
									}
								}
							},
							"tres": {
								"total": [
									{
										"type": "cpu",
										"count": 10000
									},
									{
										"type": "mem",
										"count": 1048576
									}
								]
							}
						},
						"usage_threshold": {
							"set": true,
							"infinite": false,
							"number": 0.8
						},
						"grace_time": 600,
						"usage_factor": {
							"set": true,
							"infinite": false,
							"number": 1.0
						}
					},
					{
						"name": "high",
						"description": "High priority QoS",
						"priority": {
							"set": true,
							"infinite": false,
							"number": 200
						},
						"flags": ["REQUIRED_RESERVATION"],
						"preempt": {
							"mode": "CANCEL",
							"list": ["normal", "low"]
						}
					}
				]
			}`))
		}))
		defer server.Close()

		// Create client
		client, err := createTestClient(server.URL)
		require.NoError(t, err)

		// Test List
		mgr := NewQoSManagerImpl(client)
		result, err := mgr.List(context.Background(), nil)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, result.Total)
		assert.Len(t, result.QoS, 2)

		// Verify first QoS data
		qos := result.QoS[0]
		assert.Equal(t, "normal", qos.Name)
		assert.Equal(t, "Normal priority QoS", qos.Description)
		assert.Equal(t, 100, qos.Priority)
		assert.Equal(t, []string{"PARTITION_QOS"}, qos.Flags)
		assert.Equal(t, "DISABLED", qos.PreemptMode)
		
		// Check limits
		require.NotNil(t, qos.Limits)
		assert.Equal(t, 1000, *qos.Limits.MaxJobsTotal)
		assert.Equal(t, 2000, *qos.Limits.MaxSubmitJobs)
		assert.Equal(t, 86400, *qos.Limits.MaxWallClockPerJob)
		assert.Equal(t, 100, *qos.Limits.MaxNodesPerJob)
		assert.Equal(t, 1000, *qos.Limits.MaxCPUsPerJob)
		assert.Contains(t, qos.Limits.MaxTRESTotal, "cpu")
		assert.Equal(t, int64(10000), qos.Limits.MaxTRESTotal["cpu"])
		assert.Contains(t, qos.Limits.MaxTRESTotal, "mem")
		assert.Equal(t, int64(1048576), qos.Limits.MaxTRESTotal["mem"])

		// Check other fields
		assert.Equal(t, 0.8, *qos.UsageThreshold)
		assert.Equal(t, 600, *qos.GraceTime)
		assert.Equal(t, 1.0, *qos.UsageFactor)

		// Verify second QoS
		qos2 := result.QoS[1]
		assert.Equal(t, "high", qos2.Name)
		assert.Equal(t, "High priority QoS", qos2.Description)
		assert.Equal(t, 200, qos2.Priority)
		assert.Equal(t, []string{"REQUIRED_RESERVATION"}, qos2.Flags)
		assert.Equal(t, "CANCEL", qos2.PreemptMode)
		assert.Equal(t, []string{"normal", "low"}, qos2.PreemptableQoS)
	})

	t.Run("list with filtering", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"qos": [
					{
						"name": "normal",
						"preempt": {
							"mode": "DISABLED"
						}
					},
					{
						"name": "high",
						"preempt": {
							"mode": "CANCEL"
						}
					},
					{
						"name": "urgent",
						"preempt": {
							"mode": "CANCEL"
						}
					}
				]
			}`))
		}))
		defer server.Close()

		// Create client
		client, err := createTestClient(server.URL)
		require.NoError(t, err)

		// Test with filtering
		mgr := NewQoSManagerImpl(client)
		opts := &interfaces.ListQoSOptions{
			PreemptMode: "CANCEL",
		}
		result, err := mgr.List(context.Background(), opts)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, result.Total) // Only high and urgent match
		assert.Equal(t, "high", result.QoS[0].Name)
		assert.Equal(t, "urgent", result.QoS[1].Name)
	})
}

func TestQoSManagerImpl_Get(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/slurmdb/v0.0.43/qos/normal", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			// Return mock response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"qos": [
					{
						"name": "normal",
						"description": "Normal priority QoS",
						"priority": {
							"set": true,
							"number": 100
						}
					}
				]
			}`))
		}))
		defer server.Close()

		// Create client
		client, err := createTestClient(server.URL)
		require.NoError(t, err)

		// Test Get
		mgr := NewQoSManagerImpl(client)
		result, err := mgr.Get(context.Background(), "normal")
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "normal", result.Name)
		assert.Equal(t, "Normal priority QoS", result.Description)
		assert.Equal(t, 100, result.Priority)
	})

	t.Run("not found", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"qos": []
			}`))
		}))
		defer server.Close()

		// Create client
		client, err := createTestClient(server.URL)
		require.NoError(t, err)

		// Test not found
		mgr := NewQoSManagerImpl(client)
		result, err := mgr.Get(context.Background(), "nonexistent")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestQoSManagerImpl_Create(t *testing.T) {
	t.Run("successful create", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/slurmdb/v0.0.43/qos", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			// Return success response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{}`))
		}))
		defer server.Close()

		// Create client
		client, err := createTestClient(server.URL)
		require.NoError(t, err)

		// Test Create
		mgr := NewQoSManagerImpl(client)
		priority := 150
		maxJobs := 500
		create := &interfaces.QoSCreate{
			Name:        "new-qos",
			Description: "New QoS for testing",
			Priority:    &priority,
			Flags:       []string{"PARTITION_QOS"},
			PreemptMode: "DISABLED",
			Limits: &interfaces.QoSLimits{
				MaxJobsTotal: &maxJobs,
			},
		}
		result, err := mgr.Create(context.Background(), create)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.Created)
		assert.Equal(t, "new-qos", result.Name)
	})

	t.Run("validation errors", func(t *testing.T) {
		client, err := createTestClient("http://localhost")
		require.NoError(t, err)
		mgr := NewQoSManagerImpl(client)

		// Test missing name
		create := &interfaces.QoSCreate{
			Description: "Test QoS",
		}
		result, err := mgr.Create(context.Background(), create)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "name is required")
	})
}

func TestQoSManagerImpl_Update(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/slurmdb/v0.0.43/qos", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			// Return success response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		}))
		defer server.Close()

		// Create client
		client, err := createTestClient(server.URL)
		require.NoError(t, err)

		// Test Update
		mgr := NewQoSManagerImpl(client)
		newPriority := 250
		update := &interfaces.QoSUpdate{
			Priority: &newPriority,
		}
		err = mgr.Update(context.Background(), "normal", update)
		require.NoError(t, err)
	})
}

func TestQoSManagerImpl_Delete(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/slurmdb/v0.0.43/qos/test-qos", r.URL.Path)
			assert.Equal(t, "DELETE", r.Method)

			// Return success response
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		// Create client
		client, err := createTestClient(server.URL)
		require.NoError(t, err)

		// Test Delete
		mgr := NewQoSManagerImpl(client)
		err = mgr.Delete(context.Background(), "test-qos")
		require.NoError(t, err)
	})
}