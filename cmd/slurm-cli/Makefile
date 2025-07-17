.PHONY: build install clean test

# Binary name
BINARY=slurm-cli

# Build directory
BUILD_DIR=../../build

# Version information
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.Commit=${COMMIT}"

# Default target
all: build

# Build the binary
build:
	@echo "Building ${BINARY}..."
	@mkdir -p ${BUILD_DIR}
	go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY} .
	@echo "Binary built at ${BUILD_DIR}/${BINARY}"

# Install the binary
install:
	@echo "Installing ${BINARY}..."
	go install ${LDFLAGS} .

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f ${BUILD_DIR}/${BINARY}
	@go clean

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Build for multiple platforms
build-all: build-linux build-darwin build-windows

build-linux:
	@echo "Building for Linux..."
	@mkdir -p ${BUILD_DIR}
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}-linux-arm64 .

build-darwin:
	@echo "Building for macOS..."
	@mkdir -p ${BUILD_DIR}
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}-darwin-arm64 .

build-windows:
	@echo "Building for Windows..."
	@mkdir -p ${BUILD_DIR}
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}-windows-amd64.exe .

# Run the CLI
run:
	go run . $(ARGS)

# Get dependencies
deps:
	go mod download
	go mod tidy

# Update dependencies
update:
	go get -u ./...
	go mod tidy

# Show help
help:
	@echo "Available targets:"
	@echo "  build       - Build the binary for current platform"
	@echo "  install     - Install the binary to GOPATH/bin"
	@echo "  clean       - Remove build artifacts"
	@echo "  test        - Run tests"
	@echo "  build-all   - Build for all platforms"
	@echo "  deps        - Download dependencies"
	@echo "  update      - Update dependencies"
	@echo "  run         - Run the CLI (use ARGS to pass arguments)"
	@echo "  help        - Show this help message"