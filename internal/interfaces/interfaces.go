// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

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

	// Accounts returns the AccountManager for this version (v0.0.43+)
	Accounts() AccountManager

	// Users returns the UserManager for this version (v0.0.43+)
	Users() UserManager

	// Clusters returns the ClusterManager for this version (v0.0.43+)
	Clusters() ClusterManager

	// Associations returns the AssociationManager for this version (v0.0.43+)
	Associations() AssociationManager

	// WCKeys returns the WCKeyManager for this version (v0.0.43+)
	WCKeys() WCKeyManager

	// === Standalone Operations ===

	// GetLicenses retrieves license information
	GetLicenses(ctx context.Context) (*LicenseList, error)

	// GetShares retrieves fairshare information with optional filtering
	GetShares(ctx context.Context, opts *GetSharesOptions) (*SharesList, error)

	// GetConfig retrieves SLURM configuration
	GetConfig(ctx context.Context) (*Config, error)

	// GetDiagnostics retrieves SLURM diagnostics information
	GetDiagnostics(ctx context.Context) (*Diagnostics, error)

	// GetDBDiagnostics retrieves SLURM database diagnostics information
	GetDBDiagnostics(ctx context.Context) (*Diagnostics, error)

	// GetInstance retrieves a specific database instance
	GetInstance(ctx context.Context, opts *GetInstanceOptions) (*Instance, error)

	// GetInstances retrieves multiple database instances with filtering
	GetInstances(ctx context.Context, opts *GetInstancesOptions) (*InstanceList, error)

	// GetTRES retrieves all TRES (Trackable RESources)
	GetTRES(ctx context.Context) (*TRESList, error)

	// CreateTRES creates a new TRES entry
	CreateTRES(ctx context.Context, req *CreateTRESRequest) (*TRES, error)

	// Reconfigure triggers a SLURM reconfiguration
	Reconfigure(ctx context.Context) (*ReconfigureResponse, error)

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

	// Allocate allocates resources for a job (v0.0.43+)
	Allocate(ctx context.Context, req *JobAllocateRequest) (*JobAllocateResponse, error)

	// Cancel cancels a job
	Cancel(ctx context.Context, jobID string) error

	// Requeue requeues a job, allowing it to run again
	Requeue(ctx context.Context, jobID string) error

	// Update updates job properties (if supported by version)
	Update(ctx context.Context, jobID string, update *JobUpdate) error

	// Steps retrieves job steps for a job
	Steps(ctx context.Context, jobID string) (*JobStepList, error)

	// Watch provides real-time job updates (if supported by version)
	Watch(ctx context.Context, opts *WatchJobsOptions) (<-chan JobEvent, error)

	// Hold holds a job (prevents it from running)
	Hold(ctx context.Context, jobID string) error

	// Release releases a held job (allows it to run)
	Release(ctx context.Context, jobID string) error

	// Signal sends a signal to a job
	Signal(ctx context.Context, jobID string, signal string) error

	// Notify sends a message to a job
	Notify(ctx context.Context, jobID string, message string) error

	// Analytics methods for resource utilization and performance
	// GetJobUtilization retrieves comprehensive resource utilization metrics for a job
	GetJobUtilization(ctx context.Context, jobID string) (*JobUtilization, error)
	// GetJobEfficiency calculates efficiency metrics for a completed job
	GetJobEfficiency(ctx context.Context, jobID string) (*ResourceUtilization, error)
	// GetJobPerformance retrieves detailed performance metrics for a job
	GetJobPerformance(ctx context.Context, jobID string) (*JobPerformance, error)

	// GetJobLiveMetrics retrieves real-time performance metrics for a running job
	GetJobLiveMetrics(ctx context.Context, jobID string) (*JobLiveMetrics, error)

	// WatchJobMetrics provides streaming performance updates for a running job
	WatchJobMetrics(ctx context.Context, jobID string, opts *WatchMetricsOptions) (<-chan JobMetricsEvent, error)

	// GetJobResourceTrends retrieves performance trends over specified time windows
	GetJobResourceTrends(ctx context.Context, jobID string, opts *ResourceTrendsOptions) (*JobResourceTrends, error)

	// GetJobStepDetails retrieves detailed information about a specific job step
	GetJobStepDetails(ctx context.Context, jobID string, stepID string) (*JobStepDetails, error)

	// GetJobStepUtilization retrieves resource utilization metrics for a specific job step
	GetJobStepUtilization(ctx context.Context, jobID string, stepID string) (*JobStepUtilization, error)

	// ListJobStepsWithMetrics retrieves all job steps with their performance metrics
	ListJobStepsWithMetrics(ctx context.Context, jobID string, opts *ListJobStepsOptions) (*JobStepMetricsList, error)

	// SLURM Integration Methods for Task 2.7

	// GetJobStepsFromAccounting retrieves job step data from SLURM's accounting database
	GetJobStepsFromAccounting(ctx context.Context, jobID string, opts *AccountingQueryOptions) (*AccountingJobSteps, error)

	// GetStepAccountingData retrieves detailed accounting information for a specific job step
	GetStepAccountingData(ctx context.Context, jobID string, stepID string) (*StepAccountingRecord, error)

	// GetJobStepAPIData integrates with SLURM's native job step APIs for real-time data
	GetJobStepAPIData(ctx context.Context, jobID string, stepID string) (*JobStepAPIData, error)

	// ListJobStepsFromSacct queries job steps using SLURM's sacct command integration
	ListJobStepsFromSacct(ctx context.Context, jobID string, opts *SacctQueryOptions) (*SacctJobStepData, error)

	// Advanced Analytics Methods
	// GetJobCPUAnalytics retrieves detailed CPU performance metrics for a job
	GetJobCPUAnalytics(ctx context.Context, jobID string) (*CPUAnalytics, error)

	// GetJobMemoryAnalytics retrieves detailed memory performance metrics for a job
	GetJobMemoryAnalytics(ctx context.Context, jobID string) (*MemoryAnalytics, error)

	// GetJobIOAnalytics retrieves detailed I/O performance metrics for a job
	GetJobIOAnalytics(ctx context.Context, jobID string) (*IOAnalytics, error)

	// GetJobComprehensiveAnalytics retrieves all performance metrics for a job
	GetJobComprehensiveAnalytics(ctx context.Context, jobID string) (*JobComprehensiveAnalytics, error)

	// Historical Performance Tracking Methods
	// GetJobPerformanceHistory retrieves historical performance data for a job
	GetJobPerformanceHistory(ctx context.Context, jobID string, opts *PerformanceHistoryOptions) (*JobPerformanceHistory, error)

	// GetPerformanceTrends analyzes cluster-wide performance trends
	GetPerformanceTrends(ctx context.Context, opts *TrendAnalysisOptions) (*PerformanceTrends, error)

	// GetUserEfficiencyTrends tracks efficiency trends for a specific user
	GetUserEfficiencyTrends(ctx context.Context, userID string, opts *EfficiencyTrendOptions) (*UserEfficiencyTrends, error)

	// AnalyzeBatchJobs performs bulk analysis on a collection of jobs
	AnalyzeBatchJobs(ctx context.Context, jobIDs []string, opts *BatchAnalysisOptions) (*BatchJobAnalysis, error)

	// GetWorkflowPerformance analyzes performance of multi-job workflows
	GetWorkflowPerformance(ctx context.Context, workflowID string, opts *WorkflowAnalysisOptions) (*WorkflowPerformance, error)

	// GenerateEfficiencyReport creates comprehensive efficiency reports
	GenerateEfficiencyReport(ctx context.Context, opts *ReportOptions) (*EfficiencyReport, error)
}

// NodeManager provides version-agnostic node operations
type NodeManager interface {
	// List nodes with optional filtering
	List(ctx context.Context, opts *ListNodesOptions) (*NodeList, error)

	// Get retrieves a specific node by name
	Get(ctx context.Context, nodeName string) (*Node, error)

	// Update updates node properties (if supported by version)
	Update(ctx context.Context, nodeName string, update *NodeUpdate) error

	// Delete removes a node from the cluster (if supported by version)
	Delete(ctx context.Context, nodeName string) error

	// Drain drains a node, preventing new jobs from being scheduled on it
	Drain(ctx context.Context, nodeName string, reason string) error

	// Resume resumes a drained node, allowing new jobs to be scheduled on it
	Resume(ctx context.Context, nodeName string) error

	// Watch provides real-time node updates (if supported by version)
	Watch(ctx context.Context, opts *WatchNodesOptions) (<-chan NodeEvent, error)
}

// PartitionManager provides version-agnostic partition operations
type PartitionManager interface {
	// List partitions with optional filtering
	List(ctx context.Context, opts *ListPartitionsOptions) (*PartitionList, error)

	// Get retrieves a specific partition by name
	Get(ctx context.Context, partitionName string) (*Partition, error)

	// Create creates a new partition (if supported by version)
	Create(ctx context.Context, partition *PartitionCreate) (*PartitionCreateResponse, error)

	// Update updates partition properties (if supported by version)
	Update(ctx context.Context, partitionName string, update *PartitionUpdate) error

	// Delete removes a partition (if supported by version)
	Delete(ctx context.Context, partitionName string) error

	// Watch provides real-time partition updates (if supported by version)
	Watch(ctx context.Context, opts *WatchPartitionsOptions) (<-chan PartitionEvent, error)
}

