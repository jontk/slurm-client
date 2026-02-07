// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_43

import (
	"time"

	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_43"
)

// enhancePartitionWithSkippedFields adds the skipped fields to a partition after base conversion
func (a *PartitionAdapter) enhancePartitionWithSkippedFields(result *types.Partition, apiObj api.V0043PartitionInfo) {
	if result == nil {
		return
	}
	// Accounts
	if apiObj.Accounts != nil {
		result.Accounts = &types.PartitionAccounts{
			Allowed: apiObj.Accounts.Allowed,
			Deny:    apiObj.Accounts.Deny,
		}
	}
	// CPUs
	if apiObj.Cpus != nil {
		result.CPUs = &types.PartitionCPUs{
			TaskBinding: apiObj.Cpus.TaskBinding,
			Total:       apiObj.Cpus.Total,
		}
	}
	// Defaults
	if apiObj.Defaults != nil {
		result.Defaults = &types.PartitionDefaults{
			Job:          apiObj.Defaults.Job,
			MemoryPerCPU: apiObj.Defaults.MemoryPerCpu,
		}
		if apiObj.Defaults.PartitionMemoryPerCpu != nil && apiObj.Defaults.PartitionMemoryPerCpu.Number != nil {
			v := uint64(*apiObj.Defaults.PartitionMemoryPerCpu.Number)
			result.Defaults.PartitionMemoryPerCPU = &v
		}
		if apiObj.Defaults.PartitionMemoryPerNode != nil && apiObj.Defaults.PartitionMemoryPerNode.Number != nil {
			v := uint64(*apiObj.Defaults.PartitionMemoryPerNode.Number)
			result.Defaults.PartitionMemoryPerNode = &v
		}
		if apiObj.Defaults.Time != nil && apiObj.Defaults.Time.Number != nil {
			v := uint32(*apiObj.Defaults.Time.Number)
			result.Defaults.Time = &v
		}
	}
	// Groups
	if apiObj.Groups != nil {
		result.Groups = &types.PartitionGroups{
			Allowed: apiObj.Groups.Allowed,
		}
	}
	// Maximums
	if apiObj.Maximums != nil {
		result.Maximums = &types.PartitionMaximums{
			MemoryPerCPU: apiObj.Maximums.MemoryPerCpu,
			Shares:       apiObj.Maximums.Shares,
		}
		if apiObj.Maximums.CpusPerNode != nil && apiObj.Maximums.CpusPerNode.Number != nil {
			v := uint32(*apiObj.Maximums.CpusPerNode.Number)
			result.Maximums.CPUsPerNode = &v
		}
		if apiObj.Maximums.CpusPerSocket != nil && apiObj.Maximums.CpusPerSocket.Number != nil {
			v := uint32(*apiObj.Maximums.CpusPerSocket.Number)
			result.Maximums.CPUsPerSocket = &v
		}
		if apiObj.Maximums.Nodes != nil && apiObj.Maximums.Nodes.Number != nil {
			v := uint32(*apiObj.Maximums.Nodes.Number)
			result.Maximums.Nodes = &v
		}
		if apiObj.Maximums.OverTimeLimit != nil && apiObj.Maximums.OverTimeLimit.Number != nil {
			v := uint16(*apiObj.Maximums.OverTimeLimit.Number)
			result.Maximums.OverTimeLimit = &v
		}
		if apiObj.Maximums.Time != nil && apiObj.Maximums.Time.Number != nil {
			v := uint32(*apiObj.Maximums.Time.Number)
			result.Maximums.Time = &v
		}
	}
	// Minimums
	if apiObj.Minimums != nil {
		result.Minimums = &types.PartitionMinimums{
			Nodes: apiObj.Minimums.Nodes,
		}
	}
	// Nodes
	if apiObj.Nodes != nil {
		result.Nodes = &types.PartitionNodes{
			AllowedAllocation: apiObj.Nodes.AllowedAllocation,
			Configured:        apiObj.Nodes.Configured,
			Total:             apiObj.Nodes.Total,
		}
	}
	// Partition (nested struct with state info)
	if apiObj.Partition != nil {
		result.Partition = &types.PartitionPartition{}
		if apiObj.Partition.State != nil {
			for _, s := range *apiObj.Partition.State {
				result.Partition.State = append(result.Partition.State, types.StateValue(s))
			}
		}
	}
	// Priority
	if apiObj.Priority != nil {
		result.Priority = &types.PartitionPriority{
			JobFactor: apiObj.Priority.JobFactor,
			Tier:      apiObj.Priority.Tier,
		}
	}
	// QoS
	if apiObj.Qos != nil {
		result.QoS = &types.PartitionQoS{
			Allowed:  apiObj.Qos.Allowed,
			Deny:     apiObj.Qos.Deny,
			Assigned: apiObj.Qos.Assigned,
		}
	}
	// SelectType
	if apiObj.SelectType != nil {
		for _, st := range *apiObj.SelectType {
			result.SelectType = append(result.SelectType, types.SelectTypeValue(st))
		}
	}
	// SuspendTime
	if apiObj.SuspendTime != nil && apiObj.SuspendTime.Number != nil && *apiObj.SuspendTime.Number > 0 {
		result.SuspendTime = time.Unix(int64(*apiObj.SuspendTime.Number), 0)
	}
	// Timeouts
	if apiObj.Timeouts != nil {
		result.Timeouts = &types.PartitionTimeouts{}
		if apiObj.Timeouts.Resume != nil && apiObj.Timeouts.Resume.Number != nil {
			v := uint16(*apiObj.Timeouts.Resume.Number)
			result.Timeouts.Resume = &v
		}
		if apiObj.Timeouts.Suspend != nil && apiObj.Timeouts.Suspend.Number != nil {
			v := uint16(*apiObj.Timeouts.Suspend.Number)
			result.Timeouts.Suspend = &v
		}
	}
	// TRES
	if apiObj.Tres != nil {
		result.TRES = &types.PartitionTRES{
			BillingWeights: apiObj.Tres.BillingWeights,
			Configured:     apiObj.Tres.Configured,
		}
	}
}
