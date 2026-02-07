// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_44

import (
	"context"
	stderrors "errors"
	"testing"

	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_44"
	"github.com/jontk/slurm-client/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobAdapter_Hold_Validation(t *testing.T) {
	adapter := NewJobAdapter(&api.ClientWithResponses{})
	tests := []struct {
		name    string
		ctx     context.Context
		req     *types.JobHoldRequest
		wantErr bool
		errCode errors.ErrorCode
	}{
		{
			name:    "nil context",
			ctx:     nil,
			req:     &types.JobHoldRequest{JobId: 123, Hold: true},
			wantErr: true,
			errCode: errors.ErrorCodeValidationFailed,
		},
		{
			name:    "invalid job ID (zero)",
			ctx:     context.Background(),
			req:     &types.JobHoldRequest{JobId: 0, Hold: true},
			wantErr: true,
			errCode: errors.ErrorCodeValidationFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.Hold(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				// Check if it's a ValidationError first
				var valErr *errors.ValidationError
				if stderrors.As(err, &valErr) {
					assert.Equal(t, tt.errCode, valErr.Code)
				} else {
					assert.Equal(t, tt.errCode, errors.GetErrorCode(err))
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestJobAdapter_Hold_ClientNotInitialized(t *testing.T) {
	adapter := NewJobAdapter(nil)
	err := adapter.Hold(context.Background(), &types.JobHoldRequest{
		JobId: 123,
		Hold:  true,
	})
	require.Error(t, err)
	var slurmErr *errors.SlurmError
	if stderrors.As(err, &slurmErr) {
		assert.Equal(t, errors.ErrorCodeClientNotInitialized, slurmErr.Code)
	} else {
		t.Errorf("Expected SlurmError, got %T", err)
	}
}
func TestJobAdapter_Signal_Validation(t *testing.T) {
	adapter := NewJobAdapter(&api.ClientWithResponses{})
	tests := []struct {
		name    string
		ctx     context.Context
		req     *types.JobSignalRequest
		wantErr bool
		errCode errors.ErrorCode
	}{
		{
			name:    "nil context",
			ctx:     nil,
			req:     &types.JobSignalRequest{JobId: 123, Signal: "TERM"},
			wantErr: true,
			errCode: errors.ErrorCodeValidationFailed,
		},
		{
			name:    "invalid job ID",
			ctx:     context.Background(),
			req:     &types.JobSignalRequest{JobId: 0, Signal: "TERM"},
			wantErr: true,
			errCode: errors.ErrorCodeValidationFailed,
		},
		{
			name:    "empty signal",
			ctx:     context.Background(),
			req:     &types.JobSignalRequest{JobId: 123, Signal: ""},
			wantErr: true,
			errCode: errors.ErrorCodeValidationFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.Signal(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				var valErr *errors.ValidationError
				if stderrors.As(err, &valErr) {
					assert.Equal(t, tt.errCode, valErr.Code)
				} else {
					assert.Equal(t, tt.errCode, errors.GetErrorCode(err))
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestJobAdapter_Signal_ClientNotInitialized(t *testing.T) {
	adapter := NewJobAdapter(nil)
	err := adapter.Signal(context.Background(), &types.JobSignalRequest{
		JobId:  123,
		Signal: "TERM",
	})
	require.Error(t, err)
	var slurmErr *errors.SlurmError
	if stderrors.As(err, &slurmErr) {
		assert.Equal(t, errors.ErrorCodeClientNotInitialized, slurmErr.Code)
	} else {
		t.Errorf("Expected SlurmError, got %T", err)
	}
}
func TestJobAdapter_Notify_Validation(t *testing.T) {
	adapter := NewJobAdapter(&api.ClientWithResponses{})
	tests := []struct {
		name    string
		ctx     context.Context
		req     *types.JobNotifyRequest
		wantErr bool
		errCode errors.ErrorCode
	}{
		{
			name:    "nil context",
			ctx:     nil,
			req:     &types.JobNotifyRequest{JobId: 123, Message: "test"},
			wantErr: true,
			errCode: errors.ErrorCodeValidationFailed,
		},
		{
			name:    "invalid job ID",
			ctx:     context.Background(),
			req:     &types.JobNotifyRequest{JobId: 0, Message: "test"},
			wantErr: true,
			errCode: errors.ErrorCodeValidationFailed,
		},
		{
			name:    "empty message",
			ctx:     context.Background(),
			req:     &types.JobNotifyRequest{JobId: 123, Message: ""},
			wantErr: true,
			errCode: errors.ErrorCodeValidationFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.Notify(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				var valErr *errors.ValidationError
				if stderrors.As(err, &valErr) {
					assert.Equal(t, tt.errCode, valErr.Code)
				} else {
					assert.Equal(t, tt.errCode, errors.GetErrorCode(err))
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestJobAdapter_Notify_NotSupported(t *testing.T) {
	adapter := NewJobAdapter(&api.ClientWithResponses{})
	// Job notification is not supported in REST API, should return unsupported operation error
	err := adapter.Notify(context.Background(), &types.JobNotifyRequest{
		JobId:   123,
		Message: "test message",
	})
	require.Error(t, err)
	var slurmErr *errors.SlurmError
	if stderrors.As(err, &slurmErr) {
		assert.Equal(t, errors.ErrorCodeUnsupportedOperation, slurmErr.Code)
	} else {
		t.Errorf("Expected SlurmError, got %T", err)
	}
	assert.Contains(t, err.Error(), "not supported via REST API")
	assert.Contains(t, err.Error(), "scontrol notify")
}
func TestJobAdapter_Requeue_Validation(t *testing.T) {
	adapter := NewJobAdapter(&api.ClientWithResponses{})
	tests := []struct {
		name    string
		ctx     context.Context
		jobID   int32
		wantErr bool
		errCode errors.ErrorCode
	}{
		{
			name:    "nil context",
			ctx:     nil,
			jobID:   123,
			wantErr: true,
			errCode: errors.ErrorCodeValidationFailed,
		},
		{
			name:    "invalid job ID (zero)",
			ctx:     context.Background(),
			jobID:   0,
			wantErr: true,
			errCode: errors.ErrorCodeValidationFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.Requeue(tt.ctx, tt.jobID)
			if tt.wantErr {
				require.Error(t, err)
				var valErr *errors.ValidationError
				if stderrors.As(err, &valErr) {
					assert.Equal(t, tt.errCode, valErr.Code)
				} else {
					assert.Equal(t, tt.errCode, errors.GetErrorCode(err))
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestJobAdapter_Requeue_ClientNotInitialized(t *testing.T) {
	adapter := NewJobAdapter(nil)
	err := adapter.Requeue(context.Background(), 123)
	require.Error(t, err)
	var slurmErr *errors.SlurmError
	if stderrors.As(err, &slurmErr) {
		assert.Equal(t, errors.ErrorCodeClientNotInitialized, slurmErr.Code)
	} else {
		t.Errorf("Expected SlurmError, got %T", err)
	}
}

// NOTE: Integration tests with priority and step_id would require a mock server
// and are better suited for end-to-end testing
