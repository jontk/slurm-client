package v0_0_43

import (
	"context"
	"strings"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// NodeAdapter implements the NodeAdapter interface for v0.0.43
type NodeAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewNodeAdapter creates a new Node adapter for v0.0.43
func NewNodeAdapter(client *api.ClientWithResponses) *NodeAdapter {
	return &NodeAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Node"),
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
	params := &api.SlurmV0043GetNodesParams{}

	// Apply filters from options
	if opts != nil {
		if len(opts.Names) > 0 {
			nameStr := strings.Join(opts.Names, ",")
			params.NodeName = &nameStr
		}
		if len(opts.States) > 0 {
			stateStrs := make([]string, len(opts.States))
			for i, state := range opts.States {
				stateStrs[i] = string(state)
			}
			stateStr := strings.Join(stateStrs, ",")
			params.States = &stateStr
		}
		if opts.UpdateTime != nil {
			updateTime := opts.UpdateTime.Unix()
			params.UpdateTime = &updateTime
		}
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0043GetNodesWithResponse(ctx, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "List Nodes"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Nodes, "List Nodes - nodes field"); err != nil {
		return nil, err
	}

	// Convert the response to common types
	nodeList := make([]types.Node, 0, len(*resp.JSON200.Nodes))
	for _, apiNode := range *resp.JSON200.Nodes {
		node, err := a.convertAPINodeToCommon(apiNode)
		if err != nil {
			return nil, a.HandleConversionError(err, apiNode.Name)
		}
		nodeList = append(nodeList, *node)
	}

	// Apply client-side filtering if needed
	if opts != nil {
		nodeList = a.filterNodeList(nodeList, opts)
	}

	// Apply pagination
	listOpts := base.ListOptions{}
	if opts != nil {
		listOpts.Limit = opts.Limit
		listOpts.Offset = opts.Offset
	}

	// Apply pagination
	start := listOpts.Offset
	if start < 0 {
		start = 0
	}
	if start >= len(nodeList) {
		return &types.NodeList{
			Nodes: []types.Node{},
			Total: len(nodeList),
		}, nil
	}

	end := len(nodeList)
	if listOpts.Limit > 0 {
		end = start + listOpts.Limit
		if end > len(nodeList) {
			end = len(nodeList)
		}
	}

	return &types.NodeList{
		Nodes: nodeList[start:end],
		Total: len(nodeList),
	}, nil
}

// Get retrieves a specific node by name
func (a *NodeAdapter) Get(ctx context.Context, nodeName string) (*types.Node, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.ValidateResourceName(nodeName, "nodeName"); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmV0043GetNodeParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0043GetNodeWithResponse(ctx, nodeName, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "Get Node"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Nodes, "Get Node - nodes field"); err != nil {
		return nil, err
	}

	// Check if we got any node entries
	if len(*resp.JSON200.Nodes) == 0 {
		return nil, common.NewResourceNotFoundError("Node", nodeName)
	}

	// Convert the first node (should be the only one)
	node, err := a.convertAPINodeToCommon((*resp.JSON200.Nodes)[0])
	if err != nil {
		return nil, a.HandleConversionError(err, nodeName)
	}

	return node, nil
}

// Update updates an existing node
func (a *NodeAdapter) Update(ctx context.Context, nodeName string, update *types.NodeUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(nodeName, "nodeName"); err != nil {
		return err
	}
	if err := a.validateNodeUpdate(update); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// First, get the existing node to merge updates
	existingNode, err := a.Get(ctx, nodeName)
	if err != nil {
		return err
	}

	// Convert to API format and apply updates
	apiNode, err := a.convertCommonNodeUpdateToAPI(existingNode, update)
	if err != nil {
		return err
	}

	// Create request body
	reqBody := api.SlurmV0043PostNodeJSONRequestBody{
		Nodes: []api.V0043NodeInfo{*apiNode},
	}

	// Prepare parameters for the API call
	params := &api.SlurmV0043PostNodeParams{}

	// Call the generated OpenAPI client (POST is used for updates in SLURM API)
	resp, err := a.client.SlurmV0043PostNodeWithResponse(ctx, params, reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "v0.0.43")
}

// Delete deletes a node
func (a *NodeAdapter) Delete(ctx context.Context, nodeName string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(nodeName, "nodeName"); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0043DeleteNodeWithResponse(ctx, nodeName)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	// Create adapter with special handling for 204 (No Content) status
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)

	// For DELETE operations, 204 is also a success
	if resp.StatusCode() == 204 {
		return nil
	}

	return common.HandleAPIResponse(responseAdapter, "v0.0.43")
}

// validateNodeUpdate validates node update request
func (a *NodeAdapter) validateNodeUpdate(update *types.NodeUpdate) error {
	if update == nil {
		return common.NewValidationError("node update data is required", "update", nil)
	}
	// At least one field should be provided for update
	if update.Comment == nil && update.CPUBinding == nil && len(update.Features) == 0 &&
	   len(update.ActiveFeatures) == 0 && update.Gres == nil && update.NextStateAfterReboot == nil &&
	   update.Reason == nil && update.ResumeAfter == nil && update.State == nil && 
	   update.Weight == nil && len(update.Extra) == 0 {
		return common.NewValidationError("at least one field must be provided for update", "update", update)
	}
	
	// Validate numeric fields if provided
	if update.CPUBinding != nil && *update.CPUBinding < 0 {
		return common.NewValidationError("CPU binding must be non-negative", "cpuBinding", *update.CPUBinding)
	}
	if update.Weight != nil && *update.Weight < 0 {
		return common.NewValidationError("weight must be non-negative", "weight", *update.Weight)
	}
	return nil
}

// filterNodeList applies client-side filtering to node list
func (a *NodeAdapter) filterNodeList(nodes []types.Node, opts *types.NodeListOptions) []types.Node {
	if opts == nil {
		return nodes
	}

	filtered := make([]types.Node, 0, len(nodes))
	for _, node := range nodes {
		if a.matchesNodeFilters(node, opts) {
			filtered = append(filtered, node)
		}
	}

	return filtered
}

// matchesNodeFilters checks if a node matches the given filters
func (a *NodeAdapter) matchesNodeFilters(node types.Node, opts *types.NodeListOptions) bool {
	// Filter by names (already handled by API, but included for completeness)
	if len(opts.Names) > 0 {
		found := false
		for _, name := range opts.Names {
			if strings.EqualFold(node.Name, name) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by states (already handled by API, but included for completeness) 
	if len(opts.States) > 0 {
		found := false
		for _, state := range opts.States {
			if node.State == state {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by partitions
	if len(opts.Partitions) > 0 {
		found := false
		for _, partition := range opts.Partitions {
			for _, nodePartition := range node.Partitions {
				if strings.EqualFold(nodePartition, partition) {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by reasons
	if len(opts.Reasons) > 0 {
		found := false
		for _, reason := range opts.Reasons {
			if strings.Contains(strings.ToLower(node.Reason), strings.ToLower(reason)) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}