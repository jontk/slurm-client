package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/tests/mocks"
)

// TestJobAnalyticsComprehensive tests comprehensive job analytics scenarios
func TestJobAnalyticsComprehensive(t *testing.T) {
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}

	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			// Setup mock server
			mockServer := mocks.NewMockSlurmServerForVersion(version)
			defer mockServer.Close()

			// Create client
			client, err := slurm.NewClientWithVersion(context.Background(), version,
				slurm.WithBaseURL(mockServer.URL()),
				slurm.WithAuth(auth.NewNoneAuth()),
			)
			require.NoError(t, err)

			// Run comprehensive analytics tests
			t.Run("ResourceUtilizationAnalytics", func(t *testing.T) {
				testResourceUtilizationAnalytics(t, client)
			})

			t.Run("LiveMetricsMonitoring", func(t *testing.T) {
				testLiveMetricsMonitoring(t, client)
			})

			t.Run("EfficiencyAnalysis", func(t *testing.T) {
				testEfficiencyAnalysis(t, client)
			})

			t.Run("PerformanceHistory", func(t *testing.T) {
				testPerformanceHistory(t, client)
			})

			t.Run("ResourceTrends", func(t *testing.T) {
				testResourceTrends(t, client)
			})

			t.Run("JobStepAnalytics", func(t *testing.T) {
				testJobStepAnalytics(t, client)
			})

			t.Run("BatchJobAnalysis", func(t *testing.T) {
				testBatchJobAnalysis(t, client)
			})

			t.Run("OptimizationRecommendations", func(t *testing.T) {
				testOptimizationRecommendations(t, client)
			})

			t.Run("WorkflowPerformance", func(t *testing.T) {
				testWorkflowPerformance(t, client)
			})

			t.Run("ErrorScenarios", func(t *testing.T) {
				testAnalyticsErrorScenarios(t, client)
			})
		})
	}
}

// testResourceUtilizationAnalytics tests resource utilization tracking
func testResourceUtilizationAnalytics(t *testing.T, client slurm.SlurmClient) {
	ctx := context.Background()

	// Test getting utilization for a running job
	jobID := "1001"
	utilization, err := client.Jobs().GetJobUtilization(ctx, jobID)
	require.NoError(t, err)
	assert.NotNil(t, utilization)

	// Validate CPU utilization
	assert.Greater(t, utilization.CPUUtilization.AllocatedCores, 0)
	assert.GreaterOrEqual(t, utilization.CPUUtilization.UsedCores, 0.0)
	assert.GreaterOrEqual(t, utilization.CPUUtilization.UtilizationPercent, 0)
	assert.LessOrEqual(t, utilization.CPUUtilization.UtilizationPercent, 100)

	// Validate Memory utilization
	assert.Greater(t, utilization.MemoryUtilization.AllocatedBytes, int64(0))
	assert.GreaterOrEqual(t, utilization.MemoryUtilization.UsedBytes, int64(0))
	assert.GreaterOrEqual(t, utilization.MemoryUtilization.UtilizationPercent, 0)
	assert.LessOrEqual(t, utilization.MemoryUtilization.UtilizationPercent, 100)

	// Validate GPU utilization (if available)
	if utilization.GPUUtilization.DeviceCount > 0 {
		assert.GreaterOrEqual(t, utilization.GPUUtilization.UtilizationPercent, 0)
		assert.LessOrEqual(t, utilization.GPUUtilization.UtilizationPercent, 100)
	}

	// Validate I/O utilization
	assert.GreaterOrEqual(t, utilization.IOUtilization.ReadBytes, int64(0))
	assert.GreaterOrEqual(t, utilization.IOUtilization.WriteBytes, int64(0))
	assert.GreaterOrEqual(t, utilization.IOUtilization.UtilizationPercent, 0)
}

