// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package analytics

import (
	"testing"
	"time"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultResourceWeights(t *testing.T) {
	weights := DefaultResourceWeights()
	
	// Verify weights sum to 1.0
	total := weights.CPU + weights.Memory + weights.GPU + weights.IO + weights.Network + weights.Energy
	assert.InDelta(t, 1.0, total, 0.001)
	
	// Verify relative importance
	assert.Greater(t, weights.CPU, weights.Memory)
	assert.Greater(t, weights.Memory, weights.GPU)
	assert.Greater(t, weights.GPU, weights.IO)
}

func TestNewEfficiencyCalculatorWithWeights(t *testing.T) {
	// Test normalization
	weights := ResourceWeights{
		CPU:    10.0,
		Memory: 5.0,
		GPU:    5.0,
		IO:     0.0,
		Network: 0.0,
		Energy: 0.0,
	}
	
	calc := NewEfficiencyCalculatorWithWeights(weights)
	
	// Weights should be normalized
	assert.InDelta(t, 0.5, calc.defaultWeights.CPU, 0.001)
	assert.InDelta(t, 0.25, calc.defaultWeights.Memory, 0.001)
	assert.InDelta(t, 0.25, calc.defaultWeights.GPU, 0.001)
}

func TestCalculateCPUEfficiency(t *testing.T) {
	calc := NewEfficiencyCalculator()
	
	tests := []struct {
		name      string
		analytics *interfaces.CPUAnalytics
		expected  float64
	}{
		{
			name: "high utilization",
			analytics: &interfaces.CPUAnalytics{
				AllocatedCores:     16,
				UsedCores:         14.4,
				UtilizationPercent: 90.0,
				Oversubscribed:    false,
			},
			expected: 90.0,
		},
		{
			name: "with oversubscription",
			analytics: &interfaces.CPUAnalytics{
				AllocatedCores:     16,
				UsedCores:         14.4,
				UtilizationPercent: 90.0,
				Oversubscribed:    true,
			},
			expected: 72.0, // 90 * 0.8
		},
		{
			name: "with thermal throttling",
			analytics: &interfaces.CPUAnalytics{
				AllocatedCores:        16,
				UsedCores:            12.8,
				UtilizationPercent:   80.0,
				ThermalThrottleEvents: 500,
			},
			expected: 76.0, // 80 * 0.95 (5% penalty)
		},
		{
			name: "with frequency scaling",
			analytics: &interfaces.CPUAnalytics{
				AllocatedCores:     16,
				UsedCores:         12.8,
				UtilizationPercent: 80.0,
				MaxFrequency:      3.6,
				AverageFrequency:  3.0,
			},
			expected: 78.67, // 80 * 0.9833 (1.67% penalty)
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			efficiency := calc.CalculateCPUEfficiency(tt.analytics)
			assert.InDelta(t, tt.expected, efficiency, 0.1)
		})
	}
}

func TestCalculateMemoryEfficiency(t *testing.T) {
	calc := NewEfficiencyCalculator()
	
	tests := []struct {
		name      string
		analytics *interfaces.MemoryAnalytics
		expected  float64
	}{
		{
			name: "high utilization no issues",
			analytics: &interfaces.MemoryAnalytics{
				AllocatedBytes:     64 * 1024 * 1024 * 1024,
				UsedBytes:         56 * 1024 * 1024 * 1024,
				UtilizationPercent: 87.5,
				SwapBytes:         0,
				MajorPageFaults:   10,
			},
			expected: 87.5,
		},
		{
			name: "with swap usage",
			analytics: &interfaces.MemoryAnalytics{
				AllocatedBytes:     64 * 1024 * 1024 * 1024,
				UsedBytes:         56 * 1024 * 1024 * 1024,
				UtilizationPercent: 87.5,
				SwapBytes:         4 * 1024 * 1024 * 1024, // 6.25% of allocated
				MajorPageFaults:   10,
			},
			expected: 85.88, // 87.5 * 0.98125
		},
		{
			name: "with memory leak",
			analytics: &interfaces.MemoryAnalytics{
				AllocatedBytes:     64 * 1024 * 1024 * 1024,
				UsedBytes:         48 * 1024 * 1024 * 1024,
				UtilizationPercent: 75.0,
				MemoryLeakDetected: true,
			},
			expected: 37.5, // 75 * 0.5
		},
		{
			name: "with good NUMA locality",
			analytics: &interfaces.MemoryAnalytics{
				AllocatedBytes:     64 * 1024 * 1024 * 1024,
				UsedBytes:         48 * 1024 * 1024 * 1024,
				UtilizationPercent: 75.0,
				NUMAMetrics: []interfaces.NUMANodeMetrics{
					{LocalAccessPercent: 95.0},
					{LocalAccessPercent: 93.0},
				},
			},
			expected: 78.75, // 75 * 1.05
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			efficiency := calc.CalculateMemoryEfficiency(tt.analytics)
			assert.InDelta(t, tt.expected, efficiency, 0.1)
		})
	}
}

