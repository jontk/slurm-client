// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"github.com/jontk/slurm-client/internal/common/types"
)

// convertAPIQoSToCommon converts a v0.0.40 API QoS to common QoS type
// Note: v0.0.40 may not have QoS types, this is a placeholder
func (a *QoSAdapter) convertAPIQoSToCommon(apiQoS interface{}) (*types.QoS, error) {
	qos := &types.QoS{}
	// Placeholder implementation
	return qos, nil
}

// convertCommonQoSCreateToAPI converts common QoSCreate to v0.0.40 API format
// Note: v0.0.40 may not have QoS types, this is a placeholder
func (a *QoSAdapter) convertCommonQoSCreateToAPI(qos *types.QoSCreate) (interface{}, error) {
	// Placeholder implementation
	return nil, nil
}

// convertCommonQoSUpdateToAPI converts common QoSUpdate to v0.0.40 API format
// Note: v0.0.40 may not have QoS types, this is a placeholder
func (a *QoSAdapter) convertCommonQoSUpdateToAPI(existingQoS *types.QoS, update *types.QoSUpdate) (interface{}, error) {
	// Placeholder implementation
	return nil, nil
}
