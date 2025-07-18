package v0_0_43

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

// TestGetJobUtilization tests the GetJobUtilization method
func TestGetJobUtilization(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Create mock job manager
		mockJob := &interfaces.Job{
			ID:        "12345",
			Name:      "test-job",
			State:     "RUNNING",
			CPUs:      16,
			Memory:    64 * 1024 * 1024 * 1024, // 64GB
			Partition: "compute",
			Nodes:     []string{"node001", "node002"},
			SubmitTime: time.Now().Add(-2 * time.Hour),
			StartTime:  &[]time.Time{time.Now().Add(-1 * time.Hour)}[0],
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

		// Test GetJobUtilization
		utilization, err := manager.GetJobUtilization(ctx, "12345")
		require.NoError(t, err)
		require.NotNil(t, utilization)

		// Verify basic fields
		assert.Equal(t, "12345", utilization.JobID)
		assert.Equal(t, "test-job", utilization.JobName)
		assert.Equal(t, mockJob.SubmitTime, utilization.StartTime)

		// Verify CPU utilization
		assert.NotNil(t, utilization.CPUUtilization)
		assert.Equal(t, 85.0, utilization.CPUUtilization.Percentage)
		assert.Equal(t, float64(16), utilization.CPUUtilization.Allocated)
		assert.Equal(t, float64(16)*0.85, utilization.CPUUtilization.Used)

		// Verify memory utilization
		assert.NotNil(t, utilization.MemoryUtilization)
		assert.Equal(t, 72.0, utilization.MemoryUtilization.Percentage)
		assert.Equal(t, float64(mockJob.Memory), utilization.MemoryUtilization.Allocated)

		// Verify GPU utilization
		assert.NotNil(t, utilization.GPUUtilization)
		assert.Equal(t, 2, utilization.GPUUtilization.TotalGPUs)
		assert.NotNil(t, utilization.GPUUtilization.AverageUtilization)
		assert.Equal(t, 90.0, utilization.GPUUtilization.AverageUtilization.Percentage)
		assert.Len(t, utilization.GPUUtilization.GPUs, 2)

		// Verify I/O utilization
		assert.NotNil(t, utilization.IOUtilization)
		assert.NotNil(t, utilization.IOUtilization.ReadBandwidth)
		assert.Equal(t, 20.0, utilization.IOUtilization.ReadBandwidth.Percentage)

		// Verify network utilization
		assert.NotNil(t, utilization.NetworkUtilization)
		assert.NotNil(t, utilization.NetworkUtilization.TotalBandwidth)
		assert.Equal(t, 10.0, utilization.NetworkUtilization.TotalBandwidth.Percentage)

		// Verify energy usage (only if job ended)
		if mockJob.EndTime != nil {
			assert.NotNil(t, utilization.EnergyUsage)
		}

		// Verify metadata
		assert.Equal(t, "v0.0.43", utilization.Metadata["version"])
		assert.Equal(t, "simulated", utilization.Metadata["source"])
	})

	t.Run("ClientNotInitialized", func(t *testing.T) {
		manager := &JobManagerImpl{
			client: &WrapperClient{},
		}

		utilization, err := manager.GetJobUtilization(ctx, "12345")
		assert.Error(t, err)
		assert.Nil(t, utilization)
		assert.True(t, errors.IsClientError(err))
	})

	t.Run("JobNotFound", func(t *testing.T) {
		manager := &JobManagerImpl{
			client: &WrapperClient{
				apiClient: &mockAPIClient{
					getJobError: errors.NewClientError(errors.ErrorCodeResourceNotFound, "Job not found"),
				},
			},
		}

		utilization, err := manager.GetJobUtilization(ctx, "99999")
		assert.Error(t, err)
		assert.Nil(t, utilization)
	})
}

