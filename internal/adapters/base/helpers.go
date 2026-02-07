// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package base

import (
	types "github.com/jontk/slurm-client/api"
)

// String pointer helpers
// derefString safely dereferences a string pointer, returning empty string if nil
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// derefBool safely dereferences a bool pointer, returning false if nil
func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// Job field helpers
// getJobID safely extracts the job ID from a Job
func getJobID(job *types.Job) int32 {
	if job == nil || job.JobID == nil {
		return 0
	}
	return *job.JobID
}

// getJobState safely extracts the first job state from a Job
func getJobState(job *types.Job) types.JobState {
	if job == nil || len(job.JobState) == 0 {
		return ""
	}
	return job.JobState[0]
}

// Account field helpers
// isAccountDeleted checks if an account has the DELETED flag
func isAccountDeleted(account types.Account) bool {
	for _, flag := range account.Flags {
		if flag == types.AccountFlagsDeleted {
			return true
		}
	}
	return false
}

// Partition field helpers
// getPartitionName safely extracts the partition name from a Partition
func getPartitionName(partition *types.Partition) string {
	if partition == nil || partition.Name == nil {
		return ""
	}
	return *partition.Name
}

// getPartitionStates safely extracts all partition states from a Partition
func getPartitionStates(partition *types.Partition) []types.PartitionState {
	if partition == nil || partition.Partition == nil {
		return nil
	}
	states := make([]types.PartitionState, len(partition.Partition.State))
	for i, s := range partition.Partition.State {
		states[i] = types.PartitionState(s)
	}
	return states
}