// testLiveMetricsMonitoring tests real-time metrics monitoring
func testLiveMetricsMonitoring(t *testing.T, client slurm.SlurmClient) {
	ctx := context.Background()

	// Test getting live metrics for a running job
	jobID := "1001"
	liveMetrics, err := client.Jobs().GetJobLiveMetrics(ctx, jobID)
	require.NoError(t, err)
	assert.NotNil(t, liveMetrics)

	// Validate CPU metrics
	assert.GreaterOrEqual(t, liveMetrics.CPUUsage.Current, 0.0)
	assert.GreaterOrEqual(t, liveMetrics.CPUUsage.Average, 0.0)
	assert.GreaterOrEqual(t, liveMetrics.CPUUsage.Peak, 0.0)
	assert.GreaterOrEqual(t, liveMetrics.CPUUsage.Peak, liveMetrics.CPUUsage.Average)

	// Validate Memory metrics
	assert.GreaterOrEqual(t, liveMetrics.MemoryUsage.Current, int64(0))
	assert.GreaterOrEqual(t, liveMetrics.MemoryUsage.Average, int64(0))
	assert.GreaterOrEqual(t, liveMetrics.MemoryUsage.Peak, int64(0))
	assert.GreaterOrEqual(t, liveMetrics.MemoryUsage.Peak, liveMetrics.MemoryUsage.Average)

	// Validate timestamp
	assert.Greater(t, liveMetrics.Timestamp, int64(0))
	
	// Test watching metrics (if supported)
	// Note: This is a simplified test as full streaming requires more setup
	metricsOpts := &interfaces.WatchMetricsOptions{
		Interval: 1 * time.Second,
		Duration: 3 * time.Second,
	}
	
	eventChan, err := client.Jobs().WatchJobMetrics(ctx, jobID, metricsOpts)
	if err == nil && eventChan != nil {
		// Collect a few events
		events := []interfaces.JobMetricsEvent{}
		timeout := time.After(5 * time.Second)
		
		for i := 0; i < 3; i++ {
			select {
			case event, ok := <-eventChan:
				if !ok {
					break
				}
				events = append(events, event)
			case <-timeout:
				break
			}
		}
		
		// Validate we got some events
		assert.NotEmpty(t, events)
	}
}

// testEfficiencyAnalysis tests job efficiency calculations
func testEfficiencyAnalysis(t *testing.T, client slurm.SlurmClient) {
	ctx := context.Background()

	// Test getting efficiency metrics
	jobID := "1001"
	efficiency, err := client.Jobs().GetJobEfficiency(ctx, jobID)
	require.NoError(t, err)
	assert.NotNil(t, efficiency)

	// Validate overall efficiency score
	assert.GreaterOrEqual(t, efficiency.OverallEfficiencyScore, 0.0)
	assert.LessOrEqual(t, efficiency.OverallEfficiencyScore, 100.0)

	// Validate resource-specific efficiency
	assert.GreaterOrEqual(t, efficiency.CPUEfficiency, 0)
	assert.LessOrEqual(t, efficiency.CPUEfficiency, 100)
	assert.GreaterOrEqual(t, efficiency.MemoryEfficiency, 0)
	assert.LessOrEqual(t, efficiency.MemoryEfficiency, 100)

	// Validate resource waste calculations
	assert.GreaterOrEqual(t, efficiency.ResourceWaste.CPUCoreHours, 0.0)
	assert.GreaterOrEqual(t, efficiency.ResourceWaste.MemoryGBHours, 0.0)

	// Validate optimization recommendations
	assert.NotNil(t, efficiency.OptimizationRecommendations)
	for _, rec := range efficiency.OptimizationRecommendations {
		assert.NotEmpty(t, rec.Resource)
		assert.NotEmpty(t, rec.Type)
		assert.GreaterOrEqual(t, rec.Confidence, 0.0)
		assert.LessOrEqual(t, rec.Confidence, 1.0)
	}
}

// testPerformanceHistory tests historical performance data retrieval
func testPerformanceHistory(t *testing.T, client slurm.SlurmClient) {
	ctx := context.Background()

	// Test getting performance history
	jobID := "1001"
	historyOpts := &interfaces.PerformanceHistoryOptions{
		Timeframe:    "24h",
		Granularity:  "1h",
		IncludeSteps: true,
	}

	history, err := client.Jobs().GetJobPerformanceHistory(ctx, jobID, historyOpts)
	require.NoError(t, err)
	assert.NotNil(t, history)

	// Validate history structure
	assert.Equal(t, jobID, history.JobID)
	assert.NotEmpty(t, history.Timeframe)
	assert.NotEmpty(t, history.DataPoints)

	// Validate data points
	for _, dp := range history.DataPoints {
		assert.NotZero(t, dp.Timestamp)
		assert.GreaterOrEqual(t, dp.CPUUsage, 0.0)
		assert.LessOrEqual(t, dp.CPUUsage, 100.0)
		assert.GreaterOrEqual(t, dp.MemoryUsage, int64(0))
	}

	// Validate summary statistics
	assert.NotNil(t, history.Summary)
	assert.GreaterOrEqual(t, history.Summary.AverageEfficiency, 0.0)
	assert.GreaterOrEqual(t, history.Summary.PeakCPUUsage, 0.0)
	assert.GreaterOrEqual(t, history.Summary.PeakMemoryUsage, int64(0))
}

