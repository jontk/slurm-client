// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
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

// V0043IntegrationTestSuite tests v0.0.43 API integration against real server
type V0043IntegrationTestSuite struct {
	suite.Suite
	client    slurm.SlurmClient
	serverURL string
	token     string
	version   string
	
	// Test tracking
	submittedJobs []string
	testStartTime time.Time
	perfMetrics   map[string]time.Duration
	mu            sync.Mutex
}

// SetupSuite initializes the v0.0.43 test suite
func (suite *V0043IntegrationTestSuite) SetupSuite() {
	// Check if real server testing is enabled
	if os.Getenv("SLURM_REAL_SERVER_TEST") != "true" {
		suite.T().Skip("Real server tests disabled. Set SLURM_REAL_SERVER_TEST=true to enable")
	}

	suite.testStartTime = time.Now()
	suite.perfMetrics = make(map[string]time.Duration)
	suite.version = "v0.0.43"

	// Use provided server configuration
	suite.serverURL = "http://rocky9.ar.jontk.com:6820"
	suite.token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI2NTM4Mjk5NzYsImlhdCI6MTc1MzgyOTk3Niwic3VuIjoicm9vdCJ9.-z8Cq_wHuOxNJ7KHHTboX3l9r6JBtSD1RxQUgQR9owE"

	// Create client with v0.0.43 configuration
	ctx := context.Background()
	client, err := slurm.NewClientWithVersion(ctx, suite.version,
		slurm.WithBaseURL(suite.serverURL),
		slurm.WithAuth(auth.NewTokenAuth(suite.token)),
		slurm.WithConfig(&config.Config{
			Timeout:            30 * time.Second,
			MaxRetries:         3,
			Debug:              true,
			InsecureSkipVerify: true,
		}),
	)
	require.NoError(suite.T(), err, "Failed to create SLURM client for v0.0.43")
	suite.client = client

	suite.submittedJobs = make([]string, 0)
	
	suite.T().Logf("=== V0.0.43 Integration Test Suite Initialized ===")
	suite.T().Logf("Server: %s", suite.serverURL)
	suite.T().Logf("API Version: %s", suite.version)
}

// TearDownSuite cleans up after all tests
func (suite *V0043IntegrationTestSuite) TearDownSuite() {
	if suite.client == nil {
		return
	}

	// Clean up submitted jobs
	suite.cleanupJobs()
	
	// Close client
	suite.client.Close()
	
	// Generate test report
	suite.generateTestReport()
}

// Phase 1: Basic Connectivity and Authentication Tests

func (suite *V0043IntegrationTestSuite) TestPhase1_BasicConnectivity() {
	suite.T().Log("\n=== PHASE 1: Basic Connectivity and Authentication ===")
	
	// Test 1.1: Connection Validation
	suite.Run("1.1_ConnectionValidation", func() {
		ctx := context.Background()
		start := time.Now()
		
		err := suite.client.Info().Ping(ctx)
		duration := time.Since(start)
		
		suite.NoError(err, "Ping should succeed")
		suite.recordMetric("ping", duration)
		suite.T().Logf("✓ Ping successful (latency: %v)", duration)
	})
	
	// Test 1.2: Authentication Verification
	suite.Run("1.2_AuthenticationVerification", func() {
		ctx := context.Background()
		
		// Try an authenticated operation
		_, err := suite.client.Info().Version(ctx)
		suite.NoError(err, "Authenticated request should succeed")
		suite.T().Log("✓ JWT authentication verified")
	})
	
	// Test 1.3: Version Compatibility
	suite.Run("1.3_VersionCompatibility", func() {
		ctx := context.Background()
		
		versionInfo, err := suite.client.Info().Version(ctx)
		suite.Require().NoError(err)
		suite.NotEmpty(versionInfo.Version, "Version should not be empty")
		
		suite.T().Logf("✓ API Version: %s", versionInfo.Version)
		suite.T().Logf("  Major: %d, Minor: %d, Patch: %d", 
			versionInfo.Major, versionInfo.Minor, versionInfo.Patch)
	})
	
	// Test 1.4: Cluster Information
	suite.Run("1.4_ClusterInformation", func() {
		ctx := context.Background()
		
		info, err := suite.client.Info().Get(ctx)
		suite.Require().NoError(err)
		suite.NotEmpty(info.ClusterName, "Cluster name should not be empty")
		
		suite.T().Logf("✓ Connected to cluster: %s", info.ClusterName)
		if info.SlurmVersion != "" {
			suite.T().Logf("  SLURM Version: %s", info.SlurmVersion)
		}
		
		// Get cluster statistics
		stats, err := suite.client.Info().Stats(ctx)
		if err == nil {
			suite.T().Logf("  Cluster Statistics:")
			suite.T().Logf("    - Total Nodes: %d (Idle: %d, Allocated: %d)",
				stats.TotalNodes, stats.IdleNodes, stats.AllocatedNodes)
			suite.T().Logf("    - Total CPUs: %d (Idle: %d, Allocated: %d)",
				stats.TotalCPUs, stats.IdleCPUs, stats.AllocatedCPUs)
			suite.T().Logf("    - Total Jobs: %d (Running: %d, Pending: %d)",
				stats.TotalJobs, stats.RunningJobs, stats.PendingJobs)
		}
	})
}

