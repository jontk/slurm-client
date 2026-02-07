// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// NOTE: Tests updated for api/ package type changes.
// Job struct fields: ID->JobID (*int32), CPUs (*uint32), Memory->MemoryPerNode (*uint64 in MB)
// StartTime/EndTime are now value types (time.Time), not pointers

package history

import (
	"context"
	"math"
	"testing"
	"time"

	types "github.com/jontk/slurm-client/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper functions for pointer types
func ptrInt32(i int32) *int32    { return &i }
func ptrUint32(u uint32) *uint32 { return &u }
func ptrUint64(u uint64) *uint64 { return &u }
func ptrString(s string) *string { return &s }

func TestPerformanceHistoryTracker_GetJobPerformanceHistory(t *testing.T) {
	tracker := NewPerformanceHistoryTracker()

	// Create test job
	startTime := time.Now().Add(-2 * time.Hour)
	endTime := time.Now()
	job := &types.Job{
		JobID:         ptrInt32(123),
		Name:          ptrString("test-job"),
		CPUs:          ptrUint32(16),
		MemoryPerNode: ptrUint64(64 * 1024), // 64GB in MB
		StartTime:     startTime,
		EndTime:       endTime,
	}

	// Create test samples
	samples := generateTestSamples(startTime, endTime, 12) // 12 samples over 2 hours

	// Test with default options
	history, err := tracker.GetJobPerformanceHistory(context.Background(), job, samples, nil)
	require.NoError(t, err)
	require.NotNil(t, history)

	// Verify basic structure
	assert.Equal(t, "123", history.JobID)
	assert.Equal(t, "test-job", history.JobName)
	assert.Equal(t, startTime, history.StartTime)
	assert.Equal(t, endTime, history.EndTime)

	// Verify time series data
	assert.NotEmpty(t, history.TimeSeriesData)

	// Verify statistics
	assert.Greater(t, history.Statistics.AverageCPU, 0.0)
	assert.Greater(t, history.Statistics.AverageMemory, 0.0)
	assert.Greater(t, history.Statistics.AverageEfficiency, 0.0)

	// Verify trends are calculated
	assert.NotNil(t, history.Trends)

	// Verify anomalies detection (should return empty slice, not nil)
	if history.Anomalies == nil {
		t.Errorf("Anomalies should not be nil, got nil")
	}
}

func TestPerformanceHistoryTracker_GetJobPerformanceHistory_WithOptions(t *testing.T) {
	tracker := NewPerformanceHistoryTracker()

	startTime := time.Now().Add(-4 * time.Hour)
	endTime := time.Now()
	job := &types.Job{
		JobID:         ptrInt32(456),
		Name:          ptrString("filtered-job"),
		CPUs:          ptrUint32(8),
		MemoryPerNode: ptrUint64(32 * 1024), // 32GB in MB
		StartTime:     startTime,
		EndTime:       endTime,
	}

	samples := generateTestSamples(startTime, endTime, 24) // 24 samples over 4 hours

	// Test with time range filtering
	filterStart := startTime.Add(1 * time.Hour)
	filterEnd := endTime.Add(-1 * time.Hour)
	opts := &types.PerformanceHistoryOptions{
		StartTime:     &filterStart,
		EndTime:       &filterEnd,
		Interval:      "hourly",
		IncludeTrends: true,
	}

	history, err := tracker.GetJobPerformanceHistory(context.Background(), job, samples, opts)
	require.NoError(t, err)

	// Verify filtering worked
	assert.NotEmpty(t, history.TimeSeriesData)

	// All snapshots should be within the filtered time range
	for _, snapshot := range history.TimeSeriesData {
		assert.True(t, snapshot.Timestamp.After(filterStart) || snapshot.Timestamp.Equal(filterStart))
		assert.True(t, snapshot.Timestamp.Before(filterEnd) || snapshot.Timestamp.Equal(filterEnd))
	}
}

func TestPerformanceHistoryTracker_FilterSamplesByTimeRange(t *testing.T) {
	tracker := NewPerformanceHistoryTracker()

	baseTime := time.Now()
	samples := []types.JobComprehensiveAnalytics{
		{StartTime: baseTime.Add(-3 * time.Hour)},
		{StartTime: baseTime.Add(-2 * time.Hour)},
		{StartTime: baseTime.Add(-1 * time.Hour)},
		{StartTime: baseTime},
	}

	// Test with start time filter
	startTime := baseTime.Add(-2*time.Hour - 30*time.Minute)
	opts := &types.PerformanceHistoryOptions{
		StartTime: &startTime,
	}

	filtered := tracker.filterSamplesByTimeRange(samples, opts)
	assert.Len(t, filtered, 3) // Should exclude the first sample

	// Test with end time filter
	endTime := baseTime.Add(-30 * time.Minute)
	opts = &types.PerformanceHistoryOptions{
		EndTime: &endTime,
	}

	filtered = tracker.filterSamplesByTimeRange(samples, opts)
	assert.Len(t, filtered, 3) // Should exclude the last sample

	// Test with both filters
	opts = &types.PerformanceHistoryOptions{
		StartTime: &startTime,
		EndTime:   &endTime,
	}

	filtered = tracker.filterSamplesByTimeRange(samples, opts)
	assert.Len(t, filtered, 2) // Should exclude first and last samples
}

func TestPerformanceHistoryTracker_DetermineInterval(t *testing.T) {
	tracker := NewPerformanceHistoryTracker()

	baseTime := time.Now()

	tests := []struct {
		name        string
		duration    time.Duration
		sampleCount int
		expected    time.Duration
	}{
		{
			name:        "short duration",
			duration:    2 * time.Hour,
			sampleCount: 10,
			expected:    time.Hour,
		},
		{
			name:        "medium duration",
			duration:    3 * 24 * time.Hour,
			sampleCount: 10,
			expected:    6 * time.Hour,
		},
		{
			name:        "long duration",
			duration:    15 * 24 * time.Hour,
			sampleCount: 10,
			expected:    24 * time.Hour,
		},
		{
			name:        "very long duration",
			duration:    60 * 24 * time.Hour,
			sampleCount: 10,
			expected:    7 * 24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			samples := make([]types.JobComprehensiveAnalytics, tt.sampleCount)
			for i := range samples {
				samples[i].StartTime = baseTime.Add(time.Duration(i) * tt.duration / time.Duration(tt.sampleCount))
			}

			interval := tracker.determineInterval(samples, nil)
			assert.Equal(t, tt.expected, interval)
		})
	}

	// Test with explicit interval option
	samples := []types.JobComprehensiveAnalytics{
		{StartTime: baseTime},
		{StartTime: baseTime.Add(time.Hour)},
	}
	opts := &types.PerformanceHistoryOptions{
		Interval: "daily",
	}

	interval := tracker.determineInterval(samples, opts)
	assert.Equal(t, 24*time.Hour, interval)
}

