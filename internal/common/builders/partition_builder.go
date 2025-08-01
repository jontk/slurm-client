// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package builders

import (
	"fmt"

	"github.com/jontk/slurm-client/internal/common/types"
)

// PartitionBuilder provides a fluent interface for building Partition objects
type PartitionBuilder struct {
	partition *types.PartitionCreate
	errors    []error
}

// NewPartitionBuilder creates a new Partition builder with the required name
func NewPartitionBuilder(name string) *PartitionBuilder {
	return &PartitionBuilder{
		partition: &types.PartitionCreate{
			Name:                 name,
			State:                types.PartitionStateUp, // Default state
			Priority:             1,                       // Default priority
			MaxNodes:             1000,                    // Default max nodes
			MaxTime:              1440,                    // Default max time (24 hours)
			DefaultTime:          60,                      // Default time (1 hour)
			GraceTime:            0,                       // Default grace time
			AllowAccounts:        []string{},              // Initialize empty
			AllowGroups:          []string{},              // Initialize empty
			AllowQoS:             []string{},              // Initialize empty
			DenyAccounts:         []string{},              // Initialize empty
			DenyQoS:              []string{},              // Initialize empty
			PreemptMode:          []string{},              // Initialize empty
			SelectTypeParameters: []string{},              // Initialize empty
			JobDefaults:          make(map[string]string), // Initialize empty
		},
		errors: []error{},
	}
}

// WithAllocNodes sets the nodes allocated to this partition
func (b *PartitionBuilder) WithAllocNodes(nodes string) *PartitionBuilder {
	b.partition.AllocNodes = nodes
	return b
}

// WithAllowAccounts sets the accounts allowed to use this partition
func (b *PartitionBuilder) WithAllowAccounts(accounts ...string) *PartitionBuilder {
	b.partition.AllowAccounts = append(b.partition.AllowAccounts, accounts...)
	return b
}

// WithAllowAllocNodes sets nodes that can be allocated from this partition
func (b *PartitionBuilder) WithAllowAllocNodes(nodes string) *PartitionBuilder {
	b.partition.AllowAllocNodes = nodes
	return b
}

// WithAllowGroups sets the groups allowed to use this partition
func (b *PartitionBuilder) WithAllowGroups(groups ...string) *PartitionBuilder {
	b.partition.AllowGroups = append(b.partition.AllowGroups, groups...)
	return b
}

// WithAllowQoS sets the QoS allowed on this partition
func (b *PartitionBuilder) WithAllowQoS(qos ...string) *PartitionBuilder {
	b.partition.AllowQoS = append(b.partition.AllowQoS, qos...)
	return b
}

// WithDenyAccounts sets the accounts denied access to this partition
func (b *PartitionBuilder) WithDenyAccounts(accounts ...string) *PartitionBuilder {
	b.partition.DenyAccounts = append(b.partition.DenyAccounts, accounts...)
	return b
}

// WithDenyQoS sets the QoS denied on this partition
func (b *PartitionBuilder) WithDenyQoS(qos ...string) *PartitionBuilder {
	b.partition.DenyQoS = append(b.partition.DenyQoS, qos...)
	return b
}

// WithDefaultMemPerCPU sets the default memory per CPU in MB
func (b *PartitionBuilder) WithDefaultMemPerCPU(mb int64) *PartitionBuilder {
	if mb < 0 {
		b.addError(fmt.Errorf("default memory per CPU must be non-negative, got %d", mb))
		return b
	}
	b.partition.DefaultMemPerCPU = mb
	return b
}

// WithDefaultMemPerNode sets the default memory per node in MB
func (b *PartitionBuilder) WithDefaultMemPerNode(mb int64) *PartitionBuilder {
	if mb < 0 {
		b.addError(fmt.Errorf("default memory per node must be non-negative, got %d", mb))
		return b
	}
	b.partition.DefaultMemPerNode = mb
	return b
}

// WithDefaultTime sets the default time limit in minutes
func (b *PartitionBuilder) WithDefaultTime(minutes int32) *PartitionBuilder {
	if minutes <= 0 {
		b.addError(fmt.Errorf("default time must be positive, got %d", minutes))
		return b
	}
	b.partition.DefaultTime = minutes
	return b
}

