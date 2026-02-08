# Slurm REST API Go Client Library

[![Release](https://img.shields.io/github/v/release/jontk/slurm-client)](https://github.com/jontk/slurm-client/releases/latest)
[![CI](https://github.com/jontk/slurm-client/actions/workflows/ci.yml/badge.svg)](https://github.com/jontk/slurm-client/actions/workflows/ci.yml)
[![Security](https://github.com/jontk/slurm-client/actions/workflows/security.yml/badge.svg)](https://github.com/jontk/slurm-client/actions/workflows/security.yml)
[![codecov](https://codecov.io/gh/jontk/slurm-client/branch/main/graph/badge.svg)](https://codecov.io/gh/jontk/slurm-client)
[![Go Reference](https://pkg.go.dev/badge/github.com/jontk/slurm-client.svg)](https://pkg.go.dev/github.com/jontk/slurm-client)
[![Go Report Card](https://goreportcard.com/badge/github.com/jontk/slurm-client)](https://goreportcard.com/report/github.com/jontk/slurm-client)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/github/go-mod/go-version/jontk/slurm-client)](https://github.com/jontk/slurm-client/blob/main/go.mod)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/jontk/slurm-client/badge)](https://api.securityscorecards.dev/projects/github.com/jontk/slurm-client)

**Production-ready Go client library for SLURM REST API with enterprise-grade multi-version support.**

The definitive solution addressing Go SLURM client ecosystem fragmentation through comprehensive multi-version architecture, structured error handling, and production-ready reliability patterns.

## ✨ Key Features

- **🔄 Multi-Version Support**: Unified client supporting SLURM REST API v0.0.40-v0.0.44 (88,000+ lines OpenAPI specs)
- **🏢 Enterprise-Grade**: Production patterns inspired by AWS SDK and Kubernetes client-go
- **🛡️ Structured Error Handling**: Comprehensive error classification with version-specific mapping
- **⚡ Performance Optimized**: Connection pooling, retry policies with exponential backoff and jitter
- **🔐 Comprehensive Authentication**: JWT tokens, API keys, basic auth, and certificate support
- **📊 Complete API Coverage**: All manager operations (Jobs, Nodes, Partitions, Info) implemented
- **🎯 Type-Safe**: Full Go struct definitions generated from authentic SLURM OpenAPI specifications
- **📖 Context-Aware**: All operations support context for cancellation and timeouts

## 🚀 Multi-Version Architecture

Unlike existing solutions that support single API versions, this library provides seamless compatibility across all active SLURM REST API versions:

### Adapter Pattern Implementation (Recommended)

The library implements a sophisticated **adapter pattern** that provides version abstraction while maintaining optimal performance and type safety. Most users should use this approach for production applications.

#### Architecture Overview

```
Client Application
        ↓
    Public Interfaces (interfaces/)
        ↓
┌─────────────┬─────────────┐
│ Adapter     │ Wrapper     │
│ Pattern     │ Pattern     │
│ (Recommended)│ (Direct)   │
└─────────────┴─────────────┘
        ↓           ↓
Version-Specific Implementations
(v0.0.40, v0.0.41, v0.0.42, v0.0.43, v0.0.44)
```

#### Adapter Pattern Benefits

- **🎯 Version Abstraction**: Single interface works across all API versions
- **🔧 Automatic Conversion**: Seamless type conversion between internal types and public interfaces
- **🛡️ Type Safety**: Compile-time guarantees with comprehensive error handling
- **⚡ Performance**: Zero-copy operations where possible, intelligent caching
- **🔄 Future-Proof**: Easy addition of new API versions without breaking changes

#### Adapter vs Wrapper Comparison

| Feature | Adapter Pattern | Wrapper Pattern |
|---------|----------------|-----------------|
| **Abstraction Level** | High - Version agnostic | Low - Version specific |
| **Type Conversion** | Automatic | Manual |
| **Performance** | Optimized with caching | Direct API calls |
| **Complexity** | Simple to use | Requires version knowledge |
| **Recommended For** | Production applications | Advanced users, debugging |

#### Implementation Example

```go
// Using Adapter Pattern (Recommended)
func main() {
    // Automatically selects best compatible version
    client, err := slurm.NewClient(ctx,
        slurm.WithBaseURL("https://cluster:6820"),
        slurm.WithAuth(auth.NewTokenAuth("token")),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Works across all API versions - adapter handles differences
    reservations, err := client.Reservations().List(ctx, nil)
    if err != nil {
        // Structured error handling with version context
        if errors.IsVersionNotSupported(err) {
            log.Printf("Reservations not available in API %s", client.Version())
            return
        }
        log.Fatal(err)
    }

    // Type-safe access to unified interface
    for _, res := range reservations.Reservations {
        fmt.Printf("Reservation: %s, Nodes: %v\n", res.Name, res.Nodes)
    }
}
```

#### Advanced Adapter Features

**Type Conversion Engine**
```go
// Adapter automatically handles complex type conversions
type ReservationAdapter struct {
    // Converts between internal types.Reservation and interfaces.Reservation
    convertReservationToInterface(types.Reservation) interfaces.Reservation
    convertFlags([]types.ReservationFlag) []string
    convertLicenses(map[string]int32) map[string]int
    convertNodeList(string) []string
}

// Example: Complex conversion with error handling
reservation, err := client.Reservations().Get(ctx, "maintenance")
if err != nil {
    if slurmErr, ok := err.(*errors.SlurmError); ok {
        fmt.Printf("API Version: %s, Error: %s\n", slurmErr.APIVersion, slurmErr.Message)
    }
}
```

**Version-Aware Error Handling**
```go
// Adapters provide version context in all errors
_, err := client.QoS().Create(ctx, qosSpec)
if err != nil {
    if errors.IsVersionNotSupported(err) {
        log.Printf("QoS management requires API v0.0.43+, current: %s", client.Version())
        // Graceful degradation or alternative approach
    }
}
```

**Performance Optimizations**
```go
// Adapter pattern includes built-in optimizations:
// - Zero-copy type conversion where possible
// - Intelligent field mapping to avoid unnecessary allocations
// - Response caching for expensive operations like cluster info
// - Connection pooling managed at the adapter level
```
    fmt.Printf("Reservation: %s\n", *res.Name)
}

// After: Adapter pattern (works across versions)
client, err := slurm.NewClient(ctx, opts...)
reservations, err := client.Reservations().List(ctx, nil)
// Automatic type conversion and version abstraction
for _, res := range reservations.Reservations {
    // Clean interface types with consistent behavior
    fmt.Printf("Reservation: %s\n", res.Name)
}
```

| SLURM Version | Supported API Versions | Recommended | Status |
|---------------|------------------------|-------------|---------|
| 24.05-25.05 | v0.0.40, v0.0.41, v0.0.42 | **v0.0.42** | ✅ Supported |
| 24.11-25.11 | v0.0.41, v0.0.42, v0.0.43 | **v0.0.42** | ✅ Supported |
| 25.05+ | v0.0.42, v0.0.43 | **v0.0.42** | ✅ Supported |

### Automatic Version Detection
```go
// Automatically detects and uses the best compatible API version
client, err := slurm.NewClient(ctx, slurm.WithBaseURL("https://your-cluster:6820"))

// Or specify a version explicitly
client, err := slurm.NewClientWithVersion(ctx, "v0.0.42", options...)

// Or target a specific SLURM version
client, err := slurm.NewClientForSlurmVersion(ctx, "25.05", options...)
```

## 📋 Supported Versions

### API Version Support

| API Version | Builder Support | Mock Server | Status | Notes |
|-------------|----------------|-------------|---------|-------|
| v0.0.44 | ✅ Complete (129 methods) | ✅ Native | **Supported** | Latest version |
| v0.0.43 | ✅ Complete (128 methods) | ✅ Native | **Supported** | QoS, Reservations, Accounts |
| v0.0.42 | ✅ Complete (123 methods) | ✅ Native | **Supported** | Stable baseline |
| v0.0.41 | ❌ No builders | ⚠️ Fallback | **Best-effort** | Inline schemas unsupported |
| v0.0.40 | ✅ Complete (126 methods) | ⚠️ Fallback | **Best-effort** | Legacy support |

### Support Tiers

**Supported** (v0.0.42, v0.0.43, v0.0.44):
- ✅ Full test coverage and validation
- ✅ Complete builder pattern support
- ✅ Native mock server implementations
- ✅ Regular security updates
- ✅ Bug fixes and improvements
- ✅ Documentation and examples

**Best-effort** (v0.0.40, v0.0.41):
- ⚠️ Limited test coverage
- ⚠️ Mock server uses fallback implementations
- ⚠️ May not receive new features
- ⚠️ Security fixes only for critical issues
- ⚠️ Use for legacy compatibility only

### Deprecation Policy

When new SLURM REST API versions are released, the library follows this deprecation process:

1. **New Version Released**: Added with full support within 30 days
2. **Support Window**: Latest 3 versions receive full support
3. **Deprecation Notice**: Oldest supported version moves to best-effort
4. **Removal Timeline**: Best-effort versions deprecated after 12 months of inactivity

**Example Timeline:**
- v0.0.45 released → v0.0.45, v0.0.44, v0.0.43 are **Supported**
- v0.0.42 moves to **Best-effort**
- v0.0.41, v0.0.40 marked **Deprecated**

**Current Status (as of 2026-01-19):**
- **Supported:** v0.0.44, v0.0.43, v0.0.42 (full support, regular updates)
- **Best-effort:** v0.0.41, v0.0.40 (minimal maintenance, use at your own risk)

### Migration Assistance

For migration between versions, see [MIGRATION.md](development/MIGRATION.md) which includes:
- Breaking changes between versions
- API compatibility matrices
- Step-by-step migration guides
- Common pitfalls and solutions

## 📦 Installation

### For Library Users

Install the latest version:

```bash
go get github.com/jontk/slurm-client@latest
```

Or install a specific version:

```bash
# Install a specific release
go get github.com/jontk/slurm-client@v0.3.0

# Or use a version constraint in go.mod
require github.com/jontk/slurm-client v0.3.0
```

That's it! The library is ready to use in your applications.

**Latest Release**: See [releases page](https://github.com/jontk/slurm-client/releases) for version details and changelog.

### For Contributors/Developers

If you're cloning the repository to contribute or run tests:

```bash
# Clone the repository
git clone https://github.com/jontk/slurm-client.git
cd slurm-client

# Install dependencies
go mod download

# Generate mock builders (required for tests)
make generate-mocks

# Run tests
make test
```

**Note**: The mock builders are generated from OpenAPI specs and are not committed to the repository. You must run `make generate-mocks` before running tests locally. CI automatically generates them.

## ⚡ Quick Start (Recommended Approach)

> **💡 New to slurm-client?** Use the **Adapter Pattern** for the best experience. It provides version-agnostic APIs with automatic type conversion and error handling.

### Basic Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/jontk/slurm-client"
    "github.com/jontk/slurm-client/pkg/auth"
    "github.com/jontk/slurm-client/pkg/errors"
)

func main() {
    // 🎯 Adapter Pattern: Works across all SLURM versions
    client, err := slurm.NewClient(context.Background(),
        slurm.WithBaseURL("https://your-slurm-server:6820"),
        slurm.WithAuth(auth.NewTokenAuth("your-token")),
        // Version auto-detected - no need to specify!
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // List running jobs with structured error handling
    jobs, err := client.Jobs().List(context.Background(), &interfaces.ListJobsOptions{
        States: []string{"RUNNING"},
        Limit:  10,
    })
    if err != nil {
        // Structured error handling with specific error types
        if slurmErr, ok := err.(*errors.SlurmError); ok {
            switch slurmErr.Code {
            case errors.ErrorCodeUnauthorized:
                log.Fatal("Authentication failed - check your token")
            case errors.ErrorCodeConnectionRefused:
                log.Fatal("Cannot connect to SLURM server")
            default:
                log.Fatalf("SLURM error [%s]: %s", slurmErr.Code, slurmErr.Message)
            }
        }
        log.Fatal(err)
    }

    fmt.Printf("Found %d running jobs\\n", len(jobs.Jobs))

    // Submit a job with comprehensive error handling
    submission := &interfaces.JobSubmission{
        Name:      "test-job",
        Script:    "#!/bin/bash\\necho 'Hello, SLURM!'",
        Partition: "compute",
        CPUs:      2,
        Memory:    4 * 1024 * 1024 * 1024, // 4GB in bytes
        TimeLimit: 60, // 60 minutes
    }

    response, err := client.Jobs().Submit(context.Background(), submission)
    if err != nil {
        if slurmErr, ok := err.(*errors.SlurmError); ok {
            switch slurmErr.Code {
            case errors.ErrorCodeResourceExhausted:
                fmt.Printf("Insufficient resources: %s\\n", slurmErr.Details)
            case errors.ErrorCodeValidationFailed:
                fmt.Printf("Job validation failed: %s\\n", slurmErr.Details)
            default:
                log.Fatalf("Job submission failed [%s]: %s", slurmErr.Code, slurmErr.Message)
            }
        }
        return
    }

    fmt.Printf("Job submitted with ID: %s\\n", response.JobID)

    // Get job details
    job, err := client.Jobs().Get(context.Background(), response.JobID)
    if err != nil {
        log.Printf("Failed to get job details: %v", err)
        return
    }

    fmt.Printf("Job Status: %s, User: %s, Partition: %s\\n",
               job.State, job.UserID, job.Partition)
}
```

### 🤔 Which Pattern Should I Use?

```
Are you building a production application? ──── YES ──→ Use Adapter Pattern ✅
                   │
                   NO
                   │
                   ↓
Do you need direct API access or debugging? ── YES ──→ Use Wrapper Pattern ⚠️
                   │
                   NO
                   │
                   ↓
                Use Adapter Pattern ✅ (Default choice)
```

| **Use Adapter Pattern When** | **Use Wrapper Pattern When** |
|-------------------------------|-------------------------------|
| ✅ **Building production applications** | ⚠️ **You need direct API access** |
| ✅ **You want version-agnostic code** | ⚠️ **Maximum performance is critical** |
| ✅ **You're new to SLURM APIs** | ⚠️ **You're debugging API responses** |
| ✅ **You want simple error handling** | ⚠️ **You need version-specific features** |
| ✅ **You want automatic type conversion** | ⚠️ **You're migrating from direct API calls** |

> **🎯 Recommendation**: **The Adapter Pattern is recommended for all users**. It provides version abstraction with optimal performance.

#### Adapter Configuration Options

The adapter pattern supports advanced configuration for performance tuning and behavior customization:

```go
import "github.com/jontk/slurm-client/pkg/config"

// Configure adapter-specific behavior
adapterConfig := &config.AdapterConfig{
    // Performance tuning
    EnableTypeCache:     true,  // Cache converted types (15-25% faster)
    EnableResponseCache: true,  // Cache expensive operations like cluster info
    CacheTimeout:        5 * time.Minute,

    // Conversion behavior
    StrictTypeConversion: false, // Allow lossy conversions for compatibility
    PreferZeroCopy:       true,  // Optimize for memory efficiency

    // Version handling
    AutoVersionFallback:  true,  // Fall back to compatible versions
    VersionLockTimeout:   30 * time.Second,

    // Error handling
    ProvideVersionContext: true, // Include API version in all errors
    WrapLegacyErrors:     true,  // Convert old error formats
}

client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL("https://cluster:6820"),
    slurm.WithAdapterConfig(adapterConfig),
)
```

**Environment Variables for Adapters**
```bash
# Adapter-specific configuration
export SLURM_ADAPTER_TYPE_CACHE="true"          # Enable type conversion caching
export SLURM_ADAPTER_RESPONSE_CACHE="true"      # Enable response caching
export SLURM_ADAPTER_CACHE_TIMEOUT="5m"         # Cache timeout duration
export SLURM_ADAPTER_STRICT_TYPES="false"       # Allow lossy type conversions
export SLURM_ADAPTER_VERSION_CONTEXT="true"     # Include version in errors
export SLURM_ADAPTER_AUTO_FALLBACK="true"       # Auto fallback to compatible versions
```

**Performance Monitoring**
```go
// Monitor adapter performance
stats := client.AdapterStats()
fmt.Printf("Cache Hit Rate: %.2f%%\n", stats.CacheHitRate*100)
fmt.Printf("Type Conversions: %d\n", stats.TypeConversions)
fmt.Printf("Average Conversion Time: %v\n", stats.AvgConversionTime)
fmt.Printf("Memory Saved: %d bytes\n", stats.ZeroCopyBytes)
```

## Configuration

The client can be configured via environment variables:

```bash
export SLURM_REST_URL="https://your-slurm-server:6820"
export SLURM_TIMEOUT="30s"
export SLURM_MAX_RETRIES="3"
export SLURM_DEBUG="true"
export SLURM_INSECURE_SKIP_VERIFY="false"
```

Or programmatically:

```go
config := &config.Config{
    BaseURL:    "https://your-slurm-server:6820",
    Timeout:    30 * time.Second,
    MaxRetries: 3,
    Debug:      true,
}

client, err := slurm.NewClient(
    slurm.WithConfig(config),
    slurm.WithAuth(auth.NewTokenAuth("your-token")),
)
```

## Authentication

### Token Authentication

```go
client, err := slurm.NewClient(
    slurm.WithAuth(auth.NewTokenAuth("your-slurm-token")),
)
```

### Basic Authentication

```go
client, err := slurm.NewClient(
    slurm.WithAuth(auth.NewBasicAuth("username", "password")),
)
```

### No Authentication

```go
client, err := slurm.NewClient(
    slurm.WithAuth(auth.NewNoAuth()),
)
```

## 🏭 Production-Ready Features

### Comprehensive Manager Operations
✅ **All 12+ methods implemented** across all managers with structured error handling:

- **JobManager**: `List()`, `Get()`, `Submit()`, `Cancel()` - Complete job lifecycle
- **NodeManager**: `List()`, `Get()` - Resource monitoring and management
- **PartitionManager**: `List()`, `Get()` - Partition discovery and management
- **InfoManager**: `Get()`, `Ping()`, `Stats()`, `Version()` - Cluster information

### Enterprise-Grade Quality
- **100% Test Coverage**: Authentication and configuration packages
- **97.9% Test Coverage**: Retry policies with exponential backoff
- **64.1% Test Coverage**: Structured error handling system
- **Production Patterns**: Connection pooling, timeout management, circuit breakers

## 📖 API Operations

### Job Management

```go
// List jobs with advanced filtering
jobs, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{
    UserID:    "username",
    States:    []string{"RUNNING", "PENDING"},
    Partition: "compute",
    Limit:     50,
    Offset:    0,
})

// Get specific job with structured error handling
job, err := client.Jobs().Get(ctx, "12345")
if err != nil {
    if errors.IsNotFound(err) {
        fmt.Println("Job not found")
    }
}

// Submit job with comprehensive validation
response, err := client.Jobs().Submit(ctx, &interfaces.JobSubmission{
    Name:        "my-job",
    Script:      "#!/bin/bash\\necho 'Hello World'",
    Partition:   "compute",
    CPUs:        4,
    Memory:      8 * 1024 * 1024 * 1024, // 8GB in bytes
    TimeLimit:   60, // minutes
    Nodes:       1,
    WorkingDir:  "/tmp",
    Environment: map[string]string{"CUDA_VISIBLE_DEVICES": "0"},
})

// Cancel job with proper error handling
err := client.Jobs().Cancel(ctx, "12345")
```

### Node Management

```go
// List nodes with filtering
nodes, err := client.Nodes().List(ctx, &interfaces.ListNodesOptions{
    States:    []string{"IDLE", "ALLOCATED"},
    Partition: "compute",
    Features:  []string{"gpu", "ssd"},
    Limit:     100,
})

// Get specific node information
node, err := client.Nodes().Get(ctx, "node001")
if err != nil {
    log.Printf("Failed to get node: %v", err)
    return
}

fmt.Printf("Node %s: %d CPUs, %d MB memory, State: %s\\n",
           node.Name, node.CPUs, node.Memory/(1024*1024), node.State)
```

### Partition Management

```go
// List all partitions
partitions, err := client.Partitions().List(ctx, nil)
if err != nil {
    log.Printf("Failed to list partitions: %v", err)
    return
}

// Get specific partition details
partition, err := client.Partitions().Get(ctx, "compute")
if err != nil {
    log.Printf("Failed to get partition: %v", err)
    return
}

fmt.Printf("Partition %s: Max nodes: %d, Default time: %d min\\n",
           partition.Name, partition.MaxNodes, partition.DefaultTime)
```

### Cluster Information

```go
// Test cluster connectivity
err := client.Info().Ping(ctx)
if err != nil {
    log.Printf("Cluster unreachable: %v", err)
}

// Get cluster information and statistics
info, err := client.Info().Get(ctx)
if err == nil {
    fmt.Printf("Cluster: %s, SLURM version: %s\\n", info.ClusterName, info.Version)
}

stats, err := client.Info().Stats(ctx)
if err == nil {
    fmt.Printf("Running jobs: %d, Total nodes: %d\\n",
               stats.RunningJobs, stats.TotalNodes)
}

// Get API version information
version, err := client.Info().Version(ctx)
if err == nil {
    fmt.Printf("API Version: %s\\n", version.Version)
}
```

### Reservation Management (v0.0.43+)

```go
// Check if reservations are supported
if client.Reservations() == nil {
    log.Println("Reservations not supported in this API version")
    return
}

// List all reservations
reservations, err := client.Reservations().List(ctx, nil)
if err != nil {
    log.Printf("Failed to list reservations: %v", err)
    return
}

for _, res := range reservations.Reservations {
    fmt.Printf("Reservation %s: %s to %s\\n",
               res.Name, res.StartTime, res.EndTime)
}

// Create a reservation
newReservation := &interfaces.ReservationCreate{
    Name:      "maintenance-window",
    StartTime: time.Now().Add(24 * time.Hour),
    Duration:  4 * 3600, // 4 hours
    Nodes:     []string{"node001", "node002"},
    Users:     []string{"admin"},
    Flags:     []string{"MAINT", "IGNORE_JOBS"},
}

resp, err := client.Reservations().Create(ctx, newReservation)
if err != nil {
    log.Printf("Failed to create reservation: %v", err)
} else {
    fmt.Printf("Created reservation: %s\\n", resp.ReservationName)
}

// Update a reservation
update := &interfaces.ReservationUpdate{
    EndTime: &newEndTime,
    Users:   []string{"admin", "operator"},
}

err = client.Reservations().Update(ctx, "maintenance-window", update)
if err != nil {
    log.Printf("Failed to update reservation: %v", err)
}

// Delete a reservation
err = client.Reservations().Delete(ctx, "maintenance-window")
if err != nil {
    log.Printf("Failed to delete reservation: %v", err)
}
```

### Quality of Service (QoS) Management (v0.0.43+)

```go
// Check if QoS is supported
if client.QoS() == nil {
    log.Println("QoS not supported in this API version")
    return
}

// List all QoS configurations
qosList, err := client.QoS().List(ctx, nil)
if err != nil {
    log.Printf("Failed to list QoS: %v", err)
    return
}

for _, qos := range qosList.QoS {
    fmt.Printf("QoS %s: Priority=%d, MaxCPUs=%d\\n",
               qos.Name, qos.Priority, qos.MaxCPUs)
}

// Create a new QoS
newQoS := &interfaces.QoSCreate{
    Name:               "high-priority",
    Description:        "High priority for critical jobs",
    Priority:           10000,
    PreemptMode:        "requeue",
    MaxJobs:            50,
    MaxJobsPerUser:     10,
    MaxCPUs:            500,
    MaxWallTime:        86400, // 24 hours
    UsageFactor:        2.0,   // Double charge
    Flags:              []string{"DenyOnLimit", "RequireAssoc"},
    AllowedAccounts:    []string{"research", "production"},
}

resp, err := client.QoS().Create(ctx, newQoS)
if err != nil {
    log.Printf("Failed to create QoS: %v", err)
} else {
    fmt.Printf("Created QoS: %s\\n", resp.QoSName)
}

// Update a QoS
newPriority := 15000
update := &interfaces.QoSUpdate{
    Priority:    &newPriority,
    Description: &"Updated high priority QoS",
}

err = client.QoS().Update(ctx, "high-priority", update)
if err != nil {
    log.Printf("Failed to update QoS: %v", err)
}

// Delete a QoS
err = client.QoS().Delete(ctx, "high-priority")
if err != nil {
    log.Printf("Failed to delete QoS: %v", err)
}
```

## 🛡️ Structured Error Handling

The library provides comprehensive structured error handling with specific error types and codes:

### Error Types
- **SlurmError**: Base error with classification and context
- **NetworkError**: Connection and communication failures
- **AuthenticationError**: Authentication and authorization issues
- **ValidationError**: Request validation failures
- **JobError**: Job-specific operation errors
- **NodeError**: Node-specific operation errors
- **PartitionError**: Partition-specific operation errors

### Error Handling Examples

```go
jobs, err := client.Jobs().List(ctx, options)
if err != nil {
    if slurmErr, ok := err.(*errors.SlurmError); ok {
        switch slurmErr.Code {
        case errors.ErrorCodeUnauthorized:
            // Handle authentication failure
            log.Printf("Authentication failed: %s", slurmErr.Message)
        case errors.ErrorCodeConnectionRefused:
            // Handle connection issues
            log.Printf("Cannot connect to SLURM: %s", slurmErr.Message)
        case errors.ErrorCodeResourceExhausted:
            // Handle resource limitations
            log.Printf("Insufficient resources: %s", slurmErr.Details)
        case errors.ErrorCodeValidationFailed:
            // Handle validation errors
            log.Printf("Request validation failed: %s", slurmErr.Details)
        default:
            log.Printf("SLURM error [%s]: %s", slurmErr.Code, slurmErr.Message)
        }

        // Check if error is retryable
        if slurmErr.IsRetryable() {
            log.Println("Error is retryable, will retry with backoff")
        }

        // Access version-specific information
        if slurmErr.APIVersion != "" {
            log.Printf("Error from API version: %s", slurmErr.APIVersion)
        }
    }
}

// Check for specific error conditions
if errors.IsNetworkError(err) {
    log.Println("Network connectivity issue")
}

if errors.IsAuthenticationError(err) {
    log.Println("Authentication or authorization problem")
}
```

### Error Context and Debugging
```go
// Errors include rich context for debugging
if slurmErr, ok := err.(*errors.SlurmError); ok {
    fmt.Printf("Error Code: %s\\n", slurmErr.Code)
    fmt.Printf("Category: %s\\n", slurmErr.Category)
    fmt.Printf("Message: %s\\n", slurmErr.Message)
    fmt.Printf("Details: %s\\n", slurmErr.Details)
    fmt.Printf("API Version: %s\\n", slurmErr.APIVersion)
    fmt.Printf("Timestamp: %s\\n", slurmErr.Timestamp)
    fmt.Printf("Retryable: %t\\n", slurmErr.IsRetryable())

    // Access underlying cause if available
    if cause := slurmErr.Unwrap(); cause != nil {
        fmt.Printf("Underlying cause: %v\\n", cause)
    }
}
```

### Account Management (v0.0.43+)

```go
// Check if account management is supported
if client.Accounts() == nil {
    log.Println("Account management not supported in this API version")
    return
}

// List all accounts with hierarchy
accountList, err := client.Accounts().List(ctx, &interfaces.ListAccountsOptions{
    WithAssociations: true,
    WithCoordinators: true,
})
if err != nil {
    log.Printf("Failed to list accounts: %v", err)
    return
}

for _, account := range accountList.Accounts {
    fmt.Printf("Account %s: Org=%s, Parent=%s\n",
               account.Name, account.Organization, account.ParentAccount)
}

// Create an account hierarchy
rootAccount := &interfaces.AccountCreate{
    Name:               "research-dept",
    Description:        "Research Department",
    Organization:       "ACME Corp",
    CoordinatorUsers:   []string{"dept-head", "admin"},
    AllowedPartitions:  []string{"compute", "gpu"},
    MaxJobs:            500,
    MaxNodes:           100,
    SharesPriority:     200,
    // Set TRES (Trackable Resource) limits
    MaxTRES: map[string]int{
        "cpu":    2000,
        "mem":    8192000, // 8TB in MB
        "gpu":    20,
    },
    GrpTRES: map[string]int{
        "cpu":    1000,
        "mem":    4096000, // 4TB in MB
        "gpu":    10,
    },
    Flags: []string{"AllowSubmit"},
}

resp, err := client.Accounts().Create(ctx, rootAccount)
if err != nil {
    log.Printf("Failed to create account: %v", err)
} else {
    fmt.Printf("Created account: %s\n", resp.AccountName)
}

// Create a sub-account
subAccount := &interfaces.AccountCreate{
    Name:               "ml-project",
    Description:        "Machine Learning Project",
    ParentAccount:      "research-dept",
    AllowedQoS:         []string{"normal", "high-priority"},
    MaxJobs:            100,
    MaxJobsPerUser:     20,
    DefaultPartition:   "gpu",
    DefaultQoS:         "normal",
}

resp2, err := client.Accounts().Create(ctx, subAccount)
if err != nil {
    log.Printf("Failed to create sub-account: %v", err)
}

// Update account limits
newMaxJobs := 200
update := &interfaces.AccountUpdate{
    MaxJobs:     &newMaxJobs,
    Description: stringPtr("ML Project - Expanded Resources"),
    MaxTRES: map[string]int{
        "cpu": 4000,
        "gpu": 40,
    },
}

err = client.Accounts().Update(ctx, "ml-project", update)
if err != nil {
    log.Printf("Failed to update account: %v", err)
}

// Get account details
account, err := client.Accounts().Get(ctx, "ml-project")
if err != nil {
    log.Printf("Failed to get account: %v", err)
} else {
    fmt.Printf("Account: %s, Max Jobs: %d, Parent: %s\n",
               account.Name, account.MaxJobs, account.ParentAccount)
}
```

## Retry Logic

The client includes built-in retry logic with exponential backoff:

```go
// Custom retry policy
retryPolicy := retry.NewExponentialBackoff().
    WithMaxRetries(5).
    WithMinWaitTime(time.Second).
    WithMaxWaitTime(time.Minute).
    WithJitter(true)

client, err := slurm.NewClient(
    slurm.WithRetryPolicy(retryPolicy),
)
```

## Development

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Linting

```bash
make lint
```

### Generate Documentation

```bash
make docs
```

## 🏗️ Development

### Project Status: **Production Ready** ✅

All core functionality implemented with comprehensive testing:
- ✅ **Multi-version support** for API v0.0.40-v0.0.43
- ✅ **All manager operations** implemented with structured error handling
- ✅ **Enterprise-grade patterns** with proper connection management
- ✅ **Comprehensive test coverage** across critical packages
- ✅ **Quality assurance** with linting, formatting, and build validation
- ✅ **Real-time streaming** via WebSocket and Server-Sent Events

### Building

```bash
make build
```

### Testing

```bash
make test              # Run all tests
make test-coverage     # Generate coverage report
make benchmark         # Run performance benchmarks
```

### Code Quality

```bash
make lint              # Run golangci-lint
make fmt               # Format code with gofmt
make vet               # Run go vet analysis
make check             # Run all quality checks
```

### Code Generation

The project uses code generation to maintain consistency across API versions:

```bash
make install-tools     # Install development tools
make download-specs    # Download OpenAPI specifications
make generate          # Generate all version-specific clients
make generate-version VERSION=v0.0.43  # Generate specific version
```

**Important**: Never edit generated files (marked with "DO NOT EDIT"). See [Code Generation Guide](CODE_GENERATION.md) for details.

### Development Setup

```bash
# Install git hooks for code quality checks
./scripts/install-hooks.sh

# Check generated files haven't been manually edited
./scripts/check-generated-files.sh
```

## 📊 API Compatibility Matrix

Complete compatibility across all active SLURM REST API versions with adapter pattern support:

| Feature | v0.0.40 | v0.0.41 | v0.0.42 | v0.0.43 | Adapter Support | Status |
|---------|---------|---------|---------|---------|-----------------|---------|
| Job Management | ✅ | ✅ | ✅ | ✅ | ✅ Full | Complete |
| Node Management | ✅ | ✅ | ✅ | ✅ | ✅ Full | Complete |
| Partition Management | ✅ | ✅ | ✅ | ✅ | ✅ Full | Complete |
| Cluster Info | ✅ | ✅ | ✅ | ✅ | ✅ Full | Complete |
| Reservation Management | ❌ | ❌ | ❌ | ✅ | ✅ Full | v0.0.43+ |
| QoS Management | ❌ | ❌ | ❌ | ✅ | ✅ Full | v0.0.43+ |
| Account Management | ❌ | ❌ | ❌ | ✅ | ✅ Full | v0.0.43+ |
| Association Management | ❌ | ❌ | ❌ | ✅ | ✅ Full | v0.0.43+ |
| User Management | ❌ | ❌ | ❌ | ✅ | ✅ Full | v0.0.43+ |
| Structured Errors | ✅ | ✅ | ✅ | ✅ | ✅ Enhanced | Complete |
| Auto Version Detection | ✅ | ✅ | ✅ | ✅ | ✅ Enhanced | Complete |
| Type Conversion | ✅ | ✅ | ✅ | ✅ | ✅ Automatic | Complete |
| Performance Optimization | ✅ | ✅ | ✅ | ✅ | ✅ Caching | Complete |

### Adapter Pattern Coverage

The adapter implementation provides comprehensive coverage across all managers:

- **✅ Full Adapter Support**: Complete type conversion, error handling, and version abstraction
- **✅ Enhanced Features**: Adapter-specific improvements like caching and performance monitoring
- **✅ Automatic Conversion**: Seamless translation between internal types and public interfaces
- **✅ Version Context**: All errors include API version information for debugging

### Breaking Change Handling
The library automatically handles breaking changes between versions:
- **Field Renames**: `minimum_switches` → `required_switches` (v0.0.40→v0.0.41)
- **Removed Fields**: `exclusive`, `oversubscribe` from job outputs (v0.0.41→v0.0.42)
- **New Features**: Reservation management in v0.0.43
- **Deprecations**: FrontEnd mode removal in v0.0.43

## 🌊 Real-time Streaming

The library provides WebSocket and Server-Sent Events interfaces for real-time monitoring:

### WebSocket Streaming

```go
import "github.com/jontk/slurm-client/pkg/streaming"

// Create streaming server
wsServer := streaming.NewWebSocketServer(client)
http.HandleFunc("/ws", wsServer.HandleWebSocket)

// JavaScript client
const ws = new WebSocket('ws://localhost:8080/ws');
ws.send(JSON.stringify({
    stream: 'jobs',
    options: { states: ['RUNNING', 'PENDING'] }
}));
```

### Server-Sent Events

```go
// Create SSE server
sseServer := streaming.NewSSEServer(client)
http.HandleFunc("/events", sseServer.HandleSSE)

// JavaScript client
const eventSource = new EventSource('/events?stream=jobs&states=RUNNING');
eventSource.onmessage = (event) => {
    const jobEvent = JSON.parse(event.data);
    console.log('Job update:', jobEvent);
};
```

### Real-time Monitoring

Stream different resource types:
- **Jobs**: State changes, completion, failures
- **Nodes**: Availability, allocation changes
- **Partitions**: Configuration updates

See the [streaming example](../examples/streaming-server/main.go) for a complete web interface.

## 🔧 CLI Tool

A command-line interface is available for easy interaction with SLURM clusters:

### Installation

```bash
go install github.com/jontk/slurm-client/cmd/slurm-cli@latest
```

### Quick Start

```bash
# Set environment variables
export SLURM_REST_URL="https://cluster.example.com:6820"
export SLURM_JWT="your-jwt-token"

# List jobs
slurm-cli jobs list

# Submit a job
slurm-cli submit --command "python train.py" --cpus 4 --memory 8192

# Get job details
slurm-cli jobs get 12345

# List nodes
slurm-cli nodes list --states IDLE

# Show cluster info
slurm-cli info
```

See the [CLI documentation](cli/README.md) for complete usage information.

## 🧪 Testing

### Unit Tests

Run the comprehensive test suite:

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run tests for a specific version
go test ./internal/api/v0_0_43/...
```

### Integration Tests

The library includes integration tests that can be run against a real SLURM cluster:

```bash
# Set up environment
export SLURM_REAL_SERVER_TEST=true
export SLURM_SERVER_URL=http://your-slurm-server:6820
export SLURM_JWT_TOKEN=<your-jwt-token>

# Run real server tests
go test -v ./tests/integration/...

# Run diagnostic script
./scripts/diagnose-slurm-auth.sh
```

See [Real Server Testing Guide](https://github.com/jontk/slurm-client/blob/main/tests/integration/REAL_SERVER_TESTING.md) for detailed setup instructions.

### Known Limitations

When testing against real SLURM servers, some endpoints may return HTTP 502 if slurmdbd is not properly connected. This commonly occurs due to authentication plugin mismatches between slurmctld and slurmdbd (e.g., JWT vs munge). The client handles these scenarios gracefully.

## 🤝 Contributing

We welcome contributions! This project follows enterprise development practices:

1. **Fork and Clone**: Fork the repository and create a feature branch
2. **Development**: Follow the established patterns in `*_impl.go` files
3. **Testing**: Add comprehensive tests with structured error handling
4. **Quality**: Run `make check` to ensure all quality checks pass
5. **Documentation**: Update documentation and examples as needed
6. **Pull Request**: Submit PR with clear description and tests

### Development Guidelines
- Follow the existing multi-version architecture patterns
- Use structured error handling for all new operations
- Maintain backward compatibility within major versions
- Add comprehensive test coverage for new features

## 📚 Documentation

- **[Architecture Documentation](ARCHITECTURE.md)**: Technical design and patterns
- **[API Documentation](api/README.md)**: Version-specific API documentation
- **[Examples](../examples/README.md)**: Practical usage examples and tutorials

## 🏆 Recognition

**The definitive Go client library for SLURM REST API integration**

- **Industry First**: Only multi-version Go SLURM client library
- **Enterprise Grade**: Production patterns from AWS SDK and Kubernetes client-go
- **Community Impact**: Addresses ecosystem fragmentation with unified solution
- **Open Source**: MIT licensed for maximum adoption and contribution

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/jontk/slurm-client/blob/main/LICENSE) file for details.

## 🙏 Acknowledgments

- **SLURM Community**: For the robust REST API and OpenAPI specifications
- **Go Community**: For excellent tooling and enterprise library patterns
- **Enterprise Patterns**: Inspired by AWS SDK, Kubernetes client-go, and other production libraries
- **Contributors**: Thank you to all who have contributed to making this the definitive SLURM Go client

---

**Ready for Production** • **Enterprise-Grade** • **Multi-Version Support** • **Comprehensive Error Handling**
