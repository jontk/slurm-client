// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jontk/slurm-client/internal/api/v0_0_43"
	"github.com/jontk/slurm-client/internal/common/types"
)

// TRESUtils provides utilities for working with TRES (Trackable RESources)
type TRESUtils struct{}

// NewTRESUtils creates a new TRESUtils instance
func NewTRESUtils() *TRESUtils {
	return &TRESUtils{}
}

// ConvertAPITRESToCommon converts API TRES list to common types
func (u *TRESUtils) ConvertAPITRESToCommon(apiTresList v0_0_43.V0043TresList) []types.TRES {
	if apiTresList == nil {
		return []types.TRES{}
	}

	tresList := make([]types.TRES, 0, len(apiTresList))
	for _, apiTres := range apiTresList {
		tres := types.TRES{}

		if apiTres.Id != nil {
			tres.ID = int(*apiTres.Id)
		}
		tres.Type = apiTres.Type
		if apiTres.Name != nil {
			tres.Name = *apiTres.Name
		}
		if apiTres.Count != nil {
			tres.Count = *apiTres.Count
		}

		tresList = append(tresList, tres)
	}

	return tresList
}

// ConvertCommonTRESToAPI converts common TRES list to API format
func (u *TRESUtils) ConvertCommonTRESToAPI(tresList []types.TRES) v0_0_43.V0043TresList {
	if len(tresList) == 0 {
		return v0_0_43.V0043TresList{}
	}

	apiTresList := make(v0_0_43.V0043TresList, 0, len(tresList))
	for _, tres := range tresList {
		apiTres := v0_0_43.V0043Tres{
			Type: tres.Type,
		}

		if tres.ID != 0 {
			id := int32(tres.ID)
			apiTres.Id = &id
		}
		if tres.Name != "" {
			apiTres.Name = &tres.Name
		}
		if tres.Count != 0 {
			apiTres.Count = &tres.Count
		}

		apiTresList = append(apiTresList, apiTres)
	}

	return apiTresList
}

// ParseTRESString parses a TRES string (e.g., "cpu=4,mem=8G,node=1") into TRES list
func (u *TRESUtils) ParseTRESString(tresStr string) ([]types.TRES, error) {
	if tresStr == "" {
		return []types.TRES{}, nil
	}

	tresList := []types.TRES{}
	parts := strings.Split(tresStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		keyValue := strings.SplitN(part, "=", 2)
		if len(keyValue) != 2 {
			return nil, fmt.Errorf("invalid TRES format: %s", part)
		}

		tresType := strings.TrimSpace(keyValue[0])
		countStr := strings.TrimSpace(keyValue[1])

		tres := types.TRES{
			Type: tresType,
		}

		// Handle different count formats
		count, err := u.parseCountValue(countStr)
		if err != nil {
			return nil, fmt.Errorf("invalid TRES count for %s: %w", tresType, err)
		}
		tres.Count = count

		// Set name if it's not a basic type
		if !u.isBasicTRESType(tresType) {
			tres.Name = tresType
		}

		tresList = append(tresList, tres)
	}

	return tresList, nil
}

// FormatTRESString formats a TRES list back to string format
func (u *TRESUtils) FormatTRESString(tresList []types.TRES) string {
	if len(tresList) == 0 {
		return ""
	}

	parts := make([]string, 0, len(tresList))
	for _, tres := range tresList {
		var part string
		if tres.Name != "" && tres.Name != tres.Type {
			part = fmt.Sprintf("%s=%d", tres.Name, tres.Count)
		} else {
			part = fmt.Sprintf("%s=%d", tres.Type, tres.Count)
		}
		parts = append(parts, part)
	}

	return strings.Join(parts, ",")
}

// ExtractTRESByType extracts specific TRES types from a list
func (u *TRESUtils) ExtractTRESByType(tresList []types.TRES, tresType string) *types.TRES {
	for _, tres := range tresList {
		if strings.EqualFold(tres.Type, tresType) {
			return &tres
		}
	}
	return nil
}

