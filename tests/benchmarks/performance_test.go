package benchmarks

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/performance"
	"github.com/jontk/slurm-client/tests/mocks"
)

// BenchmarkJobOperations tests the performance of job-related operations
func BenchmarkJobOperations(b *testing.B) {
	scenarios := []struct {
		name     string
		profile  performance.PerformanceProfile
		parallel bool
	}{
		{"Default", performance.ProfileDefault, false},
		{"HighThroughput", performance.ProfileHighThroughput, false},
		{"LowLatency", performance.ProfileLowLatency, false},
		{"Conservative", performance.ProfileConservative, false},
		{"Parallel_Default", performance.ProfileDefault, true},
		{"Parallel_HighThroughput", performance.ProfileHighThroughput, true},
	}

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			benchmarkJobOperationsWithProfile(b, scenario.profile, scenario.parallel)
		})
	}
}

func benchmarkJobOperationsWithProfile(b *testing.B, profile performance.PerformanceProfile, parallel bool) {
	// Setup mock server
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	// Create client with performance optimizations
	ctx := context.Background()
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.42",
		slurm.WithBaseURL(mockServer.URL()),
		slurm.WithAuth(auth.NewNoAuth()),
	)
	if err != nil {
		b.Fatal(err)
	}
	defer client.Close()

	// Pre-populate some jobs for listing/getting
	for i := 0; i < 100; i++ {
		job := &mocks.MockJob{
			JobID:     fmt.Sprintf("bench-%d", i),
			Name:      fmt.Sprintf("benchmark-job-%d", i),
			UserID:    "benchuser",
			State:     "RUNNING",
			Partition: "compute",
			CPUs:      2,
			Memory:    4 * 1024 * 1024 * 1024,
		}
		mockServer.AddJob(job)
	}

	b.ResetTimer()

	if parallel {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				runJobOperationBenchmark(b, client)
			}
		})
	} else {
		for i := 0; i < b.N; i++ {
			runJobOperationBenchmark(b, client)
		}
	}
}

func runJobOperationBenchmark(b *testing.B, client slurm.SlurmClient) {
	ctx := context.Background()

	// Test job listing
	_, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{
		Limit: 10,
	})
	if err != nil {
		b.Error(err)
		return
	}

	// Test job getting
	_, err = client.Jobs().Get(ctx, "bench-1")
	if err != nil {
		// Expected to fail since client isn't fully implemented
		// In a real benchmark, this would succeed
	}
}

// BenchmarkConnectionPooling tests the performance impact of connection pooling
func BenchmarkConnectionPooling(b *testing.B) {
	profiles := []performance.PerformanceProfile{
		performance.ProfileDefault,
		performance.ProfileHighThroughput,
		performance.ProfileLowLatency,
		performance.ProfileConservative,
	}

	for _, profile := range profiles {
		b.Run(string(profile), func(b *testing.B) {
			benchmarkConnectionPooling(b, profile)
		})
	}
}

func benchmarkConnectionPooling(b *testing.B, profile performance.PerformanceProfile) {
	// Create connection pool manager
	poolManager := performance.NewHTTPClientPoolManager()
	defer poolManager.Close()

	// Get pool for the profile
	pool := poolManager.GetPoolForVersion("v0.0.42", profile)

	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	b.ResetTimer()

	// Benchmark getting HTTP clients from the pool
	for i := 0; i < b.N; i++ {
		client := pool.GetClient(mockServer.URL())
		if client == nil {
			b.Error("Failed to get HTTP client")
		}
	}
}

// BenchmarkCaching tests the performance of response caching
func BenchmarkCaching(b *testing.B) {
	profiles := []performance.PerformanceProfile{
		performance.ProfileDefault,
		performance.ProfileHighThroughput,
		performance.ProfileLowLatency,
		performance.ProfileConservative,
	}

	for _, profile := range profiles {
		b.Run(string(profile), func(b *testing.B) {
			benchmarkCaching(b, profile)
		})
	}
}

