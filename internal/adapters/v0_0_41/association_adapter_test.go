// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

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
		BaseManager: base.NewBaseManager("v0.0.41", "Association"),
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
				Account:         "testacct",
				User:            "testuser",
				Cluster:         "testcluster",
				MaxJobs:         100,
				MaxCPUs:         1000,
				MaxNodes:        50,
				MaxWallDuration: "24:00:00",
				GrpCPUs:         2000,
				GrpJobs:         200,
				GrpNodes:        100,
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
		{
			name: "invalid negative group limits",
			association: &types.AssociationCreate{
				Account: "testacct",
				User:    "testuser",
				Cluster: "testcluster",
				GrpCPUs: -1,
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
		BaseManager: base.NewBaseManager("v0.0.41", "Association"),
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
				FairshareShares: 0,
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
				GrpCPUs:         2000,
				FairshareShares: 1000,
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
				GrpCPUs:         2000,
				GrpJobs:         0,
				GrpNodes:        0,
				FairshareShares: 1000,
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
		BaseManager: base.NewBaseManager("v0.0.41", "Association"),
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
			Priority:  500,
		},
		{
			ID:        2,
			Account:   "chemistry",
			User:      "user2",
			Cluster:   "cluster1",
			Partition: "gpu",
			DefaultQoS: "high",
			MaxJobs:   50,
			Priority:  1000,
		},
		{
			ID:        3,
			Account:   "physics",
			User:      "user1",
			Cluster:   "cluster2",
			Partition: "compute",
			DefaultQoS: "normal",
			MaxJobs:   200,
			Priority:  500,
		},
		{
			ID:        4,
			Account:   "biology",
			Cluster:   "cluster1",
			Partition: "long",
			DefaultQoS: "low",
			MaxJobs:   25,
			Priority:  100,
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
			name: "filter by default QoS",
			opts: &types.AssociationListOptions{
				DefaultQoS: []string{"normal"},
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
		BaseManager: base.NewBaseManager("v0.0.41", "Association"),
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

func TestAssociationAdapter_CalculateFairshareUsage(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Association"),
	}

	tests := []struct {
		name                string
		association         *types.Association
		usage               int64
		expectedFairshare   float64
	}{
		{
			name: "normal usage",
			association: &types.Association{
				FairshareShares: 1000,
				RawShares:       1000,
			},
			usage:             5000,
			expectedFairshare: 0.5, // usage/shares ratio
		},
		{
			name: "low usage",
			association: &types.Association{
				FairshareShares: 2000,
				RawShares:       2000,
			},
			usage:             1000,
			expectedFairshare: 0.75, // lower usage = higher fairshare
		},
		{
			name: "high usage",
			association: &types.Association{
				FairshareShares: 500,
				RawShares:       500,
			},
			usage:             5000,
			expectedFairshare: 0.1, // high usage = lower fairshare
		},
		{
			name: "zero shares",
			association: &types.Association{
				FairshareShares: 0,
				RawShares:       0,
			},
			usage:             1000,
			expectedFairshare: 0.0, // no shares = no fairshare
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fairshare := adapter.CalculateFairshareUsage(tt.association, tt.usage)
			assert.InDelta(t, tt.expectedFairshare, fairshare, 0.1)
		})
	}
}

func TestAssociationAdapter_ValidateQoSLimits(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Association"),
	}

	tests := []struct {
		name    string
		qosList []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty QoS list",
			qosList: []string{},
			wantErr: false,
		},
		{
			name:    "valid QoS names",
			qosList: []string{"normal", "high", "low"},
			wantErr: false,
		},
		{
			name:    "invalid QoS name (empty)",
			qosList: []string{"normal", "", "high"},
			wantErr: true,
			errMsg:  "QoS name cannot be empty",
		},
		{
			name:    "duplicate QoS names",
			qosList: []string{"normal", "high", "normal"},
			wantErr: true,
			errMsg:  "duplicate QoS name",
		},
		{
			name:    "single valid QoS",
			qosList: []string{"normal"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateQoSLimits(tt.qosList)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAssociationAdapter_GetEffectivePriority(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Association"),
	}

	parentAssoc := &types.Association{
		ID:       1,
		Account:  "parent",
		Priority: 1000,
	}

	tests := []struct {
		name             string
		association      *types.Association
		parent           *types.Association
		expectedPriority int
	}{
		{
			name: "association with explicit priority",
			association: &types.Association{
				ID:       2,
				Account:  "child",
				Priority: 1500,
			},
			parent:           parentAssoc,
			expectedPriority: 1500,
		},
		{
			name: "association inherits from parent",
			association: &types.Association{
				ID:       3,
				Account:  "child2",
				Priority: 0, // inherit
			},
			parent:           parentAssoc,
			expectedPriority: 1000,
		},
		{
			name: "no parent association",
			association: &types.Association{
				ID:       4,
				Account:  "root",
				Priority: 500,
			},
			parent:           nil,
			expectedPriority: 500,
		},
		{
			name: "priority capped by parent",
			association: &types.Association{
				ID:       5,
				Account:  "child3",
				Priority: 2000, // wants higher
			},
			parent:           parentAssoc, // but parent only has 1000
			expectedPriority: 1000,        // capped at parent level
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			priority := adapter.GetEffectivePriority(tt.association, tt.parent)
			assert.Equal(t, tt.expectedPriority, priority)
		})
	}
}

func TestAssociationAdapter_ValidateContext(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Association"),
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
