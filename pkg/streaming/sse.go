// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package streaming

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	types "github.com/jontk/slurm-client/api"
)

// SSEServer provides Server-Sent Events interface for SLURM events
// This wraps the existing polling-based Watch functionality
type SSEServer struct {
	client types.SlurmClient
}

// NewSSEServer creates a new Server-Sent Events server
func NewSSEServer(client types.SlurmClient) *SSEServer {
	return &SSEServer{
		client: client,
	}
}

// SSEEvent represents a Server-Sent Event
type SSEEvent struct {
	ID    string      `json:"id,omitempty"`
	Event string      `json:"event,omitempty"`
	Data  interface{} `json:"data"`
	Retry int         `json:"retry,omitempty"`
}

// HandleSSE handles Server-Sent Events connections
func (sse *SSEServer) HandleSSE(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	ctx := r.Context()
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Parse query parameters
	streamType := r.URL.Query().Get("stream")
	if streamType == "" {
		sse.writeSSEEvent(w, flusher, SSEEvent{
			Event: "error",
			Data:  map[string]string{"error": "stream parameter required"},
		})
		return
	}

	// Start streaming based on type
	switch StreamType(streamType) {
	case StreamTypeJobs:
		sse.streamJobsSSE(ctx, w, flusher, r)
	case StreamTypeNodes:
		sse.streamNodesSSE(ctx, w, flusher, r)
	case StreamTypePartitions:
		sse.streamPartitionsSSE(ctx, w, flusher, r)
	default:
		sse.writeSSEEvent(w, flusher, SSEEvent{
			Event: "error",
			Data:  map[string]string{"error": "unknown stream type: " + streamType},
		})
	}
}

// streamJobsSSE streams job events via SSE
func (sse *SSEServer) streamJobsSSE(ctx context.Context, w http.ResponseWriter, flusher http.Flusher, r *http.Request) {
	// Parse options from query parameters
	options := &types.WatchJobsOptions{
		UserID:    r.URL.Query().Get("user_id"),
		Partition: r.URL.Query().Get("partition"),
		States:    parseStringSlice(r.URL.Query().Get("states")),
		JobIDs:    parseStringSlice(r.URL.Query().Get("job_ids")),
	}

	// Start watching jobs
	events, err := sse.client.Jobs().Watch(ctx, options)
	if err != nil {
		sse.writeSSEEvent(w, flusher, SSEEvent{
			Event: "error",
			Data:  map[string]string{"error": "failed to start job stream: " + err.Error()},
		})
		return
	}

	// Send connection established event
	sse.writeSSEEvent(w, flusher, SSEEvent{
		Event: "connected",
		Data:  map[string]string{"stream": "jobs", "status": "connected"},
	})

	// Send events
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-events:
			if !ok {
				sse.writeSSEEvent(w, flusher, SSEEvent{
					Event: "stream_closed",
					Data:  map[string]string{"stream": "jobs", "status": "closed"},
				})
				return
			}

			sse.writeSSEEvent(w, flusher, SSEEvent{
				ID:    fmt.Sprintf("job-%d", time.Now().UnixNano()),
				Event: "job_event",
				Data:  event,
			})
		}
	}
}

// streamNodesSSE streams node events via SSE
func (sse *SSEServer) streamNodesSSE(ctx context.Context, w http.ResponseWriter, flusher http.Flusher, r *http.Request) {
	// Parse options from query parameters
	options := &types.WatchNodesOptions{
		Partition: r.URL.Query().Get("partition"),
		States:    parseStringSlice(r.URL.Query().Get("states")),
		Features:  parseStringSlice(r.URL.Query().Get("features")),
		NodeNames: parseStringSlice(r.URL.Query().Get("node_names")),
	}

	// Start watching nodes
	events, err := sse.client.Nodes().Watch(ctx, options)
	if err != nil {
		sse.writeSSEEvent(w, flusher, SSEEvent{
			Event: "error",
			Data:  map[string]string{"error": "failed to start node stream: " + err.Error()},
		})
		return
	}

	// Send connection established event
	sse.writeSSEEvent(w, flusher, SSEEvent{
		Event: "connected",
		Data:  map[string]string{"stream": "nodes", "status": "connected"},
	})

	// Send events
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-events:
			if !ok {
				sse.writeSSEEvent(w, flusher, SSEEvent{
					Event: "stream_closed",
					Data:  map[string]string{"stream": "nodes", "status": "closed"},
				})
				return
			}

			sse.writeSSEEvent(w, flusher, SSEEvent{
				ID:    fmt.Sprintf("node-%d", time.Now().UnixNano()),
				Event: "node_event",
				Data:  event,
			})
		}
	}
}

// streamPartitionsSSE streams partition events via SSE
func (sse *SSEServer) streamPartitionsSSE(ctx context.Context, w http.ResponseWriter, flusher http.Flusher, r *http.Request) {
	// Parse options from query parameters
	options := &types.WatchPartitionsOptions{
		States:         parseStringSlice(r.URL.Query().Get("states")),
		PartitionNames: parseStringSlice(r.URL.Query().Get("partition_names")),
	}

	// Start watching partitions
	events, err := sse.client.Partitions().Watch(ctx, options)
	if err != nil {
		sse.writeSSEEvent(w, flusher, SSEEvent{
			Event: "error",
			Data:  map[string]string{"error": "failed to start partition stream: " + err.Error()},
		})
		return
	}

	// Send connection established event
	sse.writeSSEEvent(w, flusher, SSEEvent{
		Event: "connected",
		Data:  map[string]string{"stream": "partitions", "status": "connected"},
	})

	// Send events
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-events:
			if !ok {
				sse.writeSSEEvent(w, flusher, SSEEvent{
					Event: "stream_closed",
					Data:  map[string]string{"stream": "partitions", "status": "closed"},
				})
				return
			}

			sse.writeSSEEvent(w, flusher, SSEEvent{
				ID:    fmt.Sprintf("partition-%d", time.Now().UnixNano()),
				Event: "partition_event",
				Data:  event,
			})
		}
	}
}

// writeSSEEvent writes an SSE event to the response
func (sse *SSEServer) writeSSEEvent(w http.ResponseWriter, flusher http.Flusher, event SSEEvent) {
	if event.ID != "" {
		fmt.Fprintf(w, "id: %s\n", event.ID)
	}
	if event.Event != "" {
		fmt.Fprintf(w, "event: %s\n", event.Event)
	}

	data, err := json.Marshal(event.Data)
	if err != nil {
		fmt.Fprintf(w, "data: {\"error\": \"failed to marshal data\"}\n")
	} else {
		fmt.Fprintf(w, "data: %s\n", string(data))
	}

	if event.Retry > 0 {
		fmt.Fprintf(w, "retry: %d\n", event.Retry)
	}

	fmt.Fprintf(w, "\n")
	flusher.Flush()
}

// parseStringSlice parses comma-separated string into slice
func parseStringSlice(s string) []string {
	if s == "" {
		return nil
	}

	var result []string
	for _, item := range splitString(s, ",") {
		if trimmed := trimSpace(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// Helper functions for string processing
func splitString(s, sep string) []string {
	if s == "" {
		return nil
	}
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)

	// Trim leading spaces
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}

	// Trim trailing spaces
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}

	return s[start:end]
}
