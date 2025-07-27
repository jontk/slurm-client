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
	reservations := make([]interfaces.Reservation, 0, len(*resp.JSON200.Reservations))
	for _, apiRes := range *resp.JSON200.Reservations {
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
	if resp.JSON200 == nil || resp.JSON200.Reservations == nil || len(*resp.JSON200.Reservations) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "Reservation not found", fmt.Sprintf("Reservation '%s' not found", reservationName))
	}

	// Convert the first reservation in the response
	reservation, err := convertAPIReservationToInterface((*resp.JSON200.Reservations)[0])
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

	// Convert interface types to API types
	apiReservation, err := convertReservationCreateToAPI(reservation)
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeInvalidRequest, "Failed to convert reservation data")
		conversionErr.Cause = err
		return nil, conversionErr
	}

	// Create the request body
	requestBody := SlurmV0042PostReservationJSONRequestBody{
		Reservation: *apiReservation,
	}

	// Call the generated OpenAPI client
	resp, err := r.client.apiClient.SlurmV0042PostReservationWithResponse(ctx, requestBody)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.42")
	}

	// Check HTTP status (200 and 201 for creation is success)
	if resp.StatusCode() != 200 && resp.StatusCode() != 201 {
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

	// Create successful response
	response := &interfaces.ReservationCreateResponse{
		Name:    reservation.Name,
		Created: true,
	}

	// If response contains additional information, include it
	if resp.JSON200 != nil && resp.JSON200.Reservations != nil && len(*resp.JSON200.Reservations) > 0 {
		createdRes := (*resp.JSON200.Reservations)[0]
		if createdRes.Name != nil {
			response.Name = *createdRes.Name
		}
	}

	return response, nil
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

	// Convert update to API format
	apiUpdate, err := convertReservationUpdateToAPI(update)
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeInvalidRequest, "Failed to convert update data")
		conversionErr.Cause = err
		return conversionErr
	}

	// Set the reservation name in the update
	apiUpdate.Name = &reservationName

	// Create the request body
	requestBody := SlurmV0042PostReservationJSONRequestBody{
		Reservation: *apiUpdate,
	}

	// Call the generated OpenAPI client (POST is used for updates in Slurm)
	resp, err := r.client.apiClient.SlurmV0042PostReservationWithResponse(ctx, requestBody)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.42")
	}

	// Check HTTP status
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
				return apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.42")
		return httpErr
	}

	return nil
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

	// State - convert from array
	if apiRes.State != nil && len(*apiRes.State) > 0 {
		reservation.State = string((*apiRes.State)[0])
	}

	// Time fields
	if apiRes.StartTime != nil && apiRes.StartTime.Set != nil && *apiRes.StartTime.Set && apiRes.StartTime.Number != nil {
		reservation.StartTime = time.Unix(*apiRes.StartTime.Number, 0)
	}

	if apiRes.EndTime != nil && apiRes.EndTime.Set != nil && *apiRes.EndTime.Set && apiRes.EndTime.Number != nil {
		reservation.EndTime = time.Unix(*apiRes.EndTime.Number, 0)
	}

	if apiRes.Duration != nil && apiRes.Duration.Set != nil && *apiRes.Duration.Set && apiRes.Duration.Number != nil {
		reservation.Duration = int(*apiRes.Duration.Number)
	}

	// Node information
	if apiRes.Node != nil {
		reservation.Nodes = strings.Split(*apiRes.Node, ",")
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
		flags := make([]string, 0, len(*apiRes.Flags))
		for _, flag := range *apiRes.Flags {
			flags = append(flags, string(flag))
		}
		reservation.Flags = flags
	}

	// Features
	if apiRes.Features != nil {
		reservation.Features = strings.Split(*apiRes.Features, ",")
	}

	// Partitions
	if apiRes.Partition != nil {
		reservation.Partitions = strings.Split(*apiRes.Partition, ",")
	}

	// Additional fields
	if apiRes.Licenses != nil {
		reservation.Licenses = *apiRes.Licenses
	}

	if apiRes.Groups != nil {
		reservation.Groups = strings.Split(*apiRes.Groups, ",")
	}

	if apiRes.MaxStartDelay != nil && apiRes.MaxStartDelay.Set != nil && *apiRes.MaxStartDelay.Set && apiRes.MaxStartDelay.Number != nil {
		maxStartDelay := int(*apiRes.MaxStartDelay.Number)
		reservation.MaxStartDelay = &maxStartDelay
	}

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

		// Filter by state
		if opts.State != "" && res.State != opts.State {
			continue
		}

		// Filter by time range
		if !opts.StartTime.IsZero() && res.StartTime.Before(opts.StartTime) {
			continue
		}

		if !opts.EndTime.IsZero() && res.EndTime.After(opts.EndTime) {
			continue
		}

		filtered = append(filtered, res)
	}

	return filtered
}

