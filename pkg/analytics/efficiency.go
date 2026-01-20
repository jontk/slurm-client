// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package analytics

import (
	"math"
	"time"

	"github.com/jontk/slurm-client/interfaces"
)

// EfficiencyCalculator provides methods for calculating job efficiency metrics
type EfficiencyCalculator struct {
	// Default weights for resource efficiency calculation
	defaultWeights ResourceWeights
}

// ResourceWeights defines the relative importance of each resource type in efficiency calculation
type ResourceWeights struct {
	CPU     float64
	Memory  float64
	GPU     float64
	IO      float64
	Network float64
	Energy  float64
}

// DefaultResourceWeights returns standard weights for efficiency calculation
func DefaultResourceWeights() ResourceWeights {
	return ResourceWeights{
		CPU:     0.35, // CPU is most important
		Memory:  0.25, // Memory is second most important
		GPU:     0.20, // GPU is important when present
		IO:      0.10, // I/O has moderate impact
		Network: 0.05, // Network has lower impact
		Energy:  0.05, // Energy efficiency is considered
	}
}

// NewEfficiencyCalculator creates a new efficiency calculator with default weights
func NewEfficiencyCalculator() *EfficiencyCalculator {
	return &EfficiencyCalculator{
		defaultWeights: DefaultResourceWeights(),
	}
}

// NewEfficiencyCalculatorWithWeights creates a calculator with custom weights
func NewEfficiencyCalculatorWithWeights(weights ResourceWeights) *EfficiencyCalculator {
	// Normalize weights to sum to 1.0
	total := weights.CPU + weights.Memory + weights.GPU + weights.IO + weights.Network + weights.Energy
	if total > 0 {
		weights.CPU /= total
		weights.Memory /= total
		weights.GPU /= total
		weights.IO /= total
		weights.Network /= total
		weights.Energy /= total
	}

	return &EfficiencyCalculator{
		defaultWeights: weights,
	}
}

// CalculateOverallEfficiency calculates the overall job efficiency score (0-100%)
func (ec *EfficiencyCalculator) CalculateOverallEfficiency(
	cpuAnalytics *interfaces.CPUAnalytics,
	memoryAnalytics *interfaces.MemoryAnalytics,
	ioAnalytics *interfaces.IOAnalytics,
	gpuUtilization *interfaces.GPUUtilization,
	networkUtilization *interfaces.NetworkUtilization,
	energyUsage *interfaces.EnergyUsage,
) float64 {
	weights := ec.defaultWeights
	totalWeight := 0.0
	weightedSum := 0.0

	// CPU efficiency
	if cpuAnalytics != nil {
		cpuEff := ec.CalculateCPUEfficiency(cpuAnalytics)
		weightedSum += cpuEff * weights.CPU
		totalWeight += weights.CPU
	}

	// Memory efficiency
	if memoryAnalytics != nil {
		memEff := ec.CalculateMemoryEfficiency(memoryAnalytics)
		weightedSum += memEff * weights.Memory
		totalWeight += weights.Memory
	}

	// GPU efficiency (only if GPUs are allocated)
	if gpuUtilization != nil && gpuUtilization.DeviceCount > 0 {
		gpuEff := ec.CalculateGPUEfficiency(gpuUtilization)
		weightedSum += gpuEff * weights.GPU
		totalWeight += weights.GPU
	}

	// I/O efficiency
	if ioAnalytics != nil {
		ioEff := ec.CalculateIOEfficiency(ioAnalytics)
		weightedSum += ioEff * weights.IO
		totalWeight += weights.IO
	}

	// Network efficiency
	if networkUtilization != nil {
		netEff := ec.CalculateNetworkEfficiency(networkUtilization)
		weightedSum += netEff * weights.Network
		totalWeight += weights.Network
	}

	// Energy efficiency
	if energyUsage != nil {
		energyEff := ec.CalculateEnergyEfficiency(energyUsage)
		weightedSum += energyEff * weights.Energy
		totalWeight += weights.Energy
	}

	// Calculate final score
	if totalWeight > 0 {
		return (weightedSum / totalWeight) * 100.0
	}

	return 0.0
}

