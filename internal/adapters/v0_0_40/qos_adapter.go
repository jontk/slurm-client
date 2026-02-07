// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_40

import (
	"context"

	types "github.com/jontk/slurm-client/api"
	adapterbase "github.com/jontk/slurm-client/internal/adapters/base"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_40"
	"github.com/jontk/slurm-client/pkg/errors"
)

// QoSAdapter implements the QoSAdapter interface for v0.0.40
type QoSAdapter struct {
	*adapterbase.BaseManager
	client *api.ClientWithResponses
}

// NewQoSAdapter creates a new QoS adapter for v0.0.40
func NewQoSAdapter(client *api.ClientWithResponses) *QoSAdapter {
	return &QoSAdapter{
		BaseManager: adapterbase.NewBaseManager("v0.0.40", "QoS"),
		client:      client,
	}
}

// List retrieves a list of QoS with optional filtering
func (a *QoSAdapter) List(ctx context.Context, opts *types.QoSListOptions) (*types.QoSList, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	// v0.0.40 may not have QoS endpoints, return empty list for now
	return &types.QoSList{
		QoS:   []types.QoS{},
		Total: 0,
	}, nil
}

// Get retrieves a specific QoS by name
func (a *QoSAdapter) Get(ctx context.Context, qosName string) (*types.QoS, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.ValidateResourceName(qosName, "QoS name"); err != nil {
		return nil, err
	}
	// v0.0.40 may not have QoS endpoints
	return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, "QoS '"+qosName+"' not found")
}

// Create creates a new QoS
func (a *QoSAdapter) Create(ctx context.Context, qos *types.QoSCreate) (*types.QoSCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	// v0.0.40 may not have QoS endpoints
	return nil, errors.NewNotImplementedError("Create QoS", "v0.0.40")
}

// Update updates an existing QoS
func (a *QoSAdapter) Update(ctx context.Context, qosName string, update *types.QoSUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(qosName, "QoS name"); err != nil {
		return err
	}
	// v0.0.40 may not have QoS endpoints
	return errors.NewNotImplementedError("Update QoS", "v0.0.40")
}

// Delete deletes a QoS
func (a *QoSAdapter) Delete(ctx context.Context, qosName string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(qosName, "QoS name"); err != nil {
		return err
	}
	// v0.0.40 may not have QoS endpoints
	return errors.NewNotImplementedError("Delete QoS", "v0.0.40")
}
