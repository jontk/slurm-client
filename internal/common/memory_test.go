// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"testing"
)

func TestParseMemory(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
		wantErr  bool
	}{
		{"plain number", "4096", 4096, false},
		{"with M suffix", "4096M", 4096, false},
		{"with m suffix", "4096m", 4096, false},
		{"with G suffix", "4G", 4096, false},
		{"with g suffix", "4g", 4096, false},
		{"with T suffix", "1T", 1048576, false}, // 1TB = 1024*1024 MB
		{"with K suffix", "1024K", 1, false},    // 1024KB = 1 MB
		{"empty string", "", 0, false},
		{"whitespace", "  4096  ", 4096, false},
		{"invalid", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMemory(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMemory(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("ParseMemory(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseMemoryToBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected uint64
		wantErr  bool
	}{
		{"plain number", "1024", 1024, false},
		{"with K suffix", "1K", 1024, false},
		{"with M suffix", "1M", 1048576, false},
		{"with G suffix", "1G", 1073741824, false},
		{"with T suffix", "1T", 1099511627776, false},
		{"empty string", "", 0, false},
		{"invalid", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMemoryToBytes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMemoryToBytes(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("ParseMemoryToBytes(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