// ExtractResourceLimits extracts common resource limits from TRES list
func (u *TRESUtils) ExtractResourceLimits(tresList []types.TRES) (cpus int64, memory int64, nodes int64) {
	for _, tres := range tresList {
		switch strings.ToLower(tres.Type) {
		case "cpu":
			cpus = tres.Count
		case "mem", "memory":
			memory = tres.Count
		case "node":
			nodes = tres.Count
		}
	}
	return cpus, memory, nodes
}

// BuildTRESFromLimits builds a TRES list from resource limits
func (u *TRESUtils) BuildTRESFromLimits(cpus, memory, nodes int64) []types.TRES {
	tresList := []types.TRES{}

	if cpus > 0 {
		tresList = append(tresList, types.TRES{
			Type:  "cpu",
			Count: cpus,
		})
	}
	if memory > 0 {
		tresList = append(tresList, types.TRES{
			Type:  "mem",
			Count: memory,
		})
	}
	if nodes > 0 {
		tresList = append(tresList, types.TRES{
			Type:  "node",
			Count: nodes,
		})
	}

	return tresList
}

// MergeTRESLists merges multiple TRES lists, with later lists taking precedence
func (u *TRESUtils) MergeTRESLists(lists ...[]types.TRES) []types.TRES {
	if len(lists) == 0 {
		return []types.TRES{}
	}

	tresMap := make(map[string]types.TRES)

	for _, list := range lists {
		for _, tres := range list {
			key := u.getTRESKey(tres)
			tresMap[key] = tres
		}
	}

	result := make([]types.TRES, 0, len(tresMap))
	for _, tres := range tresMap {
		result = append(result, tres)
	}

	return result
}

// ValidateTRES validates a TRES entry
func (u *TRESUtils) ValidateTRES(tres types.TRES) error {
	if tres.Type == "" {
		return fmt.Errorf("TRES type is required")
	}
	if tres.Count < 0 {
		return fmt.Errorf("TRES count cannot be negative")
	}
	return nil
}

// validateTRESList validates a list of TRES entries
func (u *TRESUtils) ValidateTRESList(tresList []types.TRES) error {
	seen := make(map[string]bool)
	
	for _, tres := range tresList {
		if err := u.ValidateTRES(tres); err != nil {
			return err
		}
		
		key := u.getTRESKey(tres)
		if seen[key] {
			return fmt.Errorf("duplicate TRES entry: %s", key)
		}
		seen[key] = true
	}
	
	return nil
}

// Private helper functions

// parseCountValue parses various count value formats (e.g., "4", "8G", "1024M")
func (u *TRESUtils) parseCountValue(countStr string) (int64, error) {
	countStr = strings.ToUpper(strings.TrimSpace(countStr))
	
	if countStr == "" {
		return 0, nil
	}

	// Handle suffix multipliers
	multiplier := int64(1)
	if len(countStr) > 1 {
		lastChar := countStr[len(countStr)-1:]
		switch lastChar {
		case "K":
			multiplier = 1024
			countStr = countStr[:len(countStr)-1]
		case "M":
			multiplier = 1024 * 1024
			countStr = countStr[:len(countStr)-1]
		case "G":
			multiplier = 1024 * 1024 * 1024
			countStr = countStr[:len(countStr)-1]
		case "T":
			multiplier = 1024 * 1024 * 1024 * 1024
			countStr = countStr[:len(countStr)-1]
		}
	}

	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid count value: %s", countStr)
	}

	return count * multiplier, nil
}

// isBasicTRESType checks if a TRES type is a basic type (cpu, mem, node, etc.)
func (u *TRESUtils) isBasicTRESType(tresType string) bool {
	basicTypes := []string{"cpu", "mem", "memory", "node", "energy", "gres"}
	for _, basic := range basicTypes {
		if strings.EqualFold(tresType, basic) {
			return true
		}
	}
	return false
}

// getTRESKey generates a unique key for a TRES entry
func (u *TRESUtils) getTRESKey(tres types.TRES) string {
	if tres.Name != "" && tres.Name != tres.Type {
		return fmt.Sprintf("%s:%s", tres.Type, tres.Name)
	}
	return tres.Type
}