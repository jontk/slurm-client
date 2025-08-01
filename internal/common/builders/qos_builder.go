// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package builders

import (
	"fmt"
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
)

// QoSBuilder provides a fluent interface for building QoS objects
type QoSBuilder struct {
	qos    *types.QoSCreate
	errors []error
}

// NewQoSBuilder creates a new QoS builder with the required name
func NewQoSBuilder(name string) *QoSBuilder {
	return &QoSBuilder{
		qos: &types.QoSCreate{
			Name:           name,
			Priority:       0,        // Default priority
			UsageFactor:    1.0,      // Default usage factor
			UsageThreshold: 0,        // Default usage threshold
			Flags:          []string{}, // Initialize empty
			PreemptMode:    []string{}, // Initialize empty
		},
		errors: []error{},
	}
}

// WithDescription sets the QoS description
func (b *QoSBuilder) WithDescription(description string) *QoSBuilder {
	b.qos.Description = description
	return b
}

// WithPriority sets the QoS priority
func (b *QoSBuilder) WithPriority(priority int) *QoSBuilder {
	if priority < 0 {
		b.addError(fmt.Errorf("priority must be non-negative, got %d", priority))
		return b
	}
	b.qos.Priority = priority
	return b
}

// WithFlags sets the QoS flags
func (b *QoSBuilder) WithFlags(flags ...string) *QoSBuilder {
	b.qos.Flags = append(b.qos.Flags, flags...)
	return b
}

// WithPreemptMode sets the preemption mode
func (b *QoSBuilder) WithPreemptMode(modes ...string) *QoSBuilder {
	b.qos.PreemptMode = append(b.qos.PreemptMode, modes...)
	return b
}

// WithPreemptExemptTime sets the preemption exempt time in minutes
func (b *QoSBuilder) WithPreemptExemptTime(minutes int) *QoSBuilder {
	if minutes < 0 {
		b.addError(fmt.Errorf("preempt exempt time must be non-negative, got %d", minutes))
		return b
	}
	b.qos.PreemptExemptTime = &minutes
	return b
}

// WithGraceTime sets the grace time in seconds
func (b *QoSBuilder) WithGraceTime(seconds int) *QoSBuilder {
	if seconds < 0 {
		b.addError(fmt.Errorf("grace time must be non-negative, got %d", seconds))
		return b
	}
	b.qos.GraceTime = seconds
	return b
}

// WithUsageFactor sets the usage factor
func (b *QoSBuilder) WithUsageFactor(factor float64) *QoSBuilder {
	if factor < 0 {
		b.addError(fmt.Errorf("usage factor must be non-negative, got %f", factor))
		return b
	}
	b.qos.UsageFactor = factor
	return b
}

// WithUsageThreshold sets the usage threshold
func (b *QoSBuilder) WithUsageThreshold(threshold float64) *QoSBuilder {
	if threshold < 0 || threshold > 1 {
		b.addError(fmt.Errorf("usage threshold must be between 0 and 1, got %f", threshold))
		return b
	}
	b.qos.UsageThreshold = threshold
	return b
}

// WithLimits returns a limits builder for setting resource limits
func (b *QoSBuilder) WithLimits() *QoSLimitsBuilder {
	if b.qos.Limits == nil {
		b.qos.Limits = &types.QoSLimits{}
	}
	return &QoSLimitsBuilder{
		parent: b,
		limits: b.qos.Limits,
	}
}

// AsHighPriority applies common high-priority settings
func (b *QoSBuilder) AsHighPriority() *QoSBuilder {
	return b.
		WithPriority(1000).
		WithUsageFactor(2.0).
		WithUsageThreshold(0.95).
		WithFlags("DenyOnLimit", "RequiresReservation").
		WithPreemptMode("cluster")
}

// AsBatchQueue applies common batch queue settings
func (b *QoSBuilder) AsBatchQueue() *QoSBuilder {
	return b.
		WithPriority(10).
		WithUsageFactor(0.5).
		WithFlags("NoReserve").
		WithGraceTime(3600) // 1 hour grace time
}

