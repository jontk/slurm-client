// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/pkg/errors"
)

// ValidateQoSNameUniqueness checks if a QoS name is unique (case-insensitive)
func (m *QoSBaseManager) ValidateQoSNameUniqueness(name string, existingQoS []types.QoS) error {
	for _, qos := range existingQoS {
		if strings.EqualFold(qos.Name, name) {
			return errors.NewValidationError(
				errors.ErrorCodeConflict,
				fmt.Sprintf("QoS with name '%s' already exists", name),
				"qos.Name", name, nil,
			)
		}
	}
	return nil
}

// ValidateQoSPreemptMode validates preemption mode values
func (m *QoSBaseManager) ValidateQoSPreemptMode(modes []string) error {
	validModes := map[string]bool{
		"OFF":        true,
		"CANCEL":     true,
		"CHECKPOINT": true,
		"GANG":       true,
		"REQUEUE":    true,
		"SUSPEND":    true,
	}

	for _, mode := range modes {
		upperMode := strings.ToUpper(mode)
		if !validModes[upperMode] {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				fmt.Sprintf("Invalid preempt mode: %s. Valid modes are: OFF, CANCEL, CHECKPOINT, GANG, REQUEUE, SUSPEND", mode),
				"preemptMode", mode, nil,
			)
		}
	}
	return nil
}

// ValidateQoSFlags validates QoS flag values
func (m *QoSBaseManager) ValidateQoSFlags(flags []string) error {
	validFlags := map[string]bool{
		"DENY_LIMIT":              true,
		"ENFORCE_USAGE_THRESHOLD": true,
		"NO_DECAY":                true,
		"NO_RESERVE":              true,
		"OVER_PART_QOS":           true,
		"PART_MAX_NODE":           true,
		"PART_MIN_NODE":           true,
		"PART_TIME_LIMIT":         true,
		"RELATIVE_PRIORITY":       true,
		"REQUIRES_RES":            true,
		"USAGE_FACTOR_SAFE":       true,
	}

	for _, flag := range flags {
		upperFlag := strings.ToUpper(flag)
		if !validFlags[upperFlag] {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				"Invalid QoS flag: "+flag,
				"flag", flag, nil,
			)
		}
	}
	return nil
}

// ValidateQoSGracetime validates grace time value
func (m *QoSBaseManager) ValidateQoSGracetime(seconds int) error {
	if seconds < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Grace time must be non-negative",
			"graceTime", seconds, nil,
		)
	}
	// Maximum grace time is 7 days (604800 seconds)
	if seconds > 604800 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Grace time cannot exceed 7 days (604800 seconds)",
			"graceTime", seconds, nil,
		)
	}
	return nil
}

// ValidateQoSHierarchy validates QoS parent-child relationships
func (m *QoSBaseManager) ValidateQoSHierarchy(qosName, parentName string, existingQoS []types.QoS) error {
	if parentName == "" {
		return nil // No parent is valid
	}

	// Check for self-reference
	if strings.EqualFold(qosName, parentName) {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"QoS cannot be its own parent",
			"parentQoS", parentName, nil,
		)
	}

	// Check if parent exists
	parentExists := false
	for _, qos := range existingQoS {
		if strings.EqualFold(qos.Name, parentName) {
			parentExists = true
			break
		}
	}

	if !parentExists {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			fmt.Sprintf("Parent QoS '%s' does not exist", parentName),
			"parentQoS", parentName, nil,
		)
	}

	// Check for circular dependencies
	if err := m.checkCircularDependency(qosName, parentName, existingQoS); err != nil {
		return err
	}

	return nil
}

// checkCircularDependency checks for circular parent-child relationships
func (m *QoSBaseManager) checkCircularDependency(qosName, parentName string, existingQoS []types.QoS) error {
	visited := make(map[string]bool)
	current := parentName

	for current != "" {
		if visited[strings.ToLower(current)] {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				"Circular dependency detected in QoS hierarchy",
				"qosHierarchy", fmt.Sprintf("%s -> %s", qosName, parentName), nil,
			)
		}
		visited[strings.ToLower(current)] = true

		// Find the parent of current
		found := false
		for _, qos := range existingQoS {
			if strings.EqualFold(qos.Name, current) {
				current = qos.ParentQoS
				found = true
				break
			}
		}
		if !found {
			break
		}
	}

	return nil
}

