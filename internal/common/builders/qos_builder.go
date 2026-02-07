// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package builders

import (
	"fmt"
	"time"

	types "github.com/jontk/slurm-client/api"
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
			Priority:       0,          // Default priority
			UsageFactor:    1.0,        // Default usage factor
			UsageThreshold: 0,          // Default usage threshold
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
	if src == nil {
		return nil
	}
	dst := &types.QoSLimits{}

	// Clone top-level fields
	if src.Factor != nil {
		v := *src.Factor
		dst.Factor = &v
	}
	if src.GraceTime != nil {
		v := *src.GraceTime
		dst.GraceTime = &v
	}

	// Clone Max limits
	if src.Max != nil {
		dst.Max = b.cloneLimitsMax(src.Max)
	}

	// Clone Min limits
	if src.Min != nil {
		dst.Min = b.cloneLimitsMin(src.Min)
	}

	return dst
}

// cloneLimitsMax creates a deep copy of QoSLimitsMax
func (b *QoSBuilder) cloneLimitsMax(src *types.QoSLimitsMax) *types.QoSLimitsMax {
	if src == nil {
		return nil
	}
	dst := &types.QoSLimitsMax{}

	if src.ActiveJobs != nil {
		dst.ActiveJobs = &types.QoSLimitsMaxActiveJobs{}
		if src.ActiveJobs.Accruing != nil {
			v := *src.ActiveJobs.Accruing
			dst.ActiveJobs.Accruing = &v
		}
		if src.ActiveJobs.Count != nil {
			v := *src.ActiveJobs.Count
			dst.ActiveJobs.Count = &v
		}
	}

	if src.Jobs != nil {
		dst.Jobs = &types.QoSLimitsMaxJobs{}
		if src.Jobs.Count != nil {
			v := *src.Jobs.Count
			dst.Jobs.Count = &v
		}
		if src.Jobs.ActiveJobs != nil && src.Jobs.ActiveJobs.Per != nil {
			dst.Jobs.ActiveJobs = &types.QoSLimitsMaxJobsActiveJobs{
				Per: &types.QoSLimitsMaxJobsActiveJobsPer{},
			}
			if src.Jobs.ActiveJobs.Per.User != nil {
				v := *src.Jobs.ActiveJobs.Per.User
				dst.Jobs.ActiveJobs.Per.User = &v
			}
			if src.Jobs.ActiveJobs.Per.Account != nil {
				v := *src.Jobs.ActiveJobs.Per.Account
				dst.Jobs.ActiveJobs.Per.Account = &v
			}
		}
		if src.Jobs.Per != nil {
			dst.Jobs.Per = &types.QoSLimitsMaxJobsPer{}
			if src.Jobs.Per.User != nil {
				v := *src.Jobs.Per.User
				dst.Jobs.Per.User = &v
			}
			if src.Jobs.Per.Account != nil {
				v := *src.Jobs.Per.Account
				dst.Jobs.Per.Account = &v
			}
		}
	}

	if src.WallClock != nil && src.WallClock.Per != nil {
		dst.WallClock = &types.QoSLimitsMaxWallClock{
			Per: &types.QoSLimitsMaxWallClockPer{},
		}
		if src.WallClock.Per.Job != nil {
			v := *src.WallClock.Per.Job
			dst.WallClock.Per.Job = &v
		}
		if src.WallClock.Per.QoS != nil {
			v := *src.WallClock.Per.QoS
			dst.WallClock.Per.QoS = &v
		}
	}

	return dst
}

// cloneLimitsMin creates a deep copy of QoSLimitsMin
func (b *QoSBuilder) cloneLimitsMin(src *types.QoSLimitsMin) *types.QoSLimitsMin {
	if src == nil {
		return nil
	}
	dst := &types.QoSLimitsMin{}

	if src.PriorityThreshold != nil {
		v := *src.PriorityThreshold
		dst.PriorityThreshold = &v
	}

	return dst
}