// TestGetJobEfficiency tests the GetJobEfficiency method
func TestGetJobEfficiency(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockJob := &interfaces.Job{
			ID:        "12345",
			Name:      "test-job",
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

		efficiency, err := manager.GetJobEfficiency(ctx, "12345")
		require.NoError(t, err)
		require.NotNil(t, efficiency)

		// Verify efficiency calculation
		assert.Greater(t, efficiency.Percentage, 0.0)
		assert.LessOrEqual(t, efficiency.Percentage, 100.0)

		// Verify metadata
		metadata := efficiency.Metadata
		assert.NotNil(t, metadata["cpu_efficiency"])
		assert.NotNil(t, metadata["memory_efficiency"])
		assert.Equal(t, "weighted_average", metadata["calculation_method"])
		assert.NotNil(t, metadata["weights"])
	})

	t.Run("WithGPU", func(t *testing.T) {
		mockJob := &interfaces.Job{
			ID:        "12346",
			Name:      "gpu-job",
			State:     "COMPLETED",
			CPUs:      8,
			Memory:    32 * 1024 * 1024 * 1024,
			Partition: "gpu",
			SubmitTime: time.Now().Add(-1 * time.Hour),
			EndTime:    &[]time.Time{time.Now()}[0],
			Metadata: map[string]interface{}{
				"gpu_count": 4,
			},
		}

		manager := &JobManagerImpl{
			client: &WrapperClient{
				apiClient: &mockAPIClient{
					getJobResponse: mockJob,
				},
			},
		}

		efficiency, err := manager.GetJobEfficiency(ctx, "12346")
		require.NoError(t, err)
		require.NotNil(t, efficiency)

		// Verify GPU is included in efficiency
		weights := efficiency.Metadata["weights"].(map[string]float64)
		assert.Equal(t, 0.2, weights["gpu"])
	})
}

