// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_40

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_40"
)

// WCKeyConverterGoverter defines the goverter interface for WCKey conversions.
// goverter:converter
// goverter:output:file wckey_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_40
// goverter:extend ConvertWCKeyFlags
// goverter:extend ConvertAccounting
type WCKeyConverterGoverter interface {
	// ConvertAPIWCKeyToCommon converts API V0040Wckey to common WCKey type
	// goverter:map Id ID
	// goverter:map Flags | ConvertWCKeyFlags
	// goverter:map Accounting | ConvertAccounting
	ConvertAPIWCKeyToCommon(source api.V0040Wckey) *types.WCKey
}
