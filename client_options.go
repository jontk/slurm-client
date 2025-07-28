// Package slurm provides client options for configuring the SLURM REST API client
package slurm

import (
	"context"
	"net/http"
	"time"

	"github.com/jontk/slurm-client/internal/factory"
	"github.com/jontk/slurm-client/pkg/auth"
	slurmctx "github.com/jontk/slurm-client/pkg/context"
	"github.com/jontk/slurm-client/pkg/logging"
	"github.com/jontk/slurm-client/pkg/metrics"
	"github.com/jontk/slurm-client/pkg/middleware"
	"github.com/jontk/slurm-client/pkg/pool"
	"github.com/jontk/slurm-client/pkg/retry"
)

// Additional client options that aren't in client.go

// WithToken sets the authentication token
func WithToken(token string) ClientOption {
	return func(f *factory.ClientFactory) error {
		// Create a token provider for the given token
		return factory.WithAuth(auth.NewTokenAuth(token))(f)
	}
}

// WithUserToken sets user authentication with username and token
func WithUserToken(username, token string) ClientOption {
	return func(f *factory.ClientFactory) error {
		// Create a user token provider using the TokenAuth with X-SLURM-USER-NAME header
		provider := &userTokenAuth{
			username: username,
			token:    token,
		}
		return factory.WithAuth(provider)(f)
	}
}

// userTokenAuth implements user token authentication
type userTokenAuth struct {
	username string
	token    string
}

func (u *userTokenAuth) Authenticate(ctx context.Context, req *http.Request) error {
	req.Header.Set("X-SLURM-USER-NAME", u.username)
	req.Header.Set("X-SLURM-USER-TOKEN", u.token)
	return nil
}

func (u *userTokenAuth) Type() string {
	return "user-token"
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(f *factory.ClientFactory) error {
		return factory.WithHTTPClient(client)(f)
	}
}

// WithVersion forces a specific API version
func WithVersion(version string) ClientOption {
	return func(f *factory.ClientFactory) error {
		// Version handling is done at client creation time
		// This is a no-op for now
		return nil
	}
}

// WithNoAuth disables authentication (for testing or public endpoints)
func WithNoAuth() ClientOption {
	return func(f *factory.ClientFactory) error {
		// Set a no-op auth provider
		return factory.WithAuth(&noAuth{})(f)
	}
}

// noAuth implements a no-op authentication provider
type noAuth struct{}

func (n *noAuth) Authenticate(ctx context.Context, req *http.Request) error {
	return nil
}

func (n *noAuth) Type() string {
	return "none"
}

// === Enhanced feature options ===
// NOTE: These options may need to be implemented in the factory package
// For now, they are no-ops to allow the code to compile

// WithLogger sets a custom logger for the client
func WithLogger(logger logging.Logger) ClientOption {
	return func(f *factory.ClientFactory) error {
		// TODO: Implement in factory
		return nil
	}
}

// WithLogLevel sets the logging level
func WithLogLevel(level string) ClientOption {
	return func(f *factory.ClientFactory) error {
		// TODO: Implement in factory
		return nil
	}
}

// WithMetricsCollector sets a custom metrics collector
func WithMetricsCollector(collector metrics.Collector) ClientOption {
	return func(f *factory.ClientFactory) error {
		// TODO: Implement in factory
		return nil
	}
}

// WithMiddleware adds custom middleware to the HTTP client
func WithMiddleware(middlewares ...middleware.Middleware) ClientOption {
	return func(f *factory.ClientFactory) error {
		// TODO: Implement in factory
		return nil
	}
}

// WithTimeout sets default timeout for all operations
func WithTimeout(timeout time.Duration) ClientOption {
	return func(f *factory.ClientFactory) error {
		// Create a new HTTP client with the specified timeout
		httpClient := &http.Client{
			Timeout: timeout,
		}
		return factory.WithHTTPClient(httpClient)(f)
	}
}

// WithTimeoutConfig sets custom timeout configuration for different operation types
func WithTimeoutConfig(config *slurmctx.TimeoutConfig) ClientOption {
	return func(f *factory.ClientFactory) error {
		// TODO: Implement in factory
		return nil
	}
}

// WithConnectionPool enables connection pooling with the specified configuration
func WithConnectionPool(config *pool.PoolConfig) ClientOption {
	return func(f *factory.ClientFactory) error {
		// TODO: Implement in factory
		return nil
	}
}

// WithRetryBackoff sets a custom retry backoff strategy
func WithRetryBackoff(backoff retry.BackoffStrategy) ClientOption {
	return func(f *factory.ClientFactory) error {
		// TODO: Implement in factory
		return nil
	}
}

// WithMaxRetries sets the maximum number of retry attempts
func WithMaxRetries(maxRetries int) ClientOption {
	return func(f *factory.ClientFactory) error {
		// TODO: Implement in factory
		return nil
	}
}

// WithUserAgent sets a custom User-Agent header
func WithUserAgent(userAgent string) ClientOption {
	return func(f *factory.ClientFactory) error {
		// TODO: Implement in factory
		return nil
	}
}

// WithRequestID enables request ID generation with the provided generator
func WithRequestID(generator func() string) ClientOption {
	return func(f *factory.ClientFactory) error {
		// TODO: Implement in factory
		return nil
	}
}

// WithCircuitBreaker enables circuit breaker functionality
func WithCircuitBreaker(threshold int, timeout time.Duration) ClientOption {
	return func(f *factory.ClientFactory) error {
		// TODO: Implement in factory
		return nil
	}
}

// WithCompression enables or disables HTTP compression
func WithCompression(enabled bool) ClientOption {
	return func(f *factory.ClientFactory) error {
		// TODO: Implement in factory
		return nil
	}
}

// WithKeepAlive enables or disables HTTP keep-alive
func WithKeepAlive(enabled bool) ClientOption {
	return func(f *factory.ClientFactory) error {
		// TODO: Implement in factory
		return nil
	}
}

// WithDebug enables debug mode with verbose logging
func WithDebug() ClientOption {
	return func(f *factory.ClientFactory) error {
		// TODO: Implement in factory
		return nil
	}
}

// WithTLSConfig sets custom TLS configuration
func WithTLSConfig(tlsConfig *http.Transport) ClientOption {
	return func(f *factory.ClientFactory) error {
		// TODO: Implement in factory
		return nil
	}
}