package v0_0_43

import (
	"context"
	"strings"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// ReservationAdapter implements the ReservationAdapter interface for v0.0.43
type ReservationAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewReservationAdapter creates a new Reservation adapter for v0.0.43
func NewReservationAdapter(client *api.ClientWithResponses) *ReservationAdapter {
	return &ReservationAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Reservation"),
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
	params := &api.SlurmV0043GetReservationsParams{}

	// Apply filters from options
	if opts != nil {
		if len(opts.Names) > 0 {
			nameStr := strings.Join(opts.Names, ",")
			params.ReservationName = &nameStr
		}
		if opts.UpdateTime != nil {
			updateTime := opts.UpdateTime.Unix()
			params.UpdateTime = &updateTime
		}
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0043GetReservationsWithResponse(ctx, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
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
	reservationList := make([]types.Reservation, 0, len(*resp.JSON200.Reservations))
	for _, apiReservation := range *resp.JSON200.Reservations {
		reservation, err := a.convertAPIReservationToCommon(apiReservation)
		if err != nil {
			return nil, a.HandleConversionError(err, apiReservation.Name)
		}
		reservationList = append(reservationList, *reservation)
	}

	// Apply pagination
	listOpts := base.ListOptions{}
	if opts != nil {
		listOpts.Limit = opts.Limit
		listOpts.Offset = opts.Offset
	}

	// Apply pagination
	start := listOpts.Offset
	if start < 0 {
		start = 0
	}
	if start >= len(reservationList) {
		return &types.ReservationList{
			Reservations: []types.Reservation{},
			Total:        len(reservationList),
		}, nil
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
	}, nil
}

// Get retrieves a specific reservation by name
func (a *ReservationAdapter) Get(ctx context.Context, reservationName string) (*types.Reservation, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.ValidateResourceName(reservationName, "reservationName"); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmV0043GetReservationParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0043GetReservationWithResponse(ctx, reservationName, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "Get Reservation"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Reservations, "Get Reservation - reservations field"); err != nil {
		return nil, err
	}

	// Check if we got any reservation entries
	if len(*resp.JSON200.Reservations) == 0 {
		return nil, common.NewResourceNotFoundError("Reservation", reservationName)
	}

	// Convert the first reservation (should be the only one)
	reservation, err := a.convertAPIReservationToCommon((*resp.JSON200.Reservations)[0])
	if err != nil {
		return nil, a.HandleConversionError(err, reservationName)
	}

	return reservation, nil
}

// Create creates a new reservation
func (a *ReservationAdapter) Create(ctx context.Context, reservation *types.ReservationCreate) (*types.ReservationCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.validateReservationCreate(reservation); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert to API format
	apiReservation, err := a.convertCommonReservationCreateToAPI(reservation)
	if err != nil {
		return nil, err
	}

	// Create request body
	reqBody := api.SlurmV0043PostReservationJSONRequestBody{
		Reservations: []api.V0043ReservationInfo{*apiReservation},
	}

	// Prepare parameters for the API call
	params := &api.SlurmV0043PostReservationParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0043PostReservationWithResponse(ctx, params, reqBody)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
		return nil, err
	}

	return &types.ReservationCreateResponse{
		ReservationName: reservation.Name,
	}, nil
}

// Update updates an existing reservation
func (a *ReservationAdapter) Update(ctx context.Context, reservationName string, update *types.ReservationUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(reservationName, "reservationName"); err != nil {
		return err
	}
	if err := a.validateReservationUpdate(update); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// First, get the existing reservation to merge updates
	existingReservation, err := a.Get(ctx, reservationName)
	if err != nil {
		return err
	}

	// Convert to API format and apply updates
	apiReservation, err := a.convertCommonReservationUpdateToAPI(existingReservation, update)
	if err != nil {
		return err
	}

	// Create request body
	reqBody := api.SlurmV0043PostReservationJSONRequestBody{
		Reservations: []api.V0043ReservationInfo{*apiReservation},
	}

	// Prepare parameters for the API call
	params := &api.SlurmV0043PostReservationParams{}

	// Call the generated OpenAPI client (POST is used for updates in SLURM API)
	resp, err := a.client.SlurmV0043PostReservationWithResponse(ctx, params, reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "v0.0.43")
}

// Delete deletes a reservation
func (a *ReservationAdapter) Delete(ctx context.Context, reservationName string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(reservationName, "reservationName"); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0043DeleteReservationWithResponse(ctx, reservationName)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	// Create adapter with special handling for 204 (No Content) status
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)

	// For DELETE operations, 204 is also a success
	if resp.StatusCode() == 204 {
		return nil
	}

	return common.HandleAPIResponse(responseAdapter, "v0.0.43")
}

// validateReservationCreate validates reservation creation request
func (a *ReservationAdapter) validateReservationCreate(reservation *types.ReservationCreate) error {
	if reservation == nil {
		return common.NewValidationError("reservation creation data is required", "reservation", nil)
	}
	if reservation.Name == "" {
		return common.NewValidationError("reservation name is required", "name", reservation.Name)
	}
	if reservation.StartTime.IsZero() {
		return common.NewValidationError("start time is required", "startTime", reservation.StartTime)
	}
	if reservation.EndTime.IsZero() {
		return common.NewValidationError("end time is required", "endTime", reservation.EndTime)
	}
	if reservation.StartTime.After(reservation.EndTime) {
		return common.NewValidationError("start time cannot be after end time", "startTime/endTime", nil)
	}
	return nil
}

// validateReservationUpdate validates reservation update request
func (a *ReservationAdapter) validateReservationUpdate(update *types.ReservationUpdate) error {
	if update == nil {
		return common.NewValidationError("reservation update data is required", "update", nil)
	}
	// At least one field should be provided for update
	if update.StartTime == nil && update.EndTime == nil && update.Duration == nil &&
	   len(update.Users) == 0 && len(update.Accounts) == 0 && update.NodeCount == nil {
		return common.NewValidationError("at least one field must be provided for update", "update", update)
	}
	if update.StartTime != nil && update.EndTime != nil && update.StartTime.After(*update.EndTime) {
		return common.NewValidationError("start time cannot be after end time", "startTime/endTime", nil)
	}
	return nil
}

// Simplified converter methods for reservation management
func (a *ReservationAdapter) convertAPIReservationToCommon(apiReservation api.V0043ReservationInfo) (*types.Reservation, error) {
	reservation := &types.Reservation{}
	if apiReservation.Name != nil {
		reservation.Name = *apiReservation.Name
	}
	// TODO: Add more field conversions as needed
	return reservation, nil
}

func (a *ReservationAdapter) convertCommonReservationCreateToAPI(create *types.ReservationCreate) (*api.V0043ReservationInfo, error) {
	apiReservation := &api.V0043ReservationInfo{}
	apiReservation.Name = &create.Name
	// TODO: Add more field conversions as needed
	return apiReservation, nil
}

func (a *ReservationAdapter) convertCommonReservationUpdateToAPI(existing *types.Reservation, update *types.ReservationUpdate) (*api.V0043ReservationInfo, error) {
	apiReservation := &api.V0043ReservationInfo{}
	apiReservation.Name = &existing.Name
	// TODO: Add more field conversions as needed
	return apiReservation, nil
}