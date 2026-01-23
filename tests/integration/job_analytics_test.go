// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jontk/slurm-client/tests/mocks"
)

// TestJobAnalyticsMockServer tests all job analytics mock server endpoints across API versions
func TestJobAnalyticsMockServer(t *testing.T) {
	testCases := []struct {
		name       string
		apiVersion string
	}{
		{"v0.0.40", "v0.0.40"},
		{"v0.0.41", "v0.0.41"},
		{"v0.0.42", "v0.0.42"},
		{"v0.0.43", "v0.0.43"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testJobAnalyticsEndpointsForVersion(t, tc.apiVersion)
		})
	}
}

func testJobAnalyticsEndpointsForVersion(t *testing.T, apiVersion string) {
	// Setup mock server for the specific API version
	mockServer := mocks.NewMockSlurmServerForVersion(apiVersion)
	defer mockServer.Close()

	baseURL := mockServer.URL()
	jobID := "1001" // Use existing job from mock data
	stepID := "0"   // Default step

	// Test basic analytics endpoints that should be available in all versions
	t.Run("BasicAnalytics", func(t *testing.T) {
		testBasicAnalyticsEndpoints(t, baseURL, apiVersion, jobID)
	})

	// Test advanced analytics endpoints (for supported versions)
	if mockServer.GetConfig().SupportedOperations["jobs.live_metrics"] {
		t.Run("AdvancedAnalytics", func(t *testing.T) {
			testAdvancedAnalyticsEndpoints(t, baseURL, apiVersion, jobID, stepID)
		})
	}

	// Test historical analytics endpoints (for supported versions)
	if mockServer.GetConfig().SupportedOperations["jobs.performance_history"] {
		t.Run("HistoricalAnalytics", func(t *testing.T) {
			testHistoricalAnalyticsEndpoints(t, baseURL, apiVersion)
		})
	}

	// Test comparative analytics endpoints (for supported versions)
	if mockServer.GetConfig().SupportedOperations["jobs.compare_performance"] {
		t.Run("ComparativeAnalytics", func(t *testing.T) {
			testComparativeAnalyticsEndpoints(t, baseURL, apiVersion, jobID)
		})
	}
}

func testBasicAnalyticsEndpoints(t *testing.T, baseURL, apiVersion, jobID string) {
	// Test job utilization endpoint
	t.Run("JobUtilization", func(t *testing.T) {
		url := fmt.Sprintf("%s/slurm/%s/job/%s/utilization", baseURL, apiVersion, jobID)
		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		var utilization map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&utilization)
		require.NoError(t, err)

		// Verify response structure - mock server nests data under "utilization"
		assert.Contains(t, utilization, "utilization")
		utilizationData := utilization["utilization"].(map[string]interface{})

		assert.Equal(t, jobID, utilizationData["job_id"])
		assert.Contains(t, utilizationData, "cpu_utilization")
		assert.Contains(t, utilizationData, "memory_utilization")

		// Check CPU utilization structure
		cpuUtil := utilizationData["cpu_utilization"].(map[string]interface{})
		assert.Contains(t, cpuUtil, "allocated_cores")
		assert.Contains(t, cpuUtil, "used_cores")
		assert.Contains(t, cpuUtil, "utilization_percent")

		// Check Memory utilization structure
		memUtil := utilizationData["memory_utilization"].(map[string]interface{})
		assert.Contains(t, memUtil, "allocated_bytes")
		assert.Contains(t, memUtil, "used_bytes")
		assert.Contains(t, memUtil, "utilization_percent")

		t.Logf("Job %s utilization response received successfully", jobID)
	})

	// Test job efficiency endpoint
	t.Run("JobEfficiency", func(t *testing.T) {
		url := fmt.Sprintf("%s/slurm/%s/job/%s/efficiency", baseURL, apiVersion, jobID)
		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var efficiency map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&efficiency)
		require.NoError(t, err)

		// Verify response structure - mock server nests data under "efficiency"
		assert.Contains(t, efficiency, "efficiency")
		efficiencyData := efficiency["efficiency"].(map[string]interface{})

		assert.Equal(t, jobID, efficiencyData["job_id"])
		assert.Contains(t, efficiencyData, "overall_efficiency_score")
		assert.Contains(t, efficiencyData, "cpu_efficiency")
		assert.Contains(t, efficiencyData, "memory_efficiency")
		assert.Contains(t, efficiencyData, "resource_waste")
		assert.Contains(t, efficiencyData, "optimization_recommendations")

		t.Logf("Job %s efficiency response received successfully", jobID)
	})

	// Test job performance endpoint
	t.Run("JobPerformance", func(t *testing.T) {
		url := fmt.Sprintf("%s/slurm/%s/job/%s/performance", baseURL, apiVersion, jobID)
		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var performance map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&performance)
		require.NoError(t, err)

		// Verify response structure - mock server nests data under "performance"
		assert.Contains(t, performance, "performance")
		performanceData := performance["performance"].(map[string]interface{})

		assert.Equal(t, jobID, performanceData["job_id"])
		assert.Contains(t, performanceData, "cpu_analytics")
		assert.Contains(t, performanceData, "memory_analytics")
		assert.Contains(t, performanceData, "overall_efficiency")

		t.Logf("Job %s performance response received successfully", jobID)
	})
}

