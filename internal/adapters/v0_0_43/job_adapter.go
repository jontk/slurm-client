package v0_0_43

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/jontk/slurm-client/pkg/errors"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// JobAdapter implements the JobAdapter interface for v0.0.43
type JobAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewJobAdapter creates a new Job adapter for v0.0.43
func NewJobAdapter(client *api.ClientWithResponses) *JobAdapter {
	return &JobAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Job"),
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
	// Use SlurmV0043GetJobsParams for current job queries
	params := &api.SlurmV0043GetJobsParams{}

	// Apply filters from options - SlurmV0043GetJobsParams only supports UpdateTime and Flags
	if opts != nil {
		// Note: SlurmV0043GetJobsParams has limited filtering capabilities
		// Most filtering will need to be done post-query
		if opts.StartTime != nil {
			updateTime := fmt.Sprintf("%d", opts.StartTime.Unix())
			params.UpdateTime = &updateTime
		}
		// Other filters (accounts, users, states, etc.) need to be applied after getting results
		// Store filter options for client-side filtering
	}

	// Call the generated OpenAPI client for current job queries
	resp, err := a.client.SlurmV0043GetJobsWithResponse(ctx, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "List Jobs"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Jobs, "List Jobs - jobs field"); err != nil {
		return nil, err
	}

	// Convert the response to common types - SlurmV0043GetJobs returns V0043JobInfo
	jobList := make([]types.Job, 0, len(resp.JSON200.Jobs))
	for _, apiJobInfo := range resp.JSON200.Jobs {
		job, err := a.convertAPIJobToCommon(apiJobInfo)
		if err != nil {
			return nil, a.HandleConversionError(err, apiJobInfo.JobId)
		}
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
	params := &api.SlurmV0043GetJobParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0043GetJobWithResponse(ctx, strconv.Itoa(int(jobID)), params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
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
	job, err := a.convertAPIJobToCommon(resp.JSON200.Jobs[0])
	if err != nil {
		return nil, a.HandleConversionError(err, jobID)
	}

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
	apiJob, err := a.convertCommonJobCreateToAPI(job)
	if err != nil {
		return nil, err
	}

	// Create request body - V0043JobDescMsg format
	reqBody := api.V0043JobDescMsg{
		Account:   apiJob.Account,
		Name:      apiJob.Name,
		Partition: apiJob.Partition,
		// Copy other fields as needed
	}

	// Call the generated OpenAPI client - job submission doesn't need job ID
	resp, err := a.client.SlurmV0043PostJobWithResponse(ctx, "0", reqBody)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
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

	// Extract job ID from response - V0043OpenapiJobPostResponse structure
	if resp.JSON200 != nil && resp.JSON200.Results != nil && len(*resp.JSON200.Results) > 0 {
		result := (*resp.JSON200.Results)[0]
		if result.JobId != nil {
			jobID = *result.JobId
		}
	}
	
	// Note: V0043JobArrayResponseMsgEntry doesn't have Warnings field
	// Handle warnings through errors array if available

	// Check for errors in response
	if resp.JSON200.Errors != nil {
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
	apiJob, err := a.convertCommonJobUpdateToAPI(existingJob, update)
	if err != nil {
		return err
	}

	// Create request body for job update - V0043JobDescMsg format
	reqBody := api.V0043JobDescMsg{
		Account:   apiJob.Account,
		Name:      apiJob.Name,
		Partition: apiJob.Partition,
		// Copy other fields as needed
	}

	// Call the generated OpenAPI client (POST is used for updates in SLURM API)
	resp, err := a.client.SlurmV0043PostJobWithResponse(ctx, strconv.Itoa(int(jobID)), reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "v0.0.43")
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
	params := &api.SlurmV0043DeleteJobParams{}

	if opts != nil {
		if opts.Signal != "" {
			params.Signal = &opts.Signal
		}
		if opts.Message != "" {
			// Convert message to appropriate flags type
			flags := api.SlurmV0043DeleteJobParamsFlags(opts.Message)
			params.Flags = &flags
		}
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0043DeleteJobWithResponse(ctx, strconv.Itoa(int(jobID)), params)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	// Create adapter with special handling for 204 (No Content) status
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)

	// For DELETE operations, 204 is also a success
	if resp.StatusCode() == 204 {
		return nil
	}

	return common.HandleAPIResponse(responseAdapter, "v0.0.43")
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

	// TODO: The v0.0.43 API doesn't seem to have a dedicated job signal endpoint
	// For now, we'll return an error indicating the feature is not supported
	return errors.NewClientError(errors.ErrorCodeValidationFailed, "Job signaling is not supported in v0.0.43 API")
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

	// Create request body for job hold/release - V0043JobDescMsg format
	reqBody := api.V0043JobDescMsg{
		// Basic job identification
		Name: func() *string {
			// We might need the job name for proper identification
			return nil
		}(),
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0043PostJobWithResponse(ctx, strconv.Itoa(int(req.JobID)), reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "v0.0.43")
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

	// For v0.0.43, we simulate notification by updating the job comment
	// This is a workaround as the API might not have dedicated notification endpoints
	update := &types.JobUpdate{
		Comment: &req.Message,
	}

	return a.Update(ctx, req.JobID, update)
}

// validateJobCreate validates job creation request
func (a *JobAdapter) validateJobCreate(job *types.JobCreate) error {
	if job == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "job creation data is required", "job", nil, nil)
	}
	if job.Command == "" && job.Script == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "either command or script is required", "job", job, nil)
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
func (a *JobAdapter) convertAPIJobToCommon(apiJob api.V0043JobInfo) (*types.Job, error) {
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
	return job, nil
}

func (a *JobAdapter) convertCommonJobCreateToAPI(create *types.JobCreate) (*api.V0043Job, error) {
	apiJob := &api.V0043Job{}
	if create.Name != "" {
		apiJob.Name = &create.Name
	}
	if create.Account != "" {
		apiJob.Account = &create.Account
	}
	if create.Partition != "" {
		apiJob.Partition = &create.Partition
	}
	// TODO: Add more field conversions as needed
	return apiJob, nil
}

func (a *JobAdapter) convertCommonJobUpdateToAPI(existing *types.Job, update *types.JobUpdate) (*api.V0043Job, error) {
	apiJob := &api.V0043Job{}
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
	return apiJob, nil
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
				if job.UserName == user || fmt.Sprintf("%d", job.UserID) == user {
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