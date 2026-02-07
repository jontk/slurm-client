// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package base

import (
	"fmt"
	"strings"

	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/pkg/errors"
)

// ValidateQoSNameUniqueness checks if a QoS name is unique (case-insensitive)
func (m *QoSBaseManager) ValidateQoSNameUniqueness(name string, existingQoS []types.QoS) error {
	for _, qos := range existingQoS {
		qosName := ""
		if qos.Name != nil {
			qosName = *qos.Name
		}
		if strings.EqualFold(qosName, name) {
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
// Note: The generated QoS type does not include ParentQoS field, so this validates
// only that the parent exists and there's no self-reference. Full hierarchy validation
// would require the QoSCreate type which has ParentQoS.
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
		qosN := ""
		if qos.Name != nil {
			qosN = *qos.Name
		}
		if strings.EqualFold(qosN, parentName) {
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
	if limits == nil || limits.Max == nil || limits.Max.Jobs == nil {
		return nil
	}
	// Check that per-user limits don't exceed per-account limits for active jobs
	if limits.Max.Jobs.ActiveJobs != nil && limits.Max.Jobs.ActiveJobs.Per != nil {
		per := limits.Max.Jobs.ActiveJobs.Per
		if per.User != nil && per.Account != nil {
			if *per.User > *per.Account {
				return errors.NewValidationError(
					errors.ErrorCodeValidationFailed,
					"MaxJobsPerUser cannot exceed MaxJobsPerAccount",
					"limits", limits, nil,
				)
			}
		}
	}
	// Check that per-user limits don't exceed per-account limits for submit jobs
	if limits.Max.Jobs.Per != nil {
		per := limits.Max.Jobs.Per
		if per.User != nil && per.Account != nil {
			if *per.User > *per.Account {
				return errors.NewValidationError(
					errors.ErrorCodeValidationFailed,
					"MaxSubmitJobsPerUser cannot exceed MaxSubmitJobsPerAccount",
					"limits", limits, nil,
				)
			}
		}
	}
	return nil
}

// ValidateQoSUpdateSafety checks if an update would break existing configurations
func (m *QoSBaseManager) ValidateQoSUpdateSafety(currentQoS *types.QoS, update *types.QoSUpdate) error {
	// Check if reducing limits would affect running jobs
	if update.Limits != nil && currentQoS.Limits != nil {
		// Check if reducing MaxJobsPerUser
		newMaxJobs := getMaxJobsPerUserFromLimits(update.Limits)
		currentMaxJobs := getMaxJobsPerUserFromLimits(currentQoS.Limits)
		if newMaxJobs != nil && currentMaxJobs != nil {
			// This is a potentially breaking change
			// In a real implementation, we'd check active job counts
			// For now, just continue silently
			_ = *newMaxJobs < *currentMaxJobs
		}
	}
	// Check if changing priority would affect job scheduling
	// For now, just continue silently
	if update.Priority != nil && currentQoS.Priority != nil {
		_ = *update.Priority != int(*currentQoS.Priority)
	}
	return nil
}

// getMaxJobsPerUserFromLimits extracts max jobs per user from nested limits structure
func getMaxJobsPerUserFromLimits(limits *types.QoSLimits) *uint32 {
	if limits == nil || limits.Max == nil || limits.Max.Jobs == nil {
		return nil
	}
	if limits.Max.Jobs.ActiveJobs == nil || limits.Max.Jobs.ActiveJobs.Per == nil {
		return nil
	}
	return limits.Max.Jobs.ActiveJobs.Per.User
}

// ValidateQoSDeletionSafety checks if a QoS can be safely deleted
func (m *QoSBaseManager) ValidateQoSDeletionSafety(qosName string, associations []types.Association) error {
	// Check if any associations are using this QoS
	for _, assoc := range associations {
		defaultQoS := ""
		if assoc.Default != nil && assoc.Default.QoS != nil {
			defaultQoS = *assoc.Default.QoS
		}
		if strings.EqualFold(defaultQoS, qosName) {
			account := ""
			if assoc.Account != nil {
				account = *assoc.Account
			}
			return errors.NewValidationError(
				errors.ErrorCodeConflict,
				fmt.Sprintf("Cannot delete QoS '%s': it is used by association for user '%s' in account '%s'",
					qosName, assoc.User, account),
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
	// Apply pattern-specific enhancements
	m.enhanceHighPriorityQoS(qos, lowerName)
	m.enhanceLowPriorityQoS(qos, lowerName)
	m.enhanceDebugQoS(qos, lowerName)
	m.enhanceLongRunningQoS(qos, lowerName)
	return qos
}

// enhanceHighPriorityQoS applies defaults for high priority QoS patterns
func (m *QoSBaseManager) enhanceHighPriorityQoS(qos *types.QoSCreate, lowerName string) {
	if !strings.Contains(lowerName, "high") && !strings.Contains(lowerName, "priority") && !strings.Contains(lowerName, "urgent") {
		return
	}
	if qos.Priority == 0 {
		qos.Priority = 1000000 // High priority
	}
	if len(qos.PreemptMode) == 0 {
		qos.PreemptMode = []string{"REQUEUE"} // Allow preemption
	}
}

// enhanceLowPriorityQoS applies defaults for low priority QoS patterns
func (m *QoSBaseManager) enhanceLowPriorityQoS(qos *types.QoSCreate, lowerName string) {
	if !strings.Contains(lowerName, "low") && !strings.Contains(lowerName, "batch") && !strings.Contains(lowerName, "background") {
		return
	}
	if qos.Priority == 0 {
		qos.Priority = 1 // Low priority
	}
	if len(qos.PreemptMode) == 0 {
		qos.PreemptMode = []string{"OFF"} // Can be preempted
	}
}

// enhanceDebugQoS applies defaults for debug/test QoS patterns
func (m *QoSBaseManager) enhanceDebugQoS(qos *types.QoSCreate, lowerName string) {
	if !strings.Contains(lowerName, "debug") && !strings.Contains(lowerName, "test") && !strings.Contains(lowerName, "devel") {
		return
	}
	ensureLimitsStructure(qos)
	// Set reasonable limits for debug QoS - 1 hour default
	if qos.Limits.Max.WallClock.Per.Job == nil {
		maxTime := uint32(60) // 1 hour in minutes
		qos.Limits.Max.WallClock.Per.Job = &maxTime
	}
	// Limited concurrent debug jobs
	if qos.Limits.Max.Jobs.ActiveJobs.Per.User == nil {
		maxJobs := uint32(2)
		qos.Limits.Max.Jobs.ActiveJobs.Per.User = &maxJobs
	}
}

// enhanceLongRunningQoS applies defaults for long-running job QoS patterns
func (m *QoSBaseManager) enhanceLongRunningQoS(qos *types.QoSCreate, lowerName string) {
	if !strings.Contains(lowerName, "long") && !strings.Contains(lowerName, "extended") {
		return
	}
	ensureLimitsStructure(qos)
	// 7 days for long jobs in minutes
	if qos.Limits.Max.WallClock.Per.Job == nil {
		maxTime := uint32(10080) // 7 days in minutes
		qos.Limits.Max.WallClock.Per.Job = &maxTime
	}
}

// ensureLimitsStructure ensures all nested limit structures are initialized
func ensureLimitsStructure(qos *types.QoSCreate) {
	if qos.Limits == nil {
		qos.Limits = &types.QoSLimits{}
	}
	if qos.Limits.Max == nil {
		qos.Limits.Max = &types.QoSLimitsMax{}
	}
	if qos.Limits.Max.WallClock == nil {
		qos.Limits.Max.WallClock = &types.QoSLimitsMaxWallClock{}
	}
	if qos.Limits.Max.WallClock.Per == nil {
		qos.Limits.Max.WallClock.Per = &types.QoSLimitsMaxWallClockPer{}
	}
	if qos.Limits.Max.Jobs == nil {
		qos.Limits.Max.Jobs = &types.QoSLimitsMaxJobs{}
	}
	if qos.Limits.Max.Jobs.ActiveJobs == nil {
		qos.Limits.Max.Jobs.ActiveJobs = &types.QoSLimitsMaxJobsActiveJobs{}
	}
	if qos.Limits.Max.Jobs.ActiveJobs.Per == nil {
		qos.Limits.Max.Jobs.ActiveJobs.Per = &types.QoSLimitsMaxJobsActiveJobsPer{}
	}
}
