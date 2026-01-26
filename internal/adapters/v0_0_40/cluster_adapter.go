// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"context"
	"fmt"
	"strconv"

	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
)

// ClusterAdapter implements the ClusterAdapter interface for v0.0.40
type ClusterAdapter struct {
	*base.BaseManager
	client *api.ClientWithResponses
}

// NewClusterAdapter creates a new Cluster adapter for v0.0.40
func NewClusterAdapter(client *api.ClientWithResponses) *ClusterAdapter {
	return &ClusterAdapter{
		BaseManager: base.NewBaseManager("v0.0.40", "Cluster"),
		client:      client,
	}
}

// List retrieves a list of clusters with optional filtering
func (a *ClusterAdapter) List(ctx context.Context, opts *types.ClusterListOptions) (*types.ClusterList, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmdbV0040GetClustersParams{}

	// Apply filters from options
	if opts != nil {
		if opts.UpdateTime != nil {
			updateTimeStr := strconv.FormatInt(opts.UpdateTime.Unix(), 10)
			params.UpdateTime = &updateTimeStr
		}
	}

	// Call the API
	resp, err := a.client.SlurmdbV0040GetClustersWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list clusters")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}

	// Check for API response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from API")
	}

	// Convert the response to common types
	clusters := make([]types.Cluster, 0)

	if len(resp.JSON200.Clusters) > 0 {
		for _, apiCluster := range resp.JSON200.Clusters {
			cluster, err := a.convertAPIClusterToCommon(apiCluster)
			if err != nil {
				// Log conversion error but continue
				continue
			}
			clusters = append(clusters, *cluster)
		}
	}

	return &types.ClusterList{
		Clusters: clusters,
		Meta:     a.extractMeta(resp.JSON200.Meta),
	}, nil
}

// Get retrieves a specific cluster by name
func (a *ClusterAdapter) Get(ctx context.Context, clusterName string) (*types.Cluster, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmdbV0040GetClusterParams{}

	// Call the API
	resp, err := a.client.SlurmdbV0040GetClusterWithResponse(ctx, clusterName, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to get cluster "+clusterName)
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}

	// Check for API response
	if resp.JSON200 == nil || len(resp.JSON200.Clusters) == 0 {
		return nil, fmt.Errorf("cluster %s not found", clusterName)
	}

	// Convert the first cluster in the response
	clusters := resp.JSON200.Clusters
	return a.convertAPIClusterToCommon(clusters[0])
}

// Create creates a new cluster
func (a *ClusterAdapter) Create(ctx context.Context, cluster *types.ClusterCreate) (*types.ClusterCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert common cluster create to API request
	apiCluster := &api.V0040ClusterRec{
		Name: &cluster.Name,
	}

	if cluster.ControllerHost != "" {
		apiCluster.Controller = &struct {
			Host *string `json:"host,omitempty"`
			Port *int32  `json:"port,omitempty"`
		}{
			Host: &cluster.ControllerHost,
		}
		if cluster.ControllerPort > 0 {
			apiCluster.Controller.Port = &cluster.ControllerPort
		}
	}

	if cluster.Nodes != "" {
		apiCluster.Nodes = &cluster.Nodes
	}

	if cluster.RpcVersion > 0 {
		apiCluster.RpcVersion = &cluster.RpcVersion
	}

	if cluster.SelectPlugin != "" {
		apiCluster.SelectPlugin = &cluster.SelectPlugin //nolint:staticcheck // SA1019: Deprecated upstream but required for backward compatibility
	}

	if len(cluster.Flags) > 0 {
		flags := cluster.Flags
		apiCluster.Flags = &flags
	}

	// Create request body
	apiReq := api.V0040OpenapiClustersResp{
		Clusters: []api.V0040ClusterRec{*apiCluster},
	}

	// Call the API
	params := &api.SlurmdbV0040PostClustersParams{}
	resp, err := a.client.SlurmdbV0040PostClustersWithResponse(ctx, params, apiReq)
	if err != nil {
		return nil, a.WrapError(err, "failed to create cluster")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}

	// Check for API response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from API")
	}

	// Convert response - v0.0.40 POST returns V0040OpenapiResp not V0040OpenapiClustersResp
	return a.convertAPIClusterCreateResponseToCommon(resp.JSON200, cluster.Name)
}

// Delete deletes a cluster
func (a *ClusterAdapter) Delete(ctx context.Context, clusterName string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Prepare parameters for the API call
	params := &api.SlurmdbV0040DeleteClusterParams{}

	// Call the API
	resp, err := a.client.SlurmdbV0040DeleteClusterWithResponse(ctx, clusterName, params)
	if err != nil {
		return a.WrapError(err, "failed to delete cluster "+clusterName)
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}

	return nil
}

