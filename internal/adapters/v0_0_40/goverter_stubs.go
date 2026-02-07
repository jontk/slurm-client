//go:build goverter

// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
// goverter_stubs.go provides stub methods that allow the package to compile
// during goverter generation. These are replaced by goverter_bridge.go
// in normal builds.
package v0_0_40

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_40"
)

func (a *AccountAdapter) convertAPIAccountToCommon(apiObj api.V0040Account) *types.Account {
	return nil
}
func (a *AssociationAdapter) convertAPIAssociationToCommon(apiObj api.V0040Assoc) *types.Association {
	return nil
}
func (a *ClusterAdapter) convertAPIClusterToCommon(apiObj api.V0040ClusterRec) *types.Cluster {
	return nil
}
func (a *JobAdapter) convertAPIJobToCommon(apiObj api.V0040JobInfo) *types.Job {
	return nil
}
func (a *NodeAdapter) convertAPINodeToCommon(apiObj api.V0040Node) *types.Node {
	return nil
}
func (a *PartitionAdapter) convertAPIPartitionToCommon(apiObj api.V0040PartitionInfo) *types.Partition {
	return nil
}
func (a *QoSAdapter) convertAPIQoSToCommon(apiObj api.V0040Qos) *types.QoS {
	return nil
}
func (a *ReservationAdapter) convertAPIReservationToCommon(apiObj api.V0040ReservationInfo) *types.Reservation {
	return nil
}
func (a *UserAdapter) convertAPIUserToCommon(apiObj api.V0040User) *types.User {
	return nil
}
func (a *WCKeyAdapter) convertAPIWCKeyToCommon(apiObj api.V0040Wckey) *types.WCKey {
	return nil
}
