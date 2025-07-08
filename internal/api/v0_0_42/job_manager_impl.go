package v0_0_42

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/interfaces"
)

// JobManagerImpl provides the actual implementation for JobManager methods
type JobManagerImpl struct {
	client *WrapperClient
}

// NewJobManagerImpl creates a new JobManager implementation
func NewJobManagerImpl(client *WrapperClient) *JobManagerImpl {
	return &JobManagerImpl{client: client}
}

// List jobs with optional filtering
func (m *JobManagerImpl) List(ctx context.Context, opts *interfaces.ListJobsOptions) (*interfaces.JobList, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, fmt.Errorf("API client not initialized")
	}
	
	// Prepare parameters for the API call
	params := &SlurmV0042GetJobsParams{}
	
	// Set flags to get detailed job information
	flags := SlurmV0042GetJobsParamsFlagsDETAIL
	params.Flags = &flags
	
	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0042GetJobsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}
	
	// Check HTTP status
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode(), resp.Status())
	}
	
	// Check for API errors
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response format")
	}
	
	if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		return nil, fmt.Errorf("API error: %v", (*resp.JSON200.Errors)[0])
	}
	
	// Convert the response to our interface types
	jobs := make([]interfaces.Job, 0, len(resp.JSON200.Jobs))
	for _, apiJob := range resp.JSON200.Jobs {
		job, err := convertAPIJobToInterface(apiJob)
		if err != nil {
			return nil, fmt.Errorf("failed to convert job data: %w", err)
		}
		jobs = append(jobs, *job)
	}
	
	// Apply client-side filtering if options are provided
	if opts != nil {
		jobs = filterJobs(jobs, opts)
	}
	
	return &interfaces.JobList{
		Jobs:  jobs,
		Total: len(jobs),
	}, nil
}

