// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package factory

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jontk/slurm-client/pkg/config"
	"github.com/jontk/slurm-client/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientFactory(t *testing.T) {
	tests := []struct {
		name        string
		options     []FactoryOption
		expectError bool
	}{
		{
			name: "create factory with default options",
			options: []FactoryOption{
				WithBaseURL("https://example.com"),
			},
			expectError: false,
		},
		{
			name: "create factory with custom config",
			options: []FactoryOption{
				WithConfig(config.NewDefault()),
				WithBaseURL("https://example.com"),
			},
			expectError: false,
		},
		{
			name:        "create factory with no options",
			options:     []FactoryOption{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory, err := NewClientFactory(tt.options...)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, factory)
			} else {
				helpers.AssertNoError(t, err)
				helpers.AssertNotNil(t, factory)
			}
		})
	}
}

func TestClientFactory_NewClientWithVersion(t *testing.T) {
	ctx := helpers.TestContext(t)

	factory, err := NewClientFactory(
		WithBaseURL("https://example.com"),
	)
	helpers.RequireNoError(t, err)

	tests := []struct {
		name        string
		version     string
		expectError bool
	}{
		{
			name:        "create client with v0.0.40",
			version:     "v0.0.40",
			expectError: false,
		},
		{
			name:        "create client with v0.0.41",
			version:     "v0.0.41",
			expectError: false,
		},
		{
			name:        "create client with v0.0.42",
			version:     "v0.0.42",
			expectError: false,
		},
		{
			name:        "create client with v0.0.43",
			version:     "v0.0.43",
			expectError: false,
		},
		{
			name:        "create client with v0.0.44",
			version:     "v0.0.44",
			expectError: false,
		},
		{
			name:        "create client with latest",
			version:     "latest",
			expectError: false, // Latest is v0.0.44, now implemented
		},
		{
			name:        "create client with stable",
			version:     "stable",
			expectError: false, // Stable is v0.0.42
		},
		{
			name:        "create client with unsupported version",
			version:     "v0.0.39",
			expectError: false, // Compatible version (v0.0.43) now implemented
		},
		{
			name:        "create client with invalid version",
			version:     "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := factory.NewClientWithVersion(ctx, tt.version)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				helpers.AssertNoError(t, err)
				helpers.AssertNotNil(t, client)

				// Verify client has expected interface methods
				assert.NotNil(t, client.Jobs())
				assert.NotNil(t, client.Nodes())
				assert.NotNil(t, client.Partitions())
				assert.NotNil(t, client.Info())
			}
		})
	}
}

func TestClientFactory_NewClient(t *testing.T) {
	ctx := helpers.TestContext(t)

	factory, err := NewClientFactory(
		WithBaseURL("https://example.com"),
	)
	helpers.RequireNoError(t, err)

	// Test with empty version (should auto-detect, fallback to stable)
	client, err := factory.NewClient(ctx)

	// Should succeed with stable version since we can't actually detect
	helpers.AssertNoError(t, err)
	helpers.AssertNotNil(t, client)
}

func TestClientFactory_NewClientForSlurmVersion(t *testing.T) {
	ctx := helpers.TestContext(t)

	factory, err := NewClientFactory(
		WithBaseURL("https://example.com"),
	)
	helpers.RequireNoError(t, err)

	tests := []struct {
		name         string
		slurmVersion string
		expectError  bool
	}{
		{
			name:         "Slurm 24.05",
			slurmVersion: "24.05",
			expectError:  false, // Should use v0.0.40
		},
		{
			name:         "Slurm 25.05",
			slurmVersion: "25.05",
			expectError:  false, // Would use v0.0.43 (latest), now implemented
		},
		{
			name:         "unsupported Slurm version",
			slurmVersion: "20.02",
			expectError:  true,
		},
		{
			name:         "invalid Slurm version",
			slurmVersion: "invalid",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := factory.NewClientForSlurmVersion(ctx, tt.slurmVersion)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				helpers.AssertNoError(t, err)
				helpers.AssertNotNil(t, client)
			}
		})
	}
}

