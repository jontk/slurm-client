// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package builders

import (
	"fmt"
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
)

// JobBuilder provides a fluent interface for building Job objects
type JobBuilder struct {
	job    *types.JobCreate
	errors []error
}

// NewJobBuilder creates a new Job builder with the required command or script
func NewJobBuilder(command string) *JobBuilder {
	return &JobBuilder{
		job: &types.JobCreate{
			Command:          command,
			TimeLimit:        30,                       // Default time limit in minutes
			CPUs:             1,                        // Default CPU count
			Nodes:            1,                        // Default node count
			Tasks:            1,                        // Default task count
			Environment:      make(map[string]string), // Initialize empty
			Dependencies:     []types.JobDependency{}, // Initialize empty
			ResourceRequests: types.ResourceRequests{}, // Initialize empty
			MailType:         []string{},               // Initialize empty
			ClusterFeatures:  []string{},               // Initialize empty
			Features:         []string{},               // Initialize empty
			Profile:          []string{},               // Initialize empty
		},
		errors: []error{},
	}
}

// NewJobBuilderFromScript creates a new Job builder with a batch script
func NewJobBuilderFromScript(script string) *JobBuilder {
	return &JobBuilder{
		job: &types.JobCreate{
			Script:           script,
			Command:          "", // Empty when using script
			TimeLimit:        30,
			CPUs:             1,
			Nodes:            1,
			Tasks:            1,
			Environment:      make(map[string]string),
			Dependencies:     []types.JobDependency{},
			ResourceRequests: types.ResourceRequests{},
			MailType:         []string{},
			ClusterFeatures:  []string{},
			Features:         []string{},
			Profile:          []string{},
		},
		errors: []error{},
	}
}

// WithName sets the job name
func (b *JobBuilder) WithName(name string) *JobBuilder {
	if name == "" {
		b.addError(fmt.Errorf("job name cannot be empty"))
		return b
	}
	b.job.Name = name
	return b
}

// WithAccount sets the account to charge
func (b *JobBuilder) WithAccount(account string) *JobBuilder {
	b.job.Account = account
	return b
}

// WithPartition sets the partition to submit to
func (b *JobBuilder) WithPartition(partition string) *JobBuilder {
	b.job.Partition = partition
	return b
}

// WithQoS sets the Quality of Service
func (b *JobBuilder) WithQoS(qos string) *JobBuilder {
	b.job.QoS = qos
	return b
}

// WithTimeLimit sets the wall clock time limit in minutes
func (b *JobBuilder) WithTimeLimit(minutes int32) *JobBuilder {
	if minutes <= 0 {
		b.addError(fmt.Errorf("time limit must be positive, got %d", minutes))
		return b
	}
	b.job.TimeLimit = minutes
	return b
}

// WithTimeLimitDuration sets the wall clock time limit using duration
func (b *JobBuilder) WithTimeLimitDuration(duration time.Duration) *JobBuilder {
	minutes := int32(duration.Minutes())
	return b.WithTimeLimit(minutes)
}

// WithPriority sets the job priority
func (b *JobBuilder) WithPriority(priority int32) *JobBuilder {
	if priority < 0 {
		b.addError(fmt.Errorf("priority must be non-negative, got %d", priority))
		return b
	}
	b.job.Priority = &priority
	return b
}

// WithCPUs sets the number of CPUs required
func (b *JobBuilder) WithCPUs(cpus int32) *JobBuilder {
	if cpus <= 0 {
		b.addError(fmt.Errorf("CPUs must be positive, got %d", cpus))
		return b
	}
	b.job.CPUs = cpus
	return b
}

// WithNodes sets the number of nodes required
func (b *JobBuilder) WithNodes(nodes int32) *JobBuilder {
	if nodes <= 0 {
		b.addError(fmt.Errorf("nodes must be positive, got %d", nodes))
		return b
	}
	b.job.Nodes = nodes
	return b
}

