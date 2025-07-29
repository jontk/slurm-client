package v0_0_41

import (
	"context"
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
)

// NodeAdapter implements the NodeAdapter interface for v0.0.41
type NodeAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewNodeAdapter creates a new Node adapter for v0.0.41
func NewNodeAdapter(client *api.ClientWithResponses) *NodeAdapter {
	return &NodeAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Node"),
		client:      client,
		wrapper:     nil, // We'll implement this later
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
	if opts != nil {
		if len(opts.Names) > 0 {
			nameStr := strings.Join(opts.Names, ",")
			params.NodeName = &nameStr
		}
		if opts.State != "" {
			params.State = &opts.State
		}
		if opts.UpdateTime != nil {
			updateTimeStr := fmt.Sprintf("%d", opts.UpdateTime.Unix())
			params.UpdateTime = &updateTimeStr
		}
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
		Meta: &types.ListMeta{
			Version: a.GetVersion(),
		},
	}

	for _, apiNode := range resp.JSON200.Nodes {
		node, err := a.convertAPINodeToCommon(apiNode)
		if err != nil {
			// Log the error but continue processing other nodes
			continue
		}
		nodeList.Nodes = append(nodeList.Nodes, *node)
	}

	// Extract warning messages if any
	if resp.JSON200.Warnings != nil {
		warnings := make([]string, 0, len(*resp.JSON200.Warnings))
		for _, warning := range *resp.JSON200.Warnings {
			if warning.Description != nil {
				warnings = append(warnings, *warning.Description)
			}
		}
		if len(warnings) > 0 {
			nodeList.Meta.Warnings = warnings
		}
	}

	// Extract error messages if any
	if resp.JSON200.Errors != nil {
		errors := make([]string, 0, len(*resp.JSON200.Errors))
		for _, error := range *resp.JSON200.Errors {
			if error.Description != nil {
				errors = append(errors, *error.Description)
			}
		}
		if len(errors) > 0 {
			nodeList.Meta.Errors = errors
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
		return nil, a.WrapError(err, fmt.Sprintf("failed to get node %s", name))
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}

	if resp.JSON200 == nil || len(resp.JSON200.Nodes) == 0 {
		return nil, a.HandleNotFound(fmt.Sprintf("node %s", name))
	}

	// Convert the first node in the response
	node, err := a.convertAPINodeToCommon(resp.JSON200.Nodes[0])
	if err != nil {
		return nil, a.WrapError(err, fmt.Sprintf("failed to convert node %s", name))
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
		return a.WrapError(err, fmt.Sprintf("failed to update node %s", name))
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
		State: &state,
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
		State:  &state,
		Reason: &reason,
	}
	return a.Update(ctx, name, update)
}

// Down marks a node as down
func (a *NodeAdapter) Down(ctx context.Context, name string, reason string) error {
	state := types.NodeStateDown
	update := &types.NodeUpdate{
		State:  &state,
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
		return a.WrapError(err, fmt.Sprintf("failed to delete node %s", name))
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
	}

	return nil
}