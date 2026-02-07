// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_43

import (
	"time"

	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_43"
)

// =============================================================================
// Write Converter Interfaces (common -> API)
// =============================================================================
// These interfaces define goverter converters for converting common types to
// API types, used for create and update operations.
// AccountWriteConverterGoverter defines the goverter interface for Account write conversions.
//
// goverter:converter
// goverter:output:file write_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_43
// goverter:extend ConvertCoordNamesToSlice
// goverter:extend ConvertStringPtrToString
//
//go:generate goverter gen .
type AccountWriteConverterGoverter interface {
	// ConvertCommonAccountCreateToAPI converts common AccountCreate to API V0043Account type
	// goverter:map Coordinators | ConvertCoordNamesToSlice
	// goverter:ignore Associations
	// goverter:ignore Flags
	ConvertCommonAccountCreateToAPI(source *types.AccountCreate) *api.V0043Account
	// ConvertCommonAccountUpdateToAPI converts common AccountUpdate to API V0043Account type
	// goverter:map Coordinators | ConvertCoordNamesToSlice
	// goverter:map Description | ConvertStringPtrToString
	// goverter:map Organization | ConvertStringPtrToString
	// goverter:ignore Associations
	// goverter:ignore Flags
	// goverter:ignore Name
	ConvertCommonAccountUpdateToAPI(source *types.AccountUpdate) *api.V0043Account
}

// AssociationWriteConverterGoverter defines the goverter interface for Association write conversions.
// goverter:converter
// goverter:output:file write_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_43
// goverter:extend ConvertInt32ToUint32NoValStruct
type AssociationWriteConverterGoverter interface {
	// ConvertCommonAssociationCreateToAPI converts common AssociationCreate to API V0043Assoc type
	// goverter:map Account Account
	// goverter:map Cluster Cluster
	// goverter:map Comment Comment
	// goverter:map IsDefault IsDefault
	// goverter:map ParentAccount ParentAccount
	// goverter:map Partition Partition
	// goverter:map Priority | ConvertInt32ToUint32NoValStruct
	// goverter:map SharesRaw SharesRaw
	// goverter:map User User
	// goverter:ignore Accounting
	// goverter:ignore Default
	// goverter:ignore Flags
	// goverter:ignore Id
	// goverter:ignore Lineage
	// goverter:ignore Max
	// goverter:ignore Min
	// goverter:ignore Qos
	ConvertCommonAssociationCreateToAPI(source *types.AssociationCreate) *api.V0043Assoc
	// ConvertCommonAssociationUpdateToAPI converts common AssociationUpdate to API V0043Assoc type
	// goverter:map Comment Comment
	// goverter:map IsDefault IsDefault
	// goverter:map Priority | ConvertInt32PtrToUint32NoValStruct
	// goverter:map SharesRaw SharesRaw
	// goverter:ignore Account
	// goverter:ignore Accounting
	// goverter:ignore Cluster
	// goverter:ignore Default
	// goverter:ignore Flags
	// goverter:ignore Id
	// goverter:ignore Lineage
	// goverter:ignore Max
	// goverter:ignore Min
	// goverter:ignore ParentAccount
	// goverter:ignore Partition
	// goverter:ignore Qos
	// goverter:ignore User
	ConvertCommonAssociationUpdateToAPI(source *types.AssociationUpdate) *api.V0043Assoc
}

// ClusterWriteConverterGoverter defines the goverter interface for Cluster write conversions.
// goverter:converter
// goverter:output:file write_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_43
// goverter:extend ConvertStringSliceToClusterRecFlags
type ClusterWriteConverterGoverter interface {
	// ConvertCommonClusterCreateToAPI converts common ClusterCreate to API V0043ClusterRec type
	// goverter:map Flags | ConvertStringSliceToClusterRecFlags
	// goverter:map Name Name
	// goverter:map Nodes Nodes
	// goverter:map RpcVersion RpcVersion
	// goverter:map SelectPlugin SelectPlugin
	// goverter:ignore Associations
	// goverter:ignore Controller
	// goverter:ignore Tres
	ConvertCommonClusterCreateToAPI(source *types.ClusterCreate) *api.V0043ClusterRec
	// ConvertCommonClusterUpdateToAPI converts common ClusterUpdate to API V0043ClusterRec type
	// goverter:map RPCVersion RpcVersion
	// goverter:ignore Associations
	// goverter:ignore Controller
	// goverter:ignore Flags
	// goverter:ignore Name
	// goverter:ignore Nodes
	// goverter:ignore SelectPlugin
	// goverter:ignore Tres
	ConvertCommonClusterUpdateToAPI(source *types.ClusterUpdate) *api.V0043ClusterRec
}

