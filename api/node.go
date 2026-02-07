// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Package types provides common type definitions for SLURM entities.
// Core entity types (Job, Node, User, etc.) are generated in *.gen.go files.
// This file contains operation types (Create, Update, List, etc.).
package api

import (
	"time"
)

// NodeUpdate is generated in nodeupdate.gen.go

// NodeListOptions represents options for listing nodes
type NodeListOptions struct {
	Names      []string    `json:"names,omitempty"`
	States     []NodeState `json:"states,omitempty"`
	Partitions []string   `json:"partitions,omitempty"`
	UpdateTime *time.Time `json:"update_time,omitempty"`
	Reasons    []string   `json:"reasons,omitempty"`

	// Limit specifies the maximum number of nodes to return.
	// WARNING: Due to SLURM REST API limitations, this is CLIENT-SIDE pagination.
	// The full node list is fetched from the server, then sliced. For large clusters
	// (10K+ nodes), consider using filtering options (States, Partitions, Names, etc.)
	// to reduce the dataset before pagination.
	Limit int `json:"limit,omitempty"`

	// Offset specifies the number of nodes to skip before returning results.
	// WARNING: This is CLIENT-SIDE pagination - see Limit field documentation.
	Offset int `json:"offset,omitempty"`
}

// NodeList represents a list of nodes
type NodeList struct {
	Nodes []Node `json:"nodes"`
	Total int    `json:"total"`
}

// NodeAllocation represents the current allocation status of a node
type NodeAllocation struct {
	Name        string    `json:"name"`
	JobId       int32     `json:"job_id,omitempty"`
	JobName     string    `json:"job_name,omitempty"`
	UserName    string    `json:"user_name,omitempty"`
	AllocCpus         int32     `json:"alloc_cpus"`
	AllocMemory int64     `json:"alloc_memory"`
	AllocGres   string    `json:"alloc_gres,omitempty"`
	StartTime   time.Time `json:"start_time,omitempty"`
	EndTime     time.Time `json:"end_time,omitempty"`
}

// NodeMaintenanceRequest represents a request to put nodes into maintenance
type NodeMaintenanceRequest struct {
	Nodes         []string   `json:"nodes"`
	Reason        string     `json:"reason"`
	StartTime     *time.Time `json:"start_time,omitempty"`
	EndTime       *time.Time `json:"end_time,omitempty"`
	FixedDuration int32      `json:"fixed_duration,omitempty"`
}

// NodePowerRequest represents a request to change node power state
type NodePowerRequest struct {
	Nodes        []string       `json:"nodes"`
	PowerState   NodePowerState `json:"power_state"`
	Asynchronous bool           `json:"asynchronous,omitempty"`
	Force        bool           `json:"force,omitempty"`
}

// NodePowerState represents node power states
type NodePowerState string

const (
	NodePowerDown NodePowerState = "POWER_DOWN"
	NodePowerUp   NodePowerState = "POWER_UP"
	NodePowerSave NodePowerState = "POWER_SAVE"
)