// CalculateCPUEfficiency calculates CPU efficiency percentage
func (ec *EfficiencyCalculator) CalculateCPUEfficiency(analytics *interfaces.CPUAnalytics) float64 {
	if analytics.AllocatedCores == 0 {
		return 0.0
	}

	// Base efficiency is utilization
	efficiency := analytics.UtilizationPercent

	// Penalize oversubscription
	if analytics.Oversubscribed {
		efficiency *= 0.8 // 20% penalty for oversubscription
	}

	// Penalize thermal throttling
	if analytics.ThermalThrottleEvents > 0 {
		throttlePenalty := math.Min(float64(analytics.ThermalThrottleEvents)/1000.0*0.1, 0.2)
		efficiency *= (1.0 - throttlePenalty)
	}

	// Penalize frequency scaling
	if analytics.MaxFrequency > 0 && analytics.AverageFrequency < analytics.MaxFrequency*0.9 {
		scalingPenalty := (analytics.MaxFrequency - analytics.AverageFrequency) / analytics.MaxFrequency * 0.1
		efficiency *= (1.0 - scalingPenalty)
	}

	// Consider core efficiency distribution
	if len(analytics.CoreMetrics) > 0 {
		variance := ec.calculateCoreUtilizationVariance(analytics.CoreMetrics)
		if variance > 20.0 { // High variance indicates imbalanced load
			efficiency *= 0.95 // 5% penalty for imbalanced load
		}
	}

	return math.Min(efficiency, 100.0)
}

// CalculateMemoryEfficiency calculates memory efficiency percentage
func (ec *EfficiencyCalculator) CalculateMemoryEfficiency(analytics *interfaces.MemoryAnalytics) float64 {
	if analytics.AllocatedBytes == 0 {
		return 0.0
	}

	// Base efficiency is utilization
	efficiency := analytics.UtilizationPercent

	// Penalize swap usage (using PageSwaps as indicator)
	if analytics.PageSwaps > 0 {
		swapRatio := float64(analytics.PageSwaps) / 10000.0 // Normalize swap activity
		swapPenalty := math.Min(swapRatio*0.3, 0.3)         // Up to 30% penalty for swap
		efficiency *= (1.0 - swapPenalty)
	}

	// Penalize major page faults
	if analytics.MajorPageFaults > 100 {
		faultPenalty := math.Min(float64(analytics.MajorPageFaults)/10000.0*0.1, 0.15)
		efficiency *= (1.0 - faultPenalty)
	}

	// Reward good memory locality (NUMA efficiency)
	if len(analytics.NUMANodes) > 1 {
		avgLocalAccess := 0.0
		for _, numa := range analytics.NUMANodes {
			avgLocalAccess += numa.LocalAccesses
		}
		avgLocalAccess /= float64(len(analytics.NUMANodes))

		if avgLocalAccess > 90.0 {
			efficiency *= 1.05 // 5% bonus for excellent locality
		} else if avgLocalAccess < 70.0 {
			efficiency *= 0.95 // 5% penalty for poor locality
		}
	}

	return math.Min(efficiency, 100.0)
}

// CalculateGPUEfficiency calculates GPU efficiency percentage
func (ec *EfficiencyCalculator) CalculateGPUEfficiency(utilization *interfaces.GPUUtilization) float64 {
	if utilization.DeviceCount == 0 {
		return 0.0
	}

	// Use overall utilization
	efficiency := 0.0
	if utilization.OverallUtilization != nil {
		efficiency = utilization.OverallUtilization.Percentage
	}

	// Penalize underutilized individual GPUs
	underutilizedCount := 0
	for _, gpu := range utilization.Devices {
		if gpu.Utilization != nil && gpu.Utilization.Percentage < 50.0 {
			underutilizedCount++
		}
	}

	if underutilizedCount > 0 {
		penalty := float64(underutilizedCount) / float64(utilization.DeviceCount) * 0.1
		efficiency *= (1.0 - penalty)
	}

	// Consider memory efficiency on GPUs
	avgMemoryUtilization := 0.0
	for _, gpu := range utilization.Devices {
		if gpu.MemoryUtilization != nil {
			avgMemoryUtilization += gpu.MemoryUtilization.Percentage
		}
	}
	if len(utilization.Devices) > 0 {
		avgMemoryUtilization /= float64(len(utilization.Devices))
		// Blend GPU compute and memory utilization
		efficiency = efficiency*0.7 + avgMemoryUtilization*0.3
	}

	return math.Min(efficiency, 100.0)
}