// WithDefMemPerNode sets the def memory per node in MB (for compatibility)
func (b *PartitionBuilder) WithDefMemPerNode(mb int64) *PartitionBuilder {
	if mb < 0 {
		b.addError(fmt.Errorf("def memory per node must be non-negative, got %d", mb))
		return b
	}
	b.partition.DefMemPerNode = mb
	return b
}

// WithGraceTime sets the grace time in seconds
func (b *PartitionBuilder) WithGraceTime(seconds int32) *PartitionBuilder {
	if seconds < 0 {
		b.addError(fmt.Errorf("grace time must be non-negative, got %d", seconds))
		return b
	}
	b.partition.GraceTime = seconds
	return b
}

// WithMaxCPUsPerNode sets the maximum CPUs per node
func (b *PartitionBuilder) WithMaxCPUsPerNode(cpus int32) *PartitionBuilder {
	if cpus <= 0 {
		b.addError(fmt.Errorf("max CPUs per node must be positive, got %d", cpus))
		return b
	}
	b.partition.MaxCPUsPerNode = cpus
	return b
}

// WithMaxMemPerNode sets the maximum memory per node in MB
func (b *PartitionBuilder) WithMaxMemPerNode(mb int64) *PartitionBuilder {
	if mb <= 0 {
		b.addError(fmt.Errorf("max memory per node must be positive, got %d", mb))
		return b
	}
	b.partition.MaxMemPerNode = mb
	return b
}

// WithMaxMemPerCPU sets the maximum memory per CPU in MB
func (b *PartitionBuilder) WithMaxMemPerCPU(mb int64) *PartitionBuilder {
	if mb <= 0 {
		b.addError(fmt.Errorf("max memory per CPU must be positive, got %d", mb))
		return b
	}
	b.partition.MaxMemPerCPU = mb
	return b
}

// WithMaxNodes sets the maximum number of nodes
func (b *PartitionBuilder) WithMaxNodes(nodes int32) *PartitionBuilder {
	if nodes <= 0 {
		b.addError(fmt.Errorf("max nodes must be positive, got %d", nodes))
		return b
	}
	b.partition.MaxNodes = nodes
	return b
}

// WithMaxTime sets the maximum time limit in minutes
func (b *PartitionBuilder) WithMaxTime(minutes int32) *PartitionBuilder {
	if minutes <= 0 {
		b.addError(fmt.Errorf("max time must be positive, got %d", minutes))
		return b
	}
	b.partition.MaxTime = minutes
	return b
}

// WithMinNodes sets the minimum number of nodes
func (b *PartitionBuilder) WithMinNodes(nodes int32) *PartitionBuilder {
	if nodes < 0 {
		b.addError(fmt.Errorf("min nodes must be non-negative, got %d", nodes))
		return b
	}
	b.partition.MinNodes = nodes
	return b
}

// WithNodes sets the nodes assigned to this partition
func (b *PartitionBuilder) WithNodes(nodes string) *PartitionBuilder {
	b.partition.Nodes = nodes
	return b
}

// WithOverTimeLimit sets the over time limit in minutes
func (b *PartitionBuilder) WithOverTimeLimit(minutes int32) *PartitionBuilder {
	if minutes < 0 {
		b.addError(fmt.Errorf("over time limit must be non-negative, got %d", minutes))
		return b
	}
	b.partition.OverTimeLimit = minutes
	return b
}

// WithPreemptMode sets the preemption modes
func (b *PartitionBuilder) WithPreemptMode(modes ...string) *PartitionBuilder {
	b.partition.PreemptMode = append(b.partition.PreemptMode, modes...)
	return b
}

// WithPriority sets the partition priority
func (b *PartitionBuilder) WithPriority(priority int32) *PartitionBuilder {
	if priority < 0 {
		b.addError(fmt.Errorf("priority must be non-negative, got %d", priority))
		return b
	}
	b.partition.Priority = priority
	return b
}

// WithPriorityJobFactor sets the priority job factor
func (b *PartitionBuilder) WithPriorityJobFactor(factor int32) *PartitionBuilder {
	if factor < 0 {
		b.addError(fmt.Errorf("priority job factor must be non-negative, got %d", factor))
		return b
	}
	b.partition.PriorityJobFactor = factor
	return b
}

// WithPriorityTier sets the priority tier
func (b *PartitionBuilder) WithPriorityTier(tier int32) *PartitionBuilder {
	if tier < 0 {
		b.addError(fmt.Errorf("priority tier must be non-negative, got %d", tier))
		return b
	}
	b.partition.PriorityTier = tier
	return b
}

