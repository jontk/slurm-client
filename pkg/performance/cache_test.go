// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package performance

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultCacheConfig(t *testing.T) {
	config := DefaultCacheConfig()

	require.NotNil(t, config)
	assert.Equal(t, 5*time.Minute, config.DefaultTTL)
	assert.Equal(t, 1000, config.MaxSize)
	assert.True(t, config.EnableCompression)
	assert.Equal(t, 1*time.Minute, config.CleanupInterval)
	assert.NotEmpty(t, config.TTLByOperation)

	// Test specific operation TTLs
	assert.Equal(t, 30*time.Minute, config.TTLByOperation["info.version"])
	assert.Equal(t, 30*time.Second, config.TTLByOperation["info.stats"])
	assert.Equal(t, 30*time.Second, config.TTLByOperation["jobs.list"])
}

func TestAggressiveCacheConfig(t *testing.T) {
	config := AggressiveCacheConfig()

	require.NotNil(t, config)
	assert.Equal(t, 10*time.Minute, config.DefaultTTL)
	assert.Equal(t, 5000, config.MaxSize)
	assert.Equal(t, 2*time.Minute, config.TTLByOperation["info.stats"])
	assert.Equal(t, 5*time.Minute, config.TTLByOperation["nodes.list"])
	assert.Equal(t, 2*time.Minute, config.TTLByOperation["jobs.list"])
	assert.Equal(t, 5*time.Minute, config.TTLByOperation["jobs.get"])
}

func TestConservativeCacheConfig(t *testing.T) {
	config := ConservativeCacheConfig()

	require.NotNil(t, config)
	assert.Equal(t, 1*time.Minute, config.DefaultTTL)
	assert.Equal(t, 100, config.MaxSize)
	assert.Equal(t, 10*time.Second, config.TTLByOperation["info.stats"])
	assert.Equal(t, 30*time.Second, config.TTLByOperation["nodes.list"])
	assert.Equal(t, 10*time.Second, config.TTLByOperation["jobs.list"])
	assert.Equal(t, 30*time.Second, config.TTLByOperation["jobs.get"])
}

func TestCacheItem_IsExpired(t *testing.T) {
	t.Run("not expired", func(t *testing.T) {
		item := &CacheItem{
			Expiry: time.Now().Add(5 * time.Minute),
		}
		assert.False(t, item.IsExpired())
	})

	t.Run("expired", func(t *testing.T) {
		item := &CacheItem{
			Expiry: time.Now().Add(-5 * time.Minute),
		}
		assert.True(t, item.IsExpired())
	})
}

func TestNewResponseCache(t *testing.T) {
	t.Run("with config", func(t *testing.T) {
		config := &CacheConfig{
			DefaultTTL:      1 * time.Minute,
			MaxSize:         100,
			CleanupInterval: 30 * time.Second,
		}

		cache := NewResponseCache(config)
		defer cache.Close()

		require.NotNil(t, cache)
		assert.Equal(t, config, cache.config)
		assert.NotNil(t, cache.items)
		assert.NotNil(t, cache.stopCh)
	})

	t.Run("with nil config", func(t *testing.T) {
		cache := NewResponseCache(nil)
		defer cache.Close()

		require.NotNil(t, cache)
		assert.Equal(t, DefaultCacheConfig(), cache.config)
	})

	t.Run("no cleanup when interval is zero", func(t *testing.T) {
		config := &CacheConfig{
			DefaultTTL:      1 * time.Minute,
			MaxSize:         100,
			CleanupInterval: 0, // No cleanup
		}

		cache := NewResponseCache(config)
		defer cache.Close()

		require.NotNil(t, cache)
		assert.Nil(t, cache.cleanup)
	})
}

func TestResponseCache_GenerateKey(t *testing.T) {
	cache := NewResponseCache(nil)
	defer cache.Close()

	params1 := map[string]interface{}{
		"id":   "123",
		"type": "job",
	}

	params2 := map[string]interface{}{
		"type": "job",
		"id":   "123",
	}

	params3 := map[string]interface{}{
		"id":   "456",
		"type": "job",
	}

	key1 := cache.GenerateKey("jobs.get", params1)
	key2 := cache.GenerateKey("jobs.get", params2)
	key3 := cache.GenerateKey("jobs.get", params3)

	// Same parameters in different order should generate same key
	assert.Equal(t, key1, key2)

	// Different parameters should generate different keys
	assert.NotEqual(t, key1, key3)

	// Key should start with operation name
	assert.Contains(t, key1, "jobs.get:")
}

