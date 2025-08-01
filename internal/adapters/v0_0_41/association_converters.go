// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

import (
	"fmt"
	"strconv"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
)

// convertAPIAssociationToCommon converts a v0.0.41 API association to common association
func (a *AssociationAdapter) convertAPIAssociationToCommon(apiAssoc interface{}) (*types.Association, error) {
	// The associations in v0.0.41 are returned as an array of structs
	// directly in the response
	assoc := &types.Association{}

	// Type assertion to access fields from the anonymous struct
	assocData, ok := apiAssoc.(map[string]interface{})
	if !ok {
		// Try direct struct access
		if v, ok := apiAssoc.(struct {
			Account      *string `json:"account,omitempty"`
			Cluster      *string `json:"cluster,omitempty"`
			User         *string `json:"user,omitempty"`
			Partition    *string `json:"partition,omitempty"`
			Id           *int32  `json:"id,omitempty"`
			IsDefault    *bool   `json:"is_default,omitempty"`
			Comment      *string `json:"comment,omitempty"`
			Default      *struct {
				Qos *string `json:"qos,omitempty"`
			} `json:"default,omitempty"`
			Priority     *struct {
				Number *int32 `json:"number,omitempty"`
			} `json:"priority,omitempty"`
			SharesRaw    *int32  `json:"shares_raw,omitempty"`
		}); ok {
			// Extract fields from struct
			if v.Account != nil {
				assoc.AccountName = *v.Account
			}
			if v.Cluster != nil {
				assoc.Cluster = *v.Cluster
			}
			if v.User != nil {
				assoc.UserName = *v.User
			}
			if v.Partition != nil {
				assoc.Partition = *v.Partition
			}
			if v.Id != nil {
				assoc.ID = strconv.Itoa(int(*v.Id))
			}
			if v.IsDefault != nil {
				assoc.IsDefault = *v.IsDefault
			}
			if v.Comment != nil {
				assoc.Comment = *v.Comment
			}
			if v.Default != nil && v.Default.Qos != nil {
				assoc.DefaultQoS = *v.Default.Qos
			}
			if v.Priority != nil && v.Priority.Number != nil {
				assoc.Priority = *v.Priority.Number
			}
			if v.SharesRaw != nil {
				assoc.SharesRaw = *v.SharesRaw
			}
			return assoc, nil
		}
		return nil, fmt.Errorf("unable to convert association: unexpected type")
	}

	// Map access for interface{} type
	if account, ok := assocData["account"].(string); ok {
		assoc.AccountName = account
	}
	if cluster, ok := assocData["cluster"].(string); ok {
		assoc.Cluster = cluster
	}
	if user, ok := assocData["user"].(string); ok {
		assoc.UserName = user
	}
	if partition, ok := assocData["partition"].(string); ok {
		assoc.Partition = partition
	}
	if id, ok := assocData["id"].(float64); ok {
		assoc.ID = strconv.Itoa(int(id))
	}
	if isDefault, ok := assocData["is_default"].(bool); ok {
		assoc.IsDefault = isDefault
	}
	if comment, ok := assocData["comment"].(string); ok {
		assoc.Comment = comment
	}
	if defaultData, ok := assocData["default"].(map[string]interface{}); ok {
		if qos, ok := defaultData["qos"].(string); ok {
			assoc.DefaultQoS = qos
		}
	}
	if priorityData, ok := assocData["priority"].(map[string]interface{}); ok {
		if priority, ok := priorityData["number"].(float64); ok {
			assoc.Priority = int32(priority)
		}
	}
	if sharesRaw, ok := assocData["shares_raw"].(float64); ok {
		assoc.SharesRaw = int32(sharesRaw)
	}

	// Generate ID if not provided
	if assoc.ID == "" {
		// Create composite ID from account:user:cluster:partition
		assoc.ID = fmt.Sprintf("%s:%s:%s:%s", assoc.AccountName, assoc.UserName, assoc.Cluster, assoc.Partition)
	}

	return assoc, nil
}

// convertCommonToAPIAssociation converts common association to v0.0.41 API format
func (a *AssociationAdapter) convertCommonToAPIAssociation(assoc *types.Association) *api.V0041OpenapiAssocsResp {
	// Note: V0041OpenapiAssocsResp uses anonymous structs which can't be easily constructed
	// This is a limitation of the generated API client
	// We'll need to work around this by using map[string]interface{} or similar approaches
	_ = assoc

	req := &api.V0041OpenapiAssocsResp{}
	
	// We need to build the associations array with the exact structure expected
	// Since we can't directly create the anonymous struct, we'll need to use reflection
	// or work around this limitation
	
	// For now, return an empty response as a placeholder
	// The actual implementation would need to properly construct the anonymous struct
	// This is a limitation of the generated API client
	
	return req
}

// formatTRESMapAssoc formats a TRES map to string format for associations
func formatTRESMapAssoc(tres map[string]int64) string {
	// Convert TRES map to SLURM format (e.g., "cpu=100,mem=1000M")
	result := ""
	for k, v := range tres {
		if result != "" {
			result += ","
		}
		if k == "mem" {
			result += fmt.Sprintf("%s=%dM", k, v)
		} else {
			result += fmt.Sprintf("%s=%d", k, v)
		}
	}
	return result
}
