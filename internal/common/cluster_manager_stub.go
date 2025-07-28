package common

import (
	"context"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// ClusterManagerStub provides a stub implementation for API versions that don't support cluster management
type ClusterManagerStub struct {
	Version string
}

// NewClusterManagerStub creates a new ClusterManagerStub
func NewClusterManagerStub(version string) *ClusterManagerStub {
	return &ClusterManagerStub{
		Version: version,
	}
}

// List returns not implemented error
func (c *ClusterManagerStub) List(ctx context.Context, opts *interfaces.ListClustersOptions) (*interfaces.ClusterList, error) {
	return nil, errors.NewNotImplementedError("cluster listing", c.Version)
}

// Get returns not implemented error
func (c *ClusterManagerStub) Get(ctx context.Context, clusterName string) (*interfaces.Cluster, error) {
	return nil, errors.NewNotImplementedError("cluster retrieval", c.Version)
}

// Create returns not implemented error
func (c *ClusterManagerStub) Create(ctx context.Context, cluster *interfaces.ClusterCreate) (*interfaces.ClusterCreateResponse, error) {
	return nil, errors.NewNotImplementedError("cluster creation", c.Version)
}

// Update returns not implemented error
func (c *ClusterManagerStub) Update(ctx context.Context, clusterName string, update *interfaces.ClusterUpdate) error {
	return errors.NewNotImplementedError("cluster update", c.Version)
}

// Delete returns not implemented error
func (c *ClusterManagerStub) Delete(ctx context.Context, clusterName string) error {
	return errors.NewNotImplementedError("cluster deletion", c.Version)
}