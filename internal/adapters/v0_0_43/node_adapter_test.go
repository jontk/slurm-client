// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// Mock client for testing
type MockNodeClient struct {
	mock.Mock
}

func (m *MockNodeClient) SlurmV0043GetNodesWithResponse(ctx context.Context, params *api.SlurmV0043GetNodesParams, reqEditors ...api.RequestEditorFn) (*api.SlurmV0043GetNodesResponse, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmV0043GetNodesResponse), args.Error(1)
}

func (m *MockNodeClient) SlurmV0043GetNodeWithResponse(ctx context.Context, nodeName string, params *api.SlurmV0043GetNodeParams, reqEditors ...api.RequestEditorFn) (*api.SlurmV0043GetNodeResponse, error) {
	args := m.Called(ctx, nodeName, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmV0043GetNodeResponse), args.Error(1)
}

func (m *MockNodeClient) SlurmV0043PostNodeWithResponse(ctx context.Context, nodeName string, body api.SlurmV0043PostNodeJSONRequestBody, reqEditors ...api.RequestEditorFn) (*api.SlurmV0043PostNodeResponse, error) {
	args := m.Called(ctx, nodeName, body)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmV0043PostNodeResponse), args.Error(1)
}

func (m *MockNodeClient) SlurmV0043DeleteNodeWithResponse(ctx context.Context, nodeName string, reqEditors ...api.RequestEditorFn) (*api.SlurmV0043DeleteNodeResponse, error) {
	args := m.Called(ctx, nodeName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmV0043DeleteNodeResponse), args.Error(1)
}

func TestNewNodeAdapter(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewNodeAdapter(client)
	
	assert.NotNil(t, adapter)
	assert.Equal(t, client, adapter.client)
	assert.NotNil(t, adapter.BaseManager)
	assert.Nil(t, adapter.wrapper)
}

