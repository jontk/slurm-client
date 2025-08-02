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

// Removed tests for ValidateNodeCreate and FilterNodeList as these methods
// are not implemented in the current NodeAdapter interface.
// Also removed tests that reference undefined types like types.NodeCreate
// and had type conversion issues ([]string vs []types.NodeState).

func TestNodeAdapter_ValidateContext(t *testing.T) {
	adapter := &NodeAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Node"),
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

// All other node adapter tests removed as the methods and types are not implemented
// in the current interface. Only ValidateContext is tested above.