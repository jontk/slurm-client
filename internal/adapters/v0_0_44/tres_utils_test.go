// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_44

import (
	"testing"

	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_44"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTRESUtils_ParseTRESString(t *testing.T) {
	utils := NewTRESUtils()
	tests := []struct {
		name     string
		input    string
		expected []types.TRES
		wantErr  bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []types.TRES{},
			wantErr:  false,
		},
		{
			name:  "single CPU",
			input: "cpu=4",
			expected: []types.TRES{
				{Type: "cpu", Count: ptrInt64(4)},
			},
			wantErr: false,
		},
		{
			name:  "multiple resources",
			input: "cpu=4,mem=8G,node=1",
			expected: []types.TRES{
				{Type: "cpu", Count: ptrInt64(4)},
				{Type: "mem", Count: ptrInt64(8 * 1024 * 1024 * 1024)},
				{Type: "node", Count: ptrInt64(1)},
			},
			wantErr: false,
		},
		{
			name:  "memory with different units",
			input: "mem=1024M,cpu=2",
			expected: []types.TRES{
				{Type: "mem", Count: ptrInt64(1024 * 1024 * 1024)},
				{Type: "cpu", Count: ptrInt64(2)},
			},
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "cpu4",
			wantErr: true,
		},
		{
			name:    "invalid count",
			input:   "cpu=abc",
			wantErr: true,
		},
		{
			name:  "with spaces",
			input: " cpu = 4 , mem = 8G ",
			expected: []types.TRES{
				{Type: "cpu", Count: ptrInt64(4)},
				{Type: "mem", Count: ptrInt64(8 * 1024 * 1024 * 1024)},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := utils.ParseTRESString(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
func TestTRESUtils_FormatTRESString(t *testing.T) {
	utils := NewTRESUtils()
	tests := []struct {
		name     string
		input    []types.TRES
		expected string
	}{
		{
			name:     "empty list",
			input:    []types.TRES{},
			expected: "",
		},
		{
			name: "single TRES",
			input: []types.TRES{
				{Type: "cpu", Count: ptrInt64(4)},
			},
			expected: "cpu=4",
		},
		{
			name: "multiple TRES",
			input: []types.TRES{
				{Type: "cpu", Count: ptrInt64(4)},
				{Type: "mem", Count: ptrInt64(8192)},
				{Type: "node", Count: ptrInt64(1)},
			},
			expected: "cpu=4,mem=8192,node=1",
		},
		{
			name: "TRES with names",
			input: []types.TRES{
				{Type: "gres", Name: ptrString("gpu"), Count: ptrInt64(2)},
			},
			expected: "gpu=2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.FormatTRESString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestTRESUtils_ConvertAPITRESToCommon(t *testing.T) {
	utils := NewTRESUtils()
	tests := []struct {
		name     string
		input    api.V0044TresList
		expected []types.TRES
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: []types.TRES{},
		},
		{
			name:     "empty list",
			input:    api.V0044TresList{},
			expected: []types.TRES{},
		},
		{
			name: "single TRES",
			input: api.V0044TresList{
				{
					Type:  "cpu",
					Id:    ptrInt32(1),
					Name:  ptrString("cpu"),
					Count: int64Ptr(4),
				},
			},
			expected: []types.TRES{
				{ID: ptrInt32(1), Type: "cpu", Name: ptrString("cpu"), Count: ptrInt64(4)},
			},
		},
		{
			name: "multiple TRES",
			input: api.V0044TresList{
				{
					Type:  "cpu",
					Id:    ptrInt32(1),
					Count: int64Ptr(4),
				},
				{
					Type:  "mem",
					Id:    ptrInt32(2),
					Count: int64Ptr(8192),
				},
			},
			expected: []types.TRES{
				{ID: ptrInt32(1), Type: "cpu", Count: ptrInt64(4)},
				{ID: ptrInt32(2), Type: "mem", Count: ptrInt64(8192)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ConvertAPITRESToCommon(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestTRESUtils_ConvertCommonTRESToAPI(t *testing.T) {
	utils := NewTRESUtils()
	tests := []struct {
		name     string
		input    []types.TRES
		expected api.V0044TresList
	}{
		{
			name:     "empty list",
			input:    []types.TRES{},
			expected: api.V0044TresList{},
		},
		{
			name: "single TRES",
			input: []types.TRES{
				{ID: ptrInt32(1), Type: "cpu", Name: ptrString("cpu"), Count: ptrInt64(4)},
			},
			expected: api.V0044TresList{
				{
					Type:  "cpu",
					Id:    ptrInt32(1),
					Name:  ptrString("cpu"),
					Count: int64Ptr(4),
				},
			},
		},
		{
			name: "TRES without optional fields",
			input: []types.TRES{
				{Type: "cpu", Count: ptrInt64(4)},
			},
			expected: api.V0044TresList{
				{
					Type:  "cpu",
					Count: int64Ptr(4),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ConvertCommonTRESToAPI(tt.input)
			assert.Equal(t, len(tt.expected), len(result))
			for i, expected := range tt.expected {
				if i >= len(result) {
					continue
				}
				actual := result[i]
				assert.Equal(t, expected.Type, actual.Type)
				assert.Equal(t, expected.Id, actual.Id)
				assert.Equal(t, expected.Name, actual.Name)
				assert.Equal(t, expected.Count, actual.Count)
			}
		})
	}
}
func TestTRESUtils_ExtractTRESByType(t *testing.T) {
	utils := NewTRESUtils()
	tresList := []types.TRES{
		{Type: "cpu", Count: ptrInt64(4)},
		{Type: "mem", Count: ptrInt64(8192)},
		{Type: "node", Count: ptrInt64(1)},
	}
	tests := []struct {
		name     string
		tresType string
		expected *types.TRES
	}{
		{
			name:     "found CPU",
			tresType: "cpu",
			expected: &types.TRES{Type: "cpu", Count: ptrInt64(4)},
		},
		{
			name:     "found memory case insensitive",
			tresType: "MEM",
			expected: &types.TRES{Type: "mem", Count: ptrInt64(8192)},
		},
		{
			name:     "not found",
			tresType: "gpu",
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ExtractTRESByType(tresList, tt.tresType)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			}
		})
	}
}
func TestTRESUtils_ExtractResourceLimits(t *testing.T) {
	utils := NewTRESUtils()
	tresList := []types.TRES{
		{Type: "cpu", Count: ptrInt64(8)},
		{Type: "mem", Count: ptrInt64(16384)},
		{Type: "node", Count: ptrInt64(2)},
	}
	cpus, memory, nodes := utils.ExtractResourceLimits(tresList)
	assert.Equal(t, int64(8), cpus)
	assert.Equal(t, int64(16384), memory)
	assert.Equal(t, int64(2), nodes)
}
func TestTRESUtils_BuildTRESFromLimits(t *testing.T) {
	utils := NewTRESUtils()
	result := utils.BuildTRESFromLimits(4, 8192, 1)
	expected := []types.TRES{
		{Type: "cpu", Count: ptrInt64(4)},
		{Type: "mem", Count: ptrInt64(8192)},
		{Type: "node", Count: ptrInt64(1)},
	}
	assert.Equal(t, expected, result)
}
func TestTRESUtils_ValidateTRES(t *testing.T) {
	utils := NewTRESUtils()
	tests := []struct {
		name    string
		tres    types.TRES
		wantErr bool
	}{
		{
			name:    "valid TRES",
			tres:    types.TRES{Type: "cpu", Count: ptrInt64(4)},
			wantErr: false,
		},
		{
			name:    "empty type",
			tres:    types.TRES{Type: "", Count: ptrInt64(4)},
			wantErr: true,
		},
		{
			name:    "negative count",
			tres:    types.TRES{Type: "cpu", Count: ptrInt64(-1)},
			wantErr: true,
		},
		{
			name:    "zero count is valid",
			tres:    types.TRES{Type: "cpu", Count: ptrInt64(0)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateTRES(tt.tres)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
func TestTRESUtils_MergeTRESLists(t *testing.T) {
	utils := NewTRESUtils()
	list1 := []types.TRES{
		{Type: "cpu", Count: ptrInt64(4)},
		{Type: "mem", Count: ptrInt64(8192)},
	}
	list2 := []types.TRES{
		{Type: "cpu", Count: ptrInt64(8)},  // Override
		{Type: "node", Count: ptrInt64(1)}, // New
	}
	result := utils.MergeTRESLists(list1, list2)
	// Should have 3 entries: cpu (from list2), mem (from list1), node (from list2)
	assert.Equal(t, 3, len(result))
	cpuTres := utils.ExtractTRESByType(result, "cpu")
	require.NotNil(t, cpuTres)
	require.NotNil(t, cpuTres.Count)
	assert.Equal(t, int64(8), *cpuTres.Count) // Should be the overridden value
}

// Helper function for tests (ptrString is in update_converters.go, ptrInt64 is in test_helpers.go)
func int64Ptr(i int64) *int64 {
	return &i
}
