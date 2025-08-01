// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

// AdapterPerformanceTestSuite tests performance characteristics of adapters
type AdapterPerformanceTestSuite struct {
	suite.Suite
	clients   map[string]slurm.SlurmClient
	versions  []string
	serverURL string
	token     string
}

// PerformanceMetrics holds performance measurement data
type PerformanceMetrics struct {
	Operation     string
	Version       string
	Duration      time.Duration
	Success       bool
	Error         error
	MemoryBefore  uint64
	MemoryAfter   uint64
	MemoryDelta   int64
}

// BenchmarkResult holds aggregated benchmark results
type BenchmarkResult struct {
	Operation     string
	Version       string
	Iterations    int
	TotalDuration time.Duration
	AvgDuration   time.Duration
	MinDuration   time.Duration
	MaxDuration   time.Duration
	SuccessRate   float64
	ErrorCount    int
	MemoryUsage   int64
}

// SetupSuite initializes performance testing
func (suite *AdapterPerformanceTestSuite) SetupSuite() {
	// Check if performance testing is enabled
	if os.Getenv("SLURM_PERFORMANCE_TEST") != "true" {
		suite.T().Skip("Performance tests disabled. Set SLURM_PERFORMANCE_TEST=true to enable")
	}

	// Get server configuration
	suite.serverURL = os.Getenv("SLURM_SERVER_URL")
	if suite.serverURL == "" {
		suite.serverURL = "http://rocky9:6820"
	}

	// Get JWT token
	suite.token = os.Getenv("SLURM_JWT_TOKEN")
	if suite.token == "" {
		token, err := fetchJWTTokenViaSSH()
		require.NoError(suite.T(), err, "Failed to fetch JWT token")
		suite.token = token
	}

	// Initialize clients for performance testing
	suite.versions = []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}
	suite.clients = make(map[string]slurm.SlurmClient)

	ctx := context.Background()
	for _, version := range suite.versions {
		client, err := slurm.NewClientWithVersion(ctx, version,
			slurm.WithBaseURL(suite.serverURL),
			slurm.WithAuth(auth.NewTokenAuth(suite.token)),
			slurm.WithConfig(&config.Config{
				Timeout:            30 * time.Second,
				MaxRetries:         1, // Single retry for performance testing
				Debug:              false, // Disable debug for performance
				InsecureSkipVerify: true,
			}),
		)
		
		if err != nil {
			suite.T().Logf("Failed to create client for version %s: %v", version, err)
			continue
		}
		
		suite.clients[version] = client
		suite.T().Logf("Performance client created for version %s", version)
	}

	require.NotEmpty(suite.T(), suite.clients, "At least one client must be created")
}

// TearDownSuite cleans up performance test resources
func (suite *AdapterPerformanceTestSuite) TearDownSuite() {
	for version, client := range suite.clients {
		if client != nil {
			client.Close()
			suite.T().Logf("Closed performance client for version %s", version)
		}
	}
}

// measureMemory returns current memory usage
func (suite *AdapterPerformanceTestSuite) measureMemory() uint64 {
	var m runtime.MemStats
	runtime.GC() // Force garbage collection for accurate measurement
	runtime.ReadMemStats(&m)
	return m.Alloc
}

// runOperation executes an operation and measures performance
func (suite *AdapterPerformanceTestSuite) runOperation(version string, operation string, fn func() error) PerformanceMetrics {
	memBefore := suite.measureMemory()
	start := time.Now()
	
	err := fn()
	
	duration := time.Since(start)
	memAfter := suite.measureMemory()
	
	return PerformanceMetrics{
		Operation:     operation,
		Version:       version,
		Duration:      duration,
		Success:       err == nil,
		Error:         err,
		MemoryBefore:  memBefore,
		MemoryAfter:   memAfter,
		MemoryDelta:   int64(memAfter) - int64(memBefore),
	}
}

