// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
)

func TestNewNodeAdapter(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewNodeAdapter(client)

	assert.NotNil(t, adapter)
	assert.Equal(t, client, adapter.client)
	assert.NotNil(t, adapter.BaseManager)
	assert.Equal(t, "v0.0.42", adapter.GetVersion())
}

func TestNodeAdapter_ValidateContext(t *testing.T) {
	adapter := &NodeAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Node"),
		client:      &api.ClientWithResponses{},
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
			errMsg:  "context",
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

func TestNodeAdapter_ClientValidation(t *testing.T) {
	// Test nil client validation
	adapter := NewNodeAdapter(nil)
	ctx := context.Background()

	_, err := adapter.List(ctx, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	// Test that non-nil client passes initial validation
	validAdapter := NewNodeAdapter(&api.ClientWithResponses{})
	assert.NotNil(t, validAdapter.client)
}

func TestNodeAdapter_ListOptionsHandling(t *testing.T) {
	adapter := NewNodeAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	tests := []struct {
		name string
		opts *types.NodeListOptions
	}{
		{
			name: "nil options",
			opts: nil,
		},
		{
			name: "empty options",
			opts: &types.NodeListOptions{},
		},
		{
			name: "options with names",
			opts: &types.NodeListOptions{
				Names: []string{"node1", "node2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := adapter.List(ctx, tt.opts)
			// Should get client validation error before any option processing
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "client")
		})
	}
}

func TestNodeAdapter_GetByName(t *testing.T) {
	adapter := NewNodeAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	tests := []struct {
		name     string
		nodeName string
	}{
		{
			name:     "valid name",
			nodeName: "node1",
		},
		{
			name:     "empty name",
			nodeName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := adapter.Get(ctx, tt.nodeName)
			// Should get client validation error
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "client")
		})
	}
}

func TestNodeAdapter_ErrorHandling(t *testing.T) {
	adapter := NewNodeAdapter(nil)
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
				_, err := adapter.Get(ctx, "node1")
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