// NodeWriteConverterGoverter defines the goverter interface for Node write conversions.
// Converts to V0043UpdateNodeMsg which is the API type for node update requests.
// goverter:converter
// goverter:output:file write_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_43
// goverter:extend ConvertStringSliceToCsvString
// goverter:extend ConvertStringSliceToHostlistString
// goverter:extend ConvertNodeStatesToAPIV43
// goverter:extend ConvertUint32PtrToNoValStruct
type NodeWriteConverterGoverter interface {
	// ConvertCommonNodeUpdateToAPI converts common NodeUpdate to API V0043UpdateNodeMsg type
	// Field mappings (common -> API field names)
	// goverter:map CPUBind CpuBind
	// goverter:map GRES Gres
	// goverter:map Features Features | ConvertStringSliceToCsvString
	// goverter:map FeaturesAct FeaturesAct | ConvertStringSliceToCsvString
	// goverter:map Address Address | ConvertStringSliceToHostlistString
	// goverter:map Hostname Hostname | ConvertStringSliceToHostlistString
	// goverter:map Name Name | ConvertStringSliceToHostlistString
	// goverter:map State State | ConvertNodeStatesToAPIV43
	// goverter:map Weight Weight | ConvertUint32PtrToNoValStruct
	// goverter:map ResumeAfter ResumeAfter | ConvertUint32PtrToNoValStruct
	// goverter:map ReasonUID ReasonUid
	// goverter:map TopologyStr TopologyStr
	ConvertCommonNodeUpdateToAPI(source *types.NodeUpdate) *api.V0043UpdateNodeMsg
}

// QoSWriteConverterGoverter defines the goverter interface for QoS write conversions.
// goverter:converter
// goverter:output:file write_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_43
// goverter:extend ConvertStringSliceToQosFlags
// goverter:extend ConvertIntToUint32NoValStruct
// goverter:extend ConvertFloat64ToFloat64NoValStruct
type QoSWriteConverterGoverter interface {
	// ConvertCommonQoSCreateToAPI converts common QoSCreate to API V0043Qos type
	// goverter:map Description Description
	// goverter:map Flags | ConvertStringSliceToQosFlags
	// goverter:map Name Name
	// goverter:map Priority | ConvertIntToUint32NoValStruct
	// goverter:map UsageFactor | ConvertFloat64ToFloat64NoValStruct
	// goverter:map UsageThreshold | ConvertFloat64ToFloat64NoValStruct
	// goverter:ignore Id
	// goverter:ignore Limits
	// goverter:ignore Preempt
	ConvertCommonQoSCreateToAPI(source *types.QoSCreate) *api.V0043Qos
	// ConvertCommonQoSUpdateToAPI converts common QoSUpdate to API V0043Qos type
	// goverter:map Description Description
	// goverter:map Flags | ConvertStringSlicePtrToQosFlags
	// goverter:map Priority | ConvertIntPtrToUint32NoValStruct
	// goverter:map UsageFactor | ConvertFloat64PtrToFloat64NoValStruct
	// goverter:map UsageThreshold | ConvertFloat64PtrToFloat64NoValStruct
	// goverter:ignore Id
	// goverter:ignore Limits
	// goverter:ignore Name
	// goverter:ignore Preempt
	ConvertCommonQoSUpdateToAPI(source *types.QoSUpdate) *api.V0043Qos
}

// WCKeyWriteConverterGoverter defines the goverter interface for WCKey write conversions.
// goverter:converter
// goverter:output:file write_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_43
type WCKeyWriteConverterGoverter interface {
	// ConvertCommonWCKeyCreateToAPI converts common WCKeyCreate to API V0043Wckey type
	// goverter:map Cluster Cluster
	// goverter:map Name Name
	// goverter:map User User
	// goverter:ignore Accounting
	// goverter:ignore Flags
	// goverter:ignore Id
	ConvertCommonWCKeyCreateToAPI(source *types.WCKeyCreate) *api.V0043Wckey
}

