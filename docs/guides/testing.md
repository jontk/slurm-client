# Testing and Mocking Guide

This guide covers testing strategies, mock server usage, and best practices for testing applications that use the slurm-client library.

## Table of Contents

- [Testing Strategies](#testing-strategies)
- [Mock Server](#mock-server)
- [Unit Testing](#unit-testing)
- [Integration Testing](#integration-testing)
- [Test Patterns](#test-patterns)
- [Best Practices](#best-practices)

## Testing Strategies

### Test Pyramid

```
        ┌───────────┐
        │    E2E    │ ← Few tests against real SLURM (optional)
        ├───────────┤
        │Integration│ ← Tests with mock server
        ├───────────┤
        │   Unit    │ ← Most tests, fast and isolated
        └───────────┘
```

### When to Use Each Level

**Unit Tests:**
- Business logic that doesn't call SLURM API
- Data transformations and validations
- Error handling logic
- Fast, isolated, no external dependencies

**Integration Tests (Mock Server):**
- Client usage patterns
- API request/response handling
- Version compatibility
- Error scenarios
- Most common testing approach

**E2E Tests (Real SLURM):**
- Final validation before deployment
- Catch real-world edge cases
- Performance testing
- Optional, only when you have test SLURM cluster

## Mock Server

The library includes a full-featured mock SLURM REST API server for testing.

### Features

- ✅ Multi-version support (v0.0.40 through v0.0.44)
- ✅ Comprehensive API coverage (Jobs, Nodes, Partitions, Reservations, QoS)
- ✅ Authentication simulation (token, basic auth)
- ✅ Error scenarios and edge cases
- ✅ Pagination and filtering
- ✅ Response delays for timeout testing

### Quick Start

```go
package myapp_test

import (
    "context"
    "testing"

    slurm "github.com/jontk/slurm-client"
    "github.com/jontk/slurm-client/pkg/auth"
    "github.com/jontk/slurm-client/tests/mocks"
)

func TestMyFeature(t *testing.T) {
    // Create mock server
    mockServer := mocks.NewMockServer("v0.0.43")
    defer mockServer.Close()

    // Create client pointing to mock
    client, err := slurm.NewClient(context.Background(),
        slurm.WithBaseURL(mockServer.URL),
        slurm.WithAuth(auth.NewNoAuth()),
    )
    if err != nil {
        t.Fatalf("Failed to create client: %v", err)
    }

    // Your test code here
    jobs, err := client.Jobs().List(context.Background(), nil)
    if err != nil {
        t.Fatalf("Failed to list jobs: %v", err)
    }

    if len(jobs.Jobs) == 0 {
        t.Error("Expected some jobs")
    }
}
```

### Pre-populating Mock Data

Add test data to the mock server:

```go
import (
    "github.com/jontk/slurm-client/tests/mocks"
    builder "github.com/jontk/slurm-client/tests/mocks/generated/v0_0_43"
)

func TestJobWorkflow(t *testing.T) {
    mockServer := mocks.NewMockServer("v0.0.43")
    defer mockServer.Close()

    // Pre-populate with test jobs
    job1 := builder.NewJobInfo().
        WithJobId(1001).
        WithName("test-job-1").
        WithUserId(1000).
        WithJobState("RUNNING").
        WithCpus(4).
        WithPartition("compute").
        Build()

    job2 := builder.NewJobInfo().
        WithJobId(1002).
        WithName("test-job-2").
        WithUserId(1000).
        WithJobState("PENDING").
        WithCpus(2).
        WithPartition("compute").
        Build()

    mockServer.AddJob(job1)
    mockServer.AddJob(job2)

    // Create client and test
    client, _ := slurm.NewClient(context.Background(),
        slurm.WithBaseURL(mockServer.URL),
        slurm.WithAuth(auth.NewNoAuth()),
    )

    jobs, err := client.Jobs().List(context.Background(), nil)
    if err != nil {
        t.Fatalf("Failed to list jobs: %v", err)
    }

    if len(jobs.Jobs) != 2 {
        t.Errorf("Expected 2 jobs, got %d", len(jobs.Jobs))
    }
}
```

### Testing Error Scenarios

Mock server can simulate errors:

```go
func TestErrorHandling(t *testing.T) {
    mockServer := mocks.NewMockServer("v0.0.43")
    defer mockServer.Close()

    // Configure mock to return errors
    mockServer.SetError("/slurm/v0.0.43/jobs", 500, "Internal Server Error")

    client, _ := slurm.NewClient(context.Background(),
        slurm.WithBaseURL(mockServer.URL),
        slurm.WithAuth(auth.NewNoAuth()),
    )

    _, err := client.Jobs().List(context.Background(), nil)
    if err == nil {
        t.Error("Expected error but got none")
    }

    // Check error type
    if slurmErr, ok := err.(*errors.SlurmError); ok {
        if slurmErr.Code != errors.ErrorCodeInternalError {
            t.Errorf("Expected internal error, got %s", slurmErr.Code)
        }
    }
}
```

### Testing Timeouts

Simulate slow responses:

```go
func TestTimeout(t *testing.T) {
    mockServer := mocks.NewMockServer("v0.0.43")
    defer mockServer.Close()

    // Add 2 second delay to responses
    mockServer.SetDelay(2 * time.Second)

    client, _ := slurm.NewClient(context.Background(),
        slurm.WithBaseURL(mockServer.URL),
        slurm.WithAuth(auth.NewNoAuth()),
        slurm.WithTimeout(1*time.Second), // Shorter timeout
    )

    // This should timeout
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()

    _, err := client.Jobs().List(ctx, nil)
    if err == nil {
        t.Error("Expected timeout error")
    }

    if ctx.Err() != context.DeadlineExceeded {
        t.Error("Expected context deadline exceeded")
    }
}
```

### Multi-Version Testing

Test compatibility across API versions:

```go
func TestMultiVersionCompatibility(t *testing.T) {
    versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43", "v0.0.44"}

    for _, version := range versions {
        t.Run(version, func(t *testing.T) {
            mockServer := mocks.NewMockServer(version)
            defer mockServer.Close()

            client, err := slurm.NewClientWithVersion(context.Background(), version,
                slurm.WithBaseURL(mockServer.URL),
                slurm.WithAuth(auth.NewNoAuth()),
            )
            if err != nil {
                t.Fatalf("Failed to create client for %s: %v", version, err)
            }

            // Test basic operations work across versions
            jobs, err := client.Jobs().List(context.Background(), nil)
            if err != nil {
                t.Errorf("Job listing failed for %s: %v", version, err)
            }

            if jobs == nil {
                t.Errorf("Expected job list for %s", version)
            }
        })
    }
}
```

## Unit Testing

### Testing Business Logic

Test your application logic without calling SLURM:

```go
package myapp

import (
    "testing"
    slurm "github.com/jontk/slurm-client"
)

// Business logic function to test
func FilterHighMemoryJobs(jobs []*slurm.Job, minMemory int64) []*slurm.Job {
    var filtered []*slurm.Job
    for _, job := range jobs {
        if job.Memory >= minMemory {
            filtered = append(filtered, job)
        }
    }
    return filtered
}

// Unit test - no SLURM calls
func TestFilterHighMemoryJobs(t *testing.T) {
    jobs := []*slurm.Job{
        {Name: "job1", Memory: 1 * 1024 * 1024 * 1024},  // 1GB
        {Name: "job2", Memory: 8 * 1024 * 1024 * 1024},  // 8GB
        {Name: "job3", Memory: 16 * 1024 * 1024 * 1024}, // 16GB
    }

    minMemory := int64(4 * 1024 * 1024 * 1024) // 4GB threshold
    filtered := FilterHighMemoryJobs(jobs, minMemory)

    if len(filtered) != 2 {
        t.Errorf("Expected 2 high-memory jobs, got %d", len(filtered))
    }

    if filtered[0].Name != "job2" || filtered[1].Name != "job3" {
        t.Error("Unexpected jobs in filtered list")
    }
}
```

### Interface-Based Testing

Use interfaces for easier testing:

```go
// Define interface for testability
type JobLister interface {
    List(ctx context.Context, opts *slurm.ListJobsOptions) (*slurm.JobList, error)
}

// Your application code
type JobMonitor struct {
    client JobLister
}

func (m *JobMonitor) GetRunningJobs(ctx context.Context) ([]*slurm.Job, error) {
    opts := &slurm.ListJobsOptions{
        States: []string{"RUNNING"},
    }
    result, err := m.client.List(ctx, opts)
    if err != nil {
        return nil, err
    }
    return result.Jobs, nil
}

// Mock implementation for testing
type mockJobLister struct {
    jobs []*slurm.Job
    err  error
}

func (m *mockJobLister) List(ctx context.Context, opts *slurm.ListJobsOptions) (*slurm.JobList, error) {
    if m.err != nil {
        return nil, m.err
    }
    return &slurm.JobList{Jobs: m.jobs}, nil
}

// Unit test with mock
func TestGetRunningJobs(t *testing.T) {
    mock := &mockJobLister{
        jobs: []*slurm.Job{
            {Name: "job1", State: "RUNNING"},
            {Name: "job2", State: "RUNNING"},
        },
    }

    monitor := &JobMonitor{client: mock}
    jobs, err := monitor.GetRunningJobs(context.Background())

    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }

    if len(jobs) != 2 {
        t.Errorf("Expected 2 jobs, got %d", len(jobs))
    }
}
```

## Integration Testing

### Testing Against Real SLURM

For final validation, test against a real SLURM cluster:

```go
func TestRealServer(t *testing.T) {
    // Skip if not in integration test mode
    if os.Getenv("SLURM_REAL_SERVER_TEST") != "true" {
        t.Skip("Skipping real server test. Set SLURM_REAL_SERVER_TEST=true to run")
    }

    // Get real server configuration from environment
    serverURL := os.Getenv("SLURM_SERVER_URL")
    token := os.Getenv("SLURM_JWT_TOKEN")

    if serverURL == "" || token == "" {
        t.Fatal("SLURM_SERVER_URL and SLURM_JWT_TOKEN must be set for real server tests")
    }

    client, err := slurm.NewClient(context.Background(),
        slurm.WithBaseURL(serverURL),
        slurm.WithAuth(auth.NewTokenAuth(token)),
    )
    if err != nil {
        t.Fatalf("Failed to create client: %v", err)
    }

    // Test basic operations
    t.Run("Ping", func(t *testing.T) {
        err := client.Info().Ping(context.Background())
        if err != nil {
            t.Errorf("Ping failed: %v", err)
        }
    })

    t.Run("ListJobs", func(t *testing.T) {
        jobs, err := client.Jobs().List(context.Background(), nil)
        if err != nil {
            t.Errorf("List jobs failed: %v", err)
        }
        t.Logf("Found %d jobs", len(jobs.Jobs))
    })
}
```

### Running Integration Tests

```bash
# Mock server tests (default)
go test -v ./...

# Real server tests
export SLURM_REAL_SERVER_TEST=true
export SLURM_SERVER_URL=http://your-slurm:6820
export SLURM_JWT_TOKEN=your-token
go test -v ./... -run TestRealServer
```

## Test Patterns

### Table-Driven Tests

Test multiple scenarios efficiently:

```go
func TestJobSubmission(t *testing.T) {
    tests := []struct {
        name      string
        job       *slurm.JobSubmission
        wantError bool
    }{
        {
            name: "valid job",
            job: &slurm.JobSubmission{
                Name:      "test-job",
                Script:    "#!/bin/bash\necho hello",
                Partition: "compute",
                CPUs:      4,
            },
            wantError: false,
        },
        {
            name: "missing script",
            job: &slurm.JobSubmission{
                Name:      "test-job",
                Partition: "compute",
                CPUs:      4,
            },
            wantError: true,
        },
        {
            name: "invalid partition",
            job: &slurm.JobSubmission{
                Name:      "test-job",
                Script:    "#!/bin/bash\necho hello",
                Partition: "nonexistent",
                CPUs:      4,
            },
            wantError: true,
        },
    }

    mockServer := mocks.NewMockServer("v0.0.43")
    defer mockServer.Close()

    client, _ := slurm.NewClient(context.Background(),
        slurm.WithBaseURL(mockServer.URL),
        slurm.WithAuth(auth.NewNoAuth()),
    )

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := client.Jobs().Submit(context.Background(), tt.job)

            if tt.wantError && err == nil {
                t.Error("Expected error but got none")
            }
            if !tt.wantError && err != nil {
                t.Errorf("Unexpected error: %v", err)
            }
        })
    }
}
```

### Helper Functions

Create reusable test helpers:

```go
// testhelpers/client.go
package testhelpers

import (
    "context"
    "testing"

    slurm "github.com/jontk/slurm-client"
    "github.com/jontk/slurm-client/pkg/auth"
    "github.com/jontk/slurm-client/tests/mocks"
)

// NewTestClient creates a client with mock server for testing
func NewTestClient(t *testing.T, version string) (slurm.Client, func()) {
    t.Helper()

    mockServer := mocks.NewMockServer(version)

    client, err := slurm.NewClient(context.Background(),
        slurm.WithBaseURL(mockServer.URL),
        slurm.WithAuth(auth.NewNoAuth()),
    )
    if err != nil {
        t.Fatalf("Failed to create test client: %v", err)
    }

    cleanup := func() {
        mockServer.Close()
    }

    return client, cleanup
}

// Use in tests:
func TestMyFeature(t *testing.T) {
    client, cleanup := testhelpers.NewTestClient(t, "v0.0.43")
    defer cleanup()

    // Your test code
}
```

### Subtest Isolation

Use subtests for better isolation:

```go
func TestJobOperations(t *testing.T) {
    client, cleanup := testhelpers.NewTestClient(t, "v0.0.43")
    defer cleanup()

    t.Run("List", func(t *testing.T) {
        jobs, err := client.Jobs().List(context.Background(), nil)
        if err != nil {
            t.Fatalf("List failed: %v", err)
        }
        if jobs == nil {
            t.Error("Expected non-nil job list")
        }
    })

    t.Run("Submit", func(t *testing.T) {
        job := &slurm.JobSubmission{
            Name:   "test",
            Script: "#!/bin/bash\necho test",
        }
        resp, err := client.Jobs().Submit(context.Background(), job)
        if err != nil {
            t.Fatalf("Submit failed: %v", err)
        }
        if resp.JobID == "" {
            t.Error("Expected job ID")
        }
    })
}
```

## Best Practices

### 1. Use Mock Server for Most Tests

```go
// ✅ Good: Fast, reliable, no external dependencies
func TestFeature(t *testing.T) {
    mockServer := mocks.NewMockServer("v0.0.43")
    defer mockServer.Close()
    // ... test code
}

// ❌ Bad: Slow, requires real SLURM, can fail due to external issues
func TestFeature(t *testing.T) {
    client, _ := slurm.NewClient(ctx,
        slurm.WithBaseURL("http://production-slurm:6820"),
        // ...
    )
}
```

### 2. Test Error Handling

```go
func TestErrorHandling(t *testing.T) {
    mockServer := mocks.NewMockServer("v0.0.43")
    defer mockServer.Close()

    // Test various error scenarios
    scenarios := []struct {
        name       string
        setupError func()
        operation  func() error
    }{
        {
            name: "401 Unauthorized",
            setupError: func() {
                mockServer.SetError("/slurm/v0.0.43/jobs", 401, "Unauthorized")
            },
            operation: func() error {
                _, err := client.Jobs().List(ctx, nil)
                return err
            },
        },
        // ... more scenarios
    }

    for _, sc := range scenarios {
        t.Run(sc.name, func(t *testing.T) {
            sc.setupError()
            err := sc.operation()
            if err == nil {
                t.Error("Expected error")
            }
        })
    }
}
```

### 3. Use Context for Timeouts

```go
func TestWithTimeout(t *testing.T) {
    client, cleanup := testhelpers.NewTestClient(t, "v0.0.43")
    defer cleanup()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    jobs, err := client.Jobs().List(ctx, nil)
    if err != nil {
        t.Fatalf("List failed: %v", err)
    }
    // ...
}
```

### 4. Test Version Compatibility

```go
func TestBackwardCompatibility(t *testing.T) {
    versions := []string{"v0.0.42", "v0.0.43", "v0.0.44"}

    for _, version := range versions {
        t.Run(version, func(t *testing.T) {
            client, cleanup := testhelpers.NewTestClient(t, version)
            defer cleanup()

            // Test core operations work across versions
            jobs, err := client.Jobs().List(context.Background(), nil)
            if err != nil {
                t.Errorf("Basic operation failed for %s: %v", version, err)
            }
            if jobs == nil {
                t.Errorf("Unexpected nil result for %s", version)
            }
        })
    }
}
```

### 5. Clean Up Resources

```go
func TestResourceCleanup(t *testing.T) {
    mockServer := mocks.NewMockServer("v0.0.43")
    defer mockServer.Close() // Always use defer

    client, err := slurm.NewClient(context.Background(),
        slurm.WithBaseURL(mockServer.URL),
        slurm.WithAuth(auth.NewNoAuth()),
    )
    if err != nil {
        t.Fatalf("Setup failed: %v", err)
    }
    defer client.Close() // Clean up client resources

    // Test code
}
```

### 6. Use Parallel Tests When Possible

```go
func TestParallel(t *testing.T) {
    tests := []struct{
        name string
        // ...
    }{
        // test cases
    }

    for _, tt := range tests {
        tt := tt // capture range variable
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel() // Run tests in parallel

            client, cleanup := testhelpers.NewTestClient(t, "v0.0.43")
            defer cleanup()

            // Test code
        })
    }
}
```

## See Also

- [Integration Tests README](../../tests/integration/README.md) - Full integration test documentation
- [Mock Builders](../../tests/mocks/generated/README.md) - Generated mock builders
- [Troubleshooting Guide](../troubleshooting.md) - Debugging test failures
- [Examples](../../examples/) - Real-world usage examples
