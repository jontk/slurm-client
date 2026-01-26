// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_44

import (
	"context"
	"time"

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
	if m.client == nil || m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0044GetClustersParams{}

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmdbV0044GetClustersWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
		var responseBody []byte
		if resp.JSON200 != nil {
			// Try to extract error details from response
			if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
				apiErrors := make([]errors.SlurmAPIErrorDetail, len(*resp.JSON200.Errors))
				for i, apiErr := range *resp.JSON200.Errors {
					var errorNumber int
					if apiErr.ErrorNumber != nil {
						errorNumber = int(*apiErr.ErrorNumber)
					}
					var errorCode string
					if apiErr.Error != nil {
						errorCode = *apiErr.Error
					}
					var source string
					if apiErr.Source != nil {
						source = *apiErr.Source
					}
					var description string
					if apiErr.Description != nil {
						description = *apiErr.Description
					}

					apiErrors[i] = errors.SlurmAPIErrorDetail{
						ErrorNumber: errorNumber,
						ErrorCode:   errorCode,
						Source:      source,
						Description: description,
					}
				}
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.44", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.44")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response but got nil")
	}

	// Convert the response to our interface types
	clusters := make([]*interfaces.Cluster, 0, len(resp.JSON200.Clusters))
	for _, apiCluster := range resp.JSON200.Clusters {
		cluster := convertAPIClusterToInterface(apiCluster)
		clusters = append(clusters, cluster)
	}

	// Apply client-side filtering if options are provided
	if opts != nil {
		clusters = applyClusterFiltering(clusters, opts)
	}

	return &interfaces.ClusterList{
		Clusters: clusters,
	}, nil
}

// Get retrieves a specific cluster by name
func (m *ClusterManagerImpl) Get(ctx context.Context, clusterName string) (*interfaces.Cluster, error) {
	// Get all clusters and filter by name
	list, err := m.List(ctx, &interfaces.ListClustersOptions{
		Names: []string{clusterName},
	})
	if err != nil {
		return nil, err
	}

	if len(list.Clusters) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "Cluster not found", "cluster", clusterName)
	}

	return list.Clusters[0], nil
}

// Create creates a new cluster
func (m *ClusterManagerImpl) Create(ctx context.Context, cluster *interfaces.ClusterCreate) (*interfaces.ClusterCreateResponse, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Cluster creation not yet implemented for v0.0.44")
}

// Update updates an existing cluster
func (m *ClusterManagerImpl) Update(ctx context.Context, clusterName string, update *interfaces.ClusterUpdate) error {
	return errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Cluster update not yet implemented for v0.0.44")
}

// Delete deletes a cluster
func (m *ClusterManagerImpl) Delete(ctx context.Context, clusterName string) error {
	return errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Cluster deletion not yet implemented for v0.0.44")
}

// convertAPIClusterToInterface converts v0.0.44 API cluster to interface type
func convertAPIClusterToInterface(apiCluster V0044ClusterRec) *interfaces.Cluster {
	cluster := &interfaces.Cluster{
		Metadata: make(map[string]interface{}),
	}

	// Handle required fields
	if apiCluster.Name != nil {
		cluster.Name = *apiCluster.Name
	}

	// Handle controller information (ControlHost and ControlPort)
	if apiCluster.Controller != nil {
		if apiCluster.Controller.Host != nil {
			cluster.ControlHost = *apiCluster.Controller.Host
		}
		if apiCluster.Controller.Port != nil {
			cluster.ControlPort = int(*apiCluster.Controller.Port)
		}
	}

	// RPC version - critical for slurm_clusters_rpc_version metric
	if apiCluster.RpcVersion != nil {
		cluster.RPCVersion = int(*apiCluster.RpcVersion)
	}

	// Node count
	if apiCluster.Nodes != nil {
		// Store in metadata as v0.0.44 returns nodes as string
		cluster.Metadata["nodes"] = *apiCluster.Nodes
	}

	// TRES list
	if apiCluster.Tres != nil && len(*apiCluster.Tres) > 0 {
		tresList := make([]string, 0, len(*apiCluster.Tres))
		for _, tres := range *apiCluster.Tres {
			if tres.Type != "" {
				tresList = append(tresList, tres.Type)
			}
		}
		cluster.TRESList = tresList
	}

	// Flags (features)
	if apiCluster.Flags != nil {
		cluster.Features = make([]string, len(*apiCluster.Flags))
		for i, flag := range *apiCluster.Flags {
			cluster.Features[i] = string(flag)
		}
	}

	// Set timestamps - use current time if not available
	now := time.Now()
	cluster.Created = now
	cluster.Modified = now

	return cluster
}

// applyClusterFiltering applies client-side filtering to cluster list
func applyClusterFiltering(clusters []*interfaces.Cluster, opts *interfaces.ListClustersOptions) []*interfaces.Cluster {
	if opts == nil {
		return clusters
	}

	var filtered []*interfaces.Cluster

	for _, cluster := range clusters {
		// Filter by names if specified
		if len(opts.Names) > 0 {
			found := false
			for _, name := range opts.Names {
				if cluster.Name == name {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		filtered = append(filtered, cluster)
	}

	return filtered
}
