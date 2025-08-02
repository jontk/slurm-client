// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package builders

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQoSBuilder_Basic(t *testing.T) {
	t.Run("simple QoS", func(t *testing.T) {
		qos, err := NewQoSBuilder("test-qos").
			WithDescription("Test QoS").
			WithPriority(100).
			Build()

		require.NoError(t, err)
		assert.Equal(t, "test-qos", qos.Name)
		assert.Equal(t, "Test QoS", qos.Description)
		assert.Equal(t, 100, qos.Priority)
		assert.Equal(t, 1.0, qos.UsageFactor) // Default
	})

	t.Run("empty name fails", func(t *testing.T) {
		_, err := NewQoSBuilder("").Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("negative priority fails", func(t *testing.T) {
		_, err := NewQoSBuilder("test").
			WithPriority(-1).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be non-negative")
	})
}

func TestQoSBuilder_Limits(t *testing.T) {
	t.Run("comprehensive limits", func(t *testing.T) {
		qos, err := NewQoSBuilder("limited-qos").
			WithLimits().
				WithMaxCPUsPerUser(100).
				WithMaxJobsPerUser(10).
				WithMaxNodesPerUser(5).
				WithMaxWallTime(24 * time.Hour).
				WithMaxMemoryPerNode(64 * GB).
				WithMaxMemoryPerCPU(4 * GB).
				WithMinCPUsPerJob(2).
				Done().
			Build()

		require.NoError(t, err)
		require.NotNil(t, qos.Limits)
		assert.Equal(t, 100, *qos.Limits.MaxCPUsPerUser)
		assert.Equal(t, 10, *qos.Limits.MaxJobsPerUser)
		assert.Equal(t, 5, *qos.Limits.MaxNodesPerUser)
		assert.Equal(t, 1440, *qos.Limits.MaxWallTimePerJob) // 24 hours in minutes
		assert.Equal(t, int64(64*1024), *qos.Limits.MaxMemoryPerNode)
		assert.Equal(t, int64(4*1024), *qos.Limits.MaxMemoryPerCPU)
		assert.Equal(t, 2, *qos.Limits.MinCPUsPerJob)
	})

	t.Run("negative limits fail", func(t *testing.T) {
		_, err := NewQoSBuilder("test").
			WithLimits().
				WithMaxCPUsPerUser(-1).
				Done().
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be non-negative")
	})
}

func TestQoSBuilder_Presets(t *testing.T) {
	t.Run("high priority preset", func(t *testing.T) {
		qos, err := NewQoSBuilder("high-qos").
			AsHighPriority().
			Build()

		require.NoError(t, err)
		assert.Equal(t, 1000, qos.Priority)
		assert.Equal(t, 2.0, qos.UsageFactor)
		assert.Equal(t, 0.95, qos.UsageThreshold)
		assert.Contains(t, qos.Flags, "DenyOnLimit")
		assert.Contains(t, qos.Flags, "RequiresReservation")
		assert.Contains(t, qos.PreemptMode, "cluster")
	})

	t.Run("batch queue preset", func(t *testing.T) {
		qos, err := NewQoSBuilder("batch-qos").
			AsBatchQueue().
			Build()

		require.NoError(t, err)
		assert.Equal(t, 10, qos.Priority)
		assert.Equal(t, 0.5, qos.UsageFactor)
		assert.Contains(t, qos.Flags, "NoReserve")
		assert.Equal(t, 3600, qos.GraceTime)
	})

	t.Run("interactive preset", func(t *testing.T) {
		qos, err := NewQoSBuilder("interactive-qos").
			AsInteractive().
			Build()

		require.NoError(t, err)
		assert.Equal(t, 500, qos.Priority)
		assert.Equal(t, 1.5, qos.UsageFactor)
		assert.Contains(t, qos.Flags, "DenyOnLimit")
		assert.Contains(t, qos.PreemptMode, "suspend")
	})
}

func TestQoSBuilder_BusinessRules(t *testing.T) {
	t.Run("high priority requires reservation", func(t *testing.T) {
		_, err := NewQoSBuilder("test").
			WithPriority(1001).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RequiresReservation flag")
	})

	t.Run("high priority with reservation succeeds", func(t *testing.T) {
		qos, err := NewQoSBuilder("test").
			WithPriority(1001).
			WithFlags("RequiresReservation").
			Build()
		require.NoError(t, err)
		assert.Equal(t, 1001, qos.Priority)
	})

	t.Run("usage factor limit", func(t *testing.T) {
		_, err := NewQoSBuilder("test").
			WithUsageFactor(3.1).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot exceed 3.0")
	})

	t.Run("conflicting flags", func(t *testing.T) {
		_, err := NewQoSBuilder("test").
			WithFlags("NoReserve", "RequiresReservation").
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "conflicting flags")
	})
}

func TestQoSBuilder_Clone(t *testing.T) {
	original := NewQoSBuilder("original").
		WithDescription("Original description").
		WithPriority(100).
		WithFlags("flag1", "flag2").
		WithLimits().
			WithMaxCPUsPerUser(50).
			WithMaxJobsPerUser(5).
			Done()

	// Clone and modify
	cloned := original.Clone()
	cloned.WithDescription("Cloned description").
		WithPriority(200)

	// Build both
	originalQoS, err := original.Build()
	require.NoError(t, err)
	clonedQoS, err := cloned.Build()
	require.NoError(t, err)

	// Verify original is unchanged
	assert.Equal(t, "Original description", originalQoS.Description)
	assert.Equal(t, 100, originalQoS.Priority)

	// Verify clone has modifications
	assert.Equal(t, "Cloned description", clonedQoS.Description)
	assert.Equal(t, 200, clonedQoS.Priority)

	// Verify shared attributes are copied
	assert.Equal(t, originalQoS.Flags, clonedQoS.Flags)
	assert.Equal(t, *originalQoS.Limits.MaxCPUsPerUser, *clonedQoS.Limits.MaxCPUsPerUser)
}

func TestQoSBuilder_BuildForUpdate(t *testing.T) {
	t.Run("only modified fields", func(t *testing.T) {
		update, err := NewQoSBuilder("test").
			WithDescription("New description").
			WithPriority(200).
			BuildForUpdate()

		require.NoError(t, err)
		require.NotNil(t, update.Description)
		assert.Equal(t, "New description", *update.Description)
		require.NotNil(t, update.Priority)
		assert.Equal(t, 200, *update.Priority)

		// Default values should not be included
		assert.Nil(t, update.UsageFactor) // 1.0 is default
		assert.Nil(t, update.UsageThreshold) // 0 is default
	})

	t.Run("with limits", func(t *testing.T) {
		update, err := NewQoSBuilder("test").
			WithLimits().
				WithMaxCPUsPerUser(100).
				Done().
			BuildForUpdate()

		require.NoError(t, err)
		require.NotNil(t, update.Limits)
		assert.Equal(t, 100, *update.Limits.MaxCPUsPerUser)
	})
}

func TestQoSBuilder_Validation(t *testing.T) {
	tests := []struct {
		name    string
		builder func() *QoSBuilder
		errMsg  string
	}{
		{
			name: "negative usage factor",
			builder: func() *QoSBuilder {
				return NewQoSBuilder("test").WithUsageFactor(-1)
			},
			errMsg: "must be non-negative",
		},
		{
			name: "usage threshold out of range",
			builder: func() *QoSBuilder {
				return NewQoSBuilder("test").WithUsageThreshold(1.1)
			},
			errMsg: "must be between 0 and 1",
		},
		{
			name: "negative grace time",
			builder: func() *QoSBuilder {
				return NewQoSBuilder("test").WithGraceTime(-1)
			},
			errMsg: "must be non-negative",
		},
		{
			name: "negative preempt exempt time",
			builder: func() *QoSBuilder {
				return NewQoSBuilder("test").WithPreemptExemptTime(-1)
			},
			errMsg: "must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.builder().Build()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestQoSBuilder_ComplexScenario(t *testing.T) {
	// Build a complex QoS with all features
	qos, err := NewQoSBuilder("complex-qos").
		WithDescription("Complex QoS with all features").
		WithPriority(750).
		WithFlags("DenyOnLimit", "NoReserve").
		WithPreemptMode("cluster", "cancel").
		WithPreemptExemptTime(15).
		WithGraceTime(600).
		WithUsageFactor(1.75).
		WithUsageThreshold(0.85).
		WithLimits().
			WithMaxCPUsPerUser(200).
			WithMaxJobsPerUser(25).
			WithMaxNodesPerUser(10).
			WithMaxSubmitJobsPerUser(100).
			WithMaxCPUsPerAccount(1000).
			WithMaxJobsPerAccount(100).
			WithMaxNodesPerAccount(50).
			WithMaxCPUsPerJob(64).
			WithMaxNodesPerJob(4).
			WithMaxWallTime(48 * time.Hour).
			WithMaxMemoryPerNode(256 * GB).
			WithMaxMemoryPerCPU(8 * GB).
			WithMinCPUsPerJob(4).
			WithMinNodesPerJob(1).
			Done().
		Build()

	require.NoError(t, err)

	// Verify all fields
	assert.Equal(t, "complex-qos", qos.Name)
	assert.Equal(t, "Complex QoS with all features", qos.Description)
	assert.Equal(t, 750, qos.Priority)
	assert.Len(t, qos.Flags, 2)
	assert.Len(t, qos.PreemptMode, 2)
	assert.Equal(t, 15, *qos.PreemptExemptTime)
	assert.Equal(t, 600, qos.GraceTime)
	assert.Equal(t, 1.75, qos.UsageFactor)
	assert.Equal(t, 0.85, qos.UsageThreshold)

	// Verify limits
	require.NotNil(t, qos.Limits)
	assert.Equal(t, 200, *qos.Limits.MaxCPUsPerUser)
	assert.Equal(t, 25, *qos.Limits.MaxJobsPerUser)
	assert.Equal(t, 10, *qos.Limits.MaxNodesPerUser)
	assert.Equal(t, 100, *qos.Limits.MaxSubmitJobsPerUser)
	assert.Equal(t, 1000, *qos.Limits.MaxCPUsPerAccount)
	assert.Equal(t, 100, *qos.Limits.MaxJobsPerAccount)
	assert.Equal(t, 50, *qos.Limits.MaxNodesPerAccount)
	assert.Equal(t, 64, *qos.Limits.MaxCPUsPerJob)
	assert.Equal(t, 4, *qos.Limits.MaxNodesPerJob)
	assert.Equal(t, 2880, *qos.Limits.MaxWallTimePerJob) // 48 hours
	assert.Equal(t, int64(256*1024), *qos.Limits.MaxMemoryPerNode)
	assert.Equal(t, int64(8*1024), *qos.Limits.MaxMemoryPerCPU)
	assert.Equal(t, 4, *qos.Limits.MinCPUsPerJob)
	assert.Equal(t, 1, *qos.Limits.MinNodesPerJob)
}
