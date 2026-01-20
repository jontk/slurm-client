// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package factory

import (
	"net/http"
	"time"

	slurmctx "github.com/jontk/slurm-client/pkg/context"
	"github.com/jontk/slurm-client/pkg/logging"
	"github.com/jontk/slurm-client/pkg/metrics"
	"github.com/jontk/slurm-client/pkg/middleware"
	"github.com/jontk/slurm-client/pkg/pool"
	"github.com/jontk/slurm-client/pkg/retry"
)

// EnhancedOptions holds the new configuration options
type EnhancedOptions struct {
	// Logging
	Logger   logging.Logger
	LogLevel string

	// Metrics
	MetricsCollector metrics.Collector

	// Middleware
	Middlewares []middleware.Middleware

	// Timeouts
	DefaultTimeout time.Duration
	TimeoutConfig  *slurmctx.TimeoutConfig

	// Connection pooling
	ConnectionPool *pool.HTTPClientPool
	PoolConfig     *pool.PoolConfig

	// Retry
	RetryBackoff retry.Policy
	MaxRetries   int

	// HTTP options
	UserAgent      string
	RequestIDGen   func() string
	CircuitBreaker *circuitBreakerConfig
	Compression    *bool
	KeepAlive      *bool

	// Debug mode
	Debug bool
}

type circuitBreakerConfig struct {
	Threshold int
	Timeout   time.Duration
}

// WithLogger sets a custom logger for the client
func (f *ClientFactory) WithLogger(logger logging.Logger) error {
	if f.enhanced == nil {
		f.enhanced = &EnhancedOptions{}
	}
	f.enhanced.Logger = logger
	return nil
}

// WithLogLevel sets the logging level
func (f *ClientFactory) WithLogLevel(level string) error {
	if f.enhanced == nil {
		f.enhanced = &EnhancedOptions{}
	}
	f.enhanced.LogLevel = level
	return nil
}

// WithMetricsCollector sets a custom metrics collector
func (f *ClientFactory) WithMetricsCollector(collector metrics.Collector) error {
	if f.enhanced == nil {
		f.enhanced = &EnhancedOptions{}
	}
	f.enhanced.MetricsCollector = collector
	return nil
}

// WithMiddleware adds custom middleware to the HTTP client
func (f *ClientFactory) WithMiddleware(middlewares ...middleware.Middleware) error {
	if f.enhanced == nil {
		f.enhanced = &EnhancedOptions{}
	}
	f.enhanced.Middlewares = append(f.enhanced.Middlewares, middlewares...)
	return nil
}

// WithTimeout sets default timeout for all operations
func (f *ClientFactory) WithTimeout(timeout time.Duration) error {
	if f.enhanced == nil {
		f.enhanced = &EnhancedOptions{}
	}
	f.enhanced.DefaultTimeout = timeout
	return nil
}

// WithTimeoutConfig sets custom timeout configuration
func (f *ClientFactory) WithTimeoutConfig(config *slurmctx.TimeoutConfig) error {
	if f.enhanced == nil {
		f.enhanced = &EnhancedOptions{}
	}
	f.enhanced.TimeoutConfig = config
	return nil
}

// WithConnectionPool enables connection pooling
func (f *ClientFactory) WithConnectionPool(config *pool.PoolConfig) error {
	if f.enhanced == nil {
		f.enhanced = &EnhancedOptions{}
	}
	f.enhanced.PoolConfig = config
	return nil
}

// WithRetryBackoff sets a custom retry backoff strategy
func (f *ClientFactory) WithRetryBackoff(backoff retry.Policy) error {
	if f.enhanced == nil {
		f.enhanced = &EnhancedOptions{}
	}
	f.enhanced.RetryBackoff = backoff
	return nil
}

// WithMaxRetries sets the maximum number of retry attempts
func (f *ClientFactory) WithMaxRetries(maxRetries int) error {
	if f.enhanced == nil {
		f.enhanced = &EnhancedOptions{}
	}
	f.enhanced.MaxRetries = maxRetries
	return nil
}

// WithUserAgent sets a custom User-Agent header
func (f *ClientFactory) WithUserAgent(userAgent string) error {
	if f.enhanced == nil {
		f.enhanced = &EnhancedOptions{}
	}
	f.enhanced.UserAgent = userAgent
	return nil
}

// WithRequestID enables request ID generation
func (f *ClientFactory) WithRequestID(generator func() string) error {
	if f.enhanced == nil {
		f.enhanced = &EnhancedOptions{}
	}
	f.enhanced.RequestIDGen = generator
	return nil
}

