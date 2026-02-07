//go:build !goverter

// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
// goverter_bridge.go provides bridge methods that connect adapter methods
// to goverter-generated converters.
package v0_0_44

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_44"
)

// Package-level converter instances for read (API -> common)
var (
	accountConverter     = &AccountConverterGoverterImpl{}
	associationConverter = &AssociationConverterGoverterImpl{}
	clusterConverter     = &ClusterConverterGoverterImpl{}
	jobConverter         = &JobConverterGoverterImpl{}
	nodeConverter        = &NodeConverterGoverterImpl{}
	partitionConverter   = &PartitionConverterGoverterImpl{}
	qosConverter         = &QoSConverterGoverterImpl{}
	reservationConverter = &ReservationConverterGoverterImpl{}
	userConverter        = &UserConverterGoverterImpl{}
	wckeyConverter       = &WCKeyConverterGoverterImpl{}
)

// Package-level converter instances for write (common -> API)
var (
	accountWriteConverter     = &AccountWriteConverterGoverterImpl{}
	associationWriteConverter = &AssociationWriteConverterGoverterImpl{}
	clusterWriteConverter     = &ClusterWriteConverterGoverterImpl{}
	jobWriteConverter         = &JobWriteConverterGoverterImpl{}
	nodeWriteConverter        = &NodeWriteConverterGoverterImpl{}
	qosWriteConverter         = &QoSWriteConverterGoverterImpl{}
	reservationWriteConverter = &ReservationWriteConverterGoverterImpl{}
	userWriteConverter        = &UserWriteConverterGoverterImpl{}
	wckeyWriteConverter       = &WCKeyWriteConverterGoverterImpl{}
)

func (a *AccountAdapter) convertAPIAccountToCommon(apiObj api.V0044Account) *types.Account {
	return accountConverter.ConvertAPIAccountToCommon(apiObj)
}
func (a *AccountAdapter) convertCommonAccountCreateToAPI(input *types.AccountCreate) *api.V0044Account {
	return accountWriteConverter.ConvertCommonAccountCreateToAPI(input)
}
func (a *AccountAdapter) convertCommonAccountUpdateToAPI(input *types.AccountUpdate) *api.V0044Account {
	return accountWriteConverter.ConvertCommonAccountUpdateToAPI(input)
}
func (a *AssociationAdapter) convertAPIAssociationToCommon(apiObj api.V0044Assoc) *types.Association {
	return associationConverter.ConvertAPIAssociationToCommon(apiObj)
}
func (a *AssociationAdapter) convertCommonAssociationCreateToAPI(input *types.AssociationCreate) *api.V0044Assoc {
	return associationWriteConverter.ConvertCommonAssociationCreateToAPI(input)
}
func (a *AssociationAdapter) convertCommonAssociationUpdateToAPI(input *types.AssociationUpdate) *api.V0044Assoc {
	return associationWriteConverter.ConvertCommonAssociationUpdateToAPI(input)
}
func (a *ClusterAdapter) convertAPIClusterToCommon(apiObj api.V0044ClusterRec) *types.Cluster {
	return clusterConverter.ConvertAPIClusterToCommon(apiObj)
}
func (a *ClusterAdapter) convertCommonClusterCreateToAPI(input *types.ClusterCreate) *api.V0044ClusterRec {
	return clusterWriteConverter.ConvertCommonClusterCreateToAPI(input)
}
func (a *ClusterAdapter) convertCommonClusterUpdateToAPI(input *types.ClusterUpdate) *api.V0044ClusterRec {
	return clusterWriteConverter.ConvertCommonClusterUpdateToAPI(input)
}
func (a *JobAdapter) convertAPIJobToCommon(apiObj api.V0044JobInfo) *types.Job {
	return jobConverter.ConvertAPIJobToCommon(apiObj)
}
func (a *JobAdapter) convertCommonJobCreateToAPI(input *types.JobCreate) api.SlurmV0044PostJobSubmitJSONRequestBody {
	if input == nil {
		return api.SlurmV0044PostJobSubmitJSONRequestBody{}
	}
	jobDesc := jobWriteConverter.ConvertCommonJobCreateToAPI(input)
	return api.SlurmV0044PostJobSubmitJSONRequestBody{
		Job:    jobDesc,
		Script: jobDesc.Script,
	}
}
func (a *NodeAdapter) convertAPINodeToCommon(apiObj api.V0044Node) *types.Node {
	return nodeConverter.ConvertAPINodeToCommon(apiObj)
}
func (a *NodeAdapter) convertCommonNodeUpdateToAPI(input *types.NodeUpdate) *api.V0044UpdateNodeMsg {
	return nodeWriteConverter.ConvertCommonNodeUpdateToAPI(input)
}
func (a *PartitionAdapter) convertAPIPartitionToCommon(apiObj api.V0044PartitionInfo) *types.Partition {
	return partitionConverter.ConvertAPIPartitionToCommon(apiObj)
}
func (a *QoSAdapter) convertAPIQoSToCommon(apiObj api.V0044Qos) *types.QoS {
	return qosConverter.ConvertAPIQoSToCommon(apiObj)
}
func (a *QoSAdapter) convertCommonQoSCreateToAPI(input *types.QoSCreate) *api.V0044Qos {
	return qosWriteConverter.ConvertCommonQoSCreateToAPI(input)
}
func (a *QoSAdapter) convertCommonQoSUpdateToAPI(input *types.QoSUpdate) *api.V0044Qos {
	return qosWriteConverter.ConvertCommonQoSUpdateToAPI(input)
}
func (a *ReservationAdapter) convertAPIReservationToCommon(apiObj api.V0044ReservationInfo) *types.Reservation {
	return reservationConverter.ConvertAPIReservationToCommon(apiObj)
}
func (a *ReservationAdapter) convertCommonReservationCreateToAPI(input *types.ReservationCreate) *api.V0044ReservationDescMsg {
	return reservationWriteConverter.ConvertCommonReservationCreateToAPI(input)
}
func (a *ReservationAdapter) convertCommonReservationUpdateToAPI(input *types.ReservationUpdate) *api.V0044ReservationDescMsg {
	return reservationWriteConverter.ConvertCommonReservationUpdateToAPI(input)
}
func (a *UserAdapter) convertAPIUserToCommon(apiObj api.V0044User) *types.User {
	return userConverter.ConvertAPIUserToCommon(apiObj)
}
func (a *UserAdapter) convertCommonUserCreateToAPI(input *types.UserCreate) *api.V0044User {
	return userWriteConverter.ConvertCommonUserCreateToAPI(input)
}
func (a *UserAdapter) convertCommonUserUpdateToAPI(input *types.UserUpdate) *api.V0044User {
	return userWriteConverter.ConvertCommonUserUpdateToAPI(input)
}

func (a *WCKeyAdapter) convertAPIWCKeyToCommon(apiObj api.V0044Wckey) *types.WCKey {
	return wckeyConverter.ConvertAPIWCKeyToCommon(apiObj)
}
func (a *WCKeyAdapter) convertCommonWCKeyCreateToAPI(input *types.WCKeyCreate) *api.V0044Wckey {
	return wckeyWriteConverter.ConvertCommonWCKeyCreateToAPI(input)
}
