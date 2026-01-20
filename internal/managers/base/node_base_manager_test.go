// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeBaseManager_New(t *testing.T) {
	manager := NewNodeBaseManager("v0.0.43")
	assert.NotNil(t, manager)
	assert.Equal(t, "v0.0.43", manager.GetVersion())
	assert.Equal(t, "Node", manager.GetResourceType())
}

func TestNodeBaseManager_ValidateNodeUpdate(t *testing.T) {
	manager := NewNodeBaseManager("v0.0.43")

	drainState := types.NodeStateDrained
	resumeState := types.NodeStateResuming

	tests := []struct {
		name    string
		update  *types.NodeUpdate
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil update",
			update:  nil,
			wantErr: true,
			errMsg:  "data is required",
		},
		{
			name: "valid basic update",
			update: &types.NodeUpdate{
				State: &drainState,
			},
			wantErr: false,
		},
		{
			name: "valid complex update",
			update: &types.NodeUpdate{
				State:    &resumeState,
				Reason:   stringPtr("Maintenance complete"),
				Features: []string{"gpu", "high-memory"},
				Gres:     stringPtr("gpu:tesla:2"),
				Comment:  stringPtr("Updated after maintenance"),
				Weight:   int32Ptr(100),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateNodeUpdate(tt.update)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNodeBaseManager_FilterNodeList(t *testing.T) {
	manager := NewNodeBaseManager("v0.0.43")

	nodes := []types.Node{
		{Name: "compute-01", Partitions: []string{"compute"}, State: types.NodeStateIdle},
		{Name: "compute-02", Partitions: []string{"compute"}, State: types.NodeStateAllocated},
		{Name: "gpu-01", Partitions: []string{"gpu"}, State: types.NodeStateIdle},
	}

	tests := []struct {
		name     string
		opts     *types.NodeListOptions
		expected int
	}{
		{
			name:     "no filters",
			opts:     nil,
			expected: 3,
		},
		{
			name: "filter by name",
			opts: &types.NodeListOptions{
				Names: []string{"compute-01"},
			},
			expected: 1,
		},
		{
			name: "filter by partition",
			opts: &types.NodeListOptions{
				Partitions: []string{"compute"},
			},
			expected: 2,
		},
		{
			name: "filter by state",
			opts: &types.NodeListOptions{
				States: []types.NodeState{types.NodeStateIdle},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.FilterNodeList(nodes, tt.opts)
			assert.Len(t, result, tt.expected)
		})
	}
}
