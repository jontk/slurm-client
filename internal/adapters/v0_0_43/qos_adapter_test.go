package v0_0_43

import (
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQoSAdapter_ValidateQoSCreate(t *testing.T) {
	adapter := &QoSAdapter{
		QoSBaseManager: base.NewQoSBaseManager("v0.0.43"),
	}

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
			name: "negative priority",
			qos: &types.QoSCreate{
				Name:     "test-qos",
				Priority: -1,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "negative usage factor",
			qos: &types.QoSCreate{
				Name:        "test-qos",
				Priority:    0,
				UsageFactor: -1.0,
			},
			wantErr: true,
			errMsg:  "usage factor must be non-negative",
		},
		{
			name: "usage threshold out of range",
			qos: &types.QoSCreate{
				Name:           "test-qos",
				Priority:       0,
				UsageFactor:    1.0,
				UsageThreshold: 1.5,
			},
			wantErr: true,
			errMsg:  "usage threshold must be between 0 and 1",
		},
		{
			name: "valid QoS",
			qos: &types.QoSCreate{
				Name:           "test-qos",
				Description:    "Test QoS",
				Priority:       100,
				UsageFactor:    1.5,
				UsageThreshold: 0.8,
				Flags:          []string{"DenyOnLimit"},
			},
			wantErr: false,
		},
		{
			name: "valid QoS with limits",
			qos: &types.QoSCreate{
				Name:     "test-qos",
				Priority: 100,
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
				Name:     "test-qos",
				Priority: 100,
				Limits: &types.QoSLimits{
					MaxCPUsPerUser: intPtr(-1),
				},
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateQoSCreate(tt.qos)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestQoSAdapter_ApplyQoSDefaults(t *testing.T) {
	adapter := &QoSAdapter{
		QoSBaseManager: base.NewQoSBaseManager("v0.0.43"),
	}

	tests := []struct {
		name     string
		input    *types.QoSCreate
		expected *types.QoSCreate
	}{
		{
			name: "apply defaults to minimal QoS",
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
			name: "preserve existing values",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.ApplyQoSDefaults(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestQoSAdapter_FilterQoSList(t *testing.T) {
	adapter := &QoSAdapter{
		QoSBaseManager: base.NewQoSBaseManager("v0.0.43"),
	}

	qosList := []types.QoS{
		{
			Name:            "normal",
			AllowedAccounts: []string{"physics", "chemistry"},
			AllowedUsers:    []string{"user1", "user2"},
		},
		{
			Name:            "high",
			AllowedAccounts: []string{"physics"},
			AllowedUsers:    []string{"user3"},
		},
		{
			Name:            "low",
			AllowedAccounts: []string{"biology"},
			AllowedUsers:    []string{"user1", "user4"},
		},
	}

	tests := []struct {
		name     string
		opts     *types.QoSListOptions
		expected []string // expected QoS names
	}{
		{
			name:     "no filters",
			opts:     &types.QoSListOptions{},
			expected: []string{"normal", "high", "low"},
		},
		{
			name: "filter by name",
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
			name: "filter by user",
			opts: &types.QoSListOptions{
				Users: []string{"user1"},
			},
			expected: []string{"normal", "low"},
		},
		{
			name: "combined filters",
			opts: &types.QoSListOptions{
				Accounts: []string{"physics"},
				Users:    []string{"user3"},
			},
			expected: []string{"high"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.FilterQoSList(qosList, tt.opts)
			resultNames := make([]string, len(result))
			for i, qos := range result {
				resultNames[i] = qos.Name
			}
			assert.Equal(t, tt.expected, resultNames)
		})
	}
}

// Helper function
func intPtr(i int) *int {
	return &i
}