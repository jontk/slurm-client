// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_42

import (
	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/internal/adapters/common"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_42"
)

// Adapter implements the VersionAdapter interface for API version v0.0.42
type Adapter struct {
	version            string
	client             *api.ClientWithResponses
	qosAdapter         *QoSAdapter
	jobAdapter         *JobAdapter
	partitionAdapter   *PartitionAdapter
	nodeAdapter        *NodeAdapter
	accountAdapter     *AccountAdapter
	userAdapter        *UserAdapter
	reservationAdapter *ReservationAdapter
	associationAdapter *AssociationAdapter
	standaloneAdapter  *StandaloneAdapter
	wckeyAdapter       *WCKeyAdapter
	clusterAdapter     *ClusterAdapter
	infoAdapter        *InfoAdapter
}

// NewAdapter creates a new v0.0.42 adapter
func NewAdapter(client *api.ClientWithResponses) *Adapter {
	return &Adapter{
		version:            "v0.0.42",
		client:             client,
		qosAdapter:         NewQoSAdapter(client),
		jobAdapter:         NewJobAdapter(client),
		partitionAdapter:   NewPartitionAdapter(client),
		nodeAdapter:        NewNodeAdapter(client),
		accountAdapter:     NewAccountAdapter(client),
		userAdapter:        NewUserAdapter(client),
		reservationAdapter: NewReservationAdapter(client),
		associationAdapter: NewAssociationAdapter(client),
		standaloneAdapter:  NewStandaloneAdapter(client),
		wckeyAdapter:       NewWCKeyAdapter(client),
		clusterAdapter:     NewClusterAdapter(client),
		infoAdapter:        NewInfoAdapter(client),
	}
}

// GetVersion returns the API version this adapter supports
func (a *Adapter) GetVersion() string {
	return a.version
}

// GetCapabilities returns the features supported by this API version
func (a *Adapter) GetCapabilities() types.ClientCapabilities {
	return types.ClientCapabilities{
		Version: "v0.0.42",
		// Resource Manager Support - full read support
		SupportsJobs:         true,
		SupportsNodes:        true,
		SupportsPartitions:   true,
		SupportsReservations: true,
		// Database Manager Support - full support
		SupportsAccounts:     true,
		SupportsUsers:        true,
		SupportsQoS:          true,
		SupportsClusters:     true,
		SupportsAssociations: true,
		SupportsWCKeys:       true,
		// Write Operations Support - enhanced support
		SupportsJobSubmit:        true,
		SupportsJobUpdate:        true,
		SupportsJobCancel:        true,
		SupportsNodeUpdate:       true,
		SupportsPartitionWrite:   false, // Still limited
		SupportsReservationWrite: false, // Still limited
		// Database Write Operations - enhanced support
		SupportsAccountWrite:     true,
		SupportsUserWrite:        true,
		SupportsQoSWrite:         true,
		SupportsClusterWrite:     true,
		SupportsAssociationWrite: true,
		SupportsWCKeyWrite:       true,
		// Advanced Features - full support
		SupportsTRES:        true,
		SupportsInstances:   true,
		SupportsReconfigure: true,
		SupportsDiagnostics: true,
		SupportsShares:      true,
		SupportsLicenses:    true,
		// Extended Features - not implemented in adapter pattern
		SupportsJobSteps:       false, // Steps() requires direct API access
		SupportsJobWatch:       true,  // Watch() is implemented
		SupportsNodeWatch:      true,  // Watch() is implemented
		SupportsPartitionWatch: false, // Watch() not implemented in adapter
		SupportsAnalytics:      false, // Analytics returns nil
		// Extended Account/User Operations - not implemented in adapter
		SupportsAccountHierarchy: false,
		SupportsAccountQuotas:    false,
		SupportsUserHelpers:      false,
		SupportsFairShare:        false,
		// Cluster Operations - limited in adapter pattern
		SupportsClusterCreate: false,
		SupportsClusterUpdate: false,
		SupportsClusterDelete: false,
		// Bulk Operations
		SupportsAssociationBulkDelete: false,
	}
}

// GetQoSManager returns the QoS adapter for this version
func (a *Adapter) GetQoSManager() common.QoSAdapter {
	return a.qosAdapter
}

// GetJobManager returns the Job adapter for this version
func (a *Adapter) GetJobManager() common.JobAdapter {
	return a.jobAdapter
}

// GetPartitionManager returns the Partition adapter for this version
func (a *Adapter) GetPartitionManager() common.PartitionAdapter {
	return a.partitionAdapter
}

// GetNodeManager returns the Node adapter for this version
func (a *Adapter) GetNodeManager() common.NodeAdapter {
	return a.nodeAdapter
}

// GetAccountManager returns the Account adapter for this version
func (a *Adapter) GetAccountManager() common.AccountAdapter {
	return a.accountAdapter
}

// GetUserManager returns the User adapter for this version
func (a *Adapter) GetUserManager() common.UserAdapter {
	return a.userAdapter
}

// GetReservationManager returns the Reservation adapter for this version
func (a *Adapter) GetReservationManager() common.ReservationAdapter {
	return a.reservationAdapter
}

// GetAssociationManager returns the Association adapter for this version
func (a *Adapter) GetAssociationManager() common.AssociationAdapter {
	return a.associationAdapter
}

// GetStandaloneManager returns the Standalone adapter for this version
func (a *Adapter) GetStandaloneManager() common.StandaloneAdapter {
	return a.standaloneAdapter
}

// GetWCKeyManager returns the WCKey adapter for this version
func (a *Adapter) GetWCKeyManager() common.WCKeyAdapter {
	return a.wckeyAdapter
}

// GetClusterManager returns the Cluster adapter for this version
func (a *Adapter) GetClusterManager() common.ClusterAdapter {
	return a.clusterAdapter
}

// GetInfoManager returns the Info adapter for this version
func (a *Adapter) GetInfoManager() common.InfoAdapter {
	return a.infoAdapter
}
