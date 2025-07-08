package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

func main() {
	ctx := context.Background()

	fmt.Println("=== Slurm Client Authentication Examples ===")

	// Example 1: No authentication (public endpoints)
	fmt.Println("\n--- Example 1: No Authentication ---")
	demonstrateNoAuth(ctx)

	// Example 2: Token-based authentication
	fmt.Println("\n--- Example 2: Token Authentication ---")
	demonstrateTokenAuth(ctx)

	// Example 3: Basic authentication
	fmt.Println("\n--- Example 3: Basic Authentication ---")
	demonstrateBasicAuth(ctx)

	// Example 4: Environment-based configuration
	fmt.Println("\n--- Example 4: Environment Configuration ---")
	demonstrateEnvironmentAuth(ctx)

	// Example 5: Configuration file authentication
	fmt.Println("\n--- Example 5: Configuration File Auth ---")
	demonstrateConfigFileAuth(ctx)

	// Example 6: Custom authentication provider
	fmt.Println("\n--- Example 6: Custom Authentication ---")
	demonstrateCustomAuth(ctx)
}

func demonstrateNoAuth(ctx context.Context) {
	fmt.Println("Creating client with no authentication...")
	
	// Use none authentication for public endpoints
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(auth.NewNoneAuth()),
		slurm.WithTimeout("10s"),
	)
	if err != nil {
		log.Printf("✗ Failed to create no-auth client: %v", err)
		return
	}

	// Test connection
	info, err := client.GetInfo(ctx)
	if err != nil {
		log.Printf("✗ No-auth connection failed: %v", err)
	} else {
		fmt.Printf("✓ No-auth connection successful! Server: %s\n", info.Version)
	}

	// List public information
	partitions, err := client.ListPartitions(ctx)
	if err != nil {
		log.Printf("✗ Failed to list partitions: %v", err)
	} else {
		fmt.Printf("✓ Listed %d partitions without authentication\n", len(partitions))
	}
}

func demonstrateTokenAuth(ctx context.Context) {
	fmt.Println("Creating client with token authentication...")

	// Get token from environment or use example
	token := os.Getenv("SLURM_TOKEN")
	if token == "" {
		token = "your-api-token-here"
		fmt.Printf("Using example token (set SLURM_TOKEN env var for real token)\n")
	}

	// Create token authentication provider
	tokenAuth := auth.NewTokenAuth(token)
	
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(tokenAuth),
		slurm.WithTimeout("10s"),
	)
	if err != nil {
		log.Printf("✗ Failed to create token-auth client: %v", err)
		return
	}

	// Test authenticated operations
	fmt.Printf("Testing token authentication...\n")
	err = testAuthenticatedOperations(ctx, client)
	if err != nil {
		log.Printf("✗ Token authentication failed: %v", err)
	} else {
		fmt.Printf("✓ Token authentication successful!\n")
	}

	// Demonstrate token refresh (if supported)
	fmt.Printf("Testing token refresh...\n")
	refreshableAuth := &RefreshableTokenAuth{
		currentToken: token,
		refreshURL:   "https://auth-server/refresh",
	}

	refreshClient, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(refreshableAuth),
		slurm.WithTimeout("10s"),
	)
	if err != nil {
		log.Printf("✗ Failed to create refreshable token client: %v", err)
		return
	}

	err = testAuthenticatedOperations(ctx, refreshClient)
	if err != nil {
		log.Printf("✗ Refreshable token auth failed: %v", err)
	} else {
		fmt.Printf("✓ Refreshable token authentication successful!\n")
	}
}

func demonstrateBasicAuth(ctx context.Context) {
	fmt.Println("Creating client with basic authentication...")

	// Get credentials from environment or use examples
	username := os.Getenv("SLURM_USERNAME")
	password := os.Getenv("SLURM_PASSWORD")
	
	if username == "" || password == "" {
		username = "slurmuser"
		password = "slurmpass"
		fmt.Printf("Using example credentials (set SLURM_USERNAME/SLURM_PASSWORD env vars)\n")
	}

	// Create basic authentication provider
	basicAuth := auth.NewBasicAuth(username, password)
	
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(basicAuth),
		slurm.WithTimeout("10s"),
	)
	if err != nil {
		log.Printf("✗ Failed to create basic-auth client: %v", err)
		return
	}

	// Test authenticated operations
	fmt.Printf("Testing basic authentication...\n")
	err = testAuthenticatedOperations(ctx, client)
	if err != nil {
		log.Printf("✗ Basic authentication failed: %v", err)
	} else {
		fmt.Printf("✓ Basic authentication successful!\n")
	}
}

