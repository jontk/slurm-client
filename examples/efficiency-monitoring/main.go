// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jontk/slurm-client/tests/mocks"
)

// EfficiencyThresholds defines optimization thresholds for different resources
type EfficiencyThresholds struct {
	CPU struct {
		Excellent float64 `json:"excellent"`
		Good      float64 `json:"good"`
		Poor      float64 `json:"poor"`
	} `json:"cpu"`
	Memory struct {
		Excellent float64 `json:"excellent"`
		Good      float64 `json:"good"`
		Poor      float64 `json:"poor"`
	} `json:"memory"`
	GPU struct {
		Excellent float64 `json:"excellent"`
		Good      float64 `json:"good"`
		Poor      float64 `json:"poor"`
	} `json:"gpu"`
	IO struct {
		Excellent float64 `json:"excellent"`
		Good      float64 `json:"good"`
		Poor      float64 `json:"poor"`
	} `json:"io"`
	Overall struct {
		Excellent float64 `json:"excellent"`
		Good      float64 `json:"good"`
		Poor      float64 `json:"poor"`
	} `json:"overall"`
}

// JobEfficiencyData holds complete efficiency information for a job
type JobEfficiencyData struct {
	JobID       string          `json:"job_id"`
	Timestamp   time.Time       `json:"timestamp"`
	Utilization UtilizationData `json:"utilization"`
	Efficiency  EfficiencyData  `json:"efficiency"`
	Performance PerformanceData `json:"performance"`
	LiveMetrics LiveMetricsData `json:"live_metrics,omitempty"`
}

// UtilizationData represents resource utilization metrics
type UtilizationData struct {
	CPUUtilization struct {
		AllocatedCores     int     `json:"allocated_cores"`
		UsedCores          float64 `json:"used_cores"`
		UtilizationPercent int     `json:"utilization_percent"`
		EfficiencyPercent  int     `json:"efficiency_percent"`
	} `json:"cpu_utilization"`
	MemoryUtilization struct {
		AllocatedBytes     int64 `json:"allocated_bytes"`
		UsedBytes          int64 `json:"used_bytes"`
		UtilizationPercent int   `json:"utilization_percent"`
		EfficiencyPercent  int   `json:"efficiency_percent"`
	} `json:"memory_utilization"`
	GPUUtilization struct {
		DeviceCount        int `json:"device_count"`
		UtilizationPercent int `json:"utilization_percent"`
	} `json:"gpu_utilization"`
	IOUtilization struct {
		ReadBytes          int64 `json:"read_bytes"`
		WriteBytes         int64 `json:"write_bytes"`
		UtilizationPercent int   `json:"utilization_percent"`
	} `json:"io_utilization"`
}

// EfficiencyData represents efficiency metrics and recommendations
type EfficiencyData struct {
	OverallEfficiencyScore float64 `json:"overall_efficiency_score"`
	CPUEfficiency          int     `json:"cpu_efficiency"`
	MemoryEfficiency       int     `json:"memory_efficiency"`
	GPUEfficiency          int     `json:"gpu_efficiency"`
	ResourceWaste          struct {
		CPUCoreHours  float64 `json:"cpu_core_hours"`
		CPUPercent    int     `json:"cpu_percent"`
		MemoryGBHours float64 `json:"memory_gb_hours"`
		MemoryPercent int     `json:"memory_percent"`
	} `json:"resource_waste"`
	OptimizationRecommendations []OptimizationRecommendation `json:"optimization_recommendations"`
}

// PerformanceData represents performance analytics
type PerformanceData struct {
	CPUAnalytics struct {
		AllocatedCores     int     `json:"allocated_cores"`
		UsedCores          float64 `json:"used_cores"`
		UtilizationPercent int     `json:"utilization_percent"`
		EfficiencyPercent  int     `json:"efficiency_percent"`
		AverageFrequency   int     `json:"average_frequency"`
		MaxFrequency       int     `json:"max_frequency"`
	} `json:"cpu_analytics"`
	MemoryAnalytics struct {
		AllocatedBytes     int64 `json:"allocated_bytes"`
		UsedBytes          int64 `json:"used_bytes"`
		UtilizationPercent int   `json:"utilization_percent"`
		EfficiencyPercent  int   `json:"efficiency_percent"`
	} `json:"memory_analytics"`
	IOAnalytics struct {
		ReadBytes             int64   `json:"read_bytes"`
		WriteBytes            int64   `json:"write_bytes"`
		ReadOperations        int     `json:"read_operations"`
		WriteOperations       int     `json:"write_operations"`
		AverageReadBandwidth  float64 `json:"average_read_bandwidth"`
		AverageWriteBandwidth float64 `json:"average_write_bandwidth"`
	} `json:"io_analytics"`
	OverallEfficiency float64 `json:"overall_efficiency"`
}

