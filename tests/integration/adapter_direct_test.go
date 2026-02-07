//go:build integration
// +build integration

// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jontk/slurm-client/internal/adapters/v0_0_43"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_43"
	types "github.com/jontk/slurm-client/api"
)

// TestAdapterDirectUsage tests the adapter pattern directly without the full client
func TestAdapterDirectUsage(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Return a simple success response
		w.Write([]byte(`{"qos": []}`))
	}))
	defer server.Close()

	// Create API client
	apiClient, err := api.NewClientWithResponses(server.URL)
	require.NoError(t, err)

	// Create adapter
	adapter := v0_0_43.NewQoSAdapter(apiClient)

	t.Run("Test Validation", func(t *testing.T) {
		// Test nil validation
		err := adapter.ValidateQoSCreate(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "QoS creation data is required")

		// Test empty name validation
		err = adapter.ValidateQoSCreate(&types.QoSCreate{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "QoS name is required")

		// Test negative priority validation
		err = adapter.ValidateQoSCreate(&types.QoSCreate{
			Name:     "test",
			Priority: -1,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be non-negative")

		// Test valid QoS
		err = adapter.ValidateQoSCreate(&types.QoSCreate{
			Name:        "test",
			Priority:    100,
			UsageFactor: 1.0,
		})
		assert.NoError(t, err)
	})

	// Note: ApplyQoSDefaults and FilterQoSList methods don't exist in the adapter
	// These were likely planned features that were never implemented
	// Commenting out for now

	t.Run("Test List Call", func(t *testing.T) {
		ctx := context.Background()
		_, err := adapter.List(ctx, &types.QoSListOptions{
			Limit: 10,
		})
		// The mock server returns an empty list, but the call should succeed
		assert.NoError(t, err)
	})
}
