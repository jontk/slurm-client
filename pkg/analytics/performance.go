// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package analytics

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/jontk/slurm-client/internal/interfaces"
)

// PerformanceAnalyzer provides methods for analyzing and comparing job performance
type PerformanceAnalyzer struct {
	efficiencyCalc *EfficiencyCalculator
}

// NewPerformanceAnalyzer creates a new performance analyzer
func NewPerformanceAnalyzer() *PerformanceAnalyzer {
	return &PerformanceAnalyzer{
		efficiencyCalc: NewEfficiencyCalculator(),
	}
}

// JobPerformanceComparison represents the comparison between two jobs
type JobPerformanceComparison struct {
	JobA        *interfaces.Job
	JobB        *interfaces.Job
	AnalyticsA  *interfaces.JobComprehensiveAnalytics
	AnalyticsB  *interfaces.JobComprehensiveAnalytics
	Comparison  PerformanceMetrics
	Differences ResourceDifferences
	Winner      string // JobA ID, JobB ID, or "tie"
	Summary     string
}

// PerformanceMetrics contains comparative performance metrics
type PerformanceMetrics struct {
	OverallEfficiencyDelta float64                // B - A (positive means B is better)
	CPUEfficiencyDelta     float64
	MemoryEfficiencyDelta  float64
	GPUEfficiencyDelta     float64
	IOEfficiencyDelta      float64
	RuntimeRatio           float64                // B/A (< 1 means B is faster)
	ResourceWasteRatio     map[string]float64    // B/A for each resource
	CostEfficiencyRatio    float64                // Performance per resource unit
}

// ResourceDifferences contains the differences in resource allocation
type ResourceDifferences struct {
	CPUDelta       int     // B - A
	MemoryDeltaGB  float64
	GPUDelta       int
	RuntimeDelta   time.Duration
}

// CompareJobPerformance compares the performance of two jobs
func (pa *PerformanceAnalyzer) CompareJobPerformance(
	jobA *interfaces.Job,
	analyticsA *interfaces.JobComprehensiveAnalytics,
	jobB *interfaces.Job,
	analyticsB *interfaces.JobComprehensiveAnalytics,
) (*JobPerformanceComparison, error) {
	if jobA == nil || jobB == nil {
		return nil, fmt.Errorf("both jobs must be provided")
	}
	if analyticsA == nil || analyticsB == nil {
		return nil, fmt.Errorf("analytics for both jobs must be provided")
	}
	
	comparison := &JobPerformanceComparison{
		JobA:       jobA,
		JobB:       jobB,
		AnalyticsA: analyticsA,
		AnalyticsB: analyticsB,
	}
	
	// Calculate performance metrics
	comparison.Comparison = pa.calculatePerformanceMetrics(jobA, analyticsA, jobB, analyticsB)
	
	// Calculate resource differences
	comparison.Differences = pa.calculateResourceDifferences(jobA, jobB)
	
	// Determine winner
	comparison.Winner = pa.determineWinner(comparison.Comparison)
	
	// Generate summary
	comparison.Summary = pa.generateComparisonSummary(comparison)
	
	return comparison, nil
}

