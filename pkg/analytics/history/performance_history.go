// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package history

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/analytics"
)

// PerformanceHistoryTracker provides methods for tracking and analyzing historical performance data
type PerformanceHistoryTracker struct {
	efficiencyCalc *analytics.EfficiencyCalculator
}

// NewPerformanceHistoryTracker creates a new performance history tracker
func NewPerformanceHistoryTracker() *PerformanceHistoryTracker {
	return &PerformanceHistoryTracker{
		efficiencyCalc: analytics.NewEfficiencyCalculator(),
	}
}

// GetJobPerformanceHistory retrieves historical performance data for a job
func (pht *PerformanceHistoryTracker) GetJobPerformanceHistory(
	ctx context.Context,
	job *interfaces.Job,
	samples []interfaces.JobComprehensiveAnalytics,
	opts *interfaces.PerformanceHistoryOptions,
) (*interfaces.JobPerformanceHistory, error) {
	if job == nil {
		return nil, fmt.Errorf("job cannot be nil")
	}
	if len(samples) == 0 {
		return nil, fmt.Errorf("no performance samples provided")
	}

	// Sort samples by start time
	sort.Slice(samples, func(i, j int) bool {
		return samples[i].StartTime.Before(samples[j].StartTime)
	})

	// Filter samples based on time range if specified
	filteredSamples := pht.filterSamplesByTimeRange(samples, opts)
	if len(filteredSamples) == 0 {
		return nil, fmt.Errorf("no samples found in specified time range")
	}

	// Create time series data
	timeSeriesData := pht.createTimeSeriesData(filteredSamples, opts)

	// Calculate statistics
	statistics := pht.calculatePerformanceStatistics(filteredSamples)

	// Analyze trends if requested
	var trends *interfaces.PerformanceTrendAnalysis
	if opts == nil || opts.IncludeTrends {
		trends = pht.analyzeTrends(timeSeriesData)
	}

	// Detect anomalies
	anomalies := pht.detectAnomalies(timeSeriesData, statistics)

	return &interfaces.JobPerformanceHistory{
		JobID:   job.ID,
		JobName: job.Name,
		StartTime: func() time.Time {
			if job.StartTime != nil {
				return *job.StartTime
			}
			return time.Time{}
		}(),
		EndTime: func() time.Time {
			if job.EndTime != nil {
				return *job.EndTime
			}
			return time.Time{}
		}(),
		TimeSeriesData: timeSeriesData,
		Statistics:     statistics,
		Trends:         trends,
		Anomalies:      anomalies,
	}, nil
}

// filterSamplesByTimeRange filters samples based on time range options
func (pht *PerformanceHistoryTracker) filterSamplesByTimeRange(
	samples []interfaces.JobComprehensiveAnalytics,
	opts *interfaces.PerformanceHistoryOptions,
) []interfaces.JobComprehensiveAnalytics {
	if opts == nil || (opts.StartTime == nil && opts.EndTime == nil) {
		return samples
	}

	var filtered []interfaces.JobComprehensiveAnalytics
	for _, sample := range samples {
		if opts.StartTime != nil && sample.StartTime.Before(*opts.StartTime) {
			continue
		}
		if opts.EndTime != nil && sample.StartTime.After(*opts.EndTime) {
			continue
		}
		filtered = append(filtered, sample)
	}
	return filtered
}

// createTimeSeriesData creates time series data from analytics samples
func (pht *PerformanceHistoryTracker) createTimeSeriesData(
	samples []interfaces.JobComprehensiveAnalytics,
	opts *interfaces.PerformanceHistoryOptions,
) []interfaces.PerformanceSnapshot {
	var snapshots []interfaces.PerformanceSnapshot

	// Determine interval
	interval := pht.determineInterval(samples, opts)

	// Group samples by interval
	groups := pht.groupSamplesByInterval(samples, interval)

	// Create snapshots for each interval
	for _, group := range groups {
		if len(group) == 0 {
			continue
		}

		snapshot := pht.createSnapshot(group)
		snapshots = append(snapshots, snapshot)
	}

	return snapshots
}

// determineInterval determines the appropriate interval for time series data
func (pht *PerformanceHistoryTracker) determineInterval(
	samples []interfaces.JobComprehensiveAnalytics,
	opts *interfaces.PerformanceHistoryOptions,
) time.Duration {
	if opts != nil && opts.Interval != "" {
		switch opts.Interval {
		case "hourly":
			return time.Hour
		case "daily":
			return 24 * time.Hour
		case "weekly":
			return 7 * 24 * time.Hour
		}
	}

	// Auto-determine based on sample duration
	if len(samples) < 2 {
		return time.Hour
	}

	duration := samples[len(samples)-1].StartTime.Sub(samples[0].StartTime)
	switch {
	case duration <= 24*time.Hour:
		return time.Hour
	case duration <= 7*24*time.Hour:
		return 6 * time.Hour
	case duration <= 30*24*time.Hour:
		return 24 * time.Hour
	default:
		return 7 * 24 * time.Hour
	}
}

