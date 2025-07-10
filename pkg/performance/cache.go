package performance

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"sync"
	"time"
)

// CacheConfig holds configuration for the response cache
type CacheConfig struct {
	// DefaultTTL is the default time-to-live for cached items
	DefaultTTL time.Duration

	// MaxSize is the maximum number of items to store in cache
	MaxSize int

	// EnableCompression enables compression of cached values
	EnableCompression bool

	// CleanupInterval specifies how often to run cache cleanup
	CleanupInterval time.Duration

	// TTLByOperation specifies custom TTL for specific operations
	TTLByOperation map[string]time.Duration
}

// DefaultCacheConfig returns a default cache configuration
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		DefaultTTL:        5 * time.Minute,
		MaxSize:           1000,
		EnableCompression: true,
		CleanupInterval:   1 * time.Minute,
		TTLByOperation: map[string]time.Duration{
			"info.version":     30 * time.Minute, // API version info rarely changes
			"info.stats":       30 * time.Second,  // Cluster stats change frequently
			"info.get":         10 * time.Minute,  // Cluster info changes moderately
			"partitions.list":  10 * time.Minute,  // Partitions change infrequently
			"partitions.get":   10 * time.Minute,
			"nodes.list":       2 * time.Minute,   // Node states change frequently
			"nodes.get":        2 * time.Minute,
			"jobs.list":        30 * time.Second,  // Job lists change very frequently
			"jobs.get":         1 * time.Minute,   // Individual job details change frequently
		},
	}
}

// AggressiveCacheConfig returns configuration for aggressive caching
func AggressiveCacheConfig() *CacheConfig {
	config := DefaultCacheConfig()
	config.MaxSize = 5000
	config.DefaultTTL = 10 * time.Minute
	config.TTLByOperation["info.stats"] = 2 * time.Minute
	config.TTLByOperation["nodes.list"] = 5 * time.Minute
	config.TTLByOperation["jobs.list"] = 2 * time.Minute
	config.TTLByOperation["jobs.get"] = 5 * time.Minute
	return config
}

// ConservativeCacheConfig returns configuration for minimal caching
func ConservativeCacheConfig() *CacheConfig {
	config := DefaultCacheConfig()
	config.MaxSize = 100
	config.DefaultTTL = 1 * time.Minute
	config.TTLByOperation["info.stats"] = 10 * time.Second
	config.TTLByOperation["nodes.list"] = 30 * time.Second
	config.TTLByOperation["jobs.list"] = 10 * time.Second
	config.TTLByOperation["jobs.get"] = 30 * time.Second
	return config
}

// CacheItem represents an item stored in the cache
type CacheItem struct {
	Key        string    `json:"key"`
	Value      []byte    `json:"value"`
	Expiry     time.Time `json:"expiry"`
	CreatedAt  time.Time `json:"created_at"`
	AccessedAt time.Time `json:"accessed_at"`
	HitCount   int64     `json:"hit_count"`
	Compressed bool      `json:"compressed"`
}

// IsExpired checks if the cache item has expired
func (item *CacheItem) IsExpired() bool {
	return time.Now().After(item.Expiry)
}

// ResponseCache provides intelligent caching for API responses
type ResponseCache struct {
	items   map[string]*CacheItem
	config  *CacheConfig
	mutex   sync.RWMutex
	stats   CacheStats
	cleanup *time.Ticker
	stopCh  chan struct{}
}

// NewResponseCache creates a new response cache with the given configuration
func NewResponseCache(config *CacheConfig) *ResponseCache {
	if config == nil {
		config = DefaultCacheConfig()
	}

	cache := &ResponseCache{
		items:  make(map[string]*CacheItem),
		config: config,
		stats:  CacheStats{},
		stopCh: make(chan struct{}),
	}

	// Start cleanup goroutine
	if config.CleanupInterval > 0 {
		cache.cleanup = time.NewTicker(config.CleanupInterval)
		go cache.runCleanup()
	}

	return cache
}

// GenerateKey creates a cache key from operation and parameters
func (c *ResponseCache) GenerateKey(operation string, params map[string]interface{}) string {
	// Create a deterministic key from operation and parameters
	data := map[string]interface{}{
		"operation": operation,
		"params":    params,
	}

	// Convert to JSON for consistent ordering
	jsonData, _ := json.Marshal(data)
	
	// Create MD5 hash for compact key
	hash := md5.Sum(jsonData)
	return operation + ":" + hex.EncodeToString(hash[:])
}