// CalculateIOEfficiency calculates I/O efficiency percentage
func (ec *EfficiencyCalculator) CalculateIOEfficiency(analytics *interfaces.IOAnalytics) float64 {
	// Base efficiency on bandwidth utilization vs theoretical max
	// Assume theoretical max is 1000 MB/s for reads and 500 MB/s for writes
	theoreticalReadBW := 1000.0
	theoreticalWriteBW := 500.0

	readEfficiency := math.Min(analytics.AverageReadBandwidth/theoreticalReadBW*100.0, 100.0)
	writeEfficiency := math.Min(analytics.AverageWriteBandwidth/theoreticalWriteBW*100.0, 100.0)

	// Weight reads and writes based on operation count
	totalOps := analytics.ReadOperations + analytics.WriteOperations
	if totalOps == 0 {
		return 0.0
	}

	readWeight := float64(analytics.ReadOperations) / float64(totalOps)
	writeWeight := float64(analytics.WriteOperations) / float64(totalOps)

	efficiency := readEfficiency*readWeight + writeEfficiency*writeWeight

	// Penalize high latency
	avgLatency := (analytics.AverageReadLatency*readWeight + analytics.AverageWriteLatency*writeWeight)
	if avgLatency > 20.0 {
		latencyPenalty := math.Min((avgLatency-20.0)/100.0*0.15, 0.15)
		efficiency *= (1.0 - latencyPenalty)
	}

	// Use the utilization percentage if available
	if analytics.UtilizationPercent > 0 {
		// Blend bandwidth efficiency with overall utilization
		efficiency = efficiency*0.7 + analytics.UtilizationPercent*0.3
	}

	return math.Min(efficiency, 100.0)
}

// CalculateNetworkEfficiency calculates network efficiency percentage
func (ec *EfficiencyCalculator) CalculateNetworkEfficiency(utilization *interfaces.NetworkUtilization) float64 {
	if utilization.TotalBandwidth == nil || utilization.TotalBandwidth.UsedMax == 0 {
		return 0.0
	}

	efficiency := utilization.TotalBandwidth.Efficiency

	// Consider individual interface efficiency
	underutilizedInterfaces := 0
	for _, iface := range utilization.Interfaces {
		if iface.Utilization < 30.0 {
			underutilizedInterfaces++
		}
	}

	if len(utilization.Interfaces) > 0 && underutilizedInterfaces > 0 {
		penalty := float64(underutilizedInterfaces) / float64(len(utilization.Interfaces)) * 0.1
		efficiency *= (1.0 - penalty)
	}

	// Consider protocol efficiency if available
	if utilization.Metadata != nil {
		if retransmits, ok := utilization.Metadata["tcp_retransmits"].(float64); ok && retransmits > 0.01 {
			// Penalize high retransmission rate
			efficiency *= 0.95
		}
	}

	return math.Min(efficiency, 100.0)
}

// CalculateEnergyEfficiency calculates energy efficiency percentage
func (ec *EfficiencyCalculator) CalculateEnergyEfficiency(usage *interfaces.EnergyUsage) float64 {
	if usage.TotalEnergyJoules == 0 {
		return 100.0 // No energy usage data means we can't penalize
	}

	// Calculate performance per watt metric
	avgPower := usage.AveragePowerWatts
	if avgPower == 0 {
		return 100.0
	}

	// Base efficiency on power usage patterns
	efficiency := 100.0
	// Note: TDPWatts not available in current interface, using peak power as reference
	if usage.PeakPowerWatts > 0 {
		powerRatio := avgPower / usage.PeakPowerWatts
		if powerRatio > 0.9 {
			// Running close to peak, likely inefficient
			efficiency = 70.0
		} else if powerRatio > 0.7 {
			efficiency = 85.0
		} else if powerRatio > 0.5 {
			efficiency = 95.0
		}
	}

	// Consider energy breakdown using individual energy fields
	totalBreakdown := usage.CPUEnergyJoules + usage.MemoryEnergyJoules + usage.GPUEnergyJoules
	if totalBreakdown > 0 {
		coreRatio := usage.CPUEnergyJoules / totalBreakdown
		if coreRatio < 0.6 {
			// Too much energy in non-compute components
			efficiency *= 0.9
		}
	}

	return efficiency
}

// calculateCoreUtilizationVariance calculates the variance in core utilization
func (ec *EfficiencyCalculator) calculateCoreUtilizationVariance(cores []interfaces.CPUCoreMetric) float64 {
	if len(cores) == 0 {
		return 0.0
	}

	// Calculate mean
	sum := 0.0
	for _, core := range cores {
		sum += core.Utilization
	}
	mean := sum / float64(len(cores))

	// Calculate variance
	variance := 0.0
	for _, core := range cores {
		diff := core.Utilization - mean
		variance += diff * diff
	}
	variance /= float64(len(cores))

	return math.Sqrt(variance) // Return standard deviation
}

