// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package builders

import (
	"fmt"
	"strings"
	"time"

	types "github.com/jontk/slurm-client/api"
)

// JobBuilder provides a fluent interface for building Job objects
type JobBuilder struct {
	job    *types.JobCreate
	errors []error
}

// Helper function to create a pointer to a value
func ptrString(s string) *string { return &s }
func ptrInt32(i int32) *int32    { return &i }
func ptrUint32(i uint32) *uint32 { return &i }
func ptrBool(b bool) *bool       { return &b }
func ptrInt64(i int64) *int64    { return &i }
func ptrUint64(i uint64) *uint64 { return &i }

// NewJobBuilder creates a new Job builder with the required script
// Note: The new API uses Script instead of Command
func NewJobBuilder(script string) *JobBuilder {
	return &JobBuilder{
		job: &types.JobCreate{
			Script:      ptrString(script),
			TimeLimit:   ptrUint32(30), // Default time limit in minutes
			MinimumCPUs: ptrInt32(1),   // Default CPU count
			MinimumNodes: ptrInt32(1),  // Default node count
			Tasks:       ptrInt32(1),   // Default task count
			Environment: []string{},    // Initialize empty
			MailType:    []types.MailTypeValue{}, // Initialize empty
		},
		errors: []error{},
	}
}

// NewJobBuilderFromScript creates a new Job builder with a batch script
func NewJobBuilderFromScript(script string) *JobBuilder {
	return NewJobBuilder(script)
}

// WithName sets the job name
func (b *JobBuilder) WithName(name string) *JobBuilder {
	if name == "" {
		b.addError(fmt.Errorf("job name cannot be empty"))
		return b
	}
	b.job.Name = ptrString(name)
	return b
}

// WithAccount sets the account to charge
func (b *JobBuilder) WithAccount(account string) *JobBuilder {
	b.job.Account = ptrString(account)
	return b
}

// WithPartition sets the partition to submit to
func (b *JobBuilder) WithPartition(partition string) *JobBuilder {
	b.job.Partition = ptrString(partition)
	return b
}

// WithQoS sets the Quality of Service
func (b *JobBuilder) WithQoS(qos string) *JobBuilder {
	b.job.QoS = ptrString(qos)
	return b
}

// WithTimeLimit sets the wall clock time limit in minutes
func (b *JobBuilder) WithTimeLimit(minutes int32) *JobBuilder {
	if minutes <= 0 {
		b.addError(fmt.Errorf("time limit must be positive, got %d", minutes))
		return b
	}
	b.job.TimeLimit = ptrUint32(uint32(minutes))
	return b
}

// WithTimeLimitDuration sets the wall clock time limit using duration
func (b *JobBuilder) WithTimeLimitDuration(d time.Duration) *JobBuilder {
	minutes := int32(d.Minutes())
	if minutes <= 0 {
		b.addError(fmt.Errorf("duration must be positive, got %v", d))
		return b
	}
	b.job.TimeLimit = ptrUint32(uint32(minutes))
	return b
}

// WithPriority sets the job priority
func (b *JobBuilder) WithPriority(priority int32) *JobBuilder {
	b.job.Priority = ptrUint32(uint32(priority))
	return b
}

// WithCPUs sets the number of CPUs required
func (b *JobBuilder) WithCPUs(cpus int32) *JobBuilder {
	if cpus <= 0 {
		b.addError(fmt.Errorf("CPUs must be positive, got %d", cpus))
		return b
	}
	b.job.MinimumCPUs = ptrInt32(cpus)
	return b
}

// WithNodes sets the number of nodes required
func (b *JobBuilder) WithNodes(nodes int32) *JobBuilder {
	if nodes <= 0 {
		b.addError(fmt.Errorf("nodes must be positive, got %d", nodes))
		return b
	}
	b.job.MinimumNodes = ptrInt32(nodes)
	b.job.MaximumNodes = ptrInt32(nodes)
	return b
}

// WithTasks sets the number of tasks
func (b *JobBuilder) WithTasks(tasks int32) *JobBuilder {
	if tasks <= 0 {
		b.addError(fmt.Errorf("tasks must be positive, got %d", tasks))
		return b
	}
	b.job.Tasks = ptrInt32(tasks)
	return b
}

// WithScript sets the batch script
func (b *JobBuilder) WithScript(script string) *JobBuilder {
	b.job.Script = ptrString(script)
	return b
}

// WithWorkingDirectory sets the current working directory
func (b *JobBuilder) WithWorkingDirectory(dir string) *JobBuilder {
	b.job.CurrentWorkingDirectory = ptrString(dir)
	return b
}

// WithStandardOutput sets the standard output file
func (b *JobBuilder) WithStandardOutput(path string) *JobBuilder {
	b.job.StandardOutput = ptrString(path)
	return b
}

// WithStandardError sets the standard error file
func (b *JobBuilder) WithStandardError(path string) *JobBuilder {
	b.job.StandardError = ptrString(path)
	return b
}

// WithStandardInput sets the standard input file
func (b *JobBuilder) WithStandardInput(path string) *JobBuilder {
	b.job.StandardInput = ptrString(path)
	return b
}

// WithEnvironmentVariable adds an environment variable
// Note: Environment is now []string in format "KEY=VALUE"
func (b *JobBuilder) WithEnvironmentVariable(key, value string) *JobBuilder {
	b.job.Environment = append(b.job.Environment, key+"="+value)
	return b
}

// WithEnvironment sets multiple environment variables from a map
func (b *JobBuilder) WithEnvironment(env map[string]string) *JobBuilder {
	for k, v := range env {
		b.job.Environment = append(b.job.Environment, k+"="+v)
	}
	return b
}