// WithTasks sets the number of tasks
func (b *JobBuilder) WithTasks(tasks int32) *JobBuilder {
	if tasks <= 0 {
		b.addError(fmt.Errorf("tasks must be positive, got %d", tasks))
		return b
	}
	b.job.Tasks = tasks
	return b
}

// WithWorkingDirectory sets the working directory
func (b *JobBuilder) WithWorkingDirectory(dir string) *JobBuilder {
	b.job.WorkingDirectory = dir
	return b
}

// WithStandardOutput sets the standard output file
func (b *JobBuilder) WithStandardOutput(file string) *JobBuilder {
	b.job.StandardOutput = file
	return b
}

// WithStandardError sets the standard error file
func (b *JobBuilder) WithStandardError(file string) *JobBuilder {
	b.job.StandardError = file
	return b
}

// WithStandardInput sets the standard input file
func (b *JobBuilder) WithStandardInput(file string) *JobBuilder {
	b.job.StandardInput = file
	return b
}

// WithArrayString sets the array specification
func (b *JobBuilder) WithArrayString(arrayString string) *JobBuilder {
	b.job.ArrayString = arrayString
	return b
}

// WithEnvironment sets environment variables
func (b *JobBuilder) WithEnvironment(env map[string]string) *JobBuilder {
	if b.job.Environment == nil {
		b.job.Environment = make(map[string]string)
	}
	for k, v := range env {
		b.job.Environment[k] = v
	}
	return b
}

// WithEnvironmentVar sets a single environment variable
func (b *JobBuilder) WithEnvironmentVar(key, value string) *JobBuilder {
	if b.job.Environment == nil {
		b.job.Environment = make(map[string]string)
	}
	b.job.Environment[key] = value
	return b
}

// WithMailType sets mail notification types
func (b *JobBuilder) WithMailType(mailTypes ...string) *JobBuilder {
	b.job.MailType = append(b.job.MailType, mailTypes...)
	return b
}

// WithMailUser sets the email address for notifications
func (b *JobBuilder) WithMailUser(email string) *JobBuilder {
	b.job.MailUser = email
	return b
}

// WithExcludeNodes sets nodes to exclude
func (b *JobBuilder) WithExcludeNodes(nodes string) *JobBuilder {
	b.job.ExcludeNodes = nodes
	return b
}

// WithNice sets the nice value
func (b *JobBuilder) WithNice(nice int32) *JobBuilder {
	if nice < -20 || nice > 19 {
		b.addError(fmt.Errorf("nice value must be between -20 and 19, got %d", nice))
		return b
	}
	b.job.Nice = nice
	return b
}

// WithComment sets a comment for the job
func (b *JobBuilder) WithComment(comment string) *JobBuilder {
	b.job.Comment = comment
	return b
}

// WithDeadline sets the job deadline
func (b *JobBuilder) WithDeadline(deadline time.Time) *JobBuilder {
	b.job.Deadline = &deadline
	return b
}

// WithClusterFeatures sets required cluster features
func (b *JobBuilder) WithClusterFeatures(features ...string) *JobBuilder {
	b.job.ClusterFeatures = append(b.job.ClusterFeatures, features...)
	return b
}

// WithFeatures sets required node features
func (b *JobBuilder) WithFeatures(features ...string) *JobBuilder {
	b.job.Features = append(b.job.Features, features...)
	return b
}

// WithGres sets generic resource requirements
func (b *JobBuilder) WithGres(gres string) *JobBuilder {
	b.job.Gres = gres
	return b
}

// WithShared sets shared resource policy
func (b *JobBuilder) WithShared(shared string) *JobBuilder {
	b.job.Shared = shared
	return b
}

// WithProfile sets profiling options
func (b *JobBuilder) WithProfile(profiles ...string) *JobBuilder {
	b.job.Profile = append(b.job.Profile, profiles...)
	return b
}

// WithReservation sets the reservation to use
func (b *JobBuilder) WithReservation(reservation string) *JobBuilder {
	b.job.Reservation = reservation
	return b
}

// WithResourceRequests returns a resource requests builder
func (b *JobBuilder) WithResourceRequests() *JobResourceRequestsBuilder {
	return &JobResourceRequestsBuilder{
		parent:    b,
		resources: &b.job.ResourceRequests,
	}
}

