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
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
	"github.com/jontk/slurm-client/internal/common/types"
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
		assert.Contains(t, err.Error(), "QoS data is required")

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

	t.Run("Test Defaults", func(t *testing.T) {
		qos := &types.QoSCreate{
			Name: "test",
		}
		
		result := adapter.ApplyQoSDefaults(qos)
		assert.Equal(t, "test", result.Name)
		assert.Equal(t, 0, result.Priority)
		assert.Equal(t, 1.0, result.UsageFactor)
		assert.Equal(t, 0.0, result.UsageThreshold)
		assert.NotNil(t, result.Flags)
		assert.Empty(t, result.Flags)
	})

	t.Run("Test Filtering", func(t *testing.T) {
		qosList := []types.QoS{
			{
				Name:            "normal",
				AllowedAccounts: []string{"physics", "chemistry"},
				AllowedUsers:    []string{"user1", "user2"},
			},
			{
				Name:            "high",
				AllowedAccounts: []string{"physics"},
				AllowedUsers:    []string{"user3"},
			},
			{
				Name:            "low",
				AllowedAccounts: []string{"biology"},
				AllowedUsers:    []string{"user1", "user4"},
			},
		}

		// Filter by account
		result := adapter.FilterQoSList(qosList, &types.QoSListOptions{
			Accounts: []string{"physics"},
		})
		assert.Len(t, result, 2)
		assert.Equal(t, "normal", result[0].Name)
		assert.Equal(t, "high", result[1].Name)

		// Filter by user
		result = adapter.FilterQoSList(qosList, &types.QoSListOptions{
			Users: []string{"user1"},
		})
		assert.Len(t, result, 2)
		assert.Equal(t, "normal", result[0].Name)
		assert.Equal(t, "low", result[1].Name)
	})

	t.Run("Test List Call", func(t *testing.T) {
		ctx := context.Background()
		_, err := adapter.List(ctx, &types.QoSListOptions{
			Limit: 10,
		})
		// The mock server returns an empty list, but the call should succeed
		assert.NoError(t, err)
	})
}
