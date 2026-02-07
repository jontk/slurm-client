// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"context"
	"fmt"
	"net/http"

	types "github.com/jontk/slurm-client/api"
	adapterbase "github.com/jontk/slurm-client/internal/adapters/base"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
)

// ClusterAdapter implements the ClusterManager interface for v0.0.41
type ClusterAdapter struct {
	*adapterbase.BaseManager
	client *api.ClientWithResponses
}

// NewClusterAdapter creates a new Cluster adapter for v0.0.41
func NewClusterAdapter(client *api.ClientWithResponses) *ClusterAdapter {
	return &ClusterAdapter{
		BaseManager: adapterbase.NewBaseManager("v0.0.41", "Cluster"),
		client:      client,
	}
}

// List retrieves a list of clusters
func (a *ClusterAdapter) List(ctx context.Context, opts *types.ClusterListOptions) (*types.ClusterList, error) {
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}
	params := &api.SlurmdbV0041GetClustersParams{}
	// No parameters seem to be available for filtering in the API spec for v0.0.41
	resp, err := a.client.SlurmdbV0041GetClustersWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list clusters: %w", err)
	}
	if resp.HTTPResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: failed to list clusters", resp.HTTPResponse.StatusCode)
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected nil response")
	}
	return toInternalClusterList(resp.JSON200)
}

// Get retrieves a specific cluster by name
func (a *ClusterAdapter) Get(ctx context.Context, clusterName string) (*types.Cluster, error) {
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.ValidateResourceName("cluster name", clusterName); err != nil {
		return nil, err
	}
	list, err := a.List(ctx, nil)
	if err != nil {
		return nil, err
	}
	for _, cluster := range list.Clusters {
		if cluster.Name != nil && *cluster.Name == clusterName {
			return &cluster, nil
		}
	}
	return nil, fmt.Errorf("cluster %s not found", clusterName)
}

// Create creates a new cluster
func (a *ClusterAdapter) Create(ctx context.Context, cluster *types.ClusterCreate) (*types.ClusterCreateResponse, error) {
	return nil, a.HandleNotImplemented("Create", "v0.0.41 cluster adapter")
}

// Update updates an existing cluster
func (a *ClusterAdapter) Update(ctx context.Context, clusterName string, update *types.ClusterUpdate) error {
	return a.HandleNotImplemented("Update", "v0.0.41 cluster adapter")
}

// Delete deletes a cluster
func (a *ClusterAdapter) Delete(ctx context.Context, clusterName string) error {
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName("cluster name", clusterName); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}
	resp, err := a.client.SlurmdbV0041DeleteClusterWithResponse(ctx, clusterName, &api.SlurmdbV0041DeleteClusterParams{})
	if err != nil {
		return fmt.Errorf("failed to delete cluster %s: %w", clusterName, err)
	}
	if resp.HTTPResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: failed to delete cluster %s", resp.HTTPResponse.StatusCode, clusterName)
	}
	return nil
}

// AddNodes adds nodes to a cluster
func (a *ClusterAdapter) AddNodes(ctx context.Context, clusterName string, nodeNames []string) error {
	return a.HandleNotImplemented("AddNodes", "v0.0.41 cluster adapter")
}

// RemoveNodes removes nodes from a cluster
func (a *ClusterAdapter) RemoveNodes(ctx context.Context, clusterName string, nodeNames []string) error {
	return a.HandleNotImplemented("RemoveNodes", "v0.0.41 cluster adapter")
}
