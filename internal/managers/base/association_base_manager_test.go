package base

import (
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssociationBaseManager_New(t *testing.T) {
	manager := NewAssociationBaseManager("v0.0.43")
	assert.NotNil(t, manager)
	assert.Equal(t, "v0.0.43", manager.GetVersion())
	assert.Equal(t, "Association", manager.GetResourceType())
}

func TestAssociationBaseManager_ValidateAssociationCreate(t *testing.T) {
	manager := NewAssociationBaseManager("v0.0.43")

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
			},
			wantErr: true,
			errMsg:  "account is required",
		},
		{
			name: "valid user association",
			association: &types.AssociationCreate{
				Account: "test-account",
				User:    "testuser",
			},
			wantErr: false,
		},
		{
			name: "valid account association",
			association: &types.AssociationCreate{
				Account:       "child-account",
				ParentAccount: "parent-account",
			},
			wantErr: false,
		},
		{
			name: "valid complex association",
			association: &types.AssociationCreate{
				Account:         "test-account",
				User:            "testuser",
				Partition:       "compute",
				QoS:             "normal",
				DefaultQoS:      "normal",
				Priority:        intPtr(100),
				MaxJobs:         intPtr(10),
				MaxSubmitJobs:   intPtr(20),
				MaxWallDuration: intPtr(1440), // 24 hours
				Flags:           []string{"NoDefaultQOS"},
			},
			wantErr: false,
		},
		{
			name: "negative priority",
			association: &types.AssociationCreate{
				Account:  "test-account",
				User:     "testuser",
				Priority: intPtr(-1),
			},
			wantErr: true,
			errMsg:  "priority must be non-negative",
		},
		{
			name: "negative max jobs",
			association: &types.AssociationCreate{
				Account: "test-account",
				User:    "testuser",
				MaxJobs: intPtr(-1),
			},
			wantErr: true,
			errMsg:  "max jobs must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateAssociationCreate(tt.association)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAssociationBaseManager_ValidateAssociationUpdate(t *testing.T) {
	manager := NewAssociationBaseManager("v0.0.43")

	tests := []struct {
		name    string
		update  *types.AssociationUpdate
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil update",
			update:  nil,
			wantErr: true,
			errMsg:  "association update data is required",
		},
		{
			name: "empty account",
			update: &types.AssociationUpdate{
				Account: "",
			},
			wantErr: true,
			errMsg:  "account is required",
		},
		{
			name: "valid update",
			update: &types.AssociationUpdate{
				Account:  "test-account",
				User:     "testuser",
				Priority: intPtr(200),
			},
			wantErr: false,
		},
		{
			name: "negative priority",
			update: &types.AssociationUpdate{
				Account:  "test-account",
				User:     "testuser",
				Priority: intPtr(-1),
			},
			wantErr: true,
			errMsg:  "priority must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateAssociationUpdate(tt.update)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAssociationBaseManager_ApplyAssociationDefaults(t *testing.T) {
	manager := NewAssociationBaseManager("v0.0.43")

	tests := []struct {
		name     string
		input    *types.AssociationCreate
		expected *types.AssociationCreate
	}{
		{
			name: "apply defaults to minimal association",
			input: &types.AssociationCreate{
				Account: "test-account",
				User:    "testuser",
			},
			expected: &types.AssociationCreate{
				Account:   "test-account",
				User:      "testuser",
				Priority:  intPtr(1),         // Default priority
				Flags:     []string{},        // Empty flags
				QoS:       "",                // No default QoS
				DefaultQoS: "",               // No default QoS
			},
		},
		{
			name: "preserve existing values",
			input: &types.AssociationCreate{
				Account:         "test-account",
				User:            "testuser",
				Partition:       "gpu",
				QoS:             "high",
				DefaultQoS:      "normal",
				Priority:        intPtr(500),
				MaxJobs:         intPtr(50),
				MaxSubmitJobs:   intPtr(100),
				MaxWallDuration: intPtr(2880), // 48 hours
				Flags:           []string{"NoDefaultQOS", "RequiresReservation"},
			},
			expected: &types.AssociationCreate{
				Account:         "test-account",
				User:            "testuser",
				Partition:       "gpu",
				QoS:             "high",
				DefaultQoS:      "normal",
				Priority:        intPtr(500),
				MaxJobs:         intPtr(50),
				MaxSubmitJobs:   intPtr(100),
				MaxWallDuration: intPtr(2880),
				Flags:           []string{"NoDefaultQOS", "RequiresReservation"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.ApplyAssociationDefaults(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAssociationBaseManager_FilterAssociationList(t *testing.T) {
	manager := NewAssociationBaseManager("v0.0.43")

	associations := []types.Association{
		{
			Account:       "account1",
			User:          "user1",
			Partition:     "compute",
			QoS:           "normal",
			DefaultQoS:    "normal",
			Priority:      100,
			MaxJobs:       10,
			ParentAccount: "root",
			Flags:         []string{},
		},
		{
			Account:       "account2",
			User:          "user2",
			Partition:     "gpu",
			QoS:           "high",
			DefaultQoS:    "normal",
			Priority:      200,
			MaxJobs:       5,
			ParentAccount: "account1",
			Flags:         []string{"NoDefaultQOS"},
		},
		{
			Account:       "account1",
			User:          "user3",
			Partition:     "compute",
			QoS:           "normal",
			DefaultQoS:    "normal",
			Priority:      50,
			MaxJobs:       20,
			ParentAccount: "",
			Flags:         []string{"RequiresReservation"},
		},
	}

	tests := []struct {
		name     string
		opts     *types.AssociationListOptions
		expected int // expected count of associations
	}{
		{
			name:     "no filters",
			opts:     &types.AssociationListOptions{},
			expected: 3,
		},
		{
			name: "filter by account",
			opts: &types.AssociationListOptions{
				Accounts: []string{"account1"},
			},
			expected: 2, // user1 and user3 in account1
		},
		{
			name: "filter by user",
			opts: &types.AssociationListOptions{
				Users: []string{"user1", "user2"},
			},
			expected: 2,
		},
		{
			name: "filter by partition",
			opts: &types.AssociationListOptions{
				Partitions: []string{"gpu"},
			},
			expected: 1, // Only user2 in gpu partition
		},
		{
			name: "filter by QoS",
			opts: &types.AssociationListOptions{
				QoSList: []string{"high"},
			},
			expected: 1, // Only user2 has high QoS
		},
		{
			name: "filter by parent account",
			opts: &types.AssociationListOptions{
				ParentAccounts: []string{"account1"},
			},
			expected: 1, // Only account2 has account1 as parent
		},
		{
			name: "filter by minimum priority",
			opts: &types.AssociationListOptions{
				MinPriority: intPtr(100),
			},
			expected: 2, // user1 (100) and user2 (200)
		},
		{
			name: "filter by flag",
			opts: &types.AssociationListOptions{
				WithFlags: []string{"NoDefaultQOS"},
			},
			expected: 1, // Only user2 has NoDefaultQOS flag
		},
		{
			name: "combined filters",
			opts: &types.AssociationListOptions{
				Accounts:   []string{"account1"},
				Partitions: []string{"compute"},
			},
			expected: 2, // user1 and user3 in account1 compute partition
		},
		{
			name: "no matches",
			opts: &types.AssociationListOptions{
				Users: []string{"nonexistent"},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.FilterAssociationList(associations, tt.opts)
			assert.Len(t, result, tt.expected)
		})
	}
}

func TestAssociationBaseManager_ValidateAssociationHierarchy(t *testing.T) {
	manager := NewAssociationBaseManager("v0.0.43")

	// Mock existing associations for hierarchy validation
	existingAssociations := []types.Association{
		{Account: "root", ParentAccount: ""},
		{Account: "parent1", ParentAccount: "root"},
		{Account: "child1", ParentAccount: "parent1"},
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
			account:       "new-child",
			parentAccount: "parent1",
			wantErr:       false,
		},
		{
			name:          "self as parent",
			account:       "test-account",
			parentAccount: "test-account",
			wantErr:       true,
			errMsg:        "cannot be its own parent",
		},
		{
			name:          "circular dependency",
			account:       "parent1",
			parentAccount: "child1",
			wantErr:       true,
			errMsg:        "would create circular dependency",
		},
		{
			name:          "nonexistent parent",
			account:       "new-account",
			parentAccount: "nonexistent",
			wantErr:       true,
			errMsg:        "parent account does not exist",
		},
		{
			name:          "empty parent (root account)",
			account:       "new-root",
			parentAccount: "",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateAssociationHierarchy(tt.account, tt.parentAccount, existingAssociations)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAssociationBaseManager_GetAssociationLimits(t *testing.T) {
	manager := NewAssociationBaseManager("v0.0.43")

	association := types.Association{
		Account:         "test-account",
		User:            "testuser",
		MaxJobs:         10,
		MaxSubmitJobs:   20,
		MaxWallDuration: 1440, // 24 hours
		MaxCPUs:         intPtr(64),
		MaxNodes:        intPtr(4),
	}

	limits := manager.GetAssociationLimits(association)

	assert.Equal(t, 10, limits.MaxJobs)
	assert.Equal(t, 20, limits.MaxSubmitJobs)
	assert.Equal(t, 1440, limits.MaxWallDuration)
	assert.Equal(t, 64, *limits.MaxCPUs)
	assert.Equal(t, 4, *limits.MaxNodes)
}

func TestAssociationBaseManager_ValidateAssociationLimits(t *testing.T) {
	manager := NewAssociationBaseManager("v0.0.43")

	tests := []struct {
		name        string
		association *types.AssociationCreate
		wantErr     bool
		errMsg      string
	}{
		{
			name: "valid limits",
			association: &types.AssociationCreate{
				Account:         "test-account",
				User:            "testuser",
				MaxJobs:         intPtr(10),
				MaxSubmitJobs:   intPtr(20),
				MaxWallDuration: intPtr(1440),
				MaxCPUs:         intPtr(64),
				MaxNodes:        intPtr(4),
			},
			wantErr: false,
		},
		{
			name: "max submit jobs less than max jobs",
			association: &types.AssociationCreate{
				Account:       "test-account",
				User:          "testuser",
				MaxJobs:       intPtr(20),
				MaxSubmitJobs: intPtr(10),
			},
			wantErr: true,
			errMsg:  "max submit jobs must be greater than or equal to max jobs",
		},
		{
			name: "negative max wall duration",
			association: &types.AssociationCreate{
				Account:         "test-account",
				User:            "testuser",
				MaxWallDuration: intPtr(-1),
			},
			wantErr: true,
			errMsg:  "max wall duration must be non-negative",
		},
		{
			name: "negative max CPUs",
			association: &types.AssociationCreate{
				Account: "test-account",
				User:    "testuser",
				MaxCPUs: intPtr(-1),
			},
			wantErr: true,
			errMsg:  "max CPUs must be non-negative",
		},
		{
			name: "negative max nodes",
			association: &types.AssociationCreate{
				Account:  "test-account",
				User:     "testuser",
				MaxNodes: intPtr(-1),
			},
			wantErr: true,
			errMsg:  "max nodes must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateAssociationLimits(tt.association)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAssociationBaseManager_CalculateEffectivePriority(t *testing.T) {
	manager := NewAssociationBaseManager("v0.0.43")

	// Mock account hierarchy with priorities
	accountHierarchy := map[string]types.Association{
		"root": {
			Account:  "root",
			Priority: 1000,
		},
		"parent1": {
			Account:       "parent1",
			ParentAccount: "root",
			Priority:      800,
		},
		"child1": {
			Account:       "child1",
			ParentAccount: "parent1",
			Priority:      600,
		},
	}

	tests := []struct {
		name        string
		association types.Association
		expected    int
	}{
		{
			name: "root account priority",
			association: types.Association{
				Account:  "root",
				Priority: 1000,
			},
			expected: 1000,
		},
		{
			name: "inherited priority from parent",
			association: types.Association{
				Account:       "child1",
				ParentAccount: "parent1",
				Priority:      600,
			},
			expected: 600, // Uses own priority, not inherited
		},
		{
			name: "zero priority inherits from parent",
			association: types.Association{
				Account:       "child2",
				ParentAccount: "parent1",
				Priority:      0,
			},
			expected: 800, // Inherits parent priority
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			priority := manager.CalculateEffectivePriority(tt.association, accountHierarchy)
			assert.Equal(t, tt.expected, priority)
		})
	}
}

func TestAssociationBaseManager_GetAssociationPath(t *testing.T) {
	manager := NewAssociationBaseManager("v0.0.43")

	associations := []types.Association{
		{Account: "root", ParentAccount: ""},
		{Account: "parent1", ParentAccount: "root"},
		{Account: "child1", ParentAccount: "parent1"},
		{Account: "grandchild1", ParentAccount: "child1"},
	}

	tests := []struct {
		name        string
		accountName string
		expected    []string
	}{
		{
			name:        "root account",
			accountName: "root",
			expected:    []string{"root"},
		},
		{
			name:        "parent account",
			accountName: "parent1",
			expected:    []string{"root", "parent1"},
		},
		{
			name:        "child account",
			accountName: "child1",
			expected:    []string{"root", "parent1", "child1"},
		},
		{
			name:        "grandchild account",
			accountName: "grandchild1",
			expected:    []string{"root", "parent1", "child1", "grandchild1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := manager.GetAssociationPath(tt.accountName, associations)
			assert.Equal(t, tt.expected, path)
		})
	}
}