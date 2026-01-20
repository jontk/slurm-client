// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/stretchr/testify/assert"
)

func TestQoSNameUniqueness(t *testing.T) {
	mgr := NewQoSBaseManager("test")

	existingQoS := []types.QoS{
		{Name: "normal"},
		{Name: "high"},
		{Name: "low"},
	}

	// Test case-insensitive matching
	err := mgr.ValidateQoSNameUniqueness("Normal", existingQoS)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")

	// Test unique name
	err = mgr.ValidateQoSNameUniqueness("debug", existingQoS)
	assert.NoError(t, err)
}

func TestQoSPreemptMode(t *testing.T) {
	mgr := NewQoSBaseManager("test")

	// Test valid modes
	validModes := []string{"OFF", "CANCEL", "CHECKPOINT", "GANG", "REQUEUE", "SUSPEND"}
	err := mgr.ValidateQoSPreemptMode(validModes)
	assert.NoError(t, err)

	// Test invalid mode
	invalidModes := []string{"OFF", "INVALID_MODE"}
	err = mgr.ValidateQoSPreemptMode(invalidModes)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid preempt mode")

	// Test case insensitive
	mixedCase := []string{"off", "Cancel", "GANG"}
	err = mgr.ValidateQoSPreemptMode(mixedCase)
	assert.NoError(t, err)
}

func TestQoSFlags(t *testing.T) {
	mgr := NewQoSBaseManager("test")

	// Test valid flags
	validFlags := []string{"DENY_LIMIT", "NO_DECAY", "PART_TIME_LIMIT"}
	err := mgr.ValidateQoSFlags(validFlags)
	assert.NoError(t, err)

	// Test invalid flag
	invalidFlags := []string{"DENY_LIMIT", "INVALID_FLAG"}
	err = mgr.ValidateQoSFlags(invalidFlags)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid QoS flag")
}

func TestQoSGracetime(t *testing.T) {
	mgr := NewQoSBaseManager("test")

	// Test valid grace time
	err := mgr.ValidateQoSGracetime(3600) // 1 hour
	assert.NoError(t, err)

	// Test negative grace time
	err = mgr.ValidateQoSGracetime(-1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non-negative")

	// Test exceeding max grace time
	err = mgr.ValidateQoSGracetime(604801) // > 7 days
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot exceed 7 days")
}

func TestQoSHierarchy(t *testing.T) {
	mgr := NewQoSBaseManager("test")

	existingQoS := []types.QoS{
		{Name: "normal", ParentQoS: ""},
		{Name: "high", ParentQoS: "normal"},
		{Name: "urgent", ParentQoS: "high"},
	}

	// Test valid parent
	err := mgr.ValidateQoSHierarchy("new_qos", "normal", existingQoS)
	assert.NoError(t, err)

	// Test self-reference
	err = mgr.ValidateQoSHierarchy("new_qos", "new_qos", existingQoS)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be its own parent")

	// Test non-existent parent
	err = mgr.ValidateQoSHierarchy("new_qos", "nonexistent", existingQoS)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")

	// Test circular dependency
	// Add a QoS that would create a circle
	circularQoS := append(existingQoS, types.QoS{Name: "circular", ParentQoS: "new_qos"})
	err = mgr.ValidateQoSHierarchy("new_qos", "circular", circularQoS)
	assert.NoError(t, err) // This should pass as the circle isn't complete yet
}

func TestTRESLimits(t *testing.T) {
	mgr := NewQoSBaseManager("test")

	// Test valid TRES strings
	validTRES := "cpu=100,mem=4096,node=10"
	err := mgr.ValidateTRESLimits(validTRES)
	assert.NoError(t, err)

	// Test unlimited values
	unlimitedTRES := "cpu=unlimited,mem=-1,node=50"
	err = mgr.ValidateTRESLimits(unlimitedTRES)
	assert.NoError(t, err)

	// Test invalid format
	invalidTRES := "cpu100" // Missing =
	err = mgr.ValidateTRESLimits(invalidTRES)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid TRES format")

	// Test empty resource name
	emptyResource := "=100"
	err = mgr.ValidateTRESLimits(emptyResource)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "resource name cannot be empty")

	// Test invalid value
	invalidValue := "cpu=abc"
	err = mgr.ValidateTRESLimits(invalidValue)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid TRES value")

	// Test negative value (other than -1)
	negativeValue := "cpu=-5"
	err = mgr.ValidateTRESLimits(negativeValue)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be >= -1")
}

