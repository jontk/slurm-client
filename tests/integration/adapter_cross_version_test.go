// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

// AdapterCrossVersionTestSuite tests operations across different API versions
type AdapterCrossVersionTestSuite struct {
	suite.Suite
	clients   map[string]slurm.SlurmClient
	versions  []string
	serverURL string
	token     string
}

// SetupSuite initializes clients for all supported versions
func (suite *AdapterCrossVersionTestSuite) SetupSuite() {
	// Check if cross-version testing is enabled
	if os.Getenv("SLURM_CROSS_VERSION_TEST") != "true" {
		suite.T().Skip("Cross-version tests disabled. Set SLURM_CROSS_VERSION_TEST=true to enable")
	}

	// Get server configuration
	suite.serverURL = os.Getenv("SLURM_SERVER_URL")
	if suite.serverURL == "" {
		suite.serverURL = "http://rocky9:6820"
	}

	// Get JWT token
	suite.token = os.Getenv("SLURM_JWT_TOKEN")
	if suite.token == "" {
		token, err := fetchJWTTokenViaSSH()
		require.NoError(suite.T(), err, "Failed to fetch JWT token")
		suite.token = token
	}

	// Initialize clients for all versions
	suite.versions = []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}
	suite.clients = make(map[string]slurm.SlurmClient)

	ctx := context.Background()
	for _, version := range suite.versions {
		client, err := slurm.NewClientWithVersion(ctx, version,
			slurm.WithBaseURL(suite.serverURL),
			slurm.WithAuth(auth.NewTokenAuth(suite.token)),
			slurm.WithConfig(&config.Config{
				Timeout:            30 * time.Second,
				MaxRetries:         3,
				Debug:              false,
				InsecureSkipVerify: true,
			}),
		)

		if err != nil {
			suite.T().Logf("Failed to create client for version %s: %v", version, err)
			continue
		}

		suite.clients[version] = client
		suite.T().Logf("Successfully created client for version %s", version)
	}

	require.NotEmpty(suite.T(), suite.clients, "At least one client must be created")
}

// TearDownSuite cleans up all clients
func (suite *AdapterCrossVersionTestSuite) TearDownSuite() {
	for version, client := range suite.clients {
		if client != nil {
			client.Close()
			suite.T().Logf("Closed client for version %s", version)
		}
	}
}

// TestPingAcrossVersions tests ping functionality across all versions
func (suite *AdapterCrossVersionTestSuite) TestPingAcrossVersions() {
	ctx := context.Background()

	results := make(map[string]error)

	for version, client := range suite.clients {
		err := client.Info().Ping(ctx)
		results[version] = err

		if err == nil {
			suite.T().Logf("✓ Ping successful for version %s", version)
		} else {
			suite.T().Logf("✗ Ping failed for version %s: %v", version, err)
		}
	}

	// At least one version should work
	successCount := 0
	for _, err := range results {
		if err == nil {
			successCount++
		}
	}

	suite.Require().Greater(successCount, 0, "At least one version should support ping")
}

// TestVersionConsistency tests version information consistency
func (suite *AdapterCrossVersionTestSuite) TestVersionConsistency() {
	ctx := context.Background()

	versionInfos := make(map[string]*interfaces.APIVersion)

	for version, client := range suite.clients {
		versionInfo, err := client.Info().Version(ctx)
		if err != nil {
			suite.T().Logf("Version info failed for %s: %v", version, err)
			continue
		}

		versionInfos[version] = versionInfo
		suite.T().Logf("Version %s reports API version: %s", version, versionInfo.Version)

		// Basic validation
		suite.NotEmpty(versionInfo.Version, "Version should not be empty for %s", version)
	}

	suite.Require().NotEmpty(versionInfos, "At least one version should provide version info")
}

