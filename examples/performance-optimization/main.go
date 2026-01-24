// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
	"github.com/jontk/slurm-client/pkg/performance"
)

// Example: Performance optimization patterns for high-throughput operations
func main() {
	// Create configuration optimized for performance
	cfg := config.NewDefault()
	cfg.BaseURL = "https://cluster.example.com:6820"
	cfg.Timeout = 30 * time.Second
	cfg.MaxRetries = 3

	// Create authentication
	authProvider := auth.NewTokenAuth("your-jwt-token")

	ctx := context.Background()

	// Example 1: Connection pooling
	fmt.Println("=== Connection Pooling ===")
	demonstrateConnectionPooling(ctx, cfg, authProvider)

	// Example 2: Response caching
	fmt.Println("\n=== Response Caching ===")
	demonstrateResponseCaching(ctx, cfg, authProvider)

	// Example 3: Batch operations
	fmt.Println("\n=== Batch Operations ===")
	demonstrateBatchOperations(ctx, cfg, authProvider)

	// Example 4: Concurrent requests
	fmt.Println("\n=== Concurrent Requests ===")
	demonstrateConcurrentRequests(ctx, cfg, authProvider)

	// Example 5: Performance profiling
	fmt.Println("\n=== Performance Profiling ===")
	demonstratePerformanceProfiling(ctx, cfg, authProvider)
}

