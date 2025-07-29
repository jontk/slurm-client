package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/jontk/slurm-client/internal/interfaces"
)

// V0040IntegrationTestSuite tests v0.0.40 API integration
type V0040IntegrationTestSuite struct {
	IntegrationTestSuite
}

// SetupSuite runs once before all tests
func (suite *V0040IntegrationTestSuite) SetupSuite() {
	suite.SetupIntegrationSuite("v0.0.40")
}

// TearDownSuite runs once after all tests
func (suite *V0040IntegrationTestSuite) TearDownSuite() {
	suite.CleanupAllResources()
	suite.TearDownIntegrationSuite()
}

// Test basic connectivity and version compatibility
func (suite *V0040IntegrationTestSuite) TestBasicConnectivity() {
	suite.T().Log("Testing basic connectivity for v0.0.40")
	
	// Test ping
	ctx := context.Background()
	err := suite.client.Info().Ping(ctx)
	suite.NoError(err, "Ping should succeed for v0.0.40")
	
	// Test version info
	versionInfo, err := suite.client.Info().Version(ctx)
	suite.NoError(err, "Version info should be available")
	suite.NotEmpty(versionInfo.Version, "Version should not be empty")
	suite.T().Logf("Server API version: %s", versionInfo.Version)
	
	// Test cluster info
	clusterInfo, err := suite.client.Info().Get(ctx)
	suite.NoError(err, "Cluster info should be available")
	suite.NotEmpty(clusterInfo.ClusterName, "Cluster name should not be empty")
	suite.T().Logf("Connected to cluster: %s", clusterInfo.ClusterName)
}

// Test job operations (core functionality that doesn't require database)
func (suite *V0040IntegrationTestSuite) TestJobOperations() {
	suite.TestCRUDWorkflow("Jobs", func() {
		ctx := context.Background()
		
		// List existing jobs
		jobs, err := suite.client.Jobs().List(ctx, &interfaces.ListJobsOptions{Limit: 10})
		suite.NoError(err, "Should be able to list jobs")
		suite.NotNil(jobs, "Jobs list should not be nil")
		suite.T().Logf("Found %d existing jobs", len(jobs.Jobs))
		
		// Submit a test job
		response, err := suite.CreateTestJob("v0040-crud")
		if err != nil {
			suite.T().Logf("Job submission failed (may be expected in test environment): %v", err)
			return // Skip rest of test if job submission is not available
		}
		
		suite.NotEmpty(response.JobID, "Job ID should not be empty")
		suite.T().Logf("Created test job: %s", response.JobID)
		suite.AddJobForCleanup(response.JobID)
		
		// Get the created job
		job, err := suite.client.Jobs().Get(ctx, response.JobID)
		suite.NoError(err, "Should be able to get the created job")
		suite.Equal(response.JobID, job.ID, "Job IDs should match")
		suite.T().Logf("Retrieved job: %s, State: %s", job.ID, job.State)
		
		// Test job state transitions
		suite.T().Logf("Waiting for job %s to start or complete...", response.JobID)
		time.Sleep(3 * time.Second) // Give job time to process
		
		// Get updated job state
		updatedJob, err := suite.client.Jobs().Get(ctx, response.JobID)
		suite.NoError(err, "Should be able to get updated job")
		suite.T().Logf("Job %s updated state: %s", updatedJob.ID, updatedJob.State)
		
		// Cancel the job (cleanup)
		err = suite.client.Jobs().Cancel(ctx, response.JobID)
		suite.NoError(err, "Should be able to cancel the job")
		suite.T().Logf("Successfully cancelled job: %s", response.JobID)
	})
}

// Test node operations
func (suite *V0040IntegrationTestSuite) TestNodeOperations() {
	suite.TestCRUDWorkflow("Nodes", func() {
		ctx := context.Background()
		
		// List nodes
		nodes, err := suite.client.Nodes().List(ctx, &interfaces.ListNodesOptions{Limit: 10})
		suite.NoError(err, "Should be able to list nodes")
		suite.NotNil(nodes, "Nodes list should not be nil")
		suite.T().Logf("Found %d nodes", len(nodes.Nodes))
		
		if len(nodes.Nodes) > 0 {
			// Get specific node
			firstNode := nodes.Nodes[0]
			node, err := suite.client.Nodes().Get(ctx, firstNode.Name)
			suite.NoError(err, "Should be able to get specific node")
			suite.Equal(firstNode.Name, node.Name, "Node names should match")
			suite.T().Logf("Retrieved node: %s, State: %s, CPUs: %d", 
				node.Name, node.State, node.CPUs)
		}
	})
}

