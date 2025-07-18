package v0_0_43

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
	params := &SlurmV0043GetJobsParams{}

	// Set flags to get detailed job information
	flags := SlurmV0043GetJobsParamsFlagsDETAIL
	params.Flags = &flags

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0043GetJobsWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
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

// convertAPIJobToInterface converts a V0043JobInfo to interfaces.Job
func convertAPIJobToInterface(apiJob V0043JobInfo) (*interfaces.Job, error) {
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
		job.State = string((*apiJob.JobState)[0])
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
	if apiJob.JobResources != nil && apiJob.JobResources.Nodes != nil && apiJob.JobResources.Nodes.List != nil {
		// Parse node list string into slice
		nodeListStr := *apiJob.JobResources.Nodes.List
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
	params := &SlurmV0043GetJobParams{}

	// Set flags to get detailed job information
	flags := SlurmV0043GetJobParamsFlagsDETAIL
	params.Flags = &flags

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0043GetJobWithResponse(ctx, jobID, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
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
	requestBody := SlurmV0043PostJobSubmitJSONRequestBody{
		Job: jobDesc,
	}

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0043PostJobSubmitWithResponse(ctx, requestBody)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
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

// convertJobSubmissionToAPI converts interfaces.JobSubmission to V0043JobDescMsg
func convertJobSubmissionToAPI(job *interfaces.JobSubmission) (*V0043JobDescMsg, error) {
	if job == nil {
		return nil, errors.NewClientError(errors.ErrorCodeInvalidRequest, "Job submission cannot be nil")
	}

	jobDesc := &V0043JobDescMsg{}

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
		jobDesc.MemoryPerNode = &V0043Uint64NoValStruct{
			Number: &memoryMB,
			Set:    &set,
		}
	}

	if job.TimeLimit > 0 {
		timeLimit := int32(job.TimeLimit)
		set := true
		jobDesc.TimeLimit = &V0043Uint32NoValStruct{
			Number: &timeLimit,
			Set:    &set,
		}
	}

	if job.Nodes > 0 {
		nodes := int32(job.Nodes)
		jobDesc.MinimumNodes = &nodes
	}

	if job.Priority > 0 {
		priority := int32(job.Priority)
		set := true
		jobDesc.Priority = &V0043Uint32NoValStruct{
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
	params := &SlurmV0043DeleteJobParams{}

	// Send SIGTERM signal by default (can be made configurable later)
	signal := "SIGTERM"
	params.Signal = &signal

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0043DeleteJobWithResponse(ctx, jobID, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
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
	resp, err := m.client.apiClient.SlurmV0043PostJobWithResponse(ctx, jobID, *jobDesc)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
		return httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil {
		return errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response but got nil")
	}

	return nil
}

// convertJobUpdateToAPI converts interfaces.JobUpdate to V0043JobDescMsg
func convertJobUpdateToAPI(update *interfaces.JobUpdate) (*V0043JobDescMsg, error) {
	if update == nil {
		return nil, errors.NewClientError(errors.ErrorCodeInvalidRequest, "Job update cannot be nil")
	}

	jobDesc := &V0043JobDescMsg{}

	// Only include fields that are actually being updated (non-nil values)
	if update.Priority != nil {
		priority := int32(*update.Priority)
		set := true
		jobDesc.Priority = &V0043Uint32NoValStruct{
			Number: &priority,
			Set:    &set,
		}
	}

	if update.TimeLimit != nil {
		timeLimit := int32(*update.TimeLimit)
		set := true
		jobDesc.TimeLimit = &V0043Uint32NoValStruct{
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
	params := &SlurmV0043GetJobParams{}

	// Set flags to get detailed job information including steps
	flags := SlurmV0043GetJobParamsFlagsDETAIL
	params.Flags = &flags

	// Call the generated OpenAPI client to get job details
	resp, err := m.client.apiClient.SlurmV0043GetJobWithResponse(ctx, jobID, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
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

	// Note: V0043JobInfo does not include step details in v0.0.43 API
	// Steps would need to be retrieved through a dedicated step endpoint if available
	// For now, return empty step list as V0043JobInfo doesn't contain step information
	steps := make([]interfaces.JobStep, 0)

	return &interfaces.JobStepList{
		Steps: steps,
		Total: len(steps),
	}, nil
}

// Watch provides real-time job updates through polling
// Note: v0.0.43 API does not support native streaming/WebSocket job monitoring
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

// GetJobUtilization retrieves comprehensive resource utilization metrics for a job
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

	// In v0.0.43, full utilization metrics are available through the sacct-like endpoints
	// For now, we'll simulate utilization data based on job info
	// TODO: Integrate with real SLURM accounting/statistics endpoints when available

	utilization := &interfaces.JobUtilization{
		JobID:   jobID,
		JobName: job.Name,
		StartTime: job.SubmitTime,
		EndTime: job.EndTime,
		
		// CPU Utilization
		CPUUtilization: &interfaces.ResourceUtilization{
			Used:      float64(job.CPUs) * 0.85, // Simulated 85% utilization
			Allocated: float64(job.CPUs),
			Limit:     float64(job.CPUs),
			Percentage: 85.0,
		},
		
		// Memory Utilization
		MemoryUtilization: &interfaces.ResourceUtilization{
			Used:      float64(job.Memory) * 0.72, // Simulated 72% utilization
			Allocated: float64(job.Memory),
			Limit:     float64(job.Memory),
			Percentage: 72.0,
		},
	}

	// Add metadata
	utilization.Metadata = map[string]interface{}{
		"version": "v0.0.43",
		"source": "simulated", // TODO: Change to "accounting" when real data available
		"nodes": job.Nodes,
		"partition": job.Partition,
		"state": job.State,
	}

	// GPU utilization (if applicable)
	// In a real implementation, this would come from GPU monitoring data
	if gpuCount, ok := job.Metadata["gpu_count"].(int); ok && gpuCount > 0 {
		utilization.GPUUtilization = &interfaces.GPUUtilization{
			TotalGPUs: gpuCount,
			GPUs: make(map[string]interfaces.GPUDeviceUtilization),
			AverageUtilization: &interfaces.ResourceUtilization{
				Used:      float64(gpuCount) * 0.90, // Simulated 90% GPU utilization
				Allocated: float64(gpuCount),
				Limit:     float64(gpuCount),
				Percentage: 90.0,
			},
		}
		
		// Add per-GPU metrics
		for i := 0; i < gpuCount; i++ {
			gpuID := fmt.Sprintf("gpu%d", i)
			utilization.GPUUtilization.GPUs[gpuID] = interfaces.GPUDeviceUtilization{
				DeviceID:   gpuID,
				DeviceUUID: fmt.Sprintf("GPU-%d-UUID", i),
				Utilization: 85.0 + float64(i*5), // Varying utilization
				MemoryUsed:  16 * 1024 * 1024 * 1024, // 16GB
				MemoryTotal: 24 * 1024 * 1024 * 1024, // 24GB
				Temperature: 65.0 + float64(i*2),
				PowerDraw:   250.0 + float64(i*10),
				PowerLimit:  300.0,
			}
		}
	}

	// I/O utilization (simulated)
	utilization.IOUtilization = &interfaces.IOUtilization{
		ReadBandwidth: &interfaces.ResourceUtilization{
			Used:      100 * 1024 * 1024, // 100 MB/s
			Allocated: 500 * 1024 * 1024, // 500 MB/s limit
			Limit:     500 * 1024 * 1024,
			Percentage: 20.0,
		},
		WriteBandwidth: &interfaces.ResourceUtilization{
			Used:      50 * 1024 * 1024, // 50 MB/s
			Allocated: 500 * 1024 * 1024, // 500 MB/s limit
			Limit:     500 * 1024 * 1024,
			Percentage: 10.0,
		},
		TotalBytesRead:    10 * 1024 * 1024 * 1024, // 10 GB
		TotalBytesWritten: 5 * 1024 * 1024 * 1024,  // 5 GB
	}

	// Network utilization (simulated)
	utilization.NetworkUtilization = &interfaces.NetworkUtilization{
		TotalBandwidth: &interfaces.ResourceUtilization{
			Used:      1 * 1024 * 1024 * 1024, // 1 Gbps
			Allocated: 10 * 1024 * 1024 * 1024, // 10 Gbps limit
			Limit:     10 * 1024 * 1024 * 1024,
			Percentage: 10.0,
		},
		PacketsReceived: 1000000,
		PacketsSent:     900000,
		PacketsDropped:  100,
		Errors:          5,
		Interfaces:      make(map[string]interfaces.NetworkInterfaceStats),
	}

	// Energy usage (simulated)
	if job.EndTime != nil {
		duration := job.EndTime.Sub(job.StartTime.Unix() > 0 ? *job.StartTime : job.SubmitTime).Hours()
		avgPower := 300.0 // 300W average
		utilization.EnergyUsage = &interfaces.EnergyUsage{
			TotalEnergyJoules: avgPower * duration * 3600, // Convert to joules
			AveragePowerWatts: avgPower,
			PeakPowerWatts:    450.0,
			MinPowerWatts:     200.0,
			CPUEnergyJoules:   200.0 * duration * 3600,
			GPUEnergyJoules:   100.0 * duration * 3600,
			CarbonFootprint:   avgPower * duration * 0.0004, // Approximate carbon factor
		}
	}

	return utilization, nil
}

// GetJobEfficiency calculates efficiency metrics for a completed job
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

	// Calculate overall efficiency based on resource utilization
	// This is a weighted average of different resource efficiencies
	cpuWeight := 0.4
	memWeight := 0.3
	gpuWeight := 0.2
	ioWeight := 0.1

	totalEfficiency := 0.0
	totalWeight := 0.0

	// CPU efficiency
	if utilization.CPUUtilization != nil {
		totalEfficiency += utilization.CPUUtilization.Percentage * cpuWeight
		totalWeight += cpuWeight
	}

	// Memory efficiency
	if utilization.MemoryUtilization != nil {
		totalEfficiency += utilization.MemoryUtilization.Percentage * memWeight
		totalWeight += memWeight
	}

	// GPU efficiency (if applicable)
	if utilization.GPUUtilization != nil && utilization.GPUUtilization.AverageUtilization != nil {
		totalEfficiency += utilization.GPUUtilization.AverageUtilization.Percentage * gpuWeight
		totalWeight += gpuWeight
	} else {
		// If no GPU, redistribute weight to CPU and memory
		totalWeight += gpuWeight * 0.6 // Add 60% to CPU
		totalWeight += gpuWeight * 0.4 // Add 40% to memory
	}

	// I/O efficiency
	if utilization.IOUtilization != nil && utilization.IOUtilization.ReadBandwidth != nil {
		ioEfficiency := (utilization.IOUtilization.ReadBandwidth.Percentage + 
		                utilization.IOUtilization.WriteBandwidth.Percentage) / 2
		totalEfficiency += ioEfficiency * ioWeight
		totalWeight += ioWeight
	}

	// Calculate final efficiency percentage
	efficiency := totalEfficiency / totalWeight

	return &interfaces.ResourceUtilization{
		Used:       efficiency,
		Allocated:  100.0,
		Limit:      100.0,
		Percentage: efficiency,
		Metadata: map[string]interface{}{
			"cpu_efficiency":    utilization.CPUUtilization.Percentage,
			"memory_efficiency": utilization.MemoryUtilization.Percentage,
			"calculation_method": "weighted_average",
			"weights": map[string]float64{
				"cpu":    cpuWeight,
				"memory": memWeight,
				"gpu":    gpuWeight,
				"io":     ioWeight,
			},
		},
	}, nil
}

// GetJobPerformance retrieves detailed performance metrics for a job
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

	// Build performance report
	performance := &interfaces.JobPerformance{
		JobID:     uint32(jobIDInt),
		JobName:   job.Name,
		StartTime: job.SubmitTime,
		EndTime:   job.EndTime,
		Status:    job.State,
		ExitCode:  job.ExitCode,
		
		ResourceUtilization: efficiency,
		JobUtilization:     utilization,
		
		// Step metrics would come from job step data in real implementation
		StepMetrics: []interfaces.JobStepPerformance{},
		
		// Performance trends (simulated for now)
		PerformanceTrends: generatePerformanceTrends(job),
		
		// Bottleneck analysis
		Bottlenecks: analyzeBottlenecks(utilization),
		
		// Optimization recommendations
		Recommendations: generateRecommendations(utilization, efficiency),
	}

	return performance, nil
}

// Helper function to generate performance trends (simulated)
func generatePerformanceTrends(job *interfaces.Job) *interfaces.PerformanceTrends {
	if job.StartTime == nil {
		return nil
	}

	// Generate hourly time points
	startTime := *job.StartTime
	endTime := time.Now()
	if job.EndTime != nil {
		endTime = *job.EndTime
	}

	duration := endTime.Sub(startTime)
	points := int(duration.Hours())
	if points < 1 {
		points = 1
	}
	if points > 24 {
		points = 24 // Limit to 24 data points
	}

	trends := &interfaces.PerformanceTrends{
		TimePoints:    make([]time.Time, points),
		CPUTrends:     make([]float64, points),
		MemoryTrends:  make([]float64, points),
		IOTrends:      make([]float64, points),
		NetworkTrends: make([]float64, points),
	}

	// Generate simulated trend data
	for i := 0; i < points; i++ {
		trends.TimePoints[i] = startTime.Add(time.Duration(i) * time.Hour)
		
		// Simulate varying utilization over time
		trends.CPUTrends[i] = 70.0 + float64(i%10)*2
		trends.MemoryTrends[i] = 60.0 + float64(i%8)*3
		trends.IOTrends[i] = 20.0 + float64(i%5)*4
		trends.NetworkTrends[i] = 10.0 + float64(i%3)*5
	}

	return trends
}

// Helper function to analyze bottlenecks
func analyzeBottlenecks(utilization *interfaces.JobUtilization) []interfaces.PerformanceBottleneck {
	bottlenecks := []interfaces.PerformanceBottleneck{}

	// Check CPU bottleneck
	if utilization.CPUUtilization != nil && utilization.CPUUtilization.Percentage > 95 {
		bottlenecks = append(bottlenecks, interfaces.PerformanceBottleneck{
			Type:         "cpu",
			Severity:     "high",
			Description:  "CPU utilization is at or near maximum capacity",
			Impact:       15.0, // 15% performance impact
			TimeDetected: time.Now(),
			Duration:     time.Hour, // Simulated duration
		})
	}

	// Check memory bottleneck
	if utilization.MemoryUtilization != nil && utilization.MemoryUtilization.Percentage > 90 {
		bottlenecks = append(bottlenecks, interfaces.PerformanceBottleneck{
			Type:         "memory",
			Severity:     "medium",
			Description:  "Memory utilization is high, may cause swapping",
			Impact:       10.0,
			TimeDetected: time.Now(),
			Duration:     30 * time.Minute,
		})
	}

	// Check I/O bottleneck
	if utilization.IOUtilization != nil {
		readUtil := utilization.IOUtilization.ReadBandwidth.Percentage
		writeUtil := utilization.IOUtilization.WriteBandwidth.Percentage
		if readUtil > 80 || writeUtil > 80 {
			bottlenecks = append(bottlenecks, interfaces.PerformanceBottleneck{
				Type:         "io",
				Severity:     "medium",
				Description:  "I/O bandwidth utilization is high",
				Impact:       8.0,
				TimeDetected: time.Now(),
				Duration:     20 * time.Minute,
			})
		}
	}

	return bottlenecks
}

// Helper function to generate optimization recommendations
func generateRecommendations(utilization *interfaces.JobUtilization, efficiency *interfaces.ResourceUtilization) []interfaces.OptimizationRecommendation {
	recommendations := []interfaces.OptimizationRecommendation{}

	// CPU recommendations
	if utilization.CPUUtilization != nil {
		if utilization.CPUUtilization.Percentage < 50 {
			recommendations = append(recommendations, interfaces.OptimizationRecommendation{
				Type:                "resource_adjustment",
				Priority:            "high",
				Title:               "Reduce CPU allocation",
				Description:         "CPU utilization is low. Consider reducing CPU allocation to improve resource efficiency.",
				ExpectedImprovement: 20.0,
				ResourceChanges: map[string]interface{}{
					"cpu_reduction": "50%",
				},
			})
		} else if utilization.CPUUtilization.Percentage > 95 {
			recommendations = append(recommendations, interfaces.OptimizationRecommendation{
				Type:                "resource_adjustment",
				Priority:            "high",
				Title:               "Increase CPU allocation",
				Description:         "CPU utilization is at maximum. Consider increasing CPU allocation for better performance.",
				ExpectedImprovement: 25.0,
				ResourceChanges: map[string]interface{}{
					"cpu_increase": "50%",
				},
			})
		}
	}

	// Memory recommendations
	if utilization.MemoryUtilization != nil && utilization.MemoryUtilization.Percentage < 40 {
		recommendations = append(recommendations, interfaces.OptimizationRecommendation{
			Type:                "resource_adjustment",
			Priority:            "medium",
			Title:               "Optimize memory allocation",
			Description:         "Memory utilization is low. Consider reducing memory allocation.",
			ExpectedImprovement: 15.0,
			ResourceChanges: map[string]interface{}{
				"memory_reduction": "40%",
			},
		})
	}

	// Overall efficiency recommendation
	if efficiency.Percentage < 70 {
		recommendations = append(recommendations, interfaces.OptimizationRecommendation{
			Type:        "workflow",
			Priority:    "high",
			Title:       "Improve overall job efficiency",
			Description: "Overall job efficiency is below 70%. Review job configuration and resource allocation.",
			ExpectedImprovement: 30.0,
			ConfigChanges: map[string]string{
				"review_parallelization": "true",
				"optimize_data_locality": "true",
			},
		})
	}

	return recommendations
}

// GetJobLiveMetrics retrieves real-time performance metrics for a running job
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

	// Create live metrics response
	// In v0.0.43, we would integrate with real-time monitoring endpoints
	// For now, simulate real-time data based on job info
	liveMetrics := &interfaces.JobLiveMetrics{
		JobID:          jobID,
		JobName:        job.Name,
		State:          job.State,
		RunningTime:    runningTime,
		CollectionTime: time.Now(),
		
		// Simulate CPU usage with variations
		CPUUsage: &interfaces.LiveResourceMetric{
			Current:            float64(job.CPUs) * (0.75 + float64(time.Now().Unix()%20)/100),
			Average1Min:        float64(job.CPUs) * 0.78,
			Average5Min:        float64(job.CPUs) * 0.76,
			Peak:               float64(job.CPUs) * 0.95,
			Allocated:          float64(job.CPUs),
			UtilizationPercent: 78.0,
			Trend:              determineTrend(75.0, 78.0),
			Unit:               "cores",
		},
		
		// Simulate memory usage
		MemoryUsage: &interfaces.LiveResourceMetric{
			Current:            float64(job.Memory) * (0.65 + float64(time.Now().Unix()%15)/100),
			Average1Min:        float64(job.Memory) * 0.68,
			Average5Min:        float64(job.Memory) * 0.66,
			Peak:               float64(job.Memory) * 0.85,
			Allocated:          float64(job.Memory),
			UtilizationPercent: 68.0,
			Trend:              determineTrend(65.0, 68.0),
			Unit:               "bytes",
		},
		
		// Process information (simulated)
		ProcessCount: 8 + int(time.Now().Unix()%4),
		ThreadCount:  32 + int(time.Now().Unix()%8),
		
		// Initialize maps
		NodeMetrics: make(map[string]*interfaces.NodeLiveMetrics),
		Alerts:      []interfaces.PerformanceAlert{},
	}

	// Add GPU metrics if available
	if gpuCount, ok := job.Metadata["gpu_count"].(int); ok && gpuCount > 0 {
		liveMetrics.GPUUsage = &interfaces.LiveResourceMetric{
			Current:            float64(gpuCount) * (0.85 + float64(time.Now().Unix()%10)/100),
			Average1Min:        float64(gpuCount) * 0.88,
			Average5Min:        float64(gpuCount) * 0.86,
			Peak:               float64(gpuCount) * 0.98,
			Allocated:          float64(gpuCount),
			UtilizationPercent: 88.0,
			Trend:              "stable",
			Unit:               "gpus",
		}
	}

	// Add network metrics
	liveMetrics.NetworkUsage = &interfaces.LiveResourceMetric{
		Current:            1200.0 + float64(time.Now().Unix()%200), // Mbps
		Average1Min:        1250.0,
		Average5Min:        1230.0,
		Peak:               1800.0,
		Allocated:          10000.0, // 10 Gbps
		UtilizationPercent: 12.5,
		Trend:              "increasing",
		Unit:               "mbps",
	}

	// Add I/O metrics
	liveMetrics.IOUsage = &interfaces.LiveResourceMetric{
		Current:            150.0 + float64(time.Now().Unix()%50), // MB/s
		Average1Min:        160.0,
		Average5Min:        155.0,
		Peak:               250.0,
		Allocated:          500.0, // 500 MB/s limit
		UtilizationPercent: 32.0,
		Trend:              "stable",
		Unit:               "MB/s",
	}

	// Add node-level metrics for each allocated node
	for i, nodeName := range job.Nodes {
		nodeMetrics := &interfaces.NodeLiveMetrics{
			NodeName:  nodeName,
			CPUCores:  job.CPUs / len(job.Nodes), // Distribute cores
			MemoryGB:  float64(job.Memory) / float64(len(job.Nodes)) / (1024 * 1024 * 1024),
			
			CPUUsage: &interfaces.LiveResourceMetric{
				Current:            75.0 + float64(i*5),
				Average1Min:        77.0 + float64(i*5),
				Average5Min:        76.0 + float64(i*5),
				Peak:               92.0 + float64(i*3),
				Allocated:          100.0,
				UtilizationPercent: 77.0 + float64(i*5),
				Trend:              "stable",
				Unit:               "percent",
			},
			
			MemoryUsage: &interfaces.LiveResourceMetric{
				Current:            65.0 + float64(i*3),
				Average1Min:        67.0 + float64(i*3),
				Average5Min:        66.0 + float64(i*3),
				Peak:               82.0 + float64(i*2),
				Allocated:          100.0,
				UtilizationPercent: 67.0 + float64(i*3),
				Trend:              "stable",
				Unit:               "percent",
			},
			
			LoadAverage: []float64{
				2.5 + float64(i)*0.2,
				2.3 + float64(i)*0.2,
				2.1 + float64(i)*0.2,
			},
			
			// Optional metrics
			CPUTemperature:   65.0 + float64(i*2) + float64(time.Now().Unix()%5),
			PowerConsumption: 250.0 + float64(i*20) + float64(time.Now().Unix()%30),
			NetworkInRate:    100.0 + float64(i*50),
			NetworkOutRate:   80.0 + float64(i*40),
			DiskReadRate:     50.0 + float64(i*10),
			DiskWriteRate:    30.0 + float64(i*8),
		}
		
		liveMetrics.NodeMetrics[nodeName] = nodeMetrics
	}

	// Check for alerts
	liveMetrics.Alerts = checkPerformanceAlerts(liveMetrics)

	// Add metadata
	liveMetrics.Metadata = map[string]interface{}{
		"version":           "v0.0.43",
		"collection_method": "simulated", // TODO: Change to "live" when real monitoring available
		"refresh_interval":  5,           // seconds
		"nodes_monitored":   len(job.Nodes),
	}

	return liveMetrics, nil
}

// Helper function to determine trend
func determineTrend(previous, current float64) string {
	diff := current - previous
	if diff > 2.0 {
		return "increasing"
	} else if diff < -2.0 {
		return "decreasing"
	}
	return "stable"
}

// Helper function to check for performance alerts
func checkPerformanceAlerts(metrics *interfaces.JobLiveMetrics) []interfaces.PerformanceAlert {
	alerts := []interfaces.PerformanceAlert{}

	// Check CPU usage
	if metrics.CPUUsage != nil && metrics.CPUUsage.UtilizationPercent > 90 {
		alerts = append(alerts, interfaces.PerformanceAlert{
			Type:              "warning",
			Category:          "cpu",
			Message:           "CPU utilization is above 90%",
			Severity:          "high",
			Timestamp:         time.Now(),
			CurrentValue:      metrics.CPUUsage.UtilizationPercent,
			ThresholdValue:    90.0,
			RecommendedAction: "Consider scaling up CPU resources if this persists",
		})
	}

	// Check memory usage
	if metrics.MemoryUsage != nil && metrics.MemoryUsage.UtilizationPercent > 85 {
		alerts = append(alerts, interfaces.PerformanceAlert{
			Type:              "warning",
			Category:          "memory",
			Message:           "Memory utilization is above 85%",
			Severity:          "medium",
			Timestamp:         time.Now(),
			CurrentValue:      metrics.MemoryUsage.UtilizationPercent,
			ThresholdValue:    85.0,
			RecommendedAction: "Monitor for potential out-of-memory conditions",
		})
	}

	// Check GPU usage if available
	if metrics.GPUUsage != nil && metrics.GPUUsage.UtilizationPercent > 95 {
		alerts = append(alerts, interfaces.PerformanceAlert{
			Type:              "warning",
			Category:          "gpu",
			Message:           "GPU utilization is above 95%",
			Severity:          "high",
			Timestamp:         time.Now(),
			CurrentValue:      metrics.GPUUsage.UtilizationPercent,
			ThresholdValue:    95.0,
			RecommendedAction: "GPU is near maximum capacity",
		})
	}

	// Check node-specific alerts
	for nodeName, nodeMetrics := range metrics.NodeMetrics {
		// Check CPU temperature
		if nodeMetrics.CPUTemperature > 80 {
			alerts = append(alerts, interfaces.PerformanceAlert{
				Type:              "critical",
				Category:          "cpu",
				Message:           fmt.Sprintf("High CPU temperature on node %s", nodeName),
				Severity:          "critical",
				Timestamp:         time.Now(),
				NodeName:          nodeName,
				ResourceName:      "cpu_temperature",
				CurrentValue:      nodeMetrics.CPUTemperature,
				ThresholdValue:    80.0,
				RecommendedAction: "Check cooling system and reduce load if necessary",
			})
		}

		// Check load average
		if len(nodeMetrics.LoadAverage) > 0 && nodeMetrics.LoadAverage[0] > float64(nodeMetrics.CPUCores)*1.5 {
			alerts = append(alerts, interfaces.PerformanceAlert{
				Type:              "warning",
				Category:          "cpu",
				Message:           fmt.Sprintf("High load average on node %s", nodeName),
				Severity:          "medium",
				Timestamp:         time.Now(),
				NodeName:          nodeName,
				ResourceName:      "load_average",
				CurrentValue:      nodeMetrics.LoadAverage[0],
				ThresholdValue:    float64(nodeMetrics.CPUCores) * 1.5,
				RecommendedAction: "System is overloaded, consider distributing work",
			})
		}
	}

	return alerts
}
