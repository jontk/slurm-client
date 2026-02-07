// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/internal/adapters/common"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
)

// Adapter implements the VersionAdapter interface for API version v0.0.41
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
	infoAdapter        *InfoAdapter
	clusterAdapter     *ClusterAdapter
	wckeyAdapter       *WCKeyAdapter
}

// NewAdapter creates a new v0.0.41 adapter
func NewAdapter(client *api.ClientWithResponses) *Adapter {
	return &Adapter{
		version:            "v0.0.41",
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
		infoAdapter:        NewInfoAdapter(client),
		clusterAdapter:     NewClusterAdapter(client),
		wckeyAdapter:       NewWCKeyAdapter(client),
	}
}

// GetVersion returns the API version this adapter supports
func (a *Adapter) GetVersion() string {
	return a.version
}

// GetCapabilities returns the features supported by this API version
func (a *Adapter) GetCapabilities() types.ClientCapabilities {
	return types.ClientCapabilities{
		Version: "v0.0.41",
		// Resource Manager Support - read operations
		SupportsJobs:         true,
		SupportsNodes:        true,
		SupportsPartitions:   true,
		SupportsReservations: true,
		// Database Manager Support - read operations
		SupportsAccounts:     true,
		SupportsUsers:        true,
		SupportsQoS:          true,
		SupportsClusters:     true, // Read/Delete supported
		SupportsAssociations: true,
		SupportsWCKeys:       true, // Read/Delete supported
		// Write Operations Support - very limited
		SupportsJobSubmit:        true, // Implemented using JSON marshaling workaround
		SupportsJobUpdate:        false,
		SupportsJobCancel:        true, // Delete operations work
		SupportsNodeUpdate:       false,
		SupportsPartitionWrite:   false,
		SupportsReservationWrite: false,
		// Database Write Operations - mostly not supported
		SupportsAccountWrite:     true, // Delete supported
		SupportsUserWrite:        false,
		SupportsQoSWrite:         false,
		SupportsClusterWrite:     false,
		SupportsAssociationWrite: false,
		SupportsWCKeyWrite:       false,
		// Advanced Features
		SupportsTRES:        true,
		SupportsInstances:   true,
		SupportsReconfigure: true,
		SupportsDiagnostics: true,
		SupportsShares:      true,
		SupportsLicenses:    true,
		// Extended Features - not implemented in adapter pattern
		SupportsJobSteps:       false, // Steps() requires direct API access
		SupportsJobWatch:       false, // Watch() returns "not implemented" error in v0.0.41
		SupportsNodeWatch:      false, // Watch() returns "not implemented" error in v0.0.41
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

// GetClusterManager returns the Cluster adapter for this version
func (a *Adapter) GetClusterManager() common.ClusterAdapter {
	return a.clusterAdapter
}

// GetWCKeyManager returns the WCKey adapter for this version
func (a *Adapter) GetWCKeyManager() common.WCKeyAdapter {
	return a.wckeyAdapter
}

// GetInfoManager returns the Info adapter for this version
func (a *Adapter) GetInfoManager() common.InfoAdapter {
	return a.infoAdapter
}

