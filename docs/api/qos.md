# QoS Management API

The Quality of Service (QoS) Management API provides functionality for managing SLURM QoS policies.

## Interface

```go
type QoSManager interface {
    // List all QoS with optional filters
    List(ctx context.Context, opts *ListQoSOptions) (*QoSList, error)

    // Get a specific QoS by name
    Get(ctx context.Context, qosName string) (*QoS, error)

    // Create a new QoS
    Create(ctx context.Context, qos *QoSCreate) error

    // Update QoS properties
    Update(ctx context.Context, qosName string, updates *QoSUpdate) error

    // Delete a QoS
    Delete(ctx context.Context, qosName string) error
}
```

## Types

### QoS

```go
type QoS struct {
    Name              string
    Description       string
    Priority          int
    PreemptMode       string
    PreemptExemptTime time.Duration
    GraceTime         time.Duration
    Flags             []string
    UsageFactor       float64
    GrpTRES           map[string]int64
    GrpJobs           int
    GrpSubmitJobs     int
    GrpWall           time.Duration
    MaxTRES           map[string]int64
    MaxTRESPerUser    map[string]int64
    MaxJobs           int
    MaxJobsPerUser    int
    MaxSubmitJobs     int
    MaxSubmitJobsPerUser int
    MaxWallDuration   time.Duration
    MinPriority       int
    MinTRES           map[string]int64
}
```

### QoSCreate

```go
type QoSCreate struct {
    Name            string
    Description     string
    Priority        int
    PreemptMode     string
    Flags           []string
    UsageFactor     float64
    MaxJobs         int
    MaxWallDuration time.Duration
    GrpTRES         map[string]int64
}
```

### QoSUpdate

```go
type QoSUpdate struct {
    Description     *string
    Priority        *int
    PreemptMode     *string
    Flags           []string
    UsageFactor     *float64
    MaxJobs         *int
    MaxWallDuration *time.Duration
    GrpTRES         map[string]int64
}
```

### ListQoSOptions

```go
type ListQoSOptions struct {
    // Filter by name pattern
    NamePattern string

    // Filter by priority range
    MinPriority *int
    MaxPriority *int

    // Pagination
    Limit  int
    Offset int

    // Sorting
    SortBy string
    Order  string
}
```

## Examples

### List All QoS

```go
qosList, err := client.QoS().List(ctx, nil)
if err != nil {
    return err
}

for _, qos := range qosList.QoS {
    fmt.Printf("QoS %s: Priority=%d, MaxJobs=%d\n",
        qos.Name, qos.Priority, qos.MaxJobs)
}
```

### Get QoS Details

```go
qos, err := client.QoS().Get(ctx, "high-priority")
if err != nil {
    return err
}

fmt.Printf("QoS: %s\n", qos.Name)
fmt.Printf("Description: %s\n", qos.Description)
fmt.Printf("Priority: %d\n", qos.Priority)
fmt.Printf("Max Wall Duration: %s\n", qos.MaxWallDuration)
fmt.Printf("Usage Factor: %.2f\n", qos.UsageFactor)
```

### Create a New QoS

```go
newQoS := &interfaces.QoSCreate{
    Name:        "gpu-priority",
    Description: "High priority QoS for GPU jobs",
    Priority:    10000,
    PreemptMode: "REQUEUE",
    Flags:       []string{"NO_DECAY", "PART_TIME_LIMIT"},
    UsageFactor: 2.0,
    MaxJobs:     50,
    MaxWallDuration: 7 * 24 * time.Hour,
    GrpTRES: map[string]int64{
        "gpu": 32,
        "cpu": 1024,
    },
}

err := client.QoS().Create(ctx, newQoS)
if err != nil {
    return err
}

fmt.Println("QoS created successfully")
```

### Update QoS Limits

```go
newMaxJobs := 100
newPriority := 15000

updates := &interfaces.QoSUpdate{
    MaxJobs:  &newMaxJobs,
    Priority: &newPriority,
}

err := client.QoS().Update(ctx, "gpu-priority", updates)
if err != nil {
    return err
}
```

### Create Tiered QoS System

