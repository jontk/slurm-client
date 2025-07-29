package v0_0_43

import (
	"context"
	"testing"

	"github.com/jontk/slurm-client/internal/common/builders"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/testutil"
	"github.com/jontk/slurm-client/internal/testutil/mocks"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestQoSAdapterIntegration tests the complete QoS adapter flow
func TestQoSAdapterIntegration(t *testing.T) {
	// Create mock client
	mockClient := &api.ClientWithResponses{}
	adapter := NewQoSAdapter(mockClient)

	ctx := context.Background()

	t.Run("Create QoS using builder pattern", func(t *testing.T) {
		// Use builder to create QoS
		qosCreate, err := builders.NewQoSBuilder("integration-test").
			WithDescription("Integration test QoS").
			WithPriority(100).
			WithFlags("DenyOnLimit").
			WithUsageFactor(1.5).
			WithUsageThreshold(0.8).
			WithLimits().
				WithMaxCPUsPerUser(50).
				WithMaxJobsPerUser(10).
				WithMaxNodesPerUser(5).
				Done().
			Build()

		require.NoError(t, err)
		require.NotNil(t, qosCreate)

		// Validate the creation request
		err = adapter.ValidateQoSCreate(qosCreate)
		assert.NoError(t, err)

		// Apply defaults
		qosWithDefaults := adapter.ApplyQoSDefaults(qosCreate)
		assert.Equal(t, []string{"DenyOnLimit"}, qosWithDefaults.Flags)
		assert.Equal(t, 1.5, qosWithDefaults.UsageFactor)
	})

	t.Run("Validate QoS update", func(t *testing.T) {
		// Use builder for update
		update, err := builders.NewQoSBuilder("update-test").
			WithDescription("Updated description").
			WithPriority(200).
			BuildForUpdate()

		require.NoError(t, err)
		require.NotNil(t, update)

		// Validate update
		err = adapter.ValidateQoSUpdate(update)
		assert.NoError(t, err)
	})

	t.Run("Filter QoS list", func(t *testing.T) {
		// Sample QoS list
		qosList := []types.QoS{
			{
				Name:            "normal",
				Priority:        100,
				AllowedAccounts: []string{"physics", "chemistry"},
				AllowedUsers:    []string{"user1", "user2"},
			},
			{
				Name:            "high",
				Priority:        1000,
				AllowedAccounts: []string{"physics"},
				AllowedUsers:    []string{"user3"},
			},
			{
				Name:            "batch",
				Priority:        10,
				AllowedAccounts: []string{"chemistry", "biology"},
				AllowedUsers:    []string{"user1", "user4"},
			},
		}

		// Test various filter scenarios
		testCases := []struct {
			name     string
			opts     *types.QoSListOptions
			expected []string
		}{
			{
				name:     "no filter",
				opts:     &types.QoSListOptions{},
				expected: []string{"normal", "high", "batch"},
			},
			{
				name: "filter by name",
				opts: &types.QoSListOptions{
					Names: []string{"normal", "high"},
				},
				expected: []string{"normal", "high"},
			},
			{
				name: "filter by account",
				opts: &types.QoSListOptions{
					Accounts: []string{"physics"},
				},
				expected: []string{"normal", "high"},
			},
			{
				name: "filter by user",
				opts: &types.QoSListOptions{
					Users: []string{"user1"},
				},
				expected: []string{"normal", "batch"},
			},
			{
				name: "combined filters",
				opts: &types.QoSListOptions{
					Accounts: []string{"chemistry"},
					Users:    []string{"user1"},
				},
				expected: []string{"normal", "batch"},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				filtered := adapter.FilterQoSList(qosList, tc.opts)
				resultNames := make([]string, len(filtered))
				for i, qos := range filtered {
					resultNames[i] = qos.Name
				}
				assert.Equal(t, tc.expected, resultNames)
			})
		}
	})

	t.Run("Builder pattern edge cases", func(t *testing.T) {
		// Test high priority preset
		highPriority, err := builders.NewQoSBuilder("high-priority").
			AsHighPriority().
			Build()
		require.NoError(t, err)
		assert.Equal(t, 1000, highPriority.Priority)
		assert.Contains(t, highPriority.Flags, "RequiresReservation")

		// Test batch preset
		batch, err := builders.NewQoSBuilder("batch").
			AsBatchQueue().
			Build()
		require.NoError(t, err)
		assert.Equal(t, 10, batch.Priority)
		assert.Contains(t, batch.Flags, "NoReserve")

		// Test interactive preset
		interactive, err := builders.NewQoSBuilder("interactive").
			AsInteractive().
			Build()
		require.NoError(t, err)
		assert.Equal(t, 500, interactive.Priority)
		assert.Contains(t, interactive.PreemptMode, "suspend")
	})

	t.Run("Validation error scenarios", func(t *testing.T) {
		// Empty name
		_, err := builders.NewQoSBuilder("").Build()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name is required")

		// Negative priority
		_, err = builders.NewQoSBuilder("test").
			WithPriority(-1).
			Build()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be non-negative")

		// Invalid usage factor
		_, err = builders.NewQoSBuilder("test").
			WithUsageFactor(-1.0).
			Build()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be non-negative")

		// Usage threshold out of range
		_, err = builders.NewQoSBuilder("test").
			WithUsageThreshold(1.5).
			Build()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be between 0 and 1")

		// High priority without reservation flag
		_, err = builders.NewQoSBuilder("test").
			WithPriority(1001).
			Build()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "RequiresReservation flag")

		// Conflicting flags
		_, err = builders.NewQoSBuilder("test").
			WithFlags("NoReserve", "RequiresReservation").
			Build()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "conflicting flags")
	})
}

