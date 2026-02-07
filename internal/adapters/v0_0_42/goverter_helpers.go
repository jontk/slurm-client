// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
// goverter_helpers.go provides extend functions for goverter converters.
// These functions handle complex type conversions that goverter cannot auto-generate.
package v0_0_42

import (
	"time"

	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_42"
)

// ConvertAssocShortSlice converts API AssocShort slice to common type.
// Used by goverter as an extend function.
func ConvertAssocShortSlice(source *api.V0042AssocShortList) []types.AssocShort {
	if source == nil {
		return nil
	}
	result := make([]types.AssocShort, len(*source))
	for i, assoc := range *source {
		result[i] = types.AssocShort{
			Account:   assoc.Account,
			Cluster:   assoc.Cluster,
			ID:        assoc.Id, // Note: Id in API, ID in common
			Partition: assoc.Partition,
			User:      assoc.User,
		}
	}
	return result
}

// ConvertCoordSlice converts API Coord slice to common Coord type.
// Used by goverter as an extend function.
func ConvertCoordSlice(source *api.V0042CoordList) []types.Coord {
	if source == nil {
		return nil
	}
	coords := make([]types.Coord, len(*source))
	for i, c := range *source {
		coords[i] = types.Coord{
			Name:   c.Name,
			Direct: c.Direct,
		}
	}
	return coords
}

// ConvertAccountFlags converts API AccountFlags slice to common AccountFlagsValue slice.
// Note: v0_0_42 uses []string for AccountFlags (not a typed enum).
func ConvertAccountFlags(source *api.V0042AccountFlags) []types.AccountFlagsValue {
	if source == nil {
		return nil
	}
	flags := make([]types.AccountFlagsValue, len(*source))
	for i, flag := range *source {
		flags[i] = types.AccountFlagsValue(flag)
	}
	return flags
}

// ConvertCoordNamesToSlice converts a []string (coordinator names) to API V0042CoordList.
// Used for AccountCreate where coordinators are provided as names only.
func ConvertCoordNamesToSlice(source []string) *api.V0042CoordList {
	if len(source) == 0 {
		return nil
	}
	coords := make(api.V0042CoordList, len(source))
	for i, name := range source {
		coords[i] = api.V0042Coord{Name: name}
	}
	return &coords
}

// =============================================================================
// NoValStruct Helpers - Generic converters for SLURM's NoValStruct pattern
// =============================================================================
// ConvertTimeNoVal converts a V0042Uint64NoValStruct to time.Time.
// Returns zero time if source is nil or number is 0.
func ConvertTimeNoVal(source *api.V0042Uint64NoValStruct) time.Time {
	if source == nil || source.Number == nil || *source.Number == 0 {
		return time.Time{}
	}
	return time.Unix(*source.Number, 0)
}

// ConvertUint64NoVal converts a V0042Uint64NoValStruct to *uint64.
// Returns nil if source is nil or Set is false.
func ConvertUint64NoVal(source *api.V0042Uint64NoValStruct) *uint64 {
	if source == nil || source.Set == nil || !*source.Set || source.Number == nil {
		return nil
	}
	val := uint64(*source.Number)
	return &val
}

// ConvertUint32NoVal converts a V0042Uint32NoValStruct to *uint32.
// Returns nil if source is nil or Set is false.
func ConvertUint32NoVal(source *api.V0042Uint32NoValStruct) *uint32 {
	if source == nil || source.Set == nil || !*source.Set || source.Number == nil {
		return nil
	}
	val := uint32(*source.Number)
	return &val
}

// =============================================================================
// Enum Slice Helpers
// =============================================================================
// ConvertNodeStateSlice converts API NodeState slice to common NodeState slice.
// Note: v0_0_42 uses []string for NodeState (not a typed enum).
func ConvertNodeStateSlice(source *[]string) []types.NodeState {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.NodeState, len(*source))
	for i, s := range *source {
		result[i] = types.NodeState(s)
	}
	return result
}

// ConvertNextStateAfterReboot converts API next state enum slice to common NodeState slice.
// Note: v0_0_42 uses []string for NextStateAfterReboot (not a typed enum).
func ConvertNextStateAfterReboot(source *[]string) []types.NodeState {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.NodeState, len(*source))
	for i, s := range *source {
		result[i] = types.NodeState(s)
	}
	return result
}

// =============================================================================
// Node-Specific Helpers
// =============================================================================
// ConvertNodeEnergyGoverter converts API node energy to common type.
// This is an alias for the existing convertNodeEnergy function, exported for goverter.
func ConvertNodeEnergyGoverter(source *api.V0042AcctGatherEnergy) *types.NodeEnergy {
	return convertNodeEnergy(source)
}

// ConvertResumeAfterGoverter converts resume after time.
// This is an alias for the existing convertResumeAfter function, exported for goverter.
func ConvertResumeAfterGoverter(source *api.V0042Uint64NoValStruct) *uint64 {
	return convertResumeAfter(source)
}

// ConvertCertFlagsGoverter converts API CertFlags slice to common type.
// Note: v0_0_42 uses []string for CertFlags (not a typed enum).
func ConvertCertFlagsGoverter(source *[]string) []types.CertFlagsValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.CertFlagsValue, len(*source))
	for i, flag := range *source {
		result[i] = types.CertFlagsValue(flag)
	}
	return result
}

// ConvertExternalSensors converts API ExternalSensors to common type.
func ConvertExternalSensors(source *map[string]interface{}) map[string]interface{} {
	if source == nil {
		return nil
	}
	return *source
}

// ConvertPower converts API Power to common type.
func ConvertPower(source *map[string]interface{}) map[string]interface{} {
	if source == nil {
		return nil
	}
	return *source
}

// =============================================================================
// String/Slice Helpers
// =============================================================================
// ConvertCSVStringToSlice converts a V0042CsvString pointer to a []string.
// V0042CsvString is already []string, so this just dereferences the pointer.
func ConvertCSVStringToSlice(source *api.V0042CsvString) []string {
	if source == nil {
		return nil
	}
	return *source
}