// calculatePerformanceMetrics calculates comparative metrics between two jobs
func (pa *PerformanceAnalyzer) calculatePerformanceMetrics(
	jobA *interfaces.Job,
	analyticsA *interfaces.JobComprehensiveAnalytics,
	jobB *interfaces.Job,
	analyticsB *interfaces.JobComprehensiveAnalytics,
) PerformanceMetrics {
	metrics := PerformanceMetrics{
		ResourceWasteRatio: make(map[string]float64),
	}
	
	// Efficiency deltas using available fields
	metrics.OverallEfficiencyDelta = analyticsB.OverallEfficiency - analyticsA.OverallEfficiency
	
	if analyticsA.CPUAnalytics != nil && analyticsB.CPUAnalytics != nil {
		metrics.CPUEfficiencyDelta = analyticsB.CPUAnalytics.EfficiencyPercent - analyticsA.CPUAnalytics.EfficiencyPercent
	}
	
	if analyticsA.MemoryAnalytics != nil && analyticsB.MemoryAnalytics != nil {
		metrics.MemoryEfficiencyDelta = analyticsB.MemoryAnalytics.EfficiencyPercent - analyticsA.MemoryAnalytics.EfficiencyPercent
	}
	
	if analyticsA.IOAnalytics != nil && analyticsB.IOAnalytics != nil {
		metrics.IOEfficiencyDelta = analyticsB.IOAnalytics.EfficiencyPercent - analyticsA.IOAnalytics.EfficiencyPercent
	}
	
	// Runtime ratio
	runtimeA := pa.getJobRuntime(jobA)
	runtimeB := pa.getJobRuntime(jobB)
	if runtimeA > 0 {
		metrics.RuntimeRatio = float64(runtimeB) / float64(runtimeA)
	}
	
	// Resource waste ratios
	wasteA := pa.efficiencyCalc.CalculateResourceWaste(jobA, analyticsA, runtimeA)
	wasteB := pa.efficiencyCalc.CalculateResourceWaste(jobB, analyticsB, runtimeB)
	
	for resource, wasteValueA := range wasteA {
		if wasteValueB, exists := wasteB[resource]; exists && wasteValueA > 0 {
			metrics.ResourceWasteRatio[resource] = wasteValueB / wasteValueA
		}
	}
	
	// Cost efficiency ratio (performance per resource unit)
	costA := pa.calculateResourceCost(jobA, runtimeA)
	costB := pa.calculateResourceCost(jobB, runtimeB)
	
	if costA > 0 && costB > 0 {
		perfPerCostA := analyticsA.OverallEfficiency / costA
		perfPerCostB := analyticsB.OverallEfficiency / costB
		if perfPerCostA > 0 {
			metrics.CostEfficiencyRatio = perfPerCostB / perfPerCostA
		}
	}
	
	return metrics
}

// calculateResourceDifferences calculates the differences in resource allocation
func (pa *PerformanceAnalyzer) calculateResourceDifferences(jobA, jobB *interfaces.Job) ResourceDifferences {
	diffs := ResourceDifferences{
		CPUDelta:      jobB.CPUs - jobA.CPUs,
		MemoryDeltaGB: float64(jobB.Memory-jobA.Memory) / (1024 * 1024 * 1024),
	}
	
	// Check GPU allocation from metadata
	gpuA, _ := jobA.Metadata["gpus"].(int)
	gpuB, _ := jobB.Metadata["gpus"].(int)
	diffs.GPUDelta = gpuB - gpuA
	
	runtimeA := pa.getJobRuntime(jobA)
	runtimeB := pa.getJobRuntime(jobB)
	diffs.RuntimeDelta = runtimeB - runtimeA
	
	return diffs
}

// determineWinner determines which job performed better overall
func (pa *PerformanceAnalyzer) determineWinner(metrics PerformanceMetrics) string {
	score := 0.0
	
	// Efficiency is most important (40% weight)
	if metrics.OverallEfficiencyDelta > 5.0 {
		score += 0.4
	} else if metrics.OverallEfficiencyDelta < -5.0 {
		score -= 0.4
	}
	
	// Runtime is important (30% weight)
	if metrics.RuntimeRatio < 0.9 {
		score += 0.3
	} else if metrics.RuntimeRatio > 1.1 {
		score -= 0.3
	}
	
	// Cost efficiency is important (30% weight)
	if metrics.CostEfficiencyRatio > 1.1 {
		score += 0.3
	} else if metrics.CostEfficiencyRatio < 0.9 {
		score -= 0.3
	}
	
	if score > 0.2 {
		return "B"
	} else if score < -0.2 {
		return "A"
	}
	return "tie"
}

