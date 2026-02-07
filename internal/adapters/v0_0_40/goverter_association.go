// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_40

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_40"
)

// AssociationConverterGoverter defines the goverter interface for Association conversions.
// goverter:converter
// goverter:output:file association_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_40
// goverter:extend ConvertAssocFlags
// goverter:extend ConvertUint32NoVal
// goverter:extend ConvertQosStringIdList
// goverter:extend ConvertAssocShortToID
// goverter:extend ConvertAccounting
// goverter:extend ConvertAssociationDefault
// goverter:extend ConvertAssociationMax
// goverter:extend ConvertAssociationMin
//
//go:generate goverter gen .
type AssociationConverterGoverter interface {
	// ConvertAPIAssociationToCommon converts API V0040Assoc to common Association type
	//
	// Field name mappings (Source -> Target):
	// goverter:map Id ID | ConvertAssocShortToID
	// goverter:map Qos QoS | ConvertQosStringIdList
	//
	// Complex type conversions:
	// goverter:map Flags | ConvertAssocFlags
	// goverter:map Priority | ConvertUint32NoVal
	// goverter:map Accounting | ConvertAccounting
	// goverter:map Default | ConvertAssociationDefault
	// goverter:map Max | ConvertAssociationMax
	// goverter:map Min | ConvertAssociationMin
	ConvertAPIAssociationToCommon(source api.V0040Assoc) *types.Association
}
