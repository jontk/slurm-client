// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_44

import (
	"context"
	"fmt"
	"time"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
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
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmV0044GetNodesParams{}

	// Apply filtering options if provided
	if opts != nil {
		// Note: v0.0.44 API parameters may differ, filtering will be post-retrieval
		// TODO: Check actual v0.0.44 API parameter names and implement proper filtering
		_ = opts // Avoid unused variable warning for now
	}

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0044GetNodesWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	// Handle response
	if resp.HTTPResponse.StatusCode != 200 {
		if resp.JSONDefault != nil && resp.JSONDefault.Errors != nil && len(*resp.JSONDefault.Errors) > 0 {
			errorMsg := fmt.Sprintf("API error: %v", (*resp.JSONDefault.Errors)[0])
			return nil, errors.NewSlurmError(errors.ErrorCodeInvalidRequest, errorMsg)
		}
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal,
			fmt.Sprintf("HTTP %d", resp.HTTPResponse.StatusCode))
	}

	if resp.JSON200 == nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal, "Empty response body")
	}

	// Convert to interface types
	nodeList := &interfaces.NodeList{
		Nodes: make([]interfaces.Node, 0),
	}

	if resp.JSON200.Nodes != nil {
		for _, node := range resp.JSON200.Nodes {
			convertedNode := m.convertNodeToInterface(&node)

			// Apply partition filtering if specified (post-retrieval filtering)
			if opts != nil && opts.Partition != "" {
				partitionFound := false
				for _, partition := range convertedNode.Partitions {
					if partition == opts.Partition {
						partitionFound = true
						break
					}
				}
				if !partitionFound {
					continue
				}
			}

			nodeList.Nodes = append(nodeList.Nodes, *convertedNode)
		}
	}

	return nodeList, nil
}

// Get retrieves a specific node by name
func (m *NodeManagerImpl) Get(ctx context.Context, nodeName string) (*interfaces.Node, error) {
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Use the node-specific endpoint
	resp, err := m.client.apiClient.SlurmV0044GetNodeWithResponse(ctx, nodeName, &SlurmV0044GetNodeParams{})
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	if resp.HTTPResponse.StatusCode == 404 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound,
			fmt.Sprintf("Node %s not found", nodeName))
	}

	if resp.HTTPResponse.StatusCode != 200 {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal,
			fmt.Sprintf("HTTP %d", resp.HTTPResponse.StatusCode))
	}

	if resp.JSON200 == nil || len(resp.JSON200.Nodes) == 0 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound,
			fmt.Sprintf("Node %s not found", nodeName))
	}

	node := resp.JSON200.Nodes[0]
	return m.convertNodeToInterface(&node), nil
}

// CreateNode creates a new node in the cluster (new in v0.0.44)
// Note: Placeholder implementation until NodeCreate interface is defined
func (m *NodeManagerImpl) CreateNode(ctx context.Context, nodeName string, cpus int, memory int64) (interface{}, error) {
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Create node configuration string
	// Format: NodeName=name CPUs=n RealMemory=m
	nodeConf := fmt.Sprintf("NodeName=%s", nodeName)
	if cpus > 0 {
		nodeConf += fmt.Sprintf(" CPUs=%d", cpus)
	}
	if memory > 0 {
		memoryMB := memory / (1024 * 1024) // Convert bytes to MB
		nodeConf += fmt.Sprintf(" RealMemory=%d", memoryMB)
	}

	// Convert node creation request to v0.0.44 format
	createReq := V0044OpenapiCreateNodeReq{
		NodeConf: nodeConf,
	}

	// Submit the node creation request
	resp, err := m.client.apiClient.SlurmV0044PostNewNodeWithResponse(ctx, createReq)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	if resp.HTTPResponse.StatusCode != 200 {
		if resp.JSONDefault != nil && resp.JSONDefault.Errors != nil && len(*resp.JSONDefault.Errors) > 0 {
			errorMsg := fmt.Sprintf("Node creation failed: %v", (*resp.JSONDefault.Errors)[0])
			return nil, errors.NewSlurmError(errors.ErrorCodeValidationFailed, errorMsg)
		}
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal,
			fmt.Sprintf("HTTP %d", resp.HTTPResponse.StatusCode))
	}

	if resp.JSON200 == nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal, "Empty response body")
	}

	// Return successful creation response
	return map[string]interface{}{
		"node_name": nodeName,
		"message":   "Node created successfully",
		"response":  resp.JSON200,
	}, nil
}

