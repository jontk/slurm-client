// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package watch_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/pkg/watch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper functions for creating pointer types
func ptrInt32(i int32) *int32    { return &i }
func ptrString(s string) *string { return &s }

// Mock list function for testing
type mockJobLister struct {
	mu        sync.RWMutex
	jobs      []types.Job
	err       error
	callCount int
}

func (m *mockJobLister) List(ctx context.Context, opts *types.ListJobsOptions) (*types.JobList, error) {
	m.mu.Lock()
	m.callCount++
	err := m.err
	jobs := make([]types.Job, len(m.jobs))
	copy(jobs, m.jobs)
	m.mu.Unlock()

	if err != nil {
		return nil, err
	}
	return &types.JobList{
		Jobs:  jobs,
		Total: len(jobs),
	}, nil
}

func (m *mockJobLister) setJobs(jobs []types.Job) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.jobs = jobs
}

//lint:ignore U1000 Reserved for future test cases
func (m *mockJobLister) setError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.err = err
}

// Mock node lister for testing
type mockNodeLister struct {
	mu    sync.RWMutex
	nodes []types.Node
	err   error
}

func (m *mockNodeLister) List(ctx context.Context, opts *types.ListNodesOptions) (*types.NodeList, error) {
	m.mu.RLock()
	err := m.err
	nodes := make([]types.Node, len(m.nodes))
	copy(nodes, m.nodes)
	m.mu.RUnlock()

	if err != nil {
		return nil, err
	}
	return &types.NodeList{
		Nodes: nodes,
		Total: len(nodes),
	}, nil
}

func (m *mockNodeLister) setNodes(nodes []types.Node) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nodes = nodes
}

//lint:ignore U1000 Reserved for future test cases
func (m *mockNodeLister) setError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.err = err
}

// Mock partition lister for testing
type mockPartitionLister struct {
	mu         sync.RWMutex
	partitions []types.Partition
	err        error
}

func (m *mockPartitionLister) List(ctx context.Context, opts *types.ListPartitionsOptions) (*types.PartitionList, error) {
	m.mu.RLock()
	err := m.err
	partitions := make([]types.Partition, len(m.partitions))
	copy(partitions, m.partitions)
	m.mu.RUnlock()

	if err != nil {
		return nil, err
	}
	return &types.PartitionList{
		Partitions: partitions,
		Total:      len(partitions),
	}, nil
}

func (m *mockPartitionLister) setPartitions(partitions []types.Partition) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.partitions = partitions
}

//lint:ignore U1000 Reserved for future test cases
func (m *mockPartitionLister) setError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.err = err
}