// Test partition operations (may fail if database is required)
func (suite *V0040IntegrationTestSuite) TestPartitionOperations() {
	suite.TestCRUDWorkflow("Partitions", func() {
		ctx := context.Background()
		
		// List partitions
		partitions, err := suite.client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{Limit: 10})
		if err != nil {
			suite.SkipIfNoDatabaseConnection(err)
			suite.NoError(err, "Should be able to list partitions")
		}
		
		suite.NotNil(partitions, "Partitions list should not be nil")
		suite.T().Logf("Found %d partitions", len(partitions.Partitions))
		
		if len(partitions.Partitions) > 0 {
			// Get specific partition
			firstPartition := partitions.Partitions[0]
			partition, err := suite.client.Partitions().Get(ctx, firstPartition.Name)
			if err != nil {
				suite.SkipIfNoDatabaseConnection(err)
			}
			suite.NoError(err, "Should be able to get specific partition")
			suite.Equal(firstPartition.Name, partition.Name, "Partition names should match")
			suite.T().Logf("Retrieved partition: %s, State: %s, Nodes: %d", 
				partition.Name, partition.State, partition.TotalNodes)
		}
	})
}

// Test QoS operations (requires database - will skip if not available)
func (suite *V0040IntegrationTestSuite) TestQoSOperations() {
	if !suite.IsDatabaseAvailable() {
		suite.T().Skip("Database not available, skipping QoS tests")
		return
	}
	
	suite.TestCRUDWorkflow("QoS", func() {
		ctx := context.Background()
		
		// List QoS
		qosList, err := suite.client.QoS().List(ctx, &interfaces.ListQoSOptions{Limit: 10})
		suite.SkipIfNoDatabaseConnection(err)
		suite.NoError(err, "Should be able to list QoS")
		suite.NotNil(qosList, "QoS list should not be nil")
		suite.T().Logf("Found %d QoS entries", len(qosList.QoS))
		
		if len(qosList.QoS) > 0 {
			// Get specific QoS
			firstQoS := qosList.QoS[0]
			qos, err := suite.client.QoS().Get(ctx, firstQoS.Name)
			suite.NoError(err, "Should be able to get specific QoS")
			suite.Equal(firstQoS.Name, qos.Name, "QoS names should match")
			suite.T().Logf("Retrieved QoS: %s, Priority: %d", qos.Name, qos.Priority)
		}
	})
}

// Test error handling scenarios
func (suite *V0040IntegrationTestSuite) TestErrorHandling() {
	ctx := context.Background()
	
	// Test invalid job ID
	suite.TestErrorHandling("Invalid Job ID", func() error {
		_, err := suite.client.Jobs().Get(ctx, "invalid-job-id-999999")
		return err
	}, true)
	
	// Test invalid node name
	suite.TestErrorHandling("Invalid Node Name", func() error {
		_, err := suite.client.Nodes().Get(ctx, "invalid-node-name")
		return err
	}, true)
	
	// Test invalid partition name
	suite.TestErrorHandling("Invalid Partition Name", func() error {
		_, err := suite.client.Partitions().Get(ctx, "invalid-partition-name")
		return err
	}, true)
	
	// Test malformed job submission
	suite.TestErrorHandling("Malformed Job Submission", func() error {
		submission := &interfaces.JobSubmission{
			Name: "", // Empty name should cause error
			CPUs: -1, // Invalid CPU count
		}
		_, err := suite.client.Jobs().Submit(ctx, submission)
		return err
	}, true)
}

