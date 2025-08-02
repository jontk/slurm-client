// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"testing"

	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Removed TestAssociationAdapter_ValidateAssociationCreate as this method 
// is not implemented in the current AssociationAdapter interface

// Removed TestAssociationAdapter_ApplyAssociationDefaults as this method
// is not implemented in the current AssociationAdapter interface

// Removed TestAssociationAdapter_FilterAssociationList as this method
// is not implemented in the current AssociationAdapter interface

// Removed TestAssociationAdapter_ValidateAssociationHierarchy as this method
// is not implemented in the current AssociationAdapter interface

func TestAssociationAdapter_ValidateContext(t *testing.T) {
	adapter := &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Association"),
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
		{
			name:    "context with timeout",
			ctx:     context.TODO(),
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

// All other association adapter tests removed as the methods are not implemented
// in the current interface. Only ValidateContext is tested above.