// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	goerrors "errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
	"github.com/jontk/slurm-client/pkg/errors"
)

// RealServerTestSuite tests against a real slurmrestd server
type RealServerTestSuite struct {
	suite.Suite
	client    slurm.SlurmClient
	serverURL string
	token     string
	version   string
}

// SetupSuite runs once before all tests
func (suite *RealServerTestSuite) SetupSuite() {
	// Check if real server testing is enabled
	if os.Getenv("SLURM_REAL_SERVER_TEST") != "true" {
		suite.T().Skip("Real server tests disabled. Set SLURM_REAL_SERVER_TEST=true to enable")
	}

	// Get configuration from environment
	suite.serverURL = os.Getenv("SLURM_SERVER_URL")
	if suite.serverURL == "" {
		suite.serverURL = "http://rocky9:6820"
	}

	// Get API version from environment or use default
	suite.version = os.Getenv("SLURM_API_VERSION")
	if suite.version == "" {
		suite.version = "v0.0.43" // SLURM 25.05 latest API version
	}

	// Get token via SSH if not provided
	suite.token = os.Getenv("SLURM_JWT_TOKEN")
	if suite.token == "" {
		token, err := fetchJWTTokenViaSSH()
		require.NoError(suite.T(), err, "Failed to fetch JWT token via SSH")
		suite.token = token
		suite.T().Logf("Fetched JWT token: %s...", suite.token[:50])
	}

	// Create client
	ctx := context.Background()
	client, err := slurm.NewClientWithVersion(ctx, suite.version,
		slurm.WithBaseURL(suite.serverURL),
		slurm.WithAuth(auth.NewTokenAuth(suite.token)),
		slurm.WithConfig(&config.Config{
			Timeout:            30 * time.Second,
			MaxRetries:         3,
			Debug:              true,
			InsecureSkipVerify: true, // For test servers with self-signed certs
		}),
	)
	require.NoError(suite.T(), err)
	suite.client = client
}

// TearDownSuite runs once after all tests
func (suite *RealServerTestSuite) TearDownSuite() {
	if suite.client != nil {
		suite.client.Close()
	}
}

// TestPing verifies connectivity and authentication
func (suite *RealServerTestSuite) TestPing() {
	ctx := context.Background()
	err := suite.client.Info().Ping(ctx)
	suite.NoError(err, "Ping should succeed")
}

// TestGetClusterInfo tests retrieving cluster information
func (suite *RealServerTestSuite) TestGetClusterInfo() {
	ctx := context.Background()

	// Get cluster info
	info, err := suite.client.Info().Get(ctx)
	suite.Require().NoError(err)
	suite.NotEmpty(info.ClusterName, "Cluster name should not be empty")

	suite.T().Logf("Connected to cluster: %s", info.ClusterName)
}

// TestGetVersion tests version endpoint
func (suite *RealServerTestSuite) TestGetVersion() {
	ctx := context.Background()

	versionInfo, err := suite.client.Info().Version(ctx)
	suite.Require().NoError(err)
	suite.NotEmpty(versionInfo.Version, "Version should not be empty")

	suite.T().Logf("Server API version: %s", versionInfo.Version)
}

// TestListJobs tests job listing
func (suite *RealServerTestSuite) TestListJobs() {
	ctx := context.Background()

	// List jobs with limit
	jobs, err := suite.client.Jobs().List(ctx, &interfaces.ListJobsOptions{
		Limit: 10,
	})
	suite.Require().NoError(err)
	suite.NotNil(jobs)

	suite.T().Logf("Found %d jobs", len(jobs.Jobs))
	for i, job := range jobs.Jobs {
		if i < 5 { // Log first 5 jobs
			suite.T().Logf("  Job %d: ID=%s, Name=%s, State=%s", i+1, job.ID, job.Name, job.State)
		}
	}
}

// TestListNodes tests node listing
func (suite *RealServerTestSuite) TestListNodes() {
	ctx := context.Background()

	// List nodes
	nodes, err := suite.client.Nodes().List(ctx, &interfaces.ListNodesOptions{
		Limit: 10,
	})
	suite.Require().NoError(err)
	suite.NotNil(nodes)

	suite.T().Logf("Found %d nodes", len(nodes.Nodes))
	for i, node := range nodes.Nodes {
		if i < 5 { // Log first 5 nodes
			suite.T().Logf("  Node %d: Name=%s, State=%s, CPUs=%d", i+1, node.Name, node.State, node.CPUs)
		}
	}
}

// TestListPartitions tests partition listing
func (suite *RealServerTestSuite) TestListPartitions() {
	ctx := context.Background()

	// List partitions
	partitions, err := suite.client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{
		Limit: 10,
	})
	suite.Require().NoError(err)
	suite.NotNil(partitions)

	suite.T().Logf("Found %d partitions", len(partitions.Partitions))
	for i, partition := range partitions.Partitions {
		if i < 5 { // Log first 5 partitions
			suite.T().Logf("  Partition %d: Name=%s, State=%s, Nodes=%d",
				i+1, partition.Name, partition.State, partition.TotalNodes)
		}
	}
}

// TestListQoS tests QoS listing with our new adapter pattern
func (suite *RealServerTestSuite) TestListQoS() {
	ctx := context.Background()

	// List QoS entries
	qosList, err := suite.client.QoS().List(ctx, &interfaces.ListQoSOptions{
		Limit: 10,
	})
	suite.Require().NoError(err)
	suite.NotNil(qosList)

	suite.T().Logf("Found %d QoS entries", len(qosList.QoS))
	for i, qos := range qosList.QoS {
		if i < 5 { // Log first 5 QoS entries
			suite.T().Logf("  QoS %d: Name=%s, Priority=%d, UsageFactor=%.2f",
				i+1, qos.Name, qos.Priority, qos.UsageFactor)
			suite.T().Logf("    Grace Time: %d seconds", qos.GraceTime)
			if qos.MaxJobsPerUser > 0 {
				suite.T().Logf("    Max Jobs Per User: %d", qos.MaxJobsPerUser)
			}
			if qos.MaxJobsPerAccount > 0 {
				suite.T().Logf("    Max Jobs Per Account: %d", qos.MaxJobsPerAccount)
			}
		}
	}
}