// groupSamplesByInterval groups samples into time intervals
func (pht *PerformanceHistoryTracker) groupSamplesByInterval(
	samples []interfaces.JobComprehensiveAnalytics,
	interval time.Duration,
) [][]interfaces.JobComprehensiveAnalytics {
	if len(samples) == 0 {
		return nil
	}

	var groups [][]interfaces.JobComprehensiveAnalytics
	var currentGroup []interfaces.JobComprehensiveAnalytics

	baseTime := samples[0].StartTime.Truncate(interval)
	currentTime := baseTime

	for _, sample := range samples {
		sampleTime := sample.StartTime.Truncate(interval)

		if sampleTime.After(currentTime) {
			// Start new group
			if len(currentGroup) > 0 {
				groups = append(groups, currentGroup)
			}
			currentGroup = []interfaces.JobComprehensiveAnalytics{sample}
			currentTime = sampleTime
		} else {
			// Add to current group
			currentGroup = append(currentGroup, sample)
		}
	}

	// Add final group
	if len(currentGroup) > 0 {
		groups = append(groups, currentGroup)
	}

	return groups
}

// createSnapshot creates a performance snapshot from a group of samples
func (pht *PerformanceHistoryTracker) createSnapshot(
	group []interfaces.JobComprehensiveAnalytics,
) interfaces.PerformanceSnapshot {
	// Calculate averages for the interval
	var (
		cpuSum    float64
		memSum    float64
		ioSum     float64
		gpuSum    float64
		netSum    float64
		powerSum  float64
		effSum    float64
		count     float64
		timestamp time.Time
	)

	for i, sample := range group {
		if i == 0 {
			timestamp = sample.StartTime
		}

		if sample.CPUAnalytics != nil {
			cpuSum += sample.CPUAnalytics.UtilizationPercent
		}
		if sample.MemoryAnalytics != nil {
			memSum += sample.MemoryAnalytics.UtilizationPercent
		}
		if sample.IOAnalytics != nil {
			ioSum += sample.IOAnalytics.AverageReadBandwidth + sample.IOAnalytics.AverageWriteBandwidth
		}
		effSum += sample.OverallEfficiency
		// Add network and power if available in analytics
		count++
	}

	return interfaces.PerformanceSnapshot{
		Timestamp:         timestamp,
		CPUUtilization:    cpuSum / count,
		MemoryUtilization: memSum / count,
		IOBandwidth:       ioSum / count,
		GPUUtilization:    gpuSum / count,
		NetworkBandwidth:  netSum / count,
		PowerUsage:        powerSum / count,
		Efficiency:        effSum / count,
	}
}

// calculatePerformanceStatistics calculates aggregate statistics
func (pht *PerformanceHistoryTracker) calculatePerformanceStatistics(
	samples []interfaces.JobComprehensiveAnalytics,
) interfaces.PerformanceStatistics {
	var (
		cpuValues []float64
		memValues []float64
		ioValues  []float64
		effValues []float64
	)

	for _, sample := range samples {
		if sample.CPUAnalytics != nil {
			cpuValues = append(cpuValues, sample.CPUAnalytics.UtilizationPercent)
		}
		if sample.MemoryAnalytics != nil {
			memValues = append(memValues, sample.MemoryAnalytics.UtilizationPercent)
		}
		if sample.IOAnalytics != nil {
			ioValues = append(ioValues, sample.IOAnalytics.AverageReadBandwidth+sample.IOAnalytics.AverageWriteBandwidth)
		}
		effValues = append(effValues, sample.OverallEfficiency)
	}

	return interfaces.PerformanceStatistics{
		AverageCPU:        pht.mean(cpuValues),
		AverageMemory:     pht.mean(memValues),
		AverageIO:         pht.mean(ioValues),
		AverageEfficiency: pht.mean(effValues),
		PeakCPU:           pht.max(cpuValues),
		PeakMemory:        pht.max(memValues),
		PeakIO:            pht.max(ioValues),
		MinCPU:            pht.min(cpuValues),
		MinMemory:         pht.min(memValues),
		MinIO:             pht.min(ioValues),
		StdDevCPU:         pht.stdDev(cpuValues),
		StdDevMemory:      pht.stdDev(memValues),
		StdDevIO:          pht.stdDev(ioValues),
	}
}

