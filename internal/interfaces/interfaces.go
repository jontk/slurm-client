package interfaces

import (
	"context"
	"net/http"
	"time"
)

// SlurmClient represents a version-agnostic Slurm REST API client
type SlurmClient interface {
	// Version returns the API version this client supports
	Version() string

	// Jobs returns the JobManager for this version
	Jobs() JobManager

	// Nodes returns the NodeManager for this version
	Nodes() NodeManager

	// Partitions returns the PartitionManager for this version
	Partitions() PartitionManager

	// Info returns the InfoManager for this version
	Info() InfoManager

	// Reservations returns the ReservationManager for this version (v0.0.43+)
	Reservations() ReservationManager

	// QoS returns the QoSManager for this version (v0.0.43+)
	QoS() QoSManager

	// Close closes the client and any resources
	Close() error
}

// JobManager provides version-agnostic job operations
type JobManager interface {
	// List jobs with optional filtering
	List(ctx context.Context, opts *ListJobsOptions) (*JobList, error)

	// Get retrieves a specific job by ID
	Get(ctx context.Context, jobID string) (*Job, error)

	// Submit submits a new job
	Submit(ctx context.Context, job *JobSubmission) (*JobSubmitResponse, error)

	// Cancel cancels a job
	Cancel(ctx context.Context, jobID string) error

	// Update updates job properties (if supported by version)
	Update(ctx context.Context, jobID string, update *JobUpdate) error

	// Steps retrieves job steps for a job
	Steps(ctx context.Context, jobID string) (*JobStepList, error)

	// Watch provides real-time job updates (if supported by version)
	Watch(ctx context.Context, opts *WatchJobsOptions) (<-chan JobEvent, error)
}

// NodeManager provides version-agnostic node operations
type NodeManager interface {
	// List nodes with optional filtering
	List(ctx context.Context, opts *ListNodesOptions) (*NodeList, error)

	// Get retrieves a specific node by name
	Get(ctx context.Context, nodeName string) (*Node, error)

	// Update updates node properties (if supported by version)
	Update(ctx context.Context, nodeName string, update *NodeUpdate) error

	// Watch provides real-time node updates (if supported by version)
	Watch(ctx context.Context, opts *WatchNodesOptions) (<-chan NodeEvent, error)
}

// PartitionManager provides version-agnostic partition operations
type PartitionManager interface {
	// List partitions with optional filtering
	List(ctx context.Context, opts *ListPartitionsOptions) (*PartitionList, error)

	// Get retrieves a specific partition by name
	Get(ctx context.Context, partitionName string) (*Partition, error)

	// Update updates partition properties (if supported by version)
	Update(ctx context.Context, partitionName string, update *PartitionUpdate) error

	// Watch provides real-time partition updates (if supported by version)
	Watch(ctx context.Context, opts *WatchPartitionsOptions) (<-chan PartitionEvent, error)
}

// InfoManager provides version-agnostic cluster information operations
type InfoManager interface {
	// Get retrieves cluster information
	Get(ctx context.Context) (*ClusterInfo, error)

	// Ping tests connectivity to the cluster
	Ping(ctx context.Context) error

	// Stats retrieves cluster statistics
	Stats(ctx context.Context) (*ClusterStats, error)

	// Version retrieves API version information
	Version(ctx context.Context) (*APIVersion, error)
}

// ReservationManager provides version-agnostic reservation operations
type ReservationManager interface {
	// List reservations with optional filtering
	List(ctx context.Context, opts *ListReservationsOptions) (*ReservationList, error)

	// Get retrieves a specific reservation by name
	Get(ctx context.Context, reservationName string) (*Reservation, error)

	// Create creates a new reservation
	Create(ctx context.Context, reservation *ReservationCreate) (*ReservationCreateResponse, error)

	// Update updates an existing reservation
	Update(ctx context.Context, reservationName string, update *ReservationUpdate) error

	// Delete deletes a reservation
	Delete(ctx context.Context, reservationName string) error
}

// QoSManager provides version-agnostic QoS (Quality of Service) operations
type QoSManager interface {
	// List QoS with optional filtering
	List(ctx context.Context, opts *ListQoSOptions) (*QoSList, error)

	// Get retrieves a specific QoS by name
	Get(ctx context.Context, qosName string) (*QoS, error)

	// Create creates a new QoS
	Create(ctx context.Context, qos *QoSCreate) (*QoSCreateResponse, error)

	// Update updates an existing QoS
	Update(ctx context.Context, qosName string, update *QoSUpdate) error

	// Delete deletes a QoS
	Delete(ctx context.Context, qosName string) error
}

