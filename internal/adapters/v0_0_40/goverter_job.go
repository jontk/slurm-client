// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_40

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_40"
)

// JobConverterGoverter defines the goverter interface for Job conversions.
// This is the most complex entity converter, handling many field types:
// - time_novalnumber: Unix timestamp wrapped in NoValStruct
// - novalnumber_uint32/uint16/uint64: Optional integers wrapped in NoValStruct
// - novalnumber_float64: Optional float wrapped in NoValStruct
// - state_enum_slice: Enum slices like []JobState
// goverter:converter
// goverter:output:file job_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_40
// goverter:extend ConvertTimeNoVal
// goverter:extend ConvertUint32NoVal
// goverter:extend ConvertUint16NoVal
// goverter:extend ConvertUint64NoVal
// goverter:extend ConvertFloat64NoVal
// goverter:extend ConvertJobStateSlice
// goverter:extend ConvertCSVStringToSlice
// goverter:extend ConvertJobFlags
// goverter:extend ConvertJobMailType
// goverter:extend ConvertJobProfile
// goverter:extend ConvertJobShared
// goverter:extend ConvertExitCode
// goverter:extend ConvertJobPower
// goverter:extend ConvertJobGRESDetail
// goverter:extend ConvertJobResources
type JobConverterGoverter interface {
	// ConvertAPIJobToCommon converts API V0040JobInfo to common Job type
	//
	// Time fields (use ConvertTimeNoVal):
	// goverter:map AccrueTime | ConvertTimeNoVal
	// goverter:map Deadline | ConvertTimeNoVal
	// goverter:map EligibleTime | ConvertTimeNoVal
	// goverter:map EndTime | ConvertTimeNoVal
	// goverter:map LastSchedEvaluation | ConvertTimeNoVal
	// goverter:map PreemptTime | ConvertTimeNoVal
	// goverter:map PreemptableTime | ConvertTimeNoVal
	// goverter:map ResizeTime | ConvertTimeNoVal
	// goverter:map StartTime | ConvertTimeNoVal
	// goverter:map SubmitTime | ConvertTimeNoVal
	// goverter:map SuspendTime | ConvertTimeNoVal
	//
	// NoValNumber Uint32 fields (use ConvertUint32NoVal):
	// goverter:map ArrayJobId ArrayJobID | ConvertUint32NoVal
	// goverter:map ArrayMaxTasks | ConvertUint32NoVal
	// goverter:map ArrayTaskId ArrayTaskID | ConvertUint32NoVal
	// goverter:map CpuFrequencyGovernor CPUFrequencyGovernor | ConvertUint32NoVal
	// goverter:map CpuFrequencyMaximum CPUFrequencyMaximum | ConvertUint32NoVal
	// goverter:map CpuFrequencyMinimum CPUFrequencyMinimum | ConvertUint32NoVal
	// goverter:map Cpus CPUs | ConvertUint32NoVal
	// goverter:map DelayBoot | ConvertUint32NoVal
	// goverter:map HetJobId HetJobID | ConvertUint32NoVal
	// goverter:map HetJobOffset | ConvertUint32NoVal
	// goverter:map MaxCpus MaxCPUs | ConvertUint32NoVal
	// goverter:map MaxNodes | ConvertUint32NoVal
	// goverter:map MinimumTmpDiskPerNode | ConvertUint32NoVal
	// goverter:map NodeCount | ConvertUint32NoVal
	// goverter:map Priority | ConvertUint32NoVal
	// goverter:map Tasks | ConvertUint32NoVal
	// goverter:map TimeLimit | ConvertUint32NoVal
	// goverter:map TimeMinimum | ConvertUint32NoVal
	//
	// NoValNumber Uint16 fields (use ConvertUint16NoVal):
	// goverter:map CoresPerSocket | ConvertUint16NoVal
	// goverter:map CpusPerTask CPUsPerTask | ConvertUint16NoVal
	// goverter:map MinimumCpusPerNode MinimumCPUsPerNode | ConvertUint16NoVal
	// goverter:map SocketsPerNode | ConvertUint16NoVal
	// goverter:map TasksPerBoard | ConvertUint16NoVal
	// goverter:map TasksPerCore | ConvertUint16NoVal
	// goverter:map TasksPerNode | ConvertUint16NoVal
	// goverter:map TasksPerSocket | ConvertUint16NoVal
	// goverter:map TasksPerTres TasksPerTRES | ConvertUint16NoVal
	// goverter:map ThreadsPerCore | ConvertUint16NoVal
	//
	// NoValNumber Uint64 fields (use ConvertUint64NoVal):
	// goverter:map MemoryPerCpu MemoryPerCPU | ConvertUint64NoVal
	// goverter:map MemoryPerNode | ConvertUint64NoVal
	// goverter:map PreSusTime | ConvertUint64NoVal
	//
	// NoValNumber Float64 fields (use ConvertFloat64NoVal):
	// goverter:map BillableTres BillableTRES | ConvertFloat64NoVal
	//
	// Enum slice fields:
	// goverter:map JobState | ConvertJobStateSlice
	//
	// CSV string to slice:
	// goverter:map JobSizeStr | ConvertCSVStringToSlice
	//
	// Field name mappings (Source -> Target):
	// goverter:map AssociationId AssociationID
	// goverter:map ContainerId ContainerID
	// goverter:map CpusPerTres CPUsPerTRES
	// goverter:map GresDetail GRESDetail
	// goverter:map GroupId GroupID
	// goverter:map HetJobIdSet HetJobIDSet
	// goverter:map JobId JobID
	// goverter:map McsLabel MCSLabel
	// goverter:map MemoryPerTres MemoryPerTRES
	// goverter:map Qos QoS
	// goverter:map TresAllocStr TRESAllocStr
	// goverter:map TresBind TRESBind
	// goverter:map TresFreq TRESFreq
	// goverter:map TresPerJob TRESPerJob
	// goverter:map TresPerNode TRESPerNode
	// goverter:map TresPerSocket TRESPerSocket
	// goverter:map TresPerTask TRESPerTask
	// goverter:map TresReqStr TRESReqStr
	// goverter:map UserId UserID
	//
	// Enum and nested type conversions:
	// goverter:map Flags | ConvertJobFlags
	// goverter:map MailType | ConvertJobMailType
	// goverter:map Profile | ConvertJobProfile
	// goverter:map Shared | ConvertJobShared
	// goverter:map ExitCode | ConvertExitCode
	// goverter:map DerivedExitCode | ConvertExitCode
	//
	// Job-specific complex conversions:
	// goverter:map Power | ConvertJobPower
	// goverter:map GresDetail GRESDetail | ConvertJobGRESDetail
	//
	// Complex nested type conversions:
	// goverter:map JobResources | ConvertJobResources
	//
	// Fields that don't exist in v0_0_40 (added in later versions):
	// goverter:ignore PriorityByPartition
	// goverter:ignore StepID
	// goverter:ignore LicensesAllocated
	// goverter:ignore SegmentSize
	// goverter:ignore StderrExpanded
	// goverter:ignore StdinExpanded
	// goverter:ignore StdoutExpanded
	// goverter:ignore SubmitLine
	// goverter:ignore RequiredSwitches
	ConvertAPIJobToCommon(source api.V0040JobInfo) *types.Job
}
