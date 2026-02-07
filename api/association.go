// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Package types provides common type definitions for SLURM entities.
// Core entity types (Job, Node, User, etc.) are generated in *.gen.go files.
// This file contains operation types (Create, Update, List, etc.).
package api

import (
	"time"
)

// AssociationCreate represents the data needed to create a new association
type AssociationCreate struct {
	Account              string `json:"account_name"`
	Cluster              string           `json:"cluster"`
	User                 string `json:"user_name,omitempty"`
	Partition            string           `json:"partition,omitempty"`
	ParentAccount        string           `json:"parent_account,omitempty"`
	IsDefault            bool             `json:"is_default,omitempty"`
	Comment              string           `json:"comment,omitempty"`
	DefaultQoS           string           `json:"default_qos,omitempty"`
	QoSList              []string         `json:"qos_list,omitempty"`
	SharesRaw            int32            `json:"shares_raw,omitempty"`
	Priority             int32            `json:"priority,omitempty"`
	MaxJobs              int32            `json:"max_jobs,omitempty"`
	MaxJobsAccrue        int32            `json:"max_jobs_accrue,omitempty"`
	MaxSubmitJobs        int32            `json:"max_submit_jobs,omitempty"`
	MaxWallTime          int32            `json:"max_wall_time,omitempty"`
	MaxCPUTime           int32            `json:"max_cpu_time,omitempty"`
	MaxNodes             int32            `json:"max_nodes,omitempty"`
	MaxCPUs              int32            `json:"max_cpus,omitempty"`
	MaxMemory            int64            `json:"max_memory,omitempty"`
	MinPriorityThreshold int32            `json:"min_priority_threshold,omitempty"`
	GrpJobs              int32            `json:"grp_jobs,omitempty"`
	GrpJobsAccrue        int32            `json:"grp_jobs_accrue,omitempty"`
	GrpNodes             int32            `json:"grp_nodes,omitempty"`
	GrpCPUs              int32            `json:"grp_cpus,omitempty"`
	GrpMemory            int64            `json:"grp_memory,omitempty"`
	GrpSubmitJobs        int32            `json:"grp_submit_jobs,omitempty"`
	GrpWallTime          int32            `json:"grp_wall_time,omitempty"`
	GrpCPUTime           int32            `json:"grp_cpu_time,omitempty"`
	GrpCPURunMins        int64            `json:"grp_cpu_run_mins,omitempty"`
	GrpTRES              map[string]int64 `json:"grp_tres,omitempty"`
	GrpTRESMins          map[string]int64 `json:"grp_tres_mins,omitempty"`
	GrpTRESRunMins       map[string]int64 `json:"grp_tres_run_mins,omitempty"`
	MaxTRES              map[string]int64 `json:"max_tres,omitempty"`
	MaxTRESPerNode       map[string]int64 `json:"max_tres_per_node,omitempty"`
	MaxTRESMins          map[string]int64 `json:"max_tres_mins,omitempty"`
	MinTRES              map[string]int64 `json:"min_tres,omitempty"`
}

// AssociationUpdate represents the data needed to update an association
type AssociationUpdate struct {
	// Identifier field - the association ID is required for updates
	ID *int32 `json:"id,omitempty"` // Association ID (required for updates)

	// Optional identifier fields for filtering/lookup (not used in adapter Update)
	Account   *string `json:"account,omitempty"`   // Account name
	User      *string `json:"user,omitempty"`      // User name
	Cluster   *string `json:"cluster,omitempty"`   // Cluster name
	Partition *string `json:"partition,omitempty"` // Partition name (optional)

	// Update fields
	IsDefault            *bool            `json:"is_default,omitempty"`
	Comment              *string          `json:"comment,omitempty"`
	DefaultQoS           *string          `json:"default_qos,omitempty"`
	QoSList              []string         `json:"qos_list,omitempty"`
	SharesRaw            *int32           `json:"shares_raw,omitempty"`
	Priority             *int32           `json:"priority,omitempty"`
	MaxJobs              *int32           `json:"max_jobs,omitempty"`
	MaxJobsAccrue        *int32           `json:"max_jobs_accrue,omitempty"`
	MaxSubmitJobs        *int32           `json:"max_submit_jobs,omitempty"`
	MaxWallTime          *int32           `json:"max_wall_time,omitempty"`
	MaxCPUTime           *int32           `json:"max_cpu_time,omitempty"`
	MaxNodes             *int32           `json:"max_nodes,omitempty"`
	MaxCPUs              *int32           `json:"max_cpus,omitempty"`
	MaxMemory            *int64           `json:"max_memory,omitempty"`
	MinPriorityThreshold *int32           `json:"min_priority_threshold,omitempty"`
	GrpJobs              *int32           `json:"grp_jobs,omitempty"`
	GrpJobsAccrue        *int32           `json:"grp_jobs_accrue,omitempty"`
	GrpNodes             *int32           `json:"grp_nodes,omitempty"`
	GrpCPUs              *int32           `json:"grp_cpus,omitempty"`
	GrpMemory            *int64           `json:"grp_memory,omitempty"`
	GrpSubmitJobs        *int32           `json:"grp_submit_jobs,omitempty"`
	GrpWallTime          *int32           `json:"grp_wall_time,omitempty"`
	GrpCPUTime           *int32           `json:"grp_cpu_time,omitempty"`
	GrpCPURunMins        *int64           `json:"grp_cpu_run_mins,omitempty"`
	GrpTRES              map[string]int64 `json:"grp_tres,omitempty"`
	GrpTRESMins          map[string]int64 `json:"grp_tres_mins,omitempty"`
	GrpTRESRunMins       map[string]int64 `json:"grp_tres_run_mins,omitempty"`
	MaxTRES              map[string]int64 `json:"max_tres,omitempty"`
	MaxTRESPerNode       map[string]int64 `json:"max_tres_per_node,omitempty"`
	MaxTRESMins          map[string]int64 `json:"max_tres_mins,omitempty"`
	MinTRES              map[string]int64 `json:"min_tres,omitempty"`
}

