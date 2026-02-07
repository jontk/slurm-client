// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Package types provides common type definitions for SLURM entities.
// Core entity types (Job, Node, User, etc.) are generated in *.gen.go files.
// This file contains operation types (Create, Update, List, etc.).
package api

import (
	"time"
)

// ResourceRequests represents the resource requirements for a job
type ResourceRequests struct {
	Memory        int64 `json:"memory,omitempty"`
	MemoryPerCPU  int64 `json:"memory_per_cpu,omitempty"`
	MemoryPerGPU  int64 `json:"memory_per_gpu,omitempty"`
	TemporaryDisk int64 `json:"tmp_disk,omitempty"`
	CPUsPerTask   int32 `json:"cpus_per_task,omitempty"`
	TasksPerNode  int32 `json:"tasks_per_node,omitempty"`
	TasksPerCore  int32 `json:"tasks_per_core,omitempty"`
	Threads       int32 `json:"threads_per_core,omitempty"`
}

// JobDependency represents a job dependency (helper type for user convenience)
// Note: When using the generated JobCreate, specify dependencies as a string
// in SLURM dependency format (e.g., "afterok:123:456")
type JobDependency struct {
	Type   string  `json:"type"`
	JobIDs []int32 `json:"job_ids,omitempty"`
	State  string  `json:"state,omitempty"`
}

// NOTE: JobCreate is now generated in jobcreate.gen.go from the OpenAPI spec.
// The generated type uses SLURM API types directly. Key differences from the old manual type:
// - Environment: []string (format: "KEY=VALUE") instead of map[string]string
// - Dependencies: *string (SLURM format: "afterok:123") instead of []JobDependency
// - Features/Constraints: *string instead of []string
// - TimeLimit/Priority: *uint32 wrapped in NoValStruct semantics

// JobUpdate is an alias for JobCreate since SLURM uses the same job_desc_msg
// for both create and update operations.
type JobUpdate = JobCreate

// JobSubmitResponse represents the response from job submission
type JobSubmitResponse struct {
	JobId            int32    `json:"job_id"`          // Matches OpenAPI: JobId *int32
	StepId           string   `json:"step_id,omitempty"` // Matches OpenAPI casing
	JobSubmitUserMsg string   `json:"job_submit_user_msg,omitempty"`
	Error            []string `json:"error,omitempty"`
	Warning          []string `json:"warning,omitempty"`
}

// JobCancelRequest represents the request to cancel a job
type JobCancelRequest struct {
	Signal    string `json:"signal,omitempty"`
	Message   string `json:"message,omitempty"`
	Account   string `json:"account,omitempty"`
	Name      string `json:"name,omitempty"`
	Partition string `json:"partition,omitempty"`
	QoS       string `json:"qos,omitempty"`
	State     string `json:"state,omitempty"`
	UserID    int32  `json:"user_id,omitempty"`
	UserName  string `json:"user_name,omitempty"`
	WaitTime  int32  `json:"wait_time,omitempty"`
}

// JobListOptions represents options for listing jobs
type JobListOptions struct {
	Accounts     []string   `json:"accounts,omitempty"`
	Users        []string   `json:"users,omitempty"`
	States       []JobState `json:"states,omitempty"`
	Partitions   []string   `json:"partitions,omitempty"`
	QoS          []string   `json:"qos,omitempty"`
	JobIDs       []int32    `json:"job_ids,omitempty"`
	JobNames     []string   `json:"job_names,omitempty"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`

	// Limit specifies the maximum number of jobs to return.
	// WARNING: Due to SLURM REST API limitations, this is CLIENT-SIDE pagination.
	// The full job list is fetched from the server, then sliced. For large clusters
	// (100K+ jobs), consider using filtering options (States, Accounts, Partitions, etc.)
	// to reduce the dataset before pagination.
	Limit int `json:"limit,omitempty"`

	// Offset specifies the number of jobs to skip before returning results.
	// WARNING: This is CLIENT-SIDE pagination - see Limit field documentation.
	Offset int `json:"offset,omitempty"`

	IncludeSteps bool `json:"include_steps,omitempty"`
}

// JobList represents a list of jobs
type JobList struct {
	Jobs  []Job `json:"jobs"`
	Total int   `json:"total"`
}

// JobSignalRequest represents a request to signal a job
type JobSignalRequest struct {
	Signal string `json:"signal"`
	JobId  int32  `json:"job_id"`  // Matches OpenAPI casing
	StepId string `json:"step_id,omitempty"`
}

// JobHoldRequest represents a request to hold/release a job
type JobHoldRequest struct {
	JobId    int32 `json:"job_id"`  // Matches OpenAPI casing
	Hold     bool  `json:"hold"`
	Priority int32 `json:"priority,omitempty"`
}

// JobNotifyRequest represents a request to notify a job
type JobNotifyRequest struct {
	JobId   int32  `json:"job_id"`  // Matches OpenAPI casing
	Message string `json:"message"`
}

// JobAllocateRequest represents a request to allocate resources for a job
type JobAllocateRequest struct {
	// Job specification
	Name      string `json:"name,omitempty"`
	Account   string `json:"account,omitempty"`
	Partition string `json:"partition,omitempty"`
	QoS       string `json:"qos,omitempty"`

	// Resource requirements
	Nodes  string `json:"nodes,omitempty"`  // Number or range of nodes
	Cpus   int32  `json:"cpus,omitempty"`   // Number of CPUs (matches OpenAPI casing)
	Memory string `json:"memory,omitempty"` // Memory requirement
	Gpus   string `json:"gpus,omitempty"`   // GPU requirement (matches OpenAPI casing)

	// Time limits
	TimeLimit int32 `json:"time_limit,omitempty"` // Time limit in minutes

	// Environment and execution
	Environment map[string]string `json:"environment,omitempty"`
	WorkingDir  string            `json:"working_directory,omitempty"`
	Command     []string          `json:"command,omitempty"` // Command to run

	// Advanced options
	Exclusive   bool     `json:"exclusive,omitempty"`
	Features    []string `json:"features,omitempty"`
	Constraints string   `json:"constraints,omitempty"`

	// Output handling
	StdOut string `json:"stdout,omitempty"`
	StdErr string `json:"stderr,omitempty"`
	StdIn  string `json:"stdin,omitempty"`
}

// JobAllocateResponse represents the response from a job allocation request
type JobAllocateResponse struct {
	JobId   int32  `json:"job_id"`  // Matches OpenAPI casing
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	// Allocated resources
	Nodes  []string `json:"nodes,omitempty"`
	Cpus   int32    `json:"cpus_allocated,omitempty"` // Matches OpenAPI casing
	Memory int64    `json:"memory_allocated,omitempty"`
	Gpus   int32    `json:"gpus_allocated,omitempty"` // Matches OpenAPI casing

	// Timing information
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	TimeLimit int32      `json:"time_limit,omitempty"`

	// Connection information for interactive jobs
	ConnectionInfo map[string]interface{} `json:"connection_info,omitempty"`

	Meta map[string]interface{} `json:"meta,omitempty"`
}
