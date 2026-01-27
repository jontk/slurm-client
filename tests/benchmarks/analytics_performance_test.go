// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package benchmarks

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jontk/slurm-client/tests/mocks"
)

// BenchmarkJobAnalyticsOverhead measures the performance overhead of analytics collection
// to ensure it stays under 5% of baseline job operations
func BenchmarkJobAnalyticsOverhead(b *testing.B) {
	// Setup mock server
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	baseURL := mockServer.URL()
	jobID := "1001"

	// Benchmark baseline job operation (getting job details)
	b.Run("Baseline_JobGet", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			resp, err := makeHTTPRequest(fmt.Sprintf("%s/slurm/v0.0.42/job/%s", baseURL, jobID))
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
	})

	// Benchmark analytics operations
	analyticsEndpoints := []struct {
		name     string
		endpoint string
	}{
		{"Analytics_Utilization", fmt.Sprintf("%s/slurm/v0.0.42/job/%s/utilization", baseURL, jobID)},
		{"Analytics_Efficiency", fmt.Sprintf("%s/slurm/v0.0.42/job/%s/efficiency", baseURL, jobID)},
		{"Analytics_Performance", fmt.Sprintf("%s/slurm/v0.0.42/job/%s/performance", baseURL, jobID)},
		{"Analytics_LiveMetrics", fmt.Sprintf("%s/slurm/v0.0.42/job/%s/live_metrics", baseURL, jobID)},
		{"Analytics_ResourceTrends", fmt.Sprintf("%s/slurm/v0.0.42/job/%s/resource_trends", baseURL, jobID)},
	}

	for _, endpoint := range analyticsEndpoints {
		b.Run(endpoint.name, func(b *testing.B) {
			b.ResetTimer()
			for range b.N {
				resp, err := makeHTTPRequest(endpoint.endpoint)
				if err != nil {
					b.Fatal(err)
				}
				resp.Body.Close()
			}
		})
	}
}

