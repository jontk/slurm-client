// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_42

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_42"
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
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_42
// goverter:extend ConvertCoordNamesToSlice
// goverter:extend ConvertStringPtrToString
//
//go:generate goverter gen .
type AccountWriteConverterGoverter interface {
	// ConvertCommonAccountCreateToAPI converts common AccountCreate to API V0042Account type
	// goverter:map Coordinators | ConvertCoordNamesToSlice
	// goverter:ignore Associations
	// goverter:ignore Flags
	ConvertCommonAccountCreateToAPI(source *types.AccountCreate) *api.V0042Account
	// ConvertCommonAccountUpdateToAPI converts common AccountUpdate to API V0042Account type
	// goverter:map Coordinators | ConvertCoordNamesToSlice
	// goverter:map Description | ConvertStringPtrToString
	// goverter:map Organization | ConvertStringPtrToString
	// goverter:ignore Associations
	// goverter:ignore Flags
	// goverter:ignore Name
	ConvertCommonAccountUpdateToAPI(source *types.AccountUpdate) *api.V0042Account
}

// AssociationWriteConverterGoverter defines the goverter interface for Association write conversions.
// goverter:converter
// goverter:output:file write_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_42
// goverter:extend ConvertInt32ToUint32NoValStruct
type AssociationWriteConverterGoverter interface {
	// ConvertCommonAssociationCreateToAPI converts common AssociationCreate to API V0042Assoc type
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
	ConvertCommonAssociationCreateToAPI(source *types.AssociationCreate) *api.V0042Assoc
	// ConvertCommonAssociationUpdateToAPI converts common AssociationUpdate to API V0042Assoc type
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
	ConvertCommonAssociationUpdateToAPI(source *types.AssociationUpdate) *api.V0042Assoc
}

// ClusterWriteConverterGoverter defines the goverter interface for Cluster write conversions.
// goverter:converter
// goverter:output:file write_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_42
// goverter:extend ConvertStringSliceToClusterRecFlags
type ClusterWriteConverterGoverter interface {
	// ConvertCommonClusterCreateToAPI converts common ClusterCreate to API V0042ClusterRec type
	// goverter:map Flags | ConvertStringSliceToClusterRecFlags
	// goverter:map Name Name
	// goverter:map Nodes Nodes
	// goverter:map RpcVersion RpcVersion
	// goverter:map SelectPlugin SelectPlugin
	// goverter:ignore Associations
	// goverter:ignore Controller
	// goverter:ignore Tres
	ConvertCommonClusterCreateToAPI(source *types.ClusterCreate) *api.V0042ClusterRec
	// ConvertCommonClusterUpdateToAPI converts common ClusterUpdate to API V0042ClusterRec type
	// goverter:map RPCVersion RpcVersion
	// goverter:ignore Associations
	// goverter:ignore Controller
	// goverter:ignore Flags
	// goverter:ignore Name
	// goverter:ignore Nodes
	// goverter:ignore SelectPlugin
	// goverter:ignore Tres
	ConvertCommonClusterUpdateToAPI(source *types.ClusterUpdate) *api.V0042ClusterRec
}

// NodeWriteConverterGoverter defines the goverter interface for Node write conversions.
// Converts to V0042UpdateNodeMsg which is the API type for node update requests.
// goverter:converter
// goverter:output:file write_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_42
// goverter:extend ConvertStringSliceToCsvString
// goverter:extend ConvertStringSliceToHostlistString
// goverter:extend ConvertNodeStatesToAPIV42
// goverter:extend ConvertUint32PtrToNoValStruct
type NodeWriteConverterGoverter interface {
	// ConvertCommonNodeUpdateToAPI converts common NodeUpdate to API V0042UpdateNodeMsg type
	// Field mappings (common -> API field names)
	// goverter:map CPUBind CpuBind
	// goverter:map GRES Gres
	// goverter:map Features Features | ConvertStringSliceToCsvString
	// goverter:map FeaturesAct FeaturesAct | ConvertStringSliceToCsvString
	// goverter:map Address Address | ConvertStringSliceToHostlistString
	// goverter:map Hostname Hostname | ConvertStringSliceToHostlistString
	// goverter:map Name Name | ConvertStringSliceToHostlistString
	// goverter:map State State | ConvertNodeStatesToAPIV42
	// goverter:map Weight Weight | ConvertUint32PtrToNoValStruct
	// goverter:map ResumeAfter ResumeAfter | ConvertUint32PtrToNoValStruct
	// goverter:map ReasonUID ReasonUid
	ConvertCommonNodeUpdateToAPI(source *types.NodeUpdate) *api.V0042UpdateNodeMsg
}

