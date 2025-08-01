// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

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

// TestGetJobUtilization_v42 tests the GetJobUtilization method for v0.0.42
func TestGetJobUtilization_v42(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_Enhanced", func(t *testing.T) {
		// Create mock job manager
		mockJob := &interfaces.Job{
			ID:        "54321",
			Name:      "test-job-v42",
			State:     "RUNNING",
			CPUs:      8,
			Memory:    32 * 1024 * 1024 * 1024, // 32GB
			Partition: "compute",
			Nodes:     []string{"node003"},
			SubmitTime: time.Now().Add(-1 * time.Hour),
			StartTime:  &[]time.Time{time.Now().Add(-30 * time.Minute)}[0],
			Metadata: map[string]interface{}{
				"gpu_count": 1,
			},
		}

		manager := &JobManagerImpl{
			client: &WrapperClient{
				apiClient: &mockAPIClient{
					getJobResponse: mockJob,
				},
			},
		}

		// Test GetJobUtilization
		utilization, err := manager.GetJobUtilization(ctx, "54321")
		require.NoError(t, err)
		require.NotNil(t, utilization)

		// Verify basic fields
		assert.Equal(t, "54321", utilization.JobID)
		assert.Equal(t, "test-job-v42", utilization.JobName)

		// Verify CPU utilization (v0.0.42 has 82%)
		assert.NotNil(t, utilization.CPUUtilization)
		assert.Equal(t, 82.0, utilization.CPUUtilization.Percentage)

		// Verify memory utilization (v0.0.42 has 70%)
		assert.NotNil(t, utilization.MemoryUtilization)
		assert.Equal(t, 70.0, utilization.MemoryUtilization.Percentage)

		// Verify GPU utilization (limited in v0.0.42)
		assert.NotNil(t, utilization.GPUUtilization)
		assert.Equal(t, 1, utilization.GPUUtilization.TotalGPUs)
		assert.Equal(t, 88.0, utilization.GPUUtilization.AverageUtilization.Percentage)
		assert.Empty(t, utilization.GPUUtilization.GPUs) // No per-device metrics
		assert.Equal(t, true, utilization.GPUUtilization.Metadata["aggregated_only"])

		// Verify I/O utilization (basic in v0.0.42)
		assert.NotNil(t, utilization.IOUtilization)
		assert.Equal(t, 20.0, utilization.IOUtilization.ReadBandwidth.Percentage)
		assert.Empty(t, utilization.IOUtilization.FileSystems) // No per-filesystem

		// Verify network utilization (basic in v0.0.42)
		assert.NotNil(t, utilization.NetworkUtilization)
		assert.Equal(t, 8.0, utilization.NetworkUtilization.TotalBandwidth.Percentage)
		assert.Empty(t, utilization.NetworkUtilization.Interfaces) // No per-interface

		// Verify energy usage (limited in v0.0.42)
		if mockJob.EndTime != nil {
			assert.NotNil(t, utilization.EnergyUsage)
			assert.Equal(t, 0.0, utilization.EnergyUsage.CPUEnergyJoules) // No breakdown
			assert.Equal(t, false, utilization.EnergyUsage.Metadata["breakdown_available"])
		}

		// Verify metadata
		assert.Equal(t, "v0.0.42", utilization.Metadata["version"])
		assert.Equal(t, "enhanced", utilization.Metadata["feature_level"])
	})

	t.Run("NoGPU", func(t *testing.T) {
		mockJob := &interfaces.Job{
			ID:        "54322",
			Name:      "cpu-only-job",
			State:     "RUNNING",
			CPUs:      16,
			Memory:    64 * 1024 * 1024 * 1024,
			Partition: "compute",
			SubmitTime: time.Now().Add(-2 * time.Hour),
		}

		manager := &JobManagerImpl{
			client: &WrapperClient{
				apiClient: &mockAPIClient{
					getJobResponse: mockJob,
				},
			},
		}

		utilization, err := manager.GetJobUtilization(ctx, "54322")
		require.NoError(t, err)
		
		// Should not have GPU utilization
		assert.Nil(t, utilization.GPUUtilization)
	})
}