func TestPerformanceHistoryTracker_CalculateTrend(t *testing.T) {
	tracker := NewPerformanceHistoryTracker()

	// Test increasing trend
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{10, 15, 20, 25, 30} // Clear upward trend

	trend := tracker.calculateTrend(x, y)
	assert.Equal(t, "increasing", trend.Direction)
	assert.Greater(t, trend.Slope, 0.0)
	assert.Greater(t, trend.Confidence, 0.9) // Should be very confident
	assert.Greater(t, trend.ChangeRate, 0.0)

	// Test decreasing trend
	y = []float64{30, 25, 20, 15, 10} // Clear downward trend

	trend = tracker.calculateTrend(x, y)
	assert.Equal(t, "decreasing", trend.Direction)
	assert.Less(t, trend.Slope, 0.0)
	assert.Greater(t, trend.Confidence, 0.9)
	assert.Less(t, trend.ChangeRate, 0.0)

	// Test stable trend
	y = []float64{20, 20.1, 19.9, 20.05, 19.95} // Very stable

	trend = tracker.calculateTrend(x, y)
	assert.Equal(t, "stable", trend.Direction)
	assert.InDelta(t, 0.0, trend.Slope, 0.1)

	// Test empty input
	trend = tracker.calculateTrend([]float64{}, []float64{})
	assert.Equal(t, "stable", trend.Direction)
	assert.Equal(t, 0.0, trend.Slope)
	assert.Equal(t, 0.0, trend.Confidence)
}

