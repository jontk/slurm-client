package types

import (
	"time"
)

// User represents a SLURM user with common fields across all API versions
type User struct {
	Name                string            `json:"name"`
	UID                 int32             `json:"uid,omitempty"`
	DefaultAccount      string            `json:"default_account,omitempty"`
	DefaultWCKey        string            `json:"default_wckey,omitempty"`
	AdminLevel          AdminLevel        `json:"admin_level,omitempty"`
	Accounts            []string          `json:"accounts,omitempty"`
	Coordinators        []UserCoordinator `json:"coordinators,omitempty"`
	DefaultQoS          string            `json:"default_qos,omitempty"`
	QoSList             []string          `json:"qos_list,omitempty"`
	MaxJobs             int32             `json:"max_jobs,omitempty"`
	MaxJobsPerAccount   int32             `json:"max_jobs_per_account,omitempty"`
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
	Associations        []UserAssociation `json:"associations,omitempty"`
	WCKeys              []string          `json:"wckeys,omitempty"`
	Deleted             bool              `json:"deleted,omitempty"`
}

// AdminLevel represents the administrative level of a user
type AdminLevel string

const (
	AdminLevelNone         AdminLevel = "None"
	AdminLevelOperator     AdminLevel = "Operator"
	AdminLevelAdministrator AdminLevel = "Administrator"
)

// UserCoordinator represents coordinator information for a user
type UserCoordinator struct {
	AccountName string `json:"account_name"`
	Coordinator string `json:"coordinator"`
}

// UserAssociation represents an association between a user and an account
type UserAssociation struct {
	AccountName      string           `json:"account_name"`
	Cluster          string           `json:"cluster"`
	Partition        string           `json:"partition,omitempty"`
	UserName         string           `json:"user_name"`
	DefaultQoS       string           `json:"default_qos,omitempty"`
	QoSList          []string         `json:"qos_list,omitempty"`
	SharesRaw        int32            `json:"shares_raw,omitempty"`
	Priority         int32            `json:"priority,omitempty"`
	MaxJobs          int32            `json:"max_jobs,omitempty"`
	MaxJobsAccrue    int32            `json:"max_jobs_accrue,omitempty"`
	MaxSubmitJobs    int32            `json:"max_submit_jobs,omitempty"`
	MaxWallTime      int32            `json:"max_wall_time,omitempty"`
	MaxCPUTime       int32            `json:"max_cpu_time,omitempty"`
	MaxNodes         int32            `json:"max_nodes,omitempty"`
	MaxCPUs          int32            `json:"max_cpus,omitempty"`
	MaxMemory        int64            `json:"max_memory,omitempty"`
	GrpJobs          int32            `json:"grp_jobs,omitempty"`
	GrpJobsAccrue    int32            `json:"grp_jobs_accrue,omitempty"`
	GrpNodes         int32            `json:"grp_nodes,omitempty"`
	GrpCPUs          int32            `json:"grp_cpus,omitempty"`
	GrpMemory        int64            `json:"grp_memory,omitempty"`
	GrpSubmitJobs    int32            `json:"grp_submit_jobs,omitempty"`
	GrpWallTime      int32            `json:"grp_wall_time,omitempty"`
	GrpCPUTime       int32            `json:"grp_cpu_time,omitempty"`
	GrpTRES          map[string]int64 `json:"grp_tres,omitempty"`
	GrpTRESMins      map[string]int64 `json:"grp_tres_mins,omitempty"`
	GrpTRESRunMins   map[string]int64 `json:"grp_tres_run_mins,omitempty"`
	MaxTRES          map[string]int64 `json:"max_tres,omitempty"`
	MaxTRESPerNode   map[string]int64 `json:"max_tres_per_node,omitempty"`
	MinTRES          map[string]int64 `json:"min_tres,omitempty"`
	Comment          string           `json:"comment,omitempty"`
	Deleted          bool             `json:"deleted,omitempty"`
}

