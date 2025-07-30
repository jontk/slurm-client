# SLURM Client v0.0.43 Integration Test Plan

## Test Environment Configuration

### Server Details
- **Server URL**: http://rocky9.ar.jontk.com:6820
- **API Version**: v0.0.43
- **JWT Token**: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI2NTM4Mjk5NzYsImlhdCI6MTc1MzgyOTk3Niwic3VuIjoicm9vdCJ9.-z8Cq_wHuOxNJ7KHHTboX3l9r6JBtSD1RxQUgQR9owE

### Test Configuration
```bash
export SLURM_REAL_SERVER_TEST=true
export SLURM_REAL_INTEGRATION_TEST=true
export SLURM_SERVER_URL=http://rocky9.ar.jontk.com:6820
export SLURM_API_VERSION=v0.0.43
export SLURM_JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI2NTM4Mjk5NzYsImlhdCI6MTc1MzgyOTk3Niwic3VuIjoicm9vdCJ9.-z8Cq_wHuOxNJ7KHHTboX3l9r6JBtSD1RxQUgQR9owE"
export SLURM_TEST_DEBUG=true
export SLURM_INSECURE_SKIP_VERIFY=true
```

## Test Plan Structure

### Phase 1: Basic Connectivity and Authentication Tests

#### 1.1 Connection Validation
- **Test**: Basic ping endpoint
- **Expected**: Successful response with 200 status
- **Validates**: Network connectivity, URL correctness

#### 1.2 Authentication Verification
- **Test**: JWT token validation
- **Expected**: Authenticated requests succeed
- **Validates**: Token format, expiration, authorization

#### 1.3 Version Compatibility
- **Test**: API version endpoint
- **Expected**: Returns v0.0.43 or compatible version
- **Validates**: API version matching

#### 1.4 Cluster Information
- **Test**: Get cluster metadata
- **Expected**: Valid cluster name, configuration
- **Validates**: Basic info endpoint functionality

### Phase 2: Manager Endpoint Tests

#### 2.1 Jobs Manager
##### 2.1.1 List Operations
- **Test**: List jobs with various filters
  - Default listing (limit: 10)
  - Filter by states (running, pending, completed)
  - Large limit (100+ jobs)
  - Pagination support
- **Expected**: Valid job list response
- **Performance**: < 2s for 100 jobs

##### 2.1.2 Get Operations
- **Test**: Retrieve individual job details
  - Valid job ID
  - Non-existent job ID (error handling)
- **Expected**: Complete job information or appropriate error

##### 2.1.3 Submit Operations
- **Test**: Submit various job types
  - Simple batch script
  - Job with environment variables
  - Job with resource requirements (CPU, memory)
  - Job with working directory
  - Job with dependencies
- **Expected**: Valid job ID returned

##### 2.1.4 Modify Operations
- **Test**: Update job properties
  - Hold/release job
  - Modify time limit
  - Change priority
- **Expected**: Successful modification

##### 2.1.5 Cancel Operations
- **Test**: Cancel jobs
  - Cancel pending job
  - Cancel running job
  - Cancel completed job (error case)
- **Expected**: Appropriate state transitions

#### 2.2 Nodes Manager
##### 2.2.1 List Operations
- **Test**: List nodes with filters
  - All nodes
  - Filter by state (idle, allocated, down)
  - Filter by features
- **Expected**: Valid node list
- **Performance**: < 1s for typical cluster

##### 2.2.2 Get Operations
- **Test**: Individual node details
  - Valid node name
  - Invalid node name
- **Expected**: Complete node information

##### 2.2.3 Update Operations
- **Test**: Modify node state (if permitted)
  - Set node to drain
  - Resume node
- **Expected**: State changes reflected

#### 2.3 Partitions Manager
##### 2.3.1 List Operations
- **Test**: List all partitions
  - Default listing
  - Filter by state
- **Expected**: Valid partition list

##### 2.3.2 Get Operations
- **Test**: Individual partition details
- **Expected**: Complete partition configuration

#### 2.4 QoS Manager
##### 2.4.1 List Operations
- **Test**: List QoS entries
- **Expected**: Valid QoS list (may require database)

##### 2.4.2 Get Operations
- **Test**: Individual QoS details
- **Expected**: Complete QoS configuration

#### 2.5 Users Manager
##### 2.5.1 List Operations
- **Test**: List users
- **Expected**: Valid user list (may require database)

##### 2.5.2 Get Operations
- **Test**: Individual user details
- **Expected**: User information

#### 2.6 Accounts Manager
##### 2.6.1 List Operations
- **Test**: List accounts
- **Expected**: Valid account list (may require database)

##### 2.6.2 Get Operations
- **Test**: Individual account details
- **Expected**: Account information

