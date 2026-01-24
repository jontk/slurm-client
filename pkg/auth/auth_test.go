// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"net/http"
	"testing"

	"github.com/jontk/slurm-client/tests/helpers"
)

func TestTokenAuth(t *testing.T) {
	token := "test-token-123"
	auth := NewTokenAuth(token)

	// Test Type method
	helpers.AssertEqual(t, "token", auth.Type())

	// Test Authenticate method
	ctx := helpers.TestContext(t)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", http.NoBody)
	helpers.RequireNoError(t, err)

	err = auth.Authenticate(ctx, req)
	helpers.AssertNoError(t, err)

	// Verify token was added to header
	helpers.AssertEqual(t, token, req.Header.Get("X-SLURM-USER-TOKEN"))
}

func TestBasicAuth(t *testing.T) {
	username := "testuser"
	password := "testpass"
	auth := NewBasicAuth(username, password)

	// Test Type method
	helpers.AssertEqual(t, "basic", auth.Type())

	// Test Authenticate method
	ctx := helpers.TestContext(t)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", http.NoBody)
	helpers.RequireNoError(t, err)

	err = auth.Authenticate(ctx, req)
	helpers.AssertNoError(t, err)

	// Verify basic auth was added to header
	username_from_req, password_from_req, ok := req.BasicAuth()
	helpers.AssertEqual(t, true, ok)
	helpers.AssertEqual(t, username, username_from_req)
	helpers.AssertEqual(t, password, password_from_req)
}

func TestNoAuth(t *testing.T) {
	auth := NewNoAuth()

	// Test Type method
	helpers.AssertEqual(t, "none", auth.Type())

	// Test Authenticate method
	ctx := helpers.TestContext(t)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", http.NoBody)
	helpers.RequireNoError(t, err)

	// Store original headers
	originalHeaders := make(http.Header)
	for key, values := range req.Header {
		originalHeaders[key] = values
	}

	err = auth.Authenticate(ctx, req)
	helpers.AssertNoError(t, err)

	// Verify no headers were added
	for key, values := range req.Header {
		helpers.AssertEqual(t, originalHeaders[key], values)
	}

	// Verify no auth headers were added
	helpers.AssertEqual(t, "", req.Header.Get("X-SLURM-USER-TOKEN"))
	helpers.AssertEqual(t, "", req.Header.Get("Authorization"))
}

func TestAuthProviderInterface(t *testing.T) {
	// Test that all auth types implement the Provider interface
	var _ Provider = &TokenAuth{}
	var _ Provider = &BasicAuth{}
	var _ Provider = &NoAuth{}

	// Test different auth providers
	providers := []Provider{
		NewTokenAuth("test-token"),
		NewBasicAuth("user", "pass"),
		NewNoAuth(),
	}

	for _, provider := range providers {
		// Each provider should have a type
		authType := provider.Type()
		helpers.AssertNotNil(t, authType)

		// Each provider should be able to authenticate
		ctx := helpers.TestContext(t)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", http.NoBody)
		helpers.RequireNoError(t, err)

		err = provider.Authenticate(ctx, req)
		helpers.AssertNoError(t, err)
	}
}

func TestTokenAuthWithEmptyToken(t *testing.T) {
	auth := NewTokenAuth("")

	ctx := helpers.TestContext(t)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", http.NoBody)
	helpers.RequireNoError(t, err)

	err = auth.Authenticate(ctx, req)
	helpers.AssertNoError(t, err)

	// Verify empty token is still set (it's up to the server to validate)
	helpers.AssertEqual(t, "", req.Header.Get("X-SLURM-USER-TOKEN"))
}

func TestBasicAuthWithEmptyCredentials(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
	}{
		{
			name:     "empty username",
			username: "",
			password: "password",
		},
		{
			name:     "empty password",
			username: "username",
			password: "",
		},
		{
			name:     "both empty",
			username: "",
			password: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := NewBasicAuth(tt.username, tt.password)

			ctx := helpers.TestContext(t)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", http.NoBody)
			helpers.RequireNoError(t, err)

			err = auth.Authenticate(ctx, req)
			helpers.AssertNoError(t, err)

			// Verify basic auth was set (even if empty)
			username_from_req, password_from_req, ok := req.BasicAuth()
			helpers.AssertEqual(t, true, ok)
			helpers.AssertEqual(t, tt.username, username_from_req)
			helpers.AssertEqual(t, tt.password, password_from_req)
		})
	}
}

func TestAuthenticateMultipleTimes(t *testing.T) {
	// Test that authentication can be called multiple times
	auth := NewTokenAuth("test-token")

	ctx := helpers.TestContext(t)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", http.NoBody)
	helpers.RequireNoError(t, err)

	// First authentication
	err = auth.Authenticate(ctx, req)
	helpers.AssertNoError(t, err)
	helpers.AssertEqual(t, "test-token", req.Header.Get("X-SLURM-USER-TOKEN"))

	// Second authentication (should overwrite)
	err = auth.Authenticate(ctx, req)
	helpers.AssertNoError(t, err)
	helpers.AssertEqual(t, "test-token", req.Header.Get("X-SLURM-USER-TOKEN"))

	// Verify token header exists
	tokenValue := req.Header.Get("X-SLURM-USER-TOKEN")
	helpers.AssertEqual(t, "test-token", tokenValue)
}
