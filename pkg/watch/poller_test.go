package watch_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/watch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock list function for testing
type mockJobLister struct {
	jobs      []interfaces.Job
	err       error
	callCount int
}

func (m *mockJobLister) List(ctx context.Context, opts *interfaces.ListJobsOptions) (*interfaces.JobList, error) {
	m.callCount++
	if m.err != nil {
		return nil, m.err
	}
	return &interfaces.JobList{
		Jobs:  m.jobs,
		Total: len(m.jobs),
	}, nil
}

func TestJobPoller_Watch(t *testing.T) {
	// Create a mock lister
	lister := &mockJobLister{
		jobs: []interfaces.Job{
			{ID: "1", State: "RUNNING", UserID: "1000"},
			{ID: "2", State: "PENDING", UserID: "1000"},
		},
	}

	// Create poller with short interval for testing
	poller := watch.NewJobPoller(lister.List).WithPollInterval(100 * time.Millisecond)

	// Start watching
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventChan, err := poller.Watch(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, eventChan)

	// Wait a bit to let the initial poll complete
	time.Sleep(50 * time.Millisecond)

	// Update job states
	lister.jobs = []interfaces.Job{
		{ID: "1", State: "COMPLETED", UserID: "1000"}, // State changed
		{ID: "2", State: "RUNNING", UserID: "1000"},   // State changed
		{ID: "3", State: "PENDING", UserID: "1001"},   // New job
	}

	// Collect events
	var events []interfaces.JobEvent
	timeout := time.After(500 * time.Millisecond)

	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				t.Fatal("Event channel closed unexpectedly")
			}
			events = append(events, event)
			if len(events) >= 3 { // Expecting 3 events
				goto done
			}
		case <-timeout:
			goto done
		}
	}

done:
	cancel() // Stop the poller

	// Verify we got events
	assert.GreaterOrEqual(t, len(events), 3, "Expected at least 3 events")

	// Verify event types
	stateChangeCount := 0
	newJobCount := 0
	for _, event := range events {
		switch event.Type {
		case "job_state_change":
			stateChangeCount++
		case "job_new":
			newJobCount++
		}
	}

	assert.Equal(t, 2, stateChangeCount, "Expected 2 state change events")
	assert.Equal(t, 1, newJobCount, "Expected 1 new job event")
}

func TestJobPoller_WatchWithFilter(t *testing.T) {
	// Create a mock lister
	lister := &mockJobLister{
		jobs: []interfaces.Job{
			{ID: "1", State: "RUNNING", UserID: "1000"},
			{ID: "2", State: "PENDING", UserID: "1000"},
			{ID: "3", State: "RUNNING", UserID: "1001"},
		},
	}

	// Create poller
	poller := watch.NewJobPoller(lister.List).WithPollInterval(100 * time.Millisecond)

	// Start watching with filter for specific job IDs
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := &interfaces.WatchJobsOptions{
		JobIDs: []string{"1", "2"},
	}

	eventChan, err := poller.Watch(ctx, opts)
	require.NoError(t, err)

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	// Update states
	lister.jobs = []interfaces.Job{
		{ID: "1", State: "COMPLETED", UserID: "1000"}, // State changed
		{ID: "2", State: "RUNNING", UserID: "1000"},    // State changed
		{ID: "3", State: "COMPLETED", UserID: "1001"}, // State changed but filtered out
	}

	// Collect events
	var events []interfaces.JobEvent
	timeout := time.After(300 * time.Millisecond)

	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				t.Fatal("Event channel closed unexpectedly")
			}
			if event.Type == "job_state_change" {
				events = append(events, event)
			}
			if len(events) >= 2 {
				goto done
			}
		case <-timeout:
			goto done
		}
	}

done:
	cancel()

	// Verify we only got events for job 1 and 2
	assert.Equal(t, 2, len(events))
	jobIDs := map[string]bool{}
	for _, event := range events {
		jobIDs[event.JobID] = true
	}
	assert.True(t, jobIDs["1"])
	assert.True(t, jobIDs["2"])
	assert.False(t, jobIDs["3"]) // Should not have events for job 3
}

