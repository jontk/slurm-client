// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Package slurm provides client options for configuring the SLURM REST API client
package slurm

import (
	"net/http"
	"time"

	"github.com/jontk/slurm-client/internal/factory"
	slurmctx "github.com/jontk/slurm-client/pkg/context"
	"github.com/jontk/slurm-client/pkg/logging"
	"github.com/jontk/slurm-client/pkg/metrics"
	"github.com/jontk/slurm-client/pkg/middleware"
	"github.com/jontk/slurm-client/pkg/pool"
	"github.com/jontk/slurm-client/pkg/retry"
)

// === Deprecated Options ===
// These options are superseded by the enhanced options API.
// They still work but may be removed in a future version.
// Consider using WithHTTPClient with a custom transport for advanced configuration.

// Deprecated: WithLogger is superseded by the enhanced options API. It still works but may be removed in a future version.
func WithLogger(logger logging.Logger) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithLogger(logger)
	}
}

// Deprecated: WithLogLevel is superseded by the enhanced options API. It still works but may be removed in a future version.
func WithLogLevel(level string) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithLogLevel(level)
	}
}

// Deprecated: WithMetricsCollector is superseded by the enhanced options API. It still works but may be removed in a future version.
func WithMetricsCollector(collector metrics.Collector) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithMetricsCollector(collector)
	}
}

// Deprecated: WithMiddleware is superseded by the enhanced options API. It still works but may be removed in a future version.
func WithMiddleware(middlewares ...middleware.Middleware) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithMiddleware(middlewares...)
	}
}

// Deprecated: WithTimeoutConfig is superseded by the enhanced options API. It still works but may be removed in a future version.
func WithTimeoutConfig(config *slurmctx.TimeoutConfig) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithTimeoutConfig(config)
	}
}

// Deprecated: WithConnectionPool is superseded by the enhanced options API. It still works but may be removed in a future version.
func WithConnectionPool(config *pool.PoolConfig) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithConnectionPool(config)
	}
}

// Deprecated: WithRetryBackoff is superseded by the enhanced options API. It still works but may be removed in a future version.
func WithRetryBackoff(backoff retry.Policy) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithRetryBackoff(backoff)
	}
}

// Deprecated: WithMaxRetries is superseded by the enhanced options API. It still works but may be removed in a future version.
func WithMaxRetries(maxRetries int) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithMaxRetries(maxRetries)
	}
}

// Deprecated: WithUserAgent is superseded by the enhanced options API. It still works but may be removed in a future version.
func WithUserAgent(userAgent string) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithUserAgent(userAgent)
	}
}

// Deprecated: WithRequestID is superseded by the enhanced options API. It still works but may be removed in a future version.
func WithRequestID(generator func() string) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithRequestID(generator)
	}
}

// Deprecated: WithCircuitBreaker is superseded by the enhanced options API. It still works but may be removed in a future version.
func WithCircuitBreaker(threshold int, timeout time.Duration) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithCircuitBreaker(threshold, timeout)
	}
}

// Deprecated: WithCompression is superseded by the enhanced options API. It still works but may be removed in a future version.
func WithCompression(enabled bool) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithCompression(enabled)
	}
}

// Deprecated: WithKeepAlive is superseded by the enhanced options API. It still works but may be removed in a future version.
func WithKeepAlive(enabled bool) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithKeepAlive(enabled)
	}
}

// Deprecated: WithDebug is superseded by the enhanced options API. It still works but may be removed in a future version.
func WithDebug() ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.WithDebug()
	}
}

// Deprecated: WithTLSConfig is not functional - use http.Transport in WithHTTPClient instead.
// This function is a no-op and will be removed in a future version.
func WithTLSConfig(tlsConfig *http.Transport) ClientOption {
	return func(f *factory.ClientFactory) error {
		return nil // No-op - TLS config must be set via WithHTTPClient
	}
}
