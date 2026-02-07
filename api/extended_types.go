// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Extended types for the SDK - these are SDK-specific types not in the SLURM REST API.
// Analytics types are in pkg/analytics.
package api

import (
	"net/http"
	"time"
)

// ============================================================================
// Job Submission Types
// ============================================================================

// JobSubmission represents a job submission request.
type JobSubmission struct {
	Name        string            `json:"name"`
	Account     string            `json:"account,omitempty"`
	Script      string            `json:"script,omitempty"`
	Command     string            `json:"command,omitempty"`
	Args        []string          `json:"args,omitempty"`
	Partition   string            `json:"partition,omitempty"`
	CPUs        int               `json:"cpus,omitempty"`
	Memory      int               `json:"memory,omitempty"`
	TimeLimit   int               `json:"time_limit,omitempty"`
	WorkingDir  string            `json:"working_dir,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Nodes       int               `json:"nodes,omitempty"`
	Priority    int               `json:"priority,omitempty"`
}

// JobStepList represents a list of job steps.
type JobStepList struct {
	Steps []JobStep `json:"steps"`
	Total int       `json:"total"`
}

// JobStep represents a single job step.
type JobStep struct {
	ID        string     `json:"id"`
	JobID     string     `json:"job_id"`
	Name      string     `json:"name"`
	State     string     `json:"state"`
	CPUs      int        `json:"cpus"`
	Memory    int        `json:"memory"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	ExitCode  int        `json:"exit_code"`
}

// ============================================================================
// List Options Types
// ============================================================================

// ListJobsOptions configures job listing with a simple, string-based interface.
//
// This type provides a simpler API for common use cases with string-based filtering.
// For more comprehensive filtering with type-safe JobState values and multi-value
// filtering, see JobListOptions in api/job.go.
//
// Note: Both types are maintained for different use cases:
//   - ListJobsOptions: Simple string-based filtering (UserID, Partition, []string States)
//   - JobListOptions: Type-safe filtering ([]string Users, []string Partitions, []JobState States)
type ListJobsOptions struct {
	UserID    string   `json:"user_id,omitempty"`
	States    []string `json:"states,omitempty"`
	Partition string   `json:"partition,omitempty"`
	Limit     int      `json:"limit,omitempty"`
	Offset    int      `json:"offset,omitempty"`
}

// ListNodesOptions configures node listing.
type ListNodesOptions struct {
	States    []string `json:"states,omitempty"`
	Partition string   `json:"partition,omitempty"`
	Features  []string `json:"features,omitempty"`
	Limit     int      `json:"limit,omitempty"`
	Offset    int      `json:"offset,omitempty"`
}

// ListPartitionsOptions configures partition listing.
type ListPartitionsOptions struct {
	States []string `json:"states,omitempty"`
	Limit  int      `json:"limit,omitempty"`
	Offset int      `json:"offset,omitempty"`
}

// ListReservationsOptions configures reservation listing.
type ListReservationsOptions struct {
	Names    []string `json:"names,omitempty"`
	Users    []string `json:"users,omitempty"`
	Accounts []string `json:"accounts,omitempty"`
	States   []string `json:"states,omitempty"`
	Limit    int      `json:"limit,omitempty"`
	Offset   int      `json:"offset,omitempty"`
}

// ListQoSOptions configures QoS listing.
type ListQoSOptions struct {
	Names    []string `json:"names,omitempty"`
	Accounts []string `json:"accounts,omitempty"`
	Users    []string `json:"users,omitempty"`
	Limit    int      `json:"limit,omitempty"`
	Offset   int      `json:"offset,omitempty"`
}

// ListAccountsOptions configures account listing.
type ListAccountsOptions struct {
	Names            []string `json:"names,omitempty"`
	Organizations    []string `json:"organizations,omitempty"`
	ParentAccounts   []string `json:"parent_accounts,omitempty"`
	WithAssociations bool     `json:"with_associations,omitempty"`
	WithCoordinators bool     `json:"with_coordinators,omitempty"`
	WithDeleted      bool     `json:"with_deleted,omitempty"`
	WithUsers        bool     `json:"with_users,omitempty"`
	WithQuotas       bool     `json:"with_quotas,omitempty"`
	WithUsage        bool     `json:"with_usage,omitempty"`
	Limit            int      `json:"limit,omitempty"`
	Offset           int      `json:"offset,omitempty"`
}