// Get retrieves a value from the cache
func (c *ResponseCache) Get(operation string, params map[string]interface{}) ([]byte, bool) {
	key := c.GenerateKey(operation, params)
	
	c.mutex.RLock()
	item, exists := c.items[key]
	c.mutex.RUnlock()

	if !exists {
		c.updateStats(func(s *CacheStats) { s.Misses++ })
		return nil, false
	}

	if item.IsExpired() {
		// Remove expired item
		c.mutex.Lock()
		delete(c.items, key)
		c.mutex.Unlock()
		
		c.updateStats(func(s *CacheStats) { s.Misses++; s.Evictions++ })
		return nil, false
	}

	// Update access statistics
	c.mutex.Lock()
	item.AccessedAt = time.Now()
	item.HitCount++
	c.mutex.Unlock()

	c.updateStats(func(s *CacheStats) { s.Hits++ })

	value := item.Value
	if item.Compressed && c.config.EnableCompression {
		// In a production implementation, you would decompress here
		// For this example, we'll assume the value is already decompressed
	}

	return value, true
}

// Set stores a value in the cache
func (c *ResponseCache) Set(operation string, params map[string]interface{}, value []byte) {
	key := c.GenerateKey(operation, params)
	
	// Determine TTL for this operation
	ttl := c.config.DefaultTTL
	if operationTTL, exists := c.config.TTLByOperation[operation]; exists {
		ttl = operationTTL
	}

	now := time.Now()
	item := &CacheItem{
		Key:        key,
		Value:      value,
		Expiry:     now.Add(ttl),
		CreatedAt:  now,
		AccessedAt: now,
		HitCount:   0,
		Compressed: false,
	}

	// Compress if enabled and worthwhile
	if c.config.EnableCompression && len(value) > 1024 {
		// In a production implementation, you would compress here
		// For this example, we'll skip actual compression
		item.Compressed = true
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if we need to evict items to make space
	if len(c.items) >= c.config.MaxSize {
		c.evictLRU()
	}

	c.items[key] = item
	c.updateStatsUnsafe(func(s *CacheStats) { s.Sets++ })
}

// Delete removes an item from the cache
func (c *ResponseCache) Delete(operation string, params map[string]interface{}) {
	key := c.GenerateKey(operation, params)
	
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.items[key]; exists {
		delete(c.items, key)
		c.updateStatsUnsafe(func(s *CacheStats) { s.Deletions++ })
	}
}

// InvalidatePattern removes all items matching a pattern
func (c *ResponseCache) InvalidatePattern(pattern string) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var keysToDelete []string
	for key := range c.items {
		// Simple pattern matching - in production, use a more sophisticated approach
		if matchesPattern(key, pattern) {
			keysToDelete = append(keysToDelete, key)
		}
	}

	for _, key := range keysToDelete {
		delete(c.items, key)
	}

	c.updateStatsUnsafe(func(s *CacheStats) { s.PatternInvalidations++ })
	return len(keysToDelete)
}

// Clear removes all items from the cache
func (c *ResponseCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	itemCount := len(c.items)
	c.items = make(map[string]*CacheItem)
	c.updateStatsUnsafe(func(s *CacheStats) { 
		s.Evictions += int64(itemCount)
		s.Clears++ 
	})
}

// GetStats returns cache statistics
func (c *ResponseCache) GetStats() CacheStats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	stats := c.stats
	stats.CurrentItems = int64(len(c.items))
	
	// Calculate additional metrics
	if stats.Hits+stats.Misses > 0 {
		stats.HitRatio = float64(stats.Hits) / float64(stats.Hits+stats.Misses)
	}

	return stats
}

// GetDetailedStats returns detailed cache statistics
func (c *ResponseCache) GetDetailedStats() DetailedCacheStats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	stats := DetailedCacheStats{
		Basic: c.GetStats(),
		Items: make([]CacheItemStats, 0, len(c.items)),
	}

	var totalSize int64
	for _, item := range c.items {
		itemSize := int64(len(item.Value))
		totalSize += itemSize
		
		stats.Items = append(stats.Items, CacheItemStats{
			Key:        item.Key,
			Size:       itemSize,
			TTL:        time.Until(item.Expiry),
			Age:        time.Since(item.CreatedAt),
			HitCount:   item.HitCount,
			LastAccess: item.AccessedAt,
		})
	}

	stats.Basic.TotalSize = totalSize
	return stats
}

// Close stops the cache and cleanup goroutine
func (c *ResponseCache) Close() {
	if c.cleanup != nil {
		c.cleanup.Stop()
	}
	close(c.stopCh)
}

// evictLRU removes the least recently used item
func (c *ResponseCache) evictLRU() {
	var oldestKey string
	var oldestTime time.Time

	for key, item := range c.items {
		if oldestKey == "" || item.AccessedAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.AccessedAt
		}
	}

	if oldestKey != "" {
		delete(c.items, oldestKey)
		c.updateStatsUnsafe(func(s *CacheStats) { s.Evictions++ })
	}
}

// runCleanup periodically removes expired items
func (c *ResponseCache) runCleanup() {
	for {
		select {
		case <-c.cleanup.C:
			c.cleanupExpired()
		case <-c.stopCh:
			return
		}
	}
}

