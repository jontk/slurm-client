//go:build integration
// +build integration

// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/jontk/slurm-client"
	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

// IntegrationTestSuite provides common functionality for all version integration tests
type IntegrationTestSuite struct {
	suite.Suite
	client    slurm.SlurmClient
	serverURL string
	token     string
	version   string
}

// TestConfig holds configuration for integration tests
type TestConfig struct {
	ServerURL           string
	Token               string
	Version             string
	Timeout             time.Duration
	MaxRetries          int
	Debug               bool
	InsecureSkipVerify  bool
	RequireDatabase     bool
	SkipSlowTests       bool
	TestResourceCleanup bool
}

// tokenTransport adds the SLURM JWT token to requests
type tokenTransport struct {
	token string
	base  http.RoundTripper
}

func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	r := req.Clone(req.Context())
	r.Header.Set("X-SLURM-USER-TOKEN", t.token)
	return t.base.RoundTrip(r)
}

// intPtr is a helper function to create a pointer to an int
func intPtr(i int) *int {
	return &i
}

// stringPtr is a helper function to create a pointer to a string
func stringPtr(s string) *string {
	return &s
}

// GetTestConfig loads test configuration from environment variables
func GetTestConfig(version string) *TestConfig {
	config := &TestConfig{
		ServerURL:           getEnvWithDefault("SLURM_SERVER_URL", "http://localhost
		Version:             version,
		Timeout:             30 * time.Second,
		MaxRetries:          3,
		Debug:               getEnvBool("SLURM_TEST_DEBUG", true),
		InsecureSkipVerify:  getEnvBool("SLURM_INSECURE_SKIP_VERIFY", true),
		RequireDatabase:     getEnvBool("SLURM_REQUIRE_DATABASE", false),
		SkipSlowTests:       getEnvBool("SLURM_SKIP_SLOW_TESTS", false),
		TestResourceCleanup: getEnvBool("SLURM_TEST_CLEANUP", true),
	}

	// Get token from environment or fetch via SSH
	config.Token = os.Getenv("SLURM_JWT_TOKEN")

	return config
}

// SetupIntegrationSuite initializes a test suite for a specific API version
func (suite *IntegrationTestSuite) SetupIntegrationSuite(version string) {
	// Check if real server testing is enabled
	if os.Getenv("SLURM_REAL_SERVER_TEST") != "true" {
		suite.T().Skip("Real server tests disabled. Set SLURM_REAL_SERVER_TEST=true to enable")
	}

	testConfig := GetTestConfig(version)
	suite.serverURL = testConfig.ServerURL
	suite.version = testConfig.Version

	// Get or fetch JWT token
	if testConfig.Token == "" {
		token, err := fetchJWTTokenViaSSH()
		require.NoError(suite.T(), err, "Failed to fetch JWT token via SSH")
		suite.token = token
		suite.T().Logf("Fetched JWT token for %s: %s...", version, suite.token[:50])
	} else {
		suite.token = testConfig.Token
	}

	// Create client for this version
	ctx := context.Background()
	client, err := slurm.NewClientWithVersion(ctx, suite.version,
		slurm.WithBaseURL(suite.serverURL),
		slurm.WithAuth(auth.NewTokenAuth(suite.token)),
		slurm.WithConfig(&config.Config{
			Timeout:            testConfig.Timeout,
			MaxRetries:         testConfig.MaxRetries,
			Debug:              testConfig.Debug,
			InsecureSkipVerify: testConfig.InsecureSkipVerify,
		}),
	)
	require.NoError(suite.T(), err, "Failed to create SLURM client for version %s", version)
	suite.client = client
}

// TearDownIntegrationSuite cleans up after tests
func (suite *IntegrationTestSuite) TearDownIntegrationSuite() {
	if suite.client != nil {
		_ = suite.client.Close() // Ignore error during test cleanup
	}
}

// TestPing verifies basic connectivity
func (suite *IntegrationTestSuite) TestPing() {
	ctx := context.Background()
	err := suite.client.Info().Ping(ctx)
	suite.NoError(err, "Ping should succeed for version %s", suite.version)
}

// SkipIfNoDatabaseConnection skips the test if slurmdbd is not available
func (suite *IntegrationTestSuite) SkipIfNoDatabaseConnection(err error) {
	if err != nil {
		errorStr := err.Error()
		if strings.Contains(errorStr, "Unable to connect to database") ||
			strings.Contains(errorStr, "Failed to open slurmdbd connection") ||
			strings.Contains(errorStr, "Bad Gateway") ||
			strings.Contains(errorStr, "HTTP 502") {
			suite.T().Skip("Skipping test: slurmdbd is not connected")
			return
		}
	}
}

// RequireDatabaseConnection fails the test if database operations are not available
func (suite *IntegrationTestSuite) RequireDatabaseConnection(err error) {
	if err != nil {
		errorStr := err.Error()
		if strings.Contains(errorStr, "Unable to connect to database") ||
			strings.Contains(errorStr, "Failed to open slurmdbd connection") ||
			strings.Contains(errorStr, "Bad Gateway") ||
			strings.Contains(errorStr, "HTTP 502") {
			suite.T().Fatal("Test requires slurmdbd connection but it's not available")
			return
		}
	}
}

// CreateTestJob creates a test job for testing purposes
func (suite *IntegrationTestSuite) CreateTestJob(name string) (*interfaces.JobSubmitResponse, error) {
	ctx := context.Background()

	submission := &interfaces.JobSubmission{
		Name:      fmt.Sprintf("%s-test-%s", name, time.Now().Format("20060102-150405")),
		Script:    "#!/bin/bash\necho 'Integration test job'\nhostname\ndate\nsleep 10",
		Partition: "debug", // Using debug partition which should exist on test servers
		Nodes:     1,
		Cpus:      1,
		TimeLimit: 5, // 5 minutes
	}

	suite.T().Logf("Creating test job: %s for version %s", submission.Name, suite.version)
	return suite.client.Jobs().Submit(ctx, submission)
}

// CleanupTestJob cancels and cleans up a test job
func (suite *IntegrationTestSuite) CleanupTestJob(jobID string) {
	if jobID == "" {
		return
	}

	ctx := context.Background()
	err := suite.client.Jobs().Cancel(ctx, jobID)
	if err != nil {
		suite.T().Logf("Warning: Failed to cleanup job %s: %v", jobID, err)
	} else {
		suite.T().Logf("Cleaned up test job: %s", jobID)
	}
}

// WaitForJobState waits for a job to reach a specific state
func (suite *IntegrationTestSuite) WaitForJobState(jobID string, expectedState string, timeout time.Duration) error {
	ctx := context.Background()
	start := time.Now()

	for time.Since(start) < timeout {
		job, err := suite.client.Jobs().Get(ctx, jobID)
		if err != nil {
			return fmt.Errorf("failed to get job %s: %w", jobID, err)
		}

		if job.State == expectedState {
			return nil
		}

		suite.T().Logf("Job %s state: %s (waiting for %s)", jobID, job.State, expectedState)
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("job %s did not reach state %s within %v", jobID, expectedState, timeout)
}

// TestCRUDWorkflow tests a complete CRUD workflow for a resource manager
func (suite *IntegrationTestSuite) TestCRUDWorkflow(resourceName string, testFunc func()) {
	suite.T().Logf("Starting CRUD workflow test for %s on version %s", resourceName, suite.version)
	testFunc()
	suite.T().Logf("Completed CRUD workflow test for %s on version %s", resourceName, suite.version)
}

// TestConcurrentOperations tests concurrent operations on a resource
func (suite *IntegrationTestSuite) TestConcurrentOperations(resourceName string, operation func(id int) error, concurrency int) {
	suite.T().Logf("Starting concurrent operations test for %s (concurrency: %d) on version %s",
		resourceName, concurrency, suite.version)

	errChan := make(chan error, concurrency)

	for i := range concurrency {
		go func(id int) {
			errChan <- operation(id)
		}(i)
	}

	var errors []error
	for range concurrency {
		if err := <-errChan; err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		suite.T().Errorf("Concurrent operations failed with %d errors: %v", len(errors), errors)
	} else {
		suite.T().Logf("All %d concurrent operations succeeded for %s", concurrency, resourceName)
	}
}

// TestErrorHandling tests error handling scenarios
func (suite *IntegrationTestSuite) TestErrorHandling(scenario string, testFunc func() error, expectError bool) {
	suite.T().Logf("Testing error handling scenario: %s on version %s", scenario, suite.version)

	err := testFunc()

	if expectError && err == nil {
		suite.T().Errorf("Expected error for scenario '%s' but got none", scenario)
	} else if !expectError && err != nil {
		suite.T().Errorf("Unexpected error for scenario '%s': %v", scenario, err)
	} else if expectError && err != nil {
		suite.T().Logf("Got expected error for scenario '%s': %v", scenario, err)
	} else {
		suite.T().Logf("Scenario '%s' succeeded as expected", scenario)
	}
}

// TestPerformance measures the performance of an operation
func (suite *IntegrationTestSuite) TestPerformance(operationName string, operation func() error, maxDuration time.Duration) {
	suite.T().Logf("Testing performance for %s on version %s (max duration: %v)",
		operationName, suite.version, maxDuration)

	start := time.Now()
	err := operation()
	duration := time.Since(start)

	suite.NoError(err, "Performance test operation should succeed")

	if duration > maxDuration {
		suite.T().Errorf("Operation %s took %v, which exceeds maximum %v", operationName, duration, maxDuration)
	} else {
		suite.T().Logf("Operation %s completed in %v (within limit)", operationName, duration)
	}
}

// Utility functions

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.EqualFold(value, "true") || value == "1"
}

// fetchJWTTokenViaSSH fetches a JWT token from the server via SSH
func fetchJWTTokenViaSSH() (string, error) {
	// Get SSH configuration from environment
	sshHost := getEnvWithDefault("SLURM_SSH_HOST", "localhost
	sshUser := getEnvWithDefault("SLURM_SSH_USER", "root")

	// Command to get JWT token
	// #nosec G204 -- This is test infrastructure code; SSH host/user are from controlled test environment variables
	sshTarget := sshUser + "@" + sshHost // #nosec G204
	cmd := exec.CommandContext(context.Background(), "ssh",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		sshTarget,
		"unset SLURM_JWT; /opt/slurm/current/bin/scontrol token | grep SLURM_JWT | cut -d= -f2")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("SSH command failed: %w, output: %s", err, string(output))
	}

	token := strings.TrimSpace(string(output))
	if token == "" {
		return "", fmt.Errorf("no token found in output: %s", string(output))
	}

	// Remove any "Warning:" lines from SSH
	lines := strings.Split(token, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "eyJ") { // JWT tokens start with "eyJ"
			return line, nil
		}
	}

	return "", fmt.Errorf("no valid JWT token found in output: %s", string(output))
}