// ListClustersOptions configures cluster listing.
type ListClustersOptions struct {
	Names            []string `json:"names,omitempty"`
	FederationStates []string `json:"federation_states,omitempty"`
	Features         []string `json:"features,omitempty"`
	ControlHosts     []string `json:"control_hosts,omitempty"`
	WithFederation   bool     `json:"with_federation,omitempty"`
	WithTRES         bool     `json:"with_tres,omitempty"`
	WithPlugins      bool     `json:"with_plugins,omitempty"`
	Offset           int      `json:"offset,omitempty"`
	Limit            int      `json:"limit,omitempty"`
}

// ListAssociationsOptions configures association listing.
type ListAssociationsOptions struct {
	Users           []string `json:"users,omitempty"`
	Accounts        []string `json:"accounts,omitempty"`
	Clusters        []string `json:"clusters,omitempty"`
	Partitions      []string `json:"partitions,omitempty"`
	ParentAccounts  []string `json:"parent_accounts,omitempty"`
	QoS             []string `json:"qos,omitempty"`
	WithDeleted     bool     `json:"with_deleted,omitempty"`
	WithUsage       bool     `json:"with_usage,omitempty"`
	WithTRES        bool     `json:"with_tres,omitempty"`
	WithSubAccounts bool     `json:"with_sub_accounts,omitempty"`
	OnlyDefaults    bool     `json:"only_defaults,omitempty"`
	Offset          int      `json:"offset,omitempty"`
	Limit           int      `json:"limit,omitempty"`
}

// ListUsersOptions configures user listing.
type ListUsersOptions struct {
	Names            []string `json:"names,omitempty"`
	Accounts         []string `json:"accounts,omitempty"`
	Clusters         []string `json:"clusters,omitempty"`
	AdminLevels      []string `json:"admin_levels,omitempty"`
	ActiveOnly       bool     `json:"active_only,omitempty"`
	CoordinatorsOnly bool     `json:"coordinators_only,omitempty"`
	WithAccounts     bool     `json:"with_accounts,omitempty"`
	WithQuotas       bool     `json:"with_quotas,omitempty"`
	WithFairShare    bool     `json:"with_fair_share,omitempty"`
	WithAssociations bool     `json:"with_associations,omitempty"`
	WithUsage        bool     `json:"with_usage,omitempty"`
	Limit            int      `json:"limit,omitempty"`
	Offset           int      `json:"offset,omitempty"`
	SortBy           string   `json:"sort_by,omitempty"`
	SortOrder        string   `json:"sort_order,omitempty"`
}

// ============================================================================
// Watch Options Types
// ============================================================================

// WatchJobsOptions configures job watching.
type WatchJobsOptions struct {
	UserID           string   `json:"user_id,omitempty"`
	States           []string `json:"states,omitempty"`
	Partition        string   `json:"partition,omitempty"`
	JobIDs           []string `json:"job_ids,omitempty"`
	ExcludeNew       bool     `json:"exclude_new,omitempty"`
	ExcludeCompleted bool     `json:"exclude_completed,omitempty"`
}

// WatchNodesOptions configures node watching.
type WatchNodesOptions struct {
	States    []string `json:"states,omitempty"`
	Partition string   `json:"partition,omitempty"`
	Features  []string `json:"features,omitempty"`
	NodeNames []string `json:"node_names,omitempty"`
}

// WatchPartitionsOptions configures partition watching.
type WatchPartitionsOptions struct {
	States         []string `json:"states,omitempty"`
	PartitionNames []string `json:"partition_names,omitempty"`
}

// ============================================================================
// Association Types
// ============================================================================

// GetAssociationOptions configures association retrieval.
type GetAssociationOptions struct {
	User      string `json:"user"`
	Account   string `json:"account"`
	Cluster   string `json:"cluster,omitempty"`
	Partition string `json:"partition,omitempty"`
	WithUsage bool   `json:"with_usage,omitempty"`
	WithTRES  bool   `json:"with_tres,omitempty"`
}

