package v0_0_41

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
)

// JobAdapter implements the JobAdapter interface for v0.0.41
type JobAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewJobAdapter creates a new Job adapter for v0.0.41
func NewJobAdapter(client *api.ClientWithResponses) *JobAdapter {
	return &JobAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Job"),
		client:      client,
		wrapper:     nil, // We'll implement this later
	}
}

// List retrieves a list of jobs with optional filtering
func (a *JobAdapter) List(ctx context.Context, opts *types.JobListOptions) (*types.JobList, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmV0041GetJobsParams{}

	// Apply filters from options
	// Note: v0.0.41 GetJobs doesn't support filtering parameters
	// We'll need to filter the results after fetching all jobs
	_ = opts

	// Set flags to get detailed job information
	flags := api.SlurmV0041GetJobsParamsFlagsDETAIL
	params.Flags = &flags

	// Make the API call
	resp, err := a.client.SlurmV0041GetJobsWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list jobs")
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected nil response")
	}

	// Convert response to common types
	jobList := &types.JobList{
		Jobs: make([]types.Job, 0, len(resp.JSON200.Jobs)),
	}

	for _, apiJob := range resp.JSON200.Jobs {
		job, err := a.convertAPIJobToCommon(apiJob)
		if err != nil {
			// Log the error but continue processing other jobs
			continue
		}
		// Apply filters if options were provided
		if opts != nil {
			// Filter by account
			if len(opts.Accounts) > 0 {
				match := false
				for _, account := range opts.Accounts {
					if job.Account == account {
						match = true
						break
					}
				}
				if !match {
					continue
				}
			}
			
			// Filter by partition
			if len(opts.Partitions) > 0 {
				match := false
				for _, partition := range opts.Partitions {
					if job.Partition == partition {
						match = true
						break
					}
				}
				if !match {
					continue
				}
			}
			
			// Filter by state
			if len(opts.States) > 0 {
				match := false
				for _, state := range opts.States {
					if job.State == state {
						match = true
						break
					}
				}
				if !match {
					continue
				}
			}
		}
		
		jobList.Jobs = append(jobList.Jobs, *job)
	}

	// Extract warning messages if any
	if resp.JSON200.Warnings != nil {
		warnings := make([]string, 0, len(*resp.JSON200.Warnings))
		for _, warning := range *resp.JSON200.Warnings {
			if warning.Description != nil {
				warnings = append(warnings, *warning.Description)
			}
		}
		if len(warnings) > 0 {
			// JobList doesn't have a Meta field in common types
			// Warnings are being ignored for now
		}
	}

	// Extract error messages if any
	if resp.JSON200.Errors != nil {
		errors := make([]string, 0, len(*resp.JSON200.Errors))
		for _, error := range *resp.JSON200.Errors {
			if error.Description != nil {
				errors = append(errors, *error.Description)
			}
		}
		if len(errors) > 0 {
			// JobList doesn't have a Meta field in common types
			// Errors are being ignored for now
		}
	}

	return jobList, nil
}

