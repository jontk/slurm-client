// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
	"github.com/jontk/slurm-client/pkg/errors"
)

// RealServerIntegrationTestSuite extends the existing real server tests
// with more comprehensive integration scenarios
type RealServerIntegrationTestSuite struct {
	suite.Suite
	client    slurm.SlurmClient
	serverURL string
	token     string
	version   string
	
	// Test data tracking
	submittedJobs []string
	testStartTime time.Time
}

// SetupSuite initializes the integration test suite
func (suite *RealServerIntegrationTestSuite) SetupSuite() {
	// Check if real server integration testing is enabled
	if os.Getenv("SLURM_REAL_INTEGRATION_TEST") != "true" {
		suite.T().Skip("Real server integration tests disabled. Set SLURM_REAL_INTEGRATION_TEST=true to enable")
	}

	suite.testStartTime = time.Now()

	// Get configuration from environment
	suite.serverURL = os.Getenv("SLURM_SERVER_URL")
	if suite.serverURL == "" {
		suite.serverURL = "http://rocky9:6820"
	}

	suite.version = os.Getenv("SLURM_API_VERSION")
	if suite.version == "" {
		suite.version = "v0.0.43"
	}

	// Get JWT token
	suite.token = os.Getenv("SLURM_JWT_TOKEN")
	if suite.token == "" {
		token, err := fetchJWTTokenViaSSH()
		require.NoError(suite.T(), err, "Failed to fetch JWT token via SSH")
		suite.token = token
	}

	// Create client with enhanced configuration for integration testing
	ctx := context.Background()
	client, err := slurm.NewClientWithVersion(ctx, suite.version,
		slurm.WithBaseURL(suite.serverURL),
		slurm.WithAuth(auth.NewTokenAuth(suite.token)),
		slurm.WithConfig(&config.Config{
			Timeout:              60 * time.Second, // Longer timeout for integration tests
			MaxRetries:           5,                // More retries for reliability
			Debug:                true,
			InsecureSkipVerify:   true,
		}),
	)
	require.NoError(suite.T(), err)
	suite.client = client

	suite.submittedJobs = make([]string, 0)

	suite.T().Logf("Real server integration test suite initialized for %s", suite.version)
}

// TearDownSuite cleans up test resources
func (suite *RealServerIntegrationTestSuite) TearDownSuite() {
	if suite.client == nil {
		return
	}

	ctx := context.Background()

	// Clean up any submitted jobs
	for _, jobID := range suite.submittedJobs {
		err := suite.client.Jobs().Cancel(ctx, jobID)
		if err != nil {
			suite.T().Logf("Failed to cancel job %s: %v", jobID, err)
		} else {
			suite.T().Logf("Cleaned up job %s", jobID)
		}
	}

	suite.client.Close()
	
	duration := time.Since(suite.testStartTime)
	suite.T().Logf("Real server integration tests completed in %v", duration)
}

