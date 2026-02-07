# SLURM REST API Limitations

This document describes known limitations of the SLURM REST API and how the slurm-client SDK handles them.

## Client-Side Pagination

### Problem

The SLURM REST API does **not** support server-side pagination parameters like `limit` or `offset` for list operations. This affects all resource types:

- Jobs (`GET /slurm/v0.0.44/jobs`)
- Nodes (`GET /slurm/v0.0.44/nodes`)
- Partitions (`GET /slurm/v0.0.44/partitions`)
- Reservations (`GET /slurm/v0.0.44/reservations`)
- Accounts (`GET /slurmdb/v0.0.44/accounts`)
- Users (`GET /slurmdb/v0.0.44/users`)
- QoS (`GET /slurmdb/v0.0.44/qos`)
- Clusters (`GET /slurmdb/v0.0.44/clusters`)
- Associations (`GET /slurmdb/v0.0.44/associations`)
- WCKeys (`GET /slurmdb/v0.0.44/wckeys`)

The API only supports **filtering** parameters (e.g., `update_time`, `flags`) but returns the entire filtered dataset in a single response.

### Impact

For large SLURM deployments, this can result in:

- **High memory usage**: The full dataset must be loaded into memory
- **Slow response times**: Large JSON payloads take time to serialize/deserialize
- **Network overhead**: Unnecessary data transfer when you only need a subset

**Example Scenarios:**
- A cluster with 100,000 jobs: Requesting `Limit: 10` still fetches all 100K jobs from the API
- A database with 50,000 users: Every list operation retrieves the complete user list
- Multiple partitions with 10,000+ nodes: Node lists are always complete

### SDK Implementation

The `slurm-client` SDK implements **client-side pagination** in `ListOptions` types:

```go
type JobListOptions struct {
    // ... filtering options ...

    // Limit specifies the maximum number of jobs to return.
    // WARNING: Due to SLURM REST API limitations, this is CLIENT-SIDE pagination.
    // The full job list is fetched from the server, then sliced.
    Limit  int

    // Offset specifies the number of jobs to skip before returning results.
    // WARNING: This is CLIENT-SIDE pagination.
    Offset int
}
```

**How it works:**
1. SDK calls the SLURM API with filtering parameters only
2. API returns the complete filtered dataset
3. SDK slices the result using `Limit` and `Offset` in memory
4. SDK returns the paginated subset to the caller

### Best Practices

To minimize performance impact:

#### 1. Use Filtering First, Pagination Second

Always apply filters to reduce the dataset **before** pagination:

```go
// ❌ BAD: Fetches all 100K jobs, then slices to 10
opts := &types.JobListOptions{
    Limit: 10,
}

// ✅ GOOD: Fetches only running jobs in partition "gpu", then slices to 10
opts := &types.JobListOptions{
    States:     []types.JobState{types.JobStateRunning},
    Partitions: []string{"gpu"},
    Limit:      10,
}
```

#### 2. Use Time-Based Filtering

For frequently updated resources (jobs, nodes), use `UpdateTime` to fetch only recent changes:

```go
lastCheck := time.Now().Add(-1 * time.Hour)
opts := &types.JobListOptions{
    UpdateTime: &lastCheck,  // Only jobs updated in the last hour
    Limit:      100,
}
```

#### 3. Cache Results Locally

Since the API always returns the full dataset, consider caching:

```go
// Fetch once, reuse for multiple pagination requests
allJobs, err := client.Jobs().List(ctx, &types.JobListOptions{
    States: []types.JobState{types.JobStateRunning},
})

// Paginate locally without additional API calls
page1 := allJobs.Jobs[0:10]
page2 := allJobs.Jobs[10:20]
```

#### 4. Monitor Resource Usage

For production deployments with large clusters:

- Monitor memory usage when listing resources
- Set reasonable `Limit` values even though it's client-side
- Use filtering to keep dataset sizes manageable
- Consider implementing server-side caching if fetching the same data repeatedly

### Verification

You can verify the lack of server-side pagination by examining the generated API client:

```go
// internal/api/v0_0_44/client.go
type SlurmV0044GetJobsParams struct {
    UpdateTime *string  // Supported
    Flags      *string  // Supported
    // No Limit parameter
    // No Offset parameter
}
```

### Future Considerations

**Potential improvements if SLURM API adds pagination:**

1. **Server-side pagination**: If future SLURM versions add `limit`/`offset` parameters, the SDK can be updated to pass them through
2. **Streaming/chunked processing**: Large datasets could be streamed instead of loaded entirely into memory
3. **Cursor-based pagination**: More efficient than offset-based for frequently changing data

These features would require changes to the SLURM REST API itself and cannot be implemented in the client library alone.

## Other Limitations

### Version-Specific Features

Not all API versions support all operations. Use the `Capabilities()` method to check feature availability:

```go
client := factory.NewClient("v0.0.41", config)
caps := client.Capabilities()

if !caps.SupportsJobSubmit {
    return errors.New("job submission not supported in v0.0.41")
}
```

See [Phase 2 documentation](../README.md) for details on capability discovery.

### Read-Only Operations in Early Versions

- **v0.0.40-v0.0.41**: Limited write operations (mostly read-only)
- **v0.0.42+**: Full create/update/delete support for most resources

### No WebSocket/Streaming Support

The SLURM REST API does not provide WebSocket or Server-Sent Events (SSE) endpoints for real-time updates. Applications must poll the API for changes.

---

**Last Updated**: 2025-02-03
**API Versions Covered**: v0.0.40 - v0.0.44
