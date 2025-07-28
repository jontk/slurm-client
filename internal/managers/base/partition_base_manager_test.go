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
			name: "negative max time",
			partition: &types.PartitionCreate{
				Name:    "test-partition",
				MaxTime: -1,
			},
			wantErr: true,
			errMsg:  "max time must be non-negative",
		},
		{
			name: "negative default time",
			partition: &types.PartitionCreate{
				Name:        "test-partition",
				DefaultTime: -1,
			},
			wantErr: true,
			errMsg:  "default time must be non-negative",
		},
		{
			name: "negative max nodes",
			partition: &types.PartitionCreate{
				Name:     "test-partition",
				MaxNodes: -1,
			},
			wantErr: true,
			errMsg:  "max nodes must be non-negative",
		},
		{
			name: "negative min nodes",
			partition: &types.PartitionCreate{
				Name:     "test-partition",
				MinNodes: -1,
			},
			wantErr: true,
			errMsg:  "min nodes must be non-negative",
		},
		{
			name: "negative max CPUs per node",
			partition: &types.PartitionCreate{
				Name:          "test-partition",
				MaxCPUsPerNode: -1,
			},
			wantErr: true,
			errMsg:  "max CPUs per node must be non-negative",
		},
		{
			name: "valid basic partition",
			partition: &types.PartitionCreate{
				Name: "test-partition",
			},
			wantErr: false,
		},
		{
			name: "valid complex partition",
			partition: &types.PartitionCreate{
				Name:           "complex-partition",
				Nodes:          []string{"compute-[01-10]"},
				State:          "UP",
				MaxTime:        1440, // 24 hours
				DefaultTime:    60,   // 1 hour
				MaxNodes:       10,
				MinNodes:       1,
				MaxCPUsPerNode: 32,
				Priority:       100,
				QoS:            "normal",
				AllowGroups:    []string{"users", "admins"},
				DenyGroups:     []string{"guests"},
				AllowAccounts:  []string{"account1", "account2"},
				DenyAccounts:   []string{"restricted"},
				AllowQoS:       []string{"normal", "high"},
				DenyQoS:        []string{"debug"},
				Flags:          []string{"Default", "Hidden"},
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
			errMsg:  "partition update data is required",
		},
		{
			name: "empty name",
			update: &types.PartitionUpdate{
				Name: "",
			},
			wantErr: true,
			errMsg:  "partition name is required",
		},
		{
			name: "negative max time",
			update: &types.PartitionUpdate{
				Name:    "test-partition",
				MaxTime: intPtr(-1),
			},
			wantErr: true,
			errMsg:  "max time must be non-negative",
		},
		{
			name: "valid update",
			update: &types.PartitionUpdate{
				Name:    "test-partition",
				State:   stringPtr("DOWN"),
				MaxTime: intPtr(2880), // 48 hours
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

func TestPartitionBaseManager_ApplyPartitionDefaults(t *testing.T) {
	manager := NewPartitionBaseManager("v0.0.43")

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
				Name:        "test-partition",
				State:       "UP",        // Default state
				MaxTime:     1440,        // Default 24 hours
				DefaultTime: 60,          // Default 1 hour
				MaxNodes:    1,           // Default max nodes
				MinNodes:    0,           // Default min nodes
				Priority:    1,           // Default priority
				Flags:       []string{},  // Empty flags
				AllowGroups: []string{},  // Empty groups
				DenyGroups:  []string{},  // Empty groups
			},
		},
		{
			name: "preserve existing values",
			input: &types.PartitionCreate{
				Name:           "test-partition",
				State:          "DOWN",
				MaxTime:        2880,
				DefaultTime:    120,
				MaxNodes:       100,
				MinNodes:       5,
				MaxCPUsPerNode: 64,
				Priority:       500,
				QoS:            "high",
				Flags:          []string{"Default", "Hidden"},
			},
			expected: &types.PartitionCreate{
				Name:           "test-partition",
				State:          "DOWN",
				MaxTime:        2880,
				DefaultTime:    120,
				MaxNodes:       100,
				MinNodes:       5,
				MaxCPUsPerNode: 64,
				Priority:       500,
				QoS:            "high",
				Flags:          []string{"Default", "Hidden"},
				AllowGroups:    []string{}, // Still apply empty defaults
				DenyGroups:     []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.ApplyPartitionDefaults(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPartitionBaseManager_FilterPartitionList(t *testing.T) {
	manager := NewPartitionBaseManager("v0.0.43")

	partitions := []types.Partition{
		{
			Name:     "compute",
			State:    "UP",
			Nodes:    []string{"compute-[01-10]"},
			Priority: 100,
			QoS:      "normal",
			Flags:    []string{"Default"},
		},
		{
			Name:     "gpu",
			State:    "UP",
			Nodes:    []string{"gpu-[01-02]"},
			Priority: 200,
			QoS:      "high",
			Flags:    []string{"Hidden"},
		},
		{
			Name:     "debug",
			State:    "DOWN",
			Nodes:    []string{"debug-01"},
			Priority: 50,
			QoS:      "debug",
			Flags:    []string{"Debug"},
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
			expected: []string{"compute", "gpu", "debug"},
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
			expected: []string{"compute", "gpu"},
		},
		{
			name: "filter by minimum priority",
			opts: &types.PartitionListOptions{
				MinPriority: intPtr(100),
			},
			expected: []string{"compute", "gpu"},
		},
		{
			name: "filter by QoS",
			opts: &types.PartitionListOptions{
				QoSList: []string{"normal", "high"},
			},
			expected: []string{"compute", "gpu"},
		},
		{
			name: "filter by flag",
			opts: &types.PartitionListOptions{
				WithFlags: []string{"Default"},
			},
			expected: []string{"compute"},
		},
		{
			name: "combined filters",
			opts: &types.PartitionListOptions{
				States:      []string{"UP"},
				MinPriority: intPtr(150),
			},
			expected: []string{"gpu"},
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
			result := manager.FilterPartitionList(partitions, tt.opts)
			resultNames := make([]string, len(result))
			for i, partition := range result {
				resultNames[i] = partition.Name
			}
			assert.Equal(t, tt.expected, resultNames)
		})
	}
}

func TestPartitionBaseManager_ValidatePartitionName(t *testing.T) {
	manager := NewPartitionBaseManager("v0.0.43")

	tests := []struct {
		name          string
		partitionName string
		wantErr       bool
		errMsg        string
	}{
		{
			name:          "valid partition name",
			partitionName: "compute",
			wantErr:       false,
		},
		{
			name:          "empty partition name",
			partitionName: "",
			wantErr:       true,
			errMsg:        "partition name is required",
		},
		{
			name:          "partition name with hyphens",
			partitionName: "high-memory",
			wantErr:       false,
		},
		{
			name:          "partition name with underscores",
			partitionName: "gpu_partition",
			wantErr:       false,
		},
		{
			name:          "partition name with numbers",
			partitionName: "partition1",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidatePartitionName(tt.partitionName)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPartitionBaseManager_ValidatePartitionState(t *testing.T) {
	manager := NewPartitionBaseManager("v0.0.43")

	tests := []struct {
		name  string
		state string
		valid bool
	}{
		{name: "UP state", state: "UP", valid: true},
		{name: "DOWN state", state: "DOWN", valid: true},
		{name: "DRAIN state", state: "DRAIN", valid: true},
		{name: "INACTIVE state", state: "INACTIVE", valid: true},
		{name: "invalid state", state: "INVALID", valid: false},
		{name: "empty state", state: "", valid: true}, // Empty is allowed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidatePartitionState(tt.state)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid partition state")
			}
		})
	}
}

func TestPartitionBaseManager_ValidateNodeList(t *testing.T) {
	manager := NewPartitionBaseManager("v0.0.43")

	// Mock available nodes
	availableNodes := []string{"compute-01", "compute-02", "gpu-01", "debug-01"}

	tests := []struct {
		name      string
		nodeList  []string
		wantErr   bool
		errMsg    string
	}{
		{
			name:     "valid nodes",
			nodeList: []string{"compute-01", "gpu-01"},
			wantErr:  false,
		},
		{
			name:     "empty node list",
			nodeList: []string{},
			wantErr:  false, // Empty is allowed
		},
		{
			name:     "invalid node",
			nodeList: []string{"nonexistent-node"},
			wantErr:  true,
			errMsg:   "node does not exist",
		},
		{
			name:     "mixed valid and invalid",
			nodeList: []string{"compute-01", "nonexistent"},
			wantErr:  true,
			errMsg:   "node does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateNodeList(tt.nodeList, availableNodes)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPartitionBaseManager_CalculatePartitionStats(t *testing.T) {
	manager := NewPartitionBaseManager("v0.0.43")

	partitions := []types.Partition{
		{
			Name:     "compute",
			State:    "UP",
			Priority: 100,
			Nodes:    []string{"compute-01", "compute-02"},
		},
		{
			Name:     "gpu",
			State:    "UP",
			Priority: 200,
			Nodes:    []string{"gpu-01"},
		},
		{
			Name:     "debug",
			State:    "DOWN",
			Priority: 50,
			Nodes:    []string{"debug-01"},
		},
	}

	stats := manager.CalculatePartitionStats(partitions)

	assert.Equal(t, 3, stats.TotalPartitions)
	assert.Equal(t, 2, stats.UpPartitions)
	assert.Equal(t, 1, stats.DownPartitions)
	assert.Equal(t, 4, stats.TotalNodes) // 2 + 1 + 1
	assert.Equal(t, 116.67, stats.AveragePriority) // (100+200+50)/3 ≈ 116.67
}

func TestPartitionBaseManager_GetDefaultPartition(t *testing.T) {
	manager := NewPartitionBaseManager("v0.0.43")

	partitions := []types.Partition{
		{
			Name:  "compute",
			State: "UP",
			Flags: []string{},
		},
		{
			Name:  "gpu",
			State: "UP",
			Flags: []string{"Default"},
		},
		{
			Name:  "debug",
			State: "UP",
			Flags: []string{"Hidden"},
		},
	}

	defaultPartition := manager.GetDefaultPartition(partitions)
	require.NotNil(t, defaultPartition)
	assert.Equal(t, "gpu", defaultPartition.Name)
	assert.Contains(t, defaultPartition.Flags, "Default")
}

func TestPartitionBaseManager_SortPartitionsByPriority(t *testing.T) {
	manager := NewPartitionBaseManager("v0.0.43")

	partitions := []types.Partition{
		{Name: "debug", Priority: 50},
		{Name: "gpu", Priority: 200},
		{Name: "compute", Priority: 100},
	}

	sorted := manager.SortPartitionsByPriority(partitions)

	// Should be sorted by priority descending
	assert.Equal(t, "gpu", sorted[0].Name)     // Priority 200
	assert.Equal(t, "compute", sorted[1].Name) // Priority 100
	assert.Equal(t, "debug", sorted[2].Name)   // Priority 50
}