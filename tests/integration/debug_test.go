package integration

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func TestDebugJobSubmission(t *testing.T) {
	if os.Getenv("SLURM_REAL_SERVER_TEST") != "true" {
		t.Skip("Skipping real server tests")
	}

	// Fetch JWT token
	token, err := fetchJWTTokenViaSSH()
	require.NoError(t, err, "Failed to fetch JWT token")

	// Create a simple job submission payload
	jobPayload := map[string]interface{}{
		"job": map[string]interface{}{
			"name":   "debug-test-job",
			"script": "#!/bin/bash\necho 'Hello from test job'\n",
			"current_working_directory": "/tmp",
			"minimum_cpus":              1,
			"memory_per_node": map[string]interface{}{
				"set":    true,
				"number": 1024, // MB
			},
			"time_limit": map[string]interface{}{
				"set":    true,
				"number": 5, // minutes
			},
		},
	}

	jsonData, err := json.Marshal(jobPayload)
	require.NoError(t, err)

	// Make the request
	host := getEnvOrDefault("SLURM_HOST", "rocky9")
	port := getEnvOrDefault("SLURM_PORT", "6820")
	url := fmt.Sprintf("https://%s:%s/slurm/v0.0.43/job/submit", host, port)
	req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewReader(jsonData))
	require.NoError(t, err)

	req.Header.Set("X-SLURM-USER-TOKEN", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	t.Logf("Status Code: %d", resp.StatusCode)
	t.Logf("Response Body: %s", string(body))

	// Try to parse the response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Logf("Failed to parse response as JSON: %v", err)
	} else {
		prettyJSON, _ := json.MarshalIndent(result, "", "  ")
		t.Logf("Response JSON:\n%s", string(prettyJSON))
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected successful job submission")
}

func TestDebugListPartitions(t *testing.T) {
	if os.Getenv("SLURM_REAL_SERVER_TEST") != "true" {
		t.Skip("Skipping real server tests")
	}

	// Fetch JWT token
	token, err := fetchJWTTokenViaSSH()
	require.NoError(t, err, "Failed to fetch JWT token")

	// Make the request
	host := getEnvOrDefault("SLURM_HOST", "rocky9")
	port := getEnvOrDefault("SLURM_PORT", "6820")
	url := fmt.Sprintf("https://%s:%s/slurm/v0.0.43/partitions?flags=DETAIL", host, port)
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	require.NoError(t, err)

	req.Header.Set("X-SLURM-USER-TOKEN", token)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	t.Logf("Status Code: %d", resp.StatusCode)
	t.Logf("Response Body: %s", string(body))

	// Try to parse the response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Logf("Failed to parse response as JSON: %v", err)
	} else {
		prettyJSON, _ := json.MarshalIndent(result, "", "  ")
		t.Logf("Response JSON:\n%s", string(prettyJSON))
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected successful partition list")
}