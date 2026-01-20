# Partition Management API

The Partition Management API provides functionality for managing SLURM partitions (queues).

## Interface

```go
type PartitionManager interface {
    // List all partitions with optional filters
    List(ctx context.Context, opts *ListPartitionsOptions) (*PartitionList, error)

    // Get a specific partition by name
    Get(ctx context.Context, partitionName string) (*Partition, error)

    // Create a new partition
    Create(ctx context.Context, partition *PartitionCreate) error

    // Update partition properties
    Update(ctx context.Context, partitionName string, updates *PartitionUpdate) error

    // Delete a partition
    Delete(ctx context.Context, partitionName string) error
}
```

## Types

### Partition

```go
type Partition struct {
    Name               string
    State              PartitionState
    TotalNodes         int
    TotalCPUs          int
    DefaultTime        time.Duration
    MaxTime            time.Duration
    DefaultMemPerCPU   int
    MaxMemPerCPU       int
    DefaultMemPerNode  int
    MaxMemPerNode      int
    Priority           int
    MaxNodes           int
    MinNodes           int
    Nodes              string
    AllowGroups        []string
    AllowAccounts      []string
    AllowQoS           []string
    DenyAccounts       []string
    DenyQoS            []string
    PreemptMode        string
    GraceTime          int
    DefCPUPerGPU       int
    DefMemPerGPU       int
    MaxCPUsPerNode     int
    QoS                string
    TRESBillingWeights map[string]float64
}
```

### PartitionState

```go
type PartitionState string

const (
    PartitionStateUp       PartitionState = "UP"
    PartitionStateDown     PartitionState = "DOWN"
    PartitionStateDrain    PartitionState = "DRAIN"
    PartitionStateInactive PartitionState = "INACTIVE"
)
```

### PartitionCreate

```go
type PartitionCreate struct {
    Name             string
    Nodes            string
    DefaultTime      time.Duration
    MaxTime          time.Duration
    State            PartitionState
    MaxNodes         int
    MinNodes         int
    Priority         int
    AllowGroups      []string
    AllowAccounts    []string
    QoS              string
}
```

### PartitionUpdate

```go
type PartitionUpdate struct {
    State            *PartitionState
    Nodes            *string
    DefaultTime      *time.Duration
    MaxTime          *time.Duration
    Priority         *int
    MaxNodes         *int
    MinNodes         *int
    AllowGroups      []string
    AllowAccounts    []string
}
```

### ListPartitionsOptions

```go
type ListPartitionsOptions struct {
    // Filter by state
    States []PartitionState

    // Filter by QoS
    QoS []string

    // Pagination
    Limit  int
    Offset int

    // Sorting
    SortBy string
    Order  string
}
```

## Examples

### List All Partitions

```go
partitions, err := client.Partitions().List(ctx, nil)
if err != nil {
    return err
}

for _, partition := range partitions.Partitions {
    fmt.Printf("Partition %s: %s (%d nodes)\n",
        partition.Name, partition.State, partition.TotalNodes)
}
```

### Get Partition Details

```go
partition, err := client.Partitions().Get(ctx, "compute")
if err != nil {
    return err
}

fmt.Printf("Partition: %s\n", partition.Name)
fmt.Printf("State: %s\n", partition.State)
fmt.Printf("Total Nodes: %d\n", partition.TotalNodes)
fmt.Printf("Total CPUs: %d\n", partition.TotalCPUs)
fmt.Printf("Max Time: %s\n", partition.MaxTime)
```

### Create a New Partition

```go
newPartition := &interfaces.PartitionCreate{
    Name:          "gpu-compute",
    Nodes:         "gpu[001-010]",
    DefaultTime:   1 * time.Hour,
    MaxTime:       24 * time.Hour,
    State:         interfaces.PartitionStateUp,
    Priority:      100,
    AllowAccounts: []string{"research", "engineering"},
}

err := client.Partitions().Create(ctx, newPartition)
if err != nil {
    return err
}

fmt.Println("Partition created successfully")
```

### Update Partition Configuration

```go
updates := &interfaces.PartitionUpdate{
    MaxTime:  48 * time.Hour,
    Priority: 200,
}

err := client.Partitions().Update(ctx, "gpu-compute", updates)
if err != nil {
    return err
}
```

### Drain a Partition

```go
drain := interfaces.PartitionStateDrain
updates := &interfaces.PartitionUpdate{
    State: &drain,
}

err := client.Partitions().Update(ctx, "compute", updates)
if err != nil {
    return err
}

fmt.Println("Partition set to drain state")
```

### Monitor Partition Usage

```go
partitions, err := client.Partitions().List(ctx, nil)
if err != nil {
    return err
}

fmt.Println("Partition Usage Summary:")
for _, p := range partitions.Partitions {
    // Get jobs in partition
    jobOpts := &interfaces.ListJobsOptions{
        Partitions: []string{p.Name},
        States:     []interfaces.JobState{interfaces.JobStateRunning},
    }

    jobs, err := client.Jobs().List(ctx, jobOpts)
    if err != nil {
        continue
    }

    fmt.Printf("  %s: %d running jobs\n", p.Name, len(jobs.Jobs))
}
```

### Find Partitions by Account

```go
partitions, err := client.Partitions().List(ctx, nil)
if err != nil {
    return err
}

account := "research"
fmt.Printf("Partitions available to account '%s':\n", account)

for _, p := range partitions.Partitions {
    // Check if account is allowed
    for _, allowed := range p.AllowAccounts {
        if allowed == account {
            fmt.Printf("  - %s\n", p.Name)
            break
        }
    }
}
```

### Partition Capacity Check

```go
partition, err := client.Partitions().Get(ctx, "compute")
if err != nil {
    return err
}

// Get nodes in partition
nodeOpts := &interfaces.ListNodesOptions{
    Partitions: []string{partition.Name},
}

nodes, err := client.Nodes().List(ctx, nodeOpts)
if err != nil {
    return err
}

var idleNodes, allocNodes int
for _, node := range nodes.Nodes {
    switch node.State {
    case interfaces.NodeStateIdle:
        idleNodes++
    case interfaces.NodeStateAllocated, interfaces.NodeStateMixed:
        allocNodes++
    }
}

fmt.Printf("Partition %s capacity:\n", partition.Name)
fmt.Printf("  Total nodes: %d\n", partition.TotalNodes)
fmt.Printf("  Idle nodes: %d\n", idleNodes)
fmt.Printf("  Allocated nodes: %d\n", allocNodes)
fmt.Printf("  Utilization: %.1f%%\n",
    float64(allocNodes)/float64(partition.TotalNodes)*100)
```

## Error Handling

```go
partition, err := client.Partitions().Get(ctx, "nonexistent")
if err != nil {
    var apiErr *interfaces.APIError
    if errors.As(err, &apiErr) {
        if apiErr.Code == 404 {
            fmt.Println("Partition not found")
        } else {
            fmt.Printf("API error: %s\n", apiErr.Message)
        }
    }
    return err
}
```

## Best Practices

1. **Plan Partition Structure**: Design partitions based on workload requirements
2. **Set Appropriate Limits**: Configure time limits and node counts carefully
3. **Use QoS**: Implement Quality of Service for better resource management
4. **Monitor Usage**: Regularly check partition utilization
5. **Document Purpose**: Use clear naming and maintain documentation

## See Also

- [API Overview](./README.md)
- [Job Management API](./jobs.md)
- [Node Management API](./nodes.md)
- [Resource Allocation Example](../../examples/resource-allocation) - Working with partitions for resource allocation
- [Basic Examples](../../examples/basic) - Getting started with partitions