func TestClientFactory_ListSupportedVersions(t *testing.T) {
	factory, err := NewClientFactory(
		WithBaseURL("https://example.com"),
	)
	helpers.RequireNoError(t, err)

	versions := factory.ListSupportedVersions()

	// Verify we have the expected versions
	expectedVersions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43", "v0.0.44"}

	helpers.AssertEqual(t, len(expectedVersions), len(versions))

	for i, version := range versions {
		helpers.AssertEqual(t, expectedVersions[i], version.String())
	}
}

func TestClientFactory_GetVersionCompatibility(t *testing.T) {
	factory, err := NewClientFactory(
		WithBaseURL("https://example.com"),
	)
	helpers.RequireNoError(t, err)

	matrix := factory.GetVersionCompatibility()

	helpers.AssertNotNil(t, matrix)
	helpers.AssertNotNil(t, matrix.SlurmVersions)
	helpers.AssertNotNil(t, matrix.BreakingChanges)

	// Verify all supported versions have Slurm version mappings
	for _, version := range []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43", "v0.0.44"} {
		slurmVersions, exists := matrix.SlurmVersions[version]
		assert.True(t, exists, "Version %s should have Slurm version mapping", version)
		assert.NotEmpty(t, slurmVersions, "Version %s should have at least one Slurm version", version)
	}
}

func TestClientInterfaces(t *testing.T) {
	ctx := helpers.TestContext(t)

	factory, err := NewClientFactory(
		WithBaseURL("https://example.com"),
	)
	helpers.RequireNoError(t, err)

	client, err := factory.NewClientWithVersion(ctx, "v0.0.42")
	helpers.AssertNoError(t, err)
	helpers.AssertNotNil(t, client)

	// Test that all manager interfaces are available
	jobManager := client.Jobs()
	nodeManager := client.Nodes()
	partitionManager := client.Partitions()
	infoManager := client.Info()

	helpers.AssertNotNil(t, jobManager)
	helpers.AssertNotNil(t, nodeManager)
	helpers.AssertNotNil(t, partitionManager)
	helpers.AssertNotNil(t, infoManager)

	// Test client version
	version := client.Version()
	helpers.AssertEqual(t, "v0.0.42", version)
}

// Phase 1: Unit tests for findCompatibleAPIVersion()
func TestClientFactory_findCompatibleAPIVersion(t *testing.T) {
	factory, err := NewClientFactory(
		WithBaseURL("https://example.com"),
	)
	helpers.RequireNoError(t, err)

	tests := []struct {
		name         string
		slurmVersion string
		expectedAPI  string
		expectError  bool
	}{
		{
			name:         "SLURM 24.05 maps to v0.0.40",
			slurmVersion: "24.05",
			expectedAPI:  "v0.0.40",
			expectError:  false,
		},
		{
			name:         "SLURM 24.11 maps to v0.0.41 (highest compatible)",
			slurmVersion: "24.11",
			expectedAPI:  "v0.0.41",
			expectError:  false,
		},
		{
			name:         "SLURM 25.05 maps to v0.0.43 (highest compatible)",
			slurmVersion: "25.05",
			expectedAPI:  "v0.0.43",
			expectError:  false,
		},
		{
			name:         "SLURM 25.11 maps to v0.0.44 (highest compatible)",
			slurmVersion: "25.11",
			expectedAPI:  "v0.0.44",
			expectError:  false,
		},
		{
			name:         "SLURM 25.11.1 with patch version",
			slurmVersion: "25.11.1",
			expectedAPI:  "v0.0.44",
			expectError:  false,
		},
		{
			name:         "SLURM 24.05.7 with patch version",
			slurmVersion: "24.05.7",
			expectedAPI:  "v0.0.40",
			expectError:  false,
		},
		{
			name:         "unsupported SLURM 20.11",
			slurmVersion: "20.11",
			expectedAPI:  "",
			expectError:  true,
		},
		{
			name:         "unsupported SLURM 99.99",
			slurmVersion: "99.99",
			expectedAPI:  "",
			expectError:  true,
		},
		{
			name:         "invalid version string",
			slurmVersion: "invalid",
			expectedAPI:  "",
			expectError:  true,
		},
		{
			name:         "empty version string",
			slurmVersion: "",
			expectedAPI:  "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := factory.findCompatibleAPIVersion(tt.slurmVersion)

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, version)
			} else {
				require.NoError(t, err, "unexpected error for SLURM version %s", tt.slurmVersion)
				require.NotNil(t, version)
				assert.Equal(t, tt.expectedAPI, version.String(), "SLURM %s should map to %s", tt.slurmVersion, tt.expectedAPI)
			}
		})
	}
}

