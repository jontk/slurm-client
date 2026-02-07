// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package factory

import (
	"context"
	"testing"
	"time"

	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/internal/adapters/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStandaloneManager is a mock implementation of common.StandaloneManager
type MockStandaloneManager struct {
	mock.Mock
}

func (m *MockStandaloneManager) GetLicenses(ctx context.Context) (*types.LicenseList, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.LicenseList), args.Error(1)
}

func (m *MockStandaloneManager) GetShares(ctx context.Context, opts *types.GetSharesOptions) (*types.SharesList, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.SharesList), args.Error(1)
}

func (m *MockStandaloneManager) GetConfig(ctx context.Context) (*types.Config, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Config), args.Error(1)
}

func (m *MockStandaloneManager) GetDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Diagnostics), args.Error(1)
}

func (m *MockStandaloneManager) GetDBDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Diagnostics), args.Error(1)
}

func (m *MockStandaloneManager) GetInstance(ctx context.Context, opts *types.GetInstanceOptions) (*types.Instance, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Instance), args.Error(1)
}

func (m *MockStandaloneManager) GetInstances(ctx context.Context, opts *types.GetInstancesOptions) (*types.InstanceList, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.InstanceList), args.Error(1)
}

func (m *MockStandaloneManager) GetTRES(ctx context.Context) (*types.TRESList, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.TRESList), args.Error(1)
}

func (m *MockStandaloneManager) CreateTRES(ctx context.Context, req *types.CreateTRESRequest) (*types.TRES, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.TRES), args.Error(1)
}

func (m *MockStandaloneManager) Reconfigure(ctx context.Context) (*types.ReconfigureResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ReconfigureResponse), args.Error(1)
}

func (m *MockStandaloneManager) PingDatabase(ctx context.Context) (*types.PingResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.PingResponse), args.Error(1)
}

// MockVersionAdapter is a mock implementation of common.VersionAdapter
type MockVersionAdapter struct {
	mock.Mock
	standaloneManager *MockStandaloneManager
}

func (m *MockVersionAdapter) GetVersion() string { return "v0.0.43" }

func (m *MockVersionAdapter) GetCapabilities() types.ClientCapabilities {
	return types.ClientCapabilities{
		Version:          "v0.0.43",
		SupportsJobs:     true,
		SupportsNodes:    true,
		SupportsDiagnostics: true,
		SupportsLicenses: true,
		SupportsShares:   true,
	}
}

func (m *MockVersionAdapter) GetStandaloneManager() common.StandaloneAdapter {
	if m.standaloneManager != nil {
		return m.standaloneManager
	}
	return nil
}

func (m *MockVersionAdapter) GetJobManager() common.JobAdapter                 { return nil }
func (m *MockVersionAdapter) GetNodeManager() common.NodeAdapter               { return nil }
func (m *MockVersionAdapter) GetPartitionManager() common.PartitionAdapter     { return nil }
func (m *MockVersionAdapter) GetAccountManager() common.AccountAdapter         { return nil }
func (m *MockVersionAdapter) GetUserManager() common.UserAdapter               { return nil }
func (m *MockVersionAdapter) GetQoSManager() common.QoSAdapter                 { return nil }
func (m *MockVersionAdapter) GetReservationManager() common.ReservationAdapter { return nil }
func (m *MockVersionAdapter) GetAssociationManager() common.AssociationAdapter { return nil }
func (m *MockVersionAdapter) GetWCKeyManager() common.WCKeyAdapter             { return nil }
func (m *MockVersionAdapter) GetClusterManager() common.ClusterAdapter         { return nil }
func (m *MockVersionAdapter) GetInfoManager() common.InfoAdapter               { return nil }