func TestPerformanceHistoryTracker_DetectAnomalies(t *testing.T) {
	tracker := NewPerformanceHistoryTracker()

	// Create snapshots with some anomalies
	baseTime := time.Now()
	snapshots := []types.PerformanceSnapshot{
		{Timestamp: baseTime.Add(-4 * time.Hour), CPUUtilization: 75.0, MemoryUtilization: 60.0, Efficiency: 70.0},
		{Timestamp: baseTime.Add(-3 * time.Hour), CPUUtilization: 78.0, MemoryUtilization: 62.0, Efficiency: 72.0},
		{Timestamp: baseTime.Add(-2 * time.Hour), CPUUtilization: 95.0, MemoryUtilization: 58.0, Efficiency: 68.0}, // CPU spike
		{Timestamp: baseTime.Add(-1 * time.Hour), CPUUtilization: 76.0, MemoryUtilization: 90.0, Efficiency: 65.0}, // Memory spike
		{Timestamp: baseTime, CPUUtilization: 74.0, MemoryUtilization: 61.0, Efficiency: 45.0},                     // Efficiency drop
	}

	stats := types.PerformanceStatistics{
		AverageCPU:        75.6,
		AverageMemory:     66.2,
		AverageIO:         150.0,
		AverageEfficiency: 65.0, // Average of efficiency values: (70+72+68+65+45)/5 = 64
		StdDevCPU:         8.5,
		StdDevMemory:      10.0,
		StdDevIO:          10.0,
	}

	anomalies := tracker.detectAnomalies(snapshots, stats)

	// Should detect anomalies
	assert.NotEmpty(t, anomalies)

	// Check for specific anomaly types
	foundCPUSpike := false
	foundMemorySpike := false
	foundEfficiencyDrop := false

	for _, anomaly := range anomalies {
		switch anomaly.Metric {
		case "cpu":
			if anomaly.Type == "spike" {
				foundCPUSpike = true
				assert.Equal(t, 95.0, anomaly.Value)
				assert.InDelta(t, 75.6, anomaly.Expected, 1.0)
			}
		case "memory":
			if anomaly.Type == "spike" {
				foundMemorySpike = true
				assert.Equal(t, 90.0, anomaly.Value)
			}
		case "efficiency":
			if anomaly.Type == "drop" {
				foundEfficiencyDrop = true
				assert.Equal(t, 45.0, anomaly.Value)
			}
		}
	}

	assert.True(t, foundCPUSpike, "Should detect CPU spike anomaly")
	assert.True(t, foundMemorySpike, "Should detect memory spike anomaly")
	assert.True(t, foundEfficiencyDrop, "Should detect efficiency drop anomaly")
}

