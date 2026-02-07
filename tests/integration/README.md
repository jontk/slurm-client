# Integration Test Framework

This directory contains comprehensive integration tests for the slurm-client library, including both mock server tests and real SLURM server tests.

## Test Structure

### Mock Server Infrastructure
- **`../mocks/server.go`** - Full-featured mock SLURM REST API server
- **`../mocks/handlers.go`** - HTTP handlers for all API endpoints 
- **`../mocks/versions.go`** - Version-specific configurations and behaviors

### Test Suites
- **`job_lifecycle_test.go`** - End-to-end job workflow testing (submit → monitor → cancel)
- **`multi_version_test.go`** - Cross-version compatibility and migration testing
- **`auth_test.go`** - Authentication provider validation across versions
- **`error_handling_test.go`** - Structured error handling and classification testing
- **`real_server_test.go`** - Tests against real SLURM REST API servers

## Mock Server Features

### Multi-Version Support
- **v0.0.40** - Basic operations, limited update support
- **v0.0.41** - Enhanced job updates, field renames (minimum_switches → required_switches)  
- **v0.0.42** - Full CRUD operations, removed exclusive/oversubscribe fields
- **v0.0.43** - Latest features, FrontEnd mode removal

### Comprehensive API Coverage
- **Job Management** - List, Get, Submit, Cancel, Update, Steps
- **Node Management** - List, Get, Update (version-dependent)
- **Partition Management** - List, Get, Update (version-dependent)
- **Cluster Info** - Get, Ping, Stats, Version

### Advanced Testing Features
- **Authentication Testing** - Token, Basic, and No-auth scenarios
- **Error Simulation** - Configurable error responses for testing error handling
- **Response Delays** - Simulate network latency and timeouts
- **Pagination** - Full pagination support with offset/limit
- **Filtering** - Advanced filtering by user, state, partition, etc.

## Current Status

### ✅ Completed Infrastructure
- Mock server with 66,010+ lines of OpenAPI spec coverage
- Version-specific behavior simulation
- Comprehensive error response handling
- Authentication middleware with token validation
- Structured error classification system

### 🔧 Ready for Implementation
The test framework is complete and ready to validate client implementations. Current test failures are expected because:

1. **Client methods not implemented** - Most manager methods return `nil, nil`
2. **Only JobManager.List() working** - As documented in project status
3. **Tests validate readiness** - Framework successfully detects missing implementations

### 🎯 Next Steps for Client Implementation
1. **Complete JobManager methods** - Get, Submit, Cancel, Update, Steps
2. **Implement other managers** - NodeManager, PartitionManager, InfoManager  
3. **Version-specific testing** - Validate breaking changes and compatibility
4. **Performance optimization** - Connection pooling and caching

## Running Tests

### Mock Server Tests

```bash
# Run all mock server integration tests
go test -v ./tests/integration/... -run "^Test[^R][^e][^a][^l]"

# Run specific test suites
go test -v ./tests/integration/job_lifecycle_test.go
go test -v ./tests/integration/multi_version_test.go  
go test -v ./tests/integration/auth_test.go
go test -v ./tests/integration/error_handling_test.go

# Test with timeout for network scenarios
go test -v ./tests/integration/... -timeout 60s
```

### Real Server Tests

```bash
# Use the test script (recommended)
export SLURM_REAL_SERVER_TEST=true
./scripts/test-real-server.sh

# Or run directly with configuration
export SLURM_REAL_SERVER_TEST=true
export SLURM_SERVER_URL="http://localhost
export SLURM_API_VERSION="v0.0.42"
go test -v ./tests/integration -run TestRealServer
```

