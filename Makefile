# Makefile for the Go project

# Define the Go binary name
BINARY_NAME=gotcc

# Define the source files
SRC=$(shell find ./internal -name '*.go') $(shell find ./cmd -name '*.go')

# Define the output directory for binaries
OUTPUT_DIR=bin

# Define the Go module
GO_MODULE=$(shell go list -m)

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building the application..."
	@mkdir -p $(OUTPUT_DIR)
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME) $(SRC)

# Run the application
.PHONY: run
run: build
	@echo "Running the application..."
	./$(OUTPUT_DIR)/$(BINARY_NAME)

# Clean the build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf $(OUTPUT_DIR)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# Run database migrations
.PHONY: migrate
migrate:
	@echo "Running database migrations..."
	./scripts/migrate.sh

# Show help
.PHONY: help
help:
	@echo "Makefile commands:"
	@echo "  all       - Build the application"
	@echo "  build     - Build the application"
	@echo "  run       - Run the application"
	@echo "  clean     - Clean build artifacts"
	@echo "  test      - Run tests"
	@echo "  migrate   - Run database migrations"
	@echo "  help      - Show this help message"