// convertAPIClusterToCommon converts API cluster to common type
func (a *ClusterAdapter) convertAPIClusterToCommon(apiCluster api.V0040ClusterRec) (*types.Cluster, error) {
	cluster := &types.Cluster{
		Meta: make(map[string]interface{}),
	}

	// Set basic fields
	if apiCluster.Name != nil {
		cluster.Name = *apiCluster.Name
	}

	if apiCluster.Controller != nil {
		if apiCluster.Controller.Host != nil {
			cluster.ControllerHost = *apiCluster.Controller.Host
		}
		if apiCluster.Controller.Port != nil {
			cluster.ControllerPort = *apiCluster.Controller.Port
		}
	}

	if apiCluster.Nodes != nil {
		cluster.Nodes = *apiCluster.Nodes
	}

	if apiCluster.RpcVersion != nil {
		cluster.RpcVersion = *apiCluster.RpcVersion
	}

	//nolint:staticcheck // SA1019: Deprecated upstream but required for backward compatibility
	if apiCluster.SelectPlugin != nil {
		cluster.SelectPlugin = *apiCluster.SelectPlugin
	}

	if apiCluster.Flags != nil {
		cluster.Flags = *apiCluster.Flags
	}

	// Convert TRES if available
	if apiCluster.Tres != nil {
		tres := make([]types.TRES, 0, len(*apiCluster.Tres))
		for _, apiTres := range *apiCluster.Tres {
			t := types.TRES{
				Type: apiTres.Type,
			}
			if apiTres.Id != nil {
				t.ID = int(*apiTres.Id)
			}
			if apiTres.Name != nil {
				t.Name = *apiTres.Name
			}
			if apiTres.Count != nil {
				t.Count = *apiTres.Count
			}
			tres = append(tres, t)
		}
		cluster.TRES = tres
	}

	// Convert associations if available
	if apiCluster.Associations != nil && apiCluster.Associations.Root != nil {
		assocShort := &types.AssocShort{
			User: apiCluster.Associations.Root.User, // User is not a pointer
		}

		if apiCluster.Associations.Root.Account != nil {
			assocShort.Account = *apiCluster.Associations.Root.Account
		}
		if apiCluster.Associations.Root.Cluster != nil {
			assocShort.Cluster = *apiCluster.Associations.Root.Cluster
		}
		if apiCluster.Associations.Root.Partition != nil {
			assocShort.Partition = *apiCluster.Associations.Root.Partition
		}

		cluster.Associations = &types.AssociationShort{
			Root: assocShort,
		}
	}

	return cluster, nil
}

// convertAPIClusterCreateResponseToCommon converts API create response to common type
func (a *ClusterAdapter) convertAPIClusterCreateResponseToCommon(apiResp *api.V0040OpenapiResp, name string) (*types.ClusterCreateResponse, error) {
	resp := &types.ClusterCreateResponse{
		Name:   name,
		Status: "success",
		Meta:   make(map[string]interface{}),
	}

	// V0040OpenapiResp doesn't contain clusters - it's a general response
	// We cannot extract details from the response for v0.0.40

	// Extract metadata if available
	if apiResp.Meta != nil {
		resp.Meta = a.extractMeta(apiResp.Meta)
	}

	// Handle errors in response - V0040OpenapiErrors is []V0040OpenapiError
	if apiResp.Errors != nil && len(*apiResp.Errors) > 0 {
		resp.Status = "error"
		errors := *apiResp.Errors
		if len(errors) > 0 && errors[0].Error != nil {
			resp.Message = *errors[0].Error
		} else {
			resp.Message = "Cluster creation failed"
		}
	} else {
		resp.Message = fmt.Sprintf("Cluster '%s' created successfully", name)
	}

	return resp, nil
}

// extractMeta safely extracts metadata from API response
func (a *ClusterAdapter) extractMeta(meta *api.V0040OpenapiMeta) map[string]interface{} {
	result := make(map[string]interface{})

	if meta == nil {
		return result
	}

	// Extract basic metadata
	if meta.Client != nil {
		clientInfo := make(map[string]interface{})
		if meta.Client.Source != nil {
			clientInfo["source"] = *meta.Client.Source
		}
		if meta.Client.User != nil {
			clientInfo["user"] = *meta.Client.User
		}
		if meta.Client.Group != nil {
			clientInfo["group"] = *meta.Client.Group
		}
		if len(clientInfo) > 0 {
			result["client"] = clientInfo
		}
	}

	if meta.Plugin != nil {
		pluginInfo := make(map[string]interface{})
		if meta.Plugin.AccountingStorage != nil {
			pluginInfo["accounting_storage"] = *meta.Plugin.AccountingStorage
		}
		if len(pluginInfo) > 0 {
			result["plugin"] = pluginInfo
		}
	}

	return result
}