// ConvertStringSliceToCSV converts a []string to a *string (CSV).
func ConvertStringSliceToCSV(source []string) *string {
	if len(source) == 0 {
		return nil
	}
	result := source[0]
	for _, s := range source[1:] {
		result += "," + s
	}
	return &result
}

// =============================================================================
// QoS Helpers
// =============================================================================
// ConvertQoSFlags converts API QoS flags slice to common type.
// Note: v0_0_42 uses []string for QoS flags (not a typed enum).
func ConvertQoSFlags(source *[]string) []types.QoSFlagsValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.QoSFlagsValue, len(*source))
	for i, flag := range *source {
		result[i] = types.QoSFlagsValue(flag)
	}
	return result
}

// ConvertFloat64NoVal converts a V0042Float64NoValStruct to *float64.
// Returns nil if source is nil or Set is false.
func ConvertFloat64NoVal(source *api.V0042Float64NoValStruct) *float64 {
	if source == nil || source.Set == nil || !*source.Set || source.Number == nil {
		return nil
	}
	val := *source.Number
	return &val
}

// =============================================================================
// Job Helpers
// =============================================================================
// ConvertJobStateSlice converts API JobState slice to common JobState slice.
// Note: v0_0_42 uses []string for JobState (via V0042JobState = []string).
func ConvertJobStateSlice(source *api.V0042JobState) []types.JobState {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.JobState, len(*source))
	for i, s := range *source {
		result[i] = types.JobState(s)
	}
	return result
}

// ConvertUint16NoVal converts a V0042Uint16NoValStruct to *uint16.
// Returns nil if source is nil or Set is false.
func ConvertUint16NoVal(source *api.V0042Uint16NoValStruct) *uint16 {
	if source == nil || source.Set == nil || !*source.Set || source.Number == nil {
		return nil
	}
	val := uint16(*source.Number)
	return &val
}

// =============================================================================
// User Helpers
// =============================================================================
// ConvertAdminLevelSlice converts API AdminLevel slice to common type.
// Note: v0_0_42 uses []string for AdminLevel (V0042AdminLvl = []string).
func ConvertAdminLevelSlice(source *api.V0042AdminLvl) []types.AdministratorLevelValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.AdministratorLevelValue, len(*source))
	for i, level := range *source {
		result[i] = types.AdministratorLevelValue(level)
	}
	return result
}

// ConvertWckeySlice converts API Wckey slice to common type.
func ConvertWckeySlice(source *api.V0042WckeyList) []types.WCKey {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.WCKey, len(*source))
	for i, w := range *source {
		result[i] = types.WCKey{
			Cluster: w.Cluster,
			ID:      w.Id,
			Name:    w.Name,
			User:    w.User,
		}
	}
	return result
}

// ConvertUserFlags converts API User flags slice to common type.
// Note: v0_0_42 uses []string for User flags (not a typed enum).
func ConvertUserFlags(source *[]string) []types.UserDefaultFlagsValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.UserDefaultFlagsValue, len(*source))
	for i, flag := range *source {
		result[i] = types.UserDefaultFlagsValue(flag)
	}
	return result
}

// =============================================================================
// Association Helpers
// =============================================================================
// ConvertAssocFlags converts API Association flags slice to common type.
// Note: v0_0_42 uses []string for Association flags (V0042AssocFlags = []string).
func ConvertAssocFlags(source *api.V0042AssocFlags) []types.AssociationDefaultFlagsValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.AssociationDefaultFlagsValue, len(*source))
	for i, flag := range *source {
		result[i] = types.AssociationDefaultFlagsValue(flag)
	}
	return result
}

// ConvertQosStringIdList converts API QosStringIdList to common []string.
// V0042QosStringIdList is already []string, so this just dereferences the pointer.
func ConvertQosStringIdList(source *api.V0042QosStringIdList) []string {
	if source == nil || len(*source) == 0 {
		return nil
	}
	return *source
}

// =============================================================================
// Cluster Helpers
// =============================================================================
// ConvertClusterFlags converts API Cluster flags slice to common type.
// Note: v0_0_42 uses []string for Cluster flags (V0042ClusterRecFlags = []string).
func ConvertClusterFlags(source *api.V0042ClusterRecFlags) []types.ClusterControllerFlagsValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.ClusterControllerFlagsValue, len(*source))
	for i, flag := range *source {
		result[i] = types.ClusterControllerFlagsValue(flag)
	}
	return result
}

// ConvertClusterTRES converts API Cluster TRES to common type.
func ConvertClusterTRES(source *api.V0042TresList) []types.TRES {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.TRES, len(*source))
	for i, t := range *source {
		result[i] = types.TRES{
			Type:  t.Type,
			Name:  t.Name,
			ID:    t.Id,
			Count: t.Count,
		}
	}
	return result
}

// =============================================================================
// WCKey Helpers
// =============================================================================
// ConvertWCKeyFlags converts API WCKey flags slice to common type.
// Note: v0_0_42 uses []string for WCKey flags (not a typed enum).
func ConvertWCKeyFlags(source *[]string) []types.WCKeyFlagsValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.WCKeyFlagsValue, len(*source))
	for i, flag := range *source {
		result[i] = types.WCKeyFlagsValue(flag)
	}
	return result
}

// =============================================================================
// Reservation Helpers
// =============================================================================
// ConvertReservationFlags converts API Reservation flags slice to common type.
// Note: v0_0_42 uses []string for Reservation flags (not a typed enum).
func ConvertReservationFlags(source *[]string) []types.ReservationFlagsValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.ReservationFlagsValue, len(*source))
	for i, flag := range *source {
		result[i] = types.ReservationFlagsValue(flag)
	}
	return result
}

// =============================================================================
// Job-Specific Helpers
// =============================================================================
// Note: v0_0_42 uses []string for most of these fields (not typed enums).

// ConvertJobFlags converts API JobFlags slice to common FlagsValue slice.
// Used by goverter as an extend function.
func ConvertJobFlags(source *api.V0042JobFlags) []types.FlagsValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	flags := make([]types.FlagsValue, len(*source))
	for i, flag := range *source {
		flags[i] = types.FlagsValue(flag)
	}
	return flags
}

