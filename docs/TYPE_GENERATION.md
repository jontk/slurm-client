# Type Generation from OpenAPI Specification

## Overview

This document explains the code generation approach for creating clean, maintainable types from the SLURM REST API OpenAPI specification.

## Decision: Follow OpenAPI Structure (Option B)

We've decided to **stop flattening** nested OpenAPI structures and instead:
1. ✅ Follow OpenAPI structure exactly
2. ✅ Generate types automatically from spec
3. ✅ Simplify/unwrap special OpenAPI types
4. ✅ Verify coverage automatically

**Why?**
- Industry standard (Kubernetes, AWS, Google Cloud all do this)
- Maintainable (regenerate when OpenAPI updates)
- Verifiable (can achieve 95-100% coverage)
- Simpler adapters (no complex flatten/unflatten logic)

## Tools Created

### 1. Generated Types Status

**Current Status:**
- All main entity types (Job, Node, Partition, etc.) generated and verified
- Nested types properly named with parent context to avoid collisions
- All types compile successfully
- Used by adapters across all supported versions (v0.0.40 - v0.0.44)

### 2. Type Generator (Python)
**Location:** `tools/codegen/generate_clean_types.py`

**Purpose:** Generates clean Go types from OpenAPI spec

**Usage:**
```bash
# Generate using manual entity list
python3 tools/codegen/generate_clean_types.py \
  openapi-specs/slurm-v0.0.44.json \
  internal/common/types/

# Generate with auto-discovery (recommended for new versions)
python3 tools/codegen/generate_clean_types.py --discover \
  openapi-specs/slurm-v0.0.44.json \
  internal/common/types/

# Dry-run to see what would be generated
python3 tools/codegen/generate_clean_types.py --dry-run --discover \
  openapi-specs/slurm-v0.0.44.json \
  internal/common/types/

# Specify version explicitly (if auto-detection fails)
python3 tools/codegen/generate_clean_types.py --version 0.0.44 \
  openapi-specs/slurm-v0.0.44.json \
  internal/common/types/
```

**Features:**
- ✅ Reads any OpenAPI spec version (v0.0.40 - v0.0.44+)
- ✅ Generates clean type names (`Node` not `V0044Node`)
- ✅ Preserves OpenAPI nesting with parent context in nested type names
- ✅ Unwraps special types (`V0044Uint32NoValStruct` → `uint32`)
- ✅ Handles timestamps (`uint64` → `time.Time`)
- ✅ Unwraps CSV strings (`V0044CsvString` → `[]string`)
- ✅ Unwraps list types (`V0044TresList` → `[]TRES`)
- ✅ Auto-discovers auxiliary types from base entities (`--discover`)
- ✅ Handles deprecated empty types (converts to `interface{}`)
- ✅ Uses smart pointer logic (required fields non-pointer, arrays never pointer)
- ✅ Generates proper nested type names to avoid collisions (e.g., `QoSLimitsMaxAccruing`)

### 3. Supported API Versions

The generator supports multiple SLURM API versions:

| Version | SLURM Version | Types Generated | Status |
|---------|---------------|-----------------|--------|
| v0.0.44 | 24.11.x | 23 types | Active |
| v0.0.43 | 24.05.x | ~21 types | Active |
| v0.0.42 | 23.11.x | ~20 types | Maintenance |
| v0.0.41 | 23.02.x | ~20 types | Deprecated |
| v0.0.40 | 22.05.x | 18 types | Deprecated |

## Generated Type Example

**Before (Flattened, Manual):**
```go
type Partition struct {
    Name          string   `json:"name"`
    AllowAccounts []string `json:"allow_accounts,omitempty"` // Flattened from allowed.accounts
    AllowGroups   []string `json:"allow_groups,omitempty"`   // Flattened from allowed.groups
    DenyAccounts  []string `json:"deny_accounts,omitempty"`  // Flattened from deny.accounts
}
```

**After (Following OpenAPI, Generated):**
```go
// Code generated from OpenAPI spec. DO NOT EDIT.
type Partition struct {
    Name    string            `json:"name"`
    Allowed *PartitionAllowed `json:"allowed,omitempty"` // Matches OpenAPI structure
    Deny    *PartitionDeny    `json:"deny,omitempty"`
    State   PartitionState    `json:"state"`
}

type PartitionAllowed struct {
    Accounts []string `json:"accounts,omitempty"`
    Groups   []string `json:"groups,omitempty"`
    QoS      []string `json:"qos,omitempty"`
}
```

## Type Unwrapping Map

