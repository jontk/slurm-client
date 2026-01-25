// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	"strings"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/pkg/errors"
)

// QoSManagerImpl implements the QoSManager interface for v0.0.43
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
	// Use common client initialization check
	if err := common.CheckClientInitialized(q.client); err != nil {
		return nil, err
	}
	if err := common.CheckClientInitialized(q.client.apiClient); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0043GetQosParams{}

	// Apply filters from options
	if opts != nil {
		if len(opts.Names) > 0 {
			nameStr := strings.Join(opts.Names, ",")
			params.Name = &nameStr
		}
	}

	// Call the generated OpenAPI client
	resp, err := q.client.apiClient.SlurmdbV0043GetQosWithResponse(ctx, params)
	if err != nil {
		return nil, common.WrapAndEnhanceError(err, "v0.0.43")
	}

	// Use common response error handling
	var apiErrors *V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := common.CheckNilResponse(resp.JSON200, "List QoS"); err != nil {
		return nil, err
	}
	if err := common.CheckNilResponse(resp.JSON200.Qos, "List QoS - qos field"); err != nil {
		return nil, err
	}

	// Convert the response to our interface types
	qosList := make([]interfaces.QoS, 0, len(resp.JSON200.Qos))
	for _, apiQos := range resp.JSON200.Qos {
		qos := convertAPIQoSToInterface(apiQos)
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
	// Validate input first (cheap check)
	if qosName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS name is required", "qosName", qosName, nil)
	}

	// Then check client initialization
	if err := common.CheckClientInitialized(q.client); err != nil {
		return nil, err
	}
	if err := common.CheckClientInitialized(q.client.apiClient); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0043GetSingleQosParams{}

	// Call the generated OpenAPI client
	resp, err := q.client.apiClient.SlurmdbV0043GetSingleQosWithResponse(ctx, qosName, params)
	if err != nil {
		return nil, common.WrapAndEnhanceError(err, "v0.0.43")
	}

	// Use common response error handling
	var apiErrors *V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := common.CheckNilResponse(resp.JSON200, "Get QoS"); err != nil {
		return nil, err
	}
	if err := common.CheckNilResponse(resp.JSON200.Qos, "Get QoS - qos field"); err != nil {
		return nil, err
	}

	// Check if we got any QoS entries
	if len(resp.JSON200.Qos) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "QoS not found", "QoS '"+qosName+"' not found")
	}

	// Convert the first QoS (should be the only one)
	qos := convertAPIQoSToInterface(resp.JSON200.Qos[0])
	return qos, nil
}

// Create creates a new QoS
func (q *QoSManagerImpl) Create(ctx context.Context, qos *interfaces.QoSCreate) (*interfaces.QoSCreateResponse, error) {
	// Use common client initialization check
	if err := common.CheckClientInitialized(q.client); err != nil {
		return nil, err
	}
	if err := common.CheckClientInitialized(q.client.apiClient); err != nil {
		return nil, err
	}

	if qos == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS creation data is required", "qos", qos, nil)
	}

	if qos.Name == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS name is required", "qos.Name", qos.Name, nil)
	}

	// Validate priority - must be non-negative
	if qos.Priority < 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS priority must be non-negative", "qos.Priority", qos.Priority, nil)
	}

	// Convert the QoS create request to API format
	apiQoS := convertQoSCreateToAPI(qos)

	// Create request body
	reqBody := SlurmdbV0043PostQosJSONRequestBody{
		Qos: V0043QosList{*apiQoS},
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0043PostQosParams{}

	// Call the generated OpenAPI client
	resp, err := q.client.apiClient.SlurmdbV0043PostQosWithResponse(ctx, params, reqBody)
	if err != nil {
		return nil, common.WrapAndEnhanceError(err, "v0.0.43")
	}

	// Use common response error handling
	var apiErrors *V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
		return nil, err
	}

	return &interfaces.QoSCreateResponse{
		QoSName: qos.Name,
	}, nil
}

// Update updates an existing QoS
func (q *QoSManagerImpl) Update(ctx context.Context, qosName string, update *interfaces.QoSUpdate) error {
	// Use common client initialization check
	if err := common.CheckClientInitialized(q.client); err != nil {
		return err
	}
	if err := common.CheckClientInitialized(q.client.apiClient); err != nil {
		return err
	}

	if qosName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS name is required", "qosName", qosName, nil)
	}

	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "update data is required", "update", update, nil)
	}

	// Validate update values
	if update.Priority != nil && *update.Priority < 0 {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS priority must be non-negative", "update.Priority", *update.Priority, nil)
	}

	// First, get the existing QoS to merge updates
	existingQoS, err := q.Get(ctx, qosName)
	if err != nil {
		return err
	}

	// Convert existing QoS to API format and apply updates
	apiQoS := convertQoSUpdateToAPI(existingQoS, update)

	// Create request body
	reqBody := SlurmdbV0043PostQosJSONRequestBody{
		Qos: V0043QosList{*apiQoS},
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0043PostQosParams{}

	// Call the generated OpenAPI client (POST is used for updates in SLURM API)
	resp, err := q.client.apiClient.SlurmdbV0043PostQosWithResponse(ctx, params, reqBody)
	if err != nil {
		return common.WrapAndEnhanceError(err, "v0.0.43")
	}

	// Use common response error handling
	var errors *V0043OpenapiErrors
	if resp.JSON200 != nil {
		errors = resp.JSON200.Errors
	}

	responseAdapter := NewResponseAdapter(resp.StatusCode(), errors)
	return common.HandleAPIResponse(responseAdapter, "v0.0.43")
}

