// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// Mock client for testing
type MockReservationClient struct {
	mock.Mock
}

func (m *MockReservationClient) SlurmV0043GetReservationsWithResponse(ctx context.Context, params *api.SlurmV0043GetReservationsParams, reqEditors ...api.RequestEditorFn) (*api.SlurmV0043GetReservationsResponse, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmV0043GetReservationsResponse), args.Error(1)
}

func (m *MockReservationClient) SlurmV0043GetReservationWithResponse(ctx context.Context, reservationName string, params *api.SlurmV0043GetReservationParams, reqEditors ...api.RequestEditorFn) (*api.SlurmV0043GetReservationResponse, error) {
	args := m.Called(ctx, reservationName, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmV0043GetReservationResponse), args.Error(1)
}

func (m *MockReservationClient) SlurmV0043PostReservationWithResponse(ctx context.Context, body api.SlurmV0043PostReservationJSONRequestBody, reqEditors ...api.RequestEditorFn) (*api.SlurmV0043PostReservationResponse, error) {
	args := m.Called(ctx, body)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmV0043PostReservationResponse), args.Error(1)
}

func (m *MockReservationClient) SlurmV0043DeleteReservationWithResponse(ctx context.Context, reservationName string, reqEditors ...api.RequestEditorFn) (*api.SlurmV0043DeleteReservationResponse, error) {
	args := m.Called(ctx, reservationName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmV0043DeleteReservationResponse), args.Error(1)
}

