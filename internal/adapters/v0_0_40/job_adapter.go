// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"context"
	"fmt"
	"strconv"

	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/jontk/slurm-client/pkg/errors"
)

// JobAdapter implements the JobAdapter interface for v0.0.40
type JobAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewJobAdapter creates a new Job adapter for v0.0.40
func NewJobAdapter(client *api.ClientWithResponses) *JobAdapter {
	return &JobAdapter{
		BaseManager: base.NewBaseManager("v0.0.40", "Job"),
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
	params := &api.SlurmV0040GetJobsParams{}

	// Apply filters from options
	// Note: SlurmV0040GetJobsParams only has UpdateTime and Flags fields in v0.0.40
	// Other filtering will need to be done client-side or through different endpoints
	if opts != nil {
		// The v0.0.40 API has limited parameter support for job listing
		// We'll need to filter results client-side after retrieval
		// Only set flags for now to get detailed job information
		if len(opts.JobIDs) > 0 || len(opts.States) > 0 || len(opts.Accounts) > 0 {
			// Use DETAIL flag to get comprehensive job information for filtering
			flags := api.SlurmV0040GetJobsParamsFlagsDETAIL
			params.Flags = &flags
		}
		// Set update time if available
		// Note: UpdateTime needs to be handled differently in v0.0.40
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0040GetJobsWithResponse(ctx, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.40"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "List Jobs"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Jobs, "List Jobs - jobs field"); err != nil {
		return nil, err
	}

	// Convert the response to common types
	jobList := make([]types.Job, 0, len(resp.JSON200.Jobs))
	for _, apiJob := range resp.JSON200.Jobs {
		job := a.convertAPIJobToCommon(apiJob)
		jobList = append(jobList, *job)
	}

	// Apply pagination
	listOpts := base.ListOptions{}
	if opts != nil {
		listOpts.Limit = opts.Limit
		listOpts.Offset = opts.Offset
	}

	// Apply pagination
	start := listOpts.Offset
	if start < 0 {
		start = 0
	}
	if start >= len(jobList) {
		return &types.JobList{
			Jobs:  []types.Job{},
			Total: len(jobList),
		}, nil
	}

	end := len(jobList)
	if listOpts.Limit > 0 {
		end = start + listOpts.Limit
		if end > len(jobList) {
			end = len(jobList)
		}
	}

	return &types.JobList{
		Jobs:  jobList[start:end],
		Total: len(jobList),
	}, nil
}

// Get retrieves a specific job by ID
func (a *JobAdapter) Get(ctx context.Context, jobID int32) (*types.Job, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.ValidateResourceID(jobID, "jobID"); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmV0040GetJobParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0040GetJobWithResponse(ctx, strconv.Itoa(int(jobID)), params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.40"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "Get Job"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Jobs, "Get Job - jobs field"); err != nil {
		return nil, err
	}

	// Check if we got any job entries
	if len(resp.JSON200.Jobs) == 0 {
		return nil, fmt.Errorf("job %d not found", jobID)
	}

	// Convert the first job (should be the only one)
	job := a.convertAPIJobToCommon(resp.JSON200.Jobs[0])

	return job, nil
}

// Submit submits a new job
func (a *JobAdapter) Submit(ctx context.Context, job *types.JobCreate) (*types.JobSubmitResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.validateJobCreate(job); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Create job submission structure
	jobDesc := &api.V0040JobDescMsg{
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
		timeLimit := int64(job.TimeLimit)
		setTrue := true
		jobDesc.TimeLimit = &api.V0040Uint32NoVal{
			Set:    &setTrue,
			Number: &timeLimit,
		}
	}

	// Node count
	if job.Nodes > 0 {
		nodes := job.Nodes
		jobDesc.MinimumNodes = &nodes
		jobDesc.MaximumNodes = &nodes
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

	// Set environment in job submission
	jobDesc.Environment = &envList

	// Create request body
	submitReq := api.V0040JobSubmitReq{
		Job: jobDesc,
	}

	// Make the API call
	resp, err := a.client.SlurmV0040PostJobSubmitWithResponse(ctx, submitReq)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.40"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "Submit Job"); err != nil {
		return nil, err
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

// Update updates an existing job
func (a *JobAdapter) Update(ctx context.Context, jobID int32, update *types.JobUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceID(jobID, "jobID"); err != nil {
		return err
	}
	if err := a.validateJobUpdate(update); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// First, check if job exists
	_, err := a.Get(ctx, jobID)
	if err != nil {
		return err
	}

	// Job update is not supported in v0.0.40 adapter
	// This would need to be implemented when the API supports it
	return fmt.Errorf("job update not supported in v0.0.40")
}

// Cancel cancels a job
func (a *JobAdapter) Cancel(ctx context.Context, jobID int32, opts *types.JobCancelRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceID(jobID, "jobID"); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Prepare parameters for the API call
	params := &api.SlurmV0040DeleteJobParams{}

	if opts != nil {
		if opts.Signal != "" {
			params.Signal = &opts.Signal
		}
		if opts.Message != "" {
			// Convert message to flags enum if needed
			flags := api.SlurmV0040DeleteJobParamsFlags(opts.Message)
			params.Flags = &flags
		}
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0040DeleteJobWithResponse(ctx, strconv.Itoa(int(jobID)), params)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	// Create adapter with special handling for 204 (No Content) status
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)

	// For DELETE operations, 204 is also a success
	if resp.StatusCode() == 204 {
		return nil
	}

	return common.HandleAPIResponse(responseAdapter, "v0.0.40")
}

// Signal sends a signal to a job
func (a *JobAdapter) Signal(ctx context.Context, req *types.JobSignalRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.validateJobSignalRequest(req); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// v0.0.40 doesn't have a dedicated signal endpoint, so we use cancel with signal
	cancelReq := &types.JobCancelRequest{
		Signal: req.Signal,
	}

	return a.Cancel(ctx, req.JobID, cancelReq)
}

// Hold holds or releases a job
func (a *JobAdapter) Hold(ctx context.Context, req *types.JobHoldRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.validateJobHoldRequest(req); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Get the job first to get current state
	_, err := a.Get(ctx, req.JobID)
	if err != nil {
		return err
	}

	// Prepare update with hold state
	update := &types.JobUpdate{
		Priority: func() *int32 {
			if req.Priority != 0 {
				priority := req.Priority
				return &priority
			}
			return nil
		}(),
	}

	return a.Update(ctx, req.JobID, update)
}

// Notify sends a notification to a job
func (a *JobAdapter) Notify(ctx context.Context, req *types.JobNotifyRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.validateJobNotifyRequest(req); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// For v0.0.40, we simulate notification by updating the job comment
	// This is a workaround as the API might not have dedicated notification endpoints
	update := &types.JobUpdate{
		Comment: &req.Message,
	}

	return a.Update(ctx, req.JobID, update)
}

// Requeue requeues a job (not available in v0.0.40)
func (a *JobAdapter) Requeue(ctx context.Context, jobID int32) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Requeue is not supported in v0.0.40
	return fmt.Errorf("requeue operation not supported in API v0.0.40")
}

// validateJobCreate validates job creation request
func (a *JobAdapter) validateJobCreate(job *types.JobCreate) error {
	if job == nil {
		return fmt.Errorf("job creation data is required")
	}
	if job.Command == "" && job.Script == "" {
		return fmt.Errorf("either command or script is required")
	}
	return nil
}

// validateJobUpdate validates job update request
func (a *JobAdapter) validateJobUpdate(update *types.JobUpdate) error {
	if update == nil {
		return fmt.Errorf("job update data is required")
	}
	// At least one field should be provided for update
	if update.Name == nil && update.Account == nil && update.Partition == nil &&
		update.QoS == nil && update.TimeLimit == nil && update.Priority == nil &&
		update.Nice == nil && update.Comment == nil {
		return fmt.Errorf("at least one field must be provided for update")
	}
	return nil
}

// validateJobSignalRequest validates job signal request
func (a *JobAdapter) validateJobSignalRequest(req *types.JobSignalRequest) error {
	if req == nil {
		return fmt.Errorf("job signal request is required")
	}
	if req.JobID == 0 {
		return fmt.Errorf("job ID is required")
	}
	if req.Signal == "" {
		return fmt.Errorf("signal is required")
	}
	return nil
}

// validateJobHoldRequest validates job hold request
func (a *JobAdapter) validateJobHoldRequest(req *types.JobHoldRequest) error {
	if req == nil {
		return fmt.Errorf("job hold request is required")
	}
	if req.JobID == 0 {
		return fmt.Errorf("job ID is required")
	}
	return nil
}

// validateJobNotifyRequest validates job notify request
func (a *JobAdapter) validateJobNotifyRequest(req *types.JobNotifyRequest) error {
	if req == nil {
		return fmt.Errorf("job notify request is required")
	}
	if req.JobID == 0 {
		return fmt.Errorf("job ID is required")
	}
	if req.Message == "" {
		return fmt.Errorf("message is required")
	}
	return nil
}

// Watch provides real-time job status updates (not supported in v0.0.40)
func (a *JobAdapter) Watch(ctx context.Context, opts *types.JobWatchOptions) (<-chan types.JobWatchEvent, error) {
	return nil, fmt.Errorf("watch functionality not supported in API v0.0.40")
}

// Allocate allocates resources for a job (not supported in v0.0.40)
func (a *JobAdapter) Allocate(ctx context.Context, req *types.JobAllocateRequest) (*types.JobAllocateResponse, error) {
	return nil, errors.NewNotImplementedError("Allocate", a.GetVersion())
}
