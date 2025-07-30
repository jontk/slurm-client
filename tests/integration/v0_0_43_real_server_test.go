package integration

import (
	"context"
	goerrors "errors"
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

// V0043RealServerTestSuite tests v0.0.43 specific features against a real server
type V0043RealServerTestSuite struct {
	suite.Suite
	client    slurm.SlurmClient
	serverURL string
	token     string
	version   string
}

func (suite *V0043RealServerTestSuite) SetupSuite() {
	// Use the provided configuration
	suite.serverURL = "http://rocky9.ar.jontk.com:6820"
	suite.version = "v0.0.43"
	suite.token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI2NTM4Mjk5NzYsImlhdCI6MTc1MzgyOTk3Niwic3VuIjoicm9vdCJ9.-z8Cq_wHuOxNJ7KHHTboX3l9r6JBtSD1RxQUgQR9owE"

	// Create client
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
	require.NoError(suite.T(), err)
	suite.client = client

	suite.T().Logf("Testing against server: %s with API version: %s", suite.serverURL, suite.version)
}

func (suite *V0043RealServerTestSuite) TearDownSuite() {
	if suite.client != nil {
		suite.client.Close()
	}
}

// TestBasicConnectivity tests basic server connectivity
func (suite *V0043RealServerTestSuite) TestBasicConnectivity() {
	ctx := context.Background()

	suite.Run("Ping", func() {
		err := suite.client.Info().Ping(ctx)
		if err != nil {
			suite.logError("Ping", err)
		}
		suite.NoError(err, "Ping should succeed")
	})

	suite.Run("Version", func() {
		version, err := suite.client.Info().Version(ctx)
		if err != nil {
			suite.logError("Version", err)
		}
		suite.Require().NoError(err)
		suite.NotEmpty(version.Version)
		suite.T().Logf("Server version: %s", version.Version)
	})

	suite.Run("ClusterInfo", func() {
		info, err := suite.client.Info().Get(ctx)
		if err != nil {
			suite.logError("ClusterInfo", err)
		}
		suite.Require().NoError(err)
		suite.NotEmpty(info.ClusterName)
		suite.T().Logf("Cluster name: %s", info.ClusterName)
		suite.T().Logf("API version: %s", info.APIVersion)
	})

	suite.Run("Stats", func() {
		stats, err := suite.client.Info().Stats(ctx)
		if err != nil {
			suite.logError("Stats", err)
		}
		suite.Require().NoError(err)
		suite.T().Logf("Cluster stats: Total nodes=%d, Running jobs=%d", 
			stats.TotalNodes, stats.RunningJobs)
	})
}

// TestJobManager tests job management features
func (suite *V0043RealServerTestSuite) TestJobManager() {
	ctx := context.Background()

	suite.Run("ListJobs", func() {
		jobs, err := suite.client.Jobs().List(ctx, &interfaces.ListJobsOptions{
			Limit: 5,
		})
		if err != nil {
			suite.logError("ListJobs", err)
		}
		suite.Require().NoError(err)
		suite.NotNil(jobs)
		suite.T().Logf("Found %d jobs", len(jobs.Jobs))
	})

	suite.Run("SubmitAndCancelJob", func() {
		// Submit a test job
		submission := &interfaces.JobSubmission{
			Name:      fmt.Sprintf("v0043-test-%d", time.Now().Unix()),
			Script:    "#!/bin/bash\necho 'Testing v0.0.43'\nhostname\ndate\nsleep 10",
			Partition: "debug",
			Nodes:     1,
			CPUs:      1,
			TimeLimit: 5,
		}

		resp, err := suite.client.Jobs().Submit(ctx, submission)
		if err != nil {
			suite.logError("SubmitJob", err)
			// Check if it's a partition not found error
			var slurmErr *errors.SlurmError
			if goerrors.As(err, &slurmErr) {
				if slurmErr.Code == "INVALID_PARTITION" || slurmErr.Category == "InvalidRequest" {
					suite.T().Skip("Debug partition not available, skipping job submission test")
					return
				}
			}
		}
		suite.Require().NoError(err)
		suite.NotEmpty(resp.JobID)
		suite.T().Logf("Submitted job: %s", resp.JobID)

		// Get job details
		job, err := suite.client.Jobs().Get(ctx, resp.JobID)
		suite.Require().NoError(err)
		suite.Equal(resp.JobID, job.ID)

		// Cancel the job
		err = suite.client.Jobs().Cancel(ctx, resp.JobID)
		suite.NoError(err)
		suite.T().Logf("Cancelled job: %s", resp.JobID)
	})
}

// TestNodeManager tests node management
func (suite *V0043RealServerTestSuite) TestNodeManager() {
	ctx := context.Background()

	suite.Run("ListNodes", func() {
		nodes, err := suite.client.Nodes().List(ctx, &interfaces.ListNodesOptions{
			Limit: 5,
		})
		if err != nil {
			suite.logError("ListNodes", err)
		}
		suite.Require().NoError(err)
		suite.NotNil(nodes)
		suite.T().Logf("Found %d nodes", len(nodes.Nodes))
		
		for i, node := range nodes.Nodes {
			if i < 3 {
				suite.T().Logf("  Node: %s, State: %s, CPUs: %d", 
					node.Name, node.State, node.CPUs)
			}
		}
	})

	suite.Run("GetNode", func() {
		// First list to get a valid node name
		nodes, err := suite.client.Nodes().List(ctx, &interfaces.ListNodesOptions{
			Limit: 1,
		})
		suite.Require().NoError(err)
		
		if len(nodes.Nodes) == 0 {
			suite.T().Skip("No nodes found")
			return
		}

		nodeName := nodes.Nodes[0].Name
		node, err := suite.client.Nodes().Get(ctx, nodeName)
		if err != nil {
			suite.logError("GetNode", err)
		}
		suite.Require().NoError(err)
		suite.Equal(nodeName, node.Name)
		suite.T().Logf("Retrieved node: %s, State: %s", node.Name, node.State)
	})
}

// TestPartitionManager tests partition management
func (suite *V0043RealServerTestSuite) TestPartitionManager() {
	ctx := context.Background()

	suite.Run("ListPartitions", func() {
		partitions, err := suite.client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{
			Limit: 10,
		})
		if err != nil {
			suite.logError("ListPartitions", err)
		}
		suite.Require().NoError(err)
		suite.NotNil(partitions)
		suite.T().Logf("Found %d partitions", len(partitions.Partitions))
		
		for i, partition := range partitions.Partitions {
			if i < 3 {
				suite.T().Logf("  Partition: %s, State: %s, Nodes: %d", 
					partition.Name, partition.State, partition.TotalNodes)
			}
		}
	})
}

