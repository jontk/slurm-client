package base

import (
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountBaseManager_New(t *testing.T) {
	manager := NewAccountBaseManager("v0.0.43")
	assert.NotNil(t, manager)
	assert.Equal(t, "v0.0.43", manager.GetVersion())
	assert.Equal(t, "Account", manager.GetResourceType())
}

func TestAccountBaseManager_ValidateAccountCreate(t *testing.T) {
	manager := NewAccountBaseManager("v0.0.43")

	tests := []struct {
		name    string
		account *types.AccountCreate
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil account",
			account: nil,
			wantErr: true,
			errMsg:  "account data is required",
		},
		{
			name: "empty name",
			account: &types.AccountCreate{
				Name: "",
			},
			wantErr: true,
			errMsg:  "account name is required",
		},
		{
			name: "invalid organization",
			account: &types.AccountCreate{
				Name:         "test-account",
				Organization: "",
			},
			wantErr: false, // Organization can be empty
		},
		{
			name: "valid basic account",
			account: &types.AccountCreate{
				Name:        "test-account",
				Description: "Test account",
			},
			wantErr: false,
		},
		{
			name: "valid complex account",
			account: &types.AccountCreate{
				Name:         "complex-account",
				Description:  "Complex test account",
				Organization: "Test Org",
				ParentName:   "parent-account",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateAccountCreate(tt.account)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAccountBaseManager_ValidateAccountUpdate(t *testing.T) {
	manager := NewAccountBaseManager("v0.0.43")

	tests := []struct {
		name    string
		update  *types.AccountUpdate
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil update",
			update:  nil,
			wantErr: true,
			errMsg:  "account update data is required",
		},
		{
			name: "valid update",
			update: &types.AccountUpdate{
				Description: stringPtr("Updated description"),
			},
			wantErr: false,
		},
		{
			name: "valid complex update",
			update: &types.AccountUpdate{
				Description:  stringPtr("Updated description"),
				Organization: stringPtr("Updated Org"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateAccountUpdate(tt.update)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAccountBaseManager_ApplyAccountDefaults(t *testing.T) {
	manager := NewAccountBaseManager("v0.0.43")

	tests := []struct {
		name     string
		input    *types.AccountCreate
		expected *types.AccountCreate
	}{
		{
			name: "apply defaults to minimal account",
			input: &types.AccountCreate{
				Name: "test-account",
			},
			expected: &types.AccountCreate{
				Name:        "test-account",
				Description: "",
			},
		},
		{
			name: "preserve existing values",
			input: &types.AccountCreate{
				Name:         "test-account",
				Description:  "Custom description",
				Organization: "Custom Org",
				ParentName:   "parent-account",
				Flags:        []string{"NoAccount", "NoDefaultQOS"},
			},
			expected: &types.AccountCreate{
				Name:         "test-account",
				Description:  "Custom description",
				Organization: "Custom Org",
				ParentName:   "parent-account",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.ApplyAccountDefaults(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAccountBaseManager_FilterAccountList(t *testing.T) {
	manager := NewAccountBaseManager("v0.0.43")

	accounts := []types.Account{
		{
			Name:         "account1",
			Description:  "Account 1",
			Organization: "Org A",
			ParentName:   "root",
		},
		{
			Name:         "account2",
			Description:  "Account 2",
			Organization: "Org B",
			ParentName:   "account1",
		},
		{
			Name:         "account3",
			Description:  "Account 3",
			Organization: "Org A",
			ParentName:   "root",
		},
	}

	tests := []struct {
		name     string
		opts     *types.AccountListOptions
		expected []string // expected account names
	}{
		{
			name:     "no filters",
			opts:     &types.AccountListOptions{},
			expected: []string{"account1", "account2", "account3"},
		},
		{
			name: "filter by names",
			opts: &types.AccountListOptions{
				Names: []string{"account1", "account3"},
			},
			expected: []string{"account1", "account3"},
		},
		{
			name: "filter by organization",
			opts: &types.AccountListOptions{
				Organizations: []string{"Org A"},
			},
			expected: []string{"account1", "account3"},
		},
		{
			name: "combined filters",
			opts: &types.AccountListOptions{
				Organizations: []string{"Org A"},
			},
			expected: []string{"account1", "account3"},
		},
		{
			name: "no matches",
			opts: &types.AccountListOptions{
				Names: []string{"nonexistent"},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.FilterAccountList(accounts, tt.opts)
			resultNames := make([]string, len(result))
			for i, account := range result {
				resultNames[i] = account.Name
			}
			assert.Equal(t, tt.expected, resultNames)
		})
	}
}

func TestAccountBaseManager_ValidateAccountName(t *testing.T) {
	manager := NewAccountBaseManager("v0.0.43")

	tests := []struct {
		name        string
		accountName string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid account name",
			accountName: "test-account",
			wantErr:     false,
		},
		{
			name:        "empty account name",
			accountName: "",
			wantErr:     true,
			errMsg:      "account name is required",
		},
		{
			name:        "account name with special chars",
			accountName: "test_account-123",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateAccountName(tt.accountName)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAccountBaseManager_ValidateAccountHierarchy(t *testing.T) {
	manager := NewAccountBaseManager("v0.0.43")

	// Mock existing accounts for hierarchy validation
	existingAccounts := []types.Account{
		{Name: "root", ParentName: ""},
		{Name: "parent1", ParentName: "root"},
		{Name: "child1", ParentName: "parent1"},
	}

	tests := []struct {
		name        string
		accountName string
		parentName  string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid hierarchy",
			accountName: "new-child",
			parentName:  "parent1",
			wantErr:     false,
		},
		{
			name:        "self as parent",
			accountName: "test-account",
			parentName:  "test-account",
			wantErr:     true,
			errMsg:      "cannot be its own parent",
		},
		{
			name:        "empty parent (root account)",
			accountName: "new-root",
			parentName:  "",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateAccountHierarchy(tt.accountName, tt.parentName, existingAccounts)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAccountBaseManager_BuildAccountTree(t *testing.T) {
	manager := NewAccountBaseManager("v0.0.43")

	accounts := []types.Account{
		{Name: "root", ParentName: ""},
		{Name: "parent1", ParentName: "root"},
		{Name: "parent2", ParentName: "root"},
		{Name: "child1", ParentName: "parent1"},
		{Name: "child2", ParentName: "parent1"},
		{Name: "grandchild1", ParentName: "child1"},
	}

	tree := manager.BuildAccountTree(accounts)

	// Verify tree structure
	assert.Contains(t, tree, "root")
	assert.Len(t, tree["root"], 2) // parent1, parent2
	assert.Contains(t, tree["parent1"], "child1")
	assert.Contains(t, tree["parent1"], "child2")
	assert.Contains(t, tree["child1"], "grandchild1")
}

func TestAccountBaseManager_GetAccountDepth(t *testing.T) {
	manager := NewAccountBaseManager("v0.0.43")

	tree := map[string][]string{
		"root":        {"parent1", "parent2"},
		"parent1":     {"child1", "child2"},
		"child1":      {"grandchild1"},
		"grandchild1": {},
	}

	tests := []struct {
		name        string
		accountName string
		expected    int
	}{
		{
			name:        "root account",
			accountName: "root",
			expected:    0,
		},
		{
			name:        "parent account",
			accountName: "parent1",
			expected:    1,
		},
		{
			name:        "child account",
			accountName: "child1",
			expected:    2,
		},
		{
			name:        "grandchild account",
			accountName: "grandchild1",
			expected:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			depth := manager.GetAccountDepth(tt.accountName, tree)
			assert.Equal(t, tt.expected, depth)
		})
	}
}

