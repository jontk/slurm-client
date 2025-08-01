// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

// convertAPIJobToCommon converts a v0.0.42 API Job to common Job type
func (a *JobAdapter) convertAPIJobToCommon(apiJob api.V0042JobInfo) (*types.Job, error) {
	job := &types.Job{}

	// Basic fields
	if apiJob.JobId != nil {
		job.JobID = *apiJob.JobId
	}
	if apiJob.Name != nil {
		job.Name = *apiJob.Name
	}
	if apiJob.UserId != nil {
		job.UserID = *apiJob.UserId
	}
	if apiJob.UserName != nil {
		job.UserName = *apiJob.UserName
	}
	if apiJob.GroupId != nil {
		job.GroupID = *apiJob.GroupId
	}
	if apiJob.Account != nil {
		job.Account = *apiJob.Account
	}
	if apiJob.Partition != nil {
		job.Partition = *apiJob.Partition
	}
	if apiJob.Qos != nil {
		job.QoS = *apiJob.Qos
	}

	// Job state
	if apiJob.JobState != nil && len(*apiJob.JobState) > 0 {
		job.State = types.JobState((*apiJob.JobState)[0])
	}
	if apiJob.StateReason != nil {
		job.StateReason = *apiJob.StateReason
	}

	// Time fields
	if apiJob.SubmitTime != nil && apiJob.SubmitTime.Number != nil {
		job.SubmitTime = time.Unix(*apiJob.SubmitTime.Number, 0)
	}
	if apiJob.StartTime != nil && apiJob.StartTime.Number != nil && *apiJob.StartTime.Number > 0 {
		t := time.Unix(*apiJob.StartTime.Number, 0)
		job.StartTime = &t
	}
	if apiJob.EndTime != nil && apiJob.EndTime.Number != nil && *apiJob.EndTime.Number > 0 {
		t := time.Unix(*apiJob.EndTime.Number, 0)
		job.EndTime = &t
	}
	if apiJob.TimeLimit != nil && apiJob.TimeLimit.Number != nil && *apiJob.TimeLimit.Number > 0 {
		job.TimeLimit = *apiJob.TimeLimit.Number
	}

	// Resource requirements
	if apiJob.Nodes != nil {
		// apiJob.Nodes is a string, need to parse to get count
		// job.Nodes = parseNodeString(*apiJob.Nodes)
	}
	// NodeList is from JobResources in v0.0.42
	if apiJob.JobResources != nil && apiJob.JobResources.Nodes != nil && apiJob.JobResources.Nodes.List != nil {
		job.NodeList = *apiJob.JobResources.Nodes.List
	}
	if apiJob.Cpus != nil && apiJob.Cpus.Number != nil {
		job.CPUs = *apiJob.Cpus.Number
	}
	if apiJob.CpusPerTask != nil && apiJob.CpusPerTask.Number != nil {
		job.ResourceRequests.CPUsPerTask = *apiJob.CpusPerTask.Number
	}
	if apiJob.TasksPerNode != nil && apiJob.TasksPerNode.Number != nil {
		job.ResourceRequests.TasksPerNode = *apiJob.TasksPerNode.Number
	}
	if apiJob.MemoryPerCpu != nil && apiJob.MemoryPerCpu.Number != nil {
		job.ResourceRequests.MemoryPerCPU = int64(*apiJob.MemoryPerCpu.Number)
	}
	if apiJob.MemoryPerNode != nil && apiJob.MemoryPerNode.Number != nil {
		// MemoryPerNode field doesn't exist in Job struct, store in ResourceRequests.Memory
		job.ResourceRequests.Memory = *apiJob.MemoryPerNode.Number
	}

	// Priority
	if apiJob.Priority != nil && apiJob.Priority.Number != nil {
		job.Priority = *apiJob.Priority.Number
	}
	if apiJob.Nice != nil {
		job.Nice = *apiJob.Nice
	}

	// Job script and working directory
	if apiJob.Command != nil {
		job.Command = *apiJob.Command
	}
	if apiJob.CurrentWorkingDirectory != nil {
		job.WorkingDirectory = *apiJob.CurrentWorkingDirectory
	}
	if apiJob.StandardError != nil {
		job.StandardError = *apiJob.StandardError
	}
	if apiJob.StandardOutput != nil {
		job.StandardOutput = *apiJob.StandardOutput
	}
	if apiJob.StandardInput != nil {
		job.StandardInput = *apiJob.StandardInput
	}

	// Exit code - ExitCode doesn't exist in the common Job type
	// We'll skip this field as it's typically only relevant for completed jobs

	// Array job information
	if apiJob.ArrayJobId != nil && apiJob.ArrayJobId.Number != nil && *apiJob.ArrayJobId.Number > 0 {
		arrayJobID := int32(*apiJob.ArrayJobId.Number)
		job.ArrayJobID = &arrayJobID
	}
	if apiJob.ArrayTaskId != nil && apiJob.ArrayTaskId.Number != nil {
		arrayTaskID := int32(*apiJob.ArrayTaskId.Number)
		job.ArrayTaskID = &arrayTaskID
	}
	// ArrayMaxTasks field doesn't exist in common Job type, skip it
	if apiJob.ArrayTaskString != nil {
		job.ArrayTaskString = *apiJob.ArrayTaskString
	}

	// Dependencies
	if apiJob.Dependency != nil {
		// Dependencies stored as a simple string in v0.0.42
		// We'll need to parse this into JobDependency structures later
		// For now, we'll skip the field as it doesn't match the common type
	}

	// Features and constraints
	if apiJob.Features != nil {
		// Features is a string in v0.0.42, need to convert to slice
		job.Features = []string{*apiJob.Features}
	}
	// Reservation field exists in Job struct
	if apiJob.ResvName != nil {
		job.Reservation = *apiJob.ResvName
	}
	if apiJob.ExcludedNodes != nil {
		job.ExcludeNodes = *apiJob.ExcludedNodes
	}
	// RequiredNodes field doesn't exist in common Job type
	
	// Environment variables - V0042JobInfo doesn't have Environment field
	// TODO: Check if environment is available in a different field

	// TRES (Trackable RESources) - not in common Job type
	
	// Licenses - not in common Job type
	
	// Comments
	if apiJob.Comment != nil {
		job.Comment = *apiJob.Comment
	}
	// AdminComment field doesn't exist in common Job type
	
	// Batch fields
	if apiJob.BatchHost != nil {
		job.BatchHost = *apiJob.BatchHost
	}
	// BatchFlag and Requeue fields don't exist in common Job type

	return job, nil
}

