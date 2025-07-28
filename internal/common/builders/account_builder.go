package builders

import (
	"fmt"

	"github.com/jontk/slurm-client/internal/common/types"
)

// AccountBuilder provides a fluent interface for building Account objects
type AccountBuilder struct {
	account *types.AccountCreate
	errors  []error
}

// NewAccountBuilder creates a new Account builder with the required name
func NewAccountBuilder(name string) *AccountBuilder {
	return &AccountBuilder{
		account: &types.AccountCreate{
			Name:            name,
			FairShare:       1,                         // Default fair share
			SharesRaw:       1,                         // Default raw shares
			Priority:        0,                         // Default priority
			Coordinators:    []string{},                // Initialize empty
			QoSList:         []string{},                // Initialize empty
			AllowedPartitions: []string{},              // Initialize empty
			GrpTRES:         make(map[string]int64),    // Initialize empty
			GrpTRESMins:     make(map[string]int64),    // Initialize empty
			GrpTRESRunMins:  make(map[string]int64),    // Initialize empty
			MaxTRES:         make(map[string]int64),    // Initialize empty
			MaxTRESPerNode:  make(map[string]int64),    // Initialize empty
			MinTRES:         make(map[string]int64),    // Initialize empty
		},
		errors: []error{},
	}
}

// WithDescription sets the account description
func (b *AccountBuilder) WithDescription(description string) *AccountBuilder {
	b.account.Description = description
	return b
}

// WithOrganization sets the organization name
func (b *AccountBuilder) WithOrganization(organization string) *AccountBuilder {
	b.account.Organization = organization
	return b
}

// WithCoordinators sets the account coordinators
func (b *AccountBuilder) WithCoordinators(coordinators ...string) *AccountBuilder {
	b.account.Coordinators = append(b.account.Coordinators, coordinators...)
	return b
}

// WithDefaultQoS sets the default Quality of Service
func (b *AccountBuilder) WithDefaultQoS(qos string) *AccountBuilder {
	b.account.DefaultQoS = qos
	return b
}

// WithQoSList sets the allowed QoS list
func (b *AccountBuilder) WithQoSList(qosList ...string) *AccountBuilder {
	b.account.QoSList = append(b.account.QoSList, qosList...)
	return b
}

// WithParentAccount sets the parent account name
func (b *AccountBuilder) WithParentAccount(parentName string) *AccountBuilder {
	b.account.ParentName = parentName
	return b
}

// WithAllowedPartitions sets the allowed partitions
func (b *AccountBuilder) WithAllowedPartitions(partitions ...string) *AccountBuilder {
	b.account.AllowedPartitions = append(b.account.AllowedPartitions, partitions...)
	return b
}

// WithDefaultPartition sets the default partition
func (b *AccountBuilder) WithDefaultPartition(partition string) *AccountBuilder {
	b.account.DefaultPartition = partition
	return b
}

// WithFairShare sets the fair share value
func (b *AccountBuilder) WithFairShare(shares int32) *AccountBuilder {
	if shares < 0 {
		b.addError(fmt.Errorf("fair share must be non-negative, got %d", shares))
		return b
	}
	b.account.FairShare = shares
	return b
}

// WithSharesRaw sets the raw shares value
func (b *AccountBuilder) WithSharesRaw(shares int32) *AccountBuilder {
	if shares < 0 {
		b.addError(fmt.Errorf("raw shares must be non-negative, got %d", shares))
		return b
	}
	b.account.SharesRaw = shares
	return b
}

// WithPriority sets the account priority
func (b *AccountBuilder) WithPriority(priority int32) *AccountBuilder {
	if priority < 0 {
		b.addError(fmt.Errorf("priority must be non-negative, got %d", priority))
		return b
	}
	b.account.Priority = priority
	return b
}

