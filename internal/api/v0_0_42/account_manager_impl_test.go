// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"testing"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

func TestAccountManagerImpl_List(t *testing.T) {
	tests := []struct {
		name    string
		client  *WrapperClient
		opts    *interfaces.ListAccountsOptions
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
			name: "valid input - not implemented",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			opts:    &interfaces.ListAccountsOptions{},
			wantErr: true,
			errType: "not implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccountManagerImpl{
				client: tt.client,
			}

			result, err := a.List(context.Background(), tt.opts)

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
				case "not implemented":
					if !errors.IsNotImplementedError(err) {
						t.Errorf("List() expected not implemented error, got %T: %v", err, err)
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

func TestAccountManagerImpl_Get(t *testing.T) {
	tests := []struct {
		name        string
		client      *WrapperClient
		accountName string
		wantErr     bool
		errType     string
	}{
		{
			name:        "nil client",
			client:      nil,
			accountName: "testaccount",
			wantErr:     true,
			errType:     "client error",
		},
		{
			name: "empty account name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "",
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "valid input - not implemented",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "testaccount",
			wantErr:     true,
			errType:     "not implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccountManagerImpl{
				client: tt.client,
			}

			result, err := a.Get(context.Background(), tt.accountName)

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
				case "not implemented":
					if !errors.IsNotImplementedError(err) {
						t.Errorf("Get() expected not implemented error, got %T: %v", err, err)
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

func TestAccountManagerImpl_Create(t *testing.T) {
	tests := []struct {
		name    string
		client  *WrapperClient
		account *interfaces.AccountCreate
		wantErr bool
		errType string
	}{
		{
			name:    "nil client",
			client:  nil,
			account: &interfaces.AccountCreate{Name: "testaccount"},
			wantErr: true,
			errType: "client error",
		},
		{
			name: "nil account data",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			account: nil,
			wantErr: true,
			errType: "validation error",
		},
		{
			name: "empty account name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			account: &interfaces.AccountCreate{Name: ""},
			wantErr: true,
			errType: "validation error",
		},
		{
			name: "valid input - not implemented",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			account: &interfaces.AccountCreate{
				Name:        "testaccount",
				Description: "Test Account",
			},
			wantErr: true,
			errType: "not implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccountManagerImpl{
				client: tt.client,
			}

			result, err := a.Create(context.Background(), tt.account)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Create() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("Create() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("Create() expected validation error, got %T: %v", err, err)
					}
				case "not implemented":
					if !errors.IsNotImplementedError(err) {
						t.Errorf("Create() expected not implemented error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Create() unexpected error: %v", err)
				}
			}

			if result != nil && tt.wantErr {
				t.Errorf("Create() expected nil result on error")
			}
		})
	}
}

func TestAccountManagerImpl_Update(t *testing.T) {
	tests := []struct {
		name        string
		client      *WrapperClient
		accountName string
		update      *interfaces.AccountUpdate
		wantErr     bool
		errType     string
	}{
		{
			name:        "nil client",
			client:      nil,
			accountName: "testaccount",
			update:      &interfaces.AccountUpdate{Description: stringPtr("Updated")},
			wantErr:     true,
			errType:     "client error",
		},
		{
			name: "empty account name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "",
			update:      &interfaces.AccountUpdate{Description: stringPtr("Updated")},
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "nil update data",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "testaccount",
			update:      nil,
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "valid input - not implemented",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "testaccount",
			update:      &interfaces.AccountUpdate{Description: stringPtr("Updated Account")},
			wantErr:     true,
			errType:     "not implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccountManagerImpl{
				client: tt.client,
			}

			err := a.Update(context.Background(), tt.accountName, tt.update)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Update() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("Update() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("Update() expected validation error, got %T: %v", err, err)
					}
				case "not implemented":
					if !errors.IsNotImplementedError(err) {
						t.Errorf("Update() expected not implemented error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Update() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestAccountManagerImpl_Delete(t *testing.T) {
	tests := []struct {
		name        string
		client      *WrapperClient
		accountName string
		wantErr     bool
		errType     string
	}{
		{
			name:        "nil client",
			client:      nil,
			accountName: "testaccount",
			wantErr:     true,
			errType:     "client error",
		},
		{
			name: "empty account name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "",
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "valid input - not implemented",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "testaccount",
			wantErr:     true,
			errType:     "not implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccountManagerImpl{
				client: tt.client,
			}

			err := a.Delete(context.Background(), tt.accountName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Delete() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("Delete() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("Delete() expected validation error, got %T: %v", err, err)
					}
				case "not implemented":
					if !errors.IsNotImplementedError(err) {
						t.Errorf("Delete() expected not implemented error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Delete() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestAccountManagerImpl_GetAccountHierarchy(t *testing.T) {
	tests := []struct {
		name        string
		client      *WrapperClient
		rootAccount string
		wantErr     bool
		errType     string
	}{
		{
			name:        "nil client",
			client:      nil,
			rootAccount: "root",
			wantErr:     true,
			errType:     "client error",
		},
		{
			name: "empty root account",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			rootAccount: "",
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "valid input - not implemented",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			rootAccount: "root",
			wantErr:     true,
			errType:     "not implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccountManagerImpl{
				client: tt.client,
			}

			result, err := a.GetAccountHierarchy(context.Background(), tt.rootAccount)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetAccountHierarchy() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("GetAccountHierarchy() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("GetAccountHierarchy() expected validation error, got %T: %v", err, err)
					}
				case "not implemented":
					if !errors.IsNotImplementedError(err) {
						t.Errorf("GetAccountHierarchy() expected not implemented error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("GetAccountHierarchy() unexpected error: %v", err)
				}
			}

			if result != nil && tt.wantErr {
				t.Errorf("GetAccountHierarchy() expected nil result on error")
			}
		})
	}
}

func TestAccountManagerImpl_GetParentAccounts(t *testing.T) {
	tests := []struct {
		name        string
		client      *WrapperClient
		accountName string
		wantErr     bool
		errType     string
	}{
		{
			name:        "nil client",
			client:      nil,
			accountName: "testaccount",
			wantErr:     true,
			errType:     "client error",
		},
		{
			name: "empty account name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "",
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "valid input - not implemented",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "testaccount",
			wantErr:     true,
			errType:     "not implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccountManagerImpl{
				client: tt.client,
			}

			result, err := a.GetParentAccounts(context.Background(), tt.accountName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetParentAccounts() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("GetParentAccounts() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("GetParentAccounts() expected validation error, got %T: %v", err, err)
					}
				case "not implemented":
					if !errors.IsNotImplementedError(err) {
						t.Errorf("GetParentAccounts() expected not implemented error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("GetParentAccounts() unexpected error: %v", err)
				}
			}

			if result != nil && tt.wantErr {
				t.Errorf("GetParentAccounts() expected nil result on error")
			}
		})
	}
}

func TestAccountManagerImpl_GetChildAccounts(t *testing.T) {
	tests := []struct {
		name        string
		client      *WrapperClient
		accountName string
		depth       int
		wantErr     bool
		errType     string
	}{
		{
			name:        "nil client",
			client:      nil,
			accountName: "testaccount",
			depth:       0,
			wantErr:     true,
			errType:     "client error",
		},
		{
			name: "empty account name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "",
			depth:       0,
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "negative depth",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "testaccount",
			depth:       -1,
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "valid input with unlimited depth - not implemented",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "testaccount",
			depth:       0,
			wantErr:     true,
			errType:     "not implemented",
		},
		{
			name: "valid input with limited depth - not implemented",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "testaccount",
			depth:       2,
			wantErr:     true,
			errType:     "not implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccountManagerImpl{
				client: tt.client,
			}

			result, err := a.GetChildAccounts(context.Background(), tt.accountName, tt.depth)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetChildAccounts() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("GetChildAccounts() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("GetChildAccounts() expected validation error, got %T: %v", err, err)
					}
				case "not implemented":
					if !errors.IsNotImplementedError(err) {
						t.Errorf("GetChildAccounts() expected not implemented error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("GetChildAccounts() unexpected error: %v", err)
				}
			}

			if result != nil && tt.wantErr {
				t.Errorf("GetChildAccounts() expected nil result on error")
			}
		})
	}
}

func TestAccountManagerImpl_GetAccountQuotas(t *testing.T) {
	tests := []struct {
		name        string
		client      *WrapperClient
		accountName string
		wantErr     bool
		errType     string
	}{
		{
			name:        "nil client",
			client:      nil,
			accountName: "testaccount",
			wantErr:     true,
			errType:     "client error",
		},
		{
			name: "empty account name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "",
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "valid input - not implemented",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "testaccount",
			wantErr:     true,
			errType:     "not implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccountManagerImpl{
				client: tt.client,
			}

			result, err := a.GetAccountQuotas(context.Background(), tt.accountName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetAccountQuotas() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("GetAccountQuotas() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("GetAccountQuotas() expected validation error, got %T: %v", err, err)
					}
				case "not implemented":
					if !errors.IsNotImplementedError(err) {
						t.Errorf("GetAccountQuotas() expected not implemented error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("GetAccountQuotas() unexpected error: %v", err)
				}
			}

			if result != nil && tt.wantErr {
				t.Errorf("GetAccountQuotas() expected nil result on error")
			}
		})
	}
}

func TestAccountManagerImpl_GetAccountQuotaUsage(t *testing.T) {
	tests := []struct {
		name        string
		client      *WrapperClient
		accountName string
		timeframe   string
		wantErr     bool
		errType     string
	}{
		{
			name:        "nil client",
			client:      nil,
			accountName: "testaccount",
			timeframe:   "24h",
			wantErr:     true,
			errType:     "client error",
		},
		{
			name: "empty account name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "",
			timeframe:   "24h",
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "valid input - not implemented",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "testaccount",
			timeframe:   "24h",
			wantErr:     true,
			errType:     "not implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccountManagerImpl{
				client: tt.client,
			}

			result, err := a.GetAccountQuotaUsage(context.Background(), tt.accountName, tt.timeframe)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetAccountQuotaUsage() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("GetAccountQuotaUsage() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("GetAccountQuotaUsage() expected validation error, got %T: %v", err, err)
					}
				case "not implemented":
					if !errors.IsNotImplementedError(err) {
						t.Errorf("GetAccountQuotaUsage() expected not implemented error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("GetAccountQuotaUsage() unexpected error: %v", err)
				}
			}

			if result != nil && tt.wantErr {
				t.Errorf("GetAccountQuotaUsage() expected nil result on error")
			}
		})
	}
}

func TestAccountManagerImpl_GetAccountUsers(t *testing.T) {
	tests := []struct {
		name        string
		client      *WrapperClient
		accountName string
		opts        *interfaces.ListAccountUsersOptions
		wantErr     bool
		errType     string
	}{
		{
			name:        "nil client",
			client:      nil,
			accountName: "testaccount",
			opts:        nil,
			wantErr:     true,
			errType:     "client error",
		},
		{
			name: "empty account name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "",
			opts:        nil,
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "valid input - not implemented",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "testaccount",
			opts: &interfaces.ListAccountUsersOptions{
				ActiveOnly: true,
			},
			wantErr: true,
			errType: "not implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccountManagerImpl{
				client: tt.client,
			}

			result, err := a.GetAccountUsers(context.Background(), tt.accountName, tt.opts)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetAccountUsers() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("GetAccountUsers() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("GetAccountUsers() expected validation error, got %T: %v", err, err)
					}
				case "not implemented":
					if !errors.IsNotImplementedError(err) {
						t.Errorf("GetAccountUsers() expected not implemented error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("GetAccountUsers() unexpected error: %v", err)
				}
			}

			if result != nil && tt.wantErr {
				t.Errorf("GetAccountUsers() expected nil result on error")
			}
		})
	}
}

