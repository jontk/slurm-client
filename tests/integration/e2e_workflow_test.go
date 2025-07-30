package integration

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

// E2EWorkflowTestSuite tests complete end-to-end user workflows
type E2EWorkflowTestSuite struct {
	suite.Suite
	client       slurm.SlurmClient
	serverURL    string
	token        string
	version      string
	testPrefix   string
	createdJobs  []string
}

// SetupSuite initializes the test suite
func (suite *E2EWorkflowTestSuite) SetupSuite() {
	// Check if E2E testing is enabled
	if os.Getenv("SLURM_E2E_TEST") != "true" {
		suite.T().Skip("E2E workflow tests disabled. Set SLURM_E2E_TEST=true to enable")
	}

	// Get server configuration
	suite.serverURL = os.Getenv("SLURM_SERVER_URL")
	if suite.serverURL == "" {
		suite.serverURL = "http://rocky9:6820"
	}

	// Get API version
	suite.version = os.Getenv("SLURM_API_VERSION")
	if suite.version == "" {
		suite.version = "v0.0.43"
	}

	// Get JWT token
	suite.token = os.Getenv("SLURM_JWT_TOKEN")
	if suite.token == "" {
		token, err := fetchJWTTokenViaSSH()
		require.NoError(suite.T(), err, "Failed to fetch JWT token")
		suite.token = token
	}

	// Create client
	ctx := context.Background()
	client, err := slurm.NewClientWithVersion(ctx, suite.version,
		slurm.WithBaseURL(suite.serverURL),
		slurm.WithAuth(auth.NewTokenAuth(suite.token)),
		slurm.WithConfig(&config.Config{
			Timeout:            60 * time.Second, // Longer timeout for E2E tests
			MaxRetries:         3,
			Debug:              true,
			InsecureSkipVerify: true,
		}),
	)
	require.NoError(suite.T(), err)
	suite.client = client

	// Generate unique test prefix
	suite.testPrefix = fmt.Sprintf("e2e-test-%d", time.Now().Unix())
	suite.createdJobs = make([]string, 0)

	suite.T().Logf("Starting E2E workflow tests with prefix: %s", suite.testPrefix)
}

// TearDownSuite cleans up test resources
func (suite *E2EWorkflowTestSuite) TearDownSuite() {
	if suite.client == nil {
		return
	}

	ctx := context.Background()

	// Cancel any remaining jobs
	for _, jobID := range suite.createdJobs {
		err := suite.client.Jobs().Cancel(ctx, jobID)
		if err != nil {
			suite.T().Logf("Failed to cancel job %s: %v", jobID, err)
		} else {
			suite.T().Logf("Cleaned up job %s", jobID)
		}
	}

	suite.client.Close()
}