// JobWriteConverterGoverter defines the goverter interface for Job write conversions.
// With the generated JobCreate type that aligns with the API structure, most fields
// can be auto-mapped with the help of type conversion functions.
// goverter:converter
// goverter:output:file write_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_43
// goverter:extend ConvertUint64PtrToNoValStruct
// goverter:extend ConvertUint32PtrToNoValStruct
// goverter:extend ConvertUint16PtrToNoValStruct
// goverter:extend ConvertStringSliceToStringArray
// goverter:extend ConvertStringSliceToCsvString
// goverter:extend ConvertMailTypeSliceToAPIV43
// goverter:extend ConvertFlagsSliceToAPIV43
// goverter:extend ConvertCPUBindingFlagsSliceToAPIV43
// goverter:extend ConvertKillWarningFlagsSliceToAPIV43
// goverter:extend ConvertMemoryBindingTypeSliceToAPIV43
// goverter:extend ConvertOpenModeSliceToAPIV43
// goverter:extend ConvertProfileSliceToAPIV43
// goverter:extend ConvertSharedSliceToAPIV43
// goverter:extend ConvertX11SliceToAPIV43
// goverter:extend ConvertCronEntryToAPIV43
type JobWriteConverterGoverter interface {
	// ConvertCommonJobCreateToAPI converts common JobCreate to API V0043JobDescMsg type
	// Direct field mappings (same name, needs type conversion)
	// goverter:map BeginTime BeginTime | ConvertUint64PtrToNoValStruct
	// goverter:map TimeLimit TimeLimit | ConvertUint32PtrToNoValStruct
	// goverter:map TimeMinimum TimeMinimum | ConvertUint32PtrToNoValStruct
	// goverter:map Priority Priority | ConvertUint32PtrToNoValStruct
	// goverter:map RequiredSwitches RequiredSwitches | ConvertUint32PtrToNoValStruct
	// goverter:map MemoryPerCPU MemoryPerCpu | ConvertUint64PtrToNoValStruct
	// goverter:map MemoryPerNode MemoryPerNode | ConvertUint64PtrToNoValStruct
	// goverter:map DistributionPlaneSize DistributionPlaneSize | ConvertUint16PtrToNoValStruct
	// goverter:map KillWarningDelay KillWarningDelay | ConvertUint16PtrToNoValStruct
	// goverter:map SegmentSize SegmentSize | ConvertUint16PtrToNoValStruct
	// goverter:map Argv Argv | ConvertStringSliceToStringArray
	// goverter:map Environment Environment | ConvertStringSliceToStringArray
	// goverter:map SpankEnvironment SpankEnvironment | ConvertStringSliceToStringArray
	// goverter:map ExcludedNodes ExcludedNodes | ConvertStringSliceToCsvString
	// goverter:map RequiredNodes RequiredNodes | ConvertStringSliceToCsvString
	// goverter:map MailType MailType | ConvertMailTypeSliceToAPIV43
	// goverter:map Flags Flags | ConvertFlagsSliceToAPIV43
	// goverter:map CPUBindingFlags CpuBindingFlags | ConvertCPUBindingFlagsSliceToAPIV43
	// goverter:map KillWarningFlags KillWarningFlags | ConvertKillWarningFlagsSliceToAPIV43
	// goverter:map MemoryBindingType MemoryBindingType | ConvertMemoryBindingTypeSliceToAPIV43
	// goverter:map OpenMode OpenMode | ConvertOpenModeSliceToAPIV43
	// goverter:map Profile Profile | ConvertProfileSliceToAPIV43
	// goverter:map Shared Shared | ConvertSharedSliceToAPIV43
	// goverter:map X11 X11 | ConvertX11SliceToAPIV43
	// goverter:map Crontab Crontab | ConvertCronEntryToAPIV43
	// goverter:ignore Rlimits
	// Field name casing differences (ID/Id, CPU/Cpu)
	// goverter:map ContainerID ContainerId
	// goverter:map GroupID GroupId
	// goverter:map JobID JobId
	// goverter:map UserID UserId
	// goverter:map CPUBinding CpuBinding
	// goverter:map CPUFrequency CpuFrequency
	// goverter:map CPUsPerTask CpusPerTask
	// goverter:map CPUsPerTRES CpusPerTres
	// goverter:map MaximumCPUs MaximumCpus
	// goverter:map MinimumCPUs MinimumCpus
	// goverter:map MinimumCPUsPerNode MinimumCpusPerNode
	// goverter:map MCSLabel McsLabel
	// goverter:map QoS Qos
	// goverter:map TRESBind TresBind
	// goverter:map TRESFreq TresFreq
	// goverter:map TRESPerJob TresPerJob
	// goverter:map TRESPerNode TresPerNode
	// goverter:map TRESPerSocket TresPerSocket
	// goverter:map TRESPerTask TresPerTask
	// goverter:map NtasksPerTRES NtasksPerTres
	// goverter:map MemoryPerTRES MemoryPerTres
	// Ignore deprecated/complex fields
	// goverter:ignore PowerFlags
	ConvertCommonJobCreateToAPI(source *types.JobCreate) *api.V0043JobDescMsg
}

