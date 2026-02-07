//go:build !goverter

// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
// goverter_bridge.go provides bridge methods that connect adapter methods
// to goverter-generated converters. This allows adapters to call
// a.convertAPIEntityToCommon() which delegates to the goverter implementation.
package v0_0_42

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_42"
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
	userWriteConverter        = &UserWriteConverterGoverterImpl{}
	wckeyWriteConverter       = &WCKeyWriteConverterGoverterImpl{}
)

// =============================================================================
// Account Adapter Bridge
// =============================================================================
func (a *AccountAdapter) convertAPIAccountToCommon(apiObj api.V0042Account) *types.Account {
	return accountConverter.ConvertAPIAccountToCommon(apiObj)
}
func (a *AccountAdapter) convertCommonAccountCreateToAPI(input *types.AccountCreate) *api.V0042Account {
	return accountWriteConverter.ConvertCommonAccountCreateToAPI(input)
}
func (a *AccountAdapter) convertCommonAccountUpdateToAPI(input *types.AccountUpdate) *api.V0042Account {
	return accountWriteConverter.ConvertCommonAccountUpdateToAPI(input)
}

// =============================================================================
// Association Adapter Bridge
// =============================================================================
func (a *AssociationAdapter) convertAPIAssociationToCommon(apiObj api.V0042Assoc) *types.Association {
	return associationConverter.ConvertAPIAssociationToCommon(apiObj)
}
func (a *AssociationAdapter) convertCommonAssociationCreateToAPI(input *types.AssociationCreate) *api.V0042Assoc {
	return associationWriteConverter.ConvertCommonAssociationCreateToAPI(input)
}
func (a *AssociationAdapter) convertCommonAssociationUpdateToAPI(input *types.AssociationUpdate) *api.V0042Assoc {
	return associationWriteConverter.ConvertCommonAssociationUpdateToAPI(input)
}

// =============================================================================
// Cluster Adapter Bridge
// =============================================================================
func (a *ClusterAdapter) convertAPIClusterToCommon(apiObj api.V0042ClusterRec) *types.Cluster {
	return clusterConverter.ConvertAPIClusterToCommon(apiObj)
}
func (a *ClusterAdapter) convertCommonClusterCreateToAPI(input *types.ClusterCreate) *api.V0042ClusterRec {
	return clusterWriteConverter.ConvertCommonClusterCreateToAPI(input)
}
func (a *ClusterAdapter) convertCommonClusterUpdateToAPI(input *types.ClusterUpdate) *api.V0042ClusterRec {
	return clusterWriteConverter.ConvertCommonClusterUpdateToAPI(input)
}

// =============================================================================
// Job Adapter Bridge
// =============================================================================
func (a *JobAdapter) convertAPIJobToCommon(apiObj api.V0042JobInfo) *types.Job {
	return jobConverter.ConvertAPIJobToCommon(apiObj)
}

// Note: convertCommonJobCreateToAPI is implemented manually in job_converters_extra.go
// due to the need to return the full request body type (Job + Script).
// =============================================================================
// Node Adapter Bridge
// =============================================================================
func (a *NodeAdapter) convertAPINodeToCommon(apiObj api.V0042Node) *types.Node {
	return nodeConverter.ConvertAPINodeToCommon(apiObj)
}
func (a *NodeAdapter) convertCommonNodeUpdateToAPI(input *types.NodeUpdate) *api.V0042UpdateNodeMsg {
	return nodeWriteConverter.ConvertCommonNodeUpdateToAPI(input)
}

// =============================================================================
// Partition Adapter Bridge
// =============================================================================
func (a *PartitionAdapter) convertAPIPartitionToCommon(apiObj api.V0042PartitionInfo) *types.Partition {
	return partitionConverter.ConvertAPIPartitionToCommon(apiObj)
}

// =============================================================================
// QoS Adapter Bridge
// =============================================================================
func (a *QoSAdapter) convertAPIQoSToCommon(apiObj api.V0042Qos) *types.QoS {
	return qosConverter.ConvertAPIQoSToCommon(apiObj)
}
func (a *QoSAdapter) convertCommonQoSCreateToAPI(input *types.QoSCreate) *api.V0042Qos {
	return qosWriteConverter.ConvertCommonQoSCreateToAPI(input)
}
func (a *QoSAdapter) convertCommonQoSUpdateToAPI(input *types.QoSUpdate) *api.V0042Qos {
	return qosWriteConverter.ConvertCommonQoSUpdateToAPI(input)
}

// =============================================================================
// Reservation Adapter Bridge
// =============================================================================
func (a *ReservationAdapter) convertAPIReservationToCommon(apiObj api.V0042ReservationInfo) *types.Reservation {
	return reservationConverter.ConvertAPIReservationToCommon(apiObj)
}

// =============================================================================
// User Adapter Bridge
// =============================================================================
func (a *UserAdapter) convertAPIUserToCommon(apiObj api.V0042User) *types.User {
	return userConverter.ConvertAPIUserToCommon(apiObj)
}
func (a *UserAdapter) convertCommonUserCreateToAPI(input *types.UserCreate) *api.V0042User {
	return userWriteConverter.ConvertCommonUserCreateToAPI(input)
}
func (a *UserAdapter) convertCommonUserUpdateToAPI(input *types.UserUpdate) *api.V0042User {
	return userWriteConverter.ConvertCommonUserUpdateToAPI(input)
}

// =============================================================================
// WCKey Adapter Bridge
// =============================================================================
func (a *WCKeyAdapter) convertAPIWCKeyToCommon(apiObj api.V0042Wckey) *types.WCKey {
	return wckeyConverter.ConvertAPIWCKeyToCommon(apiObj)
}
func (a *WCKeyAdapter) convertCommonWCKeyCreateToAPI(input *types.WCKeyCreate) *api.V0042Wckey {
	return wckeyWriteConverter.ConvertCommonWCKeyCreateToAPI(input)
}
