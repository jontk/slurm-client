// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
	"github.com/jontk/slurm-client/pkg/watch"
)

// NodeManagerImpl provides the actual implementation for NodeManager methods
type NodeManagerImpl struct {
	client *WrapperClient
}

// NewNodeManagerImpl creates a new NodeManager implementation
func NewNodeManagerImpl(client *WrapperClient) *NodeManagerImpl {
	return &NodeManagerImpl{client: client}
}

// List nodes with optional filtering
func (m *NodeManagerImpl) List(ctx context.Context, opts *interfaces.ListNodesOptions) (*interfaces.NodeList, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmV0042GetNodesParams{}

	// Set flags to get detailed node information
	flags := SlurmV0042GetNodesParamsFlagsDETAIL
	params.Flags = &flags

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0042GetNodesWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.42")
	}

	// Check HTTP status and handle API errors
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.42", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.42")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response but got nil")
	}

	if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		return nil, fmt.Errorf("API error: %v", (*resp.JSON200.Errors)[0])
	}

	// Convert the response to our interface types
	nodes := make([]interfaces.Node, 0, len(resp.JSON200.Nodes))
	for _, apiNode := range resp.JSON200.Nodes {
		node, err := convertAPINodeToInterface(apiNode)
		if err != nil {
			return nil, fmt.Errorf("failed to convert node data: %w", err)
		}
		nodes = append(nodes, *node)
	}

	// Apply client-side filtering if options are provided
	if opts != nil {
		nodes = filterNodes(nodes, opts)
	}

	return &interfaces.NodeList{
		Nodes: nodes,
		Total: len(nodes),
	}, nil
}

// convertAPINodeToInterface converts a V0042Node to interfaces.Node
func convertAPINodeToInterface(apiNode V0042Node) (*interfaces.Node, error) {
	node := &interfaces.Node{}

	// Node name - simple string
	if apiNode.Name != nil {
		node.Name = *apiNode.Name
	}

	// Node state - array of strings, take first one
	if apiNode.State != nil && len(*apiNode.State) > 0 {
		node.State = (*apiNode.State)[0]
	}

	// CPUs - simple int32 pointer
	if apiNode.Cpus != nil {
		node.CPUs = int(*apiNode.Cpus)
	}

	// Memory (convert MB to bytes for consistency) - simple int64 pointer
	if apiNode.RealMemory != nil {
		node.Memory = int(*apiNode.RealMemory) * 1024 * 1024
	}

	// Partitions - V0042CsvString (which is []string)
	if apiNode.Partitions != nil {
		node.Partitions = *apiNode.Partitions
	} else {
		node.Partitions = []string{}
	}

	// Features - V0042CsvString (which is []string)
	if apiNode.Features != nil {
		node.Features = *apiNode.Features
	} else {
		node.Features = []string{}
	}

	// Reason - simple string
	if apiNode.Reason != nil {
		node.Reason = *apiNode.Reason
	}

	// Last busy time - V0042Uint64NoValStruct
	if apiNode.LastBusy != nil && apiNode.LastBusy.Set != nil && *apiNode.LastBusy.Set &&
		apiNode.LastBusy.Number != nil && *apiNode.LastBusy.Number > 0 {
		lastBusy := time.Unix(int64(*apiNode.LastBusy.Number), 0)
		node.LastBusy = &lastBusy
	}

	// Architecture - simple string
	if apiNode.Architecture != nil {
		node.Architecture = *apiNode.Architecture
	}

	// CPU Load - convert int32 to float64
	if apiNode.CpuLoad != nil {
		node.CPULoad = float64(*apiNode.CpuLoad)
	}

	// Allocated CPUs
	if apiNode.AllocCpus != nil {
		node.AllocCPUs = *apiNode.AllocCpus
	}

	// Allocated Memory (already in MB from API)
	if apiNode.AllocMemory != nil {
		node.AllocMemory = *apiNode.AllocMemory
	}

	// Free Memory - V0042Uint64NoValStruct (in MB from API)
	if apiNode.FreeMem != nil && apiNode.FreeMem.Set != nil && *apiNode.FreeMem.Set &&
		apiNode.FreeMem.Number != nil {
		node.FreeMemory = int64(*apiNode.FreeMem.Number)
	}

	// Initialize metadata
	node.Metadata = make(map[string]interface{})

	// Add additional metadata from API response
	if apiNode.BootTime != nil && apiNode.BootTime.Set != nil && *apiNode.BootTime.Set && apiNode.BootTime.Number != nil {
		node.Metadata["boot_time"] = time.Unix(int64(*apiNode.BootTime.Number), 0)
	}
	if apiNode.Boards != nil {
		node.Metadata["boards"] = int(*apiNode.Boards)
	}
	if apiNode.Cores != nil {
		node.Metadata["cores_per_socket"] = int(*apiNode.Cores)
	}
	if apiNode.Sockets != nil {
		node.Metadata["sockets"] = int(*apiNode.Sockets)
	}
	if apiNode.Threads != nil {
		node.Metadata["threads_per_core"] = int(*apiNode.Threads)
	}
	if apiNode.OperatingSystem != nil {
		node.Metadata["operating_system"] = *apiNode.OperatingSystem
	}
	if apiNode.Hostname != nil {
		node.Metadata["hostname"] = *apiNode.Hostname
	}
	if apiNode.Version != nil {
		node.Metadata["version"] = *apiNode.Version
	}
	if apiNode.Gres != nil {
		node.Metadata["gres"] = *apiNode.Gres
	}
	if apiNode.GresDrained != nil {
		node.Metadata["gres_drained"] = *apiNode.GresDrained
	}
	if apiNode.GresUsed != nil {
		node.Metadata["gres_used"] = *apiNode.GresUsed
	}

	return node, nil
}

