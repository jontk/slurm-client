# Account Management API

The Account Management API provides functionality for managing SLURM accounts and associations.

## Interface

```go
type AccountManager interface {
    // List all accounts with optional filters
    List(ctx context.Context, opts *ListAccountsOptions) (*AccountList, error)

    // Get a specific account by name
    Get(ctx context.Context, accountName string) (*Account, error)

    // Create a new account
    Create(ctx context.Context, account *AccountCreate) error

    // Update account properties
    Update(ctx context.Context, accountName string, updates *AccountUpdate) error

    // Delete an account
    Delete(ctx context.Context, accountName string) error

    // Manage associations
    ListAssociations(ctx context.Context, accountName string) (*AssociationList, error)
    AddAssociation(ctx context.Context, accountName string, assoc *AssociationCreate) error
    RemoveAssociation(ctx context.Context, accountName string, assocID string) error
}
```

## Types

### Account

```go
type Account struct {
    Name             string
    Description      string
    Organization     string
    Coordinators     []string
    AllowedPartitions []string
    DefaultQoS       string
    AllowedQoS       []string
    GrpTRES          map[string]int64
    GrpJobs          int
    GrpSubmitJobs    int
    GrpWall          time.Duration
    MaxTRES          map[string]int64
    MaxJobs          int
    MaxSubmitJobs    int
    MaxWallDuration  time.Duration
    Priority         int
    FairShare        int
}
```

### Association

```go
type Association struct {
    ID               string
    Account          string
    User             string
    Partition        string
    Priority         int
    GrpTRES          map[string]int64
    GrpJobs          int
    GrpSubmitJobs    int
    GrpWall          time.Duration
    MaxTRES          map[string]int64
    MaxJobs          int
    MaxSubmitJobs    int
    MaxWallDuration  time.Duration
    QoS              []string
    DefaultQoS       string
}
```

### AccountCreate

```go
type AccountCreate struct {
    Name             string
    Description      string
    Organization     string
    Coordinators     []string
    AllowedPartitions []string
    DefaultQoS       string
    GrpJobs          int
    MaxJobs          int
    Priority         int
}
```

### AccountUpdate

```go
type AccountUpdate struct {
    Description      *string
    Organization     *string
    Coordinators     []string
    AllowedPartitions []string
    DefaultQoS       *string
    GrpJobs          *int
    MaxJobs          *int
    Priority         *int
}
```

### AssociationCreate

```go
type AssociationCreate struct {
    User            string
    Partition       string
    Priority        int
    GrpJobs         int
    MaxJobs         int
    MaxWallDuration time.Duration
    QoS             []string
}
```

## Examples

### List All Accounts

```go
accounts, err := client.Accounts().List(ctx, nil)
if err != nil {
    return err
}

for _, account := range accounts.Accounts {
    fmt.Printf("Account: %s - %s\n", account.Name, account.Description)
}
```

### Get Account Details

```go
account, err := client.Accounts().Get(ctx, "research")
if err != nil {
    return err
}

fmt.Printf("Account: %s\n", account.Name)
fmt.Printf("Organization: %s\n", account.Organization)
fmt.Printf("Priority: %d\n", account.Priority)
fmt.Printf("Max Jobs: %d\n", account.MaxJobs)
```

### Create a New Account

```go
newAccount := &interfaces.AccountCreate{
    Name:         "ml-research",
    Description:  "Machine Learning Research Group",
    Organization: "Computer Science Department",
    Coordinators: []string{"prof-smith", "admin-jones"},
    AllowedPartitions: []string{"gpu-compute", "cpu-compute"},
    DefaultQoS:   "normal",
    MaxJobs:      100,
    Priority:     1000,
}

err := client.Accounts().Create(ctx, newAccount)
if err != nil {
    return err
}

fmt.Println("Account created successfully")
```

### Update Account Limits

```go
newMaxJobs := 200
newPriority := 1500

updates := &interfaces.AccountUpdate{
    MaxJobs:  &newMaxJobs,
    Priority: &newPriority,
}

err := client.Accounts().Update(ctx, "ml-research", updates)
if err != nil {
    return err
}
```