func TestReservationAdapter_List(t *testing.T) {
	tests := []struct {
		name          string
		opts          *types.ReservationListOptions
		mockResponse  *api.SlurmV0043GetReservationsResponse
		mockError     error
		expectedError bool
		expectedCount int
	}{
		{
			name: "successful list",
			opts: &types.ReservationListOptions{},
			mockResponse: &api.SlurmV0043GetReservationsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiReservationResp{
					Reservations: []api.V0043ReservationInfo{
						{
							Name:      ptrString("maint-window"),
							StartTime: &api.V0043Uint64NoVal{Number: ptrUint64(uint64(time.Now().Unix()))},
							EndTime:   &api.V0043Uint64NoVal{Number: ptrUint64(uint64(time.Now().Add(2 * time.Hour).Unix()))},
							NodeList:  ptrString("node[1-10]"),
							Users:     ptrString("admin"),
						},
						{
							Name:      ptrString("gpu-reservation"),
							StartTime: &api.V0043Uint64NoVal{Number: ptrUint64(uint64(time.Now().Unix()))},
							EndTime:   &api.V0043Uint64NoVal{Number: ptrUint64(uint64(time.Now().Add(4 * time.Hour).Unix()))},
							NodeList:  ptrString("gpu[1-4]"),
							Accounts:  ptrString("research"),
						},
					},
				},
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name: "list with filter",
			opts: &types.ReservationListOptions{
				Names: []string{"maint-window"},
			},
			mockResponse: &api.SlurmV0043GetReservationsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiReservationResp{
					Reservations: []api.V0043ReservationInfo{
						{
							Name: ptrString("maint-window"),
						},
					},
				},
			},
			expectedError: false,
			expectedCount: 1,
		},
		{
			name:          "nil context",
			opts:          nil,
			expectedError: true,
		},
		{
			name: "API error",
			opts: &types.ReservationListOptions{},
			mockResponse: &api.SlurmV0043GetReservationsResponse{
				HTTPResponse: &http.Response{StatusCode: 500},
			},
			expectedError: true,
		},
		{
			name:          "network error",
			opts:          &types.ReservationListOptions{},
			mockError:     fmt.Errorf("connection refused"),
			expectedError: true,
		},
		{
			name: "empty response",
			opts: &types.ReservationListOptions{},
			mockResponse: &api.SlurmV0043GetReservationsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      nil,
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockReservationClient{}
			adapter := &ReservationAdapter{
				client:      mockClient,
				BaseManager: NewReservationAdapter(nil).BaseManager,
			}

			ctx := context.Background()
			if tt.name == "nil context" {
				ctx = nil
			}

			if tt.mockResponse != nil || tt.mockError != nil {
				mockClient.On("SlurmV0043GetReservationsWithResponse", mock.Anything, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			result, err := adapter.List(ctx, tt.opts)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Reservations, tt.expectedCount)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestReservationAdapter_Get(t *testing.T) {
	tests := []struct {
		name            string
		reservationName string
		mockResponse    *api.SlurmV0043GetReservationResponse
		mockError       error
		expectedError   bool
		expectedName    string
	}{
		{
			name:            "successful get",
			reservationName: "maint-window",
			mockResponse: &api.SlurmV0043GetReservationResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiReservationResp{
					Reservations: []api.V0043ReservationInfo{
						{
							Name:      ptrString("maint-window"),
							StartTime: &api.V0043Uint64NoVal{Number: ptrUint64(uint64(time.Now().Unix()))},
							EndTime:   &api.V0043Uint64NoVal{Number: ptrUint64(uint64(time.Now().Add(2 * time.Hour).Unix()))},
							NodeList:  ptrString("node[1-10]"),
							NodeCount: ptrInt32(10),
							Users:     ptrString("admin"),
							State:     ptrString("ACTIVE"),
						},
					},
				},
			},
			expectedError: false,
			expectedName:  "maint-window",
		},
		{
			name:            "reservation not found",
			reservationName: "nonexistent",
			mockResponse: &api.SlurmV0043GetReservationResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0043OpenapiReservationResp{
					Reservations: []api.V0043ReservationInfo{},
				},
			},
			expectedError: true,
		},
		{
			name:            "API error",
			reservationName: "maint-window",
			mockResponse: &api.SlurmV0043GetReservationResponse{
				HTTPResponse: &http.Response{StatusCode: 404},
			},
			expectedError: true,
		},
		{
			name:            "network error",
			reservationName: "maint-window",
			mockError:       fmt.Errorf("timeout"),
			expectedError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockReservationClient{}
			adapter := &ReservationAdapter{
				client:      mockClient,
				BaseManager: NewReservationAdapter(nil).BaseManager,
			}

			mockClient.On("SlurmV0043GetReservationWithResponse", mock.Anything, tt.reservationName, mock.Anything).
				Return(tt.mockResponse, tt.mockError)

			result, err := adapter.Get(context.Background(), tt.reservationName)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedName, result.Name)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestReservationAdapter_Create(t *testing.T) {
	now := time.Now()
	later := now.Add(2 * time.Hour)

	tests := []struct {
		name          string
		reservation   *types.ReservationCreate
		mockResponse  *api.SlurmV0043PostReservationResponse
		mockError     error
		expectedError bool
	}{
		{
			name: "successful create",
			reservation: &types.ReservationCreate{
				Name:      "new-reservation",
				StartTime: now,
				EndTime:   later,
				NodeList:  "node[1-5]",
				Users:     []string{"user1", "user2"},
			},
			mockResponse: &api.SlurmV0043PostReservationResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      &api.V0043OpenapiResp{},
			},
			expectedError: false,
		},
		{
			name: "create with accounts",
			reservation: &types.ReservationCreate{
				Name:      "account-reservation",
				StartTime: now,
				EndTime:   later,
				NodeList:  "node[6-10]",
				Accounts:  []string{"research", "engineering"},
			},
			mockResponse: &api.SlurmV0043PostReservationResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      &api.V0043OpenapiResp{},
			},
			expectedError: false,
		},
		{
			name: "API error",
			reservation: &types.ReservationCreate{
				Name:      "new-reservation",
				StartTime: now,
				EndTime:   later,
				NodeList:  "node[1-5]",
			},
			mockResponse: &api.SlurmV0043PostReservationResponse{
				HTTPResponse: &http.Response{StatusCode: 400},
			},
			expectedError: true,
		},
		{
			name: "network error",
			reservation: &types.ReservationCreate{
				Name:      "new-reservation",
				StartTime: now,
				EndTime:   later,
				NodeList:  "node[1-5]",
			},
			mockError:     fmt.Errorf("connection failed"),
			expectedError: true,
		},
		{
			name:          "nil reservation",
			reservation:   nil,
			expectedError: true,
		},
		{
			name: "missing required fields",
			reservation: &types.ReservationCreate{
				Name: "incomplete",
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockReservationClient{}
			adapter := &ReservationAdapter{
				client:      mockClient,
				BaseManager: NewReservationAdapter(nil).BaseManager,
			}

			if tt.reservation != nil && tt.reservation.Name != "" && tt.reservation.NodeList != "" {
				mockClient.On("SlurmV0043PostReservationWithResponse", mock.Anything, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			result, err := adapter.Create(context.Background(), tt.reservation)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.reservation != nil {
					assert.Equal(t, tt.reservation.Name, result.Name)
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestReservationAdapter_Update(t *testing.T) {
	later := time.Now().Add(4 * time.Hour)

	tests := []struct {
		name            string
		reservationName string
		update          *types.ReservationUpdate
		mockResponse    *api.SlurmV0043PostReservationResponse
		mockError       error
		expectedError   bool
	}{
		{
			name:            "successful update",
			reservationName: "maint-window",
			update: &types.ReservationUpdate{
				EndTime:  &later,
				NodeList: ptrString("node[1-20]"),
				Users:    []string{"admin", "operator"},
			},
			mockResponse: &api.SlurmV0043PostReservationResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      &api.V0043OpenapiResp{},
			},
			expectedError: false,
		},
		{
			name:            "update end time only",
			reservationName: "maint-window",
			update: &types.ReservationUpdate{
				EndTime: &later,
			},
			mockResponse: &api.SlurmV0043PostReservationResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      &api.V0043OpenapiResp{},
			},
			expectedError: false,
		},
		{
			name:            "API error",
			reservationName: "maint-window",
			update: &types.ReservationUpdate{
				NodeList: ptrString("invalid-nodes"),
			},
			mockResponse: &api.SlurmV0043PostReservationResponse{
				HTTPResponse: &http.Response{StatusCode: 400},
			},
			expectedError: true,
		},
		{
			name:            "nil update",
			reservationName: "maint-window",
			update:          nil,
			expectedError:   true,
		},
		{
			name:            "empty reservation name",
			reservationName: "",
			update:          &types.ReservationUpdate{},
			expectedError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockReservationClient{}
			adapter := &ReservationAdapter{
				client:      mockClient,
				BaseManager: NewReservationAdapter(nil).BaseManager,
			}

			if tt.reservationName != "" && tt.update != nil {
				mockClient.On("SlurmV0043PostReservationWithResponse", mock.Anything, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			err := adapter.Update(context.Background(), tt.reservationName, tt.update)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestReservationAdapter_Delete(t *testing.T) {
	tests := []struct {
		name            string
		reservationName string
		mockResponse    *api.SlurmV0043DeleteReservationResponse
		mockError       error
		expectedError   bool
	}{
		{
			name:            "successful delete",
			reservationName: "old-reservation",
			mockResponse: &api.SlurmV0043DeleteReservationResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
			},
			expectedError: false,
		},
		{
			name:            "reservation not found",
			reservationName: "nonexistent",
			mockResponse: &api.SlurmV0043DeleteReservationResponse{
				HTTPResponse: &http.Response{StatusCode: 404},
			},
			expectedError: true,
		},
		{
			name:            "network error",
			reservationName: "old-reservation",
			mockError:       fmt.Errorf("network error"),
			expectedError:   true,
		},
		{
			name:            "empty reservation name",
			reservationName: "",
			expectedError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockReservationClient{}
			adapter := &ReservationAdapter{
				client:      mockClient,
				BaseManager: NewReservationAdapter(nil).BaseManager,
			}

			if tt.reservationName != "" {
				mockClient.On("SlurmV0043DeleteReservationWithResponse", mock.Anything, tt.reservationName).
					Return(tt.mockResponse, tt.mockError)
			}

			err := adapter.Delete(context.Background(), tt.reservationName)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestReservationAdapter_validateReservationCreate(t *testing.T) {
	adapter := NewReservationAdapter(nil)
	now := time.Now()
	later := now.Add(2 * time.Hour)
	past := now.Add(-2 * time.Hour)

	tests := []struct {
		name          string
		reservation   *types.ReservationCreate
		expectedError bool
		errorContains string
	}{
		{
			name: "valid reservation",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: now,
				EndTime:   later,
				NodeList:  "node[1-5]",
			},
			expectedError: false,
		},
		{
			name:          "nil reservation",
			reservation:   nil,
			expectedError: true,
			errorContains: "reservation is required",
		},
		{
			name: "empty name",
			reservation: &types.ReservationCreate{
				Name:      "",
				StartTime: now,
				EndTime:   later,
				NodeList:  "node[1-5]",
			},
			expectedError: true,
			errorContains: "reservation name is required",
		},
		{
			name: "missing node list",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: now,
				EndTime:   later,
				NodeList:  "",
			},
			expectedError: true,
			errorContains: "node list is required",
		},
		{
			name: "end time before start time",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: later,
				EndTime:   now,
				NodeList:  "node[1-5]",
			},
			expectedError: true,
			errorContains: "end time must be after start time",
		},
		{
			name: "start time in past",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: past,
				EndTime:   now,
				NodeList:  "node[1-5]",
			},
			expectedError: true,
			errorContains: "start time cannot be in the past",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.validateReservationCreate(tt.reservation)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReservationAdapter_convertAPIReservationToCommon(t *testing.T) {
	adapter := NewReservationAdapter(nil)
	now := time.Now()

	tests := []struct {
		name            string
		apiReservation  api.V0043ReservationInfo
		expected        *types.Reservation
	}{
		{
			name: "full reservation",
			apiReservation: api.V0043ReservationInfo{
				Name:      ptrString("test-reservation"),
				StartTime: &api.V0043Uint64NoVal{Number: ptrUint64(uint64(now.Unix()))},
				EndTime:   &api.V0043Uint64NoVal{Number: ptrUint64(uint64(now.Add(2 * time.Hour).Unix()))},
				NodeList:  ptrString("node[1-10]"),
				NodeCount: ptrInt32(10),
				CoreCount: ptrInt32(160),
				Users:     ptrString("user1,user2"),
				Accounts:  ptrString("account1,account2"),
				State:     ptrString("ACTIVE"),
				Flags:     &[]api.V0043ReservationInfoFlags{"MAINT", "FLEX"},
			},
			expected: &types.Reservation{
				Name:      "test-reservation",
				StartTime: time.Unix(now.Unix(), 0),
				EndTime:   time.Unix(now.Add(2*time.Hour).Unix(), 0),
				NodeList:  "node[1-10]",
				NodeCount: 10,
				CoreCount: 160,
				Users:     []string{"user1", "user2"},
				Accounts:  []string{"account1", "account2"},
				State:     "ACTIVE",
				Flags:     []string{"MAINT", "FLEX"},
			},
		},
		{
			name: "minimal reservation",
			apiReservation: api.V0043ReservationInfo{
				Name:     ptrString("minimal"),
				NodeList: ptrString("node1"),
			},
			expected: &types.Reservation{
				Name:     "minimal",
				NodeList: "node1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.convertAPIReservationToCommon(tt.apiReservation)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.NodeList, result.NodeList)
			assert.Equal(t, tt.expected.NodeCount, result.NodeCount)
			assert.Equal(t, tt.expected.State, result.State)
			assert.Equal(t, tt.expected.Users, result.Users)
			assert.Equal(t, tt.expected.Accounts, result.Accounts)
			assert.Equal(t, tt.expected.Flags, result.Flags)
		})
	}
}

