// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// goverter_helpers.go provides helper functions for type conversions in v0_0_41.
// Note: v0_0_41 uses anonymous struct types in OpenAPI responses, so goverter cannot
// be used directly. These helper functions provide similar functionality to goverter
// extend functions, allowing manual converters to use consistent conversion logic.
package v0_0_41

import (
	"time"

	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
)

// =============================================================================
// Job Helpers
// =============================================================================

// ConvertJobFlags converts API JobFlags slice to common FlagsValue slice.
// Used for Job.Flags field conversion.
func ConvertJobFlags(source *[]api.V0041OpenapiJobInfoRespJobsFlags) []types.FlagsValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	flags := make([]types.FlagsValue, len(*source))
	for i, flag := range *source {
		flags[i] = types.FlagsValue(flag)
	}
	return flags
}

// ConvertJobMailType converts API JobMailType slice to common MailTypeValue slice.
// Used for Job.MailType field conversion.
func ConvertJobMailType(source *[]api.V0041OpenapiJobInfoRespJobsMailType) []types.MailTypeValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.MailTypeValue, len(*source))
	for i, mt := range *source {
		result[i] = types.MailTypeValue(mt)
	}
	return result
}

// ConvertJobProfile converts API JobProfile slice to common ProfileValue slice.
// Used for Job.Profile field conversion.
func ConvertJobProfile(source *[]api.V0041OpenapiJobInfoRespJobsProfile) []types.ProfileValue {
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
// Used for Job.Shared field conversion.
func ConvertJobShared(source *[]api.V0041OpenapiJobInfoRespJobsShared) []types.SharedValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.SharedValue, len(*source))
	for i, s := range *source {
		result[i] = types.SharedValue(s)
	}
	return result
}

// ConvertJobState converts API JobState slice to common JobState slice.
// Used for Job.JobState field conversion.
func ConvertJobState(source *[]api.V0041OpenapiJobInfoRespJobsJobState) []types.JobState {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.JobState, len(*source))
	for i, s := range *source {
		result[i] = types.JobState(s)
	}
	return result
}

// ExitCodeData represents the exit code structure from v0_0_41 API response.
// This mirrors the anonymous struct used in V0041OpenapiJobInfoResp.Jobs.ExitCode.
type ExitCodeData struct {
	ReturnCode *struct {
		Infinite *bool  `json:"infinite,omitempty"`
		Number   *int32 `json:"number,omitempty"`
		Set      *bool  `json:"set,omitempty"`
	} `json:"return_code,omitempty"`
	Signal *struct {
		Id *struct {
			Infinite *bool  `json:"infinite,omitempty"`
			Number   *int32 `json:"number,omitempty"`
			Set      *bool  `json:"set,omitempty"`
		} `json:"id,omitempty"`
		Name *string `json:"name,omitempty"`
	} `json:"signal,omitempty"`
	Status *[]string `json:"status,omitempty"`
}

// ConvertExitCode converts an exit code structure to common ExitCode type.
// Used for Job.ExitCode and Job.DerivedExitCode field conversion.
// Note: The status field uses []string in v0_0_41, unlike typed enums in later versions.
func ConvertExitCode(source *ExitCodeData) *types.ExitCode {
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
	// Convert status - v0_0_41 uses []string
	if source.Status != nil {
		for _, s := range *source.Status {
			result.Status = append(result.Status, types.StatusValue(s))
		}
	}
	return result
}

// ConvertExitCodeStatus converts an exit code status slice to common StatusValue slice.
// This handles the typed enum version for DerivedExitCode.
func ConvertExitCodeStatus(source *[]api.V0041OpenapiJobInfoRespJobsDerivedExitCodeStatus) []types.StatusValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.StatusValue, len(*source))
	for i, s := range *source {
		result[i] = types.StatusValue(s)
	}
	return result
}