// convertAPIJobToInterface converts a V0042JobInfo to interfaces.Job
func convertAPIJobToInterface(apiJob V0042JobInfo) (*interfaces.Job, error) {
	job := &interfaces.Job{}
	
	// Job ID - simple int32 pointer
	if apiJob.JobId != nil {
		job.ID = strconv.FormatInt(int64(*apiJob.JobId), 10)
	}
	
	// Job name
	if apiJob.Name != nil {
		job.Name = *apiJob.Name
	}
	
	// User ID - simple int32 pointer
	if apiJob.UserId != nil {
		job.UserID = strconv.FormatInt(int64(*apiJob.UserId), 10)
	}
	
	// Group ID - simple int32 pointer
	if apiJob.GroupId != nil {
		job.GroupID = strconv.FormatInt(int64(*apiJob.GroupId), 10)
	}
	
	// Job state - array of strings
	if apiJob.JobState != nil && len(*apiJob.JobState) > 0 {
		job.State = (*apiJob.JobState)[0]
	}
	
	// Partition
	if apiJob.Partition != nil {
		job.Partition = *apiJob.Partition
	}
	
	// Priority - NoValStruct
	if apiJob.Priority != nil && apiJob.Priority.Set != nil && *apiJob.Priority.Set && apiJob.Priority.Number != nil {
		job.Priority = int(*apiJob.Priority.Number)
	}
	
	// Submit time - NoValStruct
	if apiJob.SubmitTime != nil && apiJob.SubmitTime.Set != nil && *apiJob.SubmitTime.Set && apiJob.SubmitTime.Number != nil {
		job.SubmitTime = time.Unix(*apiJob.SubmitTime.Number, 0)
	}
	
	// Start time - NoValStruct
	if apiJob.StartTime != nil && apiJob.StartTime.Set != nil && *apiJob.StartTime.Set && apiJob.StartTime.Number != nil && *apiJob.StartTime.Number > 0 {
		startTime := time.Unix(*apiJob.StartTime.Number, 0)
		job.StartTime = &startTime
	}
	
	// End time - NoValStruct
	if apiJob.EndTime != nil && apiJob.EndTime.Set != nil && *apiJob.EndTime.Set && apiJob.EndTime.Number != nil && *apiJob.EndTime.Number > 0 {
		endTime := time.Unix(*apiJob.EndTime.Number, 0)
		job.EndTime = &endTime
	}
	
	// CPUs - NoValStruct
	if apiJob.Cpus != nil && apiJob.Cpus.Set != nil && *apiJob.Cpus.Set && apiJob.Cpus.Number != nil {
		job.CPUs = int(*apiJob.Cpus.Number)
	}
	
	// Memory (convert MB to bytes for consistency) - NoValStruct
	if apiJob.MemoryPerNode != nil && apiJob.MemoryPerNode.Set != nil && *apiJob.MemoryPerNode.Set && apiJob.MemoryPerNode.Number != nil {
		job.Memory = int(*apiJob.MemoryPerNode.Number) * 1024 * 1024
	}
	
	// Time limit (in minutes) - NoValStruct
	if apiJob.TimeLimit != nil && apiJob.TimeLimit.Set != nil && *apiJob.TimeLimit.Set && apiJob.TimeLimit.Number != nil {
		job.TimeLimit = int(*apiJob.TimeLimit.Number)
	}
	
	// Working directory
	if apiJob.CurrentWorkingDirectory != nil {
		job.WorkingDir = *apiJob.CurrentWorkingDirectory
	}
	
	// Command
	if apiJob.Command != nil {
		job.Command = *apiJob.Command
	}
	
	// Environment variables - Initialize empty map since not directly available in JobInfo
	job.Environment = make(map[string]string)
	
	// Nodes - Extract from JobResources
	if apiJob.JobResources != nil && apiJob.JobResources.Nodes != nil && apiJob.JobResources.Nodes.List != nil {
		// Parse node list string into slice
		nodeListStr := *apiJob.JobResources.Nodes.List
		if nodeListStr != "" {
			// Simple splitting by comma - real implementation might need more sophisticated parsing
			job.Nodes = strings.Split(nodeListStr, ",")
		}
	}
	
	// Exit code - ProcessExitCodeVerbose
	if apiJob.ExitCode != nil && apiJob.ExitCode.ReturnCode != nil && 
		apiJob.ExitCode.ReturnCode.Set != nil && *apiJob.ExitCode.ReturnCode.Set && 
		apiJob.ExitCode.ReturnCode.Number != nil {
		job.ExitCode = int(*apiJob.ExitCode.ReturnCode.Number)
	}
	
	// Initialize metadata
	job.Metadata = make(map[string]interface{})
	
	// Add additional metadata from API response
	if apiJob.Account != nil {
		job.Metadata["account"] = *apiJob.Account
	}
	if apiJob.AdminComment != nil {
		job.Metadata["admin_comment"] = *apiJob.AdminComment
	}
	if apiJob.AllocatingNode != nil {
		job.Metadata["allocating_node"] = *apiJob.AllocatingNode
	}
	
	return job, nil
}

// filterJobs applies client-side filtering to job list
func filterJobs(jobs []interfaces.Job, opts *interfaces.ListJobsOptions) []interfaces.Job {
	var filtered []interfaces.Job
	
	for _, job := range jobs {
		// Filter by user ID
		if opts.UserID != "" && job.UserID != opts.UserID {
			continue
		}
		
		// Filter by states
		if len(opts.States) > 0 {
			stateMatch := false
			for _, state := range opts.States {
				if strings.EqualFold(job.State, state) {
					stateMatch = true
					break
				}
			}
			if !stateMatch {
				continue
			}
		}
		
		// Filter by partition
		if opts.Partition != "" && !strings.EqualFold(job.Partition, opts.Partition) {
			continue
		}
		
		filtered = append(filtered, job)
	}
	
	// Apply limit and offset
	if opts.Offset > 0 {
		if opts.Offset >= len(filtered) {
			return []interfaces.Job{}
		}
		filtered = filtered[opts.Offset:]
	}
	
	if opts.Limit > 0 && len(filtered) > opts.Limit {
		filtered = filtered[:opts.Limit]
	}
	
	return filtered
}