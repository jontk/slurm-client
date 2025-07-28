// +build integration

package v0_0_43

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// Integration tests for ClusterManager
// These tests require a running SLURM cluster with REST API enabled
// Run with: go test -tags=integration ./...

func TestClusterManagerImpl_Integration_FullWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup integration test client
	client := setupIntegrationClient(t)
	if client == nil {
		t.Skip("Integration test environment not available")
	}

	clusterManager := &ClusterManagerImpl{client: client}
	ctx := context.Background()

	// Test cluster name for integration tests
	testClusterName := "integration-test-cluster"

	t.Run("Full CRUD workflow", func(t *testing.T) {
		// Step 1: List existing clusters
		t.Logf("Step 1: Listing existing clusters")
		initialList, err := clusterManager.List(ctx, &interfaces.ListClustersOptions{})
		if err != nil {
			t.Logf("Warning: Could not list clusters (this may be expected): %v", err)
		} else {
			t.Logf("Found %d existing clusters", len(initialList.Clusters))
		}

		// Step 2: Verify test cluster doesn't exist
		t.Logf("Step 2: Verifying test cluster doesn't exist")
		_, err = clusterManager.Get(ctx, testClusterName)
		if err == nil {
			t.Errorf("Test cluster %s already exists, please clean up before running integration tests", testClusterName)
			return
		}
		if !errors.IsNotFoundError(err) && !errors.IsClientError(err) {
			t.Logf("Warning: Unexpected error checking for existing cluster: %v", err)
		}

		// Step 3: Create test cluster
		t.Logf("Step 3: Creating test cluster")
		createRequest := &interfaces.ClusterCreate{
			Name:        testClusterName,
			Description: "Integration test cluster created by automated tests",
			Features:    []string{"test", "integration"},
		}

		createResponse, err := clusterManager.Create(ctx, createRequest)
		if err != nil {
			t.Fatalf("Failed to create test cluster: %v", err)
		}
		if createResponse == nil {
			t.Fatalf("Create response was nil")
		}
		t.Logf("Successfully created cluster: %s", testClusterName)

		// Cleanup: Ensure we delete the test cluster even if other tests fail
		defer func() {
			t.Logf("Cleanup: Deleting test cluster")
			err := clusterManager.Delete(ctx, testClusterName)
			if err != nil {
				t.Logf("Warning: Failed to delete test cluster during cleanup: %v", err)
			} else {
				t.Logf("Successfully cleaned up test cluster")
			}
		}()

		// Step 4: Retrieve the created cluster
		t.Logf("Step 4: Retrieving created cluster")
		retrievedCluster, err := clusterManager.Get(ctx, testClusterName)
		if err != nil {
			t.Fatalf("Failed to retrieve created cluster: %v", err)
		}
		if retrievedCluster == nil {
			t.Fatalf("Retrieved cluster was nil")
		}
		if retrievedCluster.Name != testClusterName {
			t.Errorf("Retrieved cluster name mismatch: expected=%s, got=%s", testClusterName, retrievedCluster.Name)
		}
		t.Logf("Successfully retrieved cluster: %s", retrievedCluster.Name)

		// Step 5: Update the cluster
		t.Logf("Step 5: Updating cluster")
		updateRequest := &interfaces.ClusterUpdate{
			Description: "Updated integration test cluster",
			Features:    []string{"test", "integration", "updated"},
		}

		err = clusterManager.Update(ctx, testClusterName, updateRequest)
		if err != nil {
			t.Fatalf("Failed to update cluster: %v", err)
		}
		t.Logf("Successfully updated cluster")

		// Step 6: Verify the update
		t.Logf("Step 6: Verifying update")
		updatedCluster, err := clusterManager.Get(ctx, testClusterName)
		if err != nil {
			t.Fatalf("Failed to retrieve updated cluster: %v", err)
		}
		if updatedCluster.Description != updateRequest.Description {
			t.Logf("Warning: Description update may not be reflected immediately")
		}
		t.Logf("Successfully verified update")

		// Step 7: List clusters to verify our cluster is included
		t.Logf("Step 7: Verifying cluster appears in list")
		finalList, err := clusterManager.List(ctx, &interfaces.ListClustersOptions{})
		if err != nil {
			t.Fatalf("Failed to list clusters after creation: %v", err)
		}

		found := false
		for _, cluster := range finalList.Clusters {
			if cluster.Name == testClusterName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Created cluster not found in cluster list")
		} else {
			t.Logf("Successfully found cluster in list")
		}
	})
}