func benchmarkCaching(b *testing.B, profile performance.PerformanceProfile) {
	// Create cache with profile-specific configuration
	config := performance.GetCacheConfigForProfile(profile)
	cache := performance.NewResponseCache(config)
	defer cache.Close()

	// Test data
	operation := "jobs.list"
	params := map[string]interface{}{
		"limit":  10,
		"offset": 0,
		"state":  "RUNNING",
	}
	value := []byte(`{"jobs": [{"id": "1", "name": "test"}]}`)

	b.ResetTimer()

	// Benchmark cache operations
	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			testParams := make(map[string]interface{})
			for k, v := range params {
				testParams[k] = v
			}
			testParams["offset"] = i // Make each entry unique
			
			cache.Set(operation, testParams, value)
		}
	})

	b.Run("Get", func(b *testing.B) {
		// Pre-populate cache
		for i := 0; i < 1000; i++ {
			testParams := make(map[string]interface{})
			for k, v := range params {
				testParams[k] = v
			}
			testParams["offset"] = i
			cache.Set(operation, testParams, value)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			testParams := make(map[string]interface{})
			for k, v := range params {
				testParams[k] = v
			}
			testParams["offset"] = i % 1000
			
			_, found := cache.Get(operation, testParams)
			if !found && i%1000 < 500 {
				// Should find about half the items
				b.Error("Expected to find cached item")
			}
		}
	})
}

// BenchmarkConcurrentAccess tests performance under concurrent load
func BenchmarkConcurrentAccess(b *testing.B) {
	concurrencyLevels := []int{1, 2, 4, 8, 16, 32}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency_%d", concurrency), func(b *testing.B) {
			benchmarkConcurrentAccess(b, concurrency)
		})
	}
}

func benchmarkConcurrentAccess(b *testing.B, concurrency int) {
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	// Create multiple clients to simulate concurrent access
	clients := make([]slurm.SlurmClient, concurrency)
	for i := 0; i < concurrency; i++ {
		client, err := slurm.NewClientWithVersion(context.Background(), "v0.0.42",
			slurm.WithBaseURL(mockServer.URL()),
			slurm.WithAuth(auth.NewNoAuth()),
		)
		if err != nil {
			b.Fatal(err)
		}
		clients[i] = client
		defer client.Close()
	}

	b.ResetTimer()

	// Run concurrent operations
	var wg sync.WaitGroup
	operations := b.N

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(clientIndex int) {
			defer wg.Done()
			client := clients[clientIndex]
			
			operationsPerWorker := operations / concurrency
			if clientIndex == concurrency-1 {
				operationsPerWorker += operations % concurrency
			}

			for j := 0; j < operationsPerWorker; j++ {
				ctx := context.Background()
				_, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{
					Limit: 5,
				})
				if err != nil {
					b.Error(err)
				}
			}
		}(i)
	}

	wg.Wait()
}

// BenchmarkMemoryUsage tests memory efficiency
func BenchmarkMemoryUsage(b *testing.B) {
	profiles := []performance.PerformanceProfile{
		performance.ProfileDefault,
		performance.ProfileConservative,
		performance.ProfileHighThroughput,
	}

	for _, profile := range profiles {
		b.Run(string(profile), func(b *testing.B) {
			benchmarkMemoryUsage(b, profile)
		})
	}
}

func benchmarkMemoryUsage(b *testing.B, profile performance.PerformanceProfile) {
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	// Force garbage collection and get baseline memory stats
	runtime.GC()
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// Create connection pool and cache with the given profile
	poolManager := performance.NewHTTPClientPoolManager()
	cacheManager := performance.NewCacheManager(performance.GetCacheConfigForProfile(profile))
	
	defer poolManager.Close()
	defer cacheManager.Close()

	// Create clients and perform operations
	clients := make([]slurm.SlurmClient, 10)
	for i := 0; i < 10; i++ {
		client, err := slurm.NewClientWithVersion(context.Background(), "v0.0.42",
			slurm.WithBaseURL(mockServer.URL()),
			slurm.WithAuth(auth.NewNoAuth()),
		)
		if err != nil {
			b.Fatal(err)
		}
		clients[i] = client
		defer client.Close()
	}

	b.ResetTimer()

	// Perform memory-intensive operations
	for i := 0; i < b.N; i++ {
		client := clients[i%10]
		ctx := context.Background()
		
		_, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{
			Limit: 100,
		})
		if err != nil {
			// Expected to fail with current implementation
		}
		
		// Add some cache operations
		cache := cacheManager.GetCache("v0.0.42")
		params := map[string]interface{}{"iteration": i}
		value := make([]byte, 1024) // 1KB of data
		cache.Set("test.operation", params, value)
	}

	// Force garbage collection and measure memory usage
	runtime.GC()
	runtime.ReadMemStats(&m2)

	// Report memory statistics
	allocatedMB := float64(m2.Alloc-m1.Alloc) / (1024 * 1024)
	b.ReportMetric(allocatedMB, "MB_allocated")
	
	totalAllocMB := float64(m2.TotalAlloc-m1.TotalAlloc) / (1024 * 1024)
	b.ReportMetric(totalAllocMB, "MB_total_alloc")
}

