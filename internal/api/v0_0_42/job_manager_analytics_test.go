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