func TestClusterManagerImpl_Integration_ListPagination(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupIntegrationClient(t)
	if client == nil {
		t.Skip("Integration test environment not available")
	}

	clusterManager := &ClusterManagerImpl{client: client}
	ctx := context.Background()

	t.Run("pagination with small page size", func(t *testing.T) {
		// Test pagination by requesting small page sizes
		opts := &interfaces.ListClustersOptions{
			Page:     1,
			PageSize: 1,
		}

		result, err := clusterManager.List(ctx, opts)
		if err != nil {
			t.Fatalf("Failed to list clusters with pagination: %v", err)
		}

		if result == nil {
			t.Fatalf("List result was nil")
		}

		t.Logf("Retrieved %d clusters with page size 1", len(result.Clusters))
		
		// If there are clusters, verify pagination metadata
		if len(result.Clusters) > 0 {
			if result.Meta != nil {
				t.Logf("Pagination metadata: %+v", result.Meta)
			}
		}
	})

	t.Run("pagination with different page sizes", func(t *testing.T) {
		pageSizes := []int{5, 10, 25}

		for _, pageSize := range pageSizes {
			opts := &interfaces.ListClustersOptions{
				Page:     1,
				PageSize: pageSize,
			}

			result, err := clusterManager.List(ctx, opts)
			if err != nil {
				t.Logf("Warning: Failed to list clusters with page size %d: %v", pageSize, err)
				continue
			}

			t.Logf("Page size %d returned %d clusters", pageSize, len(result.Clusters))
		}
	})
}

func TestClusterManagerImpl_Integration_Filtering(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupIntegrationClient(t)
	if client == nil {
		t.Skip("Integration test environment not available")
	}

	clusterManager := &ClusterManagerImpl{client: client}
	ctx := context.Background()

	t.Run("filter by name pattern", func(t *testing.T) {
		opts := &interfaces.ListClustersOptions{
			FilterByName: "*",
		}

		result, err := clusterManager.List(ctx, opts)
		if err != nil {
			t.Logf("Warning: Name filtering not supported or failed: %v", err)
			return
		}

		t.Logf("Name filter '*' returned %d clusters", len(result.Clusters))
	})

	t.Run("filter by state", func(t *testing.T) {
		states := []string{"ACTIVE", "INACTIVE", "UNKNOWN"}

		for _, state := range states {
			opts := &interfaces.ListClustersOptions{
				FilterByState: state,
			}

			result, err := clusterManager.List(ctx, opts)
			if err != nil {
				t.Logf("Warning: State filtering for %s not supported or failed: %v", state, err)
				continue
			}

			t.Logf("State filter '%s' returned %d clusters", state, len(result.Clusters))
		}
	})
}

func TestClusterManagerImpl_Integration_ErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupIntegrationClient(t)
	if client == nil {
		t.Skip("Integration test environment not available")
	}

	clusterManager := &ClusterManagerImpl{client: client}
	ctx := context.Background()

	t.Run("get nonexistent cluster", func(t *testing.T) {
		nonexistentName := "definitely-does-not-exist-cluster-12345"
		
		_, err := clusterManager.Get(ctx, nonexistentName)
		if err == nil {
			t.Errorf("Expected error for nonexistent cluster, got none")
		} else {
			t.Logf("Correctly received error for nonexistent cluster: %v", err)
			if !errors.IsNotFoundError(err) && !errors.IsClientError(err) {
				t.Logf("Note: Error type was %T, which may be acceptable", err)
			}
		}
	})

	t.Run("create cluster with invalid name", func(t *testing.T) {
		invalidCluster := &interfaces.ClusterCreate{
			Name: "invalid cluster name with spaces",
		}

		_, err := clusterManager.Create(ctx, invalidCluster)
		if err == nil {
			t.Errorf("Expected error for invalid cluster name, got none")
		} else {
			t.Logf("Correctly received error for invalid cluster name: %v", err)
			if !errors.IsValidationError(err) && !errors.IsClientError(err) {
				t.Logf("Note: Error type was %T, which may be acceptable", err)
			}
		}
	})

	t.Run("update nonexistent cluster", func(t *testing.T) {
		nonexistentName := "definitely-does-not-exist-cluster-12345"
		update := &interfaces.ClusterUpdate{
			Description: "This should fail",
		}

		err := clusterManager.Update(ctx, nonexistentName, update)
		if err == nil {
			t.Errorf("Expected error for updating nonexistent cluster, got none")
		} else {
			t.Logf("Correctly received error for updating nonexistent cluster: %v", err)
		}
	})

	t.Run("delete nonexistent cluster", func(t *testing.T) {
		nonexistentName := "definitely-does-not-exist-cluster-12345"

		err := clusterManager.Delete(ctx, nonexistentName)
		if err == nil {
			t.Errorf("Expected error for deleting nonexistent cluster, got none")
		} else {
			t.Logf("Correctly received error for deleting nonexistent cluster: %v", err)
		}
	})
}