// ConvertJobMailType converts API JobMailFlags slice to common MailTypeValue slice.
// Used by goverter as an extend function.
func ConvertJobMailType(source *api.V0042JobMailFlags) []types.MailTypeValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.MailTypeValue, len(*source))
	for i, mt := range *source {
		result[i] = types.MailTypeValue(mt)
	}
	return result
}

// ConvertJobProfile converts API AcctGatherProfile slice to common ProfileValue slice.
// Used by goverter as an extend function.
func ConvertJobProfile(source *api.V0042AcctGatherProfile) []types.ProfileValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.ProfileValue, len(*source))
	for i, p := range *source {
		result[i] = types.ProfileValue(p)
	}
	return result
}

// ConvertJobShared converts API JobShared slice to common SharedValue slice.
// Used by goverter as an extend function.
func ConvertJobShared(source *api.V0042JobShared) []types.SharedValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.SharedValue, len(*source))
	for i, s := range *source {
		result[i] = types.SharedValue(s)
	}
	return result
}

// ConvertExitCode converts API ProcessExitCodeVerbose to common ExitCode.
// Used by goverter as an extend function.
func ConvertExitCode(source *api.V0042ProcessExitCodeVerbose) *types.ExitCode {
	if source == nil {
		return nil
	}
	result := &types.ExitCode{}
	// Convert return code
	if source.ReturnCode != nil && source.ReturnCode.Number != nil {
		rc := uint32(*source.ReturnCode.Number)
		result.ReturnCode = &rc
	}
	// Convert signal
	if source.Signal != nil {
		result.Signal = &types.ExitCodeSignal{}
		if source.Signal.Id != nil && source.Signal.Id.Number != nil {
			sigID := uint16(*source.Signal.Id.Number)
			result.Signal.ID = &sigID
		}
		if source.Signal.Name != nil {
			result.Signal.Name = source.Signal.Name
		}
	}
	// Convert status
	if source.Status != nil {
		for _, s := range *source.Status {
			result.Status = append(result.Status, types.StatusValue(s))
		}
	}
	return result
}

// =============================================================================
// WCKey Accounting Helper
// =============================================================================

// ConvertAccounting converts API AccountingList to common Accounting slice.
// Used by goverter as an extend function.
func ConvertAccounting(source *api.V0042AccountingList) []types.Accounting {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.Accounting, len(*source))
	for i, acct := range *source {
		accounting := types.Accounting{
			ID:    acct.Id,
			IDAlt: acct.IdAlt,
			Start: acct.Start,
		}
		// Convert TRES
		if acct.TRES != nil {
			accounting.TRES = &types.TRES{
				Count: acct.TRES.Count,
				ID:    acct.TRES.Id,
				Name:  acct.TRES.Name,
				Type:  acct.TRES.Type,
			}
		}
		// Convert Allocated
		if acct.Allocated != nil {
			accounting.Allocated = &types.AccountingAllocated{
				Seconds: acct.Allocated.Seconds,
			}
		}
		result[i] = accounting
	}
	return result
}

// =============================================================================
// Reservation Read Helpers
// =============================================================================

// ConvertReservationFlagsRead converts API ReservationFlags slice to common ReservationFlagsValue slice.
// Note: v0_0_42 uses []string for Reservation flags.
// Used by goverter as an extend function for read conversions.
func ConvertReservationFlagsRead(source *api.V0042ReservationFlags) []types.ReservationFlagsValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.ReservationFlagsValue, len(*source))
	for i, flag := range *source {
		result[i] = types.ReservationFlagsValue(flag)
	}
	return result
}

// ConvertReservationCoreSpec converts API ReservationCoreSpec slice to common ReservationCoreSpec slice.
// Used by goverter as an extend function.
func ConvertReservationCoreSpec(source *api.V0042ReservationInfoCoreSpec) []types.ReservationCoreSpec {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.ReservationCoreSpec, len(*source))
	for i, spec := range *source {
		result[i] = types.ReservationCoreSpec{
			Core: spec.Core,
			Node: spec.Node,
		}
	}
	return result
}

// ConvertReservationPurgeCompleted converts API purge completed struct to common type.
// Used by goverter as an extend function.
func ConvertReservationPurgeCompleted(source *struct {
	Time *api.V0042Uint32NoValStruct `json:"time,omitempty"`
}) *types.ReservationPurgeCompleted {
	if source == nil {
		return nil
	}
	result := &types.ReservationPurgeCompleted{}
	if source.Time != nil {
		result.Time = ConvertUint32NoVal(source.Time)
	}
	return result
}

// =============================================================================
// Association Helpers
// =============================================================================

// ConvertAssociationDefault converts API AssocDefault to common AssociationDefault.
// Used by goverter as an extend function.
func ConvertAssociationDefault(source *struct {
	Qos *string `json:"qos,omitempty"`
}) *types.AssociationDefault {
	if source == nil {
		return nil
	}
	return &types.AssociationDefault{
		QoS: source.Qos,
	}
}

// =============================================================================
// User Helpers (Additional)
// =============================================================================

// ConvertUserDefault converts API UserDefault to common UserDefault.
// Note: v0_0_42 User Default does not have QoS field (only Account and Wckey).
// Used by goverter as an extend function.
func ConvertUserDefault(source *struct {
	Account *string `json:"account,omitempty"`
	Wckey   *string `json:"wckey,omitempty"`
}) *types.UserDefault {
	if source == nil {
		return nil
	}
	return &types.UserDefault{
		Account: source.Account,
		Wckey:   source.Wckey,
		// QoS is not available in v0_0_42
	}
}

// =============================================================================
// Cluster Helpers (Additional)
// =============================================================================

// ConvertClusterController converts API ClusterController to common ClusterController.
// Used by goverter as an extend function.
func ConvertClusterController(source *struct {
	Host *string `json:"host,omitempty"`
	Port *int32  `json:"port,omitempty"`
}) *types.ClusterController {
	if source == nil {
		return nil
	}
	return &types.ClusterController{
		Host: source.Host,
		Port: source.Port,
	}
}

