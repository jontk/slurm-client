package v0_0_41

import (
	"context"
	"net/http"
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/mock"
)

// MockClientWithResponses is a mock for the API client
type MockAccountClientWithResponses struct {
	mock.Mock
}

func (m *MockAccountClientWithResponses) SlurmdbV0041GetAccountsWithResponse(ctx context.Context, params *api.SlurmdbV0041GetAccountsParams, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0041GetAccountsResponse, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*api.SlurmdbV0041GetAccountsResponse), args.Error(1)
}

func (m *MockAccountClientWithResponses) SlurmdbV0041GetAccountWithResponse(ctx context.Context, accountName string, params *api.SlurmdbV0041GetAccountParams, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0041GetAccountResponse, error) {
	args := m.Called(ctx, accountName, params)
	return args.Get(0).(*api.SlurmdbV0041GetAccountResponse), args.Error(1)
}

func (m *MockAccountClientWithResponses) SlurmdbV0041PostAccountsWithResponse(ctx context.Context, body api.SlurmdbV0041PostAccountsJSONRequestBody, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0041PostAccountsResponse, error) {
	args := m.Called(ctx, body)
	return args.Get(0).(*api.SlurmdbV0041PostAccountsResponse), args.Error(1)
}

func (m *MockAccountClientWithResponses) SlurmdbV0041DeleteAccountWithResponse(ctx context.Context, accountName string, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0041DeleteAccountResponse, error) {
	args := m.Called(ctx, accountName)
	return args.Get(0).(*api.SlurmdbV0041DeleteAccountResponse), args.Error(1)
}

func (m *MockAccountClientWithResponses) SlurmdbV0041PostAccountsAssociationWithResponse(ctx context.Context, body api.SlurmdbV0041PostAccountsAssociationJSONBody, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0041PostAccountsAssociationResponse, error) {
	args := m.Called(ctx, body)
	return args.Get(0).(*api.SlurmdbV0041PostAccountsAssociationResponse), args.Error(1)
}

// Helper method for HTTPResponse mock
func (m *MockAccountClientWithResponses) HandleHTTPResponse(resp *http.Response, body []byte) error {
	args := m.Called(resp, body)
	return args.Error(0)
}

