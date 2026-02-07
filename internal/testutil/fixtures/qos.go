// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package fixtures

import (
	types "github.com/jontk/slurm-client/api"
)

// QoSFixtures provides pre-configured QoS objects for testing
type QoSFixtures struct{}

// NewQoSFixtures creates a new QoS fixtures instance
func NewQoSFixtures() *QoSFixtures {
	return &QoSFixtures{}
}

// SimpleQoS returns a basic QoS with minimal fields
func (f *QoSFixtures) SimpleQoS(name string) *types.QoS {
	id := int32(1)
	desc := "Test QoS " + name
	priority := uint32(100)
	usageFactor := 1.0
	return &types.QoS{
		ID:          &id,
		Name:        &name,
		Description: &desc,
		Priority:    &priority,
		UsageFactor: &usageFactor,
	}
}

// QoSWithLimits returns a QoS with comprehensive limits
func (f *QoSFixtures) QoSWithLimits(name string) *types.QoS {
	qos := f.SimpleQoS(name)

	// Build nested limits structure
	maxJobsPerUser := uint32(10)
	maxJobsPerAccount := uint32(100)
	maxSubmitJobsPerUser := uint32(20)
	maxSubmitJobsPerAccount := uint32(200)
	grpJobs := uint32(50)
	grpSubmitJobs := uint32(100)
	maxWallPerJob := uint32(1440) // 24 hours in minutes

	qos.Limits = &types.QoSLimits{
		Max: &types.QoSLimitsMax{
			ActiveJobs: &types.QoSLimitsMaxActiveJobs{
				Count: &grpJobs,
			},
			Jobs: &types.QoSLimitsMaxJobs{
				Count: &grpSubmitJobs,
				ActiveJobs: &types.QoSLimitsMaxJobsActiveJobs{
					Per: &types.QoSLimitsMaxJobsActiveJobsPer{
						User:    &maxJobsPerUser,
						Account: &maxJobsPerAccount,
					},
				},
				Per: &types.QoSLimitsMaxJobsPer{
					User:    &maxSubmitJobsPerUser,
					Account: &maxSubmitJobsPerAccount,
				},
			},
			WallClock: &types.QoSLimitsMaxWallClock{
				Per: &types.QoSLimitsMaxWallClockPer{
					Job: &maxWallPerJob,
				},
			},
		},
	}
	return qos
}

// HighPriorityQoS returns a high-priority QoS configuration
func (f *QoSFixtures) HighPriorityQoS() *types.QoS {
	qos := f.QoSWithLimits("high-priority")
	priority := uint32(1000)
	usageFactor := 2.0
	usageThreshold := 0.95
	exemptTime := uint32(300) // 5 minutes
	qos.Priority = &priority
	qos.Flags = []types.QoSFlagsValue{types.QoSFlagsDenyLimit, types.QoSFlagsRequiredReservation}
	qos.Preempt = &types.QoSPreempt{
		Mode:       []types.ModeValue{types.ModeCancel},
		ExemptTime: &exemptTime,
	}
	qos.UsageFactor = &usageFactor
	qos.UsageThreshold = &usageThreshold
	return qos
}

// BatchQoS returns a QoS suitable for batch jobs
func (f *QoSFixtures) BatchQoS() *types.QoS {
	qos := f.QoSWithLimits("batch")
	priority := uint32(10)
	usageFactor := 0.5
	graceTime := int32(3600) // 1 hour
	qos.Priority = &priority
	qos.Flags = []types.QoSFlagsValue{types.QoSFlagsNoReserve}
	qos.UsageFactor = &usageFactor
	qos.Limits.GraceTime = &graceTime

	// More relaxed limits for batch - update wall clock
	maxWall := uint32(10080) // 7 days in minutes
	maxJobsPerUser := uint32(100)
	if qos.Limits.Max != nil && qos.Limits.Max.WallClock != nil && qos.Limits.Max.WallClock.Per != nil {
		qos.Limits.Max.WallClock.Per.Job = &maxWall
	}
	if qos.Limits.Max != nil && qos.Limits.Max.Jobs != nil && qos.Limits.Max.Jobs.ActiveJobs != nil && qos.Limits.Max.Jobs.ActiveJobs.Per != nil {
		qos.Limits.Max.Jobs.ActiveJobs.Per.User = &maxJobsPerUser
	}
	return qos
}

// QoSCreateRequest returns a QoSCreate request for testing
func (f *QoSFixtures) QoSCreateRequest(name string) *types.QoSCreate {
	return &types.QoSCreate{
		Name:           name,
		Description:    "Test QoS creation",
		Priority:       100,
		Flags:          []string{"DenyOnLimit"},
		PreemptMode:    []string{"cluster"},
		GraceTime:      300,
		UsageFactor:    1.5,
		UsageThreshold: 0.8,
		// Note: Limits structure has changed, using nil for simplicity in tests
		Limits: nil,
	}
}

// QoSUpdateRequest returns a QoSUpdate request for testing
func (f *QoSFixtures) QoSUpdateRequest() *types.QoSUpdate {
	priority := 200
	usageFactor := 2.0
	usageThreshold := 0.9
	return &types.QoSUpdate{
		Description:    stringPtr("Updated description"),
		Priority:       &priority,
		Flags:          &[]string{"DenyOnLimit", "RequiresReservation"},
		PreemptMode:    &[]string{"suspend"},
		UsageFactor:    &usageFactor,
		UsageThreshold: &usageThreshold,
		// Note: Limits structure has changed, using nil for simplicity in tests
		Limits: nil,
	}
}

// QoSList returns a list of diverse QoS entries
func (f *QoSFixtures) QoSList() []types.QoS {
	return []types.QoS{
		*f.SimpleQoS("normal"),
		*f.HighPriorityQoS(),
		*f.BatchQoS(),
		*f.QoSWithLimits("gpu-jobs"),
		*f.SimpleQoS("low-priority"),
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}
