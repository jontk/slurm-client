// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"context"
	"testing"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/internal/testutil"
	"github.com/jontk/slurm-client/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestJobManager_List_Structure(t *testing.T) {
	// Test that the JobManager properly creates implementation
	jobManager := &JobManager{
		client: &WrapperClient{},
	}

	// Test that impl is created lazily
	assert.Nil(t, jobManager.impl)

	// After attempting to call List (even with nil client), impl should be created
	_, err := jobManager.List(context.Background(), nil)

	// We expect an error since there's no real API client
	assert.Error(t, err)
	// The impl should now be created
	assert.NotNil(t, jobManager.impl)
}

func TestConvertAPIJobToInterface(t *testing.T) {
	// Test the conversion function with minimal data
	jobId := int32(12345)
	name := "test-job"
	userId := int32(1000)

	apiJob := V0040JobInfo{
		JobId:  &jobId,
		Name:   &name,
		UserId: &userId,
	}

	interfaceJob, err := convertAPIJobToInterface(apiJob)

	assert.NoError(t, err)
	assert.NotNil(t, interfaceJob)
	assert.Equal(t, "12345", interfaceJob.ID)
	assert.Equal(t, "test-job", interfaceJob.Name)
	assert.Equal(t, "1000", interfaceJob.UserID)
	assert.NotNil(t, interfaceJob.Environment)
	assert.NotNil(t, interfaceJob.Metadata)
}

func TestFilterJobs(t *testing.T) {
	jobs := []interfaces.Job{
		{ID: "1", UserID: "1000", State: "RUNNING", Partition: "gpu"},
		{ID: "2", UserID: "1001", State: "PENDING", Partition: "cpu"},
		{ID: "3", UserID: "1000", State: "COMPLETED", Partition: "gpu"},
	}

	// Test filter by user ID
	opts := &interfaces.ListJobsOptions{UserID: "1000"}
	filtered := filterJobs(jobs, opts)
	assert.Len(t, filtered, 2)
	assert.Equal(t, "1", filtered[0].ID)
	assert.Equal(t, "3", filtered[1].ID)

	// Test filter by state
	opts = &interfaces.ListJobsOptions{States: []string{"RUNNING"}}
	filtered = filterJobs(jobs, opts)
	assert.Len(t, filtered, 1)
	assert.Equal(t, "1", filtered[0].ID)

	// Test filter by partition
	opts = &interfaces.ListJobsOptions{Partition: "gpu"}
	filtered = filterJobs(jobs, opts)
	assert.Len(t, filtered, 2)

	// Test limit and offset
	opts = &interfaces.ListJobsOptions{Limit: 1, Offset: 1}
	filtered = filterJobs(jobs, opts)
	assert.Len(t, filtered, 1)
	assert.Equal(t, "2", filtered[0].ID)
}

func TestJobManager_Get_Structure(t *testing.T) {
	// Test that Get method properly delegates to implementation
	jobManager := &JobManager{
		client: &WrapperClient{},
	}

	_, err := jobManager.Get(context.Background(), "12345")

	// We expect an error since there's no real API client
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")
	// The impl should now be created
	assert.NotNil(t, jobManager.impl)
}

func TestJobManager_Submit_Structure(t *testing.T) {
	// Test that Submit method properly delegates to implementation
	jobManager := &JobManager{
		client: &WrapperClient{},
	}

	jobSub := &interfaces.JobSubmission{
		Name:    "test-job",
		Command: "echo hello",
		CPUs:    2,
	}

	_, err := jobManager.Submit(context.Background(), jobSub)

	// We expect an error since there's no real API client
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")
	// The impl should now be created
	assert.NotNil(t, jobManager.impl)
}

func TestJobManager_Cancel_Structure(t *testing.T) {
	// Test that Cancel method properly delegates to implementation
	jobManager := &JobManager{
		client: &WrapperClient{},
	}

	err := jobManager.Cancel(context.Background(), "12345")

	// We expect an error since there's no real API client
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")
	// The impl should now be created
	assert.NotNil(t, jobManager.impl)
}

