# Basic Examples

This directory contains basic examples demonstrating core functionality of the SLURM REST API Go client.

## Examples Overview

### 01_connect - Connection and Authentication
Demonstrates various ways to connect to a SLURM REST API server:
- Basic JWT authentication
- Environment variable configuration
- Advanced connection options (timeouts, retries, rate limiting)
- Different authentication methods (JWT, API Key, Basic Auth, Certificates)

```bash
go run 01_connect/main.go
```

### 02_list_jobs - Listing and Filtering Jobs
Shows how to retrieve and filter job information:
- List all jobs
- Filter by job state (RUNNING, PENDING, etc.)
- Filter by user
- Advanced filtering with multiple criteria
- Pagination for large result sets

```bash
go run 02_list_jobs/main.go
```

### 03_submit_job - Submitting Jobs
Examples of submitting different types of jobs:
- Simple batch job submission
- Jobs with specific resource requirements
- Array jobs for parallel processing
- Jobs with dependencies

```bash
go run 03_submit_job/main.go
```

### 04_error_handling - Robust Error Handling
Comprehensive error handling patterns:
- Connection error handling
- Authentication error detection
- API error classification
- Timeout and context cancellation
- Retry strategies and recovery

```bash
go run 04_error_handling/main.go
```

## Prerequisites

Before running these examples, ensure you have:

1. **SLURM REST API Server**: A running SLURM cluster with REST API enabled
2. **Authentication Token**: A valid JWT token or other credentials
3. **Go Environment**: Go 1.20 or later installed

## Configuration

Set these environment variables for quick testing:

```bash
export SLURM_API_URL="https://your-cluster:6820"
export SLURM_API_TOKEN="your-jwt-token"
export SLURM_API_VERSION="v0.0.43"  # Optional: force specific version
```

Or modify the connection parameters in each example:

```go
client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL("https://your-cluster:6820"),
    slurm.WithAuth(auth.NewTokenAuth("your-jwt-token")),
)
```

## Common Patterns

### Context Usage
All examples use Go contexts for proper cancellation and timeout handling:

```go
ctx := context.Background()
// or with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

### Error Handling
Check for specific error types to handle different scenarios:

```go
var apiErr *types.APIError
if errors.As(err, &apiErr) {
    if apiErr.IsAuthError() {
        // Handle authentication errors
    }
}
```

### Resource Cleanup
Always close the client when done:

```go
defer client.Close()
```

## Next Steps

After mastering these basic examples, explore:
- [Advanced Examples](../advanced/README.md) - Complex workflows and optimizations
- [Tutorials](../tutorials/README.md) - Step-by-step guides for specific use cases
- [Full Examples List](../README.md) - Complete catalog of all examples

## Troubleshooting

### Connection Refused
- Verify SLURM REST API is running: `systemctl status slurmrestd`
- Check firewall rules for port 6820
- Ensure TLS certificates are valid

### Authentication Failed
- Verify JWT token is not expired
- Check token has necessary permissions
- Ensure correct authentication method is used

### API Version Mismatch
- Let the client auto-detect version by not specifying it
- Check SLURM version with `sinfo --version`
- Use appropriate API version for your SLURM installation