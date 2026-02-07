// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package watch

import (
	types "github.com/jontk/slurm-client/api"
)

// getJobID safely extracts the job ID from a Job, returning 0 if nil
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

// getNodeName safely extracts the node name from a Node, returning empty string if nil
func getNodeName(node *types.Node) string {
	if node == nil || node.Name == nil {
		return ""
	}
	return *node.Name
}

// getNodeState safely extracts the first node state from a Node
func getNodeState(node *types.Node) types.NodeState {
	if node == nil || len(node.State) == 0 {
		return ""
	}
	return node.State[0]
}

// getPartitionName safely extracts the partition name from a Partition
func getPartitionName(partition *types.Partition) string {
	if partition == nil || partition.Name == nil {
		return ""
	}
	return *partition.Name
}

// getPartitionState safely extracts the first partition state from a Partition
func getPartitionState(partition *types.Partition) types.PartitionState {
	if partition == nil || partition.Partition == nil || len(partition.Partition.State) == 0 {
		return ""
	}
	return types.PartitionState(partition.Partition.State[0])
}