// Get retrieves a specific job by ID
func (a *JobAdapter) Get(ctx context.Context, jobID int32) (*types.Job, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Validate job ID
	if jobID <= 0 {
		return nil, common.NewValidationError("job ID must be positive", "jobID", jobID)
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Set flags to get detailed job information
	flags := api.SlurmV0041GetJobParamsFlagsDETAIL
	params := &api.SlurmV0041GetJobParams{
		Flags: &flags,
	}

	// Make the API call
	jobIDStr := strconv.FormatUint(uint64(jobID), 10)
	resp, err := a.client.SlurmV0041GetJobWithResponse(ctx, jobIDStr, params)
	if err != nil {
		return nil, a.WrapError(err, fmt.Sprintf("failed to get job %d", jobID))
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}

	if resp.JSON200 == nil || len(resp.JSON200.Jobs) == 0 {
		return nil, a.HandleNotFound(fmt.Sprintf("job with ID %d", jobID))
	}

	// Convert the first job in the response
	job, err := a.convertAPIJobToCommon(resp.JSON200.Jobs[0])
	if err != nil {
		return nil, a.WrapError(err, fmt.Sprintf("failed to convert job %d", jobID))
	}

	return job, nil
}

// Submit submits a new job to the Slurm scheduler
func (a *JobAdapter) Submit(ctx context.Context, job *types.JobCreate) (*types.JobSubmitResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Validate job
	if job == nil {
		return nil, a.HandleValidationError("job cannot be nil")
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Create job description structure
	jobDesc := &api.V0041JobDescMsg{
		Script: &job.Script,
	}

	// Basic job properties
	if job.Name != "" {
		jobDesc.Name = &job.Name
	}
	if job.Account != "" {
		jobDesc.Account = &job.Account
	}
	if job.Partition != "" {
		jobDesc.Partition = &job.Partition
	}

	// Working directory
	if job.WorkingDirectory != "" {
		jobDesc.CurrentWorkingDirectory = &job.WorkingDirectory
	}

	// Standard output/error/input
	if job.StandardOutput != "" {
		jobDesc.StandardOutput = &job.StandardOutput
	}
	if job.StandardError != "" {
		jobDesc.StandardError = &job.StandardError
	}
	if job.StandardInput != "" {
		jobDesc.StandardInput = &job.StandardInput
	}

	// Time limit
	if job.TimeLimit > 0 {
		timeLimit := int32(job.TimeLimit)
		timeLimitStruct := &struct {
			Infinite *bool  `json:"infinite,omitempty"`
			Number   *int32 `json:"number,omitempty"`
			Set      *bool  `json:"set,omitempty"`
		}{
			Set:    &[]bool{true}[0],
			Number: &timeLimit,
		}
		jobDesc.TimeLimit = timeLimitStruct
	}

	// Node count
	if job.Nodes > 0 {
		nodes := int32(job.Nodes)
		jobDesc.MinimumNodes = &nodes
	}

	// Handle environment variables - CRITICAL for avoiding SLURM errors
	envList := make([]string, 0)
	
	// Always provide at least minimal environment to avoid SLURM write errors
	hasPath := false
	for key := range job.Environment {
		if key == "PATH" {
			hasPath = true
			break
		}
	}
	
	if !hasPath {
		envList = append(envList, "PATH=/usr/bin:/bin")
	}
	
	// Add all user-provided environment variables
	for key, value := range job.Environment {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}
	
	// Set environment in job description
	jobDesc.Environment = &envList

	// Create request body
	submitReq := api.V0041JobSubmitReq{
		Job: jobDesc,
	}

	// Make the API call
	resp, err := a.client.SlurmV0041PostJobSubmitWithResponse(ctx, submitReq)
	if err != nil {
		return nil, a.WrapError(err, "failed to submit job")
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected nil response")
	}

	// Convert response
	submitResp := &types.JobSubmitResponse{}

	// Extract job ID
	if resp.JSON200.JobId != nil {
		submitResp.JobID = *resp.JSON200.JobId
	}

	// Extract warnings if any
	if resp.JSON200.Warnings != nil {
		warnings := make([]string, 0, len(*resp.JSON200.Warnings))
		for _, warning := range *resp.JSON200.Warnings {
			if warning.Description != nil {
				warnings = append(warnings, *warning.Description)
			}
		}
		submitResp.Warning = warnings
	}

	// Extract errors if any
	if resp.JSON200.Errors != nil {
		errors := make([]string, 0, len(*resp.JSON200.Errors))
		for _, error := range *resp.JSON200.Errors {
			if error.Description != nil {
				errors = append(errors, *error.Description)
			}
		}
		if len(errors) > 0 {
			submitResp.Error = errors
		}
	}

	return submitResp, nil
}

// Cancel cancels a job
func (a *JobAdapter) Cancel(ctx context.Context, jobID int32, opts *types.JobCancelRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Validate job ID
	if jobID <= 0 {
		return a.HandleValidationError("jobID must be positive")
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Make the API call
	jobIDStr := strconv.FormatInt(int64(jobID), 10)
	params := &api.SlurmV0041DeleteJobParams{}
	resp, err := a.client.SlurmV0041DeleteJobWithResponse(ctx, jobIDStr, params)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to cancel job %d", jobID))
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
	}

	return nil
}

// Update updates job properties
func (a *JobAdapter) Update(ctx context.Context, jobID int32, update *types.JobUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Validate job ID
	if err := a.ValidateResourceID(strconv.FormatInt(int64(jobID), 10), "jobID"); err != nil {
		return err
	}

	// Validate update
	if update == nil {
		return a.HandleValidationError("job update cannot be nil")
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Convert update to API request
	updateReq := map[string]interface{}{
		"jobs": []map[string]interface{}{
			{
				"job_id": strconv.FormatInt(int64(jobID), 10),
			},
		},
	}
	
	// Add fields from update if provided
	if update.TimeLimit != nil {
		updateReq["jobs"].([]map[string]interface{})[0]["time_limit"] = *update.TimeLimit
	}
	if update.Priority != nil {
		updateReq["jobs"].([]map[string]interface{})[0]["priority"] = *update.Priority
	}

	// Note: v0.0.41 may not support job updates via API
	// This is a placeholder implementation
	_ = updateReq
	return fmt.Errorf("job updates not supported in v0.0.41 API")
}

// Watch watches for job state changes (not implemented in v0.0.41)
func (a *JobAdapter) Watch(ctx context.Context, opts *types.JobWatchOptions) (<-chan types.JobEvent, error) {
	return nil, fmt.Errorf("job watching is not supported in API v0.0.41")
}

// Signal sends a signal to a job (not implemented in v0.0.41)
func (a *JobAdapter) Signal(ctx context.Context, req *types.JobSignalRequest) error {
	return a.HandleNotImplemented("Signal", "v0.0.41")
}

// Hold holds or releases a job (not implemented in v0.0.41)
func (a *JobAdapter) Hold(ctx context.Context, req *types.JobHoldRequest) error {
	return a.HandleNotImplemented("Hold", "v0.0.41")
}

// Notify sends a notification to a job (not implemented in v0.0.41)
func (a *JobAdapter) Notify(ctx context.Context, req *types.JobNotifyRequest) error {
	return a.HandleNotImplemented("Notify", "v0.0.41")
}