// WithQoS sets the default Quality of Service for this partition
func (b *PartitionBuilder) WithQoS(qos string) *PartitionBuilder {
	b.partition.QoS = qos
	return b
}

// WithState sets the partition state
func (b *PartitionBuilder) WithState(state types.PartitionState) *PartitionBuilder {
	b.partition.State = state
	return b
}

// WithTresStr sets the TRES string
func (b *PartitionBuilder) WithTresStr(tres string) *PartitionBuilder {
	b.partition.TresStr = tres
	return b
}

// WithBillingWeightStr sets the billing weight string
func (b *PartitionBuilder) WithBillingWeightStr(weight string) *PartitionBuilder {
	b.partition.BillingWeightStr = weight
	return b
}

// WithSelectTypeParameters sets the select type parameters
func (b *PartitionBuilder) WithSelectTypeParameters(params ...string) *PartitionBuilder {
	b.partition.SelectTypeParameters = append(b.partition.SelectTypeParameters, params...)
	return b
}

// WithJobDefaults sets job default parameters
func (b *PartitionBuilder) WithJobDefaults(defaults map[string]string) *PartitionBuilder {
	if b.partition.JobDefaults == nil {
		b.partition.JobDefaults = make(map[string]string)
	}
	for k, v := range defaults {
		b.partition.JobDefaults[k] = v
	}
	return b
}

// WithJobDefault sets a single job default parameter
func (b *PartitionBuilder) WithJobDefault(key, value string) *PartitionBuilder {
	if b.partition.JobDefaults == nil {
		b.partition.JobDefaults = make(map[string]string)
	}
	b.partition.JobDefaults[key] = value
	return b
}

// WithResumeTimeout sets the resume timeout in seconds
func (b *PartitionBuilder) WithResumeTimeout(seconds int32) *PartitionBuilder {
	if seconds < 0 {
		b.addError(fmt.Errorf("resume timeout must be non-negative, got %d", seconds))
		return b
	}
	b.partition.ResumeTimeout = seconds
	return b
}

// WithSuspendTime sets the suspend time in seconds
func (b *PartitionBuilder) WithSuspendTime(seconds int32) *PartitionBuilder {
	if seconds < 0 {
		b.addError(fmt.Errorf("suspend time must be non-negative, got %d", seconds))
		return b
	}
	b.partition.SuspendTime = seconds
	return b
}

// WithSuspendTimeout sets the suspend timeout in seconds
func (b *PartitionBuilder) WithSuspendTimeout(seconds int32) *PartitionBuilder {
	if seconds < 0 {
		b.addError(fmt.Errorf("suspend timeout must be non-negative, got %d", seconds))
		return b
	}
	b.partition.SuspendTimeout = seconds
	return b
}

// AsHidden sets the partition as hidden
func (b *PartitionBuilder) AsHidden() *PartitionBuilder {
	b.partition.Hidden = true
	return b
}

// AsExclusiveUser sets the partition to require exclusive user access
func (b *PartitionBuilder) AsExclusiveUser() *PartitionBuilder {
	b.partition.ExclusiveUser = true
	return b
}

// AsLLN sets the partition to use Least Loaded Node scheduling
func (b *PartitionBuilder) AsLLN() *PartitionBuilder {
	b.partition.LLN = true
	return b
}

// AsRootOnly sets the partition to root-only access
func (b *PartitionBuilder) AsRootOnly() *PartitionBuilder {
	b.partition.RootOnly = true
	return b
}

// AsReqResv sets the partition to require reservations
func (b *PartitionBuilder) AsReqResv() *PartitionBuilder {
	b.partition.ReqResv = true
	return b
}

// AsPowerDownOnIdle sets the partition to power down nodes when idle
func (b *PartitionBuilder) AsPowerDownOnIdle() *PartitionBuilder {
	b.partition.PowerDownOnIdle = true
	return b
}

// AsDebugPartition applies common debug partition settings
func (b *PartitionBuilder) AsDebugPartition() *PartitionBuilder {
	return b.
		WithPriority(1000).        // High priority for debug
		WithMaxTime(30).           // 30 minutes max
		WithDefaultTime(10).       // 10 minutes default
		WithMaxNodes(2).           // Limited nodes
		WithPreemptMode("cancel"). // Cancel jobs for debugging
		WithQoS("debug")
}

