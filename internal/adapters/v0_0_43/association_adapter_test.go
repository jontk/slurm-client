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

func TestAssociationAdapter_ValidateAssociationCreate(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Association"),
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
				AccountName: "",
				Cluster: "main",
			},
			wantErr: true,
			errMsg:  "account is required",
		},
		{
			name: "empty cluster",
			association: &types.AssociationCreate{
				AccountName: "physics",
				Cluster: "",
			},
			wantErr: true,
			errMsg:  "cluster is required",
		},
		{
			name: "negative max jobs",
			association: &types.AssociationCreate{
				AccountName: "physics",
				Cluster: "main",
				MaxJobs: -10,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "negative max CPUs",
			association: &types.AssociationCreate{
				AccountName: "physics",
				Cluster: "main",
				MaxCPUs: -100,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "negative max nodes",
			association: &types.AssociationCreate{
				AccountName:  "physics",
				Cluster:  "main",
				MaxNodes: -5,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "negative max wall time",
			association: &types.AssociationCreate{
				AccountName:  "physics",
				Cluster:      "main",
				MaxWallTime:  -3600,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "negative shares raw",
			association: &types.AssociationCreate{
				AccountName: "physics",
				Cluster:     "main",
				SharesRaw:   -100,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "negative priority",
			association: &types.AssociationCreate{
				AccountName:  "physics",
				Cluster:  "main",
				Priority: -50,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "valid basic association",
			association: &types.AssociationCreate{
				AccountName: "physics",
				Cluster: "main",
			},
			wantErr: false,
		},
		{
			name: "valid association with user",
			association: &types.AssociationCreate{
				AccountName: "physics",
				Cluster: "main",
				User:    "researcher1",
			},
			wantErr: false,
		},
		{
			name: "valid complex association",
			association: &types.AssociationCreate{
				AccountName:         "physics",
				Cluster:         "main",
				User:            "researcher1",
				Partition:       "compute",
				ParentAccountName:   "science",
				DefaultQoS:      "normal",
				QoSList:         []string{"normal", "high", "debug"},
				MaxJobs:         100,
				MaxCPUs:         1000,
				MaxNodes:        50,
				MaxWallDuration: 86400, // 24 hours
				FairShare:       1000,
				Priority:        500,
				GrpCPUs:         2000,
				GrpJobs:         200,
				GrpNodes:        100,
				GrpWall:         172800, // 48 hours
			},
			wantErr: false,
		},
		{
			name: "invalid QoS in list",
			association: &types.AssociationCreate{
				AccountName: "physics",
				Cluster: "main",
				QoSList: []string{"normal", "", "high"},
			},
			wantErr: true,
			errMsg:  "QoS name cannot be empty",
		},
		{
			name: "duplicate QoS in list",
			association: &types.AssociationCreate{
				AccountName: "physics",
				Cluster: "main",
				QoSList: []string{"normal", "high", "normal"},
			},
			wantErr: true,
			errMsg:  "duplicate QoS",
		},
		{
			name: "default QoS not in QoS list",
			association: &types.AssociationCreate{
				AccountName:    "physics",
				Cluster:    "main",
				DefaultQoS: "premium",
				QoSList:    []string{"normal", "high"},
			},
			wantErr: true,
			errMsg:  "default QoS must be in QoS list",
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
		AssociationBaseManager: base.NewAssociationBaseManager("v0.0.43"),
	}

	tests := []struct {
		name     string
		input    *types.AssociationCreate
		expected *types.AssociationCreate
	}{
		{
			name: "apply defaults to minimal association",
			input: &types.AssociationCreate{
				AccountName: "physics",
				Cluster: "main",
			},
			expected: &types.AssociationCreate{
				AccountName:         "physics",
				Cluster:         "main",
				User:            "",                       // Empty user (account-level)
				Partition:       "",                       // No partition restriction
				ParentAccountName:   "",                       // No parent account
				DefaultQoS:      "",                       // No default QoS
				QoSList:         []string{},               // Empty QoS list
				MaxJobs:         0,                        // Unlimited
				MaxCPUs:         0,                        // Unlimited
				MaxNodes:        0,                        // Unlimited
				MaxWallDuration: 0,                        // Unlimited
				FairShare:       1,                        // Default fair share
				Priority:        0,                        // Default priority
				GrpCPUs:         0,                        // Unlimited group CPUs
				GrpJobs:         0,                        // Unlimited group jobs
				GrpNodes:        0,                        // Unlimited group nodes
				GrpWall:         0,                        // Unlimited group wall time
			},
		},
		{
			name: "preserve existing values",
			input: &types.AssociationCreate{
				AccountName:         "physics",
				Cluster:         "main",
				User:            "researcher1",
				Partition:       "compute",
				ParentAccountName:   "science",
				DefaultQoS:      "normal",
				QoSList:         []string{"normal", "high"},
				MaxJobs:         50,
				MaxCPUs:         500,
				MaxNodes:        25,
				MaxWallDuration: 43200,
				FairShare:       2000,
				Priority:        1000,
				GrpCPUs:         1000,
				GrpJobs:         100,
				GrpNodes:        50,
				GrpWall:         86400,
			},
			expected: &types.AssociationCreate{
				AccountName:         "physics",
				Cluster:         "main",
				User:            "researcher1",
				Partition:       "compute",
				ParentAccountName:   "science",
				DefaultQoS:      "normal",
				QoSList:         []string{"normal", "high"},
				MaxJobs:         50,
				MaxCPUs:         500,
				MaxNodes:        25,
				MaxWallDuration: 43200,
				FairShare:       2000,
				Priority:        1000,
				GrpCPUs:         1000,
				GrpJobs:         100,
				GrpNodes:        50,
				GrpWall:         86400,
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
		AssociationBaseManager: base.NewAssociationBaseManager("v0.0.43"),
	}

	associations := []types.Association{
		{
			AccountName:    "physics",
			Cluster:    "main",
			User:       "alice",
			Partition:  "compute",
			DefaultQoS: "normal",
			QoSList:    []string{"normal", "high"},
			MaxJobs:    100,
			MaxCPUs:    1000,
			FairShare:  1000,
		},
		{
			AccountName:    "chemistry",
			Cluster:    "main",
			User:       "bob",
			Partition:  "gpu",
			DefaultQoS: "high",
			QoSList:    []string{"high", "premium"},
			MaxJobs:    50,
			MaxCPUs:    500,
			FairShare:  2000,
		},
		{
			AccountName:    "physics",
			Cluster:    "backup",
			User:       "charlie",
			Partition:  "compute",
			DefaultQoS: "normal",
			QoSList:    []string{"normal"},
			MaxJobs:    25,
			MaxCPUs:    250,
			FairShare:  500,
		},
		{
			AccountName:   "admin",
			Cluster:   "main",
			User:      "",  // Account-level association
			Partition: "",
			MaxJobs:   0,   // Unlimited
			MaxCPUs:   0,   // Unlimited
			FairShare: 100,
		},
	}

	tests := []struct {
		name     string
		opts     *types.AssociationListOptions
		expected []string // expected association identifiers (account-cluster-user)
	}{
		{
			name:     "no filters",
			opts:     &types.AssociationListOptions{},
			expected: []string{"physics-main-alice", "chemistry-main-bob", "physics-backup-charlie", "admin-main-"},
		},
		{
			name: "filter by accounts",
			opts: &types.AssociationListOptions{
				Accounts: []string{"physics"},
			},
			expected: []string{"physics-main-alice", "physics-backup-charlie"},
		},
		{
			name: "filter by clusters",
			opts: &types.AssociationListOptions{
				Clusters: []string{"main"},
			},
			expected: []string{"physics-main-alice", "chemistry-main-bob", "admin-main-"},
		},
		{
			name: "filter by users",
			opts: &types.AssociationListOptions{
				Users: []string{"bob", "charlie"},
			},
			expected: []string{"chemistry-main-bob", "physics-backup-charlie"},
		},
		{
			name: "filter by partitions",
			opts: &types.AssociationListOptions{
				Partitions: []string{"compute"},
			},
			expected: []string{"physics-main-alice", "physics-backup-charlie"},
		},
		{
			name: "filter by default QoS",
			opts: &types.AssociationListOptions{
				DefaultQoSList: []string{"normal"},
			},
			expected: []string{"physics-main-alice", "physics-backup-charlie"},
		},
		{
			name: "filter by QoS list",
			opts: &types.AssociationListOptions{
				QoSList: []string{"premium"},
			},
			expected: []string{"chemistry-main-bob"},
		},
		{
			name: "filter account-level associations (empty user)",
			opts: &types.AssociationListOptions{
				OnlyAccounts: true,
			},
			expected: []string{"admin-main-"},
		},
		{
			name: "filter user-level associations (non-empty user)",
			opts: &types.AssociationListOptions{
				OnlyUsers: true,
			},
			expected: []string{"physics-main-alice", "chemistry-main-bob", "physics-backup-charlie"},
		},
		{
			name: "combined filters",
			opts: &types.AssociationListOptions{
				Accounts:  []string{"physics"},
				Clusters:  []string{"main"},
				OnlyUsers: true,
			},
			expected: []string{"physics-main-alice"},
		},
		{
			name: "no matches",
			opts: &types.AssociationListOptions{
				Accounts: []string{"nonexistent"},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.FilterAssociationList(associations, tt.opts)
			resultIdentifiers := make([]string, len(result))
			for i, assoc := range result {
				resultIdentifiers[i] = assoc.Account + "-" + assoc.Cluster + "-" + assoc.User
			}
			assert.Equal(t, tt.expected, resultIdentifiers)
		})
	}
}

func TestAssociationAdapter_ValidateResourceLimits(t *testing.T) {
	adapter := &AssociationAdapter{
		AssociationBaseManager: base.NewAssociationBaseManager("v0.0.43"),
	}

	tests := []struct {
		name            string
		maxJobs         int
		maxCPUs         int
		maxNodes        int
		maxWallDuration int
		wantErr         bool
		errMsg          string
	}{
		{
			name:            "valid resource limits",
			maxJobs:         100,
			maxCPUs:         1000,
			maxNodes:        50,
			maxWallDuration: 86400,
			wantErr:         false,
		},
		{
			name:            "zero values (unlimited)",
			maxJobs:         0,
			maxCPUs:         0,
			maxNodes:        0,
			maxWallDuration: 0,
			wantErr:         false,
		},
		{
			name:            "negative max jobs",
			maxJobs:         -10,
			maxCPUs:         1000,
			maxNodes:        50,
			maxWallDuration: 86400,
			wantErr:         true,
			errMsg:          "must be non-negative",
		},
		{
			name:            "negative max CPUs",
			maxJobs:         100,
			maxCPUs:         -1000,
			maxNodes:        50,
			maxWallDuration: 86400,
			wantErr:         true,
			errMsg:          "must be non-negative",
		},
		{
			name:            "negative max nodes",
			maxJobs:         100,
			maxCPUs:         1000,
			maxNodes:        -50,
			maxWallDuration: 86400,
			wantErr:         true,
			errMsg:          "must be non-negative",
		},
		{
			name:            "negative max wall duration",
			maxJobs:         100,
			maxCPUs:         1000,
			maxNodes:        50,
			maxWallDuration: -86400,
			wantErr:         true,
			errMsg:          "must be non-negative",
		},
		{
			name:            "large valid values",
			maxJobs:         10000,
			maxCPUs:         100000,
			maxNodes:        1000,
			maxWallDuration: 2592000, // 30 days
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateResourceLimits(tt.maxJobs, tt.maxCPUs, tt.maxNodes, tt.maxWallDuration)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAssociationAdapter_ValidateQoSSettings(t *testing.T) {
	adapter := &AssociationAdapter{
		AssociationBaseManager: base.NewAssociationBaseManager("v0.0.43"),
	}

	tests := []struct {
		name       string
		defaultQoS string
		qosList    []string
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "valid QoS settings",
			defaultQoS: "normal",
			qosList:    []string{"normal", "high", "debug"},
			wantErr:    false,
		},
		{
			name:       "no QoS settings",
			defaultQoS: "",
			qosList:    []string{},
			wantErr:    false,
		},
		{
			name:       "default QoS only",
			defaultQoS: "normal",
			qosList:    []string{},
			wantErr:    true,
			errMsg:     "default QoS must be in QoS list",
		},
		{
			name:       "QoS list only",
			defaultQoS: "",
			qosList:    []string{"normal", "high"},
			wantErr:    false,
		},
		{
			name:       "default QoS not in list",
			defaultQoS: "premium",
			qosList:    []string{"normal", "high", "debug"},
			wantErr:    true,
			errMsg:     "default QoS must be in QoS list",
		},
		{
			name:       "empty QoS name in list",
			defaultQoS: "normal",
			qosList:    []string{"normal", "", "high"},
			wantErr:    true,
			errMsg:     "QoS name cannot be empty",
		},
		{
			name:       "duplicate QoS in list",
			defaultQoS: "normal",
			qosList:    []string{"normal", "high", "normal"},
			wantErr:    true,
			errMsg:     "duplicate QoS",
		},
		{
			name:       "QoS name with spaces",
			defaultQoS: "normal qos",
			qosList:    []string{"normal qos"},
			wantErr:    true,
			errMsg:     "QoS name cannot contain spaces",
		},
		{
			name:       "QoS name with invalid characters",
			defaultQoS: "normal@qos",
			qosList:    []string{"normal@qos"},
			wantErr:    true,
			errMsg:     "QoS name contains invalid characters",
		},
		{
			name:       "valid QoS names with allowed characters",
			defaultQoS: "normal-qos",
			qosList:    []string{"normal-qos", "high_priority", "debug.mode"},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateQoSSettings(tt.defaultQoS, tt.qosList)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAssociationAdapter_ValidateFairShareAndPriority(t *testing.T) {
	adapter := &AssociationAdapter{
		AssociationBaseManager: base.NewAssociationBaseManager("v0.0.43"),
	}

	tests := []struct {
		name      string
		fairShare int
		priority  int
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid fair share and priority",
			fairShare: 1000,
			priority:  500,
			wantErr:   false,
		},
		{
			name:      "zero values",
			fairShare: 0,
			priority:  0,
			wantErr:   false,
		},
		{
			name:      "negative fair share",
			fairShare: -100,
			priority:  500,
			wantErr:   true,
			errMsg:    "must be non-negative",
		},
		{
			name:      "negative priority",
			fairShare: 1000,
			priority:  -500,
			wantErr:   true,
			errMsg:    "must be non-negative",
		},
		{
			name:      "large valid values",
			fairShare: 100000,
			priority:  50000,
			wantErr:   false,
		},
		{
			name:      "minimum valid values",
			fairShare: 1,
			priority:  1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateFairShareAndPriority(tt.fairShare, tt.priority)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAssociationAdapter_ValidateGroupLimits(t *testing.T) {
	adapter := &AssociationAdapter{
		AssociationBaseManager: base.NewAssociationBaseManager("v0.0.43"),
	}

	tests := []struct {
		name    string
		grpCPUs int
		grpJobs int
		grpNodes int
		grpWall int
		wantErr bool
		errMsg  string
	}{
		{
			name:     "valid group limits",
			grpCPUs:  2000,
			grpJobs:  200,
			grpNodes: 100,
			grpWall:  172800,
			wantErr:  false,
		},
		{
			name:     "zero values (unlimited)",
			grpCPUs:  0,
			grpJobs:  0,
			grpNodes: 0,
			grpWall:  0,
			wantErr:  false,
		},
		{
			name:     "negative group CPUs",
			grpCPUs:  -2000,
			grpJobs:  200,
			grpNodes: 100,
			grpWall:  172800,
			wantErr:  true,
			errMsg:   "must be non-negative",
		},
		{
			name:     "negative group jobs",
			grpCPUs:  2000,
			grpJobs:  -200,
			grpNodes: 100,
			grpWall:  172800,
			wantErr:  true,
			errMsg:   "must be non-negative",
		},
		{
			name:     "negative group nodes",
			grpCPUs:  2000,
			grpJobs:  200,
			grpNodes: -100,
			grpWall:  172800,
			wantErr:  true,
			errMsg:   "must be non-negative",
		},
		{
			name:     "negative group wall time",
			grpCPUs:  2000,
			grpJobs:  200,
			grpNodes: 100,
			grpWall:  -172800,
			wantErr:  true,
			errMsg:   "must be non-negative",
		},
		{
			name:     "large valid values",
			grpCPUs:  200000,
			grpJobs:  20000,
			grpNodes: 10000,
			grpWall:  31536000, // 1 year
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateGroupLimits(tt.grpCPUs, tt.grpJobs, tt.grpNodes, tt.grpWall)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAssociationAdapter_ValidateHierarchy(t *testing.T) {
	adapter := &AssociationAdapter{
		AssociationBaseManager: base.NewAssociationBaseManager("v0.0.43"),
	}

	// Mock existing associations for hierarchy validation
	existingAssociations := []types.Association{
		{AccountName: "root", Cluster: "main", User: "", ParentAccountName: ""},
		{AccountName: "science", Cluster: "main", User: "", ParentAccountName: "root"},
		{AccountName: "physics", Cluster: "main", User: "", ParentAccountName: "science"},
		{AccountName: "chemistry", Cluster: "main", User: "", ParentAccountName: "science"},
	}

	tests := []struct {
		name          string
		account       string
		parentAccount string
		wantErr       bool
		errMsg        string
	}{
		{
			name:          "valid hierarchy",
			account:       "biophysics",
			parentAccountName: "physics",
			wantErr:       false,
		},
		{
			name:          "root account (no parent)",
			account:       "newroot",
			parentAccountName: "",
			wantErr:       false,
		},
		{
			name:          "self as parent",
			account:       "physics",
			parentAccountName: "physics",
			wantErr:       true,
			errMsg:        "account cannot be its own parent",
		},
		{
			name:          "circular reference",
			account:       "science",
			parentAccountName: "physics", // physics is already a child of science
			wantErr:       true,
			errMsg:        "would create circular reference",
		},
		{
			name:          "parent account does not exist",
			account:       "newaccount",
			parentAccountName: "nonexistent",
			wantErr:       true,
			errMsg:        "parent account does not exist",
		},
		{
			name:          "valid existing parent",
			account:       "theoretical-physics",
			parentAccountName: "physics",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateHierarchy(tt.account, tt.parentAccount, existingAssociations)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAssociationAdapter_ValidateUserAccountCombination(t *testing.T) {
	adapter := &AssociationAdapter{
		AssociationBaseManager: base.NewAssociationBaseManager("v0.0.43"),
	}

	tests := []struct {
		name    string
		user    string
		account string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid user-account combination",
			user:    "researcher1",
			account: "physics",
			wantErr: false,
		},
		{
			name:    "account-level association (empty user)",
			user:    "",
			account: "physics",
			wantErr: false,
		},
		{
			name:    "user with empty account",
			user:    "researcher1",
			account: "",
			wantErr: true,
			errMsg:  "account is required",
		},
		{
			name:    "user name with spaces",
			user:    "researcher name",
			account: "physics",
			wantErr: true,
			errMsg:  "user name cannot contain spaces",
		},
		{
			name:    "account name with spaces",
			user:    "researcher1",
			account: "physics account",
			wantErr: true,
			errMsg:  "account name cannot contain spaces",
		},
		{
			name:    "user name with invalid characters",
			user:    "researcher@domain",
			account: "physics",
			wantErr: true,
			errMsg:  "user name contains invalid characters",
		},
		{
			name:    "account name with invalid characters",
			user:    "researcher1",
			account: "physics@domain",
			wantErr: true,
			errMsg:  "account name contains invalid characters",
		},
		{
			name:    "valid names with allowed characters",
			user:    "researcher-1",
			account: "physics_dept",
			wantErr: false,
		},
		{
			name:    "user name too long",
			user:    "verylongusernamethatexceedsthelimitofcharactersallowed",
			account: "physics",
			wantErr: true,
			errMsg:  "user name too long",
		},
		{
			name:    "account name too long",
			user:    "researcher1",
			account: "verylongaccountnamethatexceedsthelimitofcharactersallowed",
			wantErr: true,
			errMsg:  "account name too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateUserAccountCombination(tt.user, tt.account)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
