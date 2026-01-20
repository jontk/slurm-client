// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
)

func TestNewPartitionAdapter(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewPartitionAdapter(client)

	assert.NotNil(t, adapter)
	assert.Equal(t, client, adapter.client)
	assert.NotNil(t, adapter.BaseManager)
	assert.Equal(t, "v0.0.42", adapter.GetVersion())
}

func TestPartitionAdapter_ValidateContext(t *testing.T) {
	adapter := &PartitionAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Partition"),
		client:      &api.ClientWithResponses{},
	}

	tests := []struct {
		name    string
		ctx     context.Context
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil context",
			ctx:     nil,
			wantErr: true,
			errMsg:  "context",
		},
		{
			name:    "valid context",
			ctx:     context.Background(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateContext(tt.ctx)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPartitionAdapter_ClientValidation(t *testing.T) {
	// Test nil client validation
	adapter := NewPartitionAdapter(nil)
	ctx := context.Background()

	_, err := adapter.List(ctx, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	_, err = adapter.Get(ctx, "test-partition")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	// Test that non-nil client passes initial validation
	validAdapter := NewPartitionAdapter(&api.ClientWithResponses{})
	assert.NotNil(t, validAdapter.client)
}

func TestPartitionAdapter_ListOptionsHandling(t *testing.T) {
	adapter := NewPartitionAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	tests := []struct {
		name string
		opts *types.PartitionListOptions
	}{
		{
			name: "nil options",
			opts: nil,
		},
		{
			name: "empty options",
			opts: &types.PartitionListOptions{},
		},
		{
			name: "options with names",
			opts: &types.PartitionListOptions{
				Names: []string{"partition1", "partition2"},
			},
		},
		{
			name: "options with states",
			opts: &types.PartitionListOptions{
				States: []types.PartitionState{types.PartitionStateUp, types.PartitionStateDown},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := adapter.List(ctx, tt.opts)
			// Should get client validation error before any option processing
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "client")
		})
	}
}

func TestPartitionAdapter_GetByName(t *testing.T) {
	adapter := NewPartitionAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	tests := []struct {
		name          string
		partitionName string
	}{
		{
			name:          "valid name",
			partitionName: "compute",
		},
		{
			name:          "empty name",
			partitionName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := adapter.Get(ctx, tt.partitionName)
			// Should get client validation error
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "client")
		})
	}
}

func TestPartitionAdapter_ConvertAPIPartitionToCommon(t *testing.T) {
	adapter := NewPartitionAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name         string
		apiPartition api.V0042PartitionInfo
		expected     types.Partition
	}{
		{
			name: "basic partition",
			apiPartition: api.V0042PartitionInfo{
				Name: ptrString("compute"),
				Nodes: &struct {
					AllowedAllocation *string `json:"allowed_allocation,omitempty"`
					Configured        *string `json:"configured,omitempty"`
					Total             *int32  `json:"total,omitempty"`
				}{
					Total: ptrInt32(100),
				},
			},
			expected: types.Partition{
				Name:       "compute",
				State:      types.PartitionStateUp, // Default state from adapter
				TotalNodes: 100,
			},
		},
		{
			name: "partition with configured nodes",
			apiPartition: api.V0042PartitionInfo{
				Name: ptrString("gpu"),
				Nodes: &struct {
					AllowedAllocation *string `json:"allowed_allocation,omitempty"`
					Configured        *string `json:"configured,omitempty"`
					Total             *int32  `json:"total,omitempty"`
				}{
					Configured: ptrString("gpu[1-4]"),
					Total:      ptrInt32(4),
				},
			},
			expected: types.Partition{
				Name:       "gpu",
				State:      types.PartitionStateUp, // Default state from adapter
				TotalNodes: 4,
			},
		},
		{
			name: "minimal partition",
			apiPartition: api.V0042PartitionInfo{
				Name: ptrString("minimal"),
			},
			expected: types.Partition{
				Name:  "minimal",
				State: types.PartitionStateUp, // Default state from adapter
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.convertAPIPartitionToCommon(tt.apiPartition)

			require.NoError(t, err)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.State, result.State)
			assert.Equal(t, tt.expected.TotalNodes, result.TotalNodes)
		})
	}
}

func TestPartitionAdapter_ErrorHandling(t *testing.T) {
	adapter := NewPartitionAdapter(nil)
	ctx := context.Background()

	// Test various error conditions with nil client
	tests := []struct {
		name     string
		testFunc func() error
	}{
		{
			name: "List with nil client",
			testFunc: func() error {
				_, err := adapter.List(ctx, nil)
				return err
			},
		},
		{
			name: "Get with nil client",
			testFunc: func() error {
				_, err := adapter.Get(ctx, "compute")
				return err
			},
		},
		{
			name: "Create with nil client",
			testFunc: func() error {
				_, err := adapter.Create(ctx, &types.PartitionCreate{Name: "test"})
				return err
			},
		},
		{
			name: "Update with nil client",
			testFunc: func() error {
				return adapter.Update(ctx, "test", &types.PartitionUpdate{})
			},
		},
		{
			name: "Delete with nil client",
			testFunc: func() error {
				return adapter.Delete(ctx, "test")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testFunc()
			assert.Error(t, err)
			// Should get either client validation error or "not supported" error
			errorMsg := err.Error()
			assert.True(t,
				strings.Contains(errorMsg, "client") ||
					strings.Contains(errorMsg, "not supported"),
				"Expected client validation or not supported error, got: %v", err)
		})
	}
}