// ConvertClusterAssociations converts API ClusterAssociations to common ClusterAssociations.
// Used by goverter as an extend function.
func ConvertClusterAssociations(source *struct {
	Root *api.V0042AssocShort `json:"root,omitempty"`
}) *types.ClusterAssociations {
	if source == nil {
		return nil
	}
	result := &types.ClusterAssociations{}
	if source.Root != nil {
		result.Root = ConvertAssocShort(source.Root)
	}
	return result
}

// ConvertAssocShort converts API AssocShort to common AssocShort.
func ConvertAssocShort(source *api.V0042AssocShort) *types.AssocShort {
	if source == nil {
		return nil
	}
	return &types.AssocShort{
		Account:   source.Account,
		Cluster:   source.Cluster,
		Partition: source.Partition,
		User:      source.User,
	}
}

// =============================================================================
// QoS Helpers (Additional)
// =============================================================================

// ConvertQoSPreempt converts API QosPreempt to common QoSPreempt.
// Note: v0_0_42 uses []string for preempt modes (V0042QosPreemptModes).
// Used by goverter as an extend function.
func ConvertQoSPreempt(source *struct {
	ExemptTime *api.V0042Uint32NoValStruct `json:"exempt_time,omitempty"`
	List       *api.V0042QosPreemptList    `json:"list,omitempty"`
	Mode       *api.V0042QosPreemptModes   `json:"mode,omitempty"`
}) *types.QoSPreempt {
	if source == nil {
		return nil
	}
	result := &types.QoSPreempt{}
	if source.ExemptTime != nil {
		result.ExemptTime = ConvertUint32NoVal(source.ExemptTime)
	}
	if source.List != nil {
		result.List = *source.List
	}
	if source.Mode != nil {
		modes := make([]types.ModeValue, len(*source.Mode))
		for i, m := range *source.Mode {
			modes[i] = types.ModeValue(m)
		}
		result.Mode = modes
	}
	return result
}

// =============================================================================
// Job Helpers
// =============================================================================

// ConvertJobPower converts API JobInfo Power to common JobPower.
// Used by goverter as an extend function.
func ConvertJobPower(source *struct {
	Flags *[]interface{} `json:"flags,omitempty"`
}) *types.JobPower {
	if source == nil {
		return nil
	}
	result := &types.JobPower{}
	if source.Flags != nil {
		result.Flags = *source.Flags
	}
	return result
}

// ConvertJobPriorityByPartition converts API PriorityByPartition to common JobPartitionPriority slice.
// Used by goverter as an extend function.
func ConvertJobPriorityByPartition(source *api.V0042PriorityByPartition) []types.JobPartitionPriority {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.JobPartitionPriority, len(*source))
	for i, pp := range *source {
		result[i] = types.JobPartitionPriority{
			Partition: pp.Partition,
			Priority:  pp.Priority,
		}
	}
	return result
}

// ConvertJobGRESDetail converts API JobInfoGresDetail to common []string.
// Used by goverter as an extend function.
func ConvertJobGRESDetail(source *api.V0042JobInfoGresDetail) []string {
	if source == nil {
		return nil
	}
	return *source
}

// =============================================================================
// Partition Helpers
// =============================================================================

// ConvertPartitionAccounts converts API PartitionInfo Accounts to common PartitionAccounts.
// Used by goverter as an extend function.
func ConvertPartitionAccounts(source *struct {
	Allowed *string `json:"allowed,omitempty"`
	Deny    *string `json:"deny,omitempty"`
}) *types.PartitionAccounts {
	if source == nil {
		return nil
	}
	return &types.PartitionAccounts{
		Allowed: source.Allowed,
		Deny:    source.Deny,
	}
}

// ConvertPartitionCPUs converts API PartitionInfo CPUs to common PartitionCPUs.
// Used by goverter as an extend function.
func ConvertPartitionCPUs(source *struct {
	TaskBinding *int32 `json:"task_binding,omitempty"`
	Total       *int32 `json:"total,omitempty"`
}) *types.PartitionCPUs {
	if source == nil {
		return nil
	}
	return &types.PartitionCPUs{
		TaskBinding: source.TaskBinding,
		Total:       source.Total,
	}
}

// ConvertPartitionDefaults converts API PartitionInfo Defaults to common PartitionDefaults.
// Used by goverter as an extend function.
func ConvertPartitionDefaults(source *struct {
	Job                    *string                     `json:"job,omitempty"`
	MemoryPerCpu           *int64                      `json:"memory_per_cpu,omitempty"`
	PartitionMemoryPerCpu  *api.V0042Uint64NoValStruct `json:"partition_memory_per_cpu,omitempty"`
	PartitionMemoryPerNode *api.V0042Uint64NoValStruct `json:"partition_memory_per_node,omitempty"`
	Time                   *api.V0042Uint32NoValStruct `json:"time,omitempty"`
}) *types.PartitionDefaults {
	if source == nil {
		return nil
	}
	result := &types.PartitionDefaults{
		Job:          source.Job,
		MemoryPerCPU: source.MemoryPerCpu,
	}
	if source.PartitionMemoryPerCpu != nil {
		result.PartitionMemoryPerCPU = ConvertUint64NoVal(source.PartitionMemoryPerCpu)
	}
	if source.PartitionMemoryPerNode != nil {
		result.PartitionMemoryPerNode = ConvertUint64NoVal(source.PartitionMemoryPerNode)
	}
	if source.Time != nil {
		result.Time = ConvertUint32NoVal(source.Time)
	}
	return result
}

// ConvertPartitionGroups converts API PartitionInfo Groups to common PartitionGroups.
// Used by goverter as an extend function.
func ConvertPartitionGroups(source *struct {
	Allowed *string `json:"allowed,omitempty"`
}) *types.PartitionGroups {
	if source == nil {
		return nil
	}
	return &types.PartitionGroups{
		Allowed: source.Allowed,
	}
}

