package v0_0_41

import (
	"fmt"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
)

// convertAPIQoSToCommon converts a v0.0.41 API QoS to common QoS type
func (a *QoSAdapter) convertAPIQoSToCommon(apiQoS interface{}) (*types.QoS, error) {
	// Use map interface for handling anonymous structs in v0.0.41
	qosData, ok := apiQoS.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected QoS data type: %T", apiQoS)
	}

	qos := &types.QoS{}

	// Basic fields - using safe type assertions
	if v, ok := qosData["name"]; ok {
		if name, ok := v.(string); ok {
			qos.Name = name
		}
	}
	if v, ok := qosData["description"]; ok {
		if desc, ok := v.(string); ok {
			qos.Description = desc
		}
	}
	if v, ok := qosData["id"]; ok {
		if id, ok := v.(float64); ok {
			qos.ID = int32(id)
		}
	}

	// Priority
	if v, ok := qosData["priority"]; ok {
		if priority, ok := v.(float64); ok {
			qos.Priority = int(priority)
		}
	}

	// Usage factor and threshold
	if v, ok := qosData["usage_factor"]; ok {
		if factor, ok := v.(float64); ok {
			qos.UsageFactor = factor
		}
	}
	if v, ok := qosData["usage_threshold"]; ok {
		if threshold, ok := v.(float64); ok {
			qos.UsageThreshold = threshold
		}
	}

	// Grace time
	if v, ok := qosData["grace_time"]; ok {
		if graceTime, ok := v.(float64); ok {
			qos.GraceTime = int(graceTime)
		}
	}

	// Skip complex nested structures for now
	// These would need detailed mapping based on the actual API response structure

	return qos, nil
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