// QoSWriteConverterGoverter defines the goverter interface for QoS write conversions.
// goverter:converter
// goverter:output:file write_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_42
// goverter:extend ConvertStringSliceToQosFlags
// goverter:extend ConvertIntToUint32NoValStruct
// goverter:extend ConvertFloat64ToFloat64NoValStruct
type QoSWriteConverterGoverter interface {
	// ConvertCommonQoSCreateToAPI converts common QoSCreate to API V0042Qos type
	// goverter:map Description Description
	// goverter:map Flags | ConvertStringSliceToQosFlags
	// goverter:map Name Name
	// goverter:map Priority | ConvertIntToUint32NoValStruct
	// goverter:map UsageFactor | ConvertFloat64ToFloat64NoValStruct
	// goverter:map UsageThreshold | ConvertFloat64ToFloat64NoValStruct
	// goverter:ignore Id
	// goverter:ignore Limits
	// goverter:ignore Preempt
	ConvertCommonQoSCreateToAPI(source *types.QoSCreate) *api.V0042Qos
	// ConvertCommonQoSUpdateToAPI converts common QoSUpdate to API V0042Qos type
	// goverter:map Description Description
	// goverter:map Flags | ConvertStringSlicePtrToQosFlags
	// goverter:map Priority | ConvertIntPtrToUint32NoValStruct
	// goverter:map UsageFactor | ConvertFloat64PtrToFloat64NoValStruct
	// goverter:map UsageThreshold | ConvertFloat64PtrToFloat64NoValStruct
	// goverter:ignore Id
	// goverter:ignore Limits
	// goverter:ignore Name
	// goverter:ignore Preempt
	ConvertCommonQoSUpdateToAPI(source *types.QoSUpdate) *api.V0042Qos
}

// WCKeyWriteConverterGoverter defines the goverter interface for WCKey write conversions.
// goverter:converter
// goverter:output:file write_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_42
type WCKeyWriteConverterGoverter interface {
	// ConvertCommonWCKeyCreateToAPI converts common WCKeyCreate to API V0042Wckey type
	// goverter:map Cluster Cluster
	// goverter:map Name Name
	// goverter:map User User
	// goverter:ignore Accounting
	// goverter:ignore Flags
	// goverter:ignore Id
	ConvertCommonWCKeyCreateToAPI(source *types.WCKeyCreate) *api.V0042Wckey
}

// JobWriteConverterGoverter defines the goverter interface for Job write conversions.
// With the generated JobCreate type that aligns with the API structure, most fields
// can be auto-mapped with the help of type conversion functions.
// goverter:converter
// goverter:output:file write_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_42
// goverter:extend ConvertUint64PtrToNoValStruct
// goverter:extend ConvertUint32PtrToNoValStruct
// goverter:extend ConvertUint16PtrToNoValStruct
// goverter:extend ConvertStringSliceToStringArray
// goverter:extend ConvertStringSliceToCsvString
// goverter:extend ConvertMailTypeSliceToAPIV42
// goverter:extend ConvertFlagsSliceToAPIV42
// goverter:extend ConvertCPUBindingFlagsSliceToAPIV42
// goverter:extend ConvertKillWarningFlagsSliceToAPIV42
// goverter:extend ConvertMemoryBindingTypeSliceToAPIV42
// goverter:extend ConvertOpenModeSliceToAPIV42
// goverter:extend ConvertProfileSliceToAPIV42
// goverter:extend ConvertSharedSliceToAPIV42
// goverter:extend ConvertX11SliceToAPIV42
// goverter:extend ConvertCronEntryToAPIV42
type JobWriteConverterGoverter interface {
	// ConvertCommonJobCreateToAPI converts common JobCreate to API V0042JobDescMsg type
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
	// goverter:map MailType MailType | ConvertMailTypeSliceToAPIV42
	// goverter:map Flags Flags | ConvertFlagsSliceToAPIV42
	// goverter:map CPUBindingFlags CpuBindingFlags | ConvertCPUBindingFlagsSliceToAPIV42
	// goverter:map KillWarningFlags KillWarningFlags | ConvertKillWarningFlagsSliceToAPIV42
	// goverter:map MemoryBindingType MemoryBindingType | ConvertMemoryBindingTypeSliceToAPIV42
	// goverter:map OpenMode OpenMode | ConvertOpenModeSliceToAPIV42
	// goverter:map Profile Profile | ConvertProfileSliceToAPIV42
	// goverter:map Shared Shared | ConvertSharedSliceToAPIV42
	// goverter:map X11 X11 | ConvertX11SliceToAPIV42
	// goverter:map Crontab Crontab | ConvertCronEntryToAPIV42
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
	ConvertCommonJobCreateToAPI(source *types.JobCreate) *api.V0042JobDescMsg
}

