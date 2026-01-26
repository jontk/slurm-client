// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

import (
	"context"
	"fmt"

	"github.com/jontk/slurm-client/internal/adapters/common"
	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
	"github.com/jontk/slurm-client/internal/common/types"
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

// GetClusterManager returns the Cluster adapter for this version
func (a *Adapter) GetClusterManager() common.ClusterAdapter {
	// v0.0.41 has different API schema structure (inline anonymous types)
	// TODO: Implement custom v0.0.41-specific cluster adapter
	return &notImplementedClusterAdapter{}
}

// GetWCKeyManager returns the WCKey adapter for this version
func (a *Adapter) GetWCKeyManager() common.WCKeyAdapter {
	// v0.0.41 has different API schema structure (inline anonymous types)
	// TODO: Implement custom v0.0.41-specific wckey adapter
	return &notImplementedWCKeyAdapter{}
}

// GetInfoManager returns the Info adapter for this version
func (a *Adapter) GetInfoManager() common.InfoAdapter {
	return a.infoAdapter
}

// notImplementedClusterAdapter provides stub for v0.0.41 which has different API schema
type notImplementedClusterAdapter struct{}

func (n *notImplementedClusterAdapter) List(_ context.Context, _ *types.ClusterListOptions) (*types.ClusterList, error) {
	return nil, fmt.Errorf("cluster management not yet implemented for v0.0.41 (API schema limitation)")
}

func (n *notImplementedClusterAdapter) Get(_ context.Context, _ string) (*types.Cluster, error) {
	return nil, fmt.Errorf("cluster management not yet implemented for v0.0.41 (API schema limitation)")
}

func (n *notImplementedClusterAdapter) Create(_ context.Context, _ *types.ClusterCreate) (*types.ClusterCreateResponse, error) {
	return nil, fmt.Errorf("cluster management not yet implemented for v0.0.41 (API schema limitation)")
}

func (n *notImplementedClusterAdapter) Delete(_ context.Context, _ string) error {
	return fmt.Errorf("cluster management not yet implemented for v0.0.41 (API schema limitation)")
}

// notImplementedWCKeyAdapter provides stub for v0.0.41 which has different API schema
type notImplementedWCKeyAdapter struct{}

func (n *notImplementedWCKeyAdapter) List(_ context.Context, _ *types.WCKeyListOptions) (*types.WCKeyList, error) {
	return nil, fmt.Errorf("wckey management not yet implemented for v0.0.41 (API schema limitation)")
}

func (n *notImplementedWCKeyAdapter) Get(_ context.Context, _ string) (*types.WCKey, error) {
	return nil, fmt.Errorf("wckey management not yet implemented for v0.0.41 (API schema limitation)")
}

func (n *notImplementedWCKeyAdapter) Create(_ context.Context, _ *types.WCKeyCreate) (*types.WCKeyCreateResponse, error) {
	return nil, fmt.Errorf("wckey management not yet implemented for v0.0.41 (API schema limitation)")
}

func (n *notImplementedWCKeyAdapter) Delete(_ context.Context, _ string) error {
	return fmt.Errorf("wckey management not yet implemented for v0.0.41 (API schema limitation)")
}
