package base

import (
	"testing"
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReservationBaseManager_New(t *testing.T) {
	manager := NewReservationBaseManager("v0.0.43")
	assert.NotNil(t, manager)
	assert.Equal(t, "v0.0.43", manager.GetVersion())
	assert.Equal(t, "Reservation", manager.GetResourceType())
}

func TestReservationBaseManager_ValidateReservationCreate(t *testing.T) {
	manager := NewReservationBaseManager("v0.0.43")

	now := time.Now()
	future := now.Add(24 * time.Hour)

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
			name: "empty start time",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: time.Time{},
			},
			wantErr: true,
			errMsg:  "start time is required",
		},
		{
			name: "past start time",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: now.Add(-1 * time.Hour),
				Duration:  3600, // 1 hour
			},
			wantErr: true,
			errMsg:  "start time must be in the future",
		},
		{
			name: "negative duration",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: future,
				Duration:  -1,
			},
			wantErr: true,
			errMsg:  "duration must be positive",
		},
		{
			name: "empty node list and count",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: future,
				Duration:  3600,
			},
			wantErr: true,
			errMsg:  "either nodes or node count must be specified",
		},
		{
			name: "negative node count",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: future,
				Duration:  3600,
				NodeCount: -1,
			},
			wantErr: true,
			errMsg:  "node count must be positive",
		},
		{
			name: "valid basic reservation with node count",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: future,
				Duration:  3600,
				NodeCount: 2,
			},
			wantErr: false,
		},
		{
			name: "valid basic reservation with node list",
			reservation: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: future,
				Duration:  3600,
				Nodes:     []string{"compute-01", "compute-02"},
			},
			wantErr: false,
		},
		{
			name: "valid complex reservation",
			reservation: &types.ReservationCreate{
				Name:        "complex-reservation",
				StartTime:   future,
				Duration:    7200, // 2 hours
				Nodes:       []string{"gpu-[01-04]"},
				Users:       []string{"user1", "user2"},
				Accounts:    []string{"account1"},
				Partition:   "gpu",
				Features:    []string{"gpu", "high-memory"},
				Flags:       []string{"MAINT", "IGNORE_JOBS"},
				Comment:     "Maintenance reservation",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateReservationCreate(tt.reservation)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestReservationBaseManager_ValidateReservationUpdate(t *testing.T) {
	manager := NewReservationBaseManager("v0.0.43")

	future := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name    string
		update  *types.ReservationUpdate
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil update",
			update:  nil,
			wantErr: true,
			errMsg:  "reservation update data is required",
		},
		{
			name: "empty name",
			update: &types.ReservationUpdate{
				Name: "",
			},
			wantErr: true,
			errMsg:  "reservation name is required",
		},
		{
			name: "negative duration",
			update: &types.ReservationUpdate{
				Name:     "test-reservation",
				Duration: intPtr(-1),
			},
			wantErr: true,
			errMsg:  "duration must be positive",
		},
		{
			name: "negative node count",
			update: &types.ReservationUpdate{
				Name:      "test-reservation",
				NodeCount: intPtr(-1),
			},
			wantErr: true,
			errMsg:  "node count must be positive",
		},
		{
			name: "valid update",
			update: &types.ReservationUpdate{
				Name:      "test-reservation",
				StartTime: &future,
				Duration:  intPtr(7200),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateReservationUpdate(tt.update)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestReservationBaseManager_ApplyReservationDefaults(t *testing.T) {
	manager := NewReservationBaseManager("v0.0.43")

	future := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name     string
		input    *types.ReservationCreate
		expected *types.ReservationCreate
	}{
		{
			name: "apply defaults to minimal reservation",
			input: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: future,
				Duration:  3600,
				NodeCount: 2,
			},
			expected: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: future,
				Duration:  3600,
				NodeCount: 2,
				Users:     []string{},
				Accounts:  []string{},
				Features:  []string{},
				Flags:     []string{},
				Comment:   "",
			},
		},
		{
			name: "preserve existing values",
			input: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: future,
				Duration:  7200,
				Nodes:     []string{"compute-[01-04]"},
				Users:     []string{"user1", "user2"},
				Accounts:  []string{"account1"},
				Partition: "compute",
				Features:  []string{"high-memory"},
				Flags:     []string{"MAINT"},
				Comment:   "Test reservation",
			},
			expected: &types.ReservationCreate{
				Name:      "test-reservation",
				StartTime: future,
				Duration:  7200,
				Nodes:     []string{"compute-[01-04]"},
				Users:     []string{"user1", "user2"},
				Accounts:  []string{"account1"},
				Partition: "compute",
				Features:  []string{"high-memory"},
				Flags:     []string{"MAINT"},
				Comment:   "Test reservation",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.ApplyReservationDefaults(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestReservationBaseManager_FilterReservationList(t *testing.T) {
	manager := NewReservationBaseManager("v0.0.43")

	now := time.Now()
	reservations := []types.Reservation{
		{
			Name:      "active-reservation",
			StartTime: now.Add(-1 * time.Hour),
			EndTime:   now.Add(1 * time.Hour),
			State:     "ACTIVE",
			Users:     []string{"user1", "user2"},
			Accounts:  []string{"account1"},
			Partition: "compute",
			Flags:     []string{"MAINT"},
		},
		{
			Name:      "future-reservation",
			StartTime: now.Add(2 * time.Hour),
			EndTime:   now.Add(4 * time.Hour),
			State:     "INACTIVE",
			Users:     []string{"user2", "user3"},
			Accounts:  []string{"account2"},
			Partition: "gpu",
			Flags:     []string{},
		},
		{
			Name:      "completed-reservation",
			StartTime: now.Add(-4 * time.Hour),
			EndTime:   now.Add(-2 * time.Hour),
			State:     "COMPLETED",
			Users:     []string{"user1"},
			Accounts:  []string{"account1"},
			Partition: "compute",
			Flags:     []string{"IGNORE_JOBS"},
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
			expected: []string{"active-reservation", "future-reservation", "completed-reservation"},
		},
		{
			name: "filter by names",
			opts: &types.ReservationListOptions{
				Names: []string{"active-reservation", "future-reservation"},
			},
			expected: []string{"active-reservation", "future-reservation"},
		},
		{
			name: "filter by state",
			opts: &types.ReservationListOptions{
				States: []string{"ACTIVE", "INACTIVE"},
			},
			expected: []string{"active-reservation", "future-reservation"},
		},
		{
			name: "filter by users",
			opts: &types.ReservationListOptions{
				Users: []string{"user1"},
			},
			expected: []string{"active-reservation", "completed-reservation"},
		},
		{
			name: "filter by accounts",
			opts: &types.ReservationListOptions{
				Accounts: []string{"account2"},
			},
			expected: []string{"future-reservation"},
		},
		{
			name: "filter by partition",
			opts: &types.ReservationListOptions{
				Partitions: []string{"gpu"},
			},
			expected: []string{"future-reservation"},
		},
		{
			name: "filter by flag",
			opts: &types.ReservationListOptions{
				WithFlags: []string{"MAINT"},
			},
			expected: []string{"active-reservation"},
		},
		{
			name: "filter active reservations",
			opts: &types.ReservationListOptions{
				ActiveOnly: boolPtr(true),
			},
			expected: []string{"active-reservation"},
		},
		{
			name: "combined filters",
			opts: &types.ReservationListOptions{
				Users:      []string{"user1"},
				Partitions: []string{"compute"},
			},
			expected: []string{"active-reservation", "completed-reservation"},
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
			result := manager.FilterReservationList(reservations, tt.opts)
			resultNames := make([]string, len(result))
			for i, reservation := range result {
				resultNames[i] = reservation.Name
			}
			assert.Equal(t, tt.expected, resultNames)
		})
	}
}

func TestReservationBaseManager_ValidateReservationName(t *testing.T) {
	manager := NewReservationBaseManager("v0.0.43")

	tests := []struct {
		name            string
		reservationName string
		wantErr         bool
		errMsg          string
	}{
		{
			name:            "valid reservation name",
			reservationName: "test-reservation",
			wantErr:         false,
		},
		{
			name:            "empty reservation name",
			reservationName: "",
			wantErr:         true,
			errMsg:          "reservation name is required",
		},
		{
			name:            "reservation name with underscores",
			reservationName: "test_reservation_123",
			wantErr:         false,
		},
		{
			name:            "reservation name with numbers",
			reservationName: "reservation123",
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateReservationName(tt.reservationName)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestReservationBaseManager_ValidateTimeRange(t *testing.T) {
	manager := NewReservationBaseManager("v0.0.43")

	now := time.Now()

	tests := []struct {
		name      string
		startTime time.Time
		duration  int
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid future time range",
			startTime: now.Add(1 * time.Hour),
			duration:  3600,
			wantErr:   false,
		},
		{
			name:      "past start time",
			startTime: now.Add(-1 * time.Hour),
			duration:  3600,
			wantErr:   true,
			errMsg:    "start time must be in the future",
		},
		{
			name:      "zero duration",
			startTime: now.Add(1 * time.Hour),
			duration:  0,
			wantErr:   true,
			errMsg:    "duration must be positive",
		},
		{
			name:      "negative duration",
			startTime: now.Add(1 * time.Hour),
			duration:  -3600,
			wantErr:   true,
			errMsg:    "duration must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateTimeRange(tt.startTime, tt.duration)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestReservationBaseManager_CheckReservationConflicts(t *testing.T) {
	manager := NewReservationBaseManager("v0.0.43")

	now := time.Now()
	existingReservations := []types.Reservation{
		{
			Name:      "existing-1",
			StartTime: now.Add(1 * time.Hour),
			EndTime:   now.Add(3 * time.Hour),
			Nodes:     []string{"compute-01", "compute-02"},
		},
		{
			Name:      "existing-2",
			StartTime: now.Add(4 * time.Hour),
			EndTime:   now.Add(6 * time.Hour),
			Nodes:     []string{"gpu-01"},
		},
	}

	tests := []struct {
		name        string
		startTime   time.Time
		endTime     time.Time
		nodes       []string
		hasConflict bool
		conflictMsg string
	}{
		{
			name:        "no conflict - different time",
			startTime:   now.Add(7 * time.Hour),
			endTime:     now.Add(9 * time.Hour),
			nodes:       []string{"compute-01"},
			hasConflict: false,
		},
		{
			name:        "no conflict - different nodes",
			startTime:   now.Add(2 * time.Hour),
			endTime:     now.Add(2.5 * time.Hour),
			nodes:       []string{"debug-01"},
			hasConflict: false,
		},
		{
			name:        "conflict - overlapping time and nodes",
			startTime:   now.Add(2 * time.Hour),
			endTime:     now.Add(2.5 * time.Hour),
			nodes:       []string{"compute-01"},
			hasConflict: true,
			conflictMsg: "conflicts with existing reservation",
		},
		{
			name:        "conflict - exact overlap",
			startTime:   now.Add(1 * time.Hour),
			endTime:     now.Add(3 * time.Hour),
			nodes:       []string{"compute-02"},
			hasConflict: true,
			conflictMsg: "conflicts with existing reservation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conflicts := manager.CheckReservationConflicts(tt.startTime, tt.endTime, tt.nodes, existingReservations)
			if tt.hasConflict {
				assert.NotEmpty(t, conflicts)
			} else {
				assert.Empty(t, conflicts)
			}
		})
	}
}

func TestReservationBaseManager_CalculateReservationStats(t *testing.T) {
	manager := NewReservationBaseManager("v0.0.43")

	now := time.Now()
	reservations := []types.Reservation{
		{
			Name:      "active-1",
			StartTime: now.Add(-1 * time.Hour),
			EndTime:   now.Add(1 * time.Hour),
			State:     "ACTIVE",
			NodeCount: 2,
		},
		{
			Name:      "active-2",
			StartTime: now.Add(-30 * time.Minute),
			EndTime:   now.Add(30 * time.Minute),
			State:     "ACTIVE",
			NodeCount: 4,
		},
		{
			Name:      "future-1",
			StartTime: now.Add(2 * time.Hour),
			EndTime:   now.Add(4 * time.Hour),
			State:     "INACTIVE",
			NodeCount: 1,
		},
		{
			Name:      "completed-1",
			StartTime: now.Add(-4 * time.Hour),
			EndTime:   now.Add(-2 * time.Hour),
			State:     "COMPLETED",
			NodeCount: 3,
		},
	}

	stats := manager.CalculateReservationStats(reservations)

	assert.Equal(t, 4, stats.TotalReservations)
	assert.Equal(t, 2, stats.ActiveReservations)
	assert.Equal(t, 1, stats.FutureReservations)
	assert.Equal(t, 1, stats.CompletedReservations)
	assert.Equal(t, 6, stats.ActiveNodes) // 2 + 4
	assert.Equal(t, 10, stats.TotalNodes) // 2 + 4 + 1 + 3
}

func TestReservationBaseManager_GetReservationState(t *testing.T) {
	manager := NewReservationBaseManager("v0.0.43")

	now := time.Now()

	tests := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		expected  string
	}{
		{
			name:      "future reservation",
			startTime: now.Add(1 * time.Hour),
			endTime:   now.Add(3 * time.Hour),
			expected:  "INACTIVE",
		},
		{
			name:      "active reservation",
			startTime: now.Add(-1 * time.Hour),
			endTime:   now.Add(1 * time.Hour),
			expected:  "ACTIVE",
		},
		{
			name:      "completed reservation",
			startTime: now.Add(-3 * time.Hour),
			endTime:   now.Add(-1 * time.Hour),
			expected:  "COMPLETED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := manager.GetReservationState(tt.startTime, tt.endTime)
			assert.Equal(t, tt.expected, state)
		})
	}
}