func TestCalculateGPUEfficiency(t *testing.T) {
	calc := NewEfficiencyCalculator()
	
	tests := []struct {
		name        string
		utilization *interfaces.GPUUtilization
		expected    float64
	}{
		{
			name: "single GPU high utilization",
			utilization: &interfaces.GPUUtilization{
				TotalGPUs: 1,
				AverageUtilization: &interfaces.ResourceUtilization{
					Percentage: 85.0,
				},
				GPUs: []interfaces.GPUDeviceInfo{
					{
						Utilization: &interfaces.ResourceUtilization{Percentage: 85.0},
						MemoryUtilization: &interfaces.ResourceUtilization{Percentage: 90.0},
					},
				},
			},
			expected: 86.5, // 85*0.7 + 90*0.3
		},
		{
			name: "multiple GPUs with imbalance",
			utilization: &interfaces.GPUUtilization{
				TotalGPUs: 4,
				AverageUtilization: &interfaces.ResourceUtilization{
					Percentage: 70.0,
				},
				GPUs: []interfaces.GPUDeviceInfo{
					{
						Utilization: &interfaces.ResourceUtilization{Percentage: 90.0},
						MemoryUtilization: &interfaces.ResourceUtilization{Percentage: 85.0},
					},
					{
						Utilization: &interfaces.ResourceUtilization{Percentage: 80.0},
						MemoryUtilization: &interfaces.ResourceUtilization{Percentage: 75.0},
					},
					{
						Utilization: &interfaces.ResourceUtilization{Percentage: 60.0},
						MemoryUtilization: &interfaces.ResourceUtilization{Percentage: 65.0},
					},
					{
						Utilization: &interfaces.ResourceUtilization{Percentage: 40.0}, // Underutilized
						MemoryUtilization: &interfaces.ResourceUtilization{Percentage: 45.0},
					},
				},
			},
			expected: 65.275, // (70*0.7 + 67.5*0.3) * 0.975 (2.5% penalty for 1 underutilized)
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			efficiency := calc.CalculateGPUEfficiency(tt.utilization)
			assert.InDelta(t, tt.expected, efficiency, 0.1)
		})
	}
}

func TestCalculateIOEfficiency(t *testing.T) {
	calc := NewEfficiencyCalculator()
	
	tests := []struct {
		name      string
		analytics *interfaces.IOAnalytics
		expected  float64
	}{
		{
			name: "balanced read/write",
			analytics: &interfaces.IOAnalytics{
				ReadBandwidthMBps:  500.0,
				WriteBandwidthMBps: 250.0,
				ReadOperations:     10000,
				WriteOperations:    10000,
				IOWaitPercent:      5.0,
				ReadLatencyMs:      10.0,
				WriteLatencyMs:     15.0,
			},
			expected: 50.0, // (50 + 50) / 2
		},
		{
			name: "high IO wait",
			analytics: &interfaces.IOAnalytics{
				ReadBandwidthMBps:  500.0,
				WriteBandwidthMBps: 250.0,
				ReadOperations:     10000,
				WriteOperations:    10000,
				IOWaitPercent:      25.0, // High wait
				ReadLatencyMs:      10.0,
				WriteLatencyMs:     15.0,
			},
			expected: 47.0, // 50 * 0.94 (6% penalty)
		},
		{
			name: "with device metrics",
			analytics: &interfaces.IOAnalytics{
				ReadBandwidthMBps:  300.0,
				WriteBandwidthMBps: 150.0,
				ReadOperations:     10000,
				WriteOperations:    10000,
				IOWaitPercent:      5.0,
				DeviceMetrics: []interfaces.IODeviceMetrics{
					{DeviceName: "/dev/sda", UtilizationPercent: 60.0},
					{DeviceName: "/dev/sdb", UtilizationPercent: 40.0},
				},
			},
			expected: 35.0, // 30*0.7 + 50*0.3
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			efficiency := calc.CalculateIOEfficiency(tt.analytics)
			assert.InDelta(t, tt.expected, efficiency, 1.0)
		})
	}
}

