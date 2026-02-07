// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

import "time"

// Partition represents a SLURM Partition.
type Partition struct {
	Accounts *PartitionAccounts `json:"accounts,omitempty"`
	Alternate *string `json:"alternate,omitempty"` // Alternate - Partition name of alternate partition to be used if the state of...
	Cluster *string `json:"cluster,omitempty"` // Cluster name
	CPUs *PartitionCPUs `json:"cpus,omitempty"`
	Defaults *PartitionDefaults `json:"defaults,omitempty"`
	GraceTime *int32 `json:"grace_time,omitempty"` // GraceTime - Grace time in seconds to be extended to a job which has been...
	Groups *PartitionGroups `json:"groups,omitempty"`
	Maximums *PartitionMaximums `json:"maximums,omitempty"`
	Minimums *PartitionMinimums `json:"minimums,omitempty"`
	Name *string `json:"name,omitempty"` // PartitionName - Name by which the partition may be referenced
	NodeSets *string `json:"node_sets,omitempty"` // NodeSets - Comma-separated list of nodesets which are associated with this...
	Nodes *PartitionNodes `json:"nodes,omitempty"`
	Partition *PartitionPartition `json:"partition,omitempty"`
	Priority *PartitionPriority `json:"priority,omitempty"`
	QoS *PartitionQoS `json:"qos,omitempty"`
	SelectType []SelectTypeValue `json:"select_type,omitempty"` // Scheduler consumable resource selection type
	SuspendTime time.Time `json:"suspend_time,omitempty"` // SuspendTime - Nodes which remain idle or down for this number of seconds will...
	Timeouts *PartitionTimeouts `json:"timeouts,omitempty"`
	Topology *string `json:"topology,omitempty"` // Topology - Name of the topology, defined in topology.yaml, used by jobs in this...
	TRES *PartitionTRES `json:"tres,omitempty"`
}


// PartitionAccounts is a nested type within its parent.
type PartitionAccounts struct {
	Allowed *string `json:"allowed,omitempty"` // AllowAccounts - Comma-separated list of accounts which may execute jobs in the...
	Deny *string `json:"deny,omitempty"` // DenyAccounts - Comma-separated list of accounts which may not execute jobs in...
}


// PartitionCPUs is a nested type within its parent.
type PartitionCPUs struct {
	TaskBinding *int32 `json:"task_binding,omitempty"` // CpuBind - Default method controlling how tasks are bound to allocated resources
	Total *int32 `json:"total,omitempty"` // TotalCPUs - Number of CPUs available in this partition
}


// PartitionDefaults is a nested type within its parent.
type PartitionDefaults struct {
	Job *string `json:"job,omitempty"` // JobDefaults - Comma-separated list of job default values (this field is only...
	MemoryPerCPU *int64 `json:"memory_per_cpu,omitempty"` // Raw value for DefMemPerCPU or DefMemPerNode
	PartitionMemoryPerCPU *uint64 `json:"partition_memory_per_cpu,omitempty"` // DefMemPerCPU - Default real memory size available per allocated CPU in...
	PartitionMemoryPerNode *uint64 `json:"partition_memory_per_node,omitempty"` // DefMemPerNode - Default real memory size available per allocated node in...
	Time *uint32 `json:"time,omitempty"` // DefaultTime - Run time limit in minutes used for jobs that don't specify a...
}


// PartitionGroups is a nested type within its parent.
type PartitionGroups struct {
	Allowed *string `json:"allowed,omitempty"` // AllowGroups - Comma-separated list of group names which may execute jobs in...
}


// PartitionMaximums is a nested type within its parent.
type PartitionMaximums struct {
	CPUsPerNode *uint32 `json:"cpus_per_node,omitempty"` // MaxCPUsPerNode - Maximum number of CPUs on any node available to all jobs from...
	CPUsPerSocket *uint32 `json:"cpus_per_socket,omitempty"` // MaxCPUsPerSocket - Maximum number of CPUs on any node available on the all jobs...
	MemoryPerCPU *int64 `json:"memory_per_cpu,omitempty"` // Raw value for MaxMemPerCPU or MaxMemPerNode
	Nodes *uint32 `json:"nodes,omitempty"` // MaxNodes - Maximum count of nodes which may be allocated to any single job (32...
	OverTimeLimit *uint16 `json:"over_time_limit,omitempty"` // OverTimeLimit - Number of minutes by which a job can exceed its time limit...
	Oversubscribe *PartitionMaximumsOversubscribe `json:"oversubscribe,omitempty"`
	PartitionMemoryPerCPU *uint64 `json:"partition_memory_per_cpu,omitempty"` // MaxMemPerCPU - Maximum real memory size available per allocated CPU in...
	PartitionMemoryPerNode *uint64 `json:"partition_memory_per_node,omitempty"` // MaxMemPerNode - Maximum real memory size available per allocated node in a job...
	Shares *int32 `json:"shares,omitempty"` // OverSubscribe - Controls the ability of the partition to execute more than one...
	Time *uint32 `json:"time,omitempty"` // MaxTime - Maximum run time limit for jobs (32 bit integer number with flags)
}


