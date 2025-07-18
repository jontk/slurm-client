package v0_0_40

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
	params := &SlurmV0040GetJobsParams{}

	// Set flags to get detailed job information
	flags := SlurmV0040GetJobsParamsFlagsDETAIL
	params.Flags = &flags

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0040GetJobsWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.40")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.40", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.40")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response but got nil")
	}

	// Convert the response to our interface types
	jobs := make([]interfaces.Job, 0, len(resp.JSON200.Jobs))
	for _, apiJob := range resp.JSON200.Jobs {
		job, err := convertAPIJobToInterface(apiJob)
		if err != nil {
			conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert job data")
			conversionErr.Cause = err
			conversionErr.Details = fmt.Sprintf("Error converting job ID %v", apiJob.JobId)
			return nil, conversionErr
		}
		jobs = append(jobs, *job)
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

// convertAPIJobToInterface converts a V0040JobInfo to interfaces.Job
func convertAPIJobToInterface(apiJob V0040JobInfo) (*interfaces.Job, error) {
	job := &interfaces.Job{}

	// Job ID - simple int32 pointer
	if apiJob.JobId != nil {
		job.ID = strconv.FormatInt(int64(*apiJob.JobId), 10)
	}

	// Job name
	if apiJob.Name != nil {
		job.Name = *apiJob.Name
	}

	// User ID - simple int32 pointer
	if apiJob.UserId != nil {
		job.UserID = strconv.FormatInt(int64(*apiJob.UserId), 10)
	}

	// Group ID - simple int32 pointer
	if apiJob.GroupId != nil {
		job.GroupID = strconv.FormatInt(int64(*apiJob.GroupId), 10)
	}

	// Job state - array of strings
	if apiJob.JobState != nil && len(*apiJob.JobState) > 0 {
		job.State = (*apiJob.JobState)[0]
	}

	// Partition
	if apiJob.Partition != nil {
		job.Partition = *apiJob.Partition
	}

	// Priority - NoValStruct
	if apiJob.Priority != nil && apiJob.Priority.Set != nil && *apiJob.Priority.Set && apiJob.Priority.Number != nil {
		job.Priority = int(*apiJob.Priority.Number)
	}

	// Submit time - NoValStruct
	if apiJob.SubmitTime != nil && apiJob.SubmitTime.Set != nil && *apiJob.SubmitTime.Set && apiJob.SubmitTime.Number != nil {
		job.SubmitTime = time.Unix(*apiJob.SubmitTime.Number, 0)
	}

	// Start time - NoValStruct
	if apiJob.StartTime != nil && apiJob.StartTime.Set != nil && *apiJob.StartTime.Set && apiJob.StartTime.Number != nil && *apiJob.StartTime.Number > 0 {
		startTime := time.Unix(*apiJob.StartTime.Number, 0)
		job.StartTime = &startTime
	}

	// End time - NoValStruct
	if apiJob.EndTime != nil && apiJob.EndTime.Set != nil && *apiJob.EndTime.Set && apiJob.EndTime.Number != nil && *apiJob.EndTime.Number > 0 {
		endTime := time.Unix(*apiJob.EndTime.Number, 0)
		job.EndTime = &endTime
	}

	// CPUs - NoValStruct
	if apiJob.Cpus != nil && apiJob.Cpus.Set != nil && *apiJob.Cpus.Set && apiJob.Cpus.Number != nil {
		job.CPUs = int(*apiJob.Cpus.Number)
	}

	// Memory (convert MB to bytes for consistency) - NoValStruct
	if apiJob.MemoryPerNode != nil && apiJob.MemoryPerNode.Set != nil && *apiJob.MemoryPerNode.Set && apiJob.MemoryPerNode.Number != nil {
		job.Memory = int(*apiJob.MemoryPerNode.Number) * 1024 * 1024
	}

	// Time limit (in minutes) - NoValStruct
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

	// Environment variables - Initialize empty map since not directly available in JobInfo
	job.Environment = make(map[string]string)

	// Nodes - Extract from JobResources
	if apiJob.JobResources != nil && apiJob.JobResources.Nodes != nil {
		// Parse node list string into slice
		nodeListStr := *apiJob.JobResources.Nodes
		if nodeListStr != "" {
			// Simple splitting by comma - real implementation might need more sophisticated parsing
			job.Nodes = strings.Split(nodeListStr, ",")
		}
	}

	// Exit code - ProcessExitCodeVerbose
	if apiJob.ExitCode != nil && apiJob.ExitCode.ReturnCode != nil &&
		apiJob.ExitCode.ReturnCode.Set != nil && *apiJob.ExitCode.ReturnCode.Set &&
		apiJob.ExitCode.ReturnCode.Number != nil {
		job.ExitCode = int(*apiJob.ExitCode.ReturnCode.Number)
	}

	// Initialize metadata
	job.Metadata = make(map[string]interface{})

	// Add additional metadata from API response
	if apiJob.Account != nil {
		job.Metadata["account"] = *apiJob.Account
	}
	if apiJob.AdminComment != nil {
		job.Metadata["admin_comment"] = *apiJob.AdminComment
	}
	if apiJob.AllocatingNode != nil {
		job.Metadata["allocating_node"] = *apiJob.AllocatingNode
	}

	return job, nil
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
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmV0040GetJobParams{}

	// Set flags to get detailed job information
	flags := SlurmV0040GetJobParamsFlagsDETAIL
	params.Flags = &flags

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0040GetJobWithResponse(ctx, jobID, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.40")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.40", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.40")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response but got nil")
	}

	// Convert the response to our interface types
	if len(resp.JSON200.Jobs) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "Job not found", fmt.Sprintf("Job ID %s not found", jobID))
	}

	if len(resp.JSON200.Jobs) > 1 {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected multiple jobs returned", fmt.Sprintf("Expected 1 job but got %d for ID %s", len(resp.JSON200.Jobs), jobID))
	}

	job, err := convertAPIJobToInterface(resp.JSON200.Jobs[0])
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert job data")
		conversionErr.Cause = err
		conversionErr.Details = fmt.Sprintf("Error converting job ID %s", jobID)
		return nil, conversionErr
	}

	return job, nil
}

