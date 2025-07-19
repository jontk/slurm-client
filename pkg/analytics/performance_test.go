package analytics

import (
	"testing"
	"time"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompareJobPerformance(t *testing.T) {
	analyzer := NewPerformanceAnalyzer()
	
	// Create test jobs
	startTimeA := time.Now().Add(-2 * time.Hour)
	endTimeA := time.Now().Add(-1 * time.Hour)
	jobA := &interfaces.Job{
		ID:        "job-a",
		Name:      "test-job",
		CPUs:      16,
		Memory:    64 * 1024 * 1024 * 1024,
		GPUs:      2,
		StartTime: &startTimeA,
		EndTime:   &endTimeA,
	}
	
	startTimeB := time.Now().Add(-90 * time.Minute)
	endTimeB := time.Now().Add(-45 * time.Minute)
	jobB := &interfaces.Job{
		ID:        "job-b",
		Name:      "test-job-optimized",
		CPUs:      12,
		Memory:    48 * 1024 * 1024 * 1024,
		GPUs:      2,
		StartTime: &startTimeB,
		EndTime:   &endTimeB,
	}
	
	// Create analytics
	analyticsA := &interfaces.JobComprehensiveAnalytics{
		CPUAnalytics: &interfaces.CPUAnalytics{
			AllocatedCores:     16,
			UsedCores:         10.0,
			UtilizationPercent: 62.5,
		},
		MemoryAnalytics: &interfaces.MemoryAnalytics{
			AllocatedBytes:     64 * 1024 * 1024 * 1024,
			UsedBytes:         40 * 1024 * 1024 * 1024,
			UtilizationPercent: 62.5,
		},
		EfficiencyMetrics: &interfaces.JobEfficiencyMetrics{
			OverallEfficiencyScore: 60.0,
			CPUEfficiency:         62.5,
			MemoryEfficiency:      62.5,
			GPUEfficiency:         70.0,
			IOEfficiency:          50.0,
		},
	}
	
	analyticsB := &interfaces.JobComprehensiveAnalytics{
		CPUAnalytics: &interfaces.CPUAnalytics{
			AllocatedCores:     12,
			UsedCores:         10.2,
			UtilizationPercent: 85.0,
		},
		MemoryAnalytics: &interfaces.MemoryAnalytics{
			AllocatedBytes:     48 * 1024 * 1024 * 1024,
			UsedBytes:         40 * 1024 * 1024 * 1024,
			UtilizationPercent: 83.3,
		},
		EfficiencyMetrics: &interfaces.JobEfficiencyMetrics{
			OverallEfficiencyScore: 82.0,
			CPUEfficiency:         85.0,
			MemoryEfficiency:      83.3,
			GPUEfficiency:         85.0,
			IOEfficiency:          75.0,
		},
	}
	
	// Compare jobs
	comparison, err := analyzer.CompareJobPerformance(jobA, analyticsA, jobB, analyticsB)
	require.NoError(t, err)
	require.NotNil(t, comparison)
	
	// Verify comparison metrics
	assert.Equal(t, jobA, comparison.JobA)
	assert.Equal(t, jobB, comparison.JobB)
	
	// Check efficiency deltas
	assert.InDelta(t, 22.0, comparison.Comparison.OverallEfficiencyDelta, 0.1)
	assert.InDelta(t, 22.5, comparison.Comparison.CPUEfficiencyDelta, 0.1)
	assert.InDelta(t, 20.8, comparison.Comparison.MemoryEfficiencyDelta, 0.1)
	
	// Check runtime ratio (B is 45min, A is 60min)
	assert.InDelta(t, 0.75, comparison.Comparison.RuntimeRatio, 0.01)
	
	// Check resource differences
	assert.Equal(t, -4, comparison.Differences.CPUDelta)
	assert.InDelta(t, -16.0, comparison.Differences.MemoryDeltaGB, 0.1)
	assert.Equal(t, 0, comparison.Differences.GPUDelta)
	
	// Job B should be the winner (better efficiency, faster, lower resources)
	assert.Equal(t, "B", comparison.Winner)
	
	// Check summary contains key information
	assert.Contains(t, comparison.Summary, "Job B (job-b) performed better overall")
	assert.Contains(t, comparison.Summary, "B completed significantly faster")
	assert.Contains(t, comparison.Summary, "significantly better resource efficiency")
}

func TestCompareJobPerformance_Errors(t *testing.T) {
	analyzer := NewPerformanceAnalyzer()
	
	job := &interfaces.Job{ID: "test"}
	analytics := &interfaces.JobComprehensiveAnalytics{}
	
	// Test nil jobs
	_, err := analyzer.CompareJobPerformance(nil, analytics, job, analytics)
	assert.Error(t, err)
	
	_, err = analyzer.CompareJobPerformance(job, analytics, nil, analytics)
	assert.Error(t, err)
	
	// Test nil analytics
	_, err = analyzer.CompareJobPerformance(job, nil, job, analytics)
	assert.Error(t, err)
	
	_, err = analyzer.CompareJobPerformance(job, analytics, job, nil)
	assert.Error(t, err)
}

func TestGetSimilarJobsPerformance(t *testing.T) {
	analyzer := NewPerformanceAnalyzer()
	
	// Create reference job
	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now()
	referenceJob := &interfaces.Job{
		ID:        "ref-job",
		Name:      "matrix-multiply",
		UserID:    "user123",
		Partition: "compute",
		CPUs:      16,
		Memory:    64 * 1024 * 1024 * 1024,
		GPUs:      1,
		StartTime: &startTime,
		EndTime:   &endTime,
	}
	
	referenceAnalytics := &interfaces.JobComprehensiveAnalytics{
		EfficiencyMetrics: &interfaces.JobEfficiencyMetrics{
			OverallEfficiencyScore: 70.0,
		},
	}
	
	// Create candidate jobs
	candidates := []struct {
		Job       *interfaces.Job
		Analytics *interfaces.JobComprehensiveAnalytics
	}{
		{
			// Very similar job with better performance
			Job: &interfaces.Job{
				ID:        "similar-1",
				Name:      "matrix-multiply",
				UserID:    "user123",
				Partition: "compute",
				CPUs:      16,
				Memory:    64 * 1024 * 1024 * 1024,
				GPUs:      1,
				StartTime: &startTime,
				EndTime:   &[]time.Time{startTime.Add(45 * time.Minute)}[0],
			},
			Analytics: &interfaces.JobComprehensiveAnalytics{
				EfficiencyMetrics: &interfaces.JobEfficiencyMetrics{
					OverallEfficiencyScore: 85.0,
				},
			},
		},
		{
			// Similar job with different resources
			Job: &interfaces.Job{
				ID:        "similar-2",
				Name:      "matrix-multiply",
				UserID:    "user123",
				Partition: "compute",
				CPUs:      12,
				Memory:    48 * 1024 * 1024 * 1024,
				GPUs:      1,
				StartTime: &startTime,
				EndTime:   &[]time.Time{startTime.Add(55 * time.Minute)}[0],
			},
			Analytics: &interfaces.JobComprehensiveAnalytics{
				EfficiencyMetrics: &interfaces.JobEfficiencyMetrics{
					OverallEfficiencyScore: 80.0,
				},
			},
		},
		{
			// Less similar job (different user)
			Job: &interfaces.Job{
				ID:        "less-similar",
				Name:      "matrix-multiply",
				UserID:    "user456",
				Partition: "compute",
				CPUs:      16,
				Memory:    64 * 1024 * 1024 * 1024,
				GPUs:      1,
				StartTime: &startTime,
				EndTime:   &endTime,
			},
			Analytics: &interfaces.JobComprehensiveAnalytics{
				EfficiencyMetrics: &interfaces.JobEfficiencyMetrics{
					OverallEfficiencyScore: 65.0,
				},
			},
		},
		{
			// Dissimilar job (should be filtered out)
			Job: &interfaces.Job{
				ID:        "dissimilar",
				Name:      "deep-learning",
				UserID:    "user789",
				Partition: "gpu",
				CPUs:      32,
				Memory:    128 * 1024 * 1024 * 1024,
				GPUs:      4,
				StartTime: &startTime,
				EndTime:   &endTime,
			},
			Analytics: &interfaces.JobComprehensiveAnalytics{
				EfficiencyMetrics: &interfaces.JobEfficiencyMetrics{
					OverallEfficiencyScore: 90.0,
				},
			},
		},
	}
	
	// Analyze similar jobs
	analysis, err := analyzer.GetSimilarJobsPerformance(
		referenceJob,
		referenceAnalytics,
		candidates,
		0.7, // 70% similarity threshold
	)
	require.NoError(t, err)
	require.NotNil(t, analysis)
	
	// Should find 3 similar jobs (excluding the dissimilar one)
	assert.Len(t, analysis.SimilarJobs, 3)
	
	// Check that jobs are sorted by efficiency (best first)
	assert.Equal(t, "similar-1", analysis.SimilarJobs[0].Job.ID)
	assert.Equal(t, 1, analysis.SimilarJobs[0].PerformanceRank)
	assert.InDelta(t, 15.0, analysis.SimilarJobs[0].EfficiencyDelta, 0.1)
	
	// Check performance statistics
	stats := analysis.PerformanceStats
	assert.InDelta(t, 76.67, stats.AverageEfficiency, 0.1) // (85+80+65)/3
	assert.InDelta(t, 80.0, stats.MedianEfficiency, 0.1)
	assert.Equal(t, 85.0, stats.BestEfficiency)
	assert.Equal(t, 65.0, stats.WorstEfficiency)
	
	// Check optimal resources (from best performing job)
	assert.Equal(t, 16, stats.OptimalResources.CPUs)
	assert.Equal(t, 64.0, stats.OptimalResources.MemoryGB)
	assert.Equal(t, 1, stats.OptimalResources.GPUs)
	
	// Check best practices
	assert.NotEmpty(t, analysis.BestPractices)
	
	// Check recommendations
	assert.NotEmpty(t, analysis.Recommendations)
}

func TestCalculateJobSimilarity(t *testing.T) {
	analyzer := NewPerformanceAnalyzer()
	
	baseJob := &interfaces.Job{
		Name:      "test-job",
		UserID:    "user123",
		Partition: "compute",
		CPUs:      16,
		Memory:    64 * 1024 * 1024 * 1024,
		GPUs:      2,
	}
	
	tests := []struct {
		name     string
		otherJob *interfaces.Job
		minScore float64
		maxScore float64
	}{
		{
			name: "identical job",
			otherJob: &interfaces.Job{
				Name:      "test-job",
				UserID:    "user123",
				Partition: "compute",
				CPUs:      16,
				Memory:    64 * 1024 * 1024 * 1024,
				GPUs:      2,
			},
			minScore: 0.95, // Should be very high
			maxScore: 1.0,
		},
		{
			name: "same user different resources",
			otherJob: &interfaces.Job{
				Name:      "other-job",
				UserID:    "user123",
				Partition: "compute",
				CPUs:      8,
				Memory:    32 * 1024 * 1024 * 1024,
				GPUs:      1,
			},
			minScore: 0.5,
			maxScore: 0.7,
		},
		{
			name: "different user same resources",
			otherJob: &interfaces.Job{
				Name:      "test-job",
				UserID:    "user456",
				Partition: "compute",
				CPUs:      16,
				Memory:    64 * 1024 * 1024 * 1024,
				GPUs:      2,
			},
			minScore: 0.7,
			maxScore: 0.9,
		},
		{
			name: "completely different",
			otherJob: &interfaces.Job{
				Name:      "other-job",
				UserID:    "user456",
				Partition: "gpu",
				CPUs:      32,
				Memory:    128 * 1024 * 1024 * 1024,
				GPUs:      8,
			},
			minScore: 0.2,
			maxScore: 0.5,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			similarity := analyzer.calculateJobSimilarity(baseJob, tt.otherJob)
			assert.GreaterOrEqual(t, similarity, tt.minScore)
			assert.LessOrEqual(t, similarity, tt.maxScore)
		})
	}
}

