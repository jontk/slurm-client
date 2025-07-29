package v0_0_40

import (
	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
)

// convertAPIAssociationToCommon converts a v0.0.40 API Association to common Association type
func (a *AssociationAdapter) convertAPIAssociationToCommon(apiAssociation api.V0040AssociationShort) (*types.Association, error) {
	association := &types.Association{}

	// Basic fields
	if apiAssociation.Account != nil {
		association.Account = *apiAssociation.Account
	}
	if apiAssociation.User != nil {
		association.User = *apiAssociation.User
	}
	if apiAssociation.Cluster != nil {
		association.Cluster = *apiAssociation.Cluster
	}
	if apiAssociation.Partition != nil {
		association.Partition = *apiAssociation.Partition
	}

	// QoS
	if apiAssociation.Defaults != nil && apiAssociation.Defaults.Qos != nil {
		association.DefaultQoS = *apiAssociation.Defaults.Qos
	}
	if apiAssociation.Qos != nil && len(*apiAssociation.Qos) > 0 {
		association.QoSList = *apiAssociation.Qos
	}

	// Shares and priority
	if apiAssociation.Shares != nil && apiAssociation.Shares.Object != nil && apiAssociation.Shares.Object.Raw != nil {
		shares := int32(*apiAssociation.Shares.Object.Raw)
		association.RawShares = &shares
	}
	if apiAssociation.Priority != nil && apiAssociation.Priority.Number != nil {
		priority := *apiAssociation.Priority.Number
		association.Priority = &priority
	}

	// Parent account
	if apiAssociation.Parent != nil {
		association.ParentAccount = *apiAssociation.Parent
	}

	// Flags
	if apiAssociation.Flags != nil && len(*apiAssociation.Flags) > 0 {
		association.Flags = make([]string, len(*apiAssociation.Flags))
		for i, flag := range *apiAssociation.Flags {
			association.Flags[i] = string(flag)
		}
	}

	// Resource limits
	if apiAssociation.Max != nil {
		if apiAssociation.Max.PerJob != nil {
			association.MaxJobs = a.extractResourceLimits(apiAssociation.Max.PerJob)
		}
		if apiAssociation.Max.PerAccount != nil {
			association.MaxJobsAccrue = a.extractResourceLimits(apiAssociation.Max.PerAccount)
		}
		if apiAssociation.Max.Total != nil {
			association.MaxSubmitJobs = a.extractResourceLimits(apiAssociation.Max.Total)
		}
	}

	// GrpTRES
	if apiAssociation.Max != nil && apiAssociation.Max.PerAccount != nil && apiAssociation.Max.PerAccount.Tres != nil {
		association.GrpTRES = a.extractTRES(apiAssociation.Max.PerAccount.Tres)
	}

	// Is default
	if apiAssociation.IsDefault != nil {
		isDefault := *apiAssociation.IsDefault
		association.IsDefault = &isDefault
	}

	return association, nil
}

// convertCommonAssociationCreateToAPI converts common AssociationCreate to v0.0.40 API format
func (a *AssociationAdapter) convertCommonAssociationCreateToAPI(association *types.AssociationCreate) (*api.V0040AssociationShort, error) {
	apiAssociation := &api.V0040AssociationShort{}

	// Basic fields
	apiAssociation.Account = &association.Account
	apiAssociation.User = &association.User
	apiAssociation.Cluster = &association.Cluster
	
	if association.Partition != "" {
		apiAssociation.Partition = &association.Partition
	}

	// Parent account
	if association.ParentAccount != "" {
		apiAssociation.Parent = &association.ParentAccount
	}

	// QoS
	if association.DefaultQoS != "" {
		defaults := &api.V0040AssociationShortDefaults{
			Qos: &association.DefaultQoS,
		}
		apiAssociation.Defaults = defaults
	}
	if len(association.QoSList) > 0 {
		apiAssociation.Qos = &association.QoSList
	}

	// Shares and priority
	if association.RawShares != nil && *association.RawShares > 0 {
		rawShares := int(*association.RawShares)
		shares := &api.V0040AssociationShortShares{
			Object: &api.V0040AssociationShortSharesObject{
				Raw: &rawShares,
			},
		}
		apiAssociation.Shares = shares
	}
	if association.Priority != nil && *association.Priority > 0 {
		priority := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: association.Priority,
		}
		apiAssociation.Priority = &priority
	}

	// Flags
	if len(association.Flags) > 0 {
		flags := make([]api.V0040AssociationShortFlags, len(association.Flags))
		for i, flag := range association.Flags {
			flags[i] = api.V0040AssociationShortFlags(flag)
		}
		apiAssociation.Flags = &flags
	}

	// Is default
	if association.IsDefault != nil {
		apiAssociation.IsDefault = association.IsDefault
	}

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