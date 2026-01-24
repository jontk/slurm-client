# Reservation Management API

The Reservation Management API provides functionality for managing SLURM resource reservations.

## Interface

```go
type ReservationManager interface {
    // List all reservations with optional filters
    List(ctx context.Context, opts *ListReservationsOptions) (*ReservationList, error)

    // Get a specific reservation by name
    Get(ctx context.Context, reservationName string) (*Reservation, error)

    // Create a new reservation
    Create(ctx context.Context, reservation *ReservationCreate) (*ReservationCreateResponse, error)

    // Update reservation properties
    Update(ctx context.Context, reservationName string, updates *ReservationUpdate) error

    // Delete a reservation
    Delete(ctx context.Context, reservationName string) error
}
```

## Types

### Reservation

```go
type Reservation struct {
    Name         string
    State        ReservationState
    StartTime    time.Time
    EndTime      time.Time
    Duration     time.Duration
    Nodes        string
    NodeCount    int
    CoreCount    int
    Features     []string
    Flags        []string
    Licenses     map[string]int
    MaxStartTime time.Duration
    TRES         string
    Users        []string
    Groups       []string
    Accounts     []string
    BurstBuffer  string
    Watts        int
    Comment      string
}
```

### ReservationState

```go
type ReservationState string

const (
    ReservationStateActive     ReservationState = "ACTIVE"
    ReservationStateInactive   ReservationState = "INACTIVE"
    ReservationStatePending    ReservationState = "PENDING"
    ReservationStateCompleted  ReservationState = "COMPLETED"
)
```

### ReservationCreate

```go
type ReservationCreate struct {
    Name       string
    StartTime  time.Time
    Duration   time.Duration
    Nodes      string
    NodeCount  int
    CoreCount  int
    Partition  string
    Features   []string
    Flags      []string
    Users      []string
    Groups     []string
    Accounts   []string
    Licenses   map[string]int
    Comment    string
}
```

### ReservationUpdate

```go
type ReservationUpdate struct {
    StartTime  *time.Time
    EndTime    *time.Time
    Duration   *time.Duration
    Nodes      *string
    NodeCount  *int
    Users      []string
    Groups     []string
    Accounts   []string
    Comment    *string
}
```

### ListReservationsOptions

```go
type ListReservationsOptions struct {
    // Filter by state
    States []ReservationState

    // Filter by user
    Users []string

    // Filter by account
    Accounts []string

    // Time range filters
    StartAfter  *time.Time
    StartBefore *time.Time
    EndAfter    *time.Time
    EndBefore   *time.Time

    // Pagination
    Limit  int
    Offset int

    // Sorting
    SortBy string
    Order  string
}
```

## Examples

### List All Reservations

```go
reservations, err := client.Reservations().List(ctx, nil)
if err != nil {
    return err
}

for _, res := range reservations.Reservations {
    fmt.Printf("Reservation %s: %s (nodes: %s)\n",
        res.Name, res.State, res.Nodes)
}
```

### Get Reservation Details

```go
reservation, err := client.Reservations().Get(ctx, "weekly-maintenance")
if err != nil {
    return err
}

fmt.Printf("Reservation: %s\n", reservation.Name)
fmt.Printf("State: %s\n", reservation.State)
fmt.Printf("Start: %s\n", reservation.StartTime)
fmt.Printf("End: %s\n", reservation.EndTime)
fmt.Printf("Nodes: %s\n", reservation.Nodes)
```

### Create a Reservation

```go
// Reserve nodes for maintenance
reservation := &interfaces.ReservationCreate{
    Name:      "gpu-training-2024-01-15",
    StartTime: time.Now().Add(24 * time.Hour),
    Duration:  8 * time.Hour,
    Nodes:     "gpu[001-004]",
    Users:     []string{"researcher1", "researcher2"},
    Accounts:  []string{"ml-research"},
    Comment:   "Reserved for deep learning training",
}

response, err := client.Reservations().Create(ctx, reservation)
if err != nil {
    return err
}

fmt.Printf("Created reservation: %s\n", response.Name)
```

### Create a Weekly Recurring Reservation

```go
// Weekly maintenance window
reservation := &interfaces.ReservationCreate{
    Name:      "weekly-maintenance",
    StartTime: nextSunday2AM(),
    Duration:  4 * time.Hour,
    NodeCount: 10,
    Partition: "compute",
    Flags:     []string{"MAINT", "IGNORE_JOBS", "WEEKLY"},
    Users:     []string{"admin"},
    Comment:   "Weekly system maintenance",
}

_, err := client.Reservations().Create(ctx, reservation)
```

