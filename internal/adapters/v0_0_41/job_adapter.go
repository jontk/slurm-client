package v0_0_41

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
	if opts != nil {
		if len(opts.Accounts) > 0 {
			accountStr := strings.Join(opts.Accounts, ",")
			params.Account = &accountStr
		}
		if len(opts.Partitions) > 0 {
			partitionStr := strings.Join(opts.Partitions, ",")
			params.Partition = &partitionStr
		}
		if len(opts.States) > 0 {
			var stateStrs []string
			for _, state := range opts.States {
				stateStrs = append(stateStrs, string(state))
			}
			stateStr := strings.Join(stateStrs, ",")
			params.State = &stateStr
		}
	}

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
		Meta: &types.ListMeta{
			Version: a.GetVersion(),
		},
	}

	for _, apiJob := range resp.JSON200.Jobs {
		job, err := a.convertAPIJobToCommon(apiJob)
		if err != nil {
			// Log the error but continue processing other jobs
			continue
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
			jobList.Meta.Warnings = warnings
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
			jobList.Meta.Errors = errors
		}
	}

	return jobList, nil
}

// Get retrieves a specific job by ID
func (a *JobAdapter) Get(ctx context.Context, jobID uint32) (*types.Job, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Validate job ID
	if err := a.ValidateResourceID("jobID", jobID); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Set flags to get detailed job information
	flags := api.SlurmV0041GetJobsParamsFlagsDETAIL
	params := &api.SlurmV0041GetJobsParams{
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
func (a *JobAdapter) Submit(ctx context.Context, opts *types.JobSubmitOptions) (*types.JobSubmitResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Validate options
	if opts == nil {
		return nil, a.HandleValidationError("job submit options cannot be nil")
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert options to API request
	submitReq := a.convertCommonToAPIJobSubmit(opts)

	// Make the API call
	resp, err := a.client.SlurmV0041PostJobSubmitWithResponse(ctx, *submitReq)
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
		submitResp.Warnings = warnings
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
			return submitResp, fmt.Errorf("job submission errors: %v", errors)
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
	resp, err := a.client.SlurmV0041DeleteJobWithResponse(ctx, jobIDStr)
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
func (a *JobAdapter) Update(ctx context.Context, jobID uint32, update *types.JobUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Validate job ID
	if err := a.ValidateResourceID("jobID", jobID); err != nil {
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
	updateReq := a.convertCommonToAPIJobUpdate(update)

	// Make the API call
	jobIDStr := strconv.FormatUint(uint64(jobID), 10)
	resp, err := a.client.SlurmV0041PostJobUpdateWithResponse(ctx, jobIDStr, *updateReq)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to update job %d", jobID))
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
	}

	return nil
}

// Watch watches for job state changes (not implemented in v0.0.41)
func (a *JobAdapter) Watch(ctx context.Context, opts *types.JobWatchOptions) (<-chan types.JobEvent, error) {
	return nil, fmt.Errorf("job watching is not supported in API v0.0.41")
}