func TestPerformanceHistoryTracker_CalculatePerformanceStatistics(t *testing.T) {
	tracker := NewPerformanceHistoryTracker()

	samples := []types.JobComprehensiveAnalytics{
		{
			CPUAnalytics:      &types.CPUAnalytics{UtilizationPercent: 70.0},
			MemoryAnalytics:   &types.MemoryAnalytics{UtilizationPercent: 65.0},
			IOAnalytics:       &types.IOAnalytics{AverageReadBandwidth: 100.0, AverageWriteBandwidth: 50.0},
			OverallEfficiency: 68.0,
		},
		{
			CPUAnalytics:      &types.CPUAnalytics{UtilizationPercent: 80.0},
			MemoryAnalytics:   &types.MemoryAnalytics{UtilizationPercent: 75.0},
			IOAnalytics:       &types.IOAnalytics{AverageReadBandwidth: 120.0, AverageWriteBandwidth: 60.0},
			OverallEfficiency: 78.0,
		},
		{
			CPUAnalytics:      &types.CPUAnalytics{UtilizationPercent: 75.0},
			MemoryAnalytics:   &types.MemoryAnalytics{UtilizationPercent: 70.0},
			IOAnalytics:       &types.IOAnalytics{AverageReadBandwidth: 110.0, AverageWriteBandwidth: 55.0},
			OverallEfficiency: 73.0,
		},
	}

	stats := tracker.calculatePerformanceStatistics(samples)

	// Check averages
	assert.InDelta(t, 75.0, stats.AverageCPU, 0.1)        // (70+80+75)/3
	assert.InDelta(t, 70.0, stats.AverageMemory, 0.1)     // (65+75+70)/3
	assert.InDelta(t, 165.0, stats.AverageIO, 0.1)        // ((100+50)+(120+60)+(110+55))/3
	assert.InDelta(t, 73.0, stats.AverageEfficiency, 0.1) // (68+78+73)/3

	// Check peaks
	assert.Equal(t, 80.0, stats.PeakCPU)
	assert.Equal(t, 75.0, stats.PeakMemory)
	assert.Equal(t, 180.0, stats.PeakIO) // 120+60

	// Check minimums
	assert.Equal(t, 70.0, stats.MinCPU)
	assert.Equal(t, 65.0, stats.MinMemory)
	assert.Equal(t, 150.0, stats.MinIO) // 100+50

	// Check standard deviations
	assert.Greater(t, stats.StdDevCPU, 0.0)
	assert.Greater(t, stats.StdDevMemory, 0.0)
	assert.Greater(t, stats.StdDevIO, 0.0)
}

func TestPerformanceHistoryTracker_StatisticalHelpers(t *testing.T) {
	tracker := NewPerformanceHistoryTracker()

	values := []float64{10.0, 20.0, 30.0, 40.0, 50.0}

	// Test mean
	mean := tracker.mean(values)
	assert.InDelta(t, 30.0, mean, 0.01)

	// Test max
	maxVal := tracker.max(values)
	assert.Equal(t, 50.0, maxVal)

	// Test min
	minVal := tracker.min(values)
	assert.Equal(t, 10.0, minVal)

	// Test standard deviation
	stdDev := tracker.stdDev(values)
	expectedStdDev := math.Sqrt(200.0) // Variance is 200
	assert.InDelta(t, expectedStdDev, stdDev, 0.01)

	// Test with empty slice
	assert.Equal(t, 0.0, tracker.mean([]float64{}))
	assert.Equal(t, 0.0, tracker.max([]float64{}))
	assert.Equal(t, 0.0, tracker.min([]float64{}))
	assert.Equal(t, 0.0, tracker.stdDev([]float64{}))
}

