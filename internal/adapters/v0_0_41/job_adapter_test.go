// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

import (
	"context"
	"testing"
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobAdapter_ValidateJobSubmit(t *testing.T) {
	adapter := &JobAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Job"),
	}

	tests := []struct {
		name    string
		job     *types.JobSubmit
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil job",
			job:     nil,
			wantErr: true,
			errMsg:  "job submission data is required",
		},
		{
			name: "empty script",
			job: &types.JobSubmit{
				Script: "",
			},
			wantErr: true,
			errMsg:  "job script is required",
		},
		{
			name: "invalid cpu count",
			job: &types.JobSubmit{
				Script: "#!/bin/bash\necho hello",
				CPUs:   -1,
			},
			wantErr: true,
			errMsg:  "CPUs must be positive",
		},
		{
			name: "invalid node count",
			job: &types.JobSubmit{
				Script: "#!/bin/bash\necho hello",
				Nodes:  -1,
			},
			wantErr: true,
			errMsg:  "nodes must be positive",
		},
		{
			name: "invalid memory",
			job: &types.JobSubmit{
				Script: "#!/bin/bash\necho hello",
				Memory: -1,
			},
			wantErr: true,
			errMsg:  "memory must be positive",
		},
		{
			name: "valid basic job",
			job: &types.JobSubmit{
				Script: "#!/bin/bash\necho hello world",
			},
			wantErr: false,
		},
		{
			name: "valid complex job",
			job: &types.JobSubmit{
				Script:      "#!/bin/bash\n#SBATCH --job-name=test\necho hello",
				JobName:     "test-job",
				Account:     "physics",
				Partition:   "compute",
				QoS:         "normal",
				CPUs:        4,
				Nodes:       1,
				Memory:      8192,
				TimeLimit:   "01:00:00",
				WorkingDir:  "/home/user",
				Environment: map[string]string{"PATH": "/usr/bin"},
				Dependencies: []string{"afterok:12345"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateJobSubmit(tt.job)
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
		BaseManager: base.NewBaseManager("v0.0.41", "Job"),
	}

	tests := []struct {
		name     string
		input    *types.JobSubmit
		expected *types.JobSubmit
	}{
		{
			name: "apply defaults to minimal job",
			input: &types.JobSubmit{
				Script: "#!/bin/bash\necho hello",
			},
			expected: &types.JobSubmit{
				Script:      "#!/bin/bash\necho hello",
				JobName:     "",
				Account:     "",
				Partition:   "",
				QoS:         "",
				CPUs:        1,
				Nodes:       1,
				Memory:      0,
				TimeLimit:   "",
				WorkingDir:  "",
				Environment: map[string]string{},
				Dependencies: []string{},
				ArraySpec:   "",
				StandardOut: "",
				StandardErr: "",
			},
		},
		{
			name: "preserve existing values",
			input: &types.JobSubmit{
				Script:      "#!/bin/bash\necho hello",
				JobName:     "custom-job",
				Account:     "physics",
				Partition:   "gpu",
				CPUs:        8,
				Memory:      16384,
				TimeLimit:   "02:00:00",
				WorkingDir:  "/scratch",
				Environment: map[string]string{"OMP_NUM_THREADS": "8"},
			},
			expected: &types.JobSubmit{
				Script:      "#!/bin/bash\necho hello",
				JobName:     "custom-job",
				Account:     "physics",
				Partition:   "gpu",
				QoS:         "",
				CPUs:        8,
				Nodes:       1,
				Memory:      16384,
				TimeLimit:   "02:00:00",
				WorkingDir:  "/scratch",
				Environment: map[string]string{"OMP_NUM_THREADS": "8"},
				Dependencies: []string{},
				ArraySpec:   "",
				StandardOut: "",
				StandardErr: "",
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
		BaseManager: base.NewBaseManager("v0.0.41", "Job"),
	}

	jobs := []types.Job{
		{
			JobID:     12345,
			JobName:   "test-job-1",
			Account:   "physics",
			Partition: "compute",
			QoS:       "normal",
			User:      "user1",
			State:     "RUNNING",
			CPUs:      4,
			Nodes:     1,
			SubmitTime: time.Now().Add(-2 * time.Hour),
			StartTime:  time.Now().Add(-1 * time.Hour),
		},
		{
			JobID:     12346,
			JobName:   "test-job-2",
			Account:   "chemistry",
			Partition: "gpu",
			QoS:       "high",
			User:      "user2",
			State:     "PENDING",
			CPUs:      8,
			Nodes:     2,
			SubmitTime: time.Now().Add(-1 * time.Hour),
		},
		{
			JobID:     12347,
			JobName:   "batch-job",
			Account:   "physics",
			Partition: "compute",
			QoS:       "normal",
			User:      "user1",
			State:     "COMPLETED",
			CPUs:      2,
			Nodes:     1,
			SubmitTime: time.Now().Add(-3 * time.Hour),
			StartTime:  time.Now().Add(-3 * time.Hour),
			EndTime:    time.Now().Add(-30 * time.Minute),
		},
	}

	tests := []struct {
		name     string
		opts     *types.JobListOptions
		expected []int // expected job IDs
	}{
		{
			name:     "no filters",
			opts:     &types.JobListOptions{},
			expected: []int{12345, 12346, 12347},
		},
		{
			name: "filter by job IDs",
			opts: &types.JobListOptions{
				JobIDs: []int{12345, 12347},
			},
			expected: []int{12345, 12347},
		},
		{
			name: "filter by accounts",
			opts: &types.JobListOptions{
				Accounts: []string{"physics"},
			},
			expected: []int{12345, 12347},
		},
		{
			name: "filter by users",
			opts: &types.JobListOptions{
				Users: []string{"user1"},
			},
			expected: []int{12345, 12347},
		},
		{
			name: "filter by states",
			opts: &types.JobListOptions{
				States: []string{"RUNNING", "PENDING"},
			},
			expected: []int{12345, 12346},
		},
		{
			name: "filter by partitions",
			opts: &types.JobListOptions{
				Partitions: []string{"compute"},
			},
			expected: []int{12345, 12347},
		},
		{
			name: "filter by QoS",
			opts: &types.JobListOptions{
				QoSNames: []string{"normal"},
			},
			expected: []int{12345, 12347},
		},
		{
			name: "combined filters",
			opts: &types.JobListOptions{
				Accounts: []string{"physics"},
				States:   []string{"RUNNING"},
			},
			expected: []int{12345},
		},
		{
			name: "no matches",
			opts: &types.JobListOptions{
				Accounts: []string{"nonexistent"},
			},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.FilterJobList(jobs, tt.opts)
			resultIDs := make([]int, len(result))
			for i, job := range result {
				resultIDs[i] = job.JobID
			}
			assert.Equal(t, tt.expected, resultIDs)
		})
	}
}

func TestJobAdapter_ValidateJobScript(t *testing.T) {
	adapter := &JobAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Job"),
	}

	tests := []struct {
		name    string
		script  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty script",
			script:  "",
			wantErr: true,
			errMsg:  "script cannot be empty",
		},
		{
			name:    "no shebang",
			script:  "echo hello",
			wantErr: true,
			errMsg:  "script must start with shebang",
		},
		{
			name:    "valid bash script",
			script:  "#!/bin/bash\necho hello world",
			wantErr: false,
		},
		{
			name:    "valid sh script",
			script:  "#!/bin/sh\necho hello",
			wantErr: false,
		},
		{
			name:    "valid python script",
			script:  "#!/usr/bin/env python3\nprint('hello')",
			wantErr: false,
		},
		{
			name:    "script with SBATCH directives",
			script:  "#!/bin/bash\n#SBATCH --job-name=test\n#SBATCH --time=01:00:00\necho hello",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateJobScript(tt.script)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestJobAdapter_ParseJobDependencies(t *testing.T) {
	adapter := &JobAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Job"),
	}

	tests := []struct {
		name         string
		dependencies []string
		expected     map[string][]int
	}{
		{
			name:         "empty dependencies",
			dependencies: []string{},
			expected:     map[string][]int{},
		},
		{
			name:         "single afterok dependency",
			dependencies: []string{"afterok:12345"},
			expected: map[string][]int{
				"afterok": {12345},
			},
		},
		{
			name:         "multiple dependencies same type",
			dependencies: []string{"afterok:12345", "afterok:12346"},
			expected: map[string][]int{
				"afterok": {12345, 12346},
			},
		},
		{
			name:         "multiple dependency types",
			dependencies: []string{"afterok:12345", "afternotok:12346", "after:12347"},
			expected: map[string][]int{
				"afterok":    {12345},
				"afternotok": {12346},
				"after":      {12347},
			},
		},
		{
			name:         "colon-separated job lists",
			dependencies: []string{"afterok:12345:12346:12347"},
			expected: map[string][]int{
				"afterok": {12345, 12346, 12347},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.ParseJobDependencies(tt.dependencies)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJobAdapter_ValidateTimeLimit(t *testing.T) {
	adapter := &JobAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Job"),
	}

	tests := []struct {
		name      string
		timeLimit string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "empty time limit",
			timeLimit: "",
			wantErr:   false, // empty is valid (uses defaults)
		},
		{
			name:      "minutes format",
			timeLimit: "60",
			wantErr:   false,
		},
		{
			name:      "hours:minutes format",
			timeLimit: "01:30",
			wantErr:   false,
		},
		{
			name:      "days-hours:minutes:seconds format",
			timeLimit: "2-12:30:45",
			wantErr:   false,
		},
		{
			name:      "hours:minutes:seconds format",
			timeLimit: "01:30:45",
			wantErr:   false,
		},
		{
			name:      "invalid format",
			timeLimit: "invalid",
			wantErr:   true,
			errMsg:    "invalid time limit format",
		},
		{
			name:      "negative minutes",
			timeLimit: "-30",
			wantErr:   true,
			errMsg:    "time limit cannot be negative",
		},
		{
			name:      "invalid hours",
			timeLimit: "25:00",
			wantErr:   true,
			errMsg:    "invalid hours",
		},
		{
			name:      "invalid minutes",
			timeLimit: "01:65",
			wantErr:   true,
			errMsg:    "invalid minutes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateTimeLimit(tt.timeLimit)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestJobAdapter_CalculateJobPriority(t *testing.T) {
	adapter := &JobAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Job"),
	}

	tests := []struct {
		name             string
		job              *types.Job
		expectedPriority int
	}{
		{
			name: "basic job priority",
			job: &types.Job{
				QoS:        "normal",
				CPUs:       4,
				Nodes:      1,
				SubmitTime: time.Now().Add(-1 * time.Hour),
			},
			expectedPriority: 1000, // base priority
		},
		{
			name: "high QoS job",
			job: &types.Job{
				QoS:        "high",
				CPUs:       8,
				Nodes:      2,
				SubmitTime: time.Now().Add(-2 * time.Hour),
			},
			expectedPriority: 2000, // higher priority for high QoS
		},
		{
			name: "aged job gets priority boost",
			job: &types.Job{
				QoS:        "normal",
				CPUs:       2,
				Nodes:      1,
				SubmitTime: time.Now().Add(-24 * time.Hour), // old job
			},
			expectedPriority: 1500, // age boost
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			priority := adapter.CalculateJobPriority(tt.job)
			assert.Equal(t, tt.expectedPriority, priority)
		})
	}
}

func TestJobAdapter_ValidateContext(t *testing.T) {
	adapter := &JobAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Job"),
	}

	tests := []struct {
		name    string
		ctx     context.Context
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil context",
			ctx:     nil,
			wantErr: true,
			errMsg:  "context is required",
		},
		{
			name:    "valid context",
			ctx:     context.Background(),
			wantErr: false,
		},
		{
			name:    "context with timeout",
			ctx:     context.TODO(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateContext(tt.ctx)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
