// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// This example demonstrates how to connect to a SLURM REST API server
// with various authentication methods and connection options.
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/pkg/auth"
)

func main() {
	ctx := context.Background()

	// Example 1: Basic connection with JWT token
	fmt.Println("=== Example 1: Basic JWT Authentication ===")
	basicExample(ctx)

	// Example 2: Connection with environment variables
	fmt.Println("\n=== Example 2: Environment Variables ===")
	envExample(ctx)

	// Example 3: Advanced connection options
	fmt.Println("\n=== Example 3: Advanced Options ===")
	advancedExample(ctx)

	// Example 4: Multiple authentication methods
	fmt.Println("\n=== Example 4: Authentication Methods ===")
	authExample(ctx)
}

func basicExample(ctx context.Context) {
	// Create a basic client with user token authentication
	// IMPORTANT: Most SLURM deployments require both username and token
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithUserToken("your-username", "your-jwt-token"),
	)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
	defer client.Close()

	// Test the connection
	err = client.Info().Ping(ctx)
	if err != nil {
		log.Printf("Failed to ping server: %v", err)
		return
	}

	// Get cluster information
	info, err := client.Info().Get(ctx)
	if err != nil {
		log.Printf("Failed to get cluster info: %v", err)
		return
	}

	fmt.Printf("Connected to SLURM REST API\n")
	fmt.Printf("Cluster Name: %s\n", info.ClusterName)
	fmt.Printf("SLURM Version: %s\n", info.Version)
	fmt.Printf("SLURM Release: %s\n", info.Release)

	// Get version info
	versionInfo, err := client.Info().Version(ctx)
	if err != nil {
		log.Printf("Failed to get version info: %v", err)
		return
	}
	fmt.Printf("API Version: %s\n", versionInfo.Version)
	fmt.Printf("API Release: %s\n", versionInfo.Release)
}

func envExample(ctx context.Context) {
	// Set environment variables (normally these would be set externally)
	os.Setenv("SLURM_API_URL", "https://localhost:6820")
	os.Setenv("SLURM_API_TOKEN", "your-jwt-token")
	os.Setenv("SLURM_API_VERSION", "v0.0.43")

	// Create client using environment variables
	client, err := slurm.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to create client from env: %v", err)
		return
	}
	defer client.Close()

	fmt.Println("Successfully created client from environment variables")
}

func advancedExample(ctx context.Context) {
	// Create client with advanced options
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithUserToken("your-username", "your-jwt-token"),

		// Set custom timeout
		slurm.WithTimeout(30*time.Second),

		// Configure retry behavior
		slurm.WithMaxRetries(5),

		// Custom HTTP client with TLS configuration
		slurm.WithHTTPClient(&http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS12,
					// #nosec G402 - codeql[go/insecure-default-tls]: Only for development/testing examples
					InsecureSkipVerify: true,
				},
			},
		}),

		// Force specific API version
		slurm.WithVersion("v0.0.43"),

		// Custom user agent
		slurm.WithUserAgent("my-app/1.0"),
	)
	if err != nil {
		log.Printf("Failed to create advanced client: %v", err)
		return
	}
	defer client.Close()

	fmt.Println("Successfully created client with advanced options")
}

func authExample(ctx context.Context) {
	baseURL := "https://localhost:6820"

	// User Token Authentication (RECOMMENDED for most SLURM deployments)
	// Sets both X-SLURM-USER-NAME and X-SLURM-USER-TOKEN headers
	userTokenClient, err := slurm.NewClient(ctx,
		slurm.WithBaseURL(baseURL),
		slurm.WithUserToken("username", "user-token"),
	)
	if err == nil {
		defer userTokenClient.Close()
		fmt.Println("✓ User Token authentication configured (recommended)")
	}

	// JWT Token Authentication (DEPRECATED - missing username header)
	// WARNING: Only sets X-SLURM-USER-TOKEN, will fail with most slurmrestd deployments
	// This is shown for backwards compatibility only - use WithUserToken instead
	jwtClient, err := slurm.NewClient(ctx,
		slurm.WithBaseURL(baseURL),
		slurm.WithAuth(auth.NewTokenAuth("your-jwt-token")),
	)
	if err == nil {
		defer jwtClient.Close()
		fmt.Println("✓ JWT Token authentication configured (deprecated - use WithUserToken)")
	}

	// Basic Authentication
	basicAuthClient, err := slurm.NewClient(ctx,
		slurm.WithBaseURL(baseURL),
		slurm.WithAuth(auth.NewBasicAuth("username", "password")),
	)
	if err == nil {
		defer basicAuthClient.Close()
		fmt.Println("✓ Basic authentication configured")
	}

	// No Authentication (for public endpoints)
	noAuthClient, err := slurm.NewClient(ctx,
		slurm.WithBaseURL(baseURL),
		slurm.WithNoAuth(),
	)
	if err == nil {
		defer noAuthClient.Close()
		fmt.Println("✓ No authentication configured (public endpoints only)")
	}
}
