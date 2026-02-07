// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package factory

import (
	"context"
	"testing"

	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/internal/adapters/common"
	"github.com/jontk/slurm-client/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock adapter that implements common.VersionAdapter
type testVersionAdapter struct {
	version            string
	jobAdapter         common.JobAdapter
	nodeAdapter        common.NodeAdapter
	partitionAdapter   common.PartitionAdapter
	reservationAdapter common.ReservationAdapter
	accountAdapter     common.AccountAdapter
	associationAdapter common.AssociationAdapter
	qosAdapter         common.QoSAdapter
	userAdapter        common.UserAdapter
	infoAdapter        common.InfoAdapter
}

func (t *testVersionAdapter) GetVersion() string {
	return t.version
}

func (t *testVersionAdapter) GetCapabilities() types.ClientCapabilities {
	return types.ClientCapabilities{
		Version:              t.version,
		SupportsJobs:         true,
		SupportsNodes:        true,
		SupportsPartitions:   true,
		SupportsReservations: true,
		SupportsAccounts:     true,
		SupportsUsers:        true,
		SupportsQoS:          true,
		SupportsClusters:     true,
		SupportsAssociations: true,
		SupportsWCKeys:       true,
	}
}

func (t *testVersionAdapter) GetJobManager() common.JobAdapter {
	return t.jobAdapter
}

func (t *testVersionAdapter) GetNodeManager() common.NodeAdapter {
	return t.nodeAdapter
}

func (t *testVersionAdapter) GetPartitionManager() common.PartitionAdapter {
	return t.partitionAdapter
}

func (t *testVersionAdapter) GetReservationManager() common.ReservationAdapter {
	return t.reservationAdapter
}

func (t *testVersionAdapter) GetAccountManager() common.AccountAdapter {
	return t.accountAdapter
}

func (t *testVersionAdapter) GetAssociationManager() common.AssociationAdapter {
	return t.associationAdapter
}

func (t *testVersionAdapter) GetQoSManager() common.QoSAdapter {
	return t.qosAdapter
}

func (t *testVersionAdapter) GetUserManager() common.UserAdapter {
	return t.userAdapter
}

func (t *testVersionAdapter) GetStandaloneManager() common.StandaloneAdapter {
	return &mockStandaloneAdapter{}
}

func (t *testVersionAdapter) GetWCKeyManager() common.WCKeyAdapter {
	return nil
}

func (t *testVersionAdapter) GetClusterManager() common.ClusterAdapter {
	return nil
}

func (t *testVersionAdapter) GetInfoManager() common.InfoAdapter {
	return t.infoAdapter
}

// Mock standalone adapter
type mockStandaloneAdapter struct{}

func (m *mockStandaloneAdapter) GetLicenses(ctx context.Context) (*types.LicenseList, error) {
	return &types.LicenseList{}, nil
}

func (m *mockStandaloneAdapter) GetShares(ctx context.Context, opts *types.GetSharesOptions) (*types.SharesList, error) {
	return &types.SharesList{}, nil
}

func (m *mockStandaloneAdapter) GetConfig(ctx context.Context) (*types.Config, error) {
	return &types.Config{}, nil
}

func (m *mockStandaloneAdapter) GetDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	return &types.Diagnostics{}, nil
}

func (m *mockStandaloneAdapter) GetDBDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	return &types.Diagnostics{}, nil
}

func (m *mockStandaloneAdapter) GetInstance(ctx context.Context, opts *types.GetInstanceOptions) (*types.Instance, error) {
	return &types.Instance{}, nil
}

func (m *mockStandaloneAdapter) GetInstances(ctx context.Context, opts *types.GetInstancesOptions) (*types.InstanceList, error) {
	return &types.InstanceList{}, nil
}

func (m *mockStandaloneAdapter) GetTRES(ctx context.Context) (*types.TRESList, error) {
	return &types.TRESList{}, nil
}

func (m *mockStandaloneAdapter) CreateTRES(ctx context.Context, req *types.CreateTRESRequest) (*types.TRES, error) {
	return &types.TRES{}, nil
}

