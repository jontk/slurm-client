package v0_0_41

import (
	"context"
	"net/http"
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/mock"
)

// MockJobClientWithResponses is a mock for the API client
type MockJobClientWithResponses struct {
	mock.Mock
}

func (m *MockJobClientWithResponses) SlurmV0041GetJobsWithResponse(ctx context.Context, params *api.SlurmV0041GetJobsParams, reqEditors ...api.RequestEditorFn) (*api.SlurmV0041GetJobsResponse, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*api.SlurmV0041GetJobsResponse), args.Error(1)
}

func (m *MockJobClientWithResponses) SlurmV0041GetJobWithResponse(ctx context.Context, jobId string, params *api.SlurmV0041GetJobParams, reqEditors ...api.RequestEditorFn) (*api.SlurmV0041GetJobResponse, error) {
	args := m.Called(ctx, jobId, params)
	return args.Get(0).(*api.SlurmV0041GetJobResponse), args.Error(1)
}

func (m *MockJobClientWithResponses) SlurmV0041PostJobSubmitWithResponse(ctx context.Context, body api.SlurmV0041PostJobSubmitJSONRequestBody, reqEditors ...api.RequestEditorFn) (*api.SlurmV0041PostJobSubmitResponse, error) {
	args := m.Called(ctx, body)
	return args.Get(0).(*api.SlurmV0041PostJobSubmitResponse), args.Error(1)
}

func (m *MockJobClientWithResponses) SlurmV0041DeleteJobWithResponse(ctx context.Context, jobId string, reqEditors ...api.RequestEditorFn) (*api.SlurmV0041DeleteJobResponse, error) {
	args := m.Called(ctx, jobId)
	return args.Get(0).(*api.SlurmV0041DeleteJobResponse), args.Error(1)
}

