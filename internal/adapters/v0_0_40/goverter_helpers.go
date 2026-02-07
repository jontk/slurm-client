// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
// goverter_helpers.go provides extend functions for goverter converters.
// These functions handle complex type conversions that goverter cannot auto-generate.
package v0_0_40

import (
	"time"

	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_40"
)

// =============================================================================
// NoValStruct Helpers - Generic converters for SLURM's NoValStruct pattern
// =============================================================================
// ConvertTimeNoVal converts a V0040Uint64NoVal to time.Time.
// Returns zero time if source is nil or number is 0.
func ConvertTimeNoVal(source *api.V0040Uint64NoVal) time.Time {
	if source == nil || source.Number == nil || *source.Number == 0 {
		return time.Time{}
	}
	return time.Unix(*source.Number, 0)
}

// ConvertUint64NoVal converts a V0040Uint64NoVal to *uint64.
// Returns nil if source is nil or Set is false.
func ConvertUint64NoVal(source *api.V0040Uint64NoVal) *uint64 {
	if source == nil || source.Set == nil || !*source.Set || source.Number == nil {
		return nil
	}
	val := uint64(*source.Number)
	return &val
}

// ConvertUint32NoVal converts a V0040Uint32NoVal to *uint32.
// Returns nil if source is nil or Set is false.
func ConvertUint32NoVal(source *api.V0040Uint32NoVal) *uint32 {
	if source == nil || source.Set == nil || !*source.Set || source.Number == nil {
		return nil
	}
	val := uint32(*source.Number)
	return &val
}

// ConvertUint16NoVal converts a V0040Uint16NoVal to *uint16.
// Returns nil if source is nil or Set is false.
func ConvertUint16NoVal(source *api.V0040Uint16NoVal) *uint16 {
	if source == nil || source.Set == nil || !*source.Set || source.Number == nil {
		return nil
	}
	val := uint16(*source.Number)
	return &val
}

// ConvertFloat64NoVal converts a V0040Float64NoVal to *float64.
// Returns nil if source is nil or Set is false.
func ConvertFloat64NoVal(source *api.V0040Float64NoVal) *float64 {
	if source == nil || source.Set == nil || !*source.Set || source.Number == nil {
		return nil
	}
	val := *source.Number
	return &val
}

// =============================================================================
// Slice Converters
// =============================================================================
// ConvertAssocShortSlice converts API AssocShort slice to common type.
// Used by goverter as an extend function.
func ConvertAssocShortSlice(source *api.V0040AssocShortList) []types.AssocShort {
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
// Note: V0040Coord has Name string (not pointer) and Direct *bool.
func ConvertCoordSlice(source *api.V0040CoordList) []types.Coord {
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

// ConvertCSVStringToSlice converts a V0040CsvString pointer to a []string.
// V0040CsvString is already []string, so this just dereferences the pointer.
func ConvertCSVStringToSlice(source *api.V0040CsvString) []string {
	if source == nil {
		return nil
	}
	return *source
}

// ConvertCoordNamesToSlice converts a []string (coordinator names) to API V0040CoordList.
// Used for AccountCreate where coordinators are provided as names only.
func ConvertCoordNamesToSlice(source []string) *api.V0040CoordList {
	if len(source) == 0 {
		return nil
	}
	coords := make(api.V0040CoordList, len(source))
	for i, name := range source {
		coords[i] = api.V0040Coord{Name: name}
	}
	return &coords
}

// =============================================================================
// Flag Converters
// =============================================================================
// ConvertAccountFlags converts API AccountFlags slice to common AccountFlagsValue slice.
// Note: v0_0_40 uses []string for AccountFlags (not a typed enum).
func ConvertAccountFlags(source *api.V0040AccountFlags) []types.AccountFlagsValue {
	if source == nil {
		return nil
	}
	flags := make([]types.AccountFlagsValue, len(*source))
	for i, flag := range *source {
		flags[i] = types.AccountFlagsValue(flag)
	}
	return flags
}

// ConvertAssocFlags converts API Association flags slice to common type.
// Note: v0_0_40 uses []string for Association flags (V0040AssocFlags = []string).
func ConvertAssocFlags(source *api.V0040AssocFlags) []types.AssociationDefaultFlagsValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.AssociationDefaultFlagsValue, len(*source))
	for i, flag := range *source {
		result[i] = types.AssociationDefaultFlagsValue(flag)
	}
	return result
}