func (m *mockStandaloneAdapter) Reconfigure(ctx context.Context) (*types.ReconfigureResponse, error) {
	return &types.ReconfigureResponse{}, nil
}

func (m *mockStandaloneAdapter) PingDatabase(ctx context.Context) (*types.PingResponse, error) {
	return &types.PingResponse{}, nil
}

// Mock reservation adapter
type mockReservationAdapter struct {
	listFunc   func(ctx context.Context, opts *types.ReservationListOptions) (*types.ReservationList, error)
	getFunc    func(ctx context.Context, name string) (*types.Reservation, error)
	createFunc func(ctx context.Context, res *types.ReservationCreate) (*types.ReservationCreateResponse, error)
	updateFunc func(ctx context.Context, name string, update *types.ReservationUpdate) error
	deleteFunc func(ctx context.Context, name string) error
}

func (m *mockReservationAdapter) List(ctx context.Context, opts *types.ReservationListOptions) (*types.ReservationList, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, opts)
	}
	return &types.ReservationList{Reservations: []types.Reservation{}}, nil
}

func (m *mockReservationAdapter) Get(ctx context.Context, name string) (*types.Reservation, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, name)
	}
	return &types.Reservation{Name: ptrString(name)}, nil
}

func (m *mockReservationAdapter) Create(ctx context.Context, res *types.ReservationCreate) (*types.ReservationCreateResponse, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, res)
	}
	var resName string
	if res.Name != nil {
		resName = *res.Name
	}
	return &types.ReservationCreateResponse{ReservationName: resName}, nil
}

func (m *mockReservationAdapter) Update(ctx context.Context, name string, update *types.ReservationUpdate) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, name, update)
	}
	return nil
}

func (m *mockReservationAdapter) Delete(ctx context.Context, name string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, name)
	}
	return nil
}

// Mock association adapter
type mockAssociationAdapter struct {
	listFunc   func(ctx context.Context, opts *types.AssociationListOptions) (*types.AssociationList, error)
	getFunc    func(ctx context.Context, id string) (*types.Association, error)
	createFunc func(ctx context.Context, assoc *types.AssociationCreate) (*types.AssociationCreateResponse, error)
	updateFunc func(ctx context.Context, id string, update *types.AssociationUpdate) error
	deleteFunc func(ctx context.Context, id string) error
}

func (m *mockAssociationAdapter) List(ctx context.Context, opts *types.AssociationListOptions) (*types.AssociationList, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, opts)
	}
	return &types.AssociationList{Associations: []types.Association{}}, nil
}

func (m *mockAssociationAdapter) Get(ctx context.Context, id string) (*types.Association, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, id)
	}
	return &types.Association{ID: ptrInt32(1)}, nil
}

func (m *mockAssociationAdapter) Create(ctx context.Context, assoc *types.AssociationCreate) (*types.AssociationCreateResponse, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, assoc)
	}
	return &types.AssociationCreateResponse{Status: "success", Message: "Created association test-123"}, nil
}

func (m *mockAssociationAdapter) Update(ctx context.Context, id string, update *types.AssociationUpdate) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, update)
	}
	return nil
}

