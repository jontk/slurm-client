package v0_0_41

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

// TestGetJobUtilization_v41 tests the GetJobUtilization method for v0.0.41
func TestGetJobUtilization_v41(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_BasicMetrics", func(t *testing.T) {
		mockJob := &interfaces.Job{
			ID:        "41001",
			Name:      "basic-job",
			State:     "RUNNING",
			CPUs:      4,
			Memory:    16 * 1024 * 1024 * 1024, // 16GB
			Partition: "compute",
			Nodes:     []string{"node041"},
			SubmitTime: time.Now().Add(-1 * time.Hour),
		}

		manager := &JobManagerImpl{
			client: &WrapperClient{
				apiClient: &mockAPIClient{
					getJobResponse: mockJob,
				},
			},
		}

		utilization, err := manager.GetJobUtilization(ctx, "41001")
		require.NoError(t, err)
		require.NotNil(t, utilization)

		// Verify basic fields
		assert.Equal(t, "41001", utilization.JobID)
		assert.Equal(t, "basic-job", utilization.JobName)

		// Verify CPU utilization (v0.0.41 has fixed 75%)
		assert.NotNil(t, utilization.CPUUtilization)
		assert.Equal(t, 75.0, utilization.CPUUtilization.Percentage)

		// Verify memory utilization (v0.0.41 has fixed 65%)
		assert.NotNil(t, utilization.MemoryUtilization)
		assert.Equal(t, 65.0, utilization.MemoryUtilization.Percentage)

		// Verify no advanced metrics in v0.0.41
		assert.Nil(t, utilization.GPUUtilization)
		assert.Nil(t, utilization.IOUtilization)
		assert.Nil(t, utilization.NetworkUtilization)
		assert.Nil(t, utilization.EnergyUsage)

		// Verify metadata
		assert.Equal(t, "v0.0.41", utilization.Metadata["version"])
		assert.Equal(t, "basic", utilization.Metadata["feature_level"])
		
		limitations := utilization.Metadata["limitations"].([]string)
		assert.Contains(t, limitations, "no_gpu_metrics")
		assert.Contains(t, limitations, "no_io_metrics")
		assert.Contains(t, limitations, "no_network_metrics")
		assert.Contains(t, limitations, "no_energy_metrics")
		assert.Contains(t, limitations, "basic_cpu_memory_only")
	})

	t.Run("WithGPUMetadata", func(t *testing.T) {
		// Even with GPU metadata, v0.0.41 doesn't support GPU metrics
		mockJob := &interfaces.Job{
			ID:        "41002",
			Name:      "gpu-job-ignored",
			State:     "RUNNING",
			CPUs:      8,
			Memory:    32 * 1024 * 1024 * 1024,
			Partition: "gpu",
			SubmitTime: time.Now().Add(-30 * time.Minute),
			Metadata: map[string]interface{}{
				"gpu_count": 2,
			},
		}

		manager := &JobManagerImpl{
			client: &WrapperClient{
				apiClient: &mockAPIClient{
					getJobResponse: mockJob,
				},
			},
		}

		utilization, err := manager.GetJobUtilization(ctx, "41002")
		require.NoError(t, err)
		
		// Should still have no GPU utilization in v0.0.41
		assert.Nil(t, utilization.GPUUtilization)
	})
}

// TestGetJobEfficiency_v41 tests the GetJobEfficiency method for v0.0.41
func TestGetJobEfficiency_v41(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_CPUMemoryOnly", func(t *testing.T) {
		mockJob := &interfaces.Job{
			ID:        "41003",
			Name:      "efficiency-test",
			State:     "COMPLETED",
			CPUs:      8,
			Memory:    32 * 1024 * 1024 * 1024,
			Partition: "compute",
			SubmitTime: time.Now().Add(-2 * time.Hour),
			EndTime:    &[]time.Time{time.Now()}[0],
		}

		manager := &JobManagerImpl{
			client: &WrapperClient{
				apiClient: &mockAPIClient{
					getJobResponse: mockJob,
				},
			},
		}

		efficiency, err := manager.GetJobEfficiency(ctx, "41003")
		require.NoError(t, err)
		require.NotNil(t, efficiency)

		// v0.0.41: CPU=75%*0.6 + Mem=65%*0.4 = 45 + 26 = 71%
		assert.Equal(t, 71.0, efficiency.Percentage)

		// Verify metadata
		assert.Equal(t, "basic_cpu_memory_v41", efficiency.Metadata["calculation_method"])
		assert.Equal(t, "v0.0.41", efficiency.Metadata["version"])
		
		weights := efficiency.Metadata["weights"].(map[string]float64)
		assert.Equal(t, 0.6, weights["cpu"])
		assert.Equal(t, 0.4, weights["memory"])
		assert.Len(t, weights, 2) // Only CPU and memory

		limitations := efficiency.Metadata["limitations"].([]string)
		assert.Contains(t, limitations, "cpu_memory_only")
	})
}

