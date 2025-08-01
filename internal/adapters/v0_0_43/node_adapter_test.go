// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeAdapter_ValidateNodeCreate(t *testing.T) {
	adapter := &NodeAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Node"),
	}

	tests := []struct {
		name    string
		node    *types.NodeCreate
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil node",
			node:    nil,
			wantErr: true,
			errMsg:  "node data is required",
		},
		{
			name: "empty name",
			node: &types.NodeCreate{
				Name: "",
			},
			wantErr: true,
			errMsg:  "node name is required",
		},
		{
			name: "invalid memory format",
			node: &types.NodeCreate{
				Name:       "test-node",
				RealMemory: -1,
			},
			wantErr: true,
			errMsg:  "memory must be non-negative",
		},
		{
			name: "invalid cpu count",
			node: &types.NodeCreate{
				Name: "test-node",
				CPUs: -1,
			},
			wantErr: true,
			errMsg:  "CPUs must be non-negative",
		},
		{
			name: "valid basic node",
			node: &types.NodeCreate{
				Name:       "test-node",
				CPUs:       4,
				RealMemory: 8192,
			},
			wantErr: false,
		},
		{
			name: "valid complex node",
			node: &types.NodeCreate{
				Name:       "compute-01",
				CPUs:       32,
				Boards:     1,
				Sockets:    2,
				CoresPerSocket: 8,
				ThreadsPerCore: 2,
				RealMemory: 65536,
				TmpDisk:    1000000,
				Partitions: []string{"compute", "gpu"},
				Features:   []string{"avx2", "sse4"},
				Gres:       []string{"gpu:2"},
				State:      "IDLE",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateNodeCreate(tt.node)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNodeAdapter_ApplyNodeDefaults(t *testing.T) {
	adapter := &NodeAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Node"),
	}

	tests := []struct {
		name     string
		input    *types.NodeCreate
		expected *types.NodeCreate
	}{
		{
			name: "apply defaults to minimal node",
			input: &types.NodeCreate{
				Name: "test-node",
			},
			expected: &types.NodeCreate{
				Name:           "test-node",
				CPUs:           1,
				Boards:         1,
				Sockets:        1,
				CoresPerSocket: 1,
				ThreadsPerCore: 1,
				RealMemory:     1024,
				TmpDisk:        0,
				State:          "UNKNOWN",
				Features:       []string{},
				Partitions:     []string{},
				Gres:           []string{},
			},
		},
		{
			name: "preserve existing values",
			input: &types.NodeCreate{
				Name:           "compute-node",
				CPUs:           16,
				Sockets:        2,
				CoresPerSocket: 8,
				RealMemory:     32768,
				Features:       []string{"avx2"},
				State:          "IDLE",
			},
			expected: &types.NodeCreate{
				Name:           "compute-node",
				CPUs:           16,
				Boards:         1,
				Sockets:        2,
				CoresPerSocket: 8,
				ThreadsPerCore: 1,
				RealMemory:     32768,
				TmpDisk:        0,
				Features:       []string{"avx2"},
				Partitions:     []string{},
				Gres:           []string{},
				State:          "IDLE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.ApplyNodeDefaults(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNodeAdapter_FilterNodeList(t *testing.T) {
	adapter := &NodeAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Node"),
	}

	nodes := []types.Node{
		{
			Name:       "compute-01",
			CPUs:       32,
			RealMemory: 65536,
			Partitions: []string{"compute", "gpu"},
			Features:   []string{"avx2", "gpu"},
			State:      "IDLE",
		},
		{
			Name:       "compute-02",
			CPUs:       16,
			RealMemory: 32768,
			Partitions: []string{"compute"},
			Features:   []string{"avx2"},
			State:      "ALLOCATED",
		},
		{
			Name:       "login-01",
			CPUs:       8,
			RealMemory: 16384,
			Partitions: []string{"login"},
			Features:   []string{"sse4"},
			State:      "IDLE",
		},
	}

	tests := []struct {
		name     string
		opts     *types.NodeListOptions
		expected []string // expected node names
	}{
		{
			name:     "no filters",
			opts:     &types.NodeListOptions{},
			expected: []string{"compute-01", "compute-02", "login-01"},
		},
		{
			name: "filter by names",
			opts: &types.NodeListOptions{
				Names: []string{"compute-01", "login-01"},
			},
			expected: []string{"compute-01", "login-01"},
		},
		{
			name: "filter by state",
			opts: &types.NodeListOptions{
				States: []string{"IDLE"},
			},
			expected: []string{"compute-01", "login-01"},
		},
		{
			name: "filter by partition",
			opts: &types.NodeListOptions{
				Partitions: []string{"compute"},
			},
			expected: []string{"compute-01", "compute-02"},
		},
		{
			name: "filter by feature",
			opts: &types.NodeListOptions{
				Features: []string{"gpu"},
			},
			expected: []string{"compute-01"},
		},
		{
			name: "combined filters",
			opts: &types.NodeListOptions{
				States:     []string{"IDLE"},
				Partitions: []string{"compute"},
			},
			expected: []string{"compute-01"},
		},
		{
			name: "no matches",
			opts: &types.NodeListOptions{
				Names: []string{"nonexistent"},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.FilterNodeList(nodes, tt.opts)
			resultNames := make([]string, len(result))
			for i, node := range result {
				resultNames[i] = node.Name
			}
			assert.Equal(t, tt.expected, resultNames)
		})
	}
}

func TestNodeAdapter_ValidateNodeConfiguration(t *testing.T) {
	adapter := &NodeAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Node"),
	}

	tests := []struct {
		name    string
		node    *types.NodeCreate
		wantErr bool
		errMsg  string
	}{
		{
			name: "invalid CPU topology",
			node: &types.NodeCreate{
				Name:           "test-node",
				CPUs:           8,
				Sockets:        2,
				CoresPerSocket: 3, // 2*3 = 6, but CPUs = 8
				ThreadsPerCore: 1,
			},
			wantErr: true,
			errMsg:  "CPU topology mismatch",
		},
		{
			name: "valid CPU topology",
			node: &types.NodeCreate{
				Name:           "test-node",
				CPUs:           16,
				Sockets:        2,
				CoresPerSocket: 4,
				ThreadsPerCore: 2, // 2*4*2 = 16
			},
			wantErr: false,
		},
		{
			name: "valid minimal configuration",
			node: &types.NodeCreate{
				Name: "minimal-node",
				CPUs: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateNodeConfiguration(tt.node)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNodeAdapter_ParseNodeFeatures(t *testing.T) {
	adapter := &NodeAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Node"),
	}

	tests := []struct {
		name     string
		input    []string
		expected map[string]bool
	}{
		{
			name:     "empty features",
			input:    []string{},
			expected: map[string]bool{},
		},
		{
			name:  "single feature",
			input: []string{"avx2"},
			expected: map[string]bool{
				"avx2": true,
			},
		},
		{
			name:  "multiple features",
			input: []string{"avx2", "sse4", "gpu"},
			expected: map[string]bool{
				"avx2": true,
				"sse4": true,
				"gpu":  true,
			},
		},
		{
			name:  "duplicate features",
			input: []string{"avx2", "avx2", "sse4"},
			expected: map[string]bool{
				"avx2": true,
				"sse4": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.ParseNodeFeatures(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNodeAdapter_ValidateContext(t *testing.T) {
	adapter := &NodeAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Node"),
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

func TestNodeAdapter_CalculateNodeUtilization(t *testing.T) {
	adapter := &NodeAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Node"),
	}

	tests := []struct {
		name          string
		node          *types.Node
		expectedCPU   float64
		expectedMem   float64
	}{
		{
			name: "fully utilized node",
			node: &types.Node{
				Name:         "test-node",
				CPUs:         16,
				AllocCPUs:    16,
				RealMemory:   32768,
				AllocMemory:  32768,
			},
			expectedCPU: 1.0,
			expectedMem: 1.0,
		},
		{
			name: "half utilized node",
			node: &types.Node{
				Name:         "test-node",
				CPUs:         16,
				AllocCPUs:    8,
				RealMemory:   32768,
				AllocMemory:  16384,
			},
			expectedCPU: 0.5,
			expectedMem: 0.5,
		},
		{
			name: "idle node",
			node: &types.Node{
				Name:         "test-node",
				CPUs:         16,
				AllocCPUs:    0,
				RealMemory:   32768,
				AllocMemory:  0,
			},
			expectedCPU: 0.0,
			expectedMem: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpuUtil, memUtil := adapter.CalculateNodeUtilization(tt.node)
			assert.InDelta(t, tt.expectedCPU, cpuUtil, 0.001)
			assert.InDelta(t, tt.expectedMem, memUtil, 0.001)
		})
	}
}