// generateComparisonSummary generates a human-readable summary
func (pa *PerformanceAnalyzer) generateComparisonSummary(comparison *JobPerformanceComparison) string {
	winner := comparison.Winner
	metrics := comparison.Comparison
	diffs := comparison.Differences
	
	summary := fmt.Sprintf("Job Performance Comparison:\n")
	summary += fmt.Sprintf("Job A: %s vs Job B: %s\n\n", comparison.JobA.ID, comparison.JobB.ID)
	
	// Resource differences
	summary += "Resource Allocation:\n"
	summary += fmt.Sprintf("  CPU: A=%d cores, B=%d cores (diff: %+d)\n", 
		comparison.JobA.CPUs, comparison.JobB.CPUs, diffs.CPUDelta)
	summary += fmt.Sprintf("  Memory: A=%.1fGB, B=%.1fGB (diff: %+.1fGB)\n",
		float64(comparison.JobA.Memory)/(1024*1024*1024),
		float64(comparison.JobB.Memory)/(1024*1024*1024),
		diffs.MemoryDeltaGB)
	gpuA, _ := comparison.JobA.Metadata["gpus"].(int)
	gpuB, _ := comparison.JobB.Metadata["gpus"].(int)
	if gpuA > 0 || gpuB > 0 {
		summary += fmt.Sprintf("  GPU: A=%d, B=%d (diff: %+d)\n",
			gpuA, gpuB, diffs.GPUDelta)
	}
	summary += fmt.Sprintf("  Runtime: A=%v, B=%v (ratio: %.2f)\n\n",
		pa.getJobRuntime(comparison.JobA),
		pa.getJobRuntime(comparison.JobB),
		metrics.RuntimeRatio)
	
	// Efficiency comparison
	summary += "Efficiency Metrics:\n"
	summary += fmt.Sprintf("  Overall: %+.1f%% (B is %.1f%% %s)\n",
		metrics.OverallEfficiencyDelta,
		math.Abs(metrics.OverallEfficiencyDelta),
		pa.betterOrWorse(metrics.OverallEfficiencyDelta))
	summary += fmt.Sprintf("  CPU: %+.1f%%\n", metrics.CPUEfficiencyDelta)
	summary += fmt.Sprintf("  Memory: %+.1f%%\n", metrics.MemoryEfficiencyDelta)
	gpuA2, _ := comparison.JobA.Metadata["gpus"].(int)
	gpuB2, _ := comparison.JobB.Metadata["gpus"].(int)
	if gpuA2 > 0 || gpuB2 > 0 {
		summary += fmt.Sprintf("  GPU: %+.1f%%\n", metrics.GPUEfficiencyDelta)
	}
	summary += fmt.Sprintf("  I/O: %+.1f%%\n\n", metrics.IOEfficiencyDelta)
	
	// Winner
	if winner == "tie" {
		summary += "Result: Both jobs performed similarly\n"
	} else if winner == "A" {
		summary += fmt.Sprintf("Result: Job A (%s) performed better overall\n", comparison.JobA.ID)
	} else {
		summary += fmt.Sprintf("Result: Job B (%s) performed better overall\n", comparison.JobB.ID)
	}
	
	// Key insights
	summary += "\nKey Insights:\n"
	if metrics.RuntimeRatio < 0.8 {
		summary += "  - Job B completed significantly faster\n"
	} else if metrics.RuntimeRatio > 1.2 {
		summary += "  - Job A completed significantly faster\n"
	}
	
	if math.Abs(metrics.OverallEfficiencyDelta) > 10.0 {
		better := "A"
		if metrics.OverallEfficiencyDelta > 0 {
			better = "B"
		}
		summary += fmt.Sprintf("  - Job %s has significantly better resource efficiency\n", better)
	}
	
	if metrics.CostEfficiencyRatio > 1.2 {
		summary += "  - Job B provides better performance per resource unit\n"
	} else if metrics.CostEfficiencyRatio < 0.8 {
		summary += "  - Job A provides better performance per resource unit\n"
	}
	
	return summary
}

// SimilarJobsAnalysis represents analysis of similar jobs
type SimilarJobsAnalysis struct {
	ReferenceJob     *interfaces.Job
	SimilarJobs      []*SimilarJobResult
	PerformanceStats PerformanceStatistics
	BestPractices    []string
	Recommendations  []string
}

