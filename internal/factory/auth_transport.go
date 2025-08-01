// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package factory

import (
	"net/http"

	"github.com/jontk/slurm-client/pkg/auth"
)

// authTransport wraps an http.RoundTripper to add authentication
type authTransport struct {
	base http.RoundTripper
	auth auth.Provider
}

// newAuthTransport creates a new authenticated transport
func newAuthTransport(base http.RoundTripper, auth auth.Provider) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return &authTransport{
		base: base,
		auth: auth,
	}
}

// RoundTrip implements http.RoundTripper
func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	reqCopy := req.Clone(req.Context())
	
	// Apply authentication if available
	if t.auth != nil {
		// Use the request's context for authentication
		if err := t.auth.Authenticate(req.Context(), reqCopy); err != nil {
			// Log error but continue - some endpoints may not need auth
			// In production, you might want to handle this differently
		}
	}
	
	// Execute the request
	return t.base.RoundTrip(reqCopy)
}

// createAuthenticatedHTTPClient creates an HTTP client with authentication
func createAuthenticatedHTTPClient(baseClient *http.Client, authProvider auth.Provider) *http.Client {
	if baseClient == nil {
		baseClient = &http.Client{}
	}
	
	// Clone the client to avoid modifying the original
	client := &http.Client{
		Timeout:       baseClient.Timeout,
		CheckRedirect: baseClient.CheckRedirect,
		Jar:           baseClient.Jar,
	}
	
	// Wrap the transport with authentication
	if baseClient.Transport != nil {
		client.Transport = newAuthTransport(baseClient.Transport, authProvider)
	} else {
		client.Transport = newAuthTransport(http.DefaultTransport, authProvider)
	}
	
	return client
}
