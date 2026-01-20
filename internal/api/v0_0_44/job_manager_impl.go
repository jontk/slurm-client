// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_44

import (
	"context"
	"fmt"
	"time"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
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
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmV0044GetJobsParams{}

	// Set flags to get detailed job information
	flags := SlurmV0044GetJobsParamsFlagsDETAIL
	params.Flags = &flags

	// Note: v0.0.44 API has different filtering parameters
	// Filtering will be done post-retrieval for now
	_ = opts // Avoid unused variable warning

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0044GetJobsWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	// Handle response
	if resp.HTTPResponse.StatusCode != 200 {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal,
			fmt.Sprintf("HTTP %d", resp.HTTPResponse.StatusCode))
	}

	if resp.JSON200 == nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal, "Empty response body")
	}

	// Convert to interface types
	jobList := &interfaces.JobList{
		Jobs: make([]interfaces.Job, 0),
	}

	if resp.JSON200 != nil {
		for _, job := range resp.JSON200.Jobs {
			convertedJob := m.convertJobToInterface(&job)
			jobList.Jobs = append(jobList.Jobs, *convertedJob)
		}
	}

	return jobList, nil
}

// Get retrieves a specific job by ID
func (m *JobManagerImpl) Get(ctx context.Context, jobID string) (*interfaces.Job, error) {
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Use the job-specific endpoint
	resp, err := m.client.apiClient.SlurmV0044GetJobWithResponse(ctx, jobID, &SlurmV0044GetJobParams{})
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	if resp.HTTPResponse.StatusCode == 404 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound,
			fmt.Sprintf("Job %s not found", jobID))
	}

	if resp.HTTPResponse.StatusCode != 200 {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal,
			fmt.Sprintf("HTTP %d", resp.HTTPResponse.StatusCode))
	}

	if resp.JSON200 == nil || len(resp.JSON200.Jobs) == 0 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound,
			fmt.Sprintf("Job %s not found", jobID))
	}

	job := resp.JSON200.Jobs[0]
	return m.convertJobToInterface(&job), nil
}

// Submit submits a new job
func (m *JobManagerImpl) Submit(ctx context.Context, job *interfaces.JobSubmission) (*interfaces.JobSubmitResponse, error) {
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Convert job submission to v0.0.44 format
	jobDesc := &V0044JobDescMsg{
		Name:      &job.Name,
		Partition: &job.Partition,
	}

	// Set resource requirements
	if job.CPUs > 0 {
		cpus := int32(job.CPUs)
		jobDesc.CpusPerTask = &cpus
	}

	if job.Memory > 0 {
		memory := int64(job.Memory / (1024 * 1024)) // Convert bytes to MB
		setTrue := true
		jobDesc.MemoryPerNode = &V0044Uint64NoValStruct{
			Set:    &setTrue,
			Number: &memory,
		}
	}

	if job.TimeLimit > 0 {
		setTrue := true
		timeLimit := int32(job.TimeLimit)
		jobDesc.TimeLimit = &V0044Uint32NoValStruct{
			Set:    &setTrue,
			Number: &timeLimit,
		}
	}

	if job.Nodes > 0 {
		nodeSpec := fmt.Sprintf("%d", job.Nodes)
		jobDesc.Nodes = &nodeSpec
	}

	if job.WorkingDir != "" {
		jobDesc.CurrentWorkingDirectory = &job.WorkingDir
	}

	// Set environment variables
	if len(job.Environment) > 0 {
		envVars := make([]string, 0, len(job.Environment))
		for k, v := range job.Environment {
			envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
		}
		jobDesc.Environment = (*V0044StringArray)(&envVars)
	}

	submitReq := &V0044JobSubmitReq{
		Job:    jobDesc,
		Script: &job.Script,
	}

	// Submit the job
	resp, err := m.client.apiClient.SlurmV0044PostJobSubmitWithResponse(ctx, *submitReq)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	if resp.HTTPResponse.StatusCode != 200 {
		if resp.JSONDefault != nil && resp.JSONDefault.Errors != nil && len(*resp.JSONDefault.Errors) > 0 {
			errorMsg := fmt.Sprintf("Job submission failed: %v", (*resp.JSONDefault.Errors)[0])
			return nil, errors.NewSlurmError(errors.ErrorCodeValidationFailed, errorMsg)
		}
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal,
			fmt.Sprintf("HTTP %d", resp.HTTPResponse.StatusCode))
	}

	if resp.JSON200 == nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal, "Empty response body")
	}

	// Extract job ID from response
	var jobID string
	if resp.JSON200.JobId != nil {
		jobID = fmt.Sprintf("%d", *resp.JSON200.JobId)
	}

	return &interfaces.JobSubmitResponse{
		JobID: jobID,
	}, nil
}