// TestGetJobPerformance_v41 tests the GetJobPerformance method for v0.0.41
func TestGetJobPerformance_v41(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_MinimalFeatures", func(t *testing.T) {
		startTime := time.Now().Add(-1 * time.Hour)
		endTime := time.Now()
		mockJob := &interfaces.Job{
			ID:         "41004",
			Name:       "perf-test-v41",
			State:      "COMPLETED",
			CPUs:       16,
			Memory:     64 * 1024 * 1024 * 1024,
			Partition:  "compute",
			SubmitTime: time.Now().Add(-2 * time.Hour),
			StartTime:  &startTime,
			EndTime:    &endTime,
			ExitCode:   0,
		}

		manager := &JobManagerImpl{
			client: &WrapperClient{
				apiClient: &mockAPIClient{
					getJobResponse: mockJob,
				},
			},
		}

		performance, err := manager.GetJobPerformance(ctx, "41004")
		require.NoError(t, err)
		require.NotNil(t, performance)

		// Verify basic fields
		assert.Equal(t, uint32(41004), performance.JobID)
		assert.Equal(t, "perf-test-v41", performance.JobName)

		// Verify no advanced features in v0.0.41
		assert.Nil(t, performance.StepMetrics)
		assert.Nil(t, performance.PerformanceTrends)

		// Verify basic bottlenecks
		assert.NotNil(t, performance.Bottlenecks)
		// Should be empty unless utilization is high

		// Verify recommendations
		assert.NotNil(t, performance.Recommendations)
		// Should include API upgrade recommendation
		var foundUpgradeRec bool
		for _, rec := range performance.Recommendations {
			if rec.Title == "Limited analytics in API v0.0.41" {
				foundUpgradeRec = true
				assert.Equal(t, "v0.0.42_or_higher", rec.ConfigChanges["recommended_api_version"])
			}
		}
		assert.True(t, foundUpgradeRec)
	})
}

// TestHelperFunctions_v41 tests the v0.0.41-specific helper functions
func TestHelperFunctions_v41(t *testing.T) {
	t.Run("analyzeBottlenecksV41", func(t *testing.T) {
		// Test with high utilization
		utilization := &interfaces.JobUtilization{
			CPUUtilization: &interfaces.ResourceUtilization{
				Percentage: 87.0, // Above 85% threshold
			},
			MemoryUtilization: &interfaces.ResourceUtilization{
				Percentage: 82.0, // Above 80% threshold
			},
		}

		bottlenecks := analyzeBottlenecksV41(utilization)
		assert.Len(t, bottlenecks, 2)

		// Check CPU bottleneck
		assert.Equal(t, "cpu", bottlenecks[0].Type)
		assert.Equal(t, "medium", bottlenecks[0].Severity)
		assert.Equal(t, 10.0, bottlenecks[0].Impact)

		// Check memory bottleneck
		assert.Equal(t, "memory", bottlenecks[1].Type)
		assert.Equal(t, "low", bottlenecks[1].Severity)
		assert.Equal(t, 5.0, bottlenecks[1].Impact)
	})

	t.Run("generateRecommendationsV41", func(t *testing.T) {
		// Test with low efficiency
		efficiency := &interfaces.ResourceUtilization{
			Percentage: 65.0,
		}

		recommendations := generateRecommendationsV41(efficiency)
		assert.Len(t, recommendations, 2)

		// Should have resource utilization recommendation
		assert.Equal(t, "workflow", recommendations[0].Type)
		assert.Equal(t, "low", recommendations[0].Priority)

		// Should always have API upgrade recommendation
		assert.Equal(t, "configuration", recommendations[1].Type)
		assert.Equal(t, "Limited analytics in API v0.0.41", recommendations[1].Title)
	})
}

// mockAPIClient is a mock implementation for testing
type mockAPIClient struct {
	getJobResponse *interfaces.Job
	getJobError    error
}

