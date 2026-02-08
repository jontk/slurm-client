# Release Process

This document describes the release process for slurm-client.

## Overview

We use [Semantic Versioning](https://semver.org/) and automated releases via [GoReleaser](https://goreleaser.com/). The release process is triggered by pushing a git tag.

## Versioning Strategy

### Semantic Versioning

We follow semantic versioning: `MAJOR.MINOR.PATCH`

- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality in a backwards-compatible manner
- **PATCH**: Backwards-compatible bug fixes

### Pre-v1.0.0 Releases (Current)

We're currently in the v0.x.x phase, which signals:
- API may change between minor versions
- Production-ready but API not yet stable
- We're gathering feedback before committing to v1.0.0 API stability

### v1.0.0 and Beyond

When we release v1.0.0, it will signal:
- Stable, production-ready API
- Breaking changes require major version bump
- Full backward compatibility within major version

## Release Types

### Regular Releases

Standard releases with new features, bug fixes, and improvements.

**Version Format**: `v0.MINOR.PATCH` or `v1.MINOR.PATCH`

### Patch Releases

Bug fixes and security updates without new features.

**Version Format**: `v0.MINOR.PATCH` (increment PATCH)

### Pre-releases

Test releases before a major/minor version.

**Version Format**: `v0.MINOR.0-rc.1`, `v0.MINOR.0-beta.1`, `v0.MINOR.0-alpha.1`

## Release Process

### Prerequisites

1. All tests passing on main branch
2. CHANGELOG.md updated with release notes
3. No open critical bugs
4. Security scan passed
5. Documentation up to date

### Step-by-Step Release

#### 1. Prepare the Release

```bash
# Ensure you're on main and up to date
git checkout main
git pull origin main

# Verify all tests pass
go test ./...

# Run linters
golangci-lint run

# Run security scan
gosec ./...
```

#### 2. Update CHANGELOG.md

Move items from `[Unreleased]` to a new version section:

```markdown
## [Unreleased]

## [0.2.0] - 2026-01-25

### Added
- New feature X
- New feature Y

### Fixed
- Bug fix A
```

Commit the changelog:

```bash
git add CHANGELOG.md
git commit -m "chore: prepare v0.2.0 release"
git push origin main
```

#### 3. Create and Push Tag

```bash
# Create annotated tag
git tag -a v0.2.0 -m "Release v0.2.0"

# Push tag (this triggers the release workflow)
git push origin v0.2.0
```

#### 4. Automated Release

Once the tag is pushed:

1. GitHub Actions workflow is triggered
2. GoReleaser builds binaries for all platforms
3. Creates GitHub Release with:
   - Release notes from CHANGELOG.md
   - Compiled binaries for Linux, macOS, Windows
   - Checksums file
   - Source code archives

#### 5. Verify Release

1. Check GitHub Actions workflow completed successfully
2. Verify GitHub Release was created: https://github.com/jontk/slurm-client/releases
3. Download and test binaries for at least one platform
4. Verify `go get github.com/jontk/slurm-client@v0.3.0` works

#### 6. Announce Release

1. Create release announcement (optional)
2. Update documentation if needed
3. Notify users through appropriate channels

### Rollback Procedure

If a release has critical issues:

1. **Delete the tag locally and remotely**:
   ```bash
   git tag -d v0.2.0
   git push origin :refs/tags/v0.2.0
   ```

2. **Delete the GitHub Release**:
   - Go to GitHub Releases page
   - Click "Delete" on the problematic release

3. **Fix the issue**:
   - Create a hotfix branch
   - Fix the critical issue
   - Test thoroughly

4. **Create a new patch release**:
   - Follow release process for v0.2.1

## Release Checklist

Use this checklist for each release:

- [ ] All tests passing on main
- [ ] Security scans completed (gosec, govulncheck)
- [ ] CHANGELOG.md updated with all changes
- [ ] Version bumped appropriately (major/minor/patch)
- [ ] Documentation updated if needed
- [ ] No open critical bugs
- [ ] Tag created and pushed
- [ ] GitHub Actions workflow completed successfully
- [ ] GitHub Release created with all artifacts
- [ ] Binaries downloaded and smoke tested
- [ ] `go get` with new version works
- [ ] Release announcement prepared (if needed)

## Hotfix Process

For critical bugs that need immediate release:

1. **Create hotfix branch from tag**:
   ```bash
   git checkout -b hotfix/v0.2.1 v0.2.0
   ```

2. **Fix the bug**:
   ```bash
   # Make changes
   git add .
   git commit -m "fix: critical bug description"
   ```

3. **Update CHANGELOG.md**:
   ```markdown
   ## [0.2.1] - 2026-01-26

   ### Fixed
   - Critical bug that caused X
   ```

4. **Merge to main**:
   ```bash
   git checkout main
   git merge hotfix/v0.2.1
   git push origin main
   ```

5. **Create tag**:
   ```bash
   git tag -a v0.2.1 -m "Hotfix v0.2.1: Fix critical bug"
   git push origin v0.2.1
   ```

6. **Delete hotfix branch**:
   ```bash
   git branch -d hotfix/v0.2.1
   ```

## Version Support Policy

- **Latest minor version**: Full support with bug fixes and security updates
- **Previous minor version**: Security updates only (for 3 months after new minor release)
- **Older versions**: No support (users should upgrade)

## Emergency Security Release

For security vulnerabilities:

1. Follow the security policy in SECURITY.md
2. Create patch in private
3. Prepare security advisory
4. Release patch version
5. Publish security advisory
6. Notify affected users

## Automated Changelog

We use conventional commits to automatically categorize changes:

- `feat:` → New Features
- `fix:` → Bug Fixes
- `perf:` → Performance Improvements
- `sec:` → Security Updates
- `docs:` → Documentation (excluded from changelog)
- `test:` → Tests (excluded from changelog)
- `chore:` → Chores (excluded from changelog)

Example commit messages:
```bash
feat: add support for job array management
fix: resolve nil pointer dereference in partition manager
perf: optimize job list queries with caching
sec: update dependencies to address CVE-2025-XXXXX
```

## Release Artifacts

Each release includes:

1. **Source Code Archives**:
   - `.tar.gz` and `.zip` with source code

2. **Binary Archives**:
   - `slurm-client_{version}_Linux_x86_64.tar.gz`
   - `slurm-client_{version}_Linux_arm64.tar.gz`
   - `slurm-client_{version}_Darwin_x86_64.tar.gz`
   - `slurm-client_{version}_Darwin_arm64.tar.gz`
   - `slurm-client_{version}_Windows_x86_64.zip`

3. **Checksums**:
   - `checksums.txt` with SHA256 sums

4. **Release Notes**:
   - Extracted from CHANGELOG.md
   - Includes upgrade instructions if needed

## Testing Releases Locally

Test the release process without publishing:

```bash
# Install GoReleaser
go install github.com/goreleaser/goreleaser@latest

# Build snapshot (doesn't publish)
goreleaser release --snapshot --clean

# Check dist/ directory for artifacts
ls -la dist/
```

## CI/CD Pipeline

Our release pipeline:

1. **Trigger**: Git tag pushed
2. **Checkout**: Code at tag version
3. **Setup**: Go environment, dependencies
4. **Test**: Run full test suite
5. **Security**: Run security scans
6. **Build**: GoReleaser builds all platforms
7. **Publish**: Create GitHub Release
8. **Verify**: Smoke tests on artifacts

## FAQ

### Q: What if the release workflow fails?

A: Delete the tag, fix the issue, and recreate the tag:
```bash
git tag -d v0.2.0
git push origin :refs/tags/v0.2.0
# Fix the issue
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
```

### Q: Can I create a release from a branch other than main?

A: Not recommended. All releases should come from main to ensure proper testing and review.

### Q: How do I create a pre-release?

A: Use pre-release version format:
```bash
git tag -a v0.2.0-rc.1 -m "Release candidate v0.2.0-rc.1"
git push origin v0.2.0-rc.1
```
GoReleaser will automatically mark it as a pre-release.

### Q: What if I need to update release notes after publishing?

A: Edit the GitHub Release directly through the web interface.

## Resources

- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [GoReleaser Documentation](https://goreleaser.com/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github)
