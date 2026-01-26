// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"github.com/jontk/slurm-client/internal/adapters/common"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// Adapter implements the VersionAdapter interface for API version v0.0.43
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
	wcKeyAdapter       *WCKeyAdapter
	standaloneAdapter  *StandaloneAdapter
	clusterAdapter     *ClusterAdapter
	infoAdapter        *InfoAdapter
}

// NewAdapter creates a new v0.0.43 adapter
func NewAdapter(client *api.ClientWithResponses) *Adapter {
	return &Adapter{
		version:            "v0.0.43",
		client:             client,
		qosAdapter:         NewQoSAdapter(client),
		jobAdapter:         NewJobAdapter(client),
		partitionAdapter:   NewPartitionAdapter(client),
		nodeAdapter:        NewNodeAdapter(client),
		accountAdapter:     NewAccountAdapter(client),
		userAdapter:        NewUserAdapter(client),
		reservationAdapter: NewReservationAdapter(client),
		associationAdapter: NewAssociationAdapter(client),
		wcKeyAdapter:       NewWCKeyAdapter(client),
		standaloneAdapter:  NewStandaloneAdapter(client),
		clusterAdapter:     NewClusterAdapter(client),
		infoAdapter:        NewInfoAdapter(client),
	}
}

// GetVersion returns the API version this adapter supports
func (a *Adapter) GetVersion() string {
	return a.version
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

// GetWCKeyManager returns the WCKey adapter for this version
func (a *Adapter) GetWCKeyManager() common.WCKeyAdapter {
	return a.wcKeyAdapter
}

// GetStandaloneManager returns the Standalone adapter for this version
func (a *Adapter) GetStandaloneManager() common.StandaloneAdapter {
	return a.standaloneAdapter
}

// GetClusterManager returns the Cluster adapter for this version
func (a *Adapter) GetClusterManager() common.ClusterAdapter {
	return a.clusterAdapter
}

// GetInfoManager returns the Info adapter for this version
func (a *Adapter) GetInfoManager() common.InfoAdapter {
	return a.infoAdapter
}
