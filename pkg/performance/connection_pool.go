// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package performance

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"
)

// ConnectionPoolConfig holds configuration for HTTP connection pooling
type ConnectionPoolConfig struct {
	// MaxIdleConns controls the maximum number of idle connections
	MaxIdleConns int

	// MaxIdleConnsPerHost controls the maximum idle connections per host
	MaxIdleConnsPerHost int

	// MaxConnsPerHost controls the maximum connections per host
	MaxConnsPerHost int

	// IdleConnTimeout is the maximum time an idle connection will remain idle
	IdleConnTimeout time.Duration

	// TLSHandshakeTimeout specifies the maximum amount of time waiting for TLS handshake
	TLSHandshakeTimeout time.Duration

	// ExpectContinueTimeout specifies the amount of time to wait for a server's
	// first response headers after fully writing the request headers
	ExpectContinueTimeout time.Duration

	// ResponseHeaderTimeout specifies the amount of time to wait for a server's
	// response headers after fully writing the request
	ResponseHeaderTimeout time.Duration

	// DisableCompression disables compression if true
	DisableCompression bool

	// DisableKeepAlives disables HTTP keep-alives if true
	DisableKeepAlives bool

	// EnableTLSSessionResumption enables TLS session resumption for performance
	EnableTLSSessionResumption bool

	// TLSInsecureSkipVerify controls whether TLS certificates are verified
	TLSInsecureSkipVerify bool
}

// DefaultConnectionPoolConfig returns a performance-optimized configuration
func DefaultConnectionPoolConfig() *ConnectionPoolConfig {
	return &ConnectionPoolConfig{
		MaxIdleConns:               100,
		MaxIdleConnsPerHost:        30,
		MaxConnsPerHost:            100,
		IdleConnTimeout:            90 * time.Second,
		TLSHandshakeTimeout:        10 * time.Second,
		ExpectContinueTimeout:      1 * time.Second,
		ResponseHeaderTimeout:      30 * time.Second,
		DisableCompression:         false,
		DisableKeepAlives:          false,
		EnableTLSSessionResumption: true,
		TLSInsecureSkipVerify:      false,
	}
}

// HighPerformanceConnectionPoolConfig returns configuration optimized for high-throughput scenarios
func HighPerformanceConnectionPoolConfig() *ConnectionPoolConfig {
	return &ConnectionPoolConfig{
		MaxIdleConns:               200,
		MaxIdleConnsPerHost:        50,
		MaxConnsPerHost:            200,
		IdleConnTimeout:            120 * time.Second,
		TLSHandshakeTimeout:        5 * time.Second,
		ExpectContinueTimeout:      500 * time.Millisecond,
		ResponseHeaderTimeout:      15 * time.Second,
		DisableCompression:         false,
		DisableKeepAlives:          false,
		EnableTLSSessionResumption: true,
		TLSInsecureSkipVerify:      false,
	}
}

// ConservativeConnectionPoolConfig returns configuration optimized for resource-constrained environments
func ConservativeConnectionPoolConfig() *ConnectionPoolConfig {
	return &ConnectionPoolConfig{
		MaxIdleConns:               10,
		MaxIdleConnsPerHost:        5,
		MaxConnsPerHost:            20,
		IdleConnTimeout:            30 * time.Second,
		TLSHandshakeTimeout:        15 * time.Second,
		ExpectContinueTimeout:      2 * time.Second,
		ResponseHeaderTimeout:      60 * time.Second,
		DisableCompression:         false,
		DisableKeepAlives:          false,
		EnableTLSSessionResumption: true,
		TLSInsecureSkipVerify:      false,
	}
}

// HTTPClientPool manages HTTP clients with optimized connection pooling
type HTTPClientPool struct {
	clients map[string]*http.Client
	config  *ConnectionPoolConfig
	mutex   sync.RWMutex
}

// NewHTTPClientPool creates a new HTTP client pool with the given configuration
func NewHTTPClientPool(config *ConnectionPoolConfig) *HTTPClientPool {
	if config == nil {
		config = DefaultConnectionPoolConfig()
	}

	return &HTTPClientPool{
		clients: make(map[string]*http.Client),
		config:  config,
	}
}

