package common

import (
	"context"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// AssociationManagerStub provides a stub implementation for API versions that don't support association management
type AssociationManagerStub struct {
	Version string
}

// NewAssociationManagerStub creates a new AssociationManagerStub
func NewAssociationManagerStub(version string) *AssociationManagerStub {
	return &AssociationManagerStub{
		Version: version,
	}
}

// List returns not implemented error
func (a *AssociationManagerStub) List(ctx context.Context, opts *interfaces.ListAssociationsOptions) (*interfaces.AssociationList, error) {
	return nil, errors.NewNotImplementedError("association listing", a.Version)
}

// Get returns not implemented error
func (a *AssociationManagerStub) Get(ctx context.Context, opts *interfaces.GetAssociationOptions) (*interfaces.Association, error) {
	return nil, errors.NewNotImplementedError("association retrieval", a.Version)
}

// Create returns not implemented error
func (a *AssociationManagerStub) Create(ctx context.Context, associations []*interfaces.AssociationCreate) (*interfaces.AssociationCreateResponse, error) {
	return nil, errors.NewNotImplementedError("association creation", a.Version)
}

// Update returns not implemented error
func (a *AssociationManagerStub) Update(ctx context.Context, associations []*interfaces.AssociationUpdate) error {
	return errors.NewNotImplementedError("association update", a.Version)
}

// Delete returns not implemented error
func (a *AssociationManagerStub) Delete(ctx context.Context, opts *interfaces.DeleteAssociationOptions) error {
	return errors.NewNotImplementedError("association deletion", a.Version)
}

// BulkDelete returns not implemented error
func (a *AssociationManagerStub) BulkDelete(ctx context.Context, opts *interfaces.BulkDeleteOptions) (*interfaces.BulkDeleteResponse, error) {
	return nil, errors.NewNotImplementedError("bulk association deletion", a.Version)
}

// GetUserAssociations returns not implemented error
func (a *AssociationManagerStub) GetUserAssociations(ctx context.Context, userName string) ([]*interfaces.Association, error) {
	return nil, errors.NewNotImplementedError("user association retrieval", a.Version)
}

// GetAccountAssociations returns not implemented error
func (a *AssociationManagerStub) GetAccountAssociations(ctx context.Context, accountName string) ([]*interfaces.Association, error) {
	return nil, errors.NewNotImplementedError("account association retrieval", a.Version)
}

// ValidateAssociation returns not implemented error
func (a *AssociationManagerStub) ValidateAssociation(ctx context.Context, user, account, cluster string) (bool, error) {
	return false, errors.NewNotImplementedError("association validation", a.Version)
}