// TestGetJobEfficiency_v42 tests the GetJobEfficiency method for v0.0.42
func TestGetJobEfficiency_v42(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_SimpleCalculation", func(t *testing.T) {
		mockJob := &interfaces.Job{
			ID:        "54323",
			Name:      "efficiency-test",
			State:     "COMPLETED",
			CPUs:      16,
			Memory:    64 * 1024 * 1024 * 1024,
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

		efficiency, err := manager.GetJobEfficiency(ctx, "54323")
		require.NoError(t, err)
		require.NotNil(t, efficiency)

		// Verify efficiency calculation
		// v0.0.42: CPU=82%*0.5 + Mem=70%*0.3 + IO=15%*0.05 = 41 + 21 + 0.75 = 62.75%
		assert.InDelta(t, 62.75, efficiency.Percentage, 0.1)

		// Verify metadata
		assert.Equal(t, "weighted_average_v42", efficiency.Metadata["calculation_method"])
		assert.Equal(t, "v0.0.42", efficiency.Metadata["version"])
		
		weights := efficiency.Metadata["weights"].(map[string]float64)
		assert.Equal(t, 0.5, weights["cpu"])
		assert.Equal(t, 0.3, weights["memory"])
		assert.Equal(t, 0.15, weights["gpu"])
		assert.Equal(t, 0.05, weights["io"])
	})
}

// TestGetJobPerformance_v42 tests the GetJobPerformance method for v0.0.42
func TestGetJobPerformance_v42(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_LimitedFeatures", func(t *testing.T) {
		startTime := time.Now().Add(-2 * time.Hour)
		endTime := time.Now()
		mockJob := &interfaces.Job{
			ID:         "54324",
			Name:       "perf-test-v42",
			State:      "COMPLETED",
			CPUs:       16,
			Memory:     64 * 1024 * 1024 * 1024,
			Partition:  "compute",
			SubmitTime: time.Now().Add(-3 * time.Hour),
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

		performance, err := manager.GetJobPerformance(ctx, "54324")
		require.NoError(t, err)
		require.NotNil(t, performance)

		// Verify basic fields
		assert.Equal(t, uint32(54324), performance.JobID)
		assert.Equal(t, "perf-test-v42", performance.JobName)

		// Verify no step metrics in v0.0.42
		assert.Empty(t, performance.StepMetrics)

		// Verify performance trends (limited in v0.0.42)
		assert.NotNil(t, performance.PerformanceTrends)
		assert.LessOrEqual(t, len(performance.PerformanceTrends.TimePoints), 12) // Max 12 points
		assert.Nil(t, performance.PerformanceTrends.GPUTrends) // No GPU trends
		assert.Nil(t, performance.PerformanceTrends.PowerTrends) // No power trends

		// Verify bottlenecks (basic in v0.0.42)
		assert.NotNil(t, performance.Bottlenecks)

		// Verify recommendations (basic in v0.0.42)
		assert.NotNil(t, performance.Recommendations)
	})
}

