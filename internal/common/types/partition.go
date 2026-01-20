// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"time"
)

// Partition represents a SLURM partition with common fields across all API versions
type Partition struct {
	Name                 string            `json:"name"`
	AllocNodes           string            `json:"alloc_nodes,omitempty"`
	AllowAccounts        []string          `json:"allow_accounts,omitempty"`
	AllowAllocNodes      string            `json:"allow_alloc_nodes,omitempty"`
	AllowGroups          []string          `json:"allow_groups,omitempty"`
	AllowQoS             []string          `json:"allow_qos,omitempty"`
	DenyAccounts         []string          `json:"deny_accounts,omitempty"`
	DenyQoS              []string          `json:"deny_qos,omitempty"`
	DefaultMemPerCPU     int64             `json:"default_mem_per_cpu,omitempty"`
	DefaultMemPerNode    int64             `json:"default_mem_per_node,omitempty"`
	DefaultTime          int32             `json:"default_time,omitempty"`
	DefaultTimeStr       string            `json:"default_time_str,omitempty"`
	DefMemPerNode        int64             `json:"def_mem_per_node,omitempty"`
	GraceTime            int32             `json:"grace_time,omitempty"`
	MaxCPUsPerNode       int32             `json:"max_cpus_per_node,omitempty"`
	MaxMemPerNode        int64             `json:"max_mem_per_node,omitempty"`
	MaxMemPerCPU         int64             `json:"max_mem_per_cpu,omitempty"`
	MaxNodes             int32             `json:"max_nodes,omitempty"`
	MaxTime              int32             `json:"max_time,omitempty"`
	MaxTimeStr           string            `json:"max_time_str,omitempty"`
	MinNodes             int32             `json:"min_nodes,omitempty"`
	Nodes                string            `json:"nodes,omitempty"`
	NodeCount            int32             `json:"node_count,omitempty"`
	OverTimeLimit        int32             `json:"over_time_limit,omitempty"`
	PreemptMode          []string          `json:"preempt_mode,omitempty"`
	Priority             int32             `json:"priority,omitempty"`
	PriorityJobFactor    int32             `json:"priority_job_factor,omitempty"`
	PriorityTier         int32             `json:"priority_tier,omitempty"`
	QoS                  string            `json:"qos,omitempty"`
	State                PartitionState    `json:"state"`
	StateReason          string            `json:"state_reason,omitempty"`
	TotalCPUs            int32             `json:"total_cpus,omitempty"`
	TotalNodes           int32             `json:"total_nodes,omitempty"`
	TresStr              string            `json:"tres_str,omitempty"`
	BillingWeightStr     string            `json:"billing_weight_str,omitempty"`
	SelectTypeParameters []string          `json:"select_type_parameters,omitempty"`
	JobDefaults          map[string]string `json:"job_defaults,omitempty"`
	ResumeTimeout        int32             `json:"resume_timeout,omitempty"`
	SuspendTime          int32             `json:"suspend_time,omitempty"`
	SuspendTimeout       int32             `json:"suspend_timeout,omitempty"`
	Hidden               bool              `json:"hidden,omitempty"`
	ExclusiveUser        bool              `json:"exclusive_user,omitempty"`
	LLN                  bool              `json:"lln,omitempty"`
	RootOnly             bool              `json:"root_only,omitempty"`
	ReqResv              bool              `json:"req_resv,omitempty"`
	PowerDownOnIdle      bool              `json:"power_down_on_idle,omitempty"`
}

// PartitionState represents the state of a partition
type PartitionState string

const (
	PartitionStateUp       PartitionState = "UP"
	PartitionStateDown     PartitionState = "DOWN"
	PartitionStateDrain    PartitionState = "DRAIN"
	PartitionStateInactive PartitionState = "INACTIVE"
)

// PartitionCreate represents the data needed to create a new partition
type PartitionCreate struct {
	Name                 string            `json:"name"`
	AllocNodes           string            `json:"alloc_nodes,omitempty"`
	AllowAccounts        []string          `json:"allow_accounts,omitempty"`
	AllowAllocNodes      string            `json:"allow_alloc_nodes,omitempty"`
	AllowGroups          []string          `json:"allow_groups,omitempty"`
	AllowQoS             []string          `json:"allow_qos,omitempty"`
	DenyAccounts         []string          `json:"deny_accounts,omitempty"`
	DenyQoS              []string          `json:"deny_qos,omitempty"`
	DefaultMemPerCPU     int64             `json:"default_mem_per_cpu,omitempty"`
	DefaultMemPerNode    int64             `json:"default_mem_per_node,omitempty"`
	DefaultTime          int32             `json:"default_time,omitempty"`
	DefMemPerNode        int64             `json:"def_mem_per_node,omitempty"`
	GraceTime            int32             `json:"grace_time,omitempty"`
	MaxCPUsPerNode       int32             `json:"max_cpus_per_node,omitempty"`
	MaxMemPerNode        int64             `json:"max_mem_per_node,omitempty"`
	MaxMemPerCPU         int64             `json:"max_mem_per_cpu,omitempty"`
	MaxNodes             int32             `json:"max_nodes,omitempty"`
	MaxTime              int32             `json:"max_time,omitempty"`
	MinNodes             int32             `json:"min_nodes,omitempty"`
	Nodes                string            `json:"nodes,omitempty"`
	OverTimeLimit        int32             `json:"over_time_limit,omitempty"`
	PreemptMode          []string          `json:"preempt_mode,omitempty"`
	Priority             int32             `json:"priority,omitempty"`
	PriorityJobFactor    int32             `json:"priority_job_factor,omitempty"`
	PriorityTier         int32             `json:"priority_tier,omitempty"`
	QoS                  string            `json:"qos,omitempty"`
	State                PartitionState    `json:"state,omitempty"`
	TresStr              string            `json:"tres_str,omitempty"`
	BillingWeightStr     string            `json:"billing_weight_str,omitempty"`
	SelectTypeParameters []string          `json:"select_type_parameters,omitempty"`
	JobDefaults          map[string]string `json:"job_defaults,omitempty"`
	ResumeTimeout        int32             `json:"resume_timeout,omitempty"`
	SuspendTime          int32             `json:"suspend_time,omitempty"`
	SuspendTimeout       int32             `json:"suspend_timeout,omitempty"`
	Hidden               bool              `json:"hidden,omitempty"`
	ExclusiveUser        bool              `json:"exclusive_user,omitempty"`
	LLN                  bool              `json:"lln,omitempty"`
	RootOnly             bool              `json:"root_only,omitempty"`
	ReqResv              bool              `json:"req_resv,omitempty"`
	PowerDownOnIdle      bool              `json:"power_down_on_idle,omitempty"`
}