// UserCreate represents the data needed to create a new user
type UserCreate struct {
	Name                string            `json:"name"`
	UID                 int32             `json:"uid,omitempty"`
	DefaultAccount      string            `json:"default_account,omitempty"`
	DefaultWCKey        string            `json:"default_wckey,omitempty"`
	AdminLevel          AdminLevel        `json:"admin_level,omitempty"`
	Accounts            []string          `json:"accounts,omitempty"`
	DefaultQoS          string            `json:"default_qos,omitempty"`
	QoSList             []string          `json:"qos_list,omitempty"`
	MaxJobs             int32             `json:"max_jobs,omitempty"`
	MaxJobsPerAccount   int32             `json:"max_jobs_per_account,omitempty"`
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
	WCKeys              []string          `json:"wckeys,omitempty"`
}

// UserUpdate represents the data needed to update a user
type UserUpdate struct {
	DefaultAccount      *string           `json:"default_account,omitempty"`
	DefaultWCKey        *string           `json:"default_wckey,omitempty"`
	AdminLevel          *AdminLevel       `json:"admin_level,omitempty"`
	Accounts            []string          `json:"accounts,omitempty"`
	DefaultQoS          *string           `json:"default_qos,omitempty"`
	QoSList             []string          `json:"qos_list,omitempty"`
	MaxJobs             *int32            `json:"max_jobs,omitempty"`
	MaxJobsPerAccount   *int32            `json:"max_jobs_per_account,omitempty"`
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
	WCKeys              []string          `json:"wckeys,omitempty"`
}

// UserCreateResponse represents the response from user creation
type UserCreateResponse struct {
	UserName string `json:"user_name"`
}

// UserListOptions represents options for listing users
type UserListOptions struct {
	Names          []string   `json:"names,omitempty"`
	DefaultAccount string     `json:"default_account,omitempty"`
	WithDeleted    bool       `json:"with_deleted,omitempty"`
	WithAssocs     bool       `json:"with_assocs,omitempty"`
	WithCoords     bool       `json:"with_coords,omitempty"`
	WithWCKeys     bool       `json:"with_wckeys,omitempty"`
	UpdateTime     *time.Time `json:"update_time,omitempty"`
	Limit          int        `json:"limit,omitempty"`
	Offset         int        `json:"offset,omitempty"`
}

// UserList represents a list of users
type UserList struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
}

// UserUsage represents usage statistics for a user
type UserUsage struct {
	UserName         string             `json:"user_name"`
	AccountName      string             `json:"account_name,omitempty"`
	StartTime        time.Time          `json:"start_time"`
	EndTime          time.Time          `json:"end_time"`
	CPUSeconds       int64              `json:"cpu_seconds"`
	NodeHours        float64            `json:"node_hours"`
	JobCount         int32              `json:"job_count"`
	TRESUsage        map[string]float64 `json:"tres_usage,omitempty"`
	AverageJobSize   float64            `json:"average_job_size,omitempty"`
	AverageWaitTime  int32              `json:"average_wait_time,omitempty"`
	SuccessRate      float64            `json:"success_rate,omitempty"`
}

// UserPermissions represents permissions for a user
type UserPermissions struct {
	UserName            string   `json:"user_name"`
	CanSubmitJobs       bool     `json:"can_submit_jobs"`
	CanViewJobs         bool     `json:"can_view_jobs"`
	CanCancelJobs       bool     `json:"can_cancel_jobs"`
	CanModifyJobs       bool     `json:"can_modify_jobs"`
	CanViewAllJobs      bool     `json:"can_view_all_jobs"`
	CanManageReservations bool   `json:"can_manage_reservations"`
	CanManageAccounts   bool     `json:"can_manage_accounts"`
	CanManageUsers      bool     `json:"can_manage_users"`
	CanManageQoS        bool     `json:"can_manage_qos"`
	AllowedPartitions   []string `json:"allowed_partitions,omitempty"`
	AllowedQoS          []string `json:"allowed_qos,omitempty"`
	AllowedAccounts     []string `json:"allowed_accounts,omitempty"`
}