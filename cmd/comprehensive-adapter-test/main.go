// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/internal/factory"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

type TestResult struct {
	Version  string
	Endpoint string
	Method   string
	Success  bool
	Error    string
	Details  string
}

var results []TestResult

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <version|all>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Supported versions: v0.0.40, v0.0.41, v0.0.42, v0.0.43, all\n")
		os.Exit(1)
	}

	// Get JWT token from environment
	jwtToken := os.Getenv("SLURM_JWT")
	if jwtToken == "" {
		log.Fatal("SLURM_JWT environment variable is required")
	}

	var versions []string
	if os.Args[1] == "all" {
		versions = []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}
	} else {
		versions = []string{os.Args[1]}
	}

	// Test each version
	for _, version := range versions {
		fmt.Printf("\n========================================\n")
		fmt.Printf("Testing API Version: %s\n", version)
		fmt.Printf("========================================\n")

		client, err := createClient(version, jwtToken)
		if err != nil {
			log.Printf("Failed to create client for %s: %v", version, err)
			continue
		}

		// Test all endpoints
		testJobEndpoints(client, version)
		testNodeEndpoints(client, version)
		testPartitionEndpoints(client, version)
		testAccountEndpoints(client, version)
		testUserEndpoints(client, version)
		testQoSEndpoints(client, version)
		testReservationEndpoints(client, version)
		testAssociationEndpoints(client, version)
	}

	// Print summary
	printSummary()
}

func createClient(version, jwtToken string) (interfaces.SlurmClient, error) {
	// Create configuration
	cfg := config.NewDefault()
	cfg.BaseURL = "http://rocky9.ar.jontk.com:6820"
	cfg.Debug = false

	// Create JWT authentication provider
	authProvider := auth.NewTokenAuth(jwtToken)

	// Create factory with adapter option
	clientFactory, err := factory.NewClientFactory(
		factory.WithConfig(cfg),
		factory.WithAuth(authProvider),
		factory.WithBaseURL(cfg.BaseURL),
		factory.WithUseAdapters(true), // Force use of adapters
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create factory: %w", err)
	}

	// Create client with specific version
	client, err := clientFactory.NewClientWithVersion(context.Background(), version)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}

func testJobEndpoints(client interfaces.SlurmClient, version string) {
	ctx := context.Background()
	fmt.Println("\n--- Testing Job Endpoints ---")

	// Test List Jobs
	fmt.Println("Testing: List Jobs")
	listOpts := &interfaces.ListJobsOptions{
		Limit: 10,
	}
	jobs, err := client.Jobs().List(ctx, listOpts)
	var jobDetails string
	if err == nil && jobs != nil {
		jobDetails = fmt.Sprintf("Listed %d jobs", jobs.Total)
	} else {
		jobDetails = "Failed to list jobs"
	}
	recordResult(version, "Jobs", "List", err, jobDetails)

	// Test Submit Job
	fmt.Println("Testing: Submit Job")
	submitJob := &interfaces.JobSubmission{
		Name:       fmt.Sprintf("adapter-test-%s-%d", version, time.Now().Unix()),
		Account:    "root", // Required for SLURM v0.0.43
		Partition:  "normal",
		Script:     "#!/bin/bash\necho 'Adapter test job'\nhostname\ndate\nsleep 5",
		TimeLimit:  1,
		Nodes:      1,
		WorkingDir: "/tmp",
		Environment: map[string]string{
			"PATH": "/usr/bin:/bin",
			"USER": "root",
			"HOME": "/tmp",
		},
	}
	submitResp, err := client.Jobs().Submit(ctx, submitJob)
	var jobID string
	if err == nil && submitResp != nil {
		jobID = submitResp.JobID
	}
	recordResult(version, "Jobs", "Submit", err, fmt.Sprintf("Submitted job ID: %s", jobID))

	// Test Get Job
	if jobID != "" {
		fmt.Println("Testing: Get Job")
		job, err := client.Jobs().Get(ctx, jobID)
		recordResult(version, "Jobs", "Get", err, fmt.Sprintf("Retrieved job: %v", job != nil))

		// Test Update Job
		fmt.Println("Testing: Update Job")
		updateReq := &interfaces.JobUpdate{
			Priority: intPtr(100),
		}
		err = client.Jobs().Update(ctx, jobID, updateReq)
		recordResult(version, "Jobs", "Update", err, "Updated job priority")

		// Test Cancel Job
		fmt.Println("Testing: Cancel Job")
		err = client.Jobs().Cancel(ctx, jobID)
		recordResult(version, "Jobs", "Cancel", err, "Cancelled job")
	}

	// Test Watch Jobs (if supported)
	fmt.Println("Testing: Watch Jobs")
	watchOpts := &interfaces.WatchJobsOptions{
		// No Since field in this interface
	}
	eventChan, err := client.Jobs().Watch(ctx, watchOpts)
	if err == nil && eventChan != nil {
		// Just check if channel is created, don't actually watch
		recordResult(version, "Jobs", "Watch", nil, "Watch channel created")
	} else {
		recordResult(version, "Jobs", "Watch", err, "")
	}
}

