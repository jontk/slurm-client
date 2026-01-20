// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ConvertTRESListToString converts a list of TRES entries to a comma-separated string
// This is a generic version that handles various TRES value types
func ConvertTRESListToString(tresList []struct {
	Type  *string     `json:"type,omitempty"`
	Value interface{} `json:"value,omitempty"`
}) string {
	var parts []string
	for _, tres := range tresList {
		if tres.Type != nil && tres.Value != nil {
			var valueStr string

			// Handle different value types
			switch v := tres.Value.(type) {
			case *int64:
				if v != nil {
					valueStr = strconv.FormatInt(*v, 10)
				}
			case int64:
				valueStr = strconv.FormatInt(v, 10)
			case *int32:
				if v != nil {
					valueStr = strconv.FormatInt(int64(*v), 10)
				}
			case int32:
				valueStr = strconv.FormatInt(int64(v), 10)
			case *int:
				if v != nil {
					valueStr = strconv.Itoa(*v)
				}
			case int:
				valueStr = strconv.Itoa(v)
			case string:
				valueStr = v
			case *string:
				if v != nil {
					valueStr = *v
				}
			default:
				// For complex types, use fmt.Sprintf
				valueStr = fmt.Sprintf("%v", tres.Value)
			}

			if valueStr != "" && valueStr != "<nil>" {
				parts = append(parts, fmt.Sprintf("%s=%s", *tres.Type, valueStr))
			}
		}
	}
	return strings.Join(parts, ",")
}

// ConvertTRESListToStringSimple converts a simple TRES list with int64 values
// This is for the v0.0.41 association converter which uses a simpler structure
func ConvertTRESListToStringSimple(tresList []struct {
	Type  *string `json:"type,omitempty"`
	Value *int64  `json:"value,omitempty"`
}) string {
	var parts []string
	for _, tres := range tresList {
		if tres.Type != nil && tres.Value != nil {
			parts = append(parts, fmt.Sprintf("%s=%d", *tres.Type, *tres.Value))
		}
	}
	return strings.Join(parts, ",")
}

// FormatDurationForSlurm formats a Go duration for Slurm time format (HH:MM:SS)
func FormatDurationForSlurm(d time.Duration) string {
	totalSeconds := int(d.Seconds())
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

// ParseNumberField parses a number field that could be various numeric types
func ParseNumberField(field interface{}) (int32, bool) {
	switch v := field.(type) {
	case int32:
		return v, true
	case *int32:
		if v != nil {
			return *v, true
		}
	case int64:
		return int32(v), true
	case *int64:
		if v != nil {
			return int32(*v), true
		}
	case int:
		return int32(v), true
	case *int:
		if v != nil {
			return int32(*v), true
		}
	case float64:
		return int32(v), true
	case *float64:
		if v != nil {
			return int32(*v), true
		}
	}
	return 0, false
}

// ParseNumberField64 parses a number field to int64
func ParseNumberField64(field interface{}) (int64, bool) {
	switch v := field.(type) {
	case int64:
		return v, true
	case *int64:
		if v != nil {
			return *v, true
		}
	case int32:
		return int64(v), true
	case *int32:
		if v != nil {
			return int64(*v), true
		}
	case int:
		return int64(v), true
	case *int:
		if v != nil {
			return int64(*v), true
		}
	case float64:
		return int64(v), true
	case *float64:
		if v != nil {
			return int64(*v), true
		}
	}
	return 0, false
}