// AsBatchPartition applies common batch partition settings
func (b *PartitionBuilder) AsBatchPartition() *PartitionBuilder {
	return b.
		WithPriority(50).           // Medium priority
		WithMaxTime(2880).          // 48 hours max
		WithDefaultTime(60).        // 1 hour default
		WithGraceTime(300).         // 5 minutes grace time
		WithPreemptMode("suspend"). // Suspend instead of cancel
		WithQoS("normal")
}

// AsInteractivePartition applies common interactive partition settings
func (b *PartitionBuilder) AsInteractivePartition() *PartitionBuilder {
	return b.
		WithPriority(500).         // High priority for interactive
		WithMaxTime(240).          // 4 hours max
		WithDefaultTime(30).       // 30 minutes default
		WithMaxNodes(4).           // Limited nodes for interactive
		WithPreemptMode("cancel"). // Quick cancellation
		WithQoS("interactive")
}

// AsGPUPartition applies common GPU partition settings
func (b *PartitionBuilder) AsGPUPartition() *PartitionBuilder {
	return b.
		WithPriority(100).          // Medium-high priority
		WithMaxTime(1440).          // 24 hours max
		WithDefaultTime(120).       // 2 hours default
		WithTresStr("gpu:1").       // GPU resources
		WithJobDefault("gres", "gpu:1"). // Default GPU requirement
		WithQoS("gpu")
}

// AsHighMemoryPartition applies common high memory partition settings
func (b *PartitionBuilder) AsHighMemoryPartition() *PartitionBuilder {
	return b.
		WithPriority(75).               // Medium priority
		WithMaxTime(2880).              // 48 hours max
		WithDefaultTime(240).           // 4 hours default
		WithMaxMemPerNode(512 * GB).    // 512GB max memory
		WithDefaultMemPerNode(64 * GB). // 64GB default memory
		WithJobDefault("mem", "64G").   // Default memory requirement
		WithQoS("highmem")
}

// AsMaintenancePartition applies common maintenance partition settings
func (b *PartitionBuilder) AsMaintenancePartition() *PartitionBuilder {
	return b.
		WithState(types.PartitionStateDrain). // Start in drain state
		WithPriority(10).                     // Low priority
		WithMaxTime(60).                      // 1 hour max
		AsRootOnly().                         // Root access only
		AsHidden()                            // Hidden from normal users
}

// Build validates and returns the built Partition
func (b *PartitionBuilder) Build() (*types.PartitionCreate, error) {
	// Check for accumulated errors
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("builder errors: %v", b.errors)
	}

	// Validate required fields
	if b.partition.Name == "" {
		return nil, fmt.Errorf("partition name is required")
	}

	// Apply business rules
	if err := b.validateBusinessRules(); err != nil {
		return nil, err
	}

	return b.partition, nil
}

