package v0_0_40

import (
	"context"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
)

// QoSAdapter implements the QoSAdapter interface for v0.0.40
type QoSAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewQoSAdapter creates a new QoS adapter for v0.0.40
func NewQoSAdapter(client *api.ClientWithResponses) *QoSAdapter {
	return &QoSAdapter{
		BaseManager: base.NewBaseManager("v0.0.40", "QoS"),
		client:      client,
		wrapper:     nil,
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
	if err := a.ValidateResourceName(qosName, "qosName"); err != nil {
		return nil, err
	}

	// v0.0.40 may not have QoS endpoints
	return nil, common.NewResourceNotFoundError("QoS", qosName)
}

// Create creates a new QoS
func (a *QoSAdapter) Create(ctx context.Context, qos *types.QoSCreate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// v0.0.40 may not have QoS endpoints
	return common.NewNotImplementedError("Create QoS is not implemented for v0.0.40")
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

	// v0.0.40 may not have QoS endpoints
	return common.NewNotImplementedError("Update QoS is not implemented for v0.0.40")
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

	// v0.0.40 may not have QoS endpoints
	return common.NewNotImplementedError("Delete QoS is not implemented for v0.0.40")
}