// ConvertExitCodeStatusFromExitCode converts the exit code status slice to common StatusValue slice.
// This handles the typed enum version for ExitCode.
func ConvertExitCodeStatusFromExitCode(source *[]api.V0041OpenapiJobInfoRespJobsExitCodeStatus) []types.StatusValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.StatusValue, len(*source))
	for i, s := range *source {
		result[i] = types.StatusValue(s)
	}
	return result
}

// =============================================================================
// WCKey Helpers
// =============================================================================

// ConvertWCKeyFlags converts API WCKey flags slice to common WCKeyFlagsValue slice.
// Used for WCKey.Flags field conversion.
func ConvertWCKeyFlags(source *[]api.V0041OpenapiWckeyRespWckeysFlags) []types.WCKeyFlagsValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.WCKeyFlagsValue, len(*source))
	for i, flag := range *source {
		result[i] = types.WCKeyFlagsValue(flag)
	}
	return result
}

// AccountingData represents the accounting structure from v0_0_41 API response.
// This mirrors the anonymous struct used in V0041OpenapiWckeyResp.Wckeys.Accounting.
type AccountingData struct {
	TRES *struct {
		Count *int64  `json:"count,omitempty"`
		Id    *int32  `json:"id,omitempty"`
		Name  *string `json:"name,omitempty"`
		Type  string  `json:"type"`
	} `json:"TRES,omitempty"`
	Allocated *struct {
		Seconds *int64 `json:"seconds,omitempty"`
	} `json:"allocated,omitempty"`
	Id    *int32 `json:"id,omitempty"`
	Start *int64 `json:"start,omitempty"`
}

// ConvertAccounting converts an accounting data slice to common Accounting slice.
// Used for WCKey.Accounting field conversion.
func ConvertAccounting(source *[]AccountingData) []types.Accounting {
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
// Reservation Helpers
// =============================================================================

// ConvertReservationFlagsRead converts API ReservationFlags slice to common ReservationFlagsValue slice.
// Note: v0_0_41 uses typed enum V0041OpenapiReservationRespReservationsFlags for Reservation flags.
// Used for Reservation.Flags field conversion.
func ConvertReservationFlagsRead(source *[]api.V0041OpenapiReservationRespReservationsFlags) []types.ReservationFlagsValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.ReservationFlagsValue, len(*source))
	for i, flag := range *source {
		result[i] = types.ReservationFlagsValue(flag)
	}
	return result
}

// CoreSpecData represents the core specialization structure from v0_0_41 API response.
type CoreSpecData struct {
	Core *string `json:"core,omitempty"`
	Node *string `json:"node,omitempty"`
}

// ConvertReservationCoreSpec converts core specialization data to common ReservationCoreSpec slice.
// Used for Reservation.CoreSpecializations field conversion.
func ConvertReservationCoreSpec(source *[]CoreSpecData) []types.ReservationCoreSpec {
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

// PurgeCompletedData represents the purge completed structure from v0_0_41 API response.
type PurgeCompletedData struct {
	Time *struct {
		Infinite *bool  `json:"infinite,omitempty"`
		Number   *int32 `json:"number,omitempty"`
		Set      *bool  `json:"set,omitempty"`
	} `json:"time,omitempty"`
}

// ConvertReservationPurgeCompleted converts purge completed data to common ReservationPurgeCompleted.
// Used for Reservation.PurgeCompleted field conversion.
func ConvertReservationPurgeCompleted(source *PurgeCompletedData) *types.ReservationPurgeCompleted {
	if source == nil {
		return nil
	}
	result := &types.ReservationPurgeCompleted{}
	if source.Time != nil && source.Time.Number != nil {
		time := uint32(*source.Time.Number)
		result.Time = &time
	}
	return result
}

// =============================================================================
// Association Helpers
// =============================================================================

// ConvertAssocFlags converts API Association flags slice to common AssociationDefaultFlagsValue slice.
// Note: v0_0_41 uses typed enum V0041OpenapiAssocsRespAssociationsFlags for Association flags.
// Used for Association.Flags field conversion.
func ConvertAssocFlags(source *[]api.V0041OpenapiAssocsRespAssociationsFlags) []types.AssociationDefaultFlagsValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.AssociationDefaultFlagsValue, len(*source))
	for i, flag := range *source {
		result[i] = types.AssociationDefaultFlagsValue(flag)
	}
	return result
}

