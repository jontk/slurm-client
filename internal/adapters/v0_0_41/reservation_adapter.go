// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

import (
	"context"
	"fmt"
	"strconv"

	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
)

// ReservationAdapter implements the ReservationAdapter interface for v0.0.41
type ReservationAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewReservationAdapter creates a new Reservation adapter for v0.0.41
func NewReservationAdapter(client *api.ClientWithResponses) *ReservationAdapter {
	return &ReservationAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Reservation"),
		client:      client,
		wrapper:     nil, // We'll implement this later
	}
}

// List retrieves a list of reservations with optional filtering
func (a *ReservationAdapter) List(ctx context.Context, opts *types.ReservationListOptions) (*types.ReservationList, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmV0041GetReservationsParams{}

	// Apply filters from options
	if opts != nil {
		if opts.UpdateTime != nil {
			updateTimeStr := strconv.FormatInt(opts.UpdateTime.Unix(), 10)
			params.UpdateTime = &updateTimeStr
		}
	}

	// Make the API call
	resp, err := a.client.SlurmV0041GetReservationsWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list reservations")
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected nil response")
	}

	// Convert response to common types
	resList := &types.ReservationList{
		Reservations: make([]types.Reservation, 0, len(resp.JSON200.Reservations)),
		Total:        0,
	}

	for _, apiRes := range resp.JSON200.Reservations {
		res, err := a.convertAPIReservationToCommon(apiRes)
		if err != nil {
			// Log the error but continue processing other reservations
			continue
		}
		resList.Reservations = append(resList.Reservations, *res)
	}

	// Extract warning and error messages if any (but ReservationList doesn't have Meta)
	// Warnings and errors are ignored for now as ReservationList structure doesn't support them
	if resp.JSON200.Warnings != nil {
		// Log warnings if needed
		_ = resp.JSON200.Warnings
	}
	if resp.JSON200.Errors != nil {
		// Log errors if needed
		_ = resp.JSON200.Errors
	}

	// Update total count
	resList.Total = len(resList.Reservations)

	return resList, nil
}

// Get retrieves a specific reservation by name
func (a *ReservationAdapter) Get(ctx context.Context, name string) (*types.Reservation, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Validate name
	if err := a.ValidateResourceName("reservation name", name); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Make the API call
	params := &api.SlurmV0041GetReservationParams{}
	resp, err := a.client.SlurmV0041GetReservationWithResponse(ctx, name, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to get reservation "+name)
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}

	if resp.JSON200 == nil || len(resp.JSON200.Reservations) == 0 {
		return nil, a.HandleNotFound("reservation " + name)
	}

	// Convert the first reservation in the response
	res, err := a.convertAPIReservationToCommon(resp.JSON200.Reservations[0])
	if err != nil {
		return nil, a.WrapError(err, "failed to convert reservation "+name)
	}

	return res, nil
}

// Create creates a new reservation
func (a *ReservationAdapter) Create(ctx context.Context, req *types.ReservationCreate) (*types.ReservationCreateResponse, error) {
	// v0.0.41 doesn't support reservation creation through the API
	return nil, fmt.Errorf("reservation creation is not supported in API v0.0.41")
}

// Update updates an existing reservation
func (a *ReservationAdapter) Update(ctx context.Context, name string, update *types.ReservationUpdate) error {
	// v0.0.41 doesn't support reservation updates through the API
	return fmt.Errorf("reservation update is not supported in API v0.0.41")
}

// Delete deletes a reservation
func (a *ReservationAdapter) Delete(ctx context.Context, name string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Validate name
	if err := a.ValidateResourceName("reservation name", name); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Note: v0.0.41 API does not support reservation deletion
	// The SlurmV0041DeleteReservationWithResponse method does not exist in the generated client
	// TODO: Use a newer API version (v0.0.42+) for reservation deletion
	return fmt.Errorf("reservation deletion not supported in v0.0.41 - use a newer API version")
}
