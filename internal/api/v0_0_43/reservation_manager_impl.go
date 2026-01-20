// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// ReservationManagerImpl implements the ReservationManager interface for v0.0.43
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
	params := &SlurmV0043GetReservationsParams{}

	// Call the generated OpenAPI client
	resp, err := r.client.apiClient.SlurmV0043GetReservationsWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
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
	// Validate input first (cheap check)
	if reservationName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "reservation name is required", "reservationName", reservationName, nil)
	}

	// Then check client initialization
	if r.client == nil || r.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmV0043GetReservationParams{}

	// Call the generated OpenAPI client
	resp, err := r.client.apiClient.SlurmV0043GetReservationWithResponse(ctx, reservationName, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
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

	// Convert interface types to API types
	apiReservation, err := convertReservationCreateToAPI(reservation)
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeInvalidRequest, "Failed to convert reservation data")
		conversionErr.Cause = err
		return nil, conversionErr
	}

	// Create the request body - SlurmV0043PostReservationJSONRequestBody is just V0043ReservationDescMsg
	requestBody := *apiReservation

	// Call the generated OpenAPI client
	resp, err := r.client.apiClient.SlurmV0043PostReservationWithResponse(ctx, requestBody)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
		return nil, httpErr
	}

	// Create successful response
	response := &interfaces.ReservationCreateResponse{
		ReservationName: reservation.Name,
	}

	// If response contains additional information, include it
	if resp.JSON200 != nil && resp.JSON200.Reservations != nil && len(resp.JSON200.Reservations) > 0 {
		createdRes := resp.JSON200.Reservations[0]
		if createdRes.Name != nil {
			response.ReservationName = *createdRes.Name
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

	// Create the request body - SlurmV0043PostReservationJSONRequestBody is just V0043ReservationDescMsg
	requestBody := *apiUpdate

	// Call the generated OpenAPI client (POST is used for updates in Slurm)
	resp, err := r.client.apiClient.SlurmV0043PostReservationWithResponse(ctx, requestBody)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
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
	resp, err := r.client.apiClient.SlurmV0043DeleteReservationWithResponse(ctx, reservationName)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
		return httpErr
	}

	return nil
}