// ValidateTRESLimits validates TRES limit strings
func (m *QoSBaseManager) ValidateTRESLimits(tresString string) error {
	if tresString == "" {
		return nil
	}

	// TRES format: "resource1=value1,resource2=value2,..."
	parts := strings.Split(tresString, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Each part should be "resource=value"
		subParts := strings.Split(part, "=")
		if len(subParts) != 2 {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				fmt.Sprintf("Invalid TRES format: %s. Expected format: resource=value", part),
				"tres", part, nil,
			)
		}

		resource := strings.TrimSpace(subParts[0])
		value := strings.TrimSpace(subParts[1])

		// Validate resource name
		if resource == "" {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				"TRES resource name cannot be empty",
				"tres", part, nil,
			)
		}

		// Validate value (should be numeric or "unlimited")
		if value != "unlimited" && value != "-1" {
			// Try to parse as number
			var intVal int64
			if _, err := fmt.Sscanf(value, "%d", &intVal); err != nil {
				return errors.NewValidationError(
					errors.ErrorCodeValidationFailed,
					fmt.Sprintf("Invalid TRES value: %s. Must be a number, 'unlimited', or '-1'", value),
					"tres", value, nil,
				)
			}
			if intVal < -1 {
				return errors.NewValidationError(
					errors.ErrorCodeValidationFailed,
					"TRES value must be >= -1",
					"tres", value, nil,
				)
			}
		}
	}

	return nil
}

// ValidateQoSLimitsConsistency checks for logical consistency in limits
func (m *QoSBaseManager) ValidateQoSLimitsConsistency(limits *types.QoSLimits) error {
	// Check that per-user limits don't exceed per-account limits where applicable
	if limits.MaxCPUsPerUser != nil && limits.MaxCPUsPerAccount != nil {
		if *limits.MaxCPUsPerUser > *limits.MaxCPUsPerAccount {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				"MaxCPUsPerUser cannot exceed MaxCPUsPerAccount",
				"limits", limits, nil,
			)
		}
	}

	if limits.MaxJobsPerUser != nil && limits.MaxJobsPerAccount != nil {
		if *limits.MaxJobsPerUser > *limits.MaxJobsPerAccount {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				"MaxJobsPerUser cannot exceed MaxJobsPerAccount",
				"limits", limits, nil,
			)
		}
	}

	if limits.MaxNodesPerUser != nil && limits.MaxNodesPerAccount != nil {
		if *limits.MaxNodesPerUser > *limits.MaxNodesPerAccount {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				"MaxNodesPerUser cannot exceed MaxNodesPerAccount",
				"limits", limits, nil,
			)
		}
	}

	// Check that per-job limits don't exceed per-user limits
	if limits.MaxCPUsPerJob != nil && limits.MaxCPUsPerUser != nil {
		if *limits.MaxCPUsPerJob > *limits.MaxCPUsPerUser {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				"MaxCPUsPerJob cannot exceed MaxCPUsPerUser",
				"limits", limits, nil,
			)
		}
	}

	if limits.MaxNodesPerJob != nil && limits.MaxNodesPerUser != nil {
		if *limits.MaxNodesPerJob > *limits.MaxNodesPerUser {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				"MaxNodesPerJob cannot exceed MaxNodesPerUser",
				"limits", limits, nil,
			)
		}
	}

	return nil
}

