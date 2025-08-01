// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

// Example: Reservation management (v0.0.43+ only)
func main() {
	// Create configuration
	cfg := config.NewDefault()
	cfg.BaseURL = "https://cluster.example.com:6820"
	
	// Create authentication
	authProvider := auth.NewTokenAuth("your-jwt-token")

	ctx := context.Background()

	// Example 1: Check version support
	fmt.Println("=== Version Support Check ===")
	checkReservationSupport(ctx, cfg, authProvider)

	// Example 2: List reservations
	fmt.Println("\n=== List Reservations ===")
	listReservations(ctx, cfg, authProvider)

	// Example 3: Create a reservation
	fmt.Println("\n=== Create Reservation ===")
	createReservation(ctx, cfg, authProvider)

	// Example 4: Update a reservation
	fmt.Println("\n=== Update Reservation ===")
	updateReservation(ctx, cfg, authProvider)

	// Example 5: Complex reservation scenarios
	fmt.Println("\n=== Complex Reservation Scenarios ===")
	complexReservationScenarios(ctx, cfg, authProvider)
}

// checkReservationSupport checks if the cluster supports reservations
func checkReservationSupport(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	// Try different versions
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}
	
	for _, version := range versions {
		client, err := slurm.NewClientWithVersion(ctx, version,
			slurm.WithConfig(cfg),
			slurm.WithAuth(auth),
		)
		if err != nil {
			log.Printf("Failed to create %s client: %v", version, err)
			continue
		}
		defer client.Close()
		
		// Check if reservations are supported
		if client.Reservations() == nil {
			fmt.Printf("%s: Reservations NOT supported\n", version)
		} else {
			fmt.Printf("%s: Reservations supported ✓\n", version)
		}
	}
	
	// Also test auto-detection
	fmt.Println("\nAuto-detection:")
	client, err := slurm.NewClient(ctx,
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
	)
	if err != nil {
		log.Printf("Failed to create auto-detect client: %v", err)
		return
	}
	defer client.Close()
	
	fmt.Printf("Detected version: %s\n", client.Version())
	if client.Reservations() == nil {
		fmt.Println("Reservations NOT supported on this cluster")
		fmt.Println("Note: Reservations require SLURM REST API v0.0.43 or later")
	} else {
		fmt.Println("Reservations are supported!")
	}
}

// listReservations demonstrates listing reservations
func listReservations(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	// Create v0.0.43 client
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
	)
	if err != nil {
		log.Printf("Failed to create v0.0.43 client: %v", err)
		return
	}
	defer client.Close()
	
	if client.Reservations() == nil {
		fmt.Println("Reservations not supported")
		return
	}
	
	// List all reservations
	fmt.Println("Listing all reservations:")
	reservations, err := client.Reservations().List(ctx, nil)
	if err != nil {
		log.Printf("Failed to list reservations: %v", err)
		return
	}
	
	if len(reservations.Reservations) == 0 {
		fmt.Println("No reservations found")
		return
	}
	
	// Display reservations
	for _, res := range reservations.Reservations {
		fmt.Printf("\nReservation: %s\n", res.Name)
		fmt.Printf("  State: %s\n", res.State)
		fmt.Printf("  Start: %s\n", res.StartTime.Format("2006-01-02 15:04"))
		fmt.Printf("  End: %s\n", res.EndTime.Format("2006-01-02 15:04"))
		fmt.Printf("  Duration: %d hours\n", res.Duration/3600)
		fmt.Printf("  Nodes: %d (%v)\n", res.NodeCount, res.Nodes)
		fmt.Printf("  Users: %v\n", res.Users)
		fmt.Printf("  Accounts: %v\n", res.Accounts)
		
		if len(res.Flags) > 0 {
			fmt.Printf("  Flags: %v\n", res.Flags)
		}
	}
	
	// List with filters
	fmt.Println("\nListing reservations for specific users:")
	userReservations, err := client.Reservations().List(ctx, &interfaces.ListReservationsOptions{
		Users: []string{"user1", "user2"},
		States: []string{"ACTIVE"},
	})
	if err != nil {
		log.Printf("Failed to list user reservations: %v", err)
		return
	}
	
	fmt.Printf("Found %d active reservations for specified users\n", 
		len(userReservations.Reservations))
}

