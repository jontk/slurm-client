package v0_0_40

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// TestGetJobUtilization_v40 tests the GetJobUtilization method for v0.0.40
func TestGetJobUtilization_v40(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_MinimalMetrics", func(t *testing.T) {
		mockJob := &interfaces.Job{
			ID:        "40001",
			Name:      "minimal-job",
			State:     "RUNNING",
			CPUs:      2,
			Memory:    8 * 1024 * 1024 * 1024, // 8GB
			Partition: "compute",
			Nodes:     []string{"node040"},
			SubmitTime: time.Now().Add(-30 * time.Minute),
		}

		mockClient := &mockAPIClient{
			getJobResponse: mockJob,
		}
		manager := &JobManagerImpl{
			client: &WrapperClient{
				apiClient: createMockClientWithResponses(mockClient),
			},
		}

		utilization, err := manager.GetJobUtilization(ctx, "40001")
		require.NoError(t, err)
		require.NotNil(t, utilization)

		// Verify basic fields
		assert.Equal(t, "40001", utilization.JobID)
		assert.Equal(t, "minimal-job", utilization.JobName)

		// Verify CPU utilization (v0.0.40 has fixed 70%)
		assert.NotNil(t, utilization.CPUUtilization)
		assert.Equal(t, 70.0, utilization.CPUUtilization.Percentage)
		assert.Equal(t, "fixed_percentage", utilization.CPUUtilization.Metadata["estimation_method"])
		assert.Equal(t, "low", utilization.CPUUtilization.Metadata["confidence"])

		// Verify memory utilization (v0.0.40 has fixed 60%)
		assert.NotNil(t, utilization.MemoryUtilization)
		assert.Equal(t, 60.0, utilization.MemoryUtilization.Percentage)
		assert.Equal(t, "fixed_percentage", utilization.MemoryUtilization.Metadata["estimation_method"])
		assert.Equal(t, "low", utilization.MemoryUtilization.Metadata["confidence"])

		// Verify no advanced metrics in v0.0.40
		assert.Nil(t, utilization.GPUUtilization)
		assert.Nil(t, utilization.IOUtilization)
		assert.Nil(t, utilization.NetworkUtilization)
		assert.Nil(t, utilization.EnergyUsage)

		// Verify metadata
		assert.Equal(t, "v0.0.40", utilization.Metadata["version"])
		assert.Equal(t, "minimal", utilization.Metadata["feature_level"])
		assert.Equal(t, "estimated", utilization.Metadata["data_quality"])
		
		limitations := utilization.Metadata["limitations"].([]string)
		assert.Contains(t, limitations, "fixed_utilization_percentages")
		assert.Contains(t, limitations, "no_actual_measurements")
		assert.Contains(t, limitations, "no_gpu_support")
		assert.Contains(t, limitations, "no_performance_counters")
	})

	t.Run("ClientNotInitialized", func(t *testing.T) {
		manager := &JobManagerImpl{
			client: &WrapperClient{},
		}

		utilization, err := manager.GetJobUtilization(ctx, "40002")
		assert.Error(t, err)
		assert.Nil(t, utilization)
		assert.True(t, errors.IsClientError(err))
	})
}

// TestGetJobEfficiency_v40 tests the GetJobEfficiency method for v0.0.40
func TestGetJobEfficiency_v40(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_FixedEstimate", func(t *testing.T) {
		mockJob := &interfaces.Job{
			ID:        "40003",
			Name:      "efficiency-test",
			State:     "COMPLETED",
			CPUs:      4,
			Memory:    16 * 1024 * 1024 * 1024,
			Partition: "compute",
			SubmitTime: time.Now().Add(-1 * time.Hour),
			EndTime:    &[]time.Time{time.Now()}[0],
		}

		mockClient := &mockAPIClient{
			getJobResponse: mockJob,
		}
		manager := &JobManagerImpl{
			client: &WrapperClient{
				apiClient: createMockClientWithResponses(mockClient),
			},
		}

		efficiency, err := manager.GetJobEfficiency(ctx, "40003")
		require.NoError(t, err)
		require.NotNil(t, efficiency)

		// v0.0.40 always returns fixed 65% efficiency
		assert.Equal(t, 65.0, efficiency.Percentage)

		// Verify metadata
		assert.Equal(t, "fixed_estimate_v40", efficiency.Metadata["calculation_method"])
		assert.Equal(t, "v0.0.40", efficiency.Metadata["version"])
		assert.Equal(t, "very_low", efficiency.Metadata["confidence"])
		assert.Equal(t, "Efficiency is estimated, not measured in v0.0.40", efficiency.Metadata["note"])
		
		// Fixed values from assumptions
		assert.Equal(t, 70.0, efficiency.Metadata["cpu_efficiency"])
		assert.Equal(t, 60.0, efficiency.Metadata["memory_efficiency"])

		limitations := efficiency.Metadata["limitations"].([]string)
		assert.Contains(t, limitations, "no_actual_measurements")
		assert.Contains(t, limitations, "fixed_efficiency_value")
		assert.Contains(t, limitations, "no_resource_specific_data")
	})
}

