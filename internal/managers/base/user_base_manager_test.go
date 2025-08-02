// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserBaseManager_New(t *testing.T) {
	manager := NewUserBaseManager("v0.0.43")
	assert.NotNil(t, manager)
	assert.Equal(t, "v0.0.43", manager.GetVersion())
	assert.Equal(t, "User", manager.GetResourceType())
}

func TestUserBaseManager_ValidateUserCreate(t *testing.T) {
	manager := NewUserBaseManager("v0.0.43")

	tests := []struct {
		name    string
		user    *types.UserCreate
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil user",
			user:    nil,
			wantErr: true,
			errMsg:  "data is required",
		},
		{
			name: "empty name",
			user: &types.UserCreate{
				Name: "",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "valid user",
			user: &types.UserCreate{
				Name:       "testuser",
				DefaultWCKey: "default",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateUserCreate(tt.user)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserBaseManager_ValidateUserUpdate(t *testing.T) {
	manager := NewUserBaseManager("v0.0.43")

	tests := []struct {
		name    string
		update  *types.UserUpdate
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
			update: &types.UserUpdate{
				DefaultWCKey: stringPtr("new-default"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateUserUpdate(tt.update)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}