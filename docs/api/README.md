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
    "github.com/jontk/slurm-client/pkg/client/factory"
    "github.com/jontk/slurm-client/internal/interfaces"
)

// Create client
config := &interfaces.ClientConfig{
    BaseURL: "http://slurm-host:6820",
}
client, err := factory.NewClient(config)
if err != nil {
    panic(err)
}

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
config.Authentication = &interfaces.AuthConfig{
    Type:  "token",
    Token: "your-jwt-token",
}
```

### Basic Authentication
```go
config.Authentication = &interfaces.AuthConfig{
    Type:     "basic",
    Username: "username",
    Password: "password",
}
```

## Versioning

The client automatically handles API version differences. You can specify a version or let it auto-detect:

```go
// Auto-detect
config := &interfaces.ClientConfig{
    BaseURL: "http://slurm-host:6820",
}

// Specific version
config := &interfaces.ClientConfig{
    BaseURL: "http://slurm-host:6820",
    Version: "v0.0.43",
}
```

## Rate Limiting

The client includes built-in rate limiting:

```go
config.RequestConfig = &interfaces.RequestConfig{
    RateLimit: 100, // requests per second
}
```

## See Also

- [Examples](../../examples/README.md)
- [Configuration Guide](../configuration.md)
- [Troubleshooting](../troubleshooting.md)