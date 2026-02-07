// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Package analytics provides performance analytics types for HPC job analysis.
// These types are SDK-specific and provide value-added analytics not in the SLURM REST API.
package api

import "time"

// ============================================================================
// Core Analytics Types
// ============================================================================

// JobUtilization represents detailed resource utilization metrics for a job.
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

// ResourceUtilization represents metrics for a single resource type.
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

// ============================================================================
// GPU Analytics
// ============================================================================

// GPUUtilization represents GPU usage metrics.
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

// GPUDeviceUtilization represents per-GPU device metrics.
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

// GPUProcess represents a process using GPU resources.
type GPUProcess struct {
	PID         int    `json:"pid"`
	ProcessName string `json:"process_name"`
	MemoryUsed  int64  `json:"memory_used_mb"`
}

// ============================================================================
// I/O Analytics
// ============================================================================

// IOUtilization represents I/O usage metrics.
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

// IOStats represents I/O statistics for a filesystem.
type IOStats struct {
	MountPoint      string  `json:"mount_point"`
	BytesRead       int64   `json:"bytes_read"`
	BytesWritten    int64   `json:"bytes_written"`
	ReadsCompleted  int64   `json:"reads_completed"`
	WritesCompleted int64   `json:"writes_completed"`
	AvgReadLatency  float64 `json:"avg_read_latency_ms"`
	AvgWriteLatency float64 `json:"avg_write_latency_ms"`
}

// ============================================================================
// Network Analytics
// ============================================================================

// NetworkUtilization represents network usage metrics.
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

// NetworkInterfaceStats represents per-interface network statistics.
type NetworkInterfaceStats struct {
	InterfaceName   string  `json:"interface_name"`
	BytesReceived   int64   `json:"bytes_received"`
	BytesSent       int64   `json:"bytes_sent"`
	PacketsReceived int64   `json:"packets_received"`
	PacketsSent     int64   `json:"packets_sent"`
	BandwidthMbps   float64 `json:"bandwidth_mbps"`
	Utilization     float64 `json:"utilization_percentage"`
}

// ============================================================================
// Energy Analytics
// ============================================================================

// EnergyUsage represents energy consumption metrics.
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

// ============================================================================
// Performance Analytics
// ============================================================================

