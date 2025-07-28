package v0_0_43

import (
	"testing"

	"github.com/jontk/slurm-client/internal/adapters/common"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// TestAdapterIntegration verifies that the main adapter correctly implements
// the VersionAdapter interface and exposes all required managers
func TestAdapterIntegration(t *testing.T) {
	// Create a mock client for testing
	client := &api.ClientWithResponses{}
	
	// Create the adapter
	adapter := NewAdapter(client)
	
	// Verify it implements the VersionAdapter interface
	var _ common.VersionAdapter = adapter
	
	t.Run("GetVersion", func(t *testing.T) {
		version := adapter.GetVersion()
		if version != "v0.0.43" {
			t.Errorf("Expected version 'v0.0.43', got '%s'", version)
		}
	})
	
	t.Run("GetQoSManager", func(t *testing.T) {
		qosManager := adapter.GetQoSManager()
		if qosManager == nil {
			t.Error("GetQoSManager returned nil")
		}
		
		// Verify it implements the QoSAdapter interface
		var _ common.QoSAdapter = qosManager
	})
	
	t.Run("GetJobManager", func(t *testing.T) {
		jobManager := adapter.GetJobManager()
		if jobManager == nil {
			t.Error("GetJobManager returned nil")
		}
		
		// Verify it implements the JobAdapter interface
		var _ common.JobAdapter = jobManager
	})
	
	t.Run("GetPartitionManager", func(t *testing.T) {
		partitionManager := adapter.GetPartitionManager()
		if partitionManager == nil {
			t.Error("GetPartitionManager returned nil")
		}
		
		// Verify it implements the PartitionAdapter interface
		var _ common.PartitionAdapter = partitionManager
	})
	
	t.Run("GetNodeManager", func(t *testing.T) {
		nodeManager := adapter.GetNodeManager()
		if nodeManager == nil {
			t.Error("GetNodeManager returned nil")
		}
		
		// Verify it implements the NodeAdapter interface
		var _ common.NodeAdapter = nodeManager
	})
	
	t.Run("GetAccountManager", func(t *testing.T) {
		accountManager := adapter.GetAccountManager()
		if accountManager == nil {
			t.Error("GetAccountManager returned nil")
		}
		
		// Verify it implements the AccountAdapter interface
		var _ common.AccountAdapter = accountManager
	})
	
	t.Run("GetUserManager", func(t *testing.T) {
		userManager := adapter.GetUserManager()
		if userManager == nil {
			t.Error("GetUserManager returned nil")
		}
		
		// Verify it implements the UserAdapter interface
		var _ common.UserAdapter = userManager
	})
	
	t.Run("GetReservationManager", func(t *testing.T) {
		reservationManager := adapter.GetReservationManager()
		if reservationManager == nil {
			t.Error("GetReservationManager returned nil")
		}
		
		// Verify it implements the ReservationAdapter interface
		var _ common.ReservationAdapter = reservationManager
	})
	
	t.Run("GetAssociationManager", func(t *testing.T) {
		associationManager := adapter.GetAssociationManager()
		if associationManager == nil {
			t.Error("GetAssociationManager returned nil")
		}
		
		// Verify it implements the AssociationAdapter interface
		var _ common.AssociationAdapter = associationManager
	})
}

// TestAdapterManagerCreation verifies that all managers are properly initialized
func TestAdapterManagerCreation(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewAdapter(client)
	
	// Test that all managers are not nil
	managers := map[string]interface{}{
		"QoS":         adapter.qosAdapter,
		"Job":         adapter.jobAdapter,
		"Partition":   adapter.partitionAdapter,
		"Node":        adapter.nodeAdapter,
		"Account":     adapter.accountAdapter,
		"User":        adapter.userAdapter,
		"Reservation": adapter.reservationAdapter,
		"Association": adapter.associationAdapter,
	}
	
	for name, manager := range managers {
		if manager == nil {
			t.Errorf("%s manager is nil", name)
		}
	}
	
	// Verify adapter properties
	if adapter.version != "v0.0.43" {
		t.Errorf("Expected version 'v0.0.43', got '%s'", adapter.version)
	}
	
	if adapter.client != client {
		t.Error("Adapter client is not properly set")
	}
}