### Add User Association

```go
association := &interfaces.AssociationCreate{
    User:            "researcher1",
    Partition:       "gpu-compute",
    Priority:        100,
    MaxJobs:         20,
    MaxWallDuration: 24 * time.Hour,
    QoS:             []string{"normal", "high"},
}

err := client.Accounts().AddAssociation(ctx, "ml-research", association)
if err != nil {
    return err
}

fmt.Printf("Added user %s to account\n", association.User)
```

### List Account Associations

```go
associations, err := client.Accounts().ListAssociations(ctx, "ml-research")
if err != nil {
    return err
}

fmt.Printf("Users in account ml-research:\n")
for _, assoc := range associations.Associations {
    fmt.Printf("  - %s (max jobs: %d)\n", assoc.User, assoc.MaxJobs)
}
```

### Monitor Account Usage

```go
// Get account details
account, err := client.Accounts().Get(ctx, "research")
if err != nil {
    return err
}

// Get running jobs for account
jobOpts := &interfaces.ListJobsOptions{
    Accounts: []string{account.Name},
    States:   []interfaces.JobState{interfaces.JobStateRunning},
}

jobs, err := client.Jobs().List(ctx, jobOpts)
if err != nil {
    return err
}

fmt.Printf("Account %s usage:\n", account.Name)
fmt.Printf("  Running jobs: %d / %d\n", len(jobs.Jobs), account.MaxJobs)
fmt.Printf("  Usage: %.1f%%\n", float64(len(jobs.Jobs))/float64(account.MaxJobs)*100)
```

### Account Hierarchy Report

```go
accounts, err := client.Accounts().List(ctx, nil)
if err != nil {
    return err
}

for _, account := range accounts.Accounts {
    associations, err := client.Accounts().ListAssociations(ctx, account.Name)
    if err != nil {
        continue
    }

    fmt.Printf("\nAccount: %s\n", account.Name)
    fmt.Printf("  Description: %s\n", account.Description)
    fmt.Printf("  Users: %d\n", len(associations.Associations))
    fmt.Printf("  Priority: %d\n", account.Priority)
    fmt.Printf("  Max Jobs: %d\n", account.MaxJobs)
}
```

### Manage Account Coordinators

```go
// Add coordinators
updates := &interfaces.AccountUpdate{
    Coordinators: []string{"prof-smith", "dr-jones", "admin-wilson"},
}

err := client.Accounts().Update(ctx, "research", updates)
if err != nil {
    return err
}
```

### Set Account Resource Limits

```go
updates := &interfaces.AccountUpdate{
    GrpTRES: map[string]int64{
        "cpu": 1000,
        "mem": 4096000, // MB
        "gpu": 8,
    },
    GrpJobs:       50,
    GrpSubmitJobs: 100,
}

err := client.Accounts().Update(ctx, "ml-research", updates)
if err != nil {
    return err
}
```

## Error Handling

```go
account, err := client.Accounts().Get(ctx, "nonexistent")
if err != nil {
    var apiErr *interfaces.APIError
    if errors.As(err, &apiErr) {
        if apiErr.Code == 404 {
            fmt.Println("Account not found")
        } else if apiErr.Code == 409 {
            fmt.Println("Account already exists")
        } else {
            fmt.Printf("API error: %s\n", apiErr.Message)
        }
    }
    return err
}
```

## Best Practices

1. **Hierarchical Structure**: Organize accounts by department/project
2. **Set Appropriate Limits**: Configure resource limits based on needs
3. **Regular Audits**: Review account usage and associations
4. **Use Descriptive Names**: Include organization/purpose in names
5. **Document Coordinators**: Maintain up-to-date coordinator lists
6. **Monitor Usage**: Track resource consumption by account

## See Also

- [API Overview](./README.md)
- [QoS Management API](./qos.md)
- [Job Management API](./jobs.md)
- [Examples](../../examples/account-management/main.go)