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
			name: "negative node count",
			job: &types.JobCreate{
				Command: "echo 'test'",
				CPUs:    2,
				Nodes:   -1,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "valid basic job with command",
			job: &types.JobCreate{
				Command: "echo 'test'",
				CPUs:    1,
				Nodes:   1,
			},
			wantErr: false,
		},
		{
			name: "valid basic job with script",
			job: &types.JobCreate{
				Script: "#!/bin/bash\necho 'test'",
				CPUs:   2,
				Nodes:  1,
			},
			wantErr: false,
		},
		{
			name: "valid complex job",
			job: &types.JobCreate{
				Script:           "#!/bin/bash\necho 'complex test'",
				Name:             "test-job",
				Partition:        "default",
				QoS:              "normal",
				Account:          "test-account",
				TimeLimit:        120,
				CPUs:             4,
				Nodes:            2,
				Tasks:            4,
				Environment:      map[string]string{"TEST_VAR": "test_value"},
				WorkingDirectory: "/tmp/test",
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

	validJobUpdate := &types.JobUpdate{
		JobID:     intPtr(12345),
		TimeLimit: intPtr(120),
		Priority:  intPtr(1000),
	}

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
			errMsg:  "job update data is required",
		},
		{
			name: "missing job ID",
			update: &types.JobUpdate{
				TimeLimit: intPtr(120),
			},
			wantErr: true,
			errMsg:  "job ID is required",
		},
		{
			name: "zero job ID",
			update: &types.JobUpdate{
				JobID:     intPtr(0),
				TimeLimit: intPtr(120),
			},
			wantErr: true,
			errMsg:  "job ID must be greater than 0",
		},
		{
			name: "negative time limit",
			update: &types.JobUpdate{
				JobID:     intPtr(12345),
				TimeLimit: intPtr(-1),
			},
			wantErr: true,
			errMsg:  "time limit must be non-negative",
		},
		{
			name: "negative priority",
			update: &types.JobUpdate{
				JobID:    intPtr(12345),
				Priority: intPtr(-1),
			},
			wantErr: true,
			errMsg:  "priority must be non-negative",
		},
		{
			name:    "valid update",
			update:  validJobUpdate,
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

	tests := []struct {
		name     string
		input    *types.JobSubmit
		expected *types.JobSubmit
	}{
		{
			name: "apply defaults to minimal job",
			input: &types.JobSubmit{
				Script: "#!/bin/bash\necho 'test'",
			},
			expected: &types.JobSubmit{
				Script:    "#!/bin/bash\necho 'test'",
				TimeLimit: 60,  // Default 1 hour
				NodeCount: 1,   // Default 1 node
				CPUCount:  1,   // Default 1 CPU
				Partition: "",  // No default partition
				QoS:       "",  // No default QoS
				Account:   "",  // No default account
			},
		},
		{
			name: "preserve existing values",
			input: &types.JobSubmit{
				Script:    "#!/bin/bash\necho 'test'",
				JobName:   "custom-job",
				Partition: "gpu",
				QoS:       "high",
				Account:   "custom-account",
				TimeLimit: 240,
				NodeCount: 4,
				CPUCount:  8,
			},
			expected: &types.JobSubmit{
				Script:    "#!/bin/bash\necho 'test'",
				JobName:   "custom-job",
				Partition: "gpu",
				QoS:       "high",
				Account:   "custom-account",
				TimeLimit: 240,
				NodeCount: 4,
				CPUCount:  8,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.ApplyJobDefaults(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJobBaseManager_FilterJobList(t *testing.T) {
	manager := NewJobBaseManager("v0.0.43")

	jobs := []types.Job{
		{
			JobID:     12345,
			JobName:   "job1",
			User:      "user1",
			Account:   "account1",
			Partition: "partition1",
			QoS:       "normal",
			State:     "RUNNING",
		},
		{
			JobID:     12346,
			JobName:   "job2",
			User:      "user2",
			Account:   "account1",
			Partition: "partition2",
			QoS:       "high",
			State:     "PENDING",
		},
		{
			JobID:     12347,
			JobName:   "job3",
			User:      "user1",
			Account:   "account2",
			Partition: "partition1",
			QoS:       "normal",
			State:     "COMPLETED",
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
			name: "filter by user",
			opts: &types.JobListOptions{
				Users: []string{"user1"},
			},
			expected: []int{12345, 12347},
		},
		{
			name: "filter by account",
			opts: &types.JobListOptions{
				Accounts: []string{"account1"},
			},
			expected: []int{12345, 12346},
		},
		{
			name: "filter by partition",
			opts: &types.JobListOptions{
				Partitions: []string{"partition2"},
			},
			expected: []int{12346},
		},
		{
			name: "filter by QoS",
			opts: &types.JobListOptions{
				QoSList: []string{"high"},
			},
			expected: []int{12346},
		},
		{
			name: "filter by state",
			opts: &types.JobListOptions{
				States: []string{"RUNNING", "PENDING"},
			},
			expected: []int{12345, 12346},
		},
		{
			name: "combined filters",
			opts: &types.JobListOptions{
				Users:    []string{"user1"},
				States:   []string{"RUNNING"},
				Accounts: []string{"account1"},
			},
			expected: []int{12345},
		},
		{
			name: "no matches",
			opts: &types.JobListOptions{
				Users:  []string{"nonexistent"},
				States: []string{"FAILED"},
			},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.FilterJobList(jobs, tt.opts)
			resultIDs := make([]int, len(result))
			for i, job := range result {
				resultIDs[i] = job.JobID
			}
			assert.Equal(t, tt.expected, resultIDs)
		})
	}
}

func TestJobBaseManager_ValidateJobID(t *testing.T) {
	manager := NewJobBaseManager("v0.0.43")

	tests := []struct {
		name    string
		jobID   interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid job ID",
			jobID:   12345,
			wantErr: false,
		},
		{
			name:    "valid job ID pointer",
			jobID:   intPtr(12345),
			wantErr: false,
		},
		{
			name:    "zero job ID",
			jobID:   0,
			wantErr: true,
			errMsg:  "job ID must be greater than 0",
		},
		{
			name:    "negative job ID",
			jobID:   -1,
			wantErr: true,
			errMsg:  "job ID must be greater than 0",
		},
		{
			name:    "nil job ID",
			jobID:   (*int)(nil),
			wantErr: true,
			errMsg:  "job ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateJobID(tt.jobID)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestJobBaseManager_SanitizeJobScript(t *testing.T) {
	manager := NewJobBaseManager("v0.0.43")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic script",
			input:    "#!/bin/bash\necho 'hello'",
			expected: "#!/bin/bash\necho 'hello'",
		},
		{
			name:     "script with dangerous commands",
			input:    "#!/bin/bash\nrm -rf /\necho 'test'",
			expected: "#!/bin/bash\n# POTENTIALLY DANGEROUS: rm -rf /\necho 'test'",
		},
		{
			name:     "script with format command",
			input:    "#!/bin/bash\nformat c:\necho 'test'",
			expected: "#!/bin/bash\n# POTENTIALLY DANGEROUS: format c:\necho 'test'",
		},
		{
			name:     "script with sudo",
			input:    "#!/bin/bash\nsudo reboot\necho 'test'",
			expected: "#!/bin/bash\n# POTENTIALLY DANGEROUS: sudo reboot\necho 'test'",
		},
		{
			name:     "clean script",
			input:    "#!/bin/bash\necho 'safe command'\nls -la",
			expected: "#!/bin/bash\necho 'safe command'\nls -la",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.SanitizeJobScript(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