func TestCalculatePerformanceStatistics(t *testing.T) {
	analyzer := NewPerformanceAnalyzer()
	
	similarJobs := []*SimilarJobResult{
		{
			Job: &interfaces.Job{
				CPUs:   16,
				Memory: 64 * 1024 * 1024 * 1024,
			},
			Analytics: &interfaces.JobComprehensiveAnalytics{
				EfficiencyMetrics: &interfaces.JobEfficiencyMetrics{
					OverallEfficiencyScore: 90.0,
				},
			},
		},
		{
			Job: &interfaces.Job{
				CPUs:   12,
				Memory: 48 * 1024 * 1024 * 1024,
			},
			Analytics: &interfaces.JobComprehensiveAnalytics{
				EfficiencyMetrics: &interfaces.JobEfficiencyMetrics{
					OverallEfficiencyScore: 80.0,
				},
			},
		},
		{
			Job: &interfaces.Job{
				CPUs:   16,
				Memory: 64 * 1024 * 1024 * 1024,
			},
			Analytics: &interfaces.JobComprehensiveAnalytics{
				EfficiencyMetrics: &interfaces.JobEfficiencyMetrics{
					OverallEfficiencyScore: 70.0,
				},
			},
		},
	}
	
	stats := analyzer.calculatePerformanceStatistics(similarJobs)
	
	assert.InDelta(t, 80.0, stats.AverageEfficiency, 0.1)
	assert.InDelta(t, 80.0, stats.MedianEfficiency, 0.1)
	assert.InDelta(t, 8.165, stats.StdDevEfficiency, 0.1)
	assert.Equal(t, 90.0, stats.BestEfficiency)
	assert.Equal(t, 70.0, stats.WorstEfficiency)
	
	// Optimal resources should come from the best performing job
	assert.Equal(t, 16, stats.OptimalResources.CPUs)
	assert.Equal(t, 64.0, stats.OptimalResources.MemoryGB)
}

