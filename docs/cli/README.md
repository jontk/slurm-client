# CLI Reference

The SLURM CLI provides a command-line interface for interacting with SLURM clusters via the REST API.

!!! note "Auto-Generated Documentation"
    This section contains auto-generated documentation from the CLI tool itself using Cobra's documentation generator.

## Installation

```bash
go install github.com/jontk/slurm-client/cmd/slurm-cli@latest
```

## Quick Start

```bash
# Set your SLURM REST API URL
export SLURM_REST_URL="https://your-slurm-host:6820"
export SLURM_JWT="your-jwt-token"

# List jobs
slurm-cli jobs list

# Get job details
slurm-cli jobs get 12345

# List nodes
slurm-cli nodes list
```

## Configuration

The CLI can be configured using:

- Command-line flags
- Environment variables
- Configuration file

### Environment Variables

- `SLURM_REST_URL` - SLURM REST API URL (required)
- `SLURM_JWT` - JWT authentication token
- `SLURM_USERNAME` - Basic auth username
- `SLURM_PASSWORD` - Basic auth password

### Configuration File

Create `~/.slurm-cli.yaml`:

```yaml
url: https://your-slurm-host:6820
token: your-jwt-token
api-version: v0.0.43
output: table
```

## Output Formats

The CLI supports multiple output formats:

- `table` - Human-readable table (default)
- `json` - JSON output
- `yaml` - YAML output

```bash
# JSON output
slurm-cli jobs list --output json

# YAML output
slurm-cli nodes list -o yaml
```

## Commands

The complete command reference is auto-generated from the CLI tool. See the sections below for detailed documentation on each command.

---

_The following documentation is automatically generated from the CLI tool._
