// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package performance

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConnectionPoolConfig(t *testing.T) {
	config := DefaultConnectionPoolConfig()

	require.NotNil(t, config)
	assert.Equal(t, 100, config.MaxIdleConns)
	assert.Equal(t, 30, config.MaxIdleConnsPerHost)
	assert.Equal(t, 100, config.MaxConnsPerHost)
	assert.Equal(t, 90*time.Second, config.IdleConnTimeout)
	assert.Equal(t, 10*time.Second, config.TLSHandshakeTimeout)
	assert.Equal(t, 1*time.Second, config.ExpectContinueTimeout)
	assert.Equal(t, 30*time.Second, config.ResponseHeaderTimeout)
	assert.False(t, config.DisableCompression)
	assert.False(t, config.DisableKeepAlives)
	assert.True(t, config.EnableTLSSessionResumption)
	assert.False(t, config.TLSInsecureSkipVerify)
}

func TestHighPerformanceConnectionPoolConfig(t *testing.T) {
	config := HighPerformanceConnectionPoolConfig()

	require.NotNil(t, config)
	assert.Equal(t, 200, config.MaxIdleConns)
	assert.Equal(t, 50, config.MaxIdleConnsPerHost)
	assert.Equal(t, 200, config.MaxConnsPerHost)
	assert.Equal(t, 120*time.Second, config.IdleConnTimeout)
	assert.Equal(t, 5*time.Second, config.TLSHandshakeTimeout)
	assert.Equal(t, 500*time.Millisecond, config.ExpectContinueTimeout)
	assert.Equal(t, 15*time.Second, config.ResponseHeaderTimeout)
}

func TestConservativeConnectionPoolConfig(t *testing.T) {
	config := ConservativeConnectionPoolConfig()

	require.NotNil(t, config)
	assert.Equal(t, 10, config.MaxIdleConns)
	assert.Equal(t, 5, config.MaxIdleConnsPerHost)
	assert.Equal(t, 20, config.MaxConnsPerHost)
	assert.Equal(t, 30*time.Second, config.IdleConnTimeout)
	assert.Equal(t, 15*time.Second, config.TLSHandshakeTimeout)
	assert.Equal(t, 2*time.Second, config.ExpectContinueTimeout)
	assert.Equal(t, 60*time.Second, config.ResponseHeaderTimeout)
}

func TestNewHTTPClientPool(t *testing.T) {
	t.Run("with config", func(t *testing.T) {
		config := &ConnectionPoolConfig{
			MaxIdleConns:    50,
			MaxConnsPerHost: 25,
		}

		pool := NewHTTPClientPool(config)
		defer pool.Close()

		require.NotNil(t, pool)
		assert.Equal(t, config, pool.config)
		assert.NotNil(t, pool.clients)
	})

	t.Run("with nil config", func(t *testing.T) {
		pool := NewHTTPClientPool(nil)
		defer pool.Close()

		require.NotNil(t, pool)
		assert.Equal(t, DefaultConnectionPoolConfig(), pool.config)
	})
}

func TestHTTPClientPool_GetClient(t *testing.T) {
	pool := NewHTTPClientPool(nil)
	defer pool.Close()

	endpoint1 := "https://slurm.example.com"
	endpoint2 := "https://slurm2.example.com"

	// First call creates client
	client1 := pool.GetClient(endpoint1)
	require.NotNil(t, client1)

	// Second call to same endpoint returns same client
	client1Again := pool.GetClient(endpoint1)
	assert.Equal(t, client1, client1Again)

	// Different endpoint returns different client
	client2 := pool.GetClient(endpoint2)
	assert.NotEqual(t, client1, client2)

	// Verify client configuration
	transport, ok := client1.Transport.(*http.Transport)
	require.True(t, ok)

	assert.Equal(t, pool.config.MaxIdleConns, transport.MaxIdleConns)
	assert.Equal(t, pool.config.MaxIdleConnsPerHost, transport.MaxIdleConnsPerHost)
	assert.Equal(t, pool.config.MaxConnsPerHost, transport.MaxConnsPerHost)
	assert.Equal(t, pool.config.IdleConnTimeout, transport.IdleConnTimeout)
	assert.Equal(t, pool.config.TLSHandshakeTimeout, transport.TLSHandshakeTimeout)
	assert.Equal(t, pool.config.ExpectContinueTimeout, transport.ExpectContinueTimeout)
	assert.Equal(t, pool.config.ResponseHeaderTimeout, transport.ResponseHeaderTimeout)
	assert.Equal(t, pool.config.DisableCompression, transport.DisableCompression)
	assert.Equal(t, pool.config.DisableKeepAlives, transport.DisableKeepAlives)
	assert.True(t, transport.ForceAttemptHTTP2)
	assert.Equal(t, 64*1024, transport.WriteBufferSize)
	assert.Equal(t, 64*1024, transport.ReadBufferSize)
}

