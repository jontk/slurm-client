package builders

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/jontk/slurm-client/internal/common/types"
)

func TestPartitionBuilder_Basic(t *testing.T) {
	t.Run("simple partition", func(t *testing.T) {
		partition, err := NewPartitionBuilder("test-partition").
			WithNodes("node[001-010]").
			WithMaxNodes(10).
			WithMaxTime(480).
			Build()

		require.NoError(t, err)
		assert.Equal(t, "test-partition", partition.Name)
		assert.Equal(t, "node[001-010]", partition.Nodes)
		assert.Equal(t, int32(10), partition.MaxNodes)
		assert.Equal(t, int32(480), partition.MaxTime)
		assert.Equal(t, types.PartitionStateUp, partition.State) // Default
		assert.Equal(t, int32(1), partition.Priority)            // Default
	})

	t.Run("empty name fails", func(t *testing.T) {
		_, err := NewPartitionBuilder("").Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("negative values fail", func(t *testing.T) {
		_, err := NewPartitionBuilder("test").
			WithMaxNodes(-1).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be positive")
	})
}

func TestPartitionBuilder_Access(t *testing.T) {
	t.Run("account and group access", func(t *testing.T) {
		partition, err := NewPartitionBuilder("restricted").
			WithAllowAccounts("research", "compute").
			WithAllowGroups("users", "faculty").
			WithDenyAccounts("guest").
			WithAllowQoS("normal", "high").
			WithDenyQoS("debug").
			Build()

		require.NoError(t, err)
		assert.Equal(t, []string{"research", "compute"}, partition.AllowAccounts)
		assert.Equal(t, []string{"users", "faculty"}, partition.AllowGroups)
		assert.Equal(t, []string{"guest"}, partition.DenyAccounts)
		assert.Equal(t, []string{"normal", "high"}, partition.AllowQoS)
		assert.Equal(t, []string{"debug"}, partition.DenyQoS)
	})
}

func TestPartitionBuilder_Memory(t *testing.T) {
	t.Run("memory settings", func(t *testing.T) {
		partition, err := NewPartitionBuilder("highmem").
			WithDefaultMemPerCPU(4 * GB).
			WithDefaultMemPerNode(128 * GB).
			WithMaxMemPerNode(512 * GB).
			WithMaxMemPerCPU(8 * GB).
			WithDefMemPerNode(64 * GB).
			Build()

		require.NoError(t, err)
		assert.Equal(t, int64(4*GB), partition.DefaultMemPerCPU)
		assert.Equal(t, int64(128*GB), partition.DefaultMemPerNode)
		assert.Equal(t, int64(512*GB), partition.MaxMemPerNode)
		assert.Equal(t, int64(8*GB), partition.MaxMemPerCPU)
		assert.Equal(t, int64(64*GB), partition.DefMemPerNode)
	})

	t.Run("negative memory fails", func(t *testing.T) {
		_, err := NewPartitionBuilder("test").
			WithDefaultMemPerCPU(-1).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be non-negative")
	})
}

func TestPartitionBuilder_Time(t *testing.T) {
	t.Run("time limits", func(t *testing.T) {
		partition, err := NewPartitionBuilder("timed").
			WithDefaultTime(60).    // 1 hour
			WithMaxTime(1440).      // 24 hours
			WithGraceTime(300).     // 5 minutes
			WithOverTimeLimit(60).  // 1 hour overtime
			Build()

		require.NoError(t, err)
		assert.Equal(t, int32(60), partition.DefaultTime)
		assert.Equal(t, int32(1440), partition.MaxTime)
		assert.Equal(t, int32(300), partition.GraceTime)
		assert.Equal(t, int32(60), partition.OverTimeLimit)
	})
}

func TestPartitionBuilder_Priority(t *testing.T) {
	t.Run("priority settings", func(t *testing.T) {
		partition, err := NewPartitionBuilder("priority").
			WithPriority(100).
			WithPriorityJobFactor(50).
			WithPriorityTier(10).
			Build()

		require.NoError(t, err)
		assert.Equal(t, int32(100), partition.Priority)
		assert.Equal(t, int32(50), partition.PriorityJobFactor)
		assert.Equal(t, int32(10), partition.PriorityTier)
	})
}

func TestPartitionBuilder_Advanced(t *testing.T) {
	t.Run("advanced settings", func(t *testing.T) {
		jobDefaults := map[string]string{
			"mem":  "4G",
			"time": "1:00:00",
		}

		partition, err := NewPartitionBuilder("advanced").
			WithQoS("normal").
			WithTresStr("cpu=1000,mem=1000G").
			WithBillingWeightStr("cpu=1.0,mem=0.25").
			WithSelectTypeParameters("CR_Core", "CR_Memory").
			WithJobDefaults(jobDefaults).
			WithJobDefault("cpus-per-task", "1").
			WithResumeTimeout(600).
			WithSuspendTime(300).
			WithSuspendTimeout(180).
			WithPreemptMode("suspend", "cancel").
			Build()

		require.NoError(t, err)
		assert.Equal(t, "normal", partition.QoS)
		assert.Equal(t, "cpu=1000,mem=1000G", partition.TresStr)
		assert.Equal(t, "cpu=1.0,mem=0.25", partition.BillingWeightStr)
		assert.Equal(t, []string{"CR_Core", "CR_Memory"}, partition.SelectTypeParameters)
		assert.Equal(t, "4G", partition.JobDefaults["mem"])
		assert.Equal(t, "1:00:00", partition.JobDefaults["time"])
		assert.Equal(t, "1", partition.JobDefaults["cpus-per-task"])
		assert.Equal(t, int32(600), partition.ResumeTimeout)
		assert.Equal(t, int32(300), partition.SuspendTime)
		assert.Equal(t, int32(180), partition.SuspendTimeout)
		assert.Equal(t, []string{"suspend", "cancel"}, partition.PreemptMode)
	})
}

func TestPartitionBuilder_Flags(t *testing.T) {
	t.Run("partition flags", func(t *testing.T) {
		partition, err := NewPartitionBuilder("special").
			AsHidden().
			AsExclusiveUser().
			AsLLN().
			AsRootOnly().
			AsReqResv().
			AsPowerDownOnIdle().
			Build()

		require.NoError(t, err)
		assert.True(t, partition.Hidden)
		assert.True(t, partition.ExclusiveUser)
		assert.True(t, partition.LLN)
		assert.True(t, partition.RootOnly)
		assert.True(t, partition.ReqResv)
		assert.True(t, partition.PowerDownOnIdle)
	})
}

func TestPartitionBuilder_Presets(t *testing.T) {
	t.Run("debug partition preset", func(t *testing.T) {
		partition, err := NewPartitionBuilder("debug").
			AsDebugPartition().
			Build()

		require.NoError(t, err)
		assert.Equal(t, int32(1000), partition.Priority)       // High priority
		assert.Equal(t, int32(30), partition.MaxTime)          // 30 minutes max
		assert.Equal(t, int32(10), partition.DefaultTime)      // 10 minutes default
		assert.Equal(t, int32(2), partition.MaxNodes)          // Limited nodes
		assert.Contains(t, partition.PreemptMode, "cancel")    // Cancel jobs
		assert.Equal(t, "debug", partition.QoS)
	})

	t.Run("batch partition preset", func(t *testing.T) {
		partition, err := NewPartitionBuilder("batch").
			AsBatchPartition().
			Build()

		require.NoError(t, err)
		assert.Equal(t, int32(50), partition.Priority)         // Medium priority
		assert.Equal(t, int32(2880), partition.MaxTime)        // 48 hours max
		assert.Equal(t, int32(60), partition.DefaultTime)      // 1 hour default
		assert.Equal(t, int32(300), partition.GraceTime)       // 5 minutes grace
		assert.Contains(t, partition.PreemptMode, "suspend")   // Suspend jobs
		assert.Equal(t, "normal", partition.QoS)
	})

	t.Run("interactive partition preset", func(t *testing.T) {
		partition, err := NewPartitionBuilder("interactive").
			AsInteractivePartition().
			Build()

		require.NoError(t, err)
		assert.Equal(t, int32(500), partition.Priority)        // High priority
		assert.Equal(t, int32(240), partition.MaxTime)         // 4 hours max
		assert.Equal(t, int32(30), partition.DefaultTime)      // 30 minutes default
		assert.Equal(t, int32(4), partition.MaxNodes)          // Limited nodes
		assert.Contains(t, partition.PreemptMode, "cancel")    // Quick cancellation
		assert.Equal(t, "interactive", partition.QoS)
	})

	t.Run("GPU partition preset", func(t *testing.T) {
		partition, err := NewPartitionBuilder("gpu").
			AsGPUPartition().
			Build()

		require.NoError(t, err)
		assert.Equal(t, int32(100), partition.Priority)        // Medium-high priority
		assert.Equal(t, int32(1440), partition.MaxTime)        // 24 hours max
		assert.Equal(t, int32(120), partition.DefaultTime)     // 2 hours default
		assert.Equal(t, "gpu:1", partition.TresStr)            // GPU resources
		assert.Equal(t, "gpu:1", partition.JobDefaults["gres"]) // Default GPU
		assert.Equal(t, "gpu", partition.QoS)
	})

	t.Run("high memory partition preset", func(t *testing.T) {
		partition, err := NewPartitionBuilder("highmem").
			AsHighMemoryPartition().
			Build()

		require.NoError(t, err)
		assert.Equal(t, int32(75), partition.Priority)         // Medium priority
		assert.Equal(t, int32(2880), partition.MaxTime)        // 48 hours max
		assert.Equal(t, int32(240), partition.DefaultTime)     // 4 hours default
		assert.Equal(t, int64(512*GB), partition.MaxMemPerNode)    // 512GB max memory
		assert.Equal(t, int64(64*GB), partition.DefaultMemPerNode) // 64GB default memory
		assert.Equal(t, "64G", partition.JobDefaults["mem"])   // Default memory
		assert.Equal(t, "highmem", partition.QoS)
	})

	t.Run("maintenance partition preset", func(t *testing.T) {
		partition, err := NewPartitionBuilder("maintenance").
			AsMaintenancePartition().
			Build()

		require.NoError(t, err)
		assert.Equal(t, types.PartitionStateDrain, partition.State) // Drain state
		assert.Equal(t, int32(10), partition.Priority)              // Low priority
		assert.Equal(t, int32(60), partition.MaxTime)               // 1 hour max
		assert.True(t, partition.RootOnly)                          // Root only
		assert.True(t, partition.Hidden)                            // Hidden
	})
}

func TestPartitionBuilder_BusinessRules(t *testing.T) {
	t.Run("default time exceeds max time", func(t *testing.T) {
		_, err := NewPartitionBuilder("test").
			WithDefaultTime(120).
			WithMaxTime(60).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "default time (120) cannot exceed max time (60)")
	})

	t.Run("min nodes exceeds max nodes", func(t *testing.T) {
		_, err := NewPartitionBuilder("test").
			WithMinNodes(10).
			WithMaxNodes(5).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "min nodes (10) cannot exceed max nodes (5)")
	})

	t.Run("default memory exceeds max memory", func(t *testing.T) {
		_, err := NewPartitionBuilder("test").
			WithDefaultMemPerNode(128 * GB).
			WithMaxMemPerNode(64 * GB).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "default memory per node")
	})

	t.Run("conflicting account settings", func(t *testing.T) {
		_, err := NewPartitionBuilder("test").
			WithAllowAccounts("account1").
			WithDenyAccounts("account1").
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "account account1 cannot be both allowed and denied")
	})

	t.Run("conflicting QoS settings", func(t *testing.T) {
		_, err := NewPartitionBuilder("test").
			WithAllowQoS("normal").
			WithDenyQoS("normal").
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "QoS normal cannot be both allowed and denied")
	})

	t.Run("drained partition should be hidden", func(t *testing.T) {
		_, err := NewPartitionBuilder("test").
			WithState(types.PartitionStateDrain).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "drained partitions should typically be hidden")
	})

	t.Run("drained and hidden succeeds", func(t *testing.T) {
		partition, err := NewPartitionBuilder("test").
			WithState(types.PartitionStateDrain).
			AsHidden().
			Build()
		require.NoError(t, err)
		assert.Equal(t, types.PartitionStateDrain, partition.State)
		assert.True(t, partition.Hidden)
	})
}

