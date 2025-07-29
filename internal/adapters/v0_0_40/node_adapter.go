package v0_0_40

import (
	"context"
	"strings"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
)

// NodeAdapter implements the NodeAdapter interface for v0.0.40
type NodeAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewNodeAdapter creates a new Node adapter for v0.0.40
func NewNodeAdapter(client *api.ClientWithResponses) *NodeAdapter {
	return &NodeAdapter{
		BaseManager: base.NewBaseManager("v0.0.40", "Node"),
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
	params := &api.SlurmV0040GetNodesParams{}

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
	resp, err := a.client.SlurmV0040GetNodesWithResponse(ctx, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.40"); err != nil {
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
	nodeList := make([]types.Node, 0, len(resp.JSON200.Nodes))
	for _, apiNode := range resp.JSON200.Nodes {
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
	params := &api.SlurmV0040GetNodeParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0040GetNodeWithResponse(ctx, nodeName, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.40"); err != nil {
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
	if len(resp.JSON200.Nodes) == 0 {
		return nil, common.NewResourceNotFoundError("Node", nodeName)
	}

	// Convert the first node (should be the only one)
	node, err := a.convertAPINodeToCommon(resp.JSON200.Nodes[0])
	if err != nil {
		return nil, a.HandleConversionError(err, nodeName)
	}

	return node, nil
}

// Update updates a node state or properties
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

	// Convert to API format
	apiNode, err := a.convertCommonNodeUpdateToAPI(nodeName, update)
	if err != nil {
		return err
	}

	// Create request body
	reqBody := api.SlurmV0040PostNodeJSONRequestBody{
		Nodes: &[]api.V0040Node{*apiNode},
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0040PostNodeWithResponse(ctx, nodeName, reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "v0.0.40")
}

// Delete removes a node from the cluster
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
	resp, err := a.client.SlurmV0040DeleteNodeWithResponse(ctx, nodeName)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	// Create adapter with special handling for 204 (No Content) status
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)

	// For DELETE operations, 204 is also a success
	if resp.StatusCode() == 204 {
		return nil
	}

	return common.HandleAPIResponse(responseAdapter, "v0.0.40")
}

// filterNodeList applies client-side filtering to the node list
func (a *NodeAdapter) filterNodeList(nodes []types.Node, opts *types.NodeListOptions) []types.Node {
	filtered := make([]types.Node, 0, len(nodes))
	
	for _, node := range nodes {
		// Apply Partition filter
		if len(opts.Partitions) > 0 {
			found := false
			for _, partition := range opts.Partitions {
				for _, nodePartition := range node.Partitions {
					if partition == nodePartition {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				continue
			}
		}

		// Apply Features filter
		if len(opts.Features) > 0 {
			hasAllFeatures := true
			for _, feature := range opts.Features {
				found := false
				nodeFeatures := strings.Split(node.Features, ",")
				for _, nodeFeature := range nodeFeatures {
					if feature == strings.TrimSpace(nodeFeature) {
						found = true
						break
					}
				}
				if !found {
					hasAllFeatures = false
					break
				}
			}
			if !hasAllFeatures {
				continue
			}
		}

		// Apply GRES filter
		if opts.GRES != "" {
			if !strings.Contains(node.GRES, opts.GRES) {
				continue
			}
		}

		// Apply CPUs filter
		if opts.MinCPUs > 0 && node.CPUs < opts.MinCPUs {
			continue
		}

		// Apply Memory filter
		if opts.MinMemory > 0 && (node.RealMemory == nil || *node.RealMemory < opts.MinMemory) {
			continue
		}

		// Apply TmpDisk filter
		if opts.MinTmpDisk > 0 && (node.TmpDisk == nil || *node.TmpDisk < opts.MinTmpDisk) {
			continue
		}

		filtered = append(filtered, node)
	}

	return filtered
}

// validateNodeUpdate validates node update request
func (a *NodeAdapter) validateNodeUpdate(update *types.NodeUpdate) error {
	if update == nil {
		return common.NewValidationError("node update data is required", "update", nil)
	}
	// At least one field should be provided for update
	if update.State == nil && update.Reason == nil && update.Comment == nil && 
	   update.Features == nil && update.GRES == nil && update.Weight == nil {
		return common.NewValidationError("at least one field must be provided for update", "update", update)
	}
	return nil
}