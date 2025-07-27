package v0_0_43

import (
	"context"
	"testing"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

func TestUserManagerImpl_List(t *testing.T) {
	tests := []struct {
		name    string
		client  *WrapperClient
		opts    *interfaces.ListUsersOptions
		wantErr bool
		errType string
	}{
		{
			name:    "nil client",
			client:  nil,
			opts:    nil,
			wantErr: true,
			errType: "client error",
		},
		{
			name: "nil api client",
			client: &WrapperClient{
				apiClient: nil,
			},
			opts:    nil,
			wantErr: true,
			errType: "client error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			opts:    &interfaces.ListUsersOptions{},
			wantErr: true,
			errType: "client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserManagerImpl{
				client: tt.client,
			}
			
			result, err := u.List(context.Background(), tt.opts)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("List() expected error but got none")
					return
				}
				
				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("List() expected client error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("List() unexpected error: %v", err)
				}
			}
			
			if result != nil && tt.wantErr {
				t.Errorf("List() expected nil result on error")
			}
		})
	}
}

func TestUserManagerImpl_Get(t *testing.T) {
	tests := []struct {
		name     string
		client   *WrapperClient
		userName string
		wantErr  bool
		errType  string
	}{
		{
			name:     "nil client",
			client:   nil,
			userName: "testuser",
			wantErr:  true,
			errType:  "client error",
		},
		{
			name: "empty user name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName: "",
			wantErr:  true,
			errType:  "validation error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName: "testuser",
			wantErr:  true,
			errType:  "client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserManagerImpl{
				client: tt.client,
			}
			
			result, err := u.Get(context.Background(), tt.userName)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Get() expected error but got none")
					return
				}
				
				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("Get() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("Get() expected validation error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Get() unexpected error: %v", err)
				}
			}
			
			if result != nil && tt.wantErr {
				t.Errorf("Get() expected nil result on error")
			}
		})
	}
}

func TestUserManagerImpl_GetUserAccounts(t *testing.T) {
	tests := []struct {
		name     string
		client   *WrapperClient
		userName string
		wantErr  bool
		errType  string
	}{
		{
			name:     "nil client",
			client:   nil,
			userName: "testuser",
			wantErr:  true,
			errType:  "client error",
		},
		{
			name: "empty user name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName: "",
			wantErr:  true,
			errType:  "validation error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName: "testuser",
			wantErr:  true,
			errType:  "client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserManagerImpl{
				client: tt.client,
			}
			
			result, err := u.GetUserAccounts(context.Background(), tt.userName)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetUserAccounts() expected error but got none")
					return
				}
				
				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("GetUserAccounts() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("GetUserAccounts() expected validation error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("GetUserAccounts() unexpected error: %v", err)
				}
			}
			
			if result != nil && tt.wantErr {
				t.Errorf("GetUserAccounts() expected nil result on error")
			}
		})
	}
}

func TestUserManagerImpl_GetUserQuotas(t *testing.T) {
	tests := []struct {
		name     string
		client   *WrapperClient
		userName string
		wantErr  bool
		errType  string
	}{
		{
			name:     "nil client",
			client:   nil,
			userName: "testuser",
			wantErr:  true,
			errType:  "client error",
		},
		{
			name: "empty user name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName: "",
			wantErr:  true,
			errType:  "validation error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName: "testuser",
			wantErr:  true,
			errType:  "client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserManagerImpl{
				client: tt.client,
			}
			
			result, err := u.GetUserQuotas(context.Background(), tt.userName)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetUserQuotas() expected error but got none")
					return
				}
				
				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("GetUserQuotas() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("GetUserQuotas() expected validation error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("GetUserQuotas() unexpected error: %v", err)
				}
			}
			
			if result != nil && tt.wantErr {
				t.Errorf("GetUserQuotas() expected nil result on error")
			}
		})
	}
}

