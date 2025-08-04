// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// Mock client for testing
type MockAssociationClient struct {
	mock.Mock
}

func (m *MockAssociationClient) SlurmdbV0043GetAssociationsWithResponse(ctx context.Context, params *api.SlurmdbV0043GetAssociationsParams, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0043GetAssociationsResponse, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmdbV0043GetAssociationsResponse), args.Error(1)
}

func (m *MockAssociationClient) SlurmdbV0043GetAssociationWithResponse(ctx context.Context, associationId string, params *api.SlurmdbV0043GetAssociationParams, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0043GetAssociationResponse, error) {
	args := m.Called(ctx, associationId, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmdbV0043GetAssociationResponse), args.Error(1)
}

func (m *MockAssociationClient) SlurmdbV0043PostAssociationsWithResponse(ctx context.Context, params *api.SlurmdbV0043PostAssociationsParams, body api.SlurmdbV0043PostAssociationsJSONRequestBody, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0043PostAssociationsResponse, error) {
	args := m.Called(ctx, params, body)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmdbV0043PostAssociationsResponse), args.Error(1)
}

func (m *MockAssociationClient) SlurmdbV0043DeleteAssociationWithResponse(ctx context.Context, associationId string, params *api.SlurmdbV0043DeleteAssociationParams, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0043DeleteAssociationResponse, error) {
	args := m.Called(ctx, associationId, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmdbV0043DeleteAssociationResponse), args.Error(1)
}

func TestNewAssociationAdapter(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewAssociationAdapter(client)
	
	assert.NotNil(t, adapter)
	assert.Equal(t, client, adapter.client)
	assert.NotNil(t, adapter.BaseManager)
	assert.Nil(t, adapter.wrapper)
}