// GetClient returns an HTTP client optimized for the given endpoint
func (p *HTTPClientPool) GetClient(endpoint string) *http.Client {
	p.mutex.RLock()
	client, exists := p.clients[endpoint]
	p.mutex.RUnlock()

	if exists {
		return client
	}

	// Create new optimized client
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Double-check after acquiring write lock
	if client, exists := p.clients[endpoint]; exists {
		return client
	}

	client = p.createOptimizedClient()
	p.clients[endpoint] = client
	return client
}

// createOptimizedClient creates an HTTP client with optimized transport settings
func (p *HTTPClientPool) createOptimizedClient() *http.Client {
	// Create custom dialer with optimized settings
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true, // Enable both IPv4 and IPv6
	}

	// Create optimized transport
	transport := &http.Transport{
		DialContext:           dialer.DialContext,
		MaxIdleConns:          p.config.MaxIdleConns,
		MaxIdleConnsPerHost:   p.config.MaxIdleConnsPerHost,
		MaxConnsPerHost:       p.config.MaxConnsPerHost,
		IdleConnTimeout:       p.config.IdleConnTimeout,
		TLSHandshakeTimeout:   p.config.TLSHandshakeTimeout,
		ExpectContinueTimeout: p.config.ExpectContinueTimeout,
		ResponseHeaderTimeout: p.config.ResponseHeaderTimeout,
		DisableCompression:    p.config.DisableCompression,
		DisableKeepAlives:     p.config.DisableKeepAlives,
		ForceAttemptHTTP2:     true,      // Enable HTTP/2 when possible
		WriteBufferSize:       64 * 1024, // 64KB write buffer
		ReadBufferSize:        64 * 1024, // 64KB read buffer
	}

	// Configure TLS settings
	if p.config.EnableTLSSessionResumption || p.config.TLSInsecureSkipVerify {
		// #nosec G402 -- InsecureSkipVerify is user-configurable and defaults to false (secure)
		// This allows users to optionally disable TLS verification for testing/development
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: p.config.TLSInsecureSkipVerify,    // #nosec G402
			ClientSessionCache: tls.NewLRUClientSessionCache(256), // Enable session resumption
		}
	}

	return &http.Client{
		Transport: transport,
		Timeout:   0, // No timeout at client level - handle at request level
	}
}

// GetStats returns connection pool statistics
func (p *HTTPClientPool) GetStats() ConnectionPoolStats {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	stats := ConnectionPoolStats{
		ActiveClients: len(p.clients),
		Config:        *p.config,
	}

	// Collect transport statistics from each client
	for endpoint, client := range p.clients {
		if _, ok := client.Transport.(*http.Transport); ok {
			stats.Endpoints = append(stats.Endpoints, EndpointStats{
				Endpoint: endpoint,
				// Note: http.Transport doesn't expose detailed connection stats
				// In a production implementation, you might use a custom transport
				// or metrics collection to gather more detailed statistics
			})
		}
	}

	return stats
}

// Close closes all connections in the pool
func (p *HTTPClientPool) Close() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, client := range p.clients {
		if transport, ok := client.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
	}

	p.clients = make(map[string]*http.Client)
}

// ConnectionPoolStats represents statistics about the connection pool
type ConnectionPoolStats struct {
	ActiveClients int                  `json:"active_clients"`
	Endpoints     []EndpointStats      `json:"endpoints"`
	Config        ConnectionPoolConfig `json:"config"`
}

// EndpointStats represents statistics for a specific endpoint
type EndpointStats struct {
	Endpoint        string `json:"endpoint"`
	ActiveConns     int    `json:"active_connections"`
	IdleConns       int    `json:"idle_connections"`
	TotalRequests   int64  `json:"total_requests"`
	FailedRequests  int64  `json:"failed_requests"`
	AvgResponseTime string `json:"avg_response_time"`
}

// PerformanceProfile represents different performance optimization profiles
type PerformanceProfile string

