// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

/*
Package slurm provides a comprehensive Go client library for the SLURM REST API.

The library supports multiple API versions (v0.0.40 through v0.0.43) and provides
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

# Basic Usage

Create a client and perform basic operations:

	import (
	    "context"
	    "log"
	    
	    "github.com/jontk/slurm-client"
	    "github.com/jontk/slurm-client/pkg/auth"
	)
	
	func main() {
	    ctx := context.Background()
	    
	    // Create client with token authentication
	    client, err := slurm.NewClient(ctx,
	        slurm.WithBaseURL("https://cluster.example.com:6820"),
	        slurm.WithAuth(auth.NewTokenAuth("your-jwt-token")),
	    )
	    if err != nil {
	        log.Fatal(err)
	    }
	    defer client.Close()
	    
	    // List jobs
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

The library implements two complementary patterns:

1. Adapter Pattern (Recommended)
   - Version-agnostic interfaces
   - Automatic type conversion
   - Simplified error handling
   - Best for most use cases

2. Wrapper Pattern
   - Direct access to version-specific APIs
   - Full control over operations
   - Minimal overhead
   - Best for advanced use cases

# Version Support

The library supports the following SLURM REST API versions:
  - v0.0.40 (SLURM 23.02.x)
  - v0.0.41 (SLURM 23.11.x)
  - v0.0.42 (SLURM 24.05.x)
  - v0.0.43 (SLURM 24.11.x)

Version detection is automatic, but can be overridden:

	client, err := slurm.NewClient(ctx,
	    slurm.WithBaseURL("https://cluster:6820"),
	    slurm.WithVersion("v0.0.43"),
	    slurm.WithAuth(auth.NewTokenAuth("token")),
	)

# Authentication

Multiple authentication methods are supported:

JWT Token (Recommended):

	auth := auth.NewTokenAuth("your-jwt-token")

API Key:

	auth := auth.NewAPIKeyAuth("X-SLURM-Token", "your-api-key")

Basic Authentication:

	auth := auth.NewBasicAuth("username", "password")

Certificate Authentication:

	auth := auth.NewCertAuth("/path/to/cert.pem", "/path/to/key.pem")

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
	    slurm.WithRateLimiter(rate.NewLimiter(10, 1)),
	    slurm.WithTLSConfig(&tls.Config{
	        InsecureSkipVerify: false,
	    }),
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

  - SLURM_API_URL: Default base URL for the SLURM REST API
  - SLURM_API_TOKEN: Default JWT token for authentication
  - SLURM_API_VERSION: Force specific API version
  - SLURM_API_TIMEOUT: Default timeout in seconds
  - SLURM_API_INSECURE: Skip TLS verification (development only)

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