// DeleteAssociationOptions configures association deletion.
type DeleteAssociationOptions struct {
	User      string `json:"user"`
	Account   string `json:"account"`
	Cluster   string `json:"cluster,omitempty"`
	Partition string `json:"partition,omitempty"`
	Force     bool   `json:"force,omitempty"`
}

// BulkDeleteOptions configures bulk association deletion.
type BulkDeleteOptions struct {
	Users      []string `json:"users,omitempty"`
	Accounts   []string `json:"accounts,omitempty"`
	Clusters   []string `json:"clusters,omitempty"`
	Partitions []string `json:"partitions,omitempty"`
	OnlyIfIdle bool     `json:"only_if_idle,omitempty"`
	Force      bool     `json:"force,omitempty"`
}

// BulkDeleteResponse represents bulk deletion results.
type BulkDeleteResponse struct {
	Deleted             int            `json:"deleted"`
	Failed              int            `json:"failed"`
	Errors              []string       `json:"errors,omitempty"`
	DeletedAssociations []*Association `json:"deleted_associations,omitempty"`
}

// AssociationOptions configures association creation.
type AssociationOptions struct {
	Cluster   string `json:"cluster,omitempty"`
	Partition string `json:"partition,omitempty"`
	QoS       string `json:"qos,omitempty"`
	MaxJobs   *int   `json:"max_jobs,omitempty"`
	Priority  *int   `json:"priority,omitempty"`
}

// ============================================================================
// Account Types
// ============================================================================

// AccountQuota represents account resource quotas.
type AccountQuota struct {
	CPULimit           int            `json:"cpu_limit,omitempty"`
	CPUUsed            int            `json:"cpu_used,omitempty"`
	MaxJobs            int            `json:"max_jobs,omitempty"`
	JobsUsed           int            `json:"jobs_used,omitempty"`
	MaxJobsPerUser     int            `json:"max_jobs_per_user,omitempty"`
	MaxNodes           int            `json:"max_nodes,omitempty"`
	NodesUsed          int            `json:"nodes_used,omitempty"`
	MaxWallTime        int            `json:"max_wall_time,omitempty"`
	GrpTRES            map[string]int `json:"grp_tres,omitempty"`
	GrpTRESUsed        map[string]int `json:"grp_tres_used,omitempty"`
	GrpTRESMinutes     map[string]int `json:"grp_tres_minutes,omitempty"`
	GrpTRESMinutesUsed map[string]int `json:"grp_tres_minutes_used,omitempty"`
	MaxTRES            map[string]int `json:"max_tres,omitempty"`
	MaxTRESUsed        map[string]int `json:"max_tres_used,omitempty"`
	MaxTRESPerUser     map[string]int `json:"max_tres_per_user,omitempty"`
	QuotaPeriod        string         `json:"quota_period,omitempty"`
	LastUpdated        time.Time      `json:"last_updated,omitempty"`
}

// AccountHierarchy represents account hierarchy.
type AccountHierarchy struct {
	Account          *Account            `json:"account"`
	ParentAccount    *AccountHierarchy   `json:"parent_account,omitempty"`
	ChildAccounts    []*AccountHierarchy `json:"child_accounts,omitempty"`
	Level            int                 `json:"level"`
	Path             []string            `json:"path"`
	TotalUsers       int                 `json:"total_users"`
	TotalSubAccounts int                 `json:"total_sub_accounts"`
	AggregateQuota   *AccountQuota       `json:"aggregate_quota,omitempty"`
	AggregateUsage   *AccountUsage       `json:"aggregate_usage,omitempty"`
}

// ListAccountUsersOptions configures listing users in an account.
type ListAccountUsersOptions struct {
	Roles            []string `json:"roles,omitempty"`
	Permissions      []string `json:"permissions,omitempty"`
	ActiveOnly       bool     `json:"active_only,omitempty"`
	CoordinatorsOnly bool     `json:"coordinators_only,omitempty"`
	Partitions       []string `json:"partitions,omitempty"`
	QoS              []string `json:"qos,omitempty"`
	WithPermissions  bool     `json:"with_permissions,omitempty"`
	WithQuotas       bool     `json:"with_quotas,omitempty"`
	WithUsage        bool     `json:"with_usage,omitempty"`
	WithFairShare    bool     `json:"with_fair_share,omitempty"`
	Limit            int      `json:"limit,omitempty"`
	Offset           int      `json:"offset,omitempty"`
	SortBy           string   `json:"sort_by,omitempty"`
	SortOrder        string   `json:"sort_order,omitempty"`
}