func (m *mockAPIClient) SlurmV0041GetJobWithResponse(ctx context.Context, jobID string, params *SlurmV0041GetJobParams) (*SlurmV0041GetJobResponse, error) {
	if m.getJobError != nil {
		return nil, m.getJobError
	}

	// Convert jobID to int32
	jobIDInt, _ := strconv.ParseInt(jobID, 10, 32)
	jobIDInt32 := int32(jobIDInt)

	// Mock response
	resp := &SlurmV0041GetJobResponse{
		JSON200: &V0041OpenapiJobInfoResp{
			Jobs: []V0041JobInfo{
				{
					JobId:     &jobIDInt32,
					Name:      &m.getJobResponse.Name,
					JobState:  &[]string{m.getJobResponse.State},
					Partition: &m.getJobResponse.Partition,
					Cpus: &V0041Uint32NoVal{
						Number: &[]int32{int32(m.getJobResponse.CPUs)}[0],
						Set:    &[]bool{true}[0],
					},
					MemoryPerNode: &V0041Uint64NoVal{
						Number: &[]int64{m.getJobResponse.Memory / (1024 * 1024)}[0],
						Set:    &[]bool{true}[0],
					},
				},
			},
		},
	}

	return resp, nil
}

func (m *mockAPIClient) SlurmV0041GetJobsWithResponse(ctx context.Context, params *SlurmV0041GetJobsParams) (*SlurmV0041GetJobsResponse, error) {
	// Not used in these tests
	return nil, nil
}

// TestJobManager_GetJobCPUAnalytics_v41 tests the GetJobCPUAnalytics method for v0.0.41
func TestJobManager_GetJobCPUAnalytics_v41(t *testing.T) {
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
				
				// Verify basic analytics for v0.0.41
				assert.Equal(t, 8, analytics.AllocatedCores)
				assert.Equal(t, 8, analytics.RequestedCores)
				assert.Equal(t, 5.2, analytics.UsedCores)
				assert.Equal(t, 65.0, analytics.UtilizationPercent)
				assert.Equal(t, 65.0, analytics.EfficiencyPercent)
				assert.Equal(t, 2.8, analytics.IdleCores)
				assert.False(t, analytics.Oversubscribed)
				
				// Verify core-level metrics
				assert.Len(t, analytics.CoreMetrics, 8)
				for i, metric := range analytics.CoreMetrics {
					assert.Equal(t, i, metric.CoreID)
					assert.Equal(t, 60.0+float64(i), metric.UtilizationPercent)
					assert.Equal(t, 2.8, metric.Frequency)
					assert.Equal(t, 65.0, metric.Temperature)
					assert.Equal(t, 1000+i*100, metric.CacheMisses)
					assert.Equal(t, int64(2), metric.InstructionsPerCycle)
				}
			}
		})
	}
}

// TestJobManager_GetJobMemoryAnalytics_v41 tests the GetJobMemoryAnalytics method for v0.0.41
func TestJobManager_GetJobMemoryAnalytics_v41(t *testing.T) {
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
				
				// Verify basic analytics for v0.0.41
				assert.Equal(t, int64(32*1024*1024*1024), analytics.AllocatedBytes)
				assert.Equal(t, int64(32*1024*1024*1024), analytics.RequestedBytes)
				assert.InDelta(t, int64(22.4*1024*1024*1024), analytics.UsedBytes, 1024)
				assert.Equal(t, 70.0, analytics.UtilizationPercent)
				assert.InDelta(t, int64(24.32*1024*1024*1024), analytics.MaxUsedBytes, 1024)
				assert.Equal(t, int64(2*1024*1024*1024), analytics.CacheBytes)
				assert.Equal(t, int64(512*1024*1024), analytics.SwapBytes)
				assert.Equal(t, 20000, analytics.PageFaults)
				assert.Equal(t, 200, analytics.MajorPageFaults)
				
				// NUMA metrics should have better accuracy in v0.0.41
				assert.Len(t, analytics.NUMAMetrics, 2)
				if len(analytics.NUMAMetrics) >= 2 {
					// NUMA node 0
					numa0 := analytics.NUMAMetrics[0]
					assert.Equal(t, 0, numa0.NodeID)
					assert.InDelta(t, int64(11.2*1024*1024*1024), numa0.UsedBytes, 1024)
					assert.Equal(t, 90.0, numa0.LocalAccessPercent)
					assert.Equal(t, 10.0, numa0.RemoteAccessPercent)
					
					// NUMA node 1
					numa1 := analytics.NUMAMetrics[1]
					assert.Equal(t, 1, numa1.NodeID)
					assert.InDelta(t, int64(11.2*1024*1024*1024), numa1.UsedBytes, 1024)
					assert.Equal(t, 85.0, numa1.LocalAccessPercent)
					assert.Equal(t, 15.0, numa1.RemoteAccessPercent)
				}
				
				// Basic memory leak detection in v0.0.41
				assert.False(t, analytics.MemoryLeakDetected)
				assert.Equal(t, 0.0, analytics.LeakRatePerHour)
			}
		})
	}
}