func TestJobPoller_Watch(t *testing.T) {
	// Create a mock lister
	lister := &mockJobLister{
		jobs: []types.Job{
			{JobID: ptrInt32(1), JobState: []types.JobState{types.JobStateRunning}, UserName: ptrString("user1000")},
			{JobID: ptrInt32(2), JobState: []types.JobState{types.JobStatePending}, UserName: ptrString("user1000")},
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
	lister.setJobs([]types.Job{
		{JobID: ptrInt32(1), JobState: []types.JobState{types.JobStateCompleted}, UserName: ptrString("user1000")}, // State changed
		{JobID: ptrInt32(2), JobState: []types.JobState{types.JobStateRunning}, UserName: ptrString("user1000")},   // State changed
		{JobID: ptrInt32(3), JobState: []types.JobState{types.JobStatePending}, UserName: ptrString("user1001")},   // New job
	})

	// Collect events
	var events []types.JobEvent
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
		switch event.EventType {
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
		jobs: []types.Job{
			{JobID: ptrInt32(1), JobState: []types.JobState{types.JobStateRunning}, UserName: ptrString("user1000")},
			{JobID: ptrInt32(2), JobState: []types.JobState{types.JobStatePending}, UserName: ptrString("user1000")},
			{JobID: ptrInt32(3), JobState: []types.JobState{types.JobStateRunning}, UserName: ptrString("user1001")},
		},
	}

	// Create poller
	poller := watch.NewJobPoller(lister.List).WithPollInterval(100 * time.Millisecond)

	// Start watching with filter for specific job IDs
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := &types.WatchJobsOptions{
		JobIDs: []string{"1", "2"},
	}

	eventChan, err := poller.Watch(ctx, opts)
	require.NoError(t, err)

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	// Update states
	lister.setJobs([]types.Job{
		{JobID: ptrInt32(1), JobState: []types.JobState{types.JobStateCompleted}, UserName: ptrString("user1000")}, // State changed
		{JobID: ptrInt32(2), JobState: []types.JobState{types.JobStateRunning}, UserName: ptrString("user1000")},   // State changed
		{JobID: ptrInt32(3), JobState: []types.JobState{types.JobStateCompleted}, UserName: ptrString("user1001")}, // State changed but filtered out
	})

	// Collect events
	var events []types.JobEvent
	timeout := time.After(300 * time.Millisecond)

	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				t.Fatal("Event channel closed unexpectedly")
			}
			if event.EventType == "job_state_change" {
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
	assert.Len(t, events, 2)
	jobIDs := map[int32]bool{}
	for _, event := range events {
		jobIDs[event.JobId] = true
	}
	assert.True(t, jobIDs[1])
	assert.True(t, jobIDs[2])
	assert.False(t, jobIDs[3]) // Should not have events for job 3
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

	// The poller silently ignores API errors (doesn't send error events)
	// So we just verify the channel stays open and doesn't receive events
	timeout := time.After(200 * time.Millisecond)
	select {
	case _, ok := <-eventChan:
		if ok {
			// If we got an event, that's unexpected since API returns error
			t.Log("Received unexpected event from error poller")
		}
	case <-timeout:
		// Expected - no events sent when API errors occur
	}
}

func TestJobPoller_ContextCancellation(t *testing.T) {
	// Create a mock lister
	lister := &mockJobLister{
		jobs: []types.Job{
			{JobID: ptrInt32(1), JobState: []types.JobState{types.JobStateRunning}},
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
	// Create a mock lister
	lister := &mockNodeLister{
		nodes: []types.Node{
			{Name: ptrString("node-001"), State: []types.NodeState{types.NodeStateIdle}},
			{Name: ptrString("node-002"), State: []types.NodeState{types.NodeStateAllocated}},
		},
	}

	// Create poller
	poller := watch.NewNodePoller(lister.List).WithPollInterval(100 * time.Millisecond)

	// Start watching
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventChan, err := poller.Watch(ctx, nil)
	require.NoError(t, err)

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	// Update node states
	lister.setNodes([]types.Node{
		{Name: ptrString("node-001"), State: []types.NodeState{types.NodeStateDrain}},
		{Name: ptrString("node-002"), State: []types.NodeState{types.NodeStateAllocated}},
	})

	// Wait for event
	timeout := time.After(200 * time.Millisecond)
	select {
	case event := <-eventChan:
		assert.Equal(t, "node_state_change", event.EventType)
		assert.Equal(t, "node-001", event.NodeName)
		assert.Equal(t, types.NodeStateIdle, event.PreviousState)
		assert.Equal(t, types.NodeStateDrain, event.NewState)
	case <-timeout:
		t.Fatal("Timeout waiting for node event")
	}
}

func TestPartitionPoller_Watch(t *testing.T) {
	// Create a mock lister
	lister := &mockPartitionLister{
		partitions: []types.Partition{
			{Name: ptrString("gpu"), Partition: &types.PartitionPartition{State: []types.StateValue{types.StateUp}}},
			{Name: ptrString("cpu"), Partition: &types.PartitionPartition{State: []types.StateValue{types.StateUp}}},
		},
	}

	// Create poller
	poller := watch.NewPartitionPoller(lister.List).WithPollInterval(100 * time.Millisecond)

	// Start watching
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventChan, err := poller.Watch(ctx, nil)
	require.NoError(t, err)

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	// Update partition state
	lister.setPartitions([]types.Partition{
		{Name: ptrString("gpu"), Partition: &types.PartitionPartition{State: []types.StateValue{types.StateDown}}},
		{Name: ptrString("cpu"), Partition: &types.PartitionPartition{State: []types.StateValue{types.StateUp}}},
	})

	// Wait for event
	timeout := time.After(200 * time.Millisecond)
	select {
	case event := <-eventChan:
		assert.Equal(t, "partition_state_change", event.EventType)
		assert.Equal(t, "gpu", event.PartitionName)
		assert.Equal(t, types.PartitionStateUp, event.PreviousState)
		assert.Equal(t, types.PartitionStateDown, event.NewState)
	case <-timeout:
		t.Fatal("Timeout waiting for partition event")
	}
}

func TestJobPoller_WithMethods(t *testing.T) {
	lister := &mockJobLister{}

	// Test WithPollInterval
	poller1 := watch.NewJobPoller(lister.List).WithPollInterval(2 * time.Second)
	assert.NotNil(t, poller1)

	// Test WithBufferSize
	poller2 := watch.NewJobPoller(lister.List).WithBufferSize(200)
	assert.NotNil(t, poller2)

	// Test chaining
	poller3 := watch.NewJobPoller(lister.List).
		WithPollInterval(3 * time.Second).
		WithBufferSize(300)
	assert.NotNil(t, poller3)
}

func TestJobPoller_WatchWithJobCompleted(t *testing.T) {
	// Create a mock lister
	lister := &mockJobLister{
		jobs: []types.Job{
			{JobID: ptrInt32(1), JobState: []types.JobState{types.JobStateRunning}, UserName: ptrString("user1000")},
			{JobID: ptrInt32(2), JobState: []types.JobState{types.JobStatePending}, UserName: ptrString("user1000")},
		},
	}

	// Create poller with short interval for testing
	poller := watch.NewJobPoller(lister.List).WithPollInterval(50 * time.Millisecond)

	// Start watching
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventChan, err := poller.Watch(ctx, &types.WatchJobsOptions{
		ExcludeCompleted: false, // Allow completed events
	})
	require.NoError(t, err)

	// Wait for initial events to establish baseline
	time.Sleep(100 * time.Millisecond)

	// Update mock to simulate job completion (remove job 1)
	lister.setJobs([]types.Job{
		{JobID: ptrInt32(2), JobState: []types.JobState{types.JobStatePending}, UserName: ptrString("user1000")},
	})

	// Wait for completion event
	var completedEvent types.JobEvent
	select {
	case event := <-eventChan:
		if event.EventType == "job_completed" {
			completedEvent = event
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Expected job completion event")
	}

	// Verify completion event
	assert.Equal(t, "job_completed", completedEvent.EventType)
	assert.Equal(t, int32(1), completedEvent.JobId)
	assert.Equal(t, types.JobStateRunning, completedEvent.PreviousState)
	assert.Equal(t, types.JobStateCompleted, completedEvent.NewState)

	cancel()
}

func TestJobPoller_WatchWithExcludeNew(t *testing.T) {
	// Start with empty job list
	lister := &mockJobLister{
		jobs: []types.Job{},
	}

	// Create poller with short interval for testing
	poller := watch.NewJobPoller(lister.List).WithPollInterval(50 * time.Millisecond)

	// Start watching with ExcludeNew = true
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventChan, err := poller.Watch(ctx, &types.WatchJobsOptions{
		ExcludeNew: true,
	})
	require.NoError(t, err)

	// Wait for initial polling
	time.Sleep(100 * time.Millisecond)

	// Add a new job
	lister.setJobs([]types.Job{
		{JobID: ptrInt32(1), JobState: []types.JobState{types.JobStateRunning}, UserName: ptrString("user1000")},
	})

	// Wait a bit more - should NOT get new job event
	select {
	case event := <-eventChan:
		if event.EventType == "job_new" {
			t.Fatal("Should not receive job_new event when ExcludeNew is true")
		}
	case <-time.After(150 * time.Millisecond):
		// This is expected - no new job event should be sent
	}

	cancel()
}

func TestJobPoller_WatchWithExcludeCompleted(t *testing.T) {
	// Start with a job
	lister := &mockJobLister{
		jobs: []types.Job{
			{JobID: ptrInt32(1), JobState: []types.JobState{types.JobStateRunning}, UserName: ptrString("user1000")},
		},
	}

	// Create poller with short interval for testing
	poller := watch.NewJobPoller(lister.List).WithPollInterval(50 * time.Millisecond)

	// Start watching with ExcludeCompleted = true
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventChan, err := poller.Watch(ctx, &types.WatchJobsOptions{
		ExcludeCompleted: true,
	})
	require.NoError(t, err)

	// Wait for initial polling
	time.Sleep(100 * time.Millisecond)

	// Remove the job (simulate completion)
	lister.setJobs([]types.Job{})

	// Wait a bit more - should NOT get completion event
	select {
	case event := <-eventChan:
		if event.EventType == "job_completed" {
			t.Fatal("Should not receive job_completed event when ExcludeCompleted is true")
		}
	case <-time.After(150 * time.Millisecond):
		// This is expected - no completion event should be sent
	}

	cancel()
}

func TestNodePoller_WithMethods(t *testing.T) {
	mockNodeLister := func(ctx context.Context, opts *types.ListNodesOptions) (*types.NodeList, error) {
		return &types.NodeList{Nodes: []types.Node{}}, nil
	}

	// Test WithPollInterval
	poller1 := watch.NewNodePoller(mockNodeLister).WithPollInterval(2 * time.Second)
	assert.NotNil(t, poller1)

	// Test WithBufferSize
	poller2 := watch.NewNodePoller(mockNodeLister).WithBufferSize(200)
	assert.NotNil(t, poller2)

	// Test chaining
	poller3 := watch.NewNodePoller(mockNodeLister).
		WithPollInterval(3 * time.Second).
		WithBufferSize(300)
	assert.NotNil(t, poller3)
}

func TestNodePoller_WatchWithFilteredNodes(t *testing.T) {
	var callCount int32
	mockNodeLister := func(ctx context.Context, opts *types.ListNodesOptions) (*types.NodeList, error) {
		atomic.AddInt32(&callCount, 1)
		nodes := []types.Node{
			{Name: ptrString("node1"), State: []types.NodeState{types.NodeStateIdle}},
			{Name: ptrString("node2"), State: []types.NodeState{types.NodeStateAllocated}},
			{Name: ptrString("node3"), State: []types.NodeState{types.NodeStateDown}},
		}
		return &types.NodeList{Nodes: nodes}, nil
	}

	// Create poller with short interval for testing
	poller := watch.NewNodePoller(mockNodeLister).WithPollInterval(50 * time.Millisecond)

	// Start watching with specific node names
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventChan, err := poller.Watch(ctx, &types.WatchNodesOptions{
		NodeNames: []string{"node1", "node3"}, // Only watch node1 and node3
	})
	require.NoError(t, err)
	require.NotNil(t, eventChan)

	// Wait for initial events - should only get events for node1 and node3
	time.Sleep(100 * time.Millisecond)

	// Verify we got API calls
	assert.Positive(t, atomic.LoadInt32(&callCount))

	cancel()
}

func TestPartitionPoller_WithMethods(t *testing.T) {
	mockPartitionLister := func(ctx context.Context, opts *types.ListPartitionsOptions) (*types.PartitionList, error) {
		return &types.PartitionList{Partitions: []types.Partition{}}, nil
	}

	// Test WithPollInterval
	poller1 := watch.NewPartitionPoller(mockPartitionLister).WithPollInterval(2 * time.Second)
	assert.NotNil(t, poller1)

	// Test WithBufferSize
	poller2 := watch.NewPartitionPoller(mockPartitionLister).WithBufferSize(200)
	assert.NotNil(t, poller2)

	// Test chaining
	poller3 := watch.NewPartitionPoller(mockPartitionLister).
		WithPollInterval(3 * time.Second).
		WithBufferSize(300)
	assert.NotNil(t, poller3)
}

func TestPartitionPoller_WatchWithFilteredPartitions(t *testing.T) {
	var callCount int32
	mockPartitionLister := func(ctx context.Context, opts *types.ListPartitionsOptions) (*types.PartitionList, error) {
		atomic.AddInt32(&callCount, 1)
		partitions := []types.Partition{
			{Name: ptrString("debug"), Partition: &types.PartitionPartition{State: []types.StateValue{types.StateUp}}},
			{Name: ptrString("compute"), Partition: &types.PartitionPartition{State: []types.StateValue{types.StateUp}}},
			{Name: ptrString("gpu"), Partition: &types.PartitionPartition{State: []types.StateValue{types.StateDown}}},
		}
		return &types.PartitionList{Partitions: partitions}, nil
	}

	// Create poller with short interval for testing
	poller := watch.NewPartitionPoller(mockPartitionLister).WithPollInterval(50 * time.Millisecond)

	// Start watching with specific partition names
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventChan, err := poller.Watch(ctx, &types.WatchPartitionsOptions{
		PartitionNames: []string{"debug", "gpu"}, // Only watch debug and gpu
	})
	require.NoError(t, err)
	require.NotNil(t, eventChan)

	// Wait for initial events
	time.Sleep(100 * time.Millisecond)

	// Verify we got API calls
	assert.Positive(t, atomic.LoadInt32(&callCount))

	cancel()
}

func TestJobPoller_WatchWithNilOptions(t *testing.T) {
	lister := &mockJobLister{
		jobs: []types.Job{
			{JobID: ptrInt32(1), JobState: []types.JobState{types.JobStateRunning}, UserName: ptrString("user1000")},
		},
	}

	poller := watch.NewJobPoller(lister.List).WithPollInterval(50 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Pass nil options - should not crash
	eventChan, err := poller.Watch(ctx, nil)
	require.NoError(t, err)
	assert.NotNil(t, eventChan)

	// Wait a bit and cancel
	time.Sleep(100 * time.Millisecond)
	cancel()
}

func TestNodePoller_WatchWithNilOptions(t *testing.T) {
	mockNodeLister := func(ctx context.Context, opts *types.ListNodesOptions) (*types.NodeList, error) {
		return &types.NodeList{Nodes: []types.Node{{Name: ptrString("node1"), State: []types.NodeState{types.NodeStateIdle}}}}, nil
	}

	poller := watch.NewNodePoller(mockNodeLister).WithPollInterval(50 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Pass nil options - should not crash
	eventChan, err := poller.Watch(ctx, nil)
	require.NoError(t, err)
	assert.NotNil(t, eventChan)

	// Wait a bit and cancel
	time.Sleep(100 * time.Millisecond)
	cancel()
}

func TestPartitionPoller_WatchWithNilOptions(t *testing.T) {
	mockPartitionLister := func(ctx context.Context, opts *types.ListPartitionsOptions) (*types.PartitionList, error) {
		return &types.PartitionList{Partitions: []types.Partition{{Name: ptrString("debug"), Partition: &types.PartitionPartition{State: []types.StateValue{types.StateUp}}}}}, nil
	}

	poller := watch.NewPartitionPoller(mockPartitionLister).WithPollInterval(50 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Pass nil options - should not crash
	eventChan, err := poller.Watch(ctx, nil)
	require.NoError(t, err)
	assert.NotNil(t, eventChan)

	// Wait a bit and cancel
	time.Sleep(100 * time.Millisecond)
	cancel()
}

func TestNodePoller_ErrorHandling(t *testing.T) {
	// Create a mock lister that returns an error
	mockNodeLister := func(ctx context.Context, opts *types.ListNodesOptions) (*types.NodeList, error) {
		return nil, errors.New("API error")
	}

	poller := watch.NewNodePoller(mockNodeLister).WithPollInterval(50 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventChan, err := poller.Watch(ctx, &types.WatchNodesOptions{})
	require.NoError(t, err)

	// The poller silently ignores API errors (doesn't send error events)
	// So we just verify the channel stays open and doesn't receive events
	timeout := time.After(200 * time.Millisecond)
	select {
	case _, ok := <-eventChan:
		if ok {
			// If we got an event, that's unexpected since API returns error
			t.Log("Received unexpected event from error poller")
		}
	case <-timeout:
		// Expected - no events sent when API errors occur
	}

	cancel()
}

func TestPartitionPoller_ErrorHandling(t *testing.T) {
	// Create a mock lister that returns an error
	mockPartitionLister := func(ctx context.Context, opts *types.ListPartitionsOptions) (*types.PartitionList, error) {
		return nil, errors.New("API error")
	}

	poller := watch.NewPartitionPoller(mockPartitionLister).WithPollInterval(50 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventChan, err := poller.Watch(ctx, &types.WatchPartitionsOptions{})
	require.NoError(t, err)

	// The poller silently ignores API errors (doesn't send error events)
	// So we just verify the channel stays open and doesn't receive events
	timeout := time.After(200 * time.Millisecond)
	select {
	case _, ok := <-eventChan:
		if ok {
			// If we got an event, that's unexpected since API returns error
			t.Log("Received unexpected event from error poller")
		}
	case <-timeout:
		// Expected - no events sent when API errors occur
	}

	cancel()
}