func TestPerformanceHistoryTracker_ErrorCases(t *testing.T) {
	tracker := NewPerformanceHistoryTracker()

	// Test with nil job
	_, err := tracker.GetJobPerformanceHistory(context.Background(), nil, []types.JobComprehensiveAnalytics{}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "job cannot be nil")

	// Test with empty samples
	job := &types.Job{JobID: ptrInt32(999)}
	_, err = tracker.GetJobPerformanceHistory(context.Background(), job, []types.JobComprehensiveAnalytics{}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no performance samples provided")

	// Test with samples filtered out
	startTime := time.Now().Add(-2 * time.Hour)
	endTime := time.Now()
	job = &types.Job{
		JobID:     ptrInt32(999),
		StartTime: startTime,
		EndTime:   endTime,
	}

	samples := []types.JobComprehensiveAnalytics{
		{StartTime: time.Now().Add(-5 * time.Hour)}, // Outside range
	}

	filterStart := time.Now().Add(-1 * time.Hour)
	filterEnd := time.Now()
	opts := &types.PerformanceHistoryOptions{
		StartTime: &filterStart,
		EndTime:   &filterEnd,
	}

	_, err = tracker.GetJobPerformanceHistory(context.Background(), job, samples, opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no samples found in specified time range")
}

func TestPerformanceHistoryTracker_GroupSamplesByInterval(t *testing.T) {
	tracker := NewPerformanceHistoryTracker()

	baseTime := time.Now().Truncate(time.Hour)
	samples := []types.JobComprehensiveAnalytics{
		{StartTime: baseTime},
		{StartTime: baseTime.Add(30 * time.Minute)},             // Same hour
		{StartTime: baseTime.Add(1*time.Hour + 15*time.Minute)}, // Next hour
		{StartTime: baseTime.Add(1*time.Hour + 45*time.Minute)}, // Same hour as above
		{StartTime: baseTime.Add(2 * time.Hour)},                // Third hour
	}

	groups := tracker.groupSamplesByInterval(samples, time.Hour)

	// Should have 3 groups (3 different hours)
	assert.Len(t, groups, 3)

	// First group should have 2 samples
	assert.Len(t, groups[0], 2)

	// Second group should have 2 samples
	assert.Len(t, groups[1], 2)

	// Third group should have 1 sample
	assert.Len(t, groups[2], 1)
}

func TestPerformanceHistoryTracker_CreateSnapshot(t *testing.T) {
	tracker := NewPerformanceHistoryTracker()

	baseTime := time.Now()
	group := []types.JobComprehensiveAnalytics{
		{
			StartTime:         baseTime,
			CPUAnalytics:      &types.CPUAnalytics{UtilizationPercent: 70.0},
			MemoryAnalytics:   &types.MemoryAnalytics{UtilizationPercent: 60.0},
			IOAnalytics:       &types.IOAnalytics{AverageReadBandwidth: 100.0, AverageWriteBandwidth: 50.0},
			OverallEfficiency: 65.0,
		},
		{
			StartTime:         baseTime.Add(15 * time.Minute),
			CPUAnalytics:      &types.CPUAnalytics{UtilizationPercent: 80.0},
			MemoryAnalytics:   &types.MemoryAnalytics{UtilizationPercent: 70.0},
			IOAnalytics:       &types.IOAnalytics{AverageReadBandwidth: 120.0, AverageWriteBandwidth: 60.0},
			OverallEfficiency: 75.0,
		},
	}

	snapshot := tracker.createSnapshot(group)

	// Should use the first timestamp
	assert.Equal(t, baseTime, snapshot.Timestamp)

	// Should average the metrics
	assert.InDelta(t, 75.0, snapshot.CPUUtilization, 0.1)    // (70+80)/2
	assert.InDelta(t, 65.0, snapshot.MemoryUtilization, 0.1) // (60+70)/2
	assert.InDelta(t, 165.0, snapshot.IOBandwidth, 0.1)      // ((100+50)+(120+60))/2
	assert.InDelta(t, 70.0, snapshot.Efficiency, 0.1)        // (65+75)/2
}

// Helper function to generate test samples
func generateTestSamples(startTime, endTime time.Time, count int) []types.JobComprehensiveAnalytics {
	samples := make([]types.JobComprehensiveAnalytics, count)
	duration := endTime.Sub(startTime)
	interval := duration / time.Duration(count)

	for i := range count {
		timestamp := startTime.Add(time.Duration(i) * interval)
		progress := float64(i) / float64(count)

		// Simulate varying performance over time
		cpuUtil := 70.0 + 10.0*math.Sin(progress*math.Pi*2) // Oscillating between 60-80%
		memUtil := 60.0 + 15.0*progress                     // Gradually increasing 60-75%
		ioUtil := 50.0 + 20.0*(1.0-progress)                // Gradually decreasing 70-50%
		efficiency := (cpuUtil + memUtil + ioUtil) / 3.0

		samples[i] = types.JobComprehensiveAnalytics{
			JobID:     1,
			StartTime: timestamp,
			CPUAnalytics: &types.CPUAnalytics{
				UtilizationPercent: cpuUtil,
			},
			MemoryAnalytics: &types.MemoryAnalytics{
				UtilizationPercent: memUtil,
			},
			IOAnalytics: &types.IOAnalytics{
				AverageReadBandwidth:  ioUtil * 2.0,
				AverageWriteBandwidth: ioUtil * 1.0,
			},
			OverallEfficiency: efficiency,
		}
	}

	return samples
}
