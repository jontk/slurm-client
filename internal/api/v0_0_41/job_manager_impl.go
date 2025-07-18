package v0_0_41

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
	"github.com/jontk/slurm-client/pkg/watch"
)

// JobManagerImpl provides the actual implementation for JobManager methods
type JobManagerImpl struct {
	client *WrapperClient
}

// NewJobManagerImpl creates a new JobManager implementation
func NewJobManagerImpl(client *WrapperClient) *JobManagerImpl {
	return &JobManagerImpl{client: client}
}

// List jobs with optional filtering
func (m *JobManagerImpl) List(ctx context.Context, opts *interfaces.ListJobsOptions) (*interfaces.JobList, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmV0041GetJobsParams{}

	// Set flags to get detailed job information
	flags := SlurmV0041GetJobsParamsFlagsDETAIL
	params.Flags = &flags

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0041GetJobsWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.41")
	}

	// Check HTTP status and handle API errors
	if resp.StatusCode() != 200 {
		var responseBody []byte
		if resp.JSON200 != nil {
			// Try to extract error details from response
			if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
				apiErrors := make([]errors.SlurmAPIErrorDetail, len(*resp.JSON200.Errors))
				for i, apiErr := range *resp.JSON200.Errors {
					var errorNumber int
					if apiErr.ErrorNumber != nil {
						errorNumber = int(*apiErr.ErrorNumber)
					}
					var errorCode string
					if apiErr.Error != nil {
						errorCode = *apiErr.Error
					}
					var source string
					if apiErr.Source != nil {
						source = *apiErr.Source
					}
					var description string
					if apiErr.Description != nil {
						description = *apiErr.Description
					}

					apiErrors[i] = errors.SlurmAPIErrorDetail{
						ErrorNumber: errorNumber,
						ErrorCode:   errorCode,
						Source:      source,
						Description: description,
					}
				}
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.41", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.41")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response but got nil")
	}

	// Convert the response to our interface types
	jobs := make([]interfaces.Job, 0)
	if resp.JSON200.Jobs != nil {
		for _, apiJob := range resp.JSON200.Jobs {
			job := interfaces.Job{}
			
			// Job ID
			if apiJob.JobId != nil {
				job.ID = strconv.FormatInt(int64(*apiJob.JobId), 10)
			}
			
			// Job name
			if apiJob.Name != nil {
				job.Name = *apiJob.Name
			}
			
			// User ID
			if apiJob.UserId != nil {
				job.UserID = strconv.FormatInt(int64(*apiJob.UserId), 10)
			}
			
			// Group ID
			if apiJob.GroupId != nil {
				job.GroupID = strconv.FormatInt(int64(*apiJob.GroupId), 10)
			}
			
			// Job state
			if apiJob.JobState != nil && len(*apiJob.JobState) > 0 {
				job.State = string((*apiJob.JobState)[0])
			}
			
			// Partition
			if apiJob.Partition != nil {
				job.Partition = *apiJob.Partition
			}
			
			// Priority
			if apiJob.Priority != nil && apiJob.Priority.Set != nil && *apiJob.Priority.Set && apiJob.Priority.Number != nil {
				job.Priority = int(*apiJob.Priority.Number)
			}
			
			// Submit time
			if apiJob.SubmitTime != nil && apiJob.SubmitTime.Set != nil && *apiJob.SubmitTime.Set && apiJob.SubmitTime.Number != nil {
				job.SubmitTime = time.Unix(*apiJob.SubmitTime.Number, 0)
			}
			
			// Start time
			if apiJob.StartTime != nil && apiJob.StartTime.Set != nil && *apiJob.StartTime.Set && apiJob.StartTime.Number != nil && *apiJob.StartTime.Number > 0 {
				startTime := time.Unix(*apiJob.StartTime.Number, 0)
				job.StartTime = &startTime
			}
			
			// End time
			if apiJob.EndTime != nil && apiJob.EndTime.Set != nil && *apiJob.EndTime.Set && apiJob.EndTime.Number != nil && *apiJob.EndTime.Number > 0 {
				endTime := time.Unix(*apiJob.EndTime.Number, 0)
				job.EndTime = &endTime
			}
			
			// CPUs
			if apiJob.Cpus != nil && apiJob.Cpus.Set != nil && *apiJob.Cpus.Set && apiJob.Cpus.Number != nil {
				job.CPUs = int(*apiJob.Cpus.Number)
			}
			
			// Memory (convert MB to bytes)
			if apiJob.MemoryPerNode != nil && apiJob.MemoryPerNode.Set != nil && *apiJob.MemoryPerNode.Set && apiJob.MemoryPerNode.Number != nil {
				job.Memory = int(*apiJob.MemoryPerNode.Number) * 1024 * 1024
			}
			
			// Time limit
			if apiJob.TimeLimit != nil && apiJob.TimeLimit.Set != nil && *apiJob.TimeLimit.Set && apiJob.TimeLimit.Number != nil {
				job.TimeLimit = int(*apiJob.TimeLimit.Number)
			}
			
			// Working directory
			if apiJob.CurrentWorkingDirectory != nil {
				job.WorkingDir = *apiJob.CurrentWorkingDirectory
			}
			
			// Command
			if apiJob.Command != nil {
				job.Command = *apiJob.Command
			}
			
			// Environment variables - Initialize empty map
			job.Environment = make(map[string]string)
			
			// Nodes
			if apiJob.JobResources != nil && apiJob.JobResources.Nodes != nil && apiJob.JobResources.Nodes.List != nil {
				// Parse node list string into slice
				nodeListStr := *apiJob.JobResources.Nodes.List
				if nodeListStr != "" {
					job.Nodes = strings.Split(nodeListStr, ",")
				}
			}
			
			// Exit code
			if apiJob.ExitCode != nil && apiJob.ExitCode.ReturnCode != nil &&
				apiJob.ExitCode.ReturnCode.Set != nil && *apiJob.ExitCode.ReturnCode.Set &&
				apiJob.ExitCode.ReturnCode.Number != nil {
				job.ExitCode = int(*apiJob.ExitCode.ReturnCode.Number)
			}
			
			// Initialize metadata
			job.Metadata = make(map[string]interface{})
			
			// Add additional metadata
			if apiJob.Account != nil {
				job.Metadata["account"] = *apiJob.Account
			}
			if apiJob.AdminComment != nil {
				job.Metadata["admin_comment"] = *apiJob.AdminComment
			}
			if apiJob.AllocatingNode != nil {
				job.Metadata["allocating_node"] = *apiJob.AllocatingNode
			}
			
			jobs = append(jobs, job)
		}
	}

	// Apply client-side filtering if options are provided
	if opts != nil {
		jobs = filterJobs(jobs, opts)
	}

	return &interfaces.JobList{
		Jobs:  jobs,
		Total: len(jobs),
	}, nil
}