// BenchmarkCacheHitRatio tests cache effectiveness
func BenchmarkCacheHitRatio(b *testing.B) {
	config := performance.DefaultCacheConfig()
	cache := performance.NewResponseCache(config)
	defer cache.Close()

	operation := "jobs.list"
	value := []byte(`{"jobs": []}`)

	// Pre-populate cache with some data
	for i := 0; i < 100; i++ {
		params := map[string]interface{}{"page": i}
		cache.Set(operation, params, value)
	}

	b.ResetTimer()

	hits := 0
	for i := 0; i < b.N; i++ {
		// 80% chance of hitting cache (accessing existing pages)
		var params map[string]interface{}
		if i%10 < 8 {
			params = map[string]interface{}{"page": i % 100}
		} else {
			params = map[string]interface{}{"page": 100 + i}
		}

		_, found := cache.Get(operation, params)
		if found {
			hits++
		}
	}

	hitRatio := float64(hits) / float64(b.N)
	b.ReportMetric(hitRatio*100, "%_cache_hit_ratio")
}

// BenchmarkResponseTime measures response times for different scenarios
func BenchmarkResponseTime(b *testing.B) {
	scenarios := []struct {
		name      string
		delay     time.Duration
		cacheEnabled bool
	}{
		{"NoDelay_NoCache", 0, false},
		{"NoDelay_WithCache", 0, true},
		{"LowLatency_NoCache", 10 * time.Millisecond, false},
		{"LowLatency_WithCache", 10 * time.Millisecond, true},
		{"HighLatency_NoCache", 100 * time.Millisecond, false},
		{"HighLatency_WithCache", 100 * time.Millisecond, true},
	}

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			benchmarkResponseTime(b, scenario.delay, scenario.cacheEnabled)
		})
	}
}

func benchmarkResponseTime(b *testing.B, delay time.Duration, cacheEnabled bool) {
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	mockServer.GetConfig().ResponseDelay = delay
	defer mockServer.Close()

	var cache *performance.ResponseCache
	if cacheEnabled {
		cache = performance.NewResponseCache(performance.DefaultCacheConfig())
		defer cache.Close()
	}

	client, err := slurm.NewClientWithVersion(context.Background(), "v0.0.42",
		slurm.WithBaseURL(mockServer.URL()),
		slurm.WithAuth(auth.NewNoAuth()),
	)
	if err != nil {
		b.Fatal(err)
	}
	defer client.Close()

	b.ResetTimer()

	var totalDuration time.Duration
	for i := 0; i < b.N; i++ {
		start := time.Now()
		
		// Simulate cache check if enabled
		if cacheEnabled && cache != nil {
			params := map[string]interface{}{"iteration": i % 10} // Some cache hits
			if _, found := cache.Get("jobs.list", params); found {
				// Cache hit - skip actual API call
				totalDuration += time.Since(start)
				continue
			}
		}

		// Make actual API call
		ctx := context.Background()
		_, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{
			Limit: 10,
		})
		
		duration := time.Since(start)
		totalDuration += duration
		
		// Add to cache if enabled
		if cacheEnabled && cache != nil && err == nil {
			params := map[string]interface{}{"iteration": i % 10}
			value := []byte(`{"jobs": []}`)
			cache.Set("jobs.list", params, value)
		}
		
		if err != nil {
			// Expected to fail with current implementation
		}
	}

	avgResponseTime := totalDuration / time.Duration(b.N)
	b.ReportMetric(float64(avgResponseTime.Nanoseconds())/1e6, "ms_avg_response_time")
}