// ConvertPartitionMaximums converts API PartitionInfo Maximums to common PartitionMaximums.
// Used by goverter as an extend function.
func ConvertPartitionMaximums(source *struct {
	CpusPerNode            *api.V0042Uint32NoValStruct `json:"cpus_per_node,omitempty"`
	CpusPerSocket          *api.V0042Uint32NoValStruct `json:"cpus_per_socket,omitempty"`
	MemoryPerCpu           *int64                      `json:"memory_per_cpu,omitempty"`
	Nodes                  *api.V0042Uint32NoValStruct `json:"nodes,omitempty"`
	OverTimeLimit          *api.V0042Uint16NoValStruct `json:"over_time_limit,omitempty"`
	Oversubscribe          *struct {
		Flags *api.V0042OversubscribeFlags `json:"flags,omitempty"`
		Jobs  *int32                       `json:"jobs,omitempty"`
	} `json:"oversubscribe,omitempty"`
	PartitionMemoryPerCpu  *api.V0042Uint64NoValStruct `json:"partition_memory_per_cpu,omitempty"`
	PartitionMemoryPerNode *api.V0042Uint64NoValStruct `json:"partition_memory_per_node,omitempty"`
	Shares                 *int32                      `json:"shares,omitempty"`
	Time                   *api.V0042Uint32NoValStruct `json:"time,omitempty"`
}) *types.PartitionMaximums {
	if source == nil {
		return nil
	}
	result := &types.PartitionMaximums{
		MemoryPerCPU: source.MemoryPerCpu,
		Shares:       source.Shares,
	}
	if source.CpusPerNode != nil {
		result.CPUsPerNode = ConvertUint32NoVal(source.CpusPerNode)
	}
	if source.CpusPerSocket != nil {
		result.CPUsPerSocket = ConvertUint32NoVal(source.CpusPerSocket)
	}
	if source.Nodes != nil {
		result.Nodes = ConvertUint32NoVal(source.Nodes)
	}
	if source.OverTimeLimit != nil {
		result.OverTimeLimit = ConvertUint16NoVal(source.OverTimeLimit)
	}
	if source.Oversubscribe != nil {
		result.Oversubscribe = &types.PartitionMaximumsOversubscribe{
			Jobs: source.Oversubscribe.Jobs,
		}
		if source.Oversubscribe.Flags != nil {
			flags := make([]types.PartitionMaximumsOversubscribeFlagsValue, len(*source.Oversubscribe.Flags))
			for i, f := range *source.Oversubscribe.Flags {
				flags[i] = types.PartitionMaximumsOversubscribeFlagsValue(string(f))
			}
			result.Oversubscribe.Flags = flags
		}
	}
	if source.PartitionMemoryPerCpu != nil {
		result.PartitionMemoryPerCPU = ConvertUint64NoVal(source.PartitionMemoryPerCpu)
	}
	if source.PartitionMemoryPerNode != nil {
		result.PartitionMemoryPerNode = ConvertUint64NoVal(source.PartitionMemoryPerNode)
	}
	if source.Time != nil {
		result.Time = ConvertUint32NoVal(source.Time)
	}
	return result
}

// ConvertPartitionMinimums converts API PartitionInfo Minimums to common PartitionMinimums.
// Used by goverter as an extend function.
func ConvertPartitionMinimums(source *struct {
	Nodes *int32 `json:"nodes,omitempty"`
}) *types.PartitionMinimums {
	if source == nil {
		return nil
	}
	return &types.PartitionMinimums{
		Nodes: source.Nodes,
	}
}

// ConvertPartitionNodes converts API PartitionInfo Nodes to common PartitionNodes.
// Used by goverter as an extend function.
func ConvertPartitionNodes(source *struct {
	AllowedAllocation *string `json:"allowed_allocation,omitempty"`
	Configured        *string `json:"configured,omitempty"`
	Total             *int32  `json:"total,omitempty"`
}) *types.PartitionNodes {
	if source == nil {
		return nil
	}
	return &types.PartitionNodes{
		AllowedAllocation: source.AllowedAllocation,
		Configured:        source.Configured,
		Total:             source.Total,
	}
}

// ConvertPartitionPartition converts API PartitionInfo Partition to common PartitionPartition.
// Used by goverter as an extend function.
// Note: v0_0_42 uses V0042PartitionStates which is []string.
func ConvertPartitionPartition(source *struct {
	State *api.V0042PartitionStates `json:"state,omitempty"`
}) *types.PartitionPartition {
	if source == nil {
		return nil
	}
	result := &types.PartitionPartition{}
	if source.State != nil {
		states := make([]types.StateValue, len(*source.State))
		for i, s := range *source.State {
			states[i] = types.StateValue(s)
		}
		result.State = states
	}
	return result
}

// ConvertPartitionPriority converts API PartitionInfo Priority to common PartitionPriority.
// Used by goverter as an extend function.
func ConvertPartitionPriority(source *struct {
	JobFactor *int32 `json:"job_factor,omitempty"`
	Tier      *int32 `json:"tier,omitempty"`
}) *types.PartitionPriority {
	if source == nil {
		return nil
	}
	return &types.PartitionPriority{
		JobFactor: source.JobFactor,
		Tier:      source.Tier,
	}
}

// ConvertPartitionQoS converts API PartitionInfo QoS to common PartitionQoS.
// Used by goverter as an extend function.
func ConvertPartitionQoS(source *struct {
	Allowed  *string `json:"allowed,omitempty"`
	Assigned *string `json:"assigned,omitempty"`
	Deny     *string `json:"deny,omitempty"`
}) *types.PartitionQoS {
	if source == nil {
		return nil
	}
	return &types.PartitionQoS{
		Allowed:  source.Allowed,
		Assigned: source.Assigned,
		Deny:     source.Deny,
	}
}

// ConvertPartitionSelectType converts API PartitionInfo SelectType to common SelectTypeValue slice.
// Used by goverter as an extend function.
// Note: v0_0_42 uses V0042CrType which is []string.
func ConvertPartitionSelectType(source *api.V0042CrType) []types.SelectTypeValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.SelectTypeValue, len(*source))
	for i, s := range *source {
		result[i] = types.SelectTypeValue(s)
	}
	return result
}

// ConvertPartitionSuspendTime converts API PartitionInfo SuspendTime to time.Time.
// Used by goverter as an extend function.
func ConvertPartitionSuspendTime(source *api.V0042Uint32NoValStruct) time.Time {
	if source == nil || source.Number == nil || *source.Number == 0 {
		return time.Time{}
	}
	return time.Unix(int64(*source.Number), 0)
}