// TestHelperFunctions_v42 tests the v0.0.42-specific helper functions
func TestHelperFunctions_v42(t *testing.T) {
	t.Run("generatePerformanceTrendsV42", func(t *testing.T) {
		startTime := time.Now().Add(-10 * time.Hour)
		endTime := time.Now()
		job := &interfaces.Job{
			StartTime: &startTime,
			EndTime:   &endTime,
		}

		trends := generatePerformanceTrendsV42(job)
		assert.NotNil(t, trends)
		
		// v0.0.42 limits to 12 data points
		assert.LessOrEqual(t, len(trends.TimePoints), 12)
		assert.Nil(t, trends.GPUTrends)
		assert.Nil(t, trends.PowerTrends)
	})

	t.Run("analyzeBottlenecksV42", func(t *testing.T) {
		// Test with high CPU (>90% for v0.0.42)
		utilization := &interfaces.JobUtilization{
			CPUUtilization: &interfaces.ResourceUtilization{
				Percentage: 92.0,
			},
			MemoryUtilization: &interfaces.ResourceUtilization{
				Percentage: 86.0,
			},
			IOUtilization: &interfaces.IOUtilization{
				ReadBandwidth: &interfaces.ResourceUtilization{
					Percentage: 75.0,
				},
				WriteBandwidth: &interfaces.ResourceUtilization{
					Percentage: 75.0,
				},
			},
		}

		bottlenecks := analyzeBottlenecksV42(utilization)
		assert.Len(t, bottlenecks, 3) // CPU, Memory, I/O

		// Check severities
		cpuBottleneck := bottlenecks[0]
		assert.Equal(t, "cpu", cpuBottleneck.Type)
		assert.Equal(t, "high", cpuBottleneck.Severity)

		memoryBottleneck := bottlenecks[1]
		assert.Equal(t, "memory", memoryBottleneck.Type)
		assert.Equal(t, "medium", memoryBottleneck.Severity)

		ioBottleneck := bottlenecks[2]
		assert.Equal(t, "io", ioBottleneck.Type)
		assert.Equal(t, "low", ioBottleneck.Severity)
	})

	t.Run("generateRecommendationsV42", func(t *testing.T) {
		utilization := &interfaces.JobUtilization{
			CPUUtilization: &interfaces.ResourceUtilization{
				Percentage: 55.0,
			},
			MemoryUtilization: &interfaces.ResourceUtilization{
				Percentage: 45.0,
			},
		}
		efficiency := &interfaces.ResourceUtilization{
			Percentage: 72.0,
		}

		recommendations := generateRecommendationsV42(utilization, efficiency)
		assert.NotEmpty(t, recommendations)

		// Should have CPU and memory recommendations
		var foundCPU, foundMemory, foundEfficiency bool
		for _, rec := range recommendations {
			if rec.Title == "Consider reducing CPU allocation" {
				foundCPU = true
				assert.Equal(t, "medium", rec.Priority)
			}
			if rec.Title == "Memory allocation can be optimized" {
				foundMemory = true
				assert.Equal(t, "low", rec.Priority)
			}
			if rec.Title == "Job efficiency could be improved" {
				foundEfficiency = true
				assert.Equal(t, "medium", rec.Priority)
			}
		}
		assert.True(t, foundCPU)
		assert.True(t, foundMemory)
		assert.True(t, foundEfficiency)
	})
}

// mockAPIClient is a mock implementation for testing
type mockAPIClient struct {
	getJobResponse *interfaces.Job
	getJobError    error
}

func (m *mockAPIClient) SlurmV0042GetJobWithResponse(ctx context.Context, jobID string, params *SlurmV0042GetJobParams) (*SlurmV0042GetJobResponse, error) {
	if m.getJobError != nil {
		return nil, m.getJobError
	}

	// Convert jobID to int32
	jobIDInt, _ := strconv.ParseInt(jobID, 10, 32)
	jobIDInt32 := int32(jobIDInt)

	// Mock response
	resp := &SlurmV0042GetJobResponse{
		JSON200: &V0042OpenapiJobInfoResp{
			Jobs: []V0042JobInfo{
				{
					JobId:     &jobIDInt32,
					Name:      &m.getJobResponse.Name,
					JobState:  &[]string{m.getJobResponse.State},
					Partition: &m.getJobResponse.Partition,
					Cpus: &V0042Uint32NoValStruct{
						Number: &[]int32{int32(m.getJobResponse.CPUs)}[0],
						Set:    &[]bool{true}[0],
					},
					MemoryPerNode: &V0042Uint64NoValStruct{
						Number: &[]int64{m.getJobResponse.Memory / (1024 * 1024)}[0], // Convert to MB
						Set:    &[]bool{true}[0],
					},
				},
			},
		},
	}

	return resp, nil
}

func (m *mockAPIClient) SlurmV0042GetJobsWithResponse(ctx context.Context, params *SlurmV0042GetJobsParams) (*SlurmV0042GetJobsResponse, error) {
	// Not used in these tests
	return nil, nil
}