// Cancel cancels a job
func (m *JobManagerImpl) Cancel(ctx context.Context, jobID string) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Use DELETE on the job endpoint
	resp, err := m.client.apiClient.SlurmV0044DeleteJobWithResponse(ctx, jobID, &SlurmV0044DeleteJobParams{})
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	if resp.HTTPResponse.StatusCode == 404 {
		return errors.NewSlurmError(errors.ErrorCodeResourceNotFound,
			fmt.Sprintf("Job %s not found", jobID))
	}

	if resp.HTTPResponse.StatusCode != 200 {
		return errors.NewSlurmError(errors.ErrorCodeServerInternal,
			fmt.Sprintf("HTTP %d", resp.HTTPResponse.StatusCode))
	}

	return nil
}

// Update updates job properties
func (m *JobManagerImpl) Update(ctx context.Context, jobID string, update *interfaces.JobUpdate) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Convert update to v0.0.44 format - using V0044JobDescMsg directly
	updateReq := V0044JobDescMsg{}

	// Set updateable fields
	if update.TimeLimit != nil {
		setTrue := true
		timeLimit := int32(*update.TimeLimit)
		updateReq.TimeLimit = &V0044Uint32NoValStruct{
			Set:    &setTrue,
			Number: &timeLimit,
		}
	}

	if update.Name != nil {
		updateReq.Name = update.Name
	}

	// Submit the update
	resp, err := m.client.apiClient.SlurmV0044PostJobWithResponse(ctx, jobID, updateReq)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	if resp.HTTPResponse.StatusCode != 200 {
		return errors.NewSlurmError(errors.ErrorCodeServerInternal,
			fmt.Sprintf("HTTP %d", resp.HTTPResponse.StatusCode))
	}

	return nil
}

// GetResourceLayout retrieves resource layout information for a job (new in v0.0.44)
// Note: This is a placeholder implementation until JobResourceLayout interface is defined
func (m *JobManagerImpl) GetResourceLayout(ctx context.Context, jobID string) (interface{}, error) {
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Use the new resources endpoint
	resp, err := m.client.apiClient.SlurmV0044GetResourcesWithResponse(ctx, jobID)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	if resp.HTTPResponse.StatusCode == 404 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound,
			fmt.Sprintf("Job %s not found", jobID))
	}

	if resp.HTTPResponse.StatusCode != 200 {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal,
			fmt.Sprintf("HTTP %d", resp.HTTPResponse.StatusCode))
	}

	if resp.JSON200 == nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal, "Empty response body")
	}

	// Return the raw response for now
	return resp.JSON200, nil
}

// Steps retrieves job steps for a job
func (m *JobManagerImpl) Steps(ctx context.Context, jobID string) (*interfaces.JobStepList, error) {
	// Implementation would call step-related endpoints
	// For now, return a basic implementation
	return &interfaces.JobStepList{
		Steps: make([]interfaces.JobStep, 0),
	}, nil
}

// Watch provides real-time job updates
func (m *JobManagerImpl) Watch(ctx context.Context, opts *interfaces.WatchJobsOptions) (<-chan interfaces.JobEvent, error) {
	// Create a channel for job events
	eventChan := make(chan interfaces.JobEvent)

	// For now, return a basic watcher that polls for changes
	go func() {
		defer close(eventChan)

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Poll for job changes - basic implementation
				// In a real implementation, this would use WebSocket or SSE
			}
		}
	}()

	return eventChan, nil
}

