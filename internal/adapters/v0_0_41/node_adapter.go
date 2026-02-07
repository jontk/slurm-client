// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"context"
	"fmt"

	types "github.com/jontk/slurm-client/api"
	adapterbase "github.com/jontk/slurm-client/internal/adapters/base"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
)

// NodeAdapter implements the NodeAdapter interface for v0.0.41
type NodeAdapter struct {
	*adapterbase.BaseManager
	client *api.ClientWithResponses
}

// NewNodeAdapter creates a new Node adapter for v0.0.41
func NewNodeAdapter(client *api.ClientWithResponses) *NodeAdapter {
	return &NodeAdapter{
		BaseManager: adapterbase.NewBaseManager("v0.0.41", "Node"),
		client:      client,
	}
}

// List retrieves a list of nodes with optional filtering
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
	params := &api.SlurmV0041GetNodesParams{}
	// Apply filters from options
	// Note: v0.0.41 API has limited filtering support, implement client-side filtering
	if opts != nil {
		// v0.0.41 GetNodes doesn't support filtering parameters like NodeName or State
		// We'll filter results after fetching
		_ = opts
	}
	// Make the API call
	resp, err := a.client.SlurmV0041GetNodesWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list nodes")
	}
	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected nil response")
	}
	// Convert response to common types
	nodeList := &types.NodeList{
		Nodes: make([]types.Node, 0, len(resp.JSON200.Nodes)),
		Total: len(resp.JSON200.Nodes),
	}
	for _, apiNode := range resp.JSON200.Nodes {
		node, err := a.convertAPINodeToCommon(apiNode)
		if err != nil {
			// Log the error but continue processing other nodes
			continue
		}
		// Apply client-side filtering
		if opts != nil {
			// Filter by name
			if len(opts.Names) > 0 {
				match := false
				for _, name := range opts.Names {
					if node.Name != nil && *node.Name == name {
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
				for _, optState := range opts.States {
					for _, nodeState := range node.State {
						if nodeState == optState {
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
	// Update total count after filtering
	nodeList.Total = len(nodeList.Nodes)
	// Extract warning and error messages if any (but NodeList doesn't have Meta)
	// Warnings and errors are ignored for now as NodeList structure doesn't support them
	if resp.JSON200.Warnings != nil {
		// Log warnings if needed
		_ = resp.JSON200.Warnings
	}
	if resp.JSON200.Errors != nil {
		// Log errors if needed
		_ = resp.JSON200.Errors
	}
	return nodeList, nil
}

// Get retrieves a specific node by name
func (a *NodeAdapter) Get(ctx context.Context, name string) (*types.Node, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	// Validate name
	if err := a.ValidateResourceName("node name", name); err != nil {
		return nil, err
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}
	// Make the API call
	params := &api.SlurmV0041GetNodeParams{}
	resp, err := a.client.SlurmV0041GetNodeWithResponse(ctx, name, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to get node "+name)
	}
	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil || len(resp.JSON200.Nodes) == 0 {
		return nil, a.HandleNotFound("node " + name)
	}
	// Convert the first node in the response
	node, err := a.convertAPINodeToCommon(resp.JSON200.Nodes[0])
	if err != nil {
		return nil, a.WrapError(err, "failed to convert node "+name)
	}
	return node, nil
}

// Update updates a node's state or properties
func (a *NodeAdapter) Update(ctx context.Context, name string, update *types.NodeUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	// Validate name
	if err := a.ValidateResourceName("node name", name); err != nil {
		return err
	}
	// Validate update
	if update == nil {
		return a.HandleValidationError("node update cannot be nil")
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}
	// Convert update to API request
	updateReq := a.convertCommonToAPINodeUpdate(update)
	if updateReq == nil {
		return fmt.Errorf("failed to convert node update")
	}
	// Make the API call
	resp, err := a.client.SlurmV0041PostNodeWithResponse(ctx, name, *updateReq)
	if err != nil {
		return a.WrapError(err, "failed to update node "+name)
	}
	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
	}
	return nil
}

// SetState sets the state of a node
func (a *NodeAdapter) SetState(ctx context.Context, name string, state types.NodeState) error {
	update := &types.NodeUpdate{
		State: []types.NodeState{state},
	}
	return a.Update(ctx, name, update)
}

// Resume resumes a node
func (a *NodeAdapter) Resume(ctx context.Context, name string) error {
	state := types.NodeStateResume
	return a.SetState(ctx, name, state)
}

// Drain drains a node
func (a *NodeAdapter) Drain(ctx context.Context, name string, reason string) error {
	state := types.NodeStateDrain
	update := &types.NodeUpdate{
		State:  []types.NodeState{state},
		Reason: &reason,
	}
	return a.Update(ctx, name, update)
}

// Down marks a node as down
func (a *NodeAdapter) Down(ctx context.Context, name string, reason string) error {
	state := types.NodeStateDown
	update := &types.NodeUpdate{
		State:  []types.NodeState{state},
		Reason: &reason,
	}
	return a.Update(ctx, name, update)
}

// Delete deletes a node
func (a *NodeAdapter) Delete(ctx context.Context, name string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	// Validate name
	if err := a.ValidateResourceName("node name", name); err != nil {
		return err
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}
	// Make the API call
	resp, err := a.client.SlurmV0041DeleteNodeWithResponse(ctx, name)
	if err != nil {
		return a.WrapError(err, "failed to delete node "+name)
	}
	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
	}
	return nil
}

// Watch provides real-time node status updates (not fully implemented in v0.0.41)
func (a *NodeAdapter) Watch(ctx context.Context, opts *types.NodeWatchOptions) (<-chan types.NodeWatchEvent, error) {
	// For now, return not implemented error to satisfy interface
	return nil, fmt.Errorf("watch functionality not fully implemented in v0.0.41 adapter")
}
