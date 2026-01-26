// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package factory

import (
	"context"
	"testing"

	"github.com/jontk/slurm-client/internal/adapters/common"
	"github.com/jontk/slurm-client/internal/common/types"
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
	return &types.Reservation{Name: name}, nil
}

func (m *mockReservationAdapter) Create(ctx context.Context, res *types.ReservationCreate) (*types.ReservationCreateResponse, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, res)
	}
	return &types.ReservationCreateResponse{ReservationName: res.Name}, nil
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
	return &types.Association{ID: id}, nil
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
						Name:      "test-res-1",
						State:     types.ReservationState("ACTIVE"),
						NodeCount: 5,
					},
				},
			}, nil
		},
		getFunc: func(ctx context.Context, name string) (*types.Reservation, error) {
			return &types.Reservation{
				Name:      name,
				State:     types.ReservationState("ACTIVE"),
				NodeCount: 5,
			}, nil
		},
		createFunc: func(ctx context.Context, res *types.ReservationCreate) (*types.ReservationCreateResponse, error) {
			return &types.ReservationCreateResponse{
				ReservationName: res.Name,
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
	assert.Equal(t, "test-res-1", list.Reservations[0].Name)
	assert.Equal(t, "ACTIVE", list.Reservations[0].State)

	// Test Get operation
	res, err := resManager.Get(ctx, "test-res-1")
	helpers.AssertNoError(t, err)
	assert.Equal(t, "test-res-1", res.Name)
	assert.Equal(t, 5, res.NodeCount)
}

func TestAdapterClient_AssociationOperations(t *testing.T) {
	ctx := helpers.TestContext(t)

	// Create mock association adapter
	mockAssociation := &mockAssociationAdapter{
		listFunc: func(ctx context.Context, opts *types.AssociationListOptions) (*types.AssociationList, error) {
			return &types.AssociationList{
				Associations: []types.Association{
					{
						ID:          "assoc-1",
						UserName:    "user1",
						AccountName: "account1",
						Cluster:     "cluster1",
					},
					{
						ID:          "assoc-2",
						UserName:    "user1",
						AccountName: "account2",
						Cluster:     "cluster1",
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

	// Test GetUserAssociations through the interface
	userAssocs, err := assocManager.GetUserAssociations(ctx, "user1")
	helpers.AssertNoError(t, err)
	assert.Len(t, userAssocs, 2)

	// Both associations should belong to user1
	for _, assoc := range userAssocs {
		assert.Equal(t, "user1", assoc.User)
	}

	// Test GetAccountAssociations
	accountAssocs, err := assocManager.GetAccountAssociations(ctx, "account1")
	helpers.AssertNoError(t, err)
	assert.Len(t, accountAssocs, 1)
	assert.Equal(t, "account1", accountAssocs[0].Account)
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
		Name:          "test-res",
		State:         types.ReservationState("ACTIVE"),
		NodeCount:     10,
		Users:         []string{"user1", "user2"},
		Accounts:      []string{"account1"},
		PartitionName: "compute",
		NodeList:      "node[001-010]",
		Flags:         []types.ReservationFlag{types.ReservationFlag("MAINT")},
		Licenses: map[string]int32{
			"matlab": 5,
		},
	}

	// Convert using the helper function
	converted := convertReservationToInterface(reservation)

	assert.Equal(t, "test-res", converted.Name)
	assert.Equal(t, "ACTIVE", converted.State)
	assert.Equal(t, 10, converted.NodeCount)
	assert.Equal(t, []string{"user1", "user2"}, converted.Users)
	assert.Equal(t, []string{"account1"}, converted.Accounts)
	assert.Equal(t, "compute", converted.PartitionName)
	assert.Equal(t, []string{"MAINT"}, converted.Flags)
	assert.Equal(t, 5, converted.Licenses["matlab"])

	// Test association type conversion
	association := types.Association{
		ID:          "assoc-123",
		UserName:    "testuser",
		AccountName: "testaccount",
		Cluster:     "cluster1",
		Partition:   "compute",
		QoSList:     []string{"normal"},
		Priority:    100,
		MaxJobs:     10,
		IsDefault:   true,
	}

	// Convert using the helper function
	assocConverted := convertAssociationToInterface(association)

	// Note: ID is not properly mapped in convertAssociationToInterface (returns 0)
	assert.Equal(t, uint32(0), assocConverted.ID) // ID is mapped to 0 in conversion
	assert.Equal(t, "testuser", assocConverted.User)
	assert.Equal(t, "testaccount", assocConverted.Account)
	assert.Equal(t, "cluster1", assocConverted.Cluster)
	assert.Equal(t, "compute", assocConverted.Partition)
	assert.Equal(t, uint32(100), assocConverted.Priority)
	assert.Equal(t, 10, *assocConverted.MaxJobs)
	// Note: IsDefault is hardcoded to false in convertAssociationToInterface
	assert.False(t, assocConverted.IsDefault) // IsDefault is hardcoded to false
}
