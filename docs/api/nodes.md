# Node Management API

The Node Management API provides functionality for managing and monitoring SLURM compute nodes.

## Interface

```go
type NodeManager interface {
    // List all nodes with optional filters
    List(ctx context.Context, opts *ListNodesOptions) (*NodeList, error)

    // Get a specific node by name
    Get(ctx context.Context, nodeName string) (*Node, error)

    // Update node properties
    Update(ctx context.Context, nodeName string, updates *NodeUpdate) error

    // Drain a node (mark unavailable for new jobs)
    Drain(ctx context.Context, nodeName string, reason string) error

    // Resume a node (mark available for jobs)
    Resume(ctx context.Context, nodeName string) error
}
```

## Types

### Node

```go
type Node struct {
    Name            string
    State           NodeState
    Partitions      []string
    CPUs            int
    RealMemory      int
    TmpDisk         int
    Features        []string
    Gres            string
    NodeAddr        string
    NodeHostname    string
    Version         string
    OS              string
    KernelVersion   string
    Architecture    string
    CPULoad         float64
    FreeMem         int
    AllocCPUs       int
    IdleCPUs        int
    AllocMemory     int
    Reason          string
    ReasonTime      *time.Time
    ReasonUID       int
    BootTime        time.Time
    SlurmdStartTime time.Time
}
```

### NodeState

```go
type NodeState string

const (
    NodeStateUnknown     NodeState = "UNKNOWN"
    NodeStateDown        NodeState = "DOWN"
    NodeStateIdle        NodeState = "IDLE"
    NodeStateAllocated   NodeState = "ALLOCATED"
    NodeStateMixed       NodeState = "MIXED"
    NodeStateDrain       NodeState = "DRAIN"
    NodeStateDraining    NodeState = "DRAINING"
    NodeStateFail        NodeState = "FAIL"
    NodeStateFailing     NodeState = "FAILING"
    NodeStateFuture      NodeState = "FUTURE"
    NodeStateMaintenance NodeState = "MAINT"
    NodeStatePoweredDown NodeState = "POWERED_DOWN"
    NodeStatePoweringUp  NodeState = "POWERING_UP"
    NodeStateReboot      NodeState = "REBOOT"
)
```

### NodeUpdate

```go
type NodeUpdate struct {
    State       *NodeState
    Reason      *string
    Features    []string
    Weight      *int
    Gres        *string
    Comment     *string
}
```

### ListNodesOptions

```go
type ListNodesOptions struct {
    // Filter by state
    States []NodeState

    // Filter by partition
    Partitions []string

    // Filter by features
    Features []string

    // Filter by reason
    Reason string

    // Pagination
    Limit  int
    Offset int

    // Sorting
    SortBy string
    Order  string
}
```

## Examples

### List All Nodes

```go
nodes, err := client.Nodes().List(ctx, nil)
if err != nil {
    return err
}

for _, node := range nodes.Nodes {
    fmt.Printf("Node %s: %s (%d/%d CPUs)\n",
        node.Name, node.State, node.AllocCPUs, node.CPUs)
}
```

### List Idle Nodes

```go
opts := &interfaces.ListNodesOptions{
    States: []interfaces.NodeState{interfaces.NodeStateIdle},
}

nodes, err := client.Nodes().List(ctx, opts)
if err != nil {
    return err
}

fmt.Printf("Found %d idle nodes\n", len(nodes.Nodes))
```

### Get Node Details

```go
node, err := client.Nodes().Get(ctx, "compute001")
if err != nil {
    return err
}

fmt.Printf("Node: %s\n", node.Name)
fmt.Printf("State: %s\n", node.State)
fmt.Printf("CPUs: %d (allocated: %d)\n", node.CPUs, node.AllocCPUs)
fmt.Printf("Memory: %d MB (allocated: %d MB)\n", node.RealMemory, node.AllocMemory)
```

### Drain a Node for Maintenance

```go
err := client.Nodes().Drain(ctx, "compute001", "Hardware maintenance scheduled")
if err != nil {
    return err
}

fmt.Println("Node drained successfully")
```

### Resume a Node After Maintenance

```go
err := client.Nodes().Resume(ctx, "compute001")
if err != nil {
    return err
}

fmt.Println("Node resumed successfully")
```

### Update Node Features

```go
updates := &interfaces.NodeUpdate{
    Features: []string{"gpu", "infiniband", "nvme"},
}

err := client.Nodes().Update(ctx, "compute001", updates)
if err != nil {
    return err
}
```

### Monitor Node Health

```go
nodes, err := client.Nodes().List(ctx, nil)
if err != nil {
    return err
}

var down, drain, idle, allocated int

for _, node := range nodes.Nodes {
    switch node.State {
    case interfaces.NodeStateDown:
        down++
    case interfaces.NodeStateDrain, interfaces.NodeStateDraining:
        drain++
    case interfaces.NodeStateIdle:
        idle++
    case interfaces.NodeStateAllocated, interfaces.NodeStateMixed:
        allocated++
    }
}

fmt.Printf("Node Summary:\n")
fmt.Printf("  Down: %d\n", down)
fmt.Printf("  Drain: %d\n", drain)
fmt.Printf("  Idle: %d\n", idle)
fmt.Printf("  Allocated: %d\n", allocated)
```

### Find Nodes by Feature

```go
opts := &interfaces.ListNodesOptions{
    Features: []string{"gpu"},
}

nodes, err := client.Nodes().List(ctx, opts)
if err != nil {
    return err
}

fmt.Printf("Found %d GPU nodes\n", len(nodes.Nodes))
```

### Batch Node Operations

```go
// Drain multiple nodes
nodesToDrain := []string{"compute001", "compute002", "compute003"}

for _, nodeName := range nodesToDrain {
    err := client.Nodes().Drain(ctx, nodeName, "Batch maintenance")
    if err != nil {
        fmt.Printf("Failed to drain %s: %v\n", nodeName, err)
        continue
    }
    fmt.Printf("Drained %s\n", nodeName)
}
```

## Error Handling

```go
node, err := client.Nodes().Get(ctx, "compute001")
if err != nil {
    var apiErr *interfaces.APIError
    if errors.As(err, &apiErr) {
        if apiErr.Code == 404 {
            fmt.Println("Node not found")
        } else {
            fmt.Printf("API error: %s\n", apiErr.Message)
        }
    }
    return err
}
```

## Best Practices

1. **Drain Before Maintenance**: Always drain nodes before performing maintenance
2. **Check Node State**: Verify node state before operations
3. **Use Meaningful Reasons**: Provide clear reasons when draining nodes
4. **Monitor Node Health**: Regularly check for down or failing nodes
5. **Batch Operations**: Use goroutines for parallel operations on multiple nodes

## See Also

- [API Overview](./README.md)
- [Job Management API](./jobs.md)
- [Partition Management API](./partitions.md)
- [Node Monitoring Example](../../examples/watch-nodes/main.go) - Real-time node monitoring
- [Basic Examples](../../examples/basic/README.md) - Getting started with nodes