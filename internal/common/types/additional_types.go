// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package types

import "time"

// JobCancelFlags represents flags for job cancellation
type JobCancelFlags struct {
	// BatchJob cancels only the batch job
	BatchJob bool `json:"batch_job,omitempty"`
	// FullJob cancels the entire job array
	FullJob bool `json:"full_job,omitempty"`
	// HurdleJob cancels hurdle jobs
	HurdleJob bool `json:"hurdle_job,omitempty"`
	// IgnoreFederation ignores federation in job cancellation
	IgnoreFederation bool `json:"ignore_federation,omitempty"`
	// Signal to send to the job (default: SIGKILL)
	Signal string `json:"signal,omitempty"`
	// WaitTime to wait before forcefully terminating the job (in seconds)
	WaitTime int32 `json:"wait_time,omitempty"`
}

// JobSubmitOptions represents options for job submission
type JobSubmitOptions struct {
	// HoldJob holds the job immediately after submission
	HoldJob bool `json:"hold_job,omitempty"`
	// TestOnly validates the job submission without actually submitting
	TestOnly bool `json:"test_only,omitempty"`
	// Parsable returns job ID in a parsable format
	Parsable bool `json:"parsable,omitempty"`
	// WaitForCompletion waits for the job to complete
	WaitForCompletion bool `json:"wait_for_completion,omitempty"`
	// Verbose provides verbose output
	Verbose bool `json:"verbose,omitempty"`
}

// JobWatchOptions represents options for watching job events
type JobWatchOptions struct {
	// JobID is the specific job to watch (0 means all jobs)
	JobID int32 `json:"job_id,omitempty"`
	// StartTime to begin watching from
	StartTime *time.Time `json:"start_time,omitempty"`
	// EventTypes to filter (e.g., "submit", "start", "end", "fail")
	EventTypes []string `json:"event_types,omitempty"`
	// MaxEvents maximum number of events to return
	MaxEvents int32 `json:"max_events,omitempty"`
	// IncludeSteps includes job step events
	IncludeSteps bool `json:"include_steps,omitempty"`
}

// JobEvent represents a job state change event
type JobEvent struct {
	// EventTime when the event occurred
	EventTime time.Time `json:"event_time"`
	// EventType type of event (submit, start, end, fail, etc.)
	EventType string `json:"event_type"`
	// JobID of the job
	JobID int32 `json:"job_id"`
	// JobName of the job
	JobName string `json:"job_name,omitempty"`
	// UserName who owns the job
	UserName string `json:"user_name,omitempty"`
	// PreviousState before the event
	PreviousState JobState `json:"previous_state,omitempty"`
	// NewState after the event
	NewState JobState `json:"new_state"`
	// Reason for the state change
	Reason string `json:"reason,omitempty"`
	// NodeList affected by the event
	NodeList string `json:"node_list,omitempty"`
	// ExitCode if job ended
	ExitCode int32 `json:"exit_code,omitempty"`
}

// AccountUserOptions represents options for account-user operations
type AccountUserOptions struct {
	// AccountName to filter by
	AccountName string `json:"account_name,omitempty"`
	// UserNames to include
	UserNames []string `json:"user_names,omitempty"`
	// WithAssociations includes user associations
	WithAssociations bool `json:"with_associations,omitempty"`
	// WithCoordinators includes coordinator information
	WithCoordinators bool `json:"with_coordinators,omitempty"`
	// OnlyCoordinators returns only users who are coordinators
	OnlyCoordinators bool `json:"only_coordinators,omitempty"`
}

// UserAccountOptions represents options for user-account operations
type UserAccountOptions struct {
	// UserName to filter by
	UserName string `json:"user_name,omitempty"`
	// AccountNames to include
	AccountNames []string `json:"account_names,omitempty"`
	// WithLimits includes resource limits
	WithLimits bool `json:"with_limits,omitempty"`
	// WithUsage includes usage information
	WithUsage bool `json:"with_usage,omitempty"`
	// DefaultOnly returns only default account
	DefaultOnly bool `json:"default_only,omitempty"`
}