// convertJobToInterface converts a v0.0.44 job to the interface type
func (m *JobManagerImpl) convertJobToInterface(job *V0044JobInfo) *interfaces.Job {
	result := &interfaces.Job{}

	// Basic job information
	if job.JobId != nil {
		result.ID = fmt.Sprintf("%d", *job.JobId)
	}

	if job.Name != nil {
		result.Name = *job.Name
	}

	if job.UserId != nil {
		result.UserID = fmt.Sprintf("%d", *job.UserId)
	}

	if job.JobState != nil && len(*job.JobState) > 0 {
		result.State = string((*job.JobState)[0])
	}

	if job.Partition != nil {
		result.Partition = *job.Partition
	}

	if job.CurrentWorkingDirectory != nil {
		result.WorkingDir = *job.CurrentWorkingDirectory
	}

	// Resource information
	if job.Cpus != nil && job.Cpus.Number != nil {
		result.CPUs = int(*job.Cpus.Number)
	}

	// Memory handling (could be per node or per CPU)
	if job.MemoryPerNode != nil && job.MemoryPerNode.Number != nil {
		result.Memory = int(*job.MemoryPerNode.Number) // Already in MB, interfaces expect MB
	} else if job.MemoryPerCpu != nil && job.MemoryPerCpu.Number != nil && job.Cpus != nil && job.Cpus.Number != nil {
		totalMemory := int(*job.MemoryPerCpu.Number) * int(*job.Cpus.Number)
		result.Memory = totalMemory // MB
	}

	// Nodes handling - extract node count from string if possible
	if job.Nodes != nil {
		result.Nodes = []string{*job.Nodes}
	}

	if job.TimeLimit != nil && job.TimeLimit.Number != nil {
		result.TimeLimit = int(*job.TimeLimit.Number)
	}

	// Time information
	if job.SubmitTime != nil && job.SubmitTime.Number != nil {
		result.SubmitTime = time.Unix(*job.SubmitTime.Number, 0)
	}

	if job.StartTime != nil && job.StartTime.Number != nil {
		startTime := time.Unix(*job.StartTime.Number, 0)
		result.StartTime = &startTime
	}

	if job.EndTime != nil && job.EndTime.Number != nil {
		endTime := time.Unix(*job.EndTime.Number, 0)
		result.EndTime = &endTime
	}

	// Command
	if job.Command != nil {
		result.Command = *job.Command
	}

	// Environment would need to be parsed if available in v0.0.44

	return result
}

// Implementation stubs for analytics methods (maintained for compatibility)

func (m *JobManagerImpl) GetJobUtilization(ctx context.Context, jobID string) (*interfaces.JobUtilization, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "GetJobUtilization not implemented for v0.0.44")
}

func (m *JobManagerImpl) GetJobEfficiency(ctx context.Context, jobID string) (*interfaces.ResourceUtilization, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "GetJobEfficiency not implemented for v0.0.44")
}

func (m *JobManagerImpl) GetJobPerformance(ctx context.Context, jobID string) (*interfaces.JobPerformance, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "GetJobPerformance not implemented for v0.0.44")
}

func (m *JobManagerImpl) GetJobLiveMetrics(ctx context.Context, jobID string) (*interfaces.JobLiveMetrics, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "GetJobLiveMetrics not implemented for v0.0.44")
}

func (m *JobManagerImpl) WatchJobMetrics(ctx context.Context, jobID string, opts *interfaces.WatchMetricsOptions) (<-chan interfaces.JobMetricsEvent, error) {
	eventChan := make(chan interfaces.JobMetricsEvent)
	close(eventChan)
	return eventChan, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "WatchJobMetrics not implemented for v0.0.44")
}

func (m *JobManagerImpl) GetJobResourceTrends(ctx context.Context, jobID string, opts *interfaces.ResourceTrendsOptions) (*interfaces.JobResourceTrends, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "GetJobResourceTrends not implemented for v0.0.44")
}

func (m *JobManagerImpl) GetJobStepDetails(ctx context.Context, jobID string, stepID string) (*interfaces.JobStepDetails, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "GetJobStepDetails not implemented for v0.0.44")
}

func (m *JobManagerImpl) GetJobStepUtilization(ctx context.Context, jobID string, stepID string) (*interfaces.JobStepUtilization, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "GetJobStepUtilization not implemented for v0.0.44")
}

func (m *JobManagerImpl) ListJobStepsWithMetrics(ctx context.Context, jobID string, opts *interfaces.ListJobStepsOptions) (*interfaces.JobStepMetricsList, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "ListJobStepsWithMetrics not implemented for v0.0.44")
}

func (m *JobManagerImpl) GetJobStepsFromAccounting(ctx context.Context, jobID string, opts *interfaces.AccountingQueryOptions) (*interfaces.AccountingJobSteps, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "GetJobStepsFromAccounting not implemented for v0.0.44")
}

func (m *JobManagerImpl) GetStepAccountingData(ctx context.Context, jobID string, stepID string) (*interfaces.StepAccountingRecord, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "GetStepAccountingData not implemented for v0.0.44")
}

