# Makefile for slurm-client

.PHONY: build test check-mocks lint lint-staged fmt vet clean docs help install-tools install-hooks generate download-specs generate-mocks

# Variables
BINARY_NAME=slurm-client
PACKAGE=github.com/jontk/slurm-client
GO_VERSION=1.21

# Default target
all: lint test build

# Build the library
build:
	@echo "Building..."
	go build -v ./...

# Check if mock builders are generated
check-mocks:
	@if [ ! -d "tests/mocks/generated/v0_0_40" ] || \
	    [ ! -d "tests/mocks/generated/v0_0_42" ] || \
	    [ ! -d "tests/mocks/generated/v0_0_43" ] || \
	    [ ! -d "tests/mocks/generated/v0_0_44" ]; then \
		echo "❌ Error: Mock builders not found!"; \
		echo ""; \
		echo "Mock builders are required for tests but not committed to git."; \
		echo "Please generate them first:"; \
		echo ""; \
		echo "    make generate-mocks"; \
		echo ""; \
		echo "This only needs to be done once after cloning the repository."; \
		exit 1; \
	fi
	@echo "✅ Mock builders found"

# Run tests
test: check-mocks
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
test-coverage: test
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Lint the code
lint: install-tools
	@echo "Running linter..."
	golangci-lint run

# Lint only staged files (for pre-commit hooks)
lint-staged: install-tools
	@echo "Linting staged files..."
	@STAGED_GO_FILES=$$(git diff --cached --name-only --diff-filter=ACM -- '*.go' 2>/dev/null); \
	if [ -z "$$STAGED_GO_FILES" ]; then \
		echo "No staged Go files to lint"; \
	else \
		echo "Staged files: $$STAGED_GO_FILES"; \
		golangci-lint run $$STAGED_GO_FILES || exit 1; \
	fi

# Format the code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	go clean
	rm -f coverage.out coverage.html

# Generate documentation
docs:
	@echo "Generating documentation..."
	go doc -all > docs.txt
	@echo "Documentation generated: docs.txt"

# Install oapi-codegen (needed for code generation)
install-oapi-codegen:
	@command -v oapi-codegen >/dev/null 2>&1 || { \
		echo "Installing oapi-codegen..."; \
		go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest; \
	}

# Install development tools
install-tools: install-oapi-codegen
	@echo "Installing development tools..."
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	}

# Install pre-commit hooks
install-hooks:
	@echo "Installing pre-commit hooks..."
	@command -v pre-commit >/dev/null 2>&1 || { \
		echo "pre-commit not found. Installing..."; \
		if command -v pip3 >/dev/null 2>&1; then \
			pip3 install --user pre-commit; \
		elif command -v pip >/dev/null 2>&1; then \
			pip install --user pre-commit; \
		elif command -v brew >/dev/null 2>&1; then \
			brew install pre-commit; \
		else \
			echo "Error: Could not install pre-commit. Please install manually:"; \
			echo "  pip install pre-commit"; \
			echo "  or: brew install pre-commit"; \
			exit 1; \
		fi; \
	}
	@pre-commit install
	@pre-commit install --hook-type commit-msg
	@echo "✅ Pre-commit hooks installed successfully"
	@echo "   - Pre-commit hooks: code formatting, linting"
	@echo "   - Commit-msg hooks: conventional commit validation"
	@echo "Run 'pre-commit run --all-files' to check all files"

# Uninstall pre-commit hooks
uninstall-hooks:
	@echo "Uninstalling pre-commit hooks..."
	@pre-commit uninstall || true
	@echo "✅ Pre-commit hooks uninstalled"

# Tidy up dependencies
tidy:
	@echo "Tidying up dependencies..."
	go mod tidy

# Update dependencies
update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Security audit
security:
	@echo "Running security audit..."
	@command -v gosec >/dev/null 2>&1 || { \
		echo "Installing gosec..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
	}
	gosec ./...

# Run all checks (lint, vet, test)
check: fmt vet lint test

# Install the module locally
install:
	@echo "Installing module..."
	go install ./...

