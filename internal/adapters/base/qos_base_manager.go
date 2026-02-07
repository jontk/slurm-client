// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package base

import (
	"strings"

	types "github.com/jontk/slurm-client/api"
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
			"QoS creation data is required",
			"qos", qos, nil,
		)
	}
	if err := m.ValidateResourceName(qos.Name, "QoS name"); err != nil {
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
			"QoS update data is required",
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

// ValidateQoSLimits validates QoS resource limits using the nested structure
func (m *QoSBaseManager) ValidateQoSLimits(limits *types.QoSLimits) error {
	// Validate factor (limit factor)
	if limits.Factor != nil && *limits.Factor < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Limit factor must be non-negative",
			"limits.Factor", *limits.Factor, nil,
		)
	}
	// Validate grace time
	if limits.GraceTime != nil && *limits.GraceTime < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Grace time must be non-negative",
			"limits.GraceTime", *limits.GraceTime, nil,
		)
	}
	// Validate max limits if provided
	if limits.Max != nil {
		if err := m.validateQoSMaxLimits(limits.Max); err != nil {
			return err
		}
	}
	return nil
}

// validateQoSMaxLimits validates the max limits section
func (m *QoSBaseManager) validateQoSMaxLimits(maxLimits *types.QoSLimitsMax) error {
	// Validate active jobs
	if maxLimits.ActiveJobs != nil {
		if maxLimits.ActiveJobs.Count != nil {
			// Count is uint32, no need to check for negative
		}
	}
	// Validate jobs limits
	if maxLimits.Jobs != nil {
		if err := m.validateQoSJobsLimits(maxLimits.Jobs); err != nil {
			return err
		}
	}
	// Validate wall clock limits
	if maxLimits.WallClock != nil && maxLimits.WallClock.Per != nil {
		// WallClock values are uint32, no need to check for negative
	}
	return nil
}

// validateQoSJobsLimits validates job-related limits
func (m *QoSBaseManager) validateQoSJobsLimits(jobs *types.QoSLimitsMaxJobs) error {
	// Jobs limits use uint32 which are always non-negative
	// Additional validation can be added here if needed
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
	qosName := derefString(qos.Name)
	// Filter by names
	if len(opts.Names) > 0 {
		found := false
		for _, name := range opts.Names {
			if strings.EqualFold(qosName, name) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	// Note: Account and User filters removed as the generated QoS type
	// does not have AllowedAccounts or AllowedUsers fields.
	// These may need to be looked up via associations if needed.
	return true
}