func TestPartitionBuilder_Clone(t *testing.T) {
	original := NewPartitionBuilder("original").
		WithNodes("node[001-010]").
		WithMaxNodes(10).
		WithAllowAccounts("account1").
		WithJobDefault("mem", "4G")

	// Clone and modify
	cloned := original.Clone()
	cloned.WithNodes("node[011-020]").
		WithMaxNodes(20).
		WithJobDefault("time", "1:00:00")

	// Build both
	originalPartition, err := original.Build()
	require.NoError(t, err)
	clonedPartition, err := cloned.Build()
	require.NoError(t, err)

	// Verify original is unchanged
	assert.Equal(t, "node[001-010]", originalPartition.Nodes)
	assert.Equal(t, int32(10), originalPartition.MaxNodes)
	assert.Equal(t, "4G", originalPartition.JobDefaults["mem"])
	assert.Equal(t, "", originalPartition.JobDefaults["time"])

	// Verify clone has modifications
	assert.Equal(t, "node[011-020]", clonedPartition.Nodes)
	assert.Equal(t, int32(20), clonedPartition.MaxNodes)
	assert.Equal(t, "4G", clonedPartition.JobDefaults["mem"])
	assert.Equal(t, "1:00:00", clonedPartition.JobDefaults["time"])

	// Verify shared attributes are copied
	assert.Equal(t, originalPartition.Name, clonedPartition.Name)
	assert.Equal(t, originalPartition.AllowAccounts, clonedPartition.AllowAccounts)
}

