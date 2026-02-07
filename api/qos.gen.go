// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// QoS represents a SLURM QoS.
type QoS struct {
	Description *string `json:"description,omitempty"` // Arbitrary description
	Flags []QoSFlagsValue `json:"flags,omitempty"` // Flags, to avoid modifying current values specify NOT_SET
	ID *int32 `json:"id,omitempty"` // Unique ID
	Limits *QoSLimits `json:"limits,omitempty"`
	Name *string `json:"name,omitempty"` // Name
	Preempt *QoSPreempt `json:"preempt,omitempty"`
	Priority *uint32 `json:"priority,omitempty"` // Priority - QOS priority factor (32 bit integer number with flags)
	UsageFactor *float64 `json:"usage_factor,omitempty"` // UsageFactor - A float that is factored into a job's TRES usage (64 bit floating...
	UsageThreshold *float64 `json:"usage_threshold,omitempty"` // UsageThreshold - A float representing the lowest fairshare of an association...
}


// QoSLimits is a nested type within its parent.
type QoSLimits struct {
	Factor *float64 `json:"factor,omitempty"` // LimitFactor - A float that is factored into an association's [Grp|Max]TRES...
	GraceTime *int32 `json:"grace_time,omitempty"` // GraceTime - Preemption grace time in seconds to be extended to a job which has...
	Max *QoSLimitsMax `json:"max,omitempty"`
	Min *QoSLimitsMin `json:"min,omitempty"`
}


// QoSLimitsMax is a nested type within its parent.
type QoSLimitsMax struct {
	Accruing *QoSLimitsMaxAccruing `json:"accruing,omitempty"`
	ActiveJobs *QoSLimitsMaxActiveJobs `json:"active_jobs,omitempty"`
	Jobs *QoSLimitsMaxJobs `json:"jobs,omitempty"`
	TRES *QoSLimitsMaxTRES `json:"tres,omitempty"`
	WallClock *QoSLimitsMaxWallClock `json:"wall_clock,omitempty"`
}


// QoSLimitsMaxAccruing is a nested type within its parent.
type QoSLimitsMaxAccruing struct {
	Per *QoSLimitsMaxAccruingPer `json:"per,omitempty"`
}


// QoSLimitsMaxAccruingPer is a nested type within its parent.
type QoSLimitsMaxAccruingPer struct {
	Account *uint32 `json:"account,omitempty"` // MaxJobsAccruePerAccount - Maximum number of pending jobs an account (or...
	User *uint32 `json:"user,omitempty"` // MaxJobsAccruePerUser - Maximum number of pending jobs a user can have accruing...
}


// QoSFlagsValue represents possible values for QoSFlags field.
type QoSFlagsValue string

// QoSFlagsValue constants.
const (
	QoSFlagsNotSet QoSFlagsValue = "NOT_SET"
	QoSFlagsAdd QoSFlagsValue = "ADD"
	QoSFlagsRemove QoSFlagsValue = "REMOVE"
	QoSFlagsDeleted QoSFlagsValue = "DELETED"
	QoSFlagsPartitionMinimumNode QoSFlagsValue = "PARTITION_MINIMUM_NODE"
	QoSFlagsPartitionMaximumNode QoSFlagsValue = "PARTITION_MAXIMUM_NODE"
	QoSFlagsPartitionTimeLimit QoSFlagsValue = "PARTITION_TIME_LIMIT"
	QoSFlagsEnforceUsageThreshold QoSFlagsValue = "ENFORCE_USAGE_THRESHOLD"
	QoSFlagsNoReserve QoSFlagsValue = "NO_RESERVE"
	QoSFlagsRequiredReservation QoSFlagsValue = "REQUIRED_RESERVATION"
	QoSFlagsDenyLimit QoSFlagsValue = "DENY_LIMIT"
	QoSFlagsOverridePartitionQoS QoSFlagsValue = "OVERRIDE_PARTITION_QOS"
	QoSFlagsPartitionQoS QoSFlagsValue = "PARTITION_QOS"
	QoSFlagsNoDecay QoSFlagsValue = "NO_DECAY"
	QoSFlagsUsageFactorSafe QoSFlagsValue = "USAGE_FACTOR_SAFE"
	QoSFlagsRelative QoSFlagsValue = "RELATIVE"
)


// QoSLimitsMaxActiveJobs is a nested type within its parent.
type QoSLimitsMaxActiveJobs struct {
	Accruing *uint32 `json:"accruing,omitempty"` // GrpJobsAccrue - Maximum number of pending jobs able to accrue age priority (32...
	Count *uint32 `json:"count,omitempty"` // GrpJobs - Maximum number of running jobs (32 bit integer number with flags)
}


// QoSLimitsMaxJobs is a nested type within its parent.
type QoSLimitsMaxJobs struct {
	ActiveJobs *QoSLimitsMaxJobsActiveJobs `json:"active_jobs,omitempty"`
	Count *uint32 `json:"count,omitempty"` // GrpSubmitJobs - Maximum number of jobs in a pending or running state at any...
	Per *QoSLimitsMaxJobsPer `json:"per,omitempty"`
}


