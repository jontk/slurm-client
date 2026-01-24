// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/tests/helpers"
	"github.com/jontk/slurm-client/tests/mocks"
)

// TestAuthenticationFlows tests different authentication providers
func TestAuthenticationFlows(t *testing.T) {
	testCases := []struct {
		name         string
		authProvider auth.Provider
		enableAuth   bool
		expectError  bool
		description  string
	}{
		{
			name:         "NoAuth",
			authProvider: auth.NewNoAuth(),
			enableAuth:   false,
			expectError:  false,
			description:  "No authentication required",
		},
		{
			name:         "ValidTokenAuth",
			authProvider: auth.NewTokenAuth("test-token-v42"),
			enableAuth:   true,
			expectError:  false,
			description:  "Valid token authentication",
		},
		{
			name:         "InvalidTokenAuth",
			authProvider: auth.NewTokenAuth("invalid-token"),
			enableAuth:   true,
			expectError:  true,
			description:  "Invalid token should fail",
		},
		{
			name:         "EmptyTokenAuth",
			authProvider: auth.NewTokenAuth(""),
			enableAuth:   true,
			expectError:  true,
			description:  "Empty token should fail",
		},
		{
			name:         "BasicAuth",
			authProvider: auth.NewBasicAuth("admin", "password"),
			enableAuth:   false, // Mock server doesn't implement basic auth validation
			expectError:  false,
			description:  "Basic authentication",
		},
		{
			name:         "EmptyBasicAuth",
			authProvider: auth.NewBasicAuth("", ""),
			enableAuth:   false,
			expectError:  false,
			description:  "Empty basic auth credentials",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testAuthenticationFlow(t, tc.authProvider, tc.enableAuth, tc.expectError, tc.description)
		})
	}
}

func testAuthenticationFlow(t *testing.T, authProvider auth.Provider, enableAuth, expectError bool, description string) {
	t.Log(description)

	// Create mock server with authentication enabled/disabled
	config := mocks.DefaultServerConfig()
	config.EnableAuth = enableAuth
	mockServer := mocks.NewMockSlurmServer(config)
	defer mockServer.Close()

	ctx := helpers.TestContext(t)

	// Create client with the specified auth provider
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.42",
		slurm.WithBaseURL(mockServer.URL()),
		slurm.WithAuth(authProvider),
	)

	if expectError {
		// For some auth errors, client creation might succeed but operations should fail
		if err == nil {
			require.NotNil(t, client)
			defer client.Close()

			// Try a simple operation that should fail
			err = client.Info().Ping(ctx)
			assert.Error(t, err, "Operation should fail with invalid auth")
			assert.Contains(t, err.Error(), "Unauthorized")
		} else {
			// Client creation failed, which is also acceptable for auth errors
			assert.Error(t, err)
		}
		return
	}

	// For successful cases
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	// Test basic operations with authentication
	err = client.Info().Ping(ctx)
	assert.NoError(t, err, "Ping should succeed with valid auth")

	// Test job operations
	jobs, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{Limit: 5})
	assert.NoError(t, err, "Job list should succeed with valid auth")
	assert.NotNil(t, jobs)

	// Test job submission
	submission := &interfaces.JobSubmission{
		Name:      "auth-test-job",
		Script:    "#!/bin/bash\necho 'Auth test'",
		Partition: "compute",
		CPUs:      1,
	}

	response, err := client.Jobs().Submit(ctx, submission)
	assert.NoError(t, err, "Job submission should succeed with valid auth")
	assert.NotNil(t, response)
}

// TestAuthenticationAcrossVersions tests authentication across different API versions
func TestAuthenticationAcrossVersions(t *testing.T) {
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}

	for _, version := range versions {
		t.Run("Version_"+version, func(t *testing.T) {
			testAuthForVersion(t, version)
		})
	}
}

func testAuthForVersion(t *testing.T, version string) {
	// Create mock server with authentication enabled
	config, exists := mocks.VersionConfigs[version]
	require.True(t, exists, "Version config should exist for %s", version)

	configCopy := *config
	configCopy.EnableAuth = true
	configCopy.ErrorResponses = make(map[string]mocks.ErrorResponse)

	mockServer := mocks.NewMockSlurmServer(&configCopy)
	defer mockServer.Close()

	ctx := helpers.TestContext(t)

	// Test with correct token for this version
	client, err := slurm.NewClientWithVersion(ctx, version,
		slurm.WithBaseURL(mockServer.URL()),
		slurm.WithAuth(auth.NewTokenAuth(configCopy.AuthToken)),
	)
	require.NoError(t, err)
	defer client.Close()

	// Verify authentication works
	err = client.Info().Ping(ctx)
	assert.NoError(t, err, "Ping should work with correct auth for %s", version)

	// Test job operations
	jobs, err := client.Jobs().List(ctx, nil)
	assert.NoError(t, err, "Job list should work with auth for %s", version)
	assert.NotNil(t, jobs)

	t.Logf("Authentication successful for %s", version)
}

