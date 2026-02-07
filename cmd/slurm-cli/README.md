# SLURM CLI

A command-line interface for interacting with SLURM clusters via the REST API.

## Installation

```bash
go install github.com/jontk/slurm-client/cmd/slurm-cli@latest
```

Or build from source:

```bash
cd cmd/slurm-cli
go build -o slurm-cli
```

## Configuration

The CLI can be configured using command-line flags or environment variables:

- `--url` or `SLURM_REST_URL`: SLURM REST API URL (required)
- `--token` or `SLURM_JWT`: JWT authentication token
- `--username` or `SLURM_USERNAME`: Basic auth username
- `--password` or `SLURM_PASSWORD`: Basic auth password
- `--api-version`: Specific API version (e.g., v0.0.42)
- `--output`, `-o`: Output format (table, json, yaml)
- `--debug`: Enable debug logging

## Usage

### Jobs Management

List jobs:
```bash
slurm-cli jobs list
slurm-cli jobs list --user 1000 --states RUNNING,PENDING
slurm-cli jobs list --partition gpu --limit 10
```

Get job details:
```bash
slurm-cli jobs get 12345
slurm-cli jobs get 12345 -o json
```

Cancel a job:
```bash
slurm-cli jobs cancel 12345
```

### Submit Jobs

Submit a new job:
```bash
slurm-cli submit --command "python train.py" --name "training-job" \
  --cpus 4 --memory 8192 --time 120 --partition gpu
```

### Nodes Management

List nodes:
```bash
slurm-cli nodes list
slurm-cli nodes list --states IDLE,ALLOCATED
slurm-cli nodes list --partition gpu
```

Get node details:
```bash
slurm-cli nodes get node-001
slurm-cli nodes get node-001 -o json
```

### Partitions

List partitions:
```bash
slurm-cli partitions list
slurm-cli partitions list --states UP
```

### Cluster Information

Get cluster info:
```bash
slurm-cli info
slurm-cli info -o json
```

## Examples

### Using JWT Authentication

```bash
export SLURM_REST_URL="https://cluster.example.com:6820"
export SLURM_JWT="your-jwt-token-here"

slurm-cli jobs list
```

### Using Basic Authentication

```bash
slurm-cli --url https://cluster.example.com:6820 \
  --username admin --password secret \
  jobs list
```

### Specifying API Version

```bash
slurm-cli --api-version v0.0.42 jobs list
```

### JSON Output

```bash
slurm-cli jobs list -o json | jq '.jobs[] | select(.state == "RUNNING")'
```

### Watch for Job State Changes

While the CLI doesn't have built-in watch functionality, you can use it with standard Unix tools:

```bash
# Watch job status every 5 seconds
watch -n 5 'slurm-cli jobs get 12345'

# Monitor all running jobs
watch -n 10 'slurm-cli jobs list --states RUNNING'
```

## Error Handling

The CLI provides detailed error messages:

```bash
$ slurm-cli jobs get 99999
Error: job not found: Job ID 99999 not found

$ slurm-cli --url https://invalid.example.com jobs list
Error: connection refused: dial tcp: lookup invalid.example.com: no such host
```

## Exit Codes

- 0: Success
- 1: General error (configuration, network, etc.)
- 2: Resource not found
- 3: Authentication error
- 4: Invalid input/validation error