// testResourceTrends tests resource usage trend analysis
func testResourceTrends(t *testing.T, client slurm.SlurmClient) {
	ctx := context.Background()

	// Test getting resource trends
	jobID := "1001"
	trendOpts := &interfaces.ResourceTrendsOptions{
		StartTime: time.Now().Add(-24 * time.Hour),
		EndTime:   time.Now(),
		Interval:  "1h",
		Resources: []string{"cpu", "memory", "io"},
	}

	trends, err := client.Jobs().GetJobResourceTrends(ctx, jobID, trendOpts)
	require.NoError(t, err)
	assert.NotNil(t, trends)

	// Validate trends structure
	assert.Equal(t, jobID, trends.JobID)
	assert.NotEmpty(t, trends.TimeRange)
	assert.NotEmpty(t, trends.Interval)

	// Validate CPU trends
	assert.NotNil(t, trends.CPUTrends)
	assert.NotEmpty(t, trends.CPUTrends.DataPoints)
	for _, dp := range trends.CPUTrends.DataPoints {
		assert.NotZero(t, dp.Timestamp)
		assert.GreaterOrEqual(t, dp.Value, 0.0)
	}

	// Validate Memory trends
	assert.NotNil(t, trends.MemoryTrends)
	assert.NotEmpty(t, trends.MemoryTrends.DataPoints)

	// Validate trend analysis
	assert.NotEmpty(t, trends.CPUTrends.Trend)
	assert.NotZero(t, trends.CPUTrends.Average)
}

// testJobStepAnalytics tests job step-level analytics
func testJobStepAnalytics(t *testing.T, client slurm.SlurmClient) {
	ctx := context.Background()

	jobID := "1001"
	stepID := "0"

	// Test getting step details
	stepDetails, err := client.Jobs().GetJobStepDetails(ctx, jobID, stepID)
	require.NoError(t, err)
	assert.NotNil(t, stepDetails)
	assert.Equal(t, jobID, stepDetails.JobID)
	assert.Equal(t, stepID, stepDetails.StepID)

	// Test getting step utilization
	stepUtilization, err := client.Jobs().GetJobStepUtilization(ctx, jobID, stepID)
	require.NoError(t, err)
	assert.NotNil(t, stepUtilization)
	assert.Equal(t, jobID, stepUtilization.JobID)
	assert.Equal(t, stepID, stepUtilization.StepID)

	// Test listing steps with metrics
	listOpts := &interfaces.ListJobStepsOptions{
		IncludeMetrics: true,
	}
	stepsList, err := client.Jobs().ListJobStepsWithMetrics(ctx, jobID, listOpts)
	require.NoError(t, err)
	assert.NotNil(t, stepsList)
	assert.NotEmpty(t, stepsList.Steps)

	// Validate step metrics
	for _, step := range stepsList.Steps {
		assert.NotEmpty(t, step.StepID)
		if step.Metrics != nil {
			assert.GreaterOrEqual(t, step.Metrics.CPUTime, 0.0)
			assert.GreaterOrEqual(t, step.Metrics.MemoryUsed, int64(0))
		}
	}
}

// testBatchJobAnalysis tests batch job analysis functionality
func testBatchJobAnalysis(t *testing.T, client slurm.SlurmClient) {
	ctx := context.Background()

	// Test analyzing multiple jobs
	jobIDs := []string{"1001", "1002", "1003"}
	batchOpts := &interfaces.BatchAnalysisOptions{
		IncludeDetails: true,
		Parallel:       true,
	}

	batchAnalysis, err := client.Jobs().AnalyzeBatchJobs(ctx, jobIDs, batchOpts)
	require.NoError(t, err)
	assert.NotNil(t, batchAnalysis)

	// Validate batch analysis results
	assert.GreaterOrEqual(t, batchAnalysis.JobCount, len(jobIDs))
	assert.GreaterOrEqual(t, batchAnalysis.AnalyzedCount, 0)
	assert.NotNil(t, batchAnalysis.AggregateStats)

	// Validate aggregate statistics
	stats := batchAnalysis.AggregateStats
	assert.GreaterOrEqual(t, stats.AverageEfficiency, 0.0)
	assert.GreaterOrEqual(t, stats.TotalCPUTime, 0.0)
	assert.GreaterOrEqual(t, stats.TotalMemoryUsed, int64(0))

	// Validate individual job analyses
	if batchOpts.IncludeDetails {
		assert.NotEmpty(t, batchAnalysis.JobAnalyses)
		for _, jobAnalysis := range batchAnalysis.JobAnalyses {
			assert.NotEmpty(t, jobAnalysis.JobID)
			assert.GreaterOrEqual(t, jobAnalysis.EfficiencyScore, 0.0)
		}
	}
}