func TestAccountAdapter_List(t *testing.T) {
	tests := []struct {
		name           string
		opts           *types.AccountListOptions
		mockResponse   *api.SlurmdbV0041GetAccountsResponse
		mockError      error
		expectedLen    int
		expectedError  string
		setupMock      func(*MockAccountClientWithResponses)
	}{
		{
			name: "successful list with no options",
			opts: nil,
			setupMock: func(m *MockAccountClientWithResponses) {
				m.On("SlurmdbV0041GetAccountsWithResponse", mock.Anything, mock.AnythingOfType("*v0_0_41.SlurmdbV0041GetAccountsParams")).
					Return(&api.SlurmdbV0041GetAccountsResponse{
						JSON200: &api.V0041OpenapiAccountsResp{
							Accounts: []api.V0041Account{
								{Name: stringPtr("account1")},
								{Name: stringPtr("account2")},
							},
						},
						HTTPResponse: &http.Response{StatusCode: 200},
					}, nil)
				m.On("HandleHTTPResponse", mock.Anything, mock.Anything).Return(nil)
			},
			expectedLen: 2,
		},
		{
			name: "successful list with options",
			opts: &types.AccountListOptions{
				Names:        []string{"test-account"},
				Description:  "test description",
				Organization: "test org",
				WithDeleted:  true,
			},
			setupMock: func(m *MockAccountClientWithResponses) {
				m.On("SlurmdbV0041GetAccountsWithResponse", mock.Anything, mock.AnythingOfType("*v0_0_41.SlurmdbV0041GetAccountsParams")).
					Return(&api.SlurmdbV0041GetAccountsResponse{
						JSON200: &api.V0041OpenapiAccountsResp{
							Accounts: []api.V0041Account{
								{Name: stringPtr("test-account")},
							},
						},
						HTTPResponse: &http.Response{StatusCode: 200},
					}, nil)
				m.On("HandleHTTPResponse", mock.Anything, mock.Anything).Return(nil)
			},
			expectedLen: 1,
		},
		{
			name:          "API error",
			opts:          nil,
			expectedError: "failed to list accounts",
			setupMock: func(m *MockAccountClientWithResponses) {
				m.On("SlurmdbV0041GetAccountsWithResponse", mock.Anything, mock.AnythingOfType("*v0_0_41.SlurmdbV0041GetAccountsParams")).
					Return((*api.SlurmdbV0041GetAccountsResponse)(nil), assert.AnError)
			},
		},
		{
			name: "empty response",
			opts: nil,
			setupMock: func(m *MockAccountClientWithResponses) {
				m.On("SlurmdbV0041GetAccountsWithResponse", mock.Anything, mock.AnythingOfType("*v0_0_41.SlurmdbV0041GetAccountsParams")).
					Return(&api.SlurmdbV0041GetAccountsResponse{
						JSON200: &api.V0041OpenapiAccountsResp{
							Accounts: []api.V0041Account{},
						},
						HTTPResponse: &http.Response{StatusCode: 200},
					}, nil)
				m.On("HandleHTTPResponse", mock.Anything, mock.Anything).Return(nil)
			},
			expectedLen: 0,
		},
		{
			name: "response with warnings and errors",
			opts: nil,
			setupMock: func(m *MockAccountClientWithResponses) {
				warnings := []api.V0041OpenapiWarning{
					{Description: stringPtr("Test warning")},
				}
				errors := []api.V0041OpenapiError{
					{Description: stringPtr("Test error")},
				}
				m.On("SlurmdbV0041GetAccountsWithResponse", mock.Anything, mock.AnythingOfType("*v0_0_41.SlurmdbV0041GetAccountsParams")).
					Return(&api.SlurmdbV0041GetAccountsResponse{
						JSON200: &api.V0041OpenapiAccountsResp{
							Accounts: []api.V0041Account{
								{Name: stringPtr("account1")},
							},
							Warnings: &warnings,
							Errors:   &errors,
						},
						HTTPResponse: &http.Response{StatusCode: 200},
					}, nil)
				m.On("HandleHTTPResponse", mock.Anything, mock.Anything).Return(nil)
			},
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAccountClientWithResponses{}
			tt.setupMock(mockClient)

			// Create adapter with mock client
			adapter := &AccountAdapter{
				BaseManager: base.NewBaseManager("v0.0.41", "Account"),
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
				
				// Check that metadata is populated
				assert.NotNil(t, result.Meta)
				assert.Equal(t, "v0.0.41", result.Meta.Version)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestAccountAdapter_Get(t *testing.T) {
	tests := []struct {
		name          string
		accountName   string
		mockResponse  *api.SlurmdbV0041GetAccountResponse
		mockError     error
		expectedError string
		setupMock     func(*MockAccountClientWithResponses, string)
	}{
		{
			name:        "successful get",
			accountName: "test-account",
			setupMock: func(m *MockAccountClientWithResponses, name string) {
				m.On("SlurmdbV0041GetAccountWithResponse", mock.Anything, name, mock.AnythingOfType("*v0_0_41.SlurmdbV0041GetAccountParams")).
					Return(&api.SlurmdbV0041GetAccountResponse{
						JSON200: &api.V0041OpenapiAccountsResp{
							Accounts: []api.V0041Account{
								{Name: stringPtr(name)},
							},
						},
						HTTPResponse: &http.Response{StatusCode: 200},
					}, nil)
				m.On("HandleHTTPResponse", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:        "empty account name",
			accountName: "",
			setupMock: func(m *MockAccountClientWithResponses, name string) {
				// No mock setup needed as validation should fail first
			},
			expectedError: "account name",
		},
		{
			name:        "API error",
			accountName: "test-account",
			setupMock: func(m *MockAccountClientWithResponses, name string) {
				m.On("SlurmdbV0041GetAccountWithResponse", mock.Anything, name, mock.AnythingOfType("*v0_0_41.SlurmdbV0041GetAccountParams")).
					Return((*api.SlurmdbV0041GetAccountResponse)(nil), assert.AnError)
			},
			expectedError: "failed to get account",
		},
		{
			name:        "account not found",
			accountName: "nonexistent",
			setupMock: func(m *MockAccountClientWithResponses, name string) {
				m.On("SlurmdbV0041GetAccountWithResponse", mock.Anything, name, mock.AnythingOfType("*v0_0_41.SlurmdbV0041GetAccountParams")).
					Return(&api.SlurmdbV0041GetAccountResponse{
						JSON200: &api.V0041OpenapiAccountsResp{
							Accounts: []api.V0041Account{},
						},
						HTTPResponse: &http.Response{StatusCode: 200},
					}, nil)
				m.On("HandleHTTPResponse", mock.Anything, mock.Anything).Return(nil)
				m.On("HandleNotFound", mock.Anything).Return(assert.AnError)
			},
			expectedError: "account nonexistent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAccountClientWithResponses{}
			tt.setupMock(mockClient, tt.accountName)

			adapter := &AccountAdapter{
				BaseManager: base.NewBaseManager("v0.0.41", "Account"),
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
		account       *types.Account
		mockResponse  *api.SlurmdbV0041PostAccountsResponse
		mockError     error
		expectedError string
		setupMock     func(*MockAccountClientWithResponses, *types.Account)
	}{
		{
			name: "successful create",
			account: &types.Account{
				Name:        "test-account",
				Description: "Test account",
			},
			setupMock: func(m *MockAccountClientWithResponses, account *types.Account) {
				m.On("SlurmdbV0041PostAccountsWithResponse", mock.Anything, mock.AnythingOfType("v0_0_41.SlurmdbV0041PostAccountsJSONRequestBody")).
					Return(&api.SlurmdbV0041PostAccountsResponse{
						HTTPResponse: &http.Response{StatusCode: 200},
					}, nil)
				m.On("HandleHTTPResponse", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:    "nil account",
			account: nil,
			setupMock: func(m *MockAccountClientWithResponses, account *types.Account) {
				// No mock setup needed as validation should fail first
			},
			expectedError: "account cannot be nil",
		},
		{
			name: "empty name",
			account: &types.Account{
				Name: "",
			},
			setupMock: func(m *MockAccountClientWithResponses, account *types.Account) {
				// No mock setup needed as validation should fail first
			},
			expectedError: "account name",
		},
		{
			name: "API error",
			account: &types.Account{
				Name: "test-account",
			},
			setupMock: func(m *MockAccountClientWithResponses, account *types.Account) {
				m.On("SlurmdbV0041PostAccountsWithResponse", mock.Anything, mock.AnythingOfType("v0_0_41.SlurmdbV0041PostAccountsJSONRequestBody")).
					Return((*api.SlurmdbV0041PostAccountsResponse)(nil), assert.AnError)
			},
			expectedError: "failed to create account",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAccountClientWithResponses{}
			tt.setupMock(mockClient, tt.account)

			adapter := &AccountAdapter{
				BaseManager: base.NewBaseManager("v0.0.41", "Account"),
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
		mockGetResp   *api.SlurmdbV0041GetAccountResponse
		mockPostResp  *api.SlurmdbV0041PostAccountsResponse
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
				m.On("SlurmdbV0041GetAccountWithResponse", mock.Anything, name, mock.AnythingOfType("*v0_0_41.SlurmdbV0041GetAccountParams")).
					Return(&api.SlurmdbV0041GetAccountResponse{
						JSON200: &api.V0041OpenapiAccountsResp{
							Accounts: []api.V0041Account{
								{Name: stringPtr(name)},
							},
						},
						HTTPResponse: &http.Response{StatusCode: 200},
					}, nil)
				m.On("HandleHTTPResponse", mock.Anything, mock.Anything).Return(nil)
				// Mock Post call for update
				m.On("SlurmdbV0041PostAccountsWithResponse", mock.Anything, mock.AnythingOfType("v0_0_41.SlurmdbV0041PostAccountsJSONRequestBody")).
					Return(&api.SlurmdbV0041PostAccountsResponse{
						HTTPResponse: &http.Response{StatusCode: 200},
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
			expectedError: "account name",
		},
		{
			name:        "nil update",
			accountName: "test-account",
			update:      nil,
			setupMock: func(m *MockAccountClientWithResponses, name string, update *types.AccountUpdate) {
				// No mock setup needed as validation should fail first
			},
			expectedError: "account update cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAccountClientWithResponses{}
			tt.setupMock(mockClient, tt.accountName, tt.update)

			adapter := &AccountAdapter{
				BaseManager: base.NewBaseManager("v0.0.41", "Account"),
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
		mockResponse  *api.SlurmdbV0041DeleteAccountResponse
		mockError     error
		expectedError string
		setupMock     func(*MockAccountClientWithResponses, string)
	}{
		{
			name:        "successful delete",
			accountName: "test-account",
			setupMock: func(m *MockAccountClientWithResponses, name string) {
				m.On("SlurmdbV0041DeleteAccountWithResponse", mock.Anything, name).
					Return(&api.SlurmdbV0041DeleteAccountResponse{
						HTTPResponse: &http.Response{StatusCode: 200},
					}, nil)
				m.On("HandleHTTPResponse", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:        "empty account name",
			accountName: "",
			setupMock: func(m *MockAccountClientWithResponses, name string) {
				// No mock setup needed as validation should fail first
			},
			expectedError: "account name",
		},
		{
			name:        "API error",
			accountName: "test-account",
			setupMock: func(m *MockAccountClientWithResponses, name string) {
				m.On("SlurmdbV0041DeleteAccountWithResponse", mock.Anything, name).
					Return((*api.SlurmdbV0041DeleteAccountResponse)(nil), assert.AnError)
			},
			expectedError: "failed to delete account",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAccountClientWithResponses{}
			tt.setupMock(mockClient, tt.accountName)

			adapter := &AccountAdapter{
				BaseManager: base.NewBaseManager("v0.0.41", "Account"),
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

func TestAccountAdapter_AddUser(t *testing.T) {
	tests := []struct {
		name          string
		accountName   string
		userName      string
		opts          *types.AccountUserOptions
		mockResponse  *api.SlurmdbV0041PostAccountsAssociationResponse
		mockError     error
		expectedError string
		setupMock     func(*MockAccountClientWithResponses, string, string, *types.AccountUserOptions)
	}{
		{
			name:        "successful add user",
			accountName: "test-account",
			userName:    "test-user",
			opts:        nil,
			setupMock: func(m *MockAccountClientWithResponses, accName, userName string, opts *types.AccountUserOptions) {
				m.On("SlurmdbV0041PostAccountsAssociationWithResponse", mock.Anything, mock.AnythingOfType("v0_0_41.SlurmdbV0041PostAccountsAssociationJSONBody")).
					Return(&api.SlurmdbV0041PostAccountsAssociationResponse{
						HTTPResponse: &http.Response{StatusCode: 200},
					}, nil)
				m.On("HandleHTTPResponse", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:        "add user with options",
			accountName: "test-account",
			userName:    "test-user",
			opts: &types.AccountUserOptions{
				Cluster:    "test-cluster",
				Partition:  "test-partition",
				DefaultQoS: "test-qos",
			},
			setupMock: func(m *MockAccountClientWithResponses, accName, userName string, opts *types.AccountUserOptions) {
				m.On("SlurmdbV0041PostAccountsAssociationWithResponse", mock.Anything, mock.AnythingOfType("v0_0_41.SlurmdbV0041PostAccountsAssociationJSONBody")).
					Return(&api.SlurmdbV0041PostAccountsAssociationResponse{
						HTTPResponse: &http.Response{StatusCode: 200},
					}, nil)
				m.On("HandleHTTPResponse", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:        "empty account name",
			accountName: "",
			userName:    "test-user",
			setupMock: func(m *MockAccountClientWithResponses, accName, userName string, opts *types.AccountUserOptions) {
				// No mock setup needed as validation should fail first
			},
			expectedError: "account name",
		},
		{
			name:        "empty user name",
			accountName: "test-account",
			userName:    "",
			setupMock: func(m *MockAccountClientWithResponses, accName, userName string, opts *types.AccountUserOptions) {
				// No mock setup needed as validation should fail first
			},
			expectedError: "user name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAccountClientWithResponses{}
			tt.setupMock(mockClient, tt.accountName, tt.userName, tt.opts)

			adapter := &AccountAdapter{
				BaseManager: base.NewBaseManager("v0.0.41", "Account"),
				client:      mockClient,
			}

			err := adapter.AddUser(context.Background(), tt.accountName, tt.userName, tt.opts)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			if tt.accountName != "" && tt.userName != "" && tt.expectedError == "" {
				mockClient.AssertExpectations(t)
			}
		})
	}
}

func TestAccountAdapter_GetAssociations(t *testing.T) {
	adapter := &AccountAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Account"),
	}

	_, err := adapter.GetAssociations(context.Background(), "test-account")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not directly supported")
}

func TestAccountAdapter_RemoveUser(t *testing.T) {
	adapter := &AccountAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Account"),
	}

	err := adapter.RemoveUser(context.Background(), "test-account", "test-user")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not directly supported")
}

func TestAccountAdapter_SetCoordinators(t *testing.T) {
	adapter := &AccountAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Account"),
	}

	err := adapter.SetCoordinators(context.Background(), "test-account", []string{"coord1"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not supported")
}

// Test error conditions and edge cases
func TestAccountAdapter_ErrorConditions(t *testing.T) {
	t.Run("nil context", func(t *testing.T) {
		adapter := &AccountAdapter{
			BaseManager: base.NewBaseManager("v0.0.41", "Account"),
		}

		_, err := adapter.List(nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "context")
	})

	t.Run("nil client", func(t *testing.T) {
		adapter := &AccountAdapter{
			BaseManager: base.NewBaseManager("v0.0.41", "Account"),
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