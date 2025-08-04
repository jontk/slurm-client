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
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// Mock client for testing
type MockClusterClient struct {
	mock.Mock
}

func (m *MockClusterClient) SlurmdbV0043GetClustersWithResponse(ctx context.Context, params *api.SlurmdbV0043GetClustersParams, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0043GetClustersResponse, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmdbV0043GetClustersResponse), args.Error(1)
}

func (m *MockClusterClient) SlurmdbV0043GetClusterWithResponse(ctx context.Context, clusterName string, params *api.SlurmdbV0043GetClusterParams, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0043GetClusterResponse, error) {
	args := m.Called(ctx, clusterName, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmdbV0043GetClusterResponse), args.Error(1)
}

func (m *MockClusterClient) SlurmdbV0043PostClustersWithResponse(ctx context.Context, params *api.SlurmdbV0043PostClustersParams, body api.SlurmdbV0043PostClustersJSONRequestBody, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0043PostClustersResponse, error) {
	args := m.Called(ctx, params, body)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmdbV0043PostClustersResponse), args.Error(1)
}

func (m *MockClusterClient) SlurmdbV0043DeleteClusterWithResponse(ctx context.Context, clusterName string, params *api.SlurmdbV0043DeleteClusterParams, reqEditors ...api.RequestEditorFn) (*api.SlurmdbV0043DeleteClusterResponse, error) {
	args := m.Called(ctx, clusterName, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmdbV0043DeleteClusterResponse), args.Error(1)
}

