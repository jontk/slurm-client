package watch

import (
	"context"
	"sync"
	"time"

	"github.com/jontk/slurm-client/internal/interfaces"
)

// DefaultPollInterval is the default polling interval for watch operations
const DefaultPollInterval = 5 * time.Second

// JobPoller implements real-time job monitoring through polling
type JobPoller struct {
	listFunc     func(ctx context.Context, opts *interfaces.ListJobsOptions) (*interfaces.JobList, error)
	pollInterval time.Duration
	bufferSize   int
	mu           sync.RWMutex
	jobStates    map[string]string // Track job states
}

// NewJobPoller creates a new job poller
func NewJobPoller(listFunc func(ctx context.Context, opts *interfaces.ListJobsOptions) (*interfaces.JobList, error)) *JobPoller {
	return &JobPoller{
		listFunc:     listFunc,
		pollInterval: DefaultPollInterval,
		bufferSize:   100,
		jobStates:    make(map[string]string),
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
func (p *JobPoller) Watch(ctx context.Context, opts *interfaces.WatchJobsOptions) (<-chan interfaces.JobEvent, error) {
	// Create event channel
	eventChan := make(chan interfaces.JobEvent, p.bufferSize)

	// Initial state capture
	if opts == nil {
		opts = &interfaces.WatchJobsOptions{}
	}

	// Start polling goroutine
	go p.pollLoop(ctx, opts, eventChan)

	return eventChan, nil
}

// pollLoop is the main polling loop
func (p *JobPoller) pollLoop(ctx context.Context, opts *interfaces.WatchJobsOptions, eventChan chan<- interfaces.JobEvent) {
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
func (p *JobPoller) performPoll(ctx context.Context, opts *interfaces.WatchJobsOptions, eventChan chan<- interfaces.JobEvent, isInitial bool) {
	// Convert watch options to list options
	listOpts := &interfaces.ListJobsOptions{}
	if opts.JobIDs != nil && len(opts.JobIDs) > 0 {
		// For specific job IDs, we'll fetch all jobs and filter
		// (Most SLURM APIs don't support filtering by multiple job IDs directly)
		listOpts.Limit = 0 // No limit
	}
	if opts.States != nil && len(opts.States) > 0 {
		listOpts.States = opts.States
	}

	// Fetch current job list
	jobList, err := p.listFunc(ctx, listOpts)
	if err != nil {
		// Send error event
		eventChan <- interfaces.JobEvent{
			Type:      "error",
			Timestamp: time.Now(),
			Error:     err,
		}
		return
	}

	// Process jobs
	p.mu.Lock()
	defer p.mu.Unlock()

	currentJobs := make(map[string]bool)
	
	for _, job := range jobList.Jobs {
		// Filter by job IDs if specified
		if opts.JobIDs != nil && len(opts.JobIDs) > 0 {
			found := false
			for _, id := range opts.JobIDs {
				if job.ID == id {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		currentJobs[job.ID] = true
		
		previousState, exists := p.jobStates[job.ID]
		
		if !exists {
			// New job detected
			p.jobStates[job.ID] = job.State
			if !isInitial && (!opts.ExcludeNew) {
				eventChan <- interfaces.JobEvent{
					Type:      "job_new",
					JobID:     job.ID,
					NewState:  job.State,
					Timestamp: time.Now(),
					Job:       &job,
				}
			}
		} else if previousState != job.State {
			// State change detected
			p.jobStates[job.ID] = job.State
			eventChan <- interfaces.JobEvent{
				Type:      "job_state_change",
				JobID:     job.ID,
				OldState:  previousState,
				NewState:  job.State,
				Timestamp: time.Now(),
				Job:       &job,
			}
		}
	}

	// Check for completed/removed jobs
	if !opts.ExcludeCompleted {
		for jobID, state := range p.jobStates {
			if !currentJobs[jobID] {
				// Job no longer in list (completed or removed)
				delete(p.jobStates, jobID)
				eventChan <- interfaces.JobEvent{
					Type:      "job_completed",
					JobID:     jobID,
					OldState:  state,
					NewState:  "COMPLETED", // Assume completed if no longer in list
					Timestamp: time.Now(),
				}
			}
		}
	}
}

// NodePoller implements real-time node monitoring through polling
type NodePoller struct {
	listFunc     func(ctx context.Context, opts *interfaces.ListNodesOptions) (*interfaces.NodeList, error)
	pollInterval time.Duration
	bufferSize   int
	mu           sync.RWMutex
	nodeStates   map[string]string // Track node states
}

// NewNodePoller creates a new node poller
func NewNodePoller(listFunc func(ctx context.Context, opts *interfaces.ListNodesOptions) (*interfaces.NodeList, error)) *NodePoller {
	return &NodePoller{
		listFunc:     listFunc,
		pollInterval: DefaultPollInterval,
		bufferSize:   100,
		nodeStates:   make(map[string]string),
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
func (p *NodePoller) Watch(ctx context.Context, opts *interfaces.WatchNodesOptions) (<-chan interfaces.NodeEvent, error) {
	// Create event channel
	eventChan := make(chan interfaces.NodeEvent, p.bufferSize)

	// Initial state capture
	if opts == nil {
		opts = &interfaces.WatchNodesOptions{}
	}

	// Start polling goroutine
	go p.pollLoop(ctx, opts, eventChan)

	return eventChan, nil
}

// pollLoop is the main polling loop for nodes
func (p *NodePoller) pollLoop(ctx context.Context, opts *interfaces.WatchNodesOptions, eventChan chan<- interfaces.NodeEvent) {
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
func (p *NodePoller) performPoll(ctx context.Context, opts *interfaces.WatchNodesOptions, eventChan chan<- interfaces.NodeEvent, isInitial bool) {
	// Convert watch options to list options
	listOpts := &interfaces.ListNodesOptions{}
	if opts.States != nil && len(opts.States) > 0 {
		listOpts.States = opts.States
	}

	// Fetch current node list
	nodeList, err := p.listFunc(ctx, listOpts)
	if err != nil {
		// Send error event
		eventChan <- interfaces.NodeEvent{
			Type:      "error",
			Timestamp: time.Now(),
			Error:     err,
		}
		return
	}

	// Process nodes
	p.mu.Lock()
	defer p.mu.Unlock()

	currentNodes := make(map[string]bool)
	
	for _, node := range nodeList.Nodes {
		// Filter by node names if specified
		if opts.NodeNames != nil && len(opts.NodeNames) > 0 {
			found := false
			for _, name := range opts.NodeNames {
				if node.Name == name {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		currentNodes[node.Name] = true
		
		previousState, exists := p.nodeStates[node.Name]
		
		if !exists {
			// New node detected (unusual but possible)
			p.nodeStates[node.Name] = node.State
			if !isInitial {
				eventChan <- interfaces.NodeEvent{
					Type:      "node_new",
					NodeName:  node.Name,
					NewState:  node.State,
					Timestamp: time.Now(),
					Node:      &node,
				}
			}
		} else if previousState != node.State {
			// State change detected
			p.nodeStates[node.Name] = node.State
			eventChan <- interfaces.NodeEvent{
				Type:      "node_state_change",
				NodeName:  node.Name,
				OldState:  previousState,
				NewState:  node.State,
				Timestamp: time.Now(),
				Node:      &node,
			}
		}
	}
}

// PartitionPoller implements real-time partition monitoring through polling
type PartitionPoller struct {
	listFunc        func(ctx context.Context, opts *interfaces.ListPartitionsOptions) (*interfaces.PartitionList, error)
	pollInterval    time.Duration
	bufferSize      int
	mu              sync.RWMutex
	partitionStates map[string]string // Track partition states
}

// NewPartitionPoller creates a new partition poller
func NewPartitionPoller(listFunc func(ctx context.Context, opts *interfaces.ListPartitionsOptions) (*interfaces.PartitionList, error)) *PartitionPoller {
	return &PartitionPoller{
		listFunc:        listFunc,
		pollInterval:    DefaultPollInterval,
		bufferSize:      100,
		partitionStates: make(map[string]string),
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
func (p *PartitionPoller) Watch(ctx context.Context, opts *interfaces.WatchPartitionsOptions) (<-chan interfaces.PartitionEvent, error) {
	// Create event channel
	eventChan := make(chan interfaces.PartitionEvent, p.bufferSize)

	// Initial state capture
	if opts == nil {
		opts = &interfaces.WatchPartitionsOptions{}
	}

	// Start polling goroutine
	go p.pollLoop(ctx, opts, eventChan)

	return eventChan, nil
}

// pollLoop is the main polling loop for partitions
func (p *PartitionPoller) pollLoop(ctx context.Context, opts *interfaces.WatchPartitionsOptions, eventChan chan<- interfaces.PartitionEvent) {
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
func (p *PartitionPoller) performPoll(ctx context.Context, opts *interfaces.WatchPartitionsOptions, eventChan chan<- interfaces.PartitionEvent, isInitial bool) {
	// Convert watch options to list options
	listOpts := &interfaces.ListPartitionsOptions{}
	if opts.States != nil && len(opts.States) > 0 {
		listOpts.States = opts.States
	}

	// Fetch current partition list
	partitionList, err := p.listFunc(ctx, listOpts)
	if err != nil {
		// Send error event
		eventChan <- interfaces.PartitionEvent{
			Type:      "error",
			Timestamp: time.Now(),
			Error:     err,
		}
		return
	}

	// Process partitions
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, partition := range partitionList.Partitions {
		// Filter by partition names if specified
		if opts.PartitionNames != nil && len(opts.PartitionNames) > 0 {
			found := false
			for _, name := range opts.PartitionNames {
				if partition.Name == name {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		previousState, exists := p.partitionStates[partition.Name]
		
		if !exists {
			// First time seeing this partition
			p.partitionStates[partition.Name] = partition.State
			if !isInitial {
				eventChan <- interfaces.PartitionEvent{
					Type:          "partition_new",
					PartitionName: partition.Name,
					NewState:      partition.State,
					Timestamp:     time.Now(),
					Partition:     &partition,
				}
			}
		} else if previousState != partition.State {
			// State change detected
			p.partitionStates[partition.Name] = partition.State
			eventChan <- interfaces.PartitionEvent{
				Type:          "partition_state_change",
				PartitionName: partition.Name,
				OldState:      previousState,
				NewState:      partition.State,
				Timestamp:     time.Now(),
				Partition:     &partition,
			}
		}
	}
}