// convertCommonJobSubmitToAPI converts common job submission request to v0.0.42 API format
func (a *JobAdapter) convertCommonJobSubmitToAPI(req *types.JobSubmitRequest) (*api.V0042JobSubmitReq, error) {
	apiReq := &api.V0042JobSubmitReq{
		Jobs: &[]api.V0042JobDescMsg{
			{
				Name:      &req.Name,
				Account:   &req.Account,
				Partition: &req.Partition,
			},
		},
	}

	job := &(*apiReq.Jobs)[0]

	// Script
	if req.Script != "" {
		job.Script = &req.Script
	}

	// Resources
	if req.Nodes > 0 {
		// Convert nodes count to string (e.g., "4" for 4 nodes)
		nodes := strconv.Itoa(int(req.Nodes))
		job.Nodes = &nodes
	}
	if req.CPUs > 0 {
		cpuNum := int32(req.CPUs)
		job.CpusPerTask = &cpuNum
	}
	if req.ResourceRequests.Memory > 0 {
		// Check for MemoryPerNode field type
		// Since MemoryPerNode expects V0042Uint64NoValStruct, create it
		setTrue := true
		memoryMB := int64(req.ResourceRequests.Memory)
		memory := api.V0042Uint64NoValStruct{
			Set:    &setTrue,
			Number: &memoryMB,
		}
		job.MemoryPerNode = &memory
	}
	if req.TimeLimit > 0 {
		// Convert minutes to V0042Uint32NoValStruct
		setTrue := true
		timeLimitMinutes := int32(req.TimeLimit)
		timeLimit := api.V0042Uint32NoValStruct{
			Set:    &setTrue,
			Number: &timeLimitMinutes,
		}
		job.TimeLimit = &timeLimit
	}

	// Working directory
	if req.WorkingDirectory != "" {
		job.CurrentWorkingDirectory = &req.WorkingDirectory
	}

	// Output files
	if req.StandardOutput != "" {
		job.StandardOutput = &req.StandardOutput
	}
	if req.StandardError != "" {
		job.StandardError = &req.StandardError
	}

	// Environment
	if len(req.Environment) > 0 {
		// Convert map to array of "KEY=VALUE" strings
		envArray := make([]string, 0, len(req.Environment))
		for key, value := range req.Environment {
			envArray = append(envArray, fmt.Sprintf("%s=%s", key, value))
		}
		job.Environment = &envArray
	}

	// QoS
	if req.QoS != "" {
		job.Qos = &req.QoS
	}

	// Array
	if req.ArrayString != "" {
		job.Array = &req.ArrayString
	}

	// Dependencies
	if len(req.Dependencies) > 0 {
		// Convert dependencies to string format
		// TODO: Implement proper dependency conversion
		dep := fmt.Sprintf("afterok:%d", req.Dependencies[0].JobIDs[0])
		job.Dependency = &dep
	}

	// Features  
	if len(req.Features) > 0 {
		// Convert feature list to string
		features := strings.Join(req.Features, "&")
		job.Constraints = &features
	}

	// Reservation
	if req.Reservation != "" {
		job.Reservation = &req.Reservation
	}

	// Shared/Exclusive settings
	if req.Shared != "" {
		// V0042JobShared is []string
		shared := []string{req.Shared}
		job.Shared = (*api.V0042JobShared)(&shared)
	}

	// Nice
	if req.Nice != 0 {
		niceNum := int32(req.Nice)
		job.Nice = &niceNum
	}

	// Priority
	if req.Priority != nil && *req.Priority > 0 {
		setTrue := true
		priorityNum := int32(*req.Priority)
		priority := api.V0042Uint32NoValStruct{
			Set:    &setTrue,
			Number: &priorityNum,
		}
		job.Priority = &priority
	}

	return apiReq, nil
}