// Common data structures (version-agnostic)

// Job represents a job in the cluster
type Job struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	UserID      string                 `json:"user_id"`
	GroupID     string                 `json:"group_id"`
	State       string                 `json:"state"`
	Partition   string                 `json:"partition"`
	Priority    int                    `json:"priority"`
	SubmitTime  time.Time              `json:"submit_time"`
	StartTime   *time.Time             `json:"start_time,omitempty"`
	EndTime     *time.Time             `json:"end_time,omitempty"`
	CPUs        int                    `json:"cpus"`
	Memory      int                    `json:"memory"`
	TimeLimit   int                    `json:"time_limit"`
	WorkingDir  string                 `json:"working_dir"`
	Command     string                 `json:"command"`
	Environment map[string]string      `json:"environment"`
	Nodes       []string               `json:"nodes"`
	ExitCode    int                    `json:"exit_code"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// JobList represents a list of jobs
type JobList struct {
	Jobs  []Job `json:"jobs"`
	Total int   `json:"total"`
}

// JobSubmission represents a job submission request
type JobSubmission struct {
	Name        string            `json:"name"`
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

// JobSubmitResponse represents the response from job submission
type JobSubmitResponse struct {
	JobID string `json:"job_id"`
}

// JobUpdate represents a job update request
type JobUpdate struct {
	Priority  *int    `json:"priority,omitempty"`
	TimeLimit *int    `json:"time_limit,omitempty"`
	Name      *string `json:"name,omitempty"`
}

// JobStepList represents a list of job steps
type JobStepList struct {
	Steps []JobStep `json:"steps"`
	Total int       `json:"total"`
}

// JobStep represents a job step
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

// JobEvent represents a job state change event
type JobEvent struct {
	Type      string    `json:"type"`
	JobID     string    `json:"job_id"`
	OldState  string    `json:"old_state"`
	NewState  string    `json:"new_state"`
	Timestamp time.Time `json:"timestamp"`
	Job       *Job      `json:"job,omitempty"`
	Error     error     `json:"error,omitempty"`
}

// Node represents a compute node
type Node struct {
	Name         string                 `json:"name"`
	State        string                 `json:"state"`
	CPUs         int                    `json:"cpus"`
	Memory       int                    `json:"memory"`
	Partitions   []string               `json:"partitions"`
	Features     []string               `json:"features"`
	Reason       string                 `json:"reason,omitempty"`
	LastBusy     *time.Time             `json:"last_busy,omitempty"`
	Architecture string                 `json:"architecture,omitempty"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// NodeList represents a list of nodes
type NodeList struct {
	Nodes []Node `json:"nodes"`
	Total int    `json:"total"`
}

// NodeUpdate represents a node update request
type NodeUpdate struct {
	State    *string  `json:"state,omitempty"`
	Reason   *string  `json:"reason,omitempty"`
	Features []string `json:"features,omitempty"`
}

// NodeEvent represents a node state change event
type NodeEvent struct {
	Type      string    `json:"type"`
	NodeName  string    `json:"node_name"`
	OldState  string    `json:"old_state"`
	NewState  string    `json:"new_state"`
	Timestamp time.Time `json:"timestamp"`
	Node      *Node     `json:"node,omitempty"`
	Error     error     `json:"error,omitempty"`
}

// Partition represents a job partition
type Partition struct {
	Name           string   `json:"name"`
	State          string   `json:"state"`
	TotalNodes     int      `json:"total_nodes"`
	AvailableNodes int      `json:"available_nodes"`
	TotalCPUs      int      `json:"total_cpus"`
	IdleCPUs       int      `json:"idle_cpus"`
	MaxTime        int      `json:"max_time"`
	DefaultTime    int      `json:"default_time"`
	MaxMemory      int      `json:"max_memory"`
	DefaultMemory  int      `json:"default_memory"`
	AllowedUsers   []string `json:"allowed_users"`
	DeniedUsers    []string `json:"denied_users"`
	AllowedGroups  []string `json:"allowed_groups"`
	DeniedGroups   []string `json:"denied_groups"`
	Priority       int      `json:"priority"`
	Nodes          []string `json:"nodes"`
}

// PartitionList represents a list of partitions
type PartitionList struct {
	Partitions []Partition `json:"partitions"`
	Total      int         `json:"total"`
}

