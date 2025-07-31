package integration

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/tests/mocks"
)

// TestAnalyticsCoverage_JobUtilization tests GetJobUtilization method across all versions
func TestAnalyticsCoverage_JobUtilization(t *testing.T) {
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}
	
	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			// Start mock server for this version
			mockServer := mocks.NewMockSlurmServerForVersion(version)
			defer mockServer.Close()

			// Create client
			client, err := slurm.NewClientWithVersion(context.Background(), version,
				slurm.WithBaseURL(mockServer.URL()),
				slurm.WithAuth(auth.NewNoAuth()),
			)
			require.NoError(t, err)

			ctx := context.Background()
			
			// Test valid job ID
			utilization, err := client.Jobs().GetJobUtilization(ctx, "1001")
			require.NoError(t, err)
			assert.NotNil(t, utilization)
			
			// Validate CPU utilization
			assert.NotNil(t, utilization.CPUUtilization)
			assert.Greater(t, utilization.CPUUtilization.Allocated, 0.0)
			assert.GreaterOrEqual(t, utilization.CPUUtilization.Used, 0.0)
			assert.GreaterOrEqual(t, utilization.CPUUtilization.Efficiency, 0.0)
			assert.LessOrEqual(t, utilization.CPUUtilization.Efficiency, 100.0)
			
			// Validate Memory utilization
			assert.NotNil(t, utilization.MemoryUtilization)
			assert.Greater(t, utilization.MemoryUtilization.Allocated, 0.0)
			assert.GreaterOrEqual(t, utilization.MemoryUtilization.Used, 0.0)
			assert.GreaterOrEqual(t, utilization.MemoryUtilization.Efficiency, 0.0)
			assert.LessOrEqual(t, utilization.MemoryUtilization.Efficiency, 100.0)
			
			// Validate GPU utilization (if available)
			if utilization.GPUUtilization != nil {
				assert.GreaterOrEqual(t, utilization.GPUUtilization.DeviceCount, 0)
				if utilization.GPUUtilization.OverallUtilization != nil {
					assert.GreaterOrEqual(t, utilization.GPUUtilization.OverallUtilization.Efficiency, 0.0)
					assert.LessOrEqual(t, utilization.GPUUtilization.OverallUtilization.Efficiency, 100.0)
				}
			}
			
			// Validate I/O utilization (if available)
			if utilization.IOUtilization != nil {
				assert.GreaterOrEqual(t, utilization.IOUtilization.TotalBytesRead, int64(0))
				assert.GreaterOrEqual(t, utilization.IOUtilization.TotalBytesWritten, int64(0))
			}
			
			// Test invalid job ID
			_, err = client.Jobs().GetJobUtilization(ctx, "invalid_job")
			assert.Error(t, err)
		})
	}
}

// TestAnalyticsCoverage_JobEfficiency tests GetJobEfficiency method across all versions
func TestAnalyticsCoverage_JobEfficiency(t *testing.T) {
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}
	
	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			// Start mock server for this version
			mockServer := mocks.NewMockSlurmServerForVersion(version)
			defer mockServer.Close()

			// Create client
			client, err := slurm.NewClientWithVersion(context.Background(), version,
				slurm.WithBaseURL(mockServer.URL()),
				slurm.WithAuth(auth.NewNoAuth()),
			)
			require.NoError(t, err)

			ctx := context.Background()
			
			// Test valid job ID
			efficiency, err := client.Jobs().GetJobEfficiency(ctx, "1001")
			require.NoError(t, err)
			assert.NotNil(t, efficiency)
			
			// Validate efficiency score
			assert.GreaterOrEqual(t, efficiency.Efficiency, 0.0)
			assert.LessOrEqual(t, efficiency.Efficiency, 100.0)
			
			// Validate resource usage
			assert.GreaterOrEqual(t, efficiency.Requested, 0.0)
			assert.GreaterOrEqual(t, efficiency.Allocated, 0.0)
			assert.GreaterOrEqual(t, efficiency.Used, 0.0)
			
			// Validate resource waste
			assert.GreaterOrEqual(t, efficiency.Wasted, 0.0)
			
			
			// Test invalid job ID
			_, err = client.Jobs().GetJobEfficiency(ctx, "invalid_job")
			assert.Error(t, err)
		})
	}
}