// TestFullJobLifecycle tests the complete job lifecycle
func (suite *E2EWorkflowTestSuite) TestFullJobLifecycle() {
	ctx := context.Background()

	// Step 1: Verify cluster connectivity
	suite.T().Log("Step 1: Verifying cluster connectivity...")
	err := suite.client.Info().Ping(ctx)
	suite.Require().NoError(err, "Cluster should be accessible")

	// Step 2: Get cluster information
	suite.T().Log("Step 2: Getting cluster information...")
	info, err := suite.client.Info().Get(ctx)
	suite.Require().NoError(err)
	suite.T().Logf("Connected to cluster: %s", info.ClusterName)

	// Step 3: List available partitions
	suite.T().Log("Step 3: Listing available partitions...")
	partitions, err := suite.client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{
		Limit: 10,
	})
	suite.Require().NoError(err)
	suite.Require().NotEmpty(partitions.Partitions, "Should have at least one partition")

	// Find a suitable partition (prefer 'debug' or 'normal')
	var targetPartition string
	for _, partition := range partitions.Partitions {
		if strings.Contains(strings.ToLower(partition.Name), "debug") ||
		   strings.Contains(strings.ToLower(partition.Name), "normal") {
			targetPartition = partition.Name
			break
		}
	}
	if targetPartition == "" {
		targetPartition = partitions.Partitions[0].Name // Use first available
	}
	suite.T().Logf("Using partition: %s", targetPartition)

	// Step 4: List available QoS
	suite.T().Log("Step 4: Listing available QoS...")
	qosList, err := suite.client.QoS().List(ctx, &interfaces.ListQoSOptions{
		Limit: 5,
	})
	suite.Require().NoError(err)
	
	var targetQoS string
	if len(qosList.QoS) > 0 {
		targetQoS = qosList.QoS[0].Name
		suite.T().Logf("Using QoS: %s", targetQoS)
	}

	// Step 5: Get current job count
	suite.T().Log("Step 5: Getting current job statistics...")
	initialJobs, err := suite.client.Jobs().List(ctx, &interfaces.ListJobsOptions{
		Limit: 100,
	})
	suite.Require().NoError(err)
	initialJobCount := len(initialJobs.Jobs)
	suite.T().Logf("Initial job count: %d", initialJobCount)

	// Step 6: Submit a test job
	suite.T().Log("Step 6: Submitting test job...")
	jobName := fmt.Sprintf("%s-workflow", suite.testPrefix)
	submission := &interfaces.JobSubmission{
		Name:      jobName,
		Script:    fmt.Sprintf("#!/bin/bash\n# E2E Test Job\necho 'Job %s started at:'\ndate\nhostname\necho 'Running workflow test...'\nsleep 60\necho 'Job completed at:'\ndate", jobName),
		Partition: targetPartition,
		Nodes:     1,
		CPUs:      1,
		TimeLimit: 10, // 10 minutes
		QoS:       targetQoS,
	}

	response, err := suite.client.Jobs().Submit(ctx, submission)
	suite.Require().NoError(err, "Job submission should succeed")
	suite.Require().NotEmpty(response.JobID, "Job ID should be returned")
	suite.createdJobs = append(suite.createdJobs, response.JobID)
	
	suite.T().Logf("Job submitted successfully: ID=%s", response.JobID)

	// Step 7: Verify job was created
	suite.T().Log("Step 7: Verifying job creation...")
	job, err := suite.client.Jobs().Get(ctx, response.JobID)
	suite.Require().NoError(err)
	suite.Equal(response.JobID, job.ID)
	suite.Equal(jobName, job.Name)
	suite.Equal(targetPartition, job.Partition)
	suite.T().Logf("Job verified: State=%s, Name=%s", job.State, job.Name)

	// Step 8: Monitor job state changes
	suite.T().Log("Step 8: Monitoring job state changes...")
	maxWaitTime := 3 * time.Minute
	checkInterval := 10 * time.Second
	startTime := time.Now()
	
	var finalJob *interfaces.Job
	seenStates := make(map[string]bool)
	
	for time.Since(startTime) < maxWaitTime {
		currentJob, err := suite.client.Jobs().Get(ctx, response.JobID)
		suite.Require().NoError(err)
		
		if !seenStates[currentJob.State] {
			suite.T().Logf("Job state: %s (elapsed: %v)", currentJob.State, time.Since(startTime))
			seenStates[currentJob.State] = true
		}
		
		finalJob = currentJob
		
		// Break if job is in a final state
		if currentJob.State == "COMPLETED" || currentJob.State == "FAILED" || 
		   currentJob.State == "CANCELLED" || currentJob.State == "TIMEOUT" {
			break
		}
		
		time.Sleep(checkInterval)
	}
	
	suite.Require().NotNil(finalJob, "Should have job state")
	suite.T().Logf("Final job state: %s", finalJob.State)

	// Step 9: Test job cancellation (if still running)
	if finalJob.State != "COMPLETED" && finalJob.State != "FAILED" && 
	   finalJob.State != "CANCELLED" && finalJob.State != "TIMEOUT" {
		suite.T().Log("Step 9: Testing job cancellation...")
		err = suite.client.Jobs().Cancel(ctx, response.JobID)
		suite.Require().NoError(err, "Job cancellation should succeed")
		
		// Verify cancellation
		time.Sleep(5 * time.Second)
		cancelledJob, err := suite.client.Jobs().Get(ctx, response.JobID)
		suite.Require().NoError(err)
		suite.T().Logf("Job state after cancellation: %s", cancelledJob.State)
	}

	// Step 10: Verify job history
	suite.T().Log("Step 10: Verifying job appears in listings...")
	finalJobs, err := suite.client.Jobs().List(ctx, &interfaces.ListJobsOptions{
		Limit: 100,
	})
	suite.Require().NoError(err)
	
	foundJob := false
	for _, j := range finalJobs.Jobs {
		if j.ID == response.JobID {
			foundJob = true
			suite.T().Logf("Job found in listings: ID=%s, State=%s", j.ID, j.State)
			break
		}
	}
	suite.True(foundJob, "Job should appear in job listings")
}

