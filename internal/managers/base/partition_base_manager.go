// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/pkg/errors"
)

// PartitionBaseManager provides common partition management functionality
type PartitionBaseManager struct {
	*CRUDManager
}

// NewPartitionBaseManager creates a new partition base manager
func NewPartitionBaseManager(version string) *PartitionBaseManager {
	return &PartitionBaseManager{
		CRUDManager: NewCRUDManager(version, "Partition"),
	}
}

// ValidatePartitionCreate validates partition creation data
func (m *PartitionBaseManager) ValidatePartitionCreate(partition *types.PartitionCreate) error {
	if partition == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Partition data is required",
			"partition", partition, nil,
		)
	}

	if err := m.ValidateResourceName(partition.Name, "partition name"); err != nil {
		return err
	}

	// Validate numeric fields
	if err := m.ValidateNonNegative(int(partition.DefaultTime), "partition.DefaultTime"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(partition.MaxTime), "partition.MaxTime"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(partition.MinNodes), "partition.MinNodes"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(partition.MaxNodes), "partition.MaxNodes"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(partition.Priority), "partition.Priority"); err != nil {
		return err
	}

	// Validate memory fields
	if partition.DefaultMemPerCPU < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Default memory per CPU must be non-negative",
			"partition.DefaultMemPerCPU", partition.DefaultMemPerCPU, nil,
		)
	}
	if partition.MaxMemPerNode < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Max memory per node must be non-negative",
			"partition.MaxMemPerNode", partition.MaxMemPerNode, nil,
		)
	}

	// Validate node constraints
	if partition.MinNodes > partition.MaxNodes && partition.MaxNodes > 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Minimum nodes cannot be greater than maximum nodes",
			"partition.MinNodes", partition.MinNodes, nil,
		)
	}

	// Validate state if provided
	if partition.State != "" {
		if err := m.ValidatePartitionState(partition.State); err != nil {
			return err
		}
	}

	return nil
}

// ValidatePartitionUpdate validates partition update data
func (m *PartitionBaseManager) ValidatePartitionUpdate(update *types.PartitionUpdate) error {
	if update == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Update data is required",
			"update", update, nil,
		)
	}

	// Validate numeric fields if provided
	if update.DefaultTime != nil {
		if err := m.ValidateNonNegative(int(*update.DefaultTime), "update.DefaultTime"); err != nil {
			return err
		}
	}
	if update.MaxTime != nil {
		if err := m.ValidateNonNegative(int(*update.MaxTime), "update.MaxTime"); err != nil {
			return err
		}
	}
	if update.MinNodes != nil {
		if err := m.ValidateNonNegative(int(*update.MinNodes), "update.MinNodes"); err != nil {
			return err
		}
	}
	if update.MaxNodes != nil {
		if err := m.ValidateNonNegative(int(*update.MaxNodes), "update.MaxNodes"); err != nil {
			return err
		}
	}
	if update.Priority != nil {
		if err := m.ValidateNonNegative(int(*update.Priority), "update.Priority"); err != nil {
			return err
		}
	}

	// Validate memory fields if provided
	if update.DefaultMemPerCPU != nil && *update.DefaultMemPerCPU < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Default memory per CPU must be non-negative",
			"update.DefaultMemPerCPU", *update.DefaultMemPerCPU, nil,
		)
	}
	if update.MaxMemPerNode != nil && *update.MaxMemPerNode < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Max memory per node must be non-negative",
			"update.MaxMemPerNode", *update.MaxMemPerNode, nil,
		)
	}

	// Validate state if provided
	if update.State != nil {
		if err := m.ValidatePartitionState(*update.State); err != nil {
			return err
		}
	}

	return nil
}

// ValidatePartitionState validates partition state
func (m *PartitionBaseManager) ValidatePartitionState(state types.PartitionState) error {
	validStates := []types.PartitionState{
		types.PartitionStateUp,
		types.PartitionStateDown,
		types.PartitionStateDrain,
		types.PartitionStateInactive,
	}

	for _, validState := range validStates {
		if state == validState {
			return nil
		}
	}

	return errors.NewValidationError(
		errors.ErrorCodeValidationFailed,
		fmt.Sprintf("Invalid partition state: %s", state),
		"state", state, nil,
	)
}

// ApplyPartitionDefaults applies default values to partition create request
func (m *PartitionBaseManager) ApplyPartitionDefaults(partition *types.PartitionCreate) *types.PartitionCreate {
	// Apply default state if not set
	if partition.State == "" {
		partition.State = types.PartitionStateUp
	}

	// Ensure slice fields are initialized
	if partition.AllowAccounts == nil {
		partition.AllowAccounts = []string{}
	}
	if partition.AllowGroups == nil {
		partition.AllowGroups = []string{}
	}
	if partition.AllowQoS == nil {
		partition.AllowQoS = []string{}
	}
	if partition.DenyAccounts == nil {
		partition.DenyAccounts = []string{}
	}
	if partition.DenyQoS == nil {
		partition.DenyQoS = []string{}
	}
	if partition.PreemptMode == nil {
		partition.PreemptMode = []string{}
	}
	if partition.SelectTypeParameters == nil {
		partition.SelectTypeParameters = []string{}
	}

	// Initialize maps
	if partition.JobDefaults == nil {
		partition.JobDefaults = make(map[string]string)
	}

	return partition
}

// FilterPartitionList applies filtering to a partition list
func (m *PartitionBaseManager) FilterPartitionList(items []types.Partition, opts *types.PartitionListOptions) []types.Partition {
	if opts == nil {
		return items
	}

	filtered := make([]types.Partition, 0, len(items))
	for _, partition := range items {
		if m.matchesPartitionFilters(partition, opts) {
			filtered = append(filtered, partition)
		}
	}

	return filtered
}

// matchesPartitionFilters checks if a partition matches the given filters
func (m *PartitionBaseManager) matchesPartitionFilters(partition types.Partition, opts *types.PartitionListOptions) bool {
	// Filter by names
	if len(opts.Names) > 0 {
		found := false
		for _, name := range opts.Names {
			if strings.EqualFold(partition.Name, name) {
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
			if partition.State == state {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by update time
	// This would require API support to track update times
	// For now, we'll accept all items
	if opts.UpdateTime != nil {
	}

	return true
}