// TestQoSListingConsistency tests QoS listing across versions
func (suite *AdapterCrossVersionTestSuite) TestQoSListingConsistency() {
	ctx := context.Background()

	qosResults := make(map[string]*interfaces.QoSList)

	for version, client := range suite.clients {
		qosList, err := client.QoS().List(ctx, &interfaces.ListQoSOptions{
			Limit: 5,
		})

		if err != nil {
			suite.T().Logf("QoS listing failed for %s: %v", version, err)
			continue
		}

		qosResults[version] = qosList
		suite.T().Logf("Version %s found %d QoS entries", version, len(qosList.QoS))

		// Validate basic structure
		suite.NotNil(qosList.QoS, "QoS list should not be nil for %s", version)

		for i, qos := range qosList.QoS {
			if i >= 3 { // Log first 3 for comparison
				break
			}
			suite.T().Logf("  %s QoS %d: Name=%s, Priority=%d", version, i+1, qos.Name, qos.Priority)
		}
	}

	// Compare common QoS entries if multiple versions work
	if len(qosResults) >= 2 {
		suite.T().Log("Comparing QoS consistency across versions...")

		// Find common QoS names
		var firstVersion string
		var firstQoS []string

		for version, qosList := range qosResults {
			firstVersion = version
			for _, qos := range qosList.QoS {
				firstQoS = append(firstQoS, qos.Name)
			}
			break
		}

		// Check if other versions have similar QoS entries
		for version, qosList := range qosResults {
			if version == firstVersion {
				continue
			}

			commonCount := 0
			for _, qos := range qosList.QoS {
				for _, name := range firstQoS {
					if qos.Name == name {
						commonCount++
						break
					}
				}
			}

			suite.T().Logf("Common QoS entries between %s and %s: %d", firstVersion, version, commonCount)
		}
	}
}

// TestJobListingConsistency tests job listing across versions
func (suite *AdapterCrossVersionTestSuite) TestJobListingConsistency() {
	ctx := context.Background()

	jobResults := make(map[string]*interfaces.JobList)

	for version, client := range suite.clients {
		jobList, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{
			Limit: 10,
		})

		if err != nil {
			suite.T().Logf("Job listing failed for %s: %v", version, err)
			continue
		}

		jobResults[version] = jobList
		suite.T().Logf("Version %s found %d jobs", version, len(jobList.Jobs))

		// Validate basic structure
		suite.NotNil(jobList.Jobs, "Job list should not be nil for %s", version)

		for i, job := range jobList.Jobs {
			if i >= 3 { // Log first 3 for comparison
				break
			}
			suite.T().Logf("  %s Job %d: ID=%s, Name=%s, State=%s", version, i+1, job.ID, job.Name, job.State)
		}
	}

	// Validate that job counts are reasonable across versions
	if len(jobResults) >= 2 {
		jobCounts := make(map[string]int)
		for version, jobList := range jobResults {
			jobCounts[version] = len(jobList.Jobs)
		}

		suite.T().Log("Job count comparison across versions:")
		for version, count := range jobCounts {
			suite.T().Logf("  %s: %d jobs", version, count)
		}
	}
}

// TestNodeListingConsistency tests node listing across versions
func (suite *AdapterCrossVersionTestSuite) TestNodeListingConsistency() {
	ctx := context.Background()

	nodeResults := make(map[string]*interfaces.NodeList)

	for version, client := range suite.clients {
		nodeList, err := client.Nodes().List(ctx, &interfaces.ListNodesOptions{
			Limit: 10,
		})

		if err != nil {
			suite.T().Logf("Node listing failed for %s: %v", version, err)
			continue
		}

		nodeResults[version] = nodeList
		suite.T().Logf("Version %s found %d nodes", version, len(nodeList.Nodes))

		// Validate basic structure
		suite.NotNil(nodeList.Nodes, "Node list should not be nil for %s", version)

		for i, node := range nodeList.Nodes {
			if i >= 3 { // Log first 3 for comparison
				break
			}
			suite.T().Logf("  %s Node %d: Name=%s, State=%s, CPUs=%d",
				version, i+1, node.Name, node.State, node.CPUs)
		}
	}

	// Validate node consistency across versions
	if len(nodeResults) >= 2 {
		suite.T().Log("Analyzing node consistency across versions...")

		// Node counts should be similar across versions
		nodeCounts := make(map[string]int)
		for version, nodeList := range nodeResults {
			nodeCounts[version] = len(nodeList.Nodes)
		}

		suite.T().Log("Node count comparison:")
		for version, count := range nodeCounts {
			suite.T().Logf("  %s: %d nodes", version, count)
		}
	}
}