// ============================================================================
// User Types
// ============================================================================

// UserAccount represents a user's account association.
type UserAccount struct {
	AccountName   string         `json:"account_name"`
	Partition     string         `json:"partition,omitempty"`
	QoS           string         `json:"qos,omitempty"`
	DefaultQoS    string         `json:"default_qos,omitempty"`
	MaxJobs       int            `json:"max_jobs,omitempty"`
	MaxSubmitJobs int            `json:"max_submit_jobs,omitempty"`
	MaxWallTime   int            `json:"max_wall_time,omitempty"`
	Priority      int            `json:"priority,omitempty"`
	GraceTime     int            `json:"grace_time,omitempty"`
	TRES          map[string]int `json:"tres,omitempty"`
	MaxTRES       map[string]int `json:"max_tres,omitempty"`
	MinTRES       map[string]int `json:"min_tres,omitempty"`
	IsDefault     bool           `json:"is_default"`
	IsActive      bool           `json:"is_active"`
	Flags         []string       `json:"flags,omitempty"`
	Created       time.Time      `json:"created"`
	Modified      time.Time      `json:"modified"`
}

// UserQuota represents a user's resource quotas.
type UserQuota struct {
	UserName       string                       `json:"user_name"`
	DefaultAccount string                       `json:"default_account"`
	MaxJobs        int                          `json:"max_jobs"`
	MaxSubmitJobs  int                          `json:"max_submit_jobs"`
	MaxWallTime    int                          `json:"max_wall_time"`
	MaxCPUs        int                          `json:"max_cpus"`
	MaxNodes       int                          `json:"max_nodes"`
	MaxMemory      int                          `json:"max_memory"`
	TRESLimits     map[string]int               `json:"tres_limits,omitempty"`
	AccountQuotas  map[string]*UserAccountQuota `json:"account_quotas,omitempty"`
	QoSLimits      map[string]*QoSLimits        `json:"qos_limits,omitempty"`
	GraceTime      int                          `json:"grace_time,omitempty"`
	CurrentUsage   *UserUsage                   `json:"current_usage,omitempty"`
	IsActive       bool                         `json:"is_active"`
	Enforcement    string                       `json:"enforcement"`
}

// UserAccountQuota represents user-account specific quotas.
type UserAccountQuota struct {
	AccountName   string         `json:"account_name"`
	MaxJobs       int            `json:"max_jobs"`
	MaxSubmitJobs int            `json:"max_submit_jobs"`
	MaxWallTime   int            `json:"max_wall_time"`
	TRESLimits    map[string]int `json:"tres_limits,omitempty"`
	Priority      int            `json:"priority"`
	QoS           []string       `json:"qos,omitempty"`
	DefaultQoS    string         `json:"default_qos,omitempty"`
}

// UserFairShare represents user fairshare information.
type UserFairShare struct {
	UserName         string              `json:"user_name"`
	Account          string              `json:"account"`
	Cluster          string              `json:"cluster"`
	Partition        string              `json:"partition,omitempty"`
	FairShareFactor  float64             `json:"fair_share_factor"`
	NormalizedShares float64             `json:"normalized_shares"`
	EffectiveUsage   float64             `json:"effective_usage"`
	FairShareTree    *FairShareNode      `json:"fair_share_tree,omitempty"`
	PriorityFactors  *JobPriorityFactors `json:"priority_factors,omitempty"`
	RawShares        int                 `json:"raw_shares"`
	NormalizedUsage  float64             `json:"normalized_usage"`
	Level            int                 `json:"level"`
	LastDecay        time.Time           `json:"last_decay"`
}

