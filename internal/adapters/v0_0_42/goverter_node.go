// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_42

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_42"
)

// NodeConverterGoverter defines the goverter interface for Node conversions.
// This tests goverter's ability to handle complex SLURM patterns:
// - time_novalnumber: Unix timestamp wrapped in NoValStruct
// - novalnumber_uint64: Optional uint64 wrapped in NoValStruct
// - state_enum_slice: Enum slices like []NodeState
// - Custom helpers for nested structs
// goverter:converter
// goverter:output:file node_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_42
// goverter:extend ConvertTimeNoVal
// goverter:extend ConvertUint64NoVal
// goverter:extend ConvertNodeStateSlice
// goverter:extend ConvertNextStateAfterReboot
// goverter:extend ConvertNodeEnergyGoverter
// goverter:extend ConvertResumeAfterGoverter
// goverter:extend ConvertCertFlagsGoverter
// goverter:extend ConvertExternalSensors
// goverter:extend ConvertPower
// goverter:extend ConvertCSVStringToSlice
type NodeConverterGoverter interface {
	// ConvertAPINodeToCommon converts API V0042Node to common Node type
	//
	// Time fields (use ConvertTimeNoVal):
	// goverter:map BootTime | ConvertTimeNoVal
	// goverter:map LastBusy | ConvertTimeNoVal
	// goverter:map ReasonChangedAt | ConvertTimeNoVal
	// goverter:map SlurmdStartTime | ConvertTimeNoVal
	// goverter:map TlsCertLastRenewal TLSCertLastRenewal | ConvertTimeNoVal
	//
	// NoValNumber fields:
	// goverter:map FreeMem | ConvertUint64NoVal
	//
	// State enum slices:
	// goverter:map State | ConvertNodeStateSlice
	// goverter:map NextStateAfterReboot | ConvertNextStateAfterReboot
	//
	// Custom helpers:
	// goverter:map Energy | ConvertNodeEnergyGoverter
	// goverter:map ResumeAfter | ConvertResumeAfterGoverter
	// goverter:map CertFlags | ConvertCertFlagsGoverter
	// goverter:map ExternalSensors | ConvertExternalSensors
	// goverter:map Power | ConvertPower
	//
	// CSV string to slice:
	// goverter:map ActiveFeatures | ConvertCSVStringToSlice
	// goverter:map Features | ConvertCSVStringToSlice
	// goverter:map Partitions | ConvertCSVStringToSlice
	//
	// Field name mappings (Source -> Target):
	// goverter:map AllocCpus AllocCPUs
	// goverter:map AllocIdleCpus AllocIdleCPUs
	// goverter:map CpuBinding CPUBinding
	// goverter:map CpuLoad CPULoad
	// goverter:map Cpus CPUs
	// goverter:map EffectiveCpus EffectiveCPUs
	// goverter:map GpuSpec GPUSpec
	// goverter:map Gres GRES
	// goverter:map GresDrained GRESDrained
	// goverter:map GresUsed GRESUsed
	// goverter:map InstanceId InstanceID
	// goverter:map McsLabel MCSLabel
	// goverter:map ResCoresPerGpu ResCoresPerGPU
	// goverter:map SpecializedCpus SpecializedCPUs
	// goverter:map Tres TRES
	// goverter:map TresUsed TRESUsed
	// goverter:map TresWeighted TRESWeighted
	//
	// Fields that don't exist in v0_0_42:
	// goverter:ignore CertFlags
	// goverter:ignore TLSCertLastRenewal
	// goverter:ignore Topology
	ConvertAPINodeToCommon(source api.V0042Node) *types.Node
}