// cleanupExpired removes all expired items
func (c *ResponseCache) cleanupExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var expiredKeys []string
	for key, item := range c.items {
		if item.IsExpired() {
			expiredKeys = append(expiredKeys, key)
		}
	}

	for _, key := range expiredKeys {
		delete(c.items, key)
	}

	if len(expiredKeys) > 0 {
		c.updateStatsUnsafe(func(s *CacheStats) { s.Evictions += int64(len(expiredKeys)) })
	}
}

// updateStats safely updates cache statistics
func (c *ResponseCache) updateStats(fn func(*CacheStats)) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.updateStatsUnsafe(fn)
}

// updateStatsUnsafe updates cache statistics without locking
func (c *ResponseCache) updateStatsUnsafe(fn func(*CacheStats)) {
	fn(&c.stats)
}

// CacheStats represents cache statistics
type CacheStats struct {
	Hits                 int64   `json:"hits"`
	Misses               int64   `json:"misses"`
	Sets                 int64   `json:"sets"`
	Deletions            int64   `json:"deletions"`
	Evictions            int64   `json:"evictions"`
	Clears               int64   `json:"clears"`
	PatternInvalidations int64   `json:"pattern_invalidations"`
	CurrentItems         int64   `json:"current_items"`
	TotalSize            int64   `json:"total_size_bytes"`
	HitRatio             float64 `json:"hit_ratio"`
}

// DetailedCacheStats includes per-item statistics
type DetailedCacheStats struct {
	Basic CacheStats      `json:"basic"`
	Items []CacheItemStats `json:"items"`
}

// CacheItemStats represents statistics for a single cache item
type CacheItemStats struct {
	Key        string        `json:"key"`
	Size       int64         `json:"size_bytes"`
	TTL        time.Duration `json:"ttl_remaining"`
	Age        time.Duration `json:"age"`
	HitCount   int64         `json:"hit_count"`
	LastAccess time.Time     `json:"last_access"`
}

// CacheManager manages multiple caches for different purposes
type CacheManager struct {
	caches map[string]*ResponseCache
	config *CacheConfig
	mutex  sync.RWMutex
}

// NewCacheManager creates a new cache manager
func NewCacheManager(config *CacheConfig) *CacheManager {
	if config == nil {
		config = DefaultCacheConfig()
	}

	return &CacheManager{
		caches: make(map[string]*ResponseCache),
		config: config,
	}
}

// GetCache returns a cache for the given context (e.g., API version)
func (m *CacheManager) GetCache(context string) *ResponseCache {
	m.mutex.RLock()
	cache, exists := m.caches[context]
	m.mutex.RUnlock()

	if exists {
		return cache
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Double-check after acquiring write lock
	if cache, exists := m.caches[context]; exists {
		return cache
	}

	cache = NewResponseCache(m.config)
	m.caches[context] = cache
	return cache
}

// InvalidateAll invalidates all caches
func (m *CacheManager) InvalidateAll() {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, cache := range m.caches {
		cache.Clear()
	}
}

// GetGlobalStats returns statistics for all caches
func (m *CacheManager) GetGlobalStats() map[string]CacheStats {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	stats := make(map[string]CacheStats)
	for context, cache := range m.caches {
		stats[context] = cache.GetStats()
	}

	return stats
}

// Close closes all caches
func (m *CacheManager) Close() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, cache := range m.caches {
		cache.Close()
	}

	m.caches = make(map[string]*ResponseCache)
}

// Helper functions

// matchesPattern checks if a key matches a simple pattern
func matchesPattern(key, pattern string) bool {
	// Simple pattern matching - supports '*' wildcard at the end
	if pattern == "*" {
		return true
	}
	
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(key) >= len(prefix) && key[:len(prefix)] == prefix
	}
	
	return key == pattern
}

// GetCacheConfigForProfile returns cache configuration optimized for a performance profile
func GetCacheConfigForProfile(profile PerformanceProfile) *CacheConfig {
	switch profile {
	case ProfileHighThroughput:
		return AggressiveCacheConfig()
	
	case ProfileLowLatency:
		config := DefaultCacheConfig()
		config.DefaultTTL = 30 * time.Second // Shorter TTL for fresher data
		config.TTLByOperation["jobs.list"] = 5 * time.Second
		config.TTLByOperation["jobs.get"] = 15 * time.Second
		config.TTLByOperation["info.stats"] = 10 * time.Second
		return config
	
	case ProfileConservative:
		return ConservativeCacheConfig()
	
	case ProfileBatch:
		config := AggressiveCacheConfig()
		config.DefaultTTL = 30 * time.Minute // Much longer TTL for batch processing
		config.MaxSize = 10000
		return config
	
	default:
		return DefaultCacheConfig()
	}
}