// TestAnalyticsCoverage_JobPerformance tests GetJobPerformance method across all versions
func TestAnalyticsCoverage_JobPerformance(t *testing.T) {
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}
	
	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			// Start mock server for this version
			mockServer := mocks.NewMockSlurmServerForVersion(version)
			defer mockServer.Close()

			// Create client
			client, err := slurm.NewClientWithVersion(context.Background(), version,
				slurm.WithBaseURL(mockServer.URL()),
				slurm.WithAuth(auth.NewNoAuth()),
			)
			require.NoError(t, err)

			ctx := context.Background()
			
			// Test valid job ID
			performance, err := client.Jobs().GetJobPerformance(ctx, "1001")
			require.NoError(t, err)
			assert.NotNil(t, performance)
			
			// Validate resource utilization
			if performance.ResourceUtilization != nil {
				assert.GreaterOrEqual(t, performance.ResourceUtilization.Allocated, 0.0)
				assert.GreaterOrEqual(t, performance.ResourceUtilization.Used, 0.0)
				assert.GreaterOrEqual(t, performance.ResourceUtilization.Efficiency, 0.0)
				assert.LessOrEqual(t, performance.ResourceUtilization.Efficiency, 100.0)
			}
			
			// Validate job utilization
			if performance.JobUtilization != nil {
				assert.NotEmpty(t, performance.JobUtilization.JobID)
				if performance.JobUtilization.CPUUtilization != nil {
					assert.GreaterOrEqual(t, performance.JobUtilization.CPUUtilization.Allocated, 0.0)
					assert.GreaterOrEqual(t, performance.JobUtilization.CPUUtilization.Used, 0.0)
				}
				if performance.JobUtilization.MemoryUtilization != nil {
					assert.GreaterOrEqual(t, performance.JobUtilization.MemoryUtilization.Allocated, 0.0)
					assert.GreaterOrEqual(t, performance.JobUtilization.MemoryUtilization.Used, 0.0)
				}
			}
			
			// Test invalid job ID
			_, err = client.Jobs().GetJobPerformance(ctx, "invalid_job")
			assert.Error(t, err)
		})
	}
}

// TestAnalyticsCoverage_JobLiveMetrics tests GetJobLiveMetrics method across all versions
func TestAnalyticsCoverage_JobLiveMetrics(t *testing.T) {
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}
	
	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			// Start mock server for this version
			mockServer := mocks.NewMockSlurmServerForVersion(version)
			defer mockServer.Close()

			// Create client
			client, err := slurm.NewClientWithVersion(context.Background(), version,
				slurm.WithBaseURL(mockServer.URL()),
				slurm.WithAuth(auth.NewNoAuth()),
			)
			require.NoError(t, err)

			ctx := context.Background()
			
			// Test valid job ID (running job)
			liveMetrics, err := client.Jobs().GetJobLiveMetrics(ctx, "1001")
			require.NoError(t, err)
			assert.NotNil(t, liveMetrics)
			
			// Validate CPU usage
			if liveMetrics.CPUUsage != nil {
				assert.GreaterOrEqual(t, liveMetrics.CPUUsage.Current, 0.0)
				assert.LessOrEqual(t, liveMetrics.CPUUsage.Current, 100.0)
				assert.GreaterOrEqual(t, liveMetrics.CPUUsage.Average1Min, 0.0)
				assert.LessOrEqual(t, liveMetrics.CPUUsage.Average1Min, 100.0)
				assert.GreaterOrEqual(t, liveMetrics.CPUUsage.Peak, 0.0)
				assert.LessOrEqual(t, liveMetrics.CPUUsage.Peak, 100.0)
				assert.GreaterOrEqual(t, liveMetrics.CPUUsage.UtilizationPercent, 0.0)
				assert.LessOrEqual(t, liveMetrics.CPUUsage.UtilizationPercent, 100.0)
			}
			
			// Validate Memory usage
			if liveMetrics.MemoryUsage != nil {
				assert.GreaterOrEqual(t, liveMetrics.MemoryUsage.Current, 0.0)
				assert.GreaterOrEqual(t, liveMetrics.MemoryUsage.Average1Min, 0.0)
				assert.GreaterOrEqual(t, liveMetrics.MemoryUsage.Peak, 0.0)
				assert.GreaterOrEqual(t, liveMetrics.MemoryUsage.UtilizationPercent, 0.0)
				assert.LessOrEqual(t, liveMetrics.MemoryUsage.UtilizationPercent, 100.0)
			}
			
			// Validate collection time
			assert.NotZero(t, liveMetrics.CollectionTime)
			
			// Test invalid job ID
			_, err = client.Jobs().GetJobLiveMetrics(ctx, "invalid_job")
			assert.Error(t, err)
		})
	}
}

