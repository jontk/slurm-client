// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"context"

	"github.com/jontk/slurm-client/internal/common/types"
)

// VersionAdapter is the main adapter interface for a specific API version
type VersionAdapter interface {
	GetVersion() string
	GetQoSManager() QoSAdapter
	GetJobManager() JobAdapter
	GetPartitionManager() PartitionAdapter
	GetNodeManager() NodeAdapter
	GetReservationManager() ReservationAdapter
	GetAccountManager() AccountAdapter
	GetUserManager() UserAdapter
	GetAssociationManager() AssociationAdapter
	GetWCKeyManager() WCKeyAdapter
	GetClusterManager() ClusterAdapter
	GetInfoManager() InfoAdapter

	// Standalone operations (non-CRUD)
	GetStandaloneManager() StandaloneAdapter
}

// QoSAdapter defines the interface for QoS operations across all versions
type QoSAdapter interface {
	// List retrieves a list of QoS with optional filtering
	List(ctx context.Context, opts *types.QoSListOptions) (*types.QoSList, error)

	// Get retrieves a specific QoS by name
	Get(ctx context.Context, qosName string) (*types.QoS, error)

	// Create creates a new QoS
	Create(ctx context.Context, qos *types.QoSCreate) (*types.QoSCreateResponse, error)

	// Update updates an existing QoS
	Update(ctx context.Context, qosName string, update *types.QoSUpdate) error

	// Delete deletes a QoS
	Delete(ctx context.Context, qosName string) error
}

// JobAdapter defines the interface for Job management across versions
type JobAdapter interface {
	List(ctx context.Context, opts *types.JobListOptions) (*types.JobList, error)
	Get(ctx context.Context, jobID int32) (*types.Job, error)
	Submit(ctx context.Context, job *types.JobCreate) (*types.JobSubmitResponse, error)
	Update(ctx context.Context, jobID int32, update *types.JobUpdate) error
	Cancel(ctx context.Context, jobID int32, opts *types.JobCancelRequest) error
	Signal(ctx context.Context, req *types.JobSignalRequest) error
	Hold(ctx context.Context, req *types.JobHoldRequest) error
	Notify(ctx context.Context, req *types.JobNotifyRequest) error
	Requeue(ctx context.Context, jobID int32) error
	Watch(ctx context.Context, opts *types.JobWatchOptions) (<-chan types.JobWatchEvent, error)
	Allocate(ctx context.Context, req *types.JobAllocateRequest) (*types.JobAllocateResponse, error)
}

// PartitionAdapter defines the interface for Partition management across versions
type PartitionAdapter interface {
	List(ctx context.Context, opts *types.PartitionListOptions) (*types.PartitionList, error)
	Get(ctx context.Context, partitionName string) (*types.Partition, error)
	Create(ctx context.Context, partition *types.PartitionCreate) (*types.PartitionCreateResponse, error)
	Update(ctx context.Context, partitionName string, update *types.PartitionUpdate) error
	Delete(ctx context.Context, partitionName string) error
}

// NodeAdapter defines the interface for Node management across versions
type NodeAdapter interface {
	List(ctx context.Context, opts *types.NodeListOptions) (*types.NodeList, error)
	Get(ctx context.Context, nodeName string) (*types.Node, error)
	Update(ctx context.Context, nodeName string, update *types.NodeUpdate) error
	Delete(ctx context.Context, nodeName string) error
	Drain(ctx context.Context, nodeName string, reason string) error
	Resume(ctx context.Context, nodeName string) error
	Watch(ctx context.Context, opts *types.NodeWatchOptions) (<-chan types.NodeWatchEvent, error)
}

// AccountAdapter defines the interface for Account management across versions
type AccountAdapter interface {
	List(ctx context.Context, opts *types.AccountListOptions) (*types.AccountList, error)
	Get(ctx context.Context, accountName string) (*types.Account, error)
	Create(ctx context.Context, account *types.AccountCreate) (*types.AccountCreateResponse, error)
	Update(ctx context.Context, accountName string, update *types.AccountUpdate) error
	Delete(ctx context.Context, accountName string) error
	CreateAssociation(ctx context.Context, req *types.AccountAssociationRequest) (*types.AssociationCreateResponse, error)
}

