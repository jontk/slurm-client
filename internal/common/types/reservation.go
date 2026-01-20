// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"time"
)

// Reservation represents a SLURM reservation with common fields across all API versions
type Reservation struct {
	Name               string            `json:"name"`
	Accounts           []string          `json:"accounts,omitempty"`
	BurstBuffer        string            `json:"burst_buffer,omitempty"`
	CoreCount          int32             `json:"core_count,omitempty"`
	CoreSpecCount      int32             `json:"core_spec_count,omitempty"`
	EndTime            time.Time         `json:"end_time"`
	Features           []string          `json:"features,omitempty"`
	Flags              []ReservationFlag `json:"flags,omitempty"`
	Groups             []string          `json:"groups,omitempty"`
	Licenses           map[string]int32  `json:"licenses,omitempty"`
	MaxStartDelay      int32             `json:"max_start_delay,omitempty"`
	NodeCount          int32             `json:"node_count,omitempty"`
	NodeList           string            `json:"node_list,omitempty"`
	PartitionName      string            `json:"partition_name,omitempty"`
	StartTime          time.Time         `json:"start_time"`
	State              ReservationState  `json:"state"`
	TRES               map[string]int64  `json:"tres,omitempty"`
	TRESStr            string            `json:"tres_str,omitempty"`
	Users              []string          `json:"users,omitempty"`
	WattsTotal         int64             `json:"watts_total,omitempty"`
	Comment            string            `json:"comment,omitempty"`
	AssocList          []string          `json:"assoc_list,omitempty"`
	CoreIDList         []int32           `json:"core_id_list,omitempty"`
	Duration           int32             `json:"duration,omitempty"`
	PurgeCompletedTime int32             `json:"purge_completed_time,omitempty"`
	ReservationID      int32             `json:"reservation_id,omitempty"`
}

// ReservationState represents the state of a reservation
type ReservationState string

const (
	ReservationStateActive      ReservationState = "ACTIVE"
	ReservationStateInactive    ReservationState = "INACTIVE"
	ReservationStatePending     ReservationState = "PENDING"
	ReservationStateCompleted   ReservationState = "COMPLETED"
	ReservationStateCancelled   ReservationState = "CANCELLED"
	ReservationStateFailed      ReservationState = "FAILED"
	ReservationStateMaintenance ReservationState = "MAINTENANCE"
	ReservationStateOverlap     ReservationState = "OVERLAP"
	ReservationStateError       ReservationState = "ERROR"
)

// ReservationFlag represents flags for a reservation
type ReservationFlag string

const (
	ReservationFlagMaintenance     ReservationFlag = "MAINT"
	ReservationFlagOverlap         ReservationFlag = "OVERLAP"
	ReservationFlagIgnoreJobs      ReservationFlag = "IGNORE_JOBS"
	ReservationFlagDaily           ReservationFlag = "DAILY"
	ReservationFlagWeekly          ReservationFlag = "WEEKLY"
	ReservationFlagAnyNodes        ReservationFlag = "ANY_NODES"
	ReservationFlagStatic          ReservationFlag = "STATIC_ALLOC"
	ReservationFlagPartNodes       ReservationFlag = "PART_NODES"
	ReservationFlagFirstCores      ReservationFlag = "FIRST_CORES"
	ReservationFlagTimeFLoat       ReservationFlag = "TIME_FLOAT"
	ReservationFlagReplace         ReservationFlag = "REPLACE"
	ReservationFlagLicenseOnly     ReservationFlag = "LICENSE_ONLY"
	ReservationFlagNoLicenseOnly   ReservationFlag = "NO_LICENSE_ONLY"
	ReservationFlagPrompt          ReservationFlag = "PROMPT"
	ReservationFlagNoHoldJobsAfter ReservationFlag = "NO_HOLD_JOBS_AFTER"
	ReservationFlagPurgeCompleted  ReservationFlag = "PURGE_COMP"
	ReservationFlagWeekend         ReservationFlag = "WEEKEND"
	ReservationFlagFlexible        ReservationFlag = "FLEX"
	ReservationFlagMagneticCores   ReservationFlag = "MAGNETIC"
	ReservationFlagForce           ReservationFlag = "FORCE"
	ReservationFlagSkipProlog      ReservationFlag = "SKIP_PROLOG"
	ReservationFlagSkipEpilog      ReservationFlag = "SKIP_EPILOG"
)