// WithMaxJobs sets the maximum number of jobs
func (b *AccountBuilder) WithMaxJobs(jobs int32) *AccountBuilder {
	if jobs < 0 {
		b.addError(fmt.Errorf("max jobs must be non-negative, got %d", jobs))
		return b
	}
	b.account.MaxJobs = jobs
	return b
}

// WithMaxJobsPerUser sets the maximum jobs per user
func (b *AccountBuilder) WithMaxJobsPerUser(jobs int32) *AccountBuilder {
	if jobs < 0 {
		b.addError(fmt.Errorf("max jobs per user must be non-negative, got %d", jobs))
		return b
	}
	b.account.MaxJobsPerUser = jobs
	return b
}

// WithMaxSubmitJobs sets the maximum submitted jobs
func (b *AccountBuilder) WithMaxSubmitJobs(jobs int32) *AccountBuilder {
	if jobs < 0 {
		b.addError(fmt.Errorf("max submit jobs must be non-negative, got %d", jobs))
		return b
	}
	b.account.MaxSubmitJobs = jobs
	return b
}

// WithMaxWallTime sets the maximum wall time in minutes
func (b *AccountBuilder) WithMaxWallTime(minutes int32) *AccountBuilder {
	if minutes < 0 {
		b.addError(fmt.Errorf("max wall time must be non-negative, got %d", minutes))
		return b
	}
	b.account.MaxWallTime = minutes
	return b
}

// WithMaxCPUTime sets the maximum CPU time in minutes
func (b *AccountBuilder) WithMaxCPUTime(minutes int32) *AccountBuilder {
	if minutes < 0 {
		b.addError(fmt.Errorf("max CPU time must be non-negative, got %d", minutes))
		return b
	}
	b.account.MaxCPUTime = minutes
	return b
}

// WithMaxNodes sets the maximum number of nodes
func (b *AccountBuilder) WithMaxNodes(nodes int32) *AccountBuilder {
	if nodes < 0 {
		b.addError(fmt.Errorf("max nodes must be non-negative, got %d", nodes))
		return b
	}
	b.account.MaxNodes = nodes
	return b
}

// WithMaxCPUs sets the maximum number of CPUs
func (b *AccountBuilder) WithMaxCPUs(cpus int32) *AccountBuilder {
	if cpus < 0 {
		b.addError(fmt.Errorf("max CPUs must be non-negative, got %d", cpus))
		return b
	}
	b.account.MaxCPUs = cpus
	return b
}

// WithMaxMemory sets the maximum memory in MB
func (b *AccountBuilder) WithMaxMemory(mb int64) *AccountBuilder {
	if mb < 0 {
		b.addError(fmt.Errorf("max memory must be non-negative, got %d", mb))
		return b
	}
	b.account.MaxMemory = mb
	return b
}

// WithMaxMemoryGB sets the maximum memory in GB
func (b *AccountBuilder) WithMaxMemoryGB(gb int64) *AccountBuilder {
	return b.WithMaxMemory(gb * GB)
}

// WithMinPriorityThreshold sets the minimum priority threshold
func (b *AccountBuilder) WithMinPriorityThreshold(threshold int32) *AccountBuilder {
	if threshold < 0 {
		b.addError(fmt.Errorf("min priority threshold must be non-negative, got %d", threshold))
		return b
	}
	b.account.MinPriorityThreshold = threshold
	return b
}

// WithGrpJobs sets the group jobs limit
func (b *AccountBuilder) WithGrpJobs(jobs int32) *AccountBuilder {
	if jobs < 0 {
		b.addError(fmt.Errorf("group jobs must be non-negative, got %d", jobs))
		return b
	}
	b.account.GrpJobs = jobs
	return b
}

// WithGrpJobsAccrue sets the group jobs accrue limit
func (b *AccountBuilder) WithGrpJobsAccrue(jobs int32) *AccountBuilder {
	if jobs < 0 {
		b.addError(fmt.Errorf("group jobs accrue must be non-negative, got %d", jobs))
		return b
	}
	b.account.GrpJobsAccrue = jobs
	return b
}