// TestJobManager_GetJobIOAnalytics_v41 tests the GetJobIOAnalytics method for v0.0.41
func TestJobManager_GetJobIOAnalytics_v41(t *testing.T) {
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
				
				// Verify basic analytics for v0.0.41
				assert.Equal(t, int64(50*1024*1024*1024), analytics.ReadBytes)
				assert.Equal(t, int64(25*1024*1024*1024), analytics.WriteBytes)
				assert.Equal(t, int64(5000), analytics.ReadOperations)
				assert.Equal(t, int64(2500), analytics.WriteOperations)
				assert.Equal(t, 200.0, analytics.ReadBandwidthMBps)
				assert.Equal(t, 100.0, analytics.WriteBandwidthMBps)
				assert.Equal(t, 8.0, analytics.ReadLatencyMs)
				assert.Equal(t, 12.0, analytics.WriteLatencyMs)
				assert.Equal(t, 15.0, analytics.IOWaitPercent)
				
				// Device metrics should have better detail in v0.0.41
				assert.Len(t, analytics.DeviceMetrics, 2)
				if len(analytics.DeviceMetrics) >= 2 {
					// Local device
					local := analytics.DeviceMetrics[0]
					assert.Equal(t, "/dev/nvme0n1", local.DeviceName)
					assert.Equal(t, "nvme", local.DeviceType)
					assert.Equal(t, 35.0, local.UtilizationPercent)
					assert.Equal(t, 300.0, local.ReadBandwidthMBps)
					assert.Equal(t, 150.0, local.WriteBandwidthMBps)
					assert.Equal(t, int64(3000), local.ReadOps)
					assert.Equal(t, int64(1500), local.WriteOps)
					assert.Equal(t, 15.0, local.AvgQueueDepth)
					
					// Network device
					network := analytics.DeviceMetrics[1]
					assert.Equal(t, "/mnt/nfs/scratch", network.DeviceName)
					assert.Equal(t, "network", network.DeviceType)
					assert.Equal(t, 25.0, network.UtilizationPercent)
					assert.Equal(t, 100.0, network.ReadBandwidthMBps)
					assert.Equal(t, 50.0, network.WriteBandwidthMBps)
					assert.Equal(t, int64(2000), network.ReadOps)
					assert.Equal(t, int64(1000), network.WriteOps)
					assert.Equal(t, 20.0, network.AvgQueueDepth)
				}
			}
		})
	}
}

