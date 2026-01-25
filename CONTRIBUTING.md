# Contributing to SLURM REST API Client Library

Thank you for your interest in contributing to the SLURM REST API Client Library! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to abide by our code of conduct:
- Be respectful and inclusive
- Welcome newcomers and help them get started
- Focus on constructive criticism
- Respect differing viewpoints and experiences

## How to Contribute

### Reporting Issues

1. Check if the issue already exists in [GitHub Issues](https://github.com/jontk/slurm-client/issues)
2. If not, create a new issue with:
   - Clear title and description
   - Steps to reproduce (if applicable)
   - Expected vs actual behavior
   - Environment details (SLURM version, Go version, OS)

### Suggesting Features

1. Open a [GitHub Issue](https://github.com/jontk/slurm-client/issues/new) with the "enhancement" label
2. Describe the feature and its use case
3. Explain why it would be valuable to other users

### Submitting Code Changes

#### 1. Fork and Clone
```bash
# Fork the repository on GitHub, then:
git clone https://github.com/YOUR-USERNAME/slurm-client.git
cd slurm-client
git remote add upstream https://github.com/jontk/slurm-client.git
```

#### 2. Create a Branch
```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-description
```

#### 3. Make Changes
- Follow the existing code style
- Add tests for new functionality
- Update documentation as needed
- Ensure all tests pass

#### 4. Commit Changes

**⚠️ IMPORTANT: All commits MUST follow [Conventional Commits](https://www.conventionalcommits.org/) format.**

Conventional commits are **strictly enforced** by automated checks. Non-compliant commits will cause PR failures.

**Format:**
```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

**Required Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, missing semicolons, etc.)
- `refactor`: Code refactoring (no functional changes)
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `build`: Build system changes
- `ci`: CI/CD changes
- `chore`: Other changes (maintenance, dependencies, etc.)
- `revert`: Reverting a previous commit

**Scope (optional but recommended):**
- Use lowercase
- Examples: `api`, `client`, `test`, `docs`, `release`, `contrib`

**Subject:**
- Use imperative mood ("add" not "added" or "adds")
- Don't capitalize first letter
- No period at the end
- Max 100 characters for entire header

**Examples:**
```bash
git commit -m "feat(api): add new job scheduling endpoint"
git commit -m "fix(client): resolve connection timeout issue"
git commit -m "docs: update API documentation"
git commit -m "test(streaming): add WebSocket tests"
git commit -m "refactor(auth): simplify token validation"
```

**With body:**
```bash
git commit -m "feat(api): add batch job submission

Allow submitting multiple jobs in a single API call.
Improves performance by reducing HTTP overhead.

Closes #123"
```

**Breaking changes:**
```bash
git commit -m "feat(api)!: change job submission response format

BREAKING CHANGE: Response now returns job array instead of single job"
```

**Automated Validation:**
- Pre-commit hooks validate commit messages locally
- GitHub Actions validate all PR commits
- Non-compliant commits will fail CI checks

**Install hooks:**
```bash
make install-hooks  # Enables local validation
```

#### 5. Push and Create PR
```bash
git push origin feature/your-feature-name
```
Then create a Pull Request on GitHub.

## Development Setup

### Prerequisites
- Go 1.20 or later
- Access to a SLURM cluster (for integration testing)
- Git

### Building
```bash
go build ./...
```

### Testing
```bash
# Unit tests
go test ./...

# Integration tests (requires SLURM cluster)
go test -tags=integration ./tests/integration/...

# With coverage
go test -coverprofile=coverage.txt -covermode=atomic ./...
```

### Code Generation
When modifying OpenAPI specs:
```bash
# Run code generation
go generate ./...

# Verify generated files
git status
```

### Performance Testing

#### Running Benchmarks
```bash
# Run all benchmarks
go test -bench=. -benchmem ./...

# Run specific package benchmarks
go test -bench=. -benchmem ./pkg/streaming/...

# Run with more iterations for accuracy
go test -bench=. -benchmem -benchtime=5s -count=10 ./...
```

#### Comparing Performance Locally
To detect performance regressions, compare benchmarks between branches:

```bash
# Install benchstat (if not already installed)
go install golang.org/x/perf/cmd/benchstat@latest

# Run benchmarks on current branch
go test -bench=. -benchmem -benchtime=3s -count=5 ./... > new.txt

# Switch to base branch
git checkout main

# Run benchmarks on base branch
go test -bench=. -benchmem -benchtime=3s -count=5 ./... > old.txt

# Compare results
benchstat old.txt new.txt
```

**Note:** Benchmark results can vary based on system load. For reliable comparisons:
- Close unnecessary applications
- Run multiple iterations (count=5 or higher)
- Use consistent benchtime values
- Consider running on a dedicated/idle machine

Example benchstat output:
```
name                     old time/op    new time/op    delta
ParseNumberField-8         45.2ns ± 2%    46.1ns ± 1%   +2.00%  (p=0.043 n=5+5)
FormatDurationForSlurm-8    156ns ± 1%     158ns ± 2%     ~     (p=0.222 n=5+5)
```

Interpret results:
- `~` means no significant difference
- Percentages show the performance change
- `p` values indicate statistical significance (p<0.05 is significant)

### Pre-commit Hooks

We use pre-commit hooks to ensure code quality before commits. These hooks automatically run formatters and linters on your code.

#### Installing Pre-commit Hooks
```bash
# Install pre-commit hooks (one-time setup)
make install-hooks

# Or manually:
pip install pre-commit  # or: brew install pre-commit
pre-commit install
```

#### What the Hooks Do
The pre-commit hooks will automatically:
- Format Go code with `gofmt` and `goimports`
- Run `golangci-lint` with auto-fix enabled
- Tidy `go.mod` dependencies
- Remove trailing whitespace
- Fix end-of-file issues
- Check YAML syntax
- Detect merge conflicts

#### Running Hooks Manually
```bash
# Run on all files
pre-commit run --all-files

# Run on staged files only
pre-commit run

# Skip hooks for a specific commit (use sparingly)
git commit --no-verify -m "message"
```

#### Uninstalling Hooks
```bash
make uninstall-hooks
# or: pre-commit uninstall
```

**Note:** Pre-commit hooks are optional but highly recommended. They catch issues early and ensure consistent code quality.

## Code Style Guidelines

### Go Code
- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `golangci-lint` before submitting
- Add comments for exported functions

### Linting and Code Quality

For comprehensive guidance on linting, configuration, and the code review process, see [docs/LINTING.md](./docs/LINTING.md). This includes:

- How to run linters locally
- Understanding disabled linters and project philosophy
- Proposing linter changes
- CI/CD linting integration

**Quick reference:**
```bash
make lint              # Run full lint check
make lint-staged       # Check only staged changes
make fmt               # Auto-format code
```

### File Headers
All Go files must include SPDX headers:
```go
// SPDX-FileCopyrightText: 2025 Your Name
// SPDX-License-Identifier: Apache-2.0
```

### Error Handling
- Use the custom error types in `pkg/errors`
- Provide context in error messages
- Wrap errors appropriately

### Testing
- Write unit tests for new functionality
- Aim for >80% code coverage
- Use table-driven tests where appropriate
- Mock external dependencies

## Pull Request Process

1. **Before Submitting:**
   - Rebase on latest main branch
   - Ensure all tests pass
   - Run linters and fix issues
   - Update documentation

2. **PR Description:**
   - Reference any related issues
   - Describe what changes were made
   - Explain why the changes are needed
   - List any breaking changes

3. **Review Process:**
   - PRs require at least one approval
   - Address reviewer feedback
   - Keep PR focused and reasonably sized

## Release Process

Releases are managed by maintainers:
1. Version tags follow semantic versioning (v1.2.3)
2. Release notes are automatically generated
3. Packages are published to pkg.go.dev

## Getting Help

- Check [documentation](./docs)
- Review [examples](./examples)
- Ask questions in [GitHub Issues](https://github.com/jontk/slurm-client/issues)
- Review existing [Pull Requests](https://github.com/jontk/slurm-client/pulls)

## Recognition

Contributors are recognized in:
- Release notes
- README contributors section
- GitHub contributors page

Thank you for contributing to make this project better!