// filterNodes applies client-side filtering to node list
func filterNodes(nodes []interfaces.Node, opts *interfaces.ListNodesOptions) []interfaces.Node {
	var filtered []interfaces.Node

	// If no options provided, return all nodes
	if opts == nil {
		return nodes
	}

	for _, node := range nodes {
		// Filter by states
		if len(opts.States) > 0 {
			stateMatch := false
			for _, state := range opts.States {
				if strings.EqualFold(node.State, state) {
					stateMatch = true
					break
				}
			}
			if !stateMatch {
				continue
			}
		}

		// Filter by partition
		if opts.Partition != "" {
			partitionMatch := false
			for _, partition := range node.Partitions {
				if strings.EqualFold(partition, opts.Partition) {
					partitionMatch = true
					break
				}
			}
			if !partitionMatch {
				continue
			}
		}

		// Filter by features
		if len(opts.Features) > 0 {
			featureMatch := true
			for _, requiredFeature := range opts.Features {
				hasFeature := false
				for _, nodeFeature := range node.Features {
					if strings.EqualFold(nodeFeature, requiredFeature) {
						hasFeature = true
						break
					}
				}
				if !hasFeature {
					featureMatch = false
					break
				}
			}
			if !featureMatch {
				continue
			}
		}

		filtered = append(filtered, node)
	}

	// Apply limit and offset
	if opts.Offset > 0 {
		if opts.Offset >= len(filtered) {
			return []interfaces.Node{}
		}
		filtered = filtered[opts.Offset:]
	}

	if opts.Limit > 0 && len(filtered) > opts.Limit {
		filtered = filtered[:opts.Limit]
	}

	return filtered
}

