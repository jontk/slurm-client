//go:build !goverter

// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
// goverter_bridge.go provides bridge methods that connect adapter methods
// to goverter-generated converters.
package v0_0_40

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_40"
)

// Package-level converter instances
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

func (a *AccountAdapter) convertAPIAccountToCommon(apiObj api.V0040Account) *types.Account {
	return accountConverter.ConvertAPIAccountToCommon(apiObj)
}
func (a *AssociationAdapter) convertAPIAssociationToCommon(apiObj api.V0040Assoc) *types.Association {
	return associationConverter.ConvertAPIAssociationToCommon(apiObj)
}
func (a *ClusterAdapter) convertAPIClusterToCommon(apiObj api.V0040ClusterRec) *types.Cluster {
	return clusterConverter.ConvertAPIClusterToCommon(apiObj)
}
func (a *JobAdapter) convertAPIJobToCommon(apiObj api.V0040JobInfo) *types.Job {
	return jobConverter.ConvertAPIJobToCommon(apiObj)
}
func (a *NodeAdapter) convertAPINodeToCommon(apiObj api.V0040Node) *types.Node {
	return nodeConverter.ConvertAPINodeToCommon(apiObj)
}
func (a *PartitionAdapter) convertAPIPartitionToCommon(apiObj api.V0040PartitionInfo) *types.Partition {
	return partitionConverter.ConvertAPIPartitionToCommon(apiObj)
}
func (a *QoSAdapter) convertAPIQoSToCommon(apiObj api.V0040Qos) *types.QoS {
	return qosConverter.ConvertAPIQoSToCommon(apiObj)
}
func (a *ReservationAdapter) convertAPIReservationToCommon(apiObj api.V0040ReservationInfo) *types.Reservation {
	return reservationConverter.ConvertAPIReservationToCommon(apiObj)
}
func (a *UserAdapter) convertAPIUserToCommon(apiObj api.V0040User) *types.User {
	return userConverter.ConvertAPIUserToCommon(apiObj)
}
func (a *WCKeyAdapter) convertAPIWCKeyToCommon(apiObj api.V0040Wckey) *types.WCKey {
	return wckeyConverter.ConvertAPIWCKeyToCommon(apiObj)
}
