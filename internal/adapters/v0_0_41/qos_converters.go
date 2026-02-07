// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"encoding/json"
	"fmt"

	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
)

// convertAPIQoSToCommon converts a v0.0.41 API QoS to common QoS type.
// Uses JSON marshaling workaround for anonymous struct types in v0.0.41.
func (a *QoSAdapter) convertAPIQoSToCommon(apiQoS interface{}) (*types.QoS, error) {
	// Marshal the anonymous struct to JSON, then unmarshal to map
	jsonBytes, err := json.Marshal(apiQoS)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal QoS: %w", err)
	}

	var qosData map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &qosData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal QoS: %w", err)
	}
	qos := &types.QoS{}
	// Basic fields - using safe type assertions
	if v, ok := qosData["name"]; ok {
		if name, ok := v.(string); ok {
			n := name
			qos.Name = &n
		}
	}
	if v, ok := qosData["description"]; ok {
		if desc, ok := v.(string); ok {
			d := desc
			qos.Description = &d
		}
	}
	if v, ok := qosData["id"]; ok {
		if id, ok := v.(float64); ok {
			i := int32(id)
			qos.ID = &i
		}
	}
	// Priority
	if v, ok := qosData["priority"]; ok {
		if priority, ok := v.(float64); ok {
			p := uint32(priority)
			qos.Priority = &p
		}
	}
	// Usage factor and threshold
	if v, ok := qosData["usage_factor"]; ok {
		if factor, ok := v.(float64); ok {
			f := factor
			qos.UsageFactor = &f
		}
	}
	if v, ok := qosData["usage_threshold"]; ok {
		if threshold, ok := v.(float64); ok {
			t := threshold
			qos.UsageThreshold = &t
		}
	}
	// Grace time (now in Limits.GraceTime)
	if v, ok := qosData["grace_time"]; ok {
		if graceTime, ok := v.(float64); ok {
			gt := int32(graceTime)
			if qos.Limits == nil {
				qos.Limits = &types.QoSLimits{}
			}
			qos.Limits.GraceTime = &gt
		}
	}
	// Flags - convert using helper function pattern
	if v, ok := qosData["flags"]; ok {
		if flags, ok := v.([]interface{}); ok {
			qosFlags := make([]types.QoSFlagsValue, 0, len(flags))
			for _, f := range flags {
				if flag, ok := f.(string); ok {
					qosFlags = append(qosFlags, types.QoSFlagsValue(flag))
				}
			}
			qos.Flags = qosFlags
		}
	}
	// Preempt - convert using helper function pattern
	if v, ok := qosData["preempt"]; ok {
		if preemptData, ok := v.(map[string]interface{}); ok {
			qos.Preempt = convertQoSPreemptFromMap(preemptData)
		}
	}
	return qos, nil
}

// convertQoSPreemptFromMap converts preempt data from a map to common QoSPreempt type.
func convertQoSPreemptFromMap(data map[string]interface{}) *types.QoSPreempt {
	if data == nil {
		return nil
	}
	result := &types.QoSPreempt{}
	// Convert exempt_time
	if exemptTime, ok := data["exempt_time"].(map[string]interface{}); ok {
		if number, ok := exemptTime["number"].(float64); ok {
			et := uint32(number)
			result.ExemptTime = &et
		}
	}
	// Convert list - v0_0_41 uses []string for preempt list
	if list, ok := data["list"].([]interface{}); ok {
		preemptList := make([]string, 0, len(list))
		for _, l := range list {
			if name, ok := l.(string); ok {
				preemptList = append(preemptList, name)
			}
		}
		result.List = preemptList
	}
	// Convert mode - v0_0_41 uses []string for preempt modes
	if mode, ok := data["mode"].([]interface{}); ok {
		modes := make([]types.ModeValue, 0, len(mode))
		for _, m := range mode {
			if modeStr, ok := m.(string); ok {
				modes = append(modes, types.ModeValue(modeStr))
			}
		}
		result.Mode = modes
	}
	return result
}

// convertCommonToAPIQoS converts common QoS to v0.0.41 API format
func (a *QoSAdapter) convertCommonToAPIQoS(qos *types.QoS) *api.SlurmdbV0041PostQosJSONRequestBody {
	// Create a basic QoS request structure
	qosReq := &api.SlurmdbV0041PostQosJSONRequestBody{}
	// Note: The exact structure for QoS creation/update in v0.0.41 may be different
	// This is a placeholder implementation that would need to be adjusted
	// based on the actual API structure
	// For now, return an empty request
	_ = qos
	return qosReq
}
