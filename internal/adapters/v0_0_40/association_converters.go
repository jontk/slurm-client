// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
)

// convertAPIAssociationToCommon converts a v0.0.40 API Association to common Association type
func (a *AssociationAdapter) convertAPIAssociationToCommon(apiAssociation api.V0040Assoc) (*types.Association, error) {
	association := &types.Association{}

	// Basic fields from V0040Assoc
	if apiAssociation.Account != nil {
		association.AccountName = *apiAssociation.Account
	}
	// V0040Assoc has User inside the Id field (V0040AssocShort)
	if apiAssociation.Id != nil && apiAssociation.Id.User != "" {
		association.UserName = apiAssociation.Id.User
	}
	if apiAssociation.Cluster != nil {
		association.Cluster = *apiAssociation.Cluster
	}
	if apiAssociation.Partition != nil {
		association.Partition = *apiAssociation.Partition
	}
	// Set ID from V0040AssocShort nested in Id field
	if apiAssociation.Id != nil && apiAssociation.Id.Id != nil {
		association.ID = string(*apiAssociation.Id.Id)
	}

	return association, nil
}

// convertCommonAssociationCreateToAPI converts common AssociationCreate to v0.0.40 API format
func (a *AssociationAdapter) convertCommonAssociationCreateToAPI(association *types.AssociationCreate) (*api.V0040AssocShort, error) {
	apiAssociation := &api.V0040AssocShort{}

	// Basic fields for V0040AssocShort
	apiAssociation.Account = &association.AccountName
	apiAssociation.User = association.UserName // User is string, not pointer in V0040AssocShort
	apiAssociation.Cluster = &association.Cluster
	
	if association.Partition != "" {
		apiAssociation.Partition = &association.Partition
	}

	// Note: V0040AssocShort is used for creation requests

	return apiAssociation, nil
}

// extractResourceLimits extracts resource limits from API structures
func (a *AssociationAdapter) extractResourceLimits(limits interface{}) *types.ResourceLimits {
	// This would need to be implemented based on the actual API structure
	// For now, return nil as placeholder
	return nil
}

// extractTRES extracts TRES information from API structures
func (a *AssociationAdapter) extractTRES(tres interface{}) string {
	// This would need to be implemented based on the actual API structure
	// For now, return empty string as placeholder
	return ""
}