// JobPerformance represents comprehensive job performance metrics.
type JobPerformance struct {
	JobID     uint32     `json:"job_id"`
	JobName   string     `json:"job_name"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Status    string     `json:"status"`
	ExitCode  int        `json:"exit_code"`

	ResourceUtilization *ResourceUtilization `json:"resource_utilization"`
	JobUtilization      *JobUtilization      `json:"job_utilization"`
	StepMetrics         []JobStepPerformance `json:"step_metrics,omitempty"`
	PerformanceTrends   *PerformanceTrends   `json:"performance_trends,omitempty"`

	Bottlenecks     []PerformanceBottleneck      `json:"bottlenecks,omitempty"`
	Recommendations []OptimizationRecommendation `json:"recommendations,omitempty"`
}

// JobStepPerformance represents step-level performance metrics.
type JobStepPerformance struct {
	StepID    uint32        `json:"step_id"`
	StepName  string        `json:"step_name"`
	StartTime time.Time     `json:"start_time"`
	EndTime   *time.Time    `json:"end_time,omitempty"`
	Duration  time.Duration `json:"duration"`
	ExitCode  int           `json:"exit_code"`

	CPUUtilization    float64 `json:"cpu_utilization"`
	MemoryUtilization float64 `json:"memory_utilization"`
	GPUUtilization    float64 `json:"gpu_utilization,omitempty"`
	IOThroughput      float64 `json:"io_throughput"`
	NetworkThroughput float64 `json:"network_throughput"`
}

// PerformanceBottleneck identifies a resource bottleneck.
type PerformanceBottleneck struct {
	Type          string        `json:"type"`
	Resource      string        `json:"resource"`
	Severity      string        `json:"severity"`
	Description   string        `json:"description"`
	Impact        string        `json:"impact"`
	TimeDetected  time.Time     `json:"time_detected"`
	Duration      time.Duration `json:"duration"`
	AffectedNodes []string      `json:"affected_nodes,omitempty"`
}

// OptimizationRecommendation is defined in efficiency.go

// ============================================================================
// Live Metrics
// ============================================================================

// JobLiveMetrics represents real-time job metrics.
type JobLiveMetrics struct {
	JobID          string        `json:"job_id"`
	JobName        string        `json:"job_name"`
	State          string        `json:"state"`
	RunningTime    time.Duration `json:"running_time"`
	CollectionTime time.Time     `json:"collection_time"`

	CPUUsage     *LiveResourceMetric `json:"cpu_usage"`
	MemoryUsage  *LiveResourceMetric `json:"memory_usage"`
	GPUUsage     *LiveResourceMetric `json:"gpu_usage,omitempty"`
	NetworkUsage *LiveResourceMetric `json:"network_usage,omitempty"`
	IOUsage      *LiveResourceMetric `json:"io_usage,omitempty"`

	ProcessCount int `json:"process_count"`
	ThreadCount  int `json:"thread_count"`

	NodeMetrics map[string]*NodeLiveMetrics `json:"node_metrics,omitempty"`
	Alerts      []PerformanceAlert          `json:"alerts,omitempty"`
	Metadata    map[string]interface{}      `json:"metadata,omitempty"`
}

// LiveResourceMetric represents a real-time resource metric.
type LiveResourceMetric struct {
	Current            float64 `json:"current"`
	Average1Min        float64 `json:"average_1min"`
	Average5Min        float64 `json:"average_5min"`
	Peak               float64 `json:"peak"`
	Allocated          float64 `json:"allocated"`
	UtilizationPercent float64 `json:"utilization_percent"`
	Trend              string  `json:"trend"`
	Unit               string  `json:"unit"`
}

// NodeLiveMetrics represents per-node live metrics.
type NodeLiveMetrics struct {
	NodeName string  `json:"node_name"`
	CPUCores int     `json:"cpu_cores"`
	MemoryGB float64 `json:"memory_gb"`

	CPUUsage    *LiveResourceMetric `json:"cpu_usage"`
	MemoryUsage *LiveResourceMetric `json:"memory_usage"`
	LoadAverage []float64           `json:"load_average"`

	CPUTemperature   float64 `json:"cpu_temperature_celsius,omitempty"`
	PowerConsumption float64 `json:"power_consumption_watts,omitempty"`

	NetworkInRate  float64 `json:"network_in_rate_mbps,omitempty"`
	NetworkOutRate float64 `json:"network_out_rate_mbps,omitempty"`
	DiskReadRate   float64 `json:"disk_read_rate_mbps,omitempty"`
	DiskWriteRate  float64 `json:"disk_write_rate_mbps,omitempty"`
}

// PerformanceAlert represents a performance alert.
type PerformanceAlert struct {
	Type              string    `json:"type"`
	Category          string    `json:"category"`
	Message           string    `json:"message"`
	Severity          string    `json:"severity"`
	Timestamp         time.Time `json:"timestamp"`
	NodeName          string    `json:"node_name,omitempty"`
	ResourceName      string    `json:"resource_name,omitempty"`
	CurrentValue      float64   `json:"current_value,omitempty"`
	ThresholdValue    float64   `json:"threshold_value,omitempty"`
	RecommendedAction string    `json:"recommended_action,omitempty"`
}

// ============================================================================
// Watch Options and Events
// ============================================================================

// WatchMetricsOptions configures metric watching.
type WatchMetricsOptions struct {
	UpdateInterval     time.Duration `json:"update_interval,omitempty"`
	IncludeCPU         bool          `json:"include_cpu"`
	IncludeMemory      bool          `json:"include_memory"`
	IncludeGPU         bool          `json:"include_gpu"`
	IncludeNetwork     bool          `json:"include_network"`
	IncludeIO          bool          `json:"include_io"`
	IncludeEnergy      bool          `json:"include_energy"`
	IncludeNodeMetrics bool          `json:"include_node_metrics"`
	SpecificNodes      []string      `json:"specific_nodes,omitempty"`
	CPUThreshold       float64       `json:"cpu_threshold,omitempty"`
	MemoryThreshold    float64       `json:"memory_threshold,omitempty"`
	GPUThreshold       float64       `json:"gpu_threshold,omitempty"`
	StopOnCompletion   bool          `json:"stop_on_completion"`
	MaxDuration        time.Duration `json:"max_duration,omitempty"`
}

// JobMetricsEvent represents a metrics update event.
type JobMetricsEvent struct {
	Type        string            `json:"type"`
	JobID       string            `json:"job_id"`
	Timestamp   time.Time         `json:"timestamp"`
	Metrics     *JobLiveMetrics   `json:"metrics,omitempty"`
	Alert       *PerformanceAlert `json:"alert,omitempty"`
	Error       error             `json:"error,omitempty"`
	StateChange *JobStateChange   `json:"state_change,omitempty"`
}

// JobStateChange represents a job state transition.
type JobStateChange struct {
	OldState string `json:"old_state"`
	NewState string `json:"new_state"`
	Reason   string `json:"reason,omitempty"`
}

// ============================================================================
// Resource Trends
// ============================================================================

// ResourceTrendsOptions configures trend analysis.
type ResourceTrendsOptions struct {
	TimeWindow      time.Duration `json:"time_window,omitempty"`
	DataPoints      int           `json:"data_points,omitempty"`
	IncludeCPU      bool          `json:"include_cpu"`
	IncludeMemory   bool          `json:"include_memory"`
	IncludeGPU      bool          `json:"include_gpu"`
	IncludeIO       bool          `json:"include_io"`
	IncludeNetwork  bool          `json:"include_network"`
	IncludeEnergy   bool          `json:"include_energy"`
	Aggregation     string        `json:"aggregation,omitempty"`
	DetectAnomalies bool          `json:"detect_anomalies"`
}

// JobResourceTrends represents resource usage trends for a job.
type JobResourceTrends struct {
	JobID      string        `json:"job_id"`
	JobName    string        `json:"job_name"`
	StartTime  time.Time     `json:"start_time"`
	EndTime    *time.Time    `json:"end_time,omitempty"`
	TimeWindow time.Duration `json:"time_window"`
	DataPoints int           `json:"data_points"`

	TimePoints    []time.Time         `json:"time_points"`
	CPUTrends     *ResourceTimeSeries `json:"cpu_trends,omitempty"`
	MemoryTrends  *ResourceTimeSeries `json:"memory_trends,omitempty"`
	GPUTrends     *ResourceTimeSeries `json:"gpu_trends,omitempty"`
	IOTrends      *IOTimeSeries       `json:"io_trends,omitempty"`
	NetworkTrends *NetworkTimeSeries  `json:"network_trends,omitempty"`
	EnergyTrends  *EnergyTimeSeries   `json:"energy_trends,omitempty"`

	Anomalies []ResourceAnomaly      `json:"anomalies,omitempty"`
	Summary   *TrendsSummary         `json:"summary"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ResourceTimeSeries represents a time series of resource values.
type ResourceTimeSeries struct {
	Values     []float64 `json:"values"`
	Unit       string    `json:"unit"`
	Average    float64   `json:"average"`
	Min        float64   `json:"min"`
	Max        float64   `json:"max"`
	StdDev     float64   `json:"std_dev"`
	Trend      string    `json:"trend"`
	TrendSlope float64   `json:"trend_slope"`
}

// IOTimeSeries represents I/O time series data.
type IOTimeSeries struct {
	ReadBandwidth  *ResourceTimeSeries `json:"read_bandwidth,omitempty"`
	WriteBandwidth *ResourceTimeSeries `json:"write_bandwidth,omitempty"`
	ReadIOPS       *ResourceTimeSeries `json:"read_iops,omitempty"`
	WriteIOPS      *ResourceTimeSeries `json:"write_iops,omitempty"`
}

// NetworkTimeSeries represents network time series data.
type NetworkTimeSeries struct {
	IngressBandwidth *ResourceTimeSeries `json:"ingress_bandwidth,omitempty"`
	EgressBandwidth  *ResourceTimeSeries `json:"egress_bandwidth,omitempty"`
	PacketRate       *ResourceTimeSeries `json:"packet_rate,omitempty"`
}

// EnergyTimeSeries represents energy time series data.
type EnergyTimeSeries struct {
	PowerUsage        *ResourceTimeSeries `json:"power_usage,omitempty"`
	EnergyConsumption *ResourceTimeSeries `json:"energy_consumption,omitempty"`
	CarbonEmissions   *ResourceTimeSeries `json:"carbon_emissions,omitempty"`
}

// ResourceAnomaly represents a detected anomaly.
type ResourceAnomaly struct {
	Timestamp     time.Time `json:"timestamp"`
	Resource      string    `json:"resource"`
	Type          string    `json:"type"`
	Severity      string    `json:"severity"`
	Value         float64   `json:"value"`
	ExpectedValue float64   `json:"expected_value"`
	Deviation     float64   `json:"deviation_percent"`
	Description   string    `json:"description"`
}

// TrendsSummary summarizes resource trends.
type TrendsSummary struct {
	OverallTrend       string             `json:"overall_trend"`
	ResourceEfficiency float64            `json:"resource_efficiency"`
	StabilityScore     float64            `json:"stability_score"`
	VariabilityIndex   float64            `json:"variability_index"`
	PeakUtilization    map[string]float64 `json:"peak_utilization"`
	AverageUtilization map[string]float64 `json:"average_utilization"`
	ResourceBalance    string             `json:"resource_balance"`
}

// ============================================================================
// Job Step Analytics
// ============================================================================

// JobStepDetails provides detailed information about a job step.
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

	CPUAllocation    int      `json:"cpu_allocation"`
	MemoryAllocation int64    `json:"memory_allocation_bytes"`
	GPUAllocation    int      `json:"gpu_allocation,omitempty"`
	NodeList         []string `json:"node_list"`
	TaskCount        int      `json:"task_count"`
	TasksPerNode     int      `json:"tasks_per_node,omitempty"`

	TaskDistribution map[string]int    `json:"task_distribution,omitempty"`
	Command          string            `json:"command"`
	CommandLine      string            `json:"command_line"`
	WorkingDir       string            `json:"working_dir"`
	Environment      map[string]string `json:"environment,omitempty"`

	CPUTime    time.Duration `json:"cpu_time"`
	SystemTime time.Duration `json:"system_time"`
	UserTime   time.Duration `json:"user_time"`
	MaxRSS     int64         `json:"max_rss_bytes"`
	MaxVMSize  int64         `json:"max_vmsize_bytes"`
	AverageRSS int64         `json:"average_rss_bytes"`

	TotalReadBytes  int64 `json:"total_read_bytes"`
	TotalWriteBytes int64 `json:"total_write_bytes"`
	ReadOperations  int64 `json:"read_operations"`
	WriteOperations int64 `json:"write_operations"`

	NetworkBytesReceived int64 `json:"network_bytes_received"`
	NetworkBytesSent     int64 `json:"network_bytes_sent"`

	EnergyConsumed   float64 `json:"energy_consumed_joules"`
	AveragePowerDraw float64 `json:"average_power_draw_watts"`

	Tasks []StepTaskInfo `json:"tasks,omitempty"`

	StepType        string                 `json:"step_type"`
	Priority        int                    `json:"priority"`
	AccountingGroup string                 `json:"accounting_group"`
	QOSLevel        string                 `json:"qos_level"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// JobStepUtilization represents step-level utilization metrics.
type JobStepUtilization struct {
	StepID    string        `json:"step_id"`
	StepName  string        `json:"step_name"`
	JobID     string        `json:"job_id"`
	JobName   string        `json:"job_name"`
	StartTime *time.Time    `json:"start_time,omitempty"`
	EndTime   *time.Time    `json:"end_time,omitempty"`
	Duration  time.Duration `json:"duration"`

	CPUUtilization     *ResourceUtilization `json:"cpu_utilization,omitempty"`
	MemoryUtilization  *ResourceUtilization `json:"memory_utilization,omitempty"`
	GPUUtilization     *GPUUtilization      `json:"gpu_utilization,omitempty"`
	IOUtilization      *IOUtilization       `json:"io_utilization,omitempty"`
	NetworkUtilization *NetworkUtilization  `json:"network_utilization,omitempty"`
	EnergyUtilization  *ResourceUtilization `json:"energy_utilization,omitempty"`

	TaskUtilizations   []TaskUtilization       `json:"task_utilizations,omitempty"`
	PerformanceMetrics *StepPerformanceMetrics `json:"performance_metrics,omitempty"`
	Metadata           map[string]interface{}  `json:"metadata,omitempty"`
}

// TaskUtilization represents per-task utilization.
type TaskUtilization struct {
	TaskID            int     `json:"task_id"`
	NodeName          string  `json:"node_name"`
	CPUUtilization    float64 `json:"cpu_utilization"`
	MemoryUtilization float64 `json:"memory_utilization"`
	State             string  `json:"state"`
	ExitCode          int     `json:"exit_code"`
}

// StepTaskInfo provides information about a task within a step.
type StepTaskInfo struct {
	TaskID    int           `json:"task_id"`
	NodeName  string        `json:"node_name"`
	LocalID   int           `json:"local_id"`
	State     string        `json:"state"`
	ExitCode  int           `json:"exit_code"`
	CPUTime   time.Duration `json:"cpu_time"`
	MaxRSS    int64         `json:"max_rss_bytes"`
	StartTime *time.Time    `json:"start_time,omitempty"`
	EndTime   *time.Time    `json:"end_time,omitempty"`
}

// StepPerformanceMetrics provides step performance analysis.
type StepPerformanceMetrics struct {
	CPUEfficiency     float64 `json:"cpu_efficiency"`
	MemoryEfficiency  float64 `json:"memory_efficiency"`
	IOEfficiency      float64 `json:"io_efficiency"`
	OverallEfficiency float64 `json:"overall_efficiency"`

	PrimaryBottleneck  string `json:"primary_bottleneck"`
	BottleneckSeverity string `json:"bottleneck_severity"`
	ResourceBalance    string `json:"resource_balance"`

	ThroughputMBPS   float64 `json:"throughput_mbps"`
	LatencyMS        float64 `json:"latency_ms"`
	ScalabilityScore float64 `json:"scalability_score"`
}

// ListJobStepsOptions configures job step listing.
type ListJobStepsOptions struct {
	StepStates []string `json:"step_states,omitempty"`
	NodeNames  []string `json:"node_names,omitempty"`
	StepNames  []string `json:"step_names,omitempty"`
	TaskStates []string `json:"task_states,omitempty"`

	StartTimeAfter  *time.Time     `json:"start_time_after,omitempty"`
	StartTimeBefore *time.Time     `json:"start_time_before,omitempty"`
	EndTimeAfter    *time.Time     `json:"end_time_after,omitempty"`
	EndTimeBefore   *time.Time     `json:"end_time_before,omitempty"`
	MinDuration     *time.Duration `json:"min_duration,omitempty"`
	MaxDuration     *time.Duration `json:"max_duration,omitempty"`

	MinCPUEfficiency     *float64 `json:"min_cpu_efficiency,omitempty"`
	MaxCPUEfficiency     *float64 `json:"max_cpu_efficiency,omitempty"`
	MinMemoryEfficiency  *float64 `json:"min_memory_efficiency,omitempty"`
	MaxMemoryEfficiency  *float64 `json:"max_memory_efficiency,omitempty"`
	MinOverallEfficiency *float64 `json:"min_overall_efficiency,omitempty"`
	MaxOverallEfficiency *float64 `json:"max_overall_efficiency,omitempty"`

	IncludeTaskMetrics         bool `json:"include_task_metrics,omitempty"`
	IncludePerformanceAnalysis bool `json:"include_performance_analysis,omitempty"`
	IncludeResourceTrends      bool `json:"include_resource_trends,omitempty"`
	IncludeBottleneckAnalysis  bool `json:"include_bottleneck_analysis,omitempty"`

	Limit     int    `json:"limit,omitempty"`
	Offset    int    `json:"offset,omitempty"`
	SortBy    string `json:"sort_by,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
}

