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

func TestNewQoSAdapter(t *testing.T) {
	adapter := NewQoSAdapter(&api.ClientWithResponses{})
	assert.NotNil(t, adapter)
	assert.NotNil(t, adapter.BaseManager)
}

func TestQoSAdapter_ValidateContext(t *testing.T) {
	adapter := NewQoSAdapter(&api.ClientWithResponses{})

	// Test nil context
	//lint:ignore SA1012 intentionally testing nil context validation
	err := adapter.ValidateContext(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context is required")

	// Test valid context
	err = adapter.ValidateContext(context.Background())
	assert.NoError(t, err)
}

func TestQoSAdapter_List(t *testing.T) {
	adapter := NewQoSAdapter(nil) // Use nil client for testing validation logic

	// Test client initialization check (nil context validation is covered in TestQoSAdapter_ValidateContext)
	_, err := adapter.List(context.TODO(), &types.QoSListOptions{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
}

func TestQoSAdapter_Get(t *testing.T) {
	adapter := NewQoSAdapter(nil)

	// Test empty QoS name
	_, err := adapter.Get(context.TODO(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "QoS name is required")

	// Test client initialization check (nil context validation is covered in TestQoSAdapter_ValidateContext)
	_, err = adapter.Get(context.TODO(), "test-qos")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
}

func TestQoSAdapter_ConvertAPIQoSToCommon(t *testing.T) {
	adapter := NewQoSAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name         string
		apiQoS       api.V0043Qos
		expectedName string
	}{
		{
			name: "full qos",
			apiQoS: api.V0043Qos{
				Name:        ptrString("normal"),
				Description: ptrString("Normal priority queue"),
			},
			expectedName: "normal",
		},
		{
			name: "minimal qos",
			apiQoS: api.V0043Qos{
				Name: ptrString("debug"),
			},
			expectedName: "debug",
		},
		{
			name:         "empty qos",
			apiQoS:       api.V0043Qos{},
			expectedName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.convertAPIQoSToCommon(tt.apiQoS)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedName, result.Name)
		})
	}
}

func TestQoSAdapter_Create(t *testing.T) {
	adapter := NewQoSAdapter(nil)

	// Test nil QoS
	_, err := adapter.Create(context.TODO(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "QoS creation data is required")

	// Test missing required fields
	_, err = adapter.Create(context.TODO(), &types.QoSCreate{
		Name: "",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "QoS name is required")

	// Test client initialization check (nil context validation is covered in TestQoSAdapter_ValidateContext)
	_, err = adapter.Create(context.TODO(), &types.QoSCreate{
		Name: "test-qos",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
}

func TestQoSAdapter_Update(t *testing.T) {
	adapter := NewQoSAdapter(nil)

	// Test nil update
	err := adapter.Update(context.TODO(), "test-qos", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "QoS update data is required")

	// Test empty QoS name
	err = adapter.Update(context.TODO(), "", &types.QoSUpdate{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "QoS name is required")

	// Test client initialization check (nil context validation is covered in TestQoSAdapter_ValidateContext)
	err = adapter.Update(context.TODO(), "test-qos", &types.QoSUpdate{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
}

func TestQoSAdapter_Delete(t *testing.T) {
	adapter := NewQoSAdapter(nil)

	// Test empty QoS name
	err := adapter.Delete(context.TODO(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "QoS name is required")

	// Test client initialization check (nil context validation is covered in TestQoSAdapter_ValidateContext)
	err = adapter.Delete(context.TODO(), "test-qos")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
}

func TestQoSAdapter_ValidateQoSCreate(t *testing.T) {
	adapter := NewQoSAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name          string
		qos           *types.QoSCreate
		expectedError bool
		errorContains string
	}{
		{
			name: "valid qos",
			qos: &types.QoSCreate{
				Name: "test-qos",
			},
			expectedError: false,
		},
		{
			name:          "nil qos",
			qos:           nil,
			expectedError: true,
			errorContains: "QoS creation data is required",
		},
		{
			name: "missing name",
			qos: &types.QoSCreate{
				Name: "",
			},
			expectedError: true,
			errorContains: "QoS name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateQoSCreate(tt.qos)

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

func TestQoSAdapter_ValidateQoSUpdate(t *testing.T) {
	adapter := NewQoSAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name          string
		update        *types.QoSUpdate
		expectedError bool
		errorContains string
	}{
		{
			name: "valid update",
			update: &types.QoSUpdate{
				Description: ptrString("Updated description"),
			},
			expectedError: false,
		},
		{
			name:          "nil update",
			update:        nil,
			expectedError: true,
			errorContains: "QoS update data is required",
		},
		{
			name:          "empty update",
			update:        &types.QoSUpdate{},
			expectedError: false, // Empty updates are allowed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateQoSUpdate(tt.update)

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
