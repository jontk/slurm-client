// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Package pool provides HTTP connection pooling for the SLURM client
package pool

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/jontk/slurm-client/pkg/logging"
)

// HTTPClientPool manages a pool of HTTP clients with optimized transport settings
type HTTPClientPool struct {
	mu      sync.RWMutex
	clients map[string]*pooledClient
	config  *PoolConfig
	logger  logging.Logger
}

// pooledClient wraps an HTTP client with usage statistics
type pooledClient struct {
	client      *http.Client
	created     time.Time
	lastUsed    time.Time
	useCount    int64
	activeConns int32
}

// PoolConfig holds configuration for the HTTP client pool
type PoolConfig struct {
	// MaxIdleConns controls the maximum number of idle connections across all hosts
	MaxIdleConns int

	// MaxIdleConnsPerHost controls the maximum idle connections to keep per-host
	MaxIdleConnsPerHost int

	// MaxConnsPerHost limits the total connections per host
	MaxConnsPerHost int

	// IdleConnTimeout is the timeout for idle connections
	IdleConnTimeout time.Duration

	// TLSHandshakeTimeout specifies the TLS handshake timeout
	TLSHandshakeTimeout time.Duration

	// ExpectContinueTimeout is the timeout for expect-continue
	ExpectContinueTimeout time.Duration

	// DisableKeepAlives disables HTTP keep-alives
	DisableKeepAlives bool

	// DisableCompression disables transport compression
	DisableCompression bool

	// MaxResponseHeaderBytes limits response header size
	MaxResponseHeaderBytes int64
}

// DefaultPoolConfig returns a pool configuration optimized for SLURM API access
func DefaultPoolConfig() *PoolConfig {
	return &PoolConfig{
		MaxIdleConns:           100,
		MaxIdleConnsPerHost:    10,
		MaxConnsPerHost:        50,
		IdleConnTimeout:        90 * time.Second,
		TLSHandshakeTimeout:    10 * time.Second,
		ExpectContinueTimeout:  1 * time.Second,
		DisableKeepAlives:      false,
		DisableCompression:     false,
		MaxResponseHeaderBytes: 1 << 20, // 1 MB
	}
}

// NewHTTPClientPool creates a new HTTP client pool
func NewHTTPClientPool(config *PoolConfig, logger logging.Logger) *HTTPClientPool {
	if config == nil {
		config = DefaultPoolConfig()
	}
	if logger == nil {
		logger = logging.NoOpLogger{}
	}

	return &HTTPClientPool{
		clients: make(map[string]*pooledClient),
		config:  config,
		logger:  logger,
	}
}

// GetClient returns an HTTP client for the specified endpoint
func (p *HTTPClientPool) GetClient(endpoint string) *http.Client {
	p.mu.RLock()
	pc, exists := p.clients[endpoint]
	p.mu.RUnlock()

	if exists {
		// Update usage statistics
		p.mu.Lock()
		pc.lastUsed = time.Now()
		pc.useCount++
		p.mu.Unlock()

		return pc.client
	}

	// Create new client
	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check after acquiring write lock
	if pc, exists := p.clients[endpoint]; exists {
		pc.lastUsed = time.Now()
		pc.useCount++
		return pc.client
	}

	// Create new pooled client
	client := p.createHTTPClient()
	pc = &pooledClient{
		client:   client,
		created:  time.Now(),
		lastUsed: time.Now(),
		useCount: 1,
	}

	p.clients[endpoint] = pc
	p.logger.Info("created new HTTP client for endpoint", "endpoint", endpoint)

	return client
}

// createHTTPClient creates a new HTTP client with optimized transport
func (p *HTTPClientPool) createHTTPClient() *http.Client {
	// Create custom dialer
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	// Create transport with pooling configuration
	transport := &http.Transport{
		Proxy:                  http.ProxyFromEnvironment,
		DialContext:            dialer.DialContext,
		MaxIdleConns:           p.config.MaxIdleConns,
		MaxIdleConnsPerHost:    p.config.MaxIdleConnsPerHost,
		MaxConnsPerHost:        p.config.MaxConnsPerHost,
		IdleConnTimeout:        p.config.IdleConnTimeout,
		TLSHandshakeTimeout:    p.config.TLSHandshakeTimeout,
		ExpectContinueTimeout:  p.config.ExpectContinueTimeout,
		DisableKeepAlives:      p.config.DisableKeepAlives,
		DisableCompression:     p.config.DisableCompression,
		MaxResponseHeaderBytes: p.config.MaxResponseHeaderBytes,

		// Force HTTP/2 support
		ForceAttemptHTTP2: true,

		// TLS configuration
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	return &http.Client{
		Transport: transport,
		Timeout:   0, // No client-level timeout, use context instead
	}
}

// Stats returns statistics about the connection pool
func (p *HTTPClientPool) Stats() PoolStats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	stats := PoolStats{
		TotalClients: len(p.clients),
		ClientStats:  make(map[string]ClientStats),
	}

	for endpoint, pc := range p.clients {
		stats.ClientStats[endpoint] = ClientStats{
			Created:     pc.created,
			LastUsed:    pc.lastUsed,
			UseCount:    pc.useCount,
			ActiveConns: pc.activeConns,
		}
	}

	return stats
}

