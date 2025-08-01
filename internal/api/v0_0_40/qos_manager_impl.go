// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"context"
	"strings"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// QoSManagerImpl implements the QoSManager interface for v0.0.40
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
	if err := q.client.CheckContext(ctx); err != nil {
		return nil, err
	}

	// Prepare parameters
	params := &SlurmdbV0040GetQosParams{}
	if opts != nil {
		if opts.Names != nil && len(opts.Names) > 0 {
			namesStr := strings.Join(opts.Names, ",")
			params.Name = &namesStr
		}
		// WithDeleted field doesn't exist in ListQoSOptions
		// Skip this parameter
	}

	// Make API call
	resp, err := q.client.apiClient.SlurmdbV0040GetQosWithResponse(ctx, params)
	if err != nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "failed to list QoS")
	}

	// Check response
	if resp.StatusCode() != 200 {
		return nil, q.client.HandleErrorResponse(resp.StatusCode(), resp.Body)
	}

	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "received nil response")
	}

	// Convert response
	qosList := &interfaces.QoSList{
		QoS: make([]interfaces.QoS, 0),
	}

	if resp.JSON200.Qos != nil {
		for _, qosItem := range resp.JSON200.Qos {
			qos := q.convertV0040QoSToInterface(qosItem)
			qosList.QoS = append(qosList.QoS, *qos)
		}
	}

	return qosList, nil
}

// Get retrieves a specific QoS by name
func (q *QoSManagerImpl) Get(ctx context.Context, qosName string) (*interfaces.QoS, error) {
	if qosName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS name is required", "qosName", qosName, nil)
	}

	if err := q.client.CheckContext(ctx); err != nil {
		return nil, err
	}

	// Make API call
	params := &SlurmdbV0040GetSingleQosParams{}
	resp, err := q.client.apiClient.SlurmdbV0040GetSingleQosWithResponse(ctx, qosName, params)
	if err != nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "failed to get QoS")
	}

	// Check response
	if resp.StatusCode() != 200 {
		return nil, q.client.HandleErrorResponse(resp.StatusCode(), resp.Body)
	}

	if resp.JSON200 == nil || resp.JSON200.Qos == nil || len(resp.JSON200.Qos) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "QoS not found")
	}

	// Convert the first QoS
	qosList := resp.JSON200.Qos
	return q.convertV0040QoSToInterface(qosList[0]), nil
}

// Create creates a new QoS
func (q *QoSManagerImpl) Create(ctx context.Context, qos *interfaces.QoSCreate) (*interfaces.QoSCreateResponse, error) {
	if qos == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS data is required", "qos", qos, nil)
	}

	if err := q.client.CheckContext(ctx); err != nil {
		return nil, err
	}

	// Convert to API format
	apiQoS := q.convertInterfaceQoSCreateToV0040(qos)
	
	// Create request body
	reqBody := V0040OpenapiSlurmdbdQosResp{
		Qos: []V0040Qos{*apiQoS},
	}

	// Prepare parameters
	params := &SlurmdbV0040PostQosParams{}

	// Make API call
	resp, err := q.client.apiClient.SlurmdbV0040PostQosWithResponse(ctx, params, reqBody)
	if err != nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "failed to create QoS")
	}

	// Check response
	if resp.StatusCode() != 200 {
		return nil, q.client.HandleErrorResponse(resp.StatusCode(), resp.Body)
	}

	return &interfaces.QoSCreateResponse{
		QoSName: qos.Name,
	}, nil
}

// Update updates an existing QoS
func (q *QoSManagerImpl) Update(ctx context.Context, qosName string, update *interfaces.QoSUpdate) error {
	if qosName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS name is required", "qosName", qosName, nil)
	}
	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "update data is required", "update", update, nil)
	}

	if err := q.client.CheckContext(ctx); err != nil {
		return err
	}

	// Get existing QoS first
	existing, err := q.Get(ctx, qosName)
	if err != nil {
		return err
	}

	// Apply updates to existing QoS
	apiQoS := q.convertInterfaceQoSToV0040Update(existing, update)

	// Create request body
	reqBody := V0040OpenapiSlurmdbdQosResp{
		Qos: []V0040Qos{*apiQoS},
	}

	// Prepare parameters
	params := &SlurmdbV0040PostQosParams{}

	// Make API call
	resp, err := q.client.apiClient.SlurmdbV0040PostQosWithResponse(ctx, params, reqBody)
	if err != nil {
		return errors.NewClientError(errors.ErrorCodeServerInternal, "failed to update QoS")
	}

	// Check response
	if resp.StatusCode() != 200 {
		return q.client.HandleErrorResponse(resp.StatusCode(), resp.Body)
	}

	return nil
}

// Delete deletes a QoS
func (q *QoSManagerImpl) Delete(ctx context.Context, qosName string) error {
	if qosName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS name is required", "qosName", qosName, nil)
	}

	if err := q.client.CheckContext(ctx); err != nil {
		return err
	}

	// Make API call
	resp, err := q.client.apiClient.SlurmdbV0040DeleteSingleQosWithResponse(ctx, qosName)
	if err != nil {
		return errors.NewClientError(errors.ErrorCodeServerInternal, "failed to delete QoS")
	}

	// Check response
	if resp.StatusCode() != 200 && resp.StatusCode() != 204 {
		return q.client.HandleErrorResponse(resp.StatusCode(), resp.Body)
	}

	return nil
}

// convertV0040QoSToInterface converts v0.0.40 QoS to interface format
func (q *QoSManagerImpl) convertV0040QoSToInterface(qos V0040Qos) *interfaces.QoS {
	result := &interfaces.QoS{}
	
	if qos.Name != nil {
		result.Name = *qos.Name
	}
	if qos.Description != nil {
		result.Description = *qos.Description
	}
	if qos.Priority != nil {
		if qos.Priority.Number != nil {
			result.Priority = int(*qos.Priority.Number)
		}
	}
	// PreemptMode field doesn't exist in V0040Qos
	// Skip preempt mode conversion
	
	return result
}

// convertInterfaceQoSCreateToV0040 converts interface QoS create to v0.0.40 format
func (q *QoSManagerImpl) convertInterfaceQoSCreateToV0040(qos *interfaces.QoSCreate) *V0040Qos {
	apiQoS := &V0040Qos{
		Name:        &qos.Name,
		Description: &qos.Description,
	}
	
	if qos.Priority > 0 {
		set := true
		number := int64(qos.Priority)
		apiQoS.Priority = &V0040Uint32NoVal{
			Set:    &set,
			Number: &number,
		}
	}
	
	return apiQoS
}

// convertInterfaceQoSToV0040Update converts interface QoS with updates to v0.0.40 format
func (q *QoSManagerImpl) convertInterfaceQoSToV0040Update(existing *interfaces.QoS, update *interfaces.QoSUpdate) *V0040Qos {
	apiQoS := &V0040Qos{
		Name: &existing.Name,
	}
	
	// Apply updates
	if update.Description != nil {
		apiQoS.Description = update.Description
	} else {
		apiQoS.Description = &existing.Description
	}
	
	if update.Priority != nil {
		set := true
		number := int64(*update.Priority)
		apiQoS.Priority = &V0040Uint32NoVal{
			Set:    &set,
			Number: &number,
		}
	} else if existing.Priority > 0 {
		set := true
		number := int64(existing.Priority)
		apiQoS.Priority = &V0040Uint32NoVal{
			Set:    &set,
			Number: &number,
		}
	}
	
	return apiQoS
}
