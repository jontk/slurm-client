# Real SLURM Server Testing Guide

## Overview

This guide documents the setup and findings from testing the SLURM Go client against a real SLURM server.

## Test Environment

- **SLURM Version**: 25.05.0
- **Server**: localhost:6820
- **Cluster Name**: localhost
- **Authentication**: JWT tokens

## Running Real Server Tests

To run tests against a real SLURM server:

```bash
# Set environment variables
export SLURM_REAL_SERVER_TEST=true
export SLURM_SERVER_URL=http://localhost
export SLURM_JWT_TOKEN=<your-jwt-token>

# Run specific test files
go test -v ./tests/integration/adapter_ping_test.go
go test -v ./tests/integration/adapter_jobs_test.go
go test -v ./tests/integration/adapter_real_server_test.go
```

## Available vs Unavailable Endpoints

### ✅ Working Endpoints (slurmctld)

These endpoints work without slurmdbd:

1. **Ping** (`/slurm/v0.0.43/ping/`)
   - Basic connectivity test
   - Returns cluster info and version

2. **Jobs** (`/slurm/v0.0.43/jobs/`)
   - List and manage jobs
   - No database persistence required

3. **Nodes** (`/slurm/v0.0.43/nodes/`)
   - Get node information
   - Hardware details (CPUs, memory)

4. **Diagnostics** (`/slurm/v0.0.43/diag/`)
   - Server statistics
   - Performance metrics

### ❌ Unavailable Endpoints (require slurmdbd)

These endpoints fail with HTTP 502 when slurmdbd is not connected:

1. **QoS** (`/slurmdb/v0.0.43/qos/`)
   - Quality of Service management
   - Requires database backend

2. **Partitions** (`/slurm/v0.0.43/partitions/`)
   - Fails with TRES (Trackable RESources) errors
   - Needs database for resource tracking

3. **Accounts** (`/slurmdb/v0.0.43/accounts/`)
   - User account management
   - Database-backed

4. **Associations** (`/slurmdb/v0.0.43/associations/`)
   - User-account-cluster relationships
   - Database-backed

## Root Cause: Authentication Mismatch

The slurmdbd connection fails due to authentication plugin mismatch:

```
error: auth_g_unpack: authentication plugin auth/jwt(102) not found
error: slurm_unpack_received_msg: auth_g_unpack: REQUEST_PERSIST_INIT has authentication error
error: CONN:13 Failed to unpack SLURM_PERSIST_INIT message
```

### What's Happening

1. **slurmctld** (main SLURM daemon) is configured to use JWT authentication
2. **slurmdbd** (database daemon) is trying to connect but doesn't have the JWT auth plugin
3. The persistent connection between slurmctld and slurmdbd fails
4. All database-backed operations return HTTP 502 (Bad Gateway)

### Common Solutions

1. **Ensure matching auth plugins**:
   ```ini
   # slurm.conf
   AuthType=auth/jwt
   
   # slurmdbd.conf
   AuthType=auth/jwt
   ```

2. **Install JWT plugin on slurmdbd server**:
   ```bash
   # Check if JWT plugin is installed
   ls /usr/lib64/slurm/auth_jwt.so
   ```

3. **Use munge authentication** (if JWT not needed):
   ```ini
   # Both slurm.conf and slurmdbd.conf
   AuthType=auth/munge
   ```

## Test Implementation Notes

### Handling Missing slurmdbd

Our tests gracefully handle the missing database connection:

```go
// Check for slurmdbd errors
if err != nil {
    if strings.Contains(err.Error(), "Unable to connect to database") || 
       strings.Contains(err.Error(), "Failed to open slurmdbd connection") {
        t.Skip("Skipping test: slurmdbd is not connected")
        return
    }
}

// For HTTP responses
if resp.StatusCode() == 502 {
    t.Skip("Skipping test: slurmdbd is not connected (HTTP 502)")
    return
}
```

### Test Coverage

Despite the slurmdbd issue, we can test:

- ✅ Authentication flow (JWT tokens)
- ✅ Basic API connectivity
- ✅ Non-database operations
- ✅ Error handling for database failures
- ✅ Version detection and compatibility

### Future Improvements

1. **Mock slurmdbd responses** for QoS testing
2. **Docker compose setup** with properly configured slurmdbd
3. **Integration test categories**:
   - `requires_database` - Skip if slurmdbd unavailable
   - `core_functionality` - Always run
4. **Health check endpoint** to detect slurmdbd status

## Example Test Output

### Successful Ping Test
```
=== RUN   TestPingWithRealServer/Ping_SLURM_Server
    Ping response received
    SLURM Version: 25.05.0
    SLURM Release: 25.05.0
    SLURM Cluster: localhost
    Ping 1: hostname=localhost, status=UP, responding=true, latency=1296
```

### Failed QoS Test (Expected)
```
=== RUN   TestAdapterWithRealServer/List_QoS_on_Real_Server
    Error listing QoS: [SLURM_DAEMON_DOWN] HTTP 502: Bad Gateway
    Skipping test: slurmdbd is not connected. This is expected in test environments.
```

## Conclusion

The SLURM Go client successfully communicates with the SLURM REST API for non-database operations. The adapter pattern correctly handles version differences and error cases. Full functionality requires a properly configured slurmdbd with matching authentication configuration.