// TestPingPerformance benchmarks ping operations across versions
func (suite *AdapterPerformanceTestSuite) TestPingPerformance() {
	ctx := context.Background()
	iterations := 20
	
	suite.T().Logf("=== Ping Performance Test (%d iterations) ===", iterations)
	
	results := make(map[string][]PerformanceMetrics)
	
	for version, client := range suite.clients {
		suite.T().Logf("Testing ping performance for %s...", version)
		
		versionResults := make([]PerformanceMetrics, iterations)
		
		for i := 0; i < iterations; i++ {
			metric := suite.runOperation(version, "ping", func() error {
				return client.Info().Ping(ctx)
			})
			versionResults[i] = metric
			
			if !metric.Success {
				suite.T().Logf("  Iteration %d failed: %v", i+1, metric.Error)
			}
			
			// Small delay between requests
			time.Sleep(100 * time.Millisecond)
		}
		
		results[version] = versionResults
	}
	
	// Analyze and report results
	suite.analyzeAndReportResults("Ping", results)
}

// TestQoSListingPerformance benchmarks QoS listing operations
func (suite *AdapterPerformanceTestSuite) TestQoSListingPerformance() {
	ctx := context.Background()
	iterations := 10
	
	suite.T().Logf("=== QoS Listing Performance Test (%d iterations) ===", iterations)
	
	results := make(map[string][]PerformanceMetrics)
	
	for version, client := range suite.clients {
		suite.T().Logf("Testing QoS listing performance for %s...", version)
		
		versionResults := make([]PerformanceMetrics, iterations)
		
		for i := 0; i < iterations; i++ {
			metric := suite.runOperation(version, "qos_list", func() error {
				_, err := client.QoS().List(ctx, &interfaces.ListQoSOptions{
					Limit: 10,
				})
				return err
			})
			versionResults[i] = metric
			
			if !metric.Success {
				suite.T().Logf("  Iteration %d failed: %v", i+1, metric.Error)
			}
			
			time.Sleep(200 * time.Millisecond)
		}
		
		results[version] = versionResults
	}
	
	suite.analyzeAndReportResults("QoS Listing", results)
}

// TestJobListingPerformance benchmarks job listing operations
func (suite *AdapterPerformanceTestSuite) TestJobListingPerformance() {
	ctx := context.Background()
	iterations := 10
	
	suite.T().Logf("=== Job Listing Performance Test (%d iterations) ===", iterations)
	
	results := make(map[string][]PerformanceMetrics)
	
	for version, client := range suite.clients {
		suite.T().Logf("Testing job listing performance for %s...", version)
		
		versionResults := make([]PerformanceMetrics, iterations)
		
		for i := 0; i < iterations; i++ {
			metric := suite.runOperation(version, "job_list", func() error {
				_, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{
					Limit: 20,
				})
				return err
			})
			versionResults[i] = metric
			
			if !metric.Success {
				suite.T().Logf("  Iteration %d failed: %v", i+1, metric.Error)
			}
			
			time.Sleep(300 * time.Millisecond)
		}
		
		results[version] = versionResults
	}
	
	suite.analyzeAndReportResults("Job Listing", results)
}

// TestConcurrentOperations tests performance under concurrent load
func (suite *AdapterPerformanceTestSuite) TestConcurrentOperations() {
	ctx := context.Background()
	concurrency := 5
	operationsPerWorker := 10
	
	suite.T().Logf("=== Concurrent Operations Test (%d workers, %d ops each) ===", 
		concurrency, operationsPerWorker)
	
	for version, client := range suite.clients {
		suite.T().Logf("Testing concurrent operations for %s...", version)
		
		var wg sync.WaitGroup
		results := make(chan PerformanceMetrics, concurrency*operationsPerWorker)
		
		start := time.Now()
		
		for worker := 0; worker < concurrency; worker++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				
				for op := 0; op < operationsPerWorker; op++ {
					metric := suite.runOperation(version, fmt.Sprintf("concurrent_ping_w%d", workerID), func() error {
						return client.Info().Ping(ctx)
					})
					results <- metric
				}
			}(worker)
		}
		
		wg.Wait()
		close(results)
		
		totalDuration := time.Since(start)
		
		// Collect results
		var metrics []PerformanceMetrics
		successCount := 0
		
		for metric := range results {
			metrics = append(metrics, metric)
			if metric.Success {
				successCount++
			}
		}
		
		totalOps := concurrency * operationsPerWorker
		successRate := float64(successCount) / float64(totalOps) * 100
		throughput := float64(totalOps) / totalDuration.Seconds()
		
		suite.T().Logf("  Version %s concurrent results:", version)
		suite.T().Logf("    Total Operations: %d", totalOps)
		suite.T().Logf("    Success Rate: %.1f%% (%d/%d)", successRate, successCount, totalOps)
		suite.T().Logf("    Total Duration: %v", totalDuration)
		suite.T().Logf("    Throughput: %.2f ops/sec", throughput)
		
		// Calculate average latency for successful operations
		var totalLatency time.Duration
		successfulOps := 0
		
		for _, metric := range metrics {
			if metric.Success {
				totalLatency += metric.Duration
				successfulOps++
			}
		}
		
		if successfulOps > 0 {
			avgLatency := totalLatency / time.Duration(successfulOps)
			suite.T().Logf("    Average Latency: %v", avgLatency)
		}
		
		// Basic performance requirements
		suite.Greater(successRate, 80.0, "Success rate should be > 80%% for %s", version)
		suite.Greater(throughput, 1.0, "Throughput should be > 1 ops/sec for %s", version)
	}
}