// Delete deletes a QoS
func (q *QoSManagerImpl) Delete(ctx context.Context, qosName string) error {
	// Use common client initialization check
	if err := common.CheckClientInitialized(q.client); err != nil {
		return err
	}
	if err := common.CheckClientInitialized(q.client.apiClient); err != nil {
		return err
	}

	if qosName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS name is required", "qosName", qosName, nil)
	}

	// Call the generated OpenAPI client
	resp, err := q.client.apiClient.SlurmdbV0043DeleteSingleQosWithResponse(ctx, qosName)
	if err != nil {
		return common.WrapAndEnhanceError(err, "v0.0.43")
	}

	// Use common response error handling (200 or 204 for successful deletion)
	var errors *V0043OpenapiErrors
	if resp.JSON200 != nil {
		errors = resp.JSON200.Errors
	}

	// Create adapter with special handling for 204 (No Content) status
	responseAdapter := NewResponseAdapter(resp.StatusCode(), errors)

	// For DELETE operations, 204 is also a success
	if resp.StatusCode() == 204 {
		return nil
	}

	return common.HandleAPIResponse(responseAdapter, "v0.0.43")
}

// ===== HELPER FUNCTIONS (remain unchanged) =====

// The helper functions below remain the same as in the original implementation
// since they don't contain the repetitive error handling patterns

// convertAPIQoSToInterface converts a V0043Qos to interfaces.QoS
func convertAPIQoSToInterface(apiQoS V0043Qos) *interfaces.QoS {
	qos := &interfaces.QoS{}

	// Basic fields
	if apiQoS.Name != nil {
		qos.Name = *apiQoS.Name
	}
	if apiQoS.Description != nil {
		qos.Description = *apiQoS.Description
	}

	// Priority
	if apiQoS.Priority != nil && apiQoS.Priority.Set != nil && *apiQoS.Priority.Set && apiQoS.Priority.Number != nil {
		qos.Priority = int(*apiQoS.Priority.Number)
	}

	// Flags
	if apiQoS.Flags != nil {
		flags := make([]string, 0, len(*apiQoS.Flags))
		for _, flag := range *apiQoS.Flags {
			flags = append(flags, string(flag))
		}
		qos.Flags = flags
	}

	// Preempt mode
	if apiQoS.Preempt != nil && apiQoS.Preempt.Mode != nil && len(*apiQoS.Preempt.Mode) > 0 {
		qos.PreemptMode = string((*apiQoS.Preempt.Mode)[0])
	}

	// Grace time (from Limits)
	if apiQoS.Limits != nil && apiQoS.Limits.GraceTime != nil {
		qos.GraceTime = int(*apiQoS.Limits.GraceTime)
	}

	// Usage factor
	if apiQoS.UsageFactor != nil && apiQoS.UsageFactor.Set != nil && *apiQoS.UsageFactor.Set && apiQoS.UsageFactor.Number != nil {
		qos.UsageFactor = *apiQoS.UsageFactor.Number
	}

	// Usage threshold
	if apiQoS.UsageThreshold != nil && apiQoS.UsageThreshold.Set != nil && *apiQoS.UsageThreshold.Set && apiQoS.UsageThreshold.Number != nil {
		qos.UsageThreshold = *apiQoS.UsageThreshold.Number
	}

	// Max jobs per user (from Limits)
	if apiQoS.Limits != nil && apiQoS.Limits.Max != nil && apiQoS.Limits.Max.Jobs != nil && apiQoS.Limits.Max.Jobs.Per != nil && apiQoS.Limits.Max.Jobs.Per.User != nil {
		if apiQoS.Limits.Max.Jobs.Per.User.Set != nil && *apiQoS.Limits.Max.Jobs.Per.User.Set && apiQoS.Limits.Max.Jobs.Per.User.Number != nil {
			qos.MaxJobsPerUser = int(*apiQoS.Limits.Max.Jobs.Per.User.Number)
		}
	}

	// Max jobs per account (from Limits)
	if apiQoS.Limits != nil && apiQoS.Limits.Max != nil && apiQoS.Limits.Max.Jobs != nil && apiQoS.Limits.Max.Jobs.Per != nil && apiQoS.Limits.Max.Jobs.Per.Account != nil {
		if apiQoS.Limits.Max.Jobs.Per.Account.Set != nil && *apiQoS.Limits.Max.Jobs.Per.Account.Set && apiQoS.Limits.Max.Jobs.Per.Account.Number != nil {
			qos.MaxJobsPerAccount = int(*apiQoS.Limits.Max.Jobs.Per.Account.Number)
		}
	}

	return qos
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
func convertQoSCreateToAPI(create *interfaces.QoSCreate) *V0043Qos {
	apiQoS := &V0043Qos{}

	// Required fields
	apiQoS.Name = &create.Name

	// Optional fields
	if create.Description != "" {
		apiQoS.Description = &create.Description
	}

	// Priority
	if create.Priority > 0 {
		setTrue := true
		priority := int32(create.Priority)
		apiQoS.Priority = &V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &priority,
		}
	}

	// Flags
	if len(create.Flags) > 0 {
		flags := make([]V0043QosFlags, 0, len(create.Flags))
		for _, flag := range create.Flags {
			flags = append(flags, V0043QosFlags(flag))
		}
		apiQoS.Flags = &flags
	}

	// Preempt mode
	if len(create.PreemptMode) > 0 {
		modes := make([]V0043QosPreemptMode, 0, len(create.PreemptMode))
		for _, mode := range create.PreemptMode {
			modes = append(modes, V0043QosPreemptMode(mode))
		}
		if apiQoS.Preempt == nil {
			apiQoS.Preempt = &struct {
				ExemptTime *V0043Uint32NoValStruct `json:"exempt_time,omitempty"`
				List       *V0043QosPreemptList    `json:"list,omitempty"`
				Mode       *[]V0043QosPreemptMode  `json:"mode,omitempty"`
			}{}
		}
		apiQoS.Preempt.Mode = &modes
	}

	// Preempt list - not part of the interfaces, skipping

	// Grace time - simplified for now
	// TODO: Implement grace time when full Limits struct is available

	// Usage factor
	if create.UsageFactor != 0 {
		setTrue := true
		apiQoS.UsageFactor = &V0043Float64NoValStruct{
			Set:    &setTrue,
			Number: &create.UsageFactor,
		}
	}

	// Usage threshold
	if create.UsageThreshold != 0 {
		setTrue := true
		apiQoS.UsageThreshold = &V0043Float64NoValStruct{
			Set:    &setTrue,
			Number: &create.UsageThreshold,
		}
	}

	// Limits would be converted here if provided
	// This is a simplified version - full implementation would handle all limit types

	return apiQoS
}

