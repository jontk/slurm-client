package factory

import (
	"context"
	"time"
)

// SlurmClient represents a version-agnostic Slurm REST API client interface
// This is a copy of the interface from the main package to avoid import cycles
type SlurmClient interface {
	Version() string
	Jobs() JobManager
	Nodes() NodeManager
	Partitions() PartitionManager
	Info() InfoManager
	Close() error
}

// Supporting interfaces (copied to avoid import cycles)

type JobManager interface {
	List(ctx context.Context, opts *ListJobsOptions) (*JobList, error)
	Get(ctx context.Context, jobID string) (*Job, error)
	Submit(ctx context.Context, job *JobSubmission) (*JobSubmitResponse, error)
	Cancel(ctx context.Context, jobID string) error
	Update(ctx context.Context, jobID string, update *JobUpdate) error
	Steps(ctx context.Context, jobID string) (*JobStepList, error)
	Watch(ctx context.Context, opts *WatchJobsOptions) (<-chan JobEvent, error)
}

type NodeManager interface {
	List(ctx context.Context, opts *ListNodesOptions) (*NodeList, error)
	Get(ctx context.Context, nodeName string) (*Node, error)
	Update(ctx context.Context, nodeName string, update *NodeUpdate) error
	Drain(ctx context.Context, nodeName string, reason string) error
	Resume(ctx context.Context, nodeName string) error
}

type PartitionManager interface {
	List(ctx context.Context) (*PartitionList, error)
	Get(ctx context.Context, partitionName string) (*Partition, error)
	Update(ctx context.Context, partitionName string, update *PartitionUpdate) error
}

type InfoManager interface {
	Ping(ctx context.Context) error
	Version(ctx context.Context) (*VersionInfo, error)
	Configuration(ctx context.Context) (*ClusterConfig, error)
	Statistics(ctx context.Context) (*ClusterStats, error)
}

// Data types (simplified versions to avoid import cycles)

type JobState string
type NodeState string

type Job struct {
	ID          string
	Name        string
	UserID      string
	State       JobState
	Partition   string
	SubmitTime  time.Time
	StartTime   *time.Time
	EndTime     *time.Time
	CPUs        int
	Memory      int
}

type JobList struct {
	Jobs  []Job
	Total int
}

type JobSubmission struct {
	Name      string
	Script    string
	Partition string
	CPUs      int
	Memory    int
	TimeLimit int
}

type JobSubmitResponse struct {
	JobID string
}

type JobUpdate struct {
	TimeLimit *int
	Priority  *int
}

type ListJobsOptions struct {
	UserID    string
	State     JobState
	Partition string
	Limit     int
	Offset    int
}

type JobStepList struct {
	Steps []JobStep
}

type JobStep struct {
	ID    string
	JobID string
	Name  string
	State string
}

type WatchJobsOptions struct {
	UserID string
	State  JobState
}

type JobEvent struct {
	Type     string
	JobID    string
	NewState JobState
}

type Node struct {
	Name  string
	State NodeState
	CPUs  int
}

type NodeList struct {
	Nodes []Node
	Total int
}

type NodeUpdate struct {
	State  *NodeState
	Reason *string
}

type ListNodesOptions struct {
	State     NodeState
	Partition string
	Features  []string
}

type Partition struct {
	Name        string
	State       string
	TotalCPUs   int
	TotalMemory int
}

type PartitionList struct {
	Partitions []Partition
	Total      int
}

type PartitionUpdate struct {
	State *string
}

type VersionInfo struct {
	Version    string
	APIVersion string
}

type ClusterConfig struct {
	ClusterName string
}

type ClusterStats struct {
	JobsRunning int
}