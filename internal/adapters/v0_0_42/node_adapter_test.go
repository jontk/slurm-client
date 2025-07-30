package v0_0_42

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
		BaseManager: base.NewBaseManager("v0.0.42", "Node"),
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
			name: "valid basic node",
			node: &types.NodeCreate{
				Name:       "test-node",
				CPUs:       4,
				RealMemory: 8192,
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

func TestNodeAdapter_FilterNodeList(t *testing.T) {
	adapter := &NodeAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Node"),
	}

	nodes := []types.Node{
		{
			Name:       "compute-01",
			CPUs:       32,
			RealMemory: 65536,
			State:      "IDLE",
		},
		{
			Name:       "compute-02",
			CPUs:       16,
			RealMemory: 32768,
			State:      "ALLOCATED",
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
			expected: []string{"compute-01", "compute-02"},
		},
		{
			name: "filter by names",
			opts: &types.NodeListOptions{
				Names: []string{"compute-01"},
			},
			expected: []string{"compute-01"},
		},
		{
			name: "filter by state",
			opts: &types.NodeListOptions{
				States: []string{"IDLE"},
			},
			expected: []string{"compute-01"},
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

func TestNodeAdapter_ValidateContext(t *testing.T) {
	adapter := &NodeAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Node"),
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