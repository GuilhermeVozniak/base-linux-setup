# Base Linux Setup Makefile

# Application name
APP_NAME = base-linux-setup

# Build directory
BUILD_DIR = build

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod

# Binary name
BINARY_NAME = $(APP_NAME)
BINARY_PATH = $(BUILD_DIR)/$(BINARY_NAME)

# Installation directory
INSTALL_DIR = /usr/local/bin

# Version information
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Linker flags
LDFLAGS = -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.commit=$(COMMIT)"

.PHONY: all build clean test deps help install uninstall run dev

# Default target
all: build

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) .
	@echo "Build completed: $(BINARY_PATH)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "Clean completed"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Install the application
install: build
	@echo "Installing $(APP_NAME) to $(INSTALL_DIR)..."
	@sudo cp $(BINARY_PATH) $(INSTALL_DIR)/$(BINARY_NAME)
	@sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installation completed"

# Uninstall the application
uninstall:
	@echo "Uninstalling $(APP_NAME)..."
	@sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Uninstallation completed"

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	@$(BINARY_PATH)

# Development mode - build and run
dev: build
	@echo "Starting development mode..."
	@$(BINARY_PATH)

# Create release builds for multiple platforms
release:
	@echo "Creating release builds..."
	@mkdir -p $(BUILD_DIR)/release
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-amd64 .
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-arm64 .
	
	# Linux ARM
	GOOS=linux GOARCH=arm $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-arm .
	
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-amd64 .
	
	# macOS ARM64
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-arm64 .
	
	@echo "Release builds completed in $(BUILD_DIR)/release/"

# Check code quality
lint:
	@echo "Running linters..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1; }
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Update dependencies
update:
	@echo "Updating dependencies..."
	$(GOMOD) get -u
	$(GOMOD) tidy

# Show help
help:
	@echo "Available commands:"
	@echo "  build     - Build the application"
	@echo "  clean     - Clean build artifacts"
	@echo "  test      - Run tests"
	@echo "  deps      - Download dependencies"
	@echo "  install   - Install the application to $(INSTALL_DIR)"
	@echo "  uninstall - Remove the application from $(INSTALL_DIR)"
	@echo "  run       - Build and run the application"
	@echo "  dev       - Development mode (build and run)"
	@echo "  release   - Create release builds for multiple platforms"
	@echo "  lint      - Run code linters"
	@echo "  fmt       - Format code"
	@echo "  update    - Update dependencies"
	@echo "  help      - Show this help message"

# Show version information
version:
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Commit: $(COMMIT)" 