// QoSLimitsBuilder provides a fluent interface for building QoS limits
type QoSLimitsBuilder struct {
	parent *QoSBuilder
	limits *types.QoSLimits
}

// ensureMaxJobs ensures the Max.Jobs structure is initialized
func (l *QoSLimitsBuilder) ensureMaxJobs() {
	if l.limits.Max == nil {
		l.limits.Max = &types.QoSLimitsMax{}
	}
	if l.limits.Max.Jobs == nil {
		l.limits.Max.Jobs = &types.QoSLimitsMaxJobs{}
	}
}

// ensureMaxJobsActiveJobsPer ensures the Max.Jobs.ActiveJobs.Per structure is initialized
func (l *QoSLimitsBuilder) ensureMaxJobsActiveJobsPer() {
	l.ensureMaxJobs()
	if l.limits.Max.Jobs.ActiveJobs == nil {
		l.limits.Max.Jobs.ActiveJobs = &types.QoSLimitsMaxJobsActiveJobs{}
	}
	if l.limits.Max.Jobs.ActiveJobs.Per == nil {
		l.limits.Max.Jobs.ActiveJobs.Per = &types.QoSLimitsMaxJobsActiveJobsPer{}
	}
}

// ensureMaxJobsPer ensures the Max.Jobs.Per structure is initialized
func (l *QoSLimitsBuilder) ensureMaxJobsPer() {
	l.ensureMaxJobs()
	if l.limits.Max.Jobs.Per == nil {
		l.limits.Max.Jobs.Per = &types.QoSLimitsMaxJobsPer{}
	}
}

// ensureMaxWallClockPer ensures the Max.WallClock.Per structure is initialized
func (l *QoSLimitsBuilder) ensureMaxWallClockPer() {
	if l.limits.Max == nil {
		l.limits.Max = &types.QoSLimitsMax{}
	}
	if l.limits.Max.WallClock == nil {
		l.limits.Max.WallClock = &types.QoSLimitsMaxWallClock{}
	}
	if l.limits.Max.WallClock.Per == nil {
		l.limits.Max.WallClock.Per = &types.QoSLimitsMaxWallClockPer{}
	}
}

// WithMaxJobsPerUser sets the maximum running jobs per user
func (l *QoSLimitsBuilder) WithMaxJobsPerUser(jobs int) *QoSLimitsBuilder {
	if jobs < 0 {
		l.parent.addError(fmt.Errorf("max jobs per user must be non-negative, got %d", jobs))
		return l
	}
	l.ensureMaxJobsActiveJobsPer()
	v := uint32(jobs)
	l.limits.Max.Jobs.ActiveJobs.Per.User = &v
	return l
}

// WithMaxJobsPerAccount sets the maximum running jobs per account
func (l *QoSLimitsBuilder) WithMaxJobsPerAccount(jobs int) *QoSLimitsBuilder {
	if jobs < 0 {
		l.parent.addError(fmt.Errorf("max jobs per account must be non-negative, got %d", jobs))
		return l
	}
	l.ensureMaxJobsActiveJobsPer()
	v := uint32(jobs)
	l.limits.Max.Jobs.ActiveJobs.Per.Account = &v
	return l
}

// WithMaxSubmitJobsPerUser sets the maximum submitted jobs per user
func (l *QoSLimitsBuilder) WithMaxSubmitJobsPerUser(jobs int) *QoSLimitsBuilder {
	if jobs < 0 {
		l.parent.addError(fmt.Errorf("max submit jobs per user must be non-negative, got %d", jobs))
		return l
	}
	l.ensureMaxJobsPer()
	v := uint32(jobs)
	l.limits.Max.Jobs.Per.User = &v
	return l
}

