// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"testing"

	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/stretchr/testify/assert"
)

// Test error validation paths to increase coverage without complex mocking
func TestValidationErrorPaths(t *testing.T) {
	client := &api.ClientWithResponses{}

	// Test cluster adapter validation with nil context
	clusterAdapter := NewClusterAdapter(client)
	_, err := clusterAdapter.List(nil, nil) //lint:ignore SA1012 Intentionally testing nil context handling
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context")

	// Test partition adapter validation with nil context
	partitionAdapter := NewPartitionAdapter(client)
	_, err = partitionAdapter.List(nil, nil) //lint:ignore SA1012 Intentionally testing nil context handling
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context")

	// Test QoS adapter validation with nil context
	qosAdapter := NewQoSAdapter(client)
	_, err = qosAdapter.List(nil, nil) //lint:ignore SA1012 Intentionally testing nil context handling
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context")

	// Test reservation adapter validation with nil context
	reservationAdapter := NewReservationAdapter(client)
	_, err = reservationAdapter.List(nil, nil) //lint:ignore SA1012 Intentionally testing nil context handling
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context")

	// Test job adapter validation with nil context
	jobAdapter := NewJobAdapter(client)
	_, err = jobAdapter.List(nil, nil) //lint:ignore SA1012 Intentionally testing nil context handling
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context")

	// Test node adapter validation with nil context
	nodeAdapter := NewNodeAdapter(client)
	_, err = nodeAdapter.List(nil, nil) //lint:ignore SA1012 Intentionally testing nil context handling
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context")

	// Test account adapter validation with nil context
	accountAdapter := NewAccountAdapter(client)
	_, err = accountAdapter.List(nil, nil) //lint:ignore SA1012 Intentionally testing nil context handling
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context")

	// Test user adapter validation with nil context
	userAdapter := NewUserAdapter(client)
	_, err = userAdapter.List(nil, nil) //lint:ignore SA1012 Intentionally testing nil context handling
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context")

	// Test association adapter validation with nil context
	associationAdapter := NewAssociationAdapter(client)
	_, err = associationAdapter.List(nil, nil) //lint:ignore SA1012 Intentionally testing nil context handling
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context")

	// Test WCKey adapter validation with nil context
	wckeyAdapter := NewWCKeyAdapter(client)
	_, err = wckeyAdapter.List(nil, nil) //lint:ignore SA1012 Intentionally testing nil context handling
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context")
}

func TestGetOperations(t *testing.T) {
	// Test with nil client to exercise client validation paths
	clusterAdapter := NewClusterAdapter(nil)
	ctx := context.Background()

	_, err := clusterAdapter.List(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	partitionAdapter := NewPartitionAdapter(nil)
	_, err = partitionAdapter.List(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	qosAdapter := NewQoSAdapter(nil)
	_, err = qosAdapter.List(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	reservationAdapter := NewReservationAdapter(nil)
	_, err = reservationAdapter.List(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	jobAdapter := NewJobAdapter(nil)
	_, err = jobAdapter.List(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	nodeAdapter := NewNodeAdapter(nil)
	_, err = nodeAdapter.List(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	accountAdapter := NewAccountAdapter(nil)
	_, err = accountAdapter.List(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	userAdapter := NewUserAdapter(nil)
	_, err = userAdapter.List(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	wckeyAdapter := NewWCKeyAdapter(nil)
	_, err = wckeyAdapter.List(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")
}

func TestCreateOperations(t *testing.T) {
	ctx := context.Background()

	// Test create operations with nil client to exercise validation
	clusterAdapter := NewClusterAdapter(nil)
	_, err := clusterAdapter.Create(ctx, &types.ClusterCreate{Name: "test"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	partitionAdapter := NewPartitionAdapter(nil)
	_, err = partitionAdapter.Create(ctx, &types.PartitionCreate{Name: "test"})
	assert.Error(t, err)
	// Partition creation is not supported in v0.0.42, so check for appropriate error
	assert.Contains(t, err.Error(), "not supported")

	client := &api.ClientWithResponses{}
	reservationAdapter := NewReservationAdapter(client)
	_, err = reservationAdapter.Create(ctx, &types.ReservationCreate{Name: "test"})
	assert.Error(t, err)
	// Reservation creation is not supported in v0.0.42
	assert.Contains(t, err.Error(), "not supported")

	// QoS create should exercise validation
	qosAdapter := NewQoSAdapter(nil)
	_, err = qosAdapter.Create(ctx, &types.QoSCreate{Name: "test"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	// User create should exercise validation
	userAdapter := NewUserAdapter(nil)
	_, err = userAdapter.Create(ctx, &types.UserCreate{Name: "test"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")
}

func TestUpdateOperations(t *testing.T) {
	ctx := context.Background()

	// Test update operations with nil client to exercise validation
	partitionAdapter := NewPartitionAdapter(nil)
	err := partitionAdapter.Update(ctx, "test", &types.PartitionUpdate{})
	assert.Error(t, err)
	// Partition update is not supported in v0.0.42
	assert.Contains(t, err.Error(), "not supported")

	qosAdapter := NewQoSAdapter(nil)
	err = qosAdapter.Update(ctx, "test", &types.QoSUpdateRequest{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	reservationAdapter := NewReservationAdapter(nil)
	err = reservationAdapter.Update(ctx, "test", &types.ReservationUpdateRequest{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	userAdapter := NewUserAdapter(nil)
	err = userAdapter.Update(ctx, "test", &types.UserUpdateRequest{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")
}

func TestDeleteOperations(t *testing.T) {
	ctx := context.Background()

	// Test delete operations with nil client to exercise validation
	clusterAdapter := NewClusterAdapter(nil)
	err := clusterAdapter.Delete(ctx, "test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	reservationAdapter := NewReservationAdapter(nil)
	err = reservationAdapter.Delete(ctx, "test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	qosAdapter := NewQoSAdapter(nil)
	err = qosAdapter.Delete(ctx, "test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	userAdapter := NewUserAdapter(nil)
	err = userAdapter.Delete(ctx, "test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")
}

func TestSpecialOperations(t *testing.T) {
	ctx := context.Background()

	// Test job submission with nil client to exercise validation
	jobAdapter := NewJobAdapter(nil)
	_, err := jobAdapter.Submit(ctx, &types.JobCreate{
		Name:   "test-job",
		Script: "#!/bin/bash\necho 'test'",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	// Test user association creation with nil client
	userAdapter := NewUserAdapter(nil)
	_, err = userAdapter.CreateAssociation(ctx, &types.UserAssociationRequest{
		Users:   []string{"test-user"},
		Account: "test-account",
		Cluster: "test-cluster",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	// Test standalone operations with nil client
	standaloneAdapter := NewStandaloneAdapter(nil)
	_, err = standaloneAdapter.GetLicenses(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client")

	// Test error adapter
	errorAdapter := NewErrorAdapter()
	assert.NotNil(t, errorAdapter)
	// Error adapter methods need different testing approach
}