// TestCompleteClusterDiscovery tests comprehensive cluster discovery
func (suite *RealServerIntegrationTestSuite) TestCompleteClusterDiscovery() {
	ctx := context.Background()

	// Step 1: Basic connectivity
	suite.T().Log("=== Step 1: Testing basic connectivity ===")
	err := suite.client.Info().Ping(ctx)
	suite.Require().NoError(err, "Ping should succeed")

	// Step 2: Get detailed cluster information
	suite.T().Log("=== Step 2: Getting cluster information ===")
	info, err := suite.client.Info().Get(ctx)
	suite.Require().NoError(err)
	suite.NotEmpty(info.ClusterName, "Cluster name should not be empty")
	
	suite.T().Logf("Cluster: %s", info.ClusterName)
	if info.SlurmVersion != "" {
		suite.T().Logf("SLURM Version: %s", info.SlurmVersion)
	}

	// Step 3: Get version information
	suite.T().Log("=== Step 3: Getting version information ===")
	versionInfo, err := suite.client.Info().Version(ctx)
	suite.Require().NoError(err)
	suite.NotEmpty(versionInfo.Version, "Version should not be empty")
	suite.T().Logf("API Version: %s", versionInfo.Version)

	// Step 4: Get cluster statistics
	suite.T().Log("=== Step 4: Getting cluster statistics ===")
	stats, err := suite.client.Info().Stats(ctx)
	suite.Require().NoError(err)
	
	suite.T().Logf("Cluster Statistics:")
	suite.T().Logf("  Total Nodes: %d (Idle: %d, Allocated: %d)", 
		stats.TotalNodes, stats.IdleNodes, stats.AllocatedNodes)
	suite.T().Logf("  Total CPUs: %d (Idle: %d, Allocated: %d)", 
		stats.TotalCPUs, stats.IdleCPUs, stats.AllocatedCPUs)
	suite.T().Logf("  Total Jobs: %d (Running: %d, Pending: %d, Completed: %d)", 
		stats.TotalJobs, stats.RunningJobs, stats.PendingJobs, stats.CompletedJobs)

	// Validate statistics consistency
	suite.GreaterOrEqual(stats.TotalNodes, stats.IdleNodes+stats.AllocatedNodes, 
		"Total nodes should be >= idle + allocated")
	suite.GreaterOrEqual(stats.TotalCPUs, stats.IdleCPUs+stats.AllocatedCPUs, 
		"Total CPUs should be >= idle + allocated")
}

// TestComprehensiveResourceListing tests all resource listing endpoints
func (suite *RealServerIntegrationTestSuite) TestComprehensiveResourceListing() {
	ctx := context.Background()

	// Test nodes listing with various options
	suite.T().Log("=== Testing Nodes Listing ===")
	nodeTests := []struct {
		name    string
		options *interfaces.ListNodesOptions
	}{
		{"Default", &interfaces.ListNodesOptions{Limit: 10}},
		{"With States", &interfaces.ListNodesOptions{Limit: 5, States: []string{"idle", "allocated"}}},
		{"Large Limit", &interfaces.ListNodesOptions{Limit: 50}},
	}

	for _, test := range nodeTests {
		suite.T().Logf("Testing nodes: %s", test.name)
		nodes, err := suite.client.Nodes().List(ctx, test.options)
		suite.Require().NoError(err, "Node listing should succeed for %s", test.name)
		suite.NotNil(nodes.Nodes, "Nodes should not be nil")
		suite.T().Logf("  Found %d nodes", len(nodes.Nodes))
		
		// Validate node data
		for i, node := range nodes.Nodes {
			if i >= 3 { // Log first 3 nodes
				break
			}
			suite.NotEmpty(node.Name, "Node name should not be empty")
			suite.NotEmpty(node.State, "Node state should not be empty")
			suite.GreaterOrEqual(node.CPUs, int32(0), "Node CPUs should be non-negative")
			suite.T().Logf("    Node: %s, State: %s, CPUs: %d", node.Name, node.State, node.CPUs)
		}
	}

	// Test partitions listing
	suite.T().Log("=== Testing Partitions Listing ===")
	partitionTests := []struct {
		name    string
		options *interfaces.ListPartitionsOptions
	}{
		{"Default", &interfaces.ListPartitionsOptions{Limit: 10}},
		{"With States", &interfaces.ListPartitionsOptions{Limit: 5, States: []string{"up"}}},
	}

	for _, test := range partitionTests {
		suite.T().Logf("Testing partitions: %s", test.name)
		partitions, err := suite.client.Partitions().List(ctx, test.options)
		suite.Require().NoError(err, "Partition listing should succeed for %s", test.name)
		suite.NotNil(partitions.Partitions, "Partitions should not be nil")
		suite.T().Logf("  Found %d partitions", len(partitions.Partitions))

		// Validate partition data
		for i, partition := range partitions.Partitions {
			if i >= 3 { // Log first 3 partitions
				break
			}
			suite.NotEmpty(partition.Name, "Partition name should not be empty")
			suite.NotEmpty(partition.State, "Partition state should not be empty")
			suite.GreaterOrEqual(partition.TotalNodes, int32(0), "Total nodes should be non-negative")
			suite.T().Logf("    Partition: %s, State: %s, Nodes: %d", 
				partition.Name, partition.State, partition.TotalNodes)
		}
	}

	// Test QoS listing
	suite.T().Log("=== Testing QoS Listing ===")
	qosTests := []struct {
		name    string
		options *interfaces.ListQoSOptions
	}{
		{"Default", &interfaces.ListQoSOptions{Limit: 10}},
		{"Small Limit", &interfaces.ListQoSOptions{Limit: 3}},
	}

	for _, test := range qosTests {
		suite.T().Logf("Testing QoS: %s", test.name)
		qosList, err := suite.client.QoS().List(ctx, test.options)
		suite.Require().NoError(err, "QoS listing should succeed for %s", test.name)
		suite.NotNil(qosList.QoS, "QoS list should not be nil")
		suite.T().Logf("  Found %d QoS entries", len(qosList.QoS))

		// Validate QoS data
		for i, qos := range qosList.QoS {
			if i >= 3 { // Log first 3 QoS entries
				break
			}
			suite.NotEmpty(qos.Name, "QoS name should not be empty")
			suite.GreaterOrEqual(qos.Priority, int32(0), "QoS priority should be non-negative")
			suite.T().Logf("    QoS: %s, Priority: %d, UsageFactor: %.2f", 
				qos.Name, qos.Priority, qos.UsageFactor)
		}
	}

	// Test jobs listing
	suite.T().Log("=== Testing Jobs Listing ===")
	jobTests := []struct {
		name    string
		options *interfaces.ListJobsOptions
	}{
		{"Default", &interfaces.ListJobsOptions{Limit: 10}},
		{"With States", &interfaces.ListJobsOptions{Limit: 5, States: []string{"running", "pending"}}},
		{"Recent Jobs", &interfaces.ListJobsOptions{Limit: 20}},
	}

	for _, test := range jobTests {
		suite.T().Logf("Testing jobs: %s", test.name)
		jobs, err := suite.client.Jobs().List(ctx, test.options)
		suite.Require().NoError(err, "Job listing should succeed for %s", test.name)
		suite.NotNil(jobs.Jobs, "Jobs should not be nil")
		suite.T().Logf("  Found %d jobs", len(jobs.Jobs))

		// Validate job data
		for i, job := range jobs.Jobs {
			if i >= 3 { // Log first 3 jobs
				break
			}
			suite.NotEmpty(job.ID, "Job ID should not be empty")
			suite.NotEmpty(job.State, "Job state should not be empty")
			suite.T().Logf("    Job: %s, Name: %s, State: %s", job.ID, job.Name, job.State)
		}
	}
}