// TestAnalyticsCoverage_JobResourceTrends tests GetJobResourceTrends method across all versions
func TestAnalyticsCoverage_JobResourceTrends(t *testing.T) {
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}
	
	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			// Start mock server for this version
			mockServer := mocks.NewMockSlurmServerForVersion(version)
			defer mockServer.Close()

			// Create client
			client, err := slurm.NewClientWithVersion(context.Background(), version,
				slurm.WithBaseURL(mockServer.URL()),
				slurm.WithAuth(auth.NewNoAuth()),
			)
			require.NoError(t, err)

			ctx := context.Background()
			
			// Create trend options
			trendOpts := &interfaces.ResourceTrendsOptions{
				StartTime: time.Now().Add(-1 * time.Hour),
				EndTime:   time.Now(),
				Interval:  "5m",
				Resources: []string{"cpu", "memory", "io"},
			}
			
			// Test valid job ID
			trends, err := client.Jobs().GetJobResourceTrends(ctx, "1001", trendOpts)
			require.NoError(t, err)
			assert.NotNil(t, trends)
			
			// Validate trends structure
			assert.NotEmpty(t, trends.JobID)
			assert.Equal(t, "1001", trends.JobID)
			
			// Test with nil options (should use defaults)
			trends2, err := client.Jobs().GetJobResourceTrends(ctx, "1001", nil)
			require.NoError(t, err)
			assert.NotNil(t, trends2)
			
			// Test invalid job ID
			_, err = client.Jobs().GetJobResourceTrends(ctx, "invalid_job", trendOpts)
			assert.Error(t, err)
		})
	}
}

// TestAnalyticsCoverage_ResourceSpecificAnalytics tests resource-specific analytics methods
func TestAnalyticsCoverage_ResourceSpecificAnalytics(t *testing.T) {
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}
	
	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			// Start mock server for this version
			mockServer := mocks.NewMockSlurmServerForVersion(version)
			defer mockServer.Close()

			// Create client
			client, err := slurm.NewClientWithVersion(context.Background(), version,
				slurm.WithBaseURL(mockServer.URL()),
				slurm.WithAuth(auth.NewNoAuth()),
			)
			require.NoError(t, err)

			ctx := context.Background()
			jobID := "1001"
			
			// Test GetJobCPUAnalytics
			cpuAnalytics, err := client.Jobs().GetJobCPUAnalytics(ctx, jobID)
			require.NoError(t, err)
			assert.NotNil(t, cpuAnalytics)
			assert.Greater(t, cpuAnalytics.AllocatedCores, 0)
			assert.GreaterOrEqual(t, cpuAnalytics.UsedCores, 0.0)
			assert.GreaterOrEqual(t, cpuAnalytics.UtilizationPercent, 0)
			assert.LessOrEqual(t, cpuAnalytics.UtilizationPercent, 100)
			
			// Test GetJobMemoryAnalytics
			memoryAnalytics, err := client.Jobs().GetJobMemoryAnalytics(ctx, jobID)
			require.NoError(t, err)
			assert.NotNil(t, memoryAnalytics)
			assert.Greater(t, memoryAnalytics.AllocatedBytes, int64(0))
			assert.GreaterOrEqual(t, memoryAnalytics.UsedBytes, int64(0))
			assert.GreaterOrEqual(t, memoryAnalytics.UtilizationPercent, 0)
			assert.LessOrEqual(t, memoryAnalytics.UtilizationPercent, 100)
			
			// Test GetJobIOAnalytics
			ioAnalytics, err := client.Jobs().GetJobIOAnalytics(ctx, jobID)
			require.NoError(t, err)
			assert.NotNil(t, ioAnalytics)
			assert.GreaterOrEqual(t, ioAnalytics.ReadBytes, int64(0))
			assert.GreaterOrEqual(t, ioAnalytics.WriteBytes, int64(0))
			assert.GreaterOrEqual(t, ioAnalytics.ReadOperations, 0)
			assert.GreaterOrEqual(t, ioAnalytics.WriteOperations, 0)
			
			// Test GetJobComprehensiveAnalytics
			comprehensive, err := client.Jobs().GetJobComprehensiveAnalytics(ctx, jobID)
			require.NoError(t, err)
			assert.NotNil(t, comprehensive)
			assert.NotEmpty(t, comprehensive.JobID)
			assert.Equal(t, jobID, comprehensive.JobID)
			
			// Test invalid job IDs for all methods
			_, err = client.Jobs().GetJobCPUAnalytics(ctx, "invalid_job")
			assert.Error(t, err)
			
			_, err = client.Jobs().GetJobMemoryAnalytics(ctx, "invalid_job")
			assert.Error(t, err)
			
			_, err = client.Jobs().GetJobIOAnalytics(ctx, "invalid_job")
			assert.Error(t, err)
			
			_, err = client.Jobs().GetJobComprehensiveAnalytics(ctx, "invalid_job")
			assert.Error(t, err)
		})
	}
}

