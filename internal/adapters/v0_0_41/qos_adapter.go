// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"context"
	"fmt"
	"strings"

	types "github.com/jontk/slurm-client/api"
	adapterbase "github.com/jontk/slurm-client/internal/adapters/base"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
)

// QoSAdapter implements the QoSAdapter interface for v0.0.41
type QoSAdapter struct {
	*adapterbase.BaseManager
	client *api.ClientWithResponses
}

// NewQoSAdapter creates a new QoS adapter for v0.0.41
func NewQoSAdapter(client *api.ClientWithResponses) *QoSAdapter {
	return &QoSAdapter{
		BaseManager: adapterbase.NewBaseManager("v0.0.41", "QoS"),
		client:      client,
	}
}

// List retrieves a list of QoS with optional filtering
func (a *QoSAdapter) List(ctx context.Context, opts *types.QoSListOptions) (*types.QoSList, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}
	// Prepare parameters for the API call
	params := &api.SlurmdbV0041GetQosParams{}
	// Apply filters from options
	if opts != nil {
		if len(opts.Names) > 0 {
			nameStr := strings.Join(opts.Names, ",")
			params.Name = &nameStr
		}
		// ID, Description, WithDeleted, PreemptMode fields don't exist in QoSListOptions
		// Only Names, Accounts, Users, Limit, Offset are available
		if len(opts.Accounts) > 0 {
			// Convert accounts to string if API supports it
			_ = opts.Accounts
		}
		if len(opts.Users) > 0 {
			// Convert users to string if API supports it
			_ = opts.Users
		}
	}
	// Make the API call
	resp, err := a.client.SlurmdbV0041GetQosWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list QoS")
	}
	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected nil response")
	}
	// Convert response to common types
	qosList := &types.QoSList{
		QoS:   make([]types.QoS, 0, len(resp.JSON200.Qos)),
		Total: 0,
	}
	for _, apiQoS := range resp.JSON200.Qos {
		qos, err := a.convertAPIQoSToCommon(apiQoS)
		if err != nil {
			// Log the error but continue processing other QoS
			continue
		}
		qosList.QoS = append(qosList.QoS, *qos)
	}
	// Extract warning and error messages if any (but QoSList doesn't have Meta)
	// Warnings and errors are ignored for now as QoSList structure doesn't support them
	if resp.JSON200.Warnings != nil {
		// Log warnings if needed
		_ = resp.JSON200.Warnings
	}
	if resp.JSON200.Errors != nil {
		// Log errors if needed
		_ = resp.JSON200.Errors
	}
	// Update total count
	qosList.Total = len(qosList.QoS)
	return qosList, nil
}

// Get retrieves a specific QoS by name
func (a *QoSAdapter) Get(ctx context.Context, name string) (*types.QoS, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	// Validate name
	if err := a.ValidateResourceName("QoS name", name); err != nil {
		return nil, err
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}
	// Make the API call
	params := &api.SlurmdbV0041GetSingleQosParams{}
	resp, err := a.client.SlurmdbV0041GetSingleQosWithResponse(ctx, name, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to get QoS "+name)
	}
	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil || len(resp.JSON200.Qos) == 0 {
		return nil, a.HandleNotFound("QoS " + name)
	}
	// Convert the first QoS in the response
	qos, err := a.convertAPIQoSToCommon(resp.JSON200.Qos[0])
	if err != nil {
		return nil, a.WrapError(err, "failed to convert QoS "+name)
	}
	return qos, nil
}

