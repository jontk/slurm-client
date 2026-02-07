// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"strconv"
	"strings"
)

// ParseMemory converts memory specification strings to bytes.
// Supports formats like:
//   - "4096" (plain number, assumed to be MB)
//   - "4096M" or "4096m" (megabytes)
//   - "4G" or "4g" (gigabytes)
//   - "1T" or "1t" (terabytes)
//   - "1024K" or "1024k" (kilobytes)
//
// Returns the value in megabytes (MB) for SLURM compatibility.
func ParseMemory(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}

	// Check for unit suffix
	var multiplier int64 = 1 // Default: already in MB
	lastChar := s[len(s)-1]

	switch lastChar {
	case 'K', 'k':
		multiplier = 1 // KB -> MB (divide by 1024, but we'll handle after parse)
		s = s[:len(s)-1]
	case 'M', 'm':
		multiplier = 1 // MB stays as MB
		s = s[:len(s)-1]
	case 'G', 'g':
		multiplier = 1024 // GB -> MB
		s = s[:len(s)-1]
	case 'T', 't':
		multiplier = 1024 * 1024 // TB -> MB
		s = s[:len(s)-1]
	}

	// Parse the numeric part
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	// Handle KB case (divide by 1024)
	if lastChar == 'K' || lastChar == 'k' {
		return val / 1024, nil
	}

	return val * multiplier, nil
}

// ParseMemoryToBytes converts memory specification strings to bytes.
// Supports the same formats as ParseMemory but returns bytes.
func ParseMemoryToBytes(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}

	var multiplier uint64 = 1
	lastChar := s[len(s)-1]

	switch lastChar {
	case 'K', 'k':
		multiplier = 1024
		s = s[:len(s)-1]
	case 'M', 'm':
		multiplier = 1024 * 1024
		s = s[:len(s)-1]
	case 'G', 'g':
		multiplier = 1024 * 1024 * 1024
		s = s[:len(s)-1]
	case 'T', 't':
		multiplier = 1024 * 1024 * 1024 * 1024
		s = s[:len(s)-1]
		// default: No suffix, assume bytes (multiplier stays 1)
	}

	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return val * multiplier, nil
}
