// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"fmt"

	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
)

// toInternalClusterList converts API cluster response to common types
// Since both are generated from OpenAPI, we can convert the fields directly
func toInternalClusterList(from *api.V0041OpenapiClustersResp) (*types.ClusterList, error) {
	if from == nil {
		return nil, fmt.Errorf("input for cluster list conversion is nil")
	}
	clusters := make([]types.Cluster, 0, len(from.Clusters))
	for _, fromCluster := range from.Clusters {
		// Convert api.V0041OpenapiClustersRespClusters to types.Cluster
		cluster := types.Cluster{
			Name:         fromCluster.Name,
			Nodes:        fromCluster.Nodes,
			RpcVersion:   fromCluster.RpcVersion,
			SelectPlugin: fromCluster.SelectPlugin,
		}
		// Convert controller if present
		if fromCluster.Controller != nil {
			cluster.Controller = &types.ClusterController{
				Host: fromCluster.Controller.Host,
				Port: fromCluster.Controller.Port,
			}
		}
		// Convert flags
		if fromCluster.Flags != nil {
			cluster.Flags = make([]types.ClusterControllerFlagsValue, len(*fromCluster.Flags))
			for i, flag := range *fromCluster.Flags {
				cluster.Flags[i] = types.ClusterControllerFlagsValue(flag)
			}
		}
		// Convert TRES if present
		if fromCluster.Tres != nil {
			cluster.TRES = make([]types.TRES, len(*fromCluster.Tres))
			for i, tres := range *fromCluster.Tres {
				cluster.TRES[i] = types.TRES{
					Type:  tres.Type,
					Name:  tres.Name,
					ID:    tres.Id,
					Count: tres.Count,
				}
			}
		}
		// Convert Associations if present
		if fromCluster.Associations != nil {
			cluster.Associations = &types.ClusterAssociations{}
			if fromCluster.Associations.Root != nil {
				cluster.Associations.Root = &types.AssocShort{
					Account:   fromCluster.Associations.Root.Account,
					Cluster:   fromCluster.Associations.Root.Cluster,
					ID:        fromCluster.Associations.Root.Id,
					Partition: fromCluster.Associations.Root.Partition,
					User:      fromCluster.Associations.Root.User,
				}
			}
		}
		clusters = append(clusters, cluster)
	}
	return &types.ClusterList{
		Clusters: clusters,
		Total:    len(clusters),
	}, nil
}
