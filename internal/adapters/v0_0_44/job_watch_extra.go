// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_44

import (
	"context"
	"time"

	types "github.com/jontk/slurm-client/api"
)

const defaultJobPollInterval = 5 * time.Second

// watchJobsImpl provides the real implementation for job watching using polling.
// This overrides the stub in job_helpers.gen.go - generate_adapters.go must not emit
// the stub when this file exists.
func (a *JobAdapter) watchJobsImpl(ctx context.Context, opts *types.JobWatchOptions) (<-chan types.JobWatchEvent, error) {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}
	// Create event channel
	eventCh := make(chan types.JobWatchEvent, 100)
	// Determine poll interval
	pollInterval := defaultJobPollInterval
	if opts != nil && opts.PollInterval > 0 {
		pollInterval = opts.PollInterval
	}
	// Start polling goroutine
	go a.pollJobs(ctx, opts, eventCh, pollInterval)
	return eventCh, nil
}

// pollJobs polls for job state changes and emits events
func (a *JobAdapter) pollJobs(ctx context.Context, opts *types.JobWatchOptions, eventCh chan<- types.JobWatchEvent, pollInterval time.Duration) {
	defer close(eventCh)
	// Track job states - key is job ID, value is primary state
	jobStates := make(map[int32]types.JobState)
	eventCount := int32(0)
	maxEvents := int32(0)
	if opts != nil {
		maxEvents = opts.MaxEvents
	}
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	// Do initial poll
	a.pollJobsOnce(ctx, opts, eventCh, jobStates, &eventCount, maxEvents)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if maxEvents > 0 && eventCount >= maxEvents {
				return
			}
			a.pollJobsOnce(ctx, opts, eventCh, jobStates, &eventCount, maxEvents)
		}
	}
}

// pollJobsOnce performs a single poll and emits events for state changes
func (a *JobAdapter) pollJobsOnce(
	ctx context.Context,
	opts *types.JobWatchOptions,
	eventCh chan<- types.JobWatchEvent,
	jobStates map[int32]types.JobState,
	eventCount *int32,
	maxEvents int32,
) {
	// Build list options
	listOpts := &types.JobListOptions{}
	if opts != nil && opts.JobId != 0 {
		// If watching a specific job, we use Get instead
		job, err := a.Get(ctx, opts.JobId)
		if err != nil {
			// Job may have been deleted - emit event if we were tracking it
			if prevState, exists := jobStates[opts.JobId]; exists {
				event := types.JobWatchEvent{
					EventTime:     time.Now(),
					EventType:     "deleted",
					JobId:         opts.JobId,
					PreviousState: prevState,
					NewState:      "",
				}
				select {
				case eventCh <- event:
					*eventCount++
					delete(jobStates, opts.JobId)
				case <-ctx.Done():
					return
				}
			}
			return
		}
		a.processJobStatePtr(ctx, job, opts, eventCh, jobStates, eventCount, maxEvents)
		return
	}
	// List all jobs
	result, err := a.List(ctx, listOpts)
	if err != nil {
		return
	}
	// Track which jobs we've seen this poll
	seenJobs := make(map[int32]bool)
	for i := range result.Jobs {
		job := &result.Jobs[i]
		if job.JobID == nil {
			continue
		}
		seenJobs[*job.JobID] = true
		a.processJobStatePtr(ctx, job, opts, eventCh, jobStates, eventCount, maxEvents)
	}
	// Check for deleted jobs
	for jobId, prevState := range jobStates {
		if !seenJobs[jobId] {
			event := types.JobWatchEvent{
				EventTime:     time.Now(),
				EventType:     "deleted",
				JobId:         jobId,
				PreviousState: prevState,
				NewState:      "",
			}
			select {
			case eventCh <- event:
				*eventCount++
				delete(jobStates, jobId)
			case <-ctx.Done():
				return
			}
			if maxEvents > 0 && *eventCount >= maxEvents {
				return
			}
		}
	}
}

// processJobStatePtr checks for state changes and emits events (pointer version)
func (a *JobAdapter) processJobStatePtr(
	ctx context.Context,
	job *types.Job,
	opts *types.JobWatchOptions,
	eventCh chan<- types.JobWatchEvent,
	jobStates map[int32]types.JobState,
	eventCount *int32,
	maxEvents int32,
) {
	if job == nil || job.JobID == nil {
		return
	}
	jobId := *job.JobID
	// Get current state - JobState is a slice, take the first element
	var currentState types.JobState
	if len(job.JobState) > 0 {
		currentState = job.JobState[0]
	}
	// Check if state changed
	prevState, exists := jobStates[jobId]
	if !exists {
		// New job - emit "created" event
		event := types.JobWatchEvent{
			EventTime: time.Now(),
			EventType: "created",
			JobId:     jobId,
			NewState:  currentState,
		}
		if job.Name != nil {
			event.JobName = *job.Name
		}
		if job.UserName != nil {
			event.UserName = *job.UserName
		}
		// Apply event type filter
		if matchesJobEventTypes(event.EventType, opts) {
			select {
			case eventCh <- event:
				*eventCount++
			case <-ctx.Done():
				return
			}
		}
		jobStates[jobId] = currentState
	} else if currentState != prevState {
		// State changed - emit event
		eventType := determineJobEventType(prevState, currentState)
		event := types.JobWatchEvent{
			EventTime:     time.Now(),
			EventType:     eventType,
			JobId:         jobId,
			PreviousState: prevState,
			NewState:      currentState,
		}
		if job.Name != nil {
			event.JobName = *job.Name
		}
		if job.UserName != nil {
			event.UserName = *job.UserName
		}
		// Apply event type filter
		if matchesJobEventTypes(event.EventType, opts) {
			select {
			case eventCh <- event:
				*eventCount++
			case <-ctx.Done():
				return
			}
		}
		jobStates[jobId] = currentState
	}
}

// matchesJobEventTypes checks if the event type matches the filter
func matchesJobEventTypes(eventType string, opts *types.JobWatchOptions) bool {
	if opts == nil || len(opts.EventTypes) == 0 {
		return true
	}
	for _, t := range opts.EventTypes {
		if t == eventType {
			return true
		}
	}
	return false
}

// determineJobEventType determines the event type based on state transition
func determineJobEventType(prev, curr types.JobState) string {
	switch curr {
	case types.JobStatePending:
		return "pending"
	case types.JobStateRunning:
		if prev == types.JobStatePending {
			return "start"
		}
		return "running"
	case types.JobStateCompleted:
		return "end"
	case types.JobStateFailed, types.JobStateTimeout, types.JobStateCancelled, types.JobStateNodeFail:
		return "fail"
	case types.JobStateSuspended:
		return "suspend"
	default:
		return "state_change"
	}
}
