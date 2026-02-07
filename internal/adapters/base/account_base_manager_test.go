// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package base

import (
	"testing"

	types "github.com/jontk/slurm-client/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountBaseManager_New(t *testing.T) {
	manager := NewAccountBaseManager("v0.0.43")
	assert.NotNil(t, manager)
	assert.Equal(t, "v0.0.43", manager.GetVersion())
	assert.Equal(t, "Account", manager.GetResourceType())
}
func TestAccountBaseManager_ValidateAccountCreate(t *testing.T) {
	manager := NewAccountBaseManager("v0.0.43")
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
			errMsg:  "data is required",
		},
		{
			name: "empty name",
			account: &types.AccountCreate{
				Name: "",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "valid account",
			account: &types.AccountCreate{
				Name:         "test-account",
				Organization: "Test Org",
				Description:  "Test description",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateAccountCreate(tt.account)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
func TestAccountBaseManager_ValidateAccountUpdate(t *testing.T) {
	manager := NewAccountBaseManager("v0.0.43")
	tests := []struct {
		name    string
		update  *types.AccountUpdate
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
			update: &types.AccountUpdate{
				Description: stringPtr("New description"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateAccountUpdate(tt.update)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
func TestAccountBaseManager_ApplyAccountDefaults(t *testing.T) {
	manager := NewAccountBaseManager("v0.0.43")
	account := &types.AccountCreate{
		Name: "test-account",
	}
	result := manager.ApplyAccountDefaults(account)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Coordinators)
	assert.NotNil(t, result.QoSList)
	assert.NotNil(t, result.AllowedPartitions)
	assert.NotNil(t, result.GrpTRES)
	assert.NotNil(t, result.MaxTRES)
}
func TestAccountBaseManager_FilterAccountList(t *testing.T) {
	manager := NewAccountBaseManager("v0.0.43")
	accounts := []types.Account{
		{Name: "account1", Organization: "Org A"},
		{Name: "account2", Organization: "Org B"},
		{Name: "account3", Organization: "Org A"},
	}
	tests := []struct {
		name     string
		opts     *types.AccountListOptions
		expected []string
	}{
		{
			name:     "no filters",
			opts:     nil,
			expected: []string{"account1", "account2", "account3"},
		},
		{
			name: "filter by name",
			opts: &types.AccountListOptions{
				Names: []string{"account1"},
			},
			expected: []string{"account1"},
		},
		{
			name: "filter by organization",
			opts: &types.AccountListOptions{
				Organizations: []string{"Org A"},
			},
			expected: []string{"account1", "account3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.FilterAccountList(accounts, tt.opts)
			resultNames := make([]string, len(result))
			for i, account := range result {
				resultNames[i] = account.Name
			}
			assert.Equal(t, tt.expected, resultNames)
		})
	}
}
