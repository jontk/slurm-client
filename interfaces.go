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

// QoSManager provides version-agnostic QoS operations
type QoSManager = interfaces.QoSManager

// AccountManager provides version-agnostic account operations
type AccountManager = interfaces.AccountManager

// UserManager provides version-agnostic user operations
type UserManager = interfaces.UserManager

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

type QoS = interfaces.QoS
type QoSList = interfaces.QoSList
type QoSCreate = interfaces.QoSCreate
type QoSCreateResponse = interfaces.QoSCreateResponse
type QoSUpdate = interfaces.QoSUpdate

type Account = interfaces.Account
type AccountList = interfaces.AccountList
type AccountCreate = interfaces.AccountCreate
type AccountCreateResponse = interfaces.AccountCreateResponse
type AccountUpdate = interfaces.AccountUpdate
type AccountQuota = interfaces.AccountQuota
type AccountUsage = interfaces.AccountUsage
type AccountHierarchy = interfaces.AccountHierarchy

type User = interfaces.User
type UserList = interfaces.UserList
type UserAccount = interfaces.UserAccount
type UserAssociation = interfaces.UserAssociation
type UserQuota = interfaces.UserQuota
type UserAccountQuota = interfaces.UserAccountQuota
type UserUsage = interfaces.UserUsage
type AccountUsageStats = interfaces.AccountUsageStats
type UserFairShare = interfaces.UserFairShare
type FairShareNode = interfaces.FairShareNode
type JobPriorityFactors = interfaces.JobPriorityFactors
type PriorityWeights = interfaces.PriorityWeights
type JobPriorityInfo = interfaces.JobPriorityInfo
type AssociationUsage = interfaces.AssociationUsage
type QoSLimits = interfaces.QoSLimits
type UserAccountAssociation = interfaces.UserAccountAssociation
type UserAccessValidation = interfaces.UserAccessValidation
type AccountFairShare = interfaces.AccountFairShare
type FairShareHierarchy = interfaces.FairShareHierarchy

// List options
type ListJobsOptions = interfaces.ListJobsOptions
type ListNodesOptions = interfaces.ListNodesOptions
type ListPartitionsOptions = interfaces.ListPartitionsOptions
type ListReservationsOptions = interfaces.ListReservationsOptions
type ListQoSOptions = interfaces.ListQoSOptions
type ListAccountsOptions = interfaces.ListAccountsOptions
type ListUsersOptions = interfaces.ListUsersOptions
type ListAccountUsersOptions = interfaces.ListAccountUsersOptions
type ListUserAccountAssociationsOptions = interfaces.ListUserAccountAssociationsOptions

// Watch options
type WatchJobsOptions = interfaces.WatchJobsOptions
type WatchNodesOptions = interfaces.WatchNodesOptions
type WatchPartitionsOptions = interfaces.WatchPartitionsOptions