// PartitionUpdate represents the data needed to update a partition
type PartitionUpdate struct {
	AllocNodes           *string           `json:"alloc_nodes,omitempty"`
	AllowAccounts        []string          `json:"allow_accounts,omitempty"`
	AllowAllocNodes      *string           `json:"allow_alloc_nodes,omitempty"`
	AllowGroups          []string          `json:"allow_groups,omitempty"`
	AllowQoS             []string          `json:"allow_qos,omitempty"`
	DenyAccounts         []string          `json:"deny_accounts,omitempty"`
	DenyQoS              []string          `json:"deny_qos,omitempty"`
	DefaultMemPerCPU     *int64            `json:"default_mem_per_cpu,omitempty"`
	DefaultMemPerNode    *int64            `json:"default_mem_per_node,omitempty"`
	DefaultTime          *int32            `json:"default_time,omitempty"`
	DefMemPerNode        *int64            `json:"def_mem_per_node,omitempty"`
	GraceTime            *int32            `json:"grace_time,omitempty"`
	MaxCPUsPerNode       *int32            `json:"max_cpus_per_node,omitempty"`
	MaxMemPerNode        *int64            `json:"max_mem_per_node,omitempty"`
	MaxMemPerCPU         *int64            `json:"max_mem_per_cpu,omitempty"`
	MaxNodes             *int32            `json:"max_nodes,omitempty"`
	MaxTime              *int32            `json:"max_time,omitempty"`
	MinNodes             *int32            `json:"min_nodes,omitempty"`
	Nodes                *string           `json:"nodes,omitempty"`
	OverTimeLimit        *int32            `json:"over_time_limit,omitempty"`
	PreemptMode          []string          `json:"preempt_mode,omitempty"`
	Priority             *int32            `json:"priority,omitempty"`
	PriorityJobFactor    *int32            `json:"priority_job_factor,omitempty"`
	PriorityTier         *int32            `json:"priority_tier,omitempty"`
	QoS                  *string           `json:"qos,omitempty"`
	State                *PartitionState   `json:"state,omitempty"`
	TresStr              *string           `json:"tres_str,omitempty"`
	BillingWeightStr     *string           `json:"billing_weight_str,omitempty"`
	SelectTypeParameters []string          `json:"select_type_parameters,omitempty"`
	JobDefaults          map[string]string `json:"job_defaults,omitempty"`
	ResumeTimeout        *int32            `json:"resume_timeout,omitempty"`
	SuspendTime          *int32            `json:"suspend_time,omitempty"`
	SuspendTimeout       *int32            `json:"suspend_timeout,omitempty"`
	Hidden               *bool             `json:"hidden,omitempty"`
	ExclusiveUser        *bool             `json:"exclusive_user,omitempty"`
	LLN                  *bool             `json:"lln,omitempty"`
	RootOnly             *bool             `json:"root_only,omitempty"`
	ReqResv              *bool             `json:"req_resv,omitempty"`
	PowerDownOnIdle      *bool             `json:"power_down_on_idle,omitempty"`
}

// PartitionCreateResponse represents the response from partition creation
type PartitionCreateResponse struct {
	PartitionName string `json:"partition_name"`
}

// PartitionListOptions represents options for listing partitions
type PartitionListOptions struct {
	Names      []string         `json:"names,omitempty"`
	States     []PartitionState `json:"states,omitempty"`
	UpdateTime *time.Time       `json:"update_time,omitempty"`
	Limit      int              `json:"limit,omitempty"`
	Offset     int              `json:"offset,omitempty"`
}

// PartitionList represents a list of partitions
type PartitionList struct {
	Partitions []Partition `json:"partitions"`
	Total      int         `json:"total"`
}

// PartitionStatistics represents partition usage statistics
type PartitionStatistics struct {
	Name           string    `json:"name"`
	TotalNodes     int32     `json:"total_nodes"`
	AllocatedNodes int32     `json:"allocated_nodes"`
	IdleNodes      int32     `json:"idle_nodes"`
	DownNodes      int32     `json:"down_nodes"`
	TotalCPUs      int32     `json:"total_cpus"`
	AllocatedCPUs  int32     `json:"allocated_cpus"`
	IdleCPUs       int32     `json:"idle_cpus"`
	RunningJobs    int32     `json:"running_jobs"`
	PendingJobs    int32     `json:"pending_jobs"`
	SuspendedJobs  int32     `json:"suspended_jobs"`
	LastUpdateTime time.Time `json:"last_update_time"`
}