// TestMemoryUsagePatterns tests memory usage patterns
func (suite *AdapterPerformanceTestSuite) TestMemoryUsagePatterns() {
	ctx := context.Background()
	iterations := 50
	
	suite.T().Logf("=== Memory Usage Patterns Test (%d iterations) ===", iterations)
	
	for version, client := range suite.clients {
		suite.T().Logf("Testing memory usage for %s...", version)
		
		initialMemory := suite.measureMemory()
		var memoryDeltas []int64
		
		for i := 0; i < iterations; i++ {
			metric := suite.runOperation(version, "memory_test", func() error {
				// Perform multiple operations to simulate real usage
				err := client.Info().Ping(ctx)
				if err != nil {
					return err
				}
				
				_, err = client.QoS().List(ctx, &interfaces.ListQoSOptions{Limit: 5})
				if err != nil {
					return err
				}
				
				_, err = client.Jobs().List(ctx, &interfaces.ListJobsOptions{Limit: 5})
				return err
			})
			
			if metric.Success {
				memoryDeltas = append(memoryDeltas, metric.MemoryDelta)
			}
			
			// Periodic garbage collection
			if i%10 == 0 {
				runtime.GC()
			}
		}
		
		finalMemory := suite.measureMemory()
		totalMemoryGrowth := int64(finalMemory) - int64(initialMemory)
		
		// Calculate memory statistics
		var totalDelta int64
		var maxDelta int64
		var minDelta int64 = int64(^uint64(0) >> 1) // Max int64
		
		for _, delta := range memoryDeltas {
			totalDelta += delta
			if delta > maxDelta {
				maxDelta = delta
			}
			if delta < minDelta {
				minDelta = delta
			}
		}
		
		avgDelta := totalDelta / int64(len(memoryDeltas))
		
		suite.T().Logf("  Version %s memory analysis:", version)
		suite.T().Logf("    Initial Memory: %d bytes", initialMemory)
		suite.T().Logf("    Final Memory: %d bytes", finalMemory)
		suite.T().Logf("    Total Growth: %d bytes", totalMemoryGrowth)
		suite.T().Logf("    Average Delta: %d bytes", avgDelta)
		suite.T().Logf("    Max Delta: %d bytes", maxDelta)
		suite.T().Logf("    Min Delta: %d bytes", minDelta)
		
		// Memory should not grow excessively
		memoryGrowthMB := float64(totalMemoryGrowth) / (1024 * 1024)
		suite.Less(memoryGrowthMB, 100.0, "Memory growth should be < 100MB for %s", version)
	}
}

// TestLargeDataHandling tests performance with large data sets
func (suite *AdapterPerformanceTestSuite) TestLargeDataHandling() {
	ctx := context.Background()
	
	suite.T().Log("=== Large Data Handling Test ===")
	
	largeLimits := []int{50, 100, 200}
	
	for version, client := range suite.clients {
		suite.T().Logf("Testing large data handling for %s...", version)
		
		for _, limit := range largeLimits {
			suite.T().Logf("  Testing with limit %d...", limit)
			
			// Test large job listing
			jobMetric := suite.runOperation(version, fmt.Sprintf("large_jobs_%d", limit), func() error {
				_, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{
					Limit: limit,
				})
				return err
			})
			
			if jobMetric.Success {
				suite.T().Logf("    Jobs (%d): %v (memory: %d bytes)", 
					limit, jobMetric.Duration, jobMetric.MemoryDelta)
				
				// Performance expectations for large data
				suite.Less(jobMetric.Duration, 10*time.Second, 
					"Large job listing should complete within 10s for %s", version)
			} else {
				suite.T().Logf("    Jobs (%d): FAILED - %v", limit, jobMetric.Error)
			}
			
			// Test large node listing
			nodeMetric := suite.runOperation(version, fmt.Sprintf("large_nodes_%d", limit), func() error {
				_, err := client.Nodes().List(ctx, &interfaces.ListNodesOptions{
					Limit: limit,
				})
				return err
			})
			
			if nodeMetric.Success {
				suite.T().Logf("    Nodes (%d): %v (memory: %d bytes)", 
					limit, nodeMetric.Duration, nodeMetric.MemoryDelta)
			} else {
				suite.T().Logf("    Nodes (%d): FAILED - %v", limit, nodeMetric.Error)
			}
		}
	}
}

