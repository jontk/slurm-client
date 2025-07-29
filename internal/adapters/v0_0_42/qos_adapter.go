package v0_0_42

import (
	"context"
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

// QoSAdapter implements the QoSAdapter interface for v0.0.42
type QoSAdapter struct {
	*base.BaseManager
	client *api.ClientWithResponses
}

// NewQoSAdapter creates a new QoS adapter for v0.0.42
func NewQoSAdapter(client *api.ClientWithResponses) *QoSAdapter {
	return &QoSAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "QoS"),
		client:      client,
	}
}

// List retrieves a list of QoS
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
	params := &api.SlurmdbV0042GetQosParams{}

	// Apply filters from options
	if opts != nil {
		if len(opts.Names) > 0 {
			qosStr := strings.Join(opts.Names, ",")
			params.Name = &qosStr
		}
		if opts.WithDeleted {
			withDeleted := "true"
			params.WithDeleted = &withDeleted
		}
		if opts.ID != "" {
			params.Id = &opts.ID
		}
		if opts.PreemptMode != "" {
			params.PreemptMode = &opts.PreemptMode
		}
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042GetQosWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list QoS")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(resp.StatusCode(), resp.Body)
	}

	// Check for API response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from API")
	}

	// Convert the response to common types
	qosList := &types.QoSList{
		QoS: make([]*types.QoS, 0),
	}

	if resp.JSON200.Qos != nil {
		for _, apiQoS := range *resp.JSON200.Qos {
			qos, err := a.convertAPIQoSToCommon(apiQoS)
			if err != nil {
				// Log conversion error but continue
				continue
			}
			qosList.QoS = append(qosList.QoS, qos)
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

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// v0.0.42 doesn't have a single QoS get endpoint, use list with filter
	params := &api.SlurmdbV0042GetQosParams{
		Name: &name,
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042GetQosWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, fmt.Sprintf("failed to get QoS %s", name))
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(resp.StatusCode(), resp.Body)
	}

	// Check for API response
	if resp.JSON200 == nil || resp.JSON200.Qos == nil || len(*resp.JSON200.Qos) == 0 {
		return nil, fmt.Errorf("QoS %s not found", name)
	}

	// Convert the first QoS in the response
	qosList := *resp.JSON200.Qos
	for _, apiQoS := range qosList {
		if apiQoS.Name != nil && *apiQoS.Name == name {
			return a.convertAPIQoSToCommon(apiQoS)
		}
	}

	return nil, fmt.Errorf("QoS %s not found", name)
}

// Create creates a new QoS
func (a *QoSAdapter) Create(ctx context.Context, qos *types.QoSCreateRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Convert common QoS to API format
	apiQoS, err := a.convertCommonQoSCreateToAPI(qos)
	if err != nil {
		return a.WrapError(err, "failed to convert QoS create request")
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042PostQosWithResponse(ctx, apiQoS)
	if err != nil {
		return a.WrapError(err, "failed to create QoS")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleAPIError(resp.StatusCode(), resp.Body)
	}

	return nil
}

// Update updates an existing QoS
func (a *QoSAdapter) Update(ctx context.Context, name string, updates *types.QoSUpdateRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// v0.0.42 doesn't have a direct QoS update endpoint
	// Updates require delete and recreate
	return fmt.Errorf("QoS update not directly supported via v0.0.42 API - use delete and recreate")
}

// Delete deletes a QoS
func (a *QoSAdapter) Delete(ctx context.Context, name string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042DeleteQosWithResponse(ctx, name)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to delete QoS %s", name))
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleAPIError(resp.StatusCode(), resp.Body)
	}

	return nil
}