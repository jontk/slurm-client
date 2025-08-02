// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package pool

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/jontk/slurm-client/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultPoolConfig(t *testing.T) {
	config := DefaultPoolConfig()
	
	require.NotNil(t, config)
	assert.Equal(t, 100, config.MaxIdleConns)
	assert.Equal(t, 10, config.MaxIdleConnsPerHost)
	assert.Equal(t, 50, config.MaxConnsPerHost)
	assert.Equal(t, 90*time.Second, config.IdleConnTimeout)
	assert.Equal(t, 10*time.Second, config.TLSHandshakeTimeout)
	assert.Equal(t, 1*time.Second, config.ExpectContinueTimeout)
	assert.False(t, config.DisableKeepAlives)
	assert.False(t, config.DisableCompression)
	assert.Equal(t, int64(1<<20), config.MaxResponseHeaderBytes)
}

func TestNewHTTPClientPool(t *testing.T) {
	t.Run("with config and logger", func(t *testing.T) {
		config := &PoolConfig{
			MaxIdleConns: 50,
		}
		logger := logging.NoOpLogger{}
		
		pool := NewHTTPClientPool(config, logger)
		
		require.NotNil(t, pool)
		assert.Equal(t, config, pool.config)
		assert.Equal(t, logger, pool.logger)
		assert.NotNil(t, pool.clients)
	})
	
	t.Run("with nil config", func(t *testing.T) {
		pool := NewHTTPClientPool(nil, nil)
		
		require.NotNil(t, pool)
		assert.Equal(t, DefaultPoolConfig(), pool.config)
		assert.IsType(t, logging.NoOpLogger{}, pool.logger)
	})
	
	t.Run("with nil logger", func(t *testing.T) {
		config := DefaultPoolConfig()
		pool := NewHTTPClientPool(config, nil)
		
		require.NotNil(t, pool)
		assert.IsType(t, logging.NoOpLogger{}, pool.logger)
	})
}

func TestHTTPClientPool_GetClient(t *testing.T) {
	pool := NewHTTPClientPool(nil, nil)
	endpoint := "https://slurm.example.com"
	
	// First call creates client
	client1 := pool.GetClient(endpoint)
	require.NotNil(t, client1)
	
	// Second call returns same client
	client2 := pool.GetClient(endpoint)
	assert.Equal(t, client1, client2)
	
	// Verify stats
	stats := pool.Stats()
	assert.Equal(t, 1, stats.TotalClients)
	require.Contains(t, stats.ClientStats, endpoint)
	
	clientStats := stats.ClientStats[endpoint]
	assert.Equal(t, int64(2), clientStats.UseCount) // Called twice
	assert.True(t, clientStats.Created.Before(time.Now()) || clientStats.Created.Equal(time.Now()))
	assert.True(t, clientStats.LastUsed.Before(time.Now()) || clientStats.LastUsed.Equal(time.Now()))
}

func TestHTTPClientPool_GetClient_DifferentEndpoints(t *testing.T) {
	pool := NewHTTPClientPool(nil, nil)
	
	endpoint1 := "https://slurm1.example.com"
	endpoint2 := "https://slurm2.example.com"
	
	client1 := pool.GetClient(endpoint1)
	client2 := pool.GetClient(endpoint2)
	
	// Should be different clients
	assert.NotEqual(t, client1, client2)
	
	// Verify stats
	stats := pool.Stats()
	assert.Equal(t, 2, stats.TotalClients)
	assert.Contains(t, stats.ClientStats, endpoint1)
	assert.Contains(t, stats.ClientStats, endpoint2)
}

