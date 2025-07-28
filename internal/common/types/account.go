package types

import (
	"time"
)

// Account represents a SLURM account with common fields across all API versions
type Account struct {
	Name                string            `json:"name"`
	Description         string            `json:"description,omitempty"`
	Organization        string            `json:"organization,omitempty"`
	Coordinators        []string          `json:"coordinators,omitempty"`
	DefaultQoS          string            `json:"default_qos,omitempty"`
	QoSList             []string          `json:"qos_list,omitempty"`
	ParentName          string            `json:"parent_name,omitempty"`
	ChildAccounts       []string          `json:"child_accounts,omitempty"`
	AllowedPartitions   []string          `json:"allowed_partitions,omitempty"`
	DefaultPartition    string            `json:"default_partition,omitempty"`
	FairShare           int32             `json:"fair_share,omitempty"`
	SharesRaw           int32             `json:"shares_raw,omitempty"`
	Priority            int32             `json:"priority,omitempty"`
	MaxJobs             int32             `json:"max_jobs,omitempty"`
	MaxJobsPerUser      int32             `json:"max_jobs_per_user,omitempty"`
	MaxSubmitJobs       int32             `json:"max_submit_jobs,omitempty"`
	MaxWallTime         int32             `json:"max_wall_time,omitempty"`
	MaxCPUTime          int32             `json:"max_cpu_time,omitempty"`
	MaxNodes            int32             `json:"max_nodes,omitempty"`
	MaxCPUs             int32             `json:"max_cpus,omitempty"`
	MaxMemory           int64             `json:"max_memory,omitempty"`
	MinPriorityThreshold int32            `json:"min_priority_threshold,omitempty"`
	GrpJobs             int32             `json:"grp_jobs,omitempty"`
	GrpJobsAccrue       int32             `json:"grp_jobs_accrue,omitempty"`
	GrpNodes            int32             `json:"grp_nodes,omitempty"`
	GrpCPUs             int32             `json:"grp_cpus,omitempty"`
	GrpMemory           int64             `json:"grp_memory,omitempty"`
	GrpSubmitJobs       int32             `json:"grp_submit_jobs,omitempty"`
	GrpWallTime         int32             `json:"grp_wall_time,omitempty"`
	GrpCPUTime          int32             `json:"grp_cpu_time,omitempty"`
	GrpTRES             map[string]int64  `json:"grp_tres,omitempty"`
	GrpTRESMins         map[string]int64  `json:"grp_tres_mins,omitempty"`
	GrpTRESRunMins      map[string]int64  `json:"grp_tres_run_mins,omitempty"`
	MaxTRES             map[string]int64  `json:"max_tres,omitempty"`
	MaxTRESPerNode      map[string]int64  `json:"max_tres_per_node,omitempty"`
	MinTRES             map[string]int64  `json:"min_tres,omitempty"`
	IsDefault           bool              `json:"is_default,omitempty"`
	Deleted             bool              `json:"deleted,omitempty"`
}

// AccountCreate represents the data needed to create a new account
type AccountCreate struct {
	Name                string            `json:"name"`
	Description         string            `json:"description,omitempty"`
	Organization        string            `json:"organization,omitempty"`
	Coordinators        []string          `json:"coordinators,omitempty"`
	DefaultQoS          string            `json:"default_qos,omitempty"`
	QoSList             []string          `json:"qos_list,omitempty"`
	ParentName          string            `json:"parent_name,omitempty"`
	AllowedPartitions   []string          `json:"allowed_partitions,omitempty"`
	DefaultPartition    string            `json:"default_partition,omitempty"`
	FairShare           int32             `json:"fair_share,omitempty"`
	SharesRaw           int32             `json:"shares_raw,omitempty"`
	Priority            int32             `json:"priority,omitempty"`
	MaxJobs             int32             `json:"max_jobs,omitempty"`
	MaxJobsPerUser      int32             `json:"max_jobs_per_user,omitempty"`
	MaxSubmitJobs       int32             `json:"max_submit_jobs,omitempty"`
	MaxWallTime         int32             `json:"max_wall_time,omitempty"`
	MaxCPUTime          int32             `json:"max_cpu_time,omitempty"`
	MaxNodes            int32             `json:"max_nodes,omitempty"`
	MaxCPUs             int32             `json:"max_cpus,omitempty"`
	MaxMemory           int64             `json:"max_memory,omitempty"`
	MinPriorityThreshold int32            `json:"min_priority_threshold,omitempty"`
	GrpJobs             int32             `json:"grp_jobs,omitempty"`
	GrpJobsAccrue       int32             `json:"grp_jobs_accrue,omitempty"`
	GrpNodes            int32             `json:"grp_nodes,omitempty"`
	GrpCPUs             int32             `json:"grp_cpus,omitempty"`
	GrpMemory           int64             `json:"grp_memory,omitempty"`
	GrpSubmitJobs       int32             `json:"grp_submit_jobs,omitempty"`
	GrpWallTime         int32             `json:"grp_wall_time,omitempty"`
	GrpCPUTime          int32             `json:"grp_cpu_time,omitempty"`
	GrpTRES             map[string]int64  `json:"grp_tres,omitempty"`
	GrpTRESMins         map[string]int64  `json:"grp_tres_mins,omitempty"`
	GrpTRESRunMins      map[string]int64  `json:"grp_tres_run_mins,omitempty"`
	MaxTRES             map[string]int64  `json:"max_tres,omitempty"`
	MaxTRESPerNode      map[string]int64  `json:"max_tres_per_node,omitempty"`
	MinTRES             map[string]int64  `json:"min_tres,omitempty"`
}

