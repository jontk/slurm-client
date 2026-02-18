// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

/*
Package slurm provides a comprehensive Go client library for the SLURM REST API.

The library supports multiple API versions (v0.0.40 through v0.0.44) and provides
both high-level abstractions and direct API access for SLURM workload management.

# Overview

This client library addresses the fragmentation in the Go SLURM ecosystem by providing:
  - Unified interface across multiple SLURM REST API versions
  - Type-safe operations with full OpenAPI-generated structs
  - Enterprise-grade error handling and retry mechanisms
  - Comprehensive authentication support
  - Context-aware operations for cancellation and timeouts

# Installation

Install the library using Go modules:

	go get github.com/jontk/slurm-client

# Basic Usage (Recommended: Adapter Pattern)

Create a client using the adapter pattern for version-agnostic, production-ready code:

	import (
	    "context"
	    "log"

	    "github.com/jontk/slurm-client"
	    "github.com/jontk/slurm-client/pkg/auth"
	)

	func main() {
	    ctx := context.Background()

	    // 🎯 Adapter Pattern: Auto-detects version, handles conversions
	    client, err := slurm.NewClient(ctx,
	        slurm.WithBaseURL("https://cluster.example.com:6820"),
	        slurm.WithAuth(auth.NewTokenAuth("your-jwt-token")),
	        // Version auto-detected - works across all SLURM versions!
	    )
	    if err != nil {
	        log.Fatal(err)
	    }
	    defer client.Close()

	    // Version-agnostic API calls with automatic type conversion
	    jobs, err := client.Jobs().List(ctx, nil)
	    if err != nil {
	        log.Fatal(err)
	    }

	    for _, job := range jobs {
	        log.Printf("Job %d: %s (state: %s)\n",
	            job.JobID, job.Name, job.State)
	    }
	}

# Architecture

The library uses the Adapter Pattern to provide version-agnostic interfaces:

  - Version-agnostic interfaces that work across all SLURM versions
  - Automatic type conversion and validation
  - Simplified error handling with structured errors
  - Production-ready with caching and optimizations
  - Best for: Applications, automation, most use cases

# Version Support

The library supports the following SLURM REST API versions:
  - v0.0.40 (SLURM 23.02.x)
  - v0.0.41 (SLURM 23.11.x)
  - v0.0.42 (SLURM 24.05.x)
  - v0.0.43 (SLURM 24.11.x)
  - v0.0.44 (SLURM 25.02.x)

Version detection is automatic, but can be overridden using NewClientWithVersion:

	client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
	    slurm.WithBaseURL("https://cluster:6820"),
	    slurm.WithAuth(auth.NewTokenAuth("token")),
	)

# Authentication

SLURM REST API uses JWT token authentication:

User Token (RECOMMENDED - with username header):

	// Sets both X-SLURM-USER-NAME and X-SLURM-USER-TOKEN headers
	// This is required by most SLURM deployments
	client, err := slurm.NewClient(ctx,
	    slurm.WithBaseURL("https://cluster:6820"),
	    slurm.WithUserToken("username", "your-jwt-token"),
	)

JWT Token (DEPRECATED - missing username header):

	// WARNING: Only sets X-SLURM-USER-TOKEN, will fail with most slurmrestd deployments
	authProvider := auth.NewTokenAuth("your-jwt-token")  // Deprecated: Use WithUserToken instead

Environment Variable:

	// Token is read from SLURM_JWT environment variable
	export SLURM_JWT="your-jwt-token"

# Error Handling

The library provides structured error handling with typed errors:

	jobs, err := client.Jobs().List(ctx, nil)
	if err != nil {
	    var apiErr *types.APIError
	    if errors.As(err, &apiErr) {
	        log.Printf("API Error: %s (code: %s)\n",
	            apiErr.Message, apiErr.ErrorCode)

	        // Check specific error types
	        if apiErr.IsAuthError() {
	            // Handle authentication errors
	        } else if apiErr.IsRateLimitError() {
	            // Handle rate limiting
	        }
	    }
	}

# Advanced Features

Connection Options:

	client, err := slurm.NewClient(ctx,
	    slurm.WithBaseURL("https://cluster:6820"),
	    slurm.WithTimeout(30 * time.Second),
	    slurm.WithMaxRetries(3),
	    slurm.WithRetryPolicy(retry.NewExponentialBackoff(
	        100*time.Millisecond, // initial delay
	        5*time.Second,        // max delay
	        2.0,                  // multiplier
	    )),
	)

Context Usage:

	// With timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	jobs, err := client.Jobs().List(ctx, nil)

	// With cancellation
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
	    // Cancel after some condition
	    cancel()
	}()
	jobs, err := client.Jobs().List(ctx, nil)

Filtering and Pagination:

	// Filter jobs by state
	filters := &types.JobFilters{
	    States: []string{"RUNNING", "PENDING"},
	    Users:  []string{"alice", "bob"},
	}
	jobs, err := client.Jobs().List(ctx, filters)

	// Pagination support
	opts := &types.ListOptions{
	    Limit:  100,
	    Offset: 0,
	}
	jobs, err := client.Jobs().ListWithOptions(ctx, filters, opts)

# Manager Interfaces

The client provides managers for different resource types:

Jobs Manager:

	jobMgr := client.Jobs()
	jobs, _ := jobMgr.List(ctx, nil)
	job, _ := jobMgr.Get(ctx, 12345)
	jobID, _ := jobMgr.Submit(ctx, jobSpec)
	_ = jobMgr.Cancel(ctx, 12345)

Nodes Manager:

	nodeMgr := client.Nodes()
	nodes, _ := nodeMgr.List(ctx)
	node, _ := nodeMgr.Get(ctx, "node001")
	_ = nodeMgr.Update(ctx, "node001", updates)

Partitions Manager:

	partMgr := client.Partitions()
	partitions, _ := partMgr.List(ctx)
	partition, _ := partMgr.Get(ctx, "compute")

Info Manager:

	infoMgr := client.Info()
	ping, _ := infoMgr.Ping(ctx)
	diag, _ := infoMgr.Diagnostics(ctx)
	stats, _ := infoMgr.Statistics(ctx)

# Best Practices

1. Always use context for cancellation support
2. Handle errors appropriately, checking for specific error types
3. Use connection pooling by reusing client instances
4. Set appropriate timeouts for long-running operations
5. Use filters to minimize data transfer
6. Close the client when done to release resources

# Environment Variables

The client respects the following environment variables:

  - SLURM_REST_URL: Default base URL for the SLURM REST API
  - SLURM_JWT: JWT token for authentication
  - SLURM_API_VERSION: Force specific API version
  - SLURM_TIMEOUT: Request timeout duration (e.g., "30s")
  - SLURM_INSECURE_SKIP_VERIFY: Skip TLS verification (development only)
  - SLURM_MAX_RETRIES: Maximum number of retries for failed requests
  - SLURM_DEBUG: Enable debug logging

# Thread Safety

All client operations are thread-safe and can be called concurrently
from multiple goroutines. The client maintains internal connection
pooling for optimal performance.

# Contributing

Contributions are welcome! Please see CONTRIBUTING.md for guidelines.

# License

This library is licensed under the Apache License 2.0. See LICENSE for details.
*/
package slurm
