// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

import "time"

// ReservationCreate represents a SLURM ReservationCreate.
type ReservationCreate struct {
	Accounts []string `json:"accounts,omitempty"` // List of permitted accounts
	BurstBuffer *string `json:"burst_buffer,omitempty"` // BurstBuffer
	Comment *string `json:"comment,omitempty"` // Arbitrary string
	CoreCount *uint32 `json:"core_count,omitempty"` // Number of cores to reserve (32 bit integer number with flags)
	Duration *uint32 `json:"duration,omitempty"` // The length of a reservation in minutes (32 bit integer number with flags)
	EndTime time.Time `json:"end_time,omitempty"` // EndTime (UNIX timestamp) (UNIX timestamp or time string recognized by Slurm...
	Features *string `json:"features,omitempty"` // Requested node features. Multiple values may be "&" separated if all features...
	Flags []FlagsValue `json:"flags,omitempty"` // Flags associated with this reservation. Note, to remove flags use "NO_"...
	Groups []string `json:"groups,omitempty"` // List of groups permitted to use the reservation. This is mutually exclusive...
	Licenses []string `json:"licenses,omitempty"` // List of license names
	MaxStartDelay *uint32 `json:"max_start_delay,omitempty"` // MaxStartDelay in seconds (32 bit integer number with flags)
	Name *string `json:"name,omitempty"` // ReservationName
	NodeCount *uint32 `json:"node_count,omitempty"` // NodeCnt (32 bit integer number with flags)
	NodeList []string `json:"node_list,omitempty"` // The nodes to be reserved. Multiple node names may be specified using simple...
	Partition *string `json:"partition,omitempty"` // Partition used to reserve nodes from. This will attempt to allocate all nodes...
	PurgeCompleted *ReservationCreatePurgeCompleted `json:"purge_completed,omitempty"`
	StartTime time.Time `json:"start_time,omitempty"` // StartTime (UNIX timestamp) (UNIX timestamp or time string recognized by Slurm...
	TRES []TRES `json:"tres,omitempty"` // List of trackable resources
	Users []string `json:"users,omitempty"` // List of permitted users
}


// ReservationCreatePurgeCompleted is a nested type within its parent.
type ReservationCreatePurgeCompleted struct {
	Time *uint32 `json:"time,omitempty"` // If PURGE_COMP flag is set, the number of seconds this reservation will sit idle...
}
