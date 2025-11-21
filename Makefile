.PHONY: build install test clean help

# Variables
BINARY_NAME=rc
BUILD_DIR=.
INSTALL_DIR=/usr/local/bin

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/rc
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Install the binary to system
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@echo "Installation complete!"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	@go test -cover ./...
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Run the application
run: build
	@./$(BINARY_NAME)

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run || echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"

# Show help
help:
	@echo "GoCycled Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build        Build the application"
	@echo "  make install      Install to $(INSTALL_DIR)"
	@echo "  make test         Run tests"
	@echo "  make test-cover   Run tests with coverage report"
	@echo "  make clean        Remove build artifacts"
	@echo "  make run          Build and run the application"
	@echo "  make fmt          Format code"
	@echo "  make lint         Run linter"
	@echo "  make help         Show this help message"