// WithDependency adds a job dependency
func (b *JobBuilder) WithDependency(depType string, jobIDs ...int32) *JobBuilder {
	dependency := types.JobDependency{
		Type:   depType,
		JobIDs: jobIDs,
	}
	b.job.Dependencies = append(b.job.Dependencies, dependency)
	return b
}

// WithDependencyState adds a job dependency with state
func (b *JobBuilder) WithDependencyState(depType, state string, jobIDs ...int32) *JobBuilder {
	dependency := types.JobDependency{
		Type:   depType,
		State:  state,
		JobIDs: jobIDs,
	}
	b.job.Dependencies = append(b.job.Dependencies, dependency)
	return b
}

// AsInteractive applies common interactive job settings
func (b *JobBuilder) AsInteractive() *JobBuilder {
	return b.
		WithTimeLimit(60).  // 1 hour default
		WithCPUs(1).        // Single CPU
		WithNodes(1).       // Single node
		WithShared("yes")   // Allow sharing
}

// AsBatch applies common batch job settings
func (b *JobBuilder) AsBatch() *JobBuilder {
	return b.
		WithTimeLimit(1440). // 24 hours default
		WithMailType("END", "FAIL").
		WithShared("no") // Exclusive access
}

// AsArrayJob applies common array job settings
func (b *JobBuilder) AsArrayJob(arraySpec string) *JobBuilder {
	return b.
		WithArrayString(arraySpec).
		WithTimeLimit(240).  // 4 hours default
		WithCPUs(1).         // Single CPU per task
		WithMailType("END", "FAIL")
}

// AsGPUJob applies common GPU job settings
func (b *JobBuilder) AsGPUJob(gpuCount int) *JobBuilder {
	return b.
		WithGres(fmt.Sprintf("gpu:%d", gpuCount)).
		WithTimeLimit(480).  // 8 hours default
		WithCPUs(int32(gpuCount * 4)). // 4 CPUs per GPU
		WithFeatures("gpu")
}

// AsHighMemoryJob applies settings for high memory jobs
func (b *JobBuilder) AsHighMemoryJob() *JobBuilder {
	return b.
		WithFeatures("highmem").
		WithResourceRequests().
			WithMemoryPerNode(64 * GB).  // 64GB default
			Done().
		WithTimeLimit(720) // 12 hours default
}

// Build validates and returns the built Job
func (b *JobBuilder) Build() (*types.JobCreate, error) {
	// Check for accumulated errors
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("builder errors: %v", b.errors)
	}

	// Validate required fields
	if b.job.Command == "" && b.job.Script == "" {
		return nil, fmt.Errorf("either command or script is required")
	}

	if b.job.Command != "" && b.job.Script != "" {
		return nil, fmt.Errorf("cannot specify both command and script")
	}

	// Apply business rules
	if err := b.validateBusinessRules(); err != nil {
		return nil, err
	}

	return b.job, nil
}

// BuildForUpdate creates a JobUpdate object from the builder
func (b *JobBuilder) BuildForUpdate() (*types.JobUpdate, error) {
	// Check for accumulated errors
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("builder errors: %v", b.errors)
	}

	update := &types.JobUpdate{}

	// Only include fields that were explicitly set
	if b.job.Name != "" {
		update.Name = &b.job.Name
	}
	if b.job.Account != "" {
		update.Account = &b.job.Account
	}
	if b.job.Partition != "" {
		update.Partition = &b.job.Partition
	}
	if b.job.QoS != "" {
		update.QoS = &b.job.QoS
	}
	if b.job.TimeLimit != 30 {  // Not default
		update.TimeLimit = &b.job.TimeLimit
	}
	if b.job.Priority != nil {
		update.Priority = b.job.Priority
	}
	if b.job.Nice != 0 {
		update.Nice = &b.job.Nice
	}
	if b.job.Comment != "" {
		update.Comment = &b.job.Comment
	}
	if b.job.Deadline != nil {
		update.Deadline = b.job.Deadline
	}
	if len(b.job.Features) > 0 {
		update.Features = b.job.Features
	}
	if b.job.Gres != "" {
		update.Gres = &b.job.Gres
	}
	if b.job.Shared != "" {
		update.Shared = &b.job.Shared
	}
	if b.job.Reservation != "" {
		update.Reservation = &b.job.Reservation
	}
	if b.job.ExcludeNodes != "" {
		update.ExcludeNodes = &b.job.ExcludeNodes
	}

	return update, nil
}