func TestAssociationAdapter_List(t *testing.T) {
	tests := []struct {
		name          string
		opts          *types.AssociationListOptions
		mockResponse  *api.SlurmdbV0043GetAssociationsResponse
		mockError     error
		expectedError bool
		expectedCount int
	}{
		{
			name: "successful list",
			opts: &types.AssociationListOptions{},
			mockResponse: &api.SlurmdbV0043GetAssociationsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiAssocsResp{
					Associations: []api.V0043Assoc{
						{
							Account:   ptrString("account1"),
							User:      ptrString("user1"),
							Cluster:   ptrString("cluster1"),
							Partition: ptrString("normal"),
						},
						{
							Account:   ptrString("account2"),
							User:      ptrString("user2"),
							Cluster:   ptrString("cluster1"),
							Partition: ptrString("gpu"),
						},
					},
				},
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name: "list with filters",
			opts: &types.AssociationListOptions{
				Accounts:   []string{"account1", "account2"},
				Users:      []string{"user1"},
				Clusters:   []string{"cluster1"},
				Partitions: []string{"normal"},
			},
			mockResponse: &api.SlurmdbV0043GetAssociationsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiAssocsResp{
					Associations: []api.V0043Assoc{
						{
							Account: ptrString("account1"),
							User:    ptrString("user1"),
						},
					},
				},
			},
			expectedError: false,
			expectedCount: 1,
		},
		{
			name: "list with deleted",
			opts: &types.AssociationListOptions{
				WithDeleted: true,
			},
			mockResponse: &api.SlurmdbV0043GetAssociationsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiAssocsResp{
					Associations: []api.V0043Assoc{},
				},
			},
			expectedError: false,
			expectedCount: 0,
		},
		{
			name:          "nil context",
			opts:          nil,
			expectedError: true,
		},
		{
			name: "API error",
			opts: &types.AssociationListOptions{},
			mockResponse: &api.SlurmdbV0043GetAssociationsResponse{
				HTTPResponse: &http.Response{StatusCode: 500},
			},
			expectedError: true,
		},
		{
			name:          "network error",
			opts:          &types.AssociationListOptions{},
			mockError:     fmt.Errorf("connection refused"),
			expectedError: true,
		},
		{
			name: "empty response",
			opts: &types.AssociationListOptions{},
			mockResponse: &api.SlurmdbV0043GetAssociationsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      nil,
			},
			expectedError: true,
		},
		{
			name: "with pagination",
			opts: &types.AssociationListOptions{
				Limit:  10,
				Offset: 5,
			},
			mockResponse: &api.SlurmdbV0043GetAssociationsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiAssocsResp{
					Associations: []api.V0043Assoc{
						{Account: ptrString("account1")},
						{Account: ptrString("account2")},
					},
				},
			},
			expectedError: false,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAssociationClient{}
			adapter := &AssociationAdapter{
				client:      mockClient,
				BaseManager: base.NewBaseManager("v0.0.43", "Association"),
			}

			ctx := context.Background()
			if tt.name == "nil context" {
				ctx = nil
			}

			if tt.mockResponse != nil || tt.mockError != nil {
				mockClient.On("SlurmdbV0043GetAssociationsWithResponse", mock.Anything, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			result, err := adapter.List(ctx, tt.opts)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Associations, tt.expectedCount)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestAssociationAdapter_Get(t *testing.T) {
	tests := []struct {
		name          string
		associationID string
		mockResponse  *api.SlurmdbV0043GetAssociationResponse
		mockError     error
		expectedError bool
		expectedUser  string
	}{
		{
			name:          "successful get",
			associationID: "account1|user1|cluster1",
			mockResponse: &api.SlurmdbV0043GetAssociationResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiAssocsResp{
					Associations: []api.V0043Assoc{
						{
							Account:   ptrString("account1"),
							User:      ptrString("user1"),
							Cluster:   ptrString("cluster1"),
							Partition: ptrString("normal"),
							Id:        ptrInt32(123),
						},
					},
				},
			},
			expectedError: false,
			expectedUser:  "user1",
		},
		{
			name:          "association not found",
			associationID: "nonexistent",
			mockResponse: &api.SlurmdbV0043GetAssociationResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiAssocsResp{
					Associations: []api.V0043Assoc{},
				},
			},
			expectedError: true,
		},
		{
			name:          "API error",
			associationID: "account1|user1|cluster1",
			mockResponse: &api.SlurmdbV0043GetAssociationResponse{
				HTTPResponse: &http.Response{StatusCode: 404},
			},
			expectedError: true,
		},
		{
			name:          "network error",
			associationID: "account1|user1|cluster1",
			mockError:     fmt.Errorf("timeout"),
			expectedError: true,
		},
		{
			name:          "empty association ID",
			associationID: "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAssociationClient{}
			adapter := &AssociationAdapter{
				client:      mockClient,
				BaseManager: base.NewBaseManager("v0.0.43", "Association"),
			}

			if tt.associationID != "" && (tt.mockResponse != nil || tt.mockError != nil) {
				mockClient.On("SlurmdbV0043GetAssociationWithResponse", mock.Anything, tt.associationID, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			result, err := adapter.Get(context.Background(), tt.associationID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedUser, result.UserName)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestAssociationAdapter_Create(t *testing.T) {
	tests := []struct {
		name          string
		association   *types.AssociationCreate
		mockResponse  *api.SlurmdbV0043PostAssociationsResponse
		mockError     error
		expectedError bool
	}{
		{
			name: "successful create",
			association: &types.AssociationCreate{
				AccountName:    "account1",
				UserName:       "user1",
				ClusterName:    "cluster1",
				PartitionName:  "normal",
				ParentAccount:  "root",
			},
			mockResponse: &api.SlurmdbV0043PostAssociationsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      &api.V0043OpenapiResp{},
			},
			expectedError: false,
		},
		{
			name: "minimal association",
			association: &types.AssociationCreate{
				AccountName: "account1",
				UserName:    "user1",
			},
			mockResponse: &api.SlurmdbV0043PostAssociationsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      &api.V0043OpenapiResp{},
			},
			expectedError: false,
		},
		{
			name: "API error",
			association: &types.AssociationCreate{
				AccountName: "account1",
				UserName:    "user1",
			},
			mockResponse: &api.SlurmdbV0043PostAssociationsResponse{
				HTTPResponse: &http.Response{StatusCode: 400},
			},
			expectedError: true,
		},
		{
			name: "network error",
			association: &types.AssociationCreate{
				AccountName: "account1",
				UserName:    "user1",
			},
			mockError:     fmt.Errorf("connection failed"),
			expectedError: true,
		},
		{
			name:          "nil association",
			association:   nil,
			expectedError: true,
		},
		{
			name: "missing required fields",
			association: &types.AssociationCreate{
				AccountName: "",
				UserName:    "user1",
			},
			expectedError: true,
		},
		{
			name: "missing user name",
			association: &types.AssociationCreate{
				AccountName: "account1",
				UserName:    "",
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAssociationClient{}
			adapter := &AssociationAdapter{
				client:      mockClient,
				BaseManager: base.NewBaseManager("v0.0.43", "Association"),
			}

			if tt.association != nil && tt.association.AccountName != "" && tt.association.UserName != "" {
				mockClient.On("SlurmdbV0043PostAssociationsWithResponse", mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			result, err := adapter.Create(context.Background(), tt.association)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestAssociationAdapter_Update(t *testing.T) {
	tests := []struct {
		name          string
		associationID string
		update        *types.AssociationUpdate
		mockGetResp   *api.SlurmdbV0043GetAssociationResponse
		mockPostResp  *api.SlurmdbV0043PostAssociationsResponse
		mockError     error
		expectedError bool
	}{
		{
			name:          "successful update",
			associationID: "account1|user1|cluster1",
			update: &types.AssociationUpdate{
				DefaultQoS: ptrString("normal"),
				GrpCPUs:    ptrInt(100),
			},
			mockGetResp: &api.SlurmdbV0043GetAssociationResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiAssocsResp{
					Associations: []api.V0043Assoc{
						{
							Account: ptrString("account1"),
							User:    ptrString("user1"),
							Cluster: ptrString("cluster1"),
						},
					},
				},
			},
			mockPostResp: &api.SlurmdbV0043PostAssociationsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      &api.V0043OpenapiResp{},
			},
			expectedError: false,
		},
		{
			name:          "update QoS only",
			associationID: "account1|user1|cluster1",
			update: &types.AssociationUpdate{
				QoSList: []string{"normal", "high"},
			},
			mockGetResp: &api.SlurmdbV0043GetAssociationResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiAssocsResp{
					Associations: []api.V0043Assoc{
						{
							Account: ptrString("account1"),
							User:    ptrString("user1"),
							Cluster: ptrString("cluster1"),
						},
					},
				},
			},
			mockPostResp: &api.SlurmdbV0043PostAssociationsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      &api.V0043OpenapiResp{},
			},
			expectedError: false,
		},
		{
			name:          "API error on update",
			associationID: "account1|user1|cluster1",
			update: &types.AssociationUpdate{
				DefaultQoS: ptrString("invalid"),
			},
			mockGetResp: &api.SlurmdbV0043GetAssociationResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiAssocsResp{
					Associations: []api.V0043Assoc{
						{
							Account: ptrString("account1"),
							User:    ptrString("user1"),
							Cluster: ptrString("cluster1"),
						},
					},
				},
			},
			mockPostResp: &api.SlurmdbV0043PostAssociationsResponse{
				HTTPResponse: &http.Response{StatusCode: 400},
			},
			expectedError: true,
		},
		{
			name:          "nil update",
			associationID: "account1|user1|cluster1",
			update:        nil,
			expectedError: true,
		},
		{
			name:          "empty association ID",
			associationID: "",
			update:        &types.AssociationUpdate{},
			expectedError: true,
		},
		{
			name:          "association not found",
			associationID: "nonexistent",
			update:        &types.AssociationUpdate{},
			mockGetResp: &api.SlurmdbV0043GetAssociationResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiAssocsResp{
					Associations: []api.V0043Assoc{},
				},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAssociationClient{}
			adapter := &AssociationAdapter{
				client:      mockClient,
				BaseManager: base.NewBaseManager("v0.0.43", "Association"),
			}

			if tt.associationID != "" && tt.update != nil {
				if tt.mockGetResp != nil {
					mockClient.On("SlurmdbV0043GetAssociationWithResponse", mock.Anything, tt.associationID, mock.Anything).
						Return(tt.mockGetResp, nil).Once()
				}
				if tt.mockPostResp != nil && len(tt.mockGetResp.JSON200.Associations) > 0 {
					mockClient.On("SlurmdbV0043PostAssociationsWithResponse", mock.Anything, mock.Anything, mock.Anything).
						Return(tt.mockPostResp, tt.mockError).Once()
				}
			}

			err := adapter.Update(context.Background(), tt.associationID, tt.update)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestAssociationAdapter_Delete(t *testing.T) {
	tests := []struct {
		name          string
		associationID string
		mockResponse  *api.SlurmdbV0043DeleteAssociationResponse
		mockError     error
		expectedError bool
	}{
		{
			name:          "successful delete",
			associationID: "account1|user1|cluster1",
			mockResponse: &api.SlurmdbV0043DeleteAssociationResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
			},
			expectedError: false,
		},
		{
			name:          "association not found",
			associationID: "nonexistent",
			mockResponse: &api.SlurmdbV0043DeleteAssociationResponse{
				HTTPResponse: &http.Response{StatusCode: 404},
			},
			expectedError: true,
		},
		{
			name:          "network error",
			associationID: "account1|user1|cluster1",
			mockError:     fmt.Errorf("network error"),
			expectedError: true,
		},
		{
			name:          "empty association ID",
			associationID: "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAssociationClient{}
			adapter := &AssociationAdapter{
				client:      mockClient,
				BaseManager: base.NewBaseManager("v0.0.43", "Association"),
			}

			if tt.associationID != "" {
				mockClient.On("SlurmdbV0043DeleteAssociationWithResponse", mock.Anything, tt.associationID, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			err := adapter.Delete(context.Background(), tt.associationID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestAssociationAdapter_validateAssociationCreate(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Association"),
	}

	tests := []struct {
		name          string
		association   *types.AssociationCreate
		expectedError bool
		errorContains string
	}{
		{
			name: "valid association",
			association: &types.AssociationCreate{
				AccountName:   "account1",
				UserName:      "user1",
				ClusterName:   "cluster1",
				PartitionName: "normal",
			},
			expectedError: false,
		},
		{
			name:          "nil association",
			association:   nil,
			expectedError: true,
			errorContains: "association creation data is required",
		},
		{
			name: "missing account name",
			association: &types.AssociationCreate{
				AccountName: "",
				UserName:    "user1",
			},
			expectedError: true,
			errorContains: "account name is required",
		},
		{
			name: "missing user name",
			association: &types.AssociationCreate{
				AccountName: "account1",
				UserName:    "",
			},
			expectedError: true,
			errorContains: "user name is required",
		},
		{
			name: "minimal valid",
			association: &types.AssociationCreate{
				AccountName: "account1",
				UserName:    "user1",
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.validateAssociationCreate(tt.association)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAssociationAdapter_validateAssociationUpdate(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Association"),
	}

	tests := []struct {
		name          string
		update        *types.AssociationUpdate
		expectedError bool
		errorContains string
	}{
		{
			name: "valid update",
			update: &types.AssociationUpdate{
				DefaultQoS: ptrString("normal"),
				GrpCPUs:    ptrInt(100),
			},
			expectedError: false,
		},
		{
			name:          "nil update",
			update:        nil,
			expectedError: true,
			errorContains: "association update data is required",
		},
		{
			name:          "empty update",
			update:        &types.AssociationUpdate{},
			expectedError: false, // Empty updates are allowed
		},
		{
			name: "update with QoS list",
			update: &types.AssociationUpdate{
				QoSList: []string{"normal", "high", "critical"},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.validateAssociationUpdate(tt.update)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAssociationAdapter_convertAPIAssociationToCommon(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Association"),
	}

	tests := []struct {
		name           string
		apiAssociation api.V0043Assoc
		expectedUser   string
		expectedAcct   string
	}{
		{
			name: "full association",
			apiAssociation: api.V0043Assoc{
				Account:    ptrString("account1"),
				User:       ptrString("user1"),
				Cluster:    ptrString("cluster1"),
				Partition:  ptrString("normal"),
				Id:         ptrInt32(123),
				DefaultQos: ptrString("normal"),
				GrpTres: &[]api.V0043Tres{
					{
						Type:  "cpu",
						Count: ptrInt64(100),
					},
				},
			},
			expectedUser: "user1",
			expectedAcct: "account1",
		},
		{
			name: "minimal association",
			apiAssociation: api.V0043Assoc{
				Account: ptrString("account2"),
				User:    ptrString("user2"),
			},
			expectedUser: "user2",
			expectedAcct: "account2",
		},
		{
			name:           "empty association",
			apiAssociation: api.V0043Assoc{},
			expectedUser:   "",
			expectedAcct:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.convertAPIAssociationToCommon(tt.apiAssociation)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedUser, result.UserName)
			assert.Equal(t, tt.expectedAcct, result.AccountName)
		})
	}
}

func TestAssociationAdapter_getDefaultClusterName(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Association"),
	}

	// Test that it returns a non-empty default cluster name
	clusterName := adapter.getDefaultClusterName()
	assert.NotEmpty(t, clusterName)
	// Common default is "cluster"
	assert.Contains(t, []string{"cluster", "default", "main"}, clusterName)
}

func TestAssociationAdapter_ValidateContext(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Association"),
	}

	tests := []struct {
		name    string
		ctx     context.Context
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil context",
			ctx:     nil,
			wantErr: true,
			errMsg:  "context is required",
		},
		{
			name:    "valid context",
			ctx:     context.Background(),
			wantErr: false,
		},
		{
			name:    "cancelled context",
			ctx:     func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			wantErr: false, // Context validation doesn't check cancellation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateContext(tt.ctx)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