// AccountUpdate represents the data needed to update an account
type AccountUpdate struct {
	Description         *string           `json:"description,omitempty"`
	Organization        *string           `json:"organization,omitempty"`
	Coordinators        []string          `json:"coordinators,omitempty"`
	DefaultQoS          *string           `json:"default_qos,omitempty"`
	QoSList             []string          `json:"qos_list,omitempty"`
	AllowedPartitions   []string          `json:"allowed_partitions,omitempty"`
	DefaultPartition    *string           `json:"default_partition,omitempty"`
	FairShare           *int32            `json:"fair_share,omitempty"`
	SharesRaw           *int32            `json:"shares_raw,omitempty"`
	Priority            *int32            `json:"priority,omitempty"`
	MaxJobs             *int32            `json:"max_jobs,omitempty"`
	MaxJobsPerUser      *int32            `json:"max_jobs_per_user,omitempty"`
	MaxSubmitJobs       *int32            `json:"max_submit_jobs,omitempty"`
	MaxWallTime         *int32            `json:"max_wall_time,omitempty"`
	MaxCPUTime          *int32            `json:"max_cpu_time,omitempty"`
	MaxNodes            *int32            `json:"max_nodes,omitempty"`
	MaxCPUs             *int32            `json:"max_cpus,omitempty"`
	MaxMemory           *int64            `json:"max_memory,omitempty"`
	MinPriorityThreshold *int32           `json:"min_priority_threshold,omitempty"`
	GrpJobs             *int32            `json:"grp_jobs,omitempty"`
	GrpJobsAccrue       *int32            `json:"grp_jobs_accrue,omitempty"`
	GrpNodes            *int32            `json:"grp_nodes,omitempty"`
	GrpCPUs             *int32            `json:"grp_cpus,omitempty"`
	GrpMemory           *int64            `json:"grp_memory,omitempty"`
	GrpSubmitJobs       *int32            `json:"grp_submit_jobs,omitempty"`
	GrpWallTime         *int32            `json:"grp_wall_time,omitempty"`
	GrpCPUTime          *int32            `json:"grp_cpu_time,omitempty"`
	GrpTRES             map[string]int64  `json:"grp_tres,omitempty"`
	GrpTRESMins         map[string]int64  `json:"grp_tres_mins,omitempty"`
	GrpTRESRunMins      map[string]int64  `json:"grp_tres_run_mins,omitempty"`
	MaxTRES             map[string]int64  `json:"max_tres,omitempty"`
	MaxTRESPerNode      map[string]int64  `json:"max_tres_per_node,omitempty"`
	MinTRES             map[string]int64  `json:"min_tres,omitempty"`
}

// AccountCreateResponse represents the response from account creation
type AccountCreateResponse struct {
	AccountName string `json:"account_name"`
}

// AccountListOptions represents options for listing accounts
type AccountListOptions struct {
	Names         []string   `json:"names,omitempty"`
	Descriptions  []string   `json:"descriptions,omitempty"`
	Organizations []string   `json:"organizations,omitempty"`
	WithDeleted   bool       `json:"with_deleted,omitempty"`
	WithAssocs    bool       `json:"with_assocs,omitempty"`
	WithCoords    bool       `json:"with_coords,omitempty"`
	WithWCKeys    bool       `json:"with_wckeys,omitempty"`
	UpdateTime    *time.Time `json:"update_time,omitempty"`
	Limit         int        `json:"limit,omitempty"`
	Offset        int        `json:"offset,omitempty"`
}

// AccountList represents a list of accounts
type AccountList struct {
	Accounts []Account `json:"accounts"`
	Total    int       `json:"total"`
}

// AccountUsage represents usage statistics for an account
type AccountUsage struct {
	AccountName      string             `json:"account_name"`
	StartTime        time.Time          `json:"start_time"`
	EndTime          time.Time          `json:"end_time"`
	CPUSeconds       int64              `json:"cpu_seconds"`
	NodeHours        float64            `json:"node_hours"`
	JobCount         int32              `json:"job_count"`
	UserCount        int32              `json:"user_count"`
	TRESUsage        map[string]float64 `json:"tres_usage,omitempty"`
	AverageJobSize   float64            `json:"average_job_size,omitempty"`
	AverageWaitTime  int32              `json:"average_wait_time,omitempty"`
	SuccessRate      float64            `json:"success_rate,omitempty"`
}

// AccountLimits represents resource limits for an account
type AccountLimits struct {
	AccountName      string           `json:"account_name"`
	MaxJobs          int32            `json:"max_jobs"`
	MaxCPUs          int32            `json:"max_cpus"`
	MaxNodes         int32            `json:"max_nodes"`
	MaxMemory        int64            `json:"max_memory"`
	MaxWallTime      int32            `json:"max_wall_time"`
	CurrentJobs      int32            `json:"current_jobs"`
	CurrentCPUs      int32            `json:"current_cpus"`
	CurrentNodes     int32            `json:"current_nodes"`
	CurrentMemory    int64            `json:"current_memory"`
	AvailableJobs    int32            `json:"available_jobs"`
	AvailableCPUs    int32            `json:"available_cpus"`
	AvailableNodes   int32            `json:"available_nodes"`
	AvailableMemory  int64            `json:"available_memory"`
	TRES             map[string]int64 `json:"tres,omitempty"`
}