// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInMemoryCollector(t *testing.T) {
	collector := NewInMemoryCollector()

	require.NotNil(t, collector)
	assert.NotNil(t, collector.requestsByPath)
	assert.NotNil(t, collector.responsesByStatus)
	assert.NotNil(t, collector.responseTimes)
	assert.NotNil(t, collector.responseTimeByPath)
	assert.NotNil(t, collector.errorsByType)
	assert.NotNil(t, collector.errorsByPath)
	assert.False(t, collector.startTime.IsZero())
}

func TestInMemoryCollector_RecordRequest(t *testing.T) {
	collector := NewInMemoryCollector()

	collector.RecordRequest("GET", "/api/v1/jobs")
	collector.RecordRequest("POST", "/api/v1/jobs")
	collector.RecordRequest("GET", "/api/v1/jobs") // duplicate

	stats := collector.GetStats()
	assert.Equal(t, int64(3), stats.TotalRequests)
	assert.Equal(t, int64(3), stats.ActiveRequests)
	assert.Equal(t, int64(2), stats.RequestsByPath["GET /api/v1/jobs"])
	assert.Equal(t, int64(1), stats.RequestsByPath["POST /api/v1/jobs"])
}

func TestInMemoryCollector_RecordResponse(t *testing.T) {
	collector := NewInMemoryCollector()

	// Record some requests first
	collector.RecordRequest("GET", "/api/v1/jobs")
	collector.RecordRequest("POST", "/api/v1/jobs")

	collector.RecordResponse("GET", "/api/v1/jobs", 200, 100*time.Millisecond)
	collector.RecordResponse("POST", "/api/v1/jobs", 201, 200*time.Millisecond)

	stats := collector.GetStats()
	assert.Equal(t, int64(2), stats.TotalResponses)
	assert.Equal(t, int64(0), stats.ActiveRequests) // Both completed
	assert.Equal(t, int64(1), stats.ResponsesByStatus[200])
	assert.Equal(t, int64(1), stats.ResponsesByStatus[201])

	// Check overall response time stats
	assert.Equal(t, int64(2), stats.ResponseTimeStats.Count)
	assert.Equal(t, 300*time.Millisecond, stats.ResponseTimeStats.Total)
	assert.Equal(t, 100*time.Millisecond, stats.ResponseTimeStats.Min)
	assert.Equal(t, 200*time.Millisecond, stats.ResponseTimeStats.Max)
	assert.Equal(t, 150*time.Millisecond, stats.ResponseTimeStats.Average)

	// Check per-path response time stats
	getStats := stats.ResponseTimeByPath["GET /api/v1/jobs"]
	assert.Equal(t, int64(1), getStats.Count)
	assert.Equal(t, 100*time.Millisecond, getStats.Total)
	assert.Equal(t, 100*time.Millisecond, getStats.Average)

	postStats := stats.ResponseTimeByPath["POST /api/v1/jobs"]
	assert.Equal(t, int64(1), postStats.Count)
	assert.Equal(t, 200*time.Millisecond, postStats.Total)
	assert.Equal(t, 200*time.Millisecond, postStats.Average)
}

func TestInMemoryCollector_RecordError(t *testing.T) {
	collector := NewInMemoryCollector()

	// Record some requests first
	collector.RecordRequest("GET", "/api/v1/jobs")
	collector.RecordRequest("POST", "/api/v1/jobs")

	err1 := errors.New("connection timeout")
	err2 := errors.New("unauthorized")

	collector.RecordError("GET", "/api/v1/jobs", err1)
	collector.RecordError("POST", "/api/v1/jobs", err2)
	collector.RecordError("GET", "/api/v1/jobs", err1) // duplicate error type

	stats := collector.GetStats()
	assert.Equal(t, int64(3), stats.TotalErrors)
	assert.Equal(t, int64(-1), stats.ActiveRequests) // One extra error recorded
	assert.Equal(t, int64(2), stats.ErrorsByType["connection timeout"])
	assert.Equal(t, int64(1), stats.ErrorsByType["unauthorized"])
	assert.Equal(t, int64(2), stats.ErrorsByPath["GET /api/v1/jobs"])
	assert.Equal(t, int64(1), stats.ErrorsByPath["POST /api/v1/jobs"])
}

