# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] - 2026-02-08

### Added
- **SDK Architectural Refactor** (#94): Complete restructure with unified adapter pattern
  - Squashed 300+ commits of development work into production-ready SDK
  - Clean type system in `api/` package with generated and hand-written types
  - Version-agnostic client interface via `NewClient`/`NewClientWithVersion`
  - JWT-based authentication (`NewTokenAuth`, `NewJWTAuth`)
  - Connection pooling, retry middleware, streaming support
  - Comprehensive analytics (job efficiency, performance metrics, resource utilization)
  - Builder pattern for jobs, partitions, accounts, users, QoS
  - Multi-version support with automatic version detection for v0.0.40-v0.0.44
- **Documentation Enhancements** (#99): Comprehensive review and expansion
  - New adapter pattern deep dive guide (526 lines)
  - New testing and mocking guide (721 lines)
  - Expanded troubleshooting guide (+305 lines)
  - Fixed all import paths to use public API only
  - Standardized client creation with functional options pattern
  - Updated version support information and EOL dates
  - 13 files changed, 1742 insertions, 1157 deletions

### Changed
- **Package Restructure** (#94): Complete SDK reorganization
  - `api/` - Public types and interfaces for consumers
  - `internal/adapters/` - Version-specific adapter implementations
  - `internal/factory/` - Client factory and adapter client
  - `pkg/` - Middleware, streaming, analytics, connection pooling
  - `tools/codegen/` - Code generation for types and converters
- **Code Quality** (#95): Post-refactor CI improvements
  - Added lint exclusions for goverter-generated files
  - Added lint exclusions for examples and converter helpers
  - Reduced lint errors from 1523 to 333 (78% reduction)
  - Fixed build constraints for examples with `//go:build ignore`

### Fixed
- **Circuit Breaker Concurrency** (#100): Added `sync.RWMutex` for thread-safe concurrent access
  - Prevents data races in middleware circuit breaker during concurrent requests
  - Added mutex locking in `Allow()`, `RecordFailure()`, and `RecordSuccess()` methods
  - Critical fix for production environments with high concurrency
- **Pagination Metadata** (#100): Fixed `Total` field returning page size instead of total count
  - Pagination now correctly returns total count before pagination in `adapter_client.go`
  - Clients can now accurately display total available records
- **Job List Filters** (#100): Fixed filters being dropped before pagination
  - Root cause: Generator template was missing filter application step
  - Updated adapter generator to apply `FilterJobList` before pagination
  - Regenerated all adapter versions (v0.0.42, v0.0.43, v0.0.44)
  - Future generator runs will preserve the fix (no more manual patches needed)
- **Watch Options Semantics** (#100): Removed incorrect States to EventTypes mapping
  - States (job states like RUNNING, PENDING) and EventTypes (event names like start, end, fail) are fundamentally different
  - Added clarifying comment explaining the semantic difference
  - State filtering must be implemented at higher level through event filtering
- **Account Hierarchy Error Suppression** (#100): Added proper error handling to `buildHierarchyNode`
  - Changed function signature to return `(*types.AccountHierarchy, error)`
  - Errors from `adapter.Get()` are now surfaced instead of being silently ignored
  - Prevents returning incomplete account hierarchy data

### Improved
- **Code Generator** (#100): Fixed test generator producing incorrect types
  - Node State field: Now generates `[]types.NodeState{...}` instead of `ptrNodeState(...)`
  - Reservation Name field: Now generates `&testName` instead of `testName`
  - Future test regenerations will produce correct types
  - Eliminates CI failures from type mismatches in generated tests

### Technical Debt
- **Generator Maintenance**: All manual fixes converted to permanent generator improvements
  - No more `.gen.go` files requiring manual patches after regeneration
  - Closed tracking issue `slurm-client-433`

## [0.2.4] - 2026-01-30

### Fixed
- **Job Filtering**: Fixed empty partition filter blocking all jobs (#86)
  - Empty partition/user filter strings were creating arrays like `[""]` instead of empty arrays
  - Jobs now show correctly in listings when filter parameters are empty
  - Fixed in job_adapter.go and adapter_client.go
- **Node State Handling**: Preserve all node state flags when converting API response (#86)
  - Multi-valued node state arrays were being truncated to first element only
  - DRAIN flags were lost, causing drained nodes to appear as idle
  - Now concatenates all state elements with "+" separator (e.g. "IDLE+DRAIN")
  - Fixed across all SLURM API versions (v0.0.40 through v0.0.44)
- **Job Submission**: Removed overly strict account validation for job submission (#86)
  - Job creation was blocking on missing account field unnecessarily
  - Now lets SLURM API handle account requirements
- **Job Field Mapping**: Fixed incorrect Command field mapping from API response (#86)
  - Command field from API is typically null/empty for most jobs
  - Removed Command field mapping, kept WorkingDirectory field with explanatory comment
- **Code Quality**: Replaced string concatenation in loops with strings.Join (#86)
  - Fixed 5 linter warnings (4x perfsprint concatenation-loop, 1x unconvert)
  - Improved performance for state concatenation across all node adapters

### Added
- **Node Updates**: Include State and Reason fields in node update requests (#86)
  - Enables proper node state management and drain operations
  - Node update operations now support changing state and providing reasons
- **User Management**: Extract AdminLevel from user API response (#86)
  - User objects now include complete admin level information
  - Enables proper user permission management
- **Job Information**: Expose UserName field from SLURM API (#86)
  - Job listings now show actual username instead of just user ID
  - Improves job visibility and management

### Changed
- **Debug Output**: Removed debug printf statements from job_adapter (#86)
  - Cleaner production code without debug noise

## [0.2.3] - 2026-01-28

### Fixed
- **Release Workflow**: Added missing pull-requests write permission (#83)
  - Documentation update job needs `pull-requests: write` to create PRs
  - Fixes "Resource not accessible by integration" error when creating documentation PRs
  - Completes the fix started in #81 for documentation update automation

## [0.2.2] - 2026-01-28

### Fixed
- **Release Workflow**: Fixed documentation update job failing with duplicate Authorization headers (#81)
  - Added `persist-credentials: false` to checkout action in update-documentation job
  - Prevents conflict between checkout's git auth config and create-pull-request's auth
  - Resolves 400 error that occurred even after removing explicit token parameter

## [0.2.1] - 2026-01-27

### Fixed
- **Go Module Availability**: Re-released as v0.2.1 to ensure Go module proxy has correct version
  - The v0.2.0 tag was force-pushed multiple times during release troubleshooting
  - Go module proxy (proxy.golang.org) cached the initial version and didn't update
  - v0.2.1 is identical to v0.2.0 content but ensures users get all features via `go get`
  - All v0.2.0 features are included, particularly WithUseAdapters() API

## [0.2.0] - 2026-01-27

### Added
- **Adapter Pattern Completion**: InfoAdapter implementation across all API versions (v0.0.40-v0.0.44)
  - Complete InfoManager with real API calls (Get, Ping, PingDatabase, Stats, Version)
  - Proper error handling and type conversions for all versions
  - Eliminated hybrid wrapper/adapter approach in v0.0.44
- **Performance Testing Infrastructure**: Nightly validation workflow (#74)
  - Removed benchmarks from PR CI to eliminate flaky failures
  - Added nightly performance validation workflow (4 AM UTC)
  - Smoke tests for catastrophic regression detection (30s threshold)
  - Platform-agnostic approach without macOS-specific skips
- **Standalone Operations**: Full support for standalone SLURM operations (v0.0.40/v0.0.42)
  - Licenses management
  - TRES (Trackable RESources)
  - Shares
  - Diagnostics
- **New Adapters**: Additional resource management for v0.0.40
  - WCKeyAdapter - WCKey management
  - ClusterAdapter - Cluster management
- **Linting Infrastructure**: Comprehensive code quality governance
  - goheader configuration for SPDX Apache-2.0 license enforcement
  - depguard rules for deprecated packages
  - Documented staticcheck exclusions with violation counts
  - Refined gosec exclusions with path-based rules
  - Enabled revive linter with 3-phase rollout (40+ rules)
  - Enabled errcheck type assertion checks
  - Enabled intrange linter for Go 1.22+ improvements
  - Enabled unparam linter with strategic exclusions
  - Enabled noctx and forcetypeassert linters
  - Enabled contextcheck, usetesting, dupword, and exhaustive
- **Documentation**: MkDocs Material site with autogenerated CLI and changelog
- **Release Infrastructure**:
  - GoReleaser configuration for automated releases
  - SECURITY.md with security policy and vulnerability reporting
  - RELEASE.md with detailed release process documentation
  - Automated release workflow triggered by git tags

### Changed
- **Architecture**: Removed hybrid wrapper/adapter approach from v0.0.44
  - AdapterClient no longer creates dual clients internally
  - Factory properly delegates Info() to adapter.GetInfoManager()
  - Added type converters between types.* and interfaces.*
- **Code Quality**: Reduced cyclomatic complexity thresholds from 50 to 25 to 20
  - Extracted helper methods across multiple adapters
  - Improved code readability and maintainability
- **Performance Testing Strategy**: Restructured to prevent flaky CI failures (#74)
  - Removed benchmarks from PR workflow entirely
  - Converted timing assertions to data-collection only
  - Benchmarks now run exclusively via nightly validation
  - Reduced CI time and resource usage on every PR

### Fixed
- **Critical Bug**: Fixed nil pointer dereferences in all adapter List methods (#71)
  - JobAdapter, NodeAdapter, PartitionAdapter, AccountAdapter, UserAdapter
  - QOSAdapter, ReservationAdapter, LicenseAdapter, TRESAdapter
  - ShareAdapter, DiagnosticsAdapter, WCKeyAdapter, ClusterAdapter
  - All adapters now properly check for nil slices before conversion
- **Code Quality**: Fixed 7 cyclomatic complexity violations by extracting helper methods
  - v0.0.40-42: InfoAdapter.Stats() reduced from complexity 25 to 6
  - v0.0.43-44: InfoAdapter.Get() reduced from 21 to 7, Stats() reduced from 26 to 7
- **Type Conversions**: Fixed 4 unnecessary type conversion warnings (unconvert)
  - Removed redundant string() and int64() casts
- **Linter Violations**: Fixed 58 unparam violations across all API versions
  - Removed unused error returns from adapter converters
  - Removed unused parameters from API managers
- **Unused Code**: Fixed unused append operations in adapters (SA4010)
- **Gocritic**: Fixed nestingReduce violations and removed text exclusions
- **Module Tidiness**: Updated go.mod/go.sum for cmd/slurm-cli and examples
- **CI Failures**: Resolved staticcheck and gosec warnings
  - Fixed unused linter directives
  - Added safe integer overflow checks
- **Performance Tests**: Eliminated flaky timing-based test failures (#74, #70)
  - Removed platform-specific thresholds (macOS 150% overhead)
  - Tests now report metrics without failing on timing variability
  - Smoke tests catch only catastrophic regressions (>10x slower)

## [0.1.0] - 2026-01-20

### Added
- **CI/CD Infrastructure**: GitHub Actions workflows for comprehensive testing
  - Multi-version Go testing (1.21, 1.22, 1.23)
  - Cross-platform testing (Linux, macOS, Windows)
  - Security scanning (gosec, govulncheck, CodeQL)
  - Dependency review and license compliance
  - Documentation generation and validation
- **Test Infrastructure**: Comprehensive test suite with mock server
  - Version-specific mock handler generation (v0.0.40-v0.0.44)
  - Integration tests covering all API versions
  - Performance benchmarks with platform-specific thresholds
  - Mock builder smoke tests
- **Documentation**: Comprehensive migration and setup guides
  - MIGRATION.md with version upgrade paths
  - Enhanced README with setup instructions
  - API documentation for jobs, nodes, and partitions
- **Code Quality**: Linting and code analysis tools
  - golangci-lint with comprehensive rule set
  - gosec security scanner configuration
  - cspell for documentation spell-checking
  - License header compliance checks
- **API Support**: Full support for SLURM REST API v0.0.44 (SLURM 25.11.0)
  - Complete adapter implementation for all resource managers
  - Comprehensive test coverage with 100% passing tests
  - Type-safe conversions between v0.0.44 API types and common types
  - Error handling with status code inclusion in error messages
  - Validation for all CRUD operations

### Changed
- **Architecture**: Refactored interfaces to dedicated package
  - Moved all interface definitions to `interfaces` package
  - Improved code organization and modularity
  - Prevented import cycles
- **Test Organization**: Moved test helpers to `*_test.go` files
  - Resolved import cycle issues
  - Improved test package isolation
- **API Version Support**: Updated multi-version support range from v0.0.40-v0.0.43 to v0.0.40-v0.0.44
- **Error Messages**: Enhanced error adapter to include status codes (404, 409, 500) in error messages
- **Ignore Patterns**: Improved .gitignore patterns for temporary files and reports

### Fixed
- **Critical Bugs**: Resolved panics and nil pointer dereferences
  - v0.0.41 PartitionManager nil client checks
  - v0.0.41 InfoManager nil client checks
  - Divide-by-zero in job step task generation
  - Nil pointer in resource trends with empty time series
- **API Issues**: Fixed routing and response handling
  - Analytics endpoint trailing slash compatibility
  - Dual route registration for `/jobs` endpoint
  - Response structure unwrapping in analytics tests
  - Query parameter naming (reference_job_id)
- **Error Handling**: Enhanced error classification
  - Proper deadline exceeded error handling
  - Network timeout error classification
  - Context deadline support in timeout tests
- **Test Reliability**: Platform-specific test adjustments
  - macOS performance test thresholds (80%)
  - Windows analytics overhead thresholds (10%)
  - Reduced flakiness on fast systems
- **Validation**: Improved consistency across adapters
  - Enhanced base manager validation
  - Consistent error messages across versions
  - Proper field mapping and type conversions
  - Case sensitivity issues in validation error messages
  - Type conversion for PartitionState in partition adapter
  - Empty update validation to require at least one field
  - Reservation validation to require both StartTime and EndTime

### Security
- **Dependency Updates**: Critical security patches
  - Updated kin-openapi from v0.128.0 to v0.133.0
  - Addressed CVE-2025-30153 in OpenAPI validation library
- **Security Infrastructure**: Automated vulnerability scanning
  - CodeQL analysis for Go code
  - Trivy container and dependency scanning
  - Nancy dependency vulnerability checks
  - OpenSSF Scorecard integration

---

## Version Support Matrix

| SLURM Version | REST API Version | Support Status |
|---------------|------------------|----------------|
| 25.11.0       | v0.0.44          | ✅ Supported    |
| 25.05.0       | v0.0.43          | ✅ Supported    |
| 24.11.0       | v0.0.42          | ✅ Supported    |
| 24.05.0       | v0.0.41          | ✅ Supported    |
| 23.11.0       | v0.0.40          | ✅ Supported    |

---

*For detailed changes, see the [commit history](https://github.com/jontk/slurm-client/commits/main).*
