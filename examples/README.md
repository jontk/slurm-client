# SLURM Client Examples

This directory contains comprehensive examples demonstrating various features and use cases of the slurm-client library.

## Quick Start Examples

### 🚀 Basic Examples
Start here if you're new to the library:

- [`basic/`](basic/) - **Essential examples for beginners**
  - [Connection & Authentication](basic/01_connect/main.go) - Various ways to connect to SLURM
  - [Listing Jobs](basic/02_list_jobs/main.go) - Query and filter job information
  - [Submitting Jobs](basic/03_submit_job/main.go) - Submit different types of jobs
  - [Error Handling](basic/04_error_handling/main.go) - Robust error handling patterns

### Additional Quick Start
- [`basic-usage/`](multi-version/basic-usage/main.go) - Simple job submission and monitoring
- [`job_list_example.go`](job_list_example.go) - Basic job listing example

### Authentication
- [`authentication/`](multi-version/authentication/main.go) - Different authentication methods (JWT, basic auth, no auth)

### Error Handling
- [`error-handling/`](multi-version/error-handling/main.go) - Comprehensive error handling patterns

## Advanced Examples

### Job Management

#### Array Jobs
- [`array-jobs/`](array-jobs/main.go) - Array job submission and management
  - Submit array jobs with multiple tasks
  - Monitor array job progress
  - Selective task cancellation
  - Task-specific environment variables

#### Job Dependencies
- [`job-dependencies/`](job-dependencies/main.go) - Complex job workflows
  - Simple job chains (A → B → C)
  - Fan-out/fan-in patterns
  - Complex DAG workflows
  - Conditional execution with error handling

#### Batch Operations
- [`batch-operations/`](batch-operations/main.go) - Concurrent job management
  - Batch job submission
  - Progress monitoring
  - Results collection and reporting
  - Cleanup operations

<!--
#### Job Allocation (v0.0.43+)
- [`job-allocation/`](job-allocation/) - Direct resource allocation without scripts
  - Basic resource allocation
  - GPU allocation for interactive work
  - High-memory allocation
  - Constrained allocation with specific requirements
  - Allocation workflows with dependencies
-->

#### Job Workflows
- [`job-workflow/`](multi-version/job-workflow/main.go) - Complete job lifecycle management

### Resource Management

#### Resource Allocation
- [`resource-allocation/`](resource-allocation/main.go) - Advanced resource allocation patterns
  - GPU allocation (single, multi, type-specific)
  - Memory-intensive jobs
  - Node-specific constraints
  - Resource sharing patterns
  - Dynamic resource discovery

### Workload Management

<!--
#### WCKey Management (v0.0.43+)
- [`wckey-management/`](wckey-management/) - Workload Characterization Key management
  - Creating WCKeys for different workload types
  - Listing and filtering WCKeys
  - Using WCKeys in job submission
  - WCKey-based job tracking and accounting
  - Managing WCKeys across users and clusters
-->

### Real-time Monitoring

#### Watch Jobs
- [`watch-jobs/`](watch-jobs/main.go) - Real-time job monitoring
  - State change detection
  - Event-driven notifications
  - Configurable polling intervals

#### Watch Nodes
- [`watch-nodes/`](watch-nodes/main.go) - Real-time node monitoring
  - Node state changes
  - Resource availability tracking
  - Maintenance notifications

## Multi-Version Support

The [`multi-version/`](multi-version/) directory contains examples specifically demonstrating:
- Version-specific features
- API compatibility across versions
- Version detection and selection
- Handling breaking changes between versions

## Running the Examples

### Prerequisites

1. Set up your SLURM REST API endpoint:
```bash
export SLURM_REST_URL="https://your-cluster:6820"
```

2. Configure authentication:
```bash
# For JWT token auth
export SLURM_JWT="your-jwt-token"

# For basic auth
export SLURM_USERNAME="your-username"
export SLURM_PASSWORD="your-password"
```

### Running an Example

```bash
# Navigate to the example directory
cd array-jobs/

# Run the example
go run main.go
```

### Modifying Examples

Each example includes configuration that you should modify:
- `cfg.BaseURL` - Your SLURM REST API endpoint
- Authentication credentials
- Partition names
- Resource requirements (CPUs, memory, time limits)
- Script paths and commands

## Example Categories

### By Use Case

1. **HPC Workflows**
   - Array jobs for parameter sweeps
   - Job dependencies for multi-stage pipelines
   - Resource allocation for GPU/high-memory tasks

2. **Monitoring & Operations**
   - Real-time job tracking
   - Node availability monitoring
   - Batch operations management

3. **Development & Testing**
   - Error handling patterns
   - Multi-version compatibility
   - Authentication methods

### By Complexity

1. **Beginner**
   - Basic usage
   - Simple job submission
   - Job listing

2. **Intermediate**
   - Array jobs
   - Job dependencies
   - Error handling

3. **Advanced**
   - Complex DAG workflows
   - Dynamic resource allocation
   - Real-time monitoring

## Common Patterns

### Job Submission
```go
job := &interfaces.JobSubmission{
    Name:      "my-job",
    Script:    "#!/bin/bash\necho 'Hello SLURM'",
    Partition: "compute",
    CPUs:      4,
    Memory:    8192, // 8GB in MB
    TimeLimit: 60,   // minutes
}
resp, err := client.Jobs().Submit(ctx, job)
```

### Error Handling
```go
if slurmErr, ok := err.(*errors.SlurmError); ok {
    switch slurmErr.Code {
    case errors.ErrorCodeUnauthorized:
        log.Fatal("Authentication failed")
    case errors.ErrorCodeResourceExhausted:
        log.Printf("Insufficient resources: %s", slurmErr.Details)
    }
}
```

### Resource Constraints
```go
job.Metadata = map[string]interface{}{
    "gres":       "gpu:2",        // Request 2 GPUs
    "constraint": "v100&ib",      // V100 GPUs with InfiniBand
    "exclusive":  true,           // Exclusive node access
}
```

## Tips and Best Practices

1. **Always handle errors** - Check for specific error types and codes
2. **Use contexts** - Pass contexts for cancellation and timeouts
3. **Monitor resources** - Check partition and node availability before submission
4. **Set appropriate timeouts** - Don't request more time than needed
5. **Use array jobs** - For similar tasks with different parameters
6. **Plan dependencies** - Use job dependencies for complex workflows
7. **Clean up resources** - Cancel jobs and clean temporary files

## Contributing Examples

When adding new examples:
1. Create a descriptive directory name
2. Include a focused `main.go` demonstrating one concept
3. Add comprehensive comments explaining the code
4. Update this README with your example
5. Test with different SLURM configurations

## Additional Resources

- [SLURM Documentation](https://slurm.schedmd.com/)
- [SLURM REST API Reference](https://slurm.schedmd.com/rest_api.html)
- [slurm-client Documentation](https://pkg.go.dev/github.com/jontk/slurm-client)