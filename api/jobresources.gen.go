// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// JobResources represents a SLURM JobResources.
type JobResources struct {
	CPUs int32 `json:"cpus"` // Number of allocated CPUs
	Nodes *JobResourcesNodes `json:"nodes,omitempty"`
	SelectType []SelectTypeValue `json:"select_type"` // Scheduler consumable resource selection type
	ThreadsPerCore uint16 `json:"threads_per_core"` // Number of processor threads per CPU core (16 bit integer number with flags)
}


// JobResourcesNodes is a nested type within its parent.
type JobResourcesNodes struct {
	Allocation []JobResNode `json:"allocation,omitempty"` // Allocated node resources (Job resources for a node)
	Count *int32 `json:"count,omitempty"` // Number of allocated nodes
	List *string `json:"list,omitempty"` // Node(s) allocated to the job
	SelectType []JobResourcesNodesSelectTypeValue `json:"select_type,omitempty"` // Node scheduling selection method
	Whole *bool `json:"whole,omitempty"` // Whether whole nodes were allocated
}


// JobResourcesNodesSelectTypeValue represents possible values for JobResourcesNodesSelectType field.
type JobResourcesNodesSelectTypeValue string

// JobResourcesNodesSelectTypeValue constants.
const (
	JobResourcesNodesSelectTypeAvailable JobResourcesNodesSelectTypeValue = "AVAILABLE"
	JobResourcesNodesSelectTypeOneRow JobResourcesNodesSelectTypeValue = "ONE_ROW"
	JobResourcesNodesSelectTypeReserved JobResourcesNodesSelectTypeValue = "RESERVED"
)