func TestPartitionBuilder_BuildForUpdate(t *testing.T) {
	t.Run("only modified fields", func(t *testing.T) {
		update, err := NewPartitionBuilder("test").
			WithNodes("node[001-100]").
			WithMaxTime(2880).
			WithPriority(200).
			BuildForUpdate()

		require.NoError(t, err)
		require.NotNil(t, update.Nodes)
		assert.Equal(t, "node[001-100]", *update.Nodes)
		require.NotNil(t, update.MaxTime)
		assert.Equal(t, int32(2880), *update.MaxTime)
		require.NotNil(t, update.Priority)
		assert.Equal(t, int32(200), *update.Priority)

		// Default values should not be included
		assert.Nil(t, update.DefaultTime) // Still default 60
	})

	t.Run("with flags and collections", func(t *testing.T) {
		update, err := NewPartitionBuilder("test").
			WithAllowAccounts("research").
			WithPreemptMode("suspend").
			AsHidden().
			BuildForUpdate()

		require.NoError(t, err)
		assert.Equal(t, []string{"research"}, update.AllowAccounts)
		assert.Equal(t, []string{"suspend"}, update.PreemptMode)
		require.NotNil(t, update.Hidden)
		assert.True(t, *update.Hidden)
	})
}

func TestPartitionBuilder_Validation(t *testing.T) {
	tests := []struct {
		name    string
		builder func() *PartitionBuilder
		errMsg  string
	}{
		{
			name: "negative max time",
			builder: func() *PartitionBuilder {
				return NewPartitionBuilder("test").WithMaxTime(-1)
			},
			errMsg: "must be positive",
		},
		{
			name: "negative priority",
			builder: func() *PartitionBuilder {
				return NewPartitionBuilder("test").WithPriority(-1)
			},
			errMsg: "must be non-negative",
		},
		{
			name: "negative grace time",
			builder: func() *PartitionBuilder {
				return NewPartitionBuilder("test").WithGraceTime(-1)
			},
			errMsg: "must be non-negative",
		},
		{
			name: "zero max memory per node",
			builder: func() *PartitionBuilder {
				return NewPartitionBuilder("test").WithMaxMemPerNode(0)
			},
			errMsg: "must be positive",
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

func TestPartitionBuilder_ComplexScenario(t *testing.T) {
	// Build a complex partition with all features
	jobDefaults := map[string]string{
		"mem":           "8G",
		"time":          "2:00:00",
		"cpus-per-task": "2",
	}

	partition, err := NewPartitionBuilder("complex-partition").
		WithNodes("node[001-100]").
		WithAllocNodes("node[001-050]").
		WithAllowAccounts("research", "compute").
		WithAllowGroups("faculty", "students").
		WithAllowQoS("normal", "high", "urgent").
		WithDenyAccounts("guest").
		WithDenyQoS("debug").
		WithDefaultMemPerCPU(4 * GB).
		WithDefaultMemPerNode(128 * GB).
		WithDefaultTime(120).
		WithGraceTime(600).
		WithMaxCPUsPerNode(48).
		WithMaxMemPerNode(512 * GB).
		WithMaxMemPerCPU(16 * GB).
		WithMaxNodes(50).
		WithMaxTime(2880).
		WithMinNodes(1).
		WithOverTimeLimit(120).
		WithPreemptMode("suspend", "cancel").
		WithPriority(100).
		WithPriorityJobFactor(75).
		WithPriorityTier(5).
		WithQoS("normal").
		WithState(types.PartitionStateUp).
		WithTresStr("cpu=4800,mem=25600G,gres/gpu=100").
		WithBillingWeightStr("cpu=1.0,mem=0.25,gres/gpu=2.0").
		WithSelectTypeParameters("CR_Core", "CR_Memory", "CR_GPU").
		WithJobDefaults(jobDefaults).
		WithJobDefault("gres", "gpu:0").
		WithResumeTimeout(300).
		WithSuspendTime(600).
		WithSuspendTimeout(300).
		Build()

	require.NoError(t, err)

	// Verify all fields
	assert.Equal(t, "complex-partition", partition.Name)
	assert.Equal(t, "node[001-100]", partition.Nodes)
	assert.Equal(t, "node[001-050]", partition.AllocNodes)
	assert.Equal(t, []string{"research", "compute"}, partition.AllowAccounts)
	assert.Equal(t, []string{"faculty", "students"}, partition.AllowGroups)
	assert.Equal(t, []string{"normal", "high", "urgent"}, partition.AllowQoS)
	assert.Equal(t, []string{"guest"}, partition.DenyAccounts)
	assert.Equal(t, []string{"debug"}, partition.DenyQoS)
	assert.Equal(t, int64(4*GB), partition.DefaultMemPerCPU)
	assert.Equal(t, int64(128*GB), partition.DefaultMemPerNode)
	assert.Equal(t, int32(120), partition.DefaultTime)
	assert.Equal(t, int32(600), partition.GraceTime)
	assert.Equal(t, int32(48), partition.MaxCPUsPerNode)
	assert.Equal(t, int64(512*GB), partition.MaxMemPerNode)
	assert.Equal(t, int64(16*GB), partition.MaxMemPerCPU)
	assert.Equal(t, int32(50), partition.MaxNodes)
	assert.Equal(t, int32(2880), partition.MaxTime)
	assert.Equal(t, int32(1), partition.MinNodes)
	assert.Equal(t, int32(120), partition.OverTimeLimit)
	assert.Equal(t, []string{"suspend", "cancel"}, partition.PreemptMode)
	assert.Equal(t, int32(100), partition.Priority)
	assert.Equal(t, int32(75), partition.PriorityJobFactor)
	assert.Equal(t, int32(5), partition.PriorityTier)
	assert.Equal(t, "normal", partition.QoS)
	assert.Equal(t, types.PartitionStateUp, partition.State)
	assert.Equal(t, "cpu=4800,mem=25600G,gres/gpu=100", partition.TresStr)
	assert.Equal(t, "cpu=1.0,mem=0.25,gres/gpu=2.0", partition.BillingWeightStr)
	assert.Equal(t, []string{"CR_Core", "CR_Memory", "CR_GPU"}, partition.SelectTypeParameters)
	assert.Equal(t, "8G", partition.JobDefaults["mem"])
	assert.Equal(t, "2:00:00", partition.JobDefaults["time"])
	assert.Equal(t, "2", partition.JobDefaults["cpus-per-task"])
	assert.Equal(t, "gpu:0", partition.JobDefaults["gres"])
	assert.Equal(t, int32(300), partition.ResumeTimeout)
	assert.Equal(t, int32(600), partition.SuspendTime)
	assert.Equal(t, int32(300), partition.SuspendTimeout)
}