// TestGetQoS tests retrieving a specific QoS
func (suite *RealServerTestSuite) TestGetQoS() {
	ctx := context.Background()

	// First list QoS to get a valid name
	qosList, err := suite.client.QoS().List(ctx, &interfaces.ListQoSOptions{
		Limit: 1,
	})
	suite.Require().NoError(err)

	if len(qosList.QoS) == 0 {
		suite.T().Skip("No QoS entries found to test")
		return
	}

	qosName := qosList.QoS[0].Name
	suite.T().Logf("Testing Get with QoS: %s", qosName)

	// Get specific QoS
	qos, err := suite.client.QoS().Get(ctx, qosName)
	suite.Require().NoError(err)
	suite.NotNil(qos)
	suite.Equal(qosName, qos.Name)

	suite.T().Logf("Retrieved QoS: %s", qos.Name)
	suite.T().Logf("  Description: %s", qos.Description)
	suite.T().Logf("  Priority: %d", qos.Priority)
	suite.T().Logf("  UsageFactor: %.2f", qos.UsageFactor)
	suite.T().Logf("  GraceTime: %d", qos.GraceTime)
}

// TestJobSubmission tests submitting and canceling a job
func (suite *RealServerTestSuite) TestJobSubmission() {
	ctx := context.Background()

	// Submit a test job
	submission := &interfaces.JobSubmission{
		Name:      "go-client-test-" + time.Now().Format("20060102-150405"),
		Script:    "#!/bin/bash\necho 'Hello from Go SLURM client test'\nhostname\ndate\nsleep 30",
		Partition: "debug", // Using debug partition which exists on the test server
		Nodes:     1,
		CPUs:      1,
		TimeLimit: 5, // 5 minutes
	}

	suite.T().Logf("Submitting job: %s", submission.Name)
	response, err := suite.client.Jobs().Submit(ctx, submission)
	if err != nil {
		suite.T().Logf("Job submission error: %v", err)
		// Import the errors package if not already imported
		var slurmErr *errors.SlurmError
		if goerrors.As(err, &slurmErr) {
			suite.T().Logf("Error Code: %s", slurmErr.Code)
			suite.T().Logf("Error Category: %s", slurmErr.Category)
			suite.T().Logf("Error Message: %s", slurmErr.Message)
			suite.T().Logf("Error Details: %s", slurmErr.Details)
			suite.T().Logf("Status Code: %d", slurmErr.StatusCode)
		}

		// Check if it's a SlurmAPIError with more details
		var apiErr *errors.SlurmAPIError
		if goerrors.As(err, &apiErr) {
			suite.T().Logf("API Error Number: %d", apiErr.ErrorNumber)
			suite.T().Logf("API Error Code: %s", apiErr.ErrorCode)
			suite.T().Logf("API Error Source: %s", apiErr.Source)
			if len(apiErr.Errors) > 0 {
				for i, detail := range apiErr.Errors {
					suite.T().Logf("API Error Detail %d: [%d] %s - %s (source: %s)",
						i+1, detail.ErrorNumber, detail.ErrorCode, detail.Description, detail.Source)
				}
			}
		}
	}
	suite.Require().NoError(err)
	suite.NotEmpty(response.JobID)

	suite.T().Logf("Job submitted successfully: ID=%s", response.JobID)

	// Get job details
	job, err := suite.client.Jobs().Get(ctx, response.JobID)
	suite.Require().NoError(err)
	suite.Equal(response.JobID, job.ID)
	suite.T().Logf("Job state: %s", job.State)

	// Cancel the job
	err = suite.client.Jobs().Cancel(ctx, response.JobID)
	suite.NoError(err)
	suite.T().Logf("Job cancelled successfully")
}

// TestGetStats tests retrieving cluster statistics
func (suite *RealServerTestSuite) TestGetStats() {
	ctx := context.Background()

	stats, err := suite.client.Info().Stats(ctx)
	suite.Require().NoError(err)

	suite.T().Logf("Cluster Statistics:")
	suite.T().Logf("  Node Statistics:")
	suite.T().Logf("    Total Nodes: %d", stats.TotalNodes)
	suite.T().Logf("    Idle Nodes: %d", stats.IdleNodes)
	suite.T().Logf("    Allocated Nodes: %d", stats.AllocatedNodes)
	suite.T().Logf("  CPU Statistics:")
	suite.T().Logf("    Total CPUs: %d", stats.TotalCPUs)
	suite.T().Logf("    Idle CPUs: %d", stats.IdleCPUs)
	suite.T().Logf("    Allocated CPUs: %d", stats.AllocatedCPUs)
	suite.T().Logf("  Job Statistics:")
	suite.T().Logf("    Total Jobs: %d", stats.TotalJobs)
	suite.T().Logf("    Running Jobs: %d", stats.RunningJobs)
	suite.T().Logf("    Pending Jobs: %d", stats.PendingJobs)
	suite.T().Logf("    Completed Jobs: %d", stats.CompletedJobs)
}

// TestRealServer runs the test suite
func TestRealServer(t *testing.T) {
	suite.Run(t, new(RealServerTestSuite))
}
