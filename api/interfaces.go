// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Package api contains all type definitions and interfaces for the slurm-client SDK.
// This is the single source of truth for the SDK contract.
// The root package re-exports these as type aliases for user convenience.
package api

import (
	"context"
)

// ============================================================================
// Core Client Interface
// ============================================================================

// SlurmClient represents a version-agnostic Slurm REST API client
type SlurmClient interface {
	// Version returns the API version this client supports
	Version() string

	// Capabilities returns the features supported by this client version
	Capabilities() ClientCapabilities

	// Jobs returns the JobManager for this version
	Jobs() JobManager

	// Nodes returns the NodeManager for this version
	Nodes() NodeManager

	// Partitions returns the PartitionManager for this version
	Partitions() PartitionManager

	// Info returns the InfoManager for this version
	Info() InfoManager

	// Reservations returns the ReservationManager for this version (v0.0.43+)
	Reservations() ReservationManager

	// QoS returns the QoSManager for this version (v0.0.43+)
	QoS() QoSManager

	// Accounts returns the AccountManager for this version (v0.0.43+)
	Accounts() AccountManager

	// Users returns the UserManager for this version (v0.0.43+)
	Users() UserManager

	// Clusters returns the ClusterManager for this version (v0.0.43+)
	Clusters() ClusterManager

	// Associations returns the AssociationManager for this version (v0.0.43+)
	Associations() AssociationManager

	// WCKeys returns the WCKeyManager for this version (v0.0.43+)
	WCKeys() WCKeyManager

	// Analytics returns the AnalyticsManager (optional value-added feature)
	// Returns nil if analytics is not implemented in this client version.
	Analytics() AnalyticsManager

	// === Standalone Operations ===

	// GetLicenses retrieves license information
	GetLicenses(ctx context.Context) (*LicenseList, error)

	// GetShares retrieves fairshare information with optional filtering
	GetShares(ctx context.Context, opts *GetSharesOptions) (*SharesList, error)

	// GetConfig retrieves SLURM configuration
	GetConfig(ctx context.Context) (*Config, error)

	// GetDiagnostics retrieves SLURM diagnostics information
	GetDiagnostics(ctx context.Context) (*Diagnostics, error)

	// GetDBDiagnostics retrieves SLURM database diagnostics information
	GetDBDiagnostics(ctx context.Context) (*Diagnostics, error)

	// GetInstance retrieves a specific database instance
	GetInstance(ctx context.Context, opts *GetInstanceOptions) (*Instance, error)

	// GetInstances retrieves multiple database instances with filtering
	GetInstances(ctx context.Context, opts *GetInstancesOptions) (*InstanceList, error)

	// GetTRES retrieves all TRES (Trackable RESources)
	GetTRES(ctx context.Context) (*TRESList, error)

	// CreateTRES creates a new TRES entry
	CreateTRES(ctx context.Context, req *CreateTRESRequest) (*TRES, error)

	// Reconfigure triggers a SLURM reconfiguration
	Reconfigure(ctx context.Context) (*ReconfigureResponse, error)

	// Close closes the client and any resources
	Close() error
}

// ============================================================================
// Job Interfaces
// ============================================================================

// JobReader provides read-only job operations
type JobReader interface {
	List(ctx context.Context, opts *ListJobsOptions) (*JobList, error)
	Get(ctx context.Context, jobID string) (*Job, error)
	// Note: Job steps are available via Job.Steps field from Get() - no separate endpoint exists
}

// JobWriter provides job mutation operations
type JobWriter interface {
	Submit(ctx context.Context, job *JobSubmission) (*JobSubmitResponse, error)
	Update(ctx context.Context, jobID string, update *JobUpdate) error
}

// JobController provides job control operations
type JobController interface {
	Cancel(ctx context.Context, jobID string) error
	Hold(ctx context.Context, jobID string) error
	Release(ctx context.Context, jobID string) error
	Signal(ctx context.Context, jobID string, signal string) error
	Notify(ctx context.Context, jobID string, message string) error
	Requeue(ctx context.Context, jobID string) error
}