// WithGrpNodes sets the group nodes limit
func (b *AccountBuilder) WithGrpNodes(nodes int32) *AccountBuilder {
	if nodes < 0 {
		b.addError(fmt.Errorf("group nodes must be non-negative, got %d", nodes))
		return b
	}
	b.account.GrpNodes = nodes
	return b
}

// WithGrpCPUs sets the group CPUs limit
func (b *AccountBuilder) WithGrpCPUs(cpus int32) *AccountBuilder {
	if cpus < 0 {
		b.addError(fmt.Errorf("group CPUs must be non-negative, got %d", cpus))
		return b
	}
	b.account.GrpCPUs = cpus
	return b
}

// WithGrpMemory sets the group memory limit in MB
func (b *AccountBuilder) WithGrpMemory(mb int64) *AccountBuilder {
	if mb < 0 {
		b.addError(fmt.Errorf("group memory must be non-negative, got %d", mb))
		return b
	}
	b.account.GrpMemory = mb
	return b
}

// WithGrpMemoryGB sets the group memory limit in GB
func (b *AccountBuilder) WithGrpMemoryGB(gb int64) *AccountBuilder {
	return b.WithGrpMemory(gb * GB)
}

// WithGrpSubmitJobs sets the group submit jobs limit
func (b *AccountBuilder) WithGrpSubmitJobs(jobs int32) *AccountBuilder {
	if jobs < 0 {
		b.addError(fmt.Errorf("group submit jobs must be non-negative, got %d", jobs))
		return b
	}
	b.account.GrpSubmitJobs = jobs
	return b
}

// WithGrpWallTime sets the group wall time limit in minutes
func (b *AccountBuilder) WithGrpWallTime(minutes int32) *AccountBuilder {
	if minutes < 0 {
		b.addError(fmt.Errorf("group wall time must be non-negative, got %d", minutes))
		return b
	}
	b.account.GrpWallTime = minutes
	return b
}

// WithGrpCPUTime sets the group CPU time limit in minutes
func (b *AccountBuilder) WithGrpCPUTime(minutes int32) *AccountBuilder {
	if minutes < 0 {
		b.addError(fmt.Errorf("group CPU time must be non-negative, got %d", minutes))
		return b
	}
	b.account.GrpCPUTime = minutes
	return b
}

// WithGrpTRES sets a group TRES limit
func (b *AccountBuilder) WithGrpTRES(resource string, limit int64) *AccountBuilder {
	if limit < 0 {
		b.addError(fmt.Errorf("group TRES %s must be non-negative, got %d", resource, limit))
		return b
	}
	if b.account.GrpTRES == nil {
		b.account.GrpTRES = make(map[string]int64)
	}
	b.account.GrpTRES[resource] = limit
	return b
}

// WithGrpTRESMins sets a group TRES minutes limit
func (b *AccountBuilder) WithGrpTRESMins(resource string, limit int64) *AccountBuilder {
	if limit < 0 {
		b.addError(fmt.Errorf("group TRES mins %s must be non-negative, got %d", resource, limit))
		return b
	}
	if b.account.GrpTRESMins == nil {
		b.account.GrpTRESMins = make(map[string]int64)
	}
	b.account.GrpTRESMins[resource] = limit
	return b
}

// WithGrpTRESRunMins sets a group TRES run minutes limit
func (b *AccountBuilder) WithGrpTRESRunMins(resource string, limit int64) *AccountBuilder {
	if limit < 0 {
		b.addError(fmt.Errorf("group TRES run mins %s must be non-negative, got %d", resource, limit))
		return b
	}
	if b.account.GrpTRESRunMins == nil {
		b.account.GrpTRESRunMins = make(map[string]int64)
	}
	b.account.GrpTRESRunMins[resource] = limit
	return b
}

