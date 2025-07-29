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
		job.JobID = uint32(apiJob.JobId.Number)
	}
	if apiJob.Name != nil {
		job.Name = *apiJob.Name
	}
	if apiJob.UserId != nil {
		job.UserID = uint32(apiJob.UserId.Number)
	}
	if apiJob.UserName != nil {
		job.UserName = *apiJob.UserName
	}
	if apiJob.GroupId != nil {
		job.GroupID = uint32(apiJob.GroupId.Number)
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
	if apiJob.SubmitTime != nil {
		job.SubmitTime = time.Unix(int64(apiJob.SubmitTime.Number), 0)
	}
	if apiJob.StartTime != nil && apiJob.StartTime.Number > 0 {
		t := time.Unix(int64(apiJob.StartTime.Number), 0)
		job.StartTime = &t
	}
	if apiJob.EndTime != nil && apiJob.EndTime.Number > 0 {
		t := time.Unix(int64(apiJob.EndTime.Number), 0)
		job.EndTime = &t
	}
	if apiJob.TimeLimit != nil && apiJob.TimeLimit.Number > 0 {
		job.TimeLimit = uint32(apiJob.TimeLimit.Number)
	}
	if apiJob.TimeUsed != nil {
		job.TimeUsed = uint32(apiJob.TimeUsed.Number)
	}

	// Resource requirements
	if apiJob.Nodes != nil {
		job.NodeCount = uint32(*apiJob.Nodes)
	}
	if apiJob.NodeList != nil {
		job.NodeList = *apiJob.NodeList
	}
	if apiJob.Cpus != nil {
		job.CPUs = uint32(apiJob.Cpus.Number)
	}
	if apiJob.CpusPerTask != nil {
		job.CPUsPerTask = uint32(apiJob.CpusPerTask.Number)
	}
	if apiJob.TasksPerNode != nil {
		job.TasksPerNode = uint32(apiJob.TasksPerNode.Number)
	}
	if apiJob.MemoryPerCpu != nil {
		job.MemoryPerCPU = uint64(apiJob.MemoryPerCpu.Number)
	}
	if apiJob.MemoryPerNode != nil {
		job.MemoryPerNode = uint64(apiJob.MemoryPerNode.Number)
	}

	// Priority
	if apiJob.Priority != nil {
		job.Priority = uint32(apiJob.Priority.Number)
	}
	if apiJob.Nice != nil {
		job.Nice = int32(apiJob.Nice.Number)
	}

	// Job script and working directory
	if apiJob.Command != nil {
		job.Command = *apiJob.Command
	}
	if apiJob.CurrentWorkingDirectory != nil {
		job.WorkingDirectory = *apiJob.CurrentWorkingDirectory
	}
	if apiJob.StandardError != nil {
		job.StdErr = *apiJob.StandardError
	}
	if apiJob.StandardOutput != nil {
		job.StdOut = *apiJob.StandardOutput
	}
	if apiJob.StandardInput != nil {
		job.StdIn = *apiJob.StandardInput
	}

	// Exit code
	if apiJob.ExitCode != nil {
		if apiJob.ExitCode.ReturnCode != nil {
			job.ExitCode = int32(apiJob.ExitCode.ReturnCode.Number)
		}
	}

	// Array job information
	if apiJob.ArrayJobId != nil && apiJob.ArrayJobId.Number > 0 {
		arrayJobID := uint32(apiJob.ArrayJobId.Number)
		job.ArrayJobID = &arrayJobID
	}
	if apiJob.ArrayTaskId != nil && apiJob.ArrayTaskId.Number < ^uint32(0) {
		arrayTaskID := uint32(apiJob.ArrayTaskId.Number)
		job.ArrayTaskID = &arrayTaskID
	}
	if apiJob.ArrayMaxTasks != nil && apiJob.ArrayMaxTasks.Number > 0 {
		arrayMaxTasks := uint32(apiJob.ArrayMaxTasks.Number)
		job.ArrayMaxTasks = &arrayMaxTasks
	}
	if apiJob.ArrayTaskString != nil {
		job.ArrayTaskString = *apiJob.ArrayTaskString
	}

	// Dependencies
	if apiJob.Dependency != nil {
		job.Dependency = *apiJob.Dependency
	}

	// Features and constraints
	if apiJob.Features != nil {
		job.Features = *apiJob.Features
	}
	if apiJob.Reservation != nil {
		job.Reservation = *apiJob.Reservation
	}
	if apiJob.ExcludedNodes != nil {
		job.ExcludeNodes = *apiJob.ExcludedNodes
	}
	if apiJob.RequiredNodes != nil {
		job.RequiredNodes = *apiJob.RequiredNodes
	}

	// Environment variables
	if apiJob.Environment != nil {
		job.Environment = *apiJob.Environment
	}

	// TRES (Trackable RESources)
	if apiJob.TresAllocStr != nil {
		job.TRESAlloc = *apiJob.TresAllocStr
	}
	if apiJob.TresReqStr != nil {
		job.TRESReq = *apiJob.TresReqStr
	}

	// Licenses
	if apiJob.Licenses != nil {
		job.Licenses = *apiJob.Licenses
	}

	// Comments
	if apiJob.Comment != nil {
		job.Comment = *apiJob.Comment
	}
	if apiJob.AdminComment != nil {
		job.AdminComment = *apiJob.AdminComment
	}

	// Batch flags
	if apiJob.BatchFlag != nil {
		job.BatchFlag = *apiJob.BatchFlag
	}
	if apiJob.BatchHost != nil {
		job.BatchHost = *apiJob.BatchHost
	}

	// Requeue
	if apiJob.Requeue != nil {
		job.Requeue = *apiJob.Requeue
	}

	return job, nil
}

