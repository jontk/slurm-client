// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/jontk/slurm-client/internal/adapters/common"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

func TestNewAdapter(t *testing.T) {
	tests := []struct {
		name   string
		client *api.ClientWithResponses
	}{
		{
			name:   "with nil client",
			client: nil,
		},
		{
			name:   "with valid client",
			client: &api.ClientWithResponses{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewAdapter(tt.client)

			// Verify adapter is created
			assert.NotNil(t, adapter)
			assert.Equal(t, "v0.0.43", adapter.version)
			assert.Equal(t, tt.client, adapter.client)

			// Verify all sub-adapters are initialized
			assert.NotNil(t, adapter.qosAdapter)
			assert.NotNil(t, adapter.jobAdapter)
			assert.NotNil(t, adapter.partitionAdapter)
			assert.NotNil(t, adapter.nodeAdapter)
			assert.NotNil(t, adapter.accountAdapter)
			assert.NotNil(t, adapter.userAdapter)
			assert.NotNil(t, adapter.reservationAdapter)
			assert.NotNil(t, adapter.associationAdapter)
			assert.NotNil(t, adapter.wcKeyAdapter)
			assert.NotNil(t, adapter.standaloneAdapter)
			assert.NotNil(t, adapter.clusterAdapter)
		})
	}
}

func TestAdapter_GetVersion(t *testing.T) {
	adapter := NewAdapter(nil)
	version := adapter.GetVersion()
	assert.Equal(t, "v0.0.43", version)
}

func TestAdapter_GetQoSManager(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewAdapter(client)

	manager := adapter.GetQoSManager()
	assert.NotNil(t, manager)
	assert.IsType(t, &QoSAdapter{}, manager)

	// Verify it implements the interface
	var _ common.QoSAdapter = manager
}

func TestAdapter_GetJobManager(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewAdapter(client)

	manager := adapter.GetJobManager()
	assert.NotNil(t, manager)
	assert.IsType(t, &JobAdapter{}, manager)

	// Verify it implements the interface
	var _ common.JobAdapter = manager
}

func TestAdapter_GetPartitionManager(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewAdapter(client)

	manager := adapter.GetPartitionManager()
	assert.NotNil(t, manager)
	assert.IsType(t, &PartitionAdapter{}, manager)

	// Verify it implements the interface
	var _ common.PartitionAdapter = manager
}

func TestAdapter_GetNodeManager(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewAdapter(client)

	manager := adapter.GetNodeManager()
	assert.NotNil(t, manager)
	assert.IsType(t, &NodeAdapter{}, manager)

	// Verify it implements the interface
	var _ common.NodeAdapter = manager
}

func TestAdapter_GetAccountManager(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewAdapter(client)

	manager := adapter.GetAccountManager()
	assert.NotNil(t, manager)
	assert.IsType(t, &AccountAdapter{}, manager)

	// Verify it implements the interface
	var _ common.AccountAdapter = manager
}

func TestAdapter_GetUserManager(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewAdapter(client)

	manager := adapter.GetUserManager()
	assert.NotNil(t, manager)
	assert.IsType(t, &UserAdapter{}, manager)

	// Verify it implements the interface
	var _ common.UserAdapter = manager
}

func TestAdapter_GetReservationManager(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewAdapter(client)

	manager := adapter.GetReservationManager()
	assert.NotNil(t, manager)
	assert.IsType(t, &ReservationAdapter{}, manager)

	// Verify it implements the interface
	var _ common.ReservationAdapter = manager
}

func TestAdapter_GetAssociationManager(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewAdapter(client)

	manager := adapter.GetAssociationManager()
	assert.NotNil(t, manager)
	assert.IsType(t, &AssociationAdapter{}, manager)

	// Verify it implements the interface
	var _ common.AssociationAdapter = manager
}

func TestAdapter_GetWCKeyManager(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewAdapter(client)

	manager := adapter.GetWCKeyManager()
	assert.NotNil(t, manager)
	assert.IsType(t, &WCKeyAdapter{}, manager)

	// Verify it implements the interface
	var _ common.WCKeyAdapter = manager
}

func TestAdapter_GetStandaloneManager(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewAdapter(client)

	manager := adapter.GetStandaloneManager()
	assert.NotNil(t, manager)
	assert.IsType(t, &StandaloneAdapter{}, manager)

	// Verify it implements the interface
	var _ common.StandaloneAdapter = manager
}

func TestAdapter_GetClusterManager(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewAdapter(client)

	manager := adapter.GetClusterManager()
	assert.NotNil(t, manager)
	assert.IsType(t, &ClusterAdapter{}, manager)

	// Verify it implements the interface
	var _ common.ClusterAdapter = manager
}

