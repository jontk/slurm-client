// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jontk/slurm-client/tests/mocks"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// JobAnalyticsData represents the analytics data for a job
type JobAnalyticsData struct {
	JobID       string           `json:"job_id"`
	Utilization *UtilizationData `json:"utilization,omitempty"`
	Efficiency  *EfficiencyData  `json:"efficiency,omitempty"`
	Performance *PerformanceData `json:"performance,omitempty"`
	LiveMetrics *LiveMetricsData `json:"live_metrics,omitempty"`
	Trends      *TrendsData      `json:"trends,omitempty"`
}

type UtilizationData struct {
	CPU    CPUUtilization    `json:"cpu_utilization"`
	Memory MemoryUtilization `json:"memory_utilization"`
	GPU    GPUUtilization    `json:"gpu_utilization"`
	IO     IOUtilization     `json:"io_utilization"`
}

type EfficiencyData struct {
	OverallEfficiency float64                      `json:"overall_efficiency_score"`
	CPUEfficiency     float64                      `json:"cpu_efficiency"`
	MemoryEfficiency  float64                      `json:"memory_efficiency"`
	GPUEfficiency     float64                      `json:"gpu_efficiency"`
	ResourceWaste     ResourceWaste                `json:"resource_waste"`
	Recommendations   []OptimizationRecommendation `json:"optimization_recommendations"`
}

type PerformanceData struct {
	CPUAnalytics      CPUAnalytics    `json:"cpu_analytics"`
	MemoryAnalytics   MemoryAnalytics `json:"memory_analytics"`
	IOAnalytics       IOAnalytics     `json:"io_analytics"`
	OverallEfficiency float64         `json:"overall_efficiency"`
}

type LiveMetricsData struct {
	CPUUsage     CPUUsage     `json:"cpu_usage"`
	MemoryUsage  MemoryUsage  `json:"memory_usage"`
	DiskUsage    DiskUsage    `json:"disk_usage"`
	NetworkUsage NetworkUsage `json:"network_usage"`
	Timestamp    int64        `json:"timestamp"`
}

type TrendsData struct {
	TimeWindow  string      `json:"time_window"`
	CPUTrend    TrendPoints `json:"cpu_trend"`
	MemoryTrend TrendPoints `json:"memory_trend"`
}

// Utility structures
type CPUUtilization struct {
	AllocatedCores     int     `json:"allocated_cores"`
	UsedCores          float64 `json:"used_cores"`
	UtilizationPercent float64 `json:"utilization_percent"`
	EfficiencyPercent  float64 `json:"efficiency_percent"`
}

type MemoryUtilization struct {
	AllocatedBytes     int64   `json:"allocated_bytes"`
	UsedBytes          int64   `json:"used_bytes"`
	UtilizationPercent float64 `json:"utilization_percent"`
	EfficiencyPercent  float64 `json:"efficiency_percent"`
}

type GPUUtilization struct {
	DeviceCount        int     `json:"device_count"`
	UtilizationPercent float64 `json:"utilization_percent"`
}

type IOUtilization struct {
	ReadBytes          int64   `json:"read_bytes"`
	WriteBytes         int64   `json:"write_bytes"`
	UtilizationPercent float64 `json:"utilization_percent"`
}

