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
Follow [Conventional Commits](https://www.conventionalcommits.org/):
```bash
git commit -m "feat: add new job scheduling feature"
git commit -m "fix: resolve connection timeout issue"
git commit -m "docs: update API documentation"
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

## Code Style Guidelines

### Go Code
- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `golangci-lint` before submitting
- Add comments for exported functions

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