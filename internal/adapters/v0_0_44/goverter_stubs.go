//go:build goverter

// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
// goverter_stubs.go provides stub methods that allow the package to compile
// during goverter generation. These are replaced by goverter_bridge.go
// in normal builds.
package v0_0_44

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_44"
)

// Stub converter instances for goverter generation
// These are placeholders that allow helper files to compile during goverter runs
type stubJobWriteConverter struct{}

func (s *stubJobWriteConverter) ConvertCommonJobCreateToAPI(source *types.JobCreate) *api.V0044JobDescMsg {
	return nil
}

type stubNodeWriteConverter struct{}

func (s *stubNodeWriteConverter) ConvertCommonNodeUpdateToAPI(source *types.NodeUpdate) *api.V0044UpdateNodeMsg {
	return nil
}

var (
	jobWriteConverter  = &stubJobWriteConverter{}
	nodeWriteConverter = &stubNodeWriteConverter{}
)

func (a *AccountAdapter) convertAPIAccountToCommon(apiObj api.V0044Account) *types.Account {
	return nil
}
func (a *AssociationAdapter) convertAPIAssociationToCommon(apiObj api.V0044Assoc) *types.Association {
	return nil
}
func (a *ClusterAdapter) convertAPIClusterToCommon(apiObj api.V0044ClusterRec) *types.Cluster {
	return nil
}
func (a *JobAdapter) convertAPIJobToCommon(apiObj api.V0044JobInfo) *types.Job {
	return nil
}
func (a *NodeAdapter) convertAPINodeToCommon(apiObj api.V0044Node) *types.Node {
	return nil
}
func (a *PartitionAdapter) convertAPIPartitionToCommon(apiObj api.V0044PartitionInfo) *types.Partition {
	return nil
}
func (a *QoSAdapter) convertAPIQoSToCommon(apiObj api.V0044Qos) *types.QoS {
	return nil
}
func (a *ReservationAdapter) convertAPIReservationToCommon(apiObj api.V0044ReservationInfo) *types.Reservation {
	return nil
}
func (a *UserAdapter) convertAPIUserToCommon(apiObj api.V0044User) *types.User {
	return nil
}
func (a *WCKeyAdapter) convertAPIWCKeyToCommon(apiObj api.V0044Wckey) *types.WCKey {
	return nil
}

// =============================================================================
// Write Converter Stubs (common -> API)
// =============================================================================
func (a *AccountAdapter) convertCommonAccountCreateToAPI(input *types.AccountCreate) *api.V0044Account {
	return nil
}
func (a *AccountAdapter) convertCommonAccountUpdateToAPI(input *types.AccountUpdate) *api.V0044Account {
	return nil
}
func (a *AssociationAdapter) convertCommonAssociationCreateToAPI(input *types.AssociationCreate) *api.V0044Assoc {
	return nil
}
func (a *AssociationAdapter) convertCommonAssociationUpdateToAPI(input *types.AssociationUpdate) *api.V0044Assoc {
	return nil
}
func (a *ClusterAdapter) convertCommonClusterCreateToAPI(input *types.ClusterCreate) *api.V0044ClusterRec {
	return nil
}
func (a *ClusterAdapter) convertCommonClusterUpdateToAPI(input *types.ClusterUpdate) *api.V0044ClusterRec {
	return nil
}
func (a *NodeAdapter) convertCommonNodeUpdateToAPI(input *types.NodeUpdate) *api.V0044Node {
	return nil
}
func (a *QoSAdapter) convertCommonQoSCreateToAPI(input *types.QoSCreate) *api.V0044Qos {
	return nil
}
func (a *QoSAdapter) convertCommonQoSUpdateToAPI(input *types.QoSUpdate) *api.V0044Qos {
	return nil
}
func (a *WCKeyAdapter) convertCommonWCKeyCreateToAPI(input *types.WCKeyCreate) *api.V0044Wckey {
	return nil
}
func (a *JobAdapter) convertCommonJobCreateToAPI(input *types.JobCreate) api.SlurmV0044PostJobSubmitJSONRequestBody {
	return api.SlurmV0044PostJobSubmitJSONRequestBody{}
}
func (a *UserAdapter) convertCommonUserCreateToAPI(input *types.UserCreate) *api.V0044User {
	return nil
}
func (a *UserAdapter) convertCommonUserUpdateToAPI(input *types.UserUpdate) *api.V0044User {
	return nil
}
func (a *ReservationAdapter) convertCommonReservationCreateToAPI(input *types.ReservationCreate) *api.V0044ReservationDescMsg {
	return nil
}
func (a *ReservationAdapter) convertCommonReservationUpdateToAPI(input *types.ReservationUpdate) *api.V0044ReservationDescMsg {
	return nil
}