// BenchmarkAnalyticsCollectionOverhead measures the overhead of collecting multiple analytics
// compared to single operations
func BenchmarkAnalyticsCollectionOverhead(b *testing.B) {
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	baseURL := mockServer.URL()
	jobID := "1001"

	// Benchmark single analytics operation
	b.Run("Single_UtilizationCall", func(b *testing.B) {
		endpoint := fmt.Sprintf("%s/slurm/v0.0.42/job/%s/utilization", baseURL, jobID)
		b.ResetTimer()
		for range b.N {
			resp, err := makeHTTPRequest(endpoint)
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
	})

	// Benchmark collecting all analytics for a job
	b.Run("Complete_JobAnalytics", func(b *testing.B) {
		endpoints := []string{
			fmt.Sprintf("%s/slurm/v0.0.42/job/%s/utilization", baseURL, jobID),
			fmt.Sprintf("%s/slurm/v0.0.42/job/%s/efficiency", baseURL, jobID),
			fmt.Sprintf("%s/slurm/v0.0.42/job/%s/performance", baseURL, jobID),
			fmt.Sprintf("%s/slurm/v0.0.42/job/%s/live_metrics", baseURL, jobID),
		}

		b.ResetTimer()
		for range b.N {
			for _, endpoint := range endpoints {
				resp, err := makeHTTPRequest(endpoint)
				if err != nil {
					b.Fatal(err)
				}
				resp.Body.Close()
			}
		}
	})

	// Benchmark parallel analytics collection
	b.Run("Parallel_JobAnalytics", func(b *testing.B) {
		endpoints := []string{
			fmt.Sprintf("%s/slurm/v0.0.42/job/%s/utilization", baseURL, jobID),
			fmt.Sprintf("%s/slurm/v0.0.42/job/%s/efficiency", baseURL, jobID),
			fmt.Sprintf("%s/slurm/v0.0.42/job/%s/performance", baseURL, jobID),
			fmt.Sprintf("%s/slurm/v0.0.42/job/%s/live_metrics", baseURL, jobID),
		}

		b.ResetTimer()
		for range b.N {
			var wg sync.WaitGroup
			errChan := make(chan error, len(endpoints))
			for _, endpoint := range endpoints {
				wg.Add(1)
				go func(url string) {
					defer wg.Done()
					resp, err := makeHTTPRequest(url)
					if err != nil {
						errChan <- err
						return
					}
					_ = resp.Body.Close()
				}(endpoint)
			}
			wg.Wait()
			close(errChan)
			// Check for errors after goroutines complete
			if err := <-errChan; err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkAnalyticsMemoryUsage measures memory overhead of analytics operations
func BenchmarkAnalyticsMemoryUsage(b *testing.B) {
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	baseURL := mockServer.URL()
	jobID := "1001"

	b.Run("Memory_UtilizationAnalytics", func(b *testing.B) {
		endpoint := fmt.Sprintf("%s/slurm/v0.0.42/job/%s/utilization", baseURL, jobID)
		b.ResetTimer()
		b.ReportAllocs()
		for range b.N {
			resp, err := makeHTTPRequest(endpoint)
			if err != nil {
				b.Fatal(err)
			}
			// Read and process response to measure actual memory usage
			_, err = readResponseBody(resp)
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
	})

	b.Run("Memory_EfficiencyAnalytics", func(b *testing.B) {
		endpoint := fmt.Sprintf("%s/slurm/v0.0.42/job/%s/efficiency", baseURL, jobID)
		b.ResetTimer()
		b.ReportAllocs()
		for range b.N {
			resp, err := makeHTTPRequest(endpoint)
			if err != nil {
				b.Fatal(err)
			}
			_, err = readResponseBody(resp)
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
	})

	b.Run("Memory_PerformanceAnalytics", func(b *testing.B) {
		endpoint := fmt.Sprintf("%s/slurm/v0.0.42/job/%s/performance", baseURL, jobID)
		b.ResetTimer()
		b.ReportAllocs()
		for range b.N {
			resp, err := makeHTTPRequest(endpoint)
			if err != nil {
				b.Fatal(err)
			}
			_, err = readResponseBody(resp)
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
	})
}

// BenchmarkAnalyticsConcurrentLoad tests analytics performance under concurrent load
func BenchmarkAnalyticsConcurrentLoad(b *testing.B) {
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	baseURL := mockServer.URL()
	jobID := "1001"

	concurrencyLevels := []int{1, 5, 10, 25, 50, 100}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrent_%d_Analytics", concurrency), func(b *testing.B) {
			endpoint := fmt.Sprintf("%s/slurm/v0.0.42/job/%s/utilization", baseURL, jobID)

			b.SetParallelism(concurrency)
			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					resp, err := makeHTTPRequest(endpoint)
					if err != nil {
						b.Fatal(err)
					}
					resp.Body.Close()
				}
			})
		})
	}
}

// BenchmarkAnalyticsLatencyDistribution measures latency distribution for analytics operations
func BenchmarkAnalyticsLatencyDistribution(b *testing.B) {
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	baseURL := mockServer.URL()
	jobID := "1001"

	// Note: For meaningful distribution analysis, run with -benchtime=100x or higher
	// Example: go test -bench=BenchmarkAnalyticsLatencyDistribution -benchtime=100x

	b.Run("Latency_Distribution", func(b *testing.B) {
		endpoint := fmt.Sprintf("%s/slurm/v0.0.42/job/%s/utilization", baseURL, jobID)
		latencies := make([]time.Duration, b.N)

		b.ResetTimer()
		start := time.Now()

		for i := range b.N {
			reqStart := time.Now()
			resp, err := makeHTTPRequest(endpoint)
			latencies[i] = time.Since(reqStart)

			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}

		totalTime := time.Since(start)

		b.StopTimer()

		// Calculate latency statistics
		minLatency, maxLatency, avg, p95, p99 := calculateLatencyStats(latencies)

		b.ReportMetric(float64(totalTime.Nanoseconds())/float64(b.N), "ns/op")
		b.ReportMetric(float64(minLatency.Nanoseconds()), "min-ns")
		b.ReportMetric(float64(maxLatency.Nanoseconds()), "max-ns")
		b.ReportMetric(float64(avg.Nanoseconds()), "avg-ns")
		b.ReportMetric(float64(p95.Nanoseconds()), "p95-ns")
		b.ReportMetric(float64(p99.Nanoseconds()), "p99-ns")
	})
}

// BenchmarkAnalyticsVersionComparison compares performance across API versions
func BenchmarkAnalyticsVersionComparison(b *testing.B) {
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}

	for _, version := range versions {
		b.Run("Version_"+version, func(b *testing.B) {
			mockServer := mocks.NewMockSlurmServerForVersion(version)
			defer mockServer.Close()

			baseURL := mockServer.URL()
			jobID := "1001"
			endpoint := fmt.Sprintf("%s/slurm/%s/job/%s/utilization", baseURL, version, jobID)

			b.ResetTimer()
			for range b.N {
				resp, err := makeHTTPRequest(endpoint)
				if err != nil {
					b.Fatal(err)
				}
				resp.Body.Close()
			}
		})
	}
}

// BenchmarkAnalyticsDataProcessing measures the overhead of processing analytics data
func BenchmarkAnalyticsDataProcessing(b *testing.B) {
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	baseURL := mockServer.URL()
	jobID := "1001"

	// Benchmark raw request without processing
	b.Run("Raw_Request", func(b *testing.B) {
		endpoint := fmt.Sprintf("%s/slurm/v0.0.42/job/%s/utilization", baseURL, jobID)
		b.ResetTimer()
		for range b.N {
			resp, err := makeHTTPRequest(endpoint)
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
	})

	// Benchmark request with JSON parsing
	b.Run("With_JSON_Parsing", func(b *testing.B) {
		endpoint := fmt.Sprintf("%s/slurm/v0.0.42/job/%s/utilization", baseURL, jobID)
		b.ResetTimer()
		for range b.N {
			resp, err := makeHTTPRequest(endpoint)
			if err != nil {
				b.Fatal(err)
			}

			// Parse JSON response
			_, err = parseJSONResponse(resp)
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
	})

	// Benchmark request with full data processing
	b.Run("With_Full_Processing", func(b *testing.B) {
		endpoint := fmt.Sprintf("%s/slurm/v0.0.42/job/%s/efficiency", baseURL, jobID)
		b.ResetTimer()
		for range b.N {
			resp, err := makeHTTPRequest(endpoint)
			if err != nil {
				b.Fatal(err)
			}

			// Parse and process efficiency data
			data, err := parseJSONResponse(resp)
			if err != nil {
				b.Fatal(err)
			}

			// Simulate processing efficiency calculations
			processEfficiencyData(data)
			resp.Body.Close()
		}
	})
}

// TestAnalyticsOverheadCompliance reports analytics overhead without assertions
// This test collects performance data for trend analysis but does not fail on timing
func TestAnalyticsOverheadCompliance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping overhead compliance test in short mode")
	}

	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	baseURL := mockServer.URL()
	jobID := "1001"

	// Measure baseline performance
	baselineTime := measureOperationTime(t, func() error {
		resp, err := makeHTTPRequest(fmt.Sprintf("%s/slurm/v0.0.42/job/%s", baseURL, jobID))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return nil
	}, 100) // Run 100 iterations for stable measurement

	// Measure analytics operations
	analyticsOperations := map[string]string{
		"utilization":  fmt.Sprintf("%s/slurm/v0.0.42/job/%s/utilization", baseURL, jobID),
		"efficiency":   fmt.Sprintf("%s/slurm/v0.0.42/job/%s/efficiency", baseURL, jobID),
		"performance":  fmt.Sprintf("%s/slurm/v0.0.42/job/%s/performance", baseURL, jobID),
		"live_metrics": fmt.Sprintf("%s/slurm/v0.0.42/job/%s/live_metrics", baseURL, jobID),
	}

	for opName, endpoint := range analyticsOperations {
		t.Run("Overhead_"+opName, func(t *testing.T) {
			analyticsTime := measureOperationTime(t, func() error {
				resp, err := makeHTTPRequest(endpoint)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				// Include JSON parsing in overhead measurement
				_, err = parseJSONResponse(resp)
				return err
			}, 100)

			// Calculate overhead percentage
			overhead := float64(analyticsTime-baselineTime) / float64(baselineTime) * 100

			// Report results without failing - data collection only
			t.Logf("Operation: %s", opName)
			t.Logf("Baseline time: %v", baselineTime)
			t.Logf("Analytics time: %v", analyticsTime)
			t.Logf("Overhead: %.2f%%", overhead)
		})
	}

	// Test combined analytics overhead
	t.Run("Combined_Analytics_Overhead", func(t *testing.T) {
		combinedAnalyticsTime := measureOperationTime(t, func() error {
			for _, endpoint := range analyticsOperations {
				resp, err := makeHTTPRequest(endpoint)
				if err != nil {
					return err
				}
				_, err = parseJSONResponse(resp)
				resp.Body.Close()
				if err != nil {
					return err
				}
			}
			return nil
		}, 50) // Fewer iterations for combined test

		// For combined operations, we expect the overhead per operation to be similar
		avgAnalyticsTime := combinedAnalyticsTime / time.Duration(len(analyticsOperations))
		overhead := float64(avgAnalyticsTime-baselineTime) / float64(baselineTime) * 100

		// Report results without failing - data collection only
		t.Logf("Combined analytics average time per operation: %v", avgAnalyticsTime)
		t.Logf("Combined analytics overhead: %.2f%%", overhead)
	})
}

// TestPerformanceSmokeTest catches catastrophic performance regressions (>10x slower)
// This is a sanity check that only fails on major issues, not timing variability
func TestPerformanceSmokeTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance smoke test in short mode")
	}

	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	baseURL := mockServer.URL()
	jobID := "1001"
	endpoint := fmt.Sprintf("%s/slurm/v0.0.42/job/%s", baseURL, jobID)

	// Perform 100 requests to ensure stability
	start := time.Now()
	for i := 0; i < 100; i++ {
		resp, err := makeHTTPRequest(endpoint)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}
		resp.Body.Close()
	}
	elapsed := time.Since(start)

	// Only fail if catastrophically slow (>30s for 100 requests)
	// Normal should be <1s, but we allow huge margin for CI variability
	if elapsed > 30*time.Second {
		t.Fatalf("Performance catastrophically degraded: %v for 100 requests (threshold: 30s)", elapsed)
	}

	t.Logf("✅ 100 requests completed in %v (smoke test passed)", elapsed)
}

