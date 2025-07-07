package factory

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	v042 "github.com/jontk/slurm-client/internal/api/v0_0_42"
	"github.com/jontk/slurm-client/internal/versioning"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
	"github.com/jontk/slurm-client/pkg/retry"
)

// ClientFactory creates version-specific Slurm clients
type ClientFactory struct {
	config       *config.Config
	httpClient   *http.Client
	auth         auth.Provider
	retryPolicy  retry.Policy
	baseURL      string
	
	// Version detection cache
	detectedVersion *versioning.APIVersion
	compatibility   *versioning.VersionCompatibilityMatrix
}

// NewClientFactory creates a new client factory
func NewClientFactory(options ...FactoryOption) (*ClientFactory, error) {
	factory := &ClientFactory{
		config: config.NewDefault(),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		retryPolicy:   retry.NewExponentialBackoff(),
		compatibility: versioning.DefaultCompatibilityMatrix(),
	}
	
	for _, option := range options {
		if err := option(factory); err != nil {
			return nil, err
		}
	}
	
	if factory.baseURL == "" {
		factory.baseURL = factory.config.BaseURL
	}
	
	return factory, nil
}

// FactoryOption represents a configuration option for the ClientFactory
type FactoryOption func(*ClientFactory) error

// WithConfig sets the factory configuration
func WithConfig(cfg *config.Config) FactoryOption {
	return func(f *ClientFactory) error {
		f.config = cfg
		return nil
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) FactoryOption {
	return func(f *ClientFactory) error {
		f.httpClient = httpClient
		return nil
	}
}

// WithAuth sets the authentication provider
func WithAuth(auth auth.Provider) FactoryOption {
	return func(f *ClientFactory) error {
		f.auth = auth
		return nil
	}
}

// WithRetryPolicy sets the retry policy
func WithRetryPolicy(policy retry.Policy) FactoryOption {
	return func(f *ClientFactory) error {
		f.retryPolicy = policy
		return nil
	}
}

// WithBaseURL sets the base URL for the Slurm REST API
func WithBaseURL(baseURL string) FactoryOption {
	return func(f *ClientFactory) error {
		f.baseURL = baseURL
		return nil
	}
}

// NewClient creates a new Slurm client with automatic version detection
func (f *ClientFactory) NewClient(ctx context.Context) (SlurmClient, error) {
	return f.NewClientWithVersion(ctx, "")
}

// NewClientWithVersion creates a new Slurm client for a specific version
func (f *ClientFactory) NewClientWithVersion(ctx context.Context, version string) (SlurmClient, error) {
	var targetVersion *versioning.APIVersion
	var err error
	
	if version == "" {
		// Auto-detect version
		targetVersion, err = f.detectVersion(ctx)
		if err != nil {
			// Fallback to stable version
			if f.config.Debug {
				fmt.Printf("Version detection failed, using stable version: %v\n", err)
			}
			targetVersion = versioning.StableVersion()
		}
	} else {
		// Use specified version
		targetVersion, err = versioning.FindBestVersion(version)
		if err != nil {
			return nil, fmt.Errorf("invalid version %s: %w", version, err)
		}
	}
	
	return f.createClient(targetVersion)
}

// NewClientForSlurmVersion creates a client compatible with a specific Slurm version
func (f *ClientFactory) NewClientForSlurmVersion(ctx context.Context, slurmVersion string) (SlurmClient, error) {
	// Find compatible API version for the Slurm version
	var compatibleVersion *versioning.APIVersion
	
	for _, apiVersion := range versioning.SupportedVersions {
		if f.compatibility.IsSlurmVersionSupported(apiVersion.String(), slurmVersion) {
			if compatibleVersion == nil || apiVersion.Compare(compatibleVersion) > 0 {
				compatibleVersion = apiVersion
			}
		}
	}
	
	if compatibleVersion == nil {
		return nil, fmt.Errorf("no compatible API version found for Slurm %s", slurmVersion)
	}
	
	return f.createClient(compatibleVersion)
}

// ListSupportedVersions returns all supported API versions
func (f *ClientFactory) ListSupportedVersions() []*versioning.APIVersion {
	return versioning.SupportedVersions
}

// GetVersionCompatibility returns version compatibility information
func (f *ClientFactory) GetVersionCompatibility() *versioning.VersionCompatibilityMatrix {
	return f.compatibility
}

// detectVersion detects the API version by querying the OpenAPI endpoint
func (f *ClientFactory) detectVersion(ctx context.Context) (*versioning.APIVersion, error) {
	if f.detectedVersion != nil {
		return f.detectedVersion, nil
	}
	
	// Try to get OpenAPI spec to detect version
	req, err := http.NewRequestWithContext(ctx, "GET", f.baseURL+"/openapi/v3", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create version detection request: %w", err)
	}
	
	// Add authentication if available
	if f.auth != nil {
		if err := f.auth.Authenticate(ctx, req); err != nil {
			if f.config.Debug {
				fmt.Printf("Authentication failed during version detection: %v\n", err)
			}
		}
	}
	
	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to detect version: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("version detection failed with status %d", resp.StatusCode)
	}
	
	var openAPISpec struct {
		Info struct {
			Version string `json:"version"`
		} `json:"info"`
		Servers []struct {
			URL string `json:"url"`
		} `json:"servers"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&openAPISpec); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}
	
	// Extract version from server URLs or info
	var detectedVersionStr string
	if openAPISpec.Info.Version != "" {
		detectedVersionStr = openAPISpec.Info.Version
	} else if len(openAPISpec.Servers) > 0 {
		// Try to extract version from server URL
		// Example: /slurm/v0.0.42/ -> v0.0.42
		for _, server := range openAPISpec.Servers {
			if version := extractVersionFromURL(server.URL); version != "" {
				detectedVersionStr = version
				break
			}
		}
	}
	
	if detectedVersionStr == "" {
		return nil, fmt.Errorf("could not determine API version from OpenAPI spec")
	}
	
	version, err := versioning.ParseVersion(detectedVersionStr)
	if err != nil {
		return nil, fmt.Errorf("invalid detected version %s: %w", detectedVersionStr, err)
	}
	
	// Verify this version is supported
	supported := false
	for _, supportedVersion := range versioning.SupportedVersions {
		if version.Compare(supportedVersion) == 0 {
			supported = true
			break
		}
	}
	
	if !supported {
		return nil, fmt.Errorf("detected version %s is not supported", version.String())
	}
	
	f.detectedVersion = version
	return version, nil
}

// createClient creates a version-specific client implementation
func (f *ClientFactory) createClient(version *versioning.APIVersion) (SlurmClient, error) {
	switch version.String() {
	case "v0.0.40":
		return f.createV0_0_40Client()
	case "v0.0.41":
		return f.createV0_0_41Client()
	case "v0.0.42":
		return f.createV0_0_42Client()
	case "v0.0.43":
		return f.createV0_0_43Client()
	default:
		return nil, fmt.Errorf("unsupported API version: %s", version.String())
	}
}

// Version-specific client creation methods (to be implemented with generated code)

func (f *ClientFactory) createV0_0_40Client() (SlurmClient, error) {
	// TODO: Implement with generated v0.0.40 client
	return nil, fmt.Errorf("v0.0.40 client not yet implemented")
}

func (f *ClientFactory) createV0_0_41Client() (SlurmClient, error) {
	// TODO: Implement with generated v0.0.41 client
	return nil, fmt.Errorf("v0.0.41 client not yet implemented")
}

func (f *ClientFactory) createV0_0_42Client() (SlurmClient, error) {
	config := &v042.ClientConfig{
		BaseURL:    f.baseURL,
		HTTPClient: f.httpClient,
		APIKey:     "", // TODO: Extract from f.auth if needed
		Debug:      false, // TODO: Extract from f.config if needed
	}
	
	wrapperClient, err := v042.NewWrapperClient(config)
	if err != nil {
		return nil, err
	}
	
	// Create bridge adapter to convert concrete types to interfaces
	return &v042Bridge{client: wrapperClient}, nil
}

func (f *ClientFactory) createV0_0_43Client() (SlurmClient, error) {
	// TODO: Implement with generated v0.0.43 client
	return nil, fmt.Errorf("v0.0.43 client not yet implemented")
}

// extractVersionFromURL extracts version from a URL like "/slurm/v0.0.42/"
func extractVersionFromURL(url string) string {
	parts := strings.Split(strings.Trim(url, "/"), "/")
	for _, part := range parts {
		if strings.HasPrefix(part, "v") && strings.Count(part, ".") == 2 {
			return part
		}
	}
	return ""
}

// Helper method to create common client configuration
func (f *ClientFactory) createClientConfig(version *versioning.APIVersion) *ClientConfig {
	return &ClientConfig{
		Version:     version,
		BaseURL:     f.baseURL,
		HTTPClient:  f.httpClient,
		Auth:        f.auth,
		RetryPolicy: f.retryPolicy,
		Config:      f.config,
	}
}

// ClientConfig represents common configuration for all client versions
type ClientConfig struct {
	Version     *versioning.APIVersion
	BaseURL     string
	HTTPClient  *http.Client
	Auth        auth.Provider
	RetryPolicy retry.Policy
	Config      *config.Config
}

// v042Bridge adapts the v0.0.42 WrapperClient to implement the factory.SlurmClient interface
type v042Bridge struct {
	client *v042.WrapperClient
}

func (b *v042Bridge) Version() string {
	return b.client.Version()
}

func (b *v042Bridge) Jobs() JobManager {
	return &v042JobManagerBridge{mgr: b.client.Jobs()}
}

func (b *v042Bridge) Nodes() NodeManager {
	return &v042NodeManagerBridge{mgr: b.client.Nodes()}
}

func (b *v042Bridge) Partitions() PartitionManager {
	return &v042PartitionManagerBridge{mgr: b.client.Partitions()}
}

func (b *v042Bridge) Info() InfoManager {
	return &v042InfoManagerBridge{mgr: b.client.Info()}
}

func (b *v042Bridge) Close() error {
	return b.client.Close()
}

// Bridge adapters for managers
type v042JobManagerBridge struct {
	mgr *v042.JobManager
}

func (b *v042JobManagerBridge) List(ctx context.Context, opts *ListJobsOptions) (*JobList, error) {
	// Convert from factory types to v042 types
	v042Opts := &v042.ListJobsOptions{
		UserID:    opts.UserID,
		State:     v042.JobState(opts.State),
		Partition: opts.Partition,
		Limit:     opts.Limit,
		Offset:    opts.Offset,
	}
	
	result, err := b.mgr.List(ctx, v042Opts)
	if err != nil {
		return nil, err
	}
	
	// Convert back to factory types
	jobs := make([]Job, len(result.Jobs))
	for i, job := range result.Jobs {
		jobs[i] = Job{
			ID:         job.ID,
			Name:       job.Name,
			UserID:     job.UserID,
			State:      JobState(job.State),
			Partition:  job.Partition,
			SubmitTime: job.SubmitTime,
			StartTime:  job.StartTime,
			EndTime:    job.EndTime,
			CPUs:       job.CPUs,
			Memory:     job.Memory,
		}
	}
	
	return &JobList{
		Jobs:  jobs,
		Total: result.Total,
	}, nil
}

func (b *v042JobManagerBridge) Get(ctx context.Context, jobID string) (*Job, error) {
	result, err := b.mgr.Get(ctx, jobID)
	if err != nil {
		return nil, err
	}
	
	return &Job{
		ID:         result.ID,
		Name:       result.Name,
		UserID:     result.UserID,
		State:      JobState(result.State),
		Partition:  result.Partition,
		SubmitTime: result.SubmitTime,
		StartTime:  result.StartTime,
		EndTime:    result.EndTime,
		CPUs:       result.CPUs,
		Memory:     result.Memory,
	}, nil
}

func (b *v042JobManagerBridge) Submit(ctx context.Context, job *JobSubmission) (*JobSubmitResponse, error) {
	v042Job := &v042.JobSubmission{
		Name:      job.Name,
		Script:    job.Script,
		Partition: job.Partition,
		CPUs:      job.CPUs,
		Memory:    job.Memory,
		TimeLimit: job.TimeLimit,
	}
	
	result, err := b.mgr.Submit(ctx, v042Job)
	if err != nil {
		return nil, err
	}
	
	return &JobSubmitResponse{
		JobID: result.JobID,
	}, nil
}

func (b *v042JobManagerBridge) Cancel(ctx context.Context, jobID string) error {
	return b.mgr.Cancel(ctx, jobID)
}

func (b *v042JobManagerBridge) Update(ctx context.Context, jobID string, update *JobUpdate) error {
	v042Update := &v042.JobUpdate{
		TimeLimit: update.TimeLimit,
		Priority:  update.Priority,
	}
	return b.mgr.Update(ctx, jobID, v042Update)
}

func (b *v042JobManagerBridge) Steps(ctx context.Context, jobID string) (*JobStepList, error) {
	result, err := b.mgr.Steps(ctx, jobID)
	if err != nil {
		return nil, err
	}
	
	steps := make([]JobStep, len(result.Steps))
	for i, step := range result.Steps {
		steps[i] = JobStep{
			ID:    step.ID,
			JobID: step.JobID,
			Name:  step.Name,
			State: step.State,
		}
	}
	
	return &JobStepList{Steps: steps}, nil
}

func (b *v042JobManagerBridge) Watch(ctx context.Context, opts *WatchJobsOptions) (<-chan JobEvent, error) {
	v042Opts := &v042.WatchJobsOptions{
		UserID: opts.UserID,
		State:  v042.JobState(opts.State),
	}
	
	v042Chan, err := b.mgr.Watch(ctx, v042Opts)
	if err != nil {
		return nil, err
	}
	
	// Convert channel events
	outChan := make(chan JobEvent)
	go func() {
		defer close(outChan)
		for event := range v042Chan {
			outChan <- JobEvent{
				Type:     event.Type,
				JobID:    event.JobID,
				NewState: JobState(event.NewState),
			}
		}
	}()
	
	return outChan, nil
}

// Simplified implementations for other managers
type v042NodeManagerBridge struct {
	mgr *v042.NodeManager
}

func (b *v042NodeManagerBridge) List(ctx context.Context, opts *ListNodesOptions) (*NodeList, error) {
	// Basic implementation - would need full conversion
	return &NodeList{Nodes: []Node{}, Total: 0}, nil
}

func (b *v042NodeManagerBridge) Get(ctx context.Context, nodeName string) (*Node, error) {
	return &Node{}, nil
}

func (b *v042NodeManagerBridge) Update(ctx context.Context, nodeName string, update *NodeUpdate) error {
	return nil
}

func (b *v042NodeManagerBridge) Drain(ctx context.Context, nodeName string, reason string) error {
	return nil
}

func (b *v042NodeManagerBridge) Resume(ctx context.Context, nodeName string) error {
	return nil
}

type v042PartitionManagerBridge struct {
	mgr *v042.PartitionManager
}

func (b *v042PartitionManagerBridge) List(ctx context.Context) (*PartitionList, error) {
	return &PartitionList{Partitions: []Partition{}, Total: 0}, nil
}

func (b *v042PartitionManagerBridge) Get(ctx context.Context, partitionName string) (*Partition, error) {
	return &Partition{}, nil
}

func (b *v042PartitionManagerBridge) Update(ctx context.Context, partitionName string, update *PartitionUpdate) error {
	return nil
}

type v042InfoManagerBridge struct {
	mgr *v042.InfoManager
}

func (b *v042InfoManagerBridge) Ping(ctx context.Context) error {
	return b.mgr.Ping(ctx)
}

func (b *v042InfoManagerBridge) Version(ctx context.Context) (*VersionInfo, error) {
	result, err := b.mgr.Version(ctx)
	if err != nil {
		return nil, err
	}
	
	return &VersionInfo{
		Version:    result.Version,
		APIVersion: result.APIVersion,
	}, nil
}

func (b *v042InfoManagerBridge) Configuration(ctx context.Context) (*ClusterConfig, error) {
	result, err := b.mgr.Configuration(ctx)
	if err != nil {
		return nil, err
	}
	
	return &ClusterConfig{
		ClusterName: result.ClusterName,
	}, nil
}

func (b *v042InfoManagerBridge) Statistics(ctx context.Context) (*ClusterStats, error) {
	result, err := b.mgr.Statistics(ctx)
	if err != nil {
		return nil, err
	}
	
	return &ClusterStats{
		JobsRunning: result.JobsRunning,
	}, nil
}