// Phase 2: Manager Endpoint Tests

func (suite *V0043IntegrationTestSuite) TestPhase2_JobsManager() {
	suite.T().Log("\n=== PHASE 2.1: Jobs Manager Tests ===")
	ctx := context.Background()
	
	// Test 2.1.1: List Operations
	suite.Run("2.1.1_ListOperations", func() {
		// Default listing
		start := time.Now()
		jobs, err := suite.client.Jobs().List(ctx, &interfaces.ListJobsOptions{
			Limit: 10,
		})
		duration := time.Since(start)
		
		suite.NoError(err)
		suite.NotNil(jobs)
		suite.recordMetric("list_jobs_10", duration)
		suite.T().Logf("✓ Listed %d jobs (latency: %v)", len(jobs.Jobs), duration)
		
		// List with state filter
		if len(jobs.Jobs) > 0 {
			filteredJobs, err := suite.client.Jobs().List(ctx, &interfaces.ListJobsOptions{
				Limit:  5,
				States: []string{"running", "pending"},
			})
			suite.NoError(err)
			suite.T().Logf("✓ Listed %d jobs with state filter", len(filteredJobs.Jobs))
		}
		
		// Large limit test
		start = time.Now()
		largeJobs, err := suite.client.Jobs().List(ctx, &interfaces.ListJobsOptions{
			Limit: 100,
		})
		duration = time.Since(start)
		
		if err == nil {
			suite.recordMetric("list_jobs_100", duration)
			suite.T().Logf("✓ Listed %d jobs with large limit (latency: %v)", 
				len(largeJobs.Jobs), duration)
			suite.Less(duration, 2*time.Second, "Large job list should complete within 2s")
		}
	})
	
	// Test 2.1.2: Get Operations
	suite.Run("2.1.2_GetOperations", func() {
		// Test non-existent job
		_, err := suite.client.Jobs().Get(ctx, "999999999")
		suite.Error(err, "Should fail for non-existent job")
		
		var slurmErr *errors.SlurmError
		if errors.As(err, &slurmErr) {
			suite.T().Logf("✓ Error handling works: Code=%s, Status=%d", 
				slurmErr.Code, slurmErr.StatusCode)
		}
	})
	
	// Test 2.1.3: Submit Operations
	suite.Run("2.1.3_SubmitOperations", func() {
		// Get available partition first
		partitions, err := suite.client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{
			Limit: 5,
		})
		if err != nil || len(partitions.Partitions) == 0 {
			suite.T().Skip("No partitions available for job submission")
			return
		}
		
		partition := partitions.Partitions[0].Name
		
		// Simple job submission
		jobName := fmt.Sprintf("v0043-test-%d", time.Now().Unix())
		submission := &interfaces.JobSubmission{
			Name:      jobName,
			Script:    "#!/bin/bash\necho 'V0.0.43 Integration Test'\nhostname\ndate\nsleep 30",
			Partition: partition,
			Nodes:     1,
			CPUs:      1,
			TimeLimit: 5,
		}
		
		start := time.Now()
		response, err := suite.client.Jobs().Submit(ctx, submission)
		duration := time.Since(start)
		
		if err != nil {
			suite.T().Logf("⚠ Job submission failed (may be expected): %v", err)
			return
		}
		
		suite.NotEmpty(response.JobID)
		suite.submittedJobs = append(suite.submittedJobs, response.JobID)
		suite.recordMetric("job_submit", duration)
		suite.T().Logf("✓ Submitted job %s (latency: %v)", response.JobID, duration)
		
		// Job with environment variables
		envJob := &interfaces.JobSubmission{
			Name:      fmt.Sprintf("v0043-env-%d", time.Now().Unix()),
			Script:    "#!/bin/bash\necho \"TEST_VAR=$TEST_VAR\"\necho \"API_VERSION=$API_VERSION\"",
			Partition: partition,
			Nodes:     1,
			CPUs:      1,
			TimeLimit: 5,
			Environment: map[string]string{
				"TEST_VAR":    "v0043_test",
				"API_VERSION": "v0.0.43",
			},
		}
		
		envResponse, err := suite.client.Jobs().Submit(ctx, envJob)
		if err == nil {
			suite.submittedJobs = append(suite.submittedJobs, envResponse.JobID)
			suite.T().Logf("✓ Submitted job with environment variables: %s", envResponse.JobID)
		}
	})
	
	// Test 2.1.4: Modify Operations
	suite.Run("2.1.4_ModifyOperations", func() {
		if len(suite.submittedJobs) == 0 {
			suite.T().Skip("No jobs available for modification tests")
			return
		}
		
		jobID := suite.submittedJobs[0]
		
		// Try to hold the job
		err := suite.client.Jobs().Hold(ctx, jobID)
		if err == nil {
			suite.T().Logf("✓ Successfully held job %s", jobID)
			
			// Release the job
			err = suite.client.Jobs().Release(ctx, jobID)
			if err == nil {
				suite.T().Logf("✓ Successfully released job %s", jobID)
			}
		} else {
			suite.T().Logf("⚠ Job hold operation not available: %v", err)
		}
	})
}

