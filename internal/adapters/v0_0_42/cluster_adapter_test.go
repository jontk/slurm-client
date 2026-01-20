// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
)

func TestNewClusterAdapter(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewClusterAdapter(client)

	assert.NotNil(t, adapter)
	assert.Equal(t, client, adapter.client)
	assert.NotNil(t, adapter.BaseManager)
	assert.Equal(t, "v0.0.42", adapter.GetVersion())
}

func TestClusterAdapter_ValidateContext(t *testing.T) {
	adapter := &ClusterAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Cluster"),
		client:      &api.ClientWithResponses{},
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
			errMsg:  "context",
		},
		{
			name:    "valid context",
			ctx:     context.Background(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateContext(tt.ctx)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestClusterAdapter_ClientValidation(t *testing.T) {
	// Test nil client validation
	adapter := NewClusterAdapter(nil)
	ctx := context.Background()

	_, err := adapter.List(ctx, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	// Test that non-nil client passes initial validation
	validAdapter := NewClusterAdapter(&api.ClientWithResponses{})
	assert.NotNil(t, validAdapter.client)
}

func TestClusterAdapter_ConvertAPIClusterToCommon(t *testing.T) {
	adapter := NewClusterAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name        string
		apiCluster  api.V0042ClusterRec
		expected    types.Cluster
		expectError bool
	}{
		{
			name: "basic conversion",
			apiCluster: api.V0042ClusterRec{
				Name: ptrString("test-cluster"),
				Controller: &struct {
					Host *string `json:"host,omitempty"`
					Port *int32  `json:"port,omitempty"`
				}{
					Host: ptrString("controller.example.com"),
					Port: ptrInt32(6817),
				},
				Nodes:        ptrString("node[1-10]"),
				RpcVersion:   ptrInt32(40),
				SelectPlugin: ptrString("select/cons_tres"),
			},
			expected: types.Cluster{
				Name:           "test-cluster",
				ControllerHost: "controller.example.com",
				ControllerPort: 6817,
				Nodes:          "node[1-10]",
				RpcVersion:     40,
				SelectPlugin:   "select/cons_tres",
				Meta:           make(map[string]interface{}),
			},
			expectError: false,
		},
		{
			name: "minimal conversion",
			apiCluster: api.V0042ClusterRec{
				Name: ptrString("minimal-cluster"),
			},
			expected: types.Cluster{
				Name: "minimal-cluster",
				Meta: make(map[string]interface{}),
			},
			expectError: false,
		},
		{
			name: "with TRES data",
			apiCluster: api.V0042ClusterRec{
				Name: ptrString("tres-cluster"),
				Tres: &[]api.V0042Tres{
					{
						Type:  "cpu",
						Id:    ptrInt32(1),
						Name:  ptrString("cpu"),
						Count: ptrInt64(1000),
					},
					{
						Type:  "mem",
						Id:    ptrInt32(2),
						Name:  ptrString("memory"),
						Count: ptrInt64(1024000),
					},
				},
			},
			expected: types.Cluster{
				Name: "tres-cluster",
				TRES: []types.TRES{
					{
						Type:  "cpu",
						ID:    1,
						Name:  "cpu",
						Count: 1000,
					},
					{
						Type:  "mem",
						ID:    2,
						Name:  "memory",
						Count: 1024000,
					},
				},
				Meta: make(map[string]interface{}),
			},
			expectError: false,
		},
		{
			name: "with associations",
			apiCluster: api.V0042ClusterRec{
				Name: ptrString("assoc-cluster"),
				Associations: &struct {
					Root *api.V0042AssocShort `json:"root,omitempty"`
				}{
					Root: &api.V0042AssocShort{
						User:      "root",
						Account:   ptrString("root"),
						Cluster:   ptrString("assoc-cluster"),
						Partition: ptrString("main"),
					},
				},
			},
			expected: types.Cluster{
				Name: "assoc-cluster",
				Associations: &types.AssociationShort{
					Root: &types.AssocShort{
						User:      "root",
						Account:   "root",
						Cluster:   "assoc-cluster",
						Partition: "main",
					},
				},
				Meta: make(map[string]interface{}),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.convertAPIClusterToCommon(tt.apiCluster)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected.Name, result.Name)
				assert.Equal(t, tt.expected.ControllerHost, result.ControllerHost)
				assert.Equal(t, tt.expected.ControllerPort, result.ControllerPort)
				assert.Equal(t, tt.expected.Nodes, result.Nodes)
				assert.Equal(t, tt.expected.RpcVersion, result.RpcVersion)
				assert.Equal(t, tt.expected.SelectPlugin, result.SelectPlugin)

				// Check TRES if expected
				if len(tt.expected.TRES) > 0 {
					assert.Len(t, result.TRES, len(tt.expected.TRES))
					for i, expectedTres := range tt.expected.TRES {
						assert.Equal(t, expectedTres.Type, result.TRES[i].Type)
						assert.Equal(t, expectedTres.ID, result.TRES[i].ID)
						assert.Equal(t, expectedTres.Name, result.TRES[i].Name)
						assert.Equal(t, expectedTres.Count, result.TRES[i].Count)
					}
				}

				// Check associations if expected
				if tt.expected.Associations != nil {
					assert.NotNil(t, result.Associations)
					assert.NotNil(t, result.Associations.Root)
					assert.Equal(t, tt.expected.Associations.Root.User, result.Associations.Root.User)
					assert.Equal(t, tt.expected.Associations.Root.Account, result.Associations.Root.Account)
					assert.Equal(t, tt.expected.Associations.Root.Cluster, result.Associations.Root.Cluster)
					assert.Equal(t, tt.expected.Associations.Root.Partition, result.Associations.Root.Partition)
				}
			}
		})
	}
}

