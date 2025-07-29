package v0_0_43

import (
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeAdapter_ValidateNodeCreate(t *testing.T) {
	adapter := &NodeAdapter{
		NodeBaseManager: base.NewNodeBaseManager("v0.0.43"),
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
			name: "negative CPUs",
			node: &types.NodeCreate{
				Name: "test-node",
				CPUs: -1,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "negative RealMemory",
			node: &types.NodeCreate{
				Name:       "test-node",
				RealMemory: -1024,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "negative TmpDisk",
			node: &types.NodeCreate{
				Name:    "test-node",
				TmpDisk: -500,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "valid basic node",
			node: &types.NodeCreate{
				Name: "test-node",
				CPUs: 4,
			},
			wantErr: false,
		},
		{
			name: "valid complex node",
			node: &types.NodeCreate{
				Name:         "compute-node-001",
				CPUs:         64,
				RealMemory:   128000,
				TmpDisk:      1000000,
				State:        "IDLE",
				Reason:       "Node ready for jobs",
				Features:     []string{"gpu", "infiniband", "high-memory"},
				Gres:         []string{"gpu:tesla:2", "mps:400"},
				Weight:       100,
				Partitions:   []string{"compute", "gpu"},
				NodeAddr:     "192.168.1.100",
				NodeHostName: "compute-node-001.cluster.local",
				Port:         6818,
				Version:      "23.02.0",
				Arch:         "x86_64",
				OS:           "Linux",
			},
			wantErr: false,
		},
		{
			name: "invalid state",
			node: &types.NodeCreate{
				Name:  "test-node",
				State: "INVALID_STATE",
			},
			wantErr: true,
			errMsg:  "invalid node state",
		},
		{
			name: "negative weight",
			node: &types.NodeCreate{
				Name:   "test-node",
				Weight: -50,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "negative port",
			node: &types.NodeCreate{
				Name: "test-node",
				Port: -6818,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
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
		NodeBaseManager: base.NewNodeBaseManager("v0.0.43"),
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
				Name:       "test-node",
				CPUs:       1,                        // Default CPU count
				RealMemory: 1024,                     // Default 1GB memory
				TmpDisk:    0,                        // No tmp disk by default
				State:      "UNKNOWN",                // Default state
				Features:   []string{},               // Empty features
				Gres:       []string{},               // Empty gres
				Weight:     1,                        // Default weight
				Partitions: []string{},               // Empty partitions
				Port:       6818,                     // Default SLURM port
				Version:    "",                       // Empty version
				Arch:       "x86_64",                 // Default architecture
				OS:         "Linux",                  // Default OS
			},
		},
		{
			name: "preserve existing values",
			input: &types.NodeCreate{
				Name:         "compute-node-001",
				CPUs:         32,
				RealMemory:   64000,
				TmpDisk:      500000,
				State:        "IDLE",
				Reason:       "Node available",
				Features:     []string{"gpu", "infiniband"},
				Gres:         []string{"gpu:tesla:1"},
				Weight:       200,
				Partitions:   []string{"compute"},
				NodeAddr:     "10.0.1.100",
				NodeHostName: "node001.cluster.local",
				Port:         6819,
				Version:      "23.02.0",
				Arch:         "aarch64",
				OS:           "CentOS",
			},
			expected: &types.NodeCreate{
				Name:         "compute-node-001",
				CPUs:         32,
				RealMemory:   64000,
				TmpDisk:      500000,
				State:        "IDLE",
				Reason:       "Node available",
				Features:     []string{"gpu", "infiniband"},
				Gres:         []string{"gpu:tesla:1"},
				Weight:       200,
				Partitions:   []string{"compute"},
				NodeAddr:     "10.0.1.100",
				NodeHostName: "node001.cluster.local",
				Port:         6819,
				Version:      "23.02.0",
				Arch:         "aarch64",
				OS:           "CentOS",
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
		NodeBaseManager: base.NewNodeBaseManager("v0.0.43"),
	}

	nodes := []types.Node{
		{
			Name:       "compute-001",
			CPUs:       32,
			RealMemory: 64000,
			State:      "IDLE",
			Features:   []string{"cpu", "infiniband"},
			Partitions: []string{"compute"},
			Arch:       "x86_64",
			OS:         "Linux",
		},
		{
			Name:       "gpu-001",
			CPUs:       16,
			RealMemory: 32000,
			State:      "ALLOCATED",
			Features:   []string{"gpu", "cuda"},
			Partitions: []string{"gpu"},
			Arch:       "x86_64",
			OS:         "Linux",
		},
		{
			Name:       "bigmem-001",
			CPUs:       64,
			RealMemory: 256000,
			State:      "IDLE",
			Features:   []string{"bigmem", "high-memory"},
			Partitions: []string{"bigmem"},
			Arch:       "x86_64",
			OS:         "Linux",
		},
		{
			Name:       "debug-001",
			CPUs:       8,
			RealMemory: 16000,
			State:      "DOWN",
			Features:   []string{"debug"},
			Partitions: []string{"debug"},
			Arch:       "aarch64",
			OS:         "CentOS",
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
			expected: []string{"compute-001", "gpu-001", "bigmem-001", "debug-001"},
		},
		{
			name: "filter by names",
			opts: &types.NodeListOptions{
				Names: []string{"compute-001", "gpu-001"},
			},
			expected: []string{"compute-001", "gpu-001"},
		},
		{
			name: "filter by state",
			opts: &types.NodeListOptions{
				States: []string{"IDLE"},
			},
			expected: []string{"compute-001", "bigmem-001"},
		},
		{
			name: "filter by features",
			opts: &types.NodeListOptions{
				Features: []string{"gpu"},
			},
			expected: []string{"gpu-001"},
		},
		{
			name: "filter by partitions",
			opts: &types.NodeListOptions{
				Partitions: []string{"compute", "bigmem"},
			},
			expected: []string{"compute-001", "bigmem-001"},
		},
		{
			name: "filter by architecture",
			opts: &types.NodeListOptions{
				Architectures: []string{"aarch64"},
			},
			expected: []string{"debug-001"},
		},
		{
			name: "filter by OS",
			opts: &types.NodeListOptions{
				OperatingSystems: []string{"CentOS"},
			},
			expected: []string{"debug-001"},
		},
		{
			name: "combined filters",
			opts: &types.NodeListOptions{
				States:   []string{"IDLE"},
				Features: []string{"infiniband", "high-memory"},
			},
			expected: []string{"compute-001", "bigmem-001"},
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

func TestNodeAdapter_ValidateNodeState(t *testing.T) {
	adapter := &NodeAdapter{
		NodeBaseManager: base.NewNodeBaseManager("v0.0.43"),
	}

	tests := []struct {
		name    string
		state   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid IDLE state",
			state:   "IDLE",
			wantErr: false,
		},
		{
			name:    "valid ALLOCATED state",
			state:   "ALLOCATED",
			wantErr: false,
		},
		{
			name:    "valid DOWN state",
			state:   "DOWN",
			wantErr: false,
		},
		{
			name:    "valid DRAIN state",
			state:   "DRAIN",
			wantErr: false,
		},
		{
			name:    "valid UNKNOWN state",
			state:   "UNKNOWN",
			wantErr: false,
		},
		{
			name:    "valid MIXED state",
			state:   "MIXED",
			wantErr: false,
		},
		{
			name:    "valid COMPLETING state",
			state:   "COMPLETING",
			wantErr: false,
		},
		{
			name:    "empty state (should use default)",
			state:   "",
			wantErr: false,
		},
		{
			name:    "invalid state",
			state:   "INVALID_STATE",
			wantErr: true,
			errMsg:  "invalid node state",
		},
		{
			name:    "lowercase state",
			state:   "idle",
			wantErr: true,
			errMsg:  "invalid node state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateNodeState(tt.state)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNodeAdapter_ValidateHardwareSpecs(t *testing.T) {
	adapter := &NodeAdapter{
		NodeBaseManager: base.NewNodeBaseManager("v0.0.43"),
	}

	tests := []struct {
		name       string
		cpus       int
		realMemory int
		tmpDisk    int
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "valid hardware specs",
			cpus:       32,
			realMemory: 64000,
			tmpDisk:    500000,
			wantErr:    false,
		},
		{
			name:       "minimal specs",
			cpus:       1,
			realMemory: 1024,
			tmpDisk:    0,
			wantErr:    false,
		},
		{
			name:       "zero CPUs (invalid)",
			cpus:       0,
			realMemory: 64000,
			tmpDisk:    500000,
			wantErr:    true,
			errMsg:     "CPUs must be positive",
		},
		{
			name:       "negative CPUs",
			cpus:       -4,
			realMemory: 64000,
			tmpDisk:    500000,
			wantErr:    true,
			errMsg:     "must be non-negative",
		},
		{
			name:       "negative memory",
			cpus:       32,
			realMemory: -1024,
			tmpDisk:    500000,
			wantErr:    true,
			errMsg:     "must be non-negative",
		},
		{
			name:       "negative tmp disk",
			cpus:       32,
			realMemory: 64000,
			tmpDisk:    -1000,
			wantErr:    true,
			errMsg:     "must be non-negative",
		},
		{
			name:       "very large specs",
			cpus:       256,
			realMemory: 1024000,
			tmpDisk:    10000000,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateHardwareSpecs(tt.cpus, tt.realMemory, tt.tmpDisk)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNodeAdapter_ValidateNetworkSettings(t *testing.T) {
	adapter := &NodeAdapter{
		NodeBaseManager: base.NewNodeBaseManager("v0.0.43"),
	}

	tests := []struct {
		name         string
		nodeAddr     string
		nodeHostName string
		port         int
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "valid network settings",
			nodeAddr:     "192.168.1.100",
			nodeHostName: "node001.cluster.local",
			port:         6818,
			wantErr:      false,
		},
		{
			name:         "valid IPv6 address",
			nodeAddr:     "2001:db8::1",
			nodeHostName: "node001.cluster.local",
			port:         6818,
			wantErr:      false,
		},
		{
			name:         "empty address and hostname (valid for auto-detection)",
			nodeAddr:     "",
			nodeHostName: "",
			port:         6818,
			wantErr:      false,
		},
		{
			name:         "invalid IP address",
			nodeAddr:     "999.999.999.999",
			nodeHostName: "node001.cluster.local",
			port:         6818,
			wantErr:      true,
			errMsg:       "invalid IP address",
		},
		{
			name:         "invalid hostname format",
			nodeAddr:     "192.168.1.100",
			nodeHostName: "node..invalid",
			port:         6818,
			wantErr:      true,
			errMsg:       "invalid hostname",
		},
		{
			name:         "negative port",
			nodeAddr:     "192.168.1.100",
			nodeHostName: "node001.cluster.local",
			port:         -6818,
			wantErr:      true,
			errMsg:       "must be non-negative",
		},
		{
			name:         "port too high",
			nodeAddr:     "192.168.1.100",
			nodeHostName: "node001.cluster.local",
			port:         70000,
			wantErr:      true,
			errMsg:       "port must be between 1 and 65535",
		},
		{
			name:         "valid port range",
			nodeAddr:     "192.168.1.100",
			nodeHostName: "node001.cluster.local",
			port:         65535,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateNetworkSettings(tt.nodeAddr, tt.nodeHostName, tt.port)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNodeAdapter_ValidateFeatures(t *testing.T) {
	adapter := &NodeAdapter{
		NodeBaseManager: base.NewNodeBaseManager("v0.0.43"),
	}

	tests := []struct {
		name     string
		features []string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid features",
			features: []string{"gpu", "infiniband", "high-memory"},
			wantErr:  false,
		},
		{
			name:     "empty features",
			features: []string{},
			wantErr:  false,
		},
		{
			name:     "single feature",
			features: []string{"gpu"},
			wantErr:  false,
		},
		{
			name:     "features with numbers and dashes",
			features: []string{"gpu-tesla", "cuda-11.0", "memory-512gb"},
			wantErr:  false,
		},
		{
			name:     "duplicate features",
			features: []string{"gpu", "gpu", "infiniband"},
			wantErr:  true,
			errMsg:   "duplicate feature",
		},
		{
			name:     "empty feature string",
			features: []string{"gpu", "", "infiniband"},
			wantErr:  true,
			errMsg:   "feature cannot be empty",
		},
		{
			name:     "feature with spaces",
			features: []string{"gpu tesla", "infiniband"},
			wantErr:  true,
			errMsg:   "feature cannot contain spaces",
		},
		{
			name:     "feature with invalid characters",
			features: []string{"gpu@tesla", "infiniband"},
			wantErr:  true,
			errMsg:   "feature contains invalid characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateFeatures(tt.features)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNodeAdapter_ValidateGres(t *testing.T) {
	adapter := &NodeAdapter{
		NodeBaseManager: base.NewNodeBaseManager("v0.0.43"),
	}

	tests := []struct {
		name    string
		gres    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid gres",
			gres:    []string{"gpu:tesla:2", "mps:400"},
			wantErr: false,
		},
		{
			name:    "empty gres",
			gres:    []string{},
			wantErr: false,
		},
		{
			name:    "simple gres",
			gres:    []string{"gpu:1"},
			wantErr: false,
		},
		{
			name:    "complex gres with IDs",
			gres:    []string{"gpu:tesla:2(IDX:0-1)", "mic:2"},
			wantErr: false,
		},
		{
			name:    "duplicate gres types",
			gres:    []string{"gpu:tesla:2", "gpu:tesla:1"},
			wantErr: true,
			errMsg:  "duplicate gres type",
		},
		{
			name:    "empty gres string",
			gres:    []string{"gpu:tesla:2", ""},
			wantErr: true,
			errMsg:  "gres cannot be empty",
		},
		{
			name:    "invalid gres format",
			gres:    []string{"gpu-tesla-2"},
			wantErr: true,
			errMsg:  "invalid gres format",
		},
		{
			name:    "gres with invalid count",
			gres:    []string{"gpu:tesla:-2"},
			wantErr: true,
			errMsg:  "gres count must be positive",
		},
		{
			name:    "gres with zero count",
			gres:    []string{"gpu:tesla:0"},
			wantErr: true,
			errMsg:  "gres count must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateGres(tt.gres)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNodeAdapter_ValidatePartitionAssignments(t *testing.T) {
	adapter := &NodeAdapter{
		NodeBaseManager: base.NewNodeBaseManager("v0.0.43"),
	}

	tests := []struct {
		name       string
		partitions []string
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "valid partitions",
			partitions: []string{"compute", "gpu", "bigmem"},
			wantErr:    false,
		},
		{
			name:       "empty partitions",
			partitions: []string{},
			wantErr:    false,
		},
		{
			name:       "single partition",
			partitions: []string{"compute"},
			wantErr:    false,
		},
		{
			name:       "duplicate partitions",
			partitions: []string{"compute", "compute", "gpu"},
			wantErr:    true,
			errMsg:     "duplicate partition",
		},
		{
			name:       "empty partition string",
			partitions: []string{"compute", "", "gpu"},
			wantErr:    true,
			errMsg:     "partition cannot be empty",
		},
		{
			name:       "partition with spaces",
			partitions: []string{"compute partition", "gpu"},
			wantErr:    true,
			errMsg:     "partition cannot contain spaces",
		},
		{
			name:       "partition with invalid characters",
			partitions: []string{"compute@cluster", "gpu"},
			wantErr:    true,
			errMsg:     "partition contains invalid characters",
		},
		{
			name:       "valid partition names with dashes and numbers",
			partitions: []string{"compute-1", "gpu-v100", "bigmem-512"},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidatePartitionAssignments(tt.partitions)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}