// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"fmt"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
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
		return nil, a.WrapError(err, fmt.Sprintf("failed to get node %s", name))
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
		return a.WrapError(err, fmt.Sprintf("failed to update node %s", name))
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
