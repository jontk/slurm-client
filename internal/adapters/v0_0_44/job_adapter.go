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

	// DEBUG: Log the request details
	fmt.Printf("[DEBUG JobAdapter] Calling SlurmV0044GetJobsWithResponse with params: %+v\n", params)

	// Call the generated OpenAPI client for current job queries
	resp, err := a.client.SlurmV0044GetJobsWithResponse(ctx, params)
	if err != nil {
		fmt.Printf("[DEBUG JobAdapter] API call failed with error: %v\n", err)
		return nil, a.HandleAPIError(err)
	}

	// DEBUG: Log the response status
	fmt.Printf("[DEBUG JobAdapter] API response status: %d\n", resp.StatusCode())

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
		fmt.Printf("[DEBUG JobAdapter] CheckNilResponse failed for JSON200: %v\n", err)
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Jobs, "List Jobs - jobs field"); err != nil {
		fmt.Printf("[DEBUG JobAdapter] CheckNilResponse failed for Jobs field: %v\n", err)
		return nil, err
	}

	// DEBUG: Log the response data
	fmt.Printf("[DEBUG JobAdapter] Response contains %d jobs\n", len(resp.JSON200.Jobs))

	// Convert the response to common types - SlurmV0044GetJobs returns V0044JobInfo
	jobList := make([]types.Job, 0, len(resp.JSON200.Jobs))
	for _, apiJobInfo := range resp.JSON200.Jobs {
		job := a.convertAPIJobToCommon(apiJobInfo)
		jobList = append(jobList, *job)
	}

	// Apply client-side filtering since API has limited filter support
	if opts != nil {
		fmt.Printf("[DEBUG JobAdapter] Before applyClientSideFilters: %d jobs\n", len(jobList))
		fmt.Printf("[DEBUG JobAdapter] Filter options: Accounts=%v, Users=%v, States=%v, Partitions=%v, JobIDs=%v, JobNames=%v\n",
			opts.Accounts, opts.Users, opts.States, opts.Partitions, opts.JobIDs, opts.JobNames)
		jobList = a.applyClientSideFilters(jobList, opts)
		fmt.Printf("[DEBUG JobAdapter] After applyClientSideFilters: %d jobs\n", len(jobList))
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
// buildEnvironmentListV0044 builds the environment variable list with defaults
func (a *JobAdapter) buildEnvironmentListV0044(jobEnv map[string]string) []string {
	envList := make([]string, 0)

	// Always provide at least minimal environment to avoid SLURM write errors
	hasPath := false
	for key := range jobEnv {
		if key == "PATH" {
			hasPath = true
			break
		}
	}

	if !hasPath {
		envList = append(envList, "PATH=/usr/bin:/bin")
	}

	// Add all user-provided environment variables
	for key, value := range jobEnv {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}

	return envList
}

// setV0044JobIOProperties sets I/O properties for v0.0.44
func (a *JobAdapter) setV0044JobIOProperties(reqBody *api.V0044JobDescMsg, job *types.JobCreate) {
	if job.Script != "" {
		reqBody.Script = &job.Script
	}

	if job.WorkingDirectory != "" {
		reqBody.CurrentWorkingDirectory = &job.WorkingDirectory
	} else {
		defaultWorkDir := "/tmp"
		reqBody.CurrentWorkingDirectory = &defaultWorkDir
	}

	if job.StandardOutput != "" {
		reqBody.StandardOutput = &job.StandardOutput
	}
	if job.StandardError != "" {
		reqBody.StandardError = &job.StandardError
	}
	if job.StandardInput != "" {
		reqBody.StandardInput = &job.StandardInput
	}
}

// setV0044JobResources sets resource properties for v0.0.44
func (a *JobAdapter) setV0044JobResources(reqBody *api.V0044JobDescMsg, job *types.JobCreate) {
	if job.TimeLimit > 0 {
		timeLimit := job.TimeLimit
		setTrue := true
		reqBody.TimeLimit = &api.V0044Uint32NoValStruct{
			Set:    &setTrue,
			Number: &timeLimit,
		}
	}
	if job.Nodes > 0 {
		nodeMin := job.Nodes
		reqBody.MinimumNodes = &nodeMin
		reqBody.MaximumNodes = &nodeMin
	}
}

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

	// Create and populate request body
	reqBody := api.V0044JobDescMsg{
		Account:   apiJob.Account,
		Name:      apiJob.Name,
		Partition: apiJob.Partition,
	}
	a.setV0044JobIOProperties(&reqBody, job)
	a.setV0044JobResources(&reqBody, job)

	// Build and set environment variables
	envVars := a.buildEnvironmentListV0044(job.Environment)
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
// matchesEventTypeFilter checks if an event type passes the filter
func (a *JobAdapter) matchesEventTypeFilter(eventType string, opts *types.JobWatchOptions) bool {
	if opts == nil || len(opts.EventTypes) == 0 {
		return true
	}
	for _, et := range opts.EventTypes {
		if et == eventType {
			return true
		}
	}
	return false
}

// sendEventIfOpen sends an event on the channel if the context is not closed
func (a *JobAdapter) sendEventIfOpen(ctx context.Context, eventCh chan<- types.JobWatchEvent, event types.JobWatchEvent) bool {
	select {
	case eventCh <- event:
		return true
	case <-ctx.Done():
		return false
	}
}

// processJobStateChange handles a job state change and sends appropriate event
func (a *JobAdapter) processJobStateChange(ctx context.Context, job types.Job, previousState types.JobState, eventCh chan<- types.JobWatchEvent, opts *types.JobWatchOptions) bool {
	eventType := a.getEventTypeFromStateChange(previousState, job.State)

	// Filter by event types if specified
	if !a.matchesEventTypeFilter(eventType, opts) {
		return true
	}

	event := types.JobWatchEvent{
		EventTime:     time.Now(),
		EventType:     eventType,
		JobID:         job.JobID,
		JobName:       job.Name,
		UserName:      job.UserName,
		PreviousState: previousState,
		NewState:      job.State,
		NodeList:      job.NodeList,
		Reason:        a.getReasonFromStateChange(previousState, job.State),
	}

	return a.sendEventIfOpen(ctx, eventCh, event)
}

// jobExistsInList checks if a job ID exists in the job list
func (a *JobAdapter) jobExistsInList(jobID int32, jobs []types.Job) bool {
	for _, job := range jobs {
		if job.JobID == jobID {
			return true
		}
	}
	return false
}

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
		a.sendEventIfOpen(ctx, eventCh, types.JobWatchEvent{
			EventTime: time.Now(),
			EventType: "error",
			JobID:     0,
			Reason:    fmt.Sprintf("Failed to poll jobs: %v", err),
		})
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
			if !a.processJobStateChange(ctx, job, previousState, eventCh, opts) {
				return
			}
		}
	}

	// Check for completed/removed jobs (jobs that existed before but don't exist now)
	for jobID, previousState := range jobStates {
		// If job is not found and was not in a terminal state, it might have been removed
		if !a.jobExistsInList(jobID, jobList.Jobs) && !a.isTerminalState(previousState) {
			// Send completion event
			event := types.JobWatchEvent{
				EventTime:     time.Now(),
				EventType:     "completed",
				JobID:         jobID,
				PreviousState: previousState,
				NewState:      types.JobStateCompleted,
				Reason:        "Job completed and removed from active list",
			}

			if !a.sendEventIfOpen(ctx, eventCh, event) {
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

	a.setBasicJobFields(job, apiJob)
	a.setJobResourceFields(job, apiJob)
	a.setJobTimeFields(job, apiJob)
	a.setJobExitCode(job, apiJob)
	a.setJobMemory(job, apiJob)

	return job
}

// setBasicJobFields sets the basic identification and user fields
func (a *JobAdapter) setBasicJobFields(job *types.Job, apiJob api.V0044JobInfo) {
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
}

// setJobResourceFields sets resource-related fields (CPUs, state)
func (a *JobAdapter) setJobResourceFields(job *types.Job, apiJob api.V0044JobInfo) {
	// State - critical for performance metrics (filtering completed jobs)
	if apiJob.JobState != nil && len(*apiJob.JobState) > 0 {
		job.State = types.JobState((*apiJob.JobState)[0])
	}

	// Resource information - CPUs (NoValStruct)
	if apiJob.Cpus != nil && apiJob.Cpus.Set != nil && *apiJob.Cpus.Set && apiJob.Cpus.Number != nil {
		job.CPUs = *apiJob.Cpus.Number
	}
}

// setJobTimeFields sets time-related fields
func (a *JobAdapter) setJobTimeFields(job *types.Job, apiJob api.V0044JobInfo) {
	// Submit time (NoValStruct)
	if apiJob.SubmitTime != nil && apiJob.SubmitTime.Set != nil && *apiJob.SubmitTime.Set && apiJob.SubmitTime.Number != nil {
		job.SubmitTime = time.Unix(*apiJob.SubmitTime.Number, 0)
	}

	// Start time (NoValStruct, validate > 0 to avoid epoch zero)
	if apiJob.StartTime != nil && apiJob.StartTime.Set != nil && *apiJob.StartTime.Set && apiJob.StartTime.Number != nil && *apiJob.StartTime.Number > 0 {
		startTime := time.Unix(*apiJob.StartTime.Number, 0)
		job.StartTime = &startTime
	}

	// End time (NoValStruct, validate > 0 to avoid epoch zero)
	if apiJob.EndTime != nil && apiJob.EndTime.Set != nil && *apiJob.EndTime.Set && apiJob.EndTime.Number != nil && *apiJob.EndTime.Number > 0 {
		endTime := time.Unix(*apiJob.EndTime.Number, 0)
		job.EndTime = &endTime
	}
}

// setJobExitCode sets the exit code from ProcessExitCodeVerbose structure
func (a *JobAdapter) setJobExitCode(job *types.Job, apiJob api.V0044JobInfo) {
	// Exit code - ProcessExitCodeVerbose structure (critical for performance metrics)
	if apiJob.ExitCode != nil && apiJob.ExitCode.ReturnCode != nil &&
		apiJob.ExitCode.ReturnCode.Set != nil && *apiJob.ExitCode.ReturnCode.Set &&
		apiJob.ExitCode.ReturnCode.Number != nil {
		job.ExitCode = *apiJob.ExitCode.ReturnCode.Number
	}
}

// setJobMemory sets memory fields from MemoryPerNode or MemoryPerCpu (MB to bytes)
func (a *JobAdapter) setJobMemory(job *types.Job, apiJob api.V0044JobInfo) {
	// MemoryPerNode (NoValStruct, in MB - convert to bytes)
	if apiJob.MemoryPerNode != nil && apiJob.MemoryPerNode.Set != nil &&
		*apiJob.MemoryPerNode.Set && apiJob.MemoryPerNode.Number != nil {
		job.ResourceRequests.Memory = *apiJob.MemoryPerNode.Number * 1024 * 1024
		return
	}
	// MemoryPerCPU (NoValStruct, in MB - convert to bytes)
	if apiJob.MemoryPerCpu != nil && apiJob.MemoryPerCpu.Set != nil &&
		*apiJob.MemoryPerCpu.Set && apiJob.MemoryPerCpu.Number != nil {
		job.ResourceRequests.MemoryPerCPU = *apiJob.MemoryPerCpu.Number * 1024 * 1024
	}
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
// matchesAccountFilter checks if a job matches the account filter
func (a *JobAdapter) matchesAccountFilter(job types.Job, accounts []string) bool {
	if len(accounts) == 0 {
		return true
	}
	for _, account := range accounts {
		if job.Account == account {
			return true
		}
	}
	return false
}

// matchesUserFilter checks if a job matches the user filter
func (a *JobAdapter) matchesUserFilter(job types.Job, users []string) bool {
	if len(users) == 0 {
		return true
	}
	for _, user := range users {
		if job.UserName == user || strconv.Itoa(int(job.UserID)) == user {
			return true
		}
	}
	return false
}

// matchesStateFilter checks if a job matches the state filter
func (a *JobAdapter) matchesStateFilter(job types.Job, states []types.JobState) bool {
	if len(states) == 0 {
		return true
	}
	for _, state := range states {
		if job.State == state {
			return true
		}
	}
	return false
}

// matchesPartitionFilter checks if a job matches the partition filter
func (a *JobAdapter) matchesPartitionFilter(job types.Job, partitions []string) bool {
	fmt.Printf("[DEBUG partitionFilter] partitions=%v (len=%d, nil=%v), job.Partition=%q\n", partitions, len(partitions), partitions == nil, job.Partition)
	if len(partitions) == 0 {
		fmt.Printf("[DEBUG partitionFilter] Empty partitions filter, returning true\n")
		return true
	}
	for _, partition := range partitions {
		if job.Partition == partition {
			return true
		}
	}
	fmt.Printf("[DEBUG partitionFilter] No match found, returning false\n")
	return false
}

// matchesJobIDFilter checks if a job matches the job ID filter
func (a *JobAdapter) matchesJobIDFilter(job types.Job, jobIDs []int32) bool {
	if len(jobIDs) == 0 {
		return true
	}
	for _, id := range jobIDs {
		if job.JobID == id {
			return true
		}
	}
	return false
}

// matchesJobNameFilter checks if a job matches the job name filter
func (a *JobAdapter) matchesJobNameFilter(job types.Job, jobNames []string) bool {
	if len(jobNames) == 0 {
		return true
	}
	for _, name := range jobNames {
		if job.Name == name {
			return true
		}
	}
	return false
}

func (a *JobAdapter) applyClientSideFilters(jobs []types.Job, opts *types.JobListOptions) []types.Job {
	if opts == nil {
		return jobs
	}

	filtered := make([]types.Job, 0, len(jobs))

	for i, job := range jobs {
		accMatch := a.matchesAccountFilter(job, opts.Accounts)
		userMatch := a.matchesUserFilter(job, opts.Users)
		stateMatch := a.matchesStateFilter(job, opts.States)
		partMatch := a.matchesPartitionFilter(job, opts.Partitions)
		idMatch := a.matchesJobIDFilter(job, opts.JobIDs)
		nameMatch := a.matchesJobNameFilter(job, opts.JobNames)

		if i < 3 {  // Log details for first 3 jobs
			fmt.Printf("[DEBUG JobAdapter] Job %d filter results: Account=%v, User=%v, State=%v, Partition=%v, ID=%v, Name=%v\n",
				i, accMatch, userMatch, stateMatch, partMatch, idMatch, nameMatch)
		}

		if accMatch && userMatch && stateMatch && partMatch && idMatch && nameMatch {
			filtered = append(filtered, job)
		}
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
