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
			errMsg:  "node update data is required",
		},
		{
			name: "empty node name",
			update: &types.NodeUpdate{
				NodeName: "",
			},
			wantErr: true,
			errMsg:  "node name is required",
		},
		{
			name: "valid basic update",
			update: &types.NodeUpdate{
				NodeName: "compute-01",
				State:    stringPtr("DRAIN"),
			},
			wantErr: false,
		},
		{
			name: "valid complex update",
			update: &types.NodeUpdate{
				NodeName:   "compute-01",
				State:      stringPtr("RESUME"),
				Reason:     stringPtr("Maintenance complete"),
				Features:   &[]string{"gpu", "high-memory"},
				Gres:       stringPtr("gpu:tesla:2"),
				Comment:    stringPtr("Updated after maintenance"),
				Weight:     intPtr(100),
			},
			wantErr: false,
		},
		{
			name: "negative weight",
			update: &types.NodeUpdate{
				NodeName: "compute-01",
				Weight:   intPtr(-1),
			},
			wantErr: true,
			errMsg:  "weight must be non-negative",
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
		{
			NodeName:  "compute-01",
			State:     "IDLE",
			Partition: "compute",
			Features:  []string{"gpu", "high-memory"},
			Gres:      "gpu:tesla:2",
			CPUs:      32,
			Memory:    128000,
		},
		{
			NodeName:  "compute-02",
			State:     "ALLOCATED",
			Partition: "compute",
			Features:  []string{"cpu-only"},
			Gres:      "",
			CPUs:      16,
			Memory:    64000,
		},
		{
			NodeName:  "gpu-01",
			State:     "DRAIN",
			Partition: "gpu",
			Features:  []string{"gpu", "high-memory"},
			Gres:      "gpu:v100:4",
			CPUs:      64,
			Memory:    256000,
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
			expected: []string{"compute-01", "compute-02", "gpu-01"},
		},
		{
			name: "filter by names",
			opts: &types.NodeListOptions{
				NodeNames: []string{"compute-01", "gpu-01"},
			},
			expected: []string{"compute-01", "gpu-01"},
		},
		{
			name: "filter by state",
			opts: &types.NodeListOptions{
				States: []string{"IDLE", "ALLOCATED"},
			},
			expected: []string{"compute-01", "compute-02"},
		},
		{
			name: "filter by partition",
			opts: &types.NodeListOptions{
				Partitions: []string{"gpu"},
			},
			expected: []string{"gpu-01"},
		},
		{
			name: "filter by feature",
			opts: &types.NodeListOptions{
				Features: []string{"gpu"},
			},
			expected: []string{"compute-01", "gpu-01"},
		},
		{
			name: "filter by minimum CPUs",
			opts: &types.NodeListOptions{
				MinCPUs: intPtr(32),
			},
			expected: []string{"compute-01", "gpu-01"},
		},
		{
			name: "filter by minimum memory",
			opts: &types.NodeListOptions{
				MinMemory: int64Ptr(100000),
			},
			expected: []string{"compute-01", "gpu-01"},
		},
		{
			name: "combined filters",
			opts: &types.NodeListOptions{
				States:     []string{"IDLE"},
				Partitions: []string{"compute"},
				Features:   []string{"gpu"},
			},
			expected: []string{"compute-01"},
		},
		{
			name: "no matches",
			opts: &types.NodeListOptions{
				NodeNames: []string{"nonexistent"},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.FilterNodeList(nodes, tt.opts)
			resultNames := make([]string, len(result))
			for i, node := range result {
				resultNames[i] = node.NodeName
			}
			assert.Equal(t, tt.expected, resultNames)
		})
	}
}