func TestConvertJobSubmissionToAPI(t *testing.T) {
	// Test the conversion from interface to API types
	jobSub := &interfaces.JobSubmission{
		Name:        "test-job",
		Script:      "#!/bin/bash\necho hello",
		Partition:   "gpu",
		CPUs:        4,
		Memory:      8 * 1024 * 1024 * 1024, // 8GB in bytes
		TimeLimit:   60,                     // 60 minutes
		Nodes:       2,
		Priority:    100,
		WorkingDir:  "/tmp",
		Environment: map[string]string{"KEY": "value"},
		Args:        []string{"arg1", "arg2"},
	}

	apiJob, err := convertJobSubmissionToAPI(jobSub)

	assert.NoError(t, err)
	assert.NotNil(t, apiJob)
	assert.Equal(t, "test-job", *apiJob.Name)
	assert.Equal(t, "#!/bin/bash\necho hello", *apiJob.Script)
	assert.Equal(t, "gpu", *apiJob.Partition)
	assert.Equal(t, int32(4), *apiJob.MinimumCpus)
	assert.Equal(t, int64(8*1024), *apiJob.MemoryPerNode.Number) // 8GB in MB
	assert.Equal(t, true, *apiJob.MemoryPerNode.Set)
	assert.Equal(t, int64(60), *apiJob.TimeLimit.Number)
	assert.Equal(t, true, *apiJob.TimeLimit.Set)
	assert.Equal(t, int32(2), *apiJob.MinimumNodes)
	assert.Equal(t, int64(100), *apiJob.Priority.Number)
	assert.Equal(t, true, *apiJob.Priority.Set)
	assert.Equal(t, "/tmp", *apiJob.CurrentWorkingDirectory)
	assert.Contains(t, *apiJob.Environment, "KEY=value")
	assert.Equal(t, []string{"arg1", "arg2"}, *apiJob.Argv)
}

func TestConvertJobSubmissionToAPI_NilHandling(t *testing.T) {
	// Test handling of nil job submission
	_, err := convertJobSubmissionToAPI(nil)
	assert.Error(t, err)

	// Check that it returns a structured error
	var slurmErr *errors.SlurmError
	assert.ErrorAs(t, err, &slurmErr)
	assert.Equal(t, errors.ErrorCodeInvalidRequest, slurmErr.Code)
	assert.Contains(t, slurmErr.Message, "Job submission cannot be nil")

	// Test handling of empty job submission
	emptyJob := &interfaces.JobSubmission{}
	apiJob, err := convertJobSubmissionToAPI(emptyJob)
	assert.NoError(t, err)
	assert.NotNil(t, apiJob)
	// Should have default values (nil pointers)
	assert.Nil(t, apiJob.Name)
	assert.Nil(t, apiJob.Script)
	assert.Nil(t, apiJob.MinimumCpus)
}

// TestJobManager_ErrorHandling_Get tests structured error handling for Get method
func TestJobManager_ErrorHandling_Get(t *testing.T) {
	jobManager := &JobManager{
		client: &WrapperClient{}, // No API client initialized
	}

	_, err := jobManager.Get(context.Background(), "12345")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")

	// Check that it returns a structured error
	var slurmErr *errors.SlurmError
	assert.ErrorAs(t, err, &slurmErr)
	assert.Equal(t, errors.ErrorCodeClientNotInitialized, slurmErr.Code)
}