// QoSLimitsMaxJobsActiveJobs is a nested type within its parent.
type QoSLimitsMaxJobsActiveJobs struct {
	Per *QoSLimitsMaxJobsActiveJobsPer `json:"per,omitempty"`
}


// QoSLimitsMaxJobsActiveJobsPer is a nested type within its parent.
type QoSLimitsMaxJobsActiveJobsPer struct {
	Account *uint32 `json:"account,omitempty"` // MaxJobsPerAccount - Maximum number of running jobs per account (32 bit integer...
	User *uint32 `json:"user,omitempty"` // MaxJobsPerUser - Maximum number of running jobs per user (32 bit integer number...
}


// QoSLimitsMaxJobsPer is a nested type within its parent.
type QoSLimitsMaxJobsPer struct {
	Account *uint32 `json:"account,omitempty"` // MaxSubmitJobsPerAccount - Maximum number of jobs in a pending or running state...
	User *uint32 `json:"user,omitempty"` // MaxSubmitJobsPerUser - Maximum number of jobs in a pending or running state per...
}


// QoSLimitsMaxTRES is a nested type within its parent.
type QoSLimitsMaxTRES struct {
	Minutes *QoSLimitsMaxTRESMinutes `json:"minutes,omitempty"`
	Per *QoSLimitsMaxTRESPer `json:"per,omitempty"`
	Total []TRES `json:"total,omitempty"` // GrpTRES - Maximum number of TRES able to be allocated by running jobs
}


// QoSLimitsMaxTRESMinutes is a nested type within its parent.
type QoSLimitsMaxTRESMinutes struct {
	Per *QoSLimitsMaxTRESMinutesPer `json:"per,omitempty"`
	Total []TRES `json:"total,omitempty"` // GrpTRESMins - Maximum number of TRES minutes that can possibly be used by past,...
}


// QoSLimitsMaxTRESMinutesPer is a nested type within its parent.
type QoSLimitsMaxTRESMinutesPer struct {
	Account []TRES `json:"account,omitempty"` // MaxTRESRunMinsPerAccount - Maximum number of TRES minutes each account can use
	Job []TRES `json:"job,omitempty"` // MaxTRESMinsPerJob - Maximum number of TRES minutes each job can use
	QoS []TRES `json:"qos,omitempty"` // GrpTRESRunMins - Maximum number of TRES minutes able to be allocated by running...
	User []TRES `json:"user,omitempty"` // MaxTRESRunMinsPerUser - Maximum number of TRES minutes each user can use
}


// QoSLimitsMaxTRESPer is a nested type within its parent.
type QoSLimitsMaxTRESPer struct {
	Account []TRES `json:"account,omitempty"` // MaxTRESPerAccount - Maximum number of TRES each account can use
	Job []TRES `json:"job,omitempty"` // MaxTRESPerJob - Maximum number of TRES each job can use
	Node []TRES `json:"node,omitempty"` // MaxTRESPerNode - Maximum number of TRES each node in a job allocation can use
	User []TRES `json:"user,omitempty"` // MaxTRESPerUser - Maximum number of TRES each user can use
}


// QoSLimitsMaxWallClock is a nested type within its parent.
type QoSLimitsMaxWallClock struct {
	Per *QoSLimitsMaxWallClockPer `json:"per,omitempty"`
}


// QoSLimitsMaxWallClockPer is a nested type within its parent.
type QoSLimitsMaxWallClockPer struct {
	Job *uint32 `json:"job,omitempty"` // MaxWallDurationPerJob - Maximum wall clock time in minutes each job can use (32...
	QoS *uint32 `json:"qos,omitempty"` // GrpWall - Maximum wall clock time in minutes able to be allocated by running...
}


// QoSLimitsMin is a nested type within its parent.
type QoSLimitsMin struct {
	PriorityThreshold *uint32 `json:"priority_threshold,omitempty"` // MinPrioThreshold - Minimum priority required to reserve resources when...
	TRES *QoSLimitsMinTRES `json:"tres,omitempty"`
}


// QoSLimitsMinTRES is a nested type within its parent.
type QoSLimitsMinTRES struct {
	Per *QoSLimitsMinTRESPer `json:"per,omitempty"`
}


// QoSLimitsMinTRESPer is a nested type within its parent.
type QoSLimitsMinTRESPer struct {
	Job []TRES `json:"job,omitempty"` // MinTRESPerJob - Minimum number of TRES each job running under this QOS must...
}


// QoSPreempt is a nested type within its parent.
type QoSPreempt struct {
	ExemptTime *uint32 `json:"exempt_time,omitempty"` // PreemptExemptTime - Specifies a minimum run time for jobs before they are...
	List []string `json:"list,omitempty"` // Other QOS's this QOS can preempt
	Mode []ModeValue `json:"mode,omitempty"` // PreemptMode - Mechanism used to preempt jobs or enable gang scheduling
}


// ModeValue represents possible values for Mode field.
type ModeValue string

// ModeValue constants.
const (
	ModeDisabled ModeValue = "DISABLED"
	ModeSuspend ModeValue = "SUSPEND"
	ModeRequeue ModeValue = "REQUEUE"
	ModeCancel ModeValue = "CANCEL"
	ModeGang ModeValue = "GANG"
)
