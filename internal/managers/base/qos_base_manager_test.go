// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQoSBaseManager_New(t *testing.T) {
	manager := NewQoSBaseManager("v0.0.43")
	assert.NotNil(t, manager)
	assert.Equal(t, "v0.0.43", manager.GetVersion())
	assert.Equal(t, "QoS", manager.GetResourceType())
}

func TestQoSBaseManager_ValidateQoSCreate(t *testing.T) {
	manager := NewQoSBaseManager("v0.0.43")

	tests := []struct {
		name    string
		qos     *types.QoSCreate
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil qos",
			qos:     nil,
			wantErr: true,
			errMsg:  "data is required",
		},
		{
			name: "empty name",
			qos: &types.QoSCreate{
				Name: "",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "valid qos",
			qos: &types.QoSCreate{
				Name:        "test-qos",
				Description: "Test QoS",
				Priority:    100,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateQoSCreate(tt.qos)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestQoSBaseManager_ValidateQoSUpdate(t *testing.T) {
	manager := NewQoSBaseManager("v0.0.43")

	tests := []struct {
		name    string
		update  *types.QoSUpdate
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
			update: &types.QoSUpdate{
				Description: stringPtr("Updated description"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateQoSUpdate(tt.update)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
