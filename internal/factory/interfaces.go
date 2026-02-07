// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package factory

import (
	types "github.com/jontk/slurm-client/api"
)

// SlurmClient represents a version-agnostic Slurm REST API client interface
// This is a type alias to the public interface to avoid import cycles
type SlurmClient = types.SlurmClient

// Type aliases for all interfaces and data structures
type JobManager = types.JobManager
type NodeManager = types.NodeManager
type PartitionManager = types.PartitionManager
type InfoManager = types.InfoManager

// Data structure type aliases
type Job = types.Job
type JobList = types.JobList
type JobSubmission = types.JobSubmission
type JobSubmitResponse = types.JobSubmitResponse
type JobUpdate = types.JobUpdate
type JobStep = types.JobStep
type JobStepList = types.JobStepList
type JobEvent = types.JobEvent

type Node = types.Node
type NodeList = types.NodeList
type NodeUpdate = types.NodeUpdate
type NodeEvent = types.NodeEvent

type Partition = types.Partition
type PartitionList = types.PartitionList
type PartitionUpdate = types.PartitionUpdate
type PartitionEvent = types.PartitionEvent

type ClusterInfo = types.ClusterInfo
type ClusterStats = types.ClusterStats
type APIVersion = types.APIVersion

// List options
type ListJobsOptions = types.ListJobsOptions
type ListNodesOptions = types.ListNodesOptions
type ListPartitionsOptions = types.ListPartitionsOptions

// Watch options
type WatchJobsOptions = types.WatchJobsOptions
type WatchNodesOptions = types.WatchNodesOptions
type WatchPartitionsOptions = types.WatchPartitionsOptions

// ClientConfig for API client configuration
type ClientConfig = types.ClientConfig