func TestCalculateOverallEfficiency(t *testing.T) {
	calc := NewEfficiencyCalculator()
	
	// Create sample analytics data
	cpuAnalytics := &interfaces.CPUAnalytics{
		AllocatedCores:     16,
		UsedCores:         12.8,
		UtilizationPercent: 80.0,
	}
	
	memoryAnalytics := &interfaces.MemoryAnalytics{
		AllocatedBytes:     64 * 1024 * 1024 * 1024,
		UsedBytes:         48 * 1024 * 1024 * 1024,
		UtilizationPercent: 75.0,
	}
	
	ioAnalytics := &interfaces.IOAnalytics{
		ReadBandwidthMBps:  400.0,
		WriteBandwidthMBps: 200.0,
		ReadOperations:     10000,
		WriteOperations:    10000,
	}
	
	gpuUtilization := &interfaces.GPUUtilization{
		TotalGPUs: 2,
		AverageUtilization: &interfaces.ResourceUtilization{
			Percentage: 85.0,
		},
	}
	
	overall := calc.CalculateOverallEfficiency(
		cpuAnalytics,
		memoryAnalytics,
		ioAnalytics,
		gpuUtilization,
		nil, // No network data
		nil, // No energy data
	)
	
	// With default weights: CPU=0.35, Memory=0.25, GPU=0.20, IO=0.10
	// Expected: 80*0.35 + 75*0.25 + 85*0.20 + 40*0.10 = 67.75
	// Normalized by totalWeight = 0.9 (no network/energy)
	// Result: 67.75 / 0.9 * 100 = 75.28
	assert.InDelta(t, 75.28, overall, 1.0)
}

func TestCalculateResourceWaste(t *testing.T) {
	calc := NewEfficiencyCalculator()
	
	job := &interfaces.Job{
		ID:   "test-job",
		CPUs: 16,
		Memory: 64 * 1024 * 1024 * 1024,
		GPUs: 2,
	}
	
	analytics := &interfaces.JobComprehensiveAnalytics{
		CPUAnalytics: &interfaces.CPUAnalytics{
			AllocatedCores: 16,
			UsedCores:     10.0,
		},
		MemoryAnalytics: &interfaces.MemoryAnalytics{
			AllocatedBytes: 64 * 1024 * 1024 * 1024,
			UsedBytes:     40 * 1024 * 1024 * 1024,
		},
		IOAnalytics: &interfaces.IOAnalytics{
			IOWaitPercent: 25.0,
		},
		EfficiencyMetrics: &interfaces.JobEfficiencyMetrics{
			GPUEfficiency: 60.0,
		},
	}
	
	runtime := 2 * time.Hour
	
	waste := calc.CalculateResourceWaste(job, analytics, runtime)
	
	// CPU waste: 6 cores * 2 hours = 12 core-hours
	assert.InDelta(t, 12.0, waste["cpu_core_hours"], 0.1)
	assert.InDelta(t, 37.5, waste["cpu_percent"], 0.1) // 6/16 * 100
	
	// Memory waste: 24GB * 2 hours = 48 GB-hours
	assert.InDelta(t, 48.0, waste["memory_gb_hours"], 0.1)
	assert.InDelta(t, 37.5, waste["memory_percent"], 0.1) // 24/64 * 100
	
	// GPU waste: 2 GPUs * 40% inefficiency * 2 hours = 1.6 GPU-hours
	assert.InDelta(t, 1.6, waste["gpu_hours"], 0.1)
	assert.InDelta(t, 40.0, waste["gpu_percent"], 0.1)
	
	// I/O wait: 25% * 2 hours = 0.5 hours
	assert.InDelta(t, 0.5, waste["io_wait_hours"], 0.1)
	assert.InDelta(t, 25.0, waste["io_wait_percent"], 0.1)
}

