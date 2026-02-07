// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_40

import (
	"context"
	"testing"

	adapterbase "github.com/jontk/slurm-client/internal/adapters/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClusterAdapter_ValidateContext(t *testing.T) {
	adapter := &ClusterAdapter{
		BaseManager: adapterbase.NewBaseManager("v0.0.40", "Cluster"),
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
func TestClusterAdapter_ValidateResourceName(t *testing.T) {
	adapter := &ClusterAdapter{
		BaseManager: adapterbase.NewBaseManager("v0.0.40", "Cluster"),
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
			errMsg:  "cluster name is required",
		},
		{
			name:    "valid name",
			resName: "cluster1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateResourceName(tt.resName, "cluster name")
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