// JobWatcher provides real-time job operations
type JobWatcher interface {
	Watch(ctx context.Context, opts *WatchJobsOptions) (<-chan JobEvent, error)
	Allocate(ctx context.Context, req *JobAllocateRequest) (*JobAllocateResponse, error)
}

// JobManager combines all core job operations
type JobManager interface {
	JobReader
	JobWriter
	JobController
	JobWatcher
}

// ============================================================================
// Node Interface
// ============================================================================

type NodeManager interface {
	List(ctx context.Context, opts *ListNodesOptions) (*NodeList, error)
	Get(ctx context.Context, nodeName string) (*Node, error)
	Update(ctx context.Context, nodeName string, update *NodeUpdate) error
	Delete(ctx context.Context, nodeName string) error
	Drain(ctx context.Context, nodeName string, reason string) error
	Resume(ctx context.Context, nodeName string) error
	Watch(ctx context.Context, opts *WatchNodesOptions) (<-chan NodeEvent, error)
}

// ============================================================================
// Partition Interface
// ============================================================================

type PartitionManager interface {
	List(ctx context.Context, opts *ListPartitionsOptions) (*PartitionList, error)
	Get(ctx context.Context, partitionName string) (*Partition, error)
	Create(ctx context.Context, partition *PartitionCreate) (*PartitionCreateResponse, error)
	Update(ctx context.Context, partitionName string, update *PartitionUpdate) error
	Delete(ctx context.Context, partitionName string) error
	Watch(ctx context.Context, opts *WatchPartitionsOptions) (<-chan PartitionEvent, error)
}

// ============================================================================
// Info Interface
// ============================================================================

type InfoManager interface {
	Get(ctx context.Context) (*ClusterInfo, error)
	Ping(ctx context.Context) error
	PingDatabase(ctx context.Context) error
	Stats(ctx context.Context) (*ClusterStats, error)
	Version(ctx context.Context) (*APIVersion, error)
}

// ============================================================================
// Reservation Interface
// ============================================================================

type ReservationManager interface {
	List(ctx context.Context, opts *ListReservationsOptions) (*ReservationList, error)
	Get(ctx context.Context, reservationName string) (*Reservation, error)
	Create(ctx context.Context, reservation *ReservationCreate) (*ReservationCreateResponse, error)
	Update(ctx context.Context, reservationName string, update *ReservationUpdate) error
	Delete(ctx context.Context, reservationName string) error
}

// ============================================================================
// QoS Interface
// ============================================================================

type QoSManager interface {
	List(ctx context.Context, opts *ListQoSOptions) (*QoSList, error)
	Get(ctx context.Context, qosName string) (*QoS, error)
	Create(ctx context.Context, qos *QoSCreate) (*QoSCreateResponse, error)
	Update(ctx context.Context, qosName string, update *QoSUpdate) error
	Delete(ctx context.Context, qosName string) error
}

// ============================================================================
// Account Interface
// ============================================================================

type AccountManager interface {
	List(ctx context.Context, opts *ListAccountsOptions) (*AccountList, error)
	Get(ctx context.Context, accountName string) (*Account, error)
	Create(ctx context.Context, account *AccountCreate) (*AccountCreateResponse, error)
	Update(ctx context.Context, accountName string, update *AccountUpdate) error
	Delete(ctx context.Context, accountName string) error
}

// ============================================================================
// User Interface
// ============================================================================

type UserManager interface {
	List(ctx context.Context, opts *ListUsersOptions) (*UserList, error)
	Get(ctx context.Context, userName string) (*User, error)
	Create(ctx context.Context, user *UserCreate) (*UserCreateResponse, error)
	Update(ctx context.Context, userName string, update *UserUpdate) error
	Delete(ctx context.Context, userName string) error
}

// ============================================================================
// Cluster Interface
// ============================================================================

type ClusterManager interface {
	List(ctx context.Context, opts *ListClustersOptions) (*ClusterList, error)
	Get(ctx context.Context, clusterName string) (*Cluster, error)
	Create(ctx context.Context, cluster *ClusterCreate) (*ClusterCreateResponse, error)
	// Note: SLURM REST API does not support cluster updates - clusters can only be created or deleted
	Delete(ctx context.Context, clusterName string) error
}

