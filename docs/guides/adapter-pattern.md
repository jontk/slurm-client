# Adapter Pattern Deep Dive

This guide provides an in-depth look at how the adapter pattern works in the slurm-client library, enabling seamless multi-version API support.

## Table of Contents

- [Overview](#overview)
- [Why Adapters?](#why-adapters)
- [Architecture](#architecture)
- [Type Conversion](#type-conversion)
- [Performance Characteristics](#performance-characteristics)
- [Extension Patterns](#extension-patterns)
- [Internal Implementation](#internal-implementation)

## Overview

The adapter pattern is the core architectural pattern that enables this library to support multiple SLURM REST API versions (v0.0.40 through v0.0.44) with a single, unified interface.

### Key Benefits

- **Version Abstraction**: Write code once, works across all API versions
- **Type Safety**: Compile-time guarantees with proper error handling
- **Automatic Conversion**: Seamless translation between version-specific and public types
- **Performance**: Zero-copy operations where possible, intelligent caching
- **Future-Proof**: New API versions added without breaking changes

## Why Adapters?

### The Problem

SLURM REST API versions differ in:

1. **Field Names**: `minimum_switches` → `required_switches` (v0.0.40→v0.0.41)
2. **Data Types**: Some fields change from strings to enums
3. **Structure**: Response formats evolve between versions
4. **Availability**: Features like Reservations only in v0.0.43+
5. **Field Removal**: Some fields deprecated and removed

### The Solution

Adapters provide a **translation layer** that:

```
User Code (public API)
        ↓
    Adapter Layer (converts public ↔ internal types)
        ↓
Version-Specific Implementation (v0.0.40, v0.0.41, etc.)
        ↓
    OpenAPI-Generated Client
        ↓
   SLURM REST API
```

## Architecture

### Component Structure

```
pkg/slurm/
├── client.go              # Public Client interface
├── types.go               # Public types (Job, Node, etc.)
└── options.go             # Configuration options

internal/
├── adapters/
│   ├── common/
│   │   ├── interfaces.go  # Common adapter interfaces
│   │   └── base.go        # Shared adapter functionality
│   ├── v0_0_40/
│   │   ├── adapter_client.go
│   │   ├── job_adapter.go
│   │   ├── node_adapter.go
│   │   └── converters.go  # Type conversion functions
│   ├── v0_0_43/
│   │   └── ... (same structure)
│   └── v0_0_44/
│       └── ... (same structure)
└── api/
    ├── v0_0_40/           # OpenAPI-generated clients
    ├── v0_0_43/
    └── v0_0_44/
```

### Adapter Interfaces

Each adapter implements common interfaces:

```go
// JobAdapter converts between API types and public types
type JobAdapter interface {
    List(ctx context.Context, opts *ListJobsOptions) (*JobList, error)
    Get(ctx context.Context, jobID string) (*Job, error)
    Submit(ctx context.Context, job *JobSubmission) (*JobSubmitResponse, error)
    Cancel(ctx context.Context, jobID string) error
}

// NodeAdapter handles node operations
type NodeAdapter interface {
    List(ctx context.Context, opts *ListNodesOptions) (*NodeList, error)
    Get(ctx context.Context, nodeName string) (*Node, error)
}
```

### Client Initialization Flow

```go
// User creates client
client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL("https://cluster:6820"),
    slurm.WithAuth(auth.NewTokenAuth("token")),
)

// Internally:
// 1. Auto-detect API version by querying /openapi/v3
// 2. Select appropriate adapter (e.g., v0_0_43 adapter)
// 3. Initialize version-specific OpenAPI client
// 4. Wrap in adapter that implements public interfaces
// 5. Return unified Client interface to user
```

## Type Conversion

### Conversion Philosophy

Adapters perform **bidirectional type conversion**:

1. **Request Path**: Public types → Internal API types
2. **Response Path**: Internal API types → Public types

### Example: Job Submission

**User Code:**
```go
submission := &slurm.JobSubmission{
    Name:      "test-job",
    Script:    "#!/bin/bash\necho hello",
    Partition: "compute",
    CPUs:      4,
    Memory:    4 * 1024 * 1024 * 1024, // 4GB
}

response, err := client.Jobs().Submit(ctx, submission)
```

**Adapter Conversion (v0.0.43):**
```go
func (a *JobAdapter) Submit(ctx context.Context, job *slurm.JobSubmission) (*slurm.JobSubmitResponse, error) {
    // Convert public type to API-specific type
    apiJob := &apiv0043.V0043JobSubmitReq{
        Script: &job.Script,
        Job: apiv0043.V0043JobDescMsg{
            Name:          job.Name,
            Partition:     job.Partition,
            MinCpus:       int32(job.CPUs),
            MinMemoryNode: int64(job.Memory),
            // ... more fields
        },
    }

    // Call OpenAPI client
    resp, err := a.apiClient.SlurmV0043PostJob(ctx).V0043JobSubmitReq(*apiJob).Execute()
    if err != nil {
        return nil, a.convertError(err)
    }

    // Convert response back to public type
    return &slurm.JobSubmitResponse{
        JobID: fmt.Sprintf("%d", resp.JobId),
        // ... more fields
    }, nil
}
```

### Complex Conversions

#### 1. String Lists ↔ Arrays

SLURM often uses comma-separated strings for lists:

```go
// API: "node001,node002,node003" (string)
// Public: []string{"node001", "node002", "node003"}

func convertNodeList(nodeStr string) []string {
    if nodeStr == "" {
        return []string{}
    }
    return strings.Split(nodeStr, ",")
}
```

#### 2. Flags and Enums

```go
// API: []string{"MAINT", "IGNORE_JOBS"}
// Public: []ReservationFlag

func convertFlags(apiFlags []string) []ReservationFlag {
    flags := make([]ReservationFlag, len(apiFlags))
    for i, f := range apiFlags {
        flags[i] = ReservationFlag(f)
    }
    return flags
}
```

#### 3. Memory Units

```go
// API may use MB, the public API uses bytes
// Convert to consistent unit (bytes)

func convertMemoryToBytes(memMB int64) int64 {
    return memMB * 1024 * 1024
}

func convertMemoryToMB(memBytes int64) int64 {
    return memBytes / (1024 * 1024)
}
```

#### 4. Time Representations

```go
// API: Unix timestamps (int64)
// Public: time.Time

func convertUnixToTime(unix int64) time.Time {
    if unix == 0 {
        return time.Time{} // Zero value
    }
    return time.Unix(unix, 0)
}
```

#### 5. Nullable Fields

```go
// API: Pointers (*int32, *string)
// Public: Direct values with zero values for unset

func convertOptionalInt(ptr *int32) int {
    if ptr == nil {
        return 0
    }
    return int(*ptr)
}
```

### Version-Specific Conversions

Some conversions differ between API versions:

```go
// v0.0.40 uses "minimum_switches"
// v0.0.41+ uses "required_switches"

// v0.0.40 adapter:
func (a *v0040Adapter) convertJobDescriptor(job *Job) *apiv0040.JobDescMsg {
    return &apiv0040.JobDescMsg{
        MinimumSwitches: int32(job.RequiredSwitches), // Old field name
    }
}

// v0.0.41 adapter:
func (a *v0041Adapter) convertJobDescriptor(job *Job) *apiv0041.JobDescMsg {
    return &apiv0041.JobDescMsg{
        RequiredSwitches: int32(job.RequiredSwitches), // New field name
    }
}
```

## Performance Characteristics

### Zero-Copy Operations

Where possible, adapters avoid copying data:

```go
// Good: Direct reference
func (a *Adapter) convertJob(apiJob *api.Job) *slurm.Job {
    return &slurm.Job{
        Name:   apiJob.Name,   // String is immutable, no copy
        Script: apiJob.Script, // Direct reference
    }
}

// Avoid: Unnecessary copies
// Instead of creating intermediate slices, convert in-place
```

### Caching Strategies

Adapters cache expensive operations:

```go
type AdapterClient struct {
    apiClient   *openapi.APIClient
    versionInfo *VersionInfo // Cached version info
    infoCache   *InfoCache    // Cache cluster info responses
}

func (c *AdapterClient) Info() InfoManager {
    // Version info is cached after first fetch
    if c.versionInfo == nil {
        c.versionInfo = c.fetchVersionInfo()
    }
    return c.infoManager
}
```

### Batch Operations

For bulk operations, adapters optimize by:

1. **Single API call** instead of multiple
2. **Concurrent conversions** for large result sets
3. **Lazy conversion** of fields accessed rarely

```go
func (a *JobAdapter) List(ctx context.Context, opts *ListJobsOptions) (*JobList, error) {
    // Single API call for all jobs
    resp, err := a.apiClient.GetJobs(ctx, convertOptions(opts))

    // Convert results (could be parallelized for large sets)
    jobs := make([]*Job, len(resp.Jobs))
    for i, apiJob := range resp.Jobs {
        jobs[i] = a.convertJob(apiJob)
    }

    return &JobList{Jobs: jobs}, nil
}
```

### Memory Footprint

Typical adapter overhead:

- **Per-client**: ~10-20KB (adapter struct + caching)
- **Per-operation**: ~1-5KB (temporary conversion buffers)
- **Zero persistent overhead** for simple operations

**Benchmark Results** (v0.0.43 adapter):

```
BenchmarkJobList/100_jobs-8    5000  250000 ns/op   80KB alloc   500 allocs/op
BenchmarkJobList/1000_jobs-8   500   2500000 ns/op  800KB alloc  5000 allocs/op

BenchmarkNodeList/50_nodes-8   10000 120000 ns/op   40KB alloc   250 allocs/op
```

## Extension Patterns

### Adding Custom Adapters

You can create custom adapters for specialized use cases:

```go
// Custom adapter that adds retry logic
type RetryAdapter struct {
    base    JobAdapter
    retries int
}

func (r *RetryAdapter) Submit(ctx context.Context, job *JobSubmission) (*JobSubmitResponse, error) {
    var lastErr error
    for i := 0; i < r.retries; i++ {
        resp, err := r.base.Submit(ctx, job)
        if err == nil {
            return resp, nil
        }
        lastErr = err
        time.Sleep(time.Second * time.Duration(i+1))
    }
    return nil, lastErr
}
```

### Custom Type Conversions

Override conversions for specific needs:

```go
// Custom converter that handles special node name formats
type CustomNodeAdapter struct {
    *v0043.NodeAdapter
}

func (c *CustomNodeAdapter) convertNodeName(apiName string) string {
    // Custom logic: extract node ID from FQDN
    // "node001.cluster.example.com" → "node001"
    parts := strings.Split(apiName, ".")
    return parts[0]
}

func (c *CustomNodeAdapter) Get(ctx context.Context, name string) (*Node, error) {
    // Convert using custom logic before calling base
    apiName := c.expandNodeName(name)
    return c.NodeAdapter.Get(ctx, apiName)
}
```

### Middleware Pattern

Wrap adapters with middleware for cross-cutting concerns:

```go
// Logging middleware
type LoggingAdapter struct {
    base   JobAdapter
    logger *log.Logger
}

func (l *LoggingAdapter) Submit(ctx context.Context, job *JobSubmission) (*JobSubmitResponse, error) {
    l.logger.Printf("Submitting job: %s", job.Name)
    start := time.Now()

    resp, err := l.base.Submit(ctx, job)

    elapsed := time.Since(start)
    if err != nil {
        l.logger.Printf("Job submission failed after %v: %v", elapsed, err)
    } else {
        l.logger.Printf("Job submitted successfully (ID: %s) in %v", resp.JobID, elapsed)
    }

    return resp, err
}
```

## Internal Implementation

### Adapter Lifecycle

```go
// 1. Client initialization
client, _ := slurm.NewClient(ctx, opts...)

// 2. Version detection
version := detectVersion(baseURL) // e.g., "v0.0.43"

// 3. Create version-specific OpenAPI client
apiClient := apiv0043.NewAPIClient(apiConfig)

// 4. Create adapter
adapter := &v0043.AdapterClient{
    apiClient: apiClient,
    jobAdapter: &v0043.JobAdapter{client: apiClient},
    nodeAdapter: &v0043.NodeAdapter{client: apiClient},
    // ... more adapters
}

// 5. Return wrapped client
return &client{adapter: adapter}
```

### Error Handling

Adapters provide rich error context:

```go
func (a *JobAdapter) convertError(err error) error {
    if apiErr, ok := err.(*openapi.GenericOpenAPIError); ok {
        return &errors.SlurmError{
            Code:       errors.ErrorCodeFromHTTP(apiErr.StatusCode),
            Message:    apiErr.Error(),
            APIVersion: "v0.0.43",
            Details:    string(apiErr.Body()),
        }
    }
    return err
}
```

### Version Capability Detection

Adapters advertise their capabilities:

```go
type Capabilities struct {
    Reservations bool // v0.0.43+
    QoS          bool // v0.0.43+
    Accounts     bool // v0.0.43+
}

func (a *AdapterClient) Capabilities() Capabilities {
    return Capabilities{
        Reservations: true,  // v0.0.43 supports this
        QoS:          true,
        Accounts:     true,
    }
}

// Client uses capabilities to provide features
func (c *Client) Reservations() ReservationManager {
    if !c.adapter.Capabilities().Reservations {
        return nil // Feature not available
    }
    return c.adapter.Reservations()
}
```

## Best Practices

### For Library Users

1. **Trust the adapter** - Let it handle version differences
2. **Use feature detection** - Check if optional features are available
3. **Handle errors properly** - Check for version-specific errors
4. **Avoid internal types** - Only use public API types

### For Contributors

1. **Maintain consistency** - Follow existing conversion patterns
2. **Document differences** - Comment on version-specific behavior
3. **Test thoroughly** - Ensure conversions work both ways
4. **Optimize carefully** - Profile before optimizing conversions
5. **Version all changes** - Each adapter is independent

## See Also

- [Architecture Overview](ARCHITECTURE.md) - High-level system design
- [Code Generation](CODE_GENERATION.md) - How OpenAPI clients are generated
- [Contributing Guide](../development/CONTRIBUTING.md) - Adding new adapters
- [API Documentation](../api/README.md) - Public API reference