func TestAccountManagerImpl_ValidateUserAccess(t *testing.T) {
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
			name: "valid input - not implemented",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName:    "testuser",
			accountName: "testaccount",
			wantErr:     true,
			errType:     "not implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccountManagerImpl{
				client: tt.client,
			}

			result, err := a.ValidateUserAccess(context.Background(), tt.userName, tt.accountName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateUserAccess() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("ValidateUserAccess() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("ValidateUserAccess() expected validation error, got %T: %v", err, err)
					}
				case "not implemented":
					if !errors.IsNotImplementedError(err) {
						t.Errorf("ValidateUserAccess() expected not implemented error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("ValidateUserAccess() unexpected error: %v", err)
				}
			}

			if result != nil && tt.wantErr {
				t.Errorf("ValidateUserAccess() expected nil result on error")
			}
		})
	}
}

func TestAccountManagerImpl_GetAccountUsersWithPermissions(t *testing.T) {
	tests := []struct {
		name        string
		client      *WrapperClient
		accountName string
		permissions []string
		wantErr     bool
		errType     string
	}{
		{
			name:        "nil client",
			client:      nil,
			accountName: "testaccount",
			permissions: []string{"read", "write"},
			wantErr:     true,
			errType:     "client error",
		},
		{
			name: "empty account name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "",
			permissions: []string{"read", "write"},
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "empty permissions",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "testaccount",
			permissions: []string{},
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "valid input - not implemented",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "testaccount",
			permissions: []string{"read", "write", "admin"},
			wantErr:     true,
			errType:     "not implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccountManagerImpl{
				client: tt.client,
			}

			result, err := a.GetAccountUsersWithPermissions(context.Background(), tt.accountName, tt.permissions)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetAccountUsersWithPermissions() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("GetAccountUsersWithPermissions() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("GetAccountUsersWithPermissions() expected validation error, got %T: %v", err, err)
					}
				case "not implemented":
					if !errors.IsNotImplementedError(err) {
						t.Errorf("GetAccountUsersWithPermissions() expected not implemented error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("GetAccountUsersWithPermissions() unexpected error: %v", err)
				}
			}

			if result != nil && tt.wantErr {
				t.Errorf("GetAccountUsersWithPermissions() expected nil result on error")
			}
		})
	}
}

