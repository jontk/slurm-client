package v0_0_43

import (
	"testing"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReservationAdapter_ValidateReservationCreate(t *testing.T) {
	adapter := &ReservationAdapter{
		ReservationBaseManager: base.NewReservationBaseManager("v0.0.43"),
	}

	tests := []struct {
		name        string
		reservation *types.ReservationCreate
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "nil reservation",
			reservation: nil,
			wantErr:     true,
			errMsg:      "reservation data is required",
		},
		{
			name: "empty name",
			reservation: &types.ReservationCreate{
				Name: "",
			},
			wantErr: true,
			errMsg:  "reservation name is required",
		},
		{
			name: "empty nodes and node count",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: 1640995200, // 2022-01-01 00:00:00 UTC
				EndTime:   1641081600, // 2022-01-02 00:00:00 UTC
			},
			wantErr: true,
			errMsg:  "either nodes or node count must be specified",
		},
		{
			name: "invalid time range - end before start",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: 1641081600, // 2022-01-02 00:00:00 UTC
				EndTime:   1640995200, // 2022-01-01 00:00:00 UTC
				NodeCount: 2,
			},
			wantErr: true,
			errMsg:  "end time must be after start time",
		},
		{
			name: "negative node count",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: 1640995200,
				EndTime:   1641081600,
				NodeCount: -1,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "negative core count",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: 1640995200,
				EndTime:   1641081600,
				NodeCount: 2,
				CoreCount: -4,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "valid basic reservation with node count",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: 1640995200,
				EndTime:   1641081600,
				NodeCount: 2,
			},
			wantErr: false,
		},
		{
			name: "valid basic reservation with nodes",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: 1640995200,
				EndTime:   1641081600,
				Nodes:     []string{"node001", "node002"},
			},
			wantErr: false,
		},
		{
			name: "valid complex reservation",
			reservation: &types.ReservationCreate{
				Name:        "maintenance-reservation",
				StartTime:   1640995200,
				EndTime:     1641081600,
				Duration:    86400, // 24 hours
				NodeCount:   4,
				CoreCount:   16,
				Users:       []string{"admin", "maintenance"},
				Accounts:    []string{"admin"},
				Licenses:    []string{"matlab:2", "ansys:1"},
				Features:    []string{"gpu", "infiniband"},
				Flags:       []string{"MAINT", "IGNORE_JOBS"},
				Partition:   "maintenance",
				TRES:        "cpu=64,mem=256G,gres/gpu=4",
				BurstBuffer: "datawarp:100GB",
				Watts:       1000,
				MaxStartDelay: 3600, // 1 hour
			},
			wantErr: false,
		},
		{
			name: "invalid flag",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: 1640995200,
				EndTime:   1641081600,
				NodeCount: 2,
				Flags:     []string{"INVALID_FLAG"},
			},
			wantErr: true,
			errMsg:  "invalid reservation flag",
		},
		{
			name: "negative watts",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: 1640995200,
				EndTime:   1641081600,
				NodeCount: 2,
				Watts:     -500,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
		{
			name: "negative max start delay",
			reservation: &types.ReservationCreate{
				Name:          "test-reservation",
				StartTime:     1640995200,
				EndTime:       1641081600,
				NodeCount:     2,
				MaxStartDelay: -3600,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateReservationCreate(tt.reservation)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestReservationAdapter_ApplyReservationDefaults(t *testing.T) {
	adapter := &ReservationAdapter{
		ReservationBaseManager: base.NewReservationBaseManager("v0.0.43"),
	}

	tests := []struct {
		name     string
		input    *types.ReservationCreate
		expected *types.ReservationCreate
	}{
		{
			name: "apply defaults to minimal reservation",
			input: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: 1640995200,
				EndTime:   1641081600,
				NodeCount: 2,
			},
			expected: &types.ReservationCreate{
				Name:          "test-reservation",
				StartTime:     1640995200,
				EndTime:       1641081600,
				Duration:      0,                     // Calculated from start/end time
				NodeCount:     2,
				CoreCount:     0,                     // Unlimited
				Users:         []string{},            // Empty users
				Accounts:      []string{},            // Empty accounts
				Licenses:      []string{},            // Empty licenses
				Features:      []string{},            // Empty features
				Flags:         []string{},            // Empty flags
				Partition:     "",                    // No partition restriction
				TRES:          "",                    // No TRES restriction
				BurstBuffer:   "",                    // No burst buffer
				Watts:         0,                     // No power limit
				MaxStartDelay: 0,                     // No start delay
				Nodes:         []string{},            // Empty node list
			},
		},
		{
			name: "preserve existing values",
			input: &types.ReservationCreate{
				Name:        "complex-reservation",
				StartTime:   1640995200,
				EndTime:     1641081600,
				Duration:    86400,
				NodeCount:   4,
				CoreCount:   16,
				Users:       []string{"user1", "user2"},
				Accounts:    []string{"physics"},
				Licenses:    []string{"matlab:2"},
				Features:    []string{"gpu"},
				Flags:       []string{"MAINT"},
				Partition:   "compute",
				TRES:        "cpu=64,mem=256G",
				BurstBuffer: "datawarp:100GB",
				Watts:       2000,
				MaxStartDelay: 1800,
				Nodes:       []string{"node001", "node002"},
			},
			expected: &types.ReservationCreate{
				Name:        "complex-reservation",
				StartTime:   1640995200,
				EndTime:     1641081600,
				Duration:    86400,
				NodeCount:   4,
				CoreCount:   16,
				Users:       []string{"user1", "user2"},
				Accounts:    []string{"physics"},
				Licenses:    []string{"matlab:2"},
				Features:    []string{"gpu"},
				Flags:       []string{"MAINT"},
				Partition:   "compute",
				TRES:        "cpu=64,mem=256G",
				BurstBuffer: "datawarp:100GB",
				Watts:       2000,
				MaxStartDelay: 1800,
				Nodes:       []string{"node001", "node002"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.ApplyReservationDefaults(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestReservationAdapter_FilterReservationList(t *testing.T) {
	adapter := &ReservationAdapter{
		ReservationBaseManager: base.NewReservationBaseManager("v0.0.43"),
	}

	reservations := []types.Reservation{
		{
			Name:      "maint-reservation",
			StartTime: 1640995200,
			EndTime:   1641081600,
			NodeCount: 2,
			Users:     []string{"admin"},
			Accounts:  []string{"admin"},
			Features:  []string{"maintenance"},
			Flags:     []string{"MAINT"},
			Partition: "maintenance",
			State:     "ACTIVE",
		},
		{
			Name:      "gpu-reservation",
			StartTime: 1641168000,
			EndTime:   1641254400,
			NodeCount: 4,
			Users:     []string{"researcher1", "researcher2"},
			Accounts:  []string{"physics", "chemistry"},
			Features:  []string{"gpu", "cuda"},
			Flags:     []string{},
			Partition: "gpu",
			State:     "INACTIVE",
		},
		{
			Name:      "bigmem-reservation",
			StartTime: 1641340800,
			EndTime:   1641427200,
			NodeCount: 1,
			Users:     []string{"researcher3"},
			Accounts:  []string{"biology"},
			Features:  []string{"bigmem", "high-memory"},
			Flags:     []string{"SPEC_NODES"},
			Partition: "bigmem",
			State:     "ACTIVE",
			Nodes:     []string{"bigmem001"},
		},
	}

	tests := []struct {
		name     string
		opts     *types.ReservationListOptions
		expected []string // expected reservation names
	}{
		{
			name:     "no filters",
			opts:     &types.ReservationListOptions{},
			expected: []string{"maint-reservation", "gpu-reservation", "bigmem-reservation"},
		},
		{
			name: "filter by names",
			opts: &types.ReservationListOptions{
				Names: []string{"maint-reservation", "gpu-reservation"},
			},
			expected: []string{"maint-reservation", "gpu-reservation"},
		},
		{
			name: "filter by users",
			opts: &types.ReservationListOptions{
				Users: []string{"researcher1"},
			},
			expected: []string{"gpu-reservation"},
		},
		{
			name: "filter by accounts",
			opts: &types.ReservationListOptions{
				Accounts: []string{"physics"},
			},
			expected: []string{"gpu-reservation"},
		},
		{
			name: "filter by state",
			opts: &types.ReservationListOptions{
				States: []string{"ACTIVE"},
			},
			expected: []string{"maint-reservation", "bigmem-reservation"},
		},
		{
			name: "filter by features",
			opts: &types.ReservationListOptions{
				Features: []string{"gpu"},
			},
			expected: []string{"gpu-reservation"},
		},
		{
			name: "filter by partition",
			opts: &types.ReservationListOptions{
				Partitions: []string{"maintenance", "bigmem"},
			},
			expected: []string{"maint-reservation", "bigmem-reservation"},
		},
		{
			name: "filter by flags",
			opts: &types.ReservationListOptions{
				Flags: []string{"MAINT"},
			},
			expected: []string{"maint-reservation"},
		},
		{
			name: "filter by nodes",
			opts: &types.ReservationListOptions{
				Nodes: []string{"bigmem001"},
			},
			expected: []string{"bigmem-reservation"},
		},
		{
			name: "combined filters",
			opts: &types.ReservationListOptions{
				States:   []string{"ACTIVE"},
				Features: []string{"maintenance", "high-memory"},
			},
			expected: []string{"maint-reservation", "bigmem-reservation"},
		},
		{
			name: "no matches",
			opts: &types.ReservationListOptions{
				Names: []string{"nonexistent"},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.FilterReservationList(reservations, tt.opts)
			resultNames := make([]string, len(result))
			for i, reservation := range result {
				resultNames[i] = reservation.Name
			}
			assert.Equal(t, tt.expected, resultNames)
		})
	}
}

func TestReservationAdapter_ValidateTimeRange(t *testing.T) {
	adapter := &ReservationAdapter{
		ReservationBaseManager: base.NewReservationBaseManager("v0.0.43"),
	}

	tests := []struct {
		name      string
		startTime int64
		endTime   int64
		duration  int
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid time range",
			startTime: 1640995200,
			endTime:   1641081600,
			duration:  0,
			wantErr:   false,
		},
		{
			name:      "valid time range with duration",
			startTime: 1640995200,
			endTime:   0,
			duration:  86400,
			wantErr:   false,
		},
		{
			name:      "end time before start time",
			startTime: 1641081600,
			endTime:   1640995200,
			duration:  0,
			wantErr:   true,
			errMsg:    "end time must be after start time",
		},
		{
			name:      "start time in the past",
			startTime: 1000000000, // Way in the past
			endTime:   1641081600,
			duration:  0,
			wantErr:   true,
			errMsg:    "start time cannot be in the past",
		},
		{
			name:      "negative duration",
			startTime: 1640995200,
			endTime:   0,
			duration:  -3600,
			wantErr:   true,
			errMsg:    "duration must be positive",
		},
		{
			name:      "zero duration when no end time",
			startTime: 1640995200,
			endTime:   0,
			duration:  0,
			wantErr:   true,
			errMsg:    "either end time or duration must be specified",
		},
		{
			name:      "both end time and duration specified",
			startTime: 1640995200,
			endTime:   1641081600,
			duration:  86400,
			wantErr:   false, // This is valid - duration can be used for validation
		},
		{
			name:      "zero start time",
			startTime: 0,
			endTime:   1641081600,
			duration:  0,
			wantErr:   true,
			errMsg:    "start time is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateTimeRange(tt.startTime, tt.endTime, tt.duration)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestReservationAdapter_ValidateResourceRequirements(t *testing.T) {
	adapter := &ReservationAdapter{
		ReservationBaseManager: base.NewReservationBaseManager("v0.0.43"),
	}

	tests := []struct {
		name      string
		nodeCount int
		coreCount int
		nodes     []string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid node count",
			nodeCount: 4,
			coreCount: 16,
			nodes:     []string{},
			wantErr:   false,
		},
		{
			name:      "valid node list",
			nodeCount: 0,
			coreCount: 0,
			nodes:     []string{"node001", "node002"},
			wantErr:   false,
		},
		{
			name:      "both node count and node list",
			nodeCount: 2,
			coreCount: 8,
			nodes:     []string{"node001", "node002"},
			wantErr:   false, // This is valid - they should match
		},
		{
			name:      "no nodes or node count",
			nodeCount: 0,
			coreCount: 0,
			nodes:     []string{},
			wantErr:   true,
			errMsg:    "either nodes or node count must be specified",
		},
		{
			name:      "negative node count",
			nodeCount: -2,
			coreCount: 8,
			nodes:     []string{},
			wantErr:   true,
			errMsg:    "must be non-negative",
		},
		{
			name:      "negative core count",
			nodeCount: 2,
			coreCount: -8,
			nodes:     []string{},
			wantErr:   true,
			errMsg:    "must be non-negative",
		},
		{
			name:      "empty node name",
			nodeCount: 0,
			coreCount: 0,
			nodes:     []string{"node001", "", "node003"},
			wantErr:   true,
			errMsg:    "node name cannot be empty",
		},
		{
			name:      "duplicate node names",
			nodeCount: 0,
			coreCount: 0,
			nodes:     []string{"node001", "node001", "node002"},
			wantErr:   true,
			errMsg:    "duplicate node name",
		},
		{
			name:      "mismatched node count and list",
			nodeCount: 3,
			coreCount: 0,
			nodes:     []string{"node001", "node002"}, // Count is 3 but only 2 nodes
			wantErr:   true,
			errMsg:    "node count does not match node list length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateResourceRequirements(tt.nodeCount, tt.coreCount, tt.nodes)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestReservationAdapter_ValidateAccessControls(t *testing.T) {
	adapter := &ReservationAdapter{
		ReservationBaseManager: base.NewReservationBaseManager("v0.0.43"),
	}

	tests := []struct {
		name     string
		users    []string
		accounts []string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid users and accounts",
			users:    []string{"user1", "user2"},
			accounts: []string{"physics", "chemistry"},
			wantErr:  false,
		},
		{
			name:     "empty users and accounts (valid - no restrictions)",
			users:    []string{},
			accounts: []string{},
			wantErr:  false,
		},
		{
			name:     "users only",
			users:    []string{"user1", "user2"},
			accounts: []string{},
			wantErr:  false,
		},
		{
			name:     "accounts only",
			users:    []string{},
			accounts: []string{"physics", "chemistry"},
			wantErr:  false,
		},
		{
			name:     "empty user name",
			users:    []string{"user1", "", "user3"},
			accounts: []string{"physics"},
			wantErr:  true,
			errMsg:   "user name cannot be empty",
		},
		{
			name:     "duplicate user names",
			users:    []string{"user1", "user2", "user1"},
			accounts: []string{"physics"},
			wantErr:  true,
			errMsg:   "duplicate user name",
		},
		{
			name:     "empty account name",
			users:    []string{"user1"},
			accounts: []string{"physics", "", "chemistry"},
			wantErr:  true,
			errMsg:   "account name cannot be empty",
		},
		{
			name:     "duplicate account names",
			users:    []string{"user1"},
			accounts: []string{"physics", "chemistry", "physics"},
			wantErr:  true,
			errMsg:   "duplicate account name",
		},
		{
			name:     "user name with spaces",
			users:    []string{"user name"},
			accounts: []string{"physics"},
			wantErr:  true,
			errMsg:   "user name cannot contain spaces",
		},
		{
			name:     "account name with spaces",
			users:    []string{"user1"},
			accounts: []string{"physics account"},
			wantErr:  true,
			errMsg:   "account name cannot contain spaces",
		},
		{
			name:     "valid user and account names with allowed characters",
			users:    []string{"user-1", "user_2", "user.3"},
			accounts: []string{"physics-dept", "chemistry_lab"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateAccessControls(tt.users, tt.accounts)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestReservationAdapter_ValidateReservationFlags(t *testing.T) {
	adapter := &ReservationAdapter{
		ReservationBaseManager: base.NewReservationBaseManager("v0.0.43"),
	}

	tests := []struct {
		name    string
		flags   []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid flags",
			flags:   []string{"MAINT", "IGNORE_JOBS", "SPEC_NODES"},
			wantErr: false,
		},
		{
			name:    "empty flags",
			flags:   []string{},
			wantErr: false,
		},
		{
			name:    "single valid flag",
			flags:   []string{"MAINT"},
			wantErr: false,
		},
		{
			name:    "invalid flag",
			flags:   []string{"INVALID_FLAG"},
			wantErr: true,
			errMsg:  "invalid reservation flag",
		},
		{
			name:    "mixed valid and invalid flags",
			flags:   []string{"MAINT", "INVALID_FLAG", "IGNORE_JOBS"},
			wantErr: true,
			errMsg:  "invalid reservation flag",
		},
		{
			name:    "duplicate flags",
			flags:   []string{"MAINT", "MAINT", "IGNORE_JOBS"},
			wantErr: true,
			errMsg:  "duplicate reservation flag",
		},
		{
			name:    "empty flag string",
			flags:   []string{"MAINT", "", "IGNORE_JOBS"},
			wantErr: true,
			errMsg:  "reservation flag cannot be empty",
		},
		{
			name:    "lowercase flag",
			flags:   []string{"maint"},
			wantErr: true,
			errMsg:  "invalid reservation flag",
		},
		{
			name:    "all valid standard flags",
			flags:   []string{"MAINT", "IGNORE_JOBS", "SPEC_NODES", "PART_NODES", "FIRST_CORES", "TIME_FLOAT", "REPLACE", "ALL_NODES"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateReservationFlags(tt.flags)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestReservationAdapter_ValidateLicenses(t *testing.T) {
	adapter := &ReservationAdapter{
		ReservationBaseManager: base.NewReservationBaseManager("v0.0.43"),
	}

	tests := []struct {
		name     string
		licenses []string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid license specifications",
			licenses: []string{"matlab:2", "ansys:1", "fluent:4"},
			wantErr:  false,
		},
		{
			name:     "empty licenses",
			licenses: []string{},
			wantErr:  false,
		},
		{
			name:     "single license",
			licenses: []string{"matlab:1"},
			wantErr:  false,
		},
		{
			name:     "license without count (default to 1)",
			licenses: []string{"matlab"},
			wantErr:  false,
		},
		{
			name:     "invalid license format",
			licenses: []string{"matlab-2"},
			wantErr:  true,
			errMsg:   "invalid license format",
		},
		{
			name:     "license with zero count",
			licenses: []string{"matlab:0"},
			wantErr:  true,
			errMsg:   "license count must be positive",
		},
		{
			name:     "license with negative count",
			licenses: []string{"matlab:-2"},
			wantErr:  true,
			errMsg:   "license count must be positive",
		},
		{
			name:     "duplicate license types",
			licenses: []string{"matlab:2", "matlab:1"},
			wantErr:  true,
			errMsg:   "duplicate license type",
		},
		{
			name:     "empty license string",
			licenses: []string{"matlab:2", ""},
			wantErr:  true,
			errMsg:   "license cannot be empty",
		},
		{
			name:     "license with spaces in name",
			licenses: []string{"matlab software:2"},
			wantErr:  true,
			errMsg:   "license name cannot contain spaces",
		},
		{
			name:     "license with invalid characters",
			licenses: []string{"matlab@version:2"},
			wantErr:  true,
			errMsg:   "license name contains invalid characters",
		},
		{
			name:     "valid license names with allowed characters",
			licenses: []string{"matlab-r2021b:2", "ansys_fluent:1", "comsol.multiphysics:3"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.ValidateLicenses(tt.licenses)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}