// Phase 2: Unit tests for enhanced detectVersion() with mock HTTP server
func TestClientFactory_detectVersion_SlurmVersionString(t *testing.T) {
	ctx := helpers.TestContext(t)

	tests := []struct {
		name             string
		setupServer      func() *httptest.Server
		expectedVersion  string
		expectError      bool
		errorContains    string
	}{
		{
			name: "OpenAPI returns Slurm-25.11.1 version string",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/openapi/v3" {
						w.Header().Set("Content-Type", "application/json")
						resp := map[string]interface{}{
							"info": map[string]interface{}{
								"version": "Slurm-25.11.1",
							},
							"servers": []map[string]string{
								{"url": "/"},
							},
						}
						_ = json.NewEncoder(w).Encode(resp)
					} else {
						http.NotFound(w, r)
					}
				}))
			},
			expectedVersion: "v0.0.44",
			expectError:     false,
		},
		{
			name: "OpenAPI returns Slurm-24.05.0 version string",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/openapi/v3" {
						w.Header().Set("Content-Type", "application/json")
						resp := map[string]interface{}{
							"info": map[string]interface{}{
								"version": "Slurm-24.05.0",
							},
							"servers": []map[string]string{
								{"url": "/"},
							},
						}
						_ = json.NewEncoder(w).Encode(resp)
					} else {
						http.NotFound(w, r)
					}
				}))
			},
			expectedVersion: "v0.0.40",
			expectError:     false,
		},
		{
			name: "OpenAPI returns Slurm-25.05.3 version string",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/openapi/v3" {
						w.Header().Set("Content-Type", "application/json")
						resp := map[string]interface{}{
							"info": map[string]interface{}{
								"version": "Slurm-25.05.3",
							},
							"servers": []map[string]string{
								{"url": "/"},
							},
						}
						_ = json.NewEncoder(w).Encode(resp)
					} else {
						http.NotFound(w, r)
					}
				}))
			},
			expectedVersion: "v0.0.43",
			expectError:     false,
		},
		{
			name: "OpenAPI returns standard API version v0.0.42",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/openapi/v3" {
						w.Header().Set("Content-Type", "application/json")
						resp := map[string]interface{}{
							"info": map[string]interface{}{
								"version": "v0.0.42",
							},
							"servers": []map[string]string{
								{"url": "/"},
							},
						}
						_ = json.NewEncoder(w).Encode(resp)
					} else {
						http.NotFound(w, r)
					}
				}))
			},
			expectedVersion: "v0.0.42",
			expectError:     false,
		},
		{
			name: "OpenAPI returns version without v prefix",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/openapi/v3" {
						w.Header().Set("Content-Type", "application/json")
						resp := map[string]interface{}{
							"info": map[string]interface{}{
								"version": "0.0.42",
							},
							"servers": []map[string]string{
								{"url": "/"},
							},
						}
						_ = json.NewEncoder(w).Encode(resp)
					} else {
						http.NotFound(w, r)
					}
				}))
			},
			expectedVersion: "v0.0.42",
			expectError:     false,
		},
		{
			name: "OpenAPI extracts version from server URL",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/openapi/v3" {
						w.Header().Set("Content-Type", "application/json")
						resp := map[string]interface{}{
							"info": map[string]interface{}{
								"version": "",
							},
							"servers": []map[string]string{
								{"url": "/slurm/v0.0.42/"},
							},
						}
						_ = json.NewEncoder(w).Encode(resp)
					} else {
						http.NotFound(w, r)
					}
				}))
			},
			expectedVersion: "v0.0.42",
			expectError:     false,
		},
		{
			name: "unsupported SLURM version Slurm-20.11.0",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/openapi/v3" {
						w.Header().Set("Content-Type", "application/json")
						resp := map[string]interface{}{
							"info": map[string]interface{}{
								"version": "Slurm-20.11.0",
							},
							"servers": []map[string]string{
								{"url": "/"},
							},
						}
						_ = json.NewEncoder(w).Encode(resp)
					} else {
						http.NotFound(w, r)
					}
				}))
			},
			expectedVersion: "",
			expectError:     true,
			errorContains:   "no compatible API version",
		},
		{
			name: "invalid SLURM version string",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/openapi/v3" {
						w.Header().Set("Content-Type", "application/json")
						resp := map[string]interface{}{
							"info": map[string]interface{}{
								"version": "Slurm-invalid",
							},
							"servers": []map[string]string{
								{"url": "/"},
							},
						}
						_ = json.NewEncoder(w).Encode(resp)
					} else {
						http.NotFound(w, r)
					}
				}))
			},
			expectedVersion: "",
			expectError:     true,
			errorContains:   "invalid detected SLURM version",
		},
		{
			name: "HTTP 404 on /openapi/v3",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.NotFound(w, r)
				}))
			},
			expectedVersion: "",
			expectError:     true,
			errorContains:   "version detection failed with status 404",
		},
		{
			name: "malformed JSON response",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/openapi/v3" {
						w.Header().Set("Content-Type", "application/json")
						fmt.Fprintf(w, "{ invalid json")
					} else {
						http.NotFound(w, r)
					}
				}))
			},
			expectedVersion: "",
			expectError:     true,
			errorContains:   "failed to parse OpenAPI spec",
		},
		{
			name: "empty version and no server URLs",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/openapi/v3" {
						w.Header().Set("Content-Type", "application/json")
						resp := map[string]interface{}{
							"info": map[string]interface{}{
								"version": "",
							},
							"servers": []map[string]string{},
						}
						_ = json.NewEncoder(w).Encode(resp)
					} else {
						http.NotFound(w, r)
					}
				}))
			},
			expectedVersion: "",
			expectError:     true,
			errorContains:   "could not determine API version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer()
			defer server.Close()

			factory, err := NewClientFactory(
				WithBaseURL(server.URL),
			)
			require.NoError(t, err)

			version, err := factory.detectVersion(ctx)

			if tt.expectError {
				require.Error(t, err, "expected error for test case")
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains, "error message should contain expected text")
				}
				assert.Nil(t, version)
			} else {
				require.NoError(t, err, "unexpected error: %v", err)
				require.NotNil(t, version)
				assert.Equal(t, tt.expectedVersion, version.String(), "version mismatch")
			}
		})
	}
}

