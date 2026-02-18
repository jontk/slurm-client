# Installation

## Requirements

- Go 1.22 or later
- Access to a SLURM cluster with REST API enabled
- JWT token or other authentication credentials

## Install via go get

```bash
go get github.com/jontk/slurm-client@latest
```

## Install Specific Version

```bash
go get github.com/jontk/slurm-client@v0.3.1
```

## Verify Installation

Create a simple test file to verify the installation:

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
    ctx := context.Background()

    client, err := slurm.NewClient(ctx,
        slurm.WithBaseURL("https://your-slurm-host:6820"),
        slurm.WithAuth(auth.NewTokenAuth("your-token")),
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Connected to SLURM API version: %s\n", client.Version())
}
```

Run the test:

```bash
go run main.go
```

## CLI Installation

The library also includes a command-line interface for interacting with SLURM clusters:

```bash
go install github.com/jontk/slurm-client/cmd/slurm-cli@latest
```

Verify the CLI installation:

```bash
slurm-cli --version
```

## Next Steps

- [Quick Start Guide](quick-start.md) - Learn the basics
- [Configuration](../configuration.md) - Configure authentication and options
- [Examples](../../examples/README.md) - Browse code examples
