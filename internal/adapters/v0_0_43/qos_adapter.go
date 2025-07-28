package v0_0_43

import (
	"context"
	"strings"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// QoSAdapter implements the QoSAdapter interface for v0.0.43
type QoSAdapter struct {
	*base.QoSBaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewQoSAdapter creates a new QoS adapter for v0.0.43
func NewQoSAdapter(client *api.ClientWithResponses) *QoSAdapter {
	// For now, we'll use the client directly without the wrapper
	// The wrapper requires additional configuration that we'll add later
	return &QoSAdapter{
		QoSBaseManager: base.NewQoSBaseManager("v0.0.43"),
		client:         client,
		wrapper:        nil, // We'll implement this later
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
	params := &api.SlurmdbV0043GetQosParams{}

	// Apply filters from options
	if opts != nil && len(opts.Names) > 0 {
		nameStr := strings.Join(opts.Names, ",")
		params.Name = &nameStr
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0043GetQosWithResponse(ctx, params)
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
	if err := a.CheckNilResponse(resp.JSON200, "List QoS"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Qos, "List QoS - qos field"); err != nil {
		return nil, err
	}

	// Convert the response to common types
	qosList := make([]types.QoS, 0, len(resp.JSON200.Qos))
	for _, apiQos := range resp.JSON200.Qos {
		qos, err := a.convertAPIQoSToCommon(apiQos)
		if err != nil {
			return nil, a.HandleConversionError(err, apiQos.Name)
		}
		qosList = append(qosList, *qos)
	}

	// Apply client-side filtering using base manager
	if opts != nil {
		qosList = a.FilterQoSList(qosList, opts)
	}

	// Apply pagination using base manager
	listOpts := base.ListOptions{
		Limit:  opts.Limit,
		Offset: opts.Offset,
	}
	if opts != nil {
		listOpts.Limit = opts.Limit
		listOpts.Offset = opts.Offset
	}

	// Apply pagination
	start := listOpts.Offset
	if start < 0 {
		start = 0
	}
	if start >= len(qosList) {
		return &types.QoSList{
			QoS:   []types.QoS{},
			Total: len(qosList),
		}, nil
	}

	end := len(qosList)
	if listOpts.Limit > 0 {
		end = start + listOpts.Limit
		if end > len(qosList) {
			end = len(qosList)
		}
	}

	return &types.QoSList{
		QoS:   qosList[start:end],
		Total: len(qosList),
	}, nil
}

// Get retrieves a specific QoS by name
func (a *QoSAdapter) Get(ctx context.Context, qosName string) (*types.QoS, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.ValidateResourceName(qosName, "qosName"); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmdbV0043GetSingleQosParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0043GetSingleQosWithResponse(ctx, qosName, params)
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
	if err := a.CheckNilResponse(resp.JSON200, "Get QoS"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Qos, "Get QoS - qos field"); err != nil {
		return nil, err
	}

	// Check if we got any QoS entries
	if len(resp.JSON200.Qos) == 0 {
		return nil, a.ResourceNotFoundError(qosName)
	}

	// Convert the first QoS (should be the only one)
	qos, err := a.convertAPIQoSToCommon(resp.JSON200.Qos[0])
	if err != nil {
		return nil, a.HandleConversionError(err, qosName)
	}

	return qos, nil
}

// Create creates a new QoS
func (a *QoSAdapter) Create(ctx context.Context, qos *types.QoSCreate) (*types.QoSCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.ValidateQoSCreate(qos); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Apply defaults using base manager
	qos = a.ApplyQoSDefaults(qos)

	// Convert to API format
	apiQoS, err := a.convertCommonQoSCreateToAPI(qos)
	if err != nil {
		return nil, err
	}

	// Create request body
	reqBody := api.SlurmdbV0043PostQosJSONRequestBody{
		Qos: api.V0043QosList{*apiQoS},
	}

	// Prepare parameters for the API call
	params := &api.SlurmdbV0043PostQosParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0043PostQosWithResponse(ctx, params, reqBody)
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

	return &types.QoSCreateResponse{
		QoSName: qos.Name,
	}, nil
}

// Update updates an existing QoS
func (a *QoSAdapter) Update(ctx context.Context, qosName string, update *types.QoSUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(qosName, "qosName"); err != nil {
		return err
	}
	if err := a.ValidateQoSUpdate(update); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// First, get the existing QoS to merge updates
	existingQoS, err := a.Get(ctx, qosName)
	if err != nil {
		return err
	}

	// Convert to API format and apply updates
	apiQoS, err := a.convertCommonQoSUpdateToAPI(existingQoS, update)
	if err != nil {
		return err
	}

	// Create request body
	reqBody := api.SlurmdbV0043PostQosJSONRequestBody{
		Qos: api.V0043QosList{*apiQoS},
	}

	// Prepare parameters for the API call
	params := &api.SlurmdbV0043PostQosParams{}

	// Call the generated OpenAPI client (POST is used for updates in SLURM API)
	resp, err := a.client.SlurmdbV0043PostQosWithResponse(ctx, params, reqBody)
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

// Delete deletes a QoS
func (a *QoSAdapter) Delete(ctx context.Context, qosName string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(qosName, "qosName"); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0043DeleteSingleQosWithResponse(ctx, qosName)
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