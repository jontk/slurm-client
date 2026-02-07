# API Documentation

This directory contains detailed API documentation for the SLURM REST API Client Library.

## Quick Links

- [Job Management API](./jobs.md)
- [Node Management API](./nodes.md)
- [Partition Management API](./partitions.md)
- [Reservation Management API](./reservations.md)
- [Account Management API](./accounts.md)
- [QoS Management API](./qos.md)

## Overview

The SLURM REST API Client Library provides a unified interface for interacting with SLURM clusters. All API methods follow consistent patterns:

### Common Patterns

1. **Context Support**: All methods accept a `context.Context` as the first parameter
2. **Options Pattern**: List methods accept optional filter parameters
3. **Error Handling**: All methods return structured errors
4. **Type Safety**: All inputs and outputs are strongly typed

### Basic Usage

```go
import (
    "context"
    slurm "github.com/jontk/slurm-client"
    "github.com/jontk/slurm-client/pkg/auth"
)

// Create client
client, err := slurm.NewClient(context.Background(),
    slurm.WithBaseURL("http://slurm-host:6820"),
    slurm.WithAuth(auth.NewTokenAuth("your-token")),
)
if err != nil {
    panic(err)
}
defer client.Close()

// Use the API
ctx := context.Background()
jobs, err := client.Jobs().List(ctx, nil)
```

## API Reference

### Client Interface

The main client interface provides access to all resource managers:

```go
type Client interface {
    // Resource managers
    Jobs() JobManager
    Nodes() NodeManager
    Partitions() PartitionManager
    Reservations() ReservationManager
    Accounts() AccountManager
    QoS() QoSManager

    // Utility methods
    Info() InfoManager
    Config() ConfigManager
}
```

### Common Types

#### ListOptions

Most list methods accept options for filtering and pagination:

```go
type ListOptions struct {
    // Pagination
    Limit  int
    Offset int

    // Sorting
    SortBy string
    Order  string // "asc" or "desc"

    // Filtering
    Filters map[string]string
}
```

#### Error Types

All errors implement the `error` interface with additional context:

```go
type APIError struct {
    Code    int
    Message string
    Details map[string]interface{}
}
```

## Authentication

The client supports multiple authentication methods:

### Token Authentication
```go
client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL("http://slurm-host:6820"),
    slurm.WithAuth(auth.NewTokenAuth("your-jwt-token")),
)
```

### Basic Authentication
```go
client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL("http://slurm-host:6820"),
    slurm.WithAuth(auth.NewBasicAuth("username", "password")),
)
```

## Versioning

The client automatically handles API version differences. You can specify a version or let it auto-detect:

```go
// Auto-detect (recommended)
client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL("http://slurm-host:6820"),
    slurm.WithAuth(auth.NewTokenAuth("token")),
)

// Specific version
client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
    slurm.WithBaseURL("http://slurm-host:6820"),
    slurm.WithAuth(auth.NewTokenAuth("token")),
)
```

## Rate Limiting

The client includes built-in rate limiting and retry logic. Configure with client options:

```go
client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL("http://slurm-host:6820"),
    slurm.WithAuth(auth.NewTokenAuth("token")),
    slurm.WithTimeout(30*time.Second),
    slurm.WithMaxRetries(3),
)
```

## See Also

- [Examples](../../examples/README.md)
- [Configuration Guide](../configuration.md)
- [Troubleshooting](../troubleshooting.md)