// ConvertPartitionTimeouts converts API PartitionInfo Timeouts to common PartitionTimeouts.
// Used by goverter as an extend function.
func ConvertPartitionTimeouts(source *struct {
	Resume  *api.V0042Uint16NoValStruct `json:"resume,omitempty"`
	Suspend *api.V0042Uint16NoValStruct `json:"suspend,omitempty"`
}) *types.PartitionTimeouts {
	if source == nil {
		return nil
	}
	result := &types.PartitionTimeouts{}
	if source.Resume != nil {
		result.Resume = ConvertUint16NoVal(source.Resume)
	}
	if source.Suspend != nil {
		result.Suspend = ConvertUint16NoVal(source.Suspend)
	}
	return result
}

// ConvertPartitionTRES converts API PartitionInfo TRES to common PartitionTRES.
// Used by goverter as an extend function.
func ConvertPartitionTRES(source *struct {
	BillingWeights *string `json:"billing_weights,omitempty"`
	Configured     *string `json:"configured,omitempty"`
}) *types.PartitionTRES {
	if source == nil {
		return nil
	}
	return &types.PartitionTRES{
		BillingWeights: source.BillingWeights,
		Configured:     source.Configured,
	}
}

// =============================================================================
// Association Max/Min Helpers
// =============================================================================

// convertTresList is a helper to convert V0042TresList to []types.TRES.
func convertTresList(source *api.V0042TresList) []types.TRES {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.TRES, len(*source))
	for i, tres := range *source {
		result[i] = types.TRES{
			Count: tres.Count,
			ID:    tres.Id,
			Name:  tres.Name,
			Type:  tres.Type,
		}
	}
	return result
}

// ConvertAssociationMax converts API Assoc Max to common AssociationMax.
// Used by goverter as an extend function.
func ConvertAssociationMax(source *struct {
	Jobs *struct {
		Accruing *api.V0042Uint32NoValStruct `json:"accruing,omitempty"`
		Active   *api.V0042Uint32NoValStruct `json:"active,omitempty"`
		Per      *struct {
			Accruing  *api.V0042Uint32NoValStruct `json:"accruing,omitempty"`
			Count     *api.V0042Uint32NoValStruct `json:"count,omitempty"`
			Submitted *api.V0042Uint32NoValStruct `json:"submitted,omitempty"`
			WallClock *api.V0042Uint32NoValStruct `json:"wall_clock,omitempty"`
		} `json:"per,omitempty"`
		Total *api.V0042Uint32NoValStruct `json:"total,omitempty"`
	} `json:"jobs,omitempty"`
	Per *struct {
		Account *struct {
			WallClock *api.V0042Uint32NoValStruct `json:"wall_clock,omitempty"`
		} `json:"account,omitempty"`
	} `json:"per,omitempty"`
	Tres *struct {
		Group *struct {
			Active  *api.V0042TresList `json:"active,omitempty"`
			Minutes *api.V0042TresList `json:"minutes,omitempty"`
		} `json:"group,omitempty"`
		Minutes *struct {
			Per *struct {
				Job *api.V0042TresList `json:"job,omitempty"`
			} `json:"per,omitempty"`
			Total *api.V0042TresList `json:"total,omitempty"`
		} `json:"minutes,omitempty"`
		Per *struct {
			Job  *api.V0042TresList `json:"job,omitempty"`
			Node *api.V0042TresList `json:"node,omitempty"`
		} `json:"per,omitempty"`
		Total *api.V0042TresList `json:"total,omitempty"`
	} `json:"tres,omitempty"`
}) *types.AssociationMax {
	if source == nil {
		return nil
	}
	result := &types.AssociationMax{}

	// Convert Jobs
	if source.Jobs != nil {
		result.Jobs = &types.AssociationMaxJobs{
			Accruing: ConvertUint32NoVal(source.Jobs.Accruing),
			Active:   ConvertUint32NoVal(source.Jobs.Active),
			Total:    ConvertUint32NoVal(source.Jobs.Total),
		}
		if source.Jobs.Per != nil {
			result.Jobs.Per = &types.AssociationMaxJobsPer{
				Accruing:  ConvertUint32NoVal(source.Jobs.Per.Accruing),
				Count:     ConvertUint32NoVal(source.Jobs.Per.Count),
				Submitted: ConvertUint32NoVal(source.Jobs.Per.Submitted),
				WallClock: ConvertUint32NoVal(source.Jobs.Per.WallClock),
			}
		}
	}

	// Convert Per
	if source.Per != nil && source.Per.Account != nil {
		result.Per = &types.AssociationMaxPer{
			Account: &types.AssociationMaxPerAccount{
				WallClock: ConvertUint32NoVal(source.Per.Account.WallClock),
			},
		}
	}

	// Convert TRES
	if source.Tres != nil {
		result.TRES = &types.AssociationMaxTRES{
			Total: convertTresList(source.Tres.Total),
		}
		if source.Tres.Group != nil {
			result.TRES.Group = &types.AssociationMaxTRESGroup{
				Active:  convertTresList(source.Tres.Group.Active),
				Minutes: convertTresList(source.Tres.Group.Minutes),
			}
		}
		if source.Tres.Minutes != nil {
			result.TRES.Minutes = &types.AssociationMaxTRESMinutes{
				Total: convertTresList(source.Tres.Minutes.Total),
			}
			if source.Tres.Minutes.Per != nil {
				result.TRES.Minutes.Per = &types.AssociationMaxTRESMinutesPer{
					Job: convertTresList(source.Tres.Minutes.Per.Job),
				}
			}
		}
		if source.Tres.Per != nil {
			result.TRES.Per = &types.AssociationMaxTRESPer{
				Job:  convertTresList(source.Tres.Per.Job),
				Node: convertTresList(source.Tres.Per.Node),
			}
		}
	}

	return result
}

// ConvertAssociationMin converts API Assoc Min to common AssociationMin.
// Used by goverter as an extend function.
func ConvertAssociationMin(source *struct {
	PriorityThreshold *api.V0042Uint32NoValStruct `json:"priority_threshold,omitempty"`
}) *types.AssociationMin {
	if source == nil {
		return nil
	}
	return &types.AssociationMin{
		PriorityThreshold: ConvertUint32NoVal(source.PriorityThreshold),
	}
}

