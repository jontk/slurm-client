// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_43

import (
	"context"
	"time"

	types "github.com/jontk/slurm-client/api"
)

const defaultNodePollInterval = 5 * time.Second

// watchNodesImpl provides the real implementation for node watching using polling.
// This overrides the stub in node_helpers.gen.go - generate_adapters.go must not emit
// the stub when this file exists.
func (a *NodeAdapter) watchNodesImpl(ctx context.Context, opts *types.NodeWatchOptions) (<-chan types.NodeWatchEvent, error) {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}
	// Create event channel
	eventCh := make(chan types.NodeWatchEvent, 100)
	// Determine poll interval
	pollInterval := defaultNodePollInterval
	if opts != nil && opts.PollInterval > 0 {
		pollInterval = opts.PollInterval
	}
	// Start polling goroutine
	go a.pollNodes(ctx, opts, eventCh, pollInterval)
	return eventCh, nil
}

// pollNodes polls for node state changes and emits events
func (a *NodeAdapter) pollNodes(ctx context.Context, opts *types.NodeWatchOptions, eventCh chan<- types.NodeWatchEvent, pollInterval time.Duration) {
	defer close(eventCh)
	// Track node states - key is node name, value is primary state
	nodeStates := make(map[string]types.NodeState)
	eventCount := int32(0)
	maxEvents := int32(0)
	if opts != nil {
		maxEvents = opts.MaxEvents
	}
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	// Do initial poll
	a.pollNodesOnce(ctx, opts, eventCh, nodeStates, &eventCount, maxEvents)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if maxEvents > 0 && eventCount >= maxEvents {
				return
			}
			a.pollNodesOnce(ctx, opts, eventCh, nodeStates, &eventCount, maxEvents)
		}
	}
}

// pollNodesOnce performs a single poll and emits events for state changes
func (a *NodeAdapter) pollNodesOnce(
	ctx context.Context,
	opts *types.NodeWatchOptions,
	eventCh chan<- types.NodeWatchEvent,
	nodeStates map[string]types.NodeState,
	eventCount *int32,
	maxEvents int32,
) {
	// Build list options
	listOpts := &types.NodeListOptions{}
	// If watching specific nodes, we could filter later
	// SLURM API may not support filtering by node name in list
	// List all nodes
	result, err := a.List(ctx, listOpts)
	if err != nil {
		return
	}
	// Track which nodes we've seen this poll
	seenNodes := make(map[string]bool)
	for i := range result.Nodes {
		node := &result.Nodes[i]
		if node.Name == nil {
			continue
		}
		nodeName := *node.Name
		// Apply node name filter
		if opts != nil && len(opts.NodeNames) > 0 && !containsString(opts.NodeNames, nodeName) {
			continue
		}
		seenNodes[nodeName] = true
		a.processNodeState(ctx, node, opts, eventCh, nodeStates, eventCount, maxEvents)
	}
	// Check for nodes that disappeared (unlikely but handle it)
	for nodeName, prevState := range nodeStates {
		if !seenNodes[nodeName] {
			event := types.NodeWatchEvent{
				EventTime:     time.Now(),
				EventType:     "removed",
				NodeName:      nodeName,
				PreviousState: prevState,
				NewState:      types.NodeStateUnknown,
			}
			select {
			case eventCh <- event:
				*eventCount++
				delete(nodeStates, nodeName)
			case <-ctx.Done():
				return
			}
			if maxEvents > 0 && *eventCount >= maxEvents {
				return
			}
		}
	}
}

// processNodeState checks for state changes and emits events
func (a *NodeAdapter) processNodeState(
	ctx context.Context,
	node *types.Node,
	opts *types.NodeWatchOptions,
	eventCh chan<- types.NodeWatchEvent,
	nodeStates map[string]types.NodeState,
	eventCount *int32,
	maxEvents int32,
) {
	if node == nil || node.Name == nil {
		return
	}
	nodeName := *node.Name
	// Get current state - State is a slice, take the first element
	var currentState types.NodeState
	if len(node.State) > 0 {
		currentState = node.State[0]
	}
	// Apply state filter
	if opts != nil && len(opts.States) > 0 && !containsNodeState(opts.States, currentState) {
		return
	}
	// Check if state changed
	prevState, exists := nodeStates[nodeName]
	if !exists {
		// New node - emit "discovered" event
		event := types.NodeWatchEvent{
			EventTime: time.Now(),
			EventType: "discovered",
			NodeName:  nodeName,
			NewState:  currentState,
		}
		if len(node.Partitions) > 0 {
			event.Partitions = node.Partitions
		}
		select {
		case eventCh <- event:
			*eventCount++
		case <-ctx.Done():
			return
		}
		nodeStates[nodeName] = currentState
	} else if currentState != prevState {
		// State changed - emit event
		eventType := determineNodeEventType(prevState, currentState)
		event := types.NodeWatchEvent{
			EventTime:     time.Now(),
			EventType:     eventType,
			NodeName:      nodeName,
			PreviousState: prevState,
			NewState:      currentState,
		}
		if len(node.Partitions) > 0 {
			event.Partitions = node.Partitions
		}
		if node.Reason != nil {
			event.Reason = *node.Reason
		}
		select {
		case eventCh <- event:
			*eventCount++
		case <-ctx.Done():
			return
		}
		nodeStates[nodeName] = currentState
	}
}

// containsString checks if a string slice contains a value
func containsString(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

// containsNodeState checks if a NodeState slice contains a value
func containsNodeState(slice []types.NodeState, val types.NodeState) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

// determineNodeEventType determines the event type based on state transition
func determineNodeEventType(prev, curr types.NodeState) string {
	switch curr {
	case types.NodeStateDrain:
		return "drain"
	case types.NodeStateIdle:
		if prev == types.NodeStateDrain || prev == types.NodeStateDown {
			return "resume"
		}
		return "idle"
	case types.NodeStateAllocated:
		return "allocated"
	case types.NodeStateDown:
		return "down"
	case types.NodeStateMixed:
		return "mixed"
	case types.NodeStateError:
		return "error"
	default:
		return "state_change"
	}
}
