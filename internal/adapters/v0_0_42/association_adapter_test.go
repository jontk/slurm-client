package v0_0_42

import (
	"context"
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssociationAdapter_ValidateAssociationCreate(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Association"),
	}

	tests := []struct {
		name        string
		association *types.AssociationCreate
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "nil association",
			association: nil,
			wantErr:     true,
			errMsg:      "association data is required",
		},
		{
			name: "empty account",
			association: &types.AssociationCreate{
				Account: "",
				User:    "testuser",
				Cluster: "testcluster",
			},
			wantErr: true,
			errMsg:  "account is required",
		},
		{
			name: "empty cluster",
			association: &types.AssociationCreate{
				Account: "testacct",
				User:    "testuser",
				Cluster: "",
			},
			wantErr: true,
			errMsg:  "cluster is required",
		},
		{
			name: "valid user association",
			association: &types.AssociationCreate{
				Account: "testacct",
				User:    "testuser",
				Cluster: "testcluster",
			},
			wantErr: false,
		},
		{
			name: "valid account association",
			association: &types.AssociationCreate{
				Account:       "testacct",
				Cluster:       "testcluster",
				ParentAccount: "parent",
				DefaultQoS:    "normal",
			},
			wantErr: false,
		},
		{
			name: "association with limits",
			association: &types.AssociationCreate{
				Account:     "testacct",
				User:        "testuser",
				Cluster:     "testcluster",
				MaxJobs:     100,
				MaxCPUs:     1000,
				MaxNodes:    50,
				MaxWallDuration: "24:00:00",
			},
			wantErr: false,
		},
		{
			name: "invalid negative limits",
			association: &types.AssociationCreate{
				Account: "testacct",
				User:    "testuser",
				Cluster: "testcluster",
				MaxJobs: -1,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateAssociationCreate(tt.association)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAssociationAdapter_ApplyAssociationDefaults(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Association"),
	}

	tests := []struct {
		name     string
		input    *types.AssociationCreate
		expected *types.AssociationCreate
	}{
		{
			name: "apply defaults to minimal association",
			input: &types.AssociationCreate{
				Account: "testacct",
				User:    "testuser",
				Cluster: "testcluster",
			},
			expected: &types.AssociationCreate{
				Account:         "testacct",
				User:            "testuser",
				Cluster:         "testcluster",
				Partition:       "",
				DefaultQoS:      "",
				QoSList:         []string{},
				MaxJobs:         0,
				MaxCPUs:         0,
				MaxNodes:        0,
				MaxWallDuration: "",
				Priority:        0,
				GrpCPUs:         0,
				GrpJobs:         0,
				GrpNodes:        0,
			},
		},
		{
			name: "preserve existing values",
			input: &types.AssociationCreate{
				Account:         "testacct",
				User:            "testuser",
				Cluster:         "testcluster",
				Partition:       "compute",
				DefaultQoS:      "normal",
				QoSList:         []string{"normal", "high"},
				MaxJobs:         100,
				MaxCPUs:         1000,
				Priority:        500,
			},
			expected: &types.AssociationCreate{
				Account:         "testacct",
				User:            "testuser",
				Cluster:         "testcluster",
				Partition:       "compute",
				DefaultQoS:      "normal",
				QoSList:         []string{"normal", "high"},
				MaxJobs:         100,
				MaxCPUs:         1000,
				MaxNodes:        0,
				MaxWallDuration: "",
				Priority:        500,
				GrpCPUs:         0,
				GrpJobs:         0,
				GrpNodes:        0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.ApplyAssociationDefaults(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAssociationAdapter_FilterAssociationList(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Association"),
	}

	associations := []types.Association{
		{
			ID:        1,
			Account:   "physics",
			User:      "user1",
			Cluster:   "cluster1",
			Partition: "compute",
			DefaultQoS: "normal",
			MaxJobs:   100,
		},
		{
			ID:        2,
			Account:   "chemistry",
			User:      "user2",
			Cluster:   "cluster1",
			Partition: "gpu",
			DefaultQoS: "high",
			MaxJobs:   50,
		},
		{
			ID:        3,
			Account:   "physics",
			User:      "user1",
			Cluster:   "cluster2",
			Partition: "compute",
			DefaultQoS: "normal",
			MaxJobs:   200,
		},
		{
			ID:        4,
			Account:   "biology",
			Cluster:   "cluster1",
			Partition: "long",
			DefaultQoS: "low",
			MaxJobs:   25,
		},
	}

	tests := []struct {
		name     string
		opts     *types.AssociationListOptions
		expected []int // expected association IDs
	}{
		{
			name:     "no filters",
			opts:     &types.AssociationListOptions{},
			expected: []int{1, 2, 3, 4},
		},
		{
			name: "filter by accounts",
			opts: &types.AssociationListOptions{
				Accounts: []string{"physics"},
			},
			expected: []int{1, 3},
		},
		{
			name: "filter by users",
			opts: &types.AssociationListOptions{
				Users: []string{"user1"},
			},
			expected: []int{1, 3},
		},
		{
			name: "filter by clusters",
			opts: &types.AssociationListOptions{
				Clusters: []string{"cluster1"},
			},
			expected: []int{1, 2, 4},
		},
		{
			name: "filter by partitions",
			opts: &types.AssociationListOptions{
				Partitions: []string{"compute"},
			},
			expected: []int{1, 3},
		},
		{
			name: "combined filters",
			opts: &types.AssociationListOptions{
				Accounts: []string{"physics"},
				Clusters: []string{"cluster1"},
			},
			expected: []int{1},
		},
		{
			name: "filter by multiple accounts",
			opts: &types.AssociationListOptions{
				Accounts: []string{"physics", "chemistry"},
			},
			expected: []int{1, 2, 3},
		},
		{
			name: "no matches",
			opts: &types.AssociationListOptions{
				Accounts: []string{"nonexistent"},
			},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.FilterAssociationList(associations, tt.opts)
			resultIDs := make([]int, len(result))
			for i, assoc := range result {
				resultIDs[i] = assoc.ID
			}
			assert.Equal(t, tt.expected, resultIDs)
		})
	}
}

func TestAssociationAdapter_ValidateAssociationHierarchy(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Association"),
	}

	existingAssociations := []types.Association{
		{
			ID:      1,
			Account: "root",
			Cluster: "cluster1",
		},
		{
			ID:            2,
			Account:       "physics",
			ParentAccount: "root",
			Cluster:       "cluster1",
		},
		{
			ID:            3,
			Account:       "theory",
			ParentAccount: "physics",
			Cluster:       "cluster1",
		},
	}

	tests := []struct {
		name          string
		account       string
		parentAccount string
		cluster       string
		wantErr       bool
		errMsg        string
	}{
		{
			name:          "valid hierarchy",
			account:       "quantum",
			parentAccount: "physics",
			cluster:       "cluster1",
			wantErr:       false,
		},
		{
			name:          "self as parent",
			account:       "physics",
			parentAccount: "physics",
			cluster:       "cluster1",
			wantErr:       true,
			errMsg:        "cannot be its own parent",
		},
		{
			name:          "circular dependency",
			account:       "physics",
			parentAccount: "theory",
			cluster:       "cluster1",
			wantErr:       true,
			errMsg:        "circular dependency",
		},
		{
			name:          "nonexistent parent",
			account:       "newacct",
			parentAccount: "nonexistent",
			cluster:       "cluster1",
			wantErr:       true,
			errMsg:        "parent account not found",
		},
		{
			name:          "root account (no parent)",
			account:       "newroot",
			parentAccount: "",
			cluster:       "cluster1",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateAssociationHierarchy(tt.account, tt.parentAccount, tt.cluster, existingAssociations)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAssociationAdapter_ValidateContext(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Association"),
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

func TestAssociationAdapter_CalculateEffectiveLimits(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Association"),
	}

	tests := []struct {
		name           string
		association    *types.Association
		parentLimits   *types.AssociationLimits
		expectedMaxJobs int
		expectedMaxCPUs int
	}{
		{
			name: "no parent limits",
			association: &types.Association{
				MaxJobs: 100,
				MaxCPUs: 1000,
			},
			parentLimits: nil,
			expectedMaxJobs: 100,
			expectedMaxCPUs: 1000,
		},
		{
			name: "inherit from parent",
			association: &types.Association{
				MaxJobs: 0, // inherit
				MaxCPUs: 0, // inherit
			},
			parentLimits: &types.AssociationLimits{
				MaxJobs: 50,
				MaxCPUs: 500,
			},
			expectedMaxJobs: 50,
			expectedMaxCPUs: 500,
		},
		{
			name: "constrained by parent",
			association: &types.Association{
				MaxJobs: 200, // would like 200
				MaxCPUs: 2000, // would like 2000
			},
			parentLimits: &types.AssociationLimits{
				MaxJobs: 100, // but parent only allows 100
				MaxCPUs: 1000, // but parent only allows 1000
			},
			expectedMaxJobs: 100, // min(200, 100)
			expectedMaxCPUs: 1000, // min(2000, 1000)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limits := adapter.CalculateEffectiveLimits(tt.association, tt.parentLimits)
			assert.Equal(t, tt.expectedMaxJobs, limits.MaxJobs)
			assert.Equal(t, tt.expectedMaxCPUs, limits.MaxCPUs)
		})
	}
}

func TestAssociationAdapter_GetAssociationPath(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Association"),
	}

	associations := []types.Association{
		{
			ID:      1,
			Account: "root",
			Cluster: "cluster1",
		},
		{
			ID:            2,
			Account:       "physics",
			ParentAccount: "root",
			Cluster:       "cluster1",
		},
		{
			ID:            3,
			Account:       "theory",
			ParentAccount: "physics",
			Cluster:       "cluster1",
		},
		{
			ID:            4,
			Account:       "quantum",
			ParentAccount: "theory",
			Cluster:       "cluster1",
		},
	}

	tests := []struct {
		name         string
		account      string
		cluster      string
		expectedPath []string
	}{
		{
			name:         "root account",
			account:      "root",
			cluster:      "cluster1",
			expectedPath: []string{"root"},
		},
		{
			name:         "direct child",
			account:      "physics",
			cluster:      "cluster1",
			expectedPath: []string{"root", "physics"},
		},
		{
			name:         "deep hierarchy",
			account:      "quantum",
			cluster:      "cluster1",
			expectedPath: []string{"root", "physics", "theory", "quantum"},
		},
		{
			name:         "nonexistent account",
			account:      "nonexistent",
			cluster:      "cluster1",
			expectedPath: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := adapter.GetAssociationPath(tt.account, tt.cluster, associations)
			assert.Equal(t, tt.expectedPath, path)
		})
	}
}