// WithCircuitBreaker enables circuit breaker functionality
func (f *ClientFactory) WithCircuitBreaker(threshold int, timeout time.Duration) error {
	if f.enhanced == nil {
		f.enhanced = &EnhancedOptions{}
	}
	f.enhanced.CircuitBreaker = &circuitBreakerConfig{
		Threshold: threshold,
		Timeout:   timeout,
	}
	return nil
}

// WithCompression enables or disables HTTP compression
func (f *ClientFactory) WithCompression(enabled bool) error {
	if f.enhanced == nil {
		f.enhanced = &EnhancedOptions{}
	}
	f.enhanced.Compression = &enabled
	return nil
}

// WithKeepAlive enables or disables HTTP keep-alive
func (f *ClientFactory) WithKeepAlive(enabled bool) error {
	if f.enhanced == nil {
		f.enhanced = &EnhancedOptions{}
	}
	f.enhanced.KeepAlive = &enabled
	return nil
}

// WithDebug enables debug mode
func (f *ClientFactory) WithDebug() error {
	if f.enhanced == nil {
		f.enhanced = &EnhancedOptions{}
	}
	f.enhanced.Debug = true
	return nil
}

// buildEnhancedHTTPClient builds an HTTP client with all enhancements
func (f *ClientFactory) buildEnhancedHTTPClient() *http.Client {
	// Start with base client or pooled client
	var baseClient *http.Client

	// Use connection pool if configured
	if f.enhanced != nil && f.enhanced.PoolConfig != nil {
		logger := f.enhanced.Logger
		if logger == nil {
			logger = logging.NoOpLogger{}
		}

		pool := pool.NewHTTPClientPool(f.enhanced.PoolConfig, logger)
		f.enhanced.ConnectionPool = pool
		baseClient = pool.GetClient(f.baseURL)
	} else if f.httpClient != nil {
		baseClient = f.httpClient
	} else {
		baseClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	// Apply middleware if configured
	if f.enhanced != nil && len(f.enhanced.Middlewares) > 0 {
		transport := baseClient.Transport
		if transport == nil {
			transport = http.DefaultTransport
		}

		// Build middleware chain
		middlewares := f.buildMiddlewareChain()

		// Apply middleware
		for i := len(middlewares) - 1; i >= 0; i-- {
			transport = middlewares[i](transport)
		}

		baseClient.Transport = transport
	}

	return baseClient
}

// buildMiddlewareChain builds the complete middleware chain
func (f *ClientFactory) buildMiddlewareChain() []middleware.Middleware {
	var middlewares []middleware.Middleware

	if f.enhanced == nil {
		return middlewares
	}

	// Add timeout middleware
	if f.enhanced.DefaultTimeout > 0 {
		middlewares = append(middlewares, middleware.WithTimeout(f.enhanced.DefaultTimeout))
	}

	// Add logging middleware
	if f.enhanced.Logger != nil {
		middlewares = append(middlewares, middleware.WithLogging(f.enhanced.Logger))
	}

	// Add metrics middleware
	if f.enhanced.MetricsCollector != nil {
		middlewares = append(middlewares, middleware.WithMetrics(f.enhanced.MetricsCollector))
	}

	// Add retry middleware
	if f.enhanced.RetryBackoff != nil || f.enhanced.MaxRetries > 0 {
		maxRetries := f.enhanced.MaxRetries
		if maxRetries == 0 {
			maxRetries = 3 // default
		}
		middlewares = append(middlewares, middleware.WithRetry(maxRetries, middleware.DefaultShouldRetry))
	}

	// Add circuit breaker
	if f.enhanced.CircuitBreaker != nil {
		middlewares = append(middlewares, middleware.WithCircuitBreaker(
			f.enhanced.CircuitBreaker.Threshold,
			f.enhanced.CircuitBreaker.Timeout,
		))
	}

	// Add request ID
	if f.enhanced.RequestIDGen != nil {
		middlewares = append(middlewares, middleware.WithRequestID(f.enhanced.RequestIDGen))
	}

	// Add user agent
	if f.enhanced.UserAgent != "" {
		middlewares = append(middlewares, middleware.WithUserAgent(f.enhanced.UserAgent))
	}

	// Add user-provided middleware
	middlewares = append(middlewares, f.enhanced.Middlewares...)

	return middlewares
}

// GetEnhancedOptions returns the enhanced options for use by implementations
func (f *ClientFactory) GetEnhancedOptions() *EnhancedOptions {
	return f.enhanced
}