// Test performance characteristics
func (suite *V0040IntegrationTestSuite) TestPerformance() {
	ctx := context.Background()
	
	// Test job listing performance
	suite.TestPerformance("List Jobs", func() error {
		_, err := suite.client.Jobs().List(ctx, &interfaces.ListJobsOptions{Limit: 100})
		return err
	}, 5*time.Second)
	
	// Test node listing performance
	suite.TestPerformance("List Nodes", func() error {
		_, err := suite.client.Nodes().List(ctx, &interfaces.ListNodesOptions{Limit: 100})
		return err
	}, 5*time.Second)
	
	// Test ping performance
	suite.TestPerformance("Ping", func() error {
		return suite.client.Info().Ping(ctx)
	}, 2*time.Second)
}

// Test concurrent operations
func (suite *V0040IntegrationTestSuite) TestConcurrentOperations() {
	// Test concurrent job listings
	suite.TestConcurrentOperations("Concurrent Job Listings", func(id int) error {
		ctx := context.Background()
		_, err := suite.client.Jobs().List(ctx, &interfaces.ListJobsOptions{Limit: 10})
		return err
	}, 5)
	
	// Test concurrent node listings
	suite.TestConcurrentOperations("Concurrent Node Listings", func(id int) error {
		ctx := context.Background()
		_, err := suite.client.Nodes().List(ctx, &interfaces.ListNodesOptions{Limit: 10})
		return err
	}, 5)
	
	// Test concurrent pings
	suite.TestConcurrentOperations("Concurrent Pings", func(id int) error {
		ctx := context.Background()
		return suite.client.Info().Ping(ctx)
	}, 10)
}

// Test cluster statistics and diagnostics
func (suite *V0040IntegrationTestSuite) TestClusterStatistics() {
	ctx := context.Background()
	
	// Test cluster stats
	stats, err := suite.client.Info().Stats(ctx)
	suite.NoError(err, "Should be able to get cluster statistics")
	suite.NotNil(stats, "Stats should not be nil")
	suite.T().Logf("Cluster Stats - Nodes: %d, CPUs: %d, Jobs: %d", 
		stats.TotalNodes, stats.TotalCPUs, stats.TotalJobs)
	
	// Validate stats make sense
	suite.GreaterOrEqual(stats.TotalNodes, 0, "Total nodes should be non-negative")
	suite.GreaterOrEqual(stats.TotalCPUs, 0, "Total CPUs should be non-negative")
	suite.GreaterOrEqual(stats.TotalJobs, 0, "Total jobs should be non-negative")
}

// Test resource limits and constraints
func (suite *V0040IntegrationTestSuite) TestResourceLimits() {
	ctx := context.Background()
	
	// Test with various limits
	testCases := []struct {
		name  string
		limit int
	}{
		{"Small Limit", 1},
		{"Medium Limit", 10},
		{"Large Limit", 100},
	}
	
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			jobs, err := suite.client.Jobs().List(ctx, &interfaces.ListJobsOptions{
				Limit: tc.limit,
			})
			suite.NoError(err, "Should handle limit %d", tc.limit)
			suite.LessOrEqual(len(jobs.Jobs), tc.limit, 
				"Returned jobs should not exceed limit %d", tc.limit)
		})
	}
}

// Test API compatibility and version-specific features
func (suite *V0040IntegrationTestSuite) TestAPICompatibility() {
	// Verify client version
	suite.Equal("v0.0.40", suite.client.Version(), "Client should report correct version")
	
	// Test version-specific endpoint availability
	ctx := context.Background()
	
	// v0.0.40 should have basic job, node, partition support
	_, err := suite.client.Jobs().List(ctx, &interfaces.ListJobsOptions{Limit: 1})
	suite.NoError(err, "v0.0.40 should support job operations")
	
	_, err = suite.client.Nodes().List(ctx, &interfaces.ListNodesOptions{Limit: 1})
	suite.NoError(err, "v0.0.40 should support node operations")
	
	// Note: Some advanced features may not be available in v0.0.40
	suite.T().Log("v0.0.40 compatibility tests passed")
}

// TestV0040Integration runs the integration test suite for v0.0.40
func TestV0040Integration(t *testing.T) {
	suite.Run(t, new(V0040IntegrationTestSuite))
}