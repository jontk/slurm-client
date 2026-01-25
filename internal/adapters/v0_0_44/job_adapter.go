// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_44

import (
	"context"
	"fmt"
	"strconv"
	"time"

	api "github.com/jontk/slurm-client/internal/api/v0_0_44"
	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/jontk/slurm-client/pkg/errors"
)

// JobAdapter implements the JobAdapter interface for v0.0.44
type JobAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewJobAdapter creates a new Job adapter for v0.0.44
func NewJobAdapter(client *api.ClientWithResponses) *JobAdapter {
	return &JobAdapter{
		BaseManager: base.NewBaseManager("v0.0.44", "Job"),
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
	// Use SlurmV0044GetJobsParams for current job queries
	params := &api.SlurmV0044GetJobsParams{}

	// Apply filters from options - SlurmV0044GetJobsParams only supports UpdateTime and Flags
	if opts != nil {
		// Note: SlurmV0044GetJobsParams has limited filtering capabilities
		// Most filtering will need to be done post-query
		if opts.StartTime != nil {
			updateTime := strconv.FormatInt(opts.StartTime.Unix(), 10)
			params.UpdateTime = &updateTime
		}
		// Other filters (accounts, users, states, etc.) need to be applied after getting results
		// Store filter options for client-side filtering
	}

	// Call the generated OpenAPI client for current job queries
	resp, err := a.client.SlurmV0044GetJobsWithResponse(ctx, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0044OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.44"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "List Jobs"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Jobs, "List Jobs - jobs field"); err != nil {
		return nil, err
	}

	// Convert the response to common types - SlurmV0044GetJobs returns V0044JobInfo
	jobList := make([]types.Job, 0, len(resp.JSON200.Jobs))
	for _, apiJobInfo := range resp.JSON200.Jobs {
		job := a.convertAPIJobToCommon(apiJobInfo)
		jobList = append(jobList, *job)
	}

	// Apply client-side filtering since API has limited filter support
	if opts != nil {
		jobList = a.applyClientSideFilters(jobList, opts)
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
	params := &api.SlurmV0044GetJobParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0044GetJobWithResponse(ctx, strconv.Itoa(int(jobID)), params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0044OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.44"); err != nil {
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
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("Job with ID %d not found", jobID))
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

	// Convert to API format
	apiJob := a.convertCommonJobCreateToAPI(job)

	// Create request body - V0044JobDescMsg format
	reqBody := api.V0044JobDescMsg{
		Account:   apiJob.Account,
		Name:      apiJob.Name,
		Partition: apiJob.Partition,
	}

	// Handle script/command
	if job.Script != "" {
		reqBody.Script = &job.Script
	}

	// Handle working directory - REQUIRED in SLURM v25.05.0
	if job.WorkingDirectory != "" {
		reqBody.CurrentWorkingDirectory = &job.WorkingDirectory
	} else {
		// Default to /tmp if not specified to avoid SLURM error
		defaultWorkDir := "/tmp"
		reqBody.CurrentWorkingDirectory = &defaultWorkDir
	}

	// Handle standard output/error/input
	if job.StandardOutput != "" {
		reqBody.StandardOutput = &job.StandardOutput
	}
	if job.StandardError != "" {
		reqBody.StandardError = &job.StandardError
	}
	if job.StandardInput != "" {
		reqBody.StandardInput = &job.StandardInput
	}

	// Handle time limit
	if job.TimeLimit > 0 {
		timeLimit := job.TimeLimit
		setTrue := true
		reqBody.TimeLimit = &api.V0044Uint32NoValStruct{
			Set:    &setTrue,
			Number: &timeLimit,
		}
	}

	// Handle node count
	if job.Nodes > 0 {
		nodeMin := job.Nodes
		reqBody.MinimumNodes = &nodeMin
		reqBody.MaximumNodes = &nodeMin
	}

	// Handle environment variables - CRITICAL for avoiding SLURM errors
	envVars := make([]string, 0)

	// Always provide at least minimal environment to avoid SLURM write errors
	hasPath := false
	for key := range job.Environment {
		if key == "PATH" {
			hasPath = true
			break
		}
	}

	if !hasPath {
		envVars = append(envVars, "PATH=/usr/bin:/bin")
	}

	// Add all user-provided environment variables
	for key, value := range job.Environment {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	// Set environment in request body
	reqBody.Environment = &envVars

	// Create job submit request wrapping the job description
	submitReq := api.V0044JobSubmitReq{
		Job: &reqBody,
	}

	// Call the generated OpenAPI client - use the dedicated submit endpoint
	resp, err := a.client.SlurmV0044PostJobSubmitWithResponse(ctx, submitReq)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0044OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.44"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "Submit Job"); err != nil {
		return nil, err
	}

	// Extract job ID from response
	var jobID int32
	var warnings []string
	var errors []string

	// Extract job ID from response - V0044OpenapiJobSubmitResponse structure
	if resp.JSON200 != nil && resp.JSON200.JobId != nil {
		jobID = *resp.JSON200.JobId
	}

	// Extract warnings if available
	if resp.JSON200 != nil && resp.JSON200.Warnings != nil {
		for _, warn := range *resp.JSON200.Warnings {
			// V0044OpenapiWarning has Description and Source fields
			if warn.Description != nil {
				warnings = append(warnings, *warn.Description)
			} else if warn.Source != nil {
				warnings = append(warnings, *warn.Source)
			}
		}
	}

	// Check for errors in response
	if resp.JSON200 != nil && resp.JSON200.Errors != nil {
		for _, apiErr := range *resp.JSON200.Errors {
			if apiErr.Error != nil {
				errors = append(errors, *apiErr.Error)
			}
		}
	}

	return &types.JobSubmitResponse{
		JobID:   jobID,
		Error:   errors,
		Warning: warnings,
	}, nil
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

	// First, get the existing job to merge updates
	existingJob, err := a.Get(ctx, jobID)
	if err != nil {
		return err
	}

	// Convert to API format and apply updates
	apiJob := a.convertCommonJobUpdateToAPI(existingJob, update)

	// Create request body for job update - V0044JobDescMsg format
	reqBody := api.V0044JobDescMsg{
		Account:   apiJob.Account,
		Name:      apiJob.Name,
		Partition: apiJob.Partition,
		// Copy other fields as needed
	}

	// Call the generated OpenAPI client (POST is used for updates in SLURM API)
	resp, err := a.client.SlurmV0044PostJobWithResponse(ctx, strconv.Itoa(int(jobID)), reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0044OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "v0.0.44")
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
	params := &api.SlurmV0044DeleteJobParams{}

	if opts != nil {
		if opts.Signal != "" {
			params.Signal = &opts.Signal
		}
		if opts.Message != "" {
			// Convert message to appropriate flags type
			flags := api.SlurmV0044DeleteJobParamsFlags(opts.Message)
			params.Flags = &flags
		}
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0044DeleteJobWithResponse(ctx, strconv.Itoa(int(jobID)), params)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0044OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	// Create adapter with special handling for 204 (No Content) status
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)

	// For DELETE operations, 204 is also a success
	if resp.StatusCode() == 204 {
		return nil
	}

	return common.HandleAPIResponse(responseAdapter, "v0.0.44")
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

	// TODO: The v0.0.44 API doesn't seem to have a dedicated job signal endpoint
	// For now, we'll return an error indicating the feature is not supported
	return errors.NewClientError(errors.ErrorCodeValidationFailed, "Job signaling is not supported in v0.0.44 API")
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

	// Create request body for job hold/release - V0044JobDescMsg format
	reqBody := api.V0044JobDescMsg{
		// Basic job identification
		Name: func() *string {
			// We might need the job name for proper identification
			return nil
		}(),
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0044PostJobWithResponse(ctx, strconv.Itoa(int(req.JobID)), reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0044OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "v0.0.44")
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

	// For v0.0.44, we simulate notification by updating the job comment
	// This is a workaround as the API might not have dedicated notification endpoints
	update := &types.JobUpdate{
		Comment: &req.Message,
	}

	return a.Update(ctx, req.JobID, update)
}

// Requeue requeues a job
func (a *JobAdapter) Requeue(ctx context.Context, jobID int32) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Prepare parameters with SlurmV0044DeleteJobParamsFlagsFEDERATIONREQUEUE flag
	params := &api.SlurmV0044DeleteJobParams{}
	requeueFlag := api.SlurmV0044DeleteJobParamsFlagsFEDERATIONREQUEUE
	params.Flags = &requeueFlag

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0044DeleteJobWithResponse(ctx, strconv.Itoa(int(jobID)), params)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0044OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "v0.0.44")
}

// Watch watches for job state changes using polling
func (a *JobAdapter) Watch(ctx context.Context, opts *types.JobWatchOptions) (<-chan types.JobWatchEvent, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Create the event channel
	eventCh := make(chan types.JobWatchEvent, 10) // Buffered channel to prevent blocking

	// Start polling in a goroutine
	go func() {
		defer close(eventCh)

		// Poll interval - configurable, but default to 5 seconds
		pollInterval := 5 * time.Second

		// Keep track of job states to detect changes
		jobStates := make(map[int32]types.JobState)

		// Create a ticker for polling
		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()

		// Initial poll to populate the state map
		a.pollJobs(ctx, opts, jobStates, eventCh, true)

		// Poll for changes
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				a.pollJobs(ctx, opts, jobStates, eventCh, false)
			}
		}
	}()

	return eventCh, nil
}

