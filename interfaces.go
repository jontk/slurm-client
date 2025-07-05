package slurm

import (
	"context"
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
	
	// Update updates node properties
	Update(ctx context.Context, nodeName string, update *NodeUpdate) error
	
	// Drain drains a node
	Drain(ctx context.Context, nodeName string, reason string) error
	
	// Resume resumes a drained node
	Resume(ctx context.Context, nodeName string) error
}

// PartitionManager provides version-agnostic partition operations
type PartitionManager interface {
	// List partitions
	List(ctx context.Context) (*PartitionList, error)
	
	// Get retrieves a specific partition by name
	Get(ctx context.Context, partitionName string) (*Partition, error)
	
	// Update updates partition properties (if supported by version)
	Update(ctx context.Context, partitionName string, update *PartitionUpdate) error
}

// InfoManager provides version-agnostic cluster information
type InfoManager interface {
	// Ping tests connectivity to the Slurm REST API
	Ping(ctx context.Context) error
	
	// Version retrieves Slurm version information
	Version(ctx context.Context) (*VersionInfo, error)
	
	// Configuration retrieves cluster configuration
	Configuration(ctx context.Context) (*ClusterConfig, error)
	
	// Statistics retrieves cluster statistics
	Statistics(ctx context.Context) (*ClusterStats, error)
}

// Common data structures with version compatibility

