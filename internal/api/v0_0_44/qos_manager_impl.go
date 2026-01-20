// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_44

import (
	"context"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// QoSManagerImpl provides the actual implementation for QoSManager methods
type QoSManagerImpl struct {
	client *WrapperClient
}

// NewQoSManagerImpl creates a new QoSManager implementation
func NewQoSManagerImpl(client *WrapperClient) *QoSManagerImpl {
	return &QoSManagerImpl{client: client}
}

// List retrieves a list of QoS with optional filtering
func (m *QoSManagerImpl) List(ctx context.Context, opts *interfaces.ListQoSOptions) (*interfaces.QoSList, error) {
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// TODO: Implement QoS listing when endpoints are available
	return &interfaces.QoSList{QoS: make([]interfaces.QoS, 0)}, nil
}

// Get retrieves a specific QoS by name
func (m *QoSManagerImpl) Get(ctx context.Context, qosName string) (*interfaces.QoS, error) {
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// TODO: Implement QoS retrieval when endpoints are available
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Get QoS not yet implemented for v0.0.44")
}

// Create creates a new QoS
func (m *QoSManagerImpl) Create(ctx context.Context, qos *interfaces.QoSCreate) (*interfaces.QoSCreateResponse, error) {
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// TODO: Implement QoS creation when endpoints are available
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Create QoS not yet implemented for v0.0.44")
}

// Update updates an existing QoS
func (m *QoSManagerImpl) Update(ctx context.Context, qosName string, update *interfaces.QoSUpdate) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// TODO: Implement QoS update when endpoints are available
	return errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Update QoS not yet implemented for v0.0.44")
}

// Delete deletes a QoS
func (m *QoSManagerImpl) Delete(ctx context.Context, qosName string) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// TODO: Implement QoS deletion when endpoints are available
	return errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Delete QoS not yet implemented for v0.0.44")
}