func TestAdapter_InterfaceCompliance(t *testing.T) {
	// Verify the adapter implements the VersionAdapter interface
	adapter := NewAdapter(nil)
	var _ common.VersionAdapter = adapter
}

func TestAdapter_ManagerConsistency(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewAdapter(client)

	// Test that calling the same getter multiple times returns the same instance
	qos1 := adapter.GetQoSManager()
	qos2 := adapter.GetQoSManager()
	assert.Same(t, qos1, qos2, "GetQoSManager should return the same instance")

	job1 := adapter.GetJobManager()
	job2 := adapter.GetJobManager()
	assert.Same(t, job1, job2, "GetJobManager should return the same instance")

	partition1 := adapter.GetPartitionManager()
	partition2 := adapter.GetPartitionManager()
	assert.Same(t, partition1, partition2, "GetPartitionManager should return the same instance")

	node1 := adapter.GetNodeManager()
	node2 := adapter.GetNodeManager()
	assert.Same(t, node1, node2, "GetNodeManager should return the same instance")

	account1 := adapter.GetAccountManager()
	account2 := adapter.GetAccountManager()
	assert.Same(t, account1, account2, "GetAccountManager should return the same instance")

	user1 := adapter.GetUserManager()
	user2 := adapter.GetUserManager()
	assert.Same(t, user1, user2, "GetUserManager should return the same instance")

	reservation1 := adapter.GetReservationManager()
	reservation2 := adapter.GetReservationManager()
	assert.Same(t, reservation1, reservation2, "GetReservationManager should return the same instance")

	association1 := adapter.GetAssociationManager()
	association2 := adapter.GetAssociationManager()
	assert.Same(t, association1, association2, "GetAssociationManager should return the same instance")

	wckey1 := adapter.GetWCKeyManager()
	wckey2 := adapter.GetWCKeyManager()
	assert.Same(t, wckey1, wckey2, "GetWCKeyManager should return the same instance")

	standalone1 := adapter.GetStandaloneManager()
	standalone2 := adapter.GetStandaloneManager()
	assert.Same(t, standalone1, standalone2, "GetStandaloneManager should return the same instance")

	cluster1 := adapter.GetClusterManager()
	cluster2 := adapter.GetClusterManager()
	assert.Same(t, cluster1, cluster2, "GetClusterManager should return the same instance")
}

func TestAdapter_AllManagersWithClient(t *testing.T) {
	// Test that all managers are properly initialized with the client
	client := &api.ClientWithResponses{}
	adapter := NewAdapter(client)

	// Test each manager is properly initialized and has access to the client
	tests := []struct {
		name    string
		manager interface{}
		getter  func() interface{}
	}{
		{
			name:    "QoS Manager",
			manager: adapter.qosAdapter,
			getter:  func() interface{} { return adapter.GetQoSManager() },
		},
		{
			name:    "Job Manager",
			manager: adapter.jobAdapter,
			getter:  func() interface{} { return adapter.GetJobManager() },
		},
		{
			name:    "Partition Manager",
			manager: adapter.partitionAdapter,
			getter:  func() interface{} { return adapter.GetPartitionManager() },
		},
		{
			name:    "Node Manager",
			manager: adapter.nodeAdapter,
			getter:  func() interface{} { return adapter.GetNodeManager() },
		},
		{
			name:    "Account Manager",
			manager: adapter.accountAdapter,
			getter:  func() interface{} { return adapter.GetAccountManager() },
		},
		{
			name:    "User Manager",
			manager: adapter.userAdapter,
			getter:  func() interface{} { return adapter.GetUserManager() },
		},
		{
			name:    "Reservation Manager",
			manager: adapter.reservationAdapter,
			getter:  func() interface{} { return adapter.GetReservationManager() },
		},
		{
			name:    "Association Manager",
			manager: adapter.associationAdapter,
			getter:  func() interface{} { return adapter.GetAssociationManager() },
		},
		{
			name:    "WCKey Manager",
			manager: adapter.wcKeyAdapter,
			getter:  func() interface{} { return adapter.GetWCKeyManager() },
		},
		{
			name:    "Standalone Manager",
			manager: adapter.standaloneAdapter,
			getter:  func() interface{} { return adapter.GetStandaloneManager() },
		},
		{
			name:    "Cluster Manager",
			manager: adapter.clusterAdapter,
			getter:  func() interface{} { return adapter.GetClusterManager() },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Manager should not be nil
			assert.NotNil(t, tt.manager)
			
			// Getter should return the manager
			assert.NotNil(t, tt.getter())
			
			// Should be the same instance
			assert.Same(t, tt.manager, tt.getter())
		})
	}
}