// TestGetJobPerformance tests the GetJobPerformance method
func TestGetJobPerformance(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		startTime := time.Now().Add(-3 * time.Hour)
		endTime := time.Now().Add(-30 * time.Minute)
		mockJob := &interfaces.Job{
			ID:         "12347",
			Name:       "perf-test-job",
			State:      "COMPLETED",
			CPUs:       32,
			Memory:     128 * 1024 * 1024 * 1024,
			Partition:  "compute",
			SubmitTime: time.Now().Add(-4 * time.Hour),
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

		performance, err := manager.GetJobPerformance(ctx, "12347")
		require.NoError(t, err)
		require.NotNil(t, performance)

		// Verify basic fields
		assert.Equal(t, uint32(12347), performance.JobID)
		assert.Equal(t, "perf-test-job", performance.JobName)
		assert.Equal(t, "COMPLETED", performance.Status)
		assert.Equal(t, 0, performance.ExitCode)

		// Verify utilization and efficiency are included
		assert.NotNil(t, performance.ResourceUtilization)
		assert.NotNil(t, performance.JobUtilization)

		// Verify performance trends
		assert.NotNil(t, performance.PerformanceTrends)
		assert.NotEmpty(t, performance.PerformanceTrends.TimePoints)
		assert.NotEmpty(t, performance.PerformanceTrends.CPUTrends)
		assert.Equal(t, len(performance.PerformanceTrends.TimePoints), len(performance.PerformanceTrends.CPUTrends))

		// Verify bottlenecks
		assert.NotNil(t, performance.Bottlenecks)
		// May or may not have bottlenecks depending on simulated data

		// Verify recommendations
		assert.NotNil(t, performance.Recommendations)
		// Should have at least some recommendations
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

// TestHelperFunctions tests the helper functions
func TestHelperFunctions(t *testing.T) {
	t.Run("generatePerformanceTrends", func(t *testing.T) {
		startTime := time.Now().Add(-5 * time.Hour)
		endTime := time.Now()
		job := &interfaces.Job{
			StartTime: &startTime,
			EndTime:   &endTime,
		}

		trends := generatePerformanceTrends(job)
		assert.NotNil(t, trends)
		assert.NotEmpty(t, trends.TimePoints)
		assert.Greater(t, len(trends.TimePoints), 0)
		assert.LessOrEqual(t, len(trends.TimePoints), 24)

		// Verify all trend arrays have same length
		assert.Equal(t, len(trends.TimePoints), len(trends.CPUTrends))
		assert.Equal(t, len(trends.TimePoints), len(trends.MemoryTrends))
		assert.Equal(t, len(trends.TimePoints), len(trends.IOTrends))
		assert.Equal(t, len(trends.TimePoints), len(trends.NetworkTrends))
	})

	t.Run("analyzeBottlenecks", func(t *testing.T) {
		// Test CPU bottleneck
		utilization := &interfaces.JobUtilization{
			CPUUtilization: &interfaces.ResourceUtilization{
				Percentage: 96.0,
			},
			MemoryUtilization: &interfaces.ResourceUtilization{
				Percentage: 70.0,
			},
		}

		bottlenecks := analyzeBottlenecks(utilization)
		assert.NotEmpty(t, bottlenecks)
		
		// Find CPU bottleneck
		var cpuBottleneck *interfaces.PerformanceBottleneck
		for i := range bottlenecks {
			if bottlenecks[i].Type == "cpu" {
				cpuBottleneck = &bottlenecks[i]
				break
			}
		}
		assert.NotNil(t, cpuBottleneck)
		assert.Equal(t, "high", cpuBottleneck.Severity)
	})

	t.Run("generateRecommendations", func(t *testing.T) {
		// Test low CPU utilization
		utilization := &interfaces.JobUtilization{
			CPUUtilization: &interfaces.ResourceUtilization{
				Percentage: 30.0,
			},
			MemoryUtilization: &interfaces.ResourceUtilization{
				Percentage: 80.0,
			},
		}
		efficiency := &interfaces.ResourceUtilization{
			Percentage: 55.0,
		}

		recommendations := generateRecommendations(utilization, efficiency)
		assert.NotEmpty(t, recommendations)

		// Should have CPU reduction recommendation
		var cpuRecommendation *interfaces.OptimizationRecommendation
		for i := range recommendations {
			if recommendations[i].Type == "resource_adjustment" && 
			   recommendations[i].Title == "Reduce CPU allocation" {
				cpuRecommendation = &recommendations[i]
				break
			}
		}
		assert.NotNil(t, cpuRecommendation)
		assert.Equal(t, "high", cpuRecommendation.Priority)
	})
}

// mockAPIClient is a mock implementation for testing
type mockAPIClient struct {
	getJobResponse *interfaces.Job
	getJobError    error
}

func (m *mockAPIClient) SlurmV0043GetJobWithResponse(ctx context.Context, jobID string, params *SlurmV0043GetJobParams) (*SlurmV0043GetJobResponse, error) {
	if m.getJobError != nil {
		return nil, m.getJobError
	}

	// Convert jobID to int32
	jobIDInt, _ := strconv.ParseInt(jobID, 10, 32)
	jobIDInt32 := int32(jobIDInt)

	// Mock response
	resp := &SlurmV0043GetJobResponse{
		JSON200: &V0043OpenapiJobInfoResp{
			Jobs: []V0043JobInfo{
				{
					JobId:     &jobIDInt32,
					Name:      &m.getJobResponse.Name,
					JobState:  &[]string{m.getJobResponse.State},
					Partition: &m.getJobResponse.Partition,
					Cpus: &V0043Uint32NoValStruct{
						Number: &[]int32{int32(m.getJobResponse.CPUs)}[0],
						Set:    &[]bool{true}[0],
					},
					MemoryPerNode: &V0043Uint64NoValStruct{
						Number: &[]int64{m.getJobResponse.Memory / (1024 * 1024)}[0], // Convert to MB
						Set:    &[]bool{true}[0],
					},
				},
			},
		},
	}

	// Add metadata
	if gpuCount, ok := m.getJobResponse.Metadata["gpu_count"].(int); ok {
		// In real implementation, this would be in a different field
		resp.JSON200.Jobs[0].AdminComment = &[]string{"gpu_count:" + strconv.Itoa(gpuCount)}[0]
	}

	return resp, nil
}

func (m *mockAPIClient) SlurmV0043GetJobsWithResponse(ctx context.Context, params *SlurmV0043GetJobsParams) (*SlurmV0043GetJobsResponse, error) {
	// Not used in these tests
	return nil, nil
}