// TestAdvancedJobOperations tests complex job operations
func (suite *RealServerIntegrationTestSuite) TestAdvancedJobOperations() {
	ctx := context.Background()

	// Get available partitions for job submission
	partitions, err := suite.client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{
		Limit: 5,
	})
	suite.Require().NoError(err)
	suite.Require().NotEmpty(partitions.Partitions, "Should have partitions available")

	targetPartition := partitions.Partitions[0].Name

	// Test 1: Job with environment variables
	suite.T().Log("=== Test 1: Job with environment variables ===")
	envJobName := fmt.Sprintf("integration-env-%d", time.Now().Unix())
	envSubmission := &interfaces.JobSubmission{
		Name:      envJobName,
		Script:    "#!/bin/bash\necho 'Environment test:'\necho 'SLURM_JOB_ID='$SLURM_JOB_ID\necho 'SLURM_JOB_NAME='$SLURM_JOB_NAME\necho 'TEST_VAR='$TEST_VAR\nenv | grep SLURM | head -10\nsleep 45",
		Partition: targetPartition,
		Nodes:     1,
		CPUs:      1,
		TimeLimit: 5,
		Environment: map[string]string{
			"TEST_VAR": "integration_test_value",
			"CUSTOM_VAR": "custom_value",
		},
	}

	envResponse, err := suite.client.Jobs().Submit(ctx, envSubmission)
	suite.Require().NoError(err, "Environment job submission should succeed")
	suite.submittedJobs = append(suite.submittedJobs, envResponse.JobID)
	suite.T().Logf("Environment job submitted: %s", envResponse.JobID)

	// Test 2: Job with working directory
	suite.T().Log("=== Test 2: Job with working directory ===")
	wdJobName := fmt.Sprintf("integration-wd-%d", time.Now().Unix())
	wdSubmission := &interfaces.JobSubmission{
		Name:        wdJobName,
		Script:      "#!/bin/bash\necho 'Working directory test:'\npwd\nls -la\necho 'Creating test file...'\ntouch integration_test_file.txt\nls -la integration_test_file.txt\nsleep 30",
		Partition:   targetPartition,
		Nodes:       1,
		CPUs:        1,
		TimeLimit:   5,
		WorkingDirectory: "/tmp",
	}

	wdResponse, err := suite.client.Jobs().Submit(ctx, wdSubmission)
	suite.Require().NoError(err, "Working directory job submission should succeed")
	suite.submittedJobs = append(suite.submittedJobs, wdResponse.JobID)
	suite.T().Logf("Working directory job submitted: %s", wdResponse.JobID)

	// Test 3: Job with specific CPU and memory requirements
	suite.T().Log("=== Test 3: Job with resource requirements ===")
	resourceJobName := fmt.Sprintf("integration-resource-%d", time.Now().Unix())
	resourceSubmission := &interfaces.JobSubmission{
		Name:      resourceJobName,
		Script:    "#!/bin/bash\necho 'Resource test:'\necho 'CPUs allocated: '$SLURM_CPUS_PER_TASK\necho 'Memory info:'\nfree -h\necho 'CPU info:'\nlscpu | grep -E '^CPU\\(s\\)|^Model name'\nsleep 30",
		Partition: targetPartition,
		Nodes:     1,
		CPUs:      2, // Request 2 CPUs
		Memory:    1024, // Request 1GB memory
		TimeLimit: 5,
	}

	resourceResponse, err := suite.client.Jobs().Submit(ctx, resourceSubmission)
	if err != nil {
		suite.T().Logf("Resource job submission failed (may be expected): %v", err)
	} else {
		suite.submittedJobs = append(suite.submittedJobs, resourceResponse.JobID)
		suite.T().Logf("Resource job submitted: %s", resourceResponse.JobID)
	}

	// Monitor job progress for a short time
	suite.T().Log("=== Monitoring job progress ===")
	time.Sleep(30 * time.Second)

	for _, jobID := range suite.submittedJobs[len(suite.submittedJobs)-3:] { // Check last 3 jobs
		job, err := suite.client.Jobs().Get(ctx, jobID)
		if err != nil {
			suite.T().Logf("Failed to get job %s: %v", jobID, err)
			continue
		}
		suite.T().Logf("Job %s state: %s", jobID, job.State)
	}
}

