// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// TestPingWithRealServer tests basic connectivity to a real SLURM server
func TestPingWithRealServer(t *testing.T) {
	// Skip if not explicitly enabled
	if os.Getenv("SLURM_REAL_SERVER_TEST") != "true" {
		t.Skip("Real server tests disabled. Set SLURM_REAL_SERVER_TEST=true to enable")
	}

	// Get server configuration
	serverURL := os.Getenv("SLURM_SERVER_URL")
	if serverURL == "" {
		serverURL = "http://rocky9:6820"
	}

	token := os.Getenv("SLURM_JWT_TOKEN")
	if token == "" {
		// Use the token you provided
		token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTM3MzE3MjgsImlhdCI6MTc1MzcyOTkyOCwic3VuIjoicm9vdCJ9.7DGZd7hWhJQhkIx_0wMsKGM2rDipM27CGgaZFU1z_Ns"
	}

	// Create a transport that adds the token
	transport := &tokenTransport{
		token: token,
		base:  http.DefaultTransport,
	}
	
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	// Create API client
	apiClient, err := api.NewClientWithResponses(
		serverURL,
		api.WithHTTPClient(httpClient),
	)
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("Ping SLURM Server", func(t *testing.T) {
		// Use the correct ping endpoint
		resp, err := apiClient.SlurmV0043GetPingWithResponse(ctx)
		require.NoError(t, err)
		
		assert.Equal(t, 200, resp.StatusCode())
		
		if resp.JSON200 != nil {
			t.Logf("Ping response received")
			
			// Log the response
			jsonData, _ := json.MarshalIndent(resp.JSON200, "", "  ")
			t.Logf("Response: %s", string(jsonData))
			
			// Check pings
			if resp.JSON200.Pings != nil {
				for i, ping := range resp.JSON200.Pings {
					t.Logf("  Ping %d: hostname=%s, status=%s, responding=%v, latency=%d",
						i+1, 
						getStringValue(ping.Hostname), 
						getStringValue(ping.Pinged),
						ping.Responding,
						getInt64Value(ping.Latency))
				}
			}
			
			// Check meta information
			if resp.JSON200.Meta != nil && resp.JSON200.Meta.Slurm != nil {
				if resp.JSON200.Meta.Slurm.Version != nil {
					t.Logf("SLURM Version: %s.%s.%s",
						getStringValue(resp.JSON200.Meta.Slurm.Version.Major),
						getStringValue(resp.JSON200.Meta.Slurm.Version.Minor),
						getStringValue(resp.JSON200.Meta.Slurm.Version.Micro))
				}
				t.Logf("SLURM Release: %s", getStringValue(resp.JSON200.Meta.Slurm.Release))
				t.Logf("SLURM Cluster: %s", getStringValue(resp.JSON200.Meta.Slurm.Cluster))
			}
		}
	})

	t.Run("Check OpenAPI Spec", func(t *testing.T) {
		// Make a raw HTTP request for the OpenAPI spec
		req, err := http.NewRequestWithContext(ctx, "GET", serverURL+"/openapi/v3", nil)
		require.NoError(t, err)
		
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		
		assert.Equal(t, 200, resp.StatusCode)
		
		var openapi map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&openapi)
		require.NoError(t, err)
		
		// Check OpenAPI version
		if openAPIVersion, ok := openapi["openapi"].(string); ok {
			t.Logf("OpenAPI version: %s", openAPIVersion)
		}
		
		// Check available paths
		if paths, ok := openapi["paths"].(map[string]interface{}); ok {
			t.Logf("Found %d API endpoints", len(paths))
			
			// Count QoS endpoints
			qosCount := 0
			for path := range paths {
				if strings.Contains(path, "qos") {
					qosCount++
				}
			}
			t.Logf("Found %d QoS-related endpoints", qosCount)
		}
	})
}

// Helper functions to safely get values from pointers
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func getIntValue(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func getInt64Value(ptr *int64) int64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

