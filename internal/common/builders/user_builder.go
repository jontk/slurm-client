// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package builders

import (
	"fmt"

	"github.com/jontk/slurm-client/internal/common/types"
)

// UserBuilder provides a fluent interface for building User objects
type UserBuilder struct {
	user   *types.UserCreate
	errors []error
}

// NewUserBuilder creates a new User builder with the required name
func NewUserBuilder(name string) *UserBuilder {
	return &UserBuilder{
		user: &types.UserCreate{
			Name:           name,
			AdminLevel:     types.AdminLevelNone, // Default admin level
			Accounts:       []string{},            // Initialize empty
			QoSList:        []string{},            // Initialize empty
			GrpTRES:        make(map[string]int64),  // Initialize empty
			GrpTRESMins:    make(map[string]int64),  // Initialize empty
			GrpTRESRunMins: make(map[string]int64),  // Initialize empty
			MaxTRES:        make(map[string]int64),  // Initialize empty
			MaxTRESPerNode: make(map[string]int64),  // Initialize empty
			MinTRES:        make(map[string]int64),  // Initialize empty
			WCKeys:         []string{},              // Initialize empty
		},
		errors: []error{},
	}
}

// WithUID sets the user ID
func (b *UserBuilder) WithUID(uid int32) *UserBuilder {
	if uid < 0 {
		b.addError(fmt.Errorf("UID must be non-negative, got %d", uid))
		return b
	}
	b.user.UID = uid
	return b
}

// WithDefaultAccount sets the default account
func (b *UserBuilder) WithDefaultAccount(account string) *UserBuilder {
	b.user.DefaultAccount = account
	return b
}

// WithDefaultWCKey sets the default workload characterization key
func (b *UserBuilder) WithDefaultWCKey(wckey string) *UserBuilder {
	b.user.DefaultWCKey = wckey
	return b
}

// WithAdminLevel sets the administrative level
func (b *UserBuilder) WithAdminLevel(level types.AdminLevel) *UserBuilder {
	b.user.AdminLevel = level
	return b
}

// WithAccounts sets the accounts the user can access
func (b *UserBuilder) WithAccounts(accounts ...string) *UserBuilder {
	b.user.Accounts = append(b.user.Accounts, accounts...)
	return b
}

// WithDefaultQoS sets the default Quality of Service
func (b *UserBuilder) WithDefaultQoS(qos string) *UserBuilder {
	b.user.DefaultQoS = qos
	return b
}

// WithQoSList sets the allowed QoS list
func (b *UserBuilder) WithQoSList(qosList ...string) *UserBuilder {
	b.user.QoSList = append(b.user.QoSList, qosList...)
	return b
}

// WithMaxJobs sets the maximum number of jobs
func (b *UserBuilder) WithMaxJobs(jobs int32) *UserBuilder {
	if jobs < 0 {
		b.addError(fmt.Errorf("max jobs must be non-negative, got %d", jobs))
		return b
	}
	b.user.MaxJobs = jobs
	return b
}

// WithMaxJobsPerAccount sets the maximum jobs per account
func (b *UserBuilder) WithMaxJobsPerAccount(jobs int32) *UserBuilder {
	if jobs < 0 {
		b.addError(fmt.Errorf("max jobs per account must be non-negative, got %d", jobs))
		return b
	}
	b.user.MaxJobsPerAccount = jobs
	return b
}

// WithMaxSubmitJobs sets the maximum submitted jobs
func (b *UserBuilder) WithMaxSubmitJobs(jobs int32) *UserBuilder {
	if jobs < 0 {
		b.addError(fmt.Errorf("max submit jobs must be non-negative, got %d", jobs))
		return b
	}
	b.user.MaxSubmitJobs = jobs
	return b
}

// WithMaxWallTime sets the maximum wall time in minutes
func (b *UserBuilder) WithMaxWallTime(minutes int32) *UserBuilder {
	if minutes < 0 {
		b.addError(fmt.Errorf("max wall time must be non-negative, got %d", minutes))
		return b
	}
	b.user.MaxWallTime = minutes
	return b
}

