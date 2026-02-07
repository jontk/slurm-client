// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_42

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_42"
)

// ClusterConverterGoverter defines the goverter interface for Cluster conversions.
// goverter:converter
// goverter:output:file cluster_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_42
// goverter:extend ConvertClusterFlags
// goverter:extend ConvertClusterTRES
// goverter:extend ConvertClusterController
// goverter:extend ConvertClusterAssociations
type ClusterConverterGoverter interface {
	// ConvertAPIClusterToCommon converts API V0042ClusterRec to common Cluster type
	// goverter:map Flags | ConvertClusterFlags
	// goverter:map Tres TRES | ConvertClusterTRES
	// goverter:map Controller | ConvertClusterController
	// goverter:map Associations | ConvertClusterAssociations
	ConvertAPIClusterToCommon(source api.V0042ClusterRec) *types.Cluster
}
