// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package builders

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobBuilder_Basic(t *testing.T) {
	t.Run("simple job with command", func(t *testing.T) {
		job, err := NewJobBuilder("echo 'Hello World'").
			WithName("test-job").
			WithAccount("test-account").
			WithPartition("debug").
			Build()

		require.NoError(t, err)
		assert.Equal(t, "test-job", job.Name)
		assert.Equal(t, "echo 'Hello World'", job.Command)
		assert.Equal(t, "test-account", job.Account)
		assert.Equal(t, "debug", job.Partition)
		assert.Equal(t, int32(30), job.TimeLimit) // Default
		assert.Equal(t, int32(1), job.CPUs)       // Default
	})

	t.Run("job with script", func(t *testing.T) {
		script := `#!/bin/bash
echo "Running batch job"
sleep 10`
		job, err := NewJobBuilderFromScript(script).
			WithName("batch-job").
			Build()

		require.NoError(t, err)
		assert.Equal(t, "batch-job", job.Name)
		assert.Equal(t, script, job.Script)
		assert.Equal(t, "", job.Command) // Empty when using script
	})

	t.Run("empty command fails", func(t *testing.T) {
		_, err := NewJobBuilder("").Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "either command or script is required")
	})

	t.Run("both command and script fails", func(t *testing.T) {
		_, err := NewJobBuilder("echo test").
			WithName("test").
			Build()
		// Modify the job to have both command and script
		builder := NewJobBuilder("echo test")
		builder.job.Script = "#!/bin/bash\necho test"
		_, err = builder.Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot specify both command and script")
	})
}

func TestJobBuilder_Resources(t *testing.T) {
	t.Run("basic resources", func(t *testing.T) {
		job, err := NewJobBuilder("./my_program").
			WithName("compute-job").
			WithCPUs(8).
			WithNodes(2).
			WithTasks(16).
			WithTimeLimit(120).
			Build()

		require.NoError(t, err)
		assert.Equal(t, int32(8), job.CPUs)
		assert.Equal(t, int32(2), job.Nodes)
		assert.Equal(t, int32(16), job.Tasks)
		assert.Equal(t, int32(120), job.TimeLimit)
	})

	t.Run("resource requests", func(t *testing.T) {
		job, err := NewJobBuilder("./memory_intensive").
			WithName("memory-job").
			WithResourceRequests().
				WithMemoryGB(32).
				WithMemoryPerCPU(4 * GB).
				WithTmpDisk(100 * GB).
				WithCPUsPerTask(2).
				WithTasksPerNode(4).
				Done().
			Build()

		require.NoError(t, err)
		assert.Equal(t, int64(32*GB), job.ResourceRequests.Memory)
		assert.Equal(t, int64(4*GB), job.ResourceRequests.MemoryPerCPU)
		assert.Equal(t, int64(100*GB), job.ResourceRequests.TmpDisk)
		assert.Equal(t, int32(2), job.ResourceRequests.CPUsPerTask)
		assert.Equal(t, int32(4), job.ResourceRequests.TasksPerNode)
	})

	t.Run("negative resources fail", func(t *testing.T) {
		_, err := NewJobBuilder("test").
			WithCPUs(-1).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be positive")
	})
}

func TestJobBuilder_Environment(t *testing.T) {
	t.Run("environment variables", func(t *testing.T) {
		env := map[string]string{
			"OMP_NUM_THREADS": "8",
			"CUDA_VISIBLE_DEVICES": "0,1",
		}
		
		job, err := NewJobBuilder("./parallel_app").
			WithName("parallel-job").
			WithEnvironment(env).
			WithEnvironmentVar("CUSTOM_VAR", "custom_value").
			Build()

		require.NoError(t, err)
		require.NotNil(t, job.Environment)
		assert.Equal(t, "8", job.Environment["OMP_NUM_THREADS"])
		assert.Equal(t, "0,1", job.Environment["CUDA_VISIBLE_DEVICES"])
		assert.Equal(t, "custom_value", job.Environment["CUSTOM_VAR"])
	})
}

