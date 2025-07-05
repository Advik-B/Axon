# Makefile for the Axon Language Project

# --- Variables ---
# Use 'axon.exe' on Windows, 'axon' otherwise
ifeq ($(OS),Windows_NT)
    BINARY_NAME=axon.exe
else
    BINARY_NAME=axon
endif

# Path to the main CLI package
CMD_PATH=./cmd/axon

# Path to the Protobuf source file
PROTO_SRC=pkg/axon/axon.proto

# Output directory for transpiled Go code
OUTPUT_DIR=./out

# --- Targets ---

.PHONY: all build run proto test clean install-deps help

# Default target executed when you just run 'make'
all: build

# Build the Axon CLI binary
build:
	@echo "Building Axon CLI..."
	go build -v -o $(BINARY_NAME) $(CMD_PATH)
	@echo "✅ Build complete: $(BINARY_NAME)"

# Run the Axon CLI with arguments passed via the ARGS variable
# Example: make run ARGS="build examples/stdlib_example.ax"
run:
	go run $(CMD_PATH) $(ARGS)

# Generate Go code from the .proto definition file
proto:
	@echo "Generating Go code from Protobuf definition..."
	@protoc --go_out=. --go_opt=paths=source_relative $(PROTO_SRC)
	@echo "✅ Protobuf generation complete."

# Run all tests in the project
test:
	@echo "Running tests..."
	go test -v ./...

# Install necessary Go tools and tidy the go.mod file
install-deps:
	@echo "Installing Go tools and dependencies..."
	@go mod tidy
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Clean up all build artifacts and generated files
clean:
	@echo "Cleaning up build artifacts and generated directories..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(OUTPUT_DIR)
	@echo "🗑️ Clean complete."

# Display help information about the available commands
help:
	@echo "Axon Project Makefile"
	@echo "---------------------"
	@echo "Available commands:"
	@echo "  make build          - Compile the Axon CLI."
	@echo "  make run            - Run the CLI. Pass arguments with ARGS=\"...\" (e.g., make run ARGS=\"build examples/stdlib_example.ax\")."
	@echo "  make proto          - Regenerate Go code from the .proto file."
	@echo "  make test           - Run all Go tests."
	@echo "  make clean          - Remove the compiled binary and the '/out' directory."
	@echo "  make install-deps   - Install Go tools needed for development."
	@echo "  make help           - Show this help message."