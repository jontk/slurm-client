// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package streaming

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test NewSSEServer
func TestNewSSEServer(t *testing.T) {
	client := &mockSlurmClient{}
	server := NewSSEServer(client)

	require.NotNil(t, server)
	assert.Equal(t, client, server.client)
}

// Test HandleSSE with missing stream parameter
func TestHandleSSE_MissingStreamParameter(t *testing.T) {
	client := &mockSlurmClient{}
	server := NewSSEServer(client)

	req := httptest.NewRequest(http.MethodGet, "/sse", nil)
	w := httptest.NewRecorder()

	server.HandleSSE(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/event-stream", resp.Header.Get("Content-Type"))

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	assert.Contains(t, bodyStr, "event: error")
	assert.Contains(t, bodyStr, "stream parameter required")
}

// Test HandleSSE with unknown stream type
func TestHandleSSE_UnknownStreamType(t *testing.T) {
	client := &mockSlurmClient{}
	server := NewSSEServer(client)

	req := httptest.NewRequest(http.MethodGet, "/sse?stream=invalid", nil)
	w := httptest.NewRecorder()

	server.HandleSSE(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	assert.Contains(t, bodyStr, "event: error")
	assert.Contains(t, bodyStr, "unknown stream type: invalid")
}

// Test HandleSSE for jobs stream
func TestHandleSSE_JobsStream(t *testing.T) {
	eventChan := make(chan interfaces.JobEvent, 2)

	client := &mockSlurmClient{
		jobs: &mockJobManager{
			watchFunc: func(ctx context.Context, opts *interfaces.WatchJobsOptions) (<-chan interfaces.JobEvent, error) {
				// Send events synchronously before returning
				eventChan <- interfaces.JobEvent{
					Type:      "state_change",
					JobID:     "123",
					OldState:  "PENDING",
					NewState:  "RUNNING",
					Timestamp: time.Now(),
				}
				close(eventChan)
				return eventChan, nil
			},
		},
	}
	server := NewSSEServer(client)

	req := httptest.NewRequest(http.MethodGet, "/sse?stream=jobs&user_id=testuser&partition=default", nil)
	w := httptest.NewRecorder()

	// Use a context with timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)

	// HandleSSE blocks until context is done or channel closes
	server.HandleSSE(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	assert.Contains(t, bodyStr, "event: connected")
	assert.Contains(t, bodyStr, `"stream":"jobs"`)
	assert.Contains(t, bodyStr, "event: job_event")
	assert.Contains(t, bodyStr, `"job_id":"123"`)
}

// Test HandleSSE for jobs stream with Watch error
func TestHandleSSE_JobsStreamError(t *testing.T) {
	client := &mockSlurmClient{
		jobs: &mockJobManager{
			watchFunc: func(ctx context.Context, opts *interfaces.WatchJobsOptions) (<-chan interfaces.JobEvent, error) {
				return nil, fmt.Errorf("watch failed")
			},
		},
	}
	server := NewSSEServer(client)

	req := httptest.NewRequest(http.MethodGet, "/sse?stream=jobs", nil)
	w := httptest.NewRecorder()

	server.HandleSSE(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	assert.Contains(t, bodyStr, "event: error")
	assert.Contains(t, bodyStr, "failed to start job stream")
}

// Test HandleSSE for nodes stream
func TestHandleSSE_NodesStream(t *testing.T) {
	eventChan := make(chan interfaces.NodeEvent, 2)

	client := &mockSlurmClient{
		nodes: &mockNodeManager{
			watchFunc: func(ctx context.Context, opts *interfaces.WatchNodesOptions) (<-chan interfaces.NodeEvent, error) {
				// Send events synchronously before returning
				eventChan <- interfaces.NodeEvent{
					Type:      "state_change",
					NodeName:  "node01",
					OldState:  "IDLE",
					NewState:  "ALLOCATED",
					Timestamp: time.Now(),
				}
				close(eventChan)
				return eventChan, nil
			},
		},
	}
	server := NewSSEServer(client)

	req := httptest.NewRequest(http.MethodGet, "/sse?stream=nodes&partition=gpu", nil)
	w := httptest.NewRecorder()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)

	// HandleSSE blocks until context is done
	server.HandleSSE(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	assert.Contains(t, bodyStr, "event: connected")
	assert.Contains(t, bodyStr, `"stream":"nodes"`)
	assert.Contains(t, bodyStr, "event: node_event")
	assert.Contains(t, bodyStr, `"node_name":"node01"`)
}

// Test HandleSSE for partitions stream
func TestHandleSSE_PartitionsStream(t *testing.T) {
	eventChan := make(chan interfaces.PartitionEvent, 2)

	client := &mockSlurmClient{
		partitions: &mockPartitionManager{
			watchFunc: func(ctx context.Context, opts *interfaces.WatchPartitionsOptions) (<-chan interfaces.PartitionEvent, error) {
				// Send events synchronously before returning
				eventChan <- interfaces.PartitionEvent{
					Type:          "state_change",
					PartitionName: "gpu",
					OldState:      "UP",
					NewState:      "DOWN",
					Timestamp:     time.Now(),
				}
				close(eventChan)
				return eventChan, nil
			},
		},
	}
	server := NewSSEServer(client)

	req := httptest.NewRequest(http.MethodGet, "/sse?stream=partitions", nil)
	w := httptest.NewRecorder()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)

	// HandleSSE blocks until context is done
	server.HandleSSE(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	assert.Contains(t, bodyStr, "event: connected")
	assert.Contains(t, bodyStr, `"stream":"partitions"`)
	assert.Contains(t, bodyStr, "event: partition_event")
	assert.Contains(t, bodyStr, `"partition_name":"gpu"`)
}

// Test context cancellation handling
func TestHandleSSE_ContextCancellation(t *testing.T) {
	eventChan := make(chan interfaces.JobEvent)

	client := &mockSlurmClient{
		jobs: &mockJobManager{
			watchFunc: func(ctx context.Context, opts *interfaces.WatchJobsOptions) (<-chan interfaces.JobEvent, error) {
				return eventChan, nil
			},
		},
	}
	server := NewSSEServer(client)

	req := httptest.NewRequest(http.MethodGet, "/sse?stream=jobs", nil)
	w := httptest.NewRecorder()

	ctx, cancel := context.WithCancel(context.Background())
	req = req.WithContext(ctx)

	done := make(chan bool)
	go func() {
		server.HandleSSE(w, req)
		done <- true
	}()

	// Cancel context immediately
	cancel()

	// Wait for handler to finish
	select {
	case <-done:
		// Success - handler returned
	case <-time.After(2 * time.Second):
		t.Fatal("Handler did not return after context cancellation")
	}
}

// Test stream closed event
func TestHandleSSE_StreamClosedEvent(t *testing.T) {
	eventChan := make(chan interfaces.JobEvent)

	client := &mockSlurmClient{
		jobs: &mockJobManager{
			watchFunc: func(ctx context.Context, opts *interfaces.WatchJobsOptions) (<-chan interfaces.JobEvent, error) {
				// Close channel immediately to trigger stream_closed event
				close(eventChan)
				return eventChan, nil
			},
		},
	}
	server := NewSSEServer(client)

	req := httptest.NewRequest(http.MethodGet, "/sse?stream=jobs", nil)
	w := httptest.NewRecorder()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)

	// HandleSSE blocks until context is done or channel closes
	server.HandleSSE(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	assert.Contains(t, bodyStr, "event: stream_closed")
	assert.Contains(t, bodyStr, `"stream":"jobs"`)
}

// Test parseStringSlice helper function
func TestParseStringSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "single value",
			input:    "value1",
			expected: []string{"value1"},
		},
		{
			name:     "multiple values",
			input:    "value1,value2,value3",
			expected: []string{"value1", "value2", "value3"},
		},
		{
			name:     "values with spaces",
			input:    " value1 , value2 , value3 ",
			expected: []string{"value1", "value2", "value3"},
		},
		{
			name:     "empty values filtered",
			input:    "value1,,value2,  ,value3",
			expected: []string{"value1", "value2", "value3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseStringSlice(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test splitString helper function
func TestSplitString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		separator string
		expected  []string
	}{
		{
			name:      "empty string",
			input:     "",
			separator: ",",
			expected:  nil,
		},
		{
			name:      "single value",
			input:     "value1",
			separator: ",",
			expected:  []string{"value1"},
		},
		{
			name:      "multiple values",
			input:     "value1,value2,value3",
			separator: ",",
			expected:  []string{"value1", "value2", "value3"},
		},
		{
			name:      "empty parts",
			input:     "value1,,value3",
			separator: ",",
			expected:  []string{"value1", "", "value3"},
		},
		{
			name:      "different separator",
			input:     "value1|value2|value3",
			separator: "|",
			expected:  []string{"value1", "value2", "value3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitString(tt.input, tt.separator)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test trimSpace helper function
func TestTrimSpace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no spaces",
			input:    "value",
			expected: "value",
		},
		{
			name:     "leading spaces",
			input:    "  value",
			expected: "value",
		},
		{
			name:     "trailing spaces",
			input:    "value  ",
			expected: "value",
		},
		{
			name:     "both sides",
			input:    "  value  ",
			expected: "value",
		},
		{
			name:     "tabs and newlines",
			input:    "\t\nvalue\r\n",
			expected: "value",
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trimSpace(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test writeSSEEvent
func TestWriteSSEEvent(t *testing.T) {
	tests := []struct {
		name     string
		event    SSEEvent
		expected []string
	}{
		{
			name: "full event",
			event: SSEEvent{
				ID:    "123",
				Event: "test",
				Data:  map[string]string{"key": "value"},
				Retry: 5000,
			},
			expected: []string{"id: 123", "event: test", `data: {"key":"value"}`, "retry: 5000"},
		},
		{
			name: "minimal event",
			event: SSEEvent{
				Data: map[string]string{"status": "ok"},
			},
			expected: []string{`data: {"status":"ok"}`},
		},
		{
			name: "event with ID only",
			event: SSEEvent{
				ID:   "456",
				Data: "simple data",
			},
			expected: []string{"id: 456", `data: "simple data"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			server := &SSEServer{}

			server.writeSSEEvent(w, w, tt.event)

			body := w.Body.String()
			for _, exp := range tt.expected {
				assert.Contains(t, body, exp)
			}
		})
	}
}

// Test SSEEvent JSON marshalling
func TestSSEEvent_JSONMarshalling(t *testing.T) {
	event := SSEEvent{
		ID:    "test-id",
		Event: "test-event",
		Data: map[string]interface{}{
			"key":   "value",
			"count": 42,
		},
		Retry: 1000,
	}

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var decoded SSEEvent
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, event.ID, decoded.ID)
	assert.Equal(t, event.Event, decoded.Event)
	assert.Equal(t, event.Retry, decoded.Retry)
}

// Benchmark tests

func BenchmarkParseStringSlice(b *testing.B) {
	input := "value1,value2,value3,value4,value5"
	b.ResetTimer()
	for range b.N {
		parseStringSlice(input)
	}
}

func BenchmarkSplitString(b *testing.B) {
	input := "value1,value2,value3,value4,value5"
	b.ResetTimer()
	for range b.N {
		splitString(input, ",")
	}
}

func BenchmarkTrimSpace(b *testing.B) {
	input := "  value with spaces  "
	b.ResetTimer()
	for range b.N {
		trimSpace(input)
	}
}

func BenchmarkWriteSSEEvent(b *testing.B) {
	server := &SSEServer{}
	event := SSEEvent{
		ID:    "bench-id",
		Event: "bench-event",
		Data:  map[string]string{"key": "value"},
		Retry: 1000,
	}

	b.ResetTimer()
	for range b.N {
		w := httptest.NewRecorder()
		server.writeSSEEvent(w, w, event)
	}
}

func BenchmarkHandleSSE_JobsStream(b *testing.B) {
	eventChan := make(chan interfaces.JobEvent, 100)

	client := &mockSlurmClient{
		jobs: &mockJobManager{
			watchFunc: func(ctx context.Context, opts *interfaces.WatchJobsOptions) (<-chan interfaces.JobEvent, error) {
				return eventChan, nil
			},
		},
	}
	server := NewSSEServer(client)

	b.ResetTimer()
	for range b.N {
		b.StopTimer()
		req := httptest.NewRequest(http.MethodGet, "/sse?stream=jobs", nil)
		w := httptest.NewRecorder()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		req = req.WithContext(ctx)
		b.StartTimer()

		server.HandleSSE(w, req)
		cancel()
	}
}
