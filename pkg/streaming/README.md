# SLURM Streaming Package

This package provides WebSocket and Server-Sent Events (SSE) interfaces for real-time SLURM cluster monitoring. It wraps the existing polling-based Watch functionality to provide streaming interfaces suitable for web applications.

## Features

- **WebSocket Support**: Bidirectional real-time communication
- **Server-Sent Events**: Unidirectional real-time updates
- **Multiple Stream Types**: Jobs, Nodes, and Partitions
- **Filtering Support**: All Watch options supported
- **Production Ready**: Proper error handling and connection management

## Quick Start

### WebSocket Server

```go
package main

import (
    "net/http"
    "github.com/jontk/slurm-client"
    "github.com/jontk/slurm-client/pkg/streaming"
)

func main() {
    // Create SLURM client
    client, err := slurm.NewClient(ctx, options...)
    if err != nil {
        log.Fatal(err)
    }

    // Create WebSocket server
    wsServer := streaming.NewWebSocketServer(client)
    
    // Handle WebSocket connections
    http.HandleFunc("/ws", wsServer.HandleWebSocket)
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Server-Sent Events

```go
// Create SSE server
sseServer := streaming.NewSSEServer(client)

// Handle SSE connections
http.HandleFunc("/events", sseServer.HandleSSE)
```

## Usage

### WebSocket Client (JavaScript)

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
    // Request job stream with filtering
    ws.send(JSON.stringify({
        stream: 'jobs',
        options: {
            states: ['RUNNING', 'PENDING'],
            user_id: 'john',
            partition: 'compute'
        }
    }));
};

ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    
    switch (data.type) {
        case 'event':
            console.log('Received event:', data.data);
            break;
        case 'error':
            console.error('Stream error:', data.error);
            break;
        case 'stream_closed':
            console.log('Stream closed');
            break;
    }
};
```

### Server-Sent Events (JavaScript)

```javascript
// Connect to job stream
const eventSource = new EventSource('/events?stream=jobs&states=RUNNING,PENDING');

eventSource.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('Job event:', data);
};

// Listen for specific event types
eventSource.addEventListener('job_event', (event) => {
    const jobData = JSON.parse(event.data);
    console.log('Job changed:', jobData);
});

eventSource.addEventListener('connected', (event) => {
    console.log('Stream connected');
});

eventSource.addEventListener('error', (event) => {
    console.error('Stream error:', JSON.parse(event.data));
});
```

### cURL Examples

```bash
# Stream jobs via SSE
curl -N "http://localhost:8080/events?stream=jobs&states=RUNNING,PENDING"

# Stream nodes with partition filter
curl -N "http://localhost:8080/events?stream=nodes&partition=compute"

# Stream all partitions
curl -N "http://localhost:8080/events?stream=partitions"
```

## Stream Types

### Jobs (`stream=jobs`)

Monitor job state changes in real-time.

**WebSocket Options:**
```json
{
    "stream": "jobs",
    "options": {
        "user_id": "username",
        "states": ["RUNNING", "PENDING"],
        "partition": "compute",
        "job_ids": ["12345", "12346"],
        "exclude_new": false,
        "exclude_completed": false
    }
}
```

**SSE Query Parameters:**
- `user_id`: Filter by user
- `states`: Comma-separated job states
- `partition`: Filter by partition
- `job_ids`: Comma-separated job IDs

### Nodes (`stream=nodes`)

Monitor node state and availability changes.

**WebSocket Options:**
```json
{
    "stream": "nodes",
    "options": {
        "states": ["IDLE", "ALLOCATED"],
        "partition": "compute",
        "features": ["gpu", "ssd"],
        "node_names": ["node001", "node002"]
    }
}
```

**SSE Query Parameters:**
- `states`: Comma-separated node states
- `partition`: Filter by partition
- `features`: Comma-separated features
- `node_names`: Comma-separated node names

### Partitions (`stream=partitions`)

Monitor partition configuration changes.

**WebSocket Options:**
```json
{
    "stream": "partitions",
    "options": {
        "states": ["UP", "DOWN"],
        "partition_names": ["compute", "gpu"]
    }
}
```

**SSE Query Parameters:**
- `states`: Comma-separated partition states
- `partition_names`: Comma-separated partition names

## Message Format

### WebSocket Messages

All WebSocket messages follow this format:

```json
{
    "type": "event|error|stream_closed",
    "stream": "jobs|nodes|partitions",
    "data": {...},
    "timestamp": "2024-01-01T12:00:00Z",
    "error": "error message (if type=error)"
}
```

### SSE Events

SSE events include:
- `event`: Generic event (default)
- `connected`: Stream connection established
- `job_event`: Job-specific event
- `node_event`: Node-specific event
- `partition_event`: Partition-specific event
- `stream_closed`: Stream ended
- `error`: Error occurred

## Production Considerations

### Connection Management

- WebSocket connections include automatic ping/pong for keep-alive
- SSE connections handle client disconnections gracefully
- Both support proper context cancellation

### Error Handling

```go
// Handle WebSocket errors
ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    if (data.type === 'error') {
        console.error('Stream error:', data.error);
        // Implement reconnection logic
    }
};

// Handle SSE errors
eventSource.onerror = (event) => {
    console.error('SSE error:', event);
    // Implement reconnection logic
};
```

### Authentication

The streaming server inherits authentication from the underlying SLURM client. Ensure proper authentication is configured:

```go
client, err := slurm.NewClient(ctx,
    slurm.WithAuth(auth.NewTokenAuth("your-token")),
)
```

### CORS Configuration

For web applications, configure CORS appropriately:

```go
wsServer := streaming.NewWebSocketServer(client)
// The WebSocket upgrader includes basic CORS handling
// Customize the CheckOrigin function for production
```

## Examples

See the [streaming-server example](../../examples/streaming-server/) for a complete implementation with:
- WebSocket and SSE endpoints
- Interactive web interface
- Health check endpoint
- Proper error handling

## Dependencies

- `github.com/gorilla/websocket`: WebSocket implementation
- Standard library for SSE implementation

## Limitations

1. **Polling-Based**: Still uses polling under the hood (SLURM REST API limitation)
2. **Memory Usage**: Maintains connection state for each client
3. **Scalability**: Consider using a message broker for large-scale deployments

## Future Enhancements

- Redis-backed scaling for multiple server instances
- Prometheus metrics for monitoring
- Rate limiting for clients
- Binary WebSocket message support