// ConvertClusterFlags converts API Cluster flags slice to common type.
// Note: v0_0_40 uses []string for Cluster flags (V0040ClusterRecFlags = []string).
func ConvertClusterFlags(source *api.V0040ClusterRecFlags) []types.ClusterControllerFlagsValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.ClusterControllerFlagsValue, len(*source))
	for i, flag := range *source {
		result[i] = types.ClusterControllerFlagsValue(flag)
	}
	return result
}

// ConvertQoSFlags converts API QoS flags slice to common type.
// Note: v0_0_40 uses []string for QoS flags (not a typed enum).
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

// ConvertUserFlags converts API User flags slice to common type.
// Note: v0_0_40 uses []string for User flags (not a typed enum).
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

// ConvertWCKeyFlags converts API WCKey flags slice to common type.
// Note: v0_0_40 uses []string for WCKey flags (V0040WckeyFlags = []string).
func ConvertWCKeyFlags(source *api.V0040WckeyFlags) []types.WCKeyFlagsValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.WCKeyFlagsValue, len(*source))
	for i, flag := range *source {
		result[i] = types.WCKeyFlagsValue(flag)
	}
	return result
}

// ConvertReservationFlags converts API Reservation flags slice to common type.
// Note: v0_0_40 uses []string for Reservation flags (V0040ReservationFlags = []string).
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
// Enum Slice Converters
// =============================================================================
// ConvertNodeStateSlice converts API NodeState slice to common NodeState slice.
// Note: v0_0_40 uses []string for NodeState (not a typed enum).
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
// Note: v0_0_40 uses []string for NextStateAfterReboot (not a typed enum).
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

// ConvertJobStateSlice converts API JobState slice to common JobState slice.
// Note: v0_0_40 uses []string for JobState (via V0040JobState = []string).
func ConvertJobStateSlice(source *api.V0040JobState) []types.JobState {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.JobState, len(*source))
	for i, s := range *source {
		result[i] = types.JobState(s)
	}
	return result
}

// ConvertAdminLevelSlice converts API AdminLevel slice to common type.
// Note: v0_0_40 uses []string for AdminLevel (V0040AdminLvl = []string).
func ConvertAdminLevelSlice(source *api.V0040AdminLvl) []types.AdministratorLevelValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.AdministratorLevelValue, len(*source))
	for i, level := range *source {
		result[i] = types.AdministratorLevelValue(level)
	}
	return result
}

// =============================================================================
// Node-Specific Helpers
// =============================================================================
// ConvertNodeEnergyGoverter converts API node energy to common type.
// Note: V0040AcctGatherEnergy has AverageWatts, BaseConsumedEnergy, ConsumedEnergy,
// CurrentWatts, LastCollected, PreviousConsumedEnergy fields.
func ConvertNodeEnergyGoverter(source *api.V0040AcctGatherEnergy) *types.NodeEnergy {
	if source == nil {
		return nil
	}
	energy := &types.NodeEnergy{}
	if source.AverageWatts != nil {
		energy.AverageWatts = source.AverageWatts
	}
	if source.BaseConsumedEnergy != nil {
		energy.BaseConsumedEnergy = source.BaseConsumedEnergy
	}
	if source.ConsumedEnergy != nil {
		energy.ConsumedEnergy = source.ConsumedEnergy
	}
	// CurrentWatts is a V0040Uint32NoVal
	if source.CurrentWatts != nil && source.CurrentWatts.Set != nil && *source.CurrentWatts.Set && source.CurrentWatts.Number != nil {
		val := uint32(*source.CurrentWatts.Number)
		energy.CurrentWatts = &val
	}
	if source.LastCollected != nil {
		energy.LastCollected = source.LastCollected
	}
	if source.PreviousConsumedEnergy != nil {
		energy.PreviousConsumedEnergy = source.PreviousConsumedEnergy
	}
	return energy
}

// ConvertResumeAfterGoverter converts resume after time.
// Returns nil if source is nil or Set is false.
func ConvertResumeAfterGoverter(source *api.V0040Uint64NoVal) *uint64 {
	if source == nil || source.Set == nil || !*source.Set || source.Number == nil {
		return nil
	}
	val := uint64(*source.Number)
	return &val
}