// TestJobManager_GetJobComprehensiveAnalytics_v41 tests the GetJobComprehensiveAnalytics method for v0.0.41
func TestJobManager_GetJobComprehensiveAnalytics_v41(t *testing.T) {
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
				
				// Verify efficiency metrics for v0.0.41
				assert.Equal(t, 50.0, analytics.EfficiencyMetrics.OverallEfficiencyScore)
				assert.Equal(t, 65.0, analytics.EfficiencyMetrics.CPUEfficiency)
				assert.Equal(t, 70.0, analytics.EfficiencyMetrics.MemoryEfficiency)
				assert.Equal(t, 30.0, analytics.EfficiencyMetrics.IOEfficiency)
				assert.Equal(t, 0.0, analytics.EfficiencyMetrics.GPUEfficiency)
				
				// Verify resource waste
				waste := analytics.EfficiencyMetrics.ResourceWaste
				assert.Equal(t, 35.0, waste["cpu_cores"])
				assert.InDelta(t, 30.0, waste["memory_gb"], 0.1)
				assert.Equal(t, 0.0, waste["gpu_hours"])
				
				// Verify bottlenecks
				assert.Len(t, analytics.EfficiencyMetrics.Bottlenecks, 3)
				if len(analytics.EfficiencyMetrics.Bottlenecks) >= 3 {
					// CPU bottleneck
					cpuBottleneck := analytics.EfficiencyMetrics.Bottlenecks[0]
					assert.Equal(t, "cpu", cpuBottleneck.Resource)
					assert.Equal(t, "underutilization", cpuBottleneck.Type)
					assert.Equal(t, "moderate", cpuBottleneck.Severity)
					assert.Equal(t, "moderate", cpuBottleneck.Impact)
					assert.Contains(t, cpuBottleneck.Description, "65% utilization")
					
					// Memory bottleneck
					memBottleneck := analytics.EfficiencyMetrics.Bottlenecks[1]
					assert.Equal(t, "memory", memBottleneck.Resource)
					assert.Equal(t, "overallocation", memBottleneck.Type)
					assert.Equal(t, "low", memBottleneck.Severity)
					assert.Equal(t, "low", memBottleneck.Impact)
					assert.Contains(t, memBottleneck.Description, "allocated but unused")
					
					// I/O bottleneck
					ioBottleneck := analytics.EfficiencyMetrics.Bottlenecks[2]
					assert.Equal(t, "io", ioBottleneck.Resource)
					assert.Equal(t, "low_throughput", ioBottleneck.Type)
					assert.Equal(t, "moderate", ioBottleneck.Severity)
					assert.Equal(t, "moderate", ioBottleneck.Impact)
					assert.Contains(t, ioBottleneck.Description, "low I/O efficiency")
				}
				
				// Verify optimization recommendations
				assert.Len(t, analytics.OptimizationRecommendations, 3)
				if len(analytics.OptimizationRecommendations) >= 3 {
					assert.Contains(t, analytics.OptimizationRecommendations[0], "CPU cores")
					assert.Contains(t, analytics.OptimizationRecommendations[1], "memory allocation")
					assert.Contains(t, analytics.OptimizationRecommendations[2], "I/O pattern")
				}
			}
		})
	}
}

// TestJobManager_AnalyticsErrorHandling_v41 tests error handling for analytics methods
func TestJobManager_AnalyticsErrorHandling_v41(t *testing.T) {
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

// TestJobManager_ProgressiveEnhancement_v41 tests progressive enhancement features
func TestJobManager_ProgressiveEnhancement_v41(t *testing.T) {
	ctx := context.Background()
	jm := &JobManager{}

	t.Run("Enhanced CPU metrics", func(t *testing.T) {
		analytics, err := jm.GetJobCPUAnalytics(ctx, "12345")
		require.NoError(t, err)
		
		// v0.0.41 has better CPU metrics than v0.0.40
		assert.Greater(t, analytics.UtilizationPercent, 50.0)
		assert.Greater(t, len(analytics.CoreMetrics), 4)
		
		// Check for progressive values
		for _, core := range analytics.CoreMetrics {
			assert.Greater(t, core.CacheMisses, 0)
			assert.Greater(t, core.InstructionsPerCycle, int64(1))
		}
	})

	t.Run("NUMA awareness", func(t *testing.T) {
		analytics, err := jm.GetJobMemoryAnalytics(ctx, "12345")
		require.NoError(t, err)
		
		// v0.0.41 has NUMA awareness
		assert.Greater(t, len(analytics.NUMAMetrics), 1)
		for _, numa := range analytics.NUMAMetrics {
			assert.Greater(t, numa.LocalAccessPercent, 0.0)
			assert.Greater(t, numa.RemoteAccessPercent, 0.0)
		}
	})

	t.Run("Multiple device types", func(t *testing.T) {
		analytics, err := jm.GetJobIOAnalytics(ctx, "12345")
		require.NoError(t, err)
		
		// v0.0.41 supports multiple device types
		assert.Greater(t, len(analytics.DeviceMetrics), 1)
		
		// Check for different device types
		deviceTypes := make(map[string]bool)
		for _, device := range analytics.DeviceMetrics {
			deviceTypes[device.DeviceType] = true
		}
		assert.Greater(t, len(deviceTypes), 1)
	})
}

// BenchmarkJobManager_GetJobCPUAnalytics_v41 benchmarks GetJobCPUAnalytics performance
func BenchmarkJobManager_GetJobCPUAnalytics_v41(b *testing.B) {
	ctx := context.Background()
	jm := &JobManager{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = jm.GetJobCPUAnalytics(ctx, "12345")
	}
}

// BenchmarkJobManager_GetJobComprehensiveAnalytics_v41 benchmarks GetJobComprehensiveAnalytics performance
func BenchmarkJobManager_GetJobComprehensiveAnalytics_v41(b *testing.B) {
	ctx := context.Background()
	jm := &JobManager{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = jm.GetJobComprehensiveAnalytics(ctx, "12345")
	}
}