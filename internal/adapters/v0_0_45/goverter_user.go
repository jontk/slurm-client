// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_45

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_45"
)

// UserConverterGoverter defines the goverter interface for User conversions.
//
// goverter:converter
// goverter:output:file user_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_45
// goverter:extend ConvertAdminLevelSlice
// goverter:extend ConvertAssocShortSlice
// goverter:extend ConvertCoordSlice
// goverter:extend ConvertUserFlags
// goverter:extend ConvertWckeySlice
// goverter:extend ConvertUserDefault
//
//go:generate goverter gen .
type UserConverterGoverter interface {
	// ConvertAPIUserToCommon converts API V0045User to common User type
	// goverter:map AdministratorLevel | ConvertAdminLevelSlice
	// goverter:map Associations | ConvertAssocShortSlice
	// goverter:map Coordinators | ConvertCoordSlice
	// goverter:map Flags | ConvertUserFlags
	// goverter:map Wckeys | ConvertWckeySlice
	// goverter:map Default | ConvertUserDefault
	ConvertAPIUserToCommon(source api.V0045User) *types.User
}
