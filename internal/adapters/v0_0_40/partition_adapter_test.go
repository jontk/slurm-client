// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"context"
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPartitionAdapter_ValidatePartitionCreate(t *testing.T) {
	adapter := &PartitionAdapter{
		BaseManager: base.NewBaseManager("v0.0.40", "Partition"),
	}

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
			errMsg:    "partition data is required",
		},
		{
			name: "empty name",
			partition: &types.PartitionCreate{
				Name: "",
			},
			wantErr: true,
			errMsg:  "partition name is required",
		},
		{
			name: "invalid max time",
			partition: &types.PartitionCreate{
				Name:    "test-partition",
				MaxTime: "invalid",
			},
			wantErr: true,
			errMsg:  "invalid time format",
		},
		{
			name: "invalid default time",
			partition: &types.PartitionCreate{
				Name:        "test-partition",
				DefaultTime: "25:00:00", // invalid hours
			},
			wantErr: true,
			errMsg:  "invalid time format",
		},
		{
			name: "negative max CPUs per node",
			partition: &types.PartitionCreate{
				Name:           "test-partition",
				MaxCPUsPerNode: -1,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "negative max nodes",
			partition: &types.PartitionCreate{
				Name:     "test-partition",
				MaxNodes: -1,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "valid basic partition",
			partition: &types.PartitionCreate{
				Name:  "compute",
				Nodes: []string{"node001", "node002"},
			},
			wantErr: false,
		},
		{
			name: "valid complex partition",
			partition: &types.PartitionCreate{
				Name:           "gpu-partition",
				Nodes:          []string{"gpu001", "gpu002", "gpu003"},
				State:          "UP",
				DefaultTime:    "01:00:00",
				MaxTime:        "24:00:00",
				MaxCPUsPerNode: 64,
				MaxNodes:       10,
				MinNodes:       1,
				Priority:       1000,
				RootOnly:       false,
				Shared:         true,
				DefaultQoS:     "gpu",
				AllowedQoS:     []string{"gpu", "high"},
				DenyAccounts:   []string{"guest"},
				AllowAccounts:  []string{"physics", "chemistry"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidatePartitionCreate(tt.partition)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPartitionAdapter_ApplyPartitionDefaults(t *testing.T) {
	adapter := &PartitionAdapter{
		BaseManager: base.NewBaseManager("v0.0.40", "Partition"),
	}

	tests := []struct {
		name     string
		input    *types.PartitionCreate
		expected *types.PartitionCreate
	}{
		{
			name: "apply defaults to minimal partition",
			input: &types.PartitionCreate{
				Name: "test-partition",
			},
			expected: &types.PartitionCreate{
				Name:           "test-partition",
				Nodes:          []string{},
				State:          "UP",
				DefaultTime:    "",
				MaxTime:        "",
				MaxCPUsPerNode: 0,
				MaxNodes:       0,
				MinNodes:       0,
				Priority:       1,
				RootOnly:       false,
				Shared:         false,
				DefaultQoS:     "",
				AllowedQoS:     []string{},
				DenyAccounts:   []string{},
				AllowAccounts:  []string{},
				DenyUsers:      []string{},
				AllowUsers:     []string{},
			},
		},
		{
			name: "preserve existing values",
			input: &types.PartitionCreate{
				Name:           "gpu-partition",
				Nodes:          []string{"gpu001"},
				State:          "DRAIN",
				DefaultTime:    "02:00:00",
				MaxTime:        "12:00:00",
				MaxCPUsPerNode: 32,
				Priority:       500,
				RootOnly:       true,
				DefaultQoS:     "gpu",
				AllowedQoS:     []string{"gpu", "normal"},
			},
			expected: &types.PartitionCreate{
				Name:           "gpu-partition",
				Nodes:          []string{"gpu001"},
				State:          "DRAIN",
				DefaultTime:    "02:00:00",
				MaxTime:        "12:00:00",
				MaxCPUsPerNode: 32,
				MaxNodes:       0,
				MinNodes:       0,
				Priority:       500,
				RootOnly:       true,
				Shared:         false,
				DefaultQoS:     "gpu",
				AllowedQoS:     []string{"gpu", "normal"},
				DenyAccounts:   []string{},
				AllowAccounts:  []string{},
				DenyUsers:      []string{},
				AllowUsers:     []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.ApplyPartitionDefaults(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPartitionAdapter_FilterPartitionList(t *testing.T) {
	adapter := &PartitionAdapter{
		BaseManager: base.NewBaseManager("v0.0.40", "Partition"),
	}

	partitions := []types.Partition{
		{
			Name:        "compute",
			Nodes:       []string{"node001", "node002"},
			State:       "UP",
			Priority:    100,
			DefaultQoS:  "normal",
			AllowedQoS:  []string{"normal", "high"},
			MaxTime:     "24:00:00",
			MaxNodes:    50,
			RootOnly:    false,
		},
		{
			Name:        "gpu",
			Nodes:       []string{"gpu001", "gpu002"},
			State:       "UP",
			Priority:    200,
			DefaultQoS:  "gpu",
			AllowedQoS:  []string{"gpu", "high"},
			MaxTime:     "12:00:00",
			MaxNodes:    10,
			RootOnly:    false,
		},
		{
			Name:        "debug",
			Nodes:       []string{"debug001"},
			State:       "UP",
			Priority:    50,
			DefaultQoS:  "debug",
			AllowedQoS:  []string{"debug"},
			MaxTime:     "00:30:00",
			MaxNodes:    1,
			RootOnly:    false,
		},
		{
			Name:        "admin",
			Nodes:       []string{"admin001"},
			State:       "DRAIN",
			Priority:    1000,
			DefaultQoS:  "admin",
			AllowedQoS:  []string{"admin"},
			MaxTime:     "UNLIMITED",
			MaxNodes:    5,
			RootOnly:    true,
		},
	}

	tests := []struct {
		name     string
		opts     *types.PartitionListOptions
		expected []string // expected partition names
	}{
		{
			name:     "no filters",
			opts:     &types.PartitionListOptions{},
			expected: []string{"compute", "gpu", "debug", "admin"},
		},
		{
			name: "filter by names",
			opts: &types.PartitionListOptions{
				Names: []string{"compute", "gpu"},
			},
			expected: []string{"compute", "gpu"},
		},
		{
			name: "filter by state",
			opts: &types.PartitionListOptions{
				States: []string{"UP"},
			},
			expected: []string{"compute", "gpu", "debug"},
		},
		{
			name: "filter by nodes",
			opts: &types.PartitionListOptions{
				Nodes: []string{"gpu001"},
			},
			expected: []string{"gpu"},
		},
		{
			name: "filter by QoS",
			opts: &types.PartitionListOptions{
				QoSNames: []string{"high"},
			},
			expected: []string{"compute", "gpu"},
		},
		{
			name: "filter by default QoS",
			opts: &types.PartitionListOptions{
				DefaultQoS: []string{"normal"},
			},
			expected: []string{"compute"},
		},
		{
			name: "filter non-root partitions",
			opts: &types.PartitionListOptions{
				ExcludeRootOnly: true,
			},
			expected: []string{"compute", "gpu", "debug"},
		},
		{
			name: "combined filters",
			opts: &types.PartitionListOptions{
				States:   []string{"UP"},
				QoSNames: []string{"high"},
			},
			expected: []string{"compute", "gpu"},
		},
		{
			name: "no matches",
			opts: &types.PartitionListOptions{
				Names: []string{"nonexistent"},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.FilterPartitionList(partitions, tt.opts)
			resultNames := make([]string, len(result))
			for i, partition := range result {
				resultNames[i] = partition.Name
			}
			assert.Equal(t, tt.expected, resultNames)
		})
	}
}

func TestPartitionAdapter_ValidatePartitionNodes(t *testing.T) {
	adapter := &PartitionAdapter{
		BaseManager: base.NewBaseManager("v0.0.40", "Partition"),
	}

	availableNodes := []types.Node{
		{Name: "node001", State: "IDLE"},
		{Name: "node002", State: "ALLOCATED"},
		{Name: "gpu001", State: "IDLE"},
		{Name: "gpu002", State: "DOWN"},
	}

	tests := []struct {
		name      string
		nodeNames []string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "empty node list",
			nodeNames: []string{},
			wantErr:   false,
		},
		{
			name:      "valid existing nodes",
			nodeNames: []string{"node001", "gpu001"},
			wantErr:   false,
		},
		{
			name:      "nonexistent node",
			nodeNames: []string{"nonexistent"},
			wantErr:   true,
			errMsg:    "node not found",
		},
		{
			name:      "mixed valid and invalid nodes",
			nodeNames: []string{"node001", "nonexistent"},
			wantErr:   true,
			errMsg:    "node not found",
		},
		{
			name:      "duplicate nodes",
			nodeNames: []string{"node001", "node001"},
			wantErr:   true,
			errMsg:    "duplicate node",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidatePartitionNodes(tt.nodeNames, availableNodes)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPartitionAdapter_ValidateTimeFormat(t *testing.T) {
	adapter := &PartitionAdapter{
		BaseManager: base.NewBaseManager("v0.0.40", "Partition"),
	}

	tests := []struct {
		name      string
		timeStr   string
		wantErr   bool
		errMsg    string
	}{
		{
			name:    "empty time",
			timeStr: "",
			wantErr: false,
		},
		{
			name:    "unlimited time",
			timeStr: "UNLIMITED",
			wantErr: false,
		},
		{
			name:    "minutes format",
			timeStr: "60",
			wantErr: false,
		},
		{
			name:    "hours:minutes format",
			timeStr: "01:30",
			wantErr: false,
		},
		{
			name:    "hours:minutes:seconds format",
			timeStr: "01:30:45",
			wantErr: false,
		},
		{
			name:    "days-hours:minutes:seconds format",
			timeStr: "2-12:30:45",
			wantErr: false,
		},
		{
			name:    "invalid format",
			timeStr: "invalid",
			wantErr: true,
			errMsg:  "invalid time format",
		},
		{
			name:    "negative time",
			timeStr: "-30",
			wantErr: true,
			errMsg:  "time cannot be negative",
		},
		{
			name:    "invalid hours",
			timeStr: "25:00",
			wantErr: true,
			errMsg:  "invalid hours",
		},
		{
			name:    "invalid minutes",
			timeStr: "01:65",
			wantErr: true,
			errMsg:  "invalid minutes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateTimeFormat(tt.timeStr)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPartitionAdapter_CalculatePartitionCapacity(t *testing.T) {
	adapter := &PartitionAdapter{
		BaseManager: base.NewBaseManager("v0.0.40", "Partition"),
	}

	nodes := []types.Node{
		{
			Name:       "node001",
			CPUs:       32,
			RealMemory: 128000,
			State:      "IDLE",
		},
		{
			Name:       "node002",
			CPUs:       16,
			RealMemory: 64000,
			State:      "ALLOCATED",
		},
		{
			Name:       "gpu001",
			CPUs:       64,
			RealMemory: 256000,
			State:      "IDLE",
		},
	}

	tests := []struct {
		name              string
		partition         *types.Partition
		expectedTotalCPUs int
		expectedTotalMem  int
		expectedNodes     int
	}{
		{
			name: "partition with all nodes",
			partition: &types.Partition{
				Name:  "all",
				Nodes: []string{"node001", "node002", "gpu001"},
			},
			expectedTotalCPUs: 112, // 32 + 16 + 64
			expectedTotalMem:  448000, // 128000 + 64000 + 256000
			expectedNodes:     3,
		},
		{
			name: "partition with subset of nodes",
			partition: &types.Partition{
				Name:  "compute",
				Nodes: []string{"node001", "node002"},
			},
			expectedTotalCPUs: 48, // 32 + 16
			expectedTotalMem:  192000, // 128000 + 64000
			expectedNodes:     2,
		},
		{
			name: "partition with single node",
			partition: &types.Partition{
				Name:  "gpu",
				Nodes: []string{"gpu001"},
			},
			expectedTotalCPUs: 64,
			expectedTotalMem:  256000,
			expectedNodes:     1,
		},
		{
			name: "empty partition",
			partition: &types.Partition{
				Name:  "empty",
				Nodes: []string{},
			},
			expectedTotalCPUs: 0,
			expectedTotalMem:  0,
			expectedNodes:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			capacity := adapter.CalculatePartitionCapacity(tt.partition, nodes)
			assert.Equal(t, tt.expectedTotalCPUs, capacity.TotalCPUs)
			assert.Equal(t, tt.expectedTotalMem, capacity.TotalMemory)
			assert.Equal(t, tt.expectedNodes, capacity.TotalNodes)
		})
	}
}

func TestPartitionAdapter_ValidateQoSAccess(t *testing.T) {
	adapter := &PartitionAdapter{
		BaseManager: base.NewBaseManager("v0.0.40", "Partition"),
	}

	partition := &types.Partition{
		Name:        "compute",
		DefaultQoS:  "normal",
		AllowedQoS:  []string{"normal", "high", "low"},
	}

	tests := []struct {
		name    string
		qosName string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "allowed QoS",
			qosName: "normal",
			wantErr: false,
		},
		{
			name:    "another allowed QoS",
			qosName: "high",
			wantErr: false,
		},
		{
			name:    "disallowed QoS",
			qosName: "admin",
			wantErr: true,
			errMsg:  "QoS not allowed",
		},
		{
			name:    "empty QoS (use default)",
			qosName: "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateQoSAccess(partition, tt.qosName)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPartitionAdapter_ValidateContext(t *testing.T) {
	adapter := &PartitionAdapter{
		BaseManager: base.NewBaseManager("v0.0.40", "Partition"),
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
			errMsg:  "context is required",
		},
		{
			name:    "valid context",
			ctx:     context.Background(),
			wantErr: false,
		},
		{
			name:    "context with timeout",
			ctx:     context.TODO(),
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