// TestAuthenticationErrors tests various authentication error scenarios
func TestAuthenticationErrors(t *testing.T) {
	config := mocks.DefaultServerConfig()
	config.EnableAuth = true
	mockServer := mocks.NewMockSlurmServer(config)
	defer mockServer.Close()

	ctx := helpers.TestContext(t)

	errorScenarios := []struct {
		name         string
		authProvider auth.Provider
		expectStatus int
		expectMsg    string
	}{
		{
			name:         "WrongToken",
			authProvider: auth.NewTokenAuth("wrong-token"),
			expectStatus: 401,
			expectMsg:    "Unauthorized",
		},
		{
			name:         "MalformedToken",
			authProvider: auth.NewTokenAuth("malformed token with spaces"),
			expectStatus: 401,
			expectMsg:    "Unauthorized",
		},
		{
			name:         "EmptyToken",
			authProvider: auth.NewTokenAuth(""),
			expectStatus: 401,
			expectMsg:    "Unauthorized",
		},
	}

	for _, scenario := range errorScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			client, err := slurm.NewClientWithVersion(ctx, "v0.0.42",
				slurm.WithBaseURL(mockServer.URL()),
				slurm.WithAuth(scenario.authProvider),
			)

			// Client creation might succeed, but operations should fail
			if err == nil {
				require.NotNil(t, client)
				defer client.Close()

				err = client.Info().Ping(ctx)
			}

			require.Error(t, err)
			assert.Contains(t, err.Error(), scenario.expectMsg)
			t.Logf("Got expected error for %s: %v", scenario.name, err)
		})
	}
}

// TestAuthenticationHeaders tests that authentication headers are properly set
func TestAuthenticationHeaders(t *testing.T) {
	testCases := []struct {
		name         string
		authProvider auth.Provider
		expectedAuth string
		description  string
	}{
		{
			name:         "TokenAuth",
			authProvider: auth.NewTokenAuth("test-token-123"),
			expectedAuth: "test-token-123",
			description:  "Token auth should set X-SLURM-USER-TOKEN header",
		},
		{
			name:         "BasicAuth",
			authProvider: auth.NewBasicAuth("user", "pass"),
			expectedAuth: "Basic dXNlcjpwYXNz", // base64 of "user:pass"
			description:  "Basic auth should set Basic header",
		},
		{
			name:         "NoAuth",
			authProvider: auth.NewNoAuth(),
			expectedAuth: "",
			description:  "No auth should not set auth header",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.description)

			// Create a mock HTTP request to test header setting
			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", http.NoBody)
			require.NoError(t, err)

			// Apply authentication
			err = tc.authProvider.Authenticate(ctx, req)
			require.NoError(t, err)

			// Check the appropriate header based on auth type
			var authHeader string
			if tc.name == "TokenAuth" {
				authHeader = req.Header.Get("X-SLURM-USER-TOKEN")
			} else if tc.name == "BasicAuth" {
				authHeader = req.Header.Get("Authorization")
			}

			if tc.expectedAuth == "" {
				assert.Empty(t, authHeader, "No auth header should be set for %s", tc.name)
			} else {
				assert.Equal(t, tc.expectedAuth, authHeader, "Auth header should match expected for %s", tc.name)
			}
		})
	}
}

// TestConcurrentAuthenticatedRequests tests multiple concurrent requests with authentication
func TestConcurrentAuthenticatedRequests(t *testing.T) {
	config := mocks.DefaultServerConfig()
	config.EnableAuth = true
	config.ResponseDelay = 100 * time.Millisecond // Add delay to test concurrency
	mockServer := mocks.NewMockSlurmServer(config)
	defer mockServer.Close()

	ctx := helpers.TestContext(t)

	client, err := slurm.NewClientWithVersion(ctx, "v0.0.42",
		slurm.WithBaseURL(mockServer.URL()),
		slurm.WithAuth(auth.NewTokenAuth(config.AuthToken)),
	)
	require.NoError(t, err)
	defer client.Close()

	// Run multiple operations concurrently
	const numConcurrent = 5
	errors := make(chan error, numConcurrent)

	for i := range numConcurrent {
		go func(index int) {
			// Each goroutine performs multiple operations
			err := client.Info().Ping(ctx)
			if err != nil {
				errors <- err
				return
			}

			_, err = client.Jobs().List(ctx, &interfaces.ListJobsOptions{Limit: 1})
			if err != nil {
				errors <- err
				return
			}

			_, err = client.Nodes().List(ctx, &interfaces.ListNodesOptions{Limit: 1})
			errors <- err
		}(i)
	}

	// Collect results
	for i := range numConcurrent {
		err := <-errors
		assert.NoError(t, err, "Concurrent request %d should succeed", i)
	}

	t.Log("All concurrent authenticated requests completed successfully")
}