// ConvertExternalSensors converts API ExternalSensors to common type.
// V0040ExtSensorsData is defined as map[string]interface{}.
func ConvertExternalSensors(source *api.V0040ExtSensorsData) map[string]interface{} {
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
// WCKey Helpers
// =============================================================================
// ConvertWckeySlice converts API Wckey slice to common type.
// Note: V0040Wckey has Cluster string (not pointer), Name string, User string.
func ConvertWckeySlice(source *api.V0040WckeyList) []types.WCKey {
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

// ConvertQosStringIdList converts API QosStringIdList to common []string.
// V0040QosStringIdList is already []string, so this just dereferences the pointer.
func ConvertQosStringIdList(source *api.V0040QosStringIdList) []string {
	if source == nil || len(*source) == 0 {
		return nil
	}
	return *source
}

// =============================================================================
// Association Helpers
// =============================================================================
// ConvertAssocShortToID extracts the ID from a V0040AssocShort struct.
// In v0_0_40, Association.Id is a V0040AssocShort struct containing multiple fields,
// but the common Association type expects just the numeric ID (*int32).
func ConvertAssocShortToID(source *api.V0040AssocShort) *int32 {
	if source == nil || source.Id == nil {
		return nil
	}
	return source.Id
}

// =============================================================================
// Cluster Helpers
// =============================================================================
// ConvertClusterTRES converts API Cluster TRES to common type.
// Note: V0040Tres has Type string (not pointer).
func ConvertClusterTRES(source *api.V0040TresList) []types.TRES {
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
// String/Slice Helpers
// =============================================================================
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
// Job-Specific Helpers
// =============================================================================
// Note: v0_0_40 uses []string for most of these fields (not typed enums).

// ConvertJobFlags converts API JobFlags slice to common FlagsValue slice.
// Used by goverter as an extend function.
func ConvertJobFlags(source *api.V0040JobFlags) []types.FlagsValue {
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
func ConvertJobMailType(source *api.V0040JobMailFlags) []types.MailTypeValue {
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
func ConvertJobProfile(source *api.V0040AcctGatherProfile) []types.ProfileValue {
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
func ConvertJobShared(source *api.V0040JobShared) []types.SharedValue {
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
func ConvertExitCode(source *api.V0040ProcessExitCodeVerbose) *types.ExitCode {
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
func ConvertAccounting(source *api.V0040AccountingList) []types.Accounting {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.Accounting, len(*source))
	for i, acct := range *source {
		accounting := types.Accounting{
			ID:    acct.Id,
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
// Note: v0_0_40 uses []string for Reservation flags.
// Used by goverter as an extend function for read conversions.
func ConvertReservationFlagsRead(source *api.V0040ReservationFlags) []types.ReservationFlagsValue {
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
func ConvertReservationCoreSpec(source *api.V0040ReservationInfoCoreSpec) []types.ReservationCoreSpec {
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
	Time *api.V0040Uint32NoVal `json:"time,omitempty"`
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
// Association Helpers (Additional)
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
// Note: v0_0_40 User Default has Account and Wckey fields (no QoS).
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
		// QoS is not available in v0_0_40
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
	Root *api.V0040AssocShort `json:"root,omitempty"`
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
func ConvertAssocShort(source *api.V0040AssocShort) *types.AssocShort {
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
// Note: v0_0_40 uses []string for preempt modes (V0040QosPreemptModes).
// Used by goverter as an extend function.
func ConvertQoSPreempt(source *struct {
	ExemptTime *api.V0040Uint32NoVal    `json:"exempt_time,omitempty"`
	List       *api.V0040QosPreemptList `json:"list,omitempty"`
	Mode       *api.V0040QosPreemptModes `json:"mode,omitempty"`
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
// Job Helpers (Additional)
// =============================================================================

// ConvertJobPower converts API JobInfo Power to common JobPower.
func ConvertJobPower(source *struct {
	Flags *api.V0040PowerFlags `json:"flags,omitempty"`
}) *types.JobPower {
	if source == nil {
		return nil
	}
	result := &types.JobPower{}
	if source.Flags != nil {
		flags := make([]interface{}, len(*source.Flags))
		for i, f := range *source.Flags {
			flags[i] = f
		}
		result.Flags = flags
	}
	return result
}

// ConvertJobGRESDetail converts API JobInfoGresDetail to common []string.
func ConvertJobGRESDetail(source *api.V0040JobInfoGresDetail) []string {
	if source == nil {
		return nil
	}
	return *source
}

// =============================================================================
// Partition Helpers
// =============================================================================

// ConvertPartitionAccounts converts API PartitionInfo Accounts to common PartitionAccounts.
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

// ConvertPartitionCPUs converts API PartitionInfo Cpus to common PartitionCPUs.
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
func ConvertPartitionDefaults(source *struct {
	Job                    *string              `json:"job,omitempty"`
	MemoryPerCpu           *int64               `json:"memory_per_cpu,omitempty"`
	PartitionMemoryPerCpu  *api.V0040Uint64NoVal `json:"partition_memory_per_cpu,omitempty"`
	PartitionMemoryPerNode *api.V0040Uint64NoVal `json:"partition_memory_per_node,omitempty"`
	Time                   *api.V0040Uint32NoVal `json:"time,omitempty"`
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
func ConvertPartitionMaximums(source *struct {
	CpusPerNode            *api.V0040Uint32NoVal `json:"cpus_per_node,omitempty"`
	CpusPerSocket          *api.V0040Uint32NoVal `json:"cpus_per_socket,omitempty"`
	MemoryPerCpu           *int64                `json:"memory_per_cpu,omitempty"`
	Nodes                  *api.V0040Uint32NoVal `json:"nodes,omitempty"`
	OverTimeLimit          *api.V0040Uint16NoVal `json:"over_time_limit,omitempty"`
	Oversubscribe          *struct {
		Flags *api.V0040OversubscribeFlags `json:"flags,omitempty"`
		Jobs  *int32                        `json:"jobs,omitempty"`
	} `json:"oversubscribe,omitempty"`
	PartitionMemoryPerCpu  *api.V0040Uint64NoVal `json:"partition_memory_per_cpu,omitempty"`
	PartitionMemoryPerNode *api.V0040Uint64NoVal `json:"partition_memory_per_node,omitempty"`
	Shares                 *int32                `json:"shares,omitempty"`
	Time                   *api.V0040Uint32NoVal `json:"time,omitempty"`
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
				flags[i] = types.PartitionMaximumsOversubscribeFlagsValue(f)
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
// Note: v0_0_40 uses V0040PartitionStates which is []string.
func ConvertPartitionPartition(source *struct {
	State *api.V0040PartitionStates `json:"state,omitempty"`
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

// ConvertPartitionQoS converts API PartitionInfo Qos to common PartitionQoS.
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

// ConvertPartitionSuspendTime converts API SuspendTime to time.Time.
func ConvertPartitionSuspendTime(source *api.V0040Uint32NoVal) time.Time {
	if source == nil || source.Set == nil || !*source.Set || (source.Infinite != nil && *source.Infinite) {
		return time.Time{}
	}
	return time.Unix(int64(*source.Number), 0)
}

// ConvertPartitionTimeouts converts API PartitionInfo Timeouts to common PartitionTimeouts.
func ConvertPartitionTimeouts(source *struct {
	Resume  *api.V0040Uint16NoVal `json:"resume,omitempty"`
	Suspend *api.V0040Uint16NoVal `json:"suspend,omitempty"`
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

// ConvertPartitionTRES converts API PartitionInfo Tres to common PartitionTRES.
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

// convertTresList is a helper to convert V0040TresList to []types.TRES.
func convertTresList(source *api.V0040TresList) []types.TRES {
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
		Accruing *api.V0040Uint32NoVal `json:"accruing,omitempty"`
		Active   *api.V0040Uint32NoVal `json:"active,omitempty"`
		Per      *struct {
			Accruing  *api.V0040Uint32NoVal `json:"accruing,omitempty"`
			Count     *api.V0040Uint32NoVal `json:"count,omitempty"`
			Submitted *api.V0040Uint32NoVal `json:"submitted,omitempty"`
			WallClock *api.V0040Uint32NoVal `json:"wall_clock,omitempty"`
		} `json:"per,omitempty"`
		Total *api.V0040Uint32NoVal `json:"total,omitempty"`
	} `json:"jobs,omitempty"`
	Per *struct {
		Account *struct {
			WallClock *api.V0040Uint32NoVal `json:"wall_clock,omitempty"`
		} `json:"account,omitempty"`
	} `json:"per,omitempty"`
	Tres *struct {
		Group *struct {
			Active  *api.V0040TresList `json:"active,omitempty"`
			Minutes *api.V0040TresList `json:"minutes,omitempty"`
		} `json:"group,omitempty"`
		Minutes *struct {
			Per *struct {
				Job *api.V0040TresList `json:"job,omitempty"`
			} `json:"per,omitempty"`
			Total *api.V0040TresList `json:"total,omitempty"`
		} `json:"minutes,omitempty"`
		Per *struct {
			Job  *api.V0040TresList `json:"job,omitempty"`
			Node *api.V0040TresList `json:"node,omitempty"`
		} `json:"per,omitempty"`
		Total *api.V0040TresList `json:"total,omitempty"`
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
	PriorityThreshold *api.V0040Uint32NoVal `json:"priority_threshold,omitempty"`
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
	Factor    *api.V0040Float64NoVal `json:"factor,omitempty"`
	GraceTime *int32                 `json:"grace_time,omitempty"`
	Max       *struct {
		Accruing *struct {
			Per *struct {
				Account *api.V0040Uint32NoVal `json:"account,omitempty"`
				User    *api.V0040Uint32NoVal `json:"user,omitempty"`
			} `json:"per,omitempty"`
		} `json:"accruing,omitempty"`
		ActiveJobs *struct {
			Accruing *api.V0040Uint32NoVal `json:"accruing,omitempty"`
			Count    *api.V0040Uint32NoVal `json:"count,omitempty"`
		} `json:"active_jobs,omitempty"`
		Jobs *struct {
			ActiveJobs *struct {
				Per *struct {
					Account *api.V0040Uint32NoVal `json:"account,omitempty"`
					User    *api.V0040Uint32NoVal `json:"user,omitempty"`
				} `json:"per,omitempty"`
			} `json:"active_jobs,omitempty"`
			Per *struct {
				Account *api.V0040Uint32NoVal `json:"account,omitempty"`
				User    *api.V0040Uint32NoVal `json:"user,omitempty"`
			} `json:"per,omitempty"`
		} `json:"jobs,omitempty"`
		Tres *struct {
			Minutes *struct {
				Per *struct {
					Account *api.V0040TresList `json:"account,omitempty"`
					Job     *api.V0040TresList `json:"job,omitempty"`
					Qos     *api.V0040TresList `json:"qos,omitempty"`
					User    *api.V0040TresList `json:"user,omitempty"`
				} `json:"per,omitempty"`
			} `json:"minutes,omitempty"`
			Per *struct {
				Account *api.V0040TresList `json:"account,omitempty"`
				Job     *api.V0040TresList `json:"job,omitempty"`
				Node    *api.V0040TresList `json:"node,omitempty"`
				User    *api.V0040TresList `json:"user,omitempty"`
			} `json:"per,omitempty"`
			Total *api.V0040TresList `json:"total,omitempty"`
		} `json:"tres,omitempty"`
		WallClock *struct {
			Per *struct {
				Job *api.V0040Uint32NoVal `json:"job,omitempty"`
				Qos *api.V0040Uint32NoVal `json:"qos,omitempty"`
			} `json:"per,omitempty"`
		} `json:"wall_clock,omitempty"`
	} `json:"max,omitempty"`
	Min *struct {
		PriorityThreshold *api.V0040Uint32NoVal `json:"priority_threshold,omitempty"`
		Tres              *struct {
			Per *struct {
				Job *api.V0040TresList `json:"job,omitempty"`
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

		// Max.Jobs (v0_0_40 doesn't have Count field)
		if source.Max.Jobs != nil {
			result.Max.Jobs = &types.QoSLimitsMaxJobs{}
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
				result.Max.TRES.Minutes = &types.QoSLimitsMaxTRESMinutes{}
				// v0_0_40 doesn't have Minutes.Total field
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

// ConvertJobResources converts API V0040JobRes to common JobResources.
// Used by goverter as an extend function.
// Note: v0_0_40 has a different structure than later versions (uses AllocatedCpus, AllocatedNodes).
// AllocatedNodes is []interface{} which requires JSON reflection to parse.
func ConvertJobResources(source *api.V0040JobRes) *types.JobResources {
	if source == nil {
		return nil
	}

	result := &types.JobResources{}

	// Map AllocatedCpus to CPUs (both represent total CPUs)
	if source.AllocatedCpus != nil {
		result.CPUs = *source.AllocatedCpus
	}

	// Map AllocatedNodes ([]interface{}) to Nodes
	// v0_0_40 uses a completely different structure - AllocatedNodes is untyped
	// We attempt to parse it if possible, otherwise create a basic structure
	if source.AllocatedNodes != nil && len(*source.AllocatedNodes) > 0 {
		result.Nodes = convertJobResourcesNodesV40(source.AllocatedNodes, source.AllocatedHosts, source.Nodes)
	}

	return result
}

// convertJobResourcesNodesV40 attempts to convert the untyped AllocatedNodes to common JobResourcesNodes.
// The v0_0_40 API returns nodes as []interface{} which may contain various structures.
func convertJobResourcesNodesV40(nodes *api.V0040JobResNodes, count *int32, list *string) *types.JobResourcesNodes {
	if nodes == nil {
		return nil
	}

	result := &types.JobResourcesNodes{
		Count: count,
		List:  list,
	}

	// Parse the untyped nodes array
	// Each element should be a node with sockets/cores
	allocation := make([]types.JobResNode, 0, len(*nodes))
	for _, nodeInterface := range *nodes {
		if nodeMap, ok := nodeInterface.(map[string]interface{}); ok {
			node := convertJobResNodeFromMap(nodeMap)
			allocation = append(allocation, node)
		}
	}
	if len(allocation) > 0 {
		result.Allocation = allocation
	}

	return result
}

// convertJobResNodeFromMap converts a map[string]interface{} to JobResNode.
func convertJobResNodeFromMap(nodeMap map[string]interface{}) types.JobResNode {
	result := types.JobResNode{}

	// Extract name
	if name, ok := nodeMap["name"].(string); ok {
		result.Name = name
	}

	// Extract CPUs
	if cpusData, ok := nodeMap["cpus"].(map[string]interface{}); ok {
		result.CPUs = &types.JobResNodeCPUs{}
		if count, ok := cpusData["count"].(float64); ok {
			c := int32(count)
			result.CPUs.Count = &c
		}
		if used, ok := cpusData["used"].(float64); ok {
			u := int32(used)
			result.CPUs.Used = &u
		}
	}

	// Extract Memory
	if memoryData, ok := nodeMap["memory"].(map[string]interface{}); ok {
		result.Memory = &types.JobResNodeMemory{}
		if allocated, ok := memoryData["allocated"].(float64); ok {
			a := int64(allocated)
			result.Memory.Allocated = &a
		}
	}

	// Extract Sockets
	if socketsData, ok := nodeMap["sockets"].([]interface{}); ok {
		sockets := make([]types.JobResSocket, 0, len(socketsData))
		for _, socketInterface := range socketsData {
			if socketMap, ok := socketInterface.(map[string]interface{}); ok {
				socket := convertJobResSocketFromMap(socketMap)
				sockets = append(sockets, socket)
			}
		}
		result.Sockets = sockets
	}

	return result
}

// convertJobResSocketFromMap converts a map[string]interface{} to JobResSocket.
func convertJobResSocketFromMap(socketMap map[string]interface{}) types.JobResSocket {
	result := types.JobResSocket{}

	// Extract index
	if index, ok := socketMap["index"].(float64); ok {
		result.Index = int32(index)
	}

	// Extract cores
	if coresData, ok := socketMap["cores"].([]interface{}); ok {
		cores := make([]types.JobResCore, 0, len(coresData))
		for _, coreInterface := range coresData {
			if coreMap, ok := coreInterface.(map[string]interface{}); ok {
				core := convertJobResCoreFromMap(coreMap)
				cores = append(cores, core)
			}
		}
		result.Cores = cores
	}

	return result
}

// convertJobResCoreFromMap converts a map[string]interface{} to JobResCore.
func convertJobResCoreFromMap(coreMap map[string]interface{}) types.JobResCore {
	result := types.JobResCore{}

	// Extract index
	if index, ok := coreMap["index"].(float64); ok {
		result.Index = int32(index)
	}

	// Extract status (string array in v0_0_40)
	if statusData, ok := coreMap["status"].([]interface{}); ok {
		status := make([]types.JobResCoreStatusValue, 0, len(statusData))
		for _, s := range statusData {
			if statusStr, ok := s.(string); ok {
				status = append(status, types.JobResCoreStatusValue(statusStr))
			}
		}
		result.Status = status
	}

	return result
}
