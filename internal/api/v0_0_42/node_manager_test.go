// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"testing"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestNodeManager_List_Structure(t *testing.T) {
	// Test that the NodeManager properly creates implementation
	nodeManager := &NodeManager{
		client: &WrapperClient{},
	}

	// Test that impl is created lazily
	assert.Nil(t, nodeManager.impl)

	// After attempting to call List (even with nil client), impl should be created
	_, err := nodeManager.List(context.Background(), nil)

	// We expect an error since there's no real API client
	assert.Error(t, err)
	// The impl should now be created
	assert.NotNil(t, nodeManager.impl)
}

func TestNodeManager_Get_Structure(t *testing.T) {
	// Test that the NodeManager properly creates implementation for Get
	nodeManager := &NodeManager{
		client: &WrapperClient{},
	}

	// Test that impl is created lazily
	assert.Nil(t, nodeManager.impl)

	// After attempting to call Get (even with nil client), impl should be created
	_, err := nodeManager.Get(context.Background(), "test-node")

	// We expect an error since there's no real API client
	assert.Error(t, err)
	// The impl should now be created
	assert.NotNil(t, nodeManager.impl)
}

func TestConvertAPINodeToInterface(t *testing.T) {
	// Test the conversion function with minimal data
	name := "node-001"
	cpus := int32(48)
	realMemory := int64(128000) // 128GB in MB
	partitions := V0042CsvString{"gpu", "cpu"}
	features := V0042CsvString{"intel", "ib"}
	state := V0042NodeStates{"IDLE"}

	apiNode := V0042Node{
		Name:       &name,
		Cpus:       &cpus,
		RealMemory: &realMemory,
		Partitions: &partitions,
		Features:   &features,
		State:      &state,
	}

	interfaceNode, err := convertAPINodeToInterface(apiNode)

	assert.NoError(t, err)
	assert.NotNil(t, interfaceNode)
	assert.Equal(t, "node-001", interfaceNode.Name)
	assert.Equal(t, 48, interfaceNode.CPUs)
	assert.Equal(t, 128000*1024*1024, interfaceNode.Memory) // Converted to bytes
	assert.Equal(t, []string{"gpu", "cpu"}, interfaceNode.Partitions)
	assert.Equal(t, []string{"intel", "ib"}, interfaceNode.Features)
	assert.Equal(t, "IDLE", interfaceNode.State)
	assert.NotNil(t, interfaceNode.Metadata)
}

func TestConvertAPINodeToInterface_EmptyFields(t *testing.T) {
	// Test with empty/nil fields
	apiNode := V0042Node{}

	interfaceNode, err := convertAPINodeToInterface(apiNode)

	assert.NoError(t, err)
	assert.NotNil(t, interfaceNode)
	assert.Equal(t, "", interfaceNode.Name)
	assert.Equal(t, 0, interfaceNode.CPUs)
	assert.Equal(t, 0, interfaceNode.Memory)
	assert.Equal(t, []string{}, interfaceNode.Partitions)
	assert.Equal(t, []string{}, interfaceNode.Features)
	assert.Equal(t, "", interfaceNode.State)
	assert.NotNil(t, interfaceNode.Metadata)
}

func TestFilterNodes(t *testing.T) {
	nodes := []interfaces.Node{
		{Name: "node-001", State: "IDLE", Partitions: []string{"gpu", "cpu"}, Features: []string{"intel", "ib"}},
		{Name: "node-002", State: "ALLOCATED", Partitions: []string{"cpu"}, Features: []string{"amd", "ethernet"}},
		{Name: "node-003", State: "DOWN", Partitions: []string{"gpu"}, Features: []string{"intel", "ib"}},
	}

	// Test filter by state
	opts := &interfaces.ListNodesOptions{States: []string{"IDLE"}}
	filtered := filterNodes(nodes, opts)
	assert.Len(t, filtered, 1)
	assert.Equal(t, "node-001", filtered[0].Name)

	// Test filter by partition
	opts = &interfaces.ListNodesOptions{Partition: "gpu"}
	filtered = filterNodes(nodes, opts)
	assert.Len(t, filtered, 2)
	assert.Equal(t, "node-001", filtered[0].Name)
	assert.Equal(t, "node-003", filtered[1].Name)

	// Test filter by features (all required features must be present)
	opts = &interfaces.ListNodesOptions{Features: []string{"intel", "ib"}}
	filtered = filterNodes(nodes, opts)
	assert.Len(t, filtered, 2)
	assert.Equal(t, "node-001", filtered[0].Name)
	assert.Equal(t, "node-003", filtered[1].Name)

	// Test filter by features (partial match should not return node)
	opts = &interfaces.ListNodesOptions{Features: []string{"intel", "nonexistent"}}
	filtered = filterNodes(nodes, opts)
	assert.Len(t, filtered, 0)

	// Test limit and offset
	opts = &interfaces.ListNodesOptions{Limit: 1, Offset: 1}
	filtered = filterNodes(nodes, opts)
	assert.Len(t, filtered, 1)
	assert.Equal(t, "node-002", filtered[0].Name)

	// Test offset beyond available nodes
	opts = &interfaces.ListNodesOptions{Offset: 10}
	filtered = filterNodes(nodes, opts)
	assert.Len(t, filtered, 0)

	// Test combined filters
	opts = &interfaces.ListNodesOptions{
		States:    []string{"IDLE", "ALLOCATED"},
		Partition: "cpu",
	}
	filtered = filterNodes(nodes, opts)
	assert.Len(t, filtered, 2)
	assert.Equal(t, "node-001", filtered[0].Name)
	assert.Equal(t, "node-002", filtered[1].Name)
}

func TestFilterNodes_EmptyFilter(t *testing.T) {
	nodes := []interfaces.Node{
		{Name: "node-001", State: "IDLE"},
		{Name: "node-002", State: "ALLOCATED"},
	}

	// Test with nil options (should return all nodes)
	filtered := filterNodes(nodes, nil)
	assert.Len(t, filtered, 2)

	// Test with empty options (should return all nodes)
	opts := &interfaces.ListNodesOptions{}
	filtered = filterNodes(nodes, opts)
	assert.Len(t, filtered, 2)
}

func TestNewNodeManagerImpl(t *testing.T) {
	client := &WrapperClient{}
	impl := NewNodeManagerImpl(client)

	assert.NotNil(t, impl)
	assert.Equal(t, client, impl.client)
}