// convertReservationCreateToAPI converts interfaces.ReservationCreate to API format
func convertReservationCreateToAPI(create *interfaces.ReservationCreate) (*V0042ReservationInfo, error) {
	apiRes := &V0042ReservationInfo{}

	// Required fields
	apiRes.Name = &create.Name

	// Time fields
	startTime := create.StartTime.Unix()
	apiRes.StartTime = &V0042Int64NoValStruct{
		Set:    &[]bool{true}[0],
		Number: &startTime,
	}

	if !create.EndTime.IsZero() {
		endTime := create.EndTime.Unix()
		apiRes.EndTime = &V0042Int64NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &endTime,
		}
	} else if create.Duration > 0 {
		duration := int32(create.Duration)
		apiRes.Duration = &V0042Uint32NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &duration,
		}
	}

	// Node specifications
	if len(create.Nodes) > 0 {
		nodeStr := strings.Join(create.Nodes, ",")
		apiRes.Node = &nodeStr
	}

	if create.NodeCount > 0 {
		nodeCount := int32(create.NodeCount)
		apiRes.NodeCnt = &V0042Uint32NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &nodeCount,
		}
	}

	if create.CoreCount > 0 {
		coreCount := int32(create.CoreCount)
		apiRes.CoreCnt = &V0042Uint32NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &coreCount,
		}
	}

	// Users and accounts
	if len(create.Users) > 0 {
		userStr := strings.Join(create.Users, ",")
		apiRes.Users = &userStr
	}

	if len(create.Accounts) > 0 {
		accountStr := strings.Join(create.Accounts, ",")
		apiRes.Accounts = &accountStr
	}

	// Flags
	if len(create.Flags) > 0 {
		flags := make(V0042ReservationFlags, 0, len(create.Flags))
		for _, flag := range create.Flags {
			flags = append(flags, flag)
		}
		apiRes.Flags = &flags
	}

	// Features
	if len(create.Features) > 0 {
		featureStr := strings.Join(create.Features, ",")
		apiRes.Features = &featureStr
	}

	// Partitions
	if len(create.Partitions) > 0 {
		partitionStr := strings.Join(create.Partitions, ",")
		apiRes.Partition = &partitionStr
	}

	// Additional fields
	if create.Licenses != "" {
		apiRes.Licenses = &create.Licenses
	}

	if len(create.Groups) > 0 {
		groupStr := strings.Join(create.Groups, ",")
		apiRes.Groups = &groupStr
	}

	if create.MaxStartDelay != nil {
		maxStartDelay := int32(*create.MaxStartDelay)
		apiRes.MaxStartDelay = &V0042Uint32NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &maxStartDelay,
		}
	}

	return apiRes, nil
}

// convertReservationUpdateToAPI converts interfaces.ReservationUpdate to API format
func convertReservationUpdateToAPI(update *interfaces.ReservationUpdate) (*V0042ReservationInfo, error) {
	apiRes := &V0042ReservationInfo{}

	// Time fields (only if specified in update)
	if update.StartTime != nil && !update.StartTime.IsZero() {
		startTime := update.StartTime.Unix()
		apiRes.StartTime = &V0042Int64NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &startTime,
		}
	}

	if update.EndTime != nil && !update.EndTime.IsZero() {
		endTime := update.EndTime.Unix()
		apiRes.EndTime = &V0042Int64NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &endTime,
		}
	}

	if update.Duration != nil {
		duration := int32(*update.Duration)
		apiRes.Duration = &V0042Uint32NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &duration,
		}
	}

	// Node specifications
	if update.Nodes != nil {
		nodeStr := strings.Join(*update.Nodes, ",")
		apiRes.Node = &nodeStr
	}

	if update.NodeCount != nil {
		nodeCount := int32(*update.NodeCount)
		apiRes.NodeCnt = &V0042Uint32NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &nodeCount,
		}
	}

	if update.CoreCount != nil {
		coreCount := int32(*update.CoreCount)
		apiRes.CoreCnt = &V0042Uint32NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &coreCount,
		}
	}

	// Users and accounts
	if update.Users != nil {
		userStr := strings.Join(*update.Users, ",")
		apiRes.Users = &userStr
	}

	if update.Accounts != nil {
		accountStr := strings.Join(*update.Accounts, ",")
		apiRes.Accounts = &accountStr
	}

	// Flags
	if update.Flags != nil {
		flags := make(V0042ReservationFlags, 0, len(*update.Flags))
		for _, flag := range *update.Flags {
			flags = append(flags, flag)
		}
		apiRes.Flags = &flags
	}

	// Features
	if update.Features != nil {
		featureStr := strings.Join(*update.Features, ",")
		apiRes.Features = &featureStr
	}

	return apiRes, nil
}