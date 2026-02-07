# Configuration Guide

This guide covers all configuration options for the SLURM REST API Client Library.

## Client Configuration

### Basic Configuration

```go
import (
    "context"
    slurm "github.com/jontk/slurm-client"
    "github.com/jontk/slurm-client/pkg/auth"
)

// Basic client with auto-detected version
client, err := slurm.NewClient(context.Background(),
    slurm.WithBaseURL("http://your-slurm-host:6820"),
    slurm.WithAuth(auth.NewTokenAuth("your-token")),
)
```

### Full Configuration Options

```go
import (
    "net/http"
    "time"
    slurm "github.com/jontk/slurm-client"
    "github.com/jontk/slurm-client/pkg/auth"
)

// Custom HTTP client for advanced configuration
httpClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}

// Create client with all options
client, err := slurm.NewClient(context.Background(),
    slurm.WithBaseURL("http://your-slurm-host:6820"),
    slurm.WithAuth(auth.NewTokenAuth("your-jwt-token")),
    slurm.WithHTTPClient(httpClient),
    slurm.WithTimeout(30*time.Second),
    slurm.WithMaxRetries(3),
)
```

## Authentication Options

### Token Authentication (JWT)

```go
import "github.com/jontk/slurm-client/pkg/auth"

client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL("http://your-slurm-host:6820"),
    slurm.WithAuth(auth.NewTokenAuth("eyJhbGciOiJIUzI1NiIs...")),
)
```

### Basic Authentication

```go
client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL("http://your-slurm-host:6820"),
    slurm.WithAuth(auth.NewBasicAuth("slurm-user", "secure-password")),
)
```

### No Authentication

```go
client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL("http://your-slurm-host:6820"),
    slurm.WithAuth(auth.NewNoAuth()),
)
```

### Custom Authentication

```go
// Implement custom authentication via HTTP client
config.HTTPClient = &http.Client{
    Transport: &customAuthTransport{
        Base: http.DefaultTransport,
    },
}

type customAuthTransport struct {
    Base http.RoundTripper
}

func (t *customAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    // Add custom auth headers
    req.Header.Set("X-Custom-Auth", "custom-value")
    return t.Base.RoundTrip(req)
}
```

## Version Configuration

### Auto-Detection (Recommended)

Let the client automatically detect the API version:

```go
client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL("http://your-slurm-host:6820"),
    slurm.WithAuth(auth.NewTokenAuth("token")),
    // Version is auto-detected
)
```

### Manual Version Selection

Specify a specific API version:

```go
client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
    slurm.WithBaseURL("http://your-slurm-host:6820"),
    slurm.WithAuth(auth.NewTokenAuth("token")),
)
// Supported: v0.0.40, v0.0.41, v0.0.42, v0.0.43, v0.0.44
```

### Version Compatibility

| SLURM Version | API Version | Support Level |
|---------------|-------------|---------------|
| 23.11.x       | v0.0.43     | Full          |
| 23.02.x       | v0.0.42     | Full          |
| 22.05.x       | v0.0.41     | Full          |
| 21.08.x       | v0.0.40     | Basic         |

## Advanced Configuration

### Connection Pooling

```go
transport := &http.Transport{
    MaxIdleConns:        100,
    MaxIdleConnsPerHost: 10,
    MaxConnsPerHost:     10,
    IdleConnTimeout:     90 * time.Second,
    DisableKeepAlives:   false,
}

httpClient := &http.Client{
    Transport: transport,
    Timeout:   30 * time.Second,
}

client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL("http://your-slurm-host:6820"),
    slurm.WithAuth(auth.NewTokenAuth("token")),
    slurm.WithHTTPClient(httpClient),
)
```

### Proxy Configuration

```go
proxyURL, _ := url.Parse("http://proxy.example.com:8080")

httpClient := &http.Client{
    Transport: &http.Transport{
        Proxy: http.ProxyURL(proxyURL),
    },
}

client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL("http://your-slurm-host:6820"),
    slurm.WithAuth(auth.NewTokenAuth("token")),
    slurm.WithHTTPClient(httpClient),
)
```