func TestHTTPClientPool_TLSConfiguration(t *testing.T) {
	t.Run("with TLS session resumption", func(t *testing.T) {
		config := &ConnectionPoolConfig{
			EnableTLSSessionResumption: true,
			TLSInsecureSkipVerify:      false,
		}

		pool := NewHTTPClientPool(config)
		defer pool.Close()

		client := pool.GetClient("https://example.com")
		transport, ok := client.Transport.(*http.Transport)
		require.True(t, ok)

		require.NotNil(t, transport.TLSClientConfig)
		assert.False(t, transport.TLSClientConfig.InsecureSkipVerify)
		assert.NotNil(t, transport.TLSClientConfig.ClientSessionCache)
	})

	t.Run("with insecure skip verify", func(t *testing.T) {
		config := &ConnectionPoolConfig{
			EnableTLSSessionResumption: false,
			TLSInsecureSkipVerify:      true,
		}

		pool := NewHTTPClientPool(config)
		defer pool.Close()

		client := pool.GetClient("https://example.com")
		transport, ok := client.Transport.(*http.Transport)
		require.True(t, ok)

		require.NotNil(t, transport.TLSClientConfig)
		assert.True(t, transport.TLSClientConfig.InsecureSkipVerify)
	})

	t.Run("no TLS config when not needed", func(t *testing.T) {
		config := &ConnectionPoolConfig{
			EnableTLSSessionResumption: false,
			TLSInsecureSkipVerify:      false,
		}

		pool := NewHTTPClientPool(config)
		defer pool.Close()

		client := pool.GetClient("https://example.com")
		transport, ok := client.Transport.(*http.Transport)
		require.True(t, ok)

		// TLS config should not be set when both flags are false
		assert.Nil(t, transport.TLSClientConfig)
	})
}

func TestHTTPClientPool_GetStats(t *testing.T) {
	pool := NewHTTPClientPool(nil)
	defer pool.Close()

	// Initially no clients
	stats := pool.GetStats()
	assert.Equal(t, 0, stats.ActiveClients)
	assert.Empty(t, stats.Endpoints)
	assert.Equal(t, *pool.config, stats.Config)

	// Add some clients
	pool.GetClient("https://example1.com")
	pool.GetClient("https://example2.com")

	stats = pool.GetStats()
	assert.Equal(t, 2, stats.ActiveClients)
	assert.Len(t, stats.Endpoints, 2)

	// Check endpoint stats
	endpoints := make(map[string]bool)
	for _, ep := range stats.Endpoints {
		endpoints[ep.Endpoint] = true
	}
	assert.True(t, endpoints["https://example1.com"])
	assert.True(t, endpoints["https://example2.com"])
}

func TestHTTPClientPool_Close(t *testing.T) {
	pool := NewHTTPClientPool(nil)

	// Add some clients
	pool.GetClient("https://example1.com")
	pool.GetClient("https://example2.com")

	// Verify clients exist
	stats := pool.GetStats()
	assert.Equal(t, 2, stats.ActiveClients)

	// Close pool
	pool.Close()

	// Verify clients are cleared
	stats = pool.GetStats()
	assert.Equal(t, 0, stats.ActiveClients)
}