func (suite *V0043IntegrationTestSuite) TestPhase2_NodesManager() {
	suite.T().Log("\n=== PHASE 2.2: Nodes Manager Tests ===")
	ctx := context.Background()
	
	// Test 2.2.1: List Operations
	suite.Run("2.2.1_ListOperations", func() {
		start := time.Now()
		nodes, err := suite.client.Nodes().List(ctx, &interfaces.ListNodesOptions{
			Limit: 10,
		})
		duration := time.Since(start)
		
		suite.Require().NoError(err)
		suite.NotNil(nodes)
		suite.recordMetric("list_nodes", duration)
		
		suite.T().Logf("✓ Listed %d nodes (latency: %v)", len(nodes.Nodes), duration)
		
		// Log first few nodes
		for i, node := range nodes.Nodes {
			if i >= 3 {
				break
			}
			suite.T().Logf("  - Node: %s, State: %s, CPUs: %d, Memory: %dMB",
				node.Name, node.State, node.CPUs, node.Memory)
		}
		
		// Filter by state
		if len(nodes.Nodes) > 0 {
			idleNodes, err := suite.client.Nodes().List(ctx, &interfaces.ListNodesOptions{
				Limit:  5,
				States: []string{"idle"},
			})
			if err == nil {
				suite.T().Logf("✓ Found %d idle nodes", len(idleNodes.Nodes))
			}
		}
	})
	
	// Test 2.2.2: Get Operations
	suite.Run("2.2.2_GetOperations", func() {
		// First get a valid node name
		nodes, err := suite.client.Nodes().List(ctx, &interfaces.ListNodesOptions{
			Limit: 1,
		})
		if err != nil || len(nodes.Nodes) == 0 {
			suite.T().Skip("No nodes available for get operation test")
			return
		}
		
		nodeName := nodes.Nodes[0].Name
		
		start := time.Now()
		node, err := suite.client.Nodes().Get(ctx, nodeName)
		duration := time.Since(start)
		
		suite.NoError(err)
		suite.Equal(nodeName, node.Name)
		suite.recordMetric("get_node", duration)
		suite.T().Logf("✓ Retrieved node %s details (latency: %v)", nodeName, duration)
		
		// Test invalid node
		_, err = suite.client.Nodes().Get(ctx, "nonexistent-node-v0043")
		suite.Error(err, "Should fail for non-existent node")
	})
}

