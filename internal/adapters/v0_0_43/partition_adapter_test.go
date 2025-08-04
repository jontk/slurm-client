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
type MockPartitionClient struct {
	mock.Mock
}

func (m *MockPartitionClient) SlurmV0043GetPartitionsWithResponse(ctx context.Context, params *api.SlurmV0043GetPartitionsParams, reqEditors ...api.RequestEditorFn) (*api.SlurmV0043GetPartitionsResponse, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmV0043GetPartitionsResponse), args.Error(1)
}

func (m *MockPartitionClient) SlurmV0043GetPartitionWithResponse(ctx context.Context, partitionName string, params *api.SlurmV0043GetPartitionParams, reqEditors ...api.RequestEditorFn) (*api.SlurmV0043GetPartitionResponse, error) {
	args := m.Called(ctx, partitionName, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmV0043GetPartitionResponse), args.Error(1)
}

func (m *MockPartitionClient) SlurmdbV0043PostPartitionsWithResponse(ctx context.Context, body api.SlurmdbV0043PostPartitionsJSONRequestBody, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0043PostPartitionsResponse, error) {
	args := m.Called(ctx, body)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmdbV0043PostPartitionsResponse), args.Error(1)
}

func (m *MockPartitionClient) SlurmV0043PostPartitionWithResponse(ctx context.Context, partitionName string, body api.SlurmV0043PostPartitionJSONRequestBody, reqEditors ...api.RequestEditorFn) (*api.SlurmV0043PostPartitionResponse, error) {
	args := m.Called(ctx, partitionName, body)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmV0043PostPartitionResponse), args.Error(1)
}

func (m *MockPartitionClient) SlurmV0043DeletePartitionWithResponse(ctx context.Context, partitionName string, reqEditors ...api.RequestEditorFn) (*api.SlurmV0043DeletePartitionResponse, error) {
	args := m.Called(ctx, partitionName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmV0043DeletePartitionResponse), args.Error(1)
}