// createReservation demonstrates creating a reservation
func createReservation(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
	)
	if err != nil {
		log.Printf("Failed to create v0.0.43 client: %v", err)
		return
	}
	defer client.Close()
	
	if client.Reservations() == nil {
		fmt.Println("Reservations not supported")
		return
	}
	
	// Example 1: Node-based reservation
	fmt.Println("Creating node-based reservation:")
	
	nodeReservation := &interfaces.ReservationCreate{
		Name:      "maintenance-window",
		StartTime: time.Now().Add(24 * time.Hour),
		Duration:  4 * 3600, // 4 hours
		Nodes:     []string{"node001", "node002", "node003"},
		Users:     []string{"admin", "maintenance"},
		Flags:     []string{"MAINT", "IGNORE_JOBS"},
	}
	
	resp, err := client.Reservations().Create(ctx, nodeReservation)
	if err != nil {
		log.Printf("Failed to create node reservation: %v", err)
	} else {
		fmt.Printf("Created reservation: %s\n", resp.ReservationName)
	}
	
	// Example 2: Core-based reservation
	fmt.Println("\nCreating core-based reservation:")
	
	coreReservation := &interfaces.ReservationCreate{
		Name:          "weekly-bigdata",
		StartTime:     getNextMonday(),
		EndTime:       getNextMonday().Add(72 * time.Hour), // 3 days
		CoreCount:     256,
		PartitionName: "compute",
		Accounts:      []string{"bigdata-project"},
		Features:      []string{"highmem", "ssd"},
	}
	
	resp2, err := client.Reservations().Create(ctx, coreReservation)
	if err != nil {
		log.Printf("Failed to create core reservation: %v", err)
	} else {
		fmt.Printf("Created reservation: %s\n", resp2.ReservationName)
	}
	
	// Example 3: License reservation
	fmt.Println("\nCreating license reservation:")
	
	licenseReservation := &interfaces.ReservationCreate{
		Name:      "matlab-training",
		StartTime: time.Now().Add(48 * time.Hour),
		Duration:  8 * 3600, // 8 hours
		NodeCount: 10,
		Licenses: map[string]int{
			"matlab": 50,
			"simulink": 20,
		},
		Users: []string{"instructor", "student1", "student2"},
	}
	
	resp3, err := client.Reservations().Create(ctx, licenseReservation)
	if err != nil {
		log.Printf("Failed to create license reservation: %v", err)
	} else {
		fmt.Printf("Created reservation: %s\n", resp3.ReservationName)
	}
}

// updateReservation demonstrates updating a reservation
func updateReservation(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
	)
	if err != nil {
		log.Printf("Failed to create v0.0.43 client: %v", err)
		return
	}
	defer client.Close()
	
	if client.Reservations() == nil {
		fmt.Println("Reservations not supported")
		return
	}
	
	reservationName := "maintenance-window"
	
	// Get current reservation
	fmt.Printf("Getting reservation %s:\n", reservationName)
	current, err := client.Reservations().Get(ctx, reservationName)
	if err != nil {
		log.Printf("Failed to get reservation: %v", err)
		return
	}
	
	fmt.Printf("Current end time: %s\n", current.EndTime.Format("2006-01-02 15:04"))
	
	// Extend reservation by 2 hours
	newEndTime := current.EndTime.Add(2 * time.Hour)
	update := &interfaces.ReservationUpdate{
		EndTime: &newEndTime,
		Users:   append(current.Users, "newuser"),
	}
	
	fmt.Println("\nExtending reservation by 2 hours and adding user...")
	err = client.Reservations().Update(ctx, reservationName, update)
	if err != nil {
		log.Printf("Failed to update reservation: %v", err)
		return
	}
	
	fmt.Println("Reservation updated successfully")
	
	// Verify update
	updated, err := client.Reservations().Get(ctx, reservationName)
	if err != nil {
		log.Printf("Failed to get updated reservation: %v", err)
		return
	}
	
	fmt.Printf("New end time: %s\n", updated.EndTime.Format("2006-01-02 15:04"))
	fmt.Printf("Users: %v\n", updated.Users)
}

