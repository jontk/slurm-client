// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_45

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_45"
)

// NodeConverterGoverter defines the goverter interface for Node conversions.
// This tests goverter's ability to handle complex SLURM patterns:
// - time_novalnumber: Unix timestamp wrapped in NoValStruct
// - novalnumber_uint64: Optional uint64 wrapped in NoValStruct
// - state_enum_slice: Enum slices like []NodeState
// - Custom helpers for nested structs
// goverter:converter
// goverter:output:file node_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_45
// goverter:extend ConvertTimeNoVal
// goverter:extend ConvertUint64NoVal
// goverter:extend ConvertNodeStateSlice
// goverter:extend ConvertNextStateAfterReboot
// goverter:extend ConvertNodeEnergyGoverter
// goverter:extend ConvertResumeAfterGoverter
// goverter:extend ConvertCertFlagsGoverter
// goverter:extend ConvertCSVStringToSlice
type NodeConverterGoverter interface {
	// ConvertAPINodeToCommon converts API V0045Node to common Node type
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
	//
	// Fields removed from the v0.0.45 API (deprecated in v0.0.44); left unset:
	// goverter:ignore ExternalSensors Power TRESWeighted
	ConvertAPINodeToCommon(source api.V0045Node) *types.Node
}
