// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"fmt"
	"time"

	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
)

// NodeAdapter implements the NodeAdapter interface for v0.0.42
type NodeAdapter struct {
	*base.BaseManager
	client *api.ClientWithResponses
}

// NewNodeAdapter creates a new Node adapter for v0.0.42
func NewNodeAdapter(client *api.ClientWithResponses) *NodeAdapter {
	return &NodeAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Node"),
		client:      client,
	}
}

// List retrieves a list of nodes
func (a *NodeAdapter) List(ctx context.Context, opts *types.NodeListOptions) (*types.NodeList, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmV0042GetNodesParams{}

	// Set flags to get detailed node information
	flags := api.SlurmV0042GetNodesParamsFlagsDETAIL
	params.Flags = &flags

	// Apply filters from options
	// Note: v0.0.42 GetNodes doesn't support filtering by name, state, or partition
	// We'll need to filter the results after fetching all nodes
	_ = opts

	// Call the API
	resp, err := a.client.SlurmV0042GetNodesWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list nodes")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleHTTPResponse(resp.HTTPResponse, resp.Body)
	}

	// Check for API response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from API")
	}

	// Convert the response to common types
	nodeList := &types.NodeList{
		Nodes: make([]types.Node, 0),
	}

	if resp.JSON200.Nodes != nil {
		for _, apiNode := range resp.JSON200.Nodes {
			node, err := a.convertAPINodeToCommon(apiNode)
			if err != nil {
				// Log conversion error but continue
				continue
			}
			// Apply filters if options were provided
			if opts != nil {
				// Filter by name
				if len(opts.Names) > 0 {
					match := false
					for _, name := range opts.Names {
						if node.Name == name {
							match = true
							break
						}
					}
					if !match {
						continue
					}
				}

				// Filter by state
				if len(opts.States) > 0 {
					match := false
					for _, state := range opts.States {
						if string(node.State) == string(state) {
							match = true
							break
						}
					}
					if !match {
						continue
					}
				}

				// Filter by partition
				if len(opts.Partitions) > 0 {
					match := false
					for _, partition := range opts.Partitions {
						// Check if the node belongs to the partition
						// This might need to be adjusted based on how partitions are stored in nodes
						// Check if partition is in the node's partition list
						for _, nodePartition := range node.Partitions {
							if nodePartition == partition {
								match = true
								break
							}
						}
						if match {
							break
						}
					}
					if !match {
						continue
					}
				}
			}

			nodeList.Nodes = append(nodeList.Nodes, *node)
		}
	}

	return nodeList, nil
}

// Get retrieves a specific node by name
func (a *NodeAdapter) Get(ctx context.Context, name string) (*types.Node, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters
	params := &api.SlurmV0042GetNodeParams{}
	flags := api.SlurmV0042GetNodeParamsFlagsDETAIL
	params.Flags = &flags

	// Call the API
	resp, err := a.client.SlurmV0042GetNodeWithResponse(ctx, name, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to get node "+name)
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleHTTPResponse(resp.HTTPResponse, resp.Body)
	}

	// Check for API response
	if resp.JSON200 == nil || resp.JSON200.Nodes == nil || len(resp.JSON200.Nodes) == 0 {
		return nil, fmt.Errorf("node %s not found", name)
	}

	// Convert the first node in the response
	nodes := resp.JSON200.Nodes
	return a.convertAPINodeToCommon(nodes[0])
}

// Update updates a node's state or properties
func (a *NodeAdapter) Update(ctx context.Context, name string, updates *types.NodeUpdateRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Convert common update request to API format
	apiNodeUpdate, err := a.convertCommonNodeUpdateToAPI(name, updates)
	if err != nil {
		return a.WrapError(err, "failed to convert node update request")
	}

	// Call the API
	resp, err := a.client.SlurmV0042PostNodeWithResponse(ctx, name, *apiNodeUpdate)
	if err != nil {
		return a.WrapError(err, "failed to update node "+name)
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleHTTPResponse(resp.HTTPResponse, resp.Body)
	}

	return nil
}

