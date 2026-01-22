// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"net/http"
	"context"
	"fmt"
	"time"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// ClusterManagerImpl implements the ClusterManager interface for v0.0.43
type ClusterManagerImpl struct {
	client *WrapperClient
}

// NewClusterManagerImpl creates a new ClusterManagerImpl
func NewClusterManagerImpl(client *WrapperClient) *ClusterManagerImpl {
	return &ClusterManagerImpl{
		client: client,
	}
}

// List lists clusters with optional filtering
func (c *ClusterManagerImpl) List(ctx context.Context, opts *interfaces.ListClustersOptions) (*interfaces.ClusterList, error) {
	if c.client == nil || c.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0043GetClustersParams{}

	// Call the generated OpenAPI client
	resp, err := c.client.apiClient.SlurmdbV0043GetClustersWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response but got nil")
	}

	// Convert the response to our interface types
	clusters := make([]*interfaces.Cluster, 0, len(resp.JSON200.Clusters))
	for _, apiCluster := range resp.JSON200.Clusters {
		cluster, err := convertAPIClusterToInterface(apiCluster)
		if err != nil {
			conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert cluster data")
			conversionErr.Cause = err
			if apiCluster.Name != nil {
				conversionErr.Details = fmt.Sprintf("Error converting cluster %s", *apiCluster.Name)
			}
			return nil, conversionErr
		}
		clusters = append(clusters, cluster)
	}

	// Apply client-side filtering if options are provided
	if opts != nil {
		clusters = applyClusterFiltering(clusters, opts)
	}

	// Apply pagination if requested
	totalCount := len(clusters)
	if opts != nil && opts.Limit > 0 {
		start := opts.Offset
		end := start + opts.Limit
		if start >= len(clusters) {
			clusters = []*interfaces.Cluster{}
		} else {
			if end > len(clusters) {
				end = len(clusters)
			}
			clusters = clusters[start:end]
		}
	}

	return &interfaces.ClusterList{
		Clusters: clusters,
		Total:    totalCount,
		Metadata: map[string]interface{}{"version": "v0.0.43"},
	}, nil
}

// Get retrieves a specific cluster by name
func (c *ClusterManagerImpl) Get(ctx context.Context, clusterName string) (*interfaces.Cluster, error) {
	// Validate input first (cheap check)
	if clusterName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "cluster name is required", "clusterName", clusterName, nil)
	}

	// Then check client initialization
	if c.client == nil || c.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare get parameters
	params := &SlurmdbV0043GetClusterParams{}

	// Call the generated OpenAPI client
	resp, err := c.client.apiClient.SlurmdbV0043GetClusterWithResponse(ctx, clusterName, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() == 404 {
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "cluster not found", clusterName)
	}

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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response but got nil")
	}

	// The response should contain clusters
	if len(resp.JSON200.Clusters) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "cluster not found", clusterName)
	}

	// Convert the first cluster in the response
	cluster, err := convertAPIClusterToInterface(resp.JSON200.Clusters[0])
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert cluster data")
		conversionErr.Cause = err
		conversionErr.Details = fmt.Sprintf("Error converting cluster %s", clusterName)
		return nil, conversionErr
	}

	return cluster, nil
}

// Create creates a new cluster configuration
func (c *ClusterManagerImpl) Create(ctx context.Context, cluster *interfaces.ClusterCreate) (*interfaces.ClusterCreateResponse, error) {
	if c.client == nil || c.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if cluster == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "cluster data is required", "cluster", cluster, nil)
	}

	if cluster.Name == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "cluster name is required", "cluster.Name", cluster.Name, nil)
	}

	// Convert interface cluster to API cluster
	apiCluster, err := convertInterfaceClusterCreateToAPI(cluster)
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert cluster data for API")
		conversionErr.Cause = err
		return nil, conversionErr
	}

	// Prepare the request body
	requestBody := SlurmdbV0043PostClustersJSONRequestBody{
		Clusters: []V0043ClusterRec{*apiCluster},
	}

	// Prepare post parameters
	params := &SlurmdbV0043PostClustersParams{}

	// Call the generated OpenAPI client
	resp, err := c.client.apiClient.SlurmdbV0043PostClustersWithResponse(ctx, params, requestBody)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 && resp.StatusCode() != http.StatusCreated {
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
		return nil, httpErr
	}

	return &interfaces.ClusterCreateResponse{
		Name:     cluster.Name,
		Created:  time.Now(),
		Metadata: map[string]interface{}{"version": "v0.0.43"},
	}, nil
}