func testAdvancedAnalyticsEndpoints(t *testing.T, baseURL, apiVersion, jobID, stepID string) {
	// Test live metrics endpoint
	t.Run("LiveMetrics", func(t *testing.T) {
		url := fmt.Sprintf("%s/slurm/%s/job/%s/live_metrics", baseURL, apiVersion, jobID)
		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response structure - mock server wraps data under "live_metrics"
		assert.Contains(t, response, "live_metrics")
		liveMetrics := response["live_metrics"].(map[string]interface{})

		assert.Equal(t, jobID, liveMetrics["job_id"])
		assert.Contains(t, liveMetrics, "timestamp")
		assert.Contains(t, liveMetrics, "cpu_usage")
		assert.Contains(t, liveMetrics, "memory_usage")

		t.Logf("Job %s live metrics response received successfully", jobID)
	})

	// Test resource trends endpoint
	t.Run("ResourceTrends", func(t *testing.T) {
		url := fmt.Sprintf("%s/slurm/%s/job/%s/resource_trends?time_window=1h&interval=5m", baseURL, apiVersion, jobID)
		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response structure
		assert.Equal(t, jobID, response["job_id"])
		assert.Contains(t, response, "trends")
		assert.Contains(t, response, "analysis")

		// Verify analysis contains trend information
		analysis := response["analysis"].(map[string]interface{})
		assert.Contains(t, analysis, "cpu_trend")
		assert.Contains(t, analysis, "memory_trend")

		t.Logf("Job %s resource trends response received successfully", jobID)
	})

	// Test step utilization endpoint
	t.Run("StepUtilization", func(t *testing.T) {
		url := fmt.Sprintf("%s/slurm/%s/job/%s/step/%s/utilization", baseURL, apiVersion, jobID, stepID)
		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response structure - mock server wraps data under "step_utilization"
		assert.Contains(t, response, "step_utilization")
		stepUtil := response["step_utilization"].(map[string]interface{})

		assert.Equal(t, jobID, stepUtil["job_id"])
		assert.Equal(t, stepID, stepUtil["step_id"])
		assert.Contains(t, stepUtil, "cpu_utilization")
		assert.Contains(t, stepUtil, "memory_utilization")

		t.Logf("Job %s Step %s utilization response received successfully", jobID, stepID)
	})
}