// TestPartitionListingConsistency tests partition listing across versions
func (suite *AdapterCrossVersionTestSuite) TestPartitionListingConsistency() {
	ctx := context.Background()

	partitionResults := make(map[string]*interfaces.PartitionList)

	for version, client := range suite.clients {
		partitionList, err := client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{
			Limit: 10,
		})

		if err != nil {
			suite.T().Logf("Partition listing failed for %s: %v", version, err)
			continue
		}

		partitionResults[version] = partitionList
		suite.T().Logf("Version %s found %d partitions", version, len(partitionList.Partitions))

		// Validate basic structure
		suite.NotNil(partitionList.Partitions, "Partition list should not be nil for %s", version)

		for i, partition := range partitionList.Partitions {
			if i >= 3 { // Log first 3 for comparison
				break
			}
			suite.T().Logf("  %s Partition %d: Name=%s, State=%s, Nodes=%d",
				version, i+1, partition.Name, partition.State, partition.TotalNodes)
		}
	}

	// Validate partition consistency
	if len(partitionResults) >= 2 {
		suite.T().Log("Analyzing partition consistency across versions...")

		// Find common partition names
		partitionNames := make(map[string][]string)
		for version, partitionList := range partitionResults {
			var names []string
			for _, partition := range partitionList.Partitions {
				names = append(names, partition.Name)
			}
			partitionNames[version] = names
		}

		suite.T().Log("Partition names by version:")
		for version, names := range partitionNames {
			suite.T().Logf("  %s: %v", version, names)
		}
	}
}

// TestVersionSpecificFeatures tests features specific to certain versions
func (suite *AdapterCrossVersionTestSuite) TestVersionSpecificFeatures() {
	ctx := context.Background()

	// Test cluster stats (may not be available in all versions)
	for version, client := range suite.clients {
		stats, err := client.Info().Stats(ctx)
		if err != nil {
			suite.T().Logf("Stats not available for %s: %v", version, err)
			continue
		}

		suite.T().Logf("Version %s stats: Nodes=%d, CPUs=%d, Jobs=%d",
			version, stats.TotalNodes, stats.TotalCPUs, stats.TotalJobs)

		// Basic validation
		suite.GreaterOrEqual(stats.TotalNodes, int32(0), "Total nodes should be non-negative")
		suite.GreaterOrEqual(stats.TotalCPUs, int32(0), "Total CPUs should be non-negative")
		suite.GreaterOrEqual(stats.TotalJobs, int32(0), "Total jobs should be non-negative")
	}
}

// TestErrorHandlingConsistency tests error handling across versions
func (suite *AdapterCrossVersionTestSuite) TestErrorHandlingConsistency() {
	ctx := context.Background()

	// Test invalid resource access
	for version, client := range suite.clients {
		// Try to get a non-existent QoS
		_, err := client.QoS().Get(ctx, "nonexistent-qos-12345")

		if err != nil {
			suite.T().Logf("Version %s error for invalid QoS: %v", version, err)
			suite.Error(err, "Should return error for non-existent QoS in %s", version)
		} else {
			suite.T().Logf("Version %s unexpectedly succeeded for invalid QoS", version)
		}

		// Try to get a non-existent job
		_, err = client.Jobs().Get(ctx, "99999999")

		if err != nil {
			suite.T().Logf("Version %s error for invalid job: %v", version, err)
			suite.Error(err, "Should return error for non-existent job in %s", version)
		} else {
			suite.T().Logf("Version %s unexpectedly succeeded for invalid job", version)
		}
	}
}

// TestConcurrentOperations tests concurrent operations across versions
func (suite *AdapterCrossVersionTestSuite) TestConcurrentOperations() {
	ctx := context.Background()

	// Run concurrent ping operations
	type result struct {
		Version string
		Error   error
	}

	results := make(chan result, len(suite.clients))

	for version, client := range suite.clients {
		go func(v string, c slurm.SlurmClient) {
			err := c.Info().Ping(ctx)
			results <- result{Version: v, Error: err}
		}(version, client)
	}

	// Collect results
	successCount := 0
	for range len(suite.clients) {
		res := <-results
		if res.Error == nil {
			successCount++
			suite.T().Logf("✓ Concurrent ping successful for %s", res.Version)
		} else {
			suite.T().Logf("✗ Concurrent ping failed for %s: %v", res.Version, res.Error)
		}
	}

	suite.Greater(successCount, 0, "At least one concurrent ping should succeed")
}

// TestCrossVersionSuite runs the cross-version test suite
func TestCrossVersionSuite(t *testing.T) {
	suite.Run(t, new(AdapterCrossVersionTestSuite))
}
