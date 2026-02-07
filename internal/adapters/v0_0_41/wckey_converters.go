// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"fmt"

	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
)

func toInternalWCKeyList(from *api.V0041OpenapiWckeyResp) (*types.WCKeyList, error) {
	if from == nil {
		return nil, fmt.Errorf("input for wckey list conversion is nil")
	}
	wckeys := make([]types.WCKey, 0)
	for _, fromWCKey := range from.Wckeys {
		wckey := types.WCKey{
			Name:    fromWCKey.Name,
			User:    fromWCKey.User,
			Cluster: fromWCKey.Cluster,
			ID:      fromWCKey.Id,
		}
		// Convert Flags using helper function
		if fromWCKey.Flags != nil {
			wckey.Flags = ConvertWCKeyFlags(fromWCKey.Flags)
		}
		// Convert Accounting using helper function
		if fromWCKey.Accounting != nil {
			wckey.Accounting = convertWCKeyAccountingFromAPI(fromWCKey.Accounting)
		}
		wckeys = append(wckeys, wckey)
	}
	return &types.WCKeyList{
		WCKeys: wckeys,
		Total:  len(wckeys),
	}, nil
}

// convertWCKeyAccountingFromAPI converts the anonymous accounting struct from API to common Accounting slice.
// This handles the specific anonymous struct type used in V0041OpenapiWckeyResp.
func convertWCKeyAccountingFromAPI(source *[]struct {
	TRES *struct {
		Count *int64  `json:"count,omitempty"`
		Id    *int32  `json:"id,omitempty"`
		Name  *string `json:"name,omitempty"`
		Type  string  `json:"type"`
	} `json:"TRES,omitempty"`
	Allocated *struct {
		Seconds *int64 `json:"seconds,omitempty"`
	} `json:"allocated,omitempty"`
	Id    *int32 `json:"id,omitempty"`
	Start *int64 `json:"start,omitempty"`
}) []types.Accounting {
	if source == nil || len(*source) == 0 {
		return nil
	}
	result := make([]types.Accounting, len(*source))
	for i, acct := range *source {
		accounting := types.Accounting{
			ID:    acct.Id,
			Start: acct.Start,
		}
		// Convert TRES
		if acct.TRES != nil {
			accounting.TRES = &types.TRES{
				Count: acct.TRES.Count,
				ID:    acct.TRES.Id,
				Name:  acct.TRES.Name,
				Type:  acct.TRES.Type,
			}
		}
		// Convert Allocated
		if acct.Allocated != nil {
			accounting.Allocated = &types.AccountingAllocated{
				Seconds: acct.Allocated.Seconds,
			}
		}
		result[i] = accounting
	}
	return result
}
