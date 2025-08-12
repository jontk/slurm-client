# Code Generation Guide

This document explains the code generation process used in the SLURM REST API Client Library.

## Overview

The library uses OpenAPI specifications to generate type-safe Go clients for each supported SLURM API version. This ensures consistency and reduces manual coding errors.

## Architecture

```
api/openapi/                    # OpenAPI specifications
    ├── v0.0.40/
    │   └── openapi.yaml
    ├── v0.0.41/
    │   └── openapi.yaml
    ├── v0.0.42/
    │   └── openapi.yaml
    └── v0.0.43/
        └── openapi.yaml

internal/api/                   # Generated code output
    ├── v0_0_40/
    │   ├── client.gen.go      # Generated client
    │   └── types.gen.go       # Generated types
    ├── v0_0_41/
    ├── v0_0_42/
    └── v0_0_43/
```

## Code Generation Process

### 1. Prerequisites

Install the OpenAPI generator tool:
```bash
go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
```

### 2. Generate Code

Run the generation command:
```bash
go generate ./...
```

This executes `//go:generate` directives in the source files.

### 3. Generation Configuration

Each API version has a generation configuration in `internal/api/v0_0_XX/generate.go`:

```go
//go:generate oapi-codegen -package api -generate types,client -o client.gen.go ../../../api/openapi/v0.0.43/openapi.yaml
```

## Adding a New API Version

### Step 1: Add OpenAPI Specification

Place the new OpenAPI spec in `api/openapi/v0.0.XX/openapi.yaml`:

```yaml
openapi: 3.0.0
info:
  title: Slurm REST API
  version: v0.0.XX
  description: SLURM REST API v0.0.XX
paths:
  /slurm/v0.0.XX/jobs:
    get:
      summary: List jobs
      # ... rest of specification
```

### Step 2: Create Generation Directory

Create the directory structure:
```bash
mkdir -p internal/api/v0_0_XX
```

### Step 3: Add Generation Configuration

Create `internal/api/v0_0_XX/generate.go`:

```go
// SPDX-FileCopyrightText: 2025 SLURM REST API Client Contributors
// SPDX-License-Identifier: Apache-2.0

package api

//go:generate oapi-codegen -package api -generate types,client -o client.gen.go ../../../api/openapi/v0.0.XX/openapi.yaml
```

### Step 4: Generate Code

Run the generation:
```bash
go generate ./internal/api/v0_0_XX
```

### Step 5: Implement Adapters

Create adapters in `internal/adapters/v0_0_XX/`:

```go
// job_adapter.go
type JobAdapter struct {
    client *api.Client
}

func (a *JobAdapter) List(ctx context.Context, opts *types.JobListOptions) (*types.JobList, error) {
    // Convert options to API format
    // Call generated client
    // Convert response to common types
}
```

### Step 6: Implement Managers

Create managers in `internal/api/v0_0_XX/`:

```go
// job_manager_impl.go
type JobManagerImpl struct {
    adapter adapters.JobAdapter
}

func (m *JobManagerImpl) List(ctx context.Context, opts *interfaces.ListJobsOptions) (*interfaces.JobList, error) {
    // Convert interface types to adapter types
    // Call adapter
    // Convert response to interface types
}
```

### Step 7: Update Factory

Add the new version to `pkg/client/factory/factory.go`:

```go
func NewClient(config *interfaces.ClientConfig) (interfaces.Client, error) {
    switch config.Version {
    case "v0.0.XX":
        return v0_0_XX.NewClient(config)
    // ... other versions
    }
}
```

## Customizing Generation

### Custom Templates

To use custom templates:

```bash
oapi-codegen -package api \
  -generate types,client \
  -templates ./templates \
  -o client.gen.go \
  openapi.yaml
```

### Generation Options

Available generation targets:
- `types`: Generate type definitions
- `client`: Generate client code
- `server`: Generate server stubs (not used)
- `spec`: Generate embedded spec

### Type Mappings

Custom type mappings in `oapi-codegen.yaml`:

```yaml
package: api
generate:
  - types
  - client
output: client.gen.go
type-mappings:
  time.Time: time.Time
  UUID: github.com/google/uuid.UUID
```

## Post-Generation Tasks

### 1. Fix Imports

Ensure all imports are properly formatted:
```bash
goimports -w internal/api/v0_0_XX/
```

### 2. Run Tests

Verify generated code compiles:
```bash
go test ./internal/api/v0_0_XX/...
```

### 3. Update Documentation

Add the new version to:
- README.md version support table
- Configuration documentation
- Architecture documentation

## Handling Generation Issues

### Common Problems

1. **Missing Types**: Some OpenAPI specs may have incomplete type definitions
   - Solution: Add type patches before generation

2. **Naming Conflicts**: Generated names may conflict
   - Solution: Use type mappings or post-generation scripts

3. **Invalid Go Code**: Some OpenAPI constructs don't map well to Go
   - Solution: Post-process generated files

### Type Patches

For incomplete specs, create patches:

```go
// patches/v0_0_XX.go
type PatchedJobSubmission struct {
    api.JobSubmission
    MissingField string `json:"missing_field,omitempty"`
}
```

## Regeneration Workflow

When updating OpenAPI specs:

1. **Update Specification**:
   ```bash
   # Edit the OpenAPI spec
   vim api/openapi/v0.0.XX/openapi.yaml
   ```

2. **Regenerate Code**:
   ```bash
   go generate ./internal/api/v0_0_XX
   ```

3. **Check Differences**:
   ```bash
   git diff internal/api/v0_0_XX/
   ```

4. **Update Adapters**: Modify adapters if API changes require it

5. **Run Tests**: Ensure all tests pass

6. **Update Version**: If breaking changes, update version compatibility

## Best Practices

1. **Never Edit Generated Files**: They will be overwritten
2. **Keep Specs Updated**: Regularly sync with official SLURM specs
3. **Test Thoroughly**: Generated code may have edge cases
4. **Document Changes**: Note any manual interventions needed
5. **Version Control**: Commit generated files for easier reviews

## Tools and Scripts

### Validation Script

```bash
#!/bin/bash
# validate-openapi.sh
for spec in api/openapi/*/openapi.yaml; do
    echo "Validating $spec"
    oapi-codegen -package test -generate types $spec > /dev/null
    if [ $? -ne 0 ]; then
        echo "ERROR: Invalid spec $spec"
        exit 1
    fi
done
```

### Regeneration Script

```bash
#!/bin/bash
# regenerate-all.sh
echo "Regenerating all API clients..."
go generate ./...
echo "Running formatter..."
gofmt -w internal/api/
echo "Running imports..."
goimports -w internal/api/
echo "Running tests..."
go test ./...
```

## See Also

- [Architecture Overview](./ARCHITECTURE.md)
- [Contributing Guide](../CONTRIBUTING.md)
- [OpenAPI Specification](https://www.openapis.org/)