func TestAdapterClient_GetLicenses(t *testing.T) {
	ctx := context.Background()
	mockStandalone := new(MockStandaloneManager)
	mockAdapter := &MockVersionAdapter{
		standaloneManager: mockStandalone,
	}

	client := &AdapterClient{
		adapter: mockAdapter,
		version: "v0.0.43",
	}

	expectedLicenses := &types.LicenseList{
		Licenses: []types.License{
			{
				Name:     "matlab",
				Total:    100,
				Used:     25,
				Free:     75,
				Reserved: 10,
			},
		},
		Meta: map[string]interface{}{
			"count": 1,
		},
	}

	mockStandalone.On("GetLicenses", ctx).Return(expectedLicenses, nil)

	result, err := client.GetLicenses(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Licenses, 1)
	assert.Equal(t, "matlab", result.Licenses[0].Name)
	assert.Equal(t, 100, result.Licenses[0].Total)
	assert.Equal(t, 25, result.Licenses[0].Used)
	assert.Equal(t, 75, result.Licenses[0].Free)
	mockStandalone.AssertExpectations(t)
}

func TestAdapterClient_GetShares(t *testing.T) {
	ctx := context.Background()
	mockStandalone := new(MockStandaloneManager)
	mockAdapter := &MockVersionAdapter{
		standaloneManager: mockStandalone,
	}

	client := &AdapterClient{
		adapter: mockAdapter,
		version: "v0.0.43",
	}

	opts := &types.GetSharesOptions{
		Users:    []string{"user1"},
		Accounts: []string{"account1"},
	}

	expectedShares := &types.SharesList{
		Shares: []types.Share{
			{
				Account:          "account1",
				User:             "user1",
				FairshareLevel:   0.5,
				FairshareShares:  100,
				RawShares:        1000,
				NormalizedShares: 0.1,
				EffectiveUsage:   50.0,
			},
		},
	}

	mockStandalone.On("GetShares", ctx, mock.MatchedBy(func(o *types.GetSharesOptions) bool {
		return o != nil && len(o.Users) == 1 && o.Users[0] == "user1"
	})).Return(expectedShares, nil)

	result, err := client.GetShares(ctx, opts)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Shares, 1)
	assert.Equal(t, "account1", result.Shares[0].Account)
	assert.Equal(t, "user1", result.Shares[0].User)
	assert.Equal(t, 0.5, result.Shares[0].FairshareLevel)
	mockStandalone.AssertExpectations(t)
}

func TestAdapterClient_GetDiagnostics(t *testing.T) {
	ctx := context.Background()
	mockStandalone := new(MockStandaloneManager)
	mockAdapter := &MockVersionAdapter{
		standaloneManager: mockStandalone,
	}

	client := &AdapterClient{
		adapter: mockAdapter,
		version: "v0.0.43",
	}

	now := time.Now()
	expectedDiag := &types.Diagnostics{
		DataCollected:     now,
		ReqTime:           12345,
		ReqTimeStart:      12340,
		ServerThreadCount: 4,
		AgentCount:        10,
		JobsSubmitted:     100,
		JobsStarted:       90,
		JobsCompleted:     80,
		JobsFailed:        5,
		BFActive:          true,
	}

	mockStandalone.On("GetDiagnostics", ctx).Return(expectedDiag, nil)

	result, err := client.GetDiagnostics(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, now, result.DataCollected)
	assert.Equal(t, int64(12345), result.ReqTime)
	assert.Equal(t, 4, result.ServerThreadCount)
	assert.Equal(t, 100, result.JobsSubmitted)
	assert.True(t, result.BFActive)
	mockStandalone.AssertExpectations(t)
}

func TestAdapterClient_GetTRES(t *testing.T) {
	ctx := context.Background()
	mockStandalone := new(MockStandaloneManager)
	mockAdapter := &MockVersionAdapter{
		standaloneManager: mockStandalone,
	}

	client := &AdapterClient{
		adapter: mockAdapter,
		version: "v0.0.43",
	}

	expectedTRES := &types.TRESList{
		TRES: []types.TRES{
			{
				ID:    ptrInt32(1),
				Type:  "cpu",
				Name:  ptrString("cpu"),
				Count: ptrInt64(1000),
			},
			{
				ID:    ptrInt32(2),
				Type:  "mem",
				Name:  ptrString("mem"),
				Count: ptrInt64(1024000),
			},
		},
	}

	mockStandalone.On("GetTRES", ctx).Return(expectedTRES, nil)

	result, err := client.GetTRES(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.TRES, 2)
	assert.Equal(t, "cpu", result.TRES[0].Type)
	assert.Equal(t, "mem", result.TRES[1].Type)
	mockStandalone.AssertExpectations(t)
}