// testOptimizationRecommendations tests optimization recommendation generation
func testOptimizationRecommendations(t *testing.T, client slurm.SlurmClient) {
	ctx := context.Background()

	// Get efficiency with recommendations
	jobID := "1001"
	efficiency, err := client.Jobs().GetJobEfficiency(ctx, jobID)
	require.NoError(t, err)
	assert.NotNil(t, efficiency)

	// Test recommendations
	if len(efficiency.OptimizationRecommendations) > 0 {
		for _, rec := range efficiency.OptimizationRecommendations {
			// Validate recommendation structure
			assert.NotEmpty(t, rec.Type)
			assert.NotEmpty(t, rec.Resource)
			assert.NotEmpty(t, rec.Reason)
			
			// Validate recommendation values
			assert.GreaterOrEqual(t, rec.Current, 0)
			assert.GreaterOrEqual(t, rec.Recommended, 0)
			assert.GreaterOrEqual(t, rec.Confidence, 0.0)
			assert.LessOrEqual(t, rec.Confidence, 1.0)

			// Validate recommendation logic
			switch rec.Type {
			case "reduction":
				assert.Less(t, rec.Recommended, rec.Current, 
					"Reduction recommendation should suggest lower value")
			case "increase":
				assert.Greater(t, rec.Recommended, rec.Current,
					"Increase recommendation should suggest higher value")
			}
		}
	}

	// Test getting similar jobs for performance comparison
	similarJobs, err := client.Jobs().GetSimilarJobsPerformance(ctx, jobID)
	if err == nil && similarJobs != nil {
		assert.NotEmpty(t, similarJobs.JobID)
		assert.NotEmpty(t, similarJobs.SimilarJobs)
		
		for _, similar := range similarJobs.SimilarJobs {
			assert.NotEmpty(t, similar.JobID)
			assert.GreaterOrEqual(t, similar.Similarity, 0.0)
			assert.LessOrEqual(t, similar.Similarity, 1.0)
		}
	}
}

// testWorkflowPerformance tests workflow performance analysis
func testWorkflowPerformance(t *testing.T, client slurm.SlurmClient) {
	ctx := context.Background()

	// Test workflow performance analysis
	workflowOpts := &interfaces.WorkflowAnalysisOptions{
		WorkflowID: "workflow-123",
		JobIDs:     []string{"1001", "1002", "1003"},
		IncludeDependencies: true,
	}

	workflowPerf, err := client.Jobs().GetWorkflowPerformance(ctx, workflowOpts)
	if err == nil && workflowPerf != nil {
		// Validate workflow performance data
		assert.Equal(t, workflowOpts.WorkflowID, workflowPerf.WorkflowID)
		assert.NotEmpty(t, workflowPerf.Jobs)
		assert.GreaterOrEqual(t, workflowPerf.TotalDuration, time.Duration(0))
		assert.GreaterOrEqual(t, workflowPerf.TotalEfficiency, 0.0)

		// Validate critical path analysis
		if workflowPerf.CriticalPath != nil {
			assert.NotEmpty(t, workflowPerf.CriticalPath.Jobs)
			assert.GreaterOrEqual(t, workflowPerf.CriticalPath.Duration, time.Duration(0))
		}
	}
}

// testAnalyticsErrorScenarios tests error handling in analytics
func testAnalyticsErrorScenarios(t *testing.T, client slurm.SlurmClient) {
	ctx := context.Background()

	// Test with non-existent job
	_, err := client.Jobs().GetJobUtilization(ctx, "99999")
	assert.Error(t, err)

	// Test with invalid job ID
	_, err = client.Jobs().GetJobEfficiency(ctx, "invalid-id")
	assert.Error(t, err)

	// Test with cancelled context
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()
	_, err = client.Jobs().GetJobPerformance(cancelCtx, "1001")
	assert.Error(t, err)

	// Test with empty batch analysis
	_, err = client.Jobs().AnalyzeBatchJobs(ctx, []string{}, nil)
	assert.Error(t, err)

	// Test with invalid time range for trends
	invalidOpts := &interfaces.ResourceTrendsOptions{
		StartTime: time.Now().Add(24 * time.Hour), // Future start time
		EndTime:   time.Now(),
		Interval:  "1h",
	}
	_, err = client.Jobs().GetJobResourceTrends(ctx, "1001", invalidOpts)
	assert.Error(t, err)
}

