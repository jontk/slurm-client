// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

import (
	"context"
	"testing"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestNodeManager_List_NotImplemented(t *testing.T) {
	// Test that List returns client error when client is not initialized
	nodeManager := &NodeManager{
		client: &WrapperClient{},
	}

	_, err := nodeManager.List(context.Background(), nil)

	// Should return client not initialized error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
	// The impl should now be created
	assert.NotNil(t, nodeManager.impl)
}

func TestNodeManager_Get_NotImplemented(t *testing.T) {
	// Test that Get returns client error when client is not initialized
	nodeManager := &NodeManager{
		client: &WrapperClient{},
	}

	_, err := nodeManager.Get(context.Background(), "node-001")

	// Should return client not initialized error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
	// The impl should now be created
	assert.NotNil(t, nodeManager.impl)
}

func TestNodeManager_Update_NotImplemented(t *testing.T) {
	// Test that Update returns client error when client is not initialized
	nodeManager := &NodeManager{
		client: &WrapperClient{},
	}

	err := nodeManager.Update(context.Background(), "node-001", &interfaces.NodeUpdate{State: stringPtr("DRAIN")})

	// Should return client not initialized error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not initialized")
	// The impl should now be created
	assert.NotNil(t, nodeManager.impl)
}

func TestNodeManager_Watch_Structure(t *testing.T) {
	// Test that Watch method properly delegates to implementation
	nodeManager := &NodeManager{
		client: &WrapperClient{},
	}

	_, err := nodeManager.Watch(context.Background(), &interfaces.WatchNodesOptions{})

	// We expect an error since there's no real API client
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")
	// The impl should now be created
	assert.NotNil(t, nodeManager.impl)
}

// Helper functions for pointer creation removed due to redeclaration