// convertAPIReservationToInterface converts V0043ReservationInfo to interfaces.Reservation
func convertAPIReservationToInterface(apiRes V0043ReservationInfo) (*interfaces.Reservation, error) {
	reservation := &interfaces.Reservation{}

	// Basic fields
	if apiRes.Name != nil {
		reservation.Name = *apiRes.Name
	}

	// State - not available in v0.0.43 API, use default
	reservation.State = "ACTIVE"

	// Time fields
	if apiRes.StartTime != nil && apiRes.StartTime.Set != nil && *apiRes.StartTime.Set && apiRes.StartTime.Number != nil {
		reservation.StartTime = time.Unix(*apiRes.StartTime.Number, 0)
	}

	if apiRes.EndTime != nil && apiRes.EndTime.Set != nil && *apiRes.EndTime.Set && apiRes.EndTime.Number != nil {
		reservation.EndTime = time.Unix(*apiRes.EndTime.Number, 0)
	}

	// Duration - calculate from start/end times if available
	if !reservation.StartTime.IsZero() && !reservation.EndTime.IsZero() {
		reservation.Duration = int(reservation.EndTime.Sub(reservation.StartTime).Seconds())
	}

	// Node information - use NodeList field
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

	// Partition (single partition name in interface)
	if apiRes.Partition != nil {
		reservation.PartitionName = *apiRes.Partition
	}

	// Additional fields - Licenses (skip for now, complex field mapping)
	if apiRes.Licenses != nil {
		// Skip licenses mapping - API provides string, interface expects map[string]int
		_ = *apiRes.Licenses
	}

	if apiRes.Groups != nil {
		// Groups field not supported in interface, skip
		_ = *apiRes.Groups
	}

	if apiRes.MaxStartDelay != nil {
		// MaxStartDelay field not supported in interface, skip
		_ = *apiRes.MaxStartDelay
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

		// Filter by states (array in interface)
		if len(opts.States) > 0 {
			found := false
			for _, state := range opts.States {
				if res.State == state {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Time-based filtering not available in ListReservationsOptions interface

		filtered = append(filtered, res)
	}

	return filtered
}

// convertReservationCreateToAPI converts interfaces.ReservationCreate to API format
func convertReservationCreateToAPI(create *interfaces.ReservationCreate) (*V0043ReservationDescMsg, error) {
	apiRes := &V0043ReservationDescMsg{}

	// Required fields
	apiRes.Name = &create.Name

	// Time fields - use correct types
	startTime := int64(create.StartTime.Unix())
	apiRes.StartTime = &V0043Uint64NoValStruct{
		Set:    &[]bool{true}[0],
		Number: &startTime,
	}

	if !create.EndTime.IsZero() {
		endTime := int64(create.EndTime.Unix())
		apiRes.EndTime = &V0043Uint64NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &endTime,
		}
	} else if create.Duration > 0 {
		duration := int32(create.Duration)
		apiRes.Duration = &V0043Uint32NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &duration,
		}
	}

	// Node specifications - use correct field names
	if len(create.Nodes) > 0 {
		nodeList := V0043HostlistString(create.Nodes)
		apiRes.NodeList = &nodeList
	}

	if create.NodeCount > 0 {
		nodeCount := int32(create.NodeCount)
		apiRes.NodeCount = &V0043Uint32NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &nodeCount,
		}
	}

	if create.CoreCount > 0 {
		coreCount := int32(create.CoreCount)
		apiRes.CoreCount = &V0043Uint32NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &coreCount,
		}
	}

	// Users and accounts - use V0043CsvString type
	if len(create.Users) > 0 {
		users := V0043CsvString(create.Users)
		apiRes.Users = &users
	}

	if len(create.Accounts) > 0 {
		accounts := V0043CsvString(create.Accounts)
		apiRes.Accounts = &accounts
	}

	// Flags
	if len(create.Flags) > 0 {
		flags := make([]V0043ReservationDescMsgFlags, 0, len(create.Flags))
		for _, flag := range create.Flags {
			flags = append(flags, V0043ReservationDescMsgFlags(flag))
		}
		apiRes.Flags = &flags
	}

	// Features
	if len(create.Features) > 0 {
		featureStr := strings.Join(create.Features, ",")
		apiRes.Features = &featureStr
	}

	// Partition (single partition in interface)
	if create.PartitionName != "" {
		apiRes.Partition = &create.PartitionName
	}

	// Additional fields - map Licenses from map[string]int to V0043CsvString
	if len(create.Licenses) > 0 {
		// Convert map[string]int to string slice format
		licenseStrs := make([]string, 0, len(create.Licenses))
		for name, count := range create.Licenses {
			licenseStrs = append(licenseStrs, fmt.Sprintf("%s:%d", name, count))
		}
		licenses := V0043CsvString(licenseStrs)
		apiRes.Licenses = &licenses
	}

	return apiRes, nil
}

// convertReservationUpdateToAPI converts interfaces.ReservationUpdate to API format
func convertReservationUpdateToAPI(update *interfaces.ReservationUpdate) (*V0043ReservationDescMsg, error) {
	apiRes := &V0043ReservationDescMsg{}

	// Time fields (only if specified in update) - use correct types
	if update.StartTime != nil && !update.StartTime.IsZero() {
		startTime := int64(update.StartTime.Unix())
		apiRes.StartTime = &V0043Uint64NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &startTime,
		}
	}

	if update.EndTime != nil && !update.EndTime.IsZero() {
		endTime := int64(update.EndTime.Unix())
		apiRes.EndTime = &V0043Uint64NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &endTime,
		}
	}

	if update.Duration != nil {
		duration := int32(*update.Duration)
		apiRes.Duration = &V0043Uint32NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &duration,
		}
	}

	// Node specifications - use correct field names and types
	if update.Nodes != nil {
		nodeList := V0043HostlistString(update.Nodes)
		apiRes.NodeList = &nodeList
	}

	if update.NodeCount != nil {
		nodeCount := int32(*update.NodeCount)
		apiRes.NodeCount = &V0043Uint32NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &nodeCount,
		}
	}

	// CoreCount field not available in ReservationUpdate interface, skip

	// Users and accounts - use V0043CsvString type
	if update.Users != nil {
		users := V0043CsvString(update.Users)
		apiRes.Users = &users
	}

	if update.Accounts != nil {
		accounts := V0043CsvString(update.Accounts)
		apiRes.Accounts = &accounts
	}

	// Flags
	if update.Flags != nil {
		flags := make([]V0043ReservationDescMsgFlags, 0, len(update.Flags))
		for _, flag := range update.Flags {
			flags = append(flags, V0043ReservationDescMsgFlags(flag))
		}
		apiRes.Flags = &flags
	}

	// Features
	if update.Features != nil {
		featureStr := strings.Join(update.Features, ",")
		apiRes.Features = &featureStr
	}

	return apiRes, nil
}