// TestGetJobPerformance_v40 tests the GetJobPerformance method for v0.0.40
func TestGetJobPerformance_v40(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_MinimalReport", func(t *testing.T) {
		startTime := time.Now().Add(-45 * time.Minute)
		endTime := time.Now()
		mockJob := &interfaces.Job{
			ID:         "40004",
			Name:       "perf-test-v40",
			State:      "COMPLETED",
			CPUs:       8,
			Memory:     32 * 1024 * 1024 * 1024,
			Partition:  "compute",
			SubmitTime: time.Now().Add(-1 * time.Hour),
			StartTime:  &startTime,
			EndTime:    &endTime,
			ExitCode:   0,
		}

		mockClient := &mockAPIClient{
			getJobResponse: mockJob,
		}
		manager := &JobManagerImpl{
			client: &WrapperClient{
				apiClient: createMockClientWithResponses(mockClient),
			},
		}

		performance, err := manager.GetJobPerformance(ctx, "40004")
		require.NoError(t, err)
		require.NotNil(t, performance)

		// Verify basic fields
		assert.Equal(t, uint32(40004), performance.JobID)
		assert.Equal(t, "perf-test-v40", performance.JobName)

		// Verify minimal features in v0.0.40
		assert.Nil(t, performance.StepMetrics)
		assert.Nil(t, performance.PerformanceTrends)
		assert.Nil(t, performance.Bottlenecks) // No bottleneck detection

		// Verify only system recommendation
		assert.Len(t, performance.Recommendations, 1)
		rec := performance.Recommendations[0]
		assert.Equal(t, "system", rec.Type)
		assert.Equal(t, "medium", rec.Priority)
		assert.Equal(t, "Upgrade for better analytics", rec.Title)
		assert.Equal(t, "v0.0.40", rec.ConfigChanges["current_api_version"])
		assert.Equal(t, "v0.0.41", rec.ConfigChanges["minimum_recommended"])
		assert.Equal(t, "v0.0.42_or_higher", rec.ConfigChanges["optimal_version"])
	})

	t.Run("InvalidJobID", func(t *testing.T) {
		mockJob := &interfaces.Job{
			ID:    "invalid",
			Name:  "test-job",
			State: "RUNNING",
		}

		mockClient := &mockAPIClient{
			getJobResponse: mockJob,
		}
		manager := &JobManagerImpl{
			client: &WrapperClient{
				apiClient: createMockClientWithResponses(mockClient),
			},
		}

		performance, err := manager.GetJobPerformance(ctx, "invalid")
		assert.Error(t, err)
		assert.Nil(t, performance)
		assert.True(t, errors.IsClientError(err))
	})
}

// TestDataQuality_v40 tests that v0.0.40 clearly indicates estimated data
func TestDataQuality_v40(t *testing.T) {
	ctx := context.Background()
	
	mockJob := &interfaces.Job{
		ID:        "40005",
		Name:      "quality-test",
		State:     "RUNNING",
		CPUs:      16,
		Memory:    64 * 1024 * 1024 * 1024,
		Partition: "compute",
		SubmitTime: time.Now().Add(-1 * time.Hour),
	}

	mockClient := &mockAPIClient{
		getJobResponse: mockJob,
	}
	manager := &JobManagerImpl{
		client: &WrapperClient{
			apiClient: createMockClientWithResponses(mockClient),
		},
	}

	// Test all three methods to ensure they indicate estimated data
	t.Run("Utilization", func(t *testing.T) {
		utilization, err := manager.GetJobUtilization(ctx, "40005")
		require.NoError(t, err)
		
		assert.Equal(t, "estimated", utilization.Metadata["data_quality"])
		assert.Equal(t, "basic_accounting", utilization.Metadata["source"])
		assert.Equal(t, "minimal", utilization.Metadata["feature_level"])
	})

	t.Run("Efficiency", func(t *testing.T) {
		efficiency, err := manager.GetJobEfficiency(ctx, "40005")
		require.NoError(t, err)
		
		assert.Equal(t, "very_low", efficiency.Metadata["confidence"])
		assert.Contains(t, efficiency.Metadata["note"].(string), "estimated")
	})

	t.Run("Performance", func(t *testing.T) {
		performance, err := manager.GetJobPerformance(ctx, "40005")
		require.NoError(t, err)
		
		// Should recommend upgrade due to limited analytics
		assert.Len(t, performance.Recommendations, 1)
		assert.Contains(t, performance.Recommendations[0].Description, "minimal analytics")
	})
}

