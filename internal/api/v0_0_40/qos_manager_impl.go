package v0_0_40

import (
	"context"

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
	// v0.0.40 is an older API version with limited support
	// QoS management was added in later versions
	return nil, errors.NewNotImplementedError("QoS listing", "v0.0.40")
}

// Get retrieves a specific QoS by name
func (q *QoSManagerImpl) Get(ctx context.Context, qosName string) (*interfaces.QoS, error) {
	if qosName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS name is required", "qosName", qosName, nil)
	}
	return nil, errors.NewNotImplementedError("QoS retrieval", "v0.0.40")
}

// Create creates a new QoS
func (q *QoSManagerImpl) Create(ctx context.Context, qos *interfaces.QoSCreate) (*interfaces.QoSCreateResponse, error) {
	if qos == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS data is required", "qos", qos, nil)
	}
	return nil, errors.NewNotImplementedError("QoS creation", "v0.0.40")
}

// Update updates an existing QoS
func (q *QoSManagerImpl) Update(ctx context.Context, qosName string, update *interfaces.QoSUpdate) error {
	if qosName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS name is required", "qosName", qosName, nil)
	}
	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "update data is required", "update", update, nil)
	}
	return errors.NewNotImplementedError("QoS update", "v0.0.40")
}

// Delete deletes a QoS
func (q *QoSManagerImpl) Delete(ctx context.Context, qosName string) error {
	if qosName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS name is required", "qosName", qosName, nil)
	}
	return errors.NewNotImplementedError("QoS deletion", "v0.0.40")
}