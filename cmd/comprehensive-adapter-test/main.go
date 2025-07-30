package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/factory"
	"github.com/jontk/slurm-client/internal/interfaces"
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

	versions := []string{}
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
	recordResult(version, "Jobs", "List", err, fmt.Sprintf("Listed %d jobs", len(jobs)))

	// Test Submit Job
	fmt.Println("Testing: Submit Job")
	submitJob := &interfaces.JobSubmission{
		Name:      fmt.Sprintf("adapter-test-%s-%d", version, time.Now().Unix()),
		Partition: "normal",
		Script:    "#!/bin/bash\necho 'Adapter test job'\nhostname\ndate\nsleep 5",
		TimeLimit: 1,
		Nodes:     1,
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
		Since: time.Now(),
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
	recordResult(version, "Nodes", "List", err, fmt.Sprintf("Listed %d nodes", len(nodes)))

	// Test Get Node
	if len(nodes) > 0 && nodes[0].Name != "" {
		fmt.Println("Testing: Get Node")
		node, err := client.Nodes().Get(ctx, nodes[0].Name)
		recordResult(version, "Nodes", "Get", err, fmt.Sprintf("Retrieved node: %s", nodes[0].Name))

		// Test Update Node
		fmt.Println("Testing: Update Node")
		updateReq := &interfaces.NodeUpdate{
			Comment: stringPtr("Test comment"),
		}
		err = client.Nodes().Update(ctx, nodes[0].Name, updateReq)
		recordResult(version, "Nodes", "Update", err, "Updated node comment")
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
	recordResult(version, "Partitions", "List", err, fmt.Sprintf("Listed %d partitions", len(partitions)))

	// Test Get Partition
	if len(partitions) > 0 && partitions[0].Name != "" {
		fmt.Println("Testing: Get Partition")
		partition, err := client.Partitions().Get(ctx, partitions[0].Name)
		recordResult(version, "Partitions", "Get", err, fmt.Sprintf("Retrieved partition: %s", partitions[0].Name))

		// Test Update Partition
		fmt.Println("Testing: Update Partition")
		updateReq := &interfaces.PartitionUpdate{
			MaxTime: intPtr(120), // 2 hours
		}
		err = client.Partitions().Update(ctx, partitions[0].Name, updateReq)
		recordResult(version, "Partitions", "Update", err, "Updated partition max time")
	}
}

func testAccountEndpoints(client interfaces.SlurmClient, version string) {
	ctx := context.Background()
	fmt.Println("\n--- Testing Account Endpoints ---")

	// Test List Accounts
	fmt.Println("Testing: List Accounts")
	accounts, err := client.Accounts().List(ctx)
	recordResult(version, "Accounts", "List", err, fmt.Sprintf("Listed %d accounts", len(accounts)))

	// Test Create Account
	fmt.Println("Testing: Create Account")
	createReq := &interfaces.Account{
		Name:        fmt.Sprintf("test-account-%d", time.Now().Unix()),
		Description: "Test account for adapter testing",
	}
	err = client.Accounts().Create(ctx, createReq)
	recordResult(version, "Accounts", "Create", err, fmt.Sprintf("Created account: %s", createReq.Name))

	// Test Get Account
	if err == nil {
		fmt.Println("Testing: Get Account")
		account, err := client.Accounts().Get(ctx, createReq.Name)
		recordResult(version, "Accounts", "Get", err, fmt.Sprintf("Retrieved account: %v", account != nil))

		// Test Update Account
		fmt.Println("Testing: Update Account")
		updateReq := &interfaces.Account{
			Name:        createReq.Name,
			Description: "Updated test account",
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
	users, err := client.Users().List(ctx)
	recordResult(version, "Users", "List", err, fmt.Sprintf("Listed %d users", len(users)))

	// Test Get User
	if len(users) > 0 && users[0].Name != "" {
		fmt.Println("Testing: Get User")
		user, err := client.Users().Get(ctx, users[0].Name)
		recordResult(version, "Users", "Get", err, fmt.Sprintf("Retrieved user: %s", users[0].Name))
	}

	// Test Create User
	fmt.Println("Testing: Create User")
	createReq := &interfaces.User{
		Name:    fmt.Sprintf("testuser%d", time.Now().Unix()),
		Account: "root", // Assuming root account exists
	}
	err = client.Users().Create(ctx, createReq)
	recordResult(version, "Users", "Create", err, fmt.Sprintf("Created user: %s", createReq.Name))

	// Test Update and Delete if created
	if err == nil {
		// Test Update User
		fmt.Println("Testing: Update User")
		err = client.Users().Update(ctx, createReq.Name, createReq)
		recordResult(version, "Users", "Update", err, "Updated user")

		// Test Delete User
		fmt.Println("Testing: Delete User")
		err = client.Users().Delete(ctx, createReq.Name)
		recordResult(version, "Users", "Delete", err, "Deleted user")
	}
}

func testQoSEndpoints(client interfaces.SlurmClient, version string) {
	ctx := context.Background()
	fmt.Println("\n--- Testing QoS Endpoints ---")

	// Test List QoS
	fmt.Println("Testing: List QoS")
	qosList, err := client.QoS().List(ctx)
	recordResult(version, "QoS", "List", err, fmt.Sprintf("Listed %d QoS", len(qosList)))

	// Test Get QoS
	if len(qosList) > 0 && qosList[0].Name != "" {
		fmt.Println("Testing: Get QoS")
		qos, err := client.QoS().Get(ctx, qosList[0].Name)
		recordResult(version, "QoS", "Get", err, fmt.Sprintf("Retrieved QoS: %s", qosList[0].Name))
	}

	// Test Create QoS
	fmt.Println("Testing: Create QoS")
	createReq := &interfaces.QoS{
		Name:        fmt.Sprintf("test-qos-%d", time.Now().Unix()),
		Description: "Test QoS for adapter testing",
		Priority:    intPtr(10),
	}
	err = client.QoS().Create(ctx, createReq)
	recordResult(version, "QoS", "Create", err, fmt.Sprintf("Created QoS: %s", createReq.Name))

	// Test Update and Delete if created
	if err == nil {
		// Test Update QoS
		fmt.Println("Testing: Update QoS")
		updateReq := &interfaces.QoS{
			Name:        createReq.Name,
			Description: "Updated test QoS",
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
	recordResult(version, "Reservations", "List", err, fmt.Sprintf("Listed %d reservations", len(reservations)))

	// Test Create Reservation
	fmt.Println("Testing: Create Reservation")
	startTime := time.Now().Add(1 * time.Hour)
	endTime := startTime.Add(1 * time.Hour)
	createReq := &interfaces.Reservation{
		Name:      fmt.Sprintf("test-res-%d", time.Now().Unix()),
		StartTime: &startTime,
		EndTime:   &endTime,
		NodeCount: 1,
		Users:     []string{"root"},
	}
	err = client.Reservations().Create(ctx, createReq)
	recordResult(version, "Reservations", "Create", err, fmt.Sprintf("Created reservation: %s", createReq.Name))

	// Test Get, Update and Delete if created
	if err == nil {
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
	associations, err := client.Associations().List(ctx)
	recordResult(version, "Associations", "List", err, fmt.Sprintf("Listed %d associations", len(associations)))

	// Test Create Association
	fmt.Println("Testing: Create Association")
	createReq := &interfaces.Association{
		Account:   "root",
		User:      "root",
		Partition: "normal",
	}
	err = client.Associations().Create(ctx, createReq)
	recordResult(version, "Associations", "Create", err, "Created association")

	// Test Get Association
	if len(associations) > 0 && associations[0].ID != "" {
		fmt.Println("Testing: Get Association")
		assoc, err := client.Associations().Get(ctx, associations[0].ID)
		recordResult(version, "Associations", "Get", err, fmt.Sprintf("Retrieved association: %s", associations[0].ID))
	}

	// Test Update Association
	if len(associations) > 0 && associations[0].ID != "" {
		fmt.Println("Testing: Update Association")
		updateReq := &interfaces.Association{
			ID:      associations[0].ID,
			Account: associations[0].Account,
			User:    associations[0].User,
		}
		err = client.Associations().Update(ctx, associations[0].ID, updateReq)
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