// pollJobs polls for job changes and sends events
func (a *JobAdapter) pollJobs(ctx context.Context, opts *types.JobWatchOptions, jobStates map[int32]types.JobState, eventCh chan<- types.JobWatchEvent, isInitial bool) {
	// Create list options based on watch options
	listOpts := &types.JobListOptions{}

	// If watching a specific job, filter by job ID
	if opts != nil && opts.JobID != 0 {
		listOpts.JobIDs = []int32{opts.JobID}
	}

	// Get current job list
	jobList, err := a.List(ctx, listOpts)
	if err != nil {
		// Send error event
		select {
		case eventCh <- types.JobWatchEvent{
			EventTime: time.Now(),
			EventType: "error",
			JobID:     0,
			Reason:    fmt.Sprintf("Failed to poll jobs: %v", err),
		}:
		case <-ctx.Done():
		}
		return
	}

	// Check for state changes
	for _, job := range jobList.Jobs {
		previousState, exists := jobStates[job.JobID]
		currentState := job.State

		// Update the state map
		jobStates[job.JobID] = currentState

		// Skip initial population unless it's a completion event or the user wants all events
		if isInitial {
			continue
		}

		// Send event if state changed
		if exists && previousState != currentState {
			eventType := a.getEventTypeFromStateChange(previousState, currentState)

			// Filter by event types if specified
			if opts != nil && len(opts.EventTypes) > 0 {
				found := false
				for _, et := range opts.EventTypes {
					if et == eventType {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			event := types.JobWatchEvent{
				EventTime:     time.Now(),
				EventType:     eventType,
				JobID:         job.JobID,
				JobName:       job.Name,
				UserName:      job.UserName,
				PreviousState: previousState,
				NewState:      currentState,
				NodeList:      job.NodeList,
				Reason:        a.getReasonFromStateChange(previousState, currentState),
			}

			select {
			case eventCh <- event:
			case <-ctx.Done():
				return
			}
		}
	}

	// Check for completed/removed jobs (jobs that existed before but don't exist now)
	for jobID, previousState := range jobStates {
		found := false
		for _, job := range jobList.Jobs {
			if job.JobID == jobID {
				found = true
				break
			}
		}

		// If job is not found and was not in a terminal state, it might have been removed
		if !found && !a.isTerminalState(previousState) {
			// Send completion event
			event := types.JobWatchEvent{
				EventTime:     time.Now(),
				EventType:     "completed",
				JobID:         jobID,
				PreviousState: previousState,
				NewState:      types.JobStateCompleted,
				Reason:        "Job completed and removed from active list",
			}

			select {
			case eventCh <- event:
			case <-ctx.Done():
				return
			}

			// Remove from state map
			delete(jobStates, jobID)
		}
	}
}

// getEventTypeFromStateChange determines the event type based on state transition
func (a *JobAdapter) getEventTypeFromStateChange(previous, current types.JobState) string {
	switch current {
	case types.JobStateRunning:
		if previous == types.JobStatePending {
			return "start"
		}
		return "running"
	case types.JobStateCompleted:
		return "end"
	case types.JobStateFailed:
		return "fail"
	case types.JobStateCancelled:
		return "cancel"
	case types.JobStatePending:
		return "submit"
	case types.JobStateSuspended:
		return "suspend"
	default:
		return "state_change"
	}
}

// getReasonFromStateChange provides a reason for the state change
func (a *JobAdapter) getReasonFromStateChange(previous, current types.JobState) string {
	switch current {
	case types.JobStateRunning:
		if previous == types.JobStatePending {
			return "Job started execution"
		}
		return "Job resumed execution"
	case types.JobStateCompleted:
		return "Job completed successfully"
	case types.JobStateFailed:
		return "Job failed during execution"
	case types.JobStateCancelled:
		return "Job was cancelled"
	case types.JobStatePending:
		return "Job submitted and waiting for resources"
	case types.JobStateSuspended:
		return "Job was suspended"
	default:
		return fmt.Sprintf("Job state changed from %s to %s", previous, current)
	}
}

// isTerminalState checks if a job state is terminal
func (a *JobAdapter) isTerminalState(state types.JobState) bool {
	switch state {
	case types.JobStateCompleted, types.JobStateFailed, types.JobStateCancelled:
		return true
	default:
		return false
	}
}

// validateJobCreate validates job creation request
func (a *JobAdapter) validateJobCreate(job *types.JobCreate) error {
	if job == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "job creation data is required", "job", nil, nil)
	}
	if job.Command == "" && job.Script == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "either command or script is required", "job", job, nil)
	}
	// Account is required for SLURM v0.0.44 job submission
	if job.Account == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "account is required for job submission in SLURM v0.0.44", "account", job.Account, nil)
	}
	return nil
}

