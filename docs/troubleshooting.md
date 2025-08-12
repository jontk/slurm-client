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
   config := &interfaces.ClientConfig{
       BaseURL: "http://your-slurm-host:6820", // Ensure correct host and port
   }
   ```

3. Verify authentication credentials:
   - For token auth: Ensure token is valid and not expired
   - For basic auth: Verify username and password
   - For munge auth: Check munge service is running

### Version Compatibility

#### Problem: API version mismatch
**Symptoms:**
- "Unsupported API version" errors
- Method not implemented errors

**Solutions:**
1. Let the client auto-detect the version:
   ```go
   config := &interfaces.ClientConfig{
       BaseURL: "http://your-slurm-host:6820",
       // Version will be auto-detected
   }
   ```

2. Or specify the correct version:
   ```go
   config.Version = "v0.0.43" // Use your SLURM API version
   ```

### Job Submission Issues

#### Problem: Job submission fails
**Symptoms:**
- Invalid job submission errors
- Missing required fields errors

**Solutions:**
1. Ensure all required fields are provided:
   ```go
   job := &interfaces.JobSubmission{
       Name:      "my-job",
       Partition: "default",
       Script:    "#!/bin/bash\necho 'Hello World'",
       // Add other required fields
   }
   ```

2. Check partition exists and is available:
   ```go
   partitions, err := client.Partitions().List(ctx, nil)
   ```

### Performance Issues

#### Problem: Slow API responses
**Symptoms:**
- Long response times
- Timeouts on large queries

**Solutions:**
1. Use pagination for large result sets:
   ```go
   opts := &interfaces.ListJobsOptions{
       Limit:  100,
       Offset: 0,
   }
   ```

2. Set appropriate timeouts:
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   ```

### Authentication Issues

#### Problem: Authentication failures
**Symptoms:**
- 401 Unauthorized errors
- 403 Forbidden errors

**Solutions:**
1. For token authentication:
   ```go
   config.Authentication = &interfaces.AuthConfig{
       Type:  "token",
       Token: "your-jwt-token",
   }
   ```

2. For basic authentication:
   ```go
   config.Authentication = &interfaces.AuthConfig{
       Type:     "basic",
       Username: "your-username",
       Password: "your-password",
   }
   ```

## Debugging Tips

### Enable Debug Logging
Set environment variables to enable detailed logging:
```bash
export SLURM_CLIENT_DEBUG=true
export SLURM_CLIENT_LOG_LEVEL=debug
```

### Check API Health
Use the ping endpoint to verify connectivity:
```go
err := client.Info().Ping(ctx)
if err != nil {
    log.Printf("API health check failed: %v", err)
}
```

### Inspect Raw Responses
Enable HTTP request/response logging:
```go
config.HTTPClient = &http.Client{
    Transport: &loggingRoundTripper{
        Transport: http.DefaultTransport,
    },
}
```

## Getting Help

If you continue to experience issues:

1. Check the [GitHub Issues](https://github.com/jontk/slurm-client/issues) for similar problems
2. Enable debug logging and collect error messages
3. Open a new issue with:
   - SLURM version
   - API version
   - Code snippet reproducing the issue
   - Full error messages

## See Also

- [Configuration Guide](./configuration.md)
- [API Documentation](./api)
- [Examples](../examples)