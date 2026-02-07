// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_42

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_42"
)

// enhanceClusterWithSkippedFields adds complex fields that were skipped in generated converter
func (a *ClusterAdapter) enhanceClusterWithSkippedFields(cluster *types.Cluster, apiCluster api.V0042ClusterRec) {
	if cluster == nil {
		return
	}
	// Convert Associations
	cluster.Associations = convertClusterAssociations(apiCluster.Associations)
	// Convert Controller
	cluster.Controller = convertClusterController(apiCluster.Controller)
	// Convert Flags
	cluster.Flags = convertClusterFlags(apiCluster.Flags)
	// Convert TRES
	cluster.TRES = convertTRESList(apiCluster.Tres)
}

// convertClusterAssociations converts API cluster associations to common type
func convertClusterAssociations(apiAssoc *struct {
	Root *api.V0042AssocShort `json:"root,omitempty"`
}) *types.ClusterAssociations {
	if apiAssoc == nil {
		return nil
	}
	result := &types.ClusterAssociations{}
	if apiAssoc.Root != nil {
		result.Root = &types.AssocShort{
			Account:   apiAssoc.Root.Account,
			Cluster:   apiAssoc.Root.Cluster,
			ID:        apiAssoc.Root.Id, // Id in API, ID in common
			Partition: apiAssoc.Root.Partition,
			User:      apiAssoc.Root.User,
		}
	}
	return result
}

// convertClusterController converts API cluster controller to common type
func convertClusterController(apiController *struct {
	Host *string `json:"host,omitempty"`
	Port *int32  `json:"port,omitempty"`
}) *types.ClusterController {
	if apiController == nil {
		return nil
	}
	return &types.ClusterController{
		Host: apiController.Host,
		Port: apiController.Port,
	}
}

// convertClusterFlags converts API cluster flags to common type
// Note: In v0.0.42, Flags is []string, not an enum slice
func convertClusterFlags(apiFlags *api.V0042ClusterRecFlags) []types.ClusterControllerFlagsValue {
	if apiFlags == nil || len(*apiFlags) == 0 {
		return nil
	}
	result := make([]types.ClusterControllerFlagsValue, len(*apiFlags))
	for i, flag := range *apiFlags {
		result[i] = types.ClusterControllerFlagsValue(flag)
	}
	return result
}

// convertTRESList converts API TRES list to common type
func convertTRESList(apiTRES *api.V0042TresList) []types.TRES {
	if apiTRES == nil || len(*apiTRES) == 0 {
		return nil
	}
	result := make([]types.TRES, len(*apiTRES))
	for i, tres := range *apiTRES {
		result[i] = types.TRES{
			Count: tres.Count,
			ID:    tres.Id, // Id in API, ID in common
			Name:  tres.Name,
			Type:  tres.Type,
		}
	}
	return result
}
