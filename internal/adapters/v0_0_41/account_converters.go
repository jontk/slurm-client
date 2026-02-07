// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"encoding/json"
	"fmt"

	types "github.com/jontk/slurm-client/api"
)

// convertAPIAccountToCommon converts a v0.0.41 API Account to common Account type.
// Uses JSON marshaling workaround for anonymous struct types in v0.0.41.
func (a *AccountAdapter) convertAPIAccountToCommon(apiAccount interface{}) (*types.Account, error) {
	// Marshal the anonymous struct to JSON, then unmarshal to map
	jsonBytes, err := json.Marshal(apiAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal account: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account: %w", err)
	}

	account := &types.Account{}

	// Basic fields - safely extract from map
	if name, ok := data["name"].(string); ok {
		account.Name = name
	}
	if desc, ok := data["description"].(string); ok {
		account.Description = desc
	}
	if org, ok := data["organization"].(string); ok {
		account.Organization = org
	}

	// Handle flags as string array
	if flags, ok := data["flags"].([]interface{}); ok {
		for _, f := range flags {
			if flag, ok := f.(string); ok {
				if flag == "DELETED" || flag == "deleted" {
					account.Flags = append(account.Flags, types.AccountFlagsDeleted)
				}
			}
		}
	}

	// Coordinators - convert from nested objects
	if coords, ok := data["coordinators"].([]interface{}); ok {
		coordinators := make([]types.Coord, 0, len(coords))
		for _, c := range coords {
			if coordMap, ok := c.(map[string]interface{}); ok {
				if name, ok := coordMap["name"].(string); ok {
					coordinators = append(coordinators, types.Coord{Name: name})
				}
			}
		}
		account.Coordinators = coordinators
	}

	return account, nil
}

// convertCommonToAPIAccount converts common Account to v0.0.41 API request
func (a *AccountAdapter) convertCommonToAPIAccount(account *types.Account) {
	// Return a simplified structure that matches what the API expects
	// without relying on undefined specific types
	req := struct {
		Accounts []struct {
			Name         *string `json:"name,omitempty"`
			Description  *string `json:"description,omitempty"`
			Organization *string `json:"organization,omitempty"`
			Coordinators *[]struct {
				Name *string `json:"name,omitempty"`
			} `json:"coordinators,omitempty"`
			Flags *[]string `json:"flags,omitempty"`
		} `json:"accounts"`
	}{
		Accounts: []struct {
			Name         *string `json:"name,omitempty"`
			Description  *string `json:"description,omitempty"`
			Organization *string `json:"organization,omitempty"`
			Coordinators *[]struct {
				Name *string `json:"name,omitempty"`
			} `json:"coordinators,omitempty"`
			Flags *[]string `json:"flags,omitempty"`
		}{{}},
	}
	acc := &req.Accounts[0]
	// Set basic fields
	if account.Name != "" {
		acc.Name = &account.Name
	}
	if account.Description != "" {
		acc.Description = &account.Description
	}
	if account.Organization != "" {
		acc.Organization = &account.Organization
	}
	// Handle flags - check for deleted flag in account.Flags
	for _, flag := range account.Flags {
		if flag == types.AccountFlagsDeleted {
			flags := []string{"DELETED"}
			acc.Flags = &flags
			break
		}
	}
	// Convert coordinators - account.Coordinators is []types.Coord
	if len(account.Coordinators) > 0 {
		coords := make([]struct {
			Name *string `json:"name,omitempty"`
		}, 0, len(account.Coordinators))
		for _, coord := range account.Coordinators {
			coordName := coord.Name // Extract name from Coord struct
			coords = append(coords, struct {
				Name *string `json:"name,omitempty"`
			}{
				Name: &coordName,
			})
		}
		acc.Coordinators = &coords
	}
}