// PartitionMaximumsOversubscribe is a nested type within its parent.
type PartitionMaximumsOversubscribe struct {
	Flags []PartitionMaximumsOversubscribeFlagsValue `json:"flags,omitempty"` // Flags applicable to the OverSubscribe setting
	Jobs *int32 `json:"jobs,omitempty"` // Maximum number of jobs allowed to oversubscribe resources
}


// PartitionMaximumsOversubscribeFlagsValue represents possible values for PartitionMaximumsOversubscribeFlags field.
type PartitionMaximumsOversubscribeFlagsValue string

// PartitionMaximumsOversubscribeFlagsValue constants.
const (
	PartitionMaximumsOversubscribeFlagsForce PartitionMaximumsOversubscribeFlagsValue = "force"
)


// PartitionMinimums is a nested type within its parent.
type PartitionMinimums struct {
	Nodes *int32 `json:"nodes,omitempty"` // MinNodes - Minimum count of nodes which may be allocated to any single job
}


// PartitionNodes is a nested type within its parent.
type PartitionNodes struct {
	AllowedAllocation *string `json:"allowed_allocation,omitempty"` // AllocNodes - Comma-separated list of nodes from which users can submit jobs in...
	Configured *string `json:"configured,omitempty"` // Nodes - Comma-separated list of nodes which are associated with this partition
	Total *int32 `json:"total,omitempty"` // TotalNodes - Number of nodes available in this partition
}


// PartitionPartition is a nested type within its parent.
type PartitionPartition struct {
	State []StateValue `json:"state,omitempty"` // Current state(s)
}


// StateValue represents possible values for State field.
type StateValue string

// StateValue constants.
const (
	StateInactive StateValue = "INACTIVE"
	StateUnknown StateValue = "UNKNOWN"
	StateUp StateValue = "UP"
	StateDown StateValue = "DOWN"
	StateDrain StateValue = "DRAIN"
)


// PartitionPriority is a nested type within its parent.
type PartitionPriority struct {
	JobFactor *int32 `json:"job_factor,omitempty"` // PriorityJobFactor - Partition factor used by priority/multifactor plugin in...
	Tier *int32 `json:"tier,omitempty"` // PriorityTier - Controls the order in which the scheduler evaluates jobs from...
}


// PartitionQoS is a nested type within its parent.
type PartitionQoS struct {
	Allowed *string `json:"allowed,omitempty"` // AllowQOS - Comma-separated list of Qos which may execute jobs in the partition
	Assigned *string `json:"assigned,omitempty"` // QOS - QOS name containing limits that will apply to all jobs in this partition
	Deny *string `json:"deny,omitempty"` // DenyQOS - Comma-separated list of Qos which may not execute jobs in the...
}


// PartitionTimeouts is a nested type within its parent.
type PartitionTimeouts struct {
	Resume *uint16 `json:"resume,omitempty"` // ResumeTimeout - Resumed nodes which fail to respond in this time frame will be...
	Suspend *uint16 `json:"suspend,omitempty"` // SuspendTimeout - Maximum time permitted (in seconds) between when a node...
}


// SelectTypeValue represents possible values for SelectType field.
type SelectTypeValue string

// SelectTypeValue constants.
const (
	SelectTypeCPU SelectTypeValue = "CPU"
	SelectTypeSocket SelectTypeValue = "SOCKET"
	SelectTypeCore SelectTypeValue = "CORE"
	SelectTypeBoard SelectTypeValue = "BOARD"
	SelectTypeMemory SelectTypeValue = "MEMORY"
	SelectTypeOneTaskPerCore SelectTypeValue = "ONE_TASK_PER_CORE"
	SelectTypePackNodes SelectTypeValue = "PACK_NODES"
	SelectTypeCoreDefaultDistBlock SelectTypeValue = "CORE_DEFAULT_DIST_BLOCK"
	SelectTypeLln SelectTypeValue = "LLN"
	SelectTypeLinear SelectTypeValue = "LINEAR"
)


// PartitionTRES is a nested type within its parent.
type PartitionTRES struct {
	BillingWeights *string `json:"billing_weights,omitempty"` // TRESBillingWeights - Billing weights of each tracked TRES type that will be...
	Configured *string `json:"configured,omitempty"` // TRES - Number of each applicable TRES type available in this partition
}
