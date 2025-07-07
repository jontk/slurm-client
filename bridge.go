package slurm

import (
	"context"
	
	"github.com/jontk/slurm-client/internal/factory"
)

// clientBridge adapts the factory.SlurmClient to our main SlurmClient interface
type clientBridge struct {
	factoryClient factory.SlurmClient
}

// newClientBridge creates a bridge between factory and main interfaces
func newClientBridge(factoryClient factory.SlurmClient) SlurmClient {
	return &clientBridge{factoryClient: factoryClient}
}

// Version returns the API version
func (c *clientBridge) Version() string {
	return c.factoryClient.Version()
}

// Jobs returns the JobManager
func (c *clientBridge) Jobs() JobManager {
	return &jobManagerBridge{c.factoryClient.Jobs()}
}

// Nodes returns the NodeManager
func (c *clientBridge) Nodes() NodeManager {
	return &nodeManagerBridge{c.factoryClient.Nodes()}
}

// Partitions returns the PartitionManager
func (c *clientBridge) Partitions() PartitionManager {
	return &partitionManagerBridge{c.factoryClient.Partitions()}
}

// Info returns the InfoManager
func (c *clientBridge) Info() InfoManager {
	return &infoManagerBridge{c.factoryClient.Info()}
}

// Close closes the client
func (c *clientBridge) Close() error {
	return c.factoryClient.Close()
}

// Manager bridges to convert between factory and main interfaces

type jobManagerBridge struct {
	factory factory.JobManager
}

func (j *jobManagerBridge) List(ctx context.Context, opts *ListJobsOptions) (*JobList, error) {
	// Convert main options to factory options
	factoryOpts := &factory.ListJobsOptions{
		UserID:    opts.UserID,
		State:     factory.JobState(opts.State),
		Partition: opts.Partition,
		Limit:     opts.Limit,
		Offset:    opts.Offset,
	}
	
	result, err := j.factory.List(ctx, factoryOpts)
	if err != nil {
		return nil, err
	}
	
	// Convert factory result to main types
	jobs := make([]Job, len(result.Jobs))
	for i, factoryJob := range result.Jobs {
		jobs[i] = Job{
			ID:         factoryJob.ID,
			Name:       factoryJob.Name,
			UserID:     factoryJob.UserID,
			State:      JobState(factoryJob.State),
			Partition:  factoryJob.Partition,
			SubmitTime: factoryJob.SubmitTime,
			StartTime:  factoryJob.StartTime,
			EndTime:    factoryJob.EndTime,
			CPUs:       factoryJob.CPUs,
			Memory:     factoryJob.Memory,
			Metadata:   make(map[string]interface{}),
		}
	}
	
	return &JobList{
		Jobs:  jobs,
		Total: result.Total,
	}, nil
}

func (j *jobManagerBridge) Get(ctx context.Context, jobID string) (*Job, error) {
	factoryJob, err := j.factory.Get(ctx, jobID)
	if err != nil {
		return nil, err
	}
	
	return &Job{
		ID:         factoryJob.ID,
		Name:       factoryJob.Name,
		UserID:     factoryJob.UserID,
		State:      JobState(factoryJob.State),
		Partition:  factoryJob.Partition,
		SubmitTime: factoryJob.SubmitTime,
		StartTime:  factoryJob.StartTime,
		EndTime:    factoryJob.EndTime,
		CPUs:       factoryJob.CPUs,
		Memory:     factoryJob.Memory,
		Metadata:   make(map[string]interface{}),
	}, nil
}

func (j *jobManagerBridge) Submit(ctx context.Context, job *JobSubmission) (*JobSubmitResponse, error) {
	factoryJob := &factory.JobSubmission{
		Name:      job.Name,
		Script:    job.Script,
		Partition: job.Partition,
		CPUs:      job.CPUs,
		Memory:    job.Memory,
		TimeLimit: job.TimeLimit,
	}
	
	result, err := j.factory.Submit(ctx, factoryJob)
	if err != nil {
		return nil, err
	}
	
	return &JobSubmitResponse{
		JobID: result.JobID,
	}, nil
}

func (j *jobManagerBridge) Cancel(ctx context.Context, jobID string) error {
	return j.factory.Cancel(ctx, jobID)
}

func (j *jobManagerBridge) Update(ctx context.Context, jobID string, update *JobUpdate) error {
	factoryUpdate := &factory.JobUpdate{
		TimeLimit: update.TimeLimit,
		Priority:  update.Priority,
	}
	return j.factory.Update(ctx, jobID, factoryUpdate)
}

func (j *jobManagerBridge) Steps(ctx context.Context, jobID string) (*JobStepList, error) {
	result, err := j.factory.Steps(ctx, jobID)
	if err != nil {
		return nil, err
	}
	
	steps := make([]JobStep, len(result.Steps))
	for i, factoryStep := range result.Steps {
		steps[i] = JobStep{
			ID:    factoryStep.ID,
			JobID: factoryStep.JobID,
			Name:  factoryStep.Name,
			State: factoryStep.State,
		}
	}
	
	return &JobStepList{Steps: steps}, nil
}

