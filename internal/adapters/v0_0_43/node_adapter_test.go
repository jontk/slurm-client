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

func TestNewNodeAdapter(t *testing.T) {
	adapter := NewNodeAdapter(&api.ClientWithResponses{})
	assert.NotNil(t, adapter)
	assert.NotNil(t, adapter.BaseManager)
}

func TestNodeAdapter_ValidateContext(t *testing.T) {
	adapter := NewNodeAdapter(&api.ClientWithResponses{})

	// Test nil context
	//lint:ignore SA1012 intentionally testing nil context validation
	err := adapter.ValidateContext(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context is required")

	// Test valid context
	err = adapter.ValidateContext(context.Background())
	assert.NoError(t, err)
}

func TestNodeAdapter_List(t *testing.T) {
	adapter := NewNodeAdapter(nil) // Use nil client for testing validation logic

	// Test client initialization check (nil context validation is covered in TestNodeAdapter_ValidateContext)
	_, err := adapter.List(context.TODO(), &types.NodeListOptions{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
}

func TestNodeAdapter_Get(t *testing.T) {
	adapter := NewNodeAdapter(nil)

	// Test empty node name
	_, err := adapter.Get(context.TODO(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "node name is required")

	// Test client initialization check (nil context validation is covered in TestNodeAdapter_ValidateContext)
	_, err = adapter.Get(context.TODO(), "test-node")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
}

func TestNodeAdapter_ConvertAPINodeToCommon(t *testing.T) {
	adapter := NewNodeAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name         string
		apiNode      api.V0043Node
		expectedName string
		expectedCpus int32
	}{
		{
			name: "full node",
			apiNode: api.V0043Node{
				Name:         ptrString("node01"),
				Cpus:         ptrInt32(16),
				RealMemory:   ptrInt64(32768),
				Architecture: ptrString("x86_64"),
			},
			expectedName: "node01",
			expectedCpus: 16,
		},
		{
			name: "minimal node",
			apiNode: api.V0043Node{
				Name: ptrString("node02"),
			},
			expectedName: "node02",
			expectedCpus: 0,
		},
		{
			name:         "empty node",
			apiNode:      api.V0043Node{},
			expectedName: "",
			expectedCpus: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.convertAPINodeToCommon(tt.apiNode)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedName, result.Name)
			assert.Equal(t, tt.expectedCpus, result.CPUs)
		})
	}
}

func TestNodeAdapter_Update(t *testing.T) {
	adapter := NewNodeAdapter(nil)

	// Test nil update
	err := adapter.Update(context.TODO(), "test-node", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "node update data is required")

	// Test empty node name
	err = adapter.Update(context.TODO(), "", &types.NodeUpdate{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "node name is required")

	// Test client initialization check (nil context validation is covered in TestNodeAdapter_ValidateContext)
	err = adapter.Update(context.TODO(), "test-node", &types.NodeUpdate{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
}

func TestNodeAdapter_Delete(t *testing.T) {
	adapter := NewNodeAdapter(nil)

	// Test empty node name
	err := adapter.Delete(context.TODO(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "node name is required")

	// Test client initialization check (nil context validation is covered in TestNodeAdapter_ValidateContext)
	err = adapter.Delete(context.TODO(), "test-node")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
}

func TestNodeAdapter_ValidateNodeUpdate(t *testing.T) {
	adapter := NewNodeAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name          string
		update        *types.NodeUpdate
		expectedError bool
		errorContains string
	}{
		{
			name: "valid update",
			update: &types.NodeUpdate{
				Reason: ptrString("Test reason"),
			},
			expectedError: false,
		},
		{
			name:          "nil update",
			update:        nil,
			expectedError: true,
			errorContains: "node update data is required",
		},
		{
			name:          "empty update",
			update:        &types.NodeUpdate{},
			expectedError: false, // Empty updates are allowed
		},
		{
			name: "update with reason",
			update: &types.NodeUpdate{
				Reason: ptrString("Hardware failure"),
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.validateNodeUpdate(tt.update)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