// UserWriteConverterGoverter defines the goverter interface for User write conversions.
// Note: The Default field is an anonymous struct which requires manual handling via helper.
// goverter:converter
// goverter:output:file write_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_42
// goverter:extend ConvertUserCreateDefaultToAPI
// goverter:extend ConvertUserUpdateDefaultToAPI
type UserWriteConverterGoverter interface {
	// ConvertCommonUserCreateToAPI converts common UserCreate to API V0042User type
	// goverter:map Name Name
	// goverter:map . Default | ConvertUserCreateDefaultToAPI
	// goverter:ignore AdministratorLevel
	// goverter:ignore Associations
	// goverter:ignore Coordinators
	// goverter:ignore Flags
	// goverter:ignore OldName
	// goverter:ignore Wckeys
	ConvertCommonUserCreateToAPI(source *types.UserCreate) *api.V0042User
	// ConvertCommonUserUpdateToAPI converts common UserUpdate to API V0042User type
	// goverter:map . Default | ConvertUserUpdateDefaultToAPI
	// goverter:ignore AdministratorLevel
	// goverter:ignore Associations
	// goverter:ignore Coordinators
	// goverter:ignore Flags
	// goverter:ignore Name
	// goverter:ignore OldName
	// goverter:ignore Wckeys
	ConvertCommonUserUpdateToAPI(source *types.UserUpdate) *api.V0042User
}

// =============================================================================
// Write Helper Functions (common -> API)
// =============================================================================
// These helper functions handle complex type conversions that goverter cannot
// auto-generate for write operations (common types to API types).
// ConvertInt32ToUint32NoValStruct converts an int32 to API V0042Uint32NoValStruct.
// Used for Priority fields in AssociationCreate.
func ConvertInt32ToUint32NoValStruct(source int32) *api.V0042Uint32NoValStruct {
	if source == 0 {
		return nil
	}
	setTrue := true
	num := int32(source)
	return &api.V0042Uint32NoValStruct{
		Set:    &setTrue,
		Number: &num,
	}
}

// ConvertInt32PtrToUint32NoValStruct converts a *int32 to API V0042Uint32NoValStruct.
// Used for Priority fields in AssociationUpdate.
func ConvertInt32PtrToUint32NoValStruct(source *int32) *api.V0042Uint32NoValStruct {
	if source == nil {
		return nil
	}
	setTrue := true
	num := *source
	return &api.V0042Uint32NoValStruct{
		Set:    &setTrue,
		Number: &num,
	}
}

// ConvertIntToUint32NoValStruct converts an int to API V0042Uint32NoValStruct.
// Used for Priority fields in QoSCreate.
func ConvertIntToUint32NoValStruct(source int) *api.V0042Uint32NoValStruct {
	if source == 0 {
		return nil
	}
	setTrue := true
	num := int32(source)
	return &api.V0042Uint32NoValStruct{
		Set:    &setTrue,
		Number: &num,
	}
}

// ConvertIntPtrToUint32NoValStruct converts a *int to API V0042Uint32NoValStruct.
// Used for Priority fields in QoSUpdate.
func ConvertIntPtrToUint32NoValStruct(source *int) *api.V0042Uint32NoValStruct {
	if source == nil {
		return nil
	}
	setTrue := true
	num := int32(*source)
	return &api.V0042Uint32NoValStruct{
		Set:    &setTrue,
		Number: &num,
	}
}

// ConvertFloat64ToFloat64NoValStruct converts a float64 to API V0042Float64NoValStruct.
// Used for UsageFactor/UsageThreshold fields in QoSCreate.
func ConvertFloat64ToFloat64NoValStruct(source float64) *api.V0042Float64NoValStruct {
	if source == 0 {
		return nil
	}
	setTrue := true
	return &api.V0042Float64NoValStruct{
		Set:    &setTrue,
		Number: &source,
	}
}

// ConvertFloat64PtrToFloat64NoValStruct converts a *float64 to API V0042Float64NoValStruct.
// Used for UsageFactor/UsageThreshold fields in QoSUpdate.
func ConvertFloat64PtrToFloat64NoValStruct(source *float64) *api.V0042Float64NoValStruct {
	if source == nil {
		return nil
	}
	setTrue := true
	return &api.V0042Float64NoValStruct{
		Set:    &setTrue,
		Number: source,
	}
}