// WithMaxSubmitJobsPerAccount sets the maximum submitted jobs per account
func (l *QoSLimitsBuilder) WithMaxSubmitJobsPerAccount(jobs int) *QoSLimitsBuilder {
	if jobs < 0 {
		l.parent.addError(fmt.Errorf("max submit jobs per account must be non-negative, got %d", jobs))
		return l
	}
	l.ensureMaxJobsPer()
	v := uint32(jobs)
	l.limits.Max.Jobs.Per.Account = &v
	return l
}

// WithGrpJobs sets the maximum number of running jobs for this QoS
func (l *QoSLimitsBuilder) WithGrpJobs(jobs int) *QoSLimitsBuilder {
	if jobs < 0 {
		l.parent.addError(fmt.Errorf("grp jobs must be non-negative, got %d", jobs))
		return l
	}
	if l.limits.Max == nil {
		l.limits.Max = &types.QoSLimitsMax{}
	}
	if l.limits.Max.ActiveJobs == nil {
		l.limits.Max.ActiveJobs = &types.QoSLimitsMaxActiveJobs{}
	}
	v := uint32(jobs)
	l.limits.Max.ActiveJobs.Count = &v
	return l
}

// WithGrpSubmitJobs sets the maximum submitted jobs for this QoS
func (l *QoSLimitsBuilder) WithGrpSubmitJobs(jobs int) *QoSLimitsBuilder {
	if jobs < 0 {
		l.parent.addError(fmt.Errorf("grp submit jobs must be non-negative, got %d", jobs))
		return l
	}
	l.ensureMaxJobs()
	v := uint32(jobs)
	l.limits.Max.Jobs.Count = &v
	return l
}

// WithMaxWallTime sets the maximum wall time per job
func (l *QoSLimitsBuilder) WithMaxWallTime(duration time.Duration) *QoSLimitsBuilder {
	minutes := int(duration.Minutes())
	if minutes < 0 {
		l.parent.addError(fmt.Errorf("max wall time must be non-negative, got %v", duration))
		return l
	}
	l.ensureMaxWallClockPer()
	v := uint32(minutes)
	l.limits.Max.WallClock.Per.Job = &v
	return l
}

// WithGrpWallTime sets the maximum wall time for all jobs in this QoS
func (l *QoSLimitsBuilder) WithGrpWallTime(duration time.Duration) *QoSLimitsBuilder {
	minutes := int(duration.Minutes())
	if minutes < 0 {
		l.parent.addError(fmt.Errorf("grp wall time must be non-negative, got %v", duration))
		return l
	}
	l.ensureMaxWallClockPer()
	v := uint32(minutes)
	l.limits.Max.WallClock.Per.QoS = &v
	return l
}

// WithGraceTime sets the preemption grace time in seconds
func (l *QoSLimitsBuilder) WithGraceTime(seconds int) *QoSLimitsBuilder {
	if seconds < 0 {
		l.parent.addError(fmt.Errorf("grace time must be non-negative, got %d", seconds))
		return l
	}
	v := int32(seconds)
	l.limits.GraceTime = &v
	return l
}

// WithFactor sets the limit factor for TRES
func (l *QoSLimitsBuilder) WithFactor(factor float64) *QoSLimitsBuilder {
	if factor < 0 {
		l.parent.addError(fmt.Errorf("factor must be non-negative, got %f", factor))
		return l
	}
	l.limits.Factor = &factor
	return l
}

// WithMinPriorityThreshold sets the minimum priority threshold
func (l *QoSLimitsBuilder) WithMinPriorityThreshold(threshold int) *QoSLimitsBuilder {
	if threshold < 0 {
		l.parent.addError(fmt.Errorf("min priority threshold must be non-negative, got %d", threshold))
		return l
	}
	if l.limits.Min == nil {
		l.limits.Min = &types.QoSLimitsMin{}
	}
	v := uint32(threshold)
	l.limits.Min.PriorityThreshold = &v
	return l
}

// Done returns to the parent QoS builder
func (l *QoSLimitsBuilder) Done() *QoSBuilder {
	return l.parent
}