// TestQoSAdapterWithMockAdapter tests the QoS adapter with the mock adapter
func TestQoSAdapterWithMockAdapter(t *testing.T) {
	// Create mock version adapter
	mockAdapter := mocks.NewMockVersionAdapter("v0.0.43")
	qosAdapter := mockAdapter.GetQoSManager()

	ctx := context.Background()

	t.Run("List QoS entries", func(t *testing.T) {
		// List all QoS
		result, err := qosAdapter.List(ctx, nil)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Greater(t, len(result.QoS), 0)

		// List with pagination
		paginatedResult, err := qosAdapter.List(ctx, &types.QoSListOptions{
			Limit:  2,
			Offset: 1,
		})
		require.NoError(t, err)
		assert.Len(t, paginatedResult.QoS, 2)
		assert.Equal(t, result.Total, paginatedResult.Total)
	})

	t.Run("Get specific QoS", func(t *testing.T) {
		qos, err := qosAdapter.Get(ctx, "normal")
		require.NoError(t, err)
		assert.NotNil(t, qos)
		assert.Equal(t, "normal", qos.Name)

		// Non-existent QoS
		_, err = qosAdapter.Get(ctx, "non-existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Create new QoS", func(t *testing.T) {
		newQoS := &types.QoSCreate{
			Name:        "test-create",
			Description: "Test QoS creation",
			Priority:    50,
			Flags:       []string{"DenyOnLimit"},
		}

		response, err := qosAdapter.Create(ctx, newQoS)
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "test-create", response.QoSName)

		// Duplicate creation should fail
		_, err = qosAdapter.Create(ctx, newQoS)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Update QoS", func(t *testing.T) {
		update := &types.QoSUpdate{
			Description: stringPtr("Updated description"),
			Priority:    testutil.IntPtr(150),
		}

		err := qosAdapter.Update(ctx, "normal", update)
		require.NoError(t, err)

		// Verify update
		updated, err := qosAdapter.Get(ctx, "normal")
		require.NoError(t, err)
		assert.Equal(t, "Updated description", updated.Description)
		assert.Equal(t, 150, updated.Priority)

		// Update non-existent QoS
		err = qosAdapter.Update(ctx, "non-existent", update)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Delete QoS", func(t *testing.T) {
		// Create QoS to delete
		toDelete := &types.QoSCreate{
			Name: "test-delete",
		}
		_, err := qosAdapter.Create(ctx, toDelete)
		require.NoError(t, err)

		// Delete it
		err = qosAdapter.Delete(ctx, "test-delete")
		require.NoError(t, err)

		// Verify deletion
		_, err = qosAdapter.Get(ctx, "test-delete")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")

		// Delete non-existent QoS
		err = qosAdapter.Delete(ctx, "non-existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

