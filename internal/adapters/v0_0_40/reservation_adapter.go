// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_40

import (
	"context"
	"strconv"

	types "github.com/jontk/slurm-client/api"
	adapterbase "github.com/jontk/slurm-client/internal/adapters/base"
	"github.com/jontk/slurm-client/internal/common"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_40"
	"github.com/jontk/slurm-client/pkg/errors"
)

// ReservationAdapter implements the ReservationAdapter interface for v0.0.40
type ReservationAdapter struct {
	*adapterbase.BaseManager
	client *api.ClientWithResponses
}

// NewReservationAdapter creates a new Reservation adapter for v0.0.40
func NewReservationAdapter(client *api.ClientWithResponses) *ReservationAdapter {
	return &ReservationAdapter{
		BaseManager: adapterbase.NewBaseManager("v0.0.40", "Reservation"),
		client:      client,
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
	params := &api.SlurmV0040GetReservationsParams{}
	// Apply filters from options
	if opts != nil {
		// v0.0.40 doesn't support reservation name filtering in params
		// We'll need to filter client-side
		if opts.UpdateTime != nil {
			updateTimeStr := strconv.FormatInt(opts.UpdateTime.Unix(), 10)
			params.UpdateTime = &updateTimeStr
		}
	}
	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0040GetReservationsWithResponse(ctx, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}
	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.40"); err != nil {
		return nil, err
	}
	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "List Reservations"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Reservations, "List Reservations - reservations field"); err != nil {
		return nil, err
	}
	// Convert the response to common types
	reservationList := make([]types.Reservation, 0, len(resp.JSON200.Reservations))
	for _, apiReservation := range resp.JSON200.Reservations {
		reservation := a.convertAPIReservationToCommon(apiReservation)
		// Apply client-side filtering if needed
		if !a.reservationPassesNameFilter(reservation, opts) {
			continue
		}
		reservationList = append(reservationList, *reservation)
	}
	// Apply pagination
	return a.paginateReservationList(reservationList, opts), nil
}

// reservationPassesNameFilter checks if a reservation passes the name filter
func (a *ReservationAdapter) reservationPassesNameFilter(reservation *types.Reservation, opts *types.ReservationListOptions) bool {
	if opts == nil || len(opts.Names) == 0 {
		return true
	}
	for _, name := range opts.Names {
		if reservation.Name != nil && *reservation.Name == name {
			return true
		}
	}
	return false
}

// paginateReservationList applies pagination to a reservation list
func (a *ReservationAdapter) paginateReservationList(reservationList []types.Reservation, opts *types.ReservationListOptions) *types.ReservationList {
	listOpts := adapterbase.ListOptions{}
	if opts != nil {
		listOpts.Limit = opts.Limit
		listOpts.Offset = opts.Offset
	}
	start := listOpts.Offset
	if start < 0 {
		start = 0
	}
	if start >= len(reservationList) {
		return &types.ReservationList{
			Reservations: []types.Reservation{},
			Total:        len(reservationList),
		}
	}
	end := len(reservationList)
	if listOpts.Limit > 0 {
		end = start + listOpts.Limit
		if end > len(reservationList) {
			end = len(reservationList)
		}
	}
	return &types.ReservationList{
		Reservations: reservationList[start:end],
		Total:        len(reservationList),
	}
}

// Get retrieves a specific reservation by name
func (a *ReservationAdapter) Get(ctx context.Context, reservationName string) (*types.Reservation, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.ValidateResourceName(reservationName, "reservation name"); err != nil {
		return nil, err
	}
	// v0.0.40 doesn't have a single reservation GET endpoint
	// We need to list all and filter
	list, err := a.List(ctx, &types.ReservationListOptions{
		Names: []string{reservationName},
	})
	if err != nil {
		return nil, err
	}
	if len(list.Reservations) == 0 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, "Reservation '"+reservationName+"' not found")
	}
	return &list.Reservations[0], nil
}

// Create creates a new reservation
func (a *ReservationAdapter) Create(ctx context.Context, reservation *types.ReservationCreate) (*types.ReservationCreateResponse, error) {
	// v0.0.40 doesn't support reservation creation
	return nil, errors.NewNotImplementedError("reservation creation", "v0.0.40")
}

// Update updates an existing reservation
func (a *ReservationAdapter) Update(ctx context.Context, reservationName string, update *types.ReservationUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(reservationName, "reservation name"); err != nil {
		return err
	}
	// v0.0.40 doesn't support reservation updates
	return errors.NewNotImplementedError("reservation updates", "v0.0.40")
}

// Delete deletes a reservation
func (a *ReservationAdapter) Delete(ctx context.Context, reservationName string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(reservationName, "reservation name"); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}
	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0040DeleteReservationWithResponse(ctx, reservationName)
	if err != nil {
		return a.HandleAPIError(err)
	}
	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}
	// Create adapter with special handling for 204 (No Content) status
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	// For DELETE operations, 204 is also a success
	if resp.StatusCode() == 204 {
		return nil
	}
	return common.HandleAPIResponse(responseAdapter, "v0.0.40")
}