func (j *jobManagerBridge) Watch(ctx context.Context, opts *WatchJobsOptions) (<-chan JobEvent, error) {
	factoryOpts := &factory.WatchJobsOptions{
		UserID: opts.UserID,
		State:  factory.JobState(opts.State),
	}
	
	factoryChan, err := j.factory.Watch(ctx, factoryOpts)
	if err != nil {
		return nil, err
	}
	
	// Convert channel events
	outChan := make(chan JobEvent)
	go func() {
		defer close(outChan)
		for factoryEvent := range factoryChan {
			outChan <- JobEvent{
				Type:     JobEventType(factoryEvent.Type),
				JobID:    factoryEvent.JobID,
				NewState: JobState(factoryEvent.NewState),
			}
		}
	}()
	
	return outChan, nil
}

type nodeManagerBridge struct {
	factory factory.NodeManager
}

func (n *nodeManagerBridge) List(ctx context.Context, opts *ListNodesOptions) (*NodeList, error) {
	factoryOpts := &factory.ListNodesOptions{
		State:     factory.NodeState(opts.State),
		Partition: opts.Partition,
		Features:  opts.Features,
	}
	
	result, err := n.factory.List(ctx, factoryOpts)
	if err != nil {
		return nil, err
	}
	
	nodes := make([]Node, len(result.Nodes))
	for i, factoryNode := range result.Nodes {
		nodes[i] = Node{
			Name:     factoryNode.Name,
			State:    NodeState(factoryNode.State),
			CPUs:     factoryNode.CPUs,
			Metadata: make(map[string]interface{}),
		}
	}
	
	return &NodeList{
		Nodes: nodes,
		Total: result.Total,
	}, nil
}

func (n *nodeManagerBridge) Get(ctx context.Context, nodeName string) (*Node, error) {
	factoryNode, err := n.factory.Get(ctx, nodeName)
	if err != nil {
		return nil, err
	}
	
	return &Node{
		Name:     factoryNode.Name,
		State:    NodeState(factoryNode.State),
		CPUs:     factoryNode.CPUs,
		Metadata: make(map[string]interface{}),
	}, nil
}

func (n *nodeManagerBridge) Update(ctx context.Context, nodeName string, update *NodeUpdate) error {
	var factoryState *factory.NodeState
	if update.State != nil {
		state := factory.NodeState(*update.State)
		factoryState = &state
	}
	
	factoryUpdate := &factory.NodeUpdate{
		State:  factoryState,
		Reason: update.Reason,
	}
	return n.factory.Update(ctx, nodeName, factoryUpdate)
}

func (n *nodeManagerBridge) Drain(ctx context.Context, nodeName string, reason string) error {
	return n.factory.Drain(ctx, nodeName, reason)
}

func (n *nodeManagerBridge) Resume(ctx context.Context, nodeName string) error {
	return n.factory.Resume(ctx, nodeName)
}

type partitionManagerBridge struct {
	factory factory.PartitionManager
}

func (p *partitionManagerBridge) List(ctx context.Context) (*PartitionList, error) {
	result, err := p.factory.List(ctx)
	if err != nil {
		return nil, err
	}
	
	partitions := make([]Partition, len(result.Partitions))
	for i, factoryPartition := range result.Partitions {
		partitions[i] = Partition{
			Name:        factoryPartition.Name,
			State:       factoryPartition.State,
			TotalCPUs:   factoryPartition.TotalCPUs,
			TotalMemory: factoryPartition.TotalMemory,
			Metadata:    make(map[string]interface{}),
		}
	}
	
	return &PartitionList{
		Partitions: partitions,
		Total:      result.Total,
	}, nil
}

func (p *partitionManagerBridge) Get(ctx context.Context, partitionName string) (*Partition, error) {
	factoryPartition, err := p.factory.Get(ctx, partitionName)
	if err != nil {
		return nil, err
	}
	
	return &Partition{
		Name:        factoryPartition.Name,
		State:       factoryPartition.State,
		TotalCPUs:   factoryPartition.TotalCPUs,
		TotalMemory: factoryPartition.TotalMemory,
		Metadata:    make(map[string]interface{}),
	}, nil
}

func (p *partitionManagerBridge) Update(ctx context.Context, partitionName string, update *PartitionUpdate) error {
	factoryUpdate := &factory.PartitionUpdate{
		State: update.State,
	}
	return p.factory.Update(ctx, partitionName, factoryUpdate)
}

type infoManagerBridge struct {
	factory factory.InfoManager
}

func (i *infoManagerBridge) Ping(ctx context.Context) error {
	return i.factory.Ping(ctx)
}

func (i *infoManagerBridge) Version(ctx context.Context) (*VersionInfo, error) {
	result, err := i.factory.Version(ctx)
	if err != nil {
		return nil, err
	}
	
	return &VersionInfo{
		Version:    result.Version,
		APIVersion: result.APIVersion,
	}, nil
}

func (i *infoManagerBridge) Configuration(ctx context.Context) (*ClusterConfig, error) {
	result, err := i.factory.Configuration(ctx)
	if err != nil {
		return nil, err
	}
	
	return &ClusterConfig{
		ClusterName: result.ClusterName,
	}, nil
}

func (i *infoManagerBridge) Statistics(ctx context.Context) (*ClusterStats, error) {
	result, err := i.factory.Statistics(ctx)
	if err != nil {
		return nil, err
	}
	
	return &ClusterStats{
		JobsRunning: result.JobsRunning,
	}, nil
}