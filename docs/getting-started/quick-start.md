# Quick Start

This guide will get you up and running with the SLURM REST API Client Library in minutes.

## Basic Usage

### 1. Create a Client

The simplest way to create a client with automatic version detection:

```go
package main

import (
    "context"
    "log"

    "github.com/jontk/slurm-client"
    "github.com/jontk/slurm-client/pkg/auth"
)

func main() {
    ctx := context.Background()

    // Create client with automatic version detection
    client, err := slurm.NewClient(ctx,
        slurm.WithBaseURL("https://your-slurm-host:6820"),
        slurm.WithAuth(auth.NewTokenAuth("your-jwt-token")),
    )
    if err != nil {
        log.Fatal(err)
    }

    defer client.Close()

    // Your code here
}
```

### 2. List Jobs

```go
// List all jobs
jobs, err := client.Jobs().List(ctx, nil)
if err != nil {
    log.Fatal(err)
}

for _, job := range jobs.Jobs {
    fmt.Printf("Job %d: %s (State: %s)\n",
        job.JobID,
        job.Name,
        job.JobState,
    )
}
```

### 3. Get Job Details

```go
// Get specific job
job, err := client.Jobs().Get(ctx, 12345)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Job: %s\n", job.Name)
fmt.Printf("User: %s\n", job.UserName)
fmt.Printf("Nodes: %v\n", job.Nodes)
```

### 4. List Nodes

```go
// List all nodes
nodes, err := client.Nodes().List(ctx, nil)
if err != nil {
    log.Fatal(err)
}

for _, node := range nodes.Nodes {
    fmt.Printf("Node %s: %s (%d CPUs)\n",
        node.Name,
        node.State,
        node.CPUs,
    )
}
```

### 5. Submit a Job

```go
// Create job submission request
jobReq := &interfaces.JobSubmitRequest{
    Script: "#!/bin/bash\nsleep 60",
    Job: interfaces.JobDescriptor{
        Name:      "test-job",
        Partition: "compute",
        Nodes:     1,
        CPUs:      4,
    },
}

// Submit the job
jobID, err := client.Jobs().Submit(ctx, jobReq)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Submitted job ID: %d\n", jobID)
```

## Using Filters

Many list operations support filters:

```go
// List only running jobs for a specific user
filter := &interfaces.JobFilter{
    States: []string{"RUNNING"},
    Users:  []uint32{1000},
}

jobs, err := client.Jobs().List(ctx, filter)
```

## Error Handling

The library provides structured error handling:

```go
jobs, err := client.Jobs().List(ctx, nil)
if err != nil {
    // Check for specific error types
    if errors.IsNotFound(err) {
        log.Println("No jobs found")
        return
    }

    if errors.IsAuthenticationError(err) {
        log.Println("Authentication failed")
        return
    }

    if errors.IsVersionNotSupported(err) {
        log.Printf("Operation not supported in API version %s", client.Version())
        return
    }

    log.Fatal(err)
}
```

## Context Support

All operations support context for timeouts and cancellation:

```go
// Set a timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

jobs, err := client.Jobs().List(ctx, nil)
if err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        log.Println("Request timed out")
    }
    return
}
```

## Complete Example

Here's a complete working example:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/jontk/slurm-client"
    "github.com/jontk/slurm-client/pkg/auth"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Create client
    client, err := slurm.NewClient(ctx,
        slurm.WithBaseURL("https://your-slurm-host:6820"),
        slurm.WithAuth(auth.NewTokenAuth("your-jwt-token")),
        slurm.WithRetry(retry.NewExponentialBackoff(3, time.Second)),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // List jobs
    jobs, err := client.Jobs().List(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d jobs\n", len(jobs.Jobs))

    // List nodes
    nodes, err := client.Nodes().List(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d nodes\n", len(nodes.Nodes))

    // List partitions
    partitions, err := client.Partitions().List(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d partitions\n", len(partitions.Partitions))
}
```

## Next Steps

- Explore more [Examples](../../examples/README.md)
- Learn about [Configuration](../configuration.md) options
- Read the [API Reference](../api/README.md)
- Check out the [CLI Reference](../cli/README.md) for command-line usage