// TestErrorHandlingScenarios tests various error conditions
func (suite *RealServerIntegrationTestSuite) TestErrorHandlingScenarios() {
	ctx := context.Background()

	// Test 1: Invalid resource access
	suite.T().Log("=== Test 1: Invalid resource access ===")
	
	// Try to get non-existent job
	_, err := suite.client.Jobs().Get(ctx, "999999999")
	suite.Error(err, "Should fail for non-existent job")
	
	var slurmErr *errors.SlurmError
	if errors.As(err, &slurmErr) {
		suite.T().Logf("SLURM Error details: Code=%s, Category=%s, Status=%d", 
			slurmErr.Code, slurmErr.Category, slurmErr.StatusCode)
	}

	// Try to get non-existent QoS
	_, err = suite.client.QoS().Get(ctx, "nonexistent-qos-integration-test")
	suite.Error(err, "Should fail for non-existent QoS")

	// Test 2: Invalid job submission
	suite.T().Log("=== Test 2: Invalid job submission ===")
	
	// Submit job with invalid partition
	invalidSubmission := &interfaces.JobSubmission{
		Name:      "integration-invalid",
		Script:    "#!/bin/bash\necho 'This should fail'",
		Partition: "nonexistent-partition-integration-test",
		Nodes:     1,
		CPUs:      1,
		TimeLimit: 5,
	}

	_, err = suite.client.Jobs().Submit(ctx, invalidSubmission)
	suite.Error(err, "Should fail with invalid partition")
	suite.T().Logf("Invalid partition error: %v", err)

	// Test 3: Malformed requests
	suite.T().Log("=== Test 3: Testing malformed requests ===")
	
	// Try to submit job with zero time limit
	zeroTimeSubmission := &interfaces.JobSubmission{
		Name:      "integration-zero-time",
		Script:    "#!/bin/bash\necho 'Zero time limit'",
		Partition: "debug", // Assume debug partition exists
		Nodes:     1,
		CPUs:      1,
		TimeLimit: 0, // Invalid time limit
	}

	_, err = suite.client.Jobs().Submit(ctx, zeroTimeSubmission)
	if err != nil {
		suite.T().Logf("Zero time limit error (expected): %v", err)
	} else {
		suite.T().Log("Zero time limit job unexpectedly succeeded")
	}
}