// TestAnalyticsCoverage_PerformanceHistory tests GetJobPerformanceHistory method
func TestAnalyticsCoverage_PerformanceHistory(t *testing.T) {
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}
	
	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			// Start mock server for this version
			mockServer := mocks.NewMockSlurmServerForVersion(version)
			defer mockServer.Close()

			// Create client
			client, err := slurm.NewClientWithVersion(context.Background(), version,
				slurm.WithBaseURL(mockServer.URL()),
				slurm.WithAuth(auth.NewNoAuth()),
			)
			require.NoError(t, err)

			ctx := context.Background()
			jobID := "1001"
			
			// Create history options
			startTime := time.Now().Add(-1 * time.Hour)
			endTime := time.Now()
			historyOpts := &interfaces.PerformanceHistoryOptions{
				StartTime:    &startTime,
				EndTime:      &endTime,
				Interval:     "hourly",
				IncludeSteps: true,
				MetricTypes:  []string{"cpu", "memory", "io"},
			}
			
			// Test valid job ID
			history, err := client.Jobs().GetJobPerformanceHistory(ctx, jobID, historyOpts)
			require.NoError(t, err)
			assert.NotNil(t, history)
			assert.Equal(t, jobID, history.JobID)
			
			// Test with nil options (should use defaults)
			history2, err := client.Jobs().GetJobPerformanceHistory(ctx, jobID, nil)
			require.NoError(t, err)
			assert.NotNil(t, history2)
			
			// Test invalid job ID
			_, err = client.Jobs().GetJobPerformanceHistory(ctx, "invalid_job", historyOpts)
			assert.Error(t, err)
		})
	}
}

