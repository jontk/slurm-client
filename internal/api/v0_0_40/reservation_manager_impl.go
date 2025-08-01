// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"context"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// ReservationManagerImpl implements the ReservationManager interface for v0.0.40
type ReservationManagerImpl struct {
	client *WrapperClient
}

// NewReservationManagerImpl creates a new ReservationManagerImpl
func NewReservationManagerImpl(client *WrapperClient) *ReservationManagerImpl {
	return &ReservationManagerImpl{
		client: client,
	}
}

// List retrieves a list of reservations with optional filtering
func (r *ReservationManagerImpl) List(ctx context.Context, opts *interfaces.ListReservationsOptions) (*interfaces.ReservationList, error) {
	// v0.0.40 is an older API version with limited support
	// Reservation management was added in later versions
	return nil, errors.NewNotImplementedError("reservation listing", "v0.0.40")
}

// Get retrieves a specific reservation by name
func (r *ReservationManagerImpl) Get(ctx context.Context, reservationName string) (*interfaces.Reservation, error) {
	if reservationName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "reservation name is required", "reservationName", reservationName, nil)
	}
	return nil, errors.NewNotImplementedError("reservation retrieval", "v0.0.40")
}

// Create creates a new reservation
func (r *ReservationManagerImpl) Create(ctx context.Context, reservation *interfaces.ReservationCreate) (*interfaces.ReservationCreateResponse, error) {
	if reservation == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "reservation data is required", "reservation", reservation, nil)
	}
	return nil, errors.NewNotImplementedError("reservation creation", "v0.0.40")
}

// Update updates an existing reservation
func (r *ReservationManagerImpl) Update(ctx context.Context, reservationName string, update *interfaces.ReservationUpdate) error {
	if reservationName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "reservation name is required", "reservationName", reservationName, nil)
	}
	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "update data is required", "update", update, nil)
	}
	return errors.NewNotImplementedError("reservation update", "v0.0.40")
}

// Delete deletes a reservation
func (r *ReservationManagerImpl) Delete(ctx context.Context, reservationName string) error {
	if reservationName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "reservation name is required", "reservationName", reservationName, nil)
	}
	return errors.NewNotImplementedError("reservation deletion", "v0.0.40")
}