func (m *mockAssociationAdapter) Delete(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func TestAdapterClient_ReservationOperations(t *testing.T) {
	ctx := helpers.TestContext(t)

	// Create mock reservation adapter
	mockReservation := &mockReservationAdapter{
		listFunc: func(ctx context.Context, opts *types.ReservationListOptions) (*types.ReservationList, error) {
			return &types.ReservationList{
				Reservations: []types.Reservation{
					{
						Name:      ptrString("test-res-1"),
						NodeCount: ptrInt32(5),
					},
				},
			}, nil
		},
		getFunc: func(ctx context.Context, name string) (*types.Reservation, error) {
			return &types.Reservation{
				Name:      ptrString(name),
				NodeCount: ptrInt32(5),
			}, nil
		},
		createFunc: func(ctx context.Context, res *types.ReservationCreate) (*types.ReservationCreateResponse, error) {
			var resName string
			if res.Name != nil {
				resName = *res.Name
			}
			return &types.ReservationCreateResponse{
				ReservationName: resName,
			}, nil
		},
	}

	// Create test version adapter
	testAdapter := &testVersionAdapter{
		version:            "v0.0.42",
		reservationAdapter: mockReservation,
	}

	// Create adapter client
	client := &AdapterClient{
		adapter: testAdapter,
		version: testAdapter.GetVersion(),
	}

	// Test Reservations() returns the manager
	resManager := client.Reservations()
	require.NotNil(t, resManager)

	// Test List operation through the interface
	list, err := resManager.List(ctx, nil)
	helpers.AssertNoError(t, err)
	assert.Len(t, list.Reservations, 1)
	assert.Equal(t, "test-res-1", *list.Reservations[0].Name)

	// Test Get operation
	res, err := resManager.Get(ctx, "test-res-1")
	helpers.AssertNoError(t, err)
	assert.Equal(t, "test-res-1", *res.Name)
	assert.Equal(t, int32(5), *res.NodeCount)
}

func TestAdapterClient_AssociationOperations(t *testing.T) {
	ctx := helpers.TestContext(t)

	// Create mock association adapter
	mockAssociation := &mockAssociationAdapter{
		listFunc: func(ctx context.Context, opts *types.AssociationListOptions) (*types.AssociationList, error) {
			return &types.AssociationList{
				Associations: []types.Association{
					{
						ID:      ptrInt32(1),
						User:    "user1",
						Account: ptrString("account1"),
						Cluster: ptrString("cluster1"),
					},
					{
						ID:      ptrInt32(2),
						User:    "user1",
						Account: ptrString("account2"),
						Cluster: ptrString("cluster1"),
					},
				},
			}, nil
		},
		createFunc: func(ctx context.Context, assoc *types.AssociationCreate) (*types.AssociationCreateResponse, error) {
			return &types.AssociationCreateResponse{
				Status:  "success",
				Message: "Created association new-assoc-123",
			}, nil
		},
	}

	// Create test version adapter
	testAdapter := &testVersionAdapter{
		version:            "v0.0.42",
		associationAdapter: mockAssociation,
	}

	// Create adapter client
	client := &AdapterClient{
		adapter: testAdapter,
	}

	// Test Associations() returns the manager
	assocManager := client.Associations()
	require.NotNil(t, assocManager)

	// Test List operation through the interface
	list, err := assocManager.List(ctx, nil)
	helpers.AssertNoError(t, err)
	assert.Len(t, list.Associations, 2)

	// Both associations should belong to user1
	for _, assoc := range list.Associations {
		assert.Equal(t, "user1", assoc.User)
	}
}

func TestAdapterClient_Version(t *testing.T) {
	// Test with nil adapter
	client := &AdapterClient{
		adapter: nil,
	}
	assert.Equal(t, "", client.Version())

	// Test with valid adapter
	testAdapter := &testVersionAdapter{
		version: "v0.0.43",
	}

	client = &AdapterClient{
		adapter: testAdapter,
		version: testAdapter.GetVersion(),
	}
	assert.Equal(t, "v0.0.43", client.Version())
}

func TestAdapterClient_TypeConversions(t *testing.T) {
	// Test reservation type conversion
	reservation := types.Reservation{
		Name:      ptrString("test-res"),
		NodeCount: ptrInt32(10),
		Users:     ptrString("user1,user2"),
	}

	// Verify pointer dereferencing
	assert.Equal(t, "test-res", *reservation.Name)
	assert.Equal(t, int32(10), *reservation.NodeCount)
	assert.Equal(t, "user1,user2", *reservation.Users)

	// Test association type conversion
	association := types.Association{
		ID:      ptrInt32(123),
		User:    "testuser",
		Account: ptrString("testaccount"),
		Cluster: ptrString("testcluster"),
	}

	// Verify pointer dereferencing
	assert.Equal(t, int32(123), *association.ID)
	assert.Equal(t, "testuser", association.User)
	assert.Equal(t, "testaccount", *association.Account)
	assert.Equal(t, "testcluster", *association.Cluster)
}
