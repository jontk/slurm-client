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

# Basic Usage

Create a client with automatic version detection:

	import (
	    "context"
	    "log"

	    "github.com/jontk/slurm-client"
	)

	func main() {
	    ctx := context.Background()

	    client, err := slurm.NewClient(ctx,
	        slurm.WithBaseURL("https://cluster.example.com:6820"),
	        slurm.WithUserToken("username", "your-jwt-token"),
	    )
	    if err != nil {
	        log.Fatal(err)
	    }
	    defer client.Close()

	    jobs, err := client.Jobs().List(ctx, nil)
	    if err != nil {
	        log.Fatal(err)
	    }

	    for _, job := range jobs.Jobs {
	        log.Printf("Job %v: %v (state: %v)\n",
	            job.JobID, job.Name, job.JobState)
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

Version detection is automatic, but can be overridden using [NewClientWithVersion]:

	client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
	    slurm.WithBaseURL("https://cluster:6820"),
	    slurm.WithUserToken("username", "token"),
	)

# Authentication

SLURM REST API requires both a username and JWT token for authentication.
Use [WithUserToken] which sets both the X-SLURM-USER-NAME and X-SLURM-USER-TOKEN headers:

	client, err := slurm.NewClient(ctx,
	    slurm.WithBaseURL("https://cluster:6820"),
	    slurm.WithUserToken("username", "your-jwt-token"),
	)

# Error Handling

The library provides structured error handling via [SlurmError] and the
[github.com/jontk/slurm-client/pkg/errors] package:

	import "github.com/jontk/slurm-client/pkg/errors"

	jobs, err := client.Jobs().List(ctx, nil)
	if err != nil {
	    if errors.IsAuthenticationError(err) {
	        log.Fatal("Authentication failed - check your credentials")
	    }
	    if errors.IsNetworkError(err) {
	        log.Fatal("Network error - check server connectivity")
	    }
	    if errors.IsRetryableError(err) {
	        // Safe to retry this operation
	    }
	    log.Fatal(err)
	}

# Advanced Features

Connection Options:

	client, err := slurm.NewClient(ctx,
	    slurm.WithBaseURL("https://cluster:6820"),
	    slurm.WithUserToken("username", "token"),
	    slurm.WithTimeout(30 * time.Second),
	)

Context Usage:

	// With timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	jobs, err := client.Jobs().List(ctx, nil)

	// With cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	jobs, err = client.Jobs().List(ctx, nil)

Filtering:

	// Filter jobs by user
	jobs, err := client.Jobs().List(ctx, &slurm.ListJobsOptions{
	    UserID: "alice",
	})

# Manager Interfaces

The client provides managers for different resource types:

Jobs Manager:

	jobMgr := client.Jobs()
	jobs, _ := jobMgr.List(ctx, nil)
	job, _ := jobMgr.Get(ctx, "12345")
	resp, _ := jobMgr.Submit(ctx, jobSpec)
	_ = jobMgr.Cancel(ctx, "12345")

Nodes Manager:

	nodeMgr := client.Nodes()
	nodes, _ := nodeMgr.List(ctx, nil)
	node, _ := nodeMgr.Get(ctx, "node001")

Partitions Manager:

	partMgr := client.Partitions()
	partitions, _ := partMgr.List(ctx, nil)
	partition, _ := partMgr.Get(ctx, "compute")

Info Manager:

	infoMgr := client.Info()
	_ = infoMgr.Ping(ctx)
	info, _ := infoMgr.Get(ctx)
	stats, _ := infoMgr.Stats(ctx)

# Best Practices

1. Always use context for cancellation support
2. Handle errors appropriately, checking for specific error types
3. Use connection pooling by reusing client instances
4. Set appropriate timeouts for long-running operations
5. Use filters to minimize data transfer
6. Close the client when done to release resources

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
