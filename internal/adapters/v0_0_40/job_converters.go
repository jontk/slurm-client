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
		job.TimeLimit = int32(*apiJob.TimeLimit.Number)
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
		job.Priority = int32(*apiJob.Priority.Number)
	}

	// Resource allocation
	if apiJob.Cpus != nil && apiJob.Cpus.Number != nil {
		job.CPUs = int32(*apiJob.Cpus.Number)
	}
	if apiJob.Nodes != nil {
		job.NodeList = *apiJob.Nodes
	}
	if apiJob.NodeCount != nil {
		if apiJob.NodeCount.Number != nil {
			job.Nodes = int32(*apiJob.NodeCount.Number)
		}
	}

	// Job specification
	if apiJob.Command != nil && len(*apiJob.Command) > 0 {
		job.Command = string(*apiJob.Command)
	}
	// Note: WorkingDirectory not available in V0040JobInfo but exists in V0040JobDesc
	// Skip for now as this conversion is from JobInfo (API response), not JobDesc (API request)
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
	if apiJob.ArrayJobId != nil && apiJob.ArrayJobId.Number != nil {
		arrayJobID := int32(*apiJob.ArrayJobId.Number)
		job.ArrayJobID = &arrayJobID
	}
	if apiJob.ArrayTaskId != nil && apiJob.ArrayTaskId.Number != nil {
		taskID := int32(*apiJob.ArrayTaskId.Number)
		job.ArrayTaskID = &taskID
	}
	if apiJob.ArrayTaskString != nil {
		job.ArrayTaskString = *apiJob.ArrayTaskString
	}

	// Note: Environment field not available in V0040JobInfo 
	// Skip for now as this conversion is from JobInfo (API response), not JobDesc (API request)

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
		job.ExcludeNodes = string(*apiJob.ExcludedNodes)
	}
	if apiJob.Nice != nil {
		job.Nice = *apiJob.Nice
	}
	if apiJob.Comment != nil {
		job.Comment = *apiJob.Comment
	}

	// Features and GRES
	if apiJob.Features != nil && len(*apiJob.Features) > 0 {
		job.Features = []string{string(*apiJob.Features)}
	}
	if apiJob.GresDetail != nil && len(*apiJob.GresDetail) > 0 {
		job.Gres = strings.Join(*apiJob.GresDetail, ",")
	}

	// Note: MemoryPerNode and MemoryPerCPU fields don't exist in common Job type
	// These would need to be mapped to available memory fields or skipped

	// Note: Requeue field not available in common Job type
	// This would need to be added to the common Job type or skipped

	// Note: Tasks, TasksPerNode, and CPUsPerTask fields not available in common Job type
	// These would need to be added to the common Job type or mapped to existing fields

	// Note: ExitCode field not available in common Job type
	// This would need to be added to the common Job type or skipped

	// Note: Tags field not available in common Job type  
	// This would need to be added to the common Job type or mapped to WCKey

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
		// Note: QosId field not available in V0040JobDescMsg
		// QoS handling would need different approach in v0.0.40
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

	// Note: Arguments field not available in common JobCreate type
	// This would need to be added or handled differently

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
		cpus := int32(job.CPUs)
		jobDesc.CpusPerTask = &cpus
	}
	if job.Nodes > 0 {
		// Note: NodeCount field not available in V0040JobDescMsg
		// Would need to use different field like RequiredNodes
	}
	// Note: Memory fields not available in common JobCreate type
	// These would need to be added to the common types

	// Time limit
	if job.TimeLimit > 0 {
		timeLimit := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: int64Ptr(int64(job.TimeLimit)),
		}
		jobDesc.TimeLimit = &timeLimit
	}

	// Priority
	if job.Priority != nil && *job.Priority > 0 {
		priority := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: int64Ptr(int64(*job.Priority)),
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

	// Note: Features and Gres fields not available in V0040JobDescMsg 
	// These would need different handling in v0.0.40

	// Note: NodeList field not available in common JobCreate type
	// This would need to be added to the common types

	// Exclude nodes
	if job.ExcludeNodes != "" {
		excludeNodes := []string{job.ExcludeNodes}
		jobDesc.ExcludedNodes = &excludeNodes
	}

	// Note: MailType field not available in common JobCreate type
	// This would need to be added to the common types
	if job.MailUser != "" {
		jobDesc.MailUser = &job.MailUser
	}

	// Comment
	if job.Comment != "" {
		jobDesc.Comment = &job.Comment
	}

	// Note: Tags, ArrayTaskString, and Requeue fields not available in common JobCreate type  
	// These would need to be added to the common types or handled differently

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
		// Note: QosId field not available in V0040JobDescMsg
		// QoS handling would need different approach in v0.0.40
	}
	if update.TimeLimit != nil {
		timeLimit := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: int64Ptr(int64(*update.TimeLimit)),
		}
		apiJob.TimeLimit = &timeLimit
	}
	if update.Priority != nil {
		priority := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: int64Ptr(int64(*update.Priority)),
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