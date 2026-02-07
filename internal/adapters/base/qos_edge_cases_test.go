// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package base

import (
	"testing"

	types "github.com/jontk/slurm-client/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQoSNameUniqueness(t *testing.T) {
	mgr := NewQoSBaseManager("test")
	existingQoS := []types.QoS{
		{Name: stringPtr("normal")},
		{Name: stringPtr("high")},
		{Name: stringPtr("low")},
	}

	tests := []struct {
		name      string
		qosName   string
		wantErr   bool
		errSubstr string
	}{
		{
			name:      "case-insensitive match uppercase",
			qosName:   "Normal",
			wantErr:   true,
			errSubstr: "already exists",
		},
		{
			name:      "case-insensitive match all caps",
			qosName:   "HIGH",
			wantErr:   true,
			errSubstr: "already exists",
		},
		{
			name:    "unique name",
			qosName: "debug",
			wantErr: false,
		},
		{
			name:    "unique similar name",
			qosName: "normal-batch",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mgr.ValidateQoSNameUniqueness(tt.qosName, existingQoS)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errSubstr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestQoSPreemptMode(t *testing.T) {
	mgr := NewQoSBaseManager("test")

	tests := []struct {
		name      string
		modes     []string
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "all valid modes",
			modes:   []string{"OFF", "CANCEL", "CHECKPOINT", "GANG", "REQUEUE", "SUSPEND"},
			wantErr: false,
		},
		{
			name:    "single valid mode",
			modes:   []string{"REQUEUE"},
			wantErr: false,
		},
		{
			name:    "empty modes",
			modes:   []string{},
			wantErr: false,
		},
		{
			name:      "invalid mode",
			modes:     []string{"OFF", "INVALID_MODE"},
			wantErr:   true,
			errSubstr: "Invalid preempt mode",
		},
		{
			name:    "case insensitive lowercase",
			modes:   []string{"off", "cancel"},
			wantErr: false,
		},
		{
			name:    "case insensitive mixed",
			modes:   []string{"Off", "Cancel", "GANG"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mgr.ValidateQoSPreemptMode(tt.modes)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errSubstr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestQoSFlags(t *testing.T) {
	mgr := NewQoSBaseManager("test")

	tests := []struct {
		name      string
		flags     []string
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "valid flags",
			flags:   []string{"DENY_LIMIT", "NO_DECAY", "PART_TIME_LIMIT"},
			wantErr: false,
		},
		{
			name:    "all valid flags",
			flags:   []string{"DENY_LIMIT", "ENFORCE_USAGE_THRESHOLD", "NO_DECAY", "NO_RESERVE", "OVER_PART_QOS", "PART_MAX_NODE", "PART_MIN_NODE", "PART_TIME_LIMIT", "RELATIVE_PRIORITY", "REQUIRES_RES", "USAGE_FACTOR_SAFE"},
			wantErr: false,
		},
		{
			name:    "empty flags",
			flags:   []string{},
			wantErr: false,
		},
		{
			name:      "invalid flag",
			flags:     []string{"DENY_LIMIT", "INVALID_FLAG"},
			wantErr:   true,
			errSubstr: "Invalid QoS flag",
		},
		{
			name:    "case insensitive",
			flags:   []string{"deny_limit", "No_Decay"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mgr.ValidateQoSFlags(tt.flags)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errSubstr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestQoSGracetime(t *testing.T) {
	mgr := NewQoSBaseManager("test")

	tests := []struct {
		name      string
		seconds   int
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "valid grace time 1 hour",
			seconds: 3600,
			wantErr: false,
		},
		{
			name:    "zero grace time",
			seconds: 0,
			wantErr: false,
		},
		{
			name:    "max valid grace time (7 days)",
			seconds: 604800,
			wantErr: false,
		},
		{
			name:      "negative grace time",
			seconds:   -1,
			wantErr:   true,
			errSubstr: "non-negative",
		},
		{
			name:      "exceeds max (> 7 days)",
			seconds:   604801,
			wantErr:   true,
			errSubstr: "cannot exceed 7 days",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mgr.ValidateQoSGracetime(tt.seconds)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errSubstr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestQoSHierarchy(t *testing.T) {
	mgr := NewQoSBaseManager("test")
	existingQoS := []types.QoS{
		{Name: stringPtr("normal")},
		{Name: stringPtr("high")},
		{Name: stringPtr("urgent")},
	}

	tests := []struct {
		name       string
		qosName    string
		parentName string
		wantErr    bool
		errSubstr  string
	}{
		{
			name:       "valid parent",
			qosName:    "new_qos",
			parentName: "normal",
			wantErr:    false,
		},
		{
			name:       "no parent (empty)",
			qosName:    "new_qos",
			parentName: "",
			wantErr:    false,
		},
		{
			name:       "self-reference",
			qosName:    "new_qos",
			parentName: "new_qos",
			wantErr:    true,
			errSubstr:  "cannot be its own parent",
		},
		{
			name:       "non-existent parent",
			qosName:    "new_qos",
			parentName: "nonexistent",
			wantErr:    true,
			errSubstr:  "does not exist",
		},
		{
			name:       "case-insensitive parent match",
			qosName:    "new_qos",
			parentName: "NORMAL",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mgr.ValidateQoSHierarchy(tt.qosName, tt.parentName, existingQoS)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errSubstr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTRESLimits(t *testing.T) {
	mgr := NewQoSBaseManager("test")

	tests := []struct {
		name      string
		tres      string
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "valid TRES string",
			tres:    "cpu=100,mem=4096,node=10",
			wantErr: false,
		},
		{
			name:    "empty string",
			tres:    "",
			wantErr: false,
		},
		{
			name:    "unlimited values",
			tres:    "cpu=unlimited,mem=-1,node=50",
			wantErr: false,
		},
		{
			name:    "single resource",
			tres:    "gres/gpu=4",
			wantErr: false,
		},
		{
			name:      "invalid format missing equals",
			tres:      "cpu100",
			wantErr:   true,
			errSubstr: "Invalid TRES format",
		},
		{
			name:      "empty resource name",
			tres:      "=100",
			wantErr:   true,
			errSubstr: "resource name cannot be empty",
		},
		{
			name:      "invalid value non-numeric",
			tres:      "cpu=abc",
			wantErr:   true,
			errSubstr: "Invalid TRES value",
		},
		{
			name:      "invalid negative value",
			tres:      "cpu=-5",
			wantErr:   true,
			errSubstr: "must be >= -1",
		},
		{
			name:    "whitespace handling",
			tres:    "cpu = 100 , mem = 4096",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mgr.ValidateTRESLimits(tt.tres)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errSubstr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestQoSLimitsConsistency(t *testing.T) {
	mgr := NewQoSBaseManager("test")

	tests := []struct {
		name      string
		limits    *types.QoSLimits
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "nil limits",
			limits:  nil,
			wantErr: false,
		},
		{
			name:    "empty limits",
			limits:  &types.QoSLimits{},
			wantErr: false,
		},
		{
			name: "consistent active jobs limits (user < account)",
			limits: &types.QoSLimits{
				Max: &types.QoSLimitsMax{
					Jobs: &types.QoSLimitsMaxJobs{
						ActiveJobs: &types.QoSLimitsMaxJobsActiveJobs{
							Per: &types.QoSLimitsMaxJobsActiveJobsPer{
								User:    uint32Ptr(50),
								Account: uint32Ptr(100),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "inconsistent active jobs limits (user > account)",
			limits: &types.QoSLimits{
				Max: &types.QoSLimitsMax{
					Jobs: &types.QoSLimitsMaxJobs{
						ActiveJobs: &types.QoSLimitsMaxJobsActiveJobs{
							Per: &types.QoSLimitsMaxJobsActiveJobsPer{
								User:    uint32Ptr(100),
								Account: uint32Ptr(50),
							},
						},
					},
				},
			},
			wantErr:   true,
			errSubstr: "MaxJobsPerUser cannot exceed MaxJobsPerAccount",
		},
		{
			name: "consistent submit jobs limits (user < account)",
			limits: &types.QoSLimits{
				Max: &types.QoSLimitsMax{
					Jobs: &types.QoSLimitsMaxJobs{
						Per: &types.QoSLimitsMaxJobsPer{
							User:    uint32Ptr(50),
							Account: uint32Ptr(100),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "inconsistent submit jobs limits (user > account)",
			limits: &types.QoSLimits{
				Max: &types.QoSLimitsMax{
					Jobs: &types.QoSLimitsMaxJobs{
						Per: &types.QoSLimitsMaxJobsPer{
							User:    uint32Ptr(100),
							Account: uint32Ptr(50),
						},
					},
				},
			},
			wantErr:   true,
			errSubstr: "MaxSubmitJobsPerUser cannot exceed MaxSubmitJobsPerAccount",
		},
		{
			name: "equal limits are valid",
			limits: &types.QoSLimits{
				Max: &types.QoSLimitsMax{
					Jobs: &types.QoSLimitsMaxJobs{
						ActiveJobs: &types.QoSLimitsMaxJobsActiveJobs{
							Per: &types.QoSLimitsMaxJobsActiveJobsPer{
								User:    uint32Ptr(50),
								Account: uint32Ptr(50),
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mgr.ValidateQoSLimitsConsistency(tt.limits)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errSubstr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestQoSDeletionSafety(t *testing.T) {
	mgr := NewQoSBaseManager("test")

	tests := []struct {
		name         string
		qosName      string
		associations []types.Association
		wantErr      bool
		errSubstr    string
	}{
		{
			name:         "unused QoS can be deleted",
			qosName:      "unused",
			associations: []types.Association{},
			wantErr:      false,
		},
		{
			name:    "QoS in use cannot be deleted",
			qosName: "high",
			associations: []types.Association{
				{
					User:    "user1",
					Account: stringPtr("account1"),
					Default: &types.AssociationDefault{
						QoS: stringPtr("high"),
					},
				},
			},
			wantErr:   true,
			errSubstr: "Cannot delete QoS",
		},
		{
			name:         "default normal QoS cannot be deleted",
			qosName:      "normal",
			associations: []types.Association{},
			wantErr:      true,
			errSubstr:    "Cannot delete the default 'normal' QoS",
		},
		{
			name:    "case-insensitive normal check",
			qosName: "NORMAL",
			associations: []types.Association{
				{
					User:    "user1",
					Account: stringPtr("account1"),
					Default: &types.AssociationDefault{
						QoS: stringPtr("other"),
					},
				},
			},
			wantErr:   true,
			errSubstr: "Cannot delete the default 'normal' QoS",
		},
		{
			name:    "different QoS not affected",
			qosName: "debug",
			associations: []types.Association{
				{
					User:    "user1",
					Account: stringPtr("account1"),
					Default: &types.AssociationDefault{
						QoS: stringPtr("high"),
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mgr.ValidateQoSDeletionSafety(tt.qosName, tt.associations)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errSubstr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestEnhanceQoSCreateWithDefaults(t *testing.T) {
	mgr := NewQoSBaseManager("test")

	tests := []struct {
		name           string
		qos            *types.QoSCreate
		checkPriority  *int
		checkPreempt   []string
		checkWallClock *uint32
		checkMaxJobs   *uint32
	}{
		{
			name: "high priority QoS gets high priority and REQUEUE",
			qos: &types.QoSCreate{
				Name: "high-priority",
			},
			checkPriority: intPtr(1000000),
			checkPreempt:  []string{"REQUEUE"},
		},
		{
			name: "urgent QoS gets high priority",
			qos: &types.QoSCreate{
				Name: "urgent-jobs",
			},
			checkPriority: intPtr(1000000),
		},
		{
			name: "low priority QoS gets low priority and OFF",
			qos: &types.QoSCreate{
				Name: "low-batch",
			},
			checkPriority: intPtr(1),
			checkPreempt:  []string{"OFF"},
		},
		{
			name: "background QoS gets low priority",
			qos: &types.QoSCreate{
				Name: "background-jobs",
			},
			checkPriority: intPtr(1),
		},
		{
			name: "debug QoS gets time and job limits",
			qos: &types.QoSCreate{
				Name: "debug-test",
			},
			checkWallClock: uint32Ptr(60), // 1 hour in minutes
			checkMaxJobs:   uint32Ptr(2),
		},
		{
			name: "long-running QoS gets extended wall clock",
			qos: &types.QoSCreate{
				Name: "long-running",
			},
			checkWallClock: uint32Ptr(10080), // 7 days in minutes
		},
		{
			name: "normal QoS gets base defaults only",
			qos: &types.QoSCreate{
				Name: "standard",
			},
			// No special enhancements, but should have base defaults
		},
		{
			name: "existing values are preserved",
			qos: &types.QoSCreate{
				Name:        "high-custom",
				Priority:    500, // Custom priority should be preserved
				PreemptMode: []string{"SUSPEND"},
			},
			checkPriority: intPtr(500), // Should keep custom priority
			checkPreempt:  []string{"SUSPEND"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enhanced := mgr.EnhanceQoSCreateWithDefaults(tt.qos)
			require.NotNil(t, enhanced)

			if tt.checkPriority != nil {
				assert.Equal(t, *tt.checkPriority, enhanced.Priority)
			}

			if tt.checkPreempt != nil {
				assert.Equal(t, tt.checkPreempt, enhanced.PreemptMode)
			}

			if tt.checkWallClock != nil {
				require.NotNil(t, enhanced.Limits)
				require.NotNil(t, enhanced.Limits.Max)
				require.NotNil(t, enhanced.Limits.Max.WallClock)
				require.NotNil(t, enhanced.Limits.Max.WallClock.Per)
				require.NotNil(t, enhanced.Limits.Max.WallClock.Per.Job)
				assert.Equal(t, *tt.checkWallClock, *enhanced.Limits.Max.WallClock.Per.Job)
			}

			if tt.checkMaxJobs != nil {
				require.NotNil(t, enhanced.Limits)
				require.NotNil(t, enhanced.Limits.Max)
				require.NotNil(t, enhanced.Limits.Max.Jobs)
				require.NotNil(t, enhanced.Limits.Max.Jobs.ActiveJobs)
				require.NotNil(t, enhanced.Limits.Max.Jobs.ActiveJobs.Per)
				require.NotNil(t, enhanced.Limits.Max.Jobs.ActiveJobs.Per.User)
				assert.Equal(t, *tt.checkMaxJobs, *enhanced.Limits.Max.Jobs.ActiveJobs.Per.User)
			}

			// Base defaults should always be applied
			assert.NotNil(t, enhanced.Flags, "Flags should be initialized")
			assert.NotNil(t, enhanced.PreemptMode, "PreemptMode should be initialized")
		})
	}
}

func TestQoSUpdateSafety(t *testing.T) {
	mgr := NewQoSBaseManager("test")

	tests := []struct {
		name       string
		currentQoS *types.QoS
		update     *types.QoSUpdate
		wantErr    bool
	}{
		{
			name:       "update with no changes",
			currentQoS: &types.QoS{Name: stringPtr("test")},
			update:     &types.QoSUpdate{},
			wantErr:    false,
		},
		{
			name: "update priority",
			currentQoS: &types.QoS{
				Name:     stringPtr("test"),
				Priority: uint32Ptr(100),
			},
			update: &types.QoSUpdate{
				Priority: intPtr(200),
			},
			wantErr: false,
		},
		{
			name: "update limits",
			currentQoS: &types.QoS{
				Name: stringPtr("test"),
				Limits: &types.QoSLimits{
					Max: &types.QoSLimitsMax{
						Jobs: &types.QoSLimitsMaxJobs{
							ActiveJobs: &types.QoSLimitsMaxJobsActiveJobs{
								Per: &types.QoSLimitsMaxJobsActiveJobsPer{
									User: uint32Ptr(10),
								},
							},
						},
					},
				},
			},
			update: &types.QoSUpdate{
				Limits: &types.QoSLimits{
					Max: &types.QoSLimitsMax{
						Jobs: &types.QoSLimitsMaxJobs{
							ActiveJobs: &types.QoSLimitsMaxJobsActiveJobs{
								Per: &types.QoSLimitsMaxJobsActiveJobsPer{
									User: uint32Ptr(5), // Reducing limit
								},
							},
						},
					},
				},
			},
			wantErr: false, // Currently doesn't error, just logs/tracks
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mgr.ValidateQoSUpdateSafety(tt.currentQoS, tt.update)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// intPtr creates a pointer to an int value
func intPtr(i int) *int {
	return &i
}
