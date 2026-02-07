// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// JobResNode represents a SLURM JobResNode.
type JobResNode struct {
	CPUs *JobResNodeCPUs `json:"cpus,omitempty"`
	Index int32 `json:"index"` // Node index
	Memory *JobResNodeMemory `json:"memory,omitempty"`
	Name string `json:"name"` // Node name
	Sockets []JobResSocket `json:"sockets"` // Socket allocations in node
}


// JobResNodeCPUs is a nested type within its parent.
type JobResNodeCPUs struct {
	Count *int32 `json:"count,omitempty"` // Total number of CPUs assigned to job
	Used *int32 `json:"used,omitempty"` // Total number of CPUs used by job
}


// JobResNodeMemory is a nested type within its parent.
type JobResNodeMemory struct {
	Allocated *int64 `json:"allocated,omitempty"` // Total memory (MiB) allocated to job
	Used *int64 `json:"used,omitempty"` // Total memory (MiB) used by job
}
