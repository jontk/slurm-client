package types

import (
	"time"
)

// Job represents a SLURM job with common fields across all API versions
type Job struct {
	JobID             int32             `json:"job_id"`
	Name              string            `json:"name"`
	UserID            int32             `json:"user_id"`
	UserName          string            `json:"user_name"`
	GroupID           int32             `json:"group_id"`
	Account           string            `json:"account"`
	Partition         string            `json:"partition"`
	QoS               string            `json:"qos"`
	State             JobState          `json:"state"`
	StateReason       string            `json:"state_reason"`
	TimeLimit         int32             `json:"time_limit"`
	SubmitTime        time.Time         `json:"submit_time"`
	StartTime         *time.Time        `json:"start_time,omitempty"`
	EndTime           *time.Time        `json:"end_time,omitempty"`
	Priority          int32             `json:"priority"`
	CPUs              int32             `json:"cpus"`
	Nodes             int32             `json:"nodes"`
	NodeList          string            `json:"node_list"`
	Command           string            `json:"command"`
	WorkingDirectory  string            `json:"working_directory"`
	StandardInput     string            `json:"standard_input"`
	StandardOutput    string            `json:"standard_output"`
	StandardError     string            `json:"standard_error"`
	ArrayJobID        *int32            `json:"array_job_id,omitempty"`
	ArrayTaskID       *int32            `json:"array_task_id,omitempty"`
	ArrayTaskString   string            `json:"array_task_string,omitempty"`
	Dependencies      []JobDependency   `json:"dependencies,omitempty"`
	ResourceRequests  ResourceRequests  `json:"resource_requests"`
	Environment       map[string]string `json:"environment,omitempty"`
	MailType          []string          `json:"mail_type,omitempty"`
	MailUser          string            `json:"mail_user,omitempty"`
	ExcludeNodes      string            `json:"exclude_nodes,omitempty"`
	Nice              int32             `json:"nice"`
	Comment           string            `json:"comment,omitempty"`
	BatchHost         string            `json:"batch_host,omitempty"`
	BatchScript       string            `json:"batch_script,omitempty"`
	AllocatingNode    string            `json:"allocating_node,omitempty"`
	ScheduledNodes    string            `json:"scheduled_nodes,omitempty"`
	PreemptTime       *time.Time        `json:"preempt_time,omitempty"`
	SuspendTime       *time.Time        `json:"suspend_time,omitempty"`
	Deadline          *time.Time        `json:"deadline,omitempty"`
	ClusterFeatures   []string          `json:"cluster_features,omitempty"`
	Preference        []string          `json:"preference,omitempty"`
	MinCPUs           int32             `json:"min_cpus"`
	MinMemory         int64             `json:"min_memory"`
	MinTmpDisk        int64             `json:"min_tmp_disk"`
	Features          []string          `json:"features,omitempty"`
	Gres              string            `json:"gres,omitempty"`
	Shared            string            `json:"shared,omitempty"`
	Profile           []string          `json:"profile,omitempty"`
	Reservation       string            `json:"reservation,omitempty"`
	CPUFrequencyMin   int32             `json:"cpu_frequency_min,omitempty"`
	CPUFrequencyMax   int32             `json:"cpu_frequency_max,omitempty"`
	CPUFrequencyGov   string            `json:"cpu_frequency_gov,omitempty"`
}

// JobState represents the state of a job
type JobState string

const (
	JobStatePending     JobState = "PENDING"
	JobStateRunning     JobState = "RUNNING"
	JobStateSuspended   JobState = "SUSPENDED"
	JobStateCompleted   JobState = "COMPLETED"
	JobStateCancelled   JobState = "CANCELLED"
	JobStateFailed      JobState = "FAILED"
	JobStateTimeout     JobState = "TIMEOUT"
	JobStateNodeFail    JobState = "NODE_FAIL"
	JobStatePreempted   JobState = "PREEMPTED"
	JobStateBootFail    JobState = "BOOT_FAIL"
	JobStateDeadline    JobState = "DEADLINE"
	JobStateOutOfMemory JobState = "OUT_OF_MEMORY"
	JobStateCompleting  JobState = "COMPLETING"
	JobStateConfiguring JobState = "CONFIGURING"
	JobStateResizing    JobState = "RESIZING"
	JobStateRevoked     JobState = "REVOKED"
	JobStateStopped     JobState = "STOPPED"
	JobStateSignaling   JobState = "SIGNALING"
	JobStateStageOut    JobState = "STAGE_OUT"
)