// JobStepMetricsList represents a list of step metrics.
type JobStepMetricsList struct {
	JobID         string                 `json:"job_id"`
	JobName       string                 `json:"job_name"`
	Steps         []*JobStepWithMetrics  `json:"steps"`
	Summary       *JobStepsSummary       `json:"summary,omitempty"`
	TotalSteps    int                    `json:"total_steps"`
	FilteredSteps int                    `json:"filtered_steps"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// JobStepWithMetrics combines step details and utilization.
type JobStepWithMetrics struct {
	*JobStepDetails     `json:",inline"`
	*JobStepUtilization `json:",inline"`

	Trends       *StepResourceTrends          `json:"trends,omitempty"`
	Comparison   *StepComparison              `json:"comparison,omitempty"`
	Optimization *StepOptimizationSuggestions `json:"optimization,omitempty"`
}

// JobStepsSummary summarizes all steps in a job.
type JobStepsSummary struct {
	TotalSteps      int           `json:"total_steps"`
	CompletedSteps  int           `json:"completed_steps"`
	FailedSteps     int           `json:"failed_steps"`
	RunningSteps    int           `json:"running_steps"`
	TotalDuration   time.Duration `json:"total_duration"`
	AverageDuration time.Duration `json:"average_duration"`

	AverageCPUEfficiency     float64 `json:"average_cpu_efficiency"`
	AverageMemoryEfficiency  float64 `json:"average_memory_efficiency"`
	AverageIOEfficiency      float64 `json:"average_io_efficiency"`
	AverageOverallEfficiency float64 `json:"average_overall_efficiency"`

	TotalCPUHours     float64 `json:"total_cpu_hours"`
	TotalMemoryGBH    float64 `json:"total_memory_gb_hours"`
	TotalIOOperations int64   `json:"total_io_operations"`
	TotalEnergyUsed   float64 `json:"total_energy_used,omitempty"`

	PrimaryBottlenecks   map[string]int `json:"primary_bottlenecks"`
	BottleneckSeverities map[string]int `json:"bottleneck_severities"`
	MostEfficientStep    *string        `json:"most_efficient_step,omitempty"`
	LeastEfficientStep   *string        `json:"least_efficient_step,omitempty"`

	OptimizationPotential float64  `json:"optimization_potential"`
	RecommendedActions    []string `json:"recommended_actions,omitempty"`
}

// StepResourceTrends represents step resource trends.
type StepResourceTrends struct {
	StepID           string             `json:"step_id"`
	CPUTrend         *ResourceTrendData `json:"cpu_trend,omitempty"`
	MemoryTrend      *ResourceTrendData `json:"memory_trend,omitempty"`
	IOTrend          *ResourceTrendData `json:"io_trend,omitempty"`
	NetworkTrend     *ResourceTrendData `json:"network_trend,omitempty"`
	SamplingInterval time.Duration      `json:"sampling_interval"`
	TrendDirection   string             `json:"trend_direction"`
	TrendConfidence  float64            `json:"trend_confidence"`
}

// ResourceTrendData represents trend data for a resource.
type ResourceTrendData struct {
	Values       []float64   `json:"values"`
	Timestamps   []time.Time `json:"timestamps"`
	AverageValue float64     `json:"average_value"`
	MinValue     float64     `json:"min_value"`
	MaxValue     float64     `json:"max_value"`
	StandardDev  float64     `json:"standard_deviation"`
	SlopePerHour float64     `json:"slope_per_hour"`
}

// StepComparison compares a step to job averages.
type StepComparison struct {
	StepID                   string   `json:"step_id"`
	RelativeCPUEfficiency    float64  `json:"relative_cpu_efficiency"`
	RelativeMemoryEfficiency float64  `json:"relative_memory_efficiency"`
	RelativeDuration         float64  `json:"relative_duration"`
	PerformanceRank          int      `json:"performance_rank"`
	EfficiencyPercentile     float64  `json:"efficiency_percentile"`
	ComparisonNotes          []string `json:"comparison_notes,omitempty"`
}

// StepOptimizationSuggestions provides step optimization recommendations.
type StepOptimizationSuggestions struct {
	StepID               string  `json:"step_id"`
	OverallScore         float64 `json:"overall_score"`
	ImprovementPotential float64 `json:"improvement_potential"`

	CPUSuggestions     []OptimizationSuggestion `json:"cpu_suggestions,omitempty"`
	MemorySuggestions  []OptimizationSuggestion `json:"memory_suggestions,omitempty"`
	IOSuggestions      []OptimizationSuggestion `json:"io_suggestions,omitempty"`
	NetworkSuggestions []OptimizationSuggestion `json:"network_suggestions,omitempty"`

	RecommendedCPUs       *int     `json:"recommended_cpus,omitempty"`
	RecommendedMemoryMB   *int     `json:"recommended_memory_mb,omitempty"`
	RecommendedNodes      *int     `json:"recommended_nodes,omitempty"`
	AlternativePartitions []string `json:"alternative_partitions,omitempty"`

	HighPriorityActions   []string `json:"high_priority_actions,omitempty"`
	MediumPriorityActions []string `json:"medium_priority_actions,omitempty"`
	LowPriorityActions    []string `json:"low_priority_actions,omitempty"`
}

// OptimizationSuggestion provides a specific optimization suggestion.
type OptimizationSuggestion struct {
	Type                        string  `json:"type"`
	Severity                    string  `json:"severity"`
	Description                 string  `json:"description"`
	ExpectedBenefit             string  `json:"expected_benefit"`
	ImplementationComplexity    string  `json:"implementation_complexity"`
	EstimatedImprovementPercent float64 `json:"estimated_improvement_percent"`
	ActionRequired              string  `json:"action_required"`
}

// ============================================================================
// Accounting Data Types
// ============================================================================

// AccountingQueryOptions configures accounting queries.
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

// AccountingJobSteps represents job steps from accounting.
type AccountingJobSteps struct {
	JobID    string                 `json:"job_id"`
	JobName  string                 `json:"job_name"`
	User     string                 `json:"user"`
	Account  string                 `json:"account"`
	Steps    []StepAccountingRecord `json:"steps"`
	Summary  *JobAccountingSummary  `json:"summary,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// StepAccountingRecord represents accounting data for a step.
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

	SubmitTime  time.Time     `json:"submit_time"`
	StartTime   *time.Time    `json:"start_time,omitempty"`
	EndTime     *time.Time    `json:"end_time,omitempty"`
	ElapsedTime time.Duration `json:"elapsed_time"`
	CPUTime     time.Duration `json:"cpu_time"`

	AllocCPUs int   `json:"alloc_cpus"`
	ReqCPUs   int   `json:"req_cpus"`
	AllocMem  int64 `json:"alloc_mem"`
	ReqMem    int64 `json:"req_mem"`
	MaxRSS    int64 `json:"max_rss"`
	MaxVMSize int64 `json:"max_vm_size"`

	MaxDiskRead  int64   `json:"max_disk_read"`
	MaxDiskWrite int64   `json:"max_disk_write"`
	AveCPU       float64 `json:"ave_cpu"`
	AveCPUFreq   float64 `json:"ave_cpu_freq"`

	State           string `json:"state"`
	ExitCode        int    `json:"exit_code"`
	DerivedExitCode int    `json:"derived_exit_code"`

	QOS         string                 `json:"qos,omitempty"`
	Constraints string                 `json:"constraints,omitempty"`
	WorkDir     string                 `json:"work_dir,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// JobAccountingSummary summarizes job accounting data.
type JobAccountingSummary struct {
	TotalSteps     int `json:"total_steps"`
	CompletedSteps int `json:"completed_steps"`
	FailedSteps    int `json:"failed_steps"`
	RunningSteps   int `json:"running_steps"`

	TotalCPUTime    time.Duration `json:"total_cpu_time"`
	TotalElapsed    time.Duration `json:"total_elapsed"`
	TotalMemoryUsed int64         `json:"total_memory_used"`
	TotalDiskRead   int64         `json:"total_disk_read"`
	TotalDiskWrite  int64         `json:"total_disk_write"`

	OverallCPUEff float64 `json:"overall_cpu_efficiency"`
	OverallMemEff float64 `json:"overall_memory_efficiency"`
	OverallIOEff  float64 `json:"overall_io_efficiency"`

	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// JobStepAPIData represents real-time step data from the API.
type JobStepAPIData struct {
	StepID   string `json:"step_id"`
	JobID    string `json:"job_id"`
	StepName string `json:"step_name"`
	State    string `json:"state"`

	CurrentCPUUsage    float64 `json:"current_cpu_usage"`
	CurrentMemoryUsage int64   `json:"current_memory_usage"`
	CurrentIOReads     int64   `json:"current_io_reads"`
	CurrentIOWrites    int64   `json:"current_io_writes"`

	TaskCount      int `json:"task_count"`
	RunningTasks   int `json:"running_tasks"`
	CompletedTasks int `json:"completed_tasks"`
	FailedTasks    int `json:"failed_tasks"`

	NetworkBytesIn  int64 `json:"network_bytes_in"`
	NetworkBytesOut int64 `json:"network_bytes_out"`

	ContextSwitches int64 `json:"context_switches"`
	PageFaults      int64 `json:"page_faults"`

	StartTime  *time.Time `json:"start_time,omitempty"`
	LastUpdate time.Time  `json:"last_update"`

	ProcessTree     []ProcessInfo          `json:"process_tree,omitempty"`
	EnvironmentVars map[string]string      `json:"environment_vars,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// ProcessInfo represents a process in the process tree.
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

// SacctQueryOptions configures sacct queries.
type SacctQueryOptions struct {
	JobID      string     `json:"job_id,omitempty"`
	User       string     `json:"user,omitempty"`
	Account    string     `json:"account,omitempty"`
	Partition  string     `json:"partition,omitempty"`
	State      []string   `json:"state,omitempty"`
	StartTime  *time.Time `json:"start_time,omitempty"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	Format     []string   `json:"format,omitempty"`
	Delimiter  string     `json:"delimiter,omitempty"`
	NoHeader   bool       `json:"no_header,omitempty"`
	Parsable   bool       `json:"parsable,omitempty"`
	Brief      bool       `json:"brief,omitempty"`
	Verbose    bool       `json:"verbose,omitempty"`
	Units      string     `json:"units,omitempty"`
	MaxRecords int        `json:"max_records,omitempty"`
}

// SacctJobStepData represents sacct job step data.
type SacctJobStepData struct {
	JobID        string                 `json:"job_id"`
	QueryOptions *SacctQueryOptions     `json:"query_options"`
	Steps        []SacctStepRecord      `json:"steps"`
	TotalSteps   int                    `json:"total_steps"`
	QueryTime    time.Time              `json:"query_time"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// SacctStepRecord represents a sacct step record.
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

	AllocCPUs  int    `json:"alloc_cpus"`
	ReqCPUs    int    `json:"req_cpus"`
	AllocMem   string `json:"alloc_mem"`
	ReqMem     string `json:"req_mem"`
	AllocNodes int    `json:"alloc_nodes"`
	ReqNodes   string `json:"req_nodes"`

	Submit     string `json:"submit"`
	Start      string `json:"start"`
	End        string `json:"end"`
	Elapsed    string `json:"elapsed"`
	CPUTime    string `json:"cpu_time"`
	CPUTimeRAW int64  `json:"cpu_time_raw"`

	MaxRSS       string `json:"max_rss"`
	MaxVMSize    string `json:"max_vm_size"`
	MaxDiskRead  string `json:"max_disk_read"`
	MaxDiskWrite string `json:"max_disk_write"`
	MaxPages     int64  `json:"max_pages"`

	AveCPU       string  `json:"ave_cpu"`
	AveCPUFreq   string  `json:"ave_cpu_freq"`
	AvePages     float64 `json:"ave_pages"`
	AveDiskRead  string  `json:"ave_disk_read"`
	AveDiskWrite string  `json:"ave_disk_write"`

	State           string `json:"state"`
	ExitCode        string `json:"exit_code"`
	DerivedExitCode string `json:"derived_exit_code"`

	QOS          string `json:"qos,omitempty"`
	Priority     int    `json:"priority,omitempty"`
	ReqTRES      string `json:"req_tres,omitempty"`
	AllocTRES    string `json:"alloc_tres,omitempty"`
	TRESUsageIn  string `json:"tres_usage_in,omitempty"`
	TRESUsageOut string `json:"tres_usage_out,omitempty"`

	RawData  map[string]string      `json:"raw_data,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ============================================================================
// Comprehensive Analytics
// ============================================================================

// CPUAnalytics provides detailed CPU analytics.
type CPUAnalytics struct {
	AllocatedCores     int     `json:"allocated_cores"`
	RequestedCores     int     `json:"requested_cores"`
	UsedCores          float64 `json:"used_cores"`
	UtilizationPercent float64 `json:"utilization_percent"`
	EfficiencyPercent  float64 `json:"efficiency_percent"`
	IdleCores          float64 `json:"idle_cores"`
	Oversubscribed     bool    `json:"oversubscribed"`

	CoreMetrics []CPUCoreMetric `json:"core_metrics,omitempty"`

	AverageTemperature     float64 `json:"average_temperature"`
	MaxTemperature         float64 `json:"max_temperature"`
	ThermalThrottleEvents  int64   `json:"thermal_throttle_events"`
	AverageFrequency       float64 `json:"average_frequency_ghz"`
	MaxFrequency           float64 `json:"max_frequency_ghz"`
	FrequencyScalingEvents int64   `json:"frequency_scaling_events"`

	ContextSwitches  int64   `json:"context_switches"`
	Interrupts       int64   `json:"interrupts"`
	SoftInterrupts   int64   `json:"soft_interrupts"`
	LoadAverage1Min  float64 `json:"load_average_1min"`
	LoadAverage5Min  float64 `json:"load_average_5min"`
	LoadAverage15Min float64 `json:"load_average_15min"`

	L1CacheHitRate float64 `json:"l1_cache_hit_rate"`
	L2CacheHitRate float64 `json:"l2_cache_hit_rate"`
	L3CacheHitRate float64 `json:"l3_cache_hit_rate"`
	L1CacheMisses  int64   `json:"l1_cache_misses"`
	L2CacheMisses  int64   `json:"l2_cache_misses"`
	L3CacheMisses  int64   `json:"l3_cache_misses"`

	InstructionsPerCycle int64 `json:"instructions_per_cycle"`
	BranchMispredictions int64 `json:"branch_mispredictions"`
	TotalInstructions    int64 `json:"total_instructions"`

	Recommendations []OptimizationRecommendation `json:"recommendations,omitempty"`
	Bottlenecks     []PerformanceBottleneck      `json:"bottlenecks,omitempty"`

	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// CPUCoreMetric represents per-core CPU metrics.
type CPUCoreMetric struct {
	CoreID          int     `json:"core_id"`
	Utilization     float64 `json:"utilization_percent"`
	Frequency       float64 `json:"frequency_ghz"`
	Temperature     float64 `json:"temperature_celsius"`
	LoadAverage     float64 `json:"load_average"`
	ContextSwitches int64   `json:"context_switches"`
	Interrupts      int64   `json:"interrupts"`
}

// MemoryAnalytics provides detailed memory analytics.
type MemoryAnalytics struct {
	AllocatedBytes     int64   `json:"allocated_bytes"`
	RequestedBytes     int64   `json:"requested_bytes"`
	UsedBytes          int64   `json:"used_bytes"`
	UtilizationPercent float64 `json:"utilization_percent"`
	EfficiencyPercent  float64 `json:"efficiency_percent"`
	FreeBytes          int64   `json:"free_bytes"`
	Overcommitted      bool    `json:"overcommitted"`

	ResidentSetSize   int64 `json:"resident_set_size"`
	VirtualMemorySize int64 `json:"virtual_memory_size"`
	SharedMemory      int64 `json:"shared_memory"`
	BufferedMemory    int64 `json:"buffered_memory"`
	CachedMemory      int64 `json:"cached_memory"`

	NUMANodes []NUMANodeMetrics `json:"numa_nodes,omitempty"`

	BandwidthUtilization float64 `json:"bandwidth_utilization_percent"`
	MemoryBandwidthMBPS  int64   `json:"memory_bandwidth_mbps"`
	PeakBandwidthMBPS    int64   `json:"peak_bandwidth_mbps"`

	PageFaults      int64 `json:"page_faults"`
	MajorPageFaults int64 `json:"major_page_faults"`
	MinorPageFaults int64 `json:"minor_page_faults"`
	PageSwaps       int64 `json:"page_swaps"`

	RandomAccess     float64 `json:"random_access_percent"`
	SequentialAccess float64 `json:"sequential_access_percent"`
	LocalityScore    float64 `json:"locality_score"`

	MemoryLeaks []MemoryLeak `json:"memory_leaks,omitempty"`

	Recommendations []OptimizationRecommendation `json:"recommendations,omitempty"`
	Bottlenecks     []PerformanceBottleneck      `json:"bottlenecks,omitempty"`

	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NUMANodeMetrics represents NUMA node metrics.
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

// MemoryLeak represents a detected memory leak.
type MemoryLeak struct {
	LeakType    string `json:"leak_type"`
	SizeBytes   int64  `json:"size_bytes"`
	GrowthRate  int64  `json:"growth_rate"`
	Location    string `json:"location"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
}

// IOAnalytics provides detailed I/O analytics.
type IOAnalytics struct {
	ReadBytes          int64   `json:"read_bytes"`
	WriteBytes         int64   `json:"write_bytes"`
	ReadOperations     int64   `json:"read_operations"`
	WriteOperations    int64   `json:"write_operations"`
	UtilizationPercent float64 `json:"utilization_percent"`
	EfficiencyPercent  float64 `json:"efficiency_percent"`

	AverageReadBandwidth  float64 `json:"average_read_bandwidth_mbps"`
	AverageWriteBandwidth float64 `json:"average_write_bandwidth_mbps"`
	PeakReadBandwidth     float64 `json:"peak_read_bandwidth_mbps"`
	PeakWriteBandwidth    float64 `json:"peak_write_bandwidth_mbps"`

	AverageReadLatency  float64 `json:"average_read_latency_ms"`
	AverageWriteLatency float64 `json:"average_write_latency_ms"`
	MaxReadLatency      float64 `json:"max_read_latency_ms"`
	MaxWriteLatency     float64 `json:"max_write_latency_ms"`

	QueueDepth    float64 `json:"queue_depth"`
	MaxQueueDepth float64 `json:"max_queue_depth"`
	QueueTime     float64 `json:"queue_time_ms"`

	RandomAccessPercent     float64 `json:"random_access_percent"`
	SequentialAccessPercent float64 `json:"sequential_access_percent"`

	AverageIOSize int64 `json:"average_io_size_bytes"`
	MaxIOSize     int64 `json:"max_io_size_bytes"`
	MinIOSize     int64 `json:"min_io_size_bytes"`

	StorageDevices []StorageDevice `json:"storage_devices,omitempty"`

	Recommendations []OptimizationRecommendation `json:"recommendations,omitempty"`
	Bottlenecks     []PerformanceBottleneck      `json:"bottlenecks,omitempty"`

	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// StorageDevice represents storage device metrics.
type StorageDevice struct {
	DeviceName     string  `json:"device_name"`
	DeviceType     string  `json:"device_type"`
	MountPoint     string  `json:"mount_point"`
	TotalCapacity  int64   `json:"total_capacity"`
	UsedCapacity   int64   `json:"used_capacity"`
	AvailCapacity  int64   `json:"available_capacity"`
	Utilization    float64 `json:"utilization_percent"`
	IOPS           int64   `json:"iops"`
	ThroughputMBPS int64   `json:"throughput_mbps"`
}

// JobComprehensiveAnalytics provides complete job analytics.
type JobComprehensiveAnalytics struct {
	JobID     uint32        `json:"job_id"`
	JobName   string        `json:"job_name"`
	StartTime time.Time     `json:"start_time"`
	EndTime   *time.Time    `json:"end_time,omitempty"`
	Duration  time.Duration `json:"duration"`
	Status    string        `json:"status"`

	CPUAnalytics    *CPUAnalytics    `json:"cpu_analytics"`
	MemoryAnalytics *MemoryAnalytics `json:"memory_analytics"`
	IOAnalytics     *IOAnalytics     `json:"io_analytics"`

	OverallEfficiency float64 `json:"overall_efficiency_percent"`

	CrossResourceAnalysis *CrossResourceAnalysis   `json:"cross_resource_analysis"`
	OptimalConfiguration  *OptimalJobConfiguration `json:"optimal_configuration"`

	Recommendations []OptimizationRecommendation `json:"recommendations,omitempty"`
	Bottlenecks     []PerformanceBottleneck      `json:"bottlenecks,omitempty"`

	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// CrossResourceAnalysis analyzes cross-resource interactions.
type CrossResourceAnalysis struct {
	PrimaryBottleneck     string  `json:"primary_bottleneck"`
	SecondaryBottleneck   string  `json:"secondary_bottleneck"`
	BottleneckSeverity    string  `json:"bottleneck_severity"`
	ResourceBalance       string  `json:"resource_balance"`
	OptimizationPotential float64 `json:"optimization_potential_percent"`
	ScalabilityScore      float64 `json:"scalability_score"`
	ResourceWaste         float64 `json:"resource_waste_percent"`
	LoadBalanceScore      float64 `json:"load_balance_score"`
}

// OptimalJobConfiguration provides resource recommendations.
type OptimalJobConfiguration struct {
	RecommendedCPUs    int               `json:"recommended_cpus"`
	RecommendedMemory  int64             `json:"recommended_memory_bytes"`
	RecommendedNodes   int               `json:"recommended_nodes"`
	RecommendedRuntime int               `json:"recommended_runtime_minutes"`
	ExpectedSpeedup    float64           `json:"expected_speedup"`
	CostReduction      float64           `json:"cost_reduction_percent"`
	ConfigChanges      map[string]string `json:"config_changes,omitempty"`
}

// ============================================================================
// Performance History and Trends
// ============================================================================

// PerformanceHistoryOptions configures performance history queries.
type PerformanceHistoryOptions struct {
	StartTime     *time.Time `json:"start_time,omitempty"`
	EndTime       *time.Time `json:"end_time,omitempty"`
	Interval      string     `json:"interval,omitempty"`
	MetricTypes   []string   `json:"metric_types,omitempty"`
	IncludeSteps  bool       `json:"include_steps"`
	IncludeTrends bool       `json:"include_trends"`
}

// JobPerformanceHistory provides historical performance data.
type JobPerformanceHistory struct {
	JobID     string    `json:"job_id"`
	JobName   string    `json:"job_name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`

	TimeSeriesData []PerformanceSnapshot    `json:"time_series_data"`
	Statistics     PerformanceStatistics    `json:"statistics"`
	Trends         *PerformanceTrendAnalysis `json:"trends,omitempty"`
	Anomalies      []PerformanceAnomaly     `json:"anomalies,omitempty"`
}

// PerformanceSnapshot represents a point-in-time snapshot.
type PerformanceSnapshot struct {
	Timestamp         time.Time `json:"timestamp"`
	CPUUtilization    float64   `json:"cpu_utilization"`
	MemoryUtilization float64   `json:"memory_utilization"`
	IOBandwidth       float64   `json:"io_bandwidth_mbps"`
	GPUUtilization    float64   `json:"gpu_utilization,omitempty"`
	NetworkBandwidth  float64   `json:"network_bandwidth_mbps,omitempty"`
	PowerUsage        float64   `json:"power_usage_watts,omitempty"`
	Efficiency        float64   `json:"efficiency_score"`
}

// PerformanceStatistics is defined in performance.go

// PerformanceTrendAnalysis provides trend analysis. (Note: truncated to PerformanceTrendAnalysis to match interface)
type PerformanceTrendAnalysis struct {
	CPUTrend        TrendInfo `json:"cpu_trend"`
	MemoryTrend     TrendInfo `json:"memory_trend"`
	IOTrend         TrendInfo `json:"io_trend"`
	EfficiencyTrend TrendInfo `json:"efficiency_trend"`

	PredictedCPU     float64       `json:"predicted_cpu"`
	PredictedMemory  float64       `json:"predicted_memory"`
	PredictedRuntime time.Duration `json:"predicted_runtime"`
}

// TrendInfo describes a trend.
type TrendInfo struct {
	Direction  string  `json:"direction"`
	Slope      float64 `json:"slope"`
	Confidence float64 `json:"confidence"`
	ChangeRate float64 `json:"change_rate_percent"`
}

// PerformanceAnomaly represents a detected performance anomaly.
type PerformanceAnomaly struct {
	Timestamp   time.Time `json:"timestamp"`
	Type        string    `json:"type"`
	Metric      string    `json:"metric"`
	Severity    string    `json:"severity"`
	Value       float64   `json:"value"`
	Expected    float64   `json:"expected"`
	Deviation   float64   `json:"deviation_percent"`
	Description string    `json:"description"`
}

// TrendAnalysisOptions configures trend analysis.
type TrendAnalysisOptions struct {
	StartTime   *time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	Granularity string     `json:"granularity"`
	Partitions  []string   `json:"partitions,omitempty"`
	UserFilter  []string   `json:"users,omitempty"`
	JobTypes    []string   `json:"job_types,omitempty"`
	MinJobSize  *int       `json:"min_job_size,omitempty"`
}

// PerformanceTrends represents cluster-wide performance trends.
type PerformanceTrends struct {
	TimeRange   TimeRange `json:"time_range"`
	Granularity string    `json:"granularity"`

	ClusterUtilization []UtilizationPoint `json:"cluster_utilization"`
	ClusterEfficiency  []EfficiencyPoint  `json:"cluster_efficiency"`

	PartitionTrends map[string]*PartitionTrend `json:"partition_trends"`

	CPUTrends    ResourceTrend `json:"cpu_trends"`
	MemoryTrends ResourceTrend `json:"memory_trends"`
	GPUTrends    ResourceTrend `json:"gpu_trends"`

	JobSizeTrends     []JobSizeTrend     `json:"job_size_trends"`
	JobDurationTrends []JobDurationTrend `json:"job_duration_trends"`

	Insights []TrendInsight `json:"insights"`
}

// TimeRange represents a time range.
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// UtilizationPoint represents a utilization data point.
type UtilizationPoint struct {
	Timestamp   time.Time `json:"timestamp"`
	Utilization float64   `json:"utilization"`
	JobCount    int       `json:"job_count"`
}

// EfficiencyPoint represents an efficiency data point.
type EfficiencyPoint struct {
	Timestamp  time.Time `json:"timestamp"`
	Efficiency float64   `json:"efficiency"`
	JobCount   int       `json:"job_count"`
}

// PartitionTrend represents partition-specific trends.
type PartitionTrend struct {
	PartitionName string             `json:"partition_name"`
	Utilization   []UtilizationPoint `json:"utilization"`
	Efficiency    []EfficiencyPoint  `json:"efficiency"`
	JobCounts     []JobCountPoint    `json:"job_counts"`
	QueueLength   []QueueLengthPoint `json:"queue_length"`
}

// ResourceTrend represents resource trend data.
type ResourceTrend struct {
	Average    []float64   `json:"average"`
	Peak       []float64   `json:"peak"`
	Timestamps []time.Time `json:"timestamps"`
	Trend      TrendInfo   `json:"trend"`
}

// JobSizeTrend represents job size trends.
type JobSizeTrend struct {
	Timestamp   time.Time `json:"timestamp"`
	SmallJobs   int       `json:"small_jobs"`
	MediumJobs  int       `json:"medium_jobs"`
	LargeJobs   int       `json:"large_jobs"`
	AverageSize float64   `json:"average_size"`
}

// JobDurationTrend represents job duration trends.
type JobDurationTrend struct {
	Timestamp       time.Time     `json:"timestamp"`
	ShortJobs       int           `json:"short_jobs"`
	MediumJobs      int           `json:"medium_jobs"`
	LongJobs        int           `json:"long_jobs"`
	AverageDuration time.Duration `json:"average_duration"`
}

// JobCountPoint represents job count data.
type JobCountPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Running   int       `json:"running"`
	Pending   int       `json:"pending"`
	Total     int       `json:"total"`
}

// QueueLengthPoint represents queue length data.
type QueueLengthPoint struct {
	Timestamp   time.Time     `json:"timestamp"`
	QueueLength int           `json:"queue_length"`
	WaitTime    time.Duration `json:"average_wait_time"`
}

// TrendInsight represents an insight from trend analysis.
type TrendInsight struct {
	Type           string    `json:"type"`
	Category       string    `json:"category"`
	Severity       string    `json:"severity"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Timestamp      time.Time `json:"timestamp"`
	Confidence     float64   `json:"confidence"`
	Recommendation string    `json:"recommendation,omitempty"`
}

// ============================================================================
// User Efficiency Trends
// ============================================================================

// EfficiencyTrendOptions configures efficiency trend analysis.
type EfficiencyTrendOptions struct {
	StartTime        *time.Time `json:"start_time,omitempty"`
	EndTime          *time.Time `json:"end_time,omitempty"`
	Granularity      string     `json:"granularity"`
	JobTypes         []string   `json:"job_types,omitempty"`
	Partitions       []string   `json:"partitions,omitempty"`
	CompareToAverage bool       `json:"compare_to_average"`
}

// UserEfficiencyTrends represents user efficiency trends.
type UserEfficiencyTrends struct {
	UserID    string    `json:"user_id"`
	TimeRange TimeRange `json:"time_range"`

	EfficiencyHistory []EfficiencyDataPoint `json:"efficiency_history"`

	CPUUtilizationTrend    []float64 `json:"cpu_utilization_trend"`
	MemoryUtilizationTrend []float64 `json:"memory_utilization_trend"`

	JobCountTrend    []int           `json:"job_count_trend"`
	JobSizeTrend     []float64       `json:"job_size_trend"`
	JobDurationTrend []time.Duration `json:"job_duration_trend"`

	AverageEfficiency        float64 `json:"average_efficiency"`
	ClusterAverageEfficiency float64 `json:"cluster_average_efficiency"`
	EfficiencyRank           int     `json:"efficiency_rank"`
	EfficiencyPercentile     float64 `json:"efficiency_percentile"`

	ImprovementRate float64   `json:"improvement_rate"`
	BestPeriod      TimeRange `json:"best_period"`
	WorstPeriod     TimeRange `json:"worst_period"`

	Recommendations []string `json:"recommendations"`
}

// EfficiencyDataPoint represents an efficiency data point.
type EfficiencyDataPoint struct {
	Timestamp   time.Time `json:"timestamp"`
	Efficiency  float64   `json:"efficiency"`
	JobCount    int       `json:"job_count"`
	CPUHours    float64   `json:"cpu_hours"`
	MemoryGBH   float64   `json:"memory_gb_hours"`
	WastedHours float64   `json:"wasted_cpu_hours"`
}

// ============================================================================
// Batch Analysis
// ============================================================================

// BatchAnalysisOptions configures batch analysis.
type BatchAnalysisOptions struct {
	IncludeDetails     bool     `json:"include_details"`
	IncludeComparison  bool     `json:"include_comparison"`
	ComparisonBaseline string   `json:"comparison_baseline"`
	MetricTypes        []string `json:"metric_types"`
}

// BatchJobAnalysis represents batch job analysis results.
type BatchJobAnalysis struct {
	JobCount      int       `json:"job_count"`
	AnalyzedCount int       `json:"analyzed_count"`
	FailedCount   int       `json:"failed_count"`
	TimeRange     TimeRange `json:"time_range"`

	AggregateStats BatchStatistics `json:"aggregate_stats"`

	JobAnalyses []JobAnalysisSummary `json:"job_analyses,omitempty"`

	Comparison *BatchComparison `json:"comparison,omitempty"`

	Patterns []BatchPattern `json:"patterns"`
	Outliers []string       `json:"outlier_job_ids"`

	BatchRecommendations []BatchRecommendation `json:"recommendations"`
}

// BatchStatistics represents aggregate statistics.
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

// JobAnalysisSummary provides a job analysis summary.
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

// BatchComparison compares batch to baseline.
type BatchComparison struct {
	Baseline        string  `json:"baseline"`
	EfficiencyDelta float64 `json:"efficiency_delta"`
	RuntimeDelta    float64 `json:"runtime_delta"`
	WasteDelta      float64 `json:"waste_delta"`
	CostDelta       float64 `json:"cost_delta"`
}

// BatchPattern represents a detected pattern.
type BatchPattern struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	JobCount    int      `json:"job_count"`
	JobIDs      []string `json:"job_ids"`
	Impact      string   `json:"impact"`
	Confidence  float64  `json:"confidence"`
}

// BatchRecommendation provides batch-level recommendations.
type BatchRecommendation struct {
	Category     string   `json:"category"`
	Priority     string   `json:"priority"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Impact       string   `json:"impact"`
	JobsAffected []string `json:"jobs_affected"`
}

// ResourceWaste represents wasted resources.
type ResourceWaste struct {
	CPUCoreHours  float64 `json:"cpu_core_hours"`
	MemoryGBHours float64 `json:"memory_gb_hours"`
	GPUHours      float64 `json:"gpu_hours"`
	EstimatedCost float64 `json:"estimated_cost"`
}

// ============================================================================
// Workflow Analysis
// ============================================================================

// WorkflowAnalysisOptions configures workflow analysis.
type WorkflowAnalysisOptions struct {
	IncludeDependencies bool       `json:"include_dependencies"`
	IncludeBottlenecks  bool       `json:"include_bottlenecks"`
	IncludeOptimization bool       `json:"include_optimization"`
	TimeWindow          *TimeRange `json:"time_window,omitempty"`
}

// WorkflowPerformance provides workflow performance analysis.
type WorkflowPerformance struct {
	WorkflowID    string     `json:"workflow_id"`
	WorkflowName  string     `json:"workflow_name"`
	TotalJobs     int        `json:"total_jobs"`
	CompletedJobs int        `json:"completed_jobs"`
	StartTime     time.Time  `json:"start_time"`
	EndTime       *time.Time `json:"end_time,omitempty"`

	TotalDuration     time.Duration `json:"total_duration"`
	WallClockTime     time.Duration `json:"wall_clock_time"`
	Parallelization   float64       `json:"parallelization_efficiency"`
	OverallEfficiency float64       `json:"overall_efficiency"`

	Stages []WorkflowStage `json:"stages"`

	CriticalPath         []string      `json:"critical_path_job_ids"`
	CriticalPathDuration time.Duration `json:"critical_path_duration"`

	Bottlenecks []WorkflowBottleneck `json:"bottlenecks,omitempty"`

	Dependencies *WorkflowDependencies `json:"dependencies,omitempty"`

	Optimization *WorkflowOptimization `json:"optimization,omitempty"`
}

// WorkflowStage represents a workflow stage.
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

// WorkflowBottleneck represents a workflow bottleneck.
type WorkflowBottleneck struct {
	Type        string        `json:"type"`
	Location    string        `json:"location"`
	Description string        `json:"description"`
	Impact      time.Duration `json:"impact_duration"`
	Severity    string        `json:"severity"`
}

// WorkflowDependencies represents workflow dependencies.
type WorkflowDependencies struct {
	DependencyGraph map[string][]string `json:"dependency_graph"`
	MaxDepth        int                 `json:"max_depth"`
	MaxWidth        int                 `json:"max_width"`
	TotalEdges      int                 `json:"total_edges"`
}

// WorkflowOptimization provides workflow optimization suggestions.
type WorkflowOptimization struct {
	PotentialSpeedup     float64                   `json:"potential_speedup"`
	OptimizedDuration    time.Duration             `json:"optimized_duration"`
	OptimizedParallelism int                       `json:"optimized_parallelism"`
	RecommendedChanges   []string                  `json:"recommended_changes"`
	ResourceReallocation map[string]ResourceChange `json:"resource_reallocation"`
}

// ResourceChange represents a resource change recommendation.
type ResourceChange struct {
	JobID             string `json:"job_id"`
	CurrentCPUs       int    `json:"current_cpus"`
	RecommendedCPUs   int    `json:"recommended_cpus"`
	CurrentMemory     int64  `json:"current_memory"`
	RecommendedMemory int64  `json:"recommended_memory"`
	Reason            string `json:"reason"`
}

// ============================================================================
// Efficiency Reports
// ============================================================================

// ReportOptions configures efficiency report generation.
type ReportOptions struct {
	ReportType    string    `json:"report_type"`
	TimeRange     TimeRange `json:"time_range"`
	Partitions    []string  `json:"partitions,omitempty"`
	Users         []string  `json:"users,omitempty"`
	IncludeCharts bool      `json:"include_charts"`
	Format        string    `json:"format"`
}

// EfficiencyReport represents a complete efficiency report.
type EfficiencyReport struct {
	ReportID    string    `json:"report_id"`
	GeneratedAt time.Time `json:"generated_at"`
	TimeRange   TimeRange `json:"time_range"`
	ReportType  string    `json:"report_type"`

	Summary ExecutiveSummary `json:"summary"`

	ClusterOverview   *ClusterOverview    `json:"cluster_overview,omitempty"`
	PartitionAnalysis []PartitionAnalysis `json:"partition_analysis,omitempty"`
	UserAnalysis      []UserAnalysis      `json:"user_analysis,omitempty"`
	ResourceAnalysis  *ResourceAnalysis   `json:"resource_analysis,omitempty"`

	TrendAnalysis *ReportTrendAnalysis `json:"trend_analysis,omitempty"`

	Recommendations []ReportRecommendation `json:"recommendations"`

	Charts []ChartData `json:"charts,omitempty"`
}

// ExecutiveSummary provides a high-level summary.
type ExecutiveSummary struct {
	TotalJobs            int      `json:"total_jobs"`
	AverageEfficiency    float64  `json:"average_efficiency"`
	TotalCPUHours        float64  `json:"total_cpu_hours"`
	WastedCPUHours       float64  `json:"wasted_cpu_hours"`
	EstimatedCostSavings float64  `json:"estimated_cost_savings"`
	KeyFindings          []string `json:"key_findings"`
	ImprovementAreas     []string `json:"improvement_areas"`
}

// ClusterOverview provides cluster-level overview.
type ClusterOverview struct {
	ClusterUtilization float64 `json:"cluster_utilization"`
	ClusterEfficiency  float64 `json:"cluster_efficiency"`
	TotalNodes         int     `json:"total_nodes"`
	ActiveNodes        int     `json:"active_nodes"`
	TotalCPUCores      int     `json:"total_cpu_cores"`
	TotalMemoryGB      float64 `json:"total_memory_gb"`
	TotalGPUs          int     `json:"total_gpus"`
}

// PartitionAnalysis provides partition-level analysis.
type PartitionAnalysis struct {
	PartitionName  string   `json:"partition_name"`
	JobCount       int      `json:"job_count"`
	Utilization    float64  `json:"utilization"`
	Efficiency     float64  `json:"efficiency"`
	TopUsers       []string `json:"top_users"`
	CommonJobTypes []string `json:"common_job_types"`
	Issues         []string `json:"issues"`
}

// UserAnalysis provides user-level analysis.
type UserAnalysis struct {
	UserID            string  `json:"user_id"`
	JobCount          int     `json:"job_count"`
	AverageEfficiency float64 `json:"average_efficiency"`
	TotalCPUHours     float64 `json:"total_cpu_hours"`
	WastedHours       float64 `json:"wasted_hours"`
	Rank              int     `json:"rank"`
	Trend             string  `json:"trend"`
}

// ResourceAnalysis provides resource-level analysis.
type ResourceAnalysis struct {
	CPUAnalysis    ResourceTypeAnalysis `json:"cpu_analysis"`
	MemoryAnalysis ResourceTypeAnalysis `json:"memory_analysis"`
	GPUAnalysis    ResourceTypeAnalysis `json:"gpu_analysis"`
	IOAnalysis     ResourceTypeAnalysis `json:"io_analysis"`
}

// ResourceTypeAnalysis provides analysis for a resource type.
type ResourceTypeAnalysis struct {
	AverageUtilization    float64  `json:"average_utilization"`
	PeakUtilization       float64  `json:"peak_utilization"`
	WastePercentage       float64  `json:"waste_percentage"`
	TopWasters            []string `json:"top_waster_job_ids"`
	OptimizationPotential float64  `json:"optimization_potential"`
}

// ReportTrendAnalysis provides trend analysis for reports.
type ReportTrendAnalysis struct {
	EfficiencyTrend     TrendInfo `json:"efficiency_trend"`
	UtilizationTrend    TrendInfo `json:"utilization_trend"`
	JobVolumeTrend      TrendInfo `json:"job_volume_trend"`
	PredictedEfficiency float64   `json:"predicted_efficiency_next_period"`
}

// ReportRecommendation provides a report recommendation.
type ReportRecommendation struct {
	Category         string   `json:"category"`
	Priority         string   `json:"priority"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	ExpectedImpact   string   `json:"expected_impact"`
	AffectedEntities []string `json:"affected_entities"`
	Implementation   string   `json:"implementation_guidance"`
}

// ChartData represents chart data for visualization.
type ChartData struct {
	ChartID     string                 `json:"chart_id"`
	ChartType   string                 `json:"chart_type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Data        map[string]interface{} `json:"data"`
	Options     map[string]interface{} `json:"options"`
}

// ============================================================================
// Extended Diagnostics
// ============================================================================

// ExtendedDiagnostics provides extended cluster diagnostics.
type ExtendedDiagnostics struct {
	// Basic cluster stats
	JobsSubmitted int `json:"jobs_submitted"`
	JobsStarted   int `json:"jobs_started"`
	JobsCompleted int `json:"jobs_completed"`
	JobsPending   int `json:"jobs_pending"`
	JobsRunning   int `json:"jobs_running"`

	// Extended job statistics
	JobsFailed   int `json:"jobs_failed"`
	JobsCanceled int `json:"jobs_canceled"`
	JobsTimeout  int `json:"jobs_timeout"`

	// Backfill statistics
	BackfillActive        bool  `json:"backfill_active"`
	BackfillJobsTotal     int   `json:"backfill_jobs_total"`
	BackfillJobsRecent    int   `json:"backfill_jobs_recent"`
	BackfillCycleCount    int   `json:"backfill_cycle_count"`
	BackfillCycleMeanTime int64 `json:"backfill_cycle_mean_time"`

	// Server performance
	ServerThreadCount int   `json:"server_thread_count"`
	AgentQueueSize    int   `json:"agent_queue_size"`
	ScheduleCycleMax  int   `json:"schedule_cycle_max"`
	ScheduleCycleLast int   `json:"schedule_cycle_last"`
	ScheduleCycleMean int64 `json:"schedule_cycle_mean"`

	// RPC statistics
	RPCsTotal     int `json:"rpcs_total"`
	RPCsPending   int `json:"rpcs_pending"`
	RPCsCompleted int `json:"rpcs_completed"`

	// Metadata
	DiagTime time.Time              `json:"diag_time"`
	RawData  map[string]interface{} `json:"raw_data,omitempty"`
}

// OptimizationRecommendation represents a specific optimization suggestion.
type OptimizationRecommendation struct {
	Resource    string                 `json:"resource"`    // Resource type (CPU, Memory, GPU, IO, Network)
	Type        string                 `json:"type"`        // Type of recommendation (reduction, increase, configuration, pattern)
	Current     interface{}            `json:"current"`     // Current value/configuration
	Recommended interface{}            `json:"recommended"` // Recommended value/configuration
	Reason      string                 `json:"reason"`      // Why this recommendation is made
	Impact      string                 `json:"impact"`      // Expected impact of the change
	Confidence  float64                `json:"confidence"`  // Confidence level (0-1)
	Details     map[string]interface{} `json:"details,omitempty"`
}

// PerformanceStatistics contains aggregated performance statistics.
type PerformanceStatistics struct {
	// CPU statistics
	AverageCPU float64 `json:"average_cpu"`
	PeakCPU    float64 `json:"peak_cpu"`
	MinCPU     float64 `json:"min_cpu"`
	StdDevCPU  float64 `json:"stddev_cpu"`

	// Memory statistics
	AverageMemory float64 `json:"average_memory"`
	PeakMemory    float64 `json:"peak_memory"`
	MinMemory     float64 `json:"min_memory"`
	StdDevMemory  float64 `json:"stddev_memory"`

	// IO statistics
	AverageIO float64 `json:"average_io"`
	PeakIO    float64 `json:"peak_io"`
	MinIO     float64 `json:"min_io"`
	StdDevIO  float64 `json:"stddev_io"`

	// Efficiency statistics
	AverageEfficiency float64 `json:"average_efficiency"`
	MedianEfficiency  float64 `json:"median_efficiency"`
	StdDevEfficiency  float64 `json:"stddev_efficiency"`
	BestEfficiency    float64 `json:"best_efficiency"`
	WorstEfficiency   float64 `json:"worst_efficiency"`

	// Runtime statistics
	AverageRuntime time.Duration `json:"average_runtime"`
	MedianRuntime  time.Duration `json:"median_runtime"`

	// Recommendations
	OptimalResources ResourceRecommendation `json:"optimal_resources"`
}

// PerformanceTrendAnalysiss provides trend analysis for performance metrics.
type PerformanceTrendAnalysiss struct {
	CPUTrend        TrendInfo `json:"cpu_trend"`
	MemoryTrend     TrendInfo `json:"memory_trend"`
	IOTrend         TrendInfo `json:"io_trend"`
	EfficiencyTrend TrendInfo `json:"efficiency_trend"`

	PredictedCPU     float64       `json:"predicted_cpu"`
	PredictedMemory  float64       `json:"predicted_memory"`
	PredictedRuntime time.Duration `json:"predicted_runtime"`
}

// ResourceRecommendation contains recommended resource allocations.
type ResourceRecommendation struct {
	CPUs      int     `json:"cpus"`
	MemoryGB  float64 `json:"memory_gb"`
	GPUs      int     `json:"gpus"`
	Reasoning string  `json:"reasoning"`
}