func TestResponseCache_SetAndGet(t *testing.T) {
	cache := NewResponseCache(&CacheConfig{
		DefaultTTL: 1 * time.Minute,
		MaxSize:    10,
	})
	defer cache.Close()

	operation := "jobs.get"
	params := map[string]interface{}{"id": "123"}
	value := []byte("test response data")

	// Test cache miss
	result, found := cache.Get(operation, params)
	assert.False(t, found)
	assert.Nil(t, result)

	// Set value
	cache.Set(operation, params, value)

	// Test cache hit
	result, found = cache.Get(operation, params)
	assert.True(t, found)
	assert.Equal(t, value, result)

	// Verify stats
	stats := cache.GetStats()
	assert.Equal(t, int64(1), stats.Misses)
	assert.Equal(t, int64(1), stats.Sets)
	assert.Equal(t, int64(1), stats.Hits)
	assert.Equal(t, int64(1), stats.CurrentItems)
}

func TestResponseCache_ExpiredItems(t *testing.T) {
	cache := NewResponseCache(&CacheConfig{
		DefaultTTL: 1 * time.Millisecond, // Very short TTL
		MaxSize:    10,
	})
	defer cache.Close()

	operation := "jobs.get"
	params := map[string]interface{}{"id": "123"}
	value := []byte("test response data")

	// Set value
	cache.Set(operation, params, value)

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	// Should be expired and return cache miss
	result, found := cache.Get(operation, params)
	assert.False(t, found)
	assert.Nil(t, result)

	// Verify stats include eviction
	stats := cache.GetStats()
	assert.Equal(t, int64(1), stats.Evictions)
}

func TestResponseCache_OperationSpecificTTL(t *testing.T) {
	config := &CacheConfig{
		DefaultTTL: 5 * time.Minute,
		MaxSize:    10,
		TTLByOperation: map[string]time.Duration{
			"jobs.get": 10 * time.Millisecond,
		},
	}

	cache := NewResponseCache(config)
	defer cache.Close()

	// Set value with operation-specific TTL
	cache.Set("jobs.get", map[string]interface{}{"id": "123"}, []byte("data"))

	// Should be available immediately
	_, found := cache.Get("jobs.get", map[string]interface{}{"id": "123"})
	assert.True(t, found)

	// Wait for operation-specific TTL to expire
	time.Sleep(20 * time.Millisecond)

	// Should be expired
	_, found = cache.Get("jobs.get", map[string]interface{}{"id": "123"})
	assert.False(t, found)
}

func TestResponseCache_MaxSizeEviction(t *testing.T) {
	cache := NewResponseCache(&CacheConfig{
		DefaultTTL: 10 * time.Minute,
		MaxSize:    2, // Small size to trigger eviction
	})
	defer cache.Close()

	// Add first item
	cache.Set("jobs.get", map[string]interface{}{"id": "1"}, []byte("data1"))
	time.Sleep(time.Millisecond) // Ensure different access times

	// Add second item
	cache.Set("jobs.get", map[string]interface{}{"id": "2"}, []byte("data2"))
	time.Sleep(time.Millisecond)

	// Add third item - should evict the least recently used (first item)
	cache.Set("jobs.get", map[string]interface{}{"id": "3"}, []byte("data3"))

	// First item should be evicted
	_, found1 := cache.Get("jobs.get", map[string]interface{}{"id": "1"})
	assert.False(t, found1)

	// Second and third items should still be there
	_, found2 := cache.Get("jobs.get", map[string]interface{}{"id": "2"})
	_, found3 := cache.Get("jobs.get", map[string]interface{}{"id": "3"})
	assert.True(t, found2)
	assert.True(t, found3)

	// Verify eviction in stats
	stats := cache.GetStats()
	assert.Equal(t, int64(1), stats.Evictions)
}

func TestResponseCache_Delete(t *testing.T) {
	cache := NewResponseCache(nil)
	defer cache.Close()

	operation := "jobs.get"
	params := map[string]interface{}{"id": "123"}

	// Set value
	cache.Set(operation, params, []byte("data"))

	// Verify it exists
	_, found := cache.Get(operation, params)
	assert.True(t, found)

	// Delete it
	cache.Delete(operation, params)

	// Verify it's gone
	_, found = cache.Get(operation, params)
	assert.False(t, found)

	// Verify stats
	stats := cache.GetStats()
	assert.Equal(t, int64(1), stats.Deletions)
}

