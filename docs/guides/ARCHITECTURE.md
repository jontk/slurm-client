# Architecture Overview

This document describes the architecture of the SLURM REST API Client Library.

## Design Principles

1. **Version Independence**: Single interface supporting multiple SLURM API versions
2. **Type Safety**: Strongly typed Go interfaces with comprehensive error handling
3. **Modularity**: Clean separation of concerns with well-defined layers
4. **Extensibility**: Easy to add new API versions without breaking existing code
5. **Performance**: Efficient resource usage with connection pooling and caching

## Architecture Layers

```
┌─────────────────────────────────────────────────────────┐
│                    User Application                      │
├─────────────────────────────────────────────────────────┤
│                  Public Interface Layer                  │
│                 (interfaces package)                     │
├─────────────────────────────────────────────────────────┤
│                    Factory Layer                         │
│              (version detection & routing)               │
├─────────────────────────────────────────────────────────┤
│                  Adapter Interface Layer                 │
│            (common adapter interfaces)                   │
├─────────────────────────────────────────────────────────┤
│                 Implementation Layer                     │
│     ┌──────────┬──────────┬──────────┬──────────┐     │
│     │ v0.0.40  │ v0.0.41  │ v0.0.42  │ v0.0.43  │     │
│     └──────────┴──────────┴──────────┴──────────┘     │
├─────────────────────────────────────────────────────────┤
│                    OpenAPI Clients                       │
│            (auto-generated from specs)                   │
└─────────────────────────────────────────────────────────┘
```

### 1. Public Interface Layer

Located in `pkg/slurm/`, this layer provides the main API surface with type aliases and public contracts:

- **Manager Interfaces**: JobManager, NodeManager, PartitionManager, etc.
- **Data Types**: Job, Node, Partition, Reservation, etc.
- **Configuration**: ClientConfig, AuthConfig, etc.

This package serves as the stable public API that users import and interact with.

```go
type JobManager interface {
    List(ctx context.Context, opts *ListJobsOptions) (*JobList, error)
    Get(ctx context.Context, jobID string) (*Job, error)
    Submit(ctx context.Context, job *JobSubmission) (*JobSubmitResponse, error)
    // ... more methods
}
```

### 2. Factory Layer

Located in `pkg/client/factory/`, responsible for:

- API version auto-detection
- Client instantiation
- Version-specific routing

```go
func NewClient(config *ClientConfig) (Client, error) {
    // Auto-detect version if not specified
    if config.Version == "" {
        config.Version = detectAPIVersion(config.BaseURL)
    }

    // Route to appropriate implementation
    switch config.Version {
    case "v0.0.43":
        return v0_0_43.NewClient(config)
    // ... other versions
    }
}
```

### 3. Adapter Interface Layer

Located in `internal/adapters/common/`, provides:

- Common interfaces for all adapters
- Base functionality shared across versions
- Type conversions between API and public types

```go
type JobAdapter interface {
    List(ctx context.Context, opts *types.JobListOptions) (*types.JobList, error)
    Get(ctx context.Context, jobID int32) (*types.Job, error)
    // ... more methods
}
```

### 4. Implementation Layer

Version-specific implementations in `internal/api/v0_0_**/`:

- **Manager Implementations**: Implement public interfaces
- **Adapter Implementations**: Handle version-specific API calls
- **Type Conversions**: Convert between API types and public types

### 5. OpenAPI Client Layer

Auto-generated clients in `internal/api/v0_0_**/`:

- Generated from OpenAPI specifications
- Provides low-level HTTP API access
- Handles request/response marshaling

## Key Components

### Client Factory

```go
// Auto-detection flow
Client Creation → Version Detection → Implementation Selection → Client Instance
```

### Error Handling

```go
// Error hierarchy
SlurmError
├── ClientError (client-side issues)
├── APIError (server-side issues)
├── ValidationError (input validation)
└── NotImplementedError (unsupported operations)
```

### Type System

```go
// Three-tier type system
Public Types (interfaces) → Adapter Types (common) → API Types (version-specific)
```

## Data Flow

### Read Operation (e.g., List Jobs)

```
1. User calls client.Jobs().List()
2. Factory routes to version-specific JobManager
3. JobManager calls JobAdapter.List()
4. Adapter converts options to API format
5. Adapter calls OpenAPI client
6. OpenAPI client makes HTTP request
7. Response flows back with type conversions
```

### Write Operation (e.g., Submit Job)

```
1. User calls client.Jobs().Submit()
2. Validation occurs at interface layer
3. Factory routes to version-specific JobManager
4. JobManager converts to adapter types
5. Adapter validates version-specific requirements
6. Adapter calls OpenAPI client
7. Response converted back to public types
```

## Extension Points

### Adding New API Versions

1. Add OpenAPI specification to `api/openapi/`
2. Generate client code: `go generate ./...`
3. Implement adapters in `internal/adapters/v0_0_XX/`
4. Implement managers in `internal/api/v0_0_XX/`
5. Update factory to recognize new version

### Adding New Features

1. Define interface in `internal/interfaces/` and export via `pkg/slurm/`
2. Implement in each supported version
3. Provide fallback for unsupported versions
4. Update documentation

## Performance Considerations

### Connection Management

- HTTP client reuse across requests
- Connection pooling configuration
- Automatic retry with backoff

### Caching Strategy

- Version detection results cached
- Cluster configuration cached
- Optional response caching

### Concurrency

- Thread-safe client operations
- Context-based cancellation
- Concurrent request limiting

## Security Architecture

### Authentication Flow

```
Client Config → Auth Provider → Request Interceptor → API Request
```

### Supported Methods

1. **Token Authentication**: JWT tokens
2. **Basic Authentication**: Username/password
3. **Munge Authentication**: System-level auth
4. **Custom Authentication**: Via HTTP client

### Security Best Practices

- No credentials in logs
- Secure credential storage
- Token refresh handling
- TLS/SSL support

## Testing Architecture

### Unit Tests

- Mock interfaces for isolation
- Table-driven test patterns
- High code coverage targets

### Integration Tests

- Real SLURM cluster testing
- Version compatibility matrix
- Performance benchmarks

### Test Organization

```
tests/
├── unit/           # Unit tests
├── integration/    # Integration tests
└── benchmarks/     # Performance tests
```

## Build System

### Code Generation

```bash
# Generate from OpenAPI specs
go generate ./...

# Verify generated code
go build ./...
```

### Continuous Integration

- Multi-version Go testing
- Security scanning
- License compliance
- Documentation generation

## Future Considerations

### Planned Enhancements

1. WebSocket support for real-time updates
2. GraphQL interface option
3. Metric collection integration
4. Distributed tracing support

### Extensibility

- Plugin architecture for custom managers
- Middleware support for requests
- Custom serialization formats
- Event streaming capabilities

## See Also

- [API Documentation](../api/README.md)
- [Code Generation](./CODE_GENERATION.md)
- [Contributing Guide](../development/CONTRIBUTING.md)