// ============================================================================
// Association Interface
// ============================================================================

type AssociationManager interface {
	List(ctx context.Context, opts *ListAssociationsOptions) (*AssociationList, error)
	Get(ctx context.Context, associationID string) (*Association, error)
	Create(ctx context.Context, associations []*AssociationCreate) (*AssociationCreateResponse, error)
	Update(ctx context.Context, associations []*AssociationUpdate) error
	Delete(ctx context.Context, associationID string) error
}

// ============================================================================
// WCKey Interface
// ============================================================================

type WCKeyManager interface {
	List(ctx context.Context, opts *WCKeyListOptions) (*WCKeyList, error)
	Get(ctx context.Context, wckeyName, user, cluster string) (*WCKey, error)
	Create(ctx context.Context, wckey *WCKeyCreate) (*WCKeyCreateResponse, error)
	// Note: WCKey updates are not supported by the SLURM REST API
	Delete(ctx context.Context, wckeyID string) error
}

// ============================================================================
// Analytics Interface (Optional)
// ============================================================================

// AnalyticsManager provides advanced performance analytics.
// NOTE: This is NOT part of the Slurm REST API - it provides computed insights.
// Returns nil from SlurmClient.Analytics() if not implemented.
type AnalyticsManager interface {
	GetJobUtilization(ctx context.Context, jobID string) (*JobUtilization, error)
	GetJobEfficiency(ctx context.Context, jobID string) (*ResourceUtilization, error)
	GetJobPerformance(ctx context.Context, jobID string) (*JobPerformance, error)
	GetJobLiveMetrics(ctx context.Context, jobID string) (*JobLiveMetrics, error)
	WatchJobMetrics(ctx context.Context, jobID string, opts *WatchMetricsOptions) (<-chan JobMetricsEvent, error)
	GetJobResourceTrends(ctx context.Context, jobID string, opts *ResourceTrendsOptions) (*JobResourceTrends, error)
	GetJobStepDetails(ctx context.Context, jobID string, stepID string) (*JobStepDetails, error)
	GetJobStepUtilization(ctx context.Context, jobID string, stepID string) (*JobStepUtilization, error)
	ListJobStepsWithMetrics(ctx context.Context, jobID string, opts *ListJobStepsOptions) (*JobStepMetricsList, error)
	GetJobStepsFromAccounting(ctx context.Context, jobID string, opts *AccountingQueryOptions) (*AccountingJobSteps, error)
	GetStepAccountingData(ctx context.Context, jobID string, stepID string) (*StepAccountingRecord, error)
	GetJobStepAPIData(ctx context.Context, jobID string, stepID string) (*JobStepAPIData, error)
	ListJobStepsFromSacct(ctx context.Context, jobID string, opts *SacctQueryOptions) (*SacctJobStepData, error)
	GetJobCPUAnalytics(ctx context.Context, jobID string) (*CPUAnalytics, error)
	GetJobMemoryAnalytics(ctx context.Context, jobID string) (*MemoryAnalytics, error)
	GetJobIOAnalytics(ctx context.Context, jobID string) (*IOAnalytics, error)
	GetJobComprehensiveAnalytics(ctx context.Context, jobID string) (*JobComprehensiveAnalytics, error)
	GetJobPerformanceHistory(ctx context.Context, jobID string, opts *PerformanceHistoryOptions) (*JobPerformanceHistory, error)
	GetPerformanceTrends(ctx context.Context, opts *TrendAnalysisOptions) (*PerformanceTrends, error)
	GetUserEfficiencyTrends(ctx context.Context, userID string, opts *EfficiencyTrendOptions) (*UserEfficiencyTrends, error)
	AnalyzeBatchJobs(ctx context.Context, jobIDs []string, opts *BatchAnalysisOptions) (*BatchJobAnalysis, error)
	GetWorkflowPerformance(ctx context.Context, workflowID string, opts *WorkflowAnalysisOptions) (*WorkflowPerformance, error)
	GenerateEfficiencyReport(ctx context.Context, opts *ReportOptions) (*EfficiencyReport, error)
}