func TestGenerateOptimizationRecommendations(t *testing.T) {
	calc := NewEfficiencyCalculator()
	
	job := &interfaces.Job{
		ID:   "test-job",
		CPUs: 16,
		Memory: 64 * 1024 * 1024 * 1024,
		GPUs: 4,
	}
	
	analytics := &interfaces.JobComprehensiveAnalytics{
		CPUAnalytics: &interfaces.CPUAnalytics{
			AllocatedCores:        16,
			UsedCores:            6.4,
			UtilizationPercent:   40.0, // Low utilization
			ThermalThrottleEvents: 200, // Some throttling
		},
		MemoryAnalytics: &interfaces.MemoryAnalytics{
			AllocatedBytes:     64 * 1024 * 1024 * 1024,
			UsedBytes:         20 * 1024 * 1024 * 1024,
			UtilizationPercent: 31.25, // Very low utilization
			SwapBytes:         8 * 1024 * 1024 * 1024, // Significant swap
			MemoryLeakDetected: true,
			LeakRatePerHour:   100.0,
		},
		IOAnalytics: &interfaces.IOAnalytics{
			IOWaitPercent:      35.0, // High I/O wait
			ReadBytes:         1000000,
			WriteBytes:        1000000,
			ReadOperations:    1000000,
			WriteOperations:   1000000,
		},
		EfficiencyMetrics: &interfaces.JobEfficiencyMetrics{
			OverallEfficiencyScore: 45.0,
			CPUEfficiency:         40.0,
			MemoryEfficiency:      31.25,
			GPUEfficiency:         30.0,
			IOEfficiency:          50.0,
		},
	}
	
	recommendations := calc.GenerateOptimizationRecommendations(job, analytics)
	
	// Should have recommendations for:
	// 1. CPU reduction (low utilization)
	// 2. CPU configuration (thermal throttling)
	// 3. Memory reduction (low utilization)
	// 4. Memory leak fix
	// 5. Memory increase (swap usage)
	// 6. GPU reduction (low utilization)
	// 7. I/O optimization (high wait)
	// 8. Small I/O pattern
	// 9. Overall efficiency review
	
	require.True(t, len(recommendations) >= 8)
	
	// Check for specific recommendations
	foundCPUReduction := false
	foundMemoryLeak := false
	foundGPUReduction := false
	foundIOOptimization := false
	foundOverallReview := false
	
	for _, rec := range recommendations {
		switch rec.Resource {
		case "CPU":
			if rec.Type == "reduction" {
				foundCPUReduction = true
				assert.Equal(t, 16, rec.Current)
				assert.Equal(t, 8, rec.Recommended) // 6.4 * 1.2 rounded up
			}
		case "Memory":
			if rec.Type == "code_fix" {
				foundMemoryLeak = true
				assert.Equal(t, 100.0, rec.Current)
			}
		case "GPU":
			if rec.Type == "reduction" {
				foundGPUReduction = true
				assert.Equal(t, 4, rec.Current)
				assert.Equal(t, 2, rec.Recommended) // 4 * 30% rounded up
			}
		case "IO":
			if rec.Type == "optimization" {
				foundIOOptimization = true
				assert.Equal(t, 35.0, rec.Current)
			}
		case "Overall":
			if rec.Type == "review" {
				foundOverallReview = true
				assert.Equal(t, 45.0, rec.Current)
			}
		}
	}
	
	assert.True(t, foundCPUReduction, "Should recommend CPU reduction")
	assert.True(t, foundMemoryLeak, "Should recommend memory leak fix")
	assert.True(t, foundGPUReduction, "Should recommend GPU reduction")
	assert.True(t, foundIOOptimization, "Should recommend I/O optimization")
	assert.True(t, foundOverallReview, "Should recommend overall review")
}

func TestCalculateCoreUtilizationVariance(t *testing.T) {
	calc := NewEfficiencyCalculator()
	
	// Test balanced cores
	balancedCores := []interfaces.CPUCoreMetric{
		{CoreID: 0, Utilization: 80.0},
		{CoreID: 1, Utilization: 82.0},
		{CoreID: 2, Utilization: 78.0},
		{CoreID: 3, Utilization: 80.0},
	}
	
	variance := calc.calculateCoreUtilizationVariance(balancedCores)
	assert.Less(t, variance, 5.0) // Low variance
	
	// Test imbalanced cores
	imbalancedCores := []interfaces.CPUCoreMetric{
		{CoreID: 0, Utilization: 95.0},
		{CoreID: 1, Utilization: 90.0},
		{CoreID: 2, Utilization: 20.0},
		{CoreID: 3, Utilization: 25.0},
	}
	
	variance = calc.calculateCoreUtilizationVariance(imbalancedCores)
	assert.Greater(t, variance, 30.0) // High variance
}