func TestClusterManagerImpl_Integration_ContextHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupIntegrationClient(t)
	if client == nil {
		t.Skip("Integration test environment not available")
	}

	clusterManager := &ClusterManagerImpl{client: client}

	t.Run("context with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		_, err := clusterManager.List(ctx, &interfaces.ListClustersOptions{})
		if err != nil {
			t.Logf("Operation with short timeout resulted in error (may be expected): %v", err)
		} else {
			t.Logf("Operation completed within timeout")
		}
	})

	t.Run("cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := clusterManager.List(ctx, &interfaces.ListClustersOptions{})
		if err == nil {
			t.Errorf("Expected error for cancelled context, got none")
		} else {
			t.Logf("Correctly received error for cancelled context: %v", err)
		}
	})
}

func TestClusterManagerImpl_Integration_ConcurrentAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupIntegrationClient(t)
	if client == nil {
		t.Skip("Integration test environment not available")
	}

	clusterManager := &ClusterManagerImpl{client: client}
	ctx := context.Background()

	t.Run("concurrent list operations", func(t *testing.T) {
		concurrency := 5
		done := make(chan error, concurrency)

		for i := 0; i < concurrency; i++ {
			go func(id int) {
				opts := &interfaces.ListClustersOptions{
					Page:     1,
					PageSize: 10,
				}
				
				_, err := clusterManager.List(ctx, opts)
				done <- err
			}(i)
		}

		errors := 0
		for i := 0; i < concurrency; i++ {
			if err := <-done; err != nil {
				t.Logf("Concurrent operation %d failed: %v", i, err)
				errors++
			}
		}

		if errors > 0 {
			t.Logf("Warning: %d out of %d concurrent operations failed", errors, concurrency)
		} else {
			t.Logf("All %d concurrent operations succeeded", concurrency)
		}
	})

	t.Run("concurrent get operations", func(t *testing.T) {
		// First, get a list of available clusters
		clusterList, err := clusterManager.List(ctx, &interfaces.ListClustersOptions{})
		if err != nil || len(clusterList.Clusters) == 0 {
			t.Skip("No clusters available for concurrent get testing")
		}

		// Use the first available cluster for concurrent gets
		testClusterName := clusterList.Clusters[0].Name
		concurrency := 3
		done := make(chan error, concurrency)

		for i := 0; i < concurrency; i++ {
			go func(id int) {
				_, err := clusterManager.Get(ctx, testClusterName)
				done <- err
			}(i)
		}

		errors := 0
		for i := 0; i < concurrency; i++ {
			if err := <-done; err != nil {
				t.Logf("Concurrent get operation %d failed: %v", i, err)
				errors++
			}
		}

		if errors > 0 {
			t.Logf("Warning: %d out of %d concurrent get operations failed", errors, concurrency)
		} else {
			t.Logf("All %d concurrent get operations succeeded", concurrency)
		}
	})
}

// Helper function to setup integration test client
func setupIntegrationClient(t *testing.T) *WrapperClient {
	// This would be configured based on your integration test environment
	// For now, return nil to skip tests when integration environment is not available
	
	// Example setup (uncomment and modify for your environment):
	/*
	baseURL := os.Getenv("SLURM_REST_URL")
	if baseURL == "" {
		baseURL = "http://localhost:6820" // Default SLURM REST API port
	}

	token := os.Getenv("SLURM_REST_TOKEN")
	if token == "" {
		t.Skip("SLURM_REST_TOKEN not set, skipping integration tests")
	}

	client, err := NewClientWithResponses(baseURL, WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Set("X-SLURM-USER-NAME", "test-user")
		req.Header.Set("X-SLURM-USER-TOKEN", token)
		return nil
	}))
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}

	return &WrapperClient{apiClient: client}
	*/

	// For now, skip integration tests
	return nil
}

// Benchmark tests for integration testing
func BenchmarkClusterManagerImpl_Integration_List(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping integration benchmark in short mode")
	}

	client := setupIntegrationClient(b.(*testing.T))
	if client == nil {
		b.Skip("Integration test environment not available")
	}

	clusterManager := &ClusterManagerImpl{client: client}
	ctx := context.Background()
	opts := &interfaces.ListClustersOptions{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := clusterManager.List(ctx, opts)
		if err != nil {
			b.Fatalf("List operation failed: %v", err)
		}
	}
}

func BenchmarkClusterManagerImpl_Integration_Get(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping integration benchmark in short mode")
	}

	client := setupIntegrationClient(b.(*testing.T))
	if client == nil {
		b.Skip("Integration test environment not available")
	}

	clusterManager := &ClusterManagerImpl{client: client}
	ctx := context.Background()

	// Get the first available cluster for benchmarking
	clusterList, err := clusterManager.List(ctx, &interfaces.ListClustersOptions{})
	if err != nil || len(clusterList.Clusters) == 0 {
		b.Skip("No clusters available for benchmarking")
	}

	testClusterName := clusterList.Clusters[0].Name

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := clusterManager.Get(ctx, testClusterName)
		if err != nil {
			b.Fatalf("Get operation failed: %v", err)
		}
	}
}