// Clone creates a copy of the builder with the same settings
func (b *JobBuilder) Clone() *JobBuilder {
	newBuilder := &JobBuilder{
		job: &types.JobCreate{
			Name:              b.job.Name,
			Account:           b.job.Account,
			Partition:         b.job.Partition,
			QoS:              b.job.QoS,
			TimeLimit:        b.job.TimeLimit,
			CPUs:             b.job.CPUs,
			Nodes:            b.job.Nodes,
			Tasks:            b.job.Tasks,
			Command:          b.job.Command,
			Script:           b.job.Script,
			WorkingDirectory: b.job.WorkingDirectory,
			StandardInput:    b.job.StandardInput,
			StandardOutput:   b.job.StandardOutput,
			StandardError:    b.job.StandardError,
			ArrayString:      b.job.ArrayString,
			MailUser:         b.job.MailUser,
			ExcludeNodes:     b.job.ExcludeNodes,
			Nice:             b.job.Nice,
			Comment:          b.job.Comment,
			Gres:             b.job.Gres,
			Shared:           b.job.Shared,
			Reservation:      b.job.Reservation,
			ResourceRequests: b.job.ResourceRequests,
		},
		errors: append([]error{}, b.errors...),
	}

	// Copy pointer fields
	if b.job.Priority != nil {
		v := *b.job.Priority
		newBuilder.job.Priority = &v
	}
	if b.job.Deadline != nil {
		v := *b.job.Deadline
		newBuilder.job.Deadline = &v
	}

	// Deep copy slices and maps
	newBuilder.job.Dependencies = append([]types.JobDependency{}, b.job.Dependencies...)
	newBuilder.job.MailType = append([]string{}, b.job.MailType...)
	newBuilder.job.ClusterFeatures = append([]string{}, b.job.ClusterFeatures...)
	newBuilder.job.Features = append([]string{}, b.job.Features...)
	newBuilder.job.Profile = append([]string{}, b.job.Profile...)
	
	if b.job.Environment != nil {
		newBuilder.job.Environment = make(map[string]string)
		for k, v := range b.job.Environment {
			newBuilder.job.Environment[k] = v
		}
	}

	return newBuilder
}

// addError adds an error to the builder's error list
func (b *JobBuilder) addError(err error) {
	b.errors = append(b.errors, err)
}

// validateBusinessRules applies business logic validation
func (b *JobBuilder) validateBusinessRules() error {
	// Validate resource consistency
	if b.job.Tasks > b.job.Nodes*16 { // Assume max 16 cores per node
		return fmt.Errorf("tasks (%d) cannot exceed nodes * 16 (%d)", b.job.Tasks, b.job.Nodes*16)
	}

	// Validate array job limits
	if b.job.ArrayString != "" && b.job.TimeLimit > 1440 {
		return fmt.Errorf("array jobs should not exceed 24 hours time limit")
	}

	// Validate GPU requirements
	if b.job.Gres != "" && len(b.job.Features) == 0 {
		b.job.Features = append(b.job.Features, "gpu") // Auto-add gpu feature
	}

	// Validate high memory jobs
	if b.job.ResourceRequests.Memory > 32*GB && !b.hasFeature("highmem") {
		return fmt.Errorf("jobs requiring > 32GB memory must use highmem feature")
	}

	return nil
}

// hasFeature checks if a feature is set
func (b *JobBuilder) hasFeature(feature string) bool {
	for _, f := range b.job.Features {
		if f == feature {
			return true
		}
	}
	return false
}