func testHistoricalAnalyticsEndpoints(t *testing.T, baseURL, apiVersion string) {
	// Test performance history endpoint
	t.Run("PerformanceHistory", func(t *testing.T) {
		url := fmt.Sprintf("%s/slurm/%s/jobs/performance/history?job_id=1001&limit=10", baseURL, apiVersion)
		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var history map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&history)
		require.NoError(t, err)

		// Verify response structure - mock server nests data under "performance_history"
		assert.Contains(t, history, "performance_history")
		historyData := history["performance_history"].(map[string]interface{})
		assert.Contains(t, historyData, "job_id")
		assert.Contains(t, historyData, "start_time")
		assert.Contains(t, historyData, "time_series_data")

		t.Logf("Performance history response received successfully")
	})

	// Test performance trends endpoint
	t.Run("PerformanceTrends", func(t *testing.T) {
		url := fmt.Sprintf("%s/slurm/%s/jobs/performance/trends?time_window=7d&partition=compute", baseURL, apiVersion)
		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response structure - mock server wraps data under "performance_trends"
		assert.Contains(t, response, "performance_trends")
		trends := response["performance_trends"].(map[string]interface{})

		assert.Contains(t, trends, "cluster_performance")
		assert.Contains(t, trends, "resource_trends")
		assert.Contains(t, trends, "partition_trends")

		t.Logf("Performance trends response received successfully")
	})

	// Test user efficiency trends endpoint
	t.Run("UserEfficiencyTrends", func(t *testing.T) {
		url := fmt.Sprintf("%s/slurm/%s/jobs/efficiency/trends?user_id=testuser&time_window=30d", baseURL, apiVersion)
		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response structure - mock server wraps data under "user_efficiency_trends"
		assert.Contains(t, response, "user_efficiency_trends")
		userTrends := response["user_efficiency_trends"].(map[string]interface{})

		assert.Contains(t, userTrends, "user_id")
		assert.Contains(t, userTrends, "efficiency_trends")
		assert.Contains(t, userTrends, "monthly_data")

		t.Logf("User efficiency trends response received successfully")
	})
}