// WithMaxCPUTime sets the maximum CPU time in minutes
func (b *UserBuilder) WithMaxCPUTime(minutes int32) *UserBuilder {
	if minutes < 0 {
		b.addError(fmt.Errorf("max CPU time must be non-negative, got %d", minutes))
		return b
	}
	b.user.MaxCPUTime = minutes
	return b
}

// WithMaxNodes sets the maximum number of nodes
func (b *UserBuilder) WithMaxNodes(nodes int32) *UserBuilder {
	if nodes < 0 {
		b.addError(fmt.Errorf("max nodes must be non-negative, got %d", nodes))
		return b
	}
	b.user.MaxNodes = nodes
	return b
}

// WithMaxCPUs sets the maximum number of CPUs
func (b *UserBuilder) WithMaxCPUs(cpus int32) *UserBuilder {
	if cpus < 0 {
		b.addError(fmt.Errorf("max CPUs must be non-negative, got %d", cpus))
		return b
	}
	b.user.MaxCPUs = cpus
	return b
}

// WithMaxMemory sets the maximum memory in MB
func (b *UserBuilder) WithMaxMemory(mb int64) *UserBuilder {
	if mb < 0 {
		b.addError(fmt.Errorf("max memory must be non-negative, got %d", mb))
		return b
	}
	b.user.MaxMemory = mb
	return b
}

// WithMaxMemoryGB sets the maximum memory in GB
func (b *UserBuilder) WithMaxMemoryGB(gb int64) *UserBuilder {
	return b.WithMaxMemory(gb * GB)
}

// WithMinPriorityThreshold sets the minimum priority threshold
func (b *UserBuilder) WithMinPriorityThreshold(threshold int32) *UserBuilder {
	if threshold < 0 {
		b.addError(fmt.Errorf("min priority threshold must be non-negative, got %d", threshold))
		return b
	}
	b.user.MinPriorityThreshold = threshold
	return b
}

// WithGrpJobs sets the group jobs limit
func (b *UserBuilder) WithGrpJobs(jobs int32) *UserBuilder {
	if jobs < 0 {
		b.addError(fmt.Errorf("group jobs must be non-negative, got %d", jobs))
		return b
	}
	b.user.GrpJobs = jobs
	return b
}

// WithGrpJobsAccrue sets the group jobs accrue limit
func (b *UserBuilder) WithGrpJobsAccrue(jobs int32) *UserBuilder {
	if jobs < 0 {
		b.addError(fmt.Errorf("group jobs accrue must be non-negative, got %d", jobs))
		return b
	}
	b.user.GrpJobsAccrue = jobs
	return b
}

// WithGrpNodes sets the group nodes limit
func (b *UserBuilder) WithGrpNodes(nodes int32) *UserBuilder {
	if nodes < 0 {
		b.addError(fmt.Errorf("group nodes must be non-negative, got %d", nodes))
		return b
	}
	b.user.GrpNodes = nodes
	return b
}

// WithGrpCPUs sets the group CPUs limit
func (b *UserBuilder) WithGrpCPUs(cpus int32) *UserBuilder {
	if cpus < 0 {
		b.addError(fmt.Errorf("group CPUs must be non-negative, got %d", cpus))
		return b
	}
	b.user.GrpCPUs = cpus
	return b
}

// WithGrpMemory sets the group memory limit in MB
func (b *UserBuilder) WithGrpMemory(mb int64) *UserBuilder {
	if mb < 0 {
		b.addError(fmt.Errorf("group memory must be non-negative, got %d", mb))
		return b
	}
	b.user.GrpMemory = mb
	return b
}

// WithGrpMemoryGB sets the group memory limit in GB
func (b *UserBuilder) WithGrpMemoryGB(gb int64) *UserBuilder {
	return b.WithGrpMemory(gb * GB)
}

