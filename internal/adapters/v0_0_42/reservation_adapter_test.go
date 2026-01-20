// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
)

func TestNewReservationAdapter(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewReservationAdapter(client)

	assert.NotNil(t, adapter)
	assert.Equal(t, client, adapter.client)
	assert.NotNil(t, adapter.BaseManager)
	assert.Equal(t, "v0.0.42", adapter.GetVersion())
}

func TestReservationAdapter_ValidateContext(t *testing.T) {
	adapter := &ReservationAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Reservation"),
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

func TestReservationAdapter_ClientValidation(t *testing.T) {
	// Test nil client validation
	adapter := NewReservationAdapter(nil)
	ctx := context.Background()

	_, err := adapter.List(ctx, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	_, err = adapter.Get(ctx, "test-reservation")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	err = adapter.Delete(ctx, "test-reservation")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	// Test that non-nil client passes initial validation
	validAdapter := NewReservationAdapter(&api.ClientWithResponses{})
	assert.NotNil(t, validAdapter.client)
}

func TestReservationAdapter_ListOptionsHandling(t *testing.T) {
	adapter := NewReservationAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	tests := []struct {
		name string
		opts *types.ReservationListOptions
	}{
		{
			name: "nil options",
			opts: nil,
		},
		{
			name: "empty options",
			opts: &types.ReservationListOptions{},
		},
		{
			name: "options with names",
			opts: &types.ReservationListOptions{
				Names: []string{"maint-window", "gpu-reservation"},
			},
		},
		{
			name: "options with update time",
			opts: &types.ReservationListOptions{
				UpdateTime: func() *time.Time { t := time.Now(); return &t }(),
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

func TestReservationAdapter_GetByName(t *testing.T) {
	adapter := NewReservationAdapter(nil) // Use nil client to test validation path
	ctx := context.Background()

	tests := []struct {
		name            string
		reservationName string
	}{
		{
			name:            "valid name",
			reservationName: "maint-window",
		},
		{
			name:            "empty name",
			reservationName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := adapter.Get(ctx, tt.reservationName)
			// Should get client validation error
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "client")
		})
	}
}

func TestReservationAdapter_ConvertAPIReservationToCommon(t *testing.T) {
	adapter := NewReservationAdapter(&api.ClientWithResponses{})

	now := time.Now()
	later := now.Add(2 * time.Hour)

	tests := []struct {
		name            string
		apiReservation  api.V0042ReservationInfo
		expected        types.Reservation
		expectStartTime bool
		expectEndTime   bool
	}{
		{
			name: "basic reservation",
			apiReservation: api.V0042ReservationInfo{
				Name: ptrString("maint-window"),
				StartTime: &api.V0042Uint64NoValStruct{
					Set:    ptrBool(true),
					Number: ptrInt64(now.Unix()),
				},
				EndTime: &api.V0042Uint64NoValStruct{
					Set:    ptrBool(true),
					Number: ptrInt64(later.Unix()),
				},
				NodeList:  ptrString("node[1-10]"),
				NodeCount: ptrInt32(10),
			},
			expected: types.Reservation{
				Name:      "maint-window",
				NodeList:  "node[1-10]",
				NodeCount: 10,
			},
			expectStartTime: true,
			expectEndTime:   true,
		},
		{
			name: "reservation with users",
			apiReservation: api.V0042ReservationInfo{
				Name:     ptrString("user-reservation"),
				NodeList: ptrString("compute[1-5]"),
				Users:    ptrString("user1,user2,user3"),
			},
			expected: types.Reservation{
				Name:     "user-reservation",
				NodeList: "compute[1-5]",
				Users:    []string{"user1", "user2", "user3"},
			},
		},
		{
			name: "reservation with accounts",
			apiReservation: api.V0042ReservationInfo{
				Name:     ptrString("account-reservation"),
				NodeList: ptrString("gpu[1-4]"),
				Accounts: ptrString("research,development"),
			},
			expected: types.Reservation{
				Name:     "account-reservation",
				NodeList: "gpu[1-4]",
				Accounts: []string{"research", "development"},
			},
		},
		{
			name: "minimal reservation",
			apiReservation: api.V0042ReservationInfo{
				Name: ptrString("minimal"),
			},
			expected: types.Reservation{
				Name: "minimal",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.convertAPIReservationToCommon(tt.apiReservation)

			require.NoError(t, err)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.NodeList, result.NodeList)
			assert.Equal(t, tt.expected.NodeCount, result.NodeCount)

			if tt.expectStartTime {
				assert.NotZero(t, result.StartTime)
			}
			if tt.expectEndTime {
				assert.NotZero(t, result.EndTime)
			}
		})
	}
}

func TestReservationAdapter_ErrorHandling(t *testing.T) {
	adapter := NewReservationAdapter(nil)
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
				_, err := adapter.Get(ctx, "maint-window")
				return err
			},
		},
		{
			name: "Create with nil client",
			testFunc: func() error {
				now := time.Now()
				later := now.Add(2 * time.Hour)
				_, err := adapter.Create(ctx, &types.ReservationCreate{
					Name:      "test",
					StartTime: now,
					EndTime:   &later,
					NodeList:  "node[1-5]",
				})
				return err
			},
		},
		{
			name: "Update with nil client",
			testFunc: func() error {
				return adapter.Update(ctx, "test", &types.ReservationUpdate{})
			},
		},
		{
			name: "Delete with nil client",
			testFunc: func() error {
				return adapter.Delete(ctx, "test")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testFunc()
			assert.Error(t, err)
			// Should get either client validation error or "not supported" error
			errorMsg := err.Error()
			assert.True(t,
				strings.Contains(errorMsg, "client") ||
					strings.Contains(errorMsg, "not supported"),
				"Expected client validation or not supported error, got: %v", err)
		})
	}
}

func TestReservationAdapter_CreateNotSupported(t *testing.T) {
	adapter := NewReservationAdapter(&api.ClientWithResponses{})
	ctx := context.Background()

	now := time.Now()
	later := now.Add(2 * time.Hour)
	reservation := &types.ReservationCreate{
		Name:      "new-reservation",
		StartTime: now,
		EndTime:   &later,
		NodeList:  "node[1-5]",
		Users:     []string{"user1"},
	}

	_, err := adapter.Create(ctx, reservation)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not supported")
}

func TestReservationAdapter_UpdateNotSupported(t *testing.T) {
	adapter := NewReservationAdapter(&api.ClientWithResponses{})
	ctx := context.Background()

	later := time.Now().Add(4 * time.Hour)
	update := &types.ReservationUpdate{
		EndTime:  &later,
		NodeList: ptrString("node[1-20]"),
	}

	err := adapter.Update(ctx, "maint-window", update)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not supported")
}