func testComparativeAnalyticsEndpoints(t *testing.T, baseURL, apiVersion, jobID string) {
	// Test job performance comparison endpoint
	t.Run("CompareJobPerformance", func(t *testing.T) {
		url := fmt.Sprintf("%s/slurm/%s/jobs/performance/compare", baseURL, apiVersion)
		payload := `{"job_ids": ["1001", "1002"]}`

		resp, err := http.Post(url, "application/json",
			strings.NewReader(fmt.Sprintf(`{"job_ids": ["%s", "1002"]}`, jobID)))
		if err != nil {
			// Fallback to GET request if POST is not handled
			getUrl := fmt.Sprintf("%s?job_ids=%s,1002", url, jobID)
			resp, err = http.Get(getUrl)
		}
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var comparison map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&comparison)
		require.NoError(t, err)

		// Verify response structure - mock server nests data under "performance_comparison"
		assert.Contains(t, comparison, "performance_comparison")
		comparisonData := comparison["performance_comparison"].(map[string]interface{})
		assert.Contains(t, comparisonData, "metrics")

		t.Logf("Job performance comparison response received successfully, payload: %s", payload)
	})

	// Test similar jobs performance endpoint
	t.Run("SimilarJobsPerformance", func(t *testing.T) {
		url := fmt.Sprintf("%s/slurm/%s/jobs/performance/similar?reference_job_id=%s&criteria=cpus,memory&limit=5", baseURL, apiVersion, jobID)
		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var similar map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&similar)
		require.NoError(t, err)

		// Verify response structure - mock server nests data under "similar_jobs_performance"
		assert.Contains(t, similar, "similar_jobs_performance")
		similarData := similar["similar_jobs_performance"].(map[string]interface{})
		assert.Contains(t, similarData, "reference_job_id")
		assert.Contains(t, similarData, "similar_jobs")

		t.Logf("Similar jobs performance response received successfully")
	})

	// Test batch analysis endpoint
	t.Run("AnalyzeBatchJobs", func(t *testing.T) {
		url := fmt.Sprintf("%s/slurm/%s/jobs/performance/analyze_batch", baseURL, apiVersion)
		payload := fmt.Sprintf(`{"job_ids": ["%s", "1002"]}`, jobID)

		resp, err := http.Post(url, "application/json",
			strings.NewReader(fmt.Sprintf(`{"job_ids": ["%s", "1002"]}`, jobID)))
		if err != nil {
			// Fallback to GET request if POST is not handled
			getUrl := fmt.Sprintf("%s?job_ids=%s,1002", url, jobID)
			resp, err = http.Get(getUrl)
		}
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var batchAnalysis map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&batchAnalysis)
		require.NoError(t, err)

		// Verify response structure - mock server nests data under "batch_analysis"
		assert.Contains(t, batchAnalysis, "batch_analysis")
		batchData := batchAnalysis["batch_analysis"].(map[string]interface{})
		assert.Contains(t, batchData, "job_analyses")
		assert.Contains(t, batchData, "analysis_summary")

		t.Logf("Batch analysis response received successfully, payload: %s", payload)
	})

	// Test workflow performance endpoint
	t.Run("WorkflowPerformance", func(t *testing.T) {
		url := fmt.Sprintf("%s/slurm/%s/jobs/workflow/performance?workflow_id=test-workflow-1&user_id=testuser", baseURL, apiVersion)
		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var workflow map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&workflow)
		require.NoError(t, err)

		// Verify response structure - mock server nests data under "workflow_performance"
		assert.Contains(t, workflow, "workflow_performance")
		workflowData := workflow["workflow_performance"].(map[string]interface{})
		assert.Contains(t, workflowData, "workflow_id")
		assert.Contains(t, workflowData, "job_performance")
		assert.Contains(t, workflowData, "performance_summary")

		t.Logf("Workflow performance response received successfully")
	})

	// Test efficiency report generation endpoint
	t.Run("GenerateEfficiencyReport", func(t *testing.T) {
		url := fmt.Sprintf("%s/slurm/%s/jobs/efficiency/report?user_id=testuser&time_window=7d&format=detailed", baseURL, apiVersion)
		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var report map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&report)
		require.NoError(t, err)

		// Verify response structure - mock server nests data under "efficiency_report"
		assert.Contains(t, report, "efficiency_report")
		reportData := report["efficiency_report"].(map[string]interface{})
		assert.Contains(t, reportData, "report_metadata")
		assert.Contains(t, reportData, "executive_summary")
		assert.Contains(t, reportData, "detailed_metrics")

		t.Logf("Efficiency report response received successfully")
	})
}

// TestJobAnalyticsErrorHandling tests error scenarios for analytics endpoints
func TestJobAnalyticsErrorHandling(t *testing.T) {
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	baseURL := mockServer.URL()
	nonExistentJobID := "99999"

	// Test analytics methods with non-existent job
	t.Run("NonExistentJob", func(t *testing.T) {
		endpoints := []string{
			fmt.Sprintf("/slurm/v0.0.42/job/%s/utilization", nonExistentJobID),
			fmt.Sprintf("/slurm/v0.0.42/job/%s/efficiency", nonExistentJobID),
			fmt.Sprintf("/slurm/v0.0.42/job/%s/performance", nonExistentJobID),
		}

		for _, endpoint := range endpoints {
			url := baseURL + endpoint
			resp, err := http.Get(url)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Should return 404 for non-existent job
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)

			var errorResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&errorResp)
			require.NoError(t, err)
			assert.Contains(t, errorResp, "errors")

			t.Logf("Endpoint %s correctly returned 404 for non-existent job", endpoint)
		}
	})

	// Test unsupported operations for different versions
	t.Run("UnsupportedOperations", func(t *testing.T) {
		// Test v0.0.40 with advanced features that shouldn't be supported
		oldMockServer := mocks.NewMockSlurmServerForVersion("v0.0.40")
		defer oldMockServer.Close()

		oldBaseURL := oldMockServer.URL()

		// Test an endpoint that might not be supported in v0.0.40
		if !oldMockServer.GetConfig().SupportedOperations["jobs.live_metrics"] {
			url := oldBaseURL + "/slurm/v0.0.40/job/1001/live_metrics"
			resp, err := http.Get(url)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusNotImplemented, resp.StatusCode)
			t.Logf("v0.0.40 correctly returned 501 for unsupported live_metrics endpoint")
		}
	})
}