### Update a Reservation

```go
// Extend reservation by 2 hours
current, err := client.Reservations().Get(ctx, "gpu-training-2024-01-15")
if err != nil {
    return err
}

newEndTime := current.EndTime.Add(2 * time.Hour)
updates := &interfaces.ReservationUpdate{
    EndTime: &newEndTime,
    Comment: &"Extended for additional experiments",
}

err = client.Reservations().Update(ctx, "gpu-training-2024-01-15", updates)
```

### List Active Reservations

```go
opts := &interfaces.ListReservationsOptions{
    States: []interfaces.ReservationState{
        interfaces.ReservationStateActive,
    },
}

reservations, err := client.Reservations().List(ctx, opts)
if err != nil {
    return err
}

fmt.Printf("Active reservations: %d\n", len(reservations.Reservations))
```

### Find Reservations for User

```go
opts := &interfaces.ListReservationsOptions{
    Users: []string{"researcher1"},
}

reservations, err := client.Reservations().List(ctx, opts)
if err != nil {
    return err
}

for _, res := range reservations.Reservations {
    fmt.Printf("Reservation %s: %s to %s\n",
        res.Name,
        res.StartTime.Format(time.RFC3339),
        res.EndTime.Format(time.RFC3339))
}
```

### Check for Conflicts

```go
// Check if nodes are available for reservation
proposedStart := time.Now().Add(24 * time.Hour)
proposedEnd := proposedStart.Add(4 * time.Hour)

opts := &interfaces.ListReservationsOptions{
    StartBefore: &proposedEnd,
    EndAfter:    &proposedStart,
}

conflicts, err := client.Reservations().List(ctx, opts)
if err != nil {
    return err
}

if len(conflicts.Reservations) > 0 {
    fmt.Println("Conflicting reservations found:")
    for _, res := range conflicts.Reservations {
        fmt.Printf("  - %s: %s\n", res.Name, res.Nodes)
    }
}
```

### Monitor Upcoming Reservations

```go
// Get reservations starting in next 24 hours
now := time.Now()
tomorrow := now.Add(24 * time.Hour)

opts := &interfaces.ListReservationsOptions{
    StartAfter:  &now,
    StartBefore: &tomorrow,
    States: []interfaces.ReservationState{
        interfaces.ReservationStatePending,
    },
}

upcoming, err := client.Reservations().List(ctx, opts)
if err != nil {
    return err
}

fmt.Println("Reservations starting in next 24 hours:")
for _, res := range upcoming.Reservations {
    fmt.Printf("  %s: starts at %s (%s)\n",
        res.Name,
        res.StartTime.Format("15:04"),
        res.Comment)
}
```

### Delete Expired Reservations

```go
// Clean up old reservations
oneWeekAgo := time.Now().Add(-7 * 24 * time.Hour)

opts := &interfaces.ListReservationsOptions{
    EndBefore: &oneWeekAgo,
    States: []interfaces.ReservationState{
        interfaces.ReservationStateCompleted,
    },
}

old, err := client.Reservations().List(ctx, opts)
if err != nil {
    return err
}

for _, res := range old.Reservations {
    err := client.Reservations().Delete(ctx, res.Name)
    if err != nil {
        fmt.Printf("Failed to delete %s: %v\n", res.Name, err)
        continue
    }
    fmt.Printf("Deleted old reservation: %s\n", res.Name)
}
```

## Error Handling

```go
reservation, err := client.Reservations().Get(ctx, "nonexistent")
if err != nil {
    var apiErr *interfaces.APIError
    if errors.As(err, &apiErr) {
        if apiErr.Code == 404 {
            fmt.Println("Reservation not found")
        } else if apiErr.Code == 409 {
            fmt.Println("Reservation conflict")
        } else {
            fmt.Printf("API error: %s\n", apiErr.Message)
        }
    }
    return err
}
```

## Best Practices

1. **Plan Ahead**: Create reservations well in advance
2. **Use Descriptive Names**: Include date/purpose in reservation names
3. **Set Appropriate Flags**: Use MAINT, IGNORE_JOBS as needed
4. **Document Purpose**: Always include meaningful comments
5. **Check Conflicts**: Verify availability before creating reservations
6. **Clean Up**: Delete old reservations to keep the system tidy

## See Also

- [API Overview](./README.md)
- [Node Management API](./nodes.md)
- [Examples](../../examples/reservation-management/main.go)