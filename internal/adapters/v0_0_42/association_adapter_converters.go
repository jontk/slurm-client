// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

// convertCommonAssociationCreateToAPI converts common association create type to API format
func (a *AssociationAdapter) convertCommonAssociationCreateToAPI(create *types.AssociationCreate) (*api.V0042OpenapiAssocsResp, error) {
	// Create association entry
	assoc := api.V0042Assoc{
		Account: &create.AccountName,
		User:    create.UserName,
		Cluster: &create.Cluster,
	}
	
	if create.Partition != "" {
		assoc.Partition = &create.Partition
	}
	
	// Return response structure
	return &api.V0042OpenapiAssocsResp{
		Associations: []api.V0042Assoc{assoc},
	}, nil
}

// convertCommonAssociationUpdateToAPI converts common association update type to API format
func (a *AssociationAdapter) convertCommonAssociationUpdateToAPI(update *types.AssociationUpdate) (*api.V0042OpenapiAssocsResp, error) {
	// v0.0.42 requires the full association object for updates
	// This is a simplified version - in production, you'd need to get the existing association first
	assoc := api.V0042Assoc{}
	
	// Apply updates
	if update.DefaultQoS != nil {
		assoc.Default = &struct {
			Qos *string `json:"qos,omitempty"`
		}{
			Qos: update.DefaultQoS,
		}
	}
	
	// TODO: Add more field mappings as needed
	
	// Return response structure
	return &api.V0042OpenapiAssocsResp{
		Associations: []api.V0042Assoc{assoc},
	}, nil
}