func TestClusterAdapter_ExtractMeta(t *testing.T) {
	adapter := NewClusterAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name     string
		meta     *api.V0042OpenapiMeta
		expected map[string]interface{}
	}{
		{
			name:     "nil meta",
			meta:     nil,
			expected: map[string]interface{}{},
		},
		{
			name:     "empty meta",
			meta:     &api.V0042OpenapiMeta{},
			expected: map[string]interface{}{},
		},
		{
			name: "with client info",
			meta: &api.V0042OpenapiMeta{
				Client: &struct {
					Group  *string `json:"group,omitempty"`
					Source *string `json:"source,omitempty"`
					User   *string `json:"user,omitempty"`
				}{
					Source: ptrString("slurm-client"),
					User:   ptrString("testuser"),
					Group:  ptrString("testgroup"),
				},
			},
			expected: map[string]interface{}{
				"client": map[string]interface{}{
					"source": "slurm-client",
					"user":   "testuser",
					"group":  "testgroup",
				},
			},
		},
		{
			name: "with plugin info",
			meta: &api.V0042OpenapiMeta{
				Plugin: &struct {
					AccountingStorage *string `json:"accounting_storage,omitempty"`
					DataParser        *string `json:"data_parser,omitempty"`
					Name              *string `json:"name,omitempty"`
					Type              *string `json:"type,omitempty"`
				}{
					AccountingStorage: ptrString("accounting_storage/slurmdbd"),
				},
			},
			expected: map[string]interface{}{
				"plugin": map[string]interface{}{
					"accounting_storage": "accounting_storage/slurmdbd",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.extractMeta(tt.meta)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestClusterAdapter_ListOptionsHandling(t *testing.T) {
	adapter := NewClusterAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	tests := []struct {
		name string
		opts *types.ClusterListOptions
	}{
		{
			name: "nil options",
			opts: nil,
		},
		{
			name: "empty options",
			opts: &types.ClusterListOptions{},
		},
		{
			name: "options with update time",
			opts: &types.ClusterListOptions{
				UpdateTime: &time.Time{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := adapter.List(ctx, tt.opts)
			// Should get client validation error before any option processing
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "client")
		})
	}
}

func TestClusterAdapter_GetByName(t *testing.T) {
	adapter := NewClusterAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	tests := []struct {
		name        string
		clusterName string
	}{
		{
			name:        "valid name",
			clusterName: "test-cluster",
		},
		{
			name:        "empty name",
			clusterName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := adapter.Get(ctx, tt.clusterName)
			// Should get client validation error
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "client")
		})
	}
}

func TestClusterAdapter_CreateValidation(t *testing.T) {
	adapter := NewClusterAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	tests := []struct {
		name          string
		clusterCreate *types.ClusterCreate
	}{
		{
			name: "valid create",
			clusterCreate: &types.ClusterCreate{
				Name:           "test-cluster",
				ControllerHost: "controller.example.com",
				ControllerPort: 6817,
			},
		},
		{
			name: "minimal create",
			clusterCreate: &types.ClusterCreate{
				Name: "minimal-cluster",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := adapter.Create(ctx, tt.clusterCreate)
			// Should get client validation error
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "client")
		})
	}
}

func TestClusterAdapter_DeleteValidation(t *testing.T) {
	adapter := NewClusterAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	err := adapter.Delete(ctx, "test-cluster")
	// Should get client validation error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")
}