func TestResponseCache_InvalidatePattern(t *testing.T) {
	cache := NewResponseCache(nil)
	defer cache.Close()

	// Add multiple items
	cache.Set("jobs.get", map[string]interface{}{"id": "1"}, []byte("data1"))
	cache.Set("jobs.list", map[string]interface{}{}, []byte("data2"))
	cache.Set("nodes.get", map[string]interface{}{"id": "1"}, []byte("data3"))

	// Invalidate all jobs operations
	count := cache.InvalidatePattern("jobs.*")
	assert.Equal(t, 2, count)

	// Verify jobs operations are gone
	_, found1 := cache.Get("jobs.get", map[string]interface{}{"id": "1"})
	_, found2 := cache.Get("jobs.list", map[string]interface{}{})
	assert.False(t, found1)
	assert.False(t, found2)

	// Verify nodes operation is still there
	_, found3 := cache.Get("nodes.get", map[string]interface{}{"id": "1"})
	assert.True(t, found3)

	// Verify stats
	stats := cache.GetStats()
	assert.Equal(t, int64(1), stats.PatternInvalidations)
}

func TestResponseCache_Clear(t *testing.T) {
	cache := NewResponseCache(nil)
	defer cache.Close()

	// Add multiple items
	cache.Set("jobs.get", map[string]interface{}{"id": "1"}, []byte("data1"))
	cache.Set("jobs.list", map[string]interface{}{}, []byte("data2"))

	// Verify items exist
	stats := cache.GetStats()
	assert.Equal(t, int64(2), stats.CurrentItems)

	// Clear cache
	cache.Clear()

	// Verify all items are gone
	stats = cache.GetStats()
	assert.Equal(t, int64(0), stats.CurrentItems)
	assert.Equal(t, int64(1), stats.Clears)
	assert.Equal(t, int64(2), stats.Evictions) // Clear counts as evictions
}

func TestResponseCache_GetDetailedStats(t *testing.T) {
	cache := NewResponseCache(nil)
	defer cache.Close()

	// Add item
	cache.Set("jobs.get", map[string]interface{}{"id": "123"}, []byte("test data"))

	// Access it to update hit count
	cache.Get("jobs.get", map[string]interface{}{"id": "123"})

	detailedStats := cache.GetDetailedStats()

	assert.Equal(t, int64(1), detailedStats.Basic.CurrentItems)
	assert.Len(t, detailedStats.Items, 1)

	item := detailedStats.Items[0]
	assert.Contains(t, item.Key, "jobs.get:")
	assert.Equal(t, int64(9), item.Size) // "test data" = 9 bytes
	assert.Equal(t, int64(1), item.HitCount)
	assert.True(t, item.TTL > 0)
	assert.True(t, item.Age >= 0)
}

func TestMatchesPattern(t *testing.T) {
	tests := []struct {
		key     string
		pattern string
		matches bool
	}{
		{"jobs.get:abc123", "jobs.*", true},
		{"jobs.list:def456", "jobs.*", true},
		{"nodes.get:ghi789", "jobs.*", false},
		{"anything", "*", true},
		{"exact.match", "exact.match", true},
		{"not.match", "exact.match", false},
		{"prefix.something", "prefix.*", true},
		{"other.something", "prefix.*", false},
	}

	for _, tt := range tests {
		t.Run(tt.key+"_"+tt.pattern, func(t *testing.T) {
			result := matchesPattern(tt.key, tt.pattern)
			assert.Equal(t, tt.matches, result)
		})
	}
}

func TestNewCacheManager(t *testing.T) {
	t.Run("with config", func(t *testing.T) {
		config := &CacheConfig{MaxSize: 100}
		manager := NewCacheManager(config)
		defer manager.Close()

		require.NotNil(t, manager)
		assert.Equal(t, config, manager.config)
	})

	t.Run("with nil config", func(t *testing.T) {
		manager := NewCacheManager(nil)
		defer manager.Close()

		require.NotNil(t, manager)
		assert.Equal(t, DefaultCacheConfig(), manager.config)
	})
}

