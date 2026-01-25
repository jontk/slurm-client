// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"fmt"
	"strconv"
	"time"

	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/jontk/slurm-client/pkg/errors"
)

// JobAdapter implements the JobAdapter interface for v0.0.42
type JobAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewJobAdapter creates a new Job adapter for v0.0.42
func NewJobAdapter(client *api.ClientWithResponses) *JobAdapter {
	return &JobAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Job"),
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
	params := &api.SlurmV0042GetJobsParams{}

	// Apply filters from options
	// Note: v0.0.42 has very limited filtering capabilities in GetJobs
	// Only UpdateTime is supported
	if opts != nil && opts.StartTime != nil {
		updateTimeStr := strconv.FormatInt(opts.StartTime.Unix(), 10)
		params.UpdateTime = &updateTimeStr
	}

	// Set flags to get detailed job information
	flags := api.SlurmV0042GetJobsParamsFlagsDETAIL
	params.Flags = &flags

	// Call the API
	resp, err := a.client.SlurmV0042GetJobsWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list jobs")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}

	// Check for API response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from API")
	}

	// Convert the response to common types
	jobs := make([]types.Job, 0)

	if len(resp.JSON200.Jobs) > 0 {
		for _, apiJob := range resp.JSON200.Jobs {
			job, err := a.convertAPIJobToCommon(apiJob)
			if err != nil {
				// Log conversion error but continue
				continue
			}
			jobs = append(jobs, *job)
		}
	}

	// Apply client-side filtering since API has limited support
	if opts != nil {
		jobs = a.applyClientSideFilters(jobs, opts)
	}

	// Apply pagination
	start := 0
	if opts != nil && opts.Offset > 0 {
		start = opts.Offset
	}
	if start >= len(jobs) {
		return &types.JobList{
			Jobs:  []types.Job{},
			Total: len(jobs),
		}, nil
	}

	end := len(jobs)
	if opts != nil && opts.Limit > 0 {
		end = start + opts.Limit
		if end > len(jobs) {
			end = len(jobs)
		}
	}

	return &types.JobList{
		Jobs:  jobs[start:end],
		Total: len(jobs),
	}, nil
}

// Get retrieves a specific job by ID
func (a *JobAdapter) Get(ctx context.Context, jobID int32) (*types.Job, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters
	params := &api.SlurmV0042GetJobParams{}
	flags := api.SlurmV0042GetJobParamsFlagsDETAIL
	params.Flags = &flags

	// Call the API
	resp, err := a.client.SlurmV0042GetJobWithResponse(ctx, strconv.FormatInt(int64(jobID), 10), params)
	if err != nil {
		return nil, a.WrapError(err, fmt.Sprintf("failed to get job %d", jobID))
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}

	// Check for API response
	if resp.JSON200 == nil || resp.JSON200.Jobs == nil || len(resp.JSON200.Jobs) == 0 {
		return nil, fmt.Errorf("job %d not found", jobID)
	}

	// Convert the first job in the response
	jobs := resp.JSON200.Jobs
	return a.convertAPIJobToCommon(jobs[0])
}

// Cancel cancels a job
func (a *JobAdapter) Cancel(ctx context.Context, jobID int32, opts *types.JobCancelRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Prepare parameters
	params := &api.SlurmV0042DeleteJobParams{}

	// Set signal from options
	signal := "SIGTERM" // Default signal
	if opts != nil && opts.Signal != "" {
		signal = opts.Signal
	}
	params.Signal = &signal

	// TODO: Add flag support if JobCancelRequest gains a Flags field

	// Call the API
	resp, err := a.client.SlurmV0042DeleteJobWithResponse(ctx, strconv.FormatInt(int64(jobID), 10), params)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to cancel job %d", jobID))
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}

	return nil
}

