# Slurm Client Multi-Version Examples

This directory contains comprehensive examples demonstrating the multi-version Slurm REST API client library. These examples showcase production-ready patterns for different scenarios and use cases.

## 📁 Example Categories

### 🚀 [Basic Usage](./basic-usage/)
**Entry-level examples for getting started**
- Automatic version detection
- Explicit version selection  
- Slurm release compatibility
- Configuration-based setup

```bash
cd basic-usage && go run main.go
```

### 💼 [Job Workflow](./job-workflow/)
**Production job submission and monitoring patterns**
- Simple job submission
- Batch/array job processing
- Resource-aware job optimization
- Real-time job monitoring with progress tracking

```bash
cd job-workflow && go run main.go
```

### 🔄 [Version Features](./version-features/)
**Version-specific capabilities and compatibility**
- API version feature comparison
- Breaking changes demonstration
- Compatibility matrix visualization
- Version-specific implementation patterns

```bash
cd version-features && go run main.go
```

### 🛡️ [Error Handling](./error-handling/)
**Robust error handling and resilience patterns**
- Comprehensive error classification
- Advanced retry policies with exponential backoff
- Circuit breaker pattern implementation
- Graceful degradation strategies
- Timeout and cancellation handling

```bash
cd error-handling && go run main.go
```

### 🔐 [Authentication](./authentication/)
**Complete authentication provider examples**
- No authentication (public endpoints)
- Token-based authentication with refresh
- Basic authentication
- Environment-based configuration
- Configuration file management
- Custom authentication providers (OAuth, API keys)

```bash
cd authentication && go run main.go
```

## 🎯 Quick Start Guide

### 1. **Choose Your API Version Strategy**

```go
// Option A: Automatic detection (recommended)
client, err := slurm.NewClient(ctx, options...)

// Option B: Explicit version for consistency
client, err := slurm.NewClientWithVersion(ctx, "v0.0.42", options...)

// Option C: Best version for your Slurm release
client, err := slurm.NewClientForSlurmVersion(ctx, "25.05", options...)
```

### 2. **Configure Authentication**

```go
// Environment-based (production recommended)
client, err := slurm.NewClientFromEnvironment(ctx)

// Token authentication
client, err := slurm.NewClient(ctx,
    slurm.WithAuth(auth.NewTokenAuth("your-token")),
    // ... other options
)

// Basic authentication  
client, err := slurm.NewClient(ctx,
    slurm.WithAuth(auth.NewBasicAuth("user", "pass")),
    // ... other options
)
```

### 3. **Handle Errors Gracefully**

```go
client, err := slurm.NewClient(ctx,
    slurm.WithRetryPolicy(3, "1s", "10s"),  // Max retries, base delay, max delay
    slurm.WithTimeout("30s"),               // Request timeout
    slurm.WithDebug(true),                  // Enable debug logging
)
```

## 📊 API Version Compatibility

| Slurm Version | Recommended API | Status | Features |
|---------------|-----------------|---------|----------|
| 24.05         | v0.0.40        | Legacy  | Basic job management |
| 24.11         | v0.0.41        | Stable  | Enhanced filtering |
| 25.05         | v0.0.42        | **Stable** | Performance optimized |
| 25.11+        | v0.0.43        | Latest  | Reservation management |

## 🔧 Environment Configuration

Set these environment variables for seamless configuration:

```bash
export SLURM_REST_URL="https://your-slurm-server:6820"
export SLURM_AUTH_TYPE="token"                    # token, basic, none
export SLURM_TOKEN="your-api-token"               # For token auth
export SLURM_USERNAME="your-username"             # For basic auth  
export SLURM_PASSWORD="your-password"             # For basic auth
export SLURM_TIMEOUT="30s"                        # Request timeout
export SLURM_MAX_RETRIES="3"                      # Retry attempts
export SLURM_DEBUG="false"                        # Debug logging
export SLURM_INSECURE_SKIP_VERIFY="false"         # TLS verification
```

## 🏗️ Production Deployment Patterns

### High Availability Setup
```go
// Primary + fallback configuration
primary, _ := slurm.NewClientWithVersion(ctx, slurm.LatestVersion(), primaryOptions...)
fallback, _ := slurm.NewClientWithVersion(ctx, slurm.StableVersion(), fallbackOptions...)

// Use with circuit breaker pattern (see error-handling example)
```

### Resource-Optimized Workflow
```go
// Check available resources before submission
partitions, _ := client.ListPartitions(ctx)
selectedPartition := selectOptimalPartition(partitions, jobRequirements)

jobReq := &slurm.JobSubmissionRequest{
    Partition: selectedPartition,
    // ... optimized resource allocation
}
```

### Multi-Tenant Configuration
```go
// Per-tenant clients with different auth
tenantClients := make(map[string]slurm.SlurmClient)
for tenantID, tenantConfig := range tenants {
    tenantClients[tenantID], _ = slurm.NewClientFromConfig(ctx, tenantConfig)
}
```

## 📈 Performance Best Practices

1. **Connection Pooling**: Reuse client instances across requests
2. **Retry Policies**: Configure exponential backoff with jitter
3. **Timeouts**: Set appropriate timeouts for your environment
4. **Batch Operations**: Group related API calls when possible
5. **Resource Monitoring**: Monitor partition availability before job submission

## 🧪 Testing Your Integration

```bash
# Test basic connectivity
go run basic-usage/main.go

# Test job submission workflow  
go run job-workflow/main.go

# Test error handling resilience
go run error-handling/main.go

# Test authentication methods
go run authentication/main.go

# Test version-specific features
go run version-features/main.go
```

## 📚 Additional Resources

- **[Main Project Documentation](../../README.md)**: Core library documentation
- **[API Reference](../../docs/api/)**: Complete API documentation  
- **[Configuration Guide](../../docs/configuration.md)**: Detailed configuration options
- **[Deployment Guide](../../docs/deployment.md)**: Production deployment patterns
- **[Troubleshooting Guide](../../docs/troubleshooting.md)**: Common issues and solutions

## 🤝 Contributing

Found an issue or want to add more examples? See our [contribution guidelines](../../CONTRIBUTING.md).

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.

---

*These examples are designed for production use and follow enterprise-grade patterns from AWS SDK and Kubernetes client libraries.*