func testNodeEndpoints(client interfaces.SlurmClient, version string) {
	ctx := context.Background()
	fmt.Println("\n--- Testing Node Endpoints ---")

	// Test List Nodes
	fmt.Println("Testing: List Nodes")
	listOpts := &interfaces.ListNodesOptions{
		Limit: 10,
	}
	nodes, err := client.Nodes().List(ctx, listOpts)
	var nodeDetails string
	if err == nil && nodes != nil {
		nodeDetails = fmt.Sprintf("Listed %d nodes", nodes.Total)
	} else {
		nodeDetails = "Failed to list nodes"
	}
	recordResult(version, "Nodes", "List", err, nodeDetails)

	// Test Get Node
	if nodes != nil && len(nodes.Nodes) > 0 && nodes.Nodes[0].Name != "" {
		fmt.Println("Testing: Get Node")
		node, err := client.Nodes().Get(ctx, nodes.Nodes[0].Name)
		recordResult(version, "Nodes", "Get", err, fmt.Sprintf("Retrieved node: %s", nodes.Nodes[0].Name))

		// Test Update Node
		fmt.Println("Testing: Update Node")
		updateReq := &interfaces.NodeUpdate{
			Reason: stringPtr("Test reason"),
		}
		err = client.Nodes().Update(ctx, nodes.Nodes[0].Name, updateReq)
		recordResult(version, "Nodes", "Update", err, "Updated node reason")

		// Use the node variable to avoid unused variable error
		_ = node
	}

	// Test Watch Nodes
	fmt.Println("Testing: Watch Nodes")
	watchOpts := &interfaces.WatchNodesOptions{}
	eventChan, err := client.Nodes().Watch(ctx, watchOpts)
	if err == nil && eventChan != nil {
		recordResult(version, "Nodes", "Watch", nil, "Watch channel created")
	} else {
		recordResult(version, "Nodes", "Watch", err, "")
	}
}

func testPartitionEndpoints(client interfaces.SlurmClient, version string) {
	ctx := context.Background()
	fmt.Println("\n--- Testing Partition Endpoints ---")

	// Test List Partitions
	fmt.Println("Testing: List Partitions")
	listOpts := &interfaces.ListPartitionsOptions{}
	partitions, err := client.Partitions().List(ctx, listOpts)
	var partitionDetails string
	if err == nil && partitions != nil {
		partitionDetails = fmt.Sprintf("Listed %d partitions", partitions.Total)
	} else {
		partitionDetails = "Failed to list partitions"
	}
	recordResult(version, "Partitions", "List", err, partitionDetails)

	// Test Get Partition
	if partitions != nil && len(partitions.Partitions) > 0 && partitions.Partitions[0].Name != "" {
		fmt.Println("Testing: Get Partition")
		partition, err := client.Partitions().Get(ctx, partitions.Partitions[0].Name)
		recordResult(version, "Partitions", "Get", err, fmt.Sprintf("Retrieved partition: %s", partitions.Partitions[0].Name))

		// Test Update Partition
		fmt.Println("Testing: Update Partition")
		updateReq := &interfaces.PartitionUpdate{
			MaxTime: intPtr(120), // 2 hours
		}
		err = client.Partitions().Update(ctx, partitions.Partitions[0].Name, updateReq)
		recordResult(version, "Partitions", "Update", err, "Updated partition max time")

		// Use the partition variable to avoid unused variable error
		_ = partition
	}
}

