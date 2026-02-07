//go:build ignore

// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// converter_helpers.go contains shared helper functions for custom conversion logic
// These are generated into the output files and can be used by converters

package main

import (
	"fmt"
	"strings"
)

// HelperFunction represents a reusable conversion helper
type HelperFunction struct {
	Name        string
	Description string
	Code        string
}

// GetHelperFunctions returns all available helper functions
func GetHelperFunctions(version string) []HelperFunction {
	return []HelperFunction{
		{
			Name:        "convertResumeAfter",
			Description: "Convert ResumeAfter NoValStruct to uint64 timestamp",
			Code: `// convertResumeAfter converts resume after time
func convertResumeAfter(apiResume *api.{{VERSION_PREFIX}}Uint64NoValStruct) *uint64 {
	if apiResume != nil && apiResume.Number != nil && *apiResume.Number > 0 {
		val := uint64(*apiResume.Number)
		return &val
	}
	return nil
}`,
		},
		{
			Name:        "convertNodeEnergy",
			Description: "Convert nested node energy structure",
			Code: `// convertNodeEnergy converts API node energy to common type
func convertNodeEnergy(apiEnergy *api.{{VERSION_PREFIX}}AcctGatherEnergy) *types.NodeEnergy {
	if apiEnergy == nil {
		return nil
	}
	energy := &types.NodeEnergy{}
	if apiEnergy.AverageWatts != nil {
		avgWatts := int32(*apiEnergy.AverageWatts)
		energy.AverageWatts = &avgWatts
	}
	if apiEnergy.BaseConsumedEnergy != nil {
		baseEnergy := *apiEnergy.BaseConsumedEnergy
		energy.BaseConsumedEnergy = &baseEnergy
	}
	if apiEnergy.ConsumedEnergy != nil {
		consumed := *apiEnergy.ConsumedEnergy
		energy.ConsumedEnergy = &consumed
	}
	if apiEnergy.CurrentWatts != nil && apiEnergy.CurrentWatts.Number != nil {
		currentWatts := uint32(*apiEnergy.CurrentWatts.Number)
		energy.CurrentWatts = &currentWatts
	}
	if apiEnergy.LastCollected != nil {
		lastCollected := *apiEnergy.LastCollected
		energy.LastCollected = &lastCollected
	}
	return energy
}`,
		},
		{
			Name:        "convertAssocShortSlice",
			Description: "Convert API AssocShort slice to common AssocShort slice",
			Code: `// convertAssocShortSlice converts API AssocShort slice to common type
func convertAssocShortSlice(apiAssocs *api.{{VERSION_PREFIX}}AssocShortList) []types.AssocShort {
	if apiAssocs == nil {
		return nil
	}
	result := make([]types.AssocShort, len(*apiAssocs))
	for i, assoc := range *apiAssocs {
		result[i] = types.AssocShort{
			Account:   assoc.Account,
			Cluster:   assoc.Cluster,
			ID:        assoc.Id, // Note: Id in API, ID in common
			Partition: assoc.Partition,
			User:      assoc.User,
		}
	}
	return result
}`,
		},
	}
}

// GenerateHelpers generates helper function code for a specific version
func GenerateHelpers(version, apiPrefix string) string {
	helpers := GetHelperFunctions(version)
	var code strings.Builder

	code.WriteString(fmt.Sprintf(`
// Helper functions for custom conversions
// These are version-specific and handle complex conversion logic

`))

	for _, helper := range helpers {
		// Replace version placeholder
		helperCode := strings.ReplaceAll(helper.Code, "{{VERSION_PREFIX}}", apiPrefix)
		code.WriteString(helperCode)
		code.WriteString("\n\n")
	}

	return code.String()
}