// CalculateResourceWaste calculates wasted resources for each resource type
func (ec *EfficiencyCalculator) CalculateResourceWaste(
	job *interfaces.Job,
	analytics *interfaces.JobComprehensiveAnalytics,
	runtime time.Duration,
) map[string]float64 {
	waste := make(map[string]float64)

	// CPU waste (in core-hours)
	if analytics.CPUAnalytics != nil {
		idleCores := float64(analytics.CPUAnalytics.AllocatedCores) - analytics.CPUAnalytics.UsedCores
		waste["cpu_core_hours"] = idleCores * runtime.Hours()
		waste["cpu_percent"] = (idleCores / float64(analytics.CPUAnalytics.AllocatedCores)) * 100.0
	}

	// Memory waste (in GB-hours)
	if analytics.MemoryAnalytics != nil {
		unusedBytes := analytics.MemoryAnalytics.AllocatedBytes - analytics.MemoryAnalytics.UsedBytes
		waste["memory_gb_hours"] = float64(unusedBytes) / (1024 * 1024 * 1024) * runtime.Hours()
		waste["memory_percent"] = (float64(unusedBytes) / float64(analytics.MemoryAnalytics.AllocatedBytes)) * 100.0
	}

	// GPU waste (in GPU-hours) - using Metadata to check for GPU allocation
	if gpuCount, ok := job.Metadata["gpus"].(int); ok && gpuCount > 0 {
		// Estimate GPU efficiency from overall efficiency if available
		gpuEfficiency := analytics.OverallEfficiency
		if gpuEfficiency == 0 {
			gpuEfficiency = 50.0 // Default assumption
		}
		waste["gpu_hours"] = float64(gpuCount) * (100.0 - gpuEfficiency) / 100.0 * runtime.Hours()
		waste["gpu_percent"] = 100.0 - gpuEfficiency
	}

	// I/O waste (estimated from low utilization)
	if analytics.IOAnalytics != nil && analytics.IOAnalytics.UtilizationPercent < 50.0 {
		wastePercent := 100.0 - analytics.IOAnalytics.UtilizationPercent
		waste["io_underutilized_hours"] = wastePercent / 100.0 * runtime.Hours()
		waste["io_waste_percent"] = wastePercent
	}

	return waste
}

