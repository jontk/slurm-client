// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Package types provides common type definitions for SLURM entities.
// Core entity types (Job, Node, User, etc.) are generated in *.gen.go files.
// This file contains operation types (Create, Update, List, etc.).
package api

import (
	"time"
)

// ReservationState represents the state of a reservation (custom, not in OpenAPI)
type ReservationState string

// ReservationState constants
const (
	ReservationStateActive   ReservationState = "ACTIVE"
	ReservationStateInactive ReservationState = "INACTIVE"
)

// ReservationFlag alias for generated ReservationFlagsValue
type ReservationFlag = ReservationFlagsValue

// Backward compatibility constants mapping to generated values
const (
	ReservationFlagMaintenance = ReservationFlagsMaint
	ReservationFlagOverlap     = ReservationFlagsOverlap
	ReservationFlagIgnoreJobs  = ReservationFlagsIgnoreJobs
	ReservationFlagDaily       = ReservationFlagsDaily
	ReservationFlagWeekly      = ReservationFlagsWeekly
	ReservationFlagAnyNodes    = ReservationFlagsAnyNodes
	ReservationFlagStatic      = ReservationFlagsStatic
	ReservationFlagPartNodes   = ReservationFlagsPartNodes
)

// ReservationCreate is generated in reservationcreate.gen.go

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
	Partition      *string           `json:"partition_name,omitempty"`
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
	StartTime  *time.Time `json:"start_time,omitempty"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	UpdateTime *time.Time `json:"update_time,omitempty"`

	// Limit specifies the maximum number of reservations to return.
	// WARNING: Due to SLURM REST API limitations, this is CLIENT-SIDE pagination.
	// The full reservation list is fetched from the server, then sliced. Consider using
	// filtering options (Names, States, Accounts, Users, Partitions) to reduce the dataset.
	Limit int `json:"limit,omitempty"`

	// Offset specifies the number of reservations to skip before returning results.
	// WARNING: This is CLIENT-SIDE pagination - see Limit field documentation.
	Offset int `json:"offset,omitempty"`
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
