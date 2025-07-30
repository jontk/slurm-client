package v0_0_42

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
		BaseManager: base.NewBaseManager("v0.0.42", "Job"),
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
		BaseManager: base.NewBaseManager("v0.0.42", "Job"),
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
		BaseManager: base.NewBaseManager("v0.0.42", "Job"),
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
	}

	tests := []struct {
		name     string
		opts     *types.JobListOptions
		expected []int // expected job IDs
	}{
		{
			name:     "no filters",
			opts:     &types.JobListOptions{},
			expected: []int{12345, 12346},
		},
		{
			name: "filter by accounts",
			opts: &types.JobListOptions{
				Accounts: []string{"physics"},
			},
			expected: []int{12345},
		},
		{
			name: "filter by states",
			opts: &types.JobListOptions{
				States: []string{"RUNNING"},
			},
			expected: []int{12345},
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

func TestJobAdapter_ValidateContext(t *testing.T) {
	adapter := &JobAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Job"),
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