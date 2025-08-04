// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

// Mock client for testing
type MockReservationClient struct {
	mock.Mock
}

func (m *MockReservationClient) SlurmV0042GetReservationsWithResponse(ctx context.Context, params *api.SlurmV0042GetReservationsParams, reqEditors ...api.RequestEditorFn) (*api.SlurmV0042GetReservationsResponse, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmV0042GetReservationsResponse), args.Error(1)
}

func (m *MockReservationClient) SlurmV0042GetReservationWithResponse(ctx context.Context, reservationName string, params *api.SlurmV0042GetReservationParams, reqEditors ...api.RequestEditorFn) (*api.SlurmV0042GetReservationResponse, error) {
	args := m.Called(ctx, reservationName, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmV0042GetReservationResponse), args.Error(1)
}

func (m *MockReservationClient) SlurmV0042PostReservationWithResponse(ctx context.Context, params *api.SlurmV0042PostReservationParams, body api.SlurmV0042PostReservationJSONRequestBody, reqEditors ...api.RequestEditorFn) (*api.SlurmV0042PostReservationResponse, error) {
	args := m.Called(ctx, params, body)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmV0042PostReservationResponse), args.Error(1)
}

func (m *MockReservationClient) SlurmV0042DeleteReservationWithResponse(ctx context.Context, reservationName string, params *api.SlurmV0042DeleteReservationParams, reqEditors ...api.RequestEditorFn) (*api.SlurmV0042DeleteReservationResponse, error) {
	args := m.Called(ctx, reservationName, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.SlurmV0042DeleteReservationResponse), args.Error(1)
}

func TestNewReservationAdapter(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := NewReservationAdapter(client)
	
	assert.NotNil(t, adapter)
	assert.Equal(t, client, adapter.client)
	assert.NotNil(t, adapter.BaseManager)
}