// WithMaxTRES sets a maximum TRES limit
func (b *AccountBuilder) WithMaxTRES(resource string, limit int64) *AccountBuilder {
	if limit < 0 {
		b.addError(fmt.Errorf("max TRES %s must be non-negative, got %d", resource, limit))
		return b
	}
	if b.account.MaxTRES == nil {
		b.account.MaxTRES = make(map[string]int64)
	}
	b.account.MaxTRES[resource] = limit
	return b
}

// WithMaxTRESPerNode sets a maximum TRES per node limit
func (b *AccountBuilder) WithMaxTRESPerNode(resource string, limit int64) *AccountBuilder {
	if limit < 0 {
		b.addError(fmt.Errorf("max TRES per node %s must be non-negative, got %d", resource, limit))
		return b
	}
	if b.account.MaxTRESPerNode == nil {
		b.account.MaxTRESPerNode = make(map[string]int64)
	}
	b.account.MaxTRESPerNode[resource] = limit
	return b
}

// WithMinTRES sets a minimum TRES limit
func (b *AccountBuilder) WithMinTRES(resource string, limit int64) *AccountBuilder {
	if limit < 0 {
		b.addError(fmt.Errorf("min TRES %s must be non-negative, got %d", resource, limit))
		return b
	}
	if b.account.MinTRES == nil {
		b.account.MinTRES = make(map[string]int64)
	}
	b.account.MinTRES[resource] = limit
	return b
}

// AsResearchAccount applies common research account settings
func (b *AccountBuilder) AsResearchAccount() *AccountBuilder {
	return b.
		WithDescription("Research account for academic work").
		WithFairShare(100).
		WithMaxJobs(1000).
		WithMaxJobsPerUser(50).
		WithMaxWallTime(10080). // 7 days
		WithGrpJobs(500).
		WithDefaultQoS("normal").
		WithQoSList("normal", "high")
}

// AsComputeAccount applies common compute account settings
func (b *AccountBuilder) AsComputeAccount() *AccountBuilder {
	return b.
		WithDescription("General compute account").
		WithFairShare(50).
		WithMaxJobs(500).
		WithMaxJobsPerUser(25).
		WithMaxWallTime(2880). // 2 days
		WithGrpJobs(200).
		WithDefaultQoS("normal").
		WithQoSList("normal")
}

// AsStudentAccount applies common student account settings
func (b *AccountBuilder) AsStudentAccount() *AccountBuilder {
	return b.
		WithDescription("Student account for coursework").
		WithFairShare(10).
		WithMaxJobs(50).
		WithMaxJobsPerUser(10).
		WithMaxWallTime(480). // 8 hours
		WithMaxNodes(4).
		WithMaxCPUs(64).
		WithGrpJobs(25).
		WithDefaultQoS("normal").
		WithQoSList("normal")
}

// AsGuestAccount applies common guest account settings
func (b *AccountBuilder) AsGuestAccount() *AccountBuilder {
	return b.
		WithDescription("Guest account with limited access").
		WithFairShare(1).
		WithMaxJobs(10).
		WithMaxJobsPerUser(5).
		WithMaxWallTime(60). // 1 hour
		WithMaxNodes(1).
		WithMaxCPUs(8).
		WithGrpJobs(5).
		WithDefaultQoS("debug").
		WithQoSList("debug")
}

// AsHighPerformanceAccount applies settings for high-performance computing
func (b *AccountBuilder) AsHighPerformanceAccount() *AccountBuilder {
	return b.
		WithDescription("High-performance computing account").
		WithFairShare(200).
		WithMaxJobs(2000).
		WithMaxJobsPerUser(100).
		WithMaxWallTime(20160). // 14 days
		WithMaxNodes(100).
		WithMaxCPUs(2000).
		WithGrpJobs(1000).
		WithGrpNodes(50).
		WithGrpCPUs(1000).
		WithDefaultQoS("high").
		WithQoSList("normal", "high", "urgent").
		WithMaxTRES("cpu", 2000).
		WithMaxTRES("mem", 2000*GB).
		WithMaxTRES("node", 100)
}

