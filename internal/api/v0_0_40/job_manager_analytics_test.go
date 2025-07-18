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

		manager := &JobManagerImpl{
			client: &WrapperClient{
				apiClient: &mockAPIClient{
					getJobResponse: mockJob,
				},
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

		manager := &JobManagerImpl{
			client: &WrapperClient{
				apiClient: &mockAPIClient{
					getJobResponse: mockJob,
				},
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

		manager := &JobManagerImpl{
			client: &WrapperClient{
				apiClient: &mockAPIClient{
					getJobResponse: mockJob,
				},
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

		manager := &JobManagerImpl{
			client: &WrapperClient{
				apiClient: &mockAPIClient{
					getJobResponse: mockJob,
				},
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

	manager := &JobManagerImpl{
		client: &WrapperClient{
			apiClient: &mockAPIClient{
				getJobResponse: mockJob,
			},
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

// mockAPIClient is a mock implementation for testing
type mockAPIClient struct {
	getJobResponse *interfaces.Job
	getJobError    error
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