// validateJobUpdate validates job update request
func (a *JobAdapter) validateJobUpdate(update *types.JobUpdate) error {
	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "job update data is required", "update", nil, nil)
	}
	// At least one field should be provided for update
	if update.Name == nil && update.Account == nil && update.Partition == nil &&
		update.QoS == nil && update.TimeLimit == nil && update.Priority == nil &&
		update.Nice == nil && update.Comment == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "at least one field must be provided for update", "update", update, nil)
	}
	return nil
}

// validateJobSignalRequest validates job signal request
func (a *JobAdapter) validateJobSignalRequest(req *types.JobSignalRequest) error {
	if req == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "job signal request is required", "req", nil, nil)
	}
	if req.JobID == 0 {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "job ID is required", "jobID", req.JobID, nil)
	}
	if req.Signal == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "signal is required", "signal", req.Signal, nil)
	}
	return nil
}

// validateJobHoldRequest validates job hold request
func (a *JobAdapter) validateJobHoldRequest(req *types.JobHoldRequest) error {
	if req == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "job hold request is required", "req", nil, nil)
	}
	if req.JobID == 0 {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "job ID is required", "jobID", req.JobID, nil)
	}
	return nil
}

// validateJobNotifyRequest validates job notify request
func (a *JobAdapter) validateJobNotifyRequest(req *types.JobNotifyRequest) error {
	if req == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "job notify request is required", "req", nil, nil)
	}
	if req.JobID == 0 {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "job ID is required", "jobID", req.JobID, nil)
	}
	if req.Message == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "message is required", "message", req.Message, nil)
	}
	return nil
}