func TestReservationAdapter_List(t *testing.T) {
	tests := []struct {
		name          string
		opts          *types.ReservationListOptions
		mockResponse  *api.SlurmV0042GetReservationsResponse
		mockError     error
		expectedError bool
		expectedCount int
	}{
		{
			name: "successful list",
			opts: &types.ReservationListOptions{},
			mockResponse: &api.SlurmV0042GetReservationsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiReservationResp{
					Reservations: []api.V0042ReservationInfo{
						{
							Name: ptrString("reservation1"),
							Accounts: ptrString("account1,account2"),
							NodeCount: ptrInt32(10),
							StartTime: ptrInt64(1640995200), // 2022-01-01 00:00:00 UTC
							EndTime: ptrInt64(1640998800),   // 2022-01-01 01:00:00 UTC
							Users: ptrString("user1,user2"),
							Partition: ptrString("normal"),
						},
						{
							Name: ptrString("reservation2"),
							Accounts: ptrString("account3"),
							NodeCount: ptrInt32(5),
							StartTime: ptrInt64(1640995200),
							EndTime: ptrInt64(1640998800),
							Users: ptrString("user3"),
							Partition: ptrString("gpu"),
						},
					},
				},
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name: "list with name filter",
			opts: &types.ReservationListOptions{
				Names: []string{"reservation1"},
			},
			mockResponse: &api.SlurmV0042GetReservationsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiReservationResp{
					Reservations: []api.V0042ReservationInfo{
						{
							Name: ptrString("reservation1"),
							Accounts: ptrString("account1"),
							NodeCount: ptrInt32(10),
							StartTime: ptrInt64(1640995200),
							EndTime: ptrInt64(1640998800),
							Users: ptrString("user1"),
						},
						{
							Name: ptrString("reservation2"),
							Accounts: ptrString("account2"),
							NodeCount: ptrInt32(5),
						},
					},
				},
			},
			expectedError: false,
			expectedCount: 1, // Only reservation1 should be returned due to filtering
		},
		{
			name:          "nil context",
			opts:          nil,
			expectedError: true,
		},
		{
			name: "API error",
			opts: &types.ReservationListOptions{},
			mockResponse: &api.SlurmV0042GetReservationsResponse{
				HTTPResponse: &http.Response{StatusCode: 500},
			},
			expectedError: true,
		},
		{
			name:          "network error",
			opts:          &types.ReservationListOptions{},
			mockError:     fmt.Errorf("network error"),
			expectedError: true,
		},
		{
			name: "nil response",
			opts: &types.ReservationListOptions{},
			mockResponse: &api.SlurmV0042GetReservationsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      nil,
			},
			expectedError: true,
		},
		{
			name: "empty reservations list",
			opts: &types.ReservationListOptions{},
			mockResponse: &api.SlurmV0042GetReservationsResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiReservationResp{
					Reservations: []api.V0042ReservationInfo{},
				},
			},
			expectedError: false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockReservationClient{}
			adapter := &ReservationAdapter{
				client:      mockClient,
				BaseManager: base.NewBaseManager("v0.0.42", "Reservation"),
			}

			ctx := context.Background()
			if tt.name == "nil context" {
				ctx = nil
			}

			if tt.mockResponse != nil || tt.mockError != nil {
				mockClient.On("SlurmV0042GetReservationsWithResponse", mock.Anything, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			result, err := adapter.List(ctx, tt.opts)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Reservations, tt.expectedCount)
				
				// Verify filtering worked correctly
				if tt.opts != nil && len(tt.opts.Names) > 0 && !tt.expectedError {
					for _, reservation := range result.Reservations {
						found := false
						for _, name := range tt.opts.Names {
							if reservation.Name == name {
								found = true
								break
							}
						}
						assert.True(t, found, "Reservation %s should match filter", reservation.Name)
					}
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestReservationAdapter_Get(t *testing.T) {
	tests := []struct {
		name            string
		reservationName string
		mockResponse    *api.SlurmV0042GetReservationResponse
		mockError       error
		expectedError   bool
		expectedName    string
	}{
		{
			name:            "successful get",
			reservationName: "test-reservation",
			mockResponse: &api.SlurmV0042GetReservationResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiReservationResp{
					Reservations: []api.V0042ReservationInfo{
						{
							Name: ptrString("test-reservation"),
							Accounts: ptrString("testaccount"),
							NodeCount: ptrInt32(20),
							StartTime: ptrInt64(1640995200),
							EndTime: ptrInt64(1640998800),
							Users: ptrString("testuser"),
							Partition: ptrString("compute"),
							NodeList: ptrString("node[001-020]"),
							MaxStartDelay: ptrInt32(600),
						},
					},
				},
			},
			expectedError: false,
			expectedName:  "test-reservation",
		},
		{
			name:            "reservation not found",
			reservationName: "nonexistent",
			mockResponse: &api.SlurmV0042GetReservationResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiReservationResp{
					Reservations: []api.V0042ReservationInfo{},
				},
			},
			expectedError: true,
		},
		{
			name:            "API error",
			reservationName: "test-reservation",
			mockResponse: &api.SlurmV0042GetReservationResponse{
				HTTPResponse: &http.Response{StatusCode: 404},
			},
			expectedError: true,
		},
		{
			name:            "network error",
			reservationName: "test-reservation",
			mockError:       fmt.Errorf("connection refused"),
			expectedError:   true,
		},
		{
			name:            "empty reservation name",
			reservationName: "",
			expectedError:   true,
		},
		{
			name:            "nil response",
			reservationName: "test-reservation",
			mockResponse: &api.SlurmV0042GetReservationResponse{
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
				BaseManager: base.NewBaseManager("v0.0.42", "Reservation"),
			}

			if tt.reservationName != "" && (tt.mockResponse != nil || tt.mockError != nil) {
				mockClient.On("SlurmV0042GetReservationWithResponse", mock.Anything, tt.reservationName, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

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
	tests := []struct {
		name          string
		reservation   *types.ReservationCreate
		mockResponse  *api.SlurmV0042PostReservationResponse
		mockError     error
		expectedError bool
	}{
		{
			name: "successful create",
			reservation: &types.ReservationCreate{
				Name:      "new-reservation",
				Accounts:  []string{"account1", "account2"},
				NodeCount: 10,
				StartTime: 1640995200,
				EndTime:   1640998800,
				Users:     []string{"user1", "user2"},
				Partition: "normal",
			},
			mockResponse: &api.SlurmV0042PostReservationResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiResp{
					Meta: &api.V0042OpenapiMeta{},
				},
			},
			expectedError: false,
		},
		{
			name: "create with minimal fields",
			reservation: &types.ReservationCreate{
				Name:      "minimal-reservation",
				NodeCount: 5,
				StartTime: 1640995200,
				EndTime:   1640998800,
			},
			mockResponse: &api.SlurmV0042PostReservationResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200: &api.V0042OpenapiResp{
					Meta: &api.V0042OpenapiMeta{},
				},
			},
			expectedError: false,
		},
		{
			name: "API error",
			reservation: &types.ReservationCreate{
				Name:      "new-reservation",
				NodeCount: 10,
				StartTime: 1640995200,
				EndTime:   1640998800,
			},
			mockResponse: &api.SlurmV0042PostReservationResponse{
				HTTPResponse: &http.Response{StatusCode: 400},
			},
			expectedError: true,
		},
		{
			name: "network error",
			reservation: &types.ReservationCreate{
				Name:      "new-reservation",
				NodeCount: 10,
				StartTime: 1640995200,
				EndTime:   1640998800,
			},
			mockError:     fmt.Errorf("connection failed"),
			expectedError: true,
		},
		{
			name:          "nil reservation",
			reservation:   nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockReservationClient{}
			adapter := &ReservationAdapter{
				client:      mockClient,
				BaseManager: base.NewBaseManager("v0.0.42", "Reservation"),
			}

			if tt.reservation != nil {
				mockClient.On("SlurmV0042PostReservationWithResponse", mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			result, err := adapter.Create(context.Background(), tt.reservation)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.reservation != nil {
					assert.Equal(t, tt.reservation.Name, result.ReservationName)
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestReservationAdapter_Update(t *testing.T) {
	adapter := NewReservationAdapter(nil)
	
	updates := &types.ReservationUpdateRequest{
		NodeCount: ptrInt32(20),
	}
	
	err := adapter.Update(context.Background(), "test-reservation", updates)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not supported")
}

func TestReservationAdapter_Delete(t *testing.T) {
	tests := []struct {
		name            string
		reservationName string
		mockResponse    *api.SlurmV0042DeleteReservationResponse
		mockError       error
		expectedError   bool
	}{
		{
			name:            "successful delete",
			reservationName: "old-reservation",
			mockResponse: &api.SlurmV0042DeleteReservationResponse{
				HTTPResponse: &http.Response{StatusCode: 200},
			},
			expectedError: false,
		},
		{
			name:            "reservation not found",
			reservationName: "nonexistent",
			mockResponse: &api.SlurmV0042DeleteReservationResponse{
				HTTPResponse: &http.Response{StatusCode: 404},
			},
			expectedError: true,
		},
		{
			name:            "network error",
			reservationName: "old-reservation",
			mockError:       fmt.Errorf("timeout"),
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
				BaseManager: base.NewBaseManager("v0.0.42", "Reservation"),
			}

			if tt.reservationName != "" {
				mockClient.On("SlurmV0042DeleteReservationWithResponse", mock.Anything, tt.reservationName, mock.Anything).
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

func TestReservationAdapter_ContextValidation(t *testing.T) {
	adapter := NewReservationAdapter(&api.ClientWithResponses{})

	tests := []struct {
		name string
		ctx  context.Context
	}{
		{"nil context", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := adapter.List(tt.ctx, nil)
			assert.Error(t, err)

			_, err = adapter.Get(tt.ctx, "test")
			assert.Error(t, err)

			_, err = adapter.Create(tt.ctx, &types.ReservationCreate{Name: "test"})
			assert.Error(t, err)

			err = adapter.Update(tt.ctx, "test", &types.ReservationUpdateRequest{})
			assert.Error(t, err)

			err = adapter.Delete(tt.ctx, "test")
			assert.Error(t, err)
		})
	}
}

func TestReservationAdapter_ClientValidation(t *testing.T) {
	adapter := NewReservationAdapter(nil)
	ctx := context.Background()

	_, err := adapter.List(ctx, nil)
	assert.Error(t, err)

	_, err = adapter.Get(ctx, "test")
	assert.Error(t, err)

	_, err = adapter.Create(ctx, &types.ReservationCreate{Name: "test"})
	assert.Error(t, err)

	err = adapter.Update(ctx, "test", &types.ReservationUpdateRequest{})
	assert.Error(t, err)

	err = adapter.Delete(ctx, "test")
	assert.Error(t, err)
}