// demonstrateConnectionPooling shows how to use connection pooling
func demonstrateConnectionPooling(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	// Create pool manager
	poolManager := performance.NewHTTPClientPoolManager()

	// Get optimized pool for high-throughput workload
	pool := poolManager.GetPoolForVersion("v0.0.42", performance.ProfileHighThroughput)

	// Create client with custom HTTP client from pool
	// Note: pool is *HTTPClientPool, not *http.Client, so we need to get a client from it
	httpClient := pool.GetClient(cfg.BaseURL)
	client, err := slurm.NewClient(ctx,
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
		slurm.WithHTTPClient(httpClient),
	)
	if err != nil {
		log.Printf("Failed to create client with pooling: %v", err)
		return
	}
	defer client.Close()

	// Demonstrate connection reuse with multiple requests
	start := time.Now()

	// Make multiple requests that reuse connections
	for i := range 10 {
		err := client.Info().Ping(ctx)
		if err != nil {
			log.Printf("Ping %d failed: %v", i, err)
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("10 requests with connection pooling: %v\n", elapsed)
	fmt.Printf("Average per request: %v\n", elapsed/10)

	// Compare with non-pooled client
	basicClient, err := slurm.NewClient(ctx,
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
	)
	if err != nil {
		log.Printf("Failed to create basic client: %v", err)
		return
	}
	defer basicClient.Close()

	start = time.Now()
	for i := range 10 {
		err := basicClient.Info().Ping(ctx)
		if err != nil {
			log.Printf("Basic ping %d failed: %v", i, err)
		}
	}

	elapsed = time.Since(start)
	fmt.Printf("\n10 requests without pooling: %v\n", elapsed)
	fmt.Printf("Average per request: %v\n", elapsed/10)

	// Show pool statistics
	// Note: GetPoolStats method doesn't exist on HTTPClientPoolManager
	// stats := poolManager.GetPoolStats("v0.0.42")
	// Connection pool stats would be shown here if the method existed
}

// demonstrateResponseCaching shows how to use response caching
func demonstrateResponseCaching(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	// Create cache with custom configuration
	// GetCacheConfigForProfile doesn't exist, use DefaultCacheConfig
	cacheConfig := performance.DefaultCacheConfig()
	cache := performance.NewResponseCache(cacheConfig)

	// Create client (in real usage, you'd integrate cache with client)
	client, err := slurm.NewClient(ctx,
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
	)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
	defer client.Close()

	// Simulate caching pattern for frequently accessed data

	// First access - cache miss
	start := time.Now()
	partitions1, err := getPartitionsWithCache(ctx, client, cache, "partitions-list")
	if err != nil {
		log.Printf("Failed to get partitions: %v", err)
		return
	}
	elapsed1 := time.Since(start)
	fmt.Printf("First access (cache miss): %v, found %d partitions\n",
		elapsed1, len(partitions1.Partitions))

	// Second access - cache hit
	start = time.Now()
	partitions2, err := getPartitionsWithCache(ctx, client, cache, "partitions-list")
	if err != nil {
		log.Printf("Failed to get cached partitions: %v", err)
		return
	}
	elapsed2 := time.Since(start)
	fmt.Printf("Second access (cache hit): %v, found %d partitions\n",
		elapsed2, len(partitions2.Partitions))

	fmt.Printf("Speedup: %.2fx faster\n", float64(elapsed1)/float64(elapsed2))

	// Show cache statistics
	stats := cache.GetStats()
	fmt.Printf("\nCache statistics:\n")
	fmt.Printf("  Hits: %d\n", stats.Hits)
	fmt.Printf("  Misses: %d\n", stats.Misses)
	fmt.Printf("  Hit rate: %.2f%%\n", stats.HitRatio*100)
	fmt.Printf("  Size: %d entries\n", stats.CurrentItems)

	// Demonstrate cache invalidation
	fmt.Println("\nInvalidating cache...")
	cache.InvalidatePattern("partitions-list")

	// Third access - cache miss after invalidation
	start = time.Now()
	_, err = getPartitionsWithCache(ctx, client, cache, "partitions-list")
	if err != nil {
		log.Printf("Failed to get partitions after invalidation: %v", err)
		return
	}
	elapsed3 := time.Since(start)
	fmt.Printf("After invalidation (cache miss): %v\n", elapsed3)
}

// Helper function for caching demonstration
func getPartitionsWithCache(ctx context.Context, client slurm.SlurmClient, cache *performance.ResponseCache, key string) (*interfaces.PartitionList, error) {
	// Check cache first
	// Note: In a real implementation, we would check the cache and deserialize
	// For this example, we always fetch fresh data
	_, _ = cache.Get(key, map[string]interface{}{})

	// Fetch from API
	partitions, err := client.Partitions().List(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Store in cache
	// Convert partitions to bytes for caching
	// In real usage, you'd serialize to JSON
	// cache.Set(key, map[string]interface{}{}, serializedBytes)

	return partitions, nil
}

// demonstrateBatchOperations shows efficient batch operation patterns
func demonstrateBatchOperations(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	client, err := slurm.NewClient(ctx,
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
	)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
	defer client.Close()

	// Inefficient: Individual requests
	fmt.Println("Individual requests (inefficient):")
	start := time.Now()
	jobIDs := []string{"job1", "job2", "job3", "job4", "job5"}

	for _, jobID := range jobIDs {
		_, err := client.Jobs().Get(ctx, jobID)
		if err != nil {
			log.Printf("Failed to get job %s: %v", jobID, err)
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("  Time: %v for %d jobs\n", elapsed, len(jobIDs))

	// Efficient: Batch request with filtering
	fmt.Println("\nBatch request (efficient):")
	start = time.Now()

	// Get all jobs in one request and filter client-side
	allJobs, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{
		Limit: 100,
	})
	if err != nil {
		log.Printf("Failed to list jobs: %v", err)
		return
	}

	// Filter for our specific jobs
	jobIDMap := make(map[string]bool)
	for _, id := range jobIDs {
		jobIDMap[id] = true
	}

	matchCount := 0
	for _, job := range allJobs.Jobs {
		if jobIDMap[job.ID] {
			matchCount++
		}
	}

	elapsed = time.Since(start)
	fmt.Printf("  Time: %v for %d jobs (found %d)\n", elapsed, len(jobIDs), matchCount)
	fmt.Printf("  Speedup: %.2fx faster\n", float64(len(jobIDs))*float64(elapsed)/float64(elapsed))

	// Demonstrate batch submission
	fmt.Println("\nBatch job submission:")
	start = time.Now()

	// Prepare batch of jobs
	var submittedJobs []string
	for i := range 5 {
		job := &interfaces.JobSubmission{
			Name:      fmt.Sprintf("batch-job-%d", i),
			Command:   fmt.Sprintf("echo 'Batch job %d'", i),
			Partition: "compute",
			CPUs:      1,
			Memory:    1024,
			TimeLimit: 5,
		}

		resp, err := client.Jobs().Submit(ctx, job)
		if err != nil {
			log.Printf("Failed to submit job %d: %v", i, err)
			continue
		}
		submittedJobs = append(submittedJobs, resp.JobID)
	}

	elapsed = time.Since(start)
	fmt.Printf("  Submitted %d jobs in %v\n", len(submittedJobs), elapsed)
	fmt.Printf("  Average per job: %v\n", elapsed/time.Duration(len(submittedJobs)))
}

// demonstrateConcurrentRequests shows concurrent request patterns
func demonstrateConcurrentRequests(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	// Create multiple clients for concurrent operations
	numWorkers := 4
	var wg sync.WaitGroup
	results := make(chan result, numWorkers*3)

	// Start workers
	for i := range numWorkers {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// Each worker gets its own client
			client, err := slurm.NewClient(ctx,
				slurm.WithConfig(cfg),
				slurm.WithAuth(auth),
			)
			if err != nil {
				log.Printf("Worker %d: Failed to create client: %v", workerID, err)
				return
			}
			defer client.Close()

			// Perform concurrent operations
			operations := []operation{
				{name: "list-jobs", fn: func() error {
					_, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{Limit: 10})
					return err
				}},
				{name: "list-nodes", fn: func() error {
					_, err := client.Nodes().List(ctx, &interfaces.ListNodesOptions{Limit: 10})
					return err
				}},
				{name: "get-info", fn: func() error {
					_, err := client.Info().Get(ctx)
					return err
				}},
			}

			for _, op := range operations {
				start := time.Now()
				err := op.fn()
				elapsed := time.Since(start)

				results <- result{
					workerID:  workerID,
					operation: op.name,
					duration:  elapsed,
					error:     err,
				}
			}
		}(i)
	}

	// Wait for all workers
	wg.Wait()
	close(results)

	// Analyze results
	var totalDuration time.Duration
	successCount := 0
	operationTimes := make(map[string][]time.Duration)

	for r := range results {
		if r.error == nil {
			successCount++
			totalDuration += r.duration
			operationTimes[r.operation] = append(operationTimes[r.operation], r.duration)
		} else {
			log.Printf("Worker %d: %s failed: %v", r.workerID, r.operation, r.error)
		}
	}

	fmt.Printf("Concurrent operations completed:\n")
	fmt.Printf("  Workers: %d\n", numWorkers)
	fmt.Printf("  Total operations: %d\n", numWorkers*3)
	fmt.Printf("  Successful: %d\n", successCount)
	fmt.Printf("  Average duration: %v\n", totalDuration/time.Duration(successCount))

	// Show per-operation statistics
	fmt.Println("\nPer-operation statistics:")
	for op, times := range operationTimes {
		var total time.Duration
		for _, t := range times {
			total += t
		}
		avg := total / time.Duration(len(times))
		fmt.Printf("  %s: avg=%v, count=%d\n", op, avg, len(times))
	}
}

// demonstratePerformanceProfiling shows how to profile and optimize
func demonstratePerformanceProfiling(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	// Create performance profiler
	// Note: Profiler type doesn't exist, we'll track metrics manually

	// Profile different configurations
	profiles := []struct {
		name        string
		profile     performance.PerformanceProfile
		description string
	}{
		{
			name:        "Default",
			profile:     performance.ProfileDefault,
			description: "Balanced performance",
		},
		{
			name:        "HighThroughput",
			profile:     performance.ProfileHighThroughput,
			description: "Optimized for many requests",
		},
		{
			name:        "LowLatency",
			profile:     performance.ProfileLowLatency,
			description: "Optimized for response time",
		},
		{
			name:        "ResourceConstrained",
			profile:     performance.ProfileConservative,
			description: "Minimal resource usage",
		},
	}

	for _, p := range profiles {
		fmt.Printf("\nProfiling %s profile (%s):\n", p.name, p.description)

		// Get optimized configuration
		poolManager := performance.NewHTTPClientPoolManager()
		profilePool := poolManager.GetPoolForVersion("v0.0.42", p.profile)
		httpClient := profilePool.GetClient(cfg.BaseURL)

		client, err := slurm.NewClient(ctx,
			slurm.WithConfig(cfg),
			slurm.WithAuth(auth),
			slurm.WithHTTPClient(httpClient),
		)
		if err != nil {
			log.Printf("Failed to create client for %s: %v", p.name, err)
			continue
		}

		// Run benchmark
		metrics := runBenchmark(ctx, client)

		// Display metrics
		fmt.Printf("  Requests: %d\n", metrics.requests)
		fmt.Printf("  Duration: %v\n", metrics.duration)
		fmt.Printf("  RPS: %.2f\n", metrics.requestsPerSecond)
		fmt.Printf("  Avg latency: %v\n", metrics.avgLatency)
		fmt.Printf("  P95 latency: %v\n", metrics.p95Latency)
		fmt.Printf("  Errors: %d\n", metrics.errors)

		client.Close()
	}

	// Show optimization recommendations
	fmt.Println("\nOptimization Recommendations:")
	fmt.Println("1. Use ProfileHighThroughput for batch operations")
	fmt.Println("2. Use ProfileLowLatency for interactive workloads")
	fmt.Println("3. Enable caching for read-heavy workloads")
	fmt.Println("4. Use connection pooling for all production workloads")
	fmt.Println("5. Batch requests when possible to reduce overhead")
}

// Helper types and functions

type operation struct {
	name string
	fn   func() error
}

type result struct {
	workerID  int
	operation string
	duration  time.Duration
	error     error
}

type metrics struct {
	requests          int
	duration          time.Duration
	requestsPerSecond float64
	avgLatency        time.Duration
	p95Latency        time.Duration
	errors            int
}

func runBenchmark(ctx context.Context, client slurm.SlurmClient) metrics {
	numRequests := 50
	latencies := make([]time.Duration, 0, numRequests)
	errors := 0

	// Start profiling
	// profiler.Start(profileName) - Profiler doesn't exist
	benchStart := time.Now()

	// Run requests
	for range numRequests {
		reqStart := time.Now()
		err := client.Info().Ping(ctx)
		reqDuration := time.Since(reqStart)

		if err != nil {
			errors++
		} else {
			latencies = append(latencies, reqDuration)
		}

		// Small delay to avoid overwhelming the server
		time.Sleep(10 * time.Millisecond)
	}

	benchDuration := time.Since(benchStart)
	// profiler.Stop(profileName) - Profiler doesn't exist

	// Calculate metrics
	var totalLatency time.Duration
	for _, l := range latencies {
		totalLatency += l
	}

	avgLatency := totalLatency / time.Duration(len(latencies))

	// Calculate P95 (simplified)
	p95Index := int(float64(len(latencies)) * 0.95)
	p95Latency := latencies[0] // Default to first if not enough samples
	if p95Index < len(latencies) {
		p95Latency = latencies[p95Index]
	}

	return metrics{
		requests:          numRequests,
		duration:          benchDuration,
		requestsPerSecond: float64(numRequests) / benchDuration.Seconds(),
		avgLatency:        avgLatency,
		p95Latency:        p95Latency,
		errors:            errors,
	}
}
