package v0_0_40

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
)

// convertAPIJobToCommon converts a v0.0.40 API Job to common Job type
func (a *JobAdapter) convertAPIJobToCommon(apiJob api.V0040JobInfo) (*types.Job, error) {
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
	if apiJob.TimeLimit != nil && apiJob.TimeLimit.Number != nil {
		job.TimeLimit = *apiJob.TimeLimit.Number
	}
	if apiJob.SubmitTime != nil && apiJob.SubmitTime.Number != nil {
		job.SubmitTime = time.Unix(*apiJob.SubmitTime.Number, 0)
	}
	if apiJob.StartTime != nil && apiJob.StartTime.Number != nil && *apiJob.StartTime.Number > 0 {
		startTime := time.Unix(*apiJob.StartTime.Number, 0)
		job.StartTime = &startTime
	}
	if apiJob.EndTime != nil && apiJob.EndTime.Number != nil && *apiJob.EndTime.Number > 0 {
		endTime := time.Unix(*apiJob.EndTime.Number, 0)
		job.EndTime = &endTime
	}

	// Priority
	if apiJob.Priority != nil && apiJob.Priority.Number != nil {
		job.Priority = *apiJob.Priority.Number
	}

	// Resource allocation
	if apiJob.Cpus != nil && apiJob.Cpus.Number != nil {
		job.CPUs = *apiJob.Cpus.Number
	}
	if apiJob.Nodes != nil {
		job.NodeList = *apiJob.Nodes
	}
	if apiJob.NodeCount != nil {
		if apiJob.NodeCount.Number != nil {
			job.Nodes = *apiJob.NodeCount.Number
		}
	}

	// Job specification
	if apiJob.Command != nil && len(*apiJob.Command) > 0 {
		job.Command = (*apiJob.Command)[0]
	}
	if apiJob.WorkingDirectory != nil {
		job.WorkingDirectory = *apiJob.WorkingDirectory
	}
	if apiJob.StandardInput != nil {
		job.StandardInput = *apiJob.StandardInput
	}
	if apiJob.StandardOutput != nil {
		job.StandardOutput = *apiJob.StandardOutput
	}
	if apiJob.StandardError != nil {
		job.StandardError = *apiJob.StandardError
	}

	// Array job information
	if apiJob.ArrayJobId != nil {
		job.ArrayJobID = apiJob.ArrayJobId
	}
	if apiJob.ArrayTaskId != nil && apiJob.ArrayTaskId.Number != nil {
		taskID := *apiJob.ArrayTaskId.Number
		job.ArrayTaskID = &taskID
	}
	if apiJob.ArrayTaskString != nil {
		job.ArrayTaskString = *apiJob.ArrayTaskString
	}

	// Environment
	if apiJob.Environment != nil && len(*apiJob.Environment) > 0 {
		env := make(map[string]string)
		for _, envVar := range *apiJob.Environment {
			parts := strings.SplitN(envVar, "=", 2)
			if len(parts) == 2 {
				env[parts[0]] = parts[1]
			}
		}
		job.Environment = env
	}

	// Mail settings
	if apiJob.MailType != nil && len(*apiJob.MailType) > 0 {
		mailTypes := make([]string, len(*apiJob.MailType))
		for i, mt := range *apiJob.MailType {
			mailTypes[i] = string(mt)
		}
		job.MailType = mailTypes
	}
	if apiJob.MailUser != nil {
		job.MailUser = *apiJob.MailUser
	}

	// Additional fields
	if apiJob.ExcludedNodes != nil && len(*apiJob.ExcludedNodes) > 0 {
		job.ExcludeNodes = (*apiJob.ExcludedNodes)[0]
	}
	if apiJob.Nice != nil {
		job.Nice = *apiJob.Nice
	}
	if apiJob.Comment != nil {
		job.Comment = *apiJob.Comment
	}

	// Features and GRES
	if apiJob.Features != nil && len(*apiJob.Features) > 0 {
		job.Features = (*apiJob.Features)[0]
	}
	if apiJob.GresDetail != nil && len(*apiJob.GresDetail) > 0 {
		job.GRES = strings.Join(*apiJob.GresDetail, ",")
	}

	// Memory specifications
	if apiJob.MemoryPerNode != nil && apiJob.MemoryPerNode.Number != nil {
		memPerNode := int64(*apiJob.MemoryPerNode.Number)
		job.MemoryPerNode = &memPerNode
	}
	if apiJob.MemoryPerCpu != nil && apiJob.MemoryPerCpu.Number != nil {
		memPerCPU := int64(*apiJob.MemoryPerCpu.Number)
		job.MemoryPerCPU = &memPerCPU
	}

	// Requeue flag
	if apiJob.Requeue != nil {
		job.Requeue = *apiJob.Requeue
	}

	// Tasks
	if apiJob.Tasks != nil && apiJob.Tasks.Number != nil {
		job.Tasks = *apiJob.Tasks.Number
	}
	if apiJob.TasksPerNode != nil && apiJob.TasksPerNode.Number != nil {
		job.TasksPerNode = *apiJob.TasksPerNode.Number
	}
	if apiJob.CpusPerTask != nil && apiJob.CpusPerTask.Number != nil {
		job.CPUsPerTask = *apiJob.CpusPerTask.Number
	}

	// Exit code
	if apiJob.ExitCode != nil && apiJob.ExitCode.Status != nil && len(*apiJob.ExitCode.Status) > 0 {
		// Try to parse the first status as an exit code
		statusStr := (*apiJob.ExitCode.Status)[0]
		if parts := strings.Split(statusStr, ":"); len(parts) == 2 {
			if code, err := strconv.Atoi(parts[0]); err == nil {
				exitCode := int32(code)
				job.ExitCode = &exitCode
			}
		}
	}

	// Tags
	if apiJob.Wckey != nil {
		job.Tags = []string{*apiJob.Wckey}
	}

	return job, nil
}

