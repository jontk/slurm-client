# Troubleshooting Guide

This guide helps you resolve common issues when using the SLURM REST API Client Library.

## Common Issues

### Connection Issues

#### Problem: Cannot connect to SLURM REST API
**Symptoms:**
- Connection refused errors
- Timeout errors
- 401 Unauthorized errors

**Solutions:**
1. Verify the SLURM REST API service is running:
   ```bash
   systemctl status slurmrestd
   ```

2. Check the API endpoint URL:
   ```go
   import (
       slurm "github.com/jontk/slurm-client"
       "github.com/jontk/slurm-client/pkg/auth"
   )

   client, err := slurm.NewClient(ctx,
       slurm.WithBaseURL("http://your-slurm-host:6820"), // Ensure correct host and port
       slurm.WithAuth(auth.NewTokenAuth("your-token")),
   )
   ```

3. Verify authentication credentials:
   - For token auth: Ensure token is valid and not expired
   - For basic auth: Verify username and password

4. Test connectivity with curl:
   ```bash
   curl -H "X-SLURM-USER-TOKEN: your-token" http://your-slurm-host:6820/slurm/v0.0.43/ping
   ```

### Version Compatibility

#### Problem: API version mismatch
**Symptoms:**
- "Unsupported API version" errors
- Method not implemented errors

**Solutions:**
1. Let the client auto-detect the version (recommended):
   ```go
   client, err := slurm.NewClient(ctx,
       slurm.WithBaseURL("http://your-slurm-host:6820"),
       slurm.WithAuth(auth.NewTokenAuth("token")),
       // Version will be auto-detected
   )
   ```

2. Or specify a version explicitly:
   ```go
   client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
       slurm.WithBaseURL("http://your-slurm-host:6820"),
       slurm.WithAuth(auth.NewTokenAuth("token")),
   )
   ```

3. Check what version your server supports:
   ```bash
   curl http://your-slurm-host:6820/openapi/v3
   ```

### Job Submission Issues

#### Problem: Job submission fails
**Symptoms:**
- Invalid job submission errors
- Missing required fields errors

**Solutions:**
1. Ensure all required fields are provided:
   ```go
   import slurm "github.com/jontk/slurm-client"

   job := &slurm.JobSubmission{
       Name:      "my-job",
       Partition: "default",
       Script:    "#!/bin/bash\necho 'Hello World'",
       CPUs:      1,
       Memory:    1024 * 1024 * 1024, // 1GB in bytes
   }

   response, err := client.Jobs().Submit(ctx, job)
   if err != nil {
       log.Printf("Job submission failed: %v", err)
   }
   ```

2. Check partition exists and is available:
   ```go
   partitions, err := client.Partitions().List(ctx, nil)
   if err != nil {
       log.Printf("Failed to list partitions: %v", err)
   }
   ```

### Performance Issues

#### Problem: Slow API responses
**Symptoms:**
- Long response times
- Timeouts on large queries

**Solutions:**
1. Use pagination for large result sets:
   ```go
   import slurm "github.com/jontk/slurm-client"

   opts := &slurm.ListJobsOptions{
       Limit:  100,
       Offset: 0,
   }
   jobs, err := client.Jobs().List(ctx, opts)
   ```

2. Set appropriate timeouts:
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()

   jobs, err := client.Jobs().List(ctx, nil)
   ```

3. Use filters to reduce result set size:
   ```go
   opts := &slurm.ListJobsOptions{
       States: []string{"RUNNING", "PENDING"},
       Limit:  50,
   }
   ```

### Authentication Issues

#### Problem: Authentication failures
**Symptoms:**
- 401 Unauthorized errors
- 403 Forbidden errors

**Solutions:**
1. For token authentication:
   ```go
   import "github.com/jontk/slurm-client/pkg/auth"

   client, err := slurm.NewClient(ctx,
       slurm.WithBaseURL("http://your-slurm-host:6820"),
       slurm.WithAuth(auth.NewTokenAuth("your-jwt-token")),
   )
   ```

2. For basic authentication:
   ```go
   client, err := slurm.NewClient(ctx,
       slurm.WithBaseURL("http://your-slurm-host:6820"),
       slurm.WithAuth(auth.NewBasicAuth("username", "password")),
   )
   ```

3. Check token expiration:
   ```go
   // JWT tokens expire - regenerate if expired
   // Check token with: jwt.io or decode in your auth system
   ```

### Adapter Conversion Errors

#### Problem: Type conversion failures
**Symptoms:**
- "Failed to convert response" errors
- Unexpected nil values in responses
- Field type mismatches

**Solutions:**
1. Check API version compatibility:
   ```go
   version := client.Version()
   log.Printf("Using API version: %s", version)
   ```

2. Some fields may not be available in older API versions:
   ```go
   import "github.com/jontk/slurm-client/pkg/errors"

   reservations, err := client.Reservations().List(ctx, nil)
   if err != nil {
       if errors.IsVersionNotSupported(err) {
           log.Println("Reservations require API v0.0.43+")
       }
   }
   ```

3. For complex type conversions, check field mappings:
   ```go
   // Some fields are automatically converted:
   // - NodeList strings → []string
   // - Flags arrays → []string
   // - Memory values → int64 (bytes)
   ```

### Mock Server Issues

#### Problem: Mock server behaves differently than real SLURM
**Symptoms:**
- Tests pass with mock but fail with real server
- Unexpected response formats
- Missing fields in mock responses

**Solutions:**
1. Use the built-in mock server for testing:
   ```go
   import "github.com/jontk/slurm-client/tests/mocks"

   mockServer := mocks.NewMockServer("v0.0.43")
   defer mockServer.Close()

   client, err := slurm.NewClient(ctx,
       slurm.WithBaseURL(mockServer.URL),
       slurm.WithAuth(auth.NewNoAuth()),
   )
   ```

2. Mock server limitations:
   - Best-effort support for v0.0.40/v0.0.41 (use fallback implementations)
   - Native support for v0.0.42, v0.0.43, v0.0.44
   - Some edge cases may differ from real SLURM behavior

3. Test against real SLURM when possible:
   ```bash
   export SLURM_REAL_SERVER_TEST=true
   export SLURM_SERVER_URL=http://your-test-slurm:6820
   export SLURM_JWT_TOKEN=your-test-token
   go test -v ./...
   ```

## Debugging Tips

### Enable Debug Logging
Enable debug mode when creating the client:
```go
client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL("http://your-slurm-host:6820"),
    slurm.WithAuth(auth.NewTokenAuth("token")),
    slurm.WithDebug(true), // Enables verbose logging
)
```

### Check API Health
Use the ping endpoint to verify connectivity:
```go
err := client.Info().Ping(ctx)
if err != nil {
    log.Printf("API health check failed: %v", err)
} else {
    log.Println("API is healthy")
}
```

### Inspect HTTP Traffic
Use a logging HTTP client to see raw requests/responses:
```go
import (
    "net/http"
    "net/http/httputil"
)

