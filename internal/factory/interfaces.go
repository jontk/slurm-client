// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package factory

import (
	"github.com/jontk/slurm-client/internal/interfaces"
)

// SlurmClient represents a version-agnostic Slurm REST API client interface
// This is a type alias to the internal interface to avoid import cycles
type SlurmClient = interfaces.SlurmClient

// Type aliases for all interfaces and data structures
type JobManager = interfaces.JobManager
type NodeManager = interfaces.NodeManager
type PartitionManager = interfaces.PartitionManager
type InfoManager = interfaces.InfoManager

// Data structure type aliases
type Job = interfaces.Job
type JobList = interfaces.JobList
type JobSubmission = interfaces.JobSubmission
type JobSubmitResponse = interfaces.JobSubmitResponse
type JobUpdate = interfaces.JobUpdate
type JobStep = interfaces.JobStep
type JobStepList = interfaces.JobStepList
type JobEvent = interfaces.JobEvent

type Node = interfaces.Node
type NodeList = interfaces.NodeList
type NodeUpdate = interfaces.NodeUpdate
type NodeEvent = interfaces.NodeEvent

type Partition = interfaces.Partition
type PartitionList = interfaces.PartitionList
type PartitionUpdate = interfaces.PartitionUpdate
type PartitionEvent = interfaces.PartitionEvent

type ClusterInfo = interfaces.ClusterInfo
type ClusterStats = interfaces.ClusterStats
type APIVersion = interfaces.APIVersion

// List options
type ListJobsOptions = interfaces.ListJobsOptions
type ListNodesOptions = interfaces.ListNodesOptions
type ListPartitionsOptions = interfaces.ListPartitionsOptions

// Watch options
type WatchJobsOptions = interfaces.WatchJobsOptions
type WatchNodesOptions = interfaces.WatchNodesOptions
type WatchPartitionsOptions = interfaces.WatchPartitionsOptions

// ClientConfig for API client configuration
type ClientConfig = interfaces.ClientConfig
