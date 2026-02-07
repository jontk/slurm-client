// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Package watch provides polling-based watch implementations for Slurm resources.
//
// TODO(consolidation): This package and internal/adapters/v0_0_*/job_watch_extra.go
// contain similar polling logic. The adapter implementations are currently used directly
// and support additional features (MaxEvents, EventTypes filtering). Consider:
// 1. Enhancing this package to support those features
// 2. Migrating adapters to use this package
// 3. Removing the adapter-specific implementations
//
// See plan/codex_feedback_8.md R3 for details.
package watch

import (
	"context"
	"fmt"
	"sync"
	"time"

	types "github.com/jontk/slurm-client/api"
)

// DefaultPollInterval is the default polling interval for watch operations
const DefaultPollInterval = 5 * time.Second

// JobPoller implements real-time job monitoring through polling
type JobPoller struct {
	listFunc     func(ctx context.Context, opts *types.ListJobsOptions) (*types.JobList, error)
	pollInterval time.Duration
	bufferSize   int
	mu           sync.RWMutex
	jobStates    map[int32]types.JobState // Track job states by JobId
}

// NewJobPoller creates a new job poller
func NewJobPoller(listFunc func(ctx context.Context, opts *types.ListJobsOptions) (*types.JobList, error)) *JobPoller {
	return &JobPoller{
		listFunc:     listFunc,
		pollInterval: DefaultPollInterval,
		bufferSize:   100,
		jobStates:    make(map[int32]types.JobState),
	}
}

// WithPollInterval sets a custom poll interval
func (p *JobPoller) WithPollInterval(interval time.Duration) *JobPoller {
	p.pollInterval = interval
	return p
}

// WithBufferSize sets a custom buffer size for the event channel
func (p *JobPoller) WithBufferSize(size int) *JobPoller {
	p.bufferSize = size
	return p
}

// Watch starts watching for job state changes
func (p *JobPoller) Watch(ctx context.Context, opts *types.WatchJobsOptions) (<-chan types.JobEvent, error) {
	// Create event channel
	eventChan := make(chan types.JobEvent, p.bufferSize)

	// Initial state capture
	if opts == nil {
		opts = &types.WatchJobsOptions{}
	}

	// Start polling goroutine
	go p.pollLoop(ctx, opts, eventChan)

	return eventChan, nil
}

// pollLoop is the main polling loop
func (p *JobPoller) pollLoop(ctx context.Context, opts *types.WatchJobsOptions, eventChan chan<- types.JobEvent) {
	defer close(eventChan)

	// Create a ticker for polling
	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	// Perform initial poll to establish baseline
	p.performPoll(ctx, opts, eventChan, true)

	// Continue polling until context is cancelled
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.performPoll(ctx, opts, eventChan, false)
		}
	}
}

