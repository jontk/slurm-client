// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserAdapter_CreateAssociation(t *testing.T) {
	adapter := &UserAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "User"),
	}

	tests := []struct {
		name    string
		ctx     context.Context
		req     *types.UserAssociationRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil context",
			ctx:     nil,
			req:     &types.UserAssociationRequest{Users: []string{"user1"}, Cluster: "cluster1", Account: "acc1"},
			wantErr: true,
			errMsg:  "context is required",
		},
		{
			name:    "nil request",
			ctx:     context.Background(),
			req:     nil,
			wantErr: true,
			errMsg:  "association request is required",
		},
		{
			name: "empty users",
			ctx:  context.Background(),
			req: &types.UserAssociationRequest{
				Users:   []string{},
				Cluster: "cluster1",
				Account: "acc1",
			},
			wantErr: true,
			errMsg:  "at least one user is required",
		},
		{
			name: "empty cluster",
			ctx:  context.Background(),
			req: &types.UserAssociationRequest{
				Users:   []string{"user1"},
				Cluster: "",
				Account: "acc1",
			},
			wantErr: true,
			errMsg:  "cluster is required",
		},
		{
			name: "empty account",
			ctx:  context.Background(),
			req: &types.UserAssociationRequest{
				Users:   []string{"user1"},
				Cluster: "cluster1",
				Account: "",
			},
			wantErr: true,
			errMsg:  "account is required",
		},
		{
			name: "nil client",
			ctx:  context.Background(),
			req: &types.UserAssociationRequest{
				Users:   []string{"user1"},
				Cluster: "cluster1",
				Account: "acc1",
			},
			wantErr: true,
			errMsg:  "API client not initialized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := adapter.CreateAssociation(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, "success", resp.Status)
			}
		})
	}
}

func TestUserAdapter_ValidateContext(t *testing.T) {
	adapter := &UserAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "User"),
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