// AsServiceAccount applies settings for service accounts
func (b *AccountBuilder) AsServiceAccount() *AccountBuilder {
	return b.
		WithDescription("Service account for automated processes").
		WithFairShare(25).
		WithMaxJobs(100).
		WithMaxJobsPerUser(50).
		WithMaxWallTime(1440). // 1 day
		WithGrpJobs(50).
		WithDefaultQoS("normal").
		WithQoSList("normal")
}

// Build validates and returns the built Account
func (b *AccountBuilder) Build() (*types.AccountCreate, error) {
	// Check for accumulated errors
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("builder errors: %v", b.errors)
	}

	// Validate required fields
	if b.account.Name == "" {
		return nil, fmt.Errorf("account name is required")
	}

	// Apply business rules
	if err := b.validateBusinessRules(); err != nil {
		return nil, err
	}

	return b.account, nil
}

// BuildForUpdate creates an AccountUpdate object from the builder
func (b *AccountBuilder) BuildForUpdate() (*types.AccountUpdate, error) {
	// Check for accumulated errors
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("builder errors: %v", b.errors)
	}

	update := &types.AccountUpdate{}

	// Only include fields that were explicitly set
	if b.account.Description != "" {
		update.Description = &b.account.Description
	}
	if b.account.Organization != "" {
		update.Organization = &b.account.Organization
	}
	if len(b.account.Coordinators) > 0 {
		update.Coordinators = b.account.Coordinators
	}
	if b.account.DefaultQoS != "" {
		update.DefaultQoS = &b.account.DefaultQoS
	}
	if len(b.account.QoSList) > 0 {
		update.QoSList = b.account.QoSList
	}
	if len(b.account.AllowedPartitions) > 0 {
		update.AllowedPartitions = b.account.AllowedPartitions
	}
	if b.account.DefaultPartition != "" {
		update.DefaultPartition = &b.account.DefaultPartition
	}
	if b.account.FairShare != 1 { // Not default
		update.FairShare = &b.account.FairShare
	}
	if b.account.SharesRaw != 1 { // Not default
		update.SharesRaw = &b.account.SharesRaw
	}
	if b.account.Priority != 0 {
		update.Priority = &b.account.Priority
	}
	if b.account.MaxJobs != 0 {
		update.MaxJobs = &b.account.MaxJobs
	}
	if b.account.MaxJobsPerUser != 0 {
		update.MaxJobsPerUser = &b.account.MaxJobsPerUser
	}
	if b.account.MaxSubmitJobs != 0 {
		update.MaxSubmitJobs = &b.account.MaxSubmitJobs
	}
	if b.account.MaxWallTime != 0 {
		update.MaxWallTime = &b.account.MaxWallTime
	}
	if b.account.MaxCPUTime != 0 {
		update.MaxCPUTime = &b.account.MaxCPUTime
	}
	if b.account.MaxNodes != 0 {
		update.MaxNodes = &b.account.MaxNodes
	}
	if b.account.MaxCPUs != 0 {
		update.MaxCPUs = &b.account.MaxCPUs
	}
	if b.account.MaxMemory != 0 {
		update.MaxMemory = &b.account.MaxMemory
	}
	if b.account.MinPriorityThreshold != 0 {
		update.MinPriorityThreshold = &b.account.MinPriorityThreshold
	}
	if b.account.GrpJobs != 0 {
		update.GrpJobs = &b.account.GrpJobs
	}
	if b.account.GrpJobsAccrue != 0 {
		update.GrpJobsAccrue = &b.account.GrpJobsAccrue
	}
	if b.account.GrpNodes != 0 {
		update.GrpNodes = &b.account.GrpNodes
	}
	if b.account.GrpCPUs != 0 {
		update.GrpCPUs = &b.account.GrpCPUs
	}
	if b.account.GrpMemory != 0 {
		update.GrpMemory = &b.account.GrpMemory
	}
	if b.account.GrpSubmitJobs != 0 {
		update.GrpSubmitJobs = &b.account.GrpSubmitJobs
	}
	if b.account.GrpWallTime != 0 {
		update.GrpWallTime = &b.account.GrpWallTime
	}
	if b.account.GrpCPUTime != 0 {
		update.GrpCPUTime = &b.account.GrpCPUTime
	}
	if len(b.account.GrpTRES) > 0 {
		update.GrpTRES = b.account.GrpTRES
	}
	if len(b.account.GrpTRESMins) > 0 {
		update.GrpTRESMins = b.account.GrpTRESMins
	}
	if len(b.account.GrpTRESRunMins) > 0 {
		update.GrpTRESRunMins = b.account.GrpTRESRunMins
	}
	if len(b.account.MaxTRES) > 0 {
		update.MaxTRES = b.account.MaxTRES
	}
	if len(b.account.MaxTRESPerNode) > 0 {
		update.MaxTRESPerNode = b.account.MaxTRESPerNode
	}
	if len(b.account.MinTRES) > 0 {
		update.MinTRES = b.account.MinTRES
	}

	return update, nil
}