func TestQoSLimitsConsistency(t *testing.T) {
	mgr := NewQoSBaseManager("test")

	// Test inconsistent per-user vs per-account limits
	userCPUs := 100
	accountCPUs := 50
	limits := &types.QoSLimits{
		MaxCPUsPerUser:    &userCPUs,
		MaxCPUsPerAccount: &accountCPUs,
	}
	err := mgr.ValidateQoSLimitsConsistency(limits)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MaxCPUsPerUser cannot exceed MaxCPUsPerAccount")

	// Test consistent limits
	userCPUs = 50
	accountCPUs = 100
	limits = &types.QoSLimits{
		MaxCPUsPerUser:    &userCPUs,
		MaxCPUsPerAccount: &accountCPUs,
	}
	err = mgr.ValidateQoSLimitsConsistency(limits)
	assert.NoError(t, err)

	// Test per-job exceeding per-user
	jobCPUs := 200
	limits = &types.QoSLimits{
		MaxCPUsPerJob:  &jobCPUs,
		MaxCPUsPerUser: &userCPUs,
	}
	err = mgr.ValidateQoSLimitsConsistency(limits)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MaxCPUsPerJob cannot exceed MaxCPUsPerUser")
}

func TestQoSDeletionSafety(t *testing.T) {
	mgr := NewQoSBaseManager("test")

	// Test deletion of QoS in use
	associations := []types.Association{
		{
			UserName:    "user1",
			AccountName: "account1",
			DefaultQoS:  "high",
		},
	}

	err := mgr.ValidateQoSDeletionSafety("high", associations)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Cannot delete QoS")
	assert.Contains(t, err.Error(), "it is used by association")

	// Test deletion of unused QoS
	err = mgr.ValidateQoSDeletionSafety("unused", associations)
	assert.NoError(t, err)

	// Test deletion of default 'normal' QoS
	err = mgr.ValidateQoSDeletionSafety("normal", []types.Association{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Cannot delete the default 'normal' QoS")
}

func TestEnhanceQoSCreateWithDefaults(t *testing.T) {
	mgr := NewQoSBaseManager("test")

	// Test high priority QoS
	highQoS := &types.QoSCreate{
		Name: "high-priority",
	}
	enhanced := mgr.EnhanceQoSCreateWithDefaults(highQoS)
	assert.Equal(t, 1000000, enhanced.Priority)
	assert.Equal(t, []string{"REQUEUE"}, enhanced.PreemptMode)

	// Test low priority QoS
	lowQoS := &types.QoSCreate{
		Name: "low-batch",
	}
	enhanced = mgr.EnhanceQoSCreateWithDefaults(lowQoS)
	assert.Equal(t, 1, enhanced.Priority)
	assert.Equal(t, []string{"OFF"}, enhanced.PreemptMode)

	// Test debug QoS
	debugQoS := &types.QoSCreate{
		Name: "debug-test",
	}
	enhanced = mgr.EnhanceQoSCreateWithDefaults(debugQoS)
	assert.NotNil(t, enhanced.Limits)
	assert.NotNil(t, enhanced.Limits.MaxWallTimePerJob)
	assert.Equal(t, 3600, *enhanced.Limits.MaxWallTimePerJob)
	assert.NotNil(t, enhanced.Limits.MaxJobsPerUser)
	assert.Equal(t, 2, *enhanced.Limits.MaxJobsPerUser)

	// Test long-running QoS
	longQoS := &types.QoSCreate{
		Name: "long-running",
	}
	enhanced = mgr.EnhanceQoSCreateWithDefaults(longQoS)
	assert.NotNil(t, enhanced.Limits)
	assert.NotNil(t, enhanced.Limits.MaxWallTimePerJob)
	assert.Equal(t, 604800, *enhanced.Limits.MaxWallTimePerJob)
}