// TestJobManager_GetJobCPUAnalytics tests the GetJobCPUAnalytics method for v0.0.40
func TestJobManager_GetJobCPUAnalytics(t *testing.T) {
	ctx := context.Background()
	jm := &JobManager{}

	tests := []struct {
		name    string
		jobID   string
		wantErr bool
	}{
		{
			name:    "valid job ID",
			jobID:   "12345",
			wantErr: false,
		},
		{
			name:    "empty job ID",
			jobID:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analytics, err := jm.GetJobCPUAnalytics(ctx, tt.jobID)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, analytics)
			} else {
				require.NoError(t, err)
				require.NotNil(t, analytics)
				
				// Verify minimal analytics for v0.0.40
				assert.Equal(t, 4, analytics.AllocatedCores)
				assert.Equal(t, 4, analytics.RequestedCores)
				assert.Equal(t, 2.0, analytics.UsedCores)
				assert.Equal(t, 50.0, analytics.UtilizationPercent)
				assert.Equal(t, 50.0, analytics.EfficiencyPercent)
				assert.Equal(t, 2.0, analytics.IdleCores)
				assert.False(t, analytics.Oversubscribed)
				
				// Verify core-level metrics
				assert.Len(t, analytics.CoreMetrics, 4)
				for i, metric := range analytics.CoreMetrics {
					assert.Equal(t, i, metric.CoreID)
					assert.Equal(t, 50.0, metric.Utilization)
					assert.Equal(t, 2.5, metric.Frequency)
					assert.Equal(t, 60.0, metric.Temperature)
				}
			}
		})
	}
}

// TestJobManager_GetJobMemoryAnalytics tests the GetJobMemoryAnalytics method for v0.0.40
func TestJobManager_GetJobMemoryAnalytics(t *testing.T) {
	ctx := context.Background()
	jm := &JobManager{}

	tests := []struct {
		name    string
		jobID   string
		wantErr bool
	}{
		{
			name:    "valid job ID",
			jobID:   "12345",
			wantErr: false,
		},
		{
			name:    "empty job ID",
			jobID:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analytics, err := jm.GetJobMemoryAnalytics(ctx, tt.jobID)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, analytics)
			} else {
				require.NoError(t, err)
				require.NotNil(t, analytics)
				
				// Verify minimal analytics for v0.0.40
				assert.Equal(t, int64(16*1024*1024*1024), analytics.AllocatedBytes)
				assert.Equal(t, int64(16*1024*1024*1024), analytics.RequestedBytes)
				assert.InDelta(t, 9.6*1024*1024*1024, float64(analytics.UsedBytes), 1024)
				assert.Equal(t, 60.0, analytics.UtilizationPercent)
				// These fields don't exist in v0.0.40, they're part of the extended fields
				assert.Equal(t, int64(1*1024*1024*1024), analytics.CachedMemory)
				assert.Equal(t, int64(0), analytics.PageSwaps) // Different from SwapBytes
				assert.Equal(t, 10000, analytics.PageFaults)
				assert.Equal(t, 100, analytics.MajorPageFaults)
				
				// NUMA metrics should be minimal
				assert.Len(t, analytics.NUMANodes, 1)
				if len(analytics.NUMANodes) > 0 {
					numa := analytics.NUMANodes[0]
					assert.Equal(t, 0, numa.NodeID)
					assert.InDelta(t, 9.6*1024*1024*1024, float64(numa.UsedBytes), 1024)
					assert.Equal(t, 0.0, numa.LocalPercent)
					assert.Equal(t, 0.0, numa.RemotePercent)
				}
				
				// Memory leak detection not in base struct for v0.0.40
				// These fields would be in metadata or extended analytics
			}
		})
	}
}