func testAccountEndpoints(client interfaces.SlurmClient, version string) {
	ctx := context.Background()
	fmt.Println("\n--- Testing Account Endpoints ---")

	// Test List Accounts
	fmt.Println("Testing: List Accounts")
	listOpts := &interfaces.ListAccountsOptions{}
	accounts, err := client.Accounts().List(ctx, listOpts)
	var accountDetails string
	if err == nil && accounts != nil {
		accountDetails = fmt.Sprintf("Listed %d accounts", accounts.Total)
	} else {
		accountDetails = "Failed to list accounts"
	}
	recordResult(version, "Accounts", "List", err, accountDetails)

	// Test Create Account
	fmt.Println("Testing: Create Account")
	createReq := &interfaces.AccountCreate{
		Name:        fmt.Sprintf("test-account-%d", time.Now().Unix()),
		Description: "Test account for adapter testing",
	}
	createResp, err := client.Accounts().Create(ctx, createReq)
	recordResult(version, "Accounts", "Create", err, fmt.Sprintf("Created account: %s", createReq.Name))

	// Test Get Account
	if err == nil && createResp != nil {
		fmt.Println("Testing: Get Account")
		account, err := client.Accounts().Get(ctx, createReq.Name)
		recordResult(version, "Accounts", "Get", err, fmt.Sprintf("Retrieved account: %v", account != nil))

		// Test Update Account
		fmt.Println("Testing: Update Account")
		updateReq := &interfaces.AccountUpdate{
			Description: stringPtr("Updated test account"),
		}
		err = client.Accounts().Update(ctx, createReq.Name, updateReq)
		recordResult(version, "Accounts", "Update", err, "Updated account description")

		// Test Delete Account
		fmt.Println("Testing: Delete Account")
		err = client.Accounts().Delete(ctx, createReq.Name)
		recordResult(version, "Accounts", "Delete", err, "Deleted account")
	}
}

func testUserEndpoints(client interfaces.SlurmClient, version string) {
	ctx := context.Background()
	fmt.Println("\n--- Testing User Endpoints ---")

	// Test List Users
	fmt.Println("Testing: List Users")
	listOpts := &interfaces.ListUsersOptions{}
	users, err := client.Users().List(ctx, listOpts)
	var userDetails string
	if err == nil && users != nil {
		userDetails = fmt.Sprintf("Listed %d users", users.Total)
	} else {
		userDetails = "Failed to list users"
	}
	recordResult(version, "Users", "List", err, userDetails)

	// Test Get User
	if users != nil && len(users.Users) > 0 && users.Users[0].Name != "" {
		fmt.Println("Testing: Get User")
		user, err := client.Users().Get(ctx, users.Users[0].Name)
		recordResult(version, "Users", "Get", err, fmt.Sprintf("Retrieved user: %s", users.Users[0].Name))
		_ = user
	}

	// Note: UserManager interface doesn't have Create/Update/Delete methods
	// These operations are typically done through the account management interface
}