// performPoll executes a single poll operation
func (p *JobPoller) performPoll(ctx context.Context, opts *types.WatchJobsOptions, eventChan chan<- types.JobEvent, isInitial bool) {
	// Convert watch options to list options
	listOpts := &types.ListJobsOptions{}
	if len(opts.JobIDs) > 0 {
		// For specific job IDs, we'll fetch all jobs and filter
		// (Most SLURM APIs don't support filtering by multiple job IDs directly)
		listOpts.Limit = 0 // No limit
	}
	if len(opts.States) > 0 {
		listOpts.States = opts.States
	}

	// Fetch current job list
	jobList, err := p.listFunc(ctx, listOpts)
	if err != nil {
		// Error occurred - just return (errors not sent as events)
		return
	}

	// Process jobs
	p.mu.Lock()
	defer p.mu.Unlock()

	currentJobs := make(map[int32]bool)

	for _, job := range jobList.Jobs {
		job := job // Create local copy to avoid memory aliasing
		jobID := getJobID(&job)
		jobState := getJobState(&job)

		// Filter by job IDs if specified
		if len(opts.JobIDs) > 0 {
			found := false
			jobIDStr := fmt.Sprintf("%d", jobID)
			for _, id := range opts.JobIDs {
				if jobIDStr == id {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		currentJobs[jobID] = true

		previousState, exists := p.jobStates[jobID]

		if !exists {
			// New job detected
			p.jobStates[jobID] = jobState
			if !isInitial && (!opts.ExcludeNew) {
				jobCopy := job
				eventChan <- types.JobEvent{
					EventType: "job_new",
					JobId:     jobID,
					NewState:  jobState,
					EventTime: time.Now(),
					Job:       &jobCopy,
				}
			}
		} else if previousState != jobState {
			// State change detected
			p.jobStates[jobID] = jobState
			jobCopy := job
			eventChan <- types.JobEvent{
				EventType:     "job_state_change",
				JobId:         jobID,
				PreviousState: previousState,
				NewState:      jobState,
				EventTime:     time.Now(),
				Job:           &jobCopy,
			}
		}
	}

	// Check for completed/removed jobs
	if !opts.ExcludeCompleted {
		for jobID, state := range p.jobStates {
			if !currentJobs[jobID] {
				// Job no longer in list (completed or removed)
				delete(p.jobStates, jobID)
				// Note: NewState is JobState type, need to use appropriate constant
				completedState := types.JobState("COMPLETED")
				eventChan <- types.JobEvent{
					EventType:     "job_completed",
					JobId:         jobID,
					PreviousState: state,
					NewState:      completedState,
					EventTime:     time.Now(),
				}
			}
		}
	}
}

// NodePoller implements real-time node monitoring through polling
type NodePoller struct {
	listFunc     func(ctx context.Context, opts *types.ListNodesOptions) (*types.NodeList, error)
	pollInterval time.Duration
	bufferSize   int
	mu           sync.RWMutex
	nodeStates   map[string]types.NodeState // Track node states by name
}

// NewNodePoller creates a new node poller
func NewNodePoller(listFunc func(ctx context.Context, opts *types.ListNodesOptions) (*types.NodeList, error)) *NodePoller {
	return &NodePoller{
		listFunc:     listFunc,
		pollInterval: DefaultPollInterval,
		bufferSize:   100,
		nodeStates:   make(map[string]types.NodeState),
	}
}

// WithPollInterval sets a custom poll interval
func (p *NodePoller) WithPollInterval(interval time.Duration) *NodePoller {
	p.pollInterval = interval
	return p
}

// WithBufferSize sets a custom buffer size for the event channel
func (p *NodePoller) WithBufferSize(size int) *NodePoller {
	p.bufferSize = size
	return p
}

// Watch starts watching for node state changes
func (p *NodePoller) Watch(ctx context.Context, opts *types.WatchNodesOptions) (<-chan types.NodeEvent, error) {
	// Create event channel
	eventChan := make(chan types.NodeEvent, p.bufferSize)

	// Initial state capture
	if opts == nil {
		opts = &types.WatchNodesOptions{}
	}

	// Start polling goroutine
	go p.pollLoop(ctx, opts, eventChan)

	return eventChan, nil
}

// pollLoop is the main polling loop for nodes
func (p *NodePoller) pollLoop(ctx context.Context, opts *types.WatchNodesOptions, eventChan chan<- types.NodeEvent) {
	defer close(eventChan)

	// Create a ticker for polling
	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	// Perform initial poll to establish baseline
	p.performPoll(ctx, opts, eventChan, true)

	// Continue polling until context is cancelled
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.performPoll(ctx, opts, eventChan, false)
		}
	}
}

// performPoll executes a single poll operation for nodes
func (p *NodePoller) performPoll(ctx context.Context, opts *types.WatchNodesOptions, eventChan chan<- types.NodeEvent, isInitial bool) {
	// Convert watch options to list options
	listOpts := &types.ListNodesOptions{}
	if len(opts.States) > 0 {
		listOpts.States = opts.States
	}

	// Fetch current node list
	nodeList, err := p.listFunc(ctx, listOpts)
	if err != nil {
		// Error occurred - just return (errors not sent as events)
		return
	}

	// Process nodes
	p.mu.Lock()
	defer p.mu.Unlock()

	currentNodes := make(map[string]bool)

	for _, node := range nodeList.Nodes {
		node := node // Create local copy to avoid memory aliasing
		nodeName := getNodeName(&node)
		nodeState := getNodeState(&node)

		// Filter by node names if specified
		if len(opts.NodeNames) > 0 {
			found := false
			for _, name := range opts.NodeNames {
				if nodeName == name {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		currentNodes[nodeName] = true

		previousState, exists := p.nodeStates[nodeName]

		if !exists {
			// New node detected (unusual but possible)
			p.nodeStates[nodeName] = nodeState
			if !isInitial {
				nodeCopy := node
				eventChan <- types.NodeEvent{
					EventType: "node_new",
					NodeName:  nodeName,
					NewState:  nodeState,
					EventTime: time.Now(),
					Node:      &nodeCopy,
				}
			}
		} else if previousState != nodeState {
			// State change detected
			p.nodeStates[nodeName] = nodeState
			nodeCopy := node
			eventChan <- types.NodeEvent{
				EventType:     "node_state_change",
				NodeName:      nodeName,
				PreviousState: previousState,
				NewState:      nodeState,
				EventTime:     time.Now(),
				Node:          &nodeCopy,
			}
		}
	}
}

// PartitionPoller implements real-time partition monitoring through polling
type PartitionPoller struct {
	listFunc        func(ctx context.Context, opts *types.ListPartitionsOptions) (*types.PartitionList, error)
	pollInterval    time.Duration
	bufferSize      int
	mu              sync.RWMutex
	partitionStates map[string]types.PartitionState // Track partition states by name
}

// NewPartitionPoller creates a new partition poller
func NewPartitionPoller(listFunc func(ctx context.Context, opts *types.ListPartitionsOptions) (*types.PartitionList, error)) *PartitionPoller {
	return &PartitionPoller{
		listFunc:        listFunc,
		pollInterval:    DefaultPollInterval,
		bufferSize:      100,
		partitionStates: make(map[string]types.PartitionState),
	}
}

// WithPollInterval sets a custom poll interval
func (p *PartitionPoller) WithPollInterval(interval time.Duration) *PartitionPoller {
	p.pollInterval = interval
	return p
}

// WithBufferSize sets a custom buffer size for the event channel
func (p *PartitionPoller) WithBufferSize(size int) *PartitionPoller {
	p.bufferSize = size
	return p
}

// Watch starts watching for partition state changes
func (p *PartitionPoller) Watch(ctx context.Context, opts *types.WatchPartitionsOptions) (<-chan types.PartitionEvent, error) {
	// Create event channel
	eventChan := make(chan types.PartitionEvent, p.bufferSize)

	// Initial state capture
	if opts == nil {
		opts = &types.WatchPartitionsOptions{}
	}

	// Start polling goroutine
	go p.pollLoop(ctx, opts, eventChan)

	return eventChan, nil
}

// pollLoop is the main polling loop for partitions
func (p *PartitionPoller) pollLoop(ctx context.Context, opts *types.WatchPartitionsOptions, eventChan chan<- types.PartitionEvent) {
	defer close(eventChan)

	// Create a ticker for polling
	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	// Perform initial poll to establish baseline
	p.performPoll(ctx, opts, eventChan, true)

	// Continue polling until context is cancelled
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.performPoll(ctx, opts, eventChan, false)
		}
	}
}

// performPoll executes a single poll operation for partitions
func (p *PartitionPoller) performPoll(ctx context.Context, opts *types.WatchPartitionsOptions, eventChan chan<- types.PartitionEvent, isInitial bool) {
	// Convert watch options to list options
	listOpts := &types.ListPartitionsOptions{}
	if len(opts.States) > 0 {
		listOpts.States = opts.States
	}

	// Fetch current partition list
	partitionList, err := p.listFunc(ctx, listOpts)
	if err != nil {
		// Error occurred - just return (errors not sent as events)
		return
	}

	// Process partitions
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, partition := range partitionList.Partitions {
		partition := partition // Create local copy to avoid memory aliasing
		partitionName := getPartitionName(&partition)
		partitionState := getPartitionState(&partition)

		// Filter by partition names if specified
		if len(opts.PartitionNames) > 0 {
			found := false
			for _, name := range opts.PartitionNames {
				if partitionName == name {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		previousState, exists := p.partitionStates[partitionName]

		if !exists {
			// First time seeing this partition
			p.partitionStates[partitionName] = partitionState
			if !isInitial {
				partitionCopy := partition
				eventChan <- types.PartitionEvent{
					EventType:     "partition_new",
					PartitionName: partitionName,
					NewState:      partitionState,
					EventTime:     time.Now(),
					Partition:     &partitionCopy,
				}
			}
		} else if previousState != partitionState {
			// State change detected
			p.partitionStates[partitionName] = partitionState
			partitionCopy := partition
			eventChan <- types.PartitionEvent{
				EventType:     "partition_state_change",
				PartitionName: partitionName,
				PreviousState: previousState,
				NewState:      partitionState,
				EventTime:     time.Now(),
				Partition:     &partitionCopy,
			}
		}
	}
}