func TestJobPoller_ErrorHandling(t *testing.T) {
	// Create a mock lister that returns an error
	lister := &mockJobLister{
		err: errors.New("API error"),
	}

	// Create poller
	poller := watch.NewJobPoller(lister.List).WithPollInterval(100 * time.Millisecond)

	// Start watching
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventChan, err := poller.Watch(ctx, nil)
	require.NoError(t, err)

	// Wait for error event
	timeout := time.After(200 * time.Millisecond)
	select {
	case event := <-eventChan:
		assert.Equal(t, "error", event.Type)
		assert.Error(t, event.Error)
		assert.Contains(t, event.Error.Error(), "API error")
	case <-timeout:
		t.Fatal("Timeout waiting for error event")
	}
}

func TestJobPoller_ContextCancellation(t *testing.T) {
	// Create a mock lister
	lister := &mockJobLister{
		jobs: []interfaces.Job{
			{ID: "1", State: "RUNNING"},
		},
	}

	// Create poller
	poller := watch.NewJobPoller(lister.List).WithPollInterval(1 * time.Second)

	// Start watching
	ctx, cancel := context.WithCancel(context.Background())

	eventChan, err := poller.Watch(ctx, nil)
	require.NoError(t, err)

	// Cancel immediately
	cancel()

	// Channel should close quickly
	timeout := time.After(100 * time.Millisecond)
	select {
	case _, ok := <-eventChan:
		assert.False(t, ok, "Channel should be closed")
	case <-timeout:
		t.Fatal("Channel didn't close after context cancellation")
	}
}

func TestNodePoller_Watch(t *testing.T) {
	// Mock node list function
	nodeStates := []interfaces.Node{
		{Name: "node-001", State: "IDLE"},
		{Name: "node-002", State: "ALLOCATED"},
	}

	listFunc := func(ctx context.Context, opts *interfaces.ListNodesOptions) (*interfaces.NodeList, error) {
		return &interfaces.NodeList{
			Nodes: nodeStates,
			Total: len(nodeStates),
		}, nil
	}

	// Create poller
	poller := watch.NewNodePoller(listFunc).WithPollInterval(100 * time.Millisecond)

	// Start watching
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventChan, err := poller.Watch(ctx, nil)
	require.NoError(t, err)

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	// Update node states
	nodeStates[0].State = "DRAINING"

	// Wait for event
	timeout := time.After(200 * time.Millisecond)
	select {
	case event := <-eventChan:
		assert.Equal(t, "node_state_change", event.Type)
		assert.Equal(t, "node-001", event.NodeName)
		assert.Equal(t, "IDLE", event.OldState)
		assert.Equal(t, "DRAINING", event.NewState)
	case <-timeout:
		t.Fatal("Timeout waiting for node event")
	}
}

func TestPartitionPoller_Watch(t *testing.T) {
	// Mock partition list function
	partitionStates := []interfaces.Partition{
		{Name: "gpu", State: "UP"},
		{Name: "cpu", State: "UP"},
	}

	listFunc := func(ctx context.Context, opts *interfaces.ListPartitionsOptions) (*interfaces.PartitionList, error) {
		return &interfaces.PartitionList{
			Partitions: partitionStates,
			Total:      len(partitionStates),
		}, nil
	}

	// Create poller
	poller := watch.NewPartitionPoller(listFunc).WithPollInterval(100 * time.Millisecond)

	// Start watching
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventChan, err := poller.Watch(ctx, nil)
	require.NoError(t, err)

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	// Update partition state
	partitionStates[0].State = "DOWN"

	// Wait for event
	timeout := time.After(200 * time.Millisecond)
	select {
	case event := <-eventChan:
		assert.Equal(t, "partition_state_change", event.Type)
		assert.Equal(t, "gpu", event.PartitionName)
		assert.Equal(t, "UP", event.OldState)
		assert.Equal(t, "DOWN", event.NewState)
	case <-timeout:
		t.Fatal("Timeout waiting for partition event")
	}
}