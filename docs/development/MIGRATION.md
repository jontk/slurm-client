# Migration Guide

This guide helps you migrate between SLURM REST API versions supported by this library.

## Table of Contents

- [Quick Reference](#quick-reference)
- [v0.0.40 → v0.0.42](#v0040--v0042)
- [v0.0.42 → v0.0.43](#v0042--v0043)
- [v0.0.43 → v0.0.44](#v0043--v0044)
- [Testing Your Migration](#testing-your-migration)
- [Common Pitfalls](#common-pitfalls)

## Quick Reference

### Breaking Changes Summary

| Change | Affected Versions | Severity | Migration Effort |
|--------|------------------|----------|------------------|
| Field renames (switches) | v0.0.40→v0.0.41 | Low | Find & replace |
| Removed fields (exclusive) | v0.0.41→v0.0.42 | Medium | Update code logic |
| NoVal type changes | v0.0.40→v0.0.42+ | Internal | Automatic (via builders) |
| New features (QoS, Reservations) | v0.0.42→v0.0.43 | N/A | Opt-in |
| Account management | v0.0.42→v0.0.43 | N/A | Opt-in |

### Version Feature Matrix

| Feature | v0.0.40 | v0.0.41 | v0.0.42 | v0.0.43 | v0.0.44 |
|---------|---------|---------|---------|---------|---------|
| Job Management | ✅ | ✅ | ✅ | ✅ | ✅ |
| Node Management | ✅ | ✅ | ✅ | ✅ | ✅ |
| Partition Management | ✅ | ✅ | ✅ | ✅ | ✅ |
| Cluster Info | ✅ | ✅ | ✅ | ✅ | ✅ |
| Reservation Management | ❌ | ❌ | ❌ | ✅ | ✅ |
| QoS Management | ❌ | ❌ | ❌ | ✅ | ✅ |
| Account Management | ❌ | ❌ | ❌ | ✅ | ✅ |
| Association Management | ❌ | ❌ | ❌ | ✅ | ✅ |
| User Management | ❌ | ❌ | ❌ | ✅ | ✅ |

## v0.0.40 → v0.0.42

### Overview

**Recommended:** Skip v0.0.41 and migrate directly from v0.0.40 to v0.0.42.

**Why?** v0.0.41 uses inline schemas which are not well-supported by code generation tools. v0.0.42 is the stable baseline with proper schema definitions.

### Breaking Changes

#### 1. Field Renames

**Switch-related fields renamed:**

```go
// v0.0.40
job.MinimumSwitches  // Old field name

// v0.0.42
job.RequiredSwitches // New field name
```

**Migration:**
```bash
# Find and replace in your codebase
grep -r "MinimumSwitches" . | # Find all occurrences
sed -i 's/MinimumSwitches/RequiredSwitches/g' # Replace
```

#### 2. Removed Fields

**Removed from job responses:**
- `exclusive` - Use partition configuration instead
- `oversubscribe` - Use partition oversubscribe settings

```go
// v0.0.40 - These fields existed
if *job.Exclusive {
    fmt.Println("Job has exclusive access")
}

// v0.0.42 - Fields removed, check partition instead
partition, _ := client.Partitions().Get(ctx, job.Partition)
if partition.Exclusive {
    fmt.Println("Partition provides exclusive access")
}
```

#### 3. NoVal Type Changes (Internal)

**Impact:** Mostly transparent if using builders

The internal representation changed from `_no_val` to `_no_val_struct`:

```go
// v0.0.40 internal type
type V0040Uint32NoVal struct {
    Set    *bool
    Number *uint32
}

// v0.0.42 internal type
type V0042Uint32NoValStruct struct {
    Set      *bool
    Number   *uint32
    Infinite *bool  // New field
}
```

**Using Builders (Recommended):** No code changes needed
```go
// Works in both v0.0.40 and v0.0.42
job := builderv0040.NewJobInfo().
    WithCpus(4).
    WithMemoryPerNode(8 * 1024 * 1024 * 1024).
    Build()

job := builderv0042.NewJobInfo().
    WithCpus(4).
    WithMemoryPerNode(8 * 1024 * 1024 * 1024).
    Build()
```

**Direct Access (Not Recommended):** Manual updates required
```go
// v0.0.40
job.Cpus = &v0_0_40.V0040Uint32NoVal{
    Set:    &setTrue,
    Number: &cpuCount,
}

// v0.0.42 - Must use new type
job.Cpus = &v0_0_42.V0042Uint32NoValStruct{
    Set:      &setTrue,
    Number:   &cpuCount,
    Infinite: nil,  // New field
}
```

### Migration Steps

**Step 1: Update your imports**
```go
// Before
import v0_0_40 "github.com/jontk/slurm-client/internal/api/v0_0_40"

// After
import v0_0_42 "github.com/jontk/slurm-client/internal/api/v0_0_42"
```

**Step 2: Update field names**
```bash
# Find all switch-related fields
grep -rn "MinimumSwitches" your-project/

# Update to new names
sed -i 's/MinimumSwitches/RequiredSwitches/g' your-project/**/*.go
```

**Step 3: Remove references to deleted fields**
```bash
# Find all references to removed fields
grep -rn "\.Exclusive\|\.Oversubscribe" your-project/

# Update to use partition settings instead
```

**Step 4: Regenerate mock builders (if using)**
```bash
# Regenerate builders for v0.0.42
make generate-mocks VERSION=v0.0.42

# Update test imports
sed -i 's/builderv0040/builderv0042/g' tests/**/*_test.go
```

**Step 5: Test thoroughly**
```bash
# Run tests against v0.0.42
go test ./...

# Run integration tests if available
SLURM_API_VERSION=v0.0.42 go test ./tests/integration/...
```

### Example: Complete Migration

**Before (v0.0.40):**
```go
package main

import (
    "context"
    v0_0_40 "github.com/jontk/slurm-client/internal/api/v0_0_40"
    builderv0040 "github.com/jontk/slurm-client/tests/mocks/generated/v0_0_40"
)

func submitJob(ctx context.Context, client *v0_0_40.Client) {
    job := builderv0040.NewJobInfo().
        WithJobId(1001).
        WithName("my-job").
        WithCpus(4).
        WithMinimumSwitches(2).  // Old field name
        Build()

    // Check exclusive access
    if job.Exclusive != nil && *job.Exclusive {
        log.Println("Job has exclusive node access")
    }
}
```

**After (v0.0.42):**
```go
package main

import (
    "context"
    v0_0_42 "github.com/jontk/slurm-client/internal/api/v0_0_42"
    builderv0042 "github.com/jontk/slurm-client/tests/mocks/generated/v0_0_42"
)

func submitJob(ctx context.Context, client *v0_0_42.Client) {
    job := builderv0042.NewJobInfo().
        WithJobId(1001).
        WithName("my-job").
        WithCpus(4).
        WithRequiredSwitches(2).  // New field name
        Build()

    // Check exclusive access via partition
    partition, err := client.Partitions().Get(ctx, *job.Partition)
    if err == nil && partition.Exclusive {
        log.Println("Job runs on exclusive partition")
    }
}
```

## v0.0.42 → v0.0.43

### Overview

v0.0.43 adds significant new features while maintaining backward compatibility for core operations.

**New Features:**
- ✅ Reservation Management
- ✅ Quality of Service (QoS)
- ✅ Account Management
- ✅ Association Management
- ✅ User Management

### New Features (Opt-in)

These features are **additive** - existing code continues to work without changes.

#### 1. Reservation Management

```go
// Check if supported
if client.Reservations() == nil {
    log.Println("Reservations not supported in this API version")
    return
}

// Create a maintenance window
reservation := &interfaces.ReservationCreate{
    Name:      "maintenance",
    StartTime: time.Now().Add(24 * time.Hour),
    Duration:  4 * 3600, // 4 hours
    Nodes:     []string{"node001", "node002"},
    Users:     []string{"admin"},
}

resp, err := client.Reservations().Create(ctx, reservation)
```

#### 2. QoS Management

```go
// Check if supported
if client.QoS() == nil {
    log.Println("QoS not supported in this API version")
    return
}

// Create high-priority QoS
qos := &interfaces.QoSCreate{
    Name:     "high-priority",
    Priority: 10000,
    MaxJobs:  50,
    MaxCPUs:  500,
}

resp, err := client.QoS().Create(ctx, qos)
```

#### 3. Account Management

```go
// Check if supported
if client.Accounts() == nil {
    log.Println("Account management not supported")
    return
}

// Create account hierarchy
account := &interfaces.AccountCreate{
    Name:         "research-dept",
    Organization: "ACME Corp",
    MaxJobs:      500,
    MaxNodes:     100,
}

resp, err := client.Accounts().Create(ctx, account)
```

### Breaking Changes

**FrontEnd Mode Removal:**

FrontEnd mode (legacy cluster architecture) has been removed in v0.0.43.

```go
// v0.0.42 - FrontEnd fields existed
if node.FrontEnd != nil {
    // Handle front-end configuration
}

// v0.0.43 - Fields removed
// Modern clusters use standard node architecture
// No migration path needed - architecture change required
```

**Impact:** Only affects users running legacy FrontEnd clusters. Modern clusters are unaffected.

### Migration Steps

**Step 1: Update version**
```go
// Before
client, err := slurm.NewClientWithVersion(ctx, "v0.0.42", opts...)

// After
client, err := slurm.NewClientWithVersion(ctx, "v0.0.43", opts...)
```

**Step 2: Adopt new features (optional)**
```go
// Enable QoS management
if client.QoS() != nil {
    qosList, err := client.QoS().List(ctx, nil)
    // Use QoS features
}

// Enable reservation management
if client.Reservations() != nil {
    reservations, err := client.Reservations().List(ctx, nil)
    // Use reservation features
}
```

**Step 3: Remove FrontEnd references (if any)**
```bash
# Find FrontEnd references
grep -rn "FrontEnd" your-project/

# Remove or update to standard node architecture
```

**Step 4: Test new features**
```go
func TestV0043Features(t *testing.T) {
    client, err := slurm.NewClientWithVersion(ctx, "v0.0.43", opts...)
    require.NoError(t, err)

    // Test QoS
    if client.QoS() != nil {
        qosList, err := client.QoS().List(ctx, nil)
        require.NoError(t, err)
    }

    // Test Reservations
    if client.Reservations() != nil {
        reservations, err := client.Reservations().List(ctx, nil)
        require.NoError(t, err)
    }
}
```

## v0.0.43 → v0.0.44

### Overview

v0.0.44 is a minor update with incremental improvements and no breaking changes.

**Changes:**
- ✅ Additional fields in job responses
- ✅ Performance optimizations
- ✅ Bug fixes and refinements

### Migration Steps

**Step 1: Update version**
```go
// Before
client, err := slurm.NewClientWithVersion(ctx, "v0.0.43", opts...)

// After
client, err := slurm.NewClientWithVersion(ctx, "v0.0.44", opts...)
```

**Step 2: Verify compatibility**
```bash
# Run existing tests against v0.0.44
SLURM_API_VERSION=v0.0.44 go test ./...

# No code changes required
```

**Step 3: Update builders (if using)**
```bash
# Regenerate builders for v0.0.44
make generate-mocks VERSION=v0.0.44

# Update imports in tests
sed -i 's/builderv0043/builderv0044/g' tests/**/*_test.go
```

## Testing Your Migration

### Pre-Migration Checklist

- [ ] Review breaking changes for your target version
- [ ] Backup your current code
- [ ] Create a feature branch for migration
- [ ] Document version-specific dependencies

### Migration Testing Strategy

**Phase 1: Unit Tests**
```bash
# Run tests with old version
SLURM_API_VERSION=v0.0.42 go test ./...

# Update code

# Run tests with new version
SLURM_API_VERSION=v0.0.43 go test ./...
```

**Phase 2: Integration Tests**
```bash
# Test against mock server
go test ./tests/integration/... -v

# Test against real SLURM cluster (if available)
SLURM_REAL_SERVER_TEST=true \
SLURM_SERVER_URL=https://cluster:6820 \
SLURM_JWT_TOKEN=<token> \
go test ./tests/integration/... -v
```

**Phase 3: Smoke Tests**
```go
func TestMigrationSmokeTest(t *testing.T) {
    client, err := slurm.NewClientWithVersion(ctx, "v0.0.43", opts...)
    require.NoError(t, err)

    // Test core operations
    jobs, err := client.Jobs().List(ctx, nil)
    require.NoError(t, err)

    nodes, err := client.Nodes().List(ctx, nil)
    require.NoError(t, err)

    partitions, err := client.Partitions().List(ctx, nil)
    require.NoError(t, err)

    // Test new features (if applicable)
    if client.Reservations() != nil {
        _, err := client.Reservations().List(ctx, nil)
        require.NoError(t, err)
    }
}
```

### Post-Migration Validation

- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] Smoke tests validate core functionality
- [ ] New features work as expected
- [ ] Performance benchmarks show no regression
- [ ] Documentation updated

## Common Pitfalls

### Pitfall 1: Mixing Version Types

**Problem:**
```go
// ❌ Wrong: Mixing v0.0.40 and v0.0.42 types
func badMigration() {
    v40job := &v0_0_40.V0040JobInfo{}
    v42cpus := &v0_0_42.V0042Uint32NoValStruct{} // Type mismatch!
    v40job.Cpus = v42cpus // Compilation error
}
```

**Solution:**
```go
// ✅ Correct: Use consistent version types
func goodMigration() {
    v42job := &v0_0_42.V0042JobInfo{}
    v42cpus := &v0_0_42.V0042Uint32NoValStruct{}
    v42job.Cpus = v42cpus // Type safe
}
```

### Pitfall 2: Assuming Feature Availability

**Problem:**
```go
// ❌ Wrong: Assuming QoS exists in all versions
qosList, err := client.QoS().List(ctx, nil) // Panic on v0.0.42!
```

**Solution:**
```go
// ✅ Correct: Check feature availability
if client.QoS() != nil {
    qosList, err := client.QoS().List(ctx, nil)
    // Safe to use
} else {
    log.Println("QoS not supported in this API version")
}
```

### Pitfall 3: Hardcoding Version-Specific Logic

**Problem:**
```go
// ❌ Wrong: Hardcoding version checks
if client.Version() == "v0.0.43" {
    // Special handling
}
```

**Solution:**
```go
// ✅ Correct: Use capability detection
if client.Reservations() != nil {
    // Feature is available, use it
}
```

### Pitfall 4: Ignoring NoVal Changes

**Problem:**
```go
// ❌ Wrong: Directly accessing internal types
job.Cpus = &v0_0_42.V0042Uint32NoValStruct{
    Set:    &setTrue,
    Number: &cpuCount,
    // Missing: Infinite field
}
```

**Solution:**
```go
// ✅ Correct: Use builders
job := builderv0042.NewJobInfo().
    WithCpus(cpuCount).
    Build()
// Builder handles all internal fields correctly
```

### Pitfall 5: Not Testing Against Real Server

**Problem:**
Testing only against mock servers may miss version-specific behavior differences.

**Solution:**
```bash
# Test against real SLURM cluster before production deployment
SLURM_REAL_SERVER_TEST=true \
SLURM_SERVER_URL=https://test-cluster:6820 \
SLURM_JWT_TOKEN=<token> \
go test ./tests/integration/... -v
```

## Version-Specific Gotchas

### v0.0.40 → v0.0.42

- **NoVal types** changed internally - use builders to avoid manual field management
- **Field names** changed for switches - use find & replace
- **Exclusive/Oversubscribe** removed from jobs - check partition settings instead

### v0.0.42 → v0.0.43

- **New managers** (QoS, Reservations, Accounts) - check for nil before using
- **FrontEnd mode** removed - only affects legacy cluster architectures
- **Enum arrays** use different patterns - builders handle automatically

### v0.0.43 → v0.0.44

- **Minor updates only** - no significant breaking changes
- **Additional fields** in responses - test for nil when accessing new fields

## Getting Help

If you encounter issues during migration:

1. **Check version compatibility**: See [README.md](../index.md#supported-versions)
2. **Review examples**: See [examples/](examples/) for version-specific examples
3. **Run diagnostics**: Use `./scripts/diagnose-slurm-auth.sh` for auth issues
4. **Open an issue**: [GitHub Issues](https://github.com/jontk/slurm-client/issues)

## Additional Resources

- [API Compatibility Matrix](../index.md#api-compatibility-matrix)
- [Real Server Testing Guide](https://github.com/jontk/slurm-client/blob/main/tests/integration/REAL_SERVER_TESTING.md)
- [Builder Pattern Guide](https://github.com/jontk/slurm-client/blob/main/tests/mocks/generated/README.md)
- [Architecture Documentation](../guides/ARCHITECTURE.md)
