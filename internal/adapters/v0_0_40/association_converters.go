package v0_0_40

import (
	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
)

// convertAPIAssociationToCommon converts a v0.0.40 API Association to common Association type
func (a *AssociationAdapter) convertAPIAssociationToCommon(apiAssociation api.V0040AssocShort) (*types.Association, error) {
	association := &types.Association{}

	// Basic fields (V0040AssocShort only has these fields)
	if apiAssociation.Account != nil {
		association.Account = *apiAssociation.Account
	}
	// User is required in V0040AssocShort
	association.User = apiAssociation.User
	if apiAssociation.Cluster != nil {
		association.Cluster = *apiAssociation.Cluster
	}
	if apiAssociation.Partition != nil {
		association.Partition = *apiAssociation.Partition
	}
	// Set ID if available
	if apiAssociation.Id != nil {
		association.ID = uint32(*apiAssociation.Id)
	}

	return association, nil
}

// convertCommonAssociationCreateToAPI converts common AssociationCreate to v0.0.40 API format
func (a *AssociationAdapter) convertCommonAssociationCreateToAPI(association *types.AssociationCreate) (*api.V0040AssocShort, error) {
	apiAssociation := &api.V0040AssocShort{}

	// Basic fields (V0040AssocShort only supports limited fields)
	apiAssociation.Account = &association.Account
	apiAssociation.User = association.User // User is required and not a pointer
	apiAssociation.Cluster = &association.Cluster
	
	if association.Partition != "" {
		apiAssociation.Partition = &association.Partition
	}

	// Note: V0040AssocShort is a simplified type and doesn't support
	// advanced fields like QoS, shares, priority, flags, etc.
	// Those would require using V0040Assoc type instead.

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