// CleanupIdleClients removes clients that haven't been used recently
func (p *HTTPClientPool) CleanupIdleClients(maxIdleTime time.Duration) int {
	p.mu.Lock()
	defer p.mu.Unlock()

	removed := 0
	cutoff := time.Now().Add(-maxIdleTime)

	for endpoint, pc := range p.clients {
		if pc.lastUsed.Before(cutoff) && pc.activeConns == 0 {
			// Close idle connections
			if transport, ok := pc.client.Transport.(*http.Transport); ok {
				transport.CloseIdleConnections()
			}

			delete(p.clients, endpoint)
			removed++

			p.logger.Info("removed idle client",
				"endpoint", endpoint,
				"idle_duration", time.Since(pc.lastUsed),
			)
		}
	}

	return removed
}

// Close closes all connections in the pool
func (p *HTTPClientPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for endpoint, pc := range p.clients {
		if transport, ok := pc.client.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
		delete(p.clients, endpoint)
	}

	p.logger.Info("closed all HTTP clients in pool")
	return nil
}

// PoolStats contains statistics about the connection pool
type PoolStats struct {
	TotalClients int
	ClientStats  map[string]ClientStats
}

// ClientStats contains statistics for a single client
type ClientStats struct {
	Created     time.Time
	LastUsed    time.Time
	UseCount    int64
	ActiveConns int32
}

// ConnectionManager manages connection lifecycle and health
type ConnectionManager struct {
	pool            *HTTPClientPool
	healthCheckFunc HealthCheckFunc
	cleanupInterval time.Duration
	maxIdleTime     time.Duration
	ctx             context.Context
	cancel          context.CancelFunc
	wg              sync.WaitGroup
	logger          logging.Logger
}

// HealthCheckFunc is a function that checks if an endpoint is healthy
type HealthCheckFunc func(ctx context.Context, endpoint string, client *http.Client) error

// NewConnectionManager creates a new connection manager
func NewConnectionManager(pool *HTTPClientPool, healthCheck HealthCheckFunc, logger logging.Logger) *ConnectionManager {
	ctx, cancel := context.WithCancel(context.Background())

	if logger == nil {
		logger = logging.NoOpLogger{}
	}

	return &ConnectionManager{
		pool:            pool,
		healthCheckFunc: healthCheck,
		cleanupInterval: 5 * time.Minute,
		maxIdleTime:     15 * time.Minute,
		ctx:             ctx,
		cancel:          cancel,
		logger:          logger,
	}
}

// Start begins the connection management routines
func (cm *ConnectionManager) Start() {
	cm.wg.Add(1)
	go cm.cleanupRoutine()
}

// Stop stops the connection management routines
func (cm *ConnectionManager) Stop() {
	cm.cancel()
	cm.wg.Wait()
}

// cleanupRoutine periodically cleans up idle connections
func (cm *ConnectionManager) cleanupRoutine() {
	defer cm.wg.Done()

	ticker := time.NewTicker(cm.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			removed := cm.pool.CleanupIdleClients(cm.maxIdleTime)
			if removed > 0 {
				cm.logger.Info("cleaned up idle clients", "removed", removed)
			}
		case <-cm.ctx.Done():
			return
		}
	}
}

// GetHealthyClient returns a healthy HTTP client for the endpoint
func (cm *ConnectionManager) GetHealthyClient(ctx context.Context, endpoint string) (*http.Client, error) {
	client := cm.pool.GetClient(endpoint)

	// Perform health check if configured
	if cm.healthCheckFunc != nil {
		if err := cm.healthCheckFunc(ctx, endpoint, client); err != nil {
			return nil, fmt.Errorf("endpoint health check failed: %w", err)
		}
	}

	return client, nil
}
