// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

import (
	"context"
	"testing"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/internal/testutil"
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

func TestJobManager_Submit_NotImplemented(t *testing.T) {
	// Test that Submit returns not implemented error
	jobManager := &JobManager{
		client: &WrapperClient{},
	}

	jobSub := &interfaces.JobSubmission{
		Name:    "test-job",
		Command: "echo hello",
		CPUs:    2,
	}

	_, err := jobManager.Submit(context.Background(), jobSub)

	// v0.0.41 Submit is not implemented
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
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

func TestJobManager_Update_NotImplemented(t *testing.T) {
	// Test that Update returns not implemented error
	jobManager := &JobManager{
		client: &WrapperClient{},
	}

	err := jobManager.Update(context.Background(), "12345", &interfaces.JobUpdate{Priority: testutil.IntPtr(100)})

	// v0.0.41 Update is not implemented
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
	// The impl should now be created
	assert.NotNil(t, jobManager.impl)
}

func TestJobManager_Steps_EmptyList(t *testing.T) {
	// Test that Steps returns empty list (v0.0.41 doesn't support steps)
	jobManager := &JobManager{
		client: &WrapperClient{},
	}

	result, err := jobManager.Steps(context.Background(), "12345")

	// Should not error but return empty list
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.Steps)
	assert.Equal(t, 0, result.Total)
}

func TestJobManager_Watch_Structure(t *testing.T) {
	// Test that Watch method properly delegates to implementation
	jobManager := &JobManager{
		client: &WrapperClient{},
	}

	_, err := jobManager.Watch(context.Background(), &interfaces.WatchJobsOptions{})

	// We expect an error since there's no real API client
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")
	// The impl should now be created
	assert.NotNil(t, jobManager.impl)
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