// Clone creates a copy of the builder with the same settings
func (b *AccountBuilder) Clone() *AccountBuilder {
	newBuilder := &AccountBuilder{
		account: &types.AccountCreate{
			Name:                 b.account.Name,
			Description:          b.account.Description,
			Organization:         b.account.Organization,
			DefaultQoS:           b.account.DefaultQoS,
			ParentName:           b.account.ParentName,
			DefaultPartition:     b.account.DefaultPartition,
			FairShare:            b.account.FairShare,
			SharesRaw:            b.account.SharesRaw,
			Priority:             b.account.Priority,
			MaxJobs:              b.account.MaxJobs,
			MaxJobsPerUser:       b.account.MaxJobsPerUser,
			MaxSubmitJobs:        b.account.MaxSubmitJobs,
			MaxWallTime:          b.account.MaxWallTime,
			MaxCPUTime:           b.account.MaxCPUTime,
			MaxNodes:             b.account.MaxNodes,
			MaxCPUs:              b.account.MaxCPUs,
			MaxMemory:            b.account.MaxMemory,
			MinPriorityThreshold: b.account.MinPriorityThreshold,
			GrpJobs:              b.account.GrpJobs,
			GrpJobsAccrue:        b.account.GrpJobsAccrue,
			GrpNodes:             b.account.GrpNodes,
			GrpCPUs:              b.account.GrpCPUs,
			GrpMemory:            b.account.GrpMemory,
			GrpSubmitJobs:        b.account.GrpSubmitJobs,
			GrpWallTime:          b.account.GrpWallTime,
			GrpCPUTime:           b.account.GrpCPUTime,
		},
		errors: append([]error{}, b.errors...),
	}

	// Deep copy slices and maps
	newBuilder.account.Coordinators = append([]string{}, b.account.Coordinators...)
	newBuilder.account.QoSList = append([]string{}, b.account.QoSList...)
	newBuilder.account.AllowedPartitions = append([]string{}, b.account.AllowedPartitions...)

	if b.account.GrpTRES != nil {
		newBuilder.account.GrpTRES = make(map[string]int64)
		for k, v := range b.account.GrpTRES {
			newBuilder.account.GrpTRES[k] = v
		}
	}

	if b.account.GrpTRESMins != nil {
		newBuilder.account.GrpTRESMins = make(map[string]int64)
		for k, v := range b.account.GrpTRESMins {
			newBuilder.account.GrpTRESMins[k] = v
		}
	}

	if b.account.GrpTRESRunMins != nil {
		newBuilder.account.GrpTRESRunMins = make(map[string]int64)
		for k, v := range b.account.GrpTRESRunMins {
			newBuilder.account.GrpTRESRunMins[k] = v
		}
	}

	if b.account.MaxTRES != nil {
		newBuilder.account.MaxTRES = make(map[string]int64)
		for k, v := range b.account.MaxTRES {
			newBuilder.account.MaxTRES[k] = v
		}
	}

	if b.account.MaxTRESPerNode != nil {
		newBuilder.account.MaxTRESPerNode = make(map[string]int64)
		for k, v := range b.account.MaxTRESPerNode {
			newBuilder.account.MaxTRESPerNode[k] = v
		}
	}

	if b.account.MinTRES != nil {
		newBuilder.account.MinTRES = make(map[string]int64)
		for k, v := range b.account.MinTRES {
			newBuilder.account.MinTRES[k] = v
		}
	}

	return newBuilder
}