func TestNodeAdapter_List(t *testing.T) {
	tests := []struct {
		name          string
		opts          *types.NodeListOptions
		mockResponse  *api.SlurmV0043GetNodesResponse
		mockError     error
		expectedError bool
		expectedCount int
	}{
		{
			name: "successful list",
			opts: &types.NodeListOptions{},
			mockResponse: &api.SlurmV0043GetNodesResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiNodesResp{
					Nodes: []api.V0043Node{
						{
							Name:  ptrString("node01"),
							State: &[]api.V0043NodeState{"IDLE"},
							Cpus: &struct {
								Count           *int32  `json:"count,omitempty"`
								LoadAverage     *uint32 `json:"load_average,omitempty"`
							}{
								Count: ptrInt32(16),
							},
							RealMemory: ptrInt64(32768),
						},
						{
							Name:  ptrString("node02"),
							State: &[]api.V0043NodeState{"ALLOCATED"},
							Cpus: &struct {
								Count           *int32  `json:"count,omitempty"`
								LoadAverage     *uint32 `json:"load_average,omitempty"`
							}{
								Count: ptrInt32(32),
							},
							RealMemory: ptrInt64(65536),
						},
					},
				},
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name: "list with update time filter",
			opts: &types.NodeListOptions{
				UpdateTime: &time.Time{},
			},
			mockResponse: &api.SlurmV0043GetNodesResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiNodesResp{
					Nodes: []api.V0043Node{},
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
			opts: &types.NodeListOptions{},
			mockResponse: &api.SlurmV0043GetNodesResponse{
				HTTPResponse: &http.Response{StatusCode: 500},
			},
			expectedError: true,
		},
		{
			name:          "network error",
			opts:          &types.NodeListOptions{},
			mockError:     fmt.Errorf("connection refused"),
			expectedError: true,
		},
		{
			name: "empty response",
			opts: &types.NodeListOptions{},
			mockResponse: &api.SlurmV0043GetNodesResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      nil,
			},
			expectedError: true,
		},
		{
			name: "with pagination",
			opts: &types.NodeListOptions{
				Limit:  10,
				Offset: 5,
			},
			mockResponse: &api.SlurmV0043GetNodesResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiNodesResp{
					Nodes: []api.V0043Node{
						{Name: ptrString("node01")},
						{Name: ptrString("node02")},
					},
				},
			},
			expectedError: false,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockNodeClient{}
			adapter := &NodeAdapter{
				client:      mockClient,
				BaseManager: base.NewBaseManager("v0.0.43", "Node"),
			}

			ctx := context.Background()
			if tt.name == "nil context" {
				ctx = nil
			}

			if tt.mockResponse != nil || tt.mockError != nil {
				mockClient.On("SlurmV0043GetNodesWithResponse", mock.Anything, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			result, err := adapter.List(ctx, tt.opts)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Nodes, tt.expectedCount)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestNodeAdapter_Get(t *testing.T) {
	tests := []struct {
		name          string
		nodeName      string
		mockResponse  *api.SlurmV0043GetNodeResponse
		mockError     error
		expectedError bool
		expectedName  string
	}{
		{
			name:     "successful get",
			nodeName: "node01",
			mockResponse: &api.SlurmV0043GetNodeResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiNodesResp{
					Nodes: []api.V0043Node{
						{
							Name:  ptrString("node01"),
							State: &[]api.V0043NodeState{"IDLE"},
							Cpus: &struct {
								Count           *int32  `json:"count,omitempty"`
								LoadAverage     *uint32 `json:"load_average,omitempty"`
							}{
								Count: ptrInt32(16),
							},
							RealMemory: ptrInt64(32768),
							Arch:       ptrString("x86_64"),
							Address:    ptrString("192.168.1.101"),
							Partitions: &[]string{"normal", "gpu"},
						},
					},
				},
			},
			expectedError: false,
			expectedName:  "node01",
		},
		{
			name:     "node not found",
			nodeName: "nonexistent",
			mockResponse: &api.SlurmV0043GetNodeResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiNodesResp{
					Nodes: []api.V0043Node{},
				},
			},
			expectedError: true,
		},
		{
			name:     "API error",
			nodeName: "node01",
			mockResponse: &api.SlurmV0043GetNodeResponse{
				HTTPResponse: &http.Response{StatusCode: 404},
			},
			expectedError: true,
		},
		{
			name:          "network error",
			nodeName:      "node01",
			mockError:     fmt.Errorf("timeout"),
			expectedError: true,
		},
		{
			name:          "empty node name",
			nodeName:      "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockNodeClient{}
			adapter := &NodeAdapter{
				client:      mockClient,
				BaseManager: base.NewBaseManager("v0.0.43", "Node"),
			}

			if tt.nodeName != "" && (tt.mockResponse != nil || tt.mockError != nil) {
				mockClient.On("SlurmV0043GetNodeWithResponse", mock.Anything, tt.nodeName, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			result, err := adapter.Get(context.Background(), tt.nodeName)

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

func TestNodeAdapter_Update(t *testing.T) {
	tests := []struct {
		name          string
		nodeName      string
		update        *types.NodeUpdate
		mockResponse  *api.SlurmV0043PostNodeResponse
		mockError     error
		expectedError bool
	}{
		{
			name:     "successful update",
			nodeName: "node01",
			update: &types.NodeUpdate{
				State:    ptrString("DRAIN"),
				Reason:   ptrString("Maintenance"),
				Features: []string{"gpu", "nvme"},
			},
			mockResponse: &api.SlurmV0043PostNodeResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
			},
			expectedError: false,
		},
		{
			name:     "update comment only",
			nodeName: "node01",
			update: &types.NodeUpdate{
				Comment: ptrString("Updated comment"),
			},
			mockResponse: &api.SlurmV0043PostNodeResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
			},
			expectedError: false,
		},
		{
			name:     "API error",
			nodeName: "node01",
			update: &types.NodeUpdate{
				State: ptrString("INVALID_STATE"),
			},
			mockResponse: &api.SlurmV0043PostNodeResponse{
				HTTPResponse: &http.Response{StatusCode: 400},
			},
			expectedError: true,
		},
		{
			name:          "nil update",
			nodeName:      "node01",
			update:        nil,
			expectedError: true,
		},
		{
			name:          "empty node name",
			nodeName:      "",
			update:        &types.NodeUpdate{},
			expectedError: true,
		},
		{
			name:     "network error",
			nodeName: "node01",
			update: &types.NodeUpdate{
				State: ptrString("DRAIN"),
			},
			mockError:     fmt.Errorf("connection failed"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockNodeClient{}
			adapter := &NodeAdapter{
				client:      mockClient,
				BaseManager: base.NewBaseManager("v0.0.43", "Node"),
			}

			if tt.nodeName != "" && tt.update != nil && (tt.mockResponse != nil || tt.mockError != nil) {
				mockClient.On("SlurmV0043PostNodeWithResponse", mock.Anything, tt.nodeName, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			err := adapter.Update(context.Background(), tt.nodeName, tt.update)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestNodeAdapter_Delete(t *testing.T) {
	tests := []struct {
		name          string
		nodeName      string
		mockResponse  *api.SlurmV0043DeleteNodeResponse
		mockError     error
		expectedError bool
	}{
		{
			name:     "successful delete",
			nodeName: "old-node",
			mockResponse: &api.SlurmV0043DeleteNodeResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
			},
			expectedError: false,
		},
		{
			name:     "node not found",
			nodeName: "nonexistent",
			mockResponse: &api.SlurmV0043DeleteNodeResponse{
				HTTPResponse: &http.Response{StatusCode: 404},
			},
			expectedError: true,
		},
		{
			name:          "network error",
			nodeName:      "old-node",
			mockError:     fmt.Errorf("network error"),
			expectedError: true,
		},
		{
			name:          "empty node name",
			nodeName:      "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockNodeClient{}
			adapter := &NodeAdapter{
				client:      mockClient,
				BaseManager: base.NewBaseManager("v0.0.43", "Node"),
			}

			if tt.nodeName != "" && (tt.mockResponse != nil || tt.mockError != nil) {
				mockClient.On("SlurmV0043DeleteNodeWithResponse", mock.Anything, tt.nodeName).
					Return(tt.mockResponse, tt.mockError)
			}

			err := adapter.Delete(context.Background(), tt.nodeName)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestNodeAdapter_Watch(t *testing.T) {
	tests := []struct {
		name          string
		opts          *types.NodeWatchOptions
		expectedError bool
	}{
		{
			name: "successful watch",
			opts: &types.NodeWatchOptions{
				Names:        []string{"node01", "node02"},
				PollInterval: 1 * time.Second,
			},
			expectedError: false,
		},
		{
			name: "watch all nodes",
			opts: &types.NodeWatchOptions{
				PollInterval: 2 * time.Second,
			},
			expectedError: false,
		},
		{
			name:          "nil context",
			opts:          &types.NodeWatchOptions{},
			expectedError: true,
		},
		{
			name:          "nil options",
			opts:          nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockNodeClient{}
			adapter := &NodeAdapter{
				client:      mockClient,
				BaseManager: base.NewBaseManager("v0.0.43", "Node"),
			}

			ctx := context.Background()
			if tt.name == "nil context" {
				ctx = nil
			}

			if !tt.expectedError {
				// Setup mock for initial poll
				mockResponse := &api.SlurmV0043GetNodesResponse{
					HTTPResponse: &http.Response{StatusCode: 200},
					JSON200: &api.V0043OpenapiNodesResp{
						Nodes: []api.V0043Node{
							{
								Name:  ptrString("node01"),
								State: &[]api.V0043NodeState{"IDLE"},
							},
						},
					},
				}
				mockClient.On("SlurmV0043GetNodesWithResponse", mock.Anything, mock.Anything).
					Return(mockResponse, nil).Maybe()
			}

			eventCh, err := adapter.Watch(ctx, tt.opts)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, eventCh)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, eventCh)

				// Cancel context to stop watching
				if ctx != nil {
					cancelCtx, cancel := context.WithCancel(ctx)
					defer cancel()
					
					// Re-run with cancellable context
					eventCh, err = adapter.Watch(cancelCtx, tt.opts)
					assert.NoError(t, err)
					
					// Give it a moment to start
					time.Sleep(100 * time.Millisecond)
					cancel()
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestNodeAdapter_validateNodeUpdate(t *testing.T) {
	adapter := &NodeAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Node"),
	}

	tests := []struct {
		name          string
		update        *types.NodeUpdate
		expectedError bool
		errorContains string
	}{
		{
			name: "valid update",
			update: &types.NodeUpdate{
				State:  ptrString("DRAIN"),
				Reason: ptrString("Maintenance"),
			},
			expectedError: false,
		},
		{
			name:          "nil update",
			update:        nil,
			expectedError: true,
			errorContains: "node update data is required",
		},
		{
			name: "invalid state",
			update: &types.NodeUpdate{
				State: ptrString("INVALID_STATE"),
			},
			expectedError: true,
			errorContains: "invalid node state",
		},
		{
			name: "empty update",
			update: &types.NodeUpdate{},
			expectedError: false, // Empty updates are allowed
		},
		{
			name: "update with features",
			update: &types.NodeUpdate{
				Features: []string{"gpu", "nvme", "infiniband"},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.validateNodeUpdate(tt.update)

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

func TestNodeAdapter_filterNodeList(t *testing.T) {
	adapter := &NodeAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Node"),
	}

	nodes := []types.Node{
		{
			Name:       "node01",
			State:      types.NodeStateIdle,
			Partitions: []string{"normal", "gpu"},
		},
		{
			Name:       "node02",
			State:      types.NodeStateAllocated,
			Partitions: []string{"normal"},
		},
		{
			Name:       "node03",
			State:      types.NodeStateDrain,
			Partitions: []string{"gpu"},
		},
	}

	tests := []struct {
		name          string
		opts          *types.NodeListOptions
		expectedCount int
		expectedNodes []string
	}{
		{
			name:          "no filters",
			opts:          nil,
			expectedCount: 3,
			expectedNodes: []string{"node01", "node02", "node03"},
		},
		{
			name: "filter by names",
			opts: &types.NodeListOptions{
				Names: []string{"node01", "node03"},
			},
			expectedCount: 2,
			expectedNodes: []string{"node01", "node03"},
		},
		{
			name: "filter by states",
			opts: &types.NodeListOptions{
				States: []types.NodeState{types.NodeStateIdle, types.NodeStateDrain},
			},
			expectedCount: 2,
			expectedNodes: []string{"node01", "node03"},
		},
		{
			name: "filter by partitions",
			opts: &types.NodeListOptions{
				Partitions: []string{"gpu"},
			},
			expectedCount: 2,
			expectedNodes: []string{"node01", "node03"},
		},
		{
			name: "combined filters",
			opts: &types.NodeListOptions{
				Names:      []string{"node01", "node02"},
				Partitions: []string{"gpu"},
			},
			expectedCount: 1,
			expectedNodes: []string{"node01"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.filterNodeList(nodes, tt.opts)
			assert.Len(t, result, tt.expectedCount)

			resultNames := make([]string, len(result))
			for i, node := range result {
				resultNames[i] = node.Name
			}
			assert.Equal(t, tt.expectedNodes, resultNames)
		})
	}
}

func TestNodeAdapter_ValidateContext(t *testing.T) {
	adapter := &NodeAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Node"),
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