// Job represents a Slurm job (version-agnostic)
type Job struct {
	ID          string            `json:"job_id"`
	Name        string            `json:"name"`
	UserID      string            `json:"user_id"`
	GroupID     string            `json:"group_id"`
	State       JobState          `json:"job_state"`
	Partition   string            `json:"partition"`
	Priority    int               `json:"priority"`
	SubmitTime  time.Time         `json:"submit_time"`
	StartTime   *time.Time        `json:"start_time,omitempty"`
	EndTime     *time.Time        `json:"end_time,omitempty"`
	TimeLimit   int               `json:"time_limit"`
	NodeList    string            `json:"node_list"`
	CPUs        int               `json:"cpus"`
	Memory      int               `json:"memory"`
	WorkingDir  string            `json:"working_directory"`
	Environment map[string]string `json:"environment"`
	
	// Version-specific fields handled via metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// JobState represents job states with version compatibility
type JobState string

const (
	JobStatePending    JobState = "PENDING"
	JobStateRunning    JobState = "RUNNING"
	JobStateCompleted  JobState = "COMPLETED"
	JobStateCancelled  JobState = "CANCELLED"
	JobStateFailed     JobState = "FAILED"
	JobStateTimeout    JobState = "TIMEOUT"
	JobStateSuspended  JobState = "SUSPENDED"
)

// JobList represents a list of jobs
type JobList struct {
	Jobs     []Job  `json:"jobs"`
	Total    int    `json:"total,omitempty"`
	Offset   int    `json:"offset,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// JobSubmission represents a job submission request
type JobSubmission struct {
	Name        string            `json:"name"`
	Script      string            `json:"script,omitempty"`
	Command     []string          `json:"command,omitempty"`
	Args        []string          `json:"args,omitempty"`
	Partition   string            `json:"partition,omitempty"`
	CPUs        int               `json:"cpus,omitempty"`
	Memory      int               `json:"memory,omitempty"`
	TimeLimit   int               `json:"time_limit,omitempty"`
	WorkingDir  string            `json:"working_directory,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	
	// Version-specific submission options
	Options map[string]interface{} `json:"options,omitempty"`
}

// JobSubmitResponse represents the response from job submission
type JobSubmitResponse struct {
	JobID    string `json:"job_id"`
	Message  string `json:"message,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

// JobUpdate represents job update parameters
type JobUpdate struct {
	TimeLimit *int               `json:"time_limit,omitempty"`
	Priority  *int               `json:"priority,omitempty"`
	Partition *string            `json:"partition,omitempty"`
	Options   map[string]interface{} `json:"options,omitempty"`
}

// ListJobsOptions represents options for listing jobs
type ListJobsOptions struct {
	UserID    string
	State     JobState
	Partition string
	Limit     int
	Offset    int
	StartTime *time.Time
	EndTime   *time.Time
	
	// Version-specific filters
	Filters map[string]interface{}
}

// WatchJobsOptions represents options for watching job events
type WatchJobsOptions struct {
	UserID    string
	State     JobState
	Partition string
	
	// Event filtering
	EventTypes []JobEventType
}

// JobEvent represents a job state change event
type JobEvent struct {
	Type      JobEventType `json:"type"`
	JobID     string       `json:"job_id"`
	OldState  JobState     `json:"old_state,omitempty"`
	NewState  JobState     `json:"new_state"`
	Timestamp time.Time    `json:"timestamp"`
	Message   string       `json:"message,omitempty"`
}

// JobEventType represents the type of job event
type JobEventType string

const (
	JobEventSubmitted JobEventType = "SUBMITTED"
	JobEventStarted   JobEventType = "STARTED"
	JobEventCompleted JobEventType = "COMPLETED"
	JobEventCancelled JobEventType = "CANCELLED"
	JobEventFailed    JobEventType = "FAILED"
)

// Node represents a compute node (version-agnostic)
type Node struct {
	Name         string     `json:"name"`
	State        NodeState  `json:"state"`
	CPUs         int        `json:"cpus"`
	Memory       int        `json:"memory"`
	Features     []string   `json:"features"`
	Partitions   []string   `json:"partitions"`
	Architecture string     `json:"architecture"`
	OS           string     `json:"os"`
	Reason       string     `json:"reason,omitempty"`
	LastBusy     *time.Time `json:"last_busy,omitempty"`
	
	// Version-specific fields
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NodeState represents node states with version compatibility
type NodeState string

const (
	NodeStateIdle     NodeState = "IDLE"
	NodeStateAllocated NodeState = "ALLOCATED"
	NodeStateDrain    NodeState = "DRAIN"
	NodeStateDown     NodeState = "DOWN"
	NodeStateMix      NodeState = "MIX"
)

// NodeList represents a list of nodes
type NodeList struct {
	Nodes    []Node `json:"nodes"`
	Total    int    `json:"total,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NodeUpdate represents node update parameters
type NodeUpdate struct {
	State    *NodeState `json:"state,omitempty"`
	Reason   *string    `json:"reason,omitempty"`
	Features *[]string  `json:"features,omitempty"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

// ListNodesOptions represents options for listing nodes
type ListNodesOptions struct {
	State     NodeState
	Partition string
	Features  []string
	
	// Version-specific filters
	Filters map[string]interface{}
}

// Partition represents a job partition (version-agnostic)
type Partition struct {
	Name             string   `json:"name"`
	State            string   `json:"state"`
	Nodes            []string `json:"nodes"`
	TotalCPUs        int      `json:"total_cpus"`
	TotalMemory      int      `json:"total_memory"`
	MaxTimeLimit     int      `json:"max_time_limit"`
	DefaultTimeLimit int      `json:"default_time_limit"`
	Priority         int      `json:"priority"`
	
	// Version-specific fields
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// PartitionList represents a list of partitions
type PartitionList struct {
	Partitions []Partition `json:"partitions"`
	Total      int         `json:"total,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// PartitionUpdate represents partition update parameters
type PartitionUpdate struct {
	State            *string `json:"state,omitempty"`
	MaxTimeLimit     *int    `json:"max_time_limit,omitempty"`
	DefaultTimeLimit *int    `json:"default_time_limit,omitempty"`
	Priority         *int    `json:"priority,omitempty"`
	Options          map[string]interface{} `json:"options,omitempty"`
}

// JobStep represents a job step (version-agnostic)
type JobStep struct {
	ID        string     `json:"step_id"`
	JobID     string     `json:"job_id"`
	Name      string     `json:"name"`
	State     string     `json:"state"`
	CPUs      int        `json:"cpus"`
	Memory    int        `json:"memory"`
	NodeList  string     `json:"node_list"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	
	// Version-specific fields
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// JobStepList represents a list of job steps
type JobStepList struct {
	Steps    []JobStep `json:"steps"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// VersionInfo represents Slurm version information
type VersionInfo struct {
	Version string `json:"version"`
	Release string `json:"release"`
	APIVersion string `json:"api_version"`
	BuildInfo map[string]string `json:"build_info,omitempty"`
}

// ClusterConfig represents cluster configuration
type ClusterConfig struct {
	ClusterName string `json:"cluster_name"`
	NodeCount   int    `json:"node_count"`
	CPUCount    int    `json:"cpu_count"`
	MemoryTotal int    `json:"memory_total"`
	
	// Version-specific configuration
	Configuration map[string]interface{} `json:"configuration,omitempty"`
}

// ClusterStats represents cluster statistics
type ClusterStats struct {
	JobsRunning  int `json:"jobs_running"`
	JobsPending  int `json:"jobs_pending"`
	NodesIdle    int `json:"nodes_idle"`
	NodesAllocated int `json:"nodes_allocated"`
	
	// Version-specific statistics
	Statistics map[string]interface{} `json:"statistics,omitempty"`
}