// TestJobManager_ErrorHandling_Submit tests structured error handling for Submit method
func TestJobManager_ErrorHandling_Submit(t *testing.T) {
	jobManager := &JobManager{
		client: &WrapperClient{}, // No API client initialized
	}

	_, err := jobManager.Submit(context.Background(), &interfaces.JobSubmission{Name: "test"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")

	// Check that it returns a structured error
	var slurmErr *errors.SlurmError
	assert.ErrorAs(t, err, &slurmErr)
	assert.Equal(t, errors.ErrorCodeClientNotInitialized, slurmErr.Code)
}

// TestJobManager_ErrorHandling_Cancel tests structured error handling for Cancel method
func TestJobManager_ErrorHandling_Cancel(t *testing.T) {
	jobManager := &JobManager{
		client: &WrapperClient{}, // No API client initialized
	}

	err := jobManager.Cancel(context.Background(), "12345")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")

	// Check that it returns a structured error
	var slurmErr *errors.SlurmError
	assert.ErrorAs(t, err, &slurmErr)
	assert.Equal(t, errors.ErrorCodeClientNotInitialized, slurmErr.Code)
}

// TestJobManager_StructuredErrorTypes tests that methods return proper structured error types
func TestJobManager_StructuredErrorTypes(t *testing.T) {
	jobManager := &JobManager{
		client: &WrapperClient{}, // No API client initialized
	}

	// Test Get method returns SlurmError
	_, err := jobManager.Get(context.Background(), "12345")
	assert.Error(t, err)
	var slurmErr *errors.SlurmError
	assert.ErrorAs(t, err, &slurmErr)
	assert.Equal(t, errors.ErrorCodeClientNotInitialized, slurmErr.Code)

	// Test Submit method returns SlurmError
	_, err = jobManager.Submit(context.Background(), &interfaces.JobSubmission{Name: "test"})
	assert.Error(t, err)
	assert.ErrorAs(t, err, &slurmErr)
	assert.Equal(t, errors.ErrorCodeClientNotInitialized, slurmErr.Code)

	// Test Cancel method returns SlurmError
	err = jobManager.Cancel(context.Background(), "12345")
	assert.Error(t, err)
	assert.ErrorAs(t, err, &slurmErr)
	assert.Equal(t, errors.ErrorCodeClientNotInitialized, slurmErr.Code)

	// Test Update method returns SlurmError
	err = jobManager.Update(context.Background(), "12345", &interfaces.JobUpdate{Priority: testutil.IntPtr(100)})
	assert.Error(t, err)
	assert.ErrorAs(t, err, &slurmErr)
	assert.Equal(t, errors.ErrorCodeClientNotInitialized, slurmErr.Code)

	// Test Steps method returns SlurmError
	_, err = jobManager.Steps(context.Background(), "12345")
	assert.Error(t, err)
	assert.ErrorAs(t, err, &slurmErr)
	assert.Equal(t, errors.ErrorCodeClientNotInitialized, slurmErr.Code)

	// Test Watch method returns SlurmError
	_, err = jobManager.Watch(context.Background(), &interfaces.WatchJobsOptions{})
	assert.Error(t, err)
	assert.ErrorAs(t, err, &slurmErr)
	assert.Equal(t, errors.ErrorCodeClientNotInitialized, slurmErr.Code)
}

// TestJobManager_ErrorHandling_Update tests structured error handling for Update method
func TestJobManager_ErrorHandling_Update(t *testing.T) {
	jobManager := &JobManager{
		client: &WrapperClient{}, // No API client initialized
	}

	err := jobManager.Update(context.Background(), "12345", &interfaces.JobUpdate{Priority: testutil.IntPtr(100)})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")

	// Check that it returns a structured error
	var slurmErr *errors.SlurmError
	assert.ErrorAs(t, err, &slurmErr)
	assert.Equal(t, errors.ErrorCodeClientNotInitialized, slurmErr.Code)
}

// TestJobManager_Update_ValidateInputs tests input validation for Update method
func TestJobManager_Update_ValidateInputs(t *testing.T) {
	// Create a wrapper client (no API client - will trigger client not initialized error first)
	jobManager := &JobManager{
		client: &WrapperClient{}, // No API client initialized
	}

	// Test with nil update - should still fail with client not initialized since we check client first
	err := jobManager.Update(context.Background(), "12345", nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")

	// Check that it returns a structured error
	var slurmErr *errors.SlurmError
	assert.ErrorAs(t, err, &slurmErr)
	assert.Equal(t, errors.ErrorCodeClientNotInitialized, slurmErr.Code)
}

// TestJobManager_convertJobUpdateToAPI tests conversion from interface to API types
func TestJobManager_convertJobUpdateToAPI(t *testing.T) {
	tests := []struct {
		name     string
		update   *interfaces.JobUpdate
		expected *V0040JobDescMsg
		wantErr  bool
	}{
		{
			name: "Update priority only",
			update: &interfaces.JobUpdate{
				Priority: testutil.IntPtr(100),
			},
			expected: &V0040JobDescMsg{
				Priority: &V0040Uint32NoVal{
					Number: int64Ptr(100),
					Set:    boolPtr(true),
				},
			},
			wantErr: false,
		},
		{
			name: "Update time limit only",
			update: &interfaces.JobUpdate{
				TimeLimit: testutil.IntPtr(3600),
			},
			expected: &V0040JobDescMsg{
				TimeLimit: &V0040Uint32NoVal{
					Number: int64Ptr(3600),
					Set:    boolPtr(true),
				},
			},
			wantErr: false,
		},
		{
			name: "Update name only",
			update: &interfaces.JobUpdate{
				Name: stringPtr("updated-job"),
			},
			expected: &V0040JobDescMsg{
				Name: stringPtr("updated-job"),
			},
			wantErr: false,
		},
		{
			name: "Update all fields",
			update: &interfaces.JobUpdate{
				Priority:  testutil.IntPtr(200),
				TimeLimit: testutil.IntPtr(7200),
				Name:      stringPtr("new-job-name"),
			},
			expected: &V0040JobDescMsg{
				Priority: &V0040Uint32NoVal{
					Number: int64Ptr(200),
					Set:    boolPtr(true),
				},
				TimeLimit: &V0040Uint32NoVal{
					Number: int64Ptr(7200),
					Set:    boolPtr(true),
				},
				Name: stringPtr("new-job-name"),
			},
			wantErr: false,
		},
		{
			name:    "Nil update",
			update:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertJobUpdateToAPI(tt.update)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)

			// Check Priority
			if tt.expected.Priority != nil {
				assert.NotNil(t, result.Priority)
				assert.Equal(t, *tt.expected.Priority.Number, *result.Priority.Number)
				assert.Equal(t, *tt.expected.Priority.Set, *result.Priority.Set)
			} else {
				assert.Nil(t, result.Priority)
			}

			// Check TimeLimit
			if tt.expected.TimeLimit != nil {
				assert.NotNil(t, result.TimeLimit)
				assert.Equal(t, *tt.expected.TimeLimit.Number, *result.TimeLimit.Number)
				assert.Equal(t, *tt.expected.TimeLimit.Set, *result.TimeLimit.Set)
			} else {
				assert.Nil(t, result.TimeLimit)
			}

			// Check Name
			if tt.expected.Name != nil {
				assert.NotNil(t, result.Name)
				assert.Equal(t, *tt.expected.Name, *result.Name)
			} else {
				assert.Nil(t, result.Name)
			}
		})
	}
}

// Helper functions for pointer creation

func int64Ptr(i int64) *int64 {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

// TestJobManager_ErrorHandling_Steps tests structured error handling for Steps method
func TestJobManager_ErrorHandling_Steps(t *testing.T) {
	jobManager := &JobManager{
		client: &WrapperClient{}, // No API client initialized
	}

	_, err := jobManager.Steps(context.Background(), "12345")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")

	// Check that it returns a structured error
	var slurmErr *errors.SlurmError
	assert.ErrorAs(t, err, &slurmErr)
	assert.Equal(t, errors.ErrorCodeClientNotInitialized, slurmErr.Code)
}

// TestJobManager_Steps_EmptyList tests that Steps method returns empty list when no steps are available
func TestJobManager_Steps_EmptyList(t *testing.T) {
	// This is a basic test since the current v0.0.40 API implementation
	// returns empty steps list as V0040JobInfo doesn't contain step details
	// In a real implementation, this would need to be tested with mock responses

	// For now, we're just testing that the method structure is correct
	// More comprehensive testing would require mocking the API client
	jobManager := &JobManager{
		client: &WrapperClient{}, // No API client initialized
	}

	_, err := jobManager.Steps(context.Background(), "12345")

	// We expect an error since no API client is initialized
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")
}

// TestJobManager_ErrorHandling_Watch tests structured error handling for Watch method
func TestJobManager_ErrorHandling_Watch(t *testing.T) {
	jobManager := &JobManager{
		client: &WrapperClient{}, // No API client initialized
	}

	_, err := jobManager.Watch(context.Background(), &interfaces.WatchJobsOptions{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")

	// Check that it returns a structured error
	var slurmErr *errors.SlurmError
	assert.ErrorAs(t, err, &slurmErr)
	assert.Equal(t, errors.ErrorCodeClientNotInitialized, slurmErr.Code)
}