// Submit submits a new job
func (m *JobManagerImpl) Submit(ctx context.Context, job *interfaces.JobSubmission) (*interfaces.JobSubmitResponse, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Convert interface JobSubmission to API JobDescMsg
	jobDesc, err := convertJobSubmissionToAPI(job)
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeInvalidRequest, "Failed to convert job submission")
		conversionErr.Cause = err
		conversionErr.Details = "Error converting JobSubmission to API format"
		return nil, conversionErr
	}

	// Create the request body
	requestBody := SlurmV0040PostJobSubmitJSONRequestBody{
		Job: jobDesc,
	}

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0040PostJobSubmitWithResponse(ctx, requestBody)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.40")
	}

	// Check HTTP status (200 and 201 for creation is success)
	if resp.StatusCode() != 200 && resp.StatusCode() != 201 {
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.40", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.40")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response but got nil")
	}

	// Convert response to interface type
	result := &interfaces.JobSubmitResponse{}
	if resp.JSON200.JobId != nil {
		result.JobID = strconv.FormatInt(int64(*resp.JSON200.JobId), 10)
	} else {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Job submission successful but no job ID returned")
	}

	return result, nil
}

// convertJobSubmissionToAPI converts interfaces.JobSubmission to V0040JobDescMsg
func convertJobSubmissionToAPI(job *interfaces.JobSubmission) (*V0040JobDescMsg, error) {
	if job == nil {
		return nil, errors.NewClientError(errors.ErrorCodeInvalidRequest, "Job submission cannot be nil")
	}

	jobDesc := &V0040JobDescMsg{}

	// Basic fields
	if job.Name != "" {
		jobDesc.Name = &job.Name
	}

	if job.Script != "" {
		jobDesc.Script = &job.Script
	}

	if job.Partition != "" {
		jobDesc.Partition = &job.Partition
	}

	if job.WorkingDir != "" {
		jobDesc.CurrentWorkingDirectory = &job.WorkingDir
	}

	// Resource requirements
	if job.CPUs > 0 {
		cpus := int32(job.CPUs)
		jobDesc.MinimumCpus = &cpus
	}

	if job.Memory > 0 {
		// Convert bytes to MB (Slurm expects MB)
		memoryMB := int64(job.Memory / (1024 * 1024))
		set := true
		jobDesc.MemoryPerNode = &V0040Uint64NoVal{
			Number: &memoryMB,
			Set:    &set,
		}
	}

	if job.TimeLimit > 0 {
		timeLimit := int64(job.TimeLimit)
		set := true
		jobDesc.TimeLimit = &V0040Uint32NoVal{
			Number: &timeLimit,
			Set:    &set,
		}
	}

	if job.Nodes > 0 {
		nodes := int32(job.Nodes)
		jobDesc.MinimumNodes = &nodes
	}

	if job.Priority > 0 {
		priority := int64(job.Priority)
		set := true
		jobDesc.Priority = &V0040Uint32NoVal{
			Number: &priority,
			Set:    &set,
		}
	}

	// Environment variables
	if len(job.Environment) > 0 {
		envVars := make([]string, 0, len(job.Environment))
		for key, value := range job.Environment {
			envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
		}
		jobDesc.Environment = &envVars
	}

	// Args
	if len(job.Args) > 0 {
		jobDesc.Argv = &job.Args
	}

	return jobDesc, nil
}

