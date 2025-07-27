package v0_0_43

import (
	"context"
	"testing"

	"github.com/jontk/slurm-client/pkg/errors"
)

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
			name: "nil api client",
			client: &WrapperClient{
				apiClient: nil,
			},
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
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			rootAccount: "root",
			wantErr:     true,
			errType:     "client error",
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
			accountName: "test",
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
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "test",
			wantErr:     true,
			errType:     "client error",
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
			accountName: "test",
			depth:       1,
			wantErr:     true,
			errType:     "client error",
		},
		{
			name: "empty account name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "",
			depth:       1,
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "negative depth",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "test",
			depth:       -1,
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "test",
			depth:       1,
			wantErr:     true,
			errType:     "client error",
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
			accountName: "test",
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
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "test",
			wantErr:     true,
			errType:     "client error",
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
			accountName: "test",
			timeframe:   "current",
			wantErr:     true,
			errType:     "client error",
		},
		{
			name: "empty account name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "",
			timeframe:   "current",
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "invalid timeframe",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "test",
			timeframe:   "invalid",
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "test",
			timeframe:   "current",
			wantErr:     true,
			errType:     "client error",
		},
		{
			name: "empty timeframe defaults to current",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "test",
			timeframe:   "",
			wantErr:     true,
			errType:     "client error",
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

func TestValidateTRES(t *testing.T) {
	tests := []struct {
		name    string
		tres    map[string]int
		wantErr bool
	}{
		{
			name: "valid TRES",
			tres: map[string]int{
				"cpu": 100,
				"mem": 1024,
				"node": 10,
			},
			wantErr: false,
		},
		{
			name: "negative value",
			tres: map[string]int{
				"cpu": 100,
				"mem": -1,
			},
			wantErr: true,
		},
		{
			name:    "nil TRES",
			tres:    nil,
			wantErr: false,
		},
		{
			name:    "empty TRES",
			tres:    map[string]int{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTRES(tt.tres)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("validateTRES() expected error but got none")
				}
				if !errors.IsValidationError(err) {
					t.Errorf("validateTRES() expected validation error, got %T: %v", err, err)
				}
			} else {
				if err != nil {
					t.Errorf("validateTRES() unexpected error: %v", err)
				}
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
			name: "nil api client",
			client: &WrapperClient{
				apiClient: nil,
			},
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
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "testaccount",
			wantErr:     true,
			errType:     "client error",
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
			name: "nil api client",
			client: &WrapperClient{
				apiClient: nil,
			},
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
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			rootAccount: "root",
			wantErr:     true,
			errType:     "client error",
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
