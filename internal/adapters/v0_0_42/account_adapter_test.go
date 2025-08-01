// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountAdapter_ValidateAccountCreate(t *testing.T) {
	adapter := &AccountAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Account"),
	}

	tests := []struct {
		name    string
		account *types.AccountCreate
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil account",
			account: nil,
			wantErr: true,
			errMsg:  "account data is required",
		},
		{
			name: "empty name",
			account: &types.AccountCreate{
				Name: "",
			},
			wantErr: true,
			errMsg:  "account name is required",
		},
		{
			name: "valid basic account",
			account: &types.AccountCreate{
				Name:        "test-account",
				Description: "Test account",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateAccountCreate(tt.account)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAccountAdapter_FilterAccountList(t *testing.T) {
	adapter := &AccountAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Account"),
	}

	accounts := []types.Account{
		{
			Name:         "account1",
			Description:  "Account 1",
			Organization: "Org A",
		},
		{
			Name:         "account2",
			Description:  "Account 2",
			Organization: "Org B",
		},
	}

	tests := []struct {
		name     string
		opts     *types.AccountListOptions
		expected []string // expected account names
	}{
		{
			name:     "no filters",
			opts:     &types.AccountListOptions{},
			expected: []string{"account1", "account2"},
		},
		{
			name: "filter by names",
			opts: &types.AccountListOptions{
				Names: []string{"account1"},
			},
			expected: []string{"account1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.FilterAccountList(accounts, tt.opts)
			resultNames := make([]string, len(result))
			for i, account := range result {
				resultNames[i] = account.Name
			}
			assert.Equal(t, tt.expected, resultNames)
		})
	}
}

func TestAccountAdapter_ValidateContext(t *testing.T) {
	adapter := &AccountAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Account"),
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