func TestHTTPClientPool_createHTTPClient(t *testing.T) {
	config := &PoolConfig{
		MaxIdleConns:           200,
		MaxIdleConnsPerHost:    20,
		MaxConnsPerHost:        100,
		IdleConnTimeout:        120 * time.Second,
		TLSHandshakeTimeout:    15 * time.Second,
		ExpectContinueTimeout:  2 * time.Second,
		DisableKeepAlives:      true,
		DisableCompression:     true,
		MaxResponseHeaderBytes: 2 << 20, // 2 MB
	}
	
	pool := NewHTTPClientPool(config, nil)
	client := pool.createHTTPClient()
	
	require.NotNil(t, client)
	assert.Equal(t, time.Duration(0), client.Timeout) // No client timeout
	
	// Verify transport configuration
	transport, ok := client.Transport.(*http.Transport)
	require.True(t, ok)
	
	assert.Equal(t, config.MaxIdleConns, transport.MaxIdleConns)
	assert.Equal(t, config.MaxIdleConnsPerHost, transport.MaxIdleConnsPerHost)
	assert.Equal(t, config.MaxConnsPerHost, transport.MaxConnsPerHost)
	assert.Equal(t, config.IdleConnTimeout, transport.IdleConnTimeout)
	assert.Equal(t, config.TLSHandshakeTimeout, transport.TLSHandshakeTimeout)
	assert.Equal(t, config.ExpectContinueTimeout, transport.ExpectContinueTimeout)
	assert.Equal(t, config.DisableKeepAlives, transport.DisableKeepAlives)
	assert.Equal(t, config.DisableCompression, transport.DisableCompression)
	assert.Equal(t, config.MaxResponseHeaderBytes, transport.MaxResponseHeaderBytes)
	assert.True(t, transport.ForceAttemptHTTP2)
	
	// Verify TLS configuration
	require.NotNil(t, transport.TLSClientConfig)
	assert.GreaterOrEqual(t, transport.TLSClientConfig.MinVersion, uint16(0x0303)) // TLS 1.2
}

func TestHTTPClientPool_Stats(t *testing.T) {
	pool := NewHTTPClientPool(nil, nil)
	
	// Initially empty
	stats := pool.Stats()
	assert.Equal(t, 0, stats.TotalClients)
	assert.Empty(t, stats.ClientStats)
	
	// Add some clients
	pool.GetClient("https://endpoint1.com")
	pool.GetClient("https://endpoint2.com")
	pool.GetClient("https://endpoint1.com") // Same endpoint again
	
	stats = pool.Stats()
	assert.Equal(t, 2, stats.TotalClients)
	assert.Len(t, stats.ClientStats, 2)
	
	// Verify endpoint1 was used twice
	stats1 := stats.ClientStats["https://endpoint1.com"]
	assert.Equal(t, int64(2), stats1.UseCount)
	
	// Verify endpoint2 was used once
	stats2 := stats.ClientStats["https://endpoint2.com"]
	assert.Equal(t, int64(1), stats2.UseCount)
}

func TestHTTPClientPool_CleanupIdleClients(t *testing.T) {
	pool := NewHTTPClientPool(nil, nil)
	
	// Add some clients
	client1 := pool.GetClient("https://endpoint1.com")
	client2 := pool.GetClient("https://endpoint2.com")
	
	require.NotNil(t, client1)
	require.NotNil(t, client2)
	
	// Verify both clients exist
	stats := pool.Stats()
	assert.Equal(t, 2, stats.TotalClients)
	
	// Manually set one client as old
	pool.mu.Lock()
	pool.clients["https://endpoint1.com"].lastUsed = time.Now().Add(-1 * time.Hour)
	pool.mu.Unlock()
	
	// Cleanup with 30 minute threshold
	removed := pool.CleanupIdleClients(30 * time.Minute)
	assert.Equal(t, 1, removed)
	
	// Verify only one client remains
	stats = pool.Stats()
	assert.Equal(t, 1, stats.TotalClients)
	assert.Contains(t, stats.ClientStats, "https://endpoint2.com")
	assert.NotContains(t, stats.ClientStats, "https://endpoint1.com")
}

func TestHTTPClientPool_CleanupIdleClients_NoActiveConns(t *testing.T) {
	pool := NewHTTPClientPool(nil, nil)
	
	// Add client
	pool.GetClient("https://endpoint.com")
	
	// Set as old but with active connections
	pool.mu.Lock()
	pool.clients["https://endpoint.com"].lastUsed = time.Now().Add(-1 * time.Hour)
	pool.clients["https://endpoint.com"].activeConns = 5 // Has active connections
	pool.mu.Unlock()
	
	// Should not be removed due to active connections
	removed := pool.CleanupIdleClients(30 * time.Minute)
	assert.Equal(t, 0, removed)
	
	// Verify client still exists
	stats := pool.Stats()
	assert.Equal(t, 1, stats.TotalClients)
}