// TestPerformanceAndReliability tests system performance and reliability
func (suite *RealServerIntegrationTestSuite) TestPerformanceAndReliability() {
	ctx := context.Background()

	// Test 1: Rapid successive requests
	suite.T().Log("=== Test 1: Rapid successive requests ===")
	
	start := time.Now()
	successCount := 0
	requestCount := 10

	for i := 0; i < requestCount; i++ {
		err := suite.client.Info().Ping(ctx)
		if err == nil {
			successCount++
		} else {
			suite.T().Logf("Ping %d failed: %v", i+1, err)
		}
	}

	duration := time.Since(start)
	avgLatency := duration / time.Duration(requestCount)
	
	suite.T().Logf("Rapid requests: %d/%d succeeded, avg latency: %v", 
		successCount, requestCount, avgLatency)
	suite.Greater(successCount, requestCount/2, "At least half of rapid requests should succeed")

	// Test 2: Large data retrieval
	suite.T().Log("=== Test 2: Large data retrieval ===")
	
	start = time.Now()
	jobs, err := suite.client.Jobs().List(ctx, &interfaces.ListJobsOptions{
		Limit: 100, // Request many jobs
	})
	duration = time.Since(start)
	
	if err == nil {
		suite.T().Logf("Retrieved %d jobs in %v", len(jobs.Jobs), duration)
		suite.LessOrEqual(duration, 30*time.Second, "Large data retrieval should complete within 30 seconds")
	} else {
		suite.T().Logf("Large data retrieval failed: %v", err)
	}

	// Test 3: Concurrent operations
	suite.T().Log("=== Test 3: Concurrent operations ===")
	
	type concurrentResult struct {
		operation string
		duration  time.Duration
		error     error
	}

	results := make(chan concurrentResult, 4)
	
	// Start concurrent operations
	go func() {
		start := time.Now()
		err := suite.client.Info().Ping(ctx)
		results <- concurrentResult{"ping", time.Since(start), err}
	}()

	go func() {
		start := time.Now()
		_, err := suite.client.Info().Version(ctx)
		results <- concurrentResult{"version", time.Since(start), err}
	}()

	go func() {
		start := time.Now()
		_, err := suite.client.Nodes().List(ctx, &interfaces.ListNodesOptions{Limit: 5})
		results <- concurrentResult{"nodes", time.Since(start), err}
	}()

	go func() {
		start := time.Now()
		_, err := suite.client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{Limit: 5})
		results <- concurrentResult{"partitions", time.Since(start), err}
	}()

	// Collect results
	concurrentSuccessCount := 0
	for i := 0; i < 4; i++ {
		result := <-results
		if result.error == nil {
			concurrentSuccessCount++
			suite.T().Logf("Concurrent %s: SUCCESS (%v)", result.operation, result.duration)
		} else {
			suite.T().Logf("Concurrent %s: FAILED (%v) - %v", result.operation, result.duration, result.error)
		}
	}

	suite.Greater(concurrentSuccessCount, 2, "At least 3 concurrent operations should succeed")
}

// TestRealServerIntegrationSuite runs the real server integration test suite
func TestRealServerIntegrationSuite(t *testing.T) {
	suite.Run(t, new(RealServerIntegrationTestSuite))
}
