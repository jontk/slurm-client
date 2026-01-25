// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
)

func TestNewAssociationAdapter(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewAssociationAdapter(client)

	assert.NotNil(t, adapter)
	assert.Equal(t, client, adapter.client)
	assert.NotNil(t, adapter.BaseManager)
	assert.Equal(t, "v0.0.42", adapter.GetVersion())
}

func TestAssociationAdapter_ValidateContext(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Association"),
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

func TestAssociationAdapter_ClientValidation(t *testing.T) {
	// Test nil client validation
	adapter := NewAssociationAdapter(nil)
	ctx := context.Background()

	_, err := adapter.List(ctx, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	_, err = adapter.Get(ctx, "test-association")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	err = adapter.Delete(ctx, "test-association")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	// Test that non-nil client passes initial validation
	validAdapter := NewAssociationAdapter(&api.ClientWithResponses{})
	assert.NotNil(t, validAdapter.client)
}

func TestAssociationAdapter_ListOptionsHandling(t *testing.T) {
	adapter := NewAssociationAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	tests := []struct {
		name string
		opts *types.AssociationListOptions
	}{
		{
			name: "nil options",
			opts: nil,
		},
		{
			name: "empty options",
			opts: &types.AssociationListOptions{},
		},
		{
			name: "options with accounts",
			opts: &types.AssociationListOptions{
				Accounts: []string{"account1", "account2"},
			},
		},
		{
			name: "options with users",
			opts: &types.AssociationListOptions{
				Users: []string{"user1", "user2"},
			},
		},
		{
			name: "options with clusters",
			opts: &types.AssociationListOptions{
				Clusters: []string{"cluster1", "cluster2"},
			},
		},
		{
			name: "options with partitions",
			opts: &types.AssociationListOptions{
				Partitions: []string{"partition1", "partition2"},
			},
		},
		{
			name: "options with flags",
			opts: &types.AssociationListOptions{
				OnlyDefaults: true,
				WithDeleted:  true,
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

func TestAssociationAdapter_GetByID(t *testing.T) {
	adapter := NewAssociationAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	tests := []struct {
		name          string
		associationID string
	}{
		{
			name:          "valid ID",
			associationID: "12345",
		},
		{
			name:          "empty ID",
			associationID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := adapter.Get(ctx, tt.associationID)
			// Should get client validation error
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "client")
		})
	}
}

func TestAssociationAdapter_ConvertAPIAssociationToCommon(t *testing.T) {
	adapter := NewAssociationAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name           string
		apiAssociation api.V0042Assoc
		expected       types.Association
	}{
		{
			name: "basic association",
			apiAssociation: api.V0042Assoc{
				Account:   ptrString("account1"),
				User:      "user1",
				Cluster:   ptrString("cluster1"),
				SharesRaw: ptrInt32(1000),
			},
			expected: types.Association{
				AccountName: "account1",
				UserName:    "user1",
				Cluster:     "cluster1",
				SharesRaw:   1000,
			},
		},
		{
			name: "association with parent",
			apiAssociation: api.V0042Assoc{
				Account:       ptrString("subaccount"),
				User:          "user2",
				Cluster:       ptrString("cluster1"),
				ParentAccount: ptrString("parentaccount"),
				SharesRaw:     ptrInt32(500),
			},
			expected: types.Association{
				AccountName:   "subaccount",
				UserName:      "user2",
				Cluster:       "cluster1",
				ParentAccount: "parentaccount",
				SharesRaw:     500,
			},
		},
		{
			name: "minimal association",
			apiAssociation: api.V0042Assoc{
				Account: ptrString("minimal"),
			},
			expected: types.Association{
				AccountName: "minimal",
			},
		},
		{
			name: "association with default QoS",
			apiAssociation: api.V0042Assoc{
				Account: ptrString("account3"),
				User:    "user3",
				Default: &struct {
					Qos *string `json:"qos,omitempty"`
				}{
					Qos: ptrString("high"),
				},
			},
			expected: types.Association{
				AccountName: "account3",
				UserName:    "user3",
				DefaultQoS:  "high",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.convertAPIAssociationToCommon(tt.apiAssociation)
			assert.Equal(t, tt.expected.AccountName, result.AccountName)
			assert.Equal(t, tt.expected.UserName, result.UserName)
			assert.Equal(t, tt.expected.Cluster, result.Cluster)
			assert.Equal(t, tt.expected.ParentAccount, result.ParentAccount)
			assert.Equal(t, tt.expected.SharesRaw, result.SharesRaw)
			assert.Equal(t, tt.expected.DefaultQoS, result.DefaultQoS)
		})
	}
}

func TestAssociationAdapter_ErrorHandling(t *testing.T) {
	adapter := NewAssociationAdapter(nil)
	ctx := context.Background()

	// Test various error conditions with nil client
	tests := []struct {
		name     string
		testFunc func() error
	}{
		{
			name: "List with nil client",
			testFunc: func() error {
				_, err := adapter.List(ctx, nil)
				return err
			},
		},
		{
			name: "Get with nil client",
			testFunc: func() error {
				_, err := adapter.Get(ctx, "12345")
				return err
			},
		},
		{
			name: "Create with nil client",
			testFunc: func() error {
				_, err := adapter.Create(ctx, &types.AssociationCreate{
					AccountName: "test",
					UserName:    "user1",
					Cluster:     "cluster1",
				})
				return err
			},
		},
		{
			name: "Update with nil client",
			testFunc: func() error {
				return adapter.Update(ctx, "test", &types.AssociationUpdate{})
			},
		},
		{
			name: "Delete with nil client",
			testFunc: func() error {
				return adapter.Delete(ctx, "test")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testFunc()
			assert.Error(t, err)
			// Should get either client validation error or "not implemented" error
			errorMsg := err.Error()
			assert.True(t,
				strings.Contains(errorMsg, "client") ||
					strings.Contains(errorMsg, "not implemented") ||
					strings.Contains(errorMsg, "not supported"),
				"Expected client validation or not implemented error, got: %v", err)
		})
	}
}

func TestAssociationAdapter_CreateValidation(t *testing.T) {
	adapter := NewAssociationAdapter(nil) // Use nil client to test validation
	ctx := context.Background()

	tests := []struct {
		name        string
		association *types.AssociationCreate
	}{
		{
			name: "valid create",
			association: &types.AssociationCreate{
				AccountName: "new-account",
				UserName:    "user1",
				Cluster:     "cluster1",
			},
		},
		{
			name: "minimal create",
			association: &types.AssociationCreate{
				AccountName: "minimal",
				Cluster:     "cluster1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := adapter.Create(ctx, tt.association)
			// Should get client validation error
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "client")
		})
	}
}

func TestAssociationAdapter_UpdateValidation(t *testing.T) {
	adapter := NewAssociationAdapter(nil) // Use nil client to test validation
	ctx := context.Background()

	update := &types.AssociationUpdate{
		SharesRaw: ptrInt32(2000),
	}

	err := adapter.Update(ctx, "test-association", update)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")
}