// WithGrpSubmitJobs sets the group submit jobs limit
func (b *UserBuilder) WithGrpSubmitJobs(jobs int32) *UserBuilder {
	if jobs < 0 {
		b.addError(fmt.Errorf("group submit jobs must be non-negative, got %d", jobs))
		return b
	}
	b.user.GrpSubmitJobs = jobs
	return b
}

// WithGrpWallTime sets the group wall time limit in minutes
func (b *UserBuilder) WithGrpWallTime(minutes int32) *UserBuilder {
	if minutes < 0 {
		b.addError(fmt.Errorf("group wall time must be non-negative, got %d", minutes))
		return b
	}
	b.user.GrpWallTime = minutes
	return b
}

// WithGrpCPUTime sets the group CPU time limit in minutes
func (b *UserBuilder) WithGrpCPUTime(minutes int32) *UserBuilder {
	if minutes < 0 {
		b.addError(fmt.Errorf("group CPU time must be non-negative, got %d", minutes))
		return b
	}
	b.user.GrpCPUTime = minutes
	return b
}

// WithGrpTRES sets a group TRES limit
func (b *UserBuilder) WithGrpTRES(resource string, limit int64) *UserBuilder {
	if limit < 0 {
		b.addError(fmt.Errorf("group TRES %s must be non-negative, got %d", resource, limit))
		return b
	}
	if b.user.GrpTRES == nil {
		b.user.GrpTRES = make(map[string]int64)
	}
	b.user.GrpTRES[resource] = limit
	return b
}

// WithGrpTRESMins sets a group TRES minutes limit
func (b *UserBuilder) WithGrpTRESMins(resource string, limit int64) *UserBuilder {
	if limit < 0 {
		b.addError(fmt.Errorf("group TRES mins %s must be non-negative, got %d", resource, limit))
		return b
	}
	if b.user.GrpTRESMins == nil {
		b.user.GrpTRESMins = make(map[string]int64)
	}
	b.user.GrpTRESMins[resource] = limit
	return b
}

// WithGrpTRESRunMins sets a group TRES run minutes limit
func (b *UserBuilder) WithGrpTRESRunMins(resource string, limit int64) *UserBuilder {
	if limit < 0 {
		b.addError(fmt.Errorf("group TRES run mins %s must be non-negative, got %d", resource, limit))
		return b
	}
	if b.user.GrpTRESRunMins == nil {
		b.user.GrpTRESRunMins = make(map[string]int64)
	}
	b.user.GrpTRESRunMins[resource] = limit
	return b
}

// WithMaxTRES sets a maximum TRES limit
func (b *UserBuilder) WithMaxTRES(resource string, limit int64) *UserBuilder {
	if limit < 0 {
		b.addError(fmt.Errorf("max TRES %s must be non-negative, got %d", resource, limit))
		return b
	}
	if b.user.MaxTRES == nil {
		b.user.MaxTRES = make(map[string]int64)
	}
	b.user.MaxTRES[resource] = limit
	return b
}

// WithMaxTRESPerNode sets a maximum TRES per node limit
func (b *UserBuilder) WithMaxTRESPerNode(resource string, limit int64) *UserBuilder {
	if limit < 0 {
		b.addError(fmt.Errorf("max TRES per node %s must be non-negative, got %d", resource, limit))
		return b
	}
	if b.user.MaxTRESPerNode == nil {
		b.user.MaxTRESPerNode = make(map[string]int64)
	}
	b.user.MaxTRESPerNode[resource] = limit
	return b
}

// WithMinTRES sets a minimum TRES limit
func (b *UserBuilder) WithMinTRES(resource string, limit int64) *UserBuilder {
	if limit < 0 {
		b.addError(fmt.Errorf("min TRES %s must be non-negative, got %d", resource, limit))
		return b
	}
	if b.user.MinTRES == nil {
		b.user.MinTRES = make(map[string]int64)
	}
	b.user.MinTRES[resource] = limit
	return b
}