// TestQoSManager tests QoS management (requires slurmdbd)
func (suite *V0043RealServerTestSuite) TestQoSManager() {
	ctx := context.Background()

	suite.Run("ListQoS", func() {
		qosList, err := suite.client.QoS().List(ctx, &interfaces.ListQoSOptions{
			Limit: 10,
		})
		
		if err != nil {
			suite.logError("ListQoS", err)
			// Check if it's a database connection error
			if suite.isDatabaseError(err) {
				suite.T().Skip("QoS operations require slurmdbd connection")
				return
			}
		}
		
		suite.Require().NoError(err)
		suite.NotNil(qosList)
		suite.T().Logf("Found %d QoS entries", len(qosList.QoS))
	})
}

// TestUserAccountManagement tests user and account features (requires slurmdbd)
func (suite *V0043RealServerTestSuite) TestUserAccountManagement() {
	ctx := context.Background()

	suite.Run("ListUsers", func() {
		users, err := suite.client.Users().List(ctx, &interfaces.ListUsersOptions{
			Limit: 5,
		})
		
		if err != nil {
			suite.logError("ListUsers", err)
			if suite.isDatabaseError(err) {
				suite.T().Skip("User operations require slurmdbd connection")
				return
			}
		}
		
		suite.Require().NoError(err)
		suite.NotNil(users)
		suite.T().Logf("Found %d users", len(users.Users))
	})

	suite.Run("ListAccounts", func() {
		accounts, err := suite.client.Accounts().List(ctx, &interfaces.ListAccountsOptions{
			Limit: 5,
		})
		
		if err != nil {
			suite.logError("ListAccounts", err)
			if suite.isDatabaseError(err) {
				suite.T().Skip("Account operations require slurmdbd connection")
				return
			}
		}
		
		suite.Require().NoError(err)
		suite.NotNil(accounts)
		suite.T().Logf("Found %d accounts", len(accounts.Accounts))
	})
}