// PartitionUpdate represents a partition update request
type PartitionUpdate struct {
	State         *string  `json:"state,omitempty"`
	MaxTime       *int     `json:"max_time,omitempty"`
	DefaultTime   *int     `json:"default_time,omitempty"`
	MaxMemory     *int     `json:"max_memory,omitempty"`
	DefaultMemory *int     `json:"default_memory,omitempty"`
	AllowedUsers  []string `json:"allowed_users,omitempty"`
	DeniedUsers   []string `json:"denied_users,omitempty"`
	Priority      *int     `json:"priority,omitempty"`
}

// PartitionEvent represents a partition state change event
type PartitionEvent struct {
	Type          string     `json:"type"`
	PartitionName string     `json:"partition_name"`
	OldState      string     `json:"old_state"`
	NewState      string     `json:"new_state"`
	Timestamp     time.Time  `json:"timestamp"`
	Partition     *Partition `json:"partition,omitempty"`
	Error         error      `json:"error,omitempty"`
}

// ClusterInfo represents cluster information
type ClusterInfo struct {
	Version     string `json:"version"`
	Release     string `json:"release"`
	ClusterName string `json:"cluster_name"`
	APIVersion  string `json:"api_version"`
	Uptime      int    `json:"uptime"`
}

// ClusterStats represents cluster statistics
type ClusterStats struct {
	TotalNodes     int `json:"total_nodes"`
	IdleNodes      int `json:"idle_nodes"`
	AllocatedNodes int `json:"allocated_nodes"`
	TotalCPUs      int `json:"total_cpus"`
	IdleCPUs       int `json:"idle_cpus"`
	AllocatedCPUs  int `json:"allocated_cpus"`
	TotalJobs      int `json:"total_jobs"`
	RunningJobs    int `json:"running_jobs"`
	PendingJobs    int `json:"pending_jobs"`
	CompletedJobs  int `json:"completed_jobs"`
}

// APIVersion represents API version information
type APIVersion struct {
	Version     string `json:"version"`
	Release     string `json:"release"`
	Description string `json:"description"`
	Deprecated  bool   `json:"deprecated"`
}

// ExtendedDiagnostics represents detailed diagnostic information from the cluster
type ExtendedDiagnostics struct {
	// Basic statistics (same as ClusterStats)
	ClusterStats
	
	// Additional job statistics
	JobsFailed    int `json:"jobs_failed"`
	JobsCanceled  int `json:"jobs_canceled"`
	JobsTimeout   int `json:"jobs_timeout"`
	
	// Backfill scheduler statistics
	BackfillActive         bool  `json:"backfill_active"`
	BackfillJobsTotal      int   `json:"backfill_jobs_total"`
	BackfillJobsRecent     int   `json:"backfill_jobs_recent"`
	BackfillCycleCount     int   `json:"backfill_cycle_count"`
	BackfillCycleMeanTime  int64 `json:"backfill_cycle_mean_time"`
	
	// Server performance statistics
	ServerThreadCount      int   `json:"server_thread_count"`
	AgentQueueSize         int   `json:"agent_queue_size"`
	ScheduleCycleMax       int   `json:"schedule_cycle_max"`
	ScheduleCycleLast      int   `json:"schedule_cycle_last"`
	ScheduleCycleMean      int64 `json:"schedule_cycle_mean"`
	
	// RPC statistics
	RPCsTotal              int   `json:"rpcs_total"`
	RPCsPending            int   `json:"rpcs_pending"`
	RPCsCompleted          int   `json:"rpcs_completed"`
	
	// Additional metadata
	DiagTime               time.Time              `json:"diag_time"`
	RawData                map[string]interface{} `json:"raw_data,omitempty"`
}