// SimilarJobResult represents a similar job and its comparison
type SimilarJobResult struct {
	Job              *interfaces.Job
	Analytics        *interfaces.JobComprehensiveAnalytics
	SimilarityScore  float64 // 0-1, higher is more similar
	PerformanceRank  int     // 1 is best
	EfficiencyDelta  float64 // vs reference job
	RuntimeRatio     float64 // vs reference job
}

// PerformanceStatistics contains aggregate statistics for similar jobs
type PerformanceStatistics struct {
	AverageEfficiency float64
	MedianEfficiency  float64
	StdDevEfficiency  float64
	BestEfficiency    float64
	WorstEfficiency   float64
	AverageRuntime    time.Duration
	MedianRuntime     time.Duration
	OptimalResources  ResourceRecommendation
}

// ResourceRecommendation contains recommended resource allocations
type ResourceRecommendation struct {
	CPUs      int
	MemoryGB  float64
	GPUs      int
	Reasoning string
}

// GetSimilarJobsPerformance analyzes performance of similar jobs
func (pa *PerformanceAnalyzer) GetSimilarJobsPerformance(
	referenceJob *interfaces.Job,
	referenceAnalytics *interfaces.JobComprehensiveAnalytics,
	candidateJobs []struct {
		Job       *interfaces.Job
		Analytics *interfaces.JobComprehensiveAnalytics
	},
	similarityThreshold float64,
) (*SimilarJobsAnalysis, error) {
	if referenceJob == nil {
		return nil, fmt.Errorf("reference job must be provided")
	}
	if len(candidateJobs) == 0 {
		return nil, fmt.Errorf("no candidate jobs provided")
	}
	
	analysis := &SimilarJobsAnalysis{
		ReferenceJob: referenceJob,
		SimilarJobs:  make([]*SimilarJobResult, 0),
	}
	
	// Find similar jobs
	for _, candidate := range candidateJobs {
		similarity := pa.calculateJobSimilarity(referenceJob, candidate.Job)
		if similarity >= similarityThreshold {
			result := &SimilarJobResult{
				Job:             candidate.Job,
				Analytics:       candidate.Analytics,
				SimilarityScore: similarity,
			}
			
			// Calculate performance comparison
			if referenceAnalytics != nil && candidate.Analytics != nil {
				result.EfficiencyDelta = candidate.Analytics.OverallEfficiency - referenceAnalytics.OverallEfficiency
			}
			
			// Calculate runtime ratio
			refRuntime := pa.getJobRuntime(referenceJob)
			candRuntime := pa.getJobRuntime(candidate.Job)
			if refRuntime > 0 {
				result.RuntimeRatio = float64(candRuntime) / float64(refRuntime)
			}
			
			analysis.SimilarJobs = append(analysis.SimilarJobs, result)
		}
	}
	
	if len(analysis.SimilarJobs) == 0 {
		return nil, fmt.Errorf("no similar jobs found with threshold %.2f", similarityThreshold)
	}
	
	// Sort by efficiency (best first)
	sort.Slice(analysis.SimilarJobs, func(i, j int) bool {
		effI := float64(0)
		effJ := float64(0)
		if analysis.SimilarJobs[i].Analytics != nil {
			effI = analysis.SimilarJobs[i].Analytics.OverallEfficiency
		}
		if analysis.SimilarJobs[j].Analytics != nil {
			effJ = analysis.SimilarJobs[j].Analytics.OverallEfficiency
		}
		return effI > effJ
	})
	
	// Assign ranks
	for i := range analysis.SimilarJobs {
		analysis.SimilarJobs[i].PerformanceRank = i + 1
	}
	
	// Calculate statistics
	analysis.PerformanceStats = pa.calculatePerformanceStatistics(analysis.SimilarJobs)
	
	// Generate best practices
	analysis.BestPractices = pa.identifyBestPractices(analysis.SimilarJobs)
	
	// Generate recommendations
	analysis.Recommendations = pa.generateRecommendations(referenceJob, referenceAnalytics, analysis)
	
	return analysis, nil
}