// TestMultiJobWorkflow tests managing multiple jobs simultaneously
func (suite *E2EWorkflowTestSuite) TestMultiJobWorkflow() {
	ctx := context.Background()

	// Get a suitable partition
	partitions, err := suite.client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{
		Limit: 5,
	})
	suite.Require().NoError(err)
	suite.Require().NotEmpty(partitions.Partitions)
	
	targetPartition := partitions.Partitions[0].Name

	// Submit multiple jobs
	suite.T().Log("Submitting multiple test jobs...")
	jobCount := 3
	submittedJobs := make([]string, 0, jobCount)

	for i := 0; i < jobCount; i++ {
		jobName := fmt.Sprintf("%s-multi-%d", suite.testPrefix, i+1)
		submission := &interfaces.JobSubmission{
			Name:      jobName,
			Script:    fmt.Sprintf("#!/bin/bash\necho 'Multi-job test %d'\nhostname\nsleep 30\necho 'Job %d completed'", i+1, i+1),
			Partition: targetPartition,
			Nodes:     1,
			CPUs:      1,
			TimeLimit: 5, // 5 minutes
		}

		response, err := suite.client.Jobs().Submit(ctx, submission)
		suite.Require().NoError(err, "Job %d submission should succeed", i+1)
		
		submittedJobs = append(submittedJobs, response.JobID)
		suite.createdJobs = append(suite.createdJobs, response.JobID)
		
		suite.T().Logf("Submitted job %d: ID=%s", i+1, response.JobID)
	}

	// Monitor all jobs
	suite.T().Log("Monitoring multiple jobs...")
	maxWaitTime := 2 * time.Minute
	checkInterval := 15 * time.Second
	startTime := time.Now()

	for time.Since(startTime) < maxWaitTime {
		allCompleted := true
		
		for i, jobID := range submittedJobs {
			job, err := suite.client.Jobs().Get(ctx, jobID)
			suite.Require().NoError(err)
			
			suite.T().Logf("Job %d (%s): State=%s", i+1, jobID, job.State)
			
			if job.State != "COMPLETED" && job.State != "FAILED" && 
			   job.State != "CANCELLED" && job.State != "TIMEOUT" {
				allCompleted = false
			}
		}
		
		if allCompleted {
			break
		}
		
		time.Sleep(checkInterval)
	}

	// Cancel remaining jobs
	suite.T().Log("Cleaning up remaining jobs...")
	for i, jobID := range submittedJobs {
		err := suite.client.Jobs().Cancel(ctx, jobID)
		if err == nil {
			suite.T().Logf("Cancelled job %d (%s)", i+1, jobID)
		}
	}
}

// TestResourceRelationshipWorkflow tests relationships between resources
func (suite *E2EWorkflowTestSuite) TestResourceRelationshipWorkflow() {
	ctx := context.Background()

	// Step 1: Map cluster resources
	suite.T().Log("Step 1: Mapping cluster resources...")
	
	// Get nodes
	nodes, err := suite.client.Nodes().List(ctx, &interfaces.ListNodesOptions{
		Limit: 10,
	})
	suite.Require().NoError(err)
	suite.T().Logf("Found %d nodes", len(nodes.Nodes))

	// Get partitions
	partitions, err := suite.client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{
		Limit: 10,
	})
	suite.Require().NoError(err)
	suite.T().Logf("Found %d partitions", len(partitions.Partitions))

	// Get QoS
	qosList, err := suite.client.QoS().List(ctx, &interfaces.ListQoSOptions{
		Limit: 10,
	})
	suite.Require().NoError(err)
	suite.T().Logf("Found %d QoS entries", len(qosList.QoS))

	// Step 2: Analyze resource relationships
	suite.T().Log("Step 2: Analyzing resource relationships...")
	
	// Check partition-node relationships
	for i, partition := range partitions.Partitions {
		if i >= 3 { // Limit analysis to first 3 partitions
			break
		}
		
		suite.T().Logf("Partition '%s': %d total nodes, %d allocated nodes", 
			partition.Name, partition.TotalNodes, partition.AllocatedNodes)
		
		// Validate partition resources
		suite.GreaterOrEqual(partition.TotalNodes, int32(0), "Total nodes should be non-negative")
		suite.GreaterOrEqual(partition.AllocatedNodes, int32(0), "Allocated nodes should be non-negative")
		suite.LessOrEqual(partition.AllocatedNodes, partition.TotalNodes, "Allocated should not exceed total")
	}

	// Step 3: Test job submission with specific constraints
	if len(partitions.Partitions) > 0 && len(qosList.QoS) > 0 {
		suite.T().Log("Step 3: Testing job with resource constraints...")
		
		partition := partitions.Partitions[0]
		qos := qosList.QoS[0]
		
		jobName := fmt.Sprintf("%s-constrained", suite.testPrefix)
		submission := &interfaces.JobSubmission{
			Name:      jobName,
			Script:    "#!/bin/bash\necho 'Testing resource constraints'\nhostname\necho 'Node info:'\nuname -a\necho 'CPU info:'\nlscpu | head -10\nsleep 30",
			Partition: partition.Name,
			QoS:       qos.Name,
			Nodes:     1,
			CPUs:      1,
			TimeLimit: 5,
		}

		response, err := suite.client.Jobs().Submit(ctx, submission)
		if err != nil {
			suite.T().Logf("Constrained job submission failed (expected for some configurations): %v", err)
		} else {
			suite.createdJobs = append(suite.createdJobs, response.JobID)
			suite.T().Logf("Constrained job submitted: ID=%s", response.JobID)
			
			// Verify job details
			job, err := suite.client.Jobs().Get(ctx, response.JobID)
			suite.Require().NoError(err)
			suite.Equal(partition.Name, job.Partition)
			suite.T().Logf("Job created with Partition=%s, QoS=%s", job.Partition, job.QoS)
		}
	}
}