// Submit submits a new job
func (a *JobAdapter) Submit(ctx context.Context, job *types.JobCreate) (*types.JobSubmitResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Create the job submission structure
	apiJobSubmission := &api.V0042JobDescMsg{
		Name:      &job.Name,
		Account:   &job.Account,
		Partition: &job.Partition,
	}

	// Handle script/command
	if job.Script != "" {
		apiJobSubmission.Script = &job.Script
	}

	// Handle working directory
	if job.WorkingDirectory != "" {
		apiJobSubmission.CurrentWorkingDirectory = &job.WorkingDirectory
	}

	// Handle standard output/error/input
	if job.StandardOutput != "" {
		apiJobSubmission.StandardOutput = &job.StandardOutput
	}
	if job.StandardError != "" {
		apiJobSubmission.StandardError = &job.StandardError
	}
	if job.StandardInput != "" {
		apiJobSubmission.StandardInput = &job.StandardInput
	}

	// Handle time limit
	if job.TimeLimit > 0 {
		timeLimitNumber := job.TimeLimit
		apiJobSubmission.TimeLimit = &api.V0042Uint32NoValStruct{
			Number: &timeLimitNumber,
			Set:    &[]bool{true}[0],
		}
	}

	// Handle node count
	if job.Nodes > 0 {
		nodesStr := strconv.Itoa(int(job.Nodes))
		apiJobSubmission.Nodes = &nodesStr
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

	// Set environment in job submission
	apiJobSubmission.Environment = &envVars

	// Create request body
	apiJobReq := api.V0042JobSubmitReq{
		Jobs: &[]api.V0042JobDescMsg{*apiJobSubmission},
	}

	// Call the API
	resp, err := a.client.SlurmV0042PostJobSubmitWithResponse(ctx, apiJobReq)
	if err != nil {
		return nil, a.WrapError(err, "failed to submit job")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}

	// Check for API response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from API")
	}

	// Convert response
	return a.convertAPIJobSubmitResponseToCommon(resp.JSON200)
}

// Update updates an existing job
func (a *JobAdapter) Update(ctx context.Context, jobID int32, updates *types.JobUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// v0.0.42 doesn't have a job update endpoint
	return a.HandleNotImplemented("Job update", "v0.0.42")
}

// Signal sends a signal to a job
func (a *JobAdapter) Signal(ctx context.Context, req *types.JobSignalRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// This is a placeholder - v0.0.42 doesn't have a dedicated signal endpoint
	// Signaling is typically done through the delete/cancel endpoint with different signals
	return fmt.Errorf("job signaling not directly supported via v0.0.42 API - use cancel with specific signal")
}

// Hold holds a job
func (a *JobAdapter) Hold(ctx context.Context, req *types.JobHoldRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// v0.0.42 doesn't have a dedicated hold endpoint
	// Job holds are typically managed through job updates or administrative commands
	return fmt.Errorf("job hold not directly supported via v0.0.42 API - use administrative commands")
}

// Notify sends a notification for a job
func (a *JobAdapter) Notify(ctx context.Context, req *types.JobNotifyRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// This is a placeholder - v0.0.42 doesn't have a job notification endpoint
	return fmt.Errorf("job notification not supported via v0.0.42 API")
}

// Requeue requeues a job
func (a *JobAdapter) Requeue(ctx context.Context, jobID int32) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Prepare parameters with FEDERATIONREQUEUE flag
	params := &api.SlurmV0042DeleteJobParams{}
	requeueFlag := api.FEDERATIONREQUEUE
	params.Flags = &requeueFlag

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0042DeleteJobWithResponse(ctx, strconv.Itoa(int(jobID)), params)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}

	return nil
}

// Watch watches for job changes
func (a *JobAdapter) Watch(ctx context.Context, opts *types.JobWatchOptions) (<-chan types.JobWatchEvent, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// This is a placeholder - actual implementation would use the underlying API's watch mechanism
	// For now, return an error indicating it's not implemented
	return nil, fmt.Errorf("watch functionality not implemented for v0.0.42")
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
	resp, err := a.client.SlurmV0042PostJobAllocateWithResponse(ctx, *apiReq)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
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

	// Account is required for SLURM v0.0.42 job allocation
	if req.Account == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "account is required for job allocation in SLURM v0.0.42", "account", req.Account, nil)
	}

	// At least one resource requirement should be specified
	if req.Nodes == "" && req.CPUs == 0 && req.Memory == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "at least one resource requirement (nodes, cpus, or memory) must be specified", "resources", req, nil)
	}

	return nil
}

// convertCommonJobAllocateToAPI converts common allocation request to API request
func (a *JobAdapter) convertCommonJobAllocateToAPI(req *types.JobAllocateRequest) *api.V0042JobAllocReq {
	// Create the job description message
	jobDesc := &api.V0042JobDescMsg{}

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
	apiReq := &api.V0042JobAllocReq{
		Job: jobDesc,
	}

	return apiReq
}

// convertAPIJobAllocateResponseToCommon converts API allocation response to common response
func (a *JobAdapter) convertAPIJobAllocateResponseToCommon(apiResp *api.V0042OpenapiJobAllocResp) *types.JobAllocateResponse {
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

	// Extract metadata if available
	if apiResp != nil && apiResp.Meta != nil {
		resp.Meta["api_version"] = "v0.0.42"
	}

	return resp
}

