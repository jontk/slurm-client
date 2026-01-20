// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_44

import (
	"context"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// ClusterManagerImpl provides the actual implementation for ClusterManager methods
type ClusterManagerImpl struct {
	client *WrapperClient
}

// NewClusterManagerImpl creates a new ClusterManager implementation
func NewClusterManagerImpl(client *WrapperClient) *ClusterManagerImpl {
	return &ClusterManagerImpl{client: client}
}

// List clusters with optional filtering
func (m *ClusterManagerImpl) List(ctx context.Context, opts *interfaces.ListClustersOptions) (*interfaces.ClusterList, error) {
	return &interfaces.ClusterList{Clusters: make([]*interfaces.Cluster, 0)}, nil
}

// Get retrieves a specific cluster by name
func (m *ClusterManagerImpl) Get(ctx context.Context, clusterName string) (*interfaces.Cluster, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Cluster management not yet implemented for v0.0.44")
}

// Create creates a new cluster
func (m *ClusterManagerImpl) Create(ctx context.Context, cluster *interfaces.ClusterCreate) (*interfaces.ClusterCreateResponse, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Cluster management not yet implemented for v0.0.44")
}

// Update updates an existing cluster
func (m *ClusterManagerImpl) Update(ctx context.Context, clusterName string, update *interfaces.ClusterUpdate) error {
	return errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Cluster management not yet implemented for v0.0.44")
}

// Delete deletes a cluster
func (m *ClusterManagerImpl) Delete(ctx context.Context, clusterName string) error {
	return errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Cluster management not yet implemented for v0.0.44")
}
