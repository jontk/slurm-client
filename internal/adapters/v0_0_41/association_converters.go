// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"encoding/json"
	"fmt"

	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
)

// convertAPIAssociationToCommon converts a v0.0.41 API association to common association
func (a *AssociationAdapter) convertAPIAssociationToCommon(apiAssoc interface{}) (*types.Association, error) {
	assoc := &types.Association{}
	// v0.0.41 uses anonymous structs from oapi-codegen, convert via JSON round-trip
	assocData, ok := apiAssoc.(map[string]interface{})
	if !ok {
		jsonBytes, err := json.Marshal(apiAssoc)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal association data: %w", err)
		}
		if err := json.Unmarshal(jsonBytes, &assocData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal association data to map: %w", err)
		}
	}
	// Map access for interface{} type
	if account, ok := assocData["account"].(string); ok {
		assoc.Account = &account
	}
	if cluster, ok := assocData["cluster"].(string); ok {
		assoc.Cluster = &cluster
	}
	if user, ok := assocData["user"].(string); ok {
		assoc.User = user
	}
	if partition, ok := assocData["partition"].(string); ok {
		assoc.Partition = &partition
	}
	if id, ok := assocData["id"].(float64); ok {
		idVal := int32(id)
		assoc.ID = &idVal
	}
	if isDefault, ok := assocData["is_default"].(bool); ok {
		assoc.IsDefault = &isDefault
	}
	if comment, ok := assocData["comment"].(string); ok {
		assoc.Comment = &comment
	}
	if defaultData, ok := assocData["default"].(map[string]interface{}); ok {
		if qos, ok := defaultData["qos"].(string); ok {
			assoc.Default = &types.AssociationDefault{QoS: &qos}
		}
	}
	if priorityData, ok := assocData["priority"].(map[string]interface{}); ok {
		if priority, ok := priorityData["number"].(float64); ok {
			priorityVal := uint32(priority)
			assoc.Priority = &priorityVal
		}
	}
	if sharesRaw, ok := assocData["shares_raw"].(float64); ok {
		sharesRawVal := int32(sharesRaw)
		assoc.SharesRaw = &sharesRawVal
	}
	// Flags - convert using helper function pattern
	if v, ok := assocData["flags"]; ok {
		if flags, ok := v.([]interface{}); ok {
			assocFlags := make([]types.AssociationDefaultFlagsValue, 0, len(flags))
			for _, f := range flags {
				if flag, ok := f.(string); ok {
					assocFlags = append(assocFlags, types.AssociationDefaultFlagsValue(flag))
				}
			}
			assoc.Flags = assocFlags
		}
	}
	// Accounting - convert using helper function pattern
	if v, ok := assocData["accounting"]; ok {
		if accountingData, ok := v.([]interface{}); ok {
			assoc.Accounting = convertAssociationAccountingFromSlice(accountingData)
		}
	}
	// QoS list
	if v, ok := assocData["qos"]; ok {
		if qosList, ok := v.([]interface{}); ok {
			qos := make([]string, 0, len(qosList))
			for _, q := range qosList {
				if qosName, ok := q.(string); ok {
					qos = append(qos, qosName)
				}
			}
			assoc.QoS = qos
		}
	}
	// Generate ID if not provided
	if assoc.ID == nil {
		// Create composite ID from account:user:cluster:partition - use 0 as default
		idVal := int32(0)
		assoc.ID = &idVal
	}
	return assoc, nil
}

// convertAssociationAccountingFromSlice converts accounting data from a slice to common Accounting slice.
func convertAssociationAccountingFromSlice(data []interface{}) []types.Accounting {
	if len(data) == 0 {
		return nil
	}
	result := make([]types.Accounting, 0, len(data))
	for _, item := range data {
		if acctData, ok := item.(map[string]interface{}); ok {
			acct := types.Accounting{}
			if id, ok := acctData["id"].(float64); ok {
				idVal := int32(id)
				acct.ID = &idVal
			}
			if start, ok := acctData["start"].(float64); ok {
				startVal := int64(start)
				acct.Start = &startVal
			}
			// Convert TRES
			if tresData, ok := acctData["TRES"].(map[string]interface{}); ok {
				acct.TRES = &types.TRES{}
				if count, ok := tresData["count"].(float64); ok {
					countVal := int64(count)
					acct.TRES.Count = &countVal
				}
				if id, ok := tresData["id"].(float64); ok {
					idVal := int32(id)
					acct.TRES.ID = &idVal
				}
				if name, ok := tresData["name"].(string); ok {
					acct.TRES.Name = &name
				}
				if tresType, ok := tresData["type"].(string); ok {
					acct.TRES.Type = tresType
				}
			}
			// Convert Allocated
			if allocData, ok := acctData["allocated"].(map[string]interface{}); ok {
				acct.Allocated = &types.AccountingAllocated{}
				if seconds, ok := allocData["seconds"].(float64); ok {
					secsVal := int64(seconds)
					acct.Allocated.Seconds = &secsVal
				}
			}
			result = append(result, acct)
		}
	}
	return result
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