// Update updates an existing cluster configuration
func (c *ClusterManagerImpl) Update(ctx context.Context, clusterName string, update *interfaces.ClusterUpdate) error {
	if c.client == nil || c.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if clusterName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "cluster name is required", "clusterName", clusterName, nil)
	}

	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "update data is required", "update", update, nil)
	}

	// Convert interface update to API cluster
	apiCluster, err := convertInterfaceClusterUpdateToAPI(clusterName, update)
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert cluster update data for API")
		conversionErr.Cause = err
		return conversionErr
	}

	// Prepare the request body
	requestBody := SlurmdbV0043PostClustersJSONRequestBody{
		Clusters: []V0043ClusterRec{*apiCluster},
	}

	// Prepare post parameters
	params := &SlurmdbV0043PostClustersParams{}

	// Call the generated OpenAPI client (using POST for updates in SLURM API)
	resp, err := c.client.apiClient.SlurmdbV0043PostClustersWithResponse(ctx, params, requestBody)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() == 404 {
		return errors.NewClientError(errors.ErrorCodeResourceNotFound, "cluster not found", clusterName)
	}

	if resp.StatusCode() != 200 && resp.StatusCode() != http.StatusCreated {
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
		return httpErr
	}

	return nil
}

// Delete deletes a cluster configuration
func (c *ClusterManagerImpl) Delete(ctx context.Context, clusterName string) error {
	if c.client == nil || c.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if clusterName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "cluster name is required", "clusterName", clusterName, nil)
	}

	// Prepare delete parameters
	params := &SlurmdbV0043DeleteClusterParams{}

	// Call the generated OpenAPI client
	resp, err := c.client.apiClient.SlurmdbV0043DeleteClusterWithResponse(ctx, clusterName, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() == 404 {
		return errors.NewClientError(errors.ErrorCodeResourceNotFound, "cluster not found", clusterName)
	}

	if resp.StatusCode() != 200 && resp.StatusCode() != 204 {
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
		return httpErr
	}

	return nil
}

// Helper functions for data conversion

// convertAPIClusterToInterface converts an API cluster response to interface type
func convertAPIClusterToInterface(apiCluster V0043ClusterRec) (*interfaces.Cluster, error) {
	cluster := &interfaces.Cluster{
		Metadata: make(map[string]interface{}),
	}

	// Handle required fields
	if apiCluster.Name != nil {
		cluster.Name = *apiCluster.Name
	}

	// Handle optional fields
	if apiCluster.Controller != nil {
		if apiCluster.Controller.Host != nil {
			cluster.ControlHost = *apiCluster.Controller.Host
		}
		if apiCluster.Controller.Port != nil {
			cluster.ControlPort = int(*apiCluster.Controller.Port)
		}
	}
	if apiCluster.RpcVersion != nil {
		cluster.RPCVersion = int(*apiCluster.RpcVersion)
	}
	if apiCluster.Tres != nil {
		// Handle TRES conversion - this may need adjustment based on actual TRES structure
		cluster.TRESList = []string{} // Placeholder
	}
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

	return cluster, nil
}

