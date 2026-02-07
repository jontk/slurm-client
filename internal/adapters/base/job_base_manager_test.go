// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package base

import (
	"testing"

	types "github.com/jontk/slurm-client/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobBaseManager_New(t *testing.T) {
	manager := NewJobBaseManager("v0.0.43")
	assert.NotNil(t, manager)
	assert.Equal(t, "v0.0.43", manager.GetVersion())
	assert.Equal(t, "Job", manager.GetResourceType())
}
func TestJobBaseManager_ValidateJobCreate(t *testing.T) {
	manager := NewJobBaseManager("v0.0.43")
	tests := []struct {
		name    string
		job     *types.JobCreate
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil job",
			job:     nil,
			wantErr: true,
			errMsg:  "data is required",
		},
		{
			name: "empty script",
			job: &types.JobCreate{
				Script: stringPtr(""),
			},
			wantErr: true,
			errMsg:  "Script is required",
		},
		{
			name: "nil script",
			job: &types.JobCreate{
				Script: nil,
			},
			wantErr: true,
			errMsg:  "Script is required",
		},
		{
			name: "valid job with script",
			job: &types.JobCreate{
				Script:      stringPtr("#!/bin/bash\necho 'test'"),
				MinimumCPUs: int32Ptr(1),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateJobCreate(tt.job)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
func TestJobBaseManager_ValidateJobUpdate(t *testing.T) {
	manager := NewJobBaseManager("v0.0.43")
	tests := []struct {
		name    string
		update  *types.JobUpdate
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil update",
			update:  nil,
			wantErr: true,
			errMsg:  "data is required",
		},
		{
			name: "valid update",
			update: &types.JobUpdate{
				TimeLimit: uint32Ptr(120),
				Priority:  uint32Ptr(1000),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateJobUpdate(tt.update)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
func TestJobBaseManager_ApplyJobDefaults(t *testing.T) {
	manager := NewJobBaseManager("v0.0.43")
	job := &types.JobCreate{
		Script: stringPtr("echo 'test'"),
	}
	result := manager.ApplyJobDefaults(job)
	assert.NotNil(t, result)
	// Test that defaults are applied
	require.NotNil(t, result.MinimumCPUs)
	require.NotNil(t, result.MinimumNodes)
	require.NotNil(t, result.Tasks)
	require.NotNil(t, result.CurrentWorkingDirectory)
	assert.Equal(t, int32(1), *result.MinimumCPUs)
	assert.Equal(t, int32(1), *result.MinimumNodes)
	assert.Equal(t, int32(1), *result.Tasks)
	assert.Equal(t, "/tmp", *result.CurrentWorkingDirectory)
}
func TestJobBaseManager_FilterJobList(t *testing.T) {
	manager := NewJobBaseManager("v0.0.43")
	jobs := []types.Job{
		{JobID: int32Ptr(1), Name: stringPtr("job1"), Account: stringPtr("account1"), JobState: []types.JobState{types.JobStatePending}},
		{JobID: int32Ptr(2), Name: stringPtr("job2"), Account: stringPtr("account2"), JobState: []types.JobState{types.JobStateRunning}},
		{JobID: int32Ptr(3), Name: stringPtr("job3"), Account: stringPtr("account1"), JobState: []types.JobState{types.JobStateCompleted}},
	}
	tests := []struct {
		name     string
		opts     *types.JobListOptions
		expected int
	}{
		{
			name:     "no filters",
			opts:     nil,
			expected: 3,
		},
		{
			name: "filter by account",
			opts: &types.JobListOptions{
				Accounts: []string{"account1"},
			},
			expected: 2,
		},
		{
			name: "filter by state",
			opts: &types.JobListOptions{
				States: []types.JobState{types.JobStatePending},
			},
			expected: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.FilterJobList(jobs, tt.opts)
			assert.Len(t, result, tt.expected)
		})
	}
}
