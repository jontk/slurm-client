// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssociationBaseManager_New(t *testing.T) {
	manager := NewAssociationBaseManager("v0.0.43")
	assert.NotNil(t, manager)
	assert.Equal(t, "v0.0.43", manager.GetVersion())
	assert.Equal(t, "Association", manager.GetResourceType())
}

func TestAssociationBaseManager_ValidateAssociationCreate(t *testing.T) {
	manager := NewAssociationBaseManager("v0.0.43")

	tests := []struct {
		name        string
		association *types.AssociationCreate
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "nil association",
			association: nil,
			wantErr:     true,
			errMsg:      "data is required",
		},
		{
			name: "empty account name",
			association: &types.AssociationCreate{
				AccountName: "",
			},
			wantErr: true,
			errMsg:  "account name is required",
		},
		{
			name: "valid user association",
			association: &types.AssociationCreate{
				AccountName: "test-account",
				UserName:    "testuser",
				Cluster:     "test-cluster",
			},
			wantErr: false,
		},
		{
			name: "valid account association",
			association: &types.AssociationCreate{
				AccountName: "test-account",
				Cluster:     "test-cluster",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateAssociationCreate(tt.association)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAssociationBaseManager_ValidateAssociationUpdate(t *testing.T) {
	manager := NewAssociationBaseManager("v0.0.43")

	tests := []struct {
		name    string
		update  *types.AssociationUpdate
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil update",
			update:  nil,
			wantErr: true,
			errMsg:  "data is required",
		},
		{
			name: "valid update",
			update: &types.AssociationUpdate{
				SharesRaw: int32Ptr(100),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateAssociationUpdate(tt.update)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAssociationBaseManager_ApplyAssociationDefaults(t *testing.T) {
	manager := NewAssociationBaseManager("v0.0.43")

	association := &types.AssociationCreate{
		AccountName: "test-account",
		Cluster:     "test-cluster",
	}

	result := manager.ApplyAssociationDefaults(association)
	assert.NotNil(t, result)
	assert.NotNil(t, result.QoSList)
	assert.NotNil(t, result.GrpTRES)
	assert.NotNil(t, result.MaxTRES)
}

func TestAssociationBaseManager_FilterAssociationList(t *testing.T) {
	manager := NewAssociationBaseManager("v0.0.43")

	associations := []types.Association{
		{AccountName: "account1", UserName: "user1", Cluster: "cluster1"},
		{AccountName: "account2", UserName: "user2", Cluster: "cluster1"},
		{AccountName: "account1", UserName: "user1", Cluster: "cluster2"},
	}

	tests := []struct {
		name     string
		opts     *types.AssociationListOptions
		expected int
	}{
		{
			name:     "no filters",
			opts:     nil,
			expected: 3,
		},
		{
			name: "filter by account",
			opts: &types.AssociationListOptions{
				Accounts: []string{"account1"},
			},
			expected: 2,
		},
		{
			name: "filter by user",
			opts: &types.AssociationListOptions{
				Users: []string{"user1"},
			},
			expected: 2,
		},
		{
			name: "filter by cluster",
			opts: &types.AssociationListOptions{
				Clusters: []string{"cluster1"},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.FilterAssociationList(associations, tt.opts)
			assert.Len(t, result, tt.expected)
		})
	}
}
