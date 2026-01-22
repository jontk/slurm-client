// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_44

import (
	"net/http"
	"context"
	"fmt"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// ReservationManagerImpl provides the actual implementation for ReservationManager methods
type ReservationManagerImpl struct {
	client *WrapperClient
}

// NewReservationManagerImpl creates a new ReservationManager implementation
func NewReservationManagerImpl(client *WrapperClient) *ReservationManagerImpl {
	return &ReservationManagerImpl{client: client}
}

// List retrieves a list of reservations with optional filtering
func (m *ReservationManagerImpl) List(ctx context.Context, opts *interfaces.ListReservationsOptions) (*interfaces.ReservationList, error) {
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0044GetReservationsWithResponse(ctx, &SlurmV0044GetReservationsParams{})
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	if resp.HTTPResponse.StatusCode != http.StatusOK {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal,
			fmt.Sprintf("HTTP %d", resp.HTTPResponse.StatusCode))
	}

	// Convert to interface type
	reservationList := &interfaces.ReservationList{
		Reservations: make([]interfaces.Reservation, 0),
	}

	// TODO: Convert actual reservations when response structure is known
	return reservationList, nil
}

// Get retrieves a specific reservation by name
func (m *ReservationManagerImpl) Get(ctx context.Context, reservationName string) (*interfaces.Reservation, error) {
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0044GetReservationWithResponse(ctx, reservationName, &SlurmV0044GetReservationParams{})
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	if resp.HTTPResponse.StatusCode == http.StatusNotFound {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound,
			fmt.Sprintf("Reservation %s not found", reservationName))
	}

	if resp.HTTPResponse.StatusCode != http.StatusOK {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal,
			fmt.Sprintf("HTTP %d", resp.HTTPResponse.StatusCode))
	}

	// TODO: Convert to interface type when response structure is known
	return &interfaces.Reservation{Name: reservationName}, nil
}

// Create creates a new reservation
func (m *ReservationManagerImpl) Create(ctx context.Context, reservation *interfaces.ReservationCreate) (*interfaces.ReservationCreateResponse, error) {
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// TODO: Implement reservation creation
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Create reservation not yet implemented for v0.0.44")
}

// Update updates an existing reservation
func (m *ReservationManagerImpl) Update(ctx context.Context, reservationName string, update *interfaces.ReservationUpdate) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// TODO: Implement reservation update
	return errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Update reservation not yet implemented for v0.0.44")
}

// Delete deletes a reservation
func (m *ReservationManagerImpl) Delete(ctx context.Context, reservationName string) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// TODO: Implement reservation deletion
	return errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Delete reservation not yet implemented for v0.0.44")
}
