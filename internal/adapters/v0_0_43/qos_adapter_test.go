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
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// Mock client for testing
type MockQoSClient struct {
	mock.Mock
}

func (m *MockQoSClient) SlurmdbV0043GetQosWithResponse(ctx context.Context, params *api.SlurmdbV0043GetQosParams, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0043GetQosResponse, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmdbV0043GetQosResponse), args.Error(1)
}

func (m *MockQoSClient) SlurmdbV0043GetSingleQosWithResponse(ctx context.Context, qos string, params *api.SlurmdbV0043GetSingleQosParams, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0043GetSingleQosResponse, error) {
	args := m.Called(ctx, qos, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmdbV0043GetSingleQosResponse), args.Error(1)
}

func (m *MockQoSClient) SlurmdbV0043PostQosWithResponse(ctx context.Context, params *api.SlurmdbV0043PostQosParams, body api.SlurmdbV0043PostQosJSONRequestBody, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0043PostQosResponse, error) {
	args := m.Called(ctx, params, body)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmdbV0043PostQosResponse), args.Error(1)
}

func (m *MockQoSClient) SlurmdbV0043DeleteSingleQosWithResponse(ctx context.Context, qos string, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0043DeleteSingleQosResponse, error) {
	args := m.Called(ctx, qos)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmdbV0043DeleteSingleQosResponse), args.Error(1)
}

