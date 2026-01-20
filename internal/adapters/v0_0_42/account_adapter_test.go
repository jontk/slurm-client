// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

func TestNewAccountAdapter(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewAccountAdapter(client)

	assert.NotNil(t, adapter)
	assert.Equal(t, client, adapter.client)
	assert.NotNil(t, adapter.BaseManager)
}

func TestAccountAdapter_ValidateContext(t *testing.T) {
	adapter := NewAccountAdapter(&api.ClientWithResponses{})

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

func TestAccountAdapter_ClientValidation(t *testing.T) {
	// Test nil client validation
	adapter := NewAccountAdapter(nil)
	ctx := context.Background()

	_, err := adapter.List(ctx, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	_, err = adapter.Get(ctx, "test-account")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	// Test that non-nil client passes initial validation
	validAdapter := NewAccountAdapter(&api.ClientWithResponses{})
	assert.NotNil(t, validAdapter.client)
}

func TestAccountAdapter_GetByName(t *testing.T) {
	adapter := NewAccountAdapter(nil) // Use nil client to test validation path

	tests := []struct {
		name          string
		accountName   string
		expectedError bool
		expectedMsg   string
	}{
		{
			name:          "valid name",
			accountName:   "test-account",
			expectedError: true,
			expectedMsg:   "client",
		},
		{
			name:          "empty name",
			accountName:   "",
			expectedError: true,
			expectedMsg:   "client", // Client validation happens first before name validation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := adapter.Get(context.Background(), tt.accountName)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAccountAdapter_ErrorHandling(t *testing.T) {
	adapter := NewAccountAdapter(nil)
	ctx := context.Background()

	// Test various error conditions with nil client
	tests := []struct {
		name     string
		testFunc func() error
	}{
		{
			name: "List with nil client",
			testFunc: func() error {
				_, err := adapter.List(ctx, nil)
				return err
			},
		},
		{
			name: "Get with nil client",
			testFunc: func() error {
				_, err := adapter.Get(ctx, "test")
				return err
			},
		},
		{
			name: "Create with nil client",
			testFunc: func() error {
				_, err := adapter.Create(ctx, nil)
				return err
			},
		},
		{
			name: "Update with nil client",
			testFunc: func() error {
				err := adapter.Update(ctx, "test", nil)
				return err
			},
		},
		{
			name: "Delete with nil client",
			testFunc: func() error {
				err := adapter.Delete(ctx, "test")
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testFunc()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "client")
		})
	}
}
