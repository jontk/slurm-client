// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package base

import (
	"testing"
	"time"

	types "github.com/jontk/slurm-client/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReservationBaseManager_New(t *testing.T) {
	manager := NewReservationBaseManager("v0.0.43")
	assert.NotNil(t, manager)
	assert.Equal(t, "v0.0.43", manager.GetVersion())
	assert.Equal(t, "Reservation", manager.GetResourceType())
}
func TestReservationBaseManager_ValidateReservationCreate(t *testing.T) {
	manager := NewReservationBaseManager("v0.0.43")
	tests := []struct {
		name        string
		reservation *types.ReservationCreate
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "nil reservation",
			reservation: nil,
			wantErr:     true,
			errMsg:      "data is required",
		},
		{
			name: "empty name",
			reservation: &types.ReservationCreate{
				Name: stringPtr(""),
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "valid reservation",
			reservation: &types.ReservationCreate{
				Name:      stringPtr("test-reservation"),
				StartTime: time.Now(),
				Duration:  uint32Ptr(60),            // 60 minutes - only duration, no end time
				Accounts:  []string{"test-account"}, // Add required account
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateReservationCreate(tt.reservation)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
func TestReservationBaseManager_ValidateReservationUpdate(t *testing.T) {
	manager := NewReservationBaseManager("v0.0.43")
	tests := []struct {
		name    string
		update  *types.ReservationUpdate
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
			update: &types.ReservationUpdate{
				NodeCount: int32Ptr(15),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateReservationUpdate(tt.update)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