func TestInMemoryCollector_RecordErrorWithNil(t *testing.T) {
	collector := NewInMemoryCollector()

	collector.RecordRequest("GET", "/api/v1/jobs")
	collector.RecordError("GET", "/api/v1/jobs", nil)

	stats := collector.GetStats()
	assert.Equal(t, int64(1), stats.TotalErrors)
	assert.Equal(t, int64(1), stats.ErrorsByType["unknown"])
}

func TestInMemoryCollector_RecordCache(t *testing.T) {
	collector := NewInMemoryCollector()

	collector.RecordCacheHit("user:123")
	collector.RecordCacheHit("job:456")
	collector.RecordCacheMiss("user:789")
	collector.RecordCacheHit("user:123") // duplicate hit

	stats := collector.GetStats()
	assert.Equal(t, int64(3), stats.CacheHits)
	assert.Equal(t, int64(1), stats.CacheMisses)
	assert.Equal(t, 0.75, stats.CacheRatio) // 3/(3+1) = 0.75
}

func TestInMemoryCollector_Reset(t *testing.T) {
	collector := NewInMemoryCollector()

	// Add some data
	collector.RecordRequest("GET", "/api/v1/jobs")
	collector.RecordResponse("GET", "/api/v1/jobs", 200, 100*time.Millisecond)
	collector.RecordError("POST", "/api/v1/jobs", errors.New("test error"))
	collector.RecordCacheHit("test:key")
	collector.RecordCacheMiss("test:key2")

	// Verify data exists
	stats := collector.GetStats()
	assert.Positive(t, stats.TotalRequests)
	assert.Positive(t, stats.TotalResponses)
	assert.Positive(t, stats.TotalErrors)
	assert.Positive(t, stats.CacheHits)
	assert.Positive(t, stats.CacheMisses)

	// Reset
	collector.Reset()

	// Verify everything is reset
	stats = collector.GetStats()
	assert.Equal(t, int64(0), stats.TotalRequests)
	assert.Equal(t, int64(0), stats.ActiveRequests)
	assert.Equal(t, int64(0), stats.TotalResponses)
	assert.Equal(t, int64(0), stats.TotalErrors)
	assert.Equal(t, int64(0), stats.CacheHits)
	assert.Equal(t, int64(0), stats.CacheMisses)
	assert.Equal(t, 0.0, stats.CacheRatio)
	assert.Empty(t, stats.RequestsByPath)
	assert.Empty(t, stats.ResponsesByStatus)
	assert.Empty(t, stats.ErrorsByType)
	assert.Empty(t, stats.ErrorsByPath)
	assert.Empty(t, stats.ResponseTimeByPath)
	assert.Equal(t, int64(0), stats.ResponseTimeStats.Count)
}

func TestStats_CacheRatioCalculation(t *testing.T) {
	collector := NewInMemoryCollector()

	t.Run("no cache operations", func(t *testing.T) {
		stats := collector.GetStats()
		assert.Equal(t, 0.0, stats.CacheRatio)
	})

	t.Run("only hits", func(t *testing.T) {
		collector.Reset()
		collector.RecordCacheHit("key1")
		collector.RecordCacheHit("key2")

		stats := collector.GetStats()
		assert.Equal(t, 1.0, stats.CacheRatio)
	})

	t.Run("only misses", func(t *testing.T) {
		collector.Reset()
		collector.RecordCacheMiss("key1")
		collector.RecordCacheMiss("key2")

		stats := collector.GetStats()
		assert.Equal(t, 0.0, stats.CacheRatio)
	})

	t.Run("mixed hits and misses", func(t *testing.T) {
		collector.Reset()
		collector.RecordCacheHit("key1")
		collector.RecordCacheMiss("key2")
		collector.RecordCacheMiss("key3")

		stats := collector.GetStats()
		assert.Equal(t, 1.0/3.0, stats.CacheRatio)
	})
}

