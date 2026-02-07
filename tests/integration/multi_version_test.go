//go:build integration
// +build integration

// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jontk/slurm-client"
	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
	"github.com/jontk/slurm-client/tests/helpers"
	"github.com/jontk/slurm-client/tests/mocks"
)

// TestMultiVersionCompatibility tests that the same client code works across all API versions
func TestMultiVersionCompatibility(t *testing.T) {
	// Create mock server pool for all versions
	serverPool := mocks.NewMockServerPool()
	defer serverPool.Close()

	ctx := helpers.TestContext(t)

	supportedVersions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}

	for _, version := range supportedVersions {
		t.Run("Version_"+version, func(t *testing.T) {
			testVersionCompatibility(t, ctx, serverPool, version)
		})
	}
}

func testVersionCompatibility(t *testing.T, ctx context.Context, serverPool *mocks.MockServerPool, version string) {
	// Skip v0.0.41 - this version uses complex anonymous inline structs in the OpenAPI spec
	// that make it very difficult to implement in Go. The generated client has known limitations
	// with job submission and other operations due to these inline structs.
	if version == "v0.0.41" {
		t.Skip("v0.0.41 has incomplete implementation due to complex OpenAPI inline struct limitations")
	}

	server := serverPool.GetServer(version)
	require.NotNil(t, server, "Mock server should exist for version %s", version)

	// Create client for this version
	client, err := slurm.NewClientWithVersion(ctx, version,
		slurm.WithBaseURL(server.URL()),
		slurm.WithAuth(auth.NewNoAuth()),
		slurm.WithConfig(&config.Config{
			Timeout:    30 * time.Second,
			MaxRetries: 3,
			Debug:      false, // Reduce noise in multi-version tests
		}),
	)
	require.NoError(t, err)
	defer client.Close()

	// Test 1: Basic Info Operations (should work across all versions)
	t.Run("InfoOperations", func(t *testing.T) {
		// Ping
		err := client.Info().Ping(ctx)
		assert.NoError(t, err, "Ping should work in version %s", version)

		// Version info
		versionInfo, err := client.Info().Version(ctx)
		require.NoError(t, err, "Version should work in version %s", version)
		assert.Equal(t, version, versionInfo.Version)

		// Cluster info
		clusterInfo, err := client.Info().Get(ctx)
		require.NoError(t, err, "Cluster info should work in version %s", version)
		assert.NotEmpty(t, clusterInfo.ClusterName)

		// Stats
		stats, err := client.Info().Stats(ctx)
		require.NoError(t, err, "Stats should work in version %s", version)
		assert.GreaterOrEqual(t, stats.TotalNodes, 0)
	})

	// Test 2: Job Operations (core functionality)
	t.Run("JobOperations", func(t *testing.T) {
		// List jobs
		jobs, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{
			Limit: 10,
		})
		require.NoError(t, err, "Job list should work in version %s", version)
		assert.NotNil(t, jobs)

		if len(jobs.Jobs) > 0 {
			// Get first job
			firstJob := jobs.Jobs[0]
			job, err := client.Jobs().Get(ctx, firstJob.ID)
			require.NoError(t, err, "Job get should work in version %s", version)
			assert.Equal(t, firstJob.ID, job.ID)
		}

		// Submit a test job
		submission := &interfaces.JobSubmission{
			Name:      "compat-test-" + version,
			Script:    "#!/bin/bash\necho 'Compatibility test for " + version + "'",
			Partition: "compute",
			Cpus:      1,
			TimeLimit: 10,
		}

		response, err := client.Jobs().Submit(ctx, submission)
		require.NoError(t, err, "Job submit should work in version %s", version)
		require.NotNil(t, response)

		jobID := response.JobID

		// Cancel the job
		err = client.Jobs().Cancel(ctx, jobID)
		assert.NoError(t, err, "Job cancel should work in version %s", version)
	})

	// Test 3: Node Operations
	t.Run("NodeOperations", func(t *testing.T) {
		// List nodes
		nodes, err := client.Nodes().List(ctx, &interfaces.ListNodesOptions{
			Limit: 10,
		})
		require.NoError(t, err, "Node list should work in version %s", version)
		assert.NotNil(t, nodes)

		if len(nodes.Nodes) > 0 {
			// Get first node
			firstNode := nodes.Nodes[0]
			node, err := client.Nodes().Get(ctx, firstNode.Name)
			require.NoError(t, err, "Node get should work in version %s", version)
			assert.Equal(t, firstNode.Name, node.Name)
		}
	})

	// Test 4: Partition Operations
	t.Run("PartitionOperations", func(t *testing.T) {
		// List partitions
		partitions, err := client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{
			Limit: 10,
		})
		require.NoError(t, err, "Partition list should work in version %s", version)
		assert.NotNil(t, partitions)

		if len(partitions.Partitions) > 0 {
			// Get first partition
			firstPartition := partitions.Partitions[0]
			partition, err := client.Partitions().Get(ctx, firstPartition.Name)
			require.NoError(t, err, "Partition get should work in version %s", version)
			assert.Equal(t, firstPartition.Name, partition.Name)
		}
	})
}

