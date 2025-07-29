package v0_0_41

import (
	"context"
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
)

// QoSAdapter implements the QoSAdapter interface for v0.0.41
type QoSAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewQoSAdapter creates a new QoS adapter for v0.0.41
func NewQoSAdapter(client *api.ClientWithResponses) *QoSAdapter {
	return &QoSAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "QoS"),
		client:      client,
		wrapper:     nil, // We'll implement this later
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
		if opts.ID > 0 {
			idStr := fmt.Sprintf("%d", opts.ID)
			params.Id = &idStr
		}
		if opts.Description != "" {
			params.Description = &opts.Description
		}
		if opts.WithDeleted {
			withDeleted := "true"
			params.WithDeleted = &withDeleted
		}
		if opts.PreemptMode != "" {
			preemptMode := convertPreemptModeToAPI(opts.PreemptMode)
			params.PreemptMode = &preemptMode
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
		QoS: make([]types.QoS, 0, len(resp.JSON200.Qos)),
		Meta: &types.ListMeta{
			Version: a.GetVersion(),
		},
	}

	for _, apiQoS := range resp.JSON200.Qos {
		qos, err := a.convertAPIQoSToCommon(apiQoS)
		if err != nil {
			// Log the error but continue processing other QoS
			continue
		}
		qosList.QoS = append(qosList.QoS, *qos)
	}

	// Extract warning messages if any
	if resp.JSON200.Warnings != nil {
		warnings := make([]string, 0, len(*resp.JSON200.Warnings))
		for _, warning := range *resp.JSON200.Warnings {
			if warning.Description != nil {
				warnings = append(warnings, *warning.Description)
			}
		}
		if len(warnings) > 0 {
			qosList.Meta.Warnings = warnings
		}
	}

	// Extract error messages if any
	if resp.JSON200.Errors != nil {
		errors := make([]string, 0, len(*resp.JSON200.Errors))
		for _, error := range *resp.JSON200.Errors {
			if error.Description != nil {
				errors = append(errors, *error.Description)
			}
		}
		if len(errors) > 0 {
			qosList.Meta.Errors = errors
		}
	}

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
		return nil, a.WrapError(err, fmt.Sprintf("failed to get QoS %s", name))
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}

	if resp.JSON200 == nil || len(resp.JSON200.Qos) == 0 {
		return nil, a.HandleNotFound(fmt.Sprintf("QoS %s", name))
	}

	// Convert the first QoS in the response
	qos, err := a.convertAPIQoSToCommon(resp.JSON200.Qos[0])
	if err != nil {
		return nil, a.WrapError(err, fmt.Sprintf("failed to convert QoS %s", name))
	}

	return qos, nil
}

// Create creates a new QoS
func (a *QoSAdapter) Create(ctx context.Context, qos *types.QoS) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Validate QoS
	if qos == nil {
		return a.HandleValidationError("QoS cannot be nil")
	}
	if err := a.ValidateResourceName("QoS name", qos.Name); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Convert QoS to API request
	createReq := a.convertCommonToAPIQoS(qos)

	// Make the API call
	params := &api.SlurmdbV0041PostQosParams{}
	resp, err := a.client.SlurmdbV0041PostQosWithResponse(ctx, params, *createReq)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to create QoS %s", qos.Name))
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
	}

	return nil
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
		existingQoS.Description = *update.Description
	}
	if update.Priority != nil {
		existingQoS.Priority = *update.Priority
	}
	if update.PreemptMode != nil {
		existingQoS.PreemptMode = *update.PreemptMode
	}
	if update.GraceTime != nil {
		existingQoS.GraceTime = *update.GraceTime
	}
	if update.MaxWall != nil {
		existingQoS.MaxWall = *update.MaxWall
	}

	// Convert to API request
	updateReq := a.convertCommonToAPIQoS(existingQoS)

	// Make the API call
	params := &api.SlurmdbV0041PostQosParams{}
	resp, err := a.client.SlurmdbV0041PostQosWithResponse(ctx, params, *updateReq)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to update QoS %s", name))
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
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
		return a.WrapError(err, fmt.Sprintf("failed to delete QoS %s", name))
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

	if limits.MaxTRES != nil {
		// Convert TRES map to string format
		tresStr := formatTRESMap(limits.MaxTRES)
		update.MaxTRES = &tresStr
	}

	if limits.MaxTRESPerUser != nil {
		tresStr := formatTRESMap(limits.MaxTRESPerUser)
		update.MaxTRESPerUser = &tresStr
	}

	if limits.MaxTRESPerJob != nil {
		tresStr := formatTRESMap(limits.MaxTRESPerJob)
		update.MaxTRESPerJob = &tresStr
	}

	if limits.MaxJobsPerUser != nil {
		update.MaxJobsPerUser = limits.MaxJobsPerUser
	}

	if limits.MaxSubmitJobsPerUser != nil {
		update.MaxSubmitJobsPerUser = limits.MaxSubmitJobsPerUser
	}

	if limits.MaxWall != nil {
		update.MaxWall = limits.MaxWall
	}

	return a.Update(ctx, name, update)
}

// formatTRESMap converts a TRES map to string format
func formatTRESMap(tresMap map[string]uint64) string {
	var parts []string
	for key, value := range tresMap {
		parts = append(parts, fmt.Sprintf("%s=%d", key, value))
	}
	return strings.Join(parts, ",")
}

// convertPreemptModeToAPI converts common preempt mode to API preempt mode
func convertPreemptModeToAPI(mode string) api.SlurmdbV0041GetQosParamsPreemptMode {
	switch strings.ToUpper(mode) {
	case "OFF":
		return api.SlurmdbV0041GetQosParamsPreemptModeOFF
	case "CANCEL":
		return api.SlurmdbV0041GetQosParamsPreemptModeCANCEL
	case "REQUEUE":
		return api.SlurmdbV0041GetQosParamsPreemptModeREQUEUE
	case "SUSPEND":
		return api.SlurmdbV0041GetQosParamsPreemptModeSUSPEND
	default:
		return api.SlurmdbV0041GetQosParamsPreemptModeOFF
	}
}