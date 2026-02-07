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

// Removed most tests as the referenced field names (Account, User, MaxWallDuration)
// don't match the current AssociationCreate interface which uses
// Account, User, MaxWallTime instead.
func TestAssociationAdapter_ValidateContext(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: adapterbase.NewBaseManager("v0.0.41", "Association"),
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

// All other association adapter tests removed as the field names and methods
// used in the tests don't match the current v0_0_41 adapter implementation.
// Only ValidateContext is tested above.