// WithMailUser sets the mail notification recipient
func (b *JobBuilder) WithMailUser(email string) *JobBuilder {
	b.job.MailUser = ptrString(email)
	return b
}

// WithMailType sets the mail notification types
func (b *JobBuilder) WithMailType(mailTypes ...types.MailTypeValue) *JobBuilder {
	b.job.MailType = append(b.job.MailType, mailTypes...)
	return b
}

// WithReservation sets the reservation to use
func (b *JobBuilder) WithReservation(reservation string) *JobBuilder {
	b.job.Reservation = ptrString(reservation)
	return b
}

// WithComment sets a comment for the job
func (b *JobBuilder) WithComment(comment string) *JobBuilder {
	b.job.Comment = ptrString(comment)
	return b
}

// WithDeadline sets the job deadline
func (b *JobBuilder) WithDeadline(deadline time.Time) *JobBuilder {
	ts := deadline.Unix()
	b.job.Deadline = ptrInt64(ts)
	return b
}

// WithBeginTime sets the earliest start time
func (b *JobBuilder) WithBeginTime(beginTime time.Time) *JobBuilder {
	ts := uint64(beginTime.Unix())
	b.job.BeginTime = ptrUint64(ts)
	return b
}

// WithImmediate sets whether the job should fail if resources aren't immediately available
func (b *JobBuilder) WithImmediate(immediate bool) *JobBuilder {
	b.job.Immediate = ptrBool(immediate)
	return b
}

// WithContiguous sets whether contiguous nodes are required
func (b *JobBuilder) WithContiguous(contiguous bool) *JobBuilder {
	b.job.Contiguous = ptrBool(contiguous)
	return b
}

// WithExcludedNodes sets nodes to exclude
func (b *JobBuilder) WithExcludedNodes(nodes ...string) *JobBuilder {
	b.job.ExcludedNodes = nodes
	return b
}

// WithConstraints sets the feature constraints (e.g., "gpu", "nvme")
// Note: Features are now specified as a comma-separated string in Constraints field
func (b *JobBuilder) WithConstraints(constraints string) *JobBuilder {
	b.job.Constraints = ptrString(constraints)
	return b
}

// WithFeatures sets the feature constraints from a list
// Joins them with "&" for SLURM constraint syntax
func (b *JobBuilder) WithFeatures(features ...string) *JobBuilder {
	if len(features) > 0 {
		b.job.Constraints = ptrString(strings.Join(features, "&"))
	}
	return b
}

// WithDependency sets job dependencies in SLURM format
// Example: "afterok:123:456" or "afterany:789"
func (b *JobBuilder) WithDependency(dependency string) *JobBuilder {
	b.job.Dependency = ptrString(dependency)
	return b
}

// WithArray sets the job array specification
// Example: "0-15" or "0-15%4" (max 4 concurrent)
func (b *JobBuilder) WithArray(arraySpec string) *JobBuilder {
	b.job.Array = ptrString(arraySpec)
	return b
}

// WithMemoryPerNode sets memory per node in MB
func (b *JobBuilder) WithMemoryPerNode(mb uint64) *JobBuilder {
	b.job.MemoryPerNode = ptrUint64(mb)
	return b
}

// WithMemoryPerCPU sets memory per CPU in MB
func (b *JobBuilder) WithMemoryPerCPU(mb uint64) *JobBuilder {
	b.job.MemoryPerCPU = ptrUint64(mb)
	return b
}

// WithCPUBinding sets the CPU binding method
func (b *JobBuilder) WithCPUBinding(binding string) *JobBuilder {
	b.job.CPUBinding = ptrString(binding)
	return b
}

// WithMemoryBinding sets the memory binding method
func (b *JobBuilder) WithMemoryBinding(binding string) *JobBuilder {
	b.job.MemoryBinding = ptrString(binding)
	return b
}

// WithDistribution sets the task distribution method
func (b *JobBuilder) WithDistribution(dist string) *JobBuilder {
	b.job.Distribution = ptrString(dist)
	return b
}

// WithNice sets the job nice value
func (b *JobBuilder) WithNice(nice int32) *JobBuilder {
	b.job.Nice = ptrInt32(nice)
	return b
}

// Build creates the JobCreate object, returning any accumulated errors
func (b *JobBuilder) Build() (*types.JobCreate, error) {
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("validation errors: %v", b.errors)
	}
	return b.job, nil
}

// MustBuild creates the JobCreate object, panicking on error
func (b *JobBuilder) MustBuild() *types.JobCreate {
	job, err := b.Build()
	if err != nil {
		panic(err)
	}
	return job
}

// addError adds an error to the builder's error list
func (b *JobBuilder) addError(err error) {
	b.errors = append(b.errors, err)
}

// Errors returns any accumulated errors
func (b *JobBuilder) Errors() []error {
	return b.errors
}

// HasErrors returns true if there are validation errors
func (b *JobBuilder) HasErrors() bool {
	return len(b.errors) > 0
}

// ApplyUpdate applies a JobUpdate (now alias for JobCreate) to modify settings
func (b *JobBuilder) ApplyUpdate(update *types.JobUpdate) *JobBuilder {
	if update == nil {
		return b
	}
	if update.Name != nil {
		b.job.Name = update.Name
	}
	if update.Account != nil {
		b.job.Account = update.Account
	}
	if update.Partition != nil {
		b.job.Partition = update.Partition
	}
	if update.QoS != nil {
		b.job.QoS = update.QoS
	}
	if update.TimeLimit != nil {
		b.job.TimeLimit = update.TimeLimit
	}
	if update.Priority != nil {
		b.job.Priority = update.Priority
	}
	if update.Nice != nil {
		b.job.Nice = update.Nice
	}
	if update.Comment != nil {
		b.job.Comment = update.Comment
	}
	return b
}