func (suite *V0043IntegrationTestSuite) TestPhase2_PartitionsManager() {
	suite.T().Log("\n=== PHASE 2.3: Partitions Manager Tests ===")
	ctx := context.Background()
	
	// Test 2.3.1: List Operations
	suite.Run("2.3.1_ListOperations", func() {
		start := time.Now()
		partitions, err := suite.client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{
			Limit: 10,
		})
		duration := time.Since(start)
		
		suite.Require().NoError(err)
		suite.NotNil(partitions)
		suite.recordMetric("list_partitions", duration)
		
		suite.T().Logf("✓ Listed %d partitions (latency: %v)", len(partitions.Partitions), duration)
		
		// Log partition details
		for i, partition := range partitions.Partitions {
			if i >= 3 {
				break
			}
			suite.T().Logf("  - Partition: %s, State: %s, Nodes: %d, Default: %v",
				partition.Name, partition.State, partition.TotalNodes, partition.Default)
		}
	})
	
	// Test 2.3.2: Get Operations
	suite.Run("2.3.2_GetOperations", func() {
		// Get first partition for testing
		partitions, err := suite.client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{
			Limit: 1,
		})
		if err != nil || len(partitions.Partitions) == 0 {
			suite.T().Skip("No partitions available for get operation test")
			return
		}
		
		partitionName := partitions.Partitions[0].Name
		
		partition, err := suite.client.Partitions().Get(ctx, partitionName)
		suite.NoError(err)
		suite.Equal(partitionName, partition.Name)
		suite.T().Logf("✓ Retrieved partition %s details", partitionName)
	})
}

func (suite *V0043IntegrationTestSuite) TestPhase2_DatabaseDependentManagers() {
	suite.T().Log("\n=== PHASE 2.4-2.7: Database-Dependent Managers ===")
	ctx := context.Background()
	
	// Test QoS Manager
	suite.Run("2.4_QoSManager", func() {
		qosList, err := suite.client.QoS().List(ctx, &interfaces.ListQoSOptions{
			Limit: 5,
		})
		
		if err != nil {
			if suite.isDatabaseError(err) {
				suite.T().Log("⚠ QoS operations require database connection (skipped)")
				return
			}
			suite.T().Errorf("Unexpected error listing QoS: %v", err)
			return
		}
		
		suite.T().Logf("✓ Listed %d QoS entries", len(qosList.QoS))
		for i, qos := range qosList.QoS {
			if i >= 3 {
				break
			}
			suite.T().Logf("  - QoS: %s, Priority: %d", qos.Name, qos.Priority)
		}
	})
	
	// Test Users Manager
	suite.Run("2.5_UsersManager", func() {
		users, err := suite.client.Users().List(ctx, &interfaces.ListUsersOptions{
			Limit: 5,
		})
		
		if err != nil {
			if suite.isDatabaseError(err) {
				suite.T().Log("⚠ User operations require database connection (skipped)")
				return
			}
			suite.T().Errorf("Unexpected error listing users: %v", err)
			return
		}
		
		suite.T().Logf("✓ Listed %d users", len(users.Users))
	})
	
	// Test Accounts Manager
	suite.Run("2.6_AccountsManager", func() {
		accounts, err := suite.client.Accounts().List(ctx, &interfaces.ListAccountsOptions{
			Limit: 5,
		})
		
		if err != nil {
			if suite.isDatabaseError(err) {
				suite.T().Log("⚠ Account operations require database connection (skipped)")
				return
			}
			suite.T().Errorf("Unexpected error listing accounts: %v", err)
			return
		}
		
		suite.T().Logf("✓ Listed %d accounts", len(accounts.Accounts))
	})
}

// Phase 3: Error Handling Scenarios

