// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"context"
	"net/http"
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/mock"
)

// MockClientWithResponses is a mock for the API client
type MockAccountClientWithResponses struct {
	mock.Mock
}

func (m *MockAccountClientWithResponses) SlurmdbV0040GetAccountsWithResponse(ctx context.Context, params *api.SlurmdbV0040GetAccountsParams, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0040GetAccountsResponse, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*api.SlurmdbV0040GetAccountsResponse), args.Error(1)
}

func (m *MockAccountClientWithResponses) SlurmdbV0040GetAccountWithResponse(ctx context.Context, accountName string, params *api.SlurmdbV0040GetAccountParams, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0040GetAccountResponse, error) {
	args := m.Called(ctx, accountName, params)
	return args.Get(0).(*api.SlurmdbV0040GetAccountResponse), args.Error(1)
}

func (m *MockAccountClientWithResponses) SlurmdbV0040PostAccountsWithResponse(ctx context.Context, body api.SlurmdbV0040PostAccountsJSONRequestBody, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0040PostAccountsResponse, error) {
	args := m.Called(ctx, body)
	return args.Get(0).(*api.SlurmdbV0040PostAccountsResponse), args.Error(1)
}

func (m *MockAccountClientWithResponses) SlurmdbV0040DeleteAccountWithResponse(ctx context.Context, accountName string, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0040DeleteAccountResponse, error) {
	args := m.Called(ctx, accountName)
	return args.Get(0).(*api.SlurmdbV0040DeleteAccountResponse), args.Error(1)
}

