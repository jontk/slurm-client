# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- GoReleaser configuration for automated releases
- SECURITY.md with security policy and vulnerability reporting
- RELEASE.md with detailed release process documentation
- Automated release workflow triggered by git tags

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
