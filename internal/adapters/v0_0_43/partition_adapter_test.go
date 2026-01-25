// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	"testing"

	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/stretchr/testify/assert"
)

func TestNewPartitionAdapter(t *testing.T) {
	adapter := NewPartitionAdapter(&api.ClientWithResponses{})
	assert.NotNil(t, adapter)
	assert.NotNil(t, adapter.BaseManager)
}

func TestPartitionAdapter_ValidateContext(t *testing.T) {
	adapter := NewPartitionAdapter(&api.ClientWithResponses{})

	// Test nil context
	//lint:ignore SA1012 intentionally testing nil context validation
	err := adapter.ValidateContext(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context is required")

	// Test valid context
	err = adapter.ValidateContext(context.Background())
	assert.NoError(t, err)
}

func TestPartitionAdapter_List(t *testing.T) {
	adapter := NewPartitionAdapter(nil) // Use nil client for testing validation logic

	// Test client initialization check (nil context validation is covered in TestPartitionAdapter_ValidateContext)
	_, err := adapter.List(context.TODO(), &types.PartitionListOptions{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
}

func TestPartitionAdapter_Get(t *testing.T) {
	adapter := NewPartitionAdapter(nil)

	// Test empty partition name
	_, err := adapter.Get(context.TODO(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "partition name is required")

	// Test client initialization check (nil context validation is covered in TestPartitionAdapter_ValidateContext)
	_, err = adapter.Get(context.TODO(), "test-partition")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
}

func TestPartitionAdapter_ConvertAPIPartitionToCommon(t *testing.T) {
	adapter := NewPartitionAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name          string
		apiPartition  api.V0043PartitionInfo
		expectedName  string
		expectedState types.PartitionState
	}{
		{
			name: "full partition",
			apiPartition: api.V0043PartitionInfo{
				Name: ptrString("compute"),
				Cpus: &struct {
					TaskBinding *int32 `json:"task_binding,omitempty"`
					Total       *int32 `json:"total,omitempty"`
				}{
					Total: ptrInt32(10),
				},
			},
			expectedName:  "compute",
			expectedState: "",
		},
		{
			name: "minimal partition",
			apiPartition: api.V0043PartitionInfo{
				Name: ptrString("debug"),
			},
			expectedName:  "debug",
			expectedState: "",
		},
		{
			name:          "empty partition",
			apiPartition:  api.V0043PartitionInfo{},
			expectedName:  "",
			expectedState: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.convertAPIPartitionToCommon(tt.apiPartition)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedName, result.Name)
			assert.Equal(t, tt.expectedState, result.State)
		})
	}
}

// Note: v0.0.43 API doesn't support partition creation, update, or deletion operations
// so we only test the validation logic for these methods

func TestPartitionAdapter_Create_NotSupported(t *testing.T) {
	adapter := NewPartitionAdapter(&api.ClientWithResponses{})

	// v0.0.43 doesn't support partition creation
	_, err := adapter.Create(context.Background(), &types.PartitionCreate{
		Name: "test-partition",
	})
	assert.Error(t, err)
	// The error should indicate that the operation is not supported or not implemented
}

func TestPartitionAdapter_Update_NotSupported(t *testing.T) {
	adapter := NewPartitionAdapter(&api.ClientWithResponses{})

	// v0.0.43 doesn't support partition updates
	err := adapter.Update(context.Background(), "test-partition", &types.PartitionUpdate{})
	assert.Error(t, err)
	// The error should indicate that the operation is not supported or not implemented
}

func TestPartitionAdapter_Delete_NotSupported(t *testing.T) {
	adapter := NewPartitionAdapter(&api.ClientWithResponses{})

	// v0.0.43 doesn't support partition deletion
	err := adapter.Delete(context.Background(), "test-partition")
	assert.Error(t, err)
	// The error should indicate that the operation is not supported or not implemented
}
