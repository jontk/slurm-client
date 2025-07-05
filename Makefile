# Makefile for slurm-client

.PHONY: build test lint fmt vet clean docs help install-tools

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

# Run tests
test:
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

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	}

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

# Show help
help:
	@echo "Available targets:"
	@echo "  build           - Build the library"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  benchmark       - Run benchmarks"
	@echo "  lint            - Run linter"
	@echo "  fmt             - Format code"
	@echo "  vet             - Run go vet"
	@echo "  clean           - Clean build artifacts"
	@echo "  docs            - Generate documentation"
	@echo "  install-tools   - Install development tools"
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