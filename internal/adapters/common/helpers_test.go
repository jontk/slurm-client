// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Tests for ConvertTRESListToString
func TestConvertTRESListToString(t *testing.T) {
	// Helper function to create string pointer
	strPtr := func(s string) *string { return &s }
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }
	intPtr := func(i int) *int { return &i }

	tests := []struct {
		name  string
		input []struct {
			Type  *string     `json:"type,omitempty"`
			Value interface{} `json:"value,omitempty"`
		}
		expected string
	}{
		{
			name: "empty list",
			input: []struct {
				Type  *string     `json:"type,omitempty"`
				Value interface{} `json:"value,omitempty"`
			}{},
			expected: "",
		},
		{
			name: "single entry with int64 pointer",
			input: []struct {
				Type  *string     `json:"type,omitempty"`
				Value interface{} `json:"value,omitempty"`
			}{
				{Type: strPtr("cpu"), Value: int64Ptr(4)},
			},
			expected: "cpu=4",
		},
		{
			name: "multiple entries with different types",
			input: []struct {
				Type  *string     `json:"type,omitempty"`
				Value interface{} `json:"value,omitempty"`
			}{
				{Type: strPtr("cpu"), Value: int64(8)},
				{Type: strPtr("mem"), Value: int32Ptr(1024)},
				{Type: strPtr("gpu"), Value: intPtr(2)},
			},
			expected: "cpu=8,mem=1024,gpu=2",
		},
		{
			name: "string values",
			input: []struct {
				Type  *string     `json:"type,omitempty"`
				Value interface{} `json:"value,omitempty"`
			}{
				{Type: strPtr("license"), Value: "feature1"},
				{Type: strPtr("name"), Value: strPtr("value1")},
			},
			expected: "license=feature1,name=value1",
		},
		{
			name: "nil values ignored",
			input: []struct {
				Type  *string     `json:"type,omitempty"`
				Value interface{} `json:"value,omitempty"`
			}{
				{Type: strPtr("cpu"), Value: int64Ptr(4)},
				{Type: strPtr("mem"), Value: nil},
				{Type: strPtr("gpu"), Value: int32Ptr(2)},
			},
			expected: "cpu=4,gpu=2",
		},
		{
			name: "nil type ignored",
			input: []struct {
				Type  *string     `json:"type,omitempty"`
				Value interface{} `json:"value,omitempty"`
			}{
				{Type: nil, Value: int64(4)},
				{Type: strPtr("mem"), Value: int64(1024)},
			},
			expected: "mem=1024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertTRESListToString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Tests for ConvertTRESListToStringSimple
func TestConvertTRESListToStringSimple(t *testing.T) {
	strPtr := func(s string) *string { return &s }
	int64Ptr := func(i int64) *int64 { return &i }

	tests := []struct {
		name  string
		input []struct {
			Type  *string `json:"type,omitempty"`
			Value *int64  `json:"value,omitempty"`
		}
		expected string
	}{
		{
			name: "empty list",
			input: []struct {
				Type  *string `json:"type,omitempty"`
				Value *int64  `json:"value,omitempty"`
			}{},
			expected: "",
		},
		{
			name: "single entry",
			input: []struct {
				Type  *string `json:"type,omitempty"`
				Value *int64  `json:"value,omitempty"`
			}{
				{Type: strPtr("cpu"), Value: int64Ptr(4)},
			},
			expected: "cpu=4",
		},
		{
			name: "multiple entries",
			input: []struct {
				Type  *string `json:"type,omitempty"`
				Value *int64  `json:"value,omitempty"`
			}{
				{Type: strPtr("cpu"), Value: int64Ptr(8)},
				{Type: strPtr("mem"), Value: int64Ptr(1024)},
				{Type: strPtr("gpu"), Value: int64Ptr(2)},
			},
			expected: "cpu=8,mem=1024,gpu=2",
		},
		{
			name: "nil values ignored",
			input: []struct {
				Type  *string `json:"type,omitempty"`
				Value *int64  `json:"value,omitempty"`
			}{
				{Type: strPtr("cpu"), Value: int64Ptr(4)},
				{Type: strPtr("mem"), Value: nil},
				{Type: strPtr("gpu"), Value: int64Ptr(2)},
			},
			expected: "cpu=4,gpu=2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertTRESListToStringSimple(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Tests for FormatDurationForSlurm
func TestFormatDurationForSlurm(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "zero duration",
			duration: 0,
			expected: "00:00:00",
		},
		{
			name:     "seconds only",
			duration: 45 * time.Second,
			expected: "00:00:45",
		},
		{
			name:     "minutes and seconds",
			duration: 5*time.Minute + 30*time.Second,
			expected: "00:05:30",
		},
		{
			name:     "hours, minutes, and seconds",
			duration: 2*time.Hour + 15*time.Minute + 45*time.Second,
			expected: "02:15:45",
		},
		{
			name:     "exactly one hour",
			duration: 1 * time.Hour,
			expected: "01:00:00",
		},
		{
			name:     "more than 24 hours",
			duration: 25*time.Hour + 30*time.Minute + 15*time.Second,
			expected: "25:30:15",
		},
		{
			name:     "complex duration",
			duration: 72*time.Hour + 45*time.Minute + 30*time.Second,
			expected: "72:45:30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDurationForSlurm(tt.duration)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Tests for ParseNumberField
func TestParseNumberField(t *testing.T) {
	int32Val := int32(42)
	int64Val := int64(100)
	intVal := 50
	float64Val := 75.5

	tests := []struct {
		name     string
		input    interface{}
		expected int32
		expectOK bool
	}{
		{
			name:     "int32 value",
			input:    int32(42),
			expected: 42,
			expectOK: true,
		},
		{
			name:     "int32 pointer",
			input:    &int32Val,
			expected: 42,
			expectOK: true,
		},
		{
			name:     "int64 value",
			input:    int64(100),
			expected: 100,
			expectOK: true,
		},
		{
			name:     "int64 pointer",
			input:    &int64Val,
			expected: 100,
			expectOK: true,
		},
		{
			name:     "int value",
			input:    50,
			expected: 50,
			expectOK: true,
		},
		{
			name:     "int pointer",
			input:    &intVal,
			expected: 50,
			expectOK: true,
		},
		{
			name:     "float64 value",
			input:    75.5,
			expected: 75,
			expectOK: true,
		},
		{
			name:     "float64 pointer",
			input:    &float64Val,
			expected: 75,
			expectOK: true,
		},
		{
			name:     "nil int32 pointer",
			input:    (*int32)(nil),
			expected: 0,
			expectOK: false,
		},
		{
			name:     "string value",
			input:    "not a number",
			expected: 0,
			expectOK: false,
		},
		{
			name:     "nil",
			input:    nil,
			expected: 0,
			expectOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := ParseNumberField(tt.input)
			assert.Equal(t, tt.expectOK, ok)
			if ok {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// Tests for ParseNumberField64
func TestParseNumberField64(t *testing.T) {
	int32Val := int32(42)
	int64Val := int64(100)
	intVal := 50
	float64Val := 75.5

	tests := []struct {
		name     string
		input    interface{}
		expected int64
		expectOK bool
	}{
		{
			name:     "int64 value",
			input:    int64(100),
			expected: 100,
			expectOK: true,
		},
		{
			name:     "int64 pointer",
			input:    &int64Val,
			expected: 100,
			expectOK: true,
		},
		{
			name:     "int32 value",
			input:    int32(42),
			expected: 42,
			expectOK: true,
		},
		{
			name:     "int32 pointer",
			input:    &int32Val,
			expected: 42,
			expectOK: true,
		},
		{
			name:     "int value",
			input:    50,
			expected: 50,
			expectOK: true,
		},
		{
			name:     "int pointer",
			input:    &intVal,
			expected: 50,
			expectOK: true,
		},
		{
			name:     "float64 value",
			input:    75.5,
			expected: 75,
			expectOK: true,
		},
		{
			name:     "float64 pointer",
			input:    &float64Val,
			expected: 75,
			expectOK: true,
		},
		{
			name:     "large int64",
			input:    int64(9223372036854775807), // max int64
			expected: 9223372036854775807,
			expectOK: true,
		},
		{
			name:     "nil int64 pointer",
			input:    (*int64)(nil),
			expected: 0,
			expectOK: false,
		},
		{
			name:     "string value",
			input:    "not a number",
			expected: 0,
			expectOK: false,
		},
		{
			name:     "nil",
			input:    nil,
			expected: 0,
			expectOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := ParseNumberField64(tt.input)
			assert.Equal(t, tt.expectOK, ok)
			if ok {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// Benchmark tests
func BenchmarkConvertTRESListToString(b *testing.B) {
	strPtr := func(s string) *string { return &s }
	int64Ptr := func(i int64) *int64 { return &i }

	tresList := []struct {
		Type  *string     `json:"type,omitempty"`
		Value interface{} `json:"value,omitempty"`
	}{
		{Type: strPtr("cpu"), Value: int64Ptr(8)},
		{Type: strPtr("mem"), Value: int64Ptr(1024)},
		{Type: strPtr("gpu"), Value: int64Ptr(2)},
	}

	b.ResetTimer()
	for range b.N {
		ConvertTRESListToString(tresList)
	}
}

func BenchmarkFormatDurationForSlurm(b *testing.B) {
	duration := 2*time.Hour + 15*time.Minute + 45*time.Second

	b.ResetTimer()
	for range b.N {
		FormatDurationForSlurm(duration)
	}
}

func BenchmarkParseNumberField(b *testing.B) {
	value := int64(12345)

	b.ResetTimer()
	for range b.N {
		ParseNumberField(value)
	}
}

func BenchmarkParseNumberField64(b *testing.B) {
	value := int64(12345)

	b.ResetTimer()
	for range b.N {
		ParseNumberField64(value)
	}
}
