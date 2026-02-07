# v0.0.43 Integration Test Execution Guide

## Quick Start

### 1. Run All Tests
```bash
./tests/integration/run_v0043_tests.sh
```

### 2. Run Fast Tests Only (Skip Slow Tests)
```bash
./tests/integration/run_v0043_tests.sh --fast
```

### 3. Run Specific Test Phase
```bash
# Run only connectivity tests
./tests/integration/run_v0043_tests.sh "" TestPhase1_BasicConnectivity

# Run only Jobs Manager tests
./tests/integration/run_v0043_tests.sh "" TestPhase2_JobsManager

# Run only error handling tests
./tests/integration/run_v0043_tests.sh "" TestPhase3_ErrorHandling

# Run only performance benchmarks
./tests/integration/run_v0043_tests.sh "" TestPhase4_PerformanceBenchmarks
```

## Manual Execution

### Full Test Suite
```bash
export SLURM_REAL_SERVER_TEST=true
export SLURM_REAL_INTEGRATION_TEST=true
export SLURM_SERVER_URL="http://localhost:6820"
export SLURM_API_VERSION="v0.0.43"
export SLURM_JWT_TOKEN="your-jwt-token-here"
export SLURM_TEST_DEBUG=true
export SLURM_INSECURE_SKIP_VERIFY=true

go test -v ./tests/integration/v0_0_43_integration_test.go \
        ./tests/integration/test_helpers.go \
        -timeout 30m
```

### Individual Test Methods
```bash
# Test specific functionality
go test -v ./tests/integration/v0_0_43_integration_test.go \
        ./tests/integration/test_helpers.go \
        -run TestV0043IntegrationSuite/TestPhase2_JobsManager/2.1.3_SubmitOperations
```

## Test Phases Overview

### Phase 1: Basic Connectivity (Fast)
- Connection validation
- Authentication verification
- Version compatibility
- Cluster information

### Phase 2: Manager Endpoints (Medium)
- Jobs Manager (CRUD operations)
- Nodes Manager (List/Get)
- Partitions Manager (List/Get)
- Database-dependent managers (QoS, Users, Accounts)

### Phase 3: Error Handling (Fast)
- Invalid resource access
- Invalid operations
- Error response validation

### Phase 4: Performance Benchmarks (Slow)
- Latency tests
- Throughput tests
- Concurrent operations
- Large data handling

### Phase 5: Version-Specific Features (Medium)
- v0.0.43 new features
- Advanced job options

### Phase 6: Integration Workflows (Slow)
- Complete job lifecycle
- Resource discovery workflow
- Multi-step operations

## Debug Options

### Enable Verbose Logging
```bash
export SLURM_TEST_DEBUG=true
```

### Skip Database Tests
```bash
export SLURM_SKIP_DATABASE_TESTS=true
```

### Set Custom Timeout
```bash
go test -v ... -timeout 60m  # 60 minute timeout
```

## Expected Results

### Successful Test Output
```
=== PHASE 1: Basic Connectivity and Authentication ===
--- PASS: TestV0043IntegrationSuite/TestPhase1_BasicConnectivity (0.15s)
    --- PASS: TestV0043IntegrationSuite/TestPhase1_BasicConnectivity/1.1_ConnectionValidation (0.05s)
        ✓ Ping successful (latency: 45ms)
    --- PASS: TestV0043IntegrationSuite/TestPhase1_BasicConnectivity/1.2_AuthenticationVerification (0.03s)
        ✓ JWT authentication verified
...
```

### Test Report
At the end of the test run, a comprehensive report is generated:

```
===========================================================
V0.0.43 INTEGRATION TEST REPORT
===========================================================
Date: 2024-01-29T10:00:00Z
Duration: 5m30s
Server: http://localhost:6820
API Version: v0.0.43

Performance Metrics:
  - ping: 45ms
  - list_jobs_10: 120ms
  - list_jobs_100: 850ms
  - job_submit: 230ms

Average Operation Latency: 311ms
Jobs Submitted: 5
===========================================================
```

## Troubleshooting

### Connection Refused
```bash
# Check server is accessible
curl -I http://localhost:6820/slurm/v0.0.43/ping

# Check token is valid
curl -H "Authorization: Bearer $SLURM_JWT_TOKEN" \
     http://localhost:6820/slurm/v0.0.43/ping
```

### Authentication Errors
```bash
# Verify token format
echo $SLURM_JWT_TOKEN | cut -d. -f2 | base64 -d | jq .

# Check token expiration
# exp: 2653829976 (year 2054)
```

### Database Connection Errors
Many operations (QoS, Users, Accounts) require slurmdbd. If these fail with "Unable to connect to database", it's expected behavior when slurmdbd is not configured.

### Test Cleanup
The test suite automatically cleans up submitted jobs. If cleanup fails, manually cancel test jobs:

```bash
# List test jobs
squeue -u root | grep v0043

# Cancel specific job
scancel <job_id>

# Cancel all test jobs
squeue -u root | grep v0043 | awk '{print $1}' | xargs scancel
```

## Continuous Integration

### GitHub Actions Example
```yaml
name: v0.0.43 Integration Tests

on:
  push:
    branches: [ main ]
  schedule:
    - cron: '0 2 * * *'  # Daily at 2 AM

jobs:
  integration-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run v0.0.43 Integration Tests
        env:
          SLURM_REAL_SERVER_TEST: true
          SLURM_SERVER_URL: ${{ secrets.SLURM_SERVER_URL }}
          SLURM_JWT_TOKEN: ${{ secrets.SLURM_JWT_TOKEN }}
        run: |
          ./tests/integration/run_v0043_tests.sh
```

## Performance Baselines

Expected performance for v0.0.43:

| Operation | Target | Acceptable |
|-----------|--------|------------|
| Ping | < 50ms | < 100ms |
| Get Job | < 100ms | < 200ms |
| List 10 Jobs | < 200ms | < 500ms |
| List 100 Jobs | < 1s | < 2s |
| Submit Job | < 300ms | < 500ms |
| List Nodes | < 500ms | < 1s |

## Next Steps

After running the tests:

1. Review the test report for any failures
2. Check performance metrics against baselines
3. Investigate any unexpected errors
4. Update test plan based on findings
5. Consider adding new test cases for edge cases discovered