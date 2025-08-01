# Benchmarks

This directory contains performance benchmarks for the slurm-client library.

## Running Benchmarks

### Run all benchmarks
```bash
go test -bench=. -benchmem ./benchmarks/...
```

### Run specific benchmarks
```bash
# Auth benchmarks only
go test -bench=BenchmarkAuth -benchmem ./benchmarks/...

# Error handling benchmarks
go test -bench=BenchmarkError -benchmem ./benchmarks/...

# Retry mechanism benchmarks
go test -bench=BenchmarkRetry -benchmem ./benchmarks/...
```

### Generate baseline
```bash
go test -bench=. -benchmem ./benchmarks/... > benchmarks/baseline.txt
```

### Compare against baseline
```bash
# Install benchstat if not already installed
go install golang.org/x/perf/cmd/benchstat@latest

# Run new benchmarks
go test -bench=. -benchmem ./benchmarks/... > benchmarks/new.txt

# Compare
benchstat benchmarks/baseline.txt benchmarks/new.txt
```

## Benchmark Categories

### Authentication (`auth_bench_test.go`)
- Token authentication performance
- API key authentication performance
- Basic authentication performance
- Authentication chain performance

### Error Handling (`errors_bench_test.go`)
- Error creation performance
- Error wrapping performance
- Error code checking performance
- Error formatting performance

### Retry Logic (`retry_bench_test.go`)
- Retryable check performance
- Backoff calculation performance
- Retry loop simulation

## Performance Goals

- Authentication operations: < 100ns
- Error creation: < 500ns
- Error checking: < 50ns
- Retry decision: < 100ns

## Continuous Benchmarking

The CI pipeline runs benchmarks on each PR to detect performance regressions:

1. Benchmarks are run against the main branch
2. Results are compared with the PR branch
3. Significant regressions (>10%) are flagged
4. Results are posted as a comment on the PR

## Adding New Benchmarks

When adding new benchmarks:

1. Create focused benchmarks that test one specific operation
2. Use meaningful benchmark names (e.g., `BenchmarkAuthTokenApply`)
3. Include memory allocation benchmarks with `-benchmem`
4. Document what the benchmark measures
5. Set realistic performance goals

## Optimization Tips

Based on benchmark results:

1. **Authentication**: Pre-compute base64 encodings where possible
2. **Errors**: Use error pools for frequently created errors
3. **Retry**: Cache retry decisions for similar error types
4. **General**: Minimize allocations in hot paths