// TestJobAnalyticsPerformanceOverhead tests analytics collection overhead
func TestJobAnalyticsPerformanceOverhead(t *testing.T) {
	// Setup mock server
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	// Create client
	client, err := slurm.NewClientWithVersion(context.Background(), "v0.0.42",
		slurm.WithBaseURL(mockServer.URL()),
		slurm.WithAuth(auth.NewNoneAuth()),
	)
	require.NoError(t, err)

	ctx := context.Background()
	jobID := "1001"

	// Measure baseline job retrieval time
	baselineStart := time.Now()
	for i := 0; i < 100; i++ {
		_, err := client.Jobs().Get(ctx, jobID)
		require.NoError(t, err)
	}
	baselineDuration := time.Since(baselineStart)
	baselineAvg := baselineDuration / 100

	// Measure analytics collection time
	analyticsStart := time.Now()
	for i := 0; i < 100; i++ {
		_, err := client.Jobs().GetJobUtilization(ctx, jobID)
		require.NoError(t, err)
	}
	analyticsDuration := time.Since(analyticsStart)
	analyticsAvg := analyticsDuration / 100

	// Calculate overhead percentage
	overhead := float64(analyticsAvg-baselineAvg) / float64(baselineAvg) * 100

	t.Logf("Baseline average: %v", baselineAvg)
	t.Logf("Analytics average: %v", analyticsAvg)
	t.Logf("Overhead: %.2f%%", overhead)

	// Assert overhead is less than 5%
	assert.Less(t, overhead, 5.0, "Analytics overhead should be less than 5%")
}

// TestJobAnalyticsConcurrency tests concurrent analytics requests
func TestJobAnalyticsConcurrency(t *testing.T) {
	// Setup mock server
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	// Create client
	client, err := slurm.NewClientWithVersion(context.Background(), "v0.0.42",
		slurm.WithBaseURL(mockServer.URL()),
		slurm.WithAuth(auth.NewNoneAuth()),
	)
	require.NoError(t, err)

	ctx := context.Background()
	jobIDs := []string{"1001", "1002", "1003", "1004", "1005"}

	// Test concurrent analytics requests
	concurrency := 10
	errChan := make(chan error, concurrency*len(jobIDs))

	for i := 0; i < concurrency; i++ {
		go func() {
			for _, jobID := range jobIDs {
				// Test various analytics endpoints concurrently
				_, err := client.Jobs().GetJobUtilization(ctx, jobID)
				if err != nil {
					errChan <- fmt.Errorf("utilization error for job %s: %w", jobID, err)
					return
				}

				_, err = client.Jobs().GetJobEfficiency(ctx, jobID)
				if err != nil {
					errChan <- fmt.Errorf("efficiency error for job %s: %w", jobID, err)
					return
				}

				_, err = client.Jobs().GetJobPerformance(ctx, jobID)
				if err != nil {
					errChan <- fmt.Errorf("performance error for job %s: %w", jobID, err)
					return
				}
			}
			errChan <- nil
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < concurrency; i++ {
		err := <-errChan
		assert.NoError(t, err)
	}
}

// TestJobAnalyticsDataConsistency tests data consistency across analytics endpoints
func TestJobAnalyticsDataConsistency(t *testing.T) {
	// Setup mock server
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	// Create client
	client, err := slurm.NewClientWithVersion(context.Background(), "v0.0.42",
		slurm.WithBaseURL(mockServer.URL()),
		slurm.WithAuth(auth.NewNoneAuth()),
	)
	require.NoError(t, err)

	ctx := context.Background()
	jobID := "1001"

	// Get data from different endpoints
	utilization, err := client.Jobs().GetJobUtilization(ctx, jobID)
	require.NoError(t, err)

	performance, err := client.Jobs().GetJobPerformance(ctx, jobID)
	require.NoError(t, err)

	efficiency, err := client.Jobs().GetJobEfficiency(ctx, jobID)
	require.NoError(t, err)

	// Verify data consistency across endpoints
	// CPU data should be consistent
	assert.Equal(t, utilization.CPUUtilization.AllocatedCores, 
		performance.CPUAnalytics.AllocatedCores,
		"CPU allocated cores should be consistent")

	// Memory data should be consistent
	assert.Equal(t, utilization.MemoryUtilization.AllocatedBytes,
		performance.MemoryAnalytics.AllocatedBytes,
		"Memory allocation should be consistent")

	// Efficiency calculations should be based on utilization
	if utilization.CPUUtilization.UtilizationPercent > 0 {
		assert.Greater(t, efficiency.CPUEfficiency, 0,
			"CPU efficiency should be greater than 0 when utilization exists")
	}
}