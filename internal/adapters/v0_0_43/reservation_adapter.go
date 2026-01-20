// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	"fmt"
	"strings"

	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
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
		// Note: v0.0.43 doesn't have a ReservationName parameter for filtering
		// We'll have to filter client-side
		if opts.UpdateTime != nil {
			updateTimeStr := opts.UpdateTime.Format("2006-01-02T15:04:05")
			params.UpdateTime = &updateTimeStr
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
	reservationList := make([]types.Reservation, 0, len(resp.JSON200.Reservations))
	for _, apiReservation := range resp.JSON200.Reservations {
		reservation, err := a.convertAPIReservationToCommon(apiReservation)
		if err != nil {
			return nil, a.HandleConversionError(err, apiReservation.Name)
		}
		reservationList = append(reservationList, *reservation)
	}

	// Apply client-side filtering if needed
	if opts != nil {
		reservationList = a.filterReservationList(reservationList, opts)
	}

	// Apply pagination
	listOpts := base.ListOptions{}
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
	if err := a.ValidateResourceName(reservationName, "reservation name"); err != nil {
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
	if len(resp.JSON200.Reservations) == 0 {
		return nil, common.NewResourceNotFoundError("Reservation", reservationName)
	}

	// Convert the first reservation (should be the only one)
	reservation, err := a.convertAPIReservationToCommon(resp.JSON200.Reservations[0])
	if err != nil {
		return nil, a.HandleConversionError(err, reservationName)
	}

	return reservation, nil
}

// Create creates a new reservation
func (a *ReservationAdapter) Create(ctx context.Context, reservation *types.ReservationCreate) (*types.ReservationCreateResponse, error) {
	// Perform basic validation first (cheap checks)
	if reservation == nil {
		return nil, common.NewValidationError("reservation creation data is required", "reservation", nil)
	}
	if reservation.Name == "" {
		return nil, common.NewValidationError("reservation name is required", "name", reservation.Name)
	}

	// Check client initialization before expensive operations
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Perform remaining validation
	if err := a.validateReservationCreateAdvanced(reservation); err != nil {
		return nil, err
	}

	// Convert to API format
	apiReservation, err := a.convertCommonReservationCreateToAPI(reservation)
	if err != nil {
		return nil, err
	}

	// Create request body - PostReservation expects a V0043ReservationDescMsg
	apiReservationDesc, err := a.convertAPIReservationInfoToDescMsg(apiReservation)
	if err != nil {
		return nil, err
	}
	reqBody := *apiReservationDesc

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0043PostReservationWithResponse(ctx, reqBody)
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
	if err := a.ValidateResourceName(reservationName, "reservation name"); err != nil {
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

	// Create request body - PostReservation expects a V0043ReservationDescMsg
	apiReservationDesc, err := a.convertAPIReservationInfoToDescMsg(apiReservation)
	if err != nil {
		return err
	}
	reqBody := *apiReservationDesc

	// Call the generated OpenAPI client (POST is used for updates in SLURM API)
	resp, err := a.client.SlurmV0043PostReservationWithResponse(ctx, reqBody)
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
	if err := a.ValidateResourceName(reservationName, "reservation name"); err != nil {
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

// validateReservationCreate validates reservation creation request (basic checks only)
// This is kept for backwards compatibility with other callers
func (a *ReservationAdapter) validateReservationCreate(reservation *types.ReservationCreate) error {
	if reservation == nil {
		return common.NewValidationError("reservation creation data is required", "reservation", nil)
	}
	if reservation.Name == "" {
		return common.NewValidationError("reservation name is required", "name", reservation.Name)
	}
	return a.validateReservationCreateAdvanced(reservation)
}

// validateReservationCreateAdvanced validates reservation creation time fields
func (a *ReservationAdapter) validateReservationCreateAdvanced(reservation *types.ReservationCreate) error {
	if reservation.StartTime.IsZero() {
		return common.NewValidationError("start time is required", "startTime", reservation.StartTime)
	}
	if reservation.EndTime == nil || reservation.EndTime.IsZero() {
		return common.NewValidationError("end time is required", "endTime", reservation.EndTime)
	}
	if reservation.StartTime.After(*reservation.EndTime) {
		return common.NewValidationError("start time cannot be after end time", "startTime/endTime", nil)
	}
	return nil
}

// validateReservationUpdate validates reservation update request
func (a *ReservationAdapter) validateReservationUpdate(update *types.ReservationUpdate) error {
	if update == nil {
		return common.NewValidationError("reservation update data is required", "update", nil)
	}
	// Empty updates are allowed - the API will handle no-op updates
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

	// Required: Set reservation name
	apiReservation.Name = &create.Name

	// Convert start time to Unix timestamp
	if !create.StartTime.IsZero() {
		startTime := create.StartTime.Unix()
		apiReservation.StartTime = &api.V0043Uint64NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &startTime,
		}
	}

	// Convert end time to Unix timestamp
	if create.EndTime != nil && !create.EndTime.IsZero() {
		endTime := create.EndTime.Unix()
		apiReservation.EndTime = &api.V0043Uint64NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &endTime,
		}
	}

	// Set node count if specified
	if create.NodeCount > 0 {
		apiReservation.NodeCount = &create.NodeCount
	}

	// Set node list if specified
	if create.NodeList != "" {
		apiReservation.NodeList = &create.NodeList
	}

	// Set users list
	if len(create.Users) > 0 {
		usersList := strings.Join(create.Users, ",")
		apiReservation.Users = &usersList
	}

	// Set accounts if specified
	if len(create.Accounts) > 0 {
		accountsList := strings.Join(create.Accounts, ",")
		apiReservation.Accounts = &accountsList
	}

	// Note: Partition field might not exist in ReservationCreate
	// This would typically be set via node list or other mechanisms

	// Set features if specified
	if len(create.Features) > 0 {
		featuresStr := strings.Join(create.Features, "&")
		apiReservation.Features = &featuresStr
	}

	// Set licenses if specified
	if len(create.Licenses) > 0 {
		// Convert map[string]int32 to comma-separated string
		licensesList := make([]string, 0, len(create.Licenses))
		for lic, count := range create.Licenses {
			licensesList = append(licensesList, fmt.Sprintf("%s:%d", lic, count))
		}
		licensesStr := strings.Join(licensesList, ",")
		apiReservation.Licenses = &licensesStr
	}

	// Set flags if specified
	if len(create.Flags) > 0 {
		flags := make([]api.V0043ReservationInfoFlags, len(create.Flags))
		for i, flag := range create.Flags {
			flags[i] = api.V0043ReservationInfoFlags(flag)
		}
		apiReservation.Flags = &flags
	}

	return apiReservation, nil
}

func (a *ReservationAdapter) convertCommonReservationUpdateToAPI(existing *types.Reservation, update *types.ReservationUpdate) (*api.V0043ReservationInfo, error) {
	apiReservation := &api.V0043ReservationInfo{}
	apiReservation.Name = &existing.Name
	// TODO: Add more field conversions as needed
	return apiReservation, nil
}

// filterReservationList applies client-side filtering to the reservation list
func (a *ReservationAdapter) filterReservationList(reservations []types.Reservation, opts *types.ReservationListOptions) []types.Reservation {
	if opts == nil || len(opts.Names) == 0 {
		return reservations
	}

	// Create a map for quick lookup
	nameFilter := make(map[string]bool)
	for _, name := range opts.Names {
		nameFilter[name] = true
	}

	// Filter reservations by name
	var filtered []types.Reservation
	for _, reservation := range reservations {
		if nameFilter[reservation.Name] {
			filtered = append(filtered, reservation)
		}
	}

	return filtered
}

// convertAPIReservationInfoToDescMsg converts V0043ReservationInfo to V0043ReservationDescMsg
func (a *ReservationAdapter) convertAPIReservationInfoToDescMsg(info *api.V0043ReservationInfo) (*api.V0043ReservationDescMsg, error) {
	// Create a new V0043ReservationDescMsg
	descMsg := &api.V0043ReservationDescMsg{}

	// Convert fields from ReservationInfo to ReservationDescMsg
	if info.Name != nil {
		descMsg.Name = info.Name
	}

	// Convert time fields
	if info.StartTime != nil {
		descMsg.StartTime = info.StartTime
	}

	if info.EndTime != nil {
		descMsg.EndTime = info.EndTime
	}

	// Note: Duration field might not exist in ReservationInfo
	// Duration is typically calculated from start/end times

	// Convert node count
	if info.NodeCount != nil {
		nodeCount := int32(*info.NodeCount)
		setTrue := true
		descMsg.NodeCount = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &nodeCount,
		}
	}

	// Convert node list
	if info.NodeList != nil {
		// V0043HostlistString is []string
		// Convert comma-separated string to slice
		nodeList := strings.Split(*info.NodeList, ",")
		hostList := api.V0043HostlistString(nodeList)
		descMsg.NodeList = &hostList
	}

	// Convert users list
	if info.Users != nil {
		// V0043CsvString is []string
		// Convert comma-separated string to slice
		usersList := strings.Split(*info.Users, ",")
		csvUsers := api.V0043CsvString(usersList)
		descMsg.Users = &csvUsers
	}

	// Convert accounts list
	if info.Accounts != nil {
		// V0043CsvString is []string
		// Convert comma-separated string to slice
		accountsList := strings.Split(*info.Accounts, ",")
		csvAccounts := api.V0043CsvString(accountsList)
		descMsg.Accounts = &csvAccounts
	}

	// Convert partition
	if info.Partition != nil {
		descMsg.Partition = info.Partition
	}

	// Convert features
	if info.Features != nil {
		descMsg.Features = info.Features
	}

	// Convert licenses
	if info.Licenses != nil {
		// V0043CsvString is []string
		// Convert comma-separated string to slice
		licensesList := strings.Split(*info.Licenses, ",")
		csvLicenses := api.V0043CsvString(licensesList)
		descMsg.Licenses = &csvLicenses
	}

	// Convert flags
	if info.Flags != nil {
		flags := make([]api.V0043ReservationDescMsgFlags, len(*info.Flags))
		for i, flag := range *info.Flags {
			// Convert between flag types
			flags[i] = api.V0043ReservationDescMsgFlags(string(flag))
		}
		descMsg.Flags = &flags
	}

	// Convert burst buffer
	if info.BurstBuffer != nil {
		descMsg.BurstBuffer = info.BurstBuffer
	}

	// Note: Comment field might not exist in ReservationInfo
	// Comment would be set directly in the create request

	// Convert core count
	if info.CoreCount != nil {
		coreCount := int32(*info.CoreCount)
		setTrue := true
		descMsg.CoreCount = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &coreCount,
		}
	}

	// Convert groups
	if info.Groups != nil {
		// V0043CsvString is []string
		// Convert comma-separated string to slice
		groupsList := strings.Split(*info.Groups, ",")
		csvGroups := api.V0043CsvString(groupsList)
		descMsg.Groups = &csvGroups
	}

	// Convert max start delay
	if info.MaxStartDelay != nil {
		maxDelay := int32(*info.MaxStartDelay)
		setTrue := true
		descMsg.MaxStartDelay = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &maxDelay,
		}
	}

	// Convert TRES
	if info.Tres != nil {
		tresList := make(api.V0043TresList, 0)
		// Note: Need to parse TRES string format if it's a string
		// For now, just create empty list
		descMsg.Tres = &tresList
	}

	return descMsg, nil
}