OpenAPI uses special wrapper types that should be unwrapped:

| OpenAPI Type | Clean Go Type | Notes |
|--------------|---------------|-------|
| `V0044Uint32NoValStruct` | `uint32` | Simple integer |
| `V0044Uint64NoValStruct` | `uint64` or `time.Time` | Depends on field |
| `V0044Uint16NoValStruct` | `uint16` | Simple integer |
| `V0044Int32NoValStruct` | `int32` | Simple integer |
| `V0044Float64NoValStruct` | `float64` | Simple float |
| `V0044CsvString` | `[]string` | CSV parsed to slice |
| `V0044ProcessExitCodeVerbose` | `int32` | Exit code |

**Timestamp Fields** (unwrap to `time.Time`):
- `boot_time`, `last_busy`, `eligible_time`, `end_time`, `start_time`
- `submit_time`, `deadline`, `preempt_time`, `suspend_time`
- `created_time`, `modified_time`, `reason_changed_at`

## Workflow

### Initial Generation
```bash
# 1. Generate types from latest OpenAPI spec (with auto-discovery)
python3 tools/codegen/generate_clean_types.py --discover \
  openapi-specs/slurm-v0.0.44.json \
  internal/common/types/

# 2. Review generated files
ls -la internal/common/types/

# 3. Build to check for issues
go build ./...

# 4. Run tests
make test

# Target: All generated types compile and tests pass
```

### When OpenAPI Updates
```bash
# 1. Download new spec
cd tools/codegen
./download-specs.sh

# 2. Regenerate types with auto-discovery
python3 tools/codegen/generate_clean_types.py --discover \
  openapi-specs/slurm-v0.0.46.json \
  internal/common/types/

# 3. Update adapters (mostly mechanical changes)
# 4. Run tests
make test
# 5. Commit
```

## Impact on Adapters

**Before (Complex):**
```go
// Adapter had to manually flatten nested structures
func convertPartition(api *V0044PartitionInfo) *types.Partition {
    p := &types.Partition{
        Name: *api.Name,
    }

    // Complex flattening logic
    if api.Allowed != nil && api.Allowed.Accounts != nil {
        p.AllowAccounts = parseCSV(*api.Allowed.Accounts)
    }
    // ... many more lines
}
```

**After (Simple):**
```go
// Adapter just unwraps special types
func convertPartition(api *V0044PartitionInfo) *types.Partition {
    p := &types.Partition{
        Name: deref(api.Name),
    }

    // Nested structures converted directly
    if api.Allowed != nil {
        p.Allowed = convertPartitionAllowed(api.Allowed)
    }
    if api.Deny != nil {
        p.Deny = convertPartitionDeny(api.Deny)
    }

    return p
}
```

## Benefits

### Maintainability
- ✅ Regenerate when OpenAPI updates (minutes, not hours)
- ✅ No manual type maintenance
- ✅ Consistent across all versions

### Correctness
- ✅ 95-100% field coverage (verified automatically)
- ✅ No missing fields
- ✅ No type mismatches

### Developer Experience
- ✅ Clean type names (`Node` not `V0044Node`)
- ✅ No special wrapper types exposed
- ✅ Proper Go types (`time.Time`, `[]string`)
- ✅ IDE autocomplete works perfectly

### Simplicity
- ✅ Adapters are mechanical, not complex
- ✅ Less custom logic = fewer bugs
- ✅ Easy to understand codebase

## Next Steps

1. **Add to Build Process** (Future)
   - Makefile target: `make generate-types`
   - CI check: verify types are up-to-date

2. **Converter Generator** (Sprint 5)
   - Generate type converters from config
   - Reduce adapter boilerplate

3. **Contract Tests** (Sprint 5)
   - Verify type compatibility across versions
   - Automated regression testing

## Files

| File | Purpose |
|------|---------|
| `tools/codegen/generate_clean_types.py` | Type generator (Python) |
| `tools/codegen/generate_types.go` | Type generator (Go, alternative) |
| `docs/TYPE_GENERATION.md` | This document |
| `internal/common/types/*.go` | Generated types used by adapters |

## References

- OpenAPI Spec: `openapi-specs/slurm-v0.0.44.json`
- Generated Client: `internal/api/v0_0_44/client.go`
- Internal Types: `internal/common/types/*.go`
- Industry Examples:
  - [Kubernetes client-go](https://github.com/kubernetes/client-go)
  - [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2)
  - [Google Cloud Go SDK](https://github.com/googleapis/google-cloud-go)