// AssociationDefaultData represents the default structure from v0_0_41 API response.
type AssociationDefaultData struct {
	Qos *string `json:"qos,omitempty"`
}

// ConvertAssociationDefault converts association default data to common AssociationDefault.
// Used for Association.Default field conversion.
func ConvertAssociationDefault(source *AssociationDefaultData) *types.AssociationDefault {
	if source == nil {
		return nil
	}
	return &types.AssociationDefault{
		QoS: source.Qos,
	}
}

// =============================================================================
// User Helpers
// =============================================================================

// ConvertAdminLevelSlice converts API AdminLevel slice to common AdministratorLevelValue slice.
// Note: v0_0_41 uses typed enum for AdminLevel.
func ConvertAdminLevelSlice(source *[]string) []types.AdministratorLevelValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.AdministratorLevelValue, len(*source))
	for i, level := range *source {
		result[i] = types.AdministratorLevelValue(level)
	}
	return result
}

// UserDefaultData represents the default structure from v0_0_41 API response.
// Note: v0_0_41 User Default has Account and Wckey fields (no QoS field).
type UserDefaultData struct {
	Account *string `json:"account,omitempty"`
	Wckey   *string `json:"wckey,omitempty"`
}

// ConvertUserDefault converts user default data to common UserDefault.
// Used for User.Default field conversion.
func ConvertUserDefault(source *UserDefaultData) *types.UserDefault {
	if source == nil {
		return nil
	}
	return &types.UserDefault{
		Account: source.Account,
		Wckey:   source.Wckey,
		// QoS is not available in v0_0_41
	}
}

// ConvertUserFlags converts API User flags slice to common UserDefaultFlagsValue slice.
// Note: v0_0_41 uses typed enum V0041OpenapiUsersRespUsersFlags for User flags.
func ConvertUserFlags(source *[]api.V0041OpenapiUsersRespUsersFlags) []types.UserDefaultFlagsValue {
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
// Cluster Helpers
// =============================================================================

// ClusterControllerData represents the controller structure from v0_0_41 API response.
type ClusterControllerData struct {
	Host *string `json:"host,omitempty"`
	Port *int32  `json:"port,omitempty"`
}

// ConvertClusterController converts cluster controller data to common ClusterController.
// Used for Cluster.Controller field conversion.
func ConvertClusterController(source *ClusterControllerData) *types.ClusterController {
	if source == nil {
		return nil
	}
	return &types.ClusterController{
		Host: source.Host,
		Port: source.Port,
	}
}

// ConvertClusterFlags converts API Cluster flags slice to common ClusterControllerFlagsValue slice.
// Note: v0_0_41 uses typed enum V0041OpenapiClustersRespClustersFlags for Cluster flags.
func ConvertClusterFlags(source *[]api.V0041OpenapiClustersRespClustersFlags) []types.ClusterControllerFlagsValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.ClusterControllerFlagsValue, len(*source))
	for i, flag := range *source {
		result[i] = types.ClusterControllerFlagsValue(flag)
	}
	return result
}

// AssocShortData represents the short association structure from v0_0_41 API response.
type AssocShortData struct {
	Account   *string `json:"account,omitempty"`
	Cluster   *string `json:"cluster,omitempty"`
	Id        *int32  `json:"id,omitempty"`
	Partition *string `json:"partition,omitempty"`
	User      string  `json:"user"`
}

// ConvertAssocShort converts association short data to common AssocShort.
func ConvertAssocShort(source *AssocShortData) *types.AssocShort {
	if source == nil {
		return nil
	}
	return &types.AssocShort{
		Account:   source.Account,
		Cluster:   source.Cluster,
		ID:        source.Id,
		Partition: source.Partition,
		User:      source.User,
	}
}