// complexReservationScenarios demonstrates advanced reservation patterns
func complexReservationScenarios(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
	)
	if err != nil {
		log.Printf("Failed to create v0.0.43 client: %v", err)
		return
	}
	defer client.Close()
	
	if client.Reservations() == nil {
		fmt.Println("Reservations not supported")
		return
	}
	
	// Scenario 1: Recurring weekly reservation
	fmt.Println("Scenario 1: Setting up recurring weekly reservations")
	
	for week := 0; week < 4; week++ {
		startTime := getNextMonday().Add(time.Duration(week*7*24) * time.Hour)
		
		reservation := &interfaces.ReservationCreate{
			Name:      fmt.Sprintf("weekly-gpu-slot-week%d", week+1),
			StartTime: startTime.Add(9 * time.Hour), // 9 AM
			Duration:  8 * 3600, // 8 hours
			Nodes:     []string{"gpu001", "gpu002"},
			Accounts:  []string{"ml-research"},
			Features:  []string{"gpu"},
			Flags:     []string{"DAILY_9_5"}, // Custom flag
		}
		
		resp, err := client.Reservations().Create(ctx, reservation)
		if err != nil {
			log.Printf("Failed to create week %d reservation: %v", week+1, err)
		} else {
			fmt.Printf("  Created: %s for %s\n", 
				resp.ReservationName, 
				startTime.Format("2006-01-02"))
		}
	}
	
	// Scenario 2: Overlapping reservation check
	fmt.Println("\nScenario 2: Checking for reservation conflicts")
	
	// List all reservations
	allReservations, err := client.Reservations().List(ctx, nil)
	if err != nil {
		log.Printf("Failed to list reservations: %v", err)
		return
	}
	
	// Check for overlaps
	for i := 0; i < len(allReservations.Reservations); i++ {
		for j := i + 1; j < len(allReservations.Reservations); j++ {
			res1 := allReservations.Reservations[i]
			res2 := allReservations.Reservations[j]
			
			if hasNodeOverlap(res1.Nodes, res2.Nodes) &&
			   hasTimeOverlap(res1.StartTime, res1.EndTime, res2.StartTime, res2.EndTime) {
				fmt.Printf("  Conflict detected: %s and %s overlap\n", 
					res1.Name, res2.Name)
			}
		}
	}
	
	// Scenario 3: Maintenance window with job drainage
	fmt.Println("\nScenario 3: Creating maintenance window with job drainage")
	
	// First, create a pre-maintenance reservation to prevent new jobs
	preMaintenance := &interfaces.ReservationCreate{
		Name:      "pre-maintenance-drain",
		StartTime: time.Now().Add(23 * time.Hour), // 1 hour before maintenance
		Duration:  1 * 3600,
		Nodes:     []string{"node001", "node002", "node003", "node004"},
		Users:     []string{"root"},
		Flags:     []string{"NO_HOLD_JOBS_AFTER_END"},
	}
	
	_, err = client.Reservations().Create(ctx, preMaintenance)
	if err != nil {
		log.Printf("Failed to create pre-maintenance reservation: %v", err)
	}
	
	// Then create the actual maintenance window
	maintenance := &interfaces.ReservationCreate{
		Name:      "maintenance-full",
		StartTime: time.Now().Add(24 * time.Hour),
		Duration:  4 * 3600,
		Nodes:     []string{"node001", "node002", "node003", "node004"},
		Users:     []string{"root", "admin"},
		Flags:     []string{"MAINT", "IGNORE_JOBS"},
	}
	
	_, err = client.Reservations().Create(ctx, maintenance)
	if err != nil {
		log.Printf("Failed to create maintenance reservation: %v", err)
	} else {
		fmt.Println("  Created maintenance window with job drainage period")
	}
}

// Helper functions

func getNextMonday() time.Time {
	now := time.Now()
	daysUntilMonday := (8 - int(now.Weekday())) % 7
	if daysUntilMonday == 0 {
		daysUntilMonday = 7
	}
	return now.AddDate(0, 0, daysUntilMonday).Truncate(24 * time.Hour)
}

func hasNodeOverlap(nodes1, nodes2 []string) bool {
	nodeMap := make(map[string]bool)
	for _, node := range nodes1 {
		nodeMap[node] = true
	}
	for _, node := range nodes2 {
		if nodeMap[node] {
			return true
		}
	}
	return false
}

func hasTimeOverlap(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && start2.Before(end1)
}