// UserWriteConverterGoverter defines the goverter interface for User write conversions.
// Note: The Default field is an anonymous struct which requires manual handling via helper.
// goverter:converter
// goverter:output:file write_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_43
// goverter:extend ConvertUserCreateDefaultToAPI
// goverter:extend ConvertUserUpdateDefaultToAPI
type UserWriteConverterGoverter interface {
	// ConvertCommonUserCreateToAPI converts common UserCreate to API V0043User type
	// goverter:map Name Name
	// goverter:map . Default | ConvertUserCreateDefaultToAPI
	// goverter:ignore AdministratorLevel
	// goverter:ignore Associations
	// goverter:ignore Coordinators
	// goverter:ignore Flags
	// goverter:ignore OldName
	// goverter:ignore Wckeys
	ConvertCommonUserCreateToAPI(source *types.UserCreate) *api.V0043User
	// ConvertCommonUserUpdateToAPI converts common UserUpdate to API V0043User type
	// goverter:map . Default | ConvertUserUpdateDefaultToAPI
	// goverter:ignore AdministratorLevel
	// goverter:ignore Associations
	// goverter:ignore Coordinators
	// goverter:ignore Flags
	// goverter:ignore Name
	// goverter:ignore OldName
	// goverter:ignore Wckeys
	ConvertCommonUserUpdateToAPI(source *types.UserUpdate) *api.V0043User
}

// ReservationWriteConverterGoverter defines the goverter interface for Reservation write conversions.
// Note: Reservation uses V0043ReservationDescMsg for create/update operations, not V0043ReservationInfo.
// goverter:converter
// goverter:output:file write_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_43
// goverter:extend ConvertTimeToUint64NoValStruct
// goverter:extend ConvertTimePtrToUint64NoValStruct
// goverter:extend ConvertUint32PtrToNoValStruct
// goverter:extend ConvertStringSliceToCsvString
// goverter:extend ConvertStringSliceToHostlistString
// goverter:extend ConvertFlagsValueSliceToReservationDescMsgFlags
// goverter:extend ConvertReservationFlagSliceToReservationDescMsgFlags
// goverter:extend ConvertStringSliceToJoinedString
// goverter:extend ConvertStringPtrToHostlistString
type ReservationWriteConverterGoverter interface {
	// ConvertCommonReservationCreateToAPI converts common ReservationCreate to API V0043ReservationDescMsg type
	// Time fields:
	// goverter:map StartTime StartTime | ConvertTimeToUint64NoValStruct
	// goverter:map EndTime EndTime | ConvertTimeToUint64NoValStruct
	// NoValStruct fields:
	// goverter:map Duration Duration | ConvertUint32PtrToNoValStruct
	// goverter:map NodeCount NodeCount | ConvertUint32PtrToNoValStruct
	// goverter:map CoreCount CoreCount | ConvertUint32PtrToNoValStruct
	// goverter:map MaxStartDelay MaxStartDelay | ConvertUint32PtrToNoValStruct
	// CSV/Hostlist fields:
	// goverter:map Users Users | ConvertStringSliceToCsvString
	// goverter:map Accounts Accounts | ConvertStringSliceToCsvString
	// goverter:map Groups Groups | ConvertStringSliceToCsvString
	// goverter:map Licenses Licenses | ConvertStringSliceToCsvString
	// goverter:map NodeList NodeList | ConvertStringSliceToHostlistString
	// Flags:
	// goverter:map Flags Flags | ConvertFlagsValueSliceToReservationDescMsgFlags
	// Ignore complex/unsupported fields:
	// goverter:ignore PurgeCompleted
	// goverter:ignore Tres
	ConvertCommonReservationCreateToAPI(source *types.ReservationCreate) *api.V0043ReservationDescMsg

	// ConvertCommonReservationUpdateToAPI converts common ReservationUpdate to API V0043ReservationDescMsg type
	// Time fields:
	// goverter:map StartTime StartTime | ConvertTimePtrToUint64NoValStruct
	// goverter:map EndTime EndTime | ConvertTimePtrToUint64NoValStruct
	// NoValStruct fields:
	// goverter:map Duration Duration | ConvertInt32PtrToUint32NoValStruct
	// goverter:map NodeCount NodeCount | ConvertInt32PtrToUint32NoValStruct
	// goverter:map CoreCount CoreCount | ConvertInt32PtrToUint32NoValStruct
	// goverter:map MaxStartDelay MaxStartDelay | ConvertInt32PtrToUint32NoValStruct
	// CSV/Hostlist fields:
	// goverter:map Users Users | ConvertStringSliceToCsvString
	// goverter:map Accounts Accounts | ConvertStringSliceToCsvString
	// goverter:map Groups Groups | ConvertStringSliceToCsvString
	// goverter:map Features Features | ConvertStringSliceToJoinedString
	// goverter:map NodeList NodeList | ConvertStringPtrToHostlistString
	// Flags:
	// goverter:map Flags Flags | ConvertReservationFlagSliceToReservationDescMsgFlags
	// Ignore complex/unsupported fields:
	// goverter:ignore Name
	// goverter:ignore PurgeCompleted
	// goverter:ignore Tres
	// goverter:ignore Licenses
	ConvertCommonReservationUpdateToAPI(source *types.ReservationUpdate) *api.V0043ReservationDescMsg
}
// =============================================================================
// Write Helper Functions (common -> API)
// =============================================================================
// These helper functions handle complex type conversions that goverter cannot
// auto-generate for write operations (common types to API types).
// ConvertInt32ToUint32NoValStruct converts an int32 to API V0043Uint32NoValStruct.
// Used for Priority fields in AssociationCreate.
func ConvertInt32ToUint32NoValStruct(source int32) *api.V0043Uint32NoValStruct {
	if source == 0 {
		return nil
	}
	setTrue := true
	num := int32(source)
	return &api.V0043Uint32NoValStruct{
		Set:    &setTrue,
		Number: &num,
	}
}

