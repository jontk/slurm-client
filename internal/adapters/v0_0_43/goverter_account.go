// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_43

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_43"
)

// AccountConverterGoverter defines the goverter interface for Account conversions.
// This is a PoC to evaluate if goverter can replace the custom generator.
//
// goverter:converter
// goverter:output:file account_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_43
// goverter:extend ConvertAssocShortSlice
// goverter:extend ConvertCoordSlice
// goverter:extend ConvertAccountFlags
// goverter:extend ConvertCoordNamesToSlice
//
//go:generate goverter gen .
type AccountConverterGoverter interface {
	// ConvertAPIAccountToCommon converts API V0043Account to common Account type
	// goverter:map Associations | ConvertAssocShortSlice
	// goverter:map Coordinators | ConvertCoordSlice
	// goverter:map Flags | ConvertAccountFlags
	ConvertAPIAccountToCommon(source api.V0043Account) *types.Account
	// ConvertCommonAccountCreateToAPI converts common AccountCreate to API V0043Account type
	// goverter:map Coordinators | ConvertCoordNamesToSlice
	// goverter:ignore Associations
	// goverter:ignore Flags
	ConvertCommonAccountCreateToAPI(source *types.AccountCreate) *api.V0043Account
}