// GenerateOptimizationRecommendations generates specific optimization recommendations
func (ec *EfficiencyCalculator) GenerateOptimizationRecommendations(
	job *interfaces.Job,
	analytics *interfaces.JobComprehensiveAnalytics,
) []OptimizationRecommendation {
	recommendations := []OptimizationRecommendation{}

	// CPU recommendations
	if analytics.CPUAnalytics != nil {
		if analytics.CPUAnalytics.UtilizationPercent < 50.0 {
			recommendations = append(recommendations, OptimizationRecommendation{
				Resource:    "CPU",
				Type:        "reduction",
				Current:     analytics.CPUAnalytics.AllocatedCores,
				Recommended: int(math.Ceil(analytics.CPUAnalytics.UsedCores * 1.2)), // 20% buffer
				Reason:      "Low CPU utilization detected",
				Impact:      "Reduce resource waste and improve cluster efficiency",
				Confidence:  0.8,
			})
		}

		if analytics.CPUAnalytics.ThermalThrottleEvents > 100 {
			recommendations = append(recommendations, OptimizationRecommendation{
				Resource:    "CPU",
				Type:        "configuration",
				Current:     analytics.CPUAnalytics.AverageFrequency,
				Recommended: analytics.CPUAnalytics.AverageFrequency * 0.9,
				Reason:      "Thermal throttling detected",
				Impact:      "Improve performance stability",
				Confidence:  0.7,
			})
		}
	}

	// Memory recommendations
	if analytics.MemoryAnalytics != nil {
		if analytics.MemoryAnalytics.UtilizationPercent < 50.0 {
			recommendedBytes := int64(float64(analytics.MemoryAnalytics.UsedBytes) * 1.25) // 25% buffer
			recommendations = append(recommendations, OptimizationRecommendation{
				Resource:    "Memory",
				Type:        "reduction",
				Current:     analytics.MemoryAnalytics.AllocatedBytes / (1024 * 1024 * 1024), // Convert to GB
				Recommended: recommendedBytes / (1024 * 1024 * 1024),
				Reason:      "Low memory utilization detected",
				Impact:      "Free up memory for other jobs",
				Confidence:  0.85,
			})
		}

		// Check for potential memory leak by comparing virtual vs resident memory
		if analytics.MemoryAnalytics.VirtualMemorySize > analytics.MemoryAnalytics.ResidentSetSize*2 {
			recommendations = append(recommendations, OptimizationRecommendation{
				Resource:   "Memory",
				Type:       "code_review",
				Current:    float64(analytics.MemoryAnalytics.VirtualMemorySize) / float64(analytics.MemoryAnalytics.ResidentSetSize),
				Reason:     "High virtual memory usage suggests potential inefficiency",
				Impact:     "Review memory allocation patterns",
				Confidence: 0.6,
			})
		}

		// Check if memory usage is close to allocation (potential pressure)
		if analytics.MemoryAnalytics.UtilizationPercent > 90.0 {
			recommendations = append(recommendations, OptimizationRecommendation{
				Resource:    "Memory",
				Type:        "increase",
				Current:     analytics.MemoryAnalytics.AllocatedBytes / (1024 * 1024 * 1024),
				Recommended: int64(float64(analytics.MemoryAnalytics.AllocatedBytes)*1.2) / (1024 * 1024 * 1024),
				Reason:      "High memory pressure detected",
				Impact:      "Prevent potential memory constraints",
				Confidence:  0.8,
			})
		}
	}

	// GPU recommendations - use overall efficiency as proxy for GPU efficiency
	if gpuCount, ok := job.Metadata["gpus"].(int); ok && gpuCount > 0 && analytics.OverallEfficiency < 50.0 {
		recommendations = append(recommendations, OptimizationRecommendation{
			Resource:    "GPU",
			Type:        "reduction",
			Current:     gpuCount,
			Recommended: int(math.Ceil(float64(gpuCount) * analytics.OverallEfficiency / 100.0)),
			Reason:      "Low overall efficiency may indicate GPU underutilization",
			Impact:      "Free up expensive GPU resources",
			Confidence:  0.6,
		})
	}

	// I/O recommendations
	if analytics.IOAnalytics != nil {
		if analytics.IOAnalytics.UtilizationPercent < 30.0 {
			recommendations = append(recommendations, OptimizationRecommendation{
				Resource:   "IO",
				Type:       "optimization",
				Current:    analytics.IOAnalytics.UtilizationPercent,
				Reason:     "Low I/O utilization detected",
				Impact:     "Optimize storage usage",
				Confidence: 0.7,
				Details: map[string]interface{}{
					"suggestion":  "Consider adjusting I/O patterns or using different storage",
					"utilization": analytics.IOAnalytics.UtilizationPercent,
				},
			})
		}

		// Check for inefficient small I/O operations
		if analytics.IOAnalytics.ReadOperations+analytics.IOAnalytics.WriteOperations > 100000 {
			avgIOSize := float64(analytics.IOAnalytics.ReadBytes+analytics.IOAnalytics.WriteBytes) /
				float64(analytics.IOAnalytics.ReadOperations+analytics.IOAnalytics.WriteOperations)
			if avgIOSize < 4096 { // Less than 4KB average
				recommendations = append(recommendations, OptimizationRecommendation{
					Resource:   "IO",
					Type:       "pattern",
					Current:    avgIOSize,
					Reason:     "Many small I/O operations detected",
					Impact:     "Reduce I/O overhead",
					Confidence: 0.8,
					Details: map[string]interface{}{
						"suggestion":       "Batch I/O operations or use buffering",
						"avg_io_size":      avgIOSize,
						"total_operations": analytics.IOAnalytics.ReadOperations + analytics.IOAnalytics.WriteOperations,
					},
				})
			}
		}
	}

	// Overall efficiency recommendation
	if analytics.OverallEfficiency < 60.0 {
		recommendations = append(recommendations, OptimizationRecommendation{
			Resource:   "Overall",
			Type:       "review",
			Current:    analytics.OverallEfficiency,
			Reason:     "Low overall job efficiency",
			Impact:     "Significant resource waste",
			Confidence: 0.85,
			Details: map[string]interface{}{
				"cpu_efficiency":    analytics.CPUAnalytics.EfficiencyPercent,
				"memory_efficiency": analytics.MemoryAnalytics.EfficiencyPercent,
				"io_efficiency":     analytics.IOAnalytics.EfficiencyPercent,
			},
		})
	}

	return recommendations
}

// OptimizationRecommendation represents a specific optimization suggestion
type OptimizationRecommendation struct {
	Resource    string                 // Resource type (CPU, Memory, GPU, IO, Network)
	Type        string                 // Type of recommendation (reduction, increase, configuration, pattern)
	Current     interface{}            // Current value/configuration
	Recommended interface{}            // Recommended value/configuration
	Reason      string                 // Why this recommendation is made
	Impact      string                 // Expected impact of the change
	Confidence  float64                // Confidence level (0-1)
	Details     map[string]interface{} // Additional details
}
