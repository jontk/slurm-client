// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// user_converters_extra.go contains manual helper functions for User conversions
// that complement the goverter-generated converters.
// Note: User write converters (convertCommonUserCreateToAPI, convertCommonUserUpdateToAPI)
// are now implemented via goverter in goverter_write.go and goverter_bridge.go.
package v0_0_42

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_42"
)

// enhanceUserWithSkippedFields adds the skipped fields to a User after base conversion
func (a *UserAdapter) enhanceUserWithSkippedFields(result *types.User, apiObj api.V0042User) {
	if result == nil {
		return
	}
	// AdministratorLevel
	if apiObj.AdministratorLevel != nil {
		for _, level := range *apiObj.AdministratorLevel {
			result.AdministratorLevel = append(result.AdministratorLevel, types.AdministratorLevelValue(level))
		}
	}
	// Associations
	if apiObj.Associations != nil {
		result.Associations = convertAPIAssocShortListToCommon(*apiObj.Associations)
	}
	// Coordinators
	if apiObj.Coordinators != nil {
		result.Coordinators = convertAPICoordListToCommon(*apiObj.Coordinators)
	}
	// Default
	if apiObj.Default != nil {
		result.Default = &types.UserDefault{
			Account: apiObj.Default.Account,
			Wckey:   apiObj.Default.Wckey,
		}
	}
	// Flags
	if apiObj.Flags != nil {
		for _, flag := range *apiObj.Flags {
			result.Flags = append(result.Flags, types.UserDefaultFlagsValue(flag))
		}
	}
	// Wckeys
	if apiObj.Wckeys != nil {
		result.Wckeys = convertAPIWckeyListToCommon(*apiObj.Wckeys)
	}
}

// convertAPIAssocShortListToCommon converts API AssocShortList to common type
func convertAPIAssocShortListToCommon(apiList api.V0042AssocShortList) []types.AssocShort {
	if len(apiList) == 0 {
		return nil
	}
	result := make([]types.AssocShort, 0, len(apiList))
	for _, assoc := range apiList {
		result = append(result, types.AssocShort{
			Account:   assoc.Account,
			Cluster:   assoc.Cluster,
			ID:        assoc.Id,
			Partition: assoc.Partition,
			User:      assoc.User,
		})
	}
	return result
}

// convertAPICoordListToCommon converts API CoordList to common type
func convertAPICoordListToCommon(apiList api.V0042CoordList) []types.Coord {
	if len(apiList) == 0 {
		return nil
	}
	result := make([]types.Coord, 0, len(apiList))
	for _, coord := range apiList {
		result = append(result, types.Coord{
			Direct: coord.Direct,
			Name:   coord.Name,
		})
	}
	return result
}

// convertAPIWckeyListToCommon converts API WckeyList to common type
func convertAPIWckeyListToCommon(apiList api.V0042WckeyList) []types.WCKey {
	if len(apiList) == 0 {
		return nil
	}
	result := make([]types.WCKey, 0, len(apiList))
	for _, wckey := range apiList {
		converted := types.WCKey{
			Cluster: wckey.Cluster,
			ID:      wckey.Id,
			Name:    wckey.Name,
			User:    wckey.User,
		}
		// Accounting
		if wckey.Accounting != nil {
			converted.Accounting = convertAPIAccountingListToCommon(*wckey.Accounting)
		}
		// Flags
		if wckey.Flags != nil {
			for _, flag := range *wckey.Flags {
				converted.Flags = append(converted.Flags, types.WCKeyFlagsValue(flag))
			}
		}
		result = append(result, converted)
	}
	return result
}

// convertAPIAccountingListToCommon converts API AccountingList to common type
func convertAPIAccountingListToCommon(apiList api.V0042AccountingList) []types.Accounting {
	if len(apiList) == 0 {
		return nil
	}
	result := make([]types.Accounting, 0, len(apiList))
	for _, acc := range apiList {
		converted := types.Accounting{
			ID:    acc.Id,
			IDAlt: acc.IdAlt,
			Start: acc.Start,
		}
		// TRES
		if acc.TRES != nil {
			converted.TRES = &types.TRES{
				Count: acc.TRES.Count,
				ID:    acc.TRES.Id,
				Name:  acc.TRES.Name,
				Type:  acc.TRES.Type,
			}
		}
		// Allocated
		if acc.Allocated != nil {
			converted.Allocated = &types.AccountingAllocated{
				Seconds: acc.Allocated.Seconds,
			}
		}
		result = append(result, converted)
	}
	return result
}