// Reservation represents a resource reservation
type Reservation struct {
	Name         string            `json:"name"`
	State        string            `json:"state"`
	StartTime    time.Time         `json:"start_time"`
	EndTime      time.Time         `json:"end_time"`
	Duration     int               `json:"duration"`
	Nodes        []string          `json:"nodes"`
	NodeCount    int               `json:"node_count"`
	CoreCount    int               `json:"core_count"`
	Users        []string          `json:"users"`
	Accounts     []string          `json:"accounts"`
	Flags        []string          `json:"flags"`
	Features     []string          `json:"features"`
	PartitionName string           `json:"partition_name"`
	Licenses     map[string]int    `json:"licenses"`
	BurstBuffer  string            `json:"burst_buffer"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// ReservationList represents a list of reservations
type ReservationList struct {
	Reservations []Reservation `json:"reservations"`
	Total        int           `json:"total"`
}

// ReservationCreate represents a reservation creation request
type ReservationCreate struct {
	Name         string            `json:"name"`
	StartTime    time.Time         `json:"start_time"`
	EndTime      time.Time         `json:"end_time,omitempty"`
	Duration     int               `json:"duration,omitempty"`
	Nodes        []string          `json:"nodes,omitempty"`
	NodeCount    int               `json:"node_count,omitempty"`
	CoreCount    int               `json:"core_count,omitempty"`
	Users        []string          `json:"users,omitempty"`
	Accounts     []string          `json:"accounts,omitempty"`
	Flags        []string          `json:"flags,omitempty"`
	Features     []string          `json:"features,omitempty"`
	PartitionName string           `json:"partition_name,omitempty"`
	Licenses     map[string]int    `json:"licenses,omitempty"`
	BurstBuffer  string            `json:"burst_buffer,omitempty"`
}

// ReservationCreateResponse represents the response from reservation creation
type ReservationCreateResponse struct {
	ReservationName string `json:"reservation_name"`
}

// ReservationUpdate represents a reservation update request
type ReservationUpdate struct {
	StartTime    *time.Time        `json:"start_time,omitempty"`
	EndTime      *time.Time        `json:"end_time,omitempty"`
	Duration     *int              `json:"duration,omitempty"`
	Nodes        []string          `json:"nodes,omitempty"`
	NodeCount    *int              `json:"node_count,omitempty"`
	Users        []string          `json:"users,omitempty"`
	Accounts     []string          `json:"accounts,omitempty"`
	Flags        []string          `json:"flags,omitempty"`
	Features     []string          `json:"features,omitempty"`
}

// QoS represents a Quality of Service configuration
type QoS struct {
	Name               string                 `json:"name"`
	Description        string                 `json:"description"`
	Priority           int                    `json:"priority"`
	PreemptMode        string                 `json:"preempt_mode"`
	GraceTime          int                    `json:"grace_time"`
	MaxJobs            int                    `json:"max_jobs"`
	MaxJobsPerUser     int                    `json:"max_jobs_per_user"`
	MaxJobsPerAccount  int                    `json:"max_jobs_per_account"`
	MaxSubmitJobs      int                    `json:"max_submit_jobs"`
	MaxCPUs            int                    `json:"max_cpus"`
	MaxCPUsPerUser     int                    `json:"max_cpus_per_user"`
	MaxNodes           int                    `json:"max_nodes"`
	MaxWallTime        int                    `json:"max_wall_time"`
	MinCPUs            int                    `json:"min_cpus"`
	MinNodes           int                    `json:"min_nodes"`
	UsageFactor        float64                `json:"usage_factor"`
	UsageThreshold     float64                `json:"usage_threshold"`
	Flags              []string               `json:"flags"`
	AllowedAccounts    []string               `json:"allowed_accounts"`
	DeniedAccounts     []string               `json:"denied_accounts"`
	AllowedUsers       []string               `json:"allowed_users"`
	DeniedUsers        []string               `json:"denied_users"`
	Metadata           map[string]interface{} `json:"metadata"`
}

// QoSList represents a list of QoS configurations
type QoSList struct {
	QoS   []QoS `json:"qos"`
	Total int   `json:"total"`
}

// QoSCreate represents a QoS creation request
type QoSCreate struct {
	Name               string   `json:"name"`
	Description        string   `json:"description,omitempty"`
	Priority           int      `json:"priority,omitempty"`
	PreemptMode        string   `json:"preempt_mode,omitempty"`
	GraceTime          int      `json:"grace_time,omitempty"`
	MaxJobs            int      `json:"max_jobs,omitempty"`
	MaxJobsPerUser     int      `json:"max_jobs_per_user,omitempty"`
	MaxJobsPerAccount  int      `json:"max_jobs_per_account,omitempty"`
	MaxSubmitJobs      int      `json:"max_submit_jobs,omitempty"`
	MaxCPUs            int      `json:"max_cpus,omitempty"`
	MaxCPUsPerUser     int      `json:"max_cpus_per_user,omitempty"`
	MaxNodes           int      `json:"max_nodes,omitempty"`
	MaxWallTime        int      `json:"max_wall_time,omitempty"`
	MinCPUs            int      `json:"min_cpus,omitempty"`
	MinNodes           int      `json:"min_nodes,omitempty"`
	UsageFactor        float64  `json:"usage_factor,omitempty"`
	UsageThreshold     float64  `json:"usage_threshold,omitempty"`
	Flags              []string `json:"flags,omitempty"`
	AllowedAccounts    []string `json:"allowed_accounts,omitempty"`
	DeniedAccounts     []string `json:"denied_accounts,omitempty"`
	AllowedUsers       []string `json:"allowed_users,omitempty"`
	DeniedUsers        []string `json:"denied_users,omitempty"`
}

// QoSCreateResponse represents the response from QoS creation
type QoSCreateResponse struct {
	QoSName string `json:"qos_name"`
}

// QoSUpdate represents a QoS update request
type QoSUpdate struct {
	Description        *string   `json:"description,omitempty"`
	Priority           *int      `json:"priority,omitempty"`
	PreemptMode        *string   `json:"preempt_mode,omitempty"`
	GraceTime          *int      `json:"grace_time,omitempty"`
	MaxJobs            *int      `json:"max_jobs,omitempty"`
	MaxJobsPerUser     *int      `json:"max_jobs_per_user,omitempty"`
	MaxJobsPerAccount  *int      `json:"max_jobs_per_account,omitempty"`
	MaxSubmitJobs      *int      `json:"max_submit_jobs,omitempty"`
	MaxCPUs            *int      `json:"max_cpus,omitempty"`
	MaxCPUsPerUser     *int      `json:"max_cpus_per_user,omitempty"`
	MaxNodes           *int      `json:"max_nodes,omitempty"`
	MaxWallTime        *int      `json:"max_wall_time,omitempty"`
	MinCPUs            *int      `json:"min_cpus,omitempty"`
	MinNodes           *int      `json:"min_nodes,omitempty"`
	UsageFactor        *float64  `json:"usage_factor,omitempty"`
	UsageThreshold     *float64  `json:"usage_threshold,omitempty"`
	Flags              []string  `json:"flags,omitempty"`
	AllowedAccounts    []string  `json:"allowed_accounts,omitempty"`
	DeniedAccounts     []string  `json:"denied_accounts,omitempty"`
	AllowedUsers       []string  `json:"allowed_users,omitempty"`
	DeniedUsers        []string  `json:"denied_users,omitempty"`
}

// List options for filtering

// ListJobsOptions provides options for listing jobs
type ListJobsOptions struct {
	UserID    string   `json:"user_id,omitempty"`
	States    []string `json:"states,omitempty"`
	Partition string   `json:"partition,omitempty"`
	Limit     int      `json:"limit,omitempty"`
	Offset    int      `json:"offset,omitempty"`
}

// ListNodesOptions provides options for listing nodes
type ListNodesOptions struct {
	States    []string `json:"states,omitempty"`
	Partition string   `json:"partition,omitempty"`
	Features  []string `json:"features,omitempty"`
	Limit     int      `json:"limit,omitempty"`
	Offset    int      `json:"offset,omitempty"`
}

// ListPartitionsOptions provides options for listing partitions
type ListPartitionsOptions struct {
	States []string `json:"states,omitempty"`
	Limit  int      `json:"limit,omitempty"`
	Offset int      `json:"offset,omitempty"`
}

// ListReservationsOptions provides options for listing reservations
type ListReservationsOptions struct {
	Names    []string `json:"names,omitempty"`
	Users    []string `json:"users,omitempty"`
	Accounts []string `json:"accounts,omitempty"`
	States   []string `json:"states,omitempty"`
	Limit    int      `json:"limit,omitempty"`
	Offset   int      `json:"offset,omitempty"`
}

// ListQoSOptions provides options for listing QoS
type ListQoSOptions struct {
	Names    []string `json:"names,omitempty"`
	Accounts []string `json:"accounts,omitempty"`
	Users    []string `json:"users,omitempty"`
	Limit    int      `json:"limit,omitempty"`
	Offset   int      `json:"offset,omitempty"`
}

// Watch options for real-time updates

// WatchJobsOptions provides options for watching job changes
type WatchJobsOptions struct {
	UserID          string   `json:"user_id,omitempty"`
	States          []string `json:"states,omitempty"`
	Partition       string   `json:"partition,omitempty"`
	JobIDs          []string `json:"job_ids,omitempty"`
	ExcludeNew      bool     `json:"exclude_new,omitempty"`
	ExcludeCompleted bool    `json:"exclude_completed,omitempty"`
}

// WatchNodesOptions provides options for watching node changes
type WatchNodesOptions struct {
	States    []string `json:"states,omitempty"`
	Partition string   `json:"partition,omitempty"`
	Features  []string `json:"features,omitempty"`
	NodeNames []string `json:"node_names,omitempty"`
}

// WatchPartitionsOptions provides options for watching partition changes
type WatchPartitionsOptions struct {
	States         []string `json:"states,omitempty"`
	PartitionNames []string `json:"partition_names,omitempty"`
}

// ClientConfig holds configuration for the API client
type ClientConfig struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
	Debug      bool
}
