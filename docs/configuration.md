# Configuration Guide

This guide covers all configuration options for the SLURM REST API Client Library.

## Client Configuration

### Basic Configuration

```go
import (
    slurm "github.com/jontk/slurm-client"
    "github.com/jontk/slurm-client/pkg/client/factory"
)

config := &slurm.ClientConfig{
    BaseURL: "http://your-slurm-host:6820",
    Version: "v0.0.43", // Optional - auto-detected if not specified
}

client, err := factory.NewClient(config)
```

### Full Configuration Options

```go
config := &slurm.ClientConfig{
    // Required
    BaseURL: "http://your-slurm-host:6820",

    // Optional - API version (auto-detected if not specified)
    Version: "v0.0.43",

    // Optional - Custom HTTP client
    HTTPClient: &http.Client{
        Timeout: 30 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
            IdleConnTimeout:     90 * time.Second,
        },
    },

    // Optional - Authentication configuration
    Authentication: &slurm.AuthConfig{
        Type:     "token",  // "token", "basic", or "munge"
        Token:    "your-jwt-token",
        Username: "your-username",
        Password: "your-password",
    },

    // Optional - Request configuration
    RequestConfig: &slurm.RequestConfig{
        Timeout:     30 * time.Second,
        MaxRetries:  3,
        RetryDelay:  time.Second,
        RateLimit:   100, // requests per second
    },
}
```

## Authentication Options

### Token Authentication (JWT)

```go
config.Authentication = &slurm.AuthConfig{
    Type:  "token",
    Token: "eyJhbGciOiJIUzI1NiIs...", // Your JWT token
}
```

### Basic Authentication

```go
config.Authentication = &slurm.AuthConfig{
    Type:     "basic",
    Username: "slurm-user",
    Password: "secure-password",
}
```

### Munge Authentication

```go
config.Authentication = &slurm.AuthConfig{
    Type: "munge",
    // Munge credentials are handled by the system
}
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

### Auto-Detection

Let the client automatically detect the API version:

```go
config := &slurm.ClientConfig{
    BaseURL: "http://your-slurm-host:6820",
    // Version is not specified - will be auto-detected
}
```

### Manual Version Selection

Specify a specific API version:

```go
config.Version = "v0.0.43" // Supported: v0.0.40, v0.0.41, v0.0.42, v0.0.43
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

config.HTTPClient = &http.Client{
    Transport: transport,
    Timeout:   30 * time.Second,
}
```

### Proxy Configuration

```go
proxyURL, _ := url.Parse("http://proxy.example.com:8080")

config.HTTPClient = &http.Client{
    Transport: &http.Transport{
        Proxy: http.ProxyURL(proxyURL),
    },
}
```

### TLS Configuration

```go
tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS12,
    // Add certificates if needed
}

config.HTTPClient = &http.Client{
    Transport: &http.Transport{
        TLSClientConfig: tlsConfig,
    },
}
```

### Timeouts

```go
config.RequestConfig = &slurm.RequestConfig{
    Timeout: 30 * time.Second, // Overall request timeout
}

// Or via HTTP client
config.HTTPClient = &http.Client{
    Timeout: 30 * time.Second,
}

// Context-based timeouts
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
jobs, err := client.Jobs().List(ctx, nil)
```

### Retry Configuration

```go
config.RequestConfig = &slurm.RequestConfig{
    MaxRetries:  3,
    RetryDelay:  time.Second,
    RetryPolicy: func(resp *http.Response, err error) bool {
        // Custom retry logic
        if err != nil {
            return true // Retry on network errors
        }
        return resp.StatusCode >= 500 // Retry on server errors
    },
}
```

### Rate Limiting

```go
config.RequestConfig = &slurm.RequestConfig{
    RateLimit: 100, // Max 100 requests per second
}
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
)

config := &slurm.ClientConfig{
    BaseURL: os.Getenv("SLURM_API_URL"),
    Version: os.Getenv("SLURM_API_VERSION"),
}

if authType := os.Getenv("SLURM_API_AUTH_TYPE"); authType != "" {
    config.Authentication = &slurm.AuthConfig{
        Type:  authType,
        Token: os.Getenv("SLURM_API_TOKEN"),
    }
}
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
    slurm "github.com/jontk/slurm-client"
    "github.com/spf13/viper"
)

func LoadConfig() (*slurm.ClientConfig, error) {
    viper.SetConfigFile("config.yaml")
    viper.SetEnvPrefix("SLURM")
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    config := &slurm.ClientConfig{
        BaseURL: viper.GetString("api.url"),
        Version: viper.GetString("api.version"),
    }

    if viper.IsSet("auth.type") {
        config.Authentication = &slurm.AuthConfig{
            Type:  viper.GetString("auth.type"),
            Token: viper.GetString("auth.token"),
        }
    }

    return config, nil
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
type Logger interface {
    Debug(args ...interface{})
    Info(args ...interface{})
    Warn(args ...interface{})
    Error(args ...interface{})
}

// Set custom logger
client.SetLogger(customLogger)
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