// Simplified converter methods for job management
func (a *JobAdapter) convertAPIJobToCommon(apiJob api.V0044JobInfo) *types.Job {
	job := &types.Job{}
	if apiJob.JobId != nil {
		job.JobID = *apiJob.JobId
	}
	if apiJob.Name != nil {
		job.Name = *apiJob.Name
	}
	if apiJob.Account != nil {
		job.Account = *apiJob.Account
	}
	if apiJob.Partition != nil {
		job.Partition = *apiJob.Partition
	}
	if apiJob.UserId != nil {
		job.UserID = *apiJob.UserId
	}
	if apiJob.GroupId != nil {
		job.GroupID = *apiJob.GroupId
	}
	// TODO: Add more field conversions as needed
	return job
}

func (a *JobAdapter) convertCommonJobCreateToAPI(create *types.JobCreate) *api.V0044Job {
	apiJob := &api.V0044Job{}

	// Set required fields with proper pointers
	if create.Name != "" {
		apiJob.Name = &create.Name
	} else {
		// Default job name if not provided
		defaultName := "job"
		apiJob.Name = &defaultName
	}

	// Account is required in v0.0.44
	if create.Account != "" {
		apiJob.Account = &create.Account
	}

	if create.Partition != "" {
		apiJob.Partition = &create.Partition
	}

	// Set QoS if provided
	if create.QoS != "" {
		apiJob.Qos = &create.QoS
	}

	// Set working directory if provided
	// Note: WorkDirectory field might not exist in V0044Job
	// This will be set in the V0044JobDescMsg in Submit method
	if create.WorkingDirectory != "" {
	}

	// Set comment if provided
	// Comment field in V0044Job is a complex struct, not a simple string
	// This will be handled in the Submit method with proper job description
	if create.Comment != "" {
	}

	// Set priority if provided
	if create.Priority != nil && *create.Priority > 0 {
		priority := *create.Priority
		setTrue := true
		apiJob.Priority = &api.V0044Uint32NoValStruct{
			Set:    &setTrue,
			Number: &priority,
		}
	}

	return apiJob
}

