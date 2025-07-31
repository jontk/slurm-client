package v0_0_42

import (
	"context"
	"fmt"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

// ReservationAdapter implements the ReservationAdapter interface for v0.0.42
type ReservationAdapter struct {
	*base.BaseManager
	client *api.ClientWithResponses
}

// NewReservationAdapter creates a new Reservation adapter for v0.0.42
func NewReservationAdapter(client *api.ClientWithResponses) *ReservationAdapter {
	return &ReservationAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Reservation"),
		client:      client,
	}
}

// List retrieves a list of reservations
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
	params := &api.SlurmV0042GetReservationsParams{}

	// Apply filters from options
	if opts != nil && len(opts.Names) > 0 {
		// v0.0.42 doesn't support reservation name filtering in the API params,
		// we'll need to filter client-side
	}

	// Call the API
	resp, err := a.client.SlurmV0042GetReservationsWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list reservations")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleHTTPResponse(resp.HTTPResponse, resp.Body)
	}

	// Check for API response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from API")
	}

	// Convert the response to common types
	reservationList := &types.ReservationList{
		Reservations: make([]types.Reservation, 0),
	}

	if resp.JSON200.Reservations != nil {
		for _, apiReservation := range resp.JSON200.Reservations {
			reservation, err := a.convertAPIReservationToCommon(apiReservation)
			if err != nil {
				// Log conversion error but continue
				continue
			}
			
			// Apply client-side filtering if needed
			if opts != nil && len(opts.Names) > 0 {
				found := false
				for _, name := range opts.Names {
					if reservation.Name == name {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
			
			reservationList.Reservations = append(reservationList.Reservations, *reservation)
		}
	}

	return reservationList, nil
}

// Get retrieves a specific reservation by name
func (a *ReservationAdapter) Get(ctx context.Context, name string) (*types.Reservation, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters
	params := &api.SlurmV0042GetReservationParams{}

	// Call the API
	resp, err := a.client.SlurmV0042GetReservationWithResponse(ctx, name, params)
	if err != nil {
		return nil, a.WrapError(err, fmt.Sprintf("failed to get reservation %s", name))
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleHTTPResponse(resp.HTTPResponse, resp.Body)
	}

	// Check for API response
	if resp.JSON200 == nil || resp.JSON200.Reservations == nil || len(resp.JSON200.Reservations) == 0 {
		return nil, fmt.Errorf("reservation %s not found", name)
	}

	// Convert the first reservation in the response
	reservations := resp.JSON200.Reservations
	return a.convertAPIReservationToCommon(reservations[0])
}

// Create creates a new reservation
func (a *ReservationAdapter) Create(ctx context.Context, reservation *types.ReservationCreate) (*types.ReservationCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// v0.0.42 doesn't support reservation creation via API
	return nil, fmt.Errorf("reservation creation not supported via v0.0.42 REST API")
}

// Update updates an existing reservation
func (a *ReservationAdapter) Update(ctx context.Context, name string, updates *types.ReservationUpdateRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// v0.0.42 doesn't have a reservation update endpoint
	// Updates require delete and recreate
	return fmt.Errorf("reservation update not supported via v0.0.42 API - use delete and recreate")
}

// Delete deletes a reservation
func (a *ReservationAdapter) Delete(ctx context.Context, name string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the API
	resp, err := a.client.SlurmV0042DeleteReservationWithResponse(ctx, name)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to delete reservation %s", name))
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleHTTPResponse(resp.HTTPResponse, resp.Body)
	}

	return nil
}