// BuildForUpdate creates a PartitionUpdate object from the builder
func (b *PartitionBuilder) BuildForUpdate() (*types.PartitionUpdate, error) {
	// Check for accumulated errors
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("builder errors: %v", b.errors)
	}

	update := &types.PartitionUpdate{}

	// Only include fields that were explicitly set
	if b.partition.AllocNodes != "" {
		update.AllocNodes = &b.partition.AllocNodes
	}
	if len(b.partition.AllowAccounts) > 0 {
		update.AllowAccounts = b.partition.AllowAccounts
	}
	if b.partition.AllowAllocNodes != "" {
		update.AllowAllocNodes = &b.partition.AllowAllocNodes
	}
	if len(b.partition.AllowGroups) > 0 {
		update.AllowGroups = b.partition.AllowGroups
	}
	if len(b.partition.AllowQoS) > 0 {
		update.AllowQoS = b.partition.AllowQoS
	}
	if len(b.partition.DenyAccounts) > 0 {
		update.DenyAccounts = b.partition.DenyAccounts
	}
	if len(b.partition.DenyQoS) > 0 {
		update.DenyQoS = b.partition.DenyQoS
	}
	if b.partition.DefaultMemPerCPU != 0 {
		update.DefaultMemPerCPU = &b.partition.DefaultMemPerCPU
	}
	if b.partition.DefaultMemPerNode != 0 {
		update.DefaultMemPerNode = &b.partition.DefaultMemPerNode
	}
	if b.partition.DefaultTime != 60 { // Not default
		update.DefaultTime = &b.partition.DefaultTime
	}
	if b.partition.DefMemPerNode != 0 {
		update.DefMemPerNode = &b.partition.DefMemPerNode
	}
	if b.partition.GraceTime != 0 {
		update.GraceTime = &b.partition.GraceTime
	}
	if b.partition.MaxCPUsPerNode != 0 {
		update.MaxCPUsPerNode = &b.partition.MaxCPUsPerNode
	}
	if b.partition.MaxMemPerNode != 0 {
		update.MaxMemPerNode = &b.partition.MaxMemPerNode
	}
	if b.partition.MaxMemPerCPU != 0 {
		update.MaxMemPerCPU = &b.partition.MaxMemPerCPU
	}
	if b.partition.MaxNodes != 1000 { // Not default
		update.MaxNodes = &b.partition.MaxNodes
	}
	if b.partition.MaxTime != 1440 { // Not default
		update.MaxTime = &b.partition.MaxTime
	}
	if b.partition.MinNodes != 0 {
		update.MinNodes = &b.partition.MinNodes
	}
	if b.partition.Nodes != "" {
		update.Nodes = &b.partition.Nodes
	}
	if b.partition.OverTimeLimit != 0 {
		update.OverTimeLimit = &b.partition.OverTimeLimit
	}
	if len(b.partition.PreemptMode) > 0 {
		update.PreemptMode = b.partition.PreemptMode
	}
	if b.partition.Priority != 1 { // Not default
		update.Priority = &b.partition.Priority
	}
	if b.partition.PriorityJobFactor != 0 {
		update.PriorityJobFactor = &b.partition.PriorityJobFactor
	}
	if b.partition.PriorityTier != 0 {
		update.PriorityTier = &b.partition.PriorityTier
	}
	if b.partition.QoS != "" {
		update.QoS = &b.partition.QoS
	}
	if b.partition.State != types.PartitionStateUp { // Not default
		update.State = &b.partition.State
	}
	if b.partition.TresStr != "" {
		update.TresStr = &b.partition.TresStr
	}
	if b.partition.BillingWeightStr != "" {
		update.BillingWeightStr = &b.partition.BillingWeightStr
	}
	if len(b.partition.SelectTypeParameters) > 0 {
		update.SelectTypeParameters = b.partition.SelectTypeParameters
	}
	if len(b.partition.JobDefaults) > 0 {
		update.JobDefaults = b.partition.JobDefaults
	}
	if b.partition.ResumeTimeout != 0 {
		update.ResumeTimeout = &b.partition.ResumeTimeout
	}
	if b.partition.SuspendTime != 0 {
		update.SuspendTime = &b.partition.SuspendTime
	}
	if b.partition.SuspendTimeout != 0 {
		update.SuspendTimeout = &b.partition.SuspendTimeout
	}
	if b.partition.Hidden {
		update.Hidden = &b.partition.Hidden
	}
	if b.partition.ExclusiveUser {
		update.ExclusiveUser = &b.partition.ExclusiveUser
	}
	if b.partition.LLN {
		update.LLN = &b.partition.LLN
	}
	if b.partition.RootOnly {
		update.RootOnly = &b.partition.RootOnly
	}
	if b.partition.ReqResv {
		update.ReqResv = &b.partition.ReqResv
	}
	if b.partition.PowerDownOnIdle {
		update.PowerDownOnIdle = &b.partition.PowerDownOnIdle
	}

	return update, nil
}