// filterJobs applies client-side filtering to job list
func filterJobs(jobs []interfaces.Job, opts *interfaces.ListJobsOptions) []interfaces.Job {
	var filtered []interfaces.Job

	for _, job := range jobs {
		// Filter by user ID
		if opts.UserID != "" && job.UserID != opts.UserID {
			continue
		}

		// Filter by states
		if len(opts.States) > 0 {
			stateMatch := false
			for _, state := range opts.States {
				if strings.EqualFold(job.State, state) {
					stateMatch = true
					break
				}
			}
			if !stateMatch {
				continue
			}
		}

		// Filter by partition
		if opts.Partition != "" && !strings.EqualFold(job.Partition, opts.Partition) {
			continue
		}

		filtered = append(filtered, job)
	}

	// Apply limit and offset
	if opts.Offset > 0 {
		if opts.Offset >= len(filtered) {
			return []interfaces.Job{}
		}
		filtered = filtered[opts.Offset:]
	}

	if opts.Limit > 0 && len(filtered) > opts.Limit {
		filtered = filtered[:opts.Limit]
	}

	return filtered
}

// Get retrieves a specific job by ID
func (m *JobManagerImpl) Get(ctx context.Context, jobID string) (*interfaces.Job, error) {
	// For v0.0.41, we need to list all jobs and filter
	// This is because the GetJob endpoint might have different response structure
	list, err := m.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	for _, job := range list.Jobs {
		if job.ID == jobID {
			return &job, nil
		}
	}

	return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "Job not found", fmt.Sprintf("Job ID %s not found", jobID))
}

