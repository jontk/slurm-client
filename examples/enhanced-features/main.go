// Package main demonstrates the enhanced features of the SLURM client
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	slurm "github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/internal/interfaces"
	slurmctx "github.com/jontk/slurm-client/pkg/context"
	"github.com/jontk/slurm-client/pkg/logging"
	"github.com/jontk/slurm-client/pkg/metrics"
	"github.com/jontk/slurm-client/pkg/middleware"
	"github.com/jontk/slurm-client/pkg/pool"
	"github.com/jontk/slurm-client/pkg/retry"
)

func main() {
	// Get configuration from environment
	baseURL := os.Getenv("SLURM_API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:6820"
	}
	
	token := os.Getenv("SLURM_API_TOKEN")
	if token == "" {
		log.Fatal("SLURM_API_TOKEN environment variable is required")
	}

	// Demonstrate all enhanced features
	fmt.Println("🚀 SLURM Client Enhanced Features Demo")
	fmt.Println("=====================================")
	
	// 1. Create a structured logger
	logConfig := &logging.Config{
		Level:   slog.LevelDebug,
		Format:  logging.FormatJSON,
		Output:  os.Stdout,
		Version: "1.0.0",
	}
	logger := logging.NewLogger(logConfig)
	
	fmt.Println("\n1️⃣ Structured Logging Enabled (JSON format)")
	logger.Info("Starting enhanced SLURM client demo", "base_url", baseURL)
	
	// 2. Create metrics collector
	metricsCollector := metrics.NewInMemoryCollector()
	fmt.Println("\n2️⃣ Metrics Collection Enabled")
	
	// 3. Setup timeout configuration
	timeoutConfig := &slurmctx.TimeoutConfig{
		Default: 30 * time.Second,
		Read:    20 * time.Second,
		Write:   40 * time.Second,
		List:    60 * time.Second,
		Watch:   0, // No timeout for watch operations
	}
	fmt.Println("\n3️⃣ Custom Timeout Configuration Set")
	fmt.Printf("   - Read: %v, Write: %v, List: %v\n", 
		timeoutConfig.Read, timeoutConfig.Write, timeoutConfig.List)
	
	// 4. Setup connection pooling
	poolConfig := &pool.PoolConfig{
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     20,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	fmt.Println("\n4️⃣ Connection Pooling Configured")
	fmt.Printf("   - Max idle connections: %d\n", poolConfig.MaxIdleConns)
	fmt.Printf("   - Max connections per host: %d\n", poolConfig.MaxConnsPerHost)
	
	// 5. Setup retry with exponential backoff
	retryBackoff := &retry.ExponentialBackoffStrategy{
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
		Jitter:       0.1,
		MaxAttempts:  5,
	}
	fmt.Println("\n5️⃣ Retry with Exponential Backoff Configured")
	fmt.Printf("   - Max attempts: %d, Initial delay: %v\n", 
		retryBackoff.MaxAttempts, retryBackoff.InitialDelay)
	
	// 6. Create custom middleware for request tracking
	requestTracker := middleware.WithHeaders(map[string]string{
		"X-Client-Name":    "enhanced-demo",
		"X-Client-Version": "1.0.0",
	})
	fmt.Println("\n6️⃣ Custom Middleware Added")
	
	// Create client with all enhanced features
	ctx := context.Background()
	
	fmt.Println("\n📦 Creating SLURM client with all enhancements...")
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL(baseURL),
		slurm.WithToken(token),
		slurm.WithLogger(logger),
		slurm.WithMetricsCollector(metricsCollector),
		slurm.WithTimeoutConfig(timeoutConfig),
		slurm.WithConnectionPool(poolConfig),
		slurm.WithRetryBackoff(retryBackoff),
		slurm.WithUserAgent("slurm-client-enhanced/1.0"),
		slurm.WithRequestID(func() string {
			return uuid.New().String()
		}),
		slurm.WithCircuitBreaker(5, 30*time.Second),
		slurm.WithMiddleware(requestTracker),
		slurm.WithDebug(),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	
	fmt.Printf("✅ Client created successfully (version: %s)\n", client.Version())
	
	// Demonstrate features in action
	fmt.Println("\n🎯 Testing Enhanced Features...")
	
	// Test 1: Simple ping with timeout and logging
	fmt.Println("\n📌 Test 1: Ping with timeout tracking")
	ctxWithTimeout, cancel := slurmctx.WithTimeout(ctx, slurmctx.OpRead, timeoutConfig)
	defer cancel()
	
	start := time.Now()
	err = client.Info().Ping(ctxWithTimeout)
	duration := time.Since(start)
	
	if err != nil {
		fmt.Printf("❌ Ping failed: %v (duration: %v)\n", err, duration)
	} else {
		fmt.Printf("✅ Ping successful (duration: %v)\n", duration)
	}
	
	// Test 2: List jobs with metrics tracking
	fmt.Println("\n📌 Test 2: List jobs with metrics")
	jobs, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{
		Limit: 10,
	})
	if err != nil {
		fmt.Printf("❌ Failed to list jobs: %v\n", err)
	} else {
		fmt.Printf("✅ Listed %d jobs\n", len(jobs.Jobs))
	}
	
	// Test 3: Node information with retry
	fmt.Println("\n📌 Test 3: Get node info with retry")
	nodes, err := client.Nodes().List(ctx, &interfaces.ListNodesOptions{
		Limit: 5,
	})
	if err != nil {
		fmt.Printf("❌ Failed to list nodes: %v\n", err)
	} else {
		fmt.Printf("✅ Listed %d nodes\n", len(nodes.Nodes))
	}
	
	// Display metrics
	fmt.Println("\n📊 Metrics Summary")
	fmt.Println("==================")
	stats := metricsCollector.GetStats()
	
	fmt.Printf("Total Requests: %d\n", stats.TotalRequests)
	fmt.Printf("Total Responses: %d\n", stats.TotalResponses)
	fmt.Printf("Total Errors: %d\n", stats.TotalErrors)
	fmt.Printf("Active Requests: %d\n", stats.ActiveRequests)
	
	if stats.ResponseTimeStats.Count > 0 {
		fmt.Printf("\nResponse Time Statistics:\n")
		fmt.Printf("  - Count: %d\n", stats.ResponseTimeStats.Count)
		fmt.Printf("  - Average: %v\n", stats.ResponseTimeStats.Average)
		fmt.Printf("  - Min: %v\n", stats.ResponseTimeStats.Min)
		fmt.Printf("  - Max: %v\n", stats.ResponseTimeStats.Max)
	}
	
	fmt.Println("\nRequests by Path:")
	for path, count := range stats.RequestsByPath {
		fmt.Printf("  - %s: %d\n", path, count)
	}
	
	if len(stats.ResponsesByStatus) > 0 {
		fmt.Println("\nResponses by Status:")
		for status, count := range stats.ResponsesByStatus {
			fmt.Printf("  - %d: %d\n", status, count)
		}
	}
	
	// Test context cancellation
	fmt.Println("\n📌 Test 4: Context cancellation")
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	
	// Start a goroutine that will be cancelled
	done := make(chan bool)
	go func() {
		fmt.Println("Starting long-running operation...")
		_, err := client.Jobs().List(cancelCtx, &interfaces.ListJobsOptions{
			Limit: 1000,
		})
		if err != nil {
			if slurmctx.IsContextError(err) {
				fmt.Println("✅ Operation cancelled as expected")
			} else {
				fmt.Printf("❌ Unexpected error: %v\n", err)
			}
		}
		done <- true
	}()
	
	// Cancel after a short delay
	time.Sleep(100 * time.Millisecond)
	cancelFunc()
	<-done
	
	// Final metrics
	fmt.Println("\n🏁 Final Metrics")
	fmt.Println("================")
	finalStats := metricsCollector.GetStats()
	fmt.Printf("Total operations: %d\n", finalStats.TotalRequests)
	fmt.Printf("Success rate: %.2f%%\n", 
		float64(finalStats.TotalResponses)/float64(finalStats.TotalRequests)*100)
	fmt.Printf("Runtime: %v\n", finalStats.Duration)
	
	fmt.Println("\n✨ Enhanced features demo completed!")
}