func demonstrateEnvironmentAuth(ctx context.Context) {
	fmt.Println("Creating client from environment variables...")

	// Set example environment variables
	envVars := map[string]string{
		"SLURM_REST_URL":    "https://localhost:6820",
		"SLURM_AUTH_TYPE":   "token",
		"SLURM_TOKEN":       "env-token-example",
		"SLURM_TIMEOUT":     "15s",
		"SLURM_MAX_RETRIES": "3",
		"SLURM_DEBUG":       "true",
	}

	// Temporarily set environment variables
	for key, value := range envVars {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}

	// Create client from environment
	client, err := slurm.NewClientFromEnvironment(ctx)
	if err != nil {
		log.Printf("✗ Failed to create client from environment: %v", err)
		return
	}

	fmt.Printf("✓ Client created from environment variables\n")

	// Test connection
	err = testBasicConnection(ctx, client)
	if err != nil {
		log.Printf("✗ Environment-based client failed: %v", err)
	} else {
		fmt.Printf("✓ Environment-based authentication successful!\n")
	}
}

func demonstrateConfigFileAuth(ctx context.Context) {
	fmt.Println("Creating client from configuration file...")

	// Create example configuration
	cfg := &config.Config{
		BaseURL: "https://localhost:6820",
		Timeout: "20s",
		Auth: config.AuthConfig{
			Type:     "basic",
			Username: "configuser",
			Password: "configpass",
		},
		Retry: config.RetryConfig{
			MaxRetries: 5,
			BaseDelay:  "500ms",
			MaxDelay:   "30s",
		},
		TLS: config.TLSConfig{
			InsecureSkipVerify: true,
		},
		Debug: true,
	}

	// Create client from configuration
	client, err := slurm.NewClientFromConfig(ctx, cfg)
	if err != nil {
		log.Printf("✗ Failed to create client from config: %v", err)
		return
	}

	fmt.Printf("✓ Client created from configuration object\n")

	// Test connection
	err = testBasicConnection(ctx, client)
	if err != nil {
		log.Printf("✗ Config-based client failed: %v", err)
	} else {
		fmt.Printf("✓ Configuration-based authentication successful!\n")
	}

	// Demonstrate loading from file
	fmt.Printf("Demonstrating config file loading...\n")
	configPath := "/tmp/slurm-config.json"
	
	err = cfg.SaveToFile(configPath)
	if err != nil {
		log.Printf("✗ Failed to save config file: %v", err)
		return
	}
	defer os.Remove(configPath)

	loadedCfg, err := config.LoadFromFile(configPath)
	if err != nil {
		log.Printf("✗ Failed to load config file: %v", err)
		return
	}

	fileClient, err := slurm.NewClientFromConfig(ctx, loadedCfg)
	if err != nil {
		log.Printf("✗ Failed to create client from loaded config: %v", err)
		return
	}

	fmt.Printf("✓ Client created from config file: %s\n", configPath)
	
	err = testBasicConnection(ctx, fileClient)
	if err != nil {
		log.Printf("✗ File-based client failed: %v", err)
	} else {
		fmt.Printf("✓ File-based authentication successful!\n")
	}
}

