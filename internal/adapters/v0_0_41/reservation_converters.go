// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"fmt"
	"time"

	types "github.com/jontk/slurm-client/api"
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
			n := name
			res.Name = &n
		}
	}
	if v, ok := resData["node_list"]; ok {
		if nodeList, ok := v.(string); ok {
			nl := nodeList
			res.NodeList = &nl
		}
	}
	if v, ok := resData["node_count"]; ok {
		if nodeCount, ok := v.(float64); ok {
			nc := int32(nodeCount)
			res.NodeCount = &nc
		}
	}
	if v, ok := resData["core_count"]; ok {
		if coreCount, ok := v.(float64); ok {
			cc := int32(coreCount)
			res.CoreCount = &cc
		}
	}
	// String fields (stored as comma-separated strings)
	if v, ok := resData["accounts"]; ok {
		if accountsStr, ok := v.(string); ok && accountsStr != "" {
			res.Accounts = &accountsStr
		}
	}
	if v, ok := resData["users"]; ok {
		if usersStr, ok := v.(string); ok && usersStr != "" {
			res.Users = &usersStr
		}
	}
	if v, ok := resData["groups"]; ok {
		if groupsStr, ok := v.(string); ok && groupsStr != "" {
			res.Groups = &groupsStr
		}
	}
	if v, ok := resData["features"]; ok {
		if featuresStr, ok := v.(string); ok && featuresStr != "" {
			res.Features = &featuresStr
		}
	}
	// Optional string fields
	if v, ok := resData["partition"]; ok {
		if partition, ok := v.(string); ok {
			p := partition
			res.Partition = &p
		}
	}
	if v, ok := resData["burst_buffer"]; ok {
		if burstBuffer, ok := v.(string); ok {
			bb := burstBuffer
			res.BurstBuffer = &bb
		}
	}
	if v, ok := resData["tres"]; ok {
		if tres, ok := v.(string); ok {
			t := tres
			res.TRES = &t
		}
	}
	if v, ok := resData["licenses"]; ok {
		if licenses, ok := v.(string); ok {
			l := licenses
			res.Licenses = &l
		}
	}
	// Time fields
	if v, ok := resData["start_time"]; ok {
		if timeStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := timeStruct["number"].(float64); ok && number > 0 {
				res.StartTime = time.Unix(int64(number), 0)
			}
		}
	}
	if v, ok := resData["end_time"]; ok {
		if timeStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := timeStruct["number"].(float64); ok && number > 0 {
				res.EndTime = time.Unix(int64(number), 0)
			}
		}
	}
	// Numeric fields
	if v, ok := resData["watts"]; ok {
		if watts, ok := v.(float64); ok {
			w := uint32(watts)
			res.Watts = &w
		}
	}
	if v, ok := resData["max_start_delay"]; ok {
		if delay, ok := v.(float64); ok {
			d := int32(delay)
			res.MaxStartDelay = &d
		}
	}
	if v, ok := resData["purge_completed_time"]; ok {
		if purgeTime, ok := v.(float64); ok {
			pt := uint32(purgeTime)
			res.PurgeCompleted = &types.ReservationPurgeCompleted{
				Time: &pt,
			}
		}
	}
	// PurgeCompleted - also check for structured format
	if v, ok := resData["purge_completed"]; ok {
		if purgeData, ok := v.(map[string]interface{}); ok {
			res.PurgeCompleted = convertPurgeCompletedFromMap(purgeData)
		}
	}
	// Flags - convert using helper function pattern
	if v, ok := resData["flags"]; ok {
		if flags, ok := v.([]interface{}); ok {
			resFlags := make([]types.ReservationFlagsValue, 0, len(flags))
			for _, f := range flags {
				if flag, ok := f.(string); ok {
					resFlags = append(resFlags, types.ReservationFlagsValue(flag))
				}
			}
			res.Flags = resFlags
		}
	}
	// CoreSpecializations - convert core specialization structures
	if v, ok := resData["core_specializations"]; ok {
		if coreSpecs, ok := v.([]interface{}); ok {
			res.CoreSpecializations = convertCoreSpecsFromSlice(coreSpecs)
		}
	}
	return res, nil
}

// convertPurgeCompletedFromMap converts purge completed data from a map to common ReservationPurgeCompleted type.
func convertPurgeCompletedFromMap(data map[string]interface{}) *types.ReservationPurgeCompleted {
	if data == nil {
		return nil
	}
	result := &types.ReservationPurgeCompleted{}
	if timeData, ok := data["time"].(map[string]interface{}); ok {
		if number, ok := timeData["number"].(float64); ok {
			time := uint32(number)
			result.Time = &time
		}
	}
	return result
}

// convertCoreSpecsFromSlice converts core specialization data from a slice to common ReservationCoreSpec slice.
func convertCoreSpecsFromSlice(data []interface{}) []types.ReservationCoreSpec {
	if len(data) == 0 {
		return nil
	}
	result := make([]types.ReservationCoreSpec, 0, len(data))
	for _, item := range data {
		if spec, ok := item.(map[string]interface{}); ok {
			coreSpec := types.ReservationCoreSpec{}
			if core, ok := spec["core"].(string); ok {
				coreSpec.Core = &core
			}
			if node, ok := spec["node"].(string); ok {
				coreSpec.Node = &node
			}
			result = append(result, coreSpec)
		}
	}
	return result
}