// =============================================================================
// QoS Limits Helpers
// =============================================================================

// ConvertQoSLimits converts API Qos Limits to common QoSLimits.
// Used by goverter as an extend function.
func ConvertQoSLimits(source *struct {
	Factor    *api.V0042Float64NoValStruct `json:"factor,omitempty"`
	GraceTime *int32                       `json:"grace_time,omitempty"`
	Max       *struct {
		Accruing *struct {
			Per *struct {
				Account *api.V0042Uint32NoValStruct `json:"account,omitempty"`
				User    *api.V0042Uint32NoValStruct `json:"user,omitempty"`
			} `json:"per,omitempty"`
		} `json:"accruing,omitempty"`
		ActiveJobs *struct {
			Accruing *api.V0042Uint32NoValStruct `json:"accruing,omitempty"`
			Count    *api.V0042Uint32NoValStruct `json:"count,omitempty"`
		} `json:"active_jobs,omitempty"`
		Jobs *struct {
			ActiveJobs *struct {
				Per *struct {
					Account *api.V0042Uint32NoValStruct `json:"account,omitempty"`
					User    *api.V0042Uint32NoValStruct `json:"user,omitempty"`
				} `json:"per,omitempty"`
			} `json:"active_jobs,omitempty"`
			Count *api.V0042Uint32NoValStruct `json:"count,omitempty"`
			Per   *struct {
				Account *api.V0042Uint32NoValStruct `json:"account,omitempty"`
				User    *api.V0042Uint32NoValStruct `json:"user,omitempty"`
			} `json:"per,omitempty"`
		} `json:"jobs,omitempty"`
		Tres *struct {
			Minutes *struct {
				Per *struct {
					Account *api.V0042TresList `json:"account,omitempty"`
					Job     *api.V0042TresList `json:"job,omitempty"`
					Qos     *api.V0042TresList `json:"qos,omitempty"`
					User    *api.V0042TresList `json:"user,omitempty"`
				} `json:"per,omitempty"`
				Total *api.V0042TresList `json:"total,omitempty"`
			} `json:"minutes,omitempty"`
			Per *struct {
				Account *api.V0042TresList `json:"account,omitempty"`
				Job     *api.V0042TresList `json:"job,omitempty"`
				Node    *api.V0042TresList `json:"node,omitempty"`
				User    *api.V0042TresList `json:"user,omitempty"`
			} `json:"per,omitempty"`
			Total *api.V0042TresList `json:"total,omitempty"`
		} `json:"tres,omitempty"`
		WallClock *struct {
			Per *struct {
				Job *api.V0042Uint32NoValStruct `json:"job,omitempty"`
				Qos *api.V0042Uint32NoValStruct `json:"qos,omitempty"`
			} `json:"per,omitempty"`
		} `json:"wall_clock,omitempty"`
	} `json:"max,omitempty"`
	Min *struct {
		PriorityThreshold *api.V0042Uint32NoValStruct `json:"priority_threshold,omitempty"`
		Tres              *struct {
			Per *struct {
				Job *api.V0042TresList `json:"job,omitempty"`
			} `json:"per,omitempty"`
		} `json:"tres,omitempty"`
	} `json:"min,omitempty"`
}) *types.QoSLimits {
	if source == nil {
		return nil
	}
	result := &types.QoSLimits{
		Factor:    ConvertFloat64NoVal(source.Factor),
		GraceTime: source.GraceTime,
	}

	// Convert Max
	if source.Max != nil {
		result.Max = &types.QoSLimitsMax{}

		// Max.Accruing
		if source.Max.Accruing != nil && source.Max.Accruing.Per != nil {
			result.Max.Accruing = &types.QoSLimitsMaxAccruing{
				Per: &types.QoSLimitsMaxAccruingPer{
					Account: ConvertUint32NoVal(source.Max.Accruing.Per.Account),
					User:    ConvertUint32NoVal(source.Max.Accruing.Per.User),
				},
			}
		}

		// Max.ActiveJobs
		if source.Max.ActiveJobs != nil {
			result.Max.ActiveJobs = &types.QoSLimitsMaxActiveJobs{
				Accruing: ConvertUint32NoVal(source.Max.ActiveJobs.Accruing),
				Count:    ConvertUint32NoVal(source.Max.ActiveJobs.Count),
			}
		}

		// Max.Jobs
		if source.Max.Jobs != nil {
			result.Max.Jobs = &types.QoSLimitsMaxJobs{
				Count: ConvertUint32NoVal(source.Max.Jobs.Count),
			}
			if source.Max.Jobs.ActiveJobs != nil && source.Max.Jobs.ActiveJobs.Per != nil {
				result.Max.Jobs.ActiveJobs = &types.QoSLimitsMaxJobsActiveJobs{
					Per: &types.QoSLimitsMaxJobsActiveJobsPer{
						Account: ConvertUint32NoVal(source.Max.Jobs.ActiveJobs.Per.Account),
						User:    ConvertUint32NoVal(source.Max.Jobs.ActiveJobs.Per.User),
					},
				}
			}
			if source.Max.Jobs.Per != nil {
				result.Max.Jobs.Per = &types.QoSLimitsMaxJobsPer{
					Account: ConvertUint32NoVal(source.Max.Jobs.Per.Account),
					User:    ConvertUint32NoVal(source.Max.Jobs.Per.User),
				}
			}
		}

		// Max.TRES
		if source.Max.Tres != nil {
			result.Max.TRES = &types.QoSLimitsMaxTRES{
				Total: convertTresList(source.Max.Tres.Total),
			}
			if source.Max.Tres.Minutes != nil {
				result.Max.TRES.Minutes = &types.QoSLimitsMaxTRESMinutes{
					Total: convertTresList(source.Max.Tres.Minutes.Total),
				}
				if source.Max.Tres.Minutes.Per != nil {
					result.Max.TRES.Minutes.Per = &types.QoSLimitsMaxTRESMinutesPer{
						Account: convertTresList(source.Max.Tres.Minutes.Per.Account),
						Job:     convertTresList(source.Max.Tres.Minutes.Per.Job),
						QoS:     convertTresList(source.Max.Tres.Minutes.Per.Qos),
						User:    convertTresList(source.Max.Tres.Minutes.Per.User),
					}
				}
			}
			if source.Max.Tres.Per != nil {
				result.Max.TRES.Per = &types.QoSLimitsMaxTRESPer{
					Account: convertTresList(source.Max.Tres.Per.Account),
					Job:     convertTresList(source.Max.Tres.Per.Job),
					Node:    convertTresList(source.Max.Tres.Per.Node),
					User:    convertTresList(source.Max.Tres.Per.User),
				}
			}
		}

		// Max.WallClock
		if source.Max.WallClock != nil && source.Max.WallClock.Per != nil {
			result.Max.WallClock = &types.QoSLimitsMaxWallClock{
				Per: &types.QoSLimitsMaxWallClockPer{
					Job: ConvertUint32NoVal(source.Max.WallClock.Per.Job),
					QoS: ConvertUint32NoVal(source.Max.WallClock.Per.Qos),
				},
			}
		}
	}

	// Convert Min
	if source.Min != nil {
		result.Min = &types.QoSLimitsMin{
			PriorityThreshold: ConvertUint32NoVal(source.Min.PriorityThreshold),
		}
		if source.Min.Tres != nil && source.Min.Tres.Per != nil {
			result.Min.TRES = &types.QoSLimitsMinTRES{
				Per: &types.QoSLimitsMinTRESPer{
					Job: convertTresList(source.Min.Tres.Per.Job),
				},
			}
		}
	}

	return result
}

