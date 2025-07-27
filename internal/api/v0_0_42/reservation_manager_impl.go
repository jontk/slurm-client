package v0_0_42

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// ReservationManagerImpl implements the ReservationManager interface for v0.0.42
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
	if r.client == nil || r.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmV0042GetReservationsParams{}

	// Call the generated OpenAPI client
	resp, err := r.client.apiClient.SlurmV0042GetReservationsWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.42")
	}

	// Check HTTP status and handle API errors
	if resp.StatusCode() != 200 {
		var responseBody []byte
		if resp.JSON200 != nil {
			// Try to extract error details from response
			if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
				apiErrors := make([]errors.SlurmAPIErrorDetail, len(*resp.JSON200.Errors))
				for i, apiErr := range *resp.JSON200.Errors {
					var errorNumber int
					if apiErr.ErrorNumber != nil {
						errorNumber = int(*apiErr.ErrorNumber)
					}
					var errorCode string
					if apiErr.Error != nil {
						errorCode = *apiErr.Error
					}
					var source string
					if apiErr.Source != nil {
						source = *apiErr.Source
					}
					var description string
					if apiErr.Description != nil {
						description = *apiErr.Description
					}

					apiErrors[i] = errors.SlurmAPIErrorDetail{
						ErrorNumber: errorNumber,
						ErrorCode:   errorCode,
						Source:      source,
						Description: description,
					}
				}
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.42", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.42")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil || resp.JSON200.Reservations == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response with reservations but got nil")
	}

	// Convert the response to our interface types
	reservations := make([]interfaces.Reservation, 0, len(resp.JSON200.Reservations))
	for _, apiRes := range resp.JSON200.Reservations {
		reservation, err := convertAPIReservationToInterface(apiRes)
		if err != nil {
			conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert reservation data")
			conversionErr.Cause = err
			conversionErr.Details = fmt.Sprintf("Error converting reservation %v", apiRes.Name)
			return nil, conversionErr
		}
		reservations = append(reservations, *reservation)
	}

	// Apply client-side filtering if options are provided
	if opts != nil {
		reservations = filterReservations(reservations, opts)
	}

	return &interfaces.ReservationList{
		Reservations: reservations,
		Total:        len(reservations),
	}, nil
}

// Get retrieves a specific reservation by name
func (r *ReservationManagerImpl) Get(ctx context.Context, reservationName string) (*interfaces.Reservation, error) {
	if r.client == nil || r.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if reservationName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "reservation name is required", "reservationName", reservationName, nil)
	}

	// Prepare parameters for the API call
	params := &SlurmV0042GetReservationParams{}

	// Call the generated OpenAPI client
	resp, err := r.client.apiClient.SlurmV0042GetReservationWithResponse(ctx, reservationName, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.42")
	}

	// Check HTTP status and handle API errors
	if resp.StatusCode() != 200 {
		var responseBody []byte
		if resp.JSON200 != nil {
			// Try to extract error details from response
			if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
				apiErrors := make([]errors.SlurmAPIErrorDetail, len(*resp.JSON200.Errors))
				for i, apiErr := range *resp.JSON200.Errors {
					var errorNumber int
					if apiErr.ErrorNumber != nil {
						errorNumber = int(*apiErr.ErrorNumber)
					}
					var errorCode string
					if apiErr.Error != nil {
						errorCode = *apiErr.Error
					}
					var source string
					if apiErr.Source != nil {
						source = *apiErr.Source
					}
					var description string
					if apiErr.Description != nil {
						description = *apiErr.Description
					}

					apiErrors[i] = errors.SlurmAPIErrorDetail{
						ErrorNumber: errorNumber,
						ErrorCode:   errorCode,
						Source:      source,
						Description: description,
					}
				}
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.42", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.42")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil || resp.JSON200.Reservations == nil || len(resp.JSON200.Reservations) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "Reservation not found", fmt.Sprintf("Reservation '%s' not found", reservationName))
	}

	// Convert the first reservation in the response
	reservation, err := convertAPIReservationToInterface(resp.JSON200.Reservations[0])
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert reservation data")
		conversionErr.Cause = err
		conversionErr.Details = fmt.Sprintf("Error converting reservation '%s'", reservationName)
		return nil, conversionErr
	}

	return reservation, nil
}

// Create creates a new reservation
func (r *ReservationManagerImpl) Create(ctx context.Context, reservation *interfaces.ReservationCreate) (*interfaces.ReservationCreateResponse, error) {
	if r.client == nil || r.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if reservation == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "reservation data is required", "reservation", reservation, nil)
	}

	// Validate required fields
	if reservation.Name == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "reservation name is required", "reservation.Name", reservation.Name, nil)
	}
	if reservation.StartTime.IsZero() {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "reservation start time is required", "reservation.StartTime", reservation.StartTime, nil)
	}
	if reservation.EndTime.IsZero() && reservation.Duration == 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "either end time or duration is required", "reservation.EndTime/Duration", fmt.Sprintf("EndTime: %v, Duration: %v", reservation.EndTime, reservation.Duration), nil)
	}

	// For v0.0.42, reservation creation is not supported via API
	// Return appropriate error
	return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "Reservation creation not supported in v0.0.42 API")

}

// Update updates an existing reservation
func (r *ReservationManagerImpl) Update(ctx context.Context, reservationName string, update *interfaces.ReservationUpdate) error {
	if r.client == nil || r.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if reservationName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "reservation name is required", "reservationName", reservationName, nil)
	}

	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "update data is required", "update", update, nil)
	}

	// For v0.0.42, reservation updates are not supported via API
	// Return appropriate error
	return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "Reservation updates not supported in v0.0.42 API")
}

