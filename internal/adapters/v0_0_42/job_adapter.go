package v0_0_42

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
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
	if opts != nil {
		if len(opts.Accounts) > 0 {
			accountStr := strings.Join(opts.Accounts, ",")
			params.Account = &accountStr
		}
		if len(opts.Partitions) > 0 {
			partitionStr := strings.Join(opts.Partitions, ",")
			params.Partition = &partitionStr
		}
		if len(opts.Users) > 0 {
			usersStr := strings.Join(opts.Users, ",")
			params.Users = &usersStr
		}
		if len(opts.JobIDs) > 0 {
			// Convert job IDs to string
			jobIDStrs := make([]string, len(opts.JobIDs))
			for i, id := range opts.JobIDs {
				jobIDStrs[i] = strconv.FormatUint(uint64(id), 10)
			}
			jobsStr := strings.Join(jobIDStrs, ",")
			params.Job = &jobsStr
		}
		if opts.StartTime != nil {
			startTimeStr := strconv.FormatInt(opts.StartTime.Unix(), 10)
			params.StartTime = &startTimeStr
		}
		if opts.EndTime != nil {
			endTimeStr := strconv.FormatInt(opts.EndTime.Unix(), 10)
			params.EndTime = &endTimeStr
		}
		if len(opts.States) > 0 {
			// Convert states to string
			stateStr := strings.Join(opts.States, ",")
			params.State = &stateStr
		}
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
		return nil, a.HandleAPIError(resp.StatusCode(), resp.Body)
	}

	// Check for API response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from API")
	}

	// Convert the response to common types
	jobList := &types.JobList{
		Jobs: make([]*types.Job, 0),
	}

	if resp.JSON200.Jobs != nil {
		for _, apiJob := range *resp.JSON200.Jobs {
			job, err := a.convertAPIJobToCommon(apiJob)
			if err != nil {
				// Log conversion error but continue
				continue
			}
			jobList.Jobs = append(jobList.Jobs, job)
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

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters
	params := &api.SlurmV0042GetJobParams{}
	flags := api.SlurmV0042GetJobParamsFlagsDETAIL
	params.Flags = &flags

	// Call the API
	resp, err := a.client.SlurmV0042GetJobWithResponse(ctx, strconv.FormatUint(uint64(jobID), 10), params)
	if err != nil {
		return nil, a.WrapError(err, fmt.Sprintf("failed to get job %d", jobID))
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(resp.StatusCode(), resp.Body)
	}

	// Check for API response
	if resp.JSON200 == nil || resp.JSON200.Jobs == nil || len(*resp.JSON200.Jobs) == 0 {
		return nil, fmt.Errorf("job %d not found", jobID)
	}

	// Convert the first job in the response
	jobs := *resp.JSON200.Jobs
	return a.convertAPIJobToCommon(jobs[0])
}

// Cancel cancels a job
func (a *JobAdapter) Cancel(ctx context.Context, jobID uint32, signal string, flags *types.JobCancelFlags) error {
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
	
	// Set signal
	if signal != "" {
		params.Signal = &signal
	}

	// Apply flags if provided
	if flags != nil {
		var apiFlags []api.SlurmV0042DeleteJobParamsFlags
		if flags.ArrayTask {
			apiFlags = append(apiFlags, api.ARRAYTASK)
		}
		if flags.BatchJob {
			apiFlags = append(apiFlags, api.BATCHJOB)
		}
		if flags.FullJob {
			apiFlags = append(apiFlags, api.FULLJOB)
		}
		if flags.Hurry {
			apiFlags = append(apiFlags, api.HURRY)
		}
		if len(apiFlags) > 0 {
			params.Flags = &apiFlags
		}
	}

	// Call the API
	resp, err := a.client.SlurmV0042DeleteJobWithResponse(ctx, strconv.FormatUint(uint64(jobID), 10), params)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to cancel job %d", jobID))
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleAPIError(resp.StatusCode(), resp.Body)
	}

	return nil
}

// Submit submits a new job
func (a *JobAdapter) Submit(ctx context.Context, jobSpec *types.JobSubmitRequest) (*types.JobSubmitResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert common job spec to API format
	apiJobReq, err := a.convertCommonJobSubmitToAPI(jobSpec)
	if err != nil {
		return nil, a.WrapError(err, "failed to convert job submission request")
	}

	// Call the API
	resp, err := a.client.SlurmV0042PostJobSubmitWithResponse(ctx, apiJobReq)
	if err != nil {
		return nil, a.WrapError(err, "failed to submit job")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(resp.StatusCode(), resp.Body)
	}

	// Check for API response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from API")
	}

	// Convert response
	return a.convertAPIJobSubmitResponseToCommon(resp.JSON200)
}

// Update updates an existing job
func (a *JobAdapter) Update(ctx context.Context, jobID uint32, updates *types.JobUpdateRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Convert common update request to API format
	apiJobUpdate, err := a.convertCommonJobUpdateToAPI(updates)
	if err != nil {
		return a.WrapError(err, "failed to convert job update request")
	}

	// Call the API
	resp, err := a.client.SlurmV0042PostJobUpdateWithResponse(ctx, strconv.FormatUint(uint64(jobID), 10), apiJobUpdate)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to update job %d", jobID))
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleAPIError(resp.StatusCode(), resp.Body)
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