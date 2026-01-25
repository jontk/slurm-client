// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"context"
	"testing"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestPartitionManager_List_Structure(t *testing.T) {
	// Test that the PartitionManager properly creates implementation
	partitionManager := &PartitionManager{
		client: &WrapperClient{},
	}

	// Test that impl is created lazily
	assert.Nil(t, partitionManager.impl)

	// After attempting to call List (even with nil client), impl should be created
	_, err := partitionManager.List(context.Background(), nil)

	// We expect an error since there's no real API client
	assert.Error(t, err)
	// The impl should now be created
	assert.NotNil(t, partitionManager.impl)
}

func TestPartitionManager_Get_Structure(t *testing.T) {
	// Test that the PartitionManager properly creates implementation for Get
	partitionManager := &PartitionManager{
		client: &WrapperClient{},
	}

	// Test that impl is created lazily
	assert.Nil(t, partitionManager.impl)

	// After attempting to call Get (even with nil client), impl should be created
	_, err := partitionManager.Get(context.Background(), "test-partition")

	// We expect an error since there's no real API client
	assert.Error(t, err)
	// The impl should now be created
	assert.NotNil(t, partitionManager.impl)
}

func TestConvertAPIPartitionToInterface(t *testing.T) {
	// Test the conversion function with minimal data
	name := "gpu-partition"
	totalNodes := int32(10)
	totalCPUs := int32(480)
	maxTime := int32(720)         // 12 hours
	defaultTime := int32(60)      // 1 hour
	maxMemory := int64(128000)    // 128GB in MB
	defaultMemory := int64(32000) // 32GB in MB
	priorityTier := int32(100)
	allowedAccounts := "group1,group2"
	allowedGroups := "users,admins"
	nodeList := "gpu-[001-010]"
	partitionState := []string{"UP"}

	apiPartition := V0040PartitionInfo{
		Name: &name,
		Nodes: &struct {
			AllowedAllocation *string `json:"allowed_allocation,omitempty"`
			Configured        *string `json:"configured,omitempty"`
			Total             *int32  `json:"total,omitempty"`
		}{
			Total:      &totalNodes,
			Configured: &nodeList,
		},
		Cpus: &struct {
			TaskBinding *int32 `json:"task_binding,omitempty"`
			Total       *int32 `json:"total,omitempty"`
		}{
			Total: &totalCPUs,
		},
		Maximums: &struct {
			CpusPerNode   *V0040Uint32NoVal `json:"cpus_per_node,omitempty"`
			CpusPerSocket *V0040Uint32NoVal `json:"cpus_per_socket,omitempty"`
			MemoryPerCpu  *int64            `json:"memory_per_cpu,omitempty"`
			Nodes         *V0040Uint32NoVal `json:"nodes,omitempty"`
			OverTimeLimit *V0040Uint16NoVal `json:"over_time_limit,omitempty"`
			Oversubscribe *struct {
				Flags *[]string `json:"flags,omitempty"`
				Jobs  *int32    `json:"jobs,omitempty"`
			} `json:"oversubscribe,omitempty"`
			PartitionMemoryPerCpu  *V0040Uint64NoVal `json:"partition_memory_per_cpu,omitempty"`
			PartitionMemoryPerNode *V0040Uint64NoVal `json:"partition_memory_per_node,omitempty"`
			Shares                 *int32            `json:"shares,omitempty"`
			Time                   *V0040Uint32NoVal `json:"time,omitempty"`
		}{
			MemoryPerCpu: &maxMemory,
			Time: &V0040Uint32NoVal{
				Number: &[]int64{int64(maxTime)}[0],
				Set:    &[]bool{true}[0],
			},
		},
		Defaults: &struct {
			Job                    *string           `json:"job,omitempty"`
			MemoryPerCpu           *int64            `json:"memory_per_cpu,omitempty"`
			PartitionMemoryPerCpu  *V0040Uint64NoVal `json:"partition_memory_per_cpu,omitempty"`
			PartitionMemoryPerNode *V0040Uint64NoVal `json:"partition_memory_per_node,omitempty"`
			Time                   *V0040Uint32NoVal `json:"time,omitempty"`
		}{
			MemoryPerCpu: &defaultMemory,
			Time: &V0040Uint32NoVal{
				Number: &[]int64{int64(defaultTime)}[0],
				Set:    &[]bool{true}[0],
			},
		},
		Accounts: &struct {
			Allowed *string `json:"allowed,omitempty"`
			Deny    *string `json:"deny,omitempty"`
		}{
			Allowed: &allowedAccounts,
		},
		Groups: &struct {
			Allowed *string `json:"allowed,omitempty"`
		}{
			Allowed: &allowedGroups,
		},
		Priority: &struct {
			JobFactor *int32 `json:"job_factor,omitempty"`
			Tier      *int32 `json:"tier,omitempty"`
		}{
			Tier: &priorityTier,
		},
		Partition: &struct {
			State *[]string `json:"state,omitempty"`
		}{
			State: &partitionState,
		},
	}

	interfacePartition := convertAPIPartitionToInterface(apiPartition)

	assert.NotNil(t, interfacePartition)
	assert.Equal(t, "gpu-partition", interfacePartition.Name)
	assert.Equal(t, "UP", interfacePartition.State)
	assert.Equal(t, 10, interfacePartition.TotalNodes)
	assert.Equal(t, 10, interfacePartition.AvailableNodes) // Simplified: all nodes available
	assert.Equal(t, 480, interfacePartition.TotalCPUs)
	assert.Equal(t, 480, interfacePartition.IdleCPUs) // Simplified: all CPUs idle
	assert.Equal(t, 720, interfacePartition.MaxTime)
	assert.Equal(t, 60, interfacePartition.DefaultTime)
	assert.Equal(t, 128000*1024*1024, interfacePartition.MaxMemory)    // Converted to bytes
	assert.Equal(t, 32000*1024*1024, interfacePartition.DefaultMemory) // Converted to bytes
	assert.Equal(t, []string{"group1", "group2"}, interfacePartition.AllowedUsers)
	assert.Equal(t, []string{"users", "admins"}, interfacePartition.AllowedGroups)
	assert.Equal(t, []string{}, interfacePartition.DeniedUsers)
	assert.Equal(t, []string{}, interfacePartition.DeniedGroups)
	assert.Equal(t, 100, interfacePartition.Priority)
	assert.Equal(t, []string{"gpu-[001-010]"}, interfacePartition.Nodes)
}