func TestUserManagerImpl_GetUserDefaultAccount(t *testing.T) {
	tests := []struct {
		name     string
		client   *WrapperClient
		userName string
		wantErr  bool
		errType  string
	}{
		{
			name:     "nil client",
			client:   nil,
			userName: "testuser",
			wantErr:  true,
			errType:  "client error",
		},
		{
			name: "empty user name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName: "",
			wantErr:  true,
			errType:  "validation error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName: "testuser",
			wantErr:  true,
			errType:  "client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserManagerImpl{
				client: tt.client,
			}
			
			result, err := u.GetUserDefaultAccount(context.Background(), tt.userName)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetUserDefaultAccount() expected error but got none")
					return
				}
				
				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("GetUserDefaultAccount() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("GetUserDefaultAccount() expected validation error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("GetUserDefaultAccount() unexpected error: %v", err)
				}
			}
			
			if result != nil && tt.wantErr {
				t.Errorf("GetUserDefaultAccount() expected nil result on error")
			}
		})
	}
}

func TestUserManagerImpl_GetUserFairShare(t *testing.T) {
	tests := []struct {
		name     string
		client   *WrapperClient
		userName string
		wantErr  bool
		errType  string
	}{
		{
			name:     "nil client",
			client:   nil,
			userName: "testuser",
			wantErr:  true,
			errType:  "client error",
		},
		{
			name: "empty user name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName: "",
			wantErr:  true,
			errType:  "validation error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName: "testuser",
			wantErr:  true,
			errType:  "client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserManagerImpl{
				client: tt.client,
			}
			
			result, err := u.GetUserFairShare(context.Background(), tt.userName)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetUserFairShare() expected error but got none")
					return
				}
				
				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("GetUserFairShare() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("GetUserFairShare() expected validation error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("GetUserFairShare() unexpected error: %v", err)
				}
			}
			
			if result != nil && tt.wantErr {
				t.Errorf("GetUserFairShare() expected nil result on error")
			}
		})
	}
}

func TestUserManagerImpl_CalculateJobPriority(t *testing.T) {
	tests := []struct {
		name          string
		client        *WrapperClient
		userName      string
		jobSubmission *interfaces.JobSubmission
		wantErr       bool
		errType       string
	}{
		{
			name:          "nil client",
			client:        nil,
			userName:      "testuser",
			jobSubmission: &interfaces.JobSubmission{Script: "test.sh"},
			wantErr:       true,
			errType:       "client error",
		},
		{
			name: "empty user name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName:      "",
			jobSubmission: &interfaces.JobSubmission{Script: "test.sh"},
			wantErr:       true,
			errType:       "validation error",
		},
		{
			name: "nil job submission",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName:      "testuser",
			jobSubmission: nil,
			wantErr:       true,
			errType:       "validation error",
		},
		{
			name: "empty script and command",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName:      "testuser",
			jobSubmission: &interfaces.JobSubmission{},
			wantErr:       true,
			errType:       "validation error",
		},
		{
			name: "valid input with script - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName:      "testuser",
			jobSubmission: &interfaces.JobSubmission{Script: "test.sh"},
			wantErr:       true,
			errType:       "client error",
		},
		{
			name: "valid input with command - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName:      "testuser",
			jobSubmission: &interfaces.JobSubmission{Command: "echo hello"},
			wantErr:       true,
			errType:       "client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserManagerImpl{
				client: tt.client,
			}
			
			result, err := u.CalculateJobPriority(context.Background(), tt.userName, tt.jobSubmission)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("CalculateJobPriority() expected error but got none")
					return
				}
				
				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("CalculateJobPriority() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("CalculateJobPriority() expected validation error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("CalculateJobPriority() unexpected error: %v", err)
				}
			}
			
			if result != nil && tt.wantErr {
				t.Errorf("CalculateJobPriority() expected nil result on error")
			}
		})
	}
}