// TestJobManager_GetJobCPUAnalytics_v42 tests the GetJobCPUAnalytics method for v0.0.42
func TestJobManager_GetJobCPUAnalytics_v42(t *testing.T) {
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
				
				// Verify enhanced analytics for v0.0.42
				assert.Equal(t, 16, analytics.AllocatedCores)
				assert.Equal(t, 16, analytics.RequestedCores)
				assert.Equal(t, 12.0, analytics.UsedCores)
				assert.Equal(t, 75.0, analytics.UtilizationPercent)
				assert.Equal(t, 75.0, analytics.EfficiencyPercent)
				assert.Equal(t, 4.0, analytics.IdleCores)
				assert.False(t, analytics.Oversubscribed)
				
				// Verify enhanced core-level metrics
				assert.Len(t, analytics.CoreMetrics, 16)
				for i, metric := range analytics.CoreMetrics {
					assert.Equal(t, i, metric.CoreID)
					assert.Equal(t, 70.0+float64(i%4)*5, metric.UtilizationPercent)
					assert.Equal(t, 3.2, metric.Frequency)
					assert.Equal(t, 70.0, metric.Temperature)
					assert.Equal(t, 2000+i*200, metric.CacheMisses)
					assert.Equal(t, int64(3), metric.InstructionsPerCycle)
				}
			}
		})
	}
}

// TestJobManager_GetJobMemoryAnalytics_v42 tests the GetJobMemoryAnalytics method for v0.0.42
func TestJobManager_GetJobMemoryAnalytics_v42(t *testing.T) {
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
				
				// Verify enhanced analytics for v0.0.42
				assert.Equal(t, int64(64*1024*1024*1024), analytics.AllocatedBytes)
				assert.Equal(t, int64(64*1024*1024*1024), analytics.RequestedBytes)
				assert.InDelta(t, int64(51.2*1024*1024*1024), analytics.UsedBytes, 1024)
				assert.Equal(t, 80.0, analytics.UtilizationPercent)
				assert.InDelta(t, int64(57.6*1024*1024*1024), analytics.MaxUsedBytes, 1024)
				assert.Equal(t, int64(4*1024*1024*1024), analytics.CacheBytes)
				assert.Equal(t, int64(256*1024*1024), analytics.SwapBytes)
				assert.Equal(t, 50000, analytics.PageFaults)
				assert.Equal(t, 500, analytics.MajorPageFaults)
				
				// Enhanced NUMA metrics in v0.0.42
				assert.Len(t, analytics.NUMAMetrics, 4)
				if len(analytics.NUMAMetrics) >= 4 {
					for i, numa := range analytics.NUMAMetrics {
						assert.Equal(t, i, numa.NodeID)
						assert.InDelta(t, int64(12.8*1024*1024*1024), numa.UsedBytes, 1024)
						assert.Equal(t, 92.0-float64(i)*2, numa.LocalAccessPercent)
						assert.Equal(t, 8.0+float64(i)*2, numa.RemoteAccessPercent)
					}
				}
				
				// Basic memory leak detection in v0.0.42
				assert.True(t, analytics.MemoryLeakDetected)
				assert.Equal(t, 100.0, analytics.LeakRatePerHour)
			}
		})
	}
}

// TestJobManager_GetJobIOAnalytics_v42 tests the GetJobIOAnalytics method for v0.0.42
func TestJobManager_GetJobIOAnalytics_v42(t *testing.T) {
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
				
				// Verify enhanced analytics for v0.0.42
				assert.Equal(t, int64(200*1024*1024*1024), analytics.ReadBytes)
				assert.Equal(t, int64(100*1024*1024*1024), analytics.WriteBytes)
				assert.Equal(t, int64(20000), analytics.ReadOperations)
				assert.Equal(t, int64(10000), analytics.WriteOperations)
				assert.Equal(t, 500.0, analytics.ReadBandwidthMBps)
				assert.Equal(t, 250.0, analytics.WriteBandwidthMBps)
				assert.Equal(t, 5.0, analytics.ReadLatencyMs)
				assert.Equal(t, 8.0, analytics.WriteLatencyMs)
				assert.Equal(t, 8.0, analytics.IOWaitPercent)
				
				// Enhanced device metrics in v0.0.42
				assert.Len(t, analytics.DeviceMetrics, 3)
				if len(analytics.DeviceMetrics) >= 3 {
					// NVMe device
					nvme := analytics.DeviceMetrics[0]
					assert.Equal(t, "/dev/nvme0n1", nvme.DeviceName)
					assert.Equal(t, "nvme", nvme.DeviceType)
					assert.Equal(t, 45.0, nvme.UtilizationPercent)
					assert.Equal(t, 1000.0, nvme.ReadBandwidthMBps)
					assert.Equal(t, 500.0, nvme.WriteBandwidthMBps)
					assert.Equal(t, int64(10000), nvme.ReadOps)
					assert.Equal(t, int64(5000), nvme.WriteOps)
					assert.Equal(t, 20.0, nvme.AvgQueueDepth)
					
					// SSD device
					ssd := analytics.DeviceMetrics[1]
					assert.Equal(t, "/dev/sdb", ssd.DeviceName)
					assert.Equal(t, "ssd", ssd.DeviceType)
					assert.Equal(t, 35.0, ssd.UtilizationPercent)
					
					// Network device
					network := analytics.DeviceMetrics[2]
					assert.Equal(t, "/mnt/lustre", network.DeviceName)
					assert.Equal(t, "network", network.DeviceType)
					assert.Equal(t, 40.0, network.UtilizationPercent)
				}
			}
		})
	}
}