// TestAnalyticsScalability validates that analytics performance scales properly
func TestAnalyticsScalability(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping scalability test in short mode")
	}

	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	baseURL := mockServer.URL()
	endpoint := baseURL + "/slurm/v0.0.42/job/1001/utilization"

	// Test different request volumes
	volumes := []int{1, 10, 50, 100, 200}

	for _, volume := range volumes {
		t.Run(fmt.Sprintf("Volume_%d", volume), func(t *testing.T) {
			start := time.Now()

			var wg sync.WaitGroup
			errors := make(chan error, volume)

			for range volume {
				wg.Add(1)
				go func() {
					defer wg.Done()
					resp, err := makeHTTPRequest(endpoint)
					if err != nil {
						errors <- err
						return
					}
					defer resp.Body.Close()

					_, err = parseJSONResponse(resp)
					if err != nil {
						errors <- err
					}
				}()
			}

			wg.Wait()
			close(errors)

			duration := time.Since(start)

			// Check for errors
			var errorCount int
			for err := range errors {
				if err != nil {
					t.Errorf("Request failed: %v", err)
					errorCount++
				}
			}

			if errorCount > 0 {
				t.Errorf("%d out of %d requests failed", errorCount, volume)
			}

			// Calculate performance metrics
			avgTimePerRequest := duration / time.Duration(volume)
			requestsPerSecond := float64(volume) / duration.Seconds()

			t.Logf("Volume: %d requests", volume)
			t.Logf("Total time: %v", duration)
			t.Logf("Average time per request: %v", avgTimePerRequest)
			t.Logf("Requests per second: %.2f", requestsPerSecond)

			// Performance expectations
			if avgTimePerRequest > 100*time.Millisecond {
				t.Errorf("Average response time %v exceeds 100ms threshold", avgTimePerRequest)
			}

			if requestsPerSecond < 50 && volume >= 100 {
				t.Errorf("Request throughput %.2f req/s is below expected 50 req/s for volume %d",
					requestsPerSecond, volume)
			}
		})
	}
}