func TestClusterAdapter_List(t *testing.T) {
	tests := []struct {
		name          string
		opts          *types.ClusterListOptions
		mockResponse  *api.SlurmdbV0043GetClustersResponse
		mockError     error
		expectedError bool
		expectedCount int
	}{
		{
			name: "successful list",
			opts: &types.ClusterListOptions{},
			mockResponse: &api.SlurmdbV0043GetClustersResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiClustersResp{
					Clusters: []api.V0043ClusterRec{
						{
							Name: ptrString("cluster1"),
							Controller: &struct {
								Host *string `json:"host,omitempty"`
								Port *int32  `json:"port,omitempty"`
							}{
								Host: ptrString("host1"),
								Port: ptrInt32(6817),
							},
						},
						{
							Name: ptrString("cluster2"),
							Controller: &struct {
								Host *string `json:"host,omitempty"`
								Port *int32  `json:"port,omitempty"`
							}{
								Host: ptrString("host2"),
								Port: ptrInt32(6817),
							},
						},
					},
				},
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name: "list with update time filter",
			opts: &types.ClusterListOptions{
				UpdateTime: &time.Time{},
			},
			mockResponse: &api.SlurmdbV0043GetClustersResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiClustersResp{
					Clusters: []api.V0043ClusterRec{},
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
			opts: &types.ClusterListOptions{},
			mockResponse: &api.SlurmdbV0043GetClustersResponse{
				HTTPResponse: &http.Response{StatusCode: 500},
			},
			expectedError: true,
		},
		{
			name:          "network error",
			opts:          &types.ClusterListOptions{},
			mockError:     fmt.Errorf("network error"),
			expectedError: true,
		},
		{
			name: "nil response",
			opts: &types.ClusterListOptions{},
			mockResponse: &api.SlurmdbV0043GetClustersResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      nil,
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockClusterClient{}
			adapter := &ClusterAdapter{
				client:      mockClient,
				BaseManager: NewClusterAdapter(nil).BaseManager,
			}

			ctx := context.Background()
			if tt.name == "nil context" {
				ctx = nil
			}

			if tt.mockResponse != nil || tt.mockError != nil {
				mockClient.On("SlurmdbV0043GetClustersWithResponse", mock.Anything, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			result, err := adapter.List(ctx, tt.opts)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Clusters, tt.expectedCount)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestClusterAdapter_Get(t *testing.T) {
	tests := []struct {
		name          string
		clusterName   string
		mockResponse  *api.SlurmdbV0043GetClusterResponse
		mockError     error
		expectedError bool
		expectedName  string
	}{
		{
			name:        "successful get",
			clusterName: "test-cluster",
			mockResponse: &api.SlurmdbV0043GetClusterResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiClustersResp{
					Clusters: []api.V0043ClusterRec{
						{
							Name: ptrString("test-cluster"),
							Controller: &struct {
								Host *string `json:"host,omitempty"`
								Port *int32  `json:"port,omitempty"`
							}{
								Host: ptrString("controller.example.com"),
								Port: ptrInt32(6817),
							},
							Nodes:      ptrString("node[1-10]"),
							RpcVersion: ptrInt32(21),
						},
					},
				},
			},
			expectedError: false,
			expectedName:  "test-cluster",
		},
		{
			name:        "cluster not found",
			clusterName: "nonexistent",
			mockResponse: &api.SlurmdbV0043GetClusterResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiClustersResp{
					Clusters: []api.V0043ClusterRec{},
				},
			},
			expectedError: true,
		},
		{
			name:        "API error",
			clusterName: "test-cluster",
			mockResponse: &api.SlurmdbV0043GetClusterResponse{
				HTTPResponse: &http.Response{StatusCode: 404},
			},
			expectedError: true,
		},
		{
			name:          "network error",
			clusterName:   "test-cluster",
			mockError:     fmt.Errorf("connection refused"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockClusterClient{}
			adapter := &ClusterAdapter{
				client:      mockClient,
				BaseManager: NewClusterAdapter(nil).BaseManager,
			}

			mockClient.On("SlurmdbV0043GetClusterWithResponse", mock.Anything, tt.clusterName, mock.Anything).
				Return(tt.mockResponse, tt.mockError)

			result, err := adapter.Get(context.Background(), tt.clusterName)

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

func TestClusterAdapter_Create(t *testing.T) {
	tests := []struct {
		name          string
		cluster       *types.ClusterCreate
		mockResponse  *api.SlurmdbV0043PostClustersResponse
		mockError     error
		expectedError bool
	}{
		{
			name: "successful create",
			cluster: &types.ClusterCreate{
				Name:           "new-cluster",
				ControllerHost: "controller.example.com",
				ControllerPort: 6817,
				Nodes:          "node[1-5]",
				RpcVersion:     21,
				Flags:          []string{"EXTERNAL"},
			},
			mockResponse: &api.SlurmdbV0043PostClustersResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiResp{
					Meta: &api.V0043OpenapiMeta{},
				},
			},
			expectedError: false,
		},
		{
			name: "create with minimal fields",
			cluster: &types.ClusterCreate{
				Name: "minimal-cluster",
			},
			mockResponse: &api.SlurmdbV0043PostClustersResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiResp{
					Meta: &api.V0043OpenapiMeta{},
				},
			},
			expectedError: false,
		},
		{
			name: "API error",
			cluster: &types.ClusterCreate{
				Name: "new-cluster",
			},
			mockResponse: &api.SlurmdbV0043PostClustersResponse{
				HTTPResponse: &http.Response{StatusCode: 400},
			},
			expectedError: true,
		},
		{
			name: "create with error response",
			cluster: &types.ClusterCreate{
				Name: "new-cluster",
			},
			mockResponse: &api.SlurmdbV0043PostClustersResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiResp{
					Errors: &[]api.V0043OpenapiError{
						{
							Error: ptrString("Cluster already exists"),
						},
					},
				},
			},
			expectedError: false, // Response is successful but contains error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockClusterClient{}
			adapter := &ClusterAdapter{
				client:      mockClient,
				BaseManager: NewClusterAdapter(nil).BaseManager,
			}

			mockClient.On("SlurmdbV0043PostClustersWithResponse", mock.Anything, mock.Anything, mock.Anything).
				Return(tt.mockResponse, tt.mockError)

			result, err := adapter.Create(context.Background(), tt.cluster)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.cluster.Name, result.Name)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestClusterAdapter_Delete(t *testing.T) {
	tests := []struct {
		name          string
		clusterName   string
		mockResponse  *api.SlurmdbV0043DeleteClusterResponse
		mockError     error
		expectedError bool
	}{
		{
			name:        "successful delete",
			clusterName: "old-cluster",
			mockResponse: &api.SlurmdbV0043DeleteClusterResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
			},
			expectedError: false,
		},
		{
			name:        "cluster not found",
			clusterName: "nonexistent",
			mockResponse: &api.SlurmdbV0043DeleteClusterResponse{
				HTTPResponse: &http.Response{StatusCode: 404},
			},
			expectedError: true,
		},
		{
			name:          "network error",
			clusterName:   "old-cluster",
			mockError:     fmt.Errorf("timeout"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockClusterClient{}
			adapter := &ClusterAdapter{
				client:      mockClient,
				BaseManager: NewClusterAdapter(nil).BaseManager,
			}

			mockClient.On("SlurmdbV0043DeleteClusterWithResponse", mock.Anything, tt.clusterName, mock.Anything).
				Return(tt.mockResponse, tt.mockError)

			err := adapter.Delete(context.Background(), tt.clusterName)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestClusterAdapter_convertAPIClusterToCommon(t *testing.T) {
	adapter := NewClusterAdapter(nil)

	tests := []struct {
		name        string
		apiCluster  api.V0043ClusterRec
		expected    *types.Cluster
		expectError bool
	}{
		{
			name: "full cluster conversion",
			apiCluster: api.V0043ClusterRec{
				Name: ptrString("test-cluster"),
				Controller: &struct {
					Host *string `json:"host,omitempty"`
					Port *int32  `json:"port,omitempty"`
				}{
					Host: ptrString("controller.example.com"),
					Port: ptrInt32(6817),
				},
				Nodes:        ptrString("node[1-10]"),
				RpcVersion:   ptrInt32(21),
				SelectPlugin: ptrString("select/cons_tres"),
				Flags:        &[]api.V0043ClusterRecFlags{"EXTERNAL", "FEDERATION"},
				Tres: &[]api.V0043Tres{
					{
						Type:  "cpu",
						Name:  ptrString("cpu"),
						Id:    ptrInt32(1),
						Count: ptrInt64(100),
					},
				},
				Associations: &struct {
					Root *api.V0043AssocShort `json:"root,omitempty"`
				}{
					Root: &api.V0043AssocShort{
						Account:   ptrString("root"),
						Cluster:   ptrString("test-cluster"),
						Partition: ptrString("normal"),
						User:      "root",
					},
				},
			},
			expected: &types.Cluster{
				Name:           "test-cluster",
				ControllerHost: "controller.example.com",
				ControllerPort: 6817,
				Nodes:          "node[1-10]",
				RpcVersion:     21,
				SelectPlugin:   "select/cons_tres",
				Flags:          []string{"EXTERNAL", "FEDERATION"},
				TRES: []types.TRES{
					{
						Type:  "cpu",
						Name:  "cpu",
						ID:    1,
						Count: 100,
					},
				},
				Associations: &types.AssociationShort{
					Root: &types.AssocShort{
						Account:   "root",
						Cluster:   "test-cluster",
						Partition: "normal",
						User:      "root",
					},
				},
				Meta: make(map[string]interface{}),
			},
		},
		{
			name: "minimal cluster conversion",
			apiCluster: api.V0043ClusterRec{
				Name: ptrString("minimal-cluster"),
			},
			expected: &types.Cluster{
				Name: "minimal-cluster",
				Meta: make(map[string]interface{}),
			},
		},
		{
			name:       "empty cluster",
			apiCluster: api.V0043ClusterRec{},
			expected: &types.Cluster{
				Meta: make(map[string]interface{}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.convertAPIClusterToCommon(tt.apiCluster)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Name, result.Name)
				assert.Equal(t, tt.expected.ControllerHost, result.ControllerHost)
				assert.Equal(t, tt.expected.ControllerPort, result.ControllerPort)
				assert.Equal(t, tt.expected.Nodes, result.Nodes)
				assert.Equal(t, tt.expected.RpcVersion, result.RpcVersion)
				assert.Equal(t, tt.expected.Flags, result.Flags)
				assert.Equal(t, len(tt.expected.TRES), len(result.TRES))
			}
		})
	}
}

func TestClusterAdapter_extractMeta(t *testing.T) {
	adapter := NewClusterAdapter(nil)

	tests := []struct {
		name     string
		meta     *api.V0043OpenapiMeta
		expected map[string]interface{}
	}{
		{
			name: "full meta",
			meta: &api.V0043OpenapiMeta{
				Client: &struct {
					Group  *string `json:"group,omitempty"`
					Source *string `json:"source,omitempty"`
					User   *string `json:"user,omitempty"`
				}{
					Source: ptrString("slurmrestd"),
					User:   ptrString("slurm"),
					Group:  ptrString("slurm"),
				},
				Plugin: &struct {
					AccountingStorage *string `json:"accounting_storage,omitempty"`
					DataParser        *string `json:"data_parser,omitempty"`
					Name              *string `json:"name,omitempty"`
				}{
					AccountingStorage: ptrString("accounting_storage/mysql"),
					DataParser:        ptrString("data_parser/v0.0.43"),
				},
			},
			expected: map[string]interface{}{
				"client": map[string]interface{}{
					"source": "slurmrestd",
					"user":   "slurm",
					"group":  "slurm",
				},
				"plugin": map[string]interface{}{
					"accounting_storage": "accounting_storage/mysql",
					"data_parser":        "data_parser/v0.0.43",
				},
			},
		},
		{
			name:     "nil meta",
			meta:     nil,
			expected: map[string]interface{}{},
		},
		{
			name:     "empty meta",
			meta:     &api.V0043OpenapiMeta{},
			expected: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.extractMeta(tt.meta)
			assert.Equal(t, tt.expected, result)
		})
	}
}

