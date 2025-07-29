package v0_0_42

import (
	"context"
	"fmt"
	"strings"

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
	if opts != nil {
		if len(opts.Names) > 0 {
			nodeStr := strings.Join(opts.Names, ",")
			params.NodeName = &nodeStr
		}
		if len(opts.States) > 0 {
			stateStr := strings.Join(opts.States, ",")
			params.State = &stateStr
		}
		if len(opts.Partitions) > 0 {
			partitionStr := strings.Join(opts.Partitions, ",")
			params.Partition = &partitionStr
		}
	}

	// Call the API
	resp, err := a.client.SlurmV0042GetNodesWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list nodes")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(resp.StatusCode(), resp.Body)
	}

	// Check for API response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from API")
	}

	// Convert the response to common types
	nodeList := &types.NodeList{
		Nodes: make([]*types.Node, 0),
	}

	if resp.JSON200.Nodes != nil {
		for _, apiNode := range *resp.JSON200.Nodes {
			node, err := a.convertAPINodeToCommon(apiNode)
			if err != nil {
				// Log conversion error but continue
				continue
			}
			nodeList.Nodes = append(nodeList.Nodes, node)
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
		return nil, a.HandleAPIError(resp.StatusCode(), resp.Body)
	}

	// Check for API response
	if resp.JSON200 == nil || resp.JSON200.Nodes == nil || len(*resp.JSON200.Nodes) == 0 {
		return nil, fmt.Errorf("node %s not found", name)
	}

	// Convert the first node in the response
	nodes := *resp.JSON200.Nodes
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
	resp, err := a.client.SlurmV0042PostNodeWithResponse(ctx, name, apiNodeUpdate)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to update node %s", name))
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleAPIError(resp.StatusCode(), resp.Body)
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