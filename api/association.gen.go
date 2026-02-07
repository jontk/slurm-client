// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// Association represents a SLURM Association.
type Association struct {
	Account *string `json:"account,omitempty"` // Account name
	Accounting []Accounting `json:"accounting,omitempty"` // Accounting records containing related resource usage
	Cluster *string `json:"cluster,omitempty"` // Cluster name
	Comment *string `json:"comment,omitempty"` // Arbitrary comment
	Default *AssociationDefault `json:"default,omitempty"`
	Flags []AssociationDefaultFlagsValue `json:"flags,omitempty"` // Flags on the association
	ID *int32 `json:"id,omitempty"` // Unique ID (Association ID)
	IsDefault *bool `json:"is_default,omitempty"` // Is default association for user
	Lineage *string `json:"lineage,omitempty"` // Complete path up the hierarchy to the root association
	Max *AssociationMax `json:"max,omitempty"`
	Min *AssociationMin `json:"min,omitempty"`
	ParentAccount *string `json:"parent_account,omitempty"` // Name of parent account
	Partition *string `json:"partition,omitempty"` // Partition name
	Priority *uint32 `json:"priority,omitempty"` // Association priority factor (32 bit integer number with flags)
	QoS []string `json:"qos,omitempty"` // List of available QOS names (List of QOS names)
	SharesRaw *int32 `json:"shares_raw,omitempty"` // Allocated shares used for fairshare calculation
	User string `json:"user"` // User name
}


// AssociationDefault is a nested type within its parent.
type AssociationDefault struct {
	QoS *string `json:"qos,omitempty"` // Default QOS
}


// AssociationMax is a nested type within its parent.
type AssociationMax struct {
	Jobs *AssociationMaxJobs `json:"jobs,omitempty"`
	Per *AssociationMaxPer `json:"per,omitempty"`
	TRES *AssociationMaxTRES `json:"tres,omitempty"`
}


// AssociationMaxJobs is a nested type within its parent.
type AssociationMaxJobs struct {
	Accruing *uint32 `json:"accruing,omitempty"` // MaxJobsAccrue - Maximum number of pending jobs able to accrue age priority at...
	Active *uint32 `json:"active,omitempty"` // MaxJobs - Maximum number of running jobs per user in this association (32 bit...
	Per *AssociationMaxJobsPer `json:"per,omitempty"`
	Total *uint32 `json:"total,omitempty"` // MaxSubmitJobs - Maximum number of jobs in a pending or running state at any...
}


// AssociationMaxJobsPer is a nested type within its parent.
type AssociationMaxJobsPer struct {
	Accruing *uint32 `json:"accruing,omitempty"` // GrpJobsAccrue - Maximum number of pending jobs able to accrue age priority in...
	Count *uint32 `json:"count,omitempty"` // GrpJobs - Maximum number of running jobs in this association and its children...
	Submitted *uint32 `json:"submitted,omitempty"` // GrpSubmitJobs - Maximum number of jobs in a pending or running state at any...
	WallClock *uint32 `json:"wall_clock,omitempty"` // MaxWallDurationPerJob - Maximum wall clock time in minutes each job can use in...
}


// AssociationDefaultFlagsValue represents possible values for AssociationDefaultFlags field.
type AssociationDefaultFlagsValue string

// AssociationDefaultFlagsValue constants.
const (
	AssociationDefaultFlagsDeleted AssociationDefaultFlagsValue = "DELETED"
	AssociationDefaultFlagsNoupdate AssociationDefaultFlagsValue = "NoUpdate"
	AssociationDefaultFlagsExact AssociationDefaultFlagsValue = "Exact"
	AssociationDefaultFlagsNousersarecoords AssociationDefaultFlagsValue = "NoUsersAreCoords"
	AssociationDefaultFlagsUsersarecoords AssociationDefaultFlagsValue = "UsersAreCoords"
)


// AssociationMaxPer is a nested type within its parent.
type AssociationMaxPer struct {
	Account *AssociationMaxPerAccount `json:"account,omitempty"`
}


// AssociationMaxPerAccount is a nested type within its parent.
type AssociationMaxPerAccount struct {
	WallClock *uint32 `json:"wall_clock,omitempty"` // GrpWall - Maximum wall clock time in minutes able to be allocated by running...
}


// AssociationMaxTRES is a nested type within its parent.
type AssociationMaxTRES struct {
	Group *AssociationMaxTRESGroup `json:"group,omitempty"`
	Minutes *AssociationMaxTRESMinutes `json:"minutes,omitempty"`
	Per *AssociationMaxTRESPer `json:"per,omitempty"`
	Total []TRES `json:"total,omitempty"` // GrpTRES - Maximum number of TRES able to be allocated by running jobs in this...
}


// AssociationMaxTRESGroup is a nested type within its parent.
type AssociationMaxTRESGroup struct {
	Active []TRES `json:"active,omitempty"` // GrpTRESRunMins - Maximum number of TRES minutes able to be allocated by running...
	Minutes []TRES `json:"minutes,omitempty"` // GrpTRESMins - Maximum number of TRES minutes that can possibly be used by past,...
}


// AssociationMaxTRESMinutes is a nested type within its parent.
type AssociationMaxTRESMinutes struct {
	Per *AssociationMaxTRESMinutesPer `json:"per,omitempty"`
	Total []TRES `json:"total,omitempty"` // Not implemented
}


// AssociationMaxTRESMinutesPer is a nested type within its parent.
type AssociationMaxTRESMinutesPer struct {
	Job []TRES `json:"job,omitempty"` // MaxTRESMinsPerJob - Maximum number of TRES minutes each job can use in this...
}


// AssociationMaxTRESPer is a nested type within its parent.
type AssociationMaxTRESPer struct {
	Job []TRES `json:"job,omitempty"` // MaxTRESPerJob - Maximum number of TRES each job can use in this association
	Node []TRES `json:"node,omitempty"` // MaxTRESPerNode - Maximum number of TRES each node in a job allocation can use...
}


// AssociationMin is a nested type within its parent.
type AssociationMin struct {
	PriorityThreshold *uint32 `json:"priority_threshold,omitempty"` // MinPrioThreshold - Minimum priority required to reserve resources when...
}