func TestNodeBaseManager_ValidateNodeName(t *testing.T) {
	manager := NewNodeBaseManager("v0.0.43")

	tests := []struct {
		name     string
		nodeName string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid node name",
			nodeName: "compute-01",
			wantErr:  false,
		},
		{
			name:     "empty node name",
			nodeName: "",
			wantErr:  true,
			errMsg:   "node name is required",
		},
		{
			name:     "node name with numbers",
			nodeName: "node123",
			wantErr:  false,
		},
		{
			name:     "node name with hyphens",
			nodeName: "compute-node-01",
			wantErr:  false,
		},
		{
			name:     "node name with underscores",
			nodeName: "compute_node_01",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateNodeName(tt.nodeName)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNodeBaseManager_ValidateNodeState(t *testing.T) {
	manager := NewNodeBaseManager("v0.0.43")

	tests := []struct {
		name      string
		nodeState string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid IDLE state",
			nodeState: "IDLE",
			wantErr:   false,
		},
		{
			name:      "valid ALLOCATED state",
			nodeState: "ALLOCATED",
			wantErr:   false,
		},
		{
			name:      "valid DRAIN state",
			nodeState: "DRAIN",
			wantErr:   false,
		},
		{
			name:      "valid RESUME state",
			nodeState: "RESUME",
			wantErr:   false,
		},
		{
			name:      "valid DOWN state",
			nodeState: "DOWN",
			wantErr:   false,
		},
		{
			name:      "invalid state",
			nodeState: "INVALID_STATE",
			wantErr:   true,
			errMsg:    "invalid node state",
		},
		{
			name:      "empty state",
			nodeState: "",
			wantErr:   false, // Empty state is allowed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateNodeState(tt.nodeState)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNodeBaseManager_ParseNodeList(t *testing.T) {
	manager := NewNodeBaseManager("v0.0.43")

	tests := []struct {
		name     string
		nodeList string
		expected []string
	}{
		{
			name:     "single node",
			nodeList: "compute-01",
			expected: []string{"compute-01"},
		},
		{
			name:     "comma separated nodes",
			nodeList: "compute-01,compute-02,gpu-01",
			expected: []string{"compute-01", "compute-02", "gpu-01"},
		},
		{
			name:     "node range",
			nodeList: "compute-[01-03]",
			expected: []string{"compute-01", "compute-02", "compute-03"},
		},
		{
			name:     "mixed format",
			nodeList: "compute-[01-02],gpu-01",
			expected: []string{"compute-01", "compute-02", "gpu-01"},
		},
		{
			name:     "empty list",
			nodeList: "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.ParseNodeList(tt.nodeList)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNodeBaseManager_ValidateNodeFeatures(t *testing.T) {
	manager := NewNodeBaseManager("v0.0.43")

	availableFeatures := []string{"gpu", "high-memory", "cpu-only", "infiniband"}

	tests := []struct {
		name     string
		features []string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid features",
			features: []string{"gpu", "high-memory"},
			wantErr:  false,
		},
		{
			name:     "empty features",
			features: []string{},
			wantErr:  false,
		},
		{
			name:     "invalid feature",
			features: []string{"nonexistent-feature"},
			wantErr:  true,
			errMsg:   "invalid node feature",
		},
		{
			name:     "mixed valid and invalid",
			features: []string{"gpu", "invalid-feature"},
			wantErr:  true,
			errMsg:   "invalid node feature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateNodeFeatures(tt.features, availableFeatures)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNodeBaseManager_ParseGresString(t *testing.T) {
	manager := NewNodeBaseManager("v0.0.43")

	tests := []struct {
		name     string
		gresStr  string
		expected map[string]int
	}{
		{
			name:     "empty GRES",
			gresStr:  "",
			expected: map[string]int{},
		},
		{
			name:     "single GPU type",
			gresStr:  "gpu:tesla:2",
			expected: map[string]int{"gpu:tesla": 2},
		},
		{
			name:     "multiple GRES types",
			gresStr:  "gpu:tesla:2,gpu:v100:4",
			expected: map[string]int{"gpu:tesla": 2, "gpu:v100": 4},
		},
		{
			name:     "GRES without count",
			gresStr:  "gpu:tesla",
			expected: map[string]int{"gpu:tesla": 1},
		},
		{
			name:     "mixed GRES",
			gresStr:  "gpu:tesla:2,mic:1,bandwidth:lustre:10",
			expected: map[string]int{"gpu:tesla": 2, "mic": 1, "bandwidth:lustre": 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.ParseGresString(tt.gresStr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNodeBaseManager_CalculateNodeMetrics(t *testing.T) {
	manager := NewNodeBaseManager("v0.0.43")

	nodes := []types.Node{
		{
			NodeName: "compute-01",
			State:    "IDLE",
			CPUs:     32,
			Memory:   128000,
		},
		{
			NodeName: "compute-02",
			State:    "ALLOCATED",
			CPUs:     16,
			Memory:   64000,
		},
		{
			NodeName: "gpu-01",
			State:    "DRAIN",
			CPUs:     64,
			Memory:   256000,
		},
	}

	metrics := manager.CalculateNodeMetrics(nodes)

	assert.Equal(t, 3, metrics.TotalNodes)
	assert.Equal(t, 1, metrics.IdleNodes)
	assert.Equal(t, 1, metrics.AllocatedNodes)
	assert.Equal(t, 1, metrics.DrainNodes)
	assert.Equal(t, 112, metrics.TotalCPUs) // 32 + 16 + 64
	assert.Equal(t, int64(448000), metrics.TotalMemory) // 128000 + 64000 + 256000
	assert.Equal(t, 32, metrics.AvailableCPUs) // Only IDLE nodes
	assert.Equal(t, int64(128000), metrics.AvailableMemory) // Only IDLE nodes
}