// ResourceRequests represents the resource requirements for a job
type ResourceRequests struct {
	Memory       int64  `json:"memory,omitempty"`
	MemoryPerCPU int64  `json:"memory_per_cpu,omitempty"`
	MemoryPerGPU int64  `json:"memory_per_gpu,omitempty"`
	TmpDisk      int64  `json:"tmp_disk,omitempty"`
	CPUsPerTask  int32  `json:"cpus_per_task,omitempty"`
	TasksPerNode int32  `json:"tasks_per_node,omitempty"`
	TasksPerCore int32  `json:"tasks_per_core,omitempty"`
	ThreadsPerCore int32  `json:"threads_per_core,omitempty"`
}

// JobDependency represents a job dependency
type JobDependency struct {
	Type  string   `json:"type"`
	JobIDs []int32 `json:"job_ids,omitempty"`
	State string   `json:"state,omitempty"`
}

// JobCreate represents the data needed to create a new job
type JobCreate struct {
	Name              string            `json:"name,omitempty"`
	Account           string            `json:"account,omitempty"`
	Partition         string            `json:"partition,omitempty"`
	QoS               string            `json:"qos,omitempty"`
	TimeLimit         int32             `json:"time_limit,omitempty"`
	Priority          *int32            `json:"priority,omitempty"`
	CPUs              int32             `json:"cpus,omitempty"`
	Nodes             int32             `json:"nodes,omitempty"`
	Tasks             int32             `json:"tasks,omitempty"`
	Command           string            `json:"command"`
	Script            string            `json:"script,omitempty"`
	WorkingDirectory  string            `json:"working_directory,omitempty"`
	StandardInput     string            `json:"standard_input,omitempty"`
	StandardOutput    string            `json:"standard_output,omitempty"`
	StandardError     string            `json:"standard_error,omitempty"`
	ArrayString       string            `json:"array_string,omitempty"`
	Dependencies      []JobDependency   `json:"dependencies,omitempty"`
	ResourceRequests  ResourceRequests  `json:"resource_requests,omitempty"`
	Environment       map[string]string `json:"environment,omitempty"`
	MailType          []string          `json:"mail_type,omitempty"`
	MailUser          string            `json:"mail_user,omitempty"`
	ExcludeNodes      string            `json:"exclude_nodes,omitempty"`
	Nice              int32             `json:"nice,omitempty"`
	Comment           string            `json:"comment,omitempty"`
	Deadline          *time.Time        `json:"deadline,omitempty"`
	ClusterFeatures   []string          `json:"cluster_features,omitempty"`
	MinCPUs           int32             `json:"min_cpus,omitempty"`
	MinMemory         int64             `json:"min_memory,omitempty"`
	MinTmpDisk        int64             `json:"min_tmp_disk,omitempty"`
	Features          []string          `json:"features,omitempty"`
	Gres              string            `json:"gres,omitempty"`
	Shared            string            `json:"shared,omitempty"`
	Profile           []string          `json:"profile,omitempty"`
	Reservation       string            `json:"reservation,omitempty"`
	CPUFrequencyMin   int32             `json:"cpu_frequency_min,omitempty"`
	CPUFrequencyMax   int32             `json:"cpu_frequency_max,omitempty"`
	CPUFrequencyGov   string            `json:"cpu_frequency_gov,omitempty"`
	Immediate         bool              `json:"immediate,omitempty"`
	BeginTime         *time.Time        `json:"begin_time,omitempty"`
	Burst             bool              `json:"burst,omitempty"`
	Power             string            `json:"power,omitempty"`
	SpreadJob         bool              `json:"spread_job,omitempty"`
	WaitAllNodes      bool              `json:"wait_all_nodes,omitempty"`
	KillOnBadExit     bool              `json:"kill_on_bad_exit,omitempty"`
	MemoryBind        string            `json:"memory_bind,omitempty"`
	CpuBind           string            `json:"cpu_bind,omitempty"`
	NoKill            bool              `json:"no_kill,omitempty"`
	Overcommit        bool              `json:"overcommit,omitempty"`
	Contiguous        bool              `json:"contiguous,omitempty"`
	CoreSpec          string            `json:"core_spec,omitempty"`
	ThreadSpec        int32             `json:"thread_spec,omitempty"`
	Distribution      string            `json:"distribution,omitempty"`
	Switches          int32             `json:"switches,omitempty"`
	WaitTime          int32             `json:"wait_time,omitempty"`
	MinNodes          int32             `json:"min_nodes,omitempty"`
	MaxNodes          int32             `json:"max_nodes,omitempty"`
	SiteFactor        int32             `json:"site_factor,omitempty"`
	RequeuePriority   *int32            `json:"requeue_priority,omitempty"`
	RequeueHold       bool              `json:"requeue_hold,omitempty"`
	RequeueExit       string            `json:"requeue_exit,omitempty"`
	RequeueExitHold   string            `json:"requeue_exit_hold,omitempty"`
}