#### 2.7 Associations Manager
##### 2.7.1 List Operations
- **Test**: List associations
- **Expected**: Valid association list (may require database)

### Phase 3: Error Handling Scenarios

#### 3.1 Invalid Resource Access
- **Test**: Access non-existent resources
  - Invalid job ID (999999999)
  - Invalid node name
  - Invalid partition name
  - Invalid QoS name
- **Expected**: Appropriate error codes and messages

#### 3.2 Invalid Operations
- **Test**: Perform invalid operations
  - Submit job with invalid partition
  - Submit job with excessive resources
  - Cancel already completed job
  - Modify non-existent job
- **Expected**: Clear error messages

#### 3.3 Authentication Errors
- **Test**: Various auth failures
  - Invalid token
  - Expired token (if testable)
  - Missing token
- **Expected**: 401/403 status codes

#### 3.4 Network Errors
- **Test**: Network failure scenarios
  - Timeout handling
  - Connection refused
  - DNS resolution failure
- **Expected**: Appropriate error handling and retries

### Phase 4: Performance Benchmarks

#### 4.1 Latency Tests
- **Test**: Measure operation latencies
  - Ping: < 100ms
  - List jobs (10): < 500ms
  - List jobs (100): < 2s
  - Get single job: < 200ms
  - Submit job: < 500ms

#### 4.2 Throughput Tests
- **Test**: Rapid successive requests
  - 10 ping requests: < 1s total
  - 50 list operations: < 10s total
- **Expected**: No failures under load

#### 4.3 Concurrent Operations
- **Test**: Parallel requests
  - 10 concurrent job submissions
  - 20 concurrent list operations
  - Mixed read/write operations
- **Expected**: > 80% success rate

#### 4.4 Large Data Handling
- **Test**: Large result sets
  - List 1000+ jobs
  - List all nodes in large cluster
- **Expected**: Graceful handling, pagination

### Phase 5: API Version-Specific Features

#### 5.1 v0.0.43 New Features
- **Test**: Features specific to v0.0.43
  - New job submission parameters
  - Enhanced filtering options
  - Additional resource attributes
- **Expected**: Feature availability and functionality

#### 5.2 Backward Compatibility
- **Test**: Legacy feature support
  - Old parameter names
  - Deprecated endpoints
- **Expected**: Graceful handling or clear migration path

### Phase 6: Integration Workflows

#### 6.1 Complete Job Lifecycle
- **Test**: End-to-end job workflow
  1. Submit job
  2. Monitor status
  3. Check resource allocation
  4. Wait for completion
  5. Retrieve output
  6. Clean up
- **Expected**: Smooth workflow execution

#### 6.2 Resource Discovery
- **Test**: Cluster exploration workflow
  1. Get cluster info
  2. List all partitions
  3. Check node availability
  4. Identify suitable resources
  5. Submit optimized job
- **Expected**: Complete resource picture

#### 6.3 Multi-Job Dependencies
- **Test**: Job dependency chains
  1. Submit parent job
  2. Submit dependent jobs
  3. Monitor dependency resolution
  4. Verify execution order
- **Expected**: Correct dependency handling

## Test Execution Plan

### Environment Setup
```bash
# 1. Set environment variables (as shown above)
# 2. Verify connectivity
curl -H "Authorization: Bearer $SLURM_JWT_TOKEN" \
     http://rocky9.ar.jontk.com:6820/slurm/v0.0.43/ping

# 3. Run integration tests
go test -v ./tests/integration/v0_0_43_integration_test.go \
        -run TestV0043IntegrationSuite \
        -timeout 30m
```

### Test Priorities
1. **Critical**: Basic connectivity, authentication, job operations
2. **High**: Node/partition listing, error handling
3. **Medium**: QoS/User/Account operations (may require database)
4. **Low**: Performance benchmarks, edge cases

### Success Criteria
- All critical tests pass
- > 90% of high priority tests pass
- > 80% of medium priority tests pass
- Performance within specified bounds
- No data corruption or system instability

### Risk Mitigation
- All test jobs use minimal resources
- Automatic cleanup of test resources
- Timeout limits on all operations
- Non-destructive testing only
- Isolation from production workloads

## Reporting

### Test Report Format
```
Test Suite: v0.0.43 Integration Tests
Date: [timestamp]
Duration: [total time]

Summary:
- Total Tests: X
- Passed: X (X%)
- Failed: X
- Skipped: X

Critical Failures:
- [List any critical test failures]

Performance Metrics:
- Average Latency: Xms
- Throughput: X ops/sec
- Error Rate: X%

Recommendations:
- [Any findings or improvements]
```

### Continuous Integration
- Run subset of tests on every commit
- Full test suite nightly
- Performance benchmarks weekly
- Store historical metrics for trend analysis