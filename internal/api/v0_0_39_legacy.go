package slurm

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
	"github.com/jontk/slurm-client/pkg/retry"
)

// Client represents a Slurm REST API client
type Client struct {
	config     *config.Config
	httpClient *http.Client
	auth       auth.Provider
	retry      retry.Policy
	baseURL    string
}

// ClientOption represents a configuration option for the Client
type ClientOption func(*Client) error

// NewClient creates a new Slurm REST API client
func NewClient(options ...ClientOption) (*Client, error) {
	client := &Client{
		config: config.NewDefault(),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		retry: retry.NewExponentialBackoff(),
	}

	for _, option := range options {
		if err := option(client); err != nil {
			return nil, err
		}
	}

	return client, nil
}

// WithConfig sets the client configuration
func WithConfig(cfg *config.Config) ClientOption {
	return func(c *Client) error {
		c.config = cfg
		return nil
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) error {
		c.httpClient = httpClient
		return nil
	}
}

// WithAuth sets the authentication provider
func WithAuth(auth auth.Provider) ClientOption {
	return func(c *Client) error {
		c.auth = auth
		return nil
	}
}

// WithRetryPolicy sets the retry policy
func WithRetryPolicy(policy retry.Policy) ClientOption {
	return func(c *Client) error {
		c.retry = policy
		return nil
	}
}

// WithBaseURL sets the base URL for the Slurm REST API
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		c.baseURL = baseURL
		return nil
	}
}

// JobAPI provides access to job-related operations
type JobAPI interface {
	// ListJobs lists jobs with optional filtering
	ListJobs(ctx context.Context, opts *ListJobsOptions) (*JobList, error)
	
	// GetJob retrieves a specific job by ID
	GetJob(ctx context.Context, jobID string) (*Job, error)
	
	// SubmitJob submits a new job
	SubmitJob(ctx context.Context, job *JobSubmission) (*JobSubmitResponse, error)
	
	// CancelJob cancels a job
	CancelJob(ctx context.Context, jobID string) error
	
	// GetJobSteps retrieves job steps for a job
	GetJobSteps(ctx context.Context, jobID string) (*JobStepList, error)
}

// NodeAPI provides access to node-related operations
type NodeAPI interface {
	// ListNodes lists compute nodes
	ListNodes(ctx context.Context, opts *ListNodesOptions) (*NodeList, error)
	
	// GetNode retrieves a specific node by name
	GetNode(ctx context.Context, nodeName string) (*Node, error)
	
	// UpdateNode updates node properties
	UpdateNode(ctx context.Context, nodeName string, update *NodeUpdate) error
}

// PartitionAPI provides access to partition-related operations
type PartitionAPI interface {
	// ListPartitions lists partitions
	ListPartitions(ctx context.Context) (*PartitionList, error)
	
	// GetPartition retrieves a specific partition by name
	GetPartition(ctx context.Context, partitionName string) (*Partition, error)
}

// Jobs returns the JobAPI interface
func (c *Client) Jobs() JobAPI {
	return &jobAPI{client: c}
}

// Nodes returns the NodeAPI interface
func (c *Client) Nodes() NodeAPI {
	return &nodeAPI{client: c}
}

// Partitions returns the PartitionAPI interface
func (c *Client) Partitions() PartitionAPI {
	return &partitionAPI{client: c}
}

// jobAPI implements the JobAPI interface
type jobAPI struct {
	client *Client
}

// nodeAPI implements the NodeAPI interface
type nodeAPI struct {
	client *Client
}

// partitionAPI implements the PartitionAPI interface
type partitionAPI struct {
	client *Client
}

// SlurmError represents an error from the Slurm REST API
type SlurmError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Source  string `json:"source"`
}

func (e *SlurmError) Error() string {
	return e.Message
}

// Common data structures

// Job represents a Slurm job
type Job struct {
	ID          string            `json:"job_id"`
	Name        string            `json:"name"`
	UserID      string            `json:"user_id"`
	GroupID     string            `json:"group_id"`
	State       string            `json:"job_state"`
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
}