// ConvertInt32PtrToUint32NoValStruct converts a *int32 to API V0043Uint32NoValStruct.
// Used for Priority fields in AssociationUpdate.
func ConvertInt32PtrToUint32NoValStruct(source *int32) *api.V0043Uint32NoValStruct {
	if source == nil {
		return nil
	}
	setTrue := true
	num := *source
	return &api.V0043Uint32NoValStruct{
		Set:    &setTrue,
		Number: &num,
	}
}

// ConvertIntToUint32NoValStruct converts an int to API V0043Uint32NoValStruct.
// Used for Priority fields in QoSCreate.
func ConvertIntToUint32NoValStruct(source int) *api.V0043Uint32NoValStruct {
	if source == 0 {
		return nil
	}
	setTrue := true
	num := int32(source)
	return &api.V0043Uint32NoValStruct{
		Set:    &setTrue,
		Number: &num,
	}
}

// ConvertIntPtrToUint32NoValStruct converts a *int to API V0043Uint32NoValStruct.
// Used for Priority fields in QoSUpdate.
func ConvertIntPtrToUint32NoValStruct(source *int) *api.V0043Uint32NoValStruct {
	if source == nil {
		return nil
	}
	setTrue := true
	num := int32(*source)
	return &api.V0043Uint32NoValStruct{
		Set:    &setTrue,
		Number: &num,
	}
}

// ConvertFloat64ToFloat64NoValStruct converts a float64 to API V0043Float64NoValStruct.
// Used for UsageFactor/UsageThreshold fields in QoSCreate.
func ConvertFloat64ToFloat64NoValStruct(source float64) *api.V0043Float64NoValStruct {
	if source == 0 {
		return nil
	}
	setTrue := true
	return &api.V0043Float64NoValStruct{
		Set:    &setTrue,
		Number: &source,
	}
}

// ConvertFloat64PtrToFloat64NoValStruct converts a *float64 to API V0043Float64NoValStruct.
// Used for UsageFactor/UsageThreshold fields in QoSUpdate.
func ConvertFloat64PtrToFloat64NoValStruct(source *float64) *api.V0043Float64NoValStruct {
	if source == nil {
		return nil
	}
	setTrue := true
	return &api.V0043Float64NoValStruct{
		Set:    &setTrue,
		Number: source,
	}
}

// ConvertStringSliceToClusterRecFlags converts []string to *[]V0043ClusterRecFlags.
// Used for Flags in ClusterCreate.
func ConvertStringSliceToClusterRecFlags(source []string) *[]api.V0043ClusterRecFlags {
	if len(source) == 0 {
		return nil
	}
	flags := make([]api.V0043ClusterRecFlags, len(source))
	for i, f := range source {
		flags[i] = api.V0043ClusterRecFlags(f)
	}
	return &flags
}

// ConvertStringSliceToQosFlags converts []string to *[]V0043QosFlags.
// Used for Flags in QoSCreate.
func ConvertStringSliceToQosFlags(source []string) *[]api.V0043QosFlags {
	if len(source) == 0 {
		return nil
	}
	flags := make([]api.V0043QosFlags, len(source))
	for i, f := range source {
		flags[i] = api.V0043QosFlags(f)
	}
	return &flags
}

// ConvertStringSlicePtrToQosFlags converts *[]string to *[]V0043QosFlags.
// Used for Flags in QoSUpdate.
func ConvertStringSlicePtrToQosFlags(source *[]string) *[]api.V0043QosFlags {
	if source == nil || len(*source) == 0 {
		return nil
	}
	flags := make([]api.V0043QosFlags, len(*source))
	for i, f := range *source {
		flags[i] = api.V0043QosFlags(f)
	}
	return &flags
}