// AsInteractive applies common interactive job settings
func (b *QoSBuilder) AsInteractive() *QoSBuilder {
	return b.
		WithPriority(500).
		WithUsageFactor(1.5).
		WithFlags("DenyOnLimit").
		WithPreemptMode("suspend")
}

// Build validates and returns the built QoS
func (b *QoSBuilder) Build() (*types.QoSCreate, error) {
	// Check for accumulated errors
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("builder errors: %v", b.errors)
	}

	// Validate required fields
	if b.qos.Name == "" {
		return nil, fmt.Errorf("QoS name is required")
	}

	// Apply business rules
	if err := b.validateBusinessRules(); err != nil {
		return nil, err
	}

	return b.qos, nil
}

// BuildForUpdate creates a QoSUpdate object from the builder
func (b *QoSBuilder) BuildForUpdate() (*types.QoSUpdate, error) {
	// Check for accumulated errors
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("builder errors: %v", b.errors)
	}

	update := &types.QoSUpdate{}

	// Only include fields that were explicitly set
	if b.qos.Description != "" {
		update.Description = &b.qos.Description
	}
	if b.qos.Priority != 0 {
		update.Priority = &b.qos.Priority
	}
	if len(b.qos.Flags) > 0 {
		update.Flags = &b.qos.Flags
	}
	if len(b.qos.PreemptMode) > 0 && b.qos.PreemptMode[0] != "" {
		update.PreemptMode = &b.qos.PreemptMode
	}
	if b.qos.PreemptExemptTime != nil {
		update.PreemptExemptTime = b.qos.PreemptExemptTime
	}
	if b.qos.GraceTime != 0 {
		update.GraceTime = &b.qos.GraceTime
	}
	if b.qos.UsageFactor != 1.0 {
		update.UsageFactor = &b.qos.UsageFactor
	}
	if b.qos.UsageThreshold != 0 {
		update.UsageThreshold = &b.qos.UsageThreshold
	}
	if b.qos.Limits != nil {
		update.Limits = b.qos.Limits
	}

	return update, nil
}

// Clone creates a copy of the builder with the same settings
func (b *QoSBuilder) Clone() *QoSBuilder {
	newBuilder := &QoSBuilder{
		qos: &types.QoSCreate{
			Name:           b.qos.Name,
			Description:    b.qos.Description,
			Priority:       b.qos.Priority,
			UsageFactor:    b.qos.UsageFactor,
			UsageThreshold: b.qos.UsageThreshold,
			Flags:          append([]string{}, b.qos.Flags...),
			PreemptMode:    append([]string{}, b.qos.PreemptMode...),
		},
		errors: append([]error{}, b.errors...),
	}

	if b.qos.PreemptExemptTime != nil {
		v := *b.qos.PreemptExemptTime
		newBuilder.qos.PreemptExemptTime = &v
	}
	newBuilder.qos.GraceTime = b.qos.GraceTime
	if b.qos.Limits != nil {
		// Deep copy limits
		newBuilder.qos.Limits = b.cloneLimits(b.qos.Limits)
	}

	return newBuilder
}

// addError adds an error to the builder's error list
func (b *QoSBuilder) addError(err error) {
	b.errors = append(b.errors, err)
}

// validateBusinessRules applies business logic validation
func (b *QoSBuilder) validateBusinessRules() error {
	// Example business rules
	if b.qos.Priority > 1000 && !b.hasFlag("RequiresReservation") {
		return fmt.Errorf("QoS with priority > 1000 must have RequiresReservation flag")
	}

	if b.qos.UsageFactor > 3.0 {
		return fmt.Errorf("usage factor cannot exceed 3.0 for safety reasons")
	}

	// Check for conflicting flags
	if b.hasFlag("NoReserve") && b.hasFlag("RequiresReservation") {
		return fmt.Errorf("conflicting flags: NoReserve and RequiresReservation")
	}

	return nil
}

// hasFlag checks if a flag is set
func (b *QoSBuilder) hasFlag(flag string) bool {
	for _, f := range b.qos.Flags {
		if f == flag {
			return true
		}
	}
	return false
}

