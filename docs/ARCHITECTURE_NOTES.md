# Architecture Notes

## Type System Design

### Two-Layer Type System

The SDK uses a two-layer type system:

1. **Public API Types** (`interfaces/` package)
   - Simplified, stable types exposed to SDK users
   - Focused on common use cases
   - Backward compatibility is critical
   - Example: `WCKey` with just Name, User, Cluster fields

2. **Internal Types** (`internal/common/types/` package)
   - Detailed types used by version adapters
   - Include all API fields from SLURM REST API
   - Can evolve with API versions
   - Example: `WCKey` with ID, Description, Active, timestamps, etc.

### Why Not Use Type Aliases?

The original refactoring plan considered using type aliases in `interfaces/types.go` to reference `internal/common/types`. However, analysis revealed:

1. **Different Field Sets**: Public types are intentionally simplified
2. **Stability vs Evolution**: Public API needs stability, internal types track API changes
3. **Backward Compatibility**: Changing public types would break existing SDK users

### Benefits of Current Approach

- Public API remains stable across SLURM API version changes
- Internal adapters can use full API details
- Clear separation of concerns
- Flexibility to evolve internal types independently

### Trade-offs

- Some code duplication between layers
- Need to manually maintain consistency where types overlap
- Conversion logic needed in adapters (already implemented)

## Interface Organization

After Sprint 4 refactoring, interfaces are split into focused files:

- `client.go` - Core SlurmClient interface
- `job.go` - Job operations split into Reader/Writer/Controller/Watcher
- `job_analytics.go` - Analytics operations (20 methods)
- `node.go`, `partition.go`, `reservation.go`, `qos.go` - Resource managers
- `account.go`, `user.go` - Account management
- `association.go` - Associations, WCKeys, Clusters
- `cluster.go` - Cluster manager
- `info.go` - Cluster information
- `types.go` - Type definitions (public API)

This organization follows Interface Segregation Principle - clients only depend on methods they use.