func TestValidateUserName(t *testing.T) {
	tests := []struct {
		name     string
		userName string
		wantErr  bool
	}{
		{
			name:     "valid user name",
			userName: "testuser123",
			wantErr:  false,
		},
		{
			name:     "empty user name",
			userName: "",
			wantErr:  true,
		},
		{
			name:     "user name with underscore",
			userName: "test_user",
			wantErr:  false,
		},
		{
			name:     "user name with hyphen",
			userName: "test-user",
			wantErr:  false,
		},
		{
			name:     "user name with invalid characters",
			userName: "test@user",
			wantErr:  true,
		},
		{
			name:     "user name with space",
			userName: "test user",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUserName(tt.userName)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("validateUserName() expected error but got none")
				}
				if !errors.IsValidationError(err) {
					t.Errorf("validateUserName() expected validation error, got %T: %v", err, err)
				}
			} else {
				if err != nil {
					t.Errorf("validateUserName() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateAccountContext(t *testing.T) {
	tests := []struct {
		name        string
		accountName string
		wantErr     bool
	}{
		{
			name:        "valid account name",
			accountName: "testaccount123",
			wantErr:     false,
		},
		{
			name:        "empty account name",
			accountName: "",
			wantErr:     true,
		},
		{
			name:        "account name with underscore",
			accountName: "test_account",
			wantErr:     false,
		},
		{
			name:        "account name with hyphen",
			accountName: "test-account",
			wantErr:     false,
		},
		{
			name:        "account name with invalid characters",
			accountName: "test@account",
			wantErr:     true,
		},
		{
			name:        "account name with space",
			accountName: "test account",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAccountContext(tt.accountName)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("validateAccountContext() expected error but got none")
				}
				if !errors.IsValidationError(err) {
					t.Errorf("validateAccountContext() expected validation error, got %T: %v", err, err)
				}
			} else {
				if err != nil {
					t.Errorf("validateAccountContext() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestUserManagerImpl_ValidateUserAccountAccess(t *testing.T) {
	tests := []struct {
		name        string
		client      *WrapperClient
		userName    string
		accountName string
		wantErr     bool
		errType     string
	}{
		{
			name:        "nil client",
			client:      nil,
			userName:    "testuser",
			accountName: "testaccount",
			wantErr:     true,
			errType:     "client error",
		},
		{
			name: "empty user name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName:    "",
			accountName: "testaccount",
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "empty account name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName:    "testuser",
			accountName: "",
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName:    "testuser",
			accountName: "testaccount",
			wantErr:     true,
			errType:     "client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserManagerImpl{
				client: tt.client,
			}
			
			result, err := u.ValidateUserAccountAccess(context.Background(), tt.userName, tt.accountName)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateUserAccountAccess() expected error but got none")
					return
				}
				
				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("ValidateUserAccountAccess() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("ValidateUserAccountAccess() expected validation error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("ValidateUserAccountAccess() unexpected error: %v", err)
				}
			}
			
			if result != nil && tt.wantErr {
				t.Errorf("ValidateUserAccountAccess() expected nil result on error")
			}
		})
	}
}

func TestUserManagerImpl_GetUserAccountAssociations(t *testing.T) {
	tests := []struct {
		name     string
		client   *WrapperClient
		userName string
		opts     *interfaces.ListUserAccountAssociationsOptions
		wantErr  bool
		errType  string
	}{
		{
			name:     "nil client",
			client:   nil,
			userName: "testuser",
			opts:     nil,
			wantErr:  true,
			errType:  "client error",
		},
		{
			name: "empty user name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName: "",
			opts:     nil,
			wantErr:  true,
			errType:  "validation error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName: "testuser",
			opts: &interfaces.ListUserAccountAssociationsOptions{
				ActiveOnly: false,
			},
			wantErr: true,
			errType: "client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserManagerImpl{
				client: tt.client,
			}
			
			result, err := u.GetUserAccountAssociations(context.Background(), tt.userName, tt.opts)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetUserAccountAssociations() expected error but got none")
					return
				}
				
				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("GetUserAccountAssociations() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("GetUserAccountAssociations() expected validation error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("GetUserAccountAssociations() unexpected error: %v", err)
				}
			}
			
			if result != nil && tt.wantErr {
				t.Errorf("GetUserAccountAssociations() expected nil result on error")
			}
		})
	}
}