// FairShareNode represents a node in the fairshare tree.
type FairShareNode struct {
	Name             string           `json:"name"`
	Account          string           `json:"account,omitempty"`
	User             string           `json:"user,omitempty"`
	Parent           string           `json:"parent,omitempty"`
	Shares           int              `json:"shares"`
	NormalizedShares float64          `json:"normalized_shares"`
	Usage            float64          `json:"usage"`
	FairShareFactor  float64          `json:"fair_share_factor"`
	Level            int              `json:"level"`
	Children         []*FairShareNode `json:"children,omitempty"`
}

// JobPriorityFactors represents job priority factors.
type JobPriorityFactors struct {
	Age       int              `json:"age"`
	FairShare int              `json:"fair_share"`
	JobSize   int              `json:"job_size"`
	Partition int              `json:"partition"`
	QoS       int              `json:"qos"`
	TRES      int              `json:"tres"`
	Site      int              `json:"site"`
	Nice      int              `json:"nice"`
	Assoc     int              `json:"assoc"`
	Total     int              `json:"total"`
	Weights   *PriorityWeights `json:"weights,omitempty"`
}

// PriorityWeights represents priority weight configuration.
type PriorityWeights struct {
	Age       int `json:"age"`
	FairShare int `json:"fair_share"`
	JobSize   int `json:"job_size"`
	Partition int `json:"partition"`
	QoS       int `json:"qos"`
	TRES      int `json:"tres"`
	Site      int `json:"site"`
	Nice      int `json:"nice"`
	Assoc     int `json:"assoc"`
}

// JobPriorityInfo represents job priority information.
type JobPriorityInfo struct {
	JobID           uint32              `json:"job_id,omitempty"`
	UserName        string              `json:"user_name"`
	Account         string              `json:"account"`
	Partition       string              `json:"partition"`
	QoS             string              `json:"qos"`
	Priority        int                 `json:"priority"`
	Factors         *JobPriorityFactors `json:"factors"`
	Age             int                 `json:"age"`
	EligibleTime    time.Time           `json:"eligible_time"`
	EstimatedStart  time.Time           `json:"estimated_start"`
	PositionInQueue int                 `json:"position_in_queue"`
	PriorityTier    string              `json:"priority_tier"`
}

// AccountFairShare represents account fairshare information.
type AccountFairShare struct {
	AccountName      string              `json:"account_name"`
	Cluster          string              `json:"cluster"`
	Parent           string              `json:"parent,omitempty"`
	Shares           int                 `json:"shares"`
	RawShares        int                 `json:"raw_shares"`
	NormalizedShares float64             `json:"normalized_shares"`
	Usage            float64             `json:"usage"`
	EffectiveUsage   float64             `json:"effective_usage"`
	FairShareFactor  float64             `json:"fair_share_factor"`
	Level            int                 `json:"level"`
	LevelShares      int                 `json:"level_shares"`
	UserCount        int                 `json:"user_count"`
	ActiveUsers      int                 `json:"active_users"`
	JobCount         int                 `json:"job_count"`
	Children         []*AccountFairShare `json:"children,omitempty"`
	Users            []*UserFairShare    `json:"users,omitempty"`
	LastDecay        time.Time           `json:"last_decay"`
	Created          time.Time           `json:"created"`
	Modified         time.Time           `json:"modified"`
}

// FairShareHierarchy represents complete fairshare hierarchy.
type FairShareHierarchy struct {
	Cluster       string              `json:"cluster"`
	RootAccount   string              `json:"root_account"`
	Tree          *FairShareNode      `json:"tree"`
	TotalShares   int                 `json:"total_shares"`
	TotalUsage    float64             `json:"total_usage"`
	LastUpdate    time.Time           `json:"last_update"`
	DecayHalfLife int                 `json:"decay_half_life"`
	UsageWindow   int                 `json:"usage_window"`
	Algorithm     string              `json:"algorithm"`
	Accounts      []*AccountFairShare `json:"accounts,omitempty"`
	Users         []*UserFairShare    `json:"users,omitempty"`
}