// JobList represents a list of jobs
type JobList struct {
	Jobs []Job `json:"jobs"`
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
}

// JobSubmitResponse represents the response from job submission
type JobSubmitResponse struct {
	JobID string `json:"job_id"`
}

// ListJobsOptions represents options for listing jobs
type ListJobsOptions struct {
	UserID    string
	State     string
	Partition string
	Limit     int
	Offset    int
}

// Node represents a compute node
type Node struct {
	Name        string            `json:"name"`
	State       string            `json:"state"`
	CPUs        int               `json:"cpus"`
	Memory      int               `json:"memory"`
	Features    []string          `json:"features"`
	Partitions  []string          `json:"partitions"`
	Architecture string           `json:"architecture"`
	OS          string            `json:"os"`
	Reason      string            `json:"reason,omitempty"`
	LastBusy    *time.Time        `json:"last_busy,omitempty"`
}

// NodeList represents a list of nodes
type NodeList struct {
	Nodes []Node `json:"nodes"`
}

// NodeUpdate represents node update parameters
type NodeUpdate struct {
	State  string `json:"state,omitempty"`
	Reason string `json:"reason,omitempty"`
}

// ListNodesOptions represents options for listing nodes
type ListNodesOptions struct {
	State     string
	Partition string
	Features  []string
}

// Partition represents a job partition
type Partition struct {
	Name            string   `json:"name"`
	State           string   `json:"state"`
	Nodes           []string `json:"nodes"`
	TotalCPUs       int      `json:"total_cpus"`
	TotalMemory     int      `json:"total_memory"`
	MaxTimeLimit    int      `json:"max_time_limit"`
	DefaultTimeLimit int     `json:"default_time_limit"`
	Priority        int      `json:"priority"`
}

// PartitionList represents a list of partitions
type PartitionList struct {
	Partitions []Partition `json:"partitions"`
}

