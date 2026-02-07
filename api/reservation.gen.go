// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

import "time"

// Reservation represents a SLURM Reservation.
type Reservation struct {
	Accounts *string `json:"accounts,omitempty"` // Comma-separated list of permitted accounts
	BurstBuffer *string `json:"burst_buffer,omitempty"` // BurstBuffer - Burst buffer resources reserved
	CoreCount *int32 `json:"core_count,omitempty"` // CoreCnt - Number of cores reserved
	CoreSpecializations []ReservationCoreSpec `json:"core_specializations,omitempty"` // Reserved cores specification
	EndTime time.Time `json:"end_time,omitempty"` // EndTime (UNIX timestamp) (UNIX timestamp or time string recognized by Slurm...
	Features *string `json:"features,omitempty"` // Features - Expression describing the reservation's required node features
	Flags []ReservationFlagsValue `json:"flags,omitempty"` // Flags associated with this reservation
	Groups *string `json:"groups,omitempty"` // Groups - Comma-separated list of permitted groups
	Licenses *string `json:"licenses,omitempty"` // Licenses - Comma-separated list of licenses reserved
	MaxStartDelay *int32 `json:"max_start_delay,omitempty"` // MaxStartDelay - Maximum time an eligible job not requesting this reservation...
	Name *string `json:"name,omitempty"` // ReservationName - Name of the reservation
	NodeCount *int32 `json:"node_count,omitempty"` // NodeCnt - Number of nodes reserved
	NodeList *string `json:"node_list,omitempty"` // Nodes - Comma-separated list of node names and/or node ranges reserved
	Partition *string `json:"partition,omitempty"` // PartitionName - Partition used to reserve nodes from
	PurgeCompleted *ReservationPurgeCompleted `json:"purge_completed,omitempty"`
	StartTime time.Time `json:"start_time,omitempty"` // StartTime (UNIX timestamp) (UNIX timestamp or time string recognized by Slurm...
	TRES *string `json:"tres,omitempty"` // Comma-separated list of required TRES
	Users *string `json:"users,omitempty"` // Comma-separated list of permitted users
	Watts *uint32 `json:"watts,omitempty"` // 32 bit integer number with flags
}


// ReservationPurgeCompleted is a nested type within its parent.
type ReservationPurgeCompleted struct {
	Time *uint32 `json:"time,omitempty"` // If PURGE_COMP flag is set, the number of seconds this reservation will sit idle...
}


// ReservationFlagsValue represents possible values for ReservationFlags field.
type ReservationFlagsValue string

// ReservationFlagsValue constants.
const (
	ReservationFlagsMaint ReservationFlagsValue = "MAINT"
	ReservationFlagsNoMaint ReservationFlagsValue = "NO_MAINT"
	ReservationFlagsDaily ReservationFlagsValue = "DAILY"
	ReservationFlagsNoDaily ReservationFlagsValue = "NO_DAILY"
	ReservationFlagsWeekly ReservationFlagsValue = "WEEKLY"
	ReservationFlagsNoWeekly ReservationFlagsValue = "NO_WEEKLY"
	ReservationFlagsIgnoreJobs ReservationFlagsValue = "IGNORE_JOBS"
	ReservationFlagsNoIgnoreJobs ReservationFlagsValue = "NO_IGNORE_JOBS"
	ReservationFlagsAnyNodes ReservationFlagsValue = "ANY_NODES"
	ReservationFlagsNoAnyNodes ReservationFlagsValue = "NO_ANY_NODES"
	ReservationFlagsStatic ReservationFlagsValue = "STATIC"
	ReservationFlagsNoStatic ReservationFlagsValue = "NO_STATIC"
	ReservationFlagsPartNodes ReservationFlagsValue = "PART_NODES"
	ReservationFlagsNoPartNodes ReservationFlagsValue = "NO_PART_NODES"
	ReservationFlagsOverlap ReservationFlagsValue = "OVERLAP"
	ReservationFlagsSpecNodes ReservationFlagsValue = "SPEC_NODES"
	ReservationFlagsTimeFloat ReservationFlagsValue = "TIME_FLOAT"
	ReservationFlagsReplace ReservationFlagsValue = "REPLACE"
	ReservationFlagsAllNodes ReservationFlagsValue = "ALL_NODES"
	ReservationFlagsPurgeComp ReservationFlagsValue = "PURGE_COMP"
	ReservationFlagsWeekday ReservationFlagsValue = "WEEKDAY"
	ReservationFlagsNoWeekday ReservationFlagsValue = "NO_WEEKDAY"
	ReservationFlagsWeekend ReservationFlagsValue = "WEEKEND"
	ReservationFlagsNoWeekend ReservationFlagsValue = "NO_WEEKEND"
	ReservationFlagsFlex ReservationFlagsValue = "FLEX"
	ReservationFlagsNoFlex ReservationFlagsValue = "NO_FLEX"
	ReservationFlagsDurationPlus ReservationFlagsValue = "DURATION_PLUS"
	ReservationFlagsDurationMinus ReservationFlagsValue = "DURATION_MINUS"
	ReservationFlagsNoHoldJobsAfterEnd ReservationFlagsValue = "NO_HOLD_JOBS_AFTER_END"
	ReservationFlagsReplaceDown ReservationFlagsValue = "REPLACE_DOWN"
	ReservationFlagsNoPurgeComp ReservationFlagsValue = "NO_PURGE_COMP"
	ReservationFlagsMagnetic ReservationFlagsValue = "MAGNETIC"
	ReservationFlagsNoMagnetic ReservationFlagsValue = "NO_MAGNETIC"
	ReservationFlagsSkip ReservationFlagsValue = "SKIP"
	ReservationFlagsHourly ReservationFlagsValue = "HOURLY"
	ReservationFlagsNoHourly ReservationFlagsValue = "NO_HOURLY"
	ReservationFlagsUserDelete ReservationFlagsValue = "USER_DELETE"
	ReservationFlagsForceStart ReservationFlagsValue = "FORCE_START"
	ReservationFlagsNoUserDelete ReservationFlagsValue = "NO_USER_DELETE"
	ReservationFlagsReoccurring ReservationFlagsValue = "REOCCURRING"
	ReservationFlagsTRESPerNode ReservationFlagsValue = "TRES_PER_NODE"
)
