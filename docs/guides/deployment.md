# Deployment Guide

This guide covers deploying applications that use the SLURM REST API Client Library.

## Prerequisites

- Go 1.20 or later
- Access to SLURM cluster with REST API enabled
- Appropriate authentication credentials

## Building Your Application

### Basic Build
```bash
go build -o myapp main.go
```

### Production Build
```bash
# With optimizations and stripped debug info
go build -ldflags="-s -w" -o myapp main.go

# With version information
VERSION=$(git describe --tags --always --dirty)
go build -ldflags="-s -w -X main.version=$VERSION" -o myapp main.go
```

### Cross-Compilation
```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o myapp-linux-amd64 main.go

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o myapp-linux-arm64 main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o myapp-darwin-amd64 main.go
```

## Configuration Management

### Environment Variables
```bash
# Basic configuration
export SLURM_API_URL="http://slurm-host:6820"
export SLURM_API_VERSION="v0.0.43"
export SLURM_API_TOKEN="your-jwt-token"

# Advanced configuration
export SLURM_API_TIMEOUT="30s"
export SLURM_API_MAX_RETRIES="3"
export SLURM_API_DEBUG="false"
```

### Configuration File
Create a configuration file (e.g., `config.yaml`):
```yaml
slurm:
  api_url: "http://slurm-host:6820"
  api_version: "v0.0.43"
  auth:
    type: "token"
    token: "${SLURM_API_TOKEN}"
  timeout: "30s"
  max_retries: 3
```

### Loading Configuration
```go
import (
    "github.com/spf13/viper"
)

func loadConfig() (*interfaces.ClientConfig, error) {
    viper.SetConfigFile("config.yaml")
    viper.SetEnvPrefix("SLURM")
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    return &interfaces.ClientConfig{
        BaseURL: viper.GetString("api_url"),
        Version: viper.GetString("api_version"),
        Authentication: &interfaces.AuthConfig{
            Type:  viper.GetString("auth.type"),
            Token: viper.GetString("auth.token"),
        },
    }, nil
}
```

## Container Deployment

### Dockerfile
```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o myapp

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/myapp .
COPY --from=builder /app/config.yaml .

EXPOSE 8080
CMD ["./myapp"]
```

### Docker Compose
```yaml
version: '3.8'

services:
  myapp:
    build: .
    environment:
      - SLURM_API_URL=http://slurm-host:6820
      - SLURM_API_TOKEN=${SLURM_API_TOKEN}
    ports:
      - "8080:8080"
    networks:
      - slurm-network

networks:
  slurm-network:
    external: true
```

## Kubernetes Deployment

### ConfigMap
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: myapp-config
data:
  config.yaml: |
    slurm:
      api_url: "http://slurm-service:6820"
      api_version: "v0.0.43"
      timeout: "30s"
```

### Secret
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: myapp-secrets
type: Opaque
stringData:
  slurm-token: "your-jwt-token"
```

### Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: myapp
        image: myapp:latest
        ports:
        - containerPort: 8080
        env:
        - name: SLURM_API_TOKEN
          valueFrom:
            secretKeyRef:
              name: myapp-secrets
              key: slurm-token
        volumeMounts:
        - name: config
          mountPath: /app/config.yaml
          subPath: config.yaml
      volumes:
      - name: config
        configMap:
          name: myapp-config
```

## Systemd Service

Create `/etc/systemd/system/myapp.service`:
```ini
[Unit]
Description=My SLURM Application
After=network.target

[Service]
Type=simple
User=myapp
Group=myapp
WorkingDirectory=/opt/myapp
ExecStart=/opt/myapp/myapp
Restart=on-failure
RestartSec=5

# Environment
Environment="SLURM_API_URL=http://slurm-host:6820"
EnvironmentFile=/etc/myapp/myapp.env

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/log/myapp

[Install]
WantedBy=multi-user.target
```

## High Availability

### Load Balancing
```go
// Multiple SLURM API endpoints
endpoints := []string{
    "http://slurm-api-1:6820",
    "http://slurm-api-2:6820",
    "http://slurm-api-3:6820",
}

// Simple round-robin
var currentEndpoint int32
func getNextEndpoint() string {
    idx := atomic.AddInt32(&currentEndpoint, 1)
    return endpoints[int(idx) % len(endpoints)]
}
```

### Circuit Breaker
```go
import "github.com/sony/gobreaker"

cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:        "slurm-api",
    MaxRequests: 3,
    Interval:    time.Minute,
    Timeout:     30 * time.Second,
})

// Wrap API calls
result, err := cb.Execute(func() (interface{}, error) {
    return client.Jobs().List(ctx, opts)
})
```

## Monitoring

### Health Checks
```go
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    if err := client.Info().Ping(ctx); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "unhealthy",
            "error":  err.Error(),
        })
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "status": "healthy",
    })
})
```

### Metrics
```go
import "github.com/prometheus/client_golang/prometheus"

var (
    apiRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "slurm_api_requests_total",
            Help: "Total number of SLURM API requests",
        },
        []string{"method", "endpoint", "status"},
    )

    apiDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "slurm_api_request_duration_seconds",
            Help: "Duration of SLURM API requests",
        },
        []string{"method", "endpoint"},
    )
)
```

## Security Considerations

### Secrets Management
- Never hardcode credentials
- Use environment variables or secret management systems
- Rotate tokens regularly
- Use least privilege principle

### Network Security
- Use TLS/SSL for API connections
- Implement network policies in Kubernetes
- Use firewalls to restrict access

### Authentication
- Prefer token-based authentication
- Implement token refresh logic
- Monitor for authentication failures

## Troubleshooting Deployments

### Common Issues

1. **Connection Timeouts**
   - Check network connectivity
   - Verify firewall rules
   - Increase timeout values

2. **Authentication Failures**
   - Verify credentials are correct
   - Check token expiration
   - Ensure proper secret mounting

3. **Performance Issues**
   - Implement connection pooling
   - Use caching where appropriate
   - Monitor resource usage

### Debugging
Enable debug logging:
```go
if os.Getenv("DEBUG") == "true" {
    log.SetLevel(log.DebugLevel)
}
```

## See Also

- [Configuration Guide](../configuration.md)
- [Troubleshooting Guide](./troubleshooting.md)
- [Examples](../../examples/README.md)