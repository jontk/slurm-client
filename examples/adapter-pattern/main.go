package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"net/http"
	
	"github.com/jontk/slurm-client/internal/adapters/common"
	adapterV0043 "github.com/jontk/slurm-client/internal/adapters/v0_0_43"
	apiV0043 "github.com/jontk/slurm-client/internal/api/v0_0_43"
	"github.com/jontk/slurm-client/internal/common/types"
)

func main() {
	// Get configuration from environment
	host := os.Getenv("SLURM_HOST")
	if host == "" {
		host = "localhost:6820"
	}

	token := os.Getenv("SLURM_TOKEN")
	if token == "" {
		log.Fatal("SLURM_TOKEN environment variable is required")
	}

	// Create the base HTTP client with authentication
	httpClient := &http.Client{
		Transport: &tokenTransport{
			Token: token,
			Base:  http.DefaultTransport,
		},
	}

	// Create version-specific client
	apiClient, err := apiV0043.NewClientWithResponses(
		fmt.Sprintf("http://%s", host),
		apiV0043.WithHTTPClient(httpClient),
	)
	if err != nil {
		log.Fatalf("Failed to create API client: %v", err)
	}

	// Create the adapter - this is the key difference!
	// Instead of using version-specific managers directly, we use the adapter
	adapter := adapterV0043.NewAdapter(apiClient)

	// Now we can work with version-agnostic interfaces
	ctx := context.Background()
	
	// Example 1: List QoS using common types
	fmt.Println("=== Listing QoS ===")
	if err := listQoS(ctx, adapter); err != nil {
		log.Printf("Error listing QoS: %v", err)
	}

	// Example 2: Create a new QoS using common types
	fmt.Println("\n=== Creating QoS ===")
	if err := createQoS(ctx, adapter); err != nil {
		log.Printf("Error creating QoS: %v", err)
	}

	// Example 3: Update QoS using common types
	fmt.Println("\n=== Updating QoS ===")
	if err := updateQoS(ctx, adapter); err != nil {
		log.Printf("Error updating QoS: %v", err)
	}
}

func listQoS(ctx context.Context, adapter common.VersionAdapter) error {
	// Get the QoS manager from the adapter
	qosManager := adapter.GetQoSManager()

	// List QoS with options - using common types!
	opts := &types.QoSListOptions{
		Limit:  10,
		Offset: 0,
	}

	list, err := qosManager.List(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to list QoS: %w", err)
	}

	fmt.Printf("Found %d QoS entries:\n", list.Total)
	for _, qos := range list.QoS {
		fmt.Printf("  - %s (Priority: %d, Usage Factor: %.2f)\n", 
			qos.Name, qos.Priority, qos.UsageFactor)
		
		if qos.Limits != nil {
			if qos.Limits.MaxCPUsPerUser != nil {
				fmt.Printf("    Max CPUs per user: %d\n", *qos.Limits.MaxCPUsPerUser)
			}
			if qos.Limits.MaxJobsPerUser != nil {
				fmt.Printf("    Max jobs per user: %d\n", *qos.Limits.MaxJobsPerUser)
			}
		}
	}

	return nil
}

func createQoS(ctx context.Context, adapter common.VersionAdapter) error {
	qosManager := adapter.GetQoSManager()

	// Create a new QoS using common types
	newQoS := &types.QoSCreate{
		Name:           "example-qos",
		Description:    "Example QoS created via adapter pattern",
		Priority:       100,
		UsageFactor:    1.5,
		UsageThreshold: 0.8,
		Flags:          []string{"DenyOnLimit"},
		PreemptMode:    []string{"cluster"},
		Limits: &types.QoSLimits{
			MaxCPUsPerUser:  intPtr(100),
			MaxJobsPerUser:  intPtr(10),
			MaxNodesPerUser: intPtr(5),
		},
	}

	resp, err := qosManager.Create(ctx, newQoS)
	if err != nil {
		return fmt.Errorf("failed to create QoS: %w", err)
	}

	fmt.Printf("Successfully created QoS: %s\n", resp.QoSName)
	return nil
}

func updateQoS(ctx context.Context, adapter common.VersionAdapter) error {
	qosManager := adapter.GetQoSManager()

	// Update an existing QoS using common types
	update := &types.QoSUpdate{
		Description:    stringPtr("Updated description via adapter"),
		Priority:       intPtr(200),
		UsageFactor:    float64Ptr(2.0),
		UsageThreshold: float64Ptr(0.9),
		Limits: &types.QoSLimits{
			MaxCPUsPerUser: intPtr(200),
			MaxJobsPerUser: intPtr(20),
		},
	}

	err := qosManager.Update(ctx, "example-qos", update)
	if err != nil {
		return fmt.Errorf("failed to update QoS: %w", err)
	}

	fmt.Println("Successfully updated QoS")
	
	// Verify the update
	updated, err := qosManager.Get(ctx, "example-qos")
	if err != nil {
		return fmt.Errorf("failed to get updated QoS: %w", err)
	}

	fmt.Printf("Updated QoS details:\n")
	fmt.Printf("  Name: %s\n", updated.Name)
	fmt.Printf("  Description: %s\n", updated.Description)
	fmt.Printf("  Priority: %d\n", updated.Priority)
	fmt.Printf("  Usage Factor: %.2f\n", updated.UsageFactor)

	return nil
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}

// tokenTransport adds authentication token to requests
type tokenTransport struct {
	Token string
	Base  http.RoundTripper
}

func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-SLURM-USER-TOKEN", t.Token)
	return t.Base.RoundTrip(req)
}