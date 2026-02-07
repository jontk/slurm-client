//go:build !goverter

// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
// goverter_bridge.go provides bridge methods that connect adapter methods
// to goverter-generated converters.
package v0_0_43

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_43"
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

func (a *AccountAdapter) convertAPIAccountToCommon(apiObj api.V0043Account) *types.Account {
	return accountConverter.ConvertAPIAccountToCommon(apiObj)
}
func (a *AccountAdapter) convertCommonAccountCreateToAPI(input *types.AccountCreate) *api.V0043Account {
	return accountWriteConverter.ConvertCommonAccountCreateToAPI(input)
}
func (a *AccountAdapter) convertCommonAccountUpdateToAPI(input *types.AccountUpdate) *api.V0043Account {
	return accountWriteConverter.ConvertCommonAccountUpdateToAPI(input)
}
func (a *AssociationAdapter) convertAPIAssociationToCommon(apiObj api.V0043Assoc) *types.Association {
	return associationConverter.ConvertAPIAssociationToCommon(apiObj)
}
func (a *AssociationAdapter) convertCommonAssociationCreateToAPI(input *types.AssociationCreate) *api.V0043Assoc {
	return associationWriteConverter.ConvertCommonAssociationCreateToAPI(input)
}
func (a *AssociationAdapter) convertCommonAssociationUpdateToAPI(input *types.AssociationUpdate) *api.V0043Assoc {
	return associationWriteConverter.ConvertCommonAssociationUpdateToAPI(input)
}
func (a *ClusterAdapter) convertAPIClusterToCommon(apiObj api.V0043ClusterRec) *types.Cluster {
	return clusterConverter.ConvertAPIClusterToCommon(apiObj)
}
func (a *ClusterAdapter) convertCommonClusterCreateToAPI(input *types.ClusterCreate) *api.V0043ClusterRec {
	return clusterWriteConverter.ConvertCommonClusterCreateToAPI(input)
}
func (a *ClusterAdapter) convertCommonClusterUpdateToAPI(input *types.ClusterUpdate) *api.V0043ClusterRec {
	return clusterWriteConverter.ConvertCommonClusterUpdateToAPI(input)
}
func (a *JobAdapter) convertAPIJobToCommon(apiObj api.V0043JobInfo) *types.Job {
	return jobConverter.ConvertAPIJobToCommon(apiObj)
}
func (a *NodeAdapter) convertAPINodeToCommon(apiObj api.V0043Node) *types.Node {
	return nodeConverter.ConvertAPINodeToCommon(apiObj)
}
func (a *NodeAdapter) convertCommonNodeUpdateToAPI(input *types.NodeUpdate) *api.V0043Node {
	if input == nil {
		return nil
	}
	result := &api.V0043Node{
		Comment: input.Comment,
		Reason:  input.Reason,
		Gres:    input.GRES,
	}
	if input.CPUBind != nil {
		result.CpuBinding = input.CPUBind
	}
	if input.Weight != nil {
		w := int32(*input.Weight)
		result.Weight = &w
	}
	// State is handled separately in the node adapter Update method
	return result
}
func (a *PartitionAdapter) convertAPIPartitionToCommon(apiObj api.V0043PartitionInfo) *types.Partition {
	return partitionConverter.ConvertAPIPartitionToCommon(apiObj)
}
func (a *QoSAdapter) convertAPIQoSToCommon(apiObj api.V0043Qos) *types.QoS {
	return qosConverter.ConvertAPIQoSToCommon(apiObj)
}
func (a *QoSAdapter) convertCommonQoSCreateToAPI(input *types.QoSCreate) *api.V0043Qos {
	return qosWriteConverter.ConvertCommonQoSCreateToAPI(input)
}
func (a *QoSAdapter) convertCommonQoSUpdateToAPI(input *types.QoSUpdate) *api.V0043Qos {
	return qosWriteConverter.ConvertCommonQoSUpdateToAPI(input)
}
func (a *ReservationAdapter) convertAPIReservationToCommon(apiObj api.V0043ReservationInfo) *types.Reservation {
	return reservationConverter.ConvertAPIReservationToCommon(apiObj)
}
func (a *ReservationAdapter) convertCommonReservationCreateToAPI(input *types.ReservationCreate) *api.V0043ReservationDescMsg {
	return reservationWriteConverter.ConvertCommonReservationCreateToAPI(input)
}
func (a *ReservationAdapter) convertCommonReservationUpdateToAPI(input *types.ReservationUpdate) *api.V0043ReservationDescMsg {
	return reservationWriteConverter.ConvertCommonReservationUpdateToAPI(input)
}
func (a *UserAdapter) convertAPIUserToCommon(apiObj api.V0043User) *types.User {
	return userConverter.ConvertAPIUserToCommon(apiObj)
}
func (a *UserAdapter) convertCommonUserCreateToAPI(input *types.UserCreate) *api.V0043User {
	return userWriteConverter.ConvertCommonUserCreateToAPI(input)
}
func (a *UserAdapter) convertCommonUserUpdateToAPI(input *types.UserUpdate) *api.V0043User {
	return userWriteConverter.ConvertCommonUserUpdateToAPI(input)
}
func (a *WCKeyAdapter) convertAPIWCKeyToCommon(apiObj api.V0043Wckey) *types.WCKey {
	return wckeyConverter.ConvertAPIWCKeyToCommon(apiObj)
}
func (a *WCKeyAdapter) convertCommonWCKeyCreateToAPI(input *types.WCKeyCreate) *api.V0043Wckey {
	return wckeyWriteConverter.ConvertCommonWCKeyCreateToAPI(input)
}
