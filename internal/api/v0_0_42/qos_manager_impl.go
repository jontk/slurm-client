// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// QoSManagerImpl implements the QoSManager interface for v0.0.42
type QoSManagerImpl struct {
	client *WrapperClient
}

// NewQoSManagerImpl creates a new QoSManagerImpl
func NewQoSManagerImpl(client *WrapperClient) *QoSManagerImpl {
	return &QoSManagerImpl{
		client: client,
	}
}

// List retrieves a list of QoS with optional filtering
func (q *QoSManagerImpl) List(ctx context.Context, opts *interfaces.ListQoSOptions) (*interfaces.QoSList, error) {
	if q.client == nil || q.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0042GetQosParams{}

	// Apply filters from options
	if opts != nil {
		if len(opts.Names) > 0 {
			nameStr := strings.Join(opts.Names, ",")
			params.Name = &nameStr
		}
	}

	// Call the generated OpenAPI client
	resp, err := q.client.apiClient.SlurmdbV0042GetQosWithResponse(ctx, params)
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
	if resp.JSON200 == nil || resp.JSON200.Qos == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response with qos but got nil")
	}

	// Convert the response to our interface types
	qosList := make([]interfaces.QoS, 0, len(resp.JSON200.Qos))
	for _, apiQos := range resp.JSON200.Qos {
		qos, err := convertAPIQoSToInterface(apiQos)
		if err != nil {
			conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert QoS data")
			conversionErr.Cause = err
			conversionErr.Details = fmt.Sprintf("Error converting QoS %v", apiQos.Name)
			return nil, conversionErr
		}
		qosList = append(qosList, *qos)
	}

	// Apply client-side filtering if options are provided
	if opts != nil {
		qosList = filterQoS(qosList, opts)
	}

	return &interfaces.QoSList{
		QoS:   qosList,
		Total: len(qosList),
	}, nil
}

// Get retrieves a specific QoS by name
func (q *QoSManagerImpl) Get(ctx context.Context, qosName string) (*interfaces.QoS, error) {
	if q.client == nil || q.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if qosName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS name is required", "qosName", qosName, nil)
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0042GetSingleQosParams{}

	// Call the generated OpenAPI client
	resp, err := q.client.apiClient.SlurmdbV0042GetSingleQosWithResponse(ctx, qosName, params)
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
	if resp.JSON200 == nil || resp.JSON200.Qos == nil || len(resp.JSON200.Qos) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "QoS not found", fmt.Sprintf("QoS '%s' not found", qosName))
	}

	// Convert the first QoS in the response
	qos, err := convertAPIQoSToInterface(resp.JSON200.Qos[0])
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert QoS data")
		conversionErr.Cause = err
		conversionErr.Details = fmt.Sprintf("Error converting QoS '%s'", qosName)
		return nil, conversionErr
	}

	return qos, nil
}

// Create creates a new QoS
func (q *QoSManagerImpl) Create(ctx context.Context, qos *interfaces.QoSCreate) (*interfaces.QoSCreateResponse, error) {
	if q.client == nil || q.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if qos == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS data is required", "qos", qos, nil)
	}

	// Validate required fields
	if qos.Name == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS name is required", "qos.Name", qos.Name, nil)
	}

	// Convert interface types to API types
	apiQos, err := convertQoSCreateToAPI(qos)
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeInvalidRequest, "Failed to convert QoS data")
		conversionErr.Cause = err
		return nil, conversionErr
	}

	// Create the request body
	requestBody := SlurmdbV0042PostQosJSONRequestBody{
		Qos: []V0042Qos{*apiQos},
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0042PostQosParams{}

	// Call the generated OpenAPI client
	resp, err := q.client.apiClient.SlurmdbV0042PostQosWithResponse(ctx, params, requestBody)
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
	response := &interfaces.QoSCreateResponse{
		QoSName: qos.Name,
	}

	return response, nil
}

