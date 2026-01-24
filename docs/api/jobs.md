# Job Management API

The Job Management API provides comprehensive functionality for managing SLURM jobs.

## Interface

```go
type JobManager interface {
    // List all jobs with optional filters
    List(ctx context.Context, opts *ListJobsOptions) (*JobList, error)

    // Get a specific job by ID
    Get(ctx context.Context, jobID string) (*Job, error)

    // Submit a new job
    Submit(ctx context.Context, job *JobSubmission) (*JobSubmitResponse, error)

    // Cancel a job
    Cancel(ctx context.Context, jobID string) error

    // Hold a job
    Hold(ctx context.Context, jobID string) error

    // Release a held job
    Release(ctx context.Context, jobID string) error

    // Send a signal to a job
    Signal(ctx context.Context, jobID string, signal int) error

    // Notify a job
    Notify(ctx context.Context, jobID string, message string) error

    // Requeue a job
    Requeue(ctx context.Context, jobID string) error

    // Update job properties
    Update(ctx context.Context, jobID string, updates *JobUpdate) error
}
```

## Types

### Job

```go
type Job struct {
    JobID       string
    Name        string
    UserID      string
    GroupID     string
    State       JobState
    Partition   string
    TimeLimit   time.Duration
    SubmitTime  time.Time
    StartTime   *time.Time
    EndTime     *time.Time
    NodeList    string
    NumNodes    int
    NumCPUs     int
    ExitCode    *int
    WorkDir     string
    Command     string
    StdOut      string
    StdErr      string
    Priority    int
    Account     string
    QoS         string
}
```

### JobState

```go
type JobState string

const (
    JobStatePending     JobState = "PENDING"
    JobStateRunning     JobState = "RUNNING"
    JobStateCompleted   JobState = "COMPLETED"
    JobStateFailed      JobState = "FAILED"
    JobStateCancelled   JobState = "CANCELLED"
    JobStateTimeout     JobState = "TIMEOUT"
    JobStateSuspended   JobState = "SUSPENDED"
    JobStateRequeued    JobState = "REQUEUED"
    JobStateHeld        JobState = "HELD"
)
```

### JobSubmission

```go
type JobSubmission struct {
    Name         string
    Script       string
    Partition    string
    Account      string
    TimeLimit    time.Duration
    NumNodes     int
    NumTasks     int
    CPUsPerTask  int
    Memory       string
    WorkDir      string
    Environment  map[string]string
    Constraints  string
    Dependency   string
    QoS          string
    Reservation  string
    Array        string
    OutputFile   string
    ErrorFile    string
    MailUser     string
    MailType     string
}
```

### ListJobsOptions

```go
type ListJobsOptions struct {
    // Filter by state
    States []JobState

    // Filter by user
    Users []string

    // Filter by partition
    Partitions []string

    // Filter by account
    Accounts []string

    // Time range filters
    SubmittedAfter  *time.Time
    SubmittedBefore *time.Time
    StartedAfter    *time.Time
    StartedBefore   *time.Time

    // Pagination
    Limit  int
    Offset int

    // Sorting
    SortBy string
    Order  string
}
```

## Examples

### List All Jobs

```go
jobs, err := client.Jobs().List(ctx, nil)
if err != nil {
    return err
}

for _, job := range jobs.Jobs {
    fmt.Printf("Job %s: %s\n", job.JobID, job.State)
}
```

### List Running Jobs for a User

```go
opts := &interfaces.ListJobsOptions{
    States: []interfaces.JobState{interfaces.JobStateRunning},
    Users:  []string{"jdoe"},
}

jobs, err := client.Jobs().List(ctx, opts)
```

### Submit a Job

```go
job := &interfaces.JobSubmission{
    Name:        "my-analysis",
    Script:      "#!/bin/bash\n#SBATCH --nodes=2\n\nsrun hostname",
    Partition:   "compute",
    TimeLimit:   2 * time.Hour,
    NumNodes:    2,
    NumTasks:    16,
    CPUsPerTask: 4,
    Memory:      "32G",
}

response, err := client.Jobs().Submit(ctx, job)
if err != nil {
    return err
}

fmt.Printf("Submitted job ID: %s\n", response.JobID)
```

### Cancel a Job

```go
err := client.Jobs().Cancel(ctx, "12345")
if err != nil {
    return err
}
```

### Hold and Release a Job

```go
// Hold a job
err := client.Jobs().Hold(ctx, "12345")
if err != nil {
    return err
}

// Do something...

// Release the job
err = client.Jobs().Release(ctx, "12345")
if err != nil {
    return err
}
```

### Requeue a Job

```go
err := client.Jobs().Requeue(ctx, "12345")
if err != nil {
    return err
}
```

### Monitor Job Status

```go
jobID := "12345"

for {
    job, err := client.Jobs().Get(ctx, jobID)
    if err != nil {
        return err
    }

    fmt.Printf("Job %s status: %s\n", jobID, job.State)

    if job.State == interfaces.JobStateCompleted ||
       job.State == interfaces.JobStateFailed ||
       job.State == interfaces.JobStateCancelled {
        break
    }

    time.Sleep(30 * time.Second)
}
```

### Update Job Properties

```go
updates := &interfaces.JobUpdate{
    TimeLimit: 4 * time.Hour,
    Priority:  1000,
}

err := client.Jobs().Update(ctx, "12345", updates)
```

## Error Handling

```go
job, err := client.Jobs().Get(ctx, "12345")
if err != nil {
    var apiErr *interfaces.APIError
    if errors.As(err, &apiErr) {
        if apiErr.Code == 404 {
            fmt.Println("Job not found")
        } else {
            fmt.Printf("API error: %s\n", apiErr.Message)
        }
    }
    return err
}
```

## Best Practices

1. **Use Context**: Always pass appropriate contexts for timeout and cancellation
2. **Handle States**: Check job states before performing operations
3. **Pagination**: Use pagination for large job listings
4. **Error Handling**: Properly handle different error types
5. **Resource Limits**: Be mindful of cluster resource limits when submitting jobs

## See Also

- [API Overview](./README.md)
- [Node Management API](./nodes.md)
- [Job Examples](../../examples/basic/README.md) - Basic job submission and management
- [Job Analytics](../../examples/job-analytics/main.go) - Job performance and analytics examples