// TestJobManager_GetJobComprehensiveAnalytics_v42 tests the GetJobComprehensiveAnalytics method for v0.0.42
func TestJobManager_GetJobComprehensiveAnalytics_v42(t *testing.T) {
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
				
				// Verify enhanced efficiency metrics for v0.0.42
				assert.Equal(t, 61.25, analytics.EfficiencyMetrics.OverallEfficiencyScore)
				assert.Equal(t, 75.0, analytics.EfficiencyMetrics.CPUEfficiency)
				assert.Equal(t, 80.0, analytics.EfficiencyMetrics.MemoryEfficiency)
				assert.Equal(t, 40.0, analytics.EfficiencyMetrics.IOEfficiency)
				assert.Equal(t, 50.0, analytics.EfficiencyMetrics.GPUEfficiency)
				
				// Verify resource waste
				waste := analytics.EfficiencyMetrics.ResourceWaste
				assert.Equal(t, 25.0, waste["cpu_cores"])
				assert.InDelta(t, 20.0, waste["memory_gb"], 0.1)
				assert.Equal(t, 50.0, waste["gpu_hours"])
				
				// Verify comprehensive bottlenecks
				assert.Len(t, analytics.EfficiencyMetrics.Bottlenecks, 4)
				if len(analytics.EfficiencyMetrics.Bottlenecks) >= 4 {
					// Memory leak bottleneck
					leakBottleneck := analytics.EfficiencyMetrics.Bottlenecks[0]
					assert.Equal(t, "memory", leakBottleneck.Resource)
					assert.Equal(t, "memory_leak", leakBottleneck.Type)
					assert.Equal(t, "high", leakBottleneck.Severity)
					assert.Equal(t, "high", leakBottleneck.Impact)
					assert.Contains(t, leakBottleneck.Description, "Memory leak detected")
					
					// CPU bottleneck
					cpuBottleneck := analytics.EfficiencyMetrics.Bottlenecks[1]
					assert.Equal(t, "cpu", cpuBottleneck.Resource)
					assert.Equal(t, "underutilization", cpuBottleneck.Type)
					assert.Equal(t, "low", cpuBottleneck.Severity)
					assert.Equal(t, "moderate", cpuBottleneck.Impact)
					
					// I/O bottleneck
					ioBottleneck := analytics.EfficiencyMetrics.Bottlenecks[2]
					assert.Equal(t, "io", ioBottleneck.Resource)
					assert.Equal(t, "low_throughput", ioBottleneck.Type)
					assert.Equal(t, "high", ioBottleneck.Severity)
					assert.Equal(t, "high", ioBottleneck.Impact)
					
					// GPU bottleneck
					gpuBottleneck := analytics.EfficiencyMetrics.Bottlenecks[3]
					assert.Equal(t, "gpu", gpuBottleneck.Resource)
					assert.Equal(t, "underutilization", gpuBottleneck.Type)
					assert.Equal(t, "moderate", gpuBottleneck.Severity)
					assert.Equal(t, "moderate", gpuBottleneck.Impact)
				}
				
				// Verify optimization recommendations
				assert.Greater(t, len(analytics.OptimizationRecommendations), 3)
				foundMemoryLeak := false
				for _, rec := range analytics.OptimizationRecommendations {
					if contains(rec, "memory leak") {
						foundMemoryLeak = true
						break
					}
				}
				assert.True(t, foundMemoryLeak)
			}
		})
	}
}