func TestConvertAPIPartitionToInterface_EmptyFields(t *testing.T) {
	// Test with empty/nil fields
	apiPartition := V0040PartitionInfo{}

	interfacePartition := convertAPIPartitionToInterface(apiPartition)

	assert.NotNil(t, interfacePartition)
	assert.Equal(t, "", interfacePartition.Name)
	assert.Equal(t, "", interfacePartition.State)
	assert.Equal(t, 0, interfacePartition.TotalNodes)
	assert.Equal(t, 0, interfacePartition.AvailableNodes)
	assert.Equal(t, 0, interfacePartition.TotalCPUs)
	assert.Equal(t, 0, interfacePartition.IdleCPUs)
	assert.Equal(t, 0, interfacePartition.MaxTime)
	assert.Equal(t, 0, interfacePartition.DefaultTime)
	assert.Equal(t, 0, interfacePartition.MaxMemory)
	assert.Equal(t, 0, interfacePartition.DefaultMemory)
	assert.Equal(t, []string{}, interfacePartition.AllowedUsers)
	assert.Equal(t, []string{}, interfacePartition.AllowedGroups)
	assert.Equal(t, []string{}, interfacePartition.DeniedUsers)
	assert.Equal(t, []string{}, interfacePartition.DeniedGroups)
	assert.Equal(t, 0, interfacePartition.Priority)
	assert.Equal(t, []string{}, interfacePartition.Nodes)
}

func TestFilterPartitions(t *testing.T) {
	partitions := []interfaces.Partition{
		{Name: "gpu", State: "UP", TotalNodes: 10, TotalCPUs: 480},
		{Name: "cpu", State: "DOWN", TotalNodes: 20, TotalCPUs: 960},
		{Name: "debug", State: "UP", TotalNodes: 2, TotalCPUs: 96},
	}

	// Test filter by state
	opts := &interfaces.ListPartitionsOptions{States: []string{"UP"}}
	filtered := filterPartitions(partitions, opts)
	assert.Len(t, filtered, 2)
	assert.Equal(t, "gpu", filtered[0].Name)
	assert.Equal(t, "debug", filtered[1].Name)

	// Test filter by multiple states
	opts = &interfaces.ListPartitionsOptions{States: []string{"UP", "DOWN"}}
	filtered = filterPartitions(partitions, opts)
	assert.Len(t, filtered, 3) // All partitions match

	// Test limit and offset
	opts = &interfaces.ListPartitionsOptions{Limit: 1, Offset: 1}
	filtered = filterPartitions(partitions, opts)
	assert.Len(t, filtered, 1)
	assert.Equal(t, "cpu", filtered[0].Name)

	// Test offset beyond available partitions
	opts = &interfaces.ListPartitionsOptions{Offset: 10}
	filtered = filterPartitions(partitions, opts)
	assert.Len(t, filtered, 0)
}

func TestFilterPartitions_EmptyFilter(t *testing.T) {
	partitions := []interfaces.Partition{
		{Name: "gpu", State: "UP"},
		{Name: "cpu", State: "DOWN"},
	}

	// Test with nil options (should return all partitions)
	filtered := filterPartitions(partitions, nil)
	assert.Len(t, filtered, 2)

	// Test with empty options (should return all partitions)
	opts := &interfaces.ListPartitionsOptions{}
	filtered = filterPartitions(partitions, opts)
	assert.Len(t, filtered, 2)
}

func TestNewPartitionManagerImpl(t *testing.T) {
	client := &WrapperClient{}
	impl := NewPartitionManagerImpl(client)

	assert.NotNil(t, impl)
	assert.Equal(t, client, impl.client)
}
