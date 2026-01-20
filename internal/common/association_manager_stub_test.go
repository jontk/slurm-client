// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssociationManagerStub(t *testing.T) {
	stub := NewAssociationManagerStub("v0.0.40")
	require.NotNil(t, stub)

	ctx := context.Background()

	// Test List
	associations, err := stub.List(ctx, nil)
	require.Error(t, err)
	assert.Nil(t, associations)
	assert.Contains(t, err.Error(), "not implemented")

	// Test Get
	association, err := stub.Get(ctx, nil)
	require.Error(t, err)
	assert.Nil(t, association)
	assert.Contains(t, err.Error(), "not implemented")

	// Test Create
	createResp, err := stub.Create(ctx, nil)
	require.Error(t, err)
	assert.Nil(t, createResp)
	assert.Contains(t, err.Error(), "not implemented")

	// Test Update
	err = stub.Update(ctx, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")

	// Test Delete
	err = stub.Delete(ctx, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")

	// Test BulkDelete
	bulkResp, err := stub.BulkDelete(ctx, nil)
	require.Error(t, err)
	assert.Nil(t, bulkResp)
	assert.Contains(t, err.Error(), "not implemented")

	// Test GetUserAssociations
	userAssocs, err := stub.GetUserAssociations(ctx, "testuser")
	require.Error(t, err)
	assert.Nil(t, userAssocs)
	assert.Contains(t, err.Error(), "not implemented")

	// Test GetAccountAssociations
	accountAssocs, err := stub.GetAccountAssociations(ctx, "testaccount")
	require.Error(t, err)
	assert.Nil(t, accountAssocs)
	assert.Contains(t, err.Error(), "not implemented")

	// Test ValidateAssociation
	valid, err := stub.ValidateAssociation(ctx, "testuser", "testaccount", "testcluster")
	require.Error(t, err)
	assert.False(t, valid)
	assert.Contains(t, err.Error(), "not implemented")
}