// WithWCKeys sets the workload characterization keys
func (b *UserBuilder) WithWCKeys(wckeys ...string) *UserBuilder {
	b.user.WCKeys = append(b.user.WCKeys, wckeys...)
	return b
}

// AsAdministrator applies common administrator settings
func (b *UserBuilder) AsAdministrator() *UserBuilder {
	return b.
		WithAdminLevel(types.AdminLevelAdministrator).
		WithMaxJobs(10000).
		WithMaxJobsPerAccount(1000).
		WithMaxSubmitJobs(20000).
		WithMaxWallTime(43200). // 30 days
		WithDefaultQoS("high").
		WithQoSList("debug", "normal", "high", "urgent")
}

// AsOperator applies common operator settings
func (b *UserBuilder) AsOperator() *UserBuilder {
	return b.
		WithAdminLevel(types.AdminLevelOperator).
		WithMaxJobs(5000).
		WithMaxJobsPerAccount(500).
		WithMaxSubmitJobs(10000).
		WithMaxWallTime(20160). // 14 days
		WithDefaultQoS("normal").
		WithQoSList("debug", "normal", "high")
}

// AsRegularUser applies common regular user settings
func (b *UserBuilder) AsRegularUser() *UserBuilder {
	return b.
		WithAdminLevel(types.AdminLevelNone).
		WithMaxJobs(1000).
		WithMaxJobsPerAccount(100).
		WithMaxSubmitJobs(2000).
		WithMaxWallTime(10080). // 7 days
		WithMaxNodes(50).
		WithMaxCPUs(1000).
		WithDefaultQoS("normal").
		WithQoSList("normal", "high")
}

// AsStudentUser applies common student user settings
func (b *UserBuilder) AsStudentUser() *UserBuilder {
	return b.
		WithAdminLevel(types.AdminLevelNone).
		WithMaxJobs(100).
		WithMaxJobsPerAccount(25).
		WithMaxSubmitJobs(200).
		WithMaxWallTime(1440). // 1 day
		WithMaxNodes(4).
		WithMaxCPUs(64).
		WithMaxMemoryGB(256).
		WithDefaultQoS("normal").
		WithQoSList("normal")
}

// AsGuestUser applies common guest user settings
func (b *UserBuilder) AsGuestUser() *UserBuilder {
	return b.
		WithAdminLevel(types.AdminLevelNone).
		WithMaxJobs(10).
		WithMaxJobsPerAccount(5).
		WithMaxSubmitJobs(20).
		WithMaxWallTime(240). // 4 hours
		WithMaxNodes(1).
		WithMaxCPUs(8).
		WithMaxMemoryGB(32).
		WithDefaultQoS("debug").
		WithQoSList("debug")
}

// AsServiceUser applies common service user settings
func (b *UserBuilder) AsServiceUser() *UserBuilder {
	return b.
		WithAdminLevel(types.AdminLevelNone).
		WithMaxJobs(500).
		WithMaxJobsPerAccount(100).
		WithMaxSubmitJobs(1000).
		WithMaxWallTime(2880). // 2 days
		WithMaxNodes(10).
		WithMaxCPUs(200).
		WithDefaultQoS("normal").
		WithQoSList("normal")
}

// AsResearcher applies common researcher settings
func (b *UserBuilder) AsResearcher() *UserBuilder {
	return b.
		WithAdminLevel(types.AdminLevelNone).
		WithMaxJobs(2000).
		WithMaxJobsPerAccount(200).
		WithMaxSubmitJobs(4000).
		WithMaxWallTime(20160). // 14 days
		WithMaxNodes(100).
		WithMaxCPUs(2000).
		WithMaxMemoryGB(2000).
		WithDefaultQoS("normal").
		WithQoSList("normal", "high", "urgent").
		WithMaxTRES("cpu", 2000).
		WithMaxTRES("mem", 2000*GB).
		WithMaxTRES("gpu", 10)
}