// ClusterAssociationsData represents the associations structure from v0_0_41 API response.
type ClusterAssociationsData struct {
	Root *AssocShortData `json:"root,omitempty"`
}

// ConvertClusterAssociations converts cluster associations data to common ClusterAssociations.
// Used for Cluster.Associations field conversion.
func ConvertClusterAssociations(source *ClusterAssociationsData) *types.ClusterAssociations {
	if source == nil {
		return nil
	}
	result := &types.ClusterAssociations{}
	if source.Root != nil {
		result.Root = ConvertAssocShort(source.Root)
	}
	return result
}

// =============================================================================
// QoS Helpers
// =============================================================================

// ConvertQoSFlags converts API QoS flags slice to common QoSFlagsValue slice.
// Note: v0_0_41 uses typed enum V0041OpenapiSlurmdbdQosRespQosFlags for QoS flags.
func ConvertQoSFlags(source *[]api.V0041OpenapiSlurmdbdQosRespQosFlags) []types.QoSFlagsValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.QoSFlagsValue, len(*source))
	for i, flag := range *source {
		result[i] = types.QoSFlagsValue(flag)
	}
	return result
}

// QoSPreemptData represents the preempt structure from v0_0_41 API response.
type QoSPreemptData struct {
	ExemptTime *struct {
		Infinite *bool  `json:"infinite,omitempty"`
		Number   *int32 `json:"number,omitempty"`
		Set      *bool  `json:"set,omitempty"`
	} `json:"exempt_time,omitempty"`
	List *[]string `json:"list,omitempty"`
	Mode *[]string `json:"mode,omitempty"`
}

// ConvertQoSPreempt converts QoS preempt data to common QoSPreempt.
// Note: v0_0_41 uses []string for preempt modes (not typed enums).
// Used for QoS.Preempt field conversion.
func ConvertQoSPreempt(source *QoSPreemptData) *types.QoSPreempt {
	if source == nil {
		return nil
	}
	result := &types.QoSPreempt{}
	if source.ExemptTime != nil && source.ExemptTime.Number != nil {
		time := uint32(*source.ExemptTime.Number)
		result.ExemptTime = &time
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
// NoValStruct Helpers - Generic converters for SLURM's NoValStruct pattern
// =============================================================================

// NoValInt32 represents the int32 NoValStruct from v0_0_41 API response.
type NoValInt32 struct {
	Infinite *bool  `json:"infinite,omitempty"`
	Number   *int32 `json:"number,omitempty"`
	Set      *bool  `json:"set,omitempty"`
}

// NoValInt64 represents the int64 NoValStruct from v0_0_41 API response.
type NoValInt64 struct {
	Infinite *bool  `json:"infinite,omitempty"`
	Number   *int64 `json:"number,omitempty"`
	Set      *bool  `json:"set,omitempty"`
}

// NoValFloat64 represents the float64 NoValStruct from v0_0_41 API response.
type NoValFloat64 struct {
	Infinite *bool    `json:"infinite,omitempty"`
	Number   *float64 `json:"number,omitempty"`
	Set      *bool    `json:"set,omitempty"`
}

// ConvertTimeNoVal converts a NoValInt64 to time.Time.
// Returns zero time if source is nil or number is 0.
func ConvertTimeNoVal(source *NoValInt64) time.Time {
	if source == nil || source.Number == nil || *source.Number == 0 {
		return time.Time{}
	}
	return time.Unix(*source.Number, 0)
}

// ConvertUint32NoVal converts a NoValInt32 to *uint32.
// Returns nil if source is nil or Set is false.
func ConvertUint32NoVal(source *NoValInt32) *uint32 {
	if source == nil || source.Set == nil || !*source.Set || source.Number == nil {
		return nil
	}
	val := uint32(*source.Number)
	return &val
}

// ConvertUint16NoVal converts a NoValInt32 to *uint16.
// Returns nil if source is nil or Set is false.
func ConvertUint16NoVal(source *NoValInt32) *uint16 {
	if source == nil || source.Set == nil || !*source.Set || source.Number == nil {
		return nil
	}
	val := uint16(*source.Number)
	return &val
}

// ConvertUint64NoVal converts a NoValInt64 to *uint64.
// Returns nil if source is nil or Set is false.
func ConvertUint64NoVal(source *NoValInt64) *uint64 {
	if source == nil || source.Set == nil || !*source.Set || source.Number == nil {
		return nil
	}
	val := uint64(*source.Number)
	return &val
}

// ConvertFloat64NoVal converts a NoValFloat64 to *float64.
// Returns nil if source is nil or Set is false.
func ConvertFloat64NoVal(source *NoValFloat64) *float64 {
	if source == nil || source.Set == nil || !*source.Set || source.Number == nil {
		return nil
	}
	val := *source.Number
	return &val
}

// =============================================================================
// Account Helpers
// =============================================================================

// ConvertAccountFlags converts API Account flags slice to common AccountFlagsValue slice.
// Note: v0_0_41 uses typed enum V0041OpenapiAccountsRespAccountsFlags for Account flags.
func ConvertAccountFlags(source *[]api.V0041OpenapiAccountsRespAccountsFlags) []types.AccountFlagsValue {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.AccountFlagsValue, len(*source))
	for i, flag := range *source {
		result[i] = types.AccountFlagsValue(flag)
	}
	return result
}

// CoordData represents the coordinator structure from v0_0_41 API response.
type CoordData struct {
	Name   string `json:"name"`
	Direct *bool  `json:"direct,omitempty"`
}

// ConvertCoordSlice converts coordinator data slice to common Coord slice.
// Used for Account.Coordinators and User.Coordinators field conversion.
func ConvertCoordSlice(source *[]CoordData) []types.Coord {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.Coord, len(*source))
	for i, c := range *source {
		result[i] = types.Coord{
			Name:   c.Name,
			Direct: c.Direct,
		}
	}
	return result
}