// Clone creates a copy of the builder with the same settings
func (b *PartitionBuilder) Clone() *PartitionBuilder {
	newBuilder := &PartitionBuilder{
		partition: &types.PartitionCreate{
			Name:                 b.partition.Name,
			AllocNodes:           b.partition.AllocNodes,
			AllowAllocNodes:      b.partition.AllowAllocNodes,
			DefaultMemPerCPU:     b.partition.DefaultMemPerCPU,
			DefaultMemPerNode:    b.partition.DefaultMemPerNode,
			DefaultTime:          b.partition.DefaultTime,
			DefMemPerNode:        b.partition.DefMemPerNode,
			GraceTime:            b.partition.GraceTime,
			MaxCPUsPerNode:       b.partition.MaxCPUsPerNode,
			MaxMemPerNode:        b.partition.MaxMemPerNode,
			MaxMemPerCPU:         b.partition.MaxMemPerCPU,
			MaxNodes:             b.partition.MaxNodes,
			MaxTime:              b.partition.MaxTime,
			MinNodes:             b.partition.MinNodes,
			Nodes:                b.partition.Nodes,
			OverTimeLimit:        b.partition.OverTimeLimit,
			Priority:             b.partition.Priority,
			PriorityJobFactor:    b.partition.PriorityJobFactor,
			PriorityTier:         b.partition.PriorityTier,
			QoS:                  b.partition.QoS,
			State:                b.partition.State,
			TresStr:              b.partition.TresStr,
			BillingWeightStr:     b.partition.BillingWeightStr,
			ResumeTimeout:        b.partition.ResumeTimeout,
			SuspendTime:          b.partition.SuspendTime,
			SuspendTimeout:       b.partition.SuspendTimeout,
			Hidden:               b.partition.Hidden,
			ExclusiveUser:        b.partition.ExclusiveUser,
			LLN:                  b.partition.LLN,
			RootOnly:             b.partition.RootOnly,
			ReqResv:              b.partition.ReqResv,
			PowerDownOnIdle:      b.partition.PowerDownOnIdle,
		},
		errors: append([]error{}, b.errors...),
	}

	// Deep copy slices and maps
	newBuilder.partition.AllowAccounts = append([]string{}, b.partition.AllowAccounts...)
	newBuilder.partition.AllowGroups = append([]string{}, b.partition.AllowGroups...)
	newBuilder.partition.AllowQoS = append([]string{}, b.partition.AllowQoS...)
	newBuilder.partition.DenyAccounts = append([]string{}, b.partition.DenyAccounts...)
	newBuilder.partition.DenyQoS = append([]string{}, b.partition.DenyQoS...)
	newBuilder.partition.PreemptMode = append([]string{}, b.partition.PreemptMode...)
	newBuilder.partition.SelectTypeParameters = append([]string{}, b.partition.SelectTypeParameters...)

	if b.partition.JobDefaults != nil {
		newBuilder.partition.JobDefaults = make(map[string]string)
		for k, v := range b.partition.JobDefaults {
			newBuilder.partition.JobDefaults[k] = v
		}
	}

	return newBuilder
}

// addError adds an error to the builder's error list
func (b *PartitionBuilder) addError(err error) {
	b.errors = append(b.errors, err)
}

// validateBusinessRules applies business logic validation
func (b *PartitionBuilder) validateBusinessRules() error {
	// Validate time limits consistency
	if b.partition.DefaultTime > b.partition.MaxTime {
		return fmt.Errorf("default time (%d) cannot exceed max time (%d)", b.partition.DefaultTime, b.partition.MaxTime)
	}

	// Validate node limits consistency
	if b.partition.MinNodes > b.partition.MaxNodes {
		return fmt.Errorf("min nodes (%d) cannot exceed max nodes (%d)", b.partition.MinNodes, b.partition.MaxNodes)
	}

	// Validate memory limits consistency
	if b.partition.DefaultMemPerNode > 0 && b.partition.MaxMemPerNode > 0 {
		if b.partition.DefaultMemPerNode > b.partition.MaxMemPerNode {
			return fmt.Errorf("default memory per node (%d) cannot exceed max memory per node (%d)", 
				b.partition.DefaultMemPerNode, b.partition.MaxMemPerNode)
		}
	}

	// Validate conflicting account settings
	if len(b.partition.AllowAccounts) > 0 && len(b.partition.DenyAccounts) > 0 {
		for _, allow := range b.partition.AllowAccounts {
			for _, deny := range b.partition.DenyAccounts {
				if allow == deny {
					return fmt.Errorf("account %s cannot be both allowed and denied", allow)
				}
			}
		}
	}

	// Validate conflicting QoS settings
	if len(b.partition.AllowQoS) > 0 && len(b.partition.DenyQoS) > 0 {
		for _, allow := range b.partition.AllowQoS {
			for _, deny := range b.partition.DenyQoS {
				if allow == deny {
					return fmt.Errorf("QoS %s cannot be both allowed and denied", allow)
				}
			}
		}
	}

	// Validate maintenance partition constraints
	if b.partition.State == types.PartitionStateDrain && !b.partition.Hidden {
		return fmt.Errorf("drained partitions should typically be hidden")
	}

	return nil
}