func TestJobBuilder_Dependencies(t *testing.T) {
	t.Run("job dependencies", func(t *testing.T) {
		job, err := NewJobBuilder("./dependent_job").
			WithName("dependent").
			WithDependency("afterok", 12345, 12346).
			WithDependencyState("afternotok", "FAILED", 12347).
			Build()

		require.NoError(t, err)
		require.Len(t, job.Dependencies, 2)
		
		dep1 := job.Dependencies[0]
		assert.Equal(t, "afterok", dep1.Type)
		assert.Equal(t, []int32{12345, 12346}, dep1.JobIDs)
		assert.Equal(t, "", dep1.State)
		
		dep2 := job.Dependencies[1]
		assert.Equal(t, "afternotok", dep2.Type)
		assert.Equal(t, "FAILED", dep2.State)
		assert.Equal(t, []int32{12347}, dep2.JobIDs)
	})
}

func TestJobBuilder_Presets(t *testing.T) {
	t.Run("interactive preset", func(t *testing.T) {
		job, err := NewJobBuilder("python interactive.py").
			WithName("interactive").
			AsInteractive().
			Build()

		require.NoError(t, err)
		assert.Equal(t, int32(60), job.TimeLimit)  // 1 hour
		assert.Equal(t, int32(1), job.CPUs)        // Single CPU
		assert.Equal(t, int32(1), job.Nodes)       // Single node
		assert.Equal(t, "yes", job.Shared)         // Allow sharing
	})

	t.Run("batch preset", func(t *testing.T) {
		job, err := NewJobBuilder("./batch_script.sh").
			WithName("batch").
			AsBatch().
			Build()

		require.NoError(t, err)
		assert.Equal(t, int32(1440), job.TimeLimit) // 24 hours
		assert.Contains(t, job.MailType, "END")
		assert.Contains(t, job.MailType, "FAIL")
		assert.Equal(t, "no", job.Shared) // Exclusive access
	})

	t.Run("array job preset", func(t *testing.T) {
		job, err := NewJobBuilder("./array_task.sh").
			WithName("array").
			AsArrayJob("1-100").
			Build()

		require.NoError(t, err)
		assert.Equal(t, "1-100", job.ArrayString)
		assert.Equal(t, int32(240), job.TimeLimit) // 4 hours
		assert.Equal(t, int32(1), job.CPUs)        // Single CPU per task
		assert.Contains(t, job.MailType, "END")
		assert.Contains(t, job.MailType, "FAIL")
	})

	t.Run("GPU job preset", func(t *testing.T) {
		job, err := NewJobBuilder("./gpu_program").
			WithName("gpu").
			AsGPUJob(2).
			Build()

		require.NoError(t, err)
		assert.Equal(t, "gpu:2", job.Gres)
		assert.Equal(t, int32(480), job.TimeLimit) // 8 hours
		assert.Equal(t, int32(8), job.CPUs)        // 4 CPUs per GPU
		assert.Contains(t, job.Features, "gpu")
	})

	t.Run("high memory job preset", func(t *testing.T) {
		job, err := NewJobBuilder("./memory_intensive").
			WithName("highmem").
			AsHighMemoryJob().
			Build()

		require.NoError(t, err)
		assert.Contains(t, job.Features, "highmem")
		assert.Equal(t, int64(64*GB), job.ResourceRequests.Memory)
		assert.Equal(t, int32(720), job.TimeLimit) // 12 hours
	})
}

func TestJobBuilder_BusinessRules(t *testing.T) {
	t.Run("task node consistency", func(t *testing.T) {
		_, err := NewJobBuilder("test").
			WithNodes(2).
			WithTasks(40). // More than 2 * 16
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "tasks (40) cannot exceed nodes * 16 (32)")
	})

	t.Run("array job time limit", func(t *testing.T) {
		_, err := NewJobBuilder("test").
			WithArrayString("1-100").
			WithTimeLimit(2000). // More than 24 hours
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "array jobs should not exceed 24 hours")
	})

	t.Run("high memory requires feature", func(t *testing.T) {
		_, err := NewJobBuilder("test").
			WithResourceRequests().
				WithMemoryGB(40). // More than 32GB
				Done().
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "jobs requiring > 32GB memory must use highmem feature")
	})

	t.Run("GPU auto-adds feature", func(t *testing.T) {
		job, err := NewJobBuilder("test").
			WithGres("gpu:1").
			Build()
		require.NoError(t, err)
		assert.Contains(t, job.Features, "gpu")
	})
}