type ResourceWaste struct {
	CPUCoreHours  float64 `json:"cpu_core_hours"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryGBHours float64 `json:"memory_gb_hours"`
	MemoryPercent float64 `json:"memory_percent"`
}

type OptimizationRecommendation struct {
	Type        string  `json:"type"`
	Resource    string  `json:"resource"`
	Current     int     `json:"current"`
	Recommended int     `json:"recommended"`
	Reason      string  `json:"reason"`
	Confidence  float64 `json:"confidence"`
}

type CPUAnalytics struct {
	AllocatedCores     int     `json:"allocated_cores"`
	UsedCores          float64 `json:"used_cores"`
	UtilizationPercent float64 `json:"utilization_percent"`
	EfficiencyPercent  float64 `json:"efficiency_percent"`
	AverageFrequency   int     `json:"average_frequency"`
	MaxFrequency       int     `json:"max_frequency"`
}

type MemoryAnalytics struct {
	AllocatedBytes     int64   `json:"allocated_bytes"`
	UsedBytes          int64   `json:"used_bytes"`
	UtilizationPercent float64 `json:"utilization_percent"`
	EfficiencyPercent  float64 `json:"efficiency_percent"`
}

type IOAnalytics struct {
	ReadBytes             int64   `json:"read_bytes"`
	WriteBytes            int64   `json:"write_bytes"`
	ReadOperations        int     `json:"read_operations"`
	WriteOperations       int     `json:"write_operations"`
	AverageReadBandwidth  float64 `json:"average_read_bandwidth"`
	AverageWriteBandwidth float64 `json:"average_write_bandwidth"`
}

type CPUUsage struct {
	Current     float64 `json:"current"`
	Average     float64 `json:"average"`
	Peak        float64 `json:"peak"`
	Utilization float64 `json:"utilization"`
}

type MemoryUsage struct {
	Current     int64   `json:"current"`
	Average     int64   `json:"average"`
	Peak        int64   `json:"peak"`
	Utilization float64 `json:"utilization"`
}

type DiskUsage struct {
	ReadRateMBps  float64 `json:"read_rate_mbps"`
	WriteRateMBps float64 `json:"write_rate_mbps"`
}

type NetworkUsage struct {
	InRateMBps  float64 `json:"in_rate_mbps"`
	OutRateMBps float64 `json:"out_rate_mbps"`
}

type TrendPoints struct {
	DataPoints []TrendPoint `json:"data_points"`
}

type TrendPoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// Helper functions for safe type assertions
func getFloat64(data map[string]interface{}, key string) (float64, bool) {
	v, ok := data[key].(float64)
	return v, ok
}

func getString(data map[string]interface{}, key string) (string, bool) {
	v, ok := data[key].(string)
	return v, ok
}

func getMap(data map[string]interface{}, key string) (map[string]interface{}, bool) {
	v, ok := data[key].(map[string]interface{})
	return v, ok
}

func getSlice(data map[string]interface{}, key string) ([]interface{}, bool) {
	v, ok := data[key].([]interface{})
	return v, ok
}

// AnalyticsCollector collects and analyzes job analytics data
type AnalyticsCollector struct {
	baseURL    string
	httpClient *http.Client
}

// NewAnalyticsCollector creates a new analytics collector
func NewAnalyticsCollector(baseURL string) *AnalyticsCollector {
	return &AnalyticsCollector{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CollectJobAnalytics collects comprehensive analytics data for a job
func (ac *AnalyticsCollector) CollectJobAnalytics(ctx context.Context, jobID string) *JobAnalyticsData {
	analytics := &JobAnalyticsData{
		JobID: jobID,
	}

	// Collect utilization data
	utilization, err := ac.getJobUtilization(ctx, jobID)
	if err != nil {
		log.Printf("Warning: Failed to collect utilization data: %v", err)
	} else {
		analytics.Utilization = utilization
	}

	// Collect efficiency data
	efficiency, err := ac.getJobEfficiency(ctx, jobID)
	if err != nil {
		log.Printf("Warning: Failed to collect efficiency data: %v", err)
	} else {
		analytics.Efficiency = efficiency
	}

	// Collect performance data
	performance, err := ac.getJobPerformance(ctx, jobID)
	if err != nil {
		log.Printf("Warning: Failed to collect performance data: %v", err)
	} else {
		analytics.Performance = performance
	}

	// Collect live metrics data
	liveMetrics, err := ac.getJobLiveMetrics(ctx, jobID)
	if err != nil {
		log.Printf("Warning: Failed to collect live metrics data: %v", err)
	} else {
		analytics.LiveMetrics = liveMetrics
	}

	// Collect trends data
	trends, err := ac.getJobResourceTrends(ctx, jobID)
	if err != nil {
		log.Printf("Warning: Failed to collect trends data: %v", err)
	} else {
		analytics.Trends = trends
	}

	return analytics
}

func (ac *AnalyticsCollector) getJobUtilization(ctx context.Context, jobID string) (*UtilizationData, error) {
	url := fmt.Sprintf("%s/slurm/v0.0.42/job/%s/utilization", ac.baseURL, jobID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	resp, err := ac.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get utilization data: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	utilizationMap, ok := result["utilization"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid utilization response format")
	}

	// Parse the utilization data
	utilization := &UtilizationData{}

	if cpuData, ok := utilizationMap["cpu_utilization"].(map[string]interface{}); ok {
		if allocCores, ok := cpuData["allocated_cores"].(float64); ok {
			if usedCores, ok := cpuData["used_cores"].(float64); ok {
				if utilPct, ok := cpuData["utilization_percent"].(float64); ok {
					if effPct, ok := cpuData["efficiency_percent"].(float64); ok {
						utilization.CPU = CPUUtilization{
							AllocatedCores:     int(allocCores),
							UsedCores:          usedCores,
							UtilizationPercent: utilPct,
							EfficiencyPercent:  effPct,
						}
					}
				}
			}
		}
	}

	if memData, ok := utilizationMap["memory_utilization"].(map[string]interface{}); ok {
		if allocBytes, ok := memData["allocated_bytes"].(float64); ok {
			if usedBytes, ok := memData["used_bytes"].(float64); ok {
				if utilPct, ok := memData["utilization_percent"].(float64); ok {
					if effPct, ok := memData["efficiency_percent"].(float64); ok {
						utilization.Memory = MemoryUtilization{
							AllocatedBytes:     int64(allocBytes),
							UsedBytes:          int64(usedBytes),
							UtilizationPercent: utilPct,
							EfficiencyPercent:  effPct,
						}
					}
				}
			}
		}
	}

	if gpuData, ok := getMap(utilizationMap, "gpu_utilization"); ok {
		if deviceCount, ok := getFloat64(gpuData, "device_count"); ok {
			if utilPct, ok := getFloat64(gpuData, "utilization_percent"); ok {
				utilization.GPU = GPUUtilization{
					DeviceCount:        int(deviceCount),
					UtilizationPercent: utilPct,
				}
			}
		}
	}

	if ioData, ok := getMap(utilizationMap, "io_utilization"); ok {
		if readBytes, ok := getFloat64(ioData, "read_bytes"); ok {
			if writeBytes, ok := getFloat64(ioData, "write_bytes"); ok {
				if utilPct, ok := getFloat64(ioData, "utilization_percent"); ok {
					utilization.IO = IOUtilization{
						ReadBytes:          int64(readBytes),
						WriteBytes:         int64(writeBytes),
						UtilizationPercent: utilPct,
					}
				}
			}
		}
	}

	return utilization, nil
}

func (ac *AnalyticsCollector) getJobEfficiency(ctx context.Context, jobID string) (*EfficiencyData, error) {
	url := fmt.Sprintf("%s/slurm/v0.0.42/job/%s/efficiency", ac.baseURL, jobID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	resp, err := ac.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get efficiency data: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	efficiencyMap, ok := result["efficiency"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid efficiency response format")
	}

	efficiency := &EfficiencyData{}
	if overall, ok := efficiencyMap["overall_efficiency_score"].(float64); ok {
		efficiency.OverallEfficiency = overall
	}
	if cpu, ok := efficiencyMap["cpu_efficiency"].(float64); ok {
		efficiency.CPUEfficiency = cpu
	}
	if mem, ok := efficiencyMap["memory_efficiency"].(float64); ok {
		efficiency.MemoryEfficiency = mem
	}
	if gpu, ok := efficiencyMap["gpu_efficiency"].(float64); ok {
		efficiency.GPUEfficiency = gpu
	}

	// Parse resource waste
	if wasteData, ok := efficiencyMap["resource_waste"].(map[string]interface{}); ok {
		if cpuHours, ok := wasteData["cpu_core_hours"].(float64); ok {
			if cpuPct, ok := wasteData["cpu_percent"].(float64); ok {
				if memHours, ok := wasteData["memory_gb_hours"].(float64); ok {
					if memPct, ok := wasteData["memory_percent"].(float64); ok {
						efficiency.ResourceWaste = ResourceWaste{
							CPUCoreHours:  cpuHours,
							CPUPercent:    cpuPct,
							MemoryGBHours: memHours,
							MemoryPercent: memPct,
						}
					}
				}
			}
		}
	}

	// Parse recommendations
	if recsData, ok := efficiencyMap["optimization_recommendations"].([]interface{}); ok {
		for _, recData := range recsData {
			rec, ok := recData.(map[string]interface{})
			if !ok {
				continue
			}
			var recType, recResource, recReason string
			var current, recommended float64
			var confidence float64

			if v, ok := rec["type"].(string); ok {
				recType = v
			}
			if v, ok := rec["resource"].(string); ok {
				recResource = v
			}
			if v, ok := rec["current"].(float64); ok {
				current = v
			}
			if v, ok := rec["recommended"].(float64); ok {
				recommended = v
			}
			if v, ok := rec["reason"].(string); ok {
				recReason = v
			}
			if v, ok := rec["confidence"].(float64); ok {
				confidence = v
			}

			recommendation := OptimizationRecommendation{
				Type:        recType,
				Resource:    recResource,
				Current:     int(current),
				Recommended: int(recommended),
				Reason:      recReason,
				Confidence:  confidence,
			}
			efficiency.Recommendations = append(efficiency.Recommendations, recommendation)
		}
	}

	return efficiency, nil
}

func (ac *AnalyticsCollector) getJobPerformance(ctx context.Context, jobID string) (*PerformanceData, error) {
	url := fmt.Sprintf("%s/slurm/v0.0.42/job/%s/performance", ac.baseURL, jobID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	resp, err := ac.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get performance data: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	performanceMap, ok := result["performance"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid performance response format")
	}

	performance := &PerformanceData{}
	if overall, ok := getFloat64(performanceMap, "overall_efficiency"); ok {
		performance.OverallEfficiency = overall
	}

	// Parse CPU analytics
	if cpuData, ok := getMap(performanceMap, "cpu_analytics"); ok {
		if allocCores, ok := getFloat64(cpuData, "allocated_cores"); ok {
			if usedCores, ok := getFloat64(cpuData, "used_cores"); ok {
				if utilPct, ok := getFloat64(cpuData, "utilization_percent"); ok {
					if effPct, ok := getFloat64(cpuData, "efficiency_percent"); ok {
						if avgFreq, ok := getFloat64(cpuData, "average_frequency"); ok {
							if maxFreq, ok := getFloat64(cpuData, "max_frequency"); ok {
								performance.CPUAnalytics = CPUAnalytics{
									AllocatedCores:     int(allocCores),
									UsedCores:          usedCores,
									UtilizationPercent: utilPct,
									EfficiencyPercent:  effPct,
									AverageFrequency:   int(avgFreq),
									MaxFrequency:       int(maxFreq),
								}
							}
						}
					}
				}
			}
		}
	}

	// Parse Memory analytics
	if memData, ok := getMap(performanceMap, "memory_analytics"); ok {
		if allocBytes, ok := getFloat64(memData, "allocated_bytes"); ok {
			if usedBytes, ok := getFloat64(memData, "used_bytes"); ok {
				if utilPct, ok := getFloat64(memData, "utilization_percent"); ok {
					if effPct, ok := getFloat64(memData, "efficiency_percent"); ok {
						performance.MemoryAnalytics = MemoryAnalytics{
							AllocatedBytes:     int64(allocBytes),
							UsedBytes:          int64(usedBytes),
							UtilizationPercent: utilPct,
							EfficiencyPercent:  effPct,
						}
					}
				}
			}
		}
	}

	// Parse IO analytics
	if ioData, ok := getMap(performanceMap, "io_analytics"); ok {
		if readBytes, ok := getFloat64(ioData, "read_bytes"); ok {
			if writeBytes, ok := getFloat64(ioData, "write_bytes"); ok {
				if readOps, ok := getFloat64(ioData, "read_operations"); ok {
					if writeOps, ok := getFloat64(ioData, "write_operations"); ok {
						if avgReadBw, ok := getFloat64(ioData, "average_read_bandwidth"); ok {
							if avgWriteBw, ok := getFloat64(ioData, "average_write_bandwidth"); ok {
								performance.IOAnalytics = IOAnalytics{
									ReadBytes:             int64(readBytes),
									WriteBytes:            int64(writeBytes),
									ReadOperations:        int(readOps),
									WriteOperations:       int(writeOps),
									AverageReadBandwidth:  avgReadBw,
									AverageWriteBandwidth: avgWriteBw,
								}
							}
						}
					}
				}
			}
		}
	}

	return performance, nil
}

func (ac *AnalyticsCollector) getJobLiveMetrics(ctx context.Context, jobID string) (*LiveMetricsData, error) {
	url := fmt.Sprintf("%s/slurm/v0.0.42/job/%s/live_metrics", ac.baseURL, jobID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	resp, err := ac.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get live metrics data: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	metricsMap, ok := result["live_metrics"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid live metrics response format")
	}

	liveMetrics := &LiveMetricsData{}
	if ts, ok := getFloat64(metricsMap, "timestamp"); ok {
		liveMetrics.Timestamp = int64(ts)
	}

	// Parse CPU usage
	if cpuData, ok := getMap(metricsMap, "cpu_usage"); ok {
		if current, ok := getFloat64(cpuData, "current"); ok {
			if average, ok := getFloat64(cpuData, "average"); ok {
				if peak, ok := getFloat64(cpuData, "peak"); ok {
					if util, ok := getFloat64(cpuData, "utilization"); ok {
						liveMetrics.CPUUsage = CPUUsage{
							Current:     current,
							Average:     average,
							Peak:        peak,
							Utilization: util,
						}
					}
				}
			}
		}
	}

	// Parse Memory usage
	if memData, ok := getMap(metricsMap, "memory_usage"); ok {
		if current, ok := getFloat64(memData, "current"); ok {
			if average, ok := getFloat64(memData, "average"); ok {
				if peak, ok := getFloat64(memData, "peak"); ok {
					if util, ok := getFloat64(memData, "utilization"); ok {
						liveMetrics.MemoryUsage = MemoryUsage{
							Current:     int64(current),
							Average:     int64(average),
							Peak:        int64(peak),
							Utilization: util,
						}
					}
				}
			}
		}
	}

	// Parse Disk usage
	if diskData, ok := getMap(metricsMap, "disk_usage"); ok {
		if readRate, ok := getFloat64(diskData, "read_rate_mbps"); ok {
			if writeRate, ok := getFloat64(diskData, "write_rate_mbps"); ok {
				liveMetrics.DiskUsage = DiskUsage{
					ReadRateMBps:  readRate,
					WriteRateMBps: writeRate,
				}
			}
		}
	}

	// Parse Network usage
	if netData, ok := getMap(metricsMap, "network_usage"); ok {
		if inRate, ok := getFloat64(netData, "in_rate_mbps"); ok {
			if outRate, ok := getFloat64(netData, "out_rate_mbps"); ok {
				liveMetrics.NetworkUsage = NetworkUsage{
					InRateMBps:  inRate,
					OutRateMBps: outRate,
				}
			}
		}
	}

	return liveMetrics, nil
}

func (ac *AnalyticsCollector) getJobResourceTrends(ctx context.Context, jobID string) (*TrendsData, error) {
	url := fmt.Sprintf("%s/slurm/v0.0.42/job/%s/resource_trends?time_window=1h&interval=5m", ac.baseURL, jobID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	resp, err := ac.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get trends data: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	trendsMap, ok := result["trends"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid trends response format")
	}

	trends := &TrendsData{}
	if tw, ok := getString(trendsMap, "time_window"); ok {
		trends.TimeWindow = tw
	}

	// Parse CPU trend
	if cpuTrendData, ok := getMap(trendsMap, "cpu_trend"); ok {
		if dataPoints, ok := getSlice(cpuTrendData, "data_points"); ok {
			for _, pointData := range dataPoints {
				if point, ok := pointData.(map[string]interface{}); ok {
					if ts, ok := getFloat64(point, "timestamp"); ok {
						if val, ok := getFloat64(point, "value"); ok {
							trendPoint := TrendPoint{
								Timestamp: int64(ts),
								Value:     val,
							}
							trends.CPUTrend.DataPoints = append(trends.CPUTrend.DataPoints, trendPoint)
						}
					}
				}
			}
		}
	}

	// Parse Memory trend
	if memTrendData, ok := getMap(trendsMap, "memory_trend"); ok {
		if dataPoints, ok := getSlice(memTrendData, "data_points"); ok {
			for _, pointData := range dataPoints {
				if point, ok := pointData.(map[string]interface{}); ok {
					if ts, ok := getFloat64(point, "timestamp"); ok {
						if val, ok := getFloat64(point, "value"); ok {
							trendPoint := TrendPoint{
								Timestamp: int64(ts),
								Value:     val,
							}
							trends.MemoryTrend.DataPoints = append(trends.MemoryTrend.DataPoints, trendPoint)
						}
					}
				}
			}
		}
	}

	return trends, nil
}

// Analysis and Reporting Functions

// GenerateUtilizationReport generates a comprehensive utilization analysis report
func GenerateUtilizationReport(analytics *JobAnalyticsData) string {
	var report strings.Builder

	report.WriteString("=" + strings.Repeat("=", 60) + "\n")
	report.WriteString(fmt.Sprintf("Job Analytics Report for Job ID: %s\n", analytics.JobID))
	report.WriteString("=" + strings.Repeat("=", 60) + "\n\n")

	if analytics.Utilization != nil {
		report.WriteString("RESOURCE UTILIZATION ANALYSIS\n")
		report.WriteString("-" + strings.Repeat("-", 30) + "\n")

		// CPU Analysis
		cpu := analytics.Utilization.CPU
		report.WriteString("CPU Utilization:\n")
		report.WriteString(fmt.Sprintf("  • Allocated Cores: %d\n", cpu.AllocatedCores))
		report.WriteString(fmt.Sprintf("  • Used Cores: %.2f (%.1f%% utilization)\n", cpu.UsedCores, cpu.UtilizationPercent))
		report.WriteString(fmt.Sprintf("  • Efficiency: %.1f%%\n", cpu.EfficiencyPercent))

		if cpu.UtilizationPercent < 50 {
			report.WriteString("  ⚠️  WARNING: Low CPU utilization detected\n")
		} else if cpu.UtilizationPercent > 90 {
			report.WriteString("  ✅ Excellent CPU utilization\n")
		}
		report.WriteString("\n")

		// Memory Analysis
		mem := analytics.Utilization.Memory
		report.WriteString("Memory Utilization:\n")
		report.WriteString(fmt.Sprintf("  • Allocated: %s\n", formatBytes(mem.AllocatedBytes)))
		report.WriteString(fmt.Sprintf("  • Used: %s (%.1f%% utilization)\n", formatBytes(mem.UsedBytes), mem.UtilizationPercent))
		report.WriteString(fmt.Sprintf("  • Efficiency: %.1f%%\n", mem.EfficiencyPercent))

		if mem.UtilizationPercent < 30 {
			report.WriteString("  ⚠️  WARNING: Low memory utilization detected\n")
		} else if mem.UtilizationPercent > 85 {
			report.WriteString("  ✅ Good memory utilization\n")
		}
		report.WriteString("\n")

		// GPU Analysis
		if analytics.Utilization.GPU.DeviceCount > 0 {
			gpu := analytics.Utilization.GPU
			report.WriteString("GPU Utilization:\n")
			report.WriteString(fmt.Sprintf("  • GPU Devices: %d\n", gpu.DeviceCount))
			report.WriteString(fmt.Sprintf("  • Utilization: %.1f%%\n", gpu.UtilizationPercent))

			if gpu.UtilizationPercent < 60 {
				report.WriteString("  ⚠️  WARNING: Low GPU utilization detected\n")
			} else if gpu.UtilizationPercent > 90 {
				report.WriteString("  ✅ Excellent GPU utilization\n")
			}
			report.WriteString("\n")
		}

		// I/O Analysis
		io := analytics.Utilization.IO
		report.WriteString("I/O Utilization:\n")
		report.WriteString(fmt.Sprintf("  • Read: %s\n", formatBytes(io.ReadBytes)))
		report.WriteString(fmt.Sprintf("  • Write: %s\n", formatBytes(io.WriteBytes)))
		report.WriteString(fmt.Sprintf("  • I/O Utilization: %.1f%%\n", io.UtilizationPercent))
		report.WriteString("\n")
	}

	if analytics.Efficiency != nil {
		report.WriteString("EFFICIENCY ANALYSIS\n")
		report.WriteString("-" + strings.Repeat("-", 18) + "\n")

		eff := analytics.Efficiency
		report.WriteString(fmt.Sprintf("Overall Efficiency Score: %.1f%%\n", eff.OverallEfficiency))
		report.WriteString(fmt.Sprintf("CPU Efficiency: %.1f%%\n", eff.CPUEfficiency))
		report.WriteString(fmt.Sprintf("Memory Efficiency: %.1f%%\n", eff.MemoryEfficiency))
		report.WriteString(fmt.Sprintf("GPU Efficiency: %.1f%%\n", eff.GPUEfficiency))
		report.WriteString("\n")

		// Resource Waste Analysis
		waste := eff.ResourceWaste
		report.WriteString("Resource Waste Analysis:\n")
		report.WriteString(fmt.Sprintf("  • CPU Waste: %.2f core-hours (%.1f%%)\n", waste.CPUCoreHours, waste.CPUPercent))
		report.WriteString(fmt.Sprintf("  • Memory Waste: %.2f GB-hours (%.1f%%)\n", waste.MemoryGBHours, waste.MemoryPercent))
		report.WriteString("\n")

		// Optimization Recommendations
		if len(eff.Recommendations) > 0 {
			report.WriteString("OPTIMIZATION RECOMMENDATIONS\n")
			report.WriteString("-" + strings.Repeat("-", 28) + "\n")
			for i, rec := range eff.Recommendations {
				report.WriteString(fmt.Sprintf("%d. %s %s:\n", i+1, cases.Title(language.English).String(rec.Type), rec.Resource))
				report.WriteString(fmt.Sprintf("   Current: %d → Recommended: %d\n", rec.Current, rec.Recommended))
				report.WriteString(fmt.Sprintf("   Reason: %s\n", rec.Reason))
				report.WriteString(fmt.Sprintf("   Confidence: %.0f%%\n", rec.Confidence*100))
				report.WriteString("\n")
			}
		}
	}

	if analytics.Performance != nil {
		report.WriteString("PERFORMANCE ANALYSIS\n")
		report.WriteString("-" + strings.Repeat("-", 19) + "\n")

		perf := analytics.Performance
		report.WriteString(fmt.Sprintf("Overall Performance: %.1f%%\n", perf.OverallEfficiency))

		// CPU Performance
		cpu := perf.CPUAnalytics
		report.WriteString("CPU Performance:\n")
		report.WriteString(fmt.Sprintf("  • Utilization: %.1f%% (%.2f/%.d cores)\n",
			cpu.UtilizationPercent, cpu.UsedCores, cpu.AllocatedCores))
		report.WriteString(fmt.Sprintf("  • Frequency: %d MHz (max: %d MHz)\n",
			cpu.AverageFrequency, cpu.MaxFrequency))
		report.WriteString("\n")

		// Memory Performance
		mem := perf.MemoryAnalytics
		report.WriteString("Memory Performance:\n")
		report.WriteString(fmt.Sprintf("  • Utilization: %.1f%% (%s/%s)\n",
			mem.UtilizationPercent, formatBytes(mem.UsedBytes), formatBytes(mem.AllocatedBytes)))
		report.WriteString("\n")

		// I/O Performance
		io := perf.IOAnalytics
		report.WriteString("I/O Performance:\n")
		report.WriteString(fmt.Sprintf("  • Read: %s (%d ops, %.1f MB/s)\n",
			formatBytes(io.ReadBytes), io.ReadOperations, io.AverageReadBandwidth))
		report.WriteString(fmt.Sprintf("  • Write: %s (%d ops, %.1f MB/s)\n",
			formatBytes(io.WriteBytes), io.WriteOperations, io.AverageWriteBandwidth))
		report.WriteString("\n")
	}

	if analytics.LiveMetrics != nil {
		report.WriteString("LIVE METRICS SNAPSHOT\n")
		report.WriteString("-" + strings.Repeat("-", 21) + "\n")

		live := analytics.LiveMetrics
		timestamp := time.Unix(live.Timestamp, 0)
		report.WriteString(fmt.Sprintf("Snapshot Time: %s\n", timestamp.Format(time.DateTime)))

		report.WriteString("Current Usage:\n")
		report.WriteString(fmt.Sprintf("  • CPU: %.1f%% (avg: %.1f%%, peak: %.1f%%)\n",
			live.CPUUsage.Current, live.CPUUsage.Average, live.CPUUsage.Peak))
		report.WriteString(fmt.Sprintf("  • Memory: %s (avg: %s, peak: %s)\n",
			formatBytes(live.MemoryUsage.Current), formatBytes(live.MemoryUsage.Average), formatBytes(live.MemoryUsage.Peak)))
		report.WriteString(fmt.Sprintf("  • Disk I/O: %.1f MB/s read, %.1f MB/s write\n",
			live.DiskUsage.ReadRateMBps, live.DiskUsage.WriteRateMBps))
		report.WriteString(fmt.Sprintf("  • Network: %.1f MB/s in, %.1f MB/s out\n",
			live.NetworkUsage.InRateMBps, live.NetworkUsage.OutRateMBps))
		report.WriteString("\n")
	}

	if analytics.Trends != nil && len(analytics.Trends.CPUTrend.DataPoints) > 0 {
		report.WriteString("RESOURCE TRENDS ANALYSIS\n")
		report.WriteString("-" + strings.Repeat("-", 24) + "\n")

		trends := analytics.Trends
		report.WriteString(fmt.Sprintf("Time Window: %s\n", trends.TimeWindow))

		// CPU Trend Analysis
		if len(trends.CPUTrend.DataPoints) > 0 {
			cpuTrend := analyzeTrend(trends.CPUTrend.DataPoints)
			report.WriteString(fmt.Sprintf("CPU Trend: %s\n", cpuTrend))
		}

		// Memory Trend Analysis
		if len(trends.MemoryTrend.DataPoints) > 0 {
			memTrend := analyzeTrend(trends.MemoryTrend.DataPoints)
			report.WriteString(fmt.Sprintf("Memory Trend: %s\n", memTrend))
		}
		report.WriteString("\n")
	}

	// Overall Assessment
	report.WriteString("OVERALL ASSESSMENT\n")
	report.WriteString("-" + strings.Repeat("-", 18) + "\n")

	overallScore := calculateOverallScore(analytics)
	report.WriteString(fmt.Sprintf("Job Efficiency Score: %.1f/100\n", overallScore))

	if overallScore >= 80 {
		report.WriteString("✅ EXCELLENT: Job is running efficiently with optimal resource usage\n")
	} else if overallScore >= 60 {
		report.WriteString("✅ GOOD: Job is running well with minor optimization opportunities\n")
	} else if overallScore >= 40 {
		report.WriteString("⚠️  FAIR: Job has several optimization opportunities\n")
	} else {
		report.WriteString("❌ POOR: Job has significant resource waste and needs optimization\n")
	}

	return report.String()
}

// Helper functions

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func analyzeTrend(dataPoints []TrendPoint) string {
	if len(dataPoints) < 2 {
		return "Insufficient data"
	}

	first := dataPoints[0].Value
	last := dataPoints[len(dataPoints)-1].Value

	change := ((last - first) / first) * 100

	if change > 10 {
		return fmt.Sprintf("Increasing trend (+%.1f%%)", change)
	} else if change < -10 {
		return fmt.Sprintf("Decreasing trend (%.1f%%)", change)
	} else {
		return fmt.Sprintf("Stable (%.1f%% change)", change)
	}
}

func calculateOverallScore(analytics *JobAnalyticsData) float64 {
	scores := []float64{}

	if analytics.Efficiency != nil {
		scores = append(scores, analytics.Efficiency.OverallEfficiency)
	}

	if analytics.Performance != nil {
		scores = append(scores, analytics.Performance.OverallEfficiency)
	}

	if analytics.Utilization != nil {
		// Calculate utilization score
		utilizationScore := (analytics.Utilization.CPU.UtilizationPercent +
			analytics.Utilization.Memory.UtilizationPercent) / 2
		scores = append(scores, utilizationScore)
	}

	if len(scores) == 0 {
		return 0
	}

	total := 0.0
	for _, score := range scores {
		total += score
	}

	return total / float64(len(scores))
}

// CompareJobs compares analytics between multiple jobs
func CompareJobs(jobAnalytics map[string]*JobAnalyticsData) string {
	var report strings.Builder

	report.WriteString("JOB COMPARISON ANALYSIS\n")
	report.WriteString("=" + strings.Repeat("=", 23) + "\n\n")

	// Create sorted list of job IDs for consistent ordering
	jobIDs := make([]string, 0, len(jobAnalytics))
	for jobID := range jobAnalytics {
		jobIDs = append(jobIDs, jobID)
	}
	sort.Strings(jobIDs)

	// Compare efficiency scores
	report.WriteString("Efficiency Comparison:\n")
	report.WriteString("-" + strings.Repeat("-", 21) + "\n")
	for _, jobID := range jobIDs {
		analytics := jobAnalytics[jobID]
		score := calculateOverallScore(analytics)
		report.WriteString(fmt.Sprintf("Job %s: %.1f%%\n", jobID, score))
	}
	report.WriteString("\n")

	// Compare resource utilization
	report.WriteString("Resource Utilization Comparison:\n")
	report.WriteString("-" + strings.Repeat("-", 32) + "\n")
	report.WriteString(fmt.Sprintf("%-8s %-10s %-10s %-10s\n", "Job ID", "CPU %", "Memory %", "GPU %"))
	report.WriteString(strings.Repeat("-", 42) + "\n")

	for _, jobID := range jobIDs {
		analytics := jobAnalytics[jobID]
		if analytics.Utilization != nil {
			u := analytics.Utilization
			report.WriteString(fmt.Sprintf("%-8s %-10.1f %-10.1f %-10.1f\n",
				jobID, u.CPU.UtilizationPercent, u.Memory.UtilizationPercent, u.GPU.UtilizationPercent))
		}
	}
	report.WriteString("\n")

	// Find best and worst performing jobs
	bestJob, worstJob := "", ""
	bestScore, worstScore := 0.0, 100.0

	for _, jobID := range jobIDs {
		score := calculateOverallScore(jobAnalytics[jobID])
		if score > bestScore {
			bestScore = score
			bestJob = jobID
		}
		if score < worstScore {
			worstScore = score
			worstJob = jobID
		}
	}

	report.WriteString("Performance Summary:\n")
	report.WriteString("-" + strings.Repeat("-", 19) + "\n")
	report.WriteString(fmt.Sprintf("Best Performing Job: %s (%.1f%%)\n", bestJob, bestScore))
	report.WriteString(fmt.Sprintf("Worst Performing Job: %s (%.1f%%)\n", worstJob, worstScore))
	report.WriteString(fmt.Sprintf("Performance Gap: %.1f%%\n", bestScore-worstScore))

	return report.String()
}

func main() {
	// Setup mock server for demonstration
	fmt.Println("Starting Job Analytics Example...")
	fmt.Println("Setting up mock SLURM server for demonstration...")

	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")

	if err := run(mockServer); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		mockServer.Close()
		os.Exit(1)
	}
	mockServer.Close()
}

func run(mockServer *mocks.MockSlurmServer) error {
	baseURL := mockServer.URL()
	fmt.Printf("Mock server running at: %s\n\n", baseURL)

	// Create analytics collector
	collector := NewAnalyticsCollector(baseURL)

	// Demo job IDs (these exist in the mock server)
	demoJobs := []string{"1001", "1002"}

	// Collect analytics for all demo jobs
	jobAnalytics := make(map[string]*JobAnalyticsData)

	for _, jobID := range demoJobs {
		fmt.Printf("Collecting analytics for Job %s...\n", jobID)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		analytics := collector.CollectJobAnalytics(ctx, jobID)
		cancel()

		jobAnalytics[jobID] = analytics
		fmt.Printf("✅ Successfully collected analytics for Job %s\n", jobID)
	}

	if len(jobAnalytics) == 0 {
		return errors.New("failed to collect analytics for any jobs")
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("GENERATING ANALYTICS REPORTS")
	fmt.Println(strings.Repeat("=", 80))

	// Generate individual reports for each job
	for _, jobID := range demoJobs {
		if analytics, exists := jobAnalytics[jobID]; exists {
			fmt.Printf("\n")
			report := GenerateUtilizationReport(analytics)
			fmt.Print(report)
		}
	}

	// Generate comparison report if we have multiple jobs
	if len(jobAnalytics) > 1 {
		fmt.Printf("\n" + strings.Repeat("=", 80) + "\n")
		comparisonReport := CompareJobs(jobAnalytics)
		fmt.Print(comparisonReport)
	}

	// Example of how to check job efficiency and provide recommendations
	fmt.Printf("\n" + strings.Repeat("=", 80) + "\n")
	fmt.Println("EFFICIENCY RECOMMENDATIONS")
	fmt.Println(strings.Repeat("=", 80))

	for jobID, analytics := range jobAnalytics {
		if analytics.Efficiency != nil && len(analytics.Efficiency.Recommendations) > 0 {
			fmt.Printf("\nJob %s Optimization Opportunities:\n", jobID)
			for i, rec := range analytics.Efficiency.Recommendations {
				fmt.Printf("%d. %s %s: %d → %d (%.0f%% confidence)\n",
					i+1, cases.Title(language.English).String(rec.Type), rec.Resource,
					rec.Current, rec.Recommended, rec.Confidence*100)
				fmt.Printf("   Reason: %s\n", rec.Reason)
			}
		}
	}

	// Example of exporting data for further analysis
	if len(os.Args) > 1 && os.Args[1] == "--export" {
		fmt.Printf("\n" + strings.Repeat("=", 80) + "\n")
		fmt.Println("EXPORTING ANALYTICS DATA")
		fmt.Println(strings.Repeat("=", 80))

		for jobID, analytics := range jobAnalytics {
			filename := fmt.Sprintf("job_%s_analytics.json", jobID)
			file, err := os.Create(filename)
			if err != nil {
				log.Printf("Failed to create file %s: %v", filename, err)
				continue
			}

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(analytics); err != nil {
				log.Printf("Failed to encode analytics for job %s: %v", jobID, err)
			} else {
				fmt.Printf("✅ Exported analytics for Job %s to %s\n", jobID, filename)
			}
			file.Close()
		}
	}

	fmt.Printf("\n" + strings.Repeat("=", 80) + "\n")
	fmt.Println("Analytics collection and reporting completed successfully!")
	fmt.Println("Use --export flag to save analytics data to JSON files")
	fmt.Println(strings.Repeat("=", 80))
	return nil
}
