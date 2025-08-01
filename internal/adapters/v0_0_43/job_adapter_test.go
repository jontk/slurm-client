// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobAdapter_ValidateJobCreate(t *testing.T) {
	adapter := &JobAdapter{
		JobBaseManager: base.NewJobBaseManager("v0.0.43"),
	}

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
			errMsg:  "job data is required",
		},
		{
			name: "empty command and script",
			job: &types.JobCreate{
				Command: "",
				Script:  "",
			},
			wantErr: true,
			errMsg:  "either command or script is required",
		},
		{
			name: "negative CPU count",
			job: &types.JobCreate{
				Command: "echo 'test'",
				CPUs:    -1,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "valid job with command",
			job: &types.JobCreate{
				Command: "echo 'test'",
				CPUs:    1,
				Nodes:   1,
			},
			wantErr: false,
		},
		{
			name: "valid job with script",
			job: &types.JobCreate{
				Script: "#!/bin/bash\necho 'test'",
				CPUs:   2,
				Nodes:  1,
				Name:   "test-job",
			},
			wantErr: false,
		},
		{
			name: "complex job with all fields",
			job: &types.JobCreate{
				Script:           "#!/bin/bash\necho 'complex test'",
				Name:             "complex-job",
				Account:          "test-account",
				Partition:        "compute",
				QoS:              "normal",
				TimeLimit:        120,
				CPUs:             4,
				Nodes:            2,
				Tasks:            4,
				WorkingDirectory: "/tmp/test",
				StandardOutput:   "/tmp/test.out",
				StandardError:    "/tmp/test.err",
				Environment:      map[string]string{"TEST_VAR": "test_value"},
				Features:         []string{"gpu", "high-memory"},
				MailType:         []string{"BEGIN", "END", "FAIL"},
				MailUser:         "user@example.com",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateJobCreate(tt.job)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestJobAdapter_ApplyJobDefaults(t *testing.T) {
	adapter := &JobAdapter{
		JobBaseManager: base.NewJobBaseManager("v0.0.43"),
	}

	tests := []struct {
		name     string
		input    *types.JobCreate
		expected *types.JobCreate
	}{
		{
			name: "apply defaults to minimal job",
			input: &types.JobCreate{
				Command: "echo 'test'",
			},
			expected: &types.JobCreate{
				Command:   "echo 'test'",
				CPUs:      1,     // Default CPU count
				Nodes:     1,     // Default node count
				Tasks:     1,     // Default task count
				TimeLimit: 60,    // Default 1 hour
			},
		},
		{
			name: "preserve existing values",
			input: &types.JobCreate{
				Script:    "#!/bin/bash\necho 'test'",
				Name:      "custom-job",
				CPUs:      4,
				Nodes:     2,
				Tasks:     8,
				TimeLimit: 240,
				Account:   "custom-account",
				Partition: "gpu",
				QoS:       "high",
			},
			expected: &types.JobCreate{
				Script:    "#!/bin/bash\necho 'test'",
				Name:      "custom-job",
				CPUs:      4,
				Nodes:     2,
				Tasks:     8,
				TimeLimit: 240,
				Account:   "custom-account",
				Partition: "gpu",
				QoS:       "high",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.ApplyJobDefaults(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJobAdapter_FilterJobList(t *testing.T) {
	adapter := &JobAdapter{
		JobBaseManager: base.NewJobBaseManager("v0.0.43"),
	}

	jobs := []types.Job{
		{
			JobID:     12345,
			Name:      "job1",
			UserName:  "user1",
			Account:   "account1",
			Partition: "compute",
			QoS:       "normal",
			State:     types.JobStateRunning,
			CPUs:      4,
			Nodes:     1,
		},
		{
			JobID:     12346,
			Name:      "job2",
			UserName:  "user2",
			Account:   "account2",
			Partition: "gpu",
			QoS:       "high",
			State:     types.JobStatePending,
			CPUs:      8,
			Nodes:     2,
		},
		{
			JobID:     12347,
			Name:      "job3",
			UserName:  "user1",
			Account:   "account1",
			Partition: "compute",
			QoS:       "normal",
			State:     types.JobStateCompleted,
			CPUs:      2,
			Nodes:     1,
		},
	}

	tests := []struct {
		name     string
		opts     *types.JobListOptions
		expected []int32 // expected job IDs
	}{
		{
			name:     "no filters",
			opts:     &types.JobListOptions{},
			expected: []int32{12345, 12346, 12347},
		},
		{
			name: "filter by user",
			opts: &types.JobListOptions{
				Users: []string{"user1"},
			},
			expected: []int32{12345, 12347},
		},
		{
			name: "filter by account",
			opts: &types.JobListOptions{
				Accounts: []string{"account2"},
			},
			expected: []int32{12346},
		},
		{
			name: "filter by partition",
			opts: &types.JobListOptions{
				Partitions: []string{"gpu"},
			},
			expected: []int32{12346},
		},
		{
			name: "filter by state",
			opts: &types.JobListOptions{
				States: []string{"RUNNING", "PENDING"},
			},
			expected: []int32{12345, 12346},
		},
		{
			name: "filter by QoS",
			opts: &types.JobListOptions{
				QoSList: []string{"high"},
			},
			expected: []int32{12346},
		},
		{
			name: "combined filters",
			opts: &types.JobListOptions{
				Users:    []string{"user1"},
				States:   []string{"RUNNING"},
				Accounts: []string{"account1"},
			},
			expected: []int32{12345},
		},
		{
			name: "no matches",
			opts: &types.JobListOptions{
				Users: []string{"nonexistent"},
			},
			expected: []int32{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.FilterJobList(jobs, tt.opts)
			resultIDs := make([]int32, len(result))
			for i, job := range result {
				resultIDs[i] = job.JobID
			}
			assert.Equal(t, tt.expected, resultIDs)
		})
	}
}

func TestJobAdapter_ValidateResourceRequests(t *testing.T) {
	adapter := &JobAdapter{
		JobBaseManager: base.NewJobBaseManager("v0.0.43"),
	}

	tests := []struct {
		name    string
		req     *types.ResourceRequests
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid resource requests",
			req: &types.ResourceRequests{
				Memory:         1024,
				MemoryPerCPU:   256,
				CPUsPerTask:    2,
				TasksPerNode:   4,
				ThreadsPerCore: 1,
			},
			wantErr: false,
		},
		{
			name: "negative memory",
			req: &types.ResourceRequests{
				Memory: -1024,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "negative CPUs per task",
			req: &types.ResourceRequests{
				CPUsPerTask: -2,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "zero values allowed",
			req: &types.ResourceRequests{
				Memory:       0,
				CPUsPerTask:  0,
				TasksPerNode: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateResourceRequests(tt.req)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestJobAdapter_ValidateArrayString(t *testing.T) {
	adapter := &JobAdapter{
		JobBaseManager: base.NewJobBaseManager("v0.0.43"),
	}

	tests := []struct {
		name      string
		arrayStr  string
		wantErr   bool
		errMsg    string
	}{
		{
			name:     "valid simple range",
			arrayStr: "1-10",
			wantErr:  false,
		},
		{
			name:     "valid range with step",
			arrayStr: "1-10:2",
			wantErr:  false,
		},
		{
			name:     "valid complex range",
			arrayStr: "1-10,20-30:5,50",
			wantErr:  false,
		},
		{
			name:     "empty string",
			arrayStr: "",
			wantErr:  false, // Empty array string is valid
		},
		{
			name:     "invalid format",
			arrayStr: "1-10-20",
			wantErr:  true,
			errMsg:   "invalid array string format",
		},
		{
			name:     "negative values",
			arrayStr: "-1-10",
			wantErr:  true,
			errMsg:   "array indices must be positive",
		},
		{
			name:     "invalid step",
			arrayStr: "1-10:0",
			wantErr:  true,
			errMsg:   "step value must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateArrayString(tt.arrayStr)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