func TestJobBuilder_Clone(t *testing.T) {
	original := NewJobBuilder("original_command").
		WithName("original").
		WithAccount("account1").
		WithCPUs(4).
		WithEnvironmentVar("VAR1", "value1").
		WithDependency("afterok", 123)

	// Clone and modify
	cloned := original.Clone()
	cloned.WithName("cloned").
		WithAccount("account2").
		WithEnvironmentVar("VAR2", "value2")

	// Build both
	originalJob, err := original.Build()
	require.NoError(t, err)
	clonedJob, err := cloned.Build()
	require.NoError(t, err)

	// Verify original is unchanged
	assert.Equal(t, "original", originalJob.Name)
	assert.Equal(t, "account1", originalJob.Account)
	assert.Equal(t, "value1", originalJob.Environment["VAR1"])
	assert.Equal(t, "", originalJob.Environment["VAR2"])

	// Verify clone has modifications
	assert.Equal(t, "cloned", clonedJob.Name)
	assert.Equal(t, "account2", clonedJob.Account)
	assert.Equal(t, "value1", clonedJob.Environment["VAR1"])
	assert.Equal(t, "value2", clonedJob.Environment["VAR2"])

	// Verify shared attributes are copied
	assert.Equal(t, originalJob.Command, clonedJob.Command)
	assert.Equal(t, originalJob.CPUs, clonedJob.CPUs)
	assert.Equal(t, originalJob.Dependencies, clonedJob.Dependencies)
}

func TestJobBuilder_BuildForUpdate(t *testing.T) {
	t.Run("only modified fields", func(t *testing.T) {
		update, err := NewJobBuilder("test_command").
			WithName("updated-job").
			WithAccount("new-account").
			WithTimeLimit(180).
			BuildForUpdate()

		require.NoError(t, err)
		require.NotNil(t, update.Name)
		assert.Equal(t, "updated-job", *update.Name)
		require.NotNil(t, update.Account)
		assert.Equal(t, "new-account", *update.Account)
		require.NotNil(t, update.TimeLimit)
		assert.Equal(t, int32(180), *update.TimeLimit)

		// Default values should not be included for nodes, cpus, etc.
		assert.Nil(t, update.Priority) // Wasn't set
	})

	t.Run("with priority and features", func(t *testing.T) {
		priority := int32(100)
		update, err := NewJobBuilder("test").
			WithPriority(priority).
			WithFeatures("gpu", "infiniband").
			BuildForUpdate()

		require.NoError(t, err)
		require.NotNil(t, update.Priority)
		assert.Equal(t, priority, *update.Priority)
		assert.Equal(t, []string{"gpu", "infiniband"}, update.Features)
	})
}