// AssociationListOptions represents options for listing associations
type AssociationListOptions struct {
	Accounts     []string   `json:"accounts,omitempty"`
	Clusters     []string   `json:"clusters,omitempty"`
	Users        []string   `json:"users,omitempty"`
	Partitions   []string   `json:"partitions,omitempty"`
	OnlyDefaults bool       `json:"only_defaults,omitempty"`
	WithDeleted  bool       `json:"with_deleted,omitempty"`
	UpdateTime   *time.Time `json:"update_time,omitempty"`

	// Limit specifies the maximum number of associations to return.
	// WARNING: Due to SLURM REST API limitations, this is CLIENT-SIDE pagination.
	// The full association list is fetched from the server, then sliced. For large databases,
	// consider using filtering options (Accounts, Clusters, Users, Partitions) to reduce the dataset.
	Limit int `json:"limit,omitempty"`

	// Offset specifies the number of associations to skip before returning results.
	// WARNING: This is CLIENT-SIDE pagination - see Limit field documentation.
	Offset int `json:"offset,omitempty"`
}

// AssociationList represents a list of associations
type AssociationList struct {
	Associations []Association `json:"associations"`
	Total        int           `json:"total"`
}

// AssociationUsage represents usage statistics for an association
type AssociationUsage struct {
	AssociationId    string             `json:"association_id"`
	Account              string `json:"account_name"`
	User                 string `json:"user_name,omitempty"`
	Cluster          string             `json:"cluster"`
	StartTime        time.Time          `json:"start_time"`
	EndTime          time.Time          `json:"end_time"`
	AllocatedCPUTime int64              `json:"allocated_cpu_time"`
	CPUTime          int64              `json:"cpu_time"`
	WallTime         int64              `json:"wall_time"`
	Energy           int64              `json:"energy,omitempty"`
	JobCount         int32              `json:"job_count"`
	TRESUsage        map[string]float64 `json:"tres_usage,omitempty"`
}

// AssociationHierarchy represents the hierarchical structure of associations
type AssociationHierarchy struct {
	Account              string `json:"account_name"`
	Cluster       string                 `json:"cluster"`
	Association   *Association           `json:"association,omitempty"`
	Users         []AssociationUser      `json:"users,omitempty"`
	ChildAccounts []AssociationHierarchy `json:"child_accounts,omitempty"`
}

// AssociationUser represents a user within an association hierarchy
type AssociationUser struct {
	User                 string `json:"user_name"`
	Association *Association `json:"association,omitempty"`
}

// AssociationCreateResponse represents the response from creating an association
type AssociationCreateResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// AccountAssociationRequest represents a request to create account associations
type AccountAssociationRequest struct {
	Accounts     []string          `json:"accounts"`
	Cluster      string            `json:"cluster"`
	Partition    string            `json:"partition,omitempty"`
	Parent       string            `json:"parent,omitempty"`
	QoS          []string          `json:"qos,omitempty"`
	DefaultQoS   string            `json:"default_qos,omitempty"`
	Fairshare    int32             `json:"fairshare,omitempty"`
	GrpTRES      map[string]string `json:"grp_tres,omitempty"`
	MaxTRES      map[string]string `json:"max_tres,omitempty"`
	Description  string            `json:"description,omitempty"`
	Organization string            `json:"organization,omitempty"`
}

// UserAssociationRequest represents a request to create user associations
type UserAssociationRequest struct {
	Users        []string          `json:"users"`
	Account      string            `json:"account"`
	Cluster      string            `json:"cluster"`
	Partition    string            `json:"partition,omitempty"`
	QoS          []string          `json:"qos,omitempty"`
	DefaultQoS   string            `json:"default_qos,omitempty"`
	Fairshare    int32             `json:"fairshare,omitempty"`
	MaxTRES      map[string]string `json:"max_tres,omitempty"`
	DefaultWCKey string            `json:"default_wckey,omitempty"`
	AdminLevel   string            `json:"admin_level,omitempty"`
}