const (
	// ProfileDefault provides balanced performance and resource usage
	ProfileDefault PerformanceProfile = "default"

	// ProfileHighThroughput optimizes for maximum throughput
	ProfileHighThroughput PerformanceProfile = "high_throughput"

	// ProfileLowLatency optimizes for minimum latency
	ProfileLowLatency PerformanceProfile = "low_latency"

	// ProfileConservative minimizes resource usage
	ProfileConservative PerformanceProfile = "conservative"

	// ProfileBatch optimizes for batch processing workloads
	ProfileBatch PerformanceProfile = "batch"
)

// GetConnectionPoolConfigForProfile returns an optimized configuration for the given profile
func GetConnectionPoolConfigForProfile(profile PerformanceProfile) *ConnectionPoolConfig {
	switch profile {
	case ProfileHighThroughput:
		return HighPerformanceConnectionPoolConfig()

	case ProfileLowLatency:
		config := DefaultConnectionPoolConfig()
		config.ExpectContinueTimeout = 100 * time.Millisecond
		config.ResponseHeaderTimeout = 5 * time.Second
		config.TLSHandshakeTimeout = 3 * time.Second
		config.MaxConnsPerHost = 50
		return config

	case ProfileConservative:
		return ConservativeConnectionPoolConfig()

	case ProfileBatch:
		config := HighPerformanceConnectionPoolConfig()
		config.IdleConnTimeout = 300 * time.Second       // Longer idle timeout for batch jobs
		config.ResponseHeaderTimeout = 120 * time.Second // Longer response timeout
		return config

	default:
		return DefaultConnectionPoolConfig()
	}
}

// HTTPClientPoolManager manages multiple connection pools for different API versions
type HTTPClientPoolManager struct {
	pools map[string]*HTTPClientPool
	mutex sync.RWMutex
}

// NewHTTPClientPoolManager creates a new pool manager
func NewHTTPClientPoolManager() *HTTPClientPoolManager {
	return &HTTPClientPoolManager{
		pools: make(map[string]*HTTPClientPool),
	}
}

// GetPoolForVersion returns a connection pool optimized for the given API version
func (m *HTTPClientPoolManager) GetPoolForVersion(version string, profile PerformanceProfile) *HTTPClientPool {
	poolKey := version + ":" + string(profile)

	m.mutex.RLock()
	pool, exists := m.pools[poolKey]
	m.mutex.RUnlock()

	if exists {
		return pool
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Double-check after acquiring write lock
	if pool, exists := m.pools[poolKey]; exists {
		return pool
	}

	// Create new pool with version-specific optimizations
	config := GetConnectionPoolConfigForProfile(profile)

	// Version-specific optimizations
	switch version {
	case "v0.0.40":
		// Older version, more conservative settings
		config.MaxConnsPerHost = minInt(config.MaxConnsPerHost, 20)
		config.ResponseHeaderTimeout = maxDuration(config.ResponseHeaderTimeout, 30*time.Second)

	case "v0.0.41":
		// Stable version with some optimizations
		config.ExpectContinueTimeout = 750 * time.Millisecond

	case "v0.0.42":
		// Stable and performant version
		// Use default optimizations

	case "v0.0.43":
		// Latest version, enable all optimizations
		// Latest version supports HTTP/2 optimizations
		config.ExpectContinueTimeout = 250 * time.Millisecond
	}

	pool = NewHTTPClientPool(config)
	m.pools[poolKey] = pool
	return pool
}

// Close closes all connection pools
func (m *HTTPClientPoolManager) Close() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, pool := range m.pools {
		pool.Close()
	}

	m.pools = make(map[string]*HTTPClientPool)
}

// GetGlobalStats returns statistics for all connection pools
func (m *HTTPClientPoolManager) GetGlobalStats() map[string]ConnectionPoolStats {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	stats := make(map[string]ConnectionPoolStats)
	for key, pool := range m.pools {
		stats[key] = pool.GetStats()
	}

	return stats
}

// Helper functions

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxDuration(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}
