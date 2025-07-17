package slurm

import (
	"github.com/jontk/slurm-client/internal/interfaces"
)

// SlurmClient represents a version-agnostic Slurm REST API client
// This is a type alias to the internal interface to avoid import cycles
type SlurmClient = interfaces.SlurmClient

// JobManager provides version-agnostic job operations
type JobManager = interfaces.JobManager

// NodeManager provides version-agnostic node operations
type NodeManager = interfaces.NodeManager

// PartitionManager provides version-agnostic partition operations
type PartitionManager = interfaces.PartitionManager

// InfoManager provides version-agnostic cluster information operations
type InfoManager = interfaces.InfoManager

// ReservationManager provides version-agnostic reservation operations
type ReservationManager = interfaces.ReservationManager

// Type aliases for data structures
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

type Reservation = interfaces.Reservation
type ReservationList = interfaces.ReservationList
type ReservationCreate = interfaces.ReservationCreate
type ReservationCreateResponse = interfaces.ReservationCreateResponse
type ReservationUpdate = interfaces.ReservationUpdate

// List options
type ListJobsOptions = interfaces.ListJobsOptions
type ListNodesOptions = interfaces.ListNodesOptions
type ListPartitionsOptions = interfaces.ListPartitionsOptions
type ListReservationsOptions = interfaces.ListReservationsOptions

// Watch options
type WatchJobsOptions = interfaces.WatchJobsOptions
type WatchNodesOptions = interfaces.WatchNodesOptions
type WatchPartitionsOptions = interfaces.WatchPartitionsOptions