func TestHTTPClientPool_Close(t *testing.T) {
	pool := NewHTTPClientPool(nil, nil)
	
	// Add some clients
	pool.GetClient("https://endpoint1.com")
	pool.GetClient("https://endpoint2.com")
	
	// Verify clients exist
	stats := pool.Stats()
	assert.Equal(t, 2, stats.TotalClients)
	
	// Close pool
	err := pool.Close()
	assert.NoError(t, err)
	
	// Verify all clients are removed
	stats = pool.Stats()
	assert.Equal(t, 0, stats.TotalClients)
	assert.Empty(t, stats.ClientStats)
}

func TestNewConnectionManager(t *testing.T) {
	pool := NewHTTPClientPool(nil, nil)
	logger := logging.NoOpLogger{}
	
	healthCheck := func(ctx context.Context, endpoint string, client *http.Client) error {
		return nil
	}
	
	cm := NewConnectionManager(pool, healthCheck, logger)
	
	require.NotNil(t, cm)
	assert.Equal(t, pool, cm.pool)
	assert.NotNil(t, cm.healthCheckFunc)
	assert.Equal(t, logger, cm.logger)
	assert.Equal(t, 5*time.Minute, cm.cleanupInterval)
	assert.Equal(t, 15*time.Minute, cm.maxIdleTime)
	assert.NotNil(t, cm.ctx)
	assert.NotNil(t, cm.cancel)
}

func TestNewConnectionManager_NilLogger(t *testing.T) {
	pool := NewHTTPClientPool(nil, nil)
	
	cm := NewConnectionManager(pool, nil, nil)
	
	require.NotNil(t, cm)
	assert.IsType(t, logging.NoOpLogger{}, cm.logger)
}

func TestConnectionManager_StartStop(t *testing.T) {
	pool := NewHTTPClientPool(nil, nil)
	cm := NewConnectionManager(pool, nil, nil)
	
	// Start should not block
	cm.Start()
	
	// Stop should complete quickly
	done := make(chan struct{})
	go func() {
		cm.Stop()
		close(done)
	}()
	
	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Stop() took too long")
	}
}

