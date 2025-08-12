# Public Release Preparation TODO

This document outlines the tasks needed to prepare the slurm-client project for public release.

## Required Actions (High Priority)

### ✅ 1. Complete Legal Documentation
- [x] ~~The NOTICE file has a TODO for third-party licenses that needs to be completed~~
- [x] ~~The THIRD_PARTY_LICENSES.md file needs to list all dependencies and their licenses~~
- [x] ~~Update attribution for SLURM OpenAPI specifications (Apache 2.0, not GPL)~~

**Status**: COMPLETED - Both files have been updated with proper SLURM attribution

### 2. Update GitHub Actions Configuration
- [ ] Change `runs-on: self-hosted` to `runs-on: ubuntu-latest` in all workflow files:
  - [ ] `.github/workflows/ci.yml`
  - [ ] `.github/workflows/security.yml`
  - [ ] `.github/workflows/release.yml`
  - [ ] `.github/workflows/docs.yml`
- [ ] This ensures public contributors can run CI/CD without needing your self-hosted runners

### ✅ 3. Clean Up Example Output Files
- [x] ~~Remove generated files:~~
  - [x] ~~`examples/efficiency-monitoring/efficiency_report_2025-07-19_13-04-05.json`~~
  - [x] ~~`examples/job-analytics/job_1001_analytics.json`~~
  - [x] ~~`examples/job-analytics/job_1002_analytics.json`~~
- [x] ~~Add `*.json` to `.gitignore` in the examples directories to prevent future commits~~

**Status**: COMPLETED - All example output files removed and .gitignore updated

### ✅ 4. Fix README Badge
- [x] ~~The codecov badge has `?token=YOUR_TOKEN` which should be updated or removed~~
- [x] ~~Either get a proper codecov token or remove the badge until codecov is properly configured~~

**Status**: COMPLETED - Codecov badge removed from README

### ✅ 5. Update Generated Code Headers
- [x] ~~Update code generation templates to use correct SLURM license attribution:~~
  ```go
  // Original specifications: Copyright (C) SchedMD LLC, Apache-2.0
  ```
  ~~(Changed from GPL-2.0-or-later to Apache-2.0 based on actual SLURM OpenAPI license)~~

**Status**: COMPLETED - All generated files have correct Apache-2.0 attribution

## Optional Improvements (Lower Priority)

### ✅ 1. Documentation
- [x] ~~Some example directories mentioned in READMes don't exist (job-allocation/, wckey-management/)~~
- [x] ~~Either create these examples or keep them commented out as you've already done~~

**Status**: COMPLETED - Non-existent directories have been commented out in README files

### ✅ 2. Personal Information Review
- [x] ~~Your full name "Jon Thor Kristinsson" appears in copyright headers~~
  - ~~This is standard practice but verify you're comfortable with it~~
  - ~~GitHub username `jontk` in paths is normal and expected~~
- [x] ~~Consider if you want to maintain this attribution or use a different format~~

**Status**: CONFIRMED - Attribution format is acceptable as-is

### 3. Generated Files Documentation
- [ ] Consider adding a note in README about which files are generated and shouldn't be manually edited
- [ ] Add clear markers in generated files (many already have "DO NOT EDIT" headers)

### 4. Repository Settings (GitHub UI)
- [ ] Enable branch protection for main branch
- [ ] Enable security alerts
- [ ] Configure secret scanning
- [ ] Set up code scanning with CodeQL
- [ ] Enable GitHub security advisories

### 5. Documentation Enhancements
- [ ] Consider adding CHANGELOG.md for future releases
- [ ] Verify all documentation links work correctly
- [ ] Consider adding more detailed API examples

## License Compliance Notes

### ✅ SLURM OpenAPI Specifications
**Status**: RESOLVED - The OpenAPI specifications are licensed under Apache 2.0 by SchedMD, not GPL as initially assumed. This significantly simplifies the licensing situation:

- ✅ Clean Apache 2.0 licensing throughout the project
- ✅ No GPL compatibility concerns
- ✅ Commercial-friendly licensing
- ✅ Proper attribution to SchedMD LLC included in NOTICE and THIRD_PARTY_LICENSES.md

## Release Checklist

Before making the repository public:

### Critical Items
- [ ] Complete all "Required Actions" above
- [ ] Verify no sensitive information in commit history
- [ ] Test that CI/CD works with `ubuntu-latest` runners
- [ ] Verify all documentation is accurate and helpful

### Pre-Release Testing
- [ ] Run full test suite
- [ ] Test examples work correctly
- [ ] Verify generated code compiles
- [ ] Check that all documentation renders properly

### Release Process
- [ ] Tag initial release (suggest v1.0.0 since the API appears stable)
- [ ] Create GitHub release with proper release notes
- [ ] Update any badges that reference the repository
- [ ] Announce to relevant communities (r/golang, SLURM mailing lists, etc.)

## Priority Assessment

**HIGH PRIORITY** (must complete before public release):
- GitHub Actions configuration
- Example file cleanup
- README badge fix

**MEDIUM PRIORITY** (should complete for professional appearance):
- Generated code header updates
- Documentation enhancements

**LOW PRIORITY** (nice to have):
- Repository security settings
- Advanced documentation features

## Estimated Time to Complete

- **Required Actions**: 2-3 hours
- **Optional Improvements**: 1-2 days
- **Total to Minimum Viable Public Release**: 2-3 hours

## Notes

The project is already in excellent shape for public release. The code quality, documentation, and architecture are all professional-grade. The remaining tasks are primarily administrative and cleanup items rather than fundamental issues.

---

**Last Updated**: 2025-08-12
**Next Review**: After completing required actions