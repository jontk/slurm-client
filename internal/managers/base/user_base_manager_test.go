package base

import (
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserBaseManager_New(t *testing.T) {
	manager := NewUserBaseManager("v0.0.43")
	assert.NotNil(t, manager)
	assert.Equal(t, "v0.0.43", manager.GetVersion())
	assert.Equal(t, "User", manager.GetResourceType())
}

func TestUserBaseManager_ValidateUserCreate(t *testing.T) {
	manager := NewUserBaseManager("v0.0.43")

	tests := []struct {
		name    string
		user    *types.UserCreate
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil user",
			user:    nil,
			wantErr: true,
			errMsg:  "user data is required",
		},
		{
			name: "empty name",
			user: &types.UserCreate{
				Name: "",
			},
			wantErr: true,
			errMsg:  "user name is required",
		},
		{
			name: "empty default account",
			user: &types.UserCreate{
				Name:           "testuser",
				DefaultAccount: "",
			},
			wantErr: true,
			errMsg:  "default account is required",
		},
		{
			name: "valid basic user",
			user: &types.UserCreate{
				Name:           "testuser",
				DefaultAccount: "test-account",
			},
			wantErr: false,
		},
		{
			name: "valid complex user",
			user: &types.UserCreate{
				Name:           "testuser",
				DefaultAccount: "test-account",
				Accounts:       []string{"test-account", "other-account"},
				DefaultQoS:     "normal",
				QoSList:        []string{"normal", "high"},
				AdminLevel:     "Operator",
				Flags:          []string{"NoDefaultAccount"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateUserCreate(tt.user)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserBaseManager_ValidateUserUpdate(t *testing.T) {
	manager := NewUserBaseManager("v0.0.43")

	tests := []struct {
		name    string
		update  *types.UserUpdate
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil update",
			update:  nil,
			wantErr: true,
			errMsg:  "user update data is required",
		},
		{
			name: "empty name",
			update: &types.UserUpdate{
				Name: "",
			},
			wantErr: true,
			errMsg:  "user name is required",
		},
		{
			name: "valid update",
			update: &types.UserUpdate{
				Name:           "testuser",
				DefaultAccount: stringPtr("new-account"),
			},
			wantErr: false,
		},
		{
			name: "valid complex update",
			update: &types.UserUpdate{
				Name:           "testuser",
				DefaultAccount: stringPtr("new-account"),
				DefaultQoS:     stringPtr("high"),
				AdminLevel:     stringPtr("Administrator"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateUserUpdate(tt.update)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserBaseManager_ApplyUserDefaults(t *testing.T) {
	manager := NewUserBaseManager("v0.0.43")

	tests := []struct {
		name     string
		input    *types.UserCreate
		expected *types.UserCreate
	}{
		{
			name: "apply defaults to minimal user",
			input: &types.UserCreate{
				Name:           "testuser",
				DefaultAccount: "test-account",
			},
			expected: &types.UserCreate{
				Name:           "testuser",
				DefaultAccount: "test-account",
				Accounts:       []string{"test-account"}, // Default to default account
				AdminLevel:     "None",                   // Default admin level
				Flags:          []string{},               // Empty flags
				QoSList:        []string{},               // Empty QoS list
			},
		},
		{
			name: "preserve existing values",
			input: &types.UserCreate{
				Name:           "testuser",
				DefaultAccount: "test-account",
				Accounts:       []string{"test-account", "other-account"},
				DefaultQoS:     "high",
				QoSList:        []string{"normal", "high", "urgent"},
				AdminLevel:     "Operator",
				Flags:          []string{"NoDefaultAccount", "NoDefaultQOS"},
			},
			expected: &types.UserCreate{
				Name:           "testuser",
				DefaultAccount: "test-account",
				Accounts:       []string{"test-account", "other-account"},
				DefaultQoS:     "high",
				QoSList:        []string{"normal", "high", "urgent"},
				AdminLevel:     "Operator",
				Flags:          []string{"NoDefaultAccount", "NoDefaultQOS"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.ApplyUserDefaults(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUserBaseManager_FilterUserList(t *testing.T) {
	manager := NewUserBaseManager("v0.0.43")

	users := []types.User{
		{
			Name:           "user1",
			DefaultAccount: "account1",
			Accounts:       []string{"account1", "account2"},
			DefaultQoS:     "normal",
			QoSList:        []string{"normal", "high"},
			AdminLevel:     "None",
			Flags:          []string{},
		},
		{
			Name:           "user2",
			DefaultAccount: "account2",
			Accounts:       []string{"account2"},
			DefaultQoS:     "high",
			QoSList:        []string{"high", "urgent"},
			AdminLevel:     "Operator",
			Flags:          []string{"NoDefaultAccount"},
		},
		{
			Name:           "user3",
			DefaultAccount: "account1",
			Accounts:       []string{"account1"},
			DefaultQoS:     "normal",
			QoSList:        []string{"normal"},
			AdminLevel:     "Administrator",
			Flags:          []string{"NoDefaultQOS"},
		},
	}

	tests := []struct {
		name     string
		opts     *types.UserListOptions
		expected []string // expected user names
	}{
		{
			name:     "no filters",
			opts:     &types.UserListOptions{},
			expected: []string{"user1", "user2", "user3"},
		},
		{
			name: "filter by names",
			opts: &types.UserListOptions{
				Names: []string{"user1", "user3"},
			},
			expected: []string{"user1", "user3"},
		},
		{
			name: "filter by default account",
			opts: &types.UserListOptions{
				DefaultAccounts: []string{"account1"},
			},
			expected: []string{"user1", "user3"},
		},
		{
			name: "filter by account",
			opts: &types.UserListOptions{
				Accounts: []string{"account2"},
			},
			expected: []string{"user1", "user2"},
		},
		{
			name: "filter by default QoS",
			opts: &types.UserListOptions{
				DefaultQoSList: []string{"high"},
			},
			expected: []string{"user2"},
		},
		{
			name: "filter by QoS",
			opts: &types.UserListOptions{
				QoSList: []string{"urgent"},
			},
			expected: []string{"user2"},
		},
		{
			name: "filter by admin level",
			opts: &types.UserListOptions{
				AdminLevels: []string{"Operator", "Administrator"},
			},
			expected: []string{"user2", "user3"},
		},
		{
			name: "filter by flag",
			opts: &types.UserListOptions{
				WithFlags: []string{"NoDefaultAccount"},
			},
			expected: []string{"user2"},
		},
		{
			name: "combined filters",
			opts: &types.UserListOptions{
				DefaultAccounts: []string{"account1"},
				AdminLevels:     []string{"Administrator"},
			},
			expected: []string{"user3"},
		},
		{
			name: "no matches",
			opts: &types.UserListOptions{
				Names: []string{"nonexistent"},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.FilterUserList(users, tt.opts)
			resultNames := make([]string, len(result))
			for i, user := range result {
				resultNames[i] = user.Name
			}
			assert.Equal(t, tt.expected, resultNames)
		})
	}
}

func TestUserBaseManager_ValidateUserName(t *testing.T) {
	manager := NewUserBaseManager("v0.0.43")

	tests := []struct {
		name     string
		username string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid username",
			username: "testuser",
			wantErr:  false,
		},
		{
			name:     "empty username",
			username: "",
			wantErr:  true,
			errMsg:   "user name is required",
		},
		{
			name:     "username with numbers",
			username: "user123",
			wantErr:  false,
		},
		{
			name:     "username with underscores",
			username: "test_user",
			wantErr:  false,
		},
		{
			name:     "username with hyphens",
			username: "test-user",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateUserName(tt.username)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserBaseManager_ValidateAdminLevel(t *testing.T) {
	manager := NewUserBaseManager("v0.0.43")

	tests := []struct {
		name       string
		adminLevel string
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "valid None level",
			adminLevel: "None",
			wantErr:    false,
		},
		{
			name:       "valid Operator level",
			adminLevel: "Operator",
			wantErr:    false,
		},
		{
			name:       "valid Administrator level",
			adminLevel: "Administrator",
			wantErr:    false,
		},
		{
			name:       "invalid admin level",
			adminLevel: "InvalidLevel",
			wantErr:    true,
			errMsg:     "invalid admin level",
		},
		{
			name:       "empty admin level",
			adminLevel: "",
			wantErr:    false, // Empty defaults to None
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateAdminLevel(tt.adminLevel)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserBaseManager_ValidateUserAccounts(t *testing.T) {
	manager := NewUserBaseManager("v0.0.43")

	// Mock available accounts
	availableAccounts := []string{"account1", "account2", "account3"}

	tests := []struct {
		name        string
		accounts    []string
		defaultAcct string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid accounts with default",
			accounts:    []string{"account1", "account2"},
			defaultAcct: "account1",
			wantErr:     false,
		},
		{
			name:        "default account not in accounts list",
			accounts:    []string{"account2"},
			defaultAcct: "account1",
			wantErr:     true,
			errMsg:      "default account must be in accounts list",
		},
		{
			name:        "nonexistent account",
			accounts:    []string{"nonexistent"},
			defaultAcct: "nonexistent",
			wantErr:     true,
			errMsg:      "account does not exist",
		},
		{
			name:        "empty accounts list",
			accounts:    []string{},
			defaultAcct: "account1",
			wantErr:     true,
			errMsg:      "user must have at least one account",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateUserAccounts(tt.accounts, tt.defaultAcct, availableAccounts)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserBaseManager_NormalizeUserFlags(t *testing.T) {
	manager := NewUserBaseManager("v0.0.43")

	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty flags",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "valid flags",
			input:    []string{"NoDefaultAccount", "NoDefaultQOS"},
			expected: []string{"NoDefaultAccount", "NoDefaultQOS"},
		},
		{
			name:     "remove duplicates",
			input:    []string{"NoDefaultAccount", "NoDefaultQOS", "NoDefaultAccount"},
			expected: []string{"NoDefaultAccount", "NoDefaultQOS"},
		},
		{
			name:     "normalize case",
			input:    []string{"nodefaultaccount", "NODEFAULTQOS"},
			expected: []string{"NoDefaultAccount", "NoDefaultQOS"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.NormalizeUserFlags(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUserBaseManager_GetUserPermissionLevel(t *testing.T) {
	manager := NewUserBaseManager("v0.0.43")

	tests := []struct {
		name        string
		adminLevel  string
		expected    int
		description string
	}{
		{
			name:        "None level",
			adminLevel:  "None",
			expected:    0,
			description: "Regular user",
		},
		{
			name:        "Operator level",
			adminLevel:  "Operator",
			expected:    1,
			description: "Can manage users and accounts",
		},
		{
			name:        "Administrator level",
			adminLevel:  "Administrator",
			expected:    2,
			description: "Full system access",
		},
		{
			name:        "Unknown level",
			adminLevel:  "Unknown",
			expected:    0,
			description: "Default to regular user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, desc := manager.GetUserPermissionLevel(tt.adminLevel)
			assert.Equal(t, tt.expected, level)
			assert.Equal(t, tt.description, desc)
		})
	}
}