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
	"os"
	"time"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/pkg/auth"
	"golang.org/x/time/rate"
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
	// Create a basic client with JWT token authentication
	client, err := slurm.NewClient(ctx,
		slurm.WithBaseURL("https://localhost:6820"),
		slurm.WithAuth(auth.NewTokenAuth("your-jwt-token")),
	)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
	defer client.Close()

	// Test the connection
	info, err := client.Info().Ping(ctx)
	if err != nil {
		log.Printf("Failed to ping server: %v", err)
		return
	}

	fmt.Printf("Connected to SLURM REST API\n")
	fmt.Printf("API Version: %s\n", info.Meta.Plugin.Version)
	fmt.Printf("SLURM Version: %s\n", info.Meta.Slurm.Version)
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
		slurm.WithAuth(auth.NewTokenAuth("your-jwt-token")),
		
		// Set custom timeout
		slurm.WithTimeout(30*time.Second),
		
		// Configure retry behavior
		slurm.WithMaxRetries(5),
		slurm.WithRetryWaitMin(1*time.Second),
		slurm.WithRetryWaitMax(30*time.Second),
		
		// Set rate limiting
		slurm.WithRateLimiter(rate.NewLimiter(10, 1)), // 10 requests/second
		
		// Custom TLS configuration
		slurm.WithTLSConfig(&tls.Config{
			MinVersion: tls.VersionTLS12,
			// Only for development/testing
			InsecureSkipVerify: true,
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

	// JWT Token Authentication (recommended)
	jwtClient, err := slurm.NewClient(ctx,
		slurm.WithBaseURL(baseURL),
		slurm.WithAuth(auth.NewTokenAuth("your-jwt-token")),
	)
	if err == nil {
		defer jwtClient.Close()
		fmt.Println("✓ JWT Token authentication configured")
	}

	// API Key Authentication
	apiKeyClient, err := slurm.NewClient(ctx,
		slurm.WithBaseURL(baseURL),
		slurm.WithAuth(auth.NewAPIKeyAuth("X-SLURM-Token", "your-api-key")),
	)
	if err == nil {
		defer apiKeyClient.Close()
		fmt.Println("✓ API Key authentication configured")
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

	// Certificate Authentication
	certAuthClient, err := slurm.NewClient(ctx,
		slurm.WithBaseURL(baseURL),
		slurm.WithAuth(auth.NewCertAuth("/path/to/cert.pem", "/path/to/key.pem")),
	)
	if err == nil {
		defer certAuthClient.Close()
		fmt.Println("✓ Certificate authentication configured")
	}
}