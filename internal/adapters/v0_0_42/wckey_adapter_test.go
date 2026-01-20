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

func TestNewWCKeyAdapter(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewWCKeyAdapter(client)

	assert.NotNil(t, adapter)
	assert.Equal(t, client, adapter.client)
	assert.NotNil(t, adapter.BaseManager)
	assert.Equal(t, "v0.0.42", adapter.GetVersion())
}

func TestWCKeyAdapter_ValidateContext(t *testing.T) {
	adapter := &WCKeyAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "WCKey"),
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

func TestWCKeyAdapter_ClientValidation(t *testing.T) {
	// Test nil client validation
	adapter := NewWCKeyAdapter(nil)
	ctx := context.Background()

	_, err := adapter.List(ctx, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	_, err = adapter.Get(ctx, "test-wckey")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	err = adapter.Delete(ctx, "test-wckey")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	// Test that non-nil client passes initial validation
	validAdapter := NewWCKeyAdapter(&api.ClientWithResponses{})
	assert.NotNil(t, validAdapter.client)
}

func TestWCKeyAdapter_ListOptionsHandling(t *testing.T) {
	adapter := NewWCKeyAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	tests := []struct {
		name string
		opts *types.WCKeyListOptions
	}{
		{
			name: "nil options",
			opts: nil,
		},
		{
			name: "empty options",
			opts: &types.WCKeyListOptions{},
		},
		{
			name: "options with users",
			opts: &types.WCKeyListOptions{
				Users: []string{"user1", "user2"},
			},
		},
		{
			name: "options with clusters",
			opts: &types.WCKeyListOptions{
				Clusters: []string{"cluster1", "cluster2"},
			},
		},
		{
			name: "options with names",
			opts: &types.WCKeyListOptions{
				Names: []string{"default", "project1"},
			},
		},
		{
			name: "options with flags",
			opts: &types.WCKeyListOptions{
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

func TestWCKeyAdapter_GetByID(t *testing.T) {
	adapter := NewWCKeyAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	tests := []struct {
		name    string
		wckeyID string
	}{
		{
			name:    "valid ID",
			wckeyID: "default",
		},
		{
			name:    "empty ID",
			wckeyID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := adapter.Get(ctx, tt.wckeyID)
			// Should get client validation error
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "client")
		})
	}
}

func TestWCKeyAdapter_ConvertAPIWCKeyToCommon(t *testing.T) {
	adapter := NewWCKeyAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name     string
		apiWCKey api.V0042Wckey
		expected types.WCKey
	}{
		{
			name: "basic WCKey",
			apiWCKey: api.V0042Wckey{
				Id:      ptrInt32(1),
				Name:    "default",
				Cluster: "cluster1",
				User:    "user1",
			},
			expected: types.WCKey{
				ID:      "1",
				Name:    "default",
				Cluster: "cluster1",
				User:    "user1",
			},
		},
		{
			name: "WCKey with accounting",
			apiWCKey: api.V0042Wckey{
				Id:      ptrInt32(2),
				Name:    "project1",
				Cluster: "cluster2",
				User:    "user2",
			},
			expected: types.WCKey{
				ID:      "2",
				Name:    "project1",
				Cluster: "cluster2",
				User:    "user2",
			},
		},
		{
			name: "minimal WCKey",
			apiWCKey: api.V0042Wckey{
				Name: "minimal",
			},
			expected: types.WCKey{
				Name: "minimal",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.convertAPIWCKeyToCommon(tt.apiWCKey)

			require.NoError(t, err)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Cluster, result.Cluster)
			assert.Equal(t, tt.expected.User, result.User)
		})
	}
}

func TestWCKeyAdapter_ErrorHandling(t *testing.T) {
	adapter := NewWCKeyAdapter(nil)
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
				_, err := adapter.Get(ctx, "default")
				return err
			},
		},
		{
			name: "Create with nil client",
			testFunc: func() error {
				_, err := adapter.Create(ctx, &types.WCKeyCreate{
					Name:    "test",
					Cluster: "cluster1",
					User:    "user1",
				})
				return err
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
					strings.Contains(errorMsg, "not implemented"),
				"Expected client validation or not implemented error, got: %v", err)
		})
	}
}

func TestWCKeyAdapter_CreateValidation(t *testing.T) {
	adapter := NewWCKeyAdapter(nil) // Use nil client to test validation
	ctx := context.Background()

	tests := []struct {
		name  string
		wckey *types.WCKeyCreate
	}{
		{
			name: "valid create",
			wckey: &types.WCKeyCreate{
				Name:    "new-wckey",
				Cluster: "cluster1",
				User:    "user1",
			},
		},
		{
			name: "minimal create",
			wckey: &types.WCKeyCreate{
				Name: "minimal",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := adapter.Create(ctx, tt.wckey)
			// Should get client validation error
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "client")
		})
	}
}