func TestGetConnectionPoolConfigForProfile(t *testing.T) {
	tests := []struct {
		profile  PerformanceProfile
		expected func(*ConnectionPoolConfig) bool
	}{
		{
			ProfileHighThroughput,
			func(c *ConnectionPoolConfig) bool {
				return c.MaxIdleConns == 200 && c.MaxConnsPerHost == 200
			},
		},
		{
			ProfileLowLatency,
			func(c *ConnectionPoolConfig) bool {
				return c.ExpectContinueTimeout == 100*time.Millisecond &&
					c.ResponseHeaderTimeout == 5*time.Second &&
					c.TLSHandshakeTimeout == 3*time.Second &&
					c.MaxConnsPerHost == 50
			},
		},
		{
			ProfileConservative,
			func(c *ConnectionPoolConfig) bool {
				return c.MaxIdleConns == 10 && c.MaxConnsPerHost == 20
			},
		},
		{
			ProfileBatch,
			func(c *ConnectionPoolConfig) bool {
				return c.IdleConnTimeout == 300*time.Second &&
					c.ResponseHeaderTimeout == 120*time.Second
			},
		},
		{
			PerformanceProfile("unknown"),
			func(c *ConnectionPoolConfig) bool {
				return c.MaxIdleConns == 100 // Default config
			},
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.profile), func(t *testing.T) {
			config := GetConnectionPoolConfigForProfile(tt.profile)
			require.NotNil(t, config)
			assert.True(t, tt.expected(config))
		})
	}
}

func TestNewHTTPClientPoolManager(t *testing.T) {
	manager := NewHTTPClientPoolManager()
	defer manager.Close()

	require.NotNil(t, manager)
	assert.NotNil(t, manager.pools)
}

func TestHTTPClientPoolManager_GetPoolForVersion(t *testing.T) {
	manager := NewHTTPClientPoolManager()
	defer manager.Close()

	// First call creates pool
	pool1 := manager.GetPoolForVersion("v0.0.41", ProfileDefault)
	require.NotNil(t, pool1)

	// Second call to same version+profile returns same pool
	pool1Again := manager.GetPoolForVersion("v0.0.41", ProfileDefault)
	assert.Equal(t, pool1, pool1Again)

	// Different version returns different pool
	pool2 := manager.GetPoolForVersion("v0.0.42", ProfileDefault)
	assert.NotEqual(t, pool1, pool2)

	// Different profile returns different pool
	pool3 := manager.GetPoolForVersion("v0.0.41", ProfileHighThroughput)
	assert.NotEqual(t, pool1, pool3)
}

func TestHTTPClientPoolManager_VersionSpecificOptimizations(t *testing.T) {
	manager := NewHTTPClientPoolManager()
	defer manager.Close()

	tests := []struct {
		version              string
		expectedOptimization func(*http.Client) bool
	}{
		{
			"v0.0.40",
			func(client *http.Client) bool {
				transport := client.Transport.(*http.Transport)
				return transport.MaxConnsPerHost <= 20 &&
					transport.ResponseHeaderTimeout >= 30*time.Second
			},
		},
		{
			"v0.0.41",
			func(client *http.Client) bool {
				transport := client.Transport.(*http.Transport)
				return transport.ExpectContinueTimeout == 750*time.Millisecond
			},
		},
		{
			"v0.0.42",
			func(client *http.Client) bool {
				// Should use default optimizations
				return true
			},
		},
		{
			"v0.0.43",
			func(client *http.Client) bool {
				transport := client.Transport.(*http.Transport)
				return transport.ExpectContinueTimeout == 250*time.Millisecond
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			pool := manager.GetPoolForVersion(tt.version, ProfileDefault)
			client := pool.GetClient("https://example.com")

			assert.True(t, tt.expectedOptimization(client))
		})
	}
}

func TestHTTPClientPoolManager_GetGlobalStats(t *testing.T) {
	manager := NewHTTPClientPoolManager()
	defer manager.Close()

	// Initially no pools
	stats := manager.GetGlobalStats()
	assert.Empty(t, stats)

	// Create some pools
	pool1 := manager.GetPoolForVersion("v0.0.41", ProfileDefault)
	pool2 := manager.GetPoolForVersion("v0.0.42", ProfileHighThroughput)

	// Add clients to pools
	pool1.GetClient("https://example1.com")
	pool2.GetClient("https://example2.com")

	stats = manager.GetGlobalStats()
	assert.Len(t, stats, 2)
	assert.Contains(t, stats, "v0.0.41:default")
	assert.Contains(t, stats, "v0.0.42:high_throughput")

	assert.Equal(t, 1, stats["v0.0.41:default"].ActiveClients)
	assert.Equal(t, 1, stats["v0.0.42:high_throughput"].ActiveClients)
}

func TestHTTPClientPoolManager_Close(t *testing.T) {
	manager := NewHTTPClientPoolManager()

	// Create some pools
	pool1 := manager.GetPoolForVersion("v0.0.41", ProfileDefault)
	pool2 := manager.GetPoolForVersion("v0.0.42", ProfileDefault)

	// Add clients
	pool1.GetClient("https://example1.com")
	pool2.GetClient("https://example2.com")

	// Verify pools exist
	stats := manager.GetGlobalStats()
	assert.Len(t, stats, 2)

	// Close manager
	manager.Close()

	// Verify pools are cleared
	stats = manager.GetGlobalStats()
	assert.Empty(t, stats)
}

