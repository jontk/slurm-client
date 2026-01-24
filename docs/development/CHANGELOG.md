# Changelog

All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html)
and uses [Conventional Commits](https://www.conventionalcommits.org/) for commit messages.

## [Unreleased]

### ✨ Features

- enforce conventional commits and deploy docs to GitHub Pages (#36)
- add SBOM generation and code signing for releases
- add comprehensive test coverage for streaming and adapters
- add codecov configuration and code of conduct
- enable test quality linters (testifylint, thelper)
- enable gocritic linter
- enable performance linters (goconst, perfsprint) (#28)
- add pre-commit hooks for code quality enforcement (#27)
- enable performance linters (prealloc, lower gocyclo/dupl thresholds) (#26)
- enable error handling linters (nilerr, errorlint) (#25)
- enable low-risk linters (unconvert, wastedassign, usestdlibvars) (#23)

### 🐛 Bug Fixes

- add missing technical terms to spellcheck dictionary
- enable GitHub Pages deployment automatically (#37)
- enable scorecard publish_results with job-level permissions
- increase combined analytics overhead threshold to 150%
- disable scorecard publish_results to avoid permission conflict
- add id-token permission for OpenSSF Scorecard OIDC signing (#22)
- checkout main branch in update-documentation job

### 📚 Documentation

- add performance testing guide for local development

### 🔧 Maintenance

- Chore: update GitHub Actions and Go dependencies

## [v0.1.0] - 2026-01-20

### ⚠️ BREAKING CHANGES

- Expands public interface definitions. Users implementing
- Users relying on examples/adapter-pattern/ for reference

### ✨ Features

- add automated release infrastructure with GoReleaser (#20)
- add SLURM REST API v0.0.44 support
- expose node resource utilization fields in public interface
- implement job requeue functionality across all SLURM API versions
- add node control and job management methods to public API
- add Drain and Resume methods to NodeManager interface
- implement Hold/Release/Notify methods for SLURM API v0.0.42 and v0.0.43
- expose adapter-level functionality through public interfaces
- Implement full Cluster Management for v0.0.43 adapter
- Implement Cluster Management for v0.0.42 adapter
- complete v0.0.42 adapter implementation with full feature parity
- Implement Association Helper Methods for v0.0.42
- implement core v0.0.42 adapter functionality with job allocation and database ping
- Add examples and builder for new v0.0.43 features
- complete v0.0.43 SLURM REST API adapter implementation
- implement v0.0.43 SLURM REST API with WCKey management and job allocation
- add TRES resource handling for v0.0.43 adapters
- Implement missing v0.0.43 adapter functionality
- Complete open source transformation and fix all examples
- Add GitHub Actions CI/CD pipelines and quality tooling
- Transform project into best-in-class open source Go module
- Complete standalone adapter implementation with full API parity
- Add comprehensive QoS edge case handling and validation
- Enhanced SLURM error response parsing with comprehensive error codes
- Complete adapter implementation for all managers
- Enable and complete SLURM adapter implementation
- Fix adapter compilation issues
- Implement v0.0.41 API support with job and node operations
- Complete SLURM client adapter refactoring with comprehensive testing
- Add multi-version adapter support (v0.0.40, v0.0.41, v0.0.42)
- Implement comprehensive adapter pattern for SLURM API versioning
- Implement common error handling utilities and refactor QoS manager
- Add safeguards and documentation for generated code separation
- Add core packages for enhanced features implementation
- Add comprehensive performance and reliability enhancements
- Implement lazy initialization for all SLURM client managers
- Implement comprehensive AssociationManager for SLURM user-account-cluster relationships
- Implement ClusterManager and standalone operations for SLURM client
- Complete implementation of all 21 TODO items for SLURM user and account management
- Implement complete QoS and Reservation manager functionality
- Implement ReservationManager, QoSManager, AccountManager, and UserManager across all API versions
- Complete multi-version code generation with compilation fixes
- Complete 95%+ test coverage for all analytics functionality
- Complete efficiency-monitoring example with optimization recommendations
- Complete job-analytics example with comprehensive resource utilization analysis
- Task 5.2 - Add comprehensive performance benchmarks for analytics with <5% overhead validation
- Task 5.1 - Create comprehensive integration tests for job analytics using mock SLURM server
- Complete Task 4.8 - comprehensive unit tests for historical performance tracking
- Implement historical performance tracking framework (Task 4.1-4.6)
- Add efficiency calculations and optimization features (Task 3.0)
- Complete Task 2.7 - Add SLURM job step API integration with analytics
- Implement comprehensive ListJobStepsWithMetrics method across all API versions
- Implement GetJobResourceTrends for performance tracking (Task 2.3)
- Implement WatchJobMetrics for streaming performance updates (Task 2.2)
- Add GetJobLiveMetrics for real-time job performance monitoring
- Implement minimal job analytics for v0.0.40 with basic accounting
- Implement basic job analytics methods for v0.0.41
- Implement job analytics methods for v0.0.42 with enhanced metrics
- Implement GetJobUtilization, GetJobEfficiency, and GetJobPerformance in v0.0.43
- Extend JobManager interface with analytics methods
- Add job analytics data structures to interfaces.go
- Add comprehensive fair-share-analysis example (subtask 5.4)
- Add comprehensive user-account-management example (subtask 5.3)
- Complete task 4.0 - Add Fair-Share Information and Priority Calculation
- Complete task 3.0 - Implement User-Account Association Management
- Complete task 2.0 - Create UserManager Interface and Implementation
- Complete task 1.0 - Extend AccountManager Interface and Implementation
- Add real-time streaming support with WebSocket and SSE
- Add account management support for v0.0.43
- Add QoS (Quality of Service) management support for v0.0.43
- Add reservation management support for v0.0.43
- Add advanced examples for version differences, performance, and error recovery
- Add comprehensive examples for array jobs, dependencies, and resource allocation
- Implement real Watch() polling with state monitoring
- Implement v0.0.41 with basic job listing and info management
- Complete v0.0.40 implementation with all 17 interface methods
- Implement v0.0.43 manager methods and enhance testing
- Add real SLURM server integration testing framework
- Add comprehensive integration testing framework and performance optimization
- Complete v0.0.42 implementation with all 17 interface methods
- Complete NodeManager.Get() with structured error handling
- Update JobManager methods to use structured error handling
- Add comprehensive structured error handling system
- Complete NodeManager, PartitionManager, and InfoManager implementations
- Complete JobManager implementation with Get, Submit, Cancel methods
- Implement JobManager.List() for v0.0.42 with real OpenAPI integration
- Add comprehensive multi-version usage examples and documentation
- Complete v0.0.41 API bridge implementation with production-ready type conversions
- Add comprehensive test suite for production hardening
- Complete all four API versions (v0.0.40-v0.0.43) with full build support
- Complete core multi-version architecture with bridge pattern
- Complete v0.0.42 implementation with factory integration
- Complete multi-version client implementation with managers
- Implement multi-version Slurm REST API client architecture
- add Go client library for SLURM REST API

### 🐛 Bug Fixes

- sync slurm-cli with updated interfaces
- build slurm-cli from its module directory
- use GoReleaser v2 for version 2 config file support
- generate mock builders before running tests in release workflow
- make security workflow SARIF uploads non-failing
- resolve integration test compilation errors and improve code quality
- Resolve build failures and enhance test coverage across all adapter versions
- Resolve WCKey Management API type issues for v0.0.42
- Enable v0.0.43 features in factory client
- Resolve test compilation errors and remove outdated test files
- resolve compilation errors and add missing interface methods
- Resolve build failures and test compilation errors
- Remove trailing whitespace from adapter test files
- Update tests to match current API interfaces
- Resolve missing dependencies and build errors
- Remove array-jobs and batch-operations files
- Build issues and adapter implementation fixes
- Complete critical adapter fixes for v0.0.43
- Critical adapter implementation fixes for v0.0.43
- Complete high-priority adapter implementation fixes
- Complete resolution of all compilation errors across adapters and tests
- Resolve critical SLURM adapter compilation errors
- Comprehensive SLURM adapter compilation error resolution
- Fix job submission across all SLURM API versions
- Add AssociationManager stub for older API versions
- Update code generator to produce wrapper.go files with actual manager implementations
- Add authentication middleware to HTTP client
- Resolve critical import cycle issue in multi-version client architecture

### ⚡ Performance

- performance.go:
- performance summary across all supported API versions with progressive enhancement.

### 📚 Documentation

- prepare project for public release
- Enhance documentation and add comprehensive examples
- Add comprehensive adapter pattern documentation and examples
- Update README with reservation management documentation
- Update CLAUDE.md to reflect completed v0.0.42/v0.0.43 implementations
- Comprehensive README update reflecting production-ready status
- Update CLAUDE.md to reflect structured error handling completion
- Organize project documentation and create comprehensive status review

### ♻️ Code Refactoring

- Optimize adapter implementations and remove redundancy
- Consolidate adapter helpers and optimize converters
- Rename ExponentialBackoff to HTTPExponentialBackoff and remove duplicate backoff.go
- Fix deprecated API usage and improve code quality
- Convert OpenAPI specs from YAML to JSON format

### ✅ Tests

- add comprehensive test coverage for v0.0.42 and v0.0.43 adapters
- add comprehensive unit tests for v0.0.43 adapter functionality
- enhance test coverage for job analytics and fix nil pointer issues
- Significantly improve test coverage across core packages
- test success rate from 83.3% to 87.9%.
- Comprehensive SLURM adapter testing and fixes
- Add comprehensive adapter tests for reservation and association managers
- Comprehensive adapter testing against real SLURM cluster
- Add comprehensive adapter test coverage
- Build unit tests for performance monitoring (Task 2.8)
- Add comprehensive unit tests for job analytics across all API versions
- Add comprehensive test coverage for user-account management (subtask 5.8)
- Add comprehensive unit tests for v0.0.40 and v0.0.41
- Add comprehensive unit tests for v0.0.43 implementation

### 🔧 Maintenance

- CI: migrate from self-hosted to GitHub-hosted runners
- CI: switch GitHub workflows to self-hosted runners
- Chore: Update .gitignore for Claude Flow metrics and build artifacts
- Chore: remove obsolete adapter pattern examples
- Chore: Remove problematic test files moved during implementation
- Chore: Complete v0.0.40/41 implementations and v0.0.43 unit tests