// ReservationCreate represents the data needed to create a new reservation
type ReservationCreate struct {
	Name               string            `json:"name"`
	Accounts           []string          `json:"accounts,omitempty"`
	BurstBuffer        string            `json:"burst_buffer,omitempty"`
	CoreCount          int32             `json:"core_count,omitempty"`
	CoreSpecCount      int32             `json:"core_spec_count,omitempty"`
	Duration           int32             `json:"duration,omitempty"`
	EndTime            *time.Time        `json:"end_time,omitempty"`
	Features           []string          `json:"features,omitempty"`
	Flags              []ReservationFlag `json:"flags,omitempty"`
	Groups             []string          `json:"groups,omitempty"`
	Licenses           map[string]int32  `json:"licenses,omitempty"`
	MaxStartDelay      int32             `json:"max_start_delay,omitempty"`
	NodeCount          int32             `json:"node_count,omitempty"`
	NodeList           string            `json:"node_list,omitempty"`
	PartitionName      string            `json:"partition_name,omitempty"`
	StartTime          time.Time         `json:"start_time"`
	TRES               map[string]int64  `json:"tres,omitempty"`
	Users              []string          `json:"users,omitempty"`
	Watts              int64             `json:"watts,omitempty"`
	Comment            string            `json:"comment,omitempty"`
	PurgeCompletedTime int32             `json:"purge_completed_time,omitempty"`
}

// ReservationUpdate represents the data needed to update a reservation
type ReservationUpdate struct {
	Accounts           []string          `json:"accounts,omitempty"`
	BurstBuffer        *string           `json:"burst_buffer,omitempty"`
	CoreCount          *int32            `json:"core_count,omitempty"`
	CoreSpecCount      *int32            `json:"core_spec_count,omitempty"`
	Duration           *int32            `json:"duration,omitempty"`
	EndTime            *time.Time        `json:"end_time,omitempty"`
	Features           []string          `json:"features,omitempty"`
	Flags              []ReservationFlag `json:"flags,omitempty"`
	Groups             []string          `json:"groups,omitempty"`
	Licenses           map[string]int32  `json:"licenses,omitempty"`
	MaxStartDelay      *int32            `json:"max_start_delay,omitempty"`
	NodeCount          *int32            `json:"node_count,omitempty"`
	NodeList           *string           `json:"node_list,omitempty"`
	PartitionName      *string           `json:"partition_name,omitempty"`
	StartTime          *time.Time        `json:"start_time,omitempty"`
	TRES               map[string]int64  `json:"tres,omitempty"`
	Users              []string          `json:"users,omitempty"`
	Watts              *int64            `json:"watts,omitempty"`
	Comment            *string           `json:"comment,omitempty"`
	PurgeCompletedTime *int32            `json:"purge_completed_time,omitempty"`
}

// ReservationCreateResponse represents the response from reservation creation
type ReservationCreateResponse struct {
	ReservationName string `json:"reservation_name"`
}

// ReservationListOptions represents options for listing reservations
type ReservationListOptions struct {
	Names      []string           `json:"names,omitempty"`
	States     []ReservationState `json:"states,omitempty"`
	Accounts   []string           `json:"accounts,omitempty"`
	Users      []string           `json:"users,omitempty"`
	Partitions []string           `json:"partitions,omitempty"`
	StartTime  *time.Time         `json:"start_time,omitempty"`
	EndTime    *time.Time         `json:"end_time,omitempty"`
	UpdateTime *time.Time         `json:"update_time,omitempty"`
	Limit      int                `json:"limit,omitempty"`
	Offset     int                `json:"offset,omitempty"`
}

// ReservationList represents a list of reservations
type ReservationList struct {
	Reservations []Reservation `json:"reservations"`
	Total        int           `json:"total"`
}

// ReservationUsage represents usage statistics for a reservation
type ReservationUsage struct {
	ReservationName string             `json:"reservation_name"`
	StartTime       time.Time          `json:"start_time"`
	EndTime         time.Time          `json:"end_time"`
	AllocatedNodes  int32              `json:"allocated_nodes"`
	AllocatedCores  int32              `json:"allocated_cores"`
	IdleNodes       int32              `json:"idle_nodes"`
	IdleCores       int32              `json:"idle_cores"`
	JobCount        int32              `json:"job_count"`
	UserCount       int32              `json:"user_count"`
	TRESUsage       map[string]float64 `json:"tres_usage,omitempty"`
	UtilizationRate float64            `json:"utilization_rate"`
}
