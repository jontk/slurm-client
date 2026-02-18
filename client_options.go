// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Package slurm provides client options for configuring the SLURM REST API client
package slurm

import (
	"context"
	"net/http"
	"time"

	"github.com/jontk/slurm-client/internal/factory"
	"github.com/jontk/slurm-client/pkg/auth"
)

// Additional client options that aren't in client.go

// WithToken sets the authentication token (DEPRECATED)
//
// IMPORTANT: Most SLURM deployments require both X-SLURM-USER-NAME and
// X-SLURM-USER-TOKEN headers. This function only sets X-SLURM-USER-TOKEN,
// which will cause authentication failures with slurmrestd.
//
// Deprecated: Use WithUserToken(username, token) for full authentication.
// Example:
//
//	client, err := slurm.NewClient(ctx,
//	    slurm.WithBaseURL("https://cluster:6820"),
//	    slurm.WithUserToken("username", "your-jwt-token"),
//	)
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

// WithVersion is deprecated and has no effect.
//
// To specify a version, use NewClientWithVersion instead:
//
//	client, err := slurm.NewClientWithVersion(ctx, "v0.0.44", opts...)
//
// Deprecated: This function is a no-op. Use NewClientWithVersion to specify the API version.
func WithVersion(version string) ClientOption {
	return func(f *factory.ClientFactory) error {
		// This is a no-op - version must be specified via NewClientWithVersion()
		// Keeping this for backward compatibility but it will be removed in a future release
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

// WithTimeout sets default timeout for all operations
// This modifies the existing HTTP client's timeout without replacing the client,
// preserving TLS configuration and custom transport settings
func WithTimeout(timeout time.Duration) ClientOption {
	return func(f *factory.ClientFactory) error {
		return f.SetTimeout(timeout)
	}
}