type loggingTransport struct {
    Transport http.RoundTripper
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    // Log request
    reqDump, _ := httputil.DumpRequestOut(req, true)
    log.Printf("Request:\n%s\n", reqDump)

    resp, err := t.Transport.RoundTrip(req)

    // Log response
    if resp != nil {
        respDump, _ := httputil.DumpResponse(resp, true)
        log.Printf("Response:\n%s\n", respDump)
    }

    return resp, err
}

httpClient := &http.Client{
    Transport: &loggingTransport{Transport: http.DefaultTransport},
}

client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL("http://your-slurm-host:6820"),
    slurm.WithAuth(auth.NewTokenAuth("token")),
    slurm.WithHTTPClient(httpClient),
)
```

## Performance Debugging

### Identifying Bottlenecks

1. **Use context timeouts to identify slow operations:**
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()

   start := time.Now()
   jobs, err := client.Jobs().List(ctx, nil)
   elapsed := time.Since(start)

   log.Printf("Jobs().List took %v", elapsed)
   ```

2. **Profile API calls:**
   ```go
   import "runtime/pprof"

   f, _ := os.Create("cpu.prof")
   pprof.StartCPUProfile(f)
   defer pprof.StopCPUProfile()

   // Your API calls here
   ```

3. **Check connection pooling:**
   ```go
   transport := &http.Transport{
       MaxIdleConns:        100,
       MaxIdleConnsPerHost: 10,
       IdleConnTimeout:     90 * time.Second,
   }

   httpClient := &http.Client{Transport: transport}

   client, err := slurm.NewClient(ctx,
       slurm.WithHTTPClient(httpClient),
   )
   ```

### Common Performance Issues

1. **Too many API calls:**
   - Use batch operations where available
   - Cache results that don't change frequently
   - Use pagination to limit result sizes

2. **Large payload transfers:**
   ```go
   // Instead of fetching all jobs:
   allJobs, _ := client.Jobs().List(ctx, nil)

   // Fetch only what you need:
   opts := &slurm.ListJobsOptions{
       States: []string{"RUNNING"},
       Limit:  100,
   }
   jobs, _ := client.Jobs().List(ctx, opts)
   ```

3. **Network latency:**
   - Deploy client close to SLURM server
   - Use connection pooling
   - Enable HTTP keep-alive

## Common Error Messages

### "connection refused"
- SLURM REST API service (slurmrestd) is not running
- Check: `systemctl status slurmrestd`

### "401 Unauthorized"
- Invalid or expired authentication token
- Check token is correctly set in client
- Verify token with SLURM admin

### "404 Not Found"
- Endpoint doesn't exist in this API version
- Check API version compatibility
- Use version auto-detection

### "502 Bad Gateway"
- slurmdbd is not properly connected to slurmctld
- Authentication plugin mismatch (JWT vs munge)
- Check slurmctld and slurmdbd logs

### "Version not supported"
- Attempting to use features not available in current API version
- Check feature availability matrix in README
- Upgrade SLURM or use compatible features

## Getting Help

If you continue to experience issues:

1. **Check existing resources:**
   - [GitHub Issues](https://github.com/jontk/slurm-client/issues) for similar problems
   - [Examples](../examples/) for working code samples
   - [API Documentation](./api/) for method signatures

2. **Collect diagnostic information:**
   - Enable debug logging
   - Capture full error messages
   - Note SLURM and API versions
   - Try to create minimal reproduction case

3. **Open a new issue with:**
   - SLURM version (e.g., 24.05.x)
   - API version (e.g., v0.0.43)
   - Go version
   - Code snippet reproducing the issue
   - Full error messages and stack traces
   - Steps to reproduce

4. **Use the diagnostic script:**
   ```bash
   ./scripts/diagnose-slurm-auth.sh
   ```

## See Also

- [Configuration Guide](./configuration.md) - Client configuration options
- [API Documentation](./api/) - Complete API reference
- [Examples](../examples/) - Working code examples
- [Version Support](./VERSION_SUPPORT.md) - API version compatibility
- [MIGRATION.md](../MIGRATION.md) - Migrating between versions