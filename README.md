# Slurm Client

A Go client library for interacting with the Slurm REST API. This library provides a clean, type-safe interface for managing Slurm jobs, nodes, and partitions.

## Features

- **Complete API Coverage**: Jobs, nodes, and partitions management
- **Enterprise-grade**: Built with patterns from AWS SDK and other enterprise libraries
- **Configurable**: Flexible configuration via environment variables or code
- **Resilient**: Built-in retry logic with exponential backoff
- **Authentication**: Support for token-based and basic authentication
- **Type-safe**: Full Go struct definitions for all API responses
- **Context-aware**: All operations support context for cancellation and timeouts

## Installation

```bash
go get github.com/jontk/slurm-client
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/jontk/slurm-client"
    "github.com/jontk/slurm-client/pkg/auth"
)

func main() {
    // Create client with token authentication
    client, err := slurm.NewClient(
        slurm.WithBaseURL("https://your-slurm-server:6820"),
        slurm.WithAuth(auth.NewTokenAuth("your-token")),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // List jobs
    jobs, err := client.Jobs().ListJobs(context.Background(), &slurm.ListJobsOptions{
        State: "RUNNING",
        Limit: 10,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d running jobs\\n", len(jobs.Jobs))
    
    // Submit a job
    submission := &slurm.JobSubmission{
        Name:      "test-job",
        Script:    "#!/bin/bash\\necho 'Hello, Slurm!'",
        Partition: "compute",
        CPUs:      2,
        Memory:    4096,
        TimeLimit: 3600,
    }
    
    response, err := client.Jobs().SubmitJob(context.Background(), submission)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Job submitted with ID: %s\\n", response.JobID)
}
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

## API Operations

### Jobs

```go
// List jobs with filtering
jobs, err := client.Jobs().ListJobs(ctx, &slurm.ListJobsOptions{
    UserID:    "username",
    State:     "RUNNING",
    Partition: "compute",
    Limit:     50,
})

// Get specific job
job, err := client.Jobs().GetJob(ctx, "12345")

// Submit job
response, err := client.Jobs().SubmitJob(ctx, &slurm.JobSubmission{
    Name:      "my-job",
    Script:    "#!/bin/bash\\necho 'Hello World'",
    Partition: "compute",
    CPUs:      4,
    Memory:    8192,
})

// Cancel job
err := client.Jobs().CancelJob(ctx, "12345")

// Get job steps
steps, err := client.Jobs().GetJobSteps(ctx, "12345")
```

### Nodes

```go
// List nodes
nodes, err := client.Nodes().ListNodes(ctx, &slurm.ListNodesOptions{
    State:     "IDLE",
    Partition: "compute",
    Features:  []string{"gpu", "ssd"},
})

// Get specific node
node, err := client.Nodes().GetNode(ctx, "node001")

// Update node
err := client.Nodes().UpdateNode(ctx, "node001", &slurm.NodeUpdate{
    State:  "DRAIN",
    Reason: "Maintenance",
})
```

### Partitions

```go
// List partitions
partitions, err := client.Partitions().ListPartitions(ctx)

// Get specific partition
partition, err := client.Partitions().GetPartition(ctx, "compute")
```

## Error Handling

The client returns structured errors that implement the `error` interface:

```go
jobs, err := client.Jobs().ListJobs(ctx, nil)
if err != nil {
    if slurmErr, ok := err.(*slurm.SlurmError); ok {
        fmt.Printf("Slurm API error: %d - %s\\n", slurmErr.Code, slurmErr.Message)
    } else {
        fmt.Printf("Client error: %v\\n", err)
    }
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

## API Compatibility

This client is compatible with Slurm REST API version 0.0.39. The API version can be configured:

```go
config := &config.Config{
    APIVersion: "v0.0.40", // Use different API version
}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built following patterns from enterprise Go libraries like AWS SDK
- Designed for production use with proper error handling and retry logic
- Supports all major Slurm REST API operations