// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package mocks

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	v0_0_40 "github.com/jontk/slurm-client/internal/api/v0_0_40"
	v0_0_42 "github.com/jontk/slurm-client/internal/api/v0_0_42"
	v0_0_43 "github.com/jontk/slurm-client/internal/api/v0_0_43"
	v0_0_44 "github.com/jontk/slurm-client/internal/api/v0_0_44"
)

// TestVersionedHandlers_V0040 tests v0.0.40 versioned handlers
func TestVersionedHandlers_V0040(t *testing.T) {
	config := DefaultServerConfig()
	config.APIVersion = "v0.0.40"
	config.SlurmVersion = "24.05"
	server := NewMockSlurmServer(config)
	defer server.Close()

	// Test jobs list endpoint
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL()+"/slurm/v0.0.40/jobs", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to get jobs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	jobs, ok := result["jobs"].([]interface{})
	if !ok || len(jobs) == 0 {
		t.Fatal("Expected jobs array in response")
	}

	// Verify the jobs are v0.0.40 type by checking structure
	if _, ok := server.storage.Jobs["1001"].(*v0_0_40.V0040JobInfo); !ok {
		t.Error("Expected v0.0.40 job type in storage")
	}

	t.Logf("✓ v0.0.40 handler works correctly (%d jobs)", len(jobs))
}

// TestVersionedHandlers_V0042 tests v0.0.42 versioned handlers
func TestVersionedHandlers_V0042(t *testing.T) {
	config := DefaultServerConfig()
	config.APIVersion = "v0.0.42"
	config.SlurmVersion = "24.11"
	server := NewMockSlurmServer(config)
	defer server.Close()

	// Test jobs list endpoint
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL()+"/slurm/v0.0.42/jobs", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to get jobs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	jobs, ok := result["jobs"].([]interface{})
	if !ok || len(jobs) == 0 {
		t.Fatal("Expected jobs array in response")
	}

	// Verify the jobs are v0.0.42 type
	if _, ok := server.storage.Jobs["1001"].(*v0_0_42.V0042JobInfo); !ok {
		t.Error("Expected v0.0.42 job type in storage")
	}

	t.Logf("✓ v0.0.42 handler works correctly (%d jobs)", len(jobs))
}

// TestVersionedHandlers_V0043 tests v0.0.43 versioned handlers
func TestVersionedHandlers_V0043(t *testing.T) {
	config := DefaultServerConfig()
	config.APIVersion = "v0.0.43"
	config.SlurmVersion = "25.05"
	server := NewMockSlurmServer(config)
	defer server.Close()

	// Test jobs list endpoint
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL()+"/slurm/v0.0.43/jobs", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to get jobs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	jobs, ok := result["jobs"].([]interface{})
	if !ok || len(jobs) == 0 {
		t.Fatal("Expected jobs array in response")
	}

	// Verify the jobs are v0.0.43 type
	if _, ok := server.storage.Jobs["1001"].(*v0_0_43.V0043JobInfo); !ok {
		t.Error("Expected v0.0.43 job type in storage")
	}

	t.Logf("✓ v0.0.43 handler works correctly (%d jobs)", len(jobs))
}

// TestVersionedHandlers_V0044 tests v0.0.44 versioned handlers
func TestVersionedHandlers_V0044(t *testing.T) {
	config := DefaultServerConfig()
	config.APIVersion = "v0.0.44"
	config.SlurmVersion = "25.11"
	server := NewMockSlurmServer(config)
	defer server.Close()

	// Test jobs list endpoint
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL()+"/slurm/v0.0.44/jobs", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to get jobs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	jobs, ok := result["jobs"].([]interface{})
	if !ok || len(jobs) == 0 {
		t.Fatal("Expected jobs array in response")
	}

	// Verify the jobs are v0.0.44 type
	if _, ok := server.storage.Jobs["1001"].(*v0_0_44.V0044JobInfo); !ok {
		t.Error("Expected v0.0.44 job type in storage")
	}

	t.Logf("✓ v0.0.44 handler works correctly (%d jobs)", len(jobs))
}

// TestVersionedHandlers_JobGet tests getting individual jobs across versions
func TestVersionedHandlers_JobGet(t *testing.T) {
	versions := []struct {
		version   string
		checkType interface{}
	}{
		{"v0.0.40", (*v0_0_40.V0040JobInfo)(nil)},
		{"v0.0.42", (*v0_0_42.V0042JobInfo)(nil)},
		{"v0.0.43", (*v0_0_43.V0043JobInfo)(nil)},
		{"v0.0.44", (*v0_0_44.V0044JobInfo)(nil)},
	}

	for _, v := range versions {
		t.Run(v.version, func(t *testing.T) {
			config := DefaultServerConfig()
			config.APIVersion = v.version
			config.SlurmVersion = "25.05"
			server := NewMockSlurmServer(config)
			defer server.Close()

			// Get specific job
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL()+"/slurm/"+v.version+"/job/1001", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to get job: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("Expected status 200, got %d", resp.StatusCode)
			}

			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			jobs, ok := result["jobs"].([]interface{})
			if !ok || len(jobs) != 1 {
				t.Fatal("Expected single job in jobs array")
			}

			t.Logf("✓ %s job get works correctly", v.version)
		})
	}
}