// AsHighPerformanceUser applies settings for high-performance computing users
func (b *UserBuilder) AsHighPerformanceUser() *UserBuilder {
	return b.
		WithAdminLevel(types.AdminLevelNone).
		WithMaxJobs(5000).
		WithMaxJobsPerAccount(500).
		WithMaxSubmitJobs(10000).
		WithMaxWallTime(43200). // 30 days
		WithMaxNodes(500).
		WithMaxCPUs(5000).
		WithMaxMemoryGB(10000).
		WithGrpJobs(2000).
		WithGrpNodes(200).
		WithGrpCPUs(2000).
		WithGrpMemoryGB(5000).
		WithDefaultQoS("high").
		WithQoSList("normal", "high", "urgent").
		WithMaxTRES("cpu", 5000).
		WithMaxTRES("mem", 10000*GB).
		WithMaxTRES("node", 500).
		WithMaxTRES("gpu", 50)
}

// Build validates and returns the built User
func (b *UserBuilder) Build() (*types.UserCreate, error) {
	// Check for accumulated errors
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("builder errors: %v", b.errors)
	}

	// Validate required fields
	if b.user.Name == "" {
		return nil, fmt.Errorf("user name is required")
	}

	// Apply business rules
	if err := b.validateBusinessRules(); err != nil {
		return nil, err
	}

	return b.user, nil
}

// BuildForUpdate creates a UserUpdate object from the builder
func (b *UserBuilder) BuildForUpdate() (*types.UserUpdate, error) {
	// Check for accumulated errors
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("builder errors: %v", b.errors)
	}

	update := &types.UserUpdate{}

	// Only include fields that were explicitly set
	if b.user.DefaultAccount != "" {
		update.DefaultAccount = &b.user.DefaultAccount
	}
	if b.user.DefaultWCKey != "" {
		update.DefaultWCKey = &b.user.DefaultWCKey
	}
	if b.user.AdminLevel != types.AdminLevelNone { // Not default
		update.AdminLevel = &b.user.AdminLevel
	}
	if len(b.user.Accounts) > 0 {
		update.Accounts = b.user.Accounts
	}
	if b.user.DefaultQoS != "" {
		update.DefaultQoS = &b.user.DefaultQoS
	}
	if len(b.user.QoSList) > 0 {
		update.QoSList = b.user.QoSList
	}
	if b.user.MaxJobs != 0 {
		update.MaxJobs = &b.user.MaxJobs
	}
	if b.user.MaxJobsPerAccount != 0 {
		update.MaxJobsPerAccount = &b.user.MaxJobsPerAccount
	}
	if b.user.MaxSubmitJobs != 0 {
		update.MaxSubmitJobs = &b.user.MaxSubmitJobs
	}
	if b.user.MaxWallTime != 0 {
		update.MaxWallTime = &b.user.MaxWallTime
	}
	if b.user.MaxCPUTime != 0 {
		update.MaxCPUTime = &b.user.MaxCPUTime
	}
	if b.user.MaxNodes != 0 {
		update.MaxNodes = &b.user.MaxNodes
	}
	if b.user.MaxCPUs != 0 {
		update.MaxCPUs = &b.user.MaxCPUs
	}
	if b.user.MaxMemory != 0 {
		update.MaxMemory = &b.user.MaxMemory
	}
	if b.user.MinPriorityThreshold != 0 {
		update.MinPriorityThreshold = &b.user.MinPriorityThreshold
	}
	if b.user.GrpJobs != 0 {
		update.GrpJobs = &b.user.GrpJobs
	}
	if b.user.GrpJobsAccrue != 0 {
		update.GrpJobsAccrue = &b.user.GrpJobsAccrue
	}
	if b.user.GrpNodes != 0 {
		update.GrpNodes = &b.user.GrpNodes
	}
	if b.user.GrpCPUs != 0 {
		update.GrpCPUs = &b.user.GrpCPUs
	}
	if b.user.GrpMemory != 0 {
		update.GrpMemory = &b.user.GrpMemory
	}
	if b.user.GrpSubmitJobs != 0 {
		update.GrpSubmitJobs = &b.user.GrpSubmitJobs
	}
	if b.user.GrpWallTime != 0 {
		update.GrpWallTime = &b.user.GrpWallTime
	}
	if b.user.GrpCPUTime != 0 {
		update.GrpCPUTime = &b.user.GrpCPUTime
	}
	if len(b.user.GrpTRES) > 0 {
		update.GrpTRES = b.user.GrpTRES
	}
	if len(b.user.GrpTRESMins) > 0 {
		update.GrpTRESMins = b.user.GrpTRESMins
	}
	if len(b.user.GrpTRESRunMins) > 0 {
		update.GrpTRESRunMins = b.user.GrpTRESRunMins
	}
	if len(b.user.MaxTRES) > 0 {
		update.MaxTRES = b.user.MaxTRES
	}
	if len(b.user.MaxTRESPerNode) > 0 {
		update.MaxTRESPerNode = b.user.MaxTRESPerNode
	}
	if len(b.user.MinTRES) > 0 {
		update.MinTRES = b.user.MinTRES
	}
	if len(b.user.WCKeys) > 0 {
		update.WCKeys = b.user.WCKeys
	}

	return update, nil
}