// analyzeTrends analyzes performance trends
func (pht *PerformanceHistoryTracker) analyzeTrends(
	snapshots []interfaces.PerformanceSnapshot,
) *interfaces.PerformanceTrendAnalysis {
	if len(snapshots) < 2 {
		return nil
	}

	// Extract time series for each metric
	var (
		cpuValues  []float64
		memValues  []float64
		ioValues   []float64
		effValues  []float64
		timestamps []float64
	)

	baseTime := snapshots[0].Timestamp
	for _, snapshot := range snapshots {
		elapsed := snapshot.Timestamp.Sub(baseTime).Hours()
		timestamps = append(timestamps, elapsed)
		cpuValues = append(cpuValues, snapshot.CPUUtilization)
		memValues = append(memValues, snapshot.MemoryUtilization)
		ioValues = append(ioValues, snapshot.IOBandwidth)
		effValues = append(effValues, snapshot.Efficiency)
	}

	// Calculate trends
	cpuTrend := pht.calculateTrend(timestamps, cpuValues)
	memTrend := pht.calculateTrend(timestamps, memValues)
	ioTrend := pht.calculateTrend(timestamps, ioValues)
	effTrend := pht.calculateTrend(timestamps, effValues)

	// Predict future values (simple linear extrapolation)
	nextHour := timestamps[len(timestamps)-1] + 1.0
	predictedCPU := cpuTrend.Slope*nextHour + pht.mean(cpuValues)
	predictedMem := memTrend.Slope*nextHour + pht.mean(memValues)

	// Estimate runtime based on efficiency trend
	avgEfficiency := pht.mean(effValues)
	var predictedRuntime time.Duration
	if avgEfficiency > 0 {
		// Simple estimation: lower efficiency = longer runtime
		runtimeHours := 100.0 / avgEfficiency * float64(len(snapshots))
		predictedRuntime = time.Duration(runtimeHours) * time.Hour
	}

	return &interfaces.PerformanceTrendAnalysis{
		CPUTrend:         cpuTrend,
		MemoryTrend:      memTrend,
		IOTrend:          ioTrend,
		EfficiencyTrend:  effTrend,
		PredictedCPU:     math.Max(0, math.Min(100, predictedCPU)),
		PredictedMemory:  math.Max(0, math.Min(100, predictedMem)),
		PredictedRuntime: predictedRuntime,
	}
}

// calculateTrend calculates trend information for a metric
func (pht *PerformanceHistoryTracker) calculateTrend(x, y []float64) interfaces.TrendInfo {
	if len(x) != len(y) || len(x) < 2 {
		return interfaces.TrendInfo{
			Direction:  "stable",
			Slope:      0,
			Confidence: 0,
			ChangeRate: 0,
		}
	}

	// Calculate linear regression
	n := float64(len(x))
	sumX, sumY, sumXY, sumX2 := 0.0, 0.0, 0.0, 0.0

	for i := range x {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
	}

	// Calculate slope
	denominator := n*sumX2 - sumX*sumX
	if denominator == 0 {
		return interfaces.TrendInfo{
			Direction:  "stable",
			Slope:      0,
			Confidence: 0,
			ChangeRate: 0,
		}
	}

	slope := (n*sumXY - sumX*sumY) / denominator

	// Calculate R-squared for confidence
	yMean := sumY / n
	ssTotal, ssResidual := 0.0, 0.0
	for i := range y {
		predicted := slope*x[i] + (sumY-slope*sumX)/n
		ssTotal += (y[i] - yMean) * (y[i] - yMean)
		ssResidual += (y[i] - predicted) * (y[i] - predicted)
	}

	confidence := 0.0
	if ssTotal > 0 {
		confidence = 1.0 - (ssResidual / ssTotal)
	}

	// Determine direction and change rate
	direction := "stable"
	changeRate := 0.0
	if math.Abs(slope) > 0.1 {
		if slope > 0 {
			direction = "increasing"
		} else {
			direction = "decreasing"
		}
		// Calculate percentage change per hour
		//nolint:gosec // G602 - y[0] is safe, length checked at function start (len(y) >= 2)
		if y[0] != 0 {
			changeRate = (slope / y[0]) * 100
		}
	}

	return interfaces.TrendInfo{
		Direction:  direction,
		Slope:      slope,
		Confidence: math.Max(0, confidence),
		ChangeRate: changeRate,
	}
}

