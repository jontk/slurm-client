// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package benchmarks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"testing"
	"time"
)

// makeHTTPRequest performs a HTTP GET request and returns the response
func makeHTTPRequest(url string) (*http.Response, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	return client.Get(url)
}

// readResponseBody reads the entire response body and returns it as bytes
func readResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// parseJSONResponse parses JSON response body into a map
func parseJSONResponse(resp *http.Response) (map[string]interface{}, error) {
	defer resp.Body.Close()

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

// processEfficiencyData simulates processing efficiency data for benchmarking
func processEfficiencyData(data map[string]interface{}) {
	// Simulate data processing overhead
	if efficiency, ok := data["efficiency"].(map[string]interface{}); ok {
		// Simulate calculations
		if cpuEff, ok := efficiency["cpu_efficiency"].(float64); ok {
			_ = cpuEff * 1.1 // Simulate processing
		}
		if memEff, ok := efficiency["memory_efficiency"].(float64); ok {
			_ = memEff * 1.1 // Simulate processing
		}
		if overallEff, ok := efficiency["overall_efficiency_score"].(float64); ok {
			_ = overallEff * 1.1 // Simulate processing
		}
	}
}

// calculateLatencyStats calculates statistical measures for latency data
func calculateLatencyStats(latencies []time.Duration) (minLatency, maxLatency, avg, p95, p99 time.Duration) {
	if len(latencies) == 0 {
		return 0, 0, 0, 0, 0
	}

	// Sort latencies for percentile calculations
	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	// Min and max
	minLatency = sorted[0]
	maxLatency = sorted[len(sorted)-1]

	// Average
	var total time.Duration
	for _, lat := range latencies {
		total += lat
	}
	avg = total / time.Duration(len(latencies))

	// 95th percentile
	p95Index := int(float64(len(sorted)) * 0.95)
	if p95Index >= len(sorted) {
		p95Index = len(sorted) - 1
	}
	p95 = sorted[p95Index]

	// 99th percentile
	p99Index := int(float64(len(sorted)) * 0.99)
	if p99Index >= len(sorted) {
		p99Index = len(sorted) - 1
	}
	p99 = sorted[p99Index]

	return minLatency, maxLatency, avg, p95, p99
}

// measureOperationTime measures the average time for a repeated operation
func measureOperationTime(t *testing.T, operation func() error, iterations int) time.Duration {
	start := time.Now()

	for i := 0; i < iterations; i++ {
		if err := operation(); err != nil {
			t.Fatalf("Operation failed on iteration %d: %v", i, err)
		}
	}

	totalTime := time.Since(start)
	return totalTime / time.Duration(iterations)
}

// PerformanceMetrics holds performance measurement results
type PerformanceMetrics struct {
	OperationName     string
	TotalTime         time.Duration
	AverageTime       time.Duration
	MinTime           time.Duration
	MaxTime           time.Duration
	RequestsPerSecond float64
	MemoryAllocations int64
	MemoryBytes       int64
}

// MeasurePerformance performs comprehensive performance measurement of an operation
func MeasurePerformance(name string, operation func() error, iterations int) (*PerformanceMetrics, error) {
	metrics := &PerformanceMetrics{
		OperationName: name,
	}

	times := make([]time.Duration, iterations)
	start := time.Now()

	for i := 0; i < iterations; i++ {
		iterStart := time.Now()
		if err := operation(); err != nil {
			return nil, fmt.Errorf("operation failed on iteration %d: %w", i, err)
		}
		times[i] = time.Since(iterStart)
	}

	metrics.TotalTime = time.Since(start)
	metrics.AverageTime = metrics.TotalTime / time.Duration(iterations)
	metrics.RequestsPerSecond = float64(iterations) / metrics.TotalTime.Seconds()

	// Calculate min/max
	metrics.MinTime = times[0]
	metrics.MaxTime = times[0]
	for _, t := range times {
		if t < metrics.MinTime {
			metrics.MinTime = t
		}
		if t > metrics.MaxTime {
			metrics.MaxTime = t
		}
	}

	return metrics, nil
}

// ComparePerformance compares two performance metrics and calculates overhead percentage
func ComparePerformance(baseline, measured *PerformanceMetrics) float64 {
	if baseline.AverageTime == 0 {
		return 0
	}
	return ((float64(measured.AverageTime) - float64(baseline.AverageTime)) / float64(baseline.AverageTime)) * 100
}

// PerformanceReport generates a formatted performance report
func PerformanceReport(metrics *PerformanceMetrics) string {
	return fmt.Sprintf(`Performance Report: %s
  Total Time: %v
  Average Time: %v
  Min Time: %v
  Max Time: %v
  Requests/Second: %.2f
  Memory Allocations: %d
  Memory Bytes: %d`,
		metrics.OperationName,
		metrics.TotalTime,
		metrics.AverageTime,
		metrics.MinTime,
		metrics.MaxTime,
		metrics.RequestsPerSecond,
		metrics.MemoryAllocations,
		metrics.MemoryBytes)
}

// LoadTestResult holds results from load testing
type LoadTestResult struct {
	ConcurrencyLevel   int
	TotalRequests      int
	SuccessfulRequests int
	FailedRequests     int
	TotalTime          time.Duration
	AverageLatency     time.Duration
	RequestsPerSecond  float64
	ErrorRate          float64
}

// RunLoadTest performs load testing with specified concurrency and duration
func RunLoadTest(endpoint string, concurrency int, duration time.Duration) (*LoadTestResult, error) {
	result := &LoadTestResult{
		ConcurrencyLevel: concurrency,
	}

	start := time.Now()
	deadline := start.Add(duration)

	requestChan := make(chan bool, concurrency*10) // Buffer for requests
	resultChan := make(chan bool, concurrency*10)  // Buffer for results

	// Start workers
	for i := 0; i < concurrency; i++ {
		go func() {
			for range requestChan {
				success := true
				reqStart := time.Now()

				resp, err := makeHTTPRequest(endpoint)
				if err != nil {
					success = false
				} else {
					_ = resp.Body.Close() // Ignore error during cleanup
					if resp.StatusCode != http.StatusOK {
						success = false
					}
				}

				// Track latency for successful requests
				if success {
					result.AverageLatency += time.Since(reqStart)
				}

				resultChan <- success
			}
		}()
	}

	// Send requests until deadline
	go func() {
		defer close(requestChan)
		for time.Now().Before(deadline) {
			requestChan <- true
			result.TotalRequests++
		}
	}()

	// Collect results
	for i := 0; i < result.TotalRequests; i++ {
		success := <-resultChan
		if success {
			result.SuccessfulRequests++
		} else {
			result.FailedRequests++
		}
	}

	result.TotalTime = time.Since(start)

	if result.SuccessfulRequests > 0 {
		result.AverageLatency = result.AverageLatency / time.Duration(result.SuccessfulRequests)
	}

	result.RequestsPerSecond = float64(result.TotalRequests) / result.TotalTime.Seconds()
	result.ErrorRate = (float64(result.FailedRequests) / float64(result.TotalRequests)) * 100

	return result, nil
}

// ValidatePerformanceThresholds validates that performance metrics meet specified thresholds
func ValidatePerformanceThresholds(metrics *PerformanceMetrics, thresholds map[string]interface{}) []string {
	var violations []string

	if maxAvgTime, ok := thresholds["max_average_time"].(time.Duration); ok {
		if metrics.AverageTime > maxAvgTime {
			violations = append(violations, fmt.Sprintf("Average time %v exceeds threshold %v",
				metrics.AverageTime, maxAvgTime))
		}
	}

	if minRPS, ok := thresholds["min_requests_per_second"].(float64); ok {
		if metrics.RequestsPerSecond < minRPS {
			violations = append(violations, fmt.Sprintf("Requests per second %.2f below threshold %.2f",
				metrics.RequestsPerSecond, minRPS))
		}
	}

	if maxMemory, ok := thresholds["max_memory_bytes"].(int64); ok {
		if metrics.MemoryBytes > maxMemory {
			violations = append(violations, fmt.Sprintf("Memory usage %d bytes exceeds threshold %d bytes",
				metrics.MemoryBytes, maxMemory))
		}
	}

	return violations
}

// BenchmarkConfig holds configuration for benchmark execution
type BenchmarkConfig struct {
	Iterations   int
	Concurrency  int
	Duration     time.Duration
	WarmupRounds int
	Endpoints    []string
	Thresholds   map[string]interface{}
}

// DefaultBenchmarkConfig returns a default benchmark configuration
func DefaultBenchmarkConfig() *BenchmarkConfig {
	return &BenchmarkConfig{
		Iterations:   1000,
		Concurrency:  10,
		Duration:     30 * time.Second,
		WarmupRounds: 100,
		Thresholds: map[string]interface{}{
			"max_average_time":        100 * time.Millisecond,
			"min_requests_per_second": 50.0,
			"max_memory_bytes":        int64(1024 * 1024), // 1MB
		},
	}
}