// TestJobManager_AnalyticsErrorHandling_v42 tests error handling for analytics methods
func TestJobManager_AnalyticsErrorHandling_v42(t *testing.T) {
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

// TestJobManager_ProgressiveEnhancement_v42 tests progressive enhancement features
func TestJobManager_ProgressiveEnhancement_v42(t *testing.T) {
	ctx := context.Background()
	jm := &JobManager{}

	t.Run("Advanced CPU metrics", func(t *testing.T) {
		analytics, err := jm.GetJobCPUAnalytics(ctx, "12345")
		require.NoError(t, err)
		
		// v0.0.42 has advanced CPU metrics
		assert.Equal(t, 75.0, analytics.UtilizationPercent)
		assert.Equal(t, 16, len(analytics.CoreMetrics))
		
		// Check for advanced metrics
		for _, core := range analytics.CoreMetrics {
			assert.Greater(t, core.CacheMisses, 1000)
			assert.Equal(t, int64(3), core.InstructionsPerCycle)
			assert.Greater(t, core.Temperature, 60.0)
		}
	})

	t.Run("Full NUMA topology", func(t *testing.T) {
		analytics, err := jm.GetJobMemoryAnalytics(ctx, "12345")
		require.NoError(t, err)
		
		// v0.0.42 has full NUMA topology support
		assert.Equal(t, 4, len(analytics.NUMAMetrics))
		for _, numa := range analytics.NUMAMetrics {
			assert.Greater(t, numa.LocalAccessPercent, 80.0)
			assert.Less(t, numa.RemoteAccessPercent, 20.0)
		}
	})

	t.Run("Memory leak detection", func(t *testing.T) {
		analytics, err := jm.GetJobMemoryAnalytics(ctx, "12345")
		require.NoError(t, err)
		
		// v0.0.42 has memory leak detection
		assert.True(t, analytics.MemoryLeakDetected)
		assert.Greater(t, analytics.LeakRatePerHour, 0.0)
	})

	t.Run("Multiple storage device types", func(t *testing.T) {
		analytics, err := jm.GetJobIOAnalytics(ctx, "12345")
		require.NoError(t, err)
		
		// v0.0.42 supports multiple storage device types
		assert.Equal(t, 3, len(analytics.DeviceMetrics))
		
		// Check for different device types
		deviceTypes := make(map[string]bool)
		for _, device := range analytics.DeviceMetrics {
			deviceTypes[device.DeviceType] = true
		}
		assert.Equal(t, 3, len(deviceTypes))
		assert.True(t, deviceTypes["nvme"])
		assert.True(t, deviceTypes["ssd"])
		assert.True(t, deviceTypes["network"])
	})

	t.Run("GPU support", func(t *testing.T) {
		analytics, err := jm.GetJobComprehensiveAnalytics(ctx, "12345")
		require.NoError(t, err)
		
		// v0.0.42 has GPU efficiency metrics
		assert.Greater(t, analytics.EfficiencyMetrics.GPUEfficiency, 0.0)
		assert.Greater(t, analytics.EfficiencyMetrics.ResourceWaste["gpu_hours"], 0.0)
	})
}

// BenchmarkJobManager_GetJobCPUAnalytics_v42 benchmarks GetJobCPUAnalytics performance
func BenchmarkJobManager_GetJobCPUAnalytics_v42(b *testing.B) {
	ctx := context.Background()
	jm := &JobManager{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = jm.GetJobCPUAnalytics(ctx, "12345")
	}
}

// BenchmarkJobManager_GetJobComprehensiveAnalytics_v42 benchmarks GetJobComprehensiveAnalytics performance
func BenchmarkJobManager_GetJobComprehensiveAnalytics_v42(b *testing.B) {
	ctx := context.Background()
	jm := &JobManager{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = jm.GetJobComprehensiveAnalytics(ctx, "12345")
	}
}

// Helper function
func contains(s string, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || 
		   len(s) > len(substr) && containsHelper(s[1:], substr)
}

func containsHelper(s string, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	if s[:len(substr)] == substr {
		return true
	}
	return containsHelper(s[1:], substr)
}