func TestQoSAdapter_List(t *testing.T) {
	tests := []struct {
		name          string
		opts          *types.QoSListOptions
		mockResponse  *api.SlurmdbV0043GetQosResponse
		mockError     error
		expectedError bool
		expectedCount int
	}{
		{
			name: "successful list",
			opts: &types.QoSListOptions{},
			mockResponse: &api.SlurmdbV0043GetQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiSlurmdbdQosResp{
					Qos: []api.V0043Qos{
						{
							Name:        ptrString("normal"),
							Priority:    ptrUint32(100),
							Description: ptrString("Normal QoS"),
						},
						{
							Name:        ptrString("high"),
							Priority:    ptrUint32(200),
							Description: ptrString("High priority QoS"),
						},
					},
				},
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name: "list with filter",
			opts: &types.QoSListOptions{
				Names: []string{"normal"},
			},
			mockResponse: &api.SlurmdbV0043GetQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiSlurmdbdQosResp{
					Qos: []api.V0043Qos{
						{
							Name: ptrString("normal"),
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
			mockResponse: &api.SlurmdbV0043GetQosResponse{
				HTTPResponse: &http.Response{StatusCode: 500},
			},
			expectedError: true,
		},
		{
			name:          "network error",
			opts:          &types.QoSListOptions{},
			mockError:     fmt.Errorf("connection refused"),
			expectedError: true,
		},
		{
			name: "empty response",
			opts: &types.QoSListOptions{},
			mockResponse: &api.SlurmdbV0043GetQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      nil,
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockQoSClient{}
			adapter := &QoSAdapter{
				client:      mockClient,
				BaseManager: NewQoSAdapter(nil).BaseManager,
			}

			ctx := context.Background()
			if tt.name == "nil context" {
				ctx = nil
			}

			if tt.mockResponse != nil || tt.mockError != nil {
				mockClient.On("SlurmdbV0043GetQosWithResponse", mock.Anything, mock.Anything).
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
		mockResponse  *api.SlurmdbV0043GetSingleQosResponse
		mockError     error
		expectedError bool
		expectedName  string
	}{
		{
			name:    "successful get",
			qosName: "normal",
			mockResponse: &api.SlurmdbV0043GetSingleQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiSlurmdbdQosResp{
					Qos: []api.V0043Qos{
						{
							Name:        ptrString("normal"),
							Priority:    ptrUint32(100),
							Description: ptrString("Normal QoS"),
							MaxCpusPerUser: &api.V0043Uint32NoVal{
								Number: ptrUint32(1000),
							},
							MaxJobsPerUser: &api.V0043Uint32NoVal{
								Number: ptrUint32(10),
							},
						},
					},
				},
			},
			expectedError: false,
			expectedName:  "normal",
		},
		{
			name:    "qos not found",
			qosName: "nonexistent",
			mockResponse: &api.SlurmdbV0043GetSingleQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiSlurmdbdQosResp{
					Qos: []api.V0043Qos{},
				},
			},
			expectedError: true,
		},
		{
			name:    "API error",
			qosName: "normal",
			mockResponse: &api.SlurmdbV0043GetSingleQosResponse{
				HTTPResponse: &http.Response{StatusCode: 404},
			},
			expectedError: true,
		},
		{
			name:      "network error",
			qosName:   "normal",
			mockError: fmt.Errorf("timeout"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockQoSClient{}
			adapter := &QoSAdapter{
				client:      mockClient,
				BaseManager: NewQoSAdapter(nil).BaseManager,
			}

			mockClient.On("SlurmdbV0043GetSingleQosWithResponse", mock.Anything, tt.qosName, mock.Anything).
				Return(tt.mockResponse, tt.mockError)

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
		mockResponse  *api.SlurmdbV0043PostQosResponse
		mockError     error
		expectedError bool
	}{
		{
			name: "successful create",
			qos: &types.QoSCreate{
				Name:        "new-qos",
				Priority:    150,
				Description: "New QoS for testing",
				Limits: types.QoSLimits{
					MaxCPUsPerUser: ptrInt(100),
					MaxJobsPerUser: ptrInt(5),
				},
			},
			mockResponse: &api.SlurmdbV0043PostQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      &api.V0043OpenapiResp{},
			},
			expectedError: false,
		},
		{
			name: "minimal qos",
			qos: &types.QoSCreate{
				Name: "minimal",
			},
			mockResponse: &api.SlurmdbV0043PostQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      &api.V0043OpenapiResp{},
			},
			expectedError: false,
		},
		{
			name: "API error",
			qos: &types.QoSCreate{
				Name: "new-qos",
			},
			mockResponse: &api.SlurmdbV0043PostQosResponse{
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
		{
			name: "empty name",
			qos: &types.QoSCreate{
				Name: "",
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockQoSClient{}
			adapter := &QoSAdapter{
				client:      mockClient,
				BaseManager: NewQoSAdapter(nil).BaseManager,
			}

			if tt.qos != nil && tt.qos.Name != "" {
				mockClient.On("SlurmdbV0043PostQosWithResponse", mock.Anything, mock.Anything, mock.Anything).
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
	tests := []struct {
		name          string
		qosName       string
		update        *types.QoSUpdate
		mockResponse  *api.SlurmdbV0043PostQosResponse
		mockError     error
		expectedError bool
	}{
		{
			name:    "successful update",
			qosName: "normal",
			update: &types.QoSUpdate{
				Priority:    ptrInt(150),
				Description: ptrString("Updated description"),
				Limits: &types.QoSLimits{
					MaxCPUsPerUser: ptrInt(200),
				},
			},
			mockResponse: &api.SlurmdbV0043PostQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      &api.V0043OpenapiResp{},
			},
			expectedError: false,
		},
		{
			name:    "update priority only",
			qosName: "normal",
			update: &types.QoSUpdate{
				Priority: ptrInt(250),
			},
			mockResponse: &api.SlurmdbV0043PostQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      &api.V0043OpenapiResp{},
			},
			expectedError: false,
		},
		{
			name:    "API error",
			qosName: "normal",
			update: &types.QoSUpdate{
				Priority: ptrInt(-1),
			},
			mockResponse: &api.SlurmdbV0043PostQosResponse{
				HTTPResponse: &http.Response{StatusCode: 400},
			},
			expectedError: true,
		},
		{
			name:          "nil update",
			qosName:       "normal",
			update:        nil,
			expectedError: true,
		},
		{
			name:          "empty qos name",
			qosName:       "",
			update:        &types.QoSUpdate{},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockQoSClient{}
			adapter := &QoSAdapter{
				client:      mockClient,
				BaseManager: NewQoSAdapter(nil).BaseManager,
			}

			if tt.qosName != "" && tt.update != nil {
				mockClient.On("SlurmdbV0043PostQosWithResponse", mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			err := adapter.Update(context.Background(), tt.qosName, tt.update)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestQoSAdapter_Delete(t *testing.T) {
	tests := []struct {
		name          string
		qosName       string
		mockResponse  *api.SlurmdbV0043DeleteSingleQosResponse
		mockError     error
		expectedError bool
	}{
		{
			name:    "successful delete",
			qosName: "old-qos",
			mockResponse: &api.SlurmdbV0043DeleteSingleQosResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
			},
			expectedError: false,
		},
		{
			name:    "qos not found",
			qosName: "nonexistent",
			mockResponse: &api.SlurmdbV0043DeleteSingleQosResponse{
				HTTPResponse: &http.Response{StatusCode: 404},
			},
			expectedError: true,
		},
		{
			name:      "network error",
			qosName:   "old-qos",
			mockError: fmt.Errorf("network error"),
			expectedError: true,
		},
		{
			name:          "empty qos name",
			qosName:       "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockQoSClient{}
			adapter := &QoSAdapter{
				client:      mockClient,
				BaseManager: NewQoSAdapter(nil).BaseManager,
			}

			if tt.qosName != "" {
				mockClient.On("SlurmdbV0043DeleteSingleQosWithResponse", mock.Anything, tt.qosName).
					Return(tt.mockResponse, tt.mockError)
			}

			err := adapter.Delete(context.Background(), tt.qosName)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestQoSAdapter_validateQoSCreate(t *testing.T) {
	adapter := NewQoSAdapter(nil)

	tests := []struct {
		name          string
		qos           *types.QoSCreate
		expectedError bool
		errorContains string
	}{
		{
			name: "valid qos",
			qos: &types.QoSCreate{
				Name:     "test-qos",
				Priority: 100,
			},
			expectedError: false,
		},
		{
			name:          "nil qos",
			qos:           nil,
			expectedError: true,
			errorContains: "QoS is required",
		},
		{
			name: "empty name",
			qos: &types.QoSCreate{
				Name: "",
			},
			expectedError: true,
			errorContains: "QoS name is required",
		},
		{
			name: "invalid name with spaces",
			qos: &types.QoSCreate{
				Name: "test qos",
			},
			expectedError: true,
			errorContains: "QoS name cannot contain spaces",
		},
		{
			name: "invalid name with special chars",
			qos: &types.QoSCreate{
				Name: "test@qos",
			},
			expectedError: true,
			errorContains: "QoS name contains invalid characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.validateQoSCreate(tt.qos)

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

func TestQoSAdapter_validateQoSUpdate(t *testing.T) {
	adapter := NewQoSAdapter(nil)

	tests := []struct {
		name          string
		qosName       string
		update        *types.QoSUpdate
		expectedError bool
		errorContains string
	}{
		{
			name:    "valid update",
			qosName: "test-qos",
			update:  &types.QoSUpdate{Priority: ptrInt(100)},
			expectedError: false,
		},
		{
			name:          "empty qos name",
			qosName:       "",
			update:        &types.QoSUpdate{},
			expectedError: true,
			errorContains: "QoS name is required",
		},
		{
			name:          "nil update",
			qosName:       "test-qos",
			update:        nil,
			expectedError: true,
			errorContains: "update is required",
		},
		{
			name:    "negative priority",
			qosName: "test-qos",
			update:  &types.QoSUpdate{Priority: ptrInt(-1)},
			expectedError: true,
			errorContains: "priority cannot be negative",
		},
		{
			name:    "empty update object",
			qosName: "test-qos",
			update:  &types.QoSUpdate{},
			expectedError: false, // Empty updates are allowed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.validateQoSUpdate(tt.qosName, tt.update)

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