func TestDurationAggregator(t *testing.T) {
	agg := newDurationAggregator()

	t.Run("initial state", func(t *testing.T) {
		stats := agg.stats()
		assert.Equal(t, int64(0), stats.Count)
		assert.Equal(t, time.Duration(0), stats.Total)
		assert.Equal(t, time.Duration(0), stats.Min)
		assert.Equal(t, time.Duration(0), stats.Max)
		assert.Equal(t, time.Duration(0), stats.Average)
	})

	t.Run("single value", func(t *testing.T) {
		agg.add(100 * time.Millisecond)

		stats := agg.stats()
		assert.Equal(t, int64(1), stats.Count)
		assert.Equal(t, 100*time.Millisecond, stats.Total)
		assert.Equal(t, 100*time.Millisecond, stats.Min)
		assert.Equal(t, 100*time.Millisecond, stats.Max)
		assert.Equal(t, 100*time.Millisecond, stats.Average)
	})

	t.Run("multiple values", func(t *testing.T) {
		agg.add(200 * time.Millisecond)
		agg.add(50 * time.Millisecond)

		stats := agg.stats()
		assert.Equal(t, int64(3), stats.Count)
		assert.Equal(t, 350*time.Millisecond, stats.Total)
		assert.Equal(t, 50*time.Millisecond, stats.Min)
		assert.Equal(t, 200*time.Millisecond, stats.Max)
		// 350/3 = 116.666... which gets truncated to 116.666666ms due to duration precision
		expected := time.Duration(350000000 / 3) // 116.666666ms
		assert.Equal(t, expected, stats.Average)
	})
}

func TestDurationAggregator_Concurrency(t *testing.T) {
	agg := newDurationAggregator()

	const numGoroutines = 10
	const numOperations = 100

	var wg sync.WaitGroup

	// Add values concurrently
	for i := range numGoroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range numOperations {
				agg.add(time.Duration(id*numOperations+j) * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()

	stats := agg.stats()
	assert.Equal(t, int64(numGoroutines*numOperations), stats.Count)
	assert.Greater(t, stats.Total, time.Duration(0))
	assert.Greater(t, stats.Max, stats.Min)
	assert.Greater(t, stats.Average, time.Duration(0))
}

func TestInMemoryCollector_Concurrency(t *testing.T) {
	collector := NewInMemoryCollector()

	const numGoroutines = 10
	const numOperations = 100

	var wg sync.WaitGroup

	// Test concurrent operations
	for i := range numGoroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range numOperations {
				collector.RecordRequest("GET", "/api/test")
				collector.RecordResponse("GET", "/api/test", 200, time.Duration(j)*time.Millisecond)
				if j%10 == 0 {
					collector.RecordError("POST", "/api/test", errors.New("test error"))
				}
				collector.RecordCacheHit("key")
				collector.RecordCacheMiss("other-key")
			}
		}(i)
	}

	wg.Wait()

	stats := collector.GetStats()
	assert.Equal(t, int64(numGoroutines*numOperations), stats.TotalRequests)
	assert.Equal(t, int64(numGoroutines*numOperations), stats.TotalResponses)
	assert.Equal(t, int64(numGoroutines*10), stats.TotalErrors) // Every 10th operation
	assert.Equal(t, int64(numGoroutines*numOperations), stats.CacheHits)
	assert.Equal(t, int64(numGoroutines*numOperations), stats.CacheMisses)
}

func TestNoOpCollector(t *testing.T) {
	collector := NoOpCollector{}

	// All methods should not panic
	collector.RecordRequest("GET", "/api/test")
	collector.RecordResponse("GET", "/api/test", 200, 100*time.Millisecond)
	collector.RecordError("GET", "/api/test", errors.New("test error"))
	collector.RecordCacheHit("key")
	collector.RecordCacheMiss("key")

	stats := collector.GetStats()
	require.NotNil(t, stats)

	// Should return empty stats
	assert.Equal(t, int64(0), stats.TotalRequests)
	assert.Equal(t, int64(0), stats.TotalResponses)
	assert.Equal(t, int64(0), stats.TotalErrors)
	assert.Equal(t, int64(0), stats.CacheHits)
	assert.Equal(t, int64(0), stats.CacheMisses)

	// Reset should not panic
	collector.Reset()
}

