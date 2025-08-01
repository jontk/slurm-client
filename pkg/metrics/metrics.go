// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Package metrics provides metrics collection for the SLURM client
package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// Collector is the interface for metrics collection
type Collector interface {
	// RecordRequest records an API request
	RecordRequest(method, path string)
	
	// RecordResponse records an API response
	RecordResponse(method, path string, statusCode int, duration time.Duration)
	
	// RecordError records an API error
	RecordError(method, path string, err error)
	
	// RecordCacheHit records a cache hit
	RecordCacheHit(key string)
	
	// RecordCacheMiss records a cache miss
	RecordCacheMiss(key string)
	
	// GetStats returns current metrics statistics
	GetStats() *Stats
	
	// Reset resets all metrics
	Reset()
}

// Stats contains aggregated metrics statistics
type Stats struct {
	// Request metrics
	TotalRequests   int64
	ActiveRequests  int64
	RequestsByPath  map[string]int64
	
	// Response metrics
	TotalResponses      int64
	ResponsesByStatus   map[int]int64
	ResponseTimeStats   DurationStats
	ResponseTimeByPath  map[string]DurationStats
	
	// Error metrics
	TotalErrors    int64
	ErrorsByType   map[string]int64
	ErrorsByPath   map[string]int64
	
	// Cache metrics
	CacheHits   int64
	CacheMisses int64
	CacheRatio  float64
	
	// Timing
	StartTime time.Time
	Duration  time.Duration
}

// DurationStats contains statistics for duration measurements
type DurationStats struct {
	Count   int64
	Total   time.Duration
	Min     time.Duration
	Max     time.Duration
	Average time.Duration
}

// InMemoryCollector is an in-memory implementation of Collector
type InMemoryCollector struct {
	mu sync.RWMutex
	
	// Request counters
	totalRequests  int64
	activeRequests int64
	requestsByPath map[string]*int64
	
	// Response counters
	totalResponses    int64
	responsesByStatus map[int]*int64
	responseTimes     *durationAggregator
	responseTimeByPath map[string]*durationAggregator
	
	// Error counters
	totalErrors  int64
	errorsByType map[string]*int64
	errorsByPath map[string]*int64
	
	// Cache counters
	cacheHits   int64
	cacheMisses int64
	
	// Timing
	startTime time.Time
}

// NewInMemoryCollector creates a new in-memory metrics collector
func NewInMemoryCollector() *InMemoryCollector {
	return &InMemoryCollector{
		requestsByPath:     make(map[string]*int64),
		responsesByStatus:  make(map[int]*int64),
		responseTimes:      newDurationAggregator(),
		responseTimeByPath: make(map[string]*durationAggregator),
		errorsByType:       make(map[string]*int64),
		errorsByPath:       make(map[string]*int64),
		startTime:          time.Now(),
	}
}

// RecordRequest records an API request
func (c *InMemoryCollector) RecordRequest(method, path string) {
	atomic.AddInt64(&c.totalRequests, 1)
	atomic.AddInt64(&c.activeRequests, 1)
	
	key := method + " " + path
	incrementMapCounter(&c.mu, c.requestsByPath, key)
}

// RecordResponse records an API response
func (c *InMemoryCollector) RecordResponse(method, path string, statusCode int, duration time.Duration) {
	atomic.AddInt64(&c.totalResponses, 1)
	atomic.AddInt64(&c.activeRequests, -1)
	
	// Record status code
	incrementMapCounterInt(&c.mu, c.responsesByStatus, statusCode)
	
	// Record duration
	c.responseTimes.add(duration)
	
	// Record duration by path
	key := method + " " + path
	c.mu.Lock()
	agg, exists := c.responseTimeByPath[key]
	if !exists {
		agg = newDurationAggregator()
		c.responseTimeByPath[key] = agg
	}
	c.mu.Unlock()
	agg.add(duration)
}

// RecordError records an API error
func (c *InMemoryCollector) RecordError(method, path string, err error) {
	errorType := "unknown"
	if err != nil {
		errorType = err.Error()
	}
	atomic.AddInt64(&c.totalErrors, 1)
	atomic.AddInt64(&c.activeRequests, -1)
	
	incrementMapCounter(&c.mu, c.errorsByType, errorType)
	
	key := method + " " + path
	incrementMapCounter(&c.mu, c.errorsByPath, key)
}

// RecordCacheHit records a cache hit
func (c *InMemoryCollector) RecordCacheHit(key string) {
	atomic.AddInt64(&c.cacheHits, 1)
}

// RecordCacheMiss records a cache miss
func (c *InMemoryCollector) RecordCacheMiss(key string) {
	atomic.AddInt64(&c.cacheMisses, 1)
}