func TestGenerateRecommendations(t *testing.T) {
	analyzer := NewPerformanceAnalyzer()
	
	// Reference job with suboptimal resources
	referenceJob := &interfaces.Job{
		ID:     "ref-job",
		CPUs:   32, // Over-allocated
		Memory: 128 * 1024 * 1024 * 1024, // Over-allocated
		GPUs:   0,
	}
	
	referenceAnalytics := &interfaces.JobComprehensiveAnalytics{
		EfficiencyMetrics: &interfaces.JobEfficiencyMetrics{
			OverallEfficiencyScore: 55.0, // Below average
			CPUEfficiency:         50.0,
			MemoryEfficiency:      45.0,
		},
	}
	
	analysis := &SimilarJobsAnalysis{
		ReferenceJob: referenceJob,
		PerformanceStats: PerformanceStatistics{
			MedianEfficiency: 75.0,
			MedianRuntime:   45 * time.Minute,
			OptimalResources: ResourceRecommendation{
				CPUs:     16,
				MemoryGB: 64.0,
			},
		},
	}
	
	recommendations := analyzer.generateRecommendations(referenceJob, referenceAnalytics, analysis)
	
	// Should recommend reducing resources
	foundCPUReduction := false
	foundMemoryReduction := false
	foundEfficiencyImprovement := false
	
	for _, rec := range recommendations {
		if contains(rec, "reducing CPU allocation") {
			foundCPUReduction = true
			assert.Contains(t, rec, "16 cores")
		}
		if contains(rec, "reducing memory allocation") {
			foundMemoryReduction = true
			assert.Contains(t, rec, "64.0 GB")
		}
		if contains(rec, "below median") {
			foundEfficiencyImprovement = true
		}
	}
	
	assert.True(t, foundCPUReduction)
	assert.True(t, foundMemoryReduction)
	assert.True(t, foundEfficiencyImprovement)
}

func TestStatisticalFunctions(t *testing.T) {
	analyzer := NewPerformanceAnalyzer()
	
	// Test mean
	values := []float64{10, 20, 30, 40, 50}
	assert.InDelta(t, 30.0, analyzer.mean(values), 0.01)
	
	// Test median (odd count)
	assert.InDelta(t, 30.0, analyzer.median(values), 0.01)
	
	// Test median (even count)
	values = []float64{10, 20, 30, 40}
	assert.InDelta(t, 25.0, analyzer.median(values), 0.01)
	
	// Test standard deviation
	values = []float64{2, 4, 4, 4, 5, 5, 7, 9}
	assert.InDelta(t, 2.0, analyzer.stdDev(values), 0.01)
	
	// Test empty values
	empty := []float64{}
	assert.Equal(t, 0.0, analyzer.mean(empty))
	assert.Equal(t, 0.0, analyzer.median(empty))
	assert.Equal(t, 0.0, analyzer.stdDev(empty))
}

// Helper function
func contains(s string, substr string) bool {
	return len(s) >= len(substr) && (s[:len(substr)] == substr || 
		(len(s) > len(substr) && containsHelper(s[1:], substr)))
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