func TestDefaultCollector(t *testing.T) {
	// Should start with NoOpCollector
	defaultCol := GetDefaultCollector()
	assert.IsType(t, &NoOpCollector{}, defaultCol)

	// Set a new collector
	newCollector := NewInMemoryCollector()
	SetDefaultCollector(newCollector)

	assert.Equal(t, newCollector, GetDefaultCollector())

	// Set nil collector (should default to NoOpCollector)
	SetDefaultCollector(nil)
	assert.IsType(t, &NoOpCollector{}, GetDefaultCollector())

	// Restore original
	SetDefaultCollector(&NoOpCollector{})
}

func TestCollectorInterface(t *testing.T) {
	// Verify that InMemoryCollector implements Collector interface
	var _ Collector = (*InMemoryCollector)(nil)

	// Verify that NoOpCollector implements Collector interface
	var _ Collector = NoOpCollector{}
}

func TestStatsStructure(t *testing.T) {
	collector := NewInMemoryCollector()

	// Add some varied data
	collector.RecordRequest("GET", "/api/jobs")
	collector.RecordRequest("POST", "/api/jobs")
	collector.RecordResponse("GET", "/api/jobs", 200, 50*time.Millisecond)
	collector.RecordResponse("POST", "/api/jobs", 201, 150*time.Millisecond)
	collector.RecordError("DELETE", "/api/jobs", errors.New("not found"))
	collector.RecordCacheHit("job:123")
	collector.RecordCacheMiss("job:456")

	stats := collector.GetStats()

	// Verify all fields are populated correctly
	assert.NotZero(t, stats.TotalRequests)
	assert.NotZero(t, stats.TotalResponses)
	assert.NotZero(t, stats.TotalErrors)
	assert.NotZero(t, stats.CacheHits)
	assert.NotZero(t, stats.CacheMisses)
	assert.NotZero(t, stats.CacheRatio)
	assert.NotEmpty(t, stats.RequestsByPath)
	assert.NotEmpty(t, stats.ResponsesByStatus)
	assert.NotEmpty(t, stats.ErrorsByType)
	assert.NotEmpty(t, stats.ErrorsByPath)
	assert.NotEmpty(t, stats.ResponseTimeByPath)
	assert.NotZero(t, stats.ResponseTimeStats.Count)
	assert.False(t, stats.StartTime.IsZero())
	assert.GreaterOrEqual(t, stats.Duration, time.Duration(0)) // May be 0 on very fast systems
}

func TestIncrementMapCounter(t *testing.T) {
	var mu sync.RWMutex
	m := make(map[string]*int64)

	// Test creating new counter
	incrementMapCounter(&mu, m, "test-key")

	mu.RLock()
	counter, exists := m["test-key"]
	mu.RUnlock()

	assert.True(t, exists)
	assert.Equal(t, int64(1), *counter)

	// Test incrementing existing counter
	incrementMapCounter(&mu, m, "test-key")

	mu.RLock()
	counter = m["test-key"]
	mu.RUnlock()

	assert.Equal(t, int64(2), *counter)
}

func TestIncrementMapCounterInt(t *testing.T) {
	var mu sync.RWMutex
	m := make(map[int]*int64)

	// Test creating new counter
	incrementMapCounterInt(&mu, m, 200)

	mu.RLock()
	counter, exists := m[200]
	mu.RUnlock()

	assert.True(t, exists)
	assert.Equal(t, int64(1), *counter)

	// Test incrementing existing counter
	incrementMapCounterInt(&mu, m, 200)

	mu.RLock()
	counter = m[200]
	mu.RUnlock()

	assert.Equal(t, int64(2), *counter)
}
