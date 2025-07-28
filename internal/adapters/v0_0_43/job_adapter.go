package v0_0_43

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
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
	params := &api.SlurmV0043GetJobsParams{}

	// Apply filters from options
	if opts != nil {
		if len(opts.Accounts) > 0 {
			accountStr := strings.Join(opts.Accounts, ",")
			params.Account = &accountStr
		}
		if len(opts.Users) > 0 {
			userStr := strings.Join(opts.Users, ",")
			params.Users = &userStr
		}
		if len(opts.States) > 0 {
			stateStrs := make([]string, len(opts.States))
			for i, state := range opts.States {
				stateStrs[i] = string(state)
			}
			stateStr := strings.Join(stateStrs, ",")
			params.States = &stateStr
		}
		if len(opts.Partitions) > 0 {
			partitionStr := strings.Join(opts.Partitions, ",")
			params.Partition = &partitionStr
		}
		if len(opts.JobIDs) > 0 {
			jobIDStrs := make([]string, len(opts.JobIDs))
			for i, id := range opts.JobIDs {
				jobIDStrs[i] = fmt.Sprintf("%d", id)
			}
			jobIDStr := strings.Join(jobIDStrs, ",")
			params.JobIds = &jobIDStr
		}
		if len(opts.JobNames) > 0 {
			nameStr := strings.Join(opts.JobNames, ",")
			params.Name = &nameStr
		}
		if opts.StartTime != nil {
			startTime := fmt.Sprintf("%d", opts.StartTime.Unix())
			params.StartTime = &startTime
		}
		if opts.EndTime != nil {
			endTime := fmt.Sprintf("%d", opts.EndTime.Unix())
			params.EndTime = &endTime
		}
	}

	// Call the generated OpenAPI client
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

	// Convert the response to common types
	jobList := make([]types.Job, 0, len(resp.JSON200.Jobs))
	for _, apiJob := range resp.JSON200.Jobs {
		job, err := a.convertAPIJobToCommon(apiJob)
		if err != nil {
			return nil, a.HandleConversionError(err, apiJob.JobId)
		}
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
	params := &api.SlurmV0043GetJobParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0043GetJobWithResponse(ctx, jobID, params)
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
		return nil, common.NewResourceNotFoundError("Job", jobID)
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

	// Create request body
	reqBody := api.SlurmV0043PostJobJSONRequestBody{
		Jobs: []api.V0043JobProperties{*apiJob},
	}

	// Prepare parameters for the API call
	params := &api.SlurmV0043PostJobParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0043PostJobWithResponse(ctx, params, reqBody)
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

	if resp.JSON200.Results != nil && len(*resp.JSON200.Results) > 0 {
		result := (*resp.JSON200.Results)[0]
		if result.JobId != nil {
			jobID = *result.JobId
		}
		if result.Warnings != nil {
			for _, warning := range *result.Warnings {
				if warning.Warning != nil {
					warnings = append(warnings, *warning.Warning)
				}
			}
		}
	}

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

	// Create request body
	reqBody := api.SlurmV0043PostJobJSONRequestBody{
		Jobs: []api.V0043JobProperties{*apiJob},
	}

	// Prepare parameters for the API call
	params := &api.SlurmV0043PostJobParams{}

	// Call the generated OpenAPI client (POST is used for updates in SLURM API)
	resp, err := a.client.SlurmV0043PostJobWithResponse(ctx, params, reqBody)
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
			params.Flags = &opts.Message
		}
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0043DeleteJobWithResponse(ctx, jobID, params)
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

	// Create request body
	reqBody := api.SlurmV0043PostJobSignalJSONRequestBody{
		Signal: req.Signal,
	}

	// Prepare parameters for the API call
	params := &api.SlurmV0043PostJobSignalParams{}
	if req.StepID != "" {
		params.StepId = &req.StepID
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0043PostJobSignalWithResponse(ctx, req.JobID, params, reqBody)
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

	// Create request body
	var holdState string
	if req.Hold {
		holdState = "USER_HOLD"
	} else {
		holdState = "RELEASE"
	}

	reqBody := api.SlurmV0043PostJobJSONRequestBody{
		Jobs: []api.V0043JobProperties{
			{
				JobId: &req.JobID,
				JobState: &[]api.V0043JobState{
					api.V0043JobState(holdState),
				},
				Priority: func() *api.V0043Uint32NoValStruct {
					if req.Priority != 0 {
						setTrue := true
						priority := int32(req.Priority)
						return &api.V0043Uint32NoValStruct{
							Set:    &setTrue,
							Number: &priority,
						}
					}
					return nil
				}(),
			},
		},
	}

	// Prepare parameters for the API call
	params := &api.SlurmV0043PostJobParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0043PostJobWithResponse(ctx, params, reqBody)
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
		return common.NewValidationError("job creation data is required", "job", nil)
	}
	if job.Command == "" && job.Script == "" {
		return common.NewValidationError("either command or script is required", "job", job)
	}
	return nil
}

// validateJobUpdate validates job update request
func (a *JobAdapter) validateJobUpdate(update *types.JobUpdate) error {
	if update == nil {
		return common.NewValidationError("job update data is required", "update", nil)
	}
	// At least one field should be provided for update
	if update.Name == nil && update.Account == nil && update.Partition == nil && 
	   update.QoS == nil && update.TimeLimit == nil && update.Priority == nil && 
	   update.Nice == nil && update.Comment == nil {
		return common.NewValidationError("at least one field must be provided for update", "update", update)
	}
	return nil
}

// validateJobSignalRequest validates job signal request
func (a *JobAdapter) validateJobSignalRequest(req *types.JobSignalRequest) error {
	if req == nil {
		return common.NewValidationError("job signal request is required", "req", nil)
	}
	if req.JobID == 0 {
		return common.NewValidationError("job ID is required", "jobID", req.JobID)
	}
	if req.Signal == "" {
		return common.NewValidationError("signal is required", "signal", req.Signal)
	}
	return nil
}

// validateJobHoldRequest validates job hold request
func (a *JobAdapter) validateJobHoldRequest(req *types.JobHoldRequest) error {
	if req == nil {
		return common.NewValidationError("job hold request is required", "req", nil)
	}
	if req.JobID == 0 {
		return common.NewValidationError("job ID is required", "jobID", req.JobID)
	}
	return nil
}

// validateJobNotifyRequest validates job notify request
func (a *JobAdapter) validateJobNotifyRequest(req *types.JobNotifyRequest) error {
	if req == nil {
		return common.NewValidationError("job notify request is required", "req", nil)
	}
	if req.JobID == 0 {
		return common.NewValidationError("job ID is required", "jobID", req.JobID)
	}
	if req.Message == "" {
		return common.NewValidationError("message is required", "message", req.Message)
	}
	return nil
}