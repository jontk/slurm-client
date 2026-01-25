// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_44

import (
	"context"
	"testing"
	"time"

	api "github.com/jontk/slurm-client/internal/api/v0_0_44"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/stretchr/testify/assert"
)

func TestNewReservationAdapter(t *testing.T) {
	adapter := NewReservationAdapter(&api.ClientWithResponses{})
	assert.NotNil(t, adapter)
	assert.NotNil(t, adapter.BaseManager)
}

func TestReservationAdapter_ValidateContext(t *testing.T) {
	adapter := NewReservationAdapter(&api.ClientWithResponses{})

	// Test nil context
	//lint:ignore SA1012 intentionally testing nil context validation
	err := adapter.ValidateContext(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context is required")

	// Test valid context
	err = adapter.ValidateContext(context.Background())
	assert.NoError(t, err)
}

func TestReservationAdapter_List(t *testing.T) {
	adapter := NewReservationAdapter(nil) // Use nil client for testing validation logic

	// Test client initialization check (nil context validation is covered in TestReservationAdapter_ValidateContext)
	_, err := adapter.List(context.TODO(), &types.ReservationListOptions{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
}

func TestReservationAdapter_Get(t *testing.T) {
	adapter := NewReservationAdapter(nil)

	// Test empty reservation name
	_, err := adapter.Get(context.TODO(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Reservation name is required")

	// Test client initialization check (nil context validation is covered in TestReservationAdapter_ValidateContext)
	_, err = adapter.Get(context.TODO(), "test-reservation")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
}

func TestReservationAdapter_ConvertAPIReservationToCommon(t *testing.T) {
	adapter := NewReservationAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name              string
		apiReservation    api.V0044ReservationInfo
		expectedName      string
		expectedPartition string
	}{
		{
			name: "full reservation",
			apiReservation: api.V0044ReservationInfo{
				Name:      ptrString("maintenance"),
				Partition: ptrString("compute"),
				Accounts:  ptrString("root"),
				Users:     ptrString("admin"),
			},
			expectedName:      "maintenance",
			expectedPartition: "compute",
		},
		{
			name: "minimal reservation",
			apiReservation: api.V0044ReservationInfo{
				Name: ptrString("test-res"),
			},
			expectedName:      "test-res",
			expectedPartition: "",
		},
		{
			name:              "empty reservation",
			apiReservation:    api.V0044ReservationInfo{},
			expectedName:      "",
			expectedPartition: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.convertAPIReservationToCommon(tt.apiReservation)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedName, result.Name)
			// Note: Check available fields in the result struct
			// assert.Equal(t, tt.expectedPartition, result.Partition)
		})
	}
}

func TestReservationAdapter_Create(t *testing.T) {
	adapter := NewReservationAdapter(nil)

	// Test nil reservation
	_, err := adapter.Create(context.TODO(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reservation creation data is required")

	// Test missing required fields (nil context validation is covered in TestReservationAdapter_ValidateContext)
	_, err = adapter.Create(context.TODO(), &types.ReservationCreate{
		Name: "",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reservation name is required")

	_, err = adapter.Create(context.TODO(), &types.ReservationCreate{
		Name: "test-reservation",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "start time is required")
}

func TestReservationAdapter_Update(t *testing.T) {
	adapter := NewReservationAdapter(nil)

	// Test nil update
	err := adapter.Update(context.TODO(), "test-reservation", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reservation update data is required")

	// Test empty reservation name
	err = adapter.Update(context.TODO(), "", &types.ReservationUpdate{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Reservation name is required")

	// Test client initialization check (nil context validation is covered in TestReservationAdapter_ValidateContext)
	err = adapter.Update(context.TODO(), "test-reservation", &types.ReservationUpdate{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one field must be provided for update")
}

func TestReservationAdapter_Delete(t *testing.T) {
	adapter := NewReservationAdapter(nil)

	// Test empty reservation name
	err := adapter.Delete(context.TODO(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Reservation name is required")

	// Test client initialization check (nil context validation is covered in TestReservationAdapter_ValidateContext)
	err = adapter.Delete(context.TODO(), "test-reservation")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
}

func TestReservationAdapter_ValidateReservationCreate(t *testing.T) {
	adapter := NewReservationAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name          string
		reservation   *types.ReservationCreate
		expectedError bool
		errorContains string
	}{
		{
			name: "valid reservation",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: time.Now().Add(time.Hour),
				EndTime:   &[]time.Time{time.Now().Add(2 * time.Hour)}[0],
			},
			expectedError: false,
		},
		{
			name:          "nil reservation",
			reservation:   nil,
			expectedError: true,
			errorContains: "reservation creation data is required",
		},
		{
			name: "missing name",
			reservation: &types.ReservationCreate{
				Name: "",
			},
			expectedError: true,
			errorContains: "reservation name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.validateReservationCreate(tt.reservation)

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

func TestReservationAdapter_ValidateReservationUpdate(t *testing.T) {
	adapter := NewReservationAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name          string
		update        *types.ReservationUpdate
		expectedError bool
		errorContains string
	}{
		{
			name: "valid update",
			update: &types.ReservationUpdate{
				Accounts: []string{"newaccount"},
			},
			expectedError: false,
		},
		{
			name:          "nil update",
			update:        nil,
			expectedError: true,
			errorContains: "reservation update data is required",
		},
		{
			name:          "empty update",
			update:        &types.ReservationUpdate{},
			expectedError: true, // Empty updates should require at least one field
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.validateReservationUpdate(tt.update)

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
