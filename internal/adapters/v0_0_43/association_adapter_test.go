// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	"testing"

	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/stretchr/testify/assert"
)

// Helper functions are imported from test_helpers.go

func TestNewAssociationAdapter(t *testing.T) {
	adapter := NewAssociationAdapter(&api.ClientWithResponses{})
	assert.NotNil(t, adapter)
	assert.NotNil(t, adapter.BaseManager)
}

func TestAssociationAdapter_ValidateContext(t *testing.T) {
	adapter := NewAssociationAdapter(&api.ClientWithResponses{})

	// Test nil context
	//lint:ignore SA1012 intentionally testing nil context validation
	err := adapter.ValidateContext(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context is required")

	// Test valid context
	err = adapter.ValidateContext(context.Background())
	assert.NoError(t, err)
}

func TestAssociationAdapter_ConvertAPIAssociationToCommon(t *testing.T) {
	adapter := NewAssociationAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name           string
		apiAssociation api.V0043Assoc
		expectedUser   string
		expectedAcct   string
	}{
		{
			name: "full association",
			apiAssociation: api.V0043Assoc{
				Account:   ptrString("account1"),
				User:      "user1",
				Cluster:   ptrString("cluster1"),
				Partition: ptrString("normal"),
				Id:        ptrInt32(123),
			},
			expectedUser: "user1",
			expectedAcct: "account1",
		},
		{
			name: "minimal association",
			apiAssociation: api.V0043Assoc{
				Account: ptrString("account2"),
				User:    "user2",
			},
			expectedUser: "user2",
			expectedAcct: "account2",
		},
		{
			name:           "empty association",
			apiAssociation: api.V0043Assoc{},
			expectedUser:   "",
			expectedAcct:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.convertAPIAssociationToCommon(tt.apiAssociation)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedUser, result.UserName)
			assert.Equal(t, tt.expectedAcct, result.AccountName)
		})
	}
}

func TestAssociationAdapter_List(t *testing.T) {
	adapter := NewAssociationAdapter(nil) // Use nil client for testing validation logic

	// Test client initialization check (nil context validation is covered in TestAssociationAdapter_ValidateContext)
	_, err := adapter.List(context.TODO(), &types.AssociationListOptions{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
}

func TestAssociationAdapter_Get(t *testing.T) {
	adapter := NewAssociationAdapter(nil)

	// Test empty association ID
	_, err := adapter.Get(context.TODO(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "associationID is required")

	// Test client initialization check (nil context validation is covered in TestAssociationAdapter_ValidateContext)
	_, err = adapter.Get(context.TODO(), "test-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
}

func TestAssociationAdapter_Create(t *testing.T) {
	adapter := NewAssociationAdapter(nil)

	// Test nil association
	_, err := adapter.Create(context.TODO(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "association creation data is required")

	// Test missing required fields
	_, err = adapter.Create(context.TODO(), &types.AssociationCreate{
		AccountName: "",
		UserName:    "user1",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "account name is required")

	_, err = adapter.Create(context.TODO(), &types.AssociationCreate{
		AccountName: "account1",
		UserName:    "",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user name is required")

	// Nil context validation is covered in TestAssociationAdapter_ValidateContext
}

func TestAssociationAdapter_Update(t *testing.T) {
	adapter := NewAssociationAdapter(nil)

	// Test nil update
	err := adapter.Update(context.TODO(), "test-id", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "association update data is required")

	// Test empty association ID (nil context validation is covered in TestAssociationAdapter_ValidateContext)
	err = adapter.Update(context.TODO(), "", &types.AssociationUpdate{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "associationID is required")
}

func TestAssociationAdapter_Delete(t *testing.T) {
	adapter := NewAssociationAdapter(nil)

	// Test empty association ID
	err := adapter.Delete(context.TODO(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "associationID is required")

	// Test client initialization check (nil context validation is covered in TestAssociationAdapter_ValidateContext)
	err = adapter.Delete(context.TODO(), "test-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
}

func TestAssociationAdapter_ValidateAssociationCreate(t *testing.T) {
	adapter := NewAssociationAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name          string
		association   *types.AssociationCreate
		expectedError bool
		errorContains string
	}{
		{
			name: "valid association",
			association: &types.AssociationCreate{
				AccountName: "account1",
				UserName:    "user1",
			},
			expectedError: false,
		},
		{
			name:          "nil association",
			association:   nil,
			expectedError: true,
			errorContains: "association creation data is required",
		},
		{
			name: "missing account name",
			association: &types.AssociationCreate{
				AccountName: "",
				UserName:    "user1",
			},
			expectedError: true,
			errorContains: "account name is required",
		},
		{
			name: "missing user name",
			association: &types.AssociationCreate{
				AccountName: "account1",
				UserName:    "",
			},
			expectedError: true,
			errorContains: "user name is required",
		},
		{
			name: "minimal valid",
			association: &types.AssociationCreate{
				AccountName: "account1",
				UserName:    "user1",
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.validateAssociationCreate(tt.association)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAssociationAdapter_ValidateAssociationUpdate(t *testing.T) {
	adapter := NewAssociationAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name          string
		update        *types.AssociationUpdate
		expectedError bool
		errorContains string
	}{
		{
			name: "valid update",
			update: &types.AssociationUpdate{
				DefaultQoS: ptrString("normal"),
			},
			expectedError: false,
		},
		{
			name:          "nil update",
			update:        nil,
			expectedError: true,
			errorContains: "association update data is required",
		},
		{
			name:          "empty update",
			update:        &types.AssociationUpdate{},
			expectedError: false, // Empty updates are allowed
		},
		{
			name: "update with QoS list",
			update: &types.AssociationUpdate{
				QoSList: []string{"normal", "high", "critical"},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.validateAssociationUpdate(tt.update)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAssociationAdapter_GetDefaultClusterName(t *testing.T) {
	adapter := NewAssociationAdapter(&api.ClientWithResponses{})

	// Test that it returns a non-empty default cluster name
	clusterName := adapter.getDefaultClusterName()
	assert.NotEmpty(t, clusterName)
	// Common default is "cluster"
	assert.Contains(t, []string{"cluster", "default", "main"}, clusterName)
}