// Cancel cancels a job
func (m *JobManagerImpl) Cancel(ctx context.Context, jobID string) error {
	// Check if API client is available
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmV0040DeleteJobParams{}

	// Send SIGTERM signal by default (can be made configurable later)
	signal := "SIGTERM"
	params.Signal = &signal

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0040DeleteJobWithResponse(ctx, jobID, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.40")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.40", apiErrors)
				return apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.40")
		return httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil {
		return errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response but got nil")
	}

	return nil
}

// Update updates job properties
func (m *JobManagerImpl) Update(ctx context.Context, jobID string, update *interfaces.JobUpdate) error {
	// Check if API client is available
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Validate inputs
	if update == nil {
		return errors.NewClientError(errors.ErrorCodeInvalidRequest, "Job update cannot be nil")
	}

	// Convert interface JobUpdate to API JobDescMsg
	jobDesc, err := convertJobUpdateToAPI(update)
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeInvalidRequest, "Failed to convert job update")
		conversionErr.Cause = err
		conversionErr.Details = "Error converting JobUpdate to API format"
		return conversionErr
	}

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0040PostJobWithResponse(ctx, jobID, *jobDesc)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.40")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.40", apiErrors)
				return apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.40")
		return httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil {
		return errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response but got nil")
	}

	return nil
}

// convertJobUpdateToAPI converts interfaces.JobUpdate to V0040JobDescMsg
func convertJobUpdateToAPI(update *interfaces.JobUpdate) (*V0040JobDescMsg, error) {
	if update == nil {
		return nil, errors.NewClientError(errors.ErrorCodeInvalidRequest, "Job update cannot be nil")
	}

	jobDesc := &V0040JobDescMsg{}

	// Only include fields that are actually being updated (non-nil values)
	if update.Priority != nil {
		priority := int64(*update.Priority)
		set := true
		jobDesc.Priority = &V0040Uint32NoVal{
			Number: &priority,
			Set:    &set,
		}
	}

	if update.TimeLimit != nil {
		timeLimit := int64(*update.TimeLimit)
		set := true
		jobDesc.TimeLimit = &V0040Uint32NoVal{
			Number: &timeLimit,
			Set:    &set,
		}
	}

	if update.Name != nil {
		jobDesc.Name = update.Name
	}

	return jobDesc, nil
}

// Steps retrieves job steps for a job
func (m *JobManagerImpl) Steps(ctx context.Context, jobID string) (*interfaces.JobStepList, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmV0040GetJobParams{}

	// Set flags to get detailed job information including steps
	flags := SlurmV0040GetJobParamsFlagsDETAIL
	params.Flags = &flags

	// Call the generated OpenAPI client to get job details
	resp, err := m.client.apiClient.SlurmV0040GetJobWithResponse(ctx, jobID, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.40")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.40", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.40")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response but got nil")
	}

	// Find the job in the response
	if len(resp.JSON200.Jobs) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "Job not found", fmt.Sprintf("Job ID %s not found", jobID))
	}

	if len(resp.JSON200.Jobs) > 1 {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected multiple jobs returned", fmt.Sprintf("Expected 1 job but got %d for ID %s", len(resp.JSON200.Jobs), jobID))
	}

	// Note: V0040JobInfo does not include step details in v0.0.40 API
	// Steps would need to be retrieved through a dedicated step endpoint if available
	// For now, return empty step list as V0040JobInfo doesn't contain step information
	steps := make([]interfaces.JobStep, 0)

	return &interfaces.JobStepList{
		Steps: steps,
		Total: len(steps),
	}, nil
}

// Watch provides real-time job updates through polling
// Note: v0.0.40 API does not support native streaming/WebSocket job monitoring
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

