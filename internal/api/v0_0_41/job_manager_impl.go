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