// TestJobAnalyticsPerformance tests performance characteristics of analytics endpoints
func TestJobAnalyticsPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	baseURL := mockServer.URL()
	jobID := "1001"

	// Test analytics response times
	t.Run("ResponseTimes", func(t *testing.T) {
		endpoints := []string{
			fmt.Sprintf("/slurm/v0.0.42/job/%s/utilization", jobID),
			fmt.Sprintf("/slurm/v0.0.42/job/%s/efficiency", jobID),
			fmt.Sprintf("/slurm/v0.0.42/job/%s/performance", jobID),
		}

		for _, endpoint := range endpoints {
			start := time.Now()
			url := baseURL + endpoint
			resp, err := http.Get(url)
			responseTime := time.Since(start)

			require.NoError(t, err)
			resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Assert reasonable response times (should be under 1 second for mock)
			assert.Less(t, responseTime, time.Second)

			t.Logf("Endpoint %s response time: %v", endpoint, responseTime)
		}
	})

	// Test concurrent requests
	t.Run("ConcurrentRequests", func(t *testing.T) {
		numRequests := 10
		done := make(chan bool, numRequests)

		start := time.Now()

		for i := 0; i < numRequests; i++ {
			go func(requestID int) {
				url := fmt.Sprintf("%s/slurm/v0.0.42/job/%s/utilization", baseURL, jobID)
				resp, err := http.Get(url)

				assert.NoError(t, err)
				if resp != nil {
					assert.Equal(t, http.StatusOK, resp.StatusCode)
					resp.Body.Close()
				}

				done <- true
			}(i)
		}

		// Wait for all requests to complete
		for i := 0; i < numRequests; i++ {
			<-done
		}

		totalTime := time.Since(start)
		t.Logf("Completed %d concurrent requests in %v", numRequests, totalTime)

		// Should handle concurrent requests reasonably well
		assert.Less(t, totalTime, 5*time.Second)
	})
}

// TestJobAnalyticsVersionFeatures tests version-specific feature availability
func TestJobAnalyticsVersionFeatures(t *testing.T) {
	versionFeatures := map[string][]string{
		"v0.0.40": {"jobs.utilization", "jobs.efficiency", "jobs.performance"},
		"v0.0.41": {"jobs.utilization", "jobs.efficiency", "jobs.performance", "jobs.live_metrics", "jobs.resource_trends", "jobs.step_utilization"},
		"v0.0.42": {"jobs.utilization", "jobs.efficiency", "jobs.performance", "jobs.live_metrics", "jobs.resource_trends", "jobs.step_utilization", "jobs.performance_history", "jobs.performance_trends", "jobs.efficiency_trends", "jobs.compare_performance", "jobs.similar_performance", "jobs.analyze_batch", "jobs.workflow_performance", "jobs.efficiency_report"},
		"v0.0.43": {"jobs.utilization", "jobs.efficiency", "jobs.performance", "jobs.live_metrics", "jobs.resource_trends", "jobs.step_utilization", "jobs.performance_history", "jobs.performance_trends", "jobs.efficiency_trends", "jobs.compare_performance", "jobs.similar_performance", "jobs.analyze_batch", "jobs.workflow_performance", "jobs.efficiency_report"},
	}

	for version, expectedFeatures := range versionFeatures {
		t.Run(version, func(t *testing.T) {
			mockServer := mocks.NewMockSlurmServerForVersion(version)
			defer mockServer.Close()

			config := mockServer.GetConfig()

			// Verify that all expected features are supported
			for _, feature := range expectedFeatures {
				assert.True(t, config.SupportedOperations[feature],
					"Feature %s should be supported in %s", feature, version)
			}

			t.Logf("Version %s supports %d analytics features", version, len(expectedFeatures))
		})
	}
}