// ConvertStringPtrToString converts a *string to string.
// Used for fields like Description and Organization in AccountUpdate.
func ConvertStringPtrToString(source *string) string {
	if source == nil {
		return ""
	}
	return *source
}

// ConvertUserCreateDefaultToAPI converts UserCreate fields to the anonymous Default struct.
// Used by goverter to map DefaultAccount and DefaultWCKey to the nested Default field.
func ConvertUserCreateDefaultToAPI(source *types.UserCreate) *struct {
	Account *string `json:"account,omitempty"`
	Wckey   *string `json:"wckey,omitempty"`
} {
	if source == nil {
		return nil
	}
	if source.DefaultAccount == "" && source.DefaultWCKey == "" {
		return nil
	}
	result := &struct {
		Account *string `json:"account,omitempty"`
		Wckey   *string `json:"wckey,omitempty"`
	}{}
	if source.DefaultAccount != "" {
		result.Account = &source.DefaultAccount
	}
	if source.DefaultWCKey != "" {
		result.Wckey = &source.DefaultWCKey
	}
	return result
}

// ConvertUserUpdateDefaultToAPI converts UserUpdate fields to the anonymous Default struct.
// Used by goverter to map DefaultAccount and DefaultWCKey to the nested Default field.
func ConvertUserUpdateDefaultToAPI(source *types.UserUpdate) *struct {
	Account *string `json:"account,omitempty"`
	Wckey   *string `json:"wckey,omitempty"`
} {
	if source == nil {
		return nil
	}
	if source.DefaultAccount == nil && source.DefaultWCKey == nil {
		return nil
	}
	result := &struct {
		Account *string `json:"account,omitempty"`
		Wckey   *string `json:"wckey,omitempty"`
	}{}
	if source.DefaultAccount != nil {
		result.Account = source.DefaultAccount
	}
	if source.DefaultWCKey != nil {
		result.Wckey = source.DefaultWCKey
	}
	return result
}

// =============================================================================
// Node Write Helper Functions (NodeUpdate -> V0043UpdateNodeMsg)
// =============================================================================
// ConvertStringSliceToHostlistString converts []string to *V0043HostlistString.
// Used for Address, Hostname, Name fields.
func ConvertStringSliceToHostlistString(source []string) *api.V0043HostlistString {
	if len(source) == 0 {
		return nil
	}
	hl := api.V0043HostlistString(source)
	return &hl
}

// ConvertStringSliceToCsvString converts []string to *V0043CsvString.
// Used for Features, FeaturesAct fields.
func ConvertStringSliceToCsvString(source []string) *api.V0043CsvString {
	if len(source) == 0 {
		return nil
	}
	csv := api.V0043CsvString(source)
	return &csv
}

// ConvertNodeStatesToAPIV43 converts []NodeState to *[]V0043UpdateNodeMsgState (enum).
func ConvertNodeStatesToAPIV43(source []types.NodeState) *[]api.V0043UpdateNodeMsgState {
	if len(source) == 0 {
		return nil
	}
	result := make([]api.V0043UpdateNodeMsgState, len(source))
	for i, s := range source {
		result[i] = api.V0043UpdateNodeMsgState(s)
	}
	return &result
}

// ConvertUint32PtrToNoValStruct converts *uint32 to *V0043Uint32NoValStruct.
// Used for fields like Weight, ResumeAfter.
func ConvertUint32PtrToNoValStruct(source *uint32) *api.V0043Uint32NoValStruct {
	if source == nil {
		return nil
	}
	setTrue := true
	val := int32(*source)
	return &api.V0043Uint32NoValStruct{
		Set:    &setTrue,
		Number: &val,
	}
}

// =============================================================================
// Job Write Helper Functions (JobCreate -> V0043JobDescMsg)
// =============================================================================
// ConvertUint64PtrToNoValStruct converts *uint64 to *V0043Uint64NoValStruct.
// Used for fields like BeginTime, MemoryPerCPU, MemoryPerNode.
func ConvertUint64PtrToNoValStruct(source *uint64) *api.V0043Uint64NoValStruct {
	if source == nil {
		return nil
	}
	setTrue := true
	val := int64(*source)
	return &api.V0043Uint64NoValStruct{
		Set:    &setTrue,
		Number: &val,
	}
}