func demonstrateCustomAuth(ctx context.Context) {
	fmt.Println("Creating client with custom authentication provider...")

	// Create custom authentication provider
	customAuth := &CustomAuth{
		apiKey:    "custom-api-key",
		signature: "custom-signature",
	}

	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(customAuth),
		slurm.WithTimeout("10s"),
	)
	if err != nil {
		log.Printf("✗ Failed to create custom-auth client: %v", err)
		return
	}

	fmt.Printf("✓ Client created with custom authentication provider\n")

	// Test custom authentication
	err = testBasicConnection(ctx, client)
	if err != nil {
		log.Printf("✗ Custom authentication failed: %v", err)
	} else {
		fmt.Printf("✓ Custom authentication successful!\n")
	}

	// Demonstrate OAuth-like authentication
	fmt.Printf("Demonstrating OAuth-like authentication...\n")
	
	oauthAuth := &OAuthAuth{
		clientID:     "slurm-client-id",
		clientSecret: "slurm-client-secret",
		tokenURL:     "https://auth-server/oauth/token",
		scope:        "slurm:read slurm:write",
	}

	// Get OAuth token
	err = oauthAuth.authenticate(ctx)
	if err != nil {
		log.Printf("✗ OAuth authentication failed: %v", err)
		return
	}

	oauthClient, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(oauthAuth),
		slurm.WithTimeout("10s"),
	)
	if err != nil {
		log.Printf("✗ Failed to create OAuth client: %v", err)
		return
	}

	err = testBasicConnection(ctx, oauthClient)
	if err != nil {
		log.Printf("✗ OAuth client failed: %v", err)
	} else {
		fmt.Printf("✓ OAuth authentication successful!\n")
	}
}

// Helper functions

func testAuthenticatedOperations(ctx context.Context, client slurm.SlurmClient) error {
	// Test operations that require authentication
	
	// 1. Get server info
	_, err := client.GetInfo(ctx)
	if err != nil {
		return fmt.Errorf("GetInfo failed: %w", err)
	}

	// 2. List jobs (may require auth)
	_, err = client.ListJobs(ctx)
	if err != nil {
		return fmt.Errorf("ListJobs failed: %w", err)
	}

	// 3. Submit a test job (requires auth)
	jobReq := &slurm.JobSubmissionRequest{
		Script: `#!/bin/bash
#SBATCH --job-name=auth-test
#SBATCH --output=auth-test-%j.out
#SBATCH --ntasks=1
#SBATCH --time=00:01:00
echo "Authentication test successful"`,
		Name:      "auth-test",
		Partition: "debug",
		Nodes:     1,
		Tasks:     1,
		Time:      "00:01:00",
	}

	jobID, err := client.SubmitJob(ctx, jobReq)
	if err != nil {
		return fmt.Errorf("SubmitJob failed: %w", err)
	}

	fmt.Printf("  ✓ Test job submitted: %d\n", jobID)

	// 4. Cancel the test job
	err = client.CancelJob(ctx, jobID)
	if err != nil {
		return fmt.Errorf("CancelJob failed: %w", err)
	}

	fmt.Printf("  ✓ Test job cancelled: %d\n", jobID)

	return nil
}

func testBasicConnection(ctx context.Context, client slurm.SlurmClient) error {
	// Test basic connection
	info, err := client.GetInfo(ctx)
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}

	fmt.Printf("  ✓ Connected to Slurm %s\n", info.Version)
	return nil
}

// Custom authentication providers

type RefreshableTokenAuth struct {
	currentToken string
	refreshURL   string
}

func (r *RefreshableTokenAuth) AuthenticateRequest(req *http.Request) error {
	// Add current token to request
	req.Header.Set("Authorization", "Bearer "+r.currentToken)
	return nil
}

func (r *RefreshableTokenAuth) refreshToken(ctx context.Context) error {
	// Simulate token refresh
	fmt.Printf("  ⟳ Refreshing token...\n")
	r.currentToken = "refreshed-token-" + fmt.Sprintf("%d", time.Now().Unix())
	return nil
}

type CustomAuth struct {
	apiKey    string
	signature string
}

func (c *CustomAuth) AuthenticateRequest(req *http.Request) error {
	// Add custom headers
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("X-Signature", c.signature)
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	return nil
}

type OAuthAuth struct {
	clientID     string
	clientSecret string
	tokenURL     string
	scope        string
	accessToken  string
	refreshToken string
}

func (o *OAuthAuth) AuthenticateRequest(req *http.Request) error {
	if o.accessToken == "" {
		return fmt.Errorf("no access token available")
	}
	req.Header.Set("Authorization", "Bearer "+o.accessToken)
	return nil
}

func (o *OAuthAuth) authenticate(ctx context.Context) error {
	// Simulate OAuth token exchange
	fmt.Printf("  ⟳ Performing OAuth authentication...\n")
	o.accessToken = "oauth-access-token-" + fmt.Sprintf("%d", time.Now().Unix())
	o.refreshToken = "oauth-refresh-token-" + fmt.Sprintf("%d", time.Now().Unix())
	return nil
}