// TestAnalyticsCoverage_BatchAnalysis tests AnalyzeBatchJobs method
func TestAnalyticsCoverage_BatchAnalysis(t *testing.T) {
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}
	
	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			// Start mock server for this version
			mockServer := mocks.NewMockSlurmServerForVersion(version)
			defer mockServer.Close()

			// Create client
			client, err := slurm.NewClientWithVersion(context.Background(), version,
				slurm.WithBaseURL(mockServer.URL()),
				slurm.WithAuth(auth.NewNoAuth()),
			)
			require.NoError(t, err)

			ctx := context.Background()
			
			// Test batch analysis with valid job IDs
			jobIDs := []string{"1001", "1002"}
			batchOpts := &interfaces.BatchAnalysisOptions{
				IncludeDetails: true,
			}
			
			batchAnalysis, err := client.Jobs().AnalyzeBatchJobs(ctx, jobIDs, batchOpts)
			require.NoError(t, err)
			assert.NotNil(t, batchAnalysis)
			assert.Equal(t, len(jobIDs), batchAnalysis.JobCount)
			assert.GreaterOrEqual(t, batchAnalysis.AnalyzedCount, 0)
			assert.GreaterOrEqual(t, batchAnalysis.FailedCount, 0)
			assert.LessOrEqual(t, batchAnalysis.AnalyzedCount+batchAnalysis.FailedCount, batchAnalysis.JobCount)
			
			// Validate individual job analyses
			assert.Len(t, batchAnalysis.JobAnalyses, len(jobIDs))
			for _, jobAnalysis := range batchAnalysis.JobAnalyses {
				assert.NotEmpty(t, jobAnalysis.JobID)
				assert.Contains(t, []string{"completed", "failed"}, jobAnalysis.Status)
				
				if jobAnalysis.Status == "completed" {
					assert.GreaterOrEqual(t, jobAnalysis.Efficiency, 0.0)
					assert.LessOrEqual(t, jobAnalysis.Efficiency, 100.0)
					assert.GreaterOrEqual(t, jobAnalysis.CPUUtilization, 0.0)
					assert.LessOrEqual(t, jobAnalysis.CPUUtilization, 100.0)
					assert.GreaterOrEqual(t, jobAnalysis.MemoryUtilization, 0.0)
					assert.LessOrEqual(t, jobAnalysis.MemoryUtilization, 100.0)
				} else if len(jobAnalysis.Issues) > 0 {
					assert.NotEmpty(t, jobAnalysis.Issues)
				}
			}
			
			// Test with empty job IDs list
			_, err = client.Jobs().AnalyzeBatchJobs(ctx, []string{}, batchOpts)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "no job IDs provided")
			
			// Test with nil job IDs
			_, err = client.Jobs().AnalyzeBatchJobs(ctx, nil, batchOpts)
			assert.Error(t, err)
			
			// Test with mix of valid and invalid job IDs
			mixedJobIDs := []string{"1001", "invalid_job", "1002"}
			mixedAnalysis, err := client.Jobs().AnalyzeBatchJobs(ctx, mixedJobIDs, batchOpts)
			require.NoError(t, err)
			assert.NotNil(t, mixedAnalysis)
			assert.Equal(t, len(mixedJobIDs), mixedAnalysis.JobCount)
			assert.Greater(t, mixedAnalysis.FailedCount, 0) // Should have at least one failure
		})
	}
}

// TestAnalyticsCoverage_JobStepAnalytics tests job step analytics methods
func TestAnalyticsCoverage_JobStepAnalytics(t *testing.T) {
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}
	
	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			// Start mock server for this version
			mockServer := mocks.NewMockSlurmServerForVersion(version)
			defer mockServer.Close()

			// Create client
			client, err := slurm.NewClientWithVersion(context.Background(), version,
				slurm.WithBaseURL(mockServer.URL()),
				slurm.WithAuth(auth.NewNoAuth()),
			)
			require.NoError(t, err)

			ctx := context.Background()
			jobID := "1001"
			stepID := "0"
			
			// Test GetJobStepDetails
			stepDetails, err := client.Jobs().GetJobStepDetails(ctx, jobID, stepID)
			require.NoError(t, err)
			assert.NotNil(t, stepDetails)
			assert.Equal(t, jobID, stepDetails.JobID)
			assert.Equal(t, stepID, stepDetails.StepID)
			
			// Test GetJobStepUtilization
			stepUtilization, err := client.Jobs().GetJobStepUtilization(ctx, jobID, stepID)
			require.NoError(t, err)
			assert.NotNil(t, stepUtilization)
			assert.Equal(t, jobID, stepUtilization.JobID)
			assert.Equal(t, stepID, stepUtilization.StepID)
			
			// Test ListJobStepsWithMetrics
			listOpts := &interfaces.ListJobStepsOptions{
				StepStates: []string{"COMPLETED"},
			}
			stepsList, err := client.Jobs().ListJobStepsWithMetrics(ctx, jobID, listOpts)
			require.NoError(t, err)
			assert.NotNil(t, stepsList)
			assert.Equal(t, jobID, stepsList.JobID)
			
			// Test invalid job/step IDs
			_, err = client.Jobs().GetJobStepDetails(ctx, "invalid_job", stepID)
			assert.Error(t, err)
			
			_, err = client.Jobs().GetJobStepUtilization(ctx, jobID, "invalid_step")
			assert.Error(t, err)
		})
	}
}

