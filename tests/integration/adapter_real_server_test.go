// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jontk/slurm-client/internal/adapters/v0_0_43"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
	"github.com/jontk/slurm-client/internal/common/types"
)

// TestAdapterWithRealServer tests the adapter pattern against a real SLURM server
func TestAdapterWithRealServer(t *testing.T) {
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

	// Create adapter
	adapter := v0_0_43.NewQoSAdapter(apiClient)

	ctx := context.Background()

	t.Run("List QoS on Real Server", func(t *testing.T) {
		qosList, err := adapter.List(ctx, &types.QoSListOptions{
			Limit: 10,
		})

		if err != nil {
			t.Logf("Error listing QoS: %v", err)
			// Log more details about the error
			t.Logf("Server URL: %s", serverURL)
			t.Logf("Token (first 20 chars): %s...", token[:20])

			// Check if it's a database connection error
			if strings.Contains(err.Error(), "Unable to connect to database") ||
				strings.Contains(err.Error(), "Failed to open slurmdbd connection") ||
				strings.Contains(err.Error(), "SLURM_DAEMON_DOWN") {
				t.Logf("slurmdbd connection failed - likely due to auth plugin mismatch (JWT vs munge)")
				t.Logf("Check slurmdbd logs for: 'authentication plugin auth/jwt not found'")
				t.Skip("Skipping test: slurmdbd is not connected. This is expected when auth plugins don't match.")
				return
			}
		}

		require.NoError(t, err)
		assert.NotNil(t, qosList)

		t.Logf("Found %d QoS entries on the server", qosList.Total)
		for i, qos := range qosList.QoS {
			if i < 5 { // Log first 5 entries
				t.Logf("  QoS %d: Name=%s, Priority=%d, UsageFactor=%.2f, GraceTime=%d",
					i+1, qos.Name, qos.Priority, qos.UsageFactor, qos.GraceTime)
				if qos.Limits != nil {
					if qos.Limits.MaxJobsPerUser != nil {
						t.Logf("    Max Jobs Per User: %d", *qos.Limits.MaxJobsPerUser)
					}
					if qos.Limits.MaxJobsPerAccount != nil {
						t.Logf("    Max Jobs Per Account: %d", *qos.Limits.MaxJobsPerAccount)
					}
				}
			}
		}
	})

	t.Run("Get Specific QoS", func(t *testing.T) {
		// First list to get a valid QoS name
		qosList, err := adapter.List(ctx, &types.QoSListOptions{
			Limit: 1,
		})
		require.NoError(t, err)

		if len(qosList.QoS) == 0 {
			t.Skip("No QoS entries found on server")
			return
		}

		qosName := qosList.QoS[0].Name
		t.Logf("Testing Get with QoS: %s", qosName)

		qos, err := adapter.Get(ctx, qosName)
		require.NoError(t, err)
		assert.NotNil(t, qos)
		assert.Equal(t, qosName, qos.Name)

		t.Logf("Retrieved QoS Details:")
		t.Logf("  Name: %s", qos.Name)
		t.Logf("  Description: %s", qos.Description)
		t.Logf("  Priority: %d", qos.Priority)
		t.Logf("  UsageFactor: %.2f", qos.UsageFactor)
		t.Logf("  UsageThreshold: %.2f", qos.UsageThreshold)
		t.Logf("  GraceTime: %d", qos.GraceTime)
		t.Logf("  Flags: %v", qos.Flags)
	})

	// Only run create/delete tests if explicitly enabled
	if os.Getenv("SLURM_ENABLE_WRITE_TESTS") == "true" {
		t.Run("Create and Delete QoS", func(t *testing.T) {
			testQoSName := "go-test-qos-" + time.Now().Format("20060102-150405")

			// Create QoS
			createReq := &types.QoSCreate{
				Name:        testQoSName,
				Description: "Test QoS created by Go adapter pattern test",
				Priority:    500,
				Flags:       []string{"DenyOnLimit"},
				UsageFactor: 2.0,
				GraceTime:   600,
				Limits: &types.QoSLimits{
					MaxJobsPerUser:    intPtr(20),
					MaxJobsPerAccount: intPtr(100),
				},
			}

			t.Logf("Creating test QoS: %s", testQoSName)
			resp, err := adapter.Create(ctx, createReq)
			require.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, testQoSName, resp.QoSName)

			// Verify it was created
			qos, err := adapter.Get(ctx, testQoSName)
			require.NoError(t, err)
			assert.Equal(t, testQoSName, qos.Name)
			assert.Equal(t, "Test QoS created by Go adapter pattern test", qos.Description)

			// Delete it
			t.Logf("Deleting test QoS: %s", testQoSName)
			err = adapter.Delete(ctx, testQoSName)
			require.NoError(t, err)

			// Verify it was deleted
			_, err = adapter.Get(ctx, testQoSName)
			assert.Error(t, err, "QoS should not exist after deletion")
		})
	} else {
		t.Log("Skipping create/delete tests. Set SLURM_ENABLE_WRITE_TESTS=true to enable")
	}
}