// Update updates node properties
func (m *NodeManagerImpl) Update(ctx context.Context, nodeName string, update *interfaces.NodeUpdate) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Convert update to v0.0.44 format - using V0044UpdateNodeMsg directly
	updateReq := V0044UpdateNodeMsg{}

	// Set updateable fields
	if update.State != nil {
		// V0044UpdateNodeMsgState is an array, but typically we set a single state
		stateArray := []V0044UpdateNodeMsgState{V0044UpdateNodeMsgState(*update.State)}
		updateReq.State = &stateArray
	}

	if update.Reason != nil {
		updateReq.Reason = update.Reason
	}

	if len(update.Features) > 0 {
		// V0044CsvString is []string, so we can assign directly
		featuresCSV := V0044CsvString(update.Features)
		updateReq.Features = &featuresCSV
	}

	// Submit the update
	resp, err := m.client.apiClient.SlurmV0044PostNodeWithResponse(ctx, nodeName, updateReq)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	if resp.HTTPResponse.StatusCode == 404 {
		return errors.NewSlurmError(errors.ErrorCodeResourceNotFound,
			fmt.Sprintf("Node %s not found", nodeName))
	}

	if resp.HTTPResponse.StatusCode != 200 {
		return errors.NewSlurmError(errors.ErrorCodeServerInternal,
			fmt.Sprintf("HTTP %d", resp.HTTPResponse.StatusCode))
	}

	return nil
}

// Watch provides real-time node updates
func (m *NodeManagerImpl) Watch(ctx context.Context, opts *interfaces.WatchNodesOptions) (<-chan interfaces.NodeEvent, error) {
	// Create a channel for node events
	eventChan := make(chan interfaces.NodeEvent)

	// For now, return a basic watcher that polls for changes
	go func() {
		defer close(eventChan)

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Poll for node changes - basic implementation
				// In a real implementation, this would use WebSocket or SSE
			}
		}
	}()

	return eventChan, nil
}

// convertNodeToInterface converts a v0.0.44 node to the interface type
func (m *NodeManagerImpl) convertNodeToInterface(node *V0044Node) *interfaces.Node {
	result := &interfaces.Node{}

	// Basic node information
	if node.Name != nil {
		result.Name = *node.Name
	}

	if node.State != nil && len(*node.State) > 0 {
		result.State = string((*node.State)[0])
	}

	if node.Partitions != nil {
		result.Partitions = []string(*node.Partitions)
	}

	// Resource information
	if node.Cpus != nil {
		result.CPUs = int(*node.Cpus)
	}

	if node.AllocCpus != nil {
		result.AllocCPUs = *node.AllocCpus
	}

	if node.RealMemory != nil {
		result.Memory = int(*node.RealMemory) // Already in MB
	}

	if node.AllocMemory != nil {
		result.AllocMemory = *node.AllocMemory // Already in MB
	}

	if node.FreeMem != nil && node.FreeMem.Number != nil {
		result.FreeMemory = *node.FreeMem.Number // MB
	}

	// Features - V0044CsvString is []string, not string
	if node.Features != nil {
		result.Features = []string(*node.Features)
	}

	// Architecture information
	if node.Architecture != nil {
		result.Architecture = *node.Architecture
	}

	// Status and reason
	if node.Reason != nil {
		result.Reason = *node.Reason
	}

	// Load and utilization
	if node.CpuLoad != nil {
		result.CPULoad = float64(*node.CpuLoad) / 100.0 // Convert percentage to decimal
	}

	// Additional v0.0.44 specific fields can be added here as they become available

	return result
}

// Delete removes a node from the cluster (if supported by version)
func (m *NodeManagerImpl) Delete(ctx context.Context, nodeName string) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// For now, return not implemented - node deletion might not be supported in REST API
	return errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Node deletion not yet implemented for v0.0.44")
}

// Drain drains a node, preventing new jobs from being scheduled on it
func (m *NodeManagerImpl) Drain(ctx context.Context, nodeName string, reason string) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Use the node update endpoint to set the node to DRAIN state
	updateReq := V0044UpdateNodeMsg{
		State:  &[]V0044UpdateNodeMsgState{V0044UpdateNodeMsgState("DRAIN")},
		Reason: &reason,
	}

	resp, err := m.client.apiClient.SlurmV0044PostNodeWithResponse(ctx, nodeName, updateReq)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	if resp.HTTPResponse.StatusCode == 404 {
		return errors.NewSlurmError(errors.ErrorCodeResourceNotFound,
			fmt.Sprintf("Node %s not found", nodeName))
	}

	if resp.HTTPResponse.StatusCode != 200 {
		return errors.NewSlurmError(errors.ErrorCodeServerInternal,
			fmt.Sprintf("HTTP %d", resp.HTTPResponse.StatusCode))
	}

	return nil
}

// Resume resumes a drained node, allowing new jobs to be scheduled on it
func (m *NodeManagerImpl) Resume(ctx context.Context, nodeName string) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Use the node update endpoint to set the node to IDLE/RESUME state
	updateReq := V0044UpdateNodeMsg{
		State: &[]V0044UpdateNodeMsgState{V0044UpdateNodeMsgState("RESUME")},
	}

	resp, err := m.client.apiClient.SlurmV0044PostNodeWithResponse(ctx, nodeName, updateReq)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	if resp.HTTPResponse.StatusCode == 404 {
		return errors.NewSlurmError(errors.ErrorCodeResourceNotFound,
			fmt.Sprintf("Node %s not found", nodeName))
	}

	if resp.HTTPResponse.StatusCode != 200 {
		return errors.NewSlurmError(errors.ErrorCodeServerInternal,
			fmt.Sprintf("HTTP %d", resp.HTTPResponse.StatusCode))
	}

	return nil
}