// calculateJobSimilarity calculates how similar two jobs are (0-1)
func (pa *PerformanceAnalyzer) calculateJobSimilarity(jobA, jobB *interfaces.Job) float64 {
	similarity := 0.0
	weights := 0.0
	
	// Name similarity (10% weight)
	if jobA.Name == jobB.Name {
		similarity += 0.1
	}
	weights += 0.1
	
	// User similarity (15% weight)
	if jobA.UserID == jobB.UserID {
		similarity += 0.15
	}
	weights += 0.15
	
	// Partition similarity (15% weight)
	if jobA.Partition == jobB.Partition {
		similarity += 0.15
	}
	weights += 0.15
	
	// CPU similarity (20% weight)
	cpuRatio := float64(min(jobA.CPUs, jobB.CPUs)) / float64(max(jobA.CPUs, jobB.CPUs))
	similarity += 0.2 * cpuRatio
	weights += 0.2
	
	// Memory similarity (20% weight)
	memRatio := float64(min(jobA.Memory, jobB.Memory)) / float64(max(jobA.Memory, jobB.Memory))
	similarity += 0.2 * memRatio
	weights += 0.2
	
	// GPU similarity (10% weight)
	gpuA, _ := jobA.Metadata["gpus"].(int)
	gpuB, _ := jobB.Metadata["gpus"].(int)
	if gpuA == gpuB {
		similarity += 0.1
	} else if gpuA > 0 && gpuB > 0 {
		minGPU := gpuA
		maxGPU := gpuB
		if gpuB < gpuA {
			minGPU = gpuB
			maxGPU = gpuA
		}
		gpuRatio := float64(minGPU) / float64(maxGPU)
		similarity += 0.1 * gpuRatio
	}
	weights += 0.1
	
	// Runtime similarity (10% weight)
	runtimeA := pa.getJobRuntime(jobA)
	runtimeB := pa.getJobRuntime(jobB)
	if runtimeA > 0 && runtimeB > 0 {
		minRuntime := runtimeA
		maxRuntime := runtimeB
		if runtimeB < runtimeA {
			minRuntime = runtimeB
			maxRuntime = runtimeA
		}
		runtimeRatio := float64(minRuntime) / float64(maxRuntime)
		similarity += 0.1 * runtimeRatio
	}
	weights += 0.1
	
	if weights > 0 {
		return similarity / weights
	}
	return 0.0
}

// calculatePerformanceStatistics calculates aggregate statistics
func (pa *PerformanceAnalyzer) calculatePerformanceStatistics(similarJobs []*SimilarJobResult) PerformanceStatistics {
	stats := PerformanceStatistics{}
	
	if len(similarJobs) == 0 {
		return stats
	}
	
	efficiencies := make([]float64, 0)
	runtimes := make([]time.Duration, 0)
	cpuAllocations := make([]int, 0)
	memoryAllocations := make([]int, 0)
	
	for _, job := range similarJobs {
		if job.Analytics != nil {
			eff := job.Analytics.OverallEfficiency
			efficiencies = append(efficiencies, eff)
			
			if eff > stats.BestEfficiency {
				stats.BestEfficiency = eff
			}
			if stats.WorstEfficiency == 0 || eff < stats.WorstEfficiency {
				stats.WorstEfficiency = eff
			}
		}
		
		runtime := pa.getJobRuntime(job.Job)
		if runtime > 0 {
			runtimes = append(runtimes, runtime)
		}
		
		cpuAllocations = append(cpuAllocations, job.Job.CPUs)
		memoryAllocations = append(memoryAllocations, job.Job.Memory)
	}
	
	// Calculate efficiency statistics
	if len(efficiencies) > 0 {
		stats.AverageEfficiency = pa.mean(efficiencies)
		stats.MedianEfficiency = pa.median(efficiencies)
		stats.StdDevEfficiency = pa.stdDev(efficiencies)
	}
	
	// Calculate runtime statistics
	if len(runtimes) > 0 {
		totalRuntime := time.Duration(0)
		for _, rt := range runtimes {
			totalRuntime += rt
		}
		stats.AverageRuntime = totalRuntime / time.Duration(len(runtimes))
		stats.MedianRuntime = pa.medianDuration(runtimes)
	}
	
	// Find optimal resources (from best performing job)
	if len(similarJobs) > 0 {
		bestJob := similarJobs[0] // Already sorted by efficiency
		stats.OptimalResources = ResourceRecommendation{
			CPUs:      bestJob.Job.CPUs,
			MemoryGB:  float64(bestJob.Job.Memory) / (1024 * 1024 * 1024),
			GPUs:      func() int { gpus, _ := bestJob.Job.Metadata["gpus"].(int); return gpus }(),
			Reasoning: fmt.Sprintf("Based on best performing similar job with %.1f%% efficiency",
				stats.BestEfficiency),
		}
	}
	
	return stats
}

