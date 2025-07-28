package types

import (
	"time"
)

// Association represents a SLURM association with common fields across all API versions
type Association struct {
	ID                  string            `json:"id"`
	AccountName         string            `json:"account_name"`
	Cluster             string            `json:"cluster"`
	UserName            string            `json:"user_name,omitempty"`
	Partition           string            `json:"partition,omitempty"`
	ParentAccount       string            `json:"parent_account,omitempty"`
	IsDefault           bool              `json:"is_default,omitempty"`
	Comment             string            `json:"comment,omitempty"`
	DefaultQoS          string            `json:"default_qos,omitempty"`
	QoSList             []string          `json:"qos_list,omitempty"`
	SharesRaw           int32             `json:"shares_raw,omitempty"`
	SharesNormalized    float64           `json:"shares_normalized,omitempty"`
	SharesLevel         int32             `json:"shares_level,omitempty"`
	FairShareLevel      float64           `json:"fair_share_level,omitempty"`
	EffectiveUsage      float64           `json:"effective_usage,omitempty"`
	NormalizedUsage     float64           `json:"normalized_usage,omitempty"`
	RawUsage            int64             `json:"raw_usage,omitempty"`
	Priority            int32             `json:"priority,omitempty"`
	MaxJobs             int32             `json:"max_jobs,omitempty"`
	MaxJobsAccrue       int32             `json:"max_jobs_accrue,omitempty"`
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
	GrpCPURunMins       int64             `json:"grp_cpu_run_mins,omitempty"`
	GrpTRES             map[string]int64  `json:"grp_tres,omitempty"`
	GrpTRESMins         map[string]int64  `json:"grp_tres_mins,omitempty"`
	GrpTRESRunMins      map[string]int64  `json:"grp_tres_run_mins,omitempty"`
	MaxTRES             map[string]int64  `json:"max_tres,omitempty"`
	MaxTRESPerNode      map[string]int64  `json:"max_tres_per_node,omitempty"`
	MaxTRESMins         map[string]int64  `json:"max_tres_mins,omitempty"`
	MinTRES             map[string]int64  `json:"min_tres,omitempty"`
	UsedJobs            int32             `json:"used_jobs,omitempty"`
	UsedJobsAccrue      int32             `json:"used_jobs_accrue,omitempty"`
	UsedSubmitJobs      int32             `json:"used_submit_jobs,omitempty"`
	UsedCPUs            int32             `json:"used_cpus,omitempty"`
	UsedMemory          int64             `json:"used_memory,omitempty"`
	UsedNodes           int32             `json:"used_nodes,omitempty"`
	UsedWallTime        int64             `json:"used_wall_time,omitempty"`
	UsedCPUTime         int64             `json:"used_cpu_time,omitempty"`
	Deleted             bool              `json:"deleted,omitempty"`
	CreatedTime         *time.Time        `json:"created_time,omitempty"`
	ModifiedTime        *time.Time        `json:"modified_time,omitempty"`
}

// AssociationCreate represents the data needed to create a new association
type AssociationCreate struct {
	AccountName         string            `json:"account_name"`
	Cluster             string            `json:"cluster"`
	UserName            string            `json:"user_name,omitempty"`
	Partition           string            `json:"partition,omitempty"`
	ParentAccount       string            `json:"parent_account,omitempty"`
	IsDefault           bool              `json:"is_default,omitempty"`
	Comment             string            `json:"comment,omitempty"`
	DefaultQoS          string            `json:"default_qos,omitempty"`
	QoSList             []string          `json:"qos_list,omitempty"`
	SharesRaw           int32             `json:"shares_raw,omitempty"`
	Priority            int32             `json:"priority,omitempty"`
	MaxJobs             int32             `json:"max_jobs,omitempty"`
	MaxJobsAccrue       int32             `json:"max_jobs_accrue,omitempty"`
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
	GrpCPURunMins       int64             `json:"grp_cpu_run_mins,omitempty"`
	GrpTRES             map[string]int64  `json:"grp_tres,omitempty"`
	GrpTRESMins         map[string]int64  `json:"grp_tres_mins,omitempty"`
	GrpTRESRunMins      map[string]int64  `json:"grp_tres_run_mins,omitempty"`
	MaxTRES             map[string]int64  `json:"max_tres,omitempty"`
	MaxTRESPerNode      map[string]int64  `json:"max_tres_per_node,omitempty"`
	MaxTRESMins         map[string]int64  `json:"max_tres_mins,omitempty"`
	MinTRES             map[string]int64  `json:"min_tres,omitempty"`
}

// AssociationUpdate represents the data needed to update an association
type AssociationUpdate struct {
	IsDefault           *bool             `json:"is_default,omitempty"`
	Comment             *string           `json:"comment,omitempty"`
	DefaultQoS          *string           `json:"default_qos,omitempty"`
	QoSList             []string          `json:"qos_list,omitempty"`
	SharesRaw           *int32            `json:"shares_raw,omitempty"`
	Priority            *int32            `json:"priority,omitempty"`
	MaxJobs             *int32            `json:"max_jobs,omitempty"`
	MaxJobsAccrue       *int32            `json:"max_jobs_accrue,omitempty"`
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
	GrpCPURunMins       *int64            `json:"grp_cpu_run_mins,omitempty"`
	GrpTRES             map[string]int64  `json:"grp_tres,omitempty"`
	GrpTRESMins         map[string]int64  `json:"grp_tres_mins,omitempty"`
	GrpTRESRunMins      map[string]int64  `json:"grp_tres_run_mins,omitempty"`
	MaxTRES             map[string]int64  `json:"max_tres,omitempty"`
	MaxTRESPerNode      map[string]int64  `json:"max_tres_per_node,omitempty"`
	MaxTRESMins         map[string]int64  `json:"max_tres_mins,omitempty"`
	MinTRES             map[string]int64  `json:"min_tres,omitempty"`
}

// AssociationCreateResponse represents the response from association creation
type AssociationCreateResponse struct {
	AssociationID string `json:"association_id"`
	AccountName   string `json:"account_name"`
	UserName      string `json:"user_name,omitempty"`
	Cluster       string `json:"cluster"`
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
	Limit        int        `json:"limit,omitempty"`
	Offset       int        `json:"offset,omitempty"`
}

// AssociationList represents a list of associations
type AssociationList struct {
	Associations []Association `json:"associations"`
	Total        int           `json:"total"`
}

// AssociationUsage represents usage statistics for an association
type AssociationUsage struct {
	AssociationID    string             `json:"association_id"`
	AccountName      string             `json:"account_name"`
	UserName         string             `json:"user_name,omitempty"`
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
	AccountName  string                  `json:"account_name"`
	Cluster      string                  `json:"cluster"`
	Association  *Association            `json:"association,omitempty"`
	Users        []AssociationUser       `json:"users,omitempty"`
	ChildAccounts []AssociationHierarchy `json:"child_accounts,omitempty"`
}

// AssociationUser represents a user within an association hierarchy
type AssociationUser struct {
	UserName    string       `json:"user_name"`
	Association *Association `json:"association,omitempty"`
}