// LiveMetricsData represents real-time resource metrics
type LiveMetricsData struct {
	CPUUsage struct {
		Current     float64 `json:"current"`
		Average     float64 `json:"average"`
		Peak        float64 `json:"peak"`
		Utilization int     `json:"utilization"`
	} `json:"cpu_usage"`
	MemoryUsage struct {
		Current     int64 `json:"current"`
		Average     int64 `json:"average"`
		Peak        int64 `json:"peak"`
		Utilization int   `json:"utilization"`
	} `json:"memory_usage"`
	DiskUsage struct {
		ReadRateMbps  float64 `json:"read_rate_mbps"`
		WriteRateMbps float64 `json:"write_rate_mbps"`
	} `json:"disk_usage"`
	NetworkUsage struct {
		InRateMbps  float64 `json:"in_rate_mbps"`
		OutRateMbps float64 `json:"out_rate_mbps"`
	} `json:"network_usage"`
	Timestamp int64 `json:"timestamp"`
}

// OptimizationRecommendation represents a specific optimization suggestion
type OptimizationRecommendation struct {
	Type        string  `json:"type"`
	Resource    string  `json:"resource"`
	Current     int     `json:"current"`
	Recommended int     `json:"recommended"`
	Reason      string  `json:"reason"`
	Confidence  float64 `json:"confidence"`
}

// EfficiencyReport contains the complete efficiency analysis
type EfficiencyReport struct {
	Timestamp         time.Time               `json:"timestamp"`
	TotalJobsAnalyzed int                     `json:"total_jobs_analyzed"`
	OverallSummary    EfficiencySummary       `json:"overall_summary"`
	JobAnalysis       []JobEfficiencyAnalysis `json:"job_analysis"`
	Recommendations   []GlobalRecommendation  `json:"recommendations"`
	TrendAnalysis     TrendAnalysis           `json:"trend_analysis"`
}

// EfficiencySummary provides high-level efficiency statistics
type EfficiencySummary struct {
	AverageEfficiency      float64               `json:"average_efficiency"`
	EfficiencyDistribution map[string]int        `json:"efficiency_distribution"`
	ResourceWasteAnalysis  ResourceWasteAnalysis `json:"resource_waste_analysis"`
	PotentialSavings       PotentialSavings      `json:"potential_savings"`
	TopIssues              []string              `json:"top_issues"`
}

// JobEfficiencyAnalysis contains detailed analysis for individual jobs
type JobEfficiencyAnalysis struct {
	JobID             string                       `json:"job_id"`
	EfficiencyScore   float64                      `json:"efficiency_score"`
	EfficiencyGrade   string                       `json:"efficiency_grade"`
	ResourceBreakdown map[string]float64           `json:"resource_breakdown"`
	Recommendations   []OptimizationRecommendation `json:"recommendations"`
	IssuesSummary     []string                     `json:"issues_summary"`
	Optimizations     []OptimizationOpportunity    `json:"optimizations"`
}

// GlobalRecommendation represents system-wide optimization recommendations
type GlobalRecommendation struct {
	Category    string   `json:"category"`
	Impact      string   `json:"impact"`
	Description string   `json:"description"`
	Actions     []string `json:"actions"`
	Priority    string   `json:"priority"`
	Confidence  float64  `json:"confidence"`
}

// ResourceWasteAnalysis details resource waste across the system
type ResourceWasteAnalysis struct {
	TotalCPUWasteHours float64            `json:"total_cpu_waste_hours"`
	TotalMemoryWasteGB float64            `json:"total_memory_waste_gb"`
	TotalGPUWasteHours float64            `json:"total_gpu_waste_hours"`
	WasteByResource    map[string]float64 `json:"waste_by_resource"`
	MostWastefulJobs   []string           `json:"most_wasteful_jobs"`
}

// PotentialSavings estimates potential cost and resource savings
type PotentialSavings struct {
	EstimatedCostSavings  float64 `json:"estimated_cost_savings"`
	CPUCoreSavings        int     `json:"cpu_core_savings"`
	MemoryGBSavings       float64 `json:"memory_gb_savings"`
	GPUSavings            int     `json:"gpu_savings"`
	OptimizationPotential float64 `json:"optimization_potential"`
}

// TrendAnalysis provides efficiency trends over time
type TrendAnalysis struct {
	EfficiencyTrend     string   `json:"efficiency_trend"`
	TrendDirection      string   `json:"trend_direction"`
	TrendMagnitude      float64  `json:"trend_magnitude"`
	PredictedEfficiency float64  `json:"predicted_efficiency"`
	SeasonalPatterns    []string `json:"seasonal_patterns"`
	RecommendedActions  []string `json:"recommended_actions"`
}

// OptimizationOpportunity represents specific optimization opportunities
type OptimizationOpportunity struct {
	Type                 string  `json:"type"`
	Description          string  `json:"description"`
	ImpactLevel          string  `json:"impact_level"`
	ImplementationEffort string  `json:"implementation_effort"`
	EstimatedSaving      float64 `json:"estimated_saving"`
	Priority             int     `json:"priority"`
}