// identifyBestPractices identifies patterns from high-performing similar jobs
func (pa *PerformanceAnalyzer) identifyBestPractices(similarJobs []*SimilarJobResult) []string {
	practices := []string{}
	
	if len(similarJobs) < 3 {
		return practices
	}
	
	// Analyze top 25% of jobs
	topCount := max(1, len(similarJobs)/4)
	topJobs := similarJobs[:topCount]
	
	// Check for common resource allocation patterns
	avgCPUs := 0
	avgMemoryGB := 0.0
	for _, job := range topJobs {
		avgCPUs += job.Job.CPUs
		avgMemoryGB += float64(job.Job.Memory) / (1024 * 1024 * 1024)
	}
	avgCPUs /= topCount
	avgMemoryGB /= float64(topCount)
	
	practices = append(practices, fmt.Sprintf(
		"Top performing jobs typically use %d CPUs and %.1f GB memory",
		avgCPUs, avgMemoryGB))
	
	// Check for efficiency patterns
	allHighCPUEff := true
	allHighMemEff := true
	for _, job := range topJobs {
		if job.Analytics != nil {
			if job.Analytics.CPUAnalytics != nil && job.Analytics.CPUAnalytics.EfficiencyPercent < 80.0 {
				allHighCPUEff = false
			}
			if job.Analytics.MemoryAnalytics != nil && job.Analytics.MemoryAnalytics.EfficiencyPercent < 80.0 {
				allHighMemEff = false
			}
		}
	}
	
	if allHighCPUEff {
		practices = append(practices, "High CPU efficiency (>80%) is consistently achieved")
	}
	if allHighMemEff {
		practices = append(practices, "High memory efficiency (>80%) is consistently achieved")
	}
	
	// Check for runtime patterns
	runtimes := make([]time.Duration, 0)
	for _, job := range topJobs {
		runtime := pa.getJobRuntime(job.Job)
		if runtime > 0 {
			runtimes = append(runtimes, runtime)
		}
	}
	if len(runtimes) > 0 {
		avgRuntime := time.Duration(0)
		for _, rt := range runtimes {
			avgRuntime += rt
		}
		avgRuntime /= time.Duration(len(runtimes))
		practices = append(practices, fmt.Sprintf(
			"Average runtime for top performers: %v", avgRuntime))
	}
	
	return practices
}

