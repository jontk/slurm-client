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

func TestNewStandaloneAdapter(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewStandaloneAdapter(client)

	assert.NotNil(t, adapter)
	assert.Equal(t, client, adapter.client)
	assert.NotNil(t, adapter.errorAdapter)
}

func TestStandaloneAdapter_ClientValidation(t *testing.T) {
	// Test nil client validation
	adapter := NewStandaloneAdapter(nil)
	ctx := context.Background()

	_, err := adapter.GetLicenses(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	_, err = adapter.GetDiagnostics(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	_, err = adapter.GetShares(ctx, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	// Test that non-nil client passes initial validation
	validAdapter := NewStandaloneAdapter(&api.ClientWithResponses{})
	assert.NotNil(t, validAdapter.client)
}

func TestStandaloneAdapter_GetLicenses(t *testing.T) {
	adapter := NewStandaloneAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	_, err := adapter.GetLicenses(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")
}

func TestStandaloneAdapter_GetDiagnostics(t *testing.T) {
	adapter := NewStandaloneAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	_, err := adapter.GetDiagnostics(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")
}

func TestStandaloneAdapter_GetShares(t *testing.T) {
	adapter := NewStandaloneAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	_, err := adapter.GetShares(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")
}

func TestStandaloneAdapter_ErrorHandling(t *testing.T) {
	adapter := NewStandaloneAdapter(nil)
	ctx := context.Background()

	// Test various error conditions with nil client
	tests := []struct {
		name     string
		testFunc func() error
	}{
		{
			name: "GetLicenses with nil client",
			testFunc: func() error {
				_, err := adapter.GetLicenses(ctx)
				return err
			},
		},
		{
			name: "GetDiagnostics with nil client",
			testFunc: func() error {
				_, err := adapter.GetDiagnostics(ctx)
				return err
			},
		},
		{
			name: "GetShares with nil client",
			testFunc: func() error {
				_, err := adapter.GetShares(ctx, nil)
				return err
			},
		},
		{
			name: "GetConfig with nil client",
			testFunc: func() error {
				_, err := adapter.GetConfig(ctx)
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