// analyzeAndReportResults analyzes performance metrics and generates reports
func (suite *AdapterPerformanceTestSuite) analyzeAndReportResults(testName string, results map[string][]PerformanceMetrics) {
	suite.T().Logf("\n=== %s Performance Analysis ===", testName)
	
	benchmarks := make(map[string]BenchmarkResult)
	
	for version, metrics := range results {
		if len(metrics) == 0 {
			continue
		}
		
		var totalDuration time.Duration
		var minDuration = time.Duration(^uint64(0) >> 1) // Max duration
		var maxDuration time.Duration
		var successCount int
		var totalMemory int64
		
		for _, metric := range metrics {
			totalDuration += metric.Duration
			
			if metric.Duration < minDuration {
				minDuration = metric.Duration
			}
			if metric.Duration > maxDuration {
				maxDuration = metric.Duration
			}
			
			if metric.Success {
				successCount++
			}
			
			totalMemory += metric.MemoryDelta
		}
		
		benchmark := BenchmarkResult{
			Operation:     testName,
			Version:       version,
			Iterations:    len(metrics),
			TotalDuration: totalDuration,
			AvgDuration:   totalDuration / time.Duration(len(metrics)),
			MinDuration:   minDuration,
			MaxDuration:   maxDuration,
			SuccessRate:   float64(successCount) / float64(len(metrics)) * 100,
			ErrorCount:    len(metrics) - successCount,
			MemoryUsage:   totalMemory / int64(len(metrics)),
		}
		
		benchmarks[version] = benchmark
	}
	
	// Report results
	suite.T().Log("\nPerformance Summary:")
	suite.T().Log("Version\t\tAvg Duration\tMin Duration\tMax Duration\tSuccess Rate\tAvg Memory")
	suite.T().Log("-------\t\t------------\t------------\t------------\t------------\t----------")
	
	for version, benchmark := range benchmarks {
		suite.T().Logf("%s\t\t%v\t\t%v\t\t%v\t\t%.1f%%\t\t%d bytes",
			version,
			benchmark.AvgDuration.Round(time.Millisecond),
			benchmark.MinDuration.Round(time.Millisecond),
			benchmark.MaxDuration.Round(time.Millisecond),
			benchmark.SuccessRate,
			benchmark.MemoryUsage,
		)
		
		// Basic performance requirements
		suite.Greater(benchmark.SuccessRate, 90.0, "Success rate should be > 90%% for %s %s", testName, version)
		suite.Less(benchmark.AvgDuration, 5*time.Second, "Average duration should be < 5s for %s %s", testName, version)
	}
	
	// Compare versions if multiple are available
	if len(benchmarks) > 1 {
		suite.T().Log("\nVersion Comparison:")
		var fastestVersion string
		var fastestDuration time.Duration = time.Duration(^uint64(0) >> 1)
		
		for version, benchmark := range benchmarks {
			if benchmark.AvgDuration < fastestDuration && benchmark.SuccessRate > 90 {
				fastestDuration = benchmark.AvgDuration
				fastestVersion = version
			}
		}
		
		if fastestVersion != "" {
			suite.T().Logf("Fastest version for %s: %s (%v avg)", testName, fastestVersion, fastestDuration.Round(time.Millisecond))
		}
	}
}

// TestAdapterPerformanceSuite runs the adapter performance test suite
func TestAdapterPerformanceSuite(t *testing.T) {
	suite.Run(t, new(AdapterPerformanceTestSuite))
}