// convertCommonJobSubmitToAPI converts common job submission request to v0.0.42 API format
func (a *JobAdapter) convertCommonJobSubmitToAPI(req *types.JobSubmitRequest) (*api.SlurmV0042PostJobSubmitJSONRequestBody, error) {
	apiReq := &api.SlurmV0042PostJobSubmitJSONRequestBody{
		Jobs: &[]api.V0042JobSubmission{
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
		nodes := int32(req.Nodes)
		job.Nodes = &nodes
	}
	if req.CPUs > 0 {
		cpus := api.V0042Uint32NoValStruct{
			Set:    true,
			Number: uint64(req.CPUs),
		}
		job.Cpus = &cpus
	}
	if req.Memory != "" {
		job.MemoryPerNode = &req.Memory
	}
	if req.TimeLimit != "" {
		job.TimeLimit = &req.TimeLimit
	}

	// Working directory
	if req.WorkingDirectory != "" {
		job.CurrentWorkingDirectory = &req.WorkingDirectory
	}

	// Output files
	if req.Output != "" {
		job.StandardOutput = &req.Output
	}
	if req.Error != "" {
		job.StandardError = &req.Error
	}

	// Environment
	if len(req.Environment) > 0 {
		job.Environment = &req.Environment
	}

	// QoS
	if req.QoS != "" {
		job.Qos = &req.QoS
	}

	// Array
	if req.Array != "" {
		job.Array = &req.Array
	}

	// Dependencies
	if req.Dependency != "" {
		job.Dependency = &req.Dependency
	}

	// Features
	if req.Constraint != "" {
		job.Features = &req.Constraint
	}

	// Reservation
	if req.Reservation != "" {
		job.Reservation = &req.Reservation
	}

	// Exclusive
	if req.Exclusive {
		exclusive := "user"
		job.Exclusive = &exclusive
	}

	// Nice
	if req.Nice != 0 {
		nice := api.V0042Uint32NoValStruct{
			Set:    true,
			Number: uint64(req.Nice),
		}
		job.Nice = &nice
	}

	// Priority
	if req.Priority > 0 {
		priority := api.V0042Uint32NoValStruct{
			Set:    true,
			Number: uint64(req.Priority),
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
		result.JobID = uint32(*resp.JobId)
	}

	// Extract job submission details
	if resp.JobSubmitUserMsg != nil {
		result.Message = *resp.JobSubmitUserMsg
	}

	// Extract step ID if available
	if resp.StepId != nil {
		result.StepID = *resp.StepId
	}

	return result, nil
}

// convertCommonJobUpdateToAPI converts common job update request to v0.0.42 API format
// Note: v0.0.42 uses SlurmV0042PostJobJSONRequestBody for job updates
func (a *JobAdapter) convertCommonJobUpdateToAPI(req *types.JobUpdateRequest) (*api.SlurmV0042PostJobJSONRequestBody, error) {
	apiReq := &api.SlurmV0042PostJobJSONRequestBody{
		Jobs: &[]api.V0042JobInfo{
			{},
		},
	}

	job := &(*apiReq.Jobs)[0]

	// Priority
	if req.Priority != nil {
		priority := api.V0042Uint32NoValStruct{
			Set:    true,
			Number: uint64(*req.Priority),
		}
		job.Priority = &priority
	}

	// Nice
	if req.Nice != nil {
		nice := api.V0042Uint32NoValStruct{
			Set:    true,
			Number: uint64(*req.Nice),
		}
		job.Nice = &nice
	}

	// Time limit
	if req.TimeLimit != nil {
		timeLimit := api.V0042Uint32NoValStruct{
			Set:    true,
			Number: uint64(*req.TimeLimit),
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
	if req.NodeCount != nil {
		nodes := int32(*req.NodeCount)
		job.Nodes = &nodes
	}

	// Features
	if req.Features != nil {
		job.Features = req.Features
	}

	// Dependency
	if req.Dependency != nil {
		job.Dependency = req.Dependency
	}

	// Reservation
	if req.Reservation != nil {
		job.Reservation = req.Reservation
	}

	// Comment
	if req.Comment != nil {
		job.Comment = req.Comment
	}

	return apiReq, nil
}