// TestJobManager_GetJobIOAnalytics tests the GetJobIOAnalytics method for v0.0.40
func TestJobManager_GetJobIOAnalytics(t *testing.T) {
	ctx := context.Background()
	jm := &JobManager{}

	tests := []struct {
		name    string
		jobID   string
		wantErr bool
	}{
		{
			name:    "valid job ID",
			jobID:   "12345",
			wantErr: false,
		},
		{
			name:    "empty job ID",
			jobID:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analytics, err := jm.GetJobIOAnalytics(ctx, tt.jobID)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, analytics)
			} else {
				require.NoError(t, err)
				require.NotNil(t, analytics)
				
				// Verify minimal analytics for v0.0.40
				assert.Equal(t, int64(10*1024*1024*1024), analytics.ReadBytes)
				assert.Equal(t, int64(5*1024*1024*1024), analytics.WriteBytes)
				assert.Equal(t, int64(1000), analytics.ReadOperations)
				assert.Equal(t, int64(500), analytics.WriteOperations)
				assert.Equal(t, 100.0, analytics.ReadBandwidthMBps)
				assert.Equal(t, 50.0, analytics.WriteBandwidthMBps)
				assert.Equal(t, 10.0, analytics.ReadLatencyMs)
				assert.Equal(t, 15.0, analytics.WriteLatencyMs)
				assert.Equal(t, 10.0, analytics.IOWaitPercent)
				
				// Device metrics should be minimal
				assert.Len(t, analytics.DeviceMetrics, 1)
				if len(analytics.DeviceMetrics) > 0 {
					device := analytics.DeviceMetrics[0]
					assert.Equal(t, "/dev/sda", device.DeviceName)
					assert.Equal(t, "local", device.DeviceType)
					assert.Equal(t, 20.0, device.UtilizationPercent)
					assert.Equal(t, 100.0, device.ReadBandwidthMBps)
					assert.Equal(t, 50.0, device.WriteBandwidthMBps)
					assert.Equal(t, int64(1000), device.ReadOps)
					assert.Equal(t, int64(500), device.WriteOps)
					assert.Equal(t, 10.0, device.AvgQueueDepth)
				}
			}
		})
	}
}

// TestJobManager_GetJobComprehensiveAnalytics tests the GetJobComprehensiveAnalytics method for v0.0.40
func TestJobManager_GetJobComprehensiveAnalytics(t *testing.T) {
	ctx := context.Background()
	jm := &JobManager{}

	tests := []struct {
		name    string
		jobID   string
		wantErr bool
	}{
		{
			name:    "valid job ID",
			jobID:   "12345",
			wantErr: false,
		},
		{
			name:    "empty job ID",
			jobID:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analytics, err := jm.GetJobComprehensiveAnalytics(ctx, tt.jobID)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, analytics)
			} else {
				require.NoError(t, err)
				require.NotNil(t, analytics)
				
				// Verify all components are present
				assert.NotNil(t, analytics.CPUAnalytics)
				assert.NotNil(t, analytics.MemoryAnalytics)
				assert.NotNil(t, analytics.IOAnalytics)
				assert.NotNil(t, analytics.EfficiencyMetrics)
				
				// Verify efficiency metrics
				assert.Equal(t, 43.33, analytics.EfficiencyMetrics.OverallEfficiencyScore)
				assert.Equal(t, 50.0, analytics.EfficiencyMetrics.CPUEfficiency)
				assert.Equal(t, 60.0, analytics.EfficiencyMetrics.MemoryEfficiency)
				assert.Equal(t, 20.0, analytics.EfficiencyMetrics.IOEfficiency)
				assert.Equal(t, 0.0, analytics.EfficiencyMetrics.GPUEfficiency)
				
				// Verify resource waste
				waste := analytics.EfficiencyMetrics.ResourceWaste
				assert.Equal(t, 50.0, waste["cpu_cores"])
				assert.InDelta(t, 40.0, waste["memory_gb"], 0.1)
				assert.Equal(t, 0.0, waste["gpu_hours"])
				
				// Verify bottlenecks
				assert.Len(t, analytics.EfficiencyMetrics.Bottlenecks, 2)
				if len(analytics.EfficiencyMetrics.Bottlenecks) >= 2 {
					// CPU bottleneck
					cpuBottleneck := analytics.EfficiencyMetrics.Bottlenecks[0]
					assert.Equal(t, "cpu", cpuBottleneck.Resource)
					assert.Equal(t, "underutilization", cpuBottleneck.Type)
					assert.Equal(t, "moderate", cpuBottleneck.Severity)
					assert.Equal(t, "high", cpuBottleneck.Impact)
					assert.Contains(t, cpuBottleneck.Description, "50% utilization")
					
					// Memory bottleneck
					memBottleneck := analytics.EfficiencyMetrics.Bottlenecks[1]
					assert.Equal(t, "memory", memBottleneck.Resource)
					assert.Equal(t, "overallocation", memBottleneck.Type)
					assert.Equal(t, "low", memBottleneck.Severity)
					assert.Equal(t, "low", memBottleneck.Impact)
					assert.Contains(t, memBottleneck.Description, "allocated but unused")
				}
				
				// Verify optimization recommendations
				assert.Len(t, analytics.OptimizationRecommendations, 2)
				if len(analytics.OptimizationRecommendations) >= 2 {
					assert.Contains(t, analytics.OptimizationRecommendations[0], "CPU cores")
					assert.Contains(t, analytics.OptimizationRecommendations[1], "memory allocation")
				}
			}
		})
	}
}

