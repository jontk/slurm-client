package main

import (
	"fmt"
	"reflect"

	"github.com/jontk/slurm-client/internal/adapters/common"
	v0_0_43 "github.com/jontk/slurm-client/internal/adapters/v0_0_43"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// verifyIntegration checks that the adapter implements all required interfaces
func verifyIntegration() {
	fmt.Println("🔍 Verifying Adapter Integration...")
	
	// Create a test client
	client := &api.ClientWithResponses{}
	
	// Create the adapter
	adapter := v0_0_43.NewAdapter(client)
	
	// Verify it implements VersionAdapter interface
	var versionAdapter common.VersionAdapter = adapter
	_ = versionAdapter
	
	// Get the type of the adapter to inspect its methods
	adapterType := reflect.TypeOf(adapter)
	
	fmt.Printf("✅ Adapter Type: %s\n", adapterType.String())
	fmt.Printf("✅ Version: %s\n", adapter.GetVersion())
	
	// Check that all required methods exist
	requiredMethods := []string{
		"GetVersion",
		"GetQoSManager", 
		"GetJobManager",
		"GetPartitionManager",
		"GetNodeManager",
		"GetAccountManager",
		"GetUserManager",
		"GetReservationManager",
		"GetAssociationManager",
	}
	
	fmt.Println("\n📋 Checking Required Methods:")
	for _, methodName := range requiredMethods {
		if method, exists := adapterType.MethodByName(methodName); exists {
			fmt.Printf("  ✅ %s: %s\n", methodName, method.Type.String())
		} else {
			fmt.Printf("  ❌ %s: MISSING\n", methodName)
		}
	}
	
	// Check that all managers can be retrieved
	fmt.Println("\n🏭 Checking Manager Retrieval:")
	
	managers := map[string]interface{}{
		"QoS":         adapter.GetQoSManager(),
		"Job":         adapter.GetJobManager(),
		"Partition":   adapter.GetPartitionManager(),
		"Node":        adapter.GetNodeManager(),
		"Account":     adapter.GetAccountManager(),
		"User":        adapter.GetUserManager(),
		"Reservation": adapter.GetReservationManager(),
		"Association": adapter.GetAssociationManager(),
	}
	
	for name, manager := range managers {
		if manager != nil {
			fmt.Printf("  ✅ %s Manager: %s\n", name, reflect.TypeOf(manager).String())
		} else {
			fmt.Printf("  ❌ %s Manager: NIL\n", name)
		}
	}
	
	fmt.Println("\n🎉 Integration verification complete!")
}

func main() {
	verifyIntegration()
}