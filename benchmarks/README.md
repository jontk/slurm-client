# Benchmark Baseline

This directory contains the benchmark baseline used for performance comparison in CI.

## Generating the Baseline

To generate a fresh baseline, run:

```bash
go test -bench=. -benchmem -count=10 -run=^$ ./tests/benchmarks/... > benchmarks/baseline.txt
```

Note: This may take several minutes to complete as it runs each benchmark 10 times for statistical accuracy.

## Using the Baseline

The CI workflow (`.github/workflows/ci.yml`) uses `benchstat` to compare new benchmark results against this baseline. If significant performance regressions are detected (>10% slower), they will be reported in the PR summary.

## Performance Validation

For more rigorous performance validation, see the nightly performance validation workflow (`.github/workflows/performance-validation.yml`), which runs benchmarks with more iterations and checks for regressions.
