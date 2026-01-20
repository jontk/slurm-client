// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClusterManagerStub(t *testing.T) {
	stub := NewClusterManagerStub("v0.0.40")
	require.NotNil(t, stub)

	ctx := context.Background()

	// Test List
	clusters, err := stub.List(ctx, nil)
	require.Error(t, err)
	assert.Nil(t, clusters)
	assert.Contains(t, err.Error(), "not implemented")

	// Test Get
	cluster, err := stub.Get(ctx, "test-cluster")
	require.Error(t, err)
	assert.Nil(t, cluster)
	assert.Contains(t, err.Error(), "not implemented")

	// Test Create
	createResp, err := stub.Create(ctx, nil)
	require.Error(t, err)
	assert.Nil(t, createResp)
	assert.Contains(t, err.Error(), "not implemented")

	// Test Update
	err = stub.Update(ctx, "test-cluster", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")

	// Test Delete
	err = stub.Delete(ctx, "test-cluster")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}
