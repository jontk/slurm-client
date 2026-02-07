// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"context"
	"testing"

	adapterbase "github.com/jontk/slurm-client/internal/adapters/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWCKeyAdapter_ValidateContext(t *testing.T) {
	adapter := &WCKeyAdapter{
		BaseManager: adapterbase.NewBaseManager("v0.0.41", "WCKey"),
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
func TestWCKeyAdapter_ValidateResourceName(t *testing.T) {
	adapter := &WCKeyAdapter{
		BaseManager: adapterbase.NewBaseManager("v0.0.41", "WCKey"),
	}
	tests := []struct {
		name    string
		resName string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty name",
			resName: "",
			wantErr: true,
			errMsg:  "wckey name is required",
		},
		{
			name:    "valid name",
			resName: "wckey1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateResourceName(tt.resName, "wckey name")
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