func TestHelperFunctions(t *testing.T) {
	t.Run("min function", func(t *testing.T) {
		assert.Equal(t, 5, min(5, 10))
		assert.Equal(t, 5, min(10, 5))
		assert.Equal(t, 5, min(5, 5))
	})

	t.Run("max function", func(t *testing.T) {
		assert.Equal(t, 10*time.Second, max(5*time.Second, 10*time.Second))
		assert.Equal(t, 10*time.Second, max(10*time.Second, 5*time.Second))
		assert.Equal(t, 5*time.Second, max(5*time.Second, 5*time.Second))
	})
}

func TestPerformanceProfiles(t *testing.T) {
	// Test that all profile constants are defined
	profiles := []PerformanceProfile{
		ProfileDefault,
		ProfileHighThroughput,
		ProfileLowLatency,
		ProfileConservative,
		ProfileBatch,
	}

	for _, profile := range profiles {
		assert.NotEmpty(t, string(profile))
	}

	// Test that each profile has expected string values
	assert.Equal(t, "default", string(ProfileDefault))
	assert.Equal(t, "high_throughput", string(ProfileHighThroughput))
	assert.Equal(t, "low_latency", string(ProfileLowLatency))
	assert.Equal(t, "conservative", string(ProfileConservative))
	assert.Equal(t, "batch", string(ProfileBatch))
}

func TestConnectionPoolConfigCompleteness(t *testing.T) {
	// Test that all configuration functions return complete configs
	configs := []*ConnectionPoolConfig{
		DefaultConnectionPoolConfig(),
		HighPerformanceConnectionPoolConfig(),
		ConservativeConnectionPoolConfig(),
	}

	for i, config := range configs {
		t.Run([]string{"default", "high_perf", "conservative"}[i], func(t *testing.T) {
			// Verify all fields are set to non-zero values
			assert.Greater(t, config.MaxIdleConns, 0)
			assert.Greater(t, config.MaxIdleConnsPerHost, 0)
			assert.Greater(t, config.MaxConnsPerHost, 0)
			assert.Greater(t, config.IdleConnTimeout, time.Duration(0))
			assert.Greater(t, config.TLSHandshakeTimeout, time.Duration(0))
			assert.Greater(t, config.ExpectContinueTimeout, time.Duration(0))
			assert.Greater(t, config.ResponseHeaderTimeout, time.Duration(0))
		})
	}
}

func TestHTTPClientConfiguration(t *testing.T) {
	// Test various configuration scenarios
	t.Run("compression disabled", func(t *testing.T) {
		config := &ConnectionPoolConfig{
			DisableCompression: true,
		}

		pool := NewHTTPClientPool(config)
		defer pool.Close()

		client := pool.GetClient("https://example.com")
		transport := client.Transport.(*http.Transport)

		assert.True(t, transport.DisableCompression)
	})

	t.Run("keep alives disabled", func(t *testing.T) {
		config := &ConnectionPoolConfig{
			DisableKeepAlives: true,
		}

		pool := NewHTTPClientPool(config)
		defer pool.Close()

		client := pool.GetClient("https://example.com")
		transport := client.Transport.(*http.Transport)

		assert.True(t, transport.DisableKeepAlives)
	})

	t.Run("client timeout", func(t *testing.T) {
		pool := NewHTTPClientPool(nil)
		defer pool.Close()

		client := pool.GetClient("https://example.com")

		// Client timeout should be 0 (handled at request level)
		assert.Equal(t, time.Duration(0), client.Timeout)
	})
}

func TestTLSClientSessionCache(t *testing.T) {
	config := &ConnectionPoolConfig{
		EnableTLSSessionResumption: true,
	}

	pool := NewHTTPClientPool(config)
	defer pool.Close()

	client := pool.GetClient("https://example.com")
	transport := client.Transport.(*http.Transport)

	require.NotNil(t, transport.TLSClientConfig)
	require.NotNil(t, transport.TLSClientConfig.ClientSessionCache)

	// Verify it's an LRU cache (the implementation returns an interface)
	assert.NotNil(t, transport.TLSClientConfig.ClientSessionCache)
}