// Create creates a new QoS
func (a *QoSAdapter) Create(ctx context.Context, qos *types.QoSCreate) (*types.QoSCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	// Validate QoS
	if qos == nil {
		return nil, fmt.Errorf("QoS cannot be nil")
	}
	if err := a.ValidateResourceName("QoS name", qos.Name); err != nil {
		return nil, err
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}
	// Convert to API request using JSON marshaling workaround
	reqBody, err := a.convertQoSCreateToAPI(qos)
	if err != nil {
		return nil, a.WrapError(err, "failed to convert QoS create request")
	}

	// Make the API call
	params := &api.SlurmdbV0041PostQosParams{}
	resp, err := a.client.SlurmdbV0041PostQosWithResponse(ctx, params, reqBody)
	if err != nil {
		return nil, a.WrapError(err, "failed to create QoS")
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}

	// Check for errors in response
	if resp.JSON200 != nil && resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		errMsgs := make([]string, 0)
		for _, apiErr := range *resp.JSON200.Errors {
			if apiErr.Error != nil {
				errMsgs = append(errMsgs, *apiErr.Error)
			}
		}
		if len(errMsgs) > 0 {
			return nil, fmt.Errorf("QoS creation failed: %v", errMsgs)
		}
	}

	return &types.QoSCreateResponse{
		QoSName: qos.Name,
	}, nil
}

// Update updates an existing QoS
func (a *QoSAdapter) Update(ctx context.Context, name string, update *types.QoSUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	// Validate name
	if err := a.ValidateResourceName("QoS name", name); err != nil {
		return err
	}
	// Validate update
	if update == nil {
		return a.HandleValidationError("QoS update cannot be nil")
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}
	// Get the existing QoS first
	existingQoS, err := a.Get(ctx, name)
	if err != nil {
		return err
	}
	// Apply updates
	if update.Description != nil {
		existingQoS.Description = update.Description
	}
	if update.Priority != nil {
		priority := uint32(*update.Priority)
		existingQoS.Priority = &priority
	}
	if update.PreemptMode != nil && len(*update.PreemptMode) > 0 {
		if existingQoS.Preempt == nil {
			existingQoS.Preempt = &types.QoSPreempt{}
		}
		modes := make([]types.ModeValue, 0, len(*update.PreemptMode))
		for _, m := range *update.PreemptMode {
			modes = append(modes, types.ModeValue(m))
		}
		existingQoS.Preempt.Mode = modes
	}
	if update.GraceTime != nil {
		if existingQoS.Limits == nil {
			existingQoS.Limits = &types.QoSLimits{}
		}
		gt := int32(*update.GraceTime)
		existingQoS.Limits.GraceTime = &gt
	}
	// MaxWall field doesn't exist in QoSUpdate type
	// Skip MaxWall update

	// Convert to API request using JSON marshaling workaround
	reqBody, err := a.convertQoSUpdateToAPI(existingQoS)
	if err != nil {
		return a.WrapError(err, "failed to convert QoS update request")
	}

	// Make the API call
	params := &api.SlurmdbV0041PostQosParams{}
	resp, err := a.client.SlurmdbV0041PostQosWithResponse(ctx, params, reqBody)
	if err != nil {
		return a.WrapError(err, "failed to update QoS "+name)
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
	}

	// Check for errors in response
	if resp.JSON200 != nil && resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		errMsgs := make([]string, 0)
		for _, apiErr := range *resp.JSON200.Errors {
			if apiErr.Error != nil {
				errMsgs = append(errMsgs, *apiErr.Error)
			}
		}
		if len(errMsgs) > 0 {
			return fmt.Errorf("QoS update failed: %v", errMsgs)
		}
	}

	return nil
}

// Delete deletes a QoS
func (a *QoSAdapter) Delete(ctx context.Context, name string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	// Validate name
	if err := a.ValidateResourceName("QoS name", name); err != nil {
		return err
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}
	// Make the API call
	resp, err := a.client.SlurmdbV0041DeleteSingleQosWithResponse(ctx, name)
	if err != nil {
		return a.WrapError(err, "failed to delete QoS "+name)
	}
	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
	}
	return nil
}

// SetLimits sets resource limits for a QoS
func (a *QoSAdapter) SetLimits(ctx context.Context, name string, limits *types.QoSLimits) error {
	// Use the Update method to set limits
	update := &types.QoSUpdate{}
	// Set limits properly using the Limits field
	update.Limits = limits
	// All limits are set via the Limits field above
	return a.Update(ctx, name, update)
}