// cloneLimits creates a deep copy of QoS limits
func (b *QoSBuilder) cloneLimits(src *types.QoSLimits) *types.QoSLimits {
	dst := &types.QoSLimits{}

	// Clone all pointer fields
	if src.MaxCPUsPerUser != nil {
		v := *src.MaxCPUsPerUser
		dst.MaxCPUsPerUser = &v
	}
	if src.MaxJobsPerUser != nil {
		v := *src.MaxJobsPerUser
		dst.MaxJobsPerUser = &v
	}
	if src.MaxNodesPerUser != nil {
		v := *src.MaxNodesPerUser
		dst.MaxNodesPerUser = &v
	}
	if src.MaxSubmitJobsPerUser != nil {
		v := *src.MaxSubmitJobsPerUser
		dst.MaxSubmitJobsPerUser = &v
	}
	if src.MaxCPUsPerAccount != nil {
		v := *src.MaxCPUsPerAccount
		dst.MaxCPUsPerAccount = &v
	}
	if src.MaxJobsPerAccount != nil {
		v := *src.MaxJobsPerAccount
		dst.MaxJobsPerAccount = &v
	}
	if src.MaxNodesPerAccount != nil {
		v := *src.MaxNodesPerAccount
		dst.MaxNodesPerAccount = &v
	}
	if src.MaxCPUsPerJob != nil {
		v := *src.MaxCPUsPerJob
		dst.MaxCPUsPerJob = &v
	}
	if src.MaxNodesPerJob != nil {
		v := *src.MaxNodesPerJob
		dst.MaxNodesPerJob = &v
	}
	if src.MaxWallTimePerJob != nil {
		v := *src.MaxWallTimePerJob
		dst.MaxWallTimePerJob = &v
	}
	if src.MaxMemoryPerNode != nil {
		v := *src.MaxMemoryPerNode
		dst.MaxMemoryPerNode = &v
	}
	if src.MaxMemoryPerCPU != nil {
		v := *src.MaxMemoryPerCPU
		dst.MaxMemoryPerCPU = &v
	}
	if src.MaxBurstBuffer != nil {
		v := *src.MaxBurstBuffer
		dst.MaxBurstBuffer = &v
	}
	if src.MinCPUsPerJob != nil {
		v := *src.MinCPUsPerJob
		dst.MinCPUsPerJob = &v
	}
	if src.MinNodesPerJob != nil {
		v := *src.MinNodesPerJob
		dst.MinNodesPerJob = &v
	}

	return dst
}

// QoSLimitsBuilder provides a fluent interface for building QoS limits
type QoSLimitsBuilder struct {
	parent *QoSBuilder
	limits *types.QoSLimits
}

// WithMaxCPUsPerUser sets the maximum CPUs per user
func (l *QoSLimitsBuilder) WithMaxCPUsPerUser(cpus int) *QoSLimitsBuilder {
	if cpus < 0 {
		l.parent.addError(fmt.Errorf("max CPUs per user must be non-negative, got %d", cpus))
		return l
	}
	l.limits.MaxCPUsPerUser = &cpus
	return l
}

// WithMaxJobsPerUser sets the maximum jobs per user
func (l *QoSLimitsBuilder) WithMaxJobsPerUser(jobs int) *QoSLimitsBuilder {
	if jobs < 0 {
		l.parent.addError(fmt.Errorf("max jobs per user must be non-negative, got %d", jobs))
		return l
	}
	l.limits.MaxJobsPerUser = &jobs
	return l
}

// WithMaxNodesPerUser sets the maximum nodes per user
func (l *QoSLimitsBuilder) WithMaxNodesPerUser(nodes int) *QoSLimitsBuilder {
	if nodes < 0 {
		l.parent.addError(fmt.Errorf("max nodes per user must be non-negative, got %d", nodes))
		return l
	}
	l.limits.MaxNodesPerUser = &nodes
	return l
}

// WithMaxSubmitJobsPerUser sets the maximum submitted jobs per user
func (l *QoSLimitsBuilder) WithMaxSubmitJobsPerUser(jobs int) *QoSLimitsBuilder {
	if jobs < 0 {
		l.parent.addError(fmt.Errorf("max submit jobs per user must be non-negative, got %d", jobs))
		return l
	}
	l.limits.MaxSubmitJobsPerUser = &jobs
	return l
}