// detectAnomalies detects anomalies in performance data
func (pht *PerformanceHistoryTracker) detectAnomalies(
	snapshots []interfaces.PerformanceSnapshot,
	stats interfaces.PerformanceStatistics,
) []interfaces.PerformanceAnomaly {
	anomalies := make([]interfaces.PerformanceAnomaly, 0)

	// Define thresholds
	cpuThreshold := stats.StdDevCPU * 2
	memThreshold := stats.StdDevMemory * 2
	ioThreshold := stats.StdDevIO * 2

	for _, snapshot := range snapshots {
		// Check CPU anomalies
		cpuDev := math.Abs(snapshot.CPUUtilization - stats.AverageCPU)
		if cpuDev > cpuThreshold && cpuThreshold > 0 {
			severity := pht.calculateAnomalySeverity(cpuDev, cpuThreshold)
			anomalies = append(anomalies, interfaces.PerformanceAnomaly{
				Timestamp:   snapshot.Timestamp,
				Type:        pht.getAnomalyType(snapshot.CPUUtilization, stats.AverageCPU),
				Metric:      "cpu",
				Severity:    severity,
				Value:       snapshot.CPUUtilization,
				Expected:    stats.AverageCPU,
				Deviation:   (cpuDev / stats.AverageCPU) * 100,
				Description: fmt.Sprintf("CPU utilization %.1f%% (expected %.1f%%)", snapshot.CPUUtilization, stats.AverageCPU),
			})
		}

		// Check Memory anomalies
		memDev := math.Abs(snapshot.MemoryUtilization - stats.AverageMemory)
		if memDev > memThreshold && memThreshold > 0 {
			severity := pht.calculateAnomalySeverity(memDev, memThreshold)
			anomalies = append(anomalies, interfaces.PerformanceAnomaly{
				Timestamp:   snapshot.Timestamp,
				Type:        pht.getAnomalyType(snapshot.MemoryUtilization, stats.AverageMemory),
				Metric:      "memory",
				Severity:    severity,
				Value:       snapshot.MemoryUtilization,
				Expected:    stats.AverageMemory,
				Deviation:   (memDev / stats.AverageMemory) * 100,
				Description: fmt.Sprintf("Memory utilization %.1f%% (expected %.1f%%)", snapshot.MemoryUtilization, stats.AverageMemory),
			})
		}

		// Check I/O anomalies
		ioDev := math.Abs(snapshot.IOBandwidth - stats.AverageIO)
		if ioDev > ioThreshold && ioThreshold > 0 {
			severity := pht.calculateAnomalySeverity(ioDev, ioThreshold)
			anomalies = append(anomalies, interfaces.PerformanceAnomaly{
				Timestamp:   snapshot.Timestamp,
				Type:        pht.getAnomalyType(snapshot.IOBandwidth, stats.AverageIO),
				Metric:      "io",
				Severity:    severity,
				Value:       snapshot.IOBandwidth,
				Expected:    stats.AverageIO,
				Deviation:   (ioDev / stats.AverageIO) * 100,
				Description: fmt.Sprintf("I/O bandwidth %.1f MB/s (expected %.1f MB/s)", snapshot.IOBandwidth, stats.AverageIO),
			})
		}

		// Check efficiency drops
		if snapshot.Efficiency < stats.AverageEfficiency*0.7 {
			anomalies = append(anomalies, interfaces.PerformanceAnomaly{
				Timestamp:   snapshot.Timestamp,
				Type:        "drop",
				Metric:      "efficiency",
				Severity:    "high",
				Value:       snapshot.Efficiency,
				Expected:    stats.AverageEfficiency,
				Deviation:   ((stats.AverageEfficiency - snapshot.Efficiency) / stats.AverageEfficiency) * 100,
				Description: fmt.Sprintf("Efficiency dropped to %.1f%% (average %.1f%%)", snapshot.Efficiency, stats.AverageEfficiency),
			})
		}
	}

	return anomalies
}

// getAnomalyType determines the type of anomaly
func (pht *PerformanceHistoryTracker) getAnomalyType(value, expected float64) string {
	if value > expected {
		return "spike"
	}
	return "drop"
}

// calculateAnomalySeverity calculates anomaly severity
func (pht *PerformanceHistoryTracker) calculateAnomalySeverity(deviation, threshold float64) string {
	ratio := deviation / threshold
	switch {
	case ratio >= 3:
		return "critical"
	case ratio >= 2:
		return "high"
	case ratio >= 1.5:
		return "medium"
	default:
		return "low"
	}
}

// Statistical helper functions

func (pht *PerformanceHistoryTracker) mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func (pht *PerformanceHistoryTracker) max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	maxVal := values[0]
	for _, v := range values {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

func (pht *PerformanceHistoryTracker) min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	minVal := values[0]
	for _, v := range values {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

func (pht *PerformanceHistoryTracker) stdDev(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	mean := pht.mean(values)
	variance := 0.0
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(values))
	return math.Sqrt(variance)
}