// addError adds an error to the builder's error list
func (b *AccountBuilder) addError(err error) {
	b.errors = append(b.errors, err)
}

// validateBusinessRules applies business logic validation
func (b *AccountBuilder) validateBusinessRules() error {
	// Validate job limits consistency
	if b.account.MaxJobsPerUser > 0 && b.account.MaxJobs > 0 {
		if b.account.MaxJobsPerUser > b.account.MaxJobs {
			return fmt.Errorf("max jobs per user (%d) cannot exceed max jobs (%d)", 
				b.account.MaxJobsPerUser, b.account.MaxJobs)
		}
	}

	// Validate group vs individual limits
	if b.account.GrpJobs > 0 && b.account.MaxJobs > 0 {
		if b.account.GrpJobs > b.account.MaxJobs {
			return fmt.Errorf("group jobs (%d) should not exceed max jobs (%d)", 
				b.account.GrpJobs, b.account.MaxJobs)
		}
	}

	// Validate TRES consistency
	if cpuMax, exists := b.account.MaxTRES["cpu"]; exists {
		if b.account.MaxCPUs > 0 && cpuMax != int64(b.account.MaxCPUs) {
			return fmt.Errorf("MaxCPUs (%d) and MaxTRES[cpu] (%d) should be consistent", 
				b.account.MaxCPUs, cpuMax)
		}
	}

	// Validate memory consistency  
	if memMax, exists := b.account.MaxTRES["mem"]; exists {
		if b.account.MaxMemory > 0 && memMax != b.account.MaxMemory {
			return fmt.Errorf("MaxMemory (%d) and MaxTRES[mem] (%d) should be consistent", 
				b.account.MaxMemory, memMax)
		}
	}

	// Validate node consistency
	if nodeMax, exists := b.account.MaxTRES["node"]; exists {
		if b.account.MaxNodes > 0 && nodeMax != int64(b.account.MaxNodes) {
			return fmt.Errorf("MaxNodes (%d) and MaxTRES[node] (%d) should be consistent", 
				b.account.MaxNodes, nodeMax)
		}
	}

	// Validate parent-child relationship
	if b.account.ParentName != "" && b.account.ParentName == b.account.Name {
		return fmt.Errorf("account cannot be its own parent")
	}

	// Validate QoS consistency
	if b.account.DefaultQoS != "" && len(b.account.QoSList) > 0 {
		found := false
		for _, qos := range b.account.QoSList {
			if qos == b.account.DefaultQoS {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("default QoS %s must be in the allowed QoS list", b.account.DefaultQoS)
		}
	}

	// Validate partition consistency
	if b.account.DefaultPartition != "" && len(b.account.AllowedPartitions) > 0 {
		found := false
		for _, partition := range b.account.AllowedPartitions {
			if partition == b.account.DefaultPartition {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("default partition %s must be in the allowed partitions list", b.account.DefaultPartition)
		}
	}

	return nil
}