// extractTimeIfSet safely extracts a Unix timestamp from a nested structure
func (a *JobAdapter) extractTimeIfSet(timeStruct *api.V0042Uint64NoValStruct) *time.Time {
	if timeStruct == nil || timeStruct.Number == nil {
		return nil
	}
	t := time.Unix(*timeStruct.Number, 0)
	return &t
}

// extractTimeRequired safely extracts a required Unix timestamp (non-pointer return)
func (a *JobAdapter) extractTimeRequired(timeStruct *api.V0042Uint64NoValStruct) time.Time {
	if timeStruct == nil || timeStruct.Number == nil {
		return time.Time{}
	}
	return time.Unix(*timeStruct.Number, 0)
}

// extractInt32IfSet safely extracts int32 from nested structure
func (a *JobAdapter) extractInt32IfSet(valStruct *api.V0042Uint32NoValStruct) int32 {
	if valStruct == nil || valStruct.Number == nil {
		return 0
	}
	return *valStruct.Number
}

func (a *JobAdapter) convertAPIJobToCommon(apiJob api.V0042JobInfo) (*types.Job, error) {
	job := &types.Job{}

	// Set basic fields
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

	if apiJob.UserName != nil {
		job.UserName = *apiJob.UserName
	}

	// Job state conversion
	if apiJob.JobState != nil && len(*apiJob.JobState) > 0 {
		job.State = types.JobState((*apiJob.JobState)[0]) // Take the first state and convert to JobState type
	}

	// Resource information - CPU count is directly available in JobResources
	if apiJob.JobResources != nil {
		job.CPUs = apiJob.JobResources.Cpus
	}

	if apiJob.Nodes != nil {
		job.NodeList = *apiJob.Nodes
	}

	// Time fields - convert from Unix timestamp to time.Time
	job.SubmitTime = a.extractTimeRequired(apiJob.SubmitTime)
	job.StartTime = a.extractTimeIfSet(apiJob.StartTime)
	job.EndTime = a.extractTimeIfSet(apiJob.EndTime)
	job.TimeLimit = a.extractInt32IfSet(apiJob.TimeLimit)

	// Working directory is not available in v0.0.42 JobInfo structure

	// Command
	if apiJob.Command != nil {
		job.Command = *apiJob.Command
	}

	return job, nil
}

// convertAPIJobSubmitResponseToCommon converts API job submit response to common type
func (a *JobAdapter) convertAPIJobSubmitResponseToCommon(apiResp *api.V0042OpenapiJobSubmitResponse) (*types.JobSubmitResponse, error) {
	resp := &types.JobSubmitResponse{}

	// Extract job ID from the first job result if available
	if apiResp.JobId != nil {
		resp.JobID = *apiResp.JobId
	}

	// Extract user message if available
	if apiResp.JobSubmitUserMsg != nil {
		resp.JobSubmitUserMsg = *apiResp.JobSubmitUserMsg
	}

	return resp, nil
}

// applyClientSideFilters applies filtering to the job list
func (a *JobAdapter) applyClientSideFilters(jobs []types.Job, opts *types.JobListOptions) []types.Job {
	filteredJobs := make([]types.Job, 0)
	for _, job := range jobs {
		if a.matchesJobFilters(job, opts) {
			filteredJobs = append(filteredJobs, job)
		}
	}
	return filteredJobs
}

// matchesJobFilters checks if a job matches all the provided filters
func (a *JobAdapter) matchesJobFilters(job types.Job, opts *types.JobListOptions) bool {
	return a.checkAccountFilter(job, opts.Accounts) &&
		a.checkPartitionFilter(job, opts.Partitions) &&
		a.checkUserFilter(job, opts.Users) &&
		a.checkStateFilter(job, opts.States) &&
		a.checkJobIDFilter(job, opts.JobIDs)
}

func (a *JobAdapter) checkAccountFilter(job types.Job, accounts []string) bool {
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

func (a *JobAdapter) checkPartitionFilter(job types.Job, partitions []string) bool {
	if len(partitions) == 0 {
		return true
	}
	for _, partition := range partitions {
		if job.Partition == partition {
			return true
		}
	}
	return false
}

func (a *JobAdapter) checkUserFilter(job types.Job, users []string) bool {
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

func (a *JobAdapter) checkStateFilter(job types.Job, states []types.JobState) bool {
	if len(states) == 0 {
		return true
	}
	for _, state := range states {
		if string(job.State) == string(state) {
			return true
		}
	}
	return false
}

func (a *JobAdapter) checkJobIDFilter(job types.Job, jobIDs []int32) bool {
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
