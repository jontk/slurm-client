// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/tests/helpers"
	"github.com/jontk/slurm-client/tests/mocks"
)

// TestJobLifecycle tests the complete job lifecycle: submit → monitor → cancel
func TestJobLifecycle(t *testing.T) {
	testCases := []struct {
		name       string
		apiVersion string
	}{
		{"v0.0.40", "v0.0.40"},
		{"v0.0.41", "v0.0.41"},
		{"v0.0.42", "v0.0.42"},
		{"v0.0.43", "v0.0.43"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testJobLifecycleForVersion(t, tc.apiVersion)
		})
	}
}

func testJobLifecycleForVersion(t *testing.T, apiVersion string) {
	// Setup mock server for the specific API version
	mockServer := mocks.NewMockSlurmServerForVersion(apiVersion)
	defer mockServer.Close()

	// Create client
	ctx := helpers.TestContext(t)
	client, err := slurm.NewClientWithVersion(ctx, apiVersion,
		slurm.WithBaseURL(mockServer.URL()),
		slurm.WithAuth(auth.NewNoAuth()),
	)
	require.NoError(t, err)
	defer client.Close()

	// Verify client version
	assert.Equal(t, apiVersion, client.Version())

	// Phase 1: Submit Job
	t.Run("SubmitJob", func(t *testing.T) {
		jobSubmission := &interfaces.JobSubmission{
			Name:       "integration-test-job",
			Script:     "#!/bin/bash\necho 'Integration test job'\nsleep 10",
			Partition:  "compute",
			CPUs:       2,
			Memory:     2 * 1024 * 1024 * 1024, // 2GB
			TimeLimit:  30,                      // 30 minutes
			WorkingDir: "/tmp",
			Environment: map[string]string{
				"TEST_ENV": "integration",
				"PATH":     "/usr/local/bin:/usr/bin:/bin",
			},
		}

		response, err := client.Jobs().Submit(ctx, jobSubmission)
		require.NoError(t, err)
		require.NotNil(t, response)
		assert.NotEmpty(t, response.JobID)

		// Store job ID for subsequent tests
		jobID := response.JobID
		t.Logf("Submitted job with ID: %s", jobID)

		// Phase 2: Monitor Job (Get job details)
		t.Run("MonitorJob", func(t *testing.T) {
			// Get job immediately after submission
			job, err := client.Jobs().Get(ctx, jobID)
			require.NoError(t, err)
			require.NotNil(t, job)

			// Verify job details
			assert.Equal(t, jobID, job.ID)
			assert.Equal(t, "integration-test-job", job.Name)
			assert.Equal(t, "compute", job.Partition)
			assert.Equal(t, 2, job.CPUs)
			assert.Equal(t, int64(2*1024*1024*1024), job.Memory)
			assert.Equal(t, 30, job.TimeLimit)
			assert.Contains(t, []string{"PENDING", "RUNNING"}, job.State) // Job could be in either state

			t.Logf("Job state: %s", job.State)

			// Test job listing with filtering
			jobs, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{
				UserID:    job.UserID,
				States:    []string{job.State},
				Partition: job.Partition,
				Limit:     10,
			})
			require.NoError(t, err)
			require.NotNil(t, jobs)
			assert.Greater(t, len(jobs.Jobs), 0)

			// Find our job in the list
			found := false
			for _, listedJob := range jobs.Jobs {
				if listedJob.ID == jobID {
					found = true
					assert.Equal(t, job.Name, listedJob.Name)
					assert.Equal(t, job.State, listedJob.State)
					break
				}
			}
			assert.True(t, found, "Submitted job should appear in job list")

			// Test job steps (if supported in this version)
			if mockServer.GetConfig().SupportedOperations["jobs.steps"] {
				steps, err := client.Jobs().Steps(ctx, jobID)
				require.NoError(t, err)
				require.NotNil(t, steps)
				t.Logf("Job has %d steps", len(steps.Steps))
			}

			// Phase 3: Update Job (if supported)
			if mockServer.GetConfig().SupportedOperations["jobs.update"] {
				t.Run("UpdateJob", func(t *testing.T) {
					update := &interfaces.JobUpdate{
						Name:      stringPtr("integration-test-job-updated"),
						TimeLimit: intPtr(45), // Extend time limit
					}

					err := client.Jobs().Update(ctx, jobID, update)
					require.NoError(t, err)

					// Verify update
					updatedJob, err := client.Jobs().Get(ctx, jobID)
					require.NoError(t, err)
					assert.Equal(t, "integration-test-job-updated", updatedJob.Name)
					assert.Equal(t, 45, updatedJob.TimeLimit)
				})
			}

			// Phase 4: Cancel Job
			t.Run("CancelJob", func(t *testing.T) {
				err := client.Jobs().Cancel(ctx, jobID)
				require.NoError(t, err)

				// Verify job is cancelled
				cancelledJob, err := client.Jobs().Get(ctx, jobID)
				require.NoError(t, err)
				assert.Equal(t, "CANCELLED", cancelledJob.State)
				assert.NotNil(t, cancelledJob.EndTime)

				t.Logf("Job cancelled successfully")

				// Try to cancel again - should fail
				err = client.Jobs().Cancel(ctx, jobID)
				assert.Error(t, err, "Cancelling already cancelled job should fail")
			})
		})
	})
}