// Delete deletes a reservation
func (r *ReservationManagerImpl) Delete(ctx context.Context, reservationName string) error {
	if r.client == nil || r.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if reservationName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "reservation name is required", "reservationName", reservationName, nil)
	}

	// Call the generated OpenAPI client
	resp, err := r.client.apiClient.SlurmV0042DeleteReservationWithResponse(ctx, reservationName)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.42")
	}

	// Check HTTP status (200 or 204 for successful deletion)
	if resp.StatusCode() != 200 && resp.StatusCode() != 204 {
		var responseBody []byte
		if resp.JSON200 != nil {
			// Try to extract error details from response
			if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
				apiErrors := make([]errors.SlurmAPIErrorDetail, len(*resp.JSON200.Errors))
				for i, apiErr := range *resp.JSON200.Errors {
					var errorNumber int
					if apiErr.ErrorNumber != nil {
						errorNumber = int(*apiErr.ErrorNumber)
					}
					var errorCode string
					if apiErr.Error != nil {
						errorCode = *apiErr.Error
					}
					var source string
					if apiErr.Source != nil {
						source = *apiErr.Source
					}
					var description string
					if apiErr.Description != nil {
						description = *apiErr.Description
					}

					apiErrors[i] = errors.SlurmAPIErrorDetail{
						ErrorNumber: errorNumber,
						ErrorCode:   errorCode,
						Source:      source,
						Description: description,
					}
				}
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.42", apiErrors)
				return apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.42")
		return httpErr
	}

	return nil
}

// convertAPIReservationToInterface converts V0042ReservationInfo to interfaces.Reservation
func convertAPIReservationToInterface(apiRes V0042ReservationInfo) (*interfaces.Reservation, error) {
	reservation := &interfaces.Reservation{}

	// Basic fields
	if apiRes.Name != nil {
		reservation.Name = *apiRes.Name
	}

	// v0.0.42 doesn't have State field, set a default
	reservation.State = "UNKNOWN"

	// Time fields
	if apiRes.StartTime != nil && apiRes.StartTime.Set != nil && *apiRes.StartTime.Set && apiRes.StartTime.Number != nil {
		reservation.StartTime = time.Unix(int64(*apiRes.StartTime.Number), 0)
	}

	if apiRes.EndTime != nil && apiRes.EndTime.Set != nil && *apiRes.EndTime.Set && apiRes.EndTime.Number != nil {
		reservation.EndTime = time.Unix(int64(*apiRes.EndTime.Number), 0)
	}

	// v0.0.42 doesn't have Duration field - calculate from start/end if available
	if !reservation.StartTime.IsZero() && !reservation.EndTime.IsZero() {
		reservation.Duration = int(reservation.EndTime.Sub(reservation.StartTime).Seconds())
	}

	// Node information - v0.0.42 doesn't have Node field, use NodeList if available
	if apiRes.NodeList != nil {
		reservation.Nodes = strings.Split(*apiRes.NodeList, ",")
	}

	if apiRes.NodeCount != nil {
		reservation.NodeCount = int(*apiRes.NodeCount)
	}

	if apiRes.CoreCount != nil {
		reservation.CoreCount = int(*apiRes.CoreCount)
	}

	// Users and accounts
	if apiRes.Users != nil {
		reservation.Users = strings.Split(*apiRes.Users, ",")
	}

	if apiRes.Accounts != nil {
		reservation.Accounts = strings.Split(*apiRes.Accounts, ",")
	}

	// Flags
	if apiRes.Flags != nil {
		reservation.Flags = *apiRes.Flags
	}

	// Features
	if apiRes.Features != nil {
		reservation.Features = strings.Split(*apiRes.Features, ",")
	}

	// Additional fields - adapt to available v0.0.42 fields
	if apiRes.Licenses != nil {
		// v0.0.42 Licenses is a string, but interface expects map[string]int
		// Skip this field for compatibility
		_ = *apiRes.Licenses
	}

	if apiRes.Groups != nil {
		_ = *apiRes.Groups // v0.0.42 interfaces don't support Groups field, ignore
	}

	// v0.0.42 doesn't have MaxStartDelay field in interfaces, skip it
	_ = apiRes.MaxStartDelay

	return reservation, nil
}

// filterReservations applies client-side filtering to the reservation list
func filterReservations(reservations []interfaces.Reservation, opts *interfaces.ListReservationsOptions) []interfaces.Reservation {
	if opts == nil {
		return reservations
	}

	filtered := make([]interfaces.Reservation, 0, len(reservations))
	for _, res := range reservations {
		// Filter by names
		if len(opts.Names) > 0 {
			found := false
			for _, name := range opts.Names {
				if res.Name == name {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filter by users
		if len(opts.Users) > 0 {
			found := false
			for _, filterUser := range opts.Users {
				for _, resUser := range res.Users {
					if resUser == filterUser {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filter by accounts
		if len(opts.Accounts) > 0 {
			found := false
			for _, filterAccount := range opts.Accounts {
				for _, resAccount := range res.Accounts {
					if resAccount == filterAccount {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				continue
			}
		}

		// v0.0.42 interfaces don't support State, StartTime, EndTime filtering
		// Skip these filters for v0.0.42 compatibility

		filtered = append(filtered, res)
	}

	return filtered
}

// convertReservationCreateToAPI - not implemented for v0.0.42 (no POST API)
func convertReservationCreateToAPI(create *interfaces.ReservationCreate) (*V0042ReservationInfo, error) {
	return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "Reservation creation not supported in v0.0.42 API")
}

// convertReservationUpdateToAPI - not implemented for v0.0.42 (no POST API)
func convertReservationUpdateToAPI(update *interfaces.ReservationUpdate) (*V0042ReservationInfo, error) {
	return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "Reservation updates not supported in v0.0.42 API")
}