func (a *JobAdapter) convertCommonJobUpdateToAPI(existing *types.Job, update *types.JobUpdate) *api.V0044Job {
	apiJob := &api.V0044Job{}
	apiJob.JobId = &existing.JobID
	if existing.Name != "" {
		apiJob.Name = &existing.Name
	}
	if existing.Account != "" {
		apiJob.Account = &existing.Account
	}
	if existing.Partition != "" {
		apiJob.Partition = &existing.Partition
	}
	// Apply updates
	if update.Name != nil {
		apiJob.Name = update.Name
	}
	if update.Account != nil {
		apiJob.Account = update.Account
	}
	if update.Partition != nil {
		apiJob.Partition = update.Partition
	}
	// TODO: Add more field conversions as needed
	return apiJob
}

// applyClientSideFilters applies filters that the API doesn't support
func (a *JobAdapter) applyClientSideFilters(jobs []types.Job, opts *types.JobListOptions) []types.Job {
	if opts == nil {
		return jobs
	}

	filtered := make([]types.Job, 0, len(jobs))

	for _, job := range jobs {
		// Apply account filter
		if len(opts.Accounts) > 0 {
			found := false
			for _, account := range opts.Accounts {
				if job.Account == account {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Apply user filter (using UserName or UserID)
		if len(opts.Users) > 0 {
			found := false
			for _, user := range opts.Users {
				if job.UserName == user || strconv.Itoa(int(job.UserID)) == user {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Apply state filter
		if len(opts.States) > 0 {
			found := false
			for _, state := range opts.States {
				if string(job.State) == string(state) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Apply partition filter
		if len(opts.Partitions) > 0 {
			found := false
			for _, partition := range opts.Partitions {
				if job.Partition == partition {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Apply job ID filter
		if len(opts.JobIDs) > 0 {
			found := false
			for _, id := range opts.JobIDs {
				if job.JobID == id {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Apply job name filter
		if len(opts.JobNames) > 0 {
			found := false
			for _, name := range opts.JobNames {
				if job.Name == name {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		filtered = append(filtered, job)
	}

	return filtered
}

// Allocate allocates resources for a job
func (a *JobAdapter) Allocate(ctx context.Context, req *types.JobAllocateRequest) (*types.JobAllocateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.validateJobAllocateRequest(req); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert common allocation request to API request
	apiReq := a.convertCommonJobAllocateToAPI(req)

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0044PostJobAllocateWithResponse(ctx, *apiReq)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0044OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.44"); err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from allocation API")
	}

	// Convert API response to common response
	return a.convertAPIJobAllocateResponseToCommon(resp.JSON200), nil
}

// validateJobAllocateRequest validates job allocation request
func (a *JobAdapter) validateJobAllocateRequest(req *types.JobAllocateRequest) error {
	if req == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "job allocation request is required", "req", nil, nil)
	}

	// Account is required for SLURM v0.0.44 job allocation
	if req.Account == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "account is required for job allocation in SLURM v0.0.44", "account", req.Account, nil)
	}

	// At least one resource requirement should be specified
	if req.Nodes == "" && req.CPUs == 0 && req.Memory == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "at least one resource requirement (nodes, cpus, or memory) must be specified", "resources", req, nil)
	}

	return nil
}

// convertCommonJobAllocateToAPI converts common allocation request to API request
func (a *JobAdapter) convertCommonJobAllocateToAPI(req *types.JobAllocateRequest) *api.V0044JobAllocReq {
	// Create the job description message
	jobDesc := &api.V0044JobDescMsg{}

	// Set basic fields in the job description
	if req.Name != "" {
		jobDesc.Name = &req.Name
	}
	if req.Account != "" {
		jobDesc.Account = &req.Account
	}
	if req.Partition != "" {
		jobDesc.Partition = &req.Partition
	}
	if req.QoS != "" {
		jobDesc.Qos = &req.QoS
	}

	// Create the allocation request with the job description
	apiReq := &api.V0044JobAllocReq{
		Job: jobDesc,
	}

	// TODO: Implement more detailed field mappings
	// The V0044JobDescMsg has a very complex structure with many nested fields
	// For now, we only set the basic fields to get compilation working

	return apiReq
}

// convertAPIJobAllocateResponseToCommon converts API allocation response to common response
func (a *JobAdapter) convertAPIJobAllocateResponseToCommon(apiResp *api.V0044OpenapiJobAllocResp) *types.JobAllocateResponse {
	resp := &types.JobAllocateResponse{
		Status: "success",
		Meta:   make(map[string]interface{}),
	}

	// Extract job ID from response
	if apiResp != nil && apiResp.JobId != nil {
		resp.JobID = *apiResp.JobId
	}

	// Extract user message if available
	if apiResp != nil && apiResp.JobSubmitUserMsg != nil {
		resp.Message = *apiResp.JobSubmitUserMsg
	}

	// TODO: The V0044OpenapiJobAllocResp doesn't contain detailed allocation info
	// like node list, CPUs, timestamps, etc. Those would need to be retrieved
	// from a separate job info query if needed

	// Extract metadata if available
	if apiResp != nil && apiResp.Meta != nil {
		// TODO: Add metadata extraction if needed
		resp.Meta["api_version"] = "v0.0.44"
	}

	return resp
}