func TestPartitionAdapter_List(t *testing.T) {
	tests := []struct {
		name          string
		opts          *types.PartitionListOptions
		mockResponse  *api.SlurmV0043GetPartitionsResponse
		mockError     error
		expectedError bool
		expectedCount int
	}{
		{
			name: "successful list",
			opts: &types.PartitionListOptions{},
			mockResponse: &api.SlurmV0043GetPartitionsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiPartitionResp{
					Partitions: []api.V0043PartitionInfo{
						{
							Name:  ptrString("normal"),
							Nodes: &api.V0043PartitionInfoNodes{
								Total: ptrInt32(10),
							},
							State: &[]api.V0043PartitionInfoState{"UP"},
						},
						{
							Name:  ptrString("gpu"),
							Nodes: &api.V0043PartitionInfoNodes{
								Total: ptrInt32(5),
							},
							State: &[]api.V0043PartitionInfoState{"UP"},
						},
					},
				},
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name: "list with filter",
			opts: &types.PartitionListOptions{
				Names: []string{"normal"},
			},
			mockResponse: &api.SlurmV0043GetPartitionsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiPartitionResp{
					Partitions: []api.V0043PartitionInfo{
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
			opts: &types.PartitionListOptions{},
			mockResponse: &api.SlurmV0043GetPartitionsResponse{
				HTTPResponse: &http.Response{StatusCode: 500},
			},
			expectedError: true,
		},
		{
			name:          "network error",
			opts:          &types.PartitionListOptions{},
			mockError:     fmt.Errorf("connection refused"),
			expectedError: true,
		},
		{
			name: "empty response",
			opts: &types.PartitionListOptions{},
			mockResponse: &api.SlurmV0043GetPartitionsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      nil,
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockPartitionClient{}
			adapter := &PartitionAdapter{
				client:      mockClient,
				BaseManager: NewPartitionAdapter(nil).BaseManager,
			}

			ctx := context.Background()
			if tt.name == "nil context" {
				ctx = nil
			}

			if tt.mockResponse != nil || tt.mockError != nil {
				mockClient.On("SlurmV0043GetPartitionsWithResponse", mock.Anything, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			result, err := adapter.List(ctx, tt.opts)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Partitions, tt.expectedCount)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestPartitionAdapter_Get(t *testing.T) {
	tests := []struct {
		name          string
		partitionName string
		mockResponse  *api.SlurmV0043GetPartitionResponse
		mockError     error
		expectedError bool
		expectedName  string
	}{
		{
			name:          "successful get",
			partitionName: "normal",
			mockResponse: &api.SlurmV0043GetPartitionResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiPartitionResp{
					Partitions: []api.V0043PartitionInfo{
						{
							Name:        ptrString("normal"),
							State:       &[]api.V0043PartitionInfoState{"UP"},
							MaxTime:     &api.V0043Uint32NoVal{Number: ptrUint32(86400)},
							DefaultTime: &api.V0043Uint32NoVal{Number: ptrUint32(3600)},
							Nodes: &api.V0043PartitionInfoNodes{
								Total: ptrInt32(10),
							},
						},
					},
				},
			},
			expectedError: false,
			expectedName:  "normal",
		},
		{
			name:          "partition not found",
			partitionName: "nonexistent",
			mockResponse: &api.SlurmV0043GetPartitionResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiPartitionResp{
					Partitions: []api.V0043PartitionInfo{},
				},
			},
			expectedError: true,
		},
		{
			name:          "API error",
			partitionName: "normal",
			mockResponse: &api.SlurmV0043GetPartitionResponse{
				HTTPResponse: &http.Response{StatusCode: 404},
			},
			expectedError: true,
		},
		{
			name:          "network error",
			partitionName: "normal",
			mockError:     fmt.Errorf("timeout"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockPartitionClient{}
			adapter := &PartitionAdapter{
				client:      mockClient,
				BaseManager: NewPartitionAdapter(nil).BaseManager,
			}

			mockClient.On("SlurmV0043GetPartitionWithResponse", mock.Anything, tt.partitionName, mock.Anything).
				Return(tt.mockResponse, tt.mockError)

			result, err := adapter.Get(context.Background(), tt.partitionName)

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

func TestPartitionAdapter_Create(t *testing.T) {
	tests := []struct {
		name          string
		partition     *types.PartitionCreate
		mockResponse  *api.SlurmdbV0043PostPartitionsResponse
		mockError     error
		expectedError bool
	}{
		{
			name: "successful create",
			partition: &types.PartitionCreate{
				Name:        "new-partition",
				Nodes:       "node[1-10]",
				DefaultTime: 3600,
				MaxTime:     86400,
				State:       "UP",
			},
			mockResponse: &api.SlurmdbV0043PostPartitionsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      &api.V0043OpenapiResp{},
			},
			expectedError: false,
		},
		{
			name: "minimal partition",
			partition: &types.PartitionCreate{
				Name: "minimal",
			},
			mockResponse: &api.SlurmdbV0043PostPartitionsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      &api.V0043OpenapiResp{},
			},
			expectedError: false,
		},
		{
			name: "API error",
			partition: &types.PartitionCreate{
				Name: "new-partition",
			},
			mockResponse: &api.SlurmdbV0043PostPartitionsResponse{
				HTTPResponse: &http.Response{StatusCode: 400},
			},
			expectedError: true,
		},
		{
			name: "network error",
			partition: &types.PartitionCreate{
				Name: "new-partition",
			},
			mockError:     fmt.Errorf("connection failed"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockPartitionClient{}
			adapter := &PartitionAdapter{
				client:      mockClient,
				BaseManager: NewPartitionAdapter(nil).BaseManager,
			}

			mockClient.On("SlurmdbV0043PostPartitionsWithResponse", mock.Anything, mock.Anything).
				Return(tt.mockResponse, tt.mockError)

			result, err := adapter.Create(context.Background(), tt.partition)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.partition.Name, result.Name)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestPartitionAdapter_Update(t *testing.T) {
	tests := []struct {
		name          string
		partitionName string
		update        *types.PartitionUpdate
		mockResponse  *api.SlurmV0043PostPartitionResponse
		mockError     error
		expectedError bool
	}{
		{
			name:          "successful update",
			partitionName: "normal",
			update: &types.PartitionUpdate{
				State:       ptrString("DOWN"),
				MaxTime:     ptrInt(172800),
				DefaultTime: ptrInt(7200),
			},
			mockResponse: &api.SlurmV0043PostPartitionResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
			},
			expectedError: false,
		},
		{
			name:          "update state only",
			partitionName: "normal",
			update: &types.PartitionUpdate{
				State: ptrString("DRAIN"),
			},
			mockResponse: &api.SlurmV0043PostPartitionResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
			},
			expectedError: false,
		},
		{
			name:          "API error",
			partitionName: "normal",
			update: &types.PartitionUpdate{
				State: ptrString("INVALID_STATE"),
			},
			mockResponse: &api.SlurmV0043PostPartitionResponse{
				HTTPResponse: &http.Response{StatusCode: 400},
			},
			expectedError: true,
		},
		{
			name:          "nil update",
			partitionName: "normal",
			update:        nil,
			expectedError: true,
		},
		{
			name:          "empty partition name",
			partitionName: "",
			update:        &types.PartitionUpdate{},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockPartitionClient{}
			adapter := &PartitionAdapter{
				client:      mockClient,
				BaseManager: NewPartitionAdapter(nil).BaseManager,
			}

			if tt.mockResponse != nil || tt.mockError != nil {
				mockClient.On("SlurmV0043PostPartitionWithResponse", mock.Anything, tt.partitionName, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			err := adapter.Update(context.Background(), tt.partitionName, tt.update)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestPartitionAdapter_Delete(t *testing.T) {
	tests := []struct {
		name          string
		partitionName string
		mockResponse  *api.SlurmV0043DeletePartitionResponse
		mockError     error
		expectedError bool
	}{
		{
			name:          "successful delete",
			partitionName: "old-partition",
			mockResponse: &api.SlurmV0043DeletePartitionResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
			},
			expectedError: false,
		},
		{
			name:          "partition not found",
			partitionName: "nonexistent",
			mockResponse: &api.SlurmV0043DeletePartitionResponse{
				HTTPResponse: &http.Response{StatusCode: 404},
			},
			expectedError: true,
		},
		{
			name:          "network error",
			partitionName: "old-partition",
			mockError:     fmt.Errorf("network error"),
			expectedError: true,
		},
		{
			name:          "empty partition name",
			partitionName: "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockPartitionClient{}
			adapter := &PartitionAdapter{
				client:      mockClient,
				BaseManager: NewPartitionAdapter(nil).BaseManager,
			}

			if tt.partitionName != "" && (tt.mockResponse != nil || tt.mockError != nil) {
				mockClient.On("SlurmV0043DeletePartitionWithResponse", mock.Anything, tt.partitionName).
					Return(tt.mockResponse, tt.mockError)
			}

			err := adapter.Delete(context.Background(), tt.partitionName)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestPartitionAdapter_convertAPIPartitionToCommon(t *testing.T) {
	adapter := NewPartitionAdapter(nil)

	tests := []struct {
		name         string
		apiPartition api.V0043PartitionInfo
		expected     *types.Partition
	}{
		{
			name: "full partition",
			apiPartition: api.V0043PartitionInfo{
				Name:     ptrString("normal"),
				State:    &[]api.V0043PartitionInfoState{"UP"},
				Priority: &api.V0043Uint32NoVal{Number: ptrUint32(1000)},
				MaxTime:  &api.V0043Uint32NoVal{Number: ptrUint32(86400)},
				DefaultTime: &api.V0043Uint32NoVal{Number: ptrUint32(3600)},
				Nodes: &api.V0043PartitionInfoNodes{
					Total:      ptrInt32(10),
					Configured: ptrString("node[1-10]"),
				},
				Cpus: &api.V0043PartitionInfoCpus{
					Total: ptrUint32(160),
				},
				MaxNodes: &api.V0043Uint32NoVal{Number: ptrUint32(5)},
				QoS: &api.V0043PartitionInfoQos{
					Allowed: ptrString("normal,high"),
				},
			},
			expected: &types.Partition{
				Name:         "normal",
				State:        "UP",
				Priority:     1000,
				MaxTime:      86400,
				DefaultTime:  3600,
				Nodes:        "node[1-10]",
				TotalNodes:   10,
				TotalCPUs:    160,
				MaxNodesPerJob: 5,
				AllowedQoS:   "normal,high",
			},
		},
		{
			name: "minimal partition",
			apiPartition: api.V0043PartitionInfo{
				Name: ptrString("minimal"),
			},
			expected: &types.Partition{
				Name: "minimal",
			},
		},
		{
			name: "partition with multiple states",
			apiPartition: api.V0043PartitionInfo{
				Name:  ptrString("multi-state"),
				State: &[]api.V0043PartitionInfoState{"UP", "DRAIN"},
			},
			expected: &types.Partition{
				Name:  "multi-state",
				State: "UP",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.convertAPIPartitionToCommon(tt.apiPartition)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.State, result.State)
			assert.Equal(t, tt.expected.Priority, result.Priority)
			assert.Equal(t, tt.expected.MaxTime, result.MaxTime)
			assert.Equal(t, tt.expected.DefaultTime, result.DefaultTime)
			assert.Equal(t, tt.expected.TotalNodes, result.TotalNodes)
			assert.Equal(t, tt.expected.TotalCPUs, result.TotalCPUs)
		})
	}
}

func TestPartitionAdapter_validatePartitionUpdate(t *testing.T) {
	adapter := NewPartitionAdapter(nil)

	tests := []struct {
		name          string
		partitionName string
		update        *types.PartitionUpdate
		expectedError bool
		errorContains string
	}{
		{
			name:          "valid update",
			partitionName: "normal",
			update:        &types.PartitionUpdate{State: ptrString("DOWN")},
			expectedError: false,
		},
		{
			name:          "empty partition name",
			partitionName: "",
			update:        &types.PartitionUpdate{},
			expectedError: true,
			errorContains: "partition name is required",
		},
		{
			name:          "nil update",
			partitionName: "normal",
			update:        nil,
			expectedError: true,
			errorContains: "update request is required",
		},
		{
			name:          "empty update",
			partitionName: "normal",
			update:        &types.PartitionUpdate{},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.validatePartitionUpdate(tt.partitionName, tt.update)

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