func TestConnectionManager_GetHealthyClient_Success(t *testing.T) {
	pool := NewHTTPClientPool(nil, nil)
	
	healthCheck := func(ctx context.Context, endpoint string, client *http.Client) error {
		return nil // Always healthy
	}
	
	cm := NewConnectionManager(pool, healthCheck, nil)
	
	ctx := context.Background()
	client, err := cm.GetHealthyClient(ctx, "https://healthy.example.com")
	
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestConnectionManager_GetHealthyClient_HealthCheckFails(t *testing.T) {
	pool := NewHTTPClientPool(nil, nil)
	
	expectedErr := errors.New("endpoint is unhealthy")
	healthCheck := func(ctx context.Context, endpoint string, client *http.Client) error {
		return expectedErr
	}
	
	cm := NewConnectionManager(pool, healthCheck, nil)
	
	ctx := context.Background()
	client, err := cm.GetHealthyClient(ctx, "https://unhealthy.example.com")
	
	assert.Nil(t, client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "endpoint health check failed")
	assert.Contains(t, err.Error(), expectedErr.Error())
}

func TestConnectionManager_GetHealthyClient_NoHealthCheck(t *testing.T) {
	pool := NewHTTPClientPool(nil, nil)
	cm := NewConnectionManager(pool, nil, nil) // No health check
	
	ctx := context.Background()
	client, err := cm.GetHealthyClient(ctx, "https://example.com")
	
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestConnectionManager_CleanupRoutine(t *testing.T) {
	pool := NewHTTPClientPool(nil, nil)
	
	// Create connection manager with very short cleanup interval
	cm := NewConnectionManager(pool, nil, nil)
	cm.cleanupInterval = 10 * time.Millisecond
	cm.maxIdleTime = 5 * time.Millisecond
	
	// Add a client
	pool.GetClient("https://example.com")
	
	// Verify it exists
	stats := pool.Stats()
	assert.Equal(t, 1, stats.TotalClients)
	
	// Start cleanup routine
	cm.Start()
	
	// Wait for cleanup to run
	time.Sleep(50 * time.Millisecond)
	
	// Stop the routine
	cm.Stop()
	
	// The client should have been cleaned up
	stats = pool.Stats()
	assert.Equal(t, 0, stats.TotalClients)
}

func TestPooledClient(t *testing.T) {
	client := &http.Client{}
	now := time.Now()
	
	pc := &pooledClient{
		client:      client,
		created:     now,
		lastUsed:    now,
		useCount:    5,
		activeConns: 2,
	}
	
	assert.Equal(t, client, pc.client)
	assert.Equal(t, now, pc.created)
	assert.Equal(t, now, pc.lastUsed)
	assert.Equal(t, int64(5), pc.useCount)
	assert.Equal(t, int32(2), pc.activeConns)
}

func TestPoolConfig_CustomValues(t *testing.T) {
	config := &PoolConfig{
		MaxIdleConns:           200,
		MaxIdleConnsPerHost:    20,
		MaxConnsPerHost:        100,
		IdleConnTimeout:        120 * time.Second,
		TLSHandshakeTimeout:    15 * time.Second,
		ExpectContinueTimeout:  2 * time.Second,
		DisableKeepAlives:      true,
		DisableCompression:     true,
		MaxResponseHeaderBytes: 2 << 20,
	}
	
	assert.Equal(t, 200, config.MaxIdleConns)
	assert.Equal(t, 20, config.MaxIdleConnsPerHost)
	assert.Equal(t, 100, config.MaxConnsPerHost)
	assert.Equal(t, 120*time.Second, config.IdleConnTimeout)
	assert.Equal(t, 15*time.Second, config.TLSHandshakeTimeout)
	assert.Equal(t, 2*time.Second, config.ExpectContinueTimeout)
	assert.True(t, config.DisableKeepAlives)
	assert.True(t, config.DisableCompression)
	assert.Equal(t, int64(2<<20), config.MaxResponseHeaderBytes)
}

func TestClientStats(t *testing.T) {
	now := time.Now()
	stats := ClientStats{
		Created:     now,
		LastUsed:    now,
		UseCount:    10,
		ActiveConns: 3,
	}
	
	assert.Equal(t, now, stats.Created)
	assert.Equal(t, now, stats.LastUsed)
	assert.Equal(t, int64(10), stats.UseCount)
	assert.Equal(t, int32(3), stats.ActiveConns)
}

func TestPoolStats(t *testing.T) {
	stats := PoolStats{
		TotalClients: 5,
		ClientStats: map[string]ClientStats{
			"endpoint1": {UseCount: 10},
			"endpoint2": {UseCount: 20},
		},
	}
	
	assert.Equal(t, 5, stats.TotalClients)
	assert.Len(t, stats.ClientStats, 2)
	assert.Equal(t, int64(10), stats.ClientStats["endpoint1"].UseCount)
	assert.Equal(t, int64(20), stats.ClientStats["endpoint2"].UseCount)
}

func TestHealthCheckFunc(t *testing.T) {
	// Test that HealthCheckFunc type works as expected
	healthCheck := func(ctx context.Context, endpoint string, client *http.Client) error {
		if endpoint == "https://bad.example.com" {
			return errors.New("bad endpoint")
		}
		return nil
	}
	
	ctx := context.Background()
	client := &http.Client{}
	
	// Good endpoint
	err := healthCheck(ctx, "https://good.example.com", client)
	assert.NoError(t, err)
	
	// Bad endpoint
	err = healthCheck(ctx, "https://bad.example.com", client)
	assert.Error(t, err)
	assert.Equal(t, "bad endpoint", err.Error())
}

func TestHTTPClientPool_ConcurrentAccess(t *testing.T) {
	pool := NewHTTPClientPool(nil, nil)
	endpoint := "https://concurrent.example.com"
	
	// Concurrent access should be safe
	const numGoroutines = 10
	clients := make([]*http.Client, numGoroutines)
	done := make(chan int, numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			clients[index] = pool.GetClient(endpoint)
			done <- index
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	
	// All clients should be the same instance
	for i := 1; i < numGoroutines; i++ {
		assert.Equal(t, clients[0], clients[i])
	}
	
	// Verify stats
	stats := pool.Stats()
	assert.Equal(t, 1, stats.TotalClients)
	assert.Equal(t, int64(numGoroutines), stats.ClientStats[endpoint].UseCount)
}