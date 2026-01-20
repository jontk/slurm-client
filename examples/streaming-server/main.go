// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
	"github.com/jontk/slurm-client/pkg/streaming"
)

// Example: Streaming server for real-time SLURM events
func main() {
	// Create configuration
	cfg := config.NewDefault()

	// Try to get config from environment
	if url := os.Getenv("SLURM_REST_URL"); url != "" {
		cfg.BaseURL = url
	} else {
		cfg.BaseURL = "https://cluster.example.com:6820"
	}

	// Create authentication
	var authProvider auth.Provider
	if token := os.Getenv("SLURM_JWT"); token != "" {
		authProvider = auth.NewTokenAuth(token)
	} else {
		log.Println("Warning: No SLURM_JWT environment variable found, using no auth")
		authProvider = auth.NewNoAuth()
	}

	ctx := context.Background()

	// Create SLURM client
	client, err := slurm.NewClient(ctx,
		slurm.WithConfig(cfg),
		slurm.WithAuth(authProvider),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Printf("Connected to SLURM API version: %s\n", client.Version())

	// Create streaming servers
	wsServer := streaming.NewWebSocketServer(client)
	sseServer := streaming.NewSSEServer(client)

	// Set up HTTP routes
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", wsServer.HandleWebSocket)
	http.HandleFunc("/events", sseServer.HandleSSE)
	http.HandleFunc("/health", healthCheck)

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Printf("Starting streaming server on port %s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Printf("  - WebSocket: ws://localhost:%s/ws\n", port)
	fmt.Printf("  - Server-Sent Events: http://localhost:%s/events\n", port)
	fmt.Printf("  - Health Check: http://localhost:%s/health\n", port)
	fmt.Printf("  - Demo Page: http://localhost:%s/\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// serveHome serves the demo page
func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, homeHTML)
}

// healthCheck provides a health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status":"ok","service":"slurm-streaming-server"}`)
}

// homeHTML is a simple demo page
const homeHTML = `<!DOCTYPE html>
<html>
<head>
    <title>SLURM Streaming Demo</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .container { max-width: 800px; margin: 0 auto; }
        .section { margin: 20px 0; padding: 20px; border: 1px solid #ddd; border-radius: 5px; }
        .events { height: 300px; overflow-y: scroll; background: #f5f5f5; padding: 10px; border-radius: 3px; }
        button { padding: 10px 20px; margin: 5px; cursor: pointer; }
        .connected { color: green; }
        .disconnected { color: red; }
        .error { color: red; background: #ffeaea; padding: 5px; margin: 5px 0; }
        .event { margin: 2px 0; padding: 5px; background: white; border-radius: 3px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>SLURM Real-time Streaming Demo</h1>
        
        <div class="section">
            <h2>WebSocket Streaming</h2>
            <p>Status: <span id="ws-status" class="disconnected">Disconnected</span></p>
            <button onclick="connectWebSocket()">Connect WebSocket</button>
            <button onclick="disconnectWebSocket()">Disconnect</button>
            <button onclick="startJobStream()">Stream Jobs</button>
            <button onclick="startNodeStream()">Stream Nodes</button>
            <div id="ws-events" class="events"></div>
        </div>

        <div class="section">
            <h2>Server-Sent Events (SSE)</h2>
            <p>Status: <span id="sse-status" class="disconnected">Disconnected</span></p>
            <button onclick="connectSSE('jobs')">Stream Jobs (SSE)</button>
            <button onclick="connectSSE('nodes')">Stream Nodes (SSE)</button>
            <button onclick="connectSSE('partitions')">Stream Partitions (SSE)</button>
            <button onclick="disconnectSSE()">Disconnect SSE</button>
            <div id="sse-events" class="events"></div>
        </div>

        <div class="section">
            <h2>Usage Examples</h2>
            <h3>WebSocket (JavaScript)</h3>
            <pre><code>
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
    // Request job stream
    ws.send(JSON.stringify({
        stream: 'jobs',
        options: {
            states: ['RUNNING', 'PENDING'],
            user_id: 'john'
        }
    }));
};

ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('Received:', data);
};
            </code></pre>

            <h3>Server-Sent Events (JavaScript)</h3>
            <pre><code>
const eventSource = new EventSource('/events?stream=jobs&states=RUNNING,PENDING');

eventSource.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('Job event:', data);
};

eventSource.addEventListener('job_event', (event) => {
    const jobEvent = JSON.parse(event.data);
    console.log('Job changed:', jobEvent);
});
            </code></pre>

            <h3>cURL Examples</h3>
            <pre><code>
# Stream jobs via SSE
curl -N "http://localhost:8080/events?stream=jobs&states=RUNNING,PENDING"

# Stream nodes with partition filter
curl -N "http://localhost:8080/events?stream=nodes&partition=compute"

# Stream all partitions
curl -N "http://localhost:8080/events?stream=partitions"
            </code></pre>
        </div>
    </div>

    <script>
        let ws = null;
        let eventSource = null;

        function connectWebSocket() {
            if (ws) {
                ws.close();
            }

            ws = new WebSocket('ws://localhost:8080/ws');
            
            ws.onopen = () => {
                document.getElementById('ws-status').textContent = 'Connected';
                document.getElementById('ws-status').className = 'connected';
                addEvent('ws-events', 'WebSocket connected', 'info');
            };

            ws.onclose = () => {
                document.getElementById('ws-status').textContent = 'Disconnected';
                document.getElementById('ws-status').className = 'disconnected';
                addEvent('ws-events', 'WebSocket disconnected', 'info');
            };

            ws.onerror = (error) => {
                addEvent('ws-events', 'WebSocket error: ' + error, 'error');
            };

            ws.onmessage = (event) => {
                const data = JSON.parse(event.data);
                addEvent('ws-events', JSON.stringify(data, null, 2), 'event');
            };
        }

        function disconnectWebSocket() {
            if (ws) {
                ws.close();
                ws = null;
            }
        }

        function startJobStream() {
            if (!ws || ws.readyState !== WebSocket.OPEN) {
                addEvent('ws-events', 'WebSocket not connected', 'error');
                return;
            }

            ws.send(JSON.stringify({
                stream: 'jobs',
                options: {
                    states: ['RUNNING', 'PENDING', 'COMPLETED'],
                    exclude_new: false,
                    exclude_completed: false
                }
            }));
            addEvent('ws-events', 'Requested job stream', 'info');
        }

        function startNodeStream() {
            if (!ws || ws.readyState !== WebSocket.OPEN) {
                addEvent('ws-events', 'WebSocket not connected', 'error');
                return;
            }

            ws.send(JSON.stringify({
                stream: 'nodes',
                options: {
                    states: ['IDLE', 'ALLOCATED', 'DOWN']
                }
            }));
            addEvent('ws-events', 'Requested node stream', 'info');
        }

        function connectSSE(streamType) {
            if (eventSource) {
                eventSource.close();
            }

            const url = '/events?stream=' + streamType;
            eventSource = new EventSource(url);

            eventSource.onopen = () => {
                document.getElementById('sse-status').textContent = 'Connected (' + streamType + ')';
                document.getElementById('sse-status').className = 'connected';
                addEvent('sse-events', 'SSE connected for ' + streamType, 'info');
            };

            eventSource.onerror = () => {
                document.getElementById('sse-status').textContent = 'Error';
                document.getElementById('sse-status').className = 'disconnected';
                addEvent('sse-events', 'SSE connection error', 'error');
            };

            eventSource.onmessage = (event) => {
                const data = JSON.parse(event.data);
                addEvent('sse-events', JSON.stringify(data, null, 2), 'event');
            };

            // Listen for specific event types
            eventSource.addEventListener('connected', (event) => {
                const data = JSON.parse(event.data);
                addEvent('sse-events', 'Stream connected: ' + JSON.stringify(data), 'info');
            });

            eventSource.addEventListener('error', (event) => {
                const data = JSON.parse(event.data);
                addEvent('sse-events', 'Stream error: ' + JSON.stringify(data), 'error');
            });
        }

        function disconnectSSE() {
            if (eventSource) {
                eventSource.close();
                eventSource = null;
                document.getElementById('sse-status').textContent = 'Disconnected';
                document.getElementById('sse-status').className = 'disconnected';
                addEvent('sse-events', 'SSE disconnected', 'info');
            }
        }

        function addEvent(containerId, message, type) {
            const container = document.getElementById(containerId);
            const div = document.createElement('div');
            div.className = type === 'error' ? 'error' : 'event';
            div.innerHTML = '<strong>' + new Date().toLocaleTimeString() + ':</strong> ' + 
                           (typeof message === 'object' ? '<pre>' + message + '</pre>' : message);
            container.appendChild(div);
            container.scrollTop = container.scrollHeight;
        }
    </script>
</body>
</html>`
