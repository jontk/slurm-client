package base

import (
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQoSBaseManager(t *testing.T) {
	version := "v0.0.43"
	manager := NewQoSBaseManager(version)
	
	assert.NotNil(t, manager)
	assert.NotNil(t, manager.CRUDManager)
	assert.Equal(t, version, manager.GetVersion())
	assert.Equal(t, "QoS", manager.GetResourceType())
}

func TestQoSBaseManager_ValidateQoSCreate(t *testing.T) {
	manager := NewQoSBaseManager("v0.0.43")
	
	tests := []struct {
		name    string
		qos     *types.QoSCreate
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil QoS",
			qos:     nil,
			wantErr: true,
			errMsg:  "QoS data is required",
		},
		{
			name: "empty name",
			qos: &types.QoSCreate{
				Name: "",
			},
			wantErr: true,
			errMsg:  "QoS name is required",
		},
		{
			name: "invalid name with special chars",
			qos: &types.QoSCreate{
				Name: "qos@test",
			},
			wantErr: true,
			errMsg:  "contains invalid characters",
		},
		{
			name: "negative priority",
			qos: &types.QoSCreate{
				Name:     "test-qos",
				Priority: -1,
			},
			wantErr: true,
			errMsg:  "Priority must be non-negative",
		},
		{
			name: "negative usage factor",
			qos: &types.QoSCreate{
				Name:        "test-qos",
				UsageFactor: -1.0,
			},
			wantErr: true,
			errMsg:  "Usage factor must be non-negative",
		},
		{
			name: "usage threshold out of range (below)",
			qos: &types.QoSCreate{
				Name:           "test-qos",
				UsageThreshold: -0.1,
			},
			wantErr: true,
			errMsg:  "Usage threshold must be between 0 and 1",
		},
		{
			name: "usage threshold out of range (above)",
			qos: &types.QoSCreate{
				Name:           "test-qos",
				UsageThreshold: 1.1,
			},
			wantErr: true,
			errMsg:  "Usage threshold must be between 0 and 1",
		},
		{
			name: "negative grace time",
			qos: &types.QoSCreate{
				Name:      "test-qos",
				GraceTime: intPtr(-1),
			},
			wantErr: true,
			errMsg:  "Grace time must be non-negative",
		},
		{
			name: "negative preempt exempt time",
			qos: &types.QoSCreate{
				Name:               "test-qos",
				PreemptExemptTime: intPtr(-1),
			},
			wantErr: true,
			errMsg:  "Preempt exempt time must be non-negative",
		},
		{
			name: "valid QoS minimal",
			qos: &types.QoSCreate{
				Name: "test-qos",
			},
			wantErr: false,
		},
		{
			name: "valid QoS with all fields",
			qos: &types.QoSCreate{
				Name:               "test-qos",
				Description:        "Test QoS",
				Priority:           100,
				Flags:              []string{"DenyOnLimit"},
				PreemptMode:        []string{"cluster"},
				UsageFactor:        1.5,
				UsageThreshold:     0.8,
				GraceTime:          intPtr(300),
				PreemptExemptTime: intPtr(60),
			},
			wantErr: false,
		},
		{
			name: "valid QoS with limits",
			qos: &types.QoSCreate{
				Name: "test-qos",
				Limits: &types.QoSLimits{
					MaxCPUsPerUser:  intPtr(100),
					MaxJobsPerUser:  intPtr(10),
					MaxNodesPerUser: intPtr(5),
				},
			},
			wantErr: false,
		},
		{
			name: "invalid QoS with negative limits",
			qos: &types.QoSCreate{
				Name: "test-qos",
				Limits: &types.QoSLimits{
					MaxCPUsPerUser: intPtr(-1),
				},
			},
			wantErr: true,
			errMsg:  "MaxCPUsPerUser must be non-negative",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateQoSCreate(tt.qos)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestQoSBaseManager_ValidateQoSUpdate(t *testing.T) {
	manager := NewQoSBaseManager("v0.0.43")
	
	tests := []struct {
		name    string
		update  *types.QoSUpdate
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil update",
			update:  nil,
			wantErr: true,
			errMsg:  "QoS update data is required",
		},
		{
			name:    "empty update",
			update:  &types.QoSUpdate{},
			wantErr: false, // Empty updates are allowed
		},
		{
			name: "negative priority",
			update: &types.QoSUpdate{
				Priority: intPtr(-1),
			},
			wantErr: true,
			errMsg:  "Priority must be non-negative",
		},
		{
			name: "negative usage factor",
			update: &types.QoSUpdate{
				UsageFactor: float64Ptr(-1.0),
			},
			wantErr: true,
			errMsg:  "Usage factor must be non-negative",
		},
		{
			name: "usage threshold out of range",
			update: &types.QoSUpdate{
				UsageThreshold: float64Ptr(1.5),
			},
			wantErr: true,
			errMsg:  "Usage threshold must be between 0 and 1",
		},
		{
			name: "valid update with all fields",
			update: &types.QoSUpdate{
				Description:    stringPtr("Updated description"),
				Priority:       intPtr(200),
				Flags:          []string{"DenyOnLimit", "RequiresReservation"},
				PreemptMode:    stringPtr("suspend"),
				UsageFactor:    float64Ptr(2.0),
				UsageThreshold: float64Ptr(0.9),
			},
			wantErr: false,
		},
		{
			name: "valid update with limits",
			update: &types.QoSUpdate{
				Limits: &types.QoSLimits{
					MaxCPUsPerUser: intPtr(200),
					MaxJobsPerUser: intPtr(20),
				},
			},
			wantErr: false,
		},
		{
			name: "invalid update with negative limits",
			update: &types.QoSUpdate{
				Limits: &types.QoSLimits{
					MaxJobsPerUser: intPtr(-1),
				},
			},
			wantErr: true,
			errMsg:  "MaxJobsPerUser must be non-negative",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateQoSUpdate(tt.update)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestQoSBaseManager_ValidateQoSLimits(t *testing.T) {
	manager := NewQoSBaseManager("v0.0.43")
	
	tests := []struct {
		name    string
		limits  *types.QoSLimits
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil limits",
			limits:  nil,
			wantErr: false, // nil limits are allowed
		},
		{
			name:    "empty limits",
			limits:  &types.QoSLimits{},
			wantErr: false,
		},
		{
			name: "all valid limits",
			limits: &types.QoSLimits{
				MaxCPUsPerUser:        intPtr(100),
				MaxJobsPerUser:        intPtr(10),
				MaxNodesPerUser:       intPtr(5),
				MaxSubmitJobsPerUser:  intPtr(20),
				MaxCPUsPerAccount:     intPtr(1000),
				MaxJobsPerAccount:     intPtr(100),
				MaxNodesPerAccount:    intPtr(50),
				MaxCPUsPerJob:         intPtr(32),
				MaxNodesPerJob:        intPtr(2),
				MaxWallTimePerJob:     intPtr(1440),
				MaxMemoryPerNode:      int64Ptr(64000),
				MaxMemoryPerCPU:       int64Ptr(4000),
				MinCPUsPerJob:         intPtr(1),
				MinNodesPerJob:        intPtr(1),
			},
			wantErr: false,
		},
		{
			name: "negative MaxCPUsPerUser",
			limits: &types.QoSLimits{
				MaxCPUsPerUser: intPtr(-1),
			},
			wantErr: true,
			errMsg:  "MaxCPUsPerUser must be non-negative",
		},
		{
			name: "negative MaxJobsPerUser",
			limits: &types.QoSLimits{
				MaxJobsPerUser: intPtr(-10),
			},
			wantErr: true,
			errMsg:  "MaxJobsPerUser must be non-negative",
		},
		{
			name: "negative MaxNodesPerUser",
			limits: &types.QoSLimits{
				MaxNodesPerUser: intPtr(-5),
			},
			wantErr: true,
			errMsg:  "MaxNodesPerUser must be non-negative",
		},
		{
			name: "negative MaxWallTimePerJob",
			limits: &types.QoSLimits{
				MaxWallTimePerJob: intPtr(-60),
			},
			wantErr: true,
			errMsg:  "MaxWallTimePerJob must be non-negative",
		},
		{
			name: "negative MaxMemoryPerNode",
			limits: &types.QoSLimits{
				MaxMemoryPerNode: int64Ptr(-1000),
			},
			wantErr: true,
			errMsg:  "MaxMemoryPerNode must be non-negative",
		},
		{
			name: "negative MinCPUsPerJob",
			limits: &types.QoSLimits{
				MinCPUsPerJob: intPtr(-1),
			},
			wantErr: true,
			errMsg:  "MinCPUsPerJob must be non-negative",
		},
		{
			name: "multiple negative values",
			limits: &types.QoSLimits{
				MaxCPUsPerUser:  intPtr(-1),
				MaxJobsPerUser:  intPtr(-2),
				MaxNodesPerUser: intPtr(-3),
			},
			wantErr: true,
			errMsg:  "must be non-negative", // Should catch first error
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateQoSLimits(tt.limits)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestQoSBaseManager_ApplyQoSDefaults(t *testing.T) {
	manager := NewQoSBaseManager("v0.0.43")
	
	tests := []struct {
		name     string
		input    *types.QoSCreate
		expected *types.QoSCreate
	}{
		{
			name: "minimal QoS gets defaults",
			input: &types.QoSCreate{
				Name: "test-qos",
			},
			expected: &types.QoSCreate{
				Name:           "test-qos",
				Priority:       0,
				Flags:          []string{},
				PreemptMode:    []string{},
				UsageFactor:    1.0,
				UsageThreshold: 0,
			},
		},
		{
			name: "existing values are preserved",
			input: &types.QoSCreate{
				Name:           "test-qos",
				Priority:       100,
				Flags:          []string{"DenyOnLimit"},
				PreemptMode:    []string{"cluster"},
				UsageFactor:    2.0,
				UsageThreshold: 0.8,
			},
			expected: &types.QoSCreate{
				Name:           "test-qos",
				Priority:       100,
				Flags:          []string{"DenyOnLimit"},
				PreemptMode:    []string{"cluster"},
				UsageFactor:    2.0,
				UsageThreshold: 0.8,
			},
		},
		{
			name: "nil flags and preempt mode get empty slices",
			input: &types.QoSCreate{
				Name:        "test-qos",
				Priority:    50,
				Flags:       nil,
				PreemptMode: nil,
			},
			expected: &types.QoSCreate{
				Name:           "test-qos",
				Priority:       50,
				Flags:          []string{},
				PreemptMode:    []string{},
				UsageFactor:    1.0,
				UsageThreshold: 0,
			},
		},
		{
			name: "partial defaults",
			input: &types.QoSCreate{
				Name:        "test-qos",
				Priority:    75,
				UsageFactor: 1.5,
			},
			expected: &types.QoSCreate{
				Name:           "test-qos",
				Priority:       75,
				Flags:          []string{},
				PreemptMode:    []string{},
				UsageFactor:    1.5,
				UsageThreshold: 0,
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.ApplyQoSDefaults(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestQoSBaseManager_FilterQoSList(t *testing.T) {
	manager := NewQoSBaseManager("v0.0.43")
	
	// Test data
	qosList := []types.QoS{
		{
			Name:            "normal",
			AllowedAccounts: []string{"physics", "chemistry"},
			AllowedUsers:    []string{"user1", "user2"},
		},
		{
			Name:            "high",
			AllowedAccounts: []string{"physics", "biology"},
			AllowedUsers:    []string{"user3", "user4"},
		},
		{
			Name:            "low",
			AllowedAccounts: []string{"chemistry"},
			AllowedUsers:    []string{"user1", "user5"},
		},
		{
			Name:            "special",
			AllowedAccounts: []string{"admin"},
			AllowedUsers:    []string{"admin1", "admin2"},
		},
	}
	
	tests := []struct {
		name     string
		opts     *types.QoSListOptions
		expected []string // expected QoS names
	}{
		{
			name:     "nil options returns all",
			opts:     nil,
			expected: []string{"normal", "high", "low", "special"},
		},
		{
			name:     "empty options returns all",
			opts:     &types.QoSListOptions{},
			expected: []string{"normal", "high", "low", "special"},
		},
		{
			name: "filter by single name",
			opts: &types.QoSListOptions{
				Names: []string{"normal"},
			},
			expected: []string{"normal"},
		},
		{
			name: "filter by multiple names",
			opts: &types.QoSListOptions{
				Names: []string{"normal", "high"},
			},
			expected: []string{"normal", "high"},
		},
		{
			name: "filter by account",
			opts: &types.QoSListOptions{
				Accounts: []string{"physics"},
			},
			expected: []string{"normal", "high"},
		},
		{
			name: "filter by multiple accounts",
			opts: &types.QoSListOptions{
				Accounts: []string{"chemistry", "biology"},
			},
			expected: []string{"normal", "high", "low"},
		},
		{
			name: "filter by user",
			opts: &types.QoSListOptions{
				Users: []string{"user1"},
			},
			expected: []string{"normal", "low"},
		},
		{
			name: "filter by multiple users",
			opts: &types.QoSListOptions{
				Users: []string{"user3", "user5"},
			},
			expected: []string{"high", "low"},
		},
		{
			name: "combined filters (AND logic)",
			opts: &types.QoSListOptions{
				Accounts: []string{"physics"},
				Users:    []string{"user3"},
			},
			expected: []string{"high"},
		},
		{
			name: "no matches",
			opts: &types.QoSListOptions{
				Names: []string{"nonexistent"},
			},
			expected: []string{},
		},
		{
			name: "filter with names and accounts",
			opts: &types.QoSListOptions{
				Names:    []string{"normal", "high", "low"},
				Accounts: []string{"chemistry"},
			},
			expected: []string{"normal", "low"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.FilterQoSList(qosList, tt.opts)
			
			// Extract names from result
			resultNames := make([]string, len(result))
			for i, qos := range result {
				resultNames[i] = qos.Name
			}
			
			assert.Equal(t, tt.expected, resultNames)
		})
	}
}

// Helper functions for creating pointers
func intPtr(i int) *int {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}