// DatabaseAvailabilityCheck checks if slurmdbd is available
type DatabaseAvailabilityCheck struct {
	available bool
	checked   bool
}

var dbCheck DatabaseAvailabilityCheck

// IsDatabaseAvailable checks if the database (slurmdbd) is available
func (suite *IntegrationTestSuite) IsDatabaseAvailable() bool {
	if dbCheck.checked {
		return dbCheck.available
	}

	// Try a simple database operation to check availability
	ctx := context.Background()
	_, err := suite.client.QoS().List(ctx, &interfaces.ListQoSOptions{Limit: 1})

	dbCheck.available = err == nil || !suite.isDatabaseError(err)
	dbCheck.checked = true

	if dbCheck.available {
		suite.T().Logf("Database is available for version %s", suite.version)
	} else {
		suite.T().Logf("Database is not available for version %s: %v", suite.version, err)
	}

	return dbCheck.available
}

// isDatabaseError checks if an error is due to database unavailability
func (suite *IntegrationTestSuite) isDatabaseError(err error) bool {
	if err == nil {
		return false
	}

	errorStr := err.Error()
	return strings.Contains(errorStr, "Unable to connect to database") ||
		strings.Contains(errorStr, "Failed to open slurmdbd connection") ||
		strings.Contains(errorStr, "Bad Gateway") ||
		strings.Contains(errorStr, "HTTP 502")
}

// ResourceCleanupTracker tracks resources that need cleanup
type ResourceCleanupTracker struct {
	jobIDs []string
}

var cleanupTracker ResourceCleanupTracker

// AddJobForCleanup adds a job ID to the cleanup list
func (suite *IntegrationTestSuite) AddJobForCleanup(jobID string) {
	cleanupTracker.jobIDs = append(cleanupTracker.jobIDs, jobID)
}

// CleanupAllResources cleans up all tracked resources
func (suite *IntegrationTestSuite) CleanupAllResources() {
	for _, jobID := range cleanupTracker.jobIDs {
		suite.CleanupTestJob(jobID)
	}
	cleanupTracker.jobIDs = nil
}