// JobUpdate represents the data needed to update a job
type JobUpdate struct {
	Name              *string           `json:"name,omitempty"`
	Account           *string           `json:"account,omitempty"`
	Partition         *string           `json:"partition,omitempty"`
	QoS               *string           `json:"qos,omitempty"`
	TimeLimit         *int32            `json:"time_limit,omitempty"`
	Priority          *int32            `json:"priority,omitempty"`
	Nice              *int32            `json:"nice,omitempty"`
	Comment           *string           `json:"comment,omitempty"`
	Deadline          *time.Time        `json:"deadline,omitempty"`
	MinCPUs           *int32            `json:"min_cpus,omitempty"`
	MinMemory         *int64            `json:"min_memory,omitempty"`
	MinTmpDisk        *int64            `json:"min_tmp_disk,omitempty"`
	Features          []string          `json:"features,omitempty"`
	Gres              *string           `json:"gres,omitempty"`
	Shared            *string           `json:"shared,omitempty"`
	BeginTime         *time.Time        `json:"begin_time,omitempty"`
	Burst             *bool             `json:"burst,omitempty"`
	Reservation       *string           `json:"reservation,omitempty"`
	ExcludeNodes      *string           `json:"exclude_nodes,omitempty"`
	NodeList          *string           `json:"node_list,omitempty"`
	MinNodes          *int32            `json:"min_nodes,omitempty"`
	MaxNodes          *int32            `json:"max_nodes,omitempty"`
	RequeuePriority   *int32            `json:"requeue_priority,omitempty"`
}

// JobSubmitResponse represents the response from job submission
type JobSubmitResponse struct {
	JobID            int32    `json:"job_id"`
	StepID           string   `json:"step_id,omitempty"`
	JobSubmitUserMsg string   `json:"job_submit_user_msg,omitempty"`
	Error            []string `json:"error,omitempty"`
	Warning          []string `json:"warning,omitempty"`
}

// JobCancelRequest represents the request to cancel a job
type JobCancelRequest struct {
	Signal   string `json:"signal,omitempty"`
	Message  string `json:"message,omitempty"`
	Account  string `json:"account,omitempty"`
	Name     string `json:"name,omitempty"`
	Partition string `json:"partition,omitempty"`
	QoS      string `json:"qos,omitempty"`
	State    string `json:"state,omitempty"`
	UserID   int32  `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
	WaitTime int32  `json:"wait_time,omitempty"`
}

// JobListOptions represents options for listing jobs
type JobListOptions struct {
	Accounts     []string     `json:"accounts,omitempty"`
	Users        []string     `json:"users,omitempty"`
	States       []JobState   `json:"states,omitempty"`
	Partitions   []string     `json:"partitions,omitempty"`
	QoS          []string     `json:"qos,omitempty"`
	JobIDs       []int32      `json:"job_ids,omitempty"`
	JobNames     []string     `json:"job_names,omitempty"`
	StartTime    *time.Time   `json:"start_time,omitempty"`
	EndTime      *time.Time   `json:"end_time,omitempty"`
	Limit        int          `json:"limit,omitempty"`
	Offset       int          `json:"offset,omitempty"`
	IncludeSteps bool         `json:"include_steps,omitempty"`
}

// JobList represents a list of jobs
type JobList struct {
	Jobs  []Job `json:"jobs"`
	Total int   `json:"total"`
}

// JobSignalRequest represents a request to signal a job
type JobSignalRequest struct {
	Signal   string `json:"signal"`
	JobID    int32  `json:"job_id"`
	StepID   string `json:"step_id,omitempty"`
}

// JobHoldRequest represents a request to hold/release a job
type JobHoldRequest struct {
	JobID    int32  `json:"job_id"`
	Hold     bool   `json:"hold"`
	Priority int32  `json:"priority,omitempty"`
}

// JobNotifyRequest represents a request to notify a job
type JobNotifyRequest struct {
	JobID   int32  `json:"job_id"`
	Message string `json:"message"`
}