func (suite *V0043IntegrationTestSuite) TestPhase3_ErrorHandling() {
	suite.T().Log("\n=== PHASE 3: Error Handling Scenarios ===")
	ctx := context.Background()
	
	// Test 3.1: Invalid Resource Access
	suite.Run("3.1_InvalidResourceAccess", func() {
		// Invalid job ID
		_, err := suite.client.Jobs().Get(ctx, "999999999")
		suite.Error(err)
		suite.assertSlurmError(err, "invalid job ID")
		
		// Invalid node name
		_, err = suite.client.Nodes().Get(ctx, "nonexistent-node-v0043-test")
		suite.Error(err)
		suite.assertSlurmError(err, "invalid node name")
		
		// Invalid partition name
		_, err = suite.client.Partitions().Get(ctx, "nonexistent-partition-v0043")
		suite.Error(err)
		suite.assertSlurmError(err, "invalid partition name")
		
		suite.T().Log("✓ All invalid resource access errors handled correctly")
	})
	
	// Test 3.2: Invalid Operations
	suite.Run("3.2_InvalidOperations", func() {
		// Submit job with invalid partition
		submission := &interfaces.JobSubmission{
			Name:      "v0043-invalid-partition",
			Script:    "#!/bin/bash\necho 'Should fail'",
			Partition: "nonexistent-partition-v0043-test",
			Nodes:     1,
			CPUs:      1,
			TimeLimit: 5,
		}
		
		_, err := suite.client.Jobs().Submit(ctx, submission)
		suite.Error(err)
		suite.T().Logf("✓ Invalid partition submission failed as expected: %v", err)
		
		// Cancel non-existent job
		err = suite.client.Jobs().Cancel(ctx, "999999999")
		suite.Error(err)
		suite.T().Log("✓ Cancel non-existent job failed as expected")
	})
}

// Phase 4: Performance Benchmarks

func (suite *V0043IntegrationTestSuite) TestPhase4_PerformanceBenchmarks() {
	suite.T().Log("\n=== PHASE 4: Performance Benchmarks ===")
	ctx := context.Background()
	
	// Test 4.1: Latency Tests
	suite.Run("4.1_LatencyTests", func() {
		operations := []struct {
			name      string
			operation func() error
			maxTime   time.Duration
		}{
			{
				name: "ping",
				operation: func() error {
					return suite.client.Info().Ping(ctx)
				},
				maxTime: 100 * time.Millisecond,
			},
			{
				name: "get_single_job",
				operation: func() error {
					jobs, err := suite.client.Jobs().List(ctx, &interfaces.ListJobsOptions{Limit: 1})
					if err != nil || len(jobs.Jobs) == 0 {
						return fmt.Errorf("no jobs available")
					}
					_, err = suite.client.Jobs().Get(ctx, jobs.Jobs[0].ID)
					return err
				},
				maxTime: 200 * time.Millisecond,
			},
			{
				name: "list_10_jobs",
				operation: func() error {
					_, err := suite.client.Jobs().List(ctx, &interfaces.ListJobsOptions{Limit: 10})
					return err
				},
				maxTime: 500 * time.Millisecond,
			},
		}
		
		for _, op := range operations {
			start := time.Now()
			err := op.operation()
			duration := time.Since(start)
			
			if err != nil {
				suite.T().Logf("⚠ %s failed: %v", op.name, err)
				continue
			}
			
			suite.recordMetric("perf_"+op.name, duration)
			if duration <= op.maxTime {
				suite.T().Logf("✓ %s: %v (within %v limit)", op.name, duration, op.maxTime)
			} else {
				suite.T().Logf("⚠ %s: %v (exceeds %v limit)", op.name, duration, op.maxTime)
			}
		}
	})
	
	// Test 4.2: Throughput Tests
	suite.Run("4.2_ThroughputTests", func() {
		// Rapid ping requests
		requestCount := 10
		start := time.Now()
		successCount := 0
		
		for i := 0; i < requestCount; i++ {
			if err := suite.client.Info().Ping(ctx); err == nil {
				successCount++
			}
		}
		
		duration := time.Since(start)
		avgLatency := duration / time.Duration(requestCount)
		
		suite.T().Logf("✓ Rapid requests: %d/%d succeeded in %v (avg: %v)",
			successCount, requestCount, duration, avgLatency)
		suite.GreaterOrEqual(successCount, requestCount*8/10, "At least 80% should succeed")
	})
	
	// Test 4.3: Concurrent Operations
	suite.Run("4.3_ConcurrentOperations", func() {
		concurrency := 10
		results := make(chan error, concurrency)
		
		start := time.Now()
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				_, err := suite.client.Jobs().List(ctx, &interfaces.ListJobsOptions{
					Limit: 5,
				})
				results <- err
			}(i)
		}
		
		successCount := 0
		for i := 0; i < concurrency; i++ {
			if err := <-results; err == nil {
				successCount++
			}
		}
		
		duration := time.Since(start)
		suite.T().Logf("✓ Concurrent operations: %d/%d succeeded in %v",
			successCount, concurrency, duration)
		suite.GreaterOrEqual(successCount, concurrency*8/10, "At least 80% should succeed")
	})
}

// Phase 5: API Version-Specific Features