See the [Real Server Testing](#real-server-testing) section below for detailed configuration.

## Test Coverage

### Authentication Scenarios
- Valid/invalid token authentication
- Basic authentication  
- No authentication
- Cross-version auth compatibility
- Concurrent authenticated requests

### Error Handling Scenarios  
- HTTP status code mapping (400, 401, 404, 500, 503, etc.)
- Network errors (timeout, connection refused, DNS failures)
- Context cancellation and deadline exceeded
- Version-specific error responses
- Structured error classification

### Multi-Version Scenarios
- Same operations across all API versions
- Version-specific feature availability
- Migration between versions
- Concurrent version usage
- Breaking change validation

### Job Lifecycle Scenarios
- Complete workflow testing
- Validation of job submission parameters
- State transitions monitoring
- Error handling during operations
- Filtering and pagination

## Architecture Benefits

This integration test framework provides:

1. **Confidence in Implementation** - Comprehensive validation of all client features
2. **Regression Prevention** - Catches breaking changes across versions
3. **Performance Validation** - Tests network behavior and error recovery
4. **Documentation** - Tests serve as executable documentation
5. **Development Velocity** - Fast feedback during implementation

The framework is production-ready and will provide excellent validation once the client implementation is completed.

## Real Server Testing

The `real_server_test.go` provides comprehensive testing against real SLURM REST API servers.

### Current Status
⚠️ **Note**: As of now, only the Ping test will pass against a real server because most client manager methods are not yet implemented (they return `nil, nil`). Only `JobManager.List()` is fully implemented. This test suite will validate the implementations as they are completed.

### Prerequisites

1. Access to a SLURM REST API server (e.g., `localhost
2. SSH access to the server for JWT token generation
3. Network connectivity to the API endpoint

### Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `SLURM_REAL_SERVER_TEST` | `false` | Must be `true` to enable real server tests |
| `SLURM_SERVER_URL` | `http://localhost | SLURM REST API endpoint |
| `SLURM_API_VERSION` | `v0.0.43` | API version (v0.0.43 is latest for SLURM 25.05) |
| `SLURM_SSH_HOST` | `localhost | SSH host for token generation |
| `SLURM_SSH_USER` | `root` | SSH user for token generation |
| `SLURM_JWT_TOKEN` | (auto-fetched) | JWT token (optional, fetched via SSH if not provided) |

### Running Real Server Tests

```bash
# Quick start with defaults
export SLURM_REAL_SERVER_TEST=true
./scripts/test-real-server.sh

# Custom configuration
export SLURM_REAL_SERVER_TEST=true
export SLURM_SERVER_URL="http://your-slurm-server:6820"
export SLURM_API_VERSION="v0.0.40"
export SLURM_SSH_HOST="your-slurm-server"
export SLURM_SSH_USER="slurm-admin"
./scripts/test-real-server.sh

# Provide token directly (skip SSH)
export SLURM_REAL_SERVER_TEST=true
export SLURM_JWT_TOKEN="your-jwt-token-here"
go test -v ./tests/integration -run TestRealServer
```

### Test Coverage

Real server tests validate:
- **Authentication**: JWT token validation
- **Connectivity**: Ping and version endpoints
- **Job Operations**: Submit, list, get, cancel real jobs
- **Node Operations**: List and inspect cluster nodes
- **Partition Operations**: List and inspect partitions
- **Statistics**: Cluster utilization and metrics

### Security Considerations

1. **SSH Keys**: Use key-based authentication for automated testing
2. **Token Security**: Never commit JWT tokens to version control
3. **Test Isolation**: Use dedicated test partitions to avoid disrupting production
4. **Resource Cleanup**: Tests attempt to cancel submitted jobs

### Troubleshooting

```bash
# Test SSH connectivity
ssh -v root@localhost echo "OK"

# Manually fetch a token
TOKEN=$(ssh root@localhost 'unset SLURM_JWT; /opt/slurm/current/bin/scontrol token' | grep SLURM_JWT | cut -d= -f2)
echo "Token: ${TOKEN:0:50}..."

# Test API with curl
curl -H X-SLURM-USER-TOKEN:$TOKEN -X GET 'http://localhost
```