func testQoSEndpoints(client interfaces.SlurmClient, version string) {
	ctx := context.Background()
	fmt.Println("\n--- Testing QoS Endpoints ---")

	// Test List QoS
	fmt.Println("Testing: List QoS")
	listOpts := &interfaces.ListQoSOptions{}
	qosList, err := client.QoS().List(ctx, listOpts)
	var qosDetails string
	if err == nil && qosList != nil {
		qosDetails = fmt.Sprintf("Listed %d QoS", qosList.Total)
	} else {
		qosDetails = "Failed to list QoS"
	}
	recordResult(version, "QoS", "List", err, qosDetails)

	// Test Get QoS
	if qosList != nil && len(qosList.QoS) > 0 && qosList.QoS[0].Name != "" {
		fmt.Println("Testing: Get QoS")
		qos, err := client.QoS().Get(ctx, qosList.QoS[0].Name)
		recordResult(version, "QoS", "Get", err, fmt.Sprintf("Retrieved QoS: %s", qosList.QoS[0].Name))
		_ = qos
	}

	// Test Create QoS
	fmt.Println("Testing: Create QoS")
	createReq := &interfaces.QoSCreate{
		Name:        fmt.Sprintf("test-qos-%d", time.Now().Unix()),
		Description: "Test QoS for adapter testing",
		Priority:    10,
	}
	createResp, err := client.QoS().Create(ctx, createReq)
	recordResult(version, "QoS", "Create", err, fmt.Sprintf("Created QoS: %s", createReq.Name))

	// Test Update and Delete if created
	if err == nil && createResp != nil {
		// Test Update QoS
		fmt.Println("Testing: Update QoS")
		updateReq := &interfaces.QoSUpdate{
			Description: stringPtr("Updated test QoS"),
			Priority:    intPtr(20),
		}
		err = client.QoS().Update(ctx, createReq.Name, updateReq)
		recordResult(version, "QoS", "Update", err, "Updated QoS priority")

		// Test Delete QoS
		fmt.Println("Testing: Delete QoS")
		err = client.QoS().Delete(ctx, createReq.Name)
		recordResult(version, "QoS", "Delete", err, "Deleted QoS")
	}
}

func testReservationEndpoints(client interfaces.SlurmClient, version string) {
	ctx := context.Background()
	fmt.Println("\n--- Testing Reservation Endpoints ---")

	// Test List Reservations
	fmt.Println("Testing: List Reservations")
	listOpts := &interfaces.ListReservationsOptions{}
	reservations, err := client.Reservations().List(ctx, listOpts)
	var reservationDetails string
	if err == nil && reservations != nil {
		reservationDetails = fmt.Sprintf("Listed %d reservations", reservations.Total)
	} else {
		reservationDetails = "Failed to list reservations"
	}
	recordResult(version, "Reservations", "List", err, reservationDetails)

	// Test Create Reservation
	fmt.Println("Testing: Create Reservation")
	startTime := time.Now().Add(1 * time.Hour)
	endTime := startTime.Add(1 * time.Hour)
	createReq := &interfaces.ReservationCreate{
		Name:      fmt.Sprintf("test-res-%d", time.Now().Unix()),
		StartTime: startTime,
		EndTime:   endTime,
		NodeCount: 1,
		Users:     []string{"root"},
	}
	createResp, err := client.Reservations().Create(ctx, createReq)
	recordResult(version, "Reservations", "Create", err, fmt.Sprintf("Created reservation: %s", createReq.Name))

	// Test Get, Update and Delete if created
	if err == nil && createResp != nil {
		// Test Get Reservation
		fmt.Println("Testing: Get Reservation")
		reservation, err := client.Reservations().Get(ctx, createReq.Name)
		recordResult(version, "Reservations", "Get", err, fmt.Sprintf("Retrieved reservation: %v", reservation != nil))

		// Test Update Reservation
		fmt.Println("Testing: Update Reservation")
		updateReq := &interfaces.ReservationUpdate{
			NodeCount: intPtr(2),
		}
		err = client.Reservations().Update(ctx, createReq.Name, updateReq)
		recordResult(version, "Reservations", "Update", err, "Updated reservation node count")

		// Test Delete Reservation
		fmt.Println("Testing: Delete Reservation")
		err = client.Reservations().Delete(ctx, createReq.Name)
		recordResult(version, "Reservations", "Delete", err, "Deleted reservation")
	}
}