// ValidateQoSUpdateSafety checks if an update would break existing configurations
func (m *QoSBaseManager) ValidateQoSUpdateSafety(currentQoS *types.QoS, update *types.QoSUpdate) error {
	// Check if reducing limits would affect running jobs
	if update.Limits != nil {
		// If we're reducing MaxJobsPerUser, check if any user currently exceeds the new limit
		if update.Limits.MaxJobsPerUser != nil && currentQoS.Limits != nil && currentQoS.Limits.MaxJobsPerUser != nil {
			if *update.Limits.MaxJobsPerUser < *currentQoS.Limits.MaxJobsPerUser {
				// This is a potentially breaking change
				// In a real implementation, we'd check active job counts
				// Log warning - we'll add proper logging later
				// For now, just continue silently
			}
		}
	}

	// Check if changing priority would affect job scheduling
	if update.Priority != nil && *update.Priority != currentQoS.Priority {
		// Log warning - we'll add proper logging later
		// For now, just continue silently
	}

	// Check if removing allowed accounts/users would affect associations
	// Note: This check is commented out until we have AllowedAccounts field in QoSUpdate
	// if update.AllowedAccounts != nil && len(*update.AllowedAccounts) < len(currentQoS.AllowedAccounts) {
	//     // Log warning
	// }

	return nil
}

// ValidateQoSDeletionSafety checks if a QoS can be safely deleted
func (m *QoSBaseManager) ValidateQoSDeletionSafety(qosName string, associations []types.Association) error {
	// Check if any associations are using this QoS
	for _, assoc := range associations {
		if strings.EqualFold(assoc.DefaultQoS, qosName) {
			return errors.NewValidationError(
				errors.ErrorCodeConflict,
				fmt.Sprintf("Cannot delete QoS '%s': it is used by association for user '%s' in account '%s'",
					qosName, assoc.UserName, assoc.AccountName),
				"qos", qosName, nil,
			)
		}
	}

	// Check if it's the default QoS (usually "normal")
	if strings.EqualFold(qosName, "normal") {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Cannot delete the default 'normal' QoS",
			"qos", qosName, nil,
		)
	}

	return nil
}

// EnhanceQoSCreateWithDefaults adds intelligent defaults based on the QoS name and purpose
func (m *QoSBaseManager) EnhanceQoSCreateWithDefaults(qos *types.QoSCreate) *types.QoSCreate {
	// Apply base defaults first
	qos = m.ApplyQoSDefaults(qos)

	// Apply intelligent defaults based on QoS name patterns
	lowerName := strings.ToLower(qos.Name)

	// High priority QoS patterns
	if strings.Contains(lowerName, "high") || strings.Contains(lowerName, "priority") || strings.Contains(lowerName, "urgent") {
		if qos.Priority == 0 {
			qos.Priority = 1000000 // High priority
		}
		if len(qos.PreemptMode) == 0 {
			qos.PreemptMode = []string{"REQUEUE"} // Allow preemption
		}
	}

	// Low priority QoS patterns
	if strings.Contains(lowerName, "low") || strings.Contains(lowerName, "batch") || strings.Contains(lowerName, "background") {
		if qos.Priority == 0 {
			qos.Priority = 1 // Low priority
		}
		if len(qos.PreemptMode) == 0 {
			qos.PreemptMode = []string{"OFF"} // Can be preempted
		}
	}

	// Debug/test QoS patterns
	if strings.Contains(lowerName, "debug") || strings.Contains(lowerName, "test") || strings.Contains(lowerName, "devel") {
		if qos.Limits == nil {
			qos.Limits = &types.QoSLimits{}
		}
		// Set reasonable limits for debug QoS
		if qos.Limits.MaxWallTimePerJob == nil {
			maxTime := 3600 // 1 hour default for debug
			qos.Limits.MaxWallTimePerJob = &maxTime
		}
		if qos.Limits.MaxJobsPerUser == nil {
			maxJobs := 2 // Limited concurrent debug jobs
			qos.Limits.MaxJobsPerUser = &maxJobs
		}
	}

	// Long-running job QoS patterns
	if strings.Contains(lowerName, "long") || strings.Contains(lowerName, "extended") {
		if qos.Limits == nil {
			qos.Limits = &types.QoSLimits{}
		}
		if qos.Limits.MaxWallTimePerJob == nil {
			maxTime := 604800 // 7 days for long jobs
			qos.Limits.MaxWallTimePerJob = &maxTime
		}
	}

	return qos
}