// JobResourceRequestsBuilder provides a fluent interface for building resource requests
type JobResourceRequestsBuilder struct {
	parent    *JobBuilder
	resources *types.ResourceRequests
}

// WithMemory sets the total memory requirement in bytes
func (r *JobResourceRequestsBuilder) WithMemory(bytes int64) *JobResourceRequestsBuilder {
	if bytes < 0 {
		r.parent.addError(fmt.Errorf("memory must be non-negative, got %d", bytes))
		return r
	}
	r.resources.Memory = bytes
	return r
}

// WithMemoryMB sets the total memory requirement in MB
func (r *JobResourceRequestsBuilder) WithMemoryMB(mb int64) *JobResourceRequestsBuilder {
	return r.WithMemory(mb * MB)
}

// WithMemoryGB sets the total memory requirement in GB
func (r *JobResourceRequestsBuilder) WithMemoryGB(gb int64) *JobResourceRequestsBuilder {
	return r.WithMemory(gb * GB)
}

// WithMemoryPerCPU sets the memory per CPU in bytes
func (r *JobResourceRequestsBuilder) WithMemoryPerCPU(bytes int64) *JobResourceRequestsBuilder {
	if bytes < 0 {
		r.parent.addError(fmt.Errorf("memory per CPU must be non-negative, got %d", bytes))
		return r
	}
	r.resources.MemoryPerCPU = bytes
	return r
}

// WithMemoryPerNode sets the memory per node in bytes  
func (r *JobResourceRequestsBuilder) WithMemoryPerNode(bytes int64) *JobResourceRequestsBuilder {
	return r.WithMemory(bytes)
}

// WithMemoryPerGPU sets the memory per GPU in bytes
func (r *JobResourceRequestsBuilder) WithMemoryPerGPU(bytes int64) *JobResourceRequestsBuilder {
	if bytes < 0 {
		r.parent.addError(fmt.Errorf("memory per GPU must be non-negative, got %d", bytes))
		return r
	}
	r.resources.MemoryPerGPU = bytes
	return r
}

// WithTmpDisk sets the temporary disk space requirement in bytes
func (r *JobResourceRequestsBuilder) WithTmpDisk(bytes int64) *JobResourceRequestsBuilder {
	if bytes < 0 {
		r.parent.addError(fmt.Errorf("tmp disk must be non-negative, got %d", bytes))
		return r
	}
	r.resources.TmpDisk = bytes
	return r
}

// WithCPUsPerTask sets the number of CPUs per task
func (r *JobResourceRequestsBuilder) WithCPUsPerTask(cpus int32) *JobResourceRequestsBuilder {
	if cpus <= 0 {
		r.parent.addError(fmt.Errorf("CPUs per task must be positive, got %d", cpus))
		return r
	}
	r.resources.CPUsPerTask = cpus
	return r
}

// WithTasksPerNode sets the number of tasks per node
func (r *JobResourceRequestsBuilder) WithTasksPerNode(tasks int32) *JobResourceRequestsBuilder {
	if tasks <= 0 {
		r.parent.addError(fmt.Errorf("tasks per node must be positive, got %d", tasks))
		return r
	}
	r.resources.TasksPerNode = tasks
	return r
}

// WithTasksPerCore sets the number of tasks per core
func (r *JobResourceRequestsBuilder) WithTasksPerCore(tasks int32) *JobResourceRequestsBuilder {
	if tasks <= 0 {
		r.parent.addError(fmt.Errorf("tasks per core must be positive, got %d", tasks))
		return r
	}
	r.resources.TasksPerCore = tasks
	return r
}

// WithThreadsPerCore sets the number of threads per core
func (r *JobResourceRequestsBuilder) WithThreadsPerCore(threads int32) *JobResourceRequestsBuilder {
	if threads <= 0 {
		r.parent.addError(fmt.Errorf("threads per core must be positive, got %d", threads))
		return r
	}
	r.resources.ThreadsPerCore = threads
	return r
}

// Done returns to the parent Job builder
func (r *JobResourceRequestsBuilder) Done() *JobBuilder {
	return r.parent
}