// generateRecommendations generates specific recommendations
func (pa *PerformanceAnalyzer) generateRecommendations(
	referenceJob *interfaces.Job,
	referenceAnalytics *interfaces.JobComprehensiveAnalytics,
	analysis *SimilarJobsAnalysis,
) []string {
	recommendations := []string{}
	
	// Compare to optimal resources
	optimal := analysis.PerformanceStats.OptimalResources
	
	if referenceJob.CPUs > optimal.CPUs*2 {
		recommendations = append(recommendations, fmt.Sprintf(
			"Consider reducing CPU allocation to %d cores (currently %d)",
			optimal.CPUs, referenceJob.CPUs))
	} else if referenceJob.CPUs < optimal.CPUs/2 {
		recommendations = append(recommendations, fmt.Sprintf(
			"Consider increasing CPU allocation to %d cores (currently %d)",
			optimal.CPUs, referenceJob.CPUs))
	}
	
	refMemoryGB := float64(referenceJob.Memory) / (1024 * 1024 * 1024)
	if refMemoryGB > optimal.MemoryGB*2 {
		recommendations = append(recommendations, fmt.Sprintf(
			"Consider reducing memory allocation to %.1f GB (currently %.1f GB)",
			optimal.MemoryGB, refMemoryGB))
	} else if refMemoryGB < optimal.MemoryGB/2 {
		recommendations = append(recommendations, fmt.Sprintf(
			"Consider increasing memory allocation to %.1f GB (currently %.1f GB)",
			optimal.MemoryGB, refMemoryGB))
	}
	
	// Compare efficiency to statistics
	if referenceAnalytics != nil {
		refEff := referenceAnalytics.OverallEfficiency
		if refEff < analysis.PerformanceStats.MedianEfficiency-10.0 {
			recommendations = append(recommendations, fmt.Sprintf(
				"Job efficiency (%.1f%%) is below median (%.1f%%) for similar jobs",
				refEff, analysis.PerformanceStats.MedianEfficiency))
			
			// Add specific improvement suggestions
			if referenceAnalytics.CPUAnalytics != nil && referenceAnalytics.CPUAnalytics.EfficiencyPercent < 70.0 {
				recommendations = append(recommendations,
					"Focus on improving CPU utilization through better parallelization")
			}
			if referenceAnalytics.MemoryAnalytics != nil && referenceAnalytics.MemoryAnalytics.EfficiencyPercent < 70.0 {
				recommendations = append(recommendations,
					"Optimize memory usage patterns to reduce allocation needs")
			}
		}
	}
	
	// Runtime recommendations
	refRuntime := pa.getJobRuntime(referenceJob)
	if refRuntime > analysis.PerformanceStats.MedianRuntime*2 {
		recommendations = append(recommendations, fmt.Sprintf(
			"Job runtime (%v) is significantly longer than median (%v) for similar jobs",
			refRuntime, analysis.PerformanceStats.MedianRuntime))
	}
	
	return recommendations
}

// Helper functions

func (pa *PerformanceAnalyzer) getJobRuntime(job *interfaces.Job) time.Duration {
	if job.StartTime == nil || job.EndTime == nil {
		return 0
	}
	return job.EndTime.Sub(*job.StartTime)
}

func (pa *PerformanceAnalyzer) calculateResourceCost(job *interfaces.Job, runtime time.Duration) float64 {
	// Simple cost model: CPU-hours + Memory-GB-hours + GPU-hours
	hours := runtime.Hours()
	cpuCost := float64(job.CPUs) * hours
	memoryCost := float64(job.Memory) / (1024 * 1024 * 1024) * hours * 0.1 // Memory is cheaper
	gpuCount, _ := job.Metadata["gpus"].(int)
	gpuCost := float64(gpuCount) * hours * 10.0 // GPUs are expensive
	
	return cpuCost + memoryCost + gpuCost
}

func (pa *PerformanceAnalyzer) betterOrWorse(delta float64) string {
	if delta > 0 {
		return "better"
	}
	return "worse"
}

func (pa *PerformanceAnalyzer) mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func (pa *PerformanceAnalyzer) median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	
	if len(sorted)%2 == 0 {
		return (sorted[len(sorted)/2-1] + sorted[len(sorted)/2]) / 2
	}
	return sorted[len(sorted)/2]
}

func (pa *PerformanceAnalyzer) medianDuration(values []time.Duration) time.Duration {
	if len(values) == 0 {
		return 0
	}
	sorted := make([]time.Duration, len(values))
	copy(sorted, values)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})
	
	if len(sorted)%2 == 0 {
		return (sorted[len(sorted)/2-1] + sorted[len(sorted)/2]) / 2
	}
	return sorted[len(sorted)/2]
}

func (pa *PerformanceAnalyzer) stdDev(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	mean := pa.mean(values)
	variance := 0.0
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(values))
	return math.Sqrt(variance)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
