// Package slurm provides client options for configuring the SLURM REST API client
package slurm

import (
	"net/http"
	"time"

	"github.com/jontk/slurm-client/internal/factory"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
	slurmctx "github.com/jontk/slurm-client/pkg/context"
	"github.com/jontk/slurm-client/pkg/logging"
	"github.com/jontk/slurm-client/pkg/metrics"
	"github.com/jontk/slurm-client/pkg/middleware"
	"github.com/jontk/slurm-client/pkg/pool"
	"github.com/jontk/slurm-client/pkg/retry"
)

// WithConfig sets the client configuration from a config file or environment
func WithConfig(cfg *config.Config) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithConfig(cfg)
	}
}

// WithBaseURL sets the base URL for the SLURM REST API
func WithBaseURL(baseURL string) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithBaseURL(baseURL)
	}
}

// WithToken sets the authentication token
func WithToken(token string) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithToken(token)
	}
}

// WithUserToken sets user authentication with username and token
func WithUserToken(username, token string) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithUserToken(username, token)
	}
}

// WithTokenSource sets a custom token source for authentication
func WithTokenSource(tokenSource auth.TokenSource) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithTokenSource(tokenSource)
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithHTTPClient(client)
	}
}

// WithRetryPolicy sets a custom retry policy
func WithRetryPolicy(policy retry.Policy) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithRetryPolicy(policy)
	}
}

// WithVersion forces a specific API version
func WithVersion(version string) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithVersion(version)
	}
}

// WithNoAuth disables authentication (for testing or public endpoints)
func WithNoAuth() ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithNoAuth()
	}
}

// WithTLSConfig sets custom TLS configuration
func WithTLSConfig(tlsConfig *http.Transport) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithTLSConfig(tlsConfig)
	}
}

// === New options for the improved features ===

// WithLogger sets a custom logger for the client
func WithLogger(logger logging.Logger) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithLogger(logger)
	}
}

// WithLogLevel sets the logging level
func WithLogLevel(level string) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithLogLevel(level)
	}
}

// WithMetricsCollector sets a custom metrics collector
func WithMetricsCollector(collector metrics.Collector) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithMetricsCollector(collector)
	}
}

// WithMiddleware adds custom middleware to the HTTP client
func WithMiddleware(middlewares ...middleware.Middleware) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithMiddleware(middlewares...)
	}
}

// WithTimeout sets default timeout for all operations
func WithTimeout(timeout time.Duration) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithTimeout(timeout)
	}
}

// WithTimeoutConfig sets custom timeout configuration for different operation types
func WithTimeoutConfig(config *slurmctx.TimeoutConfig) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithTimeoutConfig(config)
	}
}

// WithConnectionPool enables connection pooling with the specified configuration
func WithConnectionPool(config *pool.PoolConfig) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithConnectionPool(config)
	}
}

// WithRetryBackoff sets a custom retry backoff strategy
func WithRetryBackoff(backoff retry.BackoffStrategy) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithRetryBackoff(backoff)
	}
}

// WithMaxRetries sets the maximum number of retry attempts
func WithMaxRetries(maxRetries int) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithMaxRetries(maxRetries)
	}
}

// WithUserAgent sets a custom User-Agent header
func WithUserAgent(userAgent string) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithUserAgent(userAgent)
	}
}

// WithRequestID enables request ID generation with the provided generator
func WithRequestID(generator func() string) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithRequestID(generator)
	}
}

// WithCircuitBreaker enables circuit breaker functionality
func WithCircuitBreaker(threshold int, timeout time.Duration) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithCircuitBreaker(threshold, timeout)
	}
}

// WithCompression enables or disables HTTP compression
func WithCompression(enabled bool) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithCompression(enabled)
	}
}

// WithKeepAlive enables or disables HTTP keep-alive
func WithKeepAlive(enabled bool) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithKeepAlive(enabled)
	}
}

// WithDebug enables debug mode with verbose logging
func WithDebug() ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithDebug()
	}
}