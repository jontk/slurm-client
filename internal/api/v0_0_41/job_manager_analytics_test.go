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