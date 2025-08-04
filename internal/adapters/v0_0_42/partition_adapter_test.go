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
type MockPartitionClient struct {
	mock.Mock
}

func (m *MockPartitionClient) SlurmV0042GetPartitionsWithResponse(ctx context.Context, params *api.SlurmV0042GetPartitionsParams, reqEditors ...api.RequestEditorFn) (*api.SlurmV0042GetPartitionsResponse, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmV0042GetPartitionsResponse), args.Error(1)
}

func (m *MockPartitionClient) SlurmV0042GetPartitionWithResponse(ctx context.Context, partitionName string, params *api.SlurmV0042GetPartitionParams, reqEditors ...api.RequestEditorFn) (*api.SlurmV0042GetPartitionResponse, error) {
	args := m.Called(ctx, partitionName, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmV0042GetPartitionResponse), args.Error(1)
}

func TestNewPartitionAdapter(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewPartitionAdapter(client)
	
	assert.NotNil(t, adapter)
	assert.Equal(t, client, adapter.client)
	assert.NotNil(t, adapter.BaseManager)
}

func TestPartitionAdapter_List(t *testing.T) {
	tests := []struct {
		name          string
		opts          *types.PartitionListOptions
		mockResponse  *api.SlurmV0042GetPartitionsResponse
		mockError     error
		expectedError bool
		expectedCount int
	}{
		{
			name: "successful list",
			opts: &types.PartitionListOptions{},
			mockResponse: &api.SlurmV0042GetPartitionsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiPartitionResp{
					Partitions: []api.V0042PartitionInfo{
						{
							Name: ptrString("partition1"),
							State: &api.V0042PartitionInfoState{
								State: &[]api.V0042PartitionInfoStateState{"UP"},
							},
							Nodes: &struct {
								Allowed *string `json:"allowed,omitempty"`
								Configured *string `json:"configured,omitempty"`
							}{
								Configured: ptrString("node[1-10]"),
							},
							MaxCpusPerNode:     ptrInt32(24),
							MaxMemoryPerNode:   ptrInt64(128000),
							MaxNodesPerJob:     ptrInt32(100),
							DefaultMemoryPerCpu: ptrInt64(4000),
							Default:            ptrBool(false),
						},
						{
							Name: ptrString("partition2"),
							State: &api.V0042PartitionInfoState{
								State: &[]api.V0042PartitionInfoStateState{"DOWN"},
							},
							Nodes: &struct {
								Allowed *string `json:"allowed,omitempty"`
								Configured *string `json:"configured,omitempty"`
							}{
								Configured: ptrString("node[11-20]"),
							},
							MaxCpusPerNode:      ptrInt32(32),
							MaxMemoryPerNode:    ptrInt64(256000),
							MaxNodesPerJob:      ptrInt32(50),
							DefaultMemoryPerCpu: ptrInt64(8000),
							Default:             ptrBool(true),
						},
					},
				},
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name: "list with name filter",
			opts: &types.PartitionListOptions{
				Names: []string{"partition1"},
			},
			mockResponse: &api.SlurmV0042GetPartitionsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiPartitionResp{
					Partitions: []api.V0042PartitionInfo{
						{
							Name: ptrString("partition1"),
							State: &api.V0042PartitionInfoState{
								State: &[]api.V0042PartitionInfoStateState{"UP"},
							},
							Nodes: &struct {
								Allowed *string `json:"allowed,omitempty"`
								Configured *string `json:"configured,omitempty"`
							}{
								Configured: ptrString("node[1-10]"),
							},
						},
						{
							Name: ptrString("partition2"),
							State: &api.V0042PartitionInfoState{
								State: &[]api.V0042PartitionInfoStateState{"DOWN"},
							},
							Nodes: &struct {
								Allowed *string `json:"allowed,omitempty"`
								Configured *string `json:"configured,omitempty"`
							}{
								Configured: ptrString("node[11-20]"),
							},
						},
					},
				},
			},
			expectedError: false,
			expectedCount: 1, // Only partition1 should be returned
		},
		{
			name:          "nil context",
			opts:          nil,
			expectedError: true,
		},
		{
			name: "API error",
			opts: &types.PartitionListOptions{},
			mockResponse: &api.SlurmV0042GetPartitionsResponse{
				HTTPResponse: &http.Response{StatusCode: 500},
			},
			expectedError: true,
		},
		{
			name:          "network error",
			opts:          &types.PartitionListOptions{},
			mockError:     fmt.Errorf("network error"),
			expectedError: true,
		},
		{
			name: "nil response",
			opts: &types.PartitionListOptions{},
			mockResponse: &api.SlurmV0042GetPartitionsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      nil,
			},
			expectedError: true,
		},
		{
			name: "empty partitions list",
			opts: &types.PartitionListOptions{},
			mockResponse: &api.SlurmV0042GetPartitionsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiPartitionResp{
					Partitions: []api.V0042PartitionInfo{},
				},
			},
			expectedError: false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockPartitionClient{}
			adapter := &PartitionAdapter{
				client:      mockClient,
				BaseManager: base.NewBaseManager("v0.0.42", "Partition"),
			}

			ctx := context.Background()
			if tt.name == "nil context" {
				ctx = nil
			}

			if tt.mockResponse != nil || tt.mockError != nil {
				mockClient.On("SlurmV0042GetPartitionsWithResponse", mock.Anything, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			result, err := adapter.List(ctx, tt.opts)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Partitions, tt.expectedCount)
				
				// Verify filtering worked correctly
				if tt.opts != nil && len(tt.opts.Names) > 0 && !tt.expectedError {
					for _, partition := range result.Partitions {
						found := false
						for _, name := range tt.opts.Names {
							if partition.Name == name {
								found = true
								break
							}
						}
						assert.True(t, found, "Partition %s should match filter", partition.Name)
					}
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestPartitionAdapter_Get(t *testing.T) {
	tests := []struct {
		name          string
		partitionName string
		mockResponse  *api.SlurmV0042GetPartitionResponse
		mockError     error
		expectedError bool
		expectedName  string
	}{
		{
			name:          "successful get",
			partitionName: "test-partition",
			mockResponse: &api.SlurmV0042GetPartitionResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiPartitionResp{
					Partitions: []api.V0042PartitionInfo{
						{
							Name: ptrString("test-partition"),
							State: &api.V0042PartitionInfoState{
								State: &[]api.V0042PartitionInfoStateState{"UP"},
							},
							Nodes: &struct {
								Allowed *string `json:"allowed,omitempty"`
								Configured *string `json:"configured,omitempty"`
							}{
								Configured: ptrString("node[1-100]"),
								Allowed:    ptrString("node[1-50]"),
							},
							MaxCpusPerNode:     ptrInt32(48),
							MaxMemoryPerNode:   ptrInt64(512000),
							MaxNodesPerJob:     ptrInt32(200),
							DefaultMemoryPerCpu: ptrInt64(2000),
							Default:            ptrBool(false),
							Priority:           ptrInt32(1000),
						},
					},
				},
			},
			expectedError: false,
			expectedName:  "test-partition",
		},
		{
			name:          "partition not found",
			partitionName: "nonexistent",
			mockResponse: &api.SlurmV0042GetPartitionResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiPartitionResp{
					Partitions: []api.V0042PartitionInfo{},
				},
			},
			expectedError: true,
		},
		{
			name:          "API error",
			partitionName: "test-partition",
			mockResponse: &api.SlurmV0042GetPartitionResponse{
				HTTPResponse: &http.Response{StatusCode: 404},
			},
			expectedError: true,
		},
		{
			name:          "network error",
			partitionName: "test-partition",
			mockError:     fmt.Errorf("connection refused"),
			expectedError: true,
		},
		{
			name:          "empty partition name",
			partitionName: "",
			expectedError: true,
		},
		{
			name:          "nil response",
			partitionName: "test-partition",
			mockResponse: &api.SlurmV0042GetPartitionResponse{
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
				BaseManager: base.NewBaseManager("v0.0.42", "Partition"),
			}

			if tt.partitionName != "" && (tt.mockResponse != nil || tt.mockError != nil) {
				mockClient.On("SlurmV0042GetPartitionWithResponse", mock.Anything, tt.partitionName, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

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
	adapter := NewPartitionAdapter(nil)
	
	partition := &types.PartitionCreate{
		Name: "new-partition",
	}
	
	result, err := adapter.Create(context.Background(), partition)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not supported")
}

func TestPartitionAdapter_Update(t *testing.T) {
	adapter := NewPartitionAdapter(nil)
	
	updates := &types.PartitionUpdateRequest{
		MaxNodes: ptrInt32(100),
	}
	
	err := adapter.Update(context.Background(), "test-partition", updates)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not supported")
}

func TestPartitionAdapter_Delete(t *testing.T) {
	adapter := NewPartitionAdapter(nil)
	
	err := adapter.Delete(context.Background(), "test-partition")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not supported")
}

func TestPartitionAdapter_convertAPIPartitionToCommon(t *testing.T) {
	adapter := NewPartitionAdapter(nil)

	tests := []struct {
		name         string
		apiPartition api.V0042PartitionInfo
		expected     *types.Partition
		expectError  bool
	}{
		{
			name: "full partition conversion",
			apiPartition: api.V0042PartitionInfo{
				Name: ptrString("full-partition"),
				State: &api.V0042PartitionInfoState{
					State: &[]api.V0042PartitionInfoStateState{"UP"},
				},
				Nodes: &struct {
					Allowed *string `json:"allowed,omitempty"`
					Configured *string `json:"configured,omitempty"`
				}{
					Configured: ptrString("node[1-100]"),
					Allowed:    ptrString("node[1-50]"),
				},
				MaxCpusPerNode:      ptrInt32(48),
				MaxMemoryPerNode:    ptrInt64(512000),
				MaxNodesPerJob:      ptrInt32(200),
				DefaultMemoryPerCpu: ptrInt64(2000),
				Default:             ptrBool(true),
				Priority:            ptrInt32(1000),
				Timeouts: &struct {
					DefMemPerCpu   *int32 `json:"def_mem_per_cpu,omitempty"`
					MaxMemPerCpu   *int32 `json:"max_mem_per_cpu,omitempty"`
					DefaultTime    *int32 `json:"default_time,omitempty"`
					MaxTime        *int32 `json:"max_time,omitempty"`
				}{
					DefaultTime: ptrInt32(60),
					MaxTime:     ptrInt32(1440),
				},
			},
			expected: &types.Partition{
				Name:                "full-partition",
				State:               []string{"UP"},
				AllowedNodes:        "node[1-50]",
				ConfiguredNodes:     "node[1-100]",
				MaxCPUsPerNode:      48,
				MaxMemoryPerNode:    512000,
				MaxNodesPerJob:      200,
				DefaultMemoryPerCPU: 2000,
				Default:             true,
				Priority:            1000,
				DefaultTime:         60,
				MaxTime:             1440,
				Meta:                make(map[string]interface{}),
			},
		},
		{
			name: "minimal partition conversion",
			apiPartition: api.V0042PartitionInfo{
				Name: ptrString("minimal-partition"),
			},
			expected: &types.Partition{
				Name: "minimal-partition",
				Meta: make(map[string]interface{}),
			},
		},
		{
			name:        "empty partition",
			apiPartition: api.V0042PartitionInfo{},
			expected: &types.Partition{
				Meta: make(map[string]interface{}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.convertAPIPartitionToCommon(tt.apiPartition)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Name, result.Name)
				assert.Equal(t, tt.expected.State, result.State)
				assert.Equal(t, tt.expected.AllowedNodes, result.AllowedNodes)
				assert.Equal(t, tt.expected.ConfiguredNodes, result.ConfiguredNodes)
				assert.Equal(t, tt.expected.MaxCPUsPerNode, result.MaxCPUsPerNode)
				assert.Equal(t, tt.expected.MaxMemoryPerNode, result.MaxMemoryPerNode)
				assert.Equal(t, tt.expected.MaxNodesPerJob, result.MaxNodesPerJob)
				assert.Equal(t, tt.expected.DefaultMemoryPerCPU, result.DefaultMemoryPerCPU)
				assert.Equal(t, tt.expected.Default, result.Default)
				assert.Equal(t, tt.expected.Priority, result.Priority)
				assert.Equal(t, tt.expected.DefaultTime, result.DefaultTime)
				assert.Equal(t, tt.expected.MaxTime, result.MaxTime)
				assert.NotNil(t, result.Meta)
			}
		})
	}
}

func TestPartitionAdapter_ContextValidation(t *testing.T) {
	adapter := NewPartitionAdapter(&api.ClientWithResponses{})

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

			_, err = adapter.Create(tt.ctx, &types.PartitionCreate{Name: "test"})
			assert.Error(t, err)

			err = adapter.Update(tt.ctx, "test", &types.PartitionUpdateRequest{})
			assert.Error(t, err)

			err = adapter.Delete(tt.ctx, "test")
			assert.Error(t, err)
		})
	}
}

func TestPartitionAdapter_ClientValidation(t *testing.T) {
	adapter := NewPartitionAdapter(nil)
	ctx := context.Background()

	_, err := adapter.List(ctx, nil)
	assert.Error(t, err)

	_, err = adapter.Get(ctx, "test")
	assert.Error(t, err)

	_, err = adapter.Create(ctx, &types.PartitionCreate{Name: "test"})
	assert.Error(t, err)

	err = adapter.Update(ctx, "test", &types.PartitionUpdateRequest{})
	assert.Error(t, err)

	err = adapter.Delete(ctx, "test")
	assert.Error(t, err)
}