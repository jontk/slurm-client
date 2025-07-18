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
}

// GPUUtilization represents GPU-specific utilization metrics
type GPUUtilization struct {
	DeviceCount        int                            `json:"device_count"`
	Devices            []GPUDeviceUtilization         `json:"devices"`
	OverallUtilization *ResourceUtilization           `json:"overall_utilization"`
	MemoryUtilization  *ResourceUtilization           `json:"memory_utilization"`
	PowerUsage         *ResourceUtilization           `json:"power_usage"`
	Temperature        map[string]float64             `json:"temperature"`
	ComputeMode        string                         `json:"compute_mode"`
	DriverVersion      string                         `json:"driver_version"`
	CUDAVersion        string                         `json:"cuda_version"`
	Metadata           map[string]interface{}         `json:"metadata,omitempty"`
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
	ReadBandwidth    *ResourceUtilization   `json:"read_bandwidth"`
	WriteBandwidth   *ResourceUtilization   `json:"write_bandwidth"`
	ReadIOPS         *ResourceUtilization   `json:"read_iops"`
	WriteIOPS        *ResourceUtilization   `json:"write_iops"`
	TotalBytesRead   int64                  `json:"total_bytes_read"`
	TotalBytesWritten int64                 `json:"total_bytes_written"`
	FileSystems      map[string]IOStats     `json:"file_systems,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
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
	Interfaces        map[string]NetworkInterfaceStats `json:"interfaces"`
	TotalBandwidth    *ResourceUtilization             `json:"total_bandwidth"`
	IngressBandwidth  *ResourceUtilization             `json:"ingress_bandwidth"`
	EgressBandwidth   *ResourceUtilization             `json:"egress_bandwidth"`
	PacketsReceived   int64                            `json:"packets_received"`
	PacketsSent       int64                            `json:"packets_sent"`
	PacketsDropped    int64                            `json:"packets_dropped"`
	Errors            int64                            `json:"errors"`
	ProtocolStats     map[string]int64                 `json:"protocol_stats,omitempty"`
	Metadata          map[string]interface{}           `json:"metadata,omitempty"`
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
	TotalEnergyJoules float64                `json:"total_energy_joules"`
	AveragePowerWatts float64                `json:"average_power_watts"`
	PeakPowerWatts    float64                `json:"peak_power_watts"`
	MinPowerWatts     float64                `json:"min_power_watts"`
	CPUEnergyJoules   float64                `json:"cpu_energy_joules"`
	GPUEnergyJoules   float64                `json:"gpu_energy_joules"`
	MemoryEnergyJoules float64               `json:"memory_energy_joules"`
	CarbonFootprint   float64                `json:"carbon_footprint_kg_co2"`
	PowerSources      map[string]float64     `json:"power_sources,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// JobPerformance represents comprehensive performance metrics for a job
type JobPerformance struct {
	JobID               uint32                      `json:"job_id"`
	JobName             string                      `json:"job_name"`
	StartTime           time.Time                   `json:"start_time"`
	EndTime             *time.Time                  `json:"end_time,omitempty"`
	Status              string                      `json:"status"`
	ExitCode            int                         `json:"exit_code"`
	
	// Resource utilization
	ResourceUtilization *ResourceUtilization        `json:"resource_utilization"`
	JobUtilization      *JobUtilization            `json:"job_utilization"`
	
	// Step-level metrics
	StepMetrics         []JobStepPerformance       `json:"step_metrics,omitempty"`
	
	// Performance trends over time
	PerformanceTrends   *PerformanceTrends         `json:"performance_trends,omitempty"`
	
	// Bottleneck analysis
	Bottlenecks         []PerformanceBottleneck    `json:"bottlenecks,omitempty"`
	
	// Optimization recommendations
	Recommendations     []OptimizationRecommendation `json:"recommendations,omitempty"`
}

// JobStepPerformance represents performance metrics for a job step
type JobStepPerformance struct {
	StepID              uint32                      `json:"step_id"`
	StepName            string                      `json:"step_name"`
	StartTime           time.Time                   `json:"start_time"`
	EndTime             *time.Time                  `json:"end_time,omitempty"`
	Duration            time.Duration               `json:"duration"`
	ExitCode            int                         `json:"exit_code"`
	
	// Resource metrics
	CPUUtilization      float64                     `json:"cpu_utilization"`
	MemoryUtilization   float64                     `json:"memory_utilization"`
	GPUUtilization      float64                     `json:"gpu_utilization,omitempty"`
	IOThroughput        float64                     `json:"io_throughput"`
	NetworkThroughput   float64                     `json:"network_throughput"`
}

// PerformanceTrends represents performance trends over time
type PerformanceTrends struct {
	TimePoints          []time.Time                 `json:"time_points"`
	CPUTrends           []float64                   `json:"cpu_trends"`
	MemoryTrends        []float64                   `json:"memory_trends"`
	GPUTrends           []float64                   `json:"gpu_trends,omitempty"`
	IOTrends            []float64                   `json:"io_trends"`
	NetworkTrends       []float64                   `json:"network_trends"`
	PowerTrends         []float64                   `json:"power_trends,omitempty"`
}

// PerformanceBottleneck represents a detected performance bottleneck
type PerformanceBottleneck struct {
	Type                string                      `json:"type"` // cpu, memory, gpu, io, network
	Severity            string                      `json:"severity"` // low, medium, high, critical
	Description         string                      `json:"description"`
	Impact              float64                     `json:"impact"` // Estimated performance impact percentage
	TimeDetected        time.Time                   `json:"time_detected"`
	Duration            time.Duration               `json:"duration"`
	AffectedNodes       []string                    `json:"affected_nodes,omitempty"`
}

// OptimizationRecommendation represents a performance optimization suggestion
type OptimizationRecommendation struct {
	Type                string                      `json:"type"` // resource_adjustment, configuration, workflow
	Priority            string                      `json:"priority"` // low, medium, high
	Title               string                      `json:"title"`
	Description         string                      `json:"description"`
	ExpectedImprovement float64                     `json:"expected_improvement"` // Percentage improvement
	ResourceChanges     map[string]interface{}      `json:"resource_changes,omitempty"`
	ConfigChanges       map[string]string           `json:"config_changes,omitempty"`
}

// JobLiveMetrics represents real-time performance metrics for a running job
type JobLiveMetrics struct {
	JobID               string                      `json:"job_id"`
	JobName             string                      `json:"job_name"`
	State               string                      `json:"state"`
	RunningTime         time.Duration               `json:"running_time"`
	CollectionTime      time.Time                   `json:"collection_time"`
	
	// Current resource usage
	CPUUsage            *LiveResourceMetric         `json:"cpu_usage"`
	MemoryUsage         *LiveResourceMetric         `json:"memory_usage"`
	GPUUsage            *LiveResourceMetric         `json:"gpu_usage,omitempty"`
	NetworkUsage        *LiveResourceMetric         `json:"network_usage,omitempty"`
	IOUsage             *LiveResourceMetric         `json:"io_usage,omitempty"`
	
	// Process information
	ProcessCount        int                         `json:"process_count"`
	ThreadCount         int                         `json:"thread_count"`
	
	// Node-level metrics
	NodeMetrics         map[string]*NodeLiveMetrics `json:"node_metrics,omitempty"`
	
	// Alerts and warnings
	Alerts              []PerformanceAlert          `json:"alerts,omitempty"`
	
	// Metadata
	Metadata            map[string]interface{}      `json:"metadata,omitempty"`
}

// LiveResourceMetric represents a real-time resource metric
type LiveResourceMetric struct {
	Current             float64                     `json:"current"`
	Average1Min         float64                     `json:"average_1min"`
	Average5Min         float64                     `json:"average_5min"`
	Peak                float64                     `json:"peak"`
	Allocated           float64                     `json:"allocated"`
	UtilizationPercent  float64                     `json:"utilization_percent"`
	Trend               string                      `json:"trend"` // increasing, decreasing, stable
	Unit                string                      `json:"unit"`
}

// NodeLiveMetrics represents real-time metrics for a specific node
type NodeLiveMetrics struct {
	NodeName            string                      `json:"node_name"`
	CPUCores            int                         `json:"cpu_cores"`
	MemoryGB            float64                     `json:"memory_gb"`
	
	// Resource metrics
	CPUUsage            *LiveResourceMetric         `json:"cpu_usage"`
	MemoryUsage         *LiveResourceMetric         `json:"memory_usage"`
	LoadAverage         []float64                   `json:"load_average"` // 1, 5, 15 min
	
	// Temperature and power
	CPUTemperature      float64                     `json:"cpu_temperature_celsius,omitempty"`
	PowerConsumption    float64                     `json:"power_consumption_watts,omitempty"`
	
	// Network and I/O
	NetworkInRate       float64                     `json:"network_in_rate_mbps,omitempty"`
	NetworkOutRate      float64                     `json:"network_out_rate_mbps,omitempty"`
	DiskReadRate        float64                     `json:"disk_read_rate_mbps,omitempty"`
	DiskWriteRate       float64                     `json:"disk_write_rate_mbps,omitempty"`
}

// PerformanceAlert represents a performance-related alert or warning
type PerformanceAlert struct {
	Type                string                      `json:"type"` // warning, critical
	Category            string                      `json:"category"` // cpu, memory, gpu, io, network
	Message             string                      `json:"message"`
	Severity            string                      `json:"severity"` // low, medium, high, critical
	Timestamp           time.Time                   `json:"timestamp"`
	NodeName            string                      `json:"node_name,omitempty"`
	ResourceName        string                      `json:"resource_name,omitempty"`
	CurrentValue        float64                     `json:"current_value,omitempty"`
	ThresholdValue      float64                     `json:"threshold_value,omitempty"`
	RecommendedAction   string                      `json:"recommended_action,omitempty"`
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

// AccountManager manages account operations (v0.0.43+)
type AccountManager interface {
	List(ctx context.Context, opts *ListAccountsOptions) (*AccountList, error)
	Get(ctx context.Context, accountName string) (*Account, error)
	Create(ctx context.Context, account *AccountCreate) (*AccountCreateResponse, error)
	Update(ctx context.Context, accountName string, update *AccountUpdate) error
	Delete(ctx context.Context, accountName string) error
	
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
	Name               string   `json:"name"`
	Description        string   `json:"description,omitempty"`
	Organization       string   `json:"organization,omitempty"`
	CoordinatorUsers   []string `json:"coordinator_users,omitempty"`
	AllowedPartitions  []string `json:"allowed_partitions,omitempty"`
	DefaultPartition   string   `json:"default_partition,omitempty"`
	AllowedQoS         []string `json:"allowed_qos,omitempty"`
	DefaultQoS         string   `json:"default_qos,omitempty"`
	CPULimit           int      `json:"cpu_limit,omitempty"`
	MaxJobs            int      `json:"max_jobs,omitempty"`
	MaxJobsPerUser     int      `json:"max_jobs_per_user,omitempty"`
	MaxNodes           int      `json:"max_nodes,omitempty"`
	MaxWallTime        int      `json:"max_wall_time,omitempty"`
	FairShareTRES      map[string]int `json:"fairshare_tres,omitempty"`
	GrpTRES            map[string]int `json:"grp_tres,omitempty"`
	GrpTRESMinutes     map[string]int `json:"grp_tres_minutes,omitempty"`
	MaxTRES            map[string]int `json:"max_tres,omitempty"`
	MaxTRESPerUser     map[string]int `json:"max_tres_per_user,omitempty"`
	SharesPriority     int      `json:"shares_priority,omitempty"`
	ParentAccount      string   `json:"parent_account,omitempty"`
	ChildAccounts      []string `json:"child_accounts,omitempty"`
	Users              []string `json:"users,omitempty"`
	Flags              []string `json:"flags,omitempty"`
	CreateTime         string   `json:"create_time,omitempty"`
	UpdateTime         string   `json:"update_time,omitempty"`
	
	// Enhanced fields for extended account management
	Quota              *AccountQuota     `json:"quota,omitempty"`
	Usage              *AccountUsage     `json:"usage,omitempty"`
	HierarchyLevel     int               `json:"hierarchy_level,omitempty"`
	HierarchyPath      []string          `json:"hierarchy_path,omitempty"`
	TotalSubAccounts   int               `json:"total_sub_accounts,omitempty"`
	ActiveUserCount    int               `json:"active_user_count,omitempty"`
}

// AccountList represents a list of accounts
type AccountList struct {
	Accounts []Account `json:"accounts"`
	Total    int       `json:"total"`
}

// AccountCreate represents a request to create an account
type AccountCreate struct {
	Name               string   `json:"name"`
	Description        string   `json:"description,omitempty"`
	Organization       string   `json:"organization,omitempty"`
	CoordinatorUsers   []string `json:"coordinator_users,omitempty"`
	AllowedPartitions  []string `json:"allowed_partitions,omitempty"`
	DefaultPartition   string   `json:"default_partition,omitempty"`
	AllowedQoS         []string `json:"allowed_qos,omitempty"`
	DefaultQoS         string   `json:"default_qos,omitempty"`
	CPULimit           int      `json:"cpu_limit,omitempty"`
	MaxJobs            int      `json:"max_jobs,omitempty"`
	MaxJobsPerUser     int      `json:"max_jobs_per_user,omitempty"`
	MaxNodes           int      `json:"max_nodes,omitempty"`
	MaxWallTime        int      `json:"max_wall_time,omitempty"`
	FairShareTRES      map[string]int `json:"fairshare_tres,omitempty"`
	GrpTRES            map[string]int `json:"grp_tres,omitempty"`
	GrpTRESMinutes     map[string]int `json:"grp_tres_minutes,omitempty"`
	MaxTRES            map[string]int `json:"max_tres,omitempty"`
	MaxTRESPerUser     map[string]int `json:"max_tres_per_user,omitempty"`
	SharesPriority     int      `json:"shares_priority,omitempty"`
	ParentAccount      string   `json:"parent_account,omitempty"`
	Flags              []string `json:"flags,omitempty"`
}

// AccountCreateResponse represents the response from account creation
type AccountCreateResponse struct {
	AccountName string `json:"account_name"`
}

// AccountUpdate represents an account update request
type AccountUpdate struct {
	Description        *string  `json:"description,omitempty"`
	Organization       *string  `json:"organization,omitempty"`
	CoordinatorUsers   []string `json:"coordinator_users,omitempty"`
	AllowedPartitions  []string `json:"allowed_partitions,omitempty"`
	DefaultPartition   *string  `json:"default_partition,omitempty"`
	AllowedQoS         []string `json:"allowed_qos,omitempty"`
	DefaultQoS         *string  `json:"default_qos,omitempty"`
	CPULimit           *int     `json:"cpu_limit,omitempty"`
	MaxJobs            *int     `json:"max_jobs,omitempty"`
	MaxJobsPerUser     *int     `json:"max_jobs_per_user,omitempty"`
	MaxNodes           *int     `json:"max_nodes,omitempty"`
	MaxWallTime        *int     `json:"max_wall_time,omitempty"`
	FairShareTRES      map[string]int `json:"fairshare_tres,omitempty"`
	GrpTRES            map[string]int `json:"grp_tres,omitempty"`
	GrpTRESMinutes     map[string]int `json:"grp_tres_minutes,omitempty"`
	MaxTRES            map[string]int `json:"max_tres,omitempty"`
	MaxTRESPerUser     map[string]int `json:"max_tres_per_user,omitempty"`
	SharesPriority     *int     `json:"shares_priority,omitempty"`
	Flags              []string `json:"flags,omitempty"`
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
	CPULimit           int                `json:"cpu_limit,omitempty"`
	CPUUsed            int                `json:"cpu_used,omitempty"`
	MaxJobs            int                `json:"max_jobs,omitempty"`
	JobsUsed           int                `json:"jobs_used,omitempty"`
	MaxJobsPerUser     int                `json:"max_jobs_per_user,omitempty"`
	MaxNodes           int                `json:"max_nodes,omitempty"`
	NodesUsed          int                `json:"nodes_used,omitempty"`
	MaxWallTime        int                `json:"max_wall_time,omitempty"`
	GrpTRES            map[string]int     `json:"grp_tres,omitempty"`
	GrpTRESUsed        map[string]int     `json:"grp_tres_used,omitempty"`
	GrpTRESMinutes     map[string]int     `json:"grp_tres_minutes,omitempty"`
	GrpTRESMinutesUsed map[string]int     `json:"grp_tres_minutes_used,omitempty"`
	MaxTRES            map[string]int     `json:"max_tres,omitempty"`
	MaxTRESUsed        map[string]int     `json:"max_tres_used,omitempty"`
	MaxTRESPerUser     map[string]int     `json:"max_tres_per_user,omitempty"`
	QuotaPeriod        string             `json:"quota_period,omitempty"`
	LastUpdated        time.Time          `json:"last_updated,omitempty"`
}

// AccountUsage represents usage statistics for an account
type AccountUsage struct {
	AccountName        string             `json:"account_name"`
	CPUHours           float64            `json:"cpu_hours"`
	JobCount           int                `json:"job_count"`
	JobsCompleted      int                `json:"jobs_completed"`
	JobsFailed         int                `json:"jobs_failed"`
	JobsCanceled       int                `json:"jobs_canceled"`
	TotalWallTime      float64            `json:"total_wall_time"`
	AverageWallTime    float64            `json:"average_wall_time"`
	TRESUsage          map[string]float64 `json:"tres_usage,omitempty"`
	UserCount          int                `json:"user_count"`
	ActiveUsers        []string           `json:"active_users,omitempty"`
	Period             string             `json:"period,omitempty"`
	StartTime          time.Time          `json:"start_time"`
	EndTime            time.Time          `json:"end_time"`
	EfficiencyRatio    float64            `json:"efficiency_ratio,omitempty"`
}

// AccountHierarchy represents the hierarchical structure of accounts
type AccountHierarchy struct {
	Account         *Account               `json:"account"`
	ParentAccount   *AccountHierarchy      `json:"parent_account,omitempty"`
	ChildAccounts   []*AccountHierarchy    `json:"child_accounts,omitempty"`
	Level           int                    `json:"level"`
	Path            []string               `json:"path"`
	TotalUsers      int                    `json:"total_users"`
	TotalSubAccounts int                   `json:"total_sub_accounts"`
	AggregateQuota  *AccountQuota          `json:"aggregate_quota,omitempty"`
	AggregateUsage  *AccountUsage          `json:"aggregate_usage,omitempty"`
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
	CPUThreshold     float64 `json:"cpu_threshold,omitempty"`      // Alert if CPU usage > threshold (0-100)
	MemoryThreshold  float64 `json:"memory_threshold,omitempty"`   // Alert if memory usage > threshold (0-100)
	GPUThreshold     float64 `json:"gpu_threshold,omitempty"`      // Alert if GPU usage > threshold (0-100)
	// Stop conditions
	StopOnCompletion bool `json:"stop_on_completion"` // Stop watching when job completes
	MaxDuration      time.Duration `json:"max_duration,omitempty"` // Maximum time to watch
}

// JobMetricsEvent represents a job metrics update event
type JobMetricsEvent struct {
	Type           string              `json:"type"` // "update", "alert", "error", "complete"
	JobID          string              `json:"job_id"`
	Timestamp      time.Time           `json:"timestamp"`
	Metrics        *JobLiveMetrics     `json:"metrics,omitempty"`
	Alert          *PerformanceAlert   `json:"alert,omitempty"`
	Error          error               `json:"error,omitempty"`
	StateChange    *JobStateChange     `json:"state_change,omitempty"`
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
	TimeWindow    time.Duration `json:"time_window,omitempty"`     // Default: job duration or 1 hour
	// Number of data points to collect
	DataPoints    int           `json:"data_points,omitempty"`     // Default: 24
	// Resources to include in trends
	IncludeCPU    bool          `json:"include_cpu"`
	IncludeMemory bool          `json:"include_memory"`
	IncludeGPU    bool          `json:"include_gpu"`
	IncludeIO     bool          `json:"include_io"`
	IncludeNetwork bool         `json:"include_network"`
	IncludeEnergy bool          `json:"include_energy"`
	// Aggregation method for data points
	Aggregation   string        `json:"aggregation,omitempty"`     // "avg", "max", "min" (default: "avg")
	// Include anomaly detection
	DetectAnomalies bool        `json:"detect_anomalies"`
}

// JobResourceTrends represents resource usage trends over time
type JobResourceTrends struct {
	JobID          string                       `json:"job_id"`
	JobName        string                       `json:"job_name"`
	StartTime      time.Time                    `json:"start_time"`
	EndTime        *time.Time                   `json:"end_time,omitempty"`
	TimeWindow     time.Duration                `json:"time_window"`
	DataPoints     int                          `json:"data_points"`
	
	// Time series data
	TimePoints     []time.Time                  `json:"time_points"`
	
	// Resource trends
	CPUTrends      *ResourceTimeSeries          `json:"cpu_trends,omitempty"`
	MemoryTrends   *ResourceTimeSeries          `json:"memory_trends,omitempty"`
	GPUTrends      *ResourceTimeSeries          `json:"gpu_trends,omitempty"`
	IOTrends       *IOTimeSeries                `json:"io_trends,omitempty"`
	NetworkTrends  *NetworkTimeSeries           `json:"network_trends,omitempty"`
	EnergyTrends   *EnergyTimeSeries            `json:"energy_trends,omitempty"`
	
	// Anomalies detected
	Anomalies      []ResourceAnomaly            `json:"anomalies,omitempty"`
	
	// Summary statistics
	Summary        *TrendsSummary               `json:"summary"`
	
	// Metadata
	Metadata       map[string]interface{}       `json:"metadata,omitempty"`
}

// ResourceTimeSeries represents time series data for a resource
type ResourceTimeSeries struct {
	Values         []float64                    `json:"values"`
	Unit           string                       `json:"unit"`
	Average        float64                      `json:"average"`
	Min            float64                      `json:"min"`
	Max            float64                      `json:"max"`
	StdDev         float64                      `json:"std_dev"`
	Trend          string                       `json:"trend"` // "increasing", "decreasing", "stable", "fluctuating"
	TrendSlope     float64                      `json:"trend_slope"`
}

// IOTimeSeries represents I/O time series data
type IOTimeSeries struct {
	ReadBandwidth  *ResourceTimeSeries          `json:"read_bandwidth,omitempty"`
	WriteBandwidth *ResourceTimeSeries          `json:"write_bandwidth,omitempty"`
	ReadIOPS       *ResourceTimeSeries          `json:"read_iops,omitempty"`
	WriteIOPS      *ResourceTimeSeries          `json:"write_iops,omitempty"`
}

// NetworkTimeSeries represents network time series data
type NetworkTimeSeries struct {
	IngressBandwidth  *ResourceTimeSeries       `json:"ingress_bandwidth,omitempty"`
	EgressBandwidth   *ResourceTimeSeries       `json:"egress_bandwidth,omitempty"`
	PacketRate        *ResourceTimeSeries       `json:"packet_rate,omitempty"`
}

// EnergyTimeSeries represents energy usage time series data
type EnergyTimeSeries struct {
	PowerUsage        *ResourceTimeSeries       `json:"power_usage,omitempty"`
	EnergyConsumption *ResourceTimeSeries       `json:"energy_consumption,omitempty"`
	CarbonEmissions   *ResourceTimeSeries       `json:"carbon_emissions,omitempty"`
}

// ResourceAnomaly represents an anomaly detected in resource usage
type ResourceAnomaly struct {
	Timestamp      time.Time                    `json:"timestamp"`
	Resource       string                       `json:"resource"` // "cpu", "memory", "gpu", etc.
	Type           string                       `json:"type"`     // "spike", "drop", "pattern_change"
	Severity       string                       `json:"severity"` // "low", "medium", "high"
	Value          float64                      `json:"value"`
	ExpectedValue  float64                      `json:"expected_value"`
	Deviation      float64                      `json:"deviation_percent"`
	Description    string                       `json:"description"`
}

// TrendsSummary provides summary statistics for resource trends
type TrendsSummary struct {
	OverallTrend      string                   `json:"overall_trend"`
	ResourceEfficiency float64                 `json:"resource_efficiency"`
	StabilityScore    float64                  `json:"stability_score"` // 0-100, higher is more stable
	VariabilityIndex  float64                  `json:"variability_index"`
	PeakUtilization   map[string]float64       `json:"peak_utilization"`
	AverageUtilization map[string]float64      `json:"average_utilization"`
	ResourceBalance   string                   `json:"resource_balance"` // "balanced", "cpu_heavy", "memory_heavy", etc.
}

// UserManager provides user-related operations
type UserManager interface {
	// Core user operations
	List(ctx context.Context, opts *ListUsersOptions) (*UserList, error)
	Get(ctx context.Context, userName string) (*User, error)
	
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

// User represents a SLURM user
type User struct {
	Name               string                 `json:"name"`
	UID                int                    `json:"uid"`
	DefaultAccount     string                 `json:"default_account"`
	DefaultWCKey       string                 `json:"default_wckey,omitempty"`
	AdminLevel         string                 `json:"admin_level"`
	CoordinatorAccounts []string              `json:"coordinator_accounts,omitempty"`
	Accounts           []UserAccount          `json:"accounts,omitempty"`
	Quotas             *UserQuota             `json:"quotas,omitempty"`
	FairShare          *UserFairShare         `json:"fair_share,omitempty"`
	Associations       []UserAssociation      `json:"associations,omitempty"`
	Created            time.Time              `json:"created"`
	Modified           time.Time              `json:"modified"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

// UserAccount represents a user's association with an account
type UserAccount struct {
	AccountName  string            `json:"account_name"`
	Partition    string            `json:"partition,omitempty"`
	QoS          string            `json:"qos,omitempty"`
	DefaultQoS   string            `json:"default_qos,omitempty"`
	MaxJobs      int               `json:"max_jobs,omitempty"`
	MaxSubmitJobs int              `json:"max_submit_jobs,omitempty"`
	MaxWallTime  int               `json:"max_wall_time,omitempty"`
	Priority     int               `json:"priority,omitempty"`
	GraceTime    int               `json:"grace_time,omitempty"`
	TRES         map[string]int    `json:"tres,omitempty"`
	MaxTRES      map[string]int    `json:"max_tres,omitempty"`
	MinTRES      map[string]int    `json:"min_tres,omitempty"`
	IsDefault    bool              `json:"is_default"`
	IsActive     bool              `json:"is_active"`
	Flags        []string          `json:"flags,omitempty"`
	Created      time.Time         `json:"created"`
	Modified     time.Time         `json:"modified"`
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
	UserName          string               `json:"user_name"`
	DefaultAccount    string               `json:"default_account"`
	MaxJobs           int                  `json:"max_jobs"`
	MaxSubmitJobs     int                  `json:"max_submit_jobs"`
	MaxWallTime       int                  `json:"max_wall_time"`
	MaxCPUs           int                  `json:"max_cpus"`
	MaxNodes          int                  `json:"max_nodes"`
	MaxMemory         int                  `json:"max_memory"`
	TRESLimits        map[string]int       `json:"tres_limits,omitempty"`
	AccountQuotas     map[string]*UserAccountQuota `json:"account_quotas,omitempty"`
	QoSLimits         map[string]*QoSLimits `json:"qos_limits,omitempty"`
	GraceTime         int                  `json:"grace_time,omitempty"`
	CurrentUsage      *UserUsage           `json:"current_usage,omitempty"`
	IsActive          bool                 `json:"is_active"`
	Enforcement       string               `json:"enforcement"`
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
	UserName        string                    `json:"user_name"`
	RunningJobs     int                       `json:"running_jobs"`
	PendingJobs     int                       `json:"pending_jobs"`
	UsedCPUHours    float64                   `json:"used_cpu_hours"`
	UsedGPUHours    float64                   `json:"used_gpu_hours,omitempty"`
	UsedWallTime    int                       `json:"used_wall_time"`
	TRESUsage       map[string]float64        `json:"tres_usage,omitempty"`
	AccountUsage    map[string]*AccountUsageStats `json:"account_usage,omitempty"`
	Efficiency      float64                   `json:"efficiency"`
	LastJobTime     time.Time                 `json:"last_job_time"`
	PeriodStart     time.Time                 `json:"period_start"`
	PeriodEnd       time.Time                 `json:"period_end"`
}

// AccountUsageStats represents usage statistics for a user within an account
type AccountUsageStats struct {
	AccountName    string            `json:"account_name"`
	JobCount       int               `json:"job_count"`
	CPUHours       float64           `json:"cpu_hours"`
	WallHours      float64           `json:"wall_hours"`
	TRESUsage      map[string]float64 `json:"tres_usage,omitempty"`
	AverageQueueTime float64         `json:"average_queue_time"`
	AverageRunTime   float64         `json:"average_run_time"`
	Efficiency     float64           `json:"efficiency"`
}

// UserFairShare represents fair-share information for a user
type UserFairShare struct {
	UserName         string                      `json:"user_name"`
	Account          string                      `json:"account"`
	Cluster          string                      `json:"cluster"`
	Partition        string                      `json:"partition,omitempty"`
	FairShareFactor  float64                     `json:"fair_share_factor"`
	NormalizedShares float64                     `json:"normalized_shares"`
	EffectiveUsage   float64                     `json:"effective_usage"`
	FairShareTree    *FairShareNode              `json:"fair_share_tree,omitempty"`
	PriorityFactors  *JobPriorityFactors         `json:"priority_factors,omitempty"`
	RawShares        int                         `json:"raw_shares"`
	NormalizedUsage  float64                     `json:"normalized_usage"`
	Level            int                         `json:"level"`
	LastDecay        time.Time                   `json:"last_decay"`
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
	Age       int     `json:"age"`
	FairShare int     `json:"fair_share"`
	JobSize   int     `json:"job_size"`
	Partition int     `json:"partition"`
	QoS       int     `json:"qos"`
	TRES      int     `json:"tres"`
	Site      int     `json:"site"`
	Nice      int     `json:"nice"`
	Assoc     int     `json:"assoc"`
	Total     int     `json:"total"`
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
	AccountName      string                `json:"account_name"`
	Cluster          string                `json:"cluster"`
	Parent           string                `json:"parent,omitempty"`
	Shares           int                   `json:"shares"`
	RawShares        int                   `json:"raw_shares"`
	NormalizedShares float64               `json:"normalized_shares"`
	Usage            float64               `json:"usage"`
	EffectiveUsage   float64               `json:"effective_usage"`
	FairShareFactor  float64               `json:"fair_share_factor"`
	Level            int                   `json:"level"`
	LevelShares      int                   `json:"level_shares"`
	UserCount        int                   `json:"user_count"`
	ActiveUsers      int                   `json:"active_users"`
	JobCount         int                   `json:"job_count"`
	Children         []*AccountFairShare   `json:"children,omitempty"`
	Users            []*UserFairShare      `json:"users,omitempty"`
	LastDecay        time.Time             `json:"last_decay"`
	Created          time.Time             `json:"created"`
	Modified         time.Time             `json:"modified"`
}

// FairShareHierarchy represents the complete fair-share tree structure
type FairShareHierarchy struct {
	Cluster       string            `json:"cluster"`
	RootAccount   string            `json:"root_account"`
	Tree          *FairShareNode    `json:"tree"`
	TotalShares   int               `json:"total_shares"`
	TotalUsage    float64           `json:"total_usage"`
	LastUpdate    time.Time         `json:"last_update"`
	DecayHalfLife int               `json:"decay_half_life"`
	UsageWindow   int               `json:"usage_window"`
	Algorithm     string            `json:"algorithm"`
	Accounts      []*AccountFairShare `json:"accounts,omitempty"`
	Users         []*UserFairShare  `json:"users,omitempty"`
}

// UserList represents a list of users
type UserList struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
}

// ListUsersOptions provides filtering options for listing users
type ListUsersOptions struct {
	// Basic filtering
	Names     []string `json:"names,omitempty"`
	Accounts  []string `json:"accounts,omitempty"`
	Clusters  []string `json:"clusters,omitempty"`
	AdminLevels []string `json:"admin_levels,omitempty"`
	
	// State filtering
	ActiveOnly    bool `json:"active_only,omitempty"`
	CoordinatorsOnly bool `json:"coordinators_only,omitempty"`
	
	// Include additional data
	WithAccounts   bool `json:"with_accounts,omitempty"`
	WithQuotas     bool `json:"with_quotas,omitempty"`
	WithFairShare  bool `json:"with_fair_share,omitempty"`
	WithAssociations bool `json:"with_associations,omitempty"`
	WithUsage      bool `json:"with_usage,omitempty"`
	
	// Pagination
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
	
	// Sorting
	SortBy    string `json:"sort_by,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
}

// AssociationUsage represents usage data for an association
type AssociationUsage struct {
	UsedCPUHours    float64            `json:"used_cpu_hours"`
	UsedGPUHours    float64            `json:"used_gpu_hours,omitempty"`
	UsedWallTime    int                `json:"used_wall_time"`
	TRESUsage       map[string]float64 `json:"tres_usage,omitempty"`
	JobCount        int                `json:"job_count"`
	RunningJobs     int                `json:"running_jobs"`
	PendingJobs     int                `json:"pending_jobs"`
	CompletedJobs   int                `json:"completed_jobs"`
	FailedJobs      int                `json:"failed_jobs"`
	CancelledJobs   int                `json:"cancelled_jobs"`
	Efficiency      float64            `json:"efficiency"`
	PeriodStart     time.Time          `json:"period_start"`
	PeriodEnd       time.Time          `json:"period_end"`
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
	UserName        string                 `json:"user_name"`
	AccountName     string                 `json:"account_name"`
	HasAccess       bool                   `json:"has_access"`
	AccessLevel     string                 `json:"access_level"`
	Permissions     []string               `json:"permissions"`
	Restrictions    []string               `json:"restrictions,omitempty"`
	Reason          string                 `json:"reason,omitempty"`
	ValidFrom       time.Time              `json:"valid_from"`
	ValidUntil      *time.Time             `json:"valid_until,omitempty"`
	Association     *UserAccountAssociation `json:"association,omitempty"`
	QuotaLimits     *UserAccountQuota      `json:"quota_limits,omitempty"`
	CurrentUsage    *AccountUsageStats     `json:"current_usage,omitempty"`
	ValidationTime  time.Time              `json:"validation_time"`
}

// ListAccountUsersOptions provides filtering options for listing account users
type ListAccountUsersOptions struct {
	// Basic filtering
	Roles         []string `json:"roles,omitempty"`
	Permissions   []string `json:"permissions,omitempty"`
	ActiveOnly    bool     `json:"active_only,omitempty"`
	CoordinatorsOnly bool  `json:"coordinators_only,omitempty"`
	
	// State filtering
	Partitions    []string `json:"partitions,omitempty"`
	QoS           []string `json:"qos,omitempty"`
	
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
	Accounts      []string `json:"accounts,omitempty"`
	Clusters      []string `json:"clusters,omitempty"`
	Partitions    []string `json:"partitions,omitempty"`
	Roles         []string `json:"roles,omitempty"`
	Permissions   []string `json:"permissions,omitempty"`
	
	// State filtering
	ActiveOnly    bool `json:"active_only,omitempty"`
	DefaultOnly   bool `json:"default_only,omitempty"`
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
