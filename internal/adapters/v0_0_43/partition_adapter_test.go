// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPartitionAdapter_ValidatePartitionCreate(t *testing.T) {
	adapter := &PartitionAdapter{
		PartitionBaseManager: base.NewPartitionBaseManager("v0.0.43"),
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
			name: "negative MaxTime",
			partition: &types.PartitionCreate{
				Name:    "test-partition",
				MaxTime: -1,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "negative DefaultTime",
			partition: &types.PartitionCreate{
				Name:        "test-partition",
				DefaultTime: -30,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "negative MaxCPUsPerNode",
			partition: &types.PartitionCreate{
				Name:           "test-partition",
				MaxCPUsPerNode: -4,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "valid basic partition",
			partition: &types.PartitionCreate{
				Name:  "test-partition",
				State: "UP",
			},
			wantErr: false,
		},
		{
			name: "valid complex partition",
			partition: &types.PartitionCreate{
				Name:             "compute-partition",
				State:            "UP",
				MaxTime:          1440, // 24 hours
				DefaultTime:      60,   // 1 hour
				MaxNodes:         100,
				MaxCPUsPerNode:   64,
				AllowAccounts:    []string{"physics", "chemistry"},
				DenyAccounts:     []string{"guest"},
				AllowGroups:      []string{"researchers"},
				DefaultQoS:       "normal",
				QoSList:          []string{"normal", "high"},
				PreemptMode:      "REQUEUE",
				PriorityJobFactor: 2,
				PriorityTier:     1000,
				GraceTime:        300,
				Features:         []string{"gpu", "infiniband"},
				Nodes:            []string{"node[001-100]"},
			},
			wantErr: false,
		},
		{
			name: "invalid state",
			partition: &types.PartitionCreate{
				Name:  "test-partition",
				State: "INVALID_STATE",
			},
			wantErr: true,
			errMsg:  "invalid partition state",
		},
		{
			name: "negative priority values",
			partition: &types.PartitionCreate{
				Name:              "test-partition",
				PriorityJobFactor: -1,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
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
		PartitionBaseManager: base.NewPartitionBaseManager("v0.0.43"),
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
				Name:             "test-partition",
				State:            "UP",
				MaxTime:          0,      // Unlimited
				DefaultTime:      60,     // 1 hour default
				MaxNodes:         0,      // Unlimited
				MaxCPUsPerNode:   0,      // Unlimited
				AllowAccounts:    []string{},
				DenyAccounts:     []string{},
				AllowGroups:      []string{},
				QoSList:          []string{},
				PreemptMode:      "OFF",
				PriorityJobFactor: 1,
				PriorityTier:     0,
				GraceTime:        0,
				Features:         []string{},
				Nodes:            []string{},
			},
		},
		{
			name: "preserve existing values",
			input: &types.PartitionCreate{
				Name:             "compute-partition",
				State:            "DRAIN",
				MaxTime:          2880,
				DefaultTime:      120,
				MaxNodes:         50,
				MaxCPUsPerNode:   32,
				AllowAccounts:    []string{"physics"},
				DefaultQoS:       "high",
				PreemptMode:      "REQUEUE",
				PriorityJobFactor: 3,
				PriorityTier:     2000,
				GraceTime:        600,
			},
			expected: &types.PartitionCreate{
				Name:             "compute-partition",
				State:            "DRAIN",
				MaxTime:          2880,
				DefaultTime:      120,
				MaxNodes:         50,
				MaxCPUsPerNode:   32,
				AllowAccounts:    []string{"physics"},
				DenyAccounts:     []string{},
				AllowGroups:      []string{},
				DefaultQoS:       "high",
				QoSList:          []string{},
				PreemptMode:      "REQUEUE",
				PriorityJobFactor: 3,
				PriorityTier:     2000,
				GraceTime:        600,
				Features:         []string{},
				Nodes:            []string{},
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
		PartitionBaseManager: base.NewPartitionBaseManager("v0.0.43"),
	}

	partitions := []types.Partition{
		{
			Name:        "compute",
			State:       "UP",
			MaxTime:     1440,
			DefaultTime: 60,
			MaxNodes:    100,
			DefaultQoS:  "normal",
			Features:    []string{"cpu", "infiniband"},
		},
		{
			Name:        "gpu",
			State:       "UP", 
			MaxTime:     720,
			DefaultTime: 30,
			MaxNodes:    10,
			DefaultQoS:  "high",
			Features:    []string{"gpu", "cuda"},
		},
		{
			Name:        "debug",
			State:       "DRAIN",
			MaxTime:     30,
			DefaultTime: 10,
			MaxNodes:    2,
			DefaultQoS:  "debug",
			Features:    []string{"debug"},
		},
		{
			Name:        "bigmem",
			State:       "UP",
			MaxTime:     2880,
			DefaultTime: 120,
			MaxNodes:    5,
			DefaultQoS:  "normal",
			Features:    []string{"bigmem", "high-memory"},
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
			expected: []string{"compute", "gpu", "debug", "bigmem"},
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
			expected: []string{"compute", "gpu", "bigmem"},
		},
		{
			name: "filter by features",
			opts: &types.PartitionListOptions{
				Features: []string{"gpu"},
			},
			expected: []string{"gpu"},
		},
		{
			name: "filter by default QoS",
			opts: &types.PartitionListOptions{
				DefaultQoS: []string{"normal"},
			},
			expected: []string{"compute", "bigmem"},
		},
		{
			name: "combined filters",
			opts: &types.PartitionListOptions{
				States:   []string{"UP"},
				Features: []string{"infiniband", "high-memory"},
			},
			expected: []string{"compute", "bigmem"},
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

func TestPartitionAdapter_ValidatePartitionState(t *testing.T) {
	adapter := &PartitionAdapter{
		PartitionBaseManager: base.NewPartitionBaseManager("v0.0.43"),
	}

	tests := []struct {
		name    string
		state   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid UP state",
			state:   "UP",
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
			name:    "valid INACTIVE state",
			state:   "INACTIVE",
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
			errMsg:  "invalid partition state",
		},
		{
			name:    "lowercase state",
			state:   "up",
			wantErr: true,
			errMsg:  "invalid partition state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidatePartitionState(tt.state)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPartitionAdapter_ValidateTimeFields(t *testing.T) {
	adapter := &PartitionAdapter{
		PartitionBaseManager: base.NewPartitionBaseManager("v0.0.43"),
	}

	tests := []struct {
		name        string
		maxTime     int
		defaultTime int
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid time values",
			maxTime:     1440,
			defaultTime: 60,
			wantErr:     false,
		},
		{
			name:        "zero values (unlimited)",
			maxTime:     0,
			defaultTime: 0,
			wantErr:     false,
		},
		{
			name:        "default time exceeds max time",
			maxTime:     60,
			defaultTime: 120,
			wantErr:     true,
			errMsg:      "default time cannot exceed max time",
		},
		{
			name:        "negative max time",
			maxTime:     -1,
			defaultTime: 60,
			wantErr:     true,
			errMsg:      "must be non-negative",
		},
		{
			name:        "negative default time",
			maxTime:     1440,
			defaultTime: -30,
			wantErr:     true,
			errMsg:      "must be non-negative",
		},
		{
			name:        "unlimited max time with default",
			maxTime:     0,
			defaultTime: 60,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateTimeFields(tt.maxTime, tt.defaultTime)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPartitionAdapter_ValidateResourceLimits(t *testing.T) {
	adapter := &PartitionAdapter{
		PartitionBaseManager: base.NewPartitionBaseManager("v0.0.43"),
	}

	tests := []struct {
		name           string
		maxNodes       int
		maxCPUsPerNode int
		wantErr        bool
		errMsg         string
	}{
		{
			name:           "valid resource limits",
			maxNodes:       10,
			maxCPUsPerNode: 64,
			wantErr:        false,
		},
		{
			name:           "zero values (unlimited)",
			maxNodes:       0,
			maxCPUsPerNode: 0,
			wantErr:        false,
		},
		{
			name:           "negative max nodes",
			maxNodes:       -1,
			maxCPUsPerNode: 64,
			wantErr:        true,
			errMsg:         "must be non-negative",
		},
		{
			name:           "negative max CPUs per node",
			maxNodes:       10,
			maxCPUsPerNode: -4,
			wantErr:        true,
			errMsg:         "must be non-negative",
		},
		{
			name:           "large valid values",
			maxNodes:       1000,
			maxCPUsPerNode: 256,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateResourceLimits(tt.maxNodes, tt.maxCPUsPerNode)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPartitionAdapter_ValidateAccessControls(t *testing.T) {
	adapter := &PartitionAdapter{
		PartitionBaseManager: base.NewPartitionBaseManager("v0.0.43"),
	}

	tests := []struct {
		name          string
		allowAccounts []string
		denyAccounts  []string
		allowGroups   []string
		wantErr       bool
		errMsg        string
	}{
		{
			name:          "valid access controls",
			allowAccounts: []string{"physics", "chemistry"},
			denyAccounts:  []string{"guest"},
			allowGroups:   []string{"researchers"},
			wantErr:       false,
		},
		{
			name:          "empty access controls",
			allowAccounts: []string{},
			denyAccounts:  []string{},
			allowGroups:   []string{},
			wantErr:       false,
		},
		{
			name:          "overlapping allow and deny accounts",
			allowAccounts: []string{"physics", "chemistry"},
			denyAccounts:  []string{"physics", "guest"},
			allowGroups:   []string{"researchers"},
			wantErr:       true,
			errMsg:        "account cannot be both allowed and denied",
		},
		{
			name:          "duplicate allow accounts",
			allowAccounts: []string{"physics", "physics"},
			denyAccounts:  []string{"guest"},
			allowGroups:   []string{"researchers"},
			wantErr:       true,
			errMsg:        "duplicate account in allow list",
		},
		{
			name:          "duplicate deny accounts",
			allowAccounts: []string{"physics"},
			denyAccounts:  []string{"guest", "guest"},
			allowGroups:   []string{"researchers"},
			wantErr:       true,
			errMsg:        "duplicate account in deny list",
		},
		{
			name:          "duplicate allow groups",
			allowAccounts: []string{"physics"},
			denyAccounts:  []string{"guest"},
			allowGroups:   []string{"researchers", "researchers"},
			wantErr:       true,
			errMsg:        "duplicate group in allow list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateAccessControls(tt.allowAccounts, tt.denyAccounts, tt.allowGroups)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPartitionAdapter_ValidatePreemptMode(t *testing.T) {
	adapter := &PartitionAdapter{
		PartitionBaseManager: base.NewPartitionBaseManager("v0.0.43"),
	}

	tests := []struct {
		name        string
		preemptMode string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid OFF",
			preemptMode: "OFF",
			wantErr:     false,
		},
		{
			name:        "valid CANCEL",
			preemptMode: "CANCEL",
			wantErr:     false,
		},
		{
			name:        "valid REQUEUE",
			preemptMode: "REQUEUE",
			wantErr:     false,
		},
		{
			name:        "valid SUSPEND",
			preemptMode: "SUSPEND",
			wantErr:     false,
		},
		{
			name:        "empty mode (should use default)",
			preemptMode: "",
			wantErr:     false,
		},
		{
			name:        "invalid mode",
			preemptMode: "INVALID_MODE",
			wantErr:     true,
			errMsg:      "invalid preempt mode",
		},
		{
			name:        "lowercase mode",
			preemptMode: "cancel",
			wantErr:     true,
			errMsg:      "invalid preempt mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidatePreemptMode(tt.preemptMode)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