### TLS Configuration

```go
tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS12,
    // Add certificates if needed
}

httpClient := &http.Client{
    Transport: &http.Transport{
        TLSClientConfig: tlsConfig,
    },
}

client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL("https://your-slurm-host:6820"),
    slurm.WithAuth(auth.NewTokenAuth("token")),
    slurm.WithHTTPClient(httpClient),
)
```

### Timeouts

```go
// Set global timeout for all requests
client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL("http://your-slurm-host:6820"),
    slurm.WithAuth(auth.NewTokenAuth("token")),
    slurm.WithTimeout(30*time.Second),
)

// Or use context-based timeouts for individual requests
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
jobs, err := client.Jobs().List(ctx, nil)
```

### Retry Configuration

```go
client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL("http://your-slurm-host:6820"),
    slurm.WithAuth(auth.NewTokenAuth("token")),
    slurm.WithMaxRetries(3),
    slurm.WithRetryWaitMin(1*time.Second),
    slurm.WithRetryWaitMax(30*time.Second),
)
```

## Environment Variables

The client can be configured via environment variables:

```bash
# Basic configuration
export SLURM_API_URL="http://slurm-host:6820"
export SLURM_API_VERSION="v0.0.43"

# Authentication
export SLURM_API_AUTH_TYPE="token"
export SLURM_API_TOKEN="your-jwt-token"

# Advanced
export SLURM_API_TIMEOUT="30s"
export SLURM_API_MAX_RETRIES="3"
export SLURM_API_DEBUG="true"
```

Loading from environment:

```go
import (
    "os"
    slurm "github.com/jontk/slurm-client"
    "github.com/jontk/slurm-client/pkg/auth"
)

baseURL := os.Getenv("SLURM_API_URL")
token := os.Getenv("SLURM_API_TOKEN")

client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL(baseURL),
    slurm.WithAuth(auth.NewTokenAuth(token)),
)
```

## Configuration Files

### YAML Configuration

```yaml
slurm:
  api:
    url: "http://slurm-host:6820"
    version: "v0.0.43"
  auth:
    type: "token"
    token: "${SLURM_API_TOKEN}"
  request:
    timeout: "30s"
    max_retries: 3
    rate_limit: 100
```

### Loading Configuration

```go
import (
    "context"
    slurm "github.com/jontk/slurm-client"
    "github.com/jontk/slurm-client/pkg/auth"
    "github.com/spf13/viper"
)

func LoadConfig(ctx context.Context) (slurm.Client, error) {
    viper.SetConfigFile("config.yaml")
    viper.SetEnvPrefix("SLURM")
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    baseURL := viper.GetString("api.url")
    token := viper.GetString("auth.token")

    return slurm.NewClient(ctx,
        slurm.WithBaseURL(baseURL),
        slurm.WithAuth(auth.NewTokenAuth(token)),
    )
}
```

## Logging Configuration

### Enable Debug Logging

```go
import "log"

// Enable debug logging
if os.Getenv("SLURM_CLIENT_DEBUG") == "true" {
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
    // Configure your logger for debug level
}
```

### Custom Logger

```go
// Enable debug mode for verbose logging
client, err := slurm.NewClient(ctx,
    slurm.WithBaseURL("http://your-slurm-host:6820"),
    slurm.WithAuth(auth.NewTokenAuth("token")),
    slurm.WithDebug(true),
)
```

## Best Practices

1. **Use Environment Variables** for sensitive data like tokens
2. **Enable Auto-Detection** for API version unless you need a specific version
3. **Set Appropriate Timeouts** based on your cluster size and network
4. **Implement Retry Logic** for transient failures
5. **Use Connection Pooling** for better performance
6. **Monitor API Usage** to avoid overwhelming the SLURM REST API

## See Also

- [Troubleshooting Guide](./troubleshooting.md)
- [Deployment Guide](./deployment.md)
- [API Documentation](./api)