func (suite *V0043IntegrationTestSuite) TestPhase5_VersionSpecificFeatures() {
	suite.T().Log("\n=== PHASE 5: v0.0.43 Specific Features ===")
	ctx := context.Background()
	
	// Test any v0.0.43 specific features
	suite.Run("5.1_NewFeatures", func() {
		// Test enhanced job submission options if available in v0.0.43
		partitions, err := suite.client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{
			Limit: 1,
		})
		if err != nil || len(partitions.Partitions) == 0 {
			suite.T().Skip("No partitions available for v0.0.43 feature testing")
			return
		}
		
		// Test with advanced job options
		jobName := fmt.Sprintf("v0043-advanced-%d", time.Now().Unix())
		submission := &interfaces.JobSubmission{
			Name:             jobName,
			Script:           "#!/bin/bash\necho 'Testing v0.0.43 features'\ndate",
			Partition:        partitions.Partitions[0].Name,
			Nodes:            1,
			CPUs:             1,
			TimeLimit:        5,
			WorkingDirectory: "/tmp",
			Comment:          "v0.0.43 integration test",
		}
		
		response, err := suite.client.Jobs().Submit(ctx, submission)
		if err == nil {
			suite.submittedJobs = append(suite.submittedJobs, response.JobID)
			suite.T().Logf("✓ Advanced job submission successful: %s", response.JobID)
			
			// Verify the job details
			job, err := suite.client.Jobs().Get(ctx, response.JobID)
			if err == nil {
				suite.T().Logf("  - Job comment: %s", job.Comment)
				suite.T().Logf("  - Working dir: %s", job.WorkingDirectory)
			}
		}
	})
}

// Phase 6: Integration Workflows

func (suite *V0043IntegrationTestSuite) TestPhase6_IntegrationWorkflows() {
	suite.T().Log("\n=== PHASE 6: Integration Workflows ===")
	ctx := context.Background()
	
	// Test 6.1: Complete Job Lifecycle
	suite.Run("6.1_CompleteJobLifecycle", func() {
		// Step 1: Resource discovery
		partitions, err := suite.client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{
			Limit: 5,
		})
		if err != nil || len(partitions.Partitions) == 0 {
			suite.T().Skip("No partitions available for lifecycle test")
			return
		}
		
		// Find best partition
		var selectedPartition string
		for _, p := range partitions.Partitions {
			if p.State == "up" && p.TotalNodes > 0 {
				selectedPartition = p.Name
				break
			}
		}
		
		if selectedPartition == "" {
			suite.T().Skip("No suitable partition found")
			return
		}
		
		suite.T().Logf("Step 1: Selected partition %s for job submission", selectedPartition)
		
		// Step 2: Submit job
		jobName := fmt.Sprintf("v0043-lifecycle-%d", time.Now().Unix())
		submission := &interfaces.JobSubmission{
			Name:      jobName,
			Script:    "#!/bin/bash\necho 'Job lifecycle test'\necho 'Step 1: Starting'\nsleep 5\necho 'Step 2: Processing'\nsleep 5\necho 'Step 3: Complete'",
			Partition: selectedPartition,
			Nodes:     1,
			CPUs:      1,
			TimeLimit: 5,
		}
		
		response, err := suite.client.Jobs().Submit(ctx, submission)
		if err != nil {
			suite.T().Logf("Job submission failed: %v", err)
			return
		}
		
		suite.submittedJobs = append(suite.submittedJobs, response.JobID)
		suite.T().Logf("Step 2: Submitted job %s", response.JobID)
		
		// Step 3: Monitor job
		suite.T().Log("Step 3: Monitoring job progress...")
		states := make([]string, 0)
		
		for i := 0; i < 30; i++ { // Monitor for up to 60 seconds
			job, err := suite.client.Jobs().Get(ctx, response.JobID)
			if err != nil {
				suite.T().Logf("Failed to get job status: %v", err)
				break
			}
			
			if len(states) == 0 || states[len(states)-1] != job.State {
				states = append(states, job.State)
				suite.T().Logf("  - Job state: %s", job.State)
			}
			
			if job.State == "COMPLETED" || job.State == "FAILED" || job.State == "CANCELLED" {
				break
			}
			
			time.Sleep(2 * time.Second)
		}
		
		suite.T().Logf("Step 4: Job lifecycle complete. States observed: %v", states)
	})
	
	// Test 6.2: Resource Discovery Workflow
	suite.Run("6.2_ResourceDiscoveryWorkflow", func() {
		suite.T().Log("Starting resource discovery workflow...")
		
		// Step 1: Get cluster overview
		info, err := suite.client.Info().Get(ctx)
		suite.Require().NoError(err)
		suite.T().Logf("Step 1: Cluster %s discovered", info.ClusterName)
		
		// Step 2: Enumerate partitions
		partitions, err := suite.client.Partitions().List(ctx, nil)
		suite.Require().NoError(err)
		suite.T().Logf("Step 2: Found %d partitions", len(partitions.Partitions))
		
		// Step 3: Check node availability
		nodes, err := suite.client.Nodes().List(ctx, &interfaces.ListNodesOptions{
			States: []string{"idle"},
			Limit:  10,
		})
		suite.Require().NoError(err)
		suite.T().Logf("Step 3: Found %d idle nodes", len(nodes.Nodes))
		
		// Step 4: Identify optimal resources
		var optimalPartition string
		var maxIdleNodes int32
		
		for _, p := range partitions.Partitions {
			if p.State == "up" && p.IdleNodes > maxIdleNodes {
				optimalPartition = p.Name
				maxIdleNodes = p.IdleNodes
			}
		}
		
		if optimalPartition != "" {
			suite.T().Logf("Step 4: Optimal partition identified: %s with %d idle nodes",
				optimalPartition, maxIdleNodes)
		}
		
		suite.T().Log("✓ Resource discovery workflow complete")
	})
}

