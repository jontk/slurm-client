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

func TestNewQoSAdapter(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewQoSAdapter(client)

	assert.NotNil(t, adapter)
	assert.Equal(t, client, adapter.client)
	assert.NotNil(t, adapter.BaseManager)
	assert.Equal(t, "v0.0.42", adapter.GetVersion())
}

func TestQoSAdapter_ValidateContext(t *testing.T) {
	adapter := &QoSAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "QoS"),
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

func TestQoSAdapter_ClientValidation(t *testing.T) {
	// Test nil client validation
	adapter := NewQoSAdapter(nil)
	ctx := context.Background()

	_, err := adapter.List(ctx, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	// Test that non-nil client passes initial validation
	// (actual API calls would fail due to no server, but that's different from client validation)
	validAdapter := NewQoSAdapter(&api.ClientWithResponses{})
	assert.NotNil(t, validAdapter.client) // Ensure client is set
}

func TestQoSAdapter_ConvertAPIQoSToCommon(t *testing.T) {
	adapter := NewQoSAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name     string
		apiQoS   api.V0042Qos
		expected types.QoS
	}{
		{
			name: "basic conversion",
			apiQoS: api.V0042Qos{
				Name: ptrString("test-qos"),
				Id:   ptrInt32(123),
				Priority: &api.V0042Uint32NoValStruct{
					Set:    ptrBool(true),
					Number: ptrInt32(100),
				},
				Description: ptrString("Test description"),
			},
			expected: types.QoS{
				Name:        "test-qos",
				ID:          123,
				Priority:    100,
				Description: "Test description",
			},
		},
		{
			name: "minimal conversion",
			apiQoS: api.V0042Qos{
				Name: ptrString("minimal-qos"),
			},
			expected: types.QoS{
				Name: "minimal-qos",
			},
		},
		{
			name: "priority not set",
			apiQoS: api.V0042Qos{
				Name: ptrString("no-priority-qos"),
				Priority: &api.V0042Uint32NoValStruct{
					Set:    ptrBool(false),
					Number: ptrInt32(100),
				},
			},
			expected: types.QoS{
				Name:     "no-priority-qos",
				Priority: 0, // Should be 0 when not set
			},
		},
		{
			name: "nil priority",
			apiQoS: api.V0042Qos{
				Name:     ptrString("nil-priority-qos"),
				Priority: nil,
			},
			expected: types.QoS{
				Name:     "nil-priority-qos",
				Priority: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.convertAPIQoSToCommon(tt.apiQoS)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.Priority, result.Priority)
			assert.Equal(t, tt.expected.Description, result.Description)
		})
	}
}

func TestQoSAdapter_ConvertCommonQoSCreateToAPI(t *testing.T) {
	adapter := NewQoSAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name          string
		qosCreate     *QoSCreateRequest
		expectedError bool
		expectedName  string
	}{
		{
			name: "basic conversion",
			qosCreate: &QoSCreateRequest{
				Name:        "test-qos",
				Description: "Test QoS",
				Priority:    100,
			},
			expectedError: false,
			expectedName:  "test-qos",
		},
		{
			name: "minimal conversion",
			qosCreate: &QoSCreateRequest{
				Name: "minimal-qos",
			},
			expectedError: false,
			expectedName:  "minimal-qos",
		},
		{
			name:          "nil input",
			qosCreate:     nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.convertCommonQoSCreateToAPI(tt.qosCreate)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotNil(t, result.Qos)
				assert.Len(t, result.Qos, 1)
				assert.Equal(t, tt.expectedName, *result.Qos[0].Name)
			}
		})
	}
}

func TestQoSAdapter_UpdateNotSupported(t *testing.T) {
	adapter := NewQoSAdapter(&api.ClientWithResponses{})
	ctx := context.Background()

	err := adapter.Update(ctx, "test-qos", &types.QoSUpdateRequest{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not directly supported")
}

func TestQoSAdapter_DeleteNotImplemented(t *testing.T) {
	adapter := NewQoSAdapter(&api.ClientWithResponses{})
	ctx := context.Background()

	err := adapter.Delete(ctx, "test-qos")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

func TestQoSAdapter_ListOptionsHandling(t *testing.T) {
	adapter := NewQoSAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	tests := []struct {
		name string
		opts *types.QoSListOptions
	}{
		{
			name: "nil options",
			opts: nil,
		},
		{
			name: "empty options",
			opts: &types.QoSListOptions{},
		},
		{
			name: "options with names",
			opts: &types.QoSListOptions{
				Names: []string{"qos1", "qos2"},
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

func TestQoSAdapter_GetByName(t *testing.T) {
	adapter := NewQoSAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	tests := []struct {
		name    string
		qosName string
	}{
		{
			name:    "valid name",
			qosName: "test-qos",
		},
		{
			name:    "empty name",
			qosName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := adapter.Get(ctx, tt.qosName)
			// Should get client validation error
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "client")
		})
	}
}

func TestQoSAdapter_CreateValidation(t *testing.T) {
	adapter := NewQoSAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	tests := []struct {
		name      string
		qosCreate *types.QoSCreate
	}{
		{
			name: "valid create",
			qosCreate: &types.QoSCreate{
				Name:        "test-qos",
				Description: "Test QoS",
				Priority:    100,
			},
		},
		{
			name: "minimal create",
			qosCreate: &types.QoSCreate{
				Name: "minimal-qos",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := adapter.Create(ctx, tt.qosCreate)
			// Should get client validation error
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "client")
		})
	}
}