// Get retrieves a specific node by name
func (m *NodeManagerImpl) Get(ctx context.Context, nodeName string) (*interfaces.Node, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmV0042GetNodeParams{}

	// Set flags to get detailed node information
	flags := SlurmV0042GetNodeParamsFlagsDETAIL
	params.Flags = &flags

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0042GetNodeWithResponse(ctx, nodeName, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.42")
	}

	// Check HTTP status and handle API errors
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.42", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.42")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response but got nil")
	}

	// Convert the response to our interface types
	if len(resp.JSON200.Nodes) == 0 {
		return nil, errors.NewNodeError([]string{nodeName}, "get", fmt.Errorf("node not found"))
	}

	if len(resp.JSON200.Nodes) > 1 {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected multiple nodes returned", fmt.Sprintf("Expected 1 node but got %d for name %s", len(resp.JSON200.Nodes), nodeName))
	}

	node, err := convertAPINodeToInterface(resp.JSON200.Nodes[0])
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert node data")
		conversionErr.Cause = err
		conversionErr.Details = fmt.Sprintf("Error converting node %s", nodeName)
		return nil, conversionErr
	}

	return node, nil
}

// Update updates node properties
func (m *NodeManagerImpl) Update(ctx context.Context, nodeName string, update *interfaces.NodeUpdate) error {
	// Check if API client is available
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Validate inputs
	if update == nil {
		return errors.NewClientError(errors.ErrorCodeInvalidRequest, "Node update cannot be nil")
	}

	// Convert interface NodeUpdate to API UpdateNodeMsg
	nodeUpdate, err := convertNodeUpdateToAPI(update)
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeInvalidRequest, "Failed to convert node update")
		conversionErr.Cause = err
		conversionErr.Details = "Error converting NodeUpdate to API format"
		return conversionErr
	}

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0042PostNodeWithResponse(ctx, nodeName, *nodeUpdate)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.42")
	}

	// Check HTTP status and handle API errors
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.42", apiErrors)
				return apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.42")
		return httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil {
		return errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response but got nil")
	}

	return nil
}

// convertNodeUpdateToAPI converts interfaces.NodeUpdate to V0042UpdateNodeMsg
func convertNodeUpdateToAPI(update *interfaces.NodeUpdate) (*V0042UpdateNodeMsg, error) {
	if update == nil {
		return nil, errors.NewClientError(errors.ErrorCodeInvalidRequest, "Node update cannot be nil")
	}

	nodeUpdate := &V0042UpdateNodeMsg{}

	// Set node state if provided
	if update.State != nil {
		// Convert state string to V0042NodeStates array
		state := []string{*update.State}
		nodeUpdate.State = &state
	}

	// Set reason if provided
	if update.Reason != nil {
		nodeUpdate.Reason = update.Reason
	}

	// Set features if provided
	if len(update.Features) > 0 {
		features := V0042CsvString(update.Features)
		nodeUpdate.Features = &features
	}

	return nodeUpdate, nil
}

// Watch provides real-time node updates through polling
// Note: v0.0.42 API does not support native streaming/WebSocket node monitoring
func (m *NodeManagerImpl) Watch(ctx context.Context, opts *interfaces.WatchNodesOptions) (<-chan interfaces.NodeEvent, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Create a node poller with the List function
	poller := watch.NewNodePoller(m.List)

	// Configure polling interval if needed (default is 5 seconds)
	// poller.WithPollInterval(3 * time.Second)

	// Start watching
	return poller.Watch(ctx, opts)
}

// Delete deletes a node
func (m *NodeManagerImpl) Delete(ctx context.Context, nodeName string) error {
	return errors.NewNotImplementedError("Delete", "v0.0.42")
}

// Drain drains a node, preventing new jobs from being scheduled on it
func (m *NodeManagerImpl) Drain(ctx context.Context, nodeName string, reason string) error {
	// v0.0.42 supports drain operations, but implement using Update method like v0.0.43
	drainState := "DRAIN"
	update := &interfaces.NodeUpdate{
		State:  &drainState,
		Reason: &reason,
	}
	return m.Update(ctx, nodeName, update)
}

// Resume resumes a drained node, allowing new jobs to be scheduled on it
func (m *NodeManagerImpl) Resume(ctx context.Context, nodeName string) error {
	// v0.0.42 supports resume operations, but implement using Update method like v0.0.43
	resumeState := "RESUME"
	update := &interfaces.NodeUpdate{
		State: &resumeState,
	}
	return m.Update(ctx, nodeName, update)
}
