.PHONY: build test clean run docker-build docker-run lint fmt help

# Variables
BINARY_NAME=htb-mcp-server
DOCKER_IMAGE=htb-mcp-server
VERSION?=1.0.0
BUILD_DIR=build
GO_VERSION=1.21

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -ldflags="-w -s -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go
	@echo "Multi-platform build complete"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run integration tests (requires HTB_TOKEN)
test-integration:
	@echo "Running integration tests..."
	@if [ -z "$(HTB_TOKEN)" ]; then echo "HTB_TOKEN environment variable required for integration tests"; exit 1; fi
	go test -v -tags=integration ./...

# Run the binary
run: build
	@echo "Running $(BINARY_NAME)..."
	@if [ -z "$(HTB_TOKEN)" ]; then echo "HTB_TOKEN environment variable required"; exit 1; fi
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run with debug logging
run-debug: build
	@echo "Running $(BINARY_NAME) in debug mode..."
	@if [ -z "$(HTB_TOKEN)" ]; then echo "HTB_TOKEN environment variable required"; exit 1; fi
	LOG_LEVEL=DEBUG ./$(BUILD_DIR)/$(BINARY_NAME)

# Build Docker image
docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE):$(VERSION)..."
	docker build -t $(DOCKER_IMAGE):$(VERSION) .
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest
	@echo "Docker image built: $(DOCKER_IMAGE):$(VERSION)"

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	@if [ -z "$(HTB_TOKEN)" ]; then echo "HTB_TOKEN environment variable required"; exit 1; fi
	docker run --rm -e HTB_TOKEN="$(HTB_TOKEN)" $(DOCKER_IMAGE):latest

# Lint the code
lint:
	@echo "Running linter..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run

# Format the code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@if command -v goimports &> /dev/null; then \
		goimports -w .; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	docker rmi -f $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest 2>/dev/null || true
	@echo "Clean complete"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Security audit
audit:
	@echo "Running security audit..."
	@if ! command -v gosec &> /dev/null; then \
		echo "gosec not found. Installing..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
	fi
	gosec ./...

# Generate documentation
docs:
	@echo "Generating documentation..."
	@if ! command -v godoc &> /dev/null; then \
		echo "godoc not found. Installing..."; \
		go install golang.org/x/tools/cmd/godoc@latest; \
	fi
	@echo "Documentation server available at: http://localhost:6060"
	godoc -http=:6060

# Create release
release: clean lint test build-all
	@echo "Creating release $(VERSION)..."
	@mkdir -p release
	@cd $(BUILD_DIR) && tar -czf ../release/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	@cd $(BUILD_DIR) && tar -czf ../release/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	@cd $(BUILD_DIR) && tar -czf ../release/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	@cd $(BUILD_DIR) && zip -r ../release/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	@echo "Release created in release/ directory"

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install golang.org/x/tools/cmd/godoc@latest
	@echo "Development tools installed"

# Benchmark tests
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Profile the application
profile: build
	@echo "Running with profiling enabled..."
	@if [ -z "$(HTB_TOKEN)" ]; then echo "HTB_TOKEN environment variable required"; exit 1; fi
	go tool pprof -http=:8080 $(BUILD_DIR)/$(BINARY_NAME)

# Check for updates
update:
	@echo "Checking for dependency updates..."
	go list -u -m all

# Validate project structure
validate:
	@echo "Validating project structure..."
	@echo "✓ Checking main.go exists..."
	@test -f main.go || (echo "✗ main.go not found" && exit 1)
	@echo "✓ Checking go.mod exists..."
	@test -f go.mod || (echo "✗ go.mod not found" && exit 1)
	@echo "✓ Checking package structure..."
	@test -d pkg || (echo "✗ pkg directory not found" && exit 1)
	@test -d internal || (echo "✗ internal directory not found" && exit 1)
	@echo "✓ Project structure validation passed"

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the binary"
	@echo "  build-all      - Build for multiple platforms"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  test-integration - Run integration tests (requires HTB_TOKEN)"
	@echo "  run            - Build and run the binary"
	@echo "  run-debug      - Run with debug logging"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Install dependencies"
	@echo "  audit          - Run security audit"
	@echo "  docs           - Generate documentation"
	@echo "  release        - Create release packages"
	@echo "  dev-setup      - Setup development environment"
	@echo "  benchmark      - Run benchmark tests"
	@echo "  profile        - Run with profiling"
	@echo "  update         - Check for dependency updates"
	@echo "  validate       - Validate project structure"
	@echo "  help           - Show this help"

# Default target
all: clean lint test build

# Set default target
.DEFAULT_GOAL := help