// InfoManager provides version-agnostic cluster information operations
type InfoManager interface {
	// Get retrieves cluster information
	Get(ctx context.Context) (*ClusterInfo, error)

	// Ping tests connectivity to the cluster
	Ping(ctx context.Context) error

	// PingDatabase tests connectivity to the SLURM database (v0.0.43+)
	PingDatabase(ctx context.Context) error

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
	UserName    string                 `json:"user_name"`
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

// JobSubmitResponse represents the response from job submission
type JobSubmitResponse struct {
	JobID string `json:"job_id"`
}

// JobAllocateRequest represents a job allocation request (v0.0.43+)
type JobAllocateRequest struct {
	Name      string `json:"name"`
	Partition string `json:"partition,omitempty"`
	Nodes     int    `json:"nodes"`
	CPUs      int    `json:"cpus"`
	TimeLimit int    `json:"time_limit"` // in minutes
	Account   string `json:"account,omitempty"`
	QoS       string `json:"qos,omitempty"`
}

// JobAllocateResponse represents the response from job allocation (v0.0.43+)
type JobAllocateResponse struct {
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

// JobUtilization represents resource utilization metrics for a job
type JobUtilization struct {
	JobID              string                 `json:"job_id"`
	JobName            string                 `json:"job_name"`
	UserID             string                 `json:"user_id"`
	StartTime          time.Time              `json:"start_time"`
	EndTime            *time.Time             `json:"end_time,omitempty"`
	CPUUtilization     *ResourceUtilization   `json:"cpu_utilization,omitempty"`
	MemoryUtilization  *ResourceUtilization   `json:"memory_utilization,omitempty"`
	GPUUtilization     *GPUUtilization        `json:"gpu_utilization,omitempty"`
	IOUtilization      *IOUtilization         `json:"io_utilization,omitempty"`
	NetworkUtilization *NetworkUtilization    `json:"network_utilization,omitempty"`
	EnergyUsage        *EnergyUsage           `json:"energy_usage,omitempty"`
	SamplingInterval   int                    `json:"sampling_interval_seconds"`
	LastUpdated        time.Time              `json:"last_updated"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

// ResourceUtilization represents generic resource utilization metrics
type ResourceUtilization struct {
	Requested      float64   `json:"requested"`
	Allocated      float64   `json:"allocated"`
	Used           float64   `json:"used"`
	UsedMin        float64   `json:"used_min"`
	UsedMax        float64   `json:"used_max"`
	UsedAvg        float64   `json:"used_avg"`
	UsedStdDev     float64   `json:"used_stddev"`
	Efficiency     float64   `json:"efficiency_percentage"`
	Wasted         float64   `json:"wasted"`
	SampleCount    int       `json:"sample_count"`
	Unit           string    `json:"unit"`
	LastSampleTime time.Time `json:"last_sample_time"`

	// Additional fields for step-level analytics
	Limit      float64                `json:"limit,omitempty"`
	Percentage float64                `json:"percentage,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// GPUUtilization represents GPU-specific utilization metrics
type GPUUtilization struct {
	DeviceCount        int                    `json:"device_count"`
	Devices            []GPUDeviceUtilization `json:"devices"`
	OverallUtilization *ResourceUtilization   `json:"overall_utilization"`
	MemoryUtilization  *ResourceUtilization   `json:"memory_utilization"`
	PowerUsage         *ResourceUtilization   `json:"power_usage"`
	Temperature        map[string]float64     `json:"temperature"`
	ComputeMode        string                 `json:"compute_mode"`
	DriverVersion      string                 `json:"driver_version"`
	CUDAVersion        string                 `json:"cuda_version"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

// GPUDeviceUtilization represents per-GPU device utilization
type GPUDeviceUtilization struct {
	DeviceID          string               `json:"device_id"`
	DeviceName        string               `json:"device_name"`
	DeviceUUID        string               `json:"device_uuid"`
	Utilization       *ResourceUtilization `json:"utilization"`
	MemoryUtilization *ResourceUtilization `json:"memory_utilization"`
	PowerUsage        float64              `json:"power_usage_watts"`
	Temperature       float64              `json:"temperature_celsius"`
	PCIeBandwidth     *ResourceUtilization `json:"pcie_bandwidth,omitempty"`
	Processes         []GPUProcess         `json:"processes,omitempty"`
}

// GPUProcess represents a process running on a GPU
type GPUProcess struct {
	PID         int    `json:"pid"`
	ProcessName string `json:"process_name"`
	MemoryUsed  int64  `json:"memory_used_mb"`
}

// IOUtilization represents I/O utilization metrics
type IOUtilization struct {
	ReadBandwidth     *ResourceUtilization   `json:"read_bandwidth"`
	WriteBandwidth    *ResourceUtilization   `json:"write_bandwidth"`
	ReadIOPS          *ResourceUtilization   `json:"read_iops"`
	WriteIOPS         *ResourceUtilization   `json:"write_iops"`
	TotalBytesRead    int64                  `json:"total_bytes_read"`
	TotalBytesWritten int64                  `json:"total_bytes_written"`
	FileSystems       map[string]IOStats     `json:"file_systems,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// IOStats represents I/O statistics for a specific file system
type IOStats struct {
	MountPoint      string  `json:"mount_point"`
	BytesRead       int64   `json:"bytes_read"`
	BytesWritten    int64   `json:"bytes_written"`
	ReadsCompleted  int64   `json:"reads_completed"`
	WritesCompleted int64   `json:"writes_completed"`
	AvgReadLatency  float64 `json:"avg_read_latency_ms"`
	AvgWriteLatency float64 `json:"avg_write_latency_ms"`
}

// NetworkUtilization represents network utilization metrics
type NetworkUtilization struct {
	Interfaces       map[string]NetworkInterfaceStats `json:"interfaces"`
	TotalBandwidth   *ResourceUtilization             `json:"total_bandwidth"`
	IngressBandwidth *ResourceUtilization             `json:"ingress_bandwidth"`
	EgressBandwidth  *ResourceUtilization             `json:"egress_bandwidth"`
	PacketsReceived  int64                            `json:"packets_received"`
	PacketsSent      int64                            `json:"packets_sent"`
	PacketsDropped   int64                            `json:"packets_dropped"`
	Errors           int64                            `json:"errors"`
	ProtocolStats    map[string]int64                 `json:"protocol_stats,omitempty"`
	Metadata         map[string]interface{}           `json:"metadata,omitempty"`
}

// NetworkInterfaceStats represents statistics for a specific network interface
type NetworkInterfaceStats struct {
	InterfaceName   string  `json:"interface_name"`
	BytesReceived   int64   `json:"bytes_received"`
	BytesSent       int64   `json:"bytes_sent"`
	PacketsReceived int64   `json:"packets_received"`
	PacketsSent     int64   `json:"packets_sent"`
	BandwidthMbps   float64 `json:"bandwidth_mbps"`
	Utilization     float64 `json:"utilization_percentage"`
}

// EnergyUsage represents energy consumption metrics
type EnergyUsage struct {
	TotalEnergyJoules  float64                `json:"total_energy_joules"`
	AveragePowerWatts  float64                `json:"average_power_watts"`
	PeakPowerWatts     float64                `json:"peak_power_watts"`
	MinPowerWatts      float64                `json:"min_power_watts"`
	CPUEnergyJoules    float64                `json:"cpu_energy_joules"`
	GPUEnergyJoules    float64                `json:"gpu_energy_joules"`
	MemoryEnergyJoules float64                `json:"memory_energy_joules"`
	CarbonFootprint    float64                `json:"carbon_footprint_kg_co2"`
	PowerSources       map[string]float64     `json:"power_sources,omitempty"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

// JobPerformance represents comprehensive performance metrics for a job
type JobPerformance struct {
	JobID     uint32     `json:"job_id"`
	JobName   string     `json:"job_name"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Status    string     `json:"status"`
	ExitCode  int        `json:"exit_code"`

	// Resource utilization
	ResourceUtilization *ResourceUtilization `json:"resource_utilization"`
	JobUtilization      *JobUtilization      `json:"job_utilization"`

	// Step-level metrics
	StepMetrics []JobStepPerformance `json:"step_metrics,omitempty"`

	// Performance trends over time
	PerformanceTrends *PerformanceTrends `json:"performance_trends,omitempty"`

	// Bottleneck analysis
	Bottlenecks []PerformanceBottleneck `json:"bottlenecks,omitempty"`

	// Optimization recommendations
	Recommendations []OptimizationRecommendation `json:"recommendations,omitempty"`
}

// JobStepPerformance represents performance metrics for a job step
type JobStepPerformance struct {
	StepID    uint32        `json:"step_id"`
	StepName  string        `json:"step_name"`
	StartTime time.Time     `json:"start_time"`
	EndTime   *time.Time    `json:"end_time,omitempty"`
	Duration  time.Duration `json:"duration"`
	ExitCode  int           `json:"exit_code"`

	// Resource metrics
	CPUUtilization    float64 `json:"cpu_utilization"`
	MemoryUtilization float64 `json:"memory_utilization"`
	GPUUtilization    float64 `json:"gpu_utilization,omitempty"`
	IOThroughput      float64 `json:"io_throughput"`
	NetworkThroughput float64 `json:"network_throughput"`
}

// PerformanceTrends is defined later in the file with comprehensive historical tracking structure

// PerformanceBottleneck represents a detected performance bottleneck
type PerformanceBottleneck struct {
	Type          string        `json:"type"`     // cpu, memory, gpu, io, network
	Resource      string        `json:"resource"` // Specific resource affected
	Severity      string        `json:"severity"` // low, medium, high, critical
	Description   string        `json:"description"`
	Impact        string        `json:"impact"` // Description of performance impact
	TimeDetected  time.Time     `json:"time_detected"`
	Duration      time.Duration `json:"duration"`
	AffectedNodes []string      `json:"affected_nodes,omitempty"`
}

// OptimizationRecommendation represents a performance optimization suggestion
type OptimizationRecommendation struct {
	Type                string                 `json:"type"`     // resource_adjustment, configuration, workflow
	Priority            string                 `json:"priority"` // low, medium, high
	Title               string                 `json:"title"`
	Description         string                 `json:"description"`
	ExpectedImprovement float64                `json:"expected_improvement"` // Percentage improvement
	ResourceChanges     map[string]interface{} `json:"resource_changes,omitempty"`
	ConfigChanges       map[string]string      `json:"config_changes,omitempty"`
}

// JobLiveMetrics represents real-time performance metrics for a running job
type JobLiveMetrics struct {
	JobID          string        `json:"job_id"`
	JobName        string        `json:"job_name"`
	State          string        `json:"state"`
	RunningTime    time.Duration `json:"running_time"`
	CollectionTime time.Time     `json:"collection_time"`

	// Current resource usage
	CPUUsage     *LiveResourceMetric `json:"cpu_usage"`
	MemoryUsage  *LiveResourceMetric `json:"memory_usage"`
	GPUUsage     *LiveResourceMetric `json:"gpu_usage,omitempty"`
	NetworkUsage *LiveResourceMetric `json:"network_usage,omitempty"`
	IOUsage      *LiveResourceMetric `json:"io_usage,omitempty"`

	// Process information
	ProcessCount int `json:"process_count"`
	ThreadCount  int `json:"thread_count"`

	// Node-level metrics
	NodeMetrics map[string]*NodeLiveMetrics `json:"node_metrics,omitempty"`

	// Alerts and warnings
	Alerts []PerformanceAlert `json:"alerts,omitempty"`

	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// LiveResourceMetric represents a real-time resource metric
type LiveResourceMetric struct {
	Current            float64 `json:"current"`
	Average1Min        float64 `json:"average_1min"`
	Average5Min        float64 `json:"average_5min"`
	Peak               float64 `json:"peak"`
	Allocated          float64 `json:"allocated"`
	UtilizationPercent float64 `json:"utilization_percent"`
	Trend              string  `json:"trend"` // increasing, decreasing, stable
	Unit               string  `json:"unit"`
}

// NodeLiveMetrics represents real-time metrics for a specific node
type NodeLiveMetrics struct {
	NodeName string  `json:"node_name"`
	CPUCores int     `json:"cpu_cores"`
	MemoryGB float64 `json:"memory_gb"`

	// Resource metrics
	CPUUsage    *LiveResourceMetric `json:"cpu_usage"`
	MemoryUsage *LiveResourceMetric `json:"memory_usage"`
	LoadAverage []float64           `json:"load_average"` // 1, 5, 15 min

	// Temperature and power
	CPUTemperature   float64 `json:"cpu_temperature_celsius,omitempty"`
	PowerConsumption float64 `json:"power_consumption_watts,omitempty"`

	// Network and I/O
	NetworkInRate  float64 `json:"network_in_rate_mbps,omitempty"`
	NetworkOutRate float64 `json:"network_out_rate_mbps,omitempty"`
	DiskReadRate   float64 `json:"disk_read_rate_mbps,omitempty"`
	DiskWriteRate  float64 `json:"disk_write_rate_mbps,omitempty"`
}

// PerformanceAlert represents a performance-related alert or warning
type PerformanceAlert struct {
	Type              string    `json:"type"`     // warning, critical
	Category          string    `json:"category"` // cpu, memory, gpu, io, network
	Message           string    `json:"message"`
	Severity          string    `json:"severity"` // low, medium, high, critical
	Timestamp         time.Time `json:"timestamp"`
	NodeName          string    `json:"node_name,omitempty"`
	ResourceName      string    `json:"resource_name,omitempty"`
	CurrentValue      float64   `json:"current_value,omitempty"`
	ThresholdValue    float64   `json:"threshold_value,omitempty"`
	RecommendedAction string    `json:"recommended_action,omitempty"`
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
	CPULoad      float64                `json:"cpu_load,omitempty"`
	AllocCPUs    int32                  `json:"alloc_cpus,omitempty"`
	AllocMemory  int64                  `json:"alloc_memory,omitempty"`
	FreeMemory   int64                  `json:"free_memory,omitempty"`
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

// PartitionCreate represents a partition creation request
type PartitionCreate struct {
	Name          string   `json:"name"`
	State         string   `json:"state,omitempty"`
	Nodes         []string `json:"nodes,omitempty"`
	MaxTime       int      `json:"max_time,omitempty"`
	DefaultTime   int      `json:"default_time,omitempty"`
	MaxMemory     int      `json:"max_memory,omitempty"`
	DefaultMemory int      `json:"default_memory,omitempty"`
	AllowedUsers  []string `json:"allowed_users,omitempty"`
	DeniedUsers   []string `json:"denied_users,omitempty"`
	Priority      int      `json:"priority,omitempty"`
}

// PartitionCreateResponse represents the response from partition creation
type PartitionCreateResponse struct {
	PartitionName string `json:"partition_name"`
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
	JobsFailed   int `json:"jobs_failed"`
	JobsCanceled int `json:"jobs_canceled"`
	JobsTimeout  int `json:"jobs_timeout"`

	// Backfill scheduler statistics
	BackfillActive        bool  `json:"backfill_active"`
	BackfillJobsTotal     int   `json:"backfill_jobs_total"`
	BackfillJobsRecent    int   `json:"backfill_jobs_recent"`
	BackfillCycleCount    int   `json:"backfill_cycle_count"`
	BackfillCycleMeanTime int64 `json:"backfill_cycle_mean_time"`

	// Server performance statistics
	ServerThreadCount int   `json:"server_thread_count"`
	AgentQueueSize    int   `json:"agent_queue_size"`
	ScheduleCycleMax  int   `json:"schedule_cycle_max"`
	ScheduleCycleLast int   `json:"schedule_cycle_last"`
	ScheduleCycleMean int64 `json:"schedule_cycle_mean"`

	// RPC statistics
	RPCsTotal     int `json:"rpcs_total"`
	RPCsPending   int `json:"rpcs_pending"`
	RPCsCompleted int `json:"rpcs_completed"`

	// Additional metadata
	DiagTime time.Time              `json:"diag_time"`
	RawData  map[string]interface{} `json:"raw_data,omitempty"`
}

// Reservation represents a resource reservation
type Reservation struct {
	Name          string                 `json:"name"`
	State         string                 `json:"state"`
	StartTime     time.Time              `json:"start_time"`
	EndTime       time.Time              `json:"end_time"`
	Duration      int                    `json:"duration"`
	Nodes         []string               `json:"nodes"`
	NodeCount     int                    `json:"node_count"`
	CoreCount     int                    `json:"core_count"`
	Users         []string               `json:"users"`
	Accounts      []string               `json:"accounts"`
	Flags         []string               `json:"flags"`
	Features      []string               `json:"features"`
	PartitionName string                 `json:"partition_name"`
	Licenses      map[string]int         `json:"licenses"`
	BurstBuffer   string                 `json:"burst_buffer"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// ReservationList represents a list of reservations
type ReservationList struct {
	Reservations []Reservation `json:"reservations"`
	Total        int           `json:"total"`
}

// ReservationCreate represents a reservation creation request
type ReservationCreate struct {
	Name          string         `json:"name"`
	StartTime     time.Time      `json:"start_time"`
	EndTime       time.Time      `json:"end_time,omitempty"`
	Duration      int            `json:"duration,omitempty"`
	Nodes         []string       `json:"nodes,omitempty"`
	NodeCount     int            `json:"node_count,omitempty"`
	CoreCount     int            `json:"core_count,omitempty"`
	Users         []string       `json:"users,omitempty"`
	Accounts      []string       `json:"accounts,omitempty"`
	Flags         []string       `json:"flags,omitempty"`
	Features      []string       `json:"features,omitempty"`
	PartitionName string         `json:"partition_name,omitempty"`
	Licenses      map[string]int `json:"licenses,omitempty"`
	BurstBuffer   string         `json:"burst_buffer,omitempty"`
}

// ReservationCreateResponse represents the response from reservation creation
type ReservationCreateResponse struct {
	ReservationName string `json:"reservation_name"`
}

// ReservationUpdate represents a reservation update request
type ReservationUpdate struct {
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Duration  *int       `json:"duration,omitempty"`
	Nodes     []string   `json:"nodes,omitempty"`
	NodeCount *int       `json:"node_count,omitempty"`
	Users     []string   `json:"users,omitempty"`
	Accounts  []string   `json:"accounts,omitempty"`
	Flags     []string   `json:"flags,omitempty"`
	Features  []string   `json:"features,omitempty"`
}

// QoS represents a Quality of Service configuration
type QoS struct {
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	Priority          int                    `json:"priority"`
	PreemptMode       string                 `json:"preempt_mode"`
	GraceTime         int                    `json:"grace_time"`
	MaxJobs           int                    `json:"max_jobs"`
	MaxJobsPerUser    int                    `json:"max_jobs_per_user"`
	MaxJobsPerAccount int                    `json:"max_jobs_per_account"`
	MaxSubmitJobs     int                    `json:"max_submit_jobs"`
	MaxCPUs           int                    `json:"max_cpus"`
	MaxCPUsPerUser    int                    `json:"max_cpus_per_user"`
	MaxNodes          int                    `json:"max_nodes"`
	MaxWallTime       int                    `json:"max_wall_time"`
	MinCPUs           int                    `json:"min_cpus"`
	MinNodes          int                    `json:"min_nodes"`
	UsageFactor       float64                `json:"usage_factor"`
	UsageThreshold    float64                `json:"usage_threshold"`
	Flags             []string               `json:"flags"`
	AllowedAccounts   []string               `json:"allowed_accounts"`
	DeniedAccounts    []string               `json:"denied_accounts"`
	AllowedUsers      []string               `json:"allowed_users"`
	DeniedUsers       []string               `json:"denied_users"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// QoSList represents a list of QoS configurations
type QoSList struct {
	QoS   []QoS `json:"qos"`
	Total int   `json:"total"`
}

// QoSCreate represents a QoS creation request
type QoSCreate struct {
	Name              string   `json:"name"`
	Description       string   `json:"description,omitempty"`
	Priority          int      `json:"priority,omitempty"`
	PreemptMode       string   `json:"preempt_mode,omitempty"`
	GraceTime         int      `json:"grace_time,omitempty"`
	MaxJobs           int      `json:"max_jobs,omitempty"`
	MaxJobsPerUser    int      `json:"max_jobs_per_user,omitempty"`
	MaxJobsPerAccount int      `json:"max_jobs_per_account,omitempty"`
	MaxSubmitJobs     int      `json:"max_submit_jobs,omitempty"`
	MaxCPUs           int      `json:"max_cpus,omitempty"`
	MaxCPUsPerUser    int      `json:"max_cpus_per_user,omitempty"`
	MaxNodes          int      `json:"max_nodes,omitempty"`
	MaxWallTime       int      `json:"max_wall_time,omitempty"`
	MinCPUs           int      `json:"min_cpus,omitempty"`
	MinNodes          int      `json:"min_nodes,omitempty"`
	UsageFactor       float64  `json:"usage_factor,omitempty"`
	UsageThreshold    float64  `json:"usage_threshold,omitempty"`
	Flags             []string `json:"flags,omitempty"`
	AllowedAccounts   []string `json:"allowed_accounts,omitempty"`
	DeniedAccounts    []string `json:"denied_accounts,omitempty"`
	AllowedUsers      []string `json:"allowed_users,omitempty"`
	DeniedUsers       []string `json:"denied_users,omitempty"`
}

// QoSCreateResponse represents the response from QoS creation
type QoSCreateResponse struct {
	QoSName string `json:"qos_name"`
}

// QoSUpdate represents a QoS update request
type QoSUpdate struct {
	Description       *string  `json:"description,omitempty"`
	Priority          *int     `json:"priority,omitempty"`
	PreemptMode       *string  `json:"preempt_mode,omitempty"`
	GraceTime         *int     `json:"grace_time,omitempty"`
	MaxJobs           *int     `json:"max_jobs,omitempty"`
	MaxJobsPerUser    *int     `json:"max_jobs_per_user,omitempty"`
	MaxJobsPerAccount *int     `json:"max_jobs_per_account,omitempty"`
	MaxSubmitJobs     *int     `json:"max_submit_jobs,omitempty"`
	MaxCPUs           *int     `json:"max_cpus,omitempty"`
	MaxCPUsPerUser    *int     `json:"max_cpus_per_user,omitempty"`
	MaxNodes          *int     `json:"max_nodes,omitempty"`
	MaxWallTime       *int     `json:"max_wall_time,omitempty"`
	MinCPUs           *int     `json:"min_cpus,omitempty"`
	MinNodes          *int     `json:"min_nodes,omitempty"`
	UsageFactor       *float64 `json:"usage_factor,omitempty"`
	UsageThreshold    *float64 `json:"usage_threshold,omitempty"`
	Flags             []string `json:"flags,omitempty"`
	AllowedAccounts   []string `json:"allowed_accounts,omitempty"`
	DeniedAccounts    []string `json:"denied_accounts,omitempty"`
	AllowedUsers      []string `json:"allowed_users,omitempty"`
	DeniedUsers       []string `json:"denied_users,omitempty"`
}

// AccountManager manages account operations (v0.0.43+)
type AccountManager interface {
	List(ctx context.Context, opts *ListAccountsOptions) (*AccountList, error)
	Get(ctx context.Context, accountName string) (*Account, error)
	Create(ctx context.Context, account *AccountCreate) (*AccountCreateResponse, error)
	Update(ctx context.Context, accountName string, update *AccountUpdate) error
	Delete(ctx context.Context, accountName string) error

	// CreateAssociation creates a new user-account association
	CreateAssociation(ctx context.Context, userName, accountName string, opts *AssociationOptions) (*AssociationCreateResponse, error)

	// Enhanced account hierarchy methods
	GetAccountHierarchy(ctx context.Context, rootAccount string) (*AccountHierarchy, error)
	GetParentAccounts(ctx context.Context, accountName string) ([]*Account, error)
	GetChildAccounts(ctx context.Context, accountName string, depth int) ([]*Account, error)

	// Quota-related methods
	GetAccountQuotas(ctx context.Context, accountName string) (*AccountQuota, error)
	GetAccountQuotaUsage(ctx context.Context, accountName string, timeframe string) (*AccountUsage, error)

	// User-account association methods
	GetAccountUsers(ctx context.Context, accountName string, opts *ListAccountUsersOptions) ([]*UserAccountAssociation, error)
	ValidateUserAccess(ctx context.Context, userName, accountName string) (*UserAccessValidation, error)
	GetAccountUsersWithPermissions(ctx context.Context, accountName string, permissions []string) ([]*UserAccountAssociation, error)

	// Fair-share methods
	GetAccountFairShare(ctx context.Context, accountName string) (*AccountFairShare, error)
	GetFairShareHierarchy(ctx context.Context, rootAccount string) (*FairShareHierarchy, error)
}

// Account represents a SLURM account
type Account struct {
	Name              string         `json:"name"`
	Description       string         `json:"description,omitempty"`
	Organization      string         `json:"organization,omitempty"`
	CoordinatorUsers  []string       `json:"coordinator_users,omitempty"`
	AllowedPartitions []string       `json:"allowed_partitions,omitempty"`
	DefaultPartition  string         `json:"default_partition,omitempty"`
	AllowedQoS        []string       `json:"allowed_qos,omitempty"`
	DefaultQoS        string         `json:"default_qos,omitempty"`
	CPULimit          int            `json:"cpu_limit,omitempty"`
	MaxJobs           int            `json:"max_jobs,omitempty"`
	MaxJobsPerUser    int            `json:"max_jobs_per_user,omitempty"`
	MaxNodes          int            `json:"max_nodes,omitempty"`
	MaxWallTime       int            `json:"max_wall_time,omitempty"`
	FairShareTRES     map[string]int `json:"fairshare_tres,omitempty"`
	GrpTRES           map[string]int `json:"grp_tres,omitempty"`
	GrpTRESMinutes    map[string]int `json:"grp_tres_minutes,omitempty"`
	MaxTRES           map[string]int `json:"max_tres,omitempty"`
	MaxTRESPerUser    map[string]int `json:"max_tres_per_user,omitempty"`
	SharesPriority    int            `json:"shares_priority,omitempty"`
	ParentAccount     string         `json:"parent_account,omitempty"`
	ChildAccounts     []string       `json:"child_accounts,omitempty"`
	Users             []string       `json:"users,omitempty"`
	Flags             []string       `json:"flags,omitempty"`
	CreateTime        string         `json:"create_time,omitempty"`
	UpdateTime        string         `json:"update_time,omitempty"`

	// Enhanced fields for extended account management
	Quota            *AccountQuota `json:"quota,omitempty"`
	Usage            *AccountUsage `json:"usage,omitempty"`
	HierarchyLevel   int           `json:"hierarchy_level,omitempty"`
	HierarchyPath    []string      `json:"hierarchy_path,omitempty"`
	TotalSubAccounts int           `json:"total_sub_accounts,omitempty"`
	ActiveUserCount  int           `json:"active_user_count,omitempty"`
}

// AccountList represents a list of accounts
type AccountList struct {
	Accounts []Account `json:"accounts"`
	Total    int       `json:"total"`
}

// AccountCreate represents a request to create an account
type AccountCreate struct {
	Name              string         `json:"name"`
	Description       string         `json:"description,omitempty"`
	Organization      string         `json:"organization,omitempty"`
	CoordinatorUsers  []string       `json:"coordinator_users,omitempty"`
	AllowedPartitions []string       `json:"allowed_partitions,omitempty"`
	DefaultPartition  string         `json:"default_partition,omitempty"`
	AllowedQoS        []string       `json:"allowed_qos,omitempty"`
	DefaultQoS        string         `json:"default_qos,omitempty"`
	CPULimit          int            `json:"cpu_limit,omitempty"`
	MaxJobs           int            `json:"max_jobs,omitempty"`
	MaxJobsPerUser    int            `json:"max_jobs_per_user,omitempty"`
	MaxNodes          int            `json:"max_nodes,omitempty"`
	MaxWallTime       int            `json:"max_wall_time,omitempty"`
	FairShareTRES     map[string]int `json:"fairshare_tres,omitempty"`
	GrpTRES           map[string]int `json:"grp_tres,omitempty"`
	GrpTRESMinutes    map[string]int `json:"grp_tres_minutes,omitempty"`
	MaxTRES           map[string]int `json:"max_tres,omitempty"`
	MaxTRESPerUser    map[string]int `json:"max_tres_per_user,omitempty"`
	SharesPriority    int            `json:"shares_priority,omitempty"`
	ParentAccount     string         `json:"parent_account,omitempty"`
	Flags             []string       `json:"flags,omitempty"`
}

// AccountCreateResponse represents the response from account creation
type AccountCreateResponse struct {
	AccountName string `json:"account_name"`
}

// AccountUpdate represents an account update request
type AccountUpdate struct {
	Description       *string        `json:"description,omitempty"`
	Organization      *string        `json:"organization,omitempty"`
	CoordinatorUsers  []string       `json:"coordinator_users,omitempty"`
	AllowedPartitions []string       `json:"allowed_partitions,omitempty"`
	DefaultPartition  *string        `json:"default_partition,omitempty"`
	AllowedQoS        []string       `json:"allowed_qos,omitempty"`
	DefaultQoS        *string        `json:"default_qos,omitempty"`
	CPULimit          *int           `json:"cpu_limit,omitempty"`
	MaxJobs           *int           `json:"max_jobs,omitempty"`
	MaxJobsPerUser    *int           `json:"max_jobs_per_user,omitempty"`
	MaxNodes          *int           `json:"max_nodes,omitempty"`
	MaxWallTime       *int           `json:"max_wall_time,omitempty"`
	FairShareTRES     map[string]int `json:"fairshare_tres,omitempty"`
	GrpTRES           map[string]int `json:"grp_tres,omitempty"`
	GrpTRESMinutes    map[string]int `json:"grp_tres_minutes,omitempty"`
	MaxTRES           map[string]int `json:"max_tres,omitempty"`
	MaxTRESPerUser    map[string]int `json:"max_tres_per_user,omitempty"`
	SharesPriority    *int           `json:"shares_priority,omitempty"`
	Flags             []string       `json:"flags,omitempty"`
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

// ListAccountsOptions provides options for listing accounts
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

// AccountQuota represents quota information for an account
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

// AccountUsage represents usage statistics for an account
type AccountUsage struct {
	AccountName     string             `json:"account_name"`
	CPUHours        float64            `json:"cpu_hours"`
	JobCount        int                `json:"job_count"`
	JobsCompleted   int                `json:"jobs_completed"`
	JobsFailed      int                `json:"jobs_failed"`
	JobsCanceled    int                `json:"jobs_canceled"`
	TotalWallTime   float64            `json:"total_wall_time"`
	AverageWallTime float64            `json:"average_wall_time"`
	TRESUsage       map[string]float64 `json:"tres_usage,omitempty"`
	UserCount       int                `json:"user_count"`
	ActiveUsers     []string           `json:"active_users,omitempty"`
	Period          string             `json:"period,omitempty"`
	StartTime       time.Time          `json:"start_time"`
	EndTime         time.Time          `json:"end_time"`
	EfficiencyRatio float64            `json:"efficiency_ratio,omitempty"`
}

// AccountHierarchy represents the hierarchical structure of accounts
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

// Cluster represents a SLURM cluster configuration
type Cluster struct {
	Name               string                 `json:"name"`
	ControlHost        string                 `json:"control_host,omitempty"`
	ControlPort        int                    `json:"control_port,omitempty"`
	RPCVersion         int                    `json:"rpc_version,omitempty"`
	PluginIDSelect     int                    `json:"plugin_id_select,omitempty"`
	PluginIDAuth       int                    `json:"plugin_id_auth,omitempty"`
	PluginIDAcct       int                    `json:"plugin_id_acct,omitempty"`
	TRESList           []string               `json:"tres_list,omitempty"`
	Features           []string               `json:"features,omitempty"`
	FederationFeatures []string               `json:"federation_features,omitempty"`
	FederationState    string                 `json:"federation_state,omitempty"`
	Created            time.Time              `json:"created"`
	Modified           time.Time              `json:"modified"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

// ClusterList represents a list of clusters
type ClusterList struct {
	Clusters []*Cluster             `json:"clusters"`
	Total    int                    `json:"total"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ClusterCreate represents a request to create a cluster
type ClusterCreate struct {
	Name               string                 `json:"name"`
	ControlHost        string                 `json:"control_host,omitempty"`
	ControlPort        int                    `json:"control_port,omitempty"`
	RPCVersion         int                    `json:"rpc_version,omitempty"`
	PluginIDSelect     int                    `json:"plugin_id_select,omitempty"`
	PluginIDAuth       int                    `json:"plugin_id_auth,omitempty"`
	PluginIDAcct       int                    `json:"plugin_id_acct,omitempty"`
	TRESList           []string               `json:"tres_list,omitempty"`
	Features           []string               `json:"features,omitempty"`
	FederationFeatures []string               `json:"federation_features,omitempty"`
	FederationState    string                 `json:"federation_state,omitempty"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

// ClusterCreateResponse represents the response from cluster creation
type ClusterCreateResponse struct {
	Name     string                 `json:"name"`
	Created  time.Time              `json:"created"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ClusterUpdate represents a cluster update request
type ClusterUpdate struct {
	ControlHost        *string                `json:"control_host,omitempty"`
	ControlPort        *int                   `json:"control_port,omitempty"`
	RPCVersion         *int                   `json:"rpc_version,omitempty"`
	PluginIDSelect     *int                   `json:"plugin_id_select,omitempty"`
	PluginIDAuth       *int                   `json:"plugin_id_auth,omitempty"`
	PluginIDAcct       *int                   `json:"plugin_id_acct,omitempty"`
	TRESList           []string               `json:"tres_list,omitempty"`
	Features           []string               `json:"features,omitempty"`
	FederationFeatures []string               `json:"federation_features,omitempty"`
	FederationState    *string                `json:"federation_state,omitempty"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

// ListClustersOptions provides options for listing clusters
type ListClustersOptions struct {
	// Cluster name filtering
	Names []string `json:"names,omitempty"`
	// Federation state filtering
	FederationStates []string `json:"federation_states,omitempty"`
	// Feature filtering
	Features []string `json:"features,omitempty"`
	// Control host filtering
	ControlHosts []string `json:"control_hosts,omitempty"`
	// Include federation information
	WithFederation bool `json:"with_federation,omitempty"`
	// Include TRES information
	WithTRES bool `json:"with_tres,omitempty"`
	// Include detailed plugin information
	WithPlugins bool `json:"with_plugins,omitempty"`
	// Pagination
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
}

// Association represents a SLURM association (user-account-cluster relationship)
type Association struct {
	ID            uint32 `json:"id"`
	User          string `json:"user"`
	Account       string `json:"account"`
	Cluster       string `json:"cluster"`
	Partition     string `json:"partition,omitempty"`
	ParentAccount string `json:"parent_account,omitempty"`
	IsDefault     bool   `json:"is_default"`
	Comment       string `json:"comment,omitempty"`
	// Resource limits
	SharesRaw       int    `json:"shares_raw,omitempty"`
	Priority        uint32 `json:"priority,omitempty"`
	MaxJobs         *int   `json:"max_jobs,omitempty"`
	MaxJobsAccrue   *int   `json:"max_jobs_accrue,omitempty"`
	MaxSubmitJobs   *int   `json:"max_submit_jobs,omitempty"`
	MaxWallDuration *int   `json:"max_wall_duration_per_job,omitempty"`
	GrpJobs         *int   `json:"grp_jobs,omitempty"`
	GrpJobsAccrue   *int   `json:"grp_jobs_accrue,omitempty"`
	GrpSubmitJobs   *int   `json:"grp_submit_jobs,omitempty"`
	GrpWall         *int   `json:"grp_wall,omitempty"`
	// TRES limits (Trackable RESources)
	MaxTRESPerJob  map[string]string `json:"max_tres_per_job,omitempty"`
	MaxTRESMins    map[string]string `json:"max_tres_mins,omitempty"`
	GrpTRES        map[string]string `json:"grp_tres,omitempty"`
	GrpTRESMins    map[string]string `json:"grp_tres_mins,omitempty"`
	GrpTRESRunMins map[string]string `json:"grp_tres_run_mins,omitempty"`
	// QoS
	DefaultQoS string   `json:"default_qos,omitempty"`
	QoSList    []string `json:"qos_list,omitempty"`
	// Flags
	Flags []string `json:"flags,omitempty"`
	// Usage information
	FairShare      float64 `json:"fair_share,omitempty"`
	UsageRaw       int64   `json:"usage_raw,omitempty"`
	EffectiveUsage float64 `json:"effective_usage,omitempty"`
	// Timestamps
	Created  time.Time  `json:"created"`
	Modified time.Time  `json:"modified"`
	Deleted  *time.Time `json:"deleted,omitempty"`
	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// AssociationList represents a list of associations
type AssociationList struct {
	Associations []*Association         `json:"associations"`
	Total        int                    `json:"total"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// AssociationCreate represents a request to create an association
type AssociationCreate struct {
	User          string `json:"user"`
	Account       string `json:"account"`
	Cluster       string `json:"cluster,omitempty"`
	Partition     string `json:"partition,omitempty"`
	ParentAccount string `json:"parent_account,omitempty"`
	IsDefault     bool   `json:"is_default,omitempty"`
	Comment       string `json:"comment,omitempty"`
	// Resource limits
	SharesRaw       *int    `json:"shares_raw,omitempty"`
	Priority        *uint32 `json:"priority,omitempty"`
	MaxJobs         *int    `json:"max_jobs,omitempty"`
	MaxJobsAccrue   *int    `json:"max_jobs_accrue,omitempty"`
	MaxSubmitJobs   *int    `json:"max_submit_jobs,omitempty"`
	MaxWallDuration *int    `json:"max_wall_duration_per_job,omitempty"`
	GrpJobs         *int    `json:"grp_jobs,omitempty"`
	GrpJobsAccrue   *int    `json:"grp_jobs_accrue,omitempty"`
	GrpSubmitJobs   *int    `json:"grp_submit_jobs,omitempty"`
	GrpWall         *int    `json:"grp_wall,omitempty"`
	// TRES limits
	MaxTRESPerJob  map[string]string `json:"max_tres_per_job,omitempty"`
	MaxTRESMins    map[string]string `json:"max_tres_mins,omitempty"`
	GrpTRES        map[string]string `json:"grp_tres,omitempty"`
	GrpTRESMins    map[string]string `json:"grp_tres_mins,omitempty"`
	GrpTRESRunMins map[string]string `json:"grp_tres_run_mins,omitempty"`
	// QoS
	DefaultQoS string   `json:"default_qos,omitempty"`
	QoSList    []string `json:"qos_list,omitempty"`
	// Flags
	Flags []string `json:"flags,omitempty"`
	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// AssociationCreateResponse represents the response from association creation
type AssociationCreateResponse struct {
	Associations []*Association         `json:"associations"`
	Created      int                    `json:"created"`
	Updated      int                    `json:"updated"`
	Errors       []string               `json:"errors,omitempty"`
	Warnings     []string               `json:"warnings,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// AssociationUpdate represents an association update request
type AssociationUpdate struct {
	User      string `json:"user"`
	Account   string `json:"account"`
	Cluster   string `json:"cluster,omitempty"`
	Partition string `json:"partition,omitempty"`
	// Fields that can be updated
	IsDefault       *bool                  `json:"is_default,omitempty"`
	Comment         *string                `json:"comment,omitempty"`
	SharesRaw       *int                   `json:"shares_raw,omitempty"`
	Priority        *uint32                `json:"priority,omitempty"`
	MaxJobs         *int                   `json:"max_jobs,omitempty"`
	MaxJobsAccrue   *int                   `json:"max_jobs_accrue,omitempty"`
	MaxSubmitJobs   *int                   `json:"max_submit_jobs,omitempty"`
	MaxWallDuration *int                   `json:"max_wall_duration_per_job,omitempty"`
	GrpJobs         *int                   `json:"grp_jobs,omitempty"`
	GrpJobsAccrue   *int                   `json:"grp_jobs_accrue,omitempty"`
	GrpSubmitJobs   *int                   `json:"grp_submit_jobs,omitempty"`
	GrpWall         *int                   `json:"grp_wall,omitempty"`
	MaxTRESPerJob   map[string]string      `json:"max_tres_per_job,omitempty"`
	MaxTRESMins     map[string]string      `json:"max_tres_mins,omitempty"`
	GrpTRES         map[string]string      `json:"grp_tres,omitempty"`
	GrpTRESMins     map[string]string      `json:"grp_tres_mins,omitempty"`
	GrpTRESRunMins  map[string]string      `json:"grp_tres_run_mins,omitempty"`
	DefaultQoS      *string                `json:"default_qos,omitempty"`
	QoSList         []string               `json:"qos_list,omitempty"`
	Flags           []string               `json:"flags,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// ListAssociationsOptions provides filtering options for listing associations
type ListAssociationsOptions struct {
	// Filter by user
	Users []string `json:"users,omitempty"`
	// Filter by account
	Accounts []string `json:"accounts,omitempty"`
	// Filter by cluster
	Clusters []string `json:"clusters,omitempty"`
	// Filter by partition
	Partitions []string `json:"partitions,omitempty"`
	// Filter by parent account
	ParentAccounts []string `json:"parent_accounts,omitempty"`
	// Filter by QoS
	QoS []string `json:"qos,omitempty"`
	// Include deleted associations
	WithDeleted bool `json:"with_deleted,omitempty"`
	// Include usage information
	WithUsage bool `json:"with_usage,omitempty"`
	// Include TRES information
	WithTRES bool `json:"with_tres,omitempty"`
	// Include subaccounts
	WithSubAccounts bool `json:"with_sub_accounts,omitempty"`
	// Only default associations
	OnlyDefaults bool `json:"only_defaults,omitempty"`
	// Pagination
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
}

// GetAssociationOptions provides options for retrieving a specific association
type GetAssociationOptions struct {
	User      string `json:"user"`
	Account   string `json:"account"`
	Cluster   string `json:"cluster,omitempty"`
	Partition string `json:"partition,omitempty"`
	// Include usage information
	WithUsage bool `json:"with_usage,omitempty"`
	// Include TRES information
	WithTRES bool `json:"with_tres,omitempty"`
}

// DeleteAssociationOptions provides options for deleting an association
type DeleteAssociationOptions struct {
	User      string `json:"user"`
	Account   string `json:"account"`
	Cluster   string `json:"cluster,omitempty"`
	Partition string `json:"partition,omitempty"`
	// Force deletion even if jobs are running
	Force bool `json:"force,omitempty"`
}

// BulkDeleteOptions provides options for bulk deleting associations
type BulkDeleteOptions struct {
	// Filter criteria for associations to delete
	Users      []string `json:"users,omitempty"`
	Accounts   []string `json:"accounts,omitempty"`
	Clusters   []string `json:"clusters,omitempty"`
	Partitions []string `json:"partitions,omitempty"`
	// Delete only if no jobs are running
	OnlyIfIdle bool `json:"only_if_idle,omitempty"`
	// Force deletion
	Force bool `json:"force,omitempty"`
}

// BulkDeleteResponse represents the response from bulk delete operation
type BulkDeleteResponse struct {
	Deleted             int            `json:"deleted"`
	Failed              int            `json:"failed"`
	Errors              []string       `json:"errors,omitempty"`
	DeletedAssociations []*Association `json:"deleted_associations,omitempty"`
}

// Watch options for real-time updates

// WatchJobsOptions provides options for watching job changes
type WatchJobsOptions struct {
	UserID           string   `json:"user_id,omitempty"`
	States           []string `json:"states,omitempty"`
	Partition        string   `json:"partition,omitempty"`
	JobIDs           []string `json:"job_ids,omitempty"`
	ExcludeNew       bool     `json:"exclude_new,omitempty"`
	ExcludeCompleted bool     `json:"exclude_completed,omitempty"`
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

// WatchMetricsOptions provides options for watching job metrics
type WatchMetricsOptions struct {
	// Interval between metric collections (default 5s)
	UpdateInterval time.Duration `json:"update_interval,omitempty"`
	// Which metrics to collect
	IncludeCPU     bool `json:"include_cpu"`
	IncludeMemory  bool `json:"include_memory"`
	IncludeGPU     bool `json:"include_gpu"`
	IncludeNetwork bool `json:"include_network"`
	IncludeIO      bool `json:"include_io"`
	IncludeEnergy  bool `json:"include_energy"`
	// Per-node metrics
	IncludeNodeMetrics bool     `json:"include_node_metrics"`
	SpecificNodes      []string `json:"specific_nodes,omitempty"`
	// Alert thresholds
	CPUThreshold    float64 `json:"cpu_threshold,omitempty"`    // Alert if CPU usage > threshold (0-100)
	MemoryThreshold float64 `json:"memory_threshold,omitempty"` // Alert if memory usage > threshold (0-100)
	GPUThreshold    float64 `json:"gpu_threshold,omitempty"`    // Alert if GPU usage > threshold (0-100)
	// Stop conditions
	StopOnCompletion bool          `json:"stop_on_completion"`     // Stop watching when job completes
	MaxDuration      time.Duration `json:"max_duration,omitempty"` // Maximum time to watch
}

// JobMetricsEvent represents a job metrics update event
type JobMetricsEvent struct {
	Type        string            `json:"type"` // "update", "alert", "error", "complete"
	JobID       string            `json:"job_id"`
	Timestamp   time.Time         `json:"timestamp"`
	Metrics     *JobLiveMetrics   `json:"metrics,omitempty"`
	Alert       *PerformanceAlert `json:"alert,omitempty"`
	Error       error             `json:"error,omitempty"`
	StateChange *JobStateChange   `json:"state_change,omitempty"`
}

// JobStateChange represents a job state transition in metrics events
type JobStateChange struct {
	OldState string `json:"old_state"`
	NewState string `json:"new_state"`
	Reason   string `json:"reason,omitempty"`
}

// ResourceTrendsOptions provides options for retrieving resource trends
type ResourceTrendsOptions struct {
	// Time window for trend analysis
	TimeWindow time.Duration `json:"time_window,omitempty"` // Default: job duration or 1 hour
	// Number of data points to collect
	DataPoints int `json:"data_points,omitempty"` // Default: 24
	// Resources to include in trends
	IncludeCPU     bool `json:"include_cpu"`
	IncludeMemory  bool `json:"include_memory"`
	IncludeGPU     bool `json:"include_gpu"`
	IncludeIO      bool `json:"include_io"`
	IncludeNetwork bool `json:"include_network"`
	IncludeEnergy  bool `json:"include_energy"`
	// Aggregation method for data points
	Aggregation string `json:"aggregation,omitempty"` // "avg", "max", "min" (default: "avg")
	// Include anomaly detection
	DetectAnomalies bool `json:"detect_anomalies"`
}

// JobResourceTrends represents resource usage trends over time
type JobResourceTrends struct {
	JobID      string        `json:"job_id"`
	JobName    string        `json:"job_name"`
	StartTime  time.Time     `json:"start_time"`
	EndTime    *time.Time    `json:"end_time,omitempty"`
	TimeWindow time.Duration `json:"time_window"`
	DataPoints int           `json:"data_points"`

	// Time series data
	TimePoints []time.Time `json:"time_points"`

	// Resource trends
	CPUTrends     *ResourceTimeSeries `json:"cpu_trends,omitempty"`
	MemoryTrends  *ResourceTimeSeries `json:"memory_trends,omitempty"`
	GPUTrends     *ResourceTimeSeries `json:"gpu_trends,omitempty"`
	IOTrends      *IOTimeSeries       `json:"io_trends,omitempty"`
	NetworkTrends *NetworkTimeSeries  `json:"network_trends,omitempty"`
	EnergyTrends  *EnergyTimeSeries   `json:"energy_trends,omitempty"`

	// Anomalies detected
	Anomalies []ResourceAnomaly `json:"anomalies,omitempty"`

	// Summary statistics
	Summary *TrendsSummary `json:"summary"`

	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ResourceTimeSeries represents time series data for a resource
type ResourceTimeSeries struct {
	Values     []float64 `json:"values"`
	Unit       string    `json:"unit"`
	Average    float64   `json:"average"`
	Min        float64   `json:"min"`
	Max        float64   `json:"max"`
	StdDev     float64   `json:"std_dev"`
	Trend      string    `json:"trend"` // "increasing", "decreasing", "stable", "fluctuating"
	TrendSlope float64   `json:"trend_slope"`
}

// IOTimeSeries represents I/O time series data
type IOTimeSeries struct {
	ReadBandwidth  *ResourceTimeSeries `json:"read_bandwidth,omitempty"`
	WriteBandwidth *ResourceTimeSeries `json:"write_bandwidth,omitempty"`
	ReadIOPS       *ResourceTimeSeries `json:"read_iops,omitempty"`
	WriteIOPS      *ResourceTimeSeries `json:"write_iops,omitempty"`
}

// NetworkTimeSeries represents network time series data
type NetworkTimeSeries struct {
	IngressBandwidth *ResourceTimeSeries `json:"ingress_bandwidth,omitempty"`
	EgressBandwidth  *ResourceTimeSeries `json:"egress_bandwidth,omitempty"`
	PacketRate       *ResourceTimeSeries `json:"packet_rate,omitempty"`
}

// EnergyTimeSeries represents energy usage time series data
type EnergyTimeSeries struct {
	PowerUsage        *ResourceTimeSeries `json:"power_usage,omitempty"`
	EnergyConsumption *ResourceTimeSeries `json:"energy_consumption,omitempty"`
	CarbonEmissions   *ResourceTimeSeries `json:"carbon_emissions,omitempty"`
}

// ResourceAnomaly represents an anomaly detected in resource usage
type ResourceAnomaly struct {
	Timestamp     time.Time `json:"timestamp"`
	Resource      string    `json:"resource"` // "cpu", "memory", "gpu", etc.
	Type          string    `json:"type"`     // "spike", "drop", "pattern_change"
	Severity      string    `json:"severity"` // "low", "medium", "high"
	Value         float64   `json:"value"`
	ExpectedValue float64   `json:"expected_value"`
	Deviation     float64   `json:"deviation_percent"`
	Description   string    `json:"description"`
}

// TrendsSummary provides summary statistics for resource trends
type TrendsSummary struct {
	OverallTrend       string             `json:"overall_trend"`
	ResourceEfficiency float64            `json:"resource_efficiency"`
	StabilityScore     float64            `json:"stability_score"` // 0-100, higher is more stable
	VariabilityIndex   float64            `json:"variability_index"`
	PeakUtilization    map[string]float64 `json:"peak_utilization"`
	AverageUtilization map[string]float64 `json:"average_utilization"`
	ResourceBalance    string             `json:"resource_balance"` // "balanced", "cpu_heavy", "memory_heavy", etc.
}

// JobStepDetails represents detailed information about a job step
type JobStepDetails struct {
	StepID    string        `json:"step_id"`
	StepName  string        `json:"step_name"`
	JobID     string        `json:"job_id"`
	JobName   string        `json:"job_name"`
	State     string        `json:"state"`
	StartTime *time.Time    `json:"start_time,omitempty"`
	EndTime   *time.Time    `json:"end_time,omitempty"`
	Duration  time.Duration `json:"duration"`
	ExitCode  int           `json:"exit_code"`

	// Resource allocation
	CPUAllocation    int      `json:"cpu_allocation"`
	MemoryAllocation int64    `json:"memory_allocation_bytes"`
	GPUAllocation    int      `json:"gpu_allocation,omitempty"`
	NodeList         []string `json:"node_list"`
	TaskCount        int      `json:"task_count"`
	TasksPerNode     int      `json:"tasks_per_node,omitempty"`

	// Task distribution
	TaskDistribution map[string]int `json:"task_distribution,omitempty"` // node -> task count

	// Command and environment
	Command     string            `json:"command"`
	CommandLine string            `json:"command_line"`
	WorkingDir  string            `json:"working_dir"`
	Environment map[string]string `json:"environment,omitempty"`

	// Performance summary
	CPUTime    time.Duration `json:"cpu_time"`
	SystemTime time.Duration `json:"system_time"`
	UserTime   time.Duration `json:"user_time"`
	MaxRSS     int64         `json:"max_rss_bytes"`
	MaxVMSize  int64         `json:"max_vmsize_bytes"`
	AverageRSS int64         `json:"average_rss_bytes"`

	// I/O statistics
	TotalReadBytes  int64 `json:"total_read_bytes"`
	TotalWriteBytes int64 `json:"total_write_bytes"`
	ReadOperations  int64 `json:"read_operations"`
	WriteOperations int64 `json:"write_operations"`

	// Network statistics
	NetworkBytesReceived int64 `json:"network_bytes_received"`
	NetworkBytesSent     int64 `json:"network_bytes_sent"`

	// Energy usage (for newer API versions)
	EnergyConsumed   float64 `json:"energy_consumed_joules"`
	AveragePowerDraw float64 `json:"average_power_draw_watts"`

	// Task-level information
	Tasks []StepTaskInfo `json:"tasks,omitempty"`

	// Step-specific metadata
	StepType        string `json:"step_type"` // "primary", "batch", "interactive"
	Priority        int    `json:"priority"`
	AccountingGroup string `json:"accounting_group"`
	QOSLevel        string `json:"qos_level"`

	// Additional metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// JobStepUtilization represents resource utilization metrics for a job step
type JobStepUtilization struct {
	StepID    string        `json:"step_id"`
	StepName  string        `json:"step_name"`
	JobID     string        `json:"job_id"`
	JobName   string        `json:"job_name"`
	StartTime *time.Time    `json:"start_time,omitempty"`
	EndTime   *time.Time    `json:"end_time,omitempty"`
	Duration  time.Duration `json:"duration"`

	// Resource utilization
	CPUUtilization     *ResourceUtilization `json:"cpu_utilization,omitempty"`
	MemoryUtilization  *ResourceUtilization `json:"memory_utilization,omitempty"`
	GPUUtilization     *GPUUtilization      `json:"gpu_utilization,omitempty"`
	IOUtilization      *IOUtilization       `json:"io_utilization,omitempty"`
	NetworkUtilization *NetworkUtilization  `json:"network_utilization,omitempty"`
	EnergyUtilization  *ResourceUtilization `json:"energy_utilization,omitempty"`

	// Task-level metrics
	TaskUtilizations []TaskUtilization `json:"task_utilizations,omitempty"`

	// Performance metrics
	PerformanceMetrics *StepPerformanceMetrics `json:"performance_metrics,omitempty"`

	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// TaskUtilization represents utilization metrics for a single task within a job step
type TaskUtilization struct {
	TaskID            int     `json:"task_id"`
	NodeName          string  `json:"node_name"`
	CPUUtilization    float64 `json:"cpu_utilization"`
	MemoryUtilization float64 `json:"memory_utilization"`
	State             string  `json:"state"`
	ExitCode          int     `json:"exit_code"`
}

// StepTaskInfo represents detailed information about a task within a job step
type StepTaskInfo struct {
	TaskID    int           `json:"task_id"`
	NodeName  string        `json:"node_name"`
	LocalID   int           `json:"local_id"` // Local task ID on the node
	State     string        `json:"state"`
	ExitCode  int           `json:"exit_code"`
	CPUTime   time.Duration `json:"cpu_time"`
	MaxRSS    int64         `json:"max_rss_bytes"`
	StartTime *time.Time    `json:"start_time,omitempty"`
	EndTime   *time.Time    `json:"end_time,omitempty"`
}

// StepPerformanceMetrics represents performance analysis for a job step
type StepPerformanceMetrics struct {
	CPUEfficiency     float64 `json:"cpu_efficiency"`
	MemoryEfficiency  float64 `json:"memory_efficiency"`
	IOEfficiency      float64 `json:"io_efficiency"`
	OverallEfficiency float64 `json:"overall_efficiency"`

	// Bottleneck analysis
	PrimaryBottleneck  string `json:"primary_bottleneck"`
	BottleneckSeverity string `json:"bottleneck_severity"`
	ResourceBalance    string `json:"resource_balance"`

	// Performance indicators
	ThroughputMBPS   float64 `json:"throughput_mbps"`
	LatencyMS        float64 `json:"latency_ms"`
	ScalabilityScore float64 `json:"scalability_score"`
}

// ListJobStepsOptions provides filtering options for listing job steps with metrics
type ListJobStepsOptions struct {
	// Basic filtering
	StepStates []string `json:"step_states,omitempty"`
	NodeNames  []string `json:"node_names,omitempty"`
	StepNames  []string `json:"step_names,omitempty"`
	TaskStates []string `json:"task_states,omitempty"`

	// Time-based filtering
	StartTimeAfter  *time.Time     `json:"start_time_after,omitempty"`
	StartTimeBefore *time.Time     `json:"start_time_before,omitempty"`
	EndTimeAfter    *time.Time     `json:"end_time_after,omitempty"`
	EndTimeBefore   *time.Time     `json:"end_time_before,omitempty"`
	MinDuration     *time.Duration `json:"min_duration,omitempty"`
	MaxDuration     *time.Duration `json:"max_duration,omitempty"`

	// Performance filtering
	MinCPUEfficiency     *float64 `json:"min_cpu_efficiency,omitempty"`
	MaxCPUEfficiency     *float64 `json:"max_cpu_efficiency,omitempty"`
	MinMemoryEfficiency  *float64 `json:"min_memory_efficiency,omitempty"`
	MaxMemoryEfficiency  *float64 `json:"max_memory_efficiency,omitempty"`
	MinOverallEfficiency *float64 `json:"min_overall_efficiency,omitempty"`
	MaxOverallEfficiency *float64 `json:"max_overall_efficiency,omitempty"`

	// Include options
	IncludeTaskMetrics         bool `json:"include_task_metrics,omitempty"`
	IncludePerformanceAnalysis bool `json:"include_performance_analysis,omitempty"`
	IncludeResourceTrends      bool `json:"include_resource_trends,omitempty"`
	IncludeBottleneckAnalysis  bool `json:"include_bottleneck_analysis,omitempty"`

	// Pagination
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`

	// Sorting
	SortBy    string `json:"sort_by,omitempty"`    // step_id, start_time, duration, cpu_efficiency, etc.
	SortOrder string `json:"sort_order,omitempty"` // asc, desc
}

// JobStepMetricsList represents a list of job steps with their performance metrics
type JobStepMetricsList struct {
	JobID         string                 `json:"job_id"`
	JobName       string                 `json:"job_name"`
	Steps         []*JobStepWithMetrics  `json:"steps"`
	Summary       *JobStepsSummary       `json:"summary,omitempty"`
	TotalSteps    int                    `json:"total_steps"`
	FilteredSteps int                    `json:"filtered_steps"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// JobStepWithMetrics combines job step details with comprehensive performance metrics
type JobStepWithMetrics struct {
	*JobStepDetails     `json:",inline"`
	*JobStepUtilization `json:",inline"`

	// Additional analytics not included in base structures
	Trends       *StepResourceTrends          `json:"trends,omitempty"`
	Comparison   *StepComparison              `json:"comparison,omitempty"`
	Optimization *StepOptimizationSuggestions `json:"optimization,omitempty"`
}

// JobStepsSummary provides aggregated metrics across all job steps
type JobStepsSummary struct {
	TotalSteps      int           `json:"total_steps"`
	CompletedSteps  int           `json:"completed_steps"`
	FailedSteps     int           `json:"failed_steps"`
	RunningSteps    int           `json:"running_steps"`
	TotalDuration   time.Duration `json:"total_duration"`
	AverageDuration time.Duration `json:"average_duration"`

	// Aggregated efficiency metrics
	AverageCPUEfficiency     float64 `json:"average_cpu_efficiency"`
	AverageMemoryEfficiency  float64 `json:"average_memory_efficiency"`
	AverageIOEfficiency      float64 `json:"average_io_efficiency"`
	AverageOverallEfficiency float64 `json:"average_overall_efficiency"`

	// Resource utilization summaries
	TotalCPUHours     float64 `json:"total_cpu_hours"`
	TotalMemoryGBH    float64 `json:"total_memory_gb_hours"`
	TotalIOOperations int64   `json:"total_io_operations"`
	TotalEnergyUsed   float64 `json:"total_energy_used,omitempty"`

	// Performance analysis
	PrimaryBottlenecks   map[string]int `json:"primary_bottlenecks"`
	BottleneckSeverities map[string]int `json:"bottleneck_severities"`
	MostEfficientStep    *string        `json:"most_efficient_step,omitempty"`
	LeastEfficientStep   *string        `json:"least_efficient_step,omitempty"`

	// Recommendations
	OptimizationPotential float64  `json:"optimization_potential"`
	RecommendedActions    []string `json:"recommended_actions,omitempty"`
}

// StepResourceTrends represents trending data for a job step's resource usage
type StepResourceTrends struct {
	StepID           string             `json:"step_id"`
	CPUTrend         *ResourceTrendData `json:"cpu_trend,omitempty"`
	MemoryTrend      *ResourceTrendData `json:"memory_trend,omitempty"`
	IOTrend          *ResourceTrendData `json:"io_trend,omitempty"`
	NetworkTrend     *ResourceTrendData `json:"network_trend,omitempty"`
	SamplingInterval time.Duration      `json:"sampling_interval"`
	TrendDirection   string             `json:"trend_direction"`  // increasing, decreasing, stable, variable
	TrendConfidence  float64            `json:"trend_confidence"` // 0.0-1.0
}

// ResourceTrendData represents trend data for a specific resource
type ResourceTrendData struct {
	Values       []float64   `json:"values"`
	Timestamps   []time.Time `json:"timestamps"`
	AverageValue float64     `json:"average_value"`
	MinValue     float64     `json:"min_value"`
	MaxValue     float64     `json:"max_value"`
	StandardDev  float64     `json:"standard_deviation"`
	SlopePerHour float64     `json:"slope_per_hour"`
}

// StepComparison provides comparative analysis between job steps
type StepComparison struct {
	StepID                   string   `json:"step_id"`
	RelativeCPUEfficiency    float64  `json:"relative_cpu_efficiency"`    // compared to job average
	RelativeMemoryEfficiency float64  `json:"relative_memory_efficiency"` // compared to job average
	RelativeDuration         float64  `json:"relative_duration"`          // compared to job average
	PerformanceRank          int      `json:"performance_rank"`           // 1 = best performing step
	EfficiencyPercentile     float64  `json:"efficiency_percentile"`      // 0.0-100.0
	ComparisonNotes          []string `json:"comparison_notes,omitempty"`
}

// StepOptimizationSuggestions provides specific optimization recommendations for a job step
type StepOptimizationSuggestions struct {
	StepID               string  `json:"step_id"`
	OverallScore         float64 `json:"overall_score"`         // 0.0-100.0, higher is better
	ImprovementPotential float64 `json:"improvement_potential"` // 0.0-100.0, estimated improvement

	// Resource-specific suggestions
	CPUSuggestions     []OptimizationSuggestion `json:"cpu_suggestions,omitempty"`
	MemorySuggestions  []OptimizationSuggestion `json:"memory_suggestions,omitempty"`
	IOSuggestions      []OptimizationSuggestion `json:"io_suggestions,omitempty"`
	NetworkSuggestions []OptimizationSuggestion `json:"network_suggestions,omitempty"`

	// Configuration recommendations
	RecommendedCPUs       *int     `json:"recommended_cpus,omitempty"`
	RecommendedMemoryMB   *int     `json:"recommended_memory_mb,omitempty"`
	RecommendedNodes      *int     `json:"recommended_nodes,omitempty"`
	AlternativePartitions []string `json:"alternative_partitions,omitempty"`

	// Priority recommendations
	HighPriorityActions   []string `json:"high_priority_actions,omitempty"`
	MediumPriorityActions []string `json:"medium_priority_actions,omitempty"`
	LowPriorityActions    []string `json:"low_priority_actions,omitempty"`
}

// OptimizationSuggestion represents a specific optimization recommendation
type OptimizationSuggestion struct {
	Type                        string  `json:"type"`     // cpu_scaling, memory_tuning, io_optimization, etc.
	Severity                    string  `json:"severity"` // critical, warning, info
	Description                 string  `json:"description"`
	ExpectedBenefit             string  `json:"expected_benefit"`
	ImplementationComplexity    string  `json:"implementation_complexity"` // low, medium, high
	EstimatedImprovementPercent float64 `json:"estimated_improvement_percent"`
	ActionRequired              string  `json:"action_required"`
}

// UserManager provides user-related operations
type UserManager interface {
	// Core user operations
	List(ctx context.Context, opts *ListUsersOptions) (*UserList, error)
	Get(ctx context.Context, userName string) (*User, error)
	Create(ctx context.Context, user *UserCreate) (*UserCreateResponse, error)
	Update(ctx context.Context, userName string, update *UserUpdate) error
	Delete(ctx context.Context, userName string) error

	// CreateAssociation creates a new user-account association
	CreateAssociation(ctx context.Context, accountName string, opts *AssociationOptions) (*AssociationCreateResponse, error)

	// User-account association operations
	GetUserAccounts(ctx context.Context, userName string) ([]*UserAccount, error)
	GetUserQuotas(ctx context.Context, userName string) (*UserQuota, error)
	GetUserDefaultAccount(ctx context.Context, userName string) (*Account, error)

	// Fair-share and priority operations
	GetUserFairShare(ctx context.Context, userName string) (*UserFairShare, error)
	CalculateJobPriority(ctx context.Context, userName string, jobSubmission *JobSubmission) (*JobPriorityInfo, error)

	// Enhanced user-account association operations
	ValidateUserAccountAccess(ctx context.Context, userName, accountName string) (*UserAccessValidation, error)
	GetUserAccountAssociations(ctx context.Context, userName string, opts *ListUserAccountAssociationsOptions) ([]*UserAccountAssociation, error)
	GetBulkUserAccounts(ctx context.Context, userNames []string) (map[string][]*UserAccount, error)
	GetBulkAccountUsers(ctx context.Context, accountNames []string) (map[string][]*UserAccountAssociation, error)
}

// ClusterManager provides cluster configuration and management operations (v0.0.43+)
type ClusterManager interface {
	// List clusters with optional filtering
	List(ctx context.Context, opts *ListClustersOptions) (*ClusterList, error)
	// Get retrieves a specific cluster by name
	Get(ctx context.Context, clusterName string) (*Cluster, error)
	// Create creates a new cluster configuration
	Create(ctx context.Context, cluster *ClusterCreate) (*ClusterCreateResponse, error)
	// Update updates an existing cluster configuration
	Update(ctx context.Context, clusterName string, update *ClusterUpdate) error
	// Delete deletes a cluster configuration
	Delete(ctx context.Context, clusterName string) error
}

// AssociationManager manages user-account-cluster relationships (v0.0.43+)
type AssociationManager interface {
	// List returns associations with optional filtering
	List(ctx context.Context, opts *ListAssociationsOptions) (*AssociationList, error)
	// Get retrieves a specific association
	Get(ctx context.Context, opts *GetAssociationOptions) (*Association, error)
	// Create creates new associations
	Create(ctx context.Context, associations []*AssociationCreate) (*AssociationCreateResponse, error)
	// Update updates existing associations
	Update(ctx context.Context, associations []*AssociationUpdate) error
	// Delete deletes a single association
	Delete(ctx context.Context, opts *DeleteAssociationOptions) error
	// BulkDelete deletes multiple associations
	BulkDelete(ctx context.Context, opts *BulkDeleteOptions) (*BulkDeleteResponse, error)
	// Helper methods
	GetUserAssociations(ctx context.Context, userName string) ([]*Association, error)
	GetAccountAssociations(ctx context.Context, accountName string) ([]*Association, error)
	ValidateAssociation(ctx context.Context, user, account, cluster string) (bool, error)
}

// WCKeyManager manages WCKey (Workload Characterization Key) operations (v0.0.43+)
type WCKeyManager interface {
	// List returns WCKeys with optional filtering
	List(ctx context.Context, opts *WCKeyListOptions) (*WCKeyList, error)
	// Get retrieves a specific WCKey
	Get(ctx context.Context, wckeyName, user, cluster string) (*WCKey, error)
	// Create creates a new WCKey
	Create(ctx context.Context, wckey *WCKeyCreate) (*WCKeyCreateResponse, error)
	// Update updates an existing WCKey
	Update(ctx context.Context, wckeyName, user, cluster string, update *WCKeyUpdate) error
	// Delete deletes a WCKey
	Delete(ctx context.Context, wckeyID string) error
}

// WCKey represents a Workload Characterization Key
type WCKey struct {
	Name    string `json:"name"`
	User    string `json:"user"`
	Cluster string `json:"cluster"`
}

// WCKeyList represents a list of WCKeys
type WCKeyList struct {
	WCKeys []WCKey `json:"wckeys"`
	Total  int     `json:"total"`
}

// WCKeyCreate represents a WCKey creation request
type WCKeyCreate struct {
	Name    string `json:"name"`
	User    string `json:"user"`
	Cluster string `json:"cluster"`
}

// WCKeyCreateResponse represents the response from WCKey creation
type WCKeyCreateResponse struct {
	WCKeyName string `json:"wckey_name"`
}

// WCKeyUpdate represents a WCKey update request
type WCKeyUpdate struct {
	Name *string `json:"name,omitempty"`
}

// WCKeyListOptions provides options for listing WCKeys
type WCKeyListOptions struct {
	Names    []string `json:"names,omitempty"`
	Users    []string `json:"users,omitempty"`
	Clusters []string `json:"clusters,omitempty"`
	Limit    int      `json:"limit,omitempty"`
	Offset   int      `json:"offset,omitempty"`
}

// User represents a SLURM user
type User struct {
	Name                string                 `json:"name"`
	UID                 int                    `json:"uid"`
	DefaultAccount      string                 `json:"default_account"`
	DefaultWCKey        string                 `json:"default_wckey,omitempty"`
	AdminLevel          string                 `json:"admin_level"`
	CoordinatorAccounts []string               `json:"coordinator_accounts,omitempty"`
	Accounts            []UserAccount          `json:"accounts,omitempty"`
	Quotas              *UserQuota             `json:"quotas,omitempty"`
	FairShare           *UserFairShare         `json:"fair_share,omitempty"`
	Associations        []UserAssociation      `json:"associations,omitempty"`
	Created             time.Time              `json:"created"`
	Modified            time.Time              `json:"modified"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}

// UserAccount represents a user's association with an account
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

// UserAssociation represents a user association with cluster/account/partition
type UserAssociation struct {
	ID              uint32            `json:"id"`
	Cluster         string            `json:"cluster"`
	Account         string            `json:"account"`
	Partition       string            `json:"partition,omitempty"`
	User            string            `json:"user"`
	ParentAccount   string            `json:"parent_account,omitempty"`
	Lft             int               `json:"lft"`
	Rgt             int               `json:"rgt"`
	Shares          int               `json:"shares"`
	MaxJobs         int               `json:"max_jobs,omitempty"`
	MaxJobsAccrue   int               `json:"max_jobs_accrue,omitempty"`
	MinPrioThresh   int               `json:"min_prio_thresh,omitempty"`
	MaxSubmitJobs   int               `json:"max_submit_jobs,omitempty"`
	MaxWallDuration int               `json:"max_wall_duration,omitempty"`
	MaxTRES         map[string]string `json:"max_tres,omitempty"`
	MinTRES         map[string]string `json:"min_tres,omitempty"`
	RunningJobs     int               `json:"running_jobs"`
	Usage           *AssociationUsage `json:"usage,omitempty"`
	IsDefault       bool              `json:"is_default"`
	QoS             []string          `json:"qos,omitempty"`
	DefaultQoS      string            `json:"default_qos,omitempty"`
}

// UserQuota represents quota information for a user
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

// UserAccountQuota represents quota limits for a user within a specific account
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

// UserUsage represents current usage statistics for a user
type UserUsage struct {
	UserName     string                        `json:"user_name"`
	RunningJobs  int                           `json:"running_jobs"`
	PendingJobs  int                           `json:"pending_jobs"`
	UsedCPUHours float64                       `json:"used_cpu_hours"`
	UsedGPUHours float64                       `json:"used_gpu_hours,omitempty"`
	UsedWallTime int                           `json:"used_wall_time"`
	TRESUsage    map[string]float64            `json:"tres_usage,omitempty"`
	AccountUsage map[string]*AccountUsageStats `json:"account_usage,omitempty"`
	Efficiency   float64                       `json:"efficiency"`
	LastJobTime  time.Time                     `json:"last_job_time"`
	PeriodStart  time.Time                     `json:"period_start"`
	PeriodEnd    time.Time                     `json:"period_end"`
}

// AccountUsageStats represents usage statistics for a user within an account
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

// UserFairShare represents fair-share information for a user
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

// FairShareNode represents a node in the fair-share tree
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

// JobPriorityFactors represents the individual factors that contribute to job priority
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

// PriorityWeights represents the weights used in priority calculation
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

// JobPriorityInfo represents calculated priority information for a job
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

// AccountFairShare represents fair-share configuration and state for an account
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

// FairShareHierarchy represents the complete fair-share tree structure
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

// UserList represents a list of users
type UserList struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
}

// UserCreate represents a user creation request
type UserCreate struct {
	Name           string   `json:"name"`
	UID            int      `json:"uid,omitempty"`
	DefaultAccount string   `json:"default_account"`
	DefaultWCKey   string   `json:"default_wckey,omitempty"`
	AdminLevel     string   `json:"admin_level,omitempty"`
	Accounts       []string `json:"accounts,omitempty"`
}

// UserCreateResponse represents the response from user creation
type UserCreateResponse struct {
	UserName string `json:"user_name"`
}

// UserUpdate represents a user update request
type UserUpdate struct {
	DefaultAccount      *string  `json:"default_account,omitempty"`
	DefaultWCKey        *string  `json:"default_wckey,omitempty"`
	AdminLevel          *string  `json:"admin_level,omitempty"`
	CoordinatorAccounts []string `json:"coordinator_accounts,omitempty"`
}

// AssociationOptions represents options for creating associations
type AssociationOptions struct {
	Cluster   string `json:"cluster,omitempty"`
	Partition string `json:"partition,omitempty"`
	QoS       string `json:"qos,omitempty"`
	MaxJobs   *int   `json:"max_jobs,omitempty"`
	Priority  *int   `json:"priority,omitempty"`
}

// ListUsersOptions provides filtering options for listing users
type ListUsersOptions struct {
	// Basic filtering
	Names       []string `json:"names,omitempty"`
	Accounts    []string `json:"accounts,omitempty"`
	Clusters    []string `json:"clusters,omitempty"`
	AdminLevels []string `json:"admin_levels,omitempty"`

	// State filtering
	ActiveOnly       bool `json:"active_only,omitempty"`
	CoordinatorsOnly bool `json:"coordinators_only,omitempty"`

	// Include additional data
	WithAccounts     bool `json:"with_accounts,omitempty"`
	WithQuotas       bool `json:"with_quotas,omitempty"`
	WithFairShare    bool `json:"with_fair_share,omitempty"`
	WithAssociations bool `json:"with_associations,omitempty"`
	WithUsage        bool `json:"with_usage,omitempty"`

	// Pagination
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`

	// Sorting
	SortBy    string `json:"sort_by,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
}

// AssociationUsage represents usage data for an association
type AssociationUsage struct {
	UsedCPUHours  float64            `json:"used_cpu_hours"`
	UsedGPUHours  float64            `json:"used_gpu_hours,omitempty"`
	UsedWallTime  int                `json:"used_wall_time"`
	TRESUsage     map[string]float64 `json:"tres_usage,omitempty"`
	JobCount      int                `json:"job_count"`
	RunningJobs   int                `json:"running_jobs"`
	PendingJobs   int                `json:"pending_jobs"`
	CompletedJobs int                `json:"completed_jobs"`
	FailedJobs    int                `json:"failed_jobs"`
	CancelledJobs int                `json:"cancelled_jobs"`
	Efficiency    float64            `json:"efficiency"`
	PeriodStart   time.Time          `json:"period_start"`
	PeriodEnd     time.Time          `json:"period_end"`
}

// QoSLimits represents Quality of Service limits
type QoSLimits struct {
	Name           string         `json:"name"`
	Priority       int            `json:"priority"`
	UsageFactor    float64        `json:"usage_factor"`
	UsageThreshold float64        `json:"usage_threshold"`
	GraceTime      int            `json:"grace_time"`
	MaxJobs        int            `json:"max_jobs"`
	MaxJobsPerUser int            `json:"max_jobs_per_user"`
	MaxSubmitJobs  int            `json:"max_submit_jobs"`
	MaxWallTime    int            `json:"max_wall_time"`
	TRESLimits     map[string]int `json:"tres_limits,omitempty"`
	Flags          []string       `json:"flags,omitempty"`
}

// UserAccountAssociation represents a detailed user-account association with permissions
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

// UserAccessValidation represents the result of user access validation
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

// ListAccountUsersOptions provides filtering options for listing account users
type ListAccountUsersOptions struct {
	// Basic filtering
	Roles            []string `json:"roles,omitempty"`
	Permissions      []string `json:"permissions,omitempty"`
	ActiveOnly       bool     `json:"active_only,omitempty"`
	CoordinatorsOnly bool     `json:"coordinators_only,omitempty"`

	// State filtering
	Partitions []string `json:"partitions,omitempty"`
	QoS        []string `json:"qos,omitempty"`

	// Include additional data
	WithPermissions bool `json:"with_permissions,omitempty"`
	WithQuotas      bool `json:"with_quotas,omitempty"`
	WithUsage       bool `json:"with_usage,omitempty"`
	WithFairShare   bool `json:"with_fair_share,omitempty"`

	// Pagination
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`

	// Sorting
	SortBy    string `json:"sort_by,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
}

// ListUserAccountAssociationsOptions provides filtering options for user account associations
type ListUserAccountAssociationsOptions struct {
	// Basic filtering
	Accounts    []string `json:"accounts,omitempty"`
	Clusters    []string `json:"clusters,omitempty"`
	Partitions  []string `json:"partitions,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`

	// State filtering
	ActiveOnly       bool `json:"active_only,omitempty"`
	DefaultOnly      bool `json:"default_only,omitempty"`
	CoordinatorRoles bool `json:"coordinator_roles,omitempty"`

	// Include additional data
	WithQuotas    bool `json:"with_quotas,omitempty"`
	WithUsage     bool `json:"with_usage,omitempty"`
	WithFairShare bool `json:"with_fair_share,omitempty"`

	// Pagination
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`

	// Sorting
	SortBy    string `json:"sort_by,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
}

// ClientConfig holds configuration for the API client
type ClientConfig struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
	Debug      bool
}

// AccountingQueryOptions defines options for querying SLURM's accounting database
type AccountingQueryOptions struct {
	StartTime  *time.Time `json:"start_time,omitempty"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	User       string     `json:"user,omitempty"`
	Account    string     `json:"account,omitempty"`
	Partition  string     `json:"partition,omitempty"`
	State      string     `json:"state,omitempty"`
	Format     string     `json:"format,omitempty"`
	MaxRecords int        `json:"max_records,omitempty"`
}

// AccountingJobSteps contains job step data from SLURM's accounting database
type AccountingJobSteps struct {
	JobID    string                 `json:"job_id"`
	JobName  string                 `json:"job_name"`
	User     string                 `json:"user"`
	Account  string                 `json:"account"`
	Steps    []StepAccountingRecord `json:"steps"`
	Summary  *JobAccountingSummary  `json:"summary,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// StepAccountingRecord contains detailed accounting information for a job step
type StepAccountingRecord struct {
	StepID    string   `json:"step_id"`
	StepName  string   `json:"step_name"`
	JobID     string   `json:"job_id"`
	UserID    string   `json:"user_id"`
	GroupID   string   `json:"group_id"`
	Account   string   `json:"account"`
	Partition string   `json:"partition"`
	Nodes     []string `json:"nodes"`
	CPUs      int      `json:"cpus"`
	Tasks     int      `json:"tasks"`

	// Timing information
	SubmitTime  time.Time     `json:"submit_time"`
	StartTime   *time.Time    `json:"start_time,omitempty"`
	EndTime     *time.Time    `json:"end_time,omitempty"`
	ElapsedTime time.Duration `json:"elapsed_time"`
	CPUTime     time.Duration `json:"cpu_time"`

	// Resource usage
	AllocCPUs int   `json:"alloc_cpus"`
	ReqCPUs   int   `json:"req_cpus"`
	AllocMem  int64 `json:"alloc_mem"`
	ReqMem    int64 `json:"req_mem"`
	MaxRSS    int64 `json:"max_rss"`
	MaxVMSize int64 `json:"max_vm_size"`

	// I/O and performance
	MaxDiskRead  int64   `json:"max_disk_read"`
	MaxDiskWrite int64   `json:"max_disk_write"`
	AveCPU       float64 `json:"ave_cpu"`
	AveCPUFreq   float64 `json:"ave_cpu_freq"`

	// Status and exit information
	State           string `json:"state"`
	ExitCode        int    `json:"exit_code"`
	DerivedExitCode int    `json:"derived_exit_code"`

	// Additional metadata
	QOS         string                 `json:"qos,omitempty"`
	Constraints string                 `json:"constraints,omitempty"`
	WorkDir     string                 `json:"work_dir,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// JobStepAPIData contains real-time data from SLURM's native job step APIs
type JobStepAPIData struct {
	StepID   string `json:"step_id"`
	JobID    string `json:"job_id"`
	StepName string `json:"step_name"`
	State    string `json:"state"`

	// Real-time resource usage
	CurrentCPUUsage    float64 `json:"current_cpu_usage"`
	CurrentMemoryUsage int64   `json:"current_memory_usage"`
	CurrentIOReads     int64   `json:"current_io_reads"`
	CurrentIOWrites    int64   `json:"current_io_writes"`

	// Task information
	TaskCount      int `json:"task_count"`
	RunningTasks   int `json:"running_tasks"`
	CompletedTasks int `json:"completed_tasks"`
	FailedTasks    int `json:"failed_tasks"`

	// Network information
	NetworkBytesIn  int64 `json:"network_bytes_in"`
	NetworkBytesOut int64 `json:"network_bytes_out"`

	// Performance counters
	ContextSwitches int64 `json:"context_switches"`
	PageFaults      int64 `json:"page_faults"`

	// Timing
	StartTime  *time.Time `json:"start_time,omitempty"`
	LastUpdate time.Time  `json:"last_update"`

	// Additional API data
	ProcessTree     []ProcessInfo          `json:"process_tree,omitempty"`
	EnvironmentVars map[string]string      `json:"environment_vars,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// SacctQueryOptions defines options for querying with sacct command
type SacctQueryOptions struct {
	JobID      string     `json:"job_id,omitempty"`
	User       string     `json:"user,omitempty"`
	Account    string     `json:"account,omitempty"`
	Partition  string     `json:"partition,omitempty"`
	State      []string   `json:"state,omitempty"`
	StartTime  *time.Time `json:"start_time,omitempty"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	Format     []string   `json:"format,omitempty"`      // Fields to include
	Delimiter  string     `json:"delimiter,omitempty"`   // Field delimiter
	NoHeader   bool       `json:"no_header,omitempty"`   // Skip header line
	Parsable   bool       `json:"parsable,omitempty"`    // Machine-readable format
	Brief      bool       `json:"brief,omitempty"`       // Brief output
	Verbose    bool       `json:"verbose,omitempty"`     // Verbose output
	Units      string     `json:"units,omitempty"`       // Units for resource values
	MaxRecords int        `json:"max_records,omitempty"` // Limit number of records
}

// SacctJobStepData contains job step data retrieved via sacct command
type SacctJobStepData struct {
	JobID        string                 `json:"job_id"`
	QueryOptions *SacctQueryOptions     `json:"query_options"`
	Steps        []SacctStepRecord      `json:"steps"`
	TotalSteps   int                    `json:"total_steps"`
	QueryTime    time.Time              `json:"query_time"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// SacctStepRecord represents a single step record from sacct output
type SacctStepRecord struct {
	JobID     string `json:"job_id"`
	StepID    string `json:"step_id"`
	StepName  string `json:"step_name"`
	Account   string `json:"account"`
	User      string `json:"user"`
	Group     string `json:"group"`
	JobName   string `json:"job_name"`
	Partition string `json:"partition"`
	NodeList  string `json:"node_list"`
	NodeCount int    `json:"node_count"`

	// Resource allocation
	AllocCPUs  int    `json:"alloc_cpus"`
	ReqCPUs    int    `json:"req_cpus"`
	AllocMem   string `json:"alloc_mem"` // Can be in various units
	ReqMem     string `json:"req_mem"`   // Can be in various units
	AllocNodes int    `json:"alloc_nodes"`
	ReqNodes   string `json:"req_nodes"`

	// Timing
	Submit     string `json:"submit"`       // Submit time
	Start      string `json:"start"`        // Start time
	End        string `json:"end"`          // End time
	Elapsed    string `json:"elapsed"`      // Elapsed time
	CPUTime    string `json:"cpu_time"`     // Total CPU time
	CPUTimeRAW int64  `json:"cpu_time_raw"` // CPU time in seconds

	// Resource usage
	MaxRSS       string `json:"max_rss"`        // Maximum resident set size
	MaxVMSize    string `json:"max_vm_size"`    // Maximum virtual memory size
	MaxDiskRead  string `json:"max_disk_read"`  // Maximum disk read
	MaxDiskWrite string `json:"max_disk_write"` // Maximum disk write
	MaxPages     int64  `json:"max_pages"`      // Maximum number of pages

	// Performance metrics
	AveCPU       string  `json:"ave_cpu"`        // Average CPU utilization
	AveCPUFreq   string  `json:"ave_cpu_freq"`   // Average CPU frequency
	AvePages     float64 `json:"ave_pages"`      // Average pages per second
	AveDiskRead  string  `json:"ave_disk_read"`  // Average disk read rate
	AveDiskWrite string  `json:"ave_disk_write"` // Average disk write rate

	// Status
	State           string `json:"state"`             // Job state
	ExitCode        string `json:"exit_code"`         // Exit code
	DerivedExitCode string `json:"derived_exit_code"` // Derived exit code

	// Additional fields
	QOS          string `json:"qos,omitempty"`
	Priority     int    `json:"priority,omitempty"`
	ReqTRES      string `json:"req_tres,omitempty"`       // Requested TRES
	AllocTRES    string `json:"alloc_tres,omitempty"`     // Allocated TRES
	TRESUsageIn  string `json:"tres_usage_in,omitempty"`  // TRES usage input
	TRESUsageOut string `json:"tres_usage_out,omitempty"` // TRES usage output

	// Raw data for custom parsing
	RawData  map[string]string      `json:"raw_data,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// JobAccountingSummary provides summary statistics for job accounting data
type JobAccountingSummary struct {
	TotalSteps     int `json:"total_steps"`
	CompletedSteps int `json:"completed_steps"`
	FailedSteps    int `json:"failed_steps"`
	RunningSteps   int `json:"running_steps"`

	// Resource totals
	TotalCPUTime    time.Duration `json:"total_cpu_time"`
	TotalElapsed    time.Duration `json:"total_elapsed"`
	TotalMemoryUsed int64         `json:"total_memory_used"`
	TotalDiskRead   int64         `json:"total_disk_read"`
	TotalDiskWrite  int64         `json:"total_disk_write"`

	// Efficiency metrics
	OverallCPUEff float64 `json:"overall_cpu_efficiency"`
	OverallMemEff float64 `json:"overall_memory_efficiency"`
	OverallIOEff  float64 `json:"overall_io_efficiency"`

	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ProcessInfo contains information about processes in a job step
type ProcessInfo struct {
	PID           int           `json:"pid"`
	PPID          int           `json:"ppid"`
	Command       string        `json:"command"`
	CPUPercent    float64       `json:"cpu_percent"`
	MemoryPercent float64       `json:"memory_percent"`
	RSS           int64         `json:"rss"`
	VMS           int64         `json:"vms"`
	Status        string        `json:"status"`
	StartTime     time.Time     `json:"start_time"`
	Children      []ProcessInfo `json:"children,omitempty"`
}

// CPUAnalytics provides detailed CPU performance analysis and metrics
type CPUAnalytics struct {
	AllocatedCores     int     `json:"allocated_cores"`
	RequestedCores     int     `json:"requested_cores"`
	UsedCores          float64 `json:"used_cores"`
	UtilizationPercent float64 `json:"utilization_percent"`
	EfficiencyPercent  float64 `json:"efficiency_percent"`
	IdleCores          float64 `json:"idle_cores"`
	Oversubscribed     bool    `json:"oversubscribed"`

	// Per-core metrics for detailed analysis
	CoreMetrics []CPUCoreMetric `json:"core_metrics,omitempty"`

	// Thermal and frequency information
	AverageTemperature     float64 `json:"average_temperature"`
	MaxTemperature         float64 `json:"max_temperature"`
	ThermalThrottleEvents  int64   `json:"thermal_throttle_events"`
	AverageFrequency       float64 `json:"average_frequency_ghz"`
	MaxFrequency           float64 `json:"max_frequency_ghz"`
	FrequencyScalingEvents int64   `json:"frequency_scaling_events"`

	// Threading and context metrics
	ContextSwitches  int64   `json:"context_switches"`
	Interrupts       int64   `json:"interrupts"`
	SoftInterrupts   int64   `json:"soft_interrupts"`
	LoadAverage1Min  float64 `json:"load_average_1min"`
	LoadAverage5Min  float64 `json:"load_average_5min"`
	LoadAverage15Min float64 `json:"load_average_15min"`

	// Cache performance metrics
	L1CacheHitRate float64 `json:"l1_cache_hit_rate"`
	L2CacheHitRate float64 `json:"l2_cache_hit_rate"`
	L3CacheHitRate float64 `json:"l3_cache_hit_rate"`
	L1CacheMisses  int64   `json:"l1_cache_misses"`
	L2CacheMisses  int64   `json:"l2_cache_misses"`
	L3CacheMisses  int64   `json:"l3_cache_misses"`

	// Instruction-level metrics
	InstructionsPerCycle int64 `json:"instructions_per_cycle"`
	BranchMispredictions int64 `json:"branch_mispredictions"`
	TotalInstructions    int64 `json:"total_instructions"`

	// Performance recommendations and bottleneck analysis
	Recommendations []OptimizationRecommendation `json:"recommendations,omitempty"`
	Bottlenecks     []PerformanceBottleneck      `json:"bottlenecks,omitempty"`

	// Additional metadata and context
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// CPUCoreMetric contains performance metrics for an individual CPU core
type CPUCoreMetric struct {
	CoreID          int     `json:"core_id"`
	Utilization     float64 `json:"utilization_percent"`
	Frequency       float64 `json:"frequency_ghz"`
	Temperature     float64 `json:"temperature_celsius"`
	LoadAverage     float64 `json:"load_average"`
	ContextSwitches int64   `json:"context_switches"`
	Interrupts      int64   `json:"interrupts"`
}

// MemoryAnalytics provides detailed memory performance analysis and metrics
type MemoryAnalytics struct {
	AllocatedBytes     int64   `json:"allocated_bytes"`
	RequestedBytes     int64   `json:"requested_bytes"`
	UsedBytes          int64   `json:"used_bytes"`
	UtilizationPercent float64 `json:"utilization_percent"`
	EfficiencyPercent  float64 `json:"efficiency_percent"`
	FreeBytes          int64   `json:"free_bytes"`
	Overcommitted      bool    `json:"overcommitted"`

	// Memory breakdown by type
	ResidentSetSize   int64 `json:"resident_set_size"`
	VirtualMemorySize int64 `json:"virtual_memory_size"`
	SharedMemory      int64 `json:"shared_memory"`
	BufferedMemory    int64 `json:"buffered_memory"`
	CachedMemory      int64 `json:"cached_memory"`

	// NUMA topology and memory locality
	NUMANodes []NUMANodeMetrics `json:"numa_nodes,omitempty"`

	// Memory bandwidth and performance
	BandwidthUtilization float64 `json:"bandwidth_utilization_percent"`
	MemoryBandwidthMBPS  int64   `json:"memory_bandwidth_mbps"`
	PeakBandwidthMBPS    int64   `json:"peak_bandwidth_mbps"`

	// Page and swap metrics
	PageFaults      int64 `json:"page_faults"`
	MajorPageFaults int64 `json:"major_page_faults"`
	MinorPageFaults int64 `json:"minor_page_faults"`
	PageSwaps       int64 `json:"page_swaps"`

	// Memory access patterns
	RandomAccess     float64 `json:"random_access_percent"`
	SequentialAccess float64 `json:"sequential_access_percent"`
	LocalityScore    float64 `json:"locality_score"`

	// Memory leak detection and analysis
	MemoryLeaks []MemoryLeak `json:"memory_leaks,omitempty"`

	// Performance recommendations and bottleneck analysis
	Recommendations []OptimizationRecommendation `json:"recommendations,omitempty"`
	Bottlenecks     []PerformanceBottleneck      `json:"bottlenecks,omitempty"`

	// Additional metadata and context
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NUMANodeMetrics contains memory and CPU metrics for a specific NUMA node
type NUMANodeMetrics struct {
	NodeID           int     `json:"node_id"`
	CPUCores         int     `json:"cpu_cores"`
	MemoryTotal      int64   `json:"memory_total"`
	MemoryUsed       int64   `json:"memory_used"`
	MemoryFree       int64   `json:"memory_free"`
	CPUUtilization   float64 `json:"cpu_utilization_percent"`
	MemoryBandwidth  int64   `json:"memory_bandwidth_mbps"`
	LocalAccesses    float64 `json:"local_accesses_percent"`
	RemoteAccesses   float64 `json:"remote_accesses_percent"`
	InterconnectLoad float64 `json:"interconnect_load_percent"`
}

// MemoryLeak contains information about detected memory leaks
type MemoryLeak struct {
	LeakType    string `json:"leak_type"`   // gradual, sudden, cyclic, etc.
	SizeBytes   int64  `json:"size_bytes"`  // Current size of the leak
	GrowthRate  int64  `json:"growth_rate"` // Bytes per second growth
	Location    string `json:"location"`    // Where the leak was detected
	Severity    string `json:"severity"`    // low, medium, high, critical
	Description string `json:"description"` // Human-readable description
}

// IOAnalytics provides detailed I/O performance analysis and metrics
type IOAnalytics struct {
	ReadBytes          int64   `json:"read_bytes"`
	WriteBytes         int64   `json:"write_bytes"`
	ReadOperations     int64   `json:"read_operations"`
	WriteOperations    int64   `json:"write_operations"`
	UtilizationPercent float64 `json:"utilization_percent"`
	EfficiencyPercent  float64 `json:"efficiency_percent"`

	// Bandwidth metrics
	AverageReadBandwidth  float64 `json:"average_read_bandwidth_mbps"`
	AverageWriteBandwidth float64 `json:"average_write_bandwidth_mbps"`
	PeakReadBandwidth     float64 `json:"peak_read_bandwidth_mbps"`
	PeakWriteBandwidth    float64 `json:"peak_write_bandwidth_mbps"`

	// Latency metrics
	AverageReadLatency  float64 `json:"average_read_latency_ms"`
	AverageWriteLatency float64 `json:"average_write_latency_ms"`
	MaxReadLatency      float64 `json:"max_read_latency_ms"`
	MaxWriteLatency     float64 `json:"max_write_latency_ms"`

	// Queue depth and wait times
	QueueDepth    float64 `json:"queue_depth"`
	MaxQueueDepth float64 `json:"max_queue_depth"`
	QueueTime     float64 `json:"queue_time_ms"`

	// Access patterns
	RandomAccessPercent     float64 `json:"random_access_percent"`
	SequentialAccessPercent float64 `json:"sequential_access_percent"`

	// I/O size distribution
	AverageIOSize int64 `json:"average_io_size_bytes"`
	MaxIOSize     int64 `json:"max_io_size_bytes"`
	MinIOSize     int64 `json:"min_io_size_bytes"`

	// Storage device information
	StorageDevices []StorageDevice `json:"storage_devices,omitempty"`

	// Performance recommendations and bottleneck analysis
	Recommendations []OptimizationRecommendation `json:"recommendations,omitempty"`
	Bottlenecks     []PerformanceBottleneck      `json:"bottlenecks,omitempty"`

	// Additional metadata and context
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// StorageDevice contains information and metrics for a storage device
type StorageDevice struct {
	DeviceName     string  `json:"device_name"`
	DeviceType     string  `json:"device_type"` // ssd, hdd, nvme_ssd, etc.
	MountPoint     string  `json:"mount_point"`
	TotalCapacity  int64   `json:"total_capacity"`
	UsedCapacity   int64   `json:"used_capacity"`
	AvailCapacity  int64   `json:"available_capacity"`
	Utilization    float64 `json:"utilization_percent"`
	IOPS           int64   `json:"iops"`
	ThroughputMBPS int64   `json:"throughput_mbps"`
}

// JobComprehensiveAnalytics combines all analytics into a complete performance report
type JobComprehensiveAnalytics struct {
	JobID     uint32        `json:"job_id"`
	JobName   string        `json:"job_name"`
	StartTime time.Time     `json:"start_time"`
	EndTime   *time.Time    `json:"end_time,omitempty"`
	Duration  time.Duration `json:"duration"`
	Status    string        `json:"status"`

	// Individual analytics components
	CPUAnalytics    *CPUAnalytics    `json:"cpu_analytics"`
	MemoryAnalytics *MemoryAnalytics `json:"memory_analytics"`
	IOAnalytics     *IOAnalytics     `json:"io_analytics"`

	// Overall performance metrics
	OverallEfficiency float64 `json:"overall_efficiency_percent"`

	// Cross-resource analysis
	CrossResourceAnalysis *CrossResourceAnalysis `json:"cross_resource_analysis"`

	// Optimization recommendations
	OptimalConfiguration *OptimalJobConfiguration `json:"optimal_configuration"`

	// Combined recommendations and bottlenecks
	Recommendations []OptimizationRecommendation `json:"recommendations,omitempty"`
	Bottlenecks     []PerformanceBottleneck      `json:"bottlenecks,omitempty"`

	// Additional metadata and context
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// CrossResourceAnalysis analyzes relationships and bottlenecks across different resource types
type CrossResourceAnalysis struct {
	PrimaryBottleneck     string  `json:"primary_bottleneck"`   // cpu, memory, io, network, none
	SecondaryBottleneck   string  `json:"secondary_bottleneck"` // cpu, memory, io, network, none
	BottleneckSeverity    string  `json:"bottleneck_severity"`  // none, low, medium, high, critical
	ResourceBalance       string  `json:"resource_balance"`     // optimal, cpu_bound, memory_bound, io_bound, unbalanced
	OptimizationPotential float64 `json:"optimization_potential_percent"`
	ScalabilityScore      float64 `json:"scalability_score"`
	ResourceWaste         float64 `json:"resource_waste_percent"`
	LoadBalanceScore      float64 `json:"load_balance_score"`
}

// OptimalJobConfiguration provides recommended resource allocation for future job runs
type OptimalJobConfiguration struct {
	RecommendedCPUs    int               `json:"recommended_cpus"`
	RecommendedMemory  int64             `json:"recommended_memory_bytes"`
	RecommendedNodes   int               `json:"recommended_nodes"`
	RecommendedRuntime int               `json:"recommended_runtime_minutes"`
	ExpectedSpeedup    float64           `json:"expected_speedup"`
	CostReduction      float64           `json:"cost_reduction_percent"`
	ConfigChanges      map[string]string `json:"config_changes,omitempty"`
}

// Historical Performance Tracking Data Structures

// PerformanceHistoryOptions configures historical performance queries
type PerformanceHistoryOptions struct {
	StartTime     *time.Time `json:"start_time,omitempty"`
	EndTime       *time.Time `json:"end_time,omitempty"`
	Interval      string     `json:"interval,omitempty"`     // hourly, daily, weekly
	MetricTypes   []string   `json:"metric_types,omitempty"` // cpu, memory, io, gpu
	IncludeSteps  bool       `json:"include_steps"`
	IncludeTrends bool       `json:"include_trends"`
}

// JobPerformanceHistory contains historical performance data for a job
type JobPerformanceHistory struct {
	JobID     string    `json:"job_id"`
	JobName   string    `json:"job_name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`

	// Time series data
	TimeSeriesData []PerformanceSnapshot `json:"time_series_data"`

	// Aggregated statistics
	Statistics PerformanceStatistics `json:"statistics"`

	// Trend analysis
	Trends *PerformanceTrendAnalysis `json:"trends,omitempty"`

	// Anomalies detected
	Anomalies []PerformanceAnomaly `json:"anomalies,omitempty"`
}

// PerformanceSnapshot represents performance metrics at a point in time
type PerformanceSnapshot struct {
	Timestamp         time.Time `json:"timestamp"`
	CPUUtilization    float64   `json:"cpu_utilization"`
	MemoryUtilization float64   `json:"memory_utilization"`
	IOBandwidth       float64   `json:"io_bandwidth_mbps"`
	GPUUtilization    float64   `json:"gpu_utilization,omitempty"`
	NetworkBandwidth  float64   `json:"network_bandwidth_mbps,omitempty"`
	PowerUsage        float64   `json:"power_usage_watts,omitempty"`

	// Efficiency metrics at this snapshot
	Efficiency float64 `json:"efficiency_score"`
}

// PerformanceStatistics contains statistical analysis of performance data
type PerformanceStatistics struct {
	AverageCPU        float64 `json:"average_cpu"`
	AverageMemory     float64 `json:"average_memory"`
	AverageIO         float64 `json:"average_io"`
	AverageEfficiency float64 `json:"average_efficiency"`

	PeakCPU    float64 `json:"peak_cpu"`
	PeakMemory float64 `json:"peak_memory"`
	PeakIO     float64 `json:"peak_io"`

	MinCPU    float64 `json:"min_cpu"`
	MinMemory float64 `json:"min_memory"`
	MinIO     float64 `json:"min_io"`

	StdDevCPU    float64 `json:"std_dev_cpu"`
	StdDevMemory float64 `json:"std_dev_memory"`
	StdDevIO     float64 `json:"std_dev_io"`
}

// PerformanceTrendAnalysis contains trend analysis results
type PerformanceTrendAnalysis struct {
	CPUTrend        TrendInfo `json:"cpu_trend"`
	MemoryTrend     TrendInfo `json:"memory_trend"`
	IOTrend         TrendInfo `json:"io_trend"`
	EfficiencyTrend TrendInfo `json:"efficiency_trend"`

	// Predicted future values
	PredictedCPU     float64       `json:"predicted_cpu"`
	PredictedMemory  float64       `json:"predicted_memory"`
	PredictedRuntime time.Duration `json:"predicted_runtime"`
}

// TrendInfo describes a performance trend
type TrendInfo struct {
	Direction  string  `json:"direction"` // increasing, decreasing, stable
	Slope      float64 `json:"slope"`
	Confidence float64 `json:"confidence"`
	ChangeRate float64 `json:"change_rate_percent"`
}

// PerformanceAnomaly represents an anomaly in performance metrics
type PerformanceAnomaly struct {
	Timestamp   time.Time `json:"timestamp"`
	Type        string    `json:"type"`     // spike, drop, pattern_change
	Metric      string    `json:"metric"`   // cpu, memory, io, etc.
	Severity    string    `json:"severity"` // low, medium, high, critical
	Value       float64   `json:"value"`
	Expected    float64   `json:"expected"`
	Deviation   float64   `json:"deviation_percent"`
	Description string    `json:"description"`
}

// TrendAnalysisOptions configures cluster-wide trend analysis
type TrendAnalysisOptions struct {
	StartTime   *time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	Granularity string     `json:"granularity"` // hourly, daily, weekly, monthly
	Partitions  []string   `json:"partitions,omitempty"`
	UserFilter  []string   `json:"users,omitempty"`
	JobTypes    []string   `json:"job_types,omitempty"`
	MinJobSize  *int       `json:"min_job_size,omitempty"`
}

// PerformanceTrends contains cluster-wide performance trends
type PerformanceTrends struct {
	TimeRange   TimeRange `json:"time_range"`
	Granularity string    `json:"granularity"`

	// Cluster-wide metrics
	ClusterUtilization []UtilizationPoint `json:"cluster_utilization"`
	ClusterEfficiency  []EfficiencyPoint  `json:"cluster_efficiency"`

	// Partition-specific trends
	PartitionTrends map[string]*PartitionTrend `json:"partition_trends"`

	// Resource type trends
	CPUTrends    ResourceTrend `json:"cpu_trends"`
	MemoryTrends ResourceTrend `json:"memory_trends"`
	GPUTrends    ResourceTrend `json:"gpu_trends"`

	// Job characteristics trends
	JobSizeTrends     []JobSizeTrend     `json:"job_size_trends"`
	JobDurationTrends []JobDurationTrend `json:"job_duration_trends"`

	// Top patterns and insights
	Insights []TrendInsight `json:"insights"`
}

// TimeRange represents a time period
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// UtilizationPoint represents utilization at a point in time
type UtilizationPoint struct {
	Timestamp   time.Time `json:"timestamp"`
	Utilization float64   `json:"utilization"`
	JobCount    int       `json:"job_count"`
}

// EfficiencyPoint represents efficiency at a point in time
type EfficiencyPoint struct {
	Timestamp  time.Time `json:"timestamp"`
	Efficiency float64   `json:"efficiency"`
	JobCount   int       `json:"job_count"`
}

// PartitionTrend contains trends for a specific partition
type PartitionTrend struct {
	PartitionName string             `json:"partition_name"`
	Utilization   []UtilizationPoint `json:"utilization"`
	Efficiency    []EfficiencyPoint  `json:"efficiency"`
	JobCounts     []JobCountPoint    `json:"job_counts"`
	QueueLength   []QueueLengthPoint `json:"queue_length"`
}

// ResourceTrend contains trend data for a resource type
type ResourceTrend struct {
	Average    []float64   `json:"average"`
	Peak       []float64   `json:"peak"`
	Timestamps []time.Time `json:"timestamps"`
	Trend      TrendInfo   `json:"trend"`
}

// JobSizeTrend tracks job size distribution over time
type JobSizeTrend struct {
	Timestamp   time.Time `json:"timestamp"`
	SmallJobs   int       `json:"small_jobs"`  // < 16 cores
	MediumJobs  int       `json:"medium_jobs"` // 16-128 cores
	LargeJobs   int       `json:"large_jobs"`  // > 128 cores
	AverageSize float64   `json:"average_size"`
}

// JobDurationTrend tracks job duration distribution over time
type JobDurationTrend struct {
	Timestamp       time.Time     `json:"timestamp"`
	ShortJobs       int           `json:"short_jobs"`  // < 1 hour
	MediumJobs      int           `json:"medium_jobs"` // 1-12 hours
	LongJobs        int           `json:"long_jobs"`   // > 12 hours
	AverageDuration time.Duration `json:"average_duration"`
}

// JobCountPoint represents job count at a point in time
type JobCountPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Running   int       `json:"running"`
	Pending   int       `json:"pending"`
	Total     int       `json:"total"`
}

// QueueLengthPoint represents queue length at a point in time
type QueueLengthPoint struct {
	Timestamp   time.Time     `json:"timestamp"`
	QueueLength int           `json:"queue_length"`
	WaitTime    time.Duration `json:"average_wait_time"`
}

// TrendInsight represents an insight derived from trend analysis
type TrendInsight struct {
	Type           string    `json:"type"`     // pattern, anomaly, prediction
	Category       string    `json:"category"` // utilization, efficiency, capacity
	Severity       string    `json:"severity"` // info, warning, critical
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Timestamp      time.Time `json:"timestamp"`
	Confidence     float64   `json:"confidence"`
	Recommendation string    `json:"recommendation,omitempty"`
}

// EfficiencyTrendOptions configures user efficiency trend analysis
type EfficiencyTrendOptions struct {
	StartTime        *time.Time `json:"start_time,omitempty"`
	EndTime          *time.Time `json:"end_time,omitempty"`
	Granularity      string     `json:"granularity"` // daily, weekly, monthly
	JobTypes         []string   `json:"job_types,omitempty"`
	Partitions       []string   `json:"partitions,omitempty"`
	CompareToAverage bool       `json:"compare_to_average"`
}

// UserEfficiencyTrends contains efficiency trends for a specific user
type UserEfficiencyTrends struct {
	UserID    string    `json:"user_id"`
	TimeRange TimeRange `json:"time_range"`

	// Efficiency over time
	EfficiencyHistory []EfficiencyDataPoint `json:"efficiency_history"`

	// Resource utilization trends
	CPUUtilizationTrend    []float64 `json:"cpu_utilization_trend"`
	MemoryUtilizationTrend []float64 `json:"memory_utilization_trend"`

	// Job characteristics
	JobCountTrend    []int           `json:"job_count_trend"`
	JobSizeTrend     []float64       `json:"job_size_trend"`
	JobDurationTrend []time.Duration `json:"job_duration_trend"`

	// Comparative analysis
	AverageEfficiency        float64 `json:"average_efficiency"`
	ClusterAverageEfficiency float64 `json:"cluster_average_efficiency"`
	EfficiencyRank           int     `json:"efficiency_rank"`
	EfficiencyPercentile     float64 `json:"efficiency_percentile"`

	// Improvement tracking
	ImprovementRate float64   `json:"improvement_rate"`
	BestPeriod      TimeRange `json:"best_period"`
	WorstPeriod     TimeRange `json:"worst_period"`

	// Recommendations
	Recommendations []string `json:"recommendations"`
}

// EfficiencyDataPoint represents efficiency at a specific time
type EfficiencyDataPoint struct {
	Timestamp   time.Time `json:"timestamp"`
	Efficiency  float64   `json:"efficiency"`
	JobCount    int       `json:"job_count"`
	CPUHours    float64   `json:"cpu_hours"`
	MemoryGBH   float64   `json:"memory_gb_hours"`
	WastedHours float64   `json:"wasted_cpu_hours"`
}

// BatchAnalysisOptions configures batch job analysis
type BatchAnalysisOptions struct {
	IncludeDetails     bool     `json:"include_details"`
	IncludeComparison  bool     `json:"include_comparison"`
	ComparisonBaseline string   `json:"comparison_baseline"` // average, best, worst
	MetricTypes        []string `json:"metric_types"`
}

// BatchJobAnalysis contains analysis results for a batch of jobs
type BatchJobAnalysis struct {
	JobCount      int       `json:"job_count"`
	AnalyzedCount int       `json:"analyzed_count"`
	FailedCount   int       `json:"failed_count"`
	TimeRange     TimeRange `json:"time_range"`

	// Aggregate statistics
	AggregateStats BatchStatistics `json:"aggregate_stats"`

	// Per-job analysis
	JobAnalyses []JobAnalysisSummary `json:"job_analyses,omitempty"`

	// Comparative analysis
	Comparison *BatchComparison `json:"comparison,omitempty"`

	// Patterns and insights
	Patterns []BatchPattern `json:"patterns"`
	Outliers []string       `json:"outlier_job_ids"`

	// Recommendations
	BatchRecommendations []BatchRecommendation `json:"recommendations"`
}

// BatchStatistics contains aggregate statistics for batch analysis
type BatchStatistics struct {
	TotalCPUHours     float64       `json:"total_cpu_hours"`
	TotalMemoryGBH    float64       `json:"total_memory_gb_hours"`
	TotalGPUHours     float64       `json:"total_gpu_hours"`
	AverageEfficiency float64       `json:"average_efficiency"`
	MedianEfficiency  float64       `json:"median_efficiency"`
	StdDevEfficiency  float64       `json:"std_dev_efficiency"`
	TotalWaste        ResourceWaste `json:"total_waste"`
	AverageRuntime    time.Duration `json:"average_runtime"`
	SuccessRate       float64       `json:"success_rate"`
}

// JobAnalysisSummary contains summary of individual job analysis
type JobAnalysisSummary struct {
	JobID             string        `json:"job_id"`
	JobName           string        `json:"job_name"`
	Efficiency        float64       `json:"efficiency"`
	Runtime           time.Duration `json:"runtime"`
	CPUUtilization    float64       `json:"cpu_utilization"`
	MemoryUtilization float64       `json:"memory_utilization"`
	Status            string        `json:"status"`
	Issues            []string      `json:"issues,omitempty"`
}

// BatchComparison compares batch performance against baseline
type BatchComparison struct {
	Baseline        string  `json:"baseline"`
	EfficiencyDelta float64 `json:"efficiency_delta"`
	RuntimeDelta    float64 `json:"runtime_delta"`
	WasteDelta      float64 `json:"waste_delta"`
	CostDelta       float64 `json:"cost_delta"`
}

// BatchPattern represents a pattern found in batch analysis
type BatchPattern struct {
	Type        string   `json:"type"` // efficiency_degradation, resource_imbalance, etc.
	Description string   `json:"description"`
	JobCount    int      `json:"job_count"`
	JobIDs      []string `json:"job_ids"`
	Impact      string   `json:"impact"`
	Confidence  float64  `json:"confidence"`
}

// BatchRecommendation provides recommendations for the batch
type BatchRecommendation struct {
	Category     string   `json:"category"` // resource_allocation, scheduling, configuration
	Priority     string   `json:"priority"` // high, medium, low
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Impact       string   `json:"impact"`
	JobsAffected []string `json:"jobs_affected"`
}

// ResourceWaste contains resource waste metrics
type ResourceWaste struct {
	CPUCoreHours  float64 `json:"cpu_core_hours"`
	MemoryGBHours float64 `json:"memory_gb_hours"`
	GPUHours      float64 `json:"gpu_hours"`
	EstimatedCost float64 `json:"estimated_cost"`
}

// WorkflowAnalysisOptions configures workflow performance analysis
type WorkflowAnalysisOptions struct {
	IncludeDependencies bool       `json:"include_dependencies"`
	IncludeBottlenecks  bool       `json:"include_bottlenecks"`
	IncludeOptimization bool       `json:"include_optimization"`
	TimeWindow          *TimeRange `json:"time_window,omitempty"`
}

// WorkflowPerformance contains performance analysis for a workflow
type WorkflowPerformance struct {
	WorkflowID    string     `json:"workflow_id"`
	WorkflowName  string     `json:"workflow_name"`
	TotalJobs     int        `json:"total_jobs"`
	CompletedJobs int        `json:"completed_jobs"`
	StartTime     time.Time  `json:"start_time"`
	EndTime       *time.Time `json:"end_time,omitempty"`

	// Overall metrics
	TotalDuration     time.Duration `json:"total_duration"`
	WallClockTime     time.Duration `json:"wall_clock_time"`
	Parallelization   float64       `json:"parallelization_efficiency"`
	OverallEfficiency float64       `json:"overall_efficiency"`

	// Stage analysis
	Stages []WorkflowStage `json:"stages"`

	// Critical path analysis
	CriticalPath         []string      `json:"critical_path_job_ids"`
	CriticalPathDuration time.Duration `json:"critical_path_duration"`

	// Bottleneck analysis
	Bottlenecks []WorkflowBottleneck `json:"bottlenecks,omitempty"`

	// Dependencies
	Dependencies *WorkflowDependencies `json:"dependencies,omitempty"`

	// Optimization suggestions
	Optimization *WorkflowOptimization `json:"optimization,omitempty"`
}

// WorkflowStage represents a stage in the workflow
type WorkflowStage struct {
	StageID     string        `json:"stage_id"`
	StageName   string        `json:"stage_name"`
	JobIDs      []string      `json:"job_ids"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
	Duration    time.Duration `json:"duration"`
	Efficiency  float64       `json:"efficiency"`
	Parallelism int           `json:"parallelism"`
	Status      string        `json:"status"`
}

// WorkflowBottleneck identifies a bottleneck in the workflow
type WorkflowBottleneck struct {
	Type        string        `json:"type"`     // resource, dependency, scheduling
	Location    string        `json:"location"` // job_id or stage_id
	Description string        `json:"description"`
	Impact      time.Duration `json:"impact_duration"`
	Severity    string        `json:"severity"`
}

// WorkflowDependencies describes job dependencies in the workflow
type WorkflowDependencies struct {
	DependencyGraph map[string][]string `json:"dependency_graph"`
	MaxDepth        int                 `json:"max_depth"`
	MaxWidth        int                 `json:"max_width"`
	TotalEdges      int                 `json:"total_edges"`
}

// WorkflowOptimization provides optimization suggestions for the workflow
type WorkflowOptimization struct {
	PotentialSpeedup     float64                   `json:"potential_speedup"`
	OptimizedDuration    time.Duration             `json:"optimized_duration"`
	OptimizedParallelism int                       `json:"optimized_parallelism"`
	RecommendedChanges   []string                  `json:"recommended_changes"`
	ResourceReallocation map[string]ResourceChange `json:"resource_reallocation"`
}

// ResourceChange describes a recommended resource change
type ResourceChange struct {
	JobID             string `json:"job_id"`
	CurrentCPUs       int    `json:"current_cpus"`
	RecommendedCPUs   int    `json:"recommended_cpus"`
	CurrentMemory     int64  `json:"current_memory"`
	RecommendedMemory int64  `json:"recommended_memory"`
	Reason            string `json:"reason"`
}

// ReportOptions configures efficiency report generation
type ReportOptions struct {
	ReportType    string    `json:"report_type"` // summary, detailed, executive
	TimeRange     TimeRange `json:"time_range"`
	Partitions    []string  `json:"partitions,omitempty"`
	Users         []string  `json:"users,omitempty"`
	IncludeCharts bool      `json:"include_charts"`
	Format        string    `json:"format"` // json, html, pdf
}

// EfficiencyReport contains a comprehensive efficiency report
type EfficiencyReport struct {
	ReportID    string    `json:"report_id"`
	GeneratedAt time.Time `json:"generated_at"`
	TimeRange   TimeRange `json:"time_range"`
	ReportType  string    `json:"report_type"`

	// Executive summary
	Summary ExecutiveSummary `json:"summary"`

	// Detailed sections
	ClusterOverview   *ClusterOverview    `json:"cluster_overview,omitempty"`
	PartitionAnalysis []PartitionAnalysis `json:"partition_analysis,omitempty"`
	UserAnalysis      []UserAnalysis      `json:"user_analysis,omitempty"`
	ResourceAnalysis  *ResourceAnalysis   `json:"resource_analysis,omitempty"`

	// Trends and patterns
	TrendAnalysis *ReportTrendAnalysis `json:"trend_analysis,omitempty"`

	// Recommendations
	Recommendations []ReportRecommendation `json:"recommendations"`

	// Charts and visualizations
	Charts []ChartData `json:"charts,omitempty"`
}

// ExecutiveSummary provides high-level report summary
type ExecutiveSummary struct {
	TotalJobs            int      `json:"total_jobs"`
	AverageEfficiency    float64  `json:"average_efficiency"`
	TotalCPUHours        float64  `json:"total_cpu_hours"`
	WastedCPUHours       float64  `json:"wasted_cpu_hours"`
	EstimatedCostSavings float64  `json:"estimated_cost_savings"`
	KeyFindings          []string `json:"key_findings"`
	ImprovementAreas     []string `json:"improvement_areas"`
}

// ClusterOverview provides cluster-wide efficiency metrics
type ClusterOverview struct {
	ClusterUtilization float64 `json:"cluster_utilization"`
	ClusterEfficiency  float64 `json:"cluster_efficiency"`
	TotalNodes         int     `json:"total_nodes"`
	ActiveNodes        int     `json:"active_nodes"`
	TotalCPUCores      int     `json:"total_cpu_cores"`
	TotalMemoryGB      float64 `json:"total_memory_gb"`
	TotalGPUs          int     `json:"total_gpus"`
}

// PartitionAnalysis contains efficiency analysis for a partition
type PartitionAnalysis struct {
	PartitionName  string   `json:"partition_name"`
	JobCount       int      `json:"job_count"`
	Utilization    float64  `json:"utilization"`
	Efficiency     float64  `json:"efficiency"`
	TopUsers       []string `json:"top_users"`
	CommonJobTypes []string `json:"common_job_types"`
	Issues         []string `json:"issues"`
}

// UserAnalysis contains efficiency analysis for a user
type UserAnalysis struct {
	UserID            string  `json:"user_id"`
	JobCount          int     `json:"job_count"`
	AverageEfficiency float64 `json:"average_efficiency"`
	TotalCPUHours     float64 `json:"total_cpu_hours"`
	WastedHours       float64 `json:"wasted_hours"`
	Rank              int     `json:"rank"`
	Trend             string  `json:"trend"` // improving, stable, declining
}

// ResourceAnalysis provides resource-specific efficiency analysis
type ResourceAnalysis struct {
	CPUAnalysis    ResourceTypeAnalysis `json:"cpu_analysis"`
	MemoryAnalysis ResourceTypeAnalysis `json:"memory_analysis"`
	GPUAnalysis    ResourceTypeAnalysis `json:"gpu_analysis"`
	IOAnalysis     ResourceTypeAnalysis `json:"io_analysis"`
}

// ResourceTypeAnalysis contains analysis for a specific resource type
type ResourceTypeAnalysis struct {
	AverageUtilization    float64  `json:"average_utilization"`
	PeakUtilization       float64  `json:"peak_utilization"`
	WastePercentage       float64  `json:"waste_percentage"`
	TopWasters            []string `json:"top_waster_job_ids"`
	OptimizationPotential float64  `json:"optimization_potential"`
}

// ReportTrendAnalysis contains trend analysis for the report
type ReportTrendAnalysis struct {
	EfficiencyTrend     TrendInfo `json:"efficiency_trend"`
	UtilizationTrend    TrendInfo `json:"utilization_trend"`
	JobVolumeTrend      TrendInfo `json:"job_volume_trend"`
	PredictedEfficiency float64   `json:"predicted_efficiency_next_period"`
}

// ReportRecommendation provides actionable recommendations
type ReportRecommendation struct {
	Category         string   `json:"category"`
	Priority         string   `json:"priority"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	ExpectedImpact   string   `json:"expected_impact"`
	AffectedEntities []string `json:"affected_entities"`
	Implementation   string   `json:"implementation_guidance"`
}

// ChartData represents data for visualization
type ChartData struct {
	ChartID     string                 `json:"chart_id"`
	ChartType   string                 `json:"chart_type"` // line, bar, pie, heatmap
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Data        map[string]interface{} `json:"data"`
	Options     map[string]interface{} `json:"options"`
}

// === Standalone Operation Types ===

// License represents a SLURM license
type License struct {
	Name       string    `json:"name"`
	Total      int       `json:"total"`
	Used       int       `json:"used"`
	Available  int       `json:"available"`
	Reserved   int       `json:"reserved,omitempty"`
	Remote     bool      `json:"remote"`
	Server     string    `json:"server,omitempty"`
	LastUpdate time.Time `json:"last_update"`
	Percent    float64   `json:"percent"`
}

// LicenseList represents a list of licenses
type LicenseList struct {
	Licenses []License              `json:"licenses"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
}

// Share represents a SLURM fair share entry
type Share struct {
	Name        string    `json:"name"`
	User        string    `json:"user,omitempty"`
	Account     string    `json:"account,omitempty"`
	Cluster     string    `json:"cluster,omitempty"`
	Partition   string    `json:"partition,omitempty"`
	Shares      int       `json:"shares"`
	RawShares   int       `json:"raw_shares"`
	NormShares  float64   `json:"norm_shares"`
	RawUsage    int       `json:"raw_usage"`
	NormUsage   float64   `json:"norm_usage"`
	EffectUsage float64   `json:"effect_usage"`
	FairShare   float64   `json:"fair_share"`
	LevelFS     float64   `json:"level_fs"`
	Priority    float64   `json:"priority"`
	Level       int       `json:"level"`
	LastUpdate  time.Time `json:"last_update"`
}

// SharesList represents a list of shares
type SharesList struct {
	Shares []Share                `json:"shares"`
	Meta   map[string]interface{} `json:"meta,omitempty"`
}

// GetSharesOptions provides filtering options for shares
type GetSharesOptions struct {
	Users      []string   `json:"users,omitempty"`
	Accounts   []string   `json:"accounts,omitempty"`
	Clusters   []string   `json:"clusters,omitempty"`
	UpdateTime *time.Time `json:"update_time,omitempty"`
}

// Config represents SLURM configuration
type Config struct {
	AccountingStorageType string                 `json:"accounting_storage_type,omitempty"`
	AccountingStorageHost string                 `json:"accounting_storage_host,omitempty"`
	AccountingStoragePort int                    `json:"accounting_storage_port,omitempty"`
	AccountingStorageUser string                 `json:"accounting_storage_user,omitempty"`
	ClusterName           string                 `json:"cluster_name,omitempty"`
	ControlMachine        string                 `json:"control_machine,omitempty"`
	AuthType              string                 `json:"auth_type,omitempty"`
	BackupController      string                 `json:"backup_controller,omitempty"`
	DefaultPartition      string                 `json:"default_partition,omitempty"`
	MaxJobCount           int                    `json:"max_job_count,omitempty"`
	MaxStepCount          int                    `json:"max_step_count,omitempty"`
	MaxTasksPerNode       int                    `json:"max_tasks_per_node,omitempty"`
	PluginDir             string                 `json:"plugin_dir,omitempty"`
	ReturnToService       int                    `json:"return_to_service,omitempty"`
	SlurmUser             string                 `json:"slurm_user,omitempty"`
	SlurmctldLogFile      string                 `json:"slurmctld_log_file,omitempty"`
	SlurmdLogFile         string                 `json:"slurmd_log_file,omitempty"`
	StateSaveLocation     string                 `json:"state_save_location,omitempty"`
	TmpFS                 string                 `json:"tmp_fs,omitempty"`
	UnkillableStepProgram string                 `json:"unkillable_step_program,omitempty"`
	Version               string                 `json:"version,omitempty"`
	Parameters            map[string]interface{} `json:"parameters,omitempty"`
	Nodes                 []ConfigNode           `json:"nodes,omitempty"`
	Partitions            []ConfigPartition      `json:"partitions,omitempty"`
}

// ConfigNode represents a node configuration
type ConfigNode struct {
	NodeName        string   `json:"node_name"`
	NodeAddr        string   `json:"node_addr,omitempty"`
	NodeHostname    string   `json:"node_hostname,omitempty"`
	CPUs            int      `json:"cpus"`
	Boards          int      `json:"boards,omitempty"`
	SocketsPerBoard int      `json:"sockets_per_board,omitempty"`
	CoresPerSocket  int      `json:"cores_per_socket,omitempty"`
	ThreadsPerCore  int      `json:"threads_per_core,omitempty"`
	RealMemory      int      `json:"real_memory"`
	TmpDisk         int      `json:"tmp_disk,omitempty"`
	Weight          int      `json:"weight,omitempty"`
	Features        []string `json:"features,omitempty"`
	Gres            string   `json:"gres,omitempty"`
	State           string   `json:"state,omitempty"`
	Reason          string   `json:"reason,omitempty"`
}

// ConfigPartition represents a partition configuration
type ConfigPartition struct {
	Name          string   `json:"name"`
	Nodes         []string `json:"nodes,omitempty"`
	AllowGroups   []string `json:"allow_groups,omitempty"`
	AllowAccounts []string `json:"allow_accounts,omitempty"`
	AllowQos      []string `json:"allow_qos,omitempty"`
	AllocNodes    []string `json:"alloc_nodes,omitempty"`
	Default       bool     `json:"default"`
	MaxTime       int      `json:"max_time,omitempty"`
	DefaultTime   int      `json:"default_time,omitempty"`
	MaxNodes      int      `json:"max_nodes,omitempty"`
	MinNodes      int      `json:"min_nodes,omitempty"`
	Priority      int      `json:"priority,omitempty"`
	State         string   `json:"state,omitempty"`
	TotalCPUs     int      `json:"total_cpus,omitempty"`
	TotalNodes    int      `json:"total_nodes,omitempty"`
}

// Diagnostics represents SLURM diagnostics information
type Diagnostics struct {
	DataCollected        time.Time              `json:"data_collected"`
	ReqTime              int64                  `json:"req_time"`
	ReqTimeStart         int64                  `json:"req_time_start"`
	ServerThreadCount    int                    `json:"server_thread_count"`
	AgentCount           int                    `json:"agent_count"`
	AgentThreadCount     int                    `json:"agent_thread_count"`
	DBDAgentCount        int                    `json:"dbd_agent_count,omitempty"`
	JobsSubmitted        int                    `json:"jobs_submitted"`
	JobsStarted          int                    `json:"jobs_started"`
	JobsCompleted        int                    `json:"jobs_completed"`
	JobsCanceled         int                    `json:"jobs_canceled"`
	JobsFailed           int                    `json:"jobs_failed"`
	ScheduleCycleMax     int                    `json:"schedule_cycle_max"`
	ScheduleCycleLast    int                    `json:"schedule_cycle_last"`
	ScheduleCycleTotal   int64                  `json:"schedule_cycle_total"`
	ScheduleCycleCounter int                    `json:"schedule_cycle_counter"`
	ScheduleCycleMean    float64                `json:"schedule_cycle_mean"`
	BackfillCycleMax     int                    `json:"backfill_cycle_max"`
	BackfillCycleLast    int                    `json:"backfill_cycle_last"`
	BackfillCycleTotal   int64                  `json:"backfill_cycle_total"`
	BackfillCycleCounter int                    `json:"backfill_cycle_counter"`
	BackfillCycleMean    float64                `json:"backfill_cycle_mean"`
	BfBackfilledJobs     int                    `json:"bf_backfilled_jobs"`
	BfLastBackfilledJobs int                    `json:"bf_last_backfilled_jobs"`
	BfCycleSum           int64                  `json:"bf_cycle_sum"`
	BfQueueLen           int                    `json:"bf_queue_len"`
	BfQueueLenSum        int64                  `json:"bf_queue_len_sum"`
	BfWhenLastCycle      int64                  `json:"bf_when_last_cycle"`
	BfActive             bool                   `json:"bf_active"`
	RPCsByMessageType    map[string]int         `json:"rpcs_by_message_type,omitempty"`
	RPCsByUser           map[string]int         `json:"rpcs_by_user,omitempty"`
	PendingRPCs          int                    `json:"pending_rpcs"`
	Statistics           map[string]interface{} `json:"statistics,omitempty"`
}

// Instance represents a SLURM database instance
type Instance struct {
	Cluster   string    `json:"cluster"`
	ExtraInfo string    `json:"extra,omitempty"`
	Instance  string    `json:"instance"`
	NodeName  string    `json:"node_name"`
	TimeEnd   int64     `json:"time_end"`
	TimeStart int64     `json:"time_start"`
	TRES      string    `json:"tres,omitempty"`
	Created   time.Time `json:"created"`
	Modified  time.Time `json:"modified"`
}

// InstanceList represents a list of instances
type InstanceList struct {
	Instances []Instance             `json:"instances"`
	Meta      map[string]interface{} `json:"meta,omitempty"`
}

// GetInstanceOptions provides filtering options for instances
type GetInstanceOptions struct {
	Cluster   string     `json:"cluster,omitempty"`
	Extra     string     `json:"extra,omitempty"`
	Instance  string     `json:"instance,omitempty"`
	NodeList  []string   `json:"node_list,omitempty"`
	TimeStart *time.Time `json:"time_start,omitempty"`
	TimeEnd   *time.Time `json:"time_end,omitempty"`
}

// GetInstancesOptions provides filtering options for multiple instances
type GetInstancesOptions struct {
	Clusters  []string   `json:"clusters,omitempty"`
	Extra     string     `json:"extra,omitempty"`
	Instances []string   `json:"instances,omitempty"`
	NodeList  []string   `json:"node_list,omitempty"`
	TimeStart *time.Time `json:"time_start,omitempty"`
	TimeEnd   *time.Time `json:"time_end,omitempty"`
}

// TRES represents a Trackable RESources entry
type TRES struct {
	ID          uint64    `json:"id"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Count       int64     `json:"count,omitempty"`
	AllocSecs   int64     `json:"alloc_secs,omitempty"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
	Description string    `json:"description,omitempty"`
}

// TRESList represents a list of TRES
type TRESList struct {
	TRES []TRES                 `json:"tres"`
	Meta map[string]interface{} `json:"meta,omitempty"`
}

// CreateTRESRequest represents a request to create TRES
type CreateTRESRequest struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// ReconfigureResponse represents the response from a reconfigure operation
type ReconfigureResponse struct {
	Status   string                 `json:"status"`
	Message  string                 `json:"message,omitempty"`
	Changes  []string               `json:"changes,omitempty"`
	Warnings []string               `json:"warnings,omitempty"`
	Errors   []string               `json:"errors,omitempty"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
}
