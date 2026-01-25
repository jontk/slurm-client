// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"testing"

	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
	"github.com/stretchr/testify/assert"
)

func TestBasicAdapterCreation(t *testing.T) {
	client := &api.ClientWithResponses{}

	// Test cluster adapter
	clusterAdapter := NewClusterAdapter(client)
	assert.NotNil(t, clusterAdapter)
	assert.Equal(t, "v0.0.42", clusterAdapter.GetVersion())

	// Test partition adapter
	partitionAdapter := NewPartitionAdapter(client)
	assert.NotNil(t, partitionAdapter)
	assert.Equal(t, "v0.0.42", partitionAdapter.GetVersion())

	// Test qos adapter
	qosAdapter := NewQoSAdapter(client)
	assert.NotNil(t, qosAdapter)
	assert.Equal(t, "v0.0.42", qosAdapter.GetVersion())

	// Test reservation adapter
	reservationAdapter := NewReservationAdapter(client)
	assert.NotNil(t, reservationAdapter)
	assert.Equal(t, "v0.0.42", reservationAdapter.GetVersion())

	// Test error adapter
	errorAdapter := NewErrorAdapter()
	assert.NotNil(t, errorAdapter)

	// Test job adapter
	jobAdapter := NewJobAdapter(client)
	assert.NotNil(t, jobAdapter)
	assert.Equal(t, "v0.0.42", jobAdapter.GetVersion())

	// Test node adapter
	nodeAdapter := NewNodeAdapter(client)
	assert.NotNil(t, nodeAdapter)
	assert.Equal(t, "v0.0.42", nodeAdapter.GetVersion())

	// Test account adapter
	accountAdapter := NewAccountAdapter(client)
	assert.NotNil(t, accountAdapter)
	assert.Equal(t, "v0.0.42", accountAdapter.GetVersion())

	// Test user adapter
	userAdapter := NewUserAdapter(client)
	assert.NotNil(t, userAdapter)
	assert.Equal(t, "v0.0.42", userAdapter.GetVersion())

	// Test standalone adapter
	standaloneAdapter := NewStandaloneAdapter(client)
	assert.NotNil(t, standaloneAdapter)

	// Test wckey adapter
	wckeyAdapter := NewWCKeyAdapter(client)
	assert.NotNil(t, wckeyAdapter)
	assert.Equal(t, "v0.0.42", wckeyAdapter.GetVersion())
}

func TestAdapterFactory(t *testing.T) {
	// Test main adapter factory
	client := &api.ClientWithResponses{}
	adapter := NewAdapter(client)
	assert.NotNil(t, adapter)
	assert.Equal(t, "v0.0.42", adapter.GetVersion())

	// Test all manager getters
	assert.NotNil(t, adapter.GetClusterManager())
	assert.NotNil(t, adapter.GetPartitionManager())
	assert.NotNil(t, adapter.GetQoSManager())
	assert.NotNil(t, adapter.GetReservationManager())
	assert.NotNil(t, adapter.GetJobManager())
	assert.NotNil(t, adapter.GetNodeManager())
	assert.NotNil(t, adapter.GetAccountManager())
	assert.NotNil(t, adapter.GetUserManager())
	assert.NotNil(t, adapter.GetStandaloneManager())
	assert.NotNil(t, adapter.GetWCKeyManager())
	assert.NotNil(t, adapter.GetAssociationManager())
}

func TestBasicConversions(t *testing.T) {
	// Test basic conversion functions that just create structures
	clusterAdapter := NewClusterAdapter(&api.ClientWithResponses{})
	_, err := clusterAdapter.convertAPIClusterToCommon(api.V0042ClusterRec{})
	assert.Nil(t, err)

	partitionAdapter := NewPartitionAdapter(&api.ClientWithResponses{})
	_, err = partitionAdapter.convertAPIPartitionToCommon(api.V0042PartitionInfo{})
	assert.Nil(t, err)

	qosAdapter := NewQoSAdapter(&api.ClientWithResponses{})
	_, err = qosAdapter.convertAPIQoSToCommon(api.V0042Qos{})
	assert.Nil(t, err)

	reservationAdapter := NewReservationAdapter(&api.ClientWithResponses{})
	_, err = reservationAdapter.convertAPIReservationToCommon(api.V0042ReservationInfo{})
	assert.Nil(t, err)

	jobAdapter := NewJobAdapter(&api.ClientWithResponses{})
	_, err = jobAdapter.convertAPIJobToCommon(api.V0042JobInfo{})
	assert.Nil(t, err)

	nodeAdapter := NewNodeAdapter(&api.ClientWithResponses{})
	_, err = nodeAdapter.convertAPINodeToCommon(api.V0042Node{})
	assert.Nil(t, err)

	accountAdapter := NewAccountAdapter(&api.ClientWithResponses{})
	_, err = accountAdapter.convertAPIAccountToCommon(api.V0042Account{})
	assert.Nil(t, err)

	userAdapter := NewUserAdapter(&api.ClientWithResponses{})
	userAdapter.convertAPIUserToCommon(api.V0042User{})
}

func TestConversionWithData(t *testing.T) {
	// Test conversions with some actual data
	clusterAdapter := NewClusterAdapter(&api.ClientWithResponses{})

	name := "test-cluster"
	clusterRec := api.V0042ClusterRec{
		Name: &name,
	}
	cluster, err := clusterAdapter.convertAPIClusterToCommon(clusterRec)
	assert.Nil(t, err)
	assert.Equal(t, "test-cluster", cluster.Name)

	// Test partition conversion with data
	partitionAdapter := NewPartitionAdapter(&api.ClientWithResponses{})
	partName := "test-partition"
	total := int32(10)
	partitionInfo := api.V0042PartitionInfo{
		Name: &partName,
		Nodes: &struct {
			AllowedAllocation *string `json:"allowed_allocation,omitempty"`
			Configured        *string `json:"configured,omitempty"`
			Total             *int32  `json:"total,omitempty"`
		}{
			Total: &total,
		},
	}
	partition, err := partitionAdapter.convertAPIPartitionToCommon(partitionInfo)
	assert.Nil(t, err)
	assert.Equal(t, "test-partition", partition.Name)
	assert.Equal(t, int32(10), partition.TotalNodes)

	// Test QoS conversion with data
	qosAdapter := NewQoSAdapter(&api.ClientWithResponses{})
	qosName := "test-qos"
	priority := int32(100)
	set := true
	qos := api.V0042Qos{
		Name: &qosName,
		Priority: &api.V0042Uint32NoValStruct{
			Set:    &set,
			Number: &priority,
		},
	}
	result, err := qosAdapter.convertAPIQoSToCommon(qos)
	assert.Nil(t, err)
	assert.Equal(t, "test-qos", result.Name)
	assert.Equal(t, 100, result.Priority)

	// Test user conversion with data
	userAdapter := NewUserAdapter(&api.ClientWithResponses{})
	userName := "test-user"
	defaultAccount := "test-account"
	user := api.V0042User{
		Name: userName,
		Default: &struct {
			Account *string `json:"account,omitempty"`
			Wckey   *string `json:"wckey,omitempty"`
		}{
			Account: &defaultAccount,
		},
	}
	userResult := userAdapter.convertAPIUserToCommon(user)
	assert.Equal(t, "test-user", userResult.Name)
	assert.Equal(t, "test-account", userResult.DefaultAccount)
}