func TestAccountAdapter_ValidateAccountCreate(t *testing.T) {
	adapter := &AccountAdapter{
		BaseManager: base.NewBaseManager("v0.0.40", "Account"),
	}

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
			errMsg:  "account creation data is required",
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
			name: "valid account minimal",
			account: &types.AccountCreate{
				Name: "test-account",
			},
			wantErr: false,
		},
		{
			name: "valid account with description",
			account: &types.AccountCreate{
				Name:        "test-account",
				Description: "Test account for unit tests",
			},
			wantErr: false,
		},
		{
			name: "valid account with organization",
			account: &types.AccountCreate{
				Name:         "test-account",
				Description:  "Test account for unit tests",
				Organization: "Test Organization",
			},
			wantErr: false,
		},
		{
			name: "valid account with coordinators",
			account: &types.AccountCreate{
				Name:         "test-account",
				Description:  "Test account for unit tests",
				Organization: "Test Organization",
				Coordinators: []string{"coord1", "coord2"},
			},
			wantErr: false,
		},
		{
			name: "valid account with resource limits",
			account: &types.AccountCreate{
				Name:        "test-account",
				Description: "Test account for unit tests",
				MaxJobs:     100,
				MaxCPUs:     1000,
				MaxNodes:    10,
				FairShare:   50,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.validateAccountCreate(tt.account)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAccountAdapter_ValidateAccountUpdate(t *testing.T) {
	adapter := &AccountAdapter{
		BaseManager: base.NewBaseManager("v0.0.40", "Account"),
	}

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
			name:    "empty update",
			update:  &types.AccountUpdate{},
			wantErr: true,
			errMsg:  "at least one field must be provided for update",
		},
		{
			name: "valid update with description",
			update: &types.AccountUpdate{
				Description: stringPtr("Updated description"),
			},
			wantErr: false,
		},
		{
			name: "valid update with organization",
			update: &types.AccountUpdate{
				Organization: stringPtr("Updated organization"),
			},
			wantErr: false,
		},
		{
			name: "valid update with both fields",
			update: &types.AccountUpdate{
				Description:  stringPtr("Updated description"),
				Organization: stringPtr("Updated organization"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.validateAccountUpdate(tt.update)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAccountAdapter_List(t *testing.T) {
	tests := []struct {
		name           string
		opts           *types.AccountListOptions
		mockResponse   *api.SlurmdbV0040GetAccountsResponse
		mockError      error
		expectedLen    int
		expectedError  string
		setupMock      func(*MockAccountClientWithResponses)
	}{
		{
			name: "successful list with no options",
			opts: nil,
			mockResponse: &api.SlurmdbV0040GetAccountsResponse{
				JSON200: &api.V0040OpenapiAccountsResp{
					Accounts: []api.V0040Account{
						{Name: stringPtr("account1")},
						{Name: stringPtr("account2")},
					},
				},
			},
			expectedLen: 2,
			setupMock: func(m *MockAccountClientWithResponses) {
				m.On("SlurmdbV0040GetAccountsWithResponse", mock.Anything, mock.AnythingOfType("*v0_0_40.SlurmdbV0040GetAccountsParams")).
					Return(&api.SlurmdbV0040GetAccountsResponse{
						JSON200: &api.V0040OpenapiAccountsResp{
							Accounts: []api.V0040Account{
								{Name: stringPtr("account1")},
								{Name: stringPtr("account2")},
							},
						},
					}, nil)
			},
		},
		{
			name: "successful list with options",
			opts: &types.AccountListOptions{
				Descriptions: []string{"test"},
				WithDeleted:  true,
			},
			expectedLen: 1,
			setupMock: func(m *MockAccountClientWithResponses) {
				m.On("SlurmdbV0040GetAccountsWithResponse", mock.Anything, mock.AnythingOfType("*v0_0_40.SlurmdbV0040GetAccountsParams")).
					Return(&api.SlurmdbV0040GetAccountsResponse{
						JSON200: &api.V0040OpenapiAccountsResp{
							Accounts: []api.V0040Account{
								{Name: stringPtr("test-account")},
							},
						},
					}, nil)
			},
		},
		{
			name:          "API error",
			opts:          nil,
			expectedError: "API error",
			setupMock: func(m *MockAccountClientWithResponses) {
				m.On("SlurmdbV0040GetAccountsWithResponse", mock.Anything, mock.AnythingOfType("*v0_0_40.SlurmdbV0040GetAccountsParams")).
					Return((*api.SlurmdbV0040GetAccountsResponse)(nil), assert.AnError)
			},
		},
		{
			name: "empty response",
			opts: nil,
			setupMock: func(m *MockAccountClientWithResponses) {
				m.On("SlurmdbV0040GetAccountsWithResponse", mock.Anything, mock.AnythingOfType("*v0_0_40.SlurmdbV0040GetAccountsParams")).
					Return(&api.SlurmdbV0040GetAccountsResponse{
						JSON200: &api.V0040OpenapiAccountsResp{
							Accounts: []api.V0040Account{},
						},
					}, nil)
			},
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAccountClientWithResponses{}
			tt.setupMock(mockClient)

			// Create adapter with mock client
			adapter := &AccountAdapter{
				BaseManager: base.NewBaseManager("v0.0.40", "Account"),
				client:      mockClient,
			}

			result, err := adapter.List(context.Background(), tt.opts)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Len(t, result.Accounts, tt.expectedLen)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestAccountAdapter_Get(t *testing.T) {
	tests := []struct {
		name          string
		accountName   string
		mockResponse  *api.SlurmdbV0040GetAccountResponse
		mockError     error
		expectedError string
		setupMock     func(*MockAccountClientWithResponses, string)
	}{
		{
			name:        "successful get",
			accountName: "test-account",
			setupMock: func(m *MockAccountClientWithResponses, name string) {
				m.On("SlurmdbV0040GetAccountWithResponse", mock.Anything, name, mock.AnythingOfType("*v0_0_40.SlurmdbV0040GetAccountParams")).
					Return(&api.SlurmdbV0040GetAccountResponse{
						JSON200: &api.V0040OpenapiAccountsResp{
							Accounts: []api.V0040Account{
								{Name: stringPtr(name)},
							},
						},
					}, nil)
			},
		},
		{
			name:        "empty account name",
			accountName: "",
			setupMock: func(m *MockAccountClientWithResponses, name string) {
				// No mock setup needed as validation should fail first
			},
			expectedError: "accountName",
		},
		{
			name:        "API error",
			accountName: "test-account",
			setupMock: func(m *MockAccountClientWithResponses, name string) {
				m.On("SlurmdbV0040GetAccountWithResponse", mock.Anything, name, mock.AnythingOfType("*v0_0_40.SlurmdbV0040GetAccountParams")).
					Return((*api.SlurmdbV0040GetAccountResponse)(nil), assert.AnError)
			},
			expectedError: "API error",
		},
		{
			name:        "account not found",
			accountName: "nonexistent",
			setupMock: func(m *MockAccountClientWithResponses, name string) {
				m.On("SlurmdbV0040GetAccountWithResponse", mock.Anything, name, mock.AnythingOfType("*v0_0_40.SlurmdbV0040GetAccountParams")).
					Return(&api.SlurmdbV0040GetAccountResponse{
						JSON200: &api.V0040OpenapiAccountsResp{
							Accounts: []api.V0040Account{},
						},
					}, nil)
			},
			expectedError: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAccountClientWithResponses{}
			tt.setupMock(mockClient, tt.accountName)

			adapter := &AccountAdapter{
				BaseManager: base.NewBaseManager("v0.0.40", "Account"),
				client:      mockClient,
			}

			result, err := adapter.Get(context.Background(), tt.accountName)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.accountName, result.Name)
			}

			if tt.accountName != "" && tt.expectedError == "" {
				mockClient.AssertExpectations(t)
			}
		})
	}
}

func TestAccountAdapter_Create(t *testing.T) {
	tests := []struct {
		name          string
		account       *types.AccountCreate
		mockResponse  *api.SlurmdbV0040PostAccountsResponse
		mockError     error
		expectedError string
		setupMock     func(*MockAccountClientWithResponses, *types.AccountCreate)
	}{
		{
			name: "successful create",
			account: &types.AccountCreate{
				Name:        "test-account",
				Description: "Test account",
			},
			setupMock: func(m *MockAccountClientWithResponses, account *types.AccountCreate) {
				m.On("SlurmdbV0040PostAccountsWithResponse", mock.Anything, mock.AnythingOfType("v0_0_40.SlurmdbV0040PostAccountsJSONRequestBody")).
					Return(&api.SlurmdbV0040PostAccountsResponse{
						JSON200: &api.V0040OpenapiResp{},
					}, nil)
			},
		},
		{
			name:    "nil account",
			account: nil,
			setupMock: func(m *MockAccountClientWithResponses, account *types.AccountCreate) {
				// No mock setup needed as validation should fail first
			},
			expectedError: "account creation data is required",
		},
		{
			name: "empty name",
			account: &types.AccountCreate{
				Name: "",
			},
			setupMock: func(m *MockAccountClientWithResponses, account *types.AccountCreate) {
				// No mock setup needed as validation should fail first
			},
			expectedError: "account name is required",
		},
		{
			name: "API error",
			account: &types.AccountCreate{
				Name: "test-account",
			},
			setupMock: func(m *MockAccountClientWithResponses, account *types.AccountCreate) {
				m.On("SlurmdbV0040PostAccountsWithResponse", mock.Anything, mock.AnythingOfType("v0_0_40.SlurmdbV0040PostAccountsJSONRequestBody")).
					Return((*api.SlurmdbV0040PostAccountsResponse)(nil), assert.AnError)
			},
			expectedError: "API error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAccountClientWithResponses{}
			tt.setupMock(mockClient, tt.account)

			adapter := &AccountAdapter{
				BaseManager: base.NewBaseManager("v0.0.40", "Account"),
				client:      mockClient,
			}

			err := adapter.Create(context.Background(), tt.account)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			if tt.account != nil && tt.account.Name != "" && tt.expectedError == "" {
				mockClient.AssertExpectations(t)
			}
		})
	}
}

func TestAccountAdapter_Update(t *testing.T) {
	tests := []struct {
		name          string
		accountName   string
		update        *types.AccountUpdate
		mockGetResp   *api.SlurmdbV0040GetAccountResponse
		mockPostResp  *api.SlurmdbV0040PostAccountsResponse
		mockError     error
		expectedError string
		setupMock     func(*MockAccountClientWithResponses, string, *types.AccountUpdate)
	}{
		{
			name:        "successful update",
			accountName: "test-account",
			update: &types.AccountUpdate{
				Description: stringPtr("Updated description"),
			},
			setupMock: func(m *MockAccountClientWithResponses, name string, update *types.AccountUpdate) {
				// Mock Get call first
				m.On("SlurmdbV0040GetAccountWithResponse", mock.Anything, name, mock.AnythingOfType("*v0_0_40.SlurmdbV0040GetAccountParams")).
					Return(&api.SlurmdbV0040GetAccountResponse{
						JSON200: &api.V0040OpenapiAccountsResp{
							Accounts: []api.V0040Account{
								{Name: stringPtr(name)},
							},
						},
					}, nil)
				// Mock Post call for update
				m.On("SlurmdbV0040PostAccountsWithResponse", mock.Anything, mock.AnythingOfType("v0_0_40.SlurmdbV0040PostAccountsJSONRequestBody")).
					Return(&api.SlurmdbV0040PostAccountsResponse{
						JSON200: &api.V0040OpenapiResp{},
					}, nil)
			},
		},
		{
			name:        "empty account name",
			accountName: "",
			update: &types.AccountUpdate{
				Description: stringPtr("Updated description"),
			},
			setupMock: func(m *MockAccountClientWithResponses, name string, update *types.AccountUpdate) {
				// No mock setup needed as validation should fail first
			},
			expectedError: "accountName",
		},
		{
			name:        "nil update",
			accountName: "test-account",
			update:      nil,
			setupMock: func(m *MockAccountClientWithResponses, name string, update *types.AccountUpdate) {
				// No mock setup needed as validation should fail first
			},
			expectedError: "account update data is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAccountClientWithResponses{}
			tt.setupMock(mockClient, tt.accountName, tt.update)

			adapter := &AccountAdapter{
				BaseManager: base.NewBaseManager("v0.0.40", "Account"),
				client:      mockClient,
			}

			err := adapter.Update(context.Background(), tt.accountName, tt.update)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			if tt.accountName != "" && tt.update != nil && tt.expectedError == "" {
				mockClient.AssertExpectations(t)
			}
		})
	}
}

func TestAccountAdapter_Delete(t *testing.T) {
	tests := []struct {
		name          string
		accountName   string
		mockResponse  *api.SlurmdbV0040DeleteAccountResponse
		mockError     error
		expectedError string
		setupMock     func(*MockAccountClientWithResponses, string)
	}{
		{
			name:        "successful delete",
			accountName: "test-account",
			setupMock: func(m *MockAccountClientWithResponses, name string) {
				m.On("SlurmdbV0040DeleteAccountWithResponse", mock.Anything, name).
					Return(&api.SlurmdbV0040DeleteAccountResponse{
						JSON200: &api.V0040OpenapiResp{},
					}, nil)
			},
		},
		{
			name:        "successful delete with 204 status",
			accountName: "test-account",
			setupMock: func(m *MockAccountClientWithResponses, name string) {
				resp := &api.SlurmdbV0040DeleteAccountResponse{}
				// Mock StatusCode method
				resp.HTTPResponse = &http.Response{StatusCode: 204}
				m.On("SlurmdbV0040DeleteAccountWithResponse", mock.Anything, name).
					Return(resp, nil)
			},
		},
		{
			name:        "empty account name",
			accountName: "",
			setupMock: func(m *MockAccountClientWithResponses, name string) {
				// No mock setup needed as validation should fail first
			},
			expectedError: "accountName",
		},
		{
			name:        "API error",
			accountName: "test-account",
			setupMock: func(m *MockAccountClientWithResponses, name string) {
				m.On("SlurmdbV0040DeleteAccountWithResponse", mock.Anything, name).
					Return((*api.SlurmdbV0040DeleteAccountResponse)(nil), assert.AnError)
			},
			expectedError: "API error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAccountClientWithResponses{}
			tt.setupMock(mockClient, tt.accountName)

			adapter := &AccountAdapter{
				BaseManager: base.NewBaseManager("v0.0.40", "Account"),
				client:      mockClient,
			}

			err := adapter.Delete(context.Background(), tt.accountName)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			if tt.accountName != "" && tt.expectedError == "" {
				mockClient.AssertExpectations(t)
			}
		})
	}
}

// Test error conditions and edge cases
func TestAccountAdapter_ErrorConditions(t *testing.T) {
	t.Run("nil context", func(t *testing.T) {
		adapter := &AccountAdapter{
			BaseManager: base.NewBaseManager("v0.0.40", "Account"),
		}

		_, err := adapter.List(nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "context")
	})

	t.Run("nil client", func(t *testing.T) {
		adapter := &AccountAdapter{
			BaseManager: base.NewBaseManager("v0.0.40", "Account"),
			client:      nil,
		}

		_, err := adapter.List(context.Background(), nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "client")
	})
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func int32Ptr(i int32) *int32 {
	return &i
}