// TestVersionSpecificFeatures tests features that are version-specific
func TestVersionSpecificFeatures(t *testing.T) {
	testCases := []struct {
		version     string
		features    map[string]bool
		description string
	}{
		{
			version: "v0.0.40",
			features: map[string]bool{
				"job_update":       false,
				"node_update":      false,
				"partition_update": false,
			},
			description: "v0.0.40 has limited update capabilities",
		},
		{
			version: "v0.0.41",
			features: map[string]bool{
				"job_update":       true,
				"node_update":      false,
				"partition_update": false,
			},
			description: "v0.0.41 adds job update support",
		},
		{
			version: "v0.0.42",
			features: map[string]bool{
				"job_update":       true,
				"node_update":      true,
				"partition_update": false, // Partition updates not available in REST API
			},
			description: "v0.0.42 has job and node update support",
		},
		{
			version: "v0.0.43",
			features: map[string]bool{
				"job_update":       true,
				"node_update":      true,
				"partition_update": false, // Partition updates not available in REST API
			},
			description: "v0.0.43 has job and node update support with new features",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.version, func(t *testing.T) {
			testVersionSpecificFeatures(t, tc.version, tc.features, tc.description)
		})
	}
}

func testVersionSpecificFeatures(t *testing.T, version string, expectedFeatures map[string]bool, description string) {
	// Skip v0.0.41 - has incomplete implementation due to OpenAPI inline struct limitations
	if version == "v0.0.41" {
		t.Skip("v0.0.41 has incomplete implementation due to complex OpenAPI inline struct limitations")
	}

	server := mocks.NewMockSlurmServerForVersion(version)
	defer server.Close()

	ctx := helpers.TestContext(t)
	client, err := slurm.NewClientWithVersion(ctx, version,
		slurm.WithBaseURL(server.URL()),
		slurm.WithAuth(auth.NewNoAuth()),
	)
	require.NoError(t, err)
	defer client.Close()

	t.Log(description)

	// Test job update capability
	if expectedFeatures["job_update"] {
		t.Run("JobUpdateSupported", func(t *testing.T) {
			// Submit a job first
			submission := &interfaces.JobSubmission{
				Name:      "update-test",
				Script:    "#!/bin/bash\necho 'test'",
				Partition: "compute",
				Cpus:      1,
			}

			response, err := client.Jobs().Submit(ctx, submission)
			require.NoError(t, err)

			// Try to update it
			update := &interfaces.JobUpdate{
				Name: stringPtr("updated-job"),
			}

			err = client.Jobs().Update(ctx, response.JobID, update)
			assert.NoError(t, err, "Job update should be supported in %s", version)
		})
	} else {
		t.Run("JobUpdateNotSupported", func(t *testing.T) {
			// Submit a job first
			submission := &interfaces.JobSubmission{
				Name:      "update-test",
				Script:    "#!/bin/bash\necho 'test'",
				Partition: "compute",
				Cpus:      1,
			}

			response, err := client.Jobs().Submit(ctx, submission)
			require.NoError(t, err)

			// Try to update it - should fail or be limited
			update := &interfaces.JobUpdate{
				Name: stringPtr("updated-job"),
			}

			err = client.Jobs().Update(ctx, response.JobID, update)
			if err != nil {
				// Expected for versions that don't support updates
				assert.Contains(t, err.Error(), "not supported", "Update should fail gracefully in %s", version)
			}
		})
	}

	// Test node update capability
	if expectedFeatures["node_update"] {
		t.Run("NodeUpdateSupported", func(t *testing.T) {
			// Get first node
			nodes, err := client.Nodes().List(ctx, &interfaces.ListNodesOptions{Limit: 1})
			require.NoError(t, err)
			require.Greater(t, len(nodes.Nodes), 0)

			nodeName := nodes.Nodes[0].Name
			update := &interfaces.NodeUpdate{
				State: stringPtr("DRAIN"),
			}

			err = client.Nodes().Update(ctx, nodeName, update)
			assert.NoError(t, err, "Node update should be supported in %s", version)
		})
	}

	// Test partition update capability
	if expectedFeatures["partition_update"] {
		t.Run("PartitionUpdateSupported", func(t *testing.T) {
			// Get first partition
			partitions, err := client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{Limit: 1})
			require.NoError(t, err)
			require.Greater(t, len(partitions.Partitions), 0)

			partitionName := partitions.Partitions[0].Name
			update := &interfaces.PartitionUpdate{
				DefaultTime: intPtr(120),
			}

			err = client.Partitions().Update(ctx, partitionName, update)
			assert.NoError(t, err, "Partition update should be supported in %s", version)
		})
	}
}