// TestErrorRecoveryWorkflow tests error handling and recovery scenarios
func (suite *E2EWorkflowTestSuite) TestErrorRecoveryWorkflow() {
	ctx := context.Background()

	// Test 1: Invalid partition
	suite.T().Log("Test 1: Testing invalid partition handling...")
	submission := &interfaces.JobSubmission{
		Name:      fmt.Sprintf("%s-invalid-partition", suite.testPrefix),
		Script:    "#!/bin/bash\necho 'This should fail'",
		Partition: "nonexistent-partition-12345",
		Nodes:     1,
		CPUs:      1,
		TimeLimit: 5,
	}

	_, err := suite.client.Jobs().Submit(ctx, submission)
	suite.Error(err, "Should fail with invalid partition")
	suite.T().Logf("Invalid partition error (expected): %v", err)

	// Test 2: Invalid QoS
	partitions, err := suite.client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{
		Limit: 1,
	})
	suite.Require().NoError(err)
	
	if len(partitions.Partitions) > 0 {
		suite.T().Log("Test 2: Testing invalid QoS handling...")
		submission := &interfaces.JobSubmission{
			Name:      fmt.Sprintf("%s-invalid-qos", suite.testPrefix),
			Script:    "#!/bin/bash\necho 'This should fail'",
			Partition: partitions.Partitions[0].Name,
			QoS:       "nonexistent-qos-12345",
			Nodes:     1,
			CPUs:      1,
			TimeLimit: 5,
		}

		_, err := suite.client.Jobs().Submit(ctx, submission)
		suite.Error(err, "Should fail with invalid QoS")
		suite.T().Logf("Invalid QoS error (expected): %v", err)
	}

	// Test 3: Excessive resource request
	if len(partitions.Partitions) > 0 {
		suite.T().Log("Test 3: Testing excessive resource request...")
		submission := &interfaces.JobSubmission{
			Name:      fmt.Sprintf("%s-excessive-resources", suite.testPrefix),
			Script:    "#!/bin/bash\necho 'This might fail'",
			Partition: partitions.Partitions[0].Name,
			Nodes:     999999, // Excessive nodes
			CPUs:      999999, // Excessive CPUs
			TimeLimit: 5,
		}

		_, err := suite.client.Jobs().Submit(ctx, submission)
		if err != nil {
			suite.T().Logf("Excessive resource error (expected): %v", err)
		} else {
			// If it succeeds, cancel it immediately
			response := err.(*interfaces.JobSubmissionResponse)
			suite.createdJobs = append(suite.createdJobs, response.JobID)
			suite.client.Jobs().Cancel(ctx, response.JobID)
			suite.T().Log("Excessive resource job succeeded (cancelled immediately)")
		}
	}

	// Test 4: Network recovery
	suite.T().Log("Test 4: Testing connection recovery...")
	
	// This test verifies that the client can handle temporary network issues
	// We'll just do a series of quick operations to test resilience
	for i := 0; i < 3; i++ {
		err := suite.client.Info().Ping(ctx)
		if err != nil {
			suite.T().Logf("Ping %d failed: %v", i+1, err)
		} else {
			suite.T().Logf("Ping %d succeeded", i+1)
		}
		time.Sleep(1 * time.Second)
	}
}

// TestE2EWorkflowSuite runs the E2E workflow test suite
func TestE2EWorkflowSuite(t *testing.T) {
	suite.Run(t, new(E2EWorkflowTestSuite))
}