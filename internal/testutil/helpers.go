package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/jontk/slurm-client/internal/adapters/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/testutil/mocks"
	"github.com/stretchr/testify/require"
)

// TestContext creates a context with a reasonable timeout for tests
func TestContext(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	return ctx
}

// RequireErrorContains asserts that an error occurred and contains the expected message
func RequireErrorContains(t *testing.T, err error, contains string) {
	require.Error(t, err)
	require.Contains(t, err.Error(), contains)
}

// AssertNoError asserts that no error occurred with a helpful message
func AssertNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	if err != nil {
		if len(msgAndArgs) > 0 {
			require.NoError(t, err, msgAndArgs...)
		} else {
			require.NoError(t, err, "Unexpected error")
		}
	}
}

// SetupMockAdapter creates a mock adapter with test data
func SetupMockAdapter(t *testing.T, version string) *mocks.MockVersionAdapter {
	adapter := mocks.NewMockVersionAdapter(version)
	require.NotNil(t, adapter)
	return adapter
}

// CompareQoSFields compares two QoS objects field by field for testing
type QoSComparison struct {
	IgnoreTimestamps bool
	IgnoreID         bool
}

// Ptr returns a pointer to the given value (generic helper)
func Ptr[T any](v T) *T {
	return &v
}

// PtrSlice converts a slice of values to a slice of pointers
func PtrSlice[T any](values []T) []*T {
	result := make([]*T, len(values))
	for i, v := range values {
		v := v // capture loop variable
		result[i] = &v
	}
	return result
}

// TestAdapterBehavior runs a standard set of adapter behavior tests
func TestAdapterBehavior(t *testing.T, adapter common.VersionAdapter) {
	ctx := TestContext(t)
	qosManager := adapter.GetQoSManager()
	
	t.Run("List empty", func(t *testing.T) {
		// This assumes the adapter starts empty or we can clear it
		// In real tests, you might want to set up a clean state
		list, err := qosManager.List(ctx, nil)
		AssertNoError(t, err)
		require.NotNil(t, list)
	})
	
	t.Run("Create and Get", func(t *testing.T) {
		createReq := &types.QoSCreate{
			Name:        "test-qos-" + time.Now().Format("20060102150405"),
			Description: "Test QoS",
			Priority:    100,
		}
		
		resp, err := qosManager.Create(ctx, createReq)
		AssertNoError(t, err)
		require.Equal(t, createReq.Name, resp.QoSName)
		
		// Get the created QoS
		qos, err := qosManager.Get(ctx, createReq.Name)
		AssertNoError(t, err)
		require.Equal(t, createReq.Name, qos.Name)
		require.Equal(t, createReq.Description, qos.Description)
		require.Equal(t, createReq.Priority, qos.Priority)
	})
	
	t.Run("Update", func(t *testing.T) {
		// Create a QoS first
		name := "update-test-" + time.Now().Format("20060102150405")
		createReq := &types.QoSCreate{
			Name:        name,
			Description: "Original description",
			Priority:    50,
		}
		
		_, err := qosManager.Create(ctx, createReq)
		AssertNoError(t, err)
		
		// Update it
		updateReq := &types.QoSUpdate{
			Description: Ptr("Updated description"),
			Priority:    Ptr(100),
		}
		
		err = qosManager.Update(ctx, name, updateReq)
		AssertNoError(t, err)
		
		// Verify the update
		updated, err := qosManager.Get(ctx, name)
		AssertNoError(t, err)
		require.Equal(t, "Updated description", updated.Description)
		require.Equal(t, 100, updated.Priority)
	})
	
	t.Run("Delete", func(t *testing.T) {
		// Create a QoS first
		name := "delete-test-" + time.Now().Format("20060102150405")
		createReq := &types.QoSCreate{
			Name: name,
		}
		
		_, err := qosManager.Create(ctx, createReq)
		AssertNoError(t, err)
		
		// Delete it
		err = qosManager.Delete(ctx, name)
		AssertNoError(t, err)
		
		// Verify it's gone
		_, err = qosManager.Get(ctx, name)
		require.Error(t, err)
	})
}

// TimePtr returns a pointer to a time.Time
func TimePtr(t time.Time) *time.Time {
	return &t
}

// Now returns the current time for test consistency
func Now() time.Time {
	// Round to nearest second for test consistency
	return time.Now().Round(time.Second)
}

// IntPtr returns a pointer to the given int value
func IntPtr(i int) *int {
	return &i
}

// Int32Ptr returns a pointer to the given int32 value
func Int32Ptr(i int32) *int32 {
	return &i
}

// Int64Ptr returns a pointer to the given int64 value
func Int64Ptr(i int64) *int64 {
	return &i
}

// StringPtr returns a pointer to the given string value
func StringPtr(s string) *string {
	return &s
}

// BoolPtr returns a pointer to the given bool value
func BoolPtr(b bool) *bool {
	return &b
}

// Float32Ptr returns a pointer to the given float32 value
func Float32Ptr(f float32) *float32 {
	return &f
}

// Float64Ptr returns a pointer to the given float64 value
func Float64Ptr(f float64) *float64 {
	return &f
}