func TestJobBuilder_Validation(t *testing.T) {
	tests := []struct {
		name    string
		builder func() *JobBuilder
		errMsg  string
	}{
		{
			name: "negative time limit",
			builder: func() *JobBuilder {
				return NewJobBuilder("test").WithTimeLimit(-1)
			},
			errMsg: "must be positive",
		},
		{
			name: "negative priority",
			builder: func() *JobBuilder {
				return NewJobBuilder("test").WithPriority(-1)
			},
			errMsg: "must be non-negative",
		},
		{
			name: "invalid nice value",
			builder: func() *JobBuilder {
				return NewJobBuilder("test").WithNice(25)
			},
			errMsg: "must be between -20 and 19",
		},
		{
			name: "empty job name",
			builder: func() *JobBuilder {
				return NewJobBuilder("test").WithName("")
			},
			errMsg: "cannot be empty",
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

func TestJobBuilder_TimeConvenience(t *testing.T) {
	t.Run("time limit from duration", func(t *testing.T) {
		job, err := NewJobBuilder("test").
			WithTimeLimitDuration(2 * time.Hour).
			Build()

		require.NoError(t, err)
		assert.Equal(t, int32(120), job.TimeLimit) // 2 hours in minutes
	})

	t.Run("deadline setting", func(t *testing.T) {
		deadline := time.Now().Add(24 * time.Hour)
		job, err := NewJobBuilder("test").
			WithDeadline(deadline).
			Build()

		require.NoError(t, err)
		require.NotNil(t, job.Deadline)
		assert.Equal(t, deadline, *job.Deadline)
	})
}

func TestJobResourceRequestsBuilder(t *testing.T) {
	t.Run("comprehensive resource requests", func(t *testing.T) {
		job, err := NewJobBuilder("./app").
			WithResourceRequests().
				WithMemoryGB(16).
				WithMemoryPerCPU(2 * GB).
				WithMemoryPerGPU(8 * GB).
				WithTmpDisk(50 * GB).
				WithCPUsPerTask(4).
				WithTasksPerNode(2).
				WithTasksPerCore(1).
				WithThreadsPerCore(2).
				Done().
			Build()

		require.NoError(t, err)
		res := job.ResourceRequests
		assert.Equal(t, int64(16*GB), res.Memory)
		assert.Equal(t, int64(2*GB), res.MemoryPerCPU)
		assert.Equal(t, int64(8*GB), res.MemoryPerGPU)
		assert.Equal(t, int64(50*GB), res.TmpDisk)
		assert.Equal(t, int32(4), res.CPUsPerTask)
		assert.Equal(t, int32(2), res.TasksPerNode)
		assert.Equal(t, int32(1), res.TasksPerCore)
		assert.Equal(t, int32(2), res.ThreadsPerCore)
	})

	t.Run("negative resource requests fail", func(t *testing.T) {
		_, err := NewJobBuilder("test").
			WithResourceRequests().
				WithMemoryGB(-1).
				Done().
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be non-negative")
	})
}

func TestJobBuilder_ComplexScenario(t *testing.T) {
	// Build a complex job with all features
	job, err := NewJobBuilder("./complex_application").
		WithName("complex-job").
		WithAccount("research-account").
		WithPartition("gpu-partition").
		WithQoS("high-priority").
		WithTimeLimit(480). // 8 hours
		WithCPUs(16).
		WithNodes(2).
		WithTasks(32).
		WithWorkingDirectory("/scratch/user/job").
		WithStandardOutput("output_%j.log").
		WithStandardError("error_%j.log").
		WithMailType("BEGIN", "END", "FAIL").
		WithMailUser("user@domain.com").
		WithFeatures("gpu", "infiniband", "highmem").
		WithGres("gpu:4").
		WithEnvironmentVar("OMP_NUM_THREADS", "4").
		WithEnvironmentVar("CUDA_VISIBLE_DEVICES", "0,1,2,3").
		WithResourceRequests().
			WithMemoryGB(128).
			WithMemoryPerCPU(8 * GB).
			WithTmpDisk(500 * GB).
			WithCPUsPerTask(2).
			Done().
		WithDependency("afterok", 12345).
		WithComment("Complex multi-GPU job").
		Build()

	require.NoError(t, err)

	// Verify all fields
	assert.Equal(t, "complex-job", job.Name)
	assert.Equal(t, "./complex_application", job.Command)
	assert.Equal(t, "research-account", job.Account)
	assert.Equal(t, "gpu-partition", job.Partition)
	assert.Equal(t, "high-priority", job.QoS)
	assert.Equal(t, int32(480), job.TimeLimit)
	assert.Equal(t, int32(16), job.CPUs)
	assert.Equal(t, int32(2), job.Nodes)
	assert.Equal(t, int32(32), job.Tasks)
	assert.Equal(t, "/scratch/user/job", job.WorkingDirectory)
	assert.Equal(t, "output_%j.log", job.StandardOutput)
	assert.Equal(t, "error_%j.log", job.StandardError)
	assert.Len(t, job.MailType, 3)
	assert.Equal(t, "user@domain.com", job.MailUser)
	assert.Len(t, job.Features, 3)
	assert.Equal(t, "gpu:4", job.Gres)
	assert.Equal(t, "4", job.Environment["OMP_NUM_THREADS"])
	assert.Equal(t, "0,1,2,3", job.Environment["CUDA_VISIBLE_DEVICES"])
	assert.Equal(t, int64(128*GB), job.ResourceRequests.Memory)
	assert.Equal(t, int64(8*GB), job.ResourceRequests.MemoryPerCPU)
	assert.Equal(t, int64(500*GB), job.ResourceRequests.TmpDisk)
	assert.Equal(t, int32(2), job.ResourceRequests.CPUsPerTask)
	assert.Len(t, job.Dependencies, 1)
	assert.Equal(t, "Complex multi-GPU job", job.Comment)
}