// WithMaxCPUsPerAccount sets the maximum CPUs per account
func (l *QoSLimitsBuilder) WithMaxCPUsPerAccount(cpus int) *QoSLimitsBuilder {
	if cpus < 0 {
		l.parent.addError(fmt.Errorf("max CPUs per account must be non-negative, got %d", cpus))
		return l
	}
	l.limits.MaxCPUsPerAccount = &cpus
	return l
}

// WithMaxJobsPerAccount sets the maximum jobs per account
func (l *QoSLimitsBuilder) WithMaxJobsPerAccount(jobs int) *QoSLimitsBuilder {
	if jobs < 0 {
		l.parent.addError(fmt.Errorf("max jobs per account must be non-negative, got %d", jobs))
		return l
	}
	l.limits.MaxJobsPerAccount = &jobs
	return l
}

// WithMaxNodesPerAccount sets the maximum nodes per account
func (l *QoSLimitsBuilder) WithMaxNodesPerAccount(nodes int) *QoSLimitsBuilder {
	if nodes < 0 {
		l.parent.addError(fmt.Errorf("max nodes per account must be non-negative, got %d", nodes))
		return l
	}
	l.limits.MaxNodesPerAccount = &nodes
	return l
}

// WithMaxCPUsPerJob sets the maximum CPUs per job
func (l *QoSLimitsBuilder) WithMaxCPUsPerJob(cpus int) *QoSLimitsBuilder {
	if cpus < 0 {
		l.parent.addError(fmt.Errorf("max CPUs per job must be non-negative, got %d", cpus))
		return l
	}
	l.limits.MaxCPUsPerJob = &cpus
	return l
}

// WithMaxNodesPerJob sets the maximum nodes per job
func (l *QoSLimitsBuilder) WithMaxNodesPerJob(nodes int) *QoSLimitsBuilder {
	if nodes < 0 {
		l.parent.addError(fmt.Errorf("max nodes per job must be non-negative, got %d", nodes))
		return l
	}
	l.limits.MaxNodesPerJob = &nodes
	return l
}

// WithMaxWallTime sets the maximum wall time per job
func (l *QoSLimitsBuilder) WithMaxWallTime(duration time.Duration) *QoSLimitsBuilder {
	minutes := int(duration.Minutes())
	if minutes < 0 {
		l.parent.addError(fmt.Errorf("max wall time must be non-negative, got %v", duration))
		return l
	}
	l.limits.MaxWallTimePerJob = &minutes
	return l
}

// WithMaxMemoryPerNode sets the maximum memory per node in MB
func (l *QoSLimitsBuilder) WithMaxMemoryPerNode(mb int64) *QoSLimitsBuilder {
	if mb < 0 {
		l.parent.addError(fmt.Errorf("max memory per node must be non-negative, got %d", mb))
		return l
	}
	l.limits.MaxMemoryPerNode = &mb
	return l
}

// WithMaxMemoryPerCPU sets the maximum memory per CPU in MB
func (l *QoSLimitsBuilder) WithMaxMemoryPerCPU(mb int64) *QoSLimitsBuilder {
	if mb < 0 {
		l.parent.addError(fmt.Errorf("max memory per CPU must be non-negative, got %d", mb))
		return l
	}
	l.limits.MaxMemoryPerCPU = &mb
	return l
}

// WithMinCPUsPerJob sets the minimum CPUs per job
func (l *QoSLimitsBuilder) WithMinCPUsPerJob(cpus int) *QoSLimitsBuilder {
	if cpus < 0 {
		l.parent.addError(fmt.Errorf("min CPUs per job must be non-negative, got %d", cpus))
		return l
	}
	l.limits.MinCPUsPerJob = &cpus
	return l
}

// WithMinNodesPerJob sets the minimum nodes per job
func (l *QoSLimitsBuilder) WithMinNodesPerJob(nodes int) *QoSLimitsBuilder {
	if nodes < 0 {
		l.parent.addError(fmt.Errorf("min nodes per job must be non-negative, got %d", nodes))
		return l
	}
	l.limits.MinNodesPerJob = &nodes
	return l
}

// Done returns to the parent QoS builder
func (l *QoSLimitsBuilder) Done() *QoSBuilder {
	return l.parent
}

