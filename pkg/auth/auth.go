package auth

import (
	"context"
	"net/http"
)

// Provider defines the interface for authentication providers
type Provider interface {
	// Authenticate adds authentication to the HTTP request
	Authenticate(ctx context.Context, req *http.Request) error
	
	// Type returns the authentication type
	Type() string
}

// TokenAuth implements token-based authentication
type TokenAuth struct {
	token string
}

// NewTokenAuth creates a new token-based authentication provider
func NewTokenAuth(token string) *TokenAuth {
	return &TokenAuth{token: token}
}

// Authenticate adds the token to the request
func (t *TokenAuth) Authenticate(ctx context.Context, req *http.Request) error {
	req.Header.Set("X-SLURM-USER-TOKEN", t.token)
	return nil
}

// Type returns the authentication type
func (t *TokenAuth) Type() string {
	return "token"
}

// BasicAuth implements basic authentication
type BasicAuth struct {
	username string
	password string
}

// NewBasicAuth creates a new basic authentication provider
func NewBasicAuth(username, password string) *BasicAuth {
	return &BasicAuth{
		username: username,
		password: password,
	}
}

// Authenticate adds basic auth to the request
func (b *BasicAuth) Authenticate(ctx context.Context, req *http.Request) error {
	req.SetBasicAuth(b.username, b.password)
	return nil
}

// Type returns the authentication type
func (b *BasicAuth) Type() string {
	return "basic"
}

// NoAuth implements no authentication
type NoAuth struct{}

// NewNoAuth creates a new no-auth provider
func NewNoAuth() *NoAuth {
	return &NoAuth{}
}

// Authenticate is a no-op for no authentication
func (n *NoAuth) Authenticate(ctx context.Context, req *http.Request) error {
	return nil
}

// Type returns the authentication type
func (n *NoAuth) Type() string {
	return "none"
}