// ConvertStringSliceToClusterRecFlags converts []string to API V0042ClusterRecFlags.
// Used for Flags in ClusterCreate.
func ConvertStringSliceToClusterRecFlags(source []string) *api.V0042ClusterRecFlags {
	if len(source) == 0 {
		return nil
	}
	flags := api.V0042ClusterRecFlags(source)
	return &flags
}

// ConvertStringSliceToQosFlags converts []string to API V0042QosFlags.
// Used for Flags in QoSCreate.
func ConvertStringSliceToQosFlags(source []string) *api.V0042QosFlags {
	if len(source) == 0 {
		return nil
	}
	flags := api.V0042QosFlags(source)
	return &flags
}

// ConvertStringSlicePtrToQosFlags converts *[]string to API V0042QosFlags.
// Used for Flags in QoSUpdate.
func ConvertStringSlicePtrToQosFlags(source *[]string) *api.V0042QosFlags {
	if source == nil || len(*source) == 0 {
		return nil
	}
	flags := api.V0042QosFlags(*source)
	return &flags
}

// ConvertNodeStateToNodeStates converts a *NodeState to *V0042NodeStates for API.
// Used for State and NextStateAfterReboot in NodeUpdate.
func ConvertNodeStateToNodeStates(source *types.NodeState) *api.V0042NodeStates {
	if source == nil {
		return nil
	}
	result := api.V0042NodeStates{string(*source)}
	return &result
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
// Node Write Helper Functions (NodeUpdate -> V0042UpdateNodeMsg)
// =============================================================================
// ConvertStringSliceToHostlistString converts []string to *V0042HostlistString.
// Used for Address, Hostname, Name fields.
func ConvertStringSliceToHostlistString(source []string) *api.V0042HostlistString {
	if len(source) == 0 {
		return nil
	}
	hl := api.V0042HostlistString(source)
	return &hl
}

// ConvertStringSliceToCsvString converts []string to *V0042CsvString.
// Used for Features, FeaturesAct fields.
func ConvertStringSliceToCsvString(source []string) *api.V0042CsvString {
	if len(source) == 0 {
		return nil
	}
	csv := api.V0042CsvString(source)
	return &csv
}

// ConvertNodeStatesToAPIV42 converts []NodeState to *V0042NodeStates ([]string).
func ConvertNodeStatesToAPIV42(source []types.NodeState) *api.V0042NodeStates {
	if len(source) == 0 {
		return nil
	}
	result := make([]string, len(source))
	for i, s := range source {
		result[i] = string(s)
	}
	states := api.V0042NodeStates(result)
	return &states
}

// ConvertUint32PtrToNoValStruct converts *uint32 to *V0042Uint32NoValStruct.
// Used for fields like Weight, ResumeAfter.
func ConvertUint32PtrToNoValStruct(source *uint32) *api.V0042Uint32NoValStruct {
	if source == nil {
		return nil
	}
	setTrue := true
	val := int32(*source)
	return &api.V0042Uint32NoValStruct{
		Set:    &setTrue,
		Number: &val,
	}
}

// =============================================================================
// Job Write Helper Functions (JobCreate -> V0042JobDescMsg)
// =============================================================================
// ConvertUint64PtrToNoValStruct converts *uint64 to *V0042Uint64NoValStruct.
// Used for fields like BeginTime, MemoryPerCPU, MemoryPerNode.
func ConvertUint64PtrToNoValStruct(source *uint64) *api.V0042Uint64NoValStruct {
	if source == nil {
		return nil
	}
	setTrue := true
	val := int64(*source)
	return &api.V0042Uint64NoValStruct{
		Set:    &setTrue,
		Number: &val,
	}
}

// ConvertUint16PtrToNoValStruct converts *uint16 to *V0042Uint16NoValStruct.
// Used for fields like DistributionPlaneSize, KillWarningDelay, SegmentSize.
func ConvertUint16PtrToNoValStruct(source *uint16) *api.V0042Uint16NoValStruct {
	if source == nil {
		return nil
	}
	setTrue := true
	val := int32(*source)
	return &api.V0042Uint16NoValStruct{
		Set:    &setTrue,
		Number: &val,
	}
}

// ConvertStringSliceToStringArray converts []string to *V0042StringArray.
// Used for Environment, Argv, SpankEnvironment.
func ConvertStringSliceToStringArray(source []string) *api.V0042StringArray {
	if len(source) == 0 {
		return nil
	}
	arr := api.V0042StringArray(source)
	return &arr
}

// ConvertMailTypeSliceToAPIV42 converts []MailTypeValue to *V0042JobMailFlags.
func ConvertMailTypeSliceToAPIV42(source []types.MailTypeValue) *api.V0042JobMailFlags {
	if len(source) == 0 {
		return nil
	}
	result := make([]string, len(source))
	for i, v := range source {
		result[i] = string(v)
	}
	mt := api.V0042JobMailFlags(result)
	return &mt
}

// ConvertFlagsSliceToAPIV42 converts []FlagsValue to *V0042JobFlags.
func ConvertFlagsSliceToAPIV42(source []types.FlagsValue) *api.V0042JobFlags {
	if len(source) == 0 {
		return nil
	}
	result := make([]string, len(source))
	for i, v := range source {
		result[i] = string(v)
	}
	flags := api.V0042JobFlags(result)
	return &flags
}

// ConvertCPUBindingFlagsSliceToAPIV42 converts []CPUBindingFlagsValue to *V0042CpuBindingFlags.
func ConvertCPUBindingFlagsSliceToAPIV42(source []types.CPUBindingFlagsValue) *api.V0042CpuBindingFlags {
	if len(source) == 0 {
		return nil
	}
	result := make([]string, len(source))
	for i, v := range source {
		result[i] = string(v)
	}
	flags := api.V0042CpuBindingFlags(result)
	return &flags
}

// ConvertKillWarningFlagsSliceToAPIV42 converts []KillWarningFlagsValue to *V0042WarnFlags.
func ConvertKillWarningFlagsSliceToAPIV42(source []types.KillWarningFlagsValue) *api.V0042WarnFlags {
	if len(source) == 0 {
		return nil
	}
	result := make([]string, len(source))
	for i, v := range source {
		result[i] = string(v)
	}
	flags := api.V0042WarnFlags(result)
	return &flags
}

// ConvertMemoryBindingTypeSliceToAPIV42 converts []MemoryBindingTypeValue to *V0042MemoryBindingType.
func ConvertMemoryBindingTypeSliceToAPIV42(source []types.MemoryBindingTypeValue) *api.V0042MemoryBindingType {
	if len(source) == 0 {
		return nil
	}
	result := make([]string, len(source))
	for i, v := range source {
		result[i] = string(v)
	}
	mbt := api.V0042MemoryBindingType(result)
	return &mbt
}

// ConvertOpenModeSliceToAPIV42 converts []OpenModeValue to *V0042OpenMode.
func ConvertOpenModeSliceToAPIV42(source []types.OpenModeValue) *api.V0042OpenMode {
	if len(source) == 0 {
		return nil
	}
	result := make([]string, len(source))
	for i, v := range source {
		result[i] = string(v)
	}
	om := api.V0042OpenMode(result)
	return &om
}

// ConvertProfileSliceToAPIV42 converts []ProfileValue to *V0042AcctGatherProfile.
func ConvertProfileSliceToAPIV42(source []types.ProfileValue) *api.V0042AcctGatherProfile {
	if len(source) == 0 {
		return nil
	}
	result := make([]string, len(source))
	for i, v := range source {
		result[i] = string(v)
	}
	jp := api.V0042AcctGatherProfile(result)
	return &jp
}

// ConvertSharedSliceToAPIV42 converts []SharedValue to *V0042JobShared.
func ConvertSharedSliceToAPIV42(source []types.SharedValue) *api.V0042JobShared {
	if len(source) == 0 {
		return nil
	}
	result := make([]string, len(source))
	for i, v := range source {
		result[i] = string(v)
	}
	js := api.V0042JobShared(result)
	return &js
}

// ConvertX11SliceToAPIV42 converts []X11Value to *V0042X11Flags.
func ConvertX11SliceToAPIV42(source []types.X11Value) *api.V0042X11Flags {
	if len(source) == 0 {
		return nil
	}
	result := make([]string, len(source))
	for i, v := range source {
		result[i] = string(v)
	}
	x11 := api.V0042X11Flags(result)
	return &x11
}

// ConvertCronEntryToAPIV42 converts *CronEntry to *V0042CronEntry.
func ConvertCronEntryToAPIV42(source *types.CronEntry) *api.V0042CronEntry {
	if source == nil {
		return nil
	}
	result := &api.V0042CronEntry{
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
		flags := make([]string, len(source.Flags))
		for i, f := range source.Flags {
			flags[i] = string(f)
		}
		cronFlags := api.V0042CronEntryFlags(flags)
		result.Flags = &cronFlags
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