// GetJobUtilization retrieves minimal resource utilization metrics for a job
// Note: v0.0.40 only supports very basic accounting data
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

	// In v0.0.40, only minimal accounting data is available
	// Most metrics are estimated based on allocated resources
	// TODO: Integrate with basic SLURM accounting when available

	utilization := &interfaces.JobUtilization{
		JobID:     jobID,
		JobName:   job.Name,
		StartTime: job.SubmitTime,
		EndTime:   job.EndTime,

		// CPU Utilization (minimal in v0.0.40 - assumes 70% utilization)
		CPUUtilization: &interfaces.ResourceUtilization{
			Used:       float64(job.CPUs) * 0.70, // Fixed 70% utilization assumption
			Allocated:  float64(job.CPUs),
			Limit:      float64(job.CPUs),
			Percentage: 70.0,
			Metadata: map[string]interface{}{
				"estimation_method": "fixed_percentage",
				"confidence":        "low",
			},
		},

		// Memory Utilization (minimal in v0.0.40 - assumes 60% utilization)
		MemoryUtilization: &interfaces.ResourceUtilization{
			Used:       float64(job.Memory) * 0.60, // Fixed 60% utilization assumption
			Allocated:  float64(job.Memory),
			Limit:      float64(job.Memory),
			Percentage: 60.0,
			Metadata: map[string]interface{}{
				"estimation_method": "fixed_percentage",
				"confidence":        "low",
			},
		},
	}

	// Add metadata about v0.0.40 limitations
	utilization.Metadata = map[string]interface{}{
		"version":       "v0.0.40",
		"source":        "basic_accounting",
		"nodes":         job.Nodes,
		"partition":     job.Partition,
		"state":         job.State,
		"feature_level": "minimal",   // v0.0.40 has minimal features
		"data_quality":  "estimated", // Most data is estimated, not measured
		"limitations": []string{
			"fixed_utilization_percentages",
			"no_actual_measurements",
			"no_gpu_support",
			"no_io_metrics",
			"no_network_metrics",
			"no_energy_metrics",
			"no_performance_counters",
		},
	}

	// All advanced metrics are not supported in v0.0.40
	utilization.GPUUtilization = nil
	utilization.IOUtilization = nil
	utilization.NetworkUtilization = nil
	utilization.EnergyUsage = nil

	return utilization, nil
}

// GetJobEfficiency calculates minimal efficiency metrics for a completed job
// Note: v0.0.40 only provides rough estimates based on assumptions
func (m *JobManagerImpl) GetJobEfficiency(ctx context.Context, jobID string) (*interfaces.ResourceUtilization, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Get job utilization data (not used in v0.0.40 basic calculation)
	_, err := m.GetJobUtilization(ctx, jobID)
	if err != nil {
		return nil, err
	}

	// v0.0.40 uses fixed efficiency estimate based on assumed utilization
	// This is very basic and not based on actual measurements
	efficiency := 65.0 // Fixed 65% efficiency assumption for v0.0.40

	return &interfaces.ResourceUtilization{
		Used:       efficiency,
		Allocated:  100.0,
		Limit:      100.0,
		Percentage: efficiency,
		Metadata: map[string]interface{}{
			"cpu_efficiency":     70.0, // Fixed from utilization assumption
			"memory_efficiency":  60.0, // Fixed from utilization assumption
			"calculation_method": "fixed_estimate_v40",
			"version":            "v0.0.40",
			"confidence":         "very_low",
			"note":               "Efficiency is estimated, not measured in v0.0.40",
			"limitations": []string{
				"no_actual_measurements",
				"fixed_efficiency_value",
				"no_resource_specific_data",
			},
		},
	}, nil
}

// GetJobPerformance retrieves minimal performance metrics for a job
// Note: v0.0.40 provides only the most basic information
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

	// Build minimal performance report (v0.0.40 version with very limited features)
	performance := &interfaces.JobPerformance{
		JobID:     uint32(jobIDInt),
		JobName:   job.Name,
		StartTime: job.SubmitTime,
		EndTime:   job.EndTime,
		Status:    job.State,
		ExitCode:  job.ExitCode,

		ResourceUtilization: efficiency,
		JobUtilization:      utilization,

		// No advanced features available in v0.0.40
		StepMetrics:       nil,
		PerformanceTrends: nil,
		Bottlenecks:       nil, // No bottleneck detection in v0.0.40

		// Only basic recommendation in v0.0.40
		Recommendations: []interfaces.OptimizationRecommendation{
			{
				Type:                "system",
				Priority:            "medium",
				Title:               "Upgrade for better analytics",
				Description:         "API v0.0.40 provides only minimal analytics. Upgrade to SLURM API v0.0.41+ for actual resource measurements and v0.0.42+ for comprehensive analytics.",
				ExpectedImprovement: 0.0,
				ConfigChanges: map[string]string{
					"current_api_version": "v0.0.40",
					"minimum_recommended": "v0.0.41",
					"optimal_version":     "v0.0.42_or_higher",
				},
			},
		},
	}

	return performance, nil
}

