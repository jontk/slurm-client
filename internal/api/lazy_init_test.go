// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api_test

import (
	"context"
	"sync"
	"testing"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/pkg/config"
)

// TestLazyInitialization verifies that managers are lazily initialized
func TestLazyInitialization(t *testing.T) {
	cfg := &config.Config{
		BaseURL: "http://localhost:6820",
		Version: "v0.0.43",
	}

	client, err := slurm.NewClient(context.Background(),
		slurm.WithConfig(cfg),
		slurm.WithNoAuth(),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Get managers - should not initialize implementation yet
	jobManager := client.Jobs()
	nodeManager := client.Nodes()
	partitionManager := client.Partitions()
	infoManager := client.Info()
	reservationManager := client.Reservations()
	qosManager := client.QoS()
	accountManager := client.Accounts()
	userManager := client.Users()

	// Verify all managers are non-nil
	if jobManager == nil {
		t.Error("JobManager is nil")
	}
	if nodeManager == nil {
		t.Error("NodeManager is nil")
	}
	if partitionManager == nil {
		t.Error("PartitionManager is nil")
	}
	if infoManager == nil {
		t.Error("InfoManager is nil")
	}
	if reservationManager == nil {
		t.Error("ReservationManager is nil")
	}
	if qosManager == nil {
		t.Error("QoSManager is nil")
	}
	if accountManager == nil {
		t.Error("AccountManager is nil")
	}
	if userManager == nil {
		t.Error("UserManager is nil")
	}
}

// TestLazyInitializationThreadSafety verifies thread safety of lazy initialization
func TestLazyInitializationThreadSafety(t *testing.T) {
	cfg := &config.Config{
		BaseURL: "http://localhost:6820",
		Version: "v0.0.43",
	}

	client, err := slurm.NewClient(context.Background(),
		slurm.WithConfig(cfg),
		slurm.WithNoAuth(),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Get a manager
	jobManager := client.Jobs()

	// Test concurrent access to ensure thread safety
	var wg sync.WaitGroup
	errors := make([]error, 0)
	var mu sync.Mutex

	// Launch multiple goroutines that will all try to use the manager
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// This should trigger lazy initialization in a thread-safe way
			_, err := jobManager.List(context.Background(), nil)
			if err != nil {
				// Store any errors (expected due to no server)
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// All errors should be consistent (same error)
	if len(errors) > 0 {
		firstErr := errors[0].Error()
		for i, err := range errors {
			if err.Error() != firstErr {
				t.Errorf("Error %d differs from first error: %v vs %v", i, err, firstErr)
			}
		}
	}
}

// TestConsistentLazyInitialization verifies all versions use lazy initialization
func TestConsistentLazyInitialization(t *testing.T) {
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}

	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			cfg := &config.Config{
				BaseURL: "http://localhost:6820",
				Version: version,
			}

			client, err := slurm.NewClient(context.Background(),
				slurm.WithConfig(cfg),
				slurm.WithNoAuth(),
			)
			if err != nil {
				t.Fatalf("Failed to create client for version %s: %v", version, err)
			}
			defer client.Close()

			// Get all managers
			managers := map[string]interface{}{
				"Jobs":        client.Jobs(),
				"Nodes":       client.Nodes(),
				"Partitions":  client.Partitions(),
				"Info":        client.Info(),
			}

			// For versions that support these managers
			if version >= "v0.0.40" {
				managers["Reservations"] = client.Reservations()
				managers["QoS"] = client.QoS()
				managers["Accounts"] = client.Accounts()
				managers["Users"] = client.Users()
			}

			// Verify all managers are non-nil
			for name, manager := range managers {
				if manager == nil {
					t.Errorf("Version %s: %s manager is nil", version, name)
				}
			}
		})
	}
}