// TestJobSubmissionValidation tests job submission validation
func TestJobSubmissionValidation(t *testing.T) {
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	ctx := helpers.TestContext(t)
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.42",
		slurm.WithBaseURL(mockServer.URL()),
		slurm.WithAuth(auth.NewNoAuth()),
	)
	require.NoError(t, err)
	defer client.Close()

	testCases := []struct {
		name        string
		submission  *interfaces.JobSubmission
		expectError bool
		errorMsg    string
	}{
		{
			name: "ValidSubmission",
			submission: &interfaces.JobSubmission{
				Name:      "valid-job",
				Script:    "#!/bin/bash\necho 'test'",
				Partition: "compute",
				CPUs:      1,
			},
			expectError: false,
		},
		{
			name: "MissingName",
			submission: &interfaces.JobSubmission{
				Script:    "#!/bin/bash\necho 'test'",
				Partition: "compute",
				CPUs:      1,
			},
			expectError: true,
			errorMsg:    "Job name is required",
		},
		{
			name: "MissingScript",
			submission: &interfaces.JobSubmission{
				Name:      "test-job",
				Partition: "compute",
				CPUs:      1,
			},
			expectError: true,
			errorMsg:    "Job script is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := client.Jobs().Submit(ctx, tc.submission)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, response)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.JobID)
			}
		})
	}
}

// TestJobNotFound tests error handling for non-existent jobs
func TestJobNotFound(t *testing.T) {
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	ctx := helpers.TestContext(t)
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.42",
		slurm.WithBaseURL(mockServer.URL()),
		slurm.WithAuth(auth.NewNoAuth()),
	)
	require.NoError(t, err)
	defer client.Close()

	nonExistentJobID := "99999"

	// Test Get
	job, err := client.Jobs().Get(ctx, nonExistentJobID)
	assert.Error(t, err)
	assert.Nil(t, job)
	assert.Contains(t, err.Error(), "not found")

	// Test Cancel
	err = client.Jobs().Cancel(ctx, nonExistentJobID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// Test Update (if supported)
	if mockServer.GetConfig().SupportedOperations["jobs.update"] {
		update := &interfaces.JobUpdate{
			Name: stringPtr("test"),
		}
		err = client.Jobs().Update(ctx, nonExistentJobID, update)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	}

	// Test Steps
	if mockServer.GetConfig().SupportedOperations["jobs.steps"] {
		steps, err := client.Jobs().Steps(ctx, nonExistentJobID)
		assert.Error(t, err)
		assert.Nil(t, steps)
		assert.Contains(t, err.Error(), "not found")
	}
}

// TestJobListFiltering tests job list filtering capabilities
func TestJobListFiltering(t *testing.T) {
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.42")
	defer mockServer.Close()

	ctx := helpers.TestContext(t)
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.42",
		slurm.WithBaseURL(mockServer.URL()),
		slurm.WithAuth(auth.NewNoAuth()),
	)
	require.NoError(t, err)
	defer client.Close()

	// Test different filtering options
	testCases := []struct {
		name    string
		options *interfaces.ListJobsOptions
	}{
		{
			name:    "NoFilter",
			options: nil,
		},
		{
			name: "FilterByState",
			options: &interfaces.ListJobsOptions{
				States: []string{"RUNNING"},
			},
		},
		{
			name: "FilterByPartition",
			options: &interfaces.ListJobsOptions{
				Partition: "compute",
			},
		},
		{
			name: "FilterByUser",
			options: &interfaces.ListJobsOptions{
				UserID: "testuser",
			},
		},
		{
			name: "MultipleFilters",
			options: &interfaces.ListJobsOptions{
				UserID:    "testuser",
				States:    []string{"RUNNING", "PENDING"},
				Partition: "compute",
			},
		},
		{
			name: "WithPagination",
			options: &interfaces.ListJobsOptions{
				Limit:  1,
				Offset: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jobs, err := client.Jobs().List(ctx, tc.options)
			require.NoError(t, err)
			require.NotNil(t, jobs)

			// Verify pagination info
			assert.GreaterOrEqual(t, jobs.Total, 0)

			// Verify filtering works (at least doesn't error)
			t.Logf("Found %d jobs with filter %+v", len(jobs.Jobs), tc.options)
		})
	}
}

// Helper functions

func stringPtr(s string) *string {
	return &s
}

