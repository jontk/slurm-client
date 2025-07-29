package v0_0_43

import (
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserAdapter_ValidateUserCreate(t *testing.T) {
	adapter := &UserAdapter{
		UserBaseManager: base.NewUserBaseManager("v0.0.43"),
	}

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
			name: "invalid name with spaces",
			user: &types.UserCreate{
				Name: "user name",
			},
			wantErr: true,
			errMsg:  "user name cannot contain spaces",
		},
		{
			name: "invalid name with special characters",
			user: &types.UserCreate{
				Name: "user@domain",
			},
			wantErr: true,
			errMsg:  "user name contains invalid characters",
		},
		{
			name: "name too long",
			user: &types.UserCreate{
				Name: "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz", // 52 characters
			},
			wantErr: true,
			errMsg:  "user name too long",
		},
		{
			name: "valid basic user",
			user: &types.UserCreate{
				Name: "testuser",
			},
			wantErr: false,
		},
		{
			name: "valid user with valid characters",
			user: &types.UserCreate{
				Name: "test-user_123",
			},
			wantErr: false,
		},
		{
			name: "valid complex user",
			user: &types.UserCreate{
				Name:           "jdoe",
				RealName:       "John Doe",
				Email:          "john.doe@example.com",
				DefaultAccount: "physics",
				AdminLevel:     "Operator",
				Associations: []types.UserAssociation{
					{
						Account:   "physics",
						Cluster:   "main",
						Partition: "compute",
						QoS:       "normal",
					},
				},
				Coordinators: []string{"coordinator1", "coordinator2"},
			},
			wantErr: false,
		},
		{
			name: "invalid email format",
			user: &types.UserCreate{
				Name:  "testuser",
				Email: "invalid-email",
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name: "invalid admin level",
			user: &types.UserCreate{
				Name:       "testuser",
				AdminLevel: "InvalidLevel",
			},
			wantErr: true,
			errMsg:  "invalid admin level",
		},
		{
			name: "empty default account",
			user: &types.UserCreate{
				Name:           "testuser",
				DefaultAccount: "",
			},
			wantErr: false, // Empty default account is valid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateUserCreate(tt.user)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserAdapter_ApplyUserDefaults(t *testing.T) {
	adapter := &UserAdapter{
		UserBaseManager: base.NewUserBaseManager("v0.0.43"),
	}

	tests := []struct {
		name     string
		input    *types.UserCreate
		expected *types.UserCreate
	}{
		{
			name: "apply defaults to minimal user",
			input: &types.UserCreate{
				Name: "testuser",
			},
			expected: &types.UserCreate{
				Name:           "testuser",
				RealName:       "",                           // Empty real name
				Email:          "",                           // Empty email
				DefaultAccount: "",                           // Empty default account
				AdminLevel:     "None",                       // Default admin level
				Associations:   []types.UserAssociation{},   // Empty associations
				Coordinators:   []string{},                  // Empty coordinators
			},
		},
		{
			name: "preserve existing values",
			input: &types.UserCreate{
				Name:           "jdoe",
				RealName:       "John Doe",
				Email:          "john.doe@example.com",
				DefaultAccount: "physics",
				AdminLevel:     "Administrator",
				Associations: []types.UserAssociation{
					{
						Account:   "physics",
						Cluster:   "main",
						Partition: "compute",
						QoS:       "high",
					},
				},
				Coordinators: []string{"coord1"},
			},
			expected: &types.UserCreate{
				Name:           "jdoe",
				RealName:       "John Doe",
				Email:          "john.doe@example.com",
				DefaultAccount: "physics",
				AdminLevel:     "Administrator",
				Associations: []types.UserAssociation{
					{
						Account:   "physics",
						Cluster:   "main",
						Partition: "compute",
						QoS:       "high",
					},
				},
				Coordinators: []string{"coord1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.ApplyUserDefaults(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUserAdapter_FilterUserList(t *testing.T) {
	adapter := &UserAdapter{
		UserBaseManager: base.NewUserBaseManager("v0.0.43"),
	}

	users := []types.User{
		{
			Name:           "alice",
			RealName:       "Alice Smith",
			Email:          "alice@physics.edu",
			DefaultAccount: "physics",
			AdminLevel:     "None",
			Associations: []types.UserAssociation{
				{Account: "physics", Cluster: "main", QoS: "normal"},
			},
		},
		{
			Name:           "bob",
			RealName:       "Bob Johnson",
			Email:          "bob@chemistry.edu",
			DefaultAccount: "chemistry",
			AdminLevel:     "Operator",
			Associations: []types.UserAssociation{
				{Account: "chemistry", Cluster: "main", QoS: "high"},
			},
		},
		{
			Name:           "charlie",
			RealName:       "Charlie Brown",
			Email:          "charlie@admin.edu",
			DefaultAccount: "admin",
			AdminLevel:     "Administrator",
			Associations: []types.UserAssociation{
				{Account: "admin", Cluster: "main", QoS: "normal"},
				{Account: "physics", Cluster: "main", QoS: "high"},
			},
		},
		{
			Name:           "diana",
			RealName:       "Diana Prince",
			Email:          "diana@physics.edu",
			DefaultAccount: "physics",
			AdminLevel:     "None",
			Associations: []types.UserAssociation{
				{Account: "physics", Cluster: "backup", QoS: "normal"},
			},
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
			expected: []string{"alice", "bob", "charlie", "diana"},
		},
		{
			name: "filter by names",
			opts: &types.UserListOptions{
				Names: []string{"alice", "bob"},
			},
			expected: []string{"alice", "bob"},
		},
		{
			name: "filter by default account",
			opts: &types.UserListOptions{
				DefaultAccounts: []string{"physics"},
			},
			expected: []string{"alice", "diana"},
		},
		{
			name: "filter by admin level",
			opts: &types.UserListOptions{
				AdminLevels: []string{"Administrator", "Operator"},
			},
			expected: []string{"bob", "charlie"},
		},
		{
			name: "filter by association account",
			opts: &types.UserListOptions{
				AssociationAccounts: []string{"chemistry"},
			},
			expected: []string{"bob"},
		},
		{
			name: "filter by association cluster",
			opts: &types.UserListOptions{
				AssociationClusters: []string{"backup"},
			},
			expected: []string{"diana"},
		},
		{
			name: "filter by association QoS",
			opts: &types.UserListOptions{
				AssociationQoSList: []string{"high"},
			},
			expected: []string{"bob", "charlie"},
		},
		{
			name: "combined filters",
			opts: &types.UserListOptions{
				DefaultAccounts: []string{"physics"},
				AdminLevels:     []string{"None"},
			},
			expected: []string{"alice", "diana"},
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
			result := adapter.FilterUserList(users, tt.opts)
			resultNames := make([]string, len(result))
			for i, user := range result {
				resultNames[i] = user.Name
			}
			assert.Equal(t, tt.expected, resultNames)
		})
	}
}

func TestUserAdapter_ValidateUserName(t *testing.T) {
	adapter := &UserAdapter{
		UserBaseManager: base.NewUserBaseManager("v0.0.43"),
	}

	tests := []struct {
		name     string
		userName string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid simple name",
			userName: "testuser",
			wantErr:  false,
		},
		{
			name:     "valid name with numbers",
			userName: "user123",
			wantErr:  false,
		},
		{
			name:     "valid name with underscore",
			userName: "test_user",
			wantErr:  false,
		},
		{
			name:     "valid name with dash",
			userName: "test-user",
			wantErr:  false,
		},
		{
			name:     "valid name with dot",
			userName: "test.user",
			wantErr:  false,
		},
		{
			name:     "empty name",
			userName: "",
			wantErr:  true,
			errMsg:   "user name is required",
		},
		{
			name:     "name with spaces",
			userName: "test user",
			wantErr:  true,
			errMsg:   "user name cannot contain spaces",
		},
		{
			name:     "name with special characters",
			userName: "test@user",
			wantErr:  true,
			errMsg:   "user name contains invalid characters",
		},
		{
			name:     "name starting with number",
			userName: "123user",
			wantErr:  true,
			errMsg:   "user name cannot start with a number",
		},
		{
			name:     "name starting with dash",
			userName: "-testuser",
			wantErr:  true,
			errMsg:   "user name cannot start with a dash",
		},
		{
			name:     "name too long",
			userName: "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz",
			wantErr:  true,
			errMsg:   "user name too long",
		},
		{
			name:     "valid maximum length name",
			userName: "abcdefghijklmnopqrstuvwxyzabcdefg", // 32 characters
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateUserName(tt.userName)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserAdapter_ValidateEmailFormat(t *testing.T) {
	adapter := &UserAdapter{
		UserBaseManager: base.NewUserBaseManager("v0.0.43"),
	}

	tests := []struct {
		name    string
		email   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid email",
			email:   "user@example.com",
			wantErr: false,
		},
		{
			name:    "valid email with subdomain",
			email:   "user@mail.example.com",
			wantErr: false,
		},
		{
			name:    "valid email with plus",
			email:   "user+tag@example.com",
			wantErr: false,
		},
		{
			name:    "valid email with dots",
			email:   "first.last@example.com",
			wantErr: false,
		},
		{
			name:    "valid email with numbers",
			email:   "user123@example123.com",
			wantErr: false,
		},
		{
			name:    "empty email (valid)",
			email:   "",
			wantErr: false,
		},
		{
			name:    "invalid email without @",
			email:   "userexample.com",
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name:    "invalid email without domain",
			email:   "user@",
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name:    "invalid email without local part",
			email:   "@example.com",
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name:    "invalid email with spaces",
			email:   "user @example.com",
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name:    "invalid email with invalid characters",
			email:   "user<>@example.com",
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name:    "invalid email too long",
			email:   "verylongusernamethatexceedsthelimitofcharactersallowed@verylongdomainnamethatexceedsthelimitofcharactersallowed.com",
			wantErr: true,
			errMsg:  "email too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateEmailFormat(tt.email)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserAdapter_ValidateAdminLevel(t *testing.T) {
	adapter := &UserAdapter{
		UserBaseManager: base.NewUserBaseManager("v0.0.43"),
	}

	tests := []struct {
		name       string
		adminLevel string
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "valid None",
			adminLevel: "None",
			wantErr:    false,
		},
		{
			name:       "valid Operator",
			adminLevel: "Operator",
			wantErr:    false,
		},
		{
			name:       "valid Administrator",
			adminLevel: "Administrator",
			wantErr:    false,
		},
		{
			name:       "empty admin level (should use default)",
			adminLevel: "",
			wantErr:    false,
		},
		{
			name:       "invalid admin level",
			adminLevel: "SuperUser",
			wantErr:    true,
			errMsg:     "invalid admin level",
		},
		{
			name:       "lowercase admin level",
			adminLevel: "administrator",
			wantErr:    true,
			errMsg:     "invalid admin level",
		},
		{
			name:       "mixed case admin level",
			adminLevel: "operator",
			wantErr:    true,
			errMsg:     "invalid admin level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateAdminLevel(tt.adminLevel)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserAdapter_ValidateUserAssociations(t *testing.T) {
	adapter := &UserAdapter{
		UserBaseManager: base.NewUserBaseManager("v0.0.43"),
	}

	tests := []struct {
		name         string
		associations []types.UserAssociation
		wantErr      bool
		errMsg       string
	}{
		{
			name: "valid associations",
			associations: []types.UserAssociation{
				{
					Account:   "physics",
					Cluster:   "main",
					Partition: "compute",
					QoS:       "normal",
				},
				{
					Account:   "chemistry",
					Cluster:   "backup",
					Partition: "gpu",
					QoS:       "high",
				},
			},
			wantErr: false,
		},
		{
			name:         "empty associations",
			associations: []types.UserAssociation{},
			wantErr:      false,
		},
		{
			name: "association with missing account",
			associations: []types.UserAssociation{
				{
					Account:   "",
					Cluster:   "main",
					Partition: "compute",
					QoS:       "normal",
				},
			},
			wantErr: true,
			errMsg:  "association account is required",
		},
		{
			name: "association with missing cluster",
			associations: []types.UserAssociation{
				{
					Account:   "physics",
					Cluster:   "",
					Partition: "compute",
					QoS:       "normal",
				},
			},
			wantErr: true,
			errMsg:  "association cluster is required",
		},
		{
			name: "duplicate associations",
			associations: []types.UserAssociation{
				{
					Account:   "physics",
					Cluster:   "main",
					Partition: "compute",
					QoS:       "normal",
				},
				{
					Account:   "physics",
					Cluster:   "main",
					Partition: "compute",
					QoS:       "high", // Different QoS but same account+cluster
				},
			},
			wantErr: true,
			errMsg:  "duplicate association",
		},
		{
			name: "association with empty partition (valid)",
			associations: []types.UserAssociation{
				{
					Account: "physics",
					Cluster: "main",
					QoS:     "normal",
				},
			},
			wantErr: false,
		},
		{
			name: "association with empty QoS (valid)",
			associations: []types.UserAssociation{
				{
					Account:   "physics",
					Cluster:   "main",
					Partition: "compute",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateUserAssociations(tt.associations)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserAdapter_ValidateCoordinators(t *testing.T) {
	adapter := &UserAdapter{
		UserBaseManager: base.NewUserBaseManager("v0.0.43"),
	}

	tests := []struct {
		name         string
		coordinators []string
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "valid coordinators",
			coordinators: []string{"coord1", "coord2", "coord3"},
			wantErr:      false,
		},
		{
			name:         "empty coordinators",
			coordinators: []string{},
			wantErr:      false,
		},
		{
			name:         "single coordinator",
			coordinators: []string{"coordinator"},
			wantErr:      false,
		},
		{
			name:         "duplicate coordinators",
			coordinators: []string{"coord1", "coord1", "coord2"},
			wantErr:      true,
			errMsg:       "duplicate coordinator",
		},
		{
			name:         "empty coordinator string",
			coordinators: []string{"coord1", "", "coord2"},
			wantErr:      true,
			errMsg:       "coordinator cannot be empty",
		},
		{
			name:         "coordinator with spaces",
			coordinators: []string{"coord with spaces"},
			wantErr:      true,
			errMsg:       "coordinator cannot contain spaces",
		},
		{
			name:         "coordinator with invalid characters",
			coordinators: []string{"coord@invalid"},
			wantErr:      true,
			errMsg:       "coordinator contains invalid characters",
		},
		{
			name:         "valid coordinators with allowed characters",
			coordinators: []string{"coord-1", "coord_2", "coord.3"},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateCoordinators(tt.coordinators)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}