// AssociationLimits represents resource limits for an association
type AssociationLimits struct {
	// MaxJobs maximum number of jobs
	MaxJobs int32 `json:"max_jobs,omitempty"`
	// MaxJobsAccrue maximum number of jobs that can accrue priority
	MaxJobsAccrue int32 `json:"max_jobs_accrue,omitempty"`
	// MaxSubmitJobs maximum number of submitted jobs
	MaxSubmitJobs int32 `json:"max_submit_jobs,omitempty"`
	// MaxWallTime maximum wall time per job (in minutes)
	MaxWallTime int32 `json:"max_wall_time,omitempty"`
	// MaxCPUTime maximum CPU time (in minutes)
	MaxCPUTime int32 `json:"max_cpu_time,omitempty"`
	// MaxNodes maximum number of nodes per job
	MaxNodes int32 `json:"max_nodes,omitempty"`
	// MaxCPUs maximum number of CPUs per job
	MaxCPUs int32 `json:"max_cpus,omitempty"`
	// MaxMemory maximum memory per job (in MB)
	MaxMemory int64 `json:"max_memory,omitempty"`
	// MinPriorityThreshold minimum priority threshold
	MinPriorityThreshold int32 `json:"min_priority_threshold,omitempty"`
	// GrpJobs maximum number of running jobs in aggregate
	GrpJobs int32 `json:"grp_jobs,omitempty"`
	// GrpJobsAccrue maximum number of jobs that can accrue priority in aggregate
	GrpJobsAccrue int32 `json:"grp_jobs_accrue,omitempty"`
	// GrpNodes maximum number of nodes in aggregate
	GrpNodes int32 `json:"grp_nodes,omitempty"`
	// GrpCPUs maximum number of CPUs in aggregate
	GrpCPUs int32 `json:"grp_cpus,omitempty"`
	// GrpMemory maximum memory in aggregate (in MB)
	GrpMemory int64 `json:"grp_memory,omitempty"`
	// GrpSubmitJobs maximum number of submitted jobs in aggregate
	GrpSubmitJobs int32 `json:"grp_submit_jobs,omitempty"`
	// GrpWallTime maximum wall time in aggregate (in minutes)
	GrpWallTime int32 `json:"grp_wall_time,omitempty"`
	// GrpCPUTime maximum CPU time in aggregate (in minutes)
	GrpCPUTime int32 `json:"grp_cpu_time,omitempty"`
	// GrpCPURunMins maximum CPU running minutes
	GrpCPURunMins int64 `json:"grp_cpu_run_mins,omitempty"`
	// GrpTRES group TRES limits
	GrpTRES map[string]int64 `json:"grp_tres,omitempty"`
	// GrpTRESMins group TRES minutes limits
	GrpTRESMins map[string]int64 `json:"grp_tres_mins,omitempty"`
	// GrpTRESRunMins group TRES running minutes limits
	GrpTRESRunMins map[string]int64 `json:"grp_tres_run_mins,omitempty"`
	// MaxTRES maximum TRES limits
	MaxTRES map[string]int64 `json:"max_tres,omitempty"`
	// MaxTRESPerNode maximum TRES per node limits
	MaxTRESPerNode map[string]int64 `json:"max_tres_per_node,omitempty"`
	// MaxTRESMins maximum TRES minutes
	MaxTRESMins map[string]int64 `json:"max_tres_mins,omitempty"`
	// MinTRES minimum TRES requirements
	MinTRES map[string]int64 `json:"min_tres,omitempty"`
}

// NodeWatchOptions represents options for watching node events
type NodeWatchOptions struct {
	// NodeNames is the list of specific nodes to watch (empty means all nodes)
	NodeNames []string `json:"node_names,omitempty"`
	// States to filter by (e.g., "idle", "allocated", "down")
	States []NodeState `json:"states,omitempty"`
	// Partitions to filter by
	Partitions []string `json:"partitions,omitempty"`
	// MaxEvents maximum number of events to return
	MaxEvents int32 `json:"max_events,omitempty"`
}

// NodeEvent represents a node state change event
type NodeEvent struct {
	// EventTime when the event occurred
	EventTime time.Time `json:"event_time"`
	// EventType type of event (state_change, drain, resume, etc.)
	EventType string `json:"event_type"`
	// NodeName of the node
	NodeName string `json:"node_name"`
	// PreviousState before the event
	PreviousState NodeState `json:"previous_state,omitempty"`
	// NewState after the event
	NewState NodeState `json:"new_state"`
	// Reason for the state change
	Reason string `json:"reason,omitempty"`
	// Partitions affected by the event
	Partitions []string `json:"partitions,omitempty"`
}

// Type aliases for compatibility
type JobWatchEvent = JobEvent
type NodeWatchEvent = NodeEvent
type ResourceLimits = AssociationLimits