// Clone creates a copy of the builder with the same settings
func (b *UserBuilder) Clone() *UserBuilder {
	newBuilder := &UserBuilder{
		user: &types.UserCreate{
			Name:                 b.user.Name,
			UID:                  b.user.UID,
			DefaultAccount:       b.user.DefaultAccount,
			DefaultWCKey:         b.user.DefaultWCKey,
			AdminLevel:           b.user.AdminLevel,
			DefaultQoS:           b.user.DefaultQoS,
			MaxJobs:              b.user.MaxJobs,
			MaxJobsPerAccount:    b.user.MaxJobsPerAccount,
			MaxSubmitJobs:        b.user.MaxSubmitJobs,
			MaxWallTime:          b.user.MaxWallTime,
			MaxCPUTime:           b.user.MaxCPUTime,
			MaxNodes:             b.user.MaxNodes,
			MaxCPUs:              b.user.MaxCPUs,
			MaxMemory:            b.user.MaxMemory,
			MinPriorityThreshold: b.user.MinPriorityThreshold,
			GrpJobs:              b.user.GrpJobs,
			GrpJobsAccrue:        b.user.GrpJobsAccrue,
			GrpNodes:             b.user.GrpNodes,
			GrpCPUs:              b.user.GrpCPUs,
			GrpMemory:            b.user.GrpMemory,
			GrpSubmitJobs:        b.user.GrpSubmitJobs,
			GrpWallTime:          b.user.GrpWallTime,
			GrpCPUTime:           b.user.GrpCPUTime,
		},
		errors: append([]error{}, b.errors...),
	}

	// Deep copy slices and maps
	newBuilder.user.Accounts = append([]string{}, b.user.Accounts...)
	newBuilder.user.QoSList = append([]string{}, b.user.QoSList...)
	newBuilder.user.WCKeys = append([]string{}, b.user.WCKeys...)

	if b.user.GrpTRES != nil {
		newBuilder.user.GrpTRES = make(map[string]int64)
		for k, v := range b.user.GrpTRES {
			newBuilder.user.GrpTRES[k] = v
		}
	}

	if b.user.GrpTRESMins != nil {
		newBuilder.user.GrpTRESMins = make(map[string]int64)
		for k, v := range b.user.GrpTRESMins {
			newBuilder.user.GrpTRESMins[k] = v
		}
	}

	if b.user.GrpTRESRunMins != nil {
		newBuilder.user.GrpTRESRunMins = make(map[string]int64)
		for k, v := range b.user.GrpTRESRunMins {
			newBuilder.user.GrpTRESRunMins[k] = v
		}
	}

	if b.user.MaxTRES != nil {
		newBuilder.user.MaxTRES = make(map[string]int64)
		for k, v := range b.user.MaxTRES {
			newBuilder.user.MaxTRES[k] = v
		}
	}

	if b.user.MaxTRESPerNode != nil {
		newBuilder.user.MaxTRESPerNode = make(map[string]int64)
		for k, v := range b.user.MaxTRESPerNode {
			newBuilder.user.MaxTRESPerNode[k] = v
		}
	}

	if b.user.MinTRES != nil {
		newBuilder.user.MinTRES = make(map[string]int64)
		for k, v := range b.user.MinTRES {
			newBuilder.user.MinTRES[k] = v
		}
	}

	return newBuilder
}

