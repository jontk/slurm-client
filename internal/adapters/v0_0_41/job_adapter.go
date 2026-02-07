// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"context"
	"fmt"
	"strconv"

	types "github.com/jontk/slurm-client/api"
	adapterbase "github.com/jontk/slurm-client/internal/adapters/base"
	"github.com/jontk/slurm-client/internal/common"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
	"github.com/jontk/slurm-client/pkg/errors"
)

// JobAdapter implements the JobAdapter interface for v0.0.41
type JobAdapter struct {
	*adapterbase.BaseManager
	client *api.ClientWithResponses
}

// NewJobAdapter creates a new Job adapter for v0.0.41
func NewJobAdapter(client *api.ClientWithResponses) *JobAdapter {
	return &JobAdapter{
		BaseManager: adapterbase.NewBaseManager("v0.0.41", "Job"),
		client:      client,
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
		if opts != nil && !a.jobPassesFilters(job, opts) {
			continue
		}
		jobList.Jobs = append(jobList.Jobs, *job)
	}
	// Note: JobList doesn't have a Meta field in common types
	// Warnings and errors from the response are being ignored for now
	return jobList, nil
}

// jobPassesFilters checks if a job matches all provided filters
func (a *JobAdapter) jobPassesFilters(job *types.Job, opts *types.JobListOptions) bool {
	// Filter by account
	if len(opts.Accounts) > 0 && !a.jobAccountMatches(job.Account, opts.Accounts) {
		return false
	}
	// Filter by partition
	if len(opts.Partitions) > 0 && !a.jobPartitionMatches(job.Partition, opts.Partitions) {
		return false
	}
	// Filter by state
	if len(opts.States) > 0 && !a.jobStateMatches(job.JobState, opts.States) {
		return false
	}
	return true
}

// jobAccountMatches checks if a job account is in the filter list
func (a *JobAdapter) jobAccountMatches(account *string, accounts []string) bool {
	if account == nil {
		return false
	}
	for _, acc := range accounts {
		if acc == *account {
			return true
		}
	}
	return false
}

// jobPartitionMatches checks if a job partition is in the filter list
func (a *JobAdapter) jobPartitionMatches(partition *string, partitions []string) bool {
	if partition == nil {
		return false
	}
	for _, p := range partitions {
		if p == *partition {
			return true
		}
	}
	return false
}

// jobStateMatches checks if a job state is in the filter list
func (a *JobAdapter) jobStateMatches(jobStates []types.JobState, states []types.JobState) bool {
	for _, js := range jobStates {
		for _, s := range states {
			if s == js {
				return true
			}
		}
	}
	return false
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

	// Convert to API request
	reqBody, err := a.convertCommonJobCreateToAPI(job)
	if err != nil {
		return nil, a.WrapError(err, "failed to convert job request")
	}

	// Make the API call
	resp, err := a.client.SlurmV0041PostJobSubmitWithResponse(ctx, reqBody)
	if err != nil {
		return nil, a.WrapError(err, "failed to submit job")
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected nil response from job submit")
	}

	// Check for errors in response
	if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		errMsgs := make([]string, 0, len(*resp.JSON200.Errors))
		for _, apiErr := range *resp.JSON200.Errors {
			if apiErr.Error != nil {
				errMsgs = append(errMsgs, *apiErr.Error)
			}
		}
		if len(errMsgs) > 0 {
			return nil, fmt.Errorf("job submission failed: %v", errMsgs)
		}
	}

	// Build response
	response := &types.JobSubmitResponse{}
	if resp.JSON200.JobId != nil {
		response.JobId = *resp.JSON200.JobId
	}

	return response, nil
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

	// Convert to API request using JSON marshaling workaround
	reqBody, err := a.convertJobUpdateToAPI(update)
	if err != nil {
		return a.WrapError(err, "failed to convert job update request")
	}

	// Make the API call
	jobIDStr := strconv.FormatInt(int64(jobID), 10)
	resp, err := a.client.SlurmV0041PostJobWithResponse(ctx, jobIDStr, reqBody)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to update job %d", jobID))
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return err
	}

	// Check for errors in response
	if resp.JSON200 != nil && resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		errMsgs := make([]string, 0)
		for _, apiErr := range *resp.JSON200.Errors {
			if apiErr.Error != nil {
				errMsgs = append(errMsgs, *apiErr.Error)
			}
		}
		if len(errMsgs) > 0 {
			return fmt.Errorf("job update failed: %v", errMsgs)
		}
	}

	return nil
}

// Watch watches for job state changes (not implemented in v0.0.41)
func (a *JobAdapter) Watch(ctx context.Context, opts *types.JobWatchOptions) (<-chan types.JobEvent, error) {
	return nil, errors.NewNotImplementedError("Watch Jobs", "v0.0.41")
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

// Requeue requeues a job (not available in v0.0.41)
func (a *JobAdapter) Requeue(ctx context.Context, jobID int32) error {
	return a.HandleNotImplemented("Requeue", "v0.0.41")
}

// Allocate allocates resources for a job (not supported in v0.0.41)
func (a *JobAdapter) Allocate(ctx context.Context, req *types.JobAllocateRequest) (*types.JobAllocateResponse, error) {
	return nil, a.HandleNotImplemented("Allocate", "v0.0.41")
}
