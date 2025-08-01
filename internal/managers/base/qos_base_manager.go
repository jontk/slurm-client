// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/pkg/errors"
)

// QoSBaseManager provides common QoS management functionality
type QoSBaseManager struct {
	*CRUDManager
}

// NewQoSBaseManager creates a new QoS base manager
func NewQoSBaseManager(version string) *QoSBaseManager {
	return &QoSBaseManager{
		CRUDManager: NewCRUDManager(version, "QoS"),
	}
}

// ValidateQoSCreate validates QoS creation data
func (m *QoSBaseManager) ValidateQoSCreate(qos *types.QoSCreate) error {
	if qos == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"QoS data is required",
			"qos", qos, nil,
		)
	}

	if err := m.ValidateResourceName(qos.Name, "qos.Name"); err != nil {
		return err
	}

	if err := m.ValidateNonNegative(qos.Priority, "qos.Priority"); err != nil {
		return err
	}

	// Validate usage factor
	if qos.UsageFactor < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"QoS usage factor must be non-negative",
			"qos.UsageFactor", qos.UsageFactor, nil,
		)
	}

	// Validate usage threshold
	if qos.UsageThreshold < 0 || qos.UsageThreshold > 1 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"QoS usage threshold must be between 0 and 1",
			"qos.UsageThreshold", qos.UsageThreshold, nil,
		)
	}

	// Validate limits if provided
	if qos.Limits != nil {
		if err := m.ValidateQoSLimits(qos.Limits); err != nil {
			return err
		}
	}

	return nil
}

// ValidateQoSUpdate validates QoS update data
func (m *QoSBaseManager) ValidateQoSUpdate(update *types.QoSUpdate) error {
	if update == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Update data is required",
			"update", update, nil,
		)
	}

	// Validate priority if provided
	if update.Priority != nil {
		if err := m.ValidateNonNegative(*update.Priority, "update.Priority"); err != nil {
			return err
		}
	}

	// Validate usage factor if provided
	if update.UsageFactor != nil && *update.UsageFactor < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"QoS usage factor must be non-negative",
			"update.UsageFactor", *update.UsageFactor, nil,
		)
	}

	// Validate usage threshold if provided
	if update.UsageThreshold != nil && (*update.UsageThreshold < 0 || *update.UsageThreshold > 1) {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"QoS usage threshold must be between 0 and 1",
			"update.UsageThreshold", *update.UsageThreshold, nil,
		)
	}

	// Validate limits if provided
	if update.Limits != nil {
		if err := m.ValidateQoSLimits(update.Limits); err != nil {
			return err
		}
	}

	return nil
}

// ValidateQoSLimits validates QoS resource limits
func (m *QoSBaseManager) ValidateQoSLimits(limits *types.QoSLimits) error {
	// Validate per-user limits
	if err := m.validateOptionalNonNegative(limits.MaxCPUsPerUser, "limits.MaxCPUsPerUser"); err != nil {
		return err
	}
	if err := m.validateOptionalNonNegative(limits.MaxJobsPerUser, "limits.MaxJobsPerUser"); err != nil {
		return err
	}
	if err := m.validateOptionalNonNegative(limits.MaxNodesPerUser, "limits.MaxNodesPerUser"); err != nil {
		return err
	}

	// Validate per-account limits
	if err := m.validateOptionalNonNegative(limits.MaxCPUsPerAccount, "limits.MaxCPUsPerAccount"); err != nil {
		return err
	}
	if err := m.validateOptionalNonNegative(limits.MaxJobsPerAccount, "limits.MaxJobsPerAccount"); err != nil {
		return err
	}
	if err := m.validateOptionalNonNegative(limits.MaxNodesPerAccount, "limits.MaxNodesPerAccount"); err != nil {
		return err
	}

	// Validate per-job limits
	if err := m.validateOptionalNonNegative(limits.MaxCPUsPerJob, "limits.MaxCPUsPerJob"); err != nil {
		return err
	}
	if err := m.validateOptionalNonNegative(limits.MaxNodesPerJob, "limits.MaxNodesPerJob"); err != nil {
		return err
	}
	if err := m.validateOptionalNonNegative(limits.MaxWallTimePerJob, "limits.MaxWallTimePerJob"); err != nil {
		return err
	}

	// Validate memory limits
	if limits.MaxMemoryPerNode != nil && *limits.MaxMemoryPerNode < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Max memory per node must be non-negative",
			"limits.MaxMemoryPerNode", *limits.MaxMemoryPerNode, nil,
		)
	}
	if limits.MaxMemoryPerCPU != nil && *limits.MaxMemoryPerCPU < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Max memory per CPU must be non-negative",
			"limits.MaxMemoryPerCPU", *limits.MaxMemoryPerCPU, nil,
		)
	}

	return nil
}

// ApplyQoSDefaults applies default values to QoS create request
func (m *QoSBaseManager) ApplyQoSDefaults(qos *types.QoSCreate) *types.QoSCreate {
	// Apply default priority if not set
	if qos.Priority == 0 {
		qos.Priority = 0 // Default priority is 0
	}

	// Ensure flags array is initialized
	if qos.Flags == nil {
		qos.Flags = []string{}
	}

	// Ensure preempt mode array is initialized
	if qos.PreemptMode == nil {
		qos.PreemptMode = []string{}
	}

	// Apply default usage factor if not set
	if qos.UsageFactor == 0 {
		qos.UsageFactor = 1.0 // Default usage factor is 1.0
	}

	return qos
}

// FilterQoSList applies filtering to a QoS list
func (m *QoSBaseManager) FilterQoSList(items []types.QoS, opts *types.QoSListOptions) []types.QoS {
	if opts == nil {
		return items
	}

	filtered := make([]types.QoS, 0, len(items))
	for _, qos := range items {
		if m.matchesQoSFilters(qos, opts) {
			filtered = append(filtered, qos)
		}
	}

	return filtered
}

// matchesQoSFilters checks if a QoS matches the given filters
func (m *QoSBaseManager) matchesQoSFilters(qos types.QoS, opts *types.QoSListOptions) bool {
	// Filter by names
	if len(opts.Names) > 0 {
		found := false
		for _, name := range opts.Names {
			if strings.EqualFold(qos.Name, name) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by accounts
	if len(opts.Accounts) > 0 {
		found := false
		for _, filterAccount := range opts.Accounts {
			for _, qosAccount := range qos.AllowedAccounts {
				if strings.EqualFold(qosAccount, filterAccount) {
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

	// Filter by users
	if len(opts.Users) > 0 {
		found := false
		for _, filterUser := range opts.Users {
			for _, qosUser := range qos.AllowedUsers {
				if strings.EqualFold(qosUser, filterUser) {
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

	return true
}

// validateOptionalNonNegative validates an optional non-negative integer
func (m *QoSBaseManager) validateOptionalNonNegative(value *int, fieldName string) error {
	if value != nil && *value < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			fmt.Sprintf("%s must be non-negative", fieldName),
			fieldName, *value, nil,
		)
	}
	return nil
}