// UserAdapter defines the interface for User management across versions
type UserAdapter interface {
	List(ctx context.Context, opts *types.UserListOptions) (*types.UserList, error)
	Get(ctx context.Context, userName string) (*types.User, error)
	Create(ctx context.Context, user *types.UserCreate) (*types.UserCreateResponse, error)
	Update(ctx context.Context, userName string, update *types.UserUpdate) error
	Delete(ctx context.Context, userName string) error
	CreateAssociation(ctx context.Context, req *types.UserAssociationRequest) (*types.AssociationCreateResponse, error)
}

// ReservationAdapter defines the interface for Reservation management across versions
type ReservationAdapter interface {
	List(ctx context.Context, opts *types.ReservationListOptions) (*types.ReservationList, error)
	Get(ctx context.Context, reservationName string) (*types.Reservation, error)
	Create(ctx context.Context, reservation *types.ReservationCreate) (*types.ReservationCreateResponse, error)
	Update(ctx context.Context, reservationName string, update *types.ReservationUpdate) error
	Delete(ctx context.Context, reservationName string) error
}

// AssociationAdapter defines the interface for Association management across versions
type AssociationAdapter interface {
	List(ctx context.Context, opts *types.AssociationListOptions) (*types.AssociationList, error)
	Get(ctx context.Context, associationID string) (*types.Association, error)
	Create(ctx context.Context, association *types.AssociationCreate) (*types.AssociationCreateResponse, error)
	Update(ctx context.Context, associationID string, update *types.AssociationUpdate) error
	Delete(ctx context.Context, associationID string) error
}

// StandaloneAdapter defines the interface for standalone operations (non-CRUD)
type StandaloneAdapter interface {
	// GetLicenses retrieves license information
	GetLicenses(ctx context.Context) (*types.LicenseList, error)

	// GetShares retrieves fairshare information with optional filtering
	GetShares(ctx context.Context, opts *types.GetSharesOptions) (*types.SharesList, error)

	// GetConfig retrieves SLURM configuration
	GetConfig(ctx context.Context) (*types.Config, error)

	// GetDiagnostics retrieves SLURM diagnostics information
	GetDiagnostics(ctx context.Context) (*types.Diagnostics, error)

	// GetDBDiagnostics retrieves SLURM database diagnostics information
	GetDBDiagnostics(ctx context.Context) (*types.Diagnostics, error)

	// GetInstance retrieves a specific database instance
	GetInstance(ctx context.Context, opts *types.GetInstanceOptions) (*types.Instance, error)

	// GetInstances retrieves multiple database instances with filtering
	GetInstances(ctx context.Context, opts *types.GetInstancesOptions) (*types.InstanceList, error)

	// GetTRES retrieves all TRES (Trackable RESources)
	GetTRES(ctx context.Context) (*types.TRESList, error)

	// CreateTRES creates a new TRES entry
	CreateTRES(ctx context.Context, req *types.CreateTRESRequest) (*types.TRES, error)

	// Reconfigure triggers a SLURM reconfiguration
	Reconfigure(ctx context.Context) (*types.ReconfigureResponse, error)

	// PingDatabase pings the SLURM database for health checks
	PingDatabase(ctx context.Context) (*types.PingResponse, error)
}

// WCKeyAdapter defines the interface for WCKey management across versions
type WCKeyAdapter interface {
	List(ctx context.Context, opts *types.WCKeyListOptions) (*types.WCKeyList, error)
	Get(ctx context.Context, wcKeyID string) (*types.WCKey, error)
	Create(ctx context.Context, wckey *types.WCKeyCreate) (*types.WCKeyCreateResponse, error)
	Delete(ctx context.Context, wcKeyID string) error
}

// ClusterAdapter defines the interface for Cluster management across versions
type ClusterAdapter interface {
	List(ctx context.Context, opts *types.ClusterListOptions) (*types.ClusterList, error)
	Get(ctx context.Context, clusterName string) (*types.Cluster, error)
	Create(ctx context.Context, cluster *types.ClusterCreate) (*types.ClusterCreateResponse, error)
	Delete(ctx context.Context, clusterName string) error
}

// InfoAdapter defines the interface for cluster information operations across versions
type InfoAdapter interface {
	// Get retrieves cluster information
	Get(ctx context.Context) (*types.ClusterInfo, error)

	// Ping tests connectivity to the cluster
	Ping(ctx context.Context) error

	// PingDatabase tests connectivity to the SLURM database
	PingDatabase(ctx context.Context) error

	// Stats retrieves cluster statistics
	Stats(ctx context.Context) (*types.ClusterStats, error)

	// Version retrieves API version information
	Version(ctx context.Context) (*types.APIVersion, error)
}
