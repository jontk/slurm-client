// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	"fmt"
	"strings"
	"time"

	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
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
		if opts.UpdateTime != nil {
			updateTimeStr := fmt.Sprintf("%d", opts.UpdateTime.Unix())
			params.UpdateTime = &updateTimeStr
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
	if err := a.ValidateResourceName(nodeName, "node name"); err != nil {
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

// Update updates an existing node
func (a *NodeAdapter) Update(ctx context.Context, nodeName string, update *types.NodeUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(nodeName, "node name"); err != nil {
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

	// Create request body - convert Node to UpdateNodeMsg
	reqBody := api.SlurmV0043PostNodeJSONRequestBody{
		Comment: apiNode.Comment,
	}
	if apiNode.CpuBinding != nil {
		reqBody.CpuBind = apiNode.CpuBinding
	}
	if apiNode.Gres != nil {
		reqBody.Gres = apiNode.Gres
	}

	// Call the generated OpenAPI client (POST is used for updates in SLURM API)
	resp, err := a.client.SlurmV0043PostNodeWithResponse(ctx, existingNode.Name, reqBody)
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
	if err := a.ValidateResourceName(nodeName, "node name"); err != nil {
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

// Watch watches for node state changes using polling
func (a *NodeAdapter) Watch(ctx context.Context, opts *types.NodeWatchOptions) (<-chan types.NodeWatchEvent, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Create the event channel
	eventCh := make(chan types.NodeWatchEvent, 10) // Buffered channel to prevent blocking

	// Start polling in a goroutine
	go func() {
		defer close(eventCh)

		// Poll interval - configurable, but default to 5 seconds
		pollInterval := 5 * time.Second

		// Keep track of node states to detect changes
		nodeStates := make(map[string]types.NodeState)

		// Create a ticker for polling
		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()

		// Initial poll to populate the state map
		a.pollNodes(ctx, opts, nodeStates, eventCh, true)

		// Poll for changes
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				a.pollNodes(ctx, opts, nodeStates, eventCh, false)
			}
		}
	}()

	return eventCh, nil
}

// pollNodes polls for node changes and sends events
func (a *NodeAdapter) pollNodes(ctx context.Context, opts *types.NodeWatchOptions, nodeStates map[string]types.NodeState, eventCh chan<- types.NodeWatchEvent, isInitial bool) {
	// Create list options based on watch options
	listOpts := &types.NodeListOptions{}

	// If watching specific nodes, filter by node names
	if opts != nil && len(opts.NodeNames) > 0 {
		listOpts.Names = opts.NodeNames
	}

	// If watching specific states, filter by states
	if opts != nil && len(opts.States) > 0 {
		listOpts.States = opts.States
	}

	// If watching specific partitions, filter by partitions
	if opts != nil && len(opts.Partitions) > 0 {
		listOpts.Partitions = opts.Partitions
	}

	// Get current node list
	nodeList, err := a.List(ctx, listOpts)
	if err != nil {
		// Send error event
		select {
		case eventCh <- types.NodeWatchEvent{
			EventTime: time.Now(),
			EventType: "error",
			NodeName:  "",
			Reason:    fmt.Sprintf("Failed to poll nodes: %v", err),
		}:
		case <-ctx.Done():
		}
		return
	}

	// Check for state changes
	for _, node := range nodeList.Nodes {
		previousState, exists := nodeStates[node.Name]
		currentState := node.State

		// Update the state map
		nodeStates[node.Name] = currentState

		// Skip initial population unless it's a special event
		if isInitial {
			continue
		}

		// Send event if state changed
		if exists && previousState != currentState {
			eventType := a.getEventTypeFromNodeStateChange(previousState, currentState)

			event := types.NodeWatchEvent{
				EventTime:     time.Now(),
				EventType:     eventType,
				NodeName:      node.Name,
				PreviousState: previousState,
				NewState:      currentState,
				Reason:        a.getReasonFromNodeStateChange(previousState, currentState, node.Reason),
				Partitions:    node.Partitions,
			}

			select {
			case eventCh <- event:
			case <-ctx.Done():
				return
			}
		}
	}

	// Check for removed nodes (nodes that existed before but don't exist now)
	for nodeName, previousState := range nodeStates {
		found := false
		for _, node := range nodeList.Nodes {
			if node.Name == nodeName {
				found = true
				break
			}
		}

		// If node is not found, it might have been removed
		if !found {
			// Send removal event
			event := types.NodeWatchEvent{
				EventTime:     time.Now(),
				EventType:     "removed",
				NodeName:      nodeName,
				PreviousState: previousState,
				NewState:      types.NodeStateUnknown,
				Reason:        "Node removed from cluster",
			}

			select {
			case eventCh <- event:
			case <-ctx.Done():
				return
			}

			// Remove from state map
			delete(nodeStates, nodeName)
		}
	}
}

// getEventTypeFromNodeStateChange determines the event type based on node state transition
func (a *NodeAdapter) getEventTypeFromNodeStateChange(previous, current types.NodeState) string {
	switch current {
	case types.NodeStateIdle:
		if previous == types.NodeStateAllocated || previous == types.NodeStateMixed {
			return "freed"
		}
		return "idle"
	case types.NodeStateAllocated:
		if previous == types.NodeStateIdle {
			return "allocated"
		}
		return "state_change"
	case types.NodeStateMixed:
		return "mixed"
	case types.NodeStateDown:
		return "down"
	case types.NodeStateError:
		return "error"
	case types.NodeStateDraining:
		return "draining"
	case types.NodeStateDrained:
		return "drained"
	case types.NodeStateResuming:
		return "resuming"
	case types.NodeStateFail:
		return "fail"
	case types.NodeStateMaintenance:
		return "maintenance"
	case types.NodeStateRebooting:
		return "rebooting"
	default:
		return "state_change"
	}
}

// getReasonFromNodeStateChange provides a reason for the node state change
func (a *NodeAdapter) getReasonFromNodeStateChange(previous, current types.NodeState, nodeReason string) string {
	// If the node has a specific reason, use that
	if nodeReason != "" {
		return nodeReason
	}

	// Otherwise, provide a generic reason based on state transition
	switch current {
	case types.NodeStateIdle:
		if previous == types.NodeStateAllocated || previous == types.NodeStateMixed {
			return "Node jobs completed, now idle"
		}
		return "Node is idle and available"
	case types.NodeStateAllocated:
		if previous == types.NodeStateIdle {
			return "Node allocated to jobs"
		}
		return "Node fully allocated"
	case types.NodeStateMixed:
		return "Node partially allocated"
	case types.NodeStateDown:
		return "Node is down"
	case types.NodeStateError:
		return "Node in error state"
	case types.NodeStateDraining:
		return "Node is draining jobs"
	case types.NodeStateDrained:
		return "Node has been drained"
	case types.NodeStateResuming:
		return "Node is resuming from power save"
	case types.NodeStateFail:
		return "Node has failed"
	case types.NodeStateMaintenance:
		return "Node is in maintenance mode"
	case types.NodeStateRebooting:
		return "Node is rebooting"
	default:
		return fmt.Sprintf("Node state changed from %s to %s", previous, current)
	}
}

// validateNodeUpdate validates node update request
func (a *NodeAdapter) validateNodeUpdate(update *types.NodeUpdate) error {
	if update == nil {
		return common.NewValidationError("node update data is required", "update", nil)
	}
	// Empty updates are allowed - the API will handle no-op updates

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

// Drain drains a node, preventing new jobs from being scheduled on it
func (a *NodeAdapter) Drain(ctx context.Context, nodeName string, reason string) error {
	// v0.0.43 supports drain operations through the Update method
	drainState := types.NodeState("DRAIN")
	update := &types.NodeUpdate{
		State:  &drainState,
		Reason: &reason,
	}
	return a.Update(ctx, nodeName, update)
}

// Resume resumes a drained node, allowing new jobs to be scheduled on it
func (a *NodeAdapter) Resume(ctx context.Context, nodeName string) error {
	// v0.0.43 supports resume operations through the Update method
	resumeState := types.NodeState("RESUME")
	update := &types.NodeUpdate{
		State: &resumeState,
	}
	return a.Update(ctx, nodeName, update)
}