// GetStats returns current metrics statistics
func (c *InMemoryCollector) GetStats() *Stats {
	stats := &Stats{
		TotalRequests:      atomic.LoadInt64(&c.totalRequests),
		ActiveRequests:     atomic.LoadInt64(&c.activeRequests),
		TotalResponses:     atomic.LoadInt64(&c.totalResponses),
		TotalErrors:        atomic.LoadInt64(&c.totalErrors),
		CacheHits:          atomic.LoadInt64(&c.cacheHits),
		CacheMisses:        atomic.LoadInt64(&c.cacheMisses),
		RequestsByPath:     c.copyMapCounters(c.requestsByPath),
		ResponsesByStatus:  c.copyIntMapCounters(c.responsesByStatus),
		ErrorsByType:       c.copyMapCounters(c.errorsByType),
		ErrorsByPath:       c.copyMapCounters(c.errorsByPath),
		ResponseTimeStats:  c.responseTimes.stats(),
		ResponseTimeByPath: c.copyDurationStats(c.responseTimeByPath),
		StartTime:          c.startTime,
		Duration:           time.Since(c.startTime),
	}
	
	// Calculate cache ratio
	totalCache := stats.CacheHits + stats.CacheMisses
	if totalCache > 0 {
		stats.CacheRatio = float64(stats.CacheHits) / float64(totalCache)
	}
	
	return stats
}

// Reset resets all metrics
func (c *InMemoryCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Reset atomic counters
	atomic.StoreInt64(&c.totalRequests, 0)
	atomic.StoreInt64(&c.activeRequests, 0)
	atomic.StoreInt64(&c.totalResponses, 0)
	atomic.StoreInt64(&c.totalErrors, 0)
	atomic.StoreInt64(&c.cacheHits, 0)
	atomic.StoreInt64(&c.cacheMisses, 0)
	
	// Reset maps
	c.requestsByPath = make(map[string]*int64)
	c.responsesByStatus = make(map[int]*int64)
	c.responseTimes = newDurationAggregator()
	c.responseTimeByPath = make(map[string]*durationAggregator)
	c.errorsByType = make(map[string]*int64)
	c.errorsByPath = make(map[string]*int64)
	
	c.startTime = time.Now()
}

// incrementMapCounter safely increments a counter in a map
func incrementMapCounter(mu *sync.RWMutex, m map[string]*int64, key string) {
	mu.Lock()
	counter, exists := m[key]
	if !exists {
		var v int64
		counter = &v
		m[key] = counter
	}
	mu.Unlock()
	
	atomic.AddInt64(counter, 1)
}

// incrementMapCounterInt safely increments a counter in a map with int keys
func incrementMapCounterInt(mu *sync.RWMutex, m map[int]*int64, key int) {
	mu.Lock()
	counter, exists := m[key]
	if !exists {
		var v int64
		counter = &v
		m[key] = counter
	}
	mu.Unlock()
	
	atomic.AddInt64(counter, 1)
}

// copyMapCounters creates a copy of string map counters
func (c *InMemoryCollector) copyMapCounters(m map[string]*int64) map[string]int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	result := make(map[string]int64, len(m))
	for k, v := range m {
		result[k] = atomic.LoadInt64(v)
	}
	return result
}

// copyIntMapCounters creates a copy of int map counters
func (c *InMemoryCollector) copyIntMapCounters(m map[int]*int64) map[int]int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	result := make(map[int]int64, len(m))
	for k, v := range m {
		result[k] = atomic.LoadInt64(v)
	}
	return result
}

// copyDurationStats creates a copy of duration statistics
func (c *InMemoryCollector) copyDurationStats(m map[string]*durationAggregator) map[string]DurationStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	result := make(map[string]DurationStats, len(m))
	for k, v := range m {
		result[k] = v.stats()
	}
	return result
}

// durationAggregator aggregates duration statistics
type durationAggregator struct {
	mu      sync.Mutex
	count   int64
	total   time.Duration
	min     time.Duration
	max     time.Duration
}

func newDurationAggregator() *durationAggregator {
	return &durationAggregator{
		min: time.Duration(1<<63 - 1), // MaxInt64
	}
}

func (d *durationAggregator) add(duration time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	d.count++
	d.total += duration
	
	if duration < d.min {
		d.min = duration
	}
	if duration > d.max {
		d.max = duration
	}
}

func (d *durationAggregator) stats() DurationStats {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	stats := DurationStats{
		Count: d.count,
		Total: d.total,
		Min:   d.min,
		Max:   d.max,
	}
	
	if d.count > 0 {
		stats.Average = time.Duration(int64(d.total) / d.count)
	}
	
	// Reset min if no data
	if d.count == 0 {
		stats.Min = 0
	}
	
	return stats
}

// NoOpCollector is a no-op implementation of Collector
type NoOpCollector struct{}

func (NoOpCollector) RecordRequest(method, path string) {}
func (NoOpCollector) RecordResponse(method, path string, statusCode int, duration time.Duration) {}
func (NoOpCollector) RecordError(method, path string, err error) {}
func (NoOpCollector) RecordCacheHit(key string) {}
func (NoOpCollector) RecordCacheMiss(key string) {}
func (NoOpCollector) GetStats() *Stats { return &Stats{} }
func (NoOpCollector) Reset() {}

// Global default collector
var defaultCollector Collector = &NoOpCollector{}

// SetDefaultCollector sets the default metrics collector
func SetDefaultCollector(collector Collector) {
	if collector == nil {
		collector = &NoOpCollector{}
	}
	defaultCollector = collector
}

// GetDefaultCollector returns the default metrics collector
func GetDefaultCollector() Collector {
	return defaultCollector
}