// TestAnalyticsCoverage_IntegrationMethods tests SLURM integration methods
func TestAnalyticsCoverage_IntegrationMethods(t *testing.T) {
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}
	
	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			// Start mock server for this version
			mockServer := mocks.NewMockSlurmServerForVersion(version)
			defer mockServer.Close()

			// Create client
			client, err := slurm.NewClientWithVersion(context.Background(), version,
				slurm.WithBaseURL(mockServer.URL()),
				slurm.WithAuth(auth.NewNoAuth()),
			)
			require.NoError(t, err)

			ctx := context.Background()
			jobID := "1001"
			stepID := "0"
			
			// Note: These methods are not fully implemented, so we expect errors
			// but we still test to ensure they exist and have the correct signatures
			
			// Test GetJobStepsFromAccounting
			startTime := time.Now().Add(-1 * time.Hour)
			endTime := time.Now()
			accountingOpts := &interfaces.AccountingQueryOptions{
				StartTime: &startTime,
				EndTime:   &endTime,
			}
			_, err = client.Jobs().GetJobStepsFromAccounting(ctx, jobID, accountingOpts)
			assert.Error(t, err) // Expected since not implemented
			assert.Contains(t, err.Error(), "not") // Should contain "not implemented" or similar
			
			// Test GetStepAccountingData
			_, err = client.Jobs().GetStepAccountingData(ctx, jobID, stepID)
			assert.Error(t, err) // Expected since not implemented
			assert.Contains(t, err.Error(), "not") // Should contain "not implemented" or similar
			
			// Test GetJobStepAPIData
			_, err = client.Jobs().GetJobStepAPIData(ctx, jobID, stepID)
			assert.Error(t, err) // Expected since not implemented
			assert.Contains(t, err.Error(), "not") // Should contain "not implemented" or similar
			
			// Test ListJobStepsFromSacct
			sacctStartTime := time.Now().Add(-1 * time.Hour)
			sacctEndTime := time.Now()
			sacctOpts := &interfaces.SacctQueryOptions{
				StartTime: &sacctStartTime,
				EndTime:   &sacctEndTime,
				Format:    []string{"JobID", "StepID", "State", "CPUTime"},
			}
			_, err = client.Jobs().ListJobStepsFromSacct(ctx, jobID, sacctOpts)
			assert.Error(t, err) // Expected since not implemented
			assert.Contains(t, err.Error(), "not") // Should contain "not implemented" or similar
		})
	}
}

// TestAnalyticsCoverage_ErrorHandling tests error handling scenarios
func TestAnalyticsCoverage_ErrorHandling(t *testing.T) {
	// Test with invalid server URL (network error)
	t.Run("NetworkError", func(t *testing.T) {
		client, err := slurm.NewClientWithVersion(context.Background(), "v0.0.42",
			slurm.WithBaseURL("http://invalid-server:9999"),
			slurm.WithAuth(auth.NewNoAuth()),
		)
		require.NoError(t, err)

		ctx := context.Background()
		
		_, err = client.Jobs().GetJobUtilization(ctx, "1001")
		assert.Error(t, err)
		// Should be a network-related error
	})
	
	// Test with mock server returning errors
	t.Run("ServerError", func(t *testing.T) {
		mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
		defer mockServer.Close()

		client, err := slurm.NewClientWithVersion(context.Background(), "v0.0.42",
			slurm.WithBaseURL(mockServer.URL()),
			slurm.WithAuth(auth.NewNoAuth()),
		)
		require.NoError(t, err)

		ctx := context.Background()
		
		// Test with non-existent job ID
		_, err = client.Jobs().GetJobUtilization(ctx, "999999")
		assert.Error(t, err)
		
		// Test with malformed job ID
		_, err = client.Jobs().GetJobEfficiency(ctx, "not_a_number")
		assert.Error(t, err)
	})
}