// TestJobManager_AnalyticsErrorHandling tests error handling for analytics methods
func TestJobManager_AnalyticsErrorHandling(t *testing.T) {
	ctx := context.Background()
	jm := &JobManager{}

	t.Run("nil context handling", func(t *testing.T) {
		// GetJobCPUAnalytics with nil context should handle gracefully
		analytics, err := jm.GetJobCPUAnalytics(nil, "12345")
		assert.NoError(t, err)
		assert.NotNil(t, analytics)
	})

	t.Run("invalid job ID format", func(t *testing.T) {
		// Empty job ID should return error
		analytics, err := jm.GetJobCPUAnalytics(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, analytics)
		assert.Contains(t, err.Error(), "job ID is required")
	})

	t.Run("consistent error messages", func(t *testing.T) {
		// All methods should return consistent error for empty job ID
		methods := []func(context.Context, string) (interface{}, error){
			func(ctx context.Context, jobID string) (interface{}, error) {
				return jm.GetJobCPUAnalytics(ctx, jobID)
			},
			func(ctx context.Context, jobID string) (interface{}, error) {
				return jm.GetJobMemoryAnalytics(ctx, jobID)
			},
			func(ctx context.Context, jobID string) (interface{}, error) {
				return jm.GetJobIOAnalytics(ctx, jobID)
			},
			func(ctx context.Context, jobID string) (interface{}, error) {
				return jm.GetJobComprehensiveAnalytics(ctx, jobID)
			},
		}

		for _, method := range methods {
			result, err := method(ctx, "")
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), "job ID is required")
		}
	})
}

// BenchmarkJobManager_GetJobCPUAnalytics benchmarks GetJobCPUAnalytics performance
func BenchmarkJobManager_GetJobCPUAnalytics(b *testing.B) {
	ctx := context.Background()
	jm := &JobManager{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = jm.GetJobCPUAnalytics(ctx, "12345")
	}
}

// BenchmarkJobManager_GetJobComprehensiveAnalytics benchmarks GetJobComprehensiveAnalytics performance
func BenchmarkJobManager_GetJobComprehensiveAnalytics(b *testing.B) {
	ctx := context.Background()
	jm := &JobManager{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = jm.GetJobComprehensiveAnalytics(ctx, "12345")
	}
}

// mockAPIClient is a mock implementation for testing
type mockAPIClient struct {
	getJobResponse *interfaces.Job
	getJobError    error
}

// createMockClientWithResponses wraps mockAPIClient to satisfy ClientWithResponses interface
func createMockClientWithResponses(mock *mockAPIClient) *ClientWithResponses {
	// We need to create a mock that implements the ClientInterface
	// For now, we'll use a type assertion workaround
	return &ClientWithResponses{
		ClientInterface: mock,
	}
}

func (m *mockAPIClient) SlurmV0040GetJobWithResponse(ctx context.Context, jobID string, params *SlurmV0040GetJobParams) (*SlurmV0040GetJobResponse, error) {
	if m.getJobError != nil {
		return nil, m.getJobError
	}

	// Convert jobID to int32
	jobIDInt, _ := strconv.ParseInt(jobID, 10, 32)
	jobIDInt32 := int32(jobIDInt)

	// Mock response
	resp := &SlurmV0040GetJobResponse{
		JSON200: &V0040OpenapiJobInfoResp{
			Jobs: []V0040JobInfo{
				{
					JobId:     &jobIDInt32,
					Name:      &m.getJobResponse.Name,
					JobState:  &[]string{m.getJobResponse.State},
					Partition: &m.getJobResponse.Partition,
					Cpus: &V0040Uint32NoVal{
						Number: &[]int32{int32(m.getJobResponse.CPUs)}[0],
						Set:    &[]bool{true}[0],
					},
					MemoryPerNode: &V0040Uint64NoVal{
						Number: &[]int64{m.getJobResponse.Memory / (1024 * 1024)}[0],
						Set:    &[]bool{true}[0],
					},
				},
			},
		},
	}

	return resp, nil
}

func (m *mockAPIClient) SlurmV0040GetJobsWithResponse(ctx context.Context, params *SlurmV0040GetJobsParams) (*SlurmV0040GetJobsResponse, error) {
	// Not used in these tests
	return nil, nil
}