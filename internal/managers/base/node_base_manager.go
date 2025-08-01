// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/pkg/errors"
)

// NodeBaseManager provides common node management functionality
type NodeBaseManager struct {
	*CRUDManager
}

// NewNodeBaseManager creates a new node base manager
func NewNodeBaseManager(version string) *NodeBaseManager {
	return &NodeBaseManager{
		CRUDManager: NewCRUDManager(version, "Node"),
	}
}

// ValidateNodeUpdate validates node update data
func (m *NodeBaseManager) ValidateNodeUpdate(update *types.NodeUpdate) error {
	if update == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Update data is required",
			"update", update, nil,
		)
	}

	// Validate numeric fields if provided
	if update.CPUBinding != nil {
		if err := m.ValidateNonNegative(int(*update.CPUBinding), "update.CPUBinding"); err != nil {
			return err
		}
	}
	if update.Weight != nil {
		if err := m.ValidateNonNegative(int(*update.Weight), "update.Weight"); err != nil {
			return err
		}
	}

	// Validate state if provided
	if update.State != nil {
		if err := m.ValidateNodeState(*update.State); err != nil {
			return err
		}
	}

	// Validate next state after reboot if provided
	if update.NextStateAfterReboot != nil {
		if err := m.ValidateNodeState(*update.NextStateAfterReboot); err != nil {
			return err
		}
	}

	return nil
}

// ValidateNodeState validates node state
func (m *NodeBaseManager) ValidateNodeState(state types.NodeState) error {
	validStates := []types.NodeState{
		types.NodeStateUnknown,
		types.NodeStateDown,
		types.NodeStateIdle,
		types.NodeStateAllocated,
		types.NodeStateError,
		types.NodeStateMixed,
		types.NodeStateFuture,
		types.NodeStateReserved,
		types.NodeStateUndrained,
		types.NodeStateCloud,
		types.NodeStateDraining,
		types.NodeStateDrained,
		types.NodeStateResuming,
		types.NodeStateFail,
		types.NodeStateFailing,
		types.NodeStateMaintenance,
		types.NodeStateRebooting,
		types.NodeStateCancelling,
		types.NodeStatePoweredDown,
		types.NodeStatePoweringDown,
		types.NodeStatePoweringUp,
		types.NodeStatePlanned,
	}

	for _, validState := range validStates {
		if state == validState {
			return nil
		}
	}

	return errors.NewValidationError(
		errors.ErrorCodeValidationFailed,
		fmt.Sprintf("Invalid node state: %s", state),
		"state", state, nil,
	)
}

// ValidateNodeMaintenanceRequest validates node maintenance request
func (m *NodeBaseManager) ValidateNodeMaintenanceRequest(request *types.NodeMaintenanceRequest) error {
	if request == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Maintenance request is required",
			"request", request, nil,
		)
	}

	if len(request.Nodes) == 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"At least one node must be specified",
			"request.Nodes", request.Nodes, nil,
		)
	}

	if request.Reason == "" {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Reason is required for maintenance request",
			"request.Reason", request.Reason, nil,
		)
	}

	// Validate time constraints
	if request.StartTime != nil && request.EndTime != nil {
		if request.EndTime.Before(*request.StartTime) {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				"End time cannot be before start time",
				"request.EndTime", request.EndTime, nil,
			)
		}
	}

	if request.FixedDuration < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Fixed duration must be non-negative",
			"request.FixedDuration", request.FixedDuration, nil,
		)
	}

	return nil
}

// ValidateNodePowerRequest validates node power request
func (m *NodeBaseManager) ValidateNodePowerRequest(request *types.NodePowerRequest) error {
	if request == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Power request is required",
			"request", request, nil,
		)
	}

	if len(request.Nodes) == 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"At least one node must be specified",
			"request.Nodes", request.Nodes, nil,
		)
	}

	// Validate power state
	validPowerStates := []types.NodePowerState{
		types.NodePowerDown,
		types.NodePowerUp,
		types.NodePowerSave,
	}

	found := false
	for _, validState := range validPowerStates {
		if request.PowerState == validState {
			found = true
			break
		}
	}

	if !found {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			fmt.Sprintf("Invalid power state: %s", request.PowerState),
			"request.PowerState", request.PowerState, nil,
		)
	}

	return nil
}

// ApplyNodeUpdateDefaults applies default values to node update request
func (m *NodeBaseManager) ApplyNodeUpdateDefaults(update *types.NodeUpdate) *types.NodeUpdate {
	// Ensure slice fields are initialized
	if update.Features == nil {
		update.Features = []string{}
	}
	if update.ActiveFeatures == nil {
		update.ActiveFeatures = []string{}
	}

	// Initialize maps
	if update.Extra == nil {
		update.Extra = make(map[string]string)
	}

	return update
}

// FilterNodeList applies filtering to a node list
func (m *NodeBaseManager) FilterNodeList(items []types.Node, opts *types.NodeListOptions) []types.Node {
	if opts == nil {
		return items
	}

	filtered := make([]types.Node, 0, len(items))
	for _, node := range items {
		if m.matchesNodeFilters(node, opts) {
			filtered = append(filtered, node)
		}
	}

	return filtered
}

// matchesNodeFilters checks if a node matches the given filters
func (m *NodeBaseManager) matchesNodeFilters(node types.Node, opts *types.NodeListOptions) bool {
	// Filter by names
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

	// Filter by states
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
		for _, filterPartition := range opts.Partitions {
			for _, nodePartition := range node.Partitions {
				if strings.EqualFold(nodePartition, filterPartition) {
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

	// Filter by update time
	if opts.UpdateTime != nil {
		// This would require API support to track update times
		// For now, we'll accept all items
	}

	return true
}
