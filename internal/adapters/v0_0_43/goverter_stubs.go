//go:build goverter

// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
// goverter_stubs.go provides stub methods that allow the package to compile
// during goverter generation. These are replaced by goverter_bridge.go
// in normal builds.
package v0_0_43

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_43"
)

func (a *AccountAdapter) convertAPIAccountToCommon(apiObj api.V0043Account) *types.Account {
	return nil
}
func (a *AssociationAdapter) convertAPIAssociationToCommon(apiObj api.V0043Assoc) *types.Association {
	return nil
}
func (a *ClusterAdapter) convertAPIClusterToCommon(apiObj api.V0043ClusterRec) *types.Cluster {
	return nil
}
func (a *JobAdapter) convertAPIJobToCommon(apiObj api.V0043JobInfo) *types.Job {
	return nil
}
func (a *NodeAdapter) convertAPINodeToCommon(apiObj api.V0043Node) *types.Node {
	return nil
}
func (a *PartitionAdapter) convertAPIPartitionToCommon(apiObj api.V0043PartitionInfo) *types.Partition {
	return nil
}
func (a *QoSAdapter) convertAPIQoSToCommon(apiObj api.V0043Qos) *types.QoS {
	return nil
}
func (a *ReservationAdapter) convertAPIReservationToCommon(apiObj api.V0043ReservationInfo) *types.Reservation {
	return nil
}
func (a *UserAdapter) convertAPIUserToCommon(apiObj api.V0043User) *types.User {
	return nil
}
func (a *WCKeyAdapter) convertAPIWCKeyToCommon(apiObj api.V0043Wckey) *types.WCKey {
	return nil
}

// =============================================================================
// Write Converter Stubs (common -> API)
// =============================================================================
func (a *AccountAdapter) convertCommonAccountCreateToAPI(input *types.AccountCreate) *api.V0043Account {
	return nil
}
func (a *AccountAdapter) convertCommonAccountUpdateToAPI(input *types.AccountUpdate) *api.V0043Account {
	return nil
}
func (a *AssociationAdapter) convertCommonAssociationCreateToAPI(input *types.AssociationCreate) *api.V0043Assoc {
	return nil
}
func (a *AssociationAdapter) convertCommonAssociationUpdateToAPI(input *types.AssociationUpdate) *api.V0043Assoc {
	return nil
}
func (a *ClusterAdapter) convertCommonClusterCreateToAPI(input *types.ClusterCreate) *api.V0043ClusterRec {
	return nil
}
func (a *ClusterAdapter) convertCommonClusterUpdateToAPI(input *types.ClusterUpdate) *api.V0043ClusterRec {
	return nil
}
func (a *NodeAdapter) convertCommonNodeUpdateToAPI(input *types.NodeUpdate) *api.V0043Node {
	return nil
}
func (a *QoSAdapter) convertCommonQoSCreateToAPI(input *types.QoSCreate) *api.V0043Qos {
	return nil
}
func (a *QoSAdapter) convertCommonQoSUpdateToAPI(input *types.QoSUpdate) *api.V0043Qos {
	return nil
}
func (a *WCKeyAdapter) convertCommonWCKeyCreateToAPI(input *types.WCKeyCreate) *api.V0043Wckey {
	return nil
}
func (a *UserAdapter) convertCommonUserCreateToAPI(input *types.UserCreate) *api.V0043User {
	return nil
}
func (a *UserAdapter) convertCommonUserUpdateToAPI(input *types.UserUpdate) *api.V0043User {
	return nil
}
func (a *ReservationAdapter) convertCommonReservationCreateToAPI(input *types.ReservationCreate) *api.V0043ReservationDescMsg {
	return nil
}
func (a *ReservationAdapter) convertCommonReservationUpdateToAPI(input *types.ReservationUpdate) *api.V0043ReservationDescMsg {
	return nil
}

// =============================================================================
// Package-level Write Converter Stub Vars
// =============================================================================
// These stub vars satisfy the references in job_helpers.gen.go and node_helpers.gen.go
// during goverter generation. They are replaced by real implementations in
// goverter_bridge.go during normal builds.
var (
	jobWriteConverter  = &stubJobWriteConverter{}
	nodeWriteConverter = &stubNodeWriteConverter{}
)

// Stub implementation for JobWriteConverterGoverter interface
type stubJobWriteConverter struct{}

func (s *stubJobWriteConverter) ConvertCommonJobCreateToAPI(source *types.JobCreate) *api.V0043JobDescMsg {
	return nil
}

// Stub implementation for NodeWriteConverterGoverter interface
type stubNodeWriteConverter struct{}

func (s *stubNodeWriteConverter) ConvertCommonNodeUpdateToAPI(source *types.NodeUpdate) *api.V0043UpdateNodeMsg {
	return nil
}