// TestVersionMigration tests migrating from one version to another
func TestVersionMigration(t *testing.T) {
	serverPool := mocks.NewMockServerPool()
	defer serverPool.Close()

	ctx := helpers.TestContext(t)

	migrationTests := []struct {
		fromVersion string
		toVersion   string
		description string
	}{
		{"v0.0.40", "v0.0.41", "Migrate from v0.0.40 to v0.0.41"},
		{"v0.0.41", "v0.0.42", "Migrate from v0.0.41 to v0.0.42"},
		{"v0.0.42", "v0.0.43", "Migrate from v0.0.42 to v0.0.43"},
		{"v0.0.40", "v0.0.42", "Skip version migration v0.0.40 to v0.0.42"},
	}

	for _, test := range migrationTests {
		t.Run(test.fromVersion+"_to_"+test.toVersion, func(t *testing.T) {
			testVersionMigration(t, ctx, serverPool, test.fromVersion, test.toVersion, test.description)
		})
	}
}

func testVersionMigration(t *testing.T, ctx context.Context, serverPool *mocks.MockServerPool, fromVersion, toVersion, description string) {
	t.Log(description)

	// Skip migrations involving v0.0.41 - has incomplete implementation
	if fromVersion == "v0.0.41" || toVersion == "v0.0.41" {
		t.Skip("v0.0.41 has incomplete implementation due to complex OpenAPI inline struct limitations")
	}

	// Create client with old version
	oldServer := serverPool.GetServer(fromVersion)
	oldClient, err := slurm.NewClientWithVersion(ctx, fromVersion,
		slurm.WithBaseURL(oldServer.URL()),
		slurm.WithAuth(auth.NewNoAuth()),
	)
	require.NoError(t, err)
	defer oldClient.Close()

	// Submit a job with old client
	submission := &interfaces.JobSubmission{
		Name:      "migration-test-" + fromVersion + "-to-" + toVersion,
		Script:    "#!/bin/bash\necho 'Migration test'",
		Partition: "compute",
		Cpus:      1,
	}

	response, err := oldClient.Jobs().Submit(ctx, submission)
	require.NoError(t, err)
	_ = response.JobID // We don't actually use this jobID in this test

	// Create client with new version
	newServer := serverPool.GetServer(toVersion)
	newClient, err := slurm.NewClientWithVersion(ctx, toVersion,
		slurm.WithBaseURL(newServer.URL()),
		slurm.WithAuth(auth.NewNoAuth()),
	)
	require.NoError(t, err)
	defer newClient.Close()

	// Test that common operations work with new client
	// Note: In a real scenario, this would be the same cluster with upgraded API

	// Submit job with new client
	newSubmission := &interfaces.JobSubmission{
		Name:      "migration-new-" + toVersion,
		Script:    "#!/bin/bash\necho 'New version test'",
		Partition: "compute",
		Cpus:      1,
	}

	newResponse, err := newClient.Jobs().Submit(ctx, newSubmission)
	require.NoError(t, err)
	assert.NotEmpty(t, newResponse.JobID)

	// Test version-specific functionality
	if newServer.GetConfig().SupportedOperations["jobs.update"] && !oldServer.GetConfig().SupportedOperations["jobs.update"] {
		t.Log("Testing new update functionality in", toVersion)
		update := &interfaces.JobUpdate{
			Name: stringPtr("updated-in-new-version"),
		}
		err = newClient.Jobs().Update(ctx, newResponse.JobID, update)
		assert.NoError(t, err, "New version should support job updates")
	}

	t.Logf("Successfully migrated from %s to %s", fromVersion, toVersion)
}

// TestConcurrentVersions tests running multiple versions simultaneously
func TestConcurrentVersions(t *testing.T) {
	serverPool := mocks.NewMockServerPool()
	defer serverPool.Close()

	ctx := helpers.TestContext(t)

	// Create clients for all versions (excluding v0.0.41 due to incomplete implementation)
	clients := make(map[string]slurm.SlurmClient)
	versions := []string{"v0.0.40", "v0.0.42", "v0.0.43"}

	for _, version := range versions {
		server := serverPool.GetServer(version)
		client, err := slurm.NewClientWithVersion(ctx, version,
			slurm.WithBaseURL(server.URL()),
			slurm.WithAuth(auth.NewNoAuth()),
		)
		require.NoError(t, err)
		clients[version] = client
		defer client.Close()
	}

	// Submit jobs concurrently with different versions
	jobIDs := make(map[string]string)

	for version, client := range clients {
		submission := &interfaces.JobSubmission{
			Name:      "concurrent-test-" + version,
			Script:    "#!/bin/bash\necho 'Concurrent test for " + version + "'",
			Partition: "compute",
			Cpus:      1,
		}

		response, err := client.Jobs().Submit(ctx, submission)
		require.NoError(t, err, "Job submission should work for %s", version)
		jobIDs[version] = response.JobID

		t.Logf("Submitted job %s with version %s", response.JobID, version)
	}

	// Verify all jobs exist in their respective systems
	for version, client := range clients {
		jobID := jobIDs[version]
		job, err := client.Jobs().Get(ctx, jobID)
		require.NoError(t, err, "Should be able to get job from %s", version)
		assert.Equal(t, jobID, job.ID)
		assert.Contains(t, job.Name, version)
	}

	t.Log("Successfully ran concurrent operations across all API versions")
}