// convertAPIJobSubmitResponseToCommon converts v0.0.42 API job submit response to common format
func (a *JobAdapter) convertAPIJobSubmitResponseToCommon(resp *api.V0042OpenapiJobSubmitResponse) (*types.JobSubmitResponse, error) {
	result := &types.JobSubmitResponse{}

	// Check for errors
	if resp.Errors != nil && len(*resp.Errors) > 0 {
		errMsgs := make([]string, 0, len(*resp.Errors))
		for _, e := range *resp.Errors {
			if e.Error != nil {
				errMsgs = append(errMsgs, *e.Error)
			}
		}
		if len(errMsgs) > 0 {
			return nil, fmt.Errorf("job submission failed: %s", strings.Join(errMsgs, "; "))
		}
	}

	// Extract job ID
	if resp.JobId != nil {
		result.JobID = int32(*resp.JobId)
	}

	// Extract job submission details
	if resp.JobSubmitUserMsg != nil {
		result.JobSubmitUserMsg = *resp.JobSubmitUserMsg
	}

	// Extract step ID if available
	if resp.StepId != nil {
		result.StepID = *resp.StepId
	}

	return result, nil
}

// convertCommonJobUpdateToAPI converts common job update request to v0.0.42 API format
// Note: v0.0.42 uses V0042JobDescMsg for job updates
func (a *JobAdapter) convertCommonJobUpdateToAPI(req *types.JobUpdate) (*api.V0042JobDescMsg, error) {
	job := &api.V0042JobDescMsg{}

	// Priority
	if req.Priority != nil {
		priorityNumber := int32(*req.Priority)
		priority := api.V0042Uint32NoValStruct{
			Number: &priorityNumber,
			Set:    &[]bool{true}[0],
		}
		job.Priority = &priority
	}

	// Nice
	if req.Nice != nil {
		job.Nice = req.Nice
	}

	// Time limit
	if req.TimeLimit != nil {
		timeLimitNumber := *req.TimeLimit
		timeLimit := api.V0042Uint32NoValStruct{
			Number: &timeLimitNumber,
			Set:    &[]bool{true}[0],
		}
		job.TimeLimit = &timeLimit
	}

	// Partition
	if req.Partition != nil {
		job.Partition = req.Partition
	}

	// QoS
	if req.QoS != nil {
		job.Qos = req.QoS
	}

	// Account
	if req.Account != nil {
		job.Account = req.Account
	}

	// Node count
	if req.MinNodes != nil {
		nodes := fmt.Sprintf("%d", *req.MinNodes)
		job.Nodes = &nodes
	}

	// Features -> Constraints in v0.0.42
	if len(req.Features) > 0 {
		features := strings.Join(req.Features, ",")
		job.Constraints = &features
	}

	// Dependency
	if req.Name != nil {
		// Name field exists in JobUpdate but not dependency
		job.Name = req.Name
	}

	// Reservation
	if req.Reservation != nil {
		job.Reservation = req.Reservation
	}

	// Comment
	if req.Comment != nil {
		job.Comment = req.Comment
	}

	return job, nil
}