func testAssociationEndpoints(client interfaces.SlurmClient, version string) {
	ctx := context.Background()
	fmt.Println("\n--- Testing Association Endpoints ---")

	// Test List Associations
	fmt.Println("Testing: List Associations")
	listOpts := &interfaces.ListAssociationsOptions{}
	associations, err := client.Associations().List(ctx, listOpts)
	var associationDetails string
	if err == nil && associations != nil {
		associationDetails = fmt.Sprintf("Listed %d associations", associations.Total)
	} else {
		associationDetails = "Failed to list associations"
	}
	recordResult(version, "Associations", "List", err, associationDetails)

	// Test Create Association
	fmt.Println("Testing: Create Association")
	createReq := &interfaces.AssociationCreate{
		Account:   "root",
		User:      "root",
		Partition: "normal",
	}
	createResp, err := client.Associations().Create(ctx, []*interfaces.AssociationCreate{createReq})
	recordResult(version, "Associations", "Create", err, "Created association")
	_ = createResp

	// Test Get Association
	if associations != nil && len(associations.Associations) > 0 && associations.Associations[0].ID != 0 {
		fmt.Println("Testing: Get Association")
		getOpts := &interfaces.GetAssociationOptions{
			User:    associations.Associations[0].User,
			Account: associations.Associations[0].Account,
			Cluster: associations.Associations[0].Cluster,
		}
		assoc, err := client.Associations().Get(ctx, getOpts)
		recordResult(version, "Associations", "Get", err, fmt.Sprintf("Retrieved association: ID %d", associations.Associations[0].ID))
		_ = assoc
	}

	// Test Update Association
	if associations != nil && len(associations.Associations) > 0 && associations.Associations[0].ID != 0 {
		fmt.Println("Testing: Update Association")
		updateReq := &interfaces.AssociationUpdate{
			Account: associations.Associations[0].Account,
			User:    associations.Associations[0].User,
			Comment: stringPtr("Updated association"),
		}
		err = client.Associations().Update(ctx, []*interfaces.AssociationUpdate{updateReq})
		recordResult(version, "Associations", "Update", err, "Updated association")
	}

	// Test Delete Association
	if err == nil {
		fmt.Println("Testing: Delete Association")
		// Note: We're not actually deleting real associations to avoid breaking the system
		recordResult(version, "Associations", "Delete", fmt.Errorf("skipped to avoid breaking system"), "")
	}
}

func recordResult(version, endpoint, method string, err error, details string) {
	result := TestResult{
		Version:  version,
		Endpoint: endpoint,
		Method:   method,
		Success:  err == nil,
		Details:  details,
	}

	if err != nil {
		result.Error = err.Error()
		fmt.Printf("  ❌ FAILED: %s\n", err.Error())
	} else {
		fmt.Printf("  ✅ SUCCESS: %s\n", details)
	}

	results = append(results, result)
}

func printSummary() {
	fmt.Println("\n\n========================================")
	fmt.Println("TEST SUMMARY")
	fmt.Println("========================================")

	// Group by version
	versionResults := make(map[string][]TestResult)
	for _, r := range results {
		versionResults[r.Version] = append(versionResults[r.Version], r)
	}

	// Print summary for each version
	for version, vResults := range versionResults {
		fmt.Printf("\n%s Results:\n", version)
		fmt.Println(strings.Repeat("-", 50))

		// Group by endpoint
		endpointResults := make(map[string][]TestResult)
		for _, r := range vResults {
			endpointResults[r.Endpoint] = append(endpointResults[r.Endpoint], r)
		}

		// Print results by endpoint
		for endpoint, eResults := range endpointResults {
			successCount := 0
			for _, r := range eResults {
				if r.Success {
					successCount++
				}
			}

			fmt.Printf("  %s: %d/%d passed\n", endpoint, successCount, len(eResults))

			// Show failures
			for _, r := range eResults {
				if !r.Success {
					fmt.Printf("    ❌ %s failed: %s\n", r.Method, r.Error)
				}
			}
		}
	}

	// Overall summary
	fmt.Println("\n" + strings.Repeat("=", 50))
	totalTests := len(results)
	successTests := 0
	for _, r := range results {
		if r.Success {
			successTests++
		}
	}

	fmt.Printf("TOTAL: %d/%d tests passed (%.1f%%)\n",
		successTests, totalTests, float64(successTests)/float64(totalTests)*100)

	// List all unique errors
	fmt.Println("\nUnique Errors Found:")
	errorMap := make(map[string]int)
	for _, r := range results {
		if !r.Success && r.Error != "" {
			errorMap[r.Error]++
		}
	}

	for err, count := range errorMap {
		fmt.Printf("  - %s (occurred %d times)\n", err, count)
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