// EfficiencyMonitor handles efficiency monitoring and analysis
type EfficiencyMonitor struct {
	mockServer *mocks.MockSlurmServer
	baseURL    string
	thresholds EfficiencyThresholds
	httpClient *http.Client
}

// NewEfficiencyMonitor creates a new efficiency monitoring instance
func NewEfficiencyMonitor() *EfficiencyMonitor {
	// Initialize default thresholds
	thresholds := EfficiencyThresholds{}
	thresholds.CPU.Excellent = 90.0
	thresholds.CPU.Good = 75.0
	thresholds.CPU.Poor = 50.0
	thresholds.Memory.Excellent = 85.0
	thresholds.Memory.Good = 70.0
	thresholds.Memory.Poor = 45.0
	thresholds.GPU.Excellent = 80.0
	thresholds.GPU.Good = 60.0
	thresholds.GPU.Poor = 30.0
	thresholds.IO.Excellent = 70.0
	thresholds.IO.Good = 50.0
	thresholds.IO.Poor = 25.0
	thresholds.Overall.Excellent = 85.0
	thresholds.Overall.Good = 70.0
	thresholds.Overall.Poor = 50.0

	return &EfficiencyMonitor{
		thresholds: thresholds,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// StartMonitoring initializes the mock server and begins monitoring
func (em *EfficiencyMonitor) StartMonitoring() error {
	fmt.Println("🔄 Starting SLURM Efficiency Monitoring System...")

	// Start mock SLURM server
	em.mockServer = mocks.NewMockSlurmServerForVersion("v0.0.42")
	em.baseURL = em.mockServer.URL()

	fmt.Printf("✅ Mock SLURM server started at: %s\n", em.baseURL)
	return nil
}

// StopMonitoring shuts down the monitoring system
func (em *EfficiencyMonitor) StopMonitoring() {
	if em.mockServer != nil {
		em.mockServer.Close()
		fmt.Println("🔒 Efficiency monitoring system stopped")
	}
}

// CollectJobEfficiencyData gathers comprehensive efficiency data for a job
func (em *EfficiencyMonitor) CollectJobEfficiencyData(jobID string) (*JobEfficiencyData, error) {
	fmt.Printf("📊 Collecting efficiency data for job %s...\n", jobID)

	data := &JobEfficiencyData{
		JobID:     jobID,
		Timestamp: time.Now(),
	}

	// Collect utilization data
	utilizationURL := fmt.Sprintf("%s/slurm/v0.0.42/job/%s/utilization", em.baseURL, jobID)
	if err := em.fetchAndParseJSON(utilizationURL, &data.Utilization); err != nil {
		return nil, fmt.Errorf("failed to fetch utilization data: %w", err)
	}

	// Collect efficiency data
	efficiencyURL := fmt.Sprintf("%s/slurm/v0.0.42/job/%s/efficiency", em.baseURL, jobID)
	if err := em.fetchAndParseJSON(efficiencyURL, &data.Efficiency); err != nil {
		return nil, fmt.Errorf("failed to fetch efficiency data: %w", err)
	}

	// Collect performance data
	performanceURL := fmt.Sprintf("%s/slurm/v0.0.42/job/%s/performance", em.baseURL, jobID)
	if err := em.fetchAndParseJSON(performanceURL, &data.Performance); err != nil {
		return nil, fmt.Errorf("failed to fetch performance data: %w", err)
	}

	// Collect live metrics (optional)
	liveMetricsURL := fmt.Sprintf("%s/slurm/v0.0.42/job/%s/live_metrics", em.baseURL, jobID)
	if err := em.fetchAndParseJSON(liveMetricsURL, &data.LiveMetrics); err != nil {
		// Live metrics might not be available, continue without error
		fmt.Printf("ℹ️  Live metrics not available for job %s\n", jobID)
	}

	fmt.Printf("✅ Efficiency data collection completed for job %s\n", jobID)
	return data, nil
}

// fetchAndParseJSON performs HTTP request and JSON parsing
func (em *EfficiencyMonitor) fetchAndParseJSON(url string, target interface{}) error {
	resp, err := em.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Parse the nested response structure from mock server
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	// Extract the specific data section and re-marshal for target parsing
	var dataToUnmarshal []byte
	if strings.Contains(url, "/utilization") {
		if util, ok := response["utilization"]; ok {
			dataToUnmarshal, err = json.Marshal(util)
		} else {
			return fmt.Errorf("utilization data not found in response")
		}
	} else if strings.Contains(url, "/efficiency") {
		if eff, ok := response["efficiency"]; ok {
			dataToUnmarshal, err = json.Marshal(eff)
		} else {
			return fmt.Errorf("efficiency data not found in response")
		}
	} else if strings.Contains(url, "/performance") {
		if perf, ok := response["performance"]; ok {
			dataToUnmarshal, err = json.Marshal(perf)
		} else {
			return fmt.Errorf("performance data not found in response")
		}
	} else if strings.Contains(url, "/live_metrics") {
		if live, ok := response["live_metrics"]; ok {
			dataToUnmarshal, err = json.Marshal(live)
		} else {
			return fmt.Errorf("live_metrics data not found in response")
		}
	} else {
		dataToUnmarshal = body
	}

	if err != nil {
		return err
	}

	return json.Unmarshal(dataToUnmarshal, target)
}

// AnalyzeJobEfficiency performs detailed efficiency analysis
func (em *EfficiencyMonitor) AnalyzeJobEfficiency(data *JobEfficiencyData) *JobEfficiencyAnalysis {
	fmt.Printf("🔍 Analyzing efficiency for job %s...\n", data.JobID)

	analysis := &JobEfficiencyAnalysis{
		JobID:           data.JobID,
		EfficiencyScore: data.Efficiency.OverallEfficiencyScore,
		ResourceBreakdown: map[string]float64{
			"CPU":    float64(data.Efficiency.CPUEfficiency),
			"Memory": float64(data.Efficiency.MemoryEfficiency),
			"GPU":    float64(data.Efficiency.GPUEfficiency),
			"IO":     float64(data.Utilization.IOUtilization.UtilizationPercent),
		},
		Recommendations: data.Efficiency.OptimizationRecommendations,
	}

	// Determine efficiency grade
	analysis.EfficiencyGrade = em.determineEfficiencyGrade(data.Efficiency.OverallEfficiencyScore)

	// Generate issues summary
	analysis.IssuesSummary = em.identifyEfficiencyIssues(data)

	// Generate optimization opportunities
	analysis.Optimizations = em.generateOptimizationOpportunities(data)

	fmt.Printf("✅ Efficiency analysis completed - Grade: %s (%.1f%%)\n",
		analysis.EfficiencyGrade, analysis.EfficiencyScore)

	return analysis
}

// determineEfficiencyGrade assigns a grade based on efficiency score
func (em *EfficiencyMonitor) determineEfficiencyGrade(score float64) string {
	switch {
	case score >= em.thresholds.Overall.Excellent:
		return "A (Excellent)"
	case score >= em.thresholds.Overall.Good:
		return "B (Good)"
	case score >= em.thresholds.Overall.Poor:
		return "C (Needs Improvement)"
	default:
		return "D (Poor)"
	}
}

// identifyEfficiencyIssues identifies specific efficiency problems
func (em *EfficiencyMonitor) identifyEfficiencyIssues(data *JobEfficiencyData) []string {
	var issues []string

	// Check CPU efficiency
	if float64(data.Efficiency.CPUEfficiency) < em.thresholds.CPU.Good {
		issues = append(issues, fmt.Sprintf("Low CPU efficiency (%d%% < %.0f%%)",
			data.Efficiency.CPUEfficiency, em.thresholds.CPU.Good))
	}

	// Check memory efficiency
	if float64(data.Efficiency.MemoryEfficiency) < em.thresholds.Memory.Good {
		issues = append(issues, fmt.Sprintf("Low memory efficiency (%d%% < %.0f%%)",
			data.Efficiency.MemoryEfficiency, em.thresholds.Memory.Good))
	}

	// Check GPU efficiency (if applicable)
	if data.Utilization.GPUUtilization.DeviceCount > 0 &&
		float64(data.Efficiency.GPUEfficiency) < em.thresholds.GPU.Good {
		issues = append(issues, fmt.Sprintf("Low GPU efficiency (%d%% < %.0f%%)",
			data.Efficiency.GPUEfficiency, em.thresholds.GPU.Good))
	}

	// Check I/O efficiency
	if float64(data.Utilization.IOUtilization.UtilizationPercent) < em.thresholds.IO.Good {
		issues = append(issues, fmt.Sprintf("Low I/O efficiency (%d%% < %.0f%%)",
			data.Utilization.IOUtilization.UtilizationPercent, em.thresholds.IO.Good))
	}

	// Check resource waste
	if data.Efficiency.ResourceWaste.CPUPercent > 20 {
		issues = append(issues, fmt.Sprintf("High CPU waste (%d%%)",
			data.Efficiency.ResourceWaste.CPUPercent))
	}

	if data.Efficiency.ResourceWaste.MemoryPercent > 25 {
		issues = append(issues, fmt.Sprintf("High memory waste (%d%%)",
			data.Efficiency.ResourceWaste.MemoryPercent))
	}

	if len(issues) == 0 {
		issues = append(issues, "No significant efficiency issues detected")
	}

	return issues
}

// generateOptimizationOpportunities creates specific optimization suggestions
func (em *EfficiencyMonitor) generateOptimizationOpportunities(data *JobEfficiencyData) []OptimizationOpportunity {
	var opportunities []OptimizationOpportunity

	// CPU optimization opportunities
	if float64(data.Efficiency.CPUEfficiency) < em.thresholds.CPU.Good {
		opportunities = append(opportunities, OptimizationOpportunity{
			Type: "CPU Resource Optimization",
			Description: fmt.Sprintf("Reduce CPU allocation from %d cores based on %d%% utilization",
				data.Utilization.CPUUtilization.AllocatedCores, data.Efficiency.CPUEfficiency),
			ImpactLevel:          "High",
			ImplementationEffort: "Low",
			EstimatedSaving:      float64(data.Efficiency.ResourceWaste.CPUPercent) * 0.5,
			Priority:             1,
		})
	}

	// Memory optimization opportunities
	if float64(data.Efficiency.MemoryEfficiency) < em.thresholds.Memory.Good {
		memoryWasteGB := float64(data.Utilization.MemoryUtilization.AllocatedBytes-
			data.Utilization.MemoryUtilization.UsedBytes) / (1024 * 1024 * 1024)
		opportunities = append(opportunities, OptimizationOpportunity{
			Type:                 "Memory Resource Optimization",
			Description:          fmt.Sprintf("Reduce memory allocation by %.1fGB based on usage patterns", memoryWasteGB),
			ImpactLevel:          "Medium",
			ImplementationEffort: "Low",
			EstimatedSaving:      float64(data.Efficiency.ResourceWaste.MemoryPercent) * 0.3,
			Priority:             2,
		})
	}

	// I/O optimization opportunities
	if float64(data.Utilization.IOUtilization.UtilizationPercent) < em.thresholds.IO.Poor {
		opportunities = append(opportunities, OptimizationOpportunity{
			Type:                 "I/O Performance Optimization",
			Description:          "Consider I/O optimization techniques or alternative storage solutions",
			ImpactLevel:          "Medium",
			ImplementationEffort: "Medium",
			EstimatedSaving:      15.0,
			Priority:             3,
		})
	}

	// Overall efficiency opportunity
	if data.Efficiency.OverallEfficiencyScore < em.thresholds.Overall.Good {
		opportunities = append(opportunities, OptimizationOpportunity{
			Type:                 "Job Configuration Review",
			Description:          "Comprehensive review of job resource requirements and algorithm efficiency",
			ImpactLevel:          "High",
			ImplementationEffort: "High",
			EstimatedSaving:      (em.thresholds.Overall.Good - data.Efficiency.OverallEfficiencyScore) * 0.7,
			Priority:             1,
		})
	}

	// Sort opportunities by priority
	sort.Slice(opportunities, func(i, j int) bool {
		return opportunities[i].Priority < opportunities[j].Priority
	})

	return opportunities
}

// GenerateEfficiencyReport creates comprehensive efficiency report
func (em *EfficiencyMonitor) GenerateEfficiencyReport(jobIDs []string) (*EfficiencyReport, error) {
	fmt.Printf("📈 Generating efficiency report for %d jobs...\n", len(jobIDs))

	report := &EfficiencyReport{
		Timestamp:         time.Now(),
		TotalJobsAnalyzed: len(jobIDs),
		JobAnalysis:       make([]JobEfficiencyAnalysis, 0, len(jobIDs)),
	}

	var totalEfficiency float64
	efficiencyDistribution := map[string]int{
		"Excellent": 0,
		"Good":      0,
		"Poor":      0,
		"Very Poor": 0,
	}

	resourceWaste := ResourceWasteAnalysis{
		WasteByResource:  make(map[string]float64),
		MostWastefulJobs: make([]string, 0),
	}

	// Analyze each job
	for i, jobID := range jobIDs {
		fmt.Printf("  Processing job %d/%d: %s\n", i+1, len(jobIDs), jobID)

		jobData, err := em.CollectJobEfficiencyData(jobID)
		if err != nil {
			log.Printf("Warning: Failed to collect data for job %s: %v", jobID, err)
			continue
		}

		analysis := em.AnalyzeJobEfficiency(jobData)
		report.JobAnalysis = append(report.JobAnalysis, *analysis)

		// Update statistics
		totalEfficiency += analysis.EfficiencyScore

		// Update distribution
		switch analysis.EfficiencyGrade {
		case "A (Excellent)":
			efficiencyDistribution["Excellent"]++
		case "B (Good)":
			efficiencyDistribution["Good"]++
		case "C (Needs Improvement)":
			efficiencyDistribution["Poor"]++
		default:
			efficiencyDistribution["Very Poor"]++
		}

		// Update resource waste tracking
		resourceWaste.TotalCPUWasteHours += jobData.Efficiency.ResourceWaste.CPUCoreHours
		resourceWaste.TotalMemoryWasteGB += jobData.Efficiency.ResourceWaste.MemoryGBHours

		if jobData.Efficiency.ResourceWaste.CPUPercent > 30 ||
			jobData.Efficiency.ResourceWaste.MemoryPercent > 35 {
			resourceWaste.MostWastefulJobs = append(resourceWaste.MostWastefulJobs, jobID)
		}
	}

	// Calculate summary statistics
	if len(report.JobAnalysis) > 0 {
		report.OverallSummary.AverageEfficiency = totalEfficiency / float64(len(report.JobAnalysis))
	}
	report.OverallSummary.EfficiencyDistribution = efficiencyDistribution

	// Calculate resource waste by type
	resourceWaste.WasteByResource["CPU"] = resourceWaste.TotalCPUWasteHours
	resourceWaste.WasteByResource["Memory"] = resourceWaste.TotalMemoryWasteGB
	report.OverallSummary.ResourceWasteAnalysis = resourceWaste

	// Generate potential savings
	report.OverallSummary.PotentialSavings = em.calculatePotentialSavings(report.JobAnalysis)

	// Identify top issues
	report.OverallSummary.TopIssues = em.identifyTopSystemIssues(report.JobAnalysis)

	// Generate global recommendations
	report.Recommendations = em.generateGlobalRecommendations(report)

	// Generate trend analysis
	report.TrendAnalysis = em.generateTrendAnalysis(report)

	fmt.Printf("✅ Efficiency report generated successfully\n")
	return report, nil
}

// calculatePotentialSavings estimates system-wide savings potential
func (em *EfficiencyMonitor) calculatePotentialSavings(jobAnalyses []JobEfficiencyAnalysis) PotentialSavings {
	var totalOptimizationPotential float64
	var cpuSavings, memorySavings float64

	for _, analysis := range jobAnalyses {
		for _, opt := range analysis.Optimizations {
			totalOptimizationPotential += opt.EstimatedSaving
			if strings.Contains(opt.Type, "CPU") {
				cpuSavings += opt.EstimatedSaving * 0.1 // Convert to core savings
			}
			if strings.Contains(opt.Type, "Memory") {
				memorySavings += opt.EstimatedSaving * 0.05 // Convert to GB savings
			}
		}
	}

	return PotentialSavings{
		EstimatedCostSavings:  totalOptimizationPotential * 50, // Estimate $50 per efficiency point
		CPUCoreSavings:        int(cpuSavings),
		MemoryGBSavings:       memorySavings,
		GPUSavings:            0, // No GPU jobs in current analysis
		OptimizationPotential: totalOptimizationPotential / float64(len(jobAnalyses)),
	}
}

// identifyTopSystemIssues finds the most common efficiency problems
func (em *EfficiencyMonitor) identifyTopSystemIssues(jobAnalyses []JobEfficiencyAnalysis) []string {
	issueFrequency := make(map[string]int)

	for _, analysis := range jobAnalyses {
		for _, issue := range analysis.IssuesSummary {
			if !strings.Contains(issue, "No significant") {
				issueFrequency[issue]++
			}
		}
	}

	// Sort issues by frequency
	type issueCount struct {
		issue string
		count int
	}
	var issues []issueCount
	for issue, count := range issueFrequency {
		issues = append(issues, issueCount{issue, count})
	}
	sort.Slice(issues, func(i, j int) bool {
		return issues[i].count > issues[j].count
	})

	// Return top 5 issues
	var topIssues []string
	for i, issue := range issues {
		if i >= 5 {
			break
		}
		topIssues = append(topIssues, fmt.Sprintf("%s (affects %d jobs)", issue.issue, issue.count))
	}

	return topIssues
}

// generateGlobalRecommendations creates system-wide optimization recommendations
func (em *EfficiencyMonitor) generateGlobalRecommendations(report *EfficiencyReport) []GlobalRecommendation {
	var recommendations []GlobalRecommendation

	avgEfficiency := report.OverallSummary.AverageEfficiency

	// System-wide efficiency recommendation
	if avgEfficiency < em.thresholds.Overall.Good {
		recommendations = append(recommendations, GlobalRecommendation{
			Category: "System Efficiency",
			Impact:   "High",
			Description: fmt.Sprintf("System average efficiency (%.1f%%) is below target (%.0f%%)",
				avgEfficiency, em.thresholds.Overall.Good),
			Actions: []string{
				"Implement resource request guidelines for users",
				"Provide efficiency training and best practices documentation",
				"Consider implementing resource limits based on historical usage",
				"Deploy automated efficiency monitoring and alerting",
			},
			Priority:   "High",
			Confidence: 0.9,
		})
	}

	// Resource waste recommendation
	wasteAnalysis := report.OverallSummary.ResourceWasteAnalysis
	if wasteAnalysis.TotalCPUWasteHours > 100 || wasteAnalysis.TotalMemoryWasteGB > 200 {
		recommendations = append(recommendations, GlobalRecommendation{
			Category: "Resource Waste Reduction",
			Impact:   "High",
			Description: fmt.Sprintf("Significant resource waste detected: %.1f CPU hours, %.1f GB memory",
				wasteAnalysis.TotalCPUWasteHours, wasteAnalysis.TotalMemoryWasteGB),
			Actions: []string{
				"Implement dynamic resource allocation policies",
				"Provide job profiling tools to help users optimize requests",
				"Create efficiency dashboards for user self-monitoring",
				"Consider implementing resource usage incentives",
			},
			Priority:   "High",
			Confidence: 0.85,
		})
	}

	// User education recommendation
	poorJobsCount := report.OverallSummary.EfficiencyDistribution["Poor"] +
		report.OverallSummary.EfficiencyDistribution["Very Poor"]
	if float64(poorJobsCount)/float64(report.TotalJobsAnalyzed) > 0.3 {
		recommendations = append(recommendations, GlobalRecommendation{
			Category: "User Education",
			Impact:   "Medium",
			Description: fmt.Sprintf("%d%% of jobs have poor efficiency",
				int(float64(poorJobsCount)/float64(report.TotalJobsAnalyzed)*100)),
			Actions: []string{
				"Develop efficiency best practices documentation",
				"Conduct user training sessions on resource optimization",
				"Create job templates for common workflows",
				"Implement efficiency scoring in job submission interface",
			},
			Priority:   "Medium",
			Confidence: 0.8,
		})
	}

	// Monitoring improvements
	recommendations = append(recommendations, GlobalRecommendation{
		Category:    "Monitoring Enhancement",
		Impact:      "Medium",
		Description: "Continuous efficiency monitoring will help maintain optimal resource utilization",
		Actions: []string{
			"Deploy automated efficiency reporting (weekly/monthly)",
			"Implement real-time efficiency alerts for jobs",
			"Create efficiency trends analysis and forecasting",
			"Integrate efficiency metrics into scheduler decisions",
		},
		Priority:   "Medium",
		Confidence: 0.75,
	})

	return recommendations
}

// generateTrendAnalysis creates efficiency trend analysis
func (em *EfficiencyMonitor) generateTrendAnalysis(report *EfficiencyReport) TrendAnalysis {
	avgEfficiency := report.OverallSummary.AverageEfficiency

	// Simulate trend analysis (in real implementation, this would use historical data)
	trend := TrendAnalysis{
		EfficiencyTrend:     "Stable",
		TrendDirection:      "Neutral",
		TrendMagnitude:      0.5,
		PredictedEfficiency: avgEfficiency,
	}

	if avgEfficiency > em.thresholds.Overall.Good {
		trend.EfficiencyTrend = "Improving"
		trend.TrendDirection = "Positive"
		trend.TrendMagnitude = 2.1
		trend.PredictedEfficiency = avgEfficiency + 2.0
	} else if avgEfficiency < em.thresholds.Overall.Poor {
		trend.EfficiencyTrend = "Concerning"
		trend.TrendDirection = "Negative"
		trend.TrendMagnitude = -1.5
		trend.PredictedEfficiency = avgEfficiency - 1.5
	}

	trend.SeasonalPatterns = []string{
		"Higher efficiency during regular business hours",
		"Lower efficiency during peak usage periods",
		"Resource waste increases during project deadlines",
	}

	trend.RecommendedActions = []string{
		"Monitor efficiency trends weekly",
		"Adjust scheduling policies based on efficiency patterns",
		"Provide proactive optimization recommendations",
		"Implement efficiency-based resource pricing",
	}

	return trend
}

// PrintEfficiencyReport displays the complete efficiency report
func (em *EfficiencyMonitor) PrintEfficiencyReport(report *EfficiencyReport) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 SLURM EFFICIENCY MONITORING REPORT")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("📅 Generated: %s\n", report.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("🎯 Jobs Analyzed: %d\n", report.TotalJobsAnalyzed)
	fmt.Printf("📈 Average Efficiency: %.1f%%\n", report.OverallSummary.AverageEfficiency)

	// Efficiency Distribution
	fmt.Println("\n📊 EFFICIENCY DISTRIBUTION:")
	fmt.Printf("  🟢 Excellent: %d jobs\n", report.OverallSummary.EfficiencyDistribution["Excellent"])
	fmt.Printf("  🟡 Good: %d jobs\n", report.OverallSummary.EfficiencyDistribution["Good"])
	fmt.Printf("  🟠 Poor: %d jobs\n", report.OverallSummary.EfficiencyDistribution["Poor"])
	fmt.Printf("  🔴 Very Poor: %d jobs\n", report.OverallSummary.EfficiencyDistribution["Very Poor"])

	// Resource Waste Analysis
	fmt.Println("\n💸 RESOURCE WASTE ANALYSIS:")
	waste := report.OverallSummary.ResourceWasteAnalysis
	fmt.Printf("  ⚡ Total CPU Waste: %.1f core-hours\n", waste.TotalCPUWasteHours)
	fmt.Printf("  💾 Total Memory Waste: %.1f GB-hours\n", waste.TotalMemoryWasteGB)
	if len(waste.MostWastefulJobs) > 0 {
		fmt.Printf("  🎯 Most Wasteful Jobs: %s\n", strings.Join(waste.MostWastefulJobs, ", "))
	}

	// Potential Savings
	fmt.Println("\n💰 POTENTIAL SAVINGS:")
	savings := report.OverallSummary.PotentialSavings
	fmt.Printf("  💵 Estimated Cost Savings: $%.2f\n", savings.EstimatedCostSavings)
	fmt.Printf("  ⚡ CPU Core Savings: %d cores\n", savings.CPUCoreSavings)
	fmt.Printf("  💾 Memory Savings: %.1f GB\n", savings.MemoryGBSavings)
	fmt.Printf("  📈 Optimization Potential: %.1f%%\n", savings.OptimizationPotential)

	// Top Issues
	fmt.Println("\n🔍 TOP SYSTEM ISSUES:")
	for i, issue := range report.OverallSummary.TopIssues {
		fmt.Printf("  %d. %s\n", i+1, issue)
	}

	// Global Recommendations
	fmt.Println("\n🎯 GLOBAL RECOMMENDATIONS:")
	for i, rec := range report.Recommendations {
		fmt.Printf("  %d. [%s] %s\n", i+1, rec.Priority, rec.Description)
		for j, action := range rec.Actions {
			if j < 2 { // Show only first 2 actions
				fmt.Printf("     • %s\n", action)
			}
		}
	}

	// Trend Analysis
	fmt.Println("\n📈 TREND ANALYSIS:")
	trend := report.TrendAnalysis
	fmt.Printf("  📊 Current Trend: %s (%s)\n", trend.EfficiencyTrend, trend.TrendDirection)
	fmt.Printf("  🔮 Predicted Efficiency: %.1f%%\n", trend.PredictedEfficiency)
	fmt.Printf("  📅 Key Pattern: %s\n", trend.SeasonalPatterns[0])

	// Individual Job Analysis (Top 5 best and worst)
	fmt.Println("\n🎯 INDIVIDUAL JOB ANALYSIS:")

	// Sort jobs by efficiency
	jobs := append([]JobEfficiencyAnalysis{}, report.JobAnalysis...)
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].EfficiencyScore > jobs[j].EfficiencyScore
	})

	fmt.Println("  🏆 TOP PERFORMING JOBS:")
	for i := 0; i < len(jobs) && i < 3; i++ {
		job := jobs[i]
		fmt.Printf("    %s: %.1f%% (%s)\n", job.JobID, job.EfficiencyScore, job.EfficiencyGrade)
		if len(job.IssuesSummary) > 0 && !strings.Contains(job.IssuesSummary[0], "No significant") {
			fmt.Printf("      Issue: %s\n", job.IssuesSummary[0])
		}
	}

	fmt.Println("  ⚠️  JOBS NEEDING ATTENTION:")
	for i := len(jobs) - 1; i >= 0 && i >= len(jobs)-3; i-- {
		job := jobs[i]
		fmt.Printf("    %s: %.1f%% (%s)\n", job.JobID, job.EfficiencyScore, job.EfficiencyGrade)
		if len(job.IssuesSummary) > 0 {
			fmt.Printf("      Issue: %s\n", job.IssuesSummary[0])
		}
		if len(job.Optimizations) > 0 {
			fmt.Printf("      Optimization: %s\n", job.Optimizations[0].Description)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
}

// ExportReportToJSON saves the efficiency report as JSON
func (em *EfficiencyMonitor) ExportReportToJSON(report *EfficiencyReport, filename string) error {
	fmt.Printf("💾 Exporting efficiency report to %s...\n", filename)

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("✅ Efficiency report exported successfully\n")
	return nil
}

func main() {
	fmt.Println("🚀 SLURM Efficiency Monitoring System")
	fmt.Println("=====================================")

	// Initialize efficiency monitor
	monitor := NewEfficiencyMonitor()

	// Start monitoring
	if err := monitor.StartMonitoring(); err != nil {
		log.Fatalf("Failed to start monitoring: %v", err)
	}
	defer monitor.StopMonitoring()

	// Sample job IDs for demonstration
	jobIDs := []string{"1001", "1002", "1003", "1004", "1005"}

	// Generate comprehensive efficiency report
	report, err := monitor.GenerateEfficiencyReport(jobIDs)
	if err != nil {
		log.Fatalf("Failed to generate efficiency report: %v", err)
	}

	// Print the efficiency report
	monitor.PrintEfficiencyReport(report)

	// Export detailed analysis
	reportFilename := fmt.Sprintf("efficiency_report_%s.json",
		time.Now().Format("2006-01-02_15-04-05"))

	if err := monitor.ExportReportToJSON(report, reportFilename); err != nil {
		log.Printf("Warning: Failed to export report: %v", err)
	}

	fmt.Printf("\n🎉 Efficiency monitoring completed successfully!\n")
	fmt.Printf("📁 Detailed report saved to: %s\n", reportFilename)
}
