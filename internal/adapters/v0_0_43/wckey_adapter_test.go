// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	"net/http"
	"testing"

	"github.com/jontk/slurm-client/internal/api/v0_0_43"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWCKeyAdapter_Create(t *testing.T) {
	adapter := &WCKeyAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "WCKey"),
	}

	tests := []struct {
		name    string
		ctx     context.Context
		wckey   *types.WCKeyCreate
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil context",
			ctx:     nil,
			wckey:   &types.WCKeyCreate{Name: "test", User: "user1", Cluster: "cluster1"},
			wantErr: true,
			errMsg:  "[VALIDATION_FAILED] context is required",
		},
		{
			name:    "nil wckey",
			ctx:     context.Background(),
			wckey:   nil,
			wantErr: true,
			errMsg:  "WCKey creation data is required",
		},
		{
			name: "empty name",
			ctx:  context.Background(),
			wckey: &types.WCKeyCreate{
				Name:    "",
				User:    "user1",
				Cluster: "cluster1",
			},
			wantErr: true,
			errMsg:  "WCKey name is required",
		},
		{
			name: "empty user",
			ctx:  context.Background(),
			wckey: &types.WCKeyCreate{
				Name:    "test",
				User:    "",
				Cluster: "cluster1",
			},
			wantErr: true,
			errMsg:  "API client not initialized",
		},
		{
			name: "empty cluster",
			ctx:  context.Background(),
			wckey: &types.WCKeyCreate{
				Name:    "test",
				User:    "user1",
				Cluster: "",
			},
			wantErr: true,
			errMsg:  "cluster is required for WCKey creation",
		},
		{
			name:    "nil client",
			ctx:     context.Background(),
			wckey:   &types.WCKeyCreate{Name: "test", User: "user1", Cluster: "cluster1"},
			wantErr: true,
			errMsg:  "API client not initialized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := adapter.Create(tt.ctx, tt.wckey)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestWCKeyAdapter_List(t *testing.T) {
	adapter := &WCKeyAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "WCKey"),
	}

	tests := []struct {
		name    string
		ctx     context.Context
		opts    *types.WCKeyListOptions
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil context",
			ctx:     nil,
			opts:    nil,
			wantErr: true,
			errMsg:  "[VALIDATION_FAILED] context is required",
		},
		{
			name:    "nil client",
			ctx:     context.Background(),
			opts:    nil,
			wantErr: true,
			errMsg:  "API client not initialized",
		},
		{
			name: "with filters",
			ctx:  context.Background(),
			opts: &types.WCKeyListOptions{
				Names:    []string{"key1", "key2"},
				Users:    []string{"user1", "user2"},
				Clusters: []string{"cluster1"},
			},
			wantErr: true,
			errMsg:  "API client not initialized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := adapter.List(tt.ctx, tt.opts)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestWCKeyAdapter_Delete(t *testing.T) {
	adapter := &WCKeyAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "WCKey"),
	}

	tests := []struct {
		name    string
		ctx     context.Context
		id      string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil context",
			ctx:     nil,
			id:      "test_user1_cluster1",
			wantErr: true,
			errMsg:  "[VALIDATION_FAILED] context is required",
		},
		{
			name:    "empty id",
			ctx:     context.Background(),
			id:      "",
			wantErr: true,
			errMsg:  "[VALIDATION_FAILED] WCKey ID cannot be empty",
		},
		{
			name:    "nil client",
			ctx:     context.Background(),
			id:      "test_user1_cluster1",
			wantErr: true,
			errMsg:  "API client not initialized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.Delete(tt.ctx, tt.id)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// MockWCKeyClient is a mock implementation for testing
type MockWCKeyClient struct {
	CreateFunc func(ctx context.Context, body v0_0_43.V0043OpenapiWckeyResp, reqEditors ...v0_0_43.RequestEditorFn) (*http.Response, error)
	ListFunc   func(ctx context.Context, params *v0_0_43.SlurmdbV0043GetWckeysParams, reqEditors ...v0_0_43.RequestEditorFn) (*http.Response, error)
	DeleteFunc func(ctx context.Context, wckey string, reqEditors ...v0_0_43.RequestEditorFn) (*http.Response, error)
}

func TestWCKeyAdapter_ConvertAPIWCKeyToCommon(t *testing.T) {
	adapter := &WCKeyAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "WCKey"),
	}

	tests := []struct {
		name    string
		apiKey  v0_0_43.V0043Wckey
		want    types.WCKey
		wantErr bool
	}{
		{
			name: "basic conversion",
			apiKey: v0_0_43.V0043Wckey{
				Name:    "test-key",
				User:    "user1",
				Cluster: "cluster1",
			},
			want: types.WCKey{
				Name:    "test-key",
				User:    "user1",
				Cluster: "cluster1",
			},
			wantErr: false,
		},
		{
			name: "with ID",
			apiKey: v0_0_43.V0043Wckey{
				Id:      int32Ptr(123),
				Name:    "test-key",
				User:    "user1",
				Cluster: "cluster1",
			},
			want: types.WCKey{
				ID:      "123",
				Name:    "test-key",
				User:    "user1",
				Cluster: "cluster1",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adapter.convertAPIWCKeyToCommon(tt.apiKey)
			if tt.wantErr {
				// This function doesn't return an error in the actual implementation
				t.Skip("Function doesn't return error")
			} else {
				assert.Equal(t, tt.want.Name, got.Name)
				assert.Equal(t, tt.want.User, got.User)
				assert.Equal(t, tt.want.Cluster, got.Cluster)
				assert.Equal(t, tt.want.ID, got.ID)
			}
		})
	}
}

// Helper function for int32 pointers
func int32Ptr(i int32) *int32 {
	return &i
}