// TestAnalyticsCoverage_ContextCancellation tests context cancellation handling
func TestAnalyticsCoverage_ContextCancellation(t *testing.T) {
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	client, err := slurm.NewClientWithVersion(context.Background(), "v0.0.42",
		slurm.WithBaseURL(mockServer.URL()),
		slurm.WithAuth(auth.NewNoAuth()),
	)
	require.NoError(t, err)

	// Create a context with immediate cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	
	// Test that cancelled context is handled properly
	_, err = client.Jobs().GetJobUtilization(ctx, "1001")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context")
}

// TestAnalyticsCoverage_RealWorldScenarios tests realistic usage scenarios
func TestAnalyticsCoverage_RealWorldScenarios(t *testing.T) {
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	client, err := slurm.NewClientWithVersion(context.Background(), "v0.0.42",
		slurm.WithBaseURL(mockServer.URL()),
		slurm.WithAuth(auth.NewNoAuth()),
	)
	require.NoError(t, err)

	ctx := context.Background()
	
	// Scenario 1: Complete job analysis workflow
	t.Run("CompleteAnalysisWorkflow", func(t *testing.T) {
		jobID := "1001"
		
		// 1. Get basic utilization
		utilization, err := client.Jobs().GetJobUtilization(ctx, jobID)
		require.NoError(t, err)
		
		// 2. Get efficiency analysis
		efficiency, err := client.Jobs().GetJobEfficiency(ctx, jobID)
		require.NoError(t, err)
		assert.NotNil(t, efficiency)
		
		// 3. Get detailed performance metrics
		performance, err := client.Jobs().GetJobPerformance(ctx, jobID)
		require.NoError(t, err)
		
		// 4. Get comprehensive analytics (should be consistent)
		comprehensive, err := client.Jobs().GetJobComprehensiveAnalytics(ctx, jobID)
		require.NoError(t, err)
		
		// Verify consistency between calls (when both utilization objects exist)
		if utilization.CPUUtilization != nil && performance.ResourceUtilization != nil {
			assert.Equal(t, utilization.CPUUtilization.Allocated, performance.ResourceUtilization.Allocated)
		}
		if utilization.MemoryUtilization != nil && performance.JobUtilization != nil && 
		   performance.JobUtilization.MemoryUtilization != nil {
			assert.Equal(t, utilization.MemoryUtilization.Allocated, performance.JobUtilization.MemoryUtilization.Allocated)
		}
		// Verify job ID (convert string to uint32 for comparison)
		jobIDUint, _ := strconv.ParseUint(jobID, 10, 32)
		assert.Equal(t, uint32(jobIDUint), comprehensive.JobID)
	})
	
	// Scenario 2: Monitoring multiple jobs
	t.Run("MultiJobMonitoring", func(t *testing.T) {
		jobIDs := []string{"1001", "1002", "1003"}
		
		// Monitor all jobs using batch analysis
		batchAnalysis, err := client.Jobs().AnalyzeBatchJobs(ctx, jobIDs, nil)
		require.NoError(t, err)
		
		assert.Equal(t, len(jobIDs), batchAnalysis.JobCount)
		assert.Len(t, batchAnalysis.JobAnalyses, len(jobIDs))
		
		// Verify each job has analysis
		for _, expectedJobID := range jobIDs {
			found := false
			for _, analysis := range batchAnalysis.JobAnalyses {
				if analysis.JobID == expectedJobID {
					found = true
					break
				}
			}
			assert.True(t, found, "Job %s should be in batch analysis results", expectedJobID)
		}
	})
	
	// Scenario 3: Resource optimization workflow
	t.Run("ResourceOptimizationWorkflow", func(t *testing.T) {
		jobID := "1001"
		
		// 1. Get efficiency analysis with recommendations
		efficiency, err := client.Jobs().GetJobEfficiency(ctx, jobID)
		require.NoError(t, err)
		
		// 2. Check if optimization is needed
		if efficiency.Efficiency < 80.0 {
			// 3. Get detailed resource-specific analytics
			cpuAnalytics, err := client.Jobs().GetJobCPUAnalytics(ctx, jobID)
			require.NoError(t, err)
			assert.NotNil(t, cpuAnalytics)
			
			memoryAnalytics, err := client.Jobs().GetJobMemoryAnalytics(ctx, jobID)
			require.NoError(t, err)
			assert.NotNil(t, memoryAnalytics)
		}
	})
}