// convertCommonJobCreateToAPI converts common JobCreate to v0.0.40 API format
func (a *JobAdapter) convertCommonJobCreateToAPI(job *types.JobCreate) (*api.V0040JobSubmitReq, error) {
	apiJob := &api.V0040JobSubmitReq{}

	// Set the job description
	jobDesc := api.V0040JobDescMsg{}

	// Basic fields
	if job.Name != "" {
		jobDesc.Name = &job.Name
	}
	if job.Account != "" {
		jobDesc.Account = &job.Account
	}
	if job.Partition != "" {
		jobDesc.Partition = &job.Partition
	}
	if job.QoS != "" {
		jobDesc.QosId = &job.QoS
	}

	// Command or script
	if job.Script != "" {
		apiJob.Script = &job.Script
	} else if job.Command != "" {
		// For command, we need to handle it differently
		// In v0.0.40, command goes in argv
		argv := []string{job.Command}
		jobDesc.Argv = &argv
	}

	// Arguments
	if len(job.Arguments) > 0 {
		if jobDesc.Argv == nil {
			jobDesc.Argv = &[]string{}
		}
		*jobDesc.Argv = append(*jobDesc.Argv, job.Arguments...)
	}

	// Working directory
	if job.WorkingDirectory != "" {
		jobDesc.CurrentWorkingDirectory = &job.WorkingDirectory
	}

	// I/O redirection
	if job.StandardInput != "" {
		jobDesc.StandardInput = &job.StandardInput
	}
	if job.StandardOutput != "" {
		jobDesc.StandardOutput = &job.StandardOutput
	}
	if job.StandardError != "" {
		jobDesc.StandardError = &job.StandardError
	}

	// Resources
	if job.CPUs > 0 {
		cpus := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: int32Ptr(job.CPUs),
		}
		jobDesc.CpusPerTask = &cpus
	}
	if job.Nodes > 0 {
		nodes := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: int32Ptr(job.Nodes),
		}
		jobDesc.NodeCount = &nodes
	}
	if job.Tasks > 0 {
		tasks := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: int32Ptr(job.Tasks),
		}
		jobDesc.Tasks = &tasks
	}
	if job.MemoryPerNode != nil && *job.MemoryPerNode > 0 {
		mem := api.V0040Uint64NoVal{
			Set:    boolPtr(true),
			Number: int64Ptr(*job.MemoryPerNode),
		}
		jobDesc.MemoryPerNode = &mem
	}
	if job.MemoryPerCPU != nil && *job.MemoryPerCPU > 0 {
		mem := api.V0040Uint64NoVal{
			Set:    boolPtr(true),
			Number: int64Ptr(*job.MemoryPerCPU),
		}
		jobDesc.MemoryPerCpu = &mem
	}

	// Time limit
	if job.TimeLimit > 0 {
		timeLimit := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: int32Ptr(job.TimeLimit),
		}
		jobDesc.TimeLimit = &timeLimit
	}

	// Priority
	if job.Priority > 0 {
		priority := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: int32Ptr(job.Priority),
		}
		jobDesc.Priority = &priority
	}

	// Environment
	if len(job.Environment) > 0 {
		env := make([]string, 0, len(job.Environment))
		for k, v := range job.Environment {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		jobDesc.Environment = &env
	}

	// Features
	if job.Features != "" {
		features := []string{job.Features}
		jobDesc.Features = &features
	}

	// GRES
	if job.GRES != "" {
		jobDesc.Gres = &job.GRES
	}

	// Node list
	if job.NodeList != "" {
		jobDesc.RequiredNodes = &job.NodeList
	}

	// Exclude nodes
	if job.ExcludeNodes != "" {
		excludeNodes := []string{job.ExcludeNodes}
		jobDesc.ExcludedNodes = &excludeNodes
	}

	// Mail settings
	if len(job.MailType) > 0 {
		mailFlags := make([]api.V0040JobMailFlags, len(job.MailType))
		for i, mt := range job.MailType {
			mailFlags[i] = api.V0040JobMailFlags(mt)
		}
		jobDesc.MailType = &mailFlags
	}
	if job.MailUser != "" {
		jobDesc.MailUser = &job.MailUser
	}

	// Comment
	if job.Comment != "" {
		jobDesc.Comment = &job.Comment
	}

	// Tags
	if len(job.Tags) > 0 && job.Tags[0] != "" {
		jobDesc.Wckey = &job.Tags[0]
	}

	// Array job
	if job.ArrayTaskString != "" {
		jobDesc.Array = &job.ArrayTaskString
	}

	// Requeue
	if job.Requeue {
		jobDesc.Requeue = &job.Requeue
	}

	// Set the job in the submit request
	apiJob.Job = &jobDesc

	return apiJob, nil
}

// convertCommonJobUpdateToAPI converts common JobUpdate to v0.0.40 API format
func (a *JobAdapter) convertCommonJobUpdateToAPI(existingJob *types.Job, update *types.JobUpdate) (*api.V0040JobDescMsg, error) {
	apiJob := &api.V0040JobDescMsg{}

	// Job ID
	jobID := existingJob.JobID
	apiJob.JobId = &jobID

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
	if update.QoS != nil {
		apiJob.QosId = update.QoS
	}
	if update.TimeLimit != nil {
		timeLimit := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: int32Ptr(*update.TimeLimit),
		}
		apiJob.TimeLimit = &timeLimit
	}
	if update.Priority != nil {
		priority := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: int32Ptr(*update.Priority),
		}
		apiJob.Priority = &priority
	}
	if update.Nice != nil {
		apiJob.Nice = update.Nice
	}
	if update.Comment != nil {
		apiJob.Comment = update.Comment
	}

	return apiJob, nil
}

// Helper functions
func boolPtr(b bool) *bool {
	return &b
}

func int32Ptr(i int32) *int32 {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}