// convertInterfaceClusterCreateToAPI converts interface cluster create to API type
func convertInterfaceClusterCreateToAPI(cluster *interfaces.ClusterCreate) (*V0043ClusterRec, error) {
	apiCluster := &V0043ClusterRec{}

	// Required fields
	apiCluster.Name = &cluster.Name

	// Optional fields
	if cluster.ControlHost != "" || cluster.ControlPort != 0 {
		apiCluster.Controller = &struct {
			Host *string `json:"host,omitempty"`
			Port *int32  `json:"port,omitempty"`
		}{}
		if cluster.ControlHost != "" {
			apiCluster.Controller.Host = &cluster.ControlHost
		}
		if cluster.ControlPort != 0 {
			port := int32(cluster.ControlPort)
			apiCluster.Controller.Port = &port
		}
	}
	if cluster.RPCVersion != 0 {
		version := int32(cluster.RPCVersion)
		apiCluster.RpcVersion = &version
	}
	if len(cluster.Features) > 0 {
		flags := make([]V0043ClusterRecFlags, len(cluster.Features))
		for i, feature := range cluster.Features {
			flags[i] = V0043ClusterRecFlags(feature)
		}
		apiCluster.Flags = &flags
	}

	return apiCluster, nil
}

// convertInterfaceClusterUpdateToAPI converts interface cluster update to API type
func convertInterfaceClusterUpdateToAPI(clusterName string, update *interfaces.ClusterUpdate) (*V0043ClusterRec, error) {
	apiCluster := &V0043ClusterRec{}

	// Set the cluster name
	apiCluster.Name = &clusterName

	// Apply optional updates
	if update.ControlHost != nil || update.ControlPort != nil {
		apiCluster.Controller = &struct {
			Host *string `json:"host,omitempty"`
			Port *int32  `json:"port,omitempty"`
		}{}
		if update.ControlHost != nil {
			apiCluster.Controller.Host = update.ControlHost
		}
		if update.ControlPort != nil {
			port := int32(*update.ControlPort)
			apiCluster.Controller.Port = &port
		}
	}
	if update.RPCVersion != nil {
		version := int32(*update.RPCVersion)
		apiCluster.RpcVersion = &version
	}
	if len(update.Features) > 0 {
		flags := make([]V0043ClusterRecFlags, len(update.Features))
		for i, feature := range update.Features {
			flags[i] = V0043ClusterRecFlags(feature)
		}
		apiCluster.Flags = &flags
	}

	return apiCluster, nil
}

// applyClusterFiltering applies client-side filtering to clusters
func applyClusterFiltering(clusters []*interfaces.Cluster, opts *interfaces.ListClustersOptions) []*interfaces.Cluster {
	if opts == nil {
		return clusters
	}

	var filtered []*interfaces.Cluster

	for _, cluster := range clusters {
		// Apply name filtering
		if len(opts.Names) > 0 {
			nameMatch := false
			for _, name := range opts.Names {
				if cluster.Name == name {
					nameMatch = true
					break
				}
			}
			if !nameMatch {
				continue
			}
		}

		// Apply federation state filtering
		if len(opts.FederationStates) > 0 {
			stateMatch := false
			for _, state := range opts.FederationStates {
				if cluster.FederationState == state {
					stateMatch = true
					break
				}
			}
			if !stateMatch {
				continue
			}
		}

		// Apply feature filtering
		if len(opts.Features) > 0 {
			featureMatch := false
			for _, requiredFeature := range opts.Features {
				for _, clusterFeature := range cluster.Features {
					if clusterFeature == requiredFeature {
						featureMatch = true
						break
					}
				}
				if featureMatch {
					break
				}
			}
			if !featureMatch {
				continue
			}
		}

		// Apply control host filtering
		if len(opts.ControlHosts) > 0 {
			hostMatch := false
			for _, host := range opts.ControlHosts {
				if cluster.ControlHost == host {
					hostMatch = true
					break
				}
			}
			if !hostMatch {
				continue
			}
		}

		filtered = append(filtered, cluster)
	}

	return filtered
}
