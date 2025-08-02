// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
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
			name: "empty command and script",
			job: &types.JobCreate{
				Command: "",
				Script:  "",
			},
			wantErr: true,
			errMsg:  "command or script is required",
		},
		{
			name: "valid job with command",
			job: &types.JobCreate{
				Command: "echo 'test'",
				CPUs:    1,
			},
			wantErr: false,
		},
		{
			name: "valid job with script",
			job: &types.JobCreate{
				Script: "#!/bin/bash\necho 'test'",
				CPUs:   1,
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
				TimeLimit: int32Ptr(120),
				Priority:  int32Ptr(1000),
			},
			wantErr: false,
		},
		{
			name: "negative time limit",
			update: &types.JobUpdate{
				TimeLimit: int32Ptr(-1),
			},
			wantErr: true,
			errMsg:  "must be non-negative",
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
		Command: "echo 'test'",
	}

	result := manager.ApplyJobDefaults(job)
	assert.NotNil(t, result)
	// Test that defaults are applied (basic check that it returns a result)
	assert.Equal(t, job.Command, result.Command)
}

func TestJobBaseManager_FilterJobList(t *testing.T) {
	manager := NewJobBaseManager("v0.0.43")

	jobs := []types.Job{
		{JobID: 1, Name: "job1", Account: "account1", State: types.JobStatePending},
		{JobID: 2, Name: "job2", Account: "account2", State: types.JobStateRunning},
		{JobID: 3, Name: "job3", Account: "account1", State: types.JobStateCompleted},
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