func TestJobAdapter_ValidateJobCreate(t *testing.T) {
	adapter := &JobAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Job"),
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
			errMsg:  "job creation data is required",
		},
		{
			name: "empty script and command",
			job: &types.JobCreate{
				Script:  "",
				Command: "",
			},
			wantErr: true,
			errMsg:  "either script or command is required",
		},
		{
			name: "valid job with script",
			job: &types.JobCreate{
				Script: "#!/bin/bash\necho 'Hello World'",
				Name:   "test-job",
			},
			wantErr: false,
		},
		{
			name: "valid job with command",
			job: &types.JobCreate{
				Command: "echo 'Hello World'",
				Name:    "test-job",
			},
			wantErr: false,
		},
		{
			name: "negative CPUs",
			job: &types.JobCreate{
				Script: "#!/bin/bash\necho 'test'",
				CPUs:   -1,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "negative nodes",
			job: &types.JobCreate{
				Script: "#!/bin/bash\necho 'test'",
				Nodes:  -1,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "negative time limit",
			job: &types.JobCreate{
				Script:    "#!/bin/bash\necho 'test'",
				TimeLimit: -60,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.validateJobCreate(tt.job)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestJobAdapter_List(t *testing.T) {
	tests := []struct {
		name           string
		opts           *types.JobListOptions
		mockResponse   *api.SlurmV0041GetJobsResponse
		mockError      error
		expectedLen    int
		expectedError  string
		setupMock      func(*MockJobClientWithResponses)
	}{
		{
			name: "successful list with no options",
			opts: nil,
			setupMock: func(m *MockJobClientWithResponses) {
				m.On("SlurmV0041GetJobsWithResponse", mock.Anything, mock.AnythingOfType("*v0_0_41.SlurmV0041GetJobsParams")).
					Return(&api.SlurmV0041GetJobsResponse{
						JSON200: &api.V0041OpenapiJobsResp{
							Jobs: []api.V0041JobInfo{
								{JobId: uint32Ptr(12345), Name: stringPtr("job1")},
								{JobId: uint32Ptr(12346), Name: stringPtr("job2")},
							},
						},
						HTTPResponse: &http.Response{StatusCode: 200},
					}, nil)
			},
			expectedLen: 2,
		},
		{
			name: "successful list with options",
			opts: &types.JobListOptions{
				States: []string{"RUNNING", "PENDING"},
				Users:  []string{"testuser"},
			},
			setupMock: func(m *MockJobClientWithResponses) {
				m.On("SlurmV0041GetJobsWithResponse", mock.Anything, mock.AnythingOfType("*v0_0_41.SlurmV0041GetJobsParams")).
					Return(&api.SlurmV0041GetJobsResponse{
						JSON200: &api.V0041OpenapiJobsResp{
							Jobs: []api.V0041JobInfo{
								{JobId: uint32Ptr(12345), Name: stringPtr("job1"), JobState: stringPtr("RUNNING")},
							},
						},
						HTTPResponse: &http.Response{StatusCode: 200},
					}, nil)
			},
			expectedLen: 1,
		},
		{
			name:          "API error",
			opts:          nil,
			expectedError: "failed to list jobs",
			setupMock: func(m *MockJobClientWithResponses) {
				m.On("SlurmV0041GetJobsWithResponse", mock.Anything, mock.AnythingOfType("*v0_0_41.SlurmV0041GetJobsParams")).
					Return((*api.SlurmV0041GetJobsResponse)(nil), assert.AnError)
			},
		},
		{
			name: "empty response",
			opts: nil,
			setupMock: func(m *MockJobClientWithResponses) {
				m.On("SlurmV0041GetJobsWithResponse", mock.Anything, mock.AnythingOfType("*v0_0_41.SlurmV0041GetJobsParams")).
					Return(&api.SlurmV0041GetJobsResponse{
						JSON200: &api.V0041OpenapiJobsResp{
							Jobs: []api.V0041JobInfo{},
						},
						HTTPResponse: &http.Response{StatusCode: 200},
					}, nil)
			},
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockJobClientWithResponses{}
			tt.setupMock(mockClient)

			adapter := &JobAdapter{
				BaseManager: base.NewBaseManager("v0.0.41", "Job"),
				client:      mockClient,
			}

			result, err := adapter.List(context.Background(), tt.opts)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Len(t, result.Jobs, tt.expectedLen)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestJobAdapter_Get(t *testing.T) {
	tests := []struct {
		name          string
		jobId         string
		mockResponse  *api.SlurmV0041GetJobResponse
		mockError     error
		expectedError string
		setupMock     func(*MockJobClientWithResponses, string)
	}{
		{
			name:  "successful get",
			jobId: "12345",
			setupMock: func(m *MockJobClientWithResponses, jobId string) {
				m.On("SlurmV0041GetJobWithResponse", mock.Anything, jobId, mock.AnythingOfType("*v0_0_41.SlurmV0041GetJobParams")).
					Return(&api.SlurmV0041GetJobResponse{
						JSON200: &api.V0041OpenapiJobsResp{
							Jobs: []api.V0041JobInfo{
								{JobId: uint32Ptr(12345), Name: stringPtr("test-job")},
							},
						},
						HTTPResponse: &http.Response{StatusCode: 200},
					}, nil)
			},
		},
		{
			name:  "empty job ID",
			jobId: "",
			setupMock: func(m *MockJobClientWithResponses, jobId string) {
				// No mock setup needed as validation should fail first
			},
			expectedError: "jobId",
		},
		{
			name:  "API error",
			jobId: "12345",
			setupMock: func(m *MockJobClientWithResponses, jobId string) {
				m.On("SlurmV0041GetJobWithResponse", mock.Anything, jobId, mock.AnythingOfType("*v0_0_41.SlurmV0041GetJobParams")).
					Return((*api.SlurmV0041GetJobResponse)(nil), assert.AnError)
			},
			expectedError: "failed to get job",
		},
		{
			name:  "job not found",
			jobId: "99999",
			setupMock: func(m *MockJobClientWithResponses, jobId string) {
				m.On("SlurmV0041GetJobWithResponse", mock.Anything, jobId, mock.AnythingOfType("*v0_0_41.SlurmV0041GetJobParams")).
					Return(&api.SlurmV0041GetJobResponse{
						JSON200: &api.V0041OpenapiJobsResp{
							Jobs: []api.V0041JobInfo{},
						},
						HTTPResponse: &http.Response{StatusCode: 200},
					}, nil)
			},
			expectedError: "job 99999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockJobClientWithResponses{}
			tt.setupMock(mockClient, tt.jobId)

			adapter := &JobAdapter{
				BaseManager: base.NewBaseManager("v0.0.41", "Job"),
				client:      mockClient,
			}

			result, err := adapter.Get(context.Background(), tt.jobId)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, int32(12345), result.JobID)
			}

			if tt.jobId != "" && tt.expectedError == "" {
				mockClient.AssertExpectations(t)
			}
		})
	}
}

func TestJobAdapter_Submit(t *testing.T) {
	tests := []struct {
		name          string
		job           *types.JobCreate
		mockResponse  *api.SlurmV0041PostJobSubmitResponse
		mockError     error
		expectedError string
		setupMock     func(*MockJobClientWithResponses, *types.JobCreate)
	}{
		{
			name: "successful submit",
			job: &types.JobCreate{
				Script: "#!/bin/bash\necho 'Hello World'",
				Name:   "test-job",
			},
			setupMock: func(m *MockJobClientWithResponses, job *types.JobCreate) {
				m.On("SlurmV0041PostJobSubmitWithResponse", mock.Anything, mock.AnythingOfType("v0_0_41.SlurmV0041PostJobSubmitJSONRequestBody")).
					Return(&api.SlurmV0041PostJobSubmitResponse{
						JSON200: &api.V0041OpenapiJobSubmitResponse{
							JobId: uint32Ptr(12345),
						},
						HTTPResponse: &http.Response{StatusCode: 200},
					}, nil)
			},
		},
		{
			name: "nil job",
			job:  nil,
			setupMock: func(m *MockJobClientWithResponses, job *types.JobCreate) {
				// No mock setup needed as validation should fail first
			},
			expectedError: "job creation data is required",
		},
		{
			name: "empty script and command",
			job: &types.JobCreate{
				Script:  "",
				Command: "",
			},
			setupMock: func(m *MockJobClientWithResponses, job *types.JobCreate) {
				// No mock setup needed as validation should fail first
			},
			expectedError: "either script or command is required",
		},
		{
			name: "API error",
			job: &types.JobCreate{
				Script: "#!/bin/bash\necho 'test'",
			},
			setupMock: func(m *MockJobClientWithResponses, job *types.JobCreate) {
				m.On("SlurmV0041PostJobSubmitWithResponse", mock.Anything, mock.AnythingOfType("v0_0_41.SlurmV0041PostJobSubmitJSONRequestBody")).
					Return((*api.SlurmV0041PostJobSubmitResponse)(nil), assert.AnError)
			},
			expectedError: "failed to submit job",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockJobClientWithResponses{}
			tt.setupMock(mockClient, tt.job)

			adapter := &JobAdapter{
				BaseManager: base.NewBaseManager("v0.0.41", "Job"),
				client:      mockClient,
			}

			result, err := adapter.Submit(context.Background(), tt.job)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, int32(12345), result.JobID)
			}

			if tt.job != nil && (tt.job.Script != "" || tt.job.Command != "") && tt.expectedError == "" {
				mockClient.AssertExpectations(t)
			}
		})
	}
}

func TestJobAdapter_Cancel(t *testing.T) {
	tests := []struct {
		name          string
		jobId         string
		mockResponse  *api.SlurmV0041DeleteJobResponse
		mockError     error
		expectedError string
		setupMock     func(*MockJobClientWithResponses, string)
	}{
		{
			name:  "successful cancel",
			jobId: "12345",
			setupMock: func(m *MockJobClientWithResponses, jobId string) {
				m.On("SlurmV0041DeleteJobWithResponse", mock.Anything, jobId).
					Return(&api.SlurmV0041DeleteJobResponse{
						HTTPResponse: &http.Response{StatusCode: 200},
					}, nil)
			},
		},
		{
			name:  "empty job ID",
			jobId: "",
			setupMock: func(m *MockJobClientWithResponses, jobId string) {
				// No mock setup needed as validation should fail first
			},
			expectedError: "jobId",
		},
		{
			name:  "API error",
			jobId: "12345",
			setupMock: func(m *MockJobClientWithResponses, jobId string) {
				m.On("SlurmV0041DeleteJobWithResponse", mock.Anything, jobId).
					Return((*api.SlurmV0041DeleteJobResponse)(nil), assert.AnError)
			},
			expectedError: "failed to cancel job",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockJobClientWithResponses{}
			tt.setupMock(mockClient, tt.jobId)

			adapter := &JobAdapter{
				BaseManager: base.NewBaseManager("v0.0.41", "Job"),
				client:      mockClient,
			}

			err := adapter.Cancel(context.Background(), tt.jobId)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			if tt.jobId != "" && tt.expectedError == "" {
				mockClient.AssertExpectations(t)
			}
		})
	}
}

// Test error conditions and edge cases
func TestJobAdapter_ErrorConditions(t *testing.T) {
	t.Run("nil context", func(t *testing.T) {
		adapter := &JobAdapter{
			BaseManager: base.NewBaseManager("v0.0.41", "Job"),
		}

		_, err := adapter.List(nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "context")
	})

	t.Run("nil client", func(t *testing.T) {
		adapter := &JobAdapter{
			BaseManager: base.NewBaseManager("v0.0.41", "Job"),
			client:      nil,
		}

		_, err := adapter.List(context.Background(), nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "client")
	})
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func uint32Ptr(i uint32) *uint32 {
	return &i
}

func intPtr(i int) *int {
	return &i
}