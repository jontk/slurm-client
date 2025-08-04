// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

// Mock client for testing
type MockQoSClient struct {
	mock.Mock
}

func (m *MockQoSClient) SlurmdbV0042GetQosWithResponse(ctx context.Context, params *api.SlurmdbV0042GetQosParams, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0042GetQosResponse, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmdbV0042GetQosResponse), args.Error(1)
}

func (m *MockQoSClient) SlurmdbV0042PostQosWithResponse(ctx context.Context, params *api.SlurmdbV0042PostQosParams, body api.SlurmdbV0042PostQosJSONRequestBody, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0042PostQosResponse, error) {
	args := m.Called(ctx, params, body)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmdbV0042PostQosResponse), args.Error(1)
}

func TestNewQoSAdapter(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewQoSAdapter(client)
	
	assert.NotNil(t, adapter)
	assert.Equal(t, client, adapter.client)
	assert.NotNil(t, adapter.BaseManager)
}

func TestQoSAdapter_List(t *testing.T) {
	tests := []struct {
		name          string
		opts          *types.QoSListOptions
		mockResponse  *api.SlurmdbV0042GetQosResponse
		mockError     error
		expectedError bool
		expectedCount int
	}{
		{
			name: "successful list",
			opts: &types.QoSListOptions{},
			mockResponse: &api.SlurmdbV0042GetQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiSlurmdbdQosResp{
					Qos: []api.V0042Qos{
						{
							Name: ptrString("qos1"),
							Id:   ptrInt32(1),
							Priority: ptrInt32(100),
						},
						{
							Name: ptrString("qos2"),
							Id:   ptrInt32(2),
							Priority: ptrInt32(200),
						},
					},
				},
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name: "list with name filter",
			opts: &types.QoSListOptions{
				Names: []string{"qos1"},
			},
			mockResponse: &api.SlurmdbV0042GetQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiSlurmdbdQosResp{
					Qos: []api.V0042Qos{
						{
							Name: ptrString("qos1"),
							Id:   ptrInt32(1),
							Priority: ptrInt32(100),
						},
					},
				},
			},
			expectedError: false,
			expectedCount: 1,
		},
		{
			name:          "nil context",
			opts:          nil,
			expectedError: true,
		},
		{
			name: "API error",
			opts: &types.QoSListOptions{},
			mockResponse: &api.SlurmdbV0042GetQosResponse{
				HTTPResponse: &http.Response{StatusCode: 500},
			},
			expectedError: true,
		},
		{
			name:          "network error",
			opts:          &types.QoSListOptions{},
			mockError:     fmt.Errorf("network error"),
			expectedError: true,
		},
		{
			name: "nil response",
			opts: &types.QoSListOptions{},
			mockResponse: &api.SlurmdbV0042GetQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      nil,
			},
			expectedError: true,
		},
		{
			name: "empty QoS list",
			opts: &types.QoSListOptions{},
			mockResponse: &api.SlurmdbV0042GetQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiSlurmdbdQosResp{
					Qos: []api.V0042Qos{},
				},
			},
			expectedError: false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockQoSClient{}
			adapter := &QoSAdapter{
				client:      mockClient,
				BaseManager: base.NewBaseManager("v0.0.42", "QoS"),
			}

			ctx := context.Background()
			if tt.name == "nil context" {
				ctx = nil
			}

			if tt.mockResponse != nil || tt.mockError != nil {
				mockClient.On("SlurmdbV0042GetQosWithResponse", mock.Anything, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			result, err := adapter.List(ctx, tt.opts)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.QoS, tt.expectedCount)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestQoSAdapter_Get(t *testing.T) {
	tests := []struct {
		name          string
		qosName       string
		mockResponse  *api.SlurmdbV0042GetQosResponse
		mockError     error
		expectedError bool
		expectedName  string
	}{
		{
			name:    "successful get",
			qosName: "test-qos",
			mockResponse: &api.SlurmdbV0042GetQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiSlurmdbdQosResp{
					Qos: []api.V0042Qos{
						{
							Name: ptrString("test-qos"),
							Id:   ptrInt32(42),
							Priority: ptrInt32(500),
							Description: ptrString("Test QoS for testing"),
							Limits: &api.V0042QosLimits{
								Max: &api.V0042QosLimitsMax{
									Jobs: &struct {
										PerAccount *int32 `json:"per_account,omitempty"`
										PerUser    *int32 `json:"per_user,omitempty"`
										Total      *int32 `json:"total,omitempty"`
									}{
										Total:      ptrInt32(100),
										PerAccount: ptrInt32(50),
										PerUser:    ptrInt32(10),
									},
								},
							},
						},
					},
				},
			},
			expectedError: false,
			expectedName:  "test-qos",
		},
		{
			name:    "QoS not found",
			qosName: "nonexistent",
			mockResponse: &api.SlurmdbV0042GetQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiSlurmdbdQosResp{
					Qos: []api.V0042Qos{},
				},
			},
			expectedError: true,
		},
		{
			name:    "API error",
			qosName: "test-qos",
			mockResponse: &api.SlurmdbV0042GetQosResponse{
				HTTPResponse: &http.Response{StatusCode: 404},
			},
			expectedError: true,
		},
		{
			name:          "network error",
			qosName:       "test-qos",
			mockError:     fmt.Errorf("connection refused"),
			expectedError: true,
		},
		{
			name:          "empty QoS name",
			qosName:       "",
			expectedError: false, // Will call API but find nothing
			mockResponse: &api.SlurmdbV0042GetQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiSlurmdbdQosResp{
					Qos: []api.V0042Qos{},
				},
			},
		},
		{
			name:    "nil response",
			qosName: "test-qos",
			mockResponse: &api.SlurmdbV0042GetQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      nil,
			},
			expectedError: true,
		},
		{
			name:    "QoS name mismatch",
			qosName: "expected-qos",
			mockResponse: &api.SlurmdbV0042GetQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiSlurmdbdQosResp{
					Qos: []api.V0042Qos{
						{
							Name: ptrString("different-qos"),
							Id:   ptrInt32(1),
						},
					},
				},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockQoSClient{}
			adapter := &QoSAdapter{
				client:      mockClient,
				BaseManager: base.NewBaseManager("v0.0.42", "QoS"),
			}

			if tt.mockResponse != nil || tt.mockError != nil {
				mockClient.On("SlurmdbV0042GetQosWithResponse", mock.Anything, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			result, err := adapter.Get(context.Background(), tt.qosName)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedName, result.Name)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestQoSAdapter_Create(t *testing.T) {
	tests := []struct {
		name          string
		qos           *types.QoSCreate
		mockResponse  *api.SlurmdbV0042PostQosResponse
		mockError     error
		expectedError bool
	}{
		{
			name: "successful create",
			qos: &types.QoSCreate{
				Name:        "new-qos",
				Description: "New QoS for testing",
				Priority:    100,
			},
			mockResponse: &api.SlurmdbV0042PostQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiResp{
					Meta: &api.V0042OpenapiMeta{},
				},
			},
			expectedError: false,
		},
		{
			name: "create with minimal fields",
			qos: &types.QoSCreate{
				Name: "minimal-qos",
			},
			mockResponse: &api.SlurmdbV0042PostQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiResp{
					Meta: &api.V0042OpenapiMeta{},
				},
			},
			expectedError: false,
		},
		{
			name: "API error",
			qos: &types.QoSCreate{
				Name: "new-qos",
			},
			mockResponse: &api.SlurmdbV0042PostQosResponse{
				HTTPResponse: &http.Response{StatusCode: 400},
			},
			expectedError: true,
		},
		{
			name: "network error",
			qos: &types.QoSCreate{
				Name: "new-qos",
			},
			mockError:     fmt.Errorf("connection failed"),
			expectedError: true,
		},
		{
			name:          "nil qos",
			qos:           nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockQoSClient{}
			adapter := &QoSAdapter{
				client:      mockClient,
				BaseManager: base.NewBaseManager("v0.0.42", "QoS"),
			}

			if tt.qos != nil {
				mockClient.On("SlurmdbV0042PostQosWithResponse", mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			result, err := adapter.Create(context.Background(), tt.qos)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.qos != nil {
					assert.Equal(t, tt.qos.Name, result.QoSName)
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestQoSAdapter_Update(t *testing.T) {
	adapter := NewQoSAdapter(nil)
	
	updates := &types.QoSUpdateRequest{
		Priority: ptrInt32(200),
	}
	
	err := adapter.Update(context.Background(), "test-qos", updates)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not directly supported")
}

func TestQoSAdapter_Delete(t *testing.T) {
	adapter := NewQoSAdapter(nil)
	
	err := adapter.Delete(context.Background(), "test-qos")
	
	assert.Error(t, err)
	// Should return NotImplementedError
	assert.Contains(t, err.Error(), "not implemented")
}

func TestQoSAdapter_convertAPIQoSToCommon(t *testing.T) {
	adapter := NewQoSAdapter(nil)

	tests := []struct {
		name        string
		apiQoS      api.V0042Qos
		expected    *types.QoS
		expectError bool
	}{
		{
			name: "full QoS conversion",
			apiQoS: api.V0042Qos{
				Name: ptrString("full-qos"),
				Id:   ptrInt32(42),
				Priority: ptrInt32(500),
				Description: ptrString("Full featured QoS"),
			},
			expected: &types.QoS{
				Name:        "full-qos",
				ID:          42,
				Priority:    500,
				Description: "Full featured QoS",
			},
		},
		{
			name: "minimal QoS conversion",
			apiQoS: api.V0042Qos{
				Name: ptrString("minimal-qos"),
			},
			expected: &types.QoS{
				Name: "minimal-qos",
			},
		},
		{
			name:   "empty QoS",
			apiQoS: api.V0042Qos{},
			expected: &types.QoS{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.convertAPIQoSToCommon(tt.apiQoS)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Name, result.Name)
				assert.Equal(t, tt.expected.ID, result.ID)
				assert.Equal(t, tt.expected.Priority, result.Priority)
				assert.Equal(t, tt.expected.Description, result.Description)
			}
		})
	}
}

func TestQoSAdapter_convertCommonQoSCreateToAPI(t *testing.T) {
	adapter := NewQoSAdapter(nil)

	tests := []struct {
		name        string
		qosCreate   *QoSCreateRequest
		expectError bool
	}{
		{
			name: "successful conversion",
			qosCreate: &QoSCreateRequest{
				Name:        "test-qos",
				Description: "Test QoS",
				Priority:    100,
			},
			expectError: false,
		},
		{
			name: "minimal conversion",
			qosCreate: &QoSCreateRequest{
				Name: "minimal-qos",
			},
			expectError: false,
		},
		{
			name:        "nil input",
			qosCreate:   nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.convertCommonQoSCreateToAPI(tt.qosCreate)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.qosCreate != nil && len(result.Qos) > 0 {
					assert.Equal(t, tt.qosCreate.Name, *result.Qos[0].Name)
					if tt.qosCreate.Description != "" {
						assert.Equal(t, tt.qosCreate.Description, *result.Qos[0].Description)
					}
					if tt.qosCreate.Priority > 0 {
						assert.Equal(t, tt.qosCreate.Priority, *result.Qos[0].Priority)
					}
				}
			}
		})
	}
}

func TestQoSAdapter_ContextValidation(t *testing.T) {
	adapter := NewQoSAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name string
		ctx  context.Context
	}{
		{"nil context", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := adapter.List(tt.ctx, nil)
			assert.Error(t, err)

			_, err = adapter.Get(tt.ctx, "test")
			assert.Error(t, err)

			_, err = adapter.Create(tt.ctx, &types.QoSCreate{Name: "test"})
			assert.Error(t, err)

			err = adapter.Update(tt.ctx, "test", &types.QoSUpdateRequest{})
			assert.Error(t, err)

			err = adapter.Delete(tt.ctx, "test")
			assert.Error(t, err)
		})
	}
}

func TestQoSAdapter_ClientValidation(t *testing.T) {
	adapter := NewQoSAdapter(nil)
	ctx := context.Background()

	_, err := adapter.List(ctx, nil)
	assert.Error(t, err)

	_, err = adapter.Get(ctx, "test")
	assert.Error(t, err)

	_, err = adapter.Create(ctx, &types.QoSCreate{Name: "test"})
	assert.Error(t, err)

	err = adapter.Update(ctx, "test", &types.QoSUpdateRequest{})
	assert.Error(t, err)

	err = adapter.Delete(ctx, "test")
	assert.Error(t, err)
}