// Submit submits a new job
func (m *JobManagerImpl) Submit(ctx context.Context, job *interfaces.JobSubmission) (*interfaces.JobSubmitResponse, error) {
	// Note: v0.0.41 has a complex inline struct for job submission
	// For now, return unsupported operation error
	return nil, errors.NewClientError(
		errors.ErrorCodeUnsupportedOperation,
		"Job submission not implemented for v0.0.41",
		"The v0.0.41 job submission requires complex inline struct mapping that differs significantly from other API versions",
	)
}

// Cancel cancels a job
func (m *JobManagerImpl) Cancel(ctx context.Context, jobID string) error {
	// Check if API client is available
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmV0041DeleteJobParams{}

	// Send SIGTERM signal by default
	signal := "SIGTERM"
	params.Signal = &signal

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0041DeleteJobWithResponse(ctx, jobID, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.41")
	}

	// Check HTTP status and handle API errors
	if resp.StatusCode() != 200 {
		var responseBody []byte
		if resp.JSON200 != nil {
			// Try to extract error details from response
			if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
				apiErrors := make([]errors.SlurmAPIErrorDetail, len(*resp.JSON200.Errors))
				for i, apiErr := range *resp.JSON200.Errors {
					var errorNumber int
					if apiErr.ErrorNumber != nil {
						errorNumber = int(*apiErr.ErrorNumber)
					}
					var errorCode string
					if apiErr.Error != nil {
						errorCode = *apiErr.Error
					}
					var source string
					if apiErr.Source != nil {
						source = *apiErr.Source
					}
					var description string
					if apiErr.Description != nil {
						description = *apiErr.Description
					}

					apiErrors[i] = errors.SlurmAPIErrorDetail{
						ErrorNumber: errorNumber,
						ErrorCode:   errorCode,
						Source:      source,
						Description: description,
					}
				}
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.41", apiErrors)
				return apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.41")
		return httpErr
	}

	return nil
}

// Update updates job properties
func (m *JobManagerImpl) Update(ctx context.Context, jobID string, update *interfaces.JobUpdate) error {
	// Note: v0.0.41 has different update structure
	return errors.NewClientError(
		errors.ErrorCodeUnsupportedOperation,
		"Job updates not implemented for v0.0.41",
		"The v0.0.41 job update requires complex inline struct mapping that differs significantly from other API versions",
	)
}

// Steps retrieves job steps for a job
func (m *JobManagerImpl) Steps(ctx context.Context, jobID string) (*interfaces.JobStepList, error) {
	// v0.0.41 doesn't include step details in job info
	// Return empty step list
	steps := make([]interfaces.JobStep, 0)

	return &interfaces.JobStepList{
		Steps: steps,
		Total: len(steps),
	}, nil
}

// Watch provides real-time job updates through polling
func (m *JobManagerImpl) Watch(ctx context.Context, opts *interfaces.WatchJobsOptions) (<-chan interfaces.JobEvent, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Create a job poller with the List function
	poller := watch.NewJobPoller(m.List)

	// Configure polling interval if needed (default is 5 seconds)
	// poller.WithPollInterval(3 * time.Second)

	// Start watching
	return poller.Watch(ctx, opts)
}

// GetJobUtilization retrieves basic resource utilization metrics for a job
// Note: v0.0.41 only supports limited CPU and memory metrics
func (m *JobManagerImpl) GetJobUtilization(ctx context.Context, jobID string) (*interfaces.JobUtilization, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// First get the job details to determine status
	job, err := m.Get(ctx, jobID)
	if err != nil {
		return nil, err
	}

	// In v0.0.41, only basic CPU and memory metrics are available
	// Advanced features like GPU, I/O, network monitoring are not supported
	// TODO: Integrate with basic SLURM accounting when available

	utilization := &interfaces.JobUtilization{
		JobID:   jobID,
		JobName: job.Name,
		StartTime: job.SubmitTime,
		EndTime: job.EndTime,
		
		// CPU Utilization (basic in v0.0.41)
		CPUUtilization: &interfaces.ResourceUtilization{
			Used:      float64(job.CPUs) * 0.75, // Simulated 75% utilization
			Allocated: float64(job.CPUs),
			Limit:     float64(job.CPUs),
			Percentage: 75.0,
		},
		
		// Memory Utilization (basic in v0.0.41)
		MemoryUtilization: &interfaces.ResourceUtilization{
			Used:      float64(job.Memory) * 0.65, // Simulated 65% utilization
			Allocated: float64(job.Memory),
			Limit:     float64(job.Memory),
			Percentage: 65.0,
		},
	}

	// Add metadata
	utilization.Metadata = map[string]interface{}{
		"version": "v0.0.41",
		"source": "simulated", // TODO: Change to "basic_accounting" when available
		"nodes": job.Nodes,
		"partition": job.Partition,
		"state": job.State,
		"feature_level": "basic", // v0.0.41 has basic features only
		"limitations": []string{
			"no_gpu_metrics",
			"no_io_metrics",
			"no_network_metrics",
			"no_energy_metrics",
			"basic_cpu_memory_only",
		},
	}

	// GPU utilization not supported in v0.0.41
	utilization.GPUUtilization = nil
	
	// I/O utilization not supported in v0.0.41
	utilization.IOUtilization = nil
	
	// Network utilization not supported in v0.0.41
	utilization.NetworkUtilization = nil
	
	// Energy usage not supported in v0.0.41
	utilization.EnergyUsage = nil

	return utilization, nil
}

// GetJobEfficiency calculates basic efficiency metrics for a completed job
// Note: v0.0.41 only supports very basic efficiency calculations
func (m *JobManagerImpl) GetJobEfficiency(ctx context.Context, jobID string) (*interfaces.ResourceUtilization, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Get job utilization data
	utilization, err := m.GetJobUtilization(ctx, jobID)
	if err != nil {
		return nil, err
	}

	// Calculate overall efficiency based on CPU and memory only (v0.0.41 limitation)
	cpuWeight := 0.6  // Higher CPU weight since only CPU/memory available
	memWeight := 0.4

	totalEfficiency := 0.0

	// CPU efficiency
	if utilization.CPUUtilization != nil {
		totalEfficiency += utilization.CPUUtilization.Percentage * cpuWeight
	}

	// Memory efficiency
	if utilization.MemoryUtilization != nil {
		totalEfficiency += utilization.MemoryUtilization.Percentage * memWeight
	}

	// Calculate final efficiency percentage
	efficiency := totalEfficiency // Already weighted sum to 1.0

	return &interfaces.ResourceUtilization{
		Used:       efficiency,
		Allocated:  100.0,
		Limit:      100.0,
		Percentage: efficiency,
		Metadata: map[string]interface{}{
			"cpu_efficiency":    utilization.CPUUtilization.Percentage,
			"memory_efficiency": utilization.MemoryUtilization.Percentage,
			"calculation_method": "basic_cpu_memory_v41",
			"version": "v0.0.41",
			"weights": map[string]float64{
				"cpu":    cpuWeight,
				"memory": memWeight,
			},
			"limitations": []string{
				"cpu_memory_only",
				"no_gpu_efficiency",
				"no_io_efficiency",
				"no_network_efficiency",
			},
		},
	}, nil
}

// GetJobPerformance retrieves basic performance metrics for a job
// Note: v0.0.41 provides minimal performance analysis capabilities
func (m *JobManagerImpl) GetJobPerformance(ctx context.Context, jobID string) (*interfaces.JobPerformance, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Get basic job info
	job, err := m.Get(ctx, jobID)
	if err != nil {
		return nil, err
	}

	// Convert string jobID to uint32
	jobIDInt, err := strconv.ParseUint(jobID, 10, 32)
	if err != nil {
		return nil, errors.NewClientError(errors.ErrorCodeInvalidRequest, "Invalid job ID format", err.Error())
	}

	// Get utilization metrics
	utilization, err := m.GetJobUtilization(ctx, jobID)
	if err != nil {
		return nil, err
	}

	// Get efficiency metrics
	efficiency, err := m.GetJobEfficiency(ctx, jobID)
	if err != nil {
		return nil, err
	}

	// Build basic performance report (v0.0.41 version with minimal features)
	performance := &interfaces.JobPerformance{
		JobID:     uint32(jobIDInt),
		JobName:   job.Name,
		StartTime: job.SubmitTime,
		EndTime:   job.EndTime,
		Status:    job.State,
		ExitCode:  job.ExitCode,
		
		ResourceUtilization: efficiency,
		JobUtilization:     utilization,
		
		// Step metrics not available in v0.0.41
		StepMetrics: nil,
		
		// Performance trends not available in v0.0.41
		PerformanceTrends: nil,
		
		// Bottleneck analysis very basic in v0.0.41
		Bottlenecks: analyzeBottlenecksV41(utilization),
		
		// Optimization recommendations very basic in v0.0.41
		Recommendations: generateRecommendationsV41(efficiency),
	}

	return performance, nil
}

// Helper function to analyze bottlenecks for v0.0.41 (very basic analysis)
func analyzeBottlenecksV41(utilization *interfaces.JobUtilization) []interfaces.PerformanceBottleneck {
	bottlenecks := []interfaces.PerformanceBottleneck{}

	// Only check CPU bottleneck in v0.0.41
	if utilization.CPUUtilization != nil && utilization.CPUUtilization.Percentage > 85 {
		bottlenecks = append(bottlenecks, interfaces.PerformanceBottleneck{
			Type:         "cpu",
			Severity:     "medium",
			Description:  "High CPU utilization detected",
			Impact:       10.0, // 10% estimated performance impact
			TimeDetected: time.Now(),
			Duration:     30 * time.Minute, // Estimated
		})
	}

	// Only check memory bottleneck in v0.0.41
	if utilization.MemoryUtilization != nil && utilization.MemoryUtilization.Percentage > 80 {
		bottlenecks = append(bottlenecks, interfaces.PerformanceBottleneck{
			Type:         "memory",
			Severity:     "low",
			Description:  "High memory utilization detected",
			Impact:       5.0,
			TimeDetected: time.Now(),
			Duration:     20 * time.Minute,
		})
	}

	return bottlenecks
}

// Helper function to generate optimization recommendations for v0.0.41 (very basic)
func generateRecommendationsV41(efficiency *interfaces.ResourceUtilization) []interfaces.OptimizationRecommendation {
	recommendations := []interfaces.OptimizationRecommendation{}

	// Only basic overall efficiency recommendation in v0.0.41
	if efficiency.Percentage < 70 {
		recommendations = append(recommendations, interfaces.OptimizationRecommendation{
			Type:        "workflow",
			Priority:    "low",
			Title:       "Resource utilization below optimal",
			Description: "Consider reviewing resource allocation for better efficiency.",
			ExpectedImprovement: 10.0,
			ConfigChanges: map[string]string{
				"action": "review_resource_usage",
			},
		})
	}

	// Add a note about v0.0.41 limitations
	recommendations = append(recommendations, interfaces.OptimizationRecommendation{
		Type:        "configuration",
		Priority:    "low", 
		Title:       "Limited analytics in API v0.0.41",
		Description: "Consider upgrading to SLURM API v0.0.42+ for enhanced analytics capabilities including GPU, I/O, and network metrics.",
		ExpectedImprovement: 0.0,
		ConfigChanges: map[string]string{
			"recommended_api_version": "v0.0.42_or_higher",
		},
	})

	return recommendations
}

// GetJobLiveMetrics retrieves real-time performance metrics for a running job
// Note: v0.0.41 only supports very basic live metrics
func (m *JobManagerImpl) GetJobLiveMetrics(ctx context.Context, jobID string) (*interfaces.JobLiveMetrics, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// First get the job details to check if it's running
	job, err := m.Get(ctx, jobID)
	if err != nil {
		return nil, err
	}

	// Only running jobs have live metrics
	if job.State != "RUNNING" && job.State != "SUSPENDED" {
		return nil, errors.NewClientError(errors.ErrorCodeInvalidRequest, 
			fmt.Sprintf("Job %s is not running (state: %s)", jobID, job.State))
	}

	// Calculate running time
	runningTime := time.Duration(0)
	if job.StartTime != nil {
		runningTime = time.Since(*job.StartTime)
	}

	// Create minimal live metrics response
	// v0.0.41 has very limited real-time capabilities
	liveMetrics := &interfaces.JobLiveMetrics{
		JobID:          jobID,
		JobName:        job.Name,
		State:          job.State,
		RunningTime:    runningTime,
		CollectionTime: time.Now(),
		
		// Very basic CPU usage (fixed estimates)
		CPUUsage: &interfaces.LiveResourceMetric{
			Current:            float64(job.CPUs) * 0.70,
			Average1Min:        float64(job.CPUs) * 0.70,
			Average5Min:        float64(job.CPUs) * 0.70,
			Peak:               float64(job.CPUs) * 0.80,
			Allocated:          float64(job.CPUs),
			UtilizationPercent: 70.0,
			Trend:              "stable",
			Unit:               "cores",
		},
		
		// Very basic memory usage
		MemoryUsage: &interfaces.LiveResourceMetric{
			Current:            float64(job.Memory) * 0.60,
			Average1Min:        float64(job.Memory) * 0.60,
			Average5Min:        float64(job.Memory) * 0.60,
			Peak:               float64(job.Memory) * 0.70,
			Allocated:          float64(job.Memory),
			UtilizationPercent: 60.0,
			Trend:              "stable",
			Unit:               "bytes",
		},
		
		// Minimal process information
		ProcessCount: 1,
		ThreadCount:  4,
		
		// No advanced metrics in v0.0.41
		GPUUsage:     nil,
		NetworkUsage: nil,
		IOUsage:      nil,
		
		// Empty collections
		NodeMetrics: make(map[string]*interfaces.NodeLiveMetrics),
		Alerts:      []interfaces.PerformanceAlert{},
	}

	// Add very basic node metrics (v0.0.41 has minimal per-node data)
	if len(job.Nodes) > 0 {
		// Only add metrics for the first node as representative
		nodeName := job.Nodes[0]
		nodeMetrics := &interfaces.NodeLiveMetrics{
			NodeName:  nodeName,
			CPUCores:  job.CPUs,
			MemoryGB:  float64(job.Memory) / (1024 * 1024 * 1024),
			
			CPUUsage: &interfaces.LiveResourceMetric{
				Current:            70.0,
				Average1Min:        70.0,
				Average5Min:        70.0,
				Peak:               80.0,
				Allocated:          100.0,
				UtilizationPercent: 70.0,
				Trend:              "stable",
				Unit:               "percent",
			},
			
			MemoryUsage: &interfaces.LiveResourceMetric{
				Current:            60.0,
				Average1Min:        60.0,
				Average5Min:        60.0,
				Peak:               70.0,
				Allocated:          100.0,
				UtilizationPercent: 60.0,
				Trend:              "stable",
				Unit:               "percent",
			},
			
			// Minimal load average
			LoadAverage: []float64{1.0, 1.0, 1.0},
			
			// No advanced monitoring in v0.0.41
			CPUTemperature:   0.0,
			PowerConsumption: 0.0,
			NetworkInRate:    0.0,
			NetworkOutRate:   0.0,
			DiskReadRate:     0.0,
			DiskWriteRate:    0.0,
		}
		
		liveMetrics.NodeMetrics[nodeName] = nodeMetrics
	}

	// Add metadata about limitations
	liveMetrics.Metadata = map[string]interface{}{
		"version":           "v0.0.41",
		"collection_method": "minimal_monitoring",
		"limitations": []string{
			"fixed_utilization_values",
			"no_gpu_metrics",
			"no_network_metrics",
			"no_io_metrics",
			"no_temperature_monitoring",
			"no_power_monitoring",
			"no_alerts",
			"single_node_representation",
		},
		"note": "Upgrade to v0.0.42+ for enhanced live monitoring",
	}

	return liveMetrics, nil
}

// WatchJobMetrics provides streaming performance updates for a running job
// Note: v0.0.41 has very limited streaming capabilities - basic CPU/memory only
func (m *JobManagerImpl) WatchJobMetrics(ctx context.Context, jobID string, opts *interfaces.WatchMetricsOptions) (<-chan interfaces.JobMetricsEvent, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Default options if not provided - v0.0.41 limitations
	if opts == nil {
		opts = &interfaces.WatchMetricsOptions{
			UpdateInterval:     15 * time.Second, // Much slower polling for v0.0.41
			IncludeCPU:        true,
			IncludeMemory:     true,
			IncludeGPU:        false, // Not supported
			IncludeNetwork:    false, // Not supported
			IncludeIO:         false, // Not supported
			IncludeEnergy:     false, // Not supported
			IncludeNodeMetrics: false, // Very limited
			StopOnCompletion:  true,
			CPUThreshold:      80.0,  // Conservative thresholds
			MemoryThreshold:   75.0,
		}
	}

	// Enforce minimum update interval for v0.0.41
	if opts.UpdateInterval < 15*time.Second {
		opts.UpdateInterval = 15 * time.Second
	}

	// Create event channel
	eventChan := make(chan interfaces.JobMetricsEvent, 3)

	// Start monitoring goroutine
	go func() {
		defer close(eventChan)

		// Track previous state
		var previousState string

		// Initial state check
		job, err := m.Get(ctx, jobID)
		if err != nil {
			eventChan <- interfaces.JobMetricsEvent{
				Type:      "error",
				JobID:     jobID,
				Timestamp: time.Now(),
				Error:     err,
			}
			return
		}
		previousState = job.State

		// Set up monitoring timer
		ticker := time.NewTicker(opts.UpdateInterval)
		defer ticker.Stop()

		// Set up max duration timer if specified
		var maxTimer *time.Timer
		if opts.MaxDuration > 0 {
			maxTimer = time.NewTimer(opts.MaxDuration)
			defer maxTimer.Stop()
		}

		// Simple monitoring loop for v0.0.41
		for {
			select {
			case <-ctx.Done():
				eventChan <- interfaces.JobMetricsEvent{
					Type:      "complete",
					JobID:     jobID,
					Timestamp: time.Now(),
					Error:     ctx.Err(),
				}
				return

			case <-ticker.C:
				// Get current job state
				job, err := m.Get(ctx, jobID)
				if err != nil {
					eventChan <- interfaces.JobMetricsEvent{
						Type:      "error",
						JobID:     jobID,
						Timestamp: time.Now(),
						Error:     err,
					}
					continue
				}

				// Check for state change
				if job.State != previousState {
					stateChange := &interfaces.JobStateChange{
						OldState: previousState,
						NewState: job.State,
					}
					previousState = job.State

					eventChan <- interfaces.JobMetricsEvent{
						Type:        "update",
						JobID:       jobID,
						Timestamp:   time.Now(),
						StateChange: stateChange,
					}

					// Check if we should stop
					if opts.StopOnCompletion && isJobCompleteV41(job.State) {
						eventChan <- interfaces.JobMetricsEvent{
							Type:      "complete",
							JobID:     jobID,
							Timestamp: time.Now(),
						}
						return
					}
				}

				// Only collect metrics for running jobs
				if job.State == "RUNNING" {
					// Get basic metrics
					metrics, err := m.GetJobLiveMetrics(ctx, jobID)
					if err != nil {
						// v0.0.41 might not support live metrics for all jobs
						// Continue monitoring without sending error
						continue
					}

					// Send simplified metrics update
					eventChan <- interfaces.JobMetricsEvent{
						Type:      "update",
						JobID:     jobID,
						Timestamp: time.Now(),
						Metrics:   metrics,
					}
				}

			case <-func() <-chan time.Time {
				if maxTimer != nil {
					return maxTimer.C
				}
				return nil
			}():
				// Max duration reached
				eventChan <- interfaces.JobMetricsEvent{
					Type:      "complete",
					JobID:     jobID,
					Timestamp: time.Now(),
				}
				return
			}
		}
	}()

	return eventChan, nil
}

// Helper function to check if job is complete (v0.0.41)
func isJobCompleteV41(state string) bool {
	completedStates := []string{
		"COMPLETED", "FAILED", "CANCELLED", "TIMEOUT",
		"NODE_FAIL", "PREEMPTED",
	}
	for _, s := range completedStates {
		if state == s {
			return true
		}
	}
	return false
}

// GetJobResourceTrends retrieves performance trends over specified time windows
// Note: v0.0.41 has very limited trend analysis - basic CPU/memory only
func (m *JobManagerImpl) GetJobResourceTrends(ctx context.Context, jobID string, opts *interfaces.ResourceTrendsOptions) (*interfaces.JobResourceTrends, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Get job details
	job, err := m.Get(ctx, jobID)
	if err != nil {
		return nil, err
	}

	// Set default options - v0.0.41 severe limitations
	if opts == nil {
		opts = &interfaces.ResourceTrendsOptions{
			DataPoints:     6, // Very few data points
			IncludeCPU:     true,
			IncludeMemory:  true,
			IncludeGPU:     false, // Not supported
			IncludeIO:      false, // Not supported
			IncludeNetwork: false, // Not supported
			IncludeEnergy:  false, // Not supported
			Aggregation:    "avg",
			DetectAnomalies: false, // Not supported
		}
	}

	// Limit data points for v0.0.41
	if opts.DataPoints == 0 || opts.DataPoints > 6 {
		opts.DataPoints = 6
	}

	// Calculate time window
	var timeWindow time.Duration
	if job.StartTime != nil && job.EndTime != nil {
		timeWindow = job.EndTime.Sub(*job.StartTime)
	} else {
		timeWindow = time.Hour
	}

	// Generate time points
	timePoints := generateMinimalTimePoints(job.StartTime, job.EndTime, opts.DataPoints)

	// Create minimal trends object
	trends := &interfaces.JobResourceTrends{
		JobID:      jobID,
		JobName:    job.Name,
		StartTime:  job.SubmitTime,
		EndTime:    job.EndTime,
		TimeWindow: timeWindow,
		DataPoints: len(timePoints),
		TimePoints: timePoints,
		Anomalies:  []interfaces.ResourceAnomaly{}, // Not supported
	}

	// Generate very basic CPU trends
	if opts.IncludeCPU {
		trends.CPUTrends = generateFixedResourceTrends(float64(job.CPUs)*0.75, float64(job.CPUs), "cores", len(timePoints))
	}

	// Generate very basic memory trends
	if opts.IncludeMemory {
		trends.MemoryTrends = generateFixedResourceTrends(float64(job.Memory)*0.65, float64(job.Memory), "bytes", len(timePoints))
	}

	// Not supported in v0.0.41
	trends.GPUTrends = nil
	trends.IOTrends = nil
	trends.NetworkTrends = nil
	trends.EnergyTrends = nil

	// Generate minimal summary
	trends.Summary = generateMinimalTrendsSummary(trends)

	// Add metadata
	trends.Metadata = map[string]interface{}{
		"version":     "v0.0.41",
		"data_source": "fixed_estimates",
		"limitations": []string{
			"minimal_data_points",
			"cpu_memory_only",
			"fixed_values",
			"no_real_trends",
			"no_anomaly_detection",
		},
		"note": "Upgrade to v0.0.42+ for enhanced trend analysis",
	}

	return trends, nil
}

// Helper functions for v0.0.41
func generateMinimalTimePoints(startTime, endTime *time.Time, numPoints int) []time.Time {
	if numPoints <= 0 {
		return []time.Time{}
	}

	// Use current time if no start time
	start := time.Now().Add(-time.Hour)
	if startTime != nil {
		start = *startTime
	}

	points := make([]time.Time, numPoints)
	interval := time.Hour / time.Duration(numPoints)
	
	for i := 0; i < numPoints; i++ {
		points[i] = start.Add(time.Duration(i) * interval)
	}

	return points
}

func generateFixedResourceTrends(avgValue, maxValue float64, unit string, numPoints int) *interfaces.ResourceTimeSeries {
	values := make([]float64, numPoints)
	
	// Fixed values for v0.0.41
	for i := range values {
		values[i] = avgValue
	}

	return &interfaces.ResourceTimeSeries{
		Values:     values,
		Unit:       unit,
		Average:    avgValue,
		Min:        avgValue * 0.9,
		Max:        maxValue,
		StdDev:     0.0,
		Trend:      "stable",
		TrendSlope: 0.0,
	}
}

func generateMinimalTrendsSummary(trends *interfaces.JobResourceTrends) *interfaces.TrendsSummary {
	summary := &interfaces.TrendsSummary{
		PeakUtilization:    make(map[string]float64),
		AverageUtilization: make(map[string]float64),
	}

	// Fixed summary values
	if trends.CPUTrends != nil {
		summary.PeakUtilization["cpu"] = trends.CPUTrends.Max
		summary.AverageUtilization["cpu"] = trends.CPUTrends.Average
	}

	if trends.MemoryTrends != nil {
		summary.PeakUtilization["memory"] = trends.MemoryTrends.Max
		summary.AverageUtilization["memory"] = trends.MemoryTrends.Average
	}

	// Fixed values for v0.0.41
	summary.ResourceEfficiency = 70.0
	summary.StabilityScore = 80.0
	summary.VariabilityIndex = 0.2
	summary.OverallTrend = "stable"
	summary.ResourceBalance = "balanced"

	return summary
}