func (m *JobManagerImpl) GetJobStepAPIData(ctx context.Context, jobID string, stepID string) (*interfaces.JobStepAPIData, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "GetJobStepAPIData not implemented for v0.0.44")
}

func (m *JobManagerImpl) ListJobStepsFromSacct(ctx context.Context, jobID string, opts *interfaces.SacctQueryOptions) (*interfaces.SacctJobStepData, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "ListJobStepsFromSacct not implemented for v0.0.44")
}

func (m *JobManagerImpl) GetJobCPUAnalytics(ctx context.Context, jobID string) (*interfaces.CPUAnalytics, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "GetJobCPUAnalytics not implemented for v0.0.44")
}

func (m *JobManagerImpl) GetJobMemoryAnalytics(ctx context.Context, jobID string) (*interfaces.MemoryAnalytics, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "GetJobMemoryAnalytics not implemented for v0.0.44")
}

func (m *JobManagerImpl) GetJobIOAnalytics(ctx context.Context, jobID string) (*interfaces.IOAnalytics, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "GetJobIOAnalytics not implemented for v0.0.44")
}

func (m *JobManagerImpl) GetJobComprehensiveAnalytics(ctx context.Context, jobID string) (*interfaces.JobComprehensiveAnalytics, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "GetJobComprehensiveAnalytics not implemented for v0.0.44")
}

func (m *JobManagerImpl) GetJobPerformanceHistory(ctx context.Context, jobID string, opts *interfaces.PerformanceHistoryOptions) (*interfaces.JobPerformanceHistory, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "GetJobPerformanceHistory not implemented for v0.0.44")
}

func (m *JobManagerImpl) GetPerformanceTrends(ctx context.Context, opts *interfaces.TrendAnalysisOptions) (*interfaces.PerformanceTrends, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "GetPerformanceTrends not implemented for v0.0.44")
}

func (m *JobManagerImpl) GetUserEfficiencyTrends(ctx context.Context, userID string, opts *interfaces.EfficiencyTrendOptions) (*interfaces.UserEfficiencyTrends, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "GetUserEfficiencyTrends not implemented for v0.0.44")
}

func (m *JobManagerImpl) AnalyzeBatchJobs(ctx context.Context, jobIDs []string, opts *interfaces.BatchAnalysisOptions) (*interfaces.BatchJobAnalysis, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "AnalyzeBatchJobs not implemented for v0.0.44")
}

func (m *JobManagerImpl) GetWorkflowPerformance(ctx context.Context, workflowID string, opts *interfaces.WorkflowAnalysisOptions) (*interfaces.WorkflowPerformance, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "GetWorkflowPerformance not implemented for v0.0.44")
}

func (m *JobManagerImpl) GenerateEfficiencyReport(ctx context.Context, opts *interfaces.ReportOptions) (*interfaces.EfficiencyReport, error) {
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "GenerateEfficiencyReport not implemented for v0.0.44")
}

// Allocate allocates resources for a job (v0.0.44)
func (m *JobManagerImpl) Allocate(ctx context.Context, req *interfaces.JobAllocateRequest) (*interfaces.JobAllocateResponse, error) {
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// For now, return not implemented - allocation might require different endpoints
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Job allocation not yet implemented for v0.0.44")
}

// Hold holds a job (prevents it from running)
func (m *JobManagerImpl) Hold(ctx context.Context, jobID string) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// For now, return not implemented - hold might use different endpoints or parameters
	return errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Job hold not yet implemented for v0.0.44")
}

// Release releases a held job (allows it to run)
func (m *JobManagerImpl) Release(ctx context.Context, jobID string) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// For now, return not implemented - release might use different endpoints or parameters
	return errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Job release not yet implemented for v0.0.44")
}

// Notify sends a message to a job
func (m *JobManagerImpl) Notify(ctx context.Context, jobID string, message string) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// For now, return not implemented - notify might use different endpoints or parameters
	return errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Job notify not yet implemented for v0.0.44")
}

// Requeue requeues a job, allowing it to run again
func (m *JobManagerImpl) Requeue(ctx context.Context, jobID string) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// For now, return not implemented - requeue might use different endpoints or parameters
	return errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Job requeue not yet implemented for v0.0.44")
}

// Signal sends a signal to a job
func (m *JobManagerImpl) Signal(ctx context.Context, jobID string, signal string) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// For now, return not implemented - signal might use different endpoints or parameters
	return errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Job signal not yet implemented for v0.0.44")
}
