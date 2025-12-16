.PHONY: all build test clean install run-daemon run-cli help

# Build variables
DAEMON_BINARY=ripleyd
CLI_BINARY=ripleyctl
BUILD_DIR=.
MAIN_PACKAGE=.
CLI_PACKAGE=./cmd/ripleyctl

all: test build

# Build both daemon and CLI
build:
	@echo "Building Ripley daemon..."
	go build -o $(BUILD_DIR)/$(DAEMON_BINARY) $(MAIN_PACKAGE)
	@echo "Building ripleyctl CLI..."
	go build -o $(BUILD_DIR)/$(CLI_BINARY) $(CLI_PACKAGE)
	@echo "Build complete!"

# Run all tests
test:
	@echo "Running tests..."
	go test -v ./...
	@echo "Tests complete!"

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(DAEMON_BINARY) $(CLI_BINARY)
	rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Install binaries to $GOPATH/bin
install:
	@echo "Installing binaries..."
	go install $(MAIN_PACKAGE)
	go install $(CLI_PACKAGE)
	@echo "Install complete!"

# Run the daemon (requires Claude CLI)
run-daemon: build
	@echo "Starting Ripley daemon..."
	./$(DAEMON_BINARY)

# Run the CLI tool (requires Claude CLI)
run-cli: build
	@echo "Starting ripleyctl CLI..."
	./$(CLI_BINARY)

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Format complete!"

# Run linter
lint:
	@echo "Running linter..."
	go vet ./...
	@echo "Lint complete!"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	@echo "Dependencies downloaded!"

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy
	@echo "Dependencies tidied!"

# Help
help:
	@echo "Ripley Daemon - Makefile targets:"
	@echo ""
	@echo "  make build          - Build daemon and CLI binaries"
	@echo "  make test           - Run all tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make install        - Install binaries to \$$GOPATH/bin"
	@echo "  make run-daemon     - Build and run the daemon"
	@echo "  make run-cli        - Build and run the CLI tool"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linter"
	@echo "  make deps           - Download dependencies"
	@echo "  make tidy           - Tidy dependencies"
	@echo "  make help           - Show this help message"
	@echo ""