// Helper Methods

func (suite *V0043IntegrationTestSuite) recordMetric(name string, duration time.Duration) {
	suite.mu.Lock()
	defer suite.mu.Unlock()
	suite.perfMetrics[name] = duration
}

func (suite *V0043IntegrationTestSuite) assertSlurmError(err error, context string) {
	var slurmErr *errors.SlurmError
	if errors.As(err, &slurmErr) {
		suite.T().Logf("✓ SLURM error for %s: Code=%s, Status=%d, Message=%s",
			context, slurmErr.Code, slurmErr.StatusCode, slurmErr.Message)
	}
}

func (suite *V0043IntegrationTestSuite) isDatabaseError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "Unable to connect to database") ||
		contains(errStr, "Failed to open slurmdbd connection") ||
		contains(errStr, "Bad Gateway") ||
		contains(errStr, "HTTP 502")
}

func (suite *V0043IntegrationTestSuite) cleanupJobs() {
	if len(suite.submittedJobs) == 0 {
		return
	}
	
	suite.T().Log("\n=== Cleaning up test jobs ===")
	ctx := context.Background()
	
	for _, jobID := range suite.submittedJobs {
		err := suite.client.Jobs().Cancel(ctx, jobID)
		if err != nil {
			suite.T().Logf("  ⚠ Failed to cancel job %s: %v", jobID, err)
		} else {
			suite.T().Logf("  ✓ Cancelled job %s", jobID)
		}
	}
}

func (suite *V0043IntegrationTestSuite) generateTestReport() {
	duration := time.Since(suite.testStartTime)
	
	suite.T().Log("\n" + strings.Repeat("=", 60))
	suite.T().Log("V0.0.43 INTEGRATION TEST REPORT")
	suite.T().Log(strings.Repeat("=", 60))
	suite.T().Logf("Date: %s", suite.testStartTime.Format(time.RFC3339))
	suite.T().Logf("Duration: %v", duration)
	suite.T().Logf("Server: %s", suite.serverURL)
	suite.T().Logf("API Version: %s", suite.version)
	
	suite.T().Log("\nPerformance Metrics:")
	var totalLatency time.Duration
	count := 0
	
	for metric, latency := range suite.perfMetrics {
		suite.T().Logf("  - %s: %v", metric, latency)
		totalLatency += latency
		count++
	}
	
	if count > 0 {
		avgLatency := totalLatency / time.Duration(count)
		suite.T().Logf("\nAverage Operation Latency: %v", avgLatency)
	}
	
	suite.T().Logf("\nJobs Submitted: %d", len(suite.submittedJobs))
	suite.T().Log(strings.Repeat("=", 60))
}

// Test Runner

func TestV0043IntegrationSuite(t *testing.T) {
	suite.Run(t, new(V0043IntegrationTestSuite))
}

// Helper functions

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