// addError adds an error to the builder's error list
func (b *UserBuilder) addError(err error) {
	b.errors = append(b.errors, err)
}

// validateBusinessRules applies business logic validation
func (b *UserBuilder) validateBusinessRules() error {
	// Validate job limits consistency
	if b.user.MaxJobsPerAccount > 0 && b.user.MaxJobs > 0 {
		if b.user.MaxJobsPerAccount > b.user.MaxJobs {
			return fmt.Errorf("max jobs per account (%d) cannot exceed max jobs (%d)", 
				b.user.MaxJobsPerAccount, b.user.MaxJobs)
		}
	}

	// Validate group vs individual limits
	if b.user.GrpJobs > 0 && b.user.MaxJobs > 0 {
		if b.user.GrpJobs > b.user.MaxJobs {
			return fmt.Errorf("group jobs (%d) should not exceed max jobs (%d)", 
				b.user.GrpJobs, b.user.MaxJobs)
		}
	}

	// Validate TRES consistency
	if cpuMax, exists := b.user.MaxTRES["cpu"]; exists {
		if b.user.MaxCPUs > 0 && cpuMax != int64(b.user.MaxCPUs) {
			return fmt.Errorf("MaxCPUs (%d) and MaxTRES[cpu] (%d) should be consistent", 
				b.user.MaxCPUs, cpuMax)
		}
	}

	// Validate memory consistency  
	if memMax, exists := b.user.MaxTRES["mem"]; exists {
		if b.user.MaxMemory > 0 && memMax != b.user.MaxMemory {
			return fmt.Errorf("MaxMemory (%d) and MaxTRES[mem] (%d) should be consistent", 
				b.user.MaxMemory, memMax)
		}
	}

	// Validate node consistency
	if nodeMax, exists := b.user.MaxTRES["node"]; exists {
		if b.user.MaxNodes > 0 && nodeMax != int64(b.user.MaxNodes) {
			return fmt.Errorf("MaxNodes (%d) and MaxTRES[node] (%d) should be consistent", 
				b.user.MaxNodes, nodeMax)
		}
	}

	// Validate QoS consistency
	if b.user.DefaultQoS != "" && len(b.user.QoSList) > 0 {
		found := false
		for _, qos := range b.user.QoSList {
			if qos == b.user.DefaultQoS {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("default QoS %s must be in the allowed QoS list", b.user.DefaultQoS)
		}
	}

	// Validate account consistency
	if b.user.DefaultAccount != "" && len(b.user.Accounts) > 0 {
		found := false
		for _, account := range b.user.Accounts {
			if account == b.user.DefaultAccount {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("default account %s must be in the allowed accounts list", b.user.DefaultAccount)
		}
	}

	// Validate admin level permissions
	if b.user.AdminLevel == types.AdminLevelAdministrator {
		// Administrators should have high limits
		if b.user.MaxJobs > 0 && b.user.MaxJobs < 1000 {
			return fmt.Errorf("administrators should have at least 1000 max jobs, got %d", b.user.MaxJobs)
		}
	}

	return nil
}