// Update updates an existing QoS
func (q *QoSManagerImpl) Update(ctx context.Context, qosName string, update *interfaces.QoSUpdate) error {
	if q.client == nil || q.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if qosName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS name is required", "qosName", qosName, nil)
	}

	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "update data is required", "update", update, nil)
	}

	// Convert update to API format
	apiUpdate, err := convertQoSUpdateToAPI(update)
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeInvalidRequest, "Failed to convert update data")
		conversionErr.Cause = err
		return conversionErr
	}

	// Set the QoS name in the update
	apiUpdate.Name = &qosName

	// Create the request body
	requestBody := SlurmdbV0042PostQosJSONRequestBody{
		Qos: []V0042Qos{*apiUpdate},
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0042PostQosParams{}

	// Call the generated OpenAPI client (POST is used for updates in Slurm)
	resp, err := q.client.apiClient.SlurmdbV0042PostQosWithResponse(ctx, params, requestBody)
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

// Delete deletes a QoS
func (q *QoSManagerImpl) Delete(ctx context.Context, qosName string) error {
	if q.client == nil || q.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if qosName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS name is required", "qosName", qosName, nil)
	}

	// Call the generated OpenAPI client
	resp, err := q.client.apiClient.SlurmdbV0042DeleteSingleQosWithResponse(ctx, qosName)
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

// convertAPIQoSToInterface converts V0042Qos to interfaces.QoS
func convertAPIQoSToInterface(apiQos V0042Qos) (*interfaces.QoS, error) {
	qos := &interfaces.QoS{}

	// Basic fields
	if apiQos.Name != nil {
		qos.Name = *apiQos.Name
	}

	if apiQos.Description != nil {
		qos.Description = *apiQos.Description
	}

	// Priority
	if apiQos.Priority != nil && apiQos.Priority.Set != nil && *apiQos.Priority.Set && apiQos.Priority.Number != nil {
		qos.Priority = int(*apiQos.Priority.Number)
	}

	// Usage factor and threshold
	if apiQos.UsageFactor != nil && apiQos.UsageFactor.Set != nil && *apiQos.UsageFactor.Set && apiQos.UsageFactor.Number != nil {
		qos.UsageFactor = *apiQos.UsageFactor.Number
	}

	if apiQos.UsageThreshold != nil && apiQos.UsageThreshold.Set != nil && *apiQos.UsageThreshold.Set && apiQos.UsageThreshold.Number != nil {
		qos.UsageThreshold = *apiQos.UsageThreshold.Number
	}

	// Limits
	if apiQos.Limits != nil {
		// Grace time
		if apiQos.Limits.GraceTime != nil {
			qos.GraceTime = int(*apiQos.Limits.GraceTime)
		}

		// Max limits
		if apiQos.Limits.Max != nil {
			if apiQos.Limits.Max.Jobs != nil {
				if apiQos.Limits.Max.Jobs.Count != nil && apiQos.Limits.Max.Jobs.Count.Set != nil && *apiQos.Limits.Max.Jobs.Count.Set && apiQos.Limits.Max.Jobs.Count.Number != nil {
					qos.MaxJobs = int(*apiQos.Limits.Max.Jobs.Count.Number)
				}
				if apiQos.Limits.Max.Jobs.Per != nil {
					if apiQos.Limits.Max.Jobs.Per.User != nil && apiQos.Limits.Max.Jobs.Per.User.Set != nil && *apiQos.Limits.Max.Jobs.Per.User.Set && apiQos.Limits.Max.Jobs.Per.User.Number != nil {
						qos.MaxJobsPerUser = int(*apiQos.Limits.Max.Jobs.Per.User.Number)
					}
					if apiQos.Limits.Max.Jobs.Per.Account != nil && apiQos.Limits.Max.Jobs.Per.Account.Set != nil && *apiQos.Limits.Max.Jobs.Per.Account.Set && apiQos.Limits.Max.Jobs.Per.Account.Number != nil {
						qos.MaxJobsPerAccount = int(*apiQos.Limits.Max.Jobs.Per.Account.Number)
					}
				}
			}

			// Submit jobs limit (found in active jobs)
			if apiQos.Limits.Max.ActiveJobs != nil && apiQos.Limits.Max.ActiveJobs.Count != nil && apiQos.Limits.Max.ActiveJobs.Count.Set != nil && *apiQos.Limits.Max.ActiveJobs.Count.Set && apiQos.Limits.Max.ActiveJobs.Count.Number != nil {
				qos.MaxSubmitJobs = int(*apiQos.Limits.Max.ActiveJobs.Count.Number)
			}

			// TRES limits for CPUs and nodes
			if apiQos.Limits.Max.Tres != nil && apiQos.Limits.Max.Tres.Per != nil {
				if apiQos.Limits.Max.Tres.Per.Job != nil {
					for _, tres := range *apiQos.Limits.Max.Tres.Per.Job {
						if tres.Count != nil {
							switch tres.Type {
							case "cpu":
								qos.MaxCPUs = int(*tres.Count)
							case "node":
								qos.MaxNodes = int(*tres.Count)
							}
						}
					}
				}
				if apiQos.Limits.Max.Tres.Per.User != nil {
					for _, tres := range *apiQos.Limits.Max.Tres.Per.User {
						if tres.Count != nil {
							if tres.Type == "cpu" {
								qos.MaxCPUsPerUser = int(*tres.Count)
							}
						}
					}
				}
			}

			// Wall clock limit
			if apiQos.Limits.Max.WallClock != nil && apiQos.Limits.Max.WallClock.Per != nil && apiQos.Limits.Max.WallClock.Per.Job != nil && apiQos.Limits.Max.WallClock.Per.Job.Set != nil && *apiQos.Limits.Max.WallClock.Per.Job.Set && apiQos.Limits.Max.WallClock.Per.Job.Number != nil {
				qos.MaxWallTime = int(*apiQos.Limits.Max.WallClock.Per.Job.Number)
			}
		}

		// Min limits
		if apiQos.Limits.Min != nil && apiQos.Limits.Min.Tres != nil && apiQos.Limits.Min.Tres.Per != nil && apiQos.Limits.Min.Tres.Per.Job != nil {
			for _, tres := range *apiQos.Limits.Min.Tres.Per.Job {
				if tres.Count != nil {
					switch tres.Type {
					case "cpu":
						qos.MinCPUs = int(*tres.Count)
					case "node":
						qos.MinNodes = int(*tres.Count)
					}
				}
			}
		}
	}

	// Preempt information
	if apiQos.Preempt != nil {
		if apiQos.Preempt.Mode != nil && len(*apiQos.Preempt.Mode) > 0 {
			qos.PreemptMode = (*apiQos.Preempt.Mode)[0]
		}
	}

	// Flags
	if apiQos.Flags != nil {
		qos.Flags = *apiQos.Flags
	}

	return qos, nil
}

// filterQoS applies client-side filtering to the QoS list
func filterQoS(qosList []interfaces.QoS, opts *interfaces.ListQoSOptions) []interfaces.QoS {
	if opts == nil {
		return qosList
	}

	filtered := make([]interfaces.QoS, 0, len(qosList))
	for _, qos := range qosList {
		// Filter by names
		if len(opts.Names) > 0 {
			found := false
			for _, name := range opts.Names {
				if qos.Name == name {
					found = true
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
				for _, qosAccount := range qos.AllowedAccounts {
					if qosAccount == filterAccount {
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

		// Filter by users
		if len(opts.Users) > 0 {
			found := false
			for _, filterUser := range opts.Users {
				for _, qosUser := range qos.AllowedUsers {
					if qosUser == filterUser {
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

		filtered = append(filtered, qos)
	}

	// Apply limit and offset
	if opts.Limit > 0 {
		start := opts.Offset
		if start < 0 {
			start = 0
		}
		if start >= len(filtered) {
			return []interfaces.QoS{}
		}
		end := start + opts.Limit
		if end > len(filtered) {
			end = len(filtered)
		}
		filtered = filtered[start:end]
	}

	return filtered
}

// convertQoSCreateToAPI converts interfaces.QoSCreate to API format
func convertQoSCreateToAPI(create *interfaces.QoSCreate) (*V0042Qos, error) {
	apiQos := &V0042Qos{}

	// Required fields
	apiQos.Name = &create.Name

	// Optional fields
	if create.Description != "" {
		apiQos.Description = &create.Description
	}

	if create.Priority > 0 {
		priority := int32(create.Priority)
		apiQos.Priority = &V0042Uint32NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &priority,
		}
	}

	if create.UsageFactor > 0 {
		apiQos.UsageFactor = &V0042Float64NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &create.UsageFactor,
		}
	}

	if create.UsageThreshold > 0 {
		apiQos.UsageThreshold = &V0042Float64NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &create.UsageThreshold,
		}
	}

	// Preempt mode
	if create.PreemptMode != "" {
		apiQos.Preempt = &struct {
			ExemptTime *V0042Uint32NoValStruct `json:"exempt_time,omitempty"`
			List       *V0042QosPreemptList    `json:"list,omitempty"`
			Mode       *V0042QosPreemptModes   `json:"mode,omitempty"`
		}{
			Mode: &V0042QosPreemptModes{create.PreemptMode},
		}
	}

	// Flags
	if len(create.Flags) > 0 {
		flags := V0042QosFlags(create.Flags)
		apiQos.Flags = &flags
	}

	// For now, we'll implement only basic functionality
	// The complex nested limits structure can be added later if needed
	// This provides a working QoS manager for basic operations

	return apiQos, nil
}

// convertQoSUpdateToAPI converts interfaces.QoSUpdate to API format
func convertQoSUpdateToAPI(update *interfaces.QoSUpdate) (*V0042Qos, error) {
	apiQos := &V0042Qos{}

	// Optional fields (only if specified in update)
	if update.Description != nil {
		apiQos.Description = update.Description
	}

	if update.Priority != nil {
		priority := int32(*update.Priority)
		apiQos.Priority = &V0042Uint32NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &priority,
		}
	}

	if update.UsageFactor != nil {
		apiQos.UsageFactor = &V0042Float64NoValStruct{
			Set:    &[]bool{true}[0],
			Number: update.UsageFactor,
		}
	}

	if update.UsageThreshold != nil {
		apiQos.UsageThreshold = &V0042Float64NoValStruct{
			Set:    &[]bool{true}[0],
			Number: update.UsageThreshold,
		}
	}

	// Preempt mode
	if update.PreemptMode != nil {
		apiQos.Preempt = &struct {
			ExemptTime *V0042Uint32NoValStruct `json:"exempt_time,omitempty"`
			List       *V0042QosPreemptList    `json:"list,omitempty"`
			Mode       *V0042QosPreemptModes   `json:"mode,omitempty"`
		}{
			Mode: &V0042QosPreemptModes{*update.PreemptMode},
		}
	}

	// Flags
	if update.Flags != nil {
		flags := V0042QosFlags(update.Flags)
		apiQos.Flags = &flags
	}

	// For now, we'll implement only basic functionality
	// The complex nested limits structure can be added later if needed
	// This provides a working QoS manager for basic operations

	return apiQos, nil
}