// ConvertAssocShortSlice converts association short data slice to common AssocShort slice.
// Used for User.Associations field conversion.
func ConvertAssocShortSlice(source *[]AssocShortData) []types.AssocShort {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.AssocShort, len(*source))
	for i, assoc := range *source {
		result[i] = types.AssocShort{
			Account:   assoc.Account,
			Cluster:   assoc.Cluster,
			ID:        assoc.Id,
			Partition: assoc.Partition,
			User:      assoc.User,
		}
	}
	return result
}

// =============================================================================
// TRES Helpers
// =============================================================================

// TRESData represents the TRES structure from v0_0_41 API response.
type TRESData struct {
	Count *int64  `json:"count,omitempty"`
	Id    *int32  `json:"id,omitempty"`
	Name  *string `json:"name,omitempty"`
	Type  string  `json:"type"`
}

// ConvertTRESSlice converts TRES data slice to common TRES slice.
// Used for Cluster.TRES field conversion.
func ConvertTRESSlice(source *[]TRESData) []types.TRES {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.TRES, len(*source))
	for i, t := range *source {
		result[i] = types.TRES{
			Count: t.Count,
			ID:    t.Id,
			Name:  t.Name,
			Type:  t.Type,
		}
	}
	return result
}

// =============================================================================
// WCKey from User Helpers
// =============================================================================

// WCKeyData represents the WCKey structure from v0_0_41 API response.
// Used for User.Wckeys field conversion.
type WCKeyData struct {
	Cluster string `json:"cluster"`
	Id      *int32 `json:"id,omitempty"`
	Name    string `json:"name"`
	User    string `json:"user"`
}

// ConvertWCKeySlice converts WCKey data slice to common WCKey slice.
// Used for User.Wckeys field conversion.
func ConvertWCKeySlice(source *[]WCKeyData) []types.WCKey {
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