// GetJobLiveMetrics retrieves real-time performance metrics for a running job
// Note: v0.0.40 doesn't support real-time monitoring, returns error
func (m *JobManagerImpl) GetJobLiveMetrics(ctx context.Context, jobID string) (*interfaces.JobLiveMetrics, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// v0.0.40 doesn't support live metrics
	return nil, errors.NewNotImplementedError(
		"GetJobLiveMetrics",
		"Real-time job monitoring is not supported in API v0.0.40. Please upgrade to v0.0.41 or higher.",
	)
}

// WatchJobMetrics provides streaming performance updates for a running job
// Note: v0.0.40 has minimal support - only job state changes, no performance metrics
func (m *JobManagerImpl) WatchJobMetrics(ctx context.Context, jobID string, opts *interfaces.WatchMetricsOptions) (<-chan interfaces.JobMetricsEvent, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Default options if not provided - v0.0.40 severe limitations
	if opts == nil {
		opts = &interfaces.WatchMetricsOptions{
			UpdateInterval:     30 * time.Second, // Very slow polling for v0.0.40
			IncludeCPU:         false,            // Not supported
			IncludeMemory:      false,            // Not supported
			IncludeGPU:         false,            // Not supported
			IncludeNetwork:     false,            // Not supported
			IncludeIO:          false,            // Not supported
			IncludeEnergy:      false,            // Not supported
			IncludeNodeMetrics: false,            // Not supported
			StopOnCompletion:   true,
		}
	}

	// Enforce minimum update interval for v0.0.40
	if opts.UpdateInterval < 30*time.Second {
		opts.UpdateInterval = 30 * time.Second
	}

	// Create event channel
	eventChan := make(chan interfaces.JobMetricsEvent, 2)

	// Start monitoring goroutine - only tracks state changes
	go func() {
		defer close(eventChan)

		// Send initial warning about limitations
		eventChan <- interfaces.JobMetricsEvent{
			Type:      "update",
			JobID:     jobID,
			Timestamp: time.Now(),
			Metrics: &interfaces.JobLiveMetrics{
				JobID: jobID,
				State: "UNKNOWN",
				Metadata: map[string]interface{}{
					"warning": "v0.0.40 only supports job state monitoring, no performance metrics available",
					"version": "v0.0.40",
				},
			},
		}

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

		// Minimal monitoring loop for v0.0.40
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
					if opts.StopOnCompletion && isJobCompleteV40(job.State) {
						eventChan <- interfaces.JobMetricsEvent{
							Type:      "complete",
							JobID:     jobID,
							Timestamp: time.Now(),
						}
						return
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

// Helper function to check if job is complete (v0.0.40)
func isJobCompleteV40(state string) bool {
	completedStates := []string{
		"COMPLETED", "FAILED", "CANCELLED", "TIMEOUT", "NODE_FAIL",
	}
	for _, s := range completedStates {
		if state == s {
			return true
		}
	}
	return false
}

// GetJobResourceTrends retrieves performance trends over specified time windows
// Note: v0.0.40 doesn't support trend analysis, returns minimal data
func (m *JobManagerImpl) GetJobResourceTrends(ctx context.Context, jobID string, opts *interfaces.ResourceTrendsOptions) (*interfaces.JobResourceTrends, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// v0.0.40 doesn't support resource trends
	return nil, errors.NewNotImplementedError(
		"GetJobResourceTrends",
		"Resource trend analysis is not supported in API v0.0.40. Please upgrade to v0.0.41 or higher.",
	)
}

// GetJobStepDetails retrieves minimal job step information (v0.0.40 - very limited features)
func (m *JobManagerImpl) GetJobStepDetails(ctx context.Context, jobID string, stepID string) (*interfaces.JobStepDetails, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// First get the job details to validate job exists
	job, err := m.Get(ctx, jobID)
	if err != nil {
		return nil, err
	}

	// Parse step ID (not used in v0.0.40)
	_, err = strconv.Atoi(stepID)
	if err != nil {
		return nil, errors.NewClientError(errors.ErrorCodeInvalidRequest, "Invalid step ID format", err.Error())
	}

	// v0.0.40 has very minimal step tracking
	stepDetails := &interfaces.JobStepDetails{
		StepID:    stepID,
		StepName:  fmt.Sprintf("step_%s", stepID),
		JobID:     jobID,
		JobName:   job.Name,
		State:     deriveMinimalStepState(job.State),
		StartTime: job.StartTime,
		EndTime:   job.EndTime,
		Duration:  calculateMinimalStepDuration(job.StartTime, job.EndTime),
		ExitCode:  job.ExitCode, // Simple inheritance

		// Minimal resource allocation for v0.0.40
		CPUAllocation:    job.CPUs,          // Use all job CPUs for simplicity
		MemoryAllocation: int64(job.Memory), // Use all job memory
		NodeList:         job.Nodes,
		TaskCount:        job.CPUs, // Simple 1:1 mapping

		// Minimal command info
		Command:     deriveMinimalStepCommand(job.Command),
		CommandLine: job.Command,
		WorkingDir:  job.WorkingDir,
		Environment: job.Environment,

		// Minimal performance metrics (v0.0.40)
		CPUTime:    time.Hour,        // Fixed 1 hour
		UserTime:   time.Hour,        // Fixed 1 hour
		SystemTime: time.Minute * 10, // Fixed 10 minutes

		// Minimal resource usage
		MaxRSS:     int64(job.Memory / 10), // Very conservative
		MaxVMSize:  int64(job.Memory / 5),  // Very conservative
		AverageRSS: int64(job.Memory / 20), // Very conservative

		// Minimal I/O statistics
		TotalReadBytes:  calculateMinimalStepIOBytes(job.CPUs, "read"),
		TotalWriteBytes: calculateMinimalStepIOBytes(job.CPUs, "write"),
		ReadOperations:  calculateMinimalStepIOOps(job.CPUs, "read"),
		WriteOperations: calculateMinimalStepIOOps(job.CPUs, "write"),

		// No network or energy statistics in v0.0.40
		NetworkBytesReceived: 0,
		NetworkBytesSent:     0,
		EnergyConsumed:       0,
		AveragePowerDraw:     0,

		// Minimal task-level information
		Tasks: generateMinimalStepTasks(job),

		// Minimal step-specific metadata
		StepType:        "primary", // Fixed type
		Priority:        job.Priority,
		AccountingGroup: "default",
		QOSLevel:        "normal",
	}

	// Add metadata (v0.0.40 specific)
	stepDetails.Metadata = map[string]interface{}{
		"version":          "v0.0.40",
		"data_source":      "simulated",
		"job_partition":    job.Partition,
		"job_submit_time":  job.SubmitTime,
		"minimal_tracking": true, // v0.0.40 feature level
		"very_limited":     true, // v0.0.40 limitation
		"fixed_values":     true, // Most values are fixed in v0.0.40
	}

	return stepDetails, nil
}

// GetJobStepUtilization retrieves minimal resource utilization metrics (v0.0.40 - very limited)
func (m *JobManagerImpl) GetJobStepUtilization(ctx context.Context, jobID string, stepID string) (*interfaces.JobStepUtilization, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Get step details first
	stepDetails, err := m.GetJobStepDetails(ctx, jobID, stepID)
	if err != nil {
		return nil, err
	}

	// Get job details for context
	job, err := m.Get(ctx, jobID)
	if err != nil {
		return nil, err
	}

	// Create minimal step utilization metrics (v0.0.40)
	stepUtilization := &interfaces.JobStepUtilization{
		StepID:   stepID,
		StepName: stepDetails.StepName,
		JobID:    jobID,
		JobName:  job.Name,

		// Time information
		StartTime: stepDetails.StartTime,
		EndTime:   stepDetails.EndTime,
		Duration:  stepDetails.Duration,

		// Minimal CPU utilization metrics
		CPUUtilization: &interfaces.ResourceUtilization{
			Used:       float64(stepDetails.CPUAllocation) * 0.5, // Fixed 50% utilization
			Allocated:  float64(stepDetails.CPUAllocation),
			Limit:      float64(stepDetails.CPUAllocation),
			Percentage: 50.0, // Fixed percentage for v0.0.40
			Metadata: map[string]interface{}{
				"minimal_tracking": true, // v0.0.40 limitation
				"fixed_value":      true, // v0.0.40 limitation
			},
		},

		// Minimal memory utilization metrics
		MemoryUtilization: &interfaces.ResourceUtilization{
			Used:       float64(stepDetails.AverageRSS),
			Allocated:  float64(stepDetails.MemoryAllocation),
			Limit:      float64(stepDetails.MemoryAllocation),
			Percentage: 30.0, // Fixed percentage for v0.0.40
			Metadata: map[string]interface{}{
				"minimal_tracking": true, // v0.0.40 limitation
				"fixed_value":      true, // v0.0.40 limitation
			},
		},

		// Minimal I/O utilization
		IOUtilization: &interfaces.IOUtilization{
			ReadBandwidth: &interfaces.ResourceUtilization{
				Used:       calculateMinimalIOBandwidth(stepDetails.TotalReadBytes, stepDetails.Duration),
				Allocated:  100 * 1024 * 1024, // 100 MB/s limit (very low)
				Limit:      100 * 1024 * 1024,
				Percentage: 10.0, // Fixed very low percentage
			},
			WriteBandwidth: &interfaces.ResourceUtilization{
				Used:       calculateMinimalIOBandwidth(stepDetails.TotalWriteBytes, stepDetails.Duration),
				Allocated:  100 * 1024 * 1024, // 100 MB/s limit
				Limit:      100 * 1024 * 1024,
				Percentage: 8.0, // Fixed very low percentage
			},
			TotalBytesRead:    stepDetails.TotalReadBytes,
			TotalBytesWritten: stepDetails.TotalWriteBytes,
		},

		// No network utilization in v0.0.40
		NetworkUtilization: &interfaces.NetworkUtilization{
			TotalBandwidth: &interfaces.ResourceUtilization{
				Used:       0,
				Allocated:  0,
				Limit:      0,
				Percentage: 0,
			},
			PacketsReceived: 0,
			PacketsSent:     0,
			PacketsDropped:  0,
			Errors:          0,
			Interfaces:      make(map[string]interfaces.NetworkInterfaceStats),
		},

		// No energy utilization in v0.0.40
		EnergyUtilization: &interfaces.ResourceUtilization{
			Used:       0,
			Allocated:  0,
			Limit:      0,
			Percentage: 0,
			Metadata: map[string]interface{}{
				"not_supported": true, // v0.0.40 limitation
			},
		},

		// Minimal task-level utilization
		TaskUtilizations: generateMinimalTaskUtilizations(stepDetails),

		// Minimal performance metrics
		PerformanceMetrics: &interfaces.StepPerformanceMetrics{
			CPUEfficiency:     50.0, // Fixed value for v0.0.40
			MemoryEfficiency:  30.0, // Fixed value for v0.0.40
			IOEfficiency:      25.0, // Fixed value for v0.0.40
			OverallEfficiency: 35.0, // Fixed value for v0.0.40

			// Minimal bottleneck analysis
			PrimaryBottleneck:  "cpu", // Fixed for v0.0.40
			BottleneckSeverity: "low",
			ResourceBalance:    "unbalanced",

			// Fixed minimal performance indicators
			ThroughputMBPS:   50.0, // Fixed value
			LatencyMS:        20.0, // Fixed value
			ScalabilityScore: 60.0, // Fixed value
		},
	}

	// Add metadata (v0.0.40 specific)
	stepUtilization.Metadata = map[string]interface{}{
		"version":               "v0.0.40",
		"data_source":           "simulated",
		"task_count":            stepDetails.TaskCount,
		"node_count":            len(stepDetails.NodeList),
		"minimal_features":      true, // v0.0.40 feature level
		"very_limited_accuracy": true, // v0.0.40 limitation
		"all_fixed_metrics":     true, // All metrics are fixed values in v0.0.40
		"upgrade_recommended":   true, // Recommendation to upgrade
	}

	return stepUtilization, nil
}

// Helper functions for v0.0.40 minimal step calculations

func deriveMinimalStepState(jobState string) string {
	// Very simple step state derivation for v0.0.40
	return jobState // Direct inheritance
}

func calculateMinimalStepDuration(startTime, endTime *time.Time) time.Duration {
	if startTime == nil {
		return time.Hour // Fixed 1 hour default
	}
	if endTime == nil {
		return time.Since(*startTime)
	}
	return endTime.Sub(*startTime)
}

func deriveMinimalStepCommand(jobCommand string) string {
	if jobCommand == "" {
		return "srun /bin/bash" // Basic command
	}
	return jobCommand // Direct inheritance
}

func calculateMinimalStepIOBytes(cpus int, ioType string) int64 {
	base := int64(cpus) * 256 * 1024 * 1024 // 256MB per CPU base (very low)
	if ioType == "write" {
		base = base / 3 // Write is third of read
	}
	return base
}

func calculateMinimalStepIOOps(cpus int, ioType string) int64 {
	base := int64(cpus) * 2000 // 2K ops per CPU base (very low)
	if ioType == "write" {
		base = base / 3
	}
	return base
}

func generateMinimalStepTasks(job *interfaces.Job) []interfaces.StepTaskInfo {
	taskCount := job.CPUs // Simple 1:1 mapping
	tasks := make([]interfaces.StepTaskInfo, taskCount)

	for i := 0; i < taskCount; i++ {
		// Minimal task distribution in v0.0.40
		nodeIndex := i % len(job.Nodes)
		nodeName := job.Nodes[nodeIndex]

		tasks[i] = interfaces.StepTaskInfo{
			TaskID:    i,
			NodeName:  nodeName,
			LocalID:   i, // Simple ID
			State:     job.State,
			ExitCode:  job.ExitCode,
			CPUTime:   time.Minute * 30,              // Fixed 30 minutes
			MaxRSS:    int64(job.Memory / taskCount), // Simple distribution
			StartTime: job.StartTime,
			EndTime:   job.EndTime,
		}
	}

	return tasks
}

func generateMinimalTaskUtilizations(stepDetails *interfaces.JobStepDetails) []interfaces.TaskUtilization {
	tasks := make([]interfaces.TaskUtilization, len(stepDetails.Tasks))

	for i, task := range stepDetails.Tasks {
		// Fixed minimal utilization values for v0.0.40
		tasks[i] = interfaces.TaskUtilization{
			TaskID:            task.TaskID,
			NodeName:          task.NodeName,
			CPUUtilization:    50.0, // Fixed CPU utilization
			MemoryUtilization: 30.0, // Fixed memory utilization
			State:             task.State,
			ExitCode:          task.ExitCode,
		}
	}

	return tasks
}

func calculateMinimalIOBandwidth(totalBytes int64, duration time.Duration) float64 {
	if duration == 0 {
		return 10.0 // Fixed 10 MB/s default
	}
	return float64(totalBytes) / duration.Seconds()
}

// ListJobStepsWithMetrics retrieves all job steps with minimal metrics for v0.0.40
func (m *JobManagerImpl) ListJobStepsWithMetrics(ctx context.Context, jobID string, opts *interfaces.ListJobStepsOptions) (*interfaces.JobStepMetricsList, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// First get the job details to validate the job exists
	job, err := m.Get(ctx, jobID)
	if err != nil {
		return nil, errors.WrapError(err)
	}

	// Get job steps using the existing Steps method
	stepList, err := m.Steps(ctx, jobID)
	if err != nil {
		return nil, errors.WrapError(err)
	}

	// Process steps with minimal metrics for v0.0.40
	filteredSteps := []*interfaces.JobStepWithMetrics{}
	
	for _, step := range stepList.Steps {
		// Very basic processing for v0.0.40 - no filtering to keep it simple
		stepDetails, err := m.GetJobStepDetails(ctx, jobID, step.ID)
		if err != nil {
			continue // Skip steps with errors
		}

		stepUtilization, err := m.GetJobStepUtilization(ctx, jobID, step.ID)
		if err != nil {
			continue // Skip steps with errors
		}

		// Create step with minimal metrics
		stepWithMetrics := &interfaces.JobStepWithMetrics{
			JobStepDetails:     stepDetails,
			JobStepUtilization: stepUtilization,
		}

		filteredSteps = append(filteredSteps, stepWithMetrics)
	}

	// No advanced options for v0.0.40 - just basic pagination
	if opts != nil && opts.Limit > 0 && len(filteredSteps) > opts.Limit {
		filteredSteps = filteredSteps[:opts.Limit]
	}

	// Generate very basic summary
	summary := generateVeryBasicJobStepsSummary(filteredSteps, convertToJobStepPointers(stepList.Steps))

	result := &interfaces.JobStepMetricsList{
		JobID:         jobID,
		JobName:       job.Name,
		Steps:         filteredSteps,
		Summary:       summary,
		TotalSteps:    len(stepList.Steps),
		FilteredSteps: len(filteredSteps),
		Metadata: map[string]interface{}{
			"api_version":    "v0.0.40",
			"generated_at":   time.Now(),
			"job_state":      job.State,
			"analysis_level": "minimal",
			"note":           "Limited metrics available in v0.0.40",
		},
	}

	return result, nil
}

// Helper function to generate very basic summary for v0.0.40
func generateVeryBasicJobStepsSummary(filteredSteps []*interfaces.JobStepWithMetrics, allSteps []*interfaces.JobStep) *interfaces.JobStepsSummary {
	summary := &interfaces.JobStepsSummary{
		TotalSteps: len(allSteps),
	}

	if len(filteredSteps) == 0 {
		return summary
	}

	// Very basic aggregation for v0.0.40
	completedSteps := 0
	for _, step := range filteredSteps {
		if step.State == "COMPLETED" {
			completedSteps++
		}
	}

	summary.CompletedSteps = completedSteps

	// Conservative fixed estimates for v0.0.40
	summary.AverageCPUEfficiency = 50.0
	summary.AverageMemoryEfficiency = 45.0
	summary.AverageIOEfficiency = 40.0
	summary.AverageOverallEfficiency = 45.0
	summary.OptimizationPotential = 55.0

	return summary
}

// Helper function to convert []JobStep to []*JobStep
func convertToJobStepPointers(steps []interfaces.JobStep) []*interfaces.JobStep {
	result := make([]*interfaces.JobStep, len(steps))
	for i := range steps {
		result[i] = &steps[i]
	}
	return result
}
