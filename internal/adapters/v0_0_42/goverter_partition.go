// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_42

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_42"
)

// PartitionConverterGoverter defines the goverter interface for Partition conversions.
//
// goverter:converter
// goverter:output:file partition_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_42
// goverter:extend ConvertPartitionAccounts
// goverter:extend ConvertPartitionCPUs
// goverter:extend ConvertPartitionDefaults
// goverter:extend ConvertPartitionGroups
// goverter:extend ConvertPartitionMaximums
// goverter:extend ConvertPartitionMinimums
// goverter:extend ConvertPartitionNodes
// goverter:extend ConvertPartitionPartition
// goverter:extend ConvertPartitionPriority
// goverter:extend ConvertPartitionQoS
// goverter:extend ConvertPartitionSelectType
// goverter:extend ConvertPartitionSuspendTime
// goverter:extend ConvertPartitionTimeouts
// goverter:extend ConvertPartitionTRES
//
//go:generate goverter gen .
type PartitionConverterGoverter interface {
	// ConvertAPIPartitionToCommon converts API V0042PartitionInfo to common Partition type
	// goverter:map Accounts | ConvertPartitionAccounts
	// goverter:map Cpus CPUs | ConvertPartitionCPUs
	// goverter:map Defaults | ConvertPartitionDefaults
	// goverter:map Groups | ConvertPartitionGroups
	// goverter:map Maximums | ConvertPartitionMaximums
	// goverter:map Minimums | ConvertPartitionMinimums
	// goverter:map Nodes | ConvertPartitionNodes
	// goverter:map Partition | ConvertPartitionPartition
	// goverter:map Priority | ConvertPartitionPriority
	// goverter:map Qos QoS | ConvertPartitionQoS
	// goverter:map SelectType | ConvertPartitionSelectType
	// goverter:map SuspendTime | ConvertPartitionSuspendTime
	// goverter:map Timeouts | ConvertPartitionTimeouts
	// goverter:map Tres TRES | ConvertPartitionTRES
	//
	// Fields that don't exist in v0_0_42 (added in later versions):
	// goverter:ignore Topology
	ConvertAPIPartitionToCommon(source api.V0042PartitionInfo) *types.Partition
}