func TestUserManagerImpl_GetBulkUserAccounts(t *testing.T) {
	tests := []struct {
		name      string
		client    *WrapperClient
		userNames []string
		wantErr   bool
		errType   string
	}{
		{
			name:      "nil client",
			client:    nil,
			userNames: []string{"user1", "user2"},
			wantErr:   true,
			errType:   "client error",
		},
		{
			name: "empty user names list",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userNames: []string{},
			wantErr:   true,
			errType:   "validation error",
		},
		{
			name: "empty user name in list",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userNames: []string{"user1", "", "user3"},
			wantErr:   true,
			errType:   "validation error",
		},
		{
			name: "invalid user name format",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userNames: []string{"user1", "user@invalid", "user3"},
			wantErr:   true,
			errType:   "validation error",
		},
		{
			name: "too many users (over limit)",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userNames: make([]string, 101), // Over the 100 limit
			wantErr:   true,
			errType:   "validation error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userNames: []string{"user1", "user2", "user_3", "user-4"},
			wantErr:   true,
			errType:   "client error",
		},
	}

	// Initialize the slice with valid user names for the "too many users" test
	for i := range tests[4].userNames {
		tests[4].userNames[i] = "user" + string(rune('0'+i%10))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserManagerImpl{
				client: tt.client,
			}
			
			result, err := u.GetBulkUserAccounts(context.Background(), tt.userNames)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetBulkUserAccounts() expected error but got none")
					return
				}
				
				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("GetBulkUserAccounts() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("GetBulkUserAccounts() expected validation error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("GetBulkUserAccounts() unexpected error: %v", err)
				}
			}
			
			if result != nil && tt.wantErr {
				t.Errorf("GetBulkUserAccounts() expected nil result on error")
			}
		})
	}
}

func TestUserManagerImpl_GetBulkAccountUsers(t *testing.T) {
	tests := []struct {
		name         string
		client       *WrapperClient
		accountNames []string
		wantErr      bool
		errType      string
	}{
		{
			name:         "nil client",
			client:       nil,
			accountNames: []string{"account1", "account2"},
			wantErr:      true,
			errType:      "client error",
		},
		{
			name: "empty account names list",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountNames: []string{},
			wantErr:      true,
			errType:      "validation error",
		},
		{
			name: "empty account name in list",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountNames: []string{"account1", "", "account3"},
			wantErr:      true,
			errType:      "validation error",
		},
		{
			name: "invalid account name format",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountNames: []string{"account1", "account@invalid", "account3"},
			wantErr:      true,
			errType:      "validation error",
		},
		{
			name: "too many accounts (over limit)",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountNames: make([]string, 101), // Over the 100 limit
			wantErr:      true,
			errType:      "validation error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountNames: []string{"account1", "account2", "account_3", "account-4"},
			wantErr:      true,
			errType:      "client error",
		},
	}

	// Initialize the slice with valid account names for the "too many accounts" test
	for i := range tests[4].accountNames {
		tests[4].accountNames[i] = "account" + string(rune('0'+i%10))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserManagerImpl{
				client: tt.client,
			}
			
			result, err := u.GetBulkAccountUsers(context.Background(), tt.accountNames)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetBulkAccountUsers() expected error but got none")
					return
				}
				
				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("GetBulkAccountUsers() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("GetBulkAccountUsers() expected validation error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("GetBulkAccountUsers() unexpected error: %v", err)
				}
			}
			
			if result != nil && tt.wantErr {
				t.Errorf("GetBulkAccountUsers() expected nil result on error")
			}
		})
	}
}