# Run examples
examples:
	@echo "Running examples..."
	@for example in examples/*/; do \
		echo "Running example: $$example"; \
		(cd "$$example" && go run .); \
	done

# Generate mock files for testing
mocks:
	@echo "Generating mocks..."
	@command -v mockgen >/dev/null 2>&1 || { \
		echo "Installing mockgen..."; \
		go install github.com/golang/mock/mockgen@latest; \
	}
	mockgen -source=client.go -destination=mocks/client_mock.go

# Run integration tests (requires Slurm server)
integration-test:
	@echo "Running integration tests..."
	@if [ -z "$$SLURM_REST_URL" ]; then \
		echo "SLURM_REST_URL environment variable is required for integration tests"; \
		exit 1; \
	fi
	go test -tags=integration -v ./tests/integration/...

# Check for Go version
check-version:
	@echo "Checking Go version..."
	@go version | grep -q "go1\\.2[1-9]" || { \
		echo "Go $(GO_VERSION) or higher is required"; \
		exit 1; \
	}

# Release preparation
release-prep: check-version clean tidy fmt vet lint test
	@echo "Release preparation complete"

# Build SLURM with v0.0.44 support
build-slurm:
	@echo "Building SLURM with REST API v0.0.44 support..."
	@./tools/build-slurm.sh

# Download OpenAPI specifications
download-specs:
	@echo "Downloading OpenAPI specifications..."
	@./tools/codegen/download-specs.sh

# Generate API clients from OpenAPI specs
generate: install-oapi-codegen download-specs
	@echo "Generating API clients..."
	@for version in v0.0.40 v0.0.41 v0.0.42 v0.0.43 v0.0.44; do \
		echo "Generating client for $$version..."; \
		go run tools/codegen/generate.go $$version || echo "Failed to generate $$version client"; \
	done

# Generate mock builders from OpenAPI specs
generate-mocks: install-oapi-codegen
	@echo "Generating mock builders..."
	@for version in v0.0.40 v0.0.41 v0.0.42 v0.0.43 v0.0.44; do \
		echo "Generating mock builders for $$version..."; \
		go run tools/codegen/generate_mocks.go $$version || echo "Failed to generate $$version mocks"; \
	done

# Generate both API clients and mock builders
generate-all: generate generate-mocks
	@echo "All code generation complete"

# Generate specific version client
generate-version: install-tools
	@if [ -z "$(VERSION)" ]; then \
		echo "Usage: make generate-version VERSION=v0.0.42"; \
		exit 1; \
	fi
	@echo "Generating client for $(VERSION)..."
	@go run tools/codegen/generate.go $(VERSION)

# Show help
help:
	@echo "Available targets:"
	@echo "  build           - Build the library"
	@echo "  test            - Run tests (auto-checks for mock builders)"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  check-mocks     - Verify mock builders are generated"
	@echo "  benchmark       - Run benchmarks"
	@echo "  lint            - Run linter"
	@echo "  fmt             - Format code"
	@echo "  vet             - Run go vet"
	@echo "  clean           - Clean build artifacts"
	@echo "  docs            - Generate documentation"
	@echo "  install-tools   - Install development tools"
	@echo "  install-hooks   - Install pre-commit hooks"
	@echo "  uninstall-hooks - Uninstall pre-commit hooks"
	@echo "  build-slurm     - Build SLURM with REST API v0.0.44 support"
	@echo "  download-specs  - Download OpenAPI specifications"
	@echo "  generate        - Generate all API clients from specs"
	@echo "  generate-mocks  - Generate mock builders from specs"
	@echo "  generate-all    - Generate both API clients and mock builders"
	@echo "  generate-version - Generate specific version client (VERSION=v0.0.44)"
	@echo "  tidy            - Tidy up dependencies"
	@echo "  update          - Update dependencies"
	@echo "  security        - Run security audit"
	@echo "  check           - Run all checks (fmt, vet, lint, test)"
	@echo "  install         - Install the module locally"
	@echo "  examples        - Run examples"
	@echo "  mocks           - Generate mock files"
	@echo "  integration-test - Run integration tests (requires SLURM_REST_URL)"
	@echo "  check-version   - Check Go version"
	@echo "  release-prep    - Prepare for release"
	@echo "  help            - Show this help message"