// Test that detectVersion caches the result
func TestClientFactory_detectVersion_Caching(t *testing.T) {
	ctx := helpers.TestContext(t)

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/openapi/v3" {
			callCount++
			w.Header().Set("Content-Type", "application/json")
			resp := map[string]interface{}{
				"info": map[string]interface{}{
					"version": "Slurm-25.11.1",
				},
				"servers": []map[string]string{
					{"url": "/"},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	factory, err := NewClientFactory(
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	// First call should hit the server
	version1, err := factory.detectVersion(ctx)
	require.NoError(t, err)
	require.NotNil(t, version1)
	initialCallCount := callCount

	// Second call should use cached version without hitting server
	version2, err := factory.detectVersion(ctx)
	require.NoError(t, err)
	require.NotNil(t, version2)

	// Verify call count didn't increase (cache was used)
	assert.Equal(t, initialCallCount, callCount, "second detectVersion should use cached result")

	// Verify both calls returned the same version
	assert.Equal(t, version1.String(), version2.String())
}

// Phase 3: End-to-end client creation with SLURM version detection
func TestClientFactory_NewClient_WithSlurmVersionDetection(t *testing.T) {
	ctx := helpers.TestContext(t)

	tests := []struct {
		name            string
		setupServer     func() *httptest.Server
		expectedVersion string
		expectError     bool
	}{
		{
			name: "client creation with Slurm-25.11.1 detection",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/openapi/v3" {
						w.Header().Set("Content-Type", "application/json")
						resp := map[string]interface{}{
							"info": map[string]interface{}{
								"version": "Slurm-25.11.1",
							},
							"servers": []map[string]string{
								{"url": "/"},
							},
						}
						_ = json.NewEncoder(w).Encode(resp)
					}
				}))
			},
			expectedVersion: "v0.0.44",
			expectError:     false,
		},
		{
			name: "client creation with Slurm-24.05.0 detection",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/openapi/v3" {
						w.Header().Set("Content-Type", "application/json")
						resp := map[string]interface{}{
							"info": map[string]interface{}{
								"version": "Slurm-24.05.0",
							},
							"servers": []map[string]string{
								{"url": "/"},
							},
						}
						_ = json.NewEncoder(w).Encode(resp)
					}
				}))
			},
			expectedVersion: "v0.0.40",
			expectError:     false,
		},
		{
			name: "client creation with v0.0.42 detection",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/openapi/v3" {
						w.Header().Set("Content-Type", "application/json")
						resp := map[string]interface{}{
							"info": map[string]interface{}{
								"version": "v0.0.42",
							},
							"servers": []map[string]string{
								{"url": "/"},
							},
						}
						_ = json.NewEncoder(w).Encode(resp)
					}
				}))
			},
			expectedVersion: "v0.0.42",
			expectError:     false,
		},
		{
			name: "client creation falls back to stable on detection failure",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Return 404 for all requests
					http.NotFound(w, r)
				}))
			},
			expectedVersion: "v0.0.42", // Stable version
			expectError:     false,      // Should not error, just use fallback
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer()
			defer server.Close()

			factory, err := NewClientFactory(
				WithBaseURL(server.URL),
			)
			require.NoError(t, err)

			// Call NewClient with empty version (triggers auto-detection)
			client, err := factory.NewClient(ctx)

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, client)
			} else {
				require.NoError(t, err, "expected successful client creation")
				require.NotNil(t, client)

				// Verify client has the expected version
				assert.Equal(t, tt.expectedVersion, client.Version())

				// Verify client has expected interface methods
				assert.NotNil(t, client.Jobs())
				assert.NotNil(t, client.Nodes())
				assert.NotNil(t, client.Partitions())
				assert.NotNil(t, client.Info())
			}
		})
	}
}

// Test SLURM version detection with authentication
func TestClientFactory_detectVersion_WithAuthentication(t *testing.T) {
	ctx := helpers.TestContext(t)

	// Test that version detection passes authentication if provided
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/openapi/v3" {
			w.Header().Set("Content-Type", "application/json")
			resp := map[string]interface{}{
				"info": map[string]interface{}{
					"version": "Slurm-25.11.1",
				},
				"servers": []map[string]string{
					{"url": "/"},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	factory, err := NewClientFactory(
		WithBaseURL(server.URL),
		WithConfig(config.NewDefault()),
	)
	require.NoError(t, err)

	// Should succeed without auth provider (auth is optional)
	version, err := factory.detectVersion(ctx)
	require.NoError(t, err)
	require.NotNil(t, version)
	assert.Equal(t, "v0.0.44", version.String())
}
