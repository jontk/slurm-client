// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPartitionBaseManager_New(t *testing.T) {
	manager := NewPartitionBaseManager("v0.0.43")
	assert.NotNil(t, manager)
	assert.Equal(t, "v0.0.43", manager.GetVersion())
	assert.Equal(t, "Partition", manager.GetResourceType())
}

func TestPartitionBaseManager_ValidatePartitionCreate(t *testing.T) {
	manager := NewPartitionBaseManager("v0.0.43")

	tests := []struct {
		name      string
		partition *types.PartitionCreate
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "nil partition",
			partition: nil,
			wantErr:   true,
			errMsg:    "data is required",
		},
		{
			name: "empty name",
			partition: &types.PartitionCreate{
				Name: "",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "valid partition",
			partition: &types.PartitionCreate{
				Name:       "compute",
				MaxTime:    3600,
				DefaultTime: 1800,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidatePartitionCreate(tt.partition)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPartitionBaseManager_ValidatePartitionUpdate(t *testing.T) {
	manager := NewPartitionBaseManager("v0.0.43")

	tests := []struct {
		name    string
		update  *types.PartitionUpdate
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
			name: "valid update",
			update: &types.PartitionUpdate{
				MaxTime: int32Ptr(7200),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidatePartitionUpdate(tt.update)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPartitionBaseManager_FilterPartitionList(t *testing.T) {
	manager := NewPartitionBaseManager("v0.0.43")

	partitions := []types.Partition{
		{Name: "compute", State: types.PartitionStateUp, MaxTime: 3600},
		{Name: "gpu", State: types.PartitionStateUp, MaxTime: 7200},
		{Name: "debug", State: types.PartitionStateDown, MaxTime: 1800},
	}

	tests := []struct {
		name     string
		opts     *types.PartitionListOptions
		expected int
	}{
		{
			name:     "no filters",
			opts:     nil,
			expected: 3,
		},
		{
			name: "filter by name",
			opts: &types.PartitionListOptions{
				Names: []string{"compute"},
			},
			expected: 1,
		},
		{
			name: "filter by state",
			opts: &types.PartitionListOptions{
				States: []types.PartitionState{types.PartitionStateUp},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.FilterPartitionList(partitions, tt.opts)
			assert.Len(t, result, tt.expected)
		})
	}
}