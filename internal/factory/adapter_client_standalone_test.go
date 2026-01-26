// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package factory

import (
	"context"
	"testing"
	"time"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/internal/adapters/common"
	"github.com/jontk/slurm-client/internal/common/types"
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

func (m *MockVersionAdapter) GetStandaloneManager() common.StandaloneAdapter {
	if m.standaloneManager != nil {
		return m.standaloneManager
	}
	return nil
}

func (m *MockVersionAdapter) GetJobManager() common.JobAdapter                   { return nil }
func (m *MockVersionAdapter) GetNodeManager() common.NodeAdapter                 { return nil }
func (m *MockVersionAdapter) GetPartitionManager() common.PartitionAdapter       { return nil }
func (m *MockVersionAdapter) GetAccountManager() common.AccountAdapter           { return nil }
func (m *MockVersionAdapter) GetUserManager() common.UserAdapter                 { return nil }
func (m *MockVersionAdapter) GetQoSManager() common.QoSAdapter                   { return nil }
func (m *MockVersionAdapter) GetReservationManager() common.ReservationAdapter   { return nil }
func (m *MockVersionAdapter) GetAssociationManager() common.AssociationAdapter   { return nil }
func (m *MockVersionAdapter) GetWCKeyManager() common.WCKeyAdapter               { return nil }
func (m *MockVersionAdapter) GetClusterManager() common.ClusterAdapter           { return nil }
func (m *MockVersionAdapter) GetInfoManager() common.InfoAdapter                 { return nil }

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
	assert.Equal(t, 75, result.Licenses[0].Available)
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

	opts := &interfaces.GetSharesOptions{
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
	assert.Equal(t, 0.5, result.Shares[0].FairShare)
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
	assert.True(t, result.BfActive)
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
				ID:    1,
				Type:  "cpu",
				Name:  "cpu",
				Count: 1000,
			},
			{
				ID:    2,
				Type:  "mem",
				Name:  "mem",
				Count: 512000,
			},
		},
	}

	mockStandalone.On("GetTRES", ctx).Return(expectedTRES, nil)

	result, err := client.GetTRES(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.TRES, 2)
	assert.Equal(t, uint64(1), result.TRES[0].ID)
	assert.Equal(t, "cpu", result.TRES[0].Type)
	assert.Equal(t, int64(1000), result.TRES[0].Count)
	mockStandalone.AssertExpectations(t)
}

func TestAdapterClient_NilStandaloneManager(t *testing.T) {
	ctx := context.Background()

	// Mock adapter that returns nil for standalone manager
	mockAdapter := &MockVersionAdapter{
		standaloneManager: nil,
	}

	client := &AdapterClient{
		adapter: mockAdapter,
		version: "v0.0.40",
	}

	// Test all methods return proper error when standalone manager is nil
	t.Run("GetLicenses", func(t *testing.T) {
		result, err := client.GetLicenses(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "standalone operations not supported")
	})

	t.Run("GetShares", func(t *testing.T) {
		result, err := client.GetShares(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "standalone operations not supported")
	})

	t.Run("GetConfig", func(t *testing.T) {
		result, err := client.GetConfig(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "standalone operations not supported")
	})

	t.Run("GetDiagnostics", func(t *testing.T) {
		result, err := client.GetDiagnostics(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "standalone operations not supported")
	})
}

func TestConvertLicenseListToInterface(t *testing.T) {
	input := &types.LicenseList{
		Licenses: []types.License{
			{
				Name:       "matlab",
				Total:      100,
				Used:       30,
				Free:       70,
				Reserved:   5,
				RemoteUsed: 10,
			},
		},
		Meta: map[string]interface{}{
			"server": "license-server",
		},
	}

	result := convertLicenseListToInterface(input)

	assert.NotNil(t, result)
	assert.Len(t, result.Licenses, 1)
	assert.Equal(t, "matlab", result.Licenses[0].Name)
	assert.Equal(t, 100, result.Licenses[0].Total)
	assert.Equal(t, 30, result.Licenses[0].Used)
	assert.Equal(t, 70, result.Licenses[0].Available)
	assert.Equal(t, 5, result.Licenses[0].Reserved)
	assert.True(t, result.Licenses[0].Remote) // RemoteUsed > 0
	assert.Equal(t, 30.0, result.Licenses[0].Percent)
}

func TestConvertSharesListToInterface(t *testing.T) {
	input := &types.SharesList{
		Shares: []types.Share{
			{
				Account:          "research",
				User:             "alice",
				Cluster:          "cluster1",
				FairshareShares:  100,
				RawShares:        1000,
				NormalizedShares: 0.1,
				RawUsage:         500,
				NormalizedUsage:  0.05,
				EffectiveUsage:   50.0,
				FairshareLevel:   0.5,
			},
		},
	}

	result := convertSharesListToInterface(input)

	assert.NotNil(t, result)
	assert.Len(t, result.Shares, 1)
	assert.Equal(t, "research", result.Shares[0].Account)
	assert.Equal(t, "alice", result.Shares[0].User)
	assert.Equal(t, 100, result.Shares[0].Shares)
	assert.Equal(t, 1000, result.Shares[0].RawShares)
	assert.Equal(t, 0.5, result.Shares[0].FairShare)
}

func TestConvertGetSharesOptionsToTypes(t *testing.T) {
	input := &interfaces.GetSharesOptions{
		Users:    []string{"user1", "user2"},
		Accounts: []string{"acc1"},
		Clusters: []string{"cluster1"},
	}

	result := convertGetSharesOptionsToTypes(input)

	assert.NotNil(t, result)
	assert.Equal(t, []string{"user1", "user2"}, result.Users)
	assert.Equal(t, []string{"acc1"}, result.Accounts)
	assert.Equal(t, []string{"cluster1"}, result.Clusters)
}

func TestCalculatePercentage(t *testing.T) {
	tests := []struct {
		name     string
		used     int
		total    int
		expected float64
	}{
		{"normal case", 25, 100, 25.0},
		{"zero total", 10, 0, 0.0},
		{"zero used", 0, 100, 0.0},
		{"full usage", 100, 100, 100.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculatePercentage(tt.used, tt.total)
			assert.Equal(t, tt.expected, result)
		})
	}
}
