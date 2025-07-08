package v0_0_42

import (
	"context"
	"testing"

	"github.com/jontk/slurm-client/internal/interfaces"
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
	
	apiJob := V0042JobInfo{
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