// convertQoSUpdateToAPI converts interfaces.QoSUpdate to API format
func convertQoSUpdateToAPI(existing *interfaces.QoS, update *interfaces.QoSUpdate) *V0043Qos {
	apiQoS := &V0043Qos{}
	apiQoS.Name = &existing.Name
	apiQoS.Description = &existing.Description

	// Apply updates
	if update.Description != nil {
		apiQoS.Description = update.Description
	}

	// Priority
	priority := existing.Priority
	if update.Priority != nil {
		priority = *update.Priority
	}
	if priority > 0 {
		setTrue := true
		priorityInt32 := int32(priority)
		apiQoS.Priority = &V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &priorityInt32,
		}
	}

	// Flags
	flags := existing.Flags
	if len(update.Flags) > 0 {
		flags = update.Flags
	}
	if len(flags) > 0 {
		apiFlags := make([]V0043QosFlags, 0, len(flags))
		for _, flag := range flags {
			apiFlags = append(apiFlags, V0043QosFlags(flag))
		}
		apiQoS.Flags = &apiFlags
	}

	// Preempt mode
	preemptMode := existing.PreemptMode
	if update.PreemptMode != nil {
		preemptMode = *update.PreemptMode
	}
	if preemptMode != "" {
		modes := []V0043QosPreemptMode{V0043QosPreemptMode(preemptMode)}
		if apiQoS.Preempt == nil {
			apiQoS.Preempt = &struct {
				ExemptTime *V0043Uint32NoValStruct `json:"exempt_time,omitempty"`
				List       *V0043QosPreemptList    `json:"list,omitempty"`
				Mode       *[]V0043QosPreemptMode  `json:"mode,omitempty"`
			}{}
		}
		apiQoS.Preempt.Mode = &modes
	}

	// Preempt list - not part of the interfaces, skipping

	// Grace time - simplified for now
	// TODO: Implement grace time when full Limits struct is available

	// Usage factor
	usageFactor := existing.UsageFactor
	if update.UsageFactor != nil {
		usageFactor = *update.UsageFactor
	}
	if usageFactor != 0 {
		setTrue := true
		apiQoS.UsageFactor = &V0043Float64NoValStruct{
			Set:    &setTrue,
			Number: &usageFactor,
		}
	}

	// Usage threshold
	usageThreshold := existing.UsageThreshold
	if update.UsageThreshold != nil {
		usageThreshold = *update.UsageThreshold
	}
	if usageThreshold != 0 {
		setTrue := true
		apiQoS.UsageThreshold = &V0043Float64NoValStruct{
			Set:    &setTrue,
			Number: &usageThreshold,
		}
	}

	// Limits would be updated here if provided
	// This is a simplified version - full implementation would handle all limit types

	return apiQoS
}
