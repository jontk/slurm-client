// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"github.com/jontk/slurm-client/internal/adapters/common"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

// Adapter implements the VersionAdapter interface for API version v0.0.42
type Adapter struct {
	version              string
	client               *api.ClientWithResponses
	qosAdapter           *QoSAdapter
	jobAdapter           *JobAdapter
	partitionAdapter     *PartitionAdapter
	nodeAdapter          *NodeAdapter
	accountAdapter       *AccountAdapter
	userAdapter          *UserAdapter
	reservationAdapter   *ReservationAdapter
	associationAdapter   *AssociationAdapter
	standaloneAdapter    *StandaloneAdapter
}

// NewAdapter creates a new v0.0.42 adapter
func NewAdapter(client *api.ClientWithResponses) *Adapter {
	return &Adapter{
		version:              "v0.0.42",
		client:               client,
		qosAdapter:           NewQoSAdapter(client),
		jobAdapter:           NewJobAdapter(client),
		partitionAdapter:     NewPartitionAdapter(client),
		nodeAdapter:          NewNodeAdapter(client),
		accountAdapter:       NewAccountAdapter(client),
		userAdapter:          NewUserAdapter(client),
		reservationAdapter:   NewReservationAdapter(client),
		associationAdapter:   NewAssociationAdapter(client),
		standaloneAdapter:    NewStandaloneAdapter(client),
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

// GetStandaloneManager returns the Standalone adapter for this version
func (a *Adapter) GetStandaloneManager() common.StandaloneAdapter {
	return a.standaloneAdapter
}

// GetWCKeyManager returns nil as WCKey management is not supported in v0.0.42
func (a *Adapter) GetWCKeyManager() common.WCKeyAdapter {
	return nil
}