// =============================================================================
// JobResources Helpers
// =============================================================================

// ConvertJobResources converts API V0042JobRes to common JobResources.
// Used by goverter as an extend function.
// Note: v0_0_42 uses string slices instead of enum types for SelectType and CoreStatus.
func ConvertJobResources(source *api.V0042JobRes) *types.JobResources {
	if source == nil {
		return nil
	}

	result := &types.JobResources{
		CPUs: source.Cpus,
	}

	// Convert SelectType (string slice in v0_0_42)
	if len(source.SelectType) > 0 {
		selectType := make([]types.SelectTypeValue, len(source.SelectType))
		for i, st := range source.SelectType {
			selectType[i] = types.SelectTypeValue(st)
		}
		result.SelectType = selectType
	}

	// Convert ThreadsPerCore (NoValStruct to uint16)
	if source.ThreadsPerCore.Set != nil && *source.ThreadsPerCore.Set {
		if source.ThreadsPerCore.Number != nil {
			tpc := uint16(*source.ThreadsPerCore.Number)
			result.ThreadsPerCore = tpc
		}
	}

	// Convert Nodes
	if source.Nodes != nil {
		result.Nodes = convertJobResourcesNodes(source.Nodes)
	}

	return result
}

// convertJobResourcesNodes converts the anonymous Nodes struct to common JobResourcesNodes.
func convertJobResourcesNodes(source *struct {
	Allocation *api.V0042JobResNodes `json:"allocation,omitempty"`
	Count      *int32                `json:"count,omitempty"`
	List       *string               `json:"list,omitempty"`
	SelectType *api.V0042NodeCrType  `json:"select_type,omitempty"`
	Whole      *bool                 `json:"whole,omitempty"`
}) *types.JobResourcesNodes {
	if source == nil {
		return nil
	}

	result := &types.JobResourcesNodes{
		Count: source.Count,
		List:  source.List,
		Whole: source.Whole,
	}

	// Convert SelectType (string slice in v0_0_42)
	if source.SelectType != nil && len(*source.SelectType) > 0 {
		selectType := make([]types.JobResourcesNodesSelectTypeValue, len(*source.SelectType))
		for i, st := range *source.SelectType {
			selectType[i] = types.JobResourcesNodesSelectTypeValue(st)
		}
		result.SelectType = selectType
	}

	// Convert Allocation (slice of nodes)
	if source.Allocation != nil && len(*source.Allocation) > 0 {
		allocation := make([]types.JobResNode, len(*source.Allocation))
		for i, node := range *source.Allocation {
			allocation[i] = convertJobResNode(node)
		}
		result.Allocation = allocation
	}

	return result
}

// convertJobResNode converts API V0042JobResNode to common JobResNode.
func convertJobResNode(source api.V0042JobResNode) types.JobResNode {
	result := types.JobResNode{
		Name: source.Name,
	}

	// Convert CPUs struct
	if source.Cpus != nil {
		result.CPUs = &types.JobResNodeCPUs{
			Count: source.Cpus.Count,
			Used:  source.Cpus.Used,
		}
	}

	// Convert Memory struct
	if source.Memory != nil {
		result.Memory = &types.JobResNodeMemory{
			Allocated: source.Memory.Allocated,
		}
	}

	// Convert Sockets
	if len(source.Sockets) > 0 {
		sockets := make([]types.JobResSocket, len(source.Sockets))
		for i, socket := range source.Sockets {
			sockets[i] = convertJobResSocket(socket)
		}
		result.Sockets = sockets
	}

	return result
}

// convertJobResSocket converts API V0042JobResSocket to common JobResSocket.
func convertJobResSocket(source api.V0042JobResSocket) types.JobResSocket {
	result := types.JobResSocket{
		Index: source.Index,
	}

	// Convert Cores
	if len(source.Cores) > 0 {
		cores := make([]types.JobResCore, len(source.Cores))
		for i, core := range source.Cores {
			cores[i] = convertJobResCore(core)
		}
		result.Cores = cores
	}

	return result
}

// convertJobResCore converts API V0042JobResCore to common JobResCore.
// Note: v0_0_42 uses string slice for Status instead of enum.
func convertJobResCore(source api.V0042JobResCore) types.JobResCore {
	result := types.JobResCore{
		Index: source.Index,
	}

	// Convert Status (string slice in v0_0_42)
	if len(source.Status) > 0 {
		status := make([]types.JobResCoreStatusValue, len(source.Status))
		for i, s := range source.Status {
			status[i] = types.JobResCoreStatusValue(s)
		}
		result.Status = status
	}

	return result
}
