// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package streaming

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jontk/slurm-client/interfaces"
)

// WebSocketServer provides a WebSocket interface for SLURM events
// This wraps the existing polling-based Watch functionality
type WebSocketServer struct {
	client   interfaces.SlurmClient
	upgrader websocket.Upgrader
}

// NewWebSocketServer creates a new WebSocket server
func NewWebSocketServer(client interfaces.SlurmClient) *WebSocketServer {
	return &WebSocketServer{
		client: client,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// In production, implement proper origin checking
				return true
			},
		},
	}
}

// StreamType represents the type of stream
type StreamType string

const (
	StreamTypeJobs       StreamType = "jobs"
	StreamTypeNodes      StreamType = "nodes"
	StreamTypePartitions StreamType = "partitions"
)

// StreamMessage represents a message sent over WebSocket
type StreamMessage struct {
	Type      string      `json:"type"`
	Stream    StreamType  `json:"stream"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	Error     string      `json:"error,omitempty"`
}

// StreamRequest represents a client request to start streaming
type StreamRequest struct {
	Stream  StreamType  `json:"stream"`
	Options interface{} `json:"options,omitempty"`
}

// JobStreamOptions for job streaming
type JobStreamOptions struct {
	UserID           string   `json:"user_id,omitempty"`
	States           []string `json:"states,omitempty"`
	Partition        string   `json:"partition,omitempty"`
	JobIDs           []string `json:"job_ids,omitempty"`
	ExcludeNew       bool     `json:"exclude_new,omitempty"`
	ExcludeCompleted bool     `json:"exclude_completed,omitempty"`
}

// NodeStreamOptions for node streaming
type NodeStreamOptions struct {
	States    []string `json:"states,omitempty"`
	Partition string   `json:"partition,omitempty"`
	Features  []string `json:"features,omitempty"`
	NodeNames []string `json:"node_names,omitempty"`
}

// PartitionStreamOptions for partition streaming
type PartitionStreamOptions struct {
	States         []string `json:"states,omitempty"`
	PartitionNames []string `json:"partition_names,omitempty"`
}

// HandleWebSocket handles WebSocket connections
func (ws *WebSocketServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("WebSocket close error: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Handle incoming messages
	go ws.handleIncomingMessages(ctx, conn, cancel)

	// Keep connection alive
	ws.keepAlive(ctx, conn)
}

// handleIncomingMessages processes messages from the client
func (ws *WebSocketServer) handleIncomingMessages(ctx context.Context, conn *websocket.Conn, cancel context.CancelFunc) {
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			var req StreamRequest
			err := conn.ReadJSON(&req)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				return
			}

			// Start streaming based on request
			go ws.handleStreamRequest(ctx, conn, req)
		}
	}
}

// handleStreamRequest starts the appropriate stream
func (ws *WebSocketServer) handleStreamRequest(ctx context.Context, conn *websocket.Conn, req StreamRequest) {
	switch req.Stream {
	case StreamTypeJobs:
		ws.streamJobs(ctx, conn, req.Options)
	case StreamTypeNodes:
		ws.streamNodes(ctx, conn, req.Options)
	case StreamTypePartitions:
		ws.streamPartitions(ctx, conn, req.Options)
	default:
		ws.sendError(conn, "unknown stream type: "+string(req.Stream))
	}
}

// streamJobs streams job events
func (ws *WebSocketServer) streamJobs(ctx context.Context, conn *websocket.Conn, optionsData interface{}) {
	// Convert options
	var options *interfaces.WatchJobsOptions
	if optionsData != nil {
		if optsBytes, err := json.Marshal(optionsData); err == nil {
			var jobOpts JobStreamOptions
			if err := json.Unmarshal(optsBytes, &jobOpts); err == nil {
				options = &interfaces.WatchJobsOptions{
					UserID:           jobOpts.UserID,
					States:           jobOpts.States,
					Partition:        jobOpts.Partition,
					JobIDs:           jobOpts.JobIDs,
					ExcludeNew:       jobOpts.ExcludeNew,
					ExcludeCompleted: jobOpts.ExcludeCompleted,
				}
			}
		}
	}

	// Start watching jobs
	events, err := ws.client.Jobs().Watch(ctx, options)
	if err != nil {
		ws.sendError(conn, "failed to start job stream: "+err.Error())
		return
	}

	// Send events over WebSocket
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-events:
			if !ok {
				ws.sendMessage(conn, StreamMessage{
					Type:      "stream_closed",
					Stream:    StreamTypeJobs,
					Timestamp: time.Now(),
				})
				return
			}

			ws.sendMessage(conn, StreamMessage{
				Type:      "event",
				Stream:    StreamTypeJobs,
				Data:      event,
				Timestamp: time.Now(),
			})
		}
	}
}

// streamNodes streams node events
func (ws *WebSocketServer) streamNodes(ctx context.Context, conn *websocket.Conn, optionsData interface{}) {
	// Convert options
	var options *interfaces.WatchNodesOptions
	if optionsData != nil {
		if optsBytes, err := json.Marshal(optionsData); err == nil {
			var nodeOpts NodeStreamOptions
			if err := json.Unmarshal(optsBytes, &nodeOpts); err == nil {
				options = &interfaces.WatchNodesOptions{
					States:    nodeOpts.States,
					Partition: nodeOpts.Partition,
					Features:  nodeOpts.Features,
					NodeNames: nodeOpts.NodeNames,
				}
			}
		}
	}

	// Start watching nodes
	events, err := ws.client.Nodes().Watch(ctx, options)
	if err != nil {
		ws.sendError(conn, "failed to start node stream: "+err.Error())
		return
	}

	// Send events over WebSocket
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-events:
			if !ok {
				ws.sendMessage(conn, StreamMessage{
					Type:      "stream_closed",
					Stream:    StreamTypeNodes,
					Timestamp: time.Now(),
				})
				return
			}

			ws.sendMessage(conn, StreamMessage{
				Type:      "event",
				Stream:    StreamTypeNodes,
				Data:      event,
				Timestamp: time.Now(),
			})
		}
	}
}

// streamPartitions streams partition events
func (ws *WebSocketServer) streamPartitions(ctx context.Context, conn *websocket.Conn, optionsData interface{}) {
	// Convert options
	var options *interfaces.WatchPartitionsOptions
	if optionsData != nil {
		if optsBytes, err := json.Marshal(optionsData); err == nil {
			var partOpts PartitionStreamOptions
			if err := json.Unmarshal(optsBytes, &partOpts); err == nil {
				options = &interfaces.WatchPartitionsOptions{
					States:         partOpts.States,
					PartitionNames: partOpts.PartitionNames,
				}
			}
		}
	}

	// Start watching partitions
	events, err := ws.client.Partitions().Watch(ctx, options)
	if err != nil {
		ws.sendError(conn, "failed to start partition stream: "+err.Error())
		return
	}

	// Send events over WebSocket
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-events:
			if !ok {
				ws.sendMessage(conn, StreamMessage{
					Type:      "stream_closed",
					Stream:    StreamTypePartitions,
					Timestamp: time.Now(),
				})
				return
			}

			ws.sendMessage(conn, StreamMessage{
				Type:      "event",
				Stream:    StreamTypePartitions,
				Data:      event,
				Timestamp: time.Now(),
			})
		}
	}
}

// sendMessage sends a message over the WebSocket
func (ws *WebSocketServer) sendMessage(conn *websocket.Conn, msg StreamMessage) {
	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("WebSocket write error: %v", err)
	}
}

// sendError sends an error message
func (ws *WebSocketServer) sendError(conn *websocket.Conn, message string) {
	ws.sendMessage(conn, StreamMessage{
		Type:      "error",
		Error:     message,
		Timestamp: time.Now(),
	})
}

// keepAlive maintains the WebSocket connection
func (ws *WebSocketServer) keepAlive(ctx context.Context, conn *websocket.Conn) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("WebSocket ping error: %v", err)
				return
			}
		}
	}
}