// JobStep represents a job step
type JobStep struct {
	ID       string     `json:"step_id"`
	JobID    string     `json:"job_id"`
	Name     string     `json:"name"`
	State    string     `json:"state"`
	CPUs     int        `json:"cpus"`
	Memory   int        `json:"memory"`
	NodeList string     `json:"node_list"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
}

// JobStepList represents a list of job steps
type JobStepList struct {
	Steps []JobStep `json:"steps"`
}

// API implementations

// ListJobs lists jobs with optional filtering
func (j *jobAPI) ListJobs(ctx context.Context, opts *ListJobsOptions) (*JobList, error) {
	endpoint := "/slurm/v0.0.39/jobs"
	params := url.Values{}
	
	if opts != nil {
		if opts.UserID != "" {
			params.Set("user_id", opts.UserID)
		}
		if opts.State != "" {
			params.Set("state", opts.State)
		}
		if opts.Partition != "" {
			params.Set("partition", opts.Partition)
		}
		if opts.Limit > 0 {
			params.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Offset > 0 {
			params.Set("offset", strconv.Itoa(opts.Offset))
		}
	}
	
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}
	
	var jobList JobList
	if err := j.client.makeRequest(ctx, http.MethodGet, endpoint, nil, &jobList); err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}
	
	return &jobList, nil
}

// GetJob retrieves a specific job by ID
func (j *jobAPI) GetJob(ctx context.Context, jobID string) (*Job, error) {
	endpoint := fmt.Sprintf("/slurm/v0.0.39/job/%s", jobID)
	
	var job Job
	if err := j.client.makeRequest(ctx, http.MethodGet, endpoint, nil, &job); err != nil {
		return nil, fmt.Errorf("failed to get job %s: %w", jobID, err)
	}
	
	return &job, nil
}

// SubmitJob submits a new job
func (j *jobAPI) SubmitJob(ctx context.Context, job *JobSubmission) (*JobSubmitResponse, error) {
	endpoint := "/slurm/v0.0.39/job/submit"
	
	var response JobSubmitResponse
	if err := j.client.makeRequest(ctx, http.MethodPost, endpoint, job, &response); err != nil {
		return nil, fmt.Errorf("failed to submit job: %w", err)
	}
	
	return &response, nil
}

// CancelJob cancels a job
func (j *jobAPI) CancelJob(ctx context.Context, jobID string) error {
	endpoint := fmt.Sprintf("/slurm/v0.0.39/job/%s", jobID)
	
	if err := j.client.makeRequest(ctx, http.MethodDelete, endpoint, nil, nil); err != nil {
		return fmt.Errorf("failed to cancel job %s: %w", jobID, err)
	}
	
	return nil
}

// GetJobSteps retrieves job steps for a job
func (j *jobAPI) GetJobSteps(ctx context.Context, jobID string) (*JobStepList, error) {
	endpoint := fmt.Sprintf("/slurm/v0.0.39/job/%s/steps", jobID)
	
	var stepList JobStepList
	if err := j.client.makeRequest(ctx, http.MethodGet, endpoint, nil, &stepList); err != nil {
		return nil, fmt.Errorf("failed to get job steps for job %s: %w", jobID, err)
	}
	
	return &stepList, nil
}

// ListNodes lists compute nodes
func (n *nodeAPI) ListNodes(ctx context.Context, opts *ListNodesOptions) (*NodeList, error) {
	endpoint := "/slurm/v0.0.39/nodes"
	params := url.Values{}
	
	if opts != nil {
		if opts.State != "" {
			params.Set("state", opts.State)
		}
		if opts.Partition != "" {
			params.Set("partition", opts.Partition)
		}
		if len(opts.Features) > 0 {
			params.Set("features", strings.Join(opts.Features, ","))
		}
	}
	
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}
	
	var nodeList NodeList
	if err := n.client.makeRequest(ctx, http.MethodGet, endpoint, nil, &nodeList); err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}
	
	return &nodeList, nil
}

// GetNode retrieves a specific node by name
func (n *nodeAPI) GetNode(ctx context.Context, nodeName string) (*Node, error) {
	endpoint := fmt.Sprintf("/slurm/v0.0.39/node/%s", nodeName)
	
	var node Node
	if err := n.client.makeRequest(ctx, http.MethodGet, endpoint, nil, &node); err != nil {
		return nil, fmt.Errorf("failed to get node %s: %w", nodeName, err)
	}
	
	return &node, nil
}

// UpdateNode updates node properties
func (n *nodeAPI) UpdateNode(ctx context.Context, nodeName string, update *NodeUpdate) error {
	endpoint := fmt.Sprintf("/slurm/v0.0.39/node/%s", nodeName)
	
	if err := n.client.makeRequest(ctx, http.MethodPatch, endpoint, update, nil); err != nil {
		return fmt.Errorf("failed to update node %s: %w", nodeName, err)
	}
	
	return nil
}

// ListPartitions lists partitions
func (p *partitionAPI) ListPartitions(ctx context.Context) (*PartitionList, error) {
	endpoint := "/slurm/v0.0.39/partitions"
	
	var partitionList PartitionList
	if err := p.client.makeRequest(ctx, http.MethodGet, endpoint, nil, &partitionList); err != nil {
		return nil, fmt.Errorf("failed to list partitions: %w", err)
	}
	
	return &partitionList, nil
}

// GetPartition retrieves a specific partition by name
func (p *partitionAPI) GetPartition(ctx context.Context, partitionName string) (*Partition, error) {
	endpoint := fmt.Sprintf("/slurm/v0.0.39/partition/%s", partitionName)
	
	var partition Partition
	if err := p.client.makeRequest(ctx, http.MethodGet, endpoint, nil, &partition); err != nil {
		return nil, fmt.Errorf("failed to get partition %s: %w", partitionName, err)
	}
	
	return &partition, nil
}