// UserAccountAssociation represents user-account association details.
type UserAccountAssociation struct {
	UserName        string                 `json:"user_name"`
	AccountName     string                 `json:"account_name"`
	Cluster         string                 `json:"cluster"`
	Partition       string                 `json:"partition,omitempty"`
	Role            string                 `json:"role"`
	Permissions     []string               `json:"permissions"`
	IsDefault       bool                   `json:"is_default"`
	IsActive        bool                   `json:"is_active"`
	IsCoordinator   bool                   `json:"is_coordinator"`
	MaxJobs         int                    `json:"max_jobs,omitempty"`
	MaxSubmitJobs   int                    `json:"max_submit_jobs,omitempty"`
	MaxWallTime     int                    `json:"max_wall_time,omitempty"`
	Priority        int                    `json:"priority,omitempty"`
	QoS             []string               `json:"qos,omitempty"`
	DefaultQoS      string                 `json:"default_qos,omitempty"`
	TRESLimits      map[string]int         `json:"tres_limits,omitempty"`
	SharesRaw       int                    `json:"shares_raw,omitempty"`
	FairShareFactor float64                `json:"fair_share_factor,omitempty"`
	GraceTime       int                    `json:"grace_time,omitempty"`
	Created         time.Time              `json:"created"`
	Modified        time.Time              `json:"modified"`
	LastAccessed    time.Time              `json:"last_accessed"`
	Flags           []string               `json:"flags,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// UserAccessValidation represents user access validation result.
type UserAccessValidation struct {
	UserName       string                  `json:"user_name"`
	AccountName    string                  `json:"account_name"`
	HasAccess      bool                    `json:"has_access"`
	AccessLevel    string                  `json:"access_level"`
	Permissions    []string                `json:"permissions"`
	Restrictions   []string                `json:"restrictions,omitempty"`
	Reason         string                  `json:"reason,omitempty"`
	ValidFrom      time.Time               `json:"valid_from"`
	ValidUntil     *time.Time              `json:"valid_until,omitempty"`
	Association    *UserAccountAssociation `json:"association,omitempty"`
	QuotaLimits    *UserAccountQuota       `json:"quota_limits,omitempty"`
	CurrentUsage   *AccountUsageStats      `json:"current_usage,omitempty"`
	ValidationTime time.Time               `json:"validation_time"`
}

// ListUserAccountAssociationsOptions configures listing user-account associations.
type ListUserAccountAssociationsOptions struct {
	Accounts         []string `json:"accounts,omitempty"`
	Clusters         []string `json:"clusters,omitempty"`
	Partitions       []string `json:"partitions,omitempty"`
	Roles            []string `json:"roles,omitempty"`
	Permissions      []string `json:"permissions,omitempty"`
	ActiveOnly       bool     `json:"active_only,omitempty"`
	DefaultOnly      bool     `json:"default_only,omitempty"`
	CoordinatorRoles bool     `json:"coordinator_roles,omitempty"`
	WithQuotas       bool     `json:"with_quotas,omitempty"`
	WithUsage        bool     `json:"with_usage,omitempty"`
	WithFairShare    bool     `json:"with_fair_share,omitempty"`
	Limit            int      `json:"limit,omitempty"`
	Offset           int      `json:"offset,omitempty"`
	SortBy           string   `json:"sort_by,omitempty"`
	SortOrder        string   `json:"sort_order,omitempty"`
}

// ============================================================================
// Configuration Types
// ============================================================================

// ClientConfig represents client configuration.
// HTTPDoer is the interface that wraps the Do method for HTTP requests.
// This matches the HttpRequestDoer interface expected by the OpenAPI clients.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientConfig struct {
	BaseURL    string
	HTTPClient HTTPDoer
	Debug      bool
}

// ============================================================================
// Forward Declarations for Analytics Types
// These are defined in pkg/analytics but referenced here for compatibility.
// ============================================================================

// AccountUsageStats is defined in pkg/analytics.
// This is a forward reference for types that depend on it.
type AccountUsageStats struct {
	AccountName      string             `json:"account_name"`
	JobCount         int                `json:"job_count"`
	CPUHours         float64            `json:"cpu_hours"`
	WallHours        float64            `json:"wall_hours"`
	TRESUsage        map[string]float64 `json:"tres_usage,omitempty"`
	AverageQueueTime float64            `json:"average_queue_time"`
	AverageRunTime   float64            `json:"average_run_time"`
	Efficiency       float64            `json:"efficiency"`
}