// ConvertUint16PtrToNoValStruct converts *uint16 to *V0043Uint16NoValStruct.
// Used for fields like DistributionPlaneSize, KillWarningDelay, SegmentSize.
func ConvertUint16PtrToNoValStruct(source *uint16) *api.V0043Uint16NoValStruct {
	if source == nil {
		return nil
	}
	setTrue := true
	val := int32(*source)
	return &api.V0043Uint16NoValStruct{
		Set:    &setTrue,
		Number: &val,
	}
}

// ConvertStringSliceToStringArray converts []string to *V0043StringArray.
// Used for Environment, Argv, SpankEnvironment.
func ConvertStringSliceToStringArray(source []string) *api.V0043StringArray {
	if len(source) == 0 {
		return nil
	}
	arr := api.V0043StringArray(source)
	return &arr
}

// ConvertMailTypeSliceToAPIV43 converts []MailTypeValue to *[]V0043JobDescMsgMailType.
func ConvertMailTypeSliceToAPIV43(source []types.MailTypeValue) *[]api.V0043JobDescMsgMailType {
	if len(source) == 0 {
		return nil
	}
	result := make([]api.V0043JobDescMsgMailType, len(source))
	for i, v := range source {
		result[i] = api.V0043JobDescMsgMailType(v)
	}
	return &result
}

// ConvertFlagsSliceToAPIV43 converts []FlagsValue to *[]V0043JobDescMsgFlags.
func ConvertFlagsSliceToAPIV43(source []types.FlagsValue) *[]api.V0043JobDescMsgFlags {
	if len(source) == 0 {
		return nil
	}
	result := make([]api.V0043JobDescMsgFlags, len(source))
	for i, v := range source {
		result[i] = api.V0043JobDescMsgFlags(v)
	}
	return &result
}

// ConvertCPUBindingFlagsSliceToAPIV43 converts []CPUBindingFlagsValue to *[]V0043JobDescMsgCpuBindingFlags.
func ConvertCPUBindingFlagsSliceToAPIV43(source []types.CPUBindingFlagsValue) *[]api.V0043JobDescMsgCpuBindingFlags {
	if len(source) == 0 {
		return nil
	}
	result := make([]api.V0043JobDescMsgCpuBindingFlags, len(source))
	for i, v := range source {
		result[i] = api.V0043JobDescMsgCpuBindingFlags(v)
	}
	return &result
}

// ConvertKillWarningFlagsSliceToAPIV43 converts []KillWarningFlagsValue to *[]V0043JobDescMsgKillWarningFlags.
func ConvertKillWarningFlagsSliceToAPIV43(source []types.KillWarningFlagsValue) *[]api.V0043JobDescMsgKillWarningFlags {
	if len(source) == 0 {
		return nil
	}
	result := make([]api.V0043JobDescMsgKillWarningFlags, len(source))
	for i, v := range source {
		result[i] = api.V0043JobDescMsgKillWarningFlags(v)
	}
	return &result
}

// ConvertMemoryBindingTypeSliceToAPIV43 converts []MemoryBindingTypeValue to *[]V0043JobDescMsgMemoryBindingType.
func ConvertMemoryBindingTypeSliceToAPIV43(source []types.MemoryBindingTypeValue) *[]api.V0043JobDescMsgMemoryBindingType {
	if len(source) == 0 {
		return nil
	}
	result := make([]api.V0043JobDescMsgMemoryBindingType, len(source))
	for i, v := range source {
		result[i] = api.V0043JobDescMsgMemoryBindingType(v)
	}
	return &result
}

// ConvertOpenModeSliceToAPIV43 converts []OpenModeValue to *[]V0043JobDescMsgOpenMode.
func ConvertOpenModeSliceToAPIV43(source []types.OpenModeValue) *[]api.V0043JobDescMsgOpenMode {
	if len(source) == 0 {
		return nil
	}
	result := make([]api.V0043JobDescMsgOpenMode, len(source))
	for i, v := range source {
		result[i] = api.V0043JobDescMsgOpenMode(v)
	}
	return &result
}

// ConvertProfileSliceToAPIV43 converts []ProfileValue to *[]V0043JobDescMsgProfile.
func ConvertProfileSliceToAPIV43(source []types.ProfileValue) *[]api.V0043JobDescMsgProfile {
	if len(source) == 0 {
		return nil
	}
	result := make([]api.V0043JobDescMsgProfile, len(source))
	for i, v := range source {
		result[i] = api.V0043JobDescMsgProfile(v)
	}
	return &result
}

// ConvertSharedSliceToAPIV43 converts []SharedValue to *[]V0043JobDescMsgShared.
func ConvertSharedSliceToAPIV43(source []types.SharedValue) *[]api.V0043JobDescMsgShared {
	if len(source) == 0 {
		return nil
	}
	result := make([]api.V0043JobDescMsgShared, len(source))
	for i, v := range source {
		result[i] = api.V0043JobDescMsgShared(v)
	}
	return &result
}

