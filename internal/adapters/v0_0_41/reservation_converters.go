package v0_0_41

import (
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
)

// convertAPIReservationToCommon converts a v0.0.41 API Reservation to common Reservation type
func (a *ReservationAdapter) convertAPIReservationToCommon(apiRes interface{}) (*types.Reservation, error) {
	// Use map interface for handling anonymous structs in v0.0.41
	resData, ok := apiRes.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected reservation data type: %T", apiRes)
	}

	res := &types.Reservation{}

	// Basic fields - using safe type assertions
	if v, ok := resData["name"]; ok {
		if name, ok := v.(string); ok {
			res.Name = name
		}
	}
	if v, ok := resData["node_list"]; ok {
		if nodeList, ok := v.(string); ok {
			res.NodeList = nodeList
		}
	}
	if v, ok := resData["node_count"]; ok {
		if nodeCount, ok := v.(float64); ok {
			res.NodeCount = int32(nodeCount)
		}
	}
	if v, ok := resData["core_count"]; ok {
		if coreCount, ok := v.(float64); ok {
			res.CoreCount = int32(coreCount)
		}
	}

	// String array fields
	if v, ok := resData["accounts"]; ok {
		if accountsStr, ok := v.(string); ok && accountsStr != "" {
			res.Accounts = strings.Split(accountsStr, ",")
		}
	}
	if v, ok := resData["users"]; ok {
		if usersStr, ok := v.(string); ok && usersStr != "" {
			res.Users = strings.Split(usersStr, ",")
		}
	}
	if v, ok := resData["groups"]; ok {
		if groupsStr, ok := v.(string); ok && groupsStr != "" {
			res.Groups = strings.Split(groupsStr, ",")
		}
	}
	if v, ok := resData["features"]; ok {
		if featuresStr, ok := v.(string); ok && featuresStr != "" {
			res.Features = strings.Split(featuresStr, ",")
		}
	}

	// Optional string fields
	if v, ok := resData["partition"]; ok {
		if partition, ok := v.(string); ok {
			res.PartitionName = partition
		}
	}
	if v, ok := resData["burst_buffer"]; ok {
		if burstBuffer, ok := v.(string); ok {
			res.BurstBuffer = burstBuffer
		}
	}
	if v, ok := resData["comment"]; ok {
		if comment, ok := v.(string); ok {
			res.Comment = comment
		}
	}
	if v, ok := resData["tres"]; ok {
		if tres, ok := v.(string); ok {
			res.TRESStr = tres
		}
	}

	// Time fields - skip complex time parsing for now
	// These would need detailed mapping based on the actual API response structure
	if v, ok := resData["start_time"]; ok {
		_ = v // Skip complex time parsing
	}
	if v, ok := resData["end_time"]; ok {
		_ = v // Skip complex time parsing  
	}

	// Numeric fields
	if v, ok := resData["watts"]; ok {
		if watts, ok := v.(float64); ok {
			res.WattsTotal = int64(watts)
		}
	}
	if v, ok := resData["max_start_delay"]; ok {
		if delay, ok := v.(float64); ok {
			res.MaxStartDelay = int32(delay)
		}
	}
	if v, ok := resData["purge_completed_time"]; ok {
		if purgeTime, ok := v.(float64); ok {
			res.PurgeCompletedTime = int32(purgeTime)
		}
	}

	// Skip complex nested structures for now
	// These would need detailed mapping based on the actual API response structure

	return res, nil
}