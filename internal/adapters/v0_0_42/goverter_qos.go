// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_42

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_42"
)

// QoSConverterGoverter defines the goverter interface for QoS conversions.
// goverter:converter
// goverter:output:file qos_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_42
// goverter:extend ConvertQoSFlags
// goverter:extend ConvertUint32NoVal
// goverter:extend ConvertFloat64NoVal
// goverter:extend ConvertQoSPreempt
// goverter:extend ConvertQoSLimits
//
//go:generate goverter gen .
type QoSConverterGoverter interface {
	// ConvertAPIQoSToCommon converts API V0042Qos to common QoS type
	// goverter:map Id ID
	// goverter:map Flags | ConvertQoSFlags
	// goverter:map Priority | ConvertUint32NoVal
	// goverter:map UsageFactor | ConvertFloat64NoVal
	// goverter:map UsageThreshold | ConvertFloat64NoVal
	// goverter:map Preempt | ConvertQoSPreempt
	// goverter:map Limits | ConvertQoSLimits
	ConvertAPIQoSToCommon(source api.V0042Qos) *types.QoS
}