// ConvertX11SliceToAPIV43 converts []X11Value to *[]V0043JobDescMsgX11.
func ConvertX11SliceToAPIV43(source []types.X11Value) *[]api.V0043JobDescMsgX11 {
	if len(source) == 0 {
		return nil
	}
	result := make([]api.V0043JobDescMsgX11, len(source))
	for i, v := range source {
		result[i] = api.V0043JobDescMsgX11(v)
	}
	return &result
}

// ConvertCronEntryToAPIV43 converts *CronEntry to *V0043CronEntry.
func ConvertCronEntryToAPIV43(source *types.CronEntry) *api.V0043CronEntry {
	if source == nil {
		return nil
	}
	result := &api.V0043CronEntry{
		Minute:        source.Minute,
		Hour:          source.Hour,
		DayOfMonth:    source.DayOfMonth,
		Month:         source.Month,
		DayOfWeek:     source.DayOfWeek,
		Specification: source.Specification,
		Command:       source.Command,
	}
	// Convert Flags
	if len(source.Flags) > 0 {
		flags := make([]api.V0043CronEntryFlags, len(source.Flags))
		for i, f := range source.Flags {
			flags[i] = api.V0043CronEntryFlags(f)
		}
		result.Flags = &flags
	}
	// Convert Line
	if source.Line != nil {
		result.Line = &struct {
			End   *int32 `json:"end,omitempty"`
			Start *int32 `json:"start,omitempty"`
		}{
			Start: source.Line.Start,
			End:   source.Line.End,
		}
	}
	return result
}

// =============================================================================
// Reservation Write Helper Functions (ReservationCreate/Update -> V0043ReservationDescMsg)
// =============================================================================

// ConvertTimeToUint64NoValStruct converts time.Time to *V0043Uint64NoValStruct.
// Used for StartTime, EndTime in ReservationCreate.
func ConvertTimeToUint64NoValStruct(source time.Time) *api.V0043Uint64NoValStruct {
	if source.IsZero() {
		return nil
	}
	setTrue := true
	unix := source.Unix()
	return &api.V0043Uint64NoValStruct{
		Set:    &setTrue,
		Number: &unix,
	}
}

// ConvertTimePtrToUint64NoValStruct converts *time.Time to *V0043Uint64NoValStruct.
// Used for StartTime, EndTime in ReservationUpdate.
func ConvertTimePtrToUint64NoValStruct(source *time.Time) *api.V0043Uint64NoValStruct {
	if source == nil || source.IsZero() {
		return nil
	}
	setTrue := true
	unix := source.Unix()
	return &api.V0043Uint64NoValStruct{
		Set:    &setTrue,
		Number: &unix,
	}
}

// ConvertFlagsValueSliceToReservationDescMsgFlags converts []FlagsValue to *[]V0043ReservationDescMsgFlags.
// Used for Flags in ReservationCreate.
func ConvertFlagsValueSliceToReservationDescMsgFlags(source []types.FlagsValue) *[]api.V0043ReservationDescMsgFlags {
	if len(source) == 0 {
		return nil
	}
	result := make([]api.V0043ReservationDescMsgFlags, len(source))
	for i, f := range source {
		result[i] = api.V0043ReservationDescMsgFlags(f)
	}
	return &result
}

// ConvertReservationFlagSliceToReservationDescMsgFlags converts []ReservationFlag to *[]V0043ReservationDescMsgFlags.
// Used for Flags in ReservationUpdate.
func ConvertReservationFlagSliceToReservationDescMsgFlags(source []types.ReservationFlag) *[]api.V0043ReservationDescMsgFlags {
	if len(source) == 0 {
		return nil
	}
	result := make([]api.V0043ReservationDescMsgFlags, len(source))
	for i, f := range source {
		result[i] = api.V0043ReservationDescMsgFlags(f)
	}
	return &result
}

// ConvertStringSliceToJoinedString converts []string to *string by joining with commas.
// Used for Features in ReservationUpdate where API expects a single comma-separated string.
func ConvertStringSliceToJoinedString(source []string) *string {
	if len(source) == 0 {
		return nil
	}
	joined := ""
	for i, s := range source {
		if i > 0 {
			joined += ","
		}
		joined += s
	}
	return &joined
}

// ConvertStringPtrToHostlistString converts *string to *V0043HostlistString.
// Used for NodeList in ReservationUpdate where the input is a single string.
func ConvertStringPtrToHostlistString(source *string) *api.V0043HostlistString {
	if source == nil || *source == "" {
		return nil
	}
	hl := api.V0043HostlistString([]string{*source})
	return &hl
}
