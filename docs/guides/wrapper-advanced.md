# ⚠️ DEPRECATED

**This documentation is deprecated and will be removed in a future version.**

The wrapper pattern described here has been superseded by the adapter pattern. Please refer to:
- [Architecture Guide](./ARCHITECTURE.md) for current architecture
- [Configuration Guide](./configuration.md) for client setup
- [API Documentation](../api/README.md) for API usage

---

# Advanced Wrapper Usage (Deprecated)

> ⚠️ **Advanced Users Only**: The wrapper pattern provides direct access to version-specific APIs. Most users should use the [Adapter Pattern](../index.md#quick-start-recommended-approach) instead.

## When to Use the Wrapper Pattern

The wrapper pattern is only recommended for:

- **Performance-critical applications** requiring minimal overhead
- **Debugging** API responses and troubleshooting issues
- **Direct API access** when you need version-specific features
- **Migration scenarios** from existing direct API code
- **Advanced use cases** not covered by the adapter layer

## Usage Example

```go
package main

import (
    "context"
    "log"

    "github.com/jontk/slurm-client/internal/api/v0_0_43"
    "github.com/jontk/slurm-client/internal/interfaces"
)

func main() {
    // Must specify exact version
    config := &interfaces.ClientConfig{
        BaseURL: "https://slurm-server:6820",
        Version: "v0.0.43", // Required for wrapper
        Authentication: &interfaces.AuthConfig{
            Type:  "token",
            Token: "your-jwt-token",
        },
    }

    // Create version-specific wrapper
    wrapper, err := v0_0_43.NewWrapperClient(config)
    if err != nil {
        log.Fatal(err)
    }

    // Direct API calls - you handle all conversions
    ctx := context.Background()
    response, err := wrapper.GetJobs(ctx, &v0_0_43.GetJobsParams{})
    if err != nil {
        log.Fatal(err)
    }

    // Manual response handling required
    if response.StatusCode() != 200 {
        log.Fatalf("API error: %d", response.StatusCode())
    }

    jobs := response.JSON200
    if jobs == nil || jobs.Jobs == nil {
        log.Println("No jobs found")
        return
    }

    // Manual type conversion
    for _, job := range *jobs.Jobs {
        if job.JobId != nil && job.Name != nil {
            log.Printf("Job %d: %s", *job.JobId, *job.Name)
        }
    }
}
```

## Key Differences from Adapter Pattern

| Aspect | Adapter Pattern | Wrapper Pattern |
|--------|----------------|-----------------|
| **Version Handling** | Automatic | Manual (must specify) |
| **Type Conversion** | Automatic | Manual |
| **Error Handling** | Structured errors | Raw HTTP responses |
| **API Coverage** | Unified interface | Version-specific APIs |
| **Complexity** | Simple | Complex |
| **Performance** | Optimized | Direct (minimal overhead) |

## Version-Specific Wrappers

Each SLURM API version has its own wrapper:

- `v0_0_40.NewWrapperClient()` - SLURM 21.08.x
- `v0_0_41.NewWrapperClient()` - SLURM 22.05.x
- `v0_0_42.NewWrapperClient()` - SLURM 23.02.x
- `v0_0_43.NewWrapperClient()` - SLURM 23.11.x+

## Error Handling

With the wrapper pattern, you must handle errors manually:

```go
response, err := wrapper.GetJobs(ctx, params)
if err != nil {
    // Network or connection errors
    log.Fatal(err)
}

switch response.StatusCode() {
case 200:
    // Success
    jobs := response.JSON200
case 401:
    // Authentication failed
    log.Fatal("Authentication failed")
case 404:
    // Not found
    log.Fatal("Endpoint not found")
default:
    // Other errors
    log.Fatalf("API error: %d", response.StatusCode())
}
```

## Migration Guide: Wrapper to Adapter

If you're using the wrapper pattern and want to migrate to the adapter pattern:

### Before (Wrapper)
```go
wrapper, err := v0_0_43.NewWrapperClient(config)
response, err := wrapper.GetJobs(ctx, &v0_0_43.GetJobsParams{})
jobs := response.JSON200
```

### After (Adapter)
```go
client, err := factory.NewClient(config) // Auto-detects version
jobs, err := client.Jobs().List(ctx, nil) // Structured response
```

## Best Practices

1. **Use sparingly** - Only when adapter pattern doesn't meet your needs
2. **Handle all error cases** - No automatic error conversion
3. **Version pinning** - Specify exact versions in your dependencies
4. **Testing** - Test against multiple SLURM versions if supporting them
5. **Documentation** - Document why you chose wrapper over adapter

## Performance Considerations

The wrapper pattern has minimal overhead but requires more code:

- **Pros**: Direct API access, no type conversion overhead
- **Cons**: Manual error handling, version-specific code, more complexity

For most applications, the adapter pattern's automatic optimizations outweigh the wrapper's minimal overhead advantage.

## Troubleshooting

Common issues with wrapper pattern:

1. **Version mismatch**: Ensure wrapper version matches your SLURM server
2. **Manual type handling**: Check for nil pointers in responses
3. **Error handling**: Implement comprehensive status code checking
4. **API changes**: Monitor SLURM API changes between versions

## See Also

- [Adapter Pattern Documentation](../index.md#quick-start-recommended-approach)
- [API Reference](../api/README.md)
- [Examples](../../examples/README.md)
- [Architecture Overview](./ARCHITECTURE.md)