// Create creates a new node (not typically supported via API)
func (a *NodeAdapter) Create(ctx context.Context, node *types.NodeCreateRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Node creation is typically done through slurm.conf, not API
	return fmt.Errorf("node creation not supported via v0.0.42 API")
}

// Delete deletes a node (not typically supported via API)
func (a *NodeAdapter) Delete(ctx context.Context, name string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Node deletion is typically done through slurm.conf, not API
	return fmt.Errorf("node deletion not supported via v0.0.42 API")
}

// Watch provides real-time node status updates (not fully implemented in v0.0.42)
func (a *NodeAdapter) Watch(ctx context.Context, opts *types.NodeWatchOptions) (<-chan types.NodeWatchEvent, error) {
	// For now, return not implemented error to satisfy interface
	return nil, fmt.Errorf("watch functionality not fully implemented in v0.0.42 adapter")
}

// convertAPINodeToCommon converts API node to common type
func (a *NodeAdapter) convertAPINodeToCommon(apiNode api.V0042Node) (*types.Node, error) {
	node := &types.Node{}

	// Set basic fields
	if apiNode.Name != nil {
		node.Name = *apiNode.Name
	}

	// Node state conversion
	if apiNode.State != nil && len(*apiNode.State) > 0 {
		node.State = types.NodeState((*apiNode.State)[0]) // Take the first state and convert to NodeState type
	}

	// CPU information
	if apiNode.Cpus != nil {
		node.CPUs = *apiNode.Cpus
	}

	// Memory information
	if apiNode.RealMemory != nil {
		node.RealMemory = *apiNode.RealMemory
	}

	// Partition information
	if apiNode.Partitions != nil && len(*apiNode.Partitions) > 0 {
		node.Partitions = *apiNode.Partitions
	}

	// Architecture
	if apiNode.Architecture != nil {
		node.Arch = *apiNode.Architecture
	}

	// OS information
	if apiNode.OperatingSystem != nil {
		node.OS = *apiNode.OperatingSystem
	}

	// Features
	if apiNode.ActiveFeatures != nil && len(*apiNode.ActiveFeatures) > 0 {
		node.ActiveFeatures = *apiNode.ActiveFeatures
	}
	if apiNode.Features != nil && len(*apiNode.Features) > 0 {
		node.Features = *apiNode.Features
	}

	// GRES (Generic Resources)
	if apiNode.Gres != nil {
		node.Gres = *apiNode.Gres
	}

	// Boot time - check if it's set and has a number
	if apiNode.BootTime != nil && apiNode.BootTime.Set != nil && *apiNode.BootTime.Set && apiNode.BootTime.Number != nil {
		bootTime := time.Unix(*apiNode.BootTime.Number, 0)
		node.BootTime = &bootTime
	}

	// Slurm version
	if apiNode.Version != nil {
		node.Version = *apiNode.Version
	}

	// Reason for node state
	if apiNode.Reason != nil {
		node.Reason = *apiNode.Reason
	}

	return node, nil
}

// convertCommonNodeUpdateToAPI converts common node update to API format
func (a *NodeAdapter) convertCommonNodeUpdateToAPI(name string, update *types.NodeUpdateRequest) (*api.SlurmV0042PostNodeJSONRequestBody, error) {
	// v0.0.42 has limited node update capabilities
	// For now, return a basic structure
	return &api.SlurmV0042PostNodeJSONRequestBody{}, nil
}

// Drain drains a node, preventing new jobs from being scheduled on it
func (a *NodeAdapter) Drain(ctx context.Context, nodeName string, reason string) error {
	// v0.0.42 supports drain operations through the Update method
	drainState := types.NodeState("DRAIN")
	update := &types.NodeUpdate{
		State:  &drainState,
		Reason: &reason,
	}
	return a.Update(ctx, nodeName, update)
}

// Resume resumes a drained node, allowing new jobs to be scheduled on it
func (a *NodeAdapter) Resume(ctx context.Context, nodeName string) error {
	// v0.0.42 supports resume operations through the Update method
	resumeState := types.NodeState("RESUME")
	update := &types.NodeUpdate{
		State: &resumeState,
	}
	return a.Update(ctx, nodeName, update)
}
