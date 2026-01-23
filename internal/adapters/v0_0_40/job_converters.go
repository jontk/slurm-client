// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"strings"
	"time"

	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
	"github.com/jontk/slurm-client/internal/common/types"
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
	if apiJob.Command != nil && *apiJob.Command != "" {
		job.Command = *apiJob.Command
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
		copy(mailTypes, *apiJob.MailType)
		job.MailType = mailTypes
	}
	if apiJob.MailUser != nil {
		job.MailUser = *apiJob.MailUser
	}

	// Additional fields
	if apiJob.ExcludedNodes != nil && *apiJob.ExcludedNodes != "" {
		job.ExcludeNodes = *apiJob.ExcludedNodes
	}
	if apiJob.Nice != nil {
		job.Nice = *apiJob.Nice
	}
	if apiJob.Comment != nil {
		job.Comment = *apiJob.Comment
	}

	// Features and GRES
	if apiJob.Features != nil && *apiJob.Features != "" {
		job.Features = []string{*apiJob.Features}
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

// Helper functions
func boolPtr(b bool) *bool {
	return &b
}