func TestCacheManager_GetCache(t *testing.T) {
	manager := NewCacheManager(nil)
	defer manager.Close()

	// First call creates cache
	cache1 := manager.GetCache("v0.0.41")
	require.NotNil(t, cache1)

	// Second call returns same cache
	cache2 := manager.GetCache("v0.0.41")
	assert.Equal(t, cache1, cache2)

	// Different context returns different cache
	cache3 := manager.GetCache("v0.0.42")
	assert.NotEqual(t, cache1, cache3)
}

func TestCacheManager_InvalidateAll(t *testing.T) {
	manager := NewCacheManager(nil)
	defer manager.Close()

	// Create caches and add data
	cache1 := manager.GetCache("v0.0.41")
	cache2 := manager.GetCache("v0.0.42")

	cache1.Set("jobs.get", map[string]interface{}{"id": "1"}, []byte("data1"))
	cache2.Set("jobs.get", map[string]interface{}{"id": "2"}, []byte("data2"))

	// Verify data exists
	_, found1 := cache1.Get("jobs.get", map[string]interface{}{"id": "1"})
	_, found2 := cache2.Get("jobs.get", map[string]interface{}{"id": "2"})
	assert.True(t, found1)
	assert.True(t, found2)

	// Invalidate all
	manager.InvalidateAll()

	// Verify data is gone
	_, found1 = cache1.Get("jobs.get", map[string]interface{}{"id": "1"})
	_, found2 = cache2.Get("jobs.get", map[string]interface{}{"id": "2"})
	assert.False(t, found1)
	assert.False(t, found2)
}

func TestCacheManager_GetGlobalStats(t *testing.T) {
	manager := NewCacheManager(nil)
	defer manager.Close()

	// Create caches and add data
	cache1 := manager.GetCache("v0.0.41")
	cache2 := manager.GetCache("v0.0.42")

	cache1.Set("jobs.get", map[string]interface{}{"id": "1"}, []byte("data1"))
	cache2.Set("jobs.get", map[string]interface{}{"id": "2"}, []byte("data2"))

	stats := manager.GetGlobalStats()

	assert.Len(t, stats, 2)
	assert.Contains(t, stats, "v0.0.41")
	assert.Contains(t, stats, "v0.0.42")

	assert.Equal(t, int64(1), stats["v0.0.41"].CurrentItems)
	assert.Equal(t, int64(1), stats["v0.0.42"].CurrentItems)
}

func TestGetCacheConfigForProfile(t *testing.T) {
	tests := []struct {
		profile  PerformanceProfile
		expected func(*CacheConfig) bool
	}{
		{
			ProfileHighThroughput,
			func(c *CacheConfig) bool { return c.MaxSize == 5000 },
		},
		{
			ProfileLowLatency,
			func(c *CacheConfig) bool { return c.DefaultTTL == 30*time.Second },
		},
		{
			ProfileConservative,
			func(c *CacheConfig) bool { return c.MaxSize == 100 },
		},
		{
			ProfileBatch,
			func(c *CacheConfig) bool { return c.DefaultTTL == 30*time.Minute && c.MaxSize == 10000 },
		},
		{
			PerformanceProfile("unknown"),
			func(c *CacheConfig) bool { return c.DefaultTTL == 5*time.Minute },
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.profile), func(t *testing.T) {
			config := GetCacheConfigForProfile(tt.profile)
			require.NotNil(t, config)
			assert.True(t, tt.expected(config))
		})
	}
}

func TestResponseCache_Cleanup(t *testing.T) {
	cache := NewResponseCache(&CacheConfig{
		DefaultTTL:      10 * time.Millisecond,
		MaxSize:         10,
		CleanupInterval: 5 * time.Millisecond,
	})
	defer cache.Close()

	// Add item that will expire
	cache.Set("jobs.get", map[string]interface{}{"id": "123"}, []byte("data"))

	// Verify it exists
	_, found := cache.Get("jobs.get", map[string]interface{}{"id": "123"})
	assert.True(t, found)

	// Wait for cleanup to run (expiration + cleanup interval + buffer)
	time.Sleep(30 * time.Millisecond)

	// Item should be cleaned up
	stats := cache.GetStats()
	assert.Equal(t, int64(0), stats.CurrentItems)
}

func TestResponseCache_Close(t *testing.T) {
	cache := NewResponseCache(&CacheConfig{
		CleanupInterval: 1 * time.Millisecond,
	})

	// Add some data
	cache.Set("test", map[string]interface{}{}, []byte("data"))

	// Close should not panic
	cache.Close()

	// Calling close again should not panic
	cache.Close()
}