// TestV0043SpecificFeatures tests features specific to v0.0.43
func (suite *V0043RealServerTestSuite) TestV0043SpecificFeatures() {
	ctx := context.Background()

	// Test Association Manager (v0.0.43 specific)
	suite.Run("AssociationManager", func() {
		if suite.client.Associations() == nil {
			suite.T().Skip("Association manager not available")
			return
		}

		associations, err := suite.client.Associations().List(ctx, &interfaces.ListAssociationsOptions{
			Limit: 5,
		})
		
		if err != nil {
			suite.logError("ListAssociations", err)
			if suite.isDatabaseError(err) {
				suite.T().Skip("Association operations require slurmdbd connection")
				return
			}
		}
		
		suite.Require().NoError(err)
		suite.NotNil(associations)
		suite.T().Logf("Found %d associations", len(associations.Associations))
	})

	// Test Cluster Manager (v0.0.43 specific)
	suite.Run("ClusterManager", func() {
		if suite.client.Clusters() == nil {
			suite.T().Skip("Cluster manager not available")
			return
		}

		clusters, err := suite.client.Clusters().List(ctx, nil)
		
		if err != nil {
			suite.logError("ListClusters", err)
			if suite.isDatabaseError(err) {
				suite.T().Skip("Cluster operations require slurmdbd connection")
				return
			}
		}
		
		suite.Require().NoError(err)
		suite.NotNil(clusters)
		suite.T().Logf("Found %d clusters", len(clusters.Clusters))
	})
}

// TestReservationManager tests reservation management
func (suite *V0043RealServerTestSuite) TestReservationManager() {
	ctx := context.Background()

	suite.Run("ListReservations", func() {
		reservations, err := suite.client.Reservations().List(ctx, &interfaces.ListReservationsOptions{
			Limit: 5,
		})
		if err != nil {
			suite.logError("ListReservations", err)
		}
		suite.Require().NoError(err)
		suite.NotNil(reservations)
		suite.T().Logf("Found %d reservations", len(reservations.Reservations))
	})
}

// Helper methods

func (suite *V0043RealServerTestSuite) logError(operation string, err error) {
	suite.T().Logf("Error in %s: %v", operation, err)
	
	var slurmErr *errors.SlurmError
	if goerrors.As(err, &slurmErr) {
		suite.T().Logf("  Error Code: %s", slurmErr.Code)
		suite.T().Logf("  Error Category: %s", slurmErr.Category)
		suite.T().Logf("  Error Message: %s", slurmErr.Message)
		suite.T().Logf("  Error Details: %s", slurmErr.Details)
		suite.T().Logf("  Status Code: %d", slurmErr.StatusCode)
	}
	
	var apiErr *errors.SlurmAPIError
	if goerrors.As(err, &apiErr) {
		suite.T().Logf("  API Error Number: %d", apiErr.ErrorNumber)
		suite.T().Logf("  API Error Code: %s", apiErr.ErrorCode)
		suite.T().Logf("  API Error Source: %s", apiErr.Source)
		if len(apiErr.Errors) > 0 {
			for i, detail := range apiErr.Errors {
				suite.T().Logf("  API Error Detail %d: [%d] %s - %s (source: %s)", 
					i+1, detail.ErrorNumber, detail.ErrorCode, detail.Description, detail.Source)
			}
		}
	}
}

func (suite *V0043RealServerTestSuite) isDatabaseError(err error) bool {
	if err == nil {
		return false
	}
	
	errStr := err.Error()
	return v0043Contains(errStr, "Unable to connect to database") ||
		v0043Contains(errStr, "Failed to open slurmdbd connection") ||
		v0043Contains(errStr, "SLURM_DAEMON_DOWN") ||
		v0043Contains(errStr, "slurmdbd") ||
		v0043Contains(errStr, "database")
}

func v0043Contains(str, substr string) bool {
	return len(substr) > 0 && len(str) >= len(substr) && 
		(str == substr || len(str) > len(substr) && 
		(str[:len(substr)] == substr || str[len(str)-len(substr):] == substr ||
		v0043FindSubstring(str, substr) >= 0))
}

func v0043FindSubstring(str, substr string) int {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// TestV0043RealServer runs the v0.0.43 specific test suite
func TestV0043RealServer(t *testing.T) {
	if os.Getenv("SLURM_REAL_SERVER_TEST") != "true" {
		t.Skip("Real server tests disabled. Set SLURM_REAL_SERVER_TEST=true to enable")
	}
	suite.Run(t, new(V0043RealServerTestSuite))
}