func TestAccountManagerImpl_GetAccountFairShare(t *testing.T) {
	tests := []struct {
		name        string
		client      *WrapperClient
		accountName string
		wantErr     bool
		errType     string
	}{
		{
			name:        "nil client",
			client:      nil,
			accountName: "testaccount",
			wantErr:     true,
			errType:     "client error",
		},
		{
			name: "empty account name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "",
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "valid input - not implemented",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "testaccount",
			wantErr:     true,
			errType:     "not implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccountManagerImpl{
				client: tt.client,
			}

			result, err := a.GetAccountFairShare(context.Background(), tt.accountName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetAccountFairShare() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("GetAccountFairShare() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("GetAccountFairShare() expected validation error, got %T: %v", err, err)
					}
				case "not implemented":
					if !errors.IsNotImplementedError(err) {
						t.Errorf("GetAccountFairShare() expected not implemented error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("GetAccountFairShare() unexpected error: %v", err)
				}
			}

			if result != nil && tt.wantErr {
				t.Errorf("GetAccountFairShare() expected nil result on error")
			}
		})
	}
}

func TestAccountManagerImpl_GetFairShareHierarchy(t *testing.T) {
	tests := []struct {
		name        string
		client      *WrapperClient
		rootAccount string
		wantErr     bool
		errType     string
	}{
		{
			name:        "nil client",
			client:      nil,
			rootAccount: "root",
			wantErr:     true,
			errType:     "client error",
		},
		{
			name: "empty root account",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			rootAccount: "",
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "valid input - not implemented",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			rootAccount: "root",
			wantErr:     true,
			errType:     "not implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccountManagerImpl{
				client: tt.client,
			}

			result, err := a.GetFairShareHierarchy(context.Background(), tt.rootAccount)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetFairShareHierarchy() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("GetFairShareHierarchy() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("GetFairShareHierarchy() expected validation error, got %T: %v", err, err)
					}
				case "not implemented":
					if !errors.IsNotImplementedError(err) {
						t.Errorf("GetFairShareHierarchy() expected not implemented error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("GetFairShareHierarchy() unexpected error: %v", err)
				}
			}

			if result != nil && tt.wantErr {
				t.Errorf("GetFairShareHierarchy() expected nil result on error")
			}
		})
	}
}