```go
// Create multiple QoS levels
qosLevels := []struct {
    Name        string
    Priority    int
    MaxJobs     int
    MaxWall     time.Duration
    UsageFactor float64
}{
    {"low", 100, 10, 1 * time.Hour, 0.5},
    {"normal", 1000, 50, 12 * time.Hour, 1.0},
    {"high", 5000, 100, 24 * time.Hour, 2.0},
    {"urgent", 10000, 200, 48 * time.Hour, 4.0},
}

for _, level := range qosLevels {
    qos := &interfaces.QoSCreate{
        Name:            level.Name,
        Description:     fmt.Sprintf("%s priority jobs", level.Name),
        Priority:        level.Priority,
        MaxJobs:         level.MaxJobs,
        MaxWallDuration: level.MaxWall,
        UsageFactor:     level.UsageFactor,
    }

    err := client.QoS().Create(ctx, qos)
    if err != nil {
        fmt.Printf("Failed to create QoS %s: %v\n", level.Name, err)
        continue
    }
    fmt.Printf("Created QoS: %s\n", level.Name)
}
```

### Set Preemption Rules

```go
// Configure preemptive QoS
updates := &interfaces.QoSUpdate{
    PreemptMode: &"REQUEUE",
    Flags: []string{
        "PREEMPT",
        "NO_DECAY",
    },
}

err := client.QoS().Update(ctx, "urgent", updates)
if err != nil {
    return err
}
```

### Monitor QoS Usage

```go
qosList, err := client.QoS().List(ctx, nil)
if err != nil {
    return err
}

for _, qos := range qosList.QoS {
    // Get jobs using this QoS
    jobOpts := &interfaces.ListJobsOptions{
        QoS: []string{qos.Name},
        States: []interfaces.JobState{
            interfaces.JobStateRunning,
            interfaces.JobStatePending,
        },
    }

    jobs, err := client.Jobs().List(ctx, jobOpts)
    if err != nil {
        continue
    }

    fmt.Printf("QoS %s: %d jobs (max: %d)\n",
        qos.Name, len(jobs.Jobs), qos.MaxJobs)
}
```

### QoS Resource Limits

```go
// Set comprehensive resource limits
updates := &interfaces.QoSUpdate{
    GrpTRES: map[string]int64{
        "cpu":  10000,
        "mem":  40960000, // 40TB in MB
        "gpu":  64,
        "node": 100,
    },
    MaxTRESPerUser: map[string]int64{
        "cpu": 1000,
        "gpu": 8,
    },
    MaxJobsPerUser: &20,
}

err := client.QoS().Update(ctx, "normal", updates)
if err != nil {
    return err
}
```

### Create Department-Specific QoS

```go
departments := []string{"physics", "chemistry", "biology", "engineering"}

for _, dept := range departments {
    qos := &interfaces.QoSCreate{
        Name:        fmt.Sprintf("%s-qos", dept),
        Description: fmt.Sprintf("QoS for %s department", dept),
        Priority:    1000,
        MaxJobs:     100,
        GrpTRES: map[string]int64{
            "cpu": 2000,
            "gpu": 16,
        },
    }

    err := client.QoS().Create(ctx, qos)
    if err != nil {
        fmt.Printf("Failed to create QoS for %s: %v\n", dept, err)
    }
}
```

### QoS Priority Report

```go
qosList, err := client.QoS().List(ctx, nil)
if err != nil {
    return err
}

// Sort by priority
sort.Slice(qosList.QoS, func(i, j int) bool {
    return qosList.QoS[i].Priority > qosList.QoS[j].Priority
})

fmt.Println("QoS Priority Report:")
fmt.Println("Name            Priority    Usage Factor    Max Jobs")
fmt.Println("----            --------    ------------    --------")

for _, qos := range qosList.QoS {
    fmt.Printf("%-15s %8d    %12.2f    %8d\n",
        qos.Name, qos.Priority, qos.UsageFactor, qos.MaxJobs)
}
```

## Error Handling

```go
qos, err := client.QoS().Get(ctx, "nonexistent")
if err != nil {
    var apiErr *interfaces.APIError
    if errors.As(err, &apiErr) {
        if apiErr.Code == 404 {
            fmt.Println("QoS not found")
        } else if apiErr.Code == 409 {
            fmt.Println("QoS already exists")
        } else {
            fmt.Printf("API error: %s\n", apiErr.Message)
        }
    }
    return err
}
```

## Best Practices

1. **Tiered Structure**: Create multiple QoS levels for different needs
2. **Clear Naming**: Use descriptive names indicating purpose/priority
3. **Document Policies**: Maintain documentation of QoS policies
4. **Regular Review**: Periodically review and adjust QoS settings
5. **Monitor Usage**: Track QoS utilization and adjust limits
6. **Fair Share**: Use usage factors to implement fair sharing